package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/util"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建UserHandler实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRequest 用户注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required,max=50"`
}

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ResendVerificationEmailRequest 重新发送验证邮件请求结构体
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Register 用户注册接口
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层注册用户
	if err := h.userService.RegisterUser(req.Username, req.Password, req.Email, req.Nickname); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户注册成功，请检查邮箱以激活账户",
	})
}

// Login 用户登录接口
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层认证用户
	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
	})
}

// ResendVerificationEmail 重新发送验证邮件接口
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层重新发送验证邮件
	if err := h.userService.ResendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "验证邮件已重新发送，请检查您的邮箱",
	})
}

// GetProfile 获取用户资料接口
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: 实现获取用户资料逻辑
	c.JSON(http.StatusOK, gin.H{"message": "GetProfile endpoint"})
}

// UpdateProfile 更新用户资料接口
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// TODO: 实现更新用户资料逻辑
	c.JSON(http.StatusOK, gin.H{"message": "UpdateProfile endpoint"})
}