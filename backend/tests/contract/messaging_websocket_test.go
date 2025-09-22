package contract

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T022: Contract test WebSocket /websocket - Real-time messaging WebSocket
func TestMessagingWebSocketContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// WebSocket upgrader for testing
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for testing
		},
	}

	tests := []struct {
		name        string
		authHeader  string
		setupFunc   func(*gin.Engine)
		testFunc    func(*testing.T, *websocket.Conn)
		expectError bool
		description string
	}{
		{
			name:       "successful_websocket_connection",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
						return
					}
					defer conn.Close()

					// Send welcome message
					welcomeMsg := map[string]interface{}{
						"type": "connection_established",
						"data": map[string]interface{}{
							"user_id":     "user_123e4567-e89b-12d3-a456-426614174000",
							"session_id":  "session_abc123",
							"server_time": "2023-12-01T15:30:00Z",
						},
						"timestamp": "2023-12-01T15:30:00Z",
					}
					conn.WriteJSON(welcomeMsg)

					// Listen for messages
					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						// Echo back with confirmation
						response := map[string]interface{}{
							"type": "message_received",
							"data": msg,
							"timestamp": "2023-12-01T15:30:00Z",
						}
						conn.WriteJSON(response)
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Read welcome message
				var welcomeMsg map[string]interface{}
				err := conn.ReadJSON(&welcomeMsg)
				require.NoError(t, err)

				assert.Equal(t, "connection_established", welcomeMsg["type"])
				assert.Contains(t, welcomeMsg, "data")
				assert.Contains(t, welcomeMsg, "timestamp")

				data := welcomeMsg["data"].(map[string]interface{})
				assert.Contains(t, data, "user_id")
				assert.Contains(t, data, "session_id")
			},
			expectError: false,
			description: "Should establish WebSocket connection successfully",
		},
		{
			name:       "join_dialog_message",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						msgType := msg["type"].(string)
						if msgType == "join_dialog" {
							data := msg["data"].(map[string]interface{})
							dialogID := data["dialog_id"].(string)

							response := map[string]interface{}{
								"type": "dialog_joined",
								"data": map[string]interface{}{
									"dialog_id": dialogID,
									"user_id":   "user_123",
									"status":    "success",
								},
								"timestamp": "2023-12-01T15:30:00Z",
							}
							conn.WriteJSON(response)
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Send join dialog message
				joinMsg := map[string]interface{}{
					"type": "join_dialog",
					"data": map[string]interface{}{
						"dialog_id": "dialog_123e4567-e89b-12d3-a456-426614174000",
					},
					"timestamp": "2023-12-01T15:30:00Z",
				}
				err := conn.WriteJSON(joinMsg)
				require.NoError(t, err)

				// Read response
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "dialog_joined", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "dialog_123e4567-e89b-12d3-a456-426614174000", data["dialog_id"])
				assert.Equal(t, "success", data["status"])
			},
			expectError: false,
			description: "Should handle join dialog message correctly",
		},
		{
			name:       "send_message_via_websocket",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						msgType := msg["type"].(string)
						if msgType == "send_message" {
							data := msg["data"].(map[string]interface{})

							// Simulate message creation
							newMessage := map[string]interface{}{
								"type": "new_message",
								"data": map[string]interface{}{
									"message": map[string]interface{}{
										"id":           "msg_new_123",
										"dialog_id":    data["dialog_id"],
										"sender_id":    "user_123",
										"content":      data["content"],
										"message_type": data["message_type"],
										"created_at":   "2023-12-01T15:30:00Z",
									},
								},
								"timestamp": "2023-12-01T15:30:00Z",
							}
							conn.WriteJSON(newMessage)
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Send message
				sendMsg := map[string]interface{}{
					"type": "send_message",
					"data": map[string]interface{}{
						"dialog_id":    "dialog_123",
						"content":      "Hello WebSocket!",
						"message_type": "text",
					},
					"timestamp": "2023-12-01T15:30:00Z",
				}
				err := conn.WriteJSON(sendMsg)
				require.NoError(t, err)

				// Read new message event
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "new_message", response["type"])
				data := response["data"].(map[string]interface{})
				message := data["message"].(map[string]interface{})

				assert.Equal(t, "dialog_123", message["dialog_id"])
				assert.Equal(t, "Hello WebSocket!", message["content"])
				assert.Equal(t, "text", message["message_type"])
				assert.Contains(t, message, "id")
				assert.Contains(t, message, "created_at")
			},
			expectError: false,
			description: "Should handle send message via WebSocket correctly",
		},
		{
			name:       "typing_indicator",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						msgType := msg["type"].(string)
						if msgType == "typing" {
							data := msg["data"].(map[string]interface{})

							// Broadcast typing indicator to other participants
							typingEvent := map[string]interface{}{
								"type": "user_typing",
								"data": map[string]interface{}{
									"dialog_id": data["dialog_id"],
									"user_id":   "user_123",
									"typing":    data["typing"],
								},
								"timestamp": "2023-12-01T15:30:00Z",
							}
							conn.WriteJSON(typingEvent)
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Send typing indicator
				typingMsg := map[string]interface{}{
					"type": "typing",
					"data": map[string]interface{}{
						"dialog_id": "dialog_123",
						"typing":    true,
					},
					"timestamp": "2023-12-01T15:30:00Z",
				}
				err := conn.WriteJSON(typingMsg)
				require.NoError(t, err)

				// Read typing event
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "user_typing", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "dialog_123", data["dialog_id"])
				assert.Equal(t, "user_123", data["user_id"])
				assert.Equal(t, true, data["typing"])
			},
			expectError: false,
			description: "Should handle typing indicator correctly",
		},
		{
			name:       "user_presence_update",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					// Send initial presence
					presenceEvent := map[string]interface{}{
						"type": "user_presence",
						"data": map[string]interface{}{
							"user_id":   "user_123",
							"status":    "online",
							"last_seen": "2023-12-01T15:30:00Z",
						},
						"timestamp": "2023-12-01T15:30:00Z",
					}
					conn.WriteJSON(presenceEvent)

					// Keep connection alive
					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Read presence event
				var response map[string]interface{}
				err := conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "user_presence", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "user_123", data["user_id"])
				assert.Equal(t, "online", data["status"])
				assert.Contains(t, data, "last_seen")
			},
			expectError: false,
			description: "Should handle user presence updates correctly",
		},
		{
			name:       "message_read_receipt",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						msgType := msg["type"].(string)
						if msgType == "mark_read" {
							data := msg["data"].(map[string]interface{})

							// Send read receipt event
							readEvent := map[string]interface{}{
								"type": "message_read",
								"data": map[string]interface{}{
									"dialog_id":  data["dialog_id"],
									"message_id": data["message_id"],
									"user_id":    "user_123",
									"read_at":    "2023-12-01T15:30:00Z",
								},
								"timestamp": "2023-12-01T15:30:00Z",
							}
							conn.WriteJSON(readEvent)
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Send mark read message
				markReadMsg := map[string]interface{}{
					"type": "mark_read",
					"data": map[string]interface{}{
						"dialog_id":  "dialog_123",
						"message_id": "msg_456",
					},
					"timestamp": "2023-12-01T15:30:00Z",
				}
				err := conn.WriteJSON(markReadMsg)
				require.NoError(t, err)

				// Read read receipt event
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "message_read", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "dialog_123", data["dialog_id"])
				assert.Equal(t, "msg_456", data["message_id"])
				assert.Equal(t, "user_123", data["user_id"])
				assert.Contains(t, data, "read_at")
			},
			expectError: false,
			description: "Should handle message read receipts correctly",
		},
		{
			name:       "invalid_message_format",
			authHeader: "Bearer valid_jwt_token_12345",
			setupFunc: func(router *gin.Engine) {
				router.GET("/ws/messaging", func(c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						return
					}
					defer conn.Close()

					for {
						var msg map[string]interface{}
						err := conn.ReadJSON(&msg)
						if err != nil {
							break
						}

						// Check for invalid message format
						if _, hasType := msg["type"]; !hasType {
							errorMsg := map[string]interface{}{
								"type": "error",
								"data": map[string]interface{}{
									"code":    "INVALID_MESSAGE_FORMAT",
									"message": "Message type is required",
								},
								"timestamp": "2023-12-01T15:30:00Z",
							}
							conn.WriteJSON(errorMsg)
						}
					}
				})
			},
			testFunc: func(t *testing.T, conn *websocket.Conn) {
				// Send invalid message (missing type)
				invalidMsg := map[string]interface{}{
					"data": map[string]interface{}{
						"content": "Invalid message",
					},
				}
				err := conn.WriteJSON(invalidMsg)
				require.NoError(t, err)

				// Read error response
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				require.NoError(t, err)

				assert.Equal(t, "error", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "INVALID_MESSAGE_FORMAT", data["code"])
				assert.Contains(t, data["message"], "type")
			},
			expectError: false,
			description: "Should handle invalid message formats with error responses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router
			router := gin.New()

			// Add authentication middleware
			router.Use(func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")
				if authHeader == "" {
					authHeader = c.Query("token") // Allow token in query for WebSocket
				}

				if authHeader != "Bearer valid_jwt_token_12345" && authHeader != "valid_jwt_token_12345" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					c.Abort()
					return
				}

				c.Next()
			})

			// Setup test-specific handler
			tt.setupFunc(router)

			// Create test server
			server := httptest.NewServer(router)
			defer server.Close()

			// Convert HTTP URL to WebSocket URL
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging?token=valid_jwt_token_12345"

			// Connect to WebSocket
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if tt.expectError {
				assert.Error(t, err, tt.description)
				return
			}

			require.NoError(t, err, "WebSocket connection should succeed")
			defer conn.Close()

			// Set read deadline
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))

			// Run test function
			tt.testFunc(t, conn)
		})
	}
}

