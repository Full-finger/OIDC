package model

import (
	"time"
)

// Client OAuth2客户端实体
type Client struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID    string    `gorm:"uniqueIndex;not null" json:"client_id"`
	SecretHash  string    `gorm:"not null" json:"secret_hash"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	RedirectURI string    `gorm:"not null" json:"redirect_uri"`
	Scopes      string    `gorm:"not null" json:"scopes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuthorizationCode OAuth2授权码实体
type AuthorizationCode struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code               string    `gorm:"uniqueIndex;not null" json:"code"`
	ClientID           string    `gorm:"not null" json:"client_id"`
	UserID             uint      `gorm:"not null" json:"user_id"`
	RedirectURI        string    `gorm:"not null" json:"redirect_uri"`
	Scopes             string    `gorm:"not null" json:"scopes"`
	CodeChallenge      string    `gorm:"type:text" json:"code_challenge"`
	CodeChallengeMethod string   `gorm:"type:text" json:"code_challenge_method"`
	ExpiresAt          time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt          time.Time `json:"created_at"`
}

// RefreshToken OAuth2刷新令牌实体
type RefreshToken struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenHash   string    `gorm:"uniqueIndex;not null" json:"token_hash"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	ClientID    string    `gorm:"not null" json:"client_id"`
	Scopes      string    `gorm:"type:text" json:"scopes"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
	RevokedAt   time.Time `gorm:"index" json:"revoked_at"`
	CreatedAt   time.Time `json:"created_at"`
}