package service

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// AnimeService 番剧服务接口
type AnimeService interface {
	// CreateAnime 创建番剧
	CreateAnime(ctx context.Context, anime *model.Anime) error
	
	// GetAnimeByID 根据ID获取番剧
	GetAnimeByID(ctx context.Context, id uint) (*model.Anime, error)
	
	// UpdateAnime 更新番剧
	UpdateAnime(ctx context.Context, anime *model.Anime) error
	
	// DeleteAnime 删除番剧
	DeleteAnime(ctx context.Context, id uint) error
	
	// ListAnimes 列出所有番剧
	ListAnimes(ctx context.Context) ([]*model.Anime, error)
	
	// SearchAnimes 搜索番剧
	SearchAnimes(ctx context.Context, keyword string) ([]*model.Anime, error)
	
	// ListAnimesByStatus 根据状态列出番剧
	ListAnimesByStatus(ctx context.Context, status string) ([]*model.Anime, error)
}