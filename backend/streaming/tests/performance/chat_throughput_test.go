package performance

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

// Helper function to convert google/uuid.UUID to gocql.UUID
func toGocqlUUID(googleUUID uuid.UUID) gocql.UUID {
	var gocqlUUID gocql.UUID
	copy(gocqlUUID[:], googleUUID[:])
	return gocqlUUID
}

const (
	// Performance targets
	targetThroughput      = 100000 // 100K messages/second
	targetWriteLatencyP99 = 5 * time.Millisecond
	targetQueryLatency    = 10 * time.Millisecond
	targetTTLDays         = 30
)

// MetricsCollector tracks performance metrics
type MetricsCollector struct {
	mu                sync.Mutex
	writeLatencies    []time.Duration
	queryLatencies    []time.Duration
	writeSuccesses    int64
	writeFailures     int64
	totalMessages     int64
	startTime         time.Time
	endTime           time.Time
	partitionDistrib  map[uuid.UUID]int64
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		writeLatencies:   make([]time.Duration, 0, 10000),
		queryLatencies:   make([]time.Duration, 0, 1000),
		partitionDistrib: make(map[uuid.UUID]int64),
		startTime:        time.Now(),
	}
}

// RecordWrite records a write operation
func (m *MetricsCollector) RecordWrite(latency time.Duration, streamID uuid.UUID, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if success {
		atomic.AddInt64(&m.writeSuccesses, 1)
		m.writeLatencies = append(m.writeLatencies, latency)
		m.partitionDistrib[streamID]++
	} else {
		atomic.AddInt64(&m.writeFailures, 1)
	}
	atomic.AddInt64(&m.totalMessages, 1)
}

// RecordQuery records a query operation
func (m *MetricsCollector) RecordQuery(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queryLatencies = append(m.queryLatencies, latency)
}

// CalculateStats calculates performance statistics
func (m *MetricsCollector) CalculateStats() *PerformanceStats {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.endTime = time.Now()
	duration := m.endTime.Sub(m.startTime)

	stats := &PerformanceStats{
		Duration:           duration,
		TotalMessages:      atomic.LoadInt64(&m.totalMessages),
		SuccessfulWrites:   atomic.LoadInt64(&m.writeSuccesses),
		FailedWrites:       atomic.LoadInt64(&m.writeFailures),
		ThroughputPerSec:   float64(atomic.LoadInt64(&m.totalMessages)) / duration.Seconds(),
		PartitionCount:     int64(len(m.partitionDistrib)),
	}

	if len(m.writeLatencies) > 0 {
		stats.WriteLatencyP50 = ChatPercentile(m.writeLatencies, 50)
		stats.WriteLatencyP95 = ChatPercentile(m.writeLatencies, 95)
		stats.WriteLatencyP99 = ChatPercentile(m.writeLatencies, 99)
		stats.WriteLatencyAvg = ChatAverage(m.writeLatencies)
	}

	if len(m.queryLatencies) > 0 {
		stats.QueryLatencyP50 = ChatPercentile(m.queryLatencies, 50)
		stats.QueryLatencyP95 = ChatPercentile(m.queryLatencies, 95)
		stats.QueryLatencyP99 = ChatPercentile(m.queryLatencies, 99)
		stats.QueryLatencyAvg = ChatAverage(m.queryLatencies)
	}

	// Calculate partition distribution uniformity
	if len(m.partitionDistrib) > 0 {
		var total, sumSquares float64
		for _, count := range m.partitionDistrib {
			total += float64(count)
			sumSquares += float64(count * count)
		}
		mean := total / float64(len(m.partitionDistrib))
		variance := (sumSquares / float64(len(m.partitionDistrib))) - (mean * mean)
		stats.PartitionStdDev = math.Sqrt(variance)
	}

	return stats
}

