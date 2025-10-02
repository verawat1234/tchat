package messaging

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// EventBus manages event publishing and subscription across services
type EventBus struct {
	kafka     *KafkaClient
	handlers  map[string][]EventHandler
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// EventHandler is a function that handles events
type EventHandler func(ctx context.Context, event *Event) error

// NewEventBus creates a new event bus
func NewEventBus(kafka *KafkaClient) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())

	return &EventBus{
		kafka:    kafka,
		handlers: make(map[string][]EventHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Subscribe registers an event handler for a specific topic
func (eb *EventBus) Subscribe(topic string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[topic] = append(eb.handlers[topic], handler)
	log.Printf("Subscribed handler to topic: %s", topic)
}

// Publish publishes an event to a topic
func (eb *EventBus) Publish(ctx context.Context, topic string, event *Event) error {
	return eb.kafka.PublishEvent(ctx, topic, event)
}

// PublishBatch publishes multiple events to a topic
func (eb *EventBus) PublishBatch(ctx context.Context, topic string, events []*Event) error {
	return eb.kafka.PublishEventBatch(ctx, topic, events)
}

// Start starts consuming events from subscribed topics
func (eb *EventBus) Start() error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	for topic, handlers := range eb.handlers {
		if len(handlers) > 0 {
			eb.wg.Add(1)
			go eb.consumeTopic(topic, handlers)
		}
	}

	log.Println("Event bus started")
	return nil
}

// Stop stops the event bus and waits for all consumers to finish
func (eb *EventBus) Stop() error {
	eb.cancel()
	eb.wg.Wait()
	log.Println("Event bus stopped")
	return nil
}

// consumeTopic consumes events from a topic and calls registered handlers
func (eb *EventBus) consumeTopic(topic string, handlers []EventHandler) {
	defer eb.wg.Done()

	err := eb.kafka.ConsumeEvents(eb.ctx, topic, eb.kafka.config.GroupID, func(event *Event) error {
		for _, handler := range handlers {
			if err := handler(eb.ctx, event); err != nil {
				log.Printf("Error in event handler for topic %s: %v", topic, err)
				// Continue processing other handlers
			}
		}
		return nil
	})

	if err != nil && err != context.Canceled {
		log.Printf("Error consuming topic %s: %v", topic, err)
	}
}

// User event publishers

// PublishUserRegistered publishes a user registered event
func (eb *EventBus) PublishUserRegistered(ctx context.Context, userID, phone, countryCode string) error {
	event := NewUserEvent(UserRegisteredTopic, userID, map[string]interface{}{
		"user_id":      userID,
		"phone":        phone,
		"country_code": countryCode,
	})

	return eb.Publish(ctx, UserRegisteredTopic, event)
}

// PublishUserLogin publishes a user login event
func (eb *EventBus) PublishUserLogin(ctx context.Context, userID, sessionID, ipAddress, userAgent string) error {
	event := NewUserEvent(UserLoginTopic, userID, map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
		"ip_address": ipAddress,
		"user_agent": userAgent,
	})

	return eb.Publish(ctx, UserLoginTopic, event)
}

// PublishUserLogout publishes a user logout event
func (eb *EventBus) PublishUserLogout(ctx context.Context, userID, sessionID string) error {
	event := NewUserEvent(UserLogoutTopic, userID, map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	})

	return eb.Publish(ctx, UserLogoutTopic, event)
}

// Message event publishers

// PublishMessageSent publishes a message sent event
func (eb *EventBus) PublishMessageSent(ctx context.Context, messageID, dialogID, senderID, messageType, content string) error {
	event := NewMessageEvent(MessageSentTopic, messageID, map[string]interface{}{
		"message_id":   messageID,
		"dialog_id":    dialogID,
		"sender_id":    senderID,
		"message_type": messageType,
		"content":      content,
	})

	return eb.Publish(ctx, MessageSentTopic, event)
}

