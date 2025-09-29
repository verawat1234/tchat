package providers

import (
	"context"
	"fmt"
	"log"

	"tchat.dev/notification/models"
)

// PushProvider implements ChannelProvider for push notifications
type PushProvider struct {
	fcmServerKey string
	apnsKeyID    string
	apnsTeamID   string
	apnsKeyPath  string
	bundleID     string
}

// PushConfig holds push notification provider configuration
type PushConfig struct {
	FCMServerKey string `json:"fcm_server_key"`
	APNSKeyID    string `json:"apns_key_id"`
	APNSTeamID   string `json:"apns_team_id"`
	APNSKeyPath  string `json:"apns_key_path"`
	BundleID     string `json:"bundle_id"`
}

// NewPushProvider creates a new push notification provider
func NewPushProvider(config PushConfig) *PushProvider {
	return &PushProvider{
		fcmServerKey: config.FCMServerKey,
		apnsKeyID:    config.APNSKeyID,
		apnsTeamID:   config.APNSTeamID,
		apnsKeyPath:  config.APNSKeyPath,
		bundleID:     config.BundleID,
	}
}

// Send sends a push notification
func (p *PushProvider) Send(ctx context.Context, notification *models.Notification) error {
	if notification.DeviceToken == "" {
		return fmt.Errorf("device token is required for push notifications")
	}

	// Determine platform and send accordingly
	switch notification.Platform {
	case "ios":
		return p.sendAPNS(ctx, notification)
	case "android":
		return p.sendFCM(ctx, notification)
	default:
		// Try both platforms if platform is not specified
		if err := p.sendFCM(ctx, notification); err != nil {
			return p.sendAPNS(ctx, notification)
		}
		return nil
	}
}

// SendBatch sends multiple push notifications
func (p *PushProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	// Group notifications by platform for efficient batch sending
	fcmNotifications := make([]*models.Notification, 0)
	apnsNotifications := make([]*models.Notification, 0)

	for _, notification := range notifications {
		switch notification.Platform {
		case "android":
			fcmNotifications = append(fcmNotifications, notification)
		case "ios":
			apnsNotifications = append(apnsNotifications, notification)
		default:
			// Default to FCM for unknown platforms
			fcmNotifications = append(fcmNotifications, notification)
		}
	}

	// Send FCM batch
	if len(fcmNotifications) > 0 {
		if err := p.sendFCMBatch(ctx, fcmNotifications); err != nil {
			log.Printf("Failed to send FCM batch: %v", err)
		}
	}

	// Send APNS batch
	if len(apnsNotifications) > 0 {
		if err := p.sendAPNSBatch(ctx, apnsNotifications); err != nil {
			log.Printf("Failed to send APNS batch: %v", err)
		}
	}

	return nil
}

// ValidateConfig validates the push provider configuration
func (p *PushProvider) ValidateConfig() error {
	if p.fcmServerKey == "" && p.apnsKeyID == "" {
		return fmt.Errorf("either FCM server key or APNS credentials are required")
	}

	if p.apnsKeyID != "" {
		if p.apnsTeamID == "" {
			return fmt.Errorf("APNS team ID is required when using APNS")
		}
		if p.apnsKeyPath == "" {
			return fmt.Errorf("APNS key path is required when using APNS")
		}
		if p.bundleID == "" {
			return fmt.Errorf("bundle ID is required when using APNS")
		}
	}

	return nil
}

// GetProviderName returns the provider name
func (p *PushProvider) GetProviderName() string {
	return "push"
}

// SupportsChannels returns the channels this provider supports
func (p *PushProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{
		models.NotificationChannelPush,
		models.NotificationChannelWebPush,
	}
}

