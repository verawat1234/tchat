// backend/streaming/handlers/start_stream_handler.go
// Start Stream Handler - WebRTC negotiation and stream lifecycle management
// Implements POST /api/v1/streams/{streamId}/start endpoint

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"

	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

// StartStreamHandler handles POST /api/v1/streams/{streamId}/start requests
type StartStreamHandler struct {
	liveStreamRepo   repository.LiveStreamRepositoryInterface
	webrtcService    services.WebRTCService
	signalingService *services.SignalingService
	recordingService services.RecordingService
	kycService       services.KYCServiceInterface
}

// StartStreamRequest represents the request body for starting a stream
type StartStreamRequest struct {
	SDPOffer      string   `json:"sdp_offer" binding:"required"`
	QualityLayers []string `json:"quality_layers"`
}

// StartStreamResponse represents the success response structure
type StartStreamResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    StartStreamResponseData `json:"data"`
}

// StartStreamResponseData contains the WebRTC answer and stream information
type StartStreamResponseData struct {
	SDPAnswer        string              `json:"sdp_answer"`
	ICEServers       []ICEServerResponse `json:"ice_servers"`
	WebRTCSessionID  string              `json:"webrtc_session_id"`
	Stream           StreamInfoResponse  `json:"stream"`
}

// ICEServerResponse represents an ICE server configuration
type ICEServerResponse struct {
	URLs []string `json:"urls"`
}

