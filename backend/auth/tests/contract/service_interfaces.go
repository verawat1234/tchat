package contract

import (
	"context"
	"time"

	"github.com/google/uuid"

	"tchat.dev/auth/models"
	"tchat.dev/auth/services"
	sharedModels "tchat.dev/shared/models"
)

// Service interfaces for contract testing

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	PhoneNumber string                 `json:"phone_number"`
	Email       string                 `json:"email"`
	Username    string                 `json:"username"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Country     sharedModels.Country   `json:"country"`
	Language    string                 `json:"language"`
	TimeZone    string                 `json:"timezone"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateUserProfileRequest represents a profile update request
type UpdateUserProfileRequest struct {
	Username    *string                    `json:"username,omitempty"`
	FirstName   *string                    `json:"first_name,omitempty"`
	LastName    *string                    `json:"last_name,omitempty"`
	Email       *string                    `json:"email,omitempty"`
	Language    *string                    `json:"language,omitempty"`
	TimeZone    *string                    `json:"timezone,omitempty"`
	Preferences *UserPreferences           `json:"preferences,omitempty"`
	Metadata    map[string]interface{}     `json:"metadata,omitempty"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	Theme               string `json:"theme,omitempty"`
	Language            string `json:"language,omitempty"`
	NotificationsEmail  bool   `json:"notifications_email"`
	NotificationsPush   bool   `json:"notifications_push"`
	PrivacyLevel        string `json:"privacy_level,omitempty"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	PhoneNumber string                 `json:"phone_number"`
	Code        string                 `json:"code"`
	UserAgent   string                 `json:"user_agent"`
	IPAddress   string                 `json:"ip_address"`
	DeviceInfo  map[string]interface{} `json:"device_info"`
}

// VerifyOTPResponse represents OTP verification response
type VerifyOTPResponse struct {
	Success    bool                      `json:"success"`
	OTPID      uuid.UUID                 `json:"otp_id"`
	VerifiedAt time.Time                 `json:"verified_at"`
	User       *services.UserResponse    `json:"user,omitempty"`
	Session    *SessionResponse          `json:"session,omitempty"`
}

// SessionResponse represents session information
type SessionResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// CreateSessionRequest represents session creation request
type CreateSessionRequest struct {
	UserID     uuid.UUID              `json:"user_id"`
	UserAgent  string                 `json:"user_agent"`
	IPAddress  string                 `json:"ip_address"`
	DeviceInfo map[string]interface{} `json:"device_info"`
}

// SessionDetailsResponse represents detailed session information
type SessionDetailsResponse struct {
	ID           uuid.UUID              `json:"id"`
	Status       models.SessionStatus   `json:"status"`
	UserAgent    string                 `json:"user_agent"`
	IPAddress    string                 `json:"ip_address"`
	DeviceInfo   map[string]interface{} `json:"device_info"`
	CreatedAt    time.Time              `json:"created_at"`
	LastActiveAt time.Time              `json:"last_active_at"`
	ExpiresAt    time.Time              `json:"expires_at"`
	IsCurrent    bool                   `json:"is_current"`
}

// UserClaims represents JWT claims for a user
type UserClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	CountryCode string    `json:"country_code"`
	KYCStatus   string    `json:"kyc_status"`
	KYCLevel    int       `json:"kyc_level"`
	SessionID   uuid.UUID `json:"session_id"`
	DeviceID    string    `json:"device_id"`
	Permissions []string  `json:"permissions,omitempty"`
	Scopes      []string  `json:"scopes,omitempty"`
}

// UserServiceInterface defines the contract for user operations
type UserServiceInterface interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*sharedModels.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*sharedModels.User, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *UpdateUserProfileRequest) (*sharedModels.User, error)
}

// AuthServiceInterface defines the contract for auth operations
type AuthServiceInterface interface {
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error)
}

// SessionServiceInterface defines the contract for session operations
type SessionServiceInterface interface {
	CreateSession(ctx context.Context, req *CreateSessionRequest) (*models.Session, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	TerminateSession(ctx context.Context, sessionID uuid.UUID, reason string) error
	ConvertSessionToDetailsResponse(session *models.Session, isCurrent bool) *SessionDetailsResponse
}

// JWTServiceInterface defines the contract for JWT operations
type JWTServiceInterface interface {
	ValidateAccessToken(ctx context.Context, tokenString string) (*services.UserClaims, error)
}