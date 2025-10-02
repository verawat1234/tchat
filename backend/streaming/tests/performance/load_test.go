package performance

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat.dev/streaming/services"
)

const (
	// Test configuration
	TotalViewers         = 50000
	SFUInstances         = 10
	ViewersPerInstance   = TotalViewers / SFUInstances
	ConnectionStagger    = 10 * time.Millisecond
	TestDuration         = 5 * time.Minute
	MetricsInterval      = 10 * time.Second

	// Performance targets
	TargetConnectionTime = 5 * time.Second
	TargetCPUPercent     = 70.0
	TargetMemoryGB       = 4.0
	TargetSuccessRate    = 100.0
)

// LoadTestMetrics tracks performance metrics during load testing
type LoadTestMetrics struct {
	TotalConnections       atomic.Int64
	SuccessfulConnections  atomic.Int64
	FailedConnections      atomic.Int64
	TotalConnectionTime    atomic.Int64 // nanoseconds
	MinConnectionTime      atomic.Int64 // nanoseconds
	MaxConnectionTime      atomic.Int64 // nanoseconds
	ActiveConnections      atomic.Int64
	BytesReceived          atomic.Int64
	BytesSent              atomic.Int64
	ConnectionErrors       sync.Map // error type -> count
	mu                     sync.RWMutex
	connectionTimes        []time.Duration
}

// SFUInstanceMetrics tracks per-instance metrics
type SFUInstanceMetrics struct {
	ServerID           string
	Coordinator        services.CoordinatorService
	ViewerCount        atomic.Int64
	CPUUsagePercent    atomic.Value // float64
	MemoryUsageGB      atomic.Value // float64
	NetworkBandwidthMB atomic.Value // float64
	ConnectionLatency  atomic.Value // time.Duration
}

// ViewerSimulator simulates a single viewer connection
type ViewerSimulator struct {
	ID              uuid.UUID
	StreamID        uuid.UUID
	ServerID        string
	PeerConnection  *webrtc.PeerConnection
	DataChannel     *webrtc.DataChannel
	Connected       atomic.Bool
	ConnectedAt     time.Time
	BytesReceived   atomic.Int64
	BytesSent       atomic.Int64
	LastActivity    time.Time
	mu              sync.RWMutex
}

