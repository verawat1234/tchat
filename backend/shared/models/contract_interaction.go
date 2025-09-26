package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MatchingRuleType represents the type of matching rule
type MatchingRuleType string

const (
	MatchingRuleTypeEquality MatchingRuleType = "equality"
	MatchingRuleTypeRegex    MatchingRuleType = "regex"
	MatchingRuleTypeType     MatchingRuleType = "type"
	MatchingRuleTypeNumber   MatchingRuleType = "number"
	MatchingRuleTypeInteger  MatchingRuleType = "integer"
	MatchingRuleTypeDecimal  MatchingRuleType = "decimal"
	MatchingRuleTypeDate     MatchingRuleType = "date"
	MatchingRuleTypeTime     MatchingRuleType = "time"
	MatchingRuleTypeTimestamp MatchingRuleType = "timestamp"
	MatchingRuleTypeNull     MatchingRuleType = "null"
	MatchingRuleTypeBoolean  MatchingRuleType = "boolean"
	MatchingRuleTypeInclude  MatchingRuleType = "include"
	MatchingRuleTypeArrayContains MatchingRuleType = "arrayContains"
)

// IsValid checks if the matching rule type is valid
func (mrt MatchingRuleType) IsValid() bool {
	switch mrt {
	case MatchingRuleTypeEquality, MatchingRuleTypeRegex, MatchingRuleTypeType,
		 MatchingRuleTypeNumber, MatchingRuleTypeInteger, MatchingRuleTypeDecimal,
		 MatchingRuleTypeDate, MatchingRuleTypeTime, MatchingRuleTypeTimestamp,
		 MatchingRuleTypeNull, MatchingRuleTypeBoolean, MatchingRuleTypeInclude,
		 MatchingRuleTypeArrayContains:
		return true
	default:
		return false
	}
}

// ProviderStateType represents the type of provider state
type ProviderStateType string

const (
	ProviderStateTypeSetup    ProviderStateType = "setup"
	ProviderStateTypeTeardown ProviderStateType = "teardown"
	ProviderStateTypePersist  ProviderStateType = "persist"
	ProviderStateTypeCleanup  ProviderStateType = "cleanup"
)

// IsValid checks if the provider state type is valid
func (pst ProviderStateType) IsValid() bool {
	switch pst {
	case ProviderStateTypeSetup, ProviderStateTypeTeardown, ProviderStateTypePersist, ProviderStateTypeCleanup:
		return true
	default:
		return false
	}
}

// HTTPRequest represents the expected HTTP request structure
type HTTPRequest struct {
	Method      string                 `json:"method" gorm:"column:method;size:10;not null"`
	Path        string                 `json:"path" gorm:"column:path;size:500;not null"`
	Query       map[string]interface{} `json:"query,omitempty" gorm:"column:query;type:jsonb"`
	Headers     map[string]string      `json:"headers,omitempty" gorm:"column:headers;type:jsonb"`
	Body        json.RawMessage        `json:"body,omitempty" gorm:"column:body;type:jsonb"`
	ContentType string                 `json:"content_type,omitempty" gorm:"column:content_type;size:100"`
	Encoding    string                 `json:"encoding,omitempty" gorm:"column:encoding;size:50"`
}

// HTTPResponse represents the expected HTTP response structure
type HTTPResponse struct {
	Status      int                    `json:"status" gorm:"column:status;not null"`
	StatusText  string                 `json:"status_text,omitempty" gorm:"column:status_text;size:100"`
	Headers     map[string]string      `json:"headers,omitempty" gorm:"column:headers;type:jsonb"`
	Body        json.RawMessage        `json:"body,omitempty" gorm:"column:body;type:jsonb"`
	ContentType string                 `json:"content_type,omitempty" gorm:"column:content_type;size:100"`
	Encoding    string                 `json:"encoding,omitempty" gorm:"column:encoding;size:50"`
}

