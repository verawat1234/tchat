package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	// Environment
	Environment string
	Debug       bool
	LogLevel    string
	Version     string

	// Server
	Server ServerConfig

	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// JWT
	JWT JWTConfig

	// Email
	Email EmailConfig

	// SMS
	SMS SMSConfig

	// Storage
	Storage StorageConfig

	// Payment
	Payment PaymentConfig

	// External APIs
	ExternalAPIs ExternalAPIsConfig

	// Rate Limiting
	RateLimit RateLimitConfig

	// Monitoring
	Monitoring MonitoringConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLS          TLSConfig
	CORS         CORSConfig
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	Issuer          string
	Audience        string
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider string
	SMTP     SMTPConfig
	SendGrid SendGridConfig
}

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SendGridConfig holds SendGrid configuration
type SendGridConfig struct {
	APIKey string
	From   string
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider string
	Twilio   TwilioConfig
}

// TwilioConfig holds Twilio configuration
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	Provider string
	Local    LocalStorageConfig
	S3       S3Config
}

// LocalStorageConfig holds local storage configuration
type LocalStorageConfig struct {
	Path string
}

// S3Config holds S3 configuration
type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	UseSSL          bool
}

// PaymentConfig holds payment service configuration
type PaymentConfig struct {
	Provider string
	Stripe   StripeConfig
	PayPal   PayPalConfig
}

// StripeConfig holds Stripe configuration
type StripeConfig struct {
	PublicKey    string
	SecretKey    string
	WebhookSecret string
}

// PayPalConfig holds PayPal configuration
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	Mode         string // sandbox or live
}

// ExternalAPIsConfig holds external API configurations
type ExternalAPIsConfig struct {
	GoogleMaps GoogleMapsConfig
	FCM        FCMConfig
}

// GoogleMapsConfig holds Google Maps API configuration
type GoogleMapsConfig struct {
	APIKey string
}

// FCMConfig holds Firebase Cloud Messaging configuration
type FCMConfig struct {
	ServerKey    string
	SenderID     string
	ProjectID    string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool
	RequestsPerMinute int
	BurstSize   int
	CleanupInterval time.Duration
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Enabled    bool
	MetricsPort int
	Jaeger     JaegerConfig
	Prometheus PrometheusConfig
}

// JaegerConfig holds Jaeger tracing configuration
type JaegerConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
}

