package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// AuditHook implements logrus.Hook for audit trail compliance
type AuditHook struct {
	ServiceName string
	Region      string
	Writer      *os.File
}

// Levels returns the levels this hook is interested in
func (hook *AuditHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

// Fire is called when a log entry is fired
func (hook *AuditHook) Fire(entry *logrus.Entry) error {
	// Only process audit-related entries
	if _, exists := entry.Data["audit_type"]; !exists {
		return nil
	}

	// Create audit log entry
	auditEntry := map[string]interface{}{
		"timestamp":    entry.Time.UTC().Format(time.RFC3339Nano),
		"level":        entry.Level.String(),
		"service":      hook.ServiceName,
		"region":       hook.Region,
		"message":      entry.Message,
		"data":         entry.Data,
	}

	// In production, this would write to secure audit storage
	// For development, we'll write to a local audit file
	if hook.Writer == nil {
		auditDir := "logs/audit"
		if err := os.MkdirAll(auditDir, 0750); err != nil {
			return fmt.Errorf("failed to create audit directory: %w", err)
		}

		filename := fmt.Sprintf("%s/audit_%s_%s.log",
			auditDir,
			hook.ServiceName,
			entry.Time.Format("2006-01-02"))

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			return fmt.Errorf("failed to open audit file: %w", err)
		}
		hook.Writer = file
	}

	// Write audit entry (in production, this would be encrypted)
	auditLine := fmt.Sprintf("%s [%s] %s: %s %v\n",
		auditEntry["timestamp"],
		strings.ToUpper(auditEntry["level"].(string)),
		hook.ServiceName,
		auditEntry["message"],
		auditEntry["data"])

	_, err := hook.Writer.WriteString(auditLine)
	if err != nil {
		return fmt.Errorf("failed to write audit entry: %w", err)
	}

	return hook.Writer.Sync()
}

// DataRetentionHook implements logrus.Hook for data retention compliance
type DataRetentionHook struct {
	RetentionPolicies map[string]time.Duration
	LastCleanup       time.Time
}

// Levels returns the levels this hook is interested in
func (hook *DataRetentionHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Fire is called when a log entry is fired
func (hook *DataRetentionHook) Fire(entry *logrus.Entry) error {
	// Run cleanup once per day
	if time.Since(hook.LastCleanup) < 24*time.Hour {
		return nil
	}

	// Cleanup old log files based on retention policies
	go func() {
		if err := hook.cleanupOldLogs(); err != nil {
			// Log cleanup error (but avoid infinite recursion)
			fmt.Printf("Log cleanup error: %v\n", err)
		}
		hook.LastCleanup = time.Now()
	}()

	return nil
}

// cleanupOldLogs removes log files older than retention policies allow
func (hook *DataRetentionHook) cleanupOldLogs() error {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return nil // No logs directory exists
	}

	return filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-log files
		if !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}

		// Determine retention period based on log type and country
		var retentionPeriod time.Duration = 365 * 24 * time.Hour // Default 1 year

		// Check if this is a country-specific log
		for countryCode, retention := range hook.RetentionPolicies {
			if strings.Contains(path, strings.ToLower(countryCode)) {
				retentionPeriod = retention
				break
			}
		}

		// Check if file is older than retention period
		if time.Since(info.ModTime()) > retentionPeriod {
			fmt.Printf("Removing old log file due to retention policy: %s\n", path)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove old log file %s: %w", path, err)
			}
		}

		return nil
	})
}

// SecurityHook implements logrus.Hook for security event monitoring
type SecurityHook struct {
	AlertThresholds map[string]int
	EventCounts     map[string]int
	LastReset       time.Time
}

// NewSecurityHook creates a new security monitoring hook
func NewSecurityHook() *SecurityHook {
	return &SecurityHook{
		AlertThresholds: map[string]int{
			"failed_login":      5,   // 5 failed logins in 1 hour
			"invalid_token":     10,  // 10 invalid tokens in 1 hour
			"rate_limit_hit":    20,  // 20 rate limit hits in 1 hour
			"security_violation": 1,  // Any security violation
		},
		EventCounts: make(map[string]int),
		LastReset:   time.Now(),
	}
}

