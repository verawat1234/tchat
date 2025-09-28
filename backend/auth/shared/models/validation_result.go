package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ValidationStatus represents the overall status of a validation
type ValidationStatus string

const (
	ValidationStatusPending    ValidationStatus = "pending"
	ValidationStatusRunning    ValidationStatus = "running"
	ValidationStatusSuccess    ValidationStatus = "success"
	ValidationStatusFailed     ValidationStatus = "failed"
	ValidationStatusTimeout    ValidationStatus = "timeout"
	ValidationStatusCancelled  ValidationStatus = "cancelled"
	ValidationStatusError      ValidationStatus = "error"
)

// IsValid checks if the validation status is valid
func (vs ValidationStatus) IsValid() bool {
	switch vs {
	case ValidationStatusPending, ValidationStatusRunning, ValidationStatusSuccess,
		 ValidationStatusFailed, ValidationStatusTimeout, ValidationStatusCancelled,
		 ValidationStatusError:
		return true
	default:
		return false
	}
}

// IsTerminal checks if the validation status is terminal
func (vs ValidationStatus) IsTerminal() bool {
	switch vs {
	case ValidationStatusSuccess, ValidationStatusFailed, ValidationStatusTimeout,
		 ValidationStatusCancelled, ValidationStatusError:
		return true
	default:
		return false
	}
}

// FailureType represents the type of validation failure
type FailureType string

const (
	FailureTypeRequestMismatch     FailureType = "request_mismatch"
	FailureTypeResponseMismatch    FailureType = "response_mismatch"
	FailureTypeStatusCodeMismatch  FailureType = "status_code_mismatch"
	FailureTypeHeaderMismatch      FailureType = "header_mismatch"
	FailureTypeBodyMismatch        FailureType = "body_mismatch"
	FailureTypeProviderStateError  FailureType = "provider_state_error"
	FailureTypeNetworkError        FailureType = "network_error"
	FailureTypeTimeoutError        FailureType = "timeout_error"
	FailureTypeValidationError     FailureType = "validation_error"
	FailureTypeSetupError          FailureType = "setup_error"
	FailureTypeTeardownError       FailureType = "teardown_error"
	FailureTypeMatchingRuleError   FailureType = "matching_rule_error"
	FailureTypeSchemaValidationError FailureType = "schema_validation_error"
)

// IsValid checks if the failure type is valid
func (ft FailureType) IsValid() bool {
	switch ft {
	case FailureTypeRequestMismatch, FailureTypeResponseMismatch, FailureTypeStatusCodeMismatch,
		 FailureTypeHeaderMismatch, FailureTypeBodyMismatch, FailureTypeProviderStateError,
		 FailureTypeNetworkError, FailureTypeTimeoutError, FailureTypeValidationError,
		 FailureTypeSetupError, FailureTypeTeardownError, FailureTypeMatchingRuleError,
		 FailureTypeSchemaValidationError:
		return true
	default:
		return false
	}
}

// FailureSeverity represents the severity level of a validation failure
type FailureSeverity string

const (
	FailureSeverityLow      FailureSeverity = "low"
	FailureSeverityMedium   FailureSeverity = "medium"
	FailureSeverityHigh     FailureSeverity = "high"
	FailureSeverityCritical FailureSeverity = "critical"
)

// IsValid checks if the failure severity is valid
func (fs FailureSeverity) IsValid() bool {
	switch fs {
	case FailureSeverityLow, FailureSeverityMedium, FailureSeverityHigh, FailureSeverityCritical:
		return true
	default:
		return false
	}
}

