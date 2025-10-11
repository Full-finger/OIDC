package mapper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// BangumiMapper Bangumi映射器接口
type BangumiMapper interface {
	BaseMapper
	
	// GetByUserID 根据用户ID获取Bangumi账号绑定记录
	GetByUserID(userID uint) (*model.BangumiAccount, error)
	
	// GetByBangumiUserID 根据Bangumi用户ID获取Bangumi账号绑定记录
	GetByBangumiUserID(bangumiUserID uint) (*model.BangumiAccount, error)
}