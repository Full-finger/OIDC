// Package service defines the service layer interfaces for the OIDC application.
package service

import (
	"context"

	"github.com/Full-finger/OIDC/internal/model"
)

// BangumiService defines the Bangumi service interface
type BangumiService interface {
	IBaseService
	ConvertInterface

	// CreateUserBangumiBinding creates a user Bangumi binding record
	CreateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error

	// GetUserBangumiBindingByUserID gets Bangumi binding record by user ID
	GetUserBangumiBindingByUserID(ctx context.Context, userID int64) (*model.UserBangumiBinding, error)

	// GetUserBangumiBindingByBangumiUserID gets Bangumi binding record by Bangumi user ID
	GetUserBangumiBindingByBangumiUserID(ctx context.Context, bangumiUserID int64) (*model.UserBangumiBinding, error)

	// UpdateUserBangumiBinding updates a user Bangumi binding record
	UpdateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error

	// DeleteUserBangumiBinding deletes a user Bangumi binding record
	DeleteUserBangumiBinding(ctx context.Context, userID int64) error

	// IsUserBangumiBound checks if a user is bound to Bangumi
	IsUserBangumiBound(ctx context.Context, userID int64) (bool, error)

	// RefreshBangumiToken refreshes Bangumi access token
	RefreshBangumiToken(ctx context.Context, userID int64) error

	// SyncBangumiData syncs Bangumi data
	SyncBangumiData(ctx context.Context, userID int64) (*BangumiSyncResult, error)
}

// BangumiSyncResult Bangumi data sync result
type BangumiSyncResult struct {
	NewAnimes          int `json:"new_animes"`
	UpdatedCollections int `json:"updated_collections"`
	TotalCollections   int `json:"total_collections"`
}