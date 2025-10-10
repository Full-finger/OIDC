// internal/handler/user_handler.go

package handler

import (
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	// Define request structure
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	safeUser, err := h.userService.RegisterWithVerification(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, safeUser)
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	// Define request structure
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	user, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.userService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.SafeUser(),
	})
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (in a real implementation, this would come from authentication)
	userID := int64(1) // Placeholder

	// Call service
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, user.SafeUser())
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (in a real implementation, this would come from authentication)
	userID := int64(1) // Placeholder

	// Define request structure
	var req struct {
		Nickname  *string `json:"nickname"`
		AvatarURL *string `json:"avatar_url"`
		Bio       *string `json:"bio"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	if err := h.userService.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL, req.Bio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
