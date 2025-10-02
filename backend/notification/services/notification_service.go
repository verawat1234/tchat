package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
	"tchat.dev/notification/repositories"
)

// NotificationServiceImpl provides notification delivery functionality
type NotificationServiceImpl struct {
	notificationRepo repositories.NotificationRepository
	templateRepo     repositories.TemplateRepository
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

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	templateRepo repositories.TemplateRepository,
	cache CacheService,
	events EventService,
	config *NotificationConfig,
) NotificationService {
	service := &NotificationServiceImpl{
		notificationRepo: notificationRepo,
		templateRepo:     templateRepo,
		cache:           cache,
		events:          events,
		providers:       make(map[models.NotificationChannel]ChannelProvider),
		config:          config,
	}

	if config == nil {
		service.config = DefaultNotificationConfig()
	}

	return service
}

// PushProvider interface extends ChannelProvider for push notifications
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

// RegisterProvider registers a channel provider
func (n *NotificationServiceImpl) RegisterProvider(channel models.NotificationChannel, provider ChannelProvider) error {
	if err := provider.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid provider configuration for %s: %v", channel, err)
	}

	n.providers[channel] = provider
	return nil
}

// sendNotificationInternal sends a single notification (internal method)
func (n *NotificationServiceImpl) sendNotificationInternal(ctx context.Context, req *SendNotificationRequest) (*models.Notification, error) {
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
		Title:       title,
		Body:        body,
		ScheduledAt: req.ScheduledAt,
		ExpiresAt:   req.ExpiresAt,
		Status:      models.DeliveryStatusPending,
		Metadata:    req.Data,
	}

	// Set up notification with defaults
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

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
func (n *NotificationServiceImpl) SendBulkNotification(ctx context.Context, req *BulkNotificationRequest) ([]*models.Notification, error) {
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
			Title:       title,
			Body:        body,
			ScheduledAt: req.ScheduledAt,
			ExpiresAt:   req.ExpiresAt,
			Status:      models.DeliveryStatusPending,
			Metadata:    req.Data,
		}

		// Set up notification with defaults
		if notification.ID == uuid.Nil {
			notification.ID = uuid.New()
		}
		notification.CreatedAt = time.Now()
		notification.UpdatedAt = time.Now()

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
func (n *NotificationServiceImpl) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	// Parse UUIDs
	nID, err := uuid.Parse(notificationID)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %v", err)
	}

	uID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	// Get notification
	notification, err := n.notificationRepo.GetByID(ctx, nID)
	if err != nil {
		return fmt.Errorf("notification not found: %v", err)
	}

	// Verify ownership
	if notification.UserID != uID {
		return fmt.Errorf("unauthorized to mark this notification as read")
	}

	// Mark as read
	if err := n.notificationRepo.MarkAsRead(ctx, nID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %v", err)
	}

	// Publish event
	event := &NotificationEvent{
		Type:           "notification_read",
		NotificationID: nID,
		UserID:         uID,
		Channel:        models.NotificationChannelInApp,
		Timestamp:      time.Now().UTC(),
	}
	n.events.PublishNotification(ctx, event)

	return nil
}

// GetUserNotifications retrieves notifications for a user
func (n *NotificationServiceImpl) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	return n.notificationRepo.GetByUserID(ctx, userID, limit, offset)
}

// BroadcastNotification sends a notification to a broadcast audience
func (n *NotificationServiceImpl) BroadcastNotification(ctx context.Context, params BroadcastParams) (string, int64, error) {
	// Convert broadcast params to bulk notification request
	bulkReq := &BulkNotificationRequest{
		UserIDs:     []uuid.UUID{}, // Will be populated based on segment
		Type:        params.Type,
		Title:       *params.Subject,
		Body:        params.Content,
		ScheduledAt: params.ScheduledAt,
	}

	// Send bulk notifications
	notifications, err := n.SendBulkNotification(ctx, bulkReq)
	if err != nil {
		return "", 0, err
	}

	broadcastID := uuid.New().String()
	return broadcastID, int64(len(notifications)), nil
}

