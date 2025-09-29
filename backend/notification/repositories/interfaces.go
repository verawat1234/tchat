package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
)

// NotificationRepository defines the interface for notification data operations
type NotificationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	Update(ctx context.Context, notification *models.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	GetByChannel(ctx context.Context, channel models.NotificationType, limit, offset int) ([]*models.Notification, error)
	GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error)
	GetFailedNotifications(ctx context.Context, limit int) ([]*models.Notification, error)
	GetScheduledNotifications(ctx context.Context, before time.Time, limit int) ([]*models.Notification, error)

	// Status operations
	MarkAsDelivered(ctx context.Context, id uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID, reason string) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.DeliveryStatus) error

	// Bulk operations
	CreateBatch(ctx context.Context, notifications []*models.Notification) error
	MarkBatchAsDelivered(ctx context.Context, ids []uuid.UUID) error

	// Analytics and cleanup
	GetDeliveryStats(ctx context.Context, startDate, endDate time.Time) (*models.DeliveryStats, error)
	CleanupOldNotifications(ctx context.Context, before time.Time) error
	CountByUserAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error)
}

// TemplateRepository defines the interface for notification template operations
type TemplateRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, template *models.NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error)
	Update(ctx context.Context, template *models.NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByType(ctx context.Context, notificationType string) (*models.NotificationTemplate, error)
	GetByTypeAndLanguage(ctx context.Context, notificationType, language string) (*models.NotificationTemplate, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.NotificationTemplate, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.NotificationTemplate, error)
	GetActiveTemplates(ctx context.Context, limit, offset int) ([]*models.NotificationTemplate, error)

	// Template operations
	GetVariables(ctx context.Context, templateID uuid.UUID) ([]string, error)
	ValidateTemplate(ctx context.Context, template *models.NotificationTemplate) error
}

// SubscriptionRepository defines the interface for notification subscription operations
type SubscriptionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, subscription *models.NotificationSubscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationSubscription, error)
	Update(ctx context.Context, subscription *models.NotificationSubscription) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.NotificationSubscription, error)
	GetByUserAndChannel(ctx context.Context, userID uuid.UUID, channel models.NotificationChannel) (*models.NotificationSubscription, error)
	GetByChannel(ctx context.Context, channel models.NotificationChannel, limit, offset int) ([]*models.NotificationSubscription, error)
	GetActiveSubscriptions(ctx context.Context, userID uuid.UUID) ([]*models.NotificationSubscription, error)

	// Subscription operations
	EnableSubscription(ctx context.Context, id uuid.UUID) error
	DisableSubscription(ctx context.Context, id uuid.UUID) error
	UpdateEndpoint(ctx context.Context, id uuid.UUID, endpoint string) error
}

// PreferencesRepository defines the interface for notification preferences operations
type PreferencesRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, preferences *models.NotificationPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationPreferences, error)
	Update(ctx context.Context, preferences *models.NotificationPreferences) error
	Delete(ctx context.Context, userID uuid.UUID) error

	// Preference operations
	UpdateChannelPreference(ctx context.Context, userID uuid.UUID, channel models.NotificationChannel, enabled bool) error
	UpdateCategoryPreference(ctx context.Context, userID uuid.UUID, category models.NotificationCategory, enabled bool) error
	UpdateQuietHours(ctx context.Context, userID uuid.UUID, quietHours map[string]interface{}) error
	IsChannelEnabled(ctx context.Context, userID uuid.UUID, channel models.NotificationChannel) (bool, error)
	IsCategoryEnabled(ctx context.Context, userID uuid.UUID, category models.NotificationCategory) (bool, error)
}