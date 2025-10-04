package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("OpenID Connect Discovery端点测试")
	fmt.Println("===============================")
	
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
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("\nDiscovery端点可访问。状态: %s\n", resp.Status)
	
	// 期望的响应结构
	fmt.Println("\n期望的响应结构:")
	fmt.Println("{")
	fmt.Println("  \"issuer\": \"http://localhost:8080\",")
	fmt.Println("  \"authorization_endpoint\": \"http://localhost:8080/oauth/authorize\",")
	fmt.Println("  \"token_endpoint\": \"http://localhost:8080/oauth/token\",")
	fmt.Println("  \"userinfo_endpoint\": \"http://localhost:8080/oauth/userinfo\",")
	fmt.Println("  \"jwks_uri\": \"http://localhost:8080/.well-known/jwks.json\",")
	fmt.Println("  \"scopes_supported\": [\"openid\", \"profile\", \"email\"],")
	fmt.Println("  \"response_types_supported\": [\"code\"],")
	fmt.Println("  \"grant_types_supported\": [\"authorization_code\", \"refresh_token\"],")
	fmt.Println("  \"subject_types_supported\": [\"public\"],")
	fmt.Println("  \"id_token_signing_alg_values_supported\": [\"HS256\"],")
	fmt.Println("  \"claims_supported\": [\"sub\", \"iss\", \"aud\", \"exp\", \"iat\", \"auth_time\", \"nonce\", \"name\", \"nickname\", \"preferred_username\", \"picture\", \"email\", \"email_verified\"]")
	fmt.Println("}")
}