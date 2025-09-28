package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
	sharedModels "tchat.dev/shared/models"
	"tchat.dev/shared/utils"
)

// MessagingHandler handles messaging HTTP requests and WebSocket connections
type MessagingHandler struct {
	messagingService services.MessagingService
	validator        *utils.Validator
	upgrader         websocket.Upgrader
	connections      map[uuid.UUID]*WebSocketConnection
	broadcast        chan *BroadcastMessage
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	UserID     uuid.UUID
	Connection *websocket.Conn
	Send       chan []byte
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	Type      string      `json:"type"`
	DialogID  *uuid.UUID  `json:"dialog_id,omitempty"`
	UserIDs   []uuid.UUID `json:"user_ids,omitempty"`
	Data      interface{} `json:"data"`
}

// NewMessagingHandler creates a new messaging handler
func NewMessagingHandler(messagingService services.MessagingService) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
		validator:        utils.NewValidator(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		connections: make(map[uuid.UUID]*WebSocketConnection),
		broadcast:   make(chan *BroadcastMessage, 256),
	}
}

// Start starts the WebSocket hub
func (h *MessagingHandler) Start() {
	go h.handleBroadcasts()
}

// RegisterRoutes registers messaging routes
func (h *MessagingHandler) RegisterRoutes(router *mux.Router) {
	// Messaging routes
	messaging := router.PathPrefix("/messaging").Subrouter()

	// WebSocket endpoint
	messaging.HandleFunc("/ws", h.HandleWebSocket)

	// Dialog endpoints
	messaging.HandleFunc("/dialogs", h.CreateDialog).Methods("POST")
	messaging.HandleFunc("/dialogs", h.GetDialogs).Methods("GET")
	messaging.HandleFunc("/dialogs/{id}", h.GetDialog).Methods("GET")
	messaging.HandleFunc("/dialogs/{id}", h.UpdateDialog).Methods("PUT")
	messaging.HandleFunc("/dialogs/{id}", h.DeleteDialog).Methods("DELETE")

	// Dialog participant endpoints
	messaging.HandleFunc("/dialogs/{id}/participants", h.GetParticipants).Methods("GET")
	messaging.HandleFunc("/dialogs/{id}/participants", h.AddParticipant).Methods("POST")
	messaging.HandleFunc("/dialogs/{id}/participants/{user_id}", h.RemoveParticipant).Methods("DELETE")

	// Message endpoints
	messaging.HandleFunc("/dialogs/{id}/messages", h.SendMessage).Methods("POST")
	messaging.HandleFunc("/dialogs/{id}/messages", h.GetMessages).Methods("GET")
	messaging.HandleFunc("/messages/{id}", h.GetMessage).Methods("GET")
	messaging.HandleFunc("/messages/{id}", h.UpdateMessage).Methods("PUT")
	messaging.HandleFunc("/messages/{id}", h.DeleteMessage).Methods("DELETE")

	// Message interaction endpoints
	messaging.HandleFunc("/messages/{id}/reactions", h.AddReaction).Methods("POST")
	messaging.HandleFunc("/messages/{id}/reactions", h.RemoveReaction).Methods("DELETE")
	messaging.HandleFunc("/messages/{id}/read", h.MarkAsRead).Methods("POST")

	// Unread count endpoint
	messaging.HandleFunc("/dialogs/{id}/unread", h.GetUnreadCount).Methods("GET")
}

// HandleWebSocket handles WebSocket connections
func (h *MessagingHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user from context (should be set by auth middleware)
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	// Create WebSocket connection
	wsConn := &WebSocketConnection{
		UserID:     user.ID,
		Connection: conn,
		Send:       make(chan []byte, 256),
	}

	// Register connection
	h.connections[user.ID] = wsConn

	// Start goroutines
	go h.handleWebSocketRead(wsConn)
	go h.handleWebSocketWrite(wsConn)

	// Send welcome message
	welcome := map[string]interface{}{
		"type":    "welcome",
		"user_id": user.ID,
		"message": "Connected to messaging service",
	}
	h.sendToConnection(wsConn, welcome)
}

