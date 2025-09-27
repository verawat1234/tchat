// Journey 6: Social Media & Community API Integration Tests
// Tests all API endpoints involved in social interactions, communities, and user relationships

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AuthenticatedUser represents an authenticated user session
// Note: AuthenticatedUser is now defined in types.go

type Journey06SocialCommunityAPISuite struct {
	suite.Suite
	baseURL     string
	httpClient  *http.Client
	user1       *AuthenticatedUser
	user2       *AuthenticatedUser
	user3       *AuthenticatedUser
	community   *CommunityInfo
	ctx         context.Context
}

type AuthenticatedSocialUser struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
	Profile     *UserProfile `json:"profile"`
}

type UserProfile struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Avatar      string `json:"avatar"`
	Interests   []string `json:"interests"`
}

type CommunityInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // public, private, restricted
	Category    string `json:"category"`
	Region      string `json:"region"`
	CreatorID   string `json:"creatorId"`
}

type FollowRequest struct {
	FollowerID  string `json:"followerId"`
	FollowingID string `json:"followingId"`
	Source      string `json:"source"` // discovery, suggestion, search
}

type CommunityCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Category    string   `json:"category"`
	Region      string   `json:"region"`
	Tags        []string `json:"tags"`
	Rules       []string `json:"rules"`
	Avatar      string   `json:"avatar,omitempty"`
	Banner      string   `json:"banner,omitempty"`
}

type PostRequest struct {
	CommunityID string                 `json:"communityId"`
	Content     string                 `json:"content"`
	Type        string                 `json:"type"` // text, image, video, link, poll
	Metadata    map[string]interface{} `json:"metadata"`
	Tags        []string               `json:"tags,omitempty"`
	Visibility  string                 `json:"visibility"` // public, members, private
}

type CommentRequest struct {
	PostID   string `json:"postId"`
	Content  string `json:"content"`
	ParentID string `json:"parentId,omitempty"` // for replies
}

type ReactionRequest struct {
	TargetID   string `json:"targetId"` // post or comment ID
	TargetType string `json:"targetType"` // post, comment
	Type       string `json:"type"` // like, love, laugh, angry, sad, wow
}

func (suite *Journey06SocialCommunityAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081" // API Gateway
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.ctx = context.Background()

	// Create test users for social interactions
	suite.user1 = suite.createTestUser("maya.social@tchat.com", "Maya", "Chen", "SG")
	suite.user2 = suite.createTestUser("rajesh.community@tchat.com", "Rajesh", "Kumar", "IN")
	suite.user3 = suite.createTestUser("nina.connect@tchat.com", "Nina", "Sari", "ID")
}

// Test 6.1: User Profile Management and Social Setup
func (suite *Journey06SocialCommunityAPISuite) TestUserProfileSocialSetup() {
	headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	// Step 1: PUT /api/v1/profiles/{userId}/social - Update social profile
	profileUpdateReq := map[string]interface{}{
		"displayName": "Maya | Tech Enthusiast",
		"bio":         "Software engineer passionate about Southeast Asian tech innovation 🚀",
		"location":    "Singapore",
		"interests":   []string{"technology", "startup", "innovation", "sea", "programming", "design"},
		"socialLinks": map[string]interface{}{
			"linkedin": "https://linkedin.com/in/mayachen",
			"twitter":  "@mayachen_tech",
			"github":   "mayachen",
		},
		"preferences": map[string]interface{}{
			"discoverability": "public",
			"messaging":       "followers",
			"notifications":   "active",
		},
	}

	_, statusCode := suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/profiles/%s/social", suite.user1.UserID), profileUpdateReq, headers)
	assert.Equal(suite.T(), 200, statusCode, "Social profile update should succeed")

	// Step 2: GET /api/v1/profiles/{userId}/social - Verify profile update
	profileResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/profiles/%s/social", suite.user1.UserID), nil, headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve updated social profile")

	var profile map[string]interface{}
	err := json.Unmarshal(profileResp, &profile)
	require.NoError(suite.T(), err, "Should parse profile response")
	assert.Equal(suite.T(), "Maya | Tech Enthusiast", profile["displayName"])
	assert.Contains(suite.T(), profile["interests"], "technology")

	// Step 3: POST /api/v1/social/interests - Update interests based on activity
	interestsReq := map[string]interface{}{
		"interests": []string{"fintech", "blockchain", "ai", "mobile-development"},
		"source":    "activity_based",
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/social/interests", interestsReq, headers)
	assert.Equal(suite.T(), 200, statusCode, "Interest update should succeed")
}

