package preprocessor

import (
	"errors"
	"regexp"
	"strings"
	"github.com/Full-finger/OIDC/internal/model"
)

// PreprocessUserRegistration 绑定、校验和整理前端传入的注册数据
func PreprocessUserRegistration(username, password, email, nickname string) (*model.User, error) {
	// 数据清理
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))
	nickname = strings.TrimSpace(nickname)

	// 数据校验
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validateNickname(nickname); err != nil {
		return nil, err
	}

	// 创建用户实体
	user := &model.User{
		Username: username,
		Email:    email,
		Nickname: nickname,
		IsActive: false, // 用户默认未激活，需要邮箱验证
	}

	return user, nil
}

// validateUsername 校验用户名
func validateUsername(username string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}

	if len(username) < 3 || len(username) > 30 {
		return errors.New("用户名长度必须在3-30个字符之间")
	}

	// 用户名只能包含字母、数字、下划线和连字符
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		return errors.New("用户名只能包含字母、数字、下划线和连字符")
	}

	return nil
}

// validatePassword 校验密码
func validatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	if len(password) < 6 {
		return errors.New("密码长度不能少于6个字符")
	}

	if len(password) > 128 {
		return errors.New("密码长度不能超过128个字符")
	}

	return nil
}

// validateEmail 校验邮箱
func validateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	// 简单的邮箱格式校验
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// validateNickname 校验昵称
func validateNickname(nickname string) error {
	if nickname == "" {
		return errors.New("昵称不能为空")
	}

	if len(nickname) > 50 {
		return errors.New("昵称长度不能超过50个字符")
	}

	return nil
}