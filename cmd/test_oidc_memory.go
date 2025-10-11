package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/router"
)

func main() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由器
	r := router.SetupRouter()

	fmt.Println("开始测试OIDC流程（使用内存存储）...")

	// 创建测试用户数据
	testUser := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
		"nickname": "Test User",
	}

	// 1. 用户注册
	fmt.Println("1. 测试用户注册...")
	registerData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(registerData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("注册响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("注册失败，响应内容: %s\n", w.Body.String())
		// 即使注册失败，我们仍继续测试其他功能
	} else {
		fmt.Println("用户注册成功")
	}

	// 2. 用户登录
	fmt.Println("2. 测试用户登录...")
	loginData := map[string]interface{}{
		"username": testUser["username"],
		"password": testUser["password"],
	}

	loginJSON, _ := json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("登录响应状态码: %d\n", w.Code)
	var accessToken string
	var refreshToken string

	if w.Code != http.StatusOK {
		fmt.Printf("登录失败，响应内容: %s\n", w.Body.String())
		// 使用模拟的令牌继续测试
		accessToken = "mock_access_token"
		refreshToken = "mock_refresh_token"
		fmt.Println("使用模拟令牌继续测试...")
	} else {
		var loginResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &loginResponse)
		
		accessToken, _ = loginResponse["access_token"].(string)
		refreshToken, _ = loginResponse["refresh_token"].(string)
		
		fmt.Println("用户登录成功")
		fmt.Printf("访问令牌: %s\n", accessToken[:10]+"...")
		fmt.Printf("刷新令牌: %s\n", refreshToken[:10]+"...")
	}

	// 3. 测试OIDC Discovery端点
	fmt.Println("3. 测试OIDC Discovery端点...")
	req, _ = http.NewRequest("GET", "/.well-known/openid-configuration", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("Discovery端点响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("Discovery端点测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("OIDC Discovery端点测试成功")
	}

	// 4. 测试JWKS端点
	fmt.Println("4. 测试JWKS端点...")
	req, _ = http.NewRequest("GET", "/jwks.json", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("JWKS端点响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("JWKS端点测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("JWKS端点测试成功")
	}

	// 5. 测试授权端点 (模拟)
	fmt.Println("5. 测试授权端点...")
	authURL := "/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:3000/callback&scope=openid profile email&state=test_state"
	req, _ = http.NewRequest("GET", authURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("授权端点响应状态码: %d\n", w.Code)
	// 注意：在实际应用中，这会重定向到客户端的回调URL
	if w.Code != http.StatusOK && w.Code != http.StatusFound {
		fmt.Printf("授权端点测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("授权端点测试成功")
	}

	// 6. 测试令牌端点 (模拟)
	fmt.Println("6. 测试令牌端点...")
	tokenData := "grant_type=authorization_code&code=test_code&client_id=test_client&client_secret=test_secret&redirect_uri=http://localhost:3000/callback"
	req, _ = http.NewRequest("POST", "/oauth/token", bytes.NewBufferString(tokenData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("令牌端点响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("令牌端点测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("令牌端点测试成功")
		var tokenResponse map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &tokenResponse); err == nil {
			if idToken, ok := tokenResponse["id_token"].(string); ok {
				fmt.Printf("ID令牌获取成功: %s\n", idToken[:10]+"...")
			}
			if accessToken, ok := tokenResponse["access_token"].(string); ok {
				fmt.Printf("访问令牌获取成功: %s\n", accessToken[:10]+"...")
			}
		}
	}

	// 7. 测试刷新令牌
	fmt.Println("7. 测试刷新令牌...")
	refreshData := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=test_client&client_secret=test_secret", refreshToken)
	req, _ = http.NewRequest("POST", "/oauth/token", bytes.NewBufferString(refreshData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("刷新令牌响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("刷新令牌测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("刷新令牌测试成功")
		var refreshResponse map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &refreshResponse); err == nil {
			if newAccessToken, ok := refreshResponse["access_token"].(string); ok {
				fmt.Printf("新访问令牌获取成功: %s\n", newAccessToken[:10]+"...")
			}
		}
	}

	// 8. 测试用户信息端点
	fmt.Println("8. 测试用户信息端点...")
	req, _ = http.NewRequest("GET", "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("用户信息端点响应状态码: %d\n", w.Code)
	if w.Code != http.StatusOK {
		fmt.Printf("用户信息端点测试失败，响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("用户信息端点测试成功")
		var userInfo map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &userInfo); err == nil {
			fmt.Printf("用户信息: %+v\n", userInfo)
		}
	}

	fmt.Println("\n所有OIDC测试完成！")
}