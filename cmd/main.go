// cmd/main.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Full-finger/OIDC/internal/aspect"
	"github.com/Full-finger/OIDC/internal/filter"
	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize components
	loggingAspect := aspect.NewLoggingAspect()
	preprocessingAspect := aspect.NewPreprocessingAspect()

	// Initialize helpers
	userHelper := helper.NewUserHelper()

	// Initialize mappers
	userMapper := mapper.NewUserMapper(db)

	// Initialize services
	userService := service.NewUserService(userMapper, userHelper)

	// Initialize filters
	authFilter := filter.NewAuthFilter(userService)

	// Initialize controllers
	userController := handler.NewUserController(userService)

	// Initialize Gin router
	r := gin.Default()

	// Apply aspects as middleware
	r.Use(func(c *gin.Context) {
		loggingAspect.Handle(c)
	})

	r.Use(func(c *gin.Context) {
		preprocessingAspect.Handle(c)
	})

	// Define routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		v1.GET("/verify-email", userController.VerifyEmail)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(func(c *gin.Context) {
			authFilter.Handle(c)
		})
		{
			protected.GET("/profile", userController.GetProfile)
			protected.PUT("/profile", userController.UpdateProfile)
			protected.PUT("/password", userController.ChangePassword)
			protected.POST("/request-verification", userController.RequestEmailVerification)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB initializes the database connection
func initDB() (*sql.DB, error) {
	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "oidc")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}