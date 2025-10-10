// Package helper defines the helper interfaces for the OIDC application.
package helper

import (
	model "github.com/Full-finger/OIDC/config"
)

// UserHelper defines the user helper interface
type UserHelper interface {
	IBaseHelper

	// ValidateUser validates user entity
	ValidateUser(user *model.User) error

	// FormatUser formats user entity
	FormatUser(user *model.User) *model.User

	// CanRequestEmailVerification checks if user can request email verification
	CanRequestEmailVerification(email string) bool

	// RecordEmailVerificationRequest records email verification request
	RecordEmailVerificationRequest(email string)
}