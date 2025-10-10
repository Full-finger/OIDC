// Package mapper defines the mapper interfaces for the OIDC application.
package mapper

import (
	"context"

	model "github.com/Full-finger/OIDC/config"
)

// UserMapper defines the user mapper interface
type UserMapper interface {
	IBaseMapper

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *model.User) error

	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, id int64) (*model.User, error)

	// GetUserByUsername gets a user by username
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserByEmail gets a user by email
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *model.User) error

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id int64) error

	// ListUsers lists users with pagination
	ListUsers(ctx context.Context, offset, limit int) ([]*model.User, error)

	// CountUsers counts the total number of users
	CountUsers(ctx context.Context) (int64, error)
}