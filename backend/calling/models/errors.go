package models

import "errors"

// Common model errors
var (
	ErrInvalidCallType          = errors.New("invalid call type")
	ErrInvalidCallStatus        = errors.New("invalid call status")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
	ErrInvalidParticipantRole   = errors.New("invalid participant role")
	ErrInvalidConnectionQuality = errors.New("invalid connection quality")
	ErrInvalidPresenceStatus    = errors.New("invalid presence status")
	ErrCallNotFound             = errors.New("call session not found")
	ErrParticipantNotFound      = errors.New("participant not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyInCall        = errors.New("user already in a call")
	ErrUserNotAvailable         = errors.New("user not available for calls")
	ErrCallAlreadyAnswered      = errors.New("call already answered")
	ErrCallAlreadyEnded         = errors.New("call already ended")
	ErrInvalidCallState         = errors.New("invalid call state for operation")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrPermissionDenied         = errors.New("permission denied")
)