// MatchingRule represents a rule for flexible matching in contract verification
type MatchingRule struct {
	Type       MatchingRuleType `json:"type" gorm:"column:type;size:50;not null"`
	Regex      string           `json:"regex,omitempty" gorm:"column:regex;size:500"`
	Min        *int             `json:"min,omitempty" gorm:"column:min"`
	Max        *int             `json:"max,omitempty" gorm:"column:max"`
	Format     string           `json:"format,omitempty" gorm:"column:format;size:50"`
	Example    interface{}      `json:"example,omitempty" gorm:"column:example;type:jsonb"`
	Generator  string           `json:"generator,omitempty" gorm:"column:generator;size:200"`
}

// MatchingRuleSet represents a set of matching rules for different parts of the interaction
type MatchingRuleSet struct {
	Path         []MatchingRule            `json:"path,omitempty" gorm:"column:path;type:jsonb"`
	Query        map[string][]MatchingRule `json:"query,omitempty" gorm:"column:query;type:jsonb"`
	Header       map[string][]MatchingRule `json:"header,omitempty" gorm:"column:header;type:jsonb"`
	Body         map[string][]MatchingRule `json:"body,omitempty" gorm:"column:body;type:jsonb"`
	StatusCode   []MatchingRule            `json:"status_code,omitempty" gorm:"column:status_code;type:jsonb"`
	ResponseBody map[string][]MatchingRule `json:"response_body,omitempty" gorm:"column:response_body;type:jsonb"`
}

// ProviderState represents the required provider state for the interaction
type ProviderState struct {
	Name        string                 `json:"name" gorm:"column:name;size:200;not null"`
	Type        ProviderStateType      `json:"type" gorm:"column:type;size:20;not null"`
	Description string                 `json:"description,omitempty" gorm:"column:description;size:500"`
	Params      map[string]interface{} `json:"params,omitempty" gorm:"column:params;type:jsonb"`
	Action      string                 `json:"action,omitempty" gorm:"column:action;size:200"`
	Timeout     *time.Duration         `json:"timeout,omitempty" gorm:"column:timeout"`
	Required    bool                   `json:"required" gorm:"column:required;default:true"`
}

// InteractionMetadata represents metadata for the interaction
type InteractionMetadata struct {
	Tags           []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	TestData       map[string]interface{} `json:"test_data,omitempty" gorm:"column:test_data;type:jsonb"`
	ExecutionOrder int                    `json:"execution_order,omitempty" gorm:"column:execution_order;default:0"`
	Timeout        *time.Duration         `json:"timeout,omitempty" gorm:"column:timeout"`
	RetryCount     int                    `json:"retry_count,omitempty" gorm:"column:retry_count;default:0"`
	Critical       bool                   `json:"critical" gorm:"column:critical;default:false"`
	Environment    string                 `json:"environment,omitempty" gorm:"column:environment;size:20"`
	Notes          string                 `json:"notes,omitempty" gorm:"column:notes;size:1000"`
}

