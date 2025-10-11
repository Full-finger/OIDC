package service

import (
	"context"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

// animeService 番剧服务实现
type animeService struct {
	animeRepo repository.AnimeRepository
}

// NewAnimeService 创建AnimeService实例
func NewAnimeService(animeRepo repository.AnimeRepository) AnimeService {
	return &animeService{
		animeRepo: animeRepo,
	}
}

// CreateAnime 创建番剧
func (s *animeService) CreateAnime(ctx context.Context, anime *model.Anime) error {
	return s.animeRepo.Create(ctx, anime)
}

// GetAnimeByID 根据ID获取番剧
func (s *animeService) GetAnimeByID(ctx context.Context, id uint) (*model.Anime, error) {
	return s.animeRepo.GetByID(ctx, id)
}

// UpdateAnime 更新番剧
func (s *animeService) UpdateAnime(ctx context.Context, anime *model.Anime) error {
	return s.animeRepo.Update(ctx, anime)
}

// DeleteAnime 删除番剧
func (s *animeService) DeleteAnime(ctx context.Context, id uint) error {
	return s.animeRepo.DeleteByID(ctx, id)
}

// ListAnimes 列出所有番剧
func (s *animeService) ListAnimes(ctx context.Context) ([]*model.Anime, error) {
	return s.animeRepo.ListAll(ctx)
}

// SearchAnimes 搜索番剧
func (s *animeService) SearchAnimes(ctx context.Context, keyword string) ([]*model.Anime, error) {
	return s.animeRepo.Search(ctx, keyword)
}

// ListAnimesByStatus 根据状态列出番剧
func (s *animeService) ListAnimesByStatus(ctx context.Context, status string) ([]*model.Anime, error) {
	return s.animeRepo.ListByStatus(ctx, status)
}