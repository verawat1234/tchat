package handlers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/messaging/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

type MessageHandler struct {
	messageService *services.MessageService
	validator      *validator.Validate
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		validator:      validator.New(),
	}
}

// SendMessageRequest represents the request to send a new message
type SendMessageRequest struct {
	Type           string                 `json:"type" validate:"required,oneof=text image video audio file location contact sticker gif" example:"text"`
	Content        string                 `json:"content,omitempty" validate:"omitempty,max=4000" example:"Hello everyone! How's it going?"`
	ReplyToID      *uuid.UUID             `json:"reply_to_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	ForwardFromID  *uuid.UUID             `json:"forward_from_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ClientMessageID string                `json:"client_message_id,omitempty" example:"client_123456789"`
	ScheduledAt    string                 `json:"scheduled_at,omitempty" example:"2024-01-21T10:00:00Z"`
	Silent         bool                   `json:"silent,omitempty" example:"false"`
	SelfDestruct   int                    `json:"self_destruct,omitempty" example:"0"`
}

// UpdateMessageRequest represents the request to update/edit a message
type UpdateMessageRequest struct {
	Content  string                 `json:"content" validate:"required,max=4000" example:"Updated message content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	ID              uuid.UUID              `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DialogID        uuid.UUID              `json:"dialog_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	SenderID        uuid.UUID              `json:"sender_id" example:"123e4567-e89b-12d3-a456-426614174002"`
	Type            string                 `json:"type" example:"text"`
	Content         string                 `json:"content" example:"Hello everyone! How's it going?"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	ReplyTo         *MessageResponse       `json:"reply_to,omitempty"`
	ForwardFrom     *ForwardInfo           `json:"forward_from,omitempty"`
	Attachments     []AttachmentInfo       `json:"attachments,omitempty"`
	Reactions       []ReactionInfo         `json:"reactions,omitempty"`
	ReadBy          []ReadInfo             `json:"read_by,omitempty"`
	EditHistory     []EditInfo             `json:"edit_history,omitempty"`
	IsEdited        bool                   `json:"is_edited" example:"false"`
	IsDeleted       bool                   `json:"is_deleted" example:"false"`
	IsForwarded     bool                   `json:"is_forwarded" example:"false"`
	SentAt          string                 `json:"sent_at" example:"2024-01-20T15:45:00Z"`
	EditedAt        string                 `json:"edited_at,omitempty" example:"2024-01-20T16:00:00Z"`
	DeliveredAt     string                 `json:"delivered_at,omitempty" example:"2024-01-20T15:45:01Z"`
	SenderInfo      SenderInfo             `json:"sender_info"`
	CanEdit         bool                   `json:"can_edit" example:"true"`
	CanDelete       bool                   `json:"can_delete" example:"true"`
	CanReact        bool                   `json:"can_react" example:"true"`
	CanReply        bool                   `json:"can_reply" example:"true"`
	SelfDestruct    int                    `json:"self_destruct,omitempty" example:"0"`
	ClientMessageID string                 `json:"client_message_id,omitempty" example:"client_123456789"`
}

type SenderInfo struct {
	DisplayName string `json:"display_name" example:"John Doe"`
	Avatar      string `json:"avatar,omitempty" example:"https://cdn.tchat.com/avatars/user.jpg"`
	IsOnline    bool   `json:"is_online" example:"true"`
}

type ForwardInfo struct {
	OriginalSenderID   uuid.UUID `json:"original_sender_id" example:"123e4567-e89b-12d3-a456-426614174003"`
	OriginalSenderName string    `json:"original_sender_name" example:"Jane Doe"`
	OriginalDialogID   uuid.UUID `json:"original_dialog_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174004"`
	OriginalSentAt     string    `json:"original_sent_at" example:"2024-01-19T10:30:00Z"`
}

type AttachmentInfo struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174005"`
	Type        string    `json:"type" example:"image"`
	Name        string    `json:"name" example:"photo.jpg"`
	Size        int64     `json:"size" example:"1048576"`
	MimeType    string    `json:"mime_type" example:"image/jpeg"`
	URL         string    `json:"url" example:"https://cdn.tchat.com/files/photo.jpg"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty" example:"https://cdn.tchat.com/thumbs/photo.jpg"`
	Width       int       `json:"width,omitempty" example:"1920"`
	Height      int       `json:"height,omitempty" example:"1080"`
	Duration    int       `json:"duration,omitempty" example:"120"`
}

type ReactionInfo struct {
	Emoji   string      `json:"emoji" example:"👍"`
	Count   int         `json:"count" example:"3"`
	UserIDs []uuid.UUID `json:"user_ids,omitempty"`
	HasUser bool        `json:"has_user" example:"true"`
}

