package messaging

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers         []string      `mapstructure:"brokers" validate:"required"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	SASLMechanism   string        `mapstructure:"sasl_mechanism"` // plain, scram-sha-256, scram-sha-512
	SecurityProtocol string       `mapstructure:"security_protocol"` // plaintext, sasl_plaintext, sasl_ssl
	TopicPrefix     string        `mapstructure:"topic_prefix"`
	GroupID         string        `mapstructure:"group_id"`
	BatchSize       int           `mapstructure:"batch_size"`
	BatchTimeout    time.Duration `mapstructure:"batch_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	RequiredAcks    int           `mapstructure:"required_acks"`
	Compression     string        `mapstructure:"compression"` // none, gzip, snappy, lz4, zstd
	RetryMax        int           `mapstructure:"retry_max"`
	RetryBackoff    time.Duration `mapstructure:"retry_backoff"`
	EnableTLS       bool          `mapstructure:"enable_tls"`
	EnableMetrics   bool          `mapstructure:"enable_metrics"`
}

// DefaultKafkaConfig returns default Kafka configuration
func DefaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers:          []string{"localhost:9092"},
		SecurityProtocol: "plaintext",
		TopicPrefix:      "tchat",
		GroupID:          "tchat.dev",
		BatchSize:        100,
		BatchTimeout:     10 * time.Second,
		ReadTimeout:      30 * time.Second,
		WriteTimeout:     10 * time.Second,
		RequiredAcks:     1,
		Compression:      "snappy",
		RetryMax:         3,
		RetryBackoff:     1 * time.Second,
		EnableTLS:        false,
		EnableMetrics:    true,
	}
}

// KafkaClient wraps Kafka functionality
type KafkaClient struct {
	config   *KafkaConfig
	writers  map[string]*kafka.Writer
	readers  map[string]*kafka.Reader
	dialer   *kafka.Dialer
}

// NewKafkaClient creates a new Kafka client
func NewKafkaClient(config *KafkaConfig) (*KafkaClient, error) {
	if config == nil {
		config = DefaultKafkaConfig()
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure SASL authentication
	if config.Username != "" && config.Password != "" {
		switch config.SASLMechanism {
		case "plain":
			dialer.SASLMechanism = plain.Mechanism{
				Username: config.Username,
				Password: config.Password,
			}
		case "scram-sha-256":
			mechanism, err := scram.Mechanism(scram.SHA256, config.Username, config.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to create SCRAM mechanism: %w", err)
			}
			dialer.SASLMechanism = mechanism
		case "scram-sha-512":
			mechanism, err := scram.Mechanism(scram.SHA512, config.Username, config.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to create SCRAM mechanism: %w", err)
			}
			dialer.SASLMechanism = mechanism
		}
	}

	// Configure TLS
	if config.EnableTLS {
		dialer.TLS = &tls.Config{}
	}

	client := &KafkaClient{
		config:  config,
		writers: make(map[string]*kafka.Writer),
		readers: make(map[string]*kafka.Reader),
		dialer:  dialer,
	}

	log.Printf("Kafka client initialized with brokers: %v", config.Brokers)
	return client, nil
}

// Topic names
const (
	// User events
	UserRegisteredTopic    = "user.registered"
	UserUpdatedTopic      = "user.updated"
	UserDeletedTopic      = "user.deleted"
	UserLoginTopic        = "user.login"
	UserLogoutTopic       = "user.logout"

	// Message events
	MessageSentTopic      = "message.sent"
	MessageEditedTopic    = "message.edited"
	MessageDeletedTopic   = "message.deleted"
	MessageDeliveredTopic = "message.delivered"
	MessageReadTopic      = "message.read"

	// Dialog events
	DialogCreatedTopic      = "dialog.created"
	DialogUpdatedTopic      = "dialog.updated"
	DialogDeletedTopic      = "dialog.deleted"
	ParticipantJoinedTopic  = "dialog.participant.joined"
	ParticipantLeftTopic    = "dialog.participant.left"

	// Payment events
	PaymentInitiatedTopic = "payment.initiated"
	PaymentCompletedTopic = "payment.completed"
	PaymentFailedTopic    = "payment.failed"
	PaymentRefundedTopic  = "payment.refunded"
	WalletCreatedTopic    = "wallet.created"
	WalletUpdatedTopic    = "wallet.updated"

	// Notification events
	NotificationSentTopic      = "notification.sent"
	NotificationDeliveredTopic = "notification.delivered"
	NotificationReadTopic      = "notification.read"

	// System events
	SystemAlertTopic    = "system.alert"
	AuditLogTopic      = "audit.log"
	MetricsUpdateTopic = "metrics.update"

	// Real-time events
	PresenceUpdateTopic  = "presence.update"
	TypingIndicatorTopic = "typing.indicator"
	ConnectionEventTopic = "connection.event"
)

// Event represents a domain event
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
}

