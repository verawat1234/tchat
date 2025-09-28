package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ReviewType represents the type of review
type ReviewType string

const (
	ReviewTypeProduct  ReviewType = "product"
	ReviewTypeBusiness ReviewType = "business"
	ReviewTypeOrder    ReviewType = "order"
)

// ReviewStatus represents the status of a review
type ReviewStatus string

const (
	ReviewStatusPending   ReviewStatus = "pending"
	ReviewStatusApproved  ReviewStatus = "approved"
	ReviewStatusRejected  ReviewStatus = "rejected"
	ReviewStatusFlagged   ReviewStatus = "flagged"
)

// ReviewImage represents an image in a review
type ReviewImage struct {
	URL       string `json:"url" gorm:"column:url;size:500;not null"`
	Caption   string `json:"caption,omitempty" gorm:"column:caption;size:200"`
	SortOrder int    `json:"sort_order" gorm:"column:sort_order;default:0"`
}

// ReviewReply represents a response to a review
type ReviewReply struct {
	ID          uuid.UUID `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null"`
	UserType    string    `json:"user_type" gorm:"column:user_type;size:20;not null"` // business_owner, admin
	Content     string    `json:"content" gorm:"column:content;size:1000;not null"`
	IsOfficial  bool      `json:"is_official" gorm:"column:is_official;default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// Review represents a product/business review
type Review struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	Type       ReviewType `json:"type" gorm:"column:type;type:varchar(20);not null"`
	Status     ReviewStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`

	// What is being reviewed
	ProductID  *uuid.UUID `json:"product_id,omitempty" gorm:"column:product_id;type:uuid"`
	BusinessID *uuid.UUID `json:"business_id,omitempty" gorm:"column:business_id;type:uuid"`
	OrderID    *uuid.UUID `json:"order_id,omitempty" gorm:"column:order_id;type:uuid"`

	// Who wrote the review
	UserID     uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null"`
	UserName   string    `json:"user_name" gorm:"column:user_name;size:100;not null"`
	UserEmail  string    `json:"user_email" gorm:"column:user_email;size:255;not null"`
	IsVerified bool      `json:"is_verified" gorm:"column:is_verified;default:false"`

	// Review content
	Rating       decimal.Decimal `json:"rating" gorm:"column:rating;type:decimal(2,1);not null"`
	Title        string          `json:"title" gorm:"column:title;size:200"`
	Content      string          `json:"content" gorm:"column:content;size:2000;not null"`
	Images       []ReviewImage   `json:"images,omitempty" gorm:"column:images;type:jsonb"`

	// Review metrics
	HelpfulCount    int `json:"helpful_count" gorm:"column:helpful_count;default:0"`
	NotHelpfulCount int `json:"not_helpful_count" gorm:"column:not_helpful_count;default:0"`
	ReportCount     int `json:"report_count" gorm:"column:report_count;default:0"`

	// Response from business
	Response *ReviewResponse `json:"response,omitempty" gorm:"column:response;type:jsonb"`

	// Moderation
	ModerationNotes string     `json:"moderation_notes,omitempty" gorm:"column:moderation_notes;size:500"`
	ModeratedBy     *uuid.UUID `json:"moderated_by,omitempty" gorm:"column:moderated_by;type:uuid"`
	ModeratedAt     *time.Time `json:"moderated_at,omitempty" gorm:"column:moderated_at"`

	// Regional compliance
	DataRegion string `json:"data_region" gorm:"column:data_region;size:20"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the Review model
func (Review) TableName() string {
	return "reviews"
}

// ReviewHelpful represents a helpful vote on a review
type ReviewHelpful struct {
	ID       uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	ReviewID uuid.UUID `json:"review_id" gorm:"column:review_id;type:uuid;not null;index"`
	UserID   uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index"`
	IsHelpful bool     `json:"is_helpful" gorm:"column:is_helpful;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`

	// Unique constraint on review_id + user_id
}

// TableName returns the table name for the ReviewHelpful model
func (ReviewHelpful) TableName() string {
	return "review_helpful"
}

// ReviewReport represents a report on a review
type ReviewReport struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	ReviewID  uuid.UUID `json:"review_id" gorm:"column:review_id;type:uuid;not null;index"`
	UserID    uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index"`
	Reason    string    `json:"reason" gorm:"column:reason;size:50;not null"`
	Comment   string    `json:"comment,omitempty" gorm:"column:comment;size:500"`
	Status    string    `json:"status" gorm:"column:status;size:20;not null;default:'pending'"` // pending, reviewed, resolved
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
}

// TableName returns the table name for the ReviewReport model
func (ReviewReport) TableName() string {
	return "review_reports"
}