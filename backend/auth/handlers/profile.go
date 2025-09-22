package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/auth/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

type ProfileHandler struct {
	userService *services.UserService
	kycService  *services.KYCService
	validator   *validator.Validate
}

func NewProfileHandler(userService *services.UserService, kycService *services.KYCService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
		kycService:  kycService,
		validator:   validator.New(),
	}
}

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name,omitempty" validate:"omitempty,min=2,max=50" example:"John Doe"`
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,min=1,max=30" example:"John"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,min=1,max=30" example:"Doe"`
	Avatar      string `json:"avatar,omitempty" validate:"omitempty,url" example:"https://cdn.tchat.com/avatars/user.jpg"`
	DateOfBirth string `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02" example:"1990-01-15"`
	Gender      string `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say" example:"male"`
	Language    string `json:"language,omitempty" validate:"omitempty,len=2" example:"en"`
	Timezone    string `json:"timezone,omitempty" validate:"omitempty,timezone" example:"Asia/Bangkok"`
	Bio         string `json:"bio,omitempty" validate:"omitempty,max=200" example:"Software developer from Bangkok"`
}

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID              uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PhoneNumber     string    `json:"phone_number" example:"+66812345678"`
	CountryCode     string    `json:"country_code" example:"TH"`
	DisplayName     string    `json:"display_name" example:"John Doe"`
	FirstName       string    `json:"first_name,omitempty" example:"John"`
	LastName        string    `json:"last_name,omitempty" example:"Doe"`
	Avatar          string    `json:"avatar,omitempty" example:"https://cdn.tchat.com/avatars/user.jpg"`
	DateOfBirth     string    `json:"date_of_birth,omitempty" example:"1990-01-15"`
	Gender          string    `json:"gender,omitempty" example:"male"`
	Language        string    `json:"language" example:"en"`
	Timezone        string    `json:"timezone" example:"Asia/Bangkok"`
	Bio             string    `json:"bio,omitempty" example:"Software developer from Bangkok"`
	KYCStatus       string    `json:"kyc_status" example:"verified"`
	KYCTier         string    `json:"kyc_tier" example:"tier2"`
	IsActive        bool      `json:"is_active" example:"true"`
	IsVerified      bool      `json:"is_verified" example:"true"`
	JoinedAt        string    `json:"joined_at" example:"2024-01-15T10:30:00Z"`
	LastActiveAt    string    `json:"last_active_at" example:"2024-01-20T15:45:00Z"`
	ProfileComplete float32   `json:"profile_complete" example:"85.5"`
}

// ChangePhoneRequest represents the request to change phone number
type ChangePhoneRequest struct {
	NewPhoneNumber string `json:"new_phone_number" validate:"required,e164" example:"+66887654321"`
	CountryCode    string `json:"country_code" validate:"required,len=2" example:"TH"`
	OTPRequestID   string `json:"otp_request_id" validate:"required" example:"req_987654321"`
	OTPCode        string `json:"otp_code" validate:"required,len=6,numeric" example:"654321"`
}

// DeactivateAccountRequest represents the request to deactivate account
type DeactivateAccountRequest struct {
	Reason       string `json:"reason,omitempty" validate:"omitempty,max=500" example:"Taking a break from social media"`
	FeedbackType string `json:"feedback_type,omitempty" validate:"omitempty,oneof=temporary permanent feedback" example:"temporary"`
	DeleteData   bool   `json:"delete_data" example:"false"`
	Confirmation string `json:"confirmation" validate:"required,eqfield=DEACTIVATE" example:"DEACTIVATE"`
}

// ProfileStatsResponse represents user profile statistics
type ProfileStatsResponse struct {
	TotalDialogs     int     `json:"total_dialogs" example:"25"`
	TotalMessages    int     `json:"total_messages" example:"1540"`
	TotalContacts    int     `json:"total_contacts" example:"89"`
	WalletBalance    float64 `json:"wallet_balance" example:"1250.75"`
	TransactionCount int     `json:"transaction_count" example:"47"`
	JoinedDaysAgo    int     `json:"joined_days_ago" example:"156"`
	LastActiveHours  int     `json:"last_active_hours" example:"2"`
}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=ProfileResponse}
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get user details
	user, err := h.userService.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "User lookup failed", "Failed to retrieve user information.")
		return
	}

	// Build profile response
	profile := ProfileResponse{
		ID:              user.ID,
		PhoneNumber:     maskPhoneNumber(user.PhoneNumber),
		CountryCode:     user.CountryCode,
		DisplayName:     user.DisplayName,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Avatar:          user.Avatar,
		Language:        user.Language,
		Timezone:        user.Timezone,
		Bio:             user.Bio,
		KYCStatus:       string(user.KYCStatus),
		KYCTier:         string(user.KYCTier),
		IsActive:        user.IsActive,
		IsVerified:      user.IsVerified,
		JoinedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastActiveAt:    user.LastActiveAt.Format("2006-01-02T15:04:05Z"),
		ProfileComplete: h.calculateProfileCompleteness(user),
	}

	// Add date of birth if present
	if user.DateOfBirth != nil {
		profile.DateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	// Add gender if present
	if user.Gender != "" {
		profile.Gender = string(user.Gender)
	}

	responses.DataResponse(c, profile)
}

