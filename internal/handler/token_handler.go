package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	oauthService service.OAuthService
}

func NewTokenHandler(oauthService service.OAuthService) *TokenHandler {
	return &TokenHandler{oauthService: oauthService}
}

func (h *TokenHandler) TokenHandler(c *gin.Context) {
	log.Printf("=== Token Request Debug Info ===")
	log.Printf("Remote IP: %s", c.ClientIP())
	log.Printf("Request URL: %s", c.Request.URL.Path)
	log.Printf("Request Method: %s", c.Request.Method)
	
	// 记录所有请求头
	log.Printf("=== Request Headers ===")
	for name, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("Header: %s = %s", name, value)
		}
	}
	
	// 记录所有表单参数
	log.Printf("=== Parsing Form ===")
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Unable to parse form data"})
		return
	}
	
	log.Printf("=== Form Values ===")
	for key, values := range c.Request.Form {
		for _, value := range values {
			log.Printf("Form: %s = %s", key, value)
		}
	}
	
	// 获取客户端认证信息
	clientID, clientSecret, ok := c.Request.BasicAuth()
	log.Printf("Basic Auth - ClientID: '%s', HasSecret: %t, OK: %t", clientID, clientSecret != "", ok)
	
	if !ok {
		log.Printf("Basic auth not provided or invalid, checking form values")
		// 检查是否通过表单传递客户端凭据
		formClientID := c.Request.Form.Get("client_id")
		formClientSecret := c.Request.Form.Get("client_secret")
		log.Printf("Form Auth - ClientID: '%s', HasSecret: %t", formClientID, formClientSecret != "")
		
		if formClientID == "" && formClientSecret == "" {
			log.Printf("No client credentials provided in token request (both basic auth and form are empty)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "No client credentials provided"})
			return
		}
		
		// 使用表单中的凭据
		clientID = formClientID
		clientSecret = formClientSecret
	} else {
		log.Printf("Using Basic Auth credentials")
	}

	log.Printf("Client authentication attempt - Client ID: '%s'", clientID)

	// 验证客户端
	log.Printf("Looking up client in database...")
	client, err := h.oauthService.FindClientByClientID(c.Request.Context(), clientID)
	if err != nil {
		log.Printf("Error retrieving client '%s' from database: %v", clientID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "Client not found in database"})
		return
	}
	
	log.Printf("Found client in database - ID: %d, Name: %s", client.ID, client.Name)

	// 检查客户端密钥
	log.Printf("Verifying client secret...")
	hashedSecret := hashSecret(clientSecret)
	log.Printf("Provided secret hash: %s", hashedSecret)
	log.Printf("Stored secret hash: %s", client.ClientSecretHash)
	log.Printf("Secret match: %t", client.ClientSecretHash == hashedSecret)
	
	if client.ClientSecretHash != hashedSecret {
		log.Printf("Client secret mismatch for client: '%s'", clientID)
		log.Printf("Provided secret (first 3 chars): '%s...'", clientSecret[:min(3, len(clientSecret))])
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "Invalid client credentials"})
		return
	}
	
	log.Printf("Client secret verified successfully")

	log.Printf("Client authentication successful. Client ID: %s", clientID)

	grantType := c.Request.Form.Get("grant_type")
	log.Printf("Grant type requested: %s", grantType)

	switch grantType {
	case "authorization_code":
		h.handleAuthorizationCodeGrant(c, client)
	case "refresh_token":
		h.handleRefreshTokenGrant(c, client)
	default:
		log.Printf("Unsupported grant type: %s", grantType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type", "error_description": "Unsupported grant type"})
	}
}

func (h *TokenHandler) handleRefreshTokenGrant(c *gin.Context, client *model.Client) {
	refreshToken := c.Request.Form.Get("refresh_token")
	
	log.Printf("Refresh token grant flow. Refresh Token provided: %t", refreshToken != "")

	if refreshToken == "" {
		log.Printf("Missing refresh token in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Missing refresh token"})
		return
	}

	// 使用刷新令牌获取新的访问令牌和刷新令牌
	log.Printf("Refreshing access token with refresh token...")
	result, err := h.oauthService.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Printf("Error refreshing access token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": err.Error()})
		return
	}

	log.Printf("Access token refreshed successfully")
	log.Printf("New access token length: %d", len(result.AccessToken))
	log.Printf("New refresh token length: %d", len(result.RefreshToken))
	if result.IDToken != "" {
		log.Printf("New ID token generated successfully")
	}

	// 返回令牌响应
	response := map[string]interface{}{
		"access_token":  result.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}
	
	// 只有当刷新令牌存在时才添加到响应中
	if result.RefreshToken != "" {
		response["refresh_token"] = result.RefreshToken
	}
	
	// 只有当ID Token存在时才添加到响应中
	if result.IDToken != "" {
		response["id_token"] = result.IDToken
	}

	c.JSON(http.StatusOK, response)
}

func (h *TokenHandler) handleAuthorizationCodeGrant(c *gin.Context, client *model.Client) {
	code := c.Request.Form.Get("code")
	redirectURI := c.Request.Form.Get("redirect_uri")

	log.Printf("Authorization code grant flow. Code: %s, Redirect URI: %s", code, redirectURI)

	if code == "" {
		log.Printf("Missing authorization code in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Missing authorization code"})
		return
	}

	// 验证重定向URI
	if redirectURI == "" {
		log.Printf("Missing redirect_uri in token request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Missing redirect_uri"})
		return
	}

	// 验证授权码并兑换访问令牌和刷新令牌
	log.Printf("Exchanging authorization code for access token and refresh token...")
	result, err := h.oauthService.ExchangeAuthorizationCode(c.Request.Context(), code, client.ClientID, redirectURI)
	if err != nil {
		log.Printf("Error exchanging authorization code %s: %v", code, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": err.Error()})
		return
	}

	log.Printf("Access token generated successfully: %s", result.AccessToken)
	log.Printf("Refresh token generated successfully: %s", result.RefreshToken)
	if result.IDToken != "" {
		log.Printf("ID token generated successfully: %s", result.IDToken)
	}

	// 返回令牌响应
	response := map[string]interface{}{
		"access_token":  result.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}
	
	// 只有当刷新令牌存在时才添加到响应中
	if result.RefreshToken != "" {
		response["refresh_token"] = result.RefreshToken
	}
	
	// 只有当ID Token存在时才添加到响应中
	if result.IDToken != "" {
		response["id_token"] = result.IDToken
	}

	c.JSON(http.StatusOK, response)
}

// hashSecret 对密钥进行哈希处理（与service中保持一致）
func hashSecret(secret string) string {
	// 注意：这应该与service/oauth_service.go中的hashSecret保持一致
	// 在实际项目中，应该将这个函数提取到一个公共包中
	hash := sha256.Sum256([]byte(secret))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