// Test 6.2: User Discovery and Following System
func (suite *Journey06SocialCommunityAPISuite) TestUserDiscoveryAndFollowing() {
	user1Headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}
	user2Headers := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	// Step 1: GET /api/v1/social/discover/users - Discover users with similar interests
	discoverResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/social/discover/users?region=SEA&interests=technology,startup", nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "User discovery should succeed")

	var discoveredUsers map[string]interface{}
	err := json.Unmarshal(discoverResp, &discoveredUsers)
	require.NoError(suite.T(), err, "Should parse discovery response")
	assert.Greater(suite.T(), len(discoveredUsers["users"].([]interface{})), 0, "Should discover users")

	// Step 2: POST /api/v1/social/follow - Follow another user
	followReq := FollowRequest{
		FollowerID:  suite.user1.UserID,
		FollowingID: suite.user2.UserID,
		Source:      "discovery",
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/social/follow", followReq, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Follow request should succeed")

	// Step 3: GET /api/v1/social/followers/{userId} - Check followers
	followersResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/social/followers/%s", suite.user2.UserID), nil, user2Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve followers")

	var followers map[string]interface{}
	err = json.Unmarshal(followersResp, &followers)
	require.NoError(suite.T(), err, "Should parse followers response")
	assert.Greater(suite.T(), len(followers["followers"].([]interface{})), 0, "Should have followers")

	// Step 4: POST /api/v1/social/follow - Follow back (mutual connection)
	followBackReq := FollowRequest{
		FollowerID:  suite.user2.UserID,
		FollowingID: suite.user1.UserID,
		Source:      "follow_back",
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/social/follow", followBackReq, user2Headers)
	assert.Equal(suite.T(), 200, statusCode, "Follow back should succeed")

	// Step 5: GET /api/v1/social/connections/{userId} - Check mutual connections
	connectionsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/social/connections/%s", suite.user1.UserID), nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve connections")

	var connections map[string]interface{}
	err = json.Unmarshal(connectionsResp, &connections)
	require.NoError(suite.T(), err, "Should parse connections response")
	assert.Greater(suite.T(), connections["mutualCount"], float64(0), "Should have mutual connections")
}

// Test 6.3: Community Creation and Management
func (suite *Journey06SocialCommunityAPISuite) TestCommunityCreationAndManagement() {
	creatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	// Step 1: POST /api/v1/communities - Create a community
	communityReq := CommunityCreateRequest{
		Name:        "Southeast Asian Tech Innovators",
		Description: "A vibrant community for tech professionals, entrepreneurs, and innovators across Southeast Asia. Share insights, collaborate on projects, and build the future of tech in SEA! 🌏",
		Type:        "public",
		Category:    "technology",
		Region:      "SEA",
		Tags:        []string{"technology", "startup", "innovation", "sea", "networking", "collaboration"},
		Rules: []string{
			"Be respectful and professional",
			"Keep discussions relevant to technology and innovation",
			"No spam or self-promotion without context",
			"Share knowledge and help others learn",
			"Celebrate diversity and inclusion",
		},
		Avatar: "https://example.com/community-avatar.jpg",
		Banner: "https://example.com/community-banner.jpg",
	}

	communityResp, statusCode := suite.makeAPICall("POST", "/api/v1/communities", communityReq, creatorHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Community creation should succeed")

	var communityResult map[string]interface{}
	err := json.Unmarshal(communityResp, &communityResult)
	require.NoError(suite.T(), err, "Should parse community response")

	communityID, ok := communityResult["id"].(string)
	assert.True(suite.T(), ok, "Expected id in community response: %+v", communityResult)
	assert.NotEmpty(suite.T(), communityID, "Should return community ID")

	suite.community = &CommunityInfo{
		ID:          communityID,
		Name:        communityReq.Name,
		Description: communityReq.Description,
		Type:        communityReq.Type,
		Category:    communityReq.Category,
		Region:      communityReq.Region,
		CreatorID:   suite.user1.UserID,
	}

	// Step 2: POST /api/v1/communities/{id}/join - Join community
	user2Headers := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	joinReq := map[string]interface{}{
		"reason": "Excited to connect with fellow tech innovators in SEA!",
		"source": "discovery",
	}

	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/communities/%s/join", communityID), joinReq, user2Headers)
	assert.Equal(suite.T(), 200, statusCode, "Community join should succeed")

	// Step 3: GET /api/v1/communities/{id}/members - Check members
	membersResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/communities/%s/members", communityID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve community members")

	var members map[string]interface{}
	err = json.Unmarshal(membersResp, &members)
	require.NoError(suite.T(), err, "Should parse members response")
	assert.Greater(suite.T(), len(members["members"].([]interface{})), 1, "Should have multiple members")

	// Step 4: PUT /api/v1/communities/{id}/roles - Assign moderator role
	roleReq := map[string]interface{}{
		"userId": suite.user2.UserID,
		"role":   "moderator",
		"reason": "Active contributor with great engagement",
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/communities/%s/roles", communityID), roleReq, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Role assignment should succeed")
}

