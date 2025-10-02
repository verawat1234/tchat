package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string]*userLimit
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type userLimit struct {
	count     int
	resetTime time.Time
}

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	limiter := &rateLimiter{
		requests: make(map[string]*userLimit),
		limit:    limit,
		window:   window,
	}

	// Cleanup goroutine
	go limiter.cleanup()

	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		key := userID.(string)

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		userLim, exists := limiter.requests[key]

		if !exists || now.After(userLim.resetTime) {
			limiter.requests[key] = &userLimit{
				count:     1,
				resetTime: now.Add(window),
			}
			c.Next()
			return
		}

		if userLim.count >= limit {
			c.Header("Retry-After", userLim.resetTime.Sub(now).String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": userLim.resetTime.Unix(),
			})
			c.Abort()
			return
		}

		userLim.count++
		c.Next()
	}
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, limit := range rl.requests {
			if now.After(limit.resetTime) {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}