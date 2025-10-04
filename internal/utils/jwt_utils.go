// Package utils provides utility functions for the OIDC application.
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userModel "github.com/Full-finger/OIDC/config"
	oauthModel "github.com/Full-finger/OIDC/internal/model"
)

// GenerateIDToken creates an OIDC ID Token for the given user, client, and nonce.
// It uses the golang-jwt/jwt/v5 library to create a signed JWT with all the required claims.
func GenerateIDToken(user *userModel.User, client *oauthModel.Client, nonce *string, scopes []string) (string, error) {
	// 从环境变量获取JWT密钥，如果没有则使用默认值
	jwtSecret := getEnv("JWT_SECRET", "default_secret_key")
	
	// 构建ID token claims
	claims := jwt.MapClaims{
		"iss": "http://localhost:8080", // TODO: 从配置中获取
		"sub": fmt.Sprintf("%d", user.ID),
		"aud": client.ClientID,
		"exp": time.Now().Add(time.Hour).Unix(), // 1小时过期
		"iat": time.Now().Unix(),
		"auth_time": time.Now().Unix(),
	}
	
	// 如果提供了nonce，则添加到claims中
	if nonce != nil {
		claims["nonce"] = *nonce
	}
	
	// 添加基于scope的claims
	for _, scope := range scopes {
		switch scope {
		case "profile":
			// 添加profile相关的claims
			claims["name"] = user.Username
			if user.Nickname != nil {
				claims["nickname"] = *user.Nickname
			}
			if user.AvatarURL != nil {
				claims["picture"] = *user.AvatarURL
			}
			// 可以添加更多profile claims
			claims["preferred_username"] = user.Username
		case "email":
			// 添加email相关的claims
			claims["email"] = user.Email
			claims["email_verified"] = user.EmailVerified
		}
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名token
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign ID token: %w", err)
	}
	
	return signedToken, nil
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}