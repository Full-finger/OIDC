package repository

import (
	"context"
	"database/sql"
	"github.com/Full-finger/OIDC/internal/model"
)

// BangumiRepository 定义Bangumi绑定相关的数据访问接口
type BangumiRepository interface {
	// CreateUserBangumiBinding 创建用户Bangumi绑定记录
	CreateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error
	
	// GetUserBangumiBindingByUserID 根据用户ID获取Bangumi绑定记录
	GetUserBangumiBindingByUserID(ctx context.Context, userID int64) (*model.UserBangumiBinding, error)
	
	// GetUserBangumiBindingByBangumiUserID 根据Bangumi用户ID获取Bangumi绑定记录
	GetUserBangumiBindingByBangumiUserID(ctx context.Context, bangumiUserID int64) (*model.UserBangumiBinding, error)
	
	// UpdateUserBangumiBinding 更新用户Bangumi绑定记录
	UpdateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error
	
	// DeleteUserBangumiBinding 删除用户Bangumi绑定记录
	DeleteUserBangumiBinding(ctx context.Context, userID int64) error
}

// bangumiRepository 实现BangumiRepository接口
type bangumiRepository struct {
	db *sql.DB
}

// NewBangumiRepository 创建BangumiRepository实例
func NewBangumiRepository(db *sql.DB) BangumiRepository {
	return &bangumiRepository{db: db}
}

// CreateUserBangumiBinding 创建用户Bangumi绑定记录
func (r *bangumiRepository) CreateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error {
	query := `
		INSERT INTO user_bangumi_bindings 
		(user_id, bangumi_user_id, access_token, refresh_token, token_expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		binding.UserID,
		binding.BangumiUserID,
		binding.AccessToken,
		binding.RefreshToken,
		binding.TokenExpiresAt).Scan(&binding.ID, &binding.CreatedAt, &binding.UpdatedAt)
	
	return err
}

// GetUserBangumiBindingByUserID 根据用户ID获取Bangumi绑定记录
func (r *bangumiRepository) GetUserBangumiBindingByUserID(ctx context.Context, userID int64) (*model.UserBangumiBinding, error) {
	query := `
		SELECT id, user_id, bangumi_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM user_bangumi_bindings
		WHERE user_id = $1`
	
	binding := &model.UserBangumiBinding{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&binding.ID,
		&binding.UserID,
		&binding.BangumiUserID,
		&binding.AccessToken,
		&binding.RefreshToken,
		&binding.TokenExpiresAt,
		&binding.CreatedAt,
		&binding.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return binding, nil
}

// GetUserBangumiBindingByBangumiUserID 根据Bangumi用户ID获取Bangumi绑定记录
func (r *bangumiRepository) GetUserBangumiBindingByBangumiUserID(ctx context.Context, bangumiUserID int64) (*model.UserBangumiBinding, error) {
	query := `
		SELECT id, user_id, bangumi_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM user_bangumi_bindings
		WHERE bangumi_user_id = $1`
	
	binding := &model.UserBangumiBinding{}
	err := r.db.QueryRowContext(ctx, query, bangumiUserID).Scan(
		&binding.ID,
		&binding.UserID,
		&binding.BangumiUserID,
		&binding.AccessToken,
		&binding.RefreshToken,
		&binding.TokenExpiresAt,
		&binding.CreatedAt,
		&binding.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return binding, nil
}

// UpdateUserBangumiBinding 更新用户Bangumi绑定记录
func (r *bangumiRepository) UpdateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error {
	query := `
		UPDATE user_bangumi_bindings
		SET bangumi_user_id = $1, access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = NOW()
		WHERE user_id = $5
		RETURNING updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		binding.BangumiUserID,
		binding.AccessToken,
		binding.RefreshToken,
		binding.TokenExpiresAt,
		binding.UserID).Scan(&binding.UpdatedAt)
	
	return err
}

// DeleteUserBangumiBinding 删除用户Bangumi绑定记录
func (r *bangumiRepository) DeleteUserBangumiBinding(ctx context.Context, userID int64) error {
	query := `DELETE FROM user_bangumi_bindings WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}