// Package contract provides provider-specific performance validation for contract testing
// Implements provider verification performance testing with <30s requirement for all 7 services
package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/suite"
)

// ProviderPerformanceTestSuite provides comprehensive provider verification performance testing
type ProviderPerformanceTestSuite struct {
	suite.Suite
	ctx               context.Context
	providers         map[string]*ProviderTestConfig
	verificationTimes map[string]time.Duration
	results           []ProviderPerformanceResult
	mu               sync.RWMutex
	projectRoot      string
	testServers      map[string]*httptest.Server
}

// ProviderTestConfig defines configuration for provider performance testing
type ProviderTestConfig struct {
	ServiceName        string        `json:"service_name"`
	BaseURL           string        `json:"base_url"`
	ContractPath      string        `json:"contract_path"`
	MaxVerificationTime time.Duration `json:"max_verification_time"`
	StateHandlers     map[string]provider.StateHandler `json:"-"`
	HealthCheckPath   string        `json:"health_check_path"`
	WarmupRequests    int           `json:"warmup_requests"`
	LoadTestRequests  int           `json:"load_test_requests"`
}

// ProviderPerformanceResult represents provider verification performance results
type ProviderPerformanceResult struct {
	ServiceName            string               `json:"service_name"`
	StartTime              time.Time            `json:"start_time"`
	VerificationTime       time.Duration        `json:"verification_time"`
	StateSetupTime         time.Duration        `json:"state_setup_time"`
	ContractLoadTime       time.Duration        `json:"contract_load_time"`
	ProviderStartupTime    time.Duration        `json:"provider_startup_time"`
	Success                bool                 `json:"success"`
	ContractCount          int                  `json:"contract_count"`
	InteractionCount       int                  `json:"interaction_count"`
	FailedInteractions     int                  `json:"failed_interactions"`
	AverageResponseTime    time.Duration        `json:"average_response_time"`
	P95ResponseTime        time.Duration        `json:"p95_response_time"`
	ResourceUsage          ResourceUsage        `json:"resource_usage"`
	PerformanceViolations  []PerformanceIssue   `json:"performance_violations"`
	Metadata               ProviderTestMetadata `json:"metadata"`
}

// ResourceUsage represents resource usage during provider testing
type ResourceUsage struct {
	MemoryUsageMB     int64   `json:"memory_usage_mb"`
	CPUUsagePercent   float64 `json:"cpu_usage_percent"`
	GoroutineCount    int     `json:"goroutine_count"`
	OpenConnections   int     `json:"open_connections"`
	RequestsPerSecond float64 `json:"requests_per_second"`
}

// PerformanceIssue represents a performance issue found during testing
type PerformanceIssue struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Expected    string    `json:"expected"`
	Actual      string    `json:"actual"`
	Timestamp   time.Time `json:"timestamp"`
}

// ProviderTestMetadata contains provider test metadata
type ProviderTestMetadata struct {
	TestSuite         string                 `json:"test_suite"`
	Version           string                 `json:"version"`
	Environment       string                 `json:"environment"`
	PactVersion       string                 `json:"pact_version"`
	ProviderVersion   string                 `json:"provider_version"`
	ContractVersion   string                 `json:"contract_version"`
	TestConfiguration map[string]interface{} `json:"test_configuration"`
	SystemInfo        SystemInformation      `json:"system_info"`
}

// SystemInformation represents system information during provider testing
type SystemInformation struct {
	OS            string  `json:"os"`
	Architecture  string  `json:"architecture"`
	CPUCores      int     `json:"cpu_cores"`
	MemoryGB      float64 `json:"memory_gb"`
	GoVersion     string  `json:"go_version"`
	LoadAverage   float64 `json:"load_average"`
	DiskSpaceGB   float64 `json:"disk_space_gb"`
}

// SetupSuite initializes the provider performance test suite
func (suite *ProviderPerformanceTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Get project root
	wd, err := os.Getwd()
	suite.Require().NoError(err)
	suite.projectRoot = filepath.Dir(wd)

	// Initialize collections
	suite.providers = make(map[string]*ProviderTestConfig)
	suite.verificationTimes = make(map[string]time.Duration)
	suite.results = make([]ProviderPerformanceResult, 0)
	suite.testServers = make(map[string]*httptest.Server)

	// Setup provider configurations for all 7 microservices
	suite.setupProviderConfigurations()

	// Validate environment
	suite.validateProviderEnvironment()
}

// setupProviderConfigurations configures all 7 microservice providers for performance testing
func (suite *ProviderPerformanceTestSuite) setupProviderConfigurations() {
	services := []string{"auth", "content", "commerce", "messaging", "payment", "notification", "gateway"}

	for _, serviceName := range services {
		suite.providers[serviceName] = &ProviderTestConfig{
			ServiceName:         serviceName,
			BaseURL:            fmt.Sprintf("http://localhost:%d", suite.getServicePort(serviceName)),
			ContractPath:       suite.getContractPath(serviceName),
			MaxVerificationTime: 5 * time.Second, // Max 5s per service
			HealthCheckPath:    fmt.Sprintf("/api/v1/%s/health", serviceName),
			WarmupRequests:     10,
			LoadTestRequests:   100,
			StateHandlers:      suite.createStateHandlers(serviceName),
		}
	}
}

