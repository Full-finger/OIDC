package main

import (
	"fmt"
	"time"

	userModel "github.com/Full-finger/OIDC/config"
	oauthModel "github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/utils"
)

func main() {
	fmt.Println("测试ID Token生成函数")
	
	// 创建测试用户
	user := &userModel.User{
		ID:            1,
		Username:      "testuser",
		Email:         "test@example.com",
		Nickname:      stringPtr("Test User"),
		AvatarURL:     stringPtr("http://example.com/avatar.jpg"),
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	// 创建测试客户端
	client := &oauthModel.Client{
		ClientID: "test_client",
		Name:     "Test Client",
	}
	
	// 创建测试scopes
	scopes := []string{"openid", "profile", "email"}
	
	// 创建nonce
	nonce := stringPtr("test_nonce")
	
	// 生成ID Token
	idToken, err := utils.GenerateIDToken(user, client, nonce, scopes)
	if err != nil {
		fmt.Printf("生成ID Token失败: %v\n", err)
		return
	}
	
	fmt.Printf("成功生成ID Token:\n%s\n", idToken)
	fmt.Println("ID Token生成函数工作正常!")
}

func stringPtr(s string) *string {
	return &s
}