package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadTestConfig defines configuration for load testing
type LoadTestConfig struct {
	BaseURL           string
	Concurrency       int
	Duration          time.Duration
	RampUpTime        time.Duration
	RequestsPerSecond int
	Timeout           time.Duration
	ThresholdP95      time.Duration
	ThresholdP99      time.Duration
	MaxErrorRate      float64
}

// LoadTestResult contains the results of a load test
type LoadTestResult struct {
	TotalRequests     int64
	SuccessfulReqs    int64
	FailedRequests    int64
	ErrorRate         float64
	AverageLatency    time.Duration
	P50Latency        time.Duration
	P95Latency        time.Duration
	P99Latency        time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	TotalDuration     time.Duration
	RequestsPerSecond float64
	Errors            map[string]int64
	StatusCodes       map[int]int64
}

// RequestMetric captures metrics for a single request
type RequestMetric struct {
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	StatusCode int
	Error      error
	Endpoint   string
	Method     string
}

// LoadTester orchestrates load testing
type LoadTester struct {
	config  LoadTestConfig
	client  *http.Client
	metrics []RequestMetric
	mu      sync.Mutex
}

// NewLoadTester creates a new load tester instance
func NewLoadTester(config LoadTestConfig) *LoadTester {
	return &LoadTester{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		metrics: make([]RequestMetric, 0),
	}
}

// RunLoadTest executes the load test
func (lt *LoadTester) RunLoadTest(ctx context.Context, scenarios []LoadTestScenario) (*LoadTestResult, error) {
	startTime := time.Now()

	// Reset metrics
	lt.mu.Lock()
	lt.metrics = make([]RequestMetric, 0)
	lt.mu.Unlock()

	// Create channels for coordination
	workerChan := make(chan LoadTestScenario, lt.config.Concurrency*2)
	resultChan := make(chan RequestMetric, lt.config.Concurrency*2)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < lt.config.Concurrency; i++ {
		wg.Add(1)
		go lt.worker(ctx, &wg, workerChan, resultChan)
	}

	// Start result collector
	go lt.collectResults(resultChan)

	// Schedule requests based on RPS
	go lt.scheduleRequests(ctx, scenarios, workerChan)

	// Wait for test duration
	testTimer := time.NewTimer(lt.config.Duration)
	select {
	case <-testTimer.C:
		// Test duration completed
	case <-ctx.Done():
		// Context cancelled
	}

	// Stop scheduling new requests
	close(workerChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Calculate and return results
	return lt.calculateResults(time.Since(startTime)), nil
}

// worker processes load test requests
func (lt *LoadTester) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan LoadTestScenario, results chan<- RequestMetric) {
	defer wg.Done()

	for {
		select {
		case scenario, ok := <-jobs:
			if !ok {
				return
			}
			metric := lt.executeRequest(scenario)
			results <- metric
		case <-ctx.Done():
			return
		}
	}
}

// scheduleRequests schedules requests according to RPS configuration
func (lt *LoadTester) scheduleRequests(ctx context.Context, scenarios []LoadTestScenario, workerChan chan<- LoadTestScenario) {
	ticker := time.NewTicker(time.Second / time.Duration(lt.config.RequestsPerSecond))
	defer ticker.Stop()

	rampUpInterval := lt.config.RampUpTime / time.Duration(lt.config.Concurrency)
	rampUpTicker := time.NewTicker(rampUpInterval)
	defer rampUpTicker.Stop()

	activeWorkers := 1
	requestCount := int64(0)

	for {
		select {
		case <-ticker.C:
			if activeWorkers >= lt.config.Concurrency {
				// All workers active, send requests
				scenario := lt.selectScenario(scenarios, requestCount)
				select {
				case workerChan <- scenario:
					atomic.AddInt64(&requestCount, 1)
				default:
					// Channel full, skip this request
				}
			}
		case <-rampUpTicker.C:
			if activeWorkers < lt.config.Concurrency {
				activeWorkers++
			}
		case <-ctx.Done():
			return
		}
	}
}

