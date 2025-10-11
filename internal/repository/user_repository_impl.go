package repository

import (
	"errors"

	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/model"
)

// userRepository 用户仓库实现
type userRepository struct {
	mapper mapper.UserMapper
}

// NewUserRepository 创建UserRepository实例
func NewUserRepository(mapper mapper.UserMapper) UserRepository {
	return &userRepository{mapper: mapper}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.mapper.Save(user)
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user, err := r.mapper.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user, err := r.mapper.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	user, err := r.mapper.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(user *model.User) error {
	return r.mapper.Save(user)
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.mapper.Delete(id)
}

// UpdateActivationStatus 更新用户激活状态
func (r *userRepository) UpdateActivationStatus(id uint, isActive bool) error {
	// TODO: 实现更新用户激活状态逻辑
	return nil
}