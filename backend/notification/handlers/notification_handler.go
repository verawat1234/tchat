package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"tchat/notification/services"
	"tchat/shared/utils"
	"tchat/shared/config"
	"tchat/shared/middleware"
	"tchat/shared/events"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	service  services.NotificationService
	logger   *zap.Logger
	eventBus *events.EventBus
	config   *config.NotificationConfig
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	service services.NotificationService,
	logger *zap.Logger,
	eventBus *events.EventBus,
	config *config.NotificationConfig,
) *NotificationHandler {
	return &NotificationHandler{
		service:  service,
		logger:   logger,
		eventBus: eventBus,
		config:   config,
	}
}

// RegisterRoutes registers all notification routes with middleware
func (h *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Request logging middleware
	r.Use(middleware.RequestLogger(h.logger))

	// Recovery middleware
	r.Use(gin.Recovery())

	// Public routes
	public := r.Group("/api/v1/notifications")
	{
		public.POST("/webhooks/email", h.HandleEmailWebhook)
		public.POST("/webhooks/sms", h.HandleSMSWebhook)
		public.POST("/webhooks/push", h.HandlePushWebhook)
	}

	// Protected routes requiring authentication
	protected := r.Group("/api/v1/notifications")
	protected.Use(middleware.AuthMiddleware())
	{
		// Notification operations
		protected.POST("/send", middleware.RateLimit(100, time.Minute), h.SendNotification)
		protected.POST("/send/bulk", middleware.RateLimit(10, time.Minute), h.SendBulkNotifications)
		protected.GET("/", h.GetNotifications)
		protected.GET("/:id", h.GetNotification)
		protected.PUT("/:id/read", h.MarkAsRead)
		protected.DELETE("/:id", h.DeleteNotification)

		// Subscription management
		protected.GET("/subscriptions", h.GetSubscriptions)
		protected.POST("/subscriptions", h.CreateSubscription)
		protected.PUT("/subscriptions/:id", h.UpdateSubscription)
		protected.DELETE("/subscriptions/:id", h.DeleteSubscription)

		// Template management
		protected.GET("/templates", h.GetTemplates)
		protected.POST("/templates", middleware.AdminOnly(), h.CreateTemplate)
		protected.GET("/templates/:id", h.GetTemplate)
		protected.PUT("/templates/:id", middleware.AdminOnly(), h.UpdateTemplate)
		protected.DELETE("/templates/:id", middleware.AdminOnly(), h.DeleteTemplate)

		// Preferences
		protected.GET("/preferences", h.GetPreferences)
		protected.PUT("/preferences", h.UpdatePreferences)
	}

	// Admin routes
	admin := r.Group("/api/v1/admin/notifications")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/analytics", h.GetAnalytics)
		admin.GET("/delivery-reports", h.GetDeliveryReports)
		admin.POST("/broadcast", middleware.RateLimit(1, time.Minute), h.BroadcastNotification)
		admin.GET("/queues/status", h.GetQueueStatus)
		admin.POST("/queues/retry/:id", h.RetryFailedNotification)
	}
}

// Request/Response DTOs

