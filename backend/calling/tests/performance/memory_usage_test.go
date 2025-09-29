package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tchat.dev/calling/config"
	"tchat.dev/calling/models"
	"tchat.dev/calling/repositories"
)

// MemoryUsageTestSuite defines the memory usage testing suite
type MemoryUsageTestSuite struct {
	suite.Suite
	db              *gorm.DB
	redisClient     redis.Cmdable
	redisMock       redismock.ClientMock
	callRepo        *repositories.CallRepository
	presenceService *config.UserPresenceService
	signalingService *config.SignalingService
}

// SetupSuite sets up the test environment
func (suite *MemoryUsageTestSuite) SetupSuite() {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	// Auto-migrate models
	err = db.AutoMigrate(&models.CallSession{}, &models.CallParticipant{})
	suite.Require().NoError(err)

	// Setup Redis mock
	rdb, mock := redismock.NewClientMock()
	suite.redisClient = rdb
	suite.redisMock = mock

	// Create Redis wrapper
	redisWrapper := &config.RedisClient{
		Client: rdb.(*redis.Client),
	}

	// Setup services
	suite.callRepo = repositories.NewCallRepository(suite.db)
	suite.presenceService = config.NewUserPresenceService(redisWrapper)
	suite.signalingService = config.NewSignalingService(redisWrapper)
}

// TearDownSuite cleans up the test environment
func (suite *MemoryUsageTestSuite) TearDownSuite() {
	if sqlDB, err := suite.db.DB(); err == nil {
		sqlDB.Close()
	}
	if redisClient, ok := suite.redisClient.(*redis.Client); ok {
		redisClient.Close()
	}
}

// MemoryStats holds memory usage statistics
type MemoryStats struct {
	AllocMB      float64
	TotalAllocMB float64
	SysMB        float64
	NumGC        uint32
	GCCPUPercent float64
	NumGoroutines int
}

// TestLongCallMemoryUsage tests memory usage during extended call sessions
func (suite *MemoryUsageTestSuite) TestLongCallMemoryUsage() {
	// Test configuration
	callDuration := 60 * time.Minute // 60 minutes
	samplingInterval := 1 * time.Minute
	maxMemoryMB := 100.0 // Maximum allowed memory usage in MB
	maxGoroutines := 1000 // Maximum allowed goroutines

	// Initial memory measurement
	initialStats := suite.getMemoryStats()
	suite.T().Logf("Initial Memory Stats: %+v", initialStats)

	// Create a long-running call session
	callSession := &models.CallSession{
		ID:          uuid.New().String(),
		Type:        "video",
		Status:      "active",
		InitiatedBy: uuid.New().String(),
		StartedAt:   time.Now(),
	}

	err := suite.callRepo.Create(callSession)
	suite.Require().NoError(err)

	// Setup Redis expectations for presence updates during the call
	totalSamples := int(callDuration / samplingInterval)
	for i := 0; i < totalSamples; i++ {
		userID := callSession.InitiatedBy
		key := "presence:user:" + userID
		presenceData := map[string]interface{}{
			"status":    "online",
			"in_call":   true,
			"call_id":   callSession.ID,
			"last_seen": time.Now().Unix(),
		}

		suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
		suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)
	}

	// Simulate long-running call with periodic memory sampling
	memoryHistory := make([]MemoryStats, 0, totalSamples)
	startTime := time.Now()

	for elapsed := time.Duration(0); elapsed < callDuration; elapsed += samplingInterval {
		// Simulate call activity (presence updates, signaling)
		err := suite.presenceService.SetUserOnline(callSession.InitiatedBy, callSession.ID)
		suite.Require().NoError(err)

		// Simulate some signaling activity
		// Note: We'll skip Redis expectations for signaling to focus on memory
		// In a real test, you'd set these up as well

		// Sample memory usage
		currentStats := suite.getMemoryStats()
		memoryHistory = append(memoryHistory, currentStats)

		// Log progress every 10 minutes
		if int(elapsed.Minutes())%10 == 0 {
			suite.T().Logf("Call Duration: %v, Memory: %.2f MB, Goroutines: %d",
				elapsed, currentStats.AllocMB, currentStats.NumGoroutines)
		}

		// Check memory constraints
		assert.True(suite.T(), currentStats.AllocMB < maxMemoryMB,
			"Memory usage %.2f MB exceeds limit %.2f MB at %v",
			currentStats.AllocMB, maxMemoryMB, elapsed)

		assert.True(suite.T(), currentStats.NumGoroutines < maxGoroutines,
			"Goroutine count %d exceeds limit %d at %v",
			currentStats.NumGoroutines, maxGoroutines, elapsed)

		// Sleep until next sampling interval
		time.Sleep(samplingInterval)
	}

	// Final memory measurement
	finalStats := suite.getMemoryStats()
	suite.T().Logf("Final Memory Stats: %+v", finalStats)

	// Analyze memory growth
	memoryGrowthMB := finalStats.AllocMB - initialStats.AllocMB
	maxMemoryUsed := suite.findMaxMemoryUsage(memoryHistory)

	suite.T().Logf("Memory Analysis:")
	suite.T().Logf("Initial Memory: %.2f MB", initialStats.AllocMB)
	suite.T().Logf("Final Memory: %.2f MB", finalStats.AllocMB)
	suite.T().Logf("Memory Growth: %.2f MB", memoryGrowthMB)
	suite.T().Logf("Peak Memory: %.2f MB", maxMemoryUsed)
	suite.T().Logf("Total Duration: %v", time.Since(startTime))

	// Memory usage assertions
	assert.True(suite.T(), memoryGrowthMB < 50.0,
		"Memory growth %.2f MB should be less than 50 MB for 60-minute call", memoryGrowthMB)
	assert.True(suite.T(), maxMemoryUsed < maxMemoryMB,
		"Peak memory usage %.2f MB should be less than %.2f MB", maxMemoryUsed, maxMemoryMB)

	// End the call session
	callSession.Status = "ended"
	callSession.EndedAt = &time.Time{}
	*callSession.EndedAt = time.Now()
	callSession.Duration = int64(time.Since(callSession.StartedAt).Seconds())

	err = suite.callRepo.Update(callSession)
	suite.Require().NoError(err)
}

