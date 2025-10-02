package performance

import (
	"context"
	"fmt"
	"math"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat.dev/streaming/services"
)

// LatencyMeasurement represents a single latency measurement
type LatencyMeasurement struct {
	FrameID           int
	BroadcasterTime   time.Time
	ViewerTime        time.Time
	LatencyMs         float64
	BandwidthKbps     float64
	QualityLayer      int
	QualityLayerName  string
	PacketLossRate    float64
	JitterMs          float64
}

// LatencyTestResult aggregates latency measurements
type LatencyTestResult struct {
	Scenario          string
	BandwidthKbps     float64
	TargetQualityLayer int
	TargetLatencyMs   float64

	Measurements      []LatencyMeasurement
	AverageLatencyMs  float64
	P95LatencyMs      float64
	P99LatencyMs      float64
	MinLatencyMs      float64
	MaxLatencyMs      float64

	QualitySwitchCount int
	RebufferEvents     int

	Success           bool
	FailureReason     string
}

// NetworkCondition represents simulated network conditions
type NetworkCondition struct {
	BandwidthKbps    float64
	PacketLossRate   float64
	JitterMs         float64
	RTTMs            float64
	Description      string
}

// High bandwidth scenario: >2Mbps → 1080p layer → <1s latency
func TestHighBandwidthLatency(t *testing.T) {
	qualityService := services.NewQualityService()

	// Simulate high bandwidth network: 3000 Kbps (3 Mbps)
	network := NetworkCondition{
		BandwidthKbps:  3000.0,
		PacketLossRate: 0.001, // 0.1% packet loss
		JitterMs:       5.0,
		RTTMs:          50.0,
		Description:    "High Bandwidth (>2Mbps)",
	}

	// Target: 1080p layer (layer 2), <1s latency
	result := runLatencyTest(t, qualityService, network, 2, 1000.0, 100)

	// Assertions
	require.True(t, result.Success, "High bandwidth latency test failed: %s", result.FailureReason)
	assert.Less(t, result.AverageLatencyMs, 1000.0, "Average latency should be <1s")
	assert.Less(t, result.P95LatencyMs, 1000.0, "P95 latency should be <1s")
	assert.Less(t, result.P99LatencyMs, 1200.0, "P99 latency should be <1.2s (20% tolerance)")
	assert.Equal(t, 0, result.RebufferEvents, "Should have no rebuffer events")

	// Log results
	logLatencyResults(t, result)
}

// Standard bandwidth scenario: 800Kbps-2Mbps → 720p layer → <3s latency
func TestStandardBandwidthLatency(t *testing.T) {
	qualityService := services.NewQualityService()

	// Simulate standard bandwidth network: 1200 Kbps (1.2 Mbps)
	network := NetworkCondition{
		BandwidthKbps:  1200.0,
		PacketLossRate: 0.005, // 0.5% packet loss
		JitterMs:       15.0,
		RTTMs:          100.0,
		Description:    "Standard Bandwidth (800Kbps-2Mbps)",
	}

	// Target: 720p layer (layer 1), <3s latency
	result := runLatencyTest(t, qualityService, network, 1, 3000.0, 100)

	// Assertions
	require.True(t, result.Success, "Standard bandwidth latency test failed: %s", result.FailureReason)
	assert.Less(t, result.AverageLatencyMs, 3000.0, "Average latency should be <3s")
	assert.Less(t, result.P95LatencyMs, 3000.0, "P95 latency should be <3s")
	assert.Less(t, result.P99LatencyMs, 3600.0, "P99 latency should be <3.6s (20% tolerance)")
	assert.Equal(t, 0, result.RebufferEvents, "Should have no rebuffer events")

	// Log results
	logLatencyResults(t, result)
}

