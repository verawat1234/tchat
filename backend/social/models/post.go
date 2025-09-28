package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a social post
type Post struct {
	ID          uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID    uuid.UUID              `json:"authorId" db:"author_id" gorm:"type:uuid;not null;index"`
	CommunityID *uuid.UUID             `json:"communityId" db:"community_id" gorm:"type:uuid;index"`
	Content     string                 `json:"content" db:"content" gorm:"type:text;not null"`
	Type        string                 `json:"type" db:"type" gorm:"size:20;not null;index"` // text, image, video, link, poll
	Metadata    map[string]interface{} `json:"metadata" db:"metadata" gorm:"type:jsonb"`
	Tags        []string               `json:"tags" db:"tags" gorm:"type:text[]"`
	Visibility  string                 `json:"visibility" db:"visibility" gorm:"size:20;not null;index"` // public, members, private, followers
	MediaURLs   []string               `json:"mediaUrls" db:"media_urls" gorm:"type:text[]"`
	LinkPreview map[string]interface{} `json:"linkPreview" db:"link_preview" gorm:"type:jsonb"`

	// Interaction counts
	LikesCount     int `json:"likesCount" db:"likes_count" gorm:"default:0"`
	CommentsCount  int `json:"commentsCount" db:"comments_count" gorm:"default:0"`
	SharesCount    int `json:"sharesCount" db:"shares_count" gorm:"default:0"`
	ReactionsCount int `json:"reactionsCount" db:"reactions_count" gorm:"default:0"`
	ViewsCount     int `json:"viewsCount" db:"views_count" gorm:"default:0"`

	// Status flags
	IsEdited      bool `json:"isEdited" db:"is_edited" gorm:"default:false"`
	IsPinned      bool `json:"isPinned" db:"is_pinned" gorm:"default:false;index"`
	IsDeleted     bool `json:"isDeleted" db:"is_deleted" gorm:"default:false;index"`
	IsTrending    bool `json:"isTrending" db:"is_trending" gorm:"default:false;index"`

	// Timestamps
	CreatedAt time.Time  `json:"createdAt" db:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at" gorm:"index"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PostID    uuid.UUID              `json:"postId" db:"post_id" gorm:"type:uuid;not null;index"`
	AuthorID  uuid.UUID              `json:"authorId" db:"author_id" gorm:"type:uuid;not null;index"`
	ParentID  *uuid.UUID             `json:"parentId" db:"parent_id" gorm:"type:uuid;index"` // for replies
	Content   string                 `json:"content" db:"content" gorm:"type:text;not null"`
	Metadata  map[string]interface{} `json:"metadata" db:"metadata" gorm:"type:jsonb"`

	// Interaction counts
	LikesCount     int `json:"likesCount" db:"likes_count" gorm:"default:0"`
	RepliesCount   int `json:"repliesCount" db:"replies_count" gorm:"default:0"`
	ReactionsCount int `json:"reactionsCount" db:"reactions_count" gorm:"default:0"`

	// Status flags
	IsEdited  bool `json:"isEdited" db:"is_edited" gorm:"default:false"`
	IsDeleted bool `json:"isDeleted" db:"is_deleted" gorm:"default:false;index"`

	// Timestamps
	CreatedAt time.Time  `json:"createdAt" db:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at" gorm:"index"`
}

// Reaction represents a reaction to a post or comment
type Reaction struct {
	ID         uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `json:"userId" db:"user_id" gorm:"type:uuid;not null;index"`
	TargetID   uuid.UUID `json:"targetId" db:"target_id" gorm:"type:uuid;not null;index"` // post or comment ID
	TargetType string    `json:"targetType" db:"target_type" gorm:"size:20;not null;index"` // post, comment
	Type       string    `json:"type" db:"type" gorm:"size:20;not null;index"` // like, love, laugh, angry, sad, wow
	CreatedAt  time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
}

