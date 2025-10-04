// internal/service/oauth_service.go

package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

// OAuthService 定义OAuth相关业务逻辑接口
type OAuthService interface {
	// Client相关操作
	CreateClient(ctx context.Context, clientID, clientSecret, name string, redirectURIs, scopes []string) (*model.Client, error)
	FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error)
	
	// AuthorizationCode相关操作
	GenerateAuthorizationCode(ctx context.Context, clientID string, userID int64, redirectURI string, scopes []string, codeChallenge *string, codeChallengeMethod *string) (*model.AuthorizationCode, error)
	ValidateAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*model.AuthorizationCode, error)
	ExchangeAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*ExchangeAuthorizationCodeResult, error)
	
	// Token相关操作
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	CreateRefreshToken(ctx context.Context, userID int64, clientID string, scopes []string) (string, error)
	RefreshAccessToken(ctx context.Context, refreshTokenString string) (*ExchangeAuthorizationCodeResult, error)
}

// oauthService 是 OAuthService 接口的实现
type oauthService struct {
	oauthRepo repository.OAuthRepository
	userRepo  repository.UserRepository
}

// NewOAuthService 创建一个新的 oauthService 实例
func NewOAuthService(oauthRepo repository.OAuthRepository, userRepo repository.UserRepository) OAuthService {
	return &oauthService{
		oauthRepo: oauthRepo,
		userRepo:  userRepo,
	}
}

// CreateClient 创建新的OAuth客户端
func (s *oauthService) CreateClient(ctx context.Context, clientID, clientSecret, name string, redirectURIs, scopes []string) (*model.Client, error) {
	// 对客户端密钥进行哈希处理
	clientSecretHash := hashSecret(clientSecret)
	
	client := &model.Client{
		ClientID:          clientID,
		ClientSecretHash:  clientSecretHash,
		Name:              name,
		RedirectURIs:      redirectURIs,
		Scopes:            scopes,
		CreatedAt:         time.Now(),
	}
	
	err := s.oauthRepo.CreateClient(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	
	return client, nil
}

// FindClientByClientID 根据ClientID查找客户端
func (s *oauthService) FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error) {
	log.Printf("OAuth Service: Finding client by ID: %s", clientID)
	client, err := s.oauthRepo.FindClientByClientID(ctx, clientID)
	if err != nil {
		log.Printf("OAuth Service: Client not found in repository: %s, error: %v", clientID, err)
		return nil, fmt.Errorf("client not found: %w", err)
	}
	
	log.Printf("OAuth Service: Found client - ID: %d, Name: %s", client.ID, client.Name)
	return client, nil
}

// GenerateAuthorizationCode 生成授权码
func (s *oauthService) GenerateAuthorizationCode(ctx context.Context, clientID string, userID int64, redirectURI string, scopes []string, codeChallenge *string, codeChallengeMethod *string) (*model.AuthorizationCode, error) {
	log.Printf("OAuth Service: Generating authorization code for client: %s, user: %d, redirectURI: %s", clientID, userID, redirectURI)
	
	// 验证客户端是否存在
	client, err := s.FindClientByClientID(ctx, clientID)
	if err != nil {
		log.Printf("OAuth Service: Invalid client: %s, error: %v", clientID, err)
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	
	// 验证重定向URI是否在允许列表中
	log.Printf("OAuth Service: Validating redirect URI: %s", redirectURI)
	log.Printf("OAuth Service: Allowed redirect URIs: %v", client.RedirectURIs)
	if !isValidRedirectURI(redirectURI, client.RedirectURIs) {
		log.Printf("OAuth Service: Invalid redirect URI: %s", redirectURI)
		return nil, fmt.Errorf("invalid redirect URI")
	}
	
	// 验证请求的scopes是否被客户端允许
	log.Printf("OAuth Service: Validating scopes: %v", scopes)
	log.Printf("OAuth Service: Allowed scopes: %v", client.Scopes)
	if !areScopesAllowed(scopes, client.Scopes) {
		log.Printf("OAuth Service: Invalid scopes: %v", scopes)
		return nil, fmt.Errorf("invalid scopes")
	}
	
	// 生成随机授权码
	code := generateRandomCode(64)
	log.Printf("OAuth Service: Generated authorization code: %s", code)
	
	// 创建授权码实体
	authCode := &model.AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              scopes,
		ExpiresAt:           time.Now().Add(10 * time.Minute), // 10分钟有效期
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	
	// 保存到数据库
	log.Printf("OAuth Service: Saving authorization code to database")
	err = s.oauthRepo.CreateAuthorizationCode(ctx, authCode)
	if err != nil {
		log.Printf("OAuth Service: Failed to save authorization code: %v", err)
		return nil, fmt.Errorf("failed to save authorization code: %w", err)
	}
	
	log.Printf("OAuth Service: Authorization code saved successfully")
	return authCode, nil
}

// ValidateAuthorizationCode 验证授权码
func (s *oauthService) ValidateAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*model.AuthorizationCode, error) {
	log.Printf("OAuth Service: Validating authorization code: %s for client: %s, redirectURI: %s", code, clientID, redirectURI)
	
	// 查找授权码
	authCode, err := s.oauthRepo.FindAuthorizationCode(ctx, code)
	if err != nil {
		log.Printf("OAuth Service: Authorization code not found: %s, error: %v", code, err)
		return nil, fmt.Errorf("invalid authorization code")
	}
	
	log.Printf("OAuth Service: Found authorization code for client: %s, user: %d", authCode.ClientID, authCode.UserID)
	
	// 检查是否过期
	if time.Now().After(authCode.ExpiresAt) {
		log.Printf("OAuth Service: Authorization code expired: %s, expired at: %v, current time: %v", code, authCode.ExpiresAt, time.Now())
		return nil, fmt.Errorf("authorization code expired")
	}
	
	// 验证客户端ID
	if authCode.ClientID != clientID {
		log.Printf("OAuth Service: Client ID mismatch. Expected: %s, Got: %s", authCode.ClientID, clientID)
		return nil, fmt.Errorf("invalid client")
	}
	
	// 验证重定向URI
	if authCode.RedirectURI != redirectURI {
		log.Printf("OAuth Service: Redirect URI mismatch. Expected: %s, Got: %s", authCode.RedirectURI, redirectURI)
		return nil, fmt.Errorf("invalid redirect URI")
	}
	
	log.Printf("OAuth Service: Authorization code validated successfully")
	return authCode, nil
}

