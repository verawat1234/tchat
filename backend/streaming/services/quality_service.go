package services

import (
	"fmt"
	"math"
	"time"
)

// QualityLayer represents adaptive bitrate quality tiers
type QualityLayer int

const (
	QualityLayerLow QualityLayer = iota
	QualityLayerMedium
	QualityLayerHigh
)

// QualityLayerConfig defines technical specifications for each quality layer
type QualityLayerConfig struct {
	Layer        QualityLayer
	Name         string
	Resolution   string
	Width        int
	Height       int
	Bitrate      int // Kbps
	Framerate    int
	MinBandwidth float64 // Kbps - minimum bandwidth required
	MaxBandwidth float64 // Kbps - maximum bandwidth threshold
}

// WebRTCStatsReport represents simplified WebRTC statistics
// This structure can be extended with actual Pion WebRTC stats when integrated
type WebRTCStatsReport struct {
	Timestamp        time.Time
	BytesSent        uint64
	BytesReceived    uint64
	PacketsSent      uint64
	PacketsReceived  uint64
	PacketsLost      uint64
	Jitter           float64 // milliseconds
	RoundTripTime    float64 // milliseconds
	AvailableBitrate float64 // Kbps
}

// BandwidthEstimation holds bandwidth estimation state
type BandwidthEstimation struct {
	CurrentBandwidth  float64 // Kbps
	SmoothedBandwidth float64 // Kbps (EMA smoothed)
	PacketLossRate    float64 // 0.0 - 1.0
	JitterMs          float64 // milliseconds
	RTTMs             float64 // milliseconds
	LastUpdateTime    time.Time
	EMAWeight         float64 // Exponential moving average weight (0.0 - 1.0)
}

// QualityChangeDecision represents quality layer change decision
type QualityChangeDecision struct {
	ShouldChange       bool
	NewLayer           int
	CurrentLayer       int
	Reason             string
	EstimatedBandwidth float64
	RequiredBandwidth  float64
	UpgradeMargin      float64 // 20% margin for upgrades
}

// QualityService manages adaptive bitrate quality selection
type QualityService interface {
	// EstimateBandwidth analyzes WebRTC stats to estimate available bandwidth
	EstimateBandwidth(stats *WebRTCStatsReport, previous *BandwidthEstimation) (*BandwidthEstimation, error)

	// SelectQualityLayer determines optimal quality layer based on bandwidth
	SelectQualityLayer(bandwidth float64, currentLayer int) QualityChangeDecision

	// ApplyHysteresis prevents quality thrashing with time-based delays
	ApplyHysteresis(newLayer, currentLayer int, lastChangeTime time.Time) (int, bool, string)

	// GetQualityLayerConfig retrieves configuration for specific quality layer
	GetQualityLayerConfig(layer int) (QualityLayerConfig, error)

	// GetAllQualityLayers returns all available quality layer configurations
	GetAllQualityLayers() []QualityLayerConfig

	// CalculateRequiredBandwidth estimates bandwidth needed for target layer
	CalculateRequiredBandwidth(layer int) (float64, error)
}

// qualityServiceImpl implements QualityService
type qualityServiceImpl struct {
	qualityLayers         []QualityLayerConfig
	upgradeDelaySeconds   int
	downgradeDelaySeconds int
	upgradeMarginPercent  float64
	emaWeight             float64
}

// NewQualityService creates a new adaptive bitrate quality service
func NewQualityService() QualityService {
	return &qualityServiceImpl{
		qualityLayers:         initializeQualityLayers(),
		upgradeDelaySeconds:   10,   // Conservative upgrade delay
		downgradeDelaySeconds: 5,    // Quick response to congestion
		upgradeMarginPercent:  0.20, // 20% bandwidth margin for upgrades
		emaWeight:             0.3,  // EMA smoothing weight
	}
}

// initializeQualityLayers defines the three quality tier configurations
func initializeQualityLayers() []QualityLayerConfig {
	return []QualityLayerConfig{
		{
			Layer:        QualityLayerLow,
			Name:         "Low Quality",
			Resolution:   "360p",
			Width:        640,
			Height:       360,
			Bitrate:      500,
			Framerate:    30,
			MinBandwidth: 0,
			MaxBandwidth: 800,
		},
		{
			Layer:        QualityLayerMedium,
			Name:         "Medium Quality",
			Resolution:   "720p",
			Width:        1280,
			Height:       720,
			Bitrate:      1200,
			Framerate:    30,
			MinBandwidth: 800,
			MaxBandwidth: 2000,
		},
		{
			Layer:        QualityLayerHigh,
			Name:         "High Quality",
			Resolution:   "1080p",
			Width:        1920,
			Height:       1080,
			Bitrate:      2500,
			Framerate:    30,
			MinBandwidth: 2000,
			MaxBandwidth: math.MaxFloat64,
		},
	}
}

