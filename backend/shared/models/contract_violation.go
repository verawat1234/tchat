package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ViolationType represents the type of contract violation
type ViolationType string

const (
	ViolationTypeRequestMismatch  ViolationType = "request_mismatch"
	ViolationTypeResponseMismatch ViolationType = "response_mismatch"
	ViolationTypeMissingEndpoint  ViolationType = "missing_endpoint"
	ViolationTypeSchemaViolation  ViolationType = "schema_violation"
)

// IsValid checks if the violation type is valid
func (vt ViolationType) IsValid() bool {
	switch vt {
	case ViolationTypeRequestMismatch, ViolationTypeResponseMismatch,
		 ViolationTypeMissingEndpoint, ViolationTypeSchemaViolation:
		return true
	default:
		return false
	}
}

// ViolationSeverity represents the severity level of a violation
type ViolationSeverity string

const (
	ViolationSeverityCritical ViolationSeverity = "critical"
	ViolationSeverityHigh     ViolationSeverity = "high"
	ViolationSeverityMedium   ViolationSeverity = "medium"
	ViolationSeverityLow      ViolationSeverity = "low"
)

// IsValid checks if the violation severity is valid
func (vs ViolationSeverity) IsValid() bool {
	switch vs {
	case ViolationSeverityCritical, ViolationSeverityHigh,
		 ViolationSeverityMedium, ViolationSeverityLow:
		return true
	default:
		return false
	}
}

// ViolationResolutionStatus represents the resolution status of a violation
type ViolationResolutionStatus string

const (
	ViolationResolutionStatusOpen      ViolationResolutionStatus = "open"
	ViolationResolutionStatusResolved  ViolationResolutionStatus = "resolved"
	ViolationResolutionStatusEscalated ViolationResolutionStatus = "escalated"
)

// IsValid checks if the violation resolution status is valid
func (vrs ViolationResolutionStatus) IsValid() bool {
	switch vrs {
	case ViolationResolutionStatusOpen, ViolationResolutionStatusResolved,
		 ViolationResolutionStatusEscalated:
		return true
	default:
		return false
	}
}

// ContractViolation represents a detected contract violation
type ContractViolation struct {
	// Primary key and identifiers
	ID         string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContractID string `json:"contract_id" gorm:"type:uuid;not null;index"`

	// Violation details
	ViolationType        ViolationType             `json:"violation_type" gorm:"column:violation_type;type:varchar(50);not null"`
	Severity             ViolationSeverity         `json:"severity" gorm:"column:severity;type:varchar(20);not null"`
	Description          string                    `json:"description" gorm:"column:description;type:text;not null"`
	ExpectedBehavior     string                    `json:"expected_behavior" gorm:"column:expected_behavior;type:text"`
	ActualBehavior       string                    `json:"actual_behavior" gorm:"column:actual_behavior;type:text"`
	ResolutionStatus     ViolationResolutionStatus `json:"resolution_status" gorm:"column:resolution_status;type:varchar(20);default:'open'"`

	// Context information
	InteractionID        *string                    `json:"interaction_id,omitempty" gorm:"type:uuid;index"`
	ValidationResultID   *string                    `json:"validation_result_id,omitempty" gorm:"type:uuid;index"`
	Environment          string                     `json:"environment" gorm:"column:environment;type:varchar(20);default:'development'"`
	Platform             ConsumerPlatform           `json:"platform" gorm:"column:platform;type:varchar(20)"`
	Service              ProviderService            `json:"service" gorm:"column:service;type:varchar(20)"`

	// Additional details
	ErrorCode            *string                    `json:"error_code,omitempty" gorm:"column:error_code;type:varchar(50)"`
	ErrorMessage         *string                    `json:"error_message,omitempty" gorm:"column:error_message;type:text"`
	RequestDetails       map[string]interface{}     `json:"request_details,omitempty" gorm:"column:request_details;type:jsonb"`
	ResponseDetails      map[string]interface{}     `json:"response_details,omitempty" gorm:"column:response_details;type:jsonb"`
	StackTrace           *string                    `json:"stack_trace,omitempty" gorm:"column:stack_trace;type:text"`

	// Resolution tracking
	ResolutionNotes      *string                    `json:"resolution_notes,omitempty" gorm:"column:resolution_notes;type:text"`
	ResolvedBy           *string                    `json:"resolved_by,omitempty" gorm:"column:resolved_by;type:varchar(100)"`
	ResolvedAt           *time.Time                 `json:"resolved_at,omitempty" gorm:"column:resolved_at"`
	EscalationReason     *string                    `json:"escalation_reason,omitempty" gorm:"column:escalation_reason;type:text"`
	EscalatedTo          *string                    `json:"escalated_to,omitempty" gorm:"column:escalated_to;type:varchar(100)"`
	EscalatedAt          *time.Time                 `json:"escalated_at,omitempty" gorm:"column:escalated_at"`

	// Timestamps
	DetectionTimestamp   time.Time                  `json:"detection_timestamp" gorm:"column:detection_timestamp;not null;default:now()"`
	CreatedAt            time.Time                  `json:"created_at" gorm:"column:created_at;not null;default:now()"`
	UpdatedAt            time.Time                  `json:"updated_at" gorm:"column:updated_at;not null;default:now()"`

	// Relations (will be populated by GORM when using Preload)
	Contract             *ContractSpecification     `json:"contract,omitempty" gorm:"foreignKey:ContractID"`
	Interaction          *ContractInteraction       `json:"interaction,omitempty" gorm:"foreignKey:InteractionID"`
	ValidationResult     *ValidationResult          `json:"validation_result,omitempty" gorm:"foreignKey:ValidationResultID"`
}