// PerformanceStats holds performance statistics
type PerformanceStats struct {
	Duration          time.Duration
	TotalMessages     int64
	SuccessfulWrites  int64
	FailedWrites      int64
	ThroughputPerSec  float64
	WriteLatencyP50   time.Duration
	WriteLatencyP95   time.Duration
	WriteLatencyP99   time.Duration
	WriteLatencyAvg   time.Duration
	QueryLatencyP50   time.Duration
	QueryLatencyP95   time.Duration
	QueryLatencyP99   time.Duration
	QueryLatencyAvg   time.Duration
	PartitionCount    int64
	PartitionStdDev   float64
}

// String formats performance stats for display
func (s *PerformanceStats) String() string {
	return fmt.Sprintf(`
Performance Test Results:
==========================
Duration:              %v
Total Messages:        %d
Successful Writes:     %d
Failed Writes:         %d
Throughput:            %.2f msg/s
Write Latency (p50):   %v
Write Latency (p95):   %v
Write Latency (p99):   %v
Write Latency (avg):   %v
Query Latency (p50):   %v
Query Latency (p95):   %v
Query Latency (p99):   %v
Query Latency (avg):   %v
Partitions:            %d
Partition StdDev:      %.2f
`,
		s.Duration,
		s.TotalMessages,
		s.SuccessfulWrites,
		s.FailedWrites,
		s.ThroughputPerSec,
		s.WriteLatencyP50,
		s.WriteLatencyP95,
		s.WriteLatencyP99,
		s.WriteLatencyAvg,
		s.QueryLatencyP50,
		s.QueryLatencyP95,
		s.QueryLatencyP99,
		s.QueryLatencyAvg,
		s.PartitionCount,
		s.PartitionStdDev,
	)
}

// setupScyllaDB creates a ScyllaDB session for testing
func setupScyllaDB(t *testing.T, config *TestConfig) *gocql.Session {
	session, err := SetupChatScyllaDB(config)
	if err != nil {
		t.Skipf("ScyllaDB not available: %v", err)
	}
	return session
}

// TestChatThroughput100K validates 100,000 messages/second throughput
func TestChatThroughput100K(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	metrics := NewMetricsCollector()
	ctx := context.Background()

	// Create test streams for partition distribution
	numStreams := 100
	streamIDs := make([]uuid.UUID, numStreams)
	for i := 0; i < numStreams; i++ {
		streamIDs[i] = uuid.New()
	}

	t.Logf("Starting throughput test: %d writers × %d msg/s × %v = %.0f messages",
		config.NumWriters,
		config.MessagesPerSec,
		config.TestDuration,
		float64(config.NumWriters*config.MessagesPerSec)*config.TestDuration.Seconds())

	// Launch concurrent writers
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	for writerID := 0; writerID < config.NumWriters; writerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ticker := time.NewTicker(time.Second / time.Duration(config.MessagesPerSec))
			defer ticker.Stop()

			senderID := uuid.New()
			streamIndex := id % len(streamIDs)
			streamID := streamIDs[streamIndex]

			for {
				select {
				case <-stopChan:
					return
				case <-ticker.C:
					startTime := time.Now()
					message := &models.ChatMessage{
						StreamID:          streamID,
						MessageID:         uuid.New(),
						SenderID:          senderID,
						SenderDisplayName: fmt.Sprintf("User%d", id),
						MessageText:       fmt.Sprintf("Test message from writer %d", id),
						ModerationStatus:  models.ModerationStatusVisible,
						MessageType:       models.MessageTypeText,
						Timestamp:         time.Now(),
					}

					err := repo.Create(ctx, message)
					latency := time.Since(startTime)
					metrics.RecordWrite(latency, streamID, err == nil)
				}
			}
		}(writerID)
	}

	// Run test for configured duration
	time.Sleep(config.TestDuration)
	close(stopChan)
	wg.Wait()

	// Calculate and display results
	stats := metrics.CalculateStats()
	t.Log(stats.String())

	// Assertions
	assert.Equal(t, int64(0), stats.FailedWrites, "All writes should succeed")
	assert.GreaterOrEqual(t, stats.ThroughputPerSec, float64(targetThroughput)*0.95,
		"Throughput should be >= 95%% of target (95K msg/s)")
	assert.LessOrEqual(t, stats.WriteLatencyP99, targetWriteLatencyP99,
		"P99 write latency should be <= 5ms")

	t.Logf("✓ Throughput target: %.2f msg/s (target: %d msg/s)", stats.ThroughputPerSec, targetThroughput)
	t.Logf("✓ Write success rate: 100%% (%d/%d)", stats.SuccessfulWrites, stats.TotalMessages)
	t.Logf("✓ Write latency P99: %v (target: %v)", stats.WriteLatencyP99, targetWriteLatencyP99)
}

