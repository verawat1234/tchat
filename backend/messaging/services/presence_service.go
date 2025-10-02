package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	sharedModels "tchat.dev/shared/models"
)

type PresenceRepository interface {
	Create(ctx context.Context, presence *models.Presence) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Presence, error)
	Update(ctx context.Context, presence *models.Presence) error
	GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Presence, error)
	GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error)
	CleanupStalePresence(ctx context.Context, staleThreshold time.Duration) error
	GetPresenceStats(ctx context.Context) (*PresenceStats, error)
}

type WebSocketManager interface {
	RegisterClient(userID uuid.UUID, conn interface{})
	BroadcastToUser(ctx context.Context, userID uuid.UUID, message interface{}) error
	BroadcastToUsers(ctx context.Context, userIDs []uuid.UUID, message interface{}) error
	GetConnectedUsers(ctx context.Context) []uuid.UUID
	IsUserConnected(ctx context.Context, userID uuid.UUID) bool
}

type LocationService interface {
	UpdateUserLocation(ctx context.Context, userID uuid.UUID, location models.Location) error
	GetNearbyUsers(ctx context.Context, userID uuid.UUID, radius float64) ([]uuid.UUID, error)
}

type PresenceStats struct {
	TotalUsers     int64 `json:"total_users"`
	OnlineUsers    int64 `json:"online_users"`
	AwayUsers      int64 `json:"away_users"`
	BusyUsers      int64 `json:"busy_users"`
	OfflineUsers   int64 `json:"offline_users"`
	AverageUptime  time.Duration `json:"average_uptime"`
	PeakOnlineTime time.Time `json:"peak_online_time"`
	PeakOnlineCount int64 `json:"peak_online_count"`
}

type PresenceService struct {
	presenceRepo      PresenceRepository
	wsManager         WebSocketManager
	locationService   LocationService
	eventPublisher    EventPublisher
	db                *gorm.DB
}

func NewPresenceService(
	presenceRepo PresenceRepository,
	wsManager WebSocketManager,
	locationService LocationService,
	eventPublisher EventPublisher,
	db *gorm.DB,
) *PresenceService {
	return &PresenceService{
		presenceRepo:    presenceRepo,
		wsManager:       wsManager,
		locationService: locationService,
		eventPublisher:  eventPublisher,
		db:              db,
	}
}

