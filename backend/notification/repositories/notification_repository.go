package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/notification/models"
)

// NotificationRepositoryImpl implements the NotificationRepository interface
type NotificationRepositoryImpl struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &NotificationRepositoryImpl{db: db}
}

// Create creates a new notification
func (r *NotificationRepositoryImpl) Create(ctx context.Context, notification *models.Notification) error {
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetByID retrieves a notification by ID
func (r *NotificationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// Update updates an existing notification
func (r *NotificationRepositoryImpl) Update(ctx context.Context, notification *models.Notification) error {
	notification.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete deletes a notification by ID
func (r *NotificationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, "id = ?", id).Error
}

// GetByUserID retrieves notifications for a specific user with pagination
func (r *NotificationRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetByChannel retrieves notifications by channel with pagination
func (r *NotificationRepositoryImpl) GetByChannel(ctx context.Context, channel models.NotificationType, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.WithContext(ctx).
		Where("type = ?", channel).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetPendingNotifications retrieves notifications with pending status
func (r *NotificationRepositoryImpl) GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.WithContext(ctx).
		Where("status = ?", models.DeliveryStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetFailedNotifications retrieves notifications with failed status
func (r *NotificationRepositoryImpl) GetFailedNotifications(ctx context.Context, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.WithContext(ctx).
		Where("status = ?", models.DeliveryStatusFailed).
		Where("retry_count < max_retries OR max_retries = 0").
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetScheduledNotifications retrieves notifications scheduled for delivery
func (r *NotificationRepositoryImpl) GetScheduledNotifications(ctx context.Context, before time.Time, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.WithContext(ctx).
		Where("scheduled_at IS NOT NULL").
		Where("scheduled_at <= ?", before).
		Where("status = ?", models.DeliveryStatusPending).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// MarkAsDelivered marks a notification as delivered
func (r *NotificationRepositoryImpl) MarkAsDelivered(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.DeliveryStatusDelivered,
			"delivered_at": &now,
			"updated_at":   now,
		}).Error
}

// MarkAsFailed marks a notification as failed with a reason
func (r *NotificationRepositoryImpl) MarkAsFailed(ctx context.Context, id uuid.UUID, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         models.DeliveryStatusFailed,
			"failure_reason": reason,
			"failed_at":      &now,
			"updated_at":     now,
		}).Error
}

// MarkAsRead marks a notification as read
func (r *NotificationRepositoryImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"read":       true,
			"read_at":    &now,
			"updated_at": now,
		}).Error
}

// UpdateStatus updates the status of a notification
func (r *NotificationRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.DeliveryStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}

	// Set appropriate timestamp based on status
	switch status {
	case models.DeliveryStatusSent:
		updates["sent_at"] = &now
	case models.DeliveryStatusDelivered:
		updates["delivered_at"] = &now
	case models.DeliveryStatusFailed:
		updates["failed_at"] = &now
	case models.DeliveryStatusRead:
		updates["read"] = true
		updates["read_at"] = &now
	}

	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CreateBatch creates multiple notifications in a single transaction
func (r *NotificationRepositoryImpl) CreateBatch(ctx context.Context, notifications []*models.Notification) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for _, notification := range notifications {
			if notification.ID == uuid.Nil {
				notification.ID = uuid.New()
			}
			notification.CreatedAt = now
			notification.UpdatedAt = now
		}
		return tx.CreateInBatches(notifications, 100).Error
	})
}

// MarkBatchAsDelivered marks multiple notifications as delivered
func (r *NotificationRepositoryImpl) MarkBatchAsDelivered(ctx context.Context, ids []uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":       models.DeliveryStatusDelivered,
			"delivered_at": &now,
			"updated_at":   now,
		}).Error
}

// GetDeliveryStats retrieves delivery statistics for a date range
func (r *NotificationRepositoryImpl) GetDeliveryStats(ctx context.Context, startDate, endDate time.Time) (*models.DeliveryStats, error) {
	var stats models.DeliveryStats

	// Total sent notifications
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalSent).Error
	if err != nil {
		return nil, err
	}

	// Delivered notifications
	err = r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.DeliveryStatusDelivered).
		Count(&stats.TotalDelivered).Error
	if err != nil {
		return nil, err
	}

	// Failed notifications
	err = r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.DeliveryStatusFailed).
		Count(&stats.TotalFailed).Error
	if err != nil {
		return nil, err
	}

	// Read notifications
	err = r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("read = ?", true).
		Count(&stats.TotalRead).Error
	if err != nil {
		return nil, err
	}

	// Calculate rates
	if stats.TotalSent > 0 {
		stats.DeliveryRate = float64(stats.TotalDelivered) / float64(stats.TotalSent) * 100
		stats.FailureRate = float64(stats.TotalFailed) / float64(stats.TotalSent) * 100
	}
	if stats.TotalDelivered > 0 {
		stats.ReadRate = float64(stats.TotalRead) / float64(stats.TotalDelivered) * 100
	}

	return &stats, nil
}

// CleanupOldNotifications removes old notifications
func (r *NotificationRepositoryImpl) CleanupOldNotifications(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", before).
		Where("status IN ?", []models.DeliveryStatus{
			models.DeliveryStatusDelivered,
			models.DeliveryStatusFailed,
			models.DeliveryStatusExpired,
		}).
		Delete(&models.Notification{}).Error
}

// CountByUserAndDateRange counts notifications for a user in a date range
func (r *NotificationRepositoryImpl) CountByUserAndDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count).Error
	return count, err
}