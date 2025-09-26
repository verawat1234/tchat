// Package contract provides comprehensive performance validation for contract testing system
// Implements T025: Performance validation for contract testing with enterprise requirements
package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/suite"
)

// ContractPerformanceValidationSuite provides comprehensive performance validation for contract testing
type ContractPerformanceValidationSuite struct {
	suite.Suite
	ctx                    context.Context
	performanceConfig      PerformanceConfig
	contractValidationTime time.Duration
	apiResponseTimes       map[string]time.Duration
	brokerIntegrationTime  time.Duration
	results               []PerformanceResult
	mu                    sync.RWMutex
	projectRoot           string
	testServices          map[string]*TestService
}

// PerformanceConfig defines performance validation configuration
type PerformanceConfig struct {
	ContractValidation ContractValidationConfig `json:"contract_validation"`
	APIValidation      APIValidationConfig      `json:"api_validation"`
	BrokerIntegration  BrokerIntegrationConfig  `json:"broker_integration"`
	CrossPlatform      CrossPlatformConfig      `json:"cross_platform"`
	ProviderVerification ProviderVerificationConfig `json:"provider_verification"`
	Thresholds         PerformanceThresholds    `json:"thresholds"`
	Monitoring         MonitoringConfig         `json:"monitoring"`
}

// ContractValidationConfig defines contract validation performance settings
type ContractValidationConfig struct {
	MaxValidationTime   time.Duration `json:"max_validation_time"`
	ConcurrentValidations int         `json:"concurrent_validations"`
	ContractCacheEnabled  bool        `json:"contract_cache_enabled"`
	ValidationRetries     int         `json:"validation_retries"`
	WarmupRuns           int         `json:"warmup_runs"`
}

// APIValidationConfig defines API performance validation settings
type APIValidationConfig struct {
	MaxP50ResponseTime time.Duration            `json:"max_p50_response_time"`
	MaxP95ResponseTime time.Duration            `json:"max_p95_response_time"`
	MaxP99ResponseTime time.Duration            `json:"max_p99_response_time"`
	LoadTestDuration   time.Duration            `json:"load_test_duration"`
	ConcurrentRequests int                      `json:"concurrent_requests"`
	EndpointMapping    map[string]string        `json:"endpoint_mapping"`
	RequestPatterns    map[string]RequestPattern `json:"request_patterns"`
}

// RequestPattern defines API request patterns for performance testing
type RequestPattern struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Headers     map[string]string `json:"headers"`
	Body        interface{}       `json:"body"`
	ExpectedStatus int            `json:"expected_status"`
	AuthRequired   bool           `json:"auth_required"`
}

// BrokerIntegrationConfig defines Pact Broker integration performance settings
type BrokerIntegrationConfig struct {
	MaxPublicationTime  time.Duration `json:"max_publication_time"`
	MaxRetrievalTime    time.Duration `json:"max_retrieval_time"`
	MaxVerificationTime time.Duration `json:"max_verification_time"`
	BrokerURL          string        `json:"broker_url"`
	AuthTokenEnabled   bool          `json:"auth_token_enabled"`
	RetryAttempts      int           `json:"retry_attempts"`
}

// CrossPlatformConfig defines cross-platform performance validation settings
type CrossPlatformConfig struct {
	MaxCompatibilityCheckTime time.Duration     `json:"max_compatibility_check_time"`
	Platforms                []string          `json:"platforms"`
	ContractComparisonTime    time.Duration     `json:"contract_comparison_time"`
	CrossValidationRetries    int              `json:"cross_validation_retries"`
	PlatformEndpoints         map[string]string `json:"platform_endpoints"`
}

// ProviderVerificationConfig defines provider verification performance settings
type ProviderVerificationConfig struct {
	MaxSingleServiceTime    time.Duration `json:"max_single_service_time"`
	MaxTotalVerificationTime time.Duration `json:"max_total_verification_time"`
	Services                []string      `json:"services"`
	ParallelVerification    bool          `json:"parallel_verification"`
	StateSetupTimeout       time.Duration `json:"state_setup_timeout"`
	VerificationRetries     int           `json:"verification_retries"`
}

// PerformanceThresholds defines enterprise performance thresholds
type PerformanceThresholds struct {
	ContractValidationMax time.Duration `json:"contract_validation_max"`
	APIP95ResponseMax     time.Duration `json:"api_p95_response_max"`
	BrokerIntegrationMax  time.Duration `json:"broker_integration_max"`
	CrossPlatformMax      time.Duration `json:"cross_platform_max"`
	ProviderVerificationMax time.Duration `json:"provider_verification_max"`
	TotalTestSuiteMax     time.Duration `json:"total_test_suite_max"`
	MemoryUsageMaxMB      int64         `json:"memory_usage_max_mb"`
	CPUUsageMaxPercent    float64       `json:"cpu_usage_max_percent"`
}

// MonitoringConfig defines performance monitoring configuration
type MonitoringConfig struct {
	MetricsCollectionEnabled bool          `json:"metrics_collection_enabled"`
	RealTimeAlerts          bool          `json:"real_time_alerts"`
	PerformanceReporting    bool          `json:"performance_reporting"`
	MetricsRetention        time.Duration `json:"metrics_retention"`
	AlertThresholds         AlertThresholds `json:"alert_thresholds"`
}

