package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	sharedModels "tchat.dev/shared/models"

	"tchat.dev/messaging/services"
)

// WebSocketHandler handles WebSocket connections using external.WebSocketManager
type WebSocketHandler struct {
	wsManager services.WebSocketManager
	upgrader  websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(wsManager services.WebSocketManager) *WebSocketHandler {
	return &WebSocketHandler{
		wsManager: wsManager,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// HandleWebSocket upgrades HTTP connection and registers with external.WebSocketManager
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract user from authentication context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	// Register client with external.WebSocketManager
	// This starts read/write pumps automatically
	h.wsManager.RegisterClient(user.ID, conn)

	log.Printf("WebSocket connection established for user: %s", user.ID)
}

// getUserFromContext extracts user from request context
func (h *WebSocketHandler) getUserFromContext(ctx interface{ Value(key interface{}) interface{} }) *sharedModels.User {
	userValue := ctx.Value("user")
	if userValue == nil {
		return nil
	}

	user, ok := userValue.(*sharedModels.User)
	if !ok {
		return nil
	}

	return user
}

// RegisterRoutes registers WebSocket routes
func (h *WebSocketHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/ws", h.HandleWebSocket).Methods("GET")
	log.Println("WebSocket route registered at /ws")
}
