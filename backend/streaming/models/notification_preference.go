package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationPreference represents user notification preferences for streaming content
type NotificationPreference struct {
	UserID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"user_id"`
	PushEnabled         bool           `gorm:"default:true;index:idx_push_enabled" json:"push_enabled"`
	InAppEnabled        bool           `gorm:"default:true" json:"in_app_enabled"`
	EmailEnabled        bool           `gorm:"default:true;index:idx_email_enabled" json:"email_enabled"`
	StoreStreamsEnabled bool           `gorm:"default:true" json:"store_streams_enabled"`
	VideoStreamsEnabled bool           `gorm:"default:true" json:"video_streams_enabled"`
	QuietHoursStart     *time.Time     `gorm:"type:time" json:"quiet_hours_start,omitempty"`
	QuietHoursEnd       *time.Time     `gorm:"type:time" json:"quiet_hours_end,omitempty"`
	Timezone            sql.NullString `gorm:"type:varchar(50)" json:"timezone,omitempty"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for NotificationPreference model
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// BeforeCreate hook to set default values
func (np *NotificationPreference) BeforeCreate(tx *gorm.DB) error {
	if np.UserID == uuid.Nil {
		np.UserID = uuid.New()
	}
	return nil
}

// IsQuietHoursActive checks if current time falls within quiet hours
func (np *NotificationPreference) IsQuietHoursActive(currentTime time.Time) bool {
	if np.QuietHoursStart == nil || np.QuietHoursEnd == nil {
		return false
	}

	// Extract time components only
	start := np.QuietHoursStart.Hour()*60 + np.QuietHoursStart.Minute()
	end := np.QuietHoursEnd.Hour()*60 + np.QuietHoursEnd.Minute()
	current := currentTime.Hour()*60 + currentTime.Minute()

	// Handle overnight quiet hours (e.g., 22:00 to 06:00)
	if start > end {
		return current >= start || current <= end
	}

	// Standard range (e.g., 08:00 to 17:00)
	return current >= start && current <= end
}

// ShouldNotify determines if notification should be sent based on preferences
func (np *NotificationPreference) ShouldNotify(channelType string, currentTime time.Time) bool {
	// Check if quiet hours are active
	if np.IsQuietHoursActive(currentTime) {
		return false
	}

	// Check channel-specific preferences
	switch channelType {
	case "push":
		return np.PushEnabled
	case "in_app":
		return np.InAppEnabled
	case "email":
		return np.EmailEnabled
	default:
		return false
	}
}