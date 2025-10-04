package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// AnimeHandler 处理番剧相关的HTTP请求
type AnimeHandler struct {
	animeService service.AnimeService
}

// NewAnimeHandler 创建一个新的 AnimeHandler 实例
func NewAnimeHandler(animeService service.AnimeService) *AnimeHandler {
	return &AnimeHandler{
		animeService: animeService,
	}
}

// CreateAnime 创建番剧
func (h *AnimeHandler) CreateAnime(c *gin.Context) {
	log.Printf("Anime Handler: Processing create anime request")
	
	var req struct {
		Title        string `json:"title" binding:"required"`
		EpisodeCount *int   `json:"episode_count"`
		Director     *string `json:"director"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Anime Handler: Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	anime, err := h.animeService.CreateAnime(
		c.Request.Context(),
		req.Title,
		req.EpisodeCount,
		req.Director,
	)
	
	if err != nil {
		log.Printf("Anime Handler: Failed to create anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create anime"})
		return
	}
	
	log.Printf("Anime Handler: Successfully created anime with ID: %d", anime.ID)
	c.JSON(http.StatusCreated, anime)
}

// GetAnime 获取番剧详情
func (h *AnimeHandler) GetAnime(c *gin.Context) {
	log.Printf("Anime Handler: Processing get anime request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid anime ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}
	
	anime, err := h.animeService.GetAnimeByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Anime Handler: Failed to get anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get anime"})
		return
	}
	
	if anime == nil {
		log.Printf("Anime Handler: Anime not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	
	log.Printf("Anime Handler: Successfully retrieved anime with ID: %d", id)
	c.JSON(http.StatusOK, anime)
}

// ListAnimes 获取番剧列表
func (h *AnimeHandler) ListAnimes(c *gin.Context) {
	log.Printf("Anime Handler: Processing list animes request")
	
	// 检查是否有title参数用于搜索
	title := c.Query("title")
	
	var animes []*model.Anime
	var err error
	
	if title != "" {
		// 如果提供了title参数，则进行搜索
		log.Printf("Anime Handler: Searching animes with title: %s", title)
		animes, err = h.animeService.SearchAnimes(c.Request.Context(), title)
	} else {
		// 否则列出所有番剧
		log.Printf("Anime Handler: Listing all animes")
		animes, err = h.animeService.ListAnimes(c.Request.Context())
	}
	
	if err != nil {
		log.Printf("Anime Handler: Failed to list/search animes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list animes"})
		return
	}
	
	log.Printf("Anime Handler: Successfully retrieved %d animes", len(animes))
	c.JSON(http.StatusOK, animes)
}

// UpdateAnime 更新番剧
func (h *AnimeHandler) UpdateAnime(c *gin.Context) {
	log.Printf("Anime Handler: Processing update anime request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid anime ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}
	
	var req struct {
		Title        string `json:"title" binding:"required"`
		EpisodeCount *int   `json:"episode_count"`
		Director     *string `json:"director"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Anime Handler: Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	anime, err := h.animeService.UpdateAnime(
		c.Request.Context(),
		id,
		req.Title,
		req.EpisodeCount,
		req.Director,
	)
	
	if err != nil {
		log.Printf("Anime Handler: Failed to update anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update anime"})
		return
	}
	
	if anime == nil {
		log.Printf("Anime Handler: Anime not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	
	log.Printf("Anime Handler: Successfully updated anime with ID: %d", id)
	c.JSON(http.StatusOK, anime)
}

// DeleteAnime 删除番剧
func (h *AnimeHandler) DeleteAnime(c *gin.Context) {
	log.Printf("Anime Handler: Processing delete anime request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid anime ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}
	
	err = h.animeService.DeleteAnime(c.Request.Context(), id)
	if err != nil {
		log.Printf("Anime Handler: Failed to delete anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete anime"})
		return
	}
	
	log.Printf("Anime Handler: Successfully deleted anime with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Anime deleted successfully"})
}

// CreateCollection 创建收藏
func (h *AnimeHandler) CreateCollection(c *gin.Context) {
	log.Printf("Anime Handler: Processing create collection request")
	
	// 从JWT中间件获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Anime Handler: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Anime Handler: failed to convert userID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	
	var req struct {
		AnimeID int64  `json:"anime_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Rating  *int   `json:"rating"`
		Comment *string `json:"comment"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Anime Handler: Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证评分范围
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 10) {
		log.Printf("Anime Handler: Invalid rating value: %d", *req.Rating)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 10"})
		return
	}
	
	collection, err := h.animeService.CreateCollection(
		c.Request.Context(),
		userID,
		req.AnimeID,
		req.Type,
		req.Rating,
		req.Comment,
	)
	
	if err != nil {
		log.Printf("Anime Handler: Failed to create collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Anime Handler: Collection already exists for user %d and anime %d", userID, req.AnimeID)
		c.JSON(http.StatusConflict, gin.H{"error": "Collection already exists"})
		return
	}
	
	log.Printf("Anime Handler: Successfully created collection with ID: %d", collection.ID)
	c.JSON(http.StatusCreated, collection)
}

// GetCollection 获取收藏详情
func (h *AnimeHandler) GetCollection(c *gin.Context) {
	log.Printf("Anime Handler: Processing get collection request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	collection, err := h.animeService.GetCollectionByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Anime Handler: Failed to get collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Anime Handler: Collection not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	log.Printf("Anime Handler: Successfully retrieved collection with ID: %d", id)
	c.JSON(http.StatusOK, collection)
}

// ListCollections 获取用户收藏列表
func (h *AnimeHandler) ListCollections(c *gin.Context) {
	log.Printf("Anime Handler: Processing list collections request")
	
	// 从JWT中间件获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Anime Handler: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Anime Handler: failed to convert userID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	
	collections, err := h.animeService.ListCollectionsByUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Anime Handler: Failed to list collections: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections"})
		return
	}
	
	log.Printf("Anime Handler: Successfully retrieved %d collections for user %d", len(collections), userID)
	c.JSON(http.StatusOK, collections)
}

// UpdateCollection 更新收藏
func (h *AnimeHandler) UpdateCollection(c *gin.Context) {
	log.Printf("Anime Handler: Processing update collection request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	var req struct {
		Type    string `json:"type" binding:"required"`
		Rating  *int   `json:"rating"`
		Comment *string `json:"comment"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Anime Handler: Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 验证评分范围
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 10) {
		log.Printf("Anime Handler: Invalid rating value: %d", *req.Rating)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 10"})
		return
	}
	
	collection, err := h.animeService.UpdateCollection(
		c.Request.Context(),
		id,
		req.Type,
		req.Rating,
		req.Comment,
	)
	
	if err != nil {
		log.Printf("Anime Handler: Failed to update collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}
	
	if collection == nil {
		log.Printf("Anime Handler: Collection not found with ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	
	log.Printf("Anime Handler: Successfully updated collection with ID: %d", id)
	c.JSON(http.StatusOK, collection)
}

// DeleteCollection 删除收藏
func (h *AnimeHandler) DeleteCollection(c *gin.Context) {
	log.Printf("Anime Handler: Processing delete collection request")
	
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Anime Handler: Invalid collection ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}
	
	err = h.animeService.DeleteCollection(c.Request.Context(), id)
	if err != nil {
		log.Printf("Anime Handler: Failed to delete collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}
	
	log.Printf("Anime Handler: Successfully deleted collection with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}