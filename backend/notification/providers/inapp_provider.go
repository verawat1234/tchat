package providers

import (
	"context"
	"fmt"
	"log"

	"tchat.dev/notification/models"
)

// InAppProvider implements ChannelProvider for in-app notifications
type InAppProvider struct {
	websocketManager WebSocketManager
	databaseRepo     InAppRepository
}

// InAppConfig holds in-app provider configuration
type InAppConfig struct {
	EnableRealTime     bool `json:"enable_realtime"`
	EnablePersistence  bool `json:"enable_persistence"`
	MaxRetries         int  `json:"max_retries"`
	RetryDelaySeconds  int  `json:"retry_delay_seconds"`
}

// WebSocketManager interface for real-time delivery
type WebSocketManager interface {
	SendToUser(ctx context.Context, userID string, message interface{}) error
	SendToUsers(ctx context.Context, userIDs []string, message interface{}) error
	IsUserConnected(userID string) bool
	GetConnectedUsers() []string
}

// InAppRepository interface for persistent storage
type InAppRepository interface {
	StoreNotification(ctx context.Context, notification *models.Notification) error
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, notificationID, userID string) error
	DeleteNotification(ctx context.Context, notificationID, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
}

// NewInAppProvider creates a new in-app notification provider
func NewInAppProvider(wsManager WebSocketManager, repo InAppRepository) *InAppProvider {
	return &InAppProvider{
		websocketManager: wsManager,
		databaseRepo:     repo,
	}
}

// Send sends an in-app notification
func (i *InAppProvider) Send(ctx context.Context, notification *models.Notification) error {
	if notification.UserID.String() == "" {
		return fmt.Errorf("user ID is required for in-app notifications")
	}

	userID := notification.UserID.String()

	// Store notification in database for persistence
	if i.databaseRepo != nil {
		if err := i.databaseRepo.StoreNotification(ctx, notification); err != nil {
			log.Printf("Failed to store in-app notification: %v", err)
			// Continue to try real-time delivery even if storage fails
		}
	}

	// Send real-time notification if user is connected
	if i.websocketManager != nil && i.websocketManager.IsUserConnected(userID) {
		message := i.buildInAppMessage(notification)
		if err := i.websocketManager.SendToUser(ctx, userID, message); err != nil {
			log.Printf("Failed to send real-time notification to user %s: %v", userID, err)
			// Notification is still stored in database, so this is not a critical failure
		} else {
			log.Printf("Sent real-time in-app notification to user %s", userID)
		}
	} else {
		log.Printf("User %s not connected, notification stored for later retrieval", userID)
	}

	return nil
}

// SendBatch sends multiple in-app notifications
func (i *InAppProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	// Group notifications by user for efficient delivery
	userNotifications := make(map[string][]*models.Notification)

	for _, notification := range notifications {
		userID := notification.UserID.String()
		userNotifications[userID] = append(userNotifications[userID], notification)
	}

	// Process each user's notifications
	for userID, userNotifs := range userNotifications {
		// Store all notifications for this user
		if i.databaseRepo != nil {
			for _, notification := range userNotifs {
				if err := i.databaseRepo.StoreNotification(ctx, notification); err != nil {
					log.Printf("Failed to store notification for user %s: %v", userID, err)
				}
			}
		}

		// Send real-time notifications if user is connected
		if i.websocketManager != nil && i.websocketManager.IsUserConnected(userID) {
			messages := make([]interface{}, len(userNotifs))
			for j, notification := range userNotifs {
				messages[j] = i.buildInAppMessage(notification)
			}

			// Send batch message
			batchMessage := map[string]interface{}{
				"type":          "notification_batch",
				"notifications": messages,
				"count":         len(messages),
			}

			if err := i.websocketManager.SendToUser(ctx, userID, batchMessage); err != nil {
				log.Printf("Failed to send batch notification to user %s: %v", userID, err)
			}
		}
	}

	return nil
}

// ValidateConfig validates the in-app provider configuration
func (i *InAppProvider) ValidateConfig() error {
	if i.websocketManager == nil && i.databaseRepo == nil {
		return fmt.Errorf("either WebSocket manager or database repository is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (i *InAppProvider) GetProviderName() string {
	return "in-app"
}

// SupportsChannels returns the channels this provider supports
func (i *InAppProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelInApp}
}

// buildInAppMessage builds the in-app notification message
func (i *InAppProvider) buildInAppMessage(notification *models.Notification) map[string]interface{} {
	message := map[string]interface{}{
		"type":            "notification",
		"id":              notification.ID.String(),
		"title":           notification.Title,
		"body":            notification.Body,
		"priority":        string(notification.Priority),
		"category":        string(notification.Category),
		"created_at":      notification.CreatedAt,
		"read":            notification.Read,
	}

	// Add optional fields
	if notification.MediaURL != "" {
		message["media_url"] = notification.MediaURL
	}

	if notification.ActionURL != "" {
		message["action_url"] = notification.ActionURL
	}

	if notification.ActionText != "" {
		message["action_text"] = notification.ActionText
	}

	if notification.IconURL != "" {
		message["icon_url"] = notification.IconURL
	}

	// Add metadata if present
	if notification.Metadata != nil {
		message["metadata"] = notification.Metadata
	}

	return message
}

// SimpleWebSocketManager is a basic implementation of WebSocketManager
type SimpleWebSocketManager struct {
	connections map[string]WebSocketConnection
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection interface {
	Send(message interface{}) error
	Close() error
	UserID() string
}

// NewSimpleWebSocketManager creates a new simple WebSocket manager
func NewSimpleWebSocketManager() *SimpleWebSocketManager {
	return &SimpleWebSocketManager{
		connections: make(map[string]WebSocketConnection),
	}
}

// SendToUser sends a message to a specific user
func (s *SimpleWebSocketManager) SendToUser(ctx context.Context, userID string, message interface{}) error {
	conn, exists := s.connections[userID]
	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	return conn.Send(message)
}

// SendToUsers sends a message to multiple users
func (s *SimpleWebSocketManager) SendToUsers(ctx context.Context, userIDs []string, message interface{}) error {
	var lastErr error
	successCount := 0

	for _, userID := range userIDs {
		if err := s.SendToUser(ctx, userID, message); err != nil {
			lastErr = err
			log.Printf("Failed to send message to user %s: %v", userID, err)
		} else {
			successCount++
		}
	}

	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("failed to send message to any users: %w", lastErr)
	}

	return nil
}

// IsUserConnected checks if a user is connected
func (s *SimpleWebSocketManager) IsUserConnected(userID string) bool {
	_, exists := s.connections[userID]
	return exists
}

// GetConnectedUsers returns list of connected user IDs
func (s *SimpleWebSocketManager) GetConnectedUsers() []string {
	users := make([]string, 0, len(s.connections))
	for userID := range s.connections {
		users = append(users, userID)
	}
	return users
}

// AddConnection adds a WebSocket connection for a user
func (s *SimpleWebSocketManager) AddConnection(userID string, conn WebSocketConnection) {
	s.connections[userID] = conn
}

// RemoveConnection removes a WebSocket connection for a user
func (s *SimpleWebSocketManager) RemoveConnection(userID string) {
	if conn, exists := s.connections[userID]; exists {
		conn.Close()
		delete(s.connections, userID)
	}
}

// BroadcastToAll sends a message to all connected users
func (s *SimpleWebSocketManager) BroadcastToAll(ctx context.Context, message interface{}) error {
	connectedUsers := s.GetConnectedUsers()
	return s.SendToUsers(ctx, connectedUsers, message)
}