package service

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// BangumiService Bangumi服务接口
type BangumiService interface {
	// GetAuthorizationURL 获取Bangumi授权URL
	GetAuthorizationURL(state string) string
	
	// ExchangeCodeForToken 用授权码换取访问令牌
	ExchangeCodeForToken(ctx context.Context, code string) (*BangumiTokenResponse, error)
	
	// RefreshToken 刷新访问令牌
	RefreshToken(ctx context.Context, refreshToken string) (*BangumiTokenResponse, error)
	
	// GetUserInfo 获取Bangumi用户信息
	GetUserInfo(ctx context.Context, accessToken string) (*BangumiUser, error)
	
	// BindAccount 绑定Bangumi账号
	BindAccount(ctx context.Context, userID uint, tokenResponse *BangumiTokenResponse) error
	
	// UnbindAccount 解绑Bangumi账号
	UnbindAccount(ctx context.Context, userID uint) error
	
	// GetBoundAccount 获取已绑定的Bangumi账号
	GetBoundAccount(ctx context.Context, userID uint) (*model.BangumiAccount, error)
	
	// SyncCollection 同步Bangumi收藏数据
	SyncCollection(ctx context.Context, userID uint) error
}

// BangumiTokenResponse Bangumi令牌响应
type BangumiTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UserID       uint   `json:"user_id"`
}

// BangumiUser Bangumi用户信息
type BangumiUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// BangumiCollection Bangumi收藏信息
type BangumiCollection struct {
	SubjectID   uint   `json:"subject_id"`
	Name        string `json:"name"`
	NameCN      string `json:"name_cn"`
	Summary     string `json:"summary"`
	Image       string `json:"image"`
	Episodes    int    `json:"episodes"`
	Status      string `json:"status"`
	Rating      float64 `json:"rating"`
	UserStatus  string `json:"user_status"`
	UserRating  float64 `json:"user_rating"`
	Comment     string `json:"comment"`
	UpdatedAt   string `json:"updated_at"`
}