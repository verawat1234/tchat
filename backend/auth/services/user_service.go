package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/auth/models"
	sharedModels "tchat.dev/shared/models"
	"tchat.dev/shared/utils"
)

type UserRepository interface {
	Create(ctx context.Context, user *sharedModels.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error)
	GetByEmail(ctx context.Context, email string) (*sharedModels.User, error)
	Update(ctx context.Context, user *sharedModels.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters UserFilters, pagination Pagination) ([]*sharedModels.User, int64, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*sharedModels.User, error)
	SearchByUsername(ctx context.Context, username string, limit int) ([]*sharedModels.User, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
}


type UserFilters struct {
	Country     *string          `json:"country,omitempty"`
	Status      *sharedModels.UserStatus       `json:"status,omitempty"`
	KYCTier     *models.VerificationTier `json:"kyc_tier,omitempty"`
	CreatedFrom *time.Time                     `json:"created_from,omitempty"`
	CreatedTo   *time.Time                     `json:"created_to,omitempty"`
	Search      string                         `json:"search,omitempty"`
}

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
	Order    string `json:"order"` // asc, desc
}

type UserStats struct {
	TotalUsers           int64                                     `json:"total_users"`
	ActiveUsers          int64                                     `json:"active_users"`
	VerifiedUsers        int64                                     `json:"verified_users"`
	NewUsersToday        int64                                     `json:"new_users_today"`
	NewUsersThisWeek     int64                                     `json:"new_users_this_week"`
	NewUsersThisMonth    int64                                     `json:"new_users_this_month"`
	UsersByCountry       map[string]int64            `json:"users_by_country"`
	UsersByKYCTier       map[models.VerificationTier]int64   `json:"users_by_kyc_tier"`
	AverageSessionLength time.Duration                             `json:"average_session_length"`
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

func (us *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*sharedModels.User, error) {
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

	// Create user using shared model with correct field mapping
	user := &sharedModels.User{
		ID:          uuid.New(),
		PhoneNumber: req.PhoneNumber,
		CountryCode: req.Country,
		Status:      string(sharedModels.UserStatusActive), // Default to active
		DisplayName: req.FirstName + " " + req.LastName,
		Locale:      req.Language,
		Timezone:    req.TimeZone,
		PhoneVerified: false,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := us.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Publish user registration event when event types are defined
	// if err := us.publishUserEvent(ctx, models.EventTypeUserRegistered, user.ID, map[string]interface{}{
	//	"phone_number": user.PhoneNumber,
	//	"country":      user.CountryCode,
	//	"language":     user.Locale,
	// }); err != nil {
	//	// Log error but don't fail the operation
	//	// In production, you might want to use a more robust event publishing mechanism
	//	fmt.Printf("Failed to publish user registration event: %v\n", err)
	// }

	return user, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*sharedModels.User, error) {
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

func (us *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error) {
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

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*sharedModels.User, error) {
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

func (us *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *UpdateUserProfileRequest) (*sharedModels.User, error) {
	// Get existing user
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check permissions - user must be active
	if !user.IsActive() {
		return nil, fmt.Errorf("user cannot update profile in current status: %s", user.Status)
	}

	// Validate update request
	if err := us.validateUpdateProfileRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.FirstName != nil || req.LastName != nil {
		// Combine first and last name into the DisplayName field in Profile
		newName := user.DisplayName
		if req.FirstName != nil && req.LastName != nil {
			newName = *req.FirstName + " " + *req.LastName
		} else if req.FirstName != nil {
			// Split existing name to get last name part
			parts := strings.Split(user.DisplayName, " ")
			lastName := ""
			if len(parts) > 1 {
				lastName = strings.Join(parts[1:], " ")
			}
			newName = *req.FirstName + " " + lastName
		} else if req.LastName != nil {
			// Split existing name to get first name part
			parts := strings.Split(user.DisplayName, " ")
			firstName := parts[0]
			newName = firstName + " " + *req.LastName
		}
		if newName != user.DisplayName {
			changes["display_name"] = map[string]string{"from": user.DisplayName, "to": newName}
			user.DisplayName = newName
		}
	}

	// Note: Email field is not currently available in the shared User model
	// This functionality would need to be added if email support is required
	if req.Email != nil {
		changes["email_note"] = "Email field not supported in current shared User model"
	}

	if req.Language != nil && *req.Language != user.Locale {
		changes["locale"] = map[string]string{"from": user.Locale, "to": *req.Language}
		user.Locale = *req.Language
	}

	if req.TimeZone != nil && *req.TimeZone != user.Timezone {
		changes["timezone"] = map[string]string{"from": user.Timezone, "to": *req.TimeZone}
		user.Timezone = *req.TimeZone
	}

	// DISABLED: Preferences field not available in current User model
	// if req.Preferences != nil {
	// 	// Preferences field not available in current User model
	// 	// Would need to add this field to models.User if preferences support is required
	// 	changes["preferences_note"] = "Preferences field not supported in current User model"
	// }

	if req.Metadata != nil {
		// Metadata field not available in current User model
		// Would need to add this field to models.User if metadata support is required
		changes["metadata_note"] = "Metadata field not supported in current User model"
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
		if err := us.publishUserEvent(ctx, sharedModels.EventTypeUserProfileUpdated, user.ID, map[string]interface{}{
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

	if user.PhoneVerified {
		return fmt.Errorf("phone number already verified")
	}

	user.PhoneVerified = true
	user.UpdatedAt = time.Now()

	// Update status if was pending - shared model doesn't have pending status
	// so we'll leave it as is

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

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user email verification status: %w", err)
	}

	return nil
}

func (us *UserService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status sharedModels.UserStatus, reason string) error {
	user, err := us.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	oldStatus := user.Status

	// Basic status transition validation - allow transitions to active/suspended/deleted
	if status != sharedModels.UserStatusActive && status != sharedModels.UserStatusSuspended && status != sharedModels.UserStatusDeleted {
		return fmt.Errorf("invalid status transition to: %s", status)
	}

	user.Status = string(status)
	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Publish status change event
	if err := us.publishUserEvent(ctx, sharedModels.EventTypeUserProfileUpdated, user.ID, map[string]interface{}{
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

	// Convert VerificationTier to KYCTier and validate upgrade
	kycTier := sharedModels.KYCTier(tier)
	if err := user.UpdateKYCTier(int(kycTier)); err != nil {
		return fmt.Errorf("invalid KYC tier upgrade: %w", err)
	}

	user.UpdatedAt = time.Now()

	if err := us.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user KYC tier: %w", err)
	}

	// Publish KYC tier update event
	if err := us.publishUserEvent(ctx, sharedModels.EventTypeUserKYCVerified, user.ID, map[string]interface{}{
		"kyc_tier_change": map[string]interface{}{
			"from": oldTier,
			"to":   tier,
		},
	}); err != nil {
		fmt.Printf("Failed to publish KYC tier update event: %v\n", err)
	}

	return nil
}

func (us *UserService) ListUsers(ctx context.Context, filters UserFilters, pagination Pagination) ([]*sharedModels.User, int64, error) {
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

func (us *UserService) SearchUsers(ctx context.Context, query string, limit int) ([]*sharedModels.User, error) {
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

	// Check if user can be deleted - must not be already deleted
	if user.Status == string(sharedModels.UserStatusDeleted) {
		return fmt.Errorf("user is already deleted")
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

func (us *UserService) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*sharedModels.User, error) {
	if len(userIDs) == 0 {
		return []*sharedModels.User{}, nil
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
	user.PhoneNumber = req.NewPhoneNumber
	user.CountryCode = req.CountryCode
	user.UpdatedAt = time.Now()

	err = us.userRepo.Update(ctx, user)
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

	if user.Status == string(sharedModels.UserStatusSuspended) {
		return fmt.Errorf("user is already deactivated")
	}

	// Update user status to suspended (shared model uses suspended instead of inactive)
	err = us.UpdateUserStatus(ctx, req.UserID, sharedModels.UserStatusSuspended, req.Reason)
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
	if !sharedModels.IsValidSEACountry(req.Country) {
		return fmt.Errorf("invalid country: %s", req.Country)
	}

	// Validate phone number format for country
	if !utils.IsValidPhoneNumber(req.PhoneNumber, req.Country) {
		return fmt.Errorf("invalid phone number format for country %s", req.Country)
	}

	// Validate email if provided
	if req.Email != "" && !utils.IsValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (us *UserService) validateUpdateProfileRequest(req *UpdateUserProfileRequest) error {
	if req.Email != nil && *req.Email != "" && !utils.IsValidEmail(*req.Email) {
		return fmt.Errorf("invalid email format")
	}

	if req.Username != nil && *req.Username != "" && !utils.IsValidUsername(*req.Username) {
		return fmt.Errorf("invalid username format")
	}

	return nil
}

func (us *UserService) publishUserEvent(ctx context.Context, eventType sharedModels.EventType, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("User %s: %s", userID, eventType),
		AggregateID:   userID.String(),
		AggregateType: "user",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
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
	Country     string                 `json:"country" binding:"required"`
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
	// Preferences *sharedModels.UserPreferences `json:"preferences"` // DISABLED: UserPreferences struct is disabled
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
	Country         string                  `json:"country"`
	Language        string                  `json:"language"`
	TimeZone        string                  `json:"timezone"`
	Status          string                  `json:"status"`
	IsPhoneVerified bool                    `json:"is_phone_verified"`
	IsEmailVerified bool                    `json:"is_email_verified"`
	KYCTier         int                     `json:"kyc_tier"`
	// Preferences     sharedModels.UserPreferences  `json:"preferences"` // DISABLED: UserPreferences struct is disabled
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

func ToUserResponse(user *sharedModels.User) *UserResponse {
	// Extract first and last name from display name if possible
	nameParts := strings.Split(user.DisplayName, " ")
	firstName := ""
	lastName := ""
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	return &UserResponse{
		ID:              user.ID,
		PhoneNumber:     user.PhoneNumber,
		Email:           "", // Email field needs to be added to shared model if needed
		FirstName:       firstName,
		LastName:        lastName,
		Country:         user.CountryCode,
		Language:        user.Locale,
		TimeZone:        user.Timezone,
		Status:          string(user.Status),
		IsPhoneVerified: user.PhoneVerified,
		IsEmailVerified: user.EmailVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
