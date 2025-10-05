package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/util"
	"github.com/gin-gonic/gin"
)

// VerificationHandler 处理邮箱验证相关请求
type VerificationHandler struct {
	userService    service.UserService
	emailService   *util.EmailService
	frontendURL    string
	verificationURL string
}

// NewVerificationHandler 创建验证处理器实例
func NewVerificationHandler(userService service.UserService) *VerificationHandler {
	// 从环境变量获取配置
	smtpHost := util.GetEnv("SMTP_HOST", "smtp.example.com")
	smtpPort := util.GetEnv("SMTP_PORT", "587")
	smtpUser := util.GetEnv("SMTP_USER", "user@example.com")
	smtpPassword := util.GetEnv("SMTP_PASSWORD", "password")
	
	frontendURL := util.GetEnv("FRONTEND_URL", "http://localhost:3000")
	verificationURL := util.GetEnv("VERIFICATION_URL", "http://localhost:8080/api/v1/verify")

	emailService := util.NewEmailService(smtpHost, smtpPort, smtpUser, smtpPassword)

	return &VerificationHandler{
		userService:     userService,
		emailService:    emailService,
		frontendURL:     frontendURL,
		verificationURL: verificationURL,
	}
}

// SendVerificationEmailRequest 发送验证邮件请求结构
type SendVerificationEmailRequest struct {
	Email string `json:"email"`
}

// SendVerificationEmail 发送验证邮件
func (h *VerificationHandler) SendVerificationEmail(c *gin.Context) {
	var req SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 检查频率限制
	if !h.userService.CanRequestVerificationEmail(req.Email) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后再试"})
		return
	}

	// TODO: 查找用户并生成验证令牌
	// 这里需要根据实际业务逻辑实现

	// 更新最后请求时间
	h.userService.UpdateLastEmailRequestTime(req.Email)

	// 生成验证链接
	token, err := generateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	verificationURL := fmt.Sprintf("%s?token=%s", h.verificationURL, token)

	// 发送验证邮件
	err = h.emailService.SendVerificationEmail(req.Email, verificationURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证邮件已发送，请检查您的邮箱"})
}

// VerifyEmail 验证邮箱
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?error=missing_token", h.frontendURL))
		return
	}

	// TODO: 验证令牌并激活用户
	// 这里需要根据实际业务逻辑实现

	// 重定向到前端页面
	c.Redirect(http.StatusFound, fmt.Sprintf("%s?verified=true", h.frontendURL))
}

// generateVerificationToken 生成验证令牌
func generateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}