package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tchat.dev/shared/config"
	"tchat.dev/shared/responses"
)

// DefaultHealthChecker provides standard health checking functionality
type DefaultHealthChecker struct {
	serviceName    string
	serviceVersion string
	db             *gorm.DB
	config         *config.Config
	registry       ServiceRegistry
	customChecks   map[string]func() error
}

// NewDefaultHealthChecker creates a new health checker
func NewDefaultHealthChecker(serviceName, serviceVersion string, db *gorm.DB, cfg *config.Config, registry ServiceRegistry) HealthChecker {
	return &DefaultHealthChecker{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		db:             db,
		config:         cfg,
		registry:       registry,
		customChecks:   make(map[string]func() error),
	}
}

// AddCustomCheck adds a custom health check
func (h *DefaultHealthChecker) AddCustomCheck(name string, checkFunc func() error) {
	h.customChecks[name] = checkFunc
}

// HealthCheck performs a basic health check
func (h *DefaultHealthChecker) HealthCheck(c *gin.Context) {
	healthData := h.GetHealthData()

	responses.SendSuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   h.serviceName,
		"version":   h.serviceVersion,
		"timestamp": time.Now().UTC(),
		"data":      healthData,
	})
}

// ReadinessCheck performs a comprehensive readiness check
func (h *DefaultHealthChecker) ReadinessCheck(c *gin.Context) {
	checks := make(map[string]interface{})
	allHealthy := true

	// Check database connection
	if h.db != nil {
		if err := h.checkDatabase(); err != nil {
			checks["database"] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			allHealthy = false
		} else {
			checks["database"] = map[string]interface{}{
				"status": "healthy",
			}
		}
	}

	// Check registered service components
	if h.registry != nil {
		componentChecks := h.checkServiceComponents()
		checks["components"] = componentChecks

		// Check if any component is unhealthy
		for _, component := range h.registry.List() {
			if !component.IsHealthy() {
				allHealthy = false
				break
			}
		}
	}

	// Run custom health checks
	customChecks := h.runCustomChecks()
	if len(customChecks) > 0 {
		checks["custom"] = customChecks

		// Check if any custom check failed
		for _, check := range customChecks {
			if checkMap, ok := check.(map[string]interface{}); ok {
				if status, exists := checkMap["status"]; exists && status != "healthy" {
					allHealthy = false
					break
				}
			}
		}
	}

	status := "ready"
	httpStatus := http.StatusOK

	if !allHealthy {
		status = "not ready"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"service":   h.serviceName,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	})
}

// GetHealthData returns basic health information
func (h *DefaultHealthChecker) GetHealthData() map[string]interface{} {
	data := map[string]interface{}{
		"service":     h.serviceName,
		"version":     h.serviceVersion,
		"environment": h.config.Environment,
		"debug":       h.config.Debug,
		"uptime":      time.Now().UTC(),
	}

	// Add component count if registry is available
	if h.registry != nil {
		components := h.registry.List()
		data["components_count"] = len(components)
	}

	return data
}

// checkDatabase checks database connectivity
func (h *DefaultHealthChecker) checkDatabase() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// checkServiceComponents checks all registered service components
func (h *DefaultHealthChecker) checkServiceComponents() map[string]interface{} {
	componentChecks := make(map[string]interface{})

	if h.registry == nil {
		return componentChecks
	}

	for name, component := range h.registry.List() {
		status := "healthy"
		if !component.IsHealthy() {
			status = "unhealthy"
		}

		componentChecks[name] = map[string]interface{}{
			"status": status,
		}
	}

	return componentChecks
}

// runCustomChecks runs all custom health checks
func (h *DefaultHealthChecker) runCustomChecks() map[string]interface{} {
	customChecks := make(map[string]interface{})

	for name, checkFunc := range h.customChecks {
		if err := checkFunc(); err != nil {
			customChecks[name] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			customChecks[name] = map[string]interface{}{
				"status": "healthy",
			}
		}
	}

	return customChecks
}

// ServiceHealthChecker provides service-specific health checking
type ServiceHealthChecker struct {
	*DefaultHealthChecker
	additionalInfo map[string]interface{}
}

// NewServiceHealthChecker creates a service-specific health checker
func NewServiceHealthChecker(serviceName, serviceVersion string, db *gorm.DB, cfg *config.Config, registry ServiceRegistry) HealthChecker {
	return &ServiceHealthChecker{
		DefaultHealthChecker: NewDefaultHealthChecker(serviceName, serviceVersion, db, cfg, registry).(*DefaultHealthChecker),
		additionalInfo:       make(map[string]interface{}),
	}
}

// AddInfo adds additional information to health responses
func (h *ServiceHealthChecker) AddInfo(key string, value interface{}) {
	h.additionalInfo[key] = value
}

// GetHealthData returns health data with additional service-specific information
func (h *ServiceHealthChecker) GetHealthData() map[string]interface{} {
	data := h.DefaultHealthChecker.GetHealthData()

	// Add additional service-specific information
	for key, value := range h.additionalInfo {
		data[key] = value
	}

	return data
}

// DetailedHealthChecker provides more detailed health information
type DetailedHealthChecker struct {
	*DefaultHealthChecker
	includeSystemInfo bool
}

// NewDetailedHealthChecker creates a detailed health checker
func NewDetailedHealthChecker(serviceName, serviceVersion string, db *gorm.DB, cfg *config.Config, registry ServiceRegistry, includeSystemInfo bool) HealthChecker {
	return &DetailedHealthChecker{
		DefaultHealthChecker: NewDefaultHealthChecker(serviceName, serviceVersion, db, cfg, registry).(*DefaultHealthChecker),
		includeSystemInfo:    includeSystemInfo,
	}
}

// GetHealthData returns detailed health information
func (h *DetailedHealthChecker) GetHealthData() map[string]interface{} {
	data := h.DefaultHealthChecker.GetHealthData()

	if h.includeSystemInfo {
		// Add system information
		data["system"] = map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"timezone":  time.Now().Location().String(),
		}

		// Add database information if available
		if h.db != nil {
			sqlDB, err := h.db.DB()
			if err == nil {
				stats := sqlDB.Stats()
				data["database_stats"] = map[string]interface{}{
					"open_connections":     stats.OpenConnections,
					"in_use":              stats.InUse,
					"idle":                stats.Idle,
					"wait_count":          stats.WaitCount,
					"wait_duration":       stats.WaitDuration.String(),
					"max_idle_closed":     stats.MaxIdleClosed,
					"max_idle_time_closed": stats.MaxIdleTimeClosed,
					"max_lifetime_closed": stats.MaxLifetimeClosed,
				}
			}
		}
	}

	return data
}