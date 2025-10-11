package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/model"
)

// bangumiRepository Bangumi仓库实现
type bangumiRepository struct {
	bangumiMapper mapper.BangumiMapper
}

// NewBangumiRepository 创建BangumiRepository实例
func NewBangumiRepository() BangumiRepository {
	return &bangumiRepository{
		bangumiMapper: mapper.NewBangumiMapper(),
	}
}

// Create 创建Bangumi账号绑定记录
func (r *bangumiRepository) Create(ctx context.Context, account *model.BangumiAccount) error {
	return r.bangumiMapper.Save(account)
}

// GetByID 根据ID获取Bangumi账号绑定记录
func (r *bangumiRepository) GetByID(ctx context.Context, id uint) (*model.BangumiAccount, error) {
	entity, err := r.bangumiMapper.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	if account, ok := entity.(*model.BangumiAccount); ok {
		return account, nil
	}
	
	return nil, nil
}

// GetByUserID 根据用户ID获取Bangumi账号绑定记录
func (r *bangumiRepository) GetByUserID(ctx context.Context, userID uint) (*model.BangumiAccount, error) {
	return r.bangumiMapper.GetByUserID(userID)
}

// GetByBangumiUserID 根据Bangumi用户ID获取Bangumi账号绑定记录
func (r *bangumiRepository) GetByBangumiUserID(ctx context.Context, bangumiUserID uint) (*model.BangumiAccount, error) {
	return r.bangumiMapper.GetByBangumiUserID(bangumiUserID)
}

// Update 更新Bangumi账号绑定记录
func (r *bangumiRepository) Update(ctx context.Context, account *model.BangumiAccount) error {
	return r.bangumiMapper.Update(account)
}

// DeleteByID 根据ID删除Bangumi账号绑定记录
func (r *bangumiRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.bangumiMapper.DeleteByID(id)
}