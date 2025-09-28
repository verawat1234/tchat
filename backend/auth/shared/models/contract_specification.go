package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContractStatus represents the status of a contract specification
type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusPublished  ContractStatus = "published"
	ContractStatusVerified   ContractStatus = "verified"
	ContractStatusDeprecated ContractStatus = "deprecated"
)

// IsValid checks if the contract status is valid
func (cs ContractStatus) IsValid() bool {
	switch cs {
	case ContractStatusDraft, ContractStatusPublished, ContractStatusVerified, ContractStatusDeprecated:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the status can transition to another status
func (cs ContractStatus) CanTransitionTo(newStatus ContractStatus) bool {
	validTransitions := map[ContractStatus][]ContractStatus{
		ContractStatusDraft:      {ContractStatusPublished},
		ContractStatusPublished:  {ContractStatusVerified, ContractStatusDraft},
		ContractStatusVerified:   {ContractStatusDeprecated, ContractStatusPublished},
		ContractStatusDeprecated: {}, // Terminal state
	}

	allowedTransitions, exists := validTransitions[cs]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// ConsumerPlatform represents the consumer platform type
type ConsumerPlatform string

const (
	ConsumerPlatformWeb     ConsumerPlatform = "web"
	ConsumerPlatformIOS     ConsumerPlatform = "ios"
	ConsumerPlatformAndroid ConsumerPlatform = "android"
)

// IsValid checks if the consumer platform is valid
func (cp ConsumerPlatform) IsValid() bool {
	switch cp {
	case ConsumerPlatformWeb, ConsumerPlatformIOS, ConsumerPlatformAndroid:
		return true
	default:
		return false
	}
}

// ProviderService represents the backend service type
type ProviderService string

const (
	ProviderServiceAuth         ProviderService = "auth"
	ProviderServiceContent      ProviderService = "content"
	ProviderServiceCommerce     ProviderService = "commerce"
	ProviderServiceMessaging    ProviderService = "messaging"
	ProviderServicePayment      ProviderService = "payment"
	ProviderServiceNotification ProviderService = "notification"
	ProviderServiceGateway      ProviderService = "gateway"
)

// IsValid checks if the provider service is valid
func (ps ProviderService) IsValid() bool {
	switch ps {
	case ProviderServiceAuth, ProviderServiceContent, ProviderServiceCommerce,
		 ProviderServiceMessaging, ProviderServicePayment, ProviderServiceNotification,
		 ProviderServiceGateway:
		return true
	default:
		return false
	}
}

// ContractMetadata represents contract metadata
type ContractMetadata struct {
	Author           string                 `json:"author" gorm:"column:author;size:100"`
	Description      string                 `json:"description" gorm:"column:description;size:500"`
	Tags             []string               `json:"tags" gorm:"column:tags;type:jsonb"`
	Repository       string                 `json:"repository,omitempty" gorm:"column:repository;size:255"`
	Branch           string                 `json:"branch,omitempty" gorm:"column:branch;size:100"`
	CommitHash       string                 `json:"commit_hash,omitempty" gorm:"column:commit_hash;size:40"`
	BuildNumber      string                 `json:"build_number,omitempty" gorm:"column:build_number;size:50"`
	Environment      string                 `json:"environment" gorm:"column:environment;size:20;default:'development'"`
	Notes            string                 `json:"notes,omitempty" gorm:"column:notes;size:1000"`
	CustomProperties map[string]interface{} `json:"custom_properties,omitempty" gorm:"column:custom_properties;type:jsonb"`
}

// SemanticVersion represents a semantic version
type SemanticVersion struct {
	Major int `json:"major" gorm:"column:major;not null"`
	Minor int `json:"minor" gorm:"column:minor;not null"`
	Patch int `json:"patch" gorm:"column:patch;not null"`
}

// String returns the string representation of the semantic version
func (sv SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

// Compare compares two semantic versions
// Returns: -1 if sv < other, 0 if equal, 1 if sv > other
func (sv SemanticVersion) Compare(other SemanticVersion) int {
	if sv.Major != other.Major {
		if sv.Major < other.Major {
			return -1
		}
		return 1
	}
	if sv.Minor != other.Minor {
		if sv.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if sv.Patch != other.Patch {
		if sv.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// IsCompatibleWith checks if this version is backward compatible with another
func (sv SemanticVersion) IsCompatibleWith(other SemanticVersion) bool {
	// Major version must match for backward compatibility
	if sv.Major != other.Major {
		return false
	}
	// This version should be >= other version for compatibility
	return sv.Compare(other) >= 0
}

// ContractSpecification represents a contract specification between consumer and provider
type ContractSpecification struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConsumerName ConsumerPlatform `json:"consumer_name" gorm:"column:consumer_name;type:varchar(20);not null;index"`
	ProviderName ProviderService  `json:"provider_name" gorm:"column:provider_name;type:varchar(20);not null;index"`

	// Version information
	Version         SemanticVersion `json:"version" gorm:"embedded;embeddedPrefix:version_"`
	VersionString   string          `json:"version_string" gorm:"column:version_string;size:20;not null;index"`
	PactVersion     string          `json:"pact_version" gorm:"column:pact_version;size:10;default:'4.0'"`

	// Contract status and lifecycle
	Status          ContractStatus   `json:"status" gorm:"column:status;type:varchar(20);not null;default:'draft';index"`
	StatusHistory   []StatusChange   `json:"status_history" gorm:"column:status_history;type:jsonb"`

	// Contract content
	Interactions    []ContractInteractionRef `json:"interactions" gorm:"column:interactions;type:jsonb"`
	InteractionCount int                     `json:"interaction_count" gorm:"column:interaction_count;default:0"`

	// Metadata and documentation
	Metadata        ContractMetadata `json:"metadata" gorm:"embedded;embeddedPrefix:metadata_"`

	// Contract publishing and sharing
	BrokerURL       string           `json:"broker_url,omitempty" gorm:"column:broker_url;size:500"`
	PublishedURL    string           `json:"published_url,omitempty" gorm:"column:published_url;size:500"`
	PactFileHash    string           `json:"pact_file_hash,omitempty" gorm:"column:pact_file_hash;size:64"`

	// Validation and verification
	LastValidatedAt *time.Time       `json:"last_validated_at,omitempty" gorm:"column:last_validated_at;index"`
	ValidationCount int              `json:"validation_count" gorm:"column:validation_count;default:0"`
	FailureCount    int              `json:"failure_count" gorm:"column:failure_count;default:0"`

	// Regional and compliance
	DataRegion      string           `json:"data_region" gorm:"column:data_region;size:20;default:'sea-central'"`
	ComplianceFlags map[string]bool  `json:"compliance_flags,omitempty" gorm:"column:compliance_flags;type:jsonb"`

	// Timestamps
	CreatedAt       time.Time        `json:"created_at" gorm:"column:created_at;not null;index"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"column:updated_at;not null"`
	PublishedAt     *time.Time       `json:"published_at,omitempty" gorm:"column:published_at;index"`
	VerifiedAt      *time.Time       `json:"verified_at,omitempty" gorm:"column:verified_at;index"`
	DeprecatedAt    *time.Time       `json:"deprecated_at,omitempty" gorm:"column:deprecated_at;index"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships (not stored in database, populated via joins)
	ContractInteractions []ContractInteraction `json:"contract_interactions,omitempty" gorm:"-"`
	ValidationResults    []ValidationResult    `json:"validation_results,omitempty" gorm:"-"`
}

// StatusChange represents a status change event
type StatusChange struct {
	FromStatus ContractStatus `json:"from_status"`
	ToStatus   ContractStatus `json:"to_status"`
	ChangedBy  string         `json:"changed_by"`
	ChangedAt  time.Time      `json:"changed_at"`
	Reason     string         `json:"reason,omitempty"`
}

// ContractInteractionRef represents a reference to a contract interaction
type ContractInteractionRef struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	StatusCode  int       `json:"status_code"`
}

// TableName returns the table name for the ContractSpecification model
func (ContractSpecification) TableName() string {
	return "contract_specifications"
}

// BeforeCreate sets up the contract specification before creation
func (cs *ContractSpecification) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}

	// Set version string from semantic version
	cs.VersionString = cs.Version.String()

	// Set default metadata
	if cs.Metadata.Environment == "" {
		cs.Metadata.Environment = "development"
	}

	// Set data region based on provider service
	if cs.DataRegion == "" {
		cs.DataRegion = "sea-central" // Default for Southeast Asian deployment
	}

	// Initialize status history
	if len(cs.StatusHistory) == 0 {
		cs.StatusHistory = []StatusChange{
			{
				FromStatus: "",
				ToStatus:   cs.Status,
				ChangedBy:  cs.Metadata.Author,
				ChangedAt:  time.Now(),
				Reason:     "Initial creation",
			},
		}
	}

	// Initialize compliance flags
	if cs.ComplianceFlags == nil {
		cs.ComplianceFlags = make(map[string]bool)
	}

	// Validate the contract specification
	if err := cs.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the contract specification before updating
func (cs *ContractSpecification) BeforeUpdate(tx *gorm.DB) error {
	// Update version string
	cs.VersionString = cs.Version.String()

	// Update interaction count
	cs.InteractionCount = len(cs.Interactions)

	return cs.Validate()
}

// Validate validates the contract specification data
func (cs *ContractSpecification) Validate() error {
	// Validate UUIDs
	if cs.ID == uuid.Nil {
		return fmt.Errorf("contract specification ID cannot be nil")
	}

	// Validate consumer and provider names
	if !cs.ConsumerName.IsValid() {
		return fmt.Errorf("invalid consumer platform: %s", cs.ConsumerName)
	}

	if !cs.ProviderName.IsValid() {
		return fmt.Errorf("invalid provider service: %s", cs.ProviderName)
	}

	// Validate status
	if !cs.Status.IsValid() {
		return fmt.Errorf("invalid contract status: %s", cs.Status)
	}

	// Validate semantic version
	if cs.Version.Major < 0 || cs.Version.Minor < 0 || cs.Version.Patch < 0 {
		return fmt.Errorf("semantic version components cannot be negative: %s", cs.Version.String())
	}

	// Validate version string format
	if !IsValidSemanticVersion(cs.VersionString) {
		return fmt.Errorf("invalid semantic version format: %s", cs.VersionString)
	}

	// Validate pact version
	if cs.PactVersion == "" {
		cs.PactVersion = "4.0" // Default to latest supported version
	}

	// Validate metadata
	if cs.Metadata.Author == "" {
		return fmt.Errorf("contract author is required")
	}

	// Validate environment
	if !IsValidEnvironment(cs.Metadata.Environment) {
		return fmt.Errorf("invalid environment: %s", cs.Metadata.Environment)
	}

	// Validate interactions
	if err := cs.validateInteractions(); err != nil {
		return err
	}

	return nil
}

// validateInteractions validates the contract interactions
func (cs *ContractSpecification) validateInteractions() error {
	if len(cs.Interactions) == 0 {
		return fmt.Errorf("contract specification must have at least one interaction")
	}

	// Check for duplicate interactions
	seen := make(map[string]bool)
	for i, interaction := range cs.Interactions {
		key := fmt.Sprintf("%s:%s", interaction.Method, interaction.Path)
		if seen[key] {
			return fmt.Errorf("duplicate interaction at index %d: %s %s", i, interaction.Method, interaction.Path)
		}
		seen[key] = true

		// Validate interaction reference
		if interaction.ID == uuid.Nil {
			return fmt.Errorf("interaction ID cannot be nil at index %d", i)
		}

		if interaction.Description == "" {
			return fmt.Errorf("interaction description is required at index %d", i)
		}

		if !IsValidHTTPMethod(interaction.Method) {
			return fmt.Errorf("invalid HTTP method at index %d: %s", i, interaction.Method)
		}

		if !IsValidHTTPPath(interaction.Path) {
			return fmt.Errorf("invalid HTTP path at index %d: %s", i, interaction.Path)
		}

		if !IsValidHTTPStatusCode(interaction.StatusCode) {
			return fmt.Errorf("invalid HTTP status code at index %d: %d", i, interaction.StatusCode)
		}
	}

	return nil
}

// UpdateStatus updates the contract status with proper validation and history tracking
func (cs *ContractSpecification) UpdateStatus(newStatus ContractStatus, changedBy, reason string) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid contract status: %s", newStatus)
	}

	if !cs.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", cs.Status, newStatus)
	}

	oldStatus := cs.Status
	cs.Status = newStatus

	// Add to status history
	statusChange := StatusChange{
		FromStatus: oldStatus,
		ToStatus:   newStatus,
		ChangedBy:  changedBy,
		ChangedAt:  time.Now(),
		Reason:     reason,
	}
	cs.StatusHistory = append(cs.StatusHistory, statusChange)

	// Set appropriate timestamps
	now := time.Now()
	switch newStatus {
	case ContractStatusPublished:
		if cs.PublishedAt == nil {
			cs.PublishedAt = &now
		}
	case ContractStatusVerified:
		if cs.VerifiedAt == nil {
			cs.VerifiedAt = &now
		}
		cs.LastValidatedAt = &now
		cs.ValidationCount++
	case ContractStatusDeprecated:
		if cs.DeprecatedAt == nil {
			cs.DeprecatedAt = &now
		}
	}

	cs.UpdatedAt = now
	return nil
}

// AddInteraction adds a new interaction reference to the contract
func (cs *ContractSpecification) AddInteraction(interactionRef ContractInteractionRef) error {
	// Validate interaction reference
	if interactionRef.ID == uuid.Nil {
		return fmt.Errorf("interaction ID cannot be nil")
	}

	// Check for duplicate
	for _, existing := range cs.Interactions {
		if existing.ID == interactionRef.ID {
			return fmt.Errorf("interaction with ID %s already exists", interactionRef.ID)
		}
		// Check for duplicate method+path combination
		if existing.Method == interactionRef.Method && existing.Path == interactionRef.Path {
			return fmt.Errorf("duplicate interaction: %s %s", interactionRef.Method, interactionRef.Path)
		}
	}

	cs.Interactions = append(cs.Interactions, interactionRef)
	cs.InteractionCount = len(cs.Interactions)
	cs.UpdatedAt = time.Now()

	return nil
}

// RemoveInteraction removes an interaction reference from the contract
func (cs *ContractSpecification) RemoveInteraction(interactionID uuid.UUID) error {
	for i, interaction := range cs.Interactions {
		if interaction.ID == interactionID {
			cs.Interactions = append(cs.Interactions[:i], cs.Interactions[i+1:]...)
			cs.InteractionCount = len(cs.Interactions)
			cs.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("interaction with ID %s not found", interactionID)
}

// GetUniqueKey returns a unique key for this contract specification
func (cs *ContractSpecification) GetUniqueKey() string {
	return fmt.Sprintf("%s-%s-%s", cs.ConsumerName, cs.ProviderName, cs.VersionString)
}

// IsPublished checks if the contract is published
func (cs *ContractSpecification) IsPublished() bool {
	return cs.Status == ContractStatusPublished || cs.Status == ContractStatusVerified
}

// IsVerified checks if the contract is verified
func (cs *ContractSpecification) IsVerified() bool {
	return cs.Status == ContractStatusVerified
}

// IsDeprecated checks if the contract is deprecated
func (cs *ContractSpecification) IsDeprecated() bool {
	return cs.Status == ContractStatusDeprecated
}

// CanBeModified checks if the contract can be modified
func (cs *ContractSpecification) CanBeModified() bool {
	return cs.Status == ContractStatusDraft
}

// GetSuccessRate calculates the validation success rate
func (cs *ContractSpecification) GetSuccessRate() float64 {
	if cs.ValidationCount == 0 {
		return 0.0
	}
	successCount := cs.ValidationCount - cs.FailureCount
	return float64(successCount) / float64(cs.ValidationCount)
}

// IncrementValidationCount increments validation statistics
func (cs *ContractSpecification) IncrementValidationCount(success bool) {
	cs.ValidationCount++
	if !success {
		cs.FailureCount++
	}
	now := time.Now()
	cs.LastValidatedAt = &now
	cs.UpdatedAt = now
}

// GetCompatibleVersions returns versions that are compatible with this contract
func (cs *ContractSpecification) GetCompatibleVersions(otherVersions []SemanticVersion) []SemanticVersion {
	compatible := make([]SemanticVersion, 0)
	for _, version := range otherVersions {
		if cs.Version.IsCompatibleWith(version) {
			compatible = append(compatible, version)
		}
	}
	return compatible
}

// MarshalJSON customizes JSON serialization
func (cs *ContractSpecification) MarshalJSON() ([]byte, error) {
	type Alias ContractSpecification
	return json.Marshal(&struct {
		*Alias
		UniqueKey       string  `json:"unique_key"`
		IsPublished     bool    `json:"is_published"`
		IsVerified      bool    `json:"is_verified"`
		IsDeprecated    bool    `json:"is_deprecated"`
		CanBeModified   bool    `json:"can_be_modified"`
		SuccessRate     float64 `json:"success_rate"`
		Age             string  `json:"age"`
		TimeSinceUpdate string  `json:"time_since_update"`
	}{
		Alias:           (*Alias)(cs),
		UniqueKey:       cs.GetUniqueKey(),
		IsPublished:     cs.IsPublished(),
		IsVerified:      cs.IsVerified(),
		IsDeprecated:    cs.IsDeprecated(),
		CanBeModified:   cs.CanBeModified(),
		SuccessRate:     cs.GetSuccessRate(),
		Age:             time.Since(cs.CreatedAt).String(),
		TimeSinceUpdate: time.Since(cs.UpdatedAt).String(),
	})
}

// Helper functions for validation

// IsValidSemanticVersion validates semantic version format
func IsValidSemanticVersion(version string) bool {
	semVerRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	return semVerRegex.MatchString(version)
}

// ParseSemanticVersion parses a semantic version string
func ParseSemanticVersion(version string) (SemanticVersion, error) {
	semVerRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	matches := semVerRegex.FindStringSubmatch(version)

	if len(matches) != 4 {
		return SemanticVersion{}, fmt.Errorf("invalid semantic version format: %s", version)
	}

	var major, minor, patch int
	if _, err := fmt.Sscanf(matches[1], "%d", &major); err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid major version: %s", matches[1])
	}
	if _, err := fmt.Sscanf(matches[2], "%d", &minor); err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid minor version: %s", matches[2])
	}
	if _, err := fmt.Sscanf(matches[3], "%d", &patch); err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// IsValidEnvironment validates environment name
func IsValidEnvironment(environment string) bool {
	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"testing":     true,
	}
	return validEnvironments[environment]
}

// IsValidHTTPMethod validates HTTP method
func IsValidHTTPMethod(method string) bool {
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"PATCH":   true,
		"DELETE":  true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	return validMethods[strings.ToUpper(method)]
}

// IsValidHTTPPath validates HTTP path format
func IsValidHTTPPath(path string) bool {
	if path == "" || !strings.HasPrefix(path, "/") {
		return false
	}
	// Basic validation - more complex validation can be added as needed
	pathRegex := regexp.MustCompile(`^/[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=%]*$`)
	return pathRegex.MatchString(path)
}

// IsValidHTTPStatusCode validates HTTP status code
func IsValidHTTPStatusCode(statusCode int) bool {
	return statusCode >= 100 && statusCode < 600
}

// GetSupportedConsumerPlatforms returns all supported consumer platforms
func GetSupportedConsumerPlatforms() []ConsumerPlatform {
	return []ConsumerPlatform{
		ConsumerPlatformWeb,
		ConsumerPlatformIOS,
		ConsumerPlatformAndroid,
	}
}

// GetSupportedProviderServices returns all supported provider services
func GetSupportedProviderServices() []ProviderService {
	return []ProviderService{
		ProviderServiceAuth,
		ProviderServiceContent,
		ProviderServiceCommerce,
		ProviderServiceMessaging,
		ProviderServicePayment,
		ProviderServiceNotification,
		ProviderServiceGateway,
	}
}