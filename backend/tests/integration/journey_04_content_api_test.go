// Journey 4: Content Creation & Discovery API Integration Tests
// Tests all API endpoints involved in content creation, publishing, and discovery

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Journey04ContentAPISuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	ctx        context.Context
	creator    *AuthenticatedUser // Maria (Philippines) - Content Creator
	viewer     *AuthenticatedUser // Maya (Thailand) - Content Viewer
}

// Note: AuthenticatedUser is now defined in types.go

// Note: CreateContentRequest is now defined in types.go

type MonetizationSettings struct {
	Enabled         bool    `json:"enabled"`
	TipJarEnabled   bool    `json:"tipJarEnabled"`
	AdRevenueShare  bool    `json:"adRevenueShare"`
	SubscriptionTier string `json:"subscriptionTier,omitempty"` // "free", "premium", "exclusive"
	MinimumTip      int64   `json:"minimumTip,omitempty"`       // In minor currency units
	Currency        string  `json:"currency,omitempty"`
}

// Note: LocationData is now defined in types.go

// Note: Coordinates is now defined in types.go

type VideoUploadData struct {
	Duration   int    `json:"duration"`   // Duration in seconds
	Resolution string `json:"resolution"` // "720p", "1080p", "4k"
	FileSize   int64  `json:"fileSize"`   // File size in bytes
	Format     string `json:"format"`     // "mp4", "mov", "webm"
	Codec      string `json:"codec,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"` // Bitrate in kbps
}

type ImageUploadData struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FileSize   int64  `json:"fileSize"`
	Format     string `json:"format"` // "jpg", "png", "webp", "gif"
	Quality    int    `json:"quality,omitempty"`
	IsAnimated bool   `json:"isAnimated,omitempty"`
}

type AudioUploadData struct {
	Duration int    `json:"duration"` // Duration in seconds
	FileSize int64  `json:"fileSize"`
	Format   string `json:"format"`   // "mp3", "wav", "aac", "ogg"
	Bitrate  int    `json:"bitrate"`  // Bitrate in kbps
	SampleRate int  `json:"sampleRate,omitempty"`
}

type ContentScheduling struct {
	PublishAt time.Time `json:"publishAt"`
	Timezone  string    `json:"timezone"`
}

// Note: ContentResponse is now defined in types.go

type MediaURLs struct {
	Original     string   `json:"original,omitempty"`
	Processed    string   `json:"processed,omitempty"`
	Thumbnail    string   `json:"thumbnail,omitempty"`
	Previews     []string `json:"previews,omitempty"`     // Different resolutions
	Transcripts  []string `json:"transcripts,omitempty"`  // For video/audio
	Captions     []string `json:"captions,omitempty"`     // For accessibility
}

type ContentAnalytics struct {
	ViewCount      int     `json:"viewCount"`
	LikeCount      int     `json:"likeCount"`
	CommentCount   int     `json:"commentCount"`
	ShareCount     int     `json:"shareCount"`
	SaveCount      int     `json:"saveCount"`
	EngagementRate float64 `json:"engagementRate"`
	WatchTime      int     `json:"watchTime,omitempty"`      // For video/audio in seconds
	CompletionRate float64 `json:"completionRate,omitempty"` // For video/audio
	Revenue        float64 `json:"revenue,omitempty"`
	RevenueToday   float64 `json:"revenueToday,omitempty"`
}

type DiscoverContentResponse struct {
	Content      []ContentSummary   `json:"content"`
	TotalCount   int                `json:"totalCount"`
	Page         int                `json:"page"`
	PageSize     int                `json:"pageSize"`
	Filters      DiscoveryFilters   `json:"filters"`
	Trending     []TrendingTag      `json:"trending,omitempty"`
}

type ContentSummary struct {
	ID             string           `json:"id"`
	CreatorID      string           `json:"creatorId"`
	CreatorName    string           `json:"creatorName"`
	CreatorAvatar  string           `json:"creatorAvatar"`
	Type           string           `json:"type"`
	Title          string           `json:"title"`
	Description    string           `json:"description,omitempty"`
	ThumbnailURL   string           `json:"thumbnailUrl,omitempty"`
	Duration       int              `json:"duration,omitempty"`
	ViewCount      int              `json:"viewCount"`
	LikeCount      int              `json:"likeCount"`
	CommentCount   int              `json:"commentCount"`
	PublishedAt    string           `json:"publishedAt"`
	Category       string           `json:"category"`
	Language       string           `json:"language"`
	Tags           []string         `json:"tags,omitempty"`
	Location       LocationData     `json:"location,omitempty"`
}

