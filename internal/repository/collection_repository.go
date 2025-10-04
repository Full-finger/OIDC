package repository

import (
	"context"
	"database/sql"
	"github.com/Full-finger/OIDC/internal/model"
)

type CollectionRepository interface {
	// CreateCollection 创建用户收藏
	CreateCollection(ctx context.Context, collection *model.UserCollection) error
	
	// GetCollectionByID 根据ID获取用户收藏
	GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error)
	
	// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取用户收藏
	GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error)
	
	// ListCollectionsByUser 获取用户的所有收藏
	ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error)
	
	// UpdateCollection 更新用户收藏
	UpdateCollection(ctx context.Context, collection *model.UserCollection) error
	
	// DeleteCollection 删除用户收藏
	DeleteCollection(ctx context.Context, id int64) error
}

type collectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

// CreateCollection 创建用户收藏
func (r *collectionRepository) CreateCollection(ctx context.Context, collection *model.UserCollection) error {
	query := `
		INSERT INTO user_collections (user_id, anime_id, type, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		collection.UserID,
		collection.AnimeID,
		collection.Type,
		collection.Rating,
		collection.Comment).Scan(&collection.ID, &collection.CreatedAt, &collection.UpdatedAt)
	
	return err
}

// GetCollectionByID 根据ID获取用户收藏
func (r *collectionRepository) GetCollectionByID(ctx context.Context, id int64) (*model.UserCollection, error) {
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

// GetCollectionByUserAndAnime 根据用户ID和番剧ID获取用户收藏
func (r *collectionRepository) GetCollectionByUserAndAnime(ctx context.Context, userID, animeID int64) (*model.UserCollection, error) {
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
func (r *collectionRepository) ListCollectionsByUser(ctx context.Context, userID int64) ([]*model.UserCollection, error) {
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
func (r *collectionRepository) UpdateCollection(ctx context.Context, collection *model.UserCollection) error {
	query := `
		UPDATE user_collections
		SET type = $1, rating = $2, comment = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		collection.Type,
		collection.Rating,
		collection.Comment,
		collection.ID).Scan(&collection.UpdatedAt)
	
	return err
}

// DeleteCollection 删除用户收藏
func (r *collectionRepository) DeleteCollection(ctx context.Context, id int64) error {
	query := `DELETE FROM user_collections WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}