package service

import (
	"context"
	"time"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

// collectionService 收藏服务实现
type collectionService struct {
	collectionRepo repository.CollectionRepository
	animeRepo      repository.AnimeRepository
}

// NewCollectionService 创建CollectionService实例
func NewCollectionService(collectionRepo repository.CollectionRepository, animeRepo repository.AnimeRepository) CollectionService {
	return &collectionService{
		collectionRepo: collectionRepo,
		animeRepo:      animeRepo,
	}
}

// AddToCollection 添加番剧到收藏
func (s *collectionService) AddToCollection(ctx context.Context, userID, animeID uint, status string, rating *float64, comment string) (*model.Collection, error) {
	// 检查番剧是否存在
	_, err := s.animeRepo.GetByID(ctx, animeID)
	if err != nil {
		return nil, err
	}
	
	// 创建收藏记录
	collection := &model.Collection{
		UserID:    userID,
		AnimeID:   animeID,
		Status:    status,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err = s.collectionRepo.Create(ctx, collection)
	if err != nil {
		return nil, err
	}
	
	return collection, nil
}

// GetCollection 获取用户对某个番剧的收藏
func (s *collectionService) GetCollection(ctx context.Context, userID, animeID uint) (*model.Collection, error) {
	return s.collectionRepo.GetByUserIDAndAnimeID(ctx, userID, animeID)
}

// UpdateCollection 更新收藏
func (s *collectionService) UpdateCollection(ctx context.Context, collection *model.Collection) error {
	collection.UpdatedAt = time.Now()
	return s.collectionRepo.Update(ctx, collection)
}

// RemoveFromCollection 从收藏中移除番剧
func (s *collectionService) RemoveFromCollection(ctx context.Context, userID, animeID uint) error {
	collection, err := s.collectionRepo.GetByUserIDAndAnimeID(ctx, userID, animeID)
	if err != nil {
		return err
	}
	
	if collection == nil {
		return nil // 收藏不存在，直接返回
	}
	
	return s.collectionRepo.DeleteByID(ctx, collection.ID)
}

// ListUserCollections 列出用户的所有收藏
func (s *collectionService) ListUserCollections(ctx context.Context, userID uint) ([]*model.Collection, error) {
	return s.collectionRepo.ListByUserID(ctx, userID)
}

// ListUserCollectionsByStatus 根据状态列出用户的收藏
func (s *collectionService) ListUserCollectionsByStatus(ctx context.Context, userID uint, status string) ([]*model.Collection, error) {
	return s.collectionRepo.ListByUserIDAndStatus(ctx, userID, status)
}

// ListUserFavorites 列出用户的收藏夹
func (s *collectionService) ListUserFavorites(ctx context.Context, userID uint) ([]*model.Collection, error) {
	return s.collectionRepo.ListFavorites(ctx, userID)
}

// UpdateProgress 更新观看进度
func (s *collectionService) UpdateProgress(ctx context.Context, userID, animeID uint, progress int) (*model.Collection, error) {
	collection, err := s.collectionRepo.GetByUserIDAndAnimeID(ctx, userID, animeID)
	if err != nil {
		return nil, err
	}
	
	if collection == nil {
		return nil, nil // 收藏不存在
	}
	
	collection.Progress = progress
	collection.UpdatedAt = time.Now()
	
	err = s.collectionRepo.Update(ctx, collection)
	if err != nil {
		return nil, err
	}
	
	return collection, nil
}