type DiscoveryFilters struct {
	Region        string   `json:"region,omitempty"`
	Language      string   `json:"language,omitempty"`
	Category      string   `json:"category,omitempty"`
	ContentType   string   `json:"contentType,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	SortBy        string   `json:"sortBy"` // "trending", "latest", "popular", "relevant"
	TimeRange     string   `json:"timeRange,omitempty"` // "today", "week", "month", "year", "all"
	MonetizedOnly bool     `json:"monetizedOnly,omitempty"`
}

type TrendingTag struct {
	Tag        string `json:"tag"`
	Count      int    `json:"count"`
	Growth     float64 `json:"growth"` // Percentage growth
}

type CreateCommentRequest struct {
	ContentID string `json:"contentId"`
	Text      string `json:"text"`
	Language  string `json:"language"`
	ParentID  string `json:"parentId,omitempty"` // For replies
}

type CommentResponse struct {
	ID          string    `json:"id"`
	ContentID   string    `json:"contentId"`
	UserID      string    `json:"userId"`
	UserName    string    `json:"userName"`
	UserAvatar  string    `json:"userAvatar,omitempty"`
	Text        string    `json:"text"`
	Language    string    `json:"language"`
	ParentID    string    `json:"parentId,omitempty"`
	Replies     []CommentResponse `json:"replies,omitempty"`
	LikeCount   int       `json:"likeCount"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt,omitempty"`
}

