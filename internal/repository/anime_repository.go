package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// AnimeRepository 番剧仓库接口
type AnimeRepository interface {
	// Create 创建番剧
	Create(ctx context.Context, anime *model.Anime) error
	
	// GetByID 根据ID获取番剧
	GetByID(ctx context.Context, id uint) (*model.Anime, error)
	
	// GetByTitle 根据标题获取番剧
	GetByTitle(ctx context.Context, title string) (*model.Anime, error)
	
	// Update 更新番剧
	Update(ctx context.Context, anime *model.Anime) error
	
	// DeleteByID 根据ID删除番剧
	DeleteByID(ctx context.Context, id uint) error
	
	// ListByStatus 根据状态列出番剧
	ListByStatus(ctx context.Context, status string) ([]*model.Anime, error)
	
	// Search 搜索番剧
	Search(ctx context.Context, keyword string) ([]*model.Anime, error)
	
	// ListAll 列出所有番剧
	ListAll(ctx context.Context) ([]*model.Anime, error)
}