// ExchangeAuthorizationCodeResult 兑换授权码的结果
type ExchangeAuthorizationCodeResult struct {
	AccessToken  string
	RefreshToken string
}

// ExchangeAuthorizationCode 兑换授权码获取访问令牌和刷新令牌
func (s *oauthService) ExchangeAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*ExchangeAuthorizationCodeResult, error) {
	log.Printf("OAuth Service: Exchanging authorization code: %s for client: %s, redirectURI: %s", code, clientID, redirectURI)
	
	// 验证授权码
	authCode, err := s.ValidateAuthorizationCode(ctx, code, clientID, redirectURI)
	if err != nil {
		log.Printf("OAuth Service: Failed to validate authorization code: %v", err)
		return nil, err
	}
	
	// 生成访问令牌
	log.Printf("OAuth Service: Generating access token for user: %d", authCode.UserID)
	accessToken, err := s.GenerateAccessToken(authCode.UserID)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate access token: %w", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	
	log.Printf("OAuth Service: Generated access token: %s", accessToken)
	
	// 生成刷新令牌
	log.Printf("OAuth Service: Generating refresh token for user: %d", authCode.UserID)
	refreshToken, err := s.CreateRefreshToken(ctx, authCode.UserID, authCode.ClientID, authCode.Scopes)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	log.Printf("OAuth Service: Generated refresh token: %s", refreshToken)
	
	// 删除已使用的授权码（一次性使用）
	log.Printf("OAuth Service: Deleting used authorization code: %s", code)
	err = s.oauthRepo.DeleteAuthorizationCode(ctx, code)
	if err != nil {
		// 记录日志但不中断流程
		// 在生产环境中应该使用适当的日志库
		log.Printf("OAuth Service: Warning: failed to delete authorization code: %v", err)
	}
	
	log.Printf("OAuth Service: Authorization code exchanged successfully")
	return &ExchangeAuthorizationCodeResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateAccessToken 生成访问令牌
func (s *oauthService) GenerateAccessToken(userID int64) (string, error) {
	log.Printf("OAuth Service: Generating access token for user: %d", userID)
	// 这里简化实现，实际项目中应该使用JWT库生成标准的JWT令牌
	// 并包含用户ID、过期时间等信息
	token := generateRandomCode(128)
	log.Printf("OAuth Service: Generated access token for user %d: %s", userID, token)
	return token, nil
}

// GenerateRefreshToken 生成刷新令牌
func (s *oauthService) GenerateRefreshToken(userID int64) (string, error) {
	log.Printf("OAuth Service: Generating refresh token for user: %d", userID)
	token := generateRandomCode(128)
	log.Printf("OAuth Service: Generated refresh token for user %d: %s", userID, token)
	return token, nil
}

// CreateRefreshToken 创建刷新令牌并存储到数据库
func (s *oauthService) CreateRefreshToken(ctx context.Context, userID int64, clientID string, scopes []string) (string, error) {
	log.Printf("OAuth Service: Creating refresh token for user: %d, client: %s", userID, clientID)
	
	// 生成刷新令牌
	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate refresh token: %v", err)
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	// 哈希刷新令牌用于存储
	tokenHash := hashToken(refreshToken)
	
	// 创建刷新令牌实体
	refreshTokenModel := &model.RefreshToken{
		TokenHash: tokenHash,
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30天有效期
		RevokedAt: nil,
	}
	
	// 保存到数据库
	log.Printf("OAuth Service: Saving refresh token to database")
	err = s.oauthRepo.CreateRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		log.Printf("OAuth Service: Failed to save refresh token: %v", err)
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}
	
	log.Printf("OAuth Service: Refresh token saved successfully")
	return refreshToken, nil
}

