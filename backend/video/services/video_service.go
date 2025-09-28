package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/video/models"
)

// VideoRepository defines the interface for video data access
type VideoRepository interface {
	GetVideos(limit, offset int, category string) ([]models.Video, error)
	GetVideoByID(id uuid.UUID) (*models.Video, error)
	CreateVideo(video *models.Video) error
	UpdateVideo(video *models.Video) error
	DeleteVideo(id uuid.UUID) error
	GetVideoInteractions(videoID uuid.UUID) ([]models.VideoInteraction, error)
	CreateVideoInteraction(interaction *models.VideoInteraction) error
	GetChannelByID(id uuid.UUID) (*models.Channel, error)
	CreateChannel(channel *models.Channel) error
	UpdateChannel(channel *models.Channel) error

	// Video search and filtering
	SearchVideos(query string, limit, offset int, category string) ([]models.Video, error)
	GetVideosByCategory(category string, limit, offset int) ([]models.Video, error)
	GetTrendingVideos(timeframe string, limit, offset int) ([]models.Video, error)

	// Video comments
	GetVideoComments(videoID uuid.UUID, limit, offset int) ([]models.VideoComment, error)
	CreateVideoComment(comment *models.VideoComment) error

	// Video sharing
	CreateVideoShare(share *models.VideoShare) error
}

// PostgreSQLVideoRepository implements VideoRepository for PostgreSQL
type PostgreSQLVideoRepository struct {
	db *gorm.DB
}

// NewPostgreSQLVideoRepository creates a new PostgreSQL video repository
func NewPostgreSQLVideoRepository(db *gorm.DB) VideoRepository {
	return &PostgreSQLVideoRepository{db: db}
}

// GetVideos retrieves videos with pagination and optional category filter
func (r *PostgreSQLVideoRepository) GetVideos(limit, offset int, category string) ([]models.Video, error) {
	var videos []models.Video
	query := r.db.Preload("Channel")

	if category != "" {
		query = query.Where("category = ? AND status = ?", category, "active")
	} else {
		query = query.Where("status = ?", "active")
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

// GetVideoByID retrieves a video by its ID
func (r *PostgreSQLVideoRepository) GetVideoByID(id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.Preload("Channel").Where("id = ? AND status = ?", id, "active").First(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("video not found")
		}
		return nil, err
	}
	return &video, nil
}

// CreateVideo creates a new video
func (r *PostgreSQLVideoRepository) CreateVideo(video *models.Video) error {
	return r.db.Create(video).Error
}

// UpdateVideo updates an existing video
func (r *PostgreSQLVideoRepository) UpdateVideo(video *models.Video) error {
	return r.db.Save(video).Error
}

// DeleteVideo soft deletes a video by setting status to inactive
func (r *PostgreSQLVideoRepository) DeleteVideo(id uuid.UUID) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).Update("status", "inactive").Error
}

// GetVideoInteractions retrieves interactions for a video
func (r *PostgreSQLVideoRepository) GetVideoInteractions(videoID uuid.UUID) ([]models.VideoInteraction, error) {
	var interactions []models.VideoInteraction
	err := r.db.Where("video_id = ?", videoID).Find(&interactions).Error
	return interactions, err
}

// CreateVideoInteraction creates a new video interaction
func (r *PostgreSQLVideoRepository) CreateVideoInteraction(interaction *models.VideoInteraction) error {
	return r.db.Create(interaction).Error
}

// GetChannelByID retrieves a channel by its ID
func (r *PostgreSQLVideoRepository) GetChannelByID(id uuid.UUID) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.Where("id = ?", id).First(&channel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, err
	}
	return &channel, nil
}

// CreateChannel creates a new channel
func (r *PostgreSQLVideoRepository) CreateChannel(channel *models.Channel) error {
	return r.db.Create(channel).Error
}

// UpdateChannel updates an existing channel
func (r *PostgreSQLVideoRepository) UpdateChannel(channel *models.Channel) error {
	return r.db.Save(channel).Error
}

