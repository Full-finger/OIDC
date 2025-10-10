// Package service defines the service layer interfaces for the OIDC application.
package service

import (
	"context"

	"github.com/Full-finger/OIDC/internal/model"
)

// OAuthService defines the OAuth service interface
type OAuthService interface {
	IBaseService
	ConvertInterface

	// Client related operations
	CreateClient(ctx context.Context, clientID, clientSecret, name string, redirectURIs, scopes []string) (*model.Client, error)
	FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error)
	
	// AuthorizationCode related operations
	GenerateAuthorizationCode(ctx context.Context, clientID string, userID int64, redirectURI string, scopes []string, codeChallenge *string, codeChallengeMethod *string) (*model.AuthorizationCode, error)
	ValidateAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*model.AuthorizationCode, error)
	ExchangeAuthorizationCode(ctx context.Context, code, clientID, redirectURI string) (*ExchangeAuthorizationCodeResult, error)
	
	// Token related operations
	GenerateAccessToken(userID int64, scopes []string) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	CreateRefreshToken(ctx context.Context, userID int64, clientID string, scopes []string) (string, error)
	RefreshAccessToken(ctx context.Context, refreshTokenString string) (*ExchangeAuthorizationCodeResult, error)

	// Legacy methods - to be deprecated in future versions
	// GetClientByID gets a client by client ID
	GetClientByID(ctx context.Context, clientID string) (*model.Client, error)

	// CreateAuthorizationCode creates an authorization code
	CreateAuthorizationCode(ctx context.Context, clientID string, userID int64, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod *string) (string, error)

	// ExchangeCodeForToken exchanges an authorization code for tokens
	ExchangeCodeForToken(ctx context.Context, code, codeVerifier string) (*TokenResponse, error)

	// GetUserInfo gets user information
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponse, error)

	// ValidateClientSecret validates a client secret
	ValidateClientSecret(ctx context.Context, clientID, clientSecret string) error
}

// TokenResponse represents the token response structure
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// UserInfoResponse represents the user info response structure
type UserInfoResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

// ExchangeAuthorizationCodeResult represents the result of exchanging an authorization code
type ExchangeAuthorizationCodeResult struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	IDToken      string   `json:"id_token,omitempty"`
	Scope        []string `json:"scope,omitempty"`
}