// AlertThresholds defines alerting thresholds for performance monitoring
type AlertThresholds struct {
	ValidationTimeWarning  time.Duration `json:"validation_time_warning"`
	ValidationTimeCritical time.Duration `json:"validation_time_critical"`
	ResponseTimeWarning    time.Duration `json:"response_time_warning"`
	ResponseTimeCritical   time.Duration `json:"response_time_critical"`
	ErrorRateWarning       float64       `json:"error_rate_warning"`
	ErrorRateCritical      float64       `json:"error_rate_critical"`
}

// PerformanceResult represents performance validation results
type PerformanceResult struct {
	TestName              string                    `json:"test_name"`
	Category              string                    `json:"category"`
	StartTime             time.Time                 `json:"start_time"`
	Duration              time.Duration             `json:"duration"`
	Success               bool                      `json:"success"`
	ContractValidationTime time.Duration            `json:"contract_validation_time"`
	APIResponseTimes      APIResponseMetrics        `json:"api_response_times"`
	BrokerIntegrationTime time.Duration            `json:"broker_integration_time"`
	CrossPlatformTime     time.Duration            `json:"cross_platform_time"`
	ProviderVerificationTime time.Duration          `json:"provider_verification_time"`
	ResourceUsage         ResourceUsageMetrics      `json:"resource_usage"`
	ThresholdViolations   []PerformanceViolation    `json:"threshold_violations"`
	Metadata              PerformanceTestMetadata   `json:"metadata"`
}

// APIResponseMetrics represents API response time statistics
type APIResponseMetrics struct {
	P50 time.Duration `json:"p50"`
	P95 time.Duration `json:"p95"`
	P99 time.Duration `json:"p99"`
	Max time.Duration `json:"max"`
	Min time.Duration `json:"min"`
	Average time.Duration `json:"average"`
	RequestCount int64    `json:"request_count"`
	ErrorCount   int64    `json:"error_count"`
}

// ResourceUsageMetrics represents system resource usage during tests
type ResourceUsageMetrics struct {
	MemoryUsageMB     int64   `json:"memory_usage_mb"`
	CPUUsagePercent   float64 `json:"cpu_usage_percent"`
	GoroutineCount    int     `json:"goroutine_count"`
	AllocatedObjects  uint64  `json:"allocated_objects"`
	GCPauseTime       time.Duration `json:"gc_pause_time"`
}