// TestConcurrentCallsMemoryUsage tests memory usage with multiple concurrent calls
func (suite *MemoryUsageTestSuite) TestConcurrentCallsMemoryUsage() {
	concurrentCalls := 50
	callDuration := 10 * time.Minute
	maxMemoryPerCallMB := 2.0 // Maximum memory per call

	initialStats := suite.getMemoryStats()
	suite.T().Logf("Initial Memory (Concurrent Test): %+v", initialStats)

	// Setup Redis expectations for all concurrent calls
	totalPresenceUpdates := concurrentCalls * int(callDuration.Minutes())
	for i := 0; i < totalPresenceUpdates; i++ {
		userID := fmt.Sprintf("user-%d", i%concurrentCalls)
		callID := fmt.Sprintf("call-%d", i%concurrentCalls)
		key := "presence:user:" + userID
		presenceData := map[string]interface{}{
			"status":    "online",
			"in_call":   true,
			"call_id":   callID,
			"last_seen": time.Now().Unix(),
		}

		suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
		suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)
	}

	// Create concurrent call sessions
	var wg sync.WaitGroup
	callSessions := make([]*models.CallSession, concurrentCalls)

	for i := 0; i < concurrentCalls; i++ {
		callSession := &models.CallSession{
			ID:          fmt.Sprintf("call-%d", i),
			Type:        "video",
			Status:      "active",
			InitiatedBy: fmt.Sprintf("user-%d", i),
			StartedAt:   time.Now(),
		}

		err := suite.callRepo.Create(callSession)
		suite.Require().NoError(err)
		callSessions[i] = callSession

		wg.Add(1)
		go func(index int, session *models.CallSession) {
			defer wg.Done()

			// Simulate call activity for the duration
			endTime := time.Now().Add(callDuration)
			for time.Now().Before(endTime) {
				// Update presence every minute
				err := suite.presenceService.SetUserOnline(session.InitiatedBy, session.ID)
				assert.NoError(suite.T(), err)

				time.Sleep(1 * time.Minute)
			}
		}(i, callSession)
	}

	// Monitor memory usage during concurrent calls
	monitorDone := make(chan bool)
	go func() {
		defer close(monitorDone)

		for {
			select {
			case <-time.After(2 * time.Minute):
				currentStats := suite.getMemoryStats()
				memoryPerCall := currentStats.AllocMB / float64(concurrentCalls)

				suite.T().Logf("Concurrent Calls Memory: Total=%.2f MB, Per Call=%.2f MB, Goroutines=%d",
					currentStats.AllocMB, memoryPerCall, currentStats.NumGoroutines)

				assert.True(suite.T(), memoryPerCall < maxMemoryPerCallMB,
					"Memory per call %.2f MB exceeds limit %.2f MB", memoryPerCall, maxMemoryPerCallMB)

			case <-monitorDone:
				return
			}
		}
	}()

	// Wait for all calls to complete
	wg.Wait()
	close(monitorDone)

	finalStats := suite.getMemoryStats()
	memoryGrowth := finalStats.AllocMB - initialStats.AllocMB
	memoryPerCall := memoryGrowth / float64(concurrentCalls)

	suite.T().Logf("Concurrent Calls Results:")
	suite.T().Logf("Total Memory Growth: %.2f MB", memoryGrowth)
	suite.T().Logf("Memory Per Call: %.2f MB", memoryPerCall)
	suite.T().Logf("Final Goroutines: %d", finalStats.NumGoroutines)

	// Assertions for concurrent calls
	assert.True(suite.T(), memoryPerCall < maxMemoryPerCallMB,
		"Average memory per call %.2f MB should be less than %.2f MB", memoryPerCall, maxMemoryPerCallMB)
	assert.True(suite.T(), finalStats.AllocMB < float64(concurrentCalls)*maxMemoryPerCallMB*2,
		"Total memory usage should be reasonable for %d concurrent calls", concurrentCalls)
}