type ReadInfo struct {
	UserID   uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserName string    `json:"user_name" example:"John Doe"`
	ReadAt   string    `json:"read_at" example:"2024-01-20T15:46:00Z"`
}

type EditInfo struct {
	EditedAt    string `json:"edited_at" example:"2024-01-20T16:00:00Z"`
	PrevContent string `json:"prev_content" example:"Previous message content"`
}

// MarkAsReadRequest represents the request to mark messages as read
type MarkAsReadRequest struct {
	MessageIDs []uuid.UUID `json:"message_ids" validate:"required,min=1,max=100" example:"[\"123e4567-e89b-12d3-a456-426614174000\"]"`
}

// AddReactionRequest represents the request to add a reaction
type AddReactionRequest struct {
	Emoji string `json:"emoji" validate:"required" example:"👍"`
}

// @Summary Get messages for a dialog
// @Description Get paginated list of messages for a specific dialog
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param limit query int false "Number of messages to return" default(50) minimum(1) maximum(100)
// @Param before query string false "Message ID to get messages before (for pagination)"
// @Param after query string false "Message ID to get messages after (for pagination)"
// @Param search query string false "Search term for message content"
// @Param type query string false "Message type filter" Enums(text,image,video,audio,file,location,contact,sticker,gif)
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var beforeID, afterID *uuid.UUID
	if before := c.Query("before"); before != "" {
		if id, err := uuid.Parse(before); err == nil {
			beforeID = &id
		}
	}
	if after := c.Query("after"); after != "" {
		if id, err := uuid.Parse(after); err == nil {
			afterID = &id
		}
	}

	search := c.Query("search")
	messageType := c.Query("type")

	// Build service request
	req := &services.GetMessagesRequest{
		DialogID:    dialogID,
		UserID:      userUUID,
		Limit:       limit,
		BeforeID:    beforeID,
		AfterID:     afterID,
		Search:      search,
		MessageType: messageType,
	}

	// Get messages from service
	result, err := h.messageService.GetMessages(c.Request.Context(), req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "dialog not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this dialog.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get messages", "Unable to retrieve messages.")
			return
		}
	}

	// Convert to response format
	messageResponses := make([]MessageResponse, len(result.Messages))
	for i, message := range result.Messages {
		messageResponses[i] = h.convertToMessageResponse(message, userUUID)
	}

	// Build response data
	data := gin.H{
		"messages":   messageResponses,
		"total":      result.Total,
		"limit":      limit,
		"has_more":   result.HasMore,
		"dialog_id":  dialogID,
		"unread_count": result.UnreadCount,
	}

	// Add pagination cursors
	if len(messageResponses) > 0 {
		data["first_message_id"] = messageResponses[0].ID
		data["last_message_id"] = messageResponses[len(messageResponses)-1].ID
	}

	responses.DataResponse(c, data)
}

// @Summary Send a message
// @Description Send a new message to a dialog
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param request body SendMessageRequest true "Message content"
// @Success 201 {object} responses.DataResponse{data=MessageResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 413 {object} responses.ErrorResponse "Content too large"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest

	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.SendMessageRequest{
		DialogID:        dialogID,
		SenderID:        userUUID,
		Type:            req.Type,
		Content:         req.Content,
		ReplyToID:       req.ReplyToID,
		ForwardFromID:   req.ForwardFromID,
		Metadata:        req.Metadata,
		ClientMessageID: req.ClientMessageID,
		ScheduledAt:     req.ScheduledAt,
		Silent:          req.Silent,
		SelfDestruct:    req.SelfDestruct,
	}

	// Send message
	message, err := h.messageService.SendMessage(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "dialog not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to send messages to this dialog.")
			return
		case strings.Contains(err.Error(), "content too large"):
			responses.ErrorResponse(c, http.StatusRequestEntityTooLarge, "Content too large", "Message content exceeds maximum allowed size.")
			return
		case strings.Contains(err.Error(), "reply message not found"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Reply message not found", "The message you're replying to doesn't exist.")
			return
		case strings.Contains(err.Error(), "moderation failed"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Content blocked", "Message content violates community guidelines.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to send message", "Unable to send message.")
			return
		}
	}

	// Convert to response format
	messageResponse := h.convertToMessageResponse(message, userUUID)

	// Log message sent
	middleware.LogInfo(c, "Message sent", gin.H{
		"message_id":       message.ID,
		"dialog_id":        dialogID,
		"sender_id":        userUUID,
		"type":             req.Type,
		"client_message_id": req.ClientMessageID,
	})

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    messageResponse,
	})
}

