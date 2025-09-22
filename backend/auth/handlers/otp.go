package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"tchat.dev/auth/services"
	"tchat.dev/shared/errors"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

type OTPHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewOTPHandler(authService *services.AuthService) *OTPHandler {
	return &OTPHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// SendOTPRequest represents the request to send OTP
type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164" example:"+66812345678"`
	CountryCode string `json:"country_code" validate:"required,len=2" example:"TH"`
	Language    string `json:"language,omitempty" validate:"omitempty,len=2" example:"en"`
	Purpose     string `json:"purpose,omitempty" validate:"omitempty,oneof=login registration password_reset verification" example:"login"`
}

// SendOTPResponse represents the response after sending OTP
type SendOTPResponse struct {
	Success     bool   `json:"success" example:"true"`
	Message     string `json:"message" example:"OTP sent successfully"`
	RequestID   string `json:"request_id" example:"req_123456789"`
	ExpiresIn   int    `json:"expires_in" example:"300"`
	NextAllowed int    `json:"next_allowed,omitempty" example:"60"`
}

// ResendOTPRequest represents the request to resend OTP
type ResendOTPRequest struct {
	RequestID string `json:"request_id" validate:"required" example:"req_123456789"`
}

// @Summary Send OTP
// @Description Send OTP code to phone number for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SendOTPRequest true "OTP request"
// @Success 200 {object} SendOTPResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 429 {object} responses.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/send [post]
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req SendOTPRequest

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

	// Normalize phone number and country code
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	req.CountryCode = strings.ToUpper(strings.TrimSpace(req.CountryCode))

	// Set default language if not provided
	if req.Language == "" {
		req.Language = "en"
	}

	// Set default purpose if not provided
	if req.Purpose == "" {
		req.Purpose = "login"
	}

	// Call auth service
	serviceReq := &services.SendOTPRequest{
		PhoneNumber: req.PhoneNumber,
		CountryCode: req.CountryCode,
		Language:    req.Language,
		Purpose:     req.Purpose,
	}

	result, err := h.authService.SendOTP(c.Request.Context(), serviceReq)
	if err != nil {
		// Handle specific error types
		switch {
		case strings.Contains(err.Error(), "rate limit"):
			responses.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many OTP requests. Please try again later.")
			return
		case strings.Contains(err.Error(), "invalid phone"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Invalid phone number", "Please provide a valid phone number for the specified country.")
			return
		case strings.Contains(err.Error(), "unsupported country"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Unsupported country", "The specified country is not supported.")
			return
		case strings.Contains(err.Error(), "SMS delivery failed"):
			responses.ErrorResponse(c, http.StatusServiceUnavailable, "SMS delivery failed", "Unable to send SMS. Please try again.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP", "Internal server error occurred.")
			return
		}
	}

	// Build response
	response := SendOTPResponse{
		Success:   true,
		Message:   "OTP sent successfully",
		RequestID: result.RequestID,
		ExpiresIn: result.ExpiresIn,
	}

	// Add next allowed time if rate limited
	if result.NextAllowedIn > 0 {
		response.NextAllowed = result.NextAllowedIn
	}

	// Log successful OTP request
	middleware.LogInfo(c, "OTP sent", gin.H{
		"phone_number": maskPhoneNumber(req.PhoneNumber),
		"country":      req.CountryCode,
		"purpose":      req.Purpose,
		"request_id":   result.RequestID,
	})

	responses.SuccessResponse(c, response)
}

// @Summary Resend OTP
// @Description Resend OTP code using existing request ID
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendOTPRequest true "Resend OTP request"
// @Success 200 {object} SendOTPResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 429 {object} responses.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *gin.Context) {
	var req ResendOTPRequest

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

	// Call auth service
	result, err := h.authService.ResendOTP(c.Request.Context(), req.RequestID)
	if err != nil {
		// Handle specific error types
		switch {
		case strings.Contains(err.Error(), "request not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Request not found", "Invalid or expired request ID.")
			return
		case strings.Contains(err.Error(), "rate limit"):
			responses.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many resend attempts. Please try again later.")
			return
		case strings.Contains(err.Error(), "max attempts"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Maximum attempts exceeded", "Maximum resend attempts reached for this request.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to resend OTP", "Internal server error occurred.")
			return
		}
	}

	// Build response
	response := SendOTPResponse{
		Success:     true,
		Message:     "OTP resent successfully",
		RequestID:   result.RequestID,
		ExpiresIn:   result.ExpiresIn,
		NextAllowed: result.NextAllowedIn,
	}

	// Log successful OTP resend
	middleware.LogInfo(c, "OTP resent", gin.H{
		"request_id": req.RequestID,
	})

	responses.SuccessResponse(c, response)
}

// @Summary Get OTP status
// @Description Get status and remaining attempts for OTP request
// @Tags auth
// @Produce json
// @Param request_id path string true "OTP Request ID"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/{request_id}/status [get]
func (h *OTPHandler) GetOTPStatus(c *gin.Context) {
	requestID := c.Param("request_id")
	if requestID == "" {
		responses.ErrorResponse(c, http.StatusBadRequest, "Missing request ID", "Request ID is required.")
		return
	}

	// Call auth service
	status, err := h.authService.GetOTPStatus(c.Request.Context(), requestID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.ErrorResponse(c, http.StatusNotFound, "Request not found", "Invalid or expired request ID.")
			return
		}
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get OTP status", "Internal server error occurred.")
		return
	}

	// Build response data
	data := gin.H{
		"request_id":       status.RequestID,
		"status":           status.Status,
		"attempts_left":    status.AttemptsLeft,
		"expires_at":       status.ExpiresAt,
		"can_resend":       status.CanResend,
		"next_resend_at":   status.NextResendAt,
		"phone_number":     maskPhoneNumber(status.PhoneNumber),
		"verification_count": status.VerificationCount,
	}

	responses.DataResponse(c, data)
}