// Test 6.4: Community Posts and Discussions
func (suite *Journey06SocialCommunityAPISuite) TestCommunityPostsAndDiscussions() {
	user1Headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}
	user2Headers := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	communityID := suite.community.ID

	// Step 1: POST /api/v1/communities/{id}/posts - Create a community post
	postReq := PostRequest{
		CommunityID: communityID,
		Content:     "🚀 Exciting news! Southeast Asia's tech ecosystem is booming with innovation! Just saw amazing demos from startups across Singapore, Thailand, and Indonesia. What trends are you seeing in your region? #SEATech #Innovation",
		Type:        "text",
		Metadata: map[string]interface{}{
			"location": "Singapore",
			"mood":     "excited",
			"topics":   []string{"startup", "innovation", "trends"},
		},
		Tags:       []string{"discussion", "trends", "startup", "sea"},
		Visibility: "public",
	}

	postResp, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/communities/%s/posts", communityID), postReq, user1Headers)
	assert.Equal(suite.T(), 201, statusCode, "Post creation should succeed")

	var postResult map[string]interface{}
	err := json.Unmarshal(postResp, &postResult)
	require.NoError(suite.T(), err, "Should parse post response")

	postID, ok := postResult["id"].(string)
	assert.True(suite.T(), ok, "Expected id in post response: %+v", postResult)
	assert.NotEmpty(suite.T(), postID, "Should return post ID")

	// Step 2: POST /api/v1/posts/{id}/reactions - Add reaction to post
	reactionReq := ReactionRequest{
		TargetID:   postID,
		TargetType: "post",
		Type:       "love",
	}

	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/posts/%s/reactions", postID), reactionReq, user2Headers)
	assert.Equal(suite.T(), 200, statusCode, "Post reaction should succeed")

	// Step 3: POST /api/v1/posts/{id}/comments - Comment on post
	commentReq := CommentRequest{
		PostID:  postID,
		Content: "Absolutely agree! The fintech scene in Thailand is particularly impressive. We're seeing incredible innovation in digital payments and blockchain applications. Looking forward to more collaboration across SEA! 🇹🇭",
	}

	commentResp, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/posts/%s/comments", postID), commentReq, user2Headers)
	assert.Equal(suite.T(), 201, statusCode, "Comment creation should succeed")

	var commentResult map[string]interface{}
	err = json.Unmarshal(commentResp, &commentResult)
	require.NoError(suite.T(), err, "Should parse comment response")

	commentID, ok := commentResult["id"].(string)
	assert.True(suite.T(), ok, "Expected id in comment response: %+v", commentResult)

	// Step 4: POST /api/v1/comments/{id}/reactions - React to comment
	commentReactionReq := ReactionRequest{
		TargetID:   commentID,
		TargetType: "comment",
		Type:       "like",
	}

	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/comments/%s/reactions", commentID), commentReactionReq, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Comment reaction should succeed")

	// Step 5: GET /api/v1/communities/{id}/feed - Get community feed
	feedResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/communities/%s/feed?limit=20", communityID), nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve community feed")

	var feed map[string]interface{}
	err = json.Unmarshal(feedResp, &feed)
	require.NoError(suite.T(), err, "Should parse feed response")
	assert.Greater(suite.T(), len(feed["posts"].([]interface{})), 0, "Should have posts in feed")
}

