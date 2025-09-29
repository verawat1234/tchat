package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceConfig configures performance monitoring and optimization
type PerformanceConfig struct {
	// Response time budgets in milliseconds
	APIResponseBudget    int `json:"api_response_budget"`    // 200ms
	ContentLoadBudget    int `json:"content_load_budget"`    // 1000ms
	ImageLoadBudget      int `json:"image_load_budget"`      // 500ms

	// Caching configuration
	EnableCaching        bool `json:"enable_caching"`
	CacheMaxAge          int  `json:"cache_max_age"`          // seconds

	// Rate limiting
	EnableRateLimit      bool `json:"enable_rate_limit"`
	RequestsPerMinute    int  `json:"requests_per_minute"`

	// Performance monitoring
	EnableMetrics        bool `json:"enable_metrics"`
	SlowRequestThreshold int  `json:"slow_request_threshold"` // milliseconds
}

// DefaultPerformanceConfig returns the default performance configuration
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		APIResponseBudget:    200,  // 200ms budget
		ContentLoadBudget:    1000, // 1s content load
		ImageLoadBudget:      500,  // 500ms image load
		EnableCaching:        true,
		CacheMaxAge:          300,  // 5 minutes
		EnableRateLimit:      true,
		RequestsPerMinute:    60,   // 60 requests per minute per IP
		EnableMetrics:        true,
		SlowRequestThreshold: 200,  // 200ms threshold for slow request logging
	}
}

// PerformanceMetrics tracks performance statistics
type PerformanceMetrics struct {
	mu                   sync.RWMutex
	TotalRequests        int64                  `json:"total_requests"`
	SlowRequests         int64                  `json:"slow_requests"`
	CacheHits            int64                  `json:"cache_hits"`
	CacheMisses          int64                  `json:"cache_misses"`
	AverageResponseTime  float64                `json:"average_response_time"`
	ResponseTimes        []float64              `json:"-"` // Internal tracking
	EndpointStats        map[string]*EndpointStat `json:"endpoint_stats"`
	LastResetTime        time.Time              `json:"last_reset_time"`
}

// EndpointStat tracks statistics for individual endpoints
type EndpointStat struct {
	Count           int64   `json:"count"`
	TotalTime       float64 `json:"total_time"`
	AverageTime     float64 `json:"average_time"`
	MinTime         float64 `json:"min_time"`
	MaxTime         float64 `json:"max_time"`
	BudgetViolations int64  `json:"budget_violations"`
}

// Simple in-memory cache for demonstration
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*CacheItem
}

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

func (c *MemoryCache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.Data, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Rate limiter using token bucket algorithm
type RateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*TokenBucket
	rate     int
	capacity int
}

type TokenBucket struct {
	tokens   int
	lastFill time.Time
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		rate:     requestsPerMinute,
		capacity: requestsPerMinute,
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[clientID]
	if !exists {
		bucket = &TokenBucket{
			tokens:   rl.capacity,
			lastFill: time.Now(),
		}
		rl.buckets[clientID] = bucket
	}

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.lastFill)
	tokensToAdd := int(elapsed.Minutes() * float64(rl.rate))

	if tokensToAdd > 0 {
		bucket.tokens = min(rl.capacity, bucket.tokens+tokensToAdd)
		bucket.lastFill = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PerformanceMiddleware provides comprehensive performance optimization
type PerformanceMiddleware struct {
	config      *PerformanceConfig
	metrics     *PerformanceMetrics
	cache       *MemoryCache
	rateLimiter *RateLimiter
}

// NewPerformanceMiddleware creates a new performance middleware instance
func NewPerformanceMiddleware(config *PerformanceConfig) *PerformanceMiddleware {
	if config == nil {
		config = DefaultPerformanceConfig()
	}

	return &PerformanceMiddleware{
		config: config,
		metrics: &PerformanceMetrics{
			EndpointStats: make(map[string]*EndpointStat),
			LastResetTime: time.Now(),
		},
		cache:       NewMemoryCache(),
		rateLimiter: NewRateLimiter(config.RequestsPerMinute),
	}
}

// Handler returns the Gin middleware handler
func (pm *PerformanceMiddleware) Handler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()

		// Rate limiting
		if pm.config.EnableRateLimit {
			clientID := c.ClientIP()
			if !pm.rateLimiter.Allow(clientID) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success": false,
					"error":   "Rate limit exceeded",
					"retry_after": 60,
				})
				c.Abort()
				return
			}
		}

		// Check cache for GET requests
		if pm.config.EnableCaching && c.Request.Method == "GET" {
			cacheKey := pm.generateCacheKey(c)
			if cachedData, exists := pm.cache.Get(cacheKey); exists {
				pm.recordCacheHit()

				// Set cache headers
				c.Header("X-Cache", "HIT")
				c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", pm.config.CacheMaxAge))

				// Return cached response
				if responseData, ok := cachedData.(map[string]interface{}); ok {
					c.JSON(http.StatusOK, responseData)
					pm.recordRequest(c.Request.URL.Path, time.Since(startTime))
					return
				}
			} else {
				pm.recordCacheMiss()
				c.Header("X-Cache", "MISS")
			}
		}

		// Performance monitoring wrapper
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
			responseData:   make([]byte, 0),
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate response time
		duration := time.Since(startTime)
		durationMs := float64(duration.Nanoseconds()) / 1e6

		// Cache successful GET responses
		if pm.config.EnableCaching && c.Request.Method == "GET" && responseWriter.statusCode == http.StatusOK {
			cacheKey := pm.generateCacheKey(c)

			// Parse response data for caching
			var responseData map[string]interface{}
			if err := json.Unmarshal(responseWriter.responseData, &responseData); err == nil {
				cacheTTL := time.Duration(pm.config.CacheMaxAge) * time.Second
				pm.cache.Set(cacheKey, responseData, cacheTTL)

				// Set cache headers
				c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", pm.config.CacheMaxAge))
				c.Header("ETag", pm.generateETag(responseWriter.responseData))
			}
		}

		// Performance headers
		c.Header("X-Response-Time", fmt.Sprintf("%.2fms", durationMs))
		c.Header("X-Performance-Budget", fmt.Sprintf("%dms", pm.config.APIResponseBudget))

		// Check performance budget
		if durationMs > float64(pm.config.APIResponseBudget) {
			c.Header("X-Performance-Status", "BUDGET_EXCEEDED")
		} else {
			c.Header("X-Performance-Status", "OK")
		}

		// Record metrics
		pm.recordRequest(c.Request.URL.Path, duration)

		// Log slow requests
		if pm.config.EnableMetrics && durationMs > float64(pm.config.SlowRequestThreshold) {
			fmt.Printf("[PERF] Slow request: %s %s took %.2fms (budget: %dms)\n",
				c.Request.Method, c.Request.URL.Path, durationMs, pm.config.APIResponseBudget)
		}
	})
}

