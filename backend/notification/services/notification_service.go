package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
)

// NotificationService provides notification delivery functionality
type NotificationService struct {
	notificationRepo NotificationRepository
	templateRepo     TemplateRepository
	cache           CacheService
	events          EventService
	providers       map[models.NotificationChannel]ChannelProvider
	config          *NotificationConfig
}

// NotificationConfig holds notification service configuration
type NotificationConfig struct {
	DefaultRetryAttempts   int
	RetryBackoffMultiplier float64
	MaxRetryDelay         time.Duration
	BatchSize             int
	EnableBatching        bool
	DefaultExpiry         time.Duration
	MaxNotificationsPerUser int
	EnableRateLimiting    bool
	RateLimit             int // per minute
	EnableDeduplication   bool
	DeduplicationWindow   time.Duration
}

// DefaultNotificationConfig returns default notification configuration
func DefaultNotificationConfig() *NotificationConfig {
	return &NotificationConfig{
		DefaultRetryAttempts:   3,
		RetryBackoffMultiplier: 2.0,
		MaxRetryDelay:         10 * time.Minute,
		BatchSize:             100,
		EnableBatching:        true,
		DefaultExpiry:         24 * time.Hour,
		MaxNotificationsPerUser: 1000,
		EnableRateLimiting:    true,
		RateLimit:             60, // 60 notifications per minute per user
		EnableDeduplication:   true,
		DeduplicationWindow:   5 * time.Minute,
	}
}

// Repository interfaces
type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	Update(ctx context.Context, notification *models.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPending(ctx context.Context, limit int) ([]*models.Notification, error)
	GetByStatus(ctx context.Context, status models.NotificationStatus, limit int) ([]*models.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.NotificationStatus) error
	MarkAsRead(ctx context.Context, id, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
}

type TemplateRepository interface {
	Create(ctx context.Context, template *models.NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error)
	GetByType(ctx context.Context, notificationType string) (*models.NotificationTemplate, error)
	Update(ctx context.Context, template *models.NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]*models.NotificationTemplate, error)
}

// Service interfaces
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	SetWithExpiry(ctx context.Context, key string, value interface{}, expiry time.Duration) error
}

type EventService interface {
	PublishNotification(ctx context.Context, event *NotificationEvent) error
	SubscribeToEvents(ctx context.Context, handler EventHandler) error
}

// Channel provider interfaces
type ChannelProvider interface {
	Send(ctx context.Context, notification *models.Notification) error
	SendBatch(ctx context.Context, notifications []*models.Notification) error
	GetDeliveryStatus(ctx context.Context, externalID string) (*DeliveryStatus, error)
	ValidateConfiguration() error
}

type PushProvider interface {
	ChannelProvider
	SendToToken(ctx context.Context, token string, payload *PushPayload) error
	SendToTopic(ctx context.Context, topic string, payload *PushPayload) error
}

type EmailProvider interface {
	ChannelProvider
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendBulkEmail(ctx context.Context, recipients []EmailRecipient) error
}

type SMSProvider interface {
	ChannelProvider
	SendSMS(ctx context.Context, phone, message string, country string) error
	SendBulkSMS(ctx context.Context, messages []SMSMessage) error
}

type InAppProvider interface {
	ChannelProvider
	DeliverToUser(ctx context.Context, userID uuid.UUID, notification *InAppNotification) error
}

// Event types
type NotificationEvent struct {
	Type           string                 `json:"type"`
	NotificationID uuid.UUID              `json:"notification_id"`
	UserID         uuid.UUID              `json:"user_id"`
	Channel        models.NotificationChannel `json:"channel"`
	Status         models.NotificationStatus  `json:"status"`
	Notification   *models.Notification   `json:"notification,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type EventHandler func(ctx context.Context, event *NotificationEvent) error

// Request/Response types
type SendNotificationRequest struct {
	UserID      uuid.UUID                     `json:"user_id"`
	Type        string                        `json:"type"`
	Title       string                        `json:"title"`
	Body        string                        `json:"body"`
	Channels    []models.NotificationChannel  `json:"channels"`
	Priority    models.NotificationPriority   `json:"priority"`
	Data        map[string]interface{}        `json:"data,omitempty"`
	ScheduledAt *time.Time                    `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time                    `json:"expires_at,omitempty"`
	Template    *string                       `json:"template,omitempty"`
	Variables   map[string]string             `json:"variables,omitempty"`
}

