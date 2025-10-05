package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 设置静态文件服务
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// 获取端口号，默认为3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 启动服务器
	fmt.Printf("Banggumoe server is running on http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}