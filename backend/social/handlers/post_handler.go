package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/shared/responses"
	"tchat/social/models"
	"tchat/social/services"
)

// PostHandler handles post-related social endpoints
type PostHandler struct {
	postService services.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary Create a new social post
// @Description Create a new post in a community or user feed
// @Tags social, posts
// @Accept json
// @Produce json
// @Param request body models.CreatePostRequest true "Post creation request"
// @Success 201 {object} responses.DataResponse{data=models.Post}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	authorID := userID.(uuid.UUID)

	post, err := h.postService.CreatePost(c.Request.Context(), authorID, &req)
	if err != nil {
		if err.Error() == "community not found" {
			responses.BadRequestResponse(c, "Invalid community ID")
			return
		}
		responses.InternalErrorResponse(c, "Failed to create post")
		return
	}

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success:   true,
		Data:      post,
		Timestamp: "2024-09-22T18:30:00Z",
	})
}

// GetPost godoc
// @Summary Get a specific post
// @Description Retrieve a post by its ID
// @Tags social, posts
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Success 200 {object} responses.DataResponse{data=models.Post}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/posts/{postId} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid post ID format")
		return
	}

	post, err := h.postService.GetPost(c.Request.Context(), postID)
	if err != nil {
		if err.Error() == "post not found" {
			responses.NotFoundResponse(c, "Post not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to retrieve post")
		return
	}

	responses.SendDataResponse(c, post)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post's content or metadata
// @Tags social, posts
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Param request body models.UpdatePostRequest true "Post update request"
// @Success 200 {object} responses.DataResponse{data=models.Post}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/posts/{postId} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid post ID format")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	authorID := userID.(uuid.UUID)

	post, err := h.postService.UpdatePost(c.Request.Context(), postID, authorID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			responses.NotFoundResponse(c, "Post not found")
			return
		}
		if err.Error() == "unauthorized" {
			responses.ForbiddenResponse(c, "Not authorized to update this post")
			return
		}
		responses.InternalErrorResponse(c, "Failed to update post")
		return
	}

	responses.SendDataResponse(c, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post (soft delete)
// @Tags social, posts
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/posts/{postId} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid post ID format")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	authorID := userID.(uuid.UUID)

	err = h.postService.DeletePost(c.Request.Context(), postID, authorID)
	if err != nil {
		if err.Error() == "post not found" {
			responses.NotFoundResponse(c, "Post not found")
			return
		}
		if err.Error() == "unauthorized" {
			responses.ForbiddenResponse(c, "Not authorized to delete this post")
			return
		}
		responses.InternalErrorResponse(c, "Failed to delete post")
		return
	}

	responses.SuccessMessageResponse(c, "Post deleted successfully")
}

// AddReaction godoc
// @Summary Add reaction to a post
// @Description Add a reaction (like, love, etc.) to a post
// @Tags social, posts, reactions
// @Accept json
// @Produce json
// @Param request body models.CreateReactionRequest true "Reaction request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/reactions [post]
func (h *PostHandler) AddReaction(c *gin.Context) {
	var req models.CreateReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	reactorID := userID.(uuid.UUID)

	err := h.postService.AddReaction(c.Request.Context(), reactorID, &req)
	if err != nil {
		if err.Error() == "target not found" {
			responses.NotFoundResponse(c, "Target post or comment not found")
			return
		}
		if err.Error() == "already reacted" {
			responses.ConflictResponse(c, "Already reacted to this content")
			return
		}
		responses.InternalErrorResponse(c, "Failed to add reaction")
		return
	}

	responses.SuccessMessageResponse(c, "Reaction added successfully")
}

