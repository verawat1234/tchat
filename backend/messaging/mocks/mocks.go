package mocks

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
	sharedModels "tchat.dev/shared/models"
)

// MockEventPublisher implements services.EventPublisher for testing
type MockEventPublisher struct{}

func (m *MockEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	log.Printf("Event published: %s - %s", event.Type, event.Subject)
	return nil
}

// MockMessageRepository implements services.MessageRepository for testing
type MockMessageRepository struct {
	db *gorm.DB
}

func NewMockMessageRepository(db *gorm.DB) *MockMessageRepository {
	return &MockMessageRepository{db: db}
}

func (m *MockMessageRepository) Create(ctx context.Context, message *models.Message) error {
	return m.db.WithContext(ctx).Create(message).Error
}

func (m *MockMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := m.db.WithContext(ctx).First(&message, id).Error
	return &message, err
}

func (m *MockMessageRepository) GetByDialogID(ctx context.Context, dialogID uuid.UUID, filters services.MessageFilters, pagination services.Pagination) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64
	query := m.db.WithContext(ctx).Where("dialog_id = ?", dialogID)

	// Count total
	query.Model(&models.Message{}).Count(&total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Limit(pagination.PageSize).Offset(offset).Find(&messages).Error
	return messages, total, err
}

func (m *MockMessageRepository) Update(ctx context.Context, message *models.Message) error {
	return m.db.WithContext(ctx).Save(message).Error
}

func (m *MockMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}

func (m *MockMessageRepository) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	log.Printf("Marking message %s as read for user %s", messageID, userID)
	return nil
}

func (m *MockMessageRepository) GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *MockMessageRepository) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	log.Printf("Marking message %s as delivered for user %s", messageID, userID)
	return nil
}

func (m *MockMessageRepository) SearchMessages(ctx context.Context, dialogID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	err := m.db.WithContext(ctx).Where("dialog_id = ? AND content ILIKE ?", dialogID, "%"+query+"%").Limit(limit).Find(&messages).Error
	return messages, err
}

func (m *MockMessageRepository) GetMessageStats(ctx context.Context, dialogID uuid.UUID) (*services.MessageStats, error) {
	return &services.MessageStats{
		TotalMessages:  100,
		TextMessages:   80,
		MediaMessages:  15,
		SystemMessages: 5,
		AverageLength:  50.5,
		MessagesPerDay: 10,
		ActiveSenders:  5,
	}, nil
}

// MockDialogRepository implements services.DialogRepository for testing
type MockDialogRepository struct {
	db *gorm.DB
}

func NewMockDialogRepository(db *gorm.DB) *MockDialogRepository {
	return &MockDialogRepository{db: db}
}

func (m *MockDialogRepository) Create(ctx context.Context, dialog *models.Dialog) error {
	return m.db.WithContext(ctx).Create(dialog).Error
}

func (m *MockDialogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error) {
	var dialog models.Dialog
	err := m.db.WithContext(ctx).First(&dialog, id).Error
	return &dialog, err
}

func (m *MockDialogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filters services.DialogFilters, pagination services.Pagination) ([]*models.Dialog, int64, error) {
	var dialogs []*models.Dialog
	var total int64
	query := m.db.WithContext(ctx).Where("? = ANY(participants)", userID)

	// Count total
	query.Model(&models.Dialog{}).Count(&total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Limit(pagination.PageSize).Offset(offset).Find(&dialogs).Error
	return dialogs, total, err
}

func (m *MockDialogRepository) GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var participants []*models.DialogParticipant
	err := m.db.WithContext(ctx).Where("dialog_id = ?", dialogID).Find(&participants).Error
	return participants, err
}

func (m *MockDialogRepository) AddParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return m.db.WithContext(ctx).Create(participant).Error
}

func (m *MockDialogRepository) RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error {
	return m.db.WithContext(ctx).Where("dialog_id = ? AND user_id = ?", dialogID, userID).Delete(&models.DialogParticipant{}).Error
}

func (m *MockDialogRepository) UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return m.db.WithContext(ctx).Save(participant).Error
}

func (m *MockDialogRepository) GetAdmins(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var admins []*models.DialogParticipant
	err := m.db.WithContext(ctx).Where("dialog_id = ? AND role IN ?", dialogID, []string{"admin", "owner"}).Find(&admins).Error
	return admins, err
}

func (m *MockDialogRepository) SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error) {
	var dialogs []*models.Dialog
	err := m.db.WithContext(ctx).Where("? = ANY(participants) AND name ILIKE ?", userID, "%"+query+"%").Limit(limit).Find(&dialogs).Error
	return dialogs, err
}

func (m *MockDialogRepository) Update(ctx context.Context, dialog *models.Dialog) error {
	return m.db.WithContext(ctx).Save(dialog).Error
}

func (m *MockDialogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Dialog{}, id).Error
}

// MockPresenceRepository implements services.PresenceRepository for testing
type MockPresenceRepository struct {
	db *gorm.DB
}

func NewMockPresenceRepository(db *gorm.DB) *MockPresenceRepository {
	return &MockPresenceRepository{db: db}
}

func (m *MockPresenceRepository) Create(ctx context.Context, presence *models.Presence) error {
	return m.db.WithContext(ctx).Create(presence).Error
}

func (m *MockPresenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Presence, error) {
	var presence models.Presence
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).First(&presence).Error
	return &presence, err
}

func (m *MockPresenceRepository) Update(ctx context.Context, presence *models.Presence) error {
	return m.db.WithContext(ctx).Save(presence).Error
}