// ValidationFailure represents a specific validation failure
type ValidationFailure struct {
	ID              uuid.UUID       `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	InteractionID   uuid.UUID       `json:"interaction_id" gorm:"column:interaction_id;type:uuid"`
	Type            FailureType     `json:"type" gorm:"column:type;size:50;not null"`
	Severity        FailureSeverity `json:"severity" gorm:"column:severity;size:20;not null"`
	Message         string          `json:"message" gorm:"column:message;size:1000;not null"`
	DetailedMessage string          `json:"detailed_message,omitempty" gorm:"column:detailed_message;type:text"`
	ExpectedValue   interface{}     `json:"expected_value,omitempty" gorm:"column:expected_value;type:jsonb"`
	ActualValue     interface{}     `json:"actual_value,omitempty" gorm:"column:actual_value;type:jsonb"`
	Path            string          `json:"path,omitempty" gorm:"column:path;size:500"`
	StackTrace      string          `json:"stack_trace,omitempty" gorm:"column:stack_trace;type:text"`
	Context         map[string]interface{} `json:"context,omitempty" gorm:"column:context;type:jsonb"`
	OccurredAt      time.Time       `json:"occurred_at" gorm:"column:occurred_at;not null"`
}

// PerformanceMetrics represents performance measurements during validation
type PerformanceMetrics struct {
	ValidationDuration   time.Duration `json:"validation_duration" gorm:"column:validation_duration;not null"`
	SetupDuration        time.Duration `json:"setup_duration" gorm:"column:setup_duration"`
	TeardownDuration     time.Duration `json:"teardown_duration" gorm:"column:teardown_duration"`
	TotalInteractions    int           `json:"total_interactions" gorm:"column:total_interactions;not null"`
	SuccessfulInteractions int         `json:"successful_interactions" gorm:"column:successful_interactions;not null"`
	FailedInteractions   int           `json:"failed_interactions" gorm:"column:failed_interactions;not null"`
	AverageResponseTime  time.Duration `json:"average_response_time" gorm:"column:average_response_time"`
	MinResponseTime      time.Duration `json:"min_response_time" gorm:"column:min_response_time"`
	MaxResponseTime      time.Duration `json:"max_response_time" gorm:"column:max_response_time"`
	TotalDataTransferred int64         `json:"total_data_transferred" gorm:"column:total_data_transferred"` // bytes
	MemoryUsage          int64         `json:"memory_usage" gorm:"column:memory_usage"` // bytes
	CPUUsage             float64       `json:"cpu_usage" gorm:"column:cpu_usage"`       // percentage
	ConcurrentRequests   int           `json:"concurrent_requests" gorm:"column:concurrent_requests"`
	RequestsPerSecond    float64       `json:"requests_per_second" gorm:"column:requests_per_second"`
}

// EnvironmentInfo represents information about the test environment
type EnvironmentInfo struct {
	Name            string                 `json:"name" gorm:"column:name;size:50;not null"`
	Type            string                 `json:"type" gorm:"column:type;size:20;not null"` // development, staging, production, testing
	Region          string                 `json:"region" gorm:"column:region;size:20"`
	Infrastructure  string                 `json:"infrastructure" gorm:"column:infrastructure;size:100"`
	ProviderBaseURL string                 `json:"provider_base_url" gorm:"column:provider_base_url;size:500"`
	DatabaseVersion string                 `json:"database_version,omitempty" gorm:"column:database_version;size:50"`
	ServiceVersions map[string]string      `json:"service_versions,omitempty" gorm:"column:service_versions;type:jsonb"`
	ConfigOverrides map[string]interface{} `json:"config_overrides,omitempty" gorm:"column:config_overrides;type:jsonb"`
	FeatureFlags    map[string]bool        `json:"feature_flags,omitempty" gorm:"column:feature_flags;type:jsonb"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

// VersionInfo represents version information for consumer or provider
type VersionInfo struct {
	Version     string                 `json:"version" gorm:"column:version;size:50;not null"`
	Branch      string                 `json:"branch,omitempty" gorm:"column:branch;size:100"`
	CommitHash  string                 `json:"commit_hash,omitempty" gorm:"column:commit_hash;size:40"`
	BuildNumber string                 `json:"build_number,omitempty" gorm:"column:build_number;size:50"`
	BuildDate   *time.Time             `json:"build_date,omitempty" gorm:"column:build_date"`
	Tags        []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

// ValidationConfiguration represents configuration used during validation
type ValidationConfiguration struct {
	Timeout            time.Duration          `json:"timeout" gorm:"column:timeout"`
	MaxRetries         int                    `json:"max_retries" gorm:"column:max_retries"`
	RetryDelay         time.Duration          `json:"retry_delay" gorm:"column:retry_delay"`
	ConcurrencyLevel   int                    `json:"concurrency_level" gorm:"column:concurrency_level"`
	FailFast           bool                   `json:"fail_fast" gorm:"column:fail_fast"`
	ValidateOptional   bool                   `json:"validate_optional" gorm:"column:validate_optional"`
	StrictMatching     bool                   `json:"strict_matching" gorm:"column:strict_matching"`
	LogLevel           string                 `json:"log_level" gorm:"column:log_level;size:20"`
	CustomHeaders      map[string]string      `json:"custom_headers,omitempty" gorm:"column:custom_headers;type:jsonb"`
	CustomSettings     map[string]interface{} `json:"custom_settings,omitempty" gorm:"column:custom_settings;type:jsonb"`
}

// ValidationResult represents the outcome of contract verification between consumer and provider
type ValidationResult struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContractID  uuid.UUID        `json:"contract_id" gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`

	// Validation identification and status
	Status      ValidationStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	Success     bool             `json:"success" gorm:"column:success;not null;index"`

	// Version information
	ProviderVersion VersionInfo `json:"provider_version" gorm:"embedded;embeddedPrefix:provider_"`
	ConsumerVersion VersionInfo `json:"consumer_version" gorm:"embedded;embeddedPrefix:consumer_"`

	// Validation details
	Failures       []ValidationFailure `json:"failures" gorm:"column:failures;type:jsonb"`
	FailureCount   int                 `json:"failure_count" gorm:"column:failure_count;default:0"`
	WarningCount   int                 `json:"warning_count" gorm:"column:warning_count;default:0"`

	// Performance and metrics
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics" gorm:"embedded;embeddedPrefix:perf_"`

	// Environment and configuration
	Environment         EnvironmentInfo          `json:"environment" gorm:"embedded;embeddedPrefix:env_"`
	ValidationConfig    ValidationConfiguration `json:"validation_config" gorm:"embedded;embeddedPrefix:config_"`

	// Execution details
	ValidationTimestamp time.Time            `json:"validation_timestamp" gorm:"column:validation_timestamp;not null;index"`
	StartedAt           time.Time            `json:"started_at" gorm:"column:started_at;not null"`
	CompletedAt         *time.Time           `json:"completed_at,omitempty" gorm:"column:completed_at;index"`
	ExecutionDuration   time.Duration        `json:"execution_duration" gorm:"column:execution_duration"`

	// Validation context
	TriggeredBy         string               `json:"triggered_by" gorm:"column:triggered_by;size:100"`
	TriggerReason       string               `json:"trigger_reason" gorm:"column:trigger_reason;size:200"`
	ValidationTool      string               `json:"validation_tool" gorm:"column:validation_tool;size:50"`
	ValidationToolVersion string             `json:"validation_tool_version" gorm:"column:validation_tool_version;size:20"`

	// CI/CD integration
	BuildID             string               `json:"build_id,omitempty" gorm:"column:build_id;size:100"`
	JobID               string               `json:"job_id,omitempty" gorm:"column:job_id;size:100"`
	PipelineID          string               `json:"pipeline_id,omitempty" gorm:"column:pipeline_id;size:100"`
	RepositoryURL       string               `json:"repository_url,omitempty" gorm:"column:repository_url;size:500"`

	// Result artifacts and logs
	LogsURL             string               `json:"logs_url,omitempty" gorm:"column:logs_url;size:500"`
	ReportURL           string               `json:"report_url,omitempty" gorm:"column:report_url;size:500"`
	ArtifactsURL        string               `json:"artifacts_url,omitempty" gorm:"column:artifacts_url;size:500"`

	// Regional and compliance
	DataRegion          string               `json:"data_region" gorm:"column:data_region;size:20;default:'sea-central'"`
	ComplianceFlags     map[string]bool      `json:"compliance_flags,omitempty" gorm:"column:compliance_flags;type:jsonb"`

	// Metadata and tagging
	Tags                []string             `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	Metadata            map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`

	// Timestamps
	CreatedAt           time.Time            `json:"created_at" gorm:"column:created_at;not null;index"`
	UpdatedAt           time.Time            `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt           gorm.DeletedAt       `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships (not stored in database)
	ContractSpecification *ContractSpecification `json:"contract_specification,omitempty" gorm:"-"`
}

// TableName returns the table name for the ValidationResult model
func (ValidationResult) TableName() string {
	return "validation_results"
}

// BeforeCreate sets up the validation result before creation
func (vr *ValidationResult) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if vr.ID == uuid.Nil {
		vr.ID = uuid.New()
	}

	// Set validation timestamp if not provided
	if vr.ValidationTimestamp.IsZero() {
		vr.ValidationTimestamp = time.Now()
	}

	// Set started timestamp
	if vr.StartedAt.IsZero() {
		vr.StartedAt = time.Now()
	}

	// Set data region
	if vr.DataRegion == "" {
		vr.DataRegion = "sea-central" // Default for Southeast Asian deployment
	}

	// Initialize counts
	vr.FailureCount = len(vr.Failures)
	vr.calculateWarningCount()

	// Set default validation tool if not provided
	if vr.ValidationTool == "" {
		vr.ValidationTool = "pact-foundation"
	}

	// Set default configuration values
	if vr.ValidationConfig.Timeout == 0 {
		vr.ValidationConfig.Timeout = 30 * time.Second
	}
	if vr.ValidationConfig.MaxRetries == 0 {
		vr.ValidationConfig.MaxRetries = 3
	}
	if vr.ValidationConfig.RetryDelay == 0 {
		vr.ValidationConfig.RetryDelay = 1 * time.Second
	}
	if vr.ValidationConfig.ConcurrencyLevel == 0 {
		vr.ValidationConfig.ConcurrencyLevel = 1
	}

	// Initialize maps if nil
	if vr.Environment.ServiceVersions == nil {
		vr.Environment.ServiceVersions = make(map[string]string)
	}
	if vr.Environment.FeatureFlags == nil {
		vr.Environment.FeatureFlags = make(map[string]bool)
	}
	if vr.ComplianceFlags == nil {
		vr.ComplianceFlags = make(map[string]bool)
	}

	// Validate the validation result
	if err := vr.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the validation result before updating
func (vr *ValidationResult) BeforeUpdate(tx *gorm.DB) error {
	// Update counts
	vr.FailureCount = len(vr.Failures)
	vr.calculateWarningCount()

	// Update success flag
	vr.Success = vr.FailureCount == 0 && vr.Status == ValidationStatusSuccess

	// Calculate execution duration if completed
	if vr.CompletedAt != nil && !vr.StartedAt.IsZero() {
		vr.ExecutionDuration = vr.CompletedAt.Sub(vr.StartedAt)
	}

	return vr.Validate()
}

// Validate validates the validation result data
func (vr *ValidationResult) Validate() error {
	// Validate UUIDs
	if vr.ID == uuid.Nil {
		return fmt.Errorf("validation result ID cannot be nil")
	}
	if vr.ContractID == uuid.Nil {
		return fmt.Errorf("contract ID cannot be nil")
	}

	// Validate status
	if !vr.Status.IsValid() {
		return fmt.Errorf("invalid validation status: %s", vr.Status)
	}

	// Validate timestamp bounds
	now := time.Now()
	if vr.ValidationTimestamp.After(now) {
		return fmt.Errorf("validation timestamp cannot be in the future")
	}

	// Validate that timestamp is not too old (1 year)
	oneYearAgo := now.AddDate(-1, 0, 0)
	if vr.ValidationTimestamp.Before(oneYearAgo) {
		return fmt.Errorf("validation timestamp is too old (more than 1 year)")
	}

	// Validate version information
	if err := vr.validateVersionInfo(); err != nil {
		return err
	}

	// Validate environment
	if err := vr.validateEnvironment(); err != nil {
		return err
	}

	// Validate configuration
	if err := vr.validateConfiguration(); err != nil {
		return err
	}

	// Validate failures
	if err := vr.validateFailures(); err != nil {
		return err
	}

	// Validate performance metrics
	if err := vr.validatePerformanceMetrics(); err != nil {
		return err
	}

	return nil
}

// validateVersionInfo validates provider and consumer version information
func (vr *ValidationResult) validateVersionInfo() error {
	// Validate provider version
	if vr.ProviderVersion.Version == "" {
		return fmt.Errorf("provider version is required")
	}

	// Validate consumer version
	if vr.ConsumerVersion.Version == "" {
		return fmt.Errorf("consumer version is required")
	}

	// Validate build dates if provided
	if vr.ProviderVersion.BuildDate != nil && vr.ProviderVersion.BuildDate.After(time.Now()) {
		return fmt.Errorf("provider build date cannot be in the future")
	}
	if vr.ConsumerVersion.BuildDate != nil && vr.ConsumerVersion.BuildDate.After(time.Now()) {
		return fmt.Errorf("consumer build date cannot be in the future")
	}

	return nil
}

// validateEnvironment validates environment information
func (vr *ValidationResult) validateEnvironment() error {
	// Validate environment name
	if vr.Environment.Name == "" {
		return fmt.Errorf("environment name is required")
	}

	// Validate environment type
	if !IsValidEnvironment(vr.Environment.Type) {
		return fmt.Errorf("invalid environment type: %s", vr.Environment.Type)
	}

	// Validate provider base URL if provided
	if vr.Environment.ProviderBaseURL != "" {
		if !strings.HasPrefix(vr.Environment.ProviderBaseURL, "http://") &&
		   !strings.HasPrefix(vr.Environment.ProviderBaseURL, "https://") {
			return fmt.Errorf("provider base URL must start with http:// or https://")
		}
	}

	return nil
}

// validateConfiguration validates validation configuration
func (vr *ValidationResult) validateConfiguration() error {
	// Validate timeout
	if vr.ValidationConfig.Timeout <= 0 {
		return fmt.Errorf("validation timeout must be positive")
	}
	if vr.ValidationConfig.Timeout > 10*time.Minute {
		return fmt.Errorf("validation timeout cannot exceed 10 minutes")
	}

	// Validate retry settings
	if vr.ValidationConfig.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	if vr.ValidationConfig.MaxRetries > 10 {
		return fmt.Errorf("max retries cannot exceed 10")
	}

	if vr.ValidationConfig.RetryDelay < 0 {
		return fmt.Errorf("retry delay cannot be negative")
	}

	// Validate concurrency level
	if vr.ValidationConfig.ConcurrencyLevel <= 0 {
		return fmt.Errorf("concurrency level must be positive")
	}
	if vr.ValidationConfig.ConcurrencyLevel > 100 {
		return fmt.Errorf("concurrency level cannot exceed 100")
	}

	return nil
}

// validateFailures validates validation failures
func (vr *ValidationResult) validateFailures() error {
	for i, failure := range vr.Failures {
		if failure.ID == uuid.Nil {
			return fmt.Errorf("failure ID cannot be nil at index %d", i)
		}

		if !failure.Type.IsValid() {
			return fmt.Errorf("invalid failure type at index %d: %s", i, failure.Type)
		}

		if !failure.Severity.IsValid() {
			return fmt.Errorf("invalid failure severity at index %d: %s", i, failure.Severity)
		}

		if failure.Message == "" {
			return fmt.Errorf("failure message is required at index %d", i)
		}

		if failure.OccurredAt.IsZero() {
			return fmt.Errorf("failure occurred_at timestamp is required at index %d", i)
		}
	}

	return nil
}

// validatePerformanceMetrics validates performance metrics
func (vr *ValidationResult) validatePerformanceMetrics() error {
	metrics := &vr.PerformanceMetrics

	// Validate durations are non-negative
	if metrics.ValidationDuration < 0 {
		return fmt.Errorf("validation duration cannot be negative")
	}
	if metrics.SetupDuration < 0 {
		return fmt.Errorf("setup duration cannot be negative")
	}
	if metrics.TeardownDuration < 0 {
		return fmt.Errorf("teardown duration cannot be negative")
	}

	// Validate interaction counts
	if metrics.TotalInteractions < 0 {
		return fmt.Errorf("total interactions cannot be negative")
	}
	if metrics.SuccessfulInteractions < 0 {
		return fmt.Errorf("successful interactions cannot be negative")
	}
	if metrics.FailedInteractions < 0 {
		return fmt.Errorf("failed interactions cannot be negative")
	}

	// Validate interaction count consistency
	if metrics.SuccessfulInteractions+metrics.FailedInteractions > metrics.TotalInteractions {
		return fmt.Errorf("successful + failed interactions cannot exceed total interactions")
	}

	// Validate response times
	if metrics.AverageResponseTime < 0 {
		return fmt.Errorf("average response time cannot be negative")
	}
	if metrics.MinResponseTime < 0 {
		return fmt.Errorf("min response time cannot be negative")
	}
	if metrics.MaxResponseTime < 0 {
		return fmt.Errorf("max response time cannot be negative")
	}

	// Validate response time consistency
	if metrics.MinResponseTime > 0 && metrics.MaxResponseTime > 0 && metrics.MinResponseTime > metrics.MaxResponseTime {
		return fmt.Errorf("min response time cannot be greater than max response time")
	}

	// Validate resource usage
	if metrics.TotalDataTransferred < 0 {
		return fmt.Errorf("total data transferred cannot be negative")
	}
	if metrics.MemoryUsage < 0 {
		return fmt.Errorf("memory usage cannot be negative")
	}
	if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
		return fmt.Errorf("CPU usage must be between 0 and 100 percent")
	}

	// Validate concurrency and throughput
	if metrics.ConcurrentRequests < 0 {
		return fmt.Errorf("concurrent requests cannot be negative")
	}
	if metrics.RequestsPerSecond < 0 {
		return fmt.Errorf("requests per second cannot be negative")
	}

	return nil
}

