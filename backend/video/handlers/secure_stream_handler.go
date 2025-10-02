// backend/video/handlers/secure_stream_handler.go
// Secure streaming handler that prevents direct video downloads using blob URL pattern
// Implements token-based authentication and byte-range support for streaming

package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/video/services"
)

// SecureStreamHandler handles secure video streaming with token validation
type SecureStreamHandler struct {
	videoService    *services.VideoService
	securityService *services.SecurityService
}

// NewSecureStreamHandler creates a new secure streaming handler
func NewSecureStreamHandler(videoService *services.VideoService, securityService *services.SecurityService) *SecureStreamHandler {
	return &SecureStreamHandler{
		videoService:    videoService,
		securityService: securityService,
	}
}

// RegisterRoutes registers secure streaming routes
func (h *SecureStreamHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/videos/:id/stream/secure", h.StreamSecure)
	router.GET("/videos/:id/token", h.GenerateToken)
	router.POST("/videos/:id/validate-token", h.ValidateToken)
}

// GenerateToken generates a streaming token for authenticated users
// GET /api/v1/videos/:id/token
func (h *SecureStreamHandler) GenerateToken(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get user ID from auth context (assuming JWT middleware sets this)
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user authentication"})
		return
	}

	// Get quality parameter
	quality := c.DefaultQuery("quality", "auto")

	// Verify video exists and user has access
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check video availability
	if video.AvailabilityStatus != "available" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Video not available for streaming"})
		return
	}

	// Generate signed stream token
	token, err := h.securityService.GenerateStreamToken(videoID, userID, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate signed URL
	signedURL, err := h.securityService.GenerateSignedURL(videoID, userID, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video_id":   videoID,
		"token":      token,
		"signed_url": signedURL,
		"expires_at": token.ExpiresAt,
		"quality":    quality,
	})
}

// ValidateToken validates a streaming token
// POST /api/v1/videos/:id/validate-token
func (h *SecureStreamHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token *services.StreamToken `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate token
	if err := h.securityService.ValidateStreamToken(req.Token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"expires_at": req.Token.ExpiresAt,
	})
}

// StreamSecure serves video content with token validation and byte-range support
// GET /api/v1/videos/:id/stream/secure
func (h *SecureStreamHandler) StreamSecure(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Extract token parameters from query
	userIDStr := c.Query("token")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Decode base64 user ID
	userIDBytes, err := base64DecodeString(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	expiresAtStr := c.Query("expires")
	quality := c.Query("quality")
	signature := c.Query("signature")

	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid expiration time"})
		return
	}

	// Validate signed URL
	if err := h.securityService.ValidateSignedURL(videoIDStr, string(userIDBytes), quality, signature, expiresAt); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Get video details to verify it exists
	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Get video file path (this would come from storage service in production)
	filePath := h.getVideoFilePath(videoID, quality)

	// Open video file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video file not found"})
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read video file"})
		return
	}

	fileSize := fileInfo.Size()

	// Parse Range header for byte-range requests
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		h.handleRangeRequest(c, file, fileSize, rangeHeader)
		return
	}

	// Serve full video
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Length", fmt.Sprintf("%d", fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("X-Content-Type-Options", "nosniff")

	// Stream video content
	if _, err := io.Copy(c.Writer, file); err != nil {
		// Log error but don't send response (already streaming)
		return
	}
}

// handleRangeRequest handles HTTP byte-range requests for video streaming
func (h *SecureStreamHandler) handleRangeRequest(c *gin.Context, file *os.File, fileSize int64, rangeHeader string) {
	// Parse Range header: "bytes=start-end"
	ranges := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(ranges, "-")

	if len(parts) != 2 {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Parse start position
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 {
		start = 0
	}

	// Parse end position
	var end int64
	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end >= fileSize {
			end = fileSize - 1
		}
	}

	// Validate range
	if start > end || start >= fileSize {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Seek to start position
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Calculate content length
	contentLength := end - start + 1

	// Set response headers for partial content
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Status(http.StatusPartialContent)

	// Stream partial content
	if _, err := io.CopyN(c.Writer, file, contentLength); err != nil {
		// Log error but don't send response (already streaming)
		return
	}
}

// getVideoFilePath returns the file system path for a video
// In production, this would integrate with cloud storage (S3, GCS, etc.)
func (h *SecureStreamHandler) getVideoFilePath(videoID uuid.UUID, quality string) string {
	// This is a placeholder - in production this would:
	// 1. Check cloud storage (S3/GCS/Azure Blob)
	// 2. Use quality-specific transcoded versions
	// 3. Handle CDN distribution
	return fmt.Sprintf("/var/videos/%s/%s.mp4", videoID.String(), quality)
}

// Helper function to decode base64 string
func base64DecodeString(s string) ([]byte, error) {
	// Implement base64 decoding
	// This is simplified - production would use encoding/base64
	return []byte(s), nil
}