type BulkNotificationRequest struct {
	UserIDs     []uuid.UUID                   `json:"user_ids"`
	Type        string                        `json:"type"`
	Title       string                        `json:"title"`
	Body        string                        `json:"body"`
	Channels    []models.NotificationChannel  `json:"channels"`
	Priority    models.NotificationPriority   `json:"priority"`
	Data        map[string]interface{}        `json:"data,omitempty"`
	ScheduledAt *time.Time                    `json:"scheduled_at,omitempty"`
	ExpiresAt   *time.Time                    `json:"expires_at,omitempty"`
	Template    *string                       `json:"template,omitempty"`
	Variables   map[string]string             `json:"variables,omitempty"`
}

type PushPayload struct {
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Badge    *int                   `json:"badge,omitempty"`
	Sound    *string                `json:"sound,omitempty"`
	ImageURL *string                `json:"image_url,omitempty"`
}

type EmailRecipient struct {
	Email   string            `json:"email"`
	Name    string            `json:"name"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	IsHTML  bool              `json:"is_html"`
	Data    map[string]string `json:"data,omitempty"`
}

type SMSMessage struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
	Country string `json:"country"`
}

type InAppNotification struct {
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	ImageURL *string                `json:"image_url,omitempty"`
	ActionURL *string               `json:"action_url,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

type DeliveryStatus struct {
	ExternalID    string    `json:"external_id"`
	Status        string    `json:"status"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty"`
	FailedAt      *time.Time `json:"failed_at,omitempty"`
	ErrorMessage  *string   `json:"error_message,omitempty"`
	RetryCount    int       `json:"retry_count"`
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo NotificationRepository,
	templateRepo TemplateRepository,
	cache CacheService,
	events EventService,
	config *NotificationConfig,
) *NotificationService {
	if config == nil {
		config = DefaultNotificationConfig()
	}

	return &NotificationService{
		notificationRepo: notificationRepo,
		templateRepo:     templateRepo,
		cache:           cache,
		events:          events,
		providers:       make(map[models.NotificationChannel]ChannelProvider),
		config:          config,
	}
}

// RegisterProvider registers a channel provider
func (n *NotificationService) RegisterProvider(channel models.NotificationChannel, provider ChannelProvider) error {
	if err := provider.ValidateConfiguration(); err != nil {
		return fmt.Errorf("invalid provider configuration for %s: %v", channel, err)
	}

	n.providers[channel] = provider
	return nil
}

// SendNotification sends a single notification
func (n *NotificationService) SendNotification(ctx context.Context, req *SendNotificationRequest) (*models.Notification, error) {
	// Validate request
	if err := n.validateSendRequest(req); err != nil {
		return nil, fmt.Errorf("invalid send notification request: %v", err)
	}

	// Check rate limiting
	if err := n.checkRateLimit(ctx, req.UserID); err != nil {
		return nil, err
	}

	// Check deduplication
	if n.config.EnableDeduplication {
		if isDuplicate, err := n.checkDuplication(ctx, req); err != nil {
			return nil, fmt.Errorf("deduplication check failed: %v", err)
		} else if isDuplicate {
			return nil, fmt.Errorf("duplicate notification detected within deduplication window")
		}
	}

	// Process template if provided
	title, body, err := n.processTemplate(ctx, req.Template, req.Title, req.Body, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %v", err)
	}

	// Create notification
	notification := &models.Notification{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       title,
		Body:        body,
		Channels:    req.Channels,
		Priority:    req.Priority,
		Data:        req.Data,
		ScheduledAt: req.ScheduledAt,
		ExpiresAt:   req.ExpiresAt,
		Status:      models.NotificationStatusPending,
	}

	// Set up notification before creation
	if err := notification.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("failed to prepare notification: %v", err)
	}

	// Validate notification
	if err := notification.Validate(); err != nil {
		return nil, fmt.Errorf("notification validation failed: %v", err)
	}

	// Save notification
	if err := n.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %v", err)
	}

	// Send immediately if not scheduled
	if req.ScheduledAt == nil || req.ScheduledAt.Before(time.Now().UTC()) {
		go n.deliverNotification(context.Background(), notification)
	}

	// Update rate limiting counter
	n.updateRateLimit(ctx, req.UserID)

	// Add to deduplication cache
	if n.config.EnableDeduplication {
		n.addToDeduplicationCache(ctx, req)
	}

	return notification, nil
}