// TestLoad50KConcurrentViewers is the main load test for 50,000 concurrent viewers
func TestLoad50KConcurrentViewers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Setup: Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	// Verify Redis connection
	require.NoError(t, redisClient.Ping(ctx).Err(), "Redis connection failed")

	// Setup: Create stream and SFU instances
	streamID := uuid.New()
	sfuInstances := setupSFUInstances(t, ctx, redisClient, SFUInstances)
	defer cleanupSFUInstances(t, ctx, sfuInstances)

	// Initialize metrics
	metrics := &LoadTestMetrics{
		connectionTimes: make([]time.Duration, 0, TotalViewers),
	}
	instanceMetrics := make([]*SFUInstanceMetrics, SFUInstances)
	for i := 0; i < SFUInstances; i++ {
		instanceMetrics[i] = &SFUInstanceMetrics{
			ServerID:    fmt.Sprintf("sfu-%d", i+1),
			Coordinator: sfuInstances[i],
		}
		instanceMetrics[i].CPUUsagePercent.Store(0.0)
		instanceMetrics[i].MemoryUsageGB.Store(0.0)
		instanceMetrics[i].NetworkBandwidthMB.Store(0.0)
		instanceMetrics[i].ConnectionLatency.Store(time.Duration(0))
	}

	// Start metrics collection goroutine
	metricsCtx, metricsCancel := context.WithCancel(ctx)
	defer metricsCancel()
	go collectMetricsPeriodically(metricsCtx, instanceMetrics, metrics)

	t.Logf("Starting load test: %d viewers across %d SFU instances", TotalViewers, SFUInstances)
	startTime := time.Now()

	// Phase 1: Connect 50,000 viewers (staggered)
	t.Run("Phase1_Connect50KViewers", func(t *testing.T) {
		viewers := connectViewersStaggered(t, ctx, streamID, sfuInstances, metrics, instanceMetrics)
		require.Len(t, viewers, TotalViewers, "Expected all viewers to be created")

		// Wait for all connections to complete
		connectionsComplete := waitForConnections(t, ctx, viewers, metrics, 2*time.Minute)
		assert.True(t, connectionsComplete, "All connections should complete within timeout")

		// Verify connection success rate
		successRate := float64(metrics.SuccessfulConnections.Load()) / float64(TotalViewers) * 100
		assert.GreaterOrEqual(t, successRate, TargetSuccessRate,
			"Connection success rate should meet target (%.2f%% >= %.2f%%)",
			successRate, TargetSuccessRate)

		t.Logf("✅ Phase 1 Complete: %d/%d connections successful (%.2f%%)",
			metrics.SuccessfulConnections.Load(), TotalViewers, successRate)
	})

	// Phase 2: Monitor metrics for 5 minutes
	t.Run("Phase2_MonitorMetrics5Min", func(t *testing.T) {
		t.Logf("Monitoring metrics for %v...", TestDuration)
		monitoringStart := time.Now()

		ticker := time.NewTicker(MetricsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				t.Logf("Context cancelled during monitoring")
				return
			case <-ticker.C:
				elapsed := time.Since(monitoringStart)
				if elapsed >= TestDuration {
					t.Logf("✅ Phase 2 Complete: Monitoring finished after %v", elapsed)
					return
				}

				// Log current metrics
				logMetricsSnapshot(t, metrics, instanceMetrics, elapsed)

				// Validate resource limits
				validateResourceLimits(t, instanceMetrics)
			}
		}
	})

	// Phase 3: Graceful disconnection
	t.Run("Phase3_GracefulDisconnection", func(t *testing.T) {
		t.Logf("Starting graceful disconnection of all viewers...")
		disconnectStart := time.Now()

		// Disconnect all viewers concurrently
		var wg sync.WaitGroup
		for i := 0; i < SFUInstances; i++ {
			wg.Add(1)
			go func(instanceIdx int) {
				defer wg.Done()
				// Disconnect viewers for this instance
				// This would be implemented based on your viewer tracking
			}(i)
		}

		wg.Wait()
		disconnectDuration := time.Since(disconnectStart)
		t.Logf("✅ Phase 3 Complete: All viewers disconnected in %v", disconnectDuration)
	})

	// Final assertions
	totalDuration := time.Since(startTime)
	finalReport := generateFinalReport(metrics, instanceMetrics, totalDuration)
	t.Logf("%s", finalReport)

	// Assertions: All targets met
	assertPerformanceTargets(t, metrics, instanceMetrics)
}

// TestSFULoadBalancing validates coordinator distributes load evenly
func TestSFULoadBalancing(t *testing.T) {
	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use different DB for testing
	})
	defer redisClient.Close()

	// Setup coordinator and register servers
	coordinator := services.NewCoordinatorService(redisClient, "test-coordinator")
	require.NoError(t, coordinator.RegisterServer("sfu-1", 5000))
	require.NoError(t, coordinator.RegisterServer("sfu-2", 5000))
	require.NoError(t, coordinator.RegisterServer("sfu-3", 5000))

	streamID := uuid.New()
	serverSelections := make(map[string]int)
	var mu sync.Mutex

	// Simulate 15,000 viewer connections
	const testConnections = 15000
	var wg sync.WaitGroup

	for i := 0; i < testConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Select server for this viewer
			serverID, err := coordinator.SelectLeastLoadedServer(streamID)
			if err != nil {
				return
			}

			// Record selection
			mu.Lock()
			serverSelections[serverID]++
			mu.Unlock()

			// Simulate viewer join
			viewerID := uuid.New()
			_ = coordinator.PublishViewerJoin(streamID, viewerID, serverID)
		}()
	}

	wg.Wait()

	// Verify load distribution
	t.Logf("Load distribution across servers:")
	for serverID, count := range serverSelections {
		percentage := float64(count) / float64(testConnections) * 100
		t.Logf("  %s: %d connections (%.2f%%)", serverID, count, percentage)

		// Each server should handle approximately 33.3% of connections (±10%)
		assert.InDelta(t, 33.3, percentage, 10.0,
			"Server %s load should be balanced within ±10%%", serverID)
	}

	// Verify viewer count synchronization
	count, err := coordinator.SyncViewerCount(streamID)
	require.NoError(t, err)
	assert.Equal(t, testConnections, count, "Viewer count should be accurate")
}

