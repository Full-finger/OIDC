package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Full-finger/OIDC/internal/middleware"
)

func main() {
	fmt.Println("UserInfo端点Scope测试")
	fmt.Println("====================")
	
	// 生成不同scope的测试令牌
	jwtSecret := "default_secret_key"
	
	// 1. 只有openid scope的令牌
	openidToken := generateTestToken(1, []string{"openid"}, jwtSecret)
	fmt.Printf("仅包含openid scope的访问令牌:\n%s\n\n", openidToken)
	
	// 2. 包含profile scope的令牌
	profileToken := generateTestToken(1, []string{"openid", "profile"}, jwtSecret)
	fmt.Printf("包含profile scope的访问令牌:\n%s\n\n", profileToken)
	
	// 3. 包含email scope的令牌
	emailToken := generateTestToken(1, []string{"openid", "email"}, jwtSecret)
	fmt.Printf("包含email scope的访问令牌:\n%s\n\n", emailToken)
	
	// 4. 包含所有scope的令牌
	fullToken := generateTestToken(1, []string{"openid", "profile", "email"}, jwtSecret)
	fmt.Printf("包含所有scope的访问令牌:\n%s\n\n", fullToken)
	
	// 测试请求示例
	fmt.Println("测试请求示例:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8080/oauth/userinfo\n\n", fullToken)
	
	fmt.Println("预期响应:")
	fmt.Println("1. 仅openid scope: {\"sub\": 1}")
	fmt.Println("2. 包含profile scope: {\"sub\": 1, \"name\": \"testuser\", \"nickname\": \"Test User\", \"picture\": \"http://example.com/avatar.jpg\"}")
	fmt.Println("3. 包含email scope: {\"sub\": 1, \"email\": \"test@example.com\", \"email_verified\": true}")
	fmt.Println("4. 包含所有scope: {\"sub\": 1, \"name\": \"testuser\", \"nickname\": \"Test User\", \"picture\": \"http://example.com/avatar.jpg\", \"email\": \"test@example.com\", \"email_verified\": true}")
}

func generateTestToken(userID int64, scopes []string, secret string) string {
	// 创建token声明
	claims := &middleware.JWTClaims{
		UserID: userID,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("生成令牌失败: %v\n", err)
		return ""
	}
	
	return tokenString
}