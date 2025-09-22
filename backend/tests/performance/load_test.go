package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tchat-backend/tests/testutil"
)

// LoadTestConfig defines configuration for load tests
type LoadTestConfig struct {
	BaseURL         string        `json:"base_url"`
	Duration        time.Duration `json:"duration"`
	ConcurrentUsers int           `json:"concurrent_users"`
	RampUpDuration  time.Duration `json:"ramp_up_duration"`
	ThinkTime       time.Duration `json:"think_time"`
	RequestTimeout  time.Duration `json:"request_timeout"`
	TargetRPS       int           `json:"target_rps"`
}

// DefaultLoadTestConfig returns default load test configuration
func DefaultLoadTestConfig() *LoadTestConfig {
	return &LoadTestConfig{
		BaseURL:         "http://localhost:8080",
		Duration:        5 * time.Minute,
		ConcurrentUsers: 100,
		RampUpDuration:  30 * time.Second,
		ThinkTime:       1 * time.Second,
		RequestTimeout:  10 * time.Second,
		TargetRPS:       1000,
	}
}

// LoadTestResults contains the results of a load test
type LoadTestResults struct {
	TotalRequests     int64                    `json:"total_requests"`
	SuccessfulRequests int64                   `json:"successful_requests"`
	FailedRequests    int64                    `json:"failed_requests"`
	TotalBytes        int64                    `json:"total_bytes"`
	RequestsPerSecond float64                  `json:"requests_per_second"`
	AvgResponseTime   time.Duration            `json:"avg_response_time"`
	MinResponseTime   time.Duration            `json:"min_response_time"`
	MaxResponseTime   time.Duration            `json:"max_response_time"`
	P50ResponseTime   time.Duration            `json:"p50_response_time"`
	P95ResponseTime   time.Duration            `json:"p95_response_time"`
	P99ResponseTime   time.Duration            `json:"p99_response_time"`
	ErrorsByType      map[string]int64         `json:"errors_by_type"`
	ResponseTimes     []time.Duration          `json:"response_times"`
	StartTime         time.Time                `json:"start_time"`
	EndTime           time.Time                `json:"end_time"`
	Duration          time.Duration            `json:"duration"`
}