func (ps *PresenceService) UpdatePresence(ctx context.Context, req *UpdatePresenceRequest) (*models.Presence, error) {
	// Validate request
	if err := ps.validateUpdatePresenceRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get or create presence record
	presence, err := ps.presenceRepo.GetByUserID(ctx, req.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	isNewPresence := false
	if err == gorm.ErrRecordNotFound {
		// Create new presence record
		presence = &models.Presence{
			ID:        uuid.New(),
			UserID:    req.UserID,
			Status:    models.PresenceStatusOffline,
			IsOnline:  false,
			Platform:  models.PlatformWeb, // Default platform
			Activity:  models.ActivityStatusIdle,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		isNewPresence = true
	}

	// Track previous status for change detection
	previousStatus := presence.Status
	previousActivity := presence.Activity

	// Update presence fields
	if req.Status != nil {
		// TODO: Implement proper status transition validation
		presence.Status = *req.Status
	}

	if req.Activity != nil {
		// TODO: Convert string to ActivityStatus properly
		presence.Activity = models.ActivityStatusIdle
	}

	if req.Platform != nil {
		presence.Platform = *req.Platform
	}

	if req.DeviceInfo != nil {
		presence.DeviceInfo = *req.DeviceInfo
	}

	if req.Location != nil {
		presence.Location = req.Location
		// TODO: Implement location sharing when Settings are available
	}

	// TODO: Implement Settings when PresenceSettings model is available

	// Update timestamps based on status
	now := time.Now()
	presence.UpdatedAt = now

	switch presence.Status {
	case models.PresenceStatusOnline:
		presence.IsOnline = true
		presence.LastSeen = &now
		// TODO: Track online since when OnlineSince field is available
	case models.PresenceStatusAway, models.PresenceStatusBusy:
		presence.IsOnline = true
		presence.LastSeen = &now
	case models.PresenceStatusOffline:
		presence.IsOnline = false
		if previousStatus != models.PresenceStatusOffline {
			// TODO: Calculate session duration when OnlineSince and TotalOnlineTime fields are available
		}
	}

	// TODO: Apply business hours logic when Settings are available

	// Save presence
	if isNewPresence {
		if err := ps.presenceRepo.Create(ctx, presence); err != nil {
			return nil, fmt.Errorf("failed to create presence: %w", err)
		}
	} else {
		if err := ps.presenceRepo.Update(ctx, presence); err != nil {
			return nil, fmt.Errorf("failed to update presence: %w", err)
		}
	}

	// Broadcast presence change to relevant users
	if previousStatus != presence.Status || previousActivity != presence.Activity {
		go ps.broadcastPresenceChange(context.Background(), presence, previousStatus, string(previousActivity))
	}

	// Publish presence change event
	if previousStatus != presence.Status {
		if err := ps.publishPresenceEvent(ctx, sharedModels.EventTypeUserPresenceChanged, req.UserID, map[string]interface{}{
			"previous_status": previousStatus,
			"new_status":     presence.Status,
			"activity":       presence.Activity,
			"platform":       presence.Platform,
			"is_online":      presence.IsOnline,
		}); err != nil {
			fmt.Printf("Failed to publish presence change event: %v\n", err)
		}
	}

	return presence, nil
}

func (ps *PresenceService) GetPresence(ctx context.Context, userID uuid.UUID, requestorID uuid.UUID) (*models.Presence, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	presence, err := ps.presenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("presence not found")
		}
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	// Apply privacy settings
	ps.applyPrivacySettings(presence, requestorID)

	return presence, nil
}

func (ps *PresenceService) GetMultiplePresence(ctx context.Context, userIDs []uuid.UUID, requestorID uuid.UUID) ([]*models.Presence, error) {
	if len(userIDs) == 0 {
		return []*models.Presence{}, nil
	}

	if len(userIDs) > 100 {
		return nil, fmt.Errorf("cannot fetch presence for more than 100 users at once")
	}

	presences, err := ps.presenceRepo.GetByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple presence: %w", err)
	}

	// Apply privacy settings for each presence
	for _, presence := range presences {
		ps.applyPrivacySettings(presence, requestorID)
	}

	return presences, nil
}

func (ps *PresenceService) SetOffline(ctx context.Context, userID uuid.UUID) error {
	req := &UpdatePresenceRequest{
		UserID: userID,
		Status: &[]models.PresenceStatus{models.PresenceStatusOffline}[0],
	}

	_, err := ps.UpdatePresence(ctx, req)
	return err
}

func (ps *PresenceService) SetOnline(ctx context.Context, userID uuid.UUID, platform models.Platform, deviceInfo models.DeviceInfo) error {
	req := &UpdatePresenceRequest{
		UserID:     userID,
		Status:     &[]models.PresenceStatus{models.PresenceStatusOnline}[0],
		Platform:   &platform,
		DeviceInfo: &deviceInfo,
	}

	_, err := ps.UpdatePresence(ctx, req)
	return err
}

func (ps *PresenceService) UpdateActivity(ctx context.Context, userID uuid.UUID, activity string) error {
	req := &UpdatePresenceRequest{
		UserID:   userID,
		Activity: &activity,
	}

	_, err := ps.UpdatePresence(ctx, req)
	return err
}

func (ps *PresenceService) UpdateSettings(ctx context.Context, userID uuid.UUID, settings models.PresenceSettings) error {
	req := &UpdatePresenceRequest{
		UserID:   userID,
		Settings: &settings,
	}

	_, err := ps.UpdatePresence(ctx, req)
	return err
}

