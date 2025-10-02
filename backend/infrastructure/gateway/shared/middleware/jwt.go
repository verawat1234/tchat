package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// JWTClaims represents the claims structure for Tchat JWT tokens
// Includes Southeast Asian regional compliance fields
type JWTClaims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"` // TH, SG, ID, MY, PH, VN
	Locale      string `json:"locale"`       // en, th, id, ms, fil, vi
	DeviceID    string `json:"device_id"`
	Platform    string `json:"platform"`    // web, mobile_ios, mobile_android
	SessionID   string `json:"session_id"`
	Scope       string `json:"scope"`       // user, business, admin
	DataRegion  string `json:"data_region"` // Regional data residency compliance
	jwt.RegisteredClaims
}

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	SecretKey        string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	RequiredScopes  []string
	RegionalCompliance bool
	Logger          *logrus.Logger
}

// DefaultJWTConfig returns a default configuration for development
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:       "tchat-dev-secret-key-change-in-production",
		Issuer:         "tchat.sea",
		AccessTokenTTL: 1 * time.Hour,
		RefreshTokenTTL: 30 * 24 * time.Hour, // 30 days
		RequiredScopes: []string{"user"},
		RegionalCompliance: true,
		Logger:         logrus.New(),
	}
}

// JWTMiddleware creates a new JWT authentication middleware
func JWTMiddleware(config *JWTConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultJWTConfig()
	}

	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			config.Logger.WithFields(logrus.Fields{
				"ip":         c.ClientIP(),
				"user_agent": c.GetHeader("User-Agent"),
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
			}).Warn("Missing Authorization header")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// Validate Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			config.Logger.WithFields(logrus.Fields{
				"ip":           c.ClientIP(),
				"auth_header":  authHeader,
				"path":         c.Request.URL.Path,
			}).Warn("Invalid Authorization header format")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"error":       err.Error(),
				"ip":          c.ClientIP(),
				"user_agent":  c.GetHeader("User-Agent"),
				"path":        c.Request.URL.Path,
			}).Error("JWT token validation failed")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			config.Logger.WithFields(logrus.Fields{
				"ip":   c.ClientIP(),
				"path": c.Request.URL.Path,
			}).Error("Invalid JWT claims structure")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
				"code":  "INVALID_CLAIMS",
			})
			c.Abort()
			return
		}

		// Validate issuer
		if claims.Issuer != config.Issuer {
			config.Logger.WithFields(logrus.Fields{
				"expected_issuer": config.Issuer,
				"actual_issuer":   claims.Issuer,
				"user_id":         claims.UserID,
			}).Error("JWT issuer validation failed")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token issuer",
				"code":  "INVALID_ISSUER",
			})
			c.Abort()
			return
		}

		// Southeast Asian regional compliance validation
		if config.RegionalCompliance {
			if err := validateRegionalCompliance(claims, c); err != nil {
				config.Logger.WithFields(logrus.Fields{
					"error":        err.Error(),
					"user_id":      claims.UserID,
					"country_code": claims.CountryCode,
					"data_region":  claims.DataRegion,
				}).Error("Regional compliance validation failed")

				c.JSON(http.StatusForbidden, gin.H{
					"error": err.Error(),
					"code":  "COMPLIANCE_VIOLATION",
				})
				c.Abort()
				return
			}
		}

		// Validate required scopes
		if len(config.RequiredScopes) > 0 {
			hasRequiredScope := false
			for _, requiredScope := range config.RequiredScopes {
				if claims.Scope == requiredScope {
					hasRequiredScope = true
					break
				}
			}

			if !hasRequiredScope {
				config.Logger.WithFields(logrus.Fields{
					"user_id":         claims.UserID,
					"user_scope":      claims.Scope,
					"required_scopes": config.RequiredScopes,
				}).Warn("Insufficient scope for requested operation")

				c.JSON(http.StatusForbidden, gin.H{
					"error": "Insufficient permissions",
					"code":  "INSUFFICIENT_SCOPE",
				})
				c.Abort()
				return
			}
		}

		// Add claims to context for use in handlers
		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)

		// Add user info to Gin context for easy access
		c.Set("user_id", claims.UserID)
		c.Set("country_code", claims.CountryCode)
		c.Set("locale", claims.Locale)
		c.Set("session_id", claims.SessionID)
		c.Set("scope", claims.Scope)
		c.Set("data_region", claims.DataRegion)

		// Log successful authentication for audit trail
		config.Logger.WithFields(logrus.Fields{
			"user_id":      claims.UserID,
			"country_code": claims.CountryCode,
			"session_id":   claims.SessionID,
			"ip":           c.ClientIP(),
			"user_agent":   c.GetHeader("User-Agent"),
			"path":         c.Request.URL.Path,
			"method":       c.Request.Method,
		}).Info("User authenticated successfully")

		c.Next()
	}
}