// StreamInfoResponse contains updated stream information after starting
type StreamInfoResponse struct {
	ID                  string    `json:"id"`
	Status              string    `json:"status"`
	ActualStartTime     string    `json:"actual_start_time"`
	ViewerCount         int       `json:"viewer_count"`
	PeakViewerCount     int       `json:"peak_viewer_count"`
	StreamKey           string    `json:"stream_key"`
	WebRTCSessionID     string    `json:"webrtc_session_id"`
	BroadcasterID       string    `json:"broadcaster_id"`
	StreamType          string    `json:"stream_type"`
	Title               string    `json:"title"`
	PrivacySetting      string    `json:"privacy_setting"`
	ScheduledStartTime  string    `json:"scheduled_start_time"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Success        bool               `json:"success"`
	Error          string             `json:"error"`
	Details        map[string]string  `json:"details,omitempty"`
	CurrentStatus  string             `json:"current_status,omitempty"`
	StreamID       string             `json:"stream_id,omitempty"`
	BroadcasterID  string             `json:"broadcaster_id,omitempty"`
	ActualStartTime string            `json:"actual_start_time,omitempty"`
}

// NewStartStreamHandler creates a new start stream handler instance
func NewStartStreamHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	webrtcService services.WebRTCService,
	signalingService *services.SignalingService,
	recordingService services.RecordingService,
	kycService services.KYCServiceInterface,
) *StartStreamHandler {
	return &StartStreamHandler{
		liveStreamRepo:   liveStreamRepo,
		webrtcService:    webrtcService,
		signalingService: signalingService,
		recordingService: recordingService,
		kycService:       kycService,
	}
}

// Handle processes the start stream request
func (h *StartStreamHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract broadcaster ID from JWT
	broadcasterIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse user_id to UUID
	var broadcasterID uuid.UUID
	switch v := broadcasterIDRaw.(type) {
	case uuid.UUID:
		broadcasterID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error:   "Invalid user ID format",
			})
			return
		}
		broadcasterID = parsed
	default:
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID type",
		})
		return
	}

	// Parse stream ID from URL parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid stream ID format",
		})
		return
	}

	// Fetch stream from database
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Stream not found",
			StreamID: streamID.String(),
		})
		return
	}

	// Authorization check - only broadcaster can start their stream
	if stream.BroadcasterID != broadcasterID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "Only the broadcaster can start their stream",
			StreamID: streamID.String(),
			BroadcasterID: stream.BroadcasterID.String(),
		})
		return
	}

	// Status validation - stream must be in scheduled state
	if stream.Status == "live" {
		actualStartStr := ""
		if stream.ActualStartTime != nil {
			actualStartStr = stream.ActualStartTime.Format("2006-01-02T15:04:05Z07:00")
		}
		c.JSON(http.StatusConflict, ErrorResponse{
			Success: false,
			Error:   "Stream is already live",
			CurrentStatus: "live",
			StreamID: streamID.String(),
			ActualStartTime: actualStartStr,
		})
		return
	}

	if stream.Status == "ended" || stream.Status == "terminated" {
		c.JSON(http.StatusConflict, ErrorResponse{
			Success: false,
			Error:   "Stream has already ended",
			CurrentStatus: stream.Status,
			StreamID: streamID.String(),
		})
		return
	}

	// Parse request body
	var req StartStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Details: map[string]string{
				"validation_error": err.Error(),
			},
		})
		return
	}

	// Validate SDP offer format
	if !isValidSDP(req.SDPOffer) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid WebRTC SDP offer format",
			Details: map[string]string{
				"sdp_offer": "SDP offer must be valid Session Description Protocol format",
			},
		})
		return
	}

	// Create WebRTC offer session description
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  req.SDPOffer,
	}

	// Handle WebRTC offer and generate answer
	answer, err := h.webrtcService.HandleOffer(streamID, offer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to process WebRTC offer",
			Details: map[string]string{
				"webrtc_error": err.Error(),
			},
		})
		return
	}

	// Generate WebRTC session ID
	webrtcSessionID := fmt.Sprintf("wrtc_session_%s", uuid.New().String()[:12])

	// Update stream status to live
	now := time.Now().UTC()
	updates := map[string]interface{}{
		"status":             "live",
		"actual_start_time":  now,
		"webrtc_session_id":  webrtcSessionID,
		"updated_at":         now,
	}

	if err := h.liveStreamRepo.Update(ctx, streamID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to update stream status",
			Details: map[string]string{
				"database_error": err.Error(),
			},
		})
		return
	}

	// Start recording in background (non-blocking)
	go func() {
		if err := h.recordingService.StartRecording(streamID, nil); err != nil {
			fmt.Printf("Failed to start recording for stream %s: %v\n", streamID, err)
		}
	}()

	// Start KYC monitoring in background (non-blocking)
	go func() {
		monitorCtx := context.Background()
		if err := h.kycService.StartKYCMonitoring(monitorCtx, streamID, broadcasterID, stream.StreamType); err != nil {
			fmt.Printf("Failed to start KYC monitoring for stream %s: %v\n", streamID, err)
		}
	}()

	// Get ICE servers configuration
	iceServers := h.webrtcService.GetICEServers()
	iceServerResponses := make([]ICEServerResponse, len(iceServers))
	for i, server := range iceServers {
		iceServerResponses[i] = ICEServerResponse{
			URLs: server.URLs,
		}
	}

	// Format scheduled start time
	scheduledStartStr := ""
	if stream.ScheduledStartTime != nil {
		scheduledStartStr = stream.ScheduledStartTime.Format("2006-01-02T15:04:05Z07:00")
	}

	// Construct success response
	response := StartStreamResponse{
		Success: true,
		Message: "Stream started successfully",
		Data: StartStreamResponseData{
			SDPAnswer:       answer.SDP,
			ICEServers:      iceServerResponses,
			WebRTCSessionID: webrtcSessionID,
			Stream: StreamInfoResponse{
				ID:                  streamID.String(),
				Status:              "live",
				ActualStartTime:     now.Format("2006-01-02T15:04:05Z07:00"),
				ViewerCount:         0,
				PeakViewerCount:     0,
				StreamKey:           stream.StreamKey,
				WebRTCSessionID:     webrtcSessionID,
				BroadcasterID:       stream.BroadcasterID.String(),
				StreamType:          stream.StreamType,
				Title:               stream.Title,
				PrivacySetting:      stream.PrivacySetting,
				ScheduledStartTime:  scheduledStartStr,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// isValidSDP validates that the SDP string has the basic required structure
func isValidSDP(sdp string) bool {
	if len(sdp) == 0 {
		return false
	}

	// Check for v=0 header (required SDP field)
	if len(sdp) < 3 || sdp[:3] != "v=0" {
		return false
	}

	// Basic validation passed
	return true
}