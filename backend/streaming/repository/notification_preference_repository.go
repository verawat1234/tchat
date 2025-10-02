package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/streaming/models"
)

// NotificationPreferenceRepository defines the interface for notification preference data access
type NotificationPreferenceRepository interface {
	// Create creates a new notification preference
	Create(ctx context.Context, pref *models.NotificationPreference) error

	// GetByUserID retrieves notification preferences for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error)

	// Update updates notification preferences for a specific user
	Update(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error

	// Delete removes notification preferences for a specific user
	Delete(ctx context.Context, userID uuid.UUID) error

	// GetOrCreateDefault retrieves existing preferences or creates default ones
	GetOrCreateDefault(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error)

	// Upsert creates or updates notification preferences
	Upsert(ctx context.Context, pref *models.NotificationPreference) error
}

// notificationPreferenceRepository implements the NotificationPreferenceRepository interface
type notificationPreferenceRepository struct {
	db *gorm.DB
}

// NewNotificationPreferenceRepository creates a new notification preference repository
func NewNotificationPreferenceRepository(db *gorm.DB) NotificationPreferenceRepository {
	return &notificationPreferenceRepository{
		db: db,
	}
}

// Create creates a new notification preference
func (r *notificationPreferenceRepository) Create(ctx context.Context, pref *models.NotificationPreference) error {
	// Validate quiet hours format if provided
	if err := r.validateQuietHours(pref); err != nil {
		return fmt.Errorf("invalid quiet hours: %w", err)
	}

	// Validate timezone if provided
	if err := r.validateTimezone(pref); err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(pref).Error; err != nil {
		return fmt.Errorf("failed to create notification preference: %w", err)
	}

	return nil
}

// GetByUserID retrieves notification preferences for a specific user
func (r *notificationPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification preference not found for user %s", userID.String())
		}
		return nil, fmt.Errorf("failed to get notification preference: %w", err)
	}

	return &pref, nil
}

// Update updates notification preferences for a specific user
func (r *notificationPreferenceRepository) Update(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// If updating quiet hours, validate them
	if quietStart, ok := updates["quiet_hours_start"]; ok {
		if quietStart != nil {
			if startTime, ok := quietStart.(*time.Time); ok {
				if err := r.validateTime(startTime); err != nil {
					return fmt.Errorf("invalid quiet_hours_start: %w", err)
				}
			}
		}
	}

	if quietEnd, ok := updates["quiet_hours_end"]; ok {
		if quietEnd != nil {
			if endTime, ok := quietEnd.(*time.Time); ok {
				if err := r.validateTime(endTime); err != nil {
					return fmt.Errorf("invalid quiet_hours_end: %w", err)
				}
			}
		}
	}

	// If updating timezone, validate it
	if timezone, ok := updates["timezone"]; ok {
		if timezone != nil {
			if err := r.validateTimezoneString(timezone); err != nil {
				return fmt.Errorf("invalid timezone: %w", err)
			}
		}
	}

	result := r.db.WithContext(ctx).
		Model(&models.NotificationPreference{}).
		Where("user_id = ?", userID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update notification preference: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification preference not found for user %s", userID.String())
	}

	return nil
}

// Delete removes notification preferences for a specific user
func (r *notificationPreferenceRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.NotificationPreference{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete notification preference: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification preference not found for user %s", userID.String())
	}

	return nil
}

// GetOrCreateDefault retrieves existing preferences or creates default ones
func (r *notificationPreferenceRepository) GetOrCreateDefault(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error) {
	// Try to get existing preference
	pref, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return pref, nil
	}

	// Check if error is "not found"
	if err != nil && err.Error() != fmt.Sprintf("notification preference not found for user %s", userID.String()) {
		return nil, fmt.Errorf("failed to check existing preference: %w", err)
	}

	// Create default preference
	defaultPref := &models.NotificationPreference{
		UserID:              userID,
		PushEnabled:         true,  // All channels enabled by default
		InAppEnabled:        true,
		EmailEnabled:        true,
		StoreStreamsEnabled: true,  // Both stream types enabled
		VideoStreamsEnabled: true,
		QuietHoursStart:     nil,   // No quiet hours initially
		QuietHoursEnd:       nil,
		Timezone:            sql.NullString{Valid: false},
		UpdatedAt:           time.Now(),
	}

	// Create the default preference
	if err := r.Create(ctx, defaultPref); err != nil {
		return nil, fmt.Errorf("failed to create default notification preference: %w", err)
	}

	return defaultPref, nil
}

// validateQuietHours validates quiet hours format (HH:MM:SS)
func (r *notificationPreferenceRepository) validateQuietHours(pref *models.NotificationPreference) error {
	if pref.QuietHoursStart != nil {
		if err := r.validateTime(pref.QuietHoursStart); err != nil {
			return fmt.Errorf("invalid quiet_hours_start: %w", err)
		}
	}

	if pref.QuietHoursEnd != nil {
		if err := r.validateTime(pref.QuietHoursEnd); err != nil {
			return fmt.Errorf("invalid quiet_hours_end: %w", err)
		}
	}

	return nil
}

// validateTime validates time format
func (r *notificationPreferenceRepository) validateTime(t *time.Time) error {
	if t == nil {
		return nil
	}

	// Ensure time is within valid range (00:00:00 - 23:59:59)
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23, got %d", hour)
	}

	if minute < 0 || minute > 59 {
		return fmt.Errorf("minute must be between 0 and 59, got %d", minute)
	}

	if second < 0 || second > 59 {
		return fmt.Errorf("second must be between 0 and 59, got %d", second)
	}

	return nil
}

// validateTimezone validates timezone string
func (r *notificationPreferenceRepository) validateTimezone(pref *models.NotificationPreference) error {
	if !pref.Timezone.Valid || pref.Timezone.String == "" {
		return nil
	}

	return r.validateTimezoneString(pref.Timezone.String)
}

// validateTimezoneString validates timezone string
func (r *notificationPreferenceRepository) validateTimezoneString(timezone interface{}) error {
	var tz string

	switch v := timezone.(type) {
	case string:
		tz = v
	default:
		return fmt.Errorf("timezone must be a string")
	}

	if tz == "" {
		return nil
	}

	// Validate timezone by attempting to load it
	_, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone %s: %w", tz, err)
	}

	return nil
}

// Upsert creates or updates notification preferences
func (r *notificationPreferenceRepository) Upsert(ctx context.Context, pref *models.NotificationPreference) error {
	// Validate quiet hours format if provided
	if err := r.validateQuietHours(pref); err != nil {
		return fmt.Errorf("invalid quiet hours: %w", err)
	}

	// Validate timezone if provided
	if err := r.validateTimezone(pref); err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	// Check if preference exists
	existing, err := r.GetByUserID(ctx, pref.UserID)
	if err != nil {
		// Create new preference if not found
		return r.Create(ctx, pref)
	}

	// Update existing preference
	updates := map[string]interface{}{
		"push_enabled":          pref.PushEnabled,
		"in_app_enabled":        pref.InAppEnabled,
		"email_enabled":         pref.EmailEnabled,
		"store_streams_enabled": pref.StoreStreamsEnabled,
		"video_streams_enabled": pref.VideoStreamsEnabled,
		"quiet_hours_start":     pref.QuietHoursStart,
		"quiet_hours_end":       pref.QuietHoursEnd,
		"timezone":              pref.Timezone,
		"updated_at":            time.Now(),
	}

	return r.Update(ctx, existing.UserID, updates)
}