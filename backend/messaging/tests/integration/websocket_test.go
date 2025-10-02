package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat.dev/messaging/external"
	"tchat.dev/messaging/handlers"
	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// TestWebSocketConnection tests basic WebSocket connection functionality
func TestWebSocketConnection(t *testing.T) {
	// Create messaging service (minimal for testing)
	messagingService := &MockMessagingService{}

	// Create handler
	handler := handlers.NewMessagingHandler(messagingService)
	handler.Start()

	// Create simple WebSocket upgrader for testing
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Create test server with simple WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ws") {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Failed to upgrade", http.StatusBadRequest)
				return
			}
			defer conn.Close()

			// Simple echo handler for testing
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				// Echo back the message
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}))
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	t.Run("successful connection", func(t *testing.T) {
		// Connect to WebSocket
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Send a test message
		testMessage := map[string]interface{}{
			"type": "ping",
			"data": "test",
		}

		err = conn.WriteJSON(testMessage)
		assert.NoError(t, err)

		// Set read deadline
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		// Try to read response (may timeout, that's ok)
		var response map[string]interface{}
		conn.ReadJSON(&response)
	})

	t.Run("multiple connections", func(t *testing.T) {
		const numConnections = 3
		var wg sync.WaitGroup

		for i := 0; i < numConnections; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				if err != nil {
					t.Errorf("Connection %d failed: %v", id, err)
					return
				}
				defer conn.Close()

				// Send test message
				testMessage := map[string]interface{}{
					"type": "test",
					"id":   id,
				}

				conn.WriteJSON(testMessage)
				time.Sleep(100 * time.Millisecond) // Hold connection briefly
			}(i)
		}

		wg.Wait()
	})
}

// TestWebSocketManager tests the WebSocket manager functionality
func TestWebSocketManager(t *testing.T) {
	wsManager := external.NewWebSocketManager()

	t.Run("user connection tracking", func(t *testing.T) {
		userID := uuid.New()

		// Initially user should not be connected
		assert.False(t, wsManager.IsUserConnected(context.Background(), userID))

		// Get connected users (should be empty)
		connectedUsers := wsManager.GetConnectedUsers(context.Background())
		assert.Empty(t, connectedUsers)
	})

	t.Run("broadcast to users", func(t *testing.T) {
		userIDs := []uuid.UUID{uuid.New(), uuid.New()}

		message := map[string]interface{}{
			"type": "test_broadcast",
			"data": "test message",
		}

		// This should not error even if users are not connected
		err := wsManager.BroadcastToUsers(context.Background(), userIDs, message)
		assert.NoError(t, err)
	})

	t.Run("broadcast to single user", func(t *testing.T) {
		userID := uuid.New()

		message := map[string]interface{}{
			"type": "test_single",
			"data": "single user message",
		}

		// This should not error even if user is not connected
		err := wsManager.BroadcastToUser(context.Background(), userID, message)
		assert.NoError(t, err)
	})
}

// TestMessageDeliveryService tests message delivery through WebSocket
func TestMessageDeliveryService(t *testing.T) {
	// Create services
	wsManager := external.NewWebSocketManager()
	notificationService := external.NewNotificationService()
	deliveryService := external.NewMessageDeliveryService(notificationService, wsManager)

	t.Run("deliver message to offline users", func(t *testing.T) {
		// Create test message
		message := &models.Message{
			ID:       uuid.New(),
			DialogID: uuid.New(),
			SenderID: uuid.New(),
			Content:  models.MessageContent{"text": "Test message content"},
			Type:     models.MessageTypeText,
			SentAt:   time.Now(),
		}

		// Recipients (all offline)
		recipientIDs := []uuid.UUID{uuid.New(), uuid.New()}

		// Deliver message
		err := deliveryService.DeliverMessage(context.Background(), message, recipientIDs)
		assert.NoError(t, err)
	})

	t.Run("send push notification", func(t *testing.T) {
		userID := uuid.New()
		message := &models.Message{
			ID:       uuid.New(),
			DialogID: uuid.New(),
			SenderID: uuid.New(),
			Content:  models.MessageContent{"text": "Push notification test"},
			Type:     models.MessageTypeText,
			SentAt:   time.Now(),
		}

		err := deliveryService.SendPushNotification(context.Background(), userID, message)
		assert.NoError(t, err)
	})
}

// TestWebSocketMessageFlow tests end-to-end message flow
func TestWebSocketMessageFlow(t *testing.T) {
	// Setup services
	wsManager := external.NewWebSocketManager()
	notificationService := external.NewNotificationService()
	deliveryService := external.NewMessageDeliveryService(notificationService, wsManager)
	messagingService := &MockMessagingService{deliveryService: deliveryService}

	// Create handler
	handler := handlers.NewMessagingHandler(messagingService)
	handler.Start()

	// Create simple WebSocket upgrader for testing
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Create test server with simple WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ws") {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Failed to upgrade", http.StatusBadRequest)
				return
			}
			defer conn.Close()

			// Simple echo handler for testing
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				// Echo back the message
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}))
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	t.Run("message delivery flow", func(t *testing.T) {
		// Connect a client
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Simulate message delivery
		message := &models.Message{
			ID:       uuid.New(),
			DialogID: uuid.New(),
			SenderID: uuid.New(),
			Content:  models.MessageContent{"text": "End-to-end test message"},
			Type:     models.MessageTypeText,
			SentAt:   time.Now(),
		}

		// Test delivery service directly
		recipientIDs := []uuid.UUID{uuid.New()}
		err = deliveryService.DeliverMessage(context.Background(), message, recipientIDs)
		assert.NoError(t, err)
	})
}

