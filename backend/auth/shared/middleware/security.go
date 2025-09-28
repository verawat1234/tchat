package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"tchat.dev/auth/services"
	"tchat.dev/shared/responses"
)

// SecurityConfig holds comprehensive security middleware configuration
type SecurityConfig struct {
	// JWT Configuration
	JWTService      *services.JWTService
	SessionService  *services.SessionService
	UserService     *services.UserService

	// Rate Limiting
	RateLimiter     *AdvancedRateLimiter

	// Security Headers
	EnableSecurityHeaders bool
	EnableCORS           bool

	// Input Validation
	EnableInputValidation bool
	MaxRequestSize       int64

	// Logging
	Logger              *logrus.Logger

	// Redis for distributed features
	Redis               *redis.Client

	// CORS Configuration
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	AllowCredentials    bool
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig(jwtService *services.JWTService, redis *redis.Client, logger *logrus.Logger) *SecurityConfig {
	return &SecurityConfig{
		JWTService:            jwtService,
		RateLimiter:           NewAdvancedRateLimiter(redis, logger),
		EnableSecurityHeaders: true,
		EnableCORS:           true,
		EnableInputValidation: true,
		MaxRequestSize:       10 * 1024 * 1024, // 10MB
		Logger:               logger,
		Redis:                redis,
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://tchat.co.th",
			"https://tchat.com.sg",
			"https://tchat.co.id",
			"https://tchat.com.my",
			"https://tchat.com.ph",
			"https://tchat.com.vn",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"X-Request-ID",
			"X-Country-Code",
			"X-Locale",
			"X-Device-ID",
		},
		AllowCredentials: true,
	}
}

// ComprehensiveSecurityMiddleware creates a comprehensive security middleware stack
func ComprehensiveSecurityMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Apply security headers first
		if config.EnableSecurityHeaders {
			applySecurityHeaders(c)
		}

		// Handle CORS
		if config.EnableCORS {
			if !handleCORS(c, config) {
				return // CORS rejected the request
			}
		}

		// Rate limiting check
		if config.RateLimiter != nil {
			if !config.RateLimiter.Allow(c) {
				return // Rate limit exceeded
			}
		}

		// Request size validation
		if config.EnableInputValidation && config.MaxRequestSize > 0 {
			if !validateRequestSize(c, config.MaxRequestSize) {
				return // Request too large
			}
		}

		c.Next()
	})
}

// AuthenticationMiddleware provides JWT-based authentication with revocation checking
func AuthenticationMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip authentication for public endpoints
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Extract and validate Bearer token
		token := extractBearerToken(c)
		if token == "" {
			responses.UnauthorizedResponse(c, "Authorization token required")
			c.Abort()
			return
		}

		// Validate JWT token with revocation checking
		claims, err := config.JWTService.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			if config.Logger != nil {
				config.Logger.WithFields(logrus.Fields{
					"error":      err.Error(),
					"ip":         c.ClientIP(),
					"user_agent": c.GetHeader("User-Agent"),
					"path":       c.Request.URL.Path,
				}).Warn("JWT authentication failed")
			}

			responses.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Additional session validation if configured
		if config.SessionService != nil {
			session, err := config.SessionService.ValidateSessionActive(c.Request.Context(), claims.SessionID)
			if err != nil {
				responses.UnauthorizedResponse(c, "Invalid session")
				c.Abort()
				return
			}

			// Update session activity
			config.SessionService.UpdateLastActive(c.Request.Context(), session.ID)
		}

		// Set authentication context
		setAuthenticationContext(c, claims)

		// Log successful authentication for audit trail
		if config.Logger != nil {
			config.Logger.WithFields(logrus.Fields{
				"user_id":      claims.UserID,
				"session_id":   claims.SessionID,
				"ip":           c.ClientIP(),
				"user_agent":   c.GetHeader("User-Agent"),
				"path":         c.Request.URL.Path,
				"method":       c.Request.Method,
			}).Info("User authenticated successfully")
		}

		c.Next()
	})
}

// InputValidationMiddleware provides comprehensive input validation and sanitization
func InputValidationMiddleware(config *SecurityConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if !config.EnableInputValidation {
			c.Next()
			return
		}

		// Validate content type for POST/PUT/PATCH requests
		if isWriteMethod(c.Request.Method) {
			contentType := c.GetHeader("Content-Type")
			if !isValidContentType(contentType) {
				responses.BadRequestResponse(c, "Invalid content type")
				c.Abort()
				return
			}
		}

		// SQL injection protection for query parameters
		if containsSQLInjection(c.Request.URL.RawQuery) {
			if config.Logger != nil {
				config.Logger.WithFields(logrus.Fields{
					"ip":           c.ClientIP(),
					"user_agent":   c.GetHeader("User-Agent"),
					"path":         c.Request.URL.Path,
					"query":        c.Request.URL.RawQuery,
				}).Warn("SQL injection attempt detected")
			}

			responses.BadRequestResponse(c, "Invalid query parameters")
			c.Abort()
			return
		}

		// XSS protection for headers
		if containsXSSPatterns(c) {
			if config.Logger != nil {
				config.Logger.WithFields(logrus.Fields{
					"ip":         c.ClientIP(),
					"user_agent": c.GetHeader("User-Agent"),
					"path":       c.Request.URL.Path,
				}).Warn("XSS attempt detected in headers")
			}

			responses.BadRequestResponse(c, "Invalid request headers")
			c.Abort()
			return
		}

		c.Next()
	})
}