// @Summary Update user profile
// @Description Update current user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update request"
// @Success 200 {object} responses.DataResponse{data=ProfileResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.UpdateUserRequest{
		UserID:      userUUID,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Avatar:      req.Avatar,
		Language:    req.Language,
		Timezone:    req.Timezone,
		Bio:         req.Bio,
		Gender:      req.Gender,
	}

	// Parse date of birth if provided
	if req.DateOfBirth != "" {
		serviceReq.DateOfBirth = req.DateOfBirth
	}

	// Update user profile
	user, err := h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			responses.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}
		responses.ErrorResponse(c, http.StatusInternalServerError, "Profile update failed", "Failed to update profile.")
		return
	}

	// Build updated profile response
	profile := ProfileResponse{
		ID:              user.ID,
		PhoneNumber:     maskPhoneNumber(user.PhoneNumber),
		CountryCode:     user.CountryCode,
		DisplayName:     user.DisplayName,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Avatar:          user.Avatar,
		Language:        user.Language,
		Timezone:        user.Timezone,
		Bio:             user.Bio,
		KYCStatus:       string(user.KYCStatus),
		KYCTier:         string(user.KYCTier),
		IsActive:        user.IsActive,
		IsVerified:      user.IsVerified,
		JoinedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastActiveAt:    user.LastActiveAt.Format("2006-01-02T15:04:05Z"),
		ProfileComplete: h.calculateProfileCompleteness(user),
	}

	// Add optional fields
	if user.DateOfBirth != nil {
		profile.DateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}
	if user.Gender != "" {
		profile.Gender = string(user.Gender)
	}

	// Log profile update
	middleware.LogInfo(c, "Profile updated", gin.H{
		"user_id":        userUUID,
		"fields_updated": getUpdatedFields(req),
	})

	responses.DataResponse(c, profile)
}

// @Summary Get profile statistics
// @Description Get user profile statistics and activity metrics
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=ProfileStatsResponse}
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile/stats [get]
func (h *ProfileHandler) GetProfileStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get user statistics
	stats, err := h.userService.GetUserStats(c.Request.Context(), userUUID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Stats lookup failed", "Failed to retrieve user statistics.")
		return
	}

	// Build stats response
	profileStats := ProfileStatsResponse{
		TotalDialogs:     stats.TotalDialogs,
		TotalMessages:    stats.TotalMessages,
		TotalContacts:    stats.TotalContacts,
		WalletBalance:    stats.WalletBalance,
		TransactionCount: stats.TransactionCount,
		JoinedDaysAgo:    stats.JoinedDaysAgo,
		LastActiveHours:  stats.LastActiveHours,
	}

	responses.DataResponse(c, profileStats)
}

// @Summary Change phone number
// @Description Change user's phone number with OTP verification
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePhoneRequest true "Phone number change request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 409 {object} responses.ErrorResponse "Phone number already exists"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile/change-phone [post]
func (h *ProfileHandler) ChangePhoneNumber(c *gin.Context) {
	var req ChangePhoneRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Normalize phone number and country code
	req.NewPhoneNumber = strings.TrimSpace(req.NewPhoneNumber)
	req.CountryCode = strings.ToUpper(strings.TrimSpace(req.CountryCode))

	// Build service request
	serviceReq := &services.ChangePhoneNumberRequest{
		UserID:         userUUID,
		NewPhoneNumber: req.NewPhoneNumber,
		CountryCode:    req.CountryCode,
		OTPRequestID:   req.OTPRequestID,
		OTPCode:        req.OTPCode,
	}

	// Change phone number
	err := h.userService.ChangePhoneNumber(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "phone number already exists"):
			responses.ErrorResponse(c, http.StatusConflict, "Phone number already exists", "This phone number is already registered to another account.")
			return
		case strings.Contains(err.Error(), "invalid OTP"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Invalid OTP", "The OTP code is incorrect or has expired.")
			return
		case strings.Contains(err.Error(), "OTP request not found"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Invalid OTP request", "The OTP request is invalid or has expired.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Phone number change failed", "Failed to change phone number.")
			return
		}
	}

	// Log phone number change
	middleware.LogInfo(c, "Phone number changed", gin.H{
		"user_id":          userUUID,
		"new_phone_number": maskPhoneNumber(req.NewPhoneNumber),
		"country_code":     req.CountryCode,
	})

	responses.SuccessMessageResponse(c, "Phone number changed successfully")
}

