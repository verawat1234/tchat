package main

import (
	"log"

	"gorm.io/gorm"
	"tchat.dev/video/models"
)

// ModelManager handles database migrations and model management
type ModelManager struct {
	db *gorm.DB
}

// NewModelManager creates a new ModelManager instance
func NewModelManager(db *gorm.DB) *ModelManager {
	return &ModelManager{db: db}
}

// RunVideoMigrations runs all necessary migrations for video service models
func (m *ModelManager) RunVideoMigrations() error {
	log.Println("Running video service migrations...")

	// Auto-migrate all video service models
	err := m.db.AutoMigrate(
		&models.Channel{},
		&models.Video{},
		&models.VideoInteraction{},
		&models.VideoComment{},
		&models.VideoShare{},
	)

	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Video service migrations completed successfully")
	return nil
}

// CreateIndexes creates database indexes for better performance
func (m *ModelManager) CreateIndexes() error {
	log.Println("Creating database indexes for video service...")

	// Indexes for videos table
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_channel_id ON videos(channel_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_category ON videos(category)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_type ON videos(type)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_views ON videos(views)").Error; err != nil {
		return err
	}

	// Indexes for channels table
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_channels_user_id ON channels(user_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_channels_verified ON channels(verified)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_channels_subscribers ON channels(subscribers)").Error; err != nil {
		return err
	}

	// Indexes for video_interactions table
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_interactions_video_id ON video_interactions(video_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_interactions_user_id ON video_interactions(user_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_interactions_type ON video_interactions(type)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_interactions_created_at ON video_interactions(created_at)").Error; err != nil {
		return err
	}

	// Indexes for video_comments table
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_comments_video_id ON video_comments(video_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_comments_user_id ON video_comments(user_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_comments_parent_id ON video_comments(parent_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_comments_created_at ON video_comments(created_at)").Error; err != nil {
		return err
	}

	// Indexes for video_shares table
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_shares_video_id ON video_shares(video_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_shares_user_id ON video_shares(user_id)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_shares_platform ON video_shares(platform)").Error; err != nil {
		return err
	}
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_video_shares_created_at ON video_shares(created_at)").Error; err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}