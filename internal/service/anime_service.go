package service

import (
	"context"

	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
)

// AnimeService 定义番剧相关的业务逻辑接口
type AnimeService interface {
	// Anime相关操作
	CreateAnime(ctx context.Context, title string, episodeCount *int, director *string) (*model.Anime, error)
	GetAnimeByID(ctx context.Context, id int64) (*model.Anime, error)
	GetAnimeByTitle(ctx context.Context, title string) (*model.Anime, error)
	ListAnimes(ctx context.Context) ([]*model.Anime, error)
	SearchAnimes(ctx context.Context, title string) ([]*model.Anime, error)
	UpdateAnime(ctx context.Context, id int64, title string, episodeCount *int, director *string) (*model.Anime, error)
	DeleteAnime(ctx context.Context, id int64) error

	// UserCollection相关操作
	CreateCollection(ctx context.Context, userID, animeID int64, collectionType string, rating *int, comment *string) (*model.UserCollection, error)
	GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error)
	GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error)
	ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error)
	UpdateCollection(ctx context.Context, id int64, collectionType string, rating *int, comment *string) (*model.UserCollection, error)
	DeleteCollection(ctx context.Context, id int64) error
}

// animeService 是 AnimeService 接口的实现
type animeService struct {
	animeRepo repository.AnimeRepository
}

// NewAnimeService 创建一个新的 animeService 实例
func NewAnimeService(animeRepo repository.AnimeRepository) AnimeService {
	return &animeService{
		animeRepo: animeRepo,
	}
}

// CreateAnime 创建一个新的番剧
func (s *animeService) CreateAnime(ctx context.Context, title string, episodeCount *int, director *string) (*model.Anime, error) {
	anime := &model.Anime{
		Title:        title,
		EpisodeCount: episodeCount,
		Director:     director,
	}
	
	err := s.animeRepo.CreateAnime(ctx, anime)
	if err != nil {
		return nil, err
	}
	
	return anime, nil
}

// GetAnimeByID 根据ID获取番剧
func (s *animeService) GetAnimeByID(ctx context.Context, id int64) (*model.Anime, error) {
	return s.animeRepo.GetAnimeByID(ctx, id)
}

// GetAnimeByTitle 根据标题获取番剧
func (s *animeService) GetAnimeByTitle(ctx context.Context, title string) (*model.Anime, error) {
	return s.animeRepo.GetAnimeByTitle(ctx, title)
}

// ListAnimes 获取所有番剧列表
func (s *animeService) ListAnimes(ctx context.Context) ([]*model.Anime, error) {
	return s.animeRepo.ListAnimes(ctx)
}

// SearchAnimes 根据标题搜索番剧
func (s *animeService) SearchAnimes(ctx context.Context, title string) ([]*model.Anime, error) {
	// 这里可以实现更复杂的搜索逻辑
	// 目前我们简单地通过标题进行模糊搜索
	return s.animeRepo.SearchAnimes(ctx, title)
}

// UpdateAnime 更新番剧信息
func (s *animeService) UpdateAnime(ctx context.Context, id int64, title string, episodeCount *int, director *string) (*model.Anime, error) {
	anime, err := s.animeRepo.GetAnimeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if anime == nil {
		return nil, nil
	}
	
	anime.Title = title
	anime.EpisodeCount = episodeCount
	anime.Director = director
	
	err = s.animeRepo.UpdateAnime(ctx, anime)
	if err != nil {
		return nil, err
	}
	
	return anime, nil
}

// DeleteAnime 删除番剧
func (s *animeService) DeleteAnime(ctx context.Context, id int64) error {
	return s.animeRepo.DeleteAnime(ctx, id)
}

// CreateCollection 创建用户收藏
func (s *animeService) CreateCollection(ctx context.Context, userID, animeID int64, collectionType string, rating *int, comment *string) (*model.UserCollection, error) {
	// 检查是否已经收藏
	existing, err := s.animeRepo.GetCollectionByUserAndAnime(ctx, userID, animeID)
	if err != nil {
		return nil, err
	}
	
	if existing != nil {
		return nil, nil // 已经收藏
	}
	
	collection := &model.UserCollection{
		UserID:  userID,
		AnimeID: animeID,
		Type:    collectionType,
		Rating:  rating,
		Comment: comment,
	}
	
	err = s.animeRepo.CreateCollection(ctx, collection)
	if err != nil {
		return nil, err
	}
	
	return collection, nil
}

// GetCollectionByID 根据ID获取收藏
func (s *animeService) GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error) {
	return s.animeRepo.GetCollectionByID(ctx, id)
}

// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取收藏
func (s *animeService) GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error) {
	return s.animeRepo.GetCollectionByUserAndAnime(ctx, userID, animeID)
}

// ListCollectionsByUser 获取用户的所有收藏
func (s *animeService) ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error) {
	return s.animeRepo.ListCollectionsByUser(ctx, userID)
}

// UpdateCollection 更新用户收藏
func (s *animeService) UpdateCollection(ctx context.Context, id int64, collectionType string, rating *int, comment *string) (*model.UserCollection, error) {
	collection, err := s.animeRepo.GetCollectionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if collection == nil {
		return nil, nil
	}
	
	collection.Type = collectionType
	collection.Rating = rating
	collection.Comment = comment
	
	err = s.animeRepo.UpdateCollection(ctx, collection)
	if err != nil {
		return nil, err
	}
	
	return collection, nil
}

// DeleteCollection 删除用户收藏
func (s *animeService) DeleteCollection(ctx context.Context, id int64) error {
	return s.animeRepo.DeleteCollection(ctx, id)
}