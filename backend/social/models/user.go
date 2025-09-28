package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	sharedModels "tchat.dev/shared/models"
)

// SocialProfile extends the shared User model with social-specific fields
// KMP Compatible: All fields use explicit JSON tags and nullable types where appropriate
type SocialProfile struct {
	// Embed the shared User model
	sharedModels.User

	// Social-specific fields with explicit nullability for KMP
	Interests        []string               `json:"interests,omitempty" db:"interests" gorm:"type:text[]"`
	SocialLinks      map[string]interface{} `json:"socialLinks,omitempty" db:"social_links" gorm:"type:jsonb"`
	SocialPreferences map[string]interface{} `json:"socialPreferences,omitempty" db:"social_preferences" gorm:"type:jsonb"`

	// Social metrics - always present with default values for mobile consistency
	FollowersCount int `json:"followersCount" db:"followers_count" gorm:"default:0"`
	FollowingCount int `json:"followingCount" db:"following_count" gorm:"default:0"`
	PostsCount     int `json:"postsCount" db:"posts_count" gorm:"default:0"`

	// Social verification (separate from KYC verification)
	IsSocialVerified bool `json:"isSocialVerified" db:"is_social_verified" gorm:"default:false"`

	// Social-specific timestamps - RFC3339 format for KMP compatibility
	SocialCreatedAt time.Time `json:"socialCreatedAt" db:"social_created_at" gorm:"autoCreateTime"`
	SocialUpdatedAt time.Time `json:"socialUpdatedAt" db:"social_updated_at" gorm:"autoUpdateTime"`
}

// MarshalJSON customizes JSON serialization for SocialProfile to include all fields
func (sp *SocialProfile) MarshalJSON() ([]byte, error) {
	// Create a temporary struct that combines User fields and SocialProfile fields
	return json.Marshal(&struct {
		// User fields - inline all the fields from the embedded User struct
		ID          uuid.UUID  `json:"id"`
		Username    string     `json:"username,omitempty"`
		Phone       string     `json:"phone,omitempty"`
		PhoneNumber string     `json:"phone_number,omitempty"`
		Email       string     `json:"email,omitempty"`
		Name        string     `json:"name"`
		DisplayName string     `json:"display_name,omitempty"`
		FirstName   string     `json:"first_name,omitempty"`
		LastName    string     `json:"last_name,omitempty"`
		Avatar      string     `json:"avatar,omitempty"`
		Country     string     `json:"country"`
		CountryCode string     `json:"country_code,omitempty"`
		Locale      string     `json:"locale"`
		Language    string     `json:"language,omitempty"`
		Timezone    string     `json:"timezone,omitempty"`
		TimezoneAlias string   `json:"timezone_alias,omitempty"`
		Bio         string     `json:"bio,omitempty"`
		DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
		Gender      string     `json:"gender,omitempty"`
		Active      bool       `json:"is_active"`
		KYCStatus   string     `json:"kyc_status"`
		KYCTier     int        `json:"kyc_tier"`
		Status      string     `json:"status"`
		LastSeen    *time.Time `json:"last_seen,omitempty"`
		LastActiveAt *time.Time `json:"last_active_at,omitempty"`
		Verified    bool       `json:"is_verified"`
		EmailVerified bool     `json:"is_email_verified"`
		PhoneVerified bool     `json:"is_phone_verified"`
		PrefTheme   string     `json:"pref_theme,omitempty"`
		PrefLanguage string    `json:"pref_language,omitempty"`
		PrefNotificationsEmail bool `json:"pref_notifications_email"`
		PrefNotificationsPush bool  `json:"pref_notifications_push"`
		PrefPrivacyLevel string     `json:"pref_privacy_level,omitempty"`
		Metadata    json.RawMessage `json:"metadata,omitempty"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`

		// SocialProfile specific fields
		Interests         []string               `json:"interests,omitempty"`
		SocialLinks       map[string]interface{} `json:"socialLinks,omitempty"`
		SocialPreferences map[string]interface{} `json:"socialPreferences,omitempty"`
		FollowersCount    int                    `json:"followersCount"`
		FollowingCount    int                    `json:"followingCount"`
		PostsCount        int                    `json:"postsCount"`
		IsSocialVerified  bool                   `json:"isSocialVerified"`
		SocialCreatedAt   time.Time              `json:"socialCreatedAt"`
		SocialUpdatedAt   time.Time              `json:"socialUpdatedAt"`
	}{
		// User fields
		ID:          sp.User.ID,
		Username:    sp.User.Username,
		Phone:       sp.User.Phone,
		PhoneNumber: sp.User.GetFullPhoneNumber(), // Use the formatted phone number
		Email:       sp.User.Email,
		Name:        sp.User.Name,
		DisplayName: sp.User.DisplayName,
		FirstName:   sp.User.FirstName,
		LastName:    sp.User.LastName,
		Avatar:      sp.User.Avatar,
		Country:     sp.User.Country,
		CountryCode: sp.User.CountryCode,
		Locale:      sp.User.Locale,
		Language:    sp.User.Language,
		Timezone:    sp.User.Timezone,
		TimezoneAlias: sp.User.TimezoneAlias,
		Bio:         sp.User.Bio,
		DateOfBirth: sp.User.DateOfBirth,
		Gender:      sp.User.Gender,
		Active:      sp.User.Active,
		KYCStatus:   sp.User.KYCStatus,
		KYCTier:     sp.User.KYCTier,
		Status:      sp.User.Status,
		LastSeen:    sp.User.LastSeen,
		LastActiveAt: sp.User.LastActiveAt,
		Verified:    sp.User.Verified,
		EmailVerified: sp.User.EmailVerified,
		PhoneVerified: sp.User.PhoneVerified,
		PrefTheme:   sp.User.PrefTheme,
		PrefLanguage: sp.User.PrefLanguage,
		PrefNotificationsEmail: sp.User.PrefNotificationsEmail,
		PrefNotificationsPush:  sp.User.PrefNotificationsPush,
		PrefPrivacyLevel: sp.User.PrefPrivacyLevel,
		Metadata:    sp.User.Metadata,
		CreatedAt:   sp.User.CreatedAt,
		UpdatedAt:   sp.User.UpdatedAt,

		// SocialProfile fields
		Interests:         sp.Interests,
		SocialLinks:       sp.SocialLinks,
		SocialPreferences: sp.SocialPreferences,
		FollowersCount:    sp.FollowersCount,
		FollowingCount:    sp.FollowingCount,
		PostsCount:        sp.PostsCount,
		IsSocialVerified:  sp.IsSocialVerified,
		SocialCreatedAt:   sp.SocialCreatedAt,
		SocialUpdatedAt:   sp.SocialUpdatedAt,
	})
}

