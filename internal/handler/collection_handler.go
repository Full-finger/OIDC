package handler

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
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
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	collection, err := h.collectionService.AddToCollection(
		c.Request.Context(),
		uint(userID),
		req.AnimeID,
		req.Status,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, collection)
}

// ListUserCollectionsHandler 列出用户收藏
func (h *CollectionHandler) ListUserCollectionsHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	collections, err := h.collectionService.ListUserCollections(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, collections)
}

// ListUserCollectionsByStatusHandler 根据状态列出用户的收藏
func (h *CollectionHandler) ListUserCollectionsByStatusHandler(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	collections, err := h.collectionService.ListUserCollectionsByStatus(
		c.Request.Context(),
		uint(userID),
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
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	collections, err := h.collectionService.ListUserFavorites(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list user favorites"})
		return
	}
	
	c.JSON(http.StatusOK, collections)
}

// GetCollectionHandler 获取用户的番剧收藏
func (h *CollectionHandler) GetCollectionHandler(c *gin.Context) {
	animeIDStr := c.Param("anime_id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	collection, err := h.collectionService.GetCollection(c.Request.Context(), uint(userID), uint(animeID))
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
	Status     string  `json:"status"`
	Rating     *float64 `json:"rating"`
	Comment    string  `json:"comment"`
	IsFavorite *bool    `json:"is_favorite"`
}

// UpdateCollectionHandler 更新用户的番剧收藏
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
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	// 首先获取现有的收藏记录
	collection, err := h.collectionService.GetCollection(c.Request.Context(), uint(userID), uint(animeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
		return
	}
	
	if collection == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}
	
	// 更新收藏记录的字段
	if req.Status != "" {
		collection.Status = req.Status
	}
	
	if req.Rating != nil {
		collection.Rating = req.Rating
	}
	
	if req.Comment != "" {
		collection.Comment = req.Comment
	}
	
	if req.IsFavorite != nil {
		collection.Favorite = *req.IsFavorite
	}
	
	// 调用服务更新收藏
	err = h.collectionService.UpdateCollection(c.Request.Context(), collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, collection)
}

// RemoveFromCollectionHandler 从用户的收藏中移除番剧
func (h *CollectionHandler) RemoveFromCollectionHandler(c *gin.Context) {
	animeIDStr := c.Param("anime_id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	
	err = h.collectionService.RemoveFromCollection(c.Request.Context(), uint(userID), uint(animeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "collection removed successfully"})
}