// @Summary Cancel OTP request
// @Description Cancel an active OTP request
// @Tags auth
// @Produce json
// @Param request_id path string true "OTP Request ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/{request_id}/cancel [post]
func (h *OTPHandler) CancelOTP(c *gin.Context) {
	requestID := c.Param("request_id")
	if requestID == "" {
		responses.ErrorResponse(c, http.StatusBadRequest, "Missing request ID", "Request ID is required.")
		return
	}

	// Call auth service
	err := h.authService.CancelOTP(c.Request.Context(), requestID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.ErrorResponse(c, http.StatusNotFound, "Request not found", "Invalid or expired request ID.")
			return
		}
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel OTP", "Internal server error occurred.")
		return
	}

	// Log OTP cancellation
	middleware.LogInfo(c, "OTP cancelled", gin.H{
		"request_id": requestID,
	})

	responses.SuccessMessageResponse(c, "OTP request cancelled successfully")
}

// @Summary Validate phone number
// @Description Validate if phone number is supported for OTP
// @Tags auth
// @Produce json
// @Param phone_number query string true "Phone number in E.164 format"
// @Param country_code query string true "ISO 3166-1 alpha-2 country code"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/otp/validate-phone [get]
func (h *OTPHandler) ValidatePhoneNumber(c *gin.Context) {
	phoneNumber := c.Query("phone_number")
	countryCode := c.Query("country_code")

	if phoneNumber == "" || countryCode == "" {
		responses.ErrorResponse(c, http.StatusBadRequest, "Missing parameters", "Both phone_number and country_code are required.")
		return
	}

	// Normalize inputs
	phoneNumber = strings.TrimSpace(phoneNumber)
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))

	// Call auth service
	validation, err := h.authService.ValidatePhoneNumber(c.Request.Context(), phoneNumber, countryCode)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Validation failed", "Failed to validate phone number.")
		return
	}

	// Build response data
	data := gin.H{
		"valid":            validation.Valid,
		"normalized":       validation.NormalizedNumber,
		"country_code":     validation.CountryCode,
		"country_name":     validation.CountryName,
		"carrier":          validation.Carrier,
		"line_type":        validation.LineType,
		"supports_sms":     validation.SupportsSMS,
		"risk_score":       validation.RiskScore,
		"formatted_local":  validation.FormattedLocal,
		"formatted_intl":   validation.FormattedInternational,
	}

	// Add validation messages if any
	if len(validation.ValidationMessages) > 0 {
		data["messages"] = validation.ValidationMessages
	}

	responses.DataResponse(c, data)
}

// Helper function to mask phone number for logging
func maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 4 {
		return "****"
	}

	// Keep country code and last 2 digits visible
	if strings.HasPrefix(phoneNumber, "+") && len(phoneNumber) > 6 {
		return phoneNumber[:3] + "****" + phoneNumber[len(phoneNumber)-2:]
	}

	// Fallback masking
	return phoneNumber[:2] + "****" + phoneNumber[len(phoneNumber)-2:]
}

// RegisterOTPRoutes registers all OTP-related routes
func RegisterOTPRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	handler := NewOTPHandler(authService)

	// OTP routes
	otp := router.Group("/otp")
	{
		otp.POST("/send", handler.SendOTP)
		otp.POST("/resend", handler.ResendOTP)
		otp.GET("/:request_id/status", handler.GetOTPStatus)
		otp.POST("/:request_id/cancel", handler.CancelOTP)
		otp.GET("/validate-phone", handler.ValidatePhoneNumber)
	}
}

// RegisterOTPRoutesWithMiddleware registers OTP routes with custom middleware
func RegisterOTPRoutesWithMiddleware(router *gin.RouterGroup, authService *services.AuthService, middlewares ...gin.HandlerFunc) {
	handler := NewOTPHandler(authService)

	// Apply middleware to OTP routes
	otp := router.Group("/otp")
	otp.Use(middlewares...)
	{
		otp.POST("/send", handler.SendOTP)
		otp.POST("/resend", handler.ResendOTP)
		otp.GET("/:request_id/status", handler.GetOTPStatus)
		otp.POST("/:request_id/cancel", handler.CancelOTP)
		otp.GET("/validate-phone", handler.ValidatePhoneNumber)
	}
}

// Middleware for OTP request validation
func OTPRequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add request ID to context for tracking
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = middleware.GenerateRequestID()
		}
		c.Set("request_id", requestID)

		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	}
}

// RateLimitMiddleware for OTP endpoints
func OTPRateLimitMiddleware() gin.HandlerFunc {
	return middleware.RateLimit(middleware.RateLimitConfig{
		WindowSize:   60,     // 1 minute window
		MaxRequests:  5,      // 5 requests per minute per IP
		KeyFunc:      middleware.IPKeyFunc,
		SkipFunc:     nil,
		ErrorHandler: nil,
	})
}