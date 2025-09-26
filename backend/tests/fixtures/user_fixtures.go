package fixtures

import (
	"time"

	"github.com/google/uuid"
	"tchat-backend/auth/models"
)

// UserFixtures provides test data for User models
type UserFixtures struct {
	*BaseFixture
}

// NewUserFixtures creates a new user fixtures instance
func NewUserFixtures(seed ...int64) *UserFixtures {
	return &UserFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BasicUser creates a basic user for testing
func (u *UserFixtures) BasicUser(country string) *models.User {
	countryCode := u.CountryCode(country)
	phone := u.Phone(countryCode)
	email := u.Email("test-user", "tchat-test.com")
	name := u.Name(countryCode)

	return &models.User{
		ID:         u.UUID("basic-user-" + countryCode),
		Phone:      &phone,
		Email:      &email,
		Name:       name,
		Avatar:     nil,
		Country:    models.Country(countryCode),
		Locale:     u.Locale(countryCode),
		KYCTier:    models.KYCTier(u.KYCTier()),
		Status:     models.UserStatusActive,
		LastSeen:   nil,
		IsVerified: false,
		CreatedAt:  u.PastTime(60), // Created 1 hour ago
		UpdatedAt:  u.PastTime(30), // Updated 30 minutes ago
	}
}

// VerifiedUser creates a verified user for testing
func (u *UserFixtures) VerifiedUser(country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("verified-user-" + country)
	user.IsVerified = true
	user.KYCTier = models.KYCTier2
	lastSeen := u.PastTime(5) // Last seen 5 minutes ago
	user.LastSeen = &lastSeen
	return user
}

// PremiumUser creates a premium user (KYC Tier 3) for testing
func (u *UserFixtures) PremiumUser(country string) *models.User {
	user := u.VerifiedUser(country)
	user.ID = u.UUID("premium-user-" + country)
	user.KYCTier = models.KYCTier3
	avatar := "https://example.com/avatars/premium-user.jpg"
	user.Avatar = &avatar
	return user
}

// TestUsers creates a collection of test users across different countries
func (u *UserFixtures) TestUsers() []*models.User {
	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}
	users := make([]*models.User, 0, len(countries)*3)

	for _, country := range countries {
		users = append(users, u.BasicUser(country))
		users = append(users, u.VerifiedUser(country))
		users = append(users, u.PremiumUser(country))
	}

	return users
}

// UserWithEmail creates a user with specific email
func (u *UserFixtures) UserWithEmail(email, country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("user-email-" + email)
	user.Email = &email
	return user
}

// UserWithPhone creates a user with specific phone
func (u *UserFixtures) UserWithPhone(phone, country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("user-phone-" + phone)
	user.Phone = &phone
	return user
}

// InactiveUser creates an inactive user for testing
func (u *UserFixtures) InactiveUser(country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("inactive-user-" + country)
	user.Status = models.UserStatusInactive
	user.LastSeen = nil
	return user
}

// SuspendedUser creates a suspended user for testing
func (u *UserFixtures) SuspendedUser(country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("suspended-user-" + country)
	user.Status = models.UserStatusSuspended
	lastSeen := u.PastTime(1440) // Last seen 24 hours ago
	user.LastSeen = &lastSeen
	return user
}

// DeletedUser creates a deleted user for testing
func (u *UserFixtures) DeletedUser(country string) *models.User {
	user := u.BasicUser(country)
	user.ID = u.UUID("deleted-user-" + country)
	user.Status = models.UserStatusInactive
	user.IsVerified = false
	user.Phone = nil // Privacy: phone removed
	user.Email = nil // Privacy: email removed
	return user
}

// SessionFixtures provides test data for Session models
type SessionFixtures struct {
	*BaseFixture
}

