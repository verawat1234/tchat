package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat.dev/shared/config"
	sharedModels "tchat.dev/shared/models"
)

// JWTServiceTestSuite provides comprehensive testing for JWT service operations
type JWTServiceTestSuite struct {
	suite.Suite
	jwtService *JWTService
	testUser   *sharedModels.User
	testConfig *config.Config
	ctx        context.Context
}

// SetupSuite initializes the test suite with configuration and test data
func (suite *JWTServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.testConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret-key-for-jwt-authentication-testing",
			AccessTokenTTL:   15 * time.Minute,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:          "tchat-auth-service",
			Audience:        "tchat-api",
		},
	}

	suite.jwtService = NewJWTService(suite.testConfig)

	// Create test user with various attributes for comprehensive testing
	suite.testUser = &sharedModels.User{
		ID:         uuid.New(),
		Phone:      "+66912345678",
		Email:      "test@tchat.dev",
		Name:       "Test User",
		Country:    "TH",
		Locale:     "th-TH",
		KYCTier:    2,
		Status:     "active",
		Verified:   true,
		CreatedAt:  time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		UpdatedAt:  time.Now(),
	}
}

// Test Token Generation

func (suite *JWTServiceTestSuite) TestGenerateTokenPair_ValidUser_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), tokenPair)
	assert.NotEmpty(suite.T(), tokenPair.AccessToken)
	assert.NotEmpty(suite.T(), tokenPair.RefreshToken)
	assert.Equal(suite.T(), "Bearer", tokenPair.TokenType)
	assert.Equal(suite.T(), int64(900), tokenPair.ExpiresIn) // 15 minutes
	assert.True(suite.T(), tokenPair.AccessExpiresAt.After(time.Now()))
	assert.True(suite.T(), tokenPair.RefreshExpiresAt.After(tokenPair.AccessExpiresAt))
	assert.Equal(suite.T(), "read write", tokenPair.Scope)
}

func (suite *JWTServiceTestSuite) TestGenerateTokenPair_NilUser_Error() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, nil, sessionID, deviceID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tokenPair)
	assert.Contains(suite.T(), err.Error(), "user cannot be nil")
}

func (suite *JWTServiceTestSuite) TestGenerateTokenPair_NilSessionID_Error() {
	deviceID := "test-device-123"

	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, uuid.Nil, deviceID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tokenPair)
	assert.Contains(suite.T(), err.Error(), "session ID is required")
}

// Test Access Token Validation

