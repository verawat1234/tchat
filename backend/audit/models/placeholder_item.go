package models

import (
	"time"
)

// PlaceholderItem represents a placeholder or TODO item found in the codebase
type PlaceholderItem struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ServiceID   string    `json:"serviceId" gorm:"not null;index"`
	Platform    string    `json:"platform" gorm:"not null;index"` // BACKEND, WEB, MOBILE, etc.
	FilePath    string    `json:"filePath" gorm:"not null"`
	LineNumber  int       `json:"lineNumber" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null;index"`        // TODO, FIXME, PLACEHOLDER, etc.
	Priority    string    `json:"priority" gorm:"not null;default:'MEDIUM'"` // LOW, MEDIUM, HIGH, CRITICAL
	Status      string    `json:"status" gorm:"not null;default:'OPEN'"`     // OPEN, IN_PROGRESS, COMPLETED, CANCELLED
	Content     string    `json:"content" gorm:"type:text"`                   // The actual placeholder text
	Description string    `json:"description" gorm:"type:text"`               // Additional description or context
	AssignedTo  *string   `json:"assignedTo" gorm:"index"`                    // Optional assignee
	EstimatedHours *float64 `json:"estimatedHours"`                          // Estimated completion hours
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	CompletedAt *time.Time `json:"completedAt"`

	// Related audit information
	DetectedBy    string `json:"detectedBy" gorm:"not null"`     // Detection method (scan, manual, etc.)
	DetectionDate time.Time `json:"detectionDate" gorm:"not null"`
	LastVerified  *time.Time `json:"lastVerified"`               // Last verification timestamp

	// Regional optimization context (for Southeast Asian markets)
	RegionalContext *string `json:"regionalContext"` // TH, SG, MY, ID, PH, VN specific context

	// Performance impact assessment
	PerformanceImpact string `json:"performanceImpact" gorm:"default:'UNKNOWN'"` // NONE, LOW, MEDIUM, HIGH, CRITICAL

	// Dependencies and relationships
	BlockedBy []string `json:"blockedBy" gorm:"type:json"` // IDs of other placeholder items blocking this one
	Blocks    []string `json:"blocks" gorm:"type:json"`    // IDs of other placeholder items this one blocks

	// Technical metadata
	CodeContext  string `json:"codeContext" gorm:"type:text"`   // Surrounding code context
	FunctionName *string `json:"functionName"`                  // Function/method containing the placeholder
	ClassName    *string `json:"className"`                     // Class containing the placeholder
	PackageName  *string `json:"packageName"`                   // Package/module containing the placeholder

	// Quality metrics
	TechnicalDebtScore *float64 `json:"technicalDebtScore"` // Calculated technical debt score
	ComplexityScore    *int     `json:"complexityScore"`    // Implementation complexity (1-10)
	RiskScore          *float64 `json:"riskScore"`          // Risk assessment score

	// Resolution tracking
	ResolutionNotes string `json:"resolutionNotes" gorm:"type:text"` // Notes on how the placeholder was resolved
	ReviewedBy      *string `json:"reviewedBy"`                       // Who reviewed the resolution
	ReviewedAt      *time.Time `json:"reviewedAt"`                    // When the resolution was reviewed
}

// PlaceholderItemType defines the types of placeholder items
type PlaceholderItemType string

const (
	PlaceholderTypeTODO        PlaceholderItemType = "TODO"
	PlaceholderTypeFIXME       PlaceholderItemType = "FIXME"
	PlaceholderTypePlaceholder PlaceholderItemType = "PLACEHOLDER"
	PlaceholderTypeHACK        PlaceholderItemType = "HACK"
	PlaceholderTypeNOTE        PlaceholderItemType = "NOTE"
	PlaceholderTypeWARNING     PlaceholderItemType = "WARNING"
	PlaceholderTypeBUG         PlaceholderItemType = "BUG"
	PlaceholderTypeImplementation PlaceholderItemType = "IMPLEMENTATION"
)

