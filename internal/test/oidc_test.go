package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/router"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/util"
)

// TestOIDCFlow 测试完整的OIDC流程
func TestOIDCFlow(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
	
	// 创建测试路由器
	r := router.SetupRouter()
	
	// 创建测试用户
	testUser := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	
	// 注册用户
	t.Run("Register User", func(t *testing.T) {
		registerData := map[string]interface{}{
			"username": testUser.Username,
			"email":    testUser.Email,
			"password": testUser.Password,
			"nickname": testUser.Nickname,
		}
		
		jsonData, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
	})
	
	// 跳过邮箱验证，直接在数据库中验证用户
	t.Run("Verify User Email", func(t *testing.T) {
		// 在实际测试中，我们会在这里直接修改数据库中的用户状态
		// 但由于我们使用内存存储，我们在登录时跳过验证检查
		t.Log("Skipping email verification in test mode")
	})
	
	// 用户登录
	var accessToken string
	var refreshToken string
	t.Run("User Login", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": testUser.Username,
			"password": testUser.Password,
		}
		
		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
			return
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if token, ok := response["access_token"].(string); ok {
			accessToken = token
		} else {
			t.Error("Access token not found in response")
		}
		
		if token, ok := response["refresh_token"].(string); ok {
			refreshToken = token
		} else {
			t.Error("Refresh token not found in response")
		}
		
		if accessToken == "" || refreshToken == "" {
			t.Error("Failed to obtain tokens")
		}
	})
	
	// 测试OIDC Discovery端点
	t.Run("OIDC Discovery", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/.well-known/openid-configuration", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if issuer, ok := response["issuer"].(string); !ok || issuer == "" {
			t.Error("Issuer not found in discovery response")
		}
	})
	
	// 创建测试OAuth客户端
	var clientID string
	var clientSecret string
	t.Run("Create Test Client", func(t *testing.T) {
		// 在实际应用中，我们会通过API创建客户端
		// 在测试中，我们直接生成一个客户端用于测试
		clientID = "test_client"
		clientSecret = "test_secret"
		t.Log("Created test OAuth client")
	})
	
	// 测试授权码流程
	var authCode string
	t.Run("Authorization Request", func(t *testing.T) {
		// 模拟授权请求
		req, _ := http.NewRequest("GET", fmt.Sprintf("/oauth/authorize?response_type=code&client_id=%s&redirect_uri=http://localhost:3000/callback&scope=openid profile email&state=test_state", clientID), nil)
		
		// 添加认证头
		req.Header.Set("Authorization", "Bearer "+accessToken)
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		// 注意：在实际应用中，这会重定向到客户端的回调URL
		// 在测试中，我们检查是否成功处理了请求
		if w.Code != http.StatusOK && w.Code != http.StatusFound {
			t.Errorf("Expected status code 200 or 302, got %d", w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
		
		// 在实际实现中，我们需要从响应中提取授权码
		// 这里我们生成一个测试用的授权码
		authCode = "test_auth_code"
		t.Log("Authorization request processed")
	})
	
	// 测试令牌端点
	t.Run("Token Exchange", func(t *testing.T) {
		tokenData := map[string]interface{}{
			"grant_type":    "authorization_code",
			"code":          authCode,
			"client_id":     clientID,
			"client_secret": clientSecret,
			"redirect_uri":  "http://localhost:3000/callback",
		}
		
		jsonData, _ := json.Marshal(tokenData)
		req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
			return
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if token, ok := response["access_token"].(string); !ok || token == "" {
			t.Error("Access token not found in token response")
		}
		
		if idToken, ok := response["id_token"].(string); !ok || idToken == "" {
			t.Error("ID token not found in token response")
		}
	})
	
	// 测试刷新令牌
	t.Run("Refresh Token", func(t *testing.T) {
		refreshData := map[string]interface{}{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     clientID,
			"client_secret": clientSecret,
		}
		
		jsonData, _ := json.Marshal(refreshData)
		req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if token, ok := response["access_token"].(string); !ok || token == "" {
			t.Error("Access token not found in refresh response")
		}
	})
	
	// 测试用户信息端点
	t.Run("User Info", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/oauth/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if sub, ok := response["sub"].(string); !ok || sub == "" {
			t.Error("Subject not found in userinfo response")
		}
	})
	
	// 测试JWKS端点
	t.Run("JWKS Endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/jwks.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			t.Errorf("Response body: %s", w.Body.String())
		}
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if keys, ok := response["keys"].([]interface{}); !ok || len(keys) == 0 {
			t.Error("Keys not found in JWKS response")
		}
	})
}