// Test 6.5: Social Feed and Content Discovery
func (suite *Journey06SocialCommunityAPISuite) TestSocialFeedAndDiscovery() {
	user1Headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	// Step 1: GET /api/v1/social/feed - Get personalized social feed
	feedResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/social/feed?algorithm=personalized&limit=20", nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve personalized feed")

	var feed map[string]interface{}
	err := json.Unmarshal(feedResp, &feed)
	require.NoError(suite.T(), err, "Should parse feed response")
	assert.Contains(suite.T(), feed, "posts", "Feed should contain posts")

	// Step 2: GET /api/v1/social/trending - Get trending content
	trendingResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/social/trending?region=SEA&timeframe=24h", nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve trending content")

	var trending map[string]interface{}
	err = json.Unmarshal(trendingResp, &trending)
	require.NoError(suite.T(), err, "Should parse trending response")
	assert.Contains(suite.T(), trending, "topics", "Should have trending topics")
	assert.Contains(suite.T(), trending, "posts", "Should have trending posts")

	// Step 3: GET /api/v1/communities/discover - Discover communities
	discoverResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/communities/discover?category=technology&region=SEA", nil, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Should discover communities")

	var communities map[string]interface{}
	err = json.Unmarshal(discoverResp, &communities)
	require.NoError(suite.T(), err, "Should parse communities response")
	assert.Greater(suite.T(), len(communities["communities"].([]interface{})), 0, "Should discover communities")

	// Step 4: POST /api/v1/social/share - Share content externally
	shareReq := map[string]interface{}{
		"contentId":   suite.community.ID,
		"contentType": "community",
		"platform":    "external",
		"message":     "Check out this amazing tech community! Great discussions about Southeast Asian innovation 🚀",
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/social/share", shareReq, user1Headers)
	assert.Equal(suite.T(), 200, statusCode, "Content sharing should succeed")
}

// Test 6.6: Social Analytics and Insights
func (suite *Journey06SocialCommunityAPISuite) TestSocialAnalyticsAndInsights() {
	creatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	communityID := suite.community.ID

	// Step 1: GET /api/v1/communities/{id}/analytics - Community analytics
	analyticsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/communities/%s/analytics?period=7d", communityID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve community analytics")

	var analytics map[string]interface{}
	err := json.Unmarshal(analyticsResp, &analytics)
	require.NoError(suite.T(), err, "Should parse analytics response")

	// Verify analytics structure
	assert.Contains(suite.T(), analytics, "engagement", "Should have engagement metrics")
	assert.Contains(suite.T(), analytics, "growth", "Should have growth metrics")
	assert.Contains(suite.T(), analytics, "demographics", "Should have demographic data")

	// Check engagement metrics
	engagement := analytics["engagement"].(map[string]interface{})
	assert.Contains(suite.T(), engagement, "postsCount", "Should have posts count")
	assert.Contains(suite.T(), engagement, "commentsCount", "Should have comments count")
	assert.Contains(suite.T(), engagement, "reactionsCount", "Should have reactions count")
	assert.Contains(suite.T(), engagement, "engagementRate", "Should have engagement rate")

	// Step 2: GET /api/v1/social/profile/{userId}/analytics - User social analytics
	profileAnalyticsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/social/profile/%s/analytics?period=30d", suite.user1.UserID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve profile analytics")

	var profileAnalytics map[string]interface{}
	err = json.Unmarshal(profileAnalyticsResp, &profileAnalytics)
	require.NoError(suite.T(), err, "Should parse profile analytics response")

	// Verify profile analytics structure
	assert.Contains(suite.T(), profileAnalytics, "followers", "Should have follower metrics")
	assert.Contains(suite.T(), profileAnalytics, "engagement", "Should have engagement metrics")
	assert.Contains(suite.T(), profileAnalytics, "reach", "Should have reach metrics")

	// Step 3: GET /api/v1/communities/{id}/insights - AI-powered community insights
	insightsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/communities/%s/insights", communityID), nil, creatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve community insights")

	var insights map[string]interface{}
	err = json.Unmarshal(insightsResp, &insights)
	require.NoError(suite.T(), err, "Should parse insights response")

	// Verify insights structure
	assert.Contains(suite.T(), insights, "topTopics", "Should have top topics")
	assert.Contains(suite.T(), insights, "sentimentAnalysis", "Should have sentiment analysis")
	assert.Contains(suite.T(), insights, "recommendations", "Should have recommendations")
	assert.Contains(suite.T(), insights, "regionBreakdown", "Should have regional breakdown")

	// Check regional insights for SEA focus
	regionBreakdown := insights["regionBreakdown"].(map[string]interface{})
	seaCountries := []string{"SG", "TH", "ID", "PH", "MY", "VN"}
	for _, country := range seaCountries {
		if _, exists := regionBreakdown[country]; exists {
			assert.GreaterOrEqual(suite.T(), regionBreakdown[country], float64(0),
				fmt.Sprintf("Should have data for %s", country))
		}
	}
}