// RemoveReaction godoc
// @Summary Remove reaction from a post
// @Description Remove a user's reaction from a post or comment
// @Tags social, posts, reactions
// @Accept json
// @Produce json
// @Param targetId path string true "Target ID (post or comment)"
// @Param targetType path string true "Target type (post or comment)"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/reactions/{targetType}/{targetId} [delete]
func (h *PostHandler) RemoveReaction(c *gin.Context) {
	targetIDParam := c.Param("targetId")
	targetType := c.Param("targetType")

	targetID, err := uuid.Parse(targetIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid target ID format")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	reactorID := userID.(uuid.UUID)

	err = h.postService.RemoveReaction(c.Request.Context(), reactorID, targetID, targetType)
	if err != nil {
		if err.Error() == "reaction not found" {
			responses.NotFoundResponse(c, "Reaction not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to remove reaction")
		return
	}

	responses.SuccessMessageResponse(c, "Reaction removed successfully")
}

// CreateComment godoc
// @Summary Create a comment on a post
// @Description Add a comment to a post
// @Tags social, posts, comments
// @Accept json
// @Produce json
// @Param request body models.CreateCommentRequest true "Comment request"
// @Success 201 {object} responses.DataResponse{data=models.Comment}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/comments [post]
func (h *PostHandler) CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	authorID := userID.(uuid.UUID)

	comment, err := h.postService.CreateComment(c.Request.Context(), authorID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			responses.NotFoundResponse(c, "Post not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to create comment")
		return
	}

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success:   true,
		Data:      comment,
		Timestamp: "2024-09-22T18:30:00Z",
	})
}

// GetSocialFeed godoc
// @Summary Get user's personalized social feed
// @Description Retrieve a personalized feed of posts for the user
// @Tags social, feed
// @Accept json
// @Produce json
// @Param algorithm query string false "Feed algorithm (personalized, chronological, trending)"
// @Param limit query int false "Limit results (default: 20)"
// @Param cursor query string false "Pagination cursor"
// @Param region query string false "Region filter"
// @Success 200 {object} responses.DataResponse{data=models.SocialFeed}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/feed [get]
func (h *PostHandler) GetSocialFeed(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	userUUID := userID.(uuid.UUID)

	var req models.SocialFeedRequest
	req.Algorithm = c.DefaultQuery("algorithm", "personalized")
	req.Cursor = c.Query("cursor")
	req.Region = c.Query("region")

	if limitParam := c.Query("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			req.Limit = limit
		} else {
			responses.BadRequestResponse(c, "Invalid limit parameter")
			return
		}
	} else {
		req.Limit = 20
	}

	feed, err := h.postService.GetSocialFeed(c.Request.Context(), userUUID, &req)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve social feed")
		return
	}

	responses.SendDataResponse(c, feed)
}

// GetTrendingContent godoc
// @Summary Get trending social content
// @Description Retrieve trending posts and topics
// @Tags social, trending
// @Accept json
// @Produce json
// @Param region query string false "Region filter"
// @Param timeframe query string false "Timeframe (1h, 24h, 7d)"
// @Param category query string false "Category filter"
// @Param limit query int false "Limit results (default: 20)"
// @Success 200 {object} responses.DataResponse{data=models.TrendingContent}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/trending [get]
func (h *PostHandler) GetTrendingContent(c *gin.Context) {
	var req models.TrendingRequest
	req.Region = c.DefaultQuery("region", "SEA")
	req.Timeframe = c.DefaultQuery("timeframe", "24h")
	req.Category = c.Query("category")

	if limitParam := c.Query("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			req.Limit = limit
		} else {
			responses.BadRequestResponse(c, "Invalid limit parameter")
			return
		}
	} else {
		req.Limit = 20
	}

	trending, err := h.postService.GetTrendingContent(c.Request.Context(), &req)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve trending content")
		return
	}

	responses.SendDataResponse(c, trending)
}

// ShareContent godoc
// @Summary Share content externally
// @Description Share a post or other content to external platforms
// @Tags social, sharing
// @Accept json
// @Produce json
// @Param request body models.ShareRequest true "Share request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/share [post]
func (h *PostHandler) ShareContent(c *gin.Context) {
	var req models.ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	sharerID := userID.(uuid.UUID)

	err := h.postService.ShareContent(c.Request.Context(), sharerID, &req)
	if err != nil {
		if err.Error() == "content not found" {
			responses.NotFoundResponse(c, "Content not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to share content")
		return
	}

	responses.SuccessMessageResponse(c, "Content shared successfully")
}