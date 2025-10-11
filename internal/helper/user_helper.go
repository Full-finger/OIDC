package helper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// UserHelper 用户助手接口
type UserHelper interface {
	// ValidateUser 验证用户数据
	ValidateUser(user *model.User) error

	// HashPassword 对密码进行哈希处理
	HashPassword(password string) (string, error)

	// CheckPassword 验证密码
	CheckPassword(hashedPassword, password string) bool

	// GenerateAvatarURL 生成头像URL
	GenerateAvatarURL(userID uint) string
}