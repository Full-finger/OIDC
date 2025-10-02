// internal/model/user.go

package model

import (
	"time"
)

// User 表示数据库中的用户模型，与 users 表字段一一对应
type User struct {
	ID            int64     `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	PasswordHash  string    `json:"-" db:"password_hash"` // 敏感字段，不对外暴露
	Email         string    `json:"email" db:"email"`
	Nickname      *string   `json:"nickname,omitempty" db:"nickname"` // 可为空字段使用指针
	AvatarURL     *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Bio           *string   `json:"bio,omitempty" db:"bio"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// SafeUser 返回一个不包含敏感信息（如密码）的安全用户视图
// 通常用于 API 响应
func (u *User) SafeUser() SafeUser {
	return SafeUser{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		Nickname:      u.Nickname,
		AvatarURL:     u.AvatarURL,
		Bio:           u.Bio,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// SafeUser 是 User 的安全子集，用于对外暴露
type SafeUser struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Nickname      *string   `json:"nickname,omitempty"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Bio           *string   `json:"bio,omitempty"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