// @Summary Upload and send file message
// @Description Upload a file and send it as a message
// @Tags messages
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param file formData file true "File to upload"
// @Param type formData string false "Message type" default(file)
// @Param content formData string false "Optional message content"
// @Param reply_to_id formData string false "Message ID to reply to"
// @Param client_message_id formData string false "Client-side message ID"
// @Success 201 {object} responses.DataResponse{data=MessageResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 413 {object} responses.ErrorResponse "File too large"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/messages/upload [post]
func (h *MessageHandler) UploadAndSendFile(c *gin.Context) {
	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(50 << 20) // 50MB max
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid form data", "Failed to parse multipart form.")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Missing file", "File is required.")
		return
	}
	defer file.Close()

	// Get form parameters
	messageType := c.DefaultPostForm("type", "file")
	content := c.PostForm("content")
	clientMessageID := c.PostForm("client_message_id")

	var replyToID *uuid.UUID
	if replyTo := c.PostForm("reply_to_id"); replyTo != "" {
		if id, err := uuid.Parse(replyTo); err == nil {
			replyToID = &id
		}
	}

	// Build service request
	serviceReq := &services.UploadAndSendFileRequest{
		DialogID:        dialogID,
		SenderID:        userUUID,
		Type:            messageType,
		Content:         content,
		ReplyToID:       replyToID,
		ClientMessageID: clientMessageID,
		File:            file,
		FileHeader:      header,
	}

	// Upload and send file message
	message, err := h.messageService.UploadAndSendFile(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "dialog not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to send files to this dialog.")
			return
		case strings.Contains(err.Error(), "file too large"):
			responses.ErrorResponse(c, http.StatusRequestEntityTooLarge, "File too large", "File size exceeds maximum allowed limit.")
			return
		case strings.Contains(err.Error(), "unsupported file type"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Unsupported file type", "This file type is not allowed.")
			return
		case strings.Contains(err.Error(), "upload failed"):
			responses.ErrorResponse(c, http.StatusInternalServerError, "Upload failed", "Failed to upload file.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to send file", "Unable to send file message.")
			return
		}
	}

	// Convert to response format
	messageResponse := h.convertToMessageResponse(message, userUUID)

	// Log file message sent
	middleware.LogInfo(c, "File message sent", gin.H{
		"message_id":       message.ID,
		"dialog_id":        dialogID,
		"sender_id":        userUUID,
		"type":             messageType,
		"file_name":        header.Filename,
		"file_size":        header.Size,
		"client_message_id": clientMessageID,
	})

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    messageResponse,
	})
}

// @Summary Update/edit a message
// @Description Update the content of an existing message
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Param request body UpdateMessageRequest true "Updated message content"
// @Success 200 {object} responses.DataResponse{data=MessageResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Permission denied"
// @Failure 404 {object} responses.ErrorResponse "Message not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /messages/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	var req UpdateMessageRequest

	// Parse message ID
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid message ID", "Message ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.UpdateMessageRequest{
		MessageID: messageID,
		UserID:    userUUID,
		Content:   req.Content,
		Metadata:  req.Metadata,
	}

	// Update message
	message, err := h.messageService.UpdateMessage(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "message not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Message not found", "The specified message does not exist.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to edit this message.")
			return
		case strings.Contains(err.Error(), "edit time expired"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Edit time expired", "The time limit for editing this message has expired.")
			return
		case strings.Contains(err.Error(), "cannot edit"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Cannot edit message", "This type of message cannot be edited.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to update message", "Unable to update message.")
			return
		}
	}

	// Convert to response format
	messageResponse := h.convertToMessageResponse(message, userUUID)

	// Log message update
	middleware.LogInfo(c, "Message updated", gin.H{
		"message_id": messageID,
		"user_id":    userUUID,
	})

	responses.DataResponse(c, messageResponse)
}

// @Summary Delete a message
// @Description Delete an existing message
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Permission denied"
// @Failure 404 {object} responses.ErrorResponse "Message not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	// Parse message ID
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid message ID", "Message ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Delete message
	err = h.messageService.DeleteMessage(c.Request.Context(), messageID, userUUID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "message not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Message not found", "The specified message does not exist.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to delete this message.")
			return
		case strings.Contains(err.Error(), "cannot delete"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Cannot delete message", "This message cannot be deleted.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete message", "Unable to delete message.")
			return
		}
	}

	// Log message deletion
	middleware.LogInfo(c, "Message deleted", gin.H{
		"message_id": messageID,
		"user_id":    userUUID,
	})

	responses.SuccessMessageResponse(c, "Message deleted successfully")
}

// @Summary Mark messages as read
// @Description Mark one or more messages as read
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body MarkAsReadRequest true "Messages to mark as read"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /messages/read [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	var req MarkAsReadRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Mark messages as read
	err := h.messageService.MarkAsRead(c.Request.Context(), req.MessageIDs, userUUID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark as read", "Unable to mark messages as read.")
		return
	}

	// Log messages marked as read
	middleware.LogInfo(c, "Messages marked as read", gin.H{
		"user_id":      userUUID,
		"message_count": len(req.MessageIDs),
	})

	responses.SuccessMessageResponse(c, "Messages marked as read successfully")
}