// validateRegionalCompliance ensures compliance with Southeast Asian data regulations
func validateRegionalCompliance(claims *JWTClaims, c *gin.Context) error {
	// Validate country code is supported
	supportedCountries := map[string]bool{
		"TH": true, // Thailand
		"SG": true, // Singapore
		"ID": true, // Indonesia
		"MY": true, // Malaysia
		"PH": true, // Philippines
		"VN": true, // Vietnam
	}

	if !supportedCountries[claims.CountryCode] {
		return fmt.Errorf("unsupported country code: %s", claims.CountryCode)
	}

	// Validate locale matches country requirements
	countryLocales := map[string][]string{
		"TH": {"th", "en"},
		"SG": {"en"},
		"ID": {"id", "en"},
		"MY": {"ms", "en"},
		"PH": {"fil", "en"},
		"VN": {"vi", "en"},
	}

	validLocales, exists := countryLocales[claims.CountryCode]
	if !exists {
		return fmt.Errorf("no locale mapping for country: %s", claims.CountryCode)
	}

	localeValid := false
	for _, validLocale := range validLocales {
		if claims.Locale == validLocale {
			localeValid = true
			break
		}
	}

	if !localeValid {
		return fmt.Errorf("invalid locale %s for country %s", claims.Locale, claims.CountryCode)
	}

	// Data residency validation - ensure user data is processed in correct region
	if claims.DataRegion == "" {
		return fmt.Errorf("data region not specified for compliance")
	}

	// Regional data processing rules
	regionMapping := map[string]string{
		"TH": "sea-central",  // Thailand - Central SEA
		"SG": "sea-central",  // Singapore - Central SEA
		"ID": "sea-central",  // Indonesia - Central SEA
		"MY": "sea-central",  // Malaysia - Central SEA
		"PH": "sea-east",     // Philippines - East SEA
		"VN": "sea-north",    // Vietnam - North SEA
	}

	expectedRegion := regionMapping[claims.CountryCode]
	if claims.DataRegion != expectedRegion {
		return fmt.Errorf("data region mismatch: expected %s, got %s", expectedRegion, claims.DataRegion)
	}

	return nil
}

// GenerateAccessToken creates a new JWT access token with regional compliance
func GenerateAccessToken(claims *JWTClaims, config *JWTConfig) (string, error) {
	if config == nil {
		config = DefaultJWTConfig()
	}

	// Set standard claims
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    config.Issuer,
		Subject:   claims.UserID,
		ID:        claims.SessionID,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a refresh token for session renewal
func GenerateRefreshToken(claims *JWTClaims, config *JWTConfig) (string, error) {
	if config == nil {
		config = DefaultJWTConfig()
	}

	// Simplified claims for refresh token
	refreshClaims := &JWTClaims{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		Scope:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
			Subject:   claims.UserID,
			ID:        claims.SessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ExtractClaims retrieves JWT claims from Gin context
func ExtractClaims(c *gin.Context) (*JWTClaims, error) {
	claims, exists := c.Request.Context().Value("claims").(*JWTClaims)
	if !exists {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

// RequireScope creates middleware that requires specific scopes
func RequireScope(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ExtractClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "MISSING_CLAIMS",
			})
			c.Abort()
			return
		}

		hasRequiredScope := false
		for _, requiredScope := range requiredScopes {
			if claims.Scope == requiredScope {
				hasRequiredScope = true
				break
			}
		}

		if !hasRequiredScope {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"code":  "INSUFFICIENT_SCOPE",
				"required_scopes": requiredScopes,
				"user_scope": claims.Scope,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}