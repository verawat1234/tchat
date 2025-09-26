package errors_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ErrorUtils provides utilities for error response testing
type ErrorUtils struct{}

// NewErrorUtils creates a new error utils instance
func NewErrorUtils() *ErrorUtils {
	return &ErrorUtils{}
}

// StandardErrorCodes defines all standard error codes used across services
var StandardErrorCodes = map[string]HTTPErrorCode{
	"VALIDATION_ERROR": {
		HTTPStatus: 400,
		Code:       "VALIDATION_ERROR",
		Type:       "validation_error",
		Message:    "Input validation failed",
		Category:   "client_error",
	},
	"AUTHENTICATION_FAILED": {
		HTTPStatus: 401,
		Code:       "AUTHENTICATION_FAILED",
		Type:       "authentication_error",
		Message:    "Authentication credentials are invalid",
		Category:   "client_error",
	},
	"AUTHORIZATION_FAILED": {
		HTTPStatus: 403,
		Code:       "AUTHORIZATION_FAILED",
		Type:       "authorization_error",
		Message:    "Insufficient permissions for this resource",
		Category:   "client_error",
	},
	"RESOURCE_NOT_FOUND": {
		HTTPStatus: 404,
		Code:       "RESOURCE_NOT_FOUND",
		Type:       "not_found_error",
		Message:    "The requested resource was not found",
		Category:   "client_error",
	},
	"METHOD_NOT_ALLOWED": {
		HTTPStatus: 405,
		Code:       "METHOD_NOT_ALLOWED",
		Type:       "method_error",
		Message:    "HTTP method not allowed for this endpoint",
		Category:   "client_error",
	},
	"RESOURCE_CONFLICT": {
		HTTPStatus: 409,
		Code:       "RESOURCE_CONFLICT",
		Type:       "conflict_error",
		Message:    "Resource conflict detected",
		Category:   "client_error",
	},
	"UNPROCESSABLE_ENTITY": {
		HTTPStatus: 422,
		Code:       "UNPROCESSABLE_ENTITY",
		Type:       "processing_error",
		Message:    "Request cannot be processed",
		Category:   "client_error",
	},
	"RATE_LIMIT_EXCEEDED": {
		HTTPStatus: 429,
		Code:       "RATE_LIMIT_EXCEEDED",
		Type:       "rate_limit_error",
		Message:    "Rate limit exceeded",
		Category:   "client_error",
	},
	"INTERNAL_SERVER_ERROR": {
		HTTPStatus: 500,
		Code:       "INTERNAL_SERVER_ERROR",
		Type:       "internal_error",
		Message:    "An internal server error occurred",
		Category:   "server_error",
	},
	"SERVICE_UNAVAILABLE": {
		HTTPStatus: 503,
		Code:       "SERVICE_UNAVAILABLE",
		Type:       "availability_error",
		Message:    "Service temporarily unavailable",
		Category:   "server_error",
	},
	"GATEWAY_TIMEOUT": {
		HTTPStatus: 504,
		Code:       "GATEWAY_TIMEOUT",
		Type:       "timeout_error",
		Message:    "Gateway timeout occurred",
		Category:   "server_error",
	},
}

// ErrorSuggestions provides helpful suggestions for each error type
var ErrorSuggestions = map[string]string{
	"VALIDATION_ERROR":        "Please check the highlighted fields and correct any errors before resubmitting.",
	"AUTHENTICATION_FAILED":   "Please check your credentials and try again. If you forgot your password, use the password reset feature.",
	"AUTHORIZATION_FAILED":    "You don't have permission to access this resource. Please contact an administrator if you believe this is an error.",
	"RESOURCE_NOT_FOUND":      "The resource may have been deleted or moved. Please check the ID and try again.",
	"METHOD_NOT_ALLOWED":      "Please check the API documentation for the correct HTTP method for this endpoint.",
	"RESOURCE_CONFLICT":       "The resource you're trying to create already exists or conflicts with an existing resource.",
	"UNPROCESSABLE_ENTITY":    "Please review your request data and ensure all required fields are provided correctly.",
	"RATE_LIMIT_EXCEEDED":     "Please wait before making another request. Consider upgrading your plan for higher limits.",
	"INTERNAL_SERVER_ERROR":   "An unexpected error occurred. Please try again later or contact support if the problem persists.",
	"SERVICE_UNAVAILABLE":     "The service is temporarily down for maintenance. Please try again in a few minutes.",
	"GATEWAY_TIMEOUT":         "The request timed out. Please try again or contact support if the problem persists.",
}