// getTopicName returns the full topic name with prefix
func (k *KafkaClient) getTopicName(topic string) string {
	if k.config.TopicPrefix != "" {
		return fmt.Sprintf("%s.%s", k.config.TopicPrefix, topic)
	}
	return topic
}

// getCompression returns the compression codec
func (k *KafkaClient) getCompression() kafka.Compression {
	switch k.config.Compression {
	case "gzip":
		return kafka.Gzip
	case "snappy":
		return kafka.Snappy
	case "lz4":
		return kafka.Lz4
	case "zstd":
		return kafka.Zstd
	default:
		return kafka.Snappy
	}
}

// GetWriter returns a Kafka writer for the specified topic
func (k *KafkaClient) GetWriter(topic string) *kafka.Writer {
	if writer, exists := k.writers[topic]; exists {
		return writer
	}

	fullTopic := k.getTopicName(topic)
	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.config.Brokers...),
		Topic:        fullTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    k.config.BatchSize,
		BatchTimeout: k.config.BatchTimeout,
		ReadTimeout:  k.config.ReadTimeout,
		WriteTimeout: k.config.WriteTimeout,
		RequiredAcks: kafka.RequiredAcks(k.config.RequiredAcks),
		Compression:  k.getCompression(),
		Transport: &kafka.Transport{
			Dial: k.dialer.DialFunc,
		},
	}

	k.writers[topic] = writer
	return writer
}

// GetReader returns a Kafka reader for the specified topic
func (k *KafkaClient) GetReader(topic, groupID string) *kafka.Reader {
	readerKey := fmt.Sprintf("%s:%s", topic, groupID)
	if reader, exists := k.readers[readerKey]; exists {
		return reader
	}

	fullTopic := k.getTopicName(topic)
	if groupID == "" {
		groupID = k.config.GroupID
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.config.Brokers,
		Topic:    fullTopic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Dialer:   k.dialer,
	})

	k.readers[readerKey] = reader
	return reader
}

// PublishEvent publishes an event to a topic
func (k *KafkaClient) PublishEvent(ctx context.Context, topic string, event *Event) error {
	writer := k.GetWriter(topic)

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(event.Subject),
		Value: eventData,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(event.Type)},
			{Key: "event-source", Value: []byte(event.Source)},
			{Key: "event-version", Value: []byte(event.Version)},
		},
		Time: event.Timestamp,
	}

	if err := writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	return nil
}

// PublishEventBatch publishes multiple events to a topic
func (k *KafkaClient) PublishEventBatch(ctx context.Context, topic string, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	writer := k.GetWriter(topic)
	messages := make([]kafka.Message, len(events))

	for i, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event %d: %w", i, err)
		}

		messages[i] = kafka.Message{
			Key:   []byte(event.Subject),
			Value: eventData,
			Headers: []kafka.Header{
				{Key: "event-type", Value: []byte(event.Type)},
				{Key: "event-source", Value: []byte(event.Source)},
				{Key: "event-version", Value: []byte(event.Version)},
			},
			Time: event.Timestamp,
		}
	}

	if err := writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("failed to write batch messages to topic %s: %w", topic, err)
	}

	return nil
}

