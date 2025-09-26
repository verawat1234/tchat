package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"tchat.dev/auth/models"
	sharedModels "tchat.dev/shared/models"
	"tchat.dev/shared/config"
)

// JWTService handles JWT token generation, validation, and management
type JWTService struct {
	config          *config.Config
	accessSecret    []byte
	refreshSecret   []byte
	issuer          string
	audience        string
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
}

// NewJWTService creates a new JWT service instance
func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		config:        cfg,
		accessSecret:  []byte(cfg.JWT.Secret),
		refreshSecret: []byte(cfg.JWT.Secret + "_refresh"),
		issuer:        cfg.JWT.Issuer,
		audience:      cfg.JWT.Audience,
		accessExpiry:  cfg.JWT.AccessTokenTTL,
		refreshExpiry: cfg.JWT.RefreshTokenTTL,
	}
}

// UserClaims represents the JWT claims for a user
type UserClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	CountryCode string    `json:"country_code"`
	KYCStatus   string    `json:"kyc_status"`
	KYCLevel    int       `json:"kyc_level"`
	SessionID   uuid.UUID `json:"session_id"`
	DeviceID    string    `json:"device_id"`
	Permissions []string  `json:"permissions,omitempty"`
	Scopes      []string  `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Scope            string    `json:"scope,omitempty"`
}

// GenerateTokenPair creates both access and refresh tokens for a user
func (js *JWTService) GenerateTokenPair(ctx context.Context, user *sharedModels.User, sessionID uuid.UUID, deviceID string) (*TokenPair, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	if sessionID == uuid.Nil {
		return nil, errors.New("session ID is required")
	}

	now := time.Now()
	accessExpiresAt := now.Add(js.accessExpiry)
	refreshExpiresAt := now.Add(js.refreshExpiry)

	// Generate access token
	accessToken, err := js.generateAccessToken(user, sessionID, deviceID, now, accessExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := js.generateRefreshToken(user, sessionID, deviceID, now, refreshExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(js.accessExpiry.Seconds()),
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Scope:            "read write",
	}, nil
}

// generateAccessToken creates a new access token
func (js *JWTService) generateAccessToken(user *sharedModels.User, sessionID uuid.UUID, deviceID string, issuedAt, expiresAt time.Time) (string, error) {
	claims := &UserClaims{
		UserID:      user.ID,
		PhoneNumber: getPhoneNumber(user),
		CountryCode: user.CountryCode,
		KYCStatus:   getKYCStatus(user),
		KYCLevel:    int(user.KYCTier),
		SessionID:   sessionID,
		DeviceID:    deviceID,
		Permissions: js.getUserPermissions(user),
		Scopes:      []string{"read", "write"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   user.ID.String(),
			Issuer:    js.issuer,
			Audience:  []string{js.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.accessSecret)
}

// generateRefreshToken creates a new refresh token
func (js *JWTService) generateRefreshToken(user *sharedModels.User, sessionID uuid.UUID, deviceID string, issuedAt, expiresAt time.Time) (string, error) {
	claims := &UserClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		Scopes:    []string{"refresh"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   user.ID.String(),
			Issuer:    js.issuer,
			Audience:  []string{js.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.refreshSecret)
}

// ValidateAccessToken validates and parses an access token
func (js *JWTService) ValidateAccessToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return js.accessSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Additional validation
	if err := js.validateClaims(claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (js *JWTService) ValidateRefreshToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	if tokenString == "" {
		return nil, errors.New("refresh token is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return js.refreshSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token claims")
	}

	// Verify this is a refresh token
	if !js.hasScope(claims.Scopes, "refresh") {
		return nil, errors.New("invalid refresh token scope")
	}

	// Additional validation
	if err := js.validateClaims(claims); err != nil {
		return nil, fmt.Errorf("invalid refresh token claims: %w", err)
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (js *JWTService) RefreshAccessToken(ctx context.Context, refreshTokenString string, user *sharedModels.User) (*TokenPair, error) {
	// Validate refresh token
	refreshClaims, err := js.ValidateRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify user match
	if refreshClaims.UserID != user.ID {
		return nil, errors.New("refresh token user mismatch")
	}

	// Generate new token pair
	return js.GenerateTokenPair(ctx, user, refreshClaims.SessionID, refreshClaims.DeviceID)
}

// ExtractClaims extracts claims from a token without validation (for debugging)
func (js *JWTService) ExtractClaims(tokenString string) (*UserClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenString, &UserClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token format")
	}

	return claims, nil
}

// RevokeToken adds a token to the revocation list (would require Redis/DB implementation)
func (js *JWTService) RevokeToken(ctx context.Context, tokenID string) error {
	// TODO: Implement token revocation using Redis or database
	// For now, we rely on session management for revocation
	return nil
}

// IsTokenRevoked checks if a token has been revoked (would require Redis/DB implementation)
func (js *JWTService) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	// TODO: Implement token revocation check using Redis or database
	// For now, we rely on session management
	return false, nil
}

// GetTokenExpiration returns the expiration time of a token
func (js *JWTService) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := js.ExtractClaims(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration")
	}

	return claims.ExpiresAt.Time, nil
}

// Helper functions

// validateClaims performs additional validation on JWT claims
func (js *JWTService) validateClaims(claims *UserClaims) error {
	// Verify issuer
	if claims.Issuer != js.issuer {
		return errors.New("invalid issuer")
	}

	// Verify audience
	if len(claims.Audience) == 0 || claims.Audience[0] != js.audience {
		return errors.New("invalid audience")
	}

	// Verify user ID
	if claims.UserID == uuid.Nil {
		return errors.New("missing user ID")
	}

	// Verify session ID
	if claims.SessionID == uuid.Nil {
		return errors.New("missing session ID")
	}

	// Verify expiration
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return errors.New("token has expired")
	}

	// Verify not before
	if claims.NotBefore != nil && time.Now().Before(claims.NotBefore.Time) {
		return errors.New("token not yet valid")
	}

	return nil
}

// getUserPermissions returns permissions based on user's KYC tier and status
func (js *JWTService) getUserPermissions(user *sharedModels.User) []string {
	permissions := []string{"profile:read", "profile:update"}

	// Add permissions based on KYC tier
	switch int(user.KYCTier) {
	case models.KYCTier1:
		permissions = append(permissions, "wallet:read", "payment:send:basic")
	case models.KYCTier2:
		permissions = append(permissions,
			"wallet:read", "wallet:create",
			"payment:send", "payment:receive",
			"commerce:browse")
	case models.KYCTier3:
		permissions = append(permissions,
			"wallet:read", "wallet:create", "wallet:manage",
			"payment:send", "payment:receive", "payment:business",
			"commerce:browse", "commerce:sell", "commerce:manage")
	}

	// Add verification-based permissions
	if user.PhoneVerified || user.EmailVerified {
		permissions = append(permissions, "verified:user")
	}

	// Add country-specific permissions
	switch models.Country(user.CountryCode) {
	case models.CountryThailand, models.CountrySingapore:
		permissions = append(permissions, "region:sea:premium")
	default:
		permissions = append(permissions, "region:sea:standard")
	}

	return permissions
}

// getPhoneNumber safely extracts phone number from user
func getPhoneNumber(user *sharedModels.User) string {
	return user.PhoneNumber
}

// getKYCStatus returns KYC verification status
func getKYCStatus(user *sharedModels.User) string {
	if user.PhoneVerified || user.EmailVerified {
		return "verified"
	}
	return "pending"
}

// hasScope checks if a scope exists in the scopes list
func (js *JWTService) hasScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// TokenInfo represents information about a token for debugging/admin purposes
type TokenInfo struct {
	TokenID     string            `json:"token_id"`
	UserID      uuid.UUID         `json:"user_id"`
	SessionID   uuid.UUID         `json:"session_id"`
	DeviceID    string            `json:"device_id"`
	Issuer      string            `json:"issuer"`
	Audience    []string          `json:"audience"`
	IssuedAt    time.Time         `json:"issued_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
	NotBefore   time.Time         `json:"not_before"`
	Permissions []string          `json:"permissions"`
	Scopes      []string          `json:"scopes"`
	KYCLevel    int               `json:"kyc_level"`
	CountryCode string            `json:"country_code"`
	Metadata    map[string]string `json:"metadata"`
}

