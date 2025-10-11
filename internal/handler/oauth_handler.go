package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
)

// OAuthHandler OAuth处理器
type OAuthHandler struct {
	oauthService service.OAuthService
}

// NewOAuthHandler 创建OAuthHandler实例
func NewOAuthHandler(oauthService service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// AuthorizeHandler 处理授权请求
func (h *OAuthHandler) AuthorizeHandler(c *gin.Context) {
	// 获取查询参数
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	scope := c.Query("scope")
	responseType := c.Query("response_type")
	state := c.Query("state")
	codeChallenge := c.Query("code_challenge")
	codeChallengeMethod := c.Query("code_challenge_method")

	// 验证必需参数
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing client_id"})
		return
	}

	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing redirect_uri"})
		return
	}

	if responseType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing response_type"})
		return
	}

	// 验证response_type是否为code
	if responseType != "code" {
		c.Redirect(http.StatusFound, redirectURI+"?error=unsupported_response_type&state="+state)
		return
	}

	// TODO: 验证用户是否已登录
	// 在实际实现中，应该检查用户是否已通过身份验证
	// 如果未登录，应该重定向到登录页面
	userID := "1" // 模拟用户ID

	// 解析scopes
	scopes := h.parseScopes(scope)

	// 处理授权请求
	var codeChallengePtr *string
	var codeChallengeMethodPtr *string
	
	if codeChallenge != "" {
		codeChallengePtr = &codeChallenge
	}
	
	if codeChallengeMethod != "" {
		codeChallengeMethodPtr = &codeChallengeMethod
	}

	// 调用服务层处理授权请求
	authCode, err := h.oauthService.HandleAuthorizationRequest(
		c.Request.Context(),
		clientID,
		userID,
		redirectURI,
		scopes,
		codeChallengePtr,
		codeChallengeMethodPtr,
	)
	
	if err != nil {
		c.Redirect(http.StatusFound, redirectURI+"?error=invalid_request&state="+state)
		return
	}

	// 重定向回客户端，携带授权码
	redirectURL := redirectURI + "?code=" + authCode.Code
	if state != "" {
		redirectURL += "&state=" + state
	}
	
	c.Redirect(http.StatusFound, redirectURL)
}

// parseScopes 解析scopes字符串
func (h *OAuthHandler) parseScopes(scope string) []string {
	if scope == "" {
		return []string{}
	}
	
	// 简单实现，实际应该更复杂
	scopes := []string{}
	// 这里应该解析空格分隔的scopes
	// 为简化示例，我们只返回一个scope
	scopes = append(scopes, scope)
	
	return scopes
}