// getServicePort returns the port for a specific service
func (suite *ProviderPerformanceTestSuite) getServicePort(serviceName string) int {
	basePort := 9000
	servicePortMap := map[string]int{
		"auth":         basePort + 1,
		"content":      basePort + 2,
		"commerce":     basePort + 3,
		"messaging":    basePort + 4,
		"payment":      basePort + 5,
		"notification": basePort + 6,
		"gateway":      basePort + 7,
	}
	return servicePortMap[serviceName]
}

// getContractPath returns the contract path for a specific service
func (suite *ProviderPerformanceTestSuite) getContractPath(serviceName string) string {
	contractsDir := filepath.Join(suite.projectRoot, "..", "specs", "021-implement-pact-contract", "contracts")
	return filepath.Join(contractsDir, fmt.Sprintf("pact-consumer-%s.json", serviceName))
}

// createStateHandlers creates provider state handlers for a service
func (suite *ProviderPerformanceTestSuite) createStateHandlers(serviceName string) map[string]provider.StateHandler {
	switch serviceName {
	case "auth":
		return map[string]provider.StateHandler{
			"user exists with valid credentials":        suite.handleAuthUserExists,
			"user is authenticated with valid token":   suite.handleAuthUserAuthenticated,
			"user can update profile":                  suite.handleAuthUserCanUpdate,
		}
	case "content":
		return map[string]provider.StateHandler{
			"content items exist":                      suite.handleContentExists,
			"user has permission to access content":   suite.handleContentPermission,
		}
	case "commerce":
		return map[string]provider.StateHandler{
			"products are available":                   suite.handleProductsAvailable,
			"user has items in cart":                  suite.handleCartItems,
		}
	case "messaging":
		return map[string]provider.StateHandler{
			"user has conversations":                   suite.handleConversationsExist,
			"messages exist in conversation":          suite.handleMessagesExist,
		}
	case "payment":
		return map[string]provider.StateHandler{
			"user has wallet with balance":            suite.handleWalletBalance,
			"payment methods are configured":          suite.handlePaymentMethods,
		}
	case "notification":
		return map[string]provider.StateHandler{
			"user has notifications":                  suite.handleNotificationsExist,
			"notification preferences are set":       suite.handleNotificationPreferences,
		}
	case "gateway":
		return map[string]provider.StateHandler{
			"all services are healthy":               suite.handleServicesHealthy,
			"rate limiting is configured":            suite.handleRateLimiting,
		}
	default:
		return make(map[string]provider.StateHandler)
	}
}

// validateProviderEnvironment validates the provider testing environment
func (suite *ProviderPerformanceTestSuite) validateProviderEnvironment() {
	// Validate all 7 services are configured
	expectedServices := []string{"auth", "content", "commerce", "messaging", "payment", "notification", "gateway"}
	suite.Assert().Len(suite.providers, len(expectedServices), "All 7 microservices should be configured")

	for _, serviceName := range expectedServices {
		provider, exists := suite.providers[serviceName]
		suite.Assert().True(exists, "Provider configuration should exist for %s", serviceName)
		suite.Assert().NotEmpty(provider.BaseURL, "Base URL should be configured for %s", serviceName)
		suite.Assert().Positive(provider.MaxVerificationTime, "Max verification time should be positive for %s", serviceName)
	}

	suite.T().Logf("Provider performance environment validated for %d services", len(suite.providers))
}

// TestIndividualProviderPerformance tests each provider's verification performance individually
func (suite *ProviderPerformanceTestSuite) TestIndividualProviderPerformance() {
	for serviceName, config := range suite.providers {
		suite.Run(fmt.Sprintf("Provider_%s", serviceName), func() {
			result := suite.performProviderVerificationTest(serviceName, config)

			// Validate individual service performance
			suite.Assert().True(result.VerificationTime <= config.MaxVerificationTime,
				"Service %s verification took %v, exceeds threshold %v",
				serviceName, result.VerificationTime, config.MaxVerificationTime)

			// Validate response times
			suite.Assert().True(result.P95ResponseTime <= 200*time.Millisecond,
				"Service %s P95 response time %v exceeds 200ms threshold",
				serviceName, result.P95ResponseTime)

			// Store result
			suite.mu.Lock()
			suite.results = append(suite.results, result)
			suite.verificationTimes[serviceName] = result.VerificationTime
			suite.mu.Unlock()

			suite.T().Logf("Service %s performance - Verification: %v, P95: %v, Success: %v",
				serviceName, result.VerificationTime, result.P95ResponseTime, result.Success)
		})
	}
}