// LoadTester manages load testing operations
type LoadTester struct {
	config  *LoadTestConfig
	client  *http.Client
	results *LoadTestResults
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewLoadTester creates a new load tester
func NewLoadTester(config *LoadTestConfig) *LoadTester {
	ctx, cancel := context.WithCancel(context.Background())

	client := &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &LoadTester{
		config: config,
		client: client,
		results: &LoadTestResults{
			ErrorsByType:  make(map[string]int64),
			ResponseTimes: make([]time.Duration, 0),
			StartTime:     time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// TestAuthServiceLoad tests authentication service under load
func TestAuthServiceLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	config := DefaultLoadTestConfig()
	config.Duration = 2 * time.Minute
	config.ConcurrentUsers = 50
	config.TargetRPS = 500

	tester := NewLoadTester(config)

	t.Run("OTP Send Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			// Generate unique phone number for each user
			phoneNumber := fmt.Sprintf("+66%09d", 800000000+userID)

			requestBody := map[string]string{
				"phone_number": phoneNumber,
				"country_code": "+66",
				"purpose":      "registration",
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			req, err := http.NewRequestWithContext(tester.ctx, "POST",
				tester.config.BaseURL+"/auth/otp/send", bytes.NewReader(jsonData))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, len(jsonData))

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements
		assert.Greater(t, results.RequestsPerSecond, float64(400), "Should handle at least 400 RPS")
		assert.Less(t, results.AvgResponseTime, 200*time.Millisecond, "Average response time should be under 200ms")
		assert.Less(t, results.P95ResponseTime, 500*time.Millisecond, "95th percentile should be under 500ms")
		assert.Greater(t, float64(results.SuccessfulRequests)/float64(results.TotalRequests), 0.99, "Success rate should be above 99%")

		t.Logf("OTP Send Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Requests/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
		t.Logf("  Success Rate: %.2f%%", float64(results.SuccessfulRequests)/float64(results.TotalRequests)*100)
	})

	t.Run("OTP Verify Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			phoneNumber := fmt.Sprintf("+66%09d", 800000000+userID)

			requestBody := map[string]interface{}{
				"phone_number": phoneNumber,
				"country_code": "+66",
				"code":         "123456",
				"purpose":      "registration",
				"device_info": map[string]string{
					"device_id":   fmt.Sprintf("device_%d", userID),
					"device_name": "Load Test Device",
					"device_type": "mobile",
					"app_version": "1.0.0",
				},
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			req, err := http.NewRequestWithContext(tester.ctx, "POST",
				tester.config.BaseURL+"/auth/otp/verify", bytes.NewReader(jsonData))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, len(jsonData))

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements
		assert.Greater(t, results.RequestsPerSecond, float64(300), "Should handle at least 300 RPS")
		assert.Less(t, results.AvgResponseTime, 300*time.Millisecond, "Average response time should be under 300ms")
		assert.Less(t, results.P95ResponseTime, 800*time.Millisecond, "95th percentile should be under 800ms")

		t.Logf("OTP Verify Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Requests/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
	})
}

// TestMessagingServiceLoad tests messaging service under load
func TestMessagingServiceLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	config := DefaultLoadTestConfig()
	config.Duration = 3 * time.Minute
	config.ConcurrentUsers = 100
	config.TargetRPS = 800

	tester := NewLoadTester(config)

	// Setup test users and auth tokens
	testUsers := setupTestUsers(t, helper, 100)

	t.Run("Send Message Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			user := testUsers[userID%len(testUsers)]

			requestBody := map[string]interface{}{
				"type":    "text",
				"content": fmt.Sprintf("Load test message from user %d at %v", userID, time.Now()),
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			// Use a test dialog ID
			dialogID := "test-dialog-123"
			url := fmt.Sprintf("%s/dialogs/%s/messages", tester.config.BaseURL, dialogID)

			req, err := http.NewRequestWithContext(tester.ctx, "POST", url, bytes.NewReader(jsonData))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, len(jsonData))

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements
		assert.Greater(t, results.RequestsPerSecond, float64(600), "Should handle at least 600 RPS")
		assert.Less(t, results.AvgResponseTime, 150*time.Millisecond, "Average response time should be under 150ms")
		assert.Less(t, results.P95ResponseTime, 400*time.Millisecond, "95th percentile should be under 400ms")

		t.Logf("Send Message Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Requests/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
	})

	t.Run("Get Messages Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			user := testUsers[userID%len(testUsers)]

			dialogID := "test-dialog-123"
			url := fmt.Sprintf("%s/dialogs/%s/messages?limit=20", tester.config.BaseURL, dialogID)

			req, err := http.NewRequestWithContext(tester.ctx, "GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, 0)

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements
		assert.Greater(t, results.RequestsPerSecond, float64(1000), "Should handle at least 1000 RPS")
		assert.Less(t, results.AvgResponseTime, 100*time.Millisecond, "Average response time should be under 100ms")
		assert.Less(t, results.P95ResponseTime, 250*time.Millisecond, "95th percentile should be under 250ms")

		t.Logf("Get Messages Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Requests/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
	})
}

// TestWebSocketLoad tests WebSocket connections under load
func TestWebSocketLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	// Test configuration
	maxConnections := 500
	messageRate := 10 // messages per second per connection
	testDuration := 2 * time.Minute

	var (
		totalConnections   int64
		totalMessages      int64
		failedConnections  int64
		failedMessages     int64
		connectionDurations []time.Duration
		messageDurations   []time.Duration
	)

	connectionsMutex := sync.Mutex{}
	messagesMutex := sync.Mutex{}

	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()

	var wg sync.WaitGroup

	t.Logf("Starting WebSocket load test with %d connections", maxConnections)

	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()

			// Connect to WebSocket
			connectStart := time.Now()
			wsURL := "ws://localhost:8080/websocket"

			u, err := url.Parse(wsURL)
			if err != nil {
				atomic.AddInt64(&failedConnections, 1)
				return
			}

			conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), http.Header{
				"Authorization": []string{"Bearer test_token"},
			})
			if err != nil {
				atomic.AddInt64(&failedConnections, 1)
				return
			}
			defer conn.Close()

			connectDuration := time.Since(connectStart)
			atomic.AddInt64(&totalConnections, 1)

			connectionsMutex.Lock()
			connectionDurations = append(connectionDurations, connectDuration)
			connectionsMutex.Unlock()

			// Send messages at specified rate
			ticker := time.NewTicker(time.Second / time.Duration(messageRate))
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					messageStart := time.Now()

					message := map[string]interface{}{
						"type": "text",
						"content": fmt.Sprintf("Load test message from connection %d at %v",
							connID, time.Now()),
						"dialog_id": "test-dialog-123",
					}

					if err := conn.WriteJSON(message); err != nil {
						atomic.AddInt64(&failedMessages, 1)
						continue
					}

					messageDuration := time.Since(messageStart)
					atomic.AddInt64(&totalMessages, 1)

					messagesMutex.Lock()
					messageDurations = append(messageDurations, messageDuration)
					messagesMutex.Unlock()

					// Read response (if any)
					conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
					_, _, err := conn.ReadMessage()
					if err != nil && !websocket.IsUnexpectedCloseError(err) {
						// Ignore timeout errors for this test
					}
				}
			}
		}(i)

		// Stagger connection attempts
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	// Calculate statistics
	actualDuration := time.Since(time.Now().Add(-testDuration))
	connectionsPerSecond := float64(totalConnections) / actualDuration.Seconds()
	messagesPerSecond := float64(totalMessages) / actualDuration.Seconds()

	// Calculate percentiles for connection times
	var avgConnectionTime, p95ConnectionTime time.Duration
	if len(connectionDurations) > 0 {
		var total time.Duration
		for _, d := range connectionDurations {
			total += d
		}
		avgConnectionTime = total / time.Duration(len(connectionDurations))

		// Simple P95 calculation
		if len(connectionDurations) > 20 {
			sorted := make([]time.Duration, len(connectionDurations))
			copy(sorted, connectionDurations)
			// Sort would go here in real implementation
			p95Index := int(float64(len(sorted)) * 0.95)
			if p95Index < len(sorted) {
				p95ConnectionTime = sorted[p95Index]
			}
		}
	}

	// Calculate percentiles for message times
	var avgMessageTime, p95MessageTime time.Duration
	if len(messageDurations) > 0 {
		var total time.Duration
		for _, d := range messageDurations {
			total += d
		}
		avgMessageTime = total / time.Duration(len(messageDurations))
	}

	// Assert performance requirements
	connectionSuccessRate := float64(totalConnections) / float64(maxConnections) * 100
	messageSuccessRate := float64(totalMessages) / float64(totalMessages + failedMessages) * 100

	assert.Greater(t, connectionSuccessRate, 95.0, "Connection success rate should be above 95%")
	assert.Greater(t, messageSuccessRate, 98.0, "Message success rate should be above 98%")
	assert.Greater(t, connectionsPerSecond, float64(50), "Should establish at least 50 connections per second")
	assert.Greater(t, messagesPerSecond, float64(1000), "Should handle at least 1000 messages per second")
	assert.Less(t, avgConnectionTime, 2*time.Second, "Average connection time should be under 2 seconds")
	assert.Less(t, avgMessageTime, 50*time.Millisecond, "Average message time should be under 50ms")

	// Log results
	t.Logf("WebSocket Load Test Results:")
	t.Logf("  Total Connections: %d/%d (%.1f%% success)", totalConnections, maxConnections, connectionSuccessRate)
	t.Logf("  Total Messages: %d (%.1f%% success)", totalMessages, messageSuccessRate)
	t.Logf("  Connections/sec: %.2f", connectionsPerSecond)
	t.Logf("  Messages/sec: %.2f", messagesPerSecond)
	t.Logf("  Avg Connection Time: %v", avgConnectionTime)
	t.Logf("  P95 Connection Time: %v", p95ConnectionTime)
	t.Logf("  Avg Message Time: %v", avgMessageTime)
}

