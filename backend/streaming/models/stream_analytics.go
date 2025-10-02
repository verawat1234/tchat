package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StreamAnalytics stores aggregated analytics and metrics for completed streams
// Calculated after stream ends for reporting and analytics dashboards
type StreamAnalytics struct {
	StreamID uuid.UUID `gorm:"type:uuid;primary_key" json:"stream_id"`

	// Viewer Metrics
	TotalUniqueViewers    int `gorm:"default:0" json:"total_unique_viewers"`
	PeakConcurrentViewers int `gorm:"default:0;index:idx_peak_concurrent_viewers" json:"peak_concurrent_viewers"`
	AverageWatchDurationSeconds int `gorm:"default:0" json:"average_watch_duration_seconds"`

	// Engagement Metrics
	TotalChatMessages int `gorm:"default:0" json:"total_chat_messages"`
	TotalReactions    int `gorm:"default:0" json:"total_reactions"`
	UniqueChatter     int `gorm:"default:0" json:"unique_chatter"`

	// Store Context Metrics (nullable for video streams)
	ProductsFeatured   *int     `gorm:"" json:"products_featured,omitempty"`
	TotalProductViews  *int     `gorm:"" json:"total_product_views,omitempty"`
	TotalProductClicks *int     `gorm:"" json:"total_product_clicks,omitempty"`
	TotalPurchases     *int     `gorm:"" json:"total_purchases,omitempty"`
	TotalRevenue       *float64 `gorm:"type:decimal(12,2);index:idx_total_revenue" json:"total_revenue,omitempty"`

	// Quality Metrics
	AverageViewerQuality sql.NullString `gorm:"type:varchar(20)" json:"average_viewer_quality,omitempty"`
	TotalRebufferEvents  int            `gorm:"default:0" json:"total_rebuffer_events"`

	// Geographic Distribution (stored as JSONB)
	// Example: {"US": 150, "SG": 80, "TH": 45}
	ViewerCountries []byte `gorm:"type:jsonb" json:"viewer_countries,omitempty"`

	// Calculation Timestamp
	CalculatedAt *time.Time `gorm:"" json:"calculated_at,omitempty"`
}

// TableName specifies the table name for StreamAnalytics model
func (StreamAnalytics) TableName() string {
	return "stream_analytics"
}

// BeforeCreate hook to set CalculatedAt timestamp
func (sa *StreamAnalytics) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if sa.CalculatedAt == nil {
		sa.CalculatedAt = &now
	}
	return nil
}