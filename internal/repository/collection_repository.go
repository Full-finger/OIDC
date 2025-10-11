package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// CollectionRepository 收藏仓库接口
type CollectionRepository interface {
	// Create 创建收藏
	Create(ctx context.Context, collection *model.Collection) error
	
	// GetByID 根据ID获取收藏
	GetByID(ctx context.Context, id uint) (*model.Collection, error)
	
	// GetByUserIDAndAnimeID 根据用户ID和番剧ID获取收藏
	GetByUserIDAndAnimeID(ctx context.Context, userID, animeID uint) (*model.Collection, error)
	
	// Update 更新收藏
	Update(ctx context.Context, collection *model.Collection) error
	
	// DeleteByID 根据ID删除收藏
	DeleteByID(ctx context.Context, id uint) error
	
	// ListByUserID 根据用户ID列出收藏
	ListByUserID(ctx context.Context, userID uint) ([]*model.Collection, error)
	
	// ListByUserIDAndStatus 根据用户ID和状态列出收藏
	ListByUserIDAndStatus(ctx context.Context, userID uint, status string) ([]*model.Collection, error)
	
	// ListFavorites 列出用户收藏夹
	ListFavorites(ctx context.Context, userID uint) ([]*model.Collection, error)
}