package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/util"
)

// userService 用户服务实现
type userService struct {
	userRepo    repository.UserRepository
	userHelper  helper.UserHelper
	tokenRepo   repository.VerificationTokenRepository
	emailQueue  util.EmailQueue
}

// NewUserService 创建UserService实例
func NewUserService(
	userRepo repository.UserRepository,
	userHelper helper.UserHelper,
	tokenRepo repository.VerificationTokenRepository,
	emailQueue util.EmailQueue,
) UserService {
	return &userService{
		userRepo:   userRepo,
		userHelper: userHelper,
		tokenRepo:  tokenRepo,
		emailQueue: emailQueue,
	}
}

// generateVerificationToken 生成验证令牌
func (s *userService) generateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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

	// 生成验证令牌
	tokenString, err := s.generateVerificationToken()
	if err != nil {
		return errors.New("生成验证令牌失败")
	}

	// 创建验证令牌记录
	token := &model.VerificationToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时后过期
	}

	// 保存验证令牌
	if err := s.tokenRepo.Create(token); err != nil {
		return errors.New("保存验证令牌失败")
	}

	// 将邮件发送任务加入队列
	emailItem := util.EmailQueueItem{
		Email: user.Email,
		Token: tokenString,
	}
	
	if err := s.emailQueue.Enqueue(emailItem); err != nil {
		// 如果加入队列失败，记录日志但不中断注册流程
		// 在实际项目中，可能需要更完善的错误处理机制
		// 比如重试机制或者将任务保存到数据库中
		// 这里简化处理，仅记录日志
		// 真实场景中应该有专门的邮件发送服务来处理队列中的任务
		// TODO: 实现更完善的错误处理机制
	}

	return nil
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(userID uint) error {
	// TODO: 实现用户激活逻辑
	return nil
}

// VerifyEmail 验证邮箱
func (s *userService) VerifyEmail(token string) error {
	// 根据令牌获取验证令牌记录
	verificationToken, err := s.tokenRepo.GetByToken(token)
	if err != nil {
		return errors.New("无效的验证令牌")
	}

	// 检查令牌是否过期
	if time.Now().After(verificationToken.ExpiresAt) {
		// 删除过期的令牌
		s.tokenRepo.Delete(verificationToken.ID)
		return errors.New("验证令牌已过期")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(verificationToken.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 激活用户
	user.IsActive = true
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户状态失败")
	}

	// 删除已使用的令牌
	if err := s.tokenRepo.Delete(verificationToken.ID); err != nil {
		// 记录日志但不中断流程
		// 在实际项目中，可能需要定期清理这些令牌
	}

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