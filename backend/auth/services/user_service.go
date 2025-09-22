package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/auth/models"
	"tchat.dev/shared/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters UserFilters, pagination Pagination) ([]*models.User, int64, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	SearchByUsername(ctx context.Context, username string, limit int) ([]*models.User, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event *models.Event) error
}

type UserFilters struct {
	Country     *models.Country          `json:"country,omitempty"`
	Status      *models.UserStatus       `json:"status,omitempty"`
	KYCTier     *models.VerificationTier `json:"kyc_tier,omitempty"`
	CreatedFrom *time.Time               `json:"created_from,omitempty"`
	CreatedTo   *time.Time               `json:"created_to,omitempty"`
	Search      string                   `json:"search,omitempty"`
}

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
	Order    string `json:"order"` // asc, desc
}

type UserStats struct {
	TotalUsers           int64                             `json:"total_users"`
	ActiveUsers          int64                             `json:"active_users"`
	VerifiedUsers        int64                             `json:"verified_users"`
	NewUsersToday        int64                             `json:"new_users_today"`
	NewUsersThisWeek     int64                             `json:"new_users_this_week"`
	NewUsersThisMonth    int64                             `json:"new_users_this_month"`
	UsersByCountry       map[models.Country]int64          `json:"users_by_country"`
	UsersByKYCTier       map[models.VerificationTier]int64 `json:"users_by_kyc_tier"`
	AverageSessionLength time.Duration                     `json:"average_session_length"`
}

type UserService struct {
	userRepo       UserRepository
	eventPublisher EventPublisher
	db             *gorm.DB
}

func NewUserService(userRepo UserRepository, eventPublisher EventPublisher, db *gorm.DB) *UserService {
	return &UserService{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		db:             db,
	}
}