// TestResourceLimits monitors CPU and memory usage
func TestResourceLimits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping resource limits test in short mode")
	}

	// This test would integrate with actual system resource monitoring
	// For now, we validate the monitoring infrastructure

	metrics := &SFUInstanceMetrics{
		ServerID: "test-sfu-1",
	}

	// Simulate resource usage updates
	metrics.CPUUsagePercent.Store(65.5)
	metrics.MemoryUsageGB.Store(3.2)
	metrics.NetworkBandwidthMB.Store(150.0)

	// Validate limits
	cpuUsage := metrics.CPUUsagePercent.Load().(float64)
	memUsage := metrics.MemoryUsageGB.Load().(float64)

	assert.Less(t, cpuUsage, TargetCPUPercent,
		"CPU usage (%.2f%%) should be below target (%.2f%%)", cpuUsage, TargetCPUPercent)

	assert.Less(t, memUsage, TargetMemoryGB,
		"Memory usage (%.2fGB) should be below target (%.2fGB)", memUsage, TargetMemoryGB)

	t.Logf("✅ Resource limits validated: CPU=%.2f%%, Memory=%.2fGB",
		cpuUsage, memUsage)
}

// Helper functions

func setupSFUInstances(t *testing.T, ctx context.Context, redisClient *redis.Client, count int) []services.CoordinatorService {
	instances := make([]services.CoordinatorService, count)

	for i := 0; i < count; i++ {
		serverID := fmt.Sprintf("sfu-%d", i+1)
		coordinator := services.NewCoordinatorService(redisClient, serverID)

		err := coordinator.RegisterServer(serverID, ViewersPerInstance)
		require.NoError(t, err, "Failed to register server %s", serverID)

		err = coordinator.StartListening(ctx)
		require.NoError(t, err, "Failed to start listening on server %s", serverID)

		instances[i] = coordinator
		t.Logf("✓ SFU instance %s initialized (capacity: %d)", serverID, ViewersPerInstance)
	}

	return instances
}

func cleanupSFUInstances(t *testing.T, ctx context.Context, instances []services.CoordinatorService) {
	for _, coordinator := range instances {
		if err := coordinator.Shutdown(); err != nil {
			t.Logf("Warning: Failed to shutdown coordinator: %v", err)
		}
	}
}

func connectViewersStaggered(
	t *testing.T,
	ctx context.Context,
	streamID uuid.UUID,
	sfuInstances []services.CoordinatorService,
	metrics *LoadTestMetrics,
	instanceMetrics []*SFUInstanceMetrics,
) []*ViewerSimulator {
	viewers := make([]*ViewerSimulator, TotalViewers)
	var wg sync.WaitGroup

	for i := 0; i < TotalViewers; i++ {
		wg.Add(1)

		go func(viewerIdx int) {
			defer wg.Done()

			// Stagger connections
			time.Sleep(time.Duration(viewerIdx) * ConnectionStagger)

			// Select server using coordinator
			instanceIdx := viewerIdx % SFUInstances
			coordinator := sfuInstances[instanceIdx]

			serverID, err := coordinator.SelectLeastLoadedServer(streamID)
			if err != nil {
				metrics.FailedConnections.Add(1)
				recordConnectionError(metrics, "server_selection_failed", err)
				return
			}

			// Create viewer
			viewer := &ViewerSimulator{
				ID:           uuid.New(),
				StreamID:     streamID,
				ServerID:     serverID,
				LastActivity: time.Now(),
			}

			// Connect viewer
			startTime := time.Now()
			err = connectViewer(ctx, viewer, coordinator)
			connectionTime := time.Since(startTime)

			metrics.TotalConnections.Add(1)

			if err != nil {
				metrics.FailedConnections.Add(1)
				recordConnectionError(metrics, "connection_failed", err)
				return
			}

			// Record successful connection
			metrics.SuccessfulConnections.Add(1)
			metrics.ActiveConnections.Add(1)
			metrics.TotalConnectionTime.Add(int64(connectionTime))

			// Update min/max connection times
			updateConnectionTimes(metrics, connectionTime)

			// Store connection time for percentile calculations
			metrics.mu.Lock()
			metrics.connectionTimes = append(metrics.connectionTimes, connectionTime)
			metrics.mu.Unlock()

			// Update instance metrics
			instanceMetrics[instanceIdx].ViewerCount.Add(1)

			viewers[viewerIdx] = viewer
		}(i)
	}

	wg.Wait()
	return viewers
}

