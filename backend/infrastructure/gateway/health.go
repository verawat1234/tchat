package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"tchat.dev/shared/logger"
)

// HealthChecker monitors service health
type HealthChecker struct {
	service    *ServiceInstance
	logger     *logger.TchatLogger
	interval   time.Duration
	timeout    time.Duration
	retryCount int
	maxRetries int
	stopChan   chan struct{}
	client     *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(service *ServiceInstance, log *logger.TchatLogger) *HealthChecker {
	return &HealthChecker{
		service:    service,
		logger:     log,
		interval:   30 * time.Second,
		timeout:    5 * time.Second,
		maxRetries: 3,
		stopChan:   make(chan struct{}),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Start begins health checking for the service
func (hc *HealthChecker) Start(registry *ServiceRegistry) {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	// Initial health check
	hc.checkHealth(registry)

	for {
		select {
		case <-ticker.C:
			hc.checkHealth(registry)
		case <-hc.stopChan:
			return
		}
	}
}

// Stop stops the health checker
func (hc *HealthChecker) Stop() {
	close(hc.stopChan)
}

// checkHealth performs a health check on the service
func (hc *HealthChecker) checkHealth(registry *ServiceRegistry) {
	healthURL := fmt.Sprintf("http://%s:%d/health", hc.service.Host, hc.service.Port)

	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		hc.markUnhealthy(registry, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		hc.retryCount++
		if hc.retryCount >= hc.maxRetries {
			hc.markUnhealthy(registry, fmt.Sprintf("Health check failed after %d retries: %v", hc.maxRetries, err))
			hc.retryCount = 0
		} else {
			hc.logger.WithFields(logrus.Fields{
				"service":     hc.service.Name,
				"service_id":  hc.service.ID,
				"retry_count": hc.retryCount,
				"error":       err.Error(),
			}).Warn("Health check failed, retrying")
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		hc.markHealthy(registry)
		hc.retryCount = 0
	} else {
		hc.retryCount++
		if hc.retryCount >= hc.maxRetries {
			hc.markUnhealthy(registry, fmt.Sprintf("Health check returned status %d after %d retries", resp.StatusCode, hc.maxRetries))
			hc.retryCount = 0
		}
	}
}

// markHealthy marks the service as healthy
func (hc *HealthChecker) markHealthy(registry *ServiceRegistry) {
	if hc.service.Health != string(Healthy) {
		hc.logger.WithFields(logrus.Fields{
			"service":    hc.service.Name,
			"service_id": hc.service.ID,
			"host":       hc.service.Host,
			"port":       hc.service.Port,
		}).Info("Service is now healthy")
	}

	registry.UpdateServiceHealth(hc.service.ID, Healthy)
}

// markUnhealthy marks the service as unhealthy
func (hc *HealthChecker) markUnhealthy(registry *ServiceRegistry, reason string) {
	if hc.service.Health != string(Unhealthy) {
		hc.logger.WithFields(logrus.Fields{
			"service":    hc.service.Name,
			"service_id": hc.service.ID,
			"host":       hc.service.Host,
			"port":       hc.service.Port,
			"reason":     reason,
		}).Error("Service is now unhealthy")
	}

	registry.UpdateServiceHealth(hc.service.ID, Unhealthy)
}

// CircuitBreaker implements circuit breaker pattern for service calls
type CircuitBreaker struct {
	serviceName    string
	maxFailures    int
	resetTimeout   time.Duration
	state          CircuitState
	failures       int
	lastFailTime   time.Time
	successCount   int
	halfOpenRequests int
	maxHalfOpenRequests int
}

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	Closed CircuitState = iota
	Open
	HalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(serviceName string, maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		serviceName:         serviceName,
		maxFailures:         maxFailures,
		resetTimeout:        resetTimeout,
		state:              Closed,
		maxHalfOpenRequests: 3,
	}
}

// CanExecute determines if a request can be executed
func (cb *CircuitBreaker) CanExecute() bool {
	switch cb.state {
	case Closed:
		return true
	case Open:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = HalfOpen
			cb.halfOpenRequests = 0
			return true
		}
		return false
	case HalfOpen:
		return cb.halfOpenRequests < cb.maxHalfOpenRequests
	}
	return false
}

// OnSuccess records a successful request
func (cb *CircuitBreaker) OnSuccess() {
	switch cb.state {
	case Closed:
		cb.failures = 0
	case HalfOpen:
		cb.successCount++
		if cb.successCount >= cb.maxHalfOpenRequests {
			cb.state = Closed
			cb.failures = 0
			cb.successCount = 0
		}
	}
}

// OnFailure records a failed request
func (cb *CircuitBreaker) OnFailure() {
	cb.failures++
	cb.lastFailTime = time.Now()

	switch cb.state {
	case Closed:
		if cb.failures >= cb.maxFailures {
			cb.state = Open
		}
	case HalfOpen:
		cb.state = Open
		cb.successCount = 0
	}
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// HealthMonitor aggregates health information across all services
type HealthMonitor struct {
	registry *ServiceRegistry
	logger   *logger.TchatLogger
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(registry *ServiceRegistry, logger *logger.TchatLogger) *HealthMonitor {
	return &HealthMonitor{
		registry: registry,
		logger:   logger,
	}
}

// GetOverallHealth returns the overall health status
func (hm *HealthMonitor) GetOverallHealth() map[string]interface{} {
	services := hm.registry.GetAllServices()

	totalServices := len(services)
	healthyServices := 0
	unhealthyServices := 0
	unknownServices := 0

	serviceStatus := make(map[string]interface{})

	for _, service := range services {
		switch service.Health {
		case string(Healthy):
			healthyServices++
		case string(Unhealthy):
			unhealthyServices++
		case string(Unknown):
			unknownServices++
		}

		serviceStatus[service.Name] = map[string]interface{}{
			"id":        service.ID,
			"health":    service.Health,
			"host":      service.Host,
			"port":      service.Port,
			"version":   service.Version,
			"tags":      service.Tags,
			"last_seen": service.LastSeen,
		}
	}

	overallStatus := "healthy"
	if healthyServices == 0 && totalServices > 0 {
		overallStatus = "critical"
	} else if float64(healthyServices)/float64(totalServices) < 0.5 {
		overallStatus = "degraded"
	}

	return map[string]interface{}{
		"overall_status":    overallStatus,
		"total_services":    totalServices,
		"healthy_services":  healthyServices,
		"unhealthy_services": unhealthyServices,
		"unknown_services":  unknownServices,
		"services":          serviceStatus,
		"timestamp":         time.Now().UTC(),
	}
}

// LogHealthSummary logs a summary of service health
func (hm *HealthMonitor) LogHealthSummary() {
	health := hm.GetOverallHealth()

	hm.logger.WithFields(logrus.Fields{
		"overall_status":     health["overall_status"],
		"total_services":     health["total_services"],
		"healthy_services":   health["healthy_services"],
		"unhealthy_services": health["unhealthy_services"],
		"unknown_services":   health["unknown_services"],
	}).Info("Health summary")
}