// AdvancedRateLimiter provides distributed rate limiting with Redis
type AdvancedRateLimiter struct {
	redis       *redis.Client
	logger      *logrus.Logger
	limiters    map[string]*rate.Limiter
	mu          sync.Mutex

	// Default rate limits
	requestsPerMinute int
	burstSize        int

	// Endpoint-specific limits
	endpointLimits   map[string]EndpointLimit
}

type EndpointLimit struct {
	RequestsPerMinute int
	BurstSize        int
	WindowDuration   time.Duration
}

// NewAdvancedRateLimiter creates a new advanced rate limiter
func NewAdvancedRateLimiter(redis *redis.Client, logger *logrus.Logger) *AdvancedRateLimiter {
	limiter := &AdvancedRateLimiter{
		redis:            redis,
		logger:           logger,
		limiters:         make(map[string]*rate.Limiter),
		requestsPerMinute: 60,
		burstSize:        10,
		endpointLimits:   make(map[string]EndpointLimit),
	}

	// Configure endpoint-specific limits
	limiter.endpointLimits["/auth/login"] = EndpointLimit{
		RequestsPerMinute: 5,
		BurstSize:        3,
		WindowDuration:   time.Minute,
	}
	limiter.endpointLimits["/auth/register"] = EndpointLimit{
		RequestsPerMinute: 3,
		BurstSize:        2,
		WindowDuration:   time.Minute,
	}
	limiter.endpointLimits["/auth/refresh"] = EndpointLimit{
		RequestsPerMinute: 10,
		BurstSize:        5,
		WindowDuration:   time.Minute,
	}

	return limiter
}

// Allow checks if the request should be allowed based on rate limits
func (rl *AdvancedRateLimiter) Allow(c *gin.Context) bool {
	clientID := getClientIdentifier(c)
	endpoint := c.Request.URL.Path

	// Use Redis for distributed rate limiting if available
	if rl.redis != nil {
		return rl.allowDistributed(c, clientID, endpoint)
	}

	// Fallback to in-memory rate limiting
	return rl.allowInMemory(c, clientID, endpoint)
}

// allowDistributed implements distributed rate limiting using Redis
func (rl *AdvancedRateLimiter) allowDistributed(c *gin.Context, clientID, endpoint string) bool {
	limit := rl.getEndpointLimit(endpoint)
	key := fmt.Sprintf("rate_limit:%s:%s", clientID, endpoint)

	ctx := c.Request.Context()

	// Get current count
	count, err := rl.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		// Redis error, fallback to allow (log the error)
		if rl.logger != nil {
			rl.logger.WithError(err).Warn("Redis rate limit check failed")
		}
		return true
	}

	// Check if limit exceeded
	if count >= limit.RequestsPerMinute {
		rl.sendRateLimitResponse(c, limit.WindowDuration)
		return false
	}

	// Increment counter
	pipe := rl.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, limit.WindowDuration)
	_, err = pipe.Exec(ctx)

	if err != nil && rl.logger != nil {
		rl.logger.WithError(err).Warn("Redis rate limit increment failed")
	}

	// Set rate limit headers
	c.Header("X-RateLimit-Limit", strconv.Itoa(limit.RequestsPerMinute))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(limit.RequestsPerMinute-count-1))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(limit.WindowDuration).Unix(), 10))

	return true
}

// allowInMemory implements in-memory rate limiting as fallback
func (rl *AdvancedRateLimiter) allowInMemory(c *gin.Context, clientID, endpoint string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := fmt.Sprintf("%s:%s", clientID, endpoint)
	limiter, exists := rl.limiters[key]
	if !exists {
		limit := rl.getEndpointLimit(endpoint)
		limiter = rate.NewLimiter(rate.Limit(limit.RequestsPerMinute)/60, limit.BurstSize)
		rl.limiters[key] = limiter
	}

	if !limiter.Allow() {
		limit := rl.getEndpointLimit(endpoint)
		rl.sendRateLimitResponse(c, limit.WindowDuration)
		return false
	}

	return true
}

