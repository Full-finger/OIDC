package mapper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// CollectionMapper 收藏映射器接口
type CollectionMapper interface {
	BaseMapper
	
	// GetByUserID 根据用户ID获取收藏列表
	GetByUserID(userID uint) ([]*model.Collection, error)
	
	// GetByUserIDAndAnimeID 根据用户ID和番剧ID获取收藏
	GetByUserIDAndAnimeID(userID, animeID uint) (*model.Collection, error)
	
	// GetByStatus 根据状态获取用户收藏列表
	GetByStatus(userID uint, status string) ([]*model.Collection, error)
	
	// GetFavorites 获取用户收藏夹
	GetFavorites(userID uint) ([]*model.Collection, error)
}