package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CallHistory represents a historical record of calls for a user
type CallHistory struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CallSessionID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"call_session_id" validate:"required"`
	UserID                 uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"` // The user this history record belongs to
	OtherParticipantID     uuid.UUID      `gorm:"type:uuid;not null" json:"other_participant_id" validate:"required"`
	OtherParticipantName   *string        `gorm:"type:varchar(255)" json:"other_participant_name,omitempty"`
	OtherParticipantAvatar *string        `gorm:"type:text" json:"other_participant_avatar,omitempty"`
	CallType               CallType       `gorm:"type:varchar(10);not null" json:"call_type" validate:"required,oneof=voice video"`
	Duration               int            `gorm:"not null;default:0" json:"duration"` // Duration in seconds
	InitiatedByMe          bool           `gorm:"not null" json:"initiated_by_me"`
	CallStatus             CallStatus     `gorm:"type:varchar(20);not null" json:"call_status"`
	CreatedAt              time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	CallSession CallSession `gorm:"foreignKey:CallSessionID" json:"-"`
}

// CallHistoryFilter represents filters for call history queries
type CallHistoryFilter struct {
	UserID     uuid.UUID   `json:"user_id"`
	CallType   *CallType   `json:"call_type,omitempty"`
	CallStatus *CallStatus `json:"call_status,omitempty"`
	StartDate  *time.Time  `json:"start_date,omitempty"`
	EndDate    *time.Time  `json:"end_date,omitempty"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (h *CallHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// CreateFromCallSession creates call history records for all participants
func CreateFromCallSession(callSession *CallSession) []*CallHistory {
	var historyRecords []*CallHistory

	for _, participant := range callSession.Participants {
		// Find the other participant(s)
		var otherParticipantID uuid.UUID
		for _, otherParticipant := range callSession.Participants {
			if otherParticipant.UserID != participant.UserID {
				otherParticipantID = otherParticipant.UserID
				break
			}
		}

		// Create history record for this participant
		history := &CallHistory{
			CallSessionID:      callSession.ID,
			UserID:             participant.UserID,
			OtherParticipantID: otherParticipantID,
			CallType:           callSession.Type,
			Duration:           0, // Will be updated when call ends
			InitiatedByMe:      callSession.InitiatedBy == participant.UserID,
			CallStatus:         callSession.Status,
			CreatedAt:          callSession.StartedAt,
		}

		// Set duration if call has ended
		if callSession.Duration != nil {
			history.Duration = *callSession.Duration
		}

		historyRecords = append(historyRecords, history)
	}

	return historyRecords
}

// UpdateFromCallSession updates the history record when call session changes
func (h *CallHistory) UpdateFromCallSession(callSession *CallSession) {
	h.CallStatus = callSession.Status
	if callSession.Duration != nil {
		h.Duration = *callSession.Duration
	}
	h.UpdatedAt = time.Now()
}

// GetFormattedDuration returns a human-readable duration string
func (h *CallHistory) GetFormattedDuration() string {
	if h.Duration == 0 {
		return "0s"
	}

	duration := time.Duration(h.Duration) * time.Second

	if duration < time.Minute {
		return duration.String()
	}

	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	if duration < time.Hour {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	hours := minutes / 60
	minutes = minutes % 60
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

// IsSuccessful returns true if the call was successfully completed
func (h *CallHistory) IsSuccessful() bool {
	return h.CallStatus == CallStatusEnded && h.Duration > 0
}

// IsMissed returns true if this was a missed call for the user
func (h *CallHistory) IsMissed() bool {
	return !h.InitiatedByMe && h.CallStatus == CallStatusFailed && h.Duration == 0
}

// IsOutgoing returns true if this was an outgoing call
func (h *CallHistory) IsOutgoing() bool {
	return h.InitiatedByMe
}

// IsIncoming returns true if this was an incoming call
func (h *CallHistory) IsIncoming() bool {
	return !h.InitiatedByMe
}

// GetCallDirection returns a string describing the call direction
func (h *CallHistory) GetCallDirection() string {
	if h.InitiatedByMe {
		return "outgoing"
	}
	return "incoming"
}

// GetCallOutcome returns a string describing the call outcome
func (h *CallHistory) GetCallOutcome() string {
	switch h.CallStatus {
	case CallStatusEnded:
		if h.Duration > 0 {
			return "completed"
		}
		return "ended"
	case CallStatusFailed:
		if h.IsMissed() {
			return "missed"
		}
		return "failed"
	default:
		return string(h.CallStatus)
	}
}