// getEndpointLimit returns rate limit configuration for specific endpoint
func (rl *AdvancedRateLimiter) getEndpointLimit(endpoint string) EndpointLimit {
	if limit, exists := rl.endpointLimits[endpoint]; exists {
		return limit
	}

	// Default limit
	return EndpointLimit{
		RequestsPerMinute: rl.requestsPerMinute,
		BurstSize:        rl.burstSize,
		WindowDuration:   time.Minute,
	}
}

// sendRateLimitResponse sends a rate limit exceeded response
func (rl *AdvancedRateLimiter) sendRateLimitResponse(c *gin.Context, retryAfter time.Duration) {
	c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))

	c.JSON(429, gin.H{"error": "too_many_requests", "message": "Rate limit exceeded"})
	c.Abort()

	if rl.logger != nil {
		rl.logger.WithFields(logrus.Fields{
			"client_ip": c.ClientIP(),
			"endpoint":  c.Request.URL.Path,
			"method":    c.Request.Method,
		}).Warn("Rate limit exceeded")
	}
}

// Helper functions

func applySecurityHeaders(c *gin.Context) {
	// Prevent MIME type sniffing
	c.Header("X-Content-Type-Options", "nosniff")

	// Prevent clickjacking
	c.Header("X-Frame-Options", "DENY")

	// XSS protection
	c.Header("X-XSS-Protection", "1; mode=block")

	// Force HTTPS
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	// Content Security Policy
	csp := "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: https:; " +
		"connect-src 'self' wss: https:; " +
		"font-src 'self' data:; " +
		"object-src 'none'; " +
		"media-src 'self'; " +
		"frame-src 'none'"
	c.Header("Content-Security-Policy", csp)

	// Referrer policy
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	// Permissions policy
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=()")

	// Remove server information
	c.Header("Server", "")
}

func handleCORS(c *gin.Context, config *SecurityConfig) bool {
	origin := c.GetHeader("Origin")

	// Check if origin is allowed
	allowed := false
	for _, allowedOrigin := range config.AllowedOrigins {
		if allowedOrigin == "*" || origin == allowedOrigin {
			allowed = true
			break
		}

		// Support wildcard subdomains
		if strings.Contains(allowedOrigin, "*") {
			pattern := strings.Replace(allowedOrigin, "*", "", -1)
			if strings.Contains(origin, pattern) {
				allowed = true
				break
			}
		}
	}

	if !allowed && origin != "" {
		responses.ForbiddenResponse(c, "Origin not allowed by CORS policy")
		c.Abort()
		return false
	}

	// Set CORS headers
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	if config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
	c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
	c.Header("Access-Control-Max-Age", "86400")

	// Handle preflight requests
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return false
	}

	return true
}

func validateRequestSize(c *gin.Context, maxSize int64) bool {
	if c.Request.ContentLength > maxSize {
		responses.BadRequestResponse(c, "Request entity too large")
		c.Abort()
		return false
	}
	return true
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func setAuthenticationContext(c *gin.Context, claims *services.UserClaims) {
	c.Set("user_id", claims.UserID)
	c.Set("session_id", claims.SessionID)
	c.Set("device_id", claims.DeviceID)
	c.Set("permissions", claims.Permissions)
	c.Set("scopes", claims.Scopes)
	c.Set("kyc_level", claims.KYCLevel)
	c.Set("country_code", claims.CountryCode)
	c.Set("user_claims", claims)
}

func getClientIdentifier(c *gin.Context) string {
	// Try authenticated user first
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	// Fall back to IP address
	clientIP := c.ClientIP()
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP = strings.TrimSpace(ips[0])
		}
	}

	return fmt.Sprintf("ip:%s", clientIP)
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
		"/docs",
		"/swagger",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

func isWriteMethod(method string) bool {
	return method == http.MethodPost ||
		   method == http.MethodPut ||
		   method == http.MethodPatch
}

func isValidContentType(contentType string) bool {
	validTypes := []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}

	for _, validType := range validTypes {
		if strings.Contains(contentType, validType) {
			return true
		}
	}

	return false
}

// Basic SQL injection detection patterns
func containsSQLInjection(query string) bool {
	suspiciousPatterns := []string{
		"'",
		"\"",
		";",
		"--",
		"/*",
		"*/",
		"xp_",
		"sp_",
		"union",
		"select",
		"insert",
		"update",
		"delete",
		"drop",
		"exec",
		"execute",
	}

	lowerQuery := strings.ToLower(query)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerQuery, pattern) {
			return true
		}
	}

	return false
}

// Basic XSS detection patterns
func containsXSSPatterns(c *gin.Context) bool {
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"onload=",
		"onerror=",
		"onclick=",
		"<iframe",
		"<object",
		"<embed",
	}

	// Check various headers for XSS patterns
	headersToCheck := []string{
		"User-Agent",
		"Referer",
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, header := range headersToCheck {
		value := strings.ToLower(c.GetHeader(header))
		for _, pattern := range suspiciousPatterns {
			if strings.Contains(value, pattern) {
				return true
			}
		}
	}

	return false
}