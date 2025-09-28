package external

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// MessageDeliveryService implements services.MessageDeliveryService
type MessageDeliveryService struct {
	notificationService services.NotificationService
	wsManager          services.WebSocketManager
}

// NewMessageDeliveryService creates a new message delivery service
func NewMessageDeliveryService(notificationService services.NotificationService, wsManager services.WebSocketManager) services.MessageDeliveryService {
	return &MessageDeliveryService{
		notificationService: notificationService,
		wsManager:          wsManager,
	}
}

// DeliverMessage delivers a message to recipients via WebSocket and push notifications
func (d *MessageDeliveryService) DeliverMessage(ctx context.Context, message *models.Message, recipientIDs []uuid.UUID) error {
	// Create message delivery payload
	messagePayload := map[string]interface{}{
		"type":      "new_message",
		"message":   message,
		"dialog_id": message.DialogID,
		"sender_id": message.SenderID,
	}

	// Send via WebSocket to online users
	err := d.wsManager.BroadcastToUsers(ctx, recipientIDs, messagePayload)
	if err != nil {
		log.Printf("Failed to broadcast message via WebSocket: %v", err)
	}

	// Send push notifications to offline users
	for _, recipientID := range recipientIDs {
		if !d.wsManager.IsUserConnected(ctx, recipientID) {
			notificationData := map[string]interface{}{
				"message_id":   message.ID.String(),
				"dialog_id":    message.DialogID.String(),
				"sender_id":    message.SenderID.String(),
				"content":      message.Content,
				"message_type": string(message.Type),
			}

			err := d.notificationService.SendNotification(ctx, recipientID, "message", notificationData)
			if err != nil {
				log.Printf("Failed to send push notification to user %s: %v", recipientID, err)
			}
		}
	}

	log.Printf("Message %s delivered to %d recipients", message.ID, len(recipientIDs))
	return nil
}

// SendPushNotification sends a push notification for a specific message
func (d *MessageDeliveryService) SendPushNotification(ctx context.Context, userID uuid.UUID, message *models.Message) error {
	notificationData := map[string]interface{}{
		"message_id":   message.ID.String(),
		"dialog_id":    message.DialogID.String(),
		"sender_id":    message.SenderID.String(),
		"content":      message.Content,
		"message_type": string(message.Type),
	}

	return d.notificationService.SendNotification(ctx, userID, "message", notificationData)
}

// LocationService implements services.LocationService for nearby user discovery
type LocationService struct {
	enabled bool
}

// NewLocationService creates a new location service
func NewLocationService() services.LocationService {
	return &LocationService{
		enabled: true,
	}
}

// UpdateUserLocation updates a user's location
func (l *LocationService) UpdateUserLocation(ctx context.Context, userID uuid.UUID, location models.Location) error {
	if !l.enabled {
		return nil
	}

	locationData, _ := json.Marshal(location)
	log.Printf("📍 Updated location for user %s: %s", userID, string(locationData))

	// TODO: Implement actual location storage
	// - Store in Redis with geospatial indexing
	// - Use PostGIS for PostgreSQL
	// - Integrate with mapping services

	return nil
}

// GetNearbyUsers returns users within a radius of the specified user
func (l *LocationService) GetNearbyUsers(ctx context.Context, userID uuid.UUID, radius float64) ([]uuid.UUID, error) {
	if !l.enabled {
		return []uuid.UUID{}, nil
	}

	log.Printf("🔍 Finding users within %.2f km of user %s", radius, userID)

	// TODO: Implement actual geospatial queries
	// For now, return some mock nearby users
	nearbyUsers := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	log.Printf("Found %d nearby users for user %s", len(nearbyUsers), userID)
	return nearbyUsers, nil
}

// ContentModerator implements services.ContentModerator for content moderation
type ContentModerator struct {
	enabled bool
}

// NewContentModerator creates a new content moderator
func NewContentModerator() services.ContentModerator {
	return &ContentModerator{
		enabled: true,
	}
}

// ModerateContent moderates message content
func (c *ContentModerator) ModerateContent(ctx context.Context, content string, contentType models.MessageType) (*services.ModerationResult, error) {
	if !c.enabled {
		return &services.ModerationResult{
			IsApproved:      true,
			Violations:      []string{},
			Confidence:      1.0,
			FilteredContent: content,
		}, nil
	}

	log.Printf("🔍 Moderating %s content: %s", contentType, content)

	// TODO: Implement actual content moderation
	// - Integrate with AI moderation services (OpenAI Moderation, Google Cloud AI)
	// - Implement keyword filtering
	// - Add profanity detection
	// - Image/video content analysis

	// Simple demonstration of moderation
	violations := []string{}
	isApproved := true
	filteredContent := content

	// Basic profanity check (demo)
	profanityWords := []string{"spam", "abuse", "hate"}
	for _, word := range profanityWords {
		if contains(content, word) {
			violations = append(violations, "inappropriate_language")
			isApproved = false
			filteredContent = "[Content filtered]"
			break
		}
	}

	result := &services.ModerationResult{
		IsApproved:      isApproved,
		Violations:      violations,
		Confidence:      0.95,
		FilteredContent: filteredContent,
	}

	log.Printf("✅ Moderation result: approved=%t, violations=%v", isApproved, violations)
	return result, nil
}