type SendNotificationRequest struct {
	RecipientID  string                 `json:"recipient_id" binding:"required,uuid"`
	Type         string                 `json:"type" binding:"required,oneof=email sms push in_app"`
	Channel      string                 `json:"channel" binding:"required"`
	TemplateID   *string                `json:"template_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content" binding:"required,min=1,max=5000"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Priority     string                 `json:"priority" binding:"oneof=low medium high urgent"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type BulkNotificationRequest struct {
	RecipientIDs []string               `json:"recipient_ids" binding:"required,min=1,max=1000"`
	Type         string                 `json:"type" binding:"required,oneof=email sms push in_app"`
	Channel      string                 `json:"channel" binding:"required"`
	TemplateID   *string                `json:"template_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content" binding:"required,min=1,max=5000"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Priority     string                 `json:"priority" binding:"oneof=low medium high urgent"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
}

type BroadcastRequest struct {
	UserSegment  string                 `json:"user_segment" binding:"required"`
	Type         string                 `json:"type" binding:"required,oneof=email sms push in_app"`
	Channel      string                 `json:"channel" binding:"required"`
	TemplateID   *string                `json:"template_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content" binding:"required,min=1,max=5000"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Priority     string                 `json:"priority" binding:"oneof=low medium high urgent"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
}

type CreateSubscriptionRequest struct {
	Channel      string                 `json:"channel" binding:"required"`
	Type         string                 `json:"type" binding:"required,oneof=email sms push"`
	Endpoint     string                 `json:"endpoint" binding:"required"`
	Enabled      bool                   `json:"enabled"`
	Preferences  map[string]interface{} `json:"preferences,omitempty"`
}

type UpdateSubscriptionRequest struct {
	Enabled      *bool                  `json:"enabled,omitempty"`
	Endpoint     *string                `json:"endpoint,omitempty"`
	Preferences  map[string]interface{} `json:"preferences,omitempty"`
}

type CreateTemplateRequest struct {
	Name         string                 `json:"name" binding:"required,min=1,max=100"`
	Type         string                 `json:"type" binding:"required,oneof=email sms push in_app"`
	Category     string                 `json:"category" binding:"required"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content" binding:"required,min=1,max=10000"`
	Variables    []string               `json:"variables,omitempty"`
	Locales      map[string]interface{} `json:"locales,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateTemplateRequest struct {
	Name         *string                `json:"name,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	Content      *string                `json:"content,omitempty"`
	Variables    []string               `json:"variables,omitempty"`
	Locales      map[string]interface{} `json:"locales,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Active       *bool                  `json:"active,omitempty"`
}

type UpdatePreferencesRequest struct {
	EmailEnabled    *bool                  `json:"email_enabled,omitempty"`
	SMSEnabled      *bool                  `json:"sms_enabled,omitempty"`
	PushEnabled     *bool                  `json:"push_enabled,omitempty"`
	InAppEnabled    *bool                  `json:"in_app_enabled,omitempty"`
	Categories      map[string]bool        `json:"categories,omitempty"`
	QuietHours      map[string]interface{} `json:"quiet_hours,omitempty"`
	Languages       []string               `json:"languages,omitempty"`
}

type WebhookPayload struct {
	Event       string                 `json:"event"`
	Timestamp   time.Time              `json:"timestamp"`
	MessageID   string                 `json:"message_id"`
	Status      string                 `json:"status"`
	Recipient   string                 `json:"recipient"`
	ErrorCode   *string                `json:"error_code,omitempty"`
	ErrorMsg    *string                `json:"error_message,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Response DTOs
type NotificationResponse struct {
	ID           string                 `json:"id"`
	RecipientID  string                 `json:"recipient_id"`
	Type         string                 `json:"type"`
	Channel      string                 `json:"channel"`
	Subject      *string                `json:"subject,omitempty"`
	Content      string                 `json:"content"`
	Status       string                 `json:"status"`
	Priority     string                 `json:"priority"`
	Read         bool                   `json:"read"`
	ReadAt       *time.Time             `json:"read_at,omitempty"`
	SentAt       *time.Time             `json:"sent_at,omitempty"`
	DeliveredAt  *time.Time             `json:"delivered_at,omitempty"`
	FailedAt     *time.Time             `json:"failed_at,omitempty"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	RetryCount   int                    `json:"retry_count"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type SubscriptionResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Channel     string                 `json:"channel"`
	Type        string                 `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Enabled     bool                   `json:"enabled"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TemplateResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Subject     *string                `json:"subject,omitempty"`
	Content     string                 `json:"content"`
	Variables   []string               `json:"variables,omitempty"`
	Locales     map[string]interface{} `json:"locales,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Active      bool                   `json:"active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type PreferencesResponse struct {
	UserID          string                 `json:"user_id"`
	EmailEnabled    bool                   `json:"email_enabled"`
	SMSEnabled      bool                   `json:"sms_enabled"`
	PushEnabled     bool                   `json:"push_enabled"`
	InAppEnabled    bool                   `json:"in_app_enabled"`
	Categories      map[string]bool        `json:"categories"`
	QuietHours      map[string]interface{} `json:"quiet_hours"`
	Languages       []string               `json:"languages"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type AnalyticsResponse struct {
	Period          string                 `json:"period"`
	TotalSent       int64                  `json:"total_sent"`
	TotalDelivered  int64                  `json:"total_delivered"`
	TotalFailed     int64                  `json:"total_failed"`
	TotalOpened     int64                  `json:"total_opened"`
	TotalClicked    int64                  `json:"total_clicked"`
	DeliveryRate    float64                `json:"delivery_rate"`
	OpenRate        float64                `json:"open_rate"`
	ClickRate       float64                `json:"click_rate"`
	ByType          map[string]int64       `json:"by_type"`
	ByChannel       map[string]int64       `json:"by_channel"`
	TopCategories   []map[string]interface{} `json:"top_categories"`
	RecentActivity  []map[string]interface{} `json:"recent_activity"`
}

type QueueStatusResponse struct {
	Pending    int64 `json:"pending"`
	Processing int64 `json:"processing"`
	Failed     int64 `json:"failed"`
	Retry      int64 `json:"retry"`
	Scheduled  int64 `json:"scheduled"`
}

// Notification Handlers

// SendNotification sends a single notification
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate recipient
	if !utils.IsValidUUID(req.RecipientID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
		return
	}

	// Validate content
	if err := utils.ValidateHTML(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Validate schedule time
	if req.ScheduledAt != nil && req.ScheduledAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Scheduled time must be in the future"})
		return
	}

	// Validate expiration time
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expiration time must be in the future"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create notification
	notification, err := h.service.SendNotification(c.Request.Context(), services.CreateNotificationParams{
		SenderID:     userID.(string),
		RecipientID:  req.RecipientID,
		Type:         req.Type,
		Channel:      req.Channel,
		TemplateID:   req.TemplateID,
		Subject:      req.Subject,
		Content:      req.Content,
		Variables:    req.Variables,
		Priority:     req.Priority,
		ScheduledAt:  req.ScheduledAt,
		ExpiresAt:    req.ExpiresAt,
		Metadata:     req.Metadata,
	})

	if err != nil {
		h.logger.Error("Failed to send notification", zap.Error(err))

		if strings.Contains(err.Error(), "recipient not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
			return
		}
		if strings.Contains(err.Error(), "template not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		if strings.Contains(err.Error(), "rate limit") {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	// Emit event
	h.eventBus.Publish("notification.sent", map[string]interface{}{
		"notification_id": notification.ID,
		"recipient_id":    notification.RecipientID,
		"type":           notification.Type,
		"channel":        notification.Channel,
		"timestamp":      time.Now(),
	})

	response := &NotificationResponse{
		ID:          notification.ID,
		RecipientID: notification.RecipientID,
		Type:        notification.Type,
		Channel:     notification.Channel,
		Subject:     notification.Subject,
		Content:     notification.Content,
		Status:      notification.Status,
		Priority:    notification.Priority,
		Read:        notification.Read,
		ReadAt:      notification.ReadAt,
		SentAt:      notification.SentAt,
		DeliveredAt: notification.DeliveredAt,
		FailedAt:    notification.FailedAt,
		ScheduledAt: notification.ScheduledAt,
		ExpiresAt:   notification.ExpiresAt,
		RetryCount:  notification.RetryCount,
		Metadata:    notification.Metadata,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"notification": response})
}

// SendBulkNotifications sends notifications to multiple recipients
func (h *NotificationHandler) SendBulkNotifications(c *gin.Context) {
	var req BulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate recipients
	for _, recipientID := range req.RecipientIDs {
		if !utils.IsValidUUID(recipientID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid recipient ID: %s", recipientID)})
			return
		}
	}

	// Validate content
	if err := utils.ValidateHTML(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Send bulk notifications
	notifications, err := h.service.SendBulkNotifications(c.Request.Context(), services.BulkNotificationParams{
		SenderID:     userID.(string),
		RecipientIDs: req.RecipientIDs,
		Type:         req.Type,
		Channel:      req.Channel,
		TemplateID:   req.TemplateID,
		Subject:      req.Subject,
		Content:      req.Content,
		Variables:    req.Variables,
		Priority:     req.Priority,
		ScheduledAt:  req.ScheduledAt,
		ExpiresAt:    req.ExpiresAt,
	})

	if err != nil {
		h.logger.Error("Failed to send bulk notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send bulk notifications"})
		return
	}

	// Convert to response format
	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = &NotificationResponse{
			ID:          notification.ID,
			RecipientID: notification.RecipientID,
			Type:        notification.Type,
			Channel:     notification.Channel,
			Subject:     notification.Subject,
			Content:     notification.Content,
			Status:      notification.Status,
			Priority:    notification.Priority,
			Read:        notification.Read,
			ReadAt:      notification.ReadAt,
			SentAt:      notification.SentAt,
			DeliveredAt: notification.DeliveredAt,
			FailedAt:    notification.FailedAt,
			ScheduledAt: notification.ScheduledAt,
			ExpiresAt:   notification.ExpiresAt,
			RetryCount:  notification.RetryCount,
			Metadata:    notification.Metadata,
			CreatedAt:   notification.CreatedAt,
			UpdatedAt:   notification.UpdatedAt,
		}
	}

	// Emit bulk event
	h.eventBus.Publish("notifications.bulk_sent", map[string]interface{}{
		"sender_id":      userID.(string),
		"recipient_count": len(req.RecipientIDs),
		"type":           req.Type,
		"channel":        req.Channel,
		"timestamp":      time.Now(),
	})

	c.JSON(http.StatusCreated, gin.H{
		"notifications": responses,
		"count":        len(responses),
	})
}

// BroadcastNotification sends notification to user segment
func (h *NotificationHandler) BroadcastNotification(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate content
	if err := utils.ValidateHTML(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Get authenticated admin user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Send broadcast
	broadcastID, count, err := h.service.BroadcastNotification(c.Request.Context(), services.BroadcastParams{
		SenderID:    userID.(string),
		UserSegment: req.UserSegment,
		Type:        req.Type,
		Channel:     req.Channel,
		TemplateID:  req.TemplateID,
		Subject:     req.Subject,
		Content:     req.Content,
		Variables:   req.Variables,
		Priority:    req.Priority,
		ScheduledAt: req.ScheduledAt,
	})

	if err != nil {
		h.logger.Error("Failed to broadcast notification", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to broadcast notification"})
		return
	}

	// Emit broadcast event
	h.eventBus.Publish("notification.broadcast", map[string]interface{}{
		"broadcast_id":    broadcastID,
		"sender_id":       userID.(string),
		"user_segment":    req.UserSegment,
		"recipient_count": count,
		"type":           req.Type,
		"channel":        req.Channel,
		"timestamp":      time.Now(),
	})

	c.JSON(http.StatusAccepted, gin.H{
		"broadcast_id":    broadcastID,
		"recipient_count": count,
		"message":        "Broadcast notification queued",
	})
}

// GetNotifications retrieves user's notifications with pagination
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	notificationType := c.Query("type")
	status := c.Query("status")
	unreadOnly := c.Query("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	notifications, total, err := h.service.GetNotifications(c.Request.Context(), services.GetNotificationsParams{
		UserID:     userID.(string),
		Type:       notificationType,
		Status:     status,
		UnreadOnly: unreadOnly,
		Page:       page,
		Limit:      limit,
	})

	if err != nil {
		h.logger.Error("Failed to get notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	// Convert to response format
	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = &NotificationResponse{
			ID:          notification.ID,
			RecipientID: notification.RecipientID,
			Type:        notification.Type,
			Channel:     notification.Channel,
			Subject:     notification.Subject,
			Content:     notification.Content,
			Status:      notification.Status,
			Priority:    notification.Priority,
			Read:        notification.Read,
			ReadAt:      notification.ReadAt,
			SentAt:      notification.SentAt,
			DeliveredAt: notification.DeliveredAt,
			FailedAt:    notification.FailedAt,
			ScheduledAt: notification.ScheduledAt,
			ExpiresAt:   notification.ExpiresAt,
			RetryCount:  notification.RetryCount,
			Metadata:    notification.Metadata,
			CreatedAt:   notification.CreatedAt,
			UpdatedAt:   notification.UpdatedAt,
		}
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"notifications": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < int(totalPages),
			"has_prev":    page > 1,
		},
	})
}

// GetNotification retrieves a specific notification
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	notificationID := c.Param("id")
	if !utils.IsValidUUID(notificationID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	notification, err := h.service.GetNotification(c.Request.Context(), notificationID, userID.(string))
	if err != nil {
		h.logger.Error("Failed to get notification", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notification"})
		return
	}

	response := &NotificationResponse{
		ID:          notification.ID,
		RecipientID: notification.RecipientID,
		Type:        notification.Type,
		Channel:     notification.Channel,
		Subject:     notification.Subject,
		Content:     notification.Content,
		Status:      notification.Status,
		Priority:    notification.Priority,
		Read:        notification.Read,
		ReadAt:      notification.ReadAt,
		SentAt:      notification.SentAt,
		DeliveredAt: notification.DeliveredAt,
		FailedAt:    notification.FailedAt,
		ScheduledAt: notification.ScheduledAt,
		ExpiresAt:   notification.ExpiresAt,
		RetryCount:  notification.RetryCount,
		Metadata:    notification.Metadata,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"notification": response})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	if !utils.IsValidUUID(notificationID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.service.MarkAsRead(c.Request.Context(), notificationID, userID.(string))
	if err != nil {
		h.logger.Error("Failed to mark notification as read", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	// Emit read event
	h.eventBus.Publish("notification.read", map[string]interface{}{
		"notification_id": notificationID,
		"user_id":        userID.(string),
		"timestamp":      time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationID := c.Param("id")
	if !utils.IsValidUUID(notificationID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.service.DeleteNotification(c.Request.Context(), notificationID, userID.(string))
	if err != nil {
		h.logger.Error("Failed to delete notification", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

// Subscription Handlers

// GetSubscriptions retrieves user's notification subscriptions
func (h *NotificationHandler) GetSubscriptions(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscriptions, err := h.service.GetSubscriptions(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	// Convert to response format
	responses := make([]*SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = &SubscriptionResponse{
			ID:          subscription.ID,
			UserID:      subscription.UserID,
			Channel:     subscription.Channel,
			Type:        subscription.Type,
			Endpoint:    subscription.Endpoint,
			Enabled:     subscription.Enabled,
			Preferences: subscription.Preferences,
			CreatedAt:   subscription.CreatedAt,
			UpdatedAt:   subscription.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": responses})
}

// CreateSubscription creates a new notification subscription
func (h *NotificationHandler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate endpoint based on type
	if req.Type == "email" && !utils.IsValidEmail(req.Endpoint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}
	if req.Type == "sms" && !utils.IsValidPhoneNumber(req.Endpoint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), services.CreateSubscriptionParams{
		UserID:      userID.(string),
		Channel:     req.Channel,
		Type:        req.Type,
		Endpoint:    req.Endpoint,
		Enabled:     req.Enabled,
		Preferences: req.Preferences,
	})

	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))

		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Subscription already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	response := &SubscriptionResponse{
		ID:          subscription.ID,
		UserID:      subscription.UserID,
		Channel:     subscription.Channel,
		Type:        subscription.Type,
		Endpoint:    subscription.Endpoint,
		Enabled:     subscription.Enabled,
		Preferences: subscription.Preferences,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"subscription": response})
}

// UpdateSubscription updates an existing subscription
func (h *NotificationHandler) UpdateSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if !utils.IsValidUUID(subscriptionID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate endpoint if provided
	if req.Endpoint != nil {
		// Note: We'd need to get the subscription type to validate properly
		// For now, just check if it's not empty
		if strings.TrimSpace(*req.Endpoint) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint cannot be empty"})
			return
		}
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscription, err := h.service.UpdateSubscription(c.Request.Context(), services.UpdateSubscriptionParams{
		ID:          subscriptionID,
		UserID:      userID.(string),
		Enabled:     req.Enabled,
		Endpoint:    req.Endpoint,
		Preferences: req.Preferences,
	})

	if err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	response := &SubscriptionResponse{
		ID:          subscription.ID,
		UserID:      subscription.UserID,
		Channel:     subscription.Channel,
		Type:        subscription.Type,
		Endpoint:    subscription.Endpoint,
		Enabled:     subscription.Enabled,
		Preferences: subscription.Preferences,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"subscription": response})
}

// DeleteSubscription deletes a subscription
func (h *NotificationHandler) DeleteSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if !utils.IsValidUUID(subscriptionID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.service.DeleteSubscription(c.Request.Context(), subscriptionID, userID.(string))
	if err != nil {
		h.logger.Error("Failed to delete subscription", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// Template Handlers

// GetTemplates retrieves notification templates with pagination
func (h *NotificationHandler) GetTemplates(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	templateType := c.Query("type")
	category := c.Query("category")
	activeOnly := c.Query("active_only") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	templates, total, err := h.service.GetTemplates(c.Request.Context(), services.GetTemplatesParams{
		Type:       templateType,
		Category:   category,
		ActiveOnly: activeOnly,
		Page:       page,
		Limit:      limit,
	})

	if err != nil {
		h.logger.Error("Failed to get templates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
		return
	}

	// Convert to response format
	responses := make([]*TemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = &TemplateResponse{
			ID:        template.ID,
			Name:      template.Name,
			Type:      template.Type,
			Category:  template.Category,
			Subject:   template.Subject,
			Content:   template.Content,
			Variables: template.Variables,
			Locales:   template.Locales,
			Metadata:  template.Metadata,
			Active:    template.Active,
			CreatedAt: template.CreatedAt,
			UpdatedAt: template.UpdatedAt,
		}
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"templates": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < int(totalPages),
			"has_prev":    page > 1,
		},
	})
}

// CreateTemplate creates a new notification template (admin only)
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate content
	if err := utils.ValidateHTML(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Get authenticated admin user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	template, err := h.service.CreateTemplate(c.Request.Context(), services.CreateTemplateParams{
		CreatedBy: userID.(string),
		Name:      req.Name,
		Type:      req.Type,
		Category:  req.Category,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: req.Variables,
		Locales:   req.Locales,
		Metadata:  req.Metadata,
	})

	if err != nil {
		h.logger.Error("Failed to create template", zap.Error(err))

		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Template name already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	response := &TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      template.Type,
		Category:  template.Category,
		Subject:   template.Subject,
		Content:   template.Content,
		Variables: template.Variables,
		Locales:   template.Locales,
		Metadata:  template.Metadata,
		Active:    template.Active,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"template": response})
}

// GetTemplate retrieves a specific template
func (h *NotificationHandler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if !utils.IsValidUUID(templateID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.service.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		h.logger.Error("Failed to get template", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve template"})
		return
	}

	response := &TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      template.Type,
		Category:  template.Category,
		Subject:   template.Subject,
		Content:   template.Content,
		Variables: template.Variables,
		Locales:   template.Locales,
		Metadata:  template.Metadata,
		Active:    template.Active,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"template": response})
}

// UpdateTemplate updates an existing template (admin only)
func (h *NotificationHandler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if !utils.IsValidUUID(templateID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Validate content if provided
	if req.Content != nil {
		if err := utils.ValidateHTML(*req.Content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
			return
		}
	}

	// Get authenticated admin user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	template, err := h.service.UpdateTemplate(c.Request.Context(), services.UpdateTemplateParams{
		ID:        templateID,
		UpdatedBy: userID.(string),
		Name:      req.Name,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: req.Variables,
		Locales:   req.Locales,
		Metadata:  req.Metadata,
		Active:    req.Active,
	})

	if err != nil {
		h.logger.Error("Failed to update template", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Template name already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	response := &TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      template.Type,
		Category:  template.Category,
		Subject:   template.Subject,
		Content:   template.Content,
		Variables: template.Variables,
		Locales:   template.Locales,
		Metadata:  template.Metadata,
		Active:    template.Active,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"template": response})
}

// DeleteTemplate deletes a template (admin only)
func (h *NotificationHandler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if !utils.IsValidUUID(templateID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	err := h.service.DeleteTemplate(c.Request.Context(), templateID)
	if err != nil {
		h.logger.Error("Failed to delete template", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		if strings.Contains(err.Error(), "in use") {
			c.JSON(http.StatusConflict, gin.H{"error": "Template is currently in use"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// Preferences Handlers

// GetPreferences retrieves user's notification preferences
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	preferences, err := h.service.GetPreferences(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve preferences"})
		return
	}

	response := &PreferencesResponse{
		UserID:       preferences.UserID,
		EmailEnabled: preferences.EmailEnabled,
		SMSEnabled:   preferences.SMSEnabled,
		PushEnabled:  preferences.PushEnabled,
		InAppEnabled: preferences.InAppEnabled,
		Categories:   preferences.Categories,
		QuietHours:   preferences.QuietHours,
		Languages:    preferences.Languages,
		UpdatedAt:    preferences.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"preferences": response})
}

// UpdatePreferences updates user's notification preferences
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	preferences, err := h.service.UpdatePreferences(c.Request.Context(), services.UpdatePreferencesParams{
		UserID:       userID.(string),
		EmailEnabled: req.EmailEnabled,
		SMSEnabled:   req.SMSEnabled,
		PushEnabled:  req.PushEnabled,
		InAppEnabled: req.InAppEnabled,
		Categories:   req.Categories,
		QuietHours:   req.QuietHours,
		Languages:    req.Languages,
	})

	if err != nil {
		h.logger.Error("Failed to update preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	response := &PreferencesResponse{
		UserID:       preferences.UserID,
		EmailEnabled: preferences.EmailEnabled,
		SMSEnabled:   preferences.SMSEnabled,
		PushEnabled:  preferences.PushEnabled,
		InAppEnabled: preferences.InAppEnabled,
		Categories:   preferences.Categories,
		QuietHours:   preferences.QuietHours,
		Languages:    preferences.Languages,
		UpdatedAt:    preferences.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"preferences": response})
}

// Admin Analytics Handlers

// GetAnalytics retrieves notification analytics (admin only)
func (h *NotificationHandler) GetAnalytics(c *gin.Context) {
	// Parse query parameters
	period := c.DefaultQuery("period", "7d") // 1d, 7d, 30d, 90d
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	analytics, err := h.service.GetAnalytics(c.Request.Context(), services.AnalyticsParams{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		h.logger.Error("Failed to get analytics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics"})
		return
	}

	response := &AnalyticsResponse{
		Period:         analytics.Period,
		TotalSent:      analytics.TotalSent,
		TotalDelivered: analytics.TotalDelivered,
		TotalFailed:    analytics.TotalFailed,
		TotalOpened:    analytics.TotalOpened,
		TotalClicked:   analytics.TotalClicked,
		DeliveryRate:   analytics.DeliveryRate,
		OpenRate:       analytics.OpenRate,
		ClickRate:      analytics.ClickRate,
		ByType:         analytics.ByType,
		ByChannel:      analytics.ByChannel,
		TopCategories:  analytics.TopCategories,
		RecentActivity: analytics.RecentActivity,
	}

	c.JSON(http.StatusOK, gin.H{"analytics": response})
}

// GetDeliveryReports retrieves detailed delivery reports (admin only)
func (h *NotificationHandler) GetDeliveryReports(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	status := c.Query("status")
	notificationType := c.Query("type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	reports, total, err := h.service.GetDeliveryReports(c.Request.Context(), services.DeliveryReportsParams{
		Status:    status,
		Type:      notificationType,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	})

	if err != nil {
		h.logger.Error("Failed to get delivery reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve delivery reports"})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < int(totalPages),
			"has_prev":    page > 1,
		},
	})
}

// GetQueueStatus retrieves notification queue status (admin only)
func (h *NotificationHandler) GetQueueStatus(c *gin.Context) {
	status, err := h.service.GetQueueStatus(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get queue status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve queue status"})
		return
	}

	response := &QueueStatusResponse{
		Pending:    status.Pending,
		Processing: status.Processing,
		Failed:     status.Failed,
		Retry:      status.Retry,
		Scheduled:  status.Scheduled,
	}

	c.JSON(http.StatusOK, gin.H{"queue_status": response})
}

// RetryFailedNotification retries a failed notification (admin only)
func (h *NotificationHandler) RetryFailedNotification(c *gin.Context) {
	notificationID := c.Param("id")
	if !utils.IsValidUUID(notificationID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err := h.service.RetryFailedNotification(c.Request.Context(), notificationID)
	if err != nil {
		h.logger.Error("Failed to retry notification", zap.Error(err))

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		if strings.Contains(err.Error(), "not failed") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notification is not in failed state"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification queued for retry"})
}

// Webhook Handlers

// HandleEmailWebhook handles email delivery webhooks
func (h *NotificationHandler) HandleEmailWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Invalid webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify webhook signature (implementation depends on email provider)
	// This is a simplified version
	signature := c.GetHeader("X-Signature")
	if !h.verifyWebhookSignature(c.Request, signature, "email") {
		h.logger.Warn("Invalid webhook signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	err := h.service.ProcessWebhook(c.Request.Context(), services.WebhookParams{
		Provider:  "email",
		Event:     payload.Event,
		MessageID: payload.MessageID,
		Status:    payload.Status,
		Recipient: payload.Recipient,
		ErrorCode: payload.ErrorCode,
		ErrorMsg:  payload.ErrorMsg,
		Metadata:  payload.Metadata,
		Timestamp: payload.Timestamp,
	})

	if err != nil {
		h.logger.Error("Failed to process email webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook processing failed"})
		return
	}

	// Emit webhook event
	h.eventBus.Publish("webhook.email.received", map[string]interface{}{
		"event":      payload.Event,
		"message_id": payload.MessageID,
		"status":     payload.Status,
		"recipient":  payload.Recipient,
		"timestamp":  payload.Timestamp,
	})

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// HandleSMSWebhook handles SMS delivery webhooks
func (h *NotificationHandler) HandleSMSWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Invalid webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("X-Signature")
	if !h.verifyWebhookSignature(c.Request, signature, "sms") {
		h.logger.Warn("Invalid webhook signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	err := h.service.ProcessWebhook(c.Request.Context(), services.WebhookParams{
		Provider:  "sms",
		Event:     payload.Event,
		MessageID: payload.MessageID,
		Status:    payload.Status,
		Recipient: payload.Recipient,
		ErrorCode: payload.ErrorCode,
		ErrorMsg:  payload.ErrorMsg,
		Metadata:  payload.Metadata,
		Timestamp: payload.Timestamp,
	})

	if err != nil {
		h.logger.Error("Failed to process SMS webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook processing failed"})
		return
	}

	// Emit webhook event
	h.eventBus.Publish("webhook.sms.received", map[string]interface{}{
		"event":      payload.Event,
		"message_id": payload.MessageID,
		"status":     payload.Status,
		"recipient":  payload.Recipient,
		"timestamp":  payload.Timestamp,
	})

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// HandlePushWebhook handles push notification webhooks
func (h *NotificationHandler) HandlePushWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Invalid webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("X-Signature")
	if !h.verifyWebhookSignature(c.Request, signature, "push") {
		h.logger.Warn("Invalid webhook signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	err := h.service.ProcessWebhook(c.Request.Context(), services.WebhookParams{
		Provider:  "push",
		Event:     payload.Event,
		MessageID: payload.MessageID,
		Status:    payload.Status,
		Recipient: payload.Recipient,
		ErrorCode: payload.ErrorCode,
		ErrorMsg:  payload.ErrorMsg,
		Metadata:  payload.Metadata,
		Timestamp: payload.Timestamp,
	})

	if err != nil {
		h.logger.Error("Failed to process push webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook processing failed"})
		return
	}

	// Emit webhook event
	h.eventBus.Publish("webhook.push.received", map[string]interface{}{
		"event":      payload.Event,
		"message_id": payload.MessageID,
		"status":     payload.Status,
		"recipient":  payload.Recipient,
		"timestamp":  payload.Timestamp,
	})

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// Helper Methods

// verifyWebhookSignature verifies webhook signature for security
func (h *NotificationHandler) verifyWebhookSignature(r *http.Request, signature string, provider string) bool {
	// This is a simplified implementation
	// In production, you would implement proper HMAC verification
	// based on the specific provider's webhook signature format

	if signature == "" {
		return false
	}

	// For demo purposes, just check if signature exists
	// In real implementation:
	// 1. Get the webhook secret for the provider
	// 2. Compute HMAC of request body
	// 3. Compare with provided signature

	return len(signature) > 0
}