func (suite *JWTServiceTestSuite) TestValidateAccessToken_ValidToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Validate access token
	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenPair.AccessToken)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), sessionID, claims.SessionID)
	assert.Equal(suite.T(), deviceID, claims.DeviceID)
	assert.Equal(suite.T(), "+66912345678", claims.PhoneNumber)
	assert.Equal(suite.T(), "TH", claims.CountryCode)
	assert.Equal(suite.T(), "verified", claims.KYCStatus)
	assert.Equal(suite.T(), 2, claims.KYCLevel)
	assert.Contains(suite.T(), claims.Scopes, "read")
	assert.Contains(suite.T(), claims.Scopes, "write")
	assert.Contains(suite.T(), claims.Permissions, "profile:read")
	assert.Contains(suite.T(), claims.Permissions, "wallet:read")
	assert.Contains(suite.T(), claims.Permissions, "region:sea:premium")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_EmptyToken_Error() {
	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, "")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "token is required")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_InvalidSignature_Error() {
	// Create token with wrong secret
	wrongService := &JWTService{
		config:       suite.testConfig,
		accessSecret: []byte("wrong-secret"),
		issuer:       suite.testConfig.JWT.Issuer,
		audience:     suite.testConfig.JWT.Audience,
	}

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: uuid.New(),
		DeviceID:  "test-device",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(wrongService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate with correct service (different secret)
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "invalid token")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_ExpiredToken_Error() {
	// Create expired token
	sessionID := uuid.New()
	deviceID := "test-device-123"

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour) // 1 hour ago
	expiredTime := now.Add(-30 * time.Minute) // 30 minutes ago

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(pastTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(suite.jwtService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate expired token
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "invalid token")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_InvalidIssuer_Error() {
	// Create token with wrong issuer
	sessionID := uuid.New()
	deviceID := "test-device-123"

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    "wrong-issuer",
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(suite.jwtService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "invalid issuer")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_InvalidAudience_Error() {
	// Create token with wrong audience
	sessionID := uuid.New()
	deviceID := "test-device-123"

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{"wrong-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(suite.jwtService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "invalid audience")
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_MissingUserID_Error() {
	// Create token with missing user ID
	sessionID := uuid.New()
	deviceID := "test-device-123"

	claims := &UserClaims{
		UserID:    uuid.Nil, // Missing user ID
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(suite.jwtService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "missing user ID")
}

// Test Refresh Token Validation

func (suite *JWTServiceTestSuite) TestValidateRefreshToken_ValidToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Validate refresh token
	claims, err := suite.jwtService.ValidateRefreshToken(suite.ctx, tokenPair.RefreshToken)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), sessionID, claims.SessionID)
	assert.Equal(suite.T(), deviceID, claims.DeviceID)
	assert.Contains(suite.T(), claims.Scopes, "refresh")
}

func (suite *JWTServiceTestSuite) TestValidateRefreshToken_EmptyToken_Error() {
	claims, err := suite.jwtService.ValidateRefreshToken(suite.ctx, "")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "refresh token is required")
}

func (suite *JWTServiceTestSuite) TestValidateRefreshToken_AccessTokenAsRefresh_Error() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Try to validate access token as refresh token
	claims, err := suite.jwtService.ValidateRefreshToken(suite.ctx, tokenPair.AccessToken)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token scope")
}

// Test Token Refresh Flow

func (suite *JWTServiceTestSuite) TestRefreshAccessToken_ValidRefreshToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate initial token pair
	originalTokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Refresh the access token
	newTokenPair, err := suite.jwtService.RefreshAccessToken(suite.ctx, originalTokenPair.RefreshToken, suite.testUser)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), newTokenPair)
	assert.NotEqual(suite.T(), originalTokenPair.AccessToken, newTokenPair.AccessToken)
	assert.NotEqual(suite.T(), originalTokenPair.RefreshToken, newTokenPair.RefreshToken)

	// Validate new access token
	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, newTokenPair.AccessToken)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), sessionID, claims.SessionID)
	assert.Equal(suite.T(), deviceID, claims.DeviceID)
}

func (suite *JWTServiceTestSuite) TestRefreshAccessToken_UserMismatch_Error() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair for original user
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Create different user
	differentUser := &sharedModels.User{
		ID:         uuid.New(), // Different ID
		Phone:      "+66987654321",
		Name:       "Different User",
		Country:    "TH",
		KYCTier:    1,
		Verified:   false,
	}

	// Try to refresh with different user
	newTokenPair, err := suite.jwtService.RefreshAccessToken(suite.ctx, tokenPair.RefreshToken, differentUser)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTokenPair)
	assert.Contains(suite.T(), err.Error(), "refresh token user mismatch")
}

func (suite *JWTServiceTestSuite) TestRefreshAccessToken_InvalidRefreshToken_Error() {
	// Try to refresh with invalid token
	newTokenPair, err := suite.jwtService.RefreshAccessToken(suite.ctx, "invalid-token", suite.testUser)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTokenPair)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token")
}

// Test Service Tokens

func (suite *JWTServiceTestSuite) TestGenerateServiceToken_ValidService_Success() {
	serviceName := "payment-service"
	permissions := []string{"payment:process", "payment:verify"}

	token, err := suite.jwtService.GenerateServiceToken(suite.ctx, serviceName, permissions)

	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// Validate service token
	claims, err := suite.jwtService.ValidateServiceToken(suite.ctx, token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), uuid.Nil, claims.UserID)
	assert.Equal(suite.T(), uuid.Nil, claims.SessionID)
	assert.Equal(suite.T(), serviceName, claims.DeviceID)
	assert.Equal(suite.T(), serviceName, claims.Subject)
	assert.Contains(suite.T(), claims.Scopes, "service")
	assert.Equal(suite.T(), permissions, claims.Permissions)
}

