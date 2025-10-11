package service

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// OAuthService OAuth服务接口
type OAuthService interface {
	// HandleAuthorizationRequest 处理授权请求
	HandleAuthorizationRequest(ctx context.Context, clientID, userID, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod *string) (*model.AuthorizationCode, error)
	
	// HandleTokenRequest 处理令牌请求
	HandleTokenRequest(ctx context.Context, grantType, code, clientID, clientSecret, redirectURI string, codeVerifier *string) (*TokenResponse, error)
	
	// ValidateClient 验证客户端
	ValidateClient(ctx context.Context, clientID, clientSecret, redirectURI string) (*model.Client, error)
	
	// GenerateAuthorizationCode 生成授权码
	GenerateAuthorizationCode(ctx context.Context, client *model.Client, userID uint, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod *string) (string, error)
	
	// ExchangeAuthorizationCode 兑换授权码获取访问令牌
	ExchangeAuthorizationCode(ctx context.Context, code, clientID, clientSecret, redirectURI string, codeVerifier *string) (*TokenResponse, error)
	
	// ValidateAuthorizationCode 验证授权码
	ValidateAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*model.AuthorizationCode, error)
	
	// CreateRefreshToken 创建刷新令牌
	CreateRefreshToken(ctx context.Context, userID uint, clientID string, scopes []string) (*model.RefreshToken, error)
	
	// RefreshAccessToken 刷新访问令牌
	RefreshAccessToken(ctx context.Context, refreshToken, clientID, clientSecret string) (*TokenResponse, error)
	
	// GetClientByClientID 根据客户端ID获取客户端
	GetClientByClientID(ctx context.Context, clientID string) (*model.Client, error)
}