func (ps *PresenceService) GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	presences, err := ps.presenceRepo.GetOnlineUsers(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get online users: %w", err)
	}

	return presences, nil
}

func (ps *PresenceService) GetPresenceStats(ctx context.Context) (*PresenceStats, error) {
	stats, err := ps.presenceRepo.GetPresenceStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence stats: %w", err)
	}

	return stats, nil
}

func (ps *PresenceService) CleanupStalePresence(ctx context.Context) error {
	staleThreshold := 30 * time.Minute // Consider presence stale after 30 minutes of inactivity
	return ps.presenceRepo.CleanupStalePresence(ctx, staleThreshold)
}

func (ps *PresenceService) AutoUpdateFromWebSocket(ctx context.Context, userID uuid.UUID) {
	// Called when user connects/disconnects from WebSocket
	isConnected := ps.wsManager.IsUserConnected(ctx, userID)

	if isConnected {
		// User connected - set online
		ps.SetOnline(ctx, userID, models.PlatformWeb, models.DeviceInfo{
			OS:          "web",
			DeviceName:  "unknown",
			AppVersion:  "1.0.0",
		})
	} else {
		// User disconnected - set offline
		ps.SetOffline(ctx, userID)
	}
}

func (ps *PresenceService) ScheduleAutoAway(ctx context.Context, userID uuid.UUID, timeout time.Duration) {
	// This would typically be handled by a background job scheduler
	// For now, we'll just update the auto-away timeout in settings
	presence, err := ps.presenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return
	}

	// TODO: Set AutoAwayTimeout when Settings field is available
	ps.presenceRepo.Update(ctx, presence)
}

// Private helper methods

func (ps *PresenceService) applyBusinessHoursLogic(presence *models.Presence) {
	// now := time.Now() // TODO: Uncomment when business hours logic is implemented

	// Check if current time is within business hours
	// TODO: Implement business hours logic when Settings field is available
	// isBusinessDay := ps.isBusinessDay(now, businessDays)
	// isBusinessHour := ps.isBusinessHour(now, businessHoursStart, businessHoursEnd)

	if false { // TODO: Replace with proper business hours check
		// During business hours, auto-set to online if not already
		if presence.Status == models.PresenceStatusOffline {
			presence.Status = models.PresenceStatusOnline
			presence.IsOnline = true
		}
	}
}

func (ps *PresenceService) isBusinessDay(t time.Time, businessDays []string) bool {
	dayNames := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}
	currentDay := dayNames[t.Weekday()]

	for _, day := range businessDays {
		if day == currentDay {
			return true
		}
	}
	return false
}

func (ps *PresenceService) isBusinessHour(t time.Time, startTime, endTime string) bool {
	currentTime := t.Format("15:04")
	return currentTime >= startTime && currentTime <= endTime
}

func (ps *PresenceService) applyPrivacySettings(presence *models.Presence, requestorID uuid.UUID) {
	// TODO: Apply privacy settings when Settings field is available
	// For now, show all information to all users
	_ = requestorID // Suppress unused parameter warning
}

func (ps *PresenceService) broadcastPresenceChange(ctx context.Context, presence *models.Presence, previousStatus models.PresenceStatus, previousActivity string) {
	// Create presence update message
	update := PresenceUpdateMessage{
		UserID:           presence.UserID,
		Status:           presence.Status,
		Activity:         string(presence.Activity),
		IsOnline:         presence.IsOnline,
		LastSeen:         presence.LastSeen,
		Platform:         presence.Platform,
		PreviousStatus:   previousStatus,
		PreviousActivity: previousActivity,
		UpdatedAt:        presence.UpdatedAt,
	}

	// Get connected users to broadcast to
	// This could be optimized to only broadcast to users who have this user in their contact list
	connectedUsers := ps.wsManager.GetConnectedUsers(ctx)

	// Broadcast to connected users (excluding the user whose presence changed)
	recipientUsers := make([]uuid.UUID, 0)
	for _, userID := range connectedUsers {
		if userID != presence.UserID {
			recipientUsers = append(recipientUsers, userID)
		}
	}

	if len(recipientUsers) > 0 {
		ps.wsManager.BroadcastToUsers(ctx, recipientUsers, update)
	}
}

