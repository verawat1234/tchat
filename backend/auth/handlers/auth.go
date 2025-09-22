package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/auth/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
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
	RequestID string `json:"request_id" validate:"required" example:"req_123456789"`
	Code      string `json:"code" validate:"required,len=6,numeric" example:"123456"`
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

// LogoutRequest represents the request to logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
	LogoutAll    bool   `json:"logout_all,omitempty" example:"false"`
}

// @Summary Verify OTP
// @Description Verify OTP code and complete authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP verification request"
// @Success 200 {object} VerifyOTPResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Invalid OTP"
// @Failure 429 {object} responses.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

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

	// Extract IP address if not provided in device info
	if req.DeviceInfo != nil && req.DeviceInfo.IPAddress == "" {
		req.DeviceInfo.IPAddress = c.ClientIP()
	}

	// Extract User-Agent if not provided
	if req.DeviceInfo != nil && req.DeviceInfo.UserAgent == "" {
		req.DeviceInfo.UserAgent = c.GetHeader("User-Agent")
	}

	// Call auth service to verify OTP
	serviceReq := &services.VerifyOTPRequest{
		RequestID: req.RequestID,
		Code:      req.Code,
	}

	verifyResult, err := h.authService.VerifyOTP(c.Request.Context(), serviceReq)
	if err != nil {
		// Handle specific error types
		switch {
		case strings.Contains(err.Error(), "invalid code"):
			responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid OTP", "The OTP code is incorrect or has expired.")
			return
		case strings.Contains(err.Error(), "expired"):
			responses.ErrorResponse(c, http.StatusUnauthorized, "OTP expired", "The OTP code has expired. Please request a new one.")
			return
		case strings.Contains(err.Error(), "max attempts"):
			responses.ErrorResponse(c, http.StatusTooManyRequests, "Maximum attempts exceeded", "Too many failed attempts. Please request a new OTP.")
			return
		case strings.Contains(err.Error(), "request not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Request not found", "Invalid or expired OTP request.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Verification failed", "Internal server error occurred.")
			return
		}
	}

	// Get or create user based on phone number
	user, err := h.userService.GetUserByPhoneNumber(c.Request.Context(), verifyResult.PhoneNumber)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		responses.ErrorResponse(c, http.StatusInternalServerError, "User lookup failed", "Failed to retrieve user information.")
		return
	}

	// Create new user if not exists
	if user == nil {
		createUserReq := &services.CreateUserRequest{
			PhoneNumber: verifyResult.PhoneNumber,
			CountryCode: verifyResult.CountryCode,
			Language:    getLanguageFromDeviceInfo(req.DeviceInfo),
		}

		user, err = h.userService.CreateUser(c.Request.Context(), createUserReq)
		if err != nil {
			responses.ErrorResponse(c, http.StatusInternalServerError, "User creation failed", "Failed to create user account.")
			return
		}
	}

	// Create session
	sessionReq := &services.CreateSessionRequest{
		UserID:     user.ID,
		DeviceInfo: buildSessionDeviceInfo(req.DeviceInfo),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), sessionReq)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Session creation failed", "Failed to create user session.")
		return
	}

	// Build response
	response := VerifyOTPResponse{
		Success:      true,
		Message:      "Authentication successful",
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(session.AccessTokenExpiresAt.Sub(time.Now()).Seconds()),
		User: UserInfo{
			ID:          user.ID,
			PhoneNumber: maskPhoneNumber(user.PhoneNumber),
			CountryCode: user.CountryCode,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			KYCStatus:   string(user.KYCStatus),
			KYCTier:     string(user.KYCTier),
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Session: SessionInfo{
			ID:         session.ID,
			DeviceInfo: session.DeviceInfo,
			IPAddress:  session.IPAddress,
			CreatedAt:  session.CreatedAt,
			ExpiresAt:  session.RefreshTokenExpiresAt,
		},
	}

	// Log successful authentication
	middleware.LogInfo(c, "Authentication successful", gin.H{
		"user_id":      user.ID,
		"phone_number": maskPhoneNumber(user.PhoneNumber),
		"session_id":   session.ID,
		"device_info":  session.DeviceInfo,
	})

	responses.SuccessResponse(c, response)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Token refresh request"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Invalid refresh token"
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

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

	// Get session by refresh token
	session, err := h.sessionService.GetSessionByRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", "The refresh token is invalid or has expired.")
			return
		}
		responses.ErrorResponse(c, http.StatusInternalServerError, "Token refresh failed", "Internal server error occurred.")
		return
	}

	// Refresh tokens
	newSession, err := h.sessionService.RefreshTokens(c.Request.Context(), session.ID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Token refresh failed", "Failed to refresh tokens.")
		return
	}

	// Update device info if provided
	if req.DeviceInfo != nil {
		updateReq := &services.UpdateSessionRequest{
			SessionID:  newSession.ID,
			DeviceInfo: buildSessionDeviceInfo(req.DeviceInfo),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		newSession, err = h.sessionService.UpdateSession(c.Request.Context(), updateReq)
		if err != nil {
			// Log warning but continue with response
			middleware.LogWarning(c, "Failed to update session device info", gin.H{
				"session_id": newSession.ID,
				"error":      err.Error(),
			})
		}
	}

	// Build response data
	data := gin.H{
		"access_token":  newSession.AccessToken,
		"refresh_token": newSession.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(newSession.AccessTokenExpiresAt.Sub(time.Now()).Seconds()),
		"session": gin.H{
			"id":         newSession.ID,
			"expires_at": newSession.RefreshTokenExpiresAt,
		},
	}

	// Log successful token refresh
	middleware.LogInfo(c, "Token refreshed", gin.H{
		"session_id": newSession.ID,
		"user_id":    newSession.UserID,
	})

	responses.DataResponse(c, data)
}

