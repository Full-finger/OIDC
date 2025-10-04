package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	fmt.Println("OIDC测试脚本")
	fmt.Println("============")
	
	// 测试OAuth令牌端点支持的grant_type
	fmt.Println("OAuth Token Endpoint Grant Types:")
	fmt.Println("1. authorization_code")
	fmt.Println("2. refresh_token")
	
	// 测试OIDC授权请求
	fmt.Println("\nOIDC授权请求示例:")
	fmt.Println("GET /oauth/authorize?" + url.Values{
		"response_type": {"code"},
		"client_id":     {"test_client"},
		"redirect_uri":  {"http://localhost:9999/callback"},
		"scope":         {"openid profile email"},
		"state":         {"xyz"},
	}.Encode())
	
	// 测试令牌交换请求
	fmt.Println("\n令牌交换请求示例:")
	fmt.Println("POST /oauth/token")
	fmt.Println("Content-Type: application/x-www-form-urlencoded")
	fmt.Println("Body:")
	fmt.Println("  grant_type=authorization_code")
	fmt.Println("  code=AUTH_CODE")
	fmt.Println("  redirect_uri=http://localhost:9999/callback")
	fmt.Println("  client_id=test_client")
	fmt.Println("  client_secret=CLIENT_SECRET")
	
	// 检查端点是否可访问
	resp, err := http.Get("http://localhost:8080/oauth/token")
	if err != nil {
		fmt.Println("\n注意: OAuth服务器似乎未运行。")
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("\nOAuth令牌端点可访问。状态: %s\n", resp.Status)
}