// Package filter provides filtering functionality for the OIDC application.
package filter

import (
	"net/http"
	"strings"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthFilter implements authentication and authorization filtering
type AuthFilter struct {
	userService service.UserService
	enabled     bool
}

// NewAuthFilter creates a new AuthFilter instance
func NewAuthFilter(userService service.UserService) *AuthFilter {
	return &AuthFilter{
		userService: userService,
		enabled:     true,
	}
}

// Handle handles authentication and authorization
func (af *AuthFilter) Handle(c *gin.Context) {
	if !af.enabled {
		c.Next()
		return
	}

	// Get Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		c.Abort()
		return
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token
	user, err := af.userService.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	// Store user in context for later use
	c.Set("user", user)

	// Continue with the request
	c.Next()
}

// Enable enables the filter
func (af *AuthFilter) Enable() {
	af.enabled = true
}

// Disable disables the filter
func (af *AuthFilter) Disable() {
	af.enabled = false
}