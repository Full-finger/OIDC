// Package handler defines the controller layer interfaces for the OIDC application.
package handler

import (
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// UserController defines the user controller interface
type UserController interface {
	IBaseController

	// Register handles user registration
	Register(c *gin.Context)

	// Login handles user login
	Login(c *gin.Context)

	// GetProfile handles getting user profile
	GetProfile(c *gin.Context)

	// UpdateProfile handles updating user profile
	UpdateProfile(c *gin.Context)

	// ChangePassword handles changing user password
	ChangePassword(c *gin.Context)

	// RequestEmailVerification handles requesting email verification
	RequestEmailVerification(c *gin.Context)

	// VerifyEmail handles verifying email
	VerifyEmail(c *gin.Context)
}

// GetUserControllerServiceHelper provides a helper function to get user service
func GetUserControllerServiceHelper(userService service.UserService) service.IBaseService {
	return userService
}