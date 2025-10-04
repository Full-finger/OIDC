// internal/repository/oauth_repository.go

package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Full-finger/OIDC/internal/model"
	"github.com/lib/pq"
)

// OAuthRepository 定义OAuth相关数据操作接口
type OAuthRepository interface {
	// Client相关操作
	CreateClient(ctx context.Context, client *model.Client) error
	FindClientByID(ctx context.Context, id int64) (*model.Client, error)
	FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error)
	
	// AuthorizationCode相关操作
	CreateAuthorizationCode(ctx context.Context, code *model.AuthorizationCode) error
	FindAuthorizationCode(ctx context.Context, code string) (*model.AuthorizationCode, error)
	DeleteAuthorizationCode(ctx context.Context, code string) error
	DeleteExpiredAuthorizationCodes(ctx context.Context) error
}

// oauthRepository 是 OAuthRepository 接口的 database/sql 实现
type oauthRepository struct {
	db *sql.DB
}

// NewOAuthRepository 创建一个新的 oauthRepository 实例
func NewOAuthRepository(db *sql.DB) OAuthRepository {
	return &oauthRepository{db: db}
}

// CreateClient 创建新的OAuth客户端
func (r *oauthRepository) CreateClient(ctx context.Context, client *model.Client) error {
	log.Printf("OAuth Repository: Creating client with ID: %s", client.ClientID)
	const query = `
		INSERT INTO oauth_clients (
			client_id, client_secret_hash, name, redirect_uris, scopes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		client.ClientID,
		client.ClientSecretHash,
		client.Name,
		pq.Array(client.RedirectURIs),
		pq.Array(client.Scopes),
		time.Now(),
	).Scan(&client.ID, &client.CreatedAt)
	
	if err != nil {
		log.Printf("OAuth Repository: Failed to create client: %v", err)
		return err
	}
	
	log.Printf("OAuth Repository: Client created successfully with ID: %d", client.ID)
	return err
}

// FindClientByID 根据ID查找客户端
func (r *oauthRepository) FindClientByID(ctx context.Context, id int64) (*model.Client, error) {
	log.Printf("OAuth Repository: Finding client by ID: %d", id)
	return r.findClient(ctx, "id = $1", id)
}

// FindClientByClientID 根据ClientID查找客户端
func (r *oauthRepository) FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error) {
	log.Printf("OAuth Repository: Finding client by ClientID: %s", clientID)
	return r.findClient(ctx, "client_id = $1", clientID)
}

// findClient 是内部辅助方法，避免重复代码
func (r *oauthRepository) findClient(ctx context.Context, whereClause string, args ...interface{}) (*model.Client, error) {
	query := `
		SELECT 
			id, client_id, client_secret_hash, name, redirect_uris, scopes, created_at
		FROM oauth_clients 
		WHERE ` + whereClause + ` LIMIT 1`
		
	log.Printf("OAuth Repository: Executing query: %s with args: %v", query, args)

	client := &model.Client{}
	var redirectURIs pq.StringArray
	var scopes pq.StringArray
	
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&client.ID,
		&client.ClientID,
		&client.ClientSecretHash,
		&client.Name,
		&redirectURIs,
		&scopes,
		&client.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("OAuth Repository: No client found")
			return nil, sql.ErrNoRows
		}
		log.Printf("OAuth Repository: Error finding client: %v", err)
		return nil, err
	}
	
	// 转换 pq.StringArray 到 []string
	client.RedirectURIs = []string(redirectURIs)
	client.Scopes = []string(scopes)
	
	log.Printf("OAuth Repository: Found client - ID: %d, ClientID: %s, Name: %s", client.ID, client.ClientID, client.Name)

	return client, nil
}

// CreateAuthorizationCode 创建新的授权码
func (r *oauthRepository) CreateAuthorizationCode(ctx context.Context, code *model.AuthorizationCode) error {
	log.Printf("OAuth Repository: Creating authorization code: %s for client: %s", code.Code, code.ClientID)
	const query = `
		INSERT INTO oauth_authorization_codes (
			code, client_id, user_id, redirect_uri, scopes, expires_at, code_challenge, code_challenge_method
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`

	_, err := r.db.ExecContext(ctx, query,
		code.Code,
		code.ClientID,
		code.UserID,
		code.RedirectURI,
		pq.Array(code.Scopes),
		code.ExpiresAt,
		code.CodeChallenge,
		code.CodeChallengeMethod,
	)
	
	if err != nil {
		log.Printf("OAuth Repository: Failed to create authorization code: %v", err)
		return err
	}
	
	log.Printf("OAuth Repository: Authorization code created successfully")
	return err
}

// FindAuthorizationCode 根据code查找授权码
func (r *oauthRepository) FindAuthorizationCode(ctx context.Context, code string) (*model.AuthorizationCode, error) {
	log.Printf("OAuth Repository: Finding authorization code: %s", code)
	const query = `
		SELECT 
			code, client_id, user_id, redirect_uri, scopes, expires_at, code_challenge, code_challenge_method
		FROM oauth_authorization_codes 
		WHERE code = $1 LIMIT 1`

	authCode := &model.AuthorizationCode{}
	var scopes pq.StringArray
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&authCode.Code,
		&authCode.ClientID,
		&authCode.UserID,
		&authCode.RedirectURI,
		&scopes,
		&authCode.ExpiresAt,
		&authCode.CodeChallenge,
		&authCode.CodeChallengeMethod,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("OAuth Repository: Authorization code not found: %s", code)
			return nil, sql.ErrNoRows
		}
		log.Printf("OAuth Repository: Error finding authorization code: %v", err)
		return nil, err
	}
	
	// 转换 pq.StringArray 到 []string
	authCode.Scopes = []string(scopes)
	
	log.Printf("OAuth Repository: Found authorization code for client: %s, user: %d", authCode.ClientID, authCode.UserID)

	return authCode, nil
}

// DeleteAuthorizationCode 删除授权码
func (r *oauthRepository) DeleteAuthorizationCode(ctx context.Context, code string) error {
	log.Printf("OAuth Repository: Deleting authorization code: %s", code)
	const query = `DELETE FROM oauth_authorization_codes WHERE code = $1`

	_, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		log.Printf("OAuth Repository: Failed to delete authorization code: %v", err)
		return err
	}
	
	log.Printf("OAuth Repository: Authorization code deleted successfully")
	return err
}

// DeleteExpiredAuthorizationCodes 删除过期的授权码
func (r *oauthRepository) DeleteExpiredAuthorizationCodes(ctx context.Context) error {
	log.Printf("OAuth Repository: Deleting expired authorization codes")
	const query = `DELETE FROM oauth_authorization_codes WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		log.Printf("OAuth Repository: Failed to delete expired authorization codes: %v", err)
		return err
	}
	
	log.Printf("OAuth Repository: Expired authorization codes deleted successfully")
	return err
}