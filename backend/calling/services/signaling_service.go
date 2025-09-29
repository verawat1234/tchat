package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"tchat.dev/calling/models"
)

// SignalingService handles WebRTC signaling coordination
type SignalingService struct {
	clients         map[uuid.UUID]*Client
	rooms           map[uuid.UUID]*Room
	callService     *CallService
	presenceService *PresenceService
	mu              sync.RWMutex
	upgrader        websocket.Upgrader
}

// Client represents a connected WebSocket client
type Client struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Send     chan []byte
	RoomID   *uuid.UUID
	IsVideo  bool
	lastPing time.Time
}

// Room represents a call room with participants
type Room struct {
	ID          uuid.UUID
	Clients     map[uuid.UUID]*Client
	CallSession *models.CallSession
	mu          sync.RWMutex
}

// SignalingMessage represents the base structure for all signaling messages
type SignalingMessage struct {
	Type      string      `json:"type"`
	From      *uuid.UUID  `json:"from,omitempty"`
	To        *uuid.UUID  `json:"to,omitempty"`
	RoomID    *uuid.UUID  `json:"room_id,omitempty"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	IsVideo   *bool       `json:"is_video,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// WebRTC-specific message types
type WebRTCOfferMessage struct {
	Type   string                 `json:"type"`
	From   uuid.UUID              `json:"from"`
	To     uuid.UUID              `json:"to"`
	RoomID uuid.UUID              `json:"room_id"`
	Offer  map[string]interface{} `json:"offer"`
}

type WebRTCAnswerMessage struct {
	Type   string                 `json:"type"`
	From   uuid.UUID              `json:"from"`
	To     uuid.UUID              `json:"to"`
	RoomID uuid.UUID              `json:"room_id"`
	Answer map[string]interface{} `json:"answer"`
}

type ICECandidateMessage struct {
	Type      string                 `json:"type"`
	From      uuid.UUID              `json:"from"`
	To        uuid.UUID              `json:"to"`
	RoomID    uuid.UUID              `json:"room_id"`
	Candidate map[string]interface{} `json:"candidate"`
}

type MediaToggleMessage struct {
	Type      string    `json:"type"`
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	MediaType string    `json:"media_type"` // "audio" or "video"
	Enabled   bool      `json:"enabled"`
}

// NewSignalingService creates a new SignalingService instance
func NewSignalingService(callService *CallService, presenceService *PresenceService) *SignalingService {
	return &SignalingService{
		clients:         make(map[uuid.UUID]*Client),
		rooms:           make(map[uuid.UUID]*Room),
		callService:     callService,
		presenceService: presenceService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Configure based on your security requirements
				return true // For development - restrict in production
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// HandleWebSocketConnection handles new WebSocket connections
func (s *SignalingService) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("websocket upgrade failed: %w", err)
	}

	client := &Client{
		ID:       uuid.New(),
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		lastPing: time.Now(),
	}

	s.mu.Lock()
	s.clients[client.ID] = client
	s.mu.Unlock()

	// Set user as online
	if err := s.presenceService.SetUserOnline(userID); err != nil {
		// Log error but continue
		_ = err // Avoid unused variable warning
	}

	// Start goroutines for handling the client
	go s.writeToClient(client)
	go s.readFromClient(client)

	return nil
}

// readFromClient handles messages from a WebSocket client
func (s *SignalingService) readFromClient(client *Client) {
	defer func() {
		s.disconnect(client)
	}()

	if err := client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		// Log error but continue
		_ = err // Avoid unused variable warning
	}
	client.Conn.SetPongHandler(func(string) error {
		client.lastPing = time.Now()
		if err := client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			// Log error but continue
			_ = err // Avoid unused variable warning
		}
		return nil
	})

	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if err := s.handleMessage(client, messageBytes); err != nil {
			log.Printf("Error handling message: %v", err)
			s.sendError(client, "MESSAGE_PROCESSING_ERROR", err.Error())
		}
	}
}

