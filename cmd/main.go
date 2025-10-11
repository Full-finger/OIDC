package main

import (
	"log"
	"os"

	"github.com/Full-finger/OIDC/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 无法加载 .env 文件")
	}

	// 设置端口，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务器
	log.Printf("服务器启动在端口 %s", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}