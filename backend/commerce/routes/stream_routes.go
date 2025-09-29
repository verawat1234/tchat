package routes

import (
	"github.com/gin-gonic/gin"
	"tchat.dev/commerce/handlers"
)

// RegisterStreamRoutes registers all stream-related API routes
func RegisterStreamRoutes(router *gin.RouterGroup, streamHandler *handlers.StreamHandler) {
	// Stream categories routes
	streamRoutes := router.Group("/stream")
	{
		// Categories
		streamRoutes.GET("/categories", streamHandler.GetStreamCategories)
		streamRoutes.GET("/categories/:id", streamHandler.GetStreamCategoryDetail)

		// Content browsing
		streamRoutes.GET("/content", streamHandler.GetStreamContent)
		streamRoutes.GET("/content/:id", streamHandler.GetStreamContentDetail)
		streamRoutes.GET("/featured", streamHandler.GetStreamFeatured)
		streamRoutes.GET("/search", streamHandler.SearchStreamContent)

		// Content purchasing
		streamRoutes.POST("/content/purchase", streamHandler.PostStreamContentPurchase)

		// User navigation and session management (requires authentication)
		authRoutes := streamRoutes.Group("/")
		// authRoutes.Use(middleware.AuthRequired()) // Uncomment when auth middleware is available
		{
			authRoutes.GET("/navigation", streamHandler.GetUserNavigationState)
			authRoutes.PUT("/navigation", streamHandler.UpdateUserNavigationState)
			authRoutes.PUT("/content/:id/progress", streamHandler.UpdateContentViewProgress)
			authRoutes.GET("/preferences", streamHandler.GetUserPreferences)
			authRoutes.PUT("/preferences", streamHandler.UpdateUserPreferences)
		}
	}
}

// RegisterStreamRoutesV1 registers stream routes under /api/v1 prefix
func RegisterStreamRoutesV1(router *gin.Engine, streamHandler *handlers.StreamHandler) {
	apiV1 := router.Group("/api/v1")
	RegisterStreamRoutes(apiV1, streamHandler)
}

// RegisterStreamHealthCheck registers health check endpoints for stream service
func RegisterStreamHealthCheck(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "stream",
			"version": "1.0.0",
		})
	})

	router.GET("/health/stream", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "stream",
			"component": "stream-content-management",
			"version":   "1.0.0",
			"endpoints": gin.H{
				"categories":   "/api/v1/stream/categories",
				"content":      "/api/v1/stream/content",
				"featured":     "/api/v1/stream/featured",
				"search":       "/api/v1/stream/search",
				"purchase":     "/api/v1/stream/content/purchase",
				"navigation":   "/api/v1/stream/navigation",
				"preferences":  "/api/v1/stream/preferences",
			},
		})
	})
}

