package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
)

// StreamHandler handles HTTP requests for stream functionality
type StreamHandler struct {
	categoryService *services.StreamCategoryService
	contentService  *services.StreamContentService
	sessionService  *services.StreamSessionService
	purchaseService *services.StreamPurchaseService
}

// NewStreamHandler creates a new stream handler
func NewStreamHandler(
	categoryService *services.StreamCategoryService,
	contentService *services.StreamContentService,
	sessionService *services.StreamSessionService,
	purchaseService *services.StreamPurchaseService,
) *StreamHandler {
	return &StreamHandler{
		categoryService: categoryService,
		contentService:  contentService,
		sessionService:  sessionService,
		purchaseService: purchaseService,
	}
}

// GetStreamCategories handles GET /api/v1/stream/categories
func (h *StreamHandler) GetStreamCategories(c *gin.Context) {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve categories",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
		"success":    true,
	})
}

// GetStreamCategoryDetail handles GET /api/v1/stream/categories/:id
func (h *StreamHandler) GetStreamCategoryDetail(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is required",
		})
		return
	}

	category, err := h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Category not found",
		})
		return
	}

	// Get category statistics
	stats, err := h.categoryService.GetCategoryStats(categoryID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to get category stats: %v\n", err)
		stats = map[string]interface{}{}
	}

	// Get subtabs for the category
	subtabs, err := h.categoryService.GetSubtabsByCategoryID(categoryID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to get category subtabs: %v\n", err)
		subtabs = []models.StreamSubtab{}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"subtabs":  subtabs,
		"stats":    stats,
		"success":  true,
	})
}

// GetStreamContent handles GET /api/v1/stream/content
func (h *StreamHandler) GetStreamContent(c *gin.Context) {
	categoryID := c.Query("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is required",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse optional subtab filter
	var subtabID *string
	if subtab := c.Query("subtabId"); subtab != "" {
		subtabID = &subtab
	}

	response, err := h.contentService.GetContentByCategory(categoryID, page, limit, subtabID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve content",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"items":   response.Items,
		"total":   response.Total,
		"hasMore": response.HasMore,
		"success": true,
	})
}

// GetStreamFeatured handles GET /api/v1/stream/featured
func (h *StreamHandler) GetStreamFeatured(c *gin.Context) {
	categoryID := c.Query("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.contentService.GetFeaturedContent(categoryID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve featured content",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"items":   response.Items,
		"total":   response.Total,
		"hasMore": response.HasMore,
		"success": true,
	})
}

// GetStreamContentDetail handles GET /api/v1/stream/content/:id
func (h *StreamHandler) GetStreamContentDetail(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Content ID is required",
		})
		return
	}

	content, err := h.contentService.GetContentByID(contentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Content not found",
		})
		return
	}

	// Track content view if user is authenticated
	userID := h.getUserIDFromContext(c)
	sessionID := h.getSessionIDFromContext(c)
	if userID != "" && sessionID != "" {
		_, err := h.sessionService.TrackContentView(userID, contentID, sessionID)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to track content view: %v\n", err)
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"content": content,
		"success": true,
	})
}

// SearchStreamContent handles GET /api/v1/stream/search
func (h *StreamHandler) SearchStreamContent(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Search query is required",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse optional category filter
	var categoryID *string
	if category := c.Query("categoryId"); category != "" {
		categoryID = &category
	}

	response, err := h.contentService.SearchContent(query, categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to search content",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"items":   response.Items,
		"total":   response.Total,
		"hasMore": response.HasMore,
		"success": true,
	})
}

// PostStreamContentPurchase handles POST /api/v1/stream/content/purchase
func (h *StreamHandler) PostStreamContentPurchase(c *gin.Context) {
	// Check authentication
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	var req models.StreamPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	response, err := h.purchaseService.ProcessContentPurchase(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

// GetUserNavigationState handles GET /api/v1/stream/navigation
func (h *StreamHandler) GetUserNavigationState(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	state, err := h.sessionService.GetUserNavigationState(userID)
	if err != nil {
		// Create default state if none exists
		sessionID := h.getSessionIDFromContext(c)
		if sessionID == "" {
			sessionID = fmt.Sprintf("session_%d", c.GetHeader("X-Request-ID"))
		}

		state, err = h.sessionService.UpdateUserNavigationState(userID, sessionID, "books", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get navigation state",
			})
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"state":   state,
		"success": true,
	})
}

// UpdateUserNavigationState handles PUT /api/v1/stream/navigation
func (h *StreamHandler) UpdateUserNavigationState(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	var req struct {
		CategoryID string  `json:"categoryId" binding:"required"`
		SubtabID   *string `json:"subtabId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	sessionID := h.getSessionIDFromContext(c)
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", c.GetHeader("X-Request-ID"))
	}

	state, err := h.sessionService.UpdateUserNavigationState(userID, sessionID, req.CategoryID, req.SubtabID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update navigation state",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"state":   state,
		"success": true,
	})
}

// UpdateContentViewProgress handles PUT /api/v1/stream/content/:id/progress
func (h *StreamHandler) UpdateContentViewProgress(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Content ID is required",
		})
		return
	}

	var req struct {
		Progress float64 `json:"progress" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	sessionID := h.getSessionIDFromContext(c)
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", c.GetHeader("X-Request-ID"))
	}

	err := h.sessionService.UpdateContentViewProgress(userID, contentID, sessionID, req.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update view progress",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Progress updated successfully",
	})
}

// GetUserPreferences handles GET /api/v1/stream/preferences
func (h *StreamHandler) GetUserPreferences(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	prefs, err := h.sessionService.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get user preferences",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"preferences": prefs,
		"success":     true,
	})
}

// UpdateUserPreferences handles PUT /api/v1/stream/preferences
func (h *StreamHandler) UpdateUserPreferences(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication required",
		})
		return
	}

	var prefs models.StreamUserPreference
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	if err := h.sessionService.UpdateUserPreferences(userID, &prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update user preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Preferences updated successfully",
	})
}

// Helper methods for extracting user and session information from context
func (h *StreamHandler) getUserIDFromContext(c *gin.Context) string {
	// In a real implementation, this would extract user ID from JWT token or session
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return uid
		}
	}

	// For testing purposes, check header
	return c.GetHeader("X-User-ID")
}

func (h *StreamHandler) getSessionIDFromContext(c *gin.Context) string {
	// In a real implementation, this would extract session ID from headers or generate one
	if sessionID := c.GetHeader("X-Session-ID"); sessionID != "" {
		return sessionID
	}

	// Generate a basic session ID based on request context
	return fmt.Sprintf("session_%s_%s", h.getUserIDFromContext(c), c.GetHeader("X-Request-ID"))
}