package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventBus provides event publishing and subscription functionality
type EventBus struct {
	subscribers map[string][]Subscriber
	handlers    map[string][]EventHandler
	middleware  []MiddlewareFunc
	config      *EventBusConfig
	mu          sync.RWMutex
	logger      Logger
}

// EventBusConfig holds event bus configuration
type EventBusConfig struct {
	MaxRetries       int
	RetryDelay       time.Duration
	HandlerTimeout   time.Duration
	EnableMetrics    bool
	BufferSize       int
	MaxConcurrency   int
	EnableDeadLetter bool
	DeadLetterTopic  string
}

// DefaultEventBusConfig returns default event bus configuration
func DefaultEventBusConfig() *EventBusConfig {
	return &EventBusConfig{
		MaxRetries:       3,
		RetryDelay:       1 * time.Second,
		HandlerTimeout:   30 * time.Second,
		EnableMetrics:    true,
		BufferSize:       1000,
		MaxConcurrency:   10,
		EnableDeadLetter: true,
		DeadLetterTopic:  "dead_letter",
	}
}

// Event represents a system event
type Event struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Subject     string                 `json:"subject"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version"`
	TraceID     string                 `json:"trace_id,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
}

// EventHandler represents a function that handles events
type EventHandler func(ctx context.Context, event *Event) error

// Subscriber represents an event subscriber
type Subscriber interface {
	Handle(ctx context.Context, event *Event) error
	GetEventTypes() []string
	GetID() string
}

// MiddlewareFunc represents event middleware
type MiddlewareFunc func(next EventHandler) EventHandler

// Logger interface for event bus logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// SimpleLogger provides a basic logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, err error, fields ...interface{}) {
	log.Printf("[ERROR] %s: %v %v", msg, err, fields)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

// EventMetrics tracks event bus metrics
type EventMetrics struct {
	EventsPublished   int64
	EventsProcessed   int64
	EventsFailed      int64
	HandlerErrors     int64
	AverageProcessingTime time.Duration
	mu                sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus(config *EventBusConfig, logger Logger) *EventBus {
	if config == nil {
		config = DefaultEventBusConfig()
	}

	if logger == nil {
		logger = &SimpleLogger{}
	}

	return &EventBus{
		subscribers: make(map[string][]Subscriber),
		handlers:    make(map[string][]EventHandler),
		middleware:  make([]MiddlewareFunc, 0),
		config:      config,
		logger:      logger,
	}
}

// Subscribe registers a subscriber for specific event types
func (eb *EventBus) Subscribe(subscriber Subscriber) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eventTypes := subscriber.GetEventTypes()
	if len(eventTypes) == 0 {
		return fmt.Errorf("subscriber must specify at least one event type")
	}

	for _, eventType := range eventTypes {
		eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
		eb.logger.Info("Subscriber registered", "event_type", eventType, "subscriber_id", subscriber.GetID())
	}

	return nil
}

// Unsubscribe removes a subscriber
func (eb *EventBus) Unsubscribe(subscriber Subscriber) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eventTypes := subscriber.GetEventTypes()
	for _, eventType := range eventTypes {
		subscribers := eb.subscribers[eventType]
		for i, sub := range subscribers {
			if sub.GetID() == subscriber.GetID() {
				// Remove subscriber from slice
				eb.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
				eb.logger.Info("Subscriber unregistered", "event_type", eventType, "subscriber_id", subscriber.GetID())
				break
			}
		}

		// Clean up empty event type
		if len(eb.subscribers[eventType]) == 0 {
			delete(eb.subscribers, eventType)
		}
	}

	return nil
}

// AddHandler registers an event handler for a specific event type
func (eb *EventBus) AddHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Info("Handler registered", "event_type", eventType)
}

// RemoveHandler removes an event handler (requires function comparison, limited use)
func (eb *EventBus) RemoveHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers := eb.handlers[eventType]
	handlerValue := reflect.ValueOf(handler)

	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == handlerValue.Pointer() {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			eb.logger.Info("Handler removed", "event_type", eventType)
			break
		}
	}

	if len(eb.handlers[eventType]) == 0 {
		delete(eb.handlers, eventType)
	}
}

// AddMiddleware adds middleware to the event processing pipeline
func (eb *EventBus) AddMiddleware(middleware MiddlewareFunc) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.middleware = append(eb.middleware, middleware)
	eb.logger.Info("Middleware added")
}

// Publish publishes an event to all relevant subscribers and handlers
func (eb *EventBus) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Set default values if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if event.Version == "" {
		event.Version = "1.0"
	}

	eb.logger.Info("Publishing event", "event_id", event.ID, "event_type", event.Type)

	// Get subscribers and handlers
	eb.mu.RLock()
	subscribers := eb.subscribers[event.Type]
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	// Process subscribers
	if len(subscribers) > 0 {
		for _, subscriber := range subscribers {
			go eb.processSubscriber(ctx, event, subscriber)
		}
	}

	// Process handlers
	if len(handlers) > 0 {
		for _, handler := range handlers {
			go eb.processHandler(ctx, event, handler)
		}
	}

	if len(subscribers) == 0 && len(handlers) == 0 {
		eb.logger.Warn("No subscribers or handlers for event type", "event_type", event.Type)
	}

	return nil
}

// PublishSync publishes an event synchronously
func (eb *EventBus) PublishSync(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Set default values if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if event.Version == "" {
		event.Version = "1.0"
	}

	eb.logger.Info("Publishing event synchronously", "event_id", event.ID, "event_type", event.Type)

	// Get subscribers and handlers
	eb.mu.RLock()
	subscribers := eb.subscribers[event.Type]
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	var errors []error

	// Process subscribers synchronously
	for _, subscriber := range subscribers {
		if err := eb.callSubscriber(ctx, event, subscriber); err != nil {
			errors = append(errors, err)
		}
	}

	// Process handlers synchronously
	for _, handler := range handlers {
		if err := eb.callHandler(ctx, event, handler); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("event processing failed with %d errors: %v", len(errors), errors)
	}

	return nil
}

// processSubscriber processes an event with a subscriber asynchronously
func (eb *EventBus) processSubscriber(ctx context.Context, event *Event, subscriber Subscriber) {
	err := eb.callSubscriber(ctx, event, subscriber)
	if err != nil {
		eb.logger.Error("Subscriber processing failed",
			err,
			"event_id", event.ID,
			"event_type", event.Type,
			"subscriber_id", subscriber.GetID())

		// Send to dead letter if enabled
		if eb.config.EnableDeadLetter {
			eb.sendToDeadLetter(ctx, event, err)
		}
	}
}

// processHandler processes an event with a handler asynchronously
func (eb *EventBus) processHandler(ctx context.Context, event *Event, handler EventHandler) {
	err := eb.callHandler(ctx, event, handler)
	if err != nil {
		eb.logger.Error("Handler processing failed",
			err,
			"event_id", event.ID,
			"event_type", event.Type)

		// Send to dead letter if enabled
		if eb.config.EnableDeadLetter {
			eb.sendToDeadLetter(ctx, event, err)
		}
	}
}

// callSubscriber calls a subscriber with retry logic and middleware
func (eb *EventBus) callSubscriber(ctx context.Context, event *Event, subscriber Subscriber) error {
	// Wrap subscriber call as handler
	handler := func(ctx context.Context, event *Event) error {
		return subscriber.Handle(ctx, event)
	}

	// Apply middleware
	finalHandler := eb.applyMiddleware(handler)

	// Execute with timeout and retry
	return eb.executeWithRetry(ctx, event, finalHandler)
}

// callHandler calls a handler with retry logic and middleware
func (eb *EventBus) callHandler(ctx context.Context, event *Event, handler EventHandler) error {
	// Apply middleware
	finalHandler := eb.applyMiddleware(handler)

	// Execute with timeout and retry
	return eb.executeWithRetry(ctx, event, finalHandler)
}

// applyMiddleware applies all middleware to a handler
func (eb *EventBus) applyMiddleware(handler EventHandler) EventHandler {
	// Apply middleware in reverse order
	for i := len(eb.middleware) - 1; i >= 0; i-- {
		handler = eb.middleware[i](handler)
	}
	return handler
}

// executeWithRetry executes a handler with timeout and retry logic
func (eb *EventBus) executeWithRetry(ctx context.Context, event *Event, handler EventHandler) error {
	var lastErr error

	for attempt := 0; attempt <= eb.config.MaxRetries; attempt++ {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx, eb.config.HandlerTimeout)

		// Execute handler
		err := handler(timeoutCtx, event)
		cancel()

		if err == nil {
			return nil // Success
		}

		lastErr = err
		eb.logger.Warn("Handler execution failed, retrying",
			"event_id", event.ID,
			"attempt", attempt+1,
			"error", err.Error())

		// Don't wait after the last attempt
		if attempt < eb.config.MaxRetries {
			time.Sleep(eb.config.RetryDelay * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	return fmt.Errorf("handler failed after %d attempts: %v", eb.config.MaxRetries+1, lastErr)
}

// sendToDeadLetter sends failed events to dead letter queue
func (eb *EventBus) sendToDeadLetter(ctx context.Context, event *Event, err error) {
	deadLetterEvent := &Event{
		ID:        uuid.New(),
		Type:      eb.config.DeadLetterTopic,
		Source:    "event_bus",
		Subject:   "failed_event",
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
		Data: map[string]interface{}{
			"original_event": event,
			"error":         err.Error(),
			"failed_at":     time.Now().UTC(),
		},
	}

	// Try to publish to dead letter (without retry to avoid infinite loop)
	eb.mu.RLock()
	deadLetterSubscribers := eb.subscribers[eb.config.DeadLetterTopic]
	deadLetterHandlers := eb.handlers[eb.config.DeadLetterTopic]
	eb.mu.RUnlock()

	for _, subscriber := range deadLetterSubscribers {
		go func(sub Subscriber) {
			if err := sub.Handle(ctx, deadLetterEvent); err != nil {
				eb.logger.Error("Dead letter subscriber failed", err)
			}
		}(subscriber)
	}

	for _, handler := range deadLetterHandlers {
		go func(h EventHandler) {
			if err := h(ctx, deadLetterEvent); err != nil {
				eb.logger.Error("Dead letter handler failed", err)
			}
		}(handler)
	}
}

// GetSubscribers returns all subscribers for an event type
func (eb *EventBus) GetSubscribers(eventType string) []Subscriber {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return eb.subscribers[eventType]
}

// GetHandlers returns all handlers for an event type
func (eb *EventBus) GetHandlers(eventType string) []EventHandler {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return eb.handlers[eventType]
}

// GetEventTypes returns all registered event types
func (eb *EventBus) GetEventTypes() []string {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	eventTypes := make([]string, 0, len(eb.subscribers)+len(eb.handlers))

	// Collect from subscribers
	for eventType := range eb.subscribers {
		eventTypes = append(eventTypes, eventType)
	}

	// Collect from handlers
	for eventType := range eb.handlers {
		found := false
		for _, existing := range eventTypes {
			if existing == eventType {
				found = true
				break
			}
		}
		if !found {
			eventTypes = append(eventTypes, eventType)
		}
	}

	return eventTypes
}

// Stop gracefully stops the event bus
func (eb *EventBus) Stop() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.logger.Info("Stopping event bus")

	// Clear all subscribers and handlers
	eb.subscribers = make(map[string][]Subscriber)
	eb.handlers = make(map[string][]EventHandler)
	eb.middleware = make([]MiddlewareFunc, 0)
}

// Common middleware functions

// LoggingMiddleware logs event processing
func LoggingMiddleware(logger Logger) MiddlewareFunc {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event *Event) error {
			start := time.Now()
			logger.Info("Processing event", "event_id", event.ID, "event_type", event.Type)

			err := next(ctx, event)

			duration := time.Since(start)
			if err != nil {
				logger.Error("Event processing failed", err, "event_id", event.ID, "duration", duration)
			} else {
				logger.Info("Event processed successfully", "event_id", event.ID, "duration", duration)
			}

			return err
		}
	}
}

// MetricsMiddleware collects processing metrics
func MetricsMiddleware(metrics *EventMetrics) MiddlewareFunc {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event *Event) error {
			start := time.Now()

			err := next(ctx, event)

			duration := time.Since(start)

			metrics.mu.Lock()
			metrics.EventsProcessed++
			if err != nil {
				metrics.EventsFailed++
			}
			// Simple moving average
			metrics.AverageProcessingTime = (metrics.AverageProcessingTime + duration) / 2
			metrics.mu.Unlock()

			return err
		}
	}
}

// RecoveryMiddleware recovers from panics in event handlers
func RecoveryMiddleware(logger Logger) MiddlewareFunc {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event *Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic recovered: %v", r)
					logger.Error("Event handler panicked", err, "event_id", event.ID)
				}
			}()

			return next(ctx, event)
		}
	}
}

// TracingMiddleware adds distributed tracing support
func TracingMiddleware() MiddlewareFunc {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event *Event) error {
			// Extract trace ID from context or generate new one
			traceID := getTraceIDFromContext(ctx)
			if traceID == "" {
				traceID = uuid.New().String()
			}

			// Set trace ID in event
			if event.TraceID == "" {
				event.TraceID = traceID
			}

			// Add trace ID to context
			ctx = setTraceIDInContext(ctx, traceID)

			return next(ctx, event)
		}
	}
}

// Utility functions for tracing (simplified implementation)
func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

func setTraceIDInContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

// EventBuilder helps build events with method chaining
type EventBuilder struct {
	event *Event
}

// NewEventBuilder creates a new event builder
func NewEventBuilder(eventType string) *EventBuilder {
	return &EventBuilder{
		event: &Event{
			ID:        uuid.New(),
			Type:      eventType,
			Timestamp: time.Now().UTC(),
			Version:   "1.0",
			Data:      make(map[string]interface{}),
			Metadata:  make(map[string]string),
		},
	}
}

// Source sets the event source
func (eb *EventBuilder) Source(source string) *EventBuilder {
	eb.event.Source = source
	return eb
}

// Subject sets the event subject
func (eb *EventBuilder) Subject(subject string) *EventBuilder {
	eb.event.Subject = subject
	return eb
}

// Data sets event data
func (eb *EventBuilder) Data(data map[string]interface{}) *EventBuilder {
	eb.event.Data = data
	return eb
}

// AddData adds a single data field
func (eb *EventBuilder) AddData(key string, value interface{}) *EventBuilder {
	eb.event.Data[key] = value
	return eb
}

// AddMetadata adds metadata
func (eb *EventBuilder) AddMetadata(key, value string) *EventBuilder {
	eb.event.Metadata[key] = value
	return eb
}

// TraceID sets the trace ID
func (eb *EventBuilder) TraceID(traceID string) *EventBuilder {
	eb.event.TraceID = traceID
	return eb
}

// Build returns the constructed event
func (eb *EventBuilder) Build() *Event {
	return eb.event
}

// JSON converts event to JSON
func (e *Event) JSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON creates an event from JSON
func FromJSON(data string) (*Event, error) {
	var event Event
	err := json.Unmarshal([]byte(data), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}