// CreateDialog handles dialog creation
func (h *MessagingHandler) CreateDialog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req CreateDialogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateCreateDialogRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Create service request
	serviceReq := &services.CreateDialogRequest{
		Type:           req.Type,
		Title:          "",
		Description:    "",
		CreatorID:      user.ID,
		ParticipantIDs: req.ParticipantIDs,
		Settings:       models.DialogSettings{},
	}

	// Set title (Name -> Title mapping)
	if req.Name != nil {
		serviceReq.Title = *req.Name
	}

	// Set description if provided
	if req.Description != nil {
		serviceReq.Description = *req.Description
	}

	// Set settings if provided
	if req.Settings != nil {
		serviceReq.Settings = *req.Settings
	}

	// Create dialog
	dialog, err := h.messagingService.CreateDialog(ctx, serviceReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to create dialog", err)
		return
	}

	// Get participant IDs from dialog for broadcasting
	participants, err := h.messagingService.GetDialogParticipants(ctx, dialog.ID)
	if err == nil {
		var participantIDs []uuid.UUID
		for _, participant := range participants {
			participantIDs = append(participantIDs, participant.UserID)
		}
		// Broadcast dialog creation to participants
		h.broadcastDialogEvent("dialog_created", dialog.ID, participantIDs, dialog)
	}

	// Success response
	h.respondSuccess(w, http.StatusCreated, "Dialog created successfully", h.sanitizeDialog(dialog))
}

// SendMessage handles message sending
func (h *MessagingHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get dialog ID from URL
	vars := mux.Vars(r)
	dialogIDStr := vars["id"]
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid dialog ID", err)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateSendMessageRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Create service request
	serviceReq := &services.SendMessageRequest{
		DialogID:  dialogID,
		SenderID:  user.ID,
		Type:      req.Type,
		Content:   h.convertContentToString(req.Content),
		ReplyToID: req.ReplyToID,
	}

	// Send message
	message, err := h.messagingService.SendMessage(ctx, serviceReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to send message", err)
		return
	}

	// Get dialog participants for broadcasting
	participants, err := h.messagingService.GetDialogParticipants(ctx, dialogID)
	if err == nil {
		var participantIDs []uuid.UUID
		for _, participant := range participants {
			participantIDs = append(participantIDs, participant.UserID)
		}

		// Broadcast message to participants
		h.broadcastMessageEvent("message_sent", dialogID, participantIDs, message)
	}

	// Success response
	h.respondSuccess(w, http.StatusCreated, "Message sent successfully", h.sanitizeMessage(message))
}

// GetMessages handles getting dialog messages
func (h *MessagingHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get dialog ID from URL
	vars := mux.Vars(r)
	dialogIDStr := vars["id"]
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid dialog ID", err)
		return
	}

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Get messages
	messages, err := h.messagingService.GetMessages(ctx, dialogID, user.ID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to get messages", err)
		return
	}

	// Sanitize messages
	var sanitizedMessages []map[string]interface{}
	for _, message := range messages {
		sanitizedMessages = append(sanitizedMessages, h.sanitizeMessage(message))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Messages retrieved successfully", map[string]interface{}{
		"messages": sanitizedMessages,
		"count":    len(sanitizedMessages),
		"limit":    limit,
		"offset":   offset,
	})
}

// GetDialogs handles getting user dialogs
func (h *MessagingHandler) GetDialogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Get dialogs
	dialogs, err := h.messagingService.GetUserDialogs(ctx, user.ID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get dialogs", err)
		return
	}

	// Sanitize dialogs
	var sanitizedDialogs []map[string]interface{}
	for _, dialog := range dialogs {
		sanitizedDialogs = append(sanitizedDialogs, h.sanitizeDialog(dialog))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Dialogs retrieved successfully", map[string]interface{}{
		"dialogs": sanitizedDialogs,
		"count":   len(sanitizedDialogs),
		"limit":   limit,
		"offset":  offset,
	})
}

// AddReaction handles adding reactions to messages
func (h *MessagingHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get message ID from URL
	vars := mux.Vars(r)
	messageIDStr := vars["id"]
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid message ID", err)
		return
	}

	var req AddReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateAddReactionRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Add reaction
	if err := h.messagingService.AddReaction(ctx, messageID, user.ID, req.Reaction); err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to add reaction", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Reaction added successfully", map[string]interface{}{
		"message_id": messageID,
		"reaction":   req.Reaction,
		"user_id":    user.ID,
	})
}

// MarkAsRead handles marking messages as read
func (h *MessagingHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get message ID from URL
	vars := mux.Vars(r)
	messageIDStr := vars["id"]
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid message ID", err)
		return
	}

	// Mark as read
	if err := h.messagingService.MarkAsRead(ctx, messageID, user.ID); err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to mark as read", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Message marked as read", map[string]interface{}{
		"message_id": messageID,
		"read":       true,
	})
}

// GetUnreadCount handles getting unread message count
func (h *MessagingHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get dialog ID from URL
	vars := mux.Vars(r)
	dialogIDStr := vars["id"]
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid dialog ID", err)
		return
	}

	// Get unread count
	count, err := h.messagingService.GetUnreadCount(ctx, dialogID, user.ID)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to get unread count", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Unread count retrieved successfully", map[string]interface{}{
		"dialog_id":    dialogID,
		"unread_count": count,
	})
}

