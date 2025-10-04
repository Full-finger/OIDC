package repository

import (
	"context"
	"database/sql"

	"github.com/Full-finger/OIDC/internal/model"
)

// AnimeRepository 定义番剧相关的数据访问接口
type AnimeRepository interface {
	// Anime相关操作
	CreateAnime(ctx context.Context, anime *model.Anime) error
	GetAnimeByID(ctx context.Context, id int64) (*model.Anime, error)
	GetAnimeByTitle(ctx context.Context, title string) (*model.Anime, error)
	ListAnimes(ctx context.Context) ([]*model.Anime, error)
	UpdateAnime(ctx context.Context, anime *model.Anime) error
	DeleteAnime(ctx context.Context, id int64) error

	// UserCollection相关操作
	CreateCollection(ctx context.Context, collection *model.UserCollection) error
	GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error)
	GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error)
	ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error)
	UpdateCollection(ctx context.Context, collection *model.UserCollection) error
	DeleteCollection(ctx context.Context, id int64) error
	
	// 搜索功能
	SearchAnimes(ctx context.Context, keyword string) ([]*model.Anime, error)
}

// animeRepository 是 AnimeRepository 接口的实现
type animeRepository struct {
	db *sql.DB
}

// NewAnimeRepository 创建一个新的 animeRepository 实例
func NewAnimeRepository(db *sql.DB) AnimeRepository {
	return &animeRepository{
		db: db,
	}
}

// CreateAnime 创建一个新的番剧
func (r *animeRepository) CreateAnime(ctx context.Context, anime *model.Anime) error {
	query := `
		INSERT INTO animes (title, episode_count, director)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		anime.Title,
		anime.EpisodeCount,
		anime.Director,
	).Scan(&anime.ID, &anime.CreatedAt, &anime.UpdatedAt)
}

// GetAnimeByID 根据ID获取番剧
func (r *animeRepository) GetAnimeByID(ctx context.Context, id int64) (*model.Anime, error) {
	query := `
		SELECT id, title, episode_count, director, created_at, updated_at
		FROM animes
		WHERE id = $1`
	
	anime := &model.Anime{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&anime.ID,
		&anime.Title,
		&anime.EpisodeCount,
		&anime.Director,
		&anime.CreatedAt,
		&anime.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return anime, nil
}

// GetAnimeByTitle 根据标题获取番剧
func (r *animeRepository) GetAnimeByTitle(ctx context.Context, title string) (*model.Anime, error) {
	query := `
		SELECT id, title, episode_count, director, created_at, updated_at
		FROM animes
		WHERE title = $1`
	
	anime := &model.Anime{}
	err := r.db.QueryRowContext(ctx, query, title).Scan(
		&anime.ID,
		&anime.Title,
		&anime.EpisodeCount,
		&anime.Director,
		&anime.CreatedAt,
		&anime.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return anime, nil
}

// ListAnimes 获取所有番剧列表
func (r *animeRepository) ListAnimes(ctx context.Context) ([]*model.Anime, error) {
	query := `
		SELECT id, title, episode_count, director, created_at, updated_at
		FROM animes
		ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var animes []*model.Anime
	for rows.Next() {
		anime := &model.Anime{}
		err := rows.Scan(
			&anime.ID,
			&anime.Title,
			&anime.EpisodeCount,
			&anime.Director,
			&anime.CreatedAt,
			&anime.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}
	
	return animes, rows.Err()
}

// UpdateAnime 更新番剧信息
func (r *animeRepository) UpdateAnime(ctx context.Context, anime *model.Anime) error {
	query := `
		UPDATE animes
		SET title = $1, episode_count = $2, director = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		anime.Title,
		anime.EpisodeCount,
		anime.Director,
		anime.ID,
	).Scan(&anime.UpdatedAt)
}

// DeleteAnime 删除番剧
func (r *animeRepository) DeleteAnime(ctx context.Context, id int64) error {
	query := `DELETE FROM animes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CreateCollection 创建用户收藏
func (r *animeRepository) CreateCollection(ctx context.Context, collection *model.UserCollection) error {
	query := `
		INSERT INTO user_collections (user_id, anime_id, type, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		collection.UserID,
		collection.AnimeID,
		collection.Type,
		collection.Rating,
		collection.Comment,
	).Scan(&collection.ID, &collection.CreatedAt, &collection.UpdatedAt)
}

// GetCollectionByID 根据ID获取收藏
func (r *animeRepository) GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error) {
	query := `
		SELECT id, user_id, anime_id, type, rating, comment, created_at, updated_at
		FROM user_collections
		WHERE id = $1`
	
	collection := &model.UserCollection{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.AnimeID,
		&collection.Type,
		&collection.Rating,
		&collection.Comment,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return collection, nil
}

// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取收藏
func (r *animeRepository) GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error) {
	query := `
		SELECT id, user_id, anime_id, type, rating, comment, created_at, updated_at
		FROM user_collections
		WHERE user_id = $1 AND anime_id = $2`
	
	collection := &model.UserCollection{}
	err := r.db.QueryRowContext(ctx, query, userID, animeID).Scan(
		&collection.ID,
		&collection.UserID,
		&collection.AnimeID,
		&collection.Type,
		&collection.Rating,
		&collection.Comment,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return collection, nil
}

// ListCollectionsByUser 获取用户的所有收藏
func (r *animeRepository) ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error) {
	query := `
		SELECT id, user_id, anime_id, type, rating, comment, created_at, updated_at
		FROM user_collections
		WHERE user_id = $1
		ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var collections []*model.UserCollection
	for rows.Next() {
		collection := &model.UserCollection{}
		err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.AnimeID,
			&collection.Type,
			&collection.Rating,
			&collection.Comment,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	
	return collections, rows.Err()
}

// UpdateCollection 更新用户收藏
func (r *animeRepository) UpdateCollection(ctx context.Context, collection *model.UserCollection) error {
	query := `
		UPDATE user_collections
		SET type = $1, rating = $2, comment = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		collection.Type,
		collection.Rating,
		collection.Comment,
		collection.ID,
	).Scan(&collection.UpdatedAt)
}

// DeleteCollection 删除用户收藏
func (r *animeRepository) DeleteCollection(ctx context.Context, id int64) error {
	query := `DELETE FROM user_collections WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// SearchAnimes 根据标题搜索番剧
func (r *animeRepository) SearchAnimes(ctx context.Context, keyword string) ([]*model.Anime, error) {
	query := `
		SELECT id, title, episode_count, director, created_at, updated_at
		FROM animes
		WHERE title ILIKE $1 OR director ILIKE $1
		ORDER BY created_at DESC`
	
	// 添加通配符以支持模糊搜索
	searchTerm := "%" + keyword + "%"
	
	rows, err := r.db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var animes []*model.Anime
	for rows.Next() {
		anime := &model.Anime{}
		err := rows.Scan(
			&anime.ID,
			&anime.Title,
			&anime.EpisodeCount,
			&anime.Director,
			&anime.CreatedAt,
			&anime.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}
	
	return animes, rows.Err()
}