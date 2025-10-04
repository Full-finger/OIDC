package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	collectionService service.CollectionService
}

func NewCollectionHandler(collectionService service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
	}
}

// CreateCollection 创建用户收藏
func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	log.Printf("Collection Handler: Processing create collection request")
	
	var req struct {
		AnimeID int64  `json:"anime_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Rating  *int   `json:"rating"`
		Comment *string `json:"comment"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Collection Handler: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 创建收藏对象
	collection := &model.UserCollection{
		UserID:  userID,
		AnimeID: req.AnimeID,
		Type:    req.Type,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	
	// 调用服务创建收藏
	if err := h.collectionService.CreateCollection(c.Request.Context(), collection); err != nil {
		log.Printf("Collection Handler: Failed to create collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	log.Printf("Collection Handler: Successfully created collection with ID: %d", collection.ID)
	c.JSON(http.StatusCreated, collection)
}

// GetCollection 获取用户收藏详情
func (h *CollectionHandler) GetCollection(c *gin.Context) {
	log.Printf("Collection Handler: Processing get collection request")
	
	// 获取路径参数
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Collection Handler: Failed to parse collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	// 调用服务获取收藏
	collection, err := h.collectionService.GetCollectionByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Collection Handler: Failed to get collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Collection Handler: Collection not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	log.Printf("Collection Handler: Successfully retrieved collection with ID: %d", id)
	c.JSON(http.StatusOK, collection)
}

// UpsertCollection 添加或更新用户收藏
func (h *CollectionHandler) UpsertCollection(c *gin.Context) {
	log.Printf("Collection Handler: Processing upsert collection request")
	
	var req struct {
		AnimeID int64  `json:"anime_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Rating  *int   `json:"rating"`
		Comment *string `json:"comment"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Collection Handler: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 检查是否已存在该收藏
	existing, err := h.collectionService.GetCollectionByUserAndAnime(c.Request.Context(), userID, req.AnimeID)
	if err != nil {
		log.Printf("Collection Handler: Failed to check existing collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing collection"})
		return
	}
	
	if existing != nil {
		// 如果已存在，更新收藏
		log.Printf("Collection Handler: Updating existing collection with ID: %d", existing.ID)
		
		if req.Type != "" {
			existing.Type = req.Type
		}
		
		if req.Rating != nil {
			existing.Rating = req.Rating
		}
		
		if req.Comment != nil {
			existing.Comment = req.Comment
		}
		
		if err := h.collectionService.UpdateCollection(c.Request.Context(), existing); err != nil {
			log.Printf("Collection Handler: Failed to update collection: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		log.Printf("Collection Handler: Successfully updated collection with ID: %d", existing.ID)
		c.JSON(http.StatusOK, existing)
		return
	}
	
	// 如果不存在，创建新收藏
	collection := &model.UserCollection{
		UserID:  userID,
		AnimeID: req.AnimeID,
		Type:    req.Type,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	
	if err := h.collectionService.CreateCollection(c.Request.Context(), collection); err != nil {
		log.Printf("Collection Handler: Failed to create collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	log.Printf("Collection Handler: Successfully created collection with ID: %d", collection.ID)
	c.JSON(http.StatusCreated, collection)
}

// ListCollections 获取用户收藏列表，支持筛选
func (h *CollectionHandler) ListCollections(c *gin.Context) {
	log.Printf("Collection Handler: Processing list collections request")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 调用服务获取用户的所有收藏
	collections, err := h.collectionService.ListCollectionsByUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Collection Handler: Failed to list collections: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections"})
		return
	}
	
	// 获取筛选参数
	filterType := c.Query("type")
	filterRating := c.Query("rating")
	
	// 筛选结果
	var filteredCollections []*model.UserCollection
	for _, collection := range collections {
		// 按类型筛选
		if filterType != "" && collection.Type != filterType {
			continue
		}
		
		// 按评分筛选
		if filterRating != "" {
			rating, err := strconv.Atoi(filterRating)
			if err == nil && collection.Rating != nil && *collection.Rating != rating {
				continue
			}
		}
		
		filteredCollections = append(filteredCollections, collection)
	}
	
	log.Printf("Collection Handler: Successfully retrieved %d collections for user ID: %d", len(filteredCollections), userID)
	c.JSON(http.StatusOK, filteredCollections)
}

// UpdateCollection 更新用户收藏
func (h *CollectionHandler) UpdateCollection(c *gin.Context) {
	log.Printf("Collection Handler: Processing update collection request")
	
	// 获取路径参数
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Collection Handler: Failed to parse collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	var req struct {
		Type    string `json:"type"`
		Rating  *int   `json:"rating"`
		Comment *string `json:"comment"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Collection Handler: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// 获取现有收藏
	collection, err := h.collectionService.GetCollectionByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Collection Handler: Failed to get collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Collection Handler: Collection not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	// 从上下文中获取用户ID并验证权限
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 验证用户是否有权限更新此收藏
	if collection.UserID != userID {
		log.Printf("Collection Handler: User %d does not have permission to update collection %d", userID, id)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	
	// 更新收藏字段
	if req.Type != "" {
		collection.Type = req.Type
	}
	
	if req.Rating != nil {
		collection.Rating = req.Rating
	}
	
	if req.Comment != nil {
		collection.Comment = req.Comment
	}
	
	// 调用服务更新收藏
	if err := h.collectionService.UpdateCollection(c.Request.Context(), collection); err != nil {
		log.Printf("Collection Handler: Failed to update collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	log.Printf("Collection Handler: Successfully updated collection with ID: %d", id)
	c.JSON(http.StatusOK, collection)
}

// DeleteCollection 删除用户收藏
func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	log.Printf("Collection Handler: Processing delete collection request")
	
	// 获取路径参数
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Collection Handler: Failed to parse collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	// 获取现有收藏
	collection, err := h.collectionService.GetCollectionByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Collection Handler: Failed to get collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Collection Handler: Collection not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	// 从上下文中获取用户ID并验证权限
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 验证用户是否有权限删除此收藏
	if collection.UserID != userID {
		log.Printf("Collection Handler: User %d does not have permission to delete collection %d", userID, id)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	
	// 调用服务删除收藏
	if err := h.collectionService.DeleteCollection(c.Request.Context(), id); err != nil {
		log.Printf("Collection Handler: Failed to delete collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}
	
	log.Printf("Collection Handler: Successfully deleted collection with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}

// DeleteCollectionByAnimeID 删除用户对指定番剧的收藏
func (h *CollectionHandler) DeleteCollectionByAnimeID(c *gin.Context) {
	log.Printf("Collection Handler: Processing delete collection by anime ID request")
	
	// 获取路径参数
	animeID, err := strconv.ParseInt(c.Param("anime_id"), 10, 64)
	if err != nil {
		log.Printf("Collection Handler: Failed to parse anime ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Collection Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Collection Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 获取现有收藏
	collection, err := h.collectionService.GetCollectionByUserAndAnime(c.Request.Context(), userID, animeID)
	if err != nil {
		log.Printf("Collection Handler: Failed to get collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Collection Handler: Collection not found for user ID: %d and anime ID: %d", userID, animeID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	// 调用服务删除收藏
	if err := h.collectionService.DeleteCollection(c.Request.Context(), collection.ID); err != nil {
		log.Printf("Collection Handler: Failed to delete collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}
	
	log.Printf("Collection Handler: Successfully deleted collection with ID: %d", collection.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}