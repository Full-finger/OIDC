package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

// bangumiService Bangumi服务实现
type bangumiService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	baseURL      string
	httpClient   *http.Client
	bangumiRepo  repository.BangumiRepository
}

// NewBangumiService 创建BangumiService实例
func NewBangumiService(bangumiRepo repository.BangumiRepository) BangumiService {
	return &bangumiService{
		clientID:     os.Getenv("BANGUMI_CLIENT_ID"),
		clientSecret: os.Getenv("BANGUMI_CLIENT_SECRET"),
		redirectURI:  os.Getenv("BANGUMI_REDIRECT_URI"),
		baseURL:      "https://bgm.tv",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		bangumiRepo:  bangumiRepo,
	}
}

// GetAuthorizationURL 获取Bangumi授权URL
func (s *bangumiService) GetAuthorizationURL(state string) string {
	authURL := fmt.Sprintf("%s/oauth/authorize", s.baseURL)
	
	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)
	
	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}

// ExchangeCodeForToken 用授权码换取访问令牌
func (s *bangumiService) ExchangeCodeForToken(ctx context.Context, code string) (*BangumiTokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/oauth/access_token", s.baseURL)
	
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("redirect_uri", s.redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for token, status code: %d", resp.StatusCode)
	}
	
	var tokenResponse BangumiTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	return &tokenResponse, nil
}

// RefreshToken 刷新访问令牌
func (s *bangumiService) RefreshToken(ctx context.Context, refreshToken string) (*BangumiTokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/oauth/access_token", s.baseURL)
	
	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("redirect_uri", s.redirectURI)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token, status code: %d", resp.StatusCode)
	}
	
	var tokenResponse BangumiTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	return &tokenResponse, nil
}

// GetUserInfo 获取Bangumi用户信息
func (s *bangumiService) GetUserInfo(ctx context.Context, accessToken string) (*BangumiUser, error) {
	userURL := fmt.Sprintf("%s/oauth/me", s.baseURL)
	
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status code: %d", resp.StatusCode)
	}
	
	var user BangumiUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	
	return &user, nil
}

// BindAccount 绑定Bangumi账号
func (s *bangumiService) BindAccount(ctx context.Context, userID uint, tokenResponse *BangumiTokenResponse) error {
	// 检查是否已经绑定
	existingAccount, err := s.bangumiRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing binding: %w", err)
	}
	
	if existingAccount != nil {
		// 如果已经绑定，更新令牌信息
		existingAccount.BangumiUserID = tokenResponse.UserID
		existingAccount.AccessToken = tokenResponse.AccessToken
		existingAccount.RefreshToken = tokenResponse.RefreshToken
		existingAccount.TokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
		existingAccount.Scope = tokenResponse.Scope
		existingAccount.UpdatedAt = time.Now()
		
		return s.bangumiRepo.Update(ctx, existingAccount)
	}
	
	// 创建新的绑定记录
	account := &model.BangumiAccount{
		UserID:          userID,
		BangumiUserID:   tokenResponse.UserID,
		AccessToken:     tokenResponse.AccessToken,
		RefreshToken:    tokenResponse.RefreshToken,
		TokenExpiresAt:  time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Scope:           tokenResponse.Scope,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	return s.bangumiRepo.Create(ctx, account)
}

// UnbindAccount 解绑Bangumi账号
func (s *bangumiService) UnbindAccount(ctx context.Context, userID uint) error {
	account, err := s.bangumiRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get bangumi account: %w", err)
	}
	
	if account == nil {
		return nil // 没有绑定记录，直接返回
	}
	
	return s.bangumiRepo.DeleteByID(ctx, account.ID)
}

// GetBoundAccount 获取已绑定的Bangumi账号
func (s *bangumiService) GetBoundAccount(ctx context.Context, userID uint) (*model.BangumiAccount, error) {
	return s.bangumiRepo.GetByUserID(ctx, userID)
}