// ErrorHelpURLs provides documentation URLs for each error type
var ErrorHelpURLs = map[string]string{
	"VALIDATION_ERROR":        "https://docs.tchat-backend/errors/validation",
	"AUTHENTICATION_FAILED":   "https://docs.tchat-backend/authentication",
	"AUTHORIZATION_FAILED":    "https://docs.tchat-backend/authorization",
	"RESOURCE_NOT_FOUND":      "https://docs.tchat-backend/resources",
	"METHOD_NOT_ALLOWED":      "https://docs.tchat-backend/api-methods",
	"RESOURCE_CONFLICT":       "https://docs.tchat-backend/errors/conflicts",
	"UNPROCESSABLE_ENTITY":    "https://docs.tchat-backend/validation",
	"RATE_LIMIT_EXCEEDED":     "https://docs.tchat-backend/rate-limits",
	"INTERNAL_SERVER_ERROR":   "https://docs.tchat-backend/troubleshooting",
	"SERVICE_UNAVAILABLE":     "https://status.tchat.dev",
	"GATEWAY_TIMEOUT":         "https://docs.tchat-backend/troubleshooting/timeouts",
}

// Southeast Asian Localized Error Messages
var SEALocalizedMessages = map[string]map[string]string{
	"VALIDATION_ERROR": {
		"TH": "ข้อมูลที่ป้อนไม่ถูกต้อง",
		"SG": "Input validation failed",
		"ID": "Validasi input gagal",
		"MY": "Pengesahan input gagal",
		"VN": "Xác thực đầu vào thất bại",
		"PH": "Input validation failed",
		"EN": "Input validation failed",
	},
	"AUTHENTICATION_FAILED": {
		"TH": "การยืนยันตัวตนล้มเหลว",
		"SG": "Authentication failed",
		"ID": "Autentikasi gagal",
		"MY": "Pengesahan gagal",
		"VN": "Xác thực thất bại",
		"PH": "Authentication failed",
		"EN": "Authentication failed",
	},
	"AUTHORIZATION_FAILED": {
		"TH": "ไม่มีสิทธิ์เข้าถึงข้อมูลนี้",
		"SG": "Access denied",
		"ID": "Akses ditolak",
		"MY": "Akses ditolak",
		"VN": "Truy cập bị từ chối",
		"PH": "Access denied",
		"EN": "Access denied",
	},
	"RESOURCE_NOT_FOUND": {
		"TH": "ไม่พบข้อมูลที่ร้องขอ",
		"SG": "Resource not found",
		"ID": "Sumber daya tidak ditemukan",
		"MY": "Sumber tidak dijumpai",
		"VN": "Không tìm thấy tài nguyên",
		"PH": "Resource not found",
		"EN": "Resource not found",
	},
	"RATE_LIMIT_EXCEEDED": {
		"TH": "เกินขีดจำกัดการเรียกใช้งาน",
		"SG": "Rate limit exceeded",
		"ID": "Batas tingkat terlampaui",
		"MY": "Had kadar melebihi",
		"VN": "Vượt quá giới hạn tốc độ",
		"PH": "Rate limit exceeded",
		"EN": "Rate limit exceeded",
	},
	"INTERNAL_SERVER_ERROR": {
		"TH": "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์",
		"SG": "Internal server error",
		"ID": "Kesalahan server internal",
		"MY": "Ralat pelayan dalaman",
		"VN": "Lỗi máy chủ nội bộ",
		"PH": "Internal server error",
		"EN": "Internal server error",
	},
}

// CreateStandardErrorResponse creates a standard error response
func (eu *ErrorUtils) CreateStandardErrorResponse(
	code string,
	message string,
	serviceName string,
	path string,
	method string,
) StandardErrorResponse {
	errorCode, exists := StandardErrorCodes[code]
	if !exists {
		errorCode = StandardErrorCodes["INTERNAL_SERVER_ERROR"]
	}

	if message == "" {
		message = errorCode.Message
	}

	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        code,
			Type:        errorCode.Type,
			Message:     message,
			ServiceName: serviceName,
			RequestID:   eu.GenerateRequestID(),
			Suggestion:  ErrorSuggestions[code],
			HelpURL:     ErrorHelpURLs[code],
		},
		TraceID: eu.GenerateTraceID(),
		Path:    path,
		Method:  method,
		Time:    time.Now(),
	}
}

// CreateValidationErrorResponse creates a validation error response with field details
func (eu *ErrorUtils) CreateValidationErrorResponse(
	validationErrors []ValidationError,
	serviceName string,
	path string,
	method string,
) StandardErrorResponse {
	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        "VALIDATION_ERROR",
			Type:        "validation_error",
			Message:     "Input validation failed",
			Validation:  validationErrors,
			ServiceName: serviceName,
			RequestID:   eu.GenerateRequestID(),
			Suggestion:  ErrorSuggestions["VALIDATION_ERROR"],
			HelpURL:     ErrorHelpURLs["VALIDATION_ERROR"],
			Details: map[string]interface{}{
				"failed_fields": eu.extractFailedFields(validationErrors),
				"total_errors":  len(validationErrors),
			},
		},
		TraceID: eu.GenerateTraceID(),
		Path:    path,
		Method:  method,
		Time:    time.Now(),
	}
}