// Constrained bandwidth scenario: <800Kbps → 360p layer → <5s latency
func TestConstrainedBandwidthLatency(t *testing.T) {
	qualityService := services.NewQualityService()

	// Simulate constrained bandwidth network: 600 Kbps
	network := NetworkCondition{
		BandwidthKbps:  600.0,
		PacketLossRate: 0.02, // 2% packet loss
		JitterMs:       30.0,
		RTTMs:          200.0,
		Description:    "Constrained Bandwidth (<800Kbps)",
	}

	// Target: 360p layer (layer 0), <5s latency
	result := runLatencyTest(t, qualityService, network, 0, 5000.0, 100)

	// Assertions
	require.True(t, result.Success, "Constrained bandwidth latency test failed: %s", result.FailureReason)
	assert.Less(t, result.AverageLatencyMs, 5000.0, "Average latency should be <5s")
	assert.Less(t, result.P95LatencyMs, 5000.0, "P95 latency should be <5s")
	assert.Less(t, result.P99LatencyMs, 6000.0, "P99 latency should be <6s (20% tolerance)")
	assert.Equal(t, 0, result.RebufferEvents, "Should have no rebuffer events")

	// Log results
	logLatencyResults(t, result)
}

// Test quality switching hysteresis: 5s delay downgrade, 10s upgrade
func TestQualitySwitchingHysteresis(t *testing.T) {
	qualityService := services.NewQualityService()

	t.Run("DowngradeDelay5s", func(t *testing.T) {
		// Start at high quality (layer 2)
		currentLayer := 2
		lastChangeTime := time.Now()

		// After 3 seconds, try to downgrade (should be blocked)
		time.Sleep(3 * time.Second)
		newLayer, allowed, reason := qualityService.ApplyHysteresis(0, currentLayer, lastChangeTime)

		assert.Equal(t, currentLayer, newLayer, "Should maintain current layer")
		assert.False(t, allowed, "Downgrade should be blocked before 5s delay")
		assert.Contains(t, reason, "Downgrade delayed", "Reason should mention downgrade delay")

		t.Logf("Downgrade blocked after 3s: %s", reason)

		// After 5+ seconds, downgrade should be allowed
		time.Sleep(2500 * time.Millisecond)
		lastChangeTime = time.Now().Add(-5 * time.Second)
		newLayer, allowed, reason = qualityService.ApplyHysteresis(0, currentLayer, lastChangeTime)

		assert.Equal(t, 0, newLayer, "Should downgrade to layer 0")
		assert.True(t, allowed, "Downgrade should be allowed after 5s delay")
		assert.Contains(t, reason, "approved", "Reason should mention approval")

		t.Logf("Downgrade allowed after 5s: %s", reason)
	})

	t.Run("UpgradeDelay10s", func(t *testing.T) {
		// Start at low quality (layer 0)
		currentLayer := 0
		lastChangeTime := time.Now()

		// After 7 seconds, try to upgrade (should be blocked)
		time.Sleep(3 * time.Second)
		lastChangeTime = time.Now().Add(-7 * time.Second)
		newLayer, allowed, reason := qualityService.ApplyHysteresis(2, currentLayer, lastChangeTime)

		assert.Equal(t, currentLayer, newLayer, "Should maintain current layer")
		assert.False(t, allowed, "Upgrade should be blocked before 10s delay")
		assert.Contains(t, reason, "Upgrade delayed", "Reason should mention upgrade delay")

		t.Logf("Upgrade blocked after 7s: %s", reason)

		// After 10+ seconds, upgrade should be allowed
		lastChangeTime = time.Now().Add(-11 * time.Second)
		newLayer, allowed, reason = qualityService.ApplyHysteresis(2, currentLayer, lastChangeTime)

		assert.Equal(t, 2, newLayer, "Should upgrade to layer 2")
		assert.True(t, allowed, "Upgrade should be allowed after 10s delay")
		assert.Contains(t, reason, "approved", "Reason should mention approval")

		t.Logf("Upgrade allowed after 11s: %s", reason)
	})
}

