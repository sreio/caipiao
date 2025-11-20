package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器管理
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

// getLimiter 获取或创建IP的限流器
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware 限流中间件
// rate: 每秒允许的请求数
// burst: 令牌桶容量
func RateLimitMiddleware(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(r, burst)

	// 定期清理过期的限流器（每10分钟）
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.mu.Lock()
			// 清理所有限流器（简单策略，实际可以记录最后访问时间）
			limiter.limiters = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()

		// 获取该IP的限流器
		ipLimiter := limiter.getLimiter(ip)

		// 检查是否允许请求
		if !ipLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DefaultRateLimitMiddleware 默认限流中间件
// 默认配置：每秒2个请求，突发容量5
func DefaultRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(rate.Limit(2), 5)
}

// StrictRateLimitMiddleware 严格限流中间件
// 用于敏感接口：每秒1个请求，突发容量2
func StrictRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(rate.Limit(1), 2)
}