// SendBulkNotification sends notifications to multiple users
func (n *NotificationService) SendBulkNotification(ctx context.Context, req *BulkNotificationRequest) ([]*models.Notification, error) {
	// Validate request
	if err := n.validateBulkRequest(req); err != nil {
		return nil, fmt.Errorf("invalid bulk notification request: %v", err)
	}

	var notifications []*models.Notification
	var errors []string

	// Process template once if provided
	title, body, err := n.processTemplate(ctx, req.Template, req.Title, req.Body, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %v", err)
	}

	// Create notifications for each user
	for _, userID := range req.UserIDs {
		// Check rate limiting for each user
		if err := n.checkRateLimit(ctx, userID); err != nil {
			errors = append(errors, fmt.Sprintf("rate limit exceeded for user %s", userID))
			continue
		}

		// Create notification
		notification := &models.Notification{
			UserID:      userID,
			Type:        req.Type,
			Title:       title,
			Body:        body,
			Channels:    req.Channels,
			Priority:    req.Priority,
			Data:        req.Data,
			ScheduledAt: req.ScheduledAt,
			ExpiresAt:   req.ExpiresAt,
			Status:      models.NotificationStatusPending,
		}

		// Set up notification before creation
		if err := notification.BeforeCreate(); err != nil {
			errors = append(errors, fmt.Sprintf("failed to prepare notification for user %s: %v", userID, err))
			continue
		}

		// Save notification
		if err := n.notificationRepo.Create(ctx, notification); err != nil {
			errors = append(errors, fmt.Sprintf("failed to create notification for user %s: %v", userID, err))
			continue
		}

		notifications = append(notifications, notification)

		// Update rate limiting counter
		n.updateRateLimit(ctx, userID)
	}

	// Send notifications if not scheduled
	if req.ScheduledAt == nil || req.ScheduledAt.Before(time.Now().UTC()) {
		if n.config.EnableBatching {
			go n.deliverNotificationsBatch(context.Background(), notifications)
		} else {
			for _, notification := range notifications {
				go n.deliverNotification(context.Background(), notification)
			}
		}
	}

	// Return results with any errors
	if len(errors) > 0 {
		return notifications, fmt.Errorf("some notifications failed: %s", strings.Join(errors, "; "))
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (n *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	// Get notification
	notification, err := n.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("notification not found: %v", err)
	}

	// Verify ownership
	if notification.UserID != userID {
		return fmt.Errorf("unauthorized to mark this notification as read")
	}

	// Mark as read
	if err := n.notificationRepo.MarkAsRead(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %v", err)
	}

	// Publish event
	event := &NotificationEvent{
		Type:           "notification_read",
		NotificationID: notificationID,
		UserID:         userID,
		Channel:        models.NotificationChannelInApp,
		Status:         models.NotificationStatusDelivered,
		Timestamp:      time.Now().UTC(),
	}
	n.events.PublishNotification(ctx, event)

	return nil
}

// GetUserNotifications retrieves notifications for a user
func (n *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	return n.notificationRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetUnreadCount gets the number of unread notifications for a user
func (n *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return n.notificationRepo.GetUnreadCount(ctx, userID)
}

// ProcessScheduledNotifications processes notifications that are scheduled to be sent
func (n *NotificationService) ProcessScheduledNotifications(ctx context.Context) error {
	// Get pending notifications
	notifications, err := n.notificationRepo.GetPending(ctx, n.config.BatchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %v", err)
	}

	now := time.Now().UTC()
	var readyNotifications []*models.Notification

	for _, notification := range notifications {
		// Check if notification is expired
		if notification.IsExpired() {
			notification.Status = models.NotificationStatusExpired
			n.notificationRepo.Update(ctx, notification)
			continue
		}

		// Check if notification is ready to be sent
		if notification.ScheduledAt == nil || notification.ScheduledAt.Before(now) {
			readyNotifications = append(readyNotifications, notification)
		}
	}

	// Deliver ready notifications
	if len(readyNotifications) > 0 {
		if n.config.EnableBatching {
			n.deliverNotificationsBatch(ctx, readyNotifications)
		} else {
			for _, notification := range readyNotifications {
				go n.deliverNotification(context.Background(), notification)
			}
		}
	}

	return nil
}

// RetryFailedNotifications retries failed notifications
func (n *NotificationService) RetryFailedNotifications(ctx context.Context) error {
	// Get failed notifications
	notifications, err := n.notificationRepo.GetByStatus(ctx, models.NotificationStatusFailed, n.config.BatchSize)
	if err != nil {
		return fmt.Errorf("failed to get failed notifications: %v", err)
	}

	for _, notification := range notifications {
		// Check if notification can be retried
		if notification.CanRetry(n.config.DefaultRetryAttempts) {
			// Calculate backoff delay
			delay := n.calculateBackoffDelay(notification.RetryCount)
			if delay > n.config.MaxRetryDelay {
				delay = n.config.MaxRetryDelay
			}

			// Check if enough time has passed since last attempt
			if time.Since(notification.UpdatedAt) >= delay {
				// Reset status and retry
				notification.Status = models.NotificationStatusPending
				notification.RetryCount++
				notification.UpdatedAt = time.Now().UTC()

				if err := n.notificationRepo.Update(ctx, notification); err != nil {
					fmt.Printf("Warning: failed to update notification for retry: %v\n", err)
					continue
				}

				go n.deliverNotification(context.Background(), notification)
			}
		}
	}

	return nil
}

// CleanupExpiredNotifications removes expired notifications
func (n *NotificationService) CleanupExpiredNotifications(ctx context.Context) (int64, error) {
	return n.notificationRepo.DeleteExpired(ctx)
}

// Template management methods

// CreateTemplate creates a new notification template
func (n *NotificationService) CreateTemplate(ctx context.Context, template *models.NotificationTemplate) error {
	// Validate template
	if err := template.Validate(); err != nil {
		return fmt.Errorf("template validation failed: %v", err)
	}

	return n.templateRepo.Create(ctx, template)
}

// GetTemplate retrieves a template by type
func (n *NotificationService) GetTemplate(ctx context.Context, notificationType string) (*models.NotificationTemplate, error) {
	return n.templateRepo.GetByType(ctx, notificationType)
}

// Helper methods

func (n *NotificationService) validateSendRequest(req *SendNotificationRequest) error {
	if req.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if strings.TrimSpace(req.Type) == "" {
		return errors.New("notification type is required")
	}

	if strings.TrimSpace(req.Title) == "" && req.Template == nil {
		return errors.New("title is required when not using template")
	}

	if strings.TrimSpace(req.Body) == "" && req.Template == nil {
		return errors.New("body is required when not using template")
	}

	if len(req.Channels) == 0 {
		return errors.New("at least one channel is required")
	}

	for _, channel := range req.Channels {
		if !channel.IsValid() {
			return fmt.Errorf("invalid channel: %s", channel)
		}
	}

	if !req.Priority.IsValid() {
		return fmt.Errorf("invalid priority: %s", req.Priority)
	}

	return nil
}

func (n *NotificationService) validateBulkRequest(req *BulkNotificationRequest) error {
	if len(req.UserIDs) == 0 {
		return errors.New("user_ids is required")
	}

	if strings.TrimSpace(req.Type) == "" {
		return errors.New("notification type is required")
	}

	if strings.TrimSpace(req.Title) == "" && req.Template == nil {
		return errors.New("title is required when not using template")
	}

	if strings.TrimSpace(req.Body) == "" && req.Template == nil {
		return errors.New("body is required when not using template")
	}

	if len(req.Channels) == 0 {
		return errors.New("at least one channel is required")
	}

	return nil
}

func (n *NotificationService) checkRateLimit(ctx context.Context, userID uuid.UUID) error {
	if !n.config.EnableRateLimiting {
		return nil
	}

	key := fmt.Sprintf("notification_rate_limit:%s", userID)
	count, err := n.cache.Get(ctx, key)
	if err != nil {
		// If key doesn't exist, it's the first notification
		return nil
	}

	if countInt, ok := count.(int64); ok && countInt >= int64(n.config.RateLimit) {
		return fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	return nil
}

func (n *NotificationService) updateRateLimit(ctx context.Context, userID uuid.UUID) {
	if !n.config.EnableRateLimiting {
		return
	}

	key := fmt.Sprintf("notification_rate_limit:%s", userID)
	n.cache.Incr(ctx, key)
	n.cache.SetWithExpiry(ctx, key, 1, time.Minute)
}

func (n *NotificationService) checkDuplication(ctx context.Context, req *SendNotificationRequest) (bool, error) {
	// Create deduplication key based on user, type, title, and body
	key := fmt.Sprintf("notification_dedup:%s:%s:%s:%s", req.UserID, req.Type, req.Title, req.Body)

	_, err := n.cache.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, not a duplicate
		return false, nil
	}

	// Key exists, it's a duplicate
	return true, nil
}

func (n *NotificationService) addToDeduplicationCache(ctx context.Context, req *SendNotificationRequest) {
	key := fmt.Sprintf("notification_dedup:%s:%s:%s:%s", req.UserID, req.Type, req.Title, req.Body)
	n.cache.SetWithExpiry(ctx, key, true, n.config.DeduplicationWindow)
}

func (n *NotificationService) processTemplate(ctx context.Context, templateType *string, defaultTitle, defaultBody string, variables map[string]string) (string, string, error) {
	if templateType == nil {
		return defaultTitle, defaultBody, nil
	}

	// Get template
	template, err := n.templateRepo.GetByType(ctx, *templateType)
	if err != nil {
		return "", "", fmt.Errorf("template not found: %v", err)
	}

	// Process title template
	title := template.TitleTemplate
	body := template.BodyTemplate

	// Simple variable replacement
	if variables != nil {
		for key, value := range variables {
			placeholder := fmt.Sprintf("{{%s}}", key)
			title = strings.ReplaceAll(title, placeholder, value)
			body = strings.ReplaceAll(body, placeholder, value)
		}
	}

	return title, body, nil
}

func (n *NotificationService) deliverNotification(ctx context.Context, notification *models.Notification) {
	// Update status to processing
	notification.Status = models.NotificationStatusProcessing
	n.notificationRepo.Update(ctx, notification)

	success := false
	var lastError error

	// Deliver to each channel
	for _, channel := range notification.Channels {
		provider, exists := n.providers[channel]
		if !exists {
			lastError = fmt.Errorf("no provider configured for channel %s", channel)
			continue
		}

		if err := provider.Send(ctx, notification); err != nil {
			lastError = err
			fmt.Printf("Failed to deliver notification %s via %s: %v\n", notification.ID, channel, err)
			continue
		}

		success = true
	}

	// Update final status
	if success {
		notification.Status = models.NotificationStatusDelivered
		notification.DeliveredAt = func() *time.Time { t := time.Now().UTC(); return &t }()
	} else {
		notification.Status = models.NotificationStatusFailed
		if lastError != nil {
			errorMsg := lastError.Error()
			notification.ErrorMessage = &errorMsg
		}
	}

	n.notificationRepo.Update(ctx, notification)

	// Publish delivery event
	event := &NotificationEvent{
		Type:           "notification_delivered",
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Status:         notification.Status,
		Notification:   notification,
		Timestamp:      time.Now().UTC(),
	}
	n.events.PublishNotification(ctx, event)
}

func (n *NotificationService) deliverNotificationsBatch(ctx context.Context, notifications []*models.Notification) {
	// Group notifications by channel
	channelGroups := make(map[models.NotificationChannel][]*models.Notification)

	for _, notification := range notifications {
		for _, channel := range notification.Channels {
			channelGroups[channel] = append(channelGroups[channel], notification)
		}
	}

	// Send batches to each channel
	for channel, channelNotifications := range channelGroups {
		provider, exists := n.providers[channel]
		if !exists {
			// Mark all notifications as failed
			for _, notification := range channelNotifications {
				notification.Status = models.NotificationStatusFailed
				errorMsg := fmt.Sprintf("no provider configured for channel %s", channel)
				notification.ErrorMessage = &errorMsg
				n.notificationRepo.Update(ctx, notification)
			}
			continue
		}

		// Send batch
		if err := provider.SendBatch(ctx, channelNotifications); err != nil {
			// Mark all notifications as failed
			for _, notification := range channelNotifications {
				notification.Status = models.NotificationStatusFailed
				errorMsg := err.Error()
				notification.ErrorMessage = &errorMsg
				n.notificationRepo.Update(ctx, notification)
			}
		} else {
			// Mark all notifications as delivered
			for _, notification := range channelNotifications {
				notification.Status = models.NotificationStatusDelivered
				notification.DeliveredAt = func() *time.Time { t := time.Now().UTC(); return &t }()
				n.notificationRepo.Update(ctx, notification)
			}
		}
	}
}

func (n *NotificationService) calculateBackoffDelay(retryCount int) time.Duration {
	// Exponential backoff: base_delay * (multiplier ^ retry_count)
	baseDelay := 30 * time.Second
	delay := float64(baseDelay) * (n.config.RetryBackoffMultiplier * float64(retryCount))
	return time.Duration(delay)
}