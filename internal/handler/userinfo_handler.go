// Package handler implements HTTP handlers for the OIDC application.
package handler

import (
	"log"
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// UserInfoHandler handles /userinfo endpoint requests
type UserInfoHandler struct {
	userService  service.UserService
	oauthService service.OAuthService
}

// NewUserInfoHandler creates a new UserInfoHandler instance
func NewUserInfoHandler(userService service.UserService, oauthService service.OAuthService) *UserInfoHandler {
	return &UserInfoHandler{
		userService:  userService,
		oauthService: oauthService,
	}
}

// GetUserInfo handles GET /userinfo requests
// This function must be protected by JWT middleware
func (h *UserInfoHandler) GetUserInfo(c *gin.Context) {
	log.Printf("UserInfo Handler: Processing userinfo request")
	
	// 从上下文中获取用户ID（由JWT中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("UserInfo Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 从上下文中获取scope（由JWT中间件注入）
	scopes, exists := c.Get("scope")
	if !exists {
		scopes = []string{} // 默认空scope
	}
	
	log.Printf("UserInfo Handler: Getting user info for user ID: %d with scopes: %v", userID, scopes)
	
	// 获取用户信息
	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		log.Printf("UserInfo Handler: Failed to get user by ID %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	
	log.Printf("UserInfo Handler: Successfully retrieved user info for user ID: %d", userID)
	
	// 构建响应
	response := map[string]interface{}{
		"sub": userID,
	}
	
	// 根据scope添加额外信息
	scopeList := scopes.([]string)
	for _, scope := range scopeList {
		switch scope {
		case "profile":
			response["name"] = user.Username
			if user.Nickname != nil {
				response["nickname"] = user.Nickname
			}
			if user.AvatarURL != nil {
				response["picture"] = user.AvatarURL
			}
		case "email":
			response["email"] = user.Email
			response["email_verified"] = user.EmailVerified
		}
	}
	
	log.Printf("UserInfo Handler: Returning user info response: %v", response)
	c.JSON(http.StatusOK, response)
}