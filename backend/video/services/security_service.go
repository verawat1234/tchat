// backend/video/services/security_service.go
// Security service for video streaming with signed URLs and token-based access control
// Prevents direct video downloads by implementing blob URL pattern

package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SecurityService handles video security and access control
type SecurityService struct {
	secretKey []byte
	tokenTTL  time.Duration
}

// NewSecurityService creates a new security service
func NewSecurityService(secretKey string) *SecurityService {
	return &SecurityService{
		secretKey: []byte(secretKey),
		tokenTTL:  2 * time.Hour, // Default 2 hour expiration
	}
}

// StreamToken represents a signed streaming token
type StreamToken struct {
	VideoID   uuid.UUID `json:"video_id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Quality   string    `json:"quality"`
	Signature string    `json:"signature"`
}

// GenerateStreamToken creates a signed token for streaming access
func (s *SecurityService) GenerateStreamToken(videoID, userID uuid.UUID, quality string) (*StreamToken, error) {
	expiresAt := time.Now().Add(s.tokenTTL)

	token := &StreamToken{
		VideoID:   videoID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Quality:   quality,
	}

	// Generate signature
	signature, err := s.generateSignature(token)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}

	token.Signature = signature
	return token, nil
}

// ValidateStreamToken verifies the token signature and expiration
func (s *SecurityService) ValidateStreamToken(token *StreamToken) error {
	// Check expiration
	if time.Now().After(token.ExpiresAt) {
		return fmt.Errorf("token expired")
	}

	// Verify signature
	expectedSignature, err := s.generateSignature(token)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	if token.Signature != expectedSignature {
		return fmt.Errorf("invalid token signature")
	}

	return nil
}

// generateSignature creates HMAC-SHA256 signature for token
func (s *SecurityService) generateSignature(token *StreamToken) (string, error) {
	// Create message to sign
	message := fmt.Sprintf("%s:%s:%d:%s",
		token.VideoID.String(),
		token.UserID.String(),
		token.ExpiresAt.Unix(),
		token.Quality,
	)

	// Generate HMAC-SHA256
	mac := hmac.New(sha256.New, s.secretKey)
	if _, err := mac.Write([]byte(message)); err != nil {
		return "", err
	}

	// Encode to base64
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return signature, nil
}

// GenerateSignedURL creates a signed URL with embedded token
func (s *SecurityService) GenerateSignedURL(videoID, userID uuid.UUID, quality string) (string, error) {
	token, err := s.GenerateStreamToken(videoID, userID, quality)
	if err != nil {
		return "", err
	}

	// Create signed URL with token parameters
	signedURL := fmt.Sprintf("/api/v1/videos/%s/stream/secure?token=%s&expires=%d&quality=%s&signature=%s",
		videoID.String(),
		base64.URLEncoding.EncodeToString([]byte(userID.String())),
		token.ExpiresAt.Unix(),
		quality,
		token.Signature,
	)

	return signedURL, nil
}

// ValidateSignedURL validates URL parameters and signature
func (s *SecurityService) ValidateSignedURL(videoIDStr, userIDStr, quality, signature string, expiresAt int64) error {
	// Parse IDs
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return fmt.Errorf("invalid video ID: %w", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Create token for validation
	token := &StreamToken{
		VideoID:   videoID,
		UserID:    userID,
		ExpiresAt: time.Unix(expiresAt, 0),
		Quality:   quality,
		Signature: signature,
	}

	// Validate token
	return s.ValidateStreamToken(token)
}

// TokenInfo contains decoded token information
type TokenInfo struct {
	VideoID   uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	Quality   string
}

// DecodeToken extracts token information without validation (for logging)
func (s *SecurityService) DecodeToken(token *StreamToken) *TokenInfo {
	return &TokenInfo{
		VideoID:   token.VideoID,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
		Quality:   token.Quality,
	}
}