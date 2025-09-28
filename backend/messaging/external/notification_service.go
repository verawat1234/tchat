package external

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"tchat.dev/messaging/services"
)

// NotificationService implements services.NotificationService for push notifications
type NotificationService struct {
	enabled bool
	// In a real implementation, this would integrate with:
	// - Firebase Cloud Messaging (FCM)
	// - Apple Push Notification Service (APNs)
	// - AWS SNS
	// - OneSignal, etc.
}

// NewNotificationService creates a new notification service
func NewNotificationService() services.NotificationService {
	return &NotificationService{
		enabled: true,
	}
}

// SendNotification sends a push notification to a user
func (n *NotificationService) SendNotification(ctx context.Context, userID uuid.UUID, notificationType string, data map[string]interface{}) error {
	if !n.enabled {
		return nil
	}

	// Convert data to JSON for logging
	notificationData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal notification data: %v", err)
		return err
	}

	// Log the notification for demonstration
	log.Printf("Sending %s notification to user %s:\n%s", notificationType, userID, string(notificationData))

	// TODO: Implement actual push notification integration
	// Examples:
	// - Firebase Cloud Messaging (FCM) for Android/iOS
	// - Apple Push Notification Service (APNs) for iOS
	// - Web Push for web browsers
	// - SMS/Email notifications
	// - In-app notifications

	// Simulate different notification types
	switch notificationType {
	case "message":
		log.Printf("📱 Push notification: New message for user %s", userID)
	case "mention":
		log.Printf("📱 Push notification: You were mentioned, user %s", userID)
	case "dialog_invite":
		log.Printf("📱 Push notification: Dialog invitation for user %s", userID)
	case "presence_update":
		log.Printf("📱 Push notification: Presence update for user %s", userID)
	default:
		log.Printf("📱 Push notification: %s for user %s", notificationType, userID)
	}

	return nil
}