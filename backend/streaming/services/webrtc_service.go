package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v3"
)

// WebRTCService provides WebRTC peer connection management for live streaming
type WebRTCService interface {
	// CreatePeerConnection creates a new WebRTC peer connection for a stream
	CreatePeerConnection(streamID uuid.UUID) (*webrtc.PeerConnection, error)

	// HandleOffer processes a WebRTC offer and returns an answer
	HandleOffer(streamID uuid.UUID, offer webrtc.SessionDescription) (webrtc.SessionDescription, error)

	// AddSimulcastTracks configures simulcast tracks with 3 quality layers
	AddSimulcastTracks(peerConnection *webrtc.PeerConnection) error

	// GetICEServers returns configured STUN/TURN servers
	GetICEServers() []webrtc.ICEServer

	// ClosePeerConnection gracefully closes a peer connection and cleans up resources
	ClosePeerConnection(streamID uuid.UUID) error

	// GetPeerConnection retrieves an existing peer connection by stream ID
	GetPeerConnection(streamID uuid.UUID) (*webrtc.PeerConnection, bool)
}

// webRTCServiceImpl is the concrete implementation of WebRTCService
type webRTCServiceImpl struct {
	mu              sync.RWMutex
	peerConnections map[uuid.UUID]*webrtc.PeerConnection
	api             *webrtc.API
	config          WebRTCConfig
}

// WebRTCConfig holds configuration for the WebRTC service
type WebRTCConfig struct {
	// ICE server configuration
	STUNServers []string
	TURNServers []TURNServer

	// Simulcast configuration
	EnableSimulcast bool
	SimulcastLayers []SimulcastLayer

	// Connection settings
	ICEConnectionTimeout time.Duration
	KeepAliveInterval    time.Duration
}

// TURNServer represents a TURN server configuration
type TURNServer struct {
	URLs       []string
	Username   string
	Credential string
}

// SimulcastLayer represents a quality layer for simulcast streaming
type SimulcastLayer struct {
	RID        string // Restriction Identifier (low, mid, high)
	Width      int
	Height     int
	Bitrate    int // bits per second
	Framerate  int
}

// DefaultSimulcastLayers returns the standard 3-layer simulcast configuration
// Based on research.md lines 69-72
func DefaultSimulcastLayers() []SimulcastLayer {
	return []SimulcastLayer{
		{
			RID:       "low",
			Width:     640,
			Height:    360,
			Bitrate:   500000, // 500 Kbps for constrained networks
			Framerate: 30,
		},
		{
			RID:       "mid",
			Width:     1280,
			Height:    720,
			Bitrate:   1200000, // 1.2 Mbps for standard quality
			Framerate: 30,
		},
		{
			RID:       "high",
			Width:     1920,
			Height:    1080,
			Bitrate:   2500000, // 2.5 Mbps for high bandwidth
			Framerate: 30,
		},
	}
}

// NewWebRTCService creates a new WebRTC service instance
func NewWebRTCService(config WebRTCConfig) (WebRTCService, error) {
	// Set default configuration if not provided
	if len(config.STUNServers) == 0 {
		config.STUNServers = []string{"stun:stun.l.google.com:19302"}
	}

	if config.SimulcastLayers == nil {
		config.EnableSimulcast = true
		config.SimulcastLayers = DefaultSimulcastLayers()
	}

	if config.ICEConnectionTimeout == 0 {
		config.ICEConnectionTimeout = 30 * time.Second
	}

	if config.KeepAliveInterval == 0 {
		config.KeepAliveInterval = 30 * time.Second
	}

	// Create media engine with codec support
	mediaEngine := &webrtc.MediaEngine{}

	// Register VP8 codec for video (widely supported)
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeVP8,
			ClockRate:    90000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, fmt.Errorf("failed to register VP8 codec: %w", err)
	}

	// Register H.264 codec for video (hardware acceleration support)
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:    webrtc.MimeTypeH264,
			ClockRate:   90000,
			Channels:    0,
			SDPFmtpLine: "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f",
		},
		PayloadType: 102,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, fmt.Errorf("failed to register H.264 codec: %w", err)
	}

	// Register Opus codec for audio
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeOpus,
			ClockRate: 48000,
			Channels:  2,
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		return nil, fmt.Errorf("failed to register Opus codec: %w", err)
	}

	// Create interceptor registry for RTP/RTCP processing
	interceptorRegistry := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
		return nil, fmt.Errorf("failed to register interceptors: %w", err)
	}

	// Create WebRTC API instance
	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(mediaEngine),
		webrtc.WithInterceptorRegistry(interceptorRegistry),
	)

	return &webRTCServiceImpl{
		peerConnections: make(map[uuid.UUID]*webrtc.PeerConnection),
		api:             api,
		config:          config,
	}, nil
}

