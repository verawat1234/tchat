package handlers

import (
	"tchat.dev/notification/models"
)

// NotificationToResponse converts a Notification model to NotificationResponse DTO
func NotificationToResponse(notification *models.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:          notification.GetIDString(),
		RecipientID: notification.GetRecipientID(),
		Type:        notification.GetTypeString(),
		Channel:     notification.GetChannelString(),
		Subject:     notification.GetSubject(),
		Content:     notification.GetContent(),
		Status:      notification.GetStatusString(),
		Priority:    notification.GetPriorityString(),
		Read:        notification.Read,
		ReadAt:      nil, // Model doesn't have ReadAt field directly
		SentAt:      notification.SentAt,
		DeliveredAt: nil, // Model doesn't have DeliveredAt field directly
		FailedAt:    nil, // Model doesn't have FailedAt field directly
		ScheduledAt: notification.ScheduledAt,
		ExpiresAt:   notification.ExpiresAt,
		RetryCount:  notification.RetryCount,
		Metadata:    notification.Metadata,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}
}

// NotificationsToResponses converts a slice of Notifications to NotificationResponse DTOs
func NotificationsToResponses(notifications []*models.Notification) []*NotificationResponse {
	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = NotificationToResponse(notification)
	}
	return responses
}

// TemplateToResponse converts a NotificationTemplate model to TemplateResponse DTO
func TemplateToResponse(template *models.NotificationTemplate) *TemplateResponse {
	return &TemplateResponse{
		ID:        template.GetIDString(),
		Name:      template.Name,
		Type:      template.GetTypeString(),
		Category:  template.GetCategoryString(),
		Subject:   template.GetSubjectPtr(),
		Content:   template.GetContent(),
		Variables: template.Variables,
		Locales:   template.GetLocales(),
		Metadata:  template.GetMetadata(),
		Active:    template.GetActive(),
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}
}

// TemplatesToResponses converts a slice of NotificationTemplates to TemplateResponse DTOs
func TemplatesToResponses(templates []*models.NotificationTemplate) []*TemplateResponse {
	responses := make([]*TemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = TemplateToResponse(template)
	}
	return responses
}
