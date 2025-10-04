package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BangumiTokenResponse Bangumi令牌响应
type BangumiTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// BangumiUserResponse Bangumi用户信息响应
type BangumiUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Sign     string `json:"sign"`
}

// BangumiCollection Bangumi收藏条目
type BangumiCollection struct {
	ID        int64          `json:"id"`
	SubjectID int64          `json:"subject_id"`
	Subject   BangumiSubject `json:"subject"`
	Rate      int            `json:"rate"`
	Type      int            `json:"type"`
	Comment   string         `json:"comment"`
	EpStatus  int            `json:"ep_status"`
	VolStatus int            `json:"vol_status"`
	Private   bool           `json:"private"`
	LastTouch int64          `json:"lasttouch"`
}

// BangumiSubject Bangumi条目信息
type BangumiSubject struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NameCN   string `json:"name_cn"`
	Image    string `json:"image"`
	Platform string `json:"platform"`
	Summary  string `json:"summary"`
}

// EpisodeCount 从Summary中提取集数信息（简化实现）
func (s *BangumiSubject) EpisodeCount() int {
	// 这里应该实现从Summary或其他字段中提取集数的逻辑
	// 现在返回默认值0
	return 0
}

// BangumiCollectionResponse Bangumi收藏列表响应
type BangumiCollectionResponse struct {
	Data []BangumiCollection `json:"data"`
}

// BangumiClient Bangumi API客户端
type BangumiClient struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	redirectURI  string
	tokenURL     string
	userInfoURL  string
}

// NewBangumiClient 创建Bangumi客户端实例
func NewBangumiClient(clientID, clientSecret, redirectURI, tokenURL, userInfoURL string) *BangumiClient {
	return &BangumiClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		tokenURL:     tokenURL,
		userInfoURL:  userInfoURL,
	}
}

// ExchangeCodeForToken 用授权码换取访问令牌
func (c *BangumiClient) ExchangeCodeForToken(code string) (*BangumiTokenResponse, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", c.clientID)
	params.Set("client_secret", c.clientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", c.redirectURI)

	resp, err := c.httpClient.Post(
		c.tokenURL,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(params.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp BangumiTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo 获取用户信息
func (c *BangumiClient) GetUserInfo(accessToken string) (*BangumiUserResponse, error) {
	req, err := http.NewRequest("GET", c.userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userResp BangumiUserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	return &userResp, nil
}

// GetUserCollections 获取用户收藏列表
func (c *BangumiClient) GetUserCollections(accessToken, username string) (*BangumiCollectionResponse, error) {
	// 构建获取用户收藏的URL
	// 注意：Bangumi API的用户收藏接口可能需要根据实际API文档进行调整
	url := fmt.Sprintf("https://api.bgm.tv/v0/users/%s/collections", username)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create collections request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "bangumoe-oidc/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user collections: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read collections response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user collections failed with status %d: %s", resp.StatusCode, string(body))
	}

	var collectionResp BangumiCollectionResponse
	if err := json.Unmarshal(body, &collectionResp); err != nil {
		// 如果解析失败，尝试解析为数组格式
		var collections []BangumiCollection
		if err2 := json.Unmarshal(body, &collections); err2 != nil {
			return nil, fmt.Errorf("failed to parse collections response: %w, alternative error: %w", err, err2)
		}
		collectionResp.Data = collections
	}

	return &collectionResp, nil
}