package mapper

import (
	"gorm.io/gorm"
	"github.com/Full-finger/OIDC/internal/model"
)

// userMapper 用户映射器实现
type userMapper struct {
	db *gorm.DB
}

// NewUserMapper 创建UserMapper实例
func NewUserMapper(db *gorm.DB) UserMapper {
	return &userMapper{db: db}
}

// Save 保存用户
func (m *userMapper) Save(entity interface{}) error {
	return m.db.Save(entity).Error
}

// DeleteByID 根据ID删除用户
func (m *userMapper) DeleteByID(id interface{}) error {
	return m.db.Delete(&model.User{}, id).Error
}

// GetByID 根据ID获取用户 (实现BaseMapper接口)
func (m *userMapper) GetByID(id interface{}) (interface{}, error) {
	var user model.User
	if err := m.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll 获取所有用户
func (m *userMapper) GetAll() ([]interface{}, error) {
	var users []model.User
	if err := m.db.Find(&users).Error; err != nil {
		return nil, err
	}
	
	interfaceUsers := make([]interface{}, len(users))
	for i, user := range users {
		interfaceUsers[i] = &user
	}
	
	return interfaceUsers, nil
}

// Update 更新用户
func (m *userMapper) Update(entity interface{}) error {
	return m.db.Save(entity).Error
}

// GetByUsername 根据用户名获取用户
func (m *userMapper) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := m.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (m *userMapper) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateActivationStatus 更新用户激活状态
func (m *userMapper) UpdateActivationStatus(id uint, isActive bool) error {
	return m.db.Model(&model.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}