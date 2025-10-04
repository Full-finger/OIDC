// internal/middleware/jwt.go

package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	UserID int64    `json:"user_id"`
	Scopes []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// 检查Bearer前缀
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in format 'Bearer <token>'"})
			return
		}

		// 解析token
		secretKey := getEnv("JWT_SECRET", "default_secret_key")
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 检查claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// 将用户ID和scopes存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("scopes", claims.Scopes)
		c.Next()
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ContainsScope 检查scopes数组是否包含指定的scope
func ContainsScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// ContainsOpenIDScope 检查scopes数组是否包含openid scope
func ContainsOpenIDScope(scopes []string) bool {
	return ContainsScope(scopes, "openid")
}

// ContainsProfileScope 检查scopes数组是否包含profile scope
func ContainsProfileScope(scopes []string) bool {
	return ContainsScope(scopes, "profile") || ContainsScope(scopes, "openid") // openid通常也包含基本profile信息
}

// ContainsEmailScope 检查scopes数组是否包含email scope
func ContainsEmailScope(scopes []string) bool {
	return ContainsScope(scopes, "email")
}