// TestMessagingWebSocketSecurity tests WebSocket security aspects
func TestMessagingWebSocketSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	t.Run("authentication_required", func(t *testing.T) {
		router := gin.New()

		// Add strict authentication
		router.GET("/ws/messaging", func(c *gin.Context) {
			token := c.Query("token")
			if token == "" {
				token = c.GetHeader("Authorization")
			}

			if token != "valid_token" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				return
			}

			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			// Keep connection alive
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		})

		server := httptest.NewServer(router)
		defer server.Close()

		// Try connecting without token
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging"
		_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.Error(t, err, "Connection without token should fail")

		// Try connecting with invalid token
		wsURLInvalid := wsURL + "?token=invalid_token"
		_, _, err = websocket.DefaultDialer.Dial(wsURLInvalid, nil)
		assert.Error(t, err, "Connection with invalid token should fail")

		// Connect with valid token should succeed
		wsURLValid := wsURL + "?token=valid_token"
		conn, _, err := websocket.DefaultDialer.Dial(wsURLValid, nil)
		assert.NoError(t, err, "Connection with valid token should succeed")
		if conn != nil {
			conn.Close()
		}
	})

	t.Run("message_validation", func(t *testing.T) {
		router := gin.New()
		router.GET("/ws/messaging", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			for {
				var msg map[string]interface{}
				err := conn.ReadJSON(&msg)
				if err != nil {
					break
				}

				// Validate message structure
				msgType, hasType := msg["type"].(string)
				if !hasType {
					errorMsg := map[string]interface{}{
						"type": "error",
						"data": map[string]interface{}{
							"code":    "INVALID_MESSAGE",
							"message": "Message type is required",
						},
					}
					conn.WriteJSON(errorMsg)
					continue
				}

				// Validate specific message types
				if msgType == "send_message" {
					data, hasData := msg["data"].(map[string]interface{})
					if !hasData {
						errorMsg := map[string]interface{}{
							"type": "error",
							"data": map[string]interface{}{
								"code":    "INVALID_MESSAGE",
								"message": "Message data is required",
							},
						}
						conn.WriteJSON(errorMsg)
						continue
					}

					// Check required fields
					requiredFields := []string{"dialog_id", "content", "message_type"}
					for _, field := range requiredFields {
						if _, exists := data[field]; !exists {
							errorMsg := map[string]interface{}{
								"type": "error",
								"data": map[string]interface{}{
									"code":    "MISSING_FIELD",
									"message": "Missing required field: " + field,
								},
							}
							conn.WriteJSON(errorMsg)
							break
						}
					}
				}

				// Echo valid messages
				conn.WriteJSON(map[string]interface{}{
					"type": "message_processed",
					"data": msg,
				})
			}
		})

		server := httptest.NewServer(router)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Test invalid message (missing type)
		invalidMsg := map[string]interface{}{
			"data": map[string]interface{}{
				"content": "test",
			},
		}
		err = conn.WriteJSON(invalidMsg)
		require.NoError(t, err)

		var errorResponse map[string]interface{}
		err = conn.ReadJSON(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "error", errorResponse["type"])

		// Test invalid send_message (missing required fields)
		invalidSendMsg := map[string]interface{}{
			"type": "send_message",
			"data": map[string]interface{}{
				"content": "test",
				// Missing dialog_id and message_type
			},
		}
		err = conn.WriteJSON(invalidSendMsg)
		require.NoError(t, err)

		err = conn.ReadJSON(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "error", errorResponse["type"])
	})

	t.Run("rate_limiting", func(t *testing.T) {
		router := gin.New()

		messageCount := 0
		router.GET("/ws/messaging", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			for {
				var msg map[string]interface{}
				err := conn.ReadJSON(&msg)
				if err != nil {
					break
				}

				messageCount++
				if messageCount > 10 { // Rate limit: 10 messages per connection
					errorMsg := map[string]interface{}{
						"type": "error",
						"data": map[string]interface{}{
							"code":    "RATE_LIMIT_EXCEEDED",
							"message": "Too many messages. Please slow down.",
						},
					}
					conn.WriteJSON(errorMsg)
					conn.Close()
					return
				}

				conn.WriteJSON(map[string]interface{}{
					"type": "message_processed",
					"data": msg,
				})
			}
		})

		server := httptest.NewServer(router)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Send messages beyond rate limit
		for i := 0; i < 15; i++ {
			msg := map[string]interface{}{
				"type": "test_message",
				"data": map[string]interface{}{
					"count": i,
				},
			}
			err = conn.WriteJSON(msg)
			require.NoError(t, err)

			var response map[string]interface{}
			err = conn.ReadJSON(&response)
			require.NoError(t, err)

			if i >= 10 {
				assert.Equal(t, "error", response["type"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "RATE_LIMIT_EXCEEDED", data["code"])
				break
			} else {
				assert.Equal(t, "message_processed", response["type"])
			}
		}
	})
}

