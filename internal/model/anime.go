package model

import (
	"time"
)

// Anime 番剧模型
type Anime struct {
	ID           int64      `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	EpisodeCount *int       `json:"episode_count" db:"episode_count"`
	Director     *string    `json:"director" db:"director"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// UserCollection 用户收藏模型
type UserCollection struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	AnimeID   int64      `json:"anime_id" db:"anime_id"`
	Type      string     `json:"type" db:"type"`
	Rating    *int       `json:"rating" db:"rating"`
	Comment   *string    `json:"comment" db:"comment"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}