package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Full-finger/OIDC/internal/client"
	"github.com/Full-finger/OIDC/internal/model"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/utils"
	"github.com/gin-gonic/gin"
)

// BangumiHandler 处理Bangumi绑定相关的HTTP请求
type BangumiHandler struct {
	bangumiService service.BangumiService
}

// NewBangumiHandler 创建BangumiHandler实例
func NewBangumiHandler(bangumiService service.BangumiService) *BangumiHandler {
	return &BangumiHandler{
		bangumiService: bangumiService,
	}
}

// BindBangumiAccount 绑定Bangumi账号 - 启动绑定流程
func (h *BangumiHandler) BindBangumiAccount(c *gin.Context) {
	log.Printf("Bangumi Handler: Processing bind Bangumi account request")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Bangumi Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Bangumi Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 检查用户是否已经绑定了Bangumi账号
	isBound, err := h.bangumiService.IsUserBangumiBound(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to check if user is bound: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check binding status"})
		return
	}
	
	if isBound {
		log.Printf("Bangumi Handler: User %d already bound to Bangumi account", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already bound to Bangumi account"})
		return
	}
	
	// 获取Bangumi OAuth配置
	config := utils.GetBangumiOAuthConfig()
	
	// 生成随机state参数
	state, err := utils.GenerateState()
	if err != nil {
		log.Printf("Bangumi Handler: Failed to generate state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	
	// 将state存储在session中（这里使用Gin的session）
	// 注意：实际项目中应该使用真正的session存储机制
	c.SetCookie("bangumi_oauth_state", state, 300, "/", "", false, true) // 5分钟过期
	
	// 构建授权URL
	authURL := utils.BuildAuthorizationURL(config, state)
	
	log.Printf("Bangumi Handler: Generated authorization URL for user %d", userID)
	c.JSON(http.StatusOK, gin.H{
		"authorization_url": authURL,
		"state":            state,
	})
}

// BangumiBindCallback Bangumi OAuth回调处理
func (h *BangumiHandler) BangumiBindCallback(c *gin.Context) {
	log.Printf("Bangumi Handler: Processing Bangumi OAuth callback")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Bangumi Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Bangumi Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 获取请求参数
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" || state == "" {
		log.Printf("Bangumi Handler: Missing code or state in callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}
	
	// 验证state参数
	cookieState, err := c.Cookie("bangumi_oauth_state")
	if err != nil {
		log.Printf("Bangumi Handler: Failed to get state from cookie: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}
	
	// 删除cookie中的state
	c.SetCookie("bangumi_oauth_state", "", -1, "/", "", false, true)
	
	if state != cookieState {
		log.Printf("Bangumi Handler: State mismatch. Expected: %s, Got: %s", cookieState, state)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
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
	
	// 用授权码换取访问令牌
	tokenResp, err := bangumiClient.ExchangeCodeForToken(code)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}
	
	// 获取用户信息
	userInfo, err := bangumiClient.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	
	// 创建Bangumi绑定记录
	binding := &model.UserBangumiBinding{
		UserID:         userID,
		BangumiUserID:  userInfo.ID,
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   &tokenResp.RefreshToken,
		TokenExpiresAt: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}
	
	// 检查该Bangumi账号是否已经被其他用户绑定
	existingBinding, err := h.bangumiService.GetUserBangumiBindingByBangumiUserID(c.Request.Context(), userInfo.ID)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to check existing Bangumi binding: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing binding"})
		return
	}
	
	if existingBinding != nil && existingBinding.UserID != userID {
		log.Printf("Bangumi Handler: Bangumi user %d already bound to another user", userInfo.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bangumi account already bound to another user"})
		return
	}
	
	// 保存绑定信息
	if existingBinding != nil {
		// 更新现有绑定
		err = h.bangumiService.UpdateUserBangumiBinding(c.Request.Context(), binding)
	} else {
		// 创建新绑定
		err = h.bangumiService.CreateUserBangumiBinding(c.Request.Context(), binding)
	}
	
	if err != nil {
		log.Printf("Bangumi Handler: Failed to save Bangumi binding: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Bangumi binding"})
		return
	}
	
	log.Printf("Bangumi Handler: Successfully bound Bangumi account %d to user %d", userInfo.ID, userID)
	c.JSON(http.StatusOK, gin.H{
		"message":         "Bangumi account bound successfully",
		"bangumi_user_id": userInfo.ID,
		"username":        userInfo.Username,
		"nickname":        userInfo.Nickname,
	})
}

// SyncBangumiData 同步Bangumi数据
func (h *BangumiHandler) SyncBangumiData(c *gin.Context) {
	log.Printf("Bangumi Handler: Processing Bangumi data sync request")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Bangumi Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Bangumi Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 同步Bangumi数据
	result, err := h.bangumiService.SyncBangumiData(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to sync Bangumi data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sync Bangumi data: %v", err)})
		return
	}
	
	log.Printf("Bangumi Handler: Successfully synced Bangumi data for user %d", userID)
	c.JSON(http.StatusOK, result)
}

// GetBangumiBinding 获取用户Bangumi绑定信息
func (h *BangumiHandler) GetBangumiBinding(c *gin.Context) {
	log.Printf("Bangumi Handler: Processing get Bangumi binding request")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Bangumi Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Bangumi Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 获取用户Bangumi绑定信息
	binding, err := h.bangumiService.GetUserBangumiBindingByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to get Bangumi binding: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Bangumi binding"})
		return
	}
	
	if binding == nil {
		log.Printf("Bangumi Handler: User %d not bound to Bangumi account", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not bound to Bangumi account"})
		return
	}
	
	// 不返回敏感的令牌信息
	safeBinding := struct {
		ID              int64  `json:"id"`
		UserID          int64  `json:"user_id"`
		BangumiUserID   int64  `json:"bangumi_user_id"`
		TokenExpiresAt  string `json:"token_expires_at"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
	}{
		ID:             binding.ID,
		UserID:         binding.UserID,
		BangumiUserID:  binding.BangumiUserID,
		TokenExpiresAt: binding.TokenExpiresAt.Format("2006-01-02 15:04:05"),
		CreatedAt:      binding.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      binding.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	log.Printf("Bangumi Handler: Successfully retrieved Bangumi binding for user %d", userID)
	c.JSON(http.StatusOK, safeBinding)
}

// UnbindBangumiAccount 解绑Bangumi账号
func (h *BangumiHandler) UnbindBangumiAccount(c *gin.Context) {
	log.Printf("Bangumi Handler: Processing unbind Bangumi account request")
	
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		log.Printf("Bangumi Handler: User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("Bangumi Handler: Failed to convert user ID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	// 删除用户Bangumi绑定记录
	err := h.bangumiService.DeleteUserBangumiBinding(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Bangumi Handler: Failed to unbind Bangumi account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unbind Bangumi account"})
		return
	}
	
	log.Printf("Bangumi Handler: Successfully unbound Bangumi account for user %d", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Bangumi account unbound successfully"})
}