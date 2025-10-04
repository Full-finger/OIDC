package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("OpenID Connect Discovery端点完整测试")
	fmt.Println("===================================")
	
	// 测试端点信息
	fmt.Println("Discovery端点地址:")
	fmt.Println("GET http://localhost:8080/.well-known/openid-configuration")
	
	// 测试请求示例
	fmt.Println("\n测试请求示例:")
	fmt.Println("curl http://localhost:8080/.well-known/openid-configuration")
	
	// 检查端点是否可访问
	resp, err := http.Get("http://localhost:8080/.well-known/openid-configuration")
	if err != nil {
		fmt.Println("\n注意: OIDC服务器似乎未运行。")
		fmt.Println("请启动服务器后再运行此测试:")
		fmt.Println("go run cmd/main.go")
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("\nDiscovery端点可访问。状态: %s\n", resp.Status)
	
	// 解析响应
	var discovery map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}
	
	// 输出关键字段
	fmt.Println("\n关键配置信息:")
	if issuer, ok := discovery["issuer"]; ok {
		fmt.Printf("Issuer: %s\n", issuer)
	}
	
	if authEndpoint, ok := discovery["authorization_endpoint"]; ok {
		fmt.Printf("Authorization Endpoint: %s\n", authEndpoint)
	}
	
	if tokenEndpoint, ok := discovery["token_endpoint"]; ok {
		fmt.Printf("Token Endpoint: %s\n", tokenEndpoint)
	}
	
	if userinfoEndpoint, ok := discovery["userinfo_endpoint"]; ok {
		fmt.Printf("Userinfo Endpoint: %s\n", userinfoEndpoint)
	}
	
	if jwksURI, ok := discovery["jwks_uri"]; ok {
		fmt.Printf("JWKS URI: %s\n", jwksURI)
	}
	
	// 检查必需字段
	requiredFields := []string{
		"issuer",
		"authorization_endpoint",
		"token_endpoint",
		"userinfo_endpoint",
		"jwks_uri",
		"scopes_supported",
		"response_types_supported",
		"subject_types_supported",
		"id_token_signing_alg_values_supported",
		"claims_supported",
	}
	
	fmt.Println("\n必需字段检查:")
	allPresent := true
	for _, field := range requiredFields {
		if _, ok := discovery[field]; !ok {
			fmt.Printf("❌ 缺少字段: %s\n", field)
			allPresent = false
		} else {
			fmt.Printf("✅ 字段存在: %s\n", field)
		}
	}
	
	if allPresent {
		fmt.Println("\n✅ 所有必需字段都已正确提供")
	} else {
		fmt.Println("\n❌ 一些必需字段缺失")
		os.Exit(1)
	}
	
	// 显示完整响应结构示例
	fmt.Println("\n完整响应结构示例:")
	fmt.Println("{")
	fmt.Println("  \"issuer\": \"http://localhost:8080\",")
	fmt.Println("  \"authorization_endpoint\": \"http://localhost:8080/oauth/authorize\",")
	fmt.Println("  \"token_endpoint\": \"http://localhost:8080/oauth/token\",")
	fmt.Println("  \"userinfo_endpoint\": \"http://localhost:8080/oauth/userinfo\",")
	fmt.Println("  \"jwks_uri\": \"http://localhost:8080/.well-known/jwks.json\",")
	fmt.Println("  \"scopes_supported\": [\"openid\", \"profile\", \"email\"],")
	fmt.Println("  \"response_types_supported\": [\"code\"],")
	fmt.Println("  \"response_modes_supported\": [\"query\"],")
	fmt.Println("  \"grant_types_supported\": [\"authorization_code\", \"refresh_token\"],")
	fmt.Println("  \"subject_types_supported\": [\"public\"],")
	fmt.Println("  \"id_token_signing_alg_values_supported\": [\"HS256\"],")
	fmt.Println("  \"token_endpoint_auth_methods_supported\": [\"client_secret_basic\", \"client_secret_post\"],")
	fmt.Println("  \"claims_supported\": [")
	fmt.Println("    \"sub\", \"iss\", \"aud\", \"exp\", \"iat\", \"auth_time\",")
	fmt.Println("    \"nonce\", \"name\", \"nickname\", \"preferred_username\",")
	fmt.Println("    \"picture\", \"email\", \"email_verified\"")
	fmt.Println("  ],")
	fmt.Println("  \"code_challenge_methods_supported\": [\"S256\"],")
	fmt.Println("  \"end_session_endpoint\": \"http://localhost:8080/oauth/logout\"")
	fmt.Println("}")
}