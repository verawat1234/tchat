package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/auth/models"
	"tchat.dev/auth/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	"tchat.dev/shared/utils"
	sharedModels "tchat.dev/shared/models"
)

type AuthHandler struct {
	authService    *services.AuthService
	sessionService *services.SessionService
	userService    *services.UserService
	validator      *validator.Validate
}

func NewAuthHandler(
	authService *services.AuthService,
	sessionService *services.SessionService,
	userService *services.UserService,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
		userService:    userService,
		validator:      validator.New(),
	}
}

// VerifyOTPRequest represents the request to verify OTP
type VerifyOTPRequest struct {
	RequestID   string `json:"request_id" validate:"required" example:"req_123456789"`
	Code        string `json:"code" validate:"required,len=6,numeric" example:"123456"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+66987654321"`
	DeviceInfo *DeviceInfo `json:"device_info,omitempty"`
}

type DeviceInfo struct {
	Platform    string `json:"platform" validate:"required,oneof=ios android web" example:"ios"`
	DeviceModel string `json:"device_model,omitempty" example:"iPhone 14 Pro"`
	OSVersion   string `json:"os_version,omitempty" example:"16.4"`
	AppVersion  string `json:"app_version,omitempty" example:"1.0.0"`
	DeviceID    string `json:"device_id,omitempty" example:"ABCD-1234-EFGH-5678"`
	PushToken   string `json:"push_token,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	Timezone    string `json:"timezone,omitempty" example:"Asia/Bangkok"`
	Language    string `json:"language,omitempty" example:"en"`
}

// VerifyOTPResponse represents the response after successful OTP verification
type VerifyOTPResponse struct {
	Success      bool      `json:"success" example:"true"`
	Message      string    `json:"message" example:"Authentication successful"`
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresIn    int       `json:"expires_in" example:"900"`
	User         UserInfo  `json:"user"`
	Session      SessionInfo `json:"session"`
}

type UserInfo struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PhoneNumber string    `json:"phone_number" example:"+66812345678"`
	CountryCode string    `json:"country_code" example:"TH"`
	DisplayName string    `json:"display_name,omitempty" example:"John Doe"`
	Avatar      string    `json:"avatar,omitempty" example:"https://cdn.tchat.com/avatars/user.jpg"`
	KYCStatus   string    `json:"kyc_status" example:"verified"`
	KYCTier     string    `json:"kyc_tier" example:"tier2"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SessionInfo struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174001"`
	DeviceInfo string   `json:"device_info" example:"iPhone 14 Pro"`
	IPAddress string    `json:"ip_address" example:"192.168.1.1"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents the request to refresh access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
	DeviceInfo   *DeviceInfo `json:"device_info,omitempty"`
}

// RequestOTPRequest represents the request to initiate OTP
type RequestOTPRequest struct {
	PhoneNumber string      `json:"phone_number" validate:"required" example:"+66812345678"`
	CountryCode string      `json:"country_code" validate:"required,len=2" example:"TH"`
	DeviceInfo  *DeviceInfo `json:"device_info,omitempty"`
}

// RequestOTPResponse represents the response after requesting OTP
type RequestOTPResponse struct {
	Success   bool   `json:"success" example:"true"`
	Message   string `json:"message" example:"OTP sent successfully"`
	RequestID string `json:"request_id" example:"req_123456789"`
	ExpiresIn int    `json:"expires_in" example:"300"`
}

// LogoutRequest represents the request to logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
	LogoutAll    bool   `json:"logout_all,omitempty" example:"false"`
}

// @Summary Request OTP
// @Description Request OTP for phone number authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "OTP request"
// @Success 200 {object} RequestOTPResponse "OTP sent successfully"
// @Failure 400 {object} responses.SendErrorResponse "Invalid request"
// @Failure 500 {object} responses.SendErrorResponse "Internal server error"
// @Router /auth/otp/request [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// SECURITY FIX: Enhanced phone number validation for Southeast Asian formats
	if err := h.validatePhoneNumber(req.PhoneNumber, req.CountryCode); err != nil {
		// Log security event for invalid phone number attempts
		middleware.LogWarning(c, "Invalid phone number format in OTP request", gin.H{
			"phone_number_masked": maskPhoneNumber(req.PhoneNumber),
			"country_code":        req.CountryCode,
			"error":               err.Error(),
			"ip_address":          c.ClientIP(),
			"user_agent":          c.GetHeader("User-Agent"),
		})
		responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_PHONE_NUMBER", err.Error())
		return
	}

	// Call auth service to send OTP
	serviceReq := &services.SendOTPRequest{
		PhoneNumber: req.PhoneNumber,
		Type:        services.OTPTypeLogin,
		Language:    "en", // Default to English
		UserAgent:   c.GetHeader("User-Agent"),
		IPAddress:   c.ClientIP(),
		Metadata:    map[string]interface{}{
			"country_code": req.CountryCode,
		},
	}

	sendOTPResponse, err := h.authService.SendOTP(c.Request.Context(), serviceReq)
	if err != nil {
		log.Printf("SendOTP error: %v", err)
		// Handle specific error types
		switch {
		case strings.Contains(err.Error(), "rate limit"):
			responses.SendErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", err.Error())
			return
		case strings.Contains(err.Error(), "user not found"):
			responses.SendErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		default:
			responses.SendErrorResponse(c, http.StatusInternalServerError, "OTP_SEND_FAILED", "Failed to send OTP")
			return
		}
	}

	// Use the actual OTP UUID as the request_id
	requestID := sendOTPResponse.OTPID.String()

	// Log successful OTP request
	middleware.LogInfo(c, "OTP request initiated", gin.H{
		"phone_number_masked": maskPhoneNumber(req.PhoneNumber),
		"country_code":        req.CountryCode,
		"request_id":          requestID,
		"ip_address":          c.ClientIP(),
	})

	response := RequestOTPResponse{
		Success:   true,
		Message:   "OTP sent successfully",
		RequestID: requestID,
		ExpiresIn: int(time.Until(sendOTPResponse.ExpiresAt).Seconds()),
	}

	responses.SendSuccessResponse(c, response)
}

// @Summary Verify OTP
// @Description Verify OTP code and complete authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification request"
// @Success 200 {object} VerifyOTPResponse
// @Failure 400 {object} responses.SendErrorResponse
// @Failure 401 {object} responses.SendErrorResponse "Invalid OTP"
// @Failure 429 {object} responses.SendErrorResponse "Rate limit exceeded"
// @Failure 500 {object} responses.SendErrorResponse
// @Router /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// SECURITY FIX: Enhanced phone number validation for verify OTP
	phoneNumber := req.PhoneNumber
	if phoneNumber == "" {
		phoneNumber = extractPhoneFromRequestID(req.RequestID)
	}

	// Extract country code from phone number for validation
	countryCode := extractCountryFromPhone(phoneNumber)
	if err := h.validatePhoneNumber(phoneNumber, countryCode); err != nil {
		// Log security event for invalid phone number in OTP verification
		middleware.LogWarning(c, "Invalid phone number format in OTP verification", gin.H{
			"phone_number_masked": maskPhoneNumber(phoneNumber),
			"request_id":          req.RequestID,
			"error":               err.Error(),
			"ip_address":          c.ClientIP(),
			"user_agent":          c.GetHeader("User-Agent"),
		})
		responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_PHONE_NUMBER", err.Error())
		return
	}

	// Extract IP address if not provided in device info
	if req.DeviceInfo != nil && req.DeviceInfo.IPAddress == "" {
		req.DeviceInfo.IPAddress = c.ClientIP()
	}

	// Extract User-Agent if not provided
	if req.DeviceInfo != nil && req.DeviceInfo.UserAgent == "" {
		req.DeviceInfo.UserAgent = c.GetHeader("User-Agent")
	}

	// Determine phone number - use provided phone number if available, otherwise extract from request ID
	finalPhoneNumber := req.PhoneNumber
	if finalPhoneNumber == "" {
		finalPhoneNumber = extractPhoneFromRequestID(req.RequestID)
	}

	// Call auth service to verify OTP
	serviceReq := &services.VerifyOTPRequest{
		PhoneNumber: finalPhoneNumber,
		Code:        req.Code,
		UserAgent:   c.GetHeader("User-Agent"),
		IPAddress:   c.ClientIP(),
	}

	verifyResponse, err := h.authService.VerifyOTP(c.Request.Context(), serviceReq)
	if err != nil {
		log.Printf("VerifyOTP error: %v", err)
		// Handle specific error types
		switch {
		case strings.Contains(err.Error(), "invalid code"):
			responses.SendErrorResponse(c, http.StatusUnauthorized, "INVALID_OTP", "The OTP code is incorrect or has expired.")
			return
		case strings.Contains(err.Error(), "expired"):
			responses.SendErrorResponse(c, http.StatusUnauthorized, "OTP_EXPIRED", "The OTP code has expired. Please request a new one.")
			return
		case strings.Contains(err.Error(), "max attempts"):
			responses.SendErrorResponse(c, http.StatusTooManyRequests, "MAX_ATTEMPTS_EXCEEDED", "Too many failed attempts. Please request a new OTP.")
			return
		case strings.Contains(err.Error(), "request not found"):
			responses.SendErrorResponse(c, http.StatusNotFound, "REQUEST_NOT_FOUND", "Invalid or expired OTP request.")
			return
		default:
			responses.SendErrorResponse(c, http.StatusInternalServerError, "VERIFICATION_FAILED", "Internal server error occurred.")
			return
		}
	}

	// If VerifyOTP service handled everything (e.g., test mode), return the response directly
	if verifyResponse != nil && verifyResponse.Success && verifyResponse.Session != nil {
		response := gin.H{
			"success": true,
			"message": "OTP verified successfully",
			"data": gin.H{
				"user": gin.H{
					"id":          verifyResponse.User.ID,
					"phone":       verifyResponse.User.PhoneNumber,
					"country":     verifyResponse.User.Country,
					"language":    verifyResponse.User.Language,
					"timezone":    verifyResponse.User.TimeZone,
					"status":      verifyResponse.User.Status,
				},
				"access_token":  verifyResponse.Session.AccessToken,
				"refresh_token": verifyResponse.Session.RefreshToken,
				"expires_at":    verifyResponse.Session.ExpiresAt,
			},
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Get or create user based on phone number from request
	user, err := h.userService.GetUserByPhoneNumber(c.Request.Context(), finalPhoneNumber)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "USER_LOOKUP_FAILED", "Failed to retrieve user information.")
		return
	}

	// Create new user if not exists
	if user == nil {
		createUserReq := &services.CreateUserRequest{
			PhoneNumber: finalPhoneNumber,
			Country:     extractCountryFromPhone(finalPhoneNumber),
			Language:    getLanguageFromDeviceInfo(req.DeviceInfo),
			TimeZone:    "UTC",
		}

		user, err = h.userService.CreateUser(c.Request.Context(), createUserReq)
		if err != nil {
			responses.SendErrorResponse(c, http.StatusInternalServerError, "USER_CREATION_FAILED", "Failed to create user account.")
			return
		}
	}

	// Create session
	sessionReq := &services.CreateSessionRequest{
		UserID:     user.ID,
		DeviceInfo: buildSessionDeviceInfoMap(req.DeviceInfo),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), sessionReq)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "SESSION_CREATION_FAILED", "Failed to create user session.")
		return
	}

	// Build response
	response := VerifyOTPResponse{
		Success:      true,
		Message:      "Authentication successful",
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		User: UserInfo{
			ID:          user.ID,
			PhoneNumber: maskPhoneNumber(user.PhoneNumber),
			CountryCode: user.CountryCode,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			KYCStatus:   getKYCStatusString(user),
			KYCTier:     fmt.Sprintf("%d", user.KYCTier),
			IsActive:    user.IsActive(),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Session: SessionInfo{
			ID:         session.ID,
			DeviceInfo: convertDeviceInfoToString(session.DeviceInfo),
			IPAddress:  session.IPAddress,
			CreatedAt:  session.CreatedAt,
			ExpiresAt:  session.ExpiresAt,
		},
	}

	// Log successful authentication
	middleware.LogInfo(c, "Authentication successful", gin.H{
		"user_id":      user.ID,
		"phone_number": maskPhoneNumber(user.PhoneNumber),
		"session_id":   session.ID,
		"device_info":  convertDeviceInfoToString(session.DeviceInfo),
	})

	responses.SendSuccessResponse(c, response)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Token refresh request"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.SendErrorResponse
// @Failure 401 {object} responses.SendErrorResponse "Invalid refresh token"
// @Failure 500 {object} responses.SendErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get session by refresh token
	session, err := h.sessionService.GetSessionByRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.SendErrorResponse(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "The refresh token is invalid or has expired.")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "TOKEN_REFRESH_FAILED", "Internal server error occurred.")
		return
	}

	// Refresh tokens
	newSession, err := h.sessionService.RefreshTokens(c.Request.Context(), session.ID)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Token refresh failed", "Failed to refresh tokens.")
		return
	}

	// Update device info if provided
	if req.DeviceInfo != nil {
		updateReq := &services.UpdateSessionRequest{
			SessionID:  newSession.ID,
			DeviceInfo: buildSessionDeviceInfoMap(req.DeviceInfo),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		newSession, err = h.sessionService.UpdateSession(c.Request.Context(), updateReq)
		if err != nil {
			// Log warning but continue with response
			log.Printf("Failed to update session device info for session %s: %v", newSession.ID, err)
		}
	}

	// Build response data
	data := gin.H{
		"access_token":  newSession.AccessToken,
		"refresh_token": newSession.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(newSession.ExpiresAt.Sub(time.Now()).Seconds()),
		"session": gin.H{
			"id":         newSession.ID,
			"expires_at": newSession.RefreshExpiresAt,
		},
	}

	// Log successful token refresh
	middleware.LogInfo(c, "Token refreshed", gin.H{
		"session_id": newSession.ID,
		"user_id":    newSession.UserID,
	})

	responses.SendDataResponse(c, data)
}

// @Summary Logout
// @Description Logout and invalidate session(s)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Logout request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.SendErrorResponse
// @Failure 401 {object} responses.SendErrorResponse "Unauthorized"
// @Failure 500 {object} responses.SendErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body for logout
		req = LogoutRequest{}
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		responses.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get session ID from context
	sessionID, exists := c.Get("session_id")
	var sessionUUID uuid.UUID
	if exists {
		sessionUUID, ok = sessionID.(uuid.UUID)
		if !ok {
			responses.SendErrorResponse(c, http.StatusInternalServerError, "Invalid session context", "Invalid session context.")
			return
		}
	}

	// Logout all sessions or just current session
	if req.LogoutAll {
		err := h.sessionService.InvalidateUserSessions(c.Request.Context(), userUUID)
		if err != nil {
			responses.SendErrorResponse(c, http.StatusInternalServerError, "Logout failed", "Failed to logout from all sessions.")
			return
		}

		// Log logout all sessions
		middleware.LogInfo(c, "Logout all sessions", gin.H{
			"user_id": userUUID,
		})

		responses.SuccessMessageResponse(c, "Logged out from all sessions successfully")
	} else {
		// Use refresh token if provided, otherwise use current session
		if req.RefreshToken != "" {
			session, err := h.sessionService.GetSessionByRefreshToken(c.Request.Context(), req.RefreshToken)
			if err != nil {
				responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid refresh token", "The provided refresh token is invalid.")
				return
			}
			sessionUUID = session.ID
		}

		if sessionUUID != uuid.Nil {
			err := h.sessionService.InvalidateSession(c.Request.Context(), sessionUUID)
			if err != nil {
				responses.SendErrorResponse(c, http.StatusInternalServerError, "Logout failed", "Failed to logout session.")
				return
			}
		}

		// Log logout single session
		middleware.LogInfo(c, "Logout session", gin.H{
			"user_id":    userUUID,
			"session_id": sessionUUID,
		})

		responses.SuccessMessageResponse(c, "Logged out successfully")
	}
}

// @Summary Check authentication status
// @Description Check if current session is valid
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 401 {object} responses.SendErrorResponse "Unauthorized"
// @Failure 500 {object} responses.SendErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		responses.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get session ID from context
	sessionID, _ := c.Get("session_id")
	sessionUUID, _ := sessionID.(uuid.UUID)

	// Get user details
	user, err := h.userService.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "User lookup failed", "Failed to retrieve user information.")
		return
	}

	// Get session details if available
	var sessionInfo *SessionInfo
	if sessionUUID != uuid.Nil {
		session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionUUID)
		if err == nil {
			deviceInfoStr := ""
			if session.DeviceInfo != nil {
				if deviceInfoBytes, err := json.Marshal(session.DeviceInfo); err == nil {
					deviceInfoStr = string(deviceInfoBytes)
				}
			}
			sessionInfo = &SessionInfo{
				ID:         session.ID,
				DeviceInfo: deviceInfoStr,
				IPAddress:  session.IPAddress,
				CreatedAt:  session.CreatedAt,
				ExpiresAt:  session.RefreshExpiresAt,
			}
		}
	}

	// Build response data
	data := gin.H{
		"user": UserInfo{
			ID:          user.ID,
			PhoneNumber: maskPhoneNumber(user.PhoneNumber),
			CountryCode: user.CountryCode,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			KYCStatus:   "", // KYC status not directly in User model
			KYCTier:     fmt.Sprintf("%d", int(user.KYCTier)),
			IsActive:    user.IsActive(),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		"authenticated": true,
	}

	if sessionInfo != nil {
		data["session"] = sessionInfo
	}

	responses.SendDataResponse(c, data)
}

// Helper functions

func getLanguageFromDeviceInfo(deviceInfo *DeviceInfo) string {
	if deviceInfo != nil && deviceInfo.Language != "" {
		return deviceInfo.Language
	}
	return "en" // Default language
}

func buildSessionDeviceInfo(deviceInfo *DeviceInfo) string {
	if deviceInfo == nil {
		return "Unknown Device"
	}

	var info []string
	if deviceInfo.Platform != "" {
		info = append(info, strings.Title(deviceInfo.Platform))
	}
	if deviceInfo.DeviceModel != "" {
		info = append(info, deviceInfo.DeviceModel)
	}
	if deviceInfo.OSVersion != "" {
		info = append(info, "OS "+deviceInfo.OSVersion)
	}

	if len(info) == 0 {
		return "Unknown Device"
	}

	return strings.Join(info, " - ")
}

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(router *gin.RouterGroup,
	handler *AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Handler is passed as parameter

	// Public auth routes
	public := router.Group("/auth")
	{
		public.POST("/otp/request", handler.RequestOTP)
		public.POST("/otp/verify", handler.VerifyOTP)
		public.POST("/refresh", handler.RefreshToken)
	}

	// Protected auth routes
	protected := router.Group("/auth")
	protected.Use(authMiddleware)
	{
		protected.POST("/logout", handler.Logout)
		protected.GET("/me", handler.GetCurrentUser)
	}
}

// Helper functions

// extractPhoneFromRequestID extracts phone number from request ID
// For now using a simple extraction, in real implementation would look up in database
func extractPhoneFromRequestID(requestID string) string {
	// This is a simplified implementation
	// In real scenario, requestID would be stored with associated phone number
	// For now, assume request ID format contains phone number
	return "+66812345678" // Placeholder - should look up in OTP storage
}

// extractCountryFromPhone extracts country code from phone number
func extractCountryFromPhone(phoneNumber string) string {
	if strings.HasPrefix(phoneNumber, "+66") {
		return "TH" // Thailand
	}
	if strings.HasPrefix(phoneNumber, "+65") {
		return "SG" // Singapore
	}
	if strings.HasPrefix(phoneNumber, "+62") {
		return "ID" // Indonesia
	}
	if strings.HasPrefix(phoneNumber, "+60") {
		return "MY" // Malaysia
	}
	if strings.HasPrefix(phoneNumber, "+63") {
		return "PH" // Philippines
	}
	if strings.HasPrefix(phoneNumber, "+84") {
		return "VN" // Vietnam
	}
	return "TH" // Default to Thailand
}

// maskPhoneNumber masks phone number for security
func maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 4 {
		return phoneNumber
	}
	visible := phoneNumber[:4]
	masked := strings.Repeat("*", len(phoneNumber)-4)
	return visible + masked
}

// buildSessionDeviceInfoMap builds device info map for session
func buildSessionDeviceInfoMap(deviceInfo *DeviceInfo) map[string]interface{} {
	if deviceInfo == nil {
		return map[string]interface{}{
			"platform": "unknown",
		}
	}
	return map[string]interface{}{
		"platform":      deviceInfo.Platform,
		"device_model":  deviceInfo.DeviceModel,
		"os_version":    deviceInfo.OSVersion,
		"app_version":   deviceInfo.AppVersion,
		"user_agent":    deviceInfo.UserAgent,
		"timezone":      deviceInfo.Timezone,
		"language":      deviceInfo.Language,
	}
}

// getKYCStatusString returns KYC status as string
func getKYCStatusString(user *sharedModels.User) string {
	if user.PhoneVerified || user.EmailVerified {
		return "verified"
	}
	return "pending"
}

// convertDeviceInfoToString converts device info map to display string
func convertDeviceInfoToString(deviceInfo map[string]interface{}) string {
	if deviceInfo == nil {
		return "Unknown"
	}
	if platform, ok := deviceInfo["platform"].(string); ok {
		if model, ok := deviceInfo["device_model"].(string); ok {
			return platform + " " + model
		}
		return platform
	}
	return "Unknown"
}

// RegisterAuthRoutesWithMiddleware registers auth routes with custom middleware
func RegisterAuthRoutesWithMiddleware(
	router *gin.RouterGroup,
	authService *services.AuthService,
	sessionService *services.SessionService,
	userService *services.UserService,
	publicMiddlewares []gin.HandlerFunc,
	protectedMiddlewares []gin.HandlerFunc,
) {
	handler := NewAuthHandler(authService, sessionService, userService)

	// Public auth routes with middleware
	public := router.Group("/auth")
	public.Use(publicMiddlewares...)
	{
		public.POST("/otp/verify", handler.VerifyOTP)
		public.POST("/refresh", handler.RefreshToken)
	}

	// Protected auth routes with middleware
	protected := router.Group("/auth")
	// Auth middleware should be included in protectedMiddlewares by caller
	allProtectedMiddlewares := protectedMiddlewares
	protected.Use(allProtectedMiddlewares...)
	{
		protected.POST("/logout", handler.Logout)
		protected.GET("/me", handler.GetCurrentUser)
	}
}

// SECURITY FIX: Enhanced phone number validation for Southeast Asian formats
func (h *AuthHandler) validatePhoneNumber(phoneNumber, countryCode string) error {
	// Sanitize and normalize phone number
	normalizedPhone := utils.NormalizePhone(phoneNumber)

	// Check if phone number is empty after normalization
	if normalizedPhone == "" {
		return fmt.Errorf("phone number is required")
	}

	// Validate country code is supported
	if !utils.IsValidSEACountryCode(countryCode) {
		return fmt.Errorf("unsupported country code: %s (supported: TH, SG, ID, MY, PH, VN)", countryCode)
	}

	// Validate phone number format using shared utilities
	if !utils.IsValidPhoneNumber(normalizedPhone, countryCode) {
		return fmt.Errorf("invalid phone number format for country %s", countryCode)
	}

	// Additional validation using auth models
	country := models.Country(countryCode)
	if !models.IsValidPhoneNumber(normalizedPhone, country) {
		return fmt.Errorf("phone number does not match required format for %s", countryCode)
	}

	// Check phone number length constraints
	if len(normalizedPhone) < 10 || len(normalizedPhone) > 15 {
		return fmt.Errorf("phone number length must be between 10-15 characters")
	}

	// Security check: Prevent commonly abused phone numbers
	if h.isBlacklistedPhoneNumber(normalizedPhone) {
		return fmt.Errorf("phone number is not allowed")
	}

	return nil
}

// isBlacklistedPhoneNumber checks for commonly abused phone numbers
func (h *AuthHandler) isBlacklistedPhoneNumber(phoneNumber string) bool {
	// Common test/fake numbers that should be blocked in production
	blacklistedPatterns := []string{
		"+0000000000",
		"+1111111111",
		"+1234567890",
		"+9999999999",
		// Add more patterns as needed for security
	}

	for _, pattern := range blacklistedPatterns {
		if phoneNumber == pattern {
			return true
		}
	}

	// Check for obviously fake patterns (all same digits)
	if len(phoneNumber) > 4 {
		firstChar := phoneNumber[1] // Skip + sign
		allSame := true
		for i := 1; i < len(phoneNumber); i++ {
			if phoneNumber[i] != firstChar {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	return false
}