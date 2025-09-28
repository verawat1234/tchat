package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/shared/responses"
	"tchat/social/models"
	"tchat/social/services"
)

// UserHandler handles user-related social endpoints
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetSocialProfile godoc
// @Summary Get user's social profile
// @Description Retrieve a user's social profile information
// @Tags social, users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.DataResponse{data=models.SocialProfile}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/profiles/{userId} [get]
func (h *UserHandler) GetSocialProfile(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid user ID format")
		return
	}

	profile, err := h.userService.GetSocialProfile(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			responses.NotFoundResponse(c, "User profile not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to retrieve user profile")
		return
	}

	responses.SendDataResponse(c, profile)
}

// UpdateSocialProfile godoc
// @Summary Update user's social profile
// @Description Update a user's social profile information
// @Tags social, users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param request body models.UpdateSocialProfileRequest true "Profile update request"
// @Success 200 {object} responses.DataResponse{data=models.SocialProfile}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/profiles/{userId} [put]
func (h *UserHandler) UpdateSocialProfile(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid user ID format")
		return
	}

	var req models.UpdateSocialProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	profile, err := h.userService.UpdateSocialProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			responses.NotFoundResponse(c, "User profile not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to update user profile")
		return
	}

	responses.SendDataResponse(c, profile)
}

// DiscoverUsers godoc
// @Summary Discover users with similar interests
// @Description Discover users based on region and interests
// @Tags social, users, discovery
// @Accept json
// @Produce json
// @Param region query string false "Region filter"
// @Param interests query []string false "Interests filter"
// @Param limit query int false "Limit results (default: 20)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} responses.DataResponse{data=[]models.SocialProfile}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/discover/users [get]
func (h *UserHandler) DiscoverUsers(c *gin.Context) {
	var req models.UserDiscoveryRequest

	// Parse query parameters
	req.Region = c.Query("region")
	req.Interests = c.QueryArray("interests")

	if limitParam := c.Query("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			req.Limit = limit
		} else {
			responses.BadRequestResponse(c, "Invalid limit parameter")
			return
		}
	} else {
		req.Limit = 20 // default
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if offset, err := strconv.Atoi(offsetParam); err == nil {
			req.Offset = offset
		} else {
			responses.BadRequestResponse(c, "Invalid offset parameter")
			return
		}
	}

	users, err := h.userService.DiscoverUsers(c.Request.Context(), &req)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to discover users")
		return
	}

	responses.SendDataResponse(c, users)
}

// FollowUser godoc
// @Summary Follow a user
// @Description Create a following relationship between users
// @Tags social, users, relationships
// @Accept json
// @Produce json
// @Param request body models.FollowRequest true "Follow request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/follow [post]
func (h *UserHandler) FollowUser(c *gin.Context) {
	var req models.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Prevent self-following
	if req.FollowerID == req.FollowingID {
		responses.BadRequestResponse(c, "Cannot follow yourself")
		return
	}

	err := h.userService.FollowUser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "already following" {
			responses.ConflictResponse(c, "Already following this user")
			return
		}
		if err.Error() == "user not found" {
			responses.NotFoundResponse(c, "User not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to follow user")
		return
	}

	responses.SuccessMessageResponse(c, "Successfully followed user")
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Remove a following relationship between users
// @Tags social, users, relationships
// @Accept json
// @Produce json
// @Param followerId path string true "Follower User ID"
// @Param followingId path string true "Following User ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/follow/{followerId}/{followingId} [delete]
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerIDParam := c.Param("followerId")
	followingIDParam := c.Param("followingId")

	followerID, err := uuid.Parse(followerIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid follower ID format")
		return
	}

	followingID, err := uuid.Parse(followingIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid following ID format")
		return
	}

	err = h.userService.UnfollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		if err.Error() == "relationship not found" {
			responses.NotFoundResponse(c, "Following relationship not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to unfollow user")
		return
	}

	responses.SuccessMessageResponse(c, "Successfully unfollowed user")
}

// GetFollowers godoc
// @Summary Get user's followers
// @Description Retrieve a list of users following the specified user
// @Tags social, users, relationships
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Limit results (default: 20)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/followers/{userId} [get]
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid user ID format")
		return
	}

	limit := 20
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil {
			offset = o
		}
	}

	followers, err := h.userService.GetFollowers(c.Request.Context(), userID, limit, offset)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve followers")
		return
	}

	responses.SendDataResponse(c, followers)
}

// GetFollowing godoc
// @Summary Get users that a user is following
// @Description Retrieve a list of users that the specified user is following
// @Tags social, users, relationships
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Limit results (default: 20)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/following/{userId} [get]
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid user ID format")
		return
	}

	limit := 20
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil {
			offset = o
		}
	}

	following, err := h.userService.GetFollowing(c.Request.Context(), userID, limit, offset)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve following users")
		return
	}

	responses.SendDataResponse(c, following)
}

// GetUserAnalytics godoc
// @Summary Get user's social analytics
// @Description Retrieve analytics data for a user's social activity
// @Tags social, users, analytics
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param period query string false "Analytics period (7d, 30d, 90d)"
// @Success 200 {object} responses.DataResponse{data=models.UserAnalyticsResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/social/analytics/users/{userId} [get]
func (h *UserHandler) GetUserAnalytics(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.BadRequestResponse(c, "Invalid user ID format")
		return
	}

	period := c.DefaultQuery("period", "30d")

	analytics, err := h.userService.GetUserAnalytics(c.Request.Context(), userID, period)
	if err != nil {
		if err.Error() == "user not found" {
			responses.NotFoundResponse(c, "User not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to retrieve user analytics")
		return
	}

	responses.SendDataResponse(c, analytics)
}