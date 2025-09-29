package models

import (
	"time"

	"gorm.io/gorm"
)

// TabNavigationState represents user's current navigation state in the Stream tab
type TabNavigationState struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            string    `json:"userId" gorm:"not null;index;size:255"`
	CurrentCategoryID string    `json:"currentCategoryId" gorm:"size:255"`
	CurrentSubtabID   *string   `json:"currentSubtabId" gorm:"size:255"`
	LastVisitedAt     time.Time `json:"lastVisitedAt" gorm:"not null"`
	SessionID         string    `json:"sessionId" gorm:"not null;size:255;index"`
	DevicePlatform    string    `json:"devicePlatform" gorm:"size:50;default:'web'"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// Navigation preferences
	AutoplayEnabled   bool `json:"autoplayEnabled" gorm:"default:true"`
	ShowSubtabs       bool `json:"showSubtabs" gorm:"default:true"`
	PreferredViewMode string `json:"preferredViewMode" gorm:"size:20;default:'grid'"` // grid, list
}

// StreamUserSession represents a user's session for Stream content consumption
type StreamUserSession struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          string    `json:"userId" gorm:"not null;index;size:255"`
	SessionToken    string    `json:"sessionToken" gorm:"not null;uniqueIndex;size:255"`
	StartedAt       time.Time `json:"startedAt" gorm:"not null"`
	LastActivityAt  time.Time `json:"lastActivityAt" gorm:"not null"`
	DeviceInfo      string    `json:"deviceInfo" gorm:"size:500"`
	IPAddress       string    `json:"ipAddress" gorm:"size:45"`
	UserAgent       string    `json:"userAgent" gorm:"size:1000"`
	IsActive        bool      `json:"isActive" gorm:"default:true"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Session statistics
	ContentViewCount    int           `json:"contentViewCount" gorm:"default:0"`
	TotalTimeSpent      time.Duration `json:"totalTimeSpent" gorm:"default:0"` // in seconds
	CategoriesVisited   int           `json:"categoriesVisited" gorm:"default:0"`
	LastCategoryVisited string        `json:"lastCategoryVisited" gorm:"size:255"`
}

// StreamContentView represents a user's interaction with specific content
type StreamContentView struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"userId" gorm:"not null;index;size:255"`
	ContentID     string    `json:"contentId" gorm:"not null;index;size:255"`
	SessionID     string    `json:"sessionId" gorm:"not null;size:255;index"`
	ViewStartedAt time.Time `json:"viewStartedAt" gorm:"not null"`
	ViewEndedAt   *time.Time `json:"viewEndedAt"`
	Duration      time.Duration `json:"duration" gorm:"default:0"` // in seconds
	ViewProgress  float64   `json:"viewProgress" gorm:"default:0"` // 0.0 to 1.0
	IsCompleted   bool      `json:"isCompleted" gorm:"default:false"`
	DevicePlatform string   `json:"devicePlatform" gorm:"size:50;default:'web'"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	// Relationship
	Content StreamContentItem `json:"content" gorm:"foreignKey:ContentID;references:ID"`
}

