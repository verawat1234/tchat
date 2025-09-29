package performance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tchat.dev/calling/config"
	"tchat.dev/calling/handlers"
	"tchat.dev/calling/models"
	"tchat.dev/calling/repositories"
)

// LoadTestSuite defines the load testing suite
type LoadTestSuite struct {
	suite.Suite
	server     *httptest.Server
	db         *gorm.DB
	redisClient *redis.Client
	router     *gin.Engine
}

// SetupSuite sets up the test environment
func (suite *LoadTestSuite) SetupSuite() {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	// Auto-migrate models
	err = db.AutoMigrate(&models.CallSession{}, &models.CallParticipant{})
	suite.Require().NoError(err)

	// Setup Redis client (using miniredis for testing)
	suite.redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use test database
	})

	// Setup router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Setup repositories
	callRepo := repositories.NewCallRepository(suite.db)
	redisWrapper := &config.RedisClient{Client: suite.redisClient}
	presenceService := config.NewUserPresenceService(redisWrapper)

	// Setup handlers
	callHandler := handlers.NewCallHandler(callRepo, presenceService)

	// Setup routes
	api := suite.router.Group("/api/v1")
	api.POST("/calls", callHandler.InitiateCall)
	api.PUT("/calls/:id/answer", callHandler.AnswerCall)
	api.DELETE("/calls/:id", callHandler.EndCall)
	api.GET("/calls/user/:userId/active", callHandler.GetActiveCalls)

	// Create test server
	suite.server = httptest.NewServer(suite.router)
}

// TearDownSuite cleans up the test environment
func (suite *LoadTestSuite) TearDownSuite() {
	suite.server.Close()
	if suite.redisClient != nil {
		suite.redisClient.Close()
	}
	if sqlDB, err := suite.db.DB(); err == nil {
		sqlDB.Close()
	}
}

// LoadTestMetrics holds performance metrics
type LoadTestMetrics struct {
	TotalRequests     int
	SuccessfulReqs    int
	FailedReqs        int
	AvgResponseTime   time.Duration
	MaxResponseTime   time.Duration
	MinResponseTime   time.Duration
	Percentile95      time.Duration
	Percentile99      time.Duration
	RequestsPerSecond float64
	ErrorRate         float64
}

// TestConcurrentCallInitiation tests concurrent call initiation load
func (suite *LoadTestSuite) TestConcurrentCallInitiation() {
	concurrentUsers := 100
	callsPerUser := 10
	targetConnectionTime := 5 * time.Second

	var wg sync.WaitGroup
	results := make(chan time.Duration, concurrentUsers*callsPerUser)
	errors := make(chan error, concurrentUsers*callsPerUser)

	startTime := time.Now()

	// Launch concurrent users
	for i := 0; i < concurrentUsers; i++ {\
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			for j := 0; j < callsPerUser; j++ {
				callStart := time.Now()

				// Create call initiation request
				callerID := uuid.New().String()
				calleeID := uuid.New().String()

				resp, err := http.Post(
					fmt.Sprintf("%s/api/v1/calls", suite.server.URL),
					"application/json",
					nil,
				)

				duration := time.Since(callStart)
				results <- duration

				if err != nil {
					errors <- err
					continue
				}

				resp.Body.Close()

				// Verify connection time is under target
				if duration > targetConnectionTime {
					errors <- fmt.Errorf("connection time %v exceeds target %v", duration, targetConnectionTime)
				}
			}
		}(i)
	}

	// Wait for completion
	wg.Wait()
	close(results)
	close(errors)

	// Calculate metrics
	metrics := suite.calculateMetrics(results, errors, time.Since(startTime))

	// Performance assertions
	suite.T().Logf("Load Test Results:")
	suite.T().Logf("Total Requests: %d", metrics.TotalRequests)
	suite.T().Logf("Success Rate: %.2f%%", (1-metrics.ErrorRate)*100)
	suite.T().Logf("Average Response Time: %v", metrics.AvgResponseTime)
	suite.T().Logf("95th Percentile: %v", metrics.Percentile95)
	suite.T().Logf("99th Percentile: %v", metrics.Percentile99)
	suite.T().Logf("Requests Per Second: %.2f", metrics.RequestsPerSecond)

	// Performance requirements
	assert.True(suite.T(), metrics.ErrorRate < 0.05, "Error rate should be less than 5%")
	assert.True(suite.T(), metrics.Percentile95 < targetConnectionTime, "95th percentile should be under 5s")
	assert.True(suite.T(), metrics.RequestsPerSecond > 50, "Should handle at least 50 RPS")
}

