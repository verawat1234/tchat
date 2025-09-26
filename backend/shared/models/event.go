package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventType string
type EventStatus string
type EventCategory string
type Severity string
type DataFormat string

const (
	// Event Types - Domain Events
	EventTypeUserRegistered       EventType = "user.registered"
	EventTypeUserProfileUpdated   EventType = "user.profile_updated"
	EventTypeUserKYCVerified      EventType = "user.kyc_verified"
	EventTypeUserSessionCreated   EventType = "user.session_created"
	EventTypeUserSessionExpired   EventType = "user.session_expired"

	// Chat Events
	EventTypeMessageSent          EventType = "message.sent"
	EventTypeMessageDelivered     EventType = "message.delivered"
	EventTypeMessageRead          EventType = "message.read"
	EventTypeDialogCreated        EventType = "dialog.created"
	EventTypeDialogParticipantAdded EventType = "dialog.participant_added"
	EventTypeUserPresenceChanged  EventType = "user.presence_changed"

	// Payment Events
	EventTypePaymentInitiated     EventType = "payment.initiated"
	EventTypePaymentCompleted     EventType = "payment.completed"
	EventTypePaymentFailed        EventType = "payment.failed"
	EventTypeWalletBalanceChanged EventType = "wallet.balance_changed"
	EventTypeTransactionCreated   EventType = "transaction.created"

	// Commerce Events
	EventTypeOrderCreated         EventType = "order.created"
	EventTypeOrderUpdated         EventType = "order.updated"
	EventTypeOrderFulfilled       EventType = "order.fulfilled"
	EventTypeOrderCancelled       EventType = "order.cancelled"
	EventTypeProductCreated       EventType = "product.created"
	EventTypeProductUpdated       EventType = "product.updated"
	EventTypeShopCreated          EventType = "shop.created"
	EventTypeShopStatusChanged    EventType = "shop.status_changed"

	// System Events
	EventTypeSystemStartup        EventType = "system.startup"
	EventTypeSystemShutdown       EventType = "system.shutdown"
	EventTypeSystemHealthCheck    EventType = "system.health_check"
	EventTypeSystemBackupCreated  EventType = "system.backup_created"
	EventTypeSystemMigrationStarted EventType = "system.migration_started"
	EventTypeSystemMigrationCompleted EventType = "system.migration_completed"

	// Security Events
	EventTypeSecurityLoginAttempt EventType = "security.login_attempt"
	EventTypeSecurityLoginSuccess EventType = "security.login_success"
	EventTypeSecurityLoginFailed  EventType = "security.login_failed"
	EventTypeSecurityPasswordChanged EventType = "security.password_changed"
	EventTypeSecuritySuspiciousActivity EventType = "security.suspicious_activity"
	EventTypeSecurityAccountLocked EventType = "security.account_locked"

	// Notification Events
	EventTypeNotificationSent     EventType = "notification.sent"
	EventTypeNotificationDelivered EventType = "notification.delivered"
	EventTypeNotificationFailed   EventType = "notification.failed"

	// Event Status
	EventStatusPending    EventStatus = "pending"
	EventStatusProcessing EventStatus = "processing"
	EventStatusProcessed  EventStatus = "processed"
	EventStatusFailed     EventStatus = "failed"
	EventStatusRetrying   EventStatus = "retrying"
	EventStatusSkipped    EventStatus = "skipped"
	EventStatusExpired    EventStatus = "expired"

	// Event Categories
	EventCategoryDomain      EventCategory = "domain"
	EventCategoryIntegration EventCategory = "integration"
	EventCategorySystem      EventCategory = "system"
	EventCategorySecurity    EventCategory = "security"
	EventCategoryAudit       EventCategory = "audit"
	EventCategoryMetrics     EventCategory = "metrics"
	EventCategoryError       EventCategory = "error"

	// Severity Levels
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
	SeverityDebug    Severity = "debug"

	// Data Formats
	DataFormatJSON     DataFormat = "json"
	DataFormatXML      DataFormat = "xml"
	DataFormatText     DataFormat = "text"
	DataFormatBinary   DataFormat = "binary"
	DataFormatAvro     DataFormat = "avro"
	DataFormatProtobuf DataFormat = "protobuf"
)