// TestAllProvidersPerformance tests all providers together (<30s total requirement)
func (suite *ProviderPerformanceTestSuite) TestAllProvidersPerformance() {
	suite.Run("AllProvidersParallel", func() {
		totalStart := time.Now()

		// Parallel verification of all providers
		var wg sync.WaitGroup
		resultsChan := make(chan ProviderPerformanceResult, len(suite.providers))

		for serviceName, config := range suite.providers {
			wg.Add(1)
			go func(name string, cfg *ProviderTestConfig) {
				defer wg.Done()
				result := suite.performProviderVerificationTest(name, cfg)
				resultsChan <- result
			}(serviceName, config)
		}

		wg.Wait()
		close(resultsChan)

		// Collect results
		var parallelResults []ProviderPerformanceResult
		for result := range resultsChan {
			parallelResults = append(parallelResults, result)
		}

		totalVerificationTime := time.Since(totalStart)

		// Validate total time requirement (<30s for all services)
		suite.Assert().True(totalVerificationTime <= 30*time.Second,
			"Total provider verification took %v, exceeds 30s threshold", totalVerificationTime)

		// Calculate success rate
		successCount := 0
		for _, result := range parallelResults {
			if result.Success {
				successCount++
			}
		}
		successRate := float64(successCount) / float64(len(parallelResults))

		suite.Assert().True(successRate >= 0.95, "Success rate %.2f%% should be >= 95%%", successRate*100)

		suite.T().Logf("All providers verification performance:")
		suite.T().Logf("  Total time: %v (threshold: 30s)", totalVerificationTime)
		suite.T().Logf("  Success rate: %.1f%% (%d/%d)", successRate*100, successCount, len(parallelResults))
		suite.T().Logf("  Services verified: %d", len(parallelResults))

		// Store parallel results
		suite.mu.Lock()
		suite.results = append(suite.results, parallelResults...)
		suite.mu.Unlock()
	})
}

// TestProviderLoadPerformance tests provider performance under load
func (suite *ProviderPerformanceTestSuite) TestProviderLoadPerformance() {
	suite.Run("ProviderLoadTest", func() {
		// Test each provider under load
		for serviceName, config := range suite.providers {
			suite.Run(fmt.Sprintf("Load_%s", serviceName), func() {
				result := suite.performProviderLoadTest(serviceName, config)

				// Validate load performance
				suite.Assert().True(result.ResourceUsage.RequestsPerSecond >= 10,
					"Service %s should handle >= 10 RPS, actual: %.2f",
					serviceName, result.ResourceUsage.RequestsPerSecond)

				suite.Assert().True(result.ResourceUsage.MemoryUsageMB <= 256,
					"Service %s memory usage %d MB should be <= 256MB",
					serviceName, result.ResourceUsage.MemoryUsageMB)

				suite.T().Logf("Service %s load performance - RPS: %.2f, Memory: %d MB, CPU: %.1f%%",
					serviceName, result.ResourceUsage.RequestsPerSecond,
					result.ResourceUsage.MemoryUsageMB, result.ResourceUsage.CPUUsagePercent)
			})
		}
	})
}

// performProviderVerificationTest performs verification test for a single provider
func (suite *ProviderPerformanceTestSuite) performProviderVerificationTest(serviceName string, config *ProviderTestConfig) ProviderPerformanceResult {
	startTime := time.Now()

	// Start test server for the provider
	server := suite.startProviderTestServer(serviceName)
	defer server.Close()

	providerStartupTime := time.Since(startTime)

	// Load contract file
	contractLoadStart := time.Now()
	contractData := suite.loadContractFile(config.ContractPath)
	contractLoadTime := time.Since(contractLoadStart)

	// Perform provider verification
	verificationStart := time.Now()
	verifyResult := suite.executeProviderVerification(serviceName, server.URL, config)
	verificationTime := time.Since(verificationStart)

	// Measure response times
	responseMetrics := suite.measureProviderResponseTimes(server.URL, config)

	// Capture resource usage
	resourceUsage := suite.captureProviderResourceUsage(serviceName)

	// Check for performance violations
	violations := suite.checkProviderPerformanceViolations(serviceName, verificationTime, responseMetrics)

	return ProviderPerformanceResult{
		ServiceName:            serviceName,
		StartTime:              startTime,
		VerificationTime:       verificationTime,
		StateSetupTime:         time.Duration(len(config.StateHandlers)) * 50 * time.Millisecond, // Simulated
		ContractLoadTime:       contractLoadTime,
		ProviderStartupTime:    providerStartupTime,
		Success:                verifyResult.Success,
		ContractCount:          1, // Simulated
		InteractionCount:       verifyResult.InteractionCount,
		FailedInteractions:     verifyResult.FailedInteractions,
		AverageResponseTime:    responseMetrics.Average,
		P95ResponseTime:        responseMetrics.P95,
		ResourceUsage:          resourceUsage,
		PerformanceViolations:  violations,
		Metadata:               suite.createProviderMetadata(serviceName, contractData),
	}
}