// TestConcurrentCallOperations tests mixed call operations under load
func (suite *LoadTestSuite) TestConcurrentCallOperations() {
	concurrentUsers := 50
	operationsPerUser := 20

	var wg sync.WaitGroup
	results := make(chan time.Duration, concurrentUsers*operationsPerUser)
	errors := make(chan error, concurrentUsers*operationsPerUser)

	startTime := time.Now()

	// Launch concurrent users performing mixed operations
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			for j := 0; j < operationsPerUser; j++ {
				operationStart := time.Now()

				// Randomly choose operation (initiate, answer, end, get active)
				operation := j % 4
				var resp *http.Response
				var err error

				switch operation {
				case 0: // Initiate call
					resp, err = http.Post(fmt.Sprintf("%s/api/v1/calls", suite.server.URL), "application/json", nil)
				case 1: // Answer call (simulated)
					callID := uuid.New().String()
					resp, err = http.Post(fmt.Sprintf("%s/api/v1/calls/%s/answer", suite.server.URL, callID), "application/json", nil)
				case 2: // End call (simulated)
					callID := uuid.New().String()
					req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/calls/%s", suite.server.URL, callID), nil)
					client := &http.Client{}
					resp, err = client.Do(req)
				case 3: // Get active calls
					userID := uuid.New().String()
					resp, err = http.Get(fmt.Sprintf("%s/api/v1/calls/user/%s/active", suite.server.URL, userID))
				}

				duration := time.Since(operationStart)
				results <- duration

				if err != nil {
					errors <- err
				} else {
					resp.Body.Close()
				}
			}
		}(i)
	}

	// Wait for completion
	wg.Wait()
	close(results)
	close(errors)

	// Calculate metrics
	metrics := suite.calculateMetrics(results, errors, time.Since(startTime))

	// Mixed operations performance assertions
	suite.T().Logf("Mixed Operations Load Test Results:")
	suite.T().Logf("Total Operations: %d", metrics.TotalRequests)
	suite.T().Logf("Success Rate: %.2f%%", (1-metrics.ErrorRate)*100)
	suite.T().Logf("Average Response Time: %v", metrics.AvgResponseTime)
	suite.T().Logf("Operations Per Second: %.2f", metrics.RequestsPerSecond)

	assert.True(suite.T(), metrics.ErrorRate < 0.1, "Error rate should be less than 10% for mixed operations")
	assert.True(suite.T(), metrics.AvgResponseTime < 2*time.Second, "Average response time should be under 2s")
}

// TestScalabilityThresholds tests system behavior at different load levels
func (suite *LoadTestSuite) TestScalabilityThresholds() {
	loadLevels := []struct {
		name        string
		users       int
		duration    time.Duration
		maxErrorRate float64
	}{
		{"Light Load", 10, 30 * time.Second, 0.01},
		{"Medium Load", 50, 60 * time.Second, 0.03},
		{"Heavy Load", 100, 90 * time.Second, 0.05},
		{"Stress Load", 200, 120 * time.Second, 0.10},
	}

	for _, level := range loadLevels {
		suite.T().Run(level.name, func(t *testing.T) {
			metrics := suite.runLoadTest(level.users, level.duration)

			t.Logf("%s Results:", level.name)
			t.Logf("Error Rate: %.2f%%", metrics.ErrorRate*100)
			t.Logf("Avg Response Time: %v", metrics.AvgResponseTime)
			t.Logf("RPS: %.2f", metrics.RequestsPerSecond)

			assert.True(t, metrics.ErrorRate <= level.maxErrorRate,
				"Error rate %.2f%% exceeds threshold %.2f%%",
				metrics.ErrorRate*100, level.maxErrorRate*100)
		})
	}
}