// @Summary Logout
// @Description Logout and invalidate session(s)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Logout request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
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
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Get session ID from context
	sessionID, exists := c.Get("session_id")
	var sessionUUID uuid.UUID
	if exists {
		sessionUUID, ok = sessionID.(uuid.UUID)
		if !ok {
			responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid session context", "Invalid session context.")
			return
		}
	}

	// Logout all sessions or just current session
	if req.LogoutAll {
		err := h.sessionService.InvalidateUserSessions(c.Request.Context(), userUUID)
		if err != nil {
			responses.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", "Failed to logout from all sessions.")
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
				responses.ErrorResponse(c, http.StatusBadRequest, "Invalid refresh token", "The provided refresh token is invalid.")
				return
			}
			sessionUUID = session.ID
		}

		if sessionUUID != uuid.Nil {
			err := h.sessionService.InvalidateSession(c.Request.Context(), sessionUUID)
			if err != nil {
				responses.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", "Failed to logout session.")
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
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
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

	// Get session ID from context
	sessionID, _ := c.Get("session_id")
	sessionUUID, _ := sessionID.(uuid.UUID)

	// Get user details
	user, err := h.userService.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "User lookup failed", "Failed to retrieve user information.")
		return
	}

	// Get session details if available
	var sessionInfo *SessionInfo
	if sessionUUID != uuid.Nil {
		session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionUUID)
		if err == nil {
			sessionInfo = &SessionInfo{
				ID:         session.ID,
				DeviceInfo: session.DeviceInfo,
				IPAddress:  session.IPAddress,
				CreatedAt:  session.CreatedAt,
				ExpiresAt:  session.RefreshTokenExpiresAt,
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
			KYCStatus:   string(user.KYCStatus),
			KYCTier:     string(user.KYCTier),
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		"authenticated": true,
	}

	if sessionInfo != nil {
		data["session"] = sessionInfo
	}

	responses.DataResponse(c, data)
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
	authService *services.AuthService,
	sessionService *services.SessionService,
	userService *services.UserService,
) {
	handler := NewAuthHandler(authService, sessionService, userService)

	// Public auth routes
	public := router.Group("/auth")
	{
		public.POST("/otp/verify", handler.VerifyOTP)
		public.POST("/refresh", handler.RefreshToken)
	}

	// Protected auth routes
	protected := router.Group("/auth")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/logout", handler.Logout)
		protected.GET("/me", handler.GetCurrentUser)
	}
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
	allProtectedMiddlewares := append(protectedMiddlewares, middleware.AuthRequired())
	protected.Use(allProtectedMiddlewares...)
	{
		protected.POST("/logout", handler.Logout)
		protected.GET("/me", handler.GetCurrentUser)
	}
}