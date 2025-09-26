package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeSystem     NotificationType = "system"
	NotificationTypeMessage    NotificationType = "message"
	NotificationTypeOrder      NotificationType = "order"
	NotificationTypePayment    NotificationType = "payment"
	NotificationTypePromotion  NotificationType = "promotion"
	NotificationTypeFriend     NotificationType = "friend"
	NotificationTypeReview     NotificationType = "review"
	NotificationTypeSecurity   NotificationType = "security"
	NotificationTypeMarketing  NotificationType = "marketing"
	NotificationTypeReminder   NotificationType = "reminder"
)

// IsValid checks if the notification type is valid
func (nt NotificationType) IsValid() bool {
	switch nt {
	case NotificationTypeSystem, NotificationTypeMessage, NotificationTypeOrder,
		 NotificationTypePayment, NotificationTypePromotion, NotificationTypeFriend,
		 NotificationTypeReview, NotificationTypeSecurity, NotificationTypeMarketing,
		 NotificationTypeReminder:
		return true
	default:
		return false
	}
}

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusRead      NotificationStatus = "read"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusCancelled NotificationStatus = "cancelled"
)

// IsValid checks if the notification status is valid
func (ns NotificationStatus) IsValid() bool {
	switch ns {
	case NotificationStatusPending, NotificationStatusSent, NotificationStatusDelivered,
		 NotificationStatusRead, NotificationStatusFailed, NotificationStatusCancelled:
		return true
	default:
		return false
	}
}

// IsRead checks if the notification has been read
func (ns NotificationStatus) IsRead() bool {
	return ns == NotificationStatusRead
}

// NotificationPriority represents the priority level of a notification
type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityNormal   NotificationPriority = "normal"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

// IsValid checks if the notification priority is valid
func (np NotificationPriority) IsValid() bool {
	switch np {
	case NotificationPriorityLow, NotificationPriorityNormal,
		 NotificationPriorityHigh, NotificationPriorityCritical:
		return true
	default:
		return false
	}
}

// NotificationChannel represents the delivery channel for notifications
type NotificationChannel string

const (
	NotificationChannelInApp  NotificationChannel = "in_app"
	NotificationChannelPush   NotificationChannel = "push"
	NotificationChannelEmail  NotificationChannel = "email"
	NotificationChannelSMS    NotificationChannel = "sms"
	NotificationChannelWebhook NotificationChannel = "webhook"
)

// IsValid checks if the notification channel is valid
func (nc NotificationChannel) IsValid() bool {
	switch nc {
	case NotificationChannelInApp, NotificationChannelPush, NotificationChannelEmail,
		 NotificationChannelSMS, NotificationChannelWebhook:
		return true
	default:
		return false
	}
}

