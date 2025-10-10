// Package mapper implements the mapper interfaces for the OIDC application.
package mapper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	model "github.com/Full-finger/OIDC/config"
	_ "github.com/lib/pq"
)

// userMapper implements UserMapper interface
type userMapper struct {
	db      *sql.DB
	version string
}

// NewUserMapper creates a new UserMapper instance
func NewUserMapper(db *sql.DB) UserMapper {
	return &userMapper{
		db:      db,
		version: "1.0.0",
	}
}

// CreateUser creates a new user
func (um *userMapper) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := um.db.QueryRowContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Nickname,
		user.AvatarURL,
		user.Bio,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID gets a user by ID
func (um *userMapper) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &model.User{}
	err := um.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Nickname,
		&user.AvatarURL,
		&user.Bio,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetUserByUsername gets a user by username
func (um *userMapper) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at
		FROM users
		WHERE username = $1`

	user := &model.User{}
	err := um.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Nickname,
		&user.AvatarURL,
		&user.Bio,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (um *userMapper) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &model.User{}
	err := um.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Nickname,
		&user.AvatarURL,
		&user.Bio,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user
func (um *userMapper) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, email = $3, nickname = $4, avatar_url = $5, bio = $6, email_verified = $7, updated_at = $8
		WHERE id = $9`

	user.UpdatedAt = time.Now()

	result, err := um.db.ExecContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Nickname,
		user.AvatarURL,
		user.Bio,
		user.EmailVerified,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user
func (um *userMapper) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := um.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ListUsers lists users with pagination
func (um *userMapper) ListUsers(ctx context.Context, offset, limit int) ([]*model.User, error) {
	query := `
		SELECT id, username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := um.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.Email,
			&user.Nickname,
			&user.AvatarURL,
			&user.Bio,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

// CountUsers counts the total number of users
func (um *userMapper) CountUsers(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := um.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// HealthCheck checks the health of the mapper
func (um *userMapper) HealthCheck() error {
	// Implement health check logic
	return nil
}

// GetVersion returns the version of the mapper
func (um *userMapper) GetVersion() string {
	return um.version
}