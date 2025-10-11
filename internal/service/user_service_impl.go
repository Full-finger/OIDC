package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
)

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	userHelper helper.UserHelper
}

// NewUserService 创建UserService实例
func NewUserService(userRepo repository.UserRepository, userHelper helper.UserHelper) UserService {
	return &userService{
		userRepo: userRepo,
		userHelper: userHelper,
	}
}

// RegisterUser 注册用户
func (s *userService) RegisterUser(username, password, email, nickname string) error {
	// 检查用户是否已存在（通过用户名）
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 检查用户是否已存在（通过邮箱）
	_, err = s.userRepo.GetByEmail(email)
	if err == nil {
		return errors.New("邮箱已被注册")
	}

	// 使用bcrypt哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户实体
	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Nickname:     nickname,
		IsActive:     false, // 用户默认未激活，需要邮箱验证
	}

	// 通过Repository创建用户
	if err := s.userRepo.Create(user); err != nil {
		return errors.New("用户创建失败")
	}

	return nil
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(userID uint) error {
	// TODO: 实现用户激活逻辑
	return nil
}

// AuthenticateUser 用户认证
func (s *userService) AuthenticateUser(username, password string) (*model.User, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否已激活
	if !user.IsActive {
		return nil, errors.New("用户未激活，请先验证邮箱")
	}

	// 比对密码哈希
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 认证成功，返回用户信息
	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	// TODO: 实现根据ID获取用户逻辑
	return nil, nil
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(userID uint, nickname, avatarURL, bio string) error {
	// TODO: 实现更新用户资料逻辑
	return nil
}