// @Summary Deactivate account
// @Description Deactivate user account (can be reactivated later)
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DeactivateAccountRequest true "Account deactivation request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile/deactivate [post]
func (h *ProfileHandler) DeactivateAccount(c *gin.Context) {
	var req DeactivateAccountRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.DeactivateUserRequest{
		UserID:       userUUID,
		Reason:       req.Reason,
		FeedbackType: req.FeedbackType,
		DeleteData:   req.DeleteData,
	}

	// Deactivate account
	err := h.userService.DeactivateUser(c.Request.Context(), serviceReq)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Account deactivation failed", "Failed to deactivate account.")
		return
	}

	// Log account deactivation
	middleware.LogInfo(c, "Account deactivated", gin.H{
		"user_id":       userUUID,
		"reason":        req.Reason,
		"feedback_type": req.FeedbackType,
		"delete_data":   req.DeleteData,
	})

	responses.SuccessMessageResponse(c, "Account deactivated successfully")
}

// @Summary Get KYC status
// @Description Get current KYC verification status and requirements
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/profile/kyc [get]
func (h *ProfileHandler) GetKYCStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get KYC status
	kycStatus, err := h.kycService.GetKYCStatus(c.Request.Context(), userUUID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "KYC status lookup failed", "Failed to retrieve KYC status.")
		return
	}

	responses.DataResponse(c, kycStatus)
}

// Helper functions

func (h *ProfileHandler) calculateProfileCompleteness(user interface{}) float32 {
	// This is a simplified calculation
	// In a real implementation, you'd check various profile fields
	completeness := float32(40) // Base score for having an account

	// Add scores for various fields (simplified)
	// You would check actual user fields here
	completeness += 20 // Display name
	completeness += 15 // Avatar
	completeness += 10 // Bio
	completeness += 15 // Other profile fields

	if completeness > 100 {
		completeness = 100
	}

	return completeness
}

func getUpdatedFields(req UpdateProfileRequest) []string {
	var fields []string

	if req.DisplayName != "" {
		fields = append(fields, "display_name")
	}
	if req.FirstName != "" {
		fields = append(fields, "first_name")
	}
	if req.LastName != "" {
		fields = append(fields, "last_name")
	}
	if req.Avatar != "" {
		fields = append(fields, "avatar")
	}
	if req.DateOfBirth != "" {
		fields = append(fields, "date_of_birth")
	}
	if req.Gender != "" {
		fields = append(fields, "gender")
	}
	if req.Language != "" {
		fields = append(fields, "language")
	}
	if req.Timezone != "" {
		fields = append(fields, "timezone")
	}
	if req.Bio != "" {
		fields = append(fields, "bio")
	}

	return fields
}

func maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 6 {
		return phoneNumber
	}

	// Show first 3 and last 3 characters, mask the middle
	if len(phoneNumber) <= 9 {
		return phoneNumber[:3] + "***" + phoneNumber[len(phoneNumber)-3:]
	}

	// For longer numbers, show first 4 and last 4
	return phoneNumber[:4] + "****" + phoneNumber[len(phoneNumber)-4:]
}

// RegisterProfileRoutes registers all profile-related routes
func RegisterProfileRoutes(router *gin.RouterGroup, userService *services.UserService, kycService *services.KYCService) {
	handler := NewProfileHandler(userService, kycService)

	// Protected profile routes
	profile := router.Group("/users/profile")
	profile.Use(middleware.AuthRequired())
	{
		profile.GET("", handler.GetProfile)
		profile.PUT("", handler.UpdateProfile)
		profile.GET("/stats", handler.GetProfileStats)
		profile.POST("/change-phone", handler.ChangePhoneNumber)
		profile.POST("/deactivate", handler.DeactivateAccount)
		profile.GET("/kyc", handler.GetKYCStatus)
	}
}

// RegisterProfileRoutesWithMiddleware registers profile routes with custom middleware
func RegisterProfileRoutesWithMiddleware(
	router *gin.RouterGroup,
	userService *services.UserService,
	kycService *services.KYCService,
	middlewares ...gin.HandlerFunc,
) {
	handler := NewProfileHandler(userService, kycService)

	// Protected profile routes with middleware
	profile := router.Group("/users/profile")
	allMiddlewares := append(middlewares, middleware.AuthRequired())
	profile.Use(allMiddlewares...)
	{
		profile.GET("", handler.GetProfile)
		profile.PUT("", handler.UpdateProfile)
		profile.GET("/stats", handler.GetProfileStats)
		profile.POST("/change-phone", handler.ChangePhoneNumber)
		profile.POST("/deactivate", handler.DeactivateAccount)
		profile.GET("/kyc", handler.GetKYCStatus)
	}
}
