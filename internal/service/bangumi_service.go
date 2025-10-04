package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Full-finger/OIDC/internal/client"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/utils"
)

// BangumiService 定义Bangumi绑定相关的业务逻辑接口
type BangumiService interface {
	// CreateUserBangumiBinding 创建用户Bangumi绑定记录
	CreateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error
	
	// GetUserBangumiBindingByUserID 根据用户ID获取Bangumi绑定记录
	GetUserBangumiBindingByUserID(ctx context.Context, userID int64) (*model.UserBangumiBinding, error)
	
	// GetUserBangumiBindingByBangumiUserID 根据Bangumi用户ID获取Bangumi绑定记录
	GetUserBangumiBindingByBangumiUserID(ctx context.Context, bangumiUserID int64) (*model.UserBangumiBinding, error)
	
	// UpdateUserBangumiBinding 更新用户Bangumi绑定记录
	UpdateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error
	
	// DeleteUserBangumiBinding 删除用户Bangumi绑定记录
	DeleteUserBangumiBinding(ctx context.Context, userID int64) error
	
	// IsUserBangumiBound 检查用户是否已绑定Bangumi账号
	IsUserBangumiBound(ctx context.Context, userID int64) (bool, error)
	
	// RefreshBangumiToken 刷新Bangumi访问令牌
	RefreshBangumiToken(ctx context.Context, userID int64) error
	
	// SyncBangumiData 同步Bangumi数据
	SyncBangumiData(ctx context.Context, userID int64) (*BangumiSyncResult, error)
}

// BangumiSyncResult Bangumi数据同步结果
type BangumiSyncResult struct {
	NewAnimes          int `json:"new_animes"`
	UpdatedCollections int `json:"updated_collections"`
	TotalCollections   int `json:"total_collections"`
}

// bangumiService 实现BangumiService接口
type bangumiService struct {
	bangumiRepo     repository.BangumiRepository
	animeService    AnimeService
	collectionService CollectionService
}

// NewBangumiService 创建BangumiService实例
func NewBangumiService(bangumiRepo repository.BangumiRepository, animeService AnimeService, collectionService CollectionService) BangumiService {
	return &bangumiService{
		bangumiRepo:     bangumiRepo,
		animeService:    animeService,
		collectionService: collectionService,
	}
}

// CreateUserBangumiBinding 创建用户Bangumi绑定记录
func (s *bangumiService) CreateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error {
	return s.bangumiRepo.CreateUserBangumiBinding(ctx, binding)
}

// GetUserBangumiBindingByUserID 根据用户ID获取Bangumi绑定记录
func (s *bangumiService) GetUserBangumiBindingByUserID(ctx context.Context, userID int64) (*model.UserBangumiBinding, error) {
	return s.bangumiRepo.GetUserBangumiBindingByUserID(ctx, userID)
}

// GetUserBangumiBindingByBangumiUserID 根据Bangumi用户ID获取Bangumi绑定记录
func (s *bangumiService) GetUserBangumiBindingByBangumiUserID(ctx context.Context, bangumiUserID int64) (*model.UserBangumiBinding, error) {
	return s.bangumiRepo.GetUserBangumiBindingByBangumiUserID(ctx, bangumiUserID)
}

// UpdateUserBangumiBinding 更新用户Bangumi绑定记录
func (s *bangumiService) UpdateUserBangumiBinding(ctx context.Context, binding *model.UserBangumiBinding) error {
	return s.bangumiRepo.UpdateUserBangumiBinding(ctx, binding)
}

// DeleteUserBangumiBinding 删除用户Bangumi绑定记录
func (s *bangumiService) DeleteUserBangumiBinding(ctx context.Context, userID int64) error {
	return s.bangumiRepo.DeleteUserBangumiBinding(ctx, userID)
}

// IsUserBangumiBound 检查用户是否已绑定Bangumi账号
func (s *bangumiService) IsUserBangumiBound(ctx context.Context, userID int64) (bool, error) {
	binding, err := s.bangumiRepo.GetUserBangumiBindingByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	
	return binding != nil, nil
}

// RefreshBangumiToken 刷新Bangumi访问令牌
func (s *bangumiService) RefreshBangumiToken(ctx context.Context, userID int64) error {
	// 获取当前绑定信息
	binding, err := s.bangumiRepo.GetUserBangumiBindingByUserID(ctx, userID)
	if err != nil {
		return err
	}
	
	if binding == nil {
		return nil // 用户未绑定Bangumi账号
	}
	
	// 这里应该实现实际的令牌刷新逻辑
	// 由于我们还没有完整的Bangumi OAuth实现，暂时留空
	// 后续会通过Bangumi API使用refresh_token获取新的access_token
	
	return nil
}