// Request/Response types
type CreateDialogRequest struct {
	Type           models.DialogType      `json:"type"`
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	ParticipantIDs []uuid.UUID            `json:"participant_ids,omitempty"`
	IsPublic       bool                   `json:"is_public"`
	Settings       *models.DialogSettings `json:"settings,omitempty"`
}

type SendMessageRequest struct {
	Type      models.MessageType    `json:"type"`
	Content   models.MessageContent `json:"content"`
	ParentID  *uuid.UUID            `json:"parent_id,omitempty"`
	Mentions  []uuid.UUID           `json:"mentions,omitempty"`
	ReplyToID *uuid.UUID            `json:"reply_to_id,omitempty"`
}

type AddReactionRequest struct {
	Reaction string `json:"reaction"`
}

// Standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// WebSocket message types
type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WebSocket handling methods
func (h *MessagingHandler) handleWebSocketRead(conn *WebSocketConnection) {
	defer func() {
		delete(h.connections, conn.UserID)
		conn.Connection.Close()
	}()

	conn.Connection.SetReadLimit(512)
	conn.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.Connection.SetPongHandler(func(string) error {
		conn.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		// Handle incoming WebSocket messages
		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		h.handleWebSocketMessage(conn, &wsMsg)
	}
}

func (h *MessagingHandler) handleWebSocketWrite(conn *WebSocketConnection) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			conn.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			conn.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *MessagingHandler) handleWebSocketMessage(conn *WebSocketConnection, msg *WebSocketMessage) {
	switch msg.Type {
	case "ping":
		h.sendToConnection(conn, map[string]interface{}{
			"type": "pong",
		})

	case "typing":
		// Handle typing indicators
		if data, ok := msg.Data.(map[string]interface{}); ok {
			if dialogIDStr, ok := data["dialog_id"].(string); ok {
				if dialogID, err := uuid.Parse(dialogIDStr); err == nil {
					h.broadcastTypingEvent(dialogID, conn.UserID, true)
				}
			}
		}

	case "stop_typing":
		// Handle stop typing
		if data, ok := msg.Data.(map[string]interface{}); ok {
			if dialogIDStr, ok := data["dialog_id"].(string); ok {
				if dialogID, err := uuid.Parse(dialogIDStr); err == nil {
					h.broadcastTypingEvent(dialogID, conn.UserID, false)
				}
			}
		}
	}
}

func (h *MessagingHandler) handleBroadcasts() {
	for {
		select {
		case broadcast := <-h.broadcast:
			h.processBroadcast(broadcast)
		}
	}
}

func (h *MessagingHandler) processBroadcast(broadcast *BroadcastMessage) {
	message, err := json.Marshal(broadcast)
	if err != nil {
		return
	}

	// Send to specific users if UserIDs provided
	if len(broadcast.UserIDs) > 0 {
		for _, userID := range broadcast.UserIDs {
			if conn, exists := h.connections[userID]; exists {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(h.connections, userID)
				}
			}
		}
	} else {
		// Broadcast to all connected users
		for userID, conn := range h.connections {
			select {
			case conn.Send <- message:
			default:
				close(conn.Send)
				delete(h.connections, userID)
			}
		}
	}
}

func (h *MessagingHandler) sendToConnection(conn *WebSocketConnection, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	select {
	case conn.Send <- message:
	default:
		close(conn.Send)
		delete(h.connections, conn.UserID)
	}
}

func (h *MessagingHandler) broadcastMessageEvent(eventType string, dialogID uuid.UUID, userIDs []uuid.UUID, message *models.Message) {
	broadcast := &BroadcastMessage{
		Type:     eventType,
		DialogID: &dialogID,
		UserIDs:  userIDs,
		Data: map[string]interface{}{
			"message": h.sanitizeMessage(message),
		},
	}

	select {
	case h.broadcast <- broadcast:
	default:
	}
}

func (h *MessagingHandler) broadcastDialogEvent(eventType string, dialogID uuid.UUID, userIDs []uuid.UUID, dialog *models.Dialog) {
	broadcast := &BroadcastMessage{
		Type:     eventType,
		DialogID: &dialogID,
		UserIDs:  userIDs,
		Data: map[string]interface{}{
			"dialog": h.sanitizeDialog(dialog),
		},
	}

	select {
	case h.broadcast <- broadcast:
	default:
	}
}

func (h *MessagingHandler) broadcastTypingEvent(dialogID, userID uuid.UUID, isTyping bool) {
	broadcast := &BroadcastMessage{
		Type:     "typing",
		DialogID: &dialogID,
		Data: map[string]interface{}{
			"user_id":   userID,
			"is_typing": isTyping,
		},
	}

	select {
	case h.broadcast <- broadcast:
	default:
	}
}