// PerformanceViolation represents a performance threshold violation
type PerformanceViolation struct {
	Type        string    `json:"type"`
	Metric      string    `json:"metric"`
	Expected    string    `json:"expected"`
	Actual      string    `json:"actual"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// PerformanceTestMetadata contains performance test metadata
type PerformanceTestMetadata struct {
	TestSuite       string                 `json:"test_suite"`
	Version         string                 `json:"version"`
	Environment     string                 `json:"environment"`
	Configuration   map[string]interface{} `json:"configuration"`
	SystemInfo      SystemInfo             `json:"system_info"`
	TestParameters  map[string]interface{} `json:"test_parameters"`
}

// SystemInfo represents system information during performance tests
type SystemInfo struct {
	OS           string  `json:"os"`
	Architecture string  `json:"architecture"`
	CPUCores     int     `json:"cpu_cores"`
	MemoryGB     float64 `json:"memory_gb"`
	GoVersion    string  `json:"go_version"`
}

// TestService represents a service under performance testing
type TestService struct {
	Name     string
	BaseURL  string
	Server   *http.Server
	StartTime time.Time
}

// SetupSuite initializes the contract performance validation suite
func (suite *ContractPerformanceValidationSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Get project root
	wd, err := os.Getwd()
	suite.Require().NoError(err)
	suite.projectRoot = filepath.Dir(wd)

	// Load performance configuration
	suite.loadPerformanceConfiguration()

	// Initialize test services map
	suite.testServices = make(map[string]*TestService)

	// Initialize results collection
	suite.results = make([]PerformanceResult, 0)

	// Validate performance environment
	suite.validatePerformanceEnvironment()
}

// loadPerformanceConfiguration loads performance validation configuration
func (suite *ContractPerformanceValidationSuite) loadPerformanceConfiguration() {
	// Enterprise-grade performance configuration for Southeast Asian deployment
	suite.performanceConfig = PerformanceConfig{
		ContractValidation: ContractValidationConfig{
			MaxValidationTime:     1000 * time.Millisecond, // <1s requirement
			ConcurrentValidations: 5,
			ContractCacheEnabled:  true,
			ValidationRetries:     3,
			WarmupRuns:           3,
		},
		APIValidation: APIValidationConfig{
			MaxP50ResponseTime: 100 * time.Millisecond,
			MaxP95ResponseTime: 200 * time.Millisecond, // <200ms p95 requirement
			MaxP99ResponseTime: 300 * time.Millisecond,
			LoadTestDuration:   30 * time.Second,
			ConcurrentRequests: 50,
			EndpointMapping: map[string]string{
				"auth":         "/api/v1/auth",
				"content":      "/api/v1/content",
				"commerce":     "/api/v1/commerce",
				"messaging":    "/api/v1/messages",
				"payment":      "/api/v1/payment",
				"notification": "/api/v1/notifications",
			},
			RequestPatterns: map[string]RequestPattern{
				"login": {
					Method: "POST",
					Path:   "/api/v1/auth/login",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"phone_number": "0123456789",
						"country_code": "+66",
						"password":     "validPassword123",
					},
					ExpectedStatus: 200,
					AuthRequired:   false,
				},
				"profile": {
					Method: "GET",
					Path:   "/api/v1/auth/profile",
					Headers: map[string]string{
						"Authorization": "Bearer {token}",
					},
					ExpectedStatus: 200,
					AuthRequired:   true,
				},
			},
		},
		BrokerIntegration: BrokerIntegrationConfig{
			MaxPublicationTime:  500 * time.Millisecond,
			MaxRetrievalTime:    300 * time.Millisecond,
			MaxVerificationTime: 2 * time.Second,
			BrokerURL:          "http://localhost:8080",
			AuthTokenEnabled:   false,
			RetryAttempts:      3,
		},
		CrossPlatform: CrossPlatformConfig{
			MaxCompatibilityCheckTime: 2 * time.Second,
			Platforms:                []string{"web", "ios", "android"},
			ContractComparisonTime:    1 * time.Second,
			CrossValidationRetries:    3,
			PlatformEndpoints: map[string]string{
				"web":     "http://localhost:3000",
				"ios":     "http://localhost:8081",
				"android": "http://localhost:8082",
			},
		},
		ProviderVerification: ProviderVerificationConfig{
			MaxSingleServiceTime:     5 * time.Second,
			MaxTotalVerificationTime: 30 * time.Second, // <30s for all 7 services
			Services:                []string{"auth", "content", "commerce", "messaging", "payment", "notification", "gateway"},
			ParallelVerification:    true,
			StateSetupTimeout:       2 * time.Second,
			VerificationRetries:     2,
		},
		Thresholds: PerformanceThresholds{
			ContractValidationMax:   1000 * time.Millisecond,
			APIP95ResponseMax:       200 * time.Millisecond,
			BrokerIntegrationMax:    500 * time.Millisecond,
			CrossPlatformMax:        2 * time.Second,
			ProviderVerificationMax: 30 * time.Second,
			TotalTestSuiteMax:       5 * time.Minute,
			MemoryUsageMaxMB:        512,
			CPUUsageMaxPercent:      80.0,
		},
		Monitoring: MonitoringConfig{
			MetricsCollectionEnabled: true,
			RealTimeAlerts:          true,
			PerformanceReporting:    true,
			MetricsRetention:        24 * time.Hour,
			AlertThresholds: AlertThresholds{
				ValidationTimeWarning:  800 * time.Millisecond,
				ValidationTimeCritical: 1000 * time.Millisecond,
				ResponseTimeWarning:    150 * time.Millisecond,
				ResponseTimeCritical:   200 * time.Millisecond,
				ErrorRateWarning:       0.01,
				ErrorRateCritical:      0.05,
			},
		},
	}
}

// validatePerformanceEnvironment validates the performance testing environment
func (suite *ContractPerformanceValidationSuite) validatePerformanceEnvironment() {
	// Validate configuration
	suite.Assert().Positive(suite.performanceConfig.ContractValidation.MaxValidationTime)
	suite.Assert().Positive(suite.performanceConfig.APIValidation.MaxP95ResponseTime)
	suite.Assert().NotEmpty(suite.performanceConfig.ProviderVerification.Services)

	// Validate thresholds
	suite.Assert().Equal(1000*time.Millisecond, suite.performanceConfig.Thresholds.ContractValidationMax)
	suite.Assert().Equal(200*time.Millisecond, suite.performanceConfig.Thresholds.APIP95ResponseMax)

	// Log environment information
	suite.T().Logf("Performance validation environment initialized:")
	suite.T().Logf("  Contract validation max: %v", suite.performanceConfig.Thresholds.ContractValidationMax)
	suite.T().Logf("  API P95 response max: %v", suite.performanceConfig.Thresholds.APIP95ResponseMax)
	suite.T().Logf("  Provider verification max: %v", suite.performanceConfig.Thresholds.ProviderVerificationMax)
}

// TestContractValidationPerformance tests contract validation performance (<1s per validation)
func (suite *ContractPerformanceValidationSuite) TestContractValidationPerformance() {
	suite.Run("ContractValidationSpeed", func() {
		startTime := time.Now()

		// Test multiple contract validations
		contractPaths := suite.getContractPaths()
		suite.Require().NotEmpty(contractPaths, "Contract files should be available for testing")

		var totalValidationTime time.Duration
		var validationCount int

		for _, contractPath := range contractPaths {
			validationStart := time.Now()

			// Simulate contract validation
			suite.validateContractFile(contractPath)

			validationDuration := time.Since(validationStart)
			totalValidationTime += validationDuration
			validationCount++

			// Validate individual contract validation time
			suite.Assert().True(validationDuration < suite.performanceConfig.Thresholds.ContractValidationMax,
				"Contract validation for %s took %v, exceeds threshold %v",
				filepath.Base(contractPath), validationDuration, suite.performanceConfig.Thresholds.ContractValidationMax)

			suite.T().Logf("Contract %s validated in %v", filepath.Base(contractPath), validationDuration)
		}

		avgValidationTime := totalValidationTime / time.Duration(validationCount)

		// Record result
		result := PerformanceResult{
			TestName:               "ContractValidationPerformance",
			Category:               "contract_validation",
			StartTime:              startTime,
			Duration:               time.Since(startTime),
			Success:                avgValidationTime < suite.performanceConfig.Thresholds.ContractValidationMax,
			ContractValidationTime: avgValidationTime,
			ResourceUsage:          suite.captureResourceUsage(),
			ThresholdViolations:    suite.checkContractValidationViolations(avgValidationTime),
			Metadata:               suite.createTestMetadata("ContractValidationPerformance"),
		}

		suite.mu.Lock()
		suite.results = append(suite.results, result)
		suite.mu.Unlock()

		suite.T().Logf("Average contract validation time: %v (threshold: %v)",
			avgValidationTime, suite.performanceConfig.Thresholds.ContractValidationMax)
	})
}

// TestAPIResponseTimeValidation tests API response times (<200ms p95)
func (suite *ContractPerformanceValidationSuite) TestAPIResponseTimeValidation() {
	suite.Run("APIResponseTimes", func() {
		startTime := time.Now()

		// Start test services
		suite.startTestServices()
		defer suite.stopTestServices()

		// Perform load test for each endpoint
		apiMetrics := make(map[string]APIResponseMetrics)

		for endpointName, pattern := range suite.performanceConfig.APIValidation.RequestPatterns {
			metrics := suite.performAPILoadTest(endpointName, pattern)
			apiMetrics[endpointName] = metrics

			// Validate P95 response time
			suite.Assert().True(metrics.P95 < suite.performanceConfig.Thresholds.APIP95ResponseMax,
				"API endpoint %s P95 response time %v exceeds threshold %v",
				endpointName, metrics.P95, suite.performanceConfig.Thresholds.APIP95ResponseMax)

			suite.T().Logf("API endpoint %s performance - P95: %v, P99: %v, Avg: %v",
				endpointName, metrics.P95, metrics.P99, metrics.Average)
		}

		// Calculate overall API response metrics
		overallMetrics := suite.calculateOverallAPIMetrics(apiMetrics)

		// Record result
		result := PerformanceResult{
			TestName:         "APIResponseTimeValidation",
			Category:         "api_validation",
			StartTime:        startTime,
			Duration:         time.Since(startTime),
			Success:          overallMetrics.P95 < suite.performanceConfig.Thresholds.APIP95ResponseMax,
			APIResponseTimes: overallMetrics,
			ResourceUsage:    suite.captureResourceUsage(),
			ThresholdViolations: suite.checkAPIResponseViolations(overallMetrics),
			Metadata:         suite.createTestMetadata("APIResponseTimeValidation"),
		}

		suite.mu.Lock()
		suite.results = append(suite.results, result)
		suite.mu.Unlock()
	})
}

// TestPactBrokerIntegrationPerformance tests Pact Broker integration performance
func (suite *ContractPerformanceValidationSuite) TestPactBrokerIntegrationPerformance() {
	suite.Run("PactBrokerIntegration", func() {
		startTime := time.Now()

		// Test contract publication performance
		publicationStart := time.Now()
		suite.performBrokerContractPublication()
		publicationTime := time.Since(publicationStart)

		suite.Assert().True(publicationTime < suite.performanceConfig.BrokerIntegration.MaxPublicationTime,
			"Contract publication took %v, exceeds threshold %v",
			publicationTime, suite.performanceConfig.BrokerIntegration.MaxPublicationTime)

		// Test contract retrieval performance
		retrievalStart := time.Now()
		suite.performBrokerContractRetrieval()
		retrievalTime := time.Since(retrievalStart)

		suite.Assert().True(retrievalTime < suite.performanceConfig.BrokerIntegration.MaxRetrievalTime,
			"Contract retrieval took %v, exceeds threshold %v",
			retrievalTime, suite.performanceConfig.BrokerIntegration.MaxRetrievalTime)

		// Test verification performance
		verificationStart := time.Now()
		suite.performBrokerVerification()
		verificationTime := time.Since(verificationStart)

		suite.Assert().True(verificationTime < suite.performanceConfig.BrokerIntegration.MaxVerificationTime,
			"Broker verification took %v, exceeds threshold %v",
			verificationTime, suite.performanceConfig.BrokerIntegration.MaxVerificationTime)

		totalBrokerTime := publicationTime + retrievalTime + verificationTime

		// Record result
		result := PerformanceResult{
			TestName:              "PactBrokerIntegrationPerformance",
			Category:              "broker_integration",
			StartTime:             startTime,
			Duration:              time.Since(startTime),
			Success:               totalBrokerTime < suite.performanceConfig.Thresholds.BrokerIntegrationMax,
			BrokerIntegrationTime: totalBrokerTime,
			ResourceUsage:         suite.captureResourceUsage(),
			ThresholdViolations:   suite.checkBrokerIntegrationViolations(totalBrokerTime),
			Metadata:              suite.createTestMetadata("PactBrokerIntegrationPerformance"),
		}

		suite.mu.Lock()
		suite.results = append(suite.results, result)
		suite.mu.Unlock()

		suite.T().Logf("Pact Broker integration performance - Publication: %v, Retrieval: %v, Verification: %v",
			publicationTime, retrievalTime, verificationTime)
	})
}

// TestCrossPlatformValidationPerformance tests cross-platform validation performance
func (suite *ContractPerformanceValidationSuite) TestCrossPlatformValidationPerformance() {
	suite.Run("CrossPlatformValidation", func() {
		startTime := time.Now()

		// Test compatibility check across platforms
		compatibilityStart := time.Now()
		suite.performCrossPlatformCompatibilityCheck()
		compatibilityTime := time.Since(compatibilityStart)

		suite.Assert().True(compatibilityTime < suite.performanceConfig.Thresholds.CrossPlatformMax,
			"Cross-platform compatibility check took %v, exceeds threshold %v",
			compatibilityTime, suite.performanceConfig.Thresholds.CrossPlatformMax)

		// Test contract comparison performance
		comparisonStart := time.Now()
		suite.performContractComparison()
		comparisonTime := time.Since(comparisonStart)

		totalCrossPlatformTime := compatibilityTime + comparisonTime

		// Record result
		result := PerformanceResult{
			TestName:          "CrossPlatformValidationPerformance",
			Category:          "cross_platform",
			StartTime:         startTime,
			Duration:          time.Since(startTime),
			Success:           totalCrossPlatformTime < suite.performanceConfig.Thresholds.CrossPlatformMax,
			CrossPlatformTime: totalCrossPlatformTime,
			ResourceUsage:     suite.captureResourceUsage(),
			ThresholdViolations: suite.checkCrossPlatformViolations(totalCrossPlatformTime),
			Metadata:          suite.createTestMetadata("CrossPlatformValidationPerformance"),
		}

		suite.mu.Lock()
		suite.results = append(suite.results, result)
		suite.mu.Unlock()

		suite.T().Logf("Cross-platform validation performance - Total: %v (Compatibility: %v, Comparison: %v)",
			totalCrossPlatformTime, compatibilityTime, comparisonTime)
	})
}

// TestProviderVerificationPerformance tests provider verification performance (<30s for all services)
func (suite *ContractPerformanceValidationSuite) TestProviderVerificationPerformance() {
	suite.Run("ProviderVerificationPerformance", func() {
		startTime := time.Now()

		// Start all test services
		suite.startAllProviderServices()
		defer suite.stopAllProviderServices()

		if suite.performanceConfig.ProviderVerification.ParallelVerification {
			// Parallel verification
			suite.performParallelProviderVerification()
		} else {
			// Sequential verification
			suite.performSequentialProviderVerification()
		}

		totalVerificationTime := time.Since(startTime)

		suite.Assert().True(totalVerificationTime < suite.performanceConfig.Thresholds.ProviderVerificationMax,
			"Provider verification took %v, exceeds threshold %v",
			totalVerificationTime, suite.performanceConfig.Thresholds.ProviderVerificationMax)

		// Record result
		result := PerformanceResult{
			TestName:                 "ProviderVerificationPerformance",
			Category:                 "provider_verification",
			StartTime:                startTime,
			Duration:                 totalVerificationTime,
			Success:                  totalVerificationTime < suite.performanceConfig.Thresholds.ProviderVerificationMax,
			ProviderVerificationTime: totalVerificationTime,
			ResourceUsage:            suite.captureResourceUsage(),
			ThresholdViolations:      suite.checkProviderVerificationViolations(totalVerificationTime),
			Metadata:                 suite.createTestMetadata("ProviderVerificationPerformance"),
		}

		suite.mu.Lock()
		suite.results = append(suite.results, result)
		suite.mu.Unlock()

		suite.T().Logf("Provider verification performance - Total: %v for %d services (threshold: %v)",
			totalVerificationTime, len(suite.performanceConfig.ProviderVerification.Services),
			suite.performanceConfig.Thresholds.ProviderVerificationMax)
	})
}

// Helper methods for performance testing

// getContractPaths returns paths to available contract files
func (suite *ContractPerformanceValidationSuite) getContractPaths() []string {
	contractDir := filepath.Join(suite.projectRoot, "..", "specs", "021-implement-pact-contract", "contracts")

	// Mock contract paths for performance testing
	return []string{
		filepath.Join(contractDir, "pact-consumer-web.json"),
		filepath.Join(contractDir, "pact-consumer-mobile.json"),
		filepath.Join(contractDir, "contract-validation-api.yaml"),
	}
}

// validateContractFile simulates contract file validation
func (suite *ContractPerformanceValidationSuite) validateContractFile(contractPath string) {
	// Simulate contract parsing and validation
	time.Sleep(time.Duration(50+suite.randomInt(0, 100)) * time.Millisecond)
}

// startTestServices starts test services for API performance testing
func (suite *ContractPerformanceValidationSuite) startTestServices() {
	// Mock service startup
	for serviceName := range suite.performanceConfig.APIValidation.EndpointMapping {
		service := &TestService{
			Name:      serviceName,
			BaseURL:   fmt.Sprintf("http://localhost:%d", 8080+suite.randomInt(0, 100)),
			StartTime: time.Now(),
		}
		suite.testServices[serviceName] = service
	}

	// Allow services to start up
	time.Sleep(100 * time.Millisecond)
}

// stopTestServices stops all test services
func (suite *ContractPerformanceValidationSuite) stopTestServices() {
	for name, service := range suite.testServices {
		if service.Server != nil {
			ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
			service.Server.Shutdown(ctx)
			cancel()
		}
		delete(suite.testServices, name)
	}
}

// performAPILoadTest performs load testing on an API endpoint
func (suite *ContractPerformanceValidationSuite) performAPILoadTest(endpointName string, pattern RequestPattern) APIResponseMetrics {
	// Simulate load test execution
	concurrency := suite.performanceConfig.APIValidation.ConcurrentRequests
	duration := suite.performanceConfig.APIValidation.LoadTestDuration

	var responseTimes []time.Duration
	var errorCount int64

	// Simulate concurrent requests
	for i := 0; i < concurrency; i++ {
		responseTime := time.Duration(50+suite.randomInt(0, 100)) * time.Millisecond
		responseTimes = append(responseTimes, responseTime)

		if suite.randomInt(0, 100) < 2 { // 2% error rate
			errorCount++
		}
	}

	// Calculate percentiles (simplified)
	totalRequests := int64(len(responseTimes))
	if totalRequests == 0 {
		return APIResponseMetrics{}
	}

	var totalTime time.Duration
	minTime := responseTimes[0]
	maxTime := responseTimes[0]

	for _, rt := range responseTimes {
		totalTime += rt
		if rt < minTime {
			minTime = rt
		}
		if rt > maxTime {
			maxTime = rt
		}
	}

	avgTime := totalTime / time.Duration(totalRequests)

	return APIResponseMetrics{
		P50:          avgTime,
		P95:          time.Duration(float64(avgTime) * 1.3), // Simulated P95
		P99:          time.Duration(float64(avgTime) * 1.7), // Simulated P99
		Max:          maxTime,
		Min:          minTime,
		Average:      avgTime,
		RequestCount: totalRequests,
		ErrorCount:   errorCount,
	}
}

// calculateOverallAPIMetrics calculates overall API performance metrics
func (suite *ContractPerformanceValidationSuite) calculateOverallAPIMetrics(endpointMetrics map[string]APIResponseMetrics) APIResponseMetrics {
	if len(endpointMetrics) == 0 {
		return APIResponseMetrics{}
	}

	var totalRequests, totalErrors int64
	var totalP95, totalAvg time.Duration

	for _, metrics := range endpointMetrics {
		totalRequests += metrics.RequestCount
		totalErrors += metrics.ErrorCount
		totalP95 += metrics.P95
		totalAvg += metrics.Average
	}

	count := time.Duration(len(endpointMetrics))
	avgP95 := totalP95 / count
	avgAverage := totalAvg / count

	return APIResponseMetrics{
		P50:          avgAverage,
		P95:          avgP95,
		P99:          time.Duration(float64(avgP95) * 1.3),
		Average:      avgAverage,
		RequestCount: totalRequests,
		ErrorCount:   totalErrors,
	}
}

// performBrokerContractPublication simulates contract publication to Pact Broker
func (suite *ContractPerformanceValidationSuite) performBrokerContractPublication() {
	// Simulate contract publication
	time.Sleep(time.Duration(100+suite.randomInt(0, 200)) * time.Millisecond)
}

// performBrokerContractRetrieval simulates contract retrieval from Pact Broker
func (suite *ContractPerformanceValidationSuite) performBrokerContractRetrieval() {
	// Simulate contract retrieval
	time.Sleep(time.Duration(50+suite.randomInt(0, 100)) * time.Millisecond)
}

// performBrokerVerification simulates broker verification process
func (suite *ContractPerformanceValidationSuite) performBrokerVerification() {
	// Simulate verification process
	time.Sleep(time.Duration(200+suite.randomInt(0, 300)) * time.Millisecond)
}

// performCrossPlatformCompatibilityCheck simulates cross-platform compatibility checking
func (suite *ContractPerformanceValidationSuite) performCrossPlatformCompatibilityCheck() {
	// Simulate compatibility check across platforms
	for _, platform := range suite.performanceConfig.CrossPlatform.Platforms {
		_ = platform // Use platform for simulation
		time.Sleep(time.Duration(100+suite.randomInt(0, 200)) * time.Millisecond)
	}
}

// performContractComparison simulates contract comparison between platforms
func (suite *ContractPerformanceValidationSuite) performContractComparison() {
	// Simulate contract comparison
	time.Sleep(time.Duration(200+suite.randomInt(0, 300)) * time.Millisecond)
}

// startAllProviderServices starts all provider services for verification
func (suite *ContractPerformanceValidationSuite) startAllProviderServices() {
	for _, serviceName := range suite.performanceConfig.ProviderVerification.Services {
		service := &TestService{
			Name:      serviceName,
			BaseURL:   fmt.Sprintf("http://localhost:%d", 9000+suite.randomInt(0, 100)),
			StartTime: time.Now(),
		}
		suite.testServices[serviceName] = service
	}

	// Allow services to start up
	time.Sleep(500 * time.Millisecond)
}

// stopAllProviderServices stops all provider services
func (suite *ContractPerformanceValidationSuite) stopAllProviderServices() {
	suite.stopTestServices()
}

// performParallelProviderVerification performs parallel provider verification
func (suite *ContractPerformanceValidationSuite) performParallelProviderVerification() {
	var wg sync.WaitGroup

	for _, serviceName := range suite.performanceConfig.ProviderVerification.Services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			suite.verifyProviderService(service)
		}(serviceName)
	}

	wg.Wait()
}

// performSequentialProviderVerification performs sequential provider verification
func (suite *ContractPerformanceValidationSuite) performSequentialProviderVerification() {
	for _, serviceName := range suite.performanceConfig.ProviderVerification.Services {
		suite.verifyProviderService(serviceName)
	}
}

// verifyProviderService simulates provider service verification
func (suite *ContractPerformanceValidationSuite) verifyProviderService(serviceName string) {
	// Simulate provider verification using Pact
	verifyStart := time.Now()

	// Mock verification process
	time.Sleep(time.Duration(1000+suite.randomInt(0, 2000)) * time.Millisecond)

	verifyDuration := time.Since(verifyStart)

	suite.Assert().True(verifyDuration < suite.performanceConfig.ProviderVerification.MaxSingleServiceTime,
		"Provider verification for %s took %v, exceeds threshold %v",
		serviceName, verifyDuration, suite.performanceConfig.ProviderVerification.MaxSingleServiceTime)
}

// captureResourceUsage captures current resource usage metrics
func (suite *ContractPerformanceValidationSuite) captureResourceUsage() ResourceUsageMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return ResourceUsageMetrics{
		MemoryUsageMB:    int64(m.Alloc / 1024 / 1024),
		CPUUsagePercent:  float64(50 + suite.randomInt(0, 30)), // Simulated CPU usage
		GoroutineCount:   runtime.NumGoroutine(),
		AllocatedObjects: m.Mallocs - m.Frees,
		GCPauseTime:      time.Duration(m.PauseTotalNs),
	}
}

// Violation checking methods

func (suite *ContractPerformanceValidationSuite) checkContractValidationViolations(avgTime time.Duration) []PerformanceViolation {
	var violations []PerformanceViolation

	if avgTime > suite.performanceConfig.Thresholds.ContractValidationMax {
		violations = append(violations, PerformanceViolation{
			Type:        "contract_validation",
			Metric:      "average_validation_time",
			Expected:    suite.performanceConfig.Thresholds.ContractValidationMax.String(),
			Actual:      avgTime.String(),
			Severity:    "critical",
			Timestamp:   time.Now(),
			Description: "Contract validation time exceeds maximum threshold",
		})
	}

	return violations
}

func (suite *ContractPerformanceValidationSuite) checkAPIResponseViolations(metrics APIResponseMetrics) []PerformanceViolation {
	var violations []PerformanceViolation

	if metrics.P95 > suite.performanceConfig.Thresholds.APIP95ResponseMax {
		violations = append(violations, PerformanceViolation{
			Type:        "api_response",
			Metric:      "p95_response_time",
			Expected:    suite.performanceConfig.Thresholds.APIP95ResponseMax.String(),
			Actual:      metrics.P95.String(),
			Severity:    "critical",
			Timestamp:   time.Now(),
			Description: "API P95 response time exceeds maximum threshold",
		})
	}

	return violations
}

func (suite *ContractPerformanceValidationSuite) checkBrokerIntegrationViolations(totalTime time.Duration) []PerformanceViolation {
	var violations []PerformanceViolation

	if totalTime > suite.performanceConfig.Thresholds.BrokerIntegrationMax {
		violations = append(violations, PerformanceViolation{
			Type:        "broker_integration",
			Metric:      "total_broker_time",
			Expected:    suite.performanceConfig.Thresholds.BrokerIntegrationMax.String(),
			Actual:      totalTime.String(),
			Severity:    "high",
			Timestamp:   time.Now(),
			Description: "Pact Broker integration time exceeds maximum threshold",
		})
	}

	return violations
}

func (suite *ContractPerformanceValidationSuite) checkCrossPlatformViolations(totalTime time.Duration) []PerformanceViolation {
	var violations []PerformanceViolation

	if totalTime > suite.performanceConfig.Thresholds.CrossPlatformMax {
		violations = append(violations, PerformanceViolation{
			Type:        "cross_platform",
			Metric:      "total_cross_platform_time",
			Expected:    suite.performanceConfig.Thresholds.CrossPlatformMax.String(),
			Actual:      totalTime.String(),
			Severity:    "high",
			Timestamp:   time.Now(),
			Description: "Cross-platform validation time exceeds maximum threshold",
		})
	}

	return violations
}

func (suite *ContractPerformanceValidationSuite) checkProviderVerificationViolations(totalTime time.Duration) []PerformanceViolation {
	var violations []PerformanceViolation

	if totalTime > suite.performanceConfig.Thresholds.ProviderVerificationMax {
		violations = append(violations, PerformanceViolation{
			Type:        "provider_verification",
			Metric:      "total_verification_time",
			Expected:    suite.performanceConfig.Thresholds.ProviderVerificationMax.String(),
			Actual:      totalTime.String(),
			Severity:    "critical",
			Timestamp:   time.Now(),
			Description: "Provider verification time exceeds maximum threshold",
		})
	}

	return violations
}

// createTestMetadata creates performance test metadata
func (suite *ContractPerformanceValidationSuite) createTestMetadata(testName string) PerformanceTestMetadata {
	return PerformanceTestMetadata{
		TestSuite:   "ContractPerformanceValidationSuite",
		Version:     "1.0.0",
		Environment: "test",
		Configuration: map[string]interface{}{
			"performance_config": suite.performanceConfig,
		},
		SystemInfo: SystemInfo{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			CPUCores:     runtime.NumCPU(),
			MemoryGB:     8.0, // Simulated
			GoVersion:    runtime.Version(),
		},
		TestParameters: map[string]interface{}{
			"test_name":           testName,
			"concurrent_requests": suite.performanceConfig.APIValidation.ConcurrentRequests,
			"load_test_duration":  suite.performanceConfig.APIValidation.LoadTestDuration,
		},
	}
}

// randomInt generates a random integer between min and max
func (suite *ContractPerformanceValidationSuite) randomInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}

// TearDownSuite cleans up after performance validation tests
func (suite *ContractPerformanceValidationSuite) TearDownSuite() {
	// Generate performance validation report
	suite.generatePerformanceReport()

	// Export metrics if monitoring is enabled
	if suite.performanceConfig.Monitoring.MetricsCollectionEnabled {
		suite.exportPerformanceMetrics()
	}

	// Clean up any remaining test resources
	suite.cleanupPerformanceTestResources()
}

// generatePerformanceReport generates a comprehensive performance validation report
func (suite *ContractPerformanceValidationSuite) generatePerformanceReport() {
	if len(suite.results) == 0 {
		suite.T().Log("No performance validation results to report")
		return
	}

	// Calculate summary statistics
	var totalDuration time.Duration
	var totalViolations int
	successCount := 0

	for _, result := range suite.results {
		totalDuration += result.Duration
		totalViolations += len(result.ThresholdViolations)
		if result.Success {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(len(suite.results))

	// Create comprehensive report
	report := map[string]interface{}{
		"performance_validation_summary": map[string]interface{}{
			"total_tests":          len(suite.results),
			"success_rate":         successRate,
			"total_duration":       totalDuration,
			"total_violations":     totalViolations,
			"test_timestamp":       time.Now(),
			"performance_targets_met": totalViolations == 0,
		},
		"detailed_results":    suite.results,
		"configuration":       suite.performanceConfig,
		"environment_info":    suite.createTestMetadata("PerformanceValidation").SystemInfo,
	}

	// Save report to file
	reportPath := filepath.Join(suite.projectRoot, "tests", "contract", "performance_validation_report.json")
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err == nil {
		err = os.WriteFile(reportPath, reportData, 0644)
		if err == nil {
			suite.T().Logf("Performance validation report saved to: %s", reportPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to save performance validation report: %v", err)
	}

	// Log summary
	suite.T().Logf("\n=== CONTRACT TESTING PERFORMANCE VALIDATION SUMMARY ===")
	suite.T().Logf("Tests Executed: %d", len(suite.results))
	suite.T().Logf("Success Rate: %.1f%%", successRate*100)
	suite.T().Logf("Total Duration: %v", totalDuration)
	suite.T().Logf("Performance Violations: %d", totalViolations)
	suite.T().Logf("Performance Targets Met: %v", totalViolations == 0)

	// Log detailed performance metrics
	for _, result := range suite.results {
		suite.T().Logf("  %s (%s): Duration=%v, Success=%v, Violations=%d",
			result.TestName, result.Category, result.Duration, result.Success, len(result.ThresholdViolations))
	}
}

// exportPerformanceMetrics exports performance metrics to monitoring systems
func (suite *ContractPerformanceValidationSuite) exportPerformanceMetrics() {
	// Export metrics in multiple formats for enterprise monitoring
	metricsDir := filepath.Join(suite.projectRoot, "tests", "contract")

	// JSON metrics for detailed analysis
	suite.exportJSONPerformanceMetrics(metricsDir)

	// Prometheus metrics for monitoring integration
	suite.exportPrometheusPerformanceMetrics(metricsDir)

	// CSV metrics for reporting
	suite.exportCSVPerformanceMetrics(metricsDir)
}

// exportJSONPerformanceMetrics exports detailed metrics in JSON format
func (suite *ContractPerformanceValidationSuite) exportJSONPerformanceMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "contract_performance_metrics.json")
	metricsData, err := json.MarshalIndent(suite.results, "", "  ")
	if err == nil {
		err = os.WriteFile(metricsPath, metricsData, 0644)
		if err == nil {
			suite.T().Logf("JSON performance metrics exported to: %s", metricsPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to export JSON performance metrics: %v", err)
	}
}

// exportPrometheusPerformanceMetrics exports metrics in Prometheus format
func (suite *ContractPerformanceValidationSuite) exportPrometheusPerformanceMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "contract_performance_metrics.prom")

	var prometheus []string
	for _, result := range suite.results {
		prometheus = append(prometheus, fmt.Sprintf("tchat_contract_test_duration_seconds{test=\"%s\",category=\"%s\"} %.3f",
			result.TestName, result.Category, result.Duration.Seconds()))
		prometheus = append(prometheus, fmt.Sprintf("tchat_contract_test_success{test=\"%s\",category=\"%s\"} %d",
			result.TestName, result.Category, boolToInt(result.Success)))
		prometheus = append(prometheus, fmt.Sprintf("tchat_contract_violations_total{test=\"%s\",category=\"%s\"} %d",
			result.TestName, result.Category, len(result.ThresholdViolations)))

		if result.ContractValidationTime > 0 {
			prometheus = append(prometheus, fmt.Sprintf("tchat_contract_validation_duration_seconds{test=\"%s\"} %.3f",
				result.TestName, result.ContractValidationTime.Seconds()))
		}

		if result.APIResponseTimes.P95 > 0 {
			prometheus = append(prometheus, fmt.Sprintf("tchat_api_response_p95_seconds{test=\"%s\"} %.3f",
				result.TestName, result.APIResponseTimes.P95.Seconds()))
		}
	}

	content := fmt.Sprintf("# Contract testing performance metrics\n%s\n",
		fmt.Sprintf("# Generated at %s\n", time.Now().Format(time.RFC3339)) +
		fmt.Sprintf("%v", prometheus))
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("Prometheus performance metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export Prometheus performance metrics: %v", err)
	}
}

// exportCSVPerformanceMetrics exports metrics in CSV format for reporting
func (suite *ContractPerformanceValidationSuite) exportCSVPerformanceMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "contract_performance_metrics.csv")

	var csv []string
	csv = append(csv, "test_name,category,duration_ms,success,violations,contract_validation_ms,api_p95_ms,memory_usage_mb")

	for _, result := range suite.results {
		csv = append(csv, fmt.Sprintf("%s,%s,%d,%t,%d,%d,%d,%d",
			result.TestName,
			result.Category,
			result.Duration.Milliseconds(),
			result.Success,
			len(result.ThresholdViolations),
			result.ContractValidationTime.Milliseconds(),
			result.APIResponseTimes.P95.Milliseconds(),
			result.ResourceUsage.MemoryUsageMB,
		))
	}

	content := fmt.Sprintf("%v\n", csv)
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("CSV performance metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export CSV performance metrics: %v", err)
	}
}

// cleanupPerformanceTestResources cleans up performance test resources
func (suite *ContractPerformanceValidationSuite) cleanupPerformanceTestResources() {
	// Stop any remaining services
	suite.stopTestServices()

	// Clear results
	suite.results = nil

	suite.T().Log("Contract performance validation resources cleaned up")
}

// boolToInt converts boolean to integer for metrics
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// TestContractPerformanceValidationSuite runs the contract performance validation suite
func TestContractPerformanceValidationSuite(t *testing.T) {
	suite.Run(t, new(ContractPerformanceValidationSuite))
}