// ContractInteraction represents a specific API request/response expectation within a contract
type ContractInteraction struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContractID  uuid.UUID `json:"contract_id" gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`

	// Interaction identification
	Description string `json:"description" gorm:"column:description;size:500;not null"`
	UniqueKey   string `json:"unique_key" gorm:"column:unique_key;size:200;not null;uniqueIndex"`

	// Request and response specifications
	Request  HTTPRequest  `json:"request" gorm:"embedded;embeddedPrefix:request_"`
	Response HTTPResponse `json:"response" gorm:"embedded;embeddedPrefix:response_"`

	// Provider state requirements
	ProviderState      ProviderState     `json:"provider_state" gorm:"embedded;embeddedPrefix:provider_state_"`
	ProviderStates     []ProviderState   `json:"provider_states,omitempty" gorm:"column:provider_states;type:jsonb"`
	RequiresSetup      bool              `json:"requires_setup" gorm:"column:requires_setup;default:false"`

	// Matching rules for flexible validation
	MatchingRules      MatchingRuleSet   `json:"matching_rules" gorm:"embedded;embeddedPrefix:matching_rules_"`

	// Interaction metadata
	Metadata           InteractionMetadata `json:"metadata" gorm:"embedded;embeddedPrefix:metadata_"`

	// Verification and testing
	LastTestedAt       *time.Time        `json:"last_tested_at,omitempty" gorm:"column:last_tested_at;index"`
	TestCount          int               `json:"test_count" gorm:"column:test_count;default:0"`
	FailureCount       int               `json:"failure_count" gorm:"column:failure_count;default:0"`
	AverageResponseTime *time.Duration   `json:"average_response_time,omitempty" gorm:"column:average_response_time"`

	// Versioning and history
	Version            int               `json:"version" gorm:"column:version;default:1"`
	PreviousVersions   []InteractionVersion `json:"previous_versions,omitempty" gorm:"column:previous_versions;type:jsonb"`

	// Timestamps
	CreatedAt          time.Time         `json:"created_at" gorm:"column:created_at;not null;index"`
	UpdatedAt          time.Time         `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt          gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships (not stored in database)
	ContractSpecification *ContractSpecification `json:"contract_specification,omitempty" gorm:"-"`
	ValidationResults     []ValidationResult     `json:"validation_results,omitempty" gorm:"-"`
}

// InteractionVersion represents a previous version of an interaction
type InteractionVersion struct {
	Version     int       `json:"version"`
	Description string    `json:"description"`
	ChangedAt   time.Time `json:"changed_at"`
	ChangedBy   string    `json:"changed_by"`
	Reason      string    `json:"reason,omitempty"`
}

// TableName returns the table name for the ContractInteraction model
func (ContractInteraction) TableName() string {
	return "contract_interactions"
}

// BeforeCreate sets up the contract interaction before creation
func (ci *ContractInteraction) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}

	// Generate unique key
	if ci.UniqueKey == "" {
		ci.UniqueKey = ci.generateUniqueKey()
	}

	// Set default content types
	if ci.Request.ContentType == "" && ci.Request.Body != nil {
		ci.Request.ContentType = "application/json"
	}
	if ci.Response.ContentType == "" && ci.Response.Body != nil {
		ci.Response.ContentType = "application/json"
	}

	// Set default encoding
	if ci.Request.Encoding == "" {
		ci.Request.Encoding = "utf-8"
	}
	if ci.Response.Encoding == "" {
		ci.Response.Encoding = "utf-8"
	}

	// Set response status text if not provided
	if ci.Response.StatusText == "" {
		ci.Response.StatusText = GetHTTPStatusText(ci.Response.Status)
	}

	// Initialize metadata defaults
	if ci.Metadata.Environment == "" {
		ci.Metadata.Environment = "development"
	}

	// Initialize maps if nil
	if ci.Request.Query == nil {
		ci.Request.Query = make(map[string]interface{})
	}
	if ci.Request.Headers == nil {
		ci.Request.Headers = make(map[string]string)
	}
	if ci.Response.Headers == nil {
		ci.Response.Headers = make(map[string]string)
	}

	// Validate the interaction
	if err := ci.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the contract interaction before updating
func (ci *ContractInteraction) BeforeUpdate(tx *gorm.DB) error {
	// Update unique key if request method or path changed
	newUniqueKey := ci.generateUniqueKey()
	if ci.UniqueKey != newUniqueKey {
		ci.UniqueKey = newUniqueKey
	}

	return ci.Validate()
}

// Validate validates the contract interaction data
func (ci *ContractInteraction) Validate() error {
	// Validate UUIDs
	if ci.ID == uuid.Nil {
		return fmt.Errorf("contract interaction ID cannot be nil")
	}
	if ci.ContractID == uuid.Nil {
		return fmt.Errorf("contract ID cannot be nil")
	}

	// Validate description
	if ci.Description == "" {
		return fmt.Errorf("interaction description is required")
	}

	// Validate request
	if err := ci.validateRequest(); err != nil {
		return fmt.Errorf("request validation failed: %w", err)
	}

	// Validate response
	if err := ci.validateResponse(); err != nil {
		return fmt.Errorf("response validation failed: %w", err)
	}

	// Validate provider states
	if err := ci.validateProviderStates(); err != nil {
		return fmt.Errorf("provider states validation failed: %w", err)
	}

	// Validate matching rules
	if err := ci.validateMatchingRules(); err != nil {
		return fmt.Errorf("matching rules validation failed: %w", err)
	}

	// Validate metadata
	if err := ci.validateMetadata(); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	return nil
}

// validateRequest validates the HTTP request specification
func (ci *ContractInteraction) validateRequest() error {
	// Validate HTTP method
	if !IsValidHTTPMethod(ci.Request.Method) {
		return fmt.Errorf("invalid HTTP method: %s", ci.Request.Method)
	}

	// Validate path
	if !IsValidHTTPPath(ci.Request.Path) {
		return fmt.Errorf("invalid HTTP path: %s", ci.Request.Path)
	}

	// Validate content type if body is present
	if ci.Request.Body != nil && ci.Request.ContentType == "" {
		return fmt.Errorf("content type is required when request body is present")
	}

	// Validate JSON body if content type is JSON
	if ci.Request.ContentType == "application/json" && ci.Request.Body != nil {
		var jsonData interface{}
		if err := json.Unmarshal(ci.Request.Body, &jsonData); err != nil {
			return fmt.Errorf("invalid JSON in request body: %w", err)
		}
	}

	return nil
}

// validateResponse validates the HTTP response specification
func (ci *ContractInteraction) validateResponse() error {
	// Validate status code
	if !IsValidHTTPStatusCode(ci.Response.Status) {
		return fmt.Errorf("invalid HTTP status code: %d", ci.Response.Status)
	}

	// Validate content type if body is present
	if ci.Response.Body != nil && ci.Response.ContentType == "" {
		return fmt.Errorf("content type is required when response body is present")
	}

	// Validate JSON body if content type is JSON
	if ci.Response.ContentType == "application/json" && ci.Response.Body != nil {
		var jsonData interface{}
		if err := json.Unmarshal(ci.Response.Body, &jsonData); err != nil {
			return fmt.Errorf("invalid JSON in response body: %w", err)
		}
	}

	return nil
}

// validateProviderStates validates provider state configurations
func (ci *ContractInteraction) validateProviderStates() error {
	// Validate primary provider state if present
	if ci.ProviderState.Name != "" {
		if !ci.ProviderState.Type.IsValid() {
			return fmt.Errorf("invalid provider state type: %s", ci.ProviderState.Type)
		}
	}

	// Validate additional provider states
	for i, state := range ci.ProviderStates {
		if state.Name == "" {
			return fmt.Errorf("provider state name is required at index %d", i)
		}
		if !state.Type.IsValid() {
			return fmt.Errorf("invalid provider state type at index %d: %s", i, state.Type)
		}
	}

	return nil
}

// validateMatchingRules validates matching rule configurations
func (ci *ContractInteraction) validateMatchingRules() error {
	// Validate path matching rules
	for i, rule := range ci.MatchingRules.Path {
		if err := ci.validateMatchingRule(rule, fmt.Sprintf("path[%d]", i)); err != nil {
			return err
		}
	}

	// Validate query matching rules
	for field, rules := range ci.MatchingRules.Query {
		for i, rule := range rules {
			if err := ci.validateMatchingRule(rule, fmt.Sprintf("query[%s][%d]", field, i)); err != nil {
				return err
			}
		}
	}

	// Validate header matching rules
	for field, rules := range ci.MatchingRules.Header {
		for i, rule := range rules {
			if err := ci.validateMatchingRule(rule, fmt.Sprintf("header[%s][%d]", field, i)); err != nil {
				return err
			}
		}
	}

	// Validate body matching rules
	for field, rules := range ci.MatchingRules.Body {
		for i, rule := range rules {
			if err := ci.validateMatchingRule(rule, fmt.Sprintf("body[%s][%d]", field, i)); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateMatchingRule validates a single matching rule
func (ci *ContractInteraction) validateMatchingRule(rule MatchingRule, context string) error {
	if !rule.Type.IsValid() {
		return fmt.Errorf("invalid matching rule type in %s: %s", context, rule.Type)
	}

	// Type-specific validation
	switch rule.Type {
	case MatchingRuleTypeRegex:
		if rule.Regex == "" {
			return fmt.Errorf("regex is required for regex matching rule in %s", context)
		}
		// Validate regex syntax
		if _, err := regexp.Compile(rule.Regex); err != nil {
			return fmt.Errorf("invalid regex in %s: %w", context, err)
		}

	case MatchingRuleTypeNumber, MatchingRuleTypeInteger, MatchingRuleTypeDecimal:
		if rule.Min != nil && rule.Max != nil && *rule.Min > *rule.Max {
			return fmt.Errorf("min cannot be greater than max in %s", context)
		}

	case MatchingRuleTypeDate, MatchingRuleTypeTime, MatchingRuleTypeTimestamp:
		if rule.Format == "" {
			return fmt.Errorf("format is required for %s matching rule in %s", rule.Type, context)
		}
	}

	return nil
}

// validateMetadata validates interaction metadata
func (ci *ContractInteraction) validateMetadata() error {
	// Validate environment
	if ci.Metadata.Environment != "" && !IsValidEnvironment(ci.Metadata.Environment) {
		return fmt.Errorf("invalid environment: %s", ci.Metadata.Environment)
	}

	// Validate timeout
	if ci.Metadata.Timeout != nil && *ci.Metadata.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Validate execution order
	if ci.Metadata.ExecutionOrder < 0 {
		return fmt.Errorf("execution order cannot be negative")
	}

	return nil
}

// generateUniqueKey generates a unique key for this interaction
func (ci *ContractInteraction) generateUniqueKey() string {
	return fmt.Sprintf("%s:%s:%s", ci.ContractID.String()[:8], ci.Request.Method, ci.Request.Path)
}

// UpdateVersion creates a new version of the interaction
func (ci *ContractInteraction) UpdateVersion(changedBy, reason string) {
	// Save current version to history
	previousVersion := InteractionVersion{
		Version:     ci.Version,
		Description: ci.Description,
		ChangedAt:   time.Now(),
		ChangedBy:   changedBy,
		Reason:      reason,
	}
	ci.PreviousVersions = append(ci.PreviousVersions, previousVersion)

	// Increment version
	ci.Version++
	ci.UpdatedAt = time.Now()
}

// AddProviderState adds a provider state to the interaction
func (ci *ContractInteraction) AddProviderState(state ProviderState) error {
	if state.Name == "" {
		return fmt.Errorf("provider state name is required")
	}

	if !state.Type.IsValid() {
		return fmt.Errorf("invalid provider state type: %s", state.Type)
	}

	// Check for duplicate state names
	for _, existing := range ci.ProviderStates {
		if existing.Name == state.Name {
			return fmt.Errorf("provider state with name '%s' already exists", state.Name)
		}
	}

	ci.ProviderStates = append(ci.ProviderStates, state)
	ci.RequiresSetup = true
	ci.UpdatedAt = time.Now()

	return nil
}

// RemoveProviderState removes a provider state by name
func (ci *ContractInteraction) RemoveProviderState(stateName string) error {
	for i, state := range ci.ProviderStates {
		if state.Name == stateName {
			ci.ProviderStates = append(ci.ProviderStates[:i], ci.ProviderStates[i+1:]...)
			ci.RequiresSetup = len(ci.ProviderStates) > 0 || ci.ProviderState.Name != ""
			ci.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("provider state with name '%s' not found", stateName)
}

// GetSuccessRate calculates the test success rate
func (ci *ContractInteraction) GetSuccessRate() float64 {
	if ci.TestCount == 0 {
		return 0.0
	}
	successCount := ci.TestCount - ci.FailureCount
	return float64(successCount) / float64(ci.TestCount)
}

// IncrementTestCount increments test statistics
func (ci *ContractInteraction) IncrementTestCount(success bool, responseTime time.Duration) {
	ci.TestCount++
	if !success {
		ci.FailureCount++
	}

	// Update average response time
	if ci.AverageResponseTime == nil {
		ci.AverageResponseTime = &responseTime
	} else {
		// Calculate running average
		totalTime := time.Duration(ci.TestCount-1) * (*ci.AverageResponseTime) + responseTime
		avgTime := totalTime / time.Duration(ci.TestCount)
		ci.AverageResponseTime = &avgTime
	}

	now := time.Now()
	ci.LastTestedAt = &now
	ci.UpdatedAt = now
}

// IsCompatibleWith checks if this interaction is compatible with another version
func (ci *ContractInteraction) IsCompatibleWith(other *ContractInteraction) bool {
	// Basic compatibility check - same method and path
	if ci.Request.Method != other.Request.Method || ci.Request.Path != other.Request.Path {
		return false
	}

	// Response status should match
	if ci.Response.Status != other.Response.Status {
		return false
	}

	// Additional compatibility checks can be added here
	// For example, checking if response schema is backward compatible
	return true
}

// MarshalJSON customizes JSON serialization
func (ci *ContractInteraction) MarshalJSON() ([]byte, error) {
	type Alias ContractInteraction
	return json.Marshal(&struct {
		*Alias
		UniqueKey           string  `json:"unique_key"`
		SuccessRate         float64 `json:"success_rate"`
		Age                 string  `json:"age"`
		TimeSinceLastTest   string  `json:"time_since_last_test"`
		TimeSinceUpdate     string  `json:"time_since_update"`
		AverageResponseTimeMs int64 `json:"average_response_time_ms"`
		IsCritical          bool    `json:"is_critical"`
		RequiresSetup       bool    `json:"requires_setup"`
	}{
		Alias:       (*Alias)(ci),
		UniqueKey:   ci.UniqueKey,
		SuccessRate: ci.GetSuccessRate(),
		Age:         time.Since(ci.CreatedAt).String(),
		TimeSinceLastTest: func() string {
			if ci.LastTestedAt != nil {
				return time.Since(*ci.LastTestedAt).String()
			}
			return "never"
		}(),
		TimeSinceUpdate: time.Since(ci.UpdatedAt).String(),
		AverageResponseTimeMs: func() int64 {
			if ci.AverageResponseTime != nil {
				return ci.AverageResponseTime.Milliseconds()
			}
			return 0
		}(),
		IsCritical:    ci.Metadata.Critical,
		RequiresSetup: ci.RequiresSetup,
	})
}

// Helper functions

// GetHTTPStatusText returns the standard text for an HTTP status code
func GetHTTPStatusText(statusCode int) string {
	statusTexts := map[int]string{
		200: "OK",
		201: "Created",
		202: "Accepted",
		204: "No Content",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		409: "Conflict",
		422: "Unprocessable Entity",
		429: "Too Many Requests",
		500: "Internal Server Error",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
	}

	if text, exists := statusTexts[statusCode]; exists {
		return text
	}

	// Default based on status code range
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "Success"
	case statusCode >= 300 && statusCode < 400:
		return "Redirection"
	case statusCode >= 400 && statusCode < 500:
		return "Client Error"
	case statusCode >= 500 && statusCode < 600:
		return "Server Error"
	default:
		return "Unknown"
	}
}

// GetSupportedMatchingRuleTypes returns all supported matching rule types
func GetSupportedMatchingRuleTypes() []MatchingRuleType {
	return []MatchingRuleType{
		MatchingRuleTypeEquality,
		MatchingRuleTypeRegex,
		MatchingRuleTypeType,
		MatchingRuleTypeNumber,
		MatchingRuleTypeInteger,
		MatchingRuleTypeDecimal,
		MatchingRuleTypeDate,
		MatchingRuleTypeTime,
		MatchingRuleTypeTimestamp,
		MatchingRuleTypeNull,
		MatchingRuleTypeBoolean,
		MatchingRuleTypeInclude,
		MatchingRuleTypeArrayContains,
	}
}

// GetSupportedProviderStateTypes returns all supported provider state types
func GetSupportedProviderStateTypes() []ProviderStateType {
	return []ProviderStateType{
		ProviderStateTypeSetup,
		ProviderStateTypeTeardown,
		ProviderStateTypePersist,
		ProviderStateTypeCleanup,
	}
}