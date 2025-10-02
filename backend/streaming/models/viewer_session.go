package models

import (
	"database/sql"
	"net"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ViewerSession represents a viewer's session for a live stream
type ViewerSession struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StreamID             uuid.UUID      `gorm:"type:uuid;not null;index:idx_stream_id" json:"stream_id"`
	UserID               *uuid.UUID     `gorm:"type:uuid;index:idx_user_id" json:"user_id,omitempty"` // NULL for anonymous viewers
	JoinedAt             time.Time      `gorm:"not null;index:idx_joined_at" json:"joined_at"`
	LeftAt               *time.Time     `gorm:"" json:"left_at,omitempty"`
	WatchDurationSeconds *int           `gorm:"->" json:"watch_duration_seconds,omitempty"` // GENERATED column in PostgreSQL
	AverageQualityLayer  sql.NullString `gorm:"type:varchar(20)" json:"average_quality_layer,omitempty"`
	RebufferCount        int            `gorm:"default:0" json:"rebuffer_count"`
	UserAgent            sql.NullString `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress            *net.IP        `gorm:"type:inet" json:"ip_address,omitempty"`
	CountryCode          sql.NullString `gorm:"type:varchar(5)" json:"country_code,omitempty"`
}

// TableName specifies the table name for GORM
func (ViewerSession) TableName() string {
	return "viewer_sessions"
}

// BeforeCreate hook to generate UUID if not provided
func (v *ViewerSession) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}