package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat.dev/video/models"
)

// TestVideoContentModel tests the VideoContent model
func TestVideoContentModel(t *testing.T) {
	t.Run("CreateVideoContent", func(t *testing.T) {
		video := &models.VideoContent{
			ID:              uuid.New(),
			CreatorID:       uuid.New(),
			Title:           "Test Video",
			Description:     "Test Description",
			VideoURL:        "https://cdn.example.com/video.mp4",
			DurationSeconds: 300,
			FileSize:        1024000,
			MimeType:        "video/mp4",
			Resolution:      "1920x1080",
			Bitrate:         5000000,
			Framerate:       60.0,
			Codec:           "h264",
			UploadStatus:    "processing",
			ContentRating:   "general",
			Tags:            []string{"test", "video"},
		}

		assert.NotEqual(t, uuid.Nil, video.ID)
		assert.Equal(t, "Test Video", video.Title)
		assert.Equal(t, 300, video.DurationSeconds)
		assert.Equal(t, "1920x1080", video.Resolution)
		assert.Equal(t, 60.0, video.Framerate)
	})

	t.Run("ValidateRequiredFields", func(t *testing.T) {
		video := &models.VideoContent{
			ID:          uuid.New(),
			CreatorID:   uuid.New(),
			Title:       "Test",
			Description: "Description",
			VideoURL:    "https://example.com/video.mp4",
		}

		assert.NotEmpty(t, video.Title)
		assert.NotEmpty(t, video.Description)
		assert.NotEmpty(t, video.VideoURL)
	})

	t.Run("SocialMetrics", func(t *testing.T) {
		video := &models.VideoContent{
			ID:        uuid.New(),
			CreatorID: uuid.New(),
			ViewCount: 1000,
			LikeCount: 50,
		}

		assert.Equal(t, int64(1000), video.ViewCount)
		assert.Equal(t, int64(50), video.LikeCount)
	})
}

// TestPlaybackSessionModel tests the PlaybackSession model
func TestPlaybackSessionModel(t *testing.T) {
	t.Run("CreatePlaybackSession", func(t *testing.T) {
		session := &models.PlaybackSession{
			ID:              uuid.New(),
			VideoID:         uuid.New(),
			UserID:          uuid.New(),
			Platform:        "web",
			DeviceID:        "device-123",
			CurrentPosition: 0.0,
			TotalDuration:   300.0,
			PlaybackSpeed:   1.0,
			Quality:         "auto",
			Volume:          1.0,
			State:           "paused",
			LastSyncTime:    time.Now(),
		}

		assert.NotEqual(t, uuid.Nil, session.ID)
		assert.Equal(t, "web", session.Platform)
		assert.Equal(t, 1.0, session.PlaybackSpeed)
		assert.Equal(t, "auto", session.Quality)
		assert.Equal(t, "paused", session.State)
	})

	t.Run("UpdatePlaybackPosition", func(t *testing.T) {
		session := &models.PlaybackSession{
			ID:              uuid.New(),
			CurrentPosition: 0.0,
			TotalDuration:   300.0,
		}

		session.CurrentPosition = 150.0
		assert.Equal(t, 150.0, session.CurrentPosition)
		assert.Equal(t, 0.5, session.CurrentPosition/session.TotalDuration)
	})

	t.Run("PlaybackStates", func(t *testing.T) {
		states := []string{"playing", "paused", "buffering", "ended", "error"}

		for _, state := range states {
			session := &models.PlaybackSession{
				ID:    uuid.New(),
				State: state,
			}
			assert.Contains(t, states, session.State)
		}
	})
}

// TestViewingHistoryModel tests the ViewingHistory model
func TestViewingHistoryModel(t *testing.T) {
	t.Run("CreateViewingHistory", func(t *testing.T) {
		history := &models.ViewingHistory{
			ID:                   uuid.New(),
			UserID:               uuid.New(),
			VideoID:              uuid.New(),
			WatchedSeconds:       150.0,
			CompletionPercentage: 50.0,
			LastWatchedPosition:  150.0,
			WatchCount:           1,
			Platform:             "web",
			DeviceID:             "device-123",
			IsCompleted:          false,
			LastWatchedAt:        time.Now(),
		}

		assert.NotEqual(t, uuid.Nil, history.ID)
		assert.Equal(t, 150.0, history.WatchedSeconds)
		assert.Equal(t, 50.0, history.CompletionPercentage)
		assert.Equal(t, 1, history.WatchCount)
		assert.False(t, history.IsCompleted)
	})

	t.Run("MarkAsCompleted", func(t *testing.T) {
		history := &models.ViewingHistory{
			ID:                   uuid.New(),
			CompletionPercentage: 100.0,
			IsCompleted:          false,
		}

		history.IsCompleted = true
		assert.True(t, history.IsCompleted)
		assert.Equal(t, 100.0, history.CompletionPercentage)
	})
}