// writeToClient handles sending messages to a WebSocket client
func (s *SignalingService) writeToClient(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		if err := client.Conn.Close(); err != nil {
			// Log error but continue
			_ = err // Avoid unused variable warning
		}
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				// Log error but continue
				_ = err // Avoid unused variable warning
			}
			if !ok {
				if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					// Log error but continue
					_ = err // Avoid unused variable warning
				}
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				// Log error but continue
				_ = err // Avoid unused variable warning
			}
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from clients
func (s *SignalingService) handleMessage(client *Client, messageBytes []byte) error {
	var baseMessage SignalingMessage
	if err := json.Unmarshal(messageBytes, &baseMessage); err != nil {
		return fmt.Errorf("invalid message format: %w", err)
	}

	baseMessage.Timestamp = time.Now()

	switch baseMessage.Type {
	case "join-room":
		return s.handleJoinRoom(client, messageBytes)
	case "leave-room":
		return s.handleLeaveRoom(client, baseMessage)
	case "webrtc-offer":
		return s.handleWebRTCOffer(client, messageBytes)
	case "webrtc-answer":
		return s.handleWebRTCAnswer(client, messageBytes)
	case "ice-candidate":
		return s.handleICECandidate(client, messageBytes)
	case "media-toggle":
		return s.handleMediaToggle(client, messageBytes)
	case "heartbeat":
		return s.handleHeartbeat(client, baseMessage)
	default:
		return fmt.Errorf("unknown message type: %s", baseMessage.Type)
	}
}

// handleJoinRoom handles room join requests
func (s *SignalingService) handleJoinRoom(client *Client, messageBytes []byte) error {
	var joinMessage struct {
		Type    string    `json:"type"`
		RoomID  uuid.UUID `json:"room_id"`
		UserID  uuid.UUID `json:"user_id"`
		IsVideo bool      `json:"is_video"`
	}

	if err := json.Unmarshal(messageBytes, &joinMessage); err != nil {
		return err
	}

	// Validate user matches client
	if joinMessage.UserID != client.UserID {
		return fmt.Errorf("user ID mismatch")
	}

	// Get or create room
	room := s.getOrCreateRoom(joinMessage.RoomID)

	// Add client to room
	room.mu.Lock()
	room.Clients[client.ID] = client
	room.mu.Unlock()

	client.RoomID = &joinMessage.RoomID
	client.IsVideo = joinMessage.IsVideo

	// Notify other clients in the room
	s.broadcastToRoom(joinMessage.RoomID, SignalingMessage{
		Type:   "user-joined",
		RoomID: &joinMessage.RoomID,
		UserID: &client.UserID,
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":   client.UserID,
				"name": "", // Would be populated from user service
			},
			"is_video": joinMessage.IsVideo,
		},
		Timestamp: time.Now(),
	}, &client.ID)

	return nil
}

// handleLeaveRoom handles room leave requests
func (s *SignalingService) handleLeaveRoom(client *Client, message SignalingMessage) error {
	if client.RoomID == nil {
		return fmt.Errorf("client not in any room")
	}

	roomID := *client.RoomID
	s.removeClientFromRoom(client, roomID)

	// Notify other clients
	s.broadcastToRoom(roomID, SignalingMessage{
		Type:   "user-left",
		RoomID: &roomID,
		UserID: &client.UserID,
		Data: map[string]interface{}{
			"reason": "normal",
		},
		Timestamp: time.Now(),
	}, &client.ID)

	return nil
}

// handleWebRTCOffer handles WebRTC offer messages
func (s *SignalingService) handleWebRTCOffer(client *Client, messageBytes []byte) error {
	var offerMessage WebRTCOfferMessage
	if err := json.Unmarshal(messageBytes, &offerMessage); err != nil {
		return err
	}

	// Forward offer to target client
	return s.sendToUser(offerMessage.To, SignalingMessage{
		Type:      "webrtc-offer",
		From:      &offerMessage.From,
		To:        &offerMessage.To,
		RoomID:    &offerMessage.RoomID,
		Data:      offerMessage.Offer,
		Timestamp: time.Now(),
	})
}

