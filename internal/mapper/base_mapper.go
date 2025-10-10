package mapper

import (
	"context"
)

// IBaseMapper defines the base mapper interface
type IBaseMapper interface {
	// HealthCheck checks the health of the mapper
	HealthCheck() error

	// GetVersion returns the version of the mapper
	GetVersion() string
}

// BaseMapper 定义基础映射器接口
type BaseMapper interface {
	// Create 创建新记录
	Create(ctx context.Context, entity interface{}) error

	// FindByID 根据ID查找记录
	FindByID(ctx context.Context, id int64) (interface{}, error)

	// Update 更新记录
	Update(ctx context.Context, entity interface{}) error

	// Delete 删除记录
	Delete(ctx context.Context, id int64) error

	// List 获取记录列表
	List(ctx context.Context, offset, limit int) ([]interface{}, error)

	// Count 获取记录总数
	Count(ctx context.Context) (int64, error)
}