// NewSessionFixtures creates a new session fixtures instance
func NewSessionFixtures(seed ...int64) *SessionFixtures {
	return &SessionFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// ActiveSession creates an active session for testing
func (s *SessionFixtures) ActiveSession(userID uuid.UUID, platform string) *models.Session {
	deviceID := s.DeviceID(platform)
	ipAddress := s.IPAddress("TH") // Default to Thailand IP
	userAgent := s.UserAgent(platform)

	return &models.Session{
		ID:               s.UUID("session-" + userID.String()),
		UserID:           userID,
		DeviceID:         deviceID,
		AccessToken:      s.Token(32),
		RefreshToken:     s.Token(32),
		ExpiresAt:        s.FutureTime(15), // Expires in 15 minutes
		RefreshExpiresAt: s.FutureTime(1440), // Refresh expires in 24 hours
		IsActive:         true,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		DeviceInfo:       map[string]interface{}{
			"platform": platform,
			"version":  "1.0.0",
			"os":       platform,
		},
		Metadata: map[string]interface{}{
			"login_method": "otp",
			"location":     "Bangkok, Thailand",
		},
		CreatedAt:    s.PastTime(30), // Created 30 minutes ago
		UpdatedAt:    s.PastTime(5),  // Updated 5 minutes ago
		LastActiveAt: s.PastTime(1),  // Last active 1 minute ago
		LastUsed:     s.PastTime(1),  // Last used 1 minute ago
		RevokedAt:    nil,
		Status:       models.SessionStatusActive,
	}
}

// ExpiredSession creates an expired session for testing
func (s *SessionFixtures) ExpiredSession(userID uuid.UUID, platform string) *models.Session {
	session := s.ActiveSession(userID, platform)
	session.ID = s.UUID("expired-session-" + userID.String())
	session.ExpiresAt = s.PastTime(60) // Expired 1 hour ago
	session.IsActive = false
	session.Status = models.SessionStatusExpired
	return session
}

// RevokedSession creates a revoked session for testing
func (s *SessionFixtures) RevokedSession(userID uuid.UUID, platform string) *models.Session {
	session := s.ActiveSession(userID, platform)
	session.ID = s.UUID("revoked-session-" + userID.String())
	session.IsActive = false
	session.Status = models.SessionStatusRevoked
	revokedAt := s.PastTime(30)
	session.RevokedAt = &revokedAt
	return session
}

// MultiDeviceSessions creates sessions for multiple devices
func (s *SessionFixtures) MultiDeviceSessions(userID uuid.UUID) []*models.Session {
	platforms := []string{"ios", "android", "web"}
	sessions := make([]*models.Session, 0, len(platforms))

	for _, platform := range platforms {
		session := s.ActiveSession(userID, platform)
		session.ID = s.UUID("multi-session-" + platform + "-" + userID.String())
		sessions = append(sessions, session)
	}

	return sessions
}

// KYCFixtures provides test data for KYC models
type KYCFixtures struct {
	*BaseFixture
}

// NewKYCFixtures creates a new KYC fixtures instance
func NewKYCFixtures(seed ...int64) *KYCFixtures {
	return &KYCFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BasicKYC creates basic KYC data for testing
func (k *KYCFixtures) BasicKYC(userID uuid.UUID, country string) *models.KYC {
	return &models.KYC{
		ID:           k.UUID("kyc-" + userID.String()),
		UserID:       userID,
		Tier:         models.TierBasic,
		Status:       models.KYCStatusPending,
		DocumentType: models.DocumentTypeNationalID,
		Nationality:  country,
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		FullName:     k.Name(country),
		Address: &models.Address{
			Street:     "123 Test Street",
			City:       "Bangkok",
			State:      "Bangkok",
			Country:    country,
			PostalCode: "10100",
		},
		DocumentNumber: "123456789",
		SelfieImage:    "https://example.com/selfie.jpg",
		ReviewNotes: "Automatic verification",
		CreatedAt:   k.PastTime(60),
		UpdatedAt:   k.PastTime(30),
	}
}

// VerifiedKYC creates verified KYC data for testing
func (k *KYCFixtures) VerifiedKYC(userID uuid.UUID, country string) *models.KYC {
	kyc := k.BasicKYC(userID, country)
	kyc.ID = k.UUID("verified-kyc-" + userID.String())
	kyc.Status = models.KYCStatusApproved
	kyc.Tier = models.TierIdentity
	approvedAt := k.PastTime(30)
	kyc.ApprovedAt = &approvedAt
	return kyc
}

// PremiumKYC creates premium KYC data (Tier 3) for testing
func (k *KYCFixtures) PremiumKYC(userID uuid.UUID, country string) *models.KYC {
	kyc := k.VerifiedKYC(userID, country)
	kyc.ID = k.UUID("premium-kyc-" + userID.String())
	kyc.Tier = models.TierEnhanced
	// Enhanced KYC with additional verification
	kyc.PlaceOfBirth = "Bangkok"
	kyc.Gender = "unspecified"

	return kyc
}

// RejectedKYC creates rejected KYC data for testing
func (k *KYCFixtures) RejectedKYC(userID uuid.UUID, country string) *models.KYC {
	kyc := k.BasicKYC(userID, country)
	kyc.ID = k.UUID("rejected-kyc-" + userID.String())
	kyc.Status = models.KYCStatusRejected
	kyc.RejectionReason = "Document quality insufficient, please resubmit"
	return kyc
}

// AllUserFixtures creates a complete set of user-related test data
func AllUserFixtures(seed ...int64) (*UserFixtures, *SessionFixtures, *KYCFixtures) {
	return NewUserFixtures(seed...), NewSessionFixtures(seed...), NewKYCFixtures(seed...)
}