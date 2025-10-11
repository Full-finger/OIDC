package model

import (
	"time"
)

// BangumiAccount Bangumi账号绑定实体
type BangumiAccount struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index;unique" json:"user_id"`           // 关联的用户ID
	BangumiUserID   uint      `gorm:"not null;index" json:"bangumi_user_id"`          // Bangumi平台用户ID
	AccessToken     string    `gorm:"not null" json:"access_token"`                   // Bangumi访问令牌
	RefreshToken    string    `gorm:"not null" json:"refresh_token"`                  // Bangumi刷新令牌
	TokenExpiresAt  time.Time `gorm:"not null" json:"token_expires_at"`               // 令牌过期时间
	Scope           string    `gorm:"type:text" json:"scope"`                         // 授权范围
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`               // 创建时间
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`               // 更新时间
	
	// 关联
	User User `gorm:"foreignKey:UserID" json:"user"`                  // 关联的用户
}

// TableName 指定BangumiAccount表名
func (BangumiAccount) TableName() string {
	return "bangumi_accounts"
}