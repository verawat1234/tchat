package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
)

// EventServiceImpl implements EventService
type EventServiceImpl struct {
	// In a real implementation, this could use Kafka, RabbitMQ, etc.
	// For now, we'll use a simple logging-based implementation
	enabledChannels map[string]bool
	asyncProcessing bool
}

// NewEventService creates a new event service
func NewEventService(asyncProcessing bool) EventService {
	return &EventServiceImpl{
		enabledChannels: map[string]bool{
			"notification.sent":      true,
			"notification.delivered": true,
			"notification.failed":    true,
			"notification.read":      true,
			"webhook.received":       true,
		},
		asyncProcessing: asyncProcessing,
	}
}

// Publish publishes a generic event
func (e *EventServiceImpl) Publish(ctx context.Context, event interface{}) error {
	if e.asyncProcessing {
		go e.publishAsync(event)
		return nil
	}

	return e.publishSync(event)
}

// PublishNotification publishes a notification event (alias for Publish)
func (e *EventServiceImpl) PublishNotification(ctx context.Context, event interface{}) error {
	return e.Publish(ctx, event)
}

// PublishNotificationSent publishes a notification sent event
func (e *EventServiceImpl) PublishNotificationSent(ctx context.Context, notification *models.Notification) error {
	if !e.enabledChannels["notification.sent"] {
		return nil
	}

	event := NotificationSentEvent{
		EventType:      "notification.sent",
		Timestamp:      time.Now(),
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Type:           string(notification.Type),
		Channel:        string(notification.Channel),
		Priority:       string(notification.Priority),
		Metadata: map[string]interface{}{
			"title":      notification.Title,
			"has_media":  notification.MediaURL != "",
			"scheduled":  notification.ScheduledAt != nil,
		},
	}

	return e.Publish(ctx, event)
}

// PublishNotificationDelivered publishes a notification delivered event
func (e *EventServiceImpl) PublishNotificationDelivered(ctx context.Context, notificationID uuid.UUID) error {
	if !e.enabledChannels["notification.delivered"] {
		return nil
	}

	event := NotificationDeliveredEvent{
		EventType:      "notification.delivered",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		DeliveredAt:    time.Now(),
	}

	return e.Publish(ctx, event)
}

// PublishNotificationFailed publishes a notification failed event
func (e *EventServiceImpl) PublishNotificationFailed(ctx context.Context, notificationID uuid.UUID, reason string) error {
	if !e.enabledChannels["notification.failed"] {
		return nil
	}

	event := NotificationFailedEvent{
		EventType:      "notification.failed",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		FailureReason:  reason,
		FailedAt:       time.Now(),
	}

	return e.Publish(ctx, event)
}

// PublishNotificationRead publishes a notification read event
func (e *EventServiceImpl) PublishNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	if !e.enabledChannels["notification.read"] {
		return nil
	}

	event := NotificationReadEvent{
		EventType:      "notification.read",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		UserID:         userID,
		ReadAt:         time.Now(),
	}

	return e.Publish(ctx, event)
}

// publishSync publishes an event synchronously
func (e *EventServiceImpl) publishSync(event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// In a real implementation, this would publish to a message queue
	// For now, we'll log the event
	log.Printf("Event published: %s", string(data))

	return nil
}

// publishAsync publishes an event asynchronously
func (e *EventServiceImpl) publishAsync(event interface{}) {
	if err := e.publishSync(event); err != nil {
		log.Printf("Failed to publish event asynchronously: %v", err)
	}
}

// EnableChannel enables event publishing for a specific channel
func (e *EventServiceImpl) EnableChannel(channel string) {
	e.enabledChannels[channel] = true
}

// DisableChannel disables event publishing for a specific channel
func (e *EventServiceImpl) DisableChannel(channel string) {
	e.enabledChannels[channel] = false
}

// Event types for different notification events

