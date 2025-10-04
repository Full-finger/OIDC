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
	ExchangeAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (string, error)
	
	// Token相关操作
	GenerateAccessToken(userID int64) (string, error)
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

// ExchangeAuthorizationCode 兑换授权码获取访问令牌
func (s *oauthService) ExchangeAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (string, error) {
	log.Printf("OAuth Service: Exchanging authorization code: %s for client: %s, redirectURI: %s", code, clientID, redirectURI)
	
	// 验证授权码
	authCode, err := s.ValidateAuthorizationCode(ctx, code, clientID, redirectURI)
	if err != nil {
		log.Printf("OAuth Service: Failed to validate authorization code: %v", err)
		return "", err
	}
	
	// 生成访问令牌
	log.Printf("OAuth Service: Generating access token for user: %d", authCode.UserID)
	accessToken, err := s.GenerateAccessToken(authCode.UserID)
	if err != nil {
		log.Printf("OAuth Service: Failed to generate access token: %w", err)
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	
	log.Printf("OAuth Service: Generated access token: %s", accessToken)
	
	// 删除已使用的授权码（一次性使用）
	log.Printf("OAuth Service: Deleting used authorization code: %s", code)
	err = s.oauthRepo.DeleteAuthorizationCode(ctx, code)
	if err != nil {
		// 记录日志但不中断流程
		// 在生产环境中应该使用适当的日志库
		log.Printf("OAuth Service: Warning: failed to delete authorization code: %v", err)
	}
	
	log.Printf("OAuth Service: Authorization code exchanged successfully")
	return accessToken, nil
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

// hashSecret 对密钥进行哈希处理
func hashSecret(secret string) string {
	// 简化实现，实际项目中应该使用bcrypt等安全的哈希算法
	log.Printf("OAuth Service: Hashing secret (first 3 chars): %s...", secret[:min(3, len(secret))])
	hash := sha256.Sum256([]byte(secret))
	result := base64.StdEncoding.EncodeToString(hash[:])
	log.Printf("OAuth Service: Secret hash result: %s", result[:min(10, len(result))])
	return result
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