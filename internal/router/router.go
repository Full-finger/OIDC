package router

import (
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/util"
	"github.com/Full-finger/OIDC/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化依赖
	userRepo := repository.NewUserRepository()
	userHelper := helper.NewUserHelper()
	tokenRepo := repository.NewVerificationTokenRepository()
	emailQueue := util.NewSimpleEmailQueue()
	
	userService := service.NewUserService(userRepo, userHelper, tokenRepo, emailQueue)
	userHandler := handler.NewUserHandler(userService)
	verificationHandler := handler.NewVerificationHandler(userService)

	// 初始化OAuth依赖
	oauthService := service.NewOAuthService()
	oauthHandler := handler.NewOAuthHandler(oauthService)

	// 初始化限流中间件
	rateLimiter := middleware.NewRateLimiter()
	// 设置为每5分钟最多5次请求
	rateLimiter.SetLimit(5, 5*60)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		v1.POST("/register", rateLimiter.LimitByIP(), userHandler.Register)
		v1.POST("/resend-verification", rateLimiter.LimitByUser(), userHandler.ResendVerificationEmail)
		v1.POST("/login", userHandler.Login)
		// 邮箱验证路由
		v1.GET("/verify", verificationHandler.VerifyEmail)
	}

	// OAuth 2.0 路由
	oauth := r.Group("/oauth")
	{
		// 授权端点
		oauth.GET("/authorize", oauthHandler.AuthorizeHandler)
		// TODO: 添加其他OAuth端点，如token、userinfo等
	}

	return r
}