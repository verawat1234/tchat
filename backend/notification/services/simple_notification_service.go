package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"tchat.dev/notification/models"
)

// SimpleNotificationService provides basic notification functionality
type SimpleNotificationService struct {
	templates map[string]*models.NotificationTemplate
}

// NewSimpleNotificationService creates a new simple notification service
func NewSimpleNotificationService() *SimpleNotificationService {
	return &SimpleNotificationService{
		templates: make(map[string]*models.NotificationTemplate),
	}
}

// SendNotification sends a basic notification
func (s *SimpleNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, title, body string) (*models.Notification, error) {
	notification := &models.Notification{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		Type:      models.NotificationTypePush,
		Category:  models.NotificationCategorySystem,
		Priority:  models.PriorityNormal,
		Status:    models.DeliveryStatusPending,
		SentAt:    func() *time.Time { t := time.Now(); return &t }(),
		Targeting: models.Targeting{
			AudienceType: models.AudienceTypeUser,
			UserIDs:      []uuid.UUID{userID},
		},
	}

	// In a real implementation, this would save to database
	// For now, just return the created notification
	return notification, nil
}

// GetTemplate gets a notification template
func (s *SimpleNotificationService) GetTemplate(templateType string) *models.NotificationTemplate {
	if template, exists := s.templates[templateType]; exists {
		return template
	}

	// Return a default template
	return &models.NotificationTemplate{
		ID:            uuid.New(),
		Type:          models.NotificationType(templateType),
		Subject:       "Default Subject",
		Body:          "Default notification body",
		TitleTemplate: "{{title}}",
		BodyTemplate:  "{{body}}",
	}
}