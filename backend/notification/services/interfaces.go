package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
)

// NotificationService defines the interface for notification business logic
type NotificationService interface {
	// Notification operations
	SendNotification(ctx context.Context, params CreateNotificationParams) (*models.Notification, error)
	SendBulkNotifications(ctx context.Context, params BulkNotificationParams) ([]*models.Notification, error)
	BroadcastNotification(ctx context.Context, params BroadcastParams) (string, int64, error)

	// Notification queries
	GetNotifications(ctx context.Context, params GetNotificationsParams) ([]*models.Notification, int64, error)
	GetNotification(ctx context.Context, notificationID, userID string) (*models.Notification, error)

	// Notification status
	MarkAsRead(ctx context.Context, notificationID, userID string) error
	DeleteNotification(ctx context.Context, notificationID, userID string) error
	RetryFailedNotification(ctx context.Context, notificationID string) error

	// Subscription management
	GetSubscriptions(ctx context.Context, userID string) ([]*models.NotificationSubscription, error)
	CreateSubscription(ctx context.Context, params CreateSubscriptionParams) (*models.NotificationSubscription, error)
	UpdateSubscription(ctx context.Context, params UpdateSubscriptionParams) (*models.NotificationSubscription, error)
	DeleteSubscription(ctx context.Context, subscriptionID, userID string) error

	// Template management
	GetTemplates(ctx context.Context, params GetTemplatesParams) ([]*models.NotificationTemplate, int64, error)
	GetTemplate(ctx context.Context, templateID string) (*models.NotificationTemplate, error)
	CreateTemplate(ctx context.Context, params CreateTemplateParams) (*models.NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, params UpdateTemplateParams) (*models.NotificationTemplate, error)
	DeleteTemplate(ctx context.Context, templateID string) error

	// Preferences management
	GetPreferences(ctx context.Context, userID string) (*models.NotificationPreferences, error)
	UpdatePreferences(ctx context.Context, params UpdatePreferencesParams) (*models.NotificationPreferences, error)

	// Analytics and reporting
	GetAnalytics(ctx context.Context, params AnalyticsParams) (*models.NotificationAnalytics, error)
	GetDeliveryReports(ctx context.Context, params DeliveryReportsParams) ([]*models.DeliveryReport, int64, error)
	GetQueueStatus(ctx context.Context) (*models.QueueStatus, error)

	// Webhook processing
	ProcessWebhook(ctx context.Context, params WebhookParams) error

	// Background processing
	ProcessPendingNotifications(ctx context.Context) error
	ProcessScheduledNotifications(ctx context.Context) error
	CleanupExpiredNotifications(ctx context.Context) error
}

// CacheService defines the interface for caching operations
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMulti(ctx context.Context, items map[string]interface{}, expiration time.Duration) error
	FlushAll(ctx context.Context) error
}

// EventService defines the interface for event publishing
type EventService interface {
	Publish(ctx context.Context, event interface{}) error
	PublishNotificationSent(ctx context.Context, notification *models.Notification) error
	PublishNotificationDelivered(ctx context.Context, notificationID uuid.UUID) error
	PublishNotificationFailed(ctx context.Context, notificationID uuid.UUID, reason string) error
	PublishNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error
}

// ChannelProvider defines the interface for notification channel providers
type ChannelProvider interface {
	Send(ctx context.Context, notification *models.Notification) error
	SendBatch(ctx context.Context, notifications []*models.Notification) error
	ValidateConfig() error
	GetProviderName() string
	SupportsChannels() []models.NotificationChannel
}

// Request/Response parameter structures