// handleWebRTCAnswer handles WebRTC answer messages
func (s *SignalingService) handleWebRTCAnswer(client *Client, messageBytes []byte) error {
	var answerMessage WebRTCAnswerMessage
	if err := json.Unmarshal(messageBytes, &answerMessage); err != nil {
		return err
	}

	// Forward answer to target client
	return s.sendToUser(answerMessage.To, SignalingMessage{
		Type:      "webrtc-answer",
		From:      &answerMessage.From,
		To:        &answerMessage.To,
		RoomID:    &answerMessage.RoomID,
		Data:      answerMessage.Answer,
		Timestamp: time.Now(),
	})
}

// handleICECandidate handles ICE candidate messages
func (s *SignalingService) handleICECandidate(client *Client, messageBytes []byte) error {
	var iceMessage ICECandidateMessage
	if err := json.Unmarshal(messageBytes, &iceMessage); err != nil {
		return err
	}

	// Forward ICE candidate to target client
	return s.sendToUser(iceMessage.To, SignalingMessage{
		Type:      "ice-candidate",
		From:      &iceMessage.From,
		To:        &iceMessage.To,
		RoomID:    &iceMessage.RoomID,
		Data:      iceMessage.Candidate,
		Timestamp: time.Now(),
	})
}

// handleMediaToggle handles media toggle messages
func (s *SignalingService) handleMediaToggle(client *Client, messageBytes []byte) error {
	var toggleMessage MediaToggleMessage
	if err := json.Unmarshal(messageBytes, &toggleMessage); err != nil {
		return err
	}

	// Update call service with media state
	if client.RoomID != nil {
		err := s.callService.ToggleMedia(*client.RoomID, client.UserID, toggleMessage.MediaType, toggleMessage.Enabled)
		if err != nil {
			log.Printf("Failed to update media state: %v", err)
		}
	}

	// Broadcast to room
	if client.RoomID != nil {
		s.broadcastToRoom(*client.RoomID, SignalingMessage{
			Type:   "media-toggle",
			RoomID: client.RoomID,
			UserID: &client.UserID,
			Data: map[string]interface{}{
				"media_type": toggleMessage.MediaType,
				"enabled":    toggleMessage.Enabled,
			},
			Timestamp: time.Now(),
		}, &client.ID)
	}

	return nil
}

// handleHeartbeat handles heartbeat messages
func (s *SignalingService) handleHeartbeat(client *Client, message SignalingMessage) error {
	client.lastPing = time.Now()
	return nil
}

// Helper methods

func (s *SignalingService) getOrCreateRoom(roomID uuid.UUID) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.rooms[roomID]; exists {
		return room
	}

	room := &Room{
		ID:      roomID,
		Clients: make(map[uuid.UUID]*Client),
	}
	s.rooms[roomID] = room
	return room
}

func (s *SignalingService) removeClientFromRoom(client *Client, roomID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.rooms[roomID]; exists {
		room.mu.Lock()
		delete(room.Clients, client.ID)
		isEmpty := len(room.Clients) == 0
		room.mu.Unlock()

		if isEmpty {
			delete(s.rooms, roomID)
		}
	}

	client.RoomID = nil
}

