package repository

import (
	"errors"
	"sync"
	
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/model"
	"gorm.io/gorm"
)

// userRepository 用户仓库实现
type userRepository struct {
	mapper mapper.UserMapper
	// 内存存储
	memoryStore map[string]*model.User
	mu          sync.RWMutex
}

// NewUserRepository 创建UserRepository实例
func NewUserRepository(mapper mapper.UserMapper) UserRepository {
	return &userRepository{
		mapper:      mapper,
		memoryStore: make(map[string]*model.User),
	}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	if r.mapper == nil {
		// 内存模式，存储用户信息
		r.mu.Lock()
		defer r.mu.Unlock()
		r.memoryStore[user.Username] = user
		r.memoryStore[user.Email] = user
		return nil
	}
	return r.mapper.Save(user)
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	if r.mapper == nil {
		// 内存模式，从内存中获取用户
		r.mu.RLock()
		defer r.mu.RUnlock()
		if user, exists := r.memoryStore[username]; exists {
			return user, nil
		}
		return nil, errors.New("用户不存在")
	}
	
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
	if r.mapper == nil {
		// 内存模式，从内存中获取用户
		r.mu.RLock()
		defer r.mu.RUnlock()
		if user, exists := r.memoryStore[email]; exists {
			return user, nil
		}
		return nil, errors.New("用户不存在")
	}
	
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
	if r.mapper == nil {
		// 内存模式，遍历查找用户
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, user := range r.memoryStore {
			if user.ID == id {
				return user, nil
			}
		}
		return nil, errors.New("用户不存在")
	}
	
	entity, err := r.mapper.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 类型断言
	user, ok := entity.(*model.User)
	if !ok {
		return nil, errors.New("类型转换失败")
	}
	
	return user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(user *model.User) error {
	if r.mapper == nil {
		// 内存模式，更新用户信息
		r.mu.Lock()
		defer r.mu.Unlock()
		r.memoryStore[user.Username] = user
		r.memoryStore[user.Email] = user
		return nil
	}
	return r.mapper.Save(user)
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	if r.mapper == nil {
		// 内存模式，查找并删除用户
		r.mu.Lock()
		defer r.mu.Unlock()
		for key, user := range r.memoryStore {
			if user.ID == id {
				delete(r.memoryStore, key)
				// 如果用户名和邮箱不同，也要删除邮箱键
				if key != user.Email {
					delete(r.memoryStore, user.Email)
				}
				break
			}
		}
		return nil
	}
	return r.mapper.DeleteByID(id)
}

// UpdateActivationStatus 更新用户激活状态
func (r *userRepository) UpdateActivationStatus(id uint, isActive bool) error {
	if r.mapper == nil {
		// 内存模式，更新用户激活状态
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, user := range r.memoryStore {
			if user.ID == id {
				user.IsActive = isActive
				// 同时更新存储中的键值
				r.memoryStore[user.Username] = user
				r.memoryStore[user.Email] = user
				break
			}
		}
		return nil
	}
	return r.mapper.UpdateActivationStatus(id, isActive)
}