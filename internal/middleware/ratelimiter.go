package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"context"
	"os"
	"strconv"
)

// RateLimiter 限流器结构体
type RateLimiter struct {
	redisClient *redis.Client
	limit       int64         // 时间窗口内的请求限制次数
	window      time.Duration // 时间窗口
}

// NewRateLimiter 创建一个新的限流器
func NewRateLimiter() *RateLimiter {
	// 从环境变量获取Redis配置
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDBStr := getEnv("REDIS_DB", "0")
	redisDB, _ := strconv.Atoi(redisDBStr)

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &RateLimiter{
		redisClient: rdb,
		limit:       5,                     // 默认限制5次请求
		window:      5 * time.Minute,       // 默认时间窗口5分钟
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SetLimit 设置限流参数
func (rl *RateLimiter) SetLimit(limit int64, window time.Duration) {
	rl.limit = limit
	rl.window = window
}

// LimitByIP 基于IP的限流中间件
func (rl *RateLimiter) LimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:ip:%s", clientIP)

		// 检查是否被限流
		if rl.isRateLimited(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// LimitByUser 基于用户的限流中间件
func (rl *RateLimiter) LimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户标识（这里可以是用户ID或者邮箱等）
		// 在注册场景中，我们可以通过请求参数获取邮箱
		email := c.PostForm("email")
		if email == "" {
			// 如果没有邮箱参数，则回退到IP限流
			rl.LimitByIP()(c)
			return
		}

		key := fmt.Sprintf("rate_limit:user:%s", email)

		// 检查是否被限流
		if rl.isRateLimited(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// isRateLimited 检查是否超出限流
func (rl *RateLimiter) isRateLimited(key string) bool {
	ctx := context.Background()

	// 使用Redis的INCR命令增加计数器
	count, err := rl.redisClient.Incr(ctx, key).Result()
	if err != nil {
		// Redis出错时，允许请求通过（避免Redis故障导致服务不可用）
		return false
	}

	// 如果是第一次设置key，设置过期时间
	if count == 1 {
		rl.redisClient.Expire(ctx, key, rl.window)
	}

	// 检查是否超出限制
	return count > rl.limit
}