// SyncBangumiData 同步Bangumi数据
func (s *bangumiService) SyncBangumiData(ctx context.Context, userID int64) (*BangumiSyncResult, error) {
	// 获取用户Bangumi绑定信息
	binding, err := s.bangumiRepo.GetUserBangumiBindingByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bangumi binding: %w", err)
	}
	
	if binding == nil {
		return nil, fmt.Errorf("user not bound to Bangumi account")
	}
	
	// 获取Bangumi OAuth配置
	config := utils.GetBangumiOAuthConfig()
	
	// 创建Bangumi客户端
	bangumiClient := client.NewBangumiClient(
		config.ClientID,
		config.ClientSecret,
		config.RedirectURI,
		config.TokenURL,
		config.UserInfoURL,
	)
	
	// 获取Bangumi用户信息
	userInfo, err := bangumiClient.GetUserInfo(binding.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bangumi user info: %w", err)
	}
	
	// 获取Bangumi用户收藏列表
	collections, err := bangumiClient.GetUserCollections(binding.AccessToken, userInfo.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bangumi collections: %w", err)
	}
	
	// 初始化同步结果
	result := &BangumiSyncResult{
		NewAnimes:          0,
		UpdatedCollections: 0,
		TotalCollections:   0,
	}
	
	// 遍历Bangumi收藏列表
	for _, bangumiCollection := range collections.Data {
		// 查找或创建番剧记录
		anime, err := s.animeService.GetAnimeByTitle(ctx, bangumiCollection.Subject.NameCN)
		if err != nil {
			log.Printf("Failed to get anime by title %s: %v", bangumiCollection.Subject.NameCN, err)
			continue
		}
		
		// 如果番剧不存在，则创建新番剧
		if anime == nil {
			episodeCount := bangumiCollection.Subject.EpisodeCount()
			anime, err = s.animeService.CreateAnime(ctx, bangumiCollection.Subject.NameCN, &episodeCount, nil)
			if err != nil {
				log.Printf("Failed to create anime %s: %v", bangumiCollection.Subject.NameCN, err)
				continue
			}
			result.NewAnimes++
		}
		
		// 检查用户是否已有该番剧的收藏
		existingCollection, err := s.collectionService.GetCollectionByUserAndAnime(ctx, userID, anime.ID)
		if err != nil {
			log.Printf("Failed to get existing collection for user %d and anime %d: %v", userID, anime.ID, err)
			continue
		}
		
		// 转换Bangumi收藏类型到我们的系统类型
		collectionType := s.convertBangumiTypeToCollectionType(bangumiCollection.Type)
		
		// 如果收藏不存在，则创建新收藏
		if existingCollection == nil {
			newCollection := &model.UserCollection{
				UserID:  userID,
				AnimeID: anime.ID,
				Type:    collectionType,
				Rating:  &bangumiCollection.Rate,
				Comment: &bangumiCollection.Comment,
			}
			
			err = s.collectionService.CreateCollection(ctx, newCollection)
			if err != nil {
				log.Printf("Failed to create collection for user %d and anime %d: %v", userID, anime.ID, err)
				continue
			}
			result.UpdatedCollections++
		} else {
			// 如果收藏已存在，则更新
			existingCollection.Type = collectionType
			existingCollection.Rating = &bangumiCollection.Rate
			existingCollection.Comment = &bangumiCollection.Comment
			
			err = s.collectionService.UpdateCollection(ctx, existingCollection)
			if err != nil {
				log.Printf("Failed to update collection %d: %v", existingCollection.ID, err)
				continue
			}
			result.UpdatedCollections++
		}
		
		result.TotalCollections++
	}
	
	return result, nil
}

// convertBangumiTypeToCollectionType 转换Bangumi收藏类型到我们的系统类型
func (s *bangumiService) convertBangumiTypeToCollectionType(bangumiType int) string {
	// Bangumi收藏类型映射:
	// 1: 想看
	// 2: 看过
	// 3: 在看
	// 4: 搁置
	// 5: 抛弃
	switch bangumiType {
	case 1:
		return "wish"
	case 2:
		return "completed"
	case 3:
		return "watching"
	case 4:
		return "on_hold"
	case 5:
		return "dropped"
	default:
		return "wish"
	}
}