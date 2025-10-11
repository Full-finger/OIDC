package handler

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/model"
)

// AnimeHandler 番剧处理器
type AnimeHandler struct {
	animeService service.AnimeService
}

// NewAnimeHandler 创建AnimeHandler实例
func NewAnimeHandler(animeService service.AnimeService) *AnimeHandler {
	return &AnimeHandler{
		animeService: animeService,
	}
}

// GetAnimeByIDHandler 根据ID获取番剧
func (h *AnimeHandler) GetAnimeByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	anime, err := h.animeService.GetAnimeByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get anime"})
		return
	}
	
	if anime == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "anime not found"})
		return
	}
	
	c.JSON(http.StatusOK, anime)
}

// SearchAnimesHandler 搜索番剧
func (h *AnimeHandler) SearchAnimesHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	
	animes, err := h.animeService.SearchAnimes(c.Request.Context(), keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search animes"})
		return
	}
	
	c.JSON(http.StatusOK, animes)
}

// ListAnimesHandler 列出所有番剧
func (h *AnimeHandler) ListAnimesHandler(c *gin.Context) {
	animes, err := h.animeService.ListAnimes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list animes"})
		return
	}
	
	c.JSON(http.StatusOK, animes)
}

// ListAnimesByStatusHandler 根据状态列出番剧
func (h *AnimeHandler) ListAnimesByStatusHandler(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	
	animes, err := h.animeService.ListAnimesByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list animes by status"})
		return
	}
	
	c.JSON(http.StatusOK, animes)
}

// CreateAnimeHandler 创建番剧
func (h *AnimeHandler) CreateAnimeHandler(c *gin.Context) {
	var anime model.Anime
	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.animeService.CreateAnime(c.Request.Context(), &anime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create anime"})
		return
	}
	
	c.JSON(http.StatusOK, anime)
}

// UpdateAnimeHandler 更新番剧
func (h *AnimeHandler) UpdateAnimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	var anime model.Anime
	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 设置番剧ID
	anime.ID = uint(id)
	
	if err := h.animeService.UpdateAnime(c.Request.Context(), &anime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update anime"})
		return
	}
	
	c.JSON(http.StatusOK, anime)
}

// DeleteAnimeHandler 删除番剧
func (h *AnimeHandler) DeleteAnimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime id"})
		return
	}
	
	if err := h.animeService.DeleteAnime(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete anime"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "anime deleted successfully"})
}
