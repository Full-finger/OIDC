package mapper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// UserMapper 用户映射器接口
type UserMapper interface {
	BaseMapper

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*model.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(email string) (*model.User, error)

	// UpdateActivationStatus 更新用户激活状态
	UpdateActivationStatus(id uint, isActive bool) error
	
	// GetByID 根据ID获取用户
	GetByID(id uint) (*model.User, error)
}