package mapper

import (
	"github.com/Full-finger/OIDC/internal/model"
)

// AnimeMapper 番剧映射器接口
type AnimeMapper interface {
	BaseMapper
	
	// GetByTitle 根据标题获取番剧
	GetByTitle(title string) (*model.Anime, error)
	
	// GetByStatus 根据状态获取番剧列表
	GetByStatus(status string) ([]*model.Anime, error)
	
	// Search 搜索番剧
	Search(keyword string) ([]*model.Anime, error)
}