package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/service"
)

// BangumiHandler Bangumi处理器
type BangumiHandler struct {
	bangumiService service.BangumiService
}

// NewBangumiHandler 创建BangumiHandler实例
func NewBangumiHandler(bangumiService service.BangumiService) *BangumiHandler {
	return &BangumiHandler{
		bangumiService: bangumiService,
	}
}

// AuthorizeHandler 处理Bangumi授权请求
func (h *BangumiHandler) AuthorizeHandler(c *gin.Context) {
	// 从上下文中获取用户ID（假设已经在中间件中设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 生成随机state参数，用于防止CSRF攻击
	state := h.generateState()
	
	// 将state与用户ID关联存储（在实际应用中应该存储在session或缓存中）
	// 这里简化处理，将用户ID编码到state中
	stateWithData := h.encodeState(state, userID.(uint))
	
	// 获取Bangumi授权URL
	authURL := h.bangumiService.GetAuthorizationURL(stateWithData)
	
	// 重定向到Bangumi授权页面
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler 处理Bangumi回调请求
func (h *BangumiHandler) CallbackHandler(c *gin.Context) {
	// 获取查询参数
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")
	
	// 检查是否有错误
	if errorParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed", "error_description": c.Query("error_description")})
		return
	}
	
	// 验证state参数（在实际应用中应该验证state是否与之前存储的一致）
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state parameter"})
		return
	}
	
	// 从state中解码用户ID（简化处理）
	userID, err := h.decodeState(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}
	
	// 用授权码换取访问令牌
	tokenResponse, err := h.bangumiService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for token", "details": err.Error()})
		return
	}
	
	// 绑定Bangumi账号
	err = h.bangumiService.BindAccount(c.Request.Context(), userID, tokenResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind bangumi account", "details": err.Error()})
		return
	}
	
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "bangumi account bound successfully",
		"user_id": userID,
	})
}

// UnbindHandler 解绑Bangumi账号
func (h *BangumiHandler) UnbindHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 解绑Bangumi账号
	err := h.bangumiService.UnbindAccount(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unbind bangumi account", "details": err.Error()})
		return
	}
	
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "bangumi account unbound successfully"})
}

// GetBoundAccountHandler 获取已绑定的Bangumi账号信息
func (h *BangumiHandler) GetBoundAccountHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 获取已绑定的Bangumi账号
	account, err := h.bangumiService.GetBoundAccount(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get bound bangumi account", "details": err.Error()})
		return
	}
	
	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no bangumi account bound"})
		return
	}
	
	// 返回账号信息（不包含敏感信息）
	response := map[string]interface{}{
		"bangumi_user_id":  account.BangumiUserID,
		"token_expires_at": account.TokenExpiresAt,
		"scope":            account.Scope,
		"created_at":       account.CreatedAt,
		"updated_at":       account.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, response)
}

// SyncCollectionHandler 同步Bangumi收藏数据
func (h *BangumiHandler) SyncCollectionHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	// 同步Bangumi收藏数据
	err := h.bangumiService.SyncCollection(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync bangumi collections", "details": err.Error()})
		return
	}
	
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "bangumi collections synced successfully"})
}

// generateState 生成随机state参数
func (h *BangumiHandler) generateState() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 出错时返回固定值，仅用于演示
		return "default_state"
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// encodeState 将用户ID编码到state中（简化处理）
func (h *BangumiHandler) encodeState(state string, userID uint) string {
	return state + "." + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", userID)))
}

// decodeState 从state中解码用户ID（简化处理）
func (h *BangumiHandler) decodeState(state string) (uint, error) {
	// 分割state字符串
	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid state format")
	}
	
	// 解码用户ID
	userIDBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed to decode user id: %w", err)
	}
	
	// 转换为uint
	userID, err := strconv.ParseUint(string(userIDBytes), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user id: %w", err)
	}
	
	return uint(userID), nil
}