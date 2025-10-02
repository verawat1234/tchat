package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}


// LogInfo logs an info message with request context
func LogInfo(c *gin.Context, message string, data gin.H) {
	// Simple logging implementation
	// In production, you would use a structured logger like logrus or zap
	userID, _ := c.Get("user_id")
	sessionID, _ := c.Get("session_id")

	logData := gin.H{
		"level":     "info",
		"message":   message,
		"timestamp": time.Now().UTC(),
		"request_id": c.GetHeader("X-Request-ID"),
		"user_id":   userID,
		"session_id": sessionID,
		"ip":        c.ClientIP(),
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
	}

	for k, v := range data {
		logData[k] = v
	}

	// In production, this would go to a proper logging system
	println("LOG INFO:", message)
}

// LogWarning logs a warning message with request context
func LogWarning(c *gin.Context, message string, data gin.H) {
	// Simple logging implementation
	// In production, you would use a structured logger like logrus or zap
	userID, _ := c.Get("user_id")
	sessionID, _ := c.Get("session_id")

	logData := gin.H{
		"level":     "warning",
		"message":   message,
		"timestamp": time.Now().UTC(),
		"request_id": c.GetHeader("X-Request-ID"),
		"user_id":   userID,
		"session_id": sessionID,
		"ip":        c.ClientIP(),
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
	}

	for k, v := range data {
		logData[k] = v
	}

	// In production, this would go to a proper logging system
	println("LOG WARNING:", message)
}

// Simple middleware wrappers for backward compatibility with handlers

// SimpleCORSMiddleware returns a simple CORS middleware without configuration
func SimpleCORSMiddleware() gin.HandlerFunc {
	return CORSMiddleware(DefaultCORSConfig())
}

// SimpleAuthMiddleware returns a simple auth middleware without configuration
func SimpleAuthMiddleware() gin.HandlerFunc {
	// Create a simple auth middleware that validates Bearer tokens
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		// Simple placeholder - in production this would validate JWT
		c.Set("user_id", "placeholder-user-id")
		c.Set("authenticated", true)
		c.Next()
	}
}

// SimpleRequestLogger returns a simple request logger
func SimpleRequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.Info("Request processed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
		)
	}
}

// SimpleAdminOnly returns a simple admin-only middleware
func SimpleAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, exists := c.Get("authenticated")
		if !exists || !authenticated.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Simple placeholder - in production this would check role from JWT
		isAdmin := c.GetHeader("X-Admin") == "true"
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SimpleRateLimit returns a simple rate limiting middleware
func SimpleRateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := &simpleLimiter{
		requests: make(map[string][]time.Time),
	}

	return func(c *gin.Context) {
		clientID := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			clientID = userID.(string)
		}

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		windowStart := now.Add(-window)

		requests := limiter.requests[clientID]
		validRequests := []time.Time{}
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) >= maxRequests {
			c.Header("X-Rate-Limit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		validRequests = append(validRequests, now)
		limiter.requests[clientID] = validRequests

		c.Next()
	}
}

type simpleLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
}