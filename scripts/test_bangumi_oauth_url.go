package main

import (
	"fmt"
	"github.com/Full-finger/OIDC/internal/utils"
)

func main() {
	// 获取Bangumi OAuth配置
	config := utils.GetBangumiOAuthConfig()
	
	// 生成随机state
	state, err := utils.GenerateState()
	if err != nil {
		fmt.Printf("Error generating state: %v\n", err)
		return
	}
	
	// 构建授权URL
	authURL := utils.BuildAuthorizationURL(config, state)
	
	fmt.Printf("Bangumi OAuth Configuration:\n")
	fmt.Printf("  Client ID: %s\n", config.ClientID)
	fmt.Printf("  Client Secret: %s\n", config.ClientSecret)
	fmt.Printf("  Redirect URI: %s\n", config.RedirectURI)
	fmt.Printf("  Auth URL: %s\n", config.AuthURL)
	fmt.Printf("  Token URL: %s\n", config.TokenURL)
	fmt.Printf("  User Info URL: %s\n", config.UserInfoURL)
	fmt.Printf("\nGenerated Authorization URL:\n%s\n", authURL)
	fmt.Printf("State: %s\n", state)
}