// ProcessScheduledNotifications processes notifications that are scheduled to be sent
func (n *NotificationServiceImpl) ProcessScheduledNotifications(ctx context.Context) error {
	// Get pending notifications
	notifications, err := n.notificationRepo.GetPendingNotifications(ctx, n.config.BatchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %v", err)
	}

	now := time.Now().UTC()
	var readyNotifications []*models.Notification

	for _, notification := range notifications {
		// Check if notification is expired
		if notification.IsExpired() {
			notification.Status = models.DeliveryStatusExpired
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
func (n *NotificationServiceImpl) RetryFailedNotifications(ctx context.Context) error {
	// Get failed notifications
	notifications, err := n.notificationRepo.GetFailedNotifications(ctx, n.config.BatchSize)
	if err != nil {
		return fmt.Errorf("failed to get failed notifications: %v", err)
	}

	for _, notification := range notifications {
		// Check if notification can be retried
		if notification.CanRetry() {
			// Calculate backoff delay
			delay := n.calculateBackoffDelay(notification.RetryCount)
			if delay > n.config.MaxRetryDelay {
				delay = n.config.MaxRetryDelay
			}

			// Check if enough time has passed since last attempt
			if time.Since(notification.UpdatedAt) >= delay {
				// Reset status and retry
				notification.Status = models.DeliveryStatusPending
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
func (n *NotificationServiceImpl) CleanupExpiredNotifications(ctx context.Context) error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	return n.notificationRepo.CleanupOldNotifications(ctx, thirtyDaysAgo)
}

// Template management methods (these are helper methods, not interface methods)

// Helper methods

func (n *NotificationServiceImpl) validateSendRequest(req *SendNotificationRequest) error {
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

func (n *NotificationServiceImpl) validateBulkRequest(req *BulkNotificationRequest) error {
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

func (n *NotificationServiceImpl) checkRateLimit(ctx context.Context, userID uuid.UUID) error {
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

func (n *NotificationServiceImpl) updateRateLimit(ctx context.Context, userID uuid.UUID) {
	if !n.config.EnableRateLimiting {
		return
	}

	key := fmt.Sprintf("notification_rate_limit:%s", userID)
	n.cache.Incr(ctx, key)
	n.cache.SetWithExpiry(ctx, key, 1, time.Minute)
}

func (n *NotificationServiceImpl) checkDuplication(ctx context.Context, req *SendNotificationRequest) (bool, error) {
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

func (n *NotificationServiceImpl) addToDeduplicationCache(ctx context.Context, req *SendNotificationRequest) {
	key := fmt.Sprintf("notification_dedup:%s:%s:%s:%s", req.UserID, req.Type, req.Title, req.Body)
	n.cache.SetWithExpiry(ctx, key, true, n.config.DeduplicationWindow)
}

func (n *NotificationServiceImpl) processTemplate(ctx context.Context, templateType *string, defaultTitle, defaultBody string, variables map[string]string) (string, string, error) {
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

func (n *NotificationServiceImpl) deliverNotification(ctx context.Context, notification *models.Notification) {
	// Update status to sent
	notification.Status = models.DeliveryStatusSent
	n.notificationRepo.Update(ctx, notification)

	success := false
	var lastError error

	// Deliver to configured channel if provider exists
	if notification.Channel != "" {
		provider, exists := n.providers[notification.Channel]
		if exists {
			if err := provider.Send(ctx, notification); err != nil {
				lastError = err
				fmt.Printf("Failed to deliver notification %s via %s: %v\n", notification.ID, notification.Channel, err)
			} else {
				success = true
			}
		} else {
			lastError = fmt.Errorf("no provider configured for channel %s", notification.Channel)
		}
	}

	// Update final status
	if success {
		notification.Status = models.DeliveryStatusDelivered
		now := time.Now().UTC()
		notification.SentAt = &now
	} else {
		notification.Status = models.DeliveryStatusFailed
		if lastError != nil {
			errorMsg := lastError.Error()
			notification.ErrorMessage = errorMsg
		}
	}

	n.notificationRepo.Update(ctx, notification)

	// Publish delivery event
	event := &NotificationEvent{
		Type:           "notification_delivered",
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Channel:        notification.Channel,
		Notification:   notification,
		Timestamp:      time.Now().UTC(),
	}
	n.events.PublishNotification(ctx, event)
}

func (n *NotificationServiceImpl) deliverNotificationsBatch(ctx context.Context, notifications []*models.Notification) {
	// Group notifications by channel
	channelGroups := make(map[models.NotificationChannel][]*models.Notification)

	for _, notification := range notifications {
		if notification.Channel != "" {
			channelGroups[notification.Channel] = append(channelGroups[notification.Channel], notification)
		}
	}

	// Send batches to each channel
	for channel, channelNotifications := range channelGroups {
		provider, exists := n.providers[channel]
		if !exists {
			// Mark all notifications as failed
			for _, notification := range channelNotifications {
				notification.Status = models.DeliveryStatusFailed
				errorMsg := fmt.Sprintf("no provider configured for channel %s", channel)
				notification.ErrorMessage = errorMsg
				n.notificationRepo.Update(ctx, notification)
			}
			continue
		}

		// Send batch
		if err := provider.SendBatch(ctx, channelNotifications); err != nil {
			// Mark all notifications as failed
			for _, notification := range channelNotifications {
				notification.Status = models.DeliveryStatusFailed
				errorMsg := err.Error()
				notification.ErrorMessage = errorMsg
				n.notificationRepo.Update(ctx, notification)
			}
		} else {
			// Mark all notifications as delivered
			for _, notification := range channelNotifications {
				notification.Status = models.DeliveryStatusDelivered
				now := time.Now().UTC()
				notification.SentAt = &now
				n.notificationRepo.Update(ctx, notification)
			}
		}
	}
}

func (n *NotificationServiceImpl) calculateBackoffDelay(retryCount int) time.Duration {
	// Exponential backoff: base_delay * (multiplier ^ retry_count)
	baseDelay := 30 * time.Second
	delay := float64(baseDelay) * (n.config.RetryBackoffMultiplier * float64(retryCount))
	return time.Duration(delay)
}
// Stub implementations for NotificationService interface methods
// These should be implemented with proper business logic

func (n *NotificationServiceImpl) GetNotifications(ctx context.Context, params GetNotificationsParams) ([]*models.Notification, int64, error) {
	userID, _ := uuid.Parse(params.UserID)
	notifications, err := n.notificationRepo.GetByUserID(ctx, userID, params.Limit, params.Page*params.Limit)
	return notifications, 0, err
}

func (n *NotificationServiceImpl) GetNotification(ctx context.Context, notificationID, userID string) (*models.Notification, error) {
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return nil, err
	}
	return n.notificationRepo.GetByID(ctx, id)
}

func (n *NotificationServiceImpl) DeleteNotification(ctx context.Context, notificationID, userID string) error {
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return err
	}
	return n.notificationRepo.Delete(ctx, id)
}

func (n *NotificationServiceImpl) RetryFailedNotification(ctx context.Context, notificationID string) error {
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return err
	}
	notification, err := n.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	notification.Status = models.DeliveryStatusPending
	return n.notificationRepo.Update(ctx, notification)
}

func (n *NotificationServiceImpl) GetSubscriptions(ctx context.Context, userID string) ([]*models.NotificationSubscription, error) {
	// Stub implementation
	return []*models.NotificationSubscription{}, nil
}

func (n *NotificationServiceImpl) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) (*models.NotificationSubscription, error) {
	// Stub implementation
	return &models.NotificationSubscription{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) UpdateSubscription(ctx context.Context, params UpdateSubscriptionParams) (*models.NotificationSubscription, error) {
	// Stub implementation
	return &models.NotificationSubscription{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) DeleteSubscription(ctx context.Context, subscriptionID, userID string) error {
	// Stub implementation
	return fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) GetTemplate(ctx context.Context, templateID string) (*models.NotificationTemplate, error) {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return nil, err
	}
	return n.templateRepo.GetByID(ctx, id)
}

func (n *NotificationServiceImpl) GetTemplates(ctx context.Context, params GetTemplatesParams) ([]*models.NotificationTemplate, int64, error) {
	templates, err := n.templateRepo.GetAll(ctx, params.Limit, params.Page*params.Limit)
	return templates, 0, err
}

func (n *NotificationServiceImpl) CreateTemplate(ctx context.Context, params CreateTemplateParams) (*models.NotificationTemplate, error) {
	template := &models.NotificationTemplate{
		ID:          uuid.New(),
		Name:        params.Name,
		Type:        models.NotificationType(params.Type),
		Category:    models.NotificationCategory(params.Category),
		Subject:     *params.Subject,
		Body:        params.Content,
		Variables:   params.Variables,
		CreatedBy:   uuid.MustParse(params.CreatedBy),
	}

	if err := n.templateRepo.ValidateTemplate(ctx, template); err != nil {
		return nil, err
	}

	if err := n.templateRepo.Create(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

func (n *NotificationServiceImpl) UpdateTemplate(ctx context.Context, params UpdateTemplateParams) (*models.NotificationTemplate, error) {
	// Stub implementation
	return &models.NotificationTemplate{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) DeleteTemplate(ctx context.Context, templateID string) error {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return err
	}
	return n.templateRepo.Delete(ctx, id)
}

func (n *NotificationServiceImpl) GetPreferences(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
	// Stub implementation
	return &models.NotificationPreferences{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) UpdatePreferences(ctx context.Context, params UpdatePreferencesParams) (*models.NotificationPreferences, error) {
	// Stub implementation
	return &models.NotificationPreferences{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) GetAnalytics(ctx context.Context, params AnalyticsParams) (*models.NotificationAnalytics, error) {
	// Stub implementation
	return &models.NotificationAnalytics{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) GetDeliveryReports(ctx context.Context, params DeliveryReportsParams) ([]*models.DeliveryReport, int64, error) {
	// Stub implementation
	return []*models.DeliveryReport{}, 0, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) GetQueueStatus(ctx context.Context) (*models.QueueStatus, error) {
	// Stub implementation
	return &models.QueueStatus{}, fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) ProcessWebhook(ctx context.Context, params WebhookParams) error {
	// Stub implementation
	return fmt.Errorf("not implemented")
}

func (n *NotificationServiceImpl) ProcessPendingNotifications(ctx context.Context) error {
	// Stub implementation - reuse existing method
	return n.ProcessScheduledNotifications(ctx)
}

// SendBulkNotifications implements NotificationService interface
func (n *NotificationServiceImpl) SendBulkNotifications(ctx context.Context, params BulkNotificationParams) ([]*models.Notification, error) {
	// Convert string IDs to UUIDs
	userIDs := make([]uuid.UUID, len(params.RecipientIDs))
	for i, idStr := range params.RecipientIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient ID %s: %v", idStr, err)
		}
		userIDs[i] = id
	}

	// Convert params to BulkNotificationRequest
	req := &BulkNotificationRequest{
		UserIDs:     userIDs,
		Type:        params.Type,
		Title:       *params.Subject,
		Body:        params.Content,
		Variables:   convertVariables(params.Variables),
		ScheduledAt: params.ScheduledAt,
		ExpiresAt:   params.ExpiresAt,
	}

	return n.SendBulkNotification(ctx, req)
}

// SendNotification implements NotificationService interface with CreateNotificationParams
func (n *NotificationServiceImpl) SendNotification(ctx context.Context, params CreateNotificationParams) (*models.Notification, error) {
	// Convert params to SendNotificationRequest
	userID, err := uuid.Parse(params.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient ID: %v", err)
	}

	req := &SendNotificationRequest{
		UserID:      userID,
		Type:        params.Type,
		Title:       *params.Subject,
		Body:        params.Content,
		Variables:   convertVariables(params.Variables),
		ScheduledAt: params.ScheduledAt,
		ExpiresAt:   params.ExpiresAt,
	}

	return n.sendNotificationInternal(ctx, req)
}

func convertVariables(vars map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range vars {
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
