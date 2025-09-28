package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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