// CreateLocalizedErrorResponse creates a localized error response
func (eu *ErrorUtils) CreateLocalizedErrorResponse(
	code string,
	country string,
	serviceName string,
	path string,
	method string,
) StandardErrorResponse {
	errorCode, exists := StandardErrorCodes[code]
	if !exists {
		errorCode = StandardErrorCodes["INTERNAL_SERVER_ERROR"]
	}

	// Get localized message
	message := errorCode.Message // Default to English
	if localizedMessages, exists := SEALocalizedMessages[code]; exists {
		if localizedMessage, exists := localizedMessages[country]; exists {
			message = localizedMessage
		}
	}

	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        code,
			Type:        errorCode.Type,
			Message:     message,
			ServiceName: serviceName,
			RequestID:   eu.GenerateRequestID(),
			Suggestion:  ErrorSuggestions[code],
			HelpURL:     ErrorHelpURLs[code],
			Details: map[string]interface{}{
				"locale": country,
			},
		},
		TraceID: eu.GenerateTraceID(),
		Path:    path,
		Method:  method,
		Time:    time.Now(),
	}
}

// CreateRateLimitErrorResponse creates a rate limit error response with details
func (eu *ErrorUtils) CreateRateLimitErrorResponse(
	limit int,
	window string,
	retryAfter int,
	currentUsage int,
	serviceName string,
	path string,
	method string,
) StandardErrorResponse {
	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        "RATE_LIMIT_EXCEEDED",
			Type:        "rate_limit_error",
			Message:     "Rate limit exceeded",
			ServiceName: serviceName,
			RequestID:   eu.GenerateRequestID(),
			Suggestion:  ErrorSuggestions["RATE_LIMIT_EXCEEDED"],
			HelpURL:     ErrorHelpURLs["RATE_LIMIT_EXCEEDED"],
			Details: map[string]interface{}{
				"limit":         limit,
				"window":        window,
				"retry_after":   retryAfter,
				"current_usage": currentUsage,
			},
		},
		TraceID: eu.GenerateTraceID(),
		Path:    path,
		Method:  method,
		Time:    time.Now(),
	}
}

// CreateConflictErrorResponse creates a resource conflict error response
func (eu *ErrorUtils) CreateConflictErrorResponse(
	conflictType string,
	conflictingField string,
	existingID string,
	serviceName string,
	path string,
	method string,
) StandardErrorResponse {
	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        "RESOURCE_CONFLICT",
			Type:        "conflict_error",
			Message:     "Resource conflict detected",
			ServiceName: serviceName,
			RequestID:   eu.GenerateRequestID(),
			Suggestion:  ErrorSuggestions["RESOURCE_CONFLICT"],
			HelpURL:     ErrorHelpURLs["RESOURCE_CONFLICT"],
			Details: map[string]interface{}{
				"conflict_type":     conflictType,
				"conflicting_field": conflictingField,
				"existing_id":       existingID,
				"resolution_hint":   fmt.Sprintf("use_different_%s", conflictingField),
			},
		},
		TraceID: eu.GenerateTraceID(),
		Path:    path,
		Method:  method,
		Time:    time.Now(),
	}
}

// ValidateErrorResponseStructure validates that an error response has all required fields
func (eu *ErrorUtils) ValidateErrorResponseStructure(errorResponse StandardErrorResponse) []string {
	var errors []string

	// Validate required top-level fields
	if errorResponse.Error.Code == "" {
		errors = append(errors, "error.code is required")
	}
	if errorResponse.Error.Type == "" {
		errors = append(errors, "error.type is required")
	}
	if errorResponse.Error.Message == "" {
		errors = append(errors, "error.message is required")
	}
	if errorResponse.TraceID == "" {
		errors = append(errors, "trace_id is required")
	}
	if errorResponse.Time.IsZero() {
		errors = append(errors, "timestamp is required")
	}

	// Validate error code format
	if !eu.isValidErrorCode(errorResponse.Error.Code) {
		errors = append(errors, "error.code must be uppercase with underscores")
	}

	// Validate error type format
	if !eu.isValidErrorType(errorResponse.Error.Type) {
		errors = append(errors, "error.type must be lowercase with underscores")
	}

	// Validate trace ID format
	if !eu.isValidTraceID(errorResponse.TraceID) {
		errors = append(errors, "trace_id must be valid UUID format")
	}

	return errors
}

// SerializeToJSON serializes error response to JSON
func (eu *ErrorUtils) SerializeToJSON(errorResponse StandardErrorResponse) ([]byte, error) {
	return json.MarshalIndent(errorResponse, "", "  ")
}