func (ps *PresenceService) validateUpdatePresenceRequest(req *UpdatePresenceRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if req.Status != nil {
		validStatuses := []models.PresenceStatus{
			models.PresenceStatusOnline,
			models.PresenceStatusAway,
			models.PresenceStatusBusy,
			models.PresenceStatusOffline,
		}

		isValid := false
		for _, status := range validStatuses {
			if *req.Status == status {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("invalid presence status: %s", *req.Status)
		}
	}

	return nil
}

func (ps *PresenceService) publishPresenceEvent(ctx context.Context, eventType sharedModels.EventType, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("Presence event: %s", eventType),
		AggregateID:   userID.String(),
		AggregateType: "user",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
			Source:      "messaging-service",
			Environment: "production",
			Region:      "sea",
		},
	}

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return ps.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type UpdatePresenceRequest struct {
	UserID     uuid.UUID                `json:"user_id" binding:"required"`
	Status     *models.PresenceStatus   `json:"status,omitempty"`
	Activity   *string                  `json:"activity,omitempty"`
	Platform   *models.Platform         `json:"platform,omitempty"`
	DeviceInfo *models.DeviceInfo       `json:"device_info,omitempty"`
	Location   *models.Location         `json:"location,omitempty"`
	Settings   *models.PresenceSettings `json:"settings,omitempty"`
}

type PresenceResponse struct {
	ID                uuid.UUID                `json:"id"`
	UserID            uuid.UUID                `json:"user_id"`
	Status            models.PresenceStatus    `json:"status"`
	Activity          string                   `json:"activity,omitempty"`
	IsOnline          bool                     `json:"is_online"`
	LastSeen          *time.Time               `json:"last_seen,omitempty"`
	Platform          models.Platform          `json:"platform"`
	DeviceInfo        models.DeviceInfo        `json:"device_info"`
	Location          *models.Location         `json:"location,omitempty"`
	OnlineSince       *time.Time               `json:"online_since,omitempty"`
	TotalOnlineTime   time.Duration            `json:"total_online_time"`
	Settings          models.PresenceSettings  `json:"settings"`
	LastUpdated       time.Time                `json:"last_updated"`
}

type PresenceUpdateMessage struct {
	UserID           uuid.UUID             `json:"user_id"`
	Status           models.PresenceStatus `json:"status"`
	Activity         string                `json:"activity,omitempty"`
	IsOnline         bool                  `json:"is_online"`
	LastSeen         *time.Time            `json:"last_seen,omitempty"`
	Platform         models.Platform       `json:"platform"`
	PreviousStatus   models.PresenceStatus `json:"previous_status"`
	PreviousActivity string                `json:"previous_activity"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type OnlineUsersResponse struct {
	Users     []*PresenceResponse `json:"users"`
	Total     int                 `json:"total"`
	Online    int                 `json:"online"`
	Away      int                 `json:"away"`
	Busy      int                 `json:"busy"`
}

func PresenceToResponse(presence *models.Presence) *PresenceResponse {
	return &PresenceResponse{
		ID:              presence.ID,
		UserID:          presence.UserID,
		Status:          presence.Status,
		Activity:        string(presence.Activity),
		IsOnline:        presence.IsOnline,
		LastSeen:        presence.LastSeen,
		Platform:        presence.Platform,
		DeviceInfo:      presence.DeviceInfo,
		Location:        presence.Location,
		// OnlineSince:     presence.OnlineSince, // TODO: Add when OnlineSince field is available
		// TotalOnlineTime: presence.TotalOnlineTime, // TODO: Add when TotalOnlineTime field is available
		// Settings:        presence.Settings, // TODO: Add when Settings field is available
		// LastUpdated:     presence.LastUpdated, // TODO: Add when LastUpdated field is available
	}
}