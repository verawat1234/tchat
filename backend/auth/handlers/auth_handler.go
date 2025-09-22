package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tchat.dev/auth/models"
	"tchat.dev/auth/services"
	"tchat.dev/shared/utils"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService services.AuthService
	validator   *utils.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   utils.NewValidator(),
	}
}

// RegisterRoutes registers authentication routes
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	// Authentication routes
	auth := router.PathPrefix("/auth").Subrouter()

	// OTP endpoints
	auth.HandleFunc("/otp/send", h.SendOTP).Methods("POST")
	auth.HandleFunc("/otp/verify", h.VerifyOTP).Methods("POST")

	// Token endpoints
	auth.HandleFunc("/token/refresh", h.RefreshToken).Methods("POST")
	auth.HandleFunc("/token/validate", h.ValidateToken).Methods("POST")

	// Session endpoints
	auth.HandleFunc("/sessions", h.GetSessions).Methods("GET")
	auth.HandleFunc("/sessions/{id}", h.RevokeSession).Methods("DELETE")
	auth.HandleFunc("/sessions/revoke-all", h.RevokeAllSessions).Methods("DELETE")

	// User profile endpoints
	auth.HandleFunc("/profile", h.GetProfile).Methods("GET")
	auth.HandleFunc("/profile", h.UpdateProfile).Methods("PUT")
}

// SendOTP handles OTP sending requests
func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateSendOTPRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get client IP and User-Agent
	ipAddress := h.getClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	// Create service request
	authReq := &services.AuthRequest{
		Phone:     req.Phone,
		Email:     req.Email,
		Country:   req.Country,
		DeviceID:  req.DeviceID,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	}

	// Send OTP
	if err := h.authService.SendOTP(ctx, authReq); err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to send OTP", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "OTP sent successfully", map[string]interface{}{
		"message": "OTP has been sent to your phone/email",
		"expires_in": 300, // 5 minutes
	})
}

// VerifyOTP handles OTP verification requests
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateVerifyOTPRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get client IP and User-Agent
	ipAddress := h.getClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	// Create service request
	verifyReq := &services.OTPVerificationRequest{
		Phone:     req.Phone,
		Email:     req.Email,
		Country:   req.Country,
		Code:      req.Code,
		DeviceID:  req.DeviceID,
		Name:      req.Name,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	}

	// Verify OTP
	authResp, err := h.authService.VerifyOTP(ctx, verifyReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "OTP verification failed", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Authentication successful", map[string]interface{}{
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"token_type":    authResp.TokenType,
		"expires_in":    authResp.ExpiresIn,
		"expires_at":    authResp.ExpiresAt,
		"user":          h.sanitizeUser(authResp.User),
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateRefreshTokenRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Refresh token
	authResp, err := h.authService.RefreshToken(ctx, req.RefreshToken, req.DeviceID)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Token refreshed successfully", map[string]interface{}{
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"token_type":    authResp.TokenType,
		"expires_in":    authResp.ExpiresIn,
		"expires_at":    authResp.ExpiresAt,
		"user":          h.sanitizeUser(authResp.User),
	})
}

// ValidateToken handles token validation requests
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.respondError(w, http.StatusUnauthorized, "Authorization header required", nil)
		return
	}

	// Parse Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		h.respondError(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
		return
	}

	accessToken := parts[1]

	// Validate token
	user, err := h.authService.ValidateToken(ctx, accessToken)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Token validation failed", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Token is valid", map[string]interface{}{
		"user":  h.sanitizeUser(user),
		"valid": true,
	})
}

// GetSessions handles getting user sessions
func (h *AuthHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (set by auth middleware)
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get user sessions
	sessions, err := h.authService.GetUserSessions(ctx, user.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get sessions", err)
		return
	}

	// Sanitize sessions
	var sanitizedSessions []map[string]interface{}
	for _, session := range sessions {
		sanitizedSessions = append(sanitizedSessions, session.ToSessionInfo())
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Sessions retrieved successfully", map[string]interface{}{
		"sessions": sanitizedSessions,
		"count":    len(sanitizedSessions),
	})
}

