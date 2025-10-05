// internal/service/user_service.go

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	model "github.com/Full-finger/OIDC/config"
	"github.com/Full-finger/OIDC/internal/repository"
)

// UserService 封装用户相关的业务逻辑
type UserService struct {
	userRepo repository.UserRepository
	// 用于存储最近发送邮件的时间戳，实际项目中应使用 Redis 等外部存储
	lastEmailRequest map[string]time.Time
	emailRequestMutex sync.Mutex
}

// NewUserService 创建 UserService 实例
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		lastEmailRequest: make(map[string]time.Time),
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

// RegisterWithVerification 注册新用户但需要邮箱验证
func (s *UserService) RegisterWithVerification(ctx context.Context, username, email, password string) (*model.SafeUser, error) {
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

	// 4. 创建未激活的用户模型
	user := &model.User{
		Username:      username,
		Email:         email,
		PasswordHash:  string(passwordHash),
		EmailVerified: false, // 用户初始状态为未验证
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

	// 3. 检查邮箱是否已验证
	if !user.EmailVerified {
		return nil, errors.New("email not verified")
	}

	// 4. 返回安全用户数据
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

// CanRequestVerificationEmail 检查是否可以请求验证邮件（频率限制）
func (s *UserService) CanRequestVerificationEmail(email string) bool {
	s.emailRequestMutex.Lock()
	defer s.emailRequestMutex.Unlock()

	if lastRequest, exists := s.lastEmailRequest[email]; exists {
		// 限制每5分钟只能请求一次
		if time.Since(lastRequest) < 5*time.Minute {
			return false
		}
	}

	return true
}

// UpdateLastEmailRequestTime 更新最后一次请求邮件的时间
func (s *UserService) UpdateLastEmailRequestTime(email string) {
	s.emailRequestMutex.Lock()
	defer s.emailRequestMutex.Unlock()

	s.lastEmailRequest[email] = time.Now()
}

// VerifyEmail 验证用户邮箱
func (s *UserService) VerifyEmail(ctx context.Context, userID int64) error {
	// 查找用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 更新用户状态
	user.EmailVerified = true
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}