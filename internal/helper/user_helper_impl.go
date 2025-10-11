package helper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// userHelper 用户助手实现
type userHelper struct {
	// 可以添加依赖项
}

// NewUserHelper 创建UserHelper实例
func NewUserHelper() UserHelper {
	return &userHelper{}
}

// ValidateUser 验证用户数据
func (h *userHelper) ValidateUser(user *model.User) error {
	// TODO: 实现用户数据验证逻辑
	return nil
}

// HashPassword 对密码进行哈希处理
func (h *userHelper) HashPassword(password string) (string, error) {
	// TODO: 实现密码哈希逻辑
	return "", nil
}

// CheckPassword 验证密码
func (h *userHelper) CheckPassword(hashedPassword, password string) bool {
	// TODO: 实现密码验证逻辑
	return false
}

// GenerateAvatarURL 生成头像URL
func (h *userHelper) GenerateAvatarURL(userID uint) string {
	// TODO: 实现头像URL生成逻辑
	return ""
}