func (m *MockPresenceRepository) GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := m.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&presences).Error
	return presences, err
}

func (m *MockPresenceRepository) GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := m.db.WithContext(ctx).Where("is_online = ?", true).Limit(limit).Find(&presences).Error
	return presences, err
}

func (m *MockPresenceRepository) CleanupStalePresence(ctx context.Context, staleThreshold time.Duration) error {
	threshold := time.Now().Add(-staleThreshold)
	return m.db.WithContext(ctx).Model(&models.Presence{}).
		Where("last_updated < ? AND is_online = ?", threshold, true).
		Update("is_online", false).Error
}

func (m *MockPresenceRepository) GetPresenceStats(ctx context.Context) (*services.PresenceStats, error) {
	return &services.PresenceStats{
		TotalUsers:      1000,
		OnlineUsers:     250,
		AwayUsers:       50,
		BusyUsers:       25,
		OfflineUsers:    675,
		AverageUptime:   4 * time.Hour,
		PeakOnlineTime:  time.Now().Add(-2 * time.Hour),
		PeakOnlineCount: 300,
	}, nil
}

// MockWebSocketManager implements services.WebSocketManager for testing
type MockWebSocketManager struct{}

func (m *MockWebSocketManager) BroadcastToUser(ctx context.Context, userID uuid.UUID, message interface{}) error {
	log.Printf("Broadcasting to user %s: %+v", userID, message)
	return nil
}

func (m *MockWebSocketManager) BroadcastToUsers(ctx context.Context, userIDs []uuid.UUID, message interface{}) error {
	log.Printf("Broadcasting to %d users: %+v", len(userIDs), message)
	return nil
}

func (m *MockWebSocketManager) GetConnectedUsers(ctx context.Context) []uuid.UUID {
	// Return some mock connected users
	return []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
}

func (m *MockWebSocketManager) IsUserConnected(ctx context.Context, userID uuid.UUID) bool {
	// Mock implementation - would check actual connections
	return true
}

// MockLocationService implements services.LocationService for testing
type MockLocationService struct{}

func (m *MockLocationService) UpdateUserLocation(ctx context.Context, userID uuid.UUID, location models.Location) error {
	log.Printf("Updated location for user %s: %+v", userID, location)
	return nil
}

func (m *MockLocationService) GetNearbyUsers(ctx context.Context, userID uuid.UUID, radius float64) ([]uuid.UUID, error) {
	// Return some mock nearby users
	return []uuid.UUID{uuid.New(), uuid.New()}, nil
}

// MockNotificationService implements services.NotificationService for testing
type MockNotificationService struct{}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, notificationType string, data map[string]interface{}) error {
	log.Printf("Notification sent to user %s: %s - %+v", userID, notificationType, data)
	return nil
}

// MockDeliveryService implements services.MessageDeliveryService for testing
type MockDeliveryService struct{}

func (m *MockDeliveryService) DeliverMessage(ctx context.Context, message *models.Message, recipientIDs []uuid.UUID) error {
	log.Printf("Message %s delivered to %d recipients", message.ID, len(recipientIDs))
	return nil
}

func (m *MockDeliveryService) SendPushNotification(ctx context.Context, userID uuid.UUID, message *models.Message) error {
	log.Printf("Push notification sent to user %s for message %s", userID, message.ID)
	return nil
}

// MockContentModerator implements services.ContentModerator for testing
type MockContentModerator struct{}

func (m *MockContentModerator) ModerateContent(ctx context.Context, content string, contentType models.MessageType) (*services.ModerationResult, error) {
	return &services.ModerationResult{
		IsApproved:      true,
		Violations:      []string{},
		Confidence:      0.95,
		FilteredContent: content,
	}, nil
}

func (m *MockContentModerator) DetectSpam(ctx context.Context, senderID uuid.UUID, content string) (*services.SpamDetectionResult, error) {
	return &services.SpamDetectionResult{
		IsSpam:     false,
		Confidence: 0.1,
		Reasons:    []string{},
	}, nil
}

// MockMediaProcessor implements services.MediaProcessor for testing
type MockMediaProcessor struct{}

func (m *MockMediaProcessor) ProcessImageUpload(ctx context.Context, imageData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	return &services.ProcessedMedia{
		URL:          "https://example.com/image.jpg",
		ThumbnailURL: "https://example.com/thumb.jpg",
		Size:         int64(len(imageData)),
		Width:        800,
		Height:       600,
		Format:       "jpeg",
		Metadata:     metadata,
	}, nil
}

func (m *MockMediaProcessor) ProcessVideoUpload(ctx context.Context, videoData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	duration := 30.0
	return &services.ProcessedMedia{
		URL:          "https://example.com/video.mp4",
		ThumbnailURL: "https://example.com/video_thumb.jpg",
		Size:         int64(len(videoData)),
		Width:        1920,
		Height:       1080,
		Duration:     &duration,
		Format:       "mp4",
		Metadata:     metadata,
	}, nil
}

func (m *MockMediaProcessor) ProcessAudioUpload(ctx context.Context, audioData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	duration := 60.0
	return &services.ProcessedMedia{
		URL:      "https://example.com/audio.mp3",
		Size:     int64(len(audioData)),
		Duration: &duration,
		Format:   "mp3",
		Metadata: metadata,
	}, nil
}

func (m *MockMediaProcessor) GenerateThumbnail(ctx context.Context, mediaURL string, mediaType string) (string, error) {
	return "https://example.com/thumbnail.jpg", nil
}