// Validation methods
func (h *MessagingHandler) validateCreateDialogRequest(req *CreateDialogRequest) error {
	h.validator.Reset()

	if !req.Type.IsValid() {
		h.validator.AddError("type", "invalid dialog type")
	}

	// Name required for group/channel types
	if (req.Type == models.DialogTypeGroup || req.Type == models.DialogTypeChannel) && req.Name == nil {
		h.validator.AddError("name", "name is required for group and channel dialogs")
	}

	if req.Name != nil {
		h.validator.MinLength("name", *req.Name, 1).MaxLength("name", *req.Name, 100)
	}

	if req.Description != nil {
		h.validator.MaxLength("description", *req.Description, 500)
	}

	return h.validator.GetError()
}

func (h *MessagingHandler) validateSendMessageRequest(req *SendMessageRequest) error {
	h.validator.Reset()

	if !req.Type.IsValid() {
		h.validator.AddError("type", "invalid message type")
	}

	// Validate content based on message type
	switch req.Type {
	case models.MessageTypeText:
		if textContent, exists := req.Content["text"]; !exists || textContent == "" {
			h.validator.AddError("content.text", "text content is required for text messages")
		}
	case models.MessageTypeFile, models.MessageTypeImage, models.MessageTypeVideo:
		if _, exists := req.Content["file"]; !exists {
			h.validator.AddError("content.file", "file content is required for file messages")
		}
	case models.MessageTypeVoice:
		if _, exists := req.Content["voice"]; !exists {
			h.validator.AddError("content.voice", "voice content is required for voice messages")
		}
	}

	return h.validator.GetError()
}

func (h *MessagingHandler) validateAddReactionRequest(req *AddReactionRequest) error {
	h.validator.Reset()

	h.validator.Required("reaction", req.Reaction).MinLength("reaction", req.Reaction, 1).MaxLength("reaction", req.Reaction, 10)

	return h.validator.GetError()
}

// Utility methods
func (h *MessagingHandler) respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *MessagingHandler) respondError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errorData interface{}
	if err != nil {
		errorData = err.Error()
	}

	response := APIResponse{
		Success: false,
		Message: message,
		Error:   errorData,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *MessagingHandler) getUserFromContext(ctx context.Context) *sharedModels.User {
	if user, ok := ctx.Value("user").(*sharedModels.User); ok {
		return user
	}
	return nil
}

func (h *MessagingHandler) sanitizeDialog(dialog *models.Dialog) map[string]interface{} {
	if dialog == nil {
		return nil
	}

	return map[string]interface{}{
		"id":               dialog.ID,
		"type":             dialog.Type,
		"name":             dialog.Name,
		"title":            dialog.Title,
		"description":      dialog.Description,
		"participant_count": dialog.ParticipantCount,
		"is_public":        dialog.Settings.IsPublic,
		"last_message_id":  dialog.LastMessageID,
		"settings":         dialog.Settings,
		"is_archived":      dialog.IsArchived,
		"is_muted":         dialog.IsMuted,
		"created_at":       dialog.CreatedAt,
		"updated_at":       dialog.UpdatedAt,
	}
}

func (h *MessagingHandler) sanitizeMessage(message *models.Message) map[string]interface{} {
	if message == nil {
		return nil
	}

	return map[string]interface{}{
		"id":         message.ID,
		"dialog_id":  message.DialogID,
		"sender_id":  message.SenderID,
		"type":       message.Type,
		"content":    message.Content,
		"reply_to_id": message.ReplyToID,
		"mentions":   message.Mentions,
		"reactions":  message.Reactions,
		"is_edited":  message.IsEdited,
		"is_deleted": message.IsDeleted,
		"sent_at":    message.SentAt,
		"edited_at":  message.EditedAt,
		"deleted_at": message.DeletedAt,
		"created_at": message.CreatedAt,
		"updated_at": message.UpdatedAt,
	}
}

func (h *MessagingHandler) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 50 // default for messages
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset = 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

// convertContentToString converts MessageContent to string for service layer
func (h *MessagingHandler) convertContentToString(content models.MessageContent) string {
	if textContent, exists := content["text"]; exists {
		if textStr, ok := textContent.(string); ok {
			return textStr
		}
	}
	// For non-text content, return empty string (content will be in other fields)
	return ""
}

// Additional placeholder methods that would be implemented
func (h *MessagingHandler) GetDialog(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) UpdateDialog(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) DeleteDialog(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *MessagingHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	// Implementation would go here
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}