// DetectSpam detects spam in message content
func (c *ContentModerator) DetectSpam(ctx context.Context, senderID uuid.UUID, content string) (*services.SpamDetectionResult, error) {
	if !c.enabled {
		return &services.SpamDetectionResult{
			IsSpam:     false,
			Confidence: 0.0,
			Reasons:    []string{},
		}, nil
	}

	log.Printf("🔍 Checking for spam from user %s: %s", senderID, content)

	// TODO: Implement actual spam detection
	// - Rate limiting per user
	// - Content similarity detection
	// - Machine learning models
	// - Pattern recognition

	// Simple demonstration
	reasons := []string{}
	isSpam := false
	confidence := 0.1

	// Basic spam indicators
	if len(content) > 1000 {
		reasons = append(reasons, "message_too_long")
		isSpam = true
		confidence = 0.8
	}

	result := &services.SpamDetectionResult{
		IsSpam:     isSpam,
		Confidence: confidence,
		Reasons:    reasons,
	}

	log.Printf("🚫 Spam detection result: spam=%t, confidence=%.2f", isSpam, confidence)
	return result, nil
}

// MediaProcessor implements services.MediaProcessor for media upload processing
type MediaProcessor struct {
	enabled bool
}

// NewMediaProcessor creates a new media processor
func NewMediaProcessor() services.MediaProcessor {
	return &MediaProcessor{
		enabled: true,
	}
}

// ProcessImageUpload processes an image upload
func (m *MediaProcessor) ProcessImageUpload(ctx context.Context, imageData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	if !m.enabled {
		return nil, nil
	}

	log.Printf("📷 Processing image upload: %d bytes", len(imageData))

	// TODO: Implement actual image processing
	// - Upload to cloud storage (AWS S3, Google Cloud Storage)
	// - Image optimization and resizing
	// - Thumbnail generation
	// - Format conversion
	// - Metadata extraction

	result := &services.ProcessedMedia{
		URL:          "https://cdn.tchat.app/images/" + uuid.New().String() + ".jpg",
		ThumbnailURL: "https://cdn.tchat.app/thumbs/" + uuid.New().String() + "_thumb.jpg",
		Size:         int64(len(imageData)),
		Width:        1920,
		Height:       1080,
		Format:       "jpeg",
		Metadata:     metadata,
	}

	log.Printf("✅ Image processed: %s", result.URL)
	return result, nil
}

// ProcessVideoUpload processes a video upload
func (m *MediaProcessor) ProcessVideoUpload(ctx context.Context, videoData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	if !m.enabled {
		return nil, nil
	}

	log.Printf("🎥 Processing video upload: %d bytes", len(videoData))

	// TODO: Implement actual video processing
	// - Upload to cloud storage
	// - Video transcoding (multiple resolutions)
	// - Thumbnail generation
	// - Duration extraction
	// - Format optimization

	duration := 120.0
	result := &services.ProcessedMedia{
		URL:          "https://cdn.tchat.app/videos/" + uuid.New().String() + ".mp4",
		ThumbnailURL: "https://cdn.tchat.app/thumbs/" + uuid.New().String() + "_video_thumb.jpg",
		Size:         int64(len(videoData)),
		Width:        1920,
		Height:       1080,
		Duration:     &duration,
		Format:       "mp4",
		Metadata:     metadata,
	}

	log.Printf("✅ Video processed: %s", result.URL)
	return result, nil
}

// ProcessAudioUpload processes an audio upload
func (m *MediaProcessor) ProcessAudioUpload(ctx context.Context, audioData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	if !m.enabled {
		return nil, nil
	}

	log.Printf("🎵 Processing audio upload: %d bytes", len(audioData))

	// TODO: Implement actual audio processing
	// - Upload to cloud storage
	// - Audio format conversion
	// - Duration extraction
	// - Quality optimization

	duration := 180.0
	result := &services.ProcessedMedia{
		URL:      "https://cdn.tchat.app/audio/" + uuid.New().String() + ".mp3",
		Size:     int64(len(audioData)),
		Duration: &duration,
		Format:   "mp3",
		Metadata: metadata,
	}

	log.Printf("✅ Audio processed: %s", result.URL)
	return result, nil
}

// GenerateThumbnail generates a thumbnail for media
func (m *MediaProcessor) GenerateThumbnail(ctx context.Context, mediaURL string, mediaType string) (string, error) {
	if !m.enabled {
		return "", nil
	}

	log.Printf("🖼️  Generating thumbnail for %s media: %s", mediaType, mediaURL)

	// TODO: Implement actual thumbnail generation
	// - Extract frames from video
	// - Resize images
	// - Generate preview images

	thumbnailURL := "https://cdn.tchat.app/thumbs/" + uuid.New().String() + "_thumb.jpg"
	log.Printf("✅ Thumbnail generated: %s", thumbnailURL)
	return thumbnailURL, nil
}

// Helper function
func contains(text, word string) bool {
	return len(text) >= len(word) && text[:len(word)] == word
}