// CreatePeerConnection creates a new WebRTC peer connection for a stream
func (s *webRTCServiceImpl) CreatePeerConnection(streamID uuid.UUID) (*webrtc.PeerConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if peer connection already exists
	if _, exists := s.peerConnections[streamID]; exists {
		return nil, errors.New("peer connection already exists for this stream")
	}

	// Configure peer connection
	config := webrtc.Configuration{
		ICEServers: s.GetICEServers(),
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
	}

	// Create peer connection
	peerConnection, err := s.api.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Set up ICE connection state handler
	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Printf("Stream %s: ICE connection state changed to %s", streamID, state.String())

		// Handle disconnection
		if state == webrtc.ICEConnectionStateFailed || state == webrtc.ICEConnectionStateClosed {
			log.Printf("Stream %s: ICE connection failed or closed, cleaning up", streamID)
			if err := s.ClosePeerConnection(streamID); err != nil {
				log.Printf("Error closing peer connection: %v", err)
			}
		}
	})

	// Set up peer connection state handler
	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Stream %s: Peer connection state changed to %s", streamID, state.String())
	})

	// Set up ICE candidate handler
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			log.Printf("Stream %s: New ICE candidate: %s", streamID, candidate.String())
		}
	})

	// Store peer connection
	s.peerConnections[streamID] = peerConnection

	log.Printf("Created peer connection for stream %s", streamID)
	return peerConnection, nil
}

// HandleOffer processes a WebRTC offer and returns an answer
func (s *webRTCServiceImpl) HandleOffer(streamID uuid.UUID, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
	// Get or create peer connection
	pc, exists := s.GetPeerConnection(streamID)
	if !exists {
		var err error
		pc, err = s.CreatePeerConnection(streamID)
		if err != nil {
			return webrtc.SessionDescription{}, fmt.Errorf("failed to create peer connection: %w", err)
		}
	}

	// Set remote description (offer)
	if err := pc.SetRemoteDescription(offer); err != nil {
		return webrtc.SessionDescription{}, fmt.Errorf("failed to set remote description: %w", err)
	}

	// Configure simulcast if enabled
	if s.config.EnableSimulcast {
		if err := s.AddSimulcastTracks(pc); err != nil {
			log.Printf("Warning: Failed to add simulcast tracks: %v", err)
			// Continue without simulcast
		}
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, fmt.Errorf("failed to create answer: %w", err)
	}

	// Set local description (answer)
	if err := pc.SetLocalDescription(answer); err != nil {
		return webrtc.SessionDescription{}, fmt.Errorf("failed to set local description: %w", err)
	}

	log.Printf("Created answer for stream %s", streamID)
	return answer, nil
}

// AddSimulcastTracks configures simulcast tracks with 3 quality layers
func (s *webRTCServiceImpl) AddSimulcastTracks(peerConnection *webrtc.PeerConnection) error {
	if !s.config.EnableSimulcast {
		return errors.New("simulcast is not enabled")
	}

	// Create RTP transceivers for each simulcast layer
	for _, layer := range s.config.SimulcastLayers {
		// Create track for this layer
		track, err := webrtc.NewTrackLocalStaticRTP(
			webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8},
			fmt.Sprintf("video-%s", layer.RID),
			fmt.Sprintf("stream-%s", layer.RID),
		)
		if err != nil {
			return fmt.Errorf("failed to create track for layer %s: %w", layer.RID, err)
		}

		// Add transceiver with simulcast configuration
		rtpSender, err := peerConnection.AddTrack(track)
		if err != nil {
			return fmt.Errorf("failed to add track for layer %s: %w", layer.RID, err)
		}

		// Handle RTCP packets for this sender
		go s.processRTCP(rtpSender, layer.RID)

		log.Printf("Added simulcast track: %s (%dx%d @ %d bps, %d fps)",
			layer.RID, layer.Width, layer.Height, layer.Bitrate, layer.Framerate)
	}

	return nil
}

