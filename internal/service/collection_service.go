package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

type CollectionService interface {
	// CreateCollection 创建用户收藏
	CreateCollection(ctx context.Context, collection *model.UserCollection) error
	
	// GetCollectionByID 根据ID获取用户收藏
	GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error)
	
	// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取用户收藏
	GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error)
	
	// ListCollectionsByUser 获取用户的所有收藏
	ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error)
	
	// UpdateCollection 更新用户收藏
	UpdateCollection(ctx context.Context, collection *model.UserCollection) error
	
	// DeleteCollection 删除用户收藏
	DeleteCollection(ctx context.Context, id int64) error
}

type collectionService struct {
	collectionRepo repository.CollectionRepository
}

func NewCollectionService(collectionRepo repository.CollectionRepository) CollectionService {
	return &collectionService{
		collectionRepo: collectionRepo,
	}
}

// CreateCollection 创建用户收藏
func (s *collectionService) CreateCollection(ctx context.Context, collection *model.UserCollection) error {
	// 检查用户是否已经收藏了该番剧
	existing, err := s.collectionRepo.GetCollectionByUserAndAnime(ctx, collection.UserID, collection.AnimeID)
	if err != nil {
		return err
	}
	
	if existing != nil {
		return errors.New("collection already exists for this user and anime")
	}
	
	// 验证评分范围
	if collection.Rating != nil && (*collection.Rating < 1 || *collection.Rating > 10) {
		return errors.New("rating must be between 1 and 10")
	}
	
	// 创建收藏
	return s.collectionRepo.CreateCollection(ctx, collection)
}

// GetCollectionByID 根据ID获取用户收藏
func (s *collectionService) GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error) {
	return s.collectionRepo.GetCollectionByID(ctx, id)
}

// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取用户收藏
func (s *collectionService) GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error) {
	return s.collectionRepo.GetCollectionByUserAndAnime(ctx, userID, animeID)
}

// ListCollectionsByUser 获取用户的所有收藏
func (s *collectionService) ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error) {
	return s.collectionRepo.ListCollectionsByUser(ctx, userID)
}

// UpdateCollection 更新用户收藏
func (s *collectionService) UpdateCollection(ctx context.Context, collection *model.UserCollection) error {
	// 验证评分范围
	if collection.Rating != nil && (*collection.Rating < 1 || *collection.Rating > 10) {
		return errors.New("rating must be between 1 and 10")
	}
	
	// 检查收藏是否存在
	existing, err := s.collectionRepo.GetCollectionByID(ctx, collection.ID)
	if err != nil {
		return err
	}
	
	if existing == nil {
		return sql.ErrNoRows
	}
	
	// 更新收藏
	return s.collectionRepo.UpdateCollection(ctx, collection)
}

// DeleteCollection 删除用户收藏
func (s *collectionService) DeleteCollection(ctx context.Context, id int64) error {
	// 检查收藏是否存在
	existing, err := s.collectionRepo.GetCollectionByID(ctx, id)
	if err != nil {
		return err
	}
	
	if existing == nil {
		return sql.ErrNoRows
	}
	
	// 删除收藏
	return s.collectionRepo.DeleteCollection(ctx, id)
}