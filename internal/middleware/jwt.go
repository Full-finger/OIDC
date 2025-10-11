package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Full-finger/OIDC/internal/util"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取访问令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		
		// 解析Bearer令牌
		tokenString := ""
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}
		
		// 解析访问令牌
		jwtUtil, err := util.NewJWTUtil()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize JWT utility"})
			c.Abort()
			return
		}
		
		claims, err := jwtUtil.ParseAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}
		
		// 从声明中提取用户ID (通过Subject字段)
		if claims.Subject == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token: missing subject"})
			c.Abort()
			return
		}
		
		// 将用户ID存储到上下文中
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}