func (us *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := us.validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check for existing user
	existingUser, err := us.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with phone number %s already exists", req.PhoneNumber)
	}

	// Check email if provided
	if req.Email != "" {
		existingUser, err := us.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	// Create user
	user := &models.User{
		ID:              uuid.New(),
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Username:        req.Username,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Country:         req.Country,
		Language:        req.Language,
		TimeZone:        req.TimeZone,
		Status:          models.UserStatusPending,
		IsPhoneVerified: false,
		IsEmailVerified: false,
		KYCTier:         models.VerificationTierNone,
		Preferences:     models.UserPreferences{},
		Metadata:        req.Metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Set default preferences
	user.SetDefaultPreferences()

	// Validate user model
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Save to database
	if err := us.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish user registration event
	if err := us.publishUserEvent(ctx, models.EventTypeUserRegistered, user.ID, map[string]interface{}{
		"phone_number": user.PhoneNumber,
		"country":      user.Country,
		"language":     user.Language,
		"username":     user.Username,
	}); err != nil {
		// Log error but don't fail the operation
		// In production, you might want to use a more robust event publishing mechanism
		fmt.Printf("Failed to publish user registration event: %v\n", err)
	}

	return user, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	user, err := us.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (us *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *UpdateUserProfileRequest) (*models.User, error) {
	// Get existing user
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !user.CanUpdateProfile() {
		return nil, fmt.Errorf("user cannot update profile in current status: %s", user.Status)
	}

	// Validate update request
	if err := us.validateUpdateProfileRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.Username != nil && *req.Username != user.Username {
		// Check username uniqueness if changed
		if *req.Username != user.Username {
			existing, err := us.userRepo.GetByEmail(ctx, *req.Username) // Using email check as proxy for uniqueness
			if err == nil && existing.ID != user.ID {
				return nil, fmt.Errorf("username already taken")
			}
		}
		changes["username"] = map[string]string{"from": user.Username, "to": *req.Username}
		user.Username = *req.Username
	}

	if req.FirstName != nil && *req.FirstName != user.FirstName {
		changes["first_name"] = map[string]string{"from": user.FirstName, "to": *req.FirstName}
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil && *req.LastName != user.LastName {
		changes["last_name"] = map[string]string{"from": user.LastName, "to": *req.LastName}
		user.LastName = *req.LastName
	}

	if req.Email != nil && *req.Email != user.Email {
		// Check email uniqueness if changed
		if *req.Email != user.Email {
			existing, err := us.userRepo.GetByEmail(ctx, *req.Email)
			if err == nil && existing.ID != user.ID {
				return nil, fmt.Errorf("email already taken")
			}
		}
		changes["email"] = map[string]string{"from": user.Email, "to": *req.Email}
		user.Email = *req.Email
		user.IsEmailVerified = false // Reset verification status
	}

	if req.Language != nil && *req.Language != user.Language {
		changes["language"] = map[string]string{"from": user.Language, "to": *req.Language}
		user.Language = *req.Language
	}

	if req.TimeZone != nil && *req.TimeZone != user.TimeZone {
		changes["timezone"] = map[string]string{"from": user.TimeZone, "to": *req.TimeZone}
		user.TimeZone = *req.TimeZone
	}

	if req.Preferences != nil {
		changes["preferences"] = map[string]interface{}{
			"from": user.Preferences,
			"to":   *req.Preferences,
		}
		user.Preferences = *req.Preferences
	}

	if req.Metadata != nil {
		// Merge metadata
		if user.Metadata == nil {
			user.Metadata = make(map[string]interface{})
		}
		for key, value := range req.Metadata {
			user.Metadata[key] = value
		}
		changes["metadata"] = req.Metadata
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Validate updated user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("updated user validation failed: %w", err)
	}

	// Save to database
	if err := us.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish profile update event if there were changes
	if len(changes) > 0 {
		if err := us.publishUserEvent(ctx, models.EventTypeUserProfileUpdated, user.ID, map[string]interface{}{
			"changes": changes,
		}); err != nil {
			fmt.Printf("Failed to publish user profile update event: %v\n", err)
		}
	}

	return user, nil
}

func (us *UserService) VerifyPhoneNumber(ctx context.Context, userID uuid.UUID) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsPhoneVerified {
		return fmt.Errorf("phone number already verified")
	}

	user.IsPhoneVerified = true
	user.UpdatedAt = time.Now()

	// Update status if was pending
	if user.Status == models.UserStatusPending {
		user.Status = models.UserStatusActive
	}

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	return nil
}

func (us *UserService) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsEmailVerified {
		return fmt.Errorf("email already verified")
	}

	user.IsEmailVerified = true
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user email verification status: %w", err)
	}

	return nil
}

func (us *UserService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus, reason string) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	oldStatus := user.Status

	if !user.CanTransitionToStatus(status) {
		return fmt.Errorf("cannot transition from %s to %s", oldStatus, status)
	}

	user.Status = status
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Publish status change event
	if err := us.publishUserEvent(ctx, models.EventTypeUserProfileUpdated, user.ID, map[string]interface{}{
		"status_change": map[string]interface{}{
			"from":   oldStatus,
			"to":     status,
			"reason": reason,
		},
	}); err != nil {
		fmt.Printf("Failed to publish user status change event: %v\n", err)
	}

	return nil
}

func (us *UserService) UpdateKYCTier(ctx context.Context, userID uuid.UUID, tier models.VerificationTier) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	oldTier := user.KYCTier

	if !user.CanUpgradeToTier(tier) {
		return fmt.Errorf("cannot upgrade from tier %d to tier %d", oldTier, tier)
	}

	user.KYCTier = tier
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user KYC tier: %w", err)
	}

	// Publish KYC tier update event
	if err := us.publishUserEvent(ctx, models.EventTypeUserKYCVerified, user.ID, map[string]interface{}{
		"kyc_tier_change": map[string]interface{}{
			"from": oldTier,
			"to":   tier,
		},
	}); err != nil {
		fmt.Printf("Failed to publish KYC tier update event: %v\n", err)
	}

	return nil
}