// StreamUserPreference represents user's preferences for Stream content
type StreamUserPreference struct {
	ID                    uint   `json:"id" gorm:"primaryKey"`
	UserID                string `json:"userId" gorm:"not null;uniqueIndex;size:255"`
	PreferredCategories   string `json:"preferredCategories" gorm:"type:text"` // JSON array of category IDs
	BlockedCategories     string `json:"blockedCategories" gorm:"type:text"`   // JSON array of category IDs
	AutoplayEnabled       bool   `json:"autoplayEnabled" gorm:"default:true"`
	HighQualityPreferred  bool   `json:"highQualityPreferred" gorm:"default:false"`
	OfflineDownloadEnabled bool  `json:"offlineDownloadEnabled" gorm:"default:false"`
	NotificationsEnabled  bool   `json:"notificationsEnabled" gorm:"default:true"`
	LanguagePreference    string `json:"languagePreference" gorm:"size:10;default:'en'"`
	RegionPreference      string `json:"regionPreference" gorm:"size:10;default:'US'"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// BeforeCreate hook for TabNavigationState
func (t *TabNavigationState) BeforeCreate(tx *gorm.DB) error {
	if t.LastVisitedAt.IsZero() {
		t.LastVisitedAt = time.Now()
	}
	return nil
}

// BeforeCreate hook for StreamUserSession
func (s *StreamUserSession) BeforeCreate(tx *gorm.DB) error {
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now()
	}
	if s.LastActivityAt.IsZero() {
		s.LastActivityAt = time.Now()
	}
	return nil
}

// BeforeCreate hook for StreamContentView
func (v *StreamContentView) BeforeCreate(tx *gorm.DB) error {
	if v.ViewStartedAt.IsZero() {
		v.ViewStartedAt = time.Now()
	}
	return nil
}

// UpdateActivity updates the last activity timestamp
func (s *StreamUserSession) UpdateActivity() {
	s.LastActivityAt = time.Now()
}

// IsExpired checks if the session has expired
func (s *StreamUserSession) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

// CalculateViewDuration calculates the duration of content view
func (v *StreamContentView) CalculateViewDuration() {
	if v.ViewEndedAt != nil {
		v.Duration = v.ViewEndedAt.Sub(v.ViewStartedAt)
	}
}

// MarkCompleted marks the content view as completed
func (v *StreamContentView) MarkCompleted() {
	v.IsCompleted = true
	v.ViewProgress = 1.0
	if v.ViewEndedAt == nil {
		now := time.Now()
		v.ViewEndedAt = &now
	}
	v.CalculateViewDuration()
}

// UpdateProgress updates the view progress
func (v *StreamContentView) UpdateProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	v.ViewProgress = progress

	// Mark as completed if progress reaches 95% or higher
	if progress >= 0.95 {
		v.MarkCompleted()
	}
}

// UpdateActivity updates the last activity timestamp for the view
func (v *StreamContentView) UpdateActivity() {
	// Update view timing and activity
	v.CalculateViewDuration()
}

// GetActiveSessionsCount returns count of active sessions for a user
func GetActiveSessionsCount(db *gorm.DB, userID string) (int64, error) {
	var count int64
	err := db.Model(&StreamUserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&count).Error
	return count, err
}

// GetUserNavigationState retrieves the current navigation state for a user
func GetUserNavigationState(db *gorm.DB, userID string) (*TabNavigationState, error) {
	var state TabNavigationState
	err := db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// GetOrCreateUserNavigationState gets existing or creates new navigation state
func GetOrCreateUserNavigationState(db *gorm.DB, userID, sessionID string) (*TabNavigationState, error) {
	var state TabNavigationState

	// Try to find existing state
	err := db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		First(&state).Error

	if err == gorm.ErrRecordNotFound {
		// Create new state with default values
		state = TabNavigationState{
			UserID:            userID,
			CurrentCategoryID: "books", // Default to first category
			SessionID:         sessionID,
			LastVisitedAt:     time.Now(),
			AutoplayEnabled:   true,
			ShowSubtabs:       true,
			PreferredViewMode: "grid",
		}

		err = db.Create(&state).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update session ID for existing state
		state.SessionID = sessionID
		state.LastVisitedAt = time.Now()
		err = db.Save(&state).Error
		if err != nil {
			return nil, err
		}
	}

	return &state, nil
}

// GetUserPreferences retrieves user preferences
func GetUserPreferences(db *gorm.DB, userID string) (*StreamUserPreference, error) {
	var pref StreamUserPreference
	err := db.Where("user_id = ?", userID).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		// Create default preferences
		pref = StreamUserPreference{
			UserID:                 userID,
			AutoplayEnabled:        true,
			HighQualityPreferred:   false,
			OfflineDownloadEnabled: false,
			NotificationsEnabled:   true,
			LanguagePreference:     "en",
			RegionPreference:       "US",
		}
		err = db.Create(&pref).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &pref, nil
}