package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Full-finger/OIDC/internal/model"
)

// VerificationTokenRepository 验证令牌仓库接口
type VerificationTokenRepository interface {
	Create(ctx context.Context, token *model.VerificationToken) error
	FindByToken(ctx context.Context, token string) (*model.VerificationToken, error)
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) error
}

type verificationTokenRepository struct {
	db *sql.DB
}

// NewVerificationTokenRepository 创建验证令牌仓库实例
func NewVerificationTokenRepository(db *sql.DB) VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

// Create 创建验证令牌
func (r *verificationTokenRepository) Create(ctx context.Context, token *model.VerificationToken) error {
	const query = `
		INSERT INTO verification_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&token.ID)
}

// FindByToken 根据令牌查找验证令牌记录
func (r *verificationTokenRepository) FindByToken(ctx context.Context, token string) (*model.VerificationToken, error) {
	const query = `
		SELECT id, user_id, token, expires_at, created_at
		FROM verification_tokens
		WHERE token = $1 AND expires_at > $2
		LIMIT 1`

	verificationToken := &model.VerificationToken{}
	err := r.db.QueryRowContext(ctx, query, token, time.Now()).Scan(
		&verificationToken.ID,
		&verificationToken.UserID,
		&verificationToken.Token,
		&verificationToken.ExpiresAt,
		&verificationToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return verificationToken, nil
}

// Delete 删除验证令牌
func (r *verificationTokenRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM verification_tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByUserID 根据用户ID删除验证令牌
func (r *verificationTokenRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	const query = `DELETE FROM verification_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpired 删除过期的验证令牌
func (r *verificationTokenRepository) DeleteExpired(ctx context.Context) error {
	const query = `DELETE FROM verification_tokens WHERE expires_at <= $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}