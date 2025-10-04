package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	// 数据库连接配置
	host := "172.24.192.125"
	user := "oidc_user"
	password := "oidc_password"
	dbname := "oidc_db"
	port := "5432"
	sslmode := "disable"

	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// 连接数据库
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 创建测试客户端
	clientID := "test_client"
	clientSecret := "test_secret"
	clientName := "Test OAuth Client"
	
	// 将字符串切片转换为PostgreSQL数组格式
	redirectURIs := []string{"http://localhost:9999/callback", "https://example.com/callback"}
	scopes := []string{"read", "write"}

	// 删除已存在的测试客户端
	_, err = db.Exec("DELETE FROM oauth_clients WHERE client_id = $1", clientID)
	if err != nil {
		log.Printf("Warning: failed to delete existing test client: %v", err)
	}

	// 插入测试客户端
	const insertClientQuery = `
		INSERT INTO oauth_clients (
			client_id, client_secret_hash, name, redirect_uris, scopes
		) VALUES (
			$1, $2, $3, $4, $5
		)`

	// 简单哈希客户端密钥（实际应用中应使用更安全的方法）
	clientSecretHash := hashSecret(clientSecret)

	// 使用pq.Array包装数组参数
	_, err = db.Exec(insertClientQuery, clientID, clientSecretHash, clientName, 
		pq.Array(redirectURIs), pq.Array(scopes))
	if err != nil {
		log.Fatalf("Failed to insert test client: %v", err)
	}

	fmt.Println("Test OAuth client created successfully!")
	fmt.Printf("Client ID: %s\n", clientID)
	fmt.Printf("Client Secret: %s\n", clientSecret)
	fmt.Printf("Redirect URIs: %v\n", redirectURIs)
	fmt.Printf("Scopes: %v\n", scopes)
}

// hashSecret 对密钥进行哈希处理（与oauth_service.go中的一致）
func hashSecret(secret string) string {
	// 注意：这是一个简化的实现，实际项目中应该使用bcrypt等安全的哈希算法
	// 这里为了测试方便，使用简单的SHA256哈希
	hash := sha256.Sum256([]byte(secret))
	return base64.StdEncoding.EncodeToString(hash[:])
}