// PrometheusConfig holds Prometheus metrics configuration
type PrometheusConfig struct {
	Enabled bool
	Path    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Environment
	config.Environment = getEnv("ENVIRONMENT", "development")
	config.Debug = getBoolEnv("DEBUG", true)
	config.LogLevel = getEnv("LOG_LEVEL", "info")
	config.Version = getEnv("APP_VERSION", "1.0.0")

	// Server
	config.Server = ServerConfig{
		Host:         getEnv("SERVER_HOST", "localhost"),
		Port:         getIntEnv("SERVER_PORT", 8080),
		ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		TLS: TLSConfig{
			Enabled:  getBoolEnv("TLS_ENABLED", false),
			CertFile: getEnv("TLS_CERT_FILE", ""),
			KeyFile:  getEnv("TLS_KEY_FILE", ""),
		},
		CORS: CORSConfig{
			AllowOrigins:     strings.Split(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"), ","),
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-User-ID", "X-Country-Code"},
			ExposeHeaders:    []string{"X-Request-ID", "X-Rate-Limit-Remaining", "X-Rate-Limit-Reset"},
			AllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getIntEnv("CORS_MAX_AGE", 3600),
		},
	}

	// Database
	config.Database = DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getIntEnv("DB_PORT", 5432),
		Username:        getEnv("DB_USERNAME", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		Database:        getEnv("DB_DATABASE", "tchat"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	// Redis
	config.Redis = RedisConfig{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getIntEnv("REDIS_PORT", 6379),
		Password:     getEnv("REDIS_PASSWORD", ""),
		Database:     getIntEnv("REDIS_DATABASE", 0),
		PoolSize:     getIntEnv("REDIS_POOL_SIZE", 10),
		DialTimeout:  getDurationEnv("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getDurationEnv("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: getDurationEnv("REDIS_WRITE_TIMEOUT", 3*time.Second),
	}

	// JWT
	config.JWT = JWTConfig{
		Secret:          getEnv("JWT_SECRET", "your-secret-key"),
		AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TOKEN_TTL", 30*24*time.Hour),
		Issuer:          getEnv("JWT_ISSUER", "tchat"),
		Audience:        getEnv("JWT_AUDIENCE", "tchat-users"),
	}

	// Email
	config.Email = EmailConfig{
		Provider: getEnv("EMAIL_PROVIDER", "smtp"),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getIntEnv("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@tchat.dev"),
		},
		SendGrid: SendGridConfig{
			APIKey: getEnv("SENDGRID_API_KEY", ""),
			From:   getEnv("SENDGRID_FROM", "noreply@tchat.dev"),
		},
	}

	// SMS
	config.SMS = SMSConfig{
		Provider: getEnv("SMS_PROVIDER", "twilio"),
		Twilio: TwilioConfig{
			AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
			AuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
			From:       getEnv("TWILIO_FROM", ""),
		},
	}

	// Storage
	config.Storage = StorageConfig{
		Provider: getEnv("STORAGE_PROVIDER", "local"),
		Local: LocalStorageConfig{
			Path: getEnv("STORAGE_LOCAL_PATH", "./uploads"),
		},
		S3: S3Config{
			Bucket:          getEnv("S3_BUCKET", ""),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			UseSSL:          getBoolEnv("S3_USE_SSL", true),
		},
	}

	// Payment
	config.Payment = PaymentConfig{
		Provider: getEnv("PAYMENT_PROVIDER", "stripe"),
		Stripe: StripeConfig{
			PublicKey:     getEnv("STRIPE_PUBLIC_KEY", ""),
			SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
		PayPal: PayPalConfig{
			ClientID:     getEnv("PAYPAL_CLIENT_ID", ""),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", ""),
			Mode:         getEnv("PAYPAL_MODE", "sandbox"),
		},
	}

	// External APIs
	config.ExternalAPIs = ExternalAPIsConfig{
		GoogleMaps: GoogleMapsConfig{
			APIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		},
		FCM: FCMConfig{
			ServerKey: getEnv("FCM_SERVER_KEY", ""),
			SenderID:  getEnv("FCM_SENDER_ID", ""),
			ProjectID: getEnv("FCM_PROJECT_ID", ""),
		},
	}

	// Rate Limiting
	config.RateLimit = RateLimitConfig{
		Enabled:           getBoolEnv("RATE_LIMIT_ENABLED", true),
		RequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
		BurstSize:         getIntEnv("RATE_LIMIT_BURST_SIZE", 10),
		CleanupInterval:   getDurationEnv("RATE_LIMIT_CLEANUP_INTERVAL", 5*time.Minute),
	}

	// Monitoring
	config.Monitoring = MonitoringConfig{
		Enabled:     getBoolEnv("MONITORING_ENABLED", false),
		MetricsPort: getIntEnv("MONITORING_METRICS_PORT", 9090),
		Jaeger: JaegerConfig{
			Enabled:     getBoolEnv("JAEGER_ENABLED", false),
			Endpoint:    getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName: getEnv("JAEGER_SERVICE_NAME", "tchat"),
		},
		Prometheus: PrometheusConfig{
			Enabled: getBoolEnv("PROMETHEUS_ENABLED", false),
			Path:    getEnv("PROMETHEUS_PATH", "/metrics"),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	var errors []string

	// Validate required fields based on environment
	if c.Environment == "production" {
		// Production-specific validations
		if c.JWT.Secret == "your-secret-key" {
			errors = append(errors, "JWT_SECRET must be set in production")
		}

		if c.Database.Password == "" {
			errors = append(errors, "DB_PASSWORD must be set in production")
		}

		// Validate external service configurations if providers are set
		if c.Email.Provider == "sendgrid" && c.Email.SendGrid.APIKey == "" {
			errors = append(errors, "SENDGRID_API_KEY must be set when using SendGrid")
		}

		if c.SMS.Provider == "twilio" && (c.SMS.Twilio.AccountSID == "" || c.SMS.Twilio.AuthToken == "") {
			errors = append(errors, "Twilio credentials must be set when using Twilio")
		}

		if c.Storage.Provider == "s3" && (c.Storage.S3.AccessKeyID == "" || c.Storage.S3.SecretAccessKey == "") {
			errors = append(errors, "S3 credentials must be set when using S3")
		}

		if c.Payment.Provider == "stripe" && c.Payment.Stripe.SecretKey == "" {
			errors = append(errors, "STRIPE_SECRET_KEY must be set when using Stripe")
		}
	}

	// Validate port ranges
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errors = append(errors, "SERVER_PORT must be between 1 and 65535")
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		errors = append(errors, "DB_PORT must be between 1 and 65535")
	}

	if c.Redis.Port < 1 || c.Redis.Port > 65535 {
		errors = append(errors, "REDIS_PORT must be between 1 and 65535")
	}

	// Validate timeouts
	if c.Server.ReadTimeout <= 0 {
		errors = append(errors, "SERVER_READ_TIMEOUT must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		errors = append(errors, "SERVER_WRITE_TIMEOUT must be positive")
	}

	if c.JWT.AccessTokenTTL <= 0 {
		errors = append(errors, "JWT_ACCESS_TOKEN_TTL must be positive")
	}

	if c.JWT.RefreshTokenTTL <= 0 {
		errors = append(errors, "JWT_REFRESH_TOKEN_TTL must be positive")
	}

	// Validate database connection limits
	if c.Database.MaxOpenConns <= 0 {
		errors = append(errors, "DB_MAX_OPEN_CONNS must be positive")
	}

	if c.Database.MaxIdleConns <= 0 {
		errors = append(errors, "DB_MAX_IDLE_CONNS must be positive")
	}

	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		errors = append(errors, "DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	// Validate rate limiting
	if c.RateLimit.Enabled {
		if c.RateLimit.RequestsPerMinute <= 0 {
			errors = append(errors, "RATE_LIMIT_REQUESTS_PER_MINUTE must be positive when rate limiting is enabled")
		}

		if c.RateLimit.BurstSize <= 0 {
			errors = append(errors, "RATE_LIMIT_BURST_SIZE must be positive when rate limiting is enabled")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// GetRedisURL returns the Redis connection URL
func (c *Config) GetRedisURL() string {
	if c.Redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d",
			c.Redis.Password,
			c.Redis.Host,
			c.Redis.Port,
			c.Redis.Database,
		)
	}
	return fmt.Sprintf("redis://%s:%d/%d",
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.Database,
	)
}

// Utility functions for environment variable parsing

// GetEnv returns environment variable or default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetIntEnv returns integer environment variable or default value
func GetIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// GetBoolEnv returns boolean environment variable or default value
func GetBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// GetDurationEnv returns duration environment variable or default value
func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// Keep lowercase versions for backward compatibility
func getEnv(key, defaultValue string) string {
	return GetEnv(key, defaultValue)
}

func getIntEnv(key string, defaultValue int) int {
	return GetIntEnv(key, defaultValue)
}

func getBoolEnv(key string, defaultValue bool) bool {
	return GetBoolEnv(key, defaultValue)
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	return GetDurationEnv(key, defaultValue)
}

// LoadFromFile loads configuration from a file (JSON, YAML, etc.)
func LoadFromFile(filename string) (*Config, error) {
	// This would implement file-based configuration loading
	// For now, we'll return an error indicating it's not implemented
	return nil, fmt.Errorf("file-based configuration loading not implemented")
}

// Save saves configuration to environment file
func (c *Config) Save(filename string) error {
	// This would implement saving configuration to a file
	// For now, we'll return an error indicating it's not implemented
	return fmt.Errorf("configuration saving not implemented")
}

// MustLoad loads configuration and panics on error
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	return config
}

// LoadWithServicePort loads configuration and overrides the port with service-specific environment variable
func LoadWithServicePort(serviceName string, defaultPort int) (*Config, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}

	// Try service-specific port environment variable first
	envVar := fmt.Sprintf("%s_SERVICE_PORT", strings.ToUpper(serviceName))
	if servicePort := getIntEnv(envVar, 0); servicePort != 0 {
		config.Server.Port = servicePort
	} else if config.Server.Port == 8080 && defaultPort != 8080 {
		// If still using default 8080 and we have a different default, use it
		config.Server.Port = defaultPort
	}

	return config, nil
}

// MustLoadWithServicePort loads configuration with service port and panics on error
func MustLoadWithServicePort(serviceName string, defaultPort int) *Config {
	config, err := LoadWithServicePort(serviceName, defaultPort)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration for %s: %v", serviceName, err))
	}
	return config
}