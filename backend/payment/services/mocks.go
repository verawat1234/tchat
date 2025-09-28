package services

import (
	"context"

	sharedModels "tchat.dev/shared/models"
)

// NoOpEventPublisher is a no-operation event publisher for testing
type NoOpEventPublisher struct{}

func (n *NoOpEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	return nil
}