package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"tchat.dev/messaging/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

type WebSocketHandler struct {
	messageService  *services.MessageService
	presenceService *services.PresenceService
	upgrader        websocket.Upgrader
	connections     *ConnectionManager
	eventBus        *EventBus
}

type ConnectionManager struct {
	connections map[uuid.UUID]*WebSocketConnection
	userSessions map[uuid.UUID]map[string]*WebSocketConnection // userID -> sessionID -> connection
	dialogRooms  map[uuid.UUID]map[uuid.UUID]*WebSocketConnection // dialogID -> userID -> connection
	mu           sync.RWMutex
}

type WebSocketConnection struct {
	ID       string
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Send     chan []byte
	Handler  *WebSocketHandler
	Context  context.Context
	Cancel   context.CancelFunc
	LastPing time.Time
	Dialogs  map[uuid.UUID]bool // Set of dialog IDs user is subscribed to
	mu       sync.RWMutex
}

type EventBus struct {
	subscribers map[string][]chan *WebSocketEvent
	mu          sync.RWMutex
}

type WebSocketEvent struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	DialogID  *uuid.UUID             `json:"dialog_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// WebSocket message types from client
type ClientMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	ID   string                 `json:"id,omitempty"` // For request-response pattern
}

// WebSocket message types to client
type ServerMessage struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data,omitempty"`
	ID        string                 `json:"id,omitempty"` // For request-response pattern
	Error     *ErrorInfo             `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Client message types
const (
	ClientMessageTypeAuth          = "auth"
	ClientMessageTypePing          = "ping"
	ClientMessageTypeSubscribe     = "subscribe"
	ClientMessageTypeUnsubscribe   = "unsubscribe"
	ClientMessageTypeSendMessage   = "send_message"
	ClientMessageTypeTyping        = "typing"
	ClientMessageTypeMarkRead      = "mark_read"
	ClientMessageTypeUpdatePresence = "update_presence"
	ClientMessageTypeJoinDialog    = "join_dialog"
	ClientMessageTypeLeaveDialog   = "leave_dialog"
)

// Server message types
const (
	ServerMessageTypeAuthResult      = "auth_result"
	ServerMessageTypePong           = "pong"
	ServerMessageTypeSubscribed     = "subscribed"
	ServerMessageTypeUnsubscribed   = "unsubscribed"
	ServerMessageTypeNewMessage     = "new_message"
	ServerMessageTypeMessageUpdated = "message_updated"
	ServerMessageTypeMessageDeleted = "message_deleted"
	ServerMessageTypeTyping         = "typing"
	ServerMessageTypePresenceUpdate = "presence_update"
	ServerMessageTypeDialogUpdate   = "dialog_update"
	ServerMessageTypeError          = "error"
	ServerMessageTypeUserJoined     = "user_joined"
	ServerMessageTypeUserLeft       = "user_left"
	ServerMessageTypeDeliveryStatus = "delivery_status"
	ServerMessageTypeReadReceipt    = "read_receipt"
)

func NewWebSocketHandler(
	messageService *services.MessageService,
	presenceService *services.PresenceService,
) *WebSocketHandler {
	return &WebSocketHandler{
		messageService:  messageService,
		presenceService: presenceService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections: NewConnectionManager(),
		eventBus:    NewEventBus(),
	}
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections:  make(map[uuid.UUID]*WebSocketConnection),
		userSessions: make(map[uuid.UUID]map[string]*WebSocketConnection),
		dialogRooms:  make(map[uuid.UUID]map[uuid.UUID]*WebSocketConnection),
	}
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan *WebSocketEvent),
	}
}

// @Summary WebSocket endpoint for real-time messaging
// @Description Establish WebSocket connection for real-time messaging and presence
// @Tags websocket
// @Param Authorization header string true "Bearer token for authentication"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /websocket [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		middleware.LogError(c, "WebSocket upgrade failed", gin.H{"error": err.Error()})
		responses.ErrorResponse(c, http.StatusBadRequest, "WebSocket upgrade failed", err.Error())
		return
	}

	// Create connection context
	ctx, cancel := context.WithCancel(c.Request.Context())

	// Create WebSocket connection
	wsConn := &WebSocketConnection{
		ID:       generateConnectionID(),
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Handler:  h,
		Context:  ctx,
		Cancel:   cancel,
		LastPing: time.Now(),
		Dialogs:  make(map[uuid.UUID]bool),
	}

	// Start connection goroutines
	go wsConn.readPump()
	go wsConn.writePump()

	// Log connection established
	middleware.LogInfo(c, "WebSocket connection established", gin.H{
		"connection_id": wsConn.ID,
		"remote_addr":   c.Request.RemoteAddr,
	})
}

// Read pump handles incoming messages from the WebSocket
func (c *WebSocketConnection) readPump() {
	defer func() {
		c.Handler.connections.removeConnection(c)
		c.Conn.Close()
		c.Cancel()
	}()

	// Set read deadline and configure connection
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.LastPing = time.Now()
		return nil
	})

	for {
		select {
		case <-c.Context.Done():
			return
		default:
			_, messageBytes, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket read error: %v", err)
				}
				return
			}

			var message ClientMessage
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				c.sendError("invalid_message", "Invalid message format", err.Error())
				continue
			}

			c.handleMessage(&message)
		}
	}
}

// Write pump handles outgoing messages to the WebSocket
func (c *WebSocketConnection) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-c.Context.Done():
			return
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Handle incoming client messages
func (c *WebSocketConnection) handleMessage(message *ClientMessage) {
	switch message.Type {
	case ClientMessageTypeAuth:
		c.handleAuth(message)
	case ClientMessageTypePing:
		c.handlePing(message)
	case ClientMessageTypeSubscribe:
		c.handleSubscribe(message)
	case ClientMessageTypeUnsubscribe:
		c.handleUnsubscribe(message)
	case ClientMessageTypeSendMessage:
		c.handleSendMessage(message)
	case ClientMessageTypeTyping:
		c.handleTyping(message)
	case ClientMessageTypeMarkRead:
		c.handleMarkRead(message)
	case ClientMessageTypeUpdatePresence:
		c.handleUpdatePresence(message)
	case ClientMessageTypeJoinDialog:
		c.handleJoinDialog(message)
	case ClientMessageTypeLeaveDialog:
		c.handleLeaveDialog(message)
	default:
		c.sendError("unknown_message_type", "Unknown message type", fmt.Sprintf("Type: %s", message.Type))
	}
}

// Authentication handler
func (c *WebSocketConnection) handleAuth(message *ClientMessage) {
	token, ok := message.Data["token"].(string)
	if !ok {
		c.sendError("missing_token", "Authentication token is required", "")
		return
	}

	// Validate token and get user ID
	// This would integrate with your auth service
	userID, err := c.validateToken(token)
	if err != nil {
		c.sendError("invalid_token", "Invalid authentication token", err.Error())
		return
	}

	// Set user ID and register connection
	c.UserID = userID
	c.Handler.connections.addConnection(c)

	// Update user presence
	presenceReq := &services.UpdatePresenceRequest{
		UserID:   userID,
		Status:   "online",
		Platform: c.getPlatformFromUserAgent(),
	}
	c.Handler.presenceService.UpdatePresence(c.Context, presenceReq)

	// Send authentication success
	c.sendResponse(ServerMessageTypeAuthResult, map[string]interface{}{
		"success": true,
		"user_id": userID,
	}, message.ID)

	// Subscribe to user's dialog events
	c.subscribeToUserDialogs()
}

// Ping handler
func (c *WebSocketConnection) handlePing(message *ClientMessage) {
	c.LastPing = time.Now()
	c.sendResponse(ServerMessageTypePong, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}, message.ID)
}

// Subscribe to dialog events
func (c *WebSocketConnection) handleSubscribe(message *ClientMessage) {
	if c.UserID == uuid.Nil {
		c.sendError("not_authenticated", "Authentication required", "")
		return
	}

	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	// Verify user has access to dialog
	hasAccess, err := c.verifyDialogAccess(dialogID)
	if err != nil || !hasAccess {
		c.sendError("access_denied", "Access denied to dialog", "")
		return
	}

	// Subscribe to dialog
	c.mu.Lock()
	c.Dialogs[dialogID] = true
	c.mu.Unlock()

	c.Handler.connections.addToDialogRoom(dialogID, c.UserID, c)

	c.sendResponse(ServerMessageTypeSubscribed, map[string]interface{}{
		"dialog_id": dialogID,
	}, message.ID)
}

// Unsubscribe from dialog events
func (c *WebSocketConnection) handleUnsubscribe(message *ClientMessage) {
	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	// Unsubscribe from dialog
	c.mu.Lock()
	delete(c.Dialogs, dialogID)
	c.mu.Unlock()

	c.Handler.connections.removeFromDialogRoom(dialogID, c.UserID)

	c.sendResponse(ServerMessageTypeUnsubscribed, map[string]interface{}{
		"dialog_id": dialogID,
	}, message.ID)
}

// Send message handler
func (c *WebSocketConnection) handleSendMessage(message *ClientMessage) {
	if c.UserID == uuid.Nil {
		c.sendError("not_authenticated", "Authentication required", "")
		return
	}

	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	content, ok := message.Data["content"].(string)
	if !ok {
		c.sendError("missing_content", "Message content is required", "")
		return
	}

	messageType, _ := message.Data["type"].(string)
	if messageType == "" {
		messageType = "text"
	}

	// Build service request
	serviceReq := &services.SendMessageRequest{
		DialogID: dialogID,
		SenderID: c.UserID,
		Type:     messageType,
		Content:  content,
	}

	// Add optional fields
	if replyToIDStr, ok := message.Data["reply_to_id"].(string); ok {
		if replyToID, err := uuid.Parse(replyToIDStr); err == nil {
			serviceReq.ReplyToID = &replyToID
		}
	}

	if clientMessageID, ok := message.Data["client_message_id"].(string); ok {
		serviceReq.ClientMessageID = clientMessageID
	}

	// Send message
	sentMessage, err := c.Handler.messageService.SendMessage(c.Context, serviceReq)
	if err != nil {
		c.sendError("send_failed", "Failed to send message", err.Error())
		return
	}

	// Broadcast message to dialog participants
	c.Handler.broadcastToDialog(dialogID, ServerMessageTypeNewMessage, map[string]interface{}{
		"message": c.convertMessageToMap(sentMessage),
	})

	// Send delivery confirmation to sender
	c.sendResponse(ServerMessageTypeDeliveryStatus, map[string]interface{}{
		"message_id":       sentMessage.ID,
		"client_message_id": serviceReq.ClientMessageID,
		"status":           "sent",
		"timestamp":        sentMessage.CreatedAt,
	}, message.ID)
}

// Typing indicator handler
func (c *WebSocketConnection) handleTyping(message *ClientMessage) {
	if c.UserID == uuid.Nil {
		c.sendError("not_authenticated", "Authentication required", "")
		return
	}

	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	isTyping, _ := message.Data["is_typing"].(bool)

	// Broadcast typing status to other dialog participants
	c.Handler.broadcastToDialogExcept(dialogID, c.UserID, ServerMessageTypeTyping, map[string]interface{}{
		"dialog_id":  dialogID,
		"user_id":    c.UserID,
		"is_typing":  isTyping,
		"timestamp":  time.Now(),
	})
}

// Mark read handler
func (c *WebSocketConnection) handleMarkRead(message *ClientMessage) {
	if c.UserID == uuid.Nil {
		c.sendError("not_authenticated", "Authentication required", "")
		return
	}

	messageIDsInterface, ok := message.Data["message_ids"].([]interface{})
	if !ok {
		c.sendError("missing_message_ids", "Message IDs are required", "")
		return
	}

	var messageIDs []uuid.UUID
	for _, idInterface := range messageIDsInterface {
		if idStr, ok := idInterface.(string); ok {
			if id, err := uuid.Parse(idStr); err == nil {
				messageIDs = append(messageIDs, id)
			}
		}
	}

	if len(messageIDs) == 0 {
		c.sendError("invalid_message_ids", "Valid message IDs are required", "")
		return
	}

	// Mark messages as read
	err := c.Handler.messageService.MarkAsRead(c.Context, messageIDs, c.UserID)
	if err != nil {
		c.sendError("mark_read_failed", "Failed to mark messages as read", err.Error())
		return
	}

	// Get dialog ID from first message (assuming all messages are from same dialog)
	firstMessage, err := c.Handler.messageService.GetMessageByID(c.Context, messageIDs[0])
	if err == nil {
		// Broadcast read receipt to dialog participants
		c.Handler.broadcastToDialogExcept(firstMessage.DialogID, c.UserID, ServerMessageTypeReadReceipt, map[string]interface{}{
			"message_ids": messageIDs,
			"user_id":     c.UserID,
			"read_at":     time.Now(),
		})
	}

	c.sendResponse("mark_read_success", map[string]interface{}{
		"message_ids": messageIDs,
		"read_at":     time.Now(),
	}, message.ID)
}

// Update presence handler
func (c *WebSocketConnection) handleUpdatePresence(message *ClientMessage) {
	if c.UserID == uuid.Nil {
		c.sendError("not_authenticated", "Authentication required", "")
		return
	}

	status, ok := message.Data["status"].(string)
	if !ok {
		status = "online"
	}

	// Update presence
	presenceReq := &services.UpdatePresenceRequest{
		UserID:   c.UserID,
		Status:   status,
		Platform: c.getPlatformFromUserAgent(),
	}

	presence, err := c.Handler.presenceService.UpdatePresence(c.Context, presenceReq)
	if err != nil {
		c.sendError("presence_update_failed", "Failed to update presence", err.Error())
		return
	}

	// Broadcast presence update to relevant users
	c.Handler.broadcastPresenceUpdate(c.UserID, presence)

	c.sendResponse("presence_updated", map[string]interface{}{
		"status":    presence.Status,
		"timestamp": presence.UpdatedAt,
	}, message.ID)
}

// Join dialog handler
func (c *WebSocketConnection) handleJoinDialog(message *ClientMessage) {
	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	// Subscribe to dialog automatically
	c.mu.Lock()
	c.Dialogs[dialogID] = true
	c.mu.Unlock()

	c.Handler.connections.addToDialogRoom(dialogID, c.UserID, c)

	// Notify other participants that user joined
	c.Handler.broadcastToDialogExcept(dialogID, c.UserID, ServerMessageTypeUserJoined, map[string]interface{}{
		"dialog_id": dialogID,
		"user_id":   c.UserID,
		"timestamp": time.Now(),
	})

	c.sendResponse("dialog_joined", map[string]interface{}{
		"dialog_id": dialogID,
	}, message.ID)
}

// Leave dialog handler
func (c *WebSocketConnection) handleLeaveDialog(message *ClientMessage) {
	dialogIDStr, ok := message.Data["dialog_id"].(string)
	if !ok {
		c.sendError("missing_dialog_id", "Dialog ID is required", "")
		return
	}

	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		c.sendError("invalid_dialog_id", "Invalid dialog ID format", err.Error())
		return
	}

	// Unsubscribe from dialog
	c.mu.Lock()
	delete(c.Dialogs, dialogID)
	c.mu.Unlock()

	c.Handler.connections.removeFromDialogRoom(dialogID, c.UserID)

	// Notify other participants that user left
	c.Handler.broadcastToDialogExcept(dialogID, c.UserID, ServerMessageTypeUserLeft, map[string]interface{}{
		"dialog_id": dialogID,
		"user_id":   c.UserID,
		"timestamp": time.Now(),
	})

	c.sendResponse("dialog_left", map[string]interface{}{
		"dialog_id": dialogID,
	}, message.ID)
}

// Send response to client
func (c *WebSocketConnection) sendResponse(messageType string, data map[string]interface{}, requestID string) {
	message := ServerMessage{
		Type:      messageType,
		Data:      data,
		ID:        requestID,
		Timestamp: time.Now(),
	}

	c.sendMessage(message)
}

// Send error to client
func (c *WebSocketConnection) sendError(code, message, details string) {
	errorMsg := ServerMessage{
		Type: ServerMessageTypeError,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}

	c.sendMessage(errorMsg)
}

// Send message to client
func (c *WebSocketConnection) sendMessage(message ServerMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	select {
	case c.Send <- messageBytes:
	default:
		// Channel is full, close connection
		close(c.Send)
	}
}

// Connection Manager methods
func (cm *ConnectionManager) addConnection(conn *WebSocketConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn.UserID != uuid.Nil {
		// Add to user sessions
		if cm.userSessions[conn.UserID] == nil {
			cm.userSessions[conn.UserID] = make(map[string]*WebSocketConnection)
		}
		cm.userSessions[conn.UserID][conn.ID] = conn
	}
}

func (cm *ConnectionManager) removeConnection(conn *WebSocketConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn.UserID != uuid.Nil {
		// Remove from user sessions
		if sessions, exists := cm.userSessions[conn.UserID]; exists {
			delete(sessions, conn.ID)
			if len(sessions) == 0 {
				delete(cm.userSessions, conn.UserID)
			}
		}

		// Remove from all dialog rooms
		for dialogID := range conn.Dialogs {
			if room, exists := cm.dialogRooms[dialogID]; exists {
				delete(room, conn.UserID)
				if len(room) == 0 {
					delete(cm.dialogRooms, dialogID)
				}
			}
		}
	}
}

func (cm *ConnectionManager) addToDialogRoom(dialogID, userID uuid.UUID, conn *WebSocketConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.dialogRooms[dialogID] == nil {
		cm.dialogRooms[dialogID] = make(map[uuid.UUID]*WebSocketConnection)
	}
	cm.dialogRooms[dialogID][userID] = conn
}

func (cm *ConnectionManager) removeFromDialogRoom(dialogID, userID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if room, exists := cm.dialogRooms[dialogID]; exists {
		delete(room, userID)
		if len(room) == 0 {
			delete(cm.dialogRooms, dialogID)
		}
	}
}

func (cm *ConnectionManager) getDialogConnections(dialogID uuid.UUID) []*WebSocketConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var connections []*WebSocketConnection
	if room, exists := cm.dialogRooms[dialogID]; exists {
		for _, conn := range room {
			connections = append(connections, conn)
		}
	}
	return connections
}

func (cm *ConnectionManager) getUserConnections(userID uuid.UUID) []*WebSocketConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var connections []*WebSocketConnection
	if sessions, exists := cm.userSessions[userID]; exists {
		for _, conn := range sessions {
			connections = append(connections, conn)
		}
	}
	return connections
}

// Broadcasting methods
func (h *WebSocketHandler) broadcastToDialog(dialogID uuid.UUID, messageType string, data map[string]interface{}) {
	connections := h.connections.getDialogConnections(dialogID)
	message := ServerMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, conn := range connections {
		select {
		case conn.Send <- messageBytes:
		default:
			// Connection is closed or buffer is full
			close(conn.Send)
			h.connections.removeConnection(conn)
		}
	}
}

func (h *WebSocketHandler) broadcastToDialogExcept(dialogID, excludeUserID uuid.UUID, messageType string, data map[string]interface{}) {
	connections := h.connections.getDialogConnections(dialogID)
	message := ServerMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, conn := range connections {
		if conn.UserID != excludeUserID {
			select {
			case conn.Send <- messageBytes:
			default:
				// Connection is closed or buffer is full
				close(conn.Send)
				h.connections.removeConnection(conn)
			}
		}
	}
}

func (h *WebSocketHandler) broadcastToUser(userID uuid.UUID, messageType string, data map[string]interface{}) {
	connections := h.connections.getUserConnections(userID)
	message := ServerMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, conn := range connections {
		select {
		case conn.Send <- messageBytes:
		default:
			// Connection is closed or buffer is full
			close(conn.Send)
			h.connections.removeConnection(conn)
		}
	}
}

func (h *WebSocketHandler) broadcastPresenceUpdate(userID uuid.UUID, presence interface{}) {
	// This would broadcast to users who can see this user's presence
	// Implementation depends on your privacy/contact system
}

// Helper methods
func (c *WebSocketConnection) validateToken(token string) (uuid.UUID, error) {
	// Implement token validation logic
	// This should integrate with your authentication service
	return uuid.New(), nil // Placeholder
}

func (c *WebSocketConnection) verifyDialogAccess(dialogID uuid.UUID) (bool, error) {
	// Implement dialog access verification
	// This should check if user is a participant in the dialog
	return true, nil // Placeholder
}

func (c *WebSocketConnection) subscribeToUserDialogs() {
	// Implement logic to automatically subscribe to user's active dialogs
}

func (c *WebSocketConnection) getPlatformFromUserAgent() string {
	// Parse user agent to determine platform
	return "web" // Placeholder
}

func (c *WebSocketConnection) convertMessageToMap(message interface{}) map[string]interface{} {
	// Convert message model to map for JSON serialization
	return map[string]interface{}{} // Placeholder
}

func generateConnectionID() string {
	return fmt.Sprintf("conn_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
}

// RegisterWebSocketRoutes registers WebSocket routes
func RegisterWebSocketRoutes(
	router *gin.RouterGroup,
	messageService *services.MessageService,
	presenceService *services.PresenceService,
) {
	handler := NewWebSocketHandler(messageService, presenceService)

	// WebSocket endpoint
	router.GET("/websocket", handler.HandleWebSocket)
}

// RegisterWebSocketRoutesWithMiddleware registers WebSocket routes with custom middleware
func RegisterWebSocketRoutesWithMiddleware(
	router *gin.RouterGroup,
	messageService *services.MessageService,
	presenceService *services.PresenceService,
	middlewares ...gin.HandlerFunc,
) {
	handler := NewWebSocketHandler(messageService, presenceService)

	// WebSocket endpoint with middleware
	ws := router.Group("")
	ws.Use(middlewares...)
	{
		ws.GET("/websocket", handler.HandleWebSocket)
	}
}