// calculateWarningCount calculates the number of warnings from failures
func (vr *ValidationResult) calculateWarningCount() {
	warningCount := 0
	for _, failure := range vr.Failures {
		if failure.Severity == FailureSeverityLow || failure.Severity == FailureSeverityMedium {
			warningCount++
		}
	}
	vr.WarningCount = warningCount
}

// UpdateStatus updates the validation status and handles completion logic
func (vr *ValidationResult) UpdateStatus(newStatus ValidationStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid validation status: %s", newStatus)
	}

	oldStatus := vr.Status
	vr.Status = newStatus

	// Handle status-specific logic
	switch newStatus {
	case ValidationStatusRunning:
		if vr.StartedAt.IsZero() {
			vr.StartedAt = time.Now()
		}

	case ValidationStatusSuccess, ValidationStatusFailed, ValidationStatusTimeout,
		 ValidationStatusCancelled, ValidationStatusError:
		// Terminal states
		if vr.CompletedAt == nil {
			now := time.Now()
			vr.CompletedAt = &now
		}

		if !vr.StartedAt.IsZero() && vr.CompletedAt != nil {
			vr.ExecutionDuration = vr.CompletedAt.Sub(vr.StartedAt)
		}

		// Update success flag
		vr.Success = newStatus == ValidationStatusSuccess && vr.FailureCount == 0
	}

	// Update performance metrics validation duration
	if !vr.StartedAt.IsZero() && vr.CompletedAt != nil {
		vr.PerformanceMetrics.ValidationDuration = vr.CompletedAt.Sub(vr.StartedAt)
	}

	vr.UpdatedAt = time.Now()

	// Log status change if needed
	if oldStatus != newStatus {
		vr.addStatusChangeToMetadata(oldStatus, newStatus)
	}

	return nil
}

