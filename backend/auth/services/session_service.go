package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/auth/models"
	sharedModels "tchat.dev/shared/models"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*models.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
}

type SessionService struct {
	sessionRepo    SessionRepository
	eventPublisher EventPublisher
	jwtService     *JWTService
	db             *gorm.DB
}

type SessionConfig struct {
	AccessTokenExpiry  time.Duration `json:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
	MaxSessionsPerUser int           `json:"max_sessions_per_user"`
	TokenLength        int           `json:"token_length"`
	AllowConcurrentSessions bool     `json:"allow_concurrent_sessions"`
}

func NewSessionService(sessionRepo SessionRepository, eventPublisher EventPublisher, jwtService *JWTService, db *gorm.DB) *SessionService {
	return &SessionService{
		sessionRepo:    sessionRepo,
		eventPublisher: eventPublisher,
		jwtService:     jwtService,
		db:             db,
	}
}

func (ss *SessionService) CreateSession(ctx context.Context, req *CreateSessionRequest) (*models.Session, error) {
	// Validate request
	if err := ss.validateCreateSessionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user for JWT token generation
	var user sharedModels.User
	if err := ss.db.WithContext(ctx).First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user for token generation: %w", err)
	}

	// Generate device ID from device info
	deviceID := "unknown"
	if req.DeviceInfo != nil {
		if deviceIDValue, exists := req.DeviceInfo["device_id"]; exists {
			if deviceIDStr, ok := deviceIDValue.(string); ok && deviceIDStr != "" {
				deviceID = deviceIDStr
			}
		}
	}

	// Generate session ID first
	sessionID := uuid.New()

	// Generate JWT tokens using JWT service
	tokenPair, err := ss.jwtService.GenerateTokenPair(ctx, &user, sessionID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT tokens: %w", err)
	}

	// Calculate expiry times
	now := time.Now()

	// Extract device type from device info
	deviceType := "unknown"
	if req.DeviceInfo != nil {
		if deviceTypeValue, exists := req.DeviceInfo["device_type"]; exists {
			if deviceTypeStr, ok := deviceTypeValue.(string); ok && deviceTypeStr != "" {
				deviceType = deviceTypeStr
			}
		}
	}

	// Create session - only fields that map to user_sessions table
	session := &models.Session{
		ID:           sessionID,
		UserID:       req.UserID,
		DeviceID:     deviceID,
		DeviceType:   deviceType,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    req.UserAgent,
		IPAddress:    req.IPAddress,
		ExpiresAt:    tokenPair.RefreshExpiresAt, // This is the refresh token expiry in the DB
		CreatedAt:    now,
		LastUsedAt:   now,

		// Additional fields for response (not stored in DB)
		AccessToken:      tokenPair.AccessToken,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt,
		IsActive:         true,
		DeviceInfo:       req.DeviceInfo,
		Metadata:         req.Metadata,
		Status:           models.SessionStatusActive,
	}

	// Check for existing active sessions
	existingSessions, err := ss.sessionRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing sessions: %w", err)
	}

	// Handle max sessions per user
	activeSessions := 0
	for _, existingSession := range existingSessions {
		if existingSession.Status == models.SessionStatusActive {
			activeSessions++
		}
	}

	maxSessions := 5 // Default max sessions per user
	if activeSessions >= maxSessions {
		// Terminate oldest session
		oldestSession := existingSessions[0]
		for _, s := range existingSessions {
			if s.CreatedAt.Before(oldestSession.CreatedAt) && s.Status == models.SessionStatusActive {
				oldestSession = s
			}
		}

		if err := ss.TerminateSession(ctx, oldestSession.ID, "max_sessions_exceeded"); err != nil {
			fmt.Printf("Failed to terminate oldest session: %v\n", err)
		}
	}

	// Save session
	if err := ss.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Publish session created event
	if err := ss.publishSessionEvent(ctx, sharedModels.EventTypeUserSessionCreated, session.UserID, map[string]interface{}{
		"session_id": session.ID,
		"ip_address": session.IPAddress,
		"user_agent": session.UserAgent,
		"expires_at": session.ExpiresAt,
	}); err != nil {
		fmt.Printf("Failed to publish session created event: %v\n", err)
	}

	return session, nil
}

func (ss *SessionService) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	if sessionID == uuid.Nil {
		return nil, fmt.Errorf("session ID is required")
	}

	session, err := ss.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (ss *SessionService) GetSessionByAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	session, err := ss.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (ss *SessionService) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	session, err := ss.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (ss *SessionService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	sessions, err := ss.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}

func (ss *SessionService) RotateTokens(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	// Get existing session
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is active
	if session.Status != models.SessionStatusActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Get user for JWT token generation
	var user sharedModels.User
	if err := ss.db.WithContext(ctx).First(&user, session.UserID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user for token refresh: %w", err)
	}

	// Generate new JWT tokens using JWT service
	tokenPair, err := ss.jwtService.GenerateTokenPair(ctx, &user, session.ID, session.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new JWT tokens: %w", err)
	}

	// Update session with new tokens and extended expiry
	now := time.Now()
	session.AccessToken = tokenPair.AccessToken
	session.RefreshToken = tokenPair.RefreshToken
	session.ExpiresAt = tokenPair.RefreshExpiresAt // Use JWT service expiry times
	session.RefreshExpiresAt = tokenPair.RefreshExpiresAt
	session.UpdatedAt = now
	session.LastActiveAt = now

	// Save updated session
	if err := ss.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}

func (ss *SessionService) UpdateLastActive(ctx context.Context, sessionID uuid.UUID) error {
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.LastActiveAt = time.Now()
	session.UpdatedAt = time.Now()

	if err := ss.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session last active: %w", err)
	}

	return nil
}

func (ss *SessionService) ExpireSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if !session.CanTransitionToStatus(models.SessionStatusExpired) {
		return fmt.Errorf("cannot expire session in status: %s", session.Status)
	}

	session.Status = models.SessionStatusExpired
	session.UpdatedAt = time.Now()

	if err := ss.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to expire session: %w", err)
	}

	// Publish session expired event
	if err := ss.publishSessionEvent(ctx, sharedModels.EventTypeUserSessionExpired, session.UserID, map[string]interface{}{
		"session_id": session.ID,
		"reason":     "expired",
	}); err != nil {
		fmt.Printf("Failed to publish session expired event: %v\n", err)
	}

	return nil
}

func (ss *SessionService) TerminateSession(ctx context.Context, sessionID uuid.UUID, reason string) error {
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if !session.CanTransitionToStatus(models.SessionStatusTerminated) {
		return fmt.Errorf("cannot terminate session in status: %s", session.Status)
	}

	session.Status = models.SessionStatusTerminated
	session.UpdatedAt = time.Now()

	if err := ss.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}

	// Publish session terminated event
	if err := ss.publishSessionEvent(ctx, "session.terminated", session.UserID, map[string]interface{}{
		"session_id": session.ID,
		"reason":     reason,
	}); err != nil {
		fmt.Printf("Failed to publish session terminated event: %v\n", err)
	}

	return nil
}

func (ss *SessionService) TerminateAllUserSessions(ctx context.Context, userID uuid.UUID, reason string) error {
	sessions, err := ss.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.Status == models.SessionStatusActive {
			if err := ss.TerminateSession(ctx, session.ID, reason); err != nil {
				fmt.Printf("Failed to terminate session %s: %v\n", session.ID, err)
			}
		}
	}

	return nil
}

func (ss *SessionService) RevokeSession(ctx context.Context, sessionID uuid.UUID, reason string) error {
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if !session.CanTransitionToStatus(models.SessionStatusRevoked) {
		return fmt.Errorf("cannot revoke session in status: %s", session.Status)
	}

	session.Status = models.SessionStatusRevoked
	session.UpdatedAt = time.Now()

	if err := ss.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	// Publish session revoked event
	if err := ss.publishSessionEvent(ctx, "session.revoked", session.UserID, map[string]interface{}{
		"session_id": session.ID,
		"reason":     reason,
	}); err != nil {
		fmt.Printf("Failed to publish session revoked event: %v\n", err)
	}

	return nil
}

func (ss *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	return ss.sessionRepo.CleanupExpiredSessions(ctx)
}

func (ss *SessionService) GetActiveSessionCount(ctx context.Context, userID uuid.UUID) (int, error) {
	sessions, err := ss.GetUserSessions(ctx, userID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, session := range sessions {
		if session.Status == models.SessionStatusActive && time.Now().Before(session.ExpiresAt) {
			count++
		}
	}

	return count, nil
}

func (ss *SessionService) ValidateSessionActive(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	session, err := ss.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is active
	if session.Status != models.SessionStatusActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Auto-expire the session
		if err := ss.ExpireSession(ctx, sessionID); err != nil {
			fmt.Printf("Failed to auto-expire session: %v\n", err)
		}
		return nil, fmt.Errorf("session has expired")
	}

	return session, nil
}

// Private helper methods

func (ss *SessionService) generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (ss *SessionService) validateCreateSessionRequest(req *CreateSessionRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	return nil
}

func (ss *SessionService) publishSessionEvent(ctx context.Context, eventType sharedModels.EventType, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategorySecurity,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("Session event: %s", eventType),
		AggregateID:   userID.String(),
		AggregateType: "user",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
			Source:      "auth-service",
			Environment: "production",
			Region:      "sea",
		},
	}

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return ss.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type CreateSessionRequest struct {
	UserID     uuid.UUID              `json:"user_id" binding:"required"`
	UserAgent  string                 `json:"user_agent"`
	IPAddress  string                 `json:"ip_address"`
	DeviceInfo map[string]interface{} `json:"device_info"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type UpdateSessionRequest struct {
	SessionID  uuid.UUID              `json:"session_id" binding:"required"`
	UserAgent  string                 `json:"user_agent"`
	IPAddress  string                 `json:"ip_address"`
	DeviceInfo map[string]interface{} `json:"device_info"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type SessionStatsResponse struct {
	TotalSessions    int64 `json:"total_sessions"`
	ActiveSessions   int64 `json:"active_sessions"`
	ExpiredSessions  int64 `json:"expired_sessions"`
	RevokedSessions  int64 `json:"revoked_sessions"`
}

type UserSessionsResponse struct {
	Sessions []*SessionDetailsResponse `json:"sessions"`
	Total    int                       `json:"total"`
}

type SessionDetailsResponse struct {
	ID               uuid.UUID              `json:"id"`
	Status           models.SessionStatus   `json:"status"`
	UserAgent        string                 `json:"user_agent"`
	IPAddress        string                 `json:"ip_address"`
	DeviceInfo       map[string]interface{} `json:"device_info"`
	CreatedAt        time.Time              `json:"created_at"`
	LastActiveAt     time.Time              `json:"last_active_at"`
	ExpiresAt        time.Time              `json:"expires_at"`
	IsCurrent        bool                   `json:"is_current"`
}

// UpdateSession updates session device info and metadata
func (ss *SessionService) UpdateSession(ctx context.Context, req *UpdateSessionRequest) (*models.Session, error) {
	session, err := ss.GetSessionByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	// Update session fields if provided
	updated := false
	if req.UserAgent != "" {
		session.UserAgent = req.UserAgent
		updated = true
	}
	if req.IPAddress != "" {
		session.IPAddress = req.IPAddress
		updated = true
	}
	if req.DeviceInfo != nil {
		session.DeviceInfo = req.DeviceInfo
		updated = true
	}
	if req.Metadata != nil {
		session.Metadata = req.Metadata
		updated = true
	}

	if updated {
		session.UpdatedAt = time.Now()
		if err := ss.sessionRepo.Update(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to update session: %w", err)
		}
	}

	return session, nil
}

// RefreshTokens is an alias for RotateTokens to match handler expectations
func (ss *SessionService) RefreshTokens(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	return ss.RotateTokens(ctx, sessionID)
}

// InvalidateUserSessions is an alias for TerminateAllUserSessions to match handler expectations
func (ss *SessionService) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	return ss.TerminateAllUserSessions(ctx, userID, "user_logout_all")
}

// InvalidateSession is an alias for TerminateSession to match handler expectations
func (ss *SessionService) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	return ss.TerminateSession(ctx, sessionID, "user_logout")
}

func ToSessionDetailsResponse(session *models.Session, isCurrentSession bool) *SessionDetailsResponse {
	return &SessionDetailsResponse{
		ID:           session.ID,
		Status:       session.Status,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
		DeviceInfo:   session.DeviceInfo,
		CreatedAt:    session.CreatedAt,
		LastActiveAt: session.LastActiveAt,
		ExpiresAt:    session.ExpiresAt,
		IsCurrent:    isCurrentSession,
	}
}