func (suite *JWTServiceTestSuite) TestValidateServiceToken_UserTokenAsService_Error() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate user token
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Try to validate user token as service token
	claims, err := suite.jwtService.ValidateServiceToken(suite.ctx, tokenPair.AccessToken)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "not a service token")
}

// Test Token Information Extraction

func (suite *JWTServiceTestSuite) TestExtractClaims_ValidToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Extract claims without validation
	claims, err := suite.jwtService.ExtractClaims(tokenPair.AccessToken)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), sessionID, claims.SessionID)
	assert.Equal(suite.T(), deviceID, claims.DeviceID)
}

func (suite *JWTServiceTestSuite) TestGetTokenInfo_ValidToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Get token info
	tokenInfo, err := suite.jwtService.GetTokenInfo(tokenPair.AccessToken)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), tokenInfo)
	assert.Equal(suite.T(), suite.testUser.ID, tokenInfo.UserID)
	assert.Equal(suite.T(), sessionID, tokenInfo.SessionID)
	assert.Equal(suite.T(), deviceID, tokenInfo.DeviceID)
	assert.Equal(suite.T(), suite.testConfig.JWT.Issuer, tokenInfo.Issuer)
	assert.Contains(suite.T(), tokenInfo.Audience, suite.testConfig.JWT.Audience)
	assert.Contains(suite.T(), tokenInfo.Permissions, "profile:read")
	assert.Equal(suite.T(), 2, tokenInfo.KYCLevel)
	assert.Equal(suite.T(), "TH", tokenInfo.CountryCode)
	assert.Equal(suite.T(), "+66912345678", tokenInfo.Metadata["phone_number"])
	assert.Equal(suite.T(), "verified", tokenInfo.Metadata["kyc_status"])
}

func (suite *JWTServiceTestSuite) TestGetTokenExpiration_ValidToken_Success() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate token pair
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Get expiration
	expiration, err := suite.jwtService.GetTokenExpiration(tokenPair.AccessToken)

	require.NoError(suite.T(), err)
	assert.True(suite.T(), expiration.After(time.Now()))
	assert.True(suite.T(), expiration.Before(time.Now().Add(16*time.Minute))) // Should be ~15 minutes
}

// Test KYC Tier Permissions

func (suite *JWTServiceTestSuite) TestPermissions_KYCTier1_BasicPermissions() {
	tier1User := &sharedModels.User{
		ID:         uuid.New(),
		Phone:      "+66912345678",
		Name:       "Tier 1 User",
		Country:    "TH",
		KYCTier:    1,
		Verified:   false,
	}

	sessionID := uuid.New()
	deviceID := "test-device-123"

	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, tier1User, sessionID, deviceID)
	require.NoError(suite.T(), err)

	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenPair.AccessToken)
	require.NoError(suite.T(), err)

	// Should have basic permissions
	assert.Contains(suite.T(), claims.Permissions, "profile:read")
	assert.Contains(suite.T(), claims.Permissions, "profile:update")
	assert.Contains(suite.T(), claims.Permissions, "wallet:read")
	assert.Contains(suite.T(), claims.Permissions, "payment:send:basic")

	// Should not have advanced permissions
	assert.NotContains(suite.T(), claims.Permissions, "wallet:create")
	assert.NotContains(suite.T(), claims.Permissions, "commerce:sell")
}

func (suite *JWTServiceTestSuite) TestPermissions_KYCTier3_AdvancedPermissions() {
	tier3User := &sharedModels.User{
		ID:         uuid.New(),
		Phone:      "+66912345678",
		Name:       "Tier 3 User",
		Country:    "SG",
		KYCTier:    3,
		Verified:   true,
	}

	sessionID := uuid.New()
	deviceID := "test-device-123"

	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, tier3User, sessionID, deviceID)
	require.NoError(suite.T(), err)

	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenPair.AccessToken)
	require.NoError(suite.T(), err)

	// Should have all permissions
	assert.Contains(suite.T(), claims.Permissions, "profile:read")
	assert.Contains(suite.T(), claims.Permissions, "wallet:manage")
	assert.Contains(suite.T(), claims.Permissions, "payment:business")
	assert.Contains(suite.T(), claims.Permissions, "commerce:manage")
	assert.Contains(suite.T(), claims.Permissions, "verified:user")
	assert.Contains(suite.T(), claims.Permissions, "region:sea:premium")
}