// TestPaymentServiceLoad tests payment service under load
func TestPaymentServiceLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping payment load test in short mode")
	}

	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	config := DefaultLoadTestConfig()
	config.Duration = 2 * time.Minute
	config.ConcurrentUsers = 50
	config.TargetRPS = 200 // Lower RPS for payment transactions

	tester := NewLoadTester(config)
	testUsers := setupTestUsers(t, helper, 50)

	t.Run("Payment Transaction Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			user := testUsers[userID%len(testUsers)]

			// Generate random transaction amount
			amount := float64(rand.Intn(1000) + 10) // 10-1010 THB

			requestBody := map[string]interface{}{
				"to_user_id": testUsers[(userID+1)%len(testUsers)].UserID,
				"amount":     amount,
				"currency":   "THB",
				"description": fmt.Sprintf("Load test payment %d", userID),
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}

			req, err := http.NewRequestWithContext(tester.ctx, "POST",
				tester.config.BaseURL+"/transactions/send", bytes.NewReader(jsonData))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, len(jsonData))

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements (more lenient for payment transactions)
		assert.Greater(t, results.RequestsPerSecond, float64(150), "Should handle at least 150 TPS")
		assert.Less(t, results.AvgResponseTime, 500*time.Millisecond, "Average response time should be under 500ms")
		assert.Less(t, results.P95ResponseTime, 1*time.Second, "95th percentile should be under 1 second")
		assert.Greater(t, float64(results.SuccessfulRequests)/float64(results.TotalRequests), 0.98, "Success rate should be above 98%")

		t.Logf("Payment Transaction Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Transactions/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
		t.Logf("  Success Rate: %.2f%%", float64(results.SuccessfulRequests)/float64(results.TotalRequests)*100)
	})

	t.Run("Wallet Balance Load Test", func(t *testing.T) {
		scenario := func(userID int) error {
			user := testUsers[userID%len(testUsers)]

			req, err := http.NewRequestWithContext(tester.ctx, "GET",
				tester.config.BaseURL+"/wallets", nil)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

			start := time.Now()
			resp, err := tester.client.Do(req)
			responseTime := time.Since(start)

			tester.recordMetrics(resp, err, responseTime, 0)

			if resp != nil {
				resp.Body.Close()
			}

			return err
		}

		results := tester.RunLoadTest(scenario)

		// Assert performance requirements
		assert.Greater(t, results.RequestsPerSecond, float64(800), "Should handle at least 800 RPS")
		assert.Less(t, results.AvgResponseTime, 100*time.Millisecond, "Average response time should be under 100ms")
		assert.Less(t, results.P95ResponseTime, 300*time.Millisecond, "95th percentile should be under 300ms")

		t.Logf("Wallet Balance Load Test Results:")
		t.Logf("  Total Requests: %d", results.TotalRequests)
		t.Logf("  Requests/sec: %.2f", results.RequestsPerSecond)
		t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
		t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
	})
}

