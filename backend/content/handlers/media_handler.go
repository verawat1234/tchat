package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tchat.dev/content/services"
)

// MediaHandler handles HTTP requests for media endpoints
type MediaHandler struct {
	service services.MediaService
}

// NewMediaHandler creates a new MediaHandler instance
func NewMediaHandler(service services.MediaService) *MediaHandler {
	return &MediaHandler{
		service: service,
	}
}

// GetCategories handles GET /media/categories
func (h *MediaHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// GetCategory handles GET /media/categories/{categoryId}
func (h *MediaHandler) GetCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Category ID is required",
		})
		return
	}

	category, err := h.service.GetCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		if err.Error() == "category not found: "+categoryID {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve category",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetFeaturedContent handles GET /media/featured
func (h *MediaHandler) GetFeaturedContent(c *gin.Context) {
	// Parse query parameters
	limitStr := c.Query("limit")
	categoryID := c.Query("categoryId")

	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
			if limit > 50 {
				limit = 50 // max limit
			}
		}
	}

	req := &services.GetFeaturedContentRequest{
		CategoryID: categoryID,
		Limit:      limit,
	}

	response, err := h.service.GetFeaturedContent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve featured content",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetContentByCategory handles GET /media/category/{categoryId}/content
func (h *MediaHandler) GetContentByCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Category ID is required",
		})
		return
	}

	// Parse query parameters
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	subtab := c.Query("subtab")

	page := 1
	if pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100 // max limit
			}
		}
	}

	req := &services.GetContentByCategoryRequest{
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		Subtab:     subtab,
	}

	response, err := h.service.GetContentByCategory(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "category not found: "+categoryID {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve content",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMovieSubtabs handles GET /media/movies/subtabs
func (h *MediaHandler) GetMovieSubtabs(c *gin.Context) {
	response, err := h.service.GetMovieSubtabs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve movie subtabs",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchContent handles GET /media/search
func (h *MediaHandler) SearchContent(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Search query parameter 'q' is required",
		})
		return
	}

	// Parse query parameters
	categoryID := c.Query("categoryId")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	if pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100 // max limit
			}
		}
	}

	req := &services.SearchContentRequest{
		Query:      query,
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
	}

	response, err := h.service.SearchContent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to search content",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all media routes with the router
func (h *MediaHandler) RegisterRoutes(router *gin.Engine) {
	// Create API v1 group
	v1 := router.Group("/api/v1")

	// Media routes
	media := v1.Group("/media")
	{
		media.GET("/categories", h.GetCategories)
		media.GET("/categories/:categoryId", h.GetCategory)
		media.GET("/featured", h.GetFeaturedContent)
		media.GET("/category/:categoryId/content", h.GetContentByCategory)
		media.GET("/movies/subtabs", h.GetMovieSubtabs)
		media.GET("/search", h.SearchContent)
	}
}