// RevokeSession handles session revocation
func (h *AuthHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get session ID from URL
	vars := mux.Vars(r)
	sessionIDStr := vars["id"]

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}

	// Revoke session
	if err := h.authService.RevokeSession(ctx, sessionID, user.ID); err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to revoke session", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Session revoked successfully", map[string]interface{}{
		"session_id": sessionID,
		"revoked":    true,
	})
}

// RevokeAllSessions handles revoking all user sessions
func (h *AuthHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Revoke all sessions
	if err := h.authService.RevokeAllSessions(ctx, user.ID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to revoke sessions", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "All sessions revoked successfully", map[string]interface{}{
		"user_id": user.ID,
		"revoked": true,
	})
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"user": h.sanitizeUser(user),
	})
}

// UpdateProfile handles updating user profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateUpdateProfileRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Update user fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// Update user
	if err := user.BeforeUpdate(); err != nil {
		h.respondError(w, http.StatusBadRequest, "Profile validation failed", err)
		return
	}

	// Success response (in real implementation, you'd save to database)
	h.respondSuccess(w, http.StatusOK, "Profile updated successfully", map[string]interface{}{
		"user": h.sanitizeUser(user),
	})
}

// Request/Response types
type SendOTPRequest struct {
	Phone    *string         `json:"phone,omitempty"`
	Email    *string         `json:"email,omitempty"`
	Country  models.Country  `json:"country"`
	DeviceID string          `json:"device_id"`
}