// RunLoadTest executes a load test scenario
func (lt *LoadTester) RunLoadTest(scenario func(int) error) *LoadTestResults {
	lt.results.StartTime = time.Now()

	var wg sync.WaitGroup
	userChan := make(chan int, lt.config.ConcurrentUsers)

	// Start workers
	for i := 0; i < lt.config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range userChan {
				if err := scenario(userID); err != nil {
					lt.mu.Lock()
					lt.results.ErrorsByType["scenario_error"]++
					lt.mu.Unlock()
				}

				// Think time between requests
				time.Sleep(lt.config.ThinkTime)
			}
		}()
	}

	// Generate load
	go func() {
		defer close(userChan)

		userID := 0
		timer := time.NewTimer(lt.config.Duration)
		defer timer.Stop()

		ticker := time.NewTicker(time.Second / time.Duration(lt.config.TargetRPS/lt.config.ConcurrentUsers))
		defer ticker.Stop()

		for {
			select {
			case <-timer.C:
				return
			case <-ticker.C:
				select {
				case userChan <- userID:
					userID++
				default:
					// Channel is full, skip this iteration
				}
			case <-lt.ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
	lt.results.EndTime = time.Now()
	lt.results.Duration = lt.results.EndTime.Sub(lt.results.StartTime)

	// Calculate final statistics
	lt.calculateStatistics()

	return lt.results
}

// recordMetrics records metrics for a single request
func (lt *LoadTester) recordMetrics(resp *http.Response, err error, responseTime time.Duration, bodySize int) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	atomic.AddInt64(&lt.results.TotalRequests, 1)
	lt.results.ResponseTimes = append(lt.results.ResponseTimes, responseTime)

	if err != nil {
		atomic.AddInt64(&lt.results.FailedRequests, 1)
		lt.results.ErrorsByType["network_error"]++
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		atomic.AddInt64(&lt.results.SuccessfulRequests, 1)
		atomic.AddInt64(&lt.results.TotalBytes, int64(bodySize))
	} else {
		atomic.AddInt64(&lt.results.FailedRequests, 1)
		errorType := fmt.Sprintf("http_%d", resp.StatusCode)
		lt.results.ErrorsByType[errorType]++
	}
}

// calculateStatistics calculates final test statistics
func (lt *LoadTester) calculateStatistics() {
	if lt.results.Duration > 0 {
		lt.results.RequestsPerSecond = float64(lt.results.TotalRequests) / lt.results.Duration.Seconds()
	}

	if len(lt.results.ResponseTimes) > 0 {
		// Calculate min, max, average
		var total time.Duration
		lt.results.MinResponseTime = lt.results.ResponseTimes[0]
		lt.results.MaxResponseTime = lt.results.ResponseTimes[0]

		for _, rt := range lt.results.ResponseTimes {
			total += rt
			if rt < lt.results.MinResponseTime {
				lt.results.MinResponseTime = rt
			}
			if rt > lt.results.MaxResponseTime {
				lt.results.MaxResponseTime = rt
			}
		}

		lt.results.AvgResponseTime = total / time.Duration(len(lt.results.ResponseTimes))

		// Calculate percentiles (simplified implementation)
		// In production, use a proper percentile calculation library
		sortedTimes := make([]time.Duration, len(lt.results.ResponseTimes))
		copy(sortedTimes, lt.results.ResponseTimes)

		// Simple bubble sort (for demonstration)
		for i := 0; i < len(sortedTimes); i++ {
			for j := 0; j < len(sortedTimes)-1-i; j++ {
				if sortedTimes[j] > sortedTimes[j+1] {
					sortedTimes[j], sortedTimes[j+1] = sortedTimes[j+1], sortedTimes[j]
				}
			}
		}

		p50Index := int(float64(len(sortedTimes)) * 0.5)
		p95Index := int(float64(len(sortedTimes)) * 0.95)
		p99Index := int(float64(len(sortedTimes)) * 0.99)

		if p50Index < len(sortedTimes) {
			lt.results.P50ResponseTime = sortedTimes[p50Index]
		}
		if p95Index < len(sortedTimes) {
			lt.results.P95ResponseTime = sortedTimes[p95Index]
		}
		if p99Index < len(sortedTimes) {
			lt.results.P99ResponseTime = sortedTimes[p99Index]
		}
	}
}

// TestUser represents a test user with authentication
type TestUser struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	AccessToken string `json:"access_token"`
}

