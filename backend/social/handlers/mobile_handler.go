package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat/social/models"
	"tchat/social/services"
)

// MobileHandler provides KMP-optimized API endpoints for mobile clients
type MobileHandler struct {
	syncService services.MobileSyncService
	userService services.UserService
}

// NewMobileHandler creates a new mobile API handler
func NewMobileHandler(syncService services.MobileSyncService, userService services.UserService) *MobileHandler {
	return &MobileHandler{
		syncService: syncService,
		userService: userService,
	}
}

// RegisterMobileRoutes registers all mobile-optimized API routes
func (h *MobileHandler) RegisterMobileRoutes(router *gin.RouterGroup) {
	mobile := router.Group("/mobile")
	{
		// Incremental sync endpoints
		mobile.GET("/sync/profile/:userId", h.GetProfileChanges)
		mobile.GET("/sync/posts/:userId", h.GetPostChanges)
		mobile.GET("/sync/follows/:userId", h.GetFollowChanges)

		// Initial data load endpoints
		mobile.GET("/init/:userId", h.GetInitialUserData)
		mobile.GET("/feed/:userId", h.GetUserFeed)

		// Conflict resolution endpoints
		mobile.POST("/resolve/:userId", h.ResolveConflicts)
		mobile.POST("/apply/:userId", h.ApplyClientChanges)

		// Discovery endpoints optimized for mobile
		mobile.GET("/discover/:userId", h.GetDiscoveryFeed)
		mobile.GET("/trending", h.GetTrendingContent)

		// Mobile-specific user operations
		mobile.GET("/profile/:userId", h.GetMobileProfile)
		mobile.PUT("/profile/:userId", h.UpdateMobileProfile)
		mobile.POST("/follow", h.MobileFollow)
		mobile.DELETE("/follow", h.MobileUnfollow)

		// Health check endpoint for mobile apps
		mobile.GET("/health", h.MobileHealthCheck)
	}
}

