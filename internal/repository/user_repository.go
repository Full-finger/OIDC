// internal/repository/user_repository.go

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/Full-finger/OIDC/config"
)

// UserRepository 定义用户数据仓库接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
}

// userRepository 是 UserRepository 接口的 database/sql 实现
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建一个新的 userRepository 实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建新用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (
			username, password_hash, email, nickname, avatar_url, bio, email_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, created_at, updated_at`

	var nickname, avatarURL, bio *string
	if user.Nickname != nil {
		nickname = user.Nickname
	}
	if user.AvatarURL != nil {
		avatarURL = user.AvatarURL
	}
	if user.Bio != nil {
		bio = user.Bio
	}

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		nickname,
		avatarURL,
		bio,
		user.EmailVerified,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// FindByID 根据 ID 查找用户
func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findUser(ctx, "id = $1", id)
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.findUser(ctx, "username = $1", username)
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findUser(ctx, "email = $1", email)
}

// findUser 是内部辅助方法，避免重复代码
func (r *userRepository) findUser(ctx context.Context, whereClause string, args ...interface{}) (*model.User, error) {
	query := `
		SELECT 
			id, username, password_hash, email, nickname, avatar_url, bio, 
			email_verified, created_at, updated_at
		FROM users 
		WHERE ` + whereClause + ` LIMIT 1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
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
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

// Update 更新用户信息（不包括密码）
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users 
		SET 
			nickname = $2, 
			avatar_url = $3, 
			bio = $4, 
			email_verified = $5, 
			updated_at = $6
		WHERE id = $1`

	var nickname, avatarURL, bio *string
	if user.Nickname != nil {
		nickname = user.Nickname
	}
	if user.AvatarURL != nil {
		avatarURL = user.AvatarURL
	}
	if user.Bio != nil {
		bio = user.Bio
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		nickname,
		avatarURL,
		bio,
		user.EmailVerified,
		time.Now(),
	)

	return err
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	const query = `
		UPDATE users 
		SET password_hash = $1, updated_at = $2 
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}