// PublishMessageEdited publishes a message edited event
func (eb *EventBus) PublishMessageEdited(ctx context.Context, messageID, dialogID, senderID, newContent string) error {
	event := NewMessageEvent(MessageEditedTopic, messageID, map[string]interface{}{
		"message_id":  messageID,
		"dialog_id":   dialogID,
		"sender_id":   senderID,
		"new_content": newContent,
	})

	return eb.Publish(ctx, MessageEditedTopic, event)
}

// PublishMessageDeleted publishes a message deleted event
func (eb *EventBus) PublishMessageDeleted(ctx context.Context, messageID, dialogID, senderID string) error {
	event := NewMessageEvent(MessageDeletedTopic, messageID, map[string]interface{}{
		"message_id": messageID,
		"dialog_id":  dialogID,
		"sender_id":  senderID,
	})

	return eb.Publish(ctx, MessageDeletedTopic, event)
}

// PublishMessageDelivered publishes a message delivered event
func (eb *EventBus) PublishMessageDelivered(ctx context.Context, messageID, userID string) error {
	event := NewMessageEvent(MessageDeliveredTopic, messageID, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
	})

	return eb.Publish(ctx, MessageDeliveredTopic, event)
}

// PublishMessageRead publishes a message read event
func (eb *EventBus) PublishMessageRead(ctx context.Context, messageID, userID string) error {
	event := NewMessageEvent(MessageReadTopic, messageID, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
	})

	return eb.Publish(ctx, MessageReadTopic, event)
}

// Dialog event publishers

// PublishDialogCreated publishes a dialog created event
func (eb *EventBus) PublishDialogCreated(ctx context.Context, dialogID, creatorID, dialogType, title string, participantIDs []string) error {
	event := NewDialogEvent(DialogCreatedTopic, dialogID, map[string]interface{}{
		"dialog_id":       dialogID,
		"creator_id":      creatorID,
		"dialog_type":     dialogType,
		"title":           title,
		"participant_ids": participantIDs,
	})

	return eb.Publish(ctx, DialogCreatedTopic, event)
}

// PublishParticipantJoined publishes a participant joined event
func (eb *EventBus) PublishParticipantJoined(ctx context.Context, dialogID, userID, inviterID, role string) error {
	event := NewDialogEvent(ParticipantJoinedTopic, dialogID, map[string]interface{}{
		"dialog_id":  dialogID,
		"user_id":    userID,
		"inviter_id": inviterID,
		"role":       role,
	})

	return eb.Publish(ctx, ParticipantJoinedTopic, event)
}

// PublishParticipantLeft publishes a participant left event
func (eb *EventBus) PublishParticipantLeft(ctx context.Context, dialogID, userID string) error {
	event := NewDialogEvent(ParticipantLeftTopic, dialogID, map[string]interface{}{
		"dialog_id": dialogID,
		"user_id":   userID,
	})

	return eb.Publish(ctx, ParticipantLeftTopic, event)
}

// Payment event publishers

// PublishPaymentInitiated publishes a payment initiated event
func (eb *EventBus) PublishPaymentInitiated(ctx context.Context, transactionID, fromUserID, toUserID, currency string, amount float64) error {
	event := NewPaymentEvent(PaymentInitiatedTopic, transactionID, map[string]interface{}{
		"transaction_id": transactionID,
		"from_user_id":   fromUserID,
		"to_user_id":     toUserID,
		"currency":       currency,
		"amount":         amount,
	})

	return eb.Publish(ctx, PaymentInitiatedTopic, event)
}

// PublishPaymentCompleted publishes a payment completed event
func (eb *EventBus) PublishPaymentCompleted(ctx context.Context, transactionID, fromUserID, toUserID, currency string, amount float64) error {
	event := NewPaymentEvent(PaymentCompletedTopic, transactionID, map[string]interface{}{
		"transaction_id": transactionID,
		"from_user_id":   fromUserID,
		"to_user_id":     toUserID,
		"currency":       currency,
		"amount":         amount,
	})

	return eb.Publish(ctx, PaymentCompletedTopic, event)
}