// EstimateBandwidth analyzes WebRTC stats to estimate available bandwidth
func (s *qualityServiceImpl) EstimateBandwidth(stats *WebRTCStatsReport, previous *BandwidthEstimation) (*BandwidthEstimation, error) {
	if stats == nil {
		return nil, fmt.Errorf("stats report cannot be nil")
	}

	estimation := &BandwidthEstimation{
		LastUpdateTime: stats.Timestamp,
		EMAWeight:      s.emaWeight,
	}

	// Calculate packet loss rate
	if stats.PacketsSent > 0 {
		estimation.PacketLossRate = float64(stats.PacketsLost) / float64(stats.PacketsSent)
	}

	// Extract jitter and RTT
	estimation.JitterMs = stats.Jitter
	estimation.RTTMs = stats.RoundTripTime

	// Calculate instantaneous bandwidth if available from stats
	rawBandwidth := stats.AvailableBitrate
	if rawBandwidth <= 0 && previous != nil {
		// Fallback: estimate from bytes sent and time delta
		if !previous.LastUpdateTime.IsZero() {
			timeDelta := stats.Timestamp.Sub(previous.LastUpdateTime).Seconds()
			if timeDelta > 0 {
				// Convert bytes to kilobits per second
				rawBandwidth = (float64(stats.BytesSent) * 8) / (timeDelta * 1000)
			}
		}
	}

	// Apply network condition penalties
	rawBandwidth = s.applyNetworkPenalties(rawBandwidth, estimation)

	// Apply exponential moving average (EMA) for smoothing
	if previous != nil && previous.SmoothedBandwidth > 0 {
		// EMA formula: smoothed = (weight * new) + ((1 - weight) * previous)
		estimation.SmoothedBandwidth = (s.emaWeight * rawBandwidth) +
			((1 - s.emaWeight) * previous.SmoothedBandwidth)
	} else {
		estimation.SmoothedBandwidth = rawBandwidth
	}

	estimation.CurrentBandwidth = rawBandwidth

	return estimation, nil
}

// applyNetworkPenalties adjusts bandwidth estimate based on network conditions
func (s *qualityServiceImpl) applyNetworkPenalties(bandwidth float64, estimation *BandwidthEstimation) float64 {
	// Packet loss penalty: reduce bandwidth estimate if packet loss is high
	if estimation.PacketLossRate > 0.01 { // > 1% packet loss
		lossMultiplier := 1.0 - (estimation.PacketLossRate * 0.5)
		bandwidth *= lossMultiplier
	}

	// High jitter penalty: reduce bandwidth if jitter indicates network instability
	if estimation.JitterMs > 30 { // > 30ms jitter
		jitterPenalty := math.Min(estimation.JitterMs/100, 0.3) // Max 30% penalty
		bandwidth *= (1.0 - jitterPenalty)
	}

	// High RTT penalty: reduce bandwidth for high latency connections
	if estimation.RTTMs > 200 { // > 200ms RTT
		rttPenalty := math.Min(estimation.RTTMs/1000, 0.2) // Max 20% penalty
		bandwidth *= (1.0 - rttPenalty)
	}

	return bandwidth
}

