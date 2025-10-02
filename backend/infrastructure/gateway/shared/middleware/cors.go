package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS middleware configuration for Southeast Asian markets
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
	RegionalDomains  bool // Enable regional domain support for SEA
}

// DefaultCORSConfig returns a default CORS configuration for development
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",           // Web development
			"http://localhost:5173",           // Vite development
			"https://tchat.sea",               // Production domain
			"https://*.tchat.sea",             // Regional subdomains
			"https://tchat.co.th",             // Thailand
			"https://tchat.com.sg",            // Singapore
			"https://tchat.co.id",             // Indonesia
			"https://tchat.com.my",            // Malaysia
			"https://tchat.com.ph",            // Philippines
			"https://tchat.com.vn",            // Vietnam
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"X-Requested-With",
			"X-CSRF-Token",
			"X-Request-ID",
			"X-Country-Code",      // Southeast Asian market support
			"X-Locale",           // Language preference
			"X-Timezone",         // Regional timezone
			"X-Device-ID",        // Device tracking
			"X-App-Version",      // Mobile app version
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-Response-Time",
			"X-Rate-Limit-Remaining",
			"X-Rate-Limit-Reset",
		},
		AllowCredentials: true,
		MaxAge:          86400, // 24 hours
		RegionalDomains: true,
	}
}

// CORSMiddleware creates a CORS middleware with Southeast Asian regional support
func CORSMiddleware(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowed := false
		if len(config.AllowOrigins) == 0 {
			allowed = true
		} else {
			for _, allowedOrigin := range config.AllowOrigins {
				if allowedOrigin == "*" {
					allowed = true
					break
				}

				// Support wildcard subdomains for regional domains
				if strings.Contains(allowedOrigin, "*") {
					pattern := strings.Replace(allowedOrigin, "*", "", -1)
					if strings.Contains(origin, pattern) {
						allowed = true
						break
					}
				} else if allowedOrigin == origin {
					allowed = true
					break
				}
			}
		}

		// Regional domain validation for Southeast Asian markets
		if config.RegionalDomains && !allowed {
			regionalDomains := []string{
				".tchat.co.th",    // Thailand
				".tchat.com.sg",   // Singapore
				".tchat.co.id",    // Indonesia
				".tchat.com.my",   // Malaysia
				".tchat.com.ph",   // Philippines
				".tchat.com.vn",   // Vietnam
			}

			for _, domain := range regionalDomains {
				if strings.HasSuffix(origin, domain) {
					allowed = true
					break
				}
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Set CORS headers
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		if len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}