package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
)

// BangumiOAuthConfig Bangumi OAuth配置
type BangumiOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	BaseURL      string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

// GetBangumiOAuthConfig 获取Bangumi OAuth配置
func GetBangumiOAuthConfig() *BangumiOAuthConfig {
	baseURL := getEnv("BANGUMI_BASE_URL", "https://bgm.tv")
	
	return &BangumiOAuthConfig{
		ClientID:     os.Getenv("BANGUMI_CLIENT_ID"),
		ClientSecret: os.Getenv("BANGUMI_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("BANGUMI_REDIRECT_URI"),
		BaseURL:      baseURL,
		AuthURL:      baseURL + "/oauth/authorize",
		TokenURL:     baseURL + "/oauth/access_token",
		UserInfoURL:  "https://api.bgm.tv/v0/me",
	}
}

// GenerateState 生成随机state参数
func GenerateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// BuildAuthorizationURL 构建Bangumi授权URL
func BuildAuthorizationURL(config *BangumiOAuthConfig, state string) string {
	params := url.Values{}
	params.Set("client_id", config.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", config.RedirectURI)
	params.Set("state", state)
	
	// 可以根据需要添加scope参数
	// params.Set("scope", "user_collection")
	
	return fmt.Sprintf("%s?%s", config.AuthURL, params.Encode())
}