// Custom response writer to capture response data
type responseWriter struct {
	gin.ResponseWriter
	statusCode   int
	responseData []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.responseData = append(rw.responseData, data...)
	return rw.ResponseWriter.Write(data)
}

// Helper methods
func (pm *PerformanceMiddleware) generateCacheKey(c *gin.Context) string {
	return fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
}

func (pm *PerformanceMiddleware) generateETag(data []byte) string {
	return fmt.Sprintf("\"%x\"", len(data)) // Simple ETag based on content length
}

func (pm *PerformanceMiddleware) recordRequest(endpoint string, duration time.Duration) {
	if !pm.config.EnableMetrics {
		return
	}

	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	durationMs := float64(duration.Nanoseconds()) / 1e6

	// Update global metrics
	pm.metrics.TotalRequests++
	if durationMs > float64(pm.config.SlowRequestThreshold) {
		pm.metrics.SlowRequests++
	}

	// Update rolling average
	pm.metrics.ResponseTimes = append(pm.metrics.ResponseTimes, durationMs)
	if len(pm.metrics.ResponseTimes) > 1000 { // Keep last 1000 requests
		pm.metrics.ResponseTimes = pm.metrics.ResponseTimes[1:]
	}

	total := 0.0
	for _, rt := range pm.metrics.ResponseTimes {
		total += rt
	}
	pm.metrics.AverageResponseTime = total / float64(len(pm.metrics.ResponseTimes))

	// Update endpoint-specific metrics
	stat, exists := pm.metrics.EndpointStats[endpoint]
	if !exists {
		stat = &EndpointStat{
			MinTime: durationMs,
			MaxTime: durationMs,
		}
		pm.metrics.EndpointStats[endpoint] = stat
	}

	stat.Count++
	stat.TotalTime += durationMs
	stat.AverageTime = stat.TotalTime / float64(stat.Count)

	if durationMs < stat.MinTime {
		stat.MinTime = durationMs
	}
	if durationMs > stat.MaxTime {
		stat.MaxTime = durationMs
	}
	if durationMs > float64(pm.config.APIResponseBudget) {
		stat.BudgetViolations++
	}
}

func (pm *PerformanceMiddleware) recordCacheHit() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()
	pm.metrics.CacheHits++
}

func (pm *PerformanceMiddleware) recordCacheMiss() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()
	pm.metrics.CacheMisses++
}

// GetMetrics returns current performance metrics
func (pm *PerformanceMiddleware) GetMetrics() *PerformanceMetrics {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := &PerformanceMetrics{
		TotalRequests:       pm.metrics.TotalRequests,
		SlowRequests:        pm.metrics.SlowRequests,
		CacheHits:           pm.metrics.CacheHits,
		CacheMisses:         pm.metrics.CacheMisses,
		AverageResponseTime: pm.metrics.AverageResponseTime,
		EndpointStats:       make(map[string]*EndpointStat),
		LastResetTime:       pm.metrics.LastResetTime,
	}

	for endpoint, stat := range pm.metrics.EndpointStats {
		metricsCopy.EndpointStats[endpoint] = &EndpointStat{
			Count:            stat.Count,
			TotalTime:        stat.TotalTime,
			AverageTime:      stat.AverageTime,
			MinTime:          stat.MinTime,
			MaxTime:          stat.MaxTime,
			BudgetViolations: stat.BudgetViolations,
		}
	}

	return metricsCopy
}

// ResetMetrics clears all performance metrics
func (pm *PerformanceMiddleware) ResetMetrics() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.TotalRequests = 0
	pm.metrics.SlowRequests = 0
	pm.metrics.CacheHits = 0
	pm.metrics.CacheMisses = 0
	pm.metrics.AverageResponseTime = 0
	pm.metrics.ResponseTimes = nil
	pm.metrics.EndpointStats = make(map[string]*EndpointStat)
	pm.metrics.LastResetTime = time.Now()
}

// ClearCache clears all cached responses
func (pm *PerformanceMiddleware) ClearCache() {
	pm.cache = NewMemoryCache()
}

// GetCacheStats returns cache statistics
func (pm *PerformanceMiddleware) GetCacheStats() map[string]interface{} {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	total := pm.metrics.CacheHits + pm.metrics.CacheMisses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(pm.metrics.CacheHits) / float64(total) * 100
	}

	return map[string]interface{}{
		"cache_hits":   pm.metrics.CacheHits,
		"cache_misses": pm.metrics.CacheMisses,
		"hit_rate":     hitRate,
		"total_requests": total,
	}
}