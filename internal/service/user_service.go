// internal/service/user_service.go

package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	model "github.com/Full-finger/OIDC/config"
	"github.com/Full-finger/OIDC/internal/repository"
)

// UserService 封装用户相关的业务逻辑
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建 UserService 实例
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register 注册新用户
// - 对密码进行 bcrypt 哈希
// - 检查用户名和邮箱是否已存在
// - 创建用户
func (s *UserService) Register(ctx context.Context, username, email, password string) (*model.SafeUser, error) {
	// 1. 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(ctx, username); err == nil {
		return nil, errors.New("username already exists")
	}
	// 2. 检查邮箱是否已存在
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 3. 哈希密码（使用 bcrypt）
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. 创建用户模型
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. 返回安全用户数据（不含密码）
	safeUser := user.SafeUser()
	return &safeUser, nil
}

// Login 用户登录
// - 验证用户名和密码
// - 返回安全用户信息
func (s *UserService) Login(ctx context.Context, username, password string) (*model.SafeUser, error) {
	// 1. 查找用户（通过用户名或邮箱均可，这里按用户名）
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		// 可选：也尝试用邮箱登录
		if user, err = s.userRepo.FindByEmail(ctx, username); err != nil {
			return nil, errors.New("invalid username or password")
		}
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 3. 返回安全用户数据
	safeUser := user.SafeUser()
	return &safeUser, nil
}

// UpdateProfile 更新用户资料（不包括密码）
func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID int64,
	nickname *string,
	avatarURL *string,
	bio *string,
) (*model.SafeUser, error) {
	// 1. 获取当前用户（确保用户存在）
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 更新字段
	user.Nickname = nickname
	user.AvatarURL = avatarURL
	user.Bio = bio

	// 3. 保存到数据库
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// 4. 返回安全用户数据
	safeUser := user.SafeUser()
	return &safeUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
