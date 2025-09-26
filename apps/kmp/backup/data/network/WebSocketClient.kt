package com.tchat.mobile.data.network

import io.ktor.client.*
import io.ktor.client.plugins.websocket.*
import io.ktor.websocket.*
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.encodeToString
import kotlinx.serialization.decodeFromString
import io.github.aakira.napier.Napier
import com.tchat.mobile.data.models.Message

/**
 * WebSocket client for real-time messaging and notifications
 */
class WebSocketClient(
    private val httpClient: HttpClient,
    private val json: Json = Json { ignoreUnknownKeys = true }
) {
    private var webSocketSession: DefaultClientWebSocketSession? = null
    private var connectionJob: Job? = null
    private val coroutineScope = CoroutineScope(Dispatchers.Default + SupervisorJob())

    // Channels for different event types
    private val _messageEvents = Channel<WebSocketEvent>(Channel.UNLIMITED)
    val messageEvents: Flow<WebSocketEvent> = _messageEvents.receiveAsFlow()

    private val _connectionState = MutableStateFlow(ConnectionState.DISCONNECTED)
    val connectionState: StateFlow<ConnectionState> = _connectionState.asStateFlow()

    private var authToken: String? = null
    private var reconnectAttempts = 0
    private val maxReconnectAttempts = 5
    private val reconnectDelayMs = 1000L

    /**
     * Connection states
     */
    enum class ConnectionState {
        DISCONNECTED,
        CONNECTING,
        CONNECTED,
        RECONNECTING,
        ERROR
    }

    /**
     * Connect to WebSocket server
     */
    suspend fun connect(wsUrl: String, token: String) {
        authToken = token
        _connectionState.value = ConnectionState.CONNECTING

        try {
            connectionJob?.cancel()
            connectionJob = coroutineScope.launch {
                connectWithRetry(wsUrl)
            }
        } catch (e: Exception) {
            Napier.e("WebSocket connection error: ${e.message}", e)
            _connectionState.value = ConnectionState.ERROR
        }
    }

    /**
     * Connect with automatic retry logic
     */
    private suspend fun connectWithRetry(wsUrl: String) {
        while (reconnectAttempts < maxReconnectAttempts && coroutineScope.isActive) {
            try {
                httpClient.webSocket(
                    urlString = wsUrl,
                    request = {
                        authToken?.let { token ->
                            header("Authorization", "Bearer $token")
                        }
                    }
                ) {
                    webSocketSession = this
                    _connectionState.value = ConnectionState.CONNECTED
                    reconnectAttempts = 0

                    Napier.i("WebSocket connected successfully")

                    // Handle incoming messages
                    handleIncomingMessages()
                }
            } catch (e: Exception) {
                Napier.w("WebSocket connection failed, attempt ${reconnectAttempts + 1}: ${e.message}")
                reconnectAttempts++

                if (reconnectAttempts < maxReconnectAttempts) {
                    _connectionState.value = ConnectionState.RECONNECTING
                    delay(reconnectDelayMs * reconnectAttempts) // Exponential backoff
                } else {
                    Napier.e("Max reconnection attempts reached")
                    _connectionState.value = ConnectionState.ERROR
                    break
                }
            }
        }
    }

    /**
     * Handle incoming WebSocket messages
     */
    private suspend fun DefaultClientWebSocketSession.handleIncomingMessages() {
        try {
            for (frame in incoming) {
                when (frame) {
                    is Frame.Text -> {
                        val text = frame.readText()
                        try {
                            val event = json.decodeFromString<WebSocketEvent>(text)
                            _messageEvents.trySend(event)
                            Napier.d("Received WebSocket event: ${event.type}")
                        } catch (e: Exception) {
                            Napier.w("Failed to parse WebSocket message: $text", e)
                        }
                    }
                    is Frame.Binary -> {
                        // Handle binary frames if needed
                        Napier.d("Received binary WebSocket frame")
                    }
                    is Frame.Close -> {
                        val reason = frame.readReason()
                        Napier.i("WebSocket closed: ${reason?.message}")
                        _connectionState.value = ConnectionState.DISCONNECTED
                        break
                    }
                    else -> {
                        // Handle other frame types
                    }
                }
            }
        } catch (e: Exception) {
            Napier.e("Error handling WebSocket messages: ${e.message}", e)
            _connectionState.value = ConnectionState.ERROR
        } finally {
            webSocketSession = null
        }
    }

    /**
     * Send message through WebSocket
     */
    suspend fun sendMessage(event: WebSocketEvent): Boolean {
        return try {
            val session = webSocketSession
            if (session != null && _connectionState.value == ConnectionState.CONNECTED) {
                val jsonString = json.encodeToString(event)
                session.send(Frame.Text(jsonString))
                Napier.d("Sent WebSocket message: ${event.type}")
                true
            } else {
                Napier.w("WebSocket not connected, cannot send message")
                false
            }
        } catch (e: Exception) {
            Napier.e("Failed to send WebSocket message: ${e.message}", e)
            false
        }
    }

    /**
     * Send typing indicator
     */
    suspend fun sendTypingIndicator(dialogId: String, isTyping: Boolean) {
        val event = WebSocketEvent(
            type = if (isTyping) "typing_start" else "typing_stop",
            data = mapOf(
                "dialog_id" to dialogId,
                "timestamp" to kotlinx.datetime.Clock.System.now().toString()
            )
        )
        sendMessage(event)
    }

    /**
     * Send presence update
     */
    suspend fun sendPresenceUpdate(status: String) {
        val event = WebSocketEvent(
            type = "presence_update",
            data = mapOf(
                "status" to status,
                "timestamp" to kotlinx.datetime.Clock.System.now().toString()
            )
        )
        sendMessage(event)
    }

    /**
     * Send read receipt
     */
    suspend fun sendReadReceipt(dialogId: String, messageId: String) {
        val event = WebSocketEvent(
            type = "message_read",
            data = mapOf(
                "dialog_id" to dialogId,
                "message_id" to messageId,
                "timestamp" to kotlinx.datetime.Clock.System.now().toString()
            )
        )
        sendMessage(event)
    }

    /**
     * Join dialog room for real-time updates
     */
    suspend fun joinDialogRoom(dialogId: String) {
        val event = WebSocketEvent(
            type = "join_dialog",
            data = mapOf(
                "dialog_id" to dialogId,
                "timestamp" to kotlinx.datetime.Clock.System.now().toString()
            )
        )
        sendMessage(event)
    }

    /**
     * Leave dialog room
     */
    suspend fun leaveDialogRoom(dialogId: String) {
        val event = WebSocketEvent(
            type = "leave_dialog",
            data = mapOf(
                "dialog_id" to dialogId,
                "timestamp" to kotlinx.datetime.Clock.System.now().toString()
            )
        )
        sendMessage(event)
    }

    /**
     * Disconnect WebSocket connection
     */
    suspend fun disconnect() {
        try {
            webSocketSession?.close(CloseReason(CloseReason.Codes.NORMAL, "Client disconnect"))
            connectionJob?.cancelAndJoin()
            _connectionState.value = ConnectionState.DISCONNECTED
            reconnectAttempts = 0
            Napier.i("WebSocket disconnected")
        } catch (e: Exception) {
            Napier.e("Error disconnecting WebSocket: ${e.message}", e)
        }
    }

    /**
     * Clean up resources
     */
    fun cleanup() {
        coroutineScope.launch {
            disconnect()
            _messageEvents.close()
        }
        coroutineScope.cancel()
    }
}

/**
 * WebSocket event data structure
 */
@Serializable
data class WebSocketEvent(
    val type: String,
    val data: Map<String, String> = emptyMap(),
    val timestamp: String = kotlinx.datetime.Clock.System.now().toString(),
    val id: String? = null
)

/**
 * Specific WebSocket events for messaging
 */
object WebSocketEvents {
    // Message events
    const val MESSAGE_NEW = "message_new"
    const val MESSAGE_UPDATED = "message_updated"
    const val MESSAGE_DELETED = "message_deleted"
    const val MESSAGE_REACTION_ADDED = "message_reaction_added"
    const val MESSAGE_REACTION_REMOVED = "message_reaction_removed"

    // Dialog events
    const val DIALOG_UPDATED = "dialog_updated"
    const val DIALOG_PARTICIPANT_ADDED = "dialog_participant_added"
    const val DIALOG_PARTICIPANT_REMOVED = "dialog_participant_removed"

    // Presence events
    const val USER_ONLINE = "user_online"
    const val USER_OFFLINE = "user_offline"
    const val USER_TYPING_START = "typing_start"
    const val USER_TYPING_STOP = "typing_stop"

    // Read receipt events
    const val MESSAGE_READ = "message_read"
    const val MESSAGE_DELIVERED = "message_delivered"

    // Connection events
    const val CONNECTION_ACK = "connection_ack"
    const val PING = "ping"
    const val PONG = "pong"

    // Error events
    const val ERROR = "error"
    const val AUTH_ERROR = "auth_error"
}

/**
 * WebSocket message parsers for specific event types
 */
object WebSocketEventParser {

    /**
     * Parse new message event
     */
    fun parseNewMessage(event: WebSocketEvent): Message? {
        return try {
            val messageData = event.data
            // Convert map data to Message object
            // This would need proper JSON deserialization
            null // Placeholder - implement based on actual message format
        } catch (e: Exception) {
            Napier.w("Failed to parse new message event: ${e.message}")
            null
        }
    }

    /**
     * Parse typing indicator event
     */
    fun parseTypingIndicator(event: WebSocketEvent): TypingIndicator? {
        return try {
            TypingIndicator(
                userId = event.data["user_id"] ?: return null,
                dialogId = event.data["dialog_id"] ?: return null,
                isTyping = event.type == WebSocketEvents.USER_TYPING_START,
                timestamp = event.timestamp
            )
        } catch (e: Exception) {
            Napier.w("Failed to parse typing indicator: ${e.message}")
            null
        }
    }

    /**
     * Parse presence update event
     */
    fun parsePresenceUpdate(event: WebSocketEvent): PresenceUpdate? {
        return try {
            PresenceUpdate(
                userId = event.data["user_id"] ?: return null,
                status = event.data["status"] ?: return null,
                timestamp = event.timestamp
            )
        } catch (e: Exception) {
            Napier.w("Failed to parse presence update: ${e.message}")
            null
        }
    }
}

/**
 * Data classes for parsed WebSocket events
 */
@Serializable
data class TypingIndicator(
    val userId: String,
    val dialogId: String,
    val isTyping: Boolean,
    val timestamp: String
)

@Serializable
data class PresenceUpdate(
    val userId: String,
    val status: String,
    val timestamp: String
)

@Serializable
data class ReadReceipt(
    val userId: String,
    val dialogId: String,
    val messageId: String,
    val timestamp: String
)

/**
 * WebSocket connection manager for handling multiple connections
 */
class WebSocketManager {
    private val connections = mutableMapOf<String, WebSocketClient>()

    /**
     * Get or create WebSocket client for endpoint
     */
    fun getClient(endpoint: String, httpClient: HttpClient): WebSocketClient {
        return connections.getOrPut(endpoint) {
            WebSocketClient(httpClient)
        }
    }

    /**
     * Remove client for endpoint
     */
    fun removeClient(endpoint: String) {
        connections.remove(endpoint)?.cleanup()
    }

    /**
     * Clean up all connections
     */
    fun cleanup() {
        connections.values.forEach { it.cleanup() }
        connections.clear()
    }
}