// hashSecret 对密钥进行哈希处理
func hashSecret(secret string) string {
	// 简化实现，实际项目中应该使用bcrypt等安全的哈希算法
	log.Printf("OAuth Service: Hashing secret (first 3 chars): %s...", secret[:min(3, len(secret))])
	hash := sha256.Sum256([]byte(secret))
	result := base64.StdEncoding.EncodeToString(hash[:])
	log.Printf("OAuth Service: Secret hash result: %s", result[:min(10, len(result))])
	return result
}

// hashToken 对令牌进行哈希处理
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// generateRandomCode 生成指定长度的随机码
func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// 在生产环境中应该处理此错误
		log.Printf("OAuth Service: Error generating random code: %v", err)
		panic(err)
	}
	
	result := base64.URLEncoding.EncodeToString(bytes)[:length]
	log.Printf("OAuth Service: Generated random code (first 10 chars): %s", result[:min(10, len(result))])
	return result
}

// isValidRedirectURI 验证重定向URI是否在允许列表中
func isValidRedirectURI(redirectURI string, allowedRedirectURIs []string) bool {
	log.Printf("OAuth Service: Checking if redirect URI %s is valid", redirectURI)
	for _, uri := range allowedRedirectURIs {
		log.Printf("OAuth Service: Comparing with allowed URI: %s", uri)
		if uri == redirectURI {
			log.Printf("OAuth Service: Redirect URI is valid")
			return true
		}
	}
	log.Printf("OAuth Service: Redirect URI is not valid")
	return false
}

// areScopesAllowed 验证请求的scopes是否被允许
func areScopesAllowed(requestedScopes, allowedScopes []string) bool {
	log.Printf("OAuth Service: Checking if scopes %v are allowed", requestedScopes)
	// 创建允许scopes的map以提高查找效率
	allowed := make(map[string]bool)
	for _, scope := range allowedScopes {
		allowed[scope] = true
		log.Printf("OAuth Service: Allowed scope: %s", scope)
	}
	
	// 检查所有请求的scopes是否都被允许
	for _, scope := range requestedScopes {
		log.Printf("OAuth Service: Checking requested scope: %s", scope)
		if !allowed[scope] {
			log.Printf("OAuth Service: Scope %s is not allowed", scope)
			return false
		}
	}
	
	log.Printf("OAuth Service: All scopes are allowed")
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RefreshAccessToken 使用刷新令牌获取新的访问令牌和刷新令牌
func (s *oauthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*ExchangeAuthorizationCodeResult, error) {
	log.Printf("OAuth Service: Refreshing access token with refresh token")
	
	// 哈希传入的刷新令牌以进行数据库查找
	tokenHash := hashToken(refreshTokenString)
	
	// 查找刷新令牌
	log.Printf("OAuth Service: Finding refresh token in database")
	refreshToken, err := s.oauthRepo.FindRefreshTokenByTokenHash(ctx, tokenHash)
	if err != nil {
		log.Printf("OAuth Service: Failed to find refresh token: %v", err)
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// 验证刷新令牌是否属于有效状态（未过期且未被撤销）
	// 注意：FindRefreshTokenByTokenHash已经检查了过期和撤销状态
	
	// 生成新的访问令牌
	log.Printf("OAuth Service: Generating new access token for user: %d", refreshToken.UserID)
	accessToken, err := s.GenerateAccessToken(refreshToken.UserID)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	
	// 生成新的刷新令牌
	log.Printf("OAuth Service: Generating new refresh token for user: %d", refreshToken.UserID)
	newRefreshToken, err := s.GenerateRefreshToken(refreshToken.UserID)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate new refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	// 撤销旧的刷新令牌
	log.Printf("OAuth Service: Revoking old refresh token")
	err = s.oauthRepo.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		log.Printf("OAuth Service: Failed to revoke old refresh token: %v", err)
		// 这里我们记录错误但不中断流程，因为令牌已经使用过了
	}
	
	// 存储新的刷新令牌
	log.Printf("OAuth Service: Creating new refresh token record")
	newTokenHash := hashToken(newRefreshToken)
	newRefreshTokenModel := &model.RefreshToken{
		TokenHash: newTokenHash,
		UserID:    refreshToken.UserID,
		ClientID:  refreshToken.ClientID,
		Scopes:    refreshToken.Scopes,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30天有效期
		RevokedAt: nil,
	}
	
	err = s.oauthRepo.CreateRefreshToken(ctx, newRefreshTokenModel)
	if err != nil {
		log.Printf("OAuth Service: Failed to save new refresh token: %v", err)
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}
	
	log.Printf("OAuth Service: Access token refreshed successfully")
	return &ExchangeAuthorizationCodeResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}