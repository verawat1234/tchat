package external

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"tchat.dev/messaging/services"
)

// WebSocketConnection represents a websocket connection
type WebSocketConnection struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
}

// WebSocketManager implements services.WebSocketManager for real-time messaging
type WebSocketManager struct {
	clients    map[uuid.UUID]*WebSocketConnection
	clientsMux sync.RWMutex
	register   chan *WebSocketConnection
	unregister chan *WebSocketConnection
	broadcast  chan []byte
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() services.WebSocketManager {
	manager := &WebSocketManager{
		clients:    make(map[uuid.UUID]*WebSocketConnection),
		register:   make(chan *WebSocketConnection),
		unregister: make(chan *WebSocketConnection),
		broadcast:  make(chan []byte),
	}

	// Start the hub goroutine
	go manager.run()

	return manager
}

// run handles the main hub loop for managing WebSocket connections
func (m *WebSocketManager) run() {
	for {
		select {
		case client := <-m.register:
			m.clientsMux.Lock()
			m.clients[client.UserID] = client
			m.clientsMux.Unlock()
			log.Printf("User %s connected via WebSocket", client.UserID)

		case client := <-m.unregister:
			m.clientsMux.Lock()
			if _, ok := m.clients[client.UserID]; ok {
				delete(m.clients, client.UserID)
				close(client.Send)
			}
			m.clientsMux.Unlock()
			log.Printf("User %s disconnected from WebSocket", client.UserID)

		case message := <-m.broadcast:
			m.clientsMux.RLock()
			for userID, client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, userID)
				}
			}
			m.clientsMux.RUnlock()
		}
	}
}

// RegisterClient registers a new WebSocket client
func (m *WebSocketManager) RegisterClient(userID uuid.UUID, conn interface{}) {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		log.Printf("Invalid connection type for user %s", userID)
		return
	}

	client := &WebSocketConnection{
		UserID: userID,
		Conn:   wsConn,
		Send:   make(chan []byte, 256),
	}

	m.register <- client

	// Start goroutines for this connection
	go m.writePump(client)
	go m.readPump(client)
}

// writePump pumps messages from the hub to the websocket connection
func (m *WebSocketManager) writePump(client *WebSocketConnection) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to write to WebSocket: %v", err)
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (m *WebSocketManager) readPump(client *WebSocketConnection) {
	defer func() {
		m.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received message from user %s: %s", client.UserID, string(message))
		// Handle incoming messages here if needed
	}
}

// BroadcastToUser sends a message to a specific user
func (m *WebSocketManager) BroadcastToUser(ctx context.Context, userID uuid.UUID, message interface{}) error {
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m.clientsMux.RLock()
	client, exists := m.clients[userID]
	m.clientsMux.RUnlock()

	if !exists {
		log.Printf("User %s is not connected via WebSocket", userID)
		return nil
	}

	select {
	case client.Send <- messageData:
		log.Printf("Message sent to user %s via WebSocket", userID)
	default:
		log.Printf("Failed to send message to user %s: channel full", userID)
	}

	return nil
}

// BroadcastToUsers sends a message to multiple users
func (m *WebSocketManager) BroadcastToUsers(ctx context.Context, userIDs []uuid.UUID, message interface{}) error {
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	successCount := 0
	m.clientsMux.RLock()
	for _, userID := range userIDs {
		if client, exists := m.clients[userID]; exists {
			select {
			case client.Send <- messageData:
				successCount++
			default:
				log.Printf("Failed to send message to user %s: channel full", userID)
			}
		}
	}
	m.clientsMux.RUnlock()

	log.Printf("Message broadcast to %d/%d users via WebSocket", successCount, len(userIDs))
	return nil
}

// GetConnectedUsers returns a list of currently connected user IDs
func (m *WebSocketManager) GetConnectedUsers(ctx context.Context) []uuid.UUID {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	userIDs := make([]uuid.UUID, 0, len(m.clients))
	for userID := range m.clients {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

// IsUserConnected checks if a user is currently connected via WebSocket
func (m *WebSocketManager) IsUserConnected(ctx context.Context, userID uuid.UUID) bool {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	_, exists := m.clients[userID]
	return exists
}