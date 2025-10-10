// Package handler defines the controller layer interfaces for the OIDC application.
package handler

import (
	"github.com/gin-gonic/gin"
)

// OIDCController defines the OIDC controller interface
type OIDCController interface {
	IBaseController

	// Discovery handles OIDC discovery endpoint
	Discovery(c *gin.Context)

	// JWKS handles JWKS endpoint
	JWKS(c *gin.Context)

	// Authorize handles authorization endpoint
	Authorize(c *gin.Context)

	// Token handles token endpoint
	Token(c *gin.Context)

	// UserInfo handles userinfo endpoint
	UserInfo(c *gin.Context)
}