type EventMetadata struct {
	Source      string            `json:"source"`        // Service that generated the event
	SourceIP    string            `json:"source_ip"`
	UserAgent   string            `json:"user_agent"`
	RequestID   string            `json:"request_id"`
	SessionID   string            `json:"session_id"`
	TraceID     string            `json:"trace_id"`
	SpanID      string            `json:"span_id"`
	Version     string            `json:"version"`       // Event schema version
	Environment string            `json:"environment"`   // dev, staging, production
	Region      string            `json:"region"`        // Southeast Asian region
	Country     string            `json:"country"`       // TH, SG, ID, MY, PH, VN
	Headers     map[string]string `json:"headers"`
	Tags        []string          `json:"tags"`
}

type RetryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	Interval      time.Duration `json:"interval"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxInterval   time.Duration `json:"max_interval"`
}

type ProcessingResult struct {
	Success       bool              `json:"success"`
	ProcessedAt   time.Time         `json:"processed_at"`
	ProcessingTime time.Duration    `json:"processing_time"`
	HandlerName   string            `json:"handler_name"`
	HandlerVersion string           `json:"handler_version"`
	ErrorMessage  string            `json:"error_message,omitempty"`
	ErrorCode     string            `json:"error_code,omitempty"`
	Metadata      map[string]any    `json:"metadata,omitempty"`
}

type Event struct {
	ID           uuid.UUID         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Type         EventType         `json:"type" gorm:"type:varchar(100);not null;index"`
	Category     EventCategory     `json:"category" gorm:"type:varchar(30);not null;index"`
	Severity     Severity          `json:"severity" gorm:"type:varchar(20);not null;index"`
	Status       EventStatus       `json:"status" gorm:"type:varchar(20);default:'pending';index"`

	// Event Content
	Subject      string            `json:"subject" gorm:"type:varchar(255);not null"`
	Description  string            `json:"description" gorm:"type:text"`
	Data         json.RawMessage   `json:"data" gorm:"type:json"`
	DataFormat   DataFormat        `json:"data_format" gorm:"type:varchar(20);default:'json'"`
	DataSchema   string            `json:"data_schema" gorm:"type:varchar(100)"`

	// Event Context
	AggregateID   string           `json:"aggregate_id" gorm:"type:varchar(100);index"`
	AggregateType string           `json:"aggregate_type" gorm:"type:varchar(50);index"`
	EventVersion  int              `json:"event_version" gorm:"default:1"`

	// Processing Information
	Metadata      EventMetadata    `json:"metadata" gorm:"type:json"`
	RetryConfig   RetryConfig      `json:"retry_config" gorm:"type:json"`
	RetryCount    int              `json:"retry_count" gorm:"default:0"`
	LastRetryAt   *time.Time       `json:"last_retry_at,omitempty"`
	NextRetryAt   *time.Time       `json:"next_retry_at,omitempty" gorm:"index"`

	// Processing Results
	ProcessingResults []ProcessingResult `json:"processing_results" gorm:"type:json"`

	// Timing
	OccurredAt    time.Time        `json:"occurred_at" gorm:"not null;index"`
	ProcessedAt   *time.Time       `json:"processed_at,omitempty" gorm:"index"`
	ExpiresAt     *time.Time       `json:"expires_at,omitempty" gorm:"index"`

	// Correlation and Causation
	ParentEventID   *uuid.UUID     `json:"parent_event_id,omitempty" gorm:"type:varchar(36);index"`
	CorrelationID   string         `json:"correlation_id" gorm:"type:varchar(100);index"`
	CausationID     string         `json:"causation_id" gorm:"type:varchar(100);index"`

	// Audit Fields
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type EventSubscription struct {
	ID            uuid.UUID        `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name          string           `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	ServiceName   string           `json:"service_name" gorm:"type:varchar(100);not null;index"`
	EventTypes    []string         `json:"event_types" gorm:"type:json;not null"`
	FilterQuery   string           `json:"filter_query" gorm:"type:text"`
	Endpoint      string           `json:"endpoint" gorm:"type:varchar(500);not null"`
	Method        string           `json:"method" gorm:"type:varchar(10);default:'POST'"`
	Headers       map[string]string `json:"headers" gorm:"type:json"`
	Timeout       time.Duration    `json:"timeout" gorm:"default:30000000000"` // 30 seconds in nanoseconds
	RetryConfig   RetryConfig      `json:"retry_config" gorm:"type:json"`
	IsActive      bool             `json:"is_active" gorm:"default:true;index"`

	// Health and Monitoring
	SuccessCount  int64            `json:"success_count" gorm:"default:0"`
	FailureCount  int64            `json:"failure_count" gorm:"default:0"`
	LastSuccess   *time.Time       `json:"last_success,omitempty"`
	LastFailure   *time.Time       `json:"last_failure,omitempty"`
	HealthStatus  string           `json:"health_status" gorm:"type:varchar(20);default:'healthy'"`

	// Audit Fields
	CreatedBy     uuid.UUID        `json:"created_by" gorm:"type:varchar(36);index"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type EventLog struct {
	ID               uuid.UUID        `json:"id" gorm:"primaryKey;type:varchar(36)"`
	EventID          uuid.UUID        `json:"event_id" gorm:"type:varchar(36);not null;index"`
	SubscriptionID   uuid.UUID        `json:"subscription_id" gorm:"type:varchar(36);not null;index"`
	Status           string           `json:"status" gorm:"type:varchar(20);not null;index"`
	AttemptNumber    int              `json:"attempt_number" gorm:"default:1"`
	ResponseCode     int              `json:"response_code"`
	ResponseBody     string           `json:"response_body" gorm:"type:text"`
	ErrorMessage     string           `json:"error_message" gorm:"type:text"`
	ProcessingTime   time.Duration    `json:"processing_time"`
	RequestPayload   json.RawMessage  `json:"request_payload" gorm:"type:json"`
	ResponseHeaders  map[string]string `json:"response_headers" gorm:"type:json"`
	AttemptedAt      time.Time        `json:"attempted_at" gorm:"not null;index"`
	CreatedAt        time.Time        `json:"created_at" gorm:"autoCreateTime"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	if e.RetryConfig.MaxRetries == 0 {
		e.RetryConfig = GetDefaultRetryConfig()
	}
	if e.CorrelationID == "" {
		e.CorrelationID = uuid.New().String()
	}
	if e.CausationID == "" {
		e.CausationID = e.CorrelationID
	}
	return nil
}

func (e *Event) IsExpired() bool {
	return e.ExpiresAt != nil && time.Now().After(*e.ExpiresAt)
}

func (e *Event) CanRetry() bool {
	return e.RetryCount < e.RetryConfig.MaxRetries && !e.IsExpired()
}

func (e *Event) CalculateNextRetry() time.Time {
	if e.RetryCount == 0 {
		return time.Now().Add(e.RetryConfig.Interval)
	}

	// Exponential backoff with jitter
	backoffDuration := time.Duration(float64(e.RetryConfig.Interval) *
		func() float64 {
			result := 1.0
			for i := 0; i < e.RetryCount; i++ {
				result *= e.RetryConfig.BackoffFactor
			}
			return result
		}())

	if backoffDuration > e.RetryConfig.MaxInterval {
		backoffDuration = e.RetryConfig.MaxInterval
	}

	// Add jitter (±25%)
	jitterFactor := 0.25 * (2*float64(time.Now().UnixNano()%2) - 1)
	jitter := time.Duration(float64(backoffDuration) * jitterFactor)
	return time.Now().Add(backoffDuration + jitter)
}

func (e *Event) AddProcessingResult(result ProcessingResult) {
	e.ProcessingResults = append(e.ProcessingResults, result)
	if result.Success {
		e.Status = EventStatusProcessed
		now := time.Now()
		e.ProcessedAt = &now
	} else {
		if e.CanRetry() {
			e.Status = EventStatusRetrying
			nextRetry := e.CalculateNextRetry()
			e.NextRetryAt = &nextRetry
			e.RetryCount++
		} else {
			e.Status = EventStatusFailed
		}
	}
}

func (e *Event) GetLatestProcessingResult() *ProcessingResult {
	if len(e.ProcessingResults) == 0 {
		return nil
	}
	return &e.ProcessingResults[len(e.ProcessingResults)-1]
}

func (e *Event) GetProcessingResultsForHandler(handlerName string) []ProcessingResult {
	results := make([]ProcessingResult, 0)
	for _, result := range e.ProcessingResults {
		if result.HandlerName == handlerName {
			results = append(results, result)
		}
	}
	return results
}

func (e *Event) UnmarshalData(v interface{}) error {
	return json.Unmarshal(e.Data, v)
}

func (e *Event) MarshalData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.Data = data
	return nil
}

func (e *Event) GetDataAsMap() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(e.Data, &data)
	return data, err
}

func (e *Event) MatchesFilter(filterQuery string) bool {
	// Simple filter implementation - in production, use a proper query engine
	if filterQuery == "" {
		return true
	}

	// Example filters:
	// "type:user.registered"
	// "severity:error"
	// "aggregate_type:user"
	// "country:TH"

	filters := strings.Split(filterQuery, " AND ")
	for _, filter := range filters {
		parts := strings.SplitN(strings.TrimSpace(filter), ":", 2)
		if len(parts) != 2 {
			continue
		}

		field, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch field {
		case "type":
			if string(e.Type) != value {
				return false
			}
		case "category":
			if string(e.Category) != value {
				return false
			}
		case "severity":
			if string(e.Severity) != value {
				return false
			}
		case "aggregate_type":
			if e.AggregateType != value {
				return false
			}
		case "country":
			if e.Metadata.Country != value {
				return false
			}
		case "source":
			if e.Metadata.Source != value {
				return false
			}
		}
	}

	return true
}

func (es *EventSubscription) BeforeCreate(tx *gorm.DB) error {
	if es.ID == uuid.Nil {
		es.ID = uuid.New()
	}
	if es.RetryConfig.MaxRetries == 0 {
		es.RetryConfig = GetDefaultRetryConfig()
	}
	return nil
}

func (es *EventSubscription) IsHealthy() bool {
	return es.HealthStatus == "healthy"
}

func (es *EventSubscription) UpdateHealthStatus() {
	total := es.SuccessCount + es.FailureCount
	if total == 0 {
		es.HealthStatus = "unknown"
		return
	}

	successRate := float64(es.SuccessCount) / float64(total)

	switch {
	case successRate >= 0.95:
		es.HealthStatus = "healthy"
	case successRate >= 0.80:
		es.HealthStatus = "degraded"
	default:
		es.HealthStatus = "unhealthy"
	}
}

func (es *EventSubscription) RecordSuccess() {
	es.SuccessCount++
	now := time.Now()
	es.LastSuccess = &now
	es.UpdateHealthStatus()
}

func (es *EventSubscription) RecordFailure() {
	es.FailureCount++
	now := time.Now()
	es.LastFailure = &now
	es.UpdateHealthStatus()
}

func (es *EventSubscription) MatchesEvent(event *Event) bool {
	if !es.IsActive {
		return false
	}

	// Check if event type matches
	eventTypeMatches := false
	for _, eventType := range es.EventTypes {
		if eventType == "*" || eventType == string(event.Type) {
			eventTypeMatches = true
			break
		}
		// Support wildcard patterns like "user.*"
		if strings.HasSuffix(eventType, "*") {
			prefix := strings.TrimSuffix(eventType, "*")
			if strings.HasPrefix(string(event.Type), prefix) {
				eventTypeMatches = true
				break
			}
		}
	}

	if !eventTypeMatches {
		return false
	}

	// Check filter query
	return event.MatchesFilter(es.FilterQuery)
}

func (el *EventLog) BeforeCreate(tx *gorm.DB) error {
	if el.ID == uuid.Nil {
		el.ID = uuid.New()
	}
	if el.AttemptedAt.IsZero() {
		el.AttemptedAt = time.Now()
	}
	return nil
}

// Response structures for API
type EventResponse struct {
	ID              uuid.UUID         `json:"id"`
	Type            EventType         `json:"type"`
	Category        EventCategory     `json:"category"`
	Severity        Severity          `json:"severity"`
	Status          EventStatus       `json:"status"`
	Subject         string            `json:"subject"`
	Description     string            `json:"description"`
	Data            json.RawMessage   `json:"data"`
	AggregateID     string            `json:"aggregate_id"`
	AggregateType   string            `json:"aggregate_type"`
	EventVersion    int               `json:"event_version"`
	RetryCount      int               `json:"retry_count"`
	OccurredAt      time.Time         `json:"occurred_at"`
	ProcessedAt     *time.Time        `json:"processed_at,omitempty"`
	CorrelationID   string            `json:"correlation_id"`
	CausationID     string            `json:"causation_id"`
	CreatedAt       time.Time         `json:"created_at"`
}

type EventSubscriptionResponse struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	ServiceName  string            `json:"service_name"`
	EventTypes   []string          `json:"event_types"`
	FilterQuery  string            `json:"filter_query"`
	Endpoint     string            `json:"endpoint"`
	IsActive     bool              `json:"is_active"`
	SuccessCount int64             `json:"success_count"`
	FailureCount int64             `json:"failure_count"`
	HealthStatus string            `json:"health_status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// Request structures for API
type CreateEventRequest struct {
	Type          EventType      `json:"type" binding:"required"`
	Category      EventCategory  `json:"category" binding:"required"`
	Severity      Severity       `json:"severity" binding:"required"`
	Subject       string         `json:"subject" binding:"required"`
	Description   string         `json:"description"`
	Data          interface{}    `json:"data"`
	AggregateID   string         `json:"aggregate_id"`
	AggregateType string         `json:"aggregate_type"`
	EventVersion  int            `json:"event_version"`
	ExpiresAt     *time.Time     `json:"expires_at,omitempty"`
	ParentEventID *uuid.UUID     `json:"parent_event_id,omitempty"`
	CorrelationID string         `json:"correlation_id"`
	CausationID   string         `json:"causation_id"`
	Metadata      EventMetadata  `json:"metadata"`
}

type CreateEventSubscriptionRequest struct {
	Name         string            `json:"name" binding:"required"`
	ServiceName  string            `json:"service_name" binding:"required"`
	EventTypes   []string          `json:"event_types" binding:"required"`
	FilterQuery  string            `json:"filter_query"`
	Endpoint     string            `json:"endpoint" binding:"required"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers"`
	Timeout      *time.Duration    `json:"timeout"`
	RetryConfig  *RetryConfig      `json:"retry_config"`
}

type UpdateEventSubscriptionRequest struct {
	Name         *string           `json:"name"`
	EventTypes   []string          `json:"event_types"`
	FilterQuery  *string           `json:"filter_query"`
	Endpoint     *string           `json:"endpoint"`
	Method       *string           `json:"method"`
	Headers      map[string]string `json:"headers"`
	Timeout      *time.Duration    `json:"timeout"`
	RetryConfig  *RetryConfig      `json:"retry_config"`
	IsActive     *bool             `json:"is_active"`
}

// Manager for business logic
type EventManager struct{}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) ValidateEvent(event *Event) error {
	if event.Type == "" {
		return fmt.Errorf("event type is required")
	}

	if event.Subject == "" {
		return fmt.Errorf("event subject is required")
	}

	if event.Category == "" {
		return fmt.Errorf("event category is required")
	}

	if event.Severity == "" {
		return fmt.Errorf("event severity is required")
	}

	return nil
}