// UserRelationship represents following/follower relationships
// KMP Compatible: Simplified for mobile sync patterns
type UserRelationship struct {
	ID           uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FollowerID   uuid.UUID `json:"followerId" db:"follower_id" gorm:"type:uuid;not null;index"`
	FollowingID  uuid.UUID `json:"followingId" db:"following_id" gorm:"type:uuid;not null;index"`
	Status       string    `json:"status" db:"status" gorm:"size:20;not null;default:'active'"` // pending, active, blocked
	Source       string    `json:"source" db:"source" gorm:"size:30;not null"` // discovery, suggestion, search, follow_back
	IsMutual     bool      `json:"isMutual" db:"is_mutual" gorm:"default:false"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at" gorm:"autoUpdateTime"`
}

// UserActivity represents user social activity for analytics
// KMP Compatible: Explicit nullable handling and mobile platform support
type UserActivity struct {
	ID           uuid.UUID              `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID              `json:"userId" db:"user_id" gorm:"type:uuid;not null;index"`
	ActivityType string                 `json:"activityType" db:"activity_type" gorm:"size:30;not null;index"` // post, comment, reaction, follow, share
	TargetID     *uuid.UUID             `json:"targetId,omitempty" db:"target_id" gorm:"type:uuid;index"`
	TargetType   string                 `json:"targetType,omitempty" db:"target_type" gorm:"size:20;index"` // post, comment, user, community
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	Region       string                 `json:"region" db:"region" gorm:"size:10;not null;index"`
	Platform     string                 `json:"platform" db:"platform" gorm:"size:20;not null;index"` // web, mobile, api, kmp_android, kmp_ios
	CreatedAt    time.Time              `json:"createdAt" db:"created_at" gorm:"autoCreateTime;index"`
}

// UpdateSocialProfileRequest represents a request to update social profile
// KMP Compatible: Explicit nullable pointers for optional updates
type UpdateSocialProfileRequest struct {
	DisplayName       *string                 `json:"displayName,omitempty" validate:"omitempty,max=100"`
	Bio               *string                 `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar            *string                 `json:"avatar,omitempty" validate:"omitempty,url,max=500"`
	Interests         []string                `json:"interests,omitempty" validate:"omitempty,dive,max=50"`
	SocialLinks       *map[string]interface{} `json:"socialLinks,omitempty"`
	SocialPreferences *map[string]interface{} `json:"socialPreferences,omitempty"`
}

// FollowRequest represents a request to follow/unfollow a user
// KMP Compatible: String UUIDs for easier mobile handling
type FollowRequest struct {
	FollowerID   uuid.UUID `json:"followerId" binding:"required" validate:"required"`
	FollowingID  uuid.UUID `json:"followingId" binding:"required" validate:"required"`
	Source       string    `json:"source" binding:"required" validate:"required,oneof=discovery suggestion search follow_back manual"`
}

// UserDiscoveryRequest represents parameters for discovering users
// KMP Compatible: Validation and sensible defaults for mobile
type UserDiscoveryRequest struct {
	Region    string   `json:"region,omitempty" form:"region" validate:"omitempty,oneof=TH SG ID MY PH VN"`
	Interests []string `json:"interests,omitempty" form:"interests" validate:"omitempty,dive,max=50"`
	Limit     int      `json:"limit" form:"limit" validate:"min=1,max=100" default:"20"`
	Offset    int      `json:"offset" form:"offset" validate:"min=0" default:"0"`
}

// UserAnalyticsResponse represents user social analytics
// KMP Compatible: Structured data instead of generic maps for better type safety
type UserAnalyticsResponse struct {
	UserID     uuid.UUID              `json:"userId"`
	Period     string                 `json:"period" validate:"oneof=1h 24h 7d 30d"`
	Followers  map[string]interface{} `json:"followers"`
	Following  map[string]interface{} `json:"following"`
	Engagement map[string]interface{} `json:"engagement"`
	Reach      map[string]interface{} `json:"reach"`
	Growth     map[string]interface{} `json:"growth"`
	Demographics map[string]interface{} `json:"demographics"`
	UpdatedAt  time.Time              `json:"updatedAt"`

	// KMP specific fields for better mobile handling
	IsRealTime bool `json:"isRealTime" default:"false"`
	CacheExpiry *time.Time `json:"cacheExpiry,omitempty"`
}

// Follow represents a simplified user following relationship for the database
type Follow struct {
	ID          uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FollowerID  uuid.UUID `json:"followerId" db:"follower_id" gorm:"type:uuid;not null;index"`
	FollowingID uuid.UUID `json:"followingId" db:"following_id" gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`

	// Add unique constraint to prevent duplicate follows
	// GORM will handle this via unique index: follower_id, following_id
}