// CreateNotificationParams contains parameters for creating a single notification
type CreateNotificationParams struct {
	SenderID     string                 `json:"sender_id"`
	RecipientID  string                 `json:"recipient_id"`
	Type         string                 `json:"type"`
	Channel      string                 `json:"channel"`
	TemplateID   *string                `json:"template_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Priority     string                 `json:"priority"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// BulkNotificationParams contains parameters for bulk notification sending
type BulkNotificationParams struct {
	SenderID     string                 `json:"sender_id"`
	RecipientIDs []string               `json:"recipient_ids"`
	Type         string                 `json:"type"`
	Channel      string                 `json:"channel"`
	TemplateID   *string                `json:"template_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Priority     string                 `json:"priority"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
}

// BroadcastParams contains parameters for broadcasting notifications
type BroadcastParams struct {
	SenderID    string                 `json:"sender_id"`
	UserSegment string                 `json:"user_segment"`
	Type        string                 `json:"type"`
	Channel     string                 `json:"channel"`
	TemplateID  *string                `json:"template_id,omitempty"`
	Subject     *string                `json:"subject,omitempty"`
	Content     string                 `json:"content"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Priority    string                 `json:"priority"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

// GetNotificationsParams contains parameters for querying notifications
type GetNotificationsParams struct {
	UserID     string `json:"user_id"`
	Type       string `json:"type,omitempty"`
	Status     string `json:"status,omitempty"`
	UnreadOnly bool   `json:"unread_only"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// CreateSubscriptionParams contains parameters for creating subscriptions
type CreateSubscriptionParams struct {
	UserID      string                 `json:"user_id"`
	Channel     string                 `json:"channel"`
	Type        string                 `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Enabled     bool                   `json:"enabled"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
}

// UpdateSubscriptionParams contains parameters for updating subscriptions
type UpdateSubscriptionParams struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	Endpoint    *string                `json:"endpoint,omitempty"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
}

// GetTemplatesParams contains parameters for querying templates
type GetTemplatesParams struct {
	Type       string `json:"type,omitempty"`
	Category   string `json:"category,omitempty"`
	ActiveOnly bool   `json:"active_only"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// CreateTemplateParams contains parameters for creating templates
type CreateTemplateParams struct {
	CreatedBy string                 `json:"created_by"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Category  string                 `json:"category"`
	Subject   *string                `json:"subject,omitempty"`
	Content   string                 `json:"content"`
	Variables []string               `json:"variables,omitempty"`
	Locales   map[string]interface{} `json:"locales,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateTemplateParams contains parameters for updating templates
type UpdateTemplateParams struct {
	ID        string                 `json:"id"`
	UpdatedBy string                 `json:"updated_by"`
	Name      *string                `json:"name,omitempty"`
	Subject   *string                `json:"subject,omitempty"`
	Content   *string                `json:"content,omitempty"`
	Variables []string               `json:"variables,omitempty"`
	Locales   map[string]interface{} `json:"locales,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Active    *bool                  `json:"active,omitempty"`
}

// UpdatePreferencesParams contains parameters for updating preferences
type UpdatePreferencesParams struct {
	UserID       string                 `json:"user_id"`
	EmailEnabled *bool                  `json:"email_enabled,omitempty"`
	SMSEnabled   *bool                  `json:"sms_enabled,omitempty"`
	PushEnabled  *bool                  `json:"push_enabled,omitempty"`
	InAppEnabled *bool                  `json:"in_app_enabled,omitempty"`
	Categories   map[string]bool        `json:"categories,omitempty"`
	QuietHours   map[string]interface{} `json:"quiet_hours,omitempty"`
	Languages    []string               `json:"languages,omitempty"`
}

// AnalyticsParams contains parameters for analytics queries
type AnalyticsParams struct {
	Period    string `json:"period"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// DeliveryReportsParams contains parameters for delivery report queries
type DeliveryReportsParams struct {
	Status    string `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

// WebhookParams contains parameters for webhook processing
type WebhookParams struct {
	Provider  string                 `json:"provider"`
	Event     string                 `json:"event"`
	MessageID string                 `json:"message_id"`
	Status    string                 `json:"status"`
	Recipient string                 `json:"recipient"`
	ErrorCode *string                `json:"error_code,omitempty"`
	ErrorMsg  *string                `json:"error_message,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}