// GetStreamAPISpec returns OpenAPI specification for stream endpoints
func GetStreamAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Stream Content API",
			"description": "API for managing stream content categories, content items, and user interactions",
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url":         "/api/v1",
				"description": "API v1 endpoint",
			},
		},
		"paths": map[string]interface{}{
			"/stream/categories": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get all stream categories",
					"description": "Retrieves all active stream categories ordered by display order",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"categories": map[string]interface{}{
												"type": "array",
												"items": map[string]interface{}{
													"$ref": "#/components/schemas/StreamCategory",
												},
											},
											"total":   map[string]interface{}{"type": "integer"},
											"success": map[string]interface{}{"type": "boolean"},
										},
									},
								},
							},
						},
					},
				},
			},
			"/stream/categories/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get category details",
					"description": "Retrieves detailed information about a specific category including subtabs and statistics",
					"parameters": []map[string]interface{}{
						{
							"name":        "id",
							"in":          "path",
							"required":    true,
							"description": "Category ID",
							"schema":      map[string]interface{}{"type": "string"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"category": map[string]interface{}{"$ref": "#/components/schemas/StreamCategory"},
											"subtabs": map[string]interface{}{
												"type": "array",
												"items": map[string]interface{}{
													"$ref": "#/components/schemas/StreamSubtab",
												},
											},
											"stats":   map[string]interface{}{"type": "object"},
											"success": map[string]interface{}{"type": "boolean"},
										},
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Category not found",
						},
					},
				},
			},
			"/stream/content": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get content for category",
					"description": "Retrieves content items for a specific category with pagination and optional subtab filtering",
					"parameters": []map[string]interface{}{
						{
							"name":        "categoryId",
							"in":          "query",
							"required":    true,
							"description": "Category ID",
							"schema":      map[string]interface{}{"type": "string"},
						},
						{
							"name":        "subtabId",
							"in":          "query",
							"required":    false,
							"description": "Optional subtab ID for filtering",
							"schema":      map[string]interface{}{"type": "string"},
						},
						{
							"name":        "page",
							"in":          "query",
							"required":    false,
							"description": "Page number (default: 1)",
							"schema":      map[string]interface{}{"type": "integer", "minimum": 1},
						},
						{
							"name":        "limit",
							"in":          "query",
							"required":    false,
							"description": "Items per page (default: 20, max: 100)",
							"schema":      map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 100},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ContentResponse",
									},
								},
							},
						},
					},
				},
			},
			"/stream/featured": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get featured content",
					"description": "Retrieves featured content for a specific category",
					"parameters": []map[string]interface{}{
						{
							"name":        "categoryId",
							"in":          "query",
							"required":    true,
							"description": "Category ID",
							"schema":      map[string]interface{}{"type": "string"},
						},
						{
							"name":        "limit",
							"in":          "query",
							"required":    false,
							"description": "Number of featured items (default: 10, max: 50)",
							"schema":      map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 50},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/FeaturedResponse",
									},
								},
							},
						},
					},
				},
			},
			"/stream/search": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Search content",
					"description": "Searches for content by title, description, or tags",
					"parameters": []map[string]interface{}{
						{
							"name":        "q",
							"in":          "query",
							"required":    true,
							"description": "Search query",
							"schema":      map[string]interface{}{"type": "string"},
						},
						{
							"name":        "categoryId",
							"in":          "query",
							"required":    false,
							"description": "Optional category filter",
							"schema":      map[string]interface{}{"type": "string"},
						},
						{
							"name":        "page",
							"in":          "query",
							"required":    false,
							"description": "Page number (default: 1)",
							"schema":      map[string]interface{}{"type": "integer", "minimum": 1},
						},
						{
							"name":        "limit",
							"in":          "query",
							"required":    false,
							"description": "Items per page (default: 20, max: 100)",
							"schema":      map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 100},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ContentResponse",
									},
								},
							},
						},
					},
				},
			},
			"/stream/content/purchase": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Purchase content",
					"description": "Processes a content purchase request",
					"security": []map[string]interface{}{
						{"bearerAuth": []string{}},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/StreamPurchaseRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Purchase successful",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/StreamPurchaseResponse",
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Invalid request",
						},
						"401": map[string]interface{}{
							"description": "Authentication required",
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"StreamCategory": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id":                     map[string]interface{}{"type": "string"},
						"name":                   map[string]interface{}{"type": "string"},
						"displayOrder":           map[string]interface{}{"type": "integer"},
						"iconName":               map[string]interface{}{"type": "string"},
						"isActive":               map[string]interface{}{"type": "boolean"},
						"featuredContentEnabled": map[string]interface{}{"type": "boolean"},
						"createdAt":              map[string]interface{}{"type": "string", "format": "date-time"},
						"updatedAt":              map[string]interface{}{"type": "string", "format": "date-time"},
					},
				},
				"StreamSubtab": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id":             map[string]interface{}{"type": "string"},
						"categoryId":     map[string]interface{}{"type": "string"},
						"name":           map[string]interface{}{"type": "string"},
						"displayOrder":   map[string]interface{}{"type": "integer"},
						"filterCriteria": map[string]interface{}{"type": "object"},
						"isActive":       map[string]interface{}{"type": "boolean"},
						"createdAt":      map[string]interface{}{"type": "string", "format": "date-time"},
						"updatedAt":      map[string]interface{}{"type": "string", "format": "date-time"},
					},
				},
				"StreamContentItem": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id":                 map[string]interface{}{"type": "string"},
						"categoryId":         map[string]interface{}{"type": "string"},
						"subtabId":           map[string]interface{}{"type": "string"},
						"title":              map[string]interface{}{"type": "string"},
						"description":        map[string]interface{}{"type": "string"},
						"thumbnailUrl":       map[string]interface{}{"type": "string"},
						"contentType":        map[string]interface{}{"type": "string"},
						"price":              map[string]interface{}{"type": "number"},
						"currency":           map[string]interface{}{"type": "string"},
						"availabilityStatus": map[string]interface{}{"type": "string"},
						"isFeatured":         map[string]interface{}{"type": "boolean"},
						"featuredOrder":      map[string]interface{}{"type": "integer"},
						"createdAt":          map[string]interface{}{"type": "string", "format": "date-time"},
						"updatedAt":          map[string]interface{}{"type": "string", "format": "date-time"},
					},
				},
				"ContentResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/StreamContentItem",
							},
						},
						"total":   map[string]interface{}{"type": "integer"},
						"hasMore": map[string]interface{}{"type": "boolean"},
						"success": map[string]interface{}{"type": "boolean"},
					},
				},
				"FeaturedResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/StreamContentItem",
							},
						},
						"total":   map[string]interface{}{"type": "integer"},
						"hasMore": map[string]interface{}{"type": "boolean"},
						"success": map[string]interface{}{"type": "boolean"},
					},
				},
				"StreamPurchaseRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"mediaContentId", "quantity", "mediaLicense"},
					"properties": map[string]interface{}{
						"mediaContentId": map[string]interface{}{"type": "string"},
						"quantity":       map[string]interface{}{"type": "integer", "minimum": 1},
						"mediaLicense":   map[string]interface{}{"type": "string", "enum": []string{"personal", "family", "commercial"}},
						"downloadFormat": map[string]interface{}{"type": "string", "enum": []string{"PDF", "EPUB", "MP3", "MP4", "FLAC"}},
						"cartId":         map[string]interface{}{"type": "string"},
					},
				},
				"StreamPurchaseResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"orderId":     map[string]interface{}{"type": "string"},
						"totalAmount": map[string]interface{}{"type": "number"},
						"currency":    map[string]interface{}{"type": "string"},
						"success":     map[string]interface{}{"type": "boolean"},
						"message":     map[string]interface{}{"type": "string"},
					},
				},
			},
			"securitySchemes": map[string]interface{}{
				"bearerAuth": map[string]interface{}{
					"type":   "http",
					"scheme": "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}
}