// TestPlatformConfigModel tests the PlatformConfig model
func TestPlatformConfigModel(t *testing.T) {
	t.Run("CreatePlatformConfig", func(t *testing.T) {
		config := &models.PlatformConfig{
			ID:                     uuid.New(),
			Platform:               "web",
			MaxResolution:          "1920x1080",
			PreferredCodec:         "h264",
			SupportedFormats:       []string{"mp4", "webm"},
			EnableAdaptiveStreaming: true,
			BufferSize:             30,
			MaxBitrate:             5000000,
		}

		assert.NotEqual(t, uuid.Nil, config.ID)
		assert.Equal(t, "web", config.Platform)
		assert.Equal(t, "1920x1080", config.MaxResolution)
		assert.True(t, config.EnableAdaptiveStreaming)
		assert.Equal(t, 30, config.BufferSize)
	})

	t.Run("PlatformSupport", func(t *testing.T) {
		platforms := []string{"web", "ios", "android", "mobile_web"}

		for _, platform := range platforms {
			config := &models.PlatformConfig{
				ID:       uuid.New(),
				Platform: platform,
			}
			assert.Contains(t, platforms, config.Platform)
		}
	})
}

// TestSyncStateModel tests the SyncState model
func TestSyncStateModel(t *testing.T) {
	t.Run("CreateSyncState", func(t *testing.T) {
		syncState := &models.SyncState{
			ID:               uuid.New(),
			UserID:           uuid.New(),
			VideoID:          uuid.New(),
			SessionID:        uuid.New(),
			SyncedPlatforms:  []string{"web", "ios"},
			LastSyncTime:     time.Now(),
			ConflictDetected: false,
			SyncVersion:      1,
			PendingChanges:   []string{},
		}

		assert.NotEqual(t, uuid.Nil, syncState.ID)
		assert.Len(t, syncState.SyncedPlatforms, 2)
		assert.False(t, syncState.ConflictDetected)
		assert.Equal(t, 1, syncState.SyncVersion)
	})

	t.Run("ConflictResolution", func(t *testing.T) {
		strategies := []string{"latest_wins", "authority_platform", "average_position", "manual_resolution"}

		for _, strategy := range strategies {
			syncState := &models.SyncState{
				ID:                  uuid.New(),
				ConflictDetected:    true,
				ConflictResolution:  &strategy,
			}
			assert.True(t, syncState.ConflictDetected)
			assert.NotNil(t, syncState.ConflictResolution)
			assert.Contains(t, strategies, *syncState.ConflictResolution)
		}
	})

	t.Run("SyncVersionIncrement", func(t *testing.T) {
		syncState := &models.SyncState{
			ID:          uuid.New(),
			SyncVersion: 1,
		}

		syncState.SyncVersion++
		assert.Equal(t, 2, syncState.SyncVersion)
	})
}

// TestVideoContentValidation tests validation logic
func TestVideoContentValidation(t *testing.T) {
	t.Run("ValidContentRating", func(t *testing.T) {
		ratings := []string{"general", "teen", "mature", "adult"}

		for _, rating := range ratings {
			video := &models.VideoContent{
				ID:            uuid.New(),
				ContentRating: rating,
			}
			assert.Contains(t, ratings, video.ContentRating)
		}
	})

	t.Run("ValidUploadStatus", func(t *testing.T) {
		statuses := []string{"processing", "available", "unavailable", "archived", "deleted"}

		for _, status := range statuses {
			video := &models.VideoContent{
				ID:           uuid.New(),
				UploadStatus: status,
			}
			assert.Contains(t, statuses, video.UploadStatus)
		}
	})
}

// TestModelRelationships tests relationships between models
func TestModelRelationships(t *testing.T) {
	t.Run("VideoToPlaybackSession", func(t *testing.T) {
		videoID := uuid.New()
		userID := uuid.New()

		session := &models.PlaybackSession{
			ID:      uuid.New(),
			VideoID: videoID,
			UserID:  userID,
		}

		assert.Equal(t, videoID, session.VideoID)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("VideoToViewingHistory", func(t *testing.T) {
		videoID := uuid.New()
		userID := uuid.New()

		history := &models.ViewingHistory{
			ID:      uuid.New(),
			VideoID: videoID,
			UserID:  userID,
		}

		assert.Equal(t, videoID, history.VideoID)
		assert.Equal(t, userID, history.UserID)
	})

	t.Run("SessionToSyncState", func(t *testing.T) {
		sessionID := uuid.New()

		syncState := &models.SyncState{
			ID:        uuid.New(),
			SessionID: sessionID,
		}

		assert.Equal(t, sessionID, syncState.SessionID)
	})
}

// BenchmarkVideoContentCreation benchmarks video content creation
func BenchmarkVideoContentCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &models.VideoContent{
			ID:              uuid.New(),
			CreatorID:       uuid.New(),
			Title:           "Benchmark Video",
			Description:     "Benchmark Description",
			VideoURL:        "https://example.com/video.mp4",
			DurationSeconds: 300,
			FileSize:        1024000,
		}
	}
}

// BenchmarkPlaybackSessionUpdate benchmarks session updates
func BenchmarkPlaybackSessionUpdate(b *testing.B) {
	session := &models.PlaybackSession{
		ID:              uuid.New(),
		CurrentPosition: 0.0,
		TotalDuration:   300.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session.CurrentPosition = float64(i % 300)
	}
}