package performance

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"tchat.dev/calling/config"
)

// WebRTCQualityTestSuite defines the WebRTC quality testing suite
type WebRTCQualityTestSuite struct {
	suite.Suite
	redisClient   redis.Cmdable
	redisMock     redismock.ClientMock
	signalingService *config.SignalingService
}

// SetupTest sets up test dependencies
func (suite *WebRTCQualityTestSuite) SetupTest() {
	// Create Redis mock
	db, mock := redismock.NewClientMock()
	suite.redisClient = db
	suite.redisMock = mock

	// Create Redis client wrapper
	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}

	// Create signaling service
	suite.signalingService = config.NewSignalingService(redisWrapper)
}

// TearDownTest cleans up after each test
func (suite *WebRTCQualityTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.redisMock.ExpectationsWereMet())
}

// LatencyMetrics holds latency measurement data
type LatencyMetrics struct {
	AverageLatency    time.Duration
	MaxLatency        time.Duration
	MinLatency        time.Duration
	Percentile95      time.Duration
	Percentile99      time.Duration
	MessagesProcessed int
	PacketLoss        float64
}

// TestSignalingLatency tests signaling message latency under load
func (suite *WebRTCQualityTestSuite) TestSignalingLatency() {
	targetLatency := 200 * time.Millisecond
	messageCount := 1000
	fromUserID := uuid.New().String()
	toUserID := uuid.New().String()

	latencies := make([]time.Duration, 0, messageCount)
	var latencyMutex sync.Mutex

	// Setup Redis expectations for signaling messages
	for i := 0; i < messageCount; i++ {
		key := fmt.Sprintf("signaling:%s", toUserID)
		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   fmt.Sprintf("test-message-%d", i),
		}

		suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
		suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	var wg sync.WaitGroup

	// Send signaling messages and measure latency
	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		go func(messageIndex int) {
			defer wg.Done()

			startTime := time.Now()

			// Store signaling message
			message := fmt.Sprintf("test-message-%d", messageIndex)
			err := suite.signalingService.StoreSignalingMessage(fromUserID, toUserID, message)

			latency := time.Since(startTime)

			latencyMutex.Lock()
			latencies = append(latencies, latency)
			latencyMutex.Unlock()

			assert.NoError(suite.T(), err)
		}(i)
	}

	wg.Wait()

	// Calculate latency metrics
	metrics := suite.calculateLatencyMetrics(latencies)

	suite.T().Logf("Signaling Latency Test Results:")
	suite.T().Logf("Messages Processed: %d", metrics.MessagesProcessed)
	suite.T().Logf("Average Latency: %v", metrics.AverageLatency)
	suite.T().Logf("Max Latency: %v", metrics.MaxLatency)
	suite.T().Logf("Min Latency: %v", metrics.MinLatency)
	suite.T().Logf("95th Percentile: %v", metrics.Percentile95)
	suite.T().Logf("99th Percentile: %v", metrics.Percentile99)

	// Quality assertions
	assert.True(suite.T(), metrics.AverageLatency < targetLatency,
		"Average latency %v should be less than %v", metrics.AverageLatency, targetLatency)
	assert.True(suite.T(), metrics.Percentile95 < targetLatency*2,
		"95th percentile %v should be less than %v", metrics.Percentile95, targetLatency*2)
	assert.Equal(suite.T(), messageCount, metrics.MessagesProcessed,
		"All messages should be processed")
}