// SearchVideos searches videos by title, description, or tags
func (r *PostgreSQLVideoRepository) SearchVideos(query string, limit, offset int, category string) ([]models.Video, error) {
	var videos []models.Video
	dbQuery := r.db.Preload("Channel").Where("status = ?", "active")

	// Add search conditions
	searchCondition := "title ILIKE ? OR description ILIKE ? OR ? = ANY(tags)"
	queryParam := "%" + query + "%"
	dbQuery = dbQuery.Where(searchCondition, queryParam, queryParam, query)

	// Add category filter if provided
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	err := dbQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

// GetVideosByCategory retrieves videos by category
func (r *PostgreSQLVideoRepository) GetVideosByCategory(category string, limit, offset int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Preload("Channel").
		Where("category = ? AND status = ?", category, "active").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&videos).Error
	return videos, err
}

// GetTrendingVideos retrieves trending videos based on timeframe
func (r *PostgreSQLVideoRepository) GetTrendingVideos(timeframe string, limit, offset int) ([]models.Video, error) {
	var videos []models.Video

	// Calculate time threshold based on timeframe
	var threshold time.Time
	now := time.Now()
	switch timeframe {
	case "day":
		threshold = now.AddDate(0, 0, -1)
	case "week":
		threshold = now.AddDate(0, 0, -7)
	case "month":
		threshold = now.AddDate(0, -1, 0)
	default:
		threshold = now.AddDate(0, 0, -1) // Default to day
	}

	// Order by a combination of views, likes, and recency for trending algorithm
	err := r.db.Preload("Channel").
		Where("status = ? AND created_at >= ?", "active", threshold).
		Order("(views + likes * 2) DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&videos).Error
	return videos, err
}

// GetVideoComments retrieves comments for a video
func (r *PostgreSQLVideoRepository) GetVideoComments(videoID uuid.UUID, limit, offset int) ([]models.VideoComment, error) {
	var comments []models.VideoComment
	err := r.db.Where("video_id = ? AND parent_id IS NULL", videoID).
		Preload("Replies").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error
	return comments, err
}

// CreateVideoComment creates a new video comment
func (r *PostgreSQLVideoRepository) CreateVideoComment(comment *models.VideoComment) error {
	return r.db.Create(comment).Error
}

// CreateVideoShare creates a new video share record
func (r *PostgreSQLVideoRepository) CreateVideoShare(share *models.VideoShare) error {
	return r.db.Create(share).Error
}

// VideoService provides business logic for video operations
type VideoService struct {
	videoRepo VideoRepository
	db        *gorm.DB
}

// NewVideoService creates a new video service
func NewVideoService(videoRepo VideoRepository, db *gorm.DB) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
		db:        db,
	}
}

// GetVideos retrieves videos with pagination
func (s *VideoService) GetVideos(page, perPage int, category string) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	return s.videoRepo.GetVideos(perPage, offset, category)
}

// GetVideoByID retrieves a video by its ID and increments view count
func (s *VideoService) GetVideoByID(id uuid.UUID, userID *uuid.UUID) (*models.Video, error) {
	video, err := s.videoRepo.GetVideoByID(id)
	if err != nil {
		return nil, err
	}

	// Increment view count
	video.Views++
	if err := s.videoRepo.UpdateVideo(video); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update view count for video %s: %v\n", id, err)
	}

	// Record view interaction if user is provided
	if userID != nil {
		interaction := &models.VideoInteraction{
			VideoID: id,
			UserID:  *userID,
			Type:    "view",
		}
		if err := s.videoRepo.CreateVideoInteraction(interaction); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to record view interaction for video %s: %v\n", id, err)
		}
	}

	return video, nil
}

// CreateVideo creates a new video
func (s *VideoService) CreateVideo(video *models.Video) error {
	// Validate channel exists
	if _, err := s.videoRepo.GetChannelByID(video.ChannelID); err != nil {
		return fmt.Errorf("channel not found: %w", err)
	}

	return s.videoRepo.CreateVideo(video)
}

