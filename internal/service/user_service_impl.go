// Package service implements the service layer interfaces for the OIDC application.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/mapper"
	model "github.com/Full-finger/OIDC/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// userService implements UserService interface
type userService struct {
	userMapper mapper.UserMapper
	userHelper helper.UserHelper
	version    string
}

// NewUserService creates a new UserService instance
func NewUserService(
	userMapper mapper.UserMapper,
	userHelper helper.UserHelper,
) UserService {
	return &userService{
		userMapper: userMapper,
		userHelper: userHelper,
		version:    "1.0.0",
	}
}

// GetDomainHelper returns the domain helper
func (us *userService) GetDomainHelper() interface{} {
	return us.userHelper
}

// GetDomainMapper returns the domain mapper
func (us *userService) GetDomainMapper() interface{} {
	return us.userMapper
}

// ConvertToEntity converts DTO to entity
func (us *userService) ConvertToEntity(dto interface{}) interface{} {
	// In a real implementation, you would convert DTO to entity
	return dto
}

// ConvertToDTO converts entity to DTO
func (us *userService) ConvertToDTO(entity interface{}) interface{} {
	// In a real implementation, you would convert entity to DTO
	return entity
}

// Register registers a new user
func (us *userService) Register(ctx context.Context, username, email, password string) error {
	// Check if user already exists
	existingUser, err := us.userMapper.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := us.userMapper.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// RegisterWithVerification registers a new user with email verification
func (us *userService) RegisterWithVerification(ctx context.Context, username, email, password string) (*model.SafeUser, error) {
	// Check if user already exists
	existingUser, err := us.userMapper.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := us.userMapper.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return safe user
	safeUser := user.SafeUser()
	return &safeUser, nil
}

// Login logs in a user
func (us *userService) Login(ctx context.Context, email, password string) (*model.User, error) {
	// Get user by email
	user, err := us.userMapper.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}

// GetProfile gets user profile
func (us *userService) GetProfile(ctx context.Context, userID int64) (*model.SafeUser, error) {
	// Get user by ID
	user, err := us.userMapper.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Return safe user
	safeUser := user.SafeUser()
	return &safeUser, nil
}

// UpdateProfile updates user profile
func (us *userService) UpdateProfile(ctx context.Context, userID int64, nickname, avatarURL, bio *string) error {
	// Get user by ID
	user, err := us.userMapper.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Update user fields
	if nickname != nil {
		user.Nickname = nickname
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}
	if bio != nil {
		user.Bio = bio
	}

	// Save updated user
	if err := us.userMapper.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ChangePassword changes user password
func (us *userService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// Get user by ID
	user, err := us.userMapper.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := us.userMapper.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// RequestEmailVerification requests email verification
func (us *userService) RequestEmailVerification(ctx context.Context, userID int64) error {
	// Get user by ID
	user, err := us.userMapper.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Check request frequency limit
	if !us.userHelper.CanRequestEmailVerification(user.Email) {
		return fmt.Errorf("please wait before requesting another verification email")
	}

	// Record email verification request
	us.userHelper.RecordEmailVerificationRequest(user.Email)

	// In a real implementation, you would generate a verification token and send the verification email
	// token, err := utils.GenerateVerificationToken(user.ID)
	// if err != nil {
	//     return fmt.Errorf("failed to generate verification token: %w", err)
	// }
	// sendVerificationEmail(user.Email, token)

	return nil
}

// VerifyEmail verifies email
func (us *userService) VerifyEmail(ctx context.Context, token string) error {
	// In a real implementation, you would verify the email with the token
	return nil
}

// GenerateJWT generates JWT token
func (us *userService) GenerateJWT(user *model.User) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token with a secret key (in a real app, use a secure secret)
	secretKey := "your-secret-key" // This should be loaded from environment variables
	return token.SignedString([]byte(secretKey))
}

// ValidateJWT validates JWT token
func (us *userService) ValidateJWT(tokenString string) (*model.User, error) {
	// Parse token
	secretKey := "your-secret-key" // This should be loaded from environment variables
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	// Create a minimal user object with the information from the token
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("failed to extract user ID from token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to extract email from token")
	}

	user := &model.User{
		ID:    int64(userIDFloat),
		Email: email,
	}

	return user, nil
}

// GetVersion returns the version of the service
func (us *userService) GetVersion() string {
	return us.version
}

// GetUserByID gets a user by ID
func (us *userService) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return us.userMapper.GetUserByID(ctx, userID)
}

// CanRequestVerificationEmail checks if a verification email can be requested
func (us *userService) CanRequestVerificationEmail(email string) bool {
	return us.userHelper.CanRequestEmailVerification(email)
}

// UpdateLastEmailRequestTime updates the last email request time
func (us *userService) UpdateLastEmailRequestTime(email string) {
	us.userHelper.RecordEmailVerificationRequest(email)
}