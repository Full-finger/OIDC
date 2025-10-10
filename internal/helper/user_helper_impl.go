// Package helper implements the helper interfaces for the OIDC application.
package helper

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	model "github.com/Full-finger/OIDC/config"
)

// userHelper implements UserHelper interface
type userHelper struct {
	// Store last email request times (in real application, use Redis or other external storage)
	lastEmailRequest map[string]time.Time
	emailRequestMutex sync.Mutex
	version string
}

// NewUserHelper creates a new UserHelper instance
func NewUserHelper() UserHelper {
	return &userHelper{
		lastEmailRequest: make(map[string]time.Time),
		version: "1.0.0",
	}
}

// Validate validates the entity
func (uh *userHelper) Validate(entity interface{}) error {
	if user, ok := entity.(*model.User); ok {
		return uh.ValidateUser(user)
	}
	return fmt.Errorf("invalid entity type")
}

// Format formats the entity
func (uh *userHelper) Format(entity interface{}) interface{} {
	if user, ok := entity.(*model.User); ok {
		return uh.FormatUser(user)
	}
	return entity
}

// ValidateUser validates user entity
func (uh *userHelper) ValidateUser(user *model.User) error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// FormatUser formats user entity
func (uh *userHelper) FormatUser(user *model.User) *model.User {
	// Trim whitespace from username and email
	if user.Username != "" {
		// In a real implementation, you might want to do more formatting
	}
	
	if user.Email != "" {
		// In a real implementation, you might want to do more formatting
	}
	
	return user
}

// CanRequestEmailVerification checks if user can request email verification
func (uh *userHelper) CanRequestEmailVerification(email string) bool {
	uh.emailRequestMutex.Lock()
	defer uh.emailRequestMutex.Unlock()

	lastRequest, exists := uh.lastEmailRequest[email]
	if !exists {
		return true
	}

	// Check if 1 minute has passed since last request
	return time.Since(lastRequest) >= time.Minute
}

// RecordEmailVerificationRequest records email verification request
func (uh *userHelper) RecordEmailVerificationRequest(email string) {
	uh.emailRequestMutex.Lock()
	defer uh.emailRequestMutex.Unlock()

	uh.lastEmailRequest[email] = time.Now()
}

// HealthCheck checks the health of the helper
func (uh *userHelper) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the helper
func (uh *userHelper) GetVersion() string {
	return uh.version
}