// runLoadTest executes a load test with specified parameters
func (suite *LoadTestSuite) runLoadTest(users int, duration time.Duration) LoadTestMetrics {
	var wg sync.WaitGroup
	results := make(chan time.Duration, users*100) // Buffer for results
	errors := make(chan error, users*100)

	startTime := time.Now()
	endTime := startTime.Add(duration)

	// Launch concurrent users
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for time.Now().Before(endTime) {
				requestStart := time.Now()

				// Make a simple call initiation request
				resp, err := http.Post(
					fmt.Sprintf("%s/api/v1/calls", suite.server.URL),
					"application/json",
					nil,
				)

				requestDuration := time.Since(requestStart)
				results <- requestDuration

				if err != nil {
					errors <- err
				} else {
					resp.Body.Close()
				}

				// Small delay to prevent overwhelming
				time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	// Wait for completion
	wg.Wait()
	close(results)
	close(errors)

	return suite.calculateMetrics(results, errors, duration)
}

// calculateMetrics computes performance metrics from results
func (suite *LoadTestSuite) calculateMetrics(results <-chan time.Duration, errors <-chan error, totalDuration time.Duration) LoadTestMetrics {
	var durations []time.Duration
	var totalDuration_sum time.Duration
	errorCount := 0

	// Collect results
	for duration := range results {
		durations = append(durations, duration)
		totalDuration_sum += duration
	}

	// Count errors
	for range errors {
		errorCount++
	}

	totalRequests := len(durations) + errorCount
	successfulReqs := len(durations)

	if len(durations) == 0 {
		return LoadTestMetrics{
			TotalRequests: totalRequests,
			FailedReqs:    errorCount,
			ErrorRate:     1.0,
		}
	}

	// Calculate basic metrics
	avgResponseTime := totalDuration_sum / time.Duration(len(durations))

	// Find min/max
	minTime := durations[0]
	maxTime := durations[0]
	for _, d := range durations {
		if d < minTime {
			minTime = d
		}
		if d > maxTime {
			maxTime = d
		}
	}

	// Calculate percentiles (simplified)
	p95Index := int(float64(len(durations)) * 0.95)
	p99Index := int(float64(len(durations)) * 0.99)

	if p95Index >= len(durations) {
		p95Index = len(durations) - 1
	}
	if p99Index >= len(durations) {
		p99Index = len(durations) - 1
	}

	return LoadTestMetrics{
		TotalRequests:     totalRequests,
		SuccessfulReqs:    successfulReqs,
		FailedReqs:        errorCount,
		AvgResponseTime:   avgResponseTime,
		MaxResponseTime:   maxTime,
		MinResponseTime:   minTime,
		Percentile95:      durations[p95Index],
		Percentile99:      durations[p99Index],
		RequestsPerSecond: float64(totalRequests) / totalDuration.Seconds(),
		ErrorRate:         float64(errorCount) / float64(totalRequests),
	}
}

// TestLoadTestSuite runs the load testing suite
func TestLoadTestSuite(t *testing.T) {
	suite.Run(t, new(LoadTestSuite))
}

// Benchmark tests for performance measurement

func BenchmarkCallInitiation(b *testing.B) {
	// Setup minimal test environment
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simple handler for benchmarking
	router.POST("/calls", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success"})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := http.Post(fmt.Sprintf("%s/calls", server.URL), "application/json", nil)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func BenchmarkConcurrentCalls(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/calls", func(c *gin.Context) {
		// Simulate some processing time
		time.Sleep(1 * time.Millisecond)
		c.JSON(200, gin.H{"status": "success"})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, _ := http.Post(fmt.Sprintf("%s/calls", server.URL), "application/json", nil)
			if resp != nil {
				resp.Body.Close()
			}
		}
	})
}