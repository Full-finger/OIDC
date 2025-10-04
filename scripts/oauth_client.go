package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	authServerURL     = "http://localhost:8080"
	clientID          = "test_client"
	clientSecret      = "test_secret"
	clientRedirectURI = "http://localhost:9999/callback"
)

var (
	client *http.Client
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	client = &http.Client{Timeout: 10 * time.Second}

	// 启动一个简单的Web服务器来模拟客户端应用
	router := gin.Default()

	// 首页，包含一个"登录"链接
	router.GET("/", func(c *gin.Context) {
		html := `<h1>OAuth2.0 Client Test</h1>
		<p>This is a test client for the OIDC service.</p>
		<a href="/login">Login with OIDC Service</a>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// 登录入口，重定向到授权服务器
	router.GET("/login", handleLogin)

	// 授权服务器的回调地址
	router.GET("/callback", handleCallback)

	// 使用Token访问受保护资源
	router.GET("/profile", handleProfile)

	log.Println("Client app starting on http://localhost:9999...")
	log.Println("Press Ctrl+C to exit")

	go func() {
		if err := router.Run(":9999"); err != nil {
			log.Fatalf("Could not start client server: %s\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Client server shutting down...")
}

func handleLogin(c *gin.Context) {
	// 构造授权URL
	authURL := fmt.Sprintf("%s/oauth/authorize?response_type=code&client_id=%s&scope=read&state=xyz&redirect_uri=%s",
		authServerURL, clientID, url.QueryEscape(clientRedirectURI))

	// 重定向用户到授权服务器
	c.Redirect(http.StatusFound, authURL)
}

func handleCallback(c *gin.Context) {
	// 从查询参数中获取授权码和state
	code := c.Query("code")
	state := c.Query("state")
	if code == "" {
		c.String(http.StatusBadRequest, "Error: Authorization code not found.")
		return
	}

	log.Printf("Received authorization code: %s, state: %s", code, state)

	// 使用授权码换取访问令牌
	token, err := exchangeCodeForToken(code)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error exchanging code for token: %v", err))
		return
	}

	log.Printf("Received access token: %s", token.AccessToken)

	// 将Token存储在session中（这里简化为URL参数）
	c.Redirect(http.StatusFound, "/profile?access_token="+token.AccessToken)
}

func handleProfile(c *gin.Context) {
	accessToken := c.Query("access_token")
	if accessToken == "" {
		c.String(http.StatusBadRequest, "Error: Access token not found.")
		return
	}

	// 使用访问令牌请求用户信息
	// 注意：我们的API目前还没有受保护的资源端点，这是一个下一步的提示
	// 假设我们有一个 /api/v1/profile 端点
	profileURL := fmt.Sprintf("%s/api/v1/profile", authServerURL)
	req, _ := http.NewRequest("GET", profileURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching profile: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(resp.StatusCode, fmt.Sprintf("Error: Failed to get profile. Status: %s", resp.Status))
		return
	}

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		c.String(http.StatusInternalServerError, "Error decoding profile response.")
		return
	}

	// 格式化输出用户信息
	profileJSON, _ := json.MarshalIndent(profile, "", "  ")

	html := fmt.Sprintf(`
	<h1>Protected Resource Access</h1>
	<p>Successfully accessed protected resource using OAuth2.0 access token!</p>
	<h2>User Profile:</h2>
	<pre>%s</pre>
	<a href="/">Back to Home</a>
	`, string(profileJSON))

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", clientRedirectURI)

	req, err := http.NewRequest("POST", authServerURL+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}