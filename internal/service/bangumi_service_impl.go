package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	animeRepo    repository.AnimeRepository
	collectionRepo repository.CollectionRepository
}

// NewBangumiService 创建BangumiService实例
func NewBangumiService(bangumiRepo repository.BangumiRepository, animeRepo repository.AnimeRepository, collectionRepo repository.CollectionRepository) BangumiService {
	return &bangumiService{
		clientID:     os.Getenv("BANGUMI_CLIENT_ID"),
		clientSecret: os.Getenv("BANGUMI_CLIENT_SECRET"),
		redirectURI:  os.Getenv("BANGUMI_REDIRECT_URI"),
		baseURL:      "https://bgm.tv",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		bangumiRepo:  bangumiRepo,
		animeRepo:    animeRepo,
		collectionRepo: collectionRepo,
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

// SyncCollection 同步Bangumi收藏数据
func (s *bangumiService) SyncCollection(ctx context.Context, userID uint) error {
	// 获取用户绑定的Bangumi账号
	account, err := s.GetBoundAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get bound bangumi account: %w", err)
	}
	
	if account == nil {
		return fmt.Errorf("no bangumi account bound for user %d", userID)
	}
	
	// 检查令牌是否过期，如果过期则刷新
	accessToken := account.AccessToken
	if time.Now().After(account.TokenExpiresAt) {
		tokenResponse, err := s.RefreshToken(ctx, account.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		
		// 更新令牌信息
		account.AccessToken = tokenResponse.AccessToken
		account.RefreshToken = tokenResponse.RefreshToken
		account.TokenExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
		account.UpdatedAt = time.Now()
		
		if err := s.bangumiRepo.Update(ctx, account); err != nil {
			return fmt.Errorf("failed to update account tokens: %w", err)
		}
		
		accessToken = tokenResponse.AccessToken
	}
	
	// 获取Bangumi收藏数据
	bangumiCollections, err := s.fetchBangumiCollections(ctx, accessToken, account.BangumiUserID)
	if err != nil {
		return fmt.Errorf("failed to fetch bangumi collections: %w", err)
	}
	
	// 转换并保存数据
	for _, bangumiCollection := range bangumiCollections {
		// 转换Bangumi收藏为本地番剧和收藏记录
		if err := s.convertAndSaveCollection(ctx, userID, bangumiCollection); err != nil {
			// 记录错误但继续处理其他收藏
			fmt.Printf("Warning: failed to convert and save collection for subject %d: %v\n", bangumiCollection.SubjectID, err)
		}
	}
	
	return nil
}

// fetchBangumiCollections 获取Bangumi收藏数据
func (s *bangumiService) fetchBangumiCollections(ctx context.Context, accessToken string, bangumiUserID uint) ([]*BangumiCollection, error) {
	// Bangumi API获取用户收藏的URL
	collectionsURL := fmt.Sprintf("%s/api/users/%d/collections", s.baseURL, bangumiUserID)
	
	req, err := http.NewRequest("GET", collectionsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create collections request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	// 设置Accept头以获取JSON响应
	req.Header.Set("Accept", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get collections, status code: %d", resp.StatusCode)
	}
	
	// 解析响应
	var collections []*BangumiCollection
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("failed to decode collections: %w", err)
	}
	
	return collections, nil
}

// convertAndSaveCollection 转换并保存收藏数据
func (s *bangumiService) convertAndSaveCollection(ctx context.Context, userID uint, bangumiCollection *BangumiCollection) error {
	// 查找或创建番剧
	anime, err := s.animeRepo.GetByTitle(ctx, bangumiCollection.Name)
	if err != nil {
		return fmt.Errorf("failed to get anime by title: %w", err)
	}
	
	// 如果番剧不存在，则创建
	if anime == nil {
		anime = &model.Anime{
			Title:       bangumiCollection.Name,
			Description: bangumiCollection.Summary,
			CoverImage:  bangumiCollection.Image,
			Episodes:    bangumiCollection.Episodes,
			Status:      s.convertBangumiStatus(bangumiCollection.Status),
			Rating:      bangumiCollection.Rating,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		// 设置发布日期（如果有的话）
		if bangumiCollection.UpdatedAt != "" {
			if releaseDate, err := time.Parse("2006-01-02", bangumiCollection.UpdatedAt); err == nil {
				anime.ReleaseDate = releaseDate
			}
		}
		
		if err := s.animeRepo.Create(ctx, anime); err != nil {
			return fmt.Errorf("failed to create anime: %w", err)
		}
	}
	
	// 查找现有的收藏记录
	collection, err := s.collectionRepo.GetByUserIDAndAnimeID(ctx, userID, anime.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing collection: %w", err)
	}
	
	// 转换Bangumi收藏状态为本地状态
	localStatus := s.convertBangumiUserStatus(bangumiCollection.UserStatus)
	
	// 如果收藏记录不存在，则创建
	if collection == nil {
		rating := bangumiCollection.UserRating
		collection = &model.Collection{
			UserID:    userID,
			AnimeID:   anime.ID,
			Status:    localStatus,
			Rating:    &rating,
			Comment:   bangumiCollection.Comment,
			Progress:  0, // Bangumi数据中没有进度信息，需要用户手动更新
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if err := s.collectionRepo.Create(ctx, collection); err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	} else {
		// 更新现有收藏记录
		collection.Status = localStatus
		collection.Rating = &bangumiCollection.UserRating
		collection.Comment = bangumiCollection.Comment
		collection.UpdatedAt = time.Now()
		
		if err := s.collectionRepo.Update(ctx, collection); err != nil {
			return fmt.Errorf("failed to update collection: %w", err)
		}
	}
	
	return nil
}

// convertBangumiStatus 转换Bangumi番剧状态为本地状态
func (s *bangumiService) convertBangumiStatus(status string) string {
	switch status {
	case "airing":
		return "airing"
	case "finished":
		return "finished"
	default:
		return "upcoming"
	}
}

// convertBangumiUserStatus 转换Bangumi用户收藏状态为本地状态
func (s *bangumiService) convertBangumiUserStatus(status string) string {
	switch status {
	case "watched":
		return "completed"
	case "watching":
		return "watching"
	case "want_watch":
		return "plan_to_watch"
	case "on_hold":
		return "on_hold"
	case "dropped":
		return "dropped"
	default:
		return "watching"
	}
}