type ContentEngagementRequest struct {
	ContentID string                 `json:"contentId"`
	Action    string                 `json:"action"` // "view", "like", "unlike", "save", "unsave", "share"
	Duration  int                    `json:"duration,omitempty"` // Watch duration for videos
	Quality   string                 `json:"quality,omitempty"`  // Video quality watched
	Device    string                 `json:"device,omitempty"`   // "mobile", "desktop", "tv"
	Location  LocationData           `json:"location,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type CreatorAnalyticsResponse struct {
	CreatorID         string                  `json:"creatorId"`
	Period            string                  `json:"period"` // "today", "week", "month", "year", "all"
	ContentCount      int                     `json:"contentCount"`
	TotalViews        int                     `json:"totalViews"`
	TotalLikes        int                     `json:"totalLikes"`
	TotalComments     int                     `json:"totalComments"`
	TotalShares       int                     `json:"totalShares"`
	SubscriberCount   int                     `json:"subscriberCount"`
	Revenue           CreatorRevenue          `json:"revenue"`
	TopContent        []ContentSummary        `json:"topContent"`
	AudienceBreakdown AudienceAnalytics       `json:"audience"`
	EngagementTrends  []EngagementDataPoint   `json:"engagementTrends"`
}

type CreatorRevenue struct {
	TotalRevenue  float64 `json:"totalRevenue"`
	TipsReceived  float64 `json:"tipsReceived"`
	AdRevenue     float64 `json:"adRevenue"`
	Premium       float64 `json:"premium"`
	Currency      string  `json:"currency"`
	PendingPayout float64 `json:"pendingPayout"`
}

type AudienceAnalytics struct {
	Countries    map[string]int    `json:"countries"`
	AgeGroups    map[string]int    `json:"ageGroups"`
	Gender       map[string]int    `json:"gender"`
	Devices      map[string]int    `json:"devices"`
	Languages    map[string]int    `json:"languages"`
	WatchTime    map[string]int    `json:"watchTime"` // By time of day
}

type EngagementDataPoint struct {
	Date        string  `json:"date"`
	Views       int     `json:"views"`
	Likes       int     `json:"likes"`
	Comments    int     `json:"comments"`
	Shares      int     `json:"shares"`
	Revenue     float64 `json:"revenue"`
}

func (suite *Journey04ContentAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081" // API Gateway
	suite.httpClient = &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for media uploads
	}
	suite.ctx = context.Background()

	// Create authenticated test users
	suite.creator = suite.createAuthenticatedUser("maria@test.com", "PH", "tl") // Content creator from Philippines
	suite.viewer = suite.createAuthenticatedUser("maya@test.com", "TH", "th")    // Content viewer from Thailand
}

// Test 4.1: Content Creation API
func (suite *Journey04ContentAPISuite) TestContentCreationAPI() {
	// Step 1: POST /api/v1/content - Create video content
	contentReq := CreateContentRequest{
		Type:        "video",
		Title:       "Traditional Filipino Dance Tutorial - Tinikling",
		Description: "Learn the beautiful and traditional Tinikling dance step by step! Perfect for beginners who want to explore Filipino culture through dance.",
		Tags:        []string{"dance", "tutorial", "philippines", "traditional", "culture", "tinikling", "bamboo"},
		Category:    "education",
		Language:    "tl",
		Visibility:  "public",
		Monetization: MonetizationSettings{
			Enabled:        true,
			TipJarEnabled:  true,
			AdRevenueShare: true,
			MinimumTip:     10000, // 100 PHP in centavos
			Currency:       "PHP",
		},
		Location: LocationData{
			Country: "PH",
			City:    "Manila",
			Coordinates: Coordinates{
				Latitude:  14.5995,
				Longitude: 120.9842,
			},
			PlaceName: "Cultural Center of the Philippines",
		},
		VideoData: &VideoUploadData{
			Duration:   300, // 5 minutes
			Resolution: "1080p",
			FileSize:   45000000, // 45MB
			Format:     "mp4",
			Codec:      "h264",
			Bitrate:    2000,
		},
		Metadata: map[string]interface{}{
			"difficulty":     "beginner",
			"equipment":      "bamboo poles",
			"participants":   "2-4 people",
			"culturalNote":   "Traditional dance from Leyte province",
		},
	}

	creatorHeaders := map[string]string{
		"Authorization":   "Bearer " + suite.creator.AccessToken,
		"Accept-Language": "tl",
	}

	contentResp, statusCode := suite.makeAPICall("POST", "/api/v1/content", contentReq, creatorHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Content creation should succeed")

	var content ContentResponse
	err := json.Unmarshal(contentResp, &content)
	require.NoError(suite.T(), err, "Should parse content response")

	contentID := content.ID
	assert.NotEmpty(suite.T(), contentID, "Should return content ID")
	assert.Equal(suite.T(), suite.creator.UserID, content.CreatorID, "Creator ID should match")
	assert.Equal(suite.T(), "Traditional Filipino Dance Tutorial - Tinikling", content.Title, "Title should match")
	assert.Equal(suite.T(), "video", content.Type, "Type should be video")
	assert.Equal(suite.T(), "processing", content.Status, "Status should be processing initially")

	// Step 2: Upload video file - POST /api/v1/content/{id}/upload
	videoData := suite.generateTestVideoData()
	uploadResp, statusCode := suite.uploadMediaFile(
		fmt.Sprintf("/api/v1/content/%s/upload", contentID),
		"video/mp4",
		"tinikling_tutorial.mp4",
		videoData,
		creatorHeaders,
	)
	assert.Equal(suite.T(), 200, statusCode, "Video upload should succeed")

	var uploadResult map[string]interface{}
	err = json.Unmarshal(uploadResp, &uploadResult)
	require.NoError(suite.T(), err, "Should parse upload response")
	assert.Equal(suite.T(), "uploaded", uploadResult["status"], "Upload status should be uploaded")

	// Step 3: Wait for processing and check status
	time.Sleep(5 * time.Second) // Allow processing time

	processedContentResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve processed content")

	var processedContent ContentResponse
	err = json.Unmarshal(processedContentResp, &processedContent)
	require.NoError(suite.T(), err, "Should parse processed content")

	assert.Equal(suite.T(), "published", processedContent.Status, "Status should be published")
	assert.NotEmpty(suite.T(), processedContent.MediaURLs.Original, "Should have original URL")
	assert.NotEmpty(suite.T(), processedContent.MediaURLs.Thumbnail, "Should have thumbnail URL")
	assert.NotEmpty(suite.T(), processedContent.PublishedAt, "Should have published timestamp")

	// Step 4: GET /api/v1/content/creator/{creatorId} - List creator's content
	creatorContentResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/creator/%s", suite.creator.UserID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list creator's content")

	var creatorContent []ContentResponse
	err = json.Unmarshal(creatorContentResp, &creatorContent)
	require.NoError(suite.T(), err, "Should parse creator content")
	assert.Greater(suite.T(), len(creatorContent), 0, "Creator should have content")

	// Find our content
	found := false
	for _, c := range creatorContent {
		if c.ID == contentID {
			found = true
			assert.Equal(suite.T(), "published", c.Status, "Content should be published")
			break
		}
	}
	assert.True(suite.T(), found, "Should find created content in creator's list")
}

// Test 4.2: Content Discovery API
func (suite *Journey04ContentAPISuite) TestContentDiscoveryAPI() {
	// Create test content first
	contentID := suite.createPublishedContent()

	viewerHeaders := map[string]string{
		"Authorization":   "Bearer " + suite.viewer.AccessToken,
		"Accept-Language": "th",
	}

	// Step 1: GET /api/v1/content/discover - General discovery
	discoverResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/content/discover?limit=20&sortBy=latest", nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Content discovery should succeed")

	var discoverResult DiscoverContentResponse
	err := json.Unmarshal(discoverResp, &discoverResult)
	require.NoError(suite.T(), err, "Should parse discovery response")

	assert.Greater(suite.T(), discoverResult.TotalCount, 0, "Should have discoverable content")
	assert.Greater(suite.T(), len(discoverResult.Content), 0, "Should return content items")
	assert.Equal(suite.T(), "latest", discoverResult.Filters.SortBy, "Sort filter should match")

	// Step 2: GET /api/v1/content/discover - Category-specific discovery
	categoryDiscoverResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/content/discover?category=education&region=SEA&language=any", nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Category discovery should succeed")

	var categoryResult DiscoverContentResponse
	err = json.Unmarshal(categoryDiscoverResp, &categoryResult)
	require.NoError(suite.T(), err, "Should parse category discovery")

	assert.Greater(suite.T(), categoryResult.TotalCount, 0, "Should have educational content")
	assert.Equal(suite.T(), "education", categoryResult.Filters.Category, "Category filter should match")

	// Step 3: GET /api/v1/content/discover - Cross-region discovery
	crossRegionResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/content/discover?region=PH&language=tl&sortBy=trending", nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Cross-region discovery should succeed")

	var crossRegionResult DiscoverContentResponse
	err = json.Unmarshal(crossRegionResp, &crossRegionResult)
	require.NoError(suite.T(), err, "Should parse cross-region discovery")

	// Find our content in cross-region results
	var foundContent *ContentSummary
	for _, content := range crossRegionResult.Content {
		if content.ID == contentID {
			foundContent = &content
			break
		}
	}
	require.NotNil(suite.T(), foundContent, "Should find Philippine content from Thailand")
	assert.Equal(suite.T(), "PH", foundContent.Location.Country, "Content should be from Philippines")
	assert.Equal(suite.T(), suite.creator.UserID, foundContent.CreatorID, "Creator should match")

	// Step 4: GET /api/v1/content/trending - Trending content
	trendingResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/content/trending?timeRange=week&region=SEA", nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get trending content")

	var trendingResult DiscoverContentResponse
	err = json.Unmarshal(trendingResp, &trendingResult)
	require.NoError(suite.T(), err, "Should parse trending content")

	assert.Greater(suite.T(), len(trendingResult.Trending), 0, "Should have trending tags")

	// Step 5: GET /api/v1/content/search - Search functionality
	searchResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/content/search?q=dance+tutorial&language=tl&limit=10", nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Content search should succeed")

	var searchResult DiscoverContentResponse
	err = json.Unmarshal(searchResp, &searchResult)
	require.NoError(suite.T(), err, "Should parse search results")

	// Should find our dance tutorial
	foundInSearch := false
	for _, content := range searchResult.Content {
		if strings.Contains(strings.ToLower(content.Title), "dance") &&
		   strings.Contains(strings.ToLower(content.Title), "tutorial") {
			foundInSearch = true
			break
		}
	}
	assert.True(suite.T(), foundInSearch, "Should find dance tutorial in search results")
}

// Test 4.3: Content Engagement API
func (suite *Journey04ContentAPISuite) TestContentEngagementAPI() {
	contentID := suite.createPublishedContent()

	viewerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.viewer.AccessToken,
	}

	// Step 1: POST /api/v1/content/{id}/engagement - View content
	viewEngagementReq := ContentEngagementRequest{
		ContentID: contentID,
		Action:    "view",
		Duration:  180, // Watched 3 minutes of 5 minute video
		Quality:   "1080p",
		Device:    "mobile",
		Location: LocationData{
			Country: "TH",
			City:    "Bangkok",
		},
		Metadata: map[string]interface{}{
			"referrer": "discovery_feed",
			"session": "abc123",
		},
	}

	_, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/content/%s/engagement", contentID), viewEngagementReq, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "View engagement should succeed")

	// Step 2: POST /api/v1/content/{id}/engagement - Like content
	likeEngagementReq := ContentEngagementRequest{
		ContentID: contentID,
		Action:    "like",
	}

	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/content/%s/engagement", contentID), likeEngagementReq, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Like engagement should succeed")

	// Step 3: POST /api/v1/content/{id}/comments - Add comment
	commentReq := CreateCommentRequest{
		ContentID: contentID,
		Text:      "Amazing tutorial! I learned so much about Filipino culture. The step-by-step instructions were very clear. 🇵🇭❤️",
		Language:  "en",
	}

	commentResp, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/content/%s/comments", contentID), commentReq, viewerHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Comment creation should succeed")

	var comment CommentResponse
	err := json.Unmarshal(commentResp, &comment)
	require.NoError(suite.T(), err, "Should parse comment response")

	commentID := comment.ID
	assert.NotEmpty(suite.T(), commentID, "Should return comment ID")
	assert.Equal(suite.T(), contentID, comment.ContentID, "Content ID should match")
	assert.Equal(suite.T(), suite.viewer.UserID, comment.UserID, "User ID should match")
	assert.Contains(suite.T(), comment.Text, "Amazing tutorial", "Comment text should match")

	// Step 4: GET /api/v1/content/{id}/comments - List comments
	listCommentsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s/comments", contentID), nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list comments")

	var comments []CommentResponse
	err = json.Unmarshal(listCommentsResp, &comments)
	require.NoError(suite.T(), err, "Should parse comments list")
	assert.Len(suite.T(), comments, 1, "Should have 1 comment")
	assert.Equal(suite.T(), commentID, comments[0].ID, "Comment ID should match")

	// Step 5: POST /api/v1/content/{id}/comments - Reply to comment
	replyReq := CreateCommentRequest{
		ContentID: contentID,
		Text:      "Thank you so much! I'm glad you enjoyed learning about our culture! 😊",
		Language:  "en",
		ParentID:  commentID,
	}

	creatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.creator.AccessToken,
	}

	replyResp, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/content/%s/comments", contentID), replyReq, creatorHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Reply creation should succeed")

	var reply CommentResponse
	err = json.Unmarshal(replyResp, &reply)
	require.NoError(suite.T(), err, "Should parse reply response")
	assert.Equal(suite.T(), commentID, reply.ParentID, "Parent ID should match original comment")

	// Step 6: GET /api/v1/content/{id} - Verify analytics updated
	time.Sleep(2 * time.Second) // Allow analytics processing

	contentAnalyticsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, viewerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve content with analytics")

	var contentWithAnalytics ContentResponse
	err = json.Unmarshal(contentAnalyticsResp, &contentWithAnalytics)
	require.NoError(suite.T(), err, "Should parse content with analytics")

	assert.Equal(suite.T(), 1, contentWithAnalytics.Analytics.ViewCount, "View count should be 1")
	assert.Equal(suite.T(), 1, contentWithAnalytics.Analytics.LikeCount, "Like count should be 1")
	assert.Equal(suite.T(), 2, contentWithAnalytics.Analytics.CommentCount, "Comment count should be 2 (original + reply)")
	assert.Greater(suite.T(), contentWithAnalytics.Analytics.EngagementRate, float64(0), "Should have engagement rate")
}

// Test 4.4: Creator Analytics API
func (suite *Journey04ContentAPISuite) TestCreatorAnalyticsAPI() {
	// Create and engage with content first
	contentID := suite.createPublishedContent()
	suite.simulateEngagement(contentID)

	creatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.creator.AccessToken,
	}

	// Step 1: GET /api/v1/creator/{id}/analytics - General analytics
	analyticsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/creator/%s/analytics?period=week", suite.creator.UserID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get creator analytics")

	var analytics CreatorAnalyticsResponse
	err := json.Unmarshal(analyticsResp, &analytics)
	require.NoError(suite.T(), err, "Should parse analytics response")

	assert.Equal(suite.T(), suite.creator.UserID, analytics.CreatorID, "Creator ID should match")
	assert.Equal(suite.T(), "week", analytics.Period, "Period should match")
	assert.Greater(suite.T(), analytics.ContentCount, 0, "Should have content count")
	assert.Greater(suite.T(), analytics.TotalViews, 0, "Should have total views")
	assert.Greater(suite.T(), len(analytics.TopContent), 0, "Should have top content")

	// Step 2: GET /api/v1/creator/{id}/revenue - Revenue analytics
	revenueResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/creator/%s/revenue?period=today", suite.creator.UserID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get revenue analytics")

	var revenue CreatorRevenue
	err = json.Unmarshal(revenueResp, &revenue)
	require.NoError(suite.T(), err, "Should parse revenue response")

	assert.Equal(suite.T(), "PHP", revenue.Currency, "Currency should be PHP")
	assert.GreaterOrEqual(suite.T(), revenue.TotalRevenue, float64(0), "Should have revenue data")

	// Step 3: GET /api/v1/content/{id}/analytics - Content-specific analytics
	contentAnalyticsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s/analytics", contentID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get content analytics")

	var contentAnalytics map[string]interface{}
	err = json.Unmarshal(contentAnalyticsResp, &contentAnalytics)
	require.NoError(suite.T(), err, "Should parse content analytics")

	assert.Greater(suite.T(), contentAnalytics["viewCount"], float64(0), "Should have view count")
	assert.Greater(suite.T(), contentAnalytics["engagementRate"], float64(0), "Should have engagement rate")

	// Verify audience breakdown exists
	audience, exists := contentAnalytics["audienceBreakdown"].(map[string]interface{})
	require.True(suite.T(), exists, "Should have audience breakdown")

	countries, exists := audience["countries"].(map[string]interface{})
	require.True(suite.T(), exists, "Should have countries breakdown")
	assert.Greater(suite.T(), countries["TH"], float64(0), "Should have Thai viewers")
}

// Test 4.5: Content Moderation and Management API
func (suite *Journey04ContentAPISuite) TestContentModerationAPI() {
	contentID := suite.createPublishedContent()

	creatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.creator.AccessToken,
	}

	// Step 1: PUT /api/v1/content/{id} - Update content
	updateReq := map[string]interface{}{
		"title":       "Traditional Filipino Dance Tutorial - Tinikling (Updated)",
		"description": "Learn the beautiful and traditional Tinikling dance step by step! Updated with new tips and cultural insights.",
		"tags":        []string{"dance", "tutorial", "philippines", "traditional", "culture", "tinikling", "bamboo", "updated"},
	}

	var statusCode int
	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/content/%s", contentID), updateReq, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Content update should succeed")

	// Step 2: PUT /api/v1/content/{id}/visibility - Change visibility
	visibilityReq := map[string]string{
		"visibility": "unlisted",
		"reason":     "Making temporary adjustments",
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/content/%s/visibility", contentID), visibilityReq, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Visibility update should succeed")

	// Step 3: GET /api/v1/content/{id} - Verify update
	verifyResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve updated content")

	var updatedContent ContentResponse
	err := json.Unmarshal(verifyResp, &updatedContent)
	require.NoError(suite.T(), err, "Should parse updated content")

	assert.Contains(suite.T(), updatedContent.Title, "Updated", "Title should contain 'Updated'")
	assert.Equal(suite.T(), "unlisted", updatedContent.Visibility, "Visibility should be unlisted")
	assert.Contains(suite.T(), updatedContent.Tags, "updated", "Should contain 'updated' tag")

	// Step 4: DELETE /api/v1/content/{id} - Archive content
	_, statusCode = suite.makeAPICall("DELETE",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Content archival should succeed")

	// Step 5: GET /api/v1/content/{id} - Verify archived
	archivedResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should still retrieve archived content for creator")

	var archivedContent ContentResponse
	err = json.Unmarshal(archivedResp, &archivedContent)
	require.NoError(suite.T(), err, "Should parse archived content")
	assert.Equal(suite.T(), "archived", archivedContent.Status, "Status should be archived")

	// Step 6: Verify public cannot access archived content
	viewerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.viewer.AccessToken,
	}

	_, statusCode = suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/content/%s", contentID), nil, viewerHeaders)
	assert.Equal(suite.T(), 404, statusCode, "Archived content should not be accessible to public")
}

// Helper methods
func (suite *Journey04ContentAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(suite.ctx, method, suite.baseURL+endpoint, reqBody)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			if key != "Content-Type" { // Don't override content type for multipart
				req.Header.Set(key, value)
			}
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey04ContentAPISuite) uploadMediaFile(endpoint, contentType, filename string, fileData []byte, headers map[string]string) ([]byte, int) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(suite.T(), err)
	part.Write(fileData)

	writer.Close()

	req, err := http.NewRequestWithContext(suite.ctx, "POST", suite.baseURL+endpoint, &b)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if headers != nil {
		for key, value := range headers {
			if key != "Content-Type" {
				req.Header.Set(key, value)
			}
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey04ContentAPISuite) createAuthenticatedUser(email, country, language string) *AuthenticatedUser {
	regReq := map[string]interface{}{
		"email":     email,
		"password":  "SecurePass123!",
		"firstName": "Test",
		"lastName":  "User",
		"country":   country,
		"language":  language,
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	require.Equal(suite.T(), 201, statusCode)

	var regResult map[string]interface{}
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err)

	verifyReq := map[string]string{
		"userId": regResult["userId"].(string),
		"code":   regResult["verifyCode"].(string),
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	require.Equal(suite.T(), 200, statusCode)

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err)

	return &AuthenticatedUser{
		UserID:       regResult["userId"].(string),
		Email:        email,
		AccessToken:  verifyResult["accessToken"].(string),
		RefreshToken: verifyResult["refreshToken"].(string),
		Country:      country,
		Language:     language,
	}
}

func (suite *Journey04ContentAPISuite) createPublishedContent() string {
	contentReq := CreateContentRequest{
		Type:        "video",
		Title:       "Test Filipino Dance Tutorial",
		Description: "Test content for integration testing",
		Tags:        []string{"dance", "tutorial", "test"},
		Category:    "education",
		Language:    "tl",
		Visibility:  "public",
		Monetization: MonetizationSettings{
			Enabled:        true,
			TipJarEnabled:  true,
			AdRevenueShare: true,
			Currency:       "PHP",
		},
		Location: LocationData{
			Country: "PH",
			City:    "Manila",
		},
		VideoData: &VideoUploadData{
			Duration:   300,
			Resolution: "1080p",
			FileSize:   25000000,
			Format:     "mp4",
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.creator.AccessToken,
	}

	contentResp, statusCode := suite.makeAPICall("POST", "/api/v1/content", contentReq, headers)
	require.Equal(suite.T(), 201, statusCode)

	var content ContentResponse
	err := json.Unmarshal(contentResp, &content)
	require.NoError(suite.T(), err)

	// Upload video and wait for processing
	videoData := suite.generateTestVideoData()
	suite.uploadMediaFile(
		fmt.Sprintf("/api/v1/content/%s/upload", content.ID),
		"video/mp4",
		"test_video.mp4",
		videoData,
		headers,
	)

	// Wait for processing
	time.Sleep(3 * time.Second)

	return content.ID
}

func (suite *Journey04ContentAPISuite) simulateEngagement(contentID string) {
	viewerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.viewer.AccessToken,
	}

	// View content
	viewReq := ContentEngagementRequest{
		ContentID: contentID,
		Action:    "view",
		Duration:  200,
		Device:    "mobile",
	}
	suite.makeAPICall("POST", fmt.Sprintf("/api/v1/content/%s/engagement", contentID), viewReq, viewerHeaders)

	// Like content
	likeReq := ContentEngagementRequest{
		ContentID: contentID,
		Action:    "like",
	}
	suite.makeAPICall("POST", fmt.Sprintf("/api/v1/content/%s/engagement", contentID), likeReq, viewerHeaders)

	// Add comment
	commentReq := CreateCommentRequest{
		ContentID: contentID,
		Text:      "Great content!",
		Language:  "en",
	}
	suite.makeAPICall("POST", fmt.Sprintf("/api/v1/content/%s/comments", contentID), commentReq, viewerHeaders)
}

func (suite *Journey04ContentAPISuite) generateTestVideoData() []byte {
	// Generate minimal test video data (fake MP4 header)
	return []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, // ftypes box header
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x02, 0x00, // isom major brand
		0x69, 0x73, 0x6F, 0x6D, 0x69, 0x73, 0x6F, 0x32, // compatible brands
		0x61, 0x76, 0x63, 0x31, 0x6D, 0x70, 0x34, 0x31, // more brands
	}
}

func TestJourney04ContentAPISuite(t *testing.T) {
	suite.Run(t, new(Journey04ContentAPISuite))
}