// NotificationSentEvent represents a notification sent event
type NotificationSentEvent struct {
	EventType      string                 `json:"event_type"`
	Timestamp      time.Time              `json:"timestamp"`
	NotificationID uuid.UUID              `json:"notification_id"`
	UserID         uuid.UUID              `json:"user_id"`
	Type           string                 `json:"type"`
	Channel        string                 `json:"channel"`
	Priority       string                 `json:"priority"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationDeliveredEvent represents a notification delivered event
type NotificationDeliveredEvent struct {
	EventType      string    `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	NotificationID uuid.UUID `json:"notification_id"`
	DeliveredAt    time.Time `json:"delivered_at"`
}

// NotificationFailedEvent represents a notification failed event
type NotificationFailedEvent struct {
	EventType      string    `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	NotificationID uuid.UUID `json:"notification_id"`
	FailureReason  string    `json:"failure_reason"`
	FailedAt       time.Time `json:"failed_at"`
}

// NotificationReadEvent represents a notification read event
type NotificationReadEvent struct {
	EventType      string    `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	NotificationID uuid.UUID `json:"notification_id"`
	UserID         uuid.UUID `json:"user_id"`
	ReadAt         time.Time `json:"read_at"`
}

// WebhookReceivedEvent represents a webhook received event
type WebhookReceivedEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Provider  string                 `json:"provider"`
	Event     string                 `json:"event"`
	MessageID string                 `json:"message_id"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// KafkaEventService implements EventService using Apache Kafka
// This is a placeholder for a real Kafka implementation
type KafkaEventService struct {
	// kafka.Producer
	topicPrefix string
}

// NewKafkaEventService creates a new Kafka-based event service
func NewKafkaEventService(brokers []string, topicPrefix string) (EventService, error) {
	// This would initialize a real Kafka producer
	// For now, return a placeholder
	return &KafkaEventService{
		topicPrefix: topicPrefix,
	}, nil
}

// Publish publishes an event to Kafka
func (k *KafkaEventService) Publish(ctx context.Context, event interface{}) error {
	// This would publish to Kafka topics
	// For now, fallback to logging
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event for Kafka: %w", err)
	}

	log.Printf("Kafka Event (would publish): %s", string(data))
	return nil
}

// PublishNotification publishes a notification event (alias for Publish)
func (k *KafkaEventService) PublishNotification(ctx context.Context, event interface{}) error {
	return k.Publish(ctx, event)
}

// PublishNotificationSent publishes notification sent event to Kafka
func (k *KafkaEventService) PublishNotificationSent(ctx context.Context, notification *models.Notification) error {
	event := NotificationSentEvent{
		EventType:      "notification.sent",
		Timestamp:      time.Now(),
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Type:           string(notification.Type),
		Channel:        string(notification.Channel),
		Priority:       string(notification.Priority),
	}

	return k.Publish(ctx, event)
}

// PublishNotificationDelivered publishes notification delivered event to Kafka
func (k *KafkaEventService) PublishNotificationDelivered(ctx context.Context, notificationID uuid.UUID) error {
	event := NotificationDeliveredEvent{
		EventType:      "notification.delivered",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		DeliveredAt:    time.Now(),
	}

	return k.Publish(ctx, event)
}

// PublishNotificationFailed publishes notification failed event to Kafka
func (k *KafkaEventService) PublishNotificationFailed(ctx context.Context, notificationID uuid.UUID, reason string) error {
	event := NotificationFailedEvent{
		EventType:      "notification.failed",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		FailureReason:  reason,
		FailedAt:       time.Now(),
	}

	return k.Publish(ctx, event)
}

// PublishNotificationRead publishes notification read event to Kafka
func (k *KafkaEventService) PublishNotificationRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	event := NotificationReadEvent{
		EventType:      "notification.read",
		Timestamp:      time.Now(),
		NotificationID: notificationID,
		UserID:         userID,
		ReadAt:         time.Now(),
	}

	return k.Publish(ctx, event)
}