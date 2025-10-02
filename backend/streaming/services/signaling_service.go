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
	"github.com/pion/webrtc/v3"
	"tchat.dev/streaming/models"
)

// SignalingService handles WebSocket signaling for live streaming
type SignalingService struct {
	clients  map[uuid.UUID]*StreamClient
	streams  map[uuid.UUID]*StreamRoom
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

// StreamClient represents a connected viewer or broadcaster
type StreamClient struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Conn         *websocket.Conn
	Send         chan []byte
	StreamID     *uuid.UUID
	IsBroadcaster bool
	lastPing     time.Time
}

// StreamRoom represents a live stream session with broadcaster and viewers
type StreamRoom struct {
	ID           uuid.UUID
	Broadcaster  *StreamClient
	Viewers      map[uuid.UUID]*StreamClient
	mu           sync.RWMutex
}

// SignalingMessage represents the base structure for all signaling messages
type SignalingMessage struct {
	Type      string      `json:"type"`
	StreamID  *uuid.UUID  `json:"stream_id,omitempty"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// WebRTC signaling message types
type OfferMessage struct {
	Type     string                 `json:"type"`
	StreamID uuid.UUID              `json:"stream_id"`
	UserID   uuid.UUID              `json:"user_id"`
	Offer    map[string]interface{} `json:"offer"`
}

type AnswerMessage struct {
	Type     string                 `json:"type"`
	StreamID uuid.UUID              `json:"stream_id"`
	UserID   uuid.UUID              `json:"user_id"`
	Answer   map[string]interface{} `json:"answer"`
}

type ICECandidateMessage struct {
	Type      string                 `json:"type"`
	StreamID  uuid.UUID              `json:"stream_id"`
	UserID    uuid.UUID              `json:"user_id"`
	Candidate map[string]interface{} `json:"candidate"`
}

// Real-time interaction message types
type ChatMessageBroadcast struct {
	Type         string    `json:"type"`
	StreamID     uuid.UUID `json:"stream_id"`
	MessageID    uuid.UUID `json:"message_id"`
	SenderID     uuid.UUID `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	MessageText  string    `json:"message_text"`
	MessageType  string    `json:"message_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type ReactionBroadcast struct {
	Type         string    `json:"type"`
	StreamID     uuid.UUID `json:"stream_id"`
	ReactionID   uuid.UUID `json:"reaction_id"`
	ViewerID     uuid.UUID `json:"viewer_id"`
	ReactionType string    `json:"reaction_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type ViewerPresenceMessage struct {
	Type       string    `json:"type"`
	StreamID   uuid.UUID `json:"stream_id"`
	ViewerID   uuid.UUID `json:"viewer_id"`
	ViewerName string    `json:"viewer_name"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewSignalingService creates a new SignalingService instance
func NewSignalingService() *SignalingService {
	return &SignalingService{
		clients: make(map[uuid.UUID]*StreamClient),
		streams: make(map[uuid.UUID]*StreamRoom),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Configure based on security requirements
				// For production, implement proper CORS validation
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// HandleConnection handles WebSocket connection upgrade
func (s *SignalingService) HandleConnection(w http.ResponseWriter, r *http.Request) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("websocket upgrade failed: %w", err)
	}

	// Extract user ID from query params or JWT token
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		conn.Close()
		return fmt.Errorf("missing user_id parameter")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		conn.Close()
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	client := &StreamClient{
		ID:       uuid.New(),
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		lastPing: time.Now(),
	}

	s.mu.Lock()
	s.clients[client.ID] = client
	s.mu.Unlock()

	// Start goroutines for handling the client
	go s.writeToClient(client)
	go s.readFromClient(client)

	return nil
}

// SendOffer sends WebRTC offer to a viewer
func (s *SignalingService) SendOffer(conn *websocket.Conn, offer webrtc.SessionDescription) error {
	message := SignalingMessage{
		Type: "OFFER",
		Data: map[string]interface{}{
			"type": offer.Type.String(),
			"sdp":  offer.SDP,
		},
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal offer: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("failed to write offer: %w", err)
	}

	return nil
}

// SendAnswer sends WebRTC answer to broadcaster
func (s *SignalingService) SendAnswer(conn *websocket.Conn, answer webrtc.SessionDescription) error {
	message := SignalingMessage{
		Type: "ANSWER",
		Data: map[string]interface{}{
			"type": answer.Type.String(),
			"sdp":  answer.SDP,
		},
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal answer: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("failed to write answer: %w", err)
	}

	return nil
}

// SendICECandidate sends ICE candidate to peer
func (s *SignalingService) SendICECandidate(conn *websocket.Conn, candidate webrtc.ICECandidateInit) error {
	message := SignalingMessage{
		Type: "ICE_CANDIDATE",
		Data: map[string]interface{}{
			"candidate":        candidate.Candidate,
			"sdpMid":           candidate.SDPMid,
			"sdpMLineIndex":    candidate.SDPMLineIndex,
		},
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal ICE candidate: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("failed to write ICE candidate: %w", err)
	}

	return nil
}

// BroadcastChatMessage broadcasts a chat message to all viewers in a stream
func (s *SignalingService) BroadcastChatMessage(streamID uuid.UUID, message models.ChatMessage) error {
	broadcast := ChatMessageBroadcast{
		Type:        "CHAT",
		StreamID:    streamID,
		MessageID:   message.MessageID,
		SenderID:    message.SenderID,
		SenderName:  message.SenderDisplayName,
		MessageText: message.MessageText,
		MessageType: message.MessageType,
		Timestamp:   message.Timestamp,
	}

	messageBytes, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("failed to marshal chat message: %w", err)
	}

	s.broadcastToStream(streamID, messageBytes, nil)
	return nil
}

// BroadcastReaction broadcasts a reaction to all viewers in a stream
func (s *SignalingService) BroadcastReaction(streamID uuid.UUID, reaction models.StreamReaction) error {
	broadcast := ReactionBroadcast{
		Type:         "REACTION",
		StreamID:     streamID,
		ReactionID:   reaction.ReactionID,
		ViewerID:     reaction.ViewerID,
		ReactionType: reaction.ReactionType,
		Timestamp:    reaction.Timestamp,
	}

	messageBytes, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("failed to marshal reaction: %w", err)
	}

	s.broadcastToStream(streamID, messageBytes, nil)
	return nil
}

// BroadcastViewerJoin broadcasts when a viewer joins a stream
func (s *SignalingService) BroadcastViewerJoin(streamID uuid.UUID, viewerID uuid.UUID) error {
	presence := ViewerPresenceMessage{
		Type:      "VIEWER_JOIN",
		StreamID:  streamID,
		ViewerID:  viewerID,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(presence)
	if err != nil {
		return fmt.Errorf("failed to marshal viewer join: %w", err)
	}

	s.broadcastToStream(streamID, messageBytes, &viewerID)
	return nil
}

// BroadcastViewerLeave broadcasts when a viewer leaves a stream
func (s *SignalingService) BroadcastViewerLeave(streamID uuid.UUID, viewerID uuid.UUID) error {
	presence := ViewerPresenceMessage{
		Type:      "VIEWER_LEAVE",
		StreamID:  streamID,
		ViewerID:  viewerID,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(presence)
	if err != nil {
		return fmt.Errorf("failed to marshal viewer leave: %w", err)
	}

	s.broadcastToStream(streamID, messageBytes, &viewerID)
	return nil
}

// Private methods

// readFromClient handles incoming messages from a WebSocket client
func (s *SignalingService) readFromClient(client *StreamClient) {
	defer func() {
		s.disconnect(client)
	}()

	// Set initial read deadline
	if err := client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
	}

	// Configure pong handler for keepalive
	client.Conn.SetPongHandler(func(string) error {
		client.lastPing = time.Now()
		return client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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
func (s *SignalingService) writeToClient(client *StreamClient) {
	ticker := time.NewTicker(30 * time.Second) // 30s ping interval
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
			}

			if !ok {
				// Channel closed
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping message for keepalive
			if err := client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
			}

			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from clients
func (s *SignalingService) handleMessage(client *StreamClient, messageBytes []byte) error {
	var baseMessage SignalingMessage
	if err := json.Unmarshal(messageBytes, &baseMessage); err != nil {
		return fmt.Errorf("invalid message format: %w", err)
	}

	baseMessage.Timestamp = time.Now()

	switch baseMessage.Type {
	case "join-stream":
		return s.handleJoinStream(client, messageBytes)
	case "leave-stream":
		return s.handleLeaveStream(client)
	case "OFFER":
		return s.handleOffer(client, messageBytes)
	case "ANSWER":
		return s.handleAnswer(client, messageBytes)
	case "ICE_CANDIDATE":
		return s.handleICECandidate(client, messageBytes)
	case "heartbeat":
		return s.handleHeartbeat(client)
	default:
		return fmt.Errorf("unknown message type: %s", baseMessage.Type)
	}
}

// handleJoinStream handles stream join requests
func (s *SignalingService) handleJoinStream(client *StreamClient, messageBytes []byte) error {
	var joinMessage struct {
		Type          string    `json:"type"`
		StreamID      uuid.UUID `json:"stream_id"`
		IsBroadcaster bool      `json:"is_broadcaster"`
	}

	if err := json.Unmarshal(messageBytes, &joinMessage); err != nil {
		return err
	}

	client.StreamID = &joinMessage.StreamID
	client.IsBroadcaster = joinMessage.IsBroadcaster

	// Get or create stream room
	room := s.getOrCreateStreamRoom(joinMessage.StreamID)

	room.mu.Lock()
	if joinMessage.IsBroadcaster {
		room.Broadcaster = client
	} else {
		room.Viewers[client.ID] = client
	}
	room.mu.Unlock()

	// Broadcast viewer join if not broadcaster
	if !joinMessage.IsBroadcaster {
		s.BroadcastViewerJoin(joinMessage.StreamID, client.UserID)
	}

	return nil
}

// handleLeaveStream handles stream leave requests
func (s *SignalingService) handleLeaveStream(client *StreamClient) error {
	if client.StreamID == nil {
		return fmt.Errorf("client not in any stream")
	}

	streamID := *client.StreamID
	s.removeClientFromStream(client, streamID)

	// Broadcast viewer leave if not broadcaster
	if !client.IsBroadcaster {
		s.BroadcastViewerLeave(streamID, client.UserID)
	}

	return nil
}

// handleOffer handles WebRTC offer messages
func (s *SignalingService) handleOffer(client *StreamClient, messageBytes []byte) error {
	var offerMessage OfferMessage
	if err := json.Unmarshal(messageBytes, &offerMessage); err != nil {
		return err
	}

	// In a live streaming context, offers typically come from the broadcaster
	// and should be forwarded to specific viewers or all viewers
	return nil
}

// handleAnswer handles WebRTC answer messages
func (s *SignalingService) handleAnswer(client *StreamClient, messageBytes []byte) error {
	var answerMessage AnswerMessage
	if err := json.Unmarshal(messageBytes, &answerMessage); err != nil {
		return err
	}

	// In a live streaming context, answers come from viewers
	// and should be forwarded to the broadcaster
	return nil
}

// handleICECandidate handles ICE candidate exchange
func (s *SignalingService) handleICECandidate(client *StreamClient, messageBytes []byte) error {
	var iceMessage ICECandidateMessage
	if err := json.Unmarshal(messageBytes, &iceMessage); err != nil {
		return err
	}

	// Forward ICE candidates between broadcaster and viewers
	return nil
}

// handleHeartbeat handles heartbeat messages
func (s *SignalingService) handleHeartbeat(client *StreamClient) error {
	client.lastPing = time.Now()
	return nil
}

// Helper methods

func (s *SignalingService) getOrCreateStreamRoom(streamID uuid.UUID) *StreamRoom {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.streams[streamID]; exists {
		return room
	}

	room := &StreamRoom{
		ID:      streamID,
		Viewers: make(map[uuid.UUID]*StreamClient),
	}
	s.streams[streamID] = room
	return room
}

func (s *SignalingService) removeClientFromStream(client *StreamClient, streamID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.streams[streamID]; exists {
		room.mu.Lock()
		if client.IsBroadcaster {
			room.Broadcaster = nil
		} else {
			delete(room.Viewers, client.ID)
		}
		isEmpty := room.Broadcaster == nil && len(room.Viewers) == 0
		room.mu.Unlock()

		if isEmpty {
			delete(s.streams, streamID)
		}
	}

	client.StreamID = nil
}

func (s *SignalingService) broadcastToStream(streamID uuid.UUID, messageBytes []byte, excludeClient *uuid.UUID) {
	s.mu.RLock()
	room, exists := s.streams[streamID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	// Send to broadcaster
	if room.Broadcaster != nil && (excludeClient == nil || room.Broadcaster.UserID != *excludeClient) {
		select {
		case room.Broadcaster.Send <- messageBytes:
		default:
			log.Printf("Failed to send message to broadcaster")
		}
	}

	// Send to all viewers
	for _, viewer := range room.Viewers {
		if excludeClient != nil && viewer.UserID == *excludeClient {
			continue
		}
		select {
		case viewer.Send <- messageBytes:
		default:
			log.Printf("Failed to send message to viewer %s", viewer.UserID)
		}
	}
}

func (s *SignalingService) sendError(client *StreamClient, code, message string) {
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
		close(client.Send)
	}
}

func (s *SignalingService) disconnect(client *StreamClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove from stream room
	if client.StreamID != nil {
		s.removeClientFromStream(client, *client.StreamID)
	}

	// Remove from clients map
	delete(s.clients, client.ID)

	// Close the send channel
	close(client.Send)
}