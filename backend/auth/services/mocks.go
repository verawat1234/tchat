package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/auth/models"
	sharedModels "tchat.dev/shared/models"
)

// Mock implementations for testing

// NoOpEventPublisher is a no-operation event publisher for testing
type NoOpEventPublisher struct{}

func (n *NoOpEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	return nil
}

// InMemoryOTPRepository is an in-memory OTP repository for testing
type InMemoryOTPRepository struct {
	otps map[string]string
	mu   sync.RWMutex
}

func NewInMemoryOTPRepository() *InMemoryOTPRepository {
	return &InMemoryOTPRepository{
		otps: make(map[string]string),
	}
}

func (r *InMemoryOTPRepository) Store(ctx context.Context, phoneNumber, code string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[phoneNumber] = code
	return nil
}

func (r *InMemoryOTPRepository) Verify(ctx context.Context, phoneNumber, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storedCode, exists := r.otps[phoneNumber]
	return exists && storedCode == code, nil
}

func (r *InMemoryOTPRepository) Delete(ctx context.Context, phoneNumber string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.otps, phoneNumber)
	return nil
}

func (r *InMemoryOTPRepository) SetTestCode(phoneNumber, code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[phoneNumber] = code
}

func (r *InMemoryOTPRepository) Create(ctx context.Context, otp *OTP) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[otp.PhoneNumber] = otp.Code
	return nil
}

func (r *InMemoryOTPRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*OTP, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	code, exists := r.otps[phoneNumber]
	if !exists {
		return nil, fmt.Errorf("OTP not found for phone number: %s", phoneNumber)
	}
	return &OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   time.Now().Add(10 * time.Minute), // Mock 10 minute expiry
	}, nil
}

func (r *InMemoryOTPRepository) GetByID(ctx context.Context, id uuid.UUID) (*OTP, error) {
	// Mock implementation - not applicable for in-memory storage
	return nil, fmt.Errorf("GetByID not implemented for in-memory OTP repository")
}

func (r *InMemoryOTPRepository) Update(ctx context.Context, otp *OTP) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[otp.PhoneNumber] = otp.Code
	return nil
}

func (r *InMemoryOTPRepository) DeleteExpired(ctx context.Context) error {
	// Mock implementation - no expiry tracking in simple in-memory storage
	return nil
}

func (r *InMemoryOTPRepository) GetAttemptCount(ctx context.Context, phoneNumber string, timeWindow time.Duration) (int, error) {
	// Mock implementation - always return 1 attempt
	return 1, nil
}

// MockSMSProvider is a mock SMS provider for testing
type MockSMSProvider struct {
	sentMessages []SMSMessage
	mu           sync.RWMutex
}

type SMSMessage struct {
	PhoneNumber string
	Message     string
	SentAt      time.Time
}

func NewMockSMSProvider() *MockSMSProvider {
	return &MockSMSProvider{
		sentMessages: make([]SMSMessage, 0),
	}
}

func (p *MockSMSProvider) SendSMS(ctx context.Context, phoneNumber, message string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.sentMessages = append(p.sentMessages, SMSMessage{
		PhoneNumber: phoneNumber,
		Message:     message,
		SentAt:      time.Now(),
	})

	return nil
}

func (p *MockSMSProvider) GetSentMessages() []SMSMessage {
	p.mu.RLock()
	defer p.mu.RUnlock()

	messages := make([]SMSMessage, len(p.sentMessages))
	copy(messages, p.sentMessages)
	return messages
}

// NoOpRateLimiter is a no-operation rate limiter for testing
type NoOpRateLimiter struct{}

func (r *NoOpRateLimiter) IsAllowed(ctx context.Context, key string) (bool, error) {
	return true, nil
}

func (r *NoOpRateLimiter) Increment(ctx context.Context, key string) error {
	return nil
}

// NoOpSecurityLogger is a no-operation security logger for testing
type NoOpSecurityLogger struct{}

func (l *NoOpSecurityLogger) LogSecurityEvent(ctx context.Context, event string, details map[string]interface{}) error {
	return nil
}

// PostgreSQLUserRepository is a PostgreSQL implementation of UserRepository
type PostgreSQLUserRepository struct {
	DB *gorm.DB
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *sharedModels.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *PostgreSQLUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.User, error) {
	var user sharedModels.User
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgreSQLUserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error) {
	var user sharedModels.User
	err := r.DB.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *sharedModels.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&sharedModels.User{}, id).Error
}

// PostgreSQLSessionRepository is a PostgreSQL implementation of SessionRepository
type PostgreSQLSessionRepository struct {
	DB *gorm.DB
}

func (r *PostgreSQLSessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.DB.WithContext(ctx).Create(session).Error
}

func (r *PostgreSQLSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PostgreSQLSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (r *PostgreSQLSessionRepository) Update(ctx context.Context, session *models.Session) error {
	return r.DB.WithContext(ctx).Save(session).Error
}

func (r *PostgreSQLSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Session{}, id).Error
}

func (r *PostgreSQLSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

// Extension to AuthService for testing
func (as *AuthService) SetTestOTPCode(phoneNumber, code string) error {
	if otpRepo, ok := as.otpRepo.(*InMemoryOTPRepository); ok {
		otpRepo.SetTestCode(phoneNumber, code)
		return nil
	}
	return fmt.Errorf("OTP repository does not support test codes")
}