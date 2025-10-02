package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// SyncHandler handles cross-platform playback synchronization endpoints
type SyncHandler struct {
	syncService  *services.SyncService
	videoService *services.VideoService
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(syncService *services.SyncService, videoService *services.VideoService) *SyncHandler {
	return &SyncHandler{
		syncService:  syncService,
		videoService: videoService,
	}
}

// RegisterRoutes registers sync routes with the Gin router
func (h *SyncHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/videos/:id/sync", h.SyncPlaybackPosition)
	router.POST("/videos/:id/sync/session", h.CreateSyncSession)
	router.GET("/videos/:id/sync/status", h.GetSyncStatus)
	router.POST("/videos/:id/sync/conflict", h.ResolveConflict)
	router.DELETE("/videos/:id/sync/session/:sessionId", h.TerminateSyncSession)
	router.GET("/videos/:id/sync/history", h.GetSyncHistory)
	router.POST("/videos/:id/sync/force", h.ForceSyncUpdate)
}

// SyncRequest represents playback sync request
type SyncRequest struct {
	SessionID   uuid.UUID           `json:"session_id" binding:"required"`
	Position    int                 `json:"position" binding:"required"`
	Quality     string              `json:"quality"`
	Platform    models.PlatformType `json:"platform" binding:"required"`
	PlaybackState string            `json:"playback_state"` // playing, paused, buffering
	BufferHealth  int               `json:"buffer_health"`
	Timestamp   time.Time           `json:"timestamp"`
}

// SyncResponse represents sync operation response
type SyncResponse struct {
	Success         bool                `json:"success"`
	SessionID       uuid.UUID           `json:"session_id"`
	UpdatedPosition int                 `json:"updated_position"`
	SyncedPlatforms []string            `json:"synced_platforms"`
	SyncLatency     int                 `json:"sync_latency_ms"`
	ConflictDetected bool               `json:"conflict_detected"`
	Message         string              `json:"message"`
}

// CreateSessionRequest represents sync session creation request
type CreateSessionRequest struct {
	VideoID     uuid.UUID           `json:"video_id" binding:"required"`
	UserID      uuid.UUID           `json:"user_id" binding:"required"`
	Platform    models.PlatformType `json:"platform" binding:"required"`
	InitialPosition int             `json:"initial_position"`
	Quality     string              `json:"quality"`
}

// SyncStatusResponse represents sync status information
type SyncStatusResponse struct {
	VideoID          uuid.UUID           `json:"video_id"`
	SessionID        uuid.UUID           `json:"session_id"`
	ActivePlatforms  []PlatformSyncInfo  `json:"active_platforms"`
	LastSyncTime     time.Time           `json:"last_sync_time"`
	SyncInterval     int                 `json:"sync_interval_ms"`
	ConflictCount    int                 `json:"conflict_count"`
	TotalSyncEvents  int                 `json:"total_sync_events"`
}

// PlatformSyncInfo represents sync information for a platform
type PlatformSyncInfo struct {
	Platform        models.PlatformType `json:"platform"`
	Position        int                 `json:"position"`
	Quality         string              `json:"quality"`
	PlaybackState   string              `json:"playback_state"`
	LastUpdate      time.Time           `json:"last_update"`
	SyncLatency     int                 `json:"sync_latency_ms"`
	IsHealthy       bool                `json:"is_healthy"`
}

// ConflictResolutionRequest represents conflict resolution request
type ConflictResolutionRequest struct {
	SessionID          uuid.UUID           `json:"session_id" binding:"required"`
	ResolutionStrategy string              `json:"resolution_strategy" binding:"required"`
	AuthorityPlatform  models.PlatformType `json:"authority_platform"`
}

// SyncHistoryResponse represents sync history data
type SyncHistoryResponse struct {
	VideoID     uuid.UUID       `json:"video_id"`
	SessionID   uuid.UUID       `json:"session_id"`
	SyncEvents  []SyncEvent     `json:"sync_events"`
	TotalEvents int             `json:"total_events"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
}

// SyncEvent represents a single sync event
type SyncEvent struct {
	EventID       uuid.UUID           `json:"event_id"`
	Platform      models.PlatformType `json:"platform"`
	Position      int                 `json:"position"`
	EventType     string              `json:"event_type"` // position_update, quality_change, conflict
	Timestamp     time.Time           `json:"timestamp"`
	SyncLatency   int                 `json:"sync_latency_ms"`
	ConflictInfo  *ConflictInfo       `json:"conflict_info,omitempty"`
}

// ConflictInfo represents conflict details
type ConflictInfo struct {
	ConflictType       string    `json:"conflict_type"`
	ConflictingPlatforms []string `json:"conflicting_platforms"`
	ResolutionStrategy string    `json:"resolution_strategy"`
	ResolvedPosition   int       `json:"resolved_position"`
}

// SyncPlaybackPosition synchronizes playback position across platforms
// POST /api/v1/videos/:id/sync
func (h *SyncHandler) SyncPlaybackPosition(c *gin.Context) {
	startTime := time.Now()

	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse request body
	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate video exists
	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Perform sync operation
	err = h.syncService.SyncPlaybackPosition(
		c.Request.Context(),
		req.SessionID,
		req.Position,
		req.Platform,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Sync failed: %v", err),
		})
		return
	}

	// Get updated sync states
	syncStates, err := h.syncService.GetSyncStatesBySession(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sync states"})
		return
	}

	// Build list of synced platforms
	syncedPlatforms := make([]string, 0)
	conflictDetected := false
	for _, state := range syncStates {
		syncedPlatforms = append(syncedPlatforms, string(state.PlatformType))
		if state.ConflictState.HasConflicts {
			conflictDetected = true
		}
	}

	// Calculate sync latency
	syncLatency := int(time.Since(startTime).Milliseconds())

	response := SyncResponse{
		Success:         true,
		SessionID:       req.SessionID,
		UpdatedPosition: req.Position,
		SyncedPlatforms: syncedPlatforms,
		SyncLatency:     syncLatency,
		ConflictDetected: conflictDetected,
		Message:         "Playback position synchronized successfully",
	}

	c.JSON(http.StatusOK, response)
}

// CreateSyncSession creates a new cross-platform sync session
// POST /api/v1/videos/:id/sync/session
func (h *SyncHandler) CreateSyncSession(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse request body
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate video ID matches
	if req.VideoID != videoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID mismatch"})
		return
	}

	// Validate video exists
	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Create sync session
// 	sessionID, err := h.syncService.InitializeSync(
// 		c.Request.Context(),
// 		videoID,
// 		req.UserID,
// 		req.Platform,
// 		req.InitialPosition,
// 	)
	// Initialize sync state
	deviceID := req.UserID // Use UserID as deviceID for now
	syncState, err := h.syncService.InitializeSync(c.Request.Context(), videoID, deviceID.String(), req.Platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create sync session: %v", err)})
		return
	}
	sessionID := syncState.SessionID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create sync session: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id":       sessionID,
		"video_id":         videoID,
		"user_id":          req.UserID,
		"platform":         string(req.Platform),
		"initial_position": req.InitialPosition,
		"created_at":       time.Now(),
		"message":          "Sync session created successfully",
	})
}

// GetSyncStatus retrieves current sync status for a video session
// GET /api/v1/videos/:id/sync/status
func (h *SyncHandler) GetSyncStatus(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get session ID from query
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	// Get sync states for session
	syncStates, err := h.syncService.GetSyncStatesBySession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sync session not found"})
		return
	}

	// Build platform sync info
	activePlatforms := make([]PlatformSyncInfo, 0)
	var lastSyncTime time.Time
	conflictCount := 0
	totalSyncEvents := 0

	for _, state := range syncStates {
		platformInfo := PlatformSyncInfo{
			Platform:      state.PlatformType,
			Position:      state.SyncData.Position,
			Quality:       state.SyncData.QualitySetting,
			PlaybackState: playbackState(state.SyncData.IsPlaying),
			LastUpdate:    *state.LastSyncTime,
			SyncLatency:   state.SyncMetrics.LastSyncLatency,
			IsHealthy:     state.SyncMetrics.SuccessfulSyncs > state.SyncMetrics.FailedSyncs,
		}
		activePlatforms = append(activePlatforms, platformInfo)

		if state.LastSyncTime.After(lastSyncTime) {
			lastSyncTime = *state.LastSyncTime
		}

		if state.ConflictState.HasConflicts {
			conflictCount++
		}

		totalSyncEvents += state.SyncMetrics.TotalSyncAttempts
	}

	response := SyncStatusResponse{
		VideoID:         videoID,
		SessionID:       sessionID,
		ActivePlatforms: activePlatforms,
		LastSyncTime:    lastSyncTime,
		SyncInterval:    100, // Target: <100ms sync latency
		ConflictCount:   conflictCount,
		TotalSyncEvents: totalSyncEvents,
	}

	c.JSON(http.StatusOK, response)
}

// ResolveConflict manually resolves a sync conflict
// POST /api/v1/videos/:id/sync/conflict
func (h *SyncHandler) ResolveConflict(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse request body
	var req ConflictResolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate resolution strategy
	validStrategies := map[string]models.ResolutionStrategy{
		"latest":    models.ResolutionLatest,
		"authority": models.ResolutionAuthority,
		"average":   models.ResolutionAverage,
		"manual":    models.ResolutionManual,
	}

	strategy, valid := validStrategies[req.ResolutionStrategy]
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid resolution strategy",
			"valid_strategies": []string{"latest", "authority", "average", "manual"},
		})
		return
	}

	// Resolve conflict
	err = h.syncService.ResolveConflict(
		c.Request.Context(),
		req.SessionID,
		strategy,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Conflict resolution failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"video_id":            videoID,
		"session_id":          req.SessionID,
		"resolution_strategy": req.ResolutionStrategy,
		"resolved_at":         time.Now(),
		"message":             "Conflict resolved successfully",
	})
}

// TerminateSyncSession ends a cross-platform sync session
// DELETE /api/v1/videos/:id/sync/session/:sessionId
func (h *SyncHandler) TerminateSyncSession(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse session ID
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	// Terminate sync session
	err = h.syncService.SuspendPlatformSync(c.Request.Context(), sessionID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to terminate session: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"video_id":     videoID,
		"session_id":   sessionID,
		"terminated_at": time.Now(),
		"message":      "Sync session terminated successfully",
	})
}

// GetSyncHistory retrieves sync event history for a session
// GET /api/v1/videos/:id/sync/history
func (h *SyncHandler) GetSyncHistory(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get session ID from query
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	// Get sync states for history
	syncStates, err := h.syncService.GetSyncStatesBySession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sync session not found"})
		return
	}

	// Build sync event history
	syncEvents := make([]SyncEvent, 0)
	var startTime, endTime time.Time

	for _, state := range syncStates {
		// Create event for this sync state
		event := SyncEvent{
			EventID:     uuid.New(),
			Platform:    state.PlatformType,
			Position:    state.SyncData.Position,
			EventType:   "position_update",
			Timestamp:   *state.LastSyncTime,
			SyncLatency: state.SyncMetrics.LastSyncLatency,
		}

		// Add conflict info if detected
		if state.ConflictState.HasConflicts {
			event.EventType = "conflict"
			event.ConflictInfo = &ConflictInfo{
				ConflictType:         string(state.ConflictState.ConflictType),
				ConflictingPlatforms: []string{string(state.PlatformType)},
				ResolutionStrategy:   string(state.ConflictState.ConflictResolution),
				ResolvedPosition:     state.SyncData.Position,
			}
		}

		syncEvents = append(syncEvents, event)

		// Track time range
		if startTime.IsZero() || state.LastSyncTime.Before(startTime) {
			startTime = *state.LastSyncTime
		}
		if state.LastSyncTime.After(endTime) {
			endTime = *state.LastSyncTime
		}
	}

	response := SyncHistoryResponse{
		VideoID:     videoID,
		SessionID:   sessionID,
		SyncEvents:  syncEvents,
		TotalEvents: len(syncEvents),
		StartTime:   startTime,
		EndTime:     endTime,
	}

	c.JSON(http.StatusOK, response)
}

// ForceSyncUpdate forces an immediate sync update across all platforms
// POST /api/v1/videos/:id/sync/force
func (h *SyncHandler) ForceSyncUpdate(c *gin.Context) {
	startTime := time.Now()

	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse request body
	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

// 	// Validate video exists
// 	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
// 		return
// 	}

	// Force sync across all platforms
	err = h.syncService.ForceSyncAll(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Force sync failed: %v", err),
		})
		return
	}

	// Calculate sync latency
	syncLatency := int(time.Since(startTime).Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"video_id":      videoID,
		"session_id":    req.SessionID,
		"position":      req.Position,
		"sync_latency_ms": syncLatency,
		"timestamp":     time.Now(),
		"message":       "Force sync completed successfully",
	})
}// Add this helper function at the end of sync_handler.go
func playbackState(isPlaying bool) string {
	if isPlaying {
		return "playing"
	}
	return "paused"
}