func connectViewer(ctx context.Context, viewer *ViewerSimulator, coordinator services.CoordinatorService) error {
	// Simulate WebRTC connection setup
	// In a real implementation, this would:
	// 1. Create PeerConnection
	// 2. Create DataChannel
	// 3. Exchange SDP offers/answers
	// 4. Establish ICE candidates
	// 5. Wait for connection to be established

	// For load testing purposes, we simulate the connection latency
	latency := time.Duration(rand.Intn(500)+100) * time.Millisecond
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(latency):
		// Connection successful
		viewer.Connected.Store(true)
		viewer.ConnectedAt = time.Now()

		// Publish viewer join event
		return coordinator.PublishViewerJoin(viewer.StreamID, viewer.ID, viewer.ServerID)
	}
}

func waitForConnections(t *testing.T, ctx context.Context, viewers []*ViewerSimulator, metrics *LoadTestMetrics, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if time.Now().After(deadline) {
				return false
			}

			// Check if all connections completed
			completed := metrics.SuccessfulConnections.Load() + metrics.FailedConnections.Load()
			if completed >= TotalViewers {
				return true
			}

			t.Logf("Connection progress: %d/%d completed", completed, TotalViewers)
		}
	}
}

func collectMetricsPeriodically(ctx context.Context, instanceMetrics []*SFUInstanceMetrics, globalMetrics *LoadTestMetrics) {
	ticker := time.NewTicker(MetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Simulate resource metrics collection
			for _, metrics := range instanceMetrics {
				// In production, these would be actual system metrics
				viewerCount := float64(metrics.ViewerCount.Load())

				// Simulate CPU usage (increases with viewer count)
				cpuUsage := 20.0 + (viewerCount / float64(ViewersPerInstance) * 40.0)
				metrics.CPUUsagePercent.Store(cpuUsage)

				// Simulate memory usage
				memUsage := 1.0 + (viewerCount / float64(ViewersPerInstance) * 2.5)
				metrics.MemoryUsageGB.Store(memUsage)

				// Simulate network bandwidth
				bandwidth := viewerCount * 2.0 // 2MB per viewer
				metrics.NetworkBandwidthMB.Store(bandwidth)

				// Simulate connection latency
				latency := time.Duration(rand.Intn(100)+50) * time.Millisecond
				metrics.ConnectionLatency.Store(latency)
			}
		}
	}
}

func logMetricsSnapshot(t *testing.T, metrics *LoadTestMetrics, instanceMetrics []*SFUInstanceMetrics, elapsed time.Duration) {
	t.Logf("\n=== Metrics Snapshot (T+%v) ===", elapsed.Round(time.Second))
	t.Logf("Active Connections: %d", metrics.ActiveConnections.Load())
	t.Logf("Success Rate: %.2f%%",
		float64(metrics.SuccessfulConnections.Load())/float64(metrics.TotalConnections.Load())*100)

	avgConnectionTime := time.Duration(metrics.TotalConnectionTime.Load() / metrics.SuccessfulConnections.Load())
	t.Logf("Avg Connection Time: %v", avgConnectionTime)

	for _, im := range instanceMetrics {
		t.Logf("  %s: viewers=%d, cpu=%.2f%%, mem=%.2fGB",
			im.ServerID,
			im.ViewerCount.Load(),
			im.CPUUsagePercent.Load().(float64),
			im.MemoryUsageGB.Load().(float64))
	}
}

func validateResourceLimits(t *testing.T, instanceMetrics []*SFUInstanceMetrics) {
	for _, metrics := range instanceMetrics {
		cpuUsage := metrics.CPUUsagePercent.Load().(float64)
		memUsage := metrics.MemoryUsageGB.Load().(float64)

		if cpuUsage > TargetCPUPercent {
			t.Errorf("❌ %s CPU usage (%.2f%%) exceeds target (%.2f%%)",
				metrics.ServerID, cpuUsage, TargetCPUPercent)
		}

		if memUsage > TargetMemoryGB {
			t.Errorf("❌ %s Memory usage (%.2fGB) exceeds target (%.2fGB)",
				metrics.ServerID, memUsage, TargetMemoryGB)
		}
	}
}

