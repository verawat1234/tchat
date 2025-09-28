package models

import (
	"time"

	"github.com/google/uuid"
)

// Community represents a social community
type Community struct {
	ID          uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string                 `json:"name" db:"name" gorm:"size:100;not null;index"`
	Description string                 `json:"description" db:"description" gorm:"size:500"`
	Type        string                 `json:"type" db:"type" gorm:"size:20;not null;index"` // public, private, restricted
	Category    string                 `json:"category" db:"category" gorm:"size:50;index"`
	Region      string                 `json:"region" db:"region" gorm:"size:10;index"`
	Avatar      string                 `json:"avatar" db:"avatar" gorm:"size:255"`
	Banner      string                 `json:"banner" db:"banner" gorm:"size:255"`
	Tags        []string               `json:"tags" db:"tags" gorm:"type:text[]"`
	Rules       []string               `json:"rules" db:"rules" gorm:"type:text[]"`
	CreatorID   uuid.UUID              `json:"creatorId" db:"creator_id" gorm:"type:uuid;not null;index"`
	Settings    map[string]interface{} `json:"settings" db:"settings" gorm:"type:jsonb"`
	MembersCount     int               `json:"membersCount" db:"members_count" gorm:"default:0"`
	PostsCount       int               `json:"postsCount" db:"posts_count" gorm:"default:0"`
	IsVerified       bool              `json:"isVerified" db:"is_verified" gorm:"default:false;index"`
	CreatedAt        time.Time         `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
}

// CommunityMember represents user membership in a community
type CommunityMember struct {
	ID          uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CommunityID uuid.UUID `json:"communityId" db:"community_id" gorm:"type:uuid;not null;index"`
	UserID      uuid.UUID `json:"userId" db:"user_id" gorm:"type:uuid;not null;index"`
	Role        string    `json:"role" db:"role" gorm:"size:20;not null;default:'member'"` // owner, moderator, member
	Status      string    `json:"status" db:"status" gorm:"size:20;not null;default:'active'"` // active, pending, banned
	JoinReason  string    `json:"joinReason" db:"join_reason" gorm:"size:255"`
	Source      string    `json:"source" db:"source" gorm:"size:20;not null"` // discovery, invitation, search
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
}

// CommunityAnalytics represents community analytics data
type CommunityAnalytics struct {
	CommunityID      uuid.UUID              `json:"communityId" db:"community_id"`
	Period           string                 `json:"period" db:"period"`
	Engagement       map[string]interface{} `json:"engagement" db:"engagement"`
	Growth           map[string]interface{} `json:"growth" db:"growth"`
	Demographics     map[string]interface{} `json:"demographics" db:"demographics"`
	TopTopics        []string               `json:"topTopics" db:"top_topics"`
	SentimentAnalysis map[string]interface{} `json:"sentimentAnalysis" db:"sentiment_analysis"`
	RegionBreakdown  map[string]interface{} `json:"regionBreakdown" db:"region_breakdown"`
	UpdatedAt        time.Time              `json:"updatedAt" db:"updated_at"`
}

// CommunityInsights represents AI-powered community insights
type CommunityInsights struct {
	CommunityID       uuid.UUID              `json:"communityId" db:"community_id"`
	TopTopics         []string               `json:"topTopics" db:"top_topics"`
	SentimentAnalysis map[string]interface{} `json:"sentimentAnalysis" db:"sentiment_analysis"`
	Recommendations   []string               `json:"recommendations" db:"recommendations"`
	RegionBreakdown   map[string]interface{} `json:"regionBreakdown" db:"region_breakdown"`
	TrendingHashtags  []string               `json:"trendingHashtags" db:"trending_hashtags"`
	EngagementPatterns map[string]interface{} `json:"engagementPatterns" db:"engagement_patterns"`
	UpdatedAt         time.Time              `json:"updatedAt" db:"updated_at"`
}

// CreateCommunityRequest represents a request to create a community
// KMP Compatible: Validation and proper field constraints
type CreateCommunityRequest struct {
	Name        string                 `json:"name" binding:"required" validate:"required,min=2,max=100"`
	Description string                 `json:"description" binding:"required" validate:"required,min=10,max=500"`
	Type        string                 `json:"type" binding:"required" validate:"required,oneof=public private restricted"`
	Category    string                 `json:"category" binding:"required" validate:"required,max=50"`
	Region      string                 `json:"region" binding:"required" validate:"required,oneof=TH SG ID MY PH VN"`
	Tags        []string               `json:"tags,omitempty" validate:"omitempty,dive,max=30"`
	Rules       []string               `json:"rules,omitempty" validate:"omitempty,dive,max=200"`
	Avatar      string                 `json:"avatar,omitempty" validate:"omitempty,url,max=500"`
	Banner      string                 `json:"banner,omitempty" validate:"omitempty,url,max=500"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// UpdateCommunityRequest represents a request to update community
// KMP Compatible: Nullable fields for partial updates with validation
type UpdateCommunityRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Avatar      *string `json:"avatar,omitempty" validate:"omitempty,url,max=500"`
	Banner      *string `json:"banner,omitempty" validate:"omitempty,url,max=500"`
	IsPinned    *bool   `json:"isPinned,omitempty"`
}

// JoinCommunityRequest represents a request to join a community
// KMP Compatible: Validation for join sources
type JoinCommunityRequest struct {
	Reason string `json:"reason,omitempty" validate:"omitempty,max=200"`
	Source string `json:"source" binding:"required" validate:"required,oneof=discovery invitation search manual"`
}

// CommunityRoleRequest represents a request to assign/update member role
// KMP Compatible: Role validation for community management
type CommunityRoleRequest struct {
	UserID uuid.UUID `json:"userId" binding:"required" validate:"required"`
	Role   string    `json:"role" binding:"required" validate:"required,oneof=owner moderator member"`
	Reason string    `json:"reason,omitempty" validate:"omitempty,max=200"`
}

// CommunityDiscoveryRequest represents parameters for discovering communities
// KMP Compatible: Mobile-friendly pagination and validation
type CommunityDiscoveryRequest struct {
	Category string   `json:"category,omitempty" form:"category" validate:"omitempty,max=50"`
	Region   string   `json:"region,omitempty" form:"region" validate:"omitempty,oneof=TH SG ID MY PH VN"`
	Type     string   `json:"type,omitempty" form:"type" validate:"omitempty,oneof=public private restricted"`
	Tags     []string `json:"tags,omitempty" form:"tags" validate:"omitempty,dive,max=30"`
	Limit    int      `json:"limit" form:"limit" validate:"min=1,max=50" default:"20"`
	Offset   int      `json:"offset" form:"offset" validate:"min=0" default:"0"`
}