// AddFailure adds a validation failure to the result
func (vr *ValidationResult) AddFailure(failure ValidationFailure) error {
	// Generate ID if not set
	if failure.ID == uuid.Nil {
		failure.ID = uuid.New()
	}

	// Set occurred timestamp if not set
	if failure.OccurredAt.IsZero() {
		failure.OccurredAt = time.Now()
	}

	// Validate the failure
	if !failure.Type.IsValid() {
		return fmt.Errorf("invalid failure type: %s", failure.Type)
	}
	if !failure.Severity.IsValid() {
		return fmt.Errorf("invalid failure severity: %s", failure.Severity)
	}
	if failure.Message == "" {
		return fmt.Errorf("failure message is required")
	}

	vr.Failures = append(vr.Failures, failure)
	vr.FailureCount = len(vr.Failures)
	vr.calculateWarningCount()
	vr.Success = false // Any failure makes the validation unsuccessful
	vr.UpdatedAt = time.Now()

	return nil
}

// GetCriticalFailures returns failures with critical severity
func (vr *ValidationResult) GetCriticalFailures() []ValidationFailure {
	var critical []ValidationFailure
	for _, failure := range vr.Failures {
		if failure.Severity == FailureSeverityCritical {
			critical = append(critical, failure)
		}
	}
	return critical
}