// UpdateVideo updates an existing video
func (s *VideoService) UpdateVideo(video *models.Video) error {
	// Check if video exists
	existing, err := s.videoRepo.GetVideoByID(video.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	video.UpdatedAt = time.Now()

	// Preserve some fields that shouldn't be updated
	video.CreatedAt = existing.CreatedAt
	video.Views = existing.Views
	video.Likes = existing.Likes

	return s.videoRepo.UpdateVideo(video)
}

// DeleteVideo deletes a video
func (s *VideoService) DeleteVideo(id uuid.UUID) error {
	// Check if video exists
	if _, err := s.videoRepo.GetVideoByID(id); err != nil {
		return err
	}

	return s.videoRepo.DeleteVideo(id)
}

// LikeVideo likes a video
func (s *VideoService) LikeVideo(videoID, userID uuid.UUID) error {
	// Check if video exists
	video, err := s.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Create like interaction
	interaction := &models.VideoInteraction{
		VideoID: videoID,
		UserID:  userID,
		Type:    "like",
	}

	if err := s.videoRepo.CreateVideoInteraction(interaction); err != nil {
		return err
	}

	// Increment like count
	video.Likes++
	return s.videoRepo.UpdateVideo(video)
}

// GetShortVideos retrieves short-form videos (TikTok style)
func (s *VideoService) GetShortVideos(page, perPage int) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	return s.videoRepo.GetVideos(perPage, offset, "")
}

// CreateChannel creates a new channel
func (s *VideoService) CreateChannel(channel *models.Channel) error {
	return s.videoRepo.CreateChannel(channel)
}

// GetChannelByID retrieves a channel by its ID
func (s *VideoService) GetChannelByID(id uuid.UUID) (*models.Channel, error) {
	return s.videoRepo.GetChannelByID(id)
}

// SearchVideos searches videos by query with optional category filter
func (s *VideoService) SearchVideos(query string, page, perPage int, category string) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	return s.videoRepo.SearchVideos(query, perPage, offset, category)
}

// GetVideosByCategory retrieves videos by category
func (s *VideoService) GetVideosByCategory(category string, page, perPage int) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	return s.videoRepo.GetVideosByCategory(category, perPage, offset)
}

// GetTrendingVideos retrieves trending videos
func (s *VideoService) GetTrendingVideos(timeframe string, page, perPage int) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate timeframe
	validTimeframes := []string{"day", "week", "month"}
	isValid := false
	for _, valid := range validTimeframes {
		if timeframe == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		timeframe = "day" // Default to day
	}

	offset := (page - 1) * perPage
	return s.videoRepo.GetTrendingVideos(timeframe, perPage, offset)
}

// GetVideoComments retrieves comments for a video
func (s *VideoService) GetVideoComments(videoID uuid.UUID, page, perPage int) ([]models.VideoComment, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Check if video exists
	if _, err := s.videoRepo.GetVideoByID(videoID); err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	return s.videoRepo.GetVideoComments(videoID, perPage, offset)
}

// AddVideoComment adds a comment to a video
func (s *VideoService) AddVideoComment(videoID, userID uuid.UUID, content string, parentID *uuid.UUID) (*models.VideoComment, error) {
	// Check if video exists
	if _, err := s.videoRepo.GetVideoByID(videoID); err != nil {
		return nil, err
	}

	// TODO: Get user details from user service
	// For now, using placeholder values
	comment := &models.VideoComment{
		VideoID:   videoID,
		UserID:    userID,
		UserName:  "User", // TODO: Get from user service
		Content:   content,
		ParentID:  parentID,
	}

	if err := s.videoRepo.CreateVideoComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// ShareVideo handles video sharing
func (s *VideoService) ShareVideo(videoID, userID uuid.UUID, platform, message string) (*models.VideoShare, error) {
	// Check if video exists
	video, err := s.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	// Generate share URL
	shareURL := fmt.Sprintf("https://tchat.dev/videos/%s?utm_source=%s", videoID, platform)

	share := &models.VideoShare{
		VideoID:  videoID,
		UserID:   userID,
		Platform: platform,
		Message:  message,
		ShareURL: shareURL,
	}

	if err := s.videoRepo.CreateVideoShare(share); err != nil {
		return nil, err
	}

	// TODO: Send share notification to video owner
	// TODO: Track share analytics

	// Increment share count in video interactions
	interaction := &models.VideoInteraction{
		VideoID: videoID,
		UserID:  userID,
		Type:    "share",
	}
	if err := s.videoRepo.CreateVideoInteraction(interaction); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to record share interaction for video %s: %v\n", videoID, err)
	}

	share.Video = *video
	return share, nil
}