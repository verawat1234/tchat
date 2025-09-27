package models

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSON is a custom type for handling JSON data in GORM
type JSON map[string]interface{}

// Value implements driver.Valuer interface for GORM
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface for GORM
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	return json.Unmarshal(data, j)
}

// Session represents a user authentication session
// Maps to the user_sessions table in the database
type Session struct {
	ID            uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key"`
	UserID        uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index"`
	SessionToken  string    `json:"session_token" gorm:"column:session_token;type:text;not null;unique"`
	RefreshToken  string    `json:"refresh_token" gorm:"column:refresh_token;type:text;not null;unique"`
	DeviceID      string    `json:"device_id" gorm:"column:device_id;type:text"`
	DeviceType    string    `json:"device_type" gorm:"column:device_type;type:varchar(50)"`
	IPAddress     string    `json:"ip_address" gorm:"column:ip_address;type:inet"`
	UserAgent     string    `json:"user_agent" gorm:"column:user_agent;type:text"`
	ExpiresAt     time.Time `json:"expires_at" gorm:"column:expires_at;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	LastUsedAt    time.Time `json:"last_used_at" gorm:"column:last_used_at;autoUpdateTime"`

	// Additional fields for JWT tokens (not stored in DB)
	AccessToken      string                 `json:"access_token,omitempty" gorm:"-"`
	RefreshExpiresAt time.Time              `json:"refresh_expires_at,omitempty" gorm:"-"`
	IsActive         bool                   `json:"is_active,omitempty" gorm:"-"`
	DeviceInfo       JSON `json:"device_info,omitempty" gorm:"-"`
	Metadata         JSON `json:"metadata,omitempty" gorm:"-"`
	UpdatedAt        time.Time              `json:"updated_at,omitempty" gorm:"-"`
	LastActiveAt     time.Time              `json:"last_active_at,omitempty" gorm:"-"`
	RevokedAt        *time.Time             `json:"revoked_at,omitempty" gorm:"-"`
	Status           SessionStatus          `json:"status,omitempty" gorm:"column:status;size:20;default:'active'"`
}

// TableName returns the table name for the Session model
func (Session) TableName() string {
	return "user_sessions"
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

	// RefreshToken validation (main token stored in DB)
	if s.RefreshToken == "" {
		errs = append(errs, "refresh_token is required")
	}

	// Expiry validation
	if s.ExpiresAt.IsZero() {
		errs = append(errs, "expires_at is required")
	}

	// Device ID validation (optional in DB schema, reasonable upper limit)
	if len(s.DeviceID) > 1000 {
		errs = append(errs, "device_id must not exceed 1000 characters")
	}

	// Device Type validation (optional in DB schema)
	if len(s.DeviceType) > 50 {
		errs = append(errs, "device_type must not exceed 50 characters")
	}

	// IP Address validation if provided
	if s.IPAddress != "" {
		if err := s.validateIPAddress(s.IPAddress); err != nil {
			errs = append(errs, fmt.Sprintf("invalid IP address: %v", err))
		}
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
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	s.CreatedAt = now
	s.LastUsedAt = now

	// Set default expiry if not provided (30 days for refresh token)
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = now.Add(RefreshTokenExpiry)
	}

	// Generate session token if not provided
	if s.SessionToken == "" {
		token, err := s.generateSecureToken()
		if err != nil {
			return fmt.Errorf("failed to generate session token: %v", err)
		}
		s.SessionToken = token
	}

	// Generate refresh token if not provided
	if s.RefreshToken == "" {
		token, err := s.generateSecureToken()
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %v", err)
		}
		s.RefreshToken = token
	}

	// Set default device type if not provided
	if s.DeviceType == "" {
		s.DeviceType = "unknown"
	}

	// Validate before creation
	return s.Validate()
}

// BeforeUpdate sets up the session before database update
func (s *Session) BeforeUpdate(tx *gorm.DB) error {
	// Update last used timestamp
	s.LastUsedAt = time.Now().UTC()

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

	// Generate new session token
	newSessionToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate new session token: %v", err)
	}

	// Generate new access token
	newAccessToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate new access token: %v", err)
	}

	// Update session
	now := time.Now().UTC()
	s.SessionToken = newSessionToken
	s.AccessToken = newAccessToken
	s.ExpiresAt = now.Add(AccessTokenExpiry)
	s.LastUsedAt = now
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
	s.LastUsedAt = now

	return s.Validate()
}

// Expire marks the session as expired
func (s *Session) Expire() error {
	now := time.Now().UTC()
	s.Status = SessionStatusExpired
	s.IsActive = false
	s.LastUsedAt = now

	return s.Validate()
}

// Suspend temporarily suspends the session
func (s *Session) Suspend() error {
	s.Status = SessionStatusSuspended
	s.IsActive = false
	s.LastUsedAt = time.Now().UTC()

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
	s.LastUsedAt = time.Now().UTC()

	return s.Validate()
}

// UpdateActivity updates the last used timestamp
func (s *Session) UpdateActivity() {
	s.LastUsedAt = time.Now().UTC()
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
		"last_used":  s.LastUsedAt,
		"expires_at": s.ExpiresAt,
		"ip_address": s.IPAddress,
		"user_agent": s.UserAgent,
	}
}

// SessionCreateRequest represents a request to create a new session
type SessionCreateRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	DeviceID  string    `json:"device_id" validate:"required,max=1000"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty,max=512"`
}

// ToSession converts a create request to a Session model
func (req *SessionCreateRequest) ToSession() *Session {
	var ipAddress, userAgent string
	if req.IPAddress != nil {
		ipAddress = *req.IPAddress
	}
	if req.UserAgent != nil {
		userAgent = *req.UserAgent
	}

	return &Session{
		UserID:    req.UserID,
		DeviceID:  req.DeviceID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
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

	if err := session.BeforeCreate(nil); err != nil {
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