package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/video/models"
)

// MockDataManager handles the creation and seeding of mock data for the video service
type MockDataManager struct {
	db *gorm.DB
}

// NewMockDataManager creates a new MockDataManager instance
func NewMockDataManager(db *gorm.DB) *MockDataManager {
	return &MockDataManager{db: db}
}

// SeedMockData seeds the database with comprehensive mock data
func (m *MockDataManager) SeedMockData() error {
	// Create mock channels first (required for videos)
	channels, err := m.createMockChannels()
	if err != nil {
		return err
	}

	// Create mock videos
	videos, err := m.createMockVideos(channels)
	if err != nil {
		return err
	}

	// Create mock video interactions
	if err := m.createMockVideoInteractions(videos); err != nil {
		return err
	}

	// Create mock video comments
	if err := m.createMockVideoComments(videos); err != nil {
		return err
	}

	// Create mock video shares
	if err := m.createMockVideoShares(videos); err != nil {
		return err
	}

	return nil
}

// createMockChannels creates realistic channel data
func (m *MockDataManager) createMockChannels() ([]models.Channel, error) {
	channels := []models.Channel{
		{
			ID:          uuid.New(),
			Name:        "TechTalk Thailand",
			Avatar:      "https://cdn.tchat.dev/avatars/techtalk-thailand.jpg",
			Subscribers: 125000,
			Verified:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-2, 0, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Bangkok Food Adventures",
			Avatar:      "https://cdn.tchat.dev/avatars/bangkok-food.jpg",
			Subscribers: 89500,
			Verified:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-1, -6, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Singapore Travel Guide",
			Avatar:      "https://cdn.tchat.dev/avatars/singapore-travel.jpg",
			Subscribers: 67200,
			Verified:    false,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-1, -3, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Indonesia Gaming Zone",
			Avatar:      "https://cdn.tchat.dev/avatars/indonesia-gaming.jpg",
			Subscribers: 234500,
			Verified:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-3, 0, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Malaysian Lifestyle",
			Avatar:      "https://cdn.tchat.dev/avatars/malaysian-lifestyle.jpg",
			Subscribers: 45600,
			Verified:    false,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(0, -8, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Vietnam Culture Hub",
			Avatar:      "https://cdn.tchat.dev/avatars/vietnam-culture.jpg",
			Subscribers: 78900,
			Verified:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-1, -9, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Philippines Entertainment",
			Avatar:      "https://cdn.tchat.dev/avatars/philippines-entertainment.jpg",
			Subscribers: 156700,
			Verified:    true,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(-2, -4, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Cambodia Adventures",
			Avatar:      "https://cdn.tchat.dev/avatars/cambodia-adventures.jpg",
			Subscribers: 23400,
			Verified:    false,
			UserID:      uuid.New(),
			CreatedAt:   time.Now().AddDate(0, -5, 0),
			UpdatedAt:   time.Now(),
		},
	}

	for _, channel := range channels {
		if err := m.db.Create(&channel).Error; err != nil {
			return nil, err
		}
	}

	return channels, nil
}

// createMockVideos creates realistic video content with diverse categories
func (m *MockDataManager) createMockVideos(channels []models.Channel) ([]models.Video, error) {
	videos := []models.Video{
		// Tech Content
		{
			ID:           uuid.New(),
			Title:        "React 18 Features You Must Know in 2024",
			Description:  "Comprehensive guide to React 18's new features including concurrent rendering, automatic batching, and more!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/react-18-features.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/react-18-features.mp4",
			Duration:     "12:45",
			Views:        45600,
			Likes:        3200,
			Category:     "technology",
			Tags:         []string{"react", "javascript", "frontend", "web development"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[0].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -7),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Title:        "Building REST APIs with Go Gin",
			Description:  "Learn how to build scalable REST APIs using Go and the Gin framework. Perfect for backend developers!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/go-gin-api.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/go-gin-api.mp4",
			Duration:     "18:30",
			Views:        28900,
			Likes:        2100,
			Category:     "technology",
			Tags:         []string{"golang", "api", "backend", "gin", "programming"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[0].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -14),
			UpdatedAt:    time.Now(),
		},

		// Food Content
		{
			ID:           uuid.New(),
			Title:        "Best Street Food in Bangkok 2024",
			Description:  "Exploring the most delicious and authentic street food spots in Bangkok. Don't miss these hidden gems!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/bangkok-street-food.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/bangkok-street-food.mp4",
			Duration:     "15:20",
			Views:        76500,
			Likes:        5400,
			Category:     "food",
			Tags:         []string{"bangkok", "street food", "thailand", "travel", "local cuisine"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[1].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -3),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Title:        "Pad Thai Recipe - Authentic Thai Style",
			Description:  "Learn to make authentic Pad Thai from a local Bangkok chef. Secret ingredients revealed!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/pad-thai-recipe.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/pad-thai-recipe.mp4",
			Duration:     "08:45",
			Views:        124300,
			Likes:        8900,
			Category:     "food",
			Tags:         []string{"pad thai", "recipe", "cooking", "thai cuisine", "authentic"},
			Type:         "short",
			Status:       "active",
			ChannelID:    channels[1].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -10),
			UpdatedAt:    time.Now(),
		},

		// Travel Content
		{
			ID:           uuid.New(),
			Title:        "Marina Bay Sands Complete Guide",
			Description:  "Everything you need to know about visiting Marina Bay Sands in Singapore. Tips, tricks, and hidden spots!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/marina-bay-sands.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/marina-bay-sands.mp4",
			Duration:     "22:15",
			Views:        89700,
			Likes:        6200,
			Category:     "travel",
			Tags:         []string{"singapore", "marina bay sands", "travel guide", "tourism", "attractions"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[2].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -5),
			UpdatedAt:    time.Now(),
		},

		// Gaming Content
		{
			ID:           uuid.New(),
			Title:        "Mobile Legends Pro Tips 2024",
			Description:  "Master Mobile Legends with these advanced strategies and hero guides. Rank up fast!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/mobile-legends-tips.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/mobile-legends-tips.mp4",
			Duration:     "16:40",
			Views:        156800,
			Likes:        12400,
			Category:     "gaming",
			Tags:         []string{"mobile legends", "gaming", "esports", "indonesia", "strategy"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[3].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -2),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Title:        "Epic Comeback in Ranked Match",
			Description:  "Watch this insane comeback from a 15k gold deficit! Never give up!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/epic-comeback.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/epic-comeback.mp4",
			Duration:     "03:25",
			Views:        234500,
			Likes:        18900,
			Category:     "gaming",
			Tags:         []string{"comeback", "gaming", "highlights", "epic", "clutch"},
			Type:         "short",
			Status:       "active",
			ChannelID:    channels[3].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -1),
			UpdatedAt:    time.Now(),
		},

		// Lifestyle Content
		{
			ID:           uuid.New(),
			Title:        "Kuala Lumpur Shopping Haul",
			Description:  "My latest shopping haul from Pavilion KL and Suria KLCC. Fashion finds and beauty products!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/kl-shopping-haul.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/kl-shopping-haul.mp4",
			Duration:     "19:30",
			Views:        67800,
			Likes:        4500,
			Category:     "lifestyle",
			Tags:         []string{"shopping", "fashion", "kuala lumpur", "haul", "beauty"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[4].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -6),
			UpdatedAt:    time.Now(),
		},

		// Culture Content
		{
			ID:           uuid.New(),
			Title:        "Vietnamese Coffee Culture Explained",
			Description:  "Deep dive into Vietnam's rich coffee culture. From traditional drip coffee to modern cafe trends!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/vietnamese-coffee.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/vietnamese-coffee.mp4",
			Duration:     "14:20",
			Views:        45300,
			Likes:        3100,
			Category:     "culture",
			Tags:         []string{"vietnam", "coffee", "culture", "tradition", "cafe"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[5].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -8),
			UpdatedAt:    time.Now(),
		},

		// Entertainment Content
		{
			ID:           uuid.New(),
			Title:        "Filipino Dance Challenge Compilation",
			Description:  "The best Filipino dance challenges taking over social media! Can you keep up with these moves?",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/filipino-dance.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/filipino-dance.mp4",
			Duration:     "07:15",
			Views:        189400,
			Likes:        15600,
			Category:     "entertainment",
			Tags:         []string{"dance", "philippines", "trending", "social media", "challenge"},
			Type:         "short",
			Status:       "active",
			ChannelID:    channels[6].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -4),
			UpdatedAt:    time.Now(),
		},

		// Adventure Content
		{
			ID:           uuid.New(),
			Title:        "Angkor Wat Sunrise Photography Tips",
			Description:  "Capture the perfect sunrise at Angkor Wat with these professional photography techniques and location guides!",
			ThumbnailURL: "https://cdn.tchat.dev/thumbnails/angkor-wat-sunrise.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/angkor-wat-sunrise.mp4",
			Duration:     "11:50",
			Views:        34200,
			Likes:        2800,
			Category:     "travel",
			Tags:         []string{"angkor wat", "cambodia", "photography", "sunrise", "travel tips"},
			Type:         "long",
			Status:       "active",
			ChannelID:    channels[7].ID,
			CreatedAt:    time.Now().AddDate(0, 0, -9),
			UpdatedAt:    time.Now(),
		},
	}

	for _, video := range videos {
		if err := m.db.Create(&video).Error; err != nil {
			return nil, err
		}
	}

	return videos, nil
}