// TestMemoryLeakDetection tests for memory leaks in call lifecycle
func (suite *MemoryUsageTestSuite) TestMemoryLeakDetection() {
	iterations := 100
	callsPerIteration := 10

	initialStats := suite.getMemoryStats()

	// Force garbage collection to get clean baseline
	runtime.GC()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	baselineStats := suite.getMemoryStats()
	suite.T().Logf("Baseline Memory: %.2f MB", baselineStats.AllocMB)

	// Setup Redis expectations for all iterations
	totalCalls := iterations * callsPerIteration
	for i := 0; i < totalCalls*2; i++ { // *2 for online + offline operations
		userID := fmt.Sprintf("leak-test-user-%d", i)

		if i < totalCalls {
			// Online expectations
			key := "presence:user:" + userID
			presenceData := map[string]interface{}{
				"status":    "online",
				"in_call":   true,
				"call_id":   fmt.Sprintf("leak-test-call-%d", i),
				"last_seen": time.Now().Unix(),
			}
			suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
			suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)
		} else {
			// Offline expectations
			key := "presence:user:" + fmt.Sprintf("leak-test-user-%d", i-totalCalls)
			suite.redisMock.ExpectDel(key).SetVal(1)
		}
	}

	memoryHistory := make([]float64, 0, iterations)

	// Run multiple iterations of call creation and cleanup
	for iteration := 0; iteration < iterations; iteration++ {
		// Create multiple calls
		for callIndex := 0; callIndex < callsPerIteration; callIndex++ {
			callSession := &models.CallSession{
				ID:          fmt.Sprintf("leak-test-call-%d-%d", iteration, callIndex),
				Type:        "voice",
				Status:      "active",
				InitiatedBy: fmt.Sprintf("leak-test-user-%d-%d", iteration, callIndex),
				StartedAt:   time.Now(),
			}

			err := suite.callRepo.Create(callSession)
			suite.Require().NoError(err)

			// Set user online
			err = suite.presenceService.SetUserOnline(callSession.InitiatedBy, callSession.ID)
			suite.Require().NoError(err)
		}

		// Simulate some call activity
		time.Sleep(10 * time.Millisecond)

		// End all calls and clean up
		for callIndex := 0; callIndex < callsPerIteration; callIndex++ {
			userID := fmt.Sprintf("leak-test-user-%d-%d", iteration, callIndex)

			// Set user offline
			err := suite.presenceService.SetUserOffline(userID)
			suite.Require().NoError(err)
		}

		// Force garbage collection every 10 iterations
		if iteration%10 == 0 {
			runtime.GC()
			time.Sleep(50 * time.Millisecond)
		}

		// Sample memory usage
		currentStats := suite.getMemoryStats()
		memoryHistory = append(memoryHistory, currentStats.AllocMB)

		if iteration%20 == 0 {
			suite.T().Logf("Iteration %d: Memory %.2f MB", iteration, currentStats.AllocMB)
		}
	}

	// Final garbage collection
	runtime.GC()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	finalStats := suite.getMemoryStats()
	memoryGrowth := finalStats.AllocMB - baselineStats.AllocMB

	suite.T().Logf("Memory Leak Detection Results:")
	suite.T().Logf("Baseline Memory: %.2f MB", baselineStats.AllocMB)
	suite.T().Logf("Final Memory: %.2f MB", finalStats.AllocMB)
	suite.T().Logf("Memory Growth: %.2f MB", memoryGrowth)
	suite.T().Logf("Total Calls Processed: %d", totalCalls)

	// Analyze memory trend
	memoryTrend := suite.calculateMemoryTrend(memoryHistory)
	suite.T().Logf("Memory Trend: %.4f MB/iteration", memoryTrend)

	// Memory leak assertions
	assert.True(suite.T(), memoryGrowth < 10.0,
		"Memory growth %.2f MB should be less than 10 MB after %d calls", memoryGrowth, totalCalls)
	assert.True(suite.T(), memoryTrend < 0.01,
		"Memory trend %.4f MB/iteration should be minimal (< 0.01)", memoryTrend)
}