// sendFCM sends a notification via Firebase Cloud Messaging
func (p *PushProvider) sendFCM(ctx context.Context, notification *models.Notification) error {
	// This would use the Firebase Admin SDK to send push notifications
	// For now, it's a placeholder implementation
	log.Printf("FCM: Sending push notification to device %s with title '%s'",
		notification.DeviceToken, notification.Title)

	// Build FCM payload
	payload := map[string]interface{}{
		"to": notification.DeviceToken,
		"notification": map[string]interface{}{
			"title": notification.Title,
			"body":  notification.Body,
		},
		"data": map[string]interface{}{
			"notification_id": notification.ID.String(),
			"action_url":      notification.ActionURL,
		},
	}

	// In a real implementation, this would make an HTTP request to FCM
	_ = payload
	return nil
}

// sendAPNS sends a notification via Apple Push Notification Service
func (p *PushProvider) sendAPNS(ctx context.Context, notification *models.Notification) error {
	// This would use the APNS library to send push notifications
	// For now, it's a placeholder implementation
	log.Printf("APNS: Sending push notification to device %s with title '%s'",
		notification.DeviceToken, notification.Title)

	// Build APNS payload
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]interface{}{
				"title": notification.Title,
				"body":  notification.Body,
			},
			"badge": 1,
			"sound": "default",
		},
		"notification_id": notification.ID.String(),
		"action_url":      notification.ActionURL,
	}

	// In a real implementation, this would use the APNS library
	_ = payload
	return nil
}

// sendFCMBatch sends multiple notifications via FCM
func (p *PushProvider) sendFCMBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("FCM: Sending batch of %d push notifications", len(notifications))

	// In a real implementation, this would use FCM's multicast messaging
	for _, notification := range notifications {
		if err := p.sendFCM(ctx, notification); err != nil {
			log.Printf("Failed to send FCM notification to %s: %v", notification.DeviceToken, err)
		}
	}

	return nil
}

// sendAPNSBatch sends multiple notifications via APNS
func (p *PushProvider) sendAPNSBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("APNS: Sending batch of %d push notifications", len(notifications))

	// In a real implementation, this would use APNS batch sending
	for _, notification := range notifications {
		if err := p.sendAPNS(ctx, notification); err != nil {
			log.Printf("Failed to send APNS notification to %s: %v", notification.DeviceToken, err)
		}
	}

	return nil
}

// ExpoProvider implements push notifications using Expo's push service
type ExpoProvider struct {
	accessToken string
}

// NewExpoProvider creates a new Expo push provider
func NewExpoProvider(accessToken string) *ExpoProvider {
	return &ExpoProvider{
		accessToken: accessToken,
	}
}

// Send sends a push notification via Expo
func (e *ExpoProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("Expo: Sending push notification to %s with title '%s'",
		notification.DeviceToken, notification.Title)

	// Build Expo push message
	message := map[string]interface{}{
		"to":    notification.DeviceToken,
		"title": notification.Title,
		"body":  notification.Body,
		"data": map[string]interface{}{
			"notification_id": notification.ID.String(),
			"action_url":      notification.ActionURL,
		},
	}

	// In a real implementation, this would make an HTTP request to Expo's push service
	_ = message
	return nil
}

// SendBatch sends multiple notifications via Expo
func (e *ExpoProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Expo: Sending batch of %d push notifications", len(notifications))

	// Expo supports batch sending natively
	messages := make([]map[string]interface{}, len(notifications))
	for i, notification := range notifications {
		messages[i] = map[string]interface{}{
			"to":    notification.DeviceToken,
			"title": notification.Title,
			"body":  notification.Body,
			"data": map[string]interface{}{
				"notification_id": notification.ID.String(),
				"action_url":      notification.ActionURL,
			},
		}
	}

	// In a real implementation, this would send the batch to Expo
	_ = messages
	return nil
}

// ValidateConfig validates Expo configuration
func (e *ExpoProvider) ValidateConfig() error {
	if e.accessToken == "" {
		return fmt.Errorf("Expo access token is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (e *ExpoProvider) GetProviderName() string {
	return "expo"
}

// SupportsChannels returns supported channels
func (e *ExpoProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelPush}
}