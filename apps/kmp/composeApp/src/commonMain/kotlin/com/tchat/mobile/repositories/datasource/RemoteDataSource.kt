package com.tchat.mobile.repositories.datasource

import com.tchat.mobile.models.*
import kotlinx.datetime.Instant
import kotlinx.coroutines.flow.Flow

/**
 * RemoteDataSource interface for API communication
 *
 * Handles all server-side operations for the chat application including:
 * - Message synchronization
 * - Chat management
 * - Conflict resolution
 * - Real-time updates
 */
interface RemoteDataSource {

    // Connection Management
    suspend fun connect(): Result<Unit>
    suspend fun disconnect(): Result<Unit>
    fun getConnectionState(): Flow<ConnectionState>

    // Authentication & User Management
    suspend fun authenticateUser(token: String): Result<User>
    suspend fun refreshToken(): Result<String>
    suspend fun getUserProfile(userId: String): Result<UserProfile>

    // Chat Operations
    suspend fun getChatSessions(userId: String): Result<List<ChatSession>>
    suspend fun getChatSession(chatId: String): Result<ChatSession>
    suspend fun createChatSession(chat: ChatSession): Result<ChatSession>
    suspend fun updateChatSession(chat: ChatSession): Result<ChatSession>
    suspend fun deleteChatSession(chatId: String): Result<Unit>

    // Message Operations
    suspend fun getMessages(chatId: String, limit: Int = 50, before: Instant? = null): Result<List<ChatMessage>>
    suspend fun sendMessage(message: ChatMessage): Result<ChatMessage>
    suspend fun updateMessage(message: ChatMessage): Result<ChatMessage>
    suspend fun deleteMessage(messageId: String): Result<Unit>
    suspend fun markMessageAsRead(messageId: String, userId: String): Result<Unit>

    // Sync Operations
    suspend fun getMessagesSince(chatId: String, timestamp: Instant): Result<List<ChatMessage>>
    suspend fun pushPendingOperations(operations: List<SyncOperation>): Result<List<SyncOperation>>
    suspend fun resolveConflicts(conflicts: List<MessageConflict>): Result<List<ConflictResolution>>
    suspend fun getServerTimestamp(): Result<Instant>

    // Bulk Operations
    suspend fun syncChatData(chatId: String, lastSyncTimestamp: Instant?): Result<SyncResult>
    suspend fun batchUpdateMessages(messages: List<ChatMessage>): Result<List<ChatMessage>>
    suspend fun batchDeleteMessages(messageIds: List<String>): Result<Unit>

    // Real-time Updates
    fun subscribeToMessages(chatId: String): Flow<ChatMessage>
    fun subscribeToPresence(chatId: String): Flow<List<PresenceUpdate>>
    fun subscribeToTypingIndicators(chatId: String): Flow<TypingIndicator>
    suspend fun sendTypingIndicator(chatId: String, isTyping: Boolean): Result<Unit>

    // File Upload & Media
    suspend fun uploadFile(file: ByteArray, fileName: String, mimeType: String): Result<String>
    suspend fun uploadImage(image: ByteArray, fileName: String): Result<String>
    suspend fun uploadVideo(video: ByteArray, fileName: String): Result<String>
    suspend fun uploadAudio(audio: ByteArray, fileName: String): Result<String>

    // Search Operations
    suspend fun searchMessages(query: String, chatId: String? = null): Result<List<ChatMessage>>
    suspend fun searchChats(query: String): Result<List<ChatSession>>

    // Analytics & Metrics
    suspend fun trackMessageDelivery(messageId: String, status: DeliveryStatus): Result<Unit>
    suspend fun reportMessageRead(messageId: String, readAt: Instant): Result<Unit>
    suspend fun getMessageAnalytics(chatId: String): Result<MessageAnalytics>

    // Error Recovery
    suspend fun retryFailedOperation(operation: SyncOperation): Result<SyncOperation>
    suspend fun validateDataIntegrity(chatId: String): Result<DataIntegrityReport>

    // Health & Monitoring
    suspend fun healthCheck(): Result<ServerHealth>
    suspend fun getServerCapabilities(): Result<ServerCapabilities>
}

/**
 * Data models for remote operations
 */
data class User(
    val id: String,
    val email: String,
    val displayName: String,
    val avatar: String? = null,
    val isActive: Boolean = true
)

data class UserProfile(
    val userId: String,
    val displayName: String,
    val username: String,
    val avatar: String? = null,
    val bio: String? = null,
    val isVerified: Boolean = false,
    val isOnline: Boolean = false,
    val lastSeen: Instant? = null
)

data class TypingIndicator(
    val chatId: String,
    val userId: String,
    val isTyping: Boolean,
    val timestamp: Instant
)

data class MessageAnalytics(
    val chatId: String,
    val totalMessages: Int,
    val averageResponseTime: Long,
    val activeParticipants: Int,
    val messageFrequency: Map<String, Int>
)

data class DataIntegrityReport(
    val chatId: String,
    val isValid: Boolean,
    val issues: List<String>,
    val lastChecked: Instant
)

data class ServerHealth(
    val status: String,
    val uptime: Long,
    val responseTime: Long,
    val activeConnections: Int
)

data class ServerCapabilities(
    val maxFileSize: Long,
    val supportedFileTypes: List<String>,
    val maxMessageLength: Int,
    val supportsBulkOperations: Boolean,
    val supportsRealTime: Boolean
)