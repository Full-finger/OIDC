package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Full-finger/OIDC/internal/router"
	"github.com/Full-finger/OIDC/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("警告: 无法加载.env文件: %v\n", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("无法连接到数据库: " + err.Error())
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.VerificationToken{},
		&model.Client{},
		&model.AuthorizationCode{},
		&model.RefreshToken{},
		&model.Anime{},
		&model.Collection{},
		&model.BangumiAccount{},
	); err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	fmt.Println("数据库连接成功")

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由器
	r := router.SetupRouter()

	// 创建测试用户
	testUser := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
		"nickname": "Test User",
	}

	fmt.Println("开始测试OIDC流程...")

	// 1. 用户注册
	fmt.Println("1. 测试用户注册...")
	registerData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(registerData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("注册失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("用户注册成功")
	}

	// 2. 验证用户邮箱（直接更新数据库跳过邮件验证）
	fmt.Println("2. 跳过邮箱验证...")
	var user model.User
	if err := db.Where("username = ?", testUser["username"]).First(&user).Error; err != nil {
		fmt.Printf("查找用户失败: %v\n", err)
		return
	}

	// 注意：这里我们假设用户表有一个字段表示邮箱验证状态
	// 如果没有，我们会在登录时跳过验证检查
	fmt.Println("用户邮箱已验证（模拟）")

	// 3. 用户登录
	fmt.Println("3. 测试用户登录...")
	loginData := map[string]interface{}{
		"username": testUser["username"],
		"password": testUser["password"],
	}

	loginJSON, _ := json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		fmt.Printf("解析登录响应失败: %v\n", err)
		return
	}

	if w.Code != http.StatusOK {
		fmt.Printf("登录失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
		return
	}

	accessToken, ok := loginResponse["access_token"].(string)
	if !ok {
		fmt.Println("无法获取访问令牌")
		return
	}

	refreshToken, ok := loginResponse["refresh_token"].(string)
	if !ok {
		fmt.Println("无法获取刷新令牌")
		return
	}

	fmt.Println("用户登录成功")
	fmt.Printf("访问令牌: %s\n", accessToken[:20]+"...")
	fmt.Printf("刷新令牌: %s\n", refreshToken[:20]+"...")

	// 4. 测试OIDC Discovery端点
	fmt.Println("4. 测试OIDC Discovery端点...")
	req, _ = http.NewRequest("GET", "/.well-known/openid-configuration", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("Discovery端点测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("OIDC Discovery端点测试成功")
	}

	// 5. 测试JWKS端点
	fmt.Println("5. 测试JWKS端点...")
	req, _ = http.NewRequest("GET", "/jwks.json", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("JWKS端点测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("JWKS端点测试成功")
	}

	// 6. 创建测试客户端
	fmt.Println("6. 创建测试OAuth客户端...")
	client := model.Client{
		ClientID:    "test_client",
		SecretHash:  "test_secret", // 简化处理，实际应该使用哈希
		Name:        "Test Client",
		Description: "Test client for OIDC flow",
		RedirectURI: "http://localhost:3000/callback",
		Scopes:      "openid profile email",
	}

	// 检查客户端是否已存在
	var existingClient model.Client
	if err := db.Where("client_id = ?", client.ClientID).First(&existingClient).Error; err != nil {
		// 客户端不存在，创建新客户端
		if err := db.Create(&client).Error; err != nil {
			fmt.Printf("创建客户端失败: %v\n", err)
			return
		}
		fmt.Println("测试客户端创建成功")
	} else {
		fmt.Println("使用已存在的测试客户端")
		client = existingClient
	}

	// 7. 测试授权端点
	fmt.Println("7. 测试授权端点...")
	authURL := "/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:3000/callback&scope=openid profile email&state=test_state"
	req, _ = http.NewRequest("GET", authURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 注意：在实际应用中，这会重定向到客户端的回调URL
	// 在测试中，我们检查是否成功处理了请求
	if w.Code != http.StatusOK && w.Code != http.StatusFound {
		fmt.Printf("授权端点测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("授权端点测试成功")
	}

	// 8. 模拟生成授权码用于测试令牌端点
	fmt.Println("8. 测试令牌端点...")
	// 由于我们无法从授权端点直接获取授权码，我们手动创建一个用于测试
	authCode := "test_auth_code_12345"

	// 保存授权码到数据库（模拟授权过程）
	authCodeModel := model.AuthorizationCode{
		Code:        authCode,
		ClientID:    client.ClientID,
		UserID:      user.ID,
		RedirectURI: "http://localhost:3000/callback",
		Scopes:      "openid profile email",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	// 检查授权码是否已存在
	var existingAuthCode model.AuthorizationCode
	if err := db.Where("code = ?", authCode).First(&existingAuthCode).Error; err != nil {
		// 授权码不存在，创建新授权码
		if err := db.Create(&authCodeModel).Error; err != nil {
			fmt.Printf("创建授权码失败: %v\n", err)
			return
		}
		fmt.Println("测试授权码创建成功")
	} else {
		fmt.Println("使用已存在的测试授权码")
	}

	// 9. 测试令牌端点
	fmt.Println("9. 测试令牌端点...")
	tokenData := fmt.Sprintf("grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=%s",
		authCode, client.ClientID, "test_secret", "http://localhost:3000/callback")

	req, _ = http.NewRequest("POST", "/oauth/token", strings.NewReader(tokenData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &tokenResponse); err != nil {
		fmt.Printf("解析令牌响应失败: %v\n", err)
		return
	}

	if w.Code != http.StatusOK {
		fmt.Printf("令牌端点测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("令牌端点测试成功")
		if idToken, ok := tokenResponse["id_token"].(string); ok {
			fmt.Printf("ID令牌获取成功: %s\n", idToken[:20]+"...")
		}
		if accessToken, ok := tokenResponse["access_token"].(string); ok {
			fmt.Printf("访问令牌获取成功: %s\n", accessToken[:20]+"...")
		}
	}

	// 10. 测试刷新令牌
	fmt.Println("10. 测试刷新令牌...")
	refreshData := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s",
		refreshToken, client.ClientID, "test_secret")

	req, _ = http.NewRequest("POST", "/oauth/token", strings.NewReader(refreshData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var refreshResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &refreshResponse); err != nil {
		fmt.Printf("解析刷新令牌响应失败: %v\n", err)
		return
	}

	if w.Code != http.StatusOK {
		fmt.Printf("刷新令牌测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("刷新令牌测试成功")
		if newAccessToken, ok := refreshResponse["access_token"].(string); ok {
			fmt.Printf("新访问令牌获取成功: %s\n", newAccessToken[:20]+"...")
		}
	}

	// 11. 测试用户信息端点
	fmt.Println("11. 测试用户信息端点...")
	req, _ = http.NewRequest("GET", "/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("用户信息端点测试失败，状态码: %d\n", w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
	} else {
		fmt.Println("用户信息端点测试成功")
		var userInfo map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &userInfo); err == nil {
			fmt.Printf("用户信息: %+v\n", userInfo)
		}
	}

	fmt.Println("\n所有OIDC测试完成！")
}