// SelectQualityLayer determines optimal quality layer based on bandwidth
func (s *qualityServiceImpl) SelectQualityLayer(bandwidth float64, currentLayer int) QualityChangeDecision {
	decision := QualityChangeDecision{
		ShouldChange:       false,
		NewLayer:           currentLayer,
		CurrentLayer:       currentLayer,
		EstimatedBandwidth: bandwidth,
	}

	// Find the best quality layer for available bandwidth
	selectedLayer := 0
	for i := len(s.qualityLayers) - 1; i >= 0; i-- {
		layer := s.qualityLayers[i]

		// For upgrades, require 20% bandwidth margin to prevent oscillation
		requiredBandwidth := float64(layer.Bitrate)
		if i > currentLayer {
			requiredBandwidth *= (1.0 + s.upgradeMarginPercent)
			decision.UpgradeMargin = s.upgradeMarginPercent
		}

		if bandwidth >= requiredBandwidth {
			selectedLayer = i
			decision.RequiredBandwidth = requiredBandwidth
			break
		}
	}

	// Determine if layer should change
	if selectedLayer != currentLayer {
		decision.ShouldChange = true
		decision.NewLayer = selectedLayer

		if selectedLayer > currentLayer {
			decision.Reason = fmt.Sprintf("Upgrade to %s: sufficient bandwidth (%.0f Kbps > %.0f Kbps required)",
				s.qualityLayers[selectedLayer].Name, bandwidth, decision.RequiredBandwidth)
		} else {
			decision.Reason = fmt.Sprintf("Downgrade to %s: insufficient bandwidth (%.0f Kbps < %.0f Kbps required)",
				s.qualityLayers[selectedLayer].Name, bandwidth, float64(s.qualityLayers[currentLayer].Bitrate))
		}
	} else {
		decision.Reason = fmt.Sprintf("Maintain %s: bandwidth stable (%.0f Kbps)",
			s.qualityLayers[currentLayer].Name, bandwidth)
	}

	return decision
}

// ApplyHysteresis prevents quality thrashing with time-based delays
func (s *qualityServiceImpl) ApplyHysteresis(newLayer, currentLayer int, lastChangeTime time.Time) (int, bool, string) {
	// No change needed - skip hysteresis
	if newLayer == currentLayer {
		return currentLayer, false, "No quality change required"
	}

	// First quality change - allow immediately
	if lastChangeTime.IsZero() {
		reason := "Initial quality layer selection"
		return newLayer, true, reason
	}

	timeSinceLastChange := time.Since(lastChangeTime)

	// Upgrading: require longer delay (conservative)
	if newLayer > currentLayer {
		upgradeDelay := time.Duration(s.upgradeDelaySeconds) * time.Second
		if timeSinceLastChange < upgradeDelay {
			reason := fmt.Sprintf("Upgrade delayed: %.1fs remaining (%.1fs / %ds required)",
				(upgradeDelay - timeSinceLastChange).Seconds(),
				timeSinceLastChange.Seconds(),
				s.upgradeDelaySeconds)
			return currentLayer, false, reason
		}
		reason := fmt.Sprintf("Upgrade approved after %ds delay", s.upgradeDelaySeconds)
		return newLayer, true, reason
	}

	// Downgrading: require shorter delay (quick congestion response)
	downgradeDelay := time.Duration(s.downgradeDelaySeconds) * time.Second
	if timeSinceLastChange < downgradeDelay {
		reason := fmt.Sprintf("Downgrade delayed: %.1fs remaining (%.1fs / %ds required)",
			(downgradeDelay - timeSinceLastChange).Seconds(),
			timeSinceLastChange.Seconds(),
			s.downgradeDelaySeconds)
		return currentLayer, false, reason
	}

	reason := fmt.Sprintf("Downgrade approved after %ds delay", s.downgradeDelaySeconds)
	return newLayer, true, reason
}

// GetQualityLayerConfig retrieves configuration for specific quality layer
func (s *qualityServiceImpl) GetQualityLayerConfig(layer int) (QualityLayerConfig, error) {
	if layer < 0 || layer >= len(s.qualityLayers) {
		return QualityLayerConfig{}, fmt.Errorf("invalid quality layer: %d (valid range: 0-%d)",
			layer, len(s.qualityLayers)-1)
	}
	return s.qualityLayers[layer], nil
}

// GetAllQualityLayers returns all available quality layer configurations
func (s *qualityServiceImpl) GetAllQualityLayers() []QualityLayerConfig {
	// Return copy to prevent external modification
	layers := make([]QualityLayerConfig, len(s.qualityLayers))
	copy(layers, s.qualityLayers)
	return layers
}

// CalculateRequiredBandwidth estimates bandwidth needed for target layer
func (s *qualityServiceImpl) CalculateRequiredBandwidth(layer int) (float64, error) {
	config, err := s.GetQualityLayerConfig(layer)
	if err != nil {
		return 0, err
	}

	// Return bitrate plus 10% overhead for protocol headers and retransmissions
	return float64(config.Bitrate) * 1.1, nil
}
