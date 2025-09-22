package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/auth/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

// AuthHandlers provides comprehensive authentication HTTP handlers
type AuthHandlers struct {
	authService    *services.AuthService
	userService    *services.UserService
	sessionService *services.SessionService
	jwtService     *services.JWTService
	validator      *validator.Validate
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(
	authService *services.AuthService,
	userService *services.UserService,
	sessionService *services.SessionService,
	jwtService *services.JWTService,
) *AuthHandlers {
	return &AuthHandlers{
		authService:    authService,
		userService:    userService,
		sessionService: sessionService,
		jwtService:     jwtService,
		validator:      validator.New(),
	}
}

// Login handles user login via OTP verification
// @Summary User login with OTP
// @Description Verify OTP code and create authentication session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} responses.DataResponse{data=LoginResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Verify OTP
	verifyReq := &services.VerifyOTPRequest{
		PhoneNumber: req.PhoneNumber,
		Code:        req.Code,
		UserAgent:   c.GetHeader("User-Agent"),
		IPAddress:   c.ClientIP(),
		DeviceInfo:  req.DeviceInfo,
	}

	otpResult, err := h.authService.VerifyOTP(c.Request.Context(), verifyReq)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(
		c.Request.Context(),
		otpResult.User,
		otpResult.Session.ID,
		req.DeviceID,
	)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Token generation failed", err.Error())
		return
	}

	// Log successful login
	middleware.LogInfo(c, "User login successful", gin.H{
		"user_id":      otpResult.User.ID,
		"phone_number": middleware.MaskPhoneNumber(req.PhoneNumber),
		"device_id":    req.DeviceID,
		"ip_address":   c.ClientIP(),
	})

	// Build response
	response := LoginResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		TokenType:        tokenPair.TokenType,
		ExpiresIn:        tokenPair.ExpiresIn,
		AccessExpiresAt:  tokenPair.AccessExpiresAt,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt,
		User:             otpResult.User.ToResponse(),
		Session: SessionInfo{
			ID:        otpResult.Session.ID,
			DeviceID:  req.DeviceID,
			CreatedAt: otpResult.Session.CreatedAt,
			ExpiresAt: otpResult.Session.ExpiresAt,
		},
	}

	responses.DataResponse(c, response)
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} responses.DataResponse{data=RegisterResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Create user
	createUserReq := &services.CreateUserRequest{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Country:     req.Country,
		Language:    req.Language,
		TimeZone:    req.TimeZone,
		Metadata:    req.Metadata,
	}

	user, err := h.userService.CreateUser(c.Request.Context(), createUserReq)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			responses.ErrorResponse(c, http.StatusConflict, "User already exists", err.Error())
		} else {
			responses.ErrorResponse(c, http.StatusBadRequest, "Registration failed", err.Error())
		}
		return
	}

	// Log successful registration
	middleware.LogInfo(c, "User registration successful", gin.H{
		"user_id":      user.ID,
		"phone_number": middleware.MaskPhoneNumber(req.PhoneNumber),
		"country":      req.Country,
		"ip_address":   c.ClientIP(),
	})

	// Build response
	response := RegisterResponse{
		User:      user.ToResponse(),
		Message:   "Registration successful. Please verify your phone number to activate your account.",
		NextStep:  "phone_verification",
		CreatedAt: user.CreatedAt,
	}

	responses.DataResponseWithStatus(c, http.StatusCreated, response)
}

// RefreshToken handles JWT token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} responses.DataResponse{data=TokenResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", err.Error())
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, "User not found", err.Error())
		return
	}

	// Generate new token pair
	tokenPair, err := h.jwtService.RefreshAccessToken(c.Request.Context(), req.RefreshToken, user)
	if err != nil {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	// Log token refresh
	middleware.LogInfo(c, "Token refreshed successfully", gin.H{
		"user_id":    claims.UserID,
		"session_id": claims.SessionID,
		"ip_address": c.ClientIP(),
	})

	response := TokenResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		TokenType:        tokenPair.TokenType,
		ExpiresIn:        tokenPair.ExpiresIn,
		AccessExpiresAt:  tokenPair.AccessExpiresAt,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt,
	}

	responses.DataResponse(c, response)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user and terminate session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandlers) Logout(c *gin.Context) {
	// Get user claims from middleware
	claims, exists := middleware.GetUserClaims(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Terminate session
	err := h.sessionService.TerminateSession(c.Request.Context(), claims.SessionID, "user_logout")
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", err.Error())
		return
	}

	// Log logout
	middleware.LogInfo(c, "User logout successful", gin.H{
		"user_id":    claims.UserID,
		"session_id": claims.SessionID,
		"ip_address": c.ClientIP(),
	})

	responses.SuccessMessageResponse(c, "Logout successful")
}