// TestChatWriteLatency measures write latency distribution
func TestChatWriteLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	ctx := context.Background()
	streamID := uuid.New()
	senderID := uuid.New()

	// Warm-up writes
	for i := 0; i < 100; i++ {
		message := &models.ChatMessage{
			StreamID:          streamID,
			MessageID:         uuid.New(),
			SenderID:          senderID,
			SenderDisplayName: "WarmupUser",
			MessageText:       "Warmup message",
			ModerationStatus:  models.ModerationStatusVisible,
			MessageType:       models.MessageTypeText,
		}
		repo.Create(ctx, message)
	}

	// Measure write latency
	metrics := NewMetricsCollector()
	numWrites := 1000

	for i := 0; i < numWrites; i++ {
		startTime := time.Now()
		message := &models.ChatMessage{
			StreamID:          streamID,
			MessageID:         uuid.New(),
			SenderID:          senderID,
			SenderDisplayName: "LatencyUser",
			MessageText:       fmt.Sprintf("Latency test message %d", i),
			ModerationStatus:  models.ModerationStatusVisible,
			MessageType:       models.MessageTypeText,
		}

		err := repo.Create(ctx, message)
		latency := time.Since(startTime)
		metrics.RecordWrite(latency, streamID, err == nil)
	}

	stats := metrics.CalculateStats()
	t.Logf("Write Latency Distribution (n=%d):", numWrites)
	t.Logf("  P50: %v", stats.WriteLatencyP50)
	t.Logf("  P95: %v", stats.WriteLatencyP95)
	t.Logf("  P99: %v", stats.WriteLatencyP99)
	t.Logf("  Avg: %v", stats.WriteLatencyAvg)

	assert.LessOrEqual(t, stats.WriteLatencyP99, targetWriteLatencyP99,
		"P99 write latency should be <= 5ms")
	assert.Equal(t, int64(0), stats.FailedWrites, "All writes should succeed")
}

