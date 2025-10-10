// Package handler implements the controller layer interfaces for the OIDC application.
package handler

import (
	"net/http"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// userController implements UserController interface
type userController struct {
	userService service.UserService
	version     string
}

// NewUserController creates a new UserController instance
func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
		version:     "1.0.0",
	}
}

// Register handles user registration
func (uc *userController) Register(c *gin.Context) {
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
	if err := uc.userService.Register(c.Request.Context(), req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
func (uc *userController) Login(c *gin.Context) {
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
	user, err := uc.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := uc.userService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user.SafeUser(),
	})
}

// GetProfile handles getting user profile
func (uc *userController) GetProfile(c *gin.Context) {
	// Get user ID from context (in a real implementation, this would come from authentication)
	userID := int64(1) // Placeholder

	// Call service
	user, err := uc.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles updating user profile
func (uc *userController) UpdateProfile(c *gin.Context) {
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
	if err := uc.userService.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL, req.Bio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ChangePassword handles changing user password
func (uc *userController) ChangePassword(c *gin.Context) {
	// Get user ID from context (in a real implementation, this would come from authentication)
	userID := int64(1) // Placeholder

	// Define request structure
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	if err := uc.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// RequestEmailVerification handles requesting email verification
func (uc *userController) RequestEmailVerification(c *gin.Context) {
	// Get user ID from context (in a real implementation, this would come from authentication)
	userID := int64(1) // Placeholder

	// Call service
	if err := uc.userService.RequestEmailVerification(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// VerifyEmail handles verifying email
func (uc *userController) VerifyEmail(c *gin.Context) {
	// Get token from query parameters
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	// Call service
	if err := uc.userService.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// HealthCheck checks the health of the controller
func (uc *userController) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the controller
func (uc *userController) GetVersion() string {
	return uc.version
}