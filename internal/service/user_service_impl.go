package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
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

	// 检查是否跳过邮箱验证
	skipEmailVerification := os.Getenv("SKIP_EMAIL_VERIFICATION") == "true"
	
	// 创建用户实体
	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Nickname:     nickname,
		IsActive:     skipEmailVerification, // 如果跳过邮箱验证，则用户默认激活
	}

	// 通过Repository创建用户
	if err := s.userRepo.Create(user); err != nil {
		return errors.New("用户创建失败")
	}

	// 如果跳过邮箱验证，则直接返回
	if skipEmailVerification {
		return nil
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
	// 更新用户激活状态
	if err := s.userRepo.UpdateActivationStatus(userID, true); err != nil {
		return errors.New("更新用户激活状态失败")
	}
	
	// 删除验证令牌
	// 先根据用户ID获取令牌
	token, err := s.tokenRepo.GetByUserID(userID)
	if err == nil && token != nil {
		// 如果令牌存在，则删除它
		if err := s.tokenRepo.DeleteByToken(token.Token); err != nil {
			// 记录日志但不中断激活流程
			// 在实际项目中，可能需要更完善的错误处理机制
		}
	}
	
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

// ResendVerificationEmail 重新发送验证邮件
func (s *userService) ResendVerificationEmail(email string) error {
	// 查找用户
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查用户是否已激活
	if user.IsActive {
		return errors.New("用户已激活，无需重新发送验证邮件")
	}

	// 生成新的验证令牌
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
		return errors.New("邮件发送任务加入队列失败")
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

	// 检查是否跳过邮箱验证
	skipEmailVerification := os.Getenv("SKIP_EMAIL_VERIFICATION") == "true"
	
	// 检查用户是否已激活（除非跳过邮箱验证）
	if !skipEmailVerification && !user.IsActive {
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
	return s.userRepo.GetByID(id)
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(userID uint, nickname, avatarURL, bio string) error {
	// 查找用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新用户资料
	user.Nickname = nickname
	user.AvatarURL = avatarURL
	user.Bio = bio

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户资料失败")
	}

	return nil
}