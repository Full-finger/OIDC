// cmd/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/golang-jwt/jwt/v5"

	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/middleware"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/service"
)

func main() {
	// 强制设置Gin为调试模式以显示详细日志
	gin.SetMode(gin.DebugMode)
	
	// 加载.env文件中的环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// 1. 初始化数据库连接池
	db := initDB()
	defer db.Close()

	// 2. 初始化各层依赖
	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	userService := service.NewUserService(userRepo)
	oauthService := service.NewOAuthService(oauthRepo, userRepo)
	
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(oauthService)
	tokenHandler := handler.NewTokenHandler(oauthService)
	userinfoHandler := handler.NewUserInfoHandler(userService, oauthService)
	discoveryHandler := handler.NewDiscoveryHandler()

	// 3. 设置 Gin 路由
	r := gin.Default()
	
	// 添加全局日志中间件
	r.Use(func(c *gin.Context) {
		log.Printf("=== Incoming Request ===")
		log.Printf("Method: %s", c.Request.Method)
		log.Printf("URL: %s", c.Request.URL.Path)
		log.Printf("Remote IP: %s", c.ClientIP())
		c.Next()
	})

	// OpenID Connect Discovery 端点
	r.GET("/.well-known/openid-configuration", discoveryHandler.GetDiscovery)

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
	
	// OAuth路由
	oauth := r.Group("/oauth")
	{
		oauth.GET("/authorize", authHandler.AuthorizeHandler)
		oauth.POST("/authorize", middleware.JWTAuthMiddleware(), authHandler.AuthorizePostHandler)
		oauth.POST("/token", tokenHandler.TokenHandler)
		// OIDC UserInfo端点
		oauth.GET("/userinfo", middleware.JWTAuthMiddleware(), userinfoHandler.GetUserInfo)
		oauth.POST("/userinfo", middleware.JWTAuthMiddleware(), userinfoHandler.GetUserInfo)
	}

	// 添加简单的登录页面路由
	r.GET("/login", func(c *gin.Context) {
		// 简单的登录页面HTML
		redirectParam := c.Query("redirect")
		// 转义引号以避免HTML问题
		redirectParam = strings.ReplaceAll(redirectParam, "\"", "&quot;")
		
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <meta charset="utf-8">
</head>
<body>
    <h2>Login</h2>
    <form method="POST" action="/login">
        <p>
            <label>Username:</label><br>
            <input type="text" name="username" required>
        </p>
        <p>
            <label>Password:</label><br>
            <input type="password" name="password" required>
        </p>
        <input type="hidden" name="redirect" value="` + redirectParam + `">
        <p>
            <button type="submit">Login</button>
        </p>
    </form>
</body>
</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// 添加登录处理路由
	r.POST("/login", func(c *gin.Context) {
		// 这里应该有实际的用户认证逻辑
		// 为了测试目的，我们简化处理
		
		// 模拟用户认证成功
		// 在实际应用中，你应该验证用户名和密码
		username := c.PostForm("username")
		_ = c.PostForm("password") // 忽略密码，仅为了消除编译错误
		redirect := c.PostForm("redirect")
		
		log.Printf("Login attempt - Username: %s, Redirect: %s", username, redirect)
		
		// 生成JWT token
		jwtSecret := getEnv("JWT_SECRET", "default_secret_key")
		
		// 创建token声明
		claims := &middleware.JWTClaims{
			UserID: 1, // 使用用户ID 1
		}
		
		// 创建token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Printf("Failed to generate JWT token: %v", err)
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}
		
		// 设置cookie
		c.SetCookie("token", tokenString, 3600, "/", "", false, true)
		
		// 重定向回OAuth流程
		if redirect == "" {
			// 如果没有提供redirect参数，使用默认值
			redirect = "/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:9999/callback&scope=read&state=xyz"
		}
		
		log.Printf("Login successful, redirecting to: %s", redirect)
		c.Redirect(302, redirect)
	})

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
	log.Printf("Gin mode: %s", gin.Mode())
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
	log.Println("Database connection established")

	return db
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}