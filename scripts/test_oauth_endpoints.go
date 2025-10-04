package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 测试OAuth令牌端点支持的grant_type
	fmt.Println("OAuth Token Endpoint Grant Types:")
	fmt.Println("1. authorization_code")
	fmt.Println("2. refresh_token")
	fmt.Println("\nTo test these endpoints, you can use curl commands like:")
	fmt.Println("\n# Test authorization_code grant type:")
	fmt.Println(`curl -X POST http://localhost:8080/oauth/token \`)
	fmt.Println(`  -H "Content-Type: application/x-www-form-urlencoded" \`)
	fmt.Println(`  -d "grant_type=authorization_code&code=AUTH_CODE&redirect_uri=REDIRECT_URI&client_id=CLIENT_ID&client_secret=CLIENT_SECRET"`)
	
	fmt.Println("\n# Test refresh_token grant type:")
	fmt.Println(`curl -X POST http://localhost:8080/oauth/token \`)
	fmt.Println(`  -H "Content-Type: application/x-www-form-urlencoded" \`)
	fmt.Println(`  -d "grant_type=refresh_token&refresh_token=REFRESH_TOKEN&client_id=CLIENT_ID&client_secret=CLIENT_SECRET"`)
	
	// 检查端点是否可访问
	resp, err := http.Get("http://localhost:8080/oauth/token")
	if err != nil {
		fmt.Println("\nNote: OAuth server does not appear to be running.")
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("\nOAuth token endpoint is accessible. Status: %s\n", resp.Status)
}