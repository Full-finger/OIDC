package repository

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// VerificationTokenRepository 验证令牌仓储接口
type VerificationTokenRepository interface {
	// Create 创建验证令牌
	Create(token *model.VerificationToken) error

	// GetByToken 根据令牌获取验证令牌记录
	GetByToken(token string) (*model.VerificationToken, error)

	// GetByUserID 根据用户ID获取验证令牌记录
	GetByUserID(userID uint) (*model.VerificationToken, error)

	// Delete 删除验证令牌
	Delete(id uint) error

	// DeleteByToken 根据令牌删除验证令牌记录
	DeleteByToken(token string) error

	// DeleteExpired 删除过期的验证令牌
	DeleteExpired() error
}