// getMemoryStats collects current memory statistics
func (suite *MemoryUsageTestSuite) getMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		AllocMB:       float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB:  float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:         float64(m.Sys) / 1024 / 1024,
		NumGC:         m.NumGC,
		GCCPUPercent:  m.GCCPUFraction * 100,
		NumGoroutines: runtime.NumGoroutine(),
	}
}

// findMaxMemoryUsage finds the peak memory usage from history
func (suite *MemoryUsageTestSuite) findMaxMemoryUsage(history []MemoryStats) float64 {
	if len(history) == 0 {
		return 0.0
	}

	max := history[0].AllocMB
	for _, stats := range history {
		if stats.AllocMB > max {
			max = stats.AllocMB
		}
	}
	return max
}

// calculateMemoryTrend calculates the memory growth trend
func (suite *MemoryUsageTestSuite) calculateMemoryTrend(history []float64) float64 {
	if len(history) < 2 {
		return 0.0
	}

	// Simple linear regression to find trend
	n := float64(len(history))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range history {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope (trend)
	trend := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return trend
}

// TestMemoryUsageTestSuite runs the memory usage testing suite
func TestMemoryUsageTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryUsageTestSuite))
}

// Benchmark tests for memory efficiency

func BenchmarkCallSessionCreation(b *testing.B) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.AutoMigrate(&models.CallSession{})
	if err != nil {
		b.Fatal(err)
	}

	callRepo := repositories.NewCallRepository(db)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		callSession := &models.CallSession{
			ID:          uuid.New().String(),
			Type:        "voice",
			Status:      "connecting",
			InitiatedBy: uuid.New().String(),
			StartedAt:   time.Now(),
		}

		callRepo.Create(callSession)
	}
}

func BenchmarkPresenceUpdates(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	presenceService := config.NewUserPresenceService(redisWrapper)

	userID := uuid.New().String()
	callID := uuid.New().String()

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		key := "presence:user:" + userID
		presenceData := map[string]interface{}{
			"status":    "online",
			"in_call":   true,
			"call_id":   callID,
			"last_seen": time.Now().Unix(),
		}

		mock.ExpectHMSet(key, presenceData).SetVal(true)
		mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		presenceService.SetUserOnline(userID, callID)
	}
}