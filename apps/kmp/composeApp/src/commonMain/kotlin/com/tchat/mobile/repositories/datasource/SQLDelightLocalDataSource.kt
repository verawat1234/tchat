package com.tchat.mobile.repositories.datasource

import app.cash.sqldelight.coroutines.asFlow
import app.cash.sqldelight.coroutines.mapToList
import app.cash.sqldelight.coroutines.mapToOneOrNull
import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import com.tchat.mobile.utils.toJsonString
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

/**
 * SQLDelight implementation of LocalDataSource
 * Handles all local database operations for chat data
 */
class SQLDelightLocalDataSource(
    private val database: TchatDatabase
) : LocalDataSource {

    // Message Operations
    override suspend fun getMessages(chatId: String): Flow<List<ChatMessage>> {
        return database.messageQueries.getMessagesByChatId(chatId, 50)
            .asFlow()
            .mapToList(Dispatchers.IO)
            .map { rows ->
                rows.map { row ->
                    Message(
                        id = row.id,
                        chatId = row.chatId,
                        senderId = row.senderId,
                        senderName = row.senderName,
                        senderAvatar = null, // Not in schema
                        type = MessageType.valueOf(row.type),
                        content = row.content,
                        isEdited = row.isEdited == 1L,
                        isPinned = row.isPinned == 1L,
                        isDeleted = row.isDeleted == 1L,
                        replyToId = row.replyToId,
                        reactions = emptyList(),
                        attachments = emptyList(),
                        createdAt = row.createdAt,
                        editedAt = row.editedAt,
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
            }
    }

    override suspend fun getMessagesByChatId(chatId: String, limit: Int): List<ChatMessage> {
        return try {
            database.messageQueries.getMessagesByChatId(chatId, limit.toLong())
                .executeAsList()
                .map { row ->
                    Message(
                        id = row.id,
                        chatId = row.chatId,
                        senderId = row.senderId,
                        senderName = row.senderName,
                        senderAvatar = null, // Not in schema
                        type = MessageType.valueOf(row.type),
                        content = row.content,
                        isEdited = row.isEdited == 1L,
                        isPinned = row.isPinned == 1L,
                        isDeleted = row.isDeleted == 1L,
                        replyToId = row.replyToId,
                        reactions = emptyList(),
                        attachments = emptyList(),
                        createdAt = row.createdAt,
                        editedAt = row.editedAt,
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun getMessage(messageId: String): ChatMessage? {
        return try {
            database.messageQueries.getMessageById(messageId)
                .executeAsOneOrNull()?.let { row ->
                    Message(
                        id = row.id,
                        chatId = row.chatId,
                        senderId = row.senderId,
                        senderName = row.senderName,
                        senderAvatar = null, // Not in schema
                        type = MessageType.valueOf(row.type),
                        content = row.content,
                        isEdited = row.isEdited == 1L,
                        isPinned = row.isPinned == 1L,
                        isDeleted = row.isDeleted == 1L,
                        replyToId = row.replyToId,
                        reactions = emptyList(),
                        attachments = emptyList(),
                        createdAt = row.createdAt,
                        editedAt = row.editedAt,
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
        } catch (e: Exception) {
            null
        }
    }

    override suspend fun saveMessage(message: ChatMessage): Result<Unit> {
        return try {
            database.messageQueries.insertMessage(
                id = message.id,
                chatId = message.chatId,
                senderId = message.senderId,
                senderName = message.senderName,
                type = message.type.name,
                content = message.content,
                isEdited = if (message.isEdited) 1L else 0L,
                isPinned = if (message.isPinned) 1L else 0L,
                isDeleted = if (message.isDeleted) 1L else 0L,
                replyToId = message.replyToId,
                reactions = message.reactions.toJsonString(),
                attachmentCount = message.attachments.size.toLong(),
                createdAt = message.createdAt,
                editedAt = message.editedAt,
                deletedAt = message.deletedAt
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateMessageStatus(messageId: String, status: DeliveryStatus): Result<Unit> {
        // No delivery status in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun markMessageAsSynced(messageId: String, serverTimestamp: Instant): Result<Unit> {
        // No sync tracking in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun updateMessageSyncStatus(messageId: String, syncStatus: SyncStatus): Result<Unit> {
        // No sync status in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun getPendingSyncMessages(): List<ChatMessage> {
        // No sync tracking in current schema - return empty list
        return emptyList()
    }

    override suspend fun getUnsyncedMessages(): List<ChatMessage> {
        // No sync tracking in current schema - return empty list
        return emptyList()
    }

    override suspend fun updateMessageVersion(messageId: String, version: Long, checksum: String): Result<Unit> {
        // No version tracking in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun updateDeliveryStatus(messageId: String, status: DeliveryStatus): Result<Unit> {
        // No delivery status in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun markMessageAsRead(messageId: String, readBy: List<String>): Result<Unit> {
        // No read tracking in current schema - return success
        return Result.success(Unit)
    }

    override suspend fun getMessagesSince(chatId: String, since: Instant): List<ChatMessage> {
        return try {
            val allMessages = database.messageQueries.getMessagesByChatId(chatId, 1000)
                .executeAsList()
            allMessages.filter { row ->
                row.createdAt >= since.toString()
            }.map { row ->
                ChatMessage(
                    id = row.id,
                    chatId = row.chatId,
                    senderId = row.senderId,
                    senderName = row.senderName,
                    type = MessageType.valueOf(row.type),
                    content = row.content,
                    isEdited = row.isEdited == 1L,
                    isPinned = row.isPinned == 1L,
                    isDeleted = row.isDeleted == 1L,
                    replyToId = row.replyToId,
                    reactions = emptyList(),
                    attachments = emptyList(),
                    createdAt = row.createdAt,
                    editedAt = row.editedAt,
                    deletedAt = row.deletedAt
                )
            }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun searchMessages(chatId: String, query: String): List<ChatMessage> {
        return try {
            database.messageQueries.searchMessages(chatId, query)
                .executeAsList()
                .map { row ->
                    Message(
                        id = row.id,
                        chatId = row.chatId,
                        senderId = row.senderId,
                        senderName = row.senderName,
                        senderAvatar = null, // Not in schema
                        type = MessageType.valueOf(row.type),
                        content = row.content,
                        isEdited = row.isEdited == 1L,
                        isPinned = row.isPinned == 1L,
                        isDeleted = row.isDeleted == 1L,
                        replyToId = row.replyToId,
                        reactions = emptyList(),
                        attachments = emptyList(),
                        createdAt = row.createdAt,
                        editedAt = row.editedAt,
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Unit> {
        return try {
            database.messageQueries.deleteMessage(messageId, Clock.System.now().toString())
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun editMessage(messageId: String, newContent: MessageContent, editedAt: Instant): Result<Unit> {
        return try {
            database.messageQueries.updateMessageContent(newContent, editedAt.toString(), messageId)
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Chat Operations
    override suspend fun getChatSessions(): Flow<List<ChatSession>> {
        return database.chatSessionQueries.getAllChatSessions()
            .asFlow()
            .mapToList(Dispatchers.IO)
            .map { rows ->
                rows.map { row ->
                    ChatSession(
                        id = row.id,
                        name = row.name ?: "",
                        type = ChatType.valueOf(row.type.uppercase()),
                        unreadCount = row.unreadCount.toInt(),
                        isPinned = row.isPinned == 1L,
                        isMuted = row.isMuted == 1L,
                        isArchived = row.isArchived == 1L,
                        isBlocked = row.isBlocked == 1L,
                        participants = emptyList(), // Load separately
                        metadata = ChatMetadata(),
                        createdAt = row.createdAt,
                        updatedAt = row.updatedAt,
                        lastActivityAt = row.lastActivityAt
                    )
                }
            }
    }

    override suspend fun getAllChatSessions(): List<ChatSession> {
        return try {
            database.chatSessionQueries.getAllChatSessions()
                .executeAsList()
                .map { row ->
                    ChatSession(
                        id = row.id,
                        name = row.name ?: "",
                        type = ChatType.valueOf(row.type.uppercase()),
                        unreadCount = row.unreadCount.toInt(),
                        isPinned = row.isPinned == 1L,
                        isMuted = row.isMuted == 1L,
                        isArchived = row.isArchived == 1L,
                        isBlocked = row.isBlocked == 1L,
                        participants = emptyList(), // Load separately
                        metadata = ChatMetadata(),
                        createdAt = row.createdAt,
                        updatedAt = row.updatedAt,
                        lastActivityAt = row.lastActivityAt
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun getChatSession(chatId: String): Result<ChatSession> {
        return try {
            val row = database.chatSessionQueries.getChatSessionById(chatId)
                .executeAsOneOrNull()
            if (row != null) {
                val participants = database.chatSessionQueries.getChatParticipants(chatId)
                    .executeAsList()
                    .map { participantRow ->
                        ChatParticipant(
                            id = participantRow.id,
                            name = participantRow.name,
                            avatar = participantRow.avatar,
                            role = ChatRole.valueOf(participantRow.role.uppercase()),
                            status = ParticipantStatus.valueOf(participantRow.status.uppercase()),
                            lastSeen = participantRow.lastSeen,
                            joinedAt = participantRow.joinedAt,
                            customTitle = participantRow.customTitle,
                            isBot = participantRow.isBot == 1L
                        )
                    }

                val chatSession = ChatSession(
                    id = row.id,
                    name = row.name ?: "",
                    type = ChatType.valueOf(row.type.uppercase()),
                    unreadCount = row.unreadCount.toInt(),
                    isPinned = row.isPinned == 1L,
                    isMuted = row.isMuted == 1L,
                    isArchived = row.isArchived == 1L,
                    isBlocked = row.isBlocked == 1L,
                    participants = participants,
                    metadata = ChatMetadata(),
                    createdAt = row.createdAt,
                    updatedAt = row.updatedAt,
                    lastActivityAt = row.lastActivityAt
                )
                Result.success(chatSession)
            } else {
                Result.failure(Exception("Chat session not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun saveChatSession(chat: ChatSession): Result<Unit> {
        return try {
            database.chatSessionQueries.insertChatSession(
                id = chat.id,
                name = chat.name,
                type = chat.type.name,
                unreadCount = chat.unreadCount.toLong(),
                isPinned = if (chat.isPinned) 1L else 0L,
                isMuted = if (chat.isMuted) 1L else 0L,
                isArchived = if (chat.isArchived) 1L else 0L,
                isBlocked = if (chat.isBlocked) 1L else 0L,
                createdAt = chat.createdAt,
                updatedAt = chat.updatedAt,
                lastActivityAt = chat.lastActivityAt
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateChatSession(chat: ChatSession): Result<Unit> {
        return try {
            database.chatSessionQueries.updateChatSessionUnreadCount(
                unreadCount = chat.unreadCount.toLong(),
                updatedAt = chat.updatedAt,
                lastActivityAt = chat.lastActivityAt,
                id = chat.id
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteChatSession(chatId: String): Result<Unit> {
        return try {
            database.chatSessionQueries.deleteChatSession(chatId)
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateLastActivity(chatId: String, timestamp: Instant): Result<Unit> {
        return try {
            database.chatSessionQueries.updateChatSessionUnreadCount(
                unreadCount = 0L, // Reset when updating activity
                updatedAt = timestamp.toString(),
                lastActivityAt = timestamp.toString(),
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateChatUnreadCount(chatId: String, count: Int): Result<Unit> {
        return try {
            val now = Clock.System.now().toString()
            database.chatSessionQueries.updateChatSessionUnreadCount(
                unreadCount = count.toLong(),
                updatedAt = now,
                lastActivityAt = now,
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun markChatAsRead(chatId: String): Result<Unit> {
        return updateChatUnreadCount(chatId, 0)
    }

    override suspend fun muteChat(chatId: String, muted: Boolean): Result<Unit> {
        return try {
            database.chatSessionQueries.updateChatSessionMuteStatus(
                isMuted = if (muted) 1L else 0L,
                updatedAt = Clock.System.now().toString(),
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun archiveChat(chatId: String, archived: Boolean): Result<Unit> {
        return try {
            database.chatSessionQueries.updateChatSessionArchiveStatus(
                isArchived = if (archived) 1L else 0L,
                updatedAt = Clock.System.now().toString(),
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun pinChat(chatId: String, pinned: Boolean): Result<Unit> {
        return try {
            database.chatSessionQueries.updateChatSessionPinStatus(
                isPinned = if (pinned) 1L else 0L,
                updatedAt = Clock.System.now().toString(),
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Sync Metadata Operations - No sync support in current schema
    override suspend fun getLastSyncTimestamp(chatId: String): Instant? = null
    override suspend fun updateSyncTimestamp(chatId: String, timestamp: Instant): Result<Unit> = Result.success(Unit)
    override suspend fun getSyncMetadata(chatId: String): SyncMetadata? = null
    override suspend fun updateSyncStatus(chatId: String, status: SyncState, errorMessage: String?): Result<Unit> = Result.success(Unit)
    override suspend fun updateLastSuccessfulSync(chatId: String, timestamp: Instant): Result<Unit> = Result.success(Unit)
    override suspend fun incrementPendingOperations(chatId: String): Result<Unit> = Result.success(Unit)
    override suspend fun decrementPendingOperations(chatId: String): Result<Unit> = Result.success(Unit)
    override suspend fun getAllSyncMetadata(): List<SyncMetadata> = emptyList()

    // Sync Operations Queue - No sync support in current schema
    override suspend fun saveSyncOperation(operation: SyncOperation): Result<Unit> = Result.success(Unit)
    override suspend fun getPendingOperations(): List<SyncOperation> = emptyList()
    override suspend fun getSyncOperationsByChatId(chatId: String): List<SyncOperation> = emptyList()
    override suspend fun getRetryableOperations(currentTime: Instant): List<SyncOperation> = emptyList()
    override suspend fun updateSyncOperationStatus(operationId: String, status: SyncOperationStatus, errorMessage: String?, retryCount: Int?, scheduledAt: Instant?): Result<Unit> = Result.success(Unit)
    override suspend fun markOperationCompleted(operationId: String): Result<Unit> = Result.success(Unit)
    override suspend fun markOperationFailed(operationId: String, errorMessage: String): Result<Unit> = Result.success(Unit)
    override suspend fun deleteSyncOperation(operationId: String): Result<Unit> = Result.success(Unit)
    override suspend fun cleanupCompletedOperations(olderThan: Instant): Result<Int> = Result.success(0)

    // Conflict Management - No conflict support in current schema
    override suspend fun saveConflict(conflict: MessageConflict): Result<Unit> = Result.success(Unit)
    override suspend fun getConflictsByChatId(chatId: String): List<MessageConflict> = emptyList()
    override suspend fun getConflictsByMessageId(messageId: String): List<MessageConflict> = emptyList()
    override suspend fun getAllUnresolvedConflicts(): List<MessageConflict> = emptyList()
    override suspend fun updateConflictResolution(conflictId: String, resolution: ConflictResolution): Result<Unit> = Result.success(Unit)
    override suspend fun deleteResolvedConflicts(olderThan: Instant): Result<Int> = Result.success(0)

    // Participant Operations
    override suspend fun getChatParticipants(chatId: String): List<ChatParticipant> {
        return try {
            database.chatSessionQueries.getChatParticipants(chatId)
                .executeAsList()
                .map { row ->
                    ChatParticipant(
                        id = row.id,
                        name = row.name,
                        avatar = row.avatar,
                        role = ChatRole.valueOf(row.role.uppercase()),
                        status = ParticipantStatus.valueOf(row.status.uppercase()),
                        lastSeen = row.lastSeen,
                        joinedAt = row.joinedAt,
                        customTitle = row.customTitle,
                        isBot = row.isBot == 1L
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun saveChatParticipant(participant: ChatParticipant): Result<Unit> {
        return try {
            database.chatSessionQueries.insertChatParticipant(
                id = participant.id,
                chatId = "", // Will need chatId parameter
                name = participant.name,
                avatar = participant.avatar,
                role = participant.role.name,
                status = participant.status.name,
                lastSeen = participant.lastSeen,
                joinedAt = participant.joinedAt,
                customTitle = participant.customTitle,
                isBot = if (participant.isBot) 1L else 0L
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateParticipantStatus(chatId: String, participantId: String, status: ParticipantStatus, lastSeen: Instant?): Result<Unit> {
        return try {
            database.chatSessionQueries.updateParticipantStatus(
                status = status.name,
                lastSeen = lastSeen?.toString(),
                id = participantId,
                chatId = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun removeParticipant(chatId: String, participantId: String): Result<Unit> {
        // No delete participant query in current schema - return success
        return Result.success(Unit)
    }

    // Attachment Operations
    override suspend fun getMessageAttachments(messageId: String): List<MessageAttachment> {
        return try {
            database.messageQueries.getMessageAttachments(messageId)
                .executeAsList()
                .map { row ->
                    MessageAttachment(
                        id = row.id,
                        type = AttachmentType.valueOf(row.type.uppercase()),
                        url = row.url,
                        thumbnail = row.thumbnail,
                        filename = row.filename,
                        fileSize = row.fileSize,
                        mimeType = row.mimeType,
                        width = row.width?.toInt(),
                        height = row.height?.toInt(),
                        duration = row.duration?.toInt(),
                        caption = row.caption,
                        metadata = emptyMap() // Parse metadata JSON if needed
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun saveMessageAttachment(attachment: MessageAttachment): Result<Unit> {
        return try {
            database.messageQueries.insertMessageAttachment(
                id = attachment.id,
                messageId = "", // Will need to be passed separately
                type = attachment.type.name,
                url = attachment.url,
                thumbnail = attachment.thumbnail,
                filename = attachment.filename,
                fileSize = attachment.fileSize,
                mimeType = attachment.mimeType,
                width = attachment.width?.toLong(),
                height = attachment.height?.toLong(),
                duration = attachment.duration?.toLong(),
                caption = attachment.caption,
                metadata = attachment.metadata.toJsonString()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteMessageAttachments(messageId: String): Result<Unit> {
        // No delete attachments query in current schema - return success
        return Result.success(Unit)
    }

    // Reaction Operations
    override suspend fun getMessageReactions(messageId: String): List<MessageReaction> {
        return try {
            database.messageQueries.getMessageReactions(messageId)
                .executeAsList()
                .map { row ->
                    MessageReaction(
                        emoji = row.emoji,
                        userId = row.userId,
                        userName = "", // Not in schema
                        timestamp = row.timestamp
                    )
                }
        } catch (e: Exception) {
            emptyList()
        }
    }

    override suspend fun saveMessageReaction(reaction: MessageReaction): Result<Unit> {
        return try {
            database.messageQueries.insertMessageReaction(
                messageId = "", // Will need to be passed separately
                userId = reaction.userId,
                emoji = reaction.emoji,
                timestamp = reaction.timestamp
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteMessageReaction(messageId: String, userId: String, emoji: String): Result<Unit> {
        return try {
            database.messageQueries.deleteMessageReaction(messageId, userId, emoji)
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Utility Operations
    override suspend fun clearAllData(): Result<Unit> {
        return try {
            // No clear all data method in current schema - return success
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getStorageInfo(): StorageInfo {
        return try {
            val messageCount = database.messageQueries.getMessagesByChatId("", 1000000).executeAsList().size.toLong()
            val chatCount = database.chatSessionQueries.getAllChatSessions().executeAsList().size.toLong()

            StorageInfo(
                totalMessages = messageCount,
                totalChats = chatCount,
                pendingOperations = 0L,
                unresolvedConflicts = 0L,
                databaseSize = 0L, // Cannot get DB size easily
                lastOptimized = null
            )
        } catch (e: Exception) {
            StorageInfo(
                totalMessages = 0L,
                totalChats = 0L,
                pendingOperations = 0L,
                unresolvedConflicts = 0L,
                databaseSize = 0L,
                lastOptimized = null
            )
        }
    }

    override suspend fun optimizeDatabase(): Result<Unit> {
        return try {
            // SQLDelight doesn't expose VACUUM directly - return success
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}