// GetTokenInfo returns detailed information about a token
func (js *JWTService) GetTokenInfo(tokenString string) (*TokenInfo, error) {
	claims, err := js.ExtractClaims(tokenString)
	if err != nil {
		return nil, err
	}

	info := &TokenInfo{
		TokenID:     claims.ID,
		UserID:      claims.UserID,
		SessionID:   claims.SessionID,
		DeviceID:    claims.DeviceID,
		Issuer:      claims.Issuer,
		Audience:    claims.Audience,
		Permissions: claims.Permissions,
		Scopes:      claims.Scopes,
		KYCLevel:    claims.KYCLevel,
		CountryCode: claims.CountryCode,
		Metadata:    make(map[string]string),
	}

	if claims.IssuedAt != nil {
		info.IssuedAt = claims.IssuedAt.Time
	}

	if claims.ExpiresAt != nil {
		info.ExpiresAt = claims.ExpiresAt.Time
	}

	if claims.NotBefore != nil {
		info.NotBefore = claims.NotBefore.Time
	}

	// Add metadata
	info.Metadata["phone_number"] = claims.PhoneNumber
	info.Metadata["kyc_status"] = claims.KYCStatus

	return info, nil
}

// GenerateServiceToken creates a token for service-to-service communication
func (js *JWTService) GenerateServiceToken(ctx context.Context, serviceName string, permissions []string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour) // Service tokens expire in 1 hour

	claims := &UserClaims{
		UserID:      uuid.Nil, // Service tokens don't have user ID
		SessionID:   uuid.Nil, // Service tokens don't have session ID
		DeviceID:    serviceName,
		Permissions: permissions,
		Scopes:      []string{"service"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   serviceName,
			Issuer:    js.issuer,
			Audience:  []string{"tchat-services"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.accessSecret)
}

// ValidateServiceToken validates a service-to-service token
func (js *JWTService) ValidateServiceToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	claims, err := js.ValidateAccessToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	// Verify this is a service token
	if !js.hasScope(claims.Scopes, "service") {
		return nil, errors.New("not a service token")
	}

	// Service tokens should have nil user and session IDs
	if claims.UserID != uuid.Nil || claims.SessionID != uuid.Nil {
		return nil, errors.New("invalid service token format")
	}

	return claims, nil
}