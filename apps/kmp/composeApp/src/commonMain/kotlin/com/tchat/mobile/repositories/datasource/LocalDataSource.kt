package com.tchat.mobile.repositories.datasource

import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.datetime.Instant

/**
 * Local Data Source - SQLDelight Implementation Interface
 * Responsible for: Local persistence, caching, offline access
 */
interface LocalDataSource {

    // Message Operations
    suspend fun getMessages(chatId: String): Flow<List<ChatMessage>>
    suspend fun getMessagesByChatId(chatId: String, limit: Int = 50): List<ChatMessage>
    suspend fun getMessage(messageId: String): ChatMessage?
    suspend fun saveMessage(message: ChatMessage): Result<Unit>
    suspend fun updateMessageStatus(messageId: String, status: DeliveryStatus): Result<Unit>
    suspend fun markMessageAsSynced(messageId: String, serverTimestamp: Instant): Result<Unit>
    suspend fun updateMessageSyncStatus(messageId: String, syncStatus: SyncStatus): Result<Unit>
    suspend fun getPendingSyncMessages(): List<ChatMessage>
    suspend fun getUnsyncedMessages(): List<ChatMessage>
    suspend fun updateMessageVersion(messageId: String, version: Long, checksum: String): Result<Unit>
    suspend fun updateDeliveryStatus(messageId: String, status: DeliveryStatus): Result<Unit>
    suspend fun markMessageAsRead(messageId: String, readBy: List<String>): Result<Unit>
    suspend fun getMessagesSince(chatId: String, since: Instant): List<ChatMessage>
    suspend fun searchMessages(chatId: String, query: String): List<ChatMessage>
    suspend fun deleteMessage(messageId: String): Result<Unit>
    suspend fun editMessage(messageId: String, newContent: MessageContent, editedAt: Instant): Result<Unit>

    // Chat Operations
    suspend fun getChatSessions(): Flow<List<ChatSession>>
    suspend fun getAllChatSessions(): List<ChatSession>
    suspend fun getChatSession(chatId: String): Result<ChatSession>
    suspend fun saveChatSession(chat: ChatSession): Result<Unit>
    suspend fun updateChatSession(chat: ChatSession): Result<Unit>
    suspend fun deleteChatSession(chatId: String): Result<Unit>
    suspend fun updateLastActivity(chatId: String, timestamp: Instant): Result<Unit>
    suspend fun updateChatUnreadCount(chatId: String, count: Int): Result<Unit>
    suspend fun markChatAsRead(chatId: String): Result<Unit>
    suspend fun muteChat(chatId: String, muted: Boolean): Result<Unit>
    suspend fun archiveChat(chatId: String, archived: Boolean): Result<Unit>
    suspend fun pinChat(chatId: String, pinned: Boolean): Result<Unit>

    // Sync Metadata Operations
    suspend fun getLastSyncTimestamp(chatId: String): Instant?
    suspend fun updateSyncTimestamp(chatId: String, timestamp: Instant): Result<Unit>
    suspend fun getSyncMetadata(chatId: String): SyncMetadata?
    suspend fun updateSyncStatus(chatId: String, status: SyncState, errorMessage: String? = null): Result<Unit>
    suspend fun updateLastSuccessfulSync(chatId: String, timestamp: Instant): Result<Unit>
    suspend fun incrementPendingOperations(chatId: String): Result<Unit>
    suspend fun decrementPendingOperations(chatId: String): Result<Unit>
    suspend fun getAllSyncMetadata(): List<SyncMetadata>

    // Sync Operations Queue
    suspend fun saveSyncOperation(operation: SyncOperation): Result<Unit>
    suspend fun getPendingOperations(): List<SyncOperation>
    suspend fun getSyncOperationsByChatId(chatId: String): List<SyncOperation>
    suspend fun getRetryableOperations(currentTime: Instant): List<SyncOperation>
    suspend fun updateSyncOperationStatus(
        operationId: String,
        status: SyncOperationStatus,
        errorMessage: String? = null,
        retryCount: Int? = null,
        scheduledAt: Instant? = null
    ): Result<Unit>
    suspend fun markOperationCompleted(operationId: String): Result<Unit>
    suspend fun markOperationFailed(operationId: String, errorMessage: String): Result<Unit>
    suspend fun deleteSyncOperation(operationId: String): Result<Unit>
    suspend fun cleanupCompletedOperations(olderThan: Instant): Result<Int>

    // Conflict Management
    suspend fun saveConflict(conflict: MessageConflict): Result<Unit>
    suspend fun getConflictsByChatId(chatId: String): List<MessageConflict>
    suspend fun getConflictsByMessageId(messageId: String): List<MessageConflict>
    suspend fun getAllUnresolvedConflicts(): List<MessageConflict>
    suspend fun updateConflictResolution(
        conflictId: String,
        resolution: ConflictResolution
    ): Result<Unit>
    suspend fun deleteResolvedConflicts(olderThan: Instant): Result<Int>

    // Participant Operations
    suspend fun getChatParticipants(chatId: String): List<ChatParticipant>
    suspend fun saveChatParticipant(participant: ChatParticipant): Result<Unit>
    suspend fun updateParticipantStatus(
        chatId: String,
        participantId: String,
        status: ParticipantStatus,
        lastSeen: Instant? = null
    ): Result<Unit>
    suspend fun removeParticipant(chatId: String, participantId: String): Result<Unit>

    // Attachment Operations
    suspend fun getMessageAttachments(messageId: String): List<MessageAttachment>
    suspend fun saveMessageAttachment(attachment: MessageAttachment): Result<Unit>
    suspend fun deleteMessageAttachments(messageId: String): Result<Unit>

    // Reaction Operations
    suspend fun getMessageReactions(messageId: String): List<MessageReaction>
    suspend fun saveMessageReaction(reaction: MessageReaction): Result<Unit>
    suspend fun deleteMessageReaction(messageId: String, userId: String, emoji: String): Result<Unit>

    // Utility Operations
    suspend fun clearAllData(): Result<Unit>
    suspend fun getStorageInfo(): StorageInfo
    suspend fun optimizeDatabase(): Result<Unit>
}

// Supporting Data Classes for LocalDataSource
data class SyncMetadata(
    val chatId: String,
    val lastSyncTimestamp: Instant?,
    val lastSuccessfulSync: Instant?,
    val syncStatus: SyncState,
    val pendingOperations: Int,
    val conflictCount: Int,
    val errorMessage: String?,
    val createdAt: Instant,
    val updatedAt: Instant
)

enum class SyncStatus {
    PENDING,
    SYNCING,
    SYNCED,
    FAILED
}

enum class SyncState {
    IDLE,
    SYNCING,
    SYNCED,
    ERROR,
    CONFLICT_PENDING
}

enum class SyncOperationStatus {
    PENDING,
    IN_PROGRESS,
    COMPLETED,
    FAILED
}

data class StorageInfo(
    val totalMessages: Long,
    val totalChats: Long,
    val pendingOperations: Long,
    val unresolvedConflicts: Long,
    val databaseSize: Long,
    val lastOptimized: Instant?
)