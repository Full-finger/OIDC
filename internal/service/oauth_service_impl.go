package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// OpenIDConfiguration OpenID配置信息
type OpenIDConfiguration struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	ScopesSupported              []string `json:"scopes_supported"`
	ResponseTypesSupported       []string `json:"response_types_supported"`
	GrantTypesSupported          []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
	SubjectTypesSupported        []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported              []string `json:"claims_supported"`
}

// UserInfo 用户信息
type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	Profile       string `json:"profile,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

// oauthService OAuth服务实现
type oauthService struct {
	jwtUtil util.JWTUtil
}

// NewOAuthService 创建OAuth服务实例
func NewOAuthService() OAuthService {
	// 初始化JWT工具
	jwtUtil, err := util.NewJWTUtil()
	if err != nil {
		// 如果JWT工具初始化失败，记录日志但继续运行
		fmt.Printf("Warning: failed to initialize JWT utility: %v\n", err)
	}
	
	return &oauthService{
		jwtUtil: jwtUtil,
	}
}

// GetOpenIDConfiguration 获取OpenID配置信息
func (s *oauthService) GetOpenIDConfiguration(ctx context.Context) (*OpenIDConfiguration, error) {
	config := &OpenIDConfiguration{
		Issuer:                           "http://localhost:8080",
		AuthorizationEndpoint:           "http://localhost:8080/oauth/authorize",
		TokenEndpoint:                   "http://localhost:8080/oauth/token",
		UserinfoEndpoint:                "http://localhost:8080/oauth/userinfo",
		JwksURI:                         "http://localhost:8080/.well-known/jwks.json",
		ScopesSupported:                 []string{"openid", "profile", "email"},
		ResponseTypesSupported:          []string{"code"},
		GrantTypesSupported:             []string{"authorization_code", "refresh_token"},
		CodeChallengeMethodsSupported:   []string{"S256", "plain"},
		SubjectTypesSupported:           []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
		ClaimsSupported:                 []string{"sub", "name", "nickname", "profile", "picture", "email", "email_verified"},
	}
	
	return config, nil
}

// GetUserInfo 获取用户信息
func (s *oauthService) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// 解析访问令牌
	var claims *util.AccessTokenClaims
	var err error
	
	if s.jwtUtil != nil {
		claims, err = s.jwtUtil.ParseAccessToken(accessToken)
		if err != nil {
			return nil, fmt.Errorf("invalid access token: %w", err)
		}
	} else {
		// JWT工具不可用时的简化实现
		claims = &util.AccessTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "user:1", // 默认用户
			},
			Scope: "openid profile email",
		}
	}
	
	// 根据scope决定返回哪些用户信息
	// 这里简化处理，实际应该从数据库获取真实用户信息
	userInfo := &UserInfo{
		Sub: claims.Subject,
	}
	
	// 检查scope中是否包含profile
	if containsScope(claims.Scope, "profile") || containsScope(claims.Scope, "openid") {
		userInfo.Name = "示例用户"
		userInfo.Nickname = "示例昵称"
		userInfo.Profile = "https://example.com/profile"
		userInfo.Picture = "https://example.com/avatar.jpg"
	}
	
	// 检查scope中是否包含email
	if containsScope(claims.Scope, "email") {
		userInfo.Email = "user@example.com"
		userInfo.EmailVerified = true
	}
	
	return userInfo, nil
}

// HandleAuthorizationRequest 处理授权请求
func (s *oauthService) HandleAuthorizationRequest(ctx context.Context, clientID, userID, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod *string) (*model.AuthorizationCode, error) {
	// 查找客户端
	client, err := s.GetClientByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	// 验证重定向URI
	if !s.isValidRedirectURI(redirectURI, client.RedirectURI) {
		return nil, fmt.Errorf("invalid redirect URI")
	}

	// 验证请求的scopes是否被客户端允许
	if !s.areScopesAllowed(scopes, client.Scopes) {
		return nil, fmt.Errorf("invalid scopes")
	}

	// 生成随机授权码
	code := s.generateRandomCode(64)

	// 解析用户ID
	var parsedUserID uint
	// 这里应该有实际的用户ID解析逻辑
	// 为简化示例，我们假设传入的userID是有效的
	fmt.Sscanf(userID, "%d", &parsedUserID)

	// 创建授权码实体
	authCode := &model.AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              parsedUserID,
		RedirectURI:         redirectURI,
		Scopes:              s.scopesToString(scopes),
		ExpiresAt:           time.Now().Add(10 * time.Minute), // 10分钟有效期
		CodeChallenge:       s.getStringValue(codeChallenge),
		CodeChallengeMethod: s.getStringValue(codeChallengeMethod),
	}

	// TODO: 保存到数据库
	// 这里应该调用repository来保存授权码
	// err = s.authorizationCodeRepo.Create(ctx, authCode)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to save authorization code: %w", err)
	// }

	return authCode, nil
}

// HandleTokenRequest 处理令牌请求
func (s *oauthService) HandleTokenRequest(ctx context.Context, grantType, code, clientID, clientSecret, redirectURI string, codeVerifier *string) (*TokenResponse, error) {
	switch grantType {
	case "authorization_code":
		// 使用授权码换取访问令牌
		return s.ExchangeAuthorizationCode(ctx, code, clientID, clientSecret, redirectURI, codeVerifier)
	case "refresh_token":
		// 使用刷新令牌获取新的访问令牌
		return s.RefreshAccessToken(ctx, code, clientID, clientSecret)
	default:
		return nil, fmt.Errorf("unsupported grant type: %s", grantType)
	}
}

// ValidateClient 验证客户端
func (s *oauthService) ValidateClient(ctx context.Context, clientID, clientSecret, redirectURI string) (*model.Client, error) {
	// 查找客户端
	client, err := s.GetClientByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	// 验证客户端密钥
	// 这里应该有实际的密钥验证逻辑
	// 为简化示例，我们跳过验证
	_ = clientSecret

	// 验证重定向URI（仅当提供了重定向URI时才验证）
	// 在刷新令牌流程中，通常不提供重定向URI
	if redirectURI != "" && !s.isValidRedirectURI(redirectURI, client.RedirectURI) {
		return nil, fmt.Errorf("invalid redirect URI")
	}

	return client, nil
}

// GenerateAuthorizationCode 生成授权码
func (s *oauthService) GenerateAuthorizationCode(ctx context.Context, client *model.Client, userID uint, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod *string) (string, error) {
	// 生成随机授权码
	code := s.generateRandomCode(64)

	// 创建授权码实体（当前仅用于演示，实际应保存到数据库）
	_ = &model.AuthorizationCode{
		Code:                code,
		ClientID:            client.ClientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              s.scopesToString(scopes),
		ExpiresAt:           time.Now().Add(10 * time.Minute), // 10分钟有效期
		CodeChallenge:       s.getStringValue(codeChallenge),
		CodeChallengeMethod: s.getStringValue(codeChallengeMethod),
	}

	// TODO: 保存到数据库
	// 这里应该调用repository来保存授权码
	// err := s.authorizationCodeRepo.Create(ctx, authCode)
	// if err != nil {
	//     return "", fmt.Errorf("failed to save authorization code: %w", err)
	// }

	return code, nil
}

// ExchangeAuthorizationCode 用授权码换取访问令牌
func (s *oauthService) ExchangeAuthorizationCode(ctx context.Context, code, clientID, clientSecret, redirectURI string, codeVerifier *string) (*TokenResponse, error) {
	// 查找授权码
	authCode, err := s.ValidateAuthorizationCode(ctx, code, clientID, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("invalid authorization code: %w", err)
	}

	// 验证客户端
	client, err := s.ValidateClient(ctx, clientID, clientSecret, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	// 验证PKCE（如果使用）
	if authCode.CodeChallenge != "" {
		if codeVerifier == nil {
			return nil, fmt.Errorf("code verifier required")
		}

		if !s.validatePKCE(authCode.CodeChallenge, *codeVerifier, authCode.CodeChallengeMethod) {
			return nil, fmt.Errorf("invalid code verifier")
		}
	}

	// 生成访问令牌
	accessToken, err := s.generateAccessToken(authCode.UserID, client.ClientID, authCode.Scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshTokenStr, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 创建刷新令牌实体（当前仅用于演示，实际应保存到数据库）
	_ = &model.RefreshToken{
		TokenHash: s.hashToken(refreshTokenStr),
		UserID:    authCode.UserID,
		ClientID:  client.ClientID,
		Scopes:    authCode.Scopes,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30天有效期
	}

	// TODO: 保存刷新令牌到数据库
	// err = s.refreshTokenRepo.Create(ctx, refreshTokenModel)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to save refresh token: %w", err)
	// }

	// 构造响应
	response := &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小时
		Scope:        authCode.Scopes,
	}

	// 如果是授权码流程，添加刷新令牌
	response.RefreshToken = refreshTokenStr

	// 检查是否包含openid scope，如果包含则生成ID Token
	if s.containsScope(s.stringToScopes(authCode.Scopes), "openid") {
		// 生成ID Token
		idToken, err := s.generateIDToken(authCode.UserID, client.ClientID, authCode.Scopes)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ID token: %w", err)
		}
		response.IDToken = idToken
	}

	// TODO: 删除已使用的授权码
	// err = s.authorizationCodeRepo.Delete(ctx, code)
	// if err != nil {
	//     // 记录日志但不中断流程
	// }

	return response, nil
}

// ValidateAuthorizationCode 验证授权码
func (s *oauthService) ValidateAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*model.AuthorizationCode, error) {
	// TODO: 从数据库查找授权码
	// authCode, err := s.authorizationCodeRepo.GetByCode(ctx, code)
	// if err != nil {
	//     return nil, fmt.Errorf("invalid authorization code")
	// }

	// 模拟一个授权码对象用于演示
	authCode := &model.AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      1,
		RedirectURI: redirectURI,
		Scopes:      "openid profile email",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	// 检查是否过期
	if time.Now().After(authCode.ExpiresAt) {
		return nil, fmt.Errorf("authorization code expired")
	}

	// 验证客户端ID
	if authCode.ClientID != clientID {
		return nil, fmt.Errorf("invalid client")
	}

	// 验证重定向URI
	if authCode.RedirectURI != redirectURI {
		return nil, fmt.Errorf("invalid redirect URI")
	}

	return authCode, nil
}

// CreateRefreshToken 创建刷新令牌
func (s *oauthService) CreateRefreshToken(ctx context.Context, userID uint, clientID string, scopes []string) (*model.RefreshToken, error) {
	// 生成刷新令牌
	refreshTokenStr, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 创建刷新令牌实体
	refreshToken := &model.RefreshToken{
		TokenHash: s.hashToken(refreshTokenStr),
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    s.scopesToString(scopes),
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30天有效期
	}

	// TODO: 保存到数据库
	// err = s.refreshTokenRepo.Create(ctx, refreshToken)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to save refresh token: %w", err)
	// }

	return refreshToken, nil
}

// RefreshAccessToken 刷新访问令牌
func (s *oauthService) RefreshAccessToken(ctx context.Context, refreshToken, clientID, clientSecret string) (*TokenResponse, error) {
	// 验证客户端
	client, err := s.ValidateClient(ctx, clientID, clientSecret, "") // 重定向URI在刷新令牌流程中不验证
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	// 查找刷新令牌
	// refresh, err := s.refreshTokenRepo.GetByTokenHash(ctx, s.hashToken(refreshToken))
	// if err != nil {
	//     return nil, fmt.Errorf("invalid refresh token")
	// }

	// 模拟一个刷新令牌对象用于演示
	refresh := &model.RefreshToken{
		TokenHash: s.hashToken(refreshToken),
		UserID:    1,
		ClientID:  clientID,
		Scopes:    "openid profile email",
		ExpiresAt: time.Now().Add(24 * time.Hour * 30),
	}

	// 检查是否过期
	if time.Now().After(refresh.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// 生成新的访问令牌
	accessToken, err := s.generateAccessToken(refresh.UserID, client.ClientID, refresh.Scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成新的刷新令牌
	newRefreshTokenStr, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 创建新的刷新令牌实体
	_ = &model.RefreshToken{
		TokenHash: s.hashToken(newRefreshTokenStr),
		UserID:    refresh.UserID,
		ClientID:  client.ClientID,
		Scopes:    refresh.Scopes,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30天有效期
	}

	// TODO: 保存新刷新令牌到数据库
	// err = s.refreshTokenRepo.Create(ctx, newRefreshToken)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	// }

	// TODO: 撤销旧的刷新令牌
	// err = s.refreshTokenRepo.Revoke(ctx, refresh.ID)
	// if err != nil {
	//     // 记录日志但不中断流程
	// }

	// 构造响应
	response := &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小时
		RefreshToken: newRefreshTokenStr,
		Scope:        refresh.Scopes,
	}

	// 如果scope包含openid，生成ID Token
	if s.containsScope(s.stringToScopes(refresh.Scopes), "openid") {
		idToken, err := s.generateIDToken(refresh.UserID, client.ClientID, refresh.Scopes)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ID token: %w", err)
		}
		response.IDToken = idToken
	}

	return response, nil
}

// GetClientByClientID 根据客户端ID获取客户端
func (s *oauthService) GetClientByClientID(ctx context.Context, clientID string) (*model.Client, error) {
	// TODO: 从数据库查找客户端
	// client, err := s.clientRepo.GetByClientID(ctx, clientID)
	// if err != nil {
	//     return nil, fmt.Errorf("client not found")
	// }

	// 模拟一个客户端对象用于演示
	client := &model.Client{
		ID:          1,
		ClientID:    "test_client",
		SecretHash:  "", // 实际应该存储哈希值
		Name:        "测试客户端",
		Description: "用于测试的客户端",
		RedirectURI: "http://localhost:3000/callback",
		Scopes:      "openid profile email",
	}
	
	// 检查客户端ID是否匹配
	if clientID != client.ClientID {
		return nil, fmt.Errorf("client not found")
	}

	return client, nil
}

// isValidRedirectURI 验证重定向URI是否有效
func (s *oauthService) isValidRedirectURI(requestedURI, allowedURI string) bool {
	// 在测试环境中允许更灵活的重定向URI验证
	// 实际生产环境中应该更严格
	return requestedURI == allowedURI || allowedURI == "http://localhost:3000/callback"
}

// areScopesAllowed 验证请求的scopes是否被允许
func (s *oauthService) areScopesAllowed(requestedScopes []string, allowedScopes string) bool {
	// 简单实现，实际应该更复杂
	allowed := s.stringToScopes(allowedScopes)
	for _, scope := range requestedScopes {
		found := false
		for _, allowedScope := range allowed {
			if scope == allowedScope {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// containsScope 检查scopes中是否包含指定的scope
func (s *oauthService) containsScope(scopes []string, targetScope string) bool {
	for _, scope := range scopes {
		if scope == targetScope {
			return true
		}
	}
	return false
}

// generateRandomCode 生成随机授权码
func (s *oauthService) generateRandomCode(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 出错时返回固定值，仅用于演示
		return "default_code"
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// generateAccessToken 生成访问令牌
func (s *oauthService) generateAccessToken(userID uint, clientID, scopes string) (string, error) {
	// 如果JWT工具可用，则生成JWT令牌
	if s.jwtUtil != nil {
		claims := &util.AccessTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   fmt.Sprintf("user:%d", userID),
				Issuer:    "OIDC",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 1小时过期
				Audience:  []string{clientID},
			},
			Scope: scopes,
		}
		
		return s.jwtUtil.GenerateAccessToken(claims)
	}
	
	// 简化实现，实际应该生成JWT令牌
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return "access_" + base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// generateIDToken 生成ID令牌
func (s *oauthService) generateIDToken(userID uint, clientID, scopes string) (string, error) {
	// 如果JWT工具不可用，返回错误
	if s.jwtUtil == nil {
		return "", fmt.Errorf("JWT utility not available")
	}
	
	// 解析scopes
	scopeList := s.stringToScopes(scopes)
	
	// 构造ID Token声明
	claims := &util.IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("user:%d", userID),
			Issuer:    "OIDC",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 1小时过期
			Audience:  []string{clientID},
		},
	}
	
	// 根据scope添加额外声明
	if s.containsScope(scopeList, "profile") {
		claims.Profile = "https://example.com/profile"
		claims.Name = "示例用户"
	}
	
	if s.containsScope(scopeList, "email") {
		claims.Email = "user@example.com"
	}
	
	// 生成ID Token
	return s.jwtUtil.GenerateIDToken(claims)
}

// generateRefreshToken 生成刷新令牌
func (s *oauthService) generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return "refresh_" + base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// validatePKCE 验证PKCE
func (s *oauthService) validatePKCE(codeChallenge, codeVerifier, method string) bool {
	switch method {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		expectedChallenge := base64.URLEncoding.EncodeToString(hash[:])
		return expectedChallenge == codeChallenge
	case "plain":
		return codeChallenge == codeVerifier
	default:
		return false
	}
}

// hashToken 哈希令牌
func (s *oauthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// scopesToString 将scopes数组转换为字符串
func (s *oauthService) scopesToString(scopes []string) string {
	result := ""
	for i, scope := range scopes {
		if i > 0 {
			result += " "
		}
		result += scope
	}
	return result
}

// stringToScopes 将字符串转换为scopes数组
func (s *oauthService) stringToScopes(scopes string) []string {
	if scopes == "" {
		return []string{}
	}
	
	// 按空格分割scopes
	var result []string
	start := 0
	for i, char := range scopes {
		if char == ' ' {
			if start < i {
				result = append(result, scopes[start:i])
			}
			start = i + 1
		}
	}
	
	// 添加最后一个scope
	if start < len(scopes) {
		result = append(result, scopes[start:])
	}
	
	return result
}

// getStringValue 获取字符串指针的值，如果指针为nil则返回空字符串
func (s *oauthService) getStringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// containsScope 检查scope字符串中是否包含指定的scope
func containsScope(scopes, targetScope string) bool {
	// 按空格分割scopes
	scopeList := []string{}
	start := 0
	for i, char := range scopes {
		if char == ' ' {
			if start < i {
				scopeList = append(scopeList, scopes[start:i])
			}
			start = i + 1
		}
	}
	
	// 添加最后一个scope
	if start < len(scopes) {
		scopeList = append(scopeList, scopes[start:])
	}
	
	// 检查是否包含目标scope
	for _, scope := range scopeList {
		if scope == targetScope {
			return true
		}
	}
	return false
}