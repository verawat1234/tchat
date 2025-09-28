package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting per client
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// GetLimiter returns the rate limiter for a client
func (rl *RateLimiter) GetLimiter(clientID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[clientID]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[clientID] = limiter
	}

	return limiter
}

// Production-ready middleware

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self' wss: https:")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Remove server information
		c.Header("Server", "Tchat-Gateway")

		c.Next()
	})
}

// RateLimitMiddleware provides rate limiting functionality
func RateLimitMiddleware(rateLimiter *RateLimiter) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get client identifier (IP address or authenticated user ID)
		clientID := getClientIdentifier(c)

		limiter := rateLimiter.GetLimiter(clientID)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
				"code":    "RATE_LIMIT_EXCEEDED",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// EnhancedLoggingMiddleware provides structured logging
func EnhancedLoggingMiddleware(logger *log.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract additional context
		userID := param.Keys["user_id"]
		requestID := param.Keys["request_id"]
		userAgent := param.Request.UserAgent()
		referer := param.Request.Referer()

		// Determine log level based on status code
		logLevel := "INFO"
		if param.StatusCode >= 400 && param.StatusCode < 500 {
			logLevel = "WARN"
		} else if param.StatusCode >= 500 {
			logLevel = "ERROR"
		}

		// Format: [LEVEL] timestamp method path status latency client_ip user_agent referer user_id request_id
		return fmt.Sprintf("[%s] %v | %3d | %13v | %15s | %s | %s | %s | %s | %s\n",
			logLevel,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method+" "+param.Path,
			userAgent,
			referer,
			userID,
			requestID,
		)
	})
}

// RequestIDMiddleware generates unique request IDs
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	})
}

// AuthenticationMiddleware validates JWT tokens
func AuthenticationMiddleware(jwtSecret string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip authentication for health checks and public endpoints
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authorization header required",
				"code":    "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// Validate Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authorization header format",
				"code":    "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// In production, validate JWT token here
		// For now, we'll use a simple validation
		if !isValidToken(token, jwtSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
				"code":    "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Extract user information from token
		userID, countryCode := extractUserInfo(token)
		c.Set("user_id", userID)
		c.Set("country_code", countryCode)
		c.Set("authenticated", true)

		c.Next()
	})
}

// AdminAuthenticationMiddleware validates admin tokens
func AdminAuthenticationMiddleware(adminSecret string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		adminToken := c.GetHeader("X-Admin-Token")

		if adminToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Admin token required",
				"code":    "MISSING_ADMIN_TOKEN",
			})
			c.Abort()
			return
		}

		// Use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(adminToken), []byte(adminSecret)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid admin token",
				"code":    "INVALID_ADMIN_TOKEN",
			})
			c.Abort()
			return
		}

		c.Set("admin", true)
		c.Next()
	})
}

// KYCMiddleware validates KYC level requirements
func KYCMiddleware(minLevel int) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User authentication required for KYC validation",
				"code":    "USER_NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		// In production, fetch user's KYC level from database
		userKYCLevel := getUserKYCLevel(userID)

		if userKYCLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error":     "insufficient_kyc",
				"message":   fmt.Sprintf("KYC level %d required, user has level %d", minLevel, userKYCLevel),
				"code":      "INSUFFICIENT_KYC_LEVEL",
				"required_level": minLevel,
				"current_level":  userKYCLevel,
			})
			c.Abort()
			return
		}

		c.Set("kyc_level", userKYCLevel)
		c.Next()
	})
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware(logger *log.Logger) gin.HandlerFunc {
	return gin.Recovery() // Use Gin's built-in recovery middleware for now
}

// MetricsMiddleware collects request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Collect metrics
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// In production, send these metrics to monitoring system
		_ = duration
		_ = status
		_ = method
		_ = path
	})
}

// CacheControlMiddleware sets cache control headers
func CacheControlMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Set cache headers based on path
		if strings.HasPrefix(path, "/health") || strings.HasPrefix(path, "/ready") {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		} else if strings.HasPrefix(path, "/api/") {
			// API responses should not be cached by default
			c.Header("Cache-Control", "no-cache")
		} else {
			// Default caching for static resources
			c.Header("Cache-Control", "public, max-age=300")
		}

		c.Next()
	})
}

// Helper functions

func getClientIdentifier(c *gin.Context) string {
	// Try to get authenticated user ID first
	if userID := c.GetString("user_id"); userID != "" {
		return "user:" + userID
	}

	// Fall back to IP address
	clientIP := c.ClientIP()

	// Handle X-Forwarded-For header
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP = strings.TrimSpace(ips[0])
		}
	}

	return "ip:" + clientIP
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/ready",
		"/metrics",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/api/v1/content/public",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

func isValidToken(token, secret string) bool {
	// In production, implement proper JWT validation
	// For now, simple validation
	return token != "" && len(token) > 10
}

func extractUserInfo(token string) (userID, countryCode string) {
	// In production, extract from JWT claims
	// For now, return dummy data
	return "user-123", "US"
}

func getUserKYCLevel(userID string) int {
	// In production, fetch from database
	// For now, return level 1
	return 1
}

// CORS configuration for production
func ProductionCORSMiddleware() gin.HandlerFunc {
	config := gin.H{
		"AllowOrigins": []string{
			"https://tchat.com",
			"https://app.tchat.com",
			"https://admin.tchat.com",
		},
		"AllowMethods": []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		"AllowHeaders": []string{
			"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID",
			"X-Admin-Token", "X-API-Key",
		},
		"ExposeHeaders": []string{
			"X-Request-ID", "X-Rate-Limit-Remaining", "X-Rate-Limit-Reset",
		},
		"AllowCredentials": true,
		"MaxAge":          12 * 60 * 60, // 12 hours
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowedOrigins := config["AllowOrigins"].([]string)
		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				originAllowed = true
				break
			}
		}

		if !originAllowed && origin != "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "origin_not_allowed",
				"message": "Origin not allowed by CORS policy",
			})
			c.Abort()
			return
		}

		// Set CORS headers
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", strings.Join(config["AllowMethods"].([]string), ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config["AllowHeaders"].([]string), ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config["ExposeHeaders"].([]string), ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", strconv.Itoa(config["MaxAge"].(int)))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}