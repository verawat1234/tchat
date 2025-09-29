package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CallSession represents a voice or video call session
type CallSession struct {
	ID            uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type          CallType          `gorm:"type:varchar(10);not null" json:"type" validate:"required,oneof=voice video"`
	Status        CallStatus        `gorm:"type:varchar(20);not null;default:'connecting'" json:"status"`
	InitiatedBy   uuid.UUID         `gorm:"type:uuid;not null" json:"initiated_by" validate:"required"`
	StartedAt     time.Time         `gorm:"not null;default:now()" json:"started_at"`
	EndedAt       *time.Time        `gorm:"" json:"ended_at,omitempty"`
	Duration      *int              `gorm:"" json:"duration,omitempty"` // Duration in seconds
	FailureReason *string           `gorm:"type:text" json:"failure_reason,omitempty"`
	Participants  []CallParticipant `gorm:"foreignKey:CallSessionID" json:"participants"`
	CreatedAt     time.Time         `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

// CallType represents the type of call
type CallType string

const (
	CallTypeVoice CallType = "voice"
	CallTypeVideo CallType = "video"
)

// CallStatus represents the current status of a call
type CallStatus string

const (
	CallStatusConnecting CallStatus = "connecting"
	CallStatusActive     CallStatus = "active"
	CallStatusEnded      CallStatus = "ended"
	CallStatusFailed     CallStatus = "failed"
)

// IsValid validates the call type
func (ct CallType) IsValid() bool {
	switch ct {
	case CallTypeVoice, CallTypeVideo:
		return true
	default:
		return false
	}
}

// IsValid validates the call status
func (cs CallStatus) IsValid() bool {
	switch cs {
	case CallStatusConnecting, CallStatusActive, CallStatusEnded, CallStatusFailed:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if status transition is valid
func (cs CallStatus) CanTransitionTo(newStatus CallStatus) bool {
	switch cs {
	case CallStatusConnecting:
		return newStatus == CallStatusActive || newStatus == CallStatusFailed
	case CallStatusActive:
		return newStatus == CallStatusEnded || newStatus == CallStatusFailed
	case CallStatusEnded, CallStatusFailed:
		return false // Terminal states
	default:
		return false
	}
}

// BeforeCreate is a GORM hook that runs before creating a record
func (c *CallSession) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CalculateDuration calculates and sets the duration if call has ended
func (c *CallSession) CalculateDuration() {
	if c.EndedAt != nil {
		duration := int(c.EndedAt.Sub(c.StartedAt).Seconds())
		c.Duration = &duration
	}
}

// IsActive returns true if the call is currently active
func (c *CallSession) IsActive() bool {
	return c.Status == CallStatusActive
}

// IsTerminal returns true if the call is in a terminal state
func (c *CallSession) IsTerminal() bool {
	return c.Status == CallStatusEnded || c.Status == CallStatusFailed
}

// GetParticipantByUserID returns the participant for a given user ID
func (c *CallSession) GetParticipantByUserID(userID uuid.UUID) *CallParticipant {
	for i := range c.Participants {
		if c.Participants[i].UserID == userID {
			return &c.Participants[i]
		}
	}
	return nil
}

// HasParticipant checks if a user is a participant in this call
func (c *CallSession) HasParticipant(userID uuid.UUID) bool {
	return c.GetParticipantByUserID(userID) != nil
}