// GetProfile gets current user profile
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=models.UserResponse}
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	// Get user ID from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Get user profile
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get profile", err.Error())
		return
	}

	responses.DataResponse(c, user.ToResponse())
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update current authenticated user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update request"
// @Success 200 {object} responses.DataResponse{data=models.UserResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/profile [put]
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	// Get user ID from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Update profile
	updateReq := &services.UpdateUserProfileRequest{
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Language:    req.Language,
		TimeZone:    req.TimeZone,
		Preferences: req.Preferences,
		Metadata:    req.Metadata,
	}

	user, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, updateReq)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Profile update failed", err.Error())
		return
	}

	// Log profile update
	middleware.LogInfo(c, "User profile updated", gin.H{
		"user_id":    userID,
		"ip_address": c.ClientIP(),
	})

	responses.DataResponse(c, user.ToResponse())
}

// GetSessions gets user sessions
// @Summary Get user sessions
// @Description Get all sessions for current authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=SessionListResponse}
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/sessions [get]
func (h *AuthHandlers) GetSessions(c *gin.Context) {
	// Get user ID and current session from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	claims, _ := middleware.GetUserClaims(c)
	currentSessionID := claims.SessionID

	// Get user sessions
	sessions, err := h.sessionService.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions", err.Error())
		return
	}

	// Convert to response format
	var sessionDetails []*SessionDetails
	for _, session := range sessions {
		isCurrent := session.ID == currentSessionID

		sessionDetails = append(sessionDetails, &SessionDetails{
			ID:           session.ID,
			Status:       session.Status,
			UserAgent:    session.UserAgent,
			IPAddress:    session.IPAddress,
			DeviceInfo:   session.DeviceInfo,
			CreatedAt:    session.CreatedAt,
			LastActiveAt: session.LastActiveAt,
			ExpiresAt:    session.ExpiresAt,
			IsCurrent:    isCurrent,
		})
	}

	response := SessionListResponse{
		Sessions: sessionDetails,
		Total:    len(sessionDetails),
	}

	responses.DataResponse(c, response)
}

// TerminateSession terminates a specific session
// @Summary Terminate session
// @Description Terminate a specific session by ID
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/sessions/{id} [delete]
func (h *AuthHandlers) TerminateSession(c *gin.Context) {
	// Get user ID from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Parse session ID
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID", err.Error())
		return
	}

	// Get session to verify ownership
	session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusNotFound, "Session not found", err.Error())
		return
	}

	// Verify session belongs to user
	if session.UserID != userID {
		responses.ErrorResponse(c, http.StatusForbidden, "Session does not belong to user", "")
		return
	}

	// Terminate session
	err = h.sessionService.TerminateSession(c.Request.Context(), sessionID, "user_request")
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to terminate session", err.Error())
		return
	}

	// Log session termination
	middleware.LogInfo(c, "Session terminated by user", gin.H{
		"user_id":           userID,
		"terminated_session": sessionID,
		"ip_address":        c.ClientIP(),
	})

	responses.SuccessMessageResponse(c, "Session terminated successfully")
}

// TerminateAllSessions terminates all user sessions except current
// @Summary Terminate all sessions
// @Description Terminate all sessions for current user except the current one
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/sessions/terminate-all [post]
func (h *AuthHandlers) TerminateAllSessions(c *gin.Context) {
	// Get user ID and current session from middleware
	userID, exists := middleware.GetUserID(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	claims, _ := middleware.GetUserClaims(c)
	currentSessionID := claims.SessionID

	// Get all user sessions
	sessions, err := h.sessionService.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions", err.Error())
		return
	}

	// Terminate all sessions except current
	terminatedCount := 0
	for _, session := range sessions {
		if session.ID != currentSessionID && session.Status == "active" {
			if err := h.sessionService.TerminateSession(c.Request.Context(), session.ID, "terminate_all_request"); err != nil {
				// Log error but continue with other sessions
				middleware.LogError(c, "Failed to terminate session", gin.H{
					"session_id": session.ID,
					"error":      err.Error(),
				})
			} else {
				terminatedCount++
			}
		}
	}

	// Log bulk termination
	middleware.LogInfo(c, "All other sessions terminated", gin.H{
		"user_id":           userID,
		"terminated_count":  terminatedCount,
		"current_session":   currentSessionID,
		"ip_address":        c.ClientIP(),
	})

	responses.SuccessMessageResponse(c, strconv.Itoa(terminatedCount)+" sessions terminated successfully")
}

// ValidateToken validates JWT token and returns user info
// @Summary Validate token
// @Description Validate JWT access token and return user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.DataResponse{data=TokenValidationResponse}
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandlers) ValidateToken(c *gin.Context) {
	// Get user claims from middleware (token already validated)
	claims, exists := middleware.GetUserClaims(c)
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Get user details
	user, err := h.userService.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err.Error())
		return
	}

	// Get session details
	session, err := h.sessionService.GetSessionByID(c.Request.Context(), claims.SessionID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get session", err.Error())
		return
	}

	response := TokenValidationResponse{
		Valid:       true,
		User:        user.ToResponse(),
		Claims:      *claims,
		Session:     session.ToDetailsResponse(true),
		ValidatedAt: time.Now(),
	}

	responses.DataResponse(c, response)
}