// TestChatQueryPerformance validates query performance for recent messages
func TestChatQueryPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	ctx := context.Background()
	streamID := uuid.New()
	senderID := uuid.New()

	// Insert test messages
	numMessages := 1000
	for i := 0; i < numMessages; i++ {
		message := &models.ChatMessage{
			StreamID:          streamID,
			MessageID:         uuid.New(),
			SenderID:          senderID,
			SenderDisplayName: "QueryUser",
			MessageText:       fmt.Sprintf("Query test message %d", i),
			ModerationStatus:  models.ModerationStatusVisible,
			MessageType:       models.MessageTypeText,
		}
		err := repo.Create(ctx, message)
		require.NoError(t, err, "Failed to insert message")
		time.Sleep(time.Millisecond) // Ensure timestamp ordering
	}

	metrics := NewMetricsCollector()

	// Test: Retrieve recent 50 messages (target: <10ms)
	t.Run("Recent50Messages", func(t *testing.T) {
		numQueries := 100

		for i := 0; i < numQueries; i++ {
			startTime := time.Now()
			messages, err := repo.ListByStream(ctx, streamID, 50, nil)
			latency := time.Since(startTime)

			require.NoError(t, err, "Query should succeed")
			assert.Equal(t, 50, len(messages), "Should retrieve 50 messages")
			metrics.RecordQuery(latency)
		}

		stats := metrics.CalculateStats()
		t.Logf("Query Performance (Recent 50 messages, n=%d):", numQueries)
		t.Logf("  P50: %v", stats.QueryLatencyP50)
		t.Logf("  P95: %v", stats.QueryLatencyP95)
		t.Logf("  P99: %v", stats.QueryLatencyP99)
		t.Logf("  Avg: %v", stats.QueryLatencyAvg)

		assert.LessOrEqual(t, stats.QueryLatencyP99, targetQueryLatency,
			"P99 query latency should be <= 10ms")
	})

	// Test: Retrieve last 5 minutes of chat (target: <50ms)
	t.Run("Last5MinutesOfChat", func(t *testing.T) {
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

		startTime := time.Now()
		messages, err := repo.ListByStream(ctx, streamID, 500, &fiveMinutesAgo)
		latency := time.Since(startTime)

		require.NoError(t, err, "Query should succeed")
		t.Logf("Retrieved %d messages in %v", len(messages), latency)
		assert.LessOrEqual(t, latency, 50*time.Millisecond,
			"Query latency should be <= 50ms")
	})

	// Test: Pagination performance
	t.Run("PaginationPerformance", func(t *testing.T) {
		pageSize := 50
		var cursor *time.Time = nil
		totalRetrieved := 0

		for page := 0; page < 10; page++ {
			startTime := time.Now()
			messages, err := repo.ListByStream(ctx, streamID, pageSize, cursor)
			latency := time.Since(startTime)

			require.NoError(t, err, "Pagination query should succeed")
			totalRetrieved += len(messages)

			if len(messages) == 0 {
				break
			}

			// Update cursor to last message timestamp
			lastTimestamp := messages[len(messages)-1].Timestamp
			cursor = &lastTimestamp

			t.Logf("Page %d: %d messages in %v", page+1, len(messages), latency)
			assert.LessOrEqual(t, latency, targetQueryLatency,
				"Pagination query should be <= 10ms")
		}

		t.Logf("Total messages retrieved: %d", totalRetrieved)
	})
}

// TestChatTTLEnforcement verifies 30-day TTL enforcement
func TestChatTTLEnforcement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TTL test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	ctx := context.Background()
	streamID := uuid.New()
	senderID := uuid.New()

	// Insert message with TTL
	message := &models.ChatMessage{
		StreamID:          streamID,
		MessageID:         uuid.New(),
		SenderID:          senderID,
		SenderDisplayName: "TTLUser",
		MessageText:       "TTL test message",
		ModerationStatus:  models.ModerationStatusVisible,
		MessageType:       models.MessageTypeText,
	}

	err = repo.Create(ctx, message)
	require.NoError(t, err, "Message creation should succeed")

	// Verify message exists
	retrievedMessage, err := repo.GetByID(ctx, streamID, message.MessageID)
	require.NoError(t, err, "Message should exist")
	assert.Equal(t, message.MessageText, retrievedMessage.MessageText)

	// Query TTL from ScyllaDB
	var ttl int
	query := session.Query(`SELECT TTL(message_text) FROM chat_messages WHERE stream_id = ? AND message_id = ? LIMIT 1 ALLOW FILTERING`,
		toGocqlUUID(streamID), toGocqlUUID(message.MessageID))
	err = query.Scan(&ttl)
	require.NoError(t, err, "TTL query should succeed")

	expectedTTL := models.DefaultTTL
	tolerance := 60 // Allow 60 seconds tolerance for test execution time

	t.Logf("Message TTL: %d seconds (expected: %d seconds, ~%d days)",
		ttl, expectedTTL, expectedTTL/86400)

	assert.GreaterOrEqual(t, ttl, expectedTTL-tolerance,
		"TTL should be approximately 30 days (2592000 seconds)")
	assert.LessOrEqual(t, ttl, expectedTTL,
		"TTL should not exceed configured value")

	t.Logf("✓ TTL enforcement verified: %d days (~%d seconds)", targetTTLDays, ttl)
}

