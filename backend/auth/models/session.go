package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Session represents a user authentication session
type Session struct {
	ID               uuid.UUID              `json:"id" db:"id"`
	UserID           uuid.UUID              `json:"user_id" db:"user_id"`
	DeviceID         string                 `json:"device_id" db:"device_id"`
	AccessToken      string                 `json:"access_token" db:"access_token_hash"`
	RefreshToken     string                 `json:"refresh_token" db:"refresh_token_hash"`
	ExpiresAt        time.Time              `json:"expires_at" db:"expires_at"`
	RefreshExpiresAt time.Time              `json:"refresh_expires_at" db:"refresh_expires_at"`
	IsActive         bool                   `json:"is_active" db:"is_active"`
	IPAddress        string                 `json:"ip_address" db:"ip_address"`
	UserAgent        string                 `json:"user_agent" db:"user_agent"`
	DeviceInfo       map[string]interface{} `json:"device_info" db:"device_info"`
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
	LastActiveAt     time.Time              `json:"last_active_at" db:"last_active_at"`
	LastUsed         time.Time              `json:"last_used" db:"last_used"`
	RevokedAt        *time.Time             `json:"revoked_at,omitempty" db:"revoked_at"`
	Status           SessionStatus          `json:"status" db:"status"`
}

// SessionStatus represents the current state of a session
type SessionStatus string

const (
	SessionStatusActive      SessionStatus = "active"
	SessionStatusExpired     SessionStatus = "expired"
	SessionStatusRevoked     SessionStatus = "revoked"
	SessionStatusSuspended   SessionStatus = "suspended"
	SessionStatusTerminated  SessionStatus = "terminated"
)

// Session configuration constants
const (
	AccessTokenExpiry  = 15 * time.Minute    // 15 minutes for access tokens
	RefreshTokenExpiry = 30 * 24 * time.Hour // 30 days for refresh tokens
	MaxSessionsPerUser = 5                   // Maximum concurrent sessions per user
	TokenLength        = 32                  // Token byte length
)

// SessionValidationError represents session validation errors
type SessionValidationError struct {
	Field   string
	Message string
}

func (e SessionValidationError) Error() string {
	return fmt.Sprintf("session validation error - %s: %s", e.Field, e.Message)
}

// IsValid validates if the session status is supported
func (s SessionStatus) IsValid() bool {
	validStatuses := []SessionStatus{
		SessionStatusActive,
		SessionStatusExpired,
		SessionStatusRevoked,
		SessionStatusSuspended,
		SessionStatusTerminated,
	}
	for _, valid := range validStatuses {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of SessionStatus
func (s SessionStatus) String() string {
	return string(s)
}

// Validate performs comprehensive validation on the Session model
func (s *Session) Validate() error {
	var errs []string

	// User ID validation
	if s.UserID == uuid.Nil {
		errs = append(errs, "user_id is required")
	}

	// Device ID validation
	if strings.TrimSpace(s.DeviceID) == "" {
		errs = append(errs, "device_id is required")
	}
	if len(s.DeviceID) > 255 {
		errs = append(errs, "device_id must not exceed 255 characters")
	}

	// Token validation
	if s.AccessToken == "" {
		errs = append(errs, "access_token is required")
	}
	if s.RefreshToken == "" {
		errs = append(errs, "refresh_token is required")
	}

	// Expiry validation
	if s.ExpiresAt.IsZero() {
		errs = append(errs, "expires_at is required")
	}

	// Status validation
	if !s.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid session status: %s", s.Status))
	}

	// IP Address validation if provided
	if s.IPAddress != "" {
		if err := s.validateIPAddress(s.IPAddress); err != nil {
			errs = append(errs, fmt.Sprintf("invalid IP address: %v", err))
		}
	}

	// User Agent validation if provided
	if len(s.UserAgent) > 512 {
		errs = append(errs, "user_agent must not exceed 512 characters")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateIPAddress validates IP address format (IPv4 or IPv6)
func (s *Session) validateIPAddress(ip string) error {
	// Basic IP validation - could be enhanced with net.ParseIP
	if len(ip) < 7 || len(ip) > 45 {
		return fmt.Errorf("invalid IP address length")
	}
	// Additional validation logic can be added here
	return nil
}

// BeforeCreate sets up the session before database creation
func (s *Session) BeforeCreate() error {
	// Generate UUID if not set
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	s.CreatedAt = now
	s.LastUsed = now

	// Set default status
	if s.Status == "" {
		s.Status = SessionStatusActive
	}

	// Set active flag based on status
	s.IsActive = (s.Status == SessionStatusActive)

	// Set default expiry if not provided
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = now.Add(AccessTokenExpiry)
	}

	// Generate tokens if not provided
	if s.AccessToken == "" {
		token, err := s.generateSecureToken()
		if err != nil {
			return fmt.Errorf("failed to generate access token: %v", err)
		}
		s.AccessToken = token
	}

	if s.RefreshToken == "" {
		token, err := s.generateSecureToken()
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %v", err)
		}
		s.RefreshToken = token
	}

	// Validate before creation
	return s.Validate()
}

// BeforeUpdate sets up the session before database update
func (s *Session) BeforeUpdate() error {
	// Update last used timestamp
	s.LastUsed = time.Now().UTC()

	// Update active flag based on status
	s.IsActive = (s.Status == SessionStatusActive) && !s.IsExpired()

	// Validate before update
	return s.Validate()
}

// generateSecureToken generates a cryptographically secure random token
func (s *Session) generateSecureToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// IsValid checks if the session is currently valid and active
func (s *Session) IsValid() bool {
	return s.IsActive &&
		   s.Status == SessionStatusActive &&
		   !s.IsExpired() &&
		   s.RevokedAt == nil
}

// CanRefresh checks if the session can be refreshed
func (s *Session) CanRefresh() bool {
	// Can refresh if not revoked and within refresh token expiry
	if s.RevokedAt != nil {
		return false
	}

	refreshExpiry := s.CreatedAt.Add(RefreshTokenExpiry)
	return time.Now().UTC().Before(refreshExpiry)
}

// Refresh extends the session with new tokens and expiry
func (s *Session) Refresh() error {
	if !s.CanRefresh() {
		return errors.New("session cannot be refreshed")
	}

	// Generate new access token
	newAccessToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate new access token: %v", err)
	}

	// Update session
	now := time.Now().UTC()
	s.AccessToken = newAccessToken
	s.ExpiresAt = now.Add(AccessTokenExpiry)
	s.LastUsed = now
	s.Status = SessionStatusActive
	s.IsActive = true

	return s.Validate()
}