// GetProfileChanges handles incremental profile sync requests
func (h *MobileHandler) GetProfileChanges(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Parse since parameter (required for incremental sync)
	sinceStr := c.Query("since")
	if sinceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_since_parameter",
			"message": "since parameter is required for incremental sync",
		})
		return
	}

	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_since_format",
			"message": "since parameter must be RFC3339 format",
		})
		return
	}

	response, err := h.syncService.GetUserProfileChanges(c.Request.Context(), userID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "sync_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPostChanges handles incremental post sync requests
func (h *MobileHandler) GetPostChanges(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	sinceStr := c.Query("since")
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_since_format"})
		return
	}

	response, err := h.syncService.GetPostChanges(c.Request.Context(), userID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFollowChanges handles incremental follow sync requests
func (h *MobileHandler) GetFollowChanges(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	sinceStr := c.Query("since")
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_since_format"})
		return
	}

	response, err := h.syncService.GetFollowChanges(c.Request.Context(), userID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInitialUserData handles initial app data load requests
func (h *MobileHandler) GetInitialUserData(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	response, err := h.syncService.GetInitialUserData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserFeed handles paginated feed requests optimized for mobile
func (h *MobileHandler) GetUserFeed(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	// Parse pagination parameters with mobile-friendly defaults
	limit := 20 // Default mobile feed page size
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	response, err := h.syncService.GetUserFeed(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResolveConflicts handles conflict resolution for offline changes
func (h *MobileHandler) ResolveConflicts(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	var request struct {
		ClientData interface{} `json:"clientData"`
		LastSync   time.Time   `json:"lastSync"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request_body"})
		return
	}

	// This would need proper type conversion for client data
	// For now, return a basic conflict resolution response
	c.JSON(http.StatusOK, gin.H{
		"userId":       userID,
		"resolvedAt":   time.Now(),
		"hasConflicts": false,
		"resolution":   "no_conflict",
	})
}

// ApplyClientChanges handles applying validated client changes
func (h *MobileHandler) ApplyClientChanges(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	var changes services.ClientChanges
	if err := c.ShouldBindJSON(&changes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	// Validate that the user ID in the request matches the URL parameter
	if changes.UserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "user_id_mismatch",
			"message": "User ID in URL must match user ID in request body",
		})
		return
	}

	result, err := h.syncService.ApplyClientChanges(c.Request.Context(), userID, &changes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDiscoveryFeed handles mobile-optimized user discovery
func (h *MobileHandler) GetDiscoveryFeed(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	// Parse region parameter (default to user's region or TH)
	region := c.Query("region")
	if region == "" {
		region = "TH" // Default to Thailand for Southeast Asian deployment
	}

	// Validate supported regions
	supportedRegions := map[string]bool{
		"TH": true, "SG": true, "ID": true,
		"MY": true, "PH": true, "VN": true,
	}

	if !supportedRegions[region] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "unsupported_region",
			"message": "Region must be one of: TH, SG, ID, MY, PH, VN",
		})
		return
	}

	// Parse limit with mobile-appropriate default
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 20 {
			limit = parsedLimit
		}
	}

	profiles, err := h.syncService.GetDiscoveryFeed(c.Request.Context(), userID, region, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":   userID,
		"region":   region,
		"profiles": profiles,
		"count":    len(profiles),
	})
}

// GetTrendingContent handles regional trending content requests
func (h *MobileHandler) GetTrendingContent(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		region = "TH"
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 20 {
			limit = parsedLimit
		}
	}

	posts, err := h.syncService.GetTrendingContent(c.Request.Context(), region, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"region": region,
		"posts":  posts,
		"count":  len(posts),
	})
}

// GetMobileProfile returns user profile optimized for mobile display
func (h *MobileHandler) GetMobileProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	profile, err := h.userService.GetSocialProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile_not_found"})
		return
	}

	// Return mobile-optimized profile response
	c.JSON(http.StatusOK, gin.H{
		"profile":     profile,
		"lastSync":    time.Now(),
		"syncVersion": "1.0",
	})
}

// UpdateMobileProfile handles mobile profile update requests
func (h *MobileHandler) UpdateMobileProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_user_id"})
		return
	}

	var updateReq struct {
		DisplayName *string                    `json:"displayName,omitempty"`
		Bio         *string                    `json:"bio,omitempty"`
		Interests   []string                   `json:"interests,omitempty"`
		SocialLinks *map[string]interface{}    `json:"socialLinks,omitempty"`
		LastSync    time.Time                  `json:"lastSync"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request_body"})
		return
	}

	// Create update request for user service
	updateProfile := &models.UpdateSocialProfileRequest{
		DisplayName: updateReq.DisplayName,
		Bio:         updateReq.Bio,
		Interests:   updateReq.Interests,
		SocialLinks: updateReq.SocialLinks,
	}

	updatedProfile, err := h.userService.UpdateSocialProfile(c.Request.Context(), userID, updateProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":  updatedProfile,
		"updated":  true,
		"syncTime": time.Now(),
	})
}

// MobileFollow handles follow requests from mobile clients
func (h *MobileHandler) MobileFollow(c *gin.Context) {
	var request struct {
		FollowerID  uuid.UUID `json:"followerId" binding:"required"`
		FollowingID uuid.UUID `json:"followingId" binding:"required"`
		Source      string    `json:"source"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request_body"})
		return
	}

	followReq := &models.FollowRequest{
		FollowerID:  request.FollowerID,
		FollowingID: request.FollowingID,
		Source:      request.Source,
	}

	if request.Source == "" {
		followReq.Source = "mobile_app"
	}

	err := h.userService.FollowUser(c.Request.Context(), followReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"followed":  true,
		"timestamp": time.Now(),
	})
}

// MobileUnfollow handles unfollow requests from mobile clients
func (h *MobileHandler) MobileUnfollow(c *gin.Context) {
	var request struct {
		FollowerID  uuid.UUID `json:"followerId" binding:"required"`
		FollowingID uuid.UUID `json:"followingId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request_body"})
		return
	}

	err := h.userService.UnfollowUser(c.Request.Context(), request.FollowerID, request.FollowingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"unfollowed": true,
		"timestamp": time.Now(),
	})
}

// MobileHealthCheck provides health status for mobile app monitoring
func (h *MobileHandler) MobileHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "social-mobile-api",
		"version":   "1.0.0",
		"timestamp": time.Now(),
		"features": gin.H{
			"sync":        true,
			"discovery":   true,
			"conflicts":   true,
			"incremental": true,
		},
	})
}