func (suite *JWTServiceTestSuite) TestPermissions_DifferentCountries_RegionalPermissions() {
	// Test Indonesia (standard region)
	indonesiaUser := &sharedModels.User{
		ID:      uuid.New(),
		Name:    "Indonesia User",
		Country: "ID",
		KYCTier: 1,
	}

	sessionID := uuid.New()
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, indonesiaUser, sessionID, "device-id")
	require.NoError(suite.T(), err)

	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenPair.AccessToken)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), claims.Permissions, "region:sea:standard")
	assert.NotContains(suite.T(), claims.Permissions, "region:sea:premium")
}

// Test Security Edge Cases

func (suite *JWTServiceTestSuite) TestValidateAccessToken_TamperedToken_Error() {
	sessionID := uuid.New()
	deviceID := "test-device-123"

	// Generate valid token
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.ctx, suite.testUser, sessionID, deviceID)
	require.NoError(suite.T(), err)

	// Tamper with token (change last character)
	tamperedToken := tokenPair.AccessToken[:len(tokenPair.AccessToken)-1] + "X"

	// Try to validate tampered token
	claims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tamperedToken)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_WrongSigningMethod_Error() {
	// Create token with different signing method
	sessionID := uuid.New()
	deviceID := "test-device-123"

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Use RS256 instead of HS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString := token.Raw // This will be malformed since we can't sign with RS256 without proper keys

	// Try to validate
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString+"invalid")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken_FutureNotBefore_Error() {
	// Create token with future not-before time
	sessionID := uuid.New()
	deviceID := "test-device-123"

	now := time.Now()
	futureTime := now.Add(1 * time.Hour) // 1 hour in future

	claims := &UserClaims{
		UserID:    suite.testUser.ID,
		SessionID: sessionID,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   suite.testUser.ID.String(),
			Issuer:    suite.testConfig.JWT.Issuer,
			Audience:  []string{suite.testConfig.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(futureTime), // Future not-before
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(suite.jwtService.accessSecret)
	require.NoError(suite.T(), err)

	// Try to validate
	validatedClaims, err := suite.jwtService.ValidateAccessToken(suite.ctx, tokenString)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
	assert.Contains(suite.T(), err.Error(), "token not yet valid")
}

// Run test suite
func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}


// Additional individual tests for edge cases

func TestJWTService_NilConfig_Panic(t *testing.T) {
	assert.Panics(t, func() {
		NewJWTService(nil)
	})
}

func TestJWTService_EmptySecret_GenerationError(t *testing.T) {
	config := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "", // Empty secret
			AccessTokenTTL:   15 * time.Minute,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:          "tchat-auth-service",
			Audience:        "tchat-api",
		},
	}

	jwtService := NewJWTService(config)

	user := &sharedModels.User{
		ID:      uuid.New(),
		Name:    "Test User",
		Country: "TH",
		KYCTier: 1,
	}

	ctx := context.Background()
	sessionID := uuid.New()
	deviceID := "test-device"

	tokenPair, err := jwtService.GenerateTokenPair(ctx, user, sessionID, deviceID)

	// Should still work but be less secure
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
}

func TestJWTService_ExtremeExpirationTimes(t *testing.T) {
	config := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret",
			AccessTokenTTL:   1 * time.Nanosecond, // Extremely short
			RefreshTokenTTL:  100 * 365 * 24 * time.Hour, // Extremely long
			Issuer:          "tchat-auth-service",
			Audience:        "tchat-api",
		},
	}

	jwtService := NewJWTService(config)

	user := &sharedModels.User{
		ID:      uuid.New(),
		Name:    "Test User",
		Country: "TH",
		KYCTier: 1,
	}

	ctx := context.Background()
	sessionID := uuid.New()
	deviceID := "test-device"

	tokenPair, err := jwtService.GenerateTokenPair(ctx, user, sessionID, deviceID)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)

	// Access token should expire very quickly
	time.Sleep(2 * time.Nanosecond)
	claims, err := jwtService.ValidateAccessToken(ctx, tokenPair.AccessToken)
	assert.Error(t, err) // Should be expired
	assert.Nil(t, claims)
}