package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/auth/services"
	"tchat.dev/shared/responses"
)

// AuthMiddleware handles authentication middleware
type AuthMiddleware struct {
	jwtService     *services.JWTService
	sessionService *services.SessionService
	userService    *services.UserService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(jwtService *services.JWTService, sessionService *services.SessionService, userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:     jwtService,
		sessionService: sessionService,
		userService:    userService,
	}
}

// RequireAuth middleware that requires authentication
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responses.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			responses.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := am.jwtService.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			responses.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Validate session
		session, err := am.sessionService.ValidateSessionActive(c.Request.Context(), claims.SessionID)
		if err != nil {
			responses.UnauthorizedResponse(c, "Invalid session")
			c.Abort()
			return
		}

		// Update session activity
		am.sessionService.UpdateLastActive(c.Request.Context(), session.ID)

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("session_id", claims.SessionID)
		c.Set("device_id", claims.DeviceID)
		c.Set("permissions", claims.Permissions)
		c.Set("scopes", claims.Scopes)

		c.Next()
	}
}

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

// RateLimiter represents a rate limiter
type RateLimiter struct {
	requestsPerMinute int
	burstSize         int
	cleanupInterval   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute, burstSize int, cleanupInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		requestsPerMinute: requestsPerMinute,
		burstSize:         burstSize,
		cleanupInterval:   cleanupInterval,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple rate limiting implementation
		// In production, you would use Redis or a proper rate limiting library
		clientIP := c.ClientIP()

		// For now, just log the rate limiting attempt
		LogInfo(c, "Rate limit check", gin.H{
			"client_ip": clientIP,
			"path":      c.Request.URL.Path,
		})

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
	println("LOG:", message)
}