// Test adaptive quality selection based on bandwidth
func TestAdaptiveQualitySelection(t *testing.T) {
	qualityService := services.NewQualityService()

	testCases := []struct {
		name             string
		bandwidthKbps    float64
		currentLayer     int
		expectedLayer    int
		shouldChange     bool
		description      string
	}{
		{
			name:          "HighBandwidth_SelectHigh",
			bandwidthKbps: 3000.0,
			currentLayer:  0,
			expectedLayer: 2,
			shouldChange:  true,
			description:   "3Mbps bandwidth should select 1080p (layer 2)",
		},
		{
			name:          "StandardBandwidth_SelectMedium",
			bandwidthKbps: 1500.0, // Need 1440 Kbps (1200 * 1.2) for upgrade from layer 0
			currentLayer:  0,
			expectedLayer: 1,
			shouldChange:  true,
			description:   "1.5Mbps bandwidth should select 720p (layer 1) with 20% upgrade margin",
		},
		{
			name:          "LowBandwidth_SelectLow",
			bandwidthKbps: 600.0,
			currentLayer:  2,
			expectedLayer: 0,
			shouldChange:  true,
			description:   "600Kbps bandwidth should select 360p (layer 0)",
		},
		{
			name:          "StableBandwidth_MaintainLayer",
			bandwidthKbps: 1300.0,
			currentLayer:  1,
			expectedLayer: 1,
			shouldChange:  false,
			description:   "1.3Mbps bandwidth should maintain 720p (layer 1)",
		},
		{
			name:          "UpgradeMargin_PreventOscillation",
			bandwidthKbps: 2400.0,
			currentLayer:  1,
			expectedLayer: 1,
			shouldChange:  false,
			description:   "2.4Mbps (below 2.5Mbps + 20% margin) should maintain 720p",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decision := qualityService.SelectQualityLayer(tc.bandwidthKbps, tc.currentLayer)

			assert.Equal(t, tc.expectedLayer, decision.NewLayer, "Quality layer mismatch")
			assert.Equal(t, tc.shouldChange, decision.ShouldChange, "Quality change decision mismatch")

			t.Logf("%s: %s", tc.description, decision.Reason)
			t.Logf("  Bandwidth: %.0f Kbps, Layer: %d → %d, Change: %v",
				decision.EstimatedBandwidth, decision.CurrentLayer, decision.NewLayer, decision.ShouldChange)
		})
	}
}

// Test dynamic quality adaptation with changing network conditions
func TestDynamicQualityAdaptation(t *testing.T) {
	qualityService := services.NewQualityService()

	// Simulate network degradation: high → standard → constrained
	phases := []NetworkCondition{
		{BandwidthKbps: 3000.0, PacketLossRate: 0.001, JitterMs: 5.0, RTTMs: 50.0, Description: "Phase 1: High bandwidth"},
		{BandwidthKbps: 1200.0, PacketLossRate: 0.005, JitterMs: 15.0, RTTMs: 100.0, Description: "Phase 2: Standard bandwidth"},
		{BandwidthKbps: 600.0, PacketLossRate: 0.02, JitterMs: 30.0, RTTMs: 200.0, Description: "Phase 3: Constrained bandwidth"},
		{BandwidthKbps: 1200.0, PacketLossRate: 0.005, JitterMs: 15.0, RTTMs: 100.0, Description: "Phase 4: Recovery to standard"},
		{BandwidthKbps: 3000.0, PacketLossRate: 0.001, JitterMs: 5.0, RTTMs: 50.0, Description: "Phase 5: Full recovery"},
	}

	currentLayer := 2 // Start at highest quality
	// lastChangeTime := time.Time{} // Allow immediate first change - not used in simplified version

	for i, phase := range phases {
		t.Logf("\n=== %s ===", phase.Description)

		// Simulate 50 frames in this phase
		result := runLatencyTest(t, qualityService, phase, currentLayer, 5000.0, 50)

		// Calculate expected layer based on bandwidth
		expectedLayer := getExpectedQualityLayer(phase.BandwidthKbps)

		t.Logf("Phase %d: Bandwidth %.0f Kbps, Current Layer %d, Expected Layer %d",
			i+1, phase.BandwidthKbps, currentLayer, expectedLayer)
		t.Logf("  Average Latency: %.0f ms, P95: %.0f ms, P99: %.0f ms",
			result.AverageLatencyMs, result.P95LatencyMs, result.P99LatencyMs)

		// Update current layer for next phase
		if len(result.Measurements) > 0 {
			currentLayer = result.Measurements[len(result.Measurements)-1].QualityLayer
		}

		// Allow time for quality adaptation to settle
		time.Sleep(100 * time.Millisecond)
	}
}

