// Package handler implements HTTP handlers for the OIDC application.
package handler

import (
	"log"
	"net/http"

	"github.com/Full-finger/OIDC/internal/middleware"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// UserInfoHandler 处理 /userinfo 端点请求
type UserInfoHandler struct {
	userService  *service.UserService
	oauthService service.OAuthService
}

// NewUserInfoHandler 创建一个新的 UserInfoHandler 实例
func NewUserInfoHandler(userService *service.UserService, oauthService service.OAuthService) *UserInfoHandler {
	return &UserInfoHandler{
		userService:  userService,
		oauthService: oauthService,
	}
}

// GetUserInfo 处理 GET /userinfo 请求
// 该函数必须被JWT中间件保护
func (h *UserInfoHandler) GetUserInfo(c *gin.Context) {
	log.Printf("UserInfo Handler: Processing userinfo request")
	
	// 从上下文中获取用户ID（由JWT中间件设置）
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("UserInfo Handler: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("UserInfo Handler: failed to convert userID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	
	// 从上下文中获取scopes（由JWT中间件设置）
	scopesInterface, exists := c.Get("scopes")
	var scopes []string
	if exists {
		if s, ok := scopesInterface.([]string); ok {
			scopes = s
		}
	}
	
	log.Printf("UserInfo Handler: Getting user info for user ID: %d with scopes: %v", userID, scopes)
	
	// 获取用户信息
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("UserInfo Handler: Failed to get user by ID %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	
	log.Printf("UserInfo Handler: Successfully retrieved user info for user ID: %d", userID)
	
	// 构建响应
	response := gin.H{
		"sub": user.ID,
	}
	
	// 根据scope决定返回哪些用户信息
	if middleware.ContainsProfileScope(scopes) {
		response["name"] = user.Username
		if user.Nickname != nil {
			response["nickname"] = user.Nickname
		}
		if user.AvatarURL != nil {
			response["picture"] = user.AvatarURL
		}
	}
	
	if middleware.ContainsEmailScope(scopes) {
		response["email"] = user.Email
		response["email_verified"] = user.EmailVerified
	}
	
	// 返回用户信息
	c.JSON(http.StatusOK, response)
}