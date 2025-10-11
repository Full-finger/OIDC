package service

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// UserService 用户服务接口
type UserService interface {
	// RegisterUser 注册用户
	RegisterUser(username, password, email, nickname string) error

	// ActivateUser 激活用户
	ActivateUser(userID uint) error

	// VerifyEmail 验证邮箱
	VerifyEmail(token string) error

	// ResendVerificationEmail 重新发送验证邮件
	ResendVerificationEmail(email string) error

	// AuthenticateUser 用户认证
	AuthenticateUser(username, password string) (*model.User, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(id uint) (*model.User, error)

	// UpdateUserProfile 更新用户资料
	UpdateUserProfile(userID uint, nickname, avatarURL, bio string) error
}