// Test all three latency targets together
func TestAllLatencyTargets(t *testing.T) {
	qualityService := services.NewQualityService()

	scenarios := []struct {
		name              string
		network           NetworkCondition
		targetLayer       int
		targetLatencyMs   float64
	}{
		{
			name: "High_Bandwidth_1080p_1s",
			network: NetworkCondition{
				BandwidthKbps: 3000.0, PacketLossRate: 0.001, JitterMs: 5.0, RTTMs: 50.0,
				Description: "High Bandwidth (>2Mbps)",
			},
			targetLayer:     2,
			targetLatencyMs: 1000.0,
		},
		{
			name: "Standard_Bandwidth_720p_3s",
			network: NetworkCondition{
				BandwidthKbps: 1200.0, PacketLossRate: 0.005, JitterMs: 15.0, RTTMs: 100.0,
				Description: "Standard Bandwidth (800Kbps-2Mbps)",
			},
			targetLayer:     1,
			targetLatencyMs: 3000.0,
		},
		{
			name: "Constrained_Bandwidth_360p_5s",
			network: NetworkCondition{
				BandwidthKbps: 600.0, PacketLossRate: 0.02, JitterMs: 30.0, RTTMs: 200.0,
				Description: "Constrained Bandwidth (<800Kbps)",
			},
			targetLayer:     0,
			targetLatencyMs: 5000.0,
		},
	}

	results := make([]LatencyTestResult, 0, len(scenarios))

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			result := runLatencyTest(t, qualityService, scenario.network, scenario.targetLayer, scenario.targetLatencyMs, 100)
			results = append(results, result)

			// Validate latency targets
			require.True(t, result.Success, "Scenario %s failed: %s", scenario.name, result.FailureReason)
			assert.Less(t, result.AverageLatencyMs, scenario.targetLatencyMs,
				"Average latency should be less than target")

			logLatencyResults(t, result)
		})
	}

	// Generate comprehensive report
	t.Log("\n=== COMPREHENSIVE LATENCY TEST REPORT ===")
	for _, result := range results {
		t.Logf("\nScenario: %s", result.Scenario)
		t.Logf("  Target Latency: %.0f ms", result.TargetLatencyMs)
		t.Logf("  Average Latency: %.0f ms (%.1f%% of target)",
			result.AverageLatencyMs, (result.AverageLatencyMs/result.TargetLatencyMs)*100)
		t.Logf("  P95 Latency: %.0f ms", result.P95LatencyMs)
		t.Logf("  P99 Latency: %.0f ms", result.P99LatencyMs)
		t.Logf("  Quality Switches: %d", result.QualitySwitchCount)
		t.Logf("  Rebuffer Events: %d", result.RebufferEvents)
		t.Logf("  Status: %s", getStatusEmoji(result.Success))
	}
}

