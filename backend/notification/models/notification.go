package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string
type NotificationCategory string
type NotificationChannel string
type NotificationStatus string
type NotificationPriority string
type Priority string
type DeliveryStatus string
type FailureReason string
type AudienceType string

const (
	// Notification Types
	NotificationTypePush          NotificationType = "push"
	NotificationTypeEmail         NotificationType = "email"
	NotificationTypeSMS           NotificationType = "sms"
	NotificationTypeInApp         NotificationType = "in_app"
	NotificationTypeWebPush       NotificationType = "web_push"
	NotificationTypeWebhook       NotificationType = "webhook"
	NotificationTypeSlack         NotificationType = "slack"
	NotificationTypeDiscord       NotificationType = "discord"
	NotificationTypeTelegram      NotificationType = "telegram"
	NotificationTypeWhatsApp      NotificationType = "whatsapp"
	NotificationTypeLine          NotificationType = "line" // Popular in Southeast Asia

	// Categories
	NotificationCategorySystem     NotificationCategory = "system"
	NotificationCategoryMarketing  NotificationCategory = "marketing"
	NotificationCategorySecurity   NotificationCategory = "security"
	NotificationCategoryPayment    NotificationCategory = "payment"
	NotificationCategoryOrder      NotificationCategory = "order"
	NotificationCategoryChat       NotificationCategory = "chat"
	NotificationCategoryPromotion  NotificationCategory = "promotion"
	NotificationCategoryUpdate     NotificationCategory = "update"
	NotificationCategoryAlert      NotificationCategory = "alert"
	NotificationCategoryReminder   NotificationCategory = "reminder"

	// Priority Levels
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
	PriorityUrgent   Priority = "urgent"

	// Delivery Status
	DeliveryStatusPending    DeliveryStatus = "pending"
	DeliveryStatusSent       DeliveryStatus = "sent"
	DeliveryStatusDelivered  DeliveryStatus = "delivered"
	DeliveryStatusRead       DeliveryStatus = "read"
	DeliveryStatusFailed     DeliveryStatus = "failed"
	DeliveryStatusCancelled  DeliveryStatus = "cancelled"
	DeliveryStatusExpired    DeliveryStatus = "expired"

	// Failure Reasons
	FailureReasonInvalidToken     FailureReason = "invalid_token"
	FailureReasonUserNotFound     FailureReason = "user_not_found"
	FailureReasonQuotaExceeded    FailureReason = "quota_exceeded"
	FailureReasonNetworkError     FailureReason = "network_error"
	FailureReasonProviderError    FailureReason = "provider_error"
	FailureReasonInvalidFormat    FailureReason = "invalid_format"
	FailureReasonBlocked          FailureReason = "blocked"
	FailureReasonOptedOut         FailureReason = "opted_out"
	FailureReasonRateLimited      FailureReason = "rate_limited"
	FailureReasonExpired          FailureReason = "expired"

	// Notification Channels
	NotificationChannelPush       NotificationChannel = "push"
	NotificationChannelEmail      NotificationChannel = "email"
	NotificationChannelSMS        NotificationChannel = "sms"
	NotificationChannelInApp      NotificationChannel = "in_app"
	NotificationChannelWebPush    NotificationChannel = "web_push"
	NotificationChannelWebhook    NotificationChannel = "webhook"

	// Notification Status
	NotificationStatusPending    NotificationStatus = "pending"
	NotificationStatusSent       NotificationStatus = "sent"
	NotificationStatusDelivered  NotificationStatus = "delivered"
	NotificationStatusFailed     NotificationStatus = "failed"
	NotificationStatusProcessing NotificationStatus = "processing"
	NotificationStatusExpired    NotificationStatus = "expired"

	// Notification Priority
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityNormal   NotificationPriority = "normal"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"

	// Audience Types
	AudienceTypeUser      AudienceType = "user"
	AudienceTypeSegment   AudienceType = "segment"
	AudienceTypeRole      AudienceType = "role"
	AudienceTypeBroadcast AudienceType = "broadcast"
	AudienceTypeLocation  AudienceType = "location"
	AudienceTypeDevice    AudienceType = "device"
)

type LocalizedContent struct {
	Language string `json:"language" gorm:"type:varchar(10)"`
	Title    string `json:"title" gorm:"type:text"`
	Body     string `json:"body" gorm:"type:text"`
	ImageURL string `json:"image_url,omitempty" gorm:"type:varchar(500)"`
	ActionURL string `json:"action_url,omitempty" gorm:"type:varchar(500)"`
}

type DeliveryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`
	TTL           time.Duration `json:"ttl"`
	BatchSize     int           `json:"batch_size"`
	RateLimit     int           `json:"rate_limit"` // per minute
}

type Targeting struct {
	AudienceType AudienceType   `json:"audience_type"`
	UserIDs      []uuid.UUID    `json:"user_ids,omitempty"`
	Segments     []string       `json:"segments,omitempty"`
	Roles        []string       `json:"roles,omitempty"`
	Countries    []string       `json:"countries,omitempty"`
	Cities       []string       `json:"cities,omitempty"`
	DeviceTypes  []string       `json:"device_types,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	Conditions   map[string]any `json:"conditions,omitempty"`
}

type Analytics struct {
	Sent         int64   `json:"sent"`
	Delivered    int64   `json:"delivered"`
	Read         int64   `json:"read"`
	Clicked      int64   `json:"clicked"`
	Failed       int64   `json:"failed"`
	DeliveryRate float64 `json:"delivery_rate"`
	ReadRate     float64 `json:"read_rate"`
	ClickRate    float64 `json:"click_rate"`
	FailureRate  float64 `json:"failure_rate"`
}

type Notification struct {
	ID          uuid.UUID            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title       string               `json:"title" gorm:"type:varchar(255);not null"`
	Body        string               `json:"body" gorm:"type:text"`
	ImageURL    string               `json:"image_url,omitempty" gorm:"type:varchar(500)"`
	ActionURL   string               `json:"action_url,omitempty" gorm:"type:varchar(500)"`
	Type        NotificationType     `json:"type" gorm:"type:varchar(20);not null;index"`
	Category    NotificationCategory `json:"category" gorm:"type:varchar(30);not null;index"`
	Priority    Priority             `json:"priority" gorm:"type:varchar(10);default:'normal';index"`
	Status      DeliveryStatus       `json:"status" gorm:"type:varchar(20);default:'pending';index"`

	// Targeting and Audience
	Targeting   Targeting    `json:"targeting" gorm:"type:json"`
	TotalUsers  int64        `json:"total_users" gorm:"default:0"`

	// Scheduling
	ScheduledAt *time.Time   `json:"scheduled_at,omitempty" gorm:"index"`
	SentAt      *time.Time   `json:"sent_at,omitempty" gorm:"index"`
	ExpiresAt   *time.Time   `json:"expires_at,omitempty" gorm:"index"`

	// Content
	LocalizedContent []LocalizedContent `json:"localized_content" gorm:"type:json"`
	Metadata         map[string]any     `json:"metadata" gorm:"type:json"`

	// Delivery Configuration
	Config       DeliveryConfig `json:"config" gorm:"type:json"`

	// Analytics and Tracking
	Analytics    Analytics      `json:"analytics" gorm:"type:json"`

	// Error Handling
	FailureReason FailureReason `json:"failure_reason,omitempty" gorm:"type:varchar(50)"`
	ErrorMessage  string        `json:"error_message,omitempty" gorm:"type:text"`
	RetryCount    int           `json:"retry_count" gorm:"default:0"`

	// Audit Fields
	CreatedBy   uuid.UUID  `json:"created_by" gorm:"type:varchar(36);index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type UserNotification struct {
	ID             uuid.UUID      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	NotificationID uuid.UUID      `json:"notification_id" gorm:"type:varchar(36);not null;index"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Status         DeliveryStatus `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	IsRead         bool           `json:"is_read" gorm:"default:false;index"`
	ReadAt         *time.Time     `json:"read_at,omitempty"`
	IsClicked      bool           `json:"is_clicked" gorm:"default:false"`
	ClickedAt      *time.Time     `json:"clicked_at,omitempty"`
	SentAt         *time.Time     `json:"sent_at,omitempty" gorm:"index"`
	DeliveredAt    *time.Time     `json:"delivered_at,omitempty"`
	FailureReason  FailureReason  `json:"failure_reason,omitempty" gorm:"type:varchar(50)"`
	ErrorMessage   string         `json:"error_message,omitempty" gorm:"type:text"`
	RetryCount     int            `json:"retry_count" gorm:"default:0"`
	DeviceToken    string         `json:"device_token,omitempty" gorm:"type:varchar(255)"`
	Platform       string         `json:"platform,omitempty" gorm:"type:varchar(20)"`
	Metadata       map[string]any `json:"metadata" gorm:"type:json"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

type NotificationTemplate struct {
	ID           uuid.UUID            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name         string               `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Type         NotificationType     `json:"type" gorm:"type:varchar(20);not null"`
	Category     NotificationCategory `json:"category" gorm:"type:varchar(30);not null"`
	Subject       string               `json:"subject" gorm:"type:varchar(255)"`
	Body          string               `json:"body" gorm:"type:text;not null"`
	TitleTemplate string               `json:"title_template" gorm:"type:text"`
	BodyTemplate  string               `json:"body_template" gorm:"type:text"`
	Variables     []string             `json:"variables" gorm:"type:json"`
	LocalizedVersions []LocalizedContent `json:"localized_versions" gorm:"type:json"`
	IsActive     bool                 `json:"is_active" gorm:"default:true"`
	Version      int                  `json:"version" gorm:"default:1"`
	Tags         []string             `json:"tags" gorm:"type:json"`
	CreatedBy    uuid.UUID            `json:"created_by" gorm:"type:varchar(36);index"`
	CreatedAt    time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.Config.MaxRetries == 0 {
		n.Config.MaxRetries = 3
	}
	if n.Config.RetryInterval == 0 {
		n.Config.RetryInterval = 5 * time.Minute
	}
	if n.Config.TTL == 0 {
		n.Config.TTL = 24 * time.Hour
	}
	if n.Config.BatchSize == 0 {
		n.Config.BatchSize = 1000
	}
	if n.Config.RateLimit == 0 {
		n.Config.RateLimit = 100
	}
	return nil
}

func (n *Notification) IsExpired() bool {
	return n.ExpiresAt != nil && time.Now().After(*n.ExpiresAt)
}

func (n *Notification) CanRetry() bool {
	return n.RetryCount < n.Config.MaxRetries && !n.IsExpired()
}

func (n *Notification) IsScheduled() bool {
	return n.ScheduledAt != nil && time.Now().Before(*n.ScheduledAt)
}

func (n *Notification) GetLocalizedContent(language string) *LocalizedContent {
	for _, content := range n.LocalizedContent {
		if content.Language == language {
			return &content
		}
	}
	// Fallback to English or first available
	for _, content := range n.LocalizedContent {
		if content.Language == "en" {
			return &content
		}
	}
	if len(n.LocalizedContent) > 0 {
		return &n.LocalizedContent[0]
	}
	return nil
}

func (n *Notification) UpdateAnalytics() {
	if n.TotalUsers > 0 {
		n.Analytics.DeliveryRate = float64(n.Analytics.Delivered) / float64(n.TotalUsers) * 100
		n.Analytics.FailureRate = float64(n.Analytics.Failed) / float64(n.TotalUsers) * 100
	}
	if n.Analytics.Delivered > 0 {
		n.Analytics.ReadRate = float64(n.Analytics.Read) / float64(n.Analytics.Delivered) * 100
		n.Analytics.ClickRate = float64(n.Analytics.Clicked) / float64(n.Analytics.Delivered) * 100
	}
}

func (n *Notification) AddFailure(reason FailureReason, message string) {
	n.Status = DeliveryStatusFailed
	n.FailureReason = reason
	n.ErrorMessage = message
	n.RetryCount++
	n.Analytics.Failed++
	n.UpdateAnalytics()
}

func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.Status = DeliveryStatusSent
	n.SentAt = &now
	n.Analytics.Sent++
	n.UpdateAnalytics()
}

func (n *Notification) MarkAsDelivered() {
	n.Status = DeliveryStatusDelivered
	n.Analytics.Delivered++
	n.UpdateAnalytics()
}

func (un *UserNotification) BeforeCreate(tx *gorm.DB) error {
	if un.ID == uuid.Nil {
		un.ID = uuid.New()
	}
	return nil
}

func (un *UserNotification) MarkAsRead() {
	if !un.IsRead {
		un.IsRead = true
		now := time.Now()
		un.ReadAt = &now
	}
}

func (un *UserNotification) MarkAsClicked() {
	if !un.IsClicked {
		un.IsClicked = true
		now := time.Now()
		un.ClickedAt = &now
	}
	if !un.IsRead {
		un.MarkAsRead()
	}
}

func (un *UserNotification) MarkAsSent() {
	un.Status = DeliveryStatusSent
	now := time.Now()
	un.SentAt = &now
}

func (un *UserNotification) MarkAsDelivered() {
	un.Status = DeliveryStatusDelivered
	now := time.Now()
	un.DeliveredAt = &now
}

func (un *UserNotification) MarkAsFailed(reason FailureReason, message string) {
	un.Status = DeliveryStatusFailed
	un.FailureReason = reason
	un.ErrorMessage = message
	un.RetryCount++
}

func (nt *NotificationTemplate) BeforeCreate(tx *gorm.DB) error {
	if nt.ID == uuid.Nil {
		nt.ID = uuid.New()
	}
	return nil
}

func (nt *NotificationTemplate) GetVariables() []string {
	if nt.Variables != nil {
		return nt.Variables
	}

	// Extract variables from body template
	variables := make([]string, 0)
	body := nt.Body

	for {
		start := strings.Index(body, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(body[start:], "}}")
		if end == -1 {
			break
		}
		variable := strings.TrimSpace(body[start+2 : start+end])
		variables = append(variables, variable)
		body = body[start+end+2:]
	}

	return variables
}

func (nt *NotificationTemplate) RenderContent(variables map[string]string, language string) (*LocalizedContent, error) {
	content := nt.GetLocalizedVersion(language)
	if content == nil {
		return nil, fmt.Errorf("no localized version found for language: %s", language)
	}

	title := content.Title
	body := content.Body

	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		title = strings.ReplaceAll(title, placeholder, value)
		body = strings.ReplaceAll(body, placeholder, value)
	}

	return &LocalizedContent{
		Language:  language,
		Title:     title,
		Body:      body,
		ImageURL:  content.ImageURL,
		ActionURL: content.ActionURL,
	}, nil
}

func (nt *NotificationTemplate) GetLocalizedVersion(language string) *LocalizedContent {
	for _, version := range nt.LocalizedVersions {
		if version.Language == language {
			return &version
		}
	}
	// Fallback to English
	for _, version := range nt.LocalizedVersions {
		if version.Language == "en" {
			return &version
		}
	}
	// Fallback to first available
	if len(nt.LocalizedVersions) > 0 {
		return &nt.LocalizedVersions[0]
	}
	return nil
}

func (nt *NotificationTemplate) AddLocalizedVersion(content LocalizedContent) {
	// Remove existing version for the same language
	for i, existing := range nt.LocalizedVersions {
		if existing.Language == content.Language {
			nt.LocalizedVersions = append(nt.LocalizedVersions[:i], nt.LocalizedVersions[i+1:]...)
			break
		}
	}
	nt.LocalizedVersions = append(nt.LocalizedVersions, content)
}

// Response structures for API
type NotificationResponse struct {
	ID              uuid.UUID            `json:"id"`
	Title           string               `json:"title"`
	Body            string               `json:"body"`
	ImageURL        string               `json:"image_url,omitempty"`
	ActionURL       string               `json:"action_url,omitempty"`
	Type            NotificationType     `json:"type"`
	Category        NotificationCategory `json:"category"`
	Priority        Priority             `json:"priority"`
	Status          DeliveryStatus       `json:"status"`
	TotalUsers      int64                `json:"total_users"`
	ScheduledAt     *time.Time           `json:"scheduled_at,omitempty"`
	SentAt          *time.Time           `json:"sent_at,omitempty"`
	ExpiresAt       *time.Time           `json:"expires_at,omitempty"`
	Analytics       Analytics            `json:"analytics"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type UserNotificationResponse struct {
	ID             uuid.UUID      `json:"id"`
	NotificationID uuid.UUID      `json:"notification_id"`
	Title          string         `json:"title"`
	Body           string         `json:"body"`
	ImageURL       string         `json:"image_url,omitempty"`
	ActionURL      string         `json:"action_url,omitempty"`
	Category       NotificationCategory `json:"category"`
	Priority       Priority       `json:"priority"`
	Status         DeliveryStatus `json:"status"`
	IsRead         bool           `json:"is_read"`
	ReadAt         *time.Time     `json:"read_at,omitempty"`
	IsClicked      bool           `json:"is_clicked"`
	ClickedAt      *time.Time     `json:"clicked_at,omitempty"`
	SentAt         *time.Time     `json:"sent_at,omitempty"`
	DeliveredAt    *time.Time     `json:"delivered_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

// Request structures for API
type CreateNotificationRequest struct {
	Title            string               `json:"title" binding:"required"`
	Body             string               `json:"body" binding:"required"`
	ImageURL         string               `json:"image_url,omitempty"`
	ActionURL        string               `json:"action_url,omitempty"`
	Type             NotificationType     `json:"type" binding:"required"`
	Category         NotificationCategory `json:"category" binding:"required"`
	Priority         Priority             `json:"priority"`
	Targeting        Targeting            `json:"targeting" binding:"required"`
	ScheduledAt      *time.Time           `json:"scheduled_at,omitempty"`
	ExpiresAt        *time.Time           `json:"expires_at,omitempty"`
	LocalizedContent []LocalizedContent   `json:"localized_content,omitempty"`
	Metadata         map[string]any       `json:"metadata,omitempty"`
	Config           *DeliveryConfig      `json:"config,omitempty"`
}

type NotificationTemplateRequest struct {
	Name              string             `json:"name" binding:"required"`
	Type              NotificationType   `json:"type" binding:"required"`
	Category          NotificationCategory `json:"category" binding:"required"`
	Subject           string             `json:"subject"`
	Body              string             `json:"body" binding:"required"`
	LocalizedVersions []LocalizedContent `json:"localized_versions,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
}

// Manager for business logic
type NotificationManager struct{}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{}
}

func (nm *NotificationManager) ValidateNotification(notification *Notification) error {
	if strings.TrimSpace(notification.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if strings.TrimSpace(notification.Body) == "" {
		return fmt.Errorf("body is required")
	}

	if notification.Type == "" {
		return fmt.Errorf("notification type is required")
	}

	if notification.Category == "" {
		return fmt.Errorf("notification category is required")
	}

	if notification.Targeting.AudienceType == "" {
		return fmt.Errorf("audience type is required")
	}

	if notification.ScheduledAt != nil && notification.ScheduledAt.Before(time.Now()) {
		return fmt.Errorf("scheduled time must be in the future")
	}

	if notification.ExpiresAt != nil && notification.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("expiration time must be in the future")
	}

	return nil
}

func (nm *NotificationManager) CalculateTargetUsers(targeting Targeting) ([]uuid.UUID, error) {
	userIDs := make([]uuid.UUID, 0)

	switch targeting.AudienceType {
	case AudienceTypeUser:
		userIDs = append(userIDs, targeting.UserIDs...)
	case AudienceTypeSegment:
		// Implementation would query users based on segments
	case AudienceTypeRole:
		// Implementation would query users based on roles
	case AudienceTypeBroadcast:
		// Implementation would get all active users
	case AudienceTypeLocation:
		// Implementation would query users based on location
	case AudienceTypeDevice:
		// Implementation would query users based on device type
	}

	return userIDs, nil
}

func (nm *NotificationManager) CreateUserNotifications(notification *Notification, userIDs []uuid.UUID) ([]*UserNotification, error) {
	userNotifications := make([]*UserNotification, len(userIDs))

	for i, userID := range userIDs {
		userNotifications[i] = &UserNotification{
			ID:             uuid.New(),
			NotificationID: notification.ID,
			UserID:         userID,
			Status:         DeliveryStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	return userNotifications, nil
}

func (nm *NotificationManager) GetSupportedLanguages() []string {
	return []string{"en", "th", "id", "ms", "vi", "zh", "tl"}
}

func (nm *NotificationManager) GetDefaultConfig() DeliveryConfig {
	return DeliveryConfig{
		MaxRetries:    3,
		RetryInterval: 5 * time.Minute,
		TTL:           24 * time.Hour,
		BatchSize:     1000,
		RateLimit:     100,
	}
}

// IsValid methods for validation
func (c NotificationChannel) IsValid() bool {
	switch c {
	case NotificationChannelPush, NotificationChannelEmail, NotificationChannelSMS, NotificationChannelInApp, NotificationChannelWebPush, NotificationChannelWebhook:
		return true
	default:
		return false
	}
}

func (p NotificationPriority) IsValid() bool {
	switch p {
	case NotificationPriorityLow, NotificationPriorityNormal, NotificationPriorityHigh, NotificationPriorityCritical:
		return true
	default:
		return false
	}
}