// TableName specifies the table name for ContractViolation
func (ContractViolation) TableName() string {
	return "contract_violations"
}

// BeforeCreate is called before creating a new ContractViolation record
func (cv *ContractViolation) BeforeCreate(tx *gorm.DB) (err error) {
	if cv.ID == "" {
		cv.ID = uuid.New().String()
	}

	now := time.Now()
	cv.CreatedAt = now
	cv.UpdatedAt = now
	cv.DetectionTimestamp = now

	return nil
}

// BeforeUpdate is called before updating a ContractViolation record
func (cv *ContractViolation) BeforeUpdate(tx *gorm.DB) (err error) {
	cv.UpdatedAt = time.Now()
	return nil
}

// Validate performs validation on the ContractViolation fields
func (cv *ContractViolation) Validate() error {
	if cv.ContractID == "" {
		return fmt.Errorf("contract_id is required")
	}

	if !cv.ViolationType.IsValid() {
		return fmt.Errorf("invalid violation_type: %s", cv.ViolationType)
	}

	if !cv.Severity.IsValid() {
		return fmt.Errorf("invalid severity: %s", cv.Severity)
	}

	if !cv.ResolutionStatus.IsValid() {
		return fmt.Errorf("invalid resolution_status: %s", cv.ResolutionStatus)
	}

	if cv.Description == "" {
		return fmt.Errorf("description is required")
	}

	if cv.Platform != "" && !cv.Platform.IsValid() {
		return fmt.Errorf("invalid platform: %s", cv.Platform)
	}

	if cv.Service != "" && !cv.Service.IsValid() {
		return fmt.Errorf("invalid service: %s", cv.Service)
	}

	return nil
}

// IsResolved checks if the violation has been resolved
func (cv *ContractViolation) IsResolved() bool {
	return cv.ResolutionStatus == ViolationResolutionStatusResolved
}

// IsEscalated checks if the violation has been escalated
func (cv *ContractViolation) IsEscalated() bool {
	return cv.ResolutionStatus == ViolationResolutionStatusEscalated
}

// CanEscalate checks if the violation can be escalated
func (cv *ContractViolation) CanEscalate() bool {
	return cv.ResolutionStatus == ViolationResolutionStatusOpen
}

// CanResolve checks if the violation can be resolved
func (cv *ContractViolation) CanResolve() bool {
	return cv.ResolutionStatus == ViolationResolutionStatusOpen ||
		   cv.ResolutionStatus == ViolationResolutionStatusEscalated
}

// Resolve marks the violation as resolved
func (cv *ContractViolation) Resolve(resolvedBy string, notes *string) {
	now := time.Now()
	cv.ResolutionStatus = ViolationResolutionStatusResolved
	cv.ResolvedBy = &resolvedBy
	cv.ResolvedAt = &now
	if notes != nil {
		cv.ResolutionNotes = notes
	}
	cv.UpdatedAt = now
}

// Escalate marks the violation as escalated
func (cv *ContractViolation) Escalate(escalatedTo string, reason *string) {
	now := time.Now()
	cv.ResolutionStatus = ViolationResolutionStatusEscalated
	cv.EscalatedTo = &escalatedTo
	cv.EscalatedAt = &now
	if reason != nil {
		cv.EscalationReason = reason
	}
	cv.UpdatedAt = now
}