// SocialFeed represents a user's personalized social feed
type SocialFeed struct {
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Posts     []Post    `json:"posts"`
	Algorithm string    `json:"algorithm"` // chronological, personalized, trending
	Region    string    `json:"region"`
	Cursor    string    `json:"cursor"`
	HasMore   bool      `json:"hasMore"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TrendingContent represents trending posts and topics
type TrendingContent struct {
	Region    string                 `json:"region"`
	Timeframe string                 `json:"timeframe"` // 1h, 24h, 7d
	Topics    []string               `json:"topics"`
	Posts     []Post                 `json:"posts"`
	Hashtags  []string               `json:"hashtags"`
	Metrics   map[string]interface{} `json:"metrics"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// CreatePostRequest represents a request to create a post
// KMP Compatible: Validation and proper optional field handling
type CreatePostRequest struct {
	CommunityID *uuid.UUID             `json:"communityId,omitempty"`
	Content     string                 `json:"content" binding:"required" validate:"required,min=1,max=2000"`
	Type        string                 `json:"type" binding:"required" validate:"required,oneof=text image video link poll"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
	Visibility  string                 `json:"visibility" binding:"required" validate:"required,oneof=public members private followers"`
	MediaURLs   []string               `json:"mediaUrls,omitempty" validate:"omitempty,dive,url,max=500"`
}

// UpdatePostRequest represents a request to update a post
// KMP Compatible: Explicit nullable fields for partial updates
type UpdatePostRequest struct {
	Content   *string                 `json:"content,omitempty" validate:"omitempty,min=1,max=2000"`
	Tags      []string                `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
	Metadata  *map[string]interface{} `json:"metadata,omitempty"`
	IsPinned  *bool                   `json:"isPinned,omitempty"`
}

// CreateCommentRequest represents a request to create a comment
// KMP Compatible: Validation for mobile input handling
type CreateCommentRequest struct {
	PostID   uuid.UUID              `json:"postId" binding:"required" validate:"required"`
	Content  string                 `json:"content" binding:"required" validate:"required,min=1,max=1000"`
	ParentID *uuid.UUID             `json:"parentId,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreateReactionRequest represents a request to add a reaction
// KMP Compatible: Strict validation for reaction types
type CreateReactionRequest struct {
	TargetID   uuid.UUID `json:"targetId" binding:"required" validate:"required"`
	TargetType string    `json:"targetType" binding:"required" validate:"required,oneof=post comment"`
	Type       string    `json:"type" binding:"required" validate:"required,oneof=like love laugh angry sad wow"`
}

// SocialFeedRequest represents parameters for getting social feed
// KMP Compatible: Sensible defaults and validation for mobile pagination
type SocialFeedRequest struct {
	Algorithm string `json:"algorithm,omitempty" form:"algorithm" validate:"omitempty,oneof=chronological personalized trending" default:"chronological"`
	Limit     int    `json:"limit" form:"limit" validate:"min=1,max=50" default:"20"`
	Cursor    string `json:"cursor,omitempty" form:"cursor"`
	Region    string `json:"region,omitempty" form:"region" validate:"omitempty,oneof=TH SG ID MY PH VN"`
}

// TrendingRequest represents parameters for getting trending content
// KMP Compatible: Validation and mobile-friendly defaults
type TrendingRequest struct {
	Region    string `json:"region,omitempty" form:"region" validate:"omitempty,oneof=TH SG ID MY PH VN"`
	Timeframe string `json:"timeframe,omitempty" form:"timeframe" validate:"omitempty,oneof=1h 24h 7d" default:"24h"`
	Category  string `json:"category,omitempty" form:"category" validate:"omitempty,max=50"`
	Limit     int    `json:"limit" form:"limit" validate:"min=1,max=50" default:"20"`
}

// ShareRequest represents a request to share content
type ShareRequest struct {
	ContentID   uuid.UUID              `json:"contentId" binding:"required"`
	ContentType string                 `json:"contentType" binding:"required"` // post, comment, community
	Platform    string                 `json:"platform" binding:"required"` // internal, external
	Message     string                 `json:"message,omitempty"`
	Privacy     string                 `json:"privacy,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Share represents a share action
type Share struct {
	ID          uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID              `json:"userId" db:"user_id" gorm:"type:uuid;not null;index"`
	ContentID   uuid.UUID              `json:"contentId" db:"content_id" gorm:"type:uuid;not null;index"`
	ContentType string                 `json:"contentType" db:"content_type" gorm:"size:20;not null;index"`
	Platform    string                 `json:"platform" db:"platform" gorm:"size:20;not null"`
	Message     string                 `json:"message" db:"message" gorm:"size:500"`
	Privacy     string                 `json:"privacy" db:"privacy" gorm:"size:20;not null"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata" gorm:"type:jsonb"`
	Status      string                 `json:"status" db:"status" gorm:"size:20;not null;default:'active'"`
	CreatedAt   time.Time              `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
}