// cmd/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/middleware"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/service"
)

func main() {
	// 加载.env文件中的环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// 1. 初始化数据库连接池
	db := initDB()
	defer db.Close()

	// 2. 初始化各层依赖
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 3. 设置 Gin 路由
	r := gin.Default()

	// 公共路由（无需认证）
	public := r.Group("/api/v1")
	{
		public.POST("/register", func(c *gin.Context) {
			userHandler.RegisterHandler(c.Writer, c.Request)
		})
		public.POST("/login", func(c *gin.Context) {
			userHandler.LoginHandler(c.Writer, c.Request)
		})
	}

	// 受保护路由（需认证）
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfileHandler)
		protected.PUT("/profile", userHandler.UpdateProfileHandler)
	}

	// 4. 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB 初始化 PostgreSQL 连接池
func initDB() *sql.DB {
	// 从环境变量读取配置参数，或使用默认值（开发用）
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "oidc_user")
	password := getEnv("DB_PASSWORD", "oidc_password")
	dbname := getEnv("DB_NAME", "oidc_db")
	port := getEnv("DB_PORT", "5432")
	sslmode := getEnv("DB_SSLMODE", "disable")
	timezone := getEnv("DB_TIMEZONE", "UTC")

	// 动态构建连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}