type VerifyOTPRequest struct {
	Phone    *string         `json:"phone,omitempty"`
	Email    *string         `json:"email,omitempty"`
	Country  models.Country  `json:"country"`
	Code     string          `json:"code"`
	DeviceID string          `json:"device_id"`
	Name     *string         `json:"name,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

type UpdateProfileRequest struct {
	Name   *string            `json:"name,omitempty"`
	Avatar *string            `json:"avatar,omitempty"`
	Status *models.UserStatus `json:"status,omitempty"`
}

// Standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Validation methods
func (h *AuthHandler) validateSendOTPRequest(req *SendOTPRequest) error {
	h.validator.Reset()

	// Either phone or email required
	if req.Phone == nil && req.Email == nil {
		h.validator.AddError("contact", "either phone or email is required")
	}

	if req.Phone != nil && req.Email != nil {
		h.validator.AddError("contact", "provide either phone or email, not both")
	}

	// Validate phone if provided
	if req.Phone != nil {
		h.validator.Required("phone", *req.Phone).Phone("phone", *req.Phone)
	}

	// Validate email if provided
	if req.Email != nil {
		h.validator.Required("email", *req.Email).Email("email", *req.Email)
	}

	// Validate country
	if !req.Country.IsValid() {
		h.validator.AddError("country", "invalid country code")
	}

	// Validate device ID
	h.validator.Required("device_id", req.DeviceID).MinLength("device_id", req.DeviceID, 1)

	return h.validator.GetError()
}

func (h *AuthHandler) validateVerifyOTPRequest(req *VerifyOTPRequest) error {
	h.validator.Reset()

	// Either phone or email required
	if req.Phone == nil && req.Email == nil {
		h.validator.AddError("contact", "either phone or email is required")
	}

	if req.Phone != nil && req.Email != nil {
		h.validator.AddError("contact", "provide either phone or email, not both")
	}

	// Validate phone if provided
	if req.Phone != nil {
		h.validator.Required("phone", *req.Phone).Phone("phone", *req.Phone)
	}

	// Validate email if provided
	if req.Email != nil {
		h.validator.Required("email", *req.Email).Email("email", *req.Email)
	}

	// Validate country
	if !req.Country.IsValid() {
		h.validator.AddError("country", "invalid country code")
	}

	// Validate OTP code
	h.validator.Required("code", req.Code).Length("code", req.Code, 6).Numeric("code", req.Code)

	// Validate device ID
	h.validator.Required("device_id", req.DeviceID).MinLength("device_id", req.DeviceID, 1)

	// Validate name if provided
	if req.Name != nil {
		h.validator.MinLength("name", *req.Name, 2).MaxLength("name", *req.Name, 100)
	}

	return h.validator.GetError()
}

func (h *AuthHandler) validateRefreshTokenRequest(req *RefreshTokenRequest) error {
	h.validator.Reset()

	h.validator.Required("refresh_token", req.RefreshToken).MinLength("refresh_token", req.RefreshToken, 10)
	h.validator.Required("device_id", req.DeviceID).MinLength("device_id", req.DeviceID, 1)

	return h.validator.GetError()
}

func (h *AuthHandler) validateUpdateProfileRequest(req *UpdateProfileRequest) error {
	h.validator.Reset()

	// Validate name if provided
	if req.Name != nil {
		h.validator.MinLength("name", *req.Name, 2).MaxLength("name", *req.Name, 100)
	}

	// Validate avatar URL if provided
	if req.Avatar != nil {
		h.validator.URL("avatar", *req.Avatar)
	}

	// Validate status if provided
	if req.Status != nil && !req.Status.IsValid() {
		h.validator.AddError("status", "invalid user status")
	}

	return h.validator.GetError()
}

// Utility methods
func (h *AuthHandler) respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) respondError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errorData interface{}
	if err != nil {
		errorData = err.Error()
	}

	response := APIResponse{
		Success: false,
		Message: message,
		Error:   errorData,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

func (h *AuthHandler) getUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value("user").(*models.User); ok {
		return user
	}
	return nil
}

func (h *AuthHandler) sanitizeUser(user *models.User) map[string]interface{} {
	if user == nil {
		return nil
	}

	return map[string]interface{}{
		"id":          user.ID,
		"name":        user.Name,
		"avatar":      user.Avatar,
		"country":     user.Country,
		"locale":      user.Locale,
		"kyc_tier":    user.KYCTier,
		"status":      user.Status,
		"is_verified": user.IsVerified,
		"created_at":  user.CreatedAt,
		"last_seen":   user.LastSeen,
	}
}

// Middleware for authentication
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.respondError(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		// Parse Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.respondError(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
			return
		}

		accessToken := parts[1]

		// Validate token
		user, err := h.authService.ValidateToken(r.Context(), accessToken)
		if err != nil {
			h.respondError(w, http.StatusUnauthorized, "Token validation failed", err)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Optional authentication middleware (doesn't fail if no token)
func (h *AuthHandler) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Parse Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				accessToken := parts[1]

				// Validate token
				if user, err := h.authService.ValidateToken(r.Context(), accessToken); err == nil {
					// Add user to context
					ctx := context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Rate limiting middleware
func (h *AuthHandler) RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter
	clients := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := h.getClientIP(r)
			now := time.Now()

			// Clean old entries
			if timestamps, exists := clients[clientIP]; exists {
				var validTimestamps []time.Time
				for _, timestamp := range timestamps {
					if now.Sub(timestamp) < time.Minute {
						validTimestamps = append(validTimestamps, timestamp)
					}
				}
				clients[clientIP] = validTimestamps
			}

			// Check rate limit
			if len(clients[clientIP]) >= requestsPerMinute {
				h.respondError(w, http.StatusTooManyRequests, "Rate limit exceeded", nil)
				return
			}

			// Add current request
			clients[clientIP] = append(clients[clientIP], now)

			next.ServeHTTP(w, r)
		})
	}
}

// CORS middleware
func (h *AuthHandler) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func (h *AuthHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %d %v\n",
			start.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
		)
	})
}

// Response writer wrapper to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Pagination helper
func (h *AuthHandler) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset = 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}