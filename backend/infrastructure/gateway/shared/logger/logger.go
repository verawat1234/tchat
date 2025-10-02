package logger

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

// RegionalComplianceConfig holds regional logging compliance settings
type RegionalComplianceConfig struct {
	Enabled         bool
	DataRetention   map[string]time.Duration // Country-specific retention periods
	PIIRedaction    bool                     // Personally Identifiable Information redaction
	AuditRequired   bool                     // Audit trail requirements
	EncryptionKey   string                   // Key for sensitive data encryption
	LocalStorage    bool                     // Local storage requirements for compliance
}

// TchatLogger represents the structured logger for Southeast Asian compliance
type TchatLogger struct {
	*logrus.Logger
	ServiceName     string
	ServiceVersion  string
	Environment     string
	Region          string
	ComplianceConfig *RegionalComplianceConfig
	RequestIDField  string
	UserIDField     string
	CountryField    string
}

// LoggerConfig holds configuration for the Tchat logger
type LoggerConfig struct {
	Level           LogLevel
	Format          string // json, text
	ServiceName     string
	ServiceVersion  string
	Environment     string
	Region          string
	OutputPath      string
	ComplianceConfig *RegionalComplianceConfig
}

// DefaultLoggerConfig returns a default logger configuration for development
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          InfoLevel,
		Format:         "json",
		ServiceName:    "tchat-service",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		Region:         "sea-central",
		OutputPath:     "stdout",
		ComplianceConfig: &RegionalComplianceConfig{
			Enabled: true,
			DataRetention: map[string]time.Duration{
				"TH": 7 * 365 * 24 * time.Hour,  // Thailand: 7 years
				"SG": 5 * 365 * 24 * time.Hour,  // Singapore: 5 years
				"ID": 5 * 365 * 24 * time.Hour,  // Indonesia: 5 years
				"MY": 7 * 365 * 24 * time.Hour,  // Malaysia: 7 years
				"PH": 5 * 365 * 24 * time.Hour,  // Philippines: 5 years
				"VN": 5 * 365 * 24 * time.Hour,  // Vietnam: 5 years
			},
			PIIRedaction:  true,
			AuditRequired: true,
			LocalStorage:  true,
		},
	}
}

// NewTchatLogger creates a new structured logger with Southeast Asian compliance
func NewTchatLogger(config *LoggerConfig) *TchatLogger {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	logger := logrus.New()

	// Set log level
	switch config.Level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case PanicLevel:
		logger.SetLevel(logrus.PanicLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set output format
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
				logrus.FieldKeyFile:  "file",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	}

	// Set output destination
	if config.OutputPath == "stdout" || config.OutputPath == "" {
		logger.SetOutput(os.Stdout)
	} else {
		// In production, this would be configured to write to files or log aggregation services
		logger.SetOutput(os.Stdout)
	}

	// Report caller information for debugging
	logger.SetReportCaller(true)

	tchatLogger := &TchatLogger{
		Logger:          logger,
		ServiceName:     config.ServiceName,
		ServiceVersion:  config.ServiceVersion,
		Environment:     config.Environment,
		Region:          config.Region,
		ComplianceConfig: config.ComplianceConfig,
		RequestIDField:  "request_id",
		UserIDField:     "user_id",
		CountryField:    "country_code",
	}

	// Add hooks for compliance
	if config.ComplianceConfig != nil && config.ComplianceConfig.Enabled {
		tchatLogger.AddComplianceHooks()
	}

	return tchatLogger
}

// WithFields returns a new logger entry with additional fields
func (l *TchatLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	// Always include service metadata
	enrichedFields := logrus.Fields{
		"service":         l.ServiceName,
		"service_version": l.ServiceVersion,
		"environment":     l.Environment,
		"region":          l.Region,
	}

	// Add user-provided fields
	for k, v := range fields {
		// Apply PII redaction if enabled
		if l.ComplianceConfig != nil && l.ComplianceConfig.PIIRedaction {
			enrichedFields[k] = l.redactPII(k, v)
		} else {
			enrichedFields[k] = v
		}
	}

	return l.Logger.WithFields(enrichedFields)
}

// WithContext extracts relevant information from context and adds to log fields
func (l *TchatLogger) WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	// Extract request ID if available (from Gin context)
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID := ginCtx.GetHeader("X-Request-ID"); requestID != "" {
			fields[l.RequestIDField] = requestID
		}
		if userID := ginCtx.GetString("user_id"); userID != "" {
			fields[l.UserIDField] = userID
		}
		if countryCode := ginCtx.GetString("country_code"); countryCode != "" {
			fields[l.CountryField] = countryCode
		}
		if locale := ginCtx.GetString("locale"); locale != "" {
			fields["locale"] = locale
		}
		if sessionID := ginCtx.GetString("session_id"); sessionID != "" {
			fields["session_id"] = sessionID
		}
	}

	// Extract from standard context values
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields[l.RequestIDField] = requestID
	}
	if userID := ctx.Value("user_id"); userID != nil {
		fields[l.UserIDField] = userID
	}
	if countryCode := ctx.Value("country_code"); countryCode != nil {
		fields[l.CountryField] = countryCode
	}

	return l.WithFields(fields)
}

// WithRequest logs with HTTP request context information
func (l *TchatLogger) WithRequest(c *gin.Context) *logrus.Entry {
	fields := logrus.Fields{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"remote_addr": c.ClientIP(),
		"user_agent":  c.GetHeader("User-Agent"),
	}

	// Add request ID if available
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		fields[l.RequestIDField] = requestID
	}

	// Add user context if authenticated
	if userID := c.GetString("user_id"); userID != "" {
		fields[l.UserIDField] = userID
	}
	if countryCode := c.GetString("country_code"); countryCode != "" {
		fields[l.CountryField] = countryCode
	}
	if locale := c.GetString("locale"); locale != "" {
		fields["locale"] = locale
	}

	return l.WithFields(fields)
}