// ConsumeEvents consumes events from a topic
func (k *KafkaClient) ConsumeEvents(ctx context.Context, topic, groupID string, handler func(*Event) error) error {
	reader := k.GetReader(topic, groupID)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from topic %s: %v", topic, err)
				continue
			}

			var event Event
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Error unmarshaling event from topic %s: %v", topic, err)
				continue
			}

			if err := handler(&event); err != nil {
				log.Printf("Error handling event %s from topic %s: %v", event.ID, topic, err)
				continue
			}

			log.Printf("Successfully processed event %s from topic %s", event.ID, topic)
		}
	}
}

// CreateTopics creates Kafka topics if they don't exist
func (k *KafkaClient) CreateTopics(ctx context.Context, topics []string) error {
	conn, err := k.dialer.DialLeader(ctx, "tcp", k.config.Brokers[0], "", 0)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             k.getTopicName(topic),
			NumPartitions:     6,
			ReplicationFactor: 1,
		}
	}

	if err := conn.CreateTopics(topicConfigs...); err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	log.Printf("Created Kafka topics: %v", topics)
	return nil
}

// Close closes all writers and readers
func (k *KafkaClient) Close() error {
	// Close all writers
	for topic, writer := range k.writers {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer for topic %s: %v", topic, err)
		}
	}

	// Close all readers
	for readerKey, reader := range k.readers {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing reader for %s: %v", readerKey, err)
		}
	}

	log.Println("Kafka client closed")
	return nil
}

// Helper functions for creating specific events

// NewUserEvent creates a user-related event
func NewUserEvent(eventType, userID string, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("%s-%d", userID, time.Now().UnixNano()),
		Type:      eventType,
		Source:    "auth-service",
		Subject:   userID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewMessageEvent creates a message-related event
func NewMessageEvent(eventType, messageID string, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("%s-%d", messageID, time.Now().UnixNano()),
		Type:      eventType,
		Source:    "messaging-service",
		Subject:   messageID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewDialogEvent creates a dialog-related event
func NewDialogEvent(eventType, dialogID string, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("%s-%d", dialogID, time.Now().UnixNano()),
		Type:      eventType,
		Source:    "messaging-service",
		Subject:   dialogID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewPaymentEvent creates a payment-related event
func NewPaymentEvent(eventType, transactionID string, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("%s-%d", transactionID, time.Now().UnixNano()),
		Type:      eventType,
		Source:    "payment-service",
		Subject:   transactionID,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// NewSystemEvent creates a system-related event
func NewSystemEvent(eventType, subject string, data map[string]interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("system-%d", time.Now().UnixNano()),
		Type:      eventType,
		Source:    "system",
		Subject:   subject,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// Health check
func (k *KafkaClient) HealthCheck(ctx context.Context) error {
	conn, err := k.dialer.DialContext(ctx, "tcp", k.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("Kafka health check failed: %w", err)
	}
	defer conn.Close()

	return nil
}

// GetAllTopics returns all topic names used by the application
func GetAllTopics() []string {
	return []string{
		UserRegisteredTopic,
		UserUpdatedTopic,
		UserDeletedTopic,
		UserLoginTopic,
		UserLogoutTopic,
		MessageSentTopic,
		MessageEditedTopic,
		MessageDeletedTopic,
		MessageDeliveredTopic,
		MessageReadTopic,
		DialogCreatedTopic,
		DialogUpdatedTopic,
		DialogDeletedTopic,
		ParticipantJoinedTopic,
		ParticipantLeftTopic,
		PaymentInitiatedTopic,
		PaymentCompletedTopic,
		PaymentFailedTopic,
		PaymentRefundedTopic,
		WalletCreatedTopic,
		WalletUpdatedTopic,
		NotificationSentTopic,
		NotificationDeliveredTopic,
		NotificationReadTopic,
		SystemAlertTopic,
		AuditLogTopic,
		MetricsUpdateTopic,
		PresenceUpdateTopic,
		TypingIndicatorTopic,
		ConnectionEventTopic,
	}
}