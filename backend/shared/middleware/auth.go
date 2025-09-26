package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"tchat.dev/shared/config"
)

// UserClaims represents the JWT claims for a user
type UserClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	CountryCode string    `json:"country_code"`
	KYCStatus   string    `json:"kyc_status"`
	KYCLevel    int       `json:"kyc_level"`
	SessionID   uuid.UUID `json:"session_id"`
	DeviceID    string    `json:"device_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// RequireAuth enforces JWT authentication for protected routes
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := am.extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Missing or invalid authorization header",
				"code":    "AUTH_MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		claims, err := am.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
				"code":    "AUTH_INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Set user context
		am.setUserContext(c, claims)
		c.Next()
	}
}

// RequireKYC enforces KYC verification level for protected routes
func (am *AuthMiddleware) RequireKYC(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := am.getUserClaims(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authentication required",
				"code":    "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		if claims.KYCLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Insufficient KYC verification level",
				"code":    "KYC_INSUFFICIENT_LEVEL",
				"required_level": minLevel,
				"current_level":  claims.KYCLevel,
			})
			c.Abort()
			return
		}

		if claims.KYCStatus != "verified" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "KYC verification required",
				"code":    "KYC_VERIFICATION_REQUIRED",
				"kyc_status": claims.KYCStatus,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCountry enforces country-specific access restrictions
func (am *AuthMiddleware) RequireCountry(allowedCountries ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := am.getUserClaims(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authentication required",
				"code":    "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		countryAllowed := false
		for _, country := range allowedCountries {
			if claims.CountryCode == country {
				countryAllowed = true
				break
			}
		}

		if !countryAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Access restricted for your country",
				"code":    "COUNTRY_RESTRICTED",
				"country": claims.CountryCode,
				"allowed_countries": allowedCountries,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth provides optional authentication (doesn't block on missing token)
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := am.extractToken(c)
		if err != nil {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		claims, err := am.validateToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user context if valid token provided
		am.setUserContext(c, claims)
		c.Next()
	}
}

// RefreshToken handles JWT token refresh
func (am *AuthMiddleware) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := am.extractRefreshToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Missing or invalid refresh token",
				"code":    "REFRESH_TOKEN_MISSING",
			})
			c.Abort()
			return
		}

		claims, err := am.validateRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired refresh token",
				"code":    "REFRESH_TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Generate new access token
		newAccessToken, err := am.generateAccessToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to generate new access token",
				"code":    "TOKEN_GENERATION_FAILED",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": newAccessToken,
			"token_type":   "Bearer",
			"expires_in":   int(am.config.JWT.AccessTokenTTL.Seconds()),
		})
	}
}

// extractToken extracts JWT token from Authorization header
func (am *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", jwt.ErrTokenMalformed
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", jwt.ErrTokenMalformed
	}

	return parts[1], nil
}

// extractRefreshToken extracts refresh token from request body or cookie
func (am *AuthMiddleware) extractRefreshToken(c *gin.Context) (string, error) {
	// Try to get from request body first
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err == nil && request.RefreshToken != "" {
		return request.RefreshToken, nil
	}

	// Try to get from cookie
	if refreshToken, err := c.Cookie("refresh_token"); err == nil && refreshToken != "" {
		return refreshToken, nil
	}

	return "", jwt.ErrTokenMalformed
}

// validateToken validates and parses JWT access token
func (am *AuthMiddleware) validateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(am.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// Verify token hasn't expired
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, jwt.ErrTokenExpired
		}

		// Verify issuer and audience
		if claims.Issuer != am.config.JWT.Issuer {
			return nil, jwt.ErrTokenMalformed
		}

		if claims.Audience[0] != am.config.JWT.Audience {
			return nil, jwt.ErrTokenMalformed
		}

		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// validateRefreshToken validates refresh token
func (am *AuthMiddleware) validateRefreshToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(am.config.JWT.Secret + "_refresh"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// generateAccessToken generates new access token from refresh claims
func (am *AuthMiddleware) generateAccessToken(refreshClaims *UserClaims) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserID:      refreshClaims.UserID,
		PhoneNumber: refreshClaims.PhoneNumber,
		CountryCode: refreshClaims.CountryCode,
		KYCStatus:   refreshClaims.KYCStatus,
		KYCLevel:    refreshClaims.KYCLevel,
		SessionID:   refreshClaims.SessionID,
		DeviceID:    refreshClaims.DeviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   refreshClaims.UserID.String(),
			Issuer:    am.config.JWT.Issuer,
			Audience:  []string{am.config.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(am.config.JWT.AccessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(am.config.JWT.Secret))
}

// setUserContext sets user information in Gin context
func (am *AuthMiddleware) setUserContext(c *gin.Context, claims *UserClaims) {
	c.Set("user_id", claims.UserID)
	c.Set("phone_number", claims.PhoneNumber)
	c.Set("country_code", claims.CountryCode)
	c.Set("kyc_status", claims.KYCStatus)
	c.Set("kyc_level", claims.KYCLevel)
	c.Set("session_id", claims.SessionID)
	c.Set("device_id", claims.DeviceID)
	c.Set("user_claims", claims)

	// Add to request context for use in other parts of the application
	ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_claims", claims)
	c.Request = c.Request.WithContext(ctx)
}

// getUserClaims retrieves user claims from Gin context
func (am *AuthMiddleware) getUserClaims(c *gin.Context) (*UserClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*UserClaims)
	return userClaims, ok
}

// Helper functions for accessing user data from context

// GetUserID extracts user ID from Gin context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserClaims extracts full user claims from Gin context
func GetUserClaims(c *gin.Context) (*UserClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*UserClaims)
	return userClaims, ok
}

// GetCountryCode extracts country code from Gin context
func GetCountryCode(c *gin.Context) (string, bool) {
	countryCode, exists := c.Get("country_code")
	if !exists {
		return "", false
	}

	code, ok := countryCode.(string)
	return code, ok
}

// RequiredAuth is a convenience middleware for most protected routes
func RequiredAuth(cfg *config.Config) gin.HandlerFunc {
	return NewAuthMiddleware(cfg).RequireAuth()
}

// RequiredKYC is a convenience middleware for KYC-protected routes
func RequiredKYC(cfg *config.Config, minLevel int) gin.HandlerFunc {
	middleware := NewAuthMiddleware(cfg)
	return gin.HandlerFunc(func(c *gin.Context) {
		// First require authentication
		middleware.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Then check KYC level
		middleware.RequireKYC(minLevel)(c)
	})
}

// SEARegionAuth restricts access to Southeast Asian countries only
func SEARegionAuth(cfg *config.Config) gin.HandlerFunc {
	middleware := NewAuthMiddleware(cfg)
	return gin.HandlerFunc(func(c *gin.Context) {
		// First require authentication
		middleware.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Then check country restriction
		middleware.RequireCountry("TH", "SG", "ID", "MY", "PH", "VN")(c)
	})
}