// selectScenario selects a scenario based on request count and weights
func (lt *LoadTester) selectScenario(scenarios []LoadTestScenario, requestCount int64) LoadTestScenario {
	if len(scenarios) == 1 {
		return scenarios[0]
	}

	// Calculate cumulative weights
	totalWeight := 0
	for _, scenario := range scenarios {
		totalWeight += scenario.Weight
	}

	// Select based on weight
	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, scenario := range scenarios {
		cumulative += scenario.Weight
		if r < cumulative {
			return scenario
		}
	}

	return scenarios[0] // Fallback
}

// executeRequest executes a single request and captures metrics
func (lt *LoadTester) executeRequest(scenario LoadTestScenario) RequestMetric {
	metric := RequestMetric{
		StartTime: time.Now(),
		Endpoint:  scenario.Endpoint,
		Method:    scenario.Method,
	}

	url := lt.config.BaseURL + scenario.Endpoint

	var req *http.Request
	var err error

	if scenario.Body != nil {
		bodyBytes, _ := json.Marshal(scenario.Body)
		req, err = http.NewRequest(scenario.Method, url, bytes.NewBuffer(bodyBytes))
	} else {
		req, err = http.NewRequest(scenario.Method, url, nil)
	}

	if err != nil {
		metric.EndTime = time.Now()
		metric.Duration = metric.EndTime.Sub(metric.StartTime)
		metric.Error = err
		return metric
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if scenario.Headers != nil {
		for key, value := range scenario.Headers {
			req.Header.Set(key, value)
		}
	}

	// Execute request
	resp, err := lt.client.Do(req)
	metric.EndTime = time.Now()
	metric.Duration = metric.EndTime.Sub(metric.StartTime)

	if err != nil {
		metric.Error = err
		return metric
	}
	defer resp.Body.Close()

	metric.StatusCode = resp.StatusCode

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		metric.Error = fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return metric
}

// collectResults collects metrics from workers
func (lt *LoadTester) collectResults(results <-chan RequestMetric) {
	for metric := range results {
		lt.mu.Lock()
		lt.metrics = append(lt.metrics, metric)
		lt.mu.Unlock()
	}
}

// calculateResults calculates test results from collected metrics
func (lt *LoadTester) calculateResults(totalDuration time.Duration) *LoadTestResult {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	if len(lt.metrics) == 0 {
		return &LoadTestResult{}
	}

	result := &LoadTestResult{
		TotalRequests: int64(len(lt.metrics)),
		TotalDuration: totalDuration,
		Errors:        make(map[string]int64),
		StatusCodes:   make(map[int]int64),
	}

	var totalLatency time.Duration
	var latencies []time.Duration

	for _, metric := range lt.metrics {
		if metric.Error != nil {
			result.FailedRequests++
			result.Errors[metric.Error.Error()]++
		} else {
			result.SuccessfulReqs++
		}

		result.StatusCodes[metric.StatusCode]++
		latencies = append(latencies, metric.Duration)
		totalLatency += metric.Duration
	}

	// Calculate error rate
	result.ErrorRate = float64(result.FailedRequests) / float64(result.TotalRequests) * 100

	// Calculate latency statistics
	if len(latencies) > 0 {
		result.AverageLatency = totalLatency / time.Duration(len(latencies))

		// Sort latencies for percentile calculations
		for i := 0; i < len(latencies); i++ {
			for j := i + 1; j < len(latencies); j++ {
				if latencies[i] > latencies[j] {
					latencies[i], latencies[j] = latencies[j], latencies[i]
				}
			}
		}

		result.MinLatency = latencies[0]
		result.MaxLatency = latencies[len(latencies)-1]
		result.P50Latency = latencies[len(latencies)*50/100]
		result.P95Latency = latencies[len(latencies)*95/100]
		result.P99Latency = latencies[len(latencies)*99/100]
	}

	// Calculate RPS
	result.RequestsPerSecond = float64(result.TotalRequests) / totalDuration.Seconds()

	return result
}

// LoadTestScenario defines a test scenario
type LoadTestScenario struct {
	Name     string
	Method   string
	Endpoint string
	Headers  map[string]string
	Body     interface{}
	Weight   int // Weight for scenario selection
}

// Commerce API Load Test Scenarios
func GetCommerceLoadTestScenarios() []LoadTestScenario {
	userToken := "Bearer test-token-" + uuid.New().String()

	return []LoadTestScenario{
		{
			Name:     "GetProducts",
			Method:   "GET",
			Endpoint: "/api/v1/commerce/products?page=1&limit=20",
			Headers:  map[string]string{"Authorization": userToken},
			Weight:   30, // 30% of requests
		},
		{
			Name:     "GetCart",
			Method:   "GET",
			Endpoint: "/api/v1/commerce/cart",
			Headers:  map[string]string{"Authorization": userToken},
			Weight:   20, // 20% of requests
		},
		{
			Name:     "GetCategories",
			Method:   "GET",
			Endpoint: "/api/v1/commerce/categories",
			Headers:  map[string]string{"Authorization": userToken},
			Weight:   15, // 15% of requests
		},
		{
			Name:     "AddToCart",
			Method:   "POST",
			Endpoint: "/api/v1/commerce/cart/items",
			Headers:  map[string]string{"Authorization": userToken},
			Body: map[string]interface{}{
				"productId": uuid.New().String(),
				"quantity":  rand.Intn(5) + 1,
				"unitPrice": 99.99 + float64(rand.Intn(900)),
			},
			Weight: 15, // 15% of requests
		},
		{
			Name:     "UpdateCartItem",
			Method:   "PUT",
			Endpoint: "/api/v1/commerce/cart/items/" + uuid.New().String(),
			Headers:  map[string]string{"Authorization": userToken},
			Body: map[string]interface{}{
				"quantity": rand.Intn(5) + 1,
			},
			Weight: 10, // 10% of requests
		},
		{
			Name:     "SearchProducts",
			Method:   "GET",
			Endpoint: fmt.Sprintf("/api/v1/commerce/products/search?query=%s", []string{"phone", "laptop", "tablet", "headphones"}[rand.Intn(4)]),
			Headers:  map[string]string{"Authorization": userToken},
			Weight:   10, // 10% of requests
		},
	}
}

// TestCommerceAPILoad runs load tests for commerce API
func TestCommerceAPILoad(t *testing.T) {
	config := LoadTestConfig{
		BaseURL:           "http://localhost:8080",
		Concurrency:       50,
		Duration:          60 * time.Second,
		RampUpTime:        10 * time.Second,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		ThresholdP95:      200 * time.Millisecond,
		ThresholdP99:      500 * time.Millisecond,
		MaxErrorRate:      5.0, // 5%
	}

	loadTester := NewLoadTester(config)
	scenarios := GetCommerceLoadTestScenarios()

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+30*time.Second)
	defer cancel()

	t.Logf("Starting load test with %d concurrent users for %v", config.Concurrency, config.Duration)
	t.Logf("Target RPS: %d", config.RequestsPerSecond)

	result, err := loadTester.RunLoadTest(ctx, scenarios)
	require.NoError(t, err)

	// Log results
	t.Logf("Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful Requests: %d", result.SuccessfulReqs)
	t.Logf("Failed Requests: %d", result.FailedRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("Average Latency: %v", result.AverageLatency)
	t.Logf("P50 Latency: %v", result.P50Latency)
	t.Logf("P95 Latency: %v", result.P95Latency)
	t.Logf("P99 Latency: %v", result.P99Latency)
	t.Logf("Min Latency: %v", result.MinLatency)
	t.Logf("Max Latency: %v", result.MaxLatency)
	t.Logf("Actual RPS: %.2f", result.RequestsPerSecond)
	t.Logf("Test Duration: %v", result.TotalDuration)

	if len(result.Errors) > 0 {
		t.Logf("Errors:")
		for error, count := range result.Errors {
			t.Logf("  %s: %d", error, count)
		}
	}

	t.Logf("Status Codes:")
	for code, count := range result.StatusCodes {
		t.Logf("  %d: %d", code, count)
	}

	// Performance assertions
	assert.LessOrEqual(t, result.ErrorRate, config.MaxErrorRate, "Error rate should be within threshold")
	assert.LessOrEqual(t, result.P95Latency, config.ThresholdP95, "P95 latency should be within threshold")
	assert.LessOrEqual(t, result.P99Latency, config.ThresholdP99, "P99 latency should be within threshold")
	assert.GreaterOrEqual(t, result.RequestsPerSecond, float64(config.RequestsPerSecond)*0.8, "Should achieve at least 80% of target RPS")
}

// TestCartWorkflowLoad tests cart workflow under load
func TestCartWorkflowLoad(t *testing.T) {
	config := LoadTestConfig{
		BaseURL:           "http://localhost:8080",
		Concurrency:       20,
		Duration:          30 * time.Second,
		RampUpTime:        5 * time.Second,
		RequestsPerSecond: 50,
		Timeout:           3 * time.Second,
		ThresholdP95:      150 * time.Millisecond,
		ThresholdP99:      300 * time.Millisecond,
		MaxErrorRate:      3.0,
	}

	loadTester := NewLoadTester(config)

	// Cart-focused scenarios
	scenarios := []LoadTestScenario{
		{
			Name:     "GetCart",
			Method:   "GET",
			Endpoint: "/api/v1/commerce/cart",
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Weight:   40,
		},
		{
			Name:     "AddToCart",
			Method:   "POST",
			Endpoint: "/api/v1/commerce/cart/items",
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Body: map[string]interface{}{
				"productId": "product-load-test",
				"quantity":  1,
				"unitPrice": 99.99,
			},
			Weight: 30,
		},
		{
			Name:     "UpdateCartItem",
			Method:   "PUT",
			Endpoint: "/api/v1/commerce/cart/items/item-load-test",
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Body: map[string]interface{}{
				"quantity": 2,
			},
			Weight: 20,
		},
		{
			Name:     "RemoveFromCart",
			Method:   "DELETE",
			Endpoint: "/api/v1/commerce/cart/items/item-load-test",
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Weight:   10,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+15*time.Second)
	defer cancel()

	result, err := loadTester.RunLoadTest(ctx, scenarios)
	require.NoError(t, err)

	t.Logf("Cart Workflow Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("P95 Latency: %v", result.P95Latency)
	t.Logf("P99 Latency: %v", result.P99Latency)
	t.Logf("RPS: %.2f", result.RequestsPerSecond)

	assert.LessOrEqual(t, result.ErrorRate, config.MaxErrorRate)
	assert.LessOrEqual(t, result.P95Latency, config.ThresholdP95)
	assert.LessOrEqual(t, result.P99Latency, config.ThresholdP99)
}

// TestProductSearchLoad tests product search under load
func TestProductSearchLoad(t *testing.T) {
	config := LoadTestConfig{
		BaseURL:           "http://localhost:8080",
		Concurrency:       30,
		Duration:          45 * time.Second,
		RampUpTime:        5 * time.Second,
		RequestsPerSecond: 80,
		Timeout:           5 * time.Second,
		ThresholdP95:      250 * time.Millisecond,
		ThresholdP99:      500 * time.Millisecond,
		MaxErrorRate:      2.0,
	}

	searchTerms := []string{
		"smartphone", "laptop", "tablet", "headphones", "camera",
		"watch", "keyboard", "mouse", "monitor", "speaker",
		"charger", "case", "cable", "adapter", "battery",
	}

	scenarios := make([]LoadTestScenario, len(searchTerms))
	for i, term := range searchTerms {
		scenarios[i] = LoadTestScenario{
			Name:     fmt.Sprintf("SearchProducts_%s", term),
			Method:   "GET",
			Endpoint: fmt.Sprintf("/api/v1/commerce/products/search?query=%s", term),
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Weight:   1,
		}
	}

	loadTester := NewLoadTester(config)
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+15*time.Second)
	defer cancel()

	result, err := loadTester.RunLoadTest(ctx, scenarios)
	require.NoError(t, err)

	t.Logf("Product Search Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("P95 Latency: %v", result.P95Latency)
	t.Logf("P99 Latency: %v", result.P99Latency)
	t.Logf("RPS: %.2f", result.RequestsPerSecond)

	assert.LessOrEqual(t, result.ErrorRate, config.MaxErrorRate)
	assert.LessOrEqual(t, result.P95Latency, config.ThresholdP95)
	assert.LessOrEqual(t, result.P99Latency, config.ThresholdP99)
}

// TestConcurrentCartUpdates tests concurrent cart updates
func TestConcurrentCartUpdates(t *testing.T) {
	config := LoadTestConfig{
		BaseURL:           "http://localhost:8080",
		Concurrency:       100,
		Duration:          30 * time.Second,
		RampUpTime:        2 * time.Second,
		RequestsPerSecond: 200,
		Timeout:           3 * time.Second,
		ThresholdP95:      100 * time.Millisecond,
		ThresholdP99:      200 * time.Millisecond,
		MaxErrorRate:      10.0, // Higher tolerance for concurrent updates
	}

	// All scenarios target the same cart item to test concurrent updates
	scenarios := []LoadTestScenario{
		{
			Name:     "ConcurrentCartUpdate",
			Method:   "PUT",
			Endpoint: "/api/v1/commerce/cart/items/concurrent-test-item",
			Headers:  map[string]string{"Authorization": "Bearer test-token"},
			Body: map[string]interface{}{
				"quantity": rand.Intn(10) + 1,
			},
			Weight: 100,
		},
	}

	loadTester := NewLoadTester(config)
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+10*time.Second)
	defer cancel()

	result, err := loadTester.RunLoadTest(ctx, scenarios)
	require.NoError(t, err)

	t.Logf("Concurrent Cart Updates Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("P95 Latency: %v", result.P95Latency)
	t.Logf("P99 Latency: %v", result.P99Latency)
	t.Logf("RPS: %.2f", result.RequestsPerSecond)

	// Less strict assertions for concurrent updates
	assert.LessOrEqual(t, result.ErrorRate, config.MaxErrorRate)
	assert.LessOrEqual(t, result.P95Latency, config.ThresholdP95)
}

// TestMemoryLeakDetection tests for potential memory leaks under sustained load
func TestMemoryLeakDetection(t *testing.T) {
	// This test would run longer to detect memory leaks
	if testing.Short() {
		t.Skip("Skipping memory leak test in short mode")
	}

	config := LoadTestConfig{
		BaseURL:           "http://localhost:8080",
		Concurrency:       10,
		Duration:          5 * time.Minute, // Longer duration
		RampUpTime:        10 * time.Second,
		RequestsPerSecond: 20,
		Timeout:           5 * time.Second,
		ThresholdP95:      300 * time.Millisecond,
		ThresholdP99:      600 * time.Millisecond,
		MaxErrorRate:      5.0,
	}

	scenarios := GetCommerceLoadTestScenarios()
	loadTester := NewLoadTester(config)

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+30*time.Second)
	defer cancel()

	t.Logf("Starting memory leak detection test (5 minutes)")

	result, err := loadTester.RunLoadTest(ctx, scenarios)
	require.NoError(t, err)

	t.Logf("Memory Leak Detection Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("Average Latency: %v", result.AverageLatency)
	t.Logf("P95 Latency: %v", result.P95Latency)
	t.Logf("P99 Latency: %v", result.P99Latency)

	// Check for degrading performance over time (potential memory leak indicator)
	// This would require more sophisticated analysis of latency trends
	assert.LessOrEqual(t, result.ErrorRate, config.MaxErrorRate)
	assert.LessOrEqual(t, result.P95Latency, config.ThresholdP95)
	assert.LessOrEqual(t, result.P99Latency, config.ThresholdP99)

	// Additional checks could include:
	// - Memory usage monitoring via API
	// - Latency trend analysis
	// - Connection pool monitoring
}

// BenchmarkCartOperations benchmarks cart operations
func BenchmarkCartOperations(b *testing.B) {
	client := &http.Client{Timeout: 5 * time.Second}
	baseURL := "http://localhost:8080"

	b.Run("GetCart", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", baseURL+"/api/v1/commerce/cart", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	})

	b.Run("AddToCart", func(b *testing.B) {
		body := map[string]interface{}{
			"productId": "benchmark-product",
			"quantity":  1,
			"unitPrice": 99.99,
		}
		bodyBytes, _ := json.Marshal(body)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", baseURL+"/api/v1/commerce/cart/items", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Authorization", "Bearer test-token")
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	})
}

// Helper function to generate test data
func generateTestData(count int) []map[string]interface{} {
	data := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		data[i] = map[string]interface{}{
			"id":        uuid.New().String(),
			"name":      fmt.Sprintf("Product %d", i),
			"price":     rand.Float64()*1000 + 10,
			"category":  []string{"electronics", "fashion", "home", "sports"}[rand.Intn(4)],
			"inStock":   rand.Intn(2) == 1,
			"quantity":  rand.Intn(100) + 1,
		}
	}
	return data
}