// DeserializeFromJSON deserializes error response from JSON
func (eu *ErrorUtils) DeserializeFromJSON(jsonData []byte) (StandardErrorResponse, error) {
	var errorResponse StandardErrorResponse
	err := json.Unmarshal(jsonData, &errorResponse)
	return errorResponse, err
}

// GetHTTPStatusFromErrorCode returns HTTP status code for error code
func (eu *ErrorUtils) GetHTTPStatusFromErrorCode(code string) int {
	if errorCode, exists := StandardErrorCodes[code]; exists {
		return errorCode.HTTPStatus
	}
	return 500 // Default to internal server error
}

// GetErrorCategory categorizes error based on HTTP status
func (eu *ErrorUtils) GetErrorCategory(httpStatus int) string {
	if httpStatus >= 400 && httpStatus < 500 {
		return "client_error"
	}
	if httpStatus >= 500 && httpStatus < 600 {
		return "server_error"
	}
	return "unknown"
}

// IsClientError checks if error is a client error (4xx)
func (eu *ErrorUtils) IsClientError(httpStatus int) bool {
	return httpStatus >= 400 && httpStatus < 500
}

// IsServerError checks if error is a server error (5xx)
func (eu *ErrorUtils) IsServerError(httpStatus int) bool {
	return httpStatus >= 500 && httpStatus < 600
}

// GenerateTraceID generates a unique trace ID
func (eu *ErrorUtils) GenerateTraceID() string {
	return uuid.New().String()
}

// GenerateRequestID generates a unique request ID
func (eu *ErrorUtils) GenerateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

// Helper methods

// extractFailedFields extracts field names from validation errors
func (eu *ErrorUtils) extractFailedFields(validationErrors []ValidationError) []string {
	fields := make([]string, len(validationErrors))
	for i, err := range validationErrors {
		fields[i] = err.Field
	}
	return fields
}

// isValidErrorCode validates error code format (UPPERCASE_WITH_UNDERSCORES)
func (eu *ErrorUtils) isValidErrorCode(code string) bool {
	if code == "" {
		return false
	}
	// Should be uppercase letters, numbers, and underscores only
	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}

// isValidErrorType validates error type format (lowercase_with_underscores)
func (eu *ErrorUtils) isValidErrorType(errorType string) bool {
	if errorType == "" {
		return false
	}
	// Should be lowercase letters and underscores only
	for _, char := range errorType {
		if !((char >= 'a' && char <= 'z') || char == '_') {
			return false
		}
	}
	return true
}

// isValidTraceID validates trace ID format (UUID)
func (eu *ErrorUtils) isValidTraceID(traceID string) bool {
	// Check if it's a valid UUID format
	_, err := uuid.Parse(traceID)
	if err != nil {
		// Also allow trace- prefix format for testing
		return strings.HasPrefix(traceID, "trace-") && len(traceID) > 6
	}
	return true
}

// ErrorResponseComparator provides methods to compare error responses
type ErrorResponseComparator struct{}

// NewErrorResponseComparator creates a new comparator
func NewErrorResponseComparator() *ErrorResponseComparator {
	return &ErrorResponseComparator{}
}

// CompareStructure compares the structure of two error responses
func (erc *ErrorResponseComparator) CompareStructure(response1, response2 StandardErrorResponse) bool {
	// Compare basic structure (ignoring values that should be unique like timestamps, IDs)
	return response1.Error.Code == response2.Error.Code &&
		response1.Error.Type == response2.Error.Type &&
		response1.Path == response2.Path &&
		response1.Method == response2.Method
}

// CompareSemantic compares the semantic content of two error responses
func (erc *ErrorResponseComparator) CompareSemantic(response1, response2 StandardErrorResponse) bool {
	// Compare all fields that should be semantically identical
	return response1.Error.Code == response2.Error.Code &&
		response1.Error.Type == response2.Error.Type &&
		response1.Error.Message == response2.Error.Message &&
		response1.Error.ServiceName == response2.Error.ServiceName &&
		response1.Path == response2.Path &&
		response1.Method == response2.Method
}

// IsCompliantWithStandard checks if error response complies with Tchat standards
func (erc *ErrorResponseComparator) IsCompliantWithStandard(response StandardErrorResponse) bool {
	// Check required fields
	if response.Error.Code == "" || response.Error.Type == "" || response.Error.Message == "" {
		return false
	}

	// Check if error code is in standard codes
	if _, exists := StandardErrorCodes[response.Error.Code]; !exists {
		return false
	}

	// Check trace ID format
	if response.TraceID == "" {
		return false
	}

	// Check timestamp
	if response.Time.IsZero() {
		return false
	}

	return true
}