// ChangePassword handles password change (for future email-based auth)
// @Summary Change password
// @Description Change user password (requires current password)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandlers) ChangePassword(c *gin.Context) {
	// This would be implemented when email-based authentication is added
	responses.ErrorResponse(c, http.StatusNotImplemented, "Password authentication not yet implemented", "")
}

// RegisterRoutes registers all auth routes
func RegisterAuthRoutes(router *gin.RouterGroup, handlers *AuthHandlers, authMiddleware gin.HandlerFunc) {
	// Public routes (no authentication required)
	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
		auth.POST("/refresh", handlers.RefreshToken)
	}

	// Protected routes (authentication required)
	authProtected := router.Group("/auth")
	authProtected.Use(authMiddleware)
	{
		authProtected.POST("/logout", handlers.Logout)
		authProtected.GET("/profile", handlers.GetProfile)
		authProtected.PUT("/profile", handlers.UpdateProfile)
		authProtected.GET("/sessions", handlers.GetSessions)
		authProtected.DELETE("/sessions/:id", handlers.TerminateSession)
		authProtected.POST("/sessions/terminate-all", handlers.TerminateAllSessions)
		authProtected.GET("/validate", handlers.ValidateToken)
		authProtected.POST("/change-password", handlers.ChangePassword)
	}
}

// Request/Response structures

type LoginRequest struct {
	PhoneNumber string                 `json:"phone_number" validate:"required,e164" example:"+66812345678"`
	Code        string                 `json:"code" validate:"required,len=6,numeric" example:"123456"`
	DeviceID    string                 `json:"device_id" validate:"required" example:"device_123"`
	DeviceInfo  map[string]interface{} `json:"device_info,omitempty"`
}

type RegisterRequest struct {
	PhoneNumber string                 `json:"phone_number" validate:"required,e164" example:"+66812345678"`
	Email       string                 `json:"email,omitempty" validate:"omitempty,email" example:"user@example.com"`
	Username    string                 `json:"username,omitempty" validate:"omitempty,min=3,max=30" example:"johndoe"`
	FirstName   string                 `json:"first_name,omitempty" validate:"omitempty,min=1,max=50" example:"John"`
	LastName    string                 `json:"last_name,omitempty" validate:"omitempty,min=1,max=50" example:"Doe"`
	Country     string                 `json:"country" validate:"required,len=2" example:"TH"`
	Language    string                 `json:"language" validate:"required,len=2" example:"en"`
	TimeZone    string                 `json:"timezone" validate:"required" example:"Asia/Bangkok"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UpdateProfileRequest struct {
	Username    *string                    `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	FirstName   *string                    `json:"first_name,omitempty" validate:"omitempty,min=1,max=50"`
	LastName    *string                    `json:"last_name,omitempty" validate:"omitempty,min=1,max=50"`
	Email       *string                    `json:"email,omitempty" validate:"omitempty,email"`
	Language    *string                    `json:"language,omitempty" validate:"omitempty,len=2"`
	TimeZone    *string                    `json:"timezone,omitempty"`
	Preferences *services.UserPreferences  `json:"preferences,omitempty"`
	Metadata    map[string]interface{}     `json:"metadata,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type LoginResponse struct {
	AccessToken      string                 `json:"access_token"`
	RefreshToken     string                 `json:"refresh_token"`
	TokenType        string                 `json:"token_type"`
	ExpiresIn        int64                  `json:"expires_in"`
	AccessExpiresAt  time.Time              `json:"access_expires_at"`
	RefreshExpiresAt time.Time              `json:"refresh_expires_at"`
	User             *services.UserResponse `json:"user"`
	Session          SessionInfo            `json:"session"`
}

type RegisterResponse struct {
	User      *services.UserResponse `json:"user"`
	Message   string                 `json:"message"`
	NextStep  string                 `json:"next_step"`
	CreatedAt time.Time              `json:"created_at"`
}

type TokenResponse struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type SessionInfo struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  string    `json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionDetails struct {
	ID           uuid.UUID              `json:"id"`
	Status       string                 `json:"status"`
	UserAgent    string                 `json:"user_agent"`
	IPAddress    string                 `json:"ip_address"`
	DeviceInfo   map[string]interface{} `json:"device_info"`
	CreatedAt    time.Time              `json:"created_at"`
	LastActiveAt time.Time              `json:"last_active_at"`
	ExpiresAt    time.Time              `json:"expires_at"`
	IsCurrent    bool                   `json:"is_current"`
}

type SessionListResponse struct {
	Sessions []*SessionDetails `json:"sessions"`
	Total    int               `json:"total"`
}

type TokenValidationResponse struct {
	Valid       bool                           `json:"valid"`
	User        *services.UserResponse         `json:"user"`
	Claims      services.UserClaims            `json:"claims"`
	Session     *services.SessionDetailsResponse `json:"session"`
	ValidatedAt time.Time                      `json:"validated_at"`
}