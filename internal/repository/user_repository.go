package repository

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(user *model.User) error

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*model.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(email string) (*model.User, error)

	// GetByID 根据ID获取用户
	GetByID(id uint) (*model.User, error)

	// Update 更新用户信息
	Update(user *model.User) error

	// UpdateActivationStatus 更新用户激活状态
	UpdateActivationStatus(id uint, isActive bool) error
}