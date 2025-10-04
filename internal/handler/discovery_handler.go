// Package handler implements HTTP handlers for the OIDC application.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DiscoveryHandler 处理 OpenID Connect Discovery 端点请求
type DiscoveryHandler struct{}

// NewDiscoveryHandler 创建一个新的 DiscoveryHandler 实例
func NewDiscoveryHandler() *DiscoveryHandler {
	return &DiscoveryHandler{}
}

// GetDiscovery 处理 GET /.well-known/openid-configuration 请求
func (h *DiscoveryHandler) GetDiscovery(c *gin.Context) {
	// 构建OpenID Connect Discovery文档
	discovery := map[string]interface{}{
		"issuer":                                "http://localhost:8080",
		"authorization_endpoint":                "http://localhost:8080/oauth/authorize",
		"token_endpoint":                       "http://localhost:8080/oauth/token",
		"userinfo_endpoint":                    "http://localhost:8080/oauth/userinfo",
		"jwks_uri":                             "http://localhost:8080/.well-known/jwks.json",
		"scopes_supported":                     []string{"openid", "profile", "email"},
		"response_types_supported":             []string{"code"},
		"response_modes_supported":             []string{"query"},
		"grant_types_supported":                []string{"authorization_code", "refresh_token"},
		"subject_types_supported":              []string{"public"},
		"id_token_signing_alg_values_supported": []string{"HS256"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"claims_supported": []string{
			"sub",
			"iss",
			"aud",
			"exp",
			"iat",
			"auth_time",
			"nonce",
			"name",
			"nickname",
			"preferred_username",
			"picture",
			"email",
			"email_verified",
		},
		"code_challenge_methods_supported": []string{"S256"},
		"end_session_endpoint":            "http://localhost:8080/oauth/logout",
		"revocation_endpoint":             "http://localhost:8080/oauth/revoke",
		"introspection_endpoint":          "http://localhost:8080/oauth/introspect",
	}

	c.JSON(http.StatusOK, discovery)
}