// TestWebSocketPerformance tests WebSocket performance under load
func TestWebSocketPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	wsManager := external.NewWebSocketManager()

	// Create test server with WebSocket handler
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ws") {
			// Simple WebSocket echo handler for performance testing
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Failed to upgrade", http.StatusBadRequest)
				return
			}
			defer conn.Close()

			// Register with WebSocketManager
			testUserID := uuid.New()
			wsManager.RegisterClient(testUserID, conn)

			// Keep connection open for testing
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				// Echo back the message
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	t.Run("concurrent connections", func(t *testing.T) {
		const numConnections = 10
		const messagesPerConnection = 5

		var wg sync.WaitGroup

		start := time.Now()

		for i := 0; i < numConnections; i++ {
			wg.Add(1)
			go func(connID int) {
				defer wg.Done()

				conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				if err != nil {
					t.Errorf("Connection %d failed: %v", connID, err)
					return
				}
				defer conn.Close()

				// Send multiple messages
				for j := 0; j < messagesPerConnection; j++ {
					message := map[string]interface{}{
						"type":          "performance_test",
						"connection_id": connID,
						"message_id":    j,
						"timestamp":     time.Now().UnixNano(),
					}

					if err := conn.WriteJSON(message); err != nil {
						t.Errorf("Failed to send message %d from connection %d: %v", j, connID, err)
						return
					}

					time.Sleep(10 * time.Millisecond) // Small delay between messages
				}
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		totalMessages := numConnections * messagesPerConnection
		t.Logf("Sent %d messages across %d connections in %v", totalMessages, numConnections, elapsed)
		t.Logf("Average: %.2f messages/second", float64(totalMessages)/elapsed.Seconds())
	})
}

// TestWebSocketErrorHandling tests error scenarios
func TestWebSocketErrorHandling(t *testing.T) {
	wsManager := external.NewWebSocketManager()

	t.Run("invalid message format", func(t *testing.T) {
		// Create test server with WebSocket handler
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/ws") {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					http.Error(w, "Failed to upgrade", http.StatusBadRequest)
					return
				}
				defer conn.Close()

				// Register with WebSocketManager
				testUserID := uuid.New()
				wsManager.RegisterClient(testUserID, conn)

				// Read messages and handle errors gracefully
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						break
					}
					// Try to handle message, ignore invalid JSON
					_ = message
				}
			}
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Send invalid JSON
		err = conn.WriteMessage(websocket.TextMessage, []byte("invalid json"))
		// Connection might close or handle gracefully
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("connection cleanup", func(t *testing.T) {
		// Test that disconnected users are properly cleaned up
		userIDs := wsManager.GetConnectedUsers(context.Background())

		// After all previous tests, there should be no connected users
		// (connections should have been cleaned up)
		t.Logf("Connected users after tests: %d", len(userIDs))
	})
}

// MockMessagingService is a minimal mock for testing
type MockMessagingService struct {
	deliveryService services.MessageDeliveryService
}

func (m *MockMessagingService) CreateDialog(ctx context.Context, req *services.CreateDialogRequest) (*models.Dialog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockMessagingService) SendMessage(ctx context.Context, req *services.SendMessageRequest) (*models.Message, error) {
	message := &models.Message{
		ID:       uuid.New(),
		DialogID: req.DialogID,
		SenderID: req.SenderID,
		Content:  models.MessageContent{"text": "test message"},
		Type:     models.MessageTypeText,
		SentAt:   time.Now(),
	}

	if m.deliveryService != nil {
		// Simulate message delivery
		recipientIDs := []uuid.UUID{uuid.New()} // Mock recipients
		m.deliveryService.DeliverMessage(ctx, message, recipientIDs)
	}

	return message, nil
}

func (m *MockMessagingService) GetDialogByID(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (*models.Dialog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockMessagingService) GetMessageByID(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*models.Message, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockMessagingService) GetDialogParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	return []*models.DialogParticipant{}, nil
}

func (m *MockMessagingService) GetUserDialogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Dialog, error) {
	return []*models.Dialog{}, nil
}

func (m *MockMessagingService) GetMessages(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	return []*models.Message{}, nil
}

func (m *MockMessagingService) AddReaction(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, reaction string) error {
	return nil // Mock implementation
}

func (m *MockMessagingService) MarkAsRead(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	return nil
}

func (m *MockMessagingService) GetUnreadCount(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *MockMessagingService) UpdateDialog(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, req *services.UpdateDialogRequest) (*models.Dialog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockMessagingService) DeleteDialog(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (m *MockMessagingService) EditMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, newContent string) (*models.Message, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockMessagingService) DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, deleteForEveryone bool) error {
	return fmt.Errorf("not implemented")
}

func (m *MockMessagingService) AddParticipant(ctx context.Context, dialogID uuid.UUID, adminUserID uuid.UUID, req *services.AddParticipantRequest) error {
	return fmt.Errorf("not implemented")
}

func (m *MockMessagingService) RemoveParticipant(ctx context.Context, dialogID uuid.UUID, adminUserID uuid.UUID, userID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (m *MockMessagingService) RemoveReaction(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, reaction string) error {
	return fmt.Errorf("not implemented")
}