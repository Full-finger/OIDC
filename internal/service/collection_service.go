package service

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// CollectionService 收藏服务接口
type CollectionService interface {
	// AddToCollection 添加番剧到收藏
	AddToCollection(ctx context.Context, userID, animeID uint, status string, rating *float64, comment string) (*model.Collection, error)
	
	// GetCollection 获取用户对某个番剧的收藏
	GetCollection(ctx context.Context, userID, animeID uint) (*model.Collection, error)
	
	// UpdateCollection 更新收藏
	UpdateCollection(ctx context.Context, collection *model.Collection) error
	
	// RemoveFromCollection 从收藏中移除番剧
	RemoveFromCollection(ctx context.Context, userID, animeID uint) error
	
	// ListUserCollections 列出用户的所有收藏
	ListUserCollections(ctx context.Context, userID uint) ([]*model.Collection, error)
	
	// ListUserCollectionsByStatus 根据状态列出用户的收藏
	ListUserCollectionsByStatus(ctx context.Context, userID uint, status string) ([]*model.Collection, error)
	
	// ListUserFavorites 列出用户的收藏夹
	ListUserFavorites(ctx context.Context, userID uint) ([]*model.Collection, error)
	
	// UpdateProgress 更新观看进度
	UpdateProgress(ctx context.Context, userID, animeID uint, progress int) (*model.Collection, error)
}