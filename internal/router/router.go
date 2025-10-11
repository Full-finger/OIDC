package router

import (
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/util"
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

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
	}

	return r
}