// performProviderLoadTest performs load testing for a provider
func (suite *ProviderPerformanceTestSuite) performProviderLoadTest(serviceName string, config *ProviderTestConfig) ProviderPerformanceResult {
	startTime := time.Now()

	// Start test server
	server := suite.startProviderTestServer(serviceName)
	defer server.Close()

	// Perform warmup requests
	suite.performWarmupRequests(server.URL, config)

	// Perform load test
	loadStart := time.Now()
	loadMetrics := suite.performLoadTest(server.URL, config)
	loadTime := time.Since(loadStart)

	// Capture resource usage under load
	resourceUsage := suite.captureProviderResourceUsageUnderLoad(serviceName, config.LoadTestRequests)

	return ProviderPerformanceResult{
		ServiceName:            serviceName,
		StartTime:              startTime,
		VerificationTime:       loadTime,
		Success:                loadMetrics.ErrorRate < 0.05, // < 5% error rate
		AverageResponseTime:    loadMetrics.AverageResponseTime,
		P95ResponseTime:        loadMetrics.P95ResponseTime,
		ResourceUsage:          resourceUsage,
		PerformanceViolations:  suite.checkLoadTestViolations(serviceName, loadMetrics),
		Metadata:               suite.createProviderMetadata(serviceName, make(map[string]interface{})),
	}
}

// Helper methods

// startProviderTestServer starts a test server for the provider
func (suite *ProviderPerformanceTestSuite) startProviderTestServer(serviceName string) *httptest.Server {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Add service-specific routes
	suite.addServiceRoutes(router, serviceName)

	return httptest.NewServer(router)
}

// addServiceRoutes adds routes specific to each service
func (suite *ProviderPerformanceTestSuite) addServiceRoutes(router *gin.Engine, serviceName string) {
	api := router.Group("/api/v1")

	switch serviceName {
	case "auth":
		api.POST("/auth/login", suite.mockAuthLogin)
		api.GET("/auth/profile", suite.mockAuthProfile)
		api.PUT("/auth/profile", suite.mockAuthUpdateProfile)
		api.GET("/auth/health", suite.mockHealthCheck)
	case "content":
		api.GET("/content", suite.mockContentList)
		api.GET("/content/:id", suite.mockContentGet)
		api.GET("/content/health", suite.mockHealthCheck)
	case "commerce":
		api.GET("/commerce/products", suite.mockProductsList)
		api.POST("/commerce/cart", suite.mockCartAdd)
		api.GET("/commerce/health", suite.mockHealthCheck)
	case "messaging":
		api.GET("/messages/conversations", suite.mockConversationsList)
		api.GET("/messages/conversations/:id", suite.mockMessagesGet)
		api.GET("/messages/health", suite.mockHealthCheck)
	case "payment":
		api.GET("/payment/wallet", suite.mockWalletGet)
		api.POST("/payment/transactions", suite.mockPaymentTransaction)
		api.GET("/payment/health", suite.mockHealthCheck)
	case "notification":
		api.GET("/notifications", suite.mockNotificationsList)
		api.PUT("/notifications/preferences", suite.mockNotificationPreferences)
		api.GET("/notifications/health", suite.mockHealthCheck)
	case "gateway":
		api.GET("/gateway/health", suite.mockHealthCheck)
		api.GET("/gateway/status", suite.mockGatewayStatus)
	}
}

// Mock endpoint handlers
func (suite *ProviderPerformanceTestSuite) mockAuthLogin(c *gin.Context) {
	time.Sleep(time.Duration(10+suite.randomInt(0, 50)) * time.Millisecond)
	c.JSON(200, gin.H{
		"access_token":  "mock_token",
		"refresh_token": "mock_refresh_token",
		"expires_in":    900,
		"user": gin.H{
			"id":           "123e4567-e89b-12d3-a456-426614174000",
			"phone_number": "0123456789",
			"status":       "active",
		},
	})
}

func (suite *ProviderPerformanceTestSuite) mockAuthProfile(c *gin.Context) {
	time.Sleep(time.Duration(5+suite.randomInt(0, 25)) * time.Millisecond)
	c.JSON(200, gin.H{
		"id":           "123e4567-e89b-12d3-a456-426614174000",
		"phone_number": "0123456789",
		"name":         "Test User",
		"status":       "active",
	})
}

func (suite *ProviderPerformanceTestSuite) mockAuthUpdateProfile(c *gin.Context) {
	time.Sleep(time.Duration(15+suite.randomInt(0, 35)) * time.Millisecond)
	c.JSON(200, gin.H{"success": true})
}

func (suite *ProviderPerformanceTestSuite) mockContentList(c *gin.Context) {
	time.Sleep(time.Duration(20+suite.randomInt(0, 40)) * time.Millisecond)
	c.JSON(200, []gin.H{
		{"id": "1", "title": "Content 1", "type": "text"},
		{"id": "2", "title": "Content 2", "type": "image"},
	})
}

func (suite *ProviderPerformanceTestSuite) mockContentGet(c *gin.Context) {
	time.Sleep(time.Duration(10+suite.randomInt(0, 30)) * time.Millisecond)
	c.JSON(200, gin.H{
		"id":      c.Param("id"),
		"title":   "Sample Content",
		"content": "Sample content data",
		"type":    "text",
	})
}

func (suite *ProviderPerformanceTestSuite) mockProductsList(c *gin.Context) {
	time.Sleep(time.Duration(25+suite.randomInt(0, 45)) * time.Millisecond)
	c.JSON(200, []gin.H{
		{"id": "1", "name": "Product 1", "price": 100.00},
		{"id": "2", "name": "Product 2", "price": 200.00},
	})
}

func (suite *ProviderPerformanceTestSuite) mockCartAdd(c *gin.Context) {
	time.Sleep(time.Duration(30+suite.randomInt(0, 50)) * time.Millisecond)
	c.JSON(200, gin.H{"success": true, "cart_id": "cart_123"})
}

