package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("UserInfo端点测试")
	fmt.Println("================")
	
	// 测试端点信息
	fmt.Println("UserInfo端点地址:")
	fmt.Println("GET  http://localhost:8080/oauth/userinfo")
	fmt.Println("POST http://localhost:8080/oauth/userinfo")
	
	// 测试请求示例
	fmt.Println("\n测试请求示例:")
	fmt.Println("curl -H \"Authorization: Bearer <access_token>\" http://localhost:8080/oauth/userinfo")
	
	// 检查端点是否可访问
	resp, err := http.Get("http://localhost:8080/oauth/userinfo")
	if err != nil {
		fmt.Println("\n注意: OAuth服务器似乎未运行。")
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("\nUserInfo端点可访问。状态: %s\n", resp.Status)
	
	// 期望的响应结构
	fmt.Println("\n期望的响应结构:")
	fmt.Println("{")
	fmt.Println("  \"sub\": 123,")
	fmt.Println("  \"name\": \"testuser\",")
	fmt.Println("  \"email\": \"test@example.com\",")
	fmt.Println("  \"email_verified\": true,")
	fmt.Println("  \"nickname\": \"Test User\",")
	fmt.Println("  \"picture\": \"http://example.com/avatar.jpg\",")
	fmt.Println("  \"bio\": \"This is a test user\"")
	fmt.Println("}")
}