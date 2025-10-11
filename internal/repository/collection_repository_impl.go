package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/model"
)

// collectionRepository 收藏仓库实现
type collectionRepository struct {
	collectionMapper mapper.CollectionMapper
}

// NewCollectionRepository 创建CollectionRepository实例
func NewCollectionRepository() CollectionRepository {
	return &collectionRepository{
		collectionMapper: mapper.NewCollectionMapper(),
	}
}

// Create 创建收藏
func (r *collectionRepository) Create(ctx context.Context, collection *model.Collection) error {
	return r.collectionMapper.Save(collection)
}

// GetByID 根据ID获取收藏
func (r *collectionRepository) GetByID(ctx context.Context, id uint) (*model.Collection, error) {
	entity, err := r.collectionMapper.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	if collection, ok := entity.(*model.Collection); ok {
		return collection, nil
	}
	
	return nil, nil
}

// GetByUserIDAndAnimeID 根据用户ID和番剧ID获取收藏
func (r *collectionRepository) GetByUserIDAndAnimeID(ctx context.Context, userID, animeID uint) (*model.Collection, error) {
	return r.collectionMapper.GetByUserIDAndAnimeID(userID, animeID)
}

// Update 更新收藏
func (r *collectionRepository) Update(ctx context.Context, collection *model.Collection) error {
	return r.collectionMapper.Update(collection)
}

// DeleteByID 根据ID删除收藏
func (r *collectionRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.collectionMapper.DeleteByID(id)
}

// ListByUserID 根据用户ID列出收藏
func (r *collectionRepository) ListByUserID(ctx context.Context, userID uint) ([]*model.Collection, error) {
	return r.collectionMapper.GetByUserID(userID)
}

// ListByUserIDAndStatus 根据用户ID和状态列出收藏
func (r *collectionRepository) ListByUserIDAndStatus(ctx context.Context, userID uint, status string) ([]*model.Collection, error) {
	return r.collectionMapper.GetByStatus(userID, status)
}

// ListFavorites 列出用户收藏夹
func (r *collectionRepository) ListFavorites(ctx context.Context, userID uint) ([]*model.Collection, error) {
	return r.collectionMapper.GetFavorites(userID)
}