// TestConcurrentSignaling tests signaling quality under concurrent load
func (suite *WebRTCQualityTestSuite) TestConcurrentSignaling() {
	concurrentUsers := 50
	messagesPerUser := 20
	totalMessages := concurrentUsers * messagesPerUser

	// Setup Redis expectations for all messages
	for i := 0; i < totalMessages; i++ {
		userID := fmt.Sprintf("user-%d", i%concurrentUsers)
		key := fmt.Sprintf("signaling:%s", userID)

		expectedMessageData := map[string]interface{}{
			"from":      fmt.Sprintf("from-user-%d", i),
			"timestamp": time.Now().Unix(),
			"message":   fmt.Sprintf("concurrent-message-%d", i),
		}

		suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
		suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	latencies := make([]time.Duration, 0, totalMessages)
	var latencyMutex sync.Mutex
	var wg sync.WaitGroup

	startTime := time.Now()

	// Launch concurrent users
	for userIndex := 0; userIndex < concurrentUsers; userIndex++ {
		wg.Add(1)
		go func(uIndex int) {
			defer wg.Done()

			fromUserID := fmt.Sprintf("from-user-%d", uIndex)
			toUserID := fmt.Sprintf("user-%d", uIndex)

			for msgIndex := 0; msgIndex < messagesPerUser; msgIndex++ {
				msgStartTime := time.Now()

				message := fmt.Sprintf("concurrent-message-%d", uIndex*messagesPerUser+msgIndex)
				err := suite.signalingService.StoreSignalingMessage(fromUserID, toUserID, message)

				latency := time.Since(msgStartTime)

				latencyMutex.Lock()
				latencies = append(latencies, latency)
				latencyMutex.Unlock()

				assert.NoError(suite.T(), err)
			}
		}(userIndex)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// Calculate metrics
	metrics := suite.calculateLatencyMetrics(latencies)

	suite.T().Logf("Concurrent Signaling Test Results:")
	suite.T().Logf("Concurrent Users: %d", concurrentUsers)
	suite.T().Logf("Total Messages: %d", totalMessages)
	suite.T().Logf("Total Duration: %v", totalDuration)
	suite.T().Logf("Average Latency: %v", metrics.AverageLatency)
	suite.T().Logf("Messages/Second: %.2f", float64(totalMessages)/totalDuration.Seconds())
	suite.T().Logf("95th Percentile Latency: %v", metrics.Percentile95)

	// Quality assertions for concurrent load
	assert.True(suite.T(), metrics.AverageLatency < 500*time.Millisecond,
		"Average latency under concurrent load should be reasonable")
	assert.True(suite.T(), metrics.Percentile95 < time.Second,
		"95th percentile should be under 1 second even under load")
	assert.Equal(suite.T(), totalMessages, metrics.MessagesProcessed,
		"All concurrent messages should be processed")
}

// TestWebRTCMessageTypes tests latency for different WebRTC message types
func (suite *WebRTCQualityTestSuite) TestWebRTCMessageTypes() {
	messageTypes := []struct {
		name    string
		message string
		size    int
	}{
		{"SDP Offer", `{"type":"offer","sdp":"v=0\r\no=- 123 0 IN IP4 127.0.0.1\r\n..."}`, 500},
		{"SDP Answer", `{"type":"answer","sdp":"v=0\r\no=- 456 0 IN IP4 127.0.0.1\r\n..."}`, 400},
		{"ICE Candidate", `{"type":"candidate","candidate":"candidate:1 1 UDP 2130706431 192.168.1.100 54400 typ host","sdpMLineIndex":0}`, 150},
		{"Large SDP", fmt.Sprintf(`{"type":"offer","sdp":"%s"}`, generateLargeSDP(2000)), 2000},
	}

	for _, msgType := range messageTypes {
		suite.T().Run(msgType.name, func(t *testing.T) {
			messageCount := 100
			fromUserID := uuid.New().String()
			toUserID := uuid.New().String()

			// Setup Redis expectations
			for i := 0; i < messageCount; i++ {
				key := fmt.Sprintf("signaling:%s", toUserID)
				expectedMessageData := map[string]interface{}{
					"from":      fromUserID,
					"timestamp": time.Now().Unix(),
					"message":   msgType.message,
				}

				suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
				suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
				suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)
			}

			latencies := make([]time.Duration, 0, messageCount)

			// Test message type latency
			for i := 0; i < messageCount; i++ {
				startTime := time.Now()

				err := suite.signalingService.StoreSignalingMessage(fromUserID, toUserID, msgType.message)
				latency := time.Since(startTime)

				latencies = append(latencies, latency)
				assert.NoError(t, err)
			}

			metrics := suite.calculateLatencyMetrics(latencies)

			t.Logf("%s Message Type Results:", msgType.name)
			t.Logf("Message Size: %d bytes", msgType.size)
			t.Logf("Average Latency: %v", metrics.AverageLatency)
			t.Logf("95th Percentile: %v", metrics.Percentile95)

			// Message type specific assertions
			expectedLatency := time.Duration(msgType.size/10) * time.Microsecond // Rough estimate
			if expectedLatency < 50*time.Millisecond {
				expectedLatency = 50 * time.Millisecond
			}

			assert.True(t, metrics.AverageLatency < expectedLatency*10,
				"Average latency for %s should be reasonable", msgType.name)
		})
	}
}

// TestNetworkConditions simulates different network conditions
func (suite *WebRTCQualityTestSuite) TestNetworkConditions() {
	networkConditions := []struct {
		name           string
		latencyDelay   time.Duration
		maxLatency     time.Duration
		messageCount   int
	}{
		{"Good Connection", 0, 100 * time.Millisecond, 200},
		{"Average Connection", 50 * time.Millisecond, 300 * time.Millisecond, 100},
		{"Poor Connection", 200 * time.Millisecond, 1 * time.Second, 50},
	}

	for _, condition := range networkConditions {
		suite.T().Run(condition.name, func(t *testing.T) {
			fromUserID := uuid.New().String()
			toUserID := uuid.New().String()

			// Setup Redis expectations
			for i := 0; i < condition.messageCount; i++ {
				key := fmt.Sprintf("signaling:%s", toUserID)
				expectedMessageData := map[string]interface{}{
					"from":      fromUserID,
					"timestamp": time.Now().Unix(),
					"message":   fmt.Sprintf("network-test-message-%d", i),
				}

				suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
				suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
				suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)
			}

			latencies := make([]time.Duration, 0, condition.messageCount)

			for i := 0; i < condition.messageCount; i++ {
				// Simulate network delay
				if condition.latencyDelay > 0 {
					time.Sleep(condition.latencyDelay)
				}

				startTime := time.Now()

				message := fmt.Sprintf("network-test-message-%d", i)
				err := suite.signalingService.StoreSignalingMessage(fromUserID, toUserID, message)

				latency := time.Since(startTime)
				latencies = append(latencies, latency)

				assert.NoError(t, err)
			}

			metrics := suite.calculateLatencyMetrics(latencies)

			t.Logf("%s Results:", condition.name)
			t.Logf("Average Latency: %v", metrics.AverageLatency)
			t.Logf("Max Latency: %v", metrics.MaxLatency)
			t.Logf("95th Percentile: %v", metrics.Percentile95)

			// Network condition specific assertions
			assert.True(t, metrics.MaxLatency < condition.maxLatency*2,
				"Max latency should be within reasonable bounds for %s", condition.name)
		})
	}
}

