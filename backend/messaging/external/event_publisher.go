package external

import (
	"context"
	"encoding/json"
	"log"

	"tchat.dev/messaging/services"
	sharedModels "tchat.dev/shared/models"
)

// EventPublisher implements services.EventPublisher using a message broker
type EventPublisher struct {
	enabled bool
	// In a real implementation, this would connect to Kafka, RabbitMQ, etc.
	// For now, we'll log events for demonstration
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher() services.EventPublisher {
	return &EventPublisher{
		enabled: true,
	}
}

// Publish publishes an event to the message broker
func (e *EventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	if !e.enabled {
		return nil
	}

	// Convert event to JSON for logging
	eventData, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	// In a real implementation, this would publish to Kafka, RabbitMQ, etc.
	log.Printf("Publishing event %s to message broker:\n%s", event.Type, string(eventData))

	// TODO: Implement actual message broker integration
	// Examples:
	// - Kafka producer
	// - RabbitMQ publisher
	// - AWS SNS/SQS
	// - Google Pub/Sub
	// - Redis Pub/Sub

	return nil
}