// TestScyllaDBPartitioning validates partition distribution uniformity
func TestScyllaDBPartitioning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping partitioning test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	ctx := context.Background()
	numStreams := 50
	messagesPerStream := 100

	streamIDs := make([]uuid.UUID, numStreams)
	for i := 0; i < numStreams; i++ {
		streamIDs[i] = uuid.New()
	}

	// Insert messages across multiple partitions
	metrics := NewMetricsCollector()
	senderID := uuid.New()

	for i := 0; i < numStreams; i++ {
		for j := 0; j < messagesPerStream; j++ {
			message := &models.ChatMessage{
				StreamID:          streamIDs[i],
				MessageID:         uuid.New(),
				SenderID:          senderID,
				SenderDisplayName: "PartitionUser",
				MessageText:       fmt.Sprintf("Partition test %d-%d", i, j),
				ModerationStatus:  models.ModerationStatusVisible,
				MessageType:       models.MessageTypeText,
			}

			startTime := time.Now()
			err := repo.Create(ctx, message)
			latency := time.Since(startTime)
			metrics.RecordWrite(latency, streamIDs[i], err == nil)
		}
	}

	stats := metrics.CalculateStats()

	t.Logf("Partition Distribution Analysis:")
	t.Logf("  Total Partitions: %d", stats.PartitionCount)
	t.Logf("  Messages per Partition (avg): %.2f", float64(stats.TotalMessages)/float64(stats.PartitionCount))
	t.Logf("  Partition StdDev: %.2f", stats.PartitionStdDev)

	assert.Equal(t, int64(numStreams), stats.PartitionCount,
		"Should have messages in all partitions")
	assert.LessOrEqual(t, stats.PartitionStdDev, float64(messagesPerStream)*0.1,
		"Partition distribution should be relatively uniform (StdDev <= 10%% of mean)")

	// Verify each partition has expected message count
	for i, streamID := range streamIDs {
		messages, err := repo.ListByStream(ctx, streamID, messagesPerStream+10, nil)
		require.NoError(t, err, "Query should succeed for partition %d", i)
		assert.Equal(t, messagesPerStream, len(messages),
			"Partition %d should have %d messages", i, messagesPerStream)
	}

	t.Logf("✓ Partition distribution verified: uniform across %d streams", numStreams)
}

// TestChatBatchPerformance validates batch insert performance
func TestChatBatchPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping batch performance test in short mode")
	}

	config := LoadTestConfig()
	session := setupScyllaDB(t, config)
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	require.NoError(t, err, "Failed to create repository")

	ctx := context.Background()
	streamID := uuid.New()
	senderID := uuid.New()

	// Test different batch sizes
	batchSizes := []int{10, 50, 100}

	for _, batchSize := range batchSizes {
		t.Run(fmt.Sprintf("BatchSize%d", batchSize), func(t *testing.T) {
			messages := make([]*models.ChatMessage, batchSize)
			for i := 0; i < batchSize; i++ {
				messages[i] = &models.ChatMessage{
					StreamID:          streamID,
					MessageID:         uuid.New(),
					SenderID:          senderID,
					SenderDisplayName: "BatchUser",
					MessageText:       fmt.Sprintf("Batch message %d", i),
					ModerationStatus:  models.ModerationStatusVisible,
					MessageType:       models.MessageTypeText,
				}
			}

			startTime := time.Now()
			err := repo.BatchCreate(ctx, messages)
			latency := time.Since(startTime)

			require.NoError(t, err, "Batch create should succeed")
			throughput := float64(batchSize) / latency.Seconds()

			t.Logf("Batch size %d: %v latency, %.2f msg/s throughput",
				batchSize, latency, throughput)

			// Batch writes should be faster than individual writes
			avgLatencyPerMessage := latency / time.Duration(batchSize)
			assert.LessOrEqual(t, avgLatencyPerMessage, targetWriteLatencyP99,
				"Average latency per message in batch should be <= 5ms")
		})
	}
}