// Levels returns the levels this hook is interested in
func (hook *SecurityHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

// Fire is called when a log entry is fired
func (hook *SecurityHook) Fire(entry *logrus.Entry) error {
	// Reset counters every hour
	if time.Since(hook.LastReset) > time.Hour {
		hook.EventCounts = make(map[string]int)
		hook.LastReset = time.Now()
	}

	// Check for security events
	if eventType, exists := entry.Data["security_event"]; exists {
		eventTypeStr := fmt.Sprintf("%v", eventType)
		hook.EventCounts[eventTypeStr]++

		// Check if threshold is exceeded
		if threshold, hasThreshold := hook.AlertThresholds[eventTypeStr]; hasThreshold {
			if hook.EventCounts[eventTypeStr] >= threshold {
				// In production, this would trigger security alerts
				fmt.Printf("SECURITY ALERT: %s threshold exceeded (%d events in 1 hour)\n",
					eventTypeStr, hook.EventCounts[eventTypeStr])

				// Reset counter to avoid spam
				hook.EventCounts[eventTypeStr] = 0
			}
		}
	}

	return nil
}

// PerformanceHook implements logrus.Hook for performance monitoring
type PerformanceHook struct {
	MetricsCollector MetricsCollector
}

// MetricsCollector interface for collecting performance metrics
type MetricsCollector interface {
	RecordDuration(operation string, duration time.Duration, labels map[string]string)
	RecordCounter(name string, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
}

// Levels returns the levels this hook is interested in
func (hook *PerformanceHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
	}
}

// Fire is called when a log entry is fired
func (hook *PerformanceHook) Fire(entry *logrus.Entry) error {
	// Only process performance metrics
	if _, exists := entry.Data["performance_metric"]; !exists {
		return nil
	}

	if hook.MetricsCollector == nil {
		return nil // No metrics collector configured
	}

	// Extract performance data
	operation, _ := entry.Data["operation"].(string)
	durationMs, _ := entry.Data["duration_ms"].(int64)

	if operation != "" && durationMs > 0 {
		duration := time.Duration(durationMs) * time.Millisecond

		// Create labels from log data
		labels := make(map[string]string)
		for key, value := range entry.Data {
			if key != "performance_metric" && key != "operation" && key != "duration_ms" && key != "duration_ns" {
				labels[key] = fmt.Sprintf("%v", value)
			}
		}

		// Record the metric
		hook.MetricsCollector.RecordDuration(operation, duration, labels)
	}

	return nil
}

// RegionalComplianceHook ensures logs comply with regional regulations
type RegionalComplianceHook struct {
	ComplianceRules map[string]ComplianceRule
}

// ComplianceRule defines compliance requirements for a region
type ComplianceRule struct {
	RequiredFields   []string          // Fields that must be present
	ForbiddenFields  []string          // Fields that must not be present
	DataLocalization bool              // Whether data must stay in region
	EncryptionLevel  string            // Required encryption level
	RetentionPeriod  time.Duration     // How long data must be retained
	Transformations  map[string]string // Field transformations required
}

// NewRegionalComplianceHook creates a new regional compliance hook
func NewRegionalComplianceHook() *RegionalComplianceHook {
	return &RegionalComplianceHook{
		ComplianceRules: map[string]ComplianceRule{
			"TH": { // Thailand - Personal Data Protection Act
				RequiredFields:   []string{"timestamp", "service", "user_consent"},
				ForbiddenFields:  []string{}, // No specific forbidden fields
				DataLocalization: true,
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  7 * 365 * 24 * time.Hour, // 7 years
			},
			"SG": { // Singapore - Personal Data Protection Act
				RequiredFields:   []string{"timestamp", "service", "purpose"},
				ForbiddenFields:  []string{},
				DataLocalization: false, // More flexible
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  5 * 365 * 24 * time.Hour, // 5 years
			},
			"ID": { // Indonesia - Personal Data Protection Law
				RequiredFields:   []string{"timestamp", "service", "legal_basis"},
				ForbiddenFields:  []string{},
				DataLocalization: true,
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  5 * 365 * 24 * time.Hour,
			},
			"MY": { // Malaysia - Personal Data Protection Act
				RequiredFields:   []string{"timestamp", "service", "consent_date"},
				ForbiddenFields:  []string{},
				DataLocalization: true,
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  7 * 365 * 24 * time.Hour,
			},
			"PH": { // Philippines - Data Privacy Act
				RequiredFields:   []string{"timestamp", "service", "data_category"},
				ForbiddenFields:  []string{},
				DataLocalization: true,
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  5 * 365 * 24 * time.Hour,
			},
			"VN": { // Vietnam - Personal Data Protection Decree
				RequiredFields:   []string{"timestamp", "service", "processing_purpose"},
				ForbiddenFields:  []string{},
				DataLocalization: true,
				EncryptionLevel:  "AES-256",
				RetentionPeriod:  5 * 365 * 24 * time.Hour,
			},
		},
	}
}

// Levels returns the levels this hook is interested in
func (hook *RegionalComplianceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log entry is fired
func (hook *RegionalComplianceHook) Fire(entry *logrus.Entry) error {
	// Get country code from log entry
	countryCode, exists := entry.Data["country_code"]
	if !exists {
		return nil // No country context
	}

	countryStr := fmt.Sprintf("%v", countryCode)
	rule, hasRule := hook.ComplianceRules[countryStr]
	if !hasRule {
		return nil // No specific rule for this country
	}

	// Check required fields
	for _, requiredField := range rule.RequiredFields {
		if _, hasField := entry.Data[requiredField]; !hasField {
			// In production, this would trigger compliance alerts
			fmt.Printf("COMPLIANCE WARNING: Missing required field '%s' for country %s\n",
				requiredField, countryStr)
		}
	}

	// Check forbidden fields
	for _, forbiddenField := range rule.ForbiddenFields {
		if _, hasField := entry.Data[forbiddenField]; hasField {
			// Remove forbidden field
			delete(entry.Data, forbiddenField)
			fmt.Printf("COMPLIANCE: Removed forbidden field '%s' for country %s\n",
				forbiddenField, countryStr)
		}
	}

	// Apply transformations
	for field, transformation := range rule.Transformations {
		if value, hasField := entry.Data[field]; hasField {
			// Apply transformation (simplified example)
			switch transformation {
			case "hash":
				entry.Data[field] = fmt.Sprintf("[HASHED:%x]", value)
			case "mask":
				entry.Data[field] = "[MASKED]"
			}
		}
	}

	return nil
}