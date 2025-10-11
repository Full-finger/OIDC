package handler

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/model"
)

// CollectionHandler 收藏处理器
type CollectionHandler struct {
	collectionService service.CollectionService
}

// NewCollectionHandler 创建CollectionHandler实例
func NewCollectionHandler(collectionService service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
	}
}

// AddToCollectionRequest 添加收藏请求
type AddToCollectionRequest struct {
	AnimeID uint    `json:"anime_id" binding:"required"`
	Status  string  `json:"status" binding:"required"`
	Rating  *float64 `json:"rating"`
	Comment string  `json:"comment"`
}

// AddToCollectionHandler 添加番剧到收藏
func (h *CollectionHandler) AddToCollectionHandler(c *gin.Context) {
	var req AddToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 从上下文中获取用户ID（假设已经在中间件中设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	collection, err := h.collectionService.AddToCollection(
		c.Request.Context(),
		userID.(uint),
		req.AnimeID,
		req.Status,
		req.Rating,
		req.Comment,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to collection"})
		return
	}
	
	c.JSON(http.StatusOK, collection)
}

// GetCollectionHandler 获取用户对某个番剧的收藏
func (h *CollectionHandler) GetCollectionHandler(c *gin.Context) {
	animeIDStr := c.Param("anime_id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	collection, err := h.collectionService.GetCollection(
		c.Request.Context(),
		userID.(uint),
		uint(animeID),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
		return
	}
	
	if collection == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}
	
	c.JSON(http.StatusOK, collection)
}

// UpdateCollectionRequest 更新收藏请求
type UpdateCollectionRequest struct {
	Status   string   `json:"status"`
	Rating   *float64 `json:"rating"`
	Comment  string   `json:"comment"`
	Progress int      `json:"progress"`
	Favorite bool     `json:"favorite"`
}

// UpdateCollectionHandler 更新收藏
func (h *CollectionHandler) UpdateCollectionHandler(c *gin.Context) {
	animeIDStr := c.Param("anime_id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	var req UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 获取当前收藏
	collection, err := h.collectionService.GetCollection(
		c.Request.Context(),
		userID.(uint),
		uint(animeID),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
		return
	}
	
	if collection == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}
	
	// 更新字段
	if req.Status != "" {
		collection.Status = req.Status
	}
	
	if req.Rating != nil {
		collection.Rating = req.Rating
	}
	
	if req.Comment != "" {
		collection.Comment = req.Comment
	}
	
	if req.Favorite {
		collection.Favorite = req.Favorite
	}
	
	// 如果请求包含进度更新，则使用UpdateProgress方法
	if req.Progress > 0 {
		collection, err = h.collectionService.UpdateProgress(
			c.Request.Context(),
			userID.(uint),
			uint(animeID),
			req.Progress,
		)
	} else {
		// 否则使用普通更新方法
		err = h.collectionService.UpdateCollection(c.Request.Context(), collection)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update collection"})
		return
	}
	
	c.JSON(http.StatusOK, collection)
}

// RemoveFromCollectionHandler 从收藏中移除番剧
func (h *CollectionHandler) RemoveFromCollectionHandler(c *gin.Context) {
	animeIDStr := c.Param("anime_id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	err = h.collectionService.RemoveFromCollection(
		c.Request.Context(),
		userID.(uint),
		uint(animeID),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove from collection"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "removed from collection"})
}

// ListUserCollectionsHandler 列出用户的所有收藏
func (h *CollectionHandler) ListUserCollectionsHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	collections, err := h.collectionService.ListUserCollections(
		c.Request.Context(),
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list collections"})
		return
	}
	
	// 填充番剧信息
	var result []*model.Collection
	for _, collection := range collections {
		// 这里应该从番剧服务获取番剧详细信息并填充到collection.Anime字段
		// 为简化示例，我们直接返回收藏列表
		result = append(result, collection)
	}
	
	c.JSON(http.StatusOK, result)
}

// ListUserCollectionsByStatusHandler 根据状态列出用户的收藏
func (h *CollectionHandler) ListUserCollectionsByStatusHandler(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	collections, err := h.collectionService.ListUserCollectionsByStatus(
		c.Request.Context(),
		userID.(uint),
		status,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list collections by status"})
		return
	}
	
	c.JSON(http.StatusOK, collections)
}

// ListUserFavoritesHandler 列出用户的收藏夹
func (h *CollectionHandler) ListUserFavoritesHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	collections, err := h.collectionService.ListUserFavorites(
		c.Request.Context(),
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list favorites"})
		return
	}
	
	c.JSON(http.StatusOK, collections)
}