// createMockVideoInteractions creates realistic user interactions
func (m *MockDataManager) createMockVideoInteractions(videos []models.Video) error {
	interactionTypes := []string{"like", "view", "share"}

	for _, video := range videos {
		// Create multiple interactions per video
		numInteractions := int(video.Views/1000) + 10 // Base interactions on view count
		if numInteractions > 100 {
			numInteractions = 100 // Cap interactions
		}

		for i := 0; i < numInteractions; i++ {
			interaction := models.VideoInteraction{
				ID:        uuid.New(),
				VideoID:   video.ID,
				UserID:    uuid.New(), // Random user IDs
				Type:      interactionTypes[i%len(interactionTypes)],
				CreatedAt: time.Now().AddDate(0, 0, -i%30), // Spread over last 30 days
			}

			if err := m.db.Create(&interaction).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// createMockVideoComments creates realistic comments and replies
func (m *MockDataManager) createMockVideoComments(videos []models.Video) error {
	comments := []string{
		"Great content! Thanks for sharing this valuable information.",
		"This is exactly what I was looking for. Very helpful!",
		"Amazing video quality and explanation. Subscribed!",
		"Can you make a follow-up video on this topic?",
		"Love your content! Keep up the great work.",
		"This helped me so much with my project. Thank you!",
		"Your tutorials are the best on the platform!",
		"Perfect timing for this video. Exactly what I needed.",
		"The production quality keeps getting better!",
		"Could you do a beginner's version of this?",
		"This is why I love this channel!",
		"Fantastic explanation. Very easy to understand.",
		"Please make more videos like this one!",
		"Your channel deserves more subscribers!",
		"This content is pure gold. Thank you!",
	}

	userNames := []string{
		"Alex_Chen", "Maya_Singh", "David_Wong", "Sarah_Kim",
		"Mike_Torres", "Lisa_Nguyen", "Ryan_Lee", "Emma_Tan",
		"Jake_Sato", "Nina_Raj", "Tom_Liu", "Zoe_Park",
		"Ben_Kumar", "Amy_Zhang", "Sam_Patel", "Ava_Lim",
	}

	userAvatars := []string{
		"https://cdn.tchat.dev/avatars/user1.jpg",
		"https://cdn.tchat.dev/avatars/user2.jpg",
		"https://cdn.tchat.dev/avatars/user3.jpg",
		"https://cdn.tchat.dev/avatars/user4.jpg",
		"https://cdn.tchat.dev/avatars/user5.jpg",
	}

	for _, video := range videos {
		// Create 3-8 comments per video
		numComments := 3 + int(video.Views/10000)%6
		if numComments > 15 {
			numComments = 15
		}

		for i := 0; i < numComments; i++ {
			comment := models.VideoComment{
				ID:         uuid.New(),
				VideoID:    video.ID,
				UserID:     uuid.New(),
				UserName:   userNames[i%len(userNames)],
				UserAvatar: userAvatars[i%len(userAvatars)],
				Content:    comments[i%len(comments)],
				ParentID:   nil, // Top-level comment
				Likes:      int64(i*3 + 1),
				IsEdited:   i%5 == 0, // 20% of comments are edited
				CreatedAt:  time.Now().AddDate(0, 0, -(i+1)),
				UpdatedAt:  time.Now(),
			}

			if err := m.db.Create(&comment).Error; err != nil {
				return err
			}

			// Create 1-2 replies for some comments
			if i%3 == 0 && i < numComments-1 {
				reply := models.VideoComment{
					ID:         uuid.New(),
					VideoID:    video.ID,
					UserID:     uuid.New(),
					UserName:   userNames[(i+5)%len(userNames)],
					UserAvatar: userAvatars[(i+2)%len(userAvatars)],
					Content:    "Thanks for the kind words! More content coming soon.",
					ParentID:   &comment.ID,
					Likes:      int64(i + 2),
					IsEdited:   false,
					CreatedAt:  time.Now().AddDate(0, 0, -i),
					UpdatedAt:  time.Now(),
				}

				if err := m.db.Create(&reply).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// createMockVideoShares creates realistic share data
func (m *MockDataManager) createMockVideoShares(videos []models.Video) error {
	platforms := []string{"facebook", "twitter", "instagram", "tiktok", "whatsapp", "telegram", "discord"}
	shareMessages := []string{
		"Check out this amazing video!",
		"This is so helpful, you should watch it!",
		"Found this great tutorial, thought you'd like it",
		"Perfect explanation of this topic",
		"This video changed my perspective",
		"Must watch! Really insightful content",
		"Sharing this gem with everyone",
	}

	for _, video := range videos {
		// Create shares based on video popularity
		numShares := int(video.Views/5000) + 1
		if numShares > 20 {
			numShares = 20
		}

		for i := 0; i < numShares; i++ {
			share := models.VideoShare{
				ID:        uuid.New(),
				VideoID:   video.ID,
				UserID:    uuid.New(),
				Platform:  platforms[i%len(platforms)],
				Message:   shareMessages[i%len(shareMessages)],
				ShareURL:  "https://tchat.dev/videos/" + video.ID.String() + "?utm_source=" + platforms[i%len(platforms)],
				CreatedAt: time.Now().AddDate(0, 0, -(i%15)), // Spread over last 15 days
			}

			if err := m.db.Create(&share).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// ClearMockData removes all mock data from the database
func (m *MockDataManager) ClearMockData() error {
	// Delete in reverse order of dependencies
	if err := m.db.Exec("DELETE FROM video_shares").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM video_comments").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM video_interactions").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM videos").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM channels").Error; err != nil {
		return err
	}

	return nil
}