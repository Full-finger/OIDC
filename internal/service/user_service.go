// internal/service/user_service.go

// Package service defines the service layer interfaces for the OIDC application.
package service

import (
	"context"

	model "github.com/Full-finger/OIDC/config"
)

// UserService defines the user service interface
type UserService interface {
	IBaseService
	ConvertInterface

	// Register registers a new user
	Register(ctx context.Context, username, email, password string) error

	// RegisterWithVerification registers a new user with email verification
	RegisterWithVerification(ctx context.Context, username, email, password string) (*model.SafeUser, error)

	// Login logs in a user
	Login(ctx context.Context, email, password string) (*model.User, error)

	// GetProfile gets user profile
	GetProfile(ctx context.Context, userID int64) (*model.SafeUser, error)

	// UpdateProfile updates user profile
	UpdateProfile(ctx context.Context, userID int64, nickname, avatarURL, bio *string) error

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error

	// RequestEmailVerification requests email verification
	RequestEmailVerification(ctx context.Context, userID int64) error

	// VerifyEmail verifies email
	VerifyEmail(ctx context.Context, token string) error

	// GenerateJWT generates JWT token
	GenerateJWT(user *model.User) (string, error)

	// ValidateJWT validates JWT token
	ValidateJWT(tokenString string) (*model.User, error)

	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)

	// CanRequestVerificationEmail checks if a verification email can be requested
	CanRequestVerificationEmail(email string) bool

	// UpdateLastEmailRequestTime updates the last email request time
	UpdateLastEmailRequestTime(email string)
}