// WithError logs with error context and stack trace
func (l *TchatLogger) WithError(err error) *logrus.Entry {
	fields := logrus.Fields{
		"error": err.Error(),
	}

	// Add stack trace for debugging (only in development)
	if l.Environment == "development" {
		const depth = 32
		var pcs [depth]uintptr
		n := runtime.Callers(3, pcs[:])
		frames := runtime.CallersFrames(pcs[:n])

		var stack []string
		for {
			frame, more := frames.Next()
			stack = append(stack, frame.Function+"():"+frame.File)
			if !more {
				break
			}
		}
		fields["stack_trace"] = stack
	}

	return l.WithFields(fields)
}

// AddComplianceHooks adds hooks for Southeast Asian regulatory compliance
func (l *TchatLogger) AddComplianceHooks() {
	// Add audit trail hook
	if l.ComplianceConfig.AuditRequired {
		l.Logger.AddHook(&AuditHook{
			ServiceName: l.ServiceName,
			Region:      l.Region,
		})
	}

	// Add data retention hook
	l.Logger.AddHook(&DataRetentionHook{
		RetentionPolicies: l.ComplianceConfig.DataRetention,
	})
}

// redactPII redacts personally identifiable information from log fields
func (l *TchatLogger) redactPII(key string, value interface{}) interface{} {
	// List of fields that may contain PII
	piiFields := map[string]bool{
		"email":        true,
		"phone":        true,
		"phone_number": true,
		"address":      true,
		"name":         true,
		"full_name":    true,
		"first_name":   true,
		"last_name":    true,
		"password":     true,
		"ssn":          true,
		"national_id":  true,
		"passport":     true,
		"credit_card":  true,
		"bank_account": true,
	}

	keyLower := strings.ToLower(key)
	if piiFields[keyLower] {
		if str, ok := value.(string); ok && len(str) > 0 {
			// Redact but keep format for debugging
			if len(str) <= 3 {
				return "***"
			}
			return str[:1] + strings.Repeat("*", len(str)-2) + str[len(str)-1:]
		}
		return "[REDACTED]"
	}

	return value
}

// AuditLog logs an audit trail entry for compliance
func (l *TchatLogger) AuditLog(action string, userID string, resource string, details map[string]interface{}) {
	fields := logrus.Fields{
		"audit_type":   "user_action",
		"action":       action,
		"user_id":      userID,
		"resource":     resource,
		"timestamp":    time.Now().UTC(),
		"service":      l.ServiceName,
		"environment":  l.Environment,
		"region":       l.Region,
	}

	// Add details if provided
	for k, v := range details {
		fields[k] = v
	}

	l.WithFields(fields).Info("Audit log entry")
}

// SecurityLog logs security-related events for monitoring
func (l *TchatLogger) SecurityLog(eventType string, userID string, severity string, details map[string]interface{}) {
	fields := logrus.Fields{
		"security_event": eventType,
		"user_id":        userID,
		"severity":       severity,
		"timestamp":      time.Now().UTC(),
		"service":        l.ServiceName,
		"environment":    l.Environment,
		"region":         l.Region,
	}

	// Add details if provided
	for k, v := range details {
		fields[k] = v
	}

	entry := l.WithFields(fields)
	switch severity {
	case "critical":
		entry.Error("Security event logged")
	case "high":
		entry.Warn("Security event logged")
	default:
		entry.Info("Security event logged")
	}
}

// PerformanceLog logs performance metrics for monitoring
func (l *TchatLogger) PerformanceLog(operation string, duration time.Duration, details map[string]interface{}) {
	fields := logrus.Fields{
		"performance_metric": true,
		"operation":         operation,
		"duration_ms":       duration.Milliseconds(),
		"duration_ns":       duration.Nanoseconds(),
		"timestamp":         time.Now().UTC(),
		"service":           l.ServiceName,
		"environment":       l.Environment,
		"region":            l.Region,
	}

	// Add details if provided
	for k, v := range details {
		fields[k] = v
	}

	l.WithFields(fields).Info("Performance metric logged")
}

// GetGlobalLogger returns a global logger instance for the application
var globalLogger *TchatLogger

// InitGlobalLogger initializes the global logger with the provided configuration
func InitGlobalLogger(config *LoggerConfig) {
	globalLogger = NewTchatLogger(config)
}

// GetLogger returns the global logger instance
func GetLogger() *TchatLogger {
	if globalLogger == nil {
		globalLogger = NewTchatLogger(DefaultLoggerConfig())
	}
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.Logger.WithFields(logrus.Fields{
		"service":         logger.ServiceName,
		"service_version": logger.ServiceVersion,
		"environment":     logger.Environment,
		"region":          logger.Region,
	})
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Debug(msg)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.Logger.WithFields(logrus.Fields{
		"service":         logger.ServiceName,
		"service_version": logger.ServiceVersion,
		"environment":     logger.Environment,
		"region":          logger.Region,
	})
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Info(msg)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.Logger.WithFields(logrus.Fields{
		"service":         logger.ServiceName,
		"service_version": logger.ServiceVersion,
		"environment":     logger.Environment,
		"region":          logger.Region,
	})
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Warn(msg)
}

// Error logs an error message using the global logger
func Error(msg string, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.Logger.WithFields(logrus.Fields{
		"service":         logger.ServiceName,
		"service_version": logger.ServiceVersion,
		"environment":     logger.Environment,
		"region":          logger.Region,
	})
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}