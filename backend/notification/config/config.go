package config

import (
	"fmt"

	"tchat.dev/shared/config"
)

// NotificationConfig holds notification service-specific configuration
type NotificationConfig struct {
	*config.Config

	// Notification-specific settings
	Notification struct {
		// Queue settings
		QueueSize          int
		WorkerCount        int
		RetryAttempts      int
		RetryDelay         int // seconds

		// Rate limiting
		RateLimitPerUser   int // notifications per user per minute
		RateLimitGlobal    int // total notifications per minute

		// Webhook settings
		WebhookSecret      string
		WebhookTimeout     int // seconds

		// Provider settings
		DefaultProvider    string
		EnabledProviders   []string

		// Template settings
		TemplateCache      bool
		TemplateCacheTTL   int // seconds

		// Batch settings
		BatchSize          int
		BatchTimeout       int // seconds
	}
}

// Load loads notification service configuration
func Load() (*NotificationConfig, error) {
	// Load base configuration
	baseConfig, err := config.LoadWithServicePort("notification", 8085)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	notifConfig := &NotificationConfig{
		Config: baseConfig,
	}

	// Load notification-specific settings from environment
	notifConfig.Notification.QueueSize = config.GetIntEnv("NOTIFICATION_QUEUE_SIZE", 10000)
	notifConfig.Notification.WorkerCount = config.GetIntEnv("NOTIFICATION_WORKER_COUNT", 10)
	notifConfig.Notification.RetryAttempts = config.GetIntEnv("NOTIFICATION_RETRY_ATTEMPTS", 3)
	notifConfig.Notification.RetryDelay = config.GetIntEnv("NOTIFICATION_RETRY_DELAY", 5)

	notifConfig.Notification.RateLimitPerUser = config.GetIntEnv("NOTIFICATION_RATE_LIMIT_USER", 100)
	notifConfig.Notification.RateLimitGlobal = config.GetIntEnv("NOTIFICATION_RATE_LIMIT_GLOBAL", 10000)

	notifConfig.Notification.WebhookSecret = config.GetEnv("NOTIFICATION_WEBHOOK_SECRET", "")
	notifConfig.Notification.WebhookTimeout = config.GetIntEnv("NOTIFICATION_WEBHOOK_TIMEOUT", 30)

	notifConfig.Notification.DefaultProvider = config.GetEnv("NOTIFICATION_DEFAULT_PROVIDER", "email")
	notifConfig.Notification.EnabledProviders = []string{"email", "sms", "push", "in_app"}

	notifConfig.Notification.TemplateCache = config.GetBoolEnv("NOTIFICATION_TEMPLATE_CACHE", true)
	notifConfig.Notification.TemplateCacheTTL = config.GetIntEnv("NOTIFICATION_TEMPLATE_CACHE_TTL", 3600)

	notifConfig.Notification.BatchSize = config.GetIntEnv("NOTIFICATION_BATCH_SIZE", 100)
	notifConfig.Notification.BatchTimeout = config.GetIntEnv("NOTIFICATION_BATCH_TIMEOUT", 5)

	return notifConfig, nil
}

// MustLoad loads configuration and panics on error
func MustLoad() *NotificationConfig {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load notification config: %v", err))
	}
	return cfg
}
