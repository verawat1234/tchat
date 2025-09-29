package models

import (
	"time"

	"github.com/google/uuid"
)

// UserPresence represents a user's availability status for calls
type UserPresence struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primary_key" json:"user_id" validate:"required"`
	Status    PresenceStatus `gorm:"type:varchar(10);not null;default:'offline'" json:"status"`
	LastSeen  time.Time      `gorm:"not null;default:now()" json:"last_seen"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
}

// PresenceStatus represents the user's availability status
type PresenceStatus string

const (
	PresenceStatusOnline  PresenceStatus = "online"
	PresenceStatusBusy    PresenceStatus = "busy"
	PresenceStatusOffline PresenceStatus = "offline"
	PresenceStatusInCall  PresenceStatus = "in_call"
)

// IsValid validates the presence status
func (ps PresenceStatus) IsValid() bool {
	switch ps {
	case PresenceStatusOnline, PresenceStatusBusy, PresenceStatusOffline, PresenceStatusInCall:
		return true
	default:
		return false
	}
}

// IsAvailableForCalls returns true if the user can receive calls
func (ps PresenceStatus) IsAvailableForCalls() bool {
	return ps == PresenceStatusOnline
}

// UpdateStatus updates the user's presence status
func (p *UserPresence) UpdateStatus(status PresenceStatus) error {
	if !status.IsValid() {
		return ErrInvalidPresenceStatus
	}

	p.Status = status
	p.UpdatedAt = time.Now()

	// Update LastSeen if coming online
	if status == PresenceStatusOnline {
		p.LastSeen = time.Now()
	}

	return nil
}

// SetOnline sets the user as online and available
func (p *UserPresence) SetOnline() {
	p.Status = PresenceStatusOnline
	p.LastSeen = time.Now()
	p.UpdatedAt = time.Now()
}

// SetOffline sets the user as offline
func (p *UserPresence) SetOffline() {
	p.Status = PresenceStatusOffline
	p.UpdatedAt = time.Now()
}

// SetBusy sets the user as busy (not available for calls)
func (p *UserPresence) SetBusy() {
	p.Status = PresenceStatusBusy
	p.UpdatedAt = time.Now()
}

// SetInCall sets the user as currently in a call
func (p *UserPresence) SetInCall() {
	p.Status = PresenceStatusInCall
	p.UpdatedAt = time.Now()
}

// IsOnline returns true if the user is online and available
func (p *UserPresence) IsOnline() bool {
	return p.Status == PresenceStatusOnline
}

// IsBusy returns true if the user is busy
func (p *UserPresence) IsBusy() bool {
	return p.Status == PresenceStatusBusy
}

// IsInCall returns true if the user is currently in a call
func (p *UserPresence) IsInCall() bool {
	return p.Status == PresenceStatusInCall
}

// IsAvailable returns true if the user can receive calls
func (p *UserPresence) IsAvailable() bool {
	return p.Status.IsAvailableForCalls()
}

// GetAvailabilityInfo returns structured availability information
func (p *UserPresence) GetAvailabilityInfo() map[string]interface{} {
	return map[string]interface{}{
		"user_id":   p.UserID,
		"available": p.IsAvailable(),
		"status":    string(p.Status),
		"last_seen": p.LastSeen,
	}
}

// IsRecentlyActive returns true if user was active within the specified duration
func (p *UserPresence) IsRecentlyActive(duration time.Duration) bool {
	return time.Since(p.LastSeen) <= duration
}

// GetOfflineDuration returns how long the user has been offline
func (p *UserPresence) GetOfflineDuration() time.Duration {
	if p.Status != PresenceStatusOffline {
		return 0
	}
	return time.Since(p.LastSeen)
}