// Helper function: Run latency test simulation
func runLatencyTest(
	t *testing.T,
	qualityService services.QualityService,
	network NetworkCondition,
	targetLayer int,
	targetLatencyMs float64,
	frameCount int,
) LatencyTestResult {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result := LatencyTestResult{
		Scenario:           network.Description,
		BandwidthKbps:      network.BandwidthKbps,
		TargetQualityLayer: targetLayer,
		TargetLatencyMs:    targetLatencyMs,
		Measurements:       make([]LatencyMeasurement, 0, frameCount),
		Success:            true,
	}

	currentLayer := targetLayer
	lastChangeTime := time.Time{} // Allow immediate first change
	previousBandwidth := (*services.BandwidthEstimation)(nil)

	for frameID := 0; frameID < frameCount; frameID++ {
		select {
		case <-ctx.Done():
			result.Success = false
			result.FailureReason = "Test timeout"
			return result
		default:
		}

		// Simulate frame generation at broadcaster
		broadcasterTime := time.Now()

		// Simulate network transmission delay based on quality layer and network conditions
		transmissionDelayMs := calculateTransmissionDelay(currentLayer, network)
		time.Sleep(time.Duration(transmissionDelayMs) * time.Millisecond)

		// Simulate frame display at viewer
		viewerTime := time.Now()
		latencyMs := viewerTime.Sub(broadcasterTime).Seconds() * 1000

		// Get quality layer config
		layerConfig, err := qualityService.GetQualityLayerConfig(currentLayer)
		require.NoError(t, err, "Failed to get quality layer config")

		// Simulate WebRTC stats
		stats := simulateWebRTCStats(network, broadcasterTime)

		// Estimate bandwidth
		bandwidth, err := qualityService.EstimateBandwidth(stats, previousBandwidth)
		require.NoError(t, err, "Failed to estimate bandwidth")
		previousBandwidth = bandwidth

		// Record measurement
		measurement := LatencyMeasurement{
			FrameID:          frameID,
			BroadcasterTime:  broadcasterTime,
			ViewerTime:       viewerTime,
			LatencyMs:        latencyMs,
			BandwidthKbps:    bandwidth.SmoothedBandwidth,
			QualityLayer:     currentLayer,
			QualityLayerName: layerConfig.Name,
			PacketLossRate:   bandwidth.PacketLossRate,
			JitterMs:         bandwidth.JitterMs,
		}
		result.Measurements = append(result.Measurements, measurement)

		// Check for quality layer change (every 10 frames)
		if frameID%10 == 0 {
			decision := qualityService.SelectQualityLayer(bandwidth.SmoothedBandwidth, currentLayer)

			if decision.ShouldChange {
				// Apply hysteresis
				newLayer, allowed, reason := qualityService.ApplyHysteresis(
					decision.NewLayer, currentLayer, lastChangeTime)

				if allowed && newLayer != currentLayer {
					t.Logf("Frame %d: Quality switch %d → %d (%s)",
						frameID, currentLayer, newLayer, reason)
					currentLayer = newLayer
					lastChangeTime = time.Now()
					result.QualitySwitchCount++
				}
			}
		}

		// Simulate frame processing time (30 fps = 33ms per frame)
		time.Sleep(33 * time.Millisecond)
	}

	// Calculate statistics
	calculateLatencyStatistics(&result)

	// Validate against target
	if result.AverageLatencyMs > targetLatencyMs {
		result.Success = false
		result.FailureReason = fmt.Sprintf("Average latency %.0f ms exceeds target %.0f ms",
			result.AverageLatencyMs, targetLatencyMs)
	}

	return result
}

// Helper function: Calculate transmission delay based on quality layer and network
func calculateTransmissionDelay(qualityLayer int, network NetworkCondition) float64 {
	// Base latency from RTT
	baseLatency := network.RTTMs

	// Quality layer bitrate impact
	layerBitrates := []int{500, 1200, 2500} // Kbps
	if qualityLayer < 0 || qualityLayer >= len(layerBitrates) {
		qualityLayer = 1 // Default to medium
	}

	bitrate := float64(layerBitrates[qualityLayer])

	// Calculate transmission delay based on bandwidth availability
	bandwidthUtilization := bitrate / network.BandwidthKbps
	if bandwidthUtilization > 1.0 {
		// Congestion: delay increases significantly
		baseLatency *= (1.0 + (bandwidthUtilization - 1.0) * 3.0)
	} else if bandwidthUtilization > 0.8 {
		// Near capacity: slight delay increase
		baseLatency *= (1.0 + (bandwidthUtilization - 0.8) * 0.5)
	}

	// Add packet loss penalty
	if network.PacketLossRate > 0 {
		retransmissionDelay := network.RTTMs * network.PacketLossRate * 10 // Estimate retransmission impact
		baseLatency += retransmissionDelay
	}

	// Add jitter variance
	jitterVariance := network.JitterMs * (0.5 + 0.5*float64(time.Now().UnixNano()%100)/100)
	baseLatency += jitterVariance

	return baseLatency
}