// calculateLatencyMetrics computes latency statistics
func (suite *WebRTCQualityTestSuite) calculateLatencyMetrics(latencies []time.Duration) LatencyMetrics {
	if len(latencies) == 0 {
		return LatencyMetrics{}
	}

	var total time.Duration
	minLatency := latencies[0]
	maxLatency := latencies[0]

	for _, latency := range latencies {
		total += latency
		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
	}

	avgLatency := total / time.Duration(len(latencies))

	// Calculate percentiles (simplified)
	p95Index := int(float64(len(latencies)) * 0.95)
	p99Index := int(float64(len(latencies)) * 0.99)

	if p95Index >= len(latencies) {
		p95Index = len(latencies) - 1
	}
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}

	return LatencyMetrics{
		AverageLatency:    avgLatency,
		MaxLatency:        maxLatency,
		MinLatency:        minLatency,
		Percentile95:      latencies[p95Index],
		Percentile99:      latencies[p99Index],
		MessagesProcessed: len(latencies),
		PacketLoss:        0.0, // Simplified for this test
	}
}

// generateLargeSDP creates a large SDP for testing
func generateLargeSDP(size int) string {
	sdp := "v=0\r\no=- 123 0 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	// Add media lines to reach desired size
	for len(sdp) < size {
		sdp += "m=audio 54400 RTP/SAVPF 111\r\na=rtpmap:111 OPUS/48000/2\r\n"
	}

	return sdp[:size]
}

// TestWebRTCQualityTestSuite runs the WebRTC quality testing suite
func TestWebRTCQualityTestSuite(t *testing.T) {
	suite.Run(t, new(WebRTCQualityTestSuite))
}

// Benchmark tests for WebRTC performance

func BenchmarkSignalingMessageStorage(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	signalingService := config.NewSignalingService(redisWrapper)

	fromUserID := uuid.New().String()
	toUserID := uuid.New().String()
	message := "benchmark-test-message"

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("signaling:%s", toUserID)
		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   message,
		}

		mock.ExpectLPush(key, expectedMessageData).SetVal(1)
		mock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		mock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signalingService.StoreSignalingMessage(fromUserID, toUserID, message)
	}
}

func BenchmarkConcurrentSignaling(b *testing.B) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	signalingService := config.NewSignalingService(redisWrapper)

	// Setup expectations for parallel benchmark
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("bench-user-%d", i%10)
		key := fmt.Sprintf("signaling:%s", userID)
		expectedMessageData := map[string]interface{}{
			"from":      fmt.Sprintf("from-user-%d", i),
			"timestamp": time.Now().Unix(),
			"message":   "concurrent-benchmark-message",
		}

		mock.ExpectLPush(key, expectedMessageData).SetVal(1)
		mock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		mock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			fromUserID := fmt.Sprintf("from-user-%d", i)
			toUserID := fmt.Sprintf("bench-user-%d", i%10)
			signalingService.StoreSignalingMessage(fromUserID, toUserID, "concurrent-benchmark-message")
			i++
		}
	})
}