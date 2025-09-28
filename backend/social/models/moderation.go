package models

import (
	"time"

	"github.com/google/uuid"
)

// Report represents a content moderation report
type Report struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	ReporterID   uuid.UUID              `json:"reporterId" db:"reporter_id"`
	TargetID     uuid.UUID              `json:"targetId" db:"target_id"`
	TargetType   string                 `json:"targetType" db:"target_type"` // post, comment, user, community
	Reason       string                 `json:"reason" db:"reason"` // spam, harassment, inappropriate, violation
	Details      string                 `json:"details" db:"details"`
	Status       string                 `json:"status" db:"status"` // pending, reviewing, resolved, dismissed
	Priority     string                 `json:"priority" db:"priority"` // low, medium, high, critical
	ReviewedBy   *uuid.UUID             `json:"reviewedBy" db:"reviewed_by"`
	ReviewNotes  string                 `json:"reviewNotes" db:"review_notes"`
	Resolution   string                 `json:"resolution" db:"resolution"` // no_action, warning, content_removed, user_banned
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time              `json:"updatedAt" db:"updated_at"`
	ResolvedAt   *time.Time             `json:"resolvedAt" db:"resolved_at"`
}

// ModerationAction represents a moderation action taken
type ModerationAction struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	ModeratorID uuid.UUID              `json:"moderatorId" db:"moderator_id"`
	TargetID    uuid.UUID              `json:"targetId" db:"target_id"`
	TargetType  string                 `json:"targetType" db:"target_type"` // post, comment, user, community
	Action      string                 `json:"action" db:"action"` // warning, mute, ban, content_removal, account_suspension
	Reason      string                 `json:"reason" db:"reason"`
	Duration    *int                   `json:"duration" db:"duration"` // duration in hours, null for permanent
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	IsActive    bool                   `json:"isActive" db:"is_active"`
	CreatedAt   time.Time              `json:"createdAt" db:"created_at"`
	ExpiresAt   *time.Time             `json:"expiresAt" db:"expires_at"`
}

// SafetyGuidelines represents community safety guidelines
type SafetyGuidelines struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	CommunityID   *uuid.UUID             `json:"communityId" db:"community_id"` // null for global guidelines
	Title         string                 `json:"title" db:"title"`
	Description   string                 `json:"description" db:"description"`
	Rules         []string               `json:"rules" db:"rules"`
	Reporting     map[string]interface{} `json:"reporting" db:"reporting"`
	Enforcement   map[string]interface{} `json:"enforcement" db:"enforcement"`
	Version       string                 `json:"version" db:"version"`
	IsActive      bool                   `json:"isActive" db:"is_active"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
}

// CreateReportRequest represents a request to create a report
type CreateReportRequest struct {
	TargetID     uuid.UUID              `json:"targetId" binding:"required"`
	TargetType   string                 `json:"targetType" binding:"required"`
	Reason       string                 `json:"reason" binding:"required"`
	Details      string                 `json:"details,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateReportRequest represents a request to update a report status
type UpdateReportRequest struct {
	Status      string  `json:"status" binding:"required"`
	ReviewNotes *string `json:"reviewNotes,omitempty"`
	Resolution  *string `json:"resolution,omitempty"`
}

// CreateModerationActionRequest represents a request to create a moderation action
type CreateModerationActionRequest struct {
	TargetID   uuid.UUID              `json:"targetId" binding:"required"`
	TargetType string                 `json:"targetType" binding:"required"`
	Action     string                 `json:"action" binding:"required"`
	Reason     string                 `json:"reason" binding:"required"`
	Duration   *int                   `json:"duration,omitempty"` // duration in hours
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ModerationStats represents moderation statistics
type ModerationStats struct {
	Period          string `json:"period"`
	TotalReports    int    `json:"totalReports"`
	PendingReports  int    `json:"pendingReports"`
	ResolvedReports int    `json:"resolvedReports"`
	ActionsToday    int    `json:"actionsToday"`
	ActionsThisWeek int    `json:"actionsThisWeek"`
	TopReasons      []struct {
		Reason string `json:"reason"`
		Count  int    `json:"count"`
	} `json:"topReasons"`
	ResolutionBreakdown map[string]int `json:"resolutionBreakdown"`
}