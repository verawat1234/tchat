package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CallParticipant represents a participant in a call session
type CallParticipant struct {
	ID                uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CallSessionID     uuid.UUID         `gorm:"type:uuid;not null;index" json:"call_session_id" validate:"required"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	Role              ParticipantRole   `gorm:"type:varchar(10);not null" json:"role" validate:"required,oneof=caller callee"`
	JoinedAt          time.Time         `gorm:"not null;default:now()" json:"joined_at"`
	LeftAt            *time.Time        `gorm:"" json:"left_at,omitempty"`
	AudioEnabled      bool              `gorm:"default:true" json:"audio_enabled"`
	VideoEnabled      bool              `gorm:"default:false" json:"video_enabled"`
	ConnectionQuality ConnectionQuality `gorm:"type:varchar(10);default:'good'" json:"connection_quality"`
	CreatedAt         time.Time         `gorm:"default:now()" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"default:now()" json:"updated_at"`
	DeletedAt         gorm.DeletedAt    `gorm:"index" json:"-"`

	// Relationships
	CallSession CallSession `gorm:"foreignKey:CallSessionID" json:"-"`
}

// ParticipantRole represents the role of a participant in a call
type ParticipantRole string

const (
	ParticipantRoleCaller ParticipantRole = "caller"
	ParticipantRoleCallee ParticipantRole = "callee"
)

// ConnectionQuality represents the quality of the participant's connection
type ConnectionQuality string

const (
	ConnectionQualityExcellent ConnectionQuality = "excellent"
	ConnectionQualityGood      ConnectionQuality = "good"
	ConnectionQualityFair      ConnectionQuality = "fair"
	ConnectionQualityPoor      ConnectionQuality = "poor"
)

// IsValid validates the participant role
func (pr ParticipantRole) IsValid() bool {
	switch pr {
	case ParticipantRoleCaller, ParticipantRoleCallee:
		return true
	default:
		return false
	}
}

// IsValid validates the connection quality
func (cq ConnectionQuality) IsValid() bool {
	switch cq {
	case ConnectionQualityExcellent, ConnectionQualityGood, ConnectionQualityFair, ConnectionQualityPoor:
		return true
	default:
		return false
	}
}

// BeforeCreate is a GORM hook that runs before creating a record
func (p *CallParticipant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// IsActive returns true if the participant is currently in the call
func (p *CallParticipant) IsActive() bool {
	return p.LeftAt == nil
}

// CalculateDuration calculates how long the participant was in the call
func (p *CallParticipant) CalculateDuration() time.Duration {
	endTime := time.Now()
	if p.LeftAt != nil {
		endTime = *p.LeftAt
	}
	return endTime.Sub(p.JoinedAt)
}

// Leave sets the participant as having left the call
func (p *CallParticipant) Leave() {
	now := time.Now()
	p.LeftAt = &now
}

// ToggleAudio toggles the participant's audio state
func (p *CallParticipant) ToggleAudio() bool {
	p.AudioEnabled = !p.AudioEnabled
	return p.AudioEnabled
}

// ToggleVideo toggles the participant's video state
func (p *CallParticipant) ToggleVideo() bool {
	p.VideoEnabled = !p.VideoEnabled
	return p.VideoEnabled
}

// SetConnectionQuality updates the participant's connection quality
func (p *CallParticipant) SetConnectionQuality(quality ConnectionQuality) error {
	if !quality.IsValid() {
		return ErrInvalidConnectionQuality
	}
	p.ConnectionQuality = quality
	return nil
}

// GetQualityScore returns a numeric score for connection quality
func (p *CallParticipant) GetQualityScore() int {
	switch p.ConnectionQuality {
	case ConnectionQualityExcellent:
		return 4
	case ConnectionQualityGood:
		return 3
	case ConnectionQualityFair:
		return 2
	case ConnectionQualityPoor:
		return 1
	default:
		return 0
	}
}