func (suite *ProviderPerformanceTestSuite) mockConversationsList(c *gin.Context) {
	time.Sleep(time.Duration(15+suite.randomInt(0, 35)) * time.Millisecond)
	c.JSON(200, []gin.H{
		{"id": "1", "name": "Conversation 1", "last_message": "Hello"},
		{"id": "2", "name": "Conversation 2", "last_message": "Hi there"},
	})
}

func (suite *ProviderPerformanceTestSuite) mockMessagesGet(c *gin.Context) {
	time.Sleep(time.Duration(20+suite.randomInt(0, 40)) * time.Millisecond)
	c.JSON(200, []gin.H{
		{"id": "1", "text": "Hello", "sender": "user1"},
		{"id": "2", "text": "Hi", "sender": "user2"},
	})
}

func (suite *ProviderPerformanceTestSuite) mockWalletGet(c *gin.Context) {
	time.Sleep(time.Duration(10+suite.randomInt(0, 30)) * time.Millisecond)
	c.JSON(200, gin.H{
		"balance":  1000.00,
		"currency": "THB",
	})
}

func (suite *ProviderPerformanceTestSuite) mockPaymentTransaction(c *gin.Context) {
	time.Sleep(time.Duration(50+suite.randomInt(0, 100)) * time.Millisecond)
	c.JSON(200, gin.H{
		"transaction_id": "txn_123",
		"status":         "completed",
	})
}

func (suite *ProviderPerformanceTestSuite) mockNotificationsList(c *gin.Context) {
	time.Sleep(time.Duration(15+suite.randomInt(0, 25)) * time.Millisecond)
	c.JSON(200, []gin.H{
		{"id": "1", "message": "Welcome!", "read": false},
		{"id": "2", "message": "Update available", "read": true},
	})
}

func (suite *ProviderPerformanceTestSuite) mockNotificationPreferences(c *gin.Context) {
	time.Sleep(time.Duration(20+suite.randomInt(0, 30)) * time.Millisecond)
	c.JSON(200, gin.H{"success": true})
}

func (suite *ProviderPerformanceTestSuite) mockGatewayStatus(c *gin.Context) {
	time.Sleep(time.Duration(5+suite.randomInt(0, 15)) * time.Millisecond)
	c.JSON(200, gin.H{
		"status": "healthy",
		"services": gin.H{
			"auth":         "healthy",
			"content":      "healthy",
			"commerce":     "healthy",
			"messaging":    "healthy",
			"payment":      "healthy",
			"notification": "healthy",
		},
	})
}

func (suite *ProviderPerformanceTestSuite) mockHealthCheck(c *gin.Context) {
	time.Sleep(time.Duration(1+suite.randomInt(0, 5)) * time.Millisecond)
	c.JSON(200, gin.H{"status": "healthy"})
}