func (em *EventManager) CreateEventFromRequest(req CreateEventRequest) (*Event, error) {
	event := &Event{
		ID:            uuid.New(),
		Type:          req.Type,
		Category:      req.Category,
		Severity:      req.Severity,
		Subject:       req.Subject,
		Description:   req.Description,
		AggregateID:   req.AggregateID,
		AggregateType: req.AggregateType,
		EventVersion:  req.EventVersion,
		ExpiresAt:     req.ExpiresAt,
		ParentEventID: req.ParentEventID,
		CorrelationID: req.CorrelationID,
		CausationID:   req.CausationID,
		Metadata:      req.Metadata,
		OccurredAt:    time.Now(),
		Status:        EventStatusPending,
		RetryConfig:   GetDefaultRetryConfig(),
	}

	if req.Data != nil {
		if err := event.MarshalData(req.Data); err != nil {
			return nil, fmt.Errorf("failed to marshal event data: %w", err)
		}
	}

	return event, em.ValidateEvent(event)
}

func (em *EventManager) FindMatchingSubscriptions(event *Event, subscriptions []*EventSubscription) []*EventSubscription {
	matches := make([]*EventSubscription, 0)
	for _, subscription := range subscriptions {
		if subscription.MatchesEvent(event) {
			matches = append(matches, subscription)
		}
	}
	return matches
}

func (em *EventManager) GetSupportedEventTypes() []EventType {
	return []EventType{
		EventTypeUserRegistered,
		EventTypeUserProfileUpdated,
		EventTypeUserKYCVerified,
		EventTypeMessageSent,
		EventTypeMessageDelivered,
		EventTypePaymentInitiated,
		EventTypePaymentCompleted,
		EventTypeOrderCreated,
		EventTypeOrderFulfilled,
		EventTypeSystemStartup,
		EventTypeSecurityLoginAttempt,
		EventTypeNotificationSent,
	}
}

func GetDefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		Interval:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxInterval:   5 * time.Minute,
	}
}

func GetDefaultEventMetadata(source, environment, region, country string) EventMetadata {
	return EventMetadata{
		Source:      source,
		Version:     "1.0",
		Environment: environment,
		Region:      region,
		Country:     country,
		Headers:     make(map[string]string),
		Tags:        make([]string, 0),
	}
}