func (s *SignalingService) broadcastToRoom(roomID uuid.UUID, message SignalingMessage, excludeClient *uuid.UUID) {
	s.mu.RLock()
	room, exists := s.rooms[roomID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	messageBytes, _ := json.Marshal(message)

	room.mu.RLock()
	for clientID, client := range room.Clients {
		if excludeClient != nil && clientID == *excludeClient {
			continue
		}
		select {
		case client.Send <- messageBytes:
		default:
			close(client.Send)
			delete(room.Clients, clientID)
		}
	}
	room.mu.RUnlock()
}

func (s *SignalingService) sendToUser(userID uuid.UUID, message SignalingMessage) error {
	s.mu.RLock()
	var targetClient *Client
	for _, client := range s.clients {
		if client.UserID == userID {
			targetClient = client
			break
		}
	}
	s.mu.RUnlock()

	if targetClient == nil {
		return fmt.Errorf("user %s not connected", userID)
	}

	messageBytes, _ := json.Marshal(message)
	select {
	case targetClient.Send <- messageBytes:
		return nil
	default:
		return fmt.Errorf("failed to send message to user %s", userID)
	}
}

func (s *SignalingService) sendError(client *Client, code, message string) {
	errorMessage := SignalingMessage{
		Type: "error",
		Data: map[string]interface{}{
			"code":    code,
			"message": message,
		},
		Timestamp: time.Now(),
	}

	messageBytes, _ := json.Marshal(errorMessage)
	select {
	case client.Send <- messageBytes:
	default:
		// Client channel is full, close it
		close(client.Send)
	}
}

func (s *SignalingService) disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove from rooms
	if client.RoomID != nil {
		s.removeClientFromRoom(client, *client.RoomID)
	}

	// Remove from clients map
	delete(s.clients, client.ID)

	// Set user offline
	if err := s.presenceService.SetUserOffline(client.UserID); err != nil {
		// Log error but continue
		_ = err // Avoid unused variable warning
	}

	// Close the send channel
	close(client.Send)
}

// SendCallInvitation sends a call invitation to a user
func (s *SignalingService) SendCallInvitation(callSession *models.CallSession, calleeID uuid.UUID) error {
	caller := callSession.GetParticipantByUserID(callSession.InitiatedBy)
	if caller == nil {
		return fmt.Errorf("caller not found in call session")
	}

	invitation := SignalingMessage{
		Type: "call-invitation",
		Data: map[string]interface{}{
			"call_id": callSession.ID,
			"caller": map[string]interface{}{
				"id":   callSession.InitiatedBy,
				"name": "", // Would be populated from user service
			},
			"call_type":  string(callSession.Type),
			"expires_at": time.Now().Add(30 * time.Second),
		},
		Timestamp: time.Now(),
	}

	return s.sendToUser(calleeID, invitation)
}

// NotifyCallAnswered notifies participants that a call was answered
func (s *SignalingService) NotifyCallAnswered(callSession *models.CallSession, answeredBy uuid.UUID) {
	message := SignalingMessage{
		Type: "call-answered",
		Data: map[string]interface{}{
			"call_id":     callSession.ID,
			"answered_by": answeredBy,
		},
		Timestamp: time.Now(),
	}

	// Notify all participants
	for _, participant := range callSession.Participants {
		if participant.UserID != answeredBy {
			if err := s.sendToUser(participant.UserID, message); err != nil {
				// Log error but continue
				_ = err // Avoid unused variable warning
			}
		}
	}
}

// NotifyCallEnded notifies participants that a call has ended
func (s *SignalingService) NotifyCallEnded(callSession *models.CallSession, endedBy uuid.UUID) {
	message := SignalingMessage{
		Type: "call-ended",
		Data: map[string]interface{}{
			"call_id":  callSession.ID,
			"ended_by": endedBy,
			"duration": callSession.Duration,
			"reason":   "normal",
		},
		Timestamp: time.Now(),
	}

	// Notify all participants and close any active rooms
	for _, participant := range callSession.Participants {
		if err := s.sendToUser(participant.UserID, message); err != nil {
			// Log error but continue
			_ = err // Avoid unused variable warning
		}
	}

	// Clean up room if it exists
	s.mu.Lock()
	if room, exists := s.rooms[callSession.ID]; exists {
		room.mu.Lock()
		// Close all client connections in the room
		for _, client := range room.Clients {
			close(client.Send)
		}
		room.mu.Unlock()
		delete(s.rooms, callSession.ID)
	}
	s.mu.Unlock()
}