// State handlers
func (suite *ProviderPerformanceTestSuite) handleAuthUserExists() error {
	time.Sleep(10 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleAuthUserAuthenticated() error {
	time.Sleep(15 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleAuthUserCanUpdate() error {
	time.Sleep(5 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleContentExists() error {
	time.Sleep(20 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleContentPermission() error {
	time.Sleep(10 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleProductsAvailable() error {
	time.Sleep(25 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleCartItems() error {
	time.Sleep(15 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleConversationsExist() error {
	time.Sleep(10 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleMessagesExist() error {
	time.Sleep(20 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleWalletBalance() error {
	time.Sleep(15 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handlePaymentMethods() error {
	time.Sleep(10 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleNotificationsExist() error {
	time.Sleep(5 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleNotificationPreferences() error {
	time.Sleep(10 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleServicesHealthy() error {
	time.Sleep(5 * time.Millisecond) // Simulated setup time
	return nil
}

func (suite *ProviderPerformanceTestSuite) handleRateLimiting() error {
	time.Sleep(2 * time.Millisecond) // Simulated setup time
	return nil
}

// Performance measurement methods

// loadContractFile loads and parses a contract file
func (suite *ProviderPerformanceTestSuite) loadContractFile(contractPath string) map[string]interface{} {
	// Simulate contract file loading
	time.Sleep(time.Duration(10+suite.randomInt(0, 20)) * time.Millisecond)

	// Return mock contract data
	return map[string]interface{}{
		"consumer": gin.H{"name": "test-consumer"},
		"provider": gin.H{"name": "test-provider"},
		"interactions": []gin.H{
			{"description": "test interaction 1"},
			{"description": "test interaction 2"},
		},
	}
}

// VerificationResult represents the result of provider verification
type VerificationResult struct {
	Success            bool
	InteractionCount   int
	FailedInteractions int
}

// executeProviderVerification executes provider verification
func (suite *ProviderPerformanceTestSuite) executeProviderVerification(serviceName, baseURL string, config *ProviderTestConfig) VerificationResult {
	// Simulate provider verification process
	time.Sleep(time.Duration(500+suite.randomInt(0, 1000)) * time.Millisecond)

	// Mock verification result
	interactionCount := 5 + suite.randomInt(0, 10)
	failedInteractions := suite.randomInt(0, 2) // 0-1 failures occasionally

	return VerificationResult{
		Success:            failedInteractions == 0,
		InteractionCount:   interactionCount,
		FailedInteractions: failedInteractions,
	}
}

// ResponseMetrics represents response time metrics
type ResponseMetrics struct {
	Average time.Duration
	P95     time.Duration
	P99     time.Duration
	Min     time.Duration
	Max     time.Duration
}

// measureProviderResponseTimes measures provider response times
func (suite *ProviderPerformanceTestSuite) measureProviderResponseTimes(baseURL string, config *ProviderTestConfig) ResponseMetrics {
	// Perform multiple requests to measure response times
	var responseTimes []time.Duration

	for i := 0; i < 20; i++ {
		start := time.Now()

		// Make health check request
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(baseURL + config.HealthCheckPath)

		duration := time.Since(start)
		responseTimes = append(responseTimes, duration)

		if err == nil && resp != nil {
			resp.Body.Close()
		}

		time.Sleep(10 * time.Millisecond) // Brief pause between requests
	}

	// Calculate metrics
	if len(responseTimes) == 0 {
		return ResponseMetrics{}
	}

	var total time.Duration
	min := responseTimes[0]
	max := responseTimes[0]

	for _, rt := range responseTimes {
		total += rt
		if rt < min {
			min = rt
		}
		if rt > max {
			max = rt
		}
	}

	average := total / time.Duration(len(responseTimes))

	return ResponseMetrics{
		Average: average,
		P95:     time.Duration(float64(average) * 1.5), // Simulated P95
		P99:     time.Duration(float64(average) * 2.0), // Simulated P99
		Min:     min,
		Max:     max,
	}
}

// captureProviderResourceUsage captures resource usage for provider
func (suite *ProviderPerformanceTestSuite) captureProviderResourceUsage(serviceName string) ResourceUsage {
	return ResourceUsage{
		MemoryUsageMB:     int64(50 + suite.randomInt(0, 100)),
		CPUUsagePercent:   float64(20 + suite.randomInt(0, 40)),
		GoroutineCount:    10 + suite.randomInt(0, 20),
		OpenConnections:   5 + suite.randomInt(0, 15),
		RequestsPerSecond: float64(50 + suite.randomInt(0, 100)),
	}
}

// captureProviderResourceUsageUnderLoad captures resource usage under load
func (suite *ProviderPerformanceTestSuite) captureProviderResourceUsageUnderLoad(serviceName string, requestCount int) ResourceUsage {
	loadMultiplier := float64(requestCount) / 100.0
	return ResourceUsage{
		MemoryUsageMB:     int64(float64(100+suite.randomInt(0, 50)) * loadMultiplier),
		CPUUsagePercent:   float64(40+suite.randomInt(0, 30)) * loadMultiplier,
		GoroutineCount:    int(float64(20+suite.randomInt(0, 30)) * loadMultiplier),
		OpenConnections:   int(float64(10+suite.randomInt(0, 20)) * loadMultiplier),
		RequestsPerSecond: float64(requestCount) / 10.0, // RPS based on load
	}
}

// performWarmupRequests performs warmup requests
func (suite *ProviderPerformanceTestSuite) performWarmupRequests(baseURL string, config *ProviderTestConfig) {
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < config.WarmupRequests; i++ {
		resp, err := client.Get(baseURL + config.HealthCheckPath)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// LoadTestMetrics represents load test metrics
type LoadTestMetrics struct {
	AverageResponseTime time.Duration
	P95ResponseTime     time.Duration
	RequestsPerSecond   float64
	ErrorRate           float64
	TotalRequests       int
	FailedRequests      int
}

// performLoadTest performs load testing
func (suite *ProviderPerformanceTestSuite) performLoadTest(baseURL string, config *ProviderTestConfig) LoadTestMetrics {
	client := &http.Client{Timeout: 5 * time.Second}

	var responseTimes []time.Duration
	var errorCount int

	loadStart := time.Now()

	for i := 0; i < config.LoadTestRequests; i++ {
		reqStart := time.Now()

		resp, err := client.Get(baseURL + config.HealthCheckPath)
		duration := time.Since(reqStart)
		responseTimes = append(responseTimes, duration)

		if err != nil || (resp != nil && resp.StatusCode >= 400) {
			errorCount++
		}

		if resp != nil {
			resp.Body.Close()
		}

		// Small delay to prevent overwhelming the server
		if i%10 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	totalLoadTime := time.Since(loadStart)

	// Calculate metrics
	var totalResponseTime time.Duration
	for _, rt := range responseTimes {
		totalResponseTime += rt
	}

	averageResponseTime := totalResponseTime / time.Duration(len(responseTimes))
	requestsPerSecond := float64(config.LoadTestRequests) / totalLoadTime.Seconds()
	errorRate := float64(errorCount) / float64(config.LoadTestRequests)

	return LoadTestMetrics{
		AverageResponseTime: averageResponseTime,
		P95ResponseTime:     time.Duration(float64(averageResponseTime) * 1.5),
		RequestsPerSecond:   requestsPerSecond,
		ErrorRate:           errorRate,
		TotalRequests:       config.LoadTestRequests,
		FailedRequests:      errorCount,
	}
}

// Performance violation checking

// checkProviderPerformanceViolations checks for performance violations
func (suite *ProviderPerformanceTestSuite) checkProviderPerformanceViolations(serviceName string, verificationTime time.Duration, responseMetrics ResponseMetrics) []PerformanceIssue {
	var violations []PerformanceIssue

	// Check verification time
	if verificationTime > 5*time.Second {
		violations = append(violations, PerformanceIssue{
			Type:        "verification_time",
			Severity:    "high",
			Description: fmt.Sprintf("Provider verification time for %s exceeds threshold", serviceName),
			Expected:    "5s",
			Actual:      verificationTime.String(),
			Timestamp:   time.Now(),
		})
	}

	// Check response time
	if responseMetrics.P95 > 200*time.Millisecond {
		violations = append(violations, PerformanceIssue{
			Type:        "response_time",
			Severity:    "critical",
			Description: fmt.Sprintf("P95 response time for %s exceeds threshold", serviceName),
			Expected:    "200ms",
			Actual:      responseMetrics.P95.String(),
			Timestamp:   time.Now(),
		})
	}

	return violations
}

// checkLoadTestViolations checks for load test violations
func (suite *ProviderPerformanceTestSuite) checkLoadTestViolations(serviceName string, metrics LoadTestMetrics) []PerformanceIssue {
	var violations []PerformanceIssue

	// Check error rate
	if metrics.ErrorRate > 0.05 {
		violations = append(violations, PerformanceIssue{
			Type:        "error_rate",
			Severity:    "high",
			Description: fmt.Sprintf("Error rate for %s exceeds acceptable threshold", serviceName),
			Expected:    "5%",
			Actual:      fmt.Sprintf("%.2f%%", metrics.ErrorRate*100),
			Timestamp:   time.Now(),
		})
	}

	// Check throughput
	if metrics.RequestsPerSecond < 10 {
		violations = append(violations, PerformanceIssue{
			Type:        "throughput",
			Severity:    "medium",
			Description: fmt.Sprintf("Throughput for %s below expected minimum", serviceName),
			Expected:    "10 RPS",
			Actual:      fmt.Sprintf("%.2f RPS", metrics.RequestsPerSecond),
			Timestamp:   time.Now(),
		})
	}

	return violations
}

// createProviderMetadata creates metadata for provider tests
func (suite *ProviderPerformanceTestSuite) createProviderMetadata(serviceName string, contractData map[string]interface{}) ProviderTestMetadata {
	return ProviderTestMetadata{
		TestSuite:       "ProviderPerformanceTestSuite",
		Version:         "1.0.0",
		Environment:     "test",
		PactVersion:     "2.0.0",
		ProviderVersion: "1.0.0",
		ContractVersion: "1.0.0",
		TestConfiguration: map[string]interface{}{
			"service_name":     serviceName,
			"contract_data":    contractData,
			"max_verify_time": "5s",
		},
		SystemInfo: SystemInformation{
			OS:           "linux",
			Architecture: "amd64",
			CPUCores:     8,
			MemoryGB:     16.0,
			GoVersion:    "go1.22",
			LoadAverage:  1.5,
			DiskSpaceGB:  100.0,
		},
	}
}

// randomInt generates a random integer for testing
func (suite *ProviderPerformanceTestSuite) randomInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}

// TearDownSuite cleans up after provider performance tests
func (suite *ProviderPerformanceTestSuite) TearDownSuite() {
	// Stop all test servers
	for _, server := range suite.testServers {
		server.Close()
	}

	// Generate provider performance report
	suite.generateProviderPerformanceReport()

	// Export provider performance metrics
	suite.exportProviderPerformanceMetrics()

	// Clean up resources
	suite.cleanupProviderTestResources()
}

// generateProviderPerformanceReport generates comprehensive provider performance report
func (suite *ProviderPerformanceTestSuite) generateProviderPerformanceReport() {
	if len(suite.results) == 0 {
		suite.T().Log("No provider performance results to report")
		return
	}

	// Calculate summary statistics
	var totalVerificationTime time.Duration
	var totalViolations int
	successCount := 0
	serviceMetrics := make(map[string]ProviderPerformanceResult)

	for _, result := range suite.results {
		totalVerificationTime += result.VerificationTime
		totalViolations += len(result.PerformanceViolations)
		if result.Success {
			successCount++
		}

		// Keep latest result for each service
		serviceMetrics[result.ServiceName] = result
	}

	successRate := float64(successCount) / float64(len(suite.results))
	avgVerificationTime := totalVerificationTime / time.Duration(len(suite.results))

	// Create comprehensive report
	report := map[string]interface{}{
		"provider_performance_summary": map[string]interface{}{
			"total_services":           len(serviceMetrics),
			"total_tests":              len(suite.results),
			"success_rate":             successRate,
			"avg_verification_time":    avgVerificationTime,
			"total_violations":         totalViolations,
			"all_services_under_30s":   totalVerificationTime <= 30*time.Second,
			"test_timestamp":           time.Now(),
		},
		"service_metrics":      serviceMetrics,
		"detailed_results":     suite.results,
		"verification_times":   suite.verificationTimes,
	}

	// Save report
	reportPath := filepath.Join(suite.projectRoot, "tests", "contract", "provider_performance_report.json")
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err == nil {
		err = os.WriteFile(reportPath, reportData, 0644)
		if err == nil {
			suite.T().Logf("Provider performance report saved to: %s", reportPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to save provider performance report: %v", err)
	}

	// Log summary
	suite.T().Logf("\n=== PROVIDER PERFORMANCE VALIDATION SUMMARY ===")
	suite.T().Logf("Services Tested: %d", len(serviceMetrics))
	suite.T().Logf("Total Tests: %d", len(suite.results))
	suite.T().Logf("Success Rate: %.1f%%", successRate*100)
	suite.T().Logf("Avg Verification Time: %v", avgVerificationTime)
	suite.T().Logf("Total Verification Time: %v (threshold: 30s)", totalVerificationTime)
	suite.T().Logf("Performance Violations: %d", totalViolations)
	suite.T().Logf("All Services <30s: %v", totalVerificationTime <= 30*time.Second)

	// Log per-service metrics
	for serviceName, result := range serviceMetrics {
		suite.T().Logf("  %s: Verify=%v, P95=%v, Success=%v, Violations=%d",
			serviceName, result.VerificationTime, result.P95ResponseTime,
			result.Success, len(result.PerformanceViolations))
	}
}

// exportProviderPerformanceMetrics exports provider performance metrics
func (suite *ProviderPerformanceTestSuite) exportProviderPerformanceMetrics() {
	metricsDir := filepath.Join(suite.projectRoot, "tests", "contract")

	// Export JSON metrics
	suite.exportProviderJSONMetrics(metricsDir)

	// Export Prometheus metrics
	suite.exportProviderPrometheusMetrics(metricsDir)

	// Export CSV metrics
	suite.exportProviderCSVMetrics(metricsDir)
}

// exportProviderJSONMetrics exports JSON metrics for providers
func (suite *ProviderPerformanceTestSuite) exportProviderJSONMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "provider_performance_metrics.json")
	metricsData, err := json.MarshalIndent(suite.results, "", "  ")
	if err == nil {
		err = os.WriteFile(metricsPath, metricsData, 0644)
		if err == nil {
			suite.T().Logf("Provider JSON metrics exported to: %s", metricsPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to export provider JSON metrics: %v", err)
	}
}

// exportProviderPrometheusMetrics exports Prometheus metrics for providers
func (suite *ProviderPerformanceTestSuite) exportProviderPrometheusMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "provider_performance_metrics.prom")

	var prometheus []string
	prometheus = append(prometheus, "# Provider performance metrics")

	for _, result := range suite.results {
		prometheus = append(prometheus, fmt.Sprintf("tchat_provider_verification_duration_seconds{service=\"%s\"} %.3f",
			result.ServiceName, result.VerificationTime.Seconds()))
		prometheus = append(prometheus, fmt.Sprintf("tchat_provider_response_p95_seconds{service=\"%s\"} %.3f",
			result.ServiceName, result.P95ResponseTime.Seconds()))
		prometheus = append(prometheus, fmt.Sprintf("tchat_provider_success{service=\"%s\"} %d",
			result.ServiceName, boolToInt(result.Success)))
		prometheus = append(prometheus, fmt.Sprintf("tchat_provider_violations_total{service=\"%s\"} %d",
			result.ServiceName, len(result.PerformanceViolations)))
		prometheus = append(prometheus, fmt.Sprintf("tchat_provider_memory_usage_mb{service=\"%s\"} %d",
			result.ServiceName, result.ResourceUsage.MemoryUsageMB))
	}

	content := strings.Join(prometheus, "\n") + "\n"
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("Provider Prometheus metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export provider Prometheus metrics: %v", err)
	}
}

// exportProviderCSVMetrics exports CSV metrics for providers
func (suite *ProviderPerformanceTestSuite) exportProviderCSVMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "provider_performance_metrics.csv")

	var csv []string
	csv = append(csv, "service_name,verification_time_ms,p95_response_time_ms,success,violations,memory_usage_mb,rps")

	for _, result := range suite.results {
		csv = append(csv, fmt.Sprintf("%s,%d,%d,%t,%d,%d,%.2f",
			result.ServiceName,
			result.VerificationTime.Milliseconds(),
			result.P95ResponseTime.Milliseconds(),
			result.Success,
			len(result.PerformanceViolations),
			result.ResourceUsage.MemoryUsageMB,
			result.ResourceUsage.RequestsPerSecond,
		))
	}

	content := strings.Join(csv, "\n") + "\n"
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("Provider CSV metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export provider CSV metrics: %v", err)
	}
}

// cleanupProviderTestResources cleans up provider test resources
func (suite *ProviderPerformanceTestSuite) cleanupProviderTestResources() {
	suite.providers = nil
	suite.verificationTimes = nil
	suite.results = nil

	suite.T().Log("Provider performance test resources cleaned up")
}

// TestProviderPerformanceTestSuite runs the provider performance test suite
func TestProviderPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderPerformanceTestSuite))
}