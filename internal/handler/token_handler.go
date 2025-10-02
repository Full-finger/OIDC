// internal/handler/token_handler.go

package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// TokenHandler 定义令牌相关处理函数
type TokenHandler struct {
	oauthService service.OAuthService
}

// TokenRequest 定义令牌请求结构
type TokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	Code         string `json:"code" form:"code"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
}

// TokenResponse 定义令牌响应结构
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewTokenHandler 创建一个新的TokenHandler实例
func NewTokenHandler(oauthService service.OAuthService) *TokenHandler {
	return &TokenHandler{
		oauthService: oauthService,
	}
}

// TokenHandler 处理POST /oauth/token请求
func (h *TokenHandler) TokenHandler(c *gin.Context) {
	var req TokenRequest
	
	// 绑定请求数据，支持JSON和表单格式
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}
	
	// 获取客户端认证信息
	clientID, clientSecret, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		// 使用Basic Auth
		req.ClientID = clientID
		req.ClientSecret = clientSecret
	}
	
	// 验证必需参数
	if req.GrantType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "grant_type is required"})
		return
	}
	
	// 验证grant_type是否为authorization_code
	if req.GrantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type", "error_description": "grant_type must be authorization_code"})
		return
	}
	
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "code is required"})
		return
	}
	
	if req.RedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "redirect_uri is required"})
		return
	}
	
	if req.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "client_id is required"})
		return
	}
	
	// 验证客户端身份
	client, err := h.oauthService.FindClientByClientID(c.Request.Context(), req.ClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "client not found"})
		return
	}
	
	// 如果使用POST方式传递客户端凭证，则验证客户端密钥
	if !hasAuth && req.ClientSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "client_secret is required"})
		return
	}
	
	if !hasAuth && req.ClientSecret != "" {
		// 验证客户端密钥（简化实现，实际应使用更安全的比较方法）
		if client.ClientSecretHash != hashSecret(req.ClientSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "invalid client credentials"})
			return
		}
	}
	
	// 如果使用Basic Auth，则验证客户端密钥
	if hasAuth {
		// 验证客户端密钥（简化实现，实际应使用更安全的比较方法）
		if client.ClientSecretHash != hashSecret(req.ClientSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "invalid client credentials"})
			return
		}
	}
	
	// 兑换授权码获取访问令牌
	accessToken, err := h.oauthService.ExchangeAuthorizationCode(
		c.Request.Context(),
		req.Code,
		req.ClientID,
		req.RedirectURI,
	)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": err.Error()})
		return
	}
	
	// 返回成功响应
	response := TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1小时，单位秒
	}
	
	c.JSON(http.StatusOK, response)
}

// hashSecret 对密钥进行哈希处理（与oauth_service.go中的一致）
func hashSecret(secret string) string {
	// 简化实现，实际项目中应该使用bcrypt等安全的哈希算法
	hash := sha256.Sum256([]byte(secret))
	return base64.StdEncoding.EncodeToString(hash[:])
}