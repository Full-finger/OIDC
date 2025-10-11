package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
)

// VerificationHandler 邮箱验证处理器
type VerificationHandler struct {
	userService service.UserService
}

// NewVerificationHandler 创建VerificationHandler实例
func NewVerificationHandler(userService service.UserService) *VerificationHandler {
	return &VerificationHandler{
		userService: userService,
	}
}

// VerifyEmail 处理邮箱验证请求
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	// 从查询参数中获取令牌
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少验证令牌"})
		return
	}

	// 调用服务层验证邮箱
	if err := h.userService.VerifyEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 邮箱验证成功
	c.JSON(http.StatusOK, gin.H{
		"message": "邮箱验证成功，您的账户已激活",
	})
}