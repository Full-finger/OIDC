package repository

import (
	"context"
	"github.com/Full-finger/OIDC/internal/mapper"
	"github.com/Full-finger/OIDC/internal/model"
)

// animeRepository 番剧仓库实现
type animeRepository struct {
	animeMapper mapper.AnimeMapper
}

// NewAnimeRepository 创建AnimeRepository实例
func NewAnimeRepository() AnimeRepository {
	return &animeRepository{
		animeMapper: mapper.NewAnimeMapper(),
	}
}

// Create 创建番剧
func (r *animeRepository) Create(ctx context.Context, anime *model.Anime) error {
	return r.animeMapper.Save(anime)
}

// GetByID 根据ID获取番剧
func (r *animeRepository) GetByID(ctx context.Context, id uint) (*model.Anime, error) {
	entity, err := r.animeMapper.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	if anime, ok := entity.(*model.Anime); ok {
		return anime, nil
	}
	
	return nil, nil
}

// GetByTitle 根据标题获取番剧
func (r *animeRepository) GetByTitle(ctx context.Context, title string) (*model.Anime, error) {
	return r.animeMapper.GetByTitle(title)
}

// Update 更新番剧
func (r *animeRepository) Update(ctx context.Context, anime *model.Anime) error {
	return r.animeMapper.Update(anime)
}

// DeleteByID 根据ID删除番剧
func (r *animeRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.animeMapper.DeleteByID(id)
}

// ListByStatus 根据状态列出番剧
func (r *animeRepository) ListByStatus(ctx context.Context, status string) ([]*model.Anime, error) {
	return r.animeMapper.GetByStatus(status)
}

// Search 搜索番剧
func (r *animeRepository) Search(ctx context.Context, keyword string) ([]*model.Anime, error) {
	return r.animeMapper.Search(keyword)
}

// ListAll 列出所有番剧
func (r *animeRepository) ListAll(ctx context.Context) ([]*model.Anime, error) {
	entities, err := r.animeMapper.GetAll()
	if err != nil {
		return nil, err
	}
	
	var animes []*model.Anime
	for _, entity := range entities {
		if anime, ok := entity.(*model.Anime); ok {
			animes = append(animes, anime)
		}
	}
	
	return animes, nil
}