// NotificationAction represents an action that can be taken from a notification
type NotificationAction struct {
	ID       string                 `json:"id" gorm:"column:id;size:50;not null"`
	Label    string                 `json:"label" gorm:"column:label;size:100;not null"`
	URL      string                 `json:"url,omitempty" gorm:"column:url;size:500"`
	Action   string                 `json:"action" gorm:"column:action;size:50"`
	Style    string                 `json:"style,omitempty" gorm:"column:style;size:20"`
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

// NotificationContent represents the content of a notification in multiple languages
type NotificationContent struct {
	Title       string                 `json:"title" gorm:"column:title;size:200;not null"`
	Body        string                 `json:"body" gorm:"column:body;size:1000;not null"`
	Summary     string                 `json:"summary,omitempty" gorm:"column:summary;size:300"`
	Language    string                 `json:"language" gorm:"column:language;size:5;not null;default:'en'"`
	RTL         bool                   `json:"rtl" gorm:"column:rtl;default:false"`
	Personalized map[string]interface{} `json:"personalized,omitempty" gorm:"column:personalized;type:jsonb"`
}

// NotificationMedia represents media attachments for notifications
type NotificationMedia struct {
	Type        string `json:"type" gorm:"column:type;size:20;not null"`
	URL         string `json:"url" gorm:"column:url;size:500;not null"`
	ThumbnailURL string `json:"thumbnail_url,omitempty" gorm:"column:thumbnail_url;size:500"`
	AltText     string `json:"alt_text,omitempty" gorm:"column:alt_text;size:200"`
	Width       int    `json:"width,omitempty" gorm:"column:width"`
	Height      int    `json:"height,omitempty" gorm:"column:height"`
	Size        int64  `json:"size,omitempty" gorm:"column:size"`
}

// NotificationDelivery represents delivery information for a specific channel
type NotificationDelivery struct {
	Channel       NotificationChannel `json:"channel" gorm:"column:channel;type:varchar(20);not null"`
	Status        NotificationStatus  `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`
	AttemptCount  int                 `json:"attempt_count" gorm:"column:attempt_count;default:0"`
	MaxAttempts   int                 `json:"max_attempts" gorm:"column:max_attempts;default:3"`
	SentAt        *time.Time          `json:"sent_at,omitempty" gorm:"column:sent_at"`
	DeliveredAt   *time.Time          `json:"delivered_at,omitempty" gorm:"column:delivered_at"`
	FailedAt      *time.Time          `json:"failed_at,omitempty" gorm:"column:failed_at"`
	ErrorMessage  string              `json:"error_message,omitempty" gorm:"column:error_message;size:500"`
	ExternalID    string              `json:"external_id,omitempty" gorm:"column:external_id;size:255"`
	ProviderData  map[string]interface{} `json:"provider_data,omitempty" gorm:"column:provider_data;type:jsonb"`
}

// NotificationScheduling represents scheduling information
type NotificationScheduling struct {
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty" gorm:"column:scheduled_at"`
	Timezone        string     `json:"timezone,omitempty" gorm:"column:timezone;size:50"`
	RecurringType   string     `json:"recurring_type,omitempty" gorm:"column:recurring_type;size:20"`
	RecurringConfig map[string]interface{} `json:"recurring_config,omitempty" gorm:"column:recurring_config;type:jsonb"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"`
}

// NotificationTargeting represents targeting and audience information
type NotificationTargeting struct {
	AudienceSegments []string               `json:"audience_segments,omitempty" gorm:"column:audience_segments;type:jsonb"`
	UserCriteria     map[string]interface{} `json:"user_criteria,omitempty" gorm:"column:user_criteria;type:jsonb"`
	GeoCriteria      map[string]interface{} `json:"geo_criteria,omitempty" gorm:"column:geo_criteria;type:jsonb"`
	DeviceCriteria   map[string]interface{} `json:"device_criteria,omitempty" gorm:"column:device_criteria;type:jsonb"`
	Localization     map[string]NotificationContent `json:"localization,omitempty" gorm:"column:localization;type:jsonb"`
}

// Notification represents a notification in the system
type Notification struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RecipientID uuid.UUID `json:"recipient_id" gorm:"type:uuid;not null;index"`
	SenderID   *uuid.UUID `json:"sender_id,omitempty" gorm:"type:uuid;index"`

	// Notification properties
	Type     NotificationType     `json:"type" gorm:"column:type;type:varchar(20);not null"`
	Status   NotificationStatus   `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`
	Priority NotificationPriority `json:"priority" gorm:"column:priority;type:varchar(20);not null;default:'normal'"`

	// Content
	Content NotificationContent `json:"content" gorm:"embedded;embeddedPrefix:content_"`

	// Media attachments
	Media []NotificationMedia `json:"media,omitempty" gorm:"column:media;type:jsonb"`

	// Actions
	Actions []NotificationAction `json:"actions,omitempty" gorm:"column:actions;type:jsonb"`

	// Related entities
	RelatedEntityType string     `json:"related_entity_type,omitempty" gorm:"column:related_entity_type;size:50"`
	RelatedEntityID   *uuid.UUID `json:"related_entity_id,omitempty" gorm:"column:related_entity_id;type:uuid"`

	// Delivery information
	Channels   []NotificationChannel  `json:"channels" gorm:"column:channels;type:jsonb"`
	Deliveries []NotificationDelivery `json:"deliveries,omitempty" gorm:"column:deliveries;type:jsonb"`

	// Scheduling
	Scheduling NotificationScheduling `json:"scheduling" gorm:"embedded;embeddedPrefix:scheduling_"`

	// Targeting (for broadcast notifications)
	Targeting NotificationTargeting `json:"targeting" gorm:"embedded;embeddedPrefix:targeting_"`

	// User interaction
	ReadAt       *time.Time `json:"read_at,omitempty" gorm:"column:read_at"`
	ClickedAt    *time.Time `json:"clicked_at,omitempty" gorm:"column:clicked_at"`
	DismissedAt  *time.Time `json:"dismissed_at,omitempty" gorm:"column:dismissed_at"`
	ActionTaken  string     `json:"action_taken,omitempty" gorm:"column:action_taken;size:50"`
	InteractionData map[string]interface{} `json:"interaction_data,omitempty" gorm:"column:interaction_data;type:jsonb"`

	// Campaign and tracking
	CampaignID   *uuid.UUID `json:"campaign_id,omitempty" gorm:"column:campaign_id;type:uuid"`
	TrackingID   string     `json:"tracking_id,omitempty" gorm:"column:tracking_id;size:100"`
	Source       string     `json:"source" gorm:"column:source;size:50;default:'system'"`
	SourceData   map[string]interface{} `json:"source_data,omitempty" gorm:"column:source_data;type:jsonb"`

	// Regional compliance
	DataRegion     string `json:"data_region" gorm:"column:data_region;size:20"`
	ComplianceData map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`

	// Metadata and tags
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Recipient *User `json:"recipient,omitempty" gorm:"foreignKey:RecipientID;references:ID"`
	Sender    *User `json:"sender,omitempty" gorm:"foreignKey:SenderID;references:ID"`
}

// TableName returns the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate sets up the notification before creation
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	// Set data region based on recipient
	if n.DataRegion == "" {
		var recipient User
		if err := tx.First(&recipient, n.RecipientID).Error; err == nil {
			n.DataRegion = GetDataRegionForCountry(recipient.CountryCode)
		} else {
			n.DataRegion = "sea-central" // Default region
		}
	}

	// Set default language if not provided
	if n.Content.Language == "" {
		var recipient User
		if err := tx.First(&recipient, n.RecipientID).Error; err == nil {
			n.Content.Language = recipient.Locale
		} else {
			n.Content.Language = "en" // Default language
		}
	}

	// Initialize deliveries for specified channels
	if len(n.Deliveries) == 0 && len(n.Channels) > 0 {
		for _, channel := range n.Channels {
			delivery := NotificationDelivery{
				Channel:     channel,
				Status:      NotificationStatusPending,
				MaxAttempts: 3,
			}
			n.Deliveries = append(n.Deliveries, delivery)
		}
	}

	// Generate tracking ID if not provided
	if n.TrackingID == "" {
		n.TrackingID = fmt.Sprintf("ntf_%s_%d", n.ID.String()[:8], time.Now().Unix())
	}

	// Validate the notification
	if err := n.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the notification before updating
func (n *Notification) BeforeUpdate(tx *gorm.DB) error {
	return n.Validate()
}

// Validate validates the notification data
func (n *Notification) Validate() error {
	// Validate UUIDs
	if n.ID == uuid.Nil {
		return fmt.Errorf("notification ID cannot be nil")
	}
	if n.RecipientID == uuid.Nil {
		return fmt.Errorf("recipient ID cannot be nil")
	}

	// Validate type, status, and priority
	if !n.Type.IsValid() {
		return fmt.Errorf("invalid notification type: %s", n.Type)
	}
	if !n.Status.IsValid() {
		return fmt.Errorf("invalid notification status: %s", n.Status)
	}
	if !n.Priority.IsValid() {
		return fmt.Errorf("invalid notification priority: %s", n.Priority)
	}

	// Validate content
	if n.Content.Title == "" {
		return fmt.Errorf("notification title is required")
	}
	if n.Content.Body == "" {
		return fmt.Errorf("notification body is required")
	}

	// Validate channels
	for i, channel := range n.Channels {
		if !channel.IsValid() {
			return fmt.Errorf("invalid notification channel at index %d: %s", i, channel)
		}
	}

	// Validate deliveries
	for i, delivery := range n.Deliveries {
		if !delivery.Channel.IsValid() {
			return fmt.Errorf("invalid delivery channel at index %d: %s", i, delivery.Channel)
		}
		if !delivery.Status.IsValid() {
			return fmt.Errorf("invalid delivery status at index %d: %s", i, delivery.Status)
		}
	}

	// Validate actions
	for i, action := range n.Actions {
		if action.ID == "" {
			return fmt.Errorf("action ID is required at index %d", i)
		}
		if action.Label == "" {
			return fmt.Errorf("action label is required at index %d", i)
		}
	}

	// Validate media
	for i, media := range n.Media {
		if media.Type == "" {
			return fmt.Errorf("media type is required at index %d", i)
		}
		if media.URL == "" {
			return fmt.Errorf("media URL is required at index %d", i)
		}
	}

	return nil
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	if n.Status != NotificationStatusRead {
		n.Status = NotificationStatusRead
		now := time.Now()
		n.ReadAt = &now
		n.UpdatedAt = now
	}
}

// MarkAsClicked marks the notification as clicked
func (n *Notification) MarkAsClicked(actionID string) {
	now := time.Now()
	n.ClickedAt = &now
	n.ActionTaken = actionID
	n.UpdatedAt = now

	// Also mark as read if not already
	if n.ReadAt == nil {
		n.MarkAsRead()
	}
}

// MarkAsDismissed marks the notification as dismissed
func (n *Notification) MarkAsDismissed() {
	now := time.Now()
	n.DismissedAt = &now
	n.UpdatedAt = now
}

// UpdateDeliveryStatus updates the delivery status for a specific channel
func (n *Notification) UpdateDeliveryStatus(channel NotificationChannel, status NotificationStatus, errorMessage string) error {
	for i, delivery := range n.Deliveries {
		if delivery.Channel == channel {
			n.Deliveries[i].Status = status
			now := time.Now()

			switch status {
			case NotificationStatusSent:
				n.Deliveries[i].SentAt = &now
			case NotificationStatusDelivered:
				n.Deliveries[i].DeliveredAt = &now
				// Update overall status if all deliveries are successful
				if n.areAllDeliveriesSuccessful() {
					n.Status = NotificationStatusDelivered
				}
			case NotificationStatusFailed:
				n.Deliveries[i].FailedAt = &now
				n.Deliveries[i].ErrorMessage = errorMessage
				n.Deliveries[i].AttemptCount++
				// Update overall status if max attempts reached
				if n.Deliveries[i].AttemptCount >= n.Deliveries[i].MaxAttempts {
					n.Status = NotificationStatusFailed
				}
			}

			n.UpdatedAt = now
			return nil
		}
	}

	return fmt.Errorf("delivery for channel %s not found", channel)
}

// areAllDeliveriesSuccessful checks if all deliveries are successful
func (n *Notification) areAllDeliveriesSuccessful() bool {
	for _, delivery := range n.Deliveries {
		if delivery.Status != NotificationStatusDelivered {
			return false
		}
	}
	return len(n.Deliveries) > 0
}

// IsRead checks if the notification has been read
func (n *Notification) IsRead() bool {
	return n.Status.IsRead()
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	return n.Scheduling.ExpiresAt != nil && n.Scheduling.ExpiresAt.Before(time.Now())
}

// IsScheduled checks if the notification is scheduled for future delivery
func (n *Notification) IsScheduled() bool {
	return n.Scheduling.ScheduledAt != nil && n.Scheduling.ScheduledAt.After(time.Now())
}

// CanRetry checks if the notification can be retried
func (n *Notification) CanRetry() bool {
	if n.Status != NotificationStatusFailed {
		return false
	}

	for _, delivery := range n.Deliveries {
		if delivery.Status == NotificationStatusFailed && delivery.AttemptCount < delivery.MaxAttempts {
			return true
		}
	}
	return false
}

// GetContentForLanguage returns content for a specific language
func (n *Notification) GetContentForLanguage(language string) NotificationContent {
	// Check if localized content exists
	if content, exists := n.Targeting.Localization[language]; exists {
		return content
	}

	// Return default content
	return n.Content
}

// GetFailedDeliveries returns deliveries that have failed
func (n *Notification) GetFailedDeliveries() []NotificationDelivery {
	var failed []NotificationDelivery
	for _, delivery := range n.Deliveries {
		if delivery.Status == NotificationStatusFailed {
			failed = append(failed, delivery)
		}
	}
	return failed
}

// GetSuccessfulDeliveries returns deliveries that were successful
func (n *Notification) GetSuccessfulDeliveries() []NotificationDelivery {
	var successful []NotificationDelivery
	for _, delivery := range n.Deliveries {
		if delivery.Status == NotificationStatusDelivered {
			successful = append(successful, delivery)
		}
	}
	return successful
}

// GenerateSearchKeywords generates search keywords for the notification
func (n *Notification) GenerateSearchKeywords() []string {
	keywords := []string{
		n.Content.Title,
		n.Content.Body,
		n.Content.Summary,
		string(n.Type),
		string(n.Status),
		string(n.Priority),
		n.Source,
		n.TrackingID,
		n.RelatedEntityType,
	}

	// Add action labels
	for _, action := range n.Actions {
		keywords = append(keywords, action.Label)
	}

	// Add tags
	keywords = append(keywords, n.Tags...)

	// Add channel names
	for _, channel := range n.Channels {
		keywords = append(keywords, string(channel))
	}

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		cleaned := strings.ToLower(strings.TrimSpace(keyword))
		if cleaned != "" && len(cleaned) > 2 && !seen[cleaned] {
			seen[cleaned] = true
			unique = append(unique, cleaned)
		}
	}

	return unique
}

// GetNotificationSummary returns a summary of notification information
func (n *Notification) GetNotificationSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":                 n.ID,
		"type":               n.Type,
		"status":             n.Status,
		"priority":           n.Priority,
		"title":              n.Content.Title,
		"is_read":            n.IsRead(),
		"is_expired":         n.IsExpired(),
		"is_scheduled":       n.IsScheduled(),
		"can_retry":          n.CanRetry(),
		"channels":           n.Channels,
		"related_entity_type": n.RelatedEntityType,
		"related_entity_id":  n.RelatedEntityID,
		"source":             n.Source,
		"created_at":         n.CreatedAt,
		"read_at":            n.ReadAt,
		"clicked_at":         n.ClickedAt,
		"dismissed_at":       n.DismissedAt,
	}

	// Add delivery summary
	deliverySummary := make(map[string]interface{})
	for _, delivery := range n.Deliveries {
		deliverySummary[string(delivery.Channel)] = map[string]interface{}{
			"status":        delivery.Status,
			"attempt_count": delivery.AttemptCount,
			"sent_at":       delivery.SentAt,
			"delivered_at":  delivery.DeliveredAt,
			"failed_at":     delivery.FailedAt,
		}
	}
	summary["delivery_status"] = deliverySummary

	return summary
}

// MarshalJSON customizes JSON serialization
func (n *Notification) MarshalJSON() ([]byte, error) {
	type Alias Notification
	return json.Marshal(&struct {
		*Alias
		IsRead               bool                     `json:"is_read"`
		IsExpired            bool                     `json:"is_expired"`
		IsScheduled          bool                     `json:"is_scheduled"`
		CanRetry             bool                     `json:"can_retry"`
		FailedDeliveries     []NotificationDelivery   `json:"failed_deliveries,omitempty"`
		SuccessfulDeliveries []NotificationDelivery   `json:"successful_deliveries,omitempty"`
		NotificationSummary  map[string]interface{}   `json:"notification_summary"`
		SearchKeywords       []string                 `json:"search_keywords,omitempty"`
	}{
		Alias:               (*Alias)(n),
		IsRead:              n.IsRead(),
		IsExpired:           n.IsExpired(),
		IsScheduled:         n.IsScheduled(),
		CanRetry:            n.CanRetry(),
		FailedDeliveries:    n.GetFailedDeliveries(),
		SuccessfulDeliveries: n.GetSuccessfulDeliveries(),
		NotificationSummary: n.GetNotificationSummary(),
		SearchKeywords:      n.GenerateSearchKeywords(),
	})
}

// Helper functions for notification management

// CreateSystemNotification creates a system notification
func CreateSystemNotification(recipientID uuid.UUID, title, body string, priority NotificationPriority) *Notification {
	return &Notification{
		RecipientID: recipientID,
		Type:        NotificationTypeSystem,
		Priority:    priority,
		Content: NotificationContent{
			Title: title,
			Body:  body,
		},
		Channels: []NotificationChannel{NotificationChannelInApp},
		Source:   "system",
	}
}

// CreateMessageNotification creates a message notification
func CreateMessageNotification(recipientID, senderID uuid.UUID, messageContent string, dialogID uuid.UUID) *Notification {
	return &Notification{
		RecipientID:       recipientID,
		SenderID:          &senderID,
		Type:              NotificationTypeMessage,
		Priority:          NotificationPriorityNormal,
		RelatedEntityType: "message",
		RelatedEntityID:   &dialogID,
		Content: NotificationContent{
			Title: "New Message",
			Body:  messageContent,
		},
		Channels: []NotificationChannel{NotificationChannelInApp, NotificationChannelPush},
		Source:   "messaging",
	}
}

// CreateOrderNotification creates an order notification
func CreateOrderNotification(recipientID uuid.UUID, orderID uuid.UUID, title, body string) *Notification {
	return &Notification{
		RecipientID:       recipientID,
		Type:              NotificationTypeOrder,
		Priority:          NotificationPriorityHigh,
		RelatedEntityType: "order",
		RelatedEntityID:   &orderID,
		Content: NotificationContent{
			Title: title,
			Body:  body,
		},
		Channels: []NotificationChannel{NotificationChannelInApp, NotificationChannelEmail},
		Source:   "commerce",
	}
}

// CreatePromotionNotification creates a promotion notification
func CreatePromotionNotification(recipientID uuid.UUID, title, body string, imageURL string) *Notification {
	notification := &Notification{
		RecipientID: recipientID,
		Type:        NotificationTypePromotion,
		Priority:    NotificationPriorityLow,
		Content: NotificationContent{
			Title: title,
			Body:  body,
		},
		Channels: []NotificationChannel{NotificationChannelInApp, NotificationChannelPush},
		Source:   "marketing",
	}

	if imageURL != "" {
		notification.Media = []NotificationMedia{
			{
				Type: "image",
				URL:  imageURL,
			},
		}
	}

	return notification
}