// processRTCP handles RTCP packets for a track
func (s *webRTCServiceImpl) processRTCP(rtpSender *webrtc.RTPSender, rid string) {
	rtcpBuf := make([]byte, 1500)
	for {
		_, _, err := rtpSender.Read(rtcpBuf)
		if err != nil {
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, io.EOF) {
				log.Printf("RTCP reader closed for RID %s", rid)
				return
			}
			log.Printf("Error reading RTCP for RID %s: %v", rid, err)
			continue
		}
		// RTCP packet processing would go here (e.g., bandwidth estimation, packet loss)
		// For now, we just log that we received RTCP packets
	}
}

// GetICEServers returns configured STUN/TURN servers
func (s *webRTCServiceImpl) GetICEServers() []webrtc.ICEServer {
	servers := make([]webrtc.ICEServer, 0)

	// Add STUN servers
	for _, stunURL := range s.config.STUNServers {
		servers = append(servers, webrtc.ICEServer{
			URLs: []string{stunURL},
		})
	}

	// Add TURN servers with credentials
	for _, turnServer := range s.config.TURNServers {
		servers = append(servers, webrtc.ICEServer{
			URLs:       turnServer.URLs,
			Username:   turnServer.Username,
			Credential: turnServer.Credential,
		})
	}

	return servers
}

// ClosePeerConnection gracefully closes a peer connection and cleans up resources
func (s *webRTCServiceImpl) ClosePeerConnection(streamID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pc, exists := s.peerConnections[streamID]
	if !exists {
		return fmt.Errorf("peer connection not found for stream %s", streamID)
	}

	// Close peer connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- pc.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error closing peer connection for stream %s: %v", streamID, err)
		}
	case <-ctx.Done():
		log.Printf("Timeout closing peer connection for stream %s", streamID)
	}

	// Remove from map
	delete(s.peerConnections, streamID)

	log.Printf("Closed peer connection for stream %s", streamID)
	return nil
}

// GetPeerConnection retrieves an existing peer connection by stream ID
func (s *webRTCServiceImpl) GetPeerConnection(streamID uuid.UUID) (*webrtc.PeerConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pc, exists := s.peerConnections[streamID]
	return pc, exists
}

// NewWebRTCServiceFromEnv creates a WebRTC service using environment variables
func NewWebRTCServiceFromEnv() (WebRTCService, error) {
	config := WebRTCConfig{
		EnableSimulcast:  true,
		SimulcastLayers:  DefaultSimulcastLayers(),
		STUNServers:      []string{"stun:stun.l.google.com:19302"},
		TURNServers:      parseTURNServersFromEnv(),
		ICEConnectionTimeout: 30 * time.Second,
		KeepAliveInterval:    30 * time.Second,
	}

	// Override STUN servers from environment if provided
	if stunServers := os.Getenv("WEBRTC_STUN_SERVERS"); stunServers != "" {
		config.STUNServers = []string{stunServers}
	}

	return NewWebRTCService(config)
}

// parseTURNServersFromEnv parses TURN server configuration from environment variables
func parseTURNServersFromEnv() []TURNServer {
	turnURLs := os.Getenv("WEBRTC_TURN_URLS")
	turnUsername := os.Getenv("WEBRTC_TURN_USERNAME")
	turnCredential := os.Getenv("WEBRTC_TURN_CREDENTIAL")

	if turnURLs == "" {
		return nil
	}

	return []TURNServer{
		{
			URLs:       []string{turnURLs},
			Username:   turnUsername,
			Credential: turnCredential,
		},
	}
}