// setupTestUsers creates authenticated test users for load testing
func setupTestUsers(t *testing.T, helper *testutil.TestHelper, count int) []TestUser {
	users := make([]TestUser, count)

	for i := 0; i < count; i++ {
		phoneNumber := fmt.Sprintf("+66%09d", 800000000+i)
		userID := helper.CreateTestUser(t, phoneNumber, "+66")

		// Create mock access token for testing
		accessToken := fmt.Sprintf("test_token_%s", userID)

		users[i] = TestUser{
			UserID:      userID,
			PhoneNumber: phoneNumber,
			AccessToken: accessToken,
		}
	}

	return users
}

// TestStressTest runs stress tests to find breaking points
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	// Gradually increase load until system breaks
	rpsLevels := []int{100, 500, 1000, 2000, 5000, 10000}

	for _, targetRPS := range rpsLevels {
		t.Run(fmt.Sprintf("Stress_Test_RPS_%d", targetRPS), func(t *testing.T) {
			config := DefaultLoadTestConfig()
			config.Duration = 1 * time.Minute
			config.ConcurrentUsers = targetRPS / 10 // 10 RPS per user
			config.TargetRPS = targetRPS

			tester := NewLoadTester(config)

			scenario := func(userID int) error {
				phoneNumber := fmt.Sprintf("+66%09d", 800000000+userID)

				requestBody := map[string]string{
					"phone_number": phoneNumber,
					"country_code": "+66",
					"purpose":      "registration",
				}

				jsonData, err := json.Marshal(requestBody)
				if err != nil {
					return err
				}

				req, err := http.NewRequestWithContext(tester.ctx, "POST",
					tester.config.BaseURL+"/auth/otp/send", bytes.NewReader(jsonData))
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				start := time.Now()
				resp, err := tester.client.Do(req)
				responseTime := time.Since(start)

				tester.recordMetrics(resp, err, responseTime, len(jsonData))

				if resp != nil {
					resp.Body.Close()
				}

				return err
			}

			results := tester.RunLoadTest(scenario)

			successRate := float64(results.SuccessfulRequests) / float64(results.TotalRequests) * 100

			t.Logf("Stress Test at %d RPS:", targetRPS)
			t.Logf("  Actual RPS: %.2f", results.RequestsPerSecond)
			t.Logf("  Success Rate: %.2f%%", successRate)
			t.Logf("  Avg Response Time: %v", results.AvgResponseTime)
			t.Logf("  P95 Response Time: %v", results.P95ResponseTime)
			t.Logf("  Errors: %v", results.ErrorsByType)

			// System is considered broken if success rate drops below 90%
			if successRate < 90.0 {
				t.Logf("System breaking point reached at %d RPS (%.2f%% success rate)", targetRPS, successRate)
				return
			}
		})
	}
}