// PublishPaymentFailed publishes a payment failed event
func (eb *EventBus) PublishPaymentFailed(ctx context.Context, transactionID, fromUserID, reason string) error {
	event := NewPaymentEvent(PaymentFailedTopic, transactionID, map[string]interface{}{
		"transaction_id": transactionID,
		"from_user_id":   fromUserID,
		"reason":         reason,
	})

	return eb.Publish(ctx, PaymentFailedTopic, event)
}

// PublishWalletCreated publishes a wallet created event
func (eb *EventBus) PublishWalletCreated(ctx context.Context, walletID, userID, currency string) error {
	event := NewPaymentEvent(WalletCreatedTopic, walletID, map[string]interface{}{
		"wallet_id": walletID,
		"user_id":   userID,
		"currency":  currency,
	})

	return eb.Publish(ctx, WalletCreatedTopic, event)
}

// Real-time event publishers

// PublishPresenceUpdate publishes a presence update event
func (eb *EventBus) PublishPresenceUpdate(ctx context.Context, userID, status string, lastSeen time.Time) error {
	event := NewUserEvent(PresenceUpdateTopic, userID, map[string]interface{}{
		"user_id":   userID,
		"status":    status,
		"last_seen": lastSeen.Unix(),
	})

	return eb.Publish(ctx, PresenceUpdateTopic, event)
}

// PublishTypingIndicator publishes a typing indicator event
func (eb *EventBus) PublishTypingIndicator(ctx context.Context, dialogID, userID string, isTyping bool) error {
	event := NewDialogEvent(TypingIndicatorTopic, dialogID, map[string]interface{}{
		"dialog_id":  dialogID,
		"user_id":    userID,
		"is_typing":  isTyping,
		"timestamp": time.Now().Unix(),
	})

	return eb.Publish(ctx, TypingIndicatorTopic, event)
}

// PublishConnectionEvent publishes a connection event
func (eb *EventBus) PublishConnectionEvent(ctx context.Context, userID, connectionID, eventType string, metadata map[string]interface{}) error {
	event := NewUserEvent(ConnectionEventTopic, userID, map[string]interface{}{
		"user_id":       userID,
		"connection_id": connectionID,
		"event_type":    eventType, // connected, disconnected, heartbeat
		"metadata":      metadata,
	})

	return eb.Publish(ctx, ConnectionEventTopic, event)
}

// System event publishers

// PublishSystemAlert publishes a system alert event
func (eb *EventBus) PublishSystemAlert(ctx context.Context, alertType, message string, severity string, metadata map[string]interface{}) error {
	event := NewSystemEvent(SystemAlertTopic, alertType, map[string]interface{}{
		"alert_type": alertType,
		"message":    message,
		"severity":   severity,
		"metadata":   metadata,
	})

	return eb.Publish(ctx, SystemAlertTopic, event)
}

// PublishAuditLog publishes an audit log event
func (eb *EventBus) PublishAuditLog(ctx context.Context, userID, action, resource string, metadata map[string]interface{}) error {
	event := NewSystemEvent(AuditLogTopic, fmt.Sprintf("%s:%s", userID, action), map[string]interface{}{
		"user_id":  userID,
		"action":   action,
		"resource": resource,
		"metadata": metadata,
	})

	return eb.Publish(ctx, AuditLogTopic, event)
}

// PublishMetricsUpdate publishes a metrics update event
func (eb *EventBus) PublishMetricsUpdate(ctx context.Context, metricName string, value float64, tags map[string]string) error {
	event := NewSystemEvent(MetricsUpdateTopic, metricName, map[string]interface{}{
		"metric_name": metricName,
		"value":       value,
		"tags":        tags,
	})

	return eb.Publish(ctx, MetricsUpdateTopic, event)
}

// Health check
func (eb *EventBus) HealthCheck(ctx context.Context) error {
	return eb.kafka.HealthCheck(ctx)
}

// CreateAllTopics creates all application topics
func (eb *EventBus) CreateAllTopics(ctx context.Context) error {
	topics := GetAllTopics()
	return eb.kafka.CreateTopics(ctx, topics)
}