// PlaceholderPriority defines priority levels for placeholder items
type PlaceholderPriority string

const (
	PriorityLow      PlaceholderPriority = "LOW"
	PriorityMedium   PlaceholderPriority = "MEDIUM"
	PriorityHigh     PlaceholderPriority = "HIGH"
	PriorityCritical PlaceholderPriority = "CRITICAL"
)

// PlaceholderStatus defines the status of placeholder items
type PlaceholderStatus string

const (
	StatusOpen       PlaceholderStatus = "OPEN"
	StatusInProgress PlaceholderStatus = "IN_PROGRESS"
	StatusCompleted  PlaceholderStatus = "COMPLETED"
	StatusCancelled  PlaceholderStatus = "CANCELLED"
	StatusBlocked    PlaceholderStatus = "BLOCKED"
	StatusUnderReview PlaceholderStatus = "UNDER_REVIEW"
)

// Platform defines the platforms where placeholders can be found
type Platform string

const (
	PlatformBackend Platform = "BACKEND"
	PlatformWeb     Platform = "WEB"
	PlatformIOS     Platform = "IOS"
	PlatformAndroid Platform = "ANDROID"
	PlatformKMP     Platform = "KMP"
	PlatformShared  Platform = "SHARED"
)

// PerformanceImpact defines the performance impact levels
type PerformanceImpact string

const (
	ImpactNone     PerformanceImpact = "NONE"
	ImpactLow      PerformanceImpact = "LOW"
	ImpactMedium   PerformanceImpact = "MEDIUM"
	ImpactHigh     PerformanceImpact = "HIGH"
	ImpactCritical PerformanceImpact = "CRITICAL"
	ImpactUnknown  PerformanceImpact = "UNKNOWN"
)

// GetPriorityScore returns a numeric score for priority (higher = more urgent)
func (p PlaceholderPriority) GetPriorityScore() int {
	switch p {
	case PriorityCritical:
		return 4
	case PriorityHigh:
		return 3
	case PriorityMedium:
		return 2
	case PriorityLow:
		return 1
	default:
		return 0
	}
}

// IsBlocked returns true if the placeholder item is blocked by other items
func (pi *PlaceholderItem) IsBlocked() bool {
	return len(pi.BlockedBy) > 0
}

// IsBlocking returns true if the placeholder item is blocking other items
func (pi *PlaceholderItem) IsBlocking() bool {
	return len(pi.Blocks) > 0
}

// IsCompleted returns true if the placeholder item has been completed
func (pi *PlaceholderItem) IsCompleted() bool {
	return pi.Status == string(StatusCompleted) && pi.CompletedAt != nil
}

// GetEstimatedDuration returns the estimated duration in hours, defaulting to 1 if not set
func (pi *PlaceholderItem) GetEstimatedDuration() float64 {
	if pi.EstimatedHours != nil {
		return *pi.EstimatedHours
	}
	return 1.0 // Default to 1 hour if not specified
}

// CalculateAge returns the age of the placeholder item in days
func (pi *PlaceholderItem) CalculateAge() int {
	return int(time.Since(pi.DetectionDate).Hours() / 24)
}

// IsStale returns true if the placeholder item hasn't been updated in over 30 days
func (pi *PlaceholderItem) IsStale() bool {
	return time.Since(pi.UpdatedAt).Hours() > 30*24 // 30 days
}

// GetRegionalPriority returns higher priority for Southeast Asian regional context items
func (pi *PlaceholderItem) GetRegionalPriority() int {
	if pi.RegionalContext != nil {
		regionalMarkets := []string{"TH", "SG", "MY", "ID", "PH", "VN"}
		for _, market := range regionalMarkets {
			if *pi.RegionalContext == market {
				return 1 // Higher priority for SEA markets
			}
		}
	}
	return 0
}