// Test 6.7: Social Moderation and Safety
func (suite *Journey06SocialCommunityAPISuite) TestSocialModerationAndSafety() {
	moderatorHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken, // user2 is moderator
	}
	memberHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user3.AccessToken,
	}

	communityID := suite.community.ID

	// Step 1: POST /api/v1/social/report - Report inappropriate content
	reportReq := map[string]interface{}{
		"targetId":   communityID,
		"targetType": "community",
		"reason":     "spam",
		"details":    "Test report for moderation system verification",
		"reporter":   suite.user3.UserID,
	}

	_, statusCode := suite.makeAPICall("POST", "/api/v1/social/report", reportReq, memberHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Content report should succeed")

	// Step 2: GET /api/v1/communities/{id}/reports - Get moderation reports
	reportsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/communities/%s/reports?status=pending", communityID), nil, moderatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve moderation reports")

	var reports map[string]interface{}
	err := json.Unmarshal(reportsResp, &reports)
	require.NoError(suite.T(), err, "Should parse reports response")
	assert.Greater(suite.T(), len(reports["reports"].([]interface{})), 0, "Should have reports")

	// Step 3: PUT /api/v1/communities/{id}/moderation - Take moderation action
	moderationReq := map[string]interface{}{
		"action":     "warning",
		"targetId":   communityID,
		"targetType": "community",
		"reason":     "Test moderation action",
		"moderator":  suite.user2.UserID,
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/communities/%s/moderation", communityID), moderationReq, moderatorHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Moderation action should succeed")

	// Step 4: GET /api/v1/social/safety/guidelines - Get community guidelines
	guidelinesResp, statusCode := suite.makeAPICall("GET", "/api/v1/social/safety/guidelines", nil, memberHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve safety guidelines")

	var guidelines map[string]interface{}
	err = json.Unmarshal(guidelinesResp, &guidelines)
	require.NoError(suite.T(), err, "Should parse guidelines response")
	assert.Contains(suite.T(), guidelines, "rules", "Should have community rules")
	assert.Contains(suite.T(), guidelines, "reporting", "Should have reporting guidelines")
}

// Helper methods
func (suite *Journey06SocialCommunityAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
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
			req.Header.Set(key, value)
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey06SocialCommunityAPISuite) createTestUser(email, firstName, lastName, country string) *AuthenticatedUser {
	regReq := map[string]interface{}{
		"email":     email,
		"password":  "SecurePass123!",
		"firstName": firstName,
		"lastName":  lastName,
		"country":   country,
		"language":  "en",
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	require.Equal(suite.T(), 201, statusCode)

	var regResult map[string]interface{}
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err)

	userID, ok := regResult["userId"].(string)
	assert.True(suite.T(), ok, "Expected userId in registration response: %+v", regResult)

	// Auto-verify for testing
	verifyReq := map[string]string{
		"userId": userID,
		"code":   func() string {
			code, ok := regResult["verifyCode"].(string)
			assert.True(suite.T(), ok, "Expected verifyCode in registration response: %+v", regResult)
			return code
		}(),
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	require.Equal(suite.T(), 200, statusCode)

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err)

	return &AuthenticatedUser{
		UserID:      userID,
		AccessToken: func() string {
			token, ok := verifyResult["accessToken"].(string)
			assert.True(suite.T(), ok, "Expected accessToken in verify response: %+v", verifyResult)
			return token
		}(),
		Email:       email,
	}
}

func TestJourney06SocialCommunityAPISuite(t *testing.T) {
	suite.Run(t, new(Journey06SocialCommunityAPISuite))
}