// TestMessagingWebSocketPerformance tests WebSocket performance aspects
func TestMessagingWebSocketPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	t.Run("message_throughput", func(t *testing.T) {
		router := gin.New()
		router.GET("/ws/messaging", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			messageCount := 0
			start := time.Now()

			for {
				var msg map[string]interface{}
				err := conn.ReadJSON(&msg)
				if err != nil {
					break
				}

				messageCount++

				// Echo back immediately
				response := map[string]interface{}{
					"type": "echo",
					"data": msg,
					"count": messageCount,
				}
				conn.WriteJSON(response)

				// Stop after 100 messages
				if messageCount >= 100 {
					duration := time.Since(start)
					finalMsg := map[string]interface{}{
						"type": "performance_stats",
						"data": map[string]interface{}{
							"total_messages": messageCount,
							"duration_ms":    duration.Milliseconds(),
							"messages_per_second": float64(messageCount) / duration.Seconds(),
						},
					}
					conn.WriteJSON(finalMsg)
					break
				}
			}
		})

		server := httptest.NewServer(router)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		start := time.Now()

		// Send 100 messages rapidly
		for i := 0; i < 100; i++ {
			msg := map[string]interface{}{
				"type": "performance_test",
				"data": map[string]interface{}{
					"message_id": i,
					"timestamp":  time.Now().UnixNano(),
				},
			}
			err = conn.WriteJSON(msg)
			require.NoError(t, err)

			// Read response
			var response map[string]interface{}
			err = conn.ReadJSON(&response)
			require.NoError(t, err)

			if response["type"] == "performance_stats" {
				data := response["data"].(map[string]interface{})
				messagesPerSecond := data["messages_per_second"].(float64)

				// Performance assertion: should handle at least 50 messages per second
				assert.GreaterOrEqual(t, messagesPerSecond, 50.0,
					"WebSocket should handle at least 50 messages per second, got %.2f", messagesPerSecond)
				break
			}
		}

		duration := time.Since(start)
		assert.Less(t, duration, 5*time.Second, "100 messages should be processed within 5 seconds")
	})

	t.Run("connection_scalability", func(t *testing.T) {
		// Test multiple concurrent connections
		router := gin.New()

		activeConnections := 0
		router.GET("/ws/messaging", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			activeConnections++

			// Send connection info
			info := map[string]interface{}{
				"type": "connection_info",
				"data": map[string]interface{}{
					"connection_id": activeConnections,
					"active_connections": activeConnections,
				},
			}
			conn.WriteJSON(info)

			// Keep connection alive for a short time
			time.Sleep(100 * time.Millisecond)
			activeConnections--
		})

		server := httptest.NewServer(router)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/messaging"

		// Create multiple concurrent connections
		connections := make([]*websocket.Conn, 10)
		for i := 0; i < 10; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			assert.NoError(t, err, "Connection %d should succeed", i)
			if conn != nil {
				connections[i] = conn
			}
		}

		// Clean up connections
		for _, conn := range connections {
			if conn != nil {
				conn.Close()
			}
		}
	})
}