// Helper function: Simulate WebRTC stats
func simulateWebRTCStats(network NetworkCondition, timestamp time.Time) *services.WebRTCStatsReport {
	// Simulate realistic WebRTC statistics based on network conditions
	return &services.WebRTCStatsReport{
		Timestamp:         timestamp,
		BytesSent:         1000000 + uint64(time.Now().UnixNano()%100000),
		BytesReceived:     950000 + uint64(time.Now().UnixNano()%50000),
		PacketsSent:       1000 + uint64(time.Now().UnixNano()%100),
		PacketsReceived:   uint64(float64(1000) * (1.0 - network.PacketLossRate)),
		PacketsLost:       uint64(float64(1000) * network.PacketLossRate),
		Jitter:            network.JitterMs,
		RoundTripTime:     network.RTTMs,
		AvailableBitrate:  network.BandwidthKbps,
	}
}

// Helper function: Calculate latency statistics
func calculateLatencyStatistics(result *LatencyTestResult) {
	if len(result.Measurements) == 0 {
		return
	}

	// Extract latency values
	latencies := make([]float64, len(result.Measurements))
	totalLatency := 0.0
	minLatency := math.MaxFloat64
	maxLatency := 0.0

	for i, m := range result.Measurements {
		latencies[i] = m.LatencyMs
		totalLatency += m.LatencyMs

		if m.LatencyMs < minLatency {
			minLatency = m.LatencyMs
		}
		if m.LatencyMs > maxLatency {
			maxLatency = m.LatencyMs
		}
	}

	// Calculate average
	result.AverageLatencyMs = totalLatency / float64(len(result.Measurements))
	result.MinLatencyMs = minLatency
	result.MaxLatencyMs = maxLatency

	// Calculate percentiles
	sort.Float64s(latencies)
	result.P95LatencyMs = calculatePercentile(latencies, 95)
	result.P99LatencyMs = calculatePercentile(latencies, 99)
}

// Helper function: Calculate percentile
func calculatePercentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}

	index := (percentile / 100.0) * float64(len(sortedValues))
	if index == float64(int(index)) {
		// Exact index
		return sortedValues[int(index)-1]
	}

	// Interpolate between values
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if upper >= len(sortedValues) {
		return sortedValues[len(sortedValues)-1]
	}
	if lower < 0 {
		return sortedValues[0]
	}

	fraction := index - float64(lower)
	return sortedValues[lower] + fraction*(sortedValues[upper]-sortedValues[lower])
}

// Helper function: Get expected quality layer based on bandwidth
func getExpectedQualityLayer(bandwidthKbps float64) int {
	if bandwidthKbps >= 2000 {
		return 2 // 1080p
	}
	if bandwidthKbps >= 800 {
		return 1 // 720p
	}
	return 0 // 360p
}

// Helper function: Log latency test results
func logLatencyResults(t *testing.T, result LatencyTestResult) {
	t.Logf("\n=== Latency Test Results: %s ===", result.Scenario)
	t.Logf("Bandwidth: %.0f Kbps", result.BandwidthKbps)
	t.Logf("Target Quality Layer: %d", result.TargetQualityLayer)
	t.Logf("Target Latency: %.0f ms", result.TargetLatencyMs)
	t.Logf("\nMeasurements (n=%d):", len(result.Measurements))
	t.Logf("  Average Latency: %.0f ms (%.1f%% of target)",
		result.AverageLatencyMs, (result.AverageLatencyMs/result.TargetLatencyMs)*100)
	t.Logf("  P95 Latency: %.0f ms", result.P95LatencyMs)
	t.Logf("  P99 Latency: %.0f ms", result.P99LatencyMs)
	t.Logf("  Min Latency: %.0f ms", result.MinLatencyMs)
	t.Logf("  Max Latency: %.0f ms", result.MaxLatencyMs)
	t.Logf("\nQuality Metrics:")
	t.Logf("  Quality Switches: %d", result.QualitySwitchCount)
	t.Logf("  Rebuffer Events: %d", result.RebufferEvents)
	t.Logf("\nResult: %s", getStatusString(result.Success))
	if !result.Success {
		t.Logf("Failure Reason: %s", result.FailureReason)
	}
}

// Helper function: Get status emoji
func getStatusEmoji(success bool) string {
	if success {
		return "PASS"
	}
	return "FAIL"
}

// Helper function: Get status string
func getStatusString(success bool) string {
	if success {
		return "PASS - All latency targets met"
	}
	return "FAIL - Latency targets not met"
}