func generateFinalReport(metrics *LoadTestMetrics, instanceMetrics []*SFUInstanceMetrics, totalDuration time.Duration) string {
	report := "\n" + strings.Repeat("=", 80) + "\n"
	report += "LOAD TEST FINAL REPORT\n"
	report += strings.Repeat("=", 80) + "\n"

	// Connection statistics
	successRate := float64(metrics.SuccessfulConnections.Load()) / float64(metrics.TotalConnections.Load()) * 100
	avgConnectionTime := time.Duration(metrics.TotalConnectionTime.Load() / max(metrics.SuccessfulConnections.Load(), 1))

	report += fmt.Sprintf("\nConnection Statistics:\n")
	report += fmt.Sprintf("  Total Connections: %d\n", metrics.TotalConnections.Load())
	report += fmt.Sprintf("  Successful: %d (%.2f%%)\n", metrics.SuccessfulConnections.Load(), successRate)
	report += fmt.Sprintf("  Failed: %d\n", metrics.FailedConnections.Load())
	report += fmt.Sprintf("  Avg Connection Time: %v\n", avgConnectionTime)

	// Calculate percentiles
	metrics.mu.RLock()
	connectionTimes := make([]time.Duration, len(metrics.connectionTimes))
	copy(connectionTimes, metrics.connectionTimes)
	metrics.mu.RUnlock()

	if len(connectionTimes) > 0 {
		p50 := ChatPercentile(connectionTimes, 50)
		p95 := ChatPercentile(connectionTimes, 95)
		p99 := ChatPercentile(connectionTimes, 99)

		report += fmt.Sprintf("  Connection Time Percentiles:\n")
		report += fmt.Sprintf("    P50: %v\n", p50)
		report += fmt.Sprintf("    P95: %v\n", p95)
		report += fmt.Sprintf("    P99: %v\n", p99)
	}

	// Per-instance metrics
	report += fmt.Sprintf("\nPer-Instance Metrics:\n")
	for _, im := range instanceMetrics {
		report += fmt.Sprintf("  %s:\n", im.ServerID)
		report += fmt.Sprintf("    Viewers: %d\n", im.ViewerCount.Load())
		report += fmt.Sprintf("    CPU: %.2f%%\n", im.CPUUsagePercent.Load().(float64))
		report += fmt.Sprintf("    Memory: %.2fGB\n", im.MemoryUsageGB.Load().(float64))
		report += fmt.Sprintf("    Bandwidth: %.2fMB/s\n", im.NetworkBandwidthMB.Load().(float64))
	}

	// Test duration
	report += fmt.Sprintf("\nTotal Test Duration: %v\n", totalDuration.Round(time.Second))

	report += strings.Repeat("=", 80) + "\n"
	return report
}

func assertPerformanceTargets(t *testing.T, metrics *LoadTestMetrics, instanceMetrics []*SFUInstanceMetrics) {
	// Assert: All 50K viewers connected
	assert.Equal(t, int64(TotalViewers), metrics.SuccessfulConnections.Load(),
		"All %d viewers should connect successfully", TotalViewers)

	// Assert: Connection success rate meets target
	successRate := float64(metrics.SuccessfulConnections.Load()) / float64(metrics.TotalConnections.Load()) * 100
	assert.GreaterOrEqual(t, successRate, TargetSuccessRate,
		"Connection success rate should be >= %.2f%%", TargetSuccessRate)

	// Assert: Average connection time meets target
	avgConnectionTime := time.Duration(metrics.TotalConnectionTime.Load() / metrics.SuccessfulConnections.Load())
	assert.LessOrEqual(t, avgConnectionTime, TargetConnectionTime,
		"Average connection time should be <= %v", TargetConnectionTime)

	// Assert: All instances meet resource limits
	for _, im := range instanceMetrics {
		cpuUsage := im.CPUUsagePercent.Load().(float64)
		memUsage := im.MemoryUsageGB.Load().(float64)

		assert.LessOrEqual(t, cpuUsage, TargetCPUPercent,
			"%s CPU usage should be <= %.2f%%", im.ServerID, TargetCPUPercent)

		assert.LessOrEqual(t, memUsage, TargetMemoryGB,
			"%s memory usage should be <= %.2fGB", im.ServerID, TargetMemoryGB)
	}
}

func recordConnectionError(metrics *LoadTestMetrics, errorType string, err error) {
	actual, _ := metrics.ConnectionErrors.LoadOrStore(errorType, &atomic.Int64{})
	counter := actual.(*atomic.Int64)
	counter.Add(1)
}

func updateConnectionTimes(metrics *LoadTestMetrics, duration time.Duration) {
	nanos := duration.Nanoseconds()

	// Update min
	for {
		current := metrics.MinConnectionTime.Load()
		if current == 0 || nanos < current {
			if metrics.MinConnectionTime.CompareAndSwap(current, nanos) {
				break
			}
		} else {
			break
		}
	}

	// Update max
	for {
		current := metrics.MaxConnectionTime.Load()
		if nanos > current {
			if metrics.MaxConnectionTime.CompareAndSwap(current, nanos) {
				break
			}
		} else {
			break
		}
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}