package contract

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/auth/models"
	"tchat.dev/auth/services"
	sharedModels "tchat.dev/shared/models"
)

// Test constants
var (
	testUserID    = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testSessionID = uuid.New()
)

// Global mock service instances
var (
	mockUserService    *MockUserService
	mockAuthService    *MockAuthService
	mockSessionService *MockSessionService
	mockJWTService     *MockJWTService
)

// MockUserService provides mock implementation for testing
type MockUserService struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*sharedModels.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: make(map[uuid.UUID]*sharedModels.User),
	}
}

func (m *MockUserService) StoreTestUser(user *sharedModels.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
}

func (m *MockUserService) GetTestUser(userID uuid.UUID) *sharedModels.User {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.users[userID]
}

func (m *MockUserService) UpdateTestUser(user *sharedModels.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
}

func (m *MockUserService) SetupPhoneNumberAvailable(phoneNumber string, isAvailable bool) {
	// For simplicity, we'll just track this in the mock service
	// In a real implementation, this would check database for existing users
}

func (m *MockUserService) ClearTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users = make(map[uuid.UUID]*sharedModels.User)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*sharedModels.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, exists := m.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, user := range m.users {
		if user.Phone == phoneNumber {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserService) CreateUser(ctx context.Context, req *services.CreateUserRequest) (*sharedModels.User, error) {
	user := &sharedModels.User{
		ID:          uuid.New(),
		Name:        fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		Phone:       req.PhoneNumber,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Country:     string(req.Country),
		Locale:      req.Language + "-" + string(req.Country),
		Language:    req.Language,
		Timezone:    req.TimeZone,
		KYCTier:     int(sharedModels.KYCTierBasic),
		Status:      string(sharedModels.UserStatusActive),
		Active:      true,
		Verified:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.StoreTestUser(user)
	return user, nil
}

func (m *MockUserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *services.UpdateUserProfileRequest) (*sharedModels.User, error) {
	user := m.GetTestUser(userID)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
		// Update display name
		lastName := user.LastName
		user.Name = *req.FirstName + " " + lastName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
		// Update display name
		firstName := user.FirstName
		user.Name = firstName + " " + *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Language != nil {
		user.Language = *req.Language
	}
	if req.TimeZone != nil {
		user.Timezone = *req.TimeZone
	}

	user.UpdatedAt = time.Now()
	m.UpdateTestUser(user)

	return user, nil
}

// MockAuthService provides mock authentication functionality
type MockAuthService struct {
	mu               sync.RWMutex
	mockLogins       map[string]*sharedModels.User
	mockOTPCodes     map[string]bool
	mockVerifications map[string]*services.VerifyOTPResponse
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{
		mockLogins:       make(map[string]*sharedModels.User),
		mockOTPCodes:     make(map[string]bool),
		mockVerifications: make(map[string]*services.VerifyOTPResponse),
	}
}

func (m *MockAuthService) SetupMockLogin(phoneNumber, countryCode, password string, user *sharedModels.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s:%s", phoneNumber, countryCode, password)
	m.mockLogins[key] = user
}

func (m *MockAuthService) SetupMockOTPVerification(code string, shouldSucceed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mockOTPCodes[code] = shouldSucceed
}

func (m *MockAuthService) SetupMockOTPSending(phoneNumber string, shouldSucceed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Store OTP sending capability for phone number
	m.mockOTPCodes[phoneNumber] = shouldSucceed
}

func (m *MockAuthService) ClearTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mockLogins = make(map[string]*sharedModels.User)
	m.mockOTPCodes = make(map[string]bool)
	m.mockVerifications = make(map[string]*services.VerifyOTPResponse)
}

func (m *MockAuthService) VerifyOTP(ctx context.Context, req *services.VerifyOTPRequest) (*services.VerifyOTPResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if OTP code is configured to succeed
	if success, exists := m.mockOTPCodes[req.Code]; exists && success {
		// Find user by phone number
		if mockUserService != nil {
			user, err := mockUserService.GetUserByPhoneNumber(ctx, req.PhoneNumber)
			if err != nil {
				return nil, err
			}

			// Create mock session
			session := &services.SessionResponse{
				ID:           testSessionID,
				AccessToken:  "mock-access-token",
				RefreshToken: "mock-refresh-token",
				ExpiresAt:    time.Now().Add(15 * time.Minute),
			}

			// Create user response manually since ToResponse method doesn't exist
			userResponse := &services.UserResponse{
				ID:              user.ID,
				Username:        user.Username,
				PhoneNumber:     user.Phone,
				Email:           user.Email,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				Country:         user.Country,
				Language:        user.Language,
				TimeZone:        user.Timezone,
				KYCTier:         user.KYCTier,
				Status:          user.Status,
				IsPhoneVerified: user.Verified,
				IsEmailVerified: user.Verified,
			}

			return &services.VerifyOTPResponse{
				Success:    true,
				OTPID:      uuid.New(),
				VerifiedAt: time.Now(),
				User:       userResponse,
				Session:    session,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid OTP")
}

// MockSessionService provides mock session management
type MockSessionService struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*models.Session
}

func NewMockSessionService() *MockSessionService {
	return &MockSessionService{
		sessions: make(map[uuid.UUID]*models.Session),
	}
}

func (m *MockSessionService) StoreTestSession(session *models.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.ID] = session
}

func (m *MockSessionService) ClearTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions = make(map[uuid.UUID]*models.Session)
}

func (m *MockSessionService) CreateSession(ctx context.Context, req *services.CreateSessionRequest) (*models.Session, error) {
	session := &models.Session{
		ID:               testSessionID,
		UserID:           req.UserID,
		Status:           models.SessionStatusActive,
		AccessToken:      "mock-access-token",
		RefreshToken:     "mock-refresh-token",
		ExpiresAt:        time.Now().Add(15 * time.Minute),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		UserAgent:        req.UserAgent,
		IPAddress:        req.IPAddress,
		DeviceInfo:       req.DeviceInfo,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		LastActiveAt:     time.Now(),
	}

	m.mu.Lock()
	m.sessions[session.ID] = session
	m.mu.Unlock()

	return session, nil
}

func (m *MockSessionService) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (m *MockSessionService) TerminateSession(ctx context.Context, sessionID uuid.UUID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, exists := m.sessions[sessionID]; exists {
		session.Status = models.SessionStatusTerminated
		session.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockSessionService) ConvertSessionToDetailsResponse(session *models.Session, isCurrent bool) *services.SessionDetailsResponse {
	return &services.SessionDetailsResponse{
		ID:           session.ID,
		Status:       session.Status,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
		DeviceInfo:   session.DeviceInfo,
		CreatedAt:    session.CreatedAt,
		LastActiveAt: session.LastActiveAt,
		ExpiresAt:    session.ExpiresAt,
		IsCurrent:    isCurrent,
	}
}

// MockJWTService provides mock JWT functionality
type MockJWTService struct {
	mu            sync.RWMutex
	validTokens   map[string]*sharedModels.User
	refreshTokens map[string]*sharedModels.User
	realService   *services.JWTService
}

func NewMockJWTService() *MockJWTService {
	return &MockJWTService{
		validTokens:   make(map[string]*sharedModels.User),
		refreshTokens: make(map[string]*sharedModels.User),
	}
}

func (m *MockJWTService) StoreValidToken(token string, user *sharedModels.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.validTokens[token] = user
}

func (m *MockJWTService) StoreValidRefreshToken(refreshToken string, user *sharedModels.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.refreshTokens[refreshToken] = user
}

func (m *MockJWTService) ClearTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.validTokens = make(map[string]*sharedModels.User)
	m.refreshTokens = make(map[string]*sharedModels.User)
}

func (m *MockJWTService) ValidateAccessToken(ctx context.Context, tokenString string) (*services.UserClaims, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if token is in our valid tokens
	if user, exists := m.validTokens[tokenString]; exists {
		return &services.UserClaims{
			UserID:      user.ID,
			PhoneNumber: getPhoneNumber(user),
			CountryCode: string(user.Country),
			KYCStatus:   getKYCStatus(user),
			KYCLevel:    int(user.KYCTier),
			SessionID:   testSessionID,
			DeviceID:    "test-device",
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Helper functions for mock services
func getPhoneNumber(user *sharedModels.User) string {
	return user.Phone
}

func getKYCStatus(user *sharedModels.User) string {
	if user.Verified {
		return "verified"
	}
	return "pending"
}

// MockAuthMiddleware provides authentication middleware for testing
func MockAuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from Bearer format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := authHeader[7:]

		// For testing, check if we have a mock valid token first
		if mockJWTService != nil {
			if user, exists := mockJWTService.validTokens[tokenString]; exists {
				// Set user context for handlers
				c.Set("user_id", user.ID)
				c.Set("user_claims", &services.UserClaims{
					UserID:      user.ID,
					PhoneNumber: getPhoneNumber(user),
					CountryCode: string(user.Country),
					SessionID:   testSessionID,
				})
				c.Next()
				return
			}
		}

		// Fallback to real JWT validation
		claims, err := jwtService.ValidateAccessToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set context for handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_claims", claims)
		c.Next()
	}
}