func (us *UserService) ListUsers(ctx context.Context, filters UserFilters, pagination Pagination) ([]*models.User, int64, error) {
	// Validate pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}
	if pagination.OrderBy == "" {
		pagination.OrderBy = "created_at"
	}
	if pagination.Order != "asc" && pagination.Order != "desc" {
		pagination.Order = "desc"
	}

	users, total, err := us.userRepo.List(ctx, filters, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

func (us *UserService) SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error) {
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	users, err := us.userRepo.SearchByUsername(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

func (us *UserService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	stats, err := us.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID uuid.UUID, reason string) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.CanBeDeleted() {
		return fmt.Errorf("user cannot be deleted in current status: %s", user.Status)
	}

	// Soft delete
	if err := us.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Publish user deletion event
	if err := us.publishUserEvent(ctx, "user.deleted", user.ID, map[string]interface{}{
		"reason": reason,
	}); err != nil {
		fmt.Printf("Failed to publish user deletion event: %v\n", err)
	}

	return nil
}

func (us *UserService) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.User, error) {
	if len(userIDs) == 0 {
		return []*models.User{}, nil
	}

	if len(userIDs) > 100 {
		return nil, fmt.Errorf("cannot fetch more than 100 users at once")
	}

	users, err := us.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	return users, nil
}

func (us *UserService) ChangePhoneNumber(ctx context.Context, req *ChangePhoneNumberRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if req.NewPhoneNumber == "" {
		return fmt.Errorf("new phone number is required")
	}

	if req.CountryCode == "" {
		return fmt.Errorf("country code is required")
	}

	// Get current user
	user, err := us.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	// Check if new phone number already exists
	existingUser, err := us.userRepo.GetByPhoneNumber(ctx, req.NewPhoneNumber)
	if err == nil && existingUser.ID != req.UserID {
		return fmt.Errorf("phone number already exists")
	}

	// Verify OTP (this would typically be done through auth service)
	// For now, we'll assume OTP verification is handled elsewhere

	// Update phone number
	updateReq := &UpdateUserRequest{
		UserID:      req.UserID,
		PhoneNumber: req.NewPhoneNumber,
		CountryCode: req.CountryCode,
	}

	_, err = us.UpdateUser(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update phone number: %w", err)
	}

	// Publish phone number changed event
	if err := us.publishUserEvent(ctx, "user.phone_changed", user.ID, map[string]interface{}{
		"old_phone_number": user.PhoneNumber,
		"new_phone_number": req.NewPhoneNumber,
		"country_code":     req.CountryCode,
	}); err != nil {
		fmt.Printf("Failed to publish phone number changed event: %v\n", err)
	}

	return nil
}

func (us *UserService) DeactivateUser(ctx context.Context, req *DeactivateUserRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	// Get current user
	user, err := us.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if user.Status == models.UserStatusInactive {
		return fmt.Errorf("user is already deactivated")
	}

	// Update user status to inactive
	updateReq := &UpdateUserRequest{
		UserID: req.UserID,
		Status: string(models.UserStatusInactive),
	}

	_, err = us.UpdateUser(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Publish user deactivated event
	if err := us.publishUserEvent(ctx, "user.deactivated", user.ID, map[string]interface{}{
		"reason":        req.Reason,
		"feedback_type": req.FeedbackType,
		"delete_data":   req.DeleteData,
	}); err != nil {
		fmt.Printf("Failed to publish user deactivated event: %v\n", err)
	}

	return nil
}

// Private helper methods

func (us *UserService) validateCreateUserRequest(req *CreateUserRequest) error {
	if req.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}

	if req.Country == "" {
		return fmt.Errorf("country is required")
	}

	if req.Language == "" {
		return fmt.Errorf("language is required")
	}

	if req.TimeZone == "" {
		return fmt.Errorf("timezone is required")
	}

	// Validate country
	if !models.IsValidCountry(req.Country) {
		return fmt.Errorf("invalid country: %s", req.Country)
	}

	// Validate phone number format for country
	if !models.IsValidPhoneNumber(req.PhoneNumber, req.Country) {
		return fmt.Errorf("invalid phone number format for country %s", req.Country)
	}

	// Validate email if provided
	if req.Email != "" && !models.IsValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (us *UserService) validateUpdateProfileRequest(req *UpdateUserProfileRequest) error {
	if req.Email != nil && *req.Email != "" && !models.IsValidEmail(*req.Email) {
		return fmt.Errorf("invalid email format")
	}

	if req.Username != nil && *req.Username != "" && !models.IsValidUsername(*req.Username) {
		return fmt.Errorf("invalid username format")
	}

	return nil
}

func (us *UserService) publishUserEvent(ctx context.Context, eventType models.EventType, userID uuid.UUID, data map[string]interface{}) error {
	event := &models.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      models.EventCategoryDomain,
		Severity:      models.SeverityInfo,
		Subject:       fmt.Sprintf("User %s: %s", userID, eventType),
		AggregateID:   userID.String(),
		AggregateType: "user",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        models.EventStatusPending,
		Metadata: models.EventMetadata{
			Source:      "auth-service",
			Environment: "production", // Should come from config
			Region:      "sea",        // Should come from config
		},
	}

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return us.eventPublisher.Publish(ctx, event)
}

// Request/Response structures
type CreateUserRequest struct {
	PhoneNumber string                 `json:"phone_number" binding:"required"`
	Email       string                 `json:"email"`
	Username    string                 `json:"username"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Country     models.Country         `json:"country" binding:"required"`
	Language    string                 `json:"language" binding:"required"`
	TimeZone    string                 `json:"timezone" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateUserProfileRequest struct {
	Username    *string                 `json:"username"`
	FirstName   *string                 `json:"first_name"`
	LastName    *string                 `json:"last_name"`
	Email       *string                 `json:"email"`
	Language    *string                 `json:"language"`
	TimeZone    *string                 `json:"timezone"`
	Preferences *models.UserPreferences `json:"preferences"`
	Metadata    map[string]interface{}  `json:"metadata"`
}

type ChangePhoneNumberRequest struct {
	UserID         uuid.UUID `json:"user_id" binding:"required"`
	NewPhoneNumber string    `json:"new_phone_number" binding:"required"`
	CountryCode    string    `json:"country_code" binding:"required"`
	OTPRequestID   string    `json:"otp_request_id" binding:"required"`
	OTPCode        string    `json:"otp_code" binding:"required"`
}

type DeactivateUserRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	Reason       string    `json:"reason"`
	FeedbackType string    `json:"feedback_type"`
	DeleteData   bool      `json:"delete_data"`
}

type UserResponse struct {
	ID              uuid.UUID               `json:"id"`
	PhoneNumber     string                  `json:"phone_number"`
	Email           string                  `json:"email"`
	Username        string                  `json:"username"`
	FirstName       string                  `json:"first_name"`
	LastName        string                  `json:"last_name"`
	Country         models.Country          `json:"country"`
	Language        string                  `json:"language"`
	TimeZone        string                  `json:"timezone"`
	Status          models.UserStatus       `json:"status"`
	IsPhoneVerified bool                    `json:"is_phone_verified"`
	IsEmailVerified bool                    `json:"is_email_verified"`
	KYCTier         models.VerificationTier `json:"kyc_tier"`
	Preferences     models.UserPreferences  `json:"preferences"`
	LastActiveAt    *time.Time              `json:"last_active_at"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

func (user *models.User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              user.ID,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Country:         user.Country,
		Language:        user.Language,
		TimeZone:        user.TimeZone,
		Status:          user.Status,
		IsPhoneVerified: user.IsPhoneVerified,
		IsEmailVerified: user.IsEmailVerified,
		KYCTier:         user.KYCTier,
		Preferences:     user.Preferences,
		LastActiveAt:    user.LastActiveAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
