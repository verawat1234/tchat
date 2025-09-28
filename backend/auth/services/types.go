package services

import (
	"context"

	sharedModels "tchat.dev/shared/models"
)

// EventPublisher interface for publishing events across all services
type EventPublisher interface {
	Publish(ctx context.Context, event *sharedModels.Event) error
}