package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type LiveStream struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BroadcasterID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_broadcaster_id"`
	StreamType          string         `gorm:"type:varchar(20);not null;check:stream_type IN ('store', 'video');index:idx_stream_type"`
	Title               string         `gorm:"type:varchar(200);not null"`
	Description         sql.NullString `gorm:"type:text"`
	Status              string         `gorm:"type:varchar(20);not null;default:'scheduled';check:status IN ('scheduled', 'live', 'ended', 'terminated');index:idx_status"`
	BroadcasterKYCTier  int            `gorm:"not null;check:broadcaster_kyc_tier >= 0 AND broadcaster_kyc_tier <= 3"`
	PrivacySetting      string         `gorm:"type:varchar(20);not null;default:'public';check:privacy_setting IN ('public', 'followers_only', 'private')"`
	ScheduledStartTime  *time.Time     `gorm:"index:idx_scheduled_start_time"`
	ActualStartTime     *time.Time
	EndTime             *time.Time
	StreamKey           string         `gorm:"type:varchar(100);unique;not null"`
	WebRTCSessionID     sql.NullString `gorm:"type:varchar(100)"`
	PrimaryServerID     sql.NullString `gorm:"type:varchar(50)"`
	RecordingURL        sql.NullString `gorm:"type:varchar(500)"`
	RecordingExpiryDate *time.Time     `gorm:"index:idx_recording_expiry_date"`
	ViewerCount         int            `gorm:"default:0"`
	PeakViewerCount     int            `gorm:"default:0"`
	MaxCapacity         int            `gorm:"default:50000"`
	ThumbnailURL        sql.NullString `gorm:"type:varchar(500)"`
	Language            string         `gorm:"type:varchar(10);default:'en'"`
	Tags                pq.StringArray `gorm:"type:text[]"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (LiveStream) TableName() string {
	return "live_streams"
}