// Revoke marks the session as revoked
func (s *Session) Revoke(reason string) error {
	now := time.Now().UTC()
	s.Status = SessionStatusRevoked
	s.IsActive = false
	s.RevokedAt = &now
	s.LastUsed = now

	return s.Validate()
}

// Expire marks the session as expired
func (s *Session) Expire() error {
	now := time.Now().UTC()
	s.Status = SessionStatusExpired
	s.IsActive = false
	s.LastUsed = now

	return s.Validate()
}

// Suspend temporarily suspends the session
func (s *Session) Suspend() error {
	s.Status = SessionStatusSuspended
	s.IsActive = false
	s.LastUsed = time.Now().UTC()

	return s.Validate()
}

// Reactivate reactivates a suspended session
func (s *Session) Reactivate() error {
	if s.Status != SessionStatusSuspended {
		return errors.New("can only reactivate suspended sessions")
	}

	if s.IsExpired() {
		return errors.New("cannot reactivate expired session")
	}

	if s.RevokedAt != nil {
		return errors.New("cannot reactivate revoked session")
	}

	s.Status = SessionStatusActive
	s.IsActive = true
	s.LastUsed = time.Now().UTC()

	return s.Validate()
}

// UpdateActivity updates the last used timestamp
func (s *Session) UpdateActivity() {
	s.LastUsed = time.Now().UTC()
}

// CanTransitionToStatus checks if session can transition to given status
func (s *Session) CanTransitionToStatus(newStatus SessionStatus) bool {
	switch s.Status {
	case SessionStatusActive:
		return newStatus == SessionStatusExpired ||
			   newStatus == SessionStatusRevoked ||
			   newStatus == SessionStatusSuspended ||
			   newStatus == SessionStatusTerminated
	case SessionStatusSuspended:
		return newStatus == SessionStatusActive ||
			   newStatus == SessionStatusExpired ||
			   newStatus == SessionStatusRevoked ||
			   newStatus == SessionStatusTerminated
	case SessionStatusExpired, SessionStatusRevoked, SessionStatusTerminated:
		// Terminal states - no transitions allowed
		return false
	default:
		return true // Allow any transition for unknown statuses
	}
}

// GetRemainingTime returns the remaining time before expiry
func (s *Session) GetRemainingTime() time.Duration {
	if s.IsExpired() {
		return 0
	}
	return time.Until(s.ExpiresAt)
}

// GetRefreshRemainingTime returns the remaining time before refresh expiry
func (s *Session) GetRefreshRemainingTime() time.Duration {
	refreshExpiry := s.CreatedAt.Add(RefreshTokenExpiry)
	if time.Now().UTC().After(refreshExpiry) {
		return 0
	}
	return time.Until(refreshExpiry)
}

// ToTokenResponse returns a response suitable for token endpoints
func (s *Session) ToTokenResponse() map[string]interface{} {
	return map[string]interface{}{
		"access_token":  s.AccessToken,
		"refresh_token": s.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(s.GetRemainingTime().Seconds()),
		"expires_at":    s.ExpiresAt,
	}
}

// ToSessionInfo returns session information for API responses
func (s *Session) ToSessionInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":         s.ID,
		"device_id":  s.DeviceID,
		"is_active":  s.IsActive,
		"status":     s.Status,
		"created_at": s.CreatedAt,
		"last_used":  s.LastUsed,
		"expires_at": s.ExpiresAt,
		"ip_address": s.IPAddress,
		"user_agent": s.UserAgent,
	}
}

// SessionCreateRequest represents a request to create a new session
type SessionCreateRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required,max=255"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty,max=512"`
}

// ToSession converts a create request to a Session model
func (req *SessionCreateRequest) ToSession() *Session {
	return &Session{
		UserID:    req.UserID,
		DeviceID:  req.DeviceID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
	}
}

// SessionManager provides session management utilities
type SessionManager struct {
	// Add dependencies like database, cache, etc.
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// CreateSession creates a new session with proper validation
func (sm *SessionManager) CreateSession(req *SessionCreateRequest) (*Session, error) {
	session := req.ToSession()

	if err := session.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("session creation failed: %v", err)
	}

	return session, nil
}

// ValidateToken validates a session token
func (sm *SessionManager) ValidateToken(token string) error {
	if len(token) != TokenLength*2 { // hex encoded
		return errors.New("invalid token length")
	}

	// Additional token validation logic
	return nil
}

// CleanupExpiredSessions returns filter criteria for cleanup operations
func (sm *SessionManager) CleanupExpiredSessions() map[string]interface{} {
	cutoff := time.Now().UTC().Add(-24 * time.Hour) // Cleanup after 24h
	return map[string]interface{}{
		"status":    []SessionStatus{SessionStatusExpired, SessionStatusRevoked},
		"before":    cutoff,
	}
}