// @Summary Add reaction to message
// @Description Add an emoji reaction to a message
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Param request body AddReactionRequest true "Reaction to add"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse "Message not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /messages/{id}/reactions [post]
func (h *MessageHandler) AddReaction(c *gin.Context) {
	var req AddReactionRequest

	// Parse message ID
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid message ID", "Message ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Add reaction
	err = h.messageService.AddReaction(c.Request.Context(), messageID, userUUID, req.Emoji)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "message not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Message not found", "The specified message does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have permission to react to this message.")
			return
		case strings.Contains(err.Error(), "invalid emoji"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Invalid emoji", "The provided emoji is not valid.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to add reaction", "Unable to add reaction.")
			return
		}
	}

	// Log reaction added
	middleware.LogInfo(c, "Reaction added", gin.H{
		"message_id": messageID,
		"user_id":    userUUID,
		"emoji":      req.Emoji,
	})

	responses.SuccessMessageResponse(c, "Reaction added successfully")
}

// @Summary Remove reaction from message
// @Description Remove a specific emoji reaction from a message
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Param emoji path string true "Emoji to remove"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse "Message or reaction not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /messages/{id}/reactions/{emoji} [delete]
func (h *MessageHandler) RemoveReaction(c *gin.Context) {
	// Parse message ID
	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid message ID", "Message ID must be a valid UUID.")
		return
	}

	// Get emoji from path
	emoji := c.Param("emoji")
	if emoji == "" {
		responses.ErrorResponse(c, http.StatusBadRequest, "Missing emoji", "Emoji is required.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Remove reaction
	err = h.messageService.RemoveReaction(c.Request.Context(), messageID, userUUID, emoji)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "message not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Message not found", "The specified message does not exist.")
			return
		case strings.Contains(err.Error(), "reaction not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Reaction not found", "You haven't reacted with this emoji.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove reaction", "Unable to remove reaction.")
			return
		}
	}

	// Log reaction removed
	middleware.LogInfo(c, "Reaction removed", gin.H{
		"message_id": messageID,
		"user_id":    userUUID,
		"emoji":      emoji,
	})

	responses.SuccessMessageResponse(c, "Reaction removed successfully")
}

// Helper functions

func (h *MessageHandler) convertToMessageResponse(message interface{}, userID uuid.UUID) MessageResponse {
	// This is a simplified conversion - in a real implementation,
	// you would properly convert from your message model to the response
	return MessageResponse{
		// Populate fields from message model
		// This is just a placeholder structure
	}
}

// RegisterMessageRoutes registers all message-related routes
func RegisterMessageRoutes(router *gin.RouterGroup, messageService *services.MessageService) {
	handler := NewMessageHandler(messageService)

	// Protected message routes
	dialogs := router.Group("/dialogs")
	dialogs.Use(middleware.AuthRequired())
	{
		dialogs.GET("/:id/messages", handler.GetMessages)
		dialogs.POST("/:id/messages", handler.SendMessage)
		dialogs.POST("/:id/messages/upload", handler.UploadAndSendFile)
	}

	messages := router.Group("/messages")
	messages.Use(middleware.AuthRequired())
	{
		messages.PUT("/:id", handler.UpdateMessage)
		messages.DELETE("/:id", handler.DeleteMessage)
		messages.POST("/read", handler.MarkAsRead)
		messages.POST("/:id/reactions", handler.AddReaction)
		messages.DELETE("/:id/reactions/:emoji", handler.RemoveReaction)
	}
}

// RegisterMessageRoutesWithMiddleware registers message routes with custom middleware
func RegisterMessageRoutesWithMiddleware(
	router *gin.RouterGroup,
	messageService *services.MessageService,
	middlewares ...gin.HandlerFunc,
) {
	handler := NewMessageHandler(messageService)

	// Protected message routes with middleware
	allMiddlewares := append(middlewares, middleware.AuthRequired())

	dialogs := router.Group("/dialogs")
	dialogs.Use(allMiddlewares...)
	{
		dialogs.GET("/:id/messages", handler.GetMessages)
		dialogs.POST("/:id/messages", handler.SendMessage)
		dialogs.POST("/:id/messages/upload", handler.UploadAndSendFile)
	}

	messages := router.Group("/messages")
	messages.Use(allMiddlewares...)
	{
		messages.PUT("/:id", handler.UpdateMessage)
		messages.DELETE("/:id", handler.DeleteMessage)
		messages.POST("/read", handler.MarkAsRead)
		messages.POST("/:id/reactions", handler.AddReaction)
		messages.DELETE("/:id/reactions/:emoji", handler.RemoveReaction)
	}
}