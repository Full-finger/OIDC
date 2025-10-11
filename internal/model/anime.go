package model

import (
	"time"
)

// Anime 番剧实体
type Anime struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null;size:255;index" json:"title"`              // 番剧标题
	Description string    `gorm:"type:text" json:"description"`                      // 番剧描述
	CoverImage  string    `gorm:"size:512" json:"cover_image"`                       // 封面图片URL
	ReleaseDate time.Time `gorm:"index" json:"release_date"`                         // 发布日期
	Episodes    int       `gorm:"default:0" json:"episodes"`                         // 总集数
	Status      string    `gorm:"size:50;default:'upcoming'" json:"status"`          // 状态 (upcoming, airing, finished)
	Rating      float64   `gorm:"default:0.0" json:"rating"`                         // 评分
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`                  // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新时间
}

// Collection 用户番剧收藏实体
type Collection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`                     // 用户ID
	AnimeID   uint      `gorm:"not null;index" json:"anime_id"`                    // 番剧ID
	Status    string    `gorm:"size:50;default:'watching'" json:"status"`          // 收藏状态 (watching, completed, on_hold, dropped, plan_to_watch)
	Rating    *float64  `json:"rating,omitempty"`                                  // 用户评分
	Progress  int       `gorm:"default:0" json:"progress"`                         // 观看进度（已观看集数）
	Comment   string    `gorm:"type:text" json:"comment"`                          // 用户评论
	Favorite  bool      `gorm:"default:false" json:"favorite"`                     // 是否为收藏夹
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                  // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新时间
	
	// 关联
	Anime Anime `gorm:"foreignKey:AnimeID" json:"anime"`                   // 关联的番剧
	User  User  `gorm:"foreignKey:UserID" json:"user"`                     // 关联的用户
}

// TableName 指定Anime表名
func (Anime) TableName() string {
	return "animes"
}

// TableName 指定Collection表名
func (Collection) TableName() string {
	return "collections"
}