// GetFailuresByType returns failures of a specific type
func (vr *ValidationResult) GetFailuresByType(failureType FailureType) []ValidationFailure {
	var filtered []ValidationFailure
	for _, failure := range vr.Failures {
		if failure.Type == failureType {
			filtered = append(filtered, failure)
		}
	}
	return filtered
}

// GetSuccessRate calculates the interaction success rate
func (vr *ValidationResult) GetSuccessRate() float64 {
	if vr.PerformanceMetrics.TotalInteractions == 0 {
		return 0.0
	}
	return float64(vr.PerformanceMetrics.SuccessfulInteractions) / float64(vr.PerformanceMetrics.TotalInteractions)
}

// IsCompleted checks if the validation is completed (terminal state)
func (vr *ValidationResult) IsCompleted() bool {
	return vr.Status.IsTerminal()
}

// IsSuccessful checks if the validation was successful
func (vr *ValidationResult) IsSuccessful() bool {
	return vr.Success && vr.Status == ValidationStatusSuccess
}

// GetDuration returns the validation duration
func (vr *ValidationResult) GetDuration() time.Duration {
	if vr.CompletedAt != nil && !vr.StartedAt.IsZero() {
		return vr.CompletedAt.Sub(vr.StartedAt)
	}
	if !vr.StartedAt.IsZero() {
		return time.Since(vr.StartedAt)
	}
	return 0
}

// addStatusChangeToMetadata adds status change information to metadata
func (vr *ValidationResult) addStatusChangeToMetadata(oldStatus, newStatus ValidationStatus) {
	if vr.Metadata == nil {
		vr.Metadata = make(map[string]interface{})
	}

	statusChanges, exists := vr.Metadata["status_changes"]
	if !exists {
		statusChanges = make([]map[string]interface{}, 0)
	}

	changes := statusChanges.([]map[string]interface{})
	changes = append(changes, map[string]interface{}{
		"from":       oldStatus,
		"to":         newStatus,
		"timestamp":  time.Now(),
	})

	vr.Metadata["status_changes"] = changes
}

// MarshalJSON customizes JSON serialization
func (vr *ValidationResult) MarshalJSON() ([]byte, error) {
	type Alias ValidationResult
	return json.Marshal(&struct {
		*Alias
		IsCompleted        bool     `json:"is_completed"`
		IsSuccessful       bool     `json:"is_successful"`
		SuccessRate        float64  `json:"success_rate"`
		Duration           string   `json:"duration"`
		DurationMs         int64    `json:"duration_ms"`
		CriticalFailures   int      `json:"critical_failures"`
		Age                string   `json:"age"`
		TimeSinceUpdate    string   `json:"time_since_update"`
		HasCriticalFailures bool    `json:"has_critical_failures"`
	}{
		Alias:               (*Alias)(vr),
		IsCompleted:         vr.IsCompleted(),
		IsSuccessful:        vr.IsSuccessful(),
		SuccessRate:         vr.GetSuccessRate(),
		Duration:            vr.GetDuration().String(),
		DurationMs:          vr.GetDuration().Milliseconds(),
		CriticalFailures:    len(vr.GetCriticalFailures()),
		Age:                 time.Since(vr.CreatedAt).String(),
		TimeSinceUpdate:     time.Since(vr.UpdatedAt).String(),
		HasCriticalFailures: len(vr.GetCriticalFailures()) > 0,
	})
}

// Helper functions

// GetSupportedValidationStatuses returns all supported validation statuses
func GetSupportedValidationStatuses() []ValidationStatus {
	return []ValidationStatus{
		ValidationStatusPending,
		ValidationStatusRunning,
		ValidationStatusSuccess,
		ValidationStatusFailed,
		ValidationStatusTimeout,
		ValidationStatusCancelled,
		ValidationStatusError,
	}
}

// GetSupportedFailureTypes returns all supported failure types
func GetSupportedFailureTypes() []FailureType {
	return []FailureType{
		FailureTypeRequestMismatch,
		FailureTypeResponseMismatch,
		FailureTypeStatusCodeMismatch,
		FailureTypeHeaderMismatch,
		FailureTypeBodyMismatch,
		FailureTypeProviderStateError,
		FailureTypeNetworkError,
		FailureTypeTimeoutError,
		FailureTypeValidationError,
		FailureTypeSetupError,
		FailureTypeTeardownError,
		FailureTypeMatchingRuleError,
		FailureTypeSchemaValidationError,
	}
}

// GetSupportedFailureSeverities returns all supported failure severities
func GetSupportedFailureSeverities() []FailureSeverity {
	return []FailureSeverity{
		FailureSeverityLow,
		FailureSeverityMedium,
		FailureSeverityHigh,
		FailureSeverityCritical,
	}
}