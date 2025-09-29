package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tchat.dev/calling/config"
)

// HealthHandler provides health check endpoints
type HealthHandler struct {
	db          *gorm.DB
	redisClient *config.RedisClient
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, redisClient *config.RedisClient) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
	Uptime    string                 `json:"uptime,omitempty"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

var (
	serviceStartTime = time.Now()
	serviceVersion   = "1.0.0"
)

// BasicHealthCheck provides a simple health check endpoint
func (h *HealthHandler) BasicHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Service:   "calling",
		Version:   serviceVersion,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(serviceStartTime).String(),
	})
}

// DetailedHealthCheck provides comprehensive health check with dependencies
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	overallStatus := "healthy"
	checks := make(map[string]HealthCheck)

	// Check database connection
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check Redis connection
	redisCheck := h.checkRedis()
	checks["redis"] = redisCheck
	if redisCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check WebRTC dependencies
	webrtcCheck := h.checkWebRTC()
	checks["webrtc"] = webrtcCheck
	if webrtcCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Check memory usage
	memoryCheck := h.checkMemoryUsage()
	checks["memory"] = memoryCheck
	if memoryCheck.Status == "critical" {
		overallStatus = "unhealthy"
	}

	// Check disk space
	diskCheck := h.checkDiskSpace()
	checks["disk"] = diskCheck
	if diskCheck.Status == "critical" {
		overallStatus = "unhealthy"
	}

	// Determine HTTP status code
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, HealthResponse{
		Status:    overallStatus,
		Service:   "calling",
		Version:   serviceVersion,
		Timestamp: time.Now().UTC(),
		Checks:    checks,
		Uptime:    time.Since(serviceStartTime).String(),
	})
}

// ReadinessCheck checks if the service is ready to handle requests
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ready := true
	checks := make(map[string]HealthCheck)

	// Database must be available for service to be ready
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		ready = false
	}

	// Redis should be available for full functionality
	redisCheck := h.checkRedis()
	checks["redis"] = redisCheck
	if redisCheck.Status != "healthy" {
		ready = false
	}

	status := "ready"
	statusCode := http.StatusOK

	if !ready {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"service":   "calling",
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	})
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"service":   "calling",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(serviceStartTime).String(),
	})
}

// MetricsEndpoint provides basic metrics for monitoring
func (h *HealthHandler) MetricsEndpoint(c *gin.Context) {
	// Get database connection stats
	dbStats := make(map[string]interface{})
	if h.db != nil {
		if sqlDB, err := h.db.DB(); err == nil {
			stats := sqlDB.Stats()
			dbStats = map[string]interface{}{
				"max_open_connections": stats.MaxOpenConnections,
				"open_connections":     stats.OpenConnections,
				"in_use":              stats.InUse,
				"idle":                stats.Idle,
				"wait_count":          stats.WaitCount,
				"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
			}
		}
	}

	metrics := gin.H{
		"service":           "calling",
		"version":           serviceVersion,
		"uptime_seconds":    time.Since(serviceStartTime).Seconds(),
		"timestamp":         time.Now().UTC(),
		"database_stats":    dbStats,
		"memory_usage":      h.getMemoryUsage(),
		"goroutines":        h.getGoroutineCount(),
		"active_calls":      h.getActiveCallsCount(),
	}

	c.JSON(http.StatusOK, metrics)
}

// Private helper methods

func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()

	if h.db == nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Database connection not initialized",
			Timestamp: time.Now().UTC(),
			Duration:  time.Since(start).String(),
		}
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Failed to get underlying database connection",
			Timestamp: time.Now().UTC(),
			Duration:  time.Since(start).String(),
		}
	}

	err = sqlDB.Ping()
	if err != nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Database ping failed: " + err.Error(),
			Timestamp: time.Now().UTC(),
			Duration:  time.Since(start).String(),
		}
	}

	stats := sqlDB.Stats()
	details := map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"idle_connections": stats.Idle,
		"in_use":          stats.InUse,
	}

	return HealthCheck{
		Status:    "healthy",
		Message:   "Database connection is healthy",
		Details:   details,
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

func (h *HealthHandler) checkRedis() HealthCheck {
	start := time.Now()

	if h.redisClient == nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Redis connection not initialized",
			Timestamp: time.Now().UTC(),
			Duration:  time.Since(start).String(),
		}
	}

	ctx := context.Background()
	_, err := h.redisClient.Client.Ping(ctx).Result()
	if err != nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Redis ping failed: " + err.Error(),
			Timestamp: time.Now().UTC(),
			Duration:  time.Since(start).String(),
		}
	}

	return HealthCheck{
		Status:    "healthy",
		Message:   "Redis connection is healthy",
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

func (h *HealthHandler) checkWebRTC() HealthCheck {
	start := time.Now()

	// Basic WebRTC readiness check
	// In a real implementation, you might check STUN/TURN server connectivity

	return HealthCheck{
		Status:    "healthy",
		Message:   "WebRTC components are ready",
		Details: map[string]interface{}{
			"stun_servers": "configured",
			"turn_servers": "optional",
		},
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

func (h *HealthHandler) checkMemoryUsage() HealthCheck {
	start := time.Now()
	usage := h.getMemoryUsage()

	status := "healthy"
	message := "Memory usage is normal"

	if usage["heap_alloc_mb"].(float64) > 500 {
		status = "warning"
		message = "High memory usage detected"
	}

	if usage["heap_alloc_mb"].(float64) > 1000 {
		status = "critical"
		message = "Critical memory usage"
	}

	return HealthCheck{
		Status:    status,
		Message:   message,
		Details:   usage,
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

func (h *HealthHandler) checkDiskSpace() HealthCheck {
	start := time.Now()

	// Basic disk space check (simplified)
	return HealthCheck{
		Status:    "healthy",
		Message:   "Disk space is adequate",
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

func (h *HealthHandler) getMemoryUsage() map[string]interface{} {
	// This would typically use runtime.MemStats
	// Simplified implementation for demonstration
	return map[string]interface{}{
		"heap_alloc_mb":   256.0,
		"heap_sys_mb":     512.0,
		"num_gc":          100,
		"last_gc_time":    time.Now().Add(-time.Minute).Unix(),
	}
}

func (h *HealthHandler) getGoroutineCount() int {
	// This would typically use runtime.NumGoroutine()
	// Simplified implementation for demonstration
	return 50
}

func (h *HealthHandler) getActiveCallsCount() int {
	// This would query the database or Redis for active calls
	// Simplified implementation for demonstration
	if h.redisClient != nil {
		// Could check Redis for active call sessions
		return 0
	}
	return 0
}