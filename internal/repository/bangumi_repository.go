package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
)

// BangumiRepository Bangumi仓库接口
type BangumiRepository interface {
	// Create 创建Bangumi账号绑定记录
	Create(ctx context.Context, account *model.BangumiAccount) error
	
	// GetByID 根据ID获取Bangumi账号绑定记录
	GetByID(ctx context.Context, id uint) (*model.BangumiAccount, error)
	
	// GetByUserID 根据用户ID获取Bangumi账号绑定记录
	GetByUserID(ctx context.Context, userID uint) (*model.BangumiAccount, error)
	
	// GetByBangumiUserID 根据Bangumi用户ID获取Bangumi账号绑定记录
	GetByBangumiUserID(ctx context.Context, bangumiUserID uint) (*model.BangumiAccount, error)
	
	// Update 更新Bangumi账号绑定记录
	Update(ctx context.Context, account *model.BangumiAccount) error
	
	// DeleteByID 根据ID删除Bangumi账号绑定记录
	DeleteByID(ctx context.Context, id uint) error
}