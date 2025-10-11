package model

import (
	"time"
)

// User 用户实体，包含用户基本信息
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Nickname     string    `gorm:"not null" json:"nickname"`
	AvatarURL    string    `gorm:"type:text" json:"avatar_url"`
	Bio          string    `gorm:"type:text" json:"bio"`
	IsActive     bool      `gorm:"default:false" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}