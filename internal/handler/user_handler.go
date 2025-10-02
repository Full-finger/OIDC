// internal/handler/user_handler.go

package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler 负责处理用户相关的 HTTP 请求
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRequest 定义注册请求的 JSON 结构
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest 定义登录请求的 JSON 结构
type LoginRequest struct {
	Username string `json:"username"` // 可以是用户名或邮箱
	Password string `json:"password"`
}

// RegisterHandler 处理用户注册请求
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析 JSON 请求体
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. 参数校验
	if req.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Username) < 3 || len(req.Username) > 20 {
		http.Error(w, "username must be between 3 and 20 characters", http.StatusBadRequest)
		return
	}
	
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	
	// 简单的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	
	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// 3. 调用服务层注册用户
	safeUser, err := h.userService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		// 根据错误类型返回不同状态码（简化处理）
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// 4. 返回成功响应（JSON 格式）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": safeUser,
	})
}

// LoginHandler 处理用户登录请求
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析 JSON 请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. 参数校验
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	// 3. 调用服务层登录
	safeUser, err := h.userService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	// 4. 生成JWT token
	tokenString, err := generateJWT(safeUser.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// 5. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  safeUser,
		"token": tokenString,
	})
}

// UpdateProfileRequest 定义更新资料请求
type UpdateProfileRequest struct {
	Nickname  *string `json:"nickname,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}

// UpdateProfileHandler 处理更新用户资料请求（需认证）
func (h *UserHandler) UpdateProfileHandler(c *gin.Context) {
	// 从Gin上下文中获取用户ID（由JWT中间件设置）
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	userID, ok := userIDValue.(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	safeUser, err := h.userService.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL, req.Bio)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": safeUser})
}

// GetProfileHandler 获取当前用户资料
func (h *UserHandler) GetProfileHandler(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.SafeUser()})
}

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// generateJWT 生成JWT token
func generateJWT(userID int64) (string, error) {
	// 从环境变量获取密钥，如果没有则使用默认值
	secretKey := getEnv("JWT_SECRET", "default_secret_key")
	
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}