package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8080"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

func main() {
	fmt.Println("开始测试OIDC客户端...")

	// 创建测试用户
	testUser := User{
		Username: "testuser_" + fmt.Sprintf("%d", time.Now().Unix()),
		Email:    "test_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	// 1. 用户注册
	fmt.Println("1. 测试用户注册...")
	if err := registerUser(testUser); err != nil {
		fmt.Printf("注册失败: %v\n", err)
		// 即使注册失败，我们仍继续测试其他功能
	} else {
		fmt.Println("用户注册成功")
	}

	// 2. 用户登录
	fmt.Println("2. 测试用户登录...")
	loginResp, err := loginUser(LoginRequest{
		Username: testUser.Username,
		Password: testUser.Password,
	})
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		// 使用模拟的令牌继续测试
		loginResp = &LoginResponse{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
		}
		fmt.Println("使用模拟令牌继续测试...")
	} else {
		fmt.Println("用户登录成功")
		fmt.Printf("访问令牌: %s\n", loginResp.AccessToken[:20]+"...")
		fmt.Printf("刷新令牌: %s\n", loginResp.RefreshToken[:20]+"...")
	}

	// 3. 测试OIDC Discovery端点
	fmt.Println("3. 测试OIDC Discovery端点...")
	if err := testDiscoveryEndpoint(); err != nil {
		fmt.Printf("Discovery端点测试失败: %v\n", err)
	} else {
		fmt.Println("OIDC Discovery端点测试成功")
	}

	// 4. 测试JWKS端点
	fmt.Println("4. 测试JWKS端点...")
	if err := testJWKSEndpoint(); err != nil {
		fmt.Printf("JWKS端点测试失败: %v\n", err)
	} else {
		fmt.Println("JWKS端点测试成功")
	}

	// 5. 测试授权端点
	fmt.Println("5. 测试授权端点...")
	if err := testAuthorizeEndpoint(loginResp.AccessToken); err != nil {
		fmt.Printf("授权端点测试失败: %v\n", err)
	} else {
		fmt.Println("授权端点测试成功")
	}

	// 6. 测试令牌端点
	fmt.Println("6. 测试令牌端点...")
	tokenResp, err := testTokenEndpoint()
	if err != nil {
		fmt.Printf("令牌端点测试失败: %v\n", err)
	} else {
		fmt.Println("令牌端点测试成功")
		fmt.Printf("访问令牌: %s\n", tokenResp.AccessToken[:20]+"...")
		if tokenResp.IDToken != "" {
			fmt.Printf("ID令牌: %s\n", tokenResp.IDToken[:20]+"...")
		}
	}

	// 7. 测试刷新令牌
	fmt.Println("7. 测试刷新令牌...")
	refreshResp, err := testRefreshToken(loginResp.RefreshToken)
	if err != nil {
		fmt.Printf("刷新令牌测试失败: %v\n", err)
	} else {
		fmt.Println("刷新令牌测试成功")
		fmt.Printf("新访问令牌: %s\n", refreshResp.AccessToken[:20]+"...")
	}

	// 8. 测试用户信息端点
	fmt.Println("8. 测试用户信息端点...")
	if err := testUserInfoEndpoint(loginResp.AccessToken); err != nil {
		fmt.Printf("用户信息端点测试失败: %v\n", err)
	} else {
		fmt.Println("用户信息端点测试成功")
	}

	fmt.Println("\n所有OIDC客户端测试完成！")
}

func registerUser(user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/api/v1/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("注册失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func loginUser(loginReq LoginRequest) (*LoginResponse, error) {
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("登录失败，状态码: %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func testDiscoveryEndpoint() error {
	resp, err := http.Get(baseURL + "/.well-known/openid-configuration")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Discovery端点测试失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func testJWKSEndpoint() error {
	resp, err := http.Get(baseURL + "/jwks.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS端点测试失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func testAuthorizeEndpoint(accessToken string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 阻止重定向
		},
	}

	req, err := http.NewRequest("GET", baseURL+"/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:3000/callback&scope=openid profile email&state=test_state", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 授权端点应该返回重定向响应
	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("授权端点测试失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func testTokenEndpoint() (*TokenResponse, error) {
	// 模拟授权码流程
	data := "grant_type=authorization_code&code=test_code&client_id=test_client&client_secret=test_secret&redirect_uri=http://localhost:3000/callback"

	resp, err := http.Post(baseURL+"/oauth/token", "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("令牌端点测试失败，状态码: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func testRefreshToken(refreshToken string) (*TokenResponse, error) {
	data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=test_client&client_secret=test_secret", refreshToken)

	resp, err := http.Post(baseURL+"/oauth/token", "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("刷新令牌测试失败，状态码: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func testUserInfoEndpoint(accessToken string) error {
	req, err := http.NewRequest("GET", baseURL+"/oauth/userinfo", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("用户信息端点测试失败，状态码: %d", resp.StatusCode)
	}

	return nil
}