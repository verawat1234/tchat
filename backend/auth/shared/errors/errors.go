package errors

import (
	"fmt"
)

// ApplicationError represents a custom application error
type ApplicationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error codes
const (
	InvalidPhoneNumber    = "invalid_phone_number"
	InvalidOTPCode        = "invalid_otp_code"
	OTPExpired           = "otp_expired"
	OTPNotFound          = "otp_not_found"
	UserNotFound         = "user_not_found"
	InvalidCredentials   = "invalid_credentials"
	Unauthorized         = "unauthorized"
	ValidationFailed     = "validation_failed"
	InternalError        = "internal_error"
)

// NewApplicationError creates a new application error
func NewApplicationError(code, message, errorType string) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}

// Common error constructors
func NewValidationError(message string) *ApplicationError {
	return NewApplicationError(ValidationFailed, message, "validation")
}

func NewUnauthorizedError(message string) *ApplicationError {
	return NewApplicationError(Unauthorized, message, "authorization")
}

func NewNotFoundError(message string) *ApplicationError {
	return NewApplicationError(UserNotFound, message, "not_found")
}

func NewInternalError(message string) *ApplicationError {
	return NewApplicationError(InternalError, message, "internal")
}