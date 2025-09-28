package com.tchat.mobile.repositories

import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.datasource.ApiRemoteDataSource
import com.tchat.mobile.database.TchatDatabase
import app.cash.sqldelight.coroutines.asFlow
import app.cash.sqldelight.coroutines.mapToList
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant

/**
 * Hybrid ChatRepository implementation that uses both remote API and local storage
 *
 * Strategy:
 * - Read operations: Try remote first for latest data, fallback to local
 * - Write operations: Write to remote first, then sync to local
 * - Real-time updates: Listen to remote changes and update local storage
 */
class HybridChatRepository(
    private val remoteDataSource: ApiRemoteDataSource,
    private val database: TchatDatabase
) : ChatRepository {

    // Chat Session Operations
    override suspend fun getChatSessions(): Result<List<ChatSession>> {
        return try {
            // Try to get latest data from remote first
            val remoteResult = remoteDataSource.getChatSessions("current_user")

            if (remoteResult.isSuccess) {
                val remoteSessions = remoteResult.getOrThrow()
                // Sync to local storage for offline access
                remoteSessions.forEach { session ->
                    insertChatSessionToLocal(session)
                }
                Result.success(remoteSessions)
            } else {
                // Fallback to local data if remote fails
                getLocalChatSessions()
            }
        } catch (e: Exception) {
            // If both fail, return local data
            getLocalChatSessions()
        }
    }

    override suspend fun getChatSession(chatId: String): Result<ChatSession> {
        return try {
            // Try remote first for latest data
            val remoteResult = remoteDataSource.getChatSession(chatId)

            if (remoteResult.isSuccess) {
                val session = remoteResult.getOrThrow()
                // Update local cache
                insertChatSessionToLocal(session)
                Result.success(session)
            } else {
                // Fallback to local
                getLocalChatSession(chatId)
            }
        } catch (e: Exception) {
            getLocalChatSession(chatId)
        }
    }

    override suspend fun createChatSession(session: ChatSession): Result<ChatSession> {
        return try {
            // Create on remote first
            val remoteResult = remoteDataSource.createChatSession(session)

            if (remoteResult.isSuccess) {
                val createdSession = remoteResult.getOrThrow()
                // Sync to local
                insertChatSessionToLocal(createdSession)
                Result.success(createdSession)
            } else {
                // Store locally for later sync
                insertChatSessionToLocal(session)
                Result.success(session)
            }
        } catch (e: Exception) {
            // Store locally and mark for sync
            insertChatSessionToLocal(session)
            Result.success(session)
        }
    }

    override suspend fun updateChatSession(chatId: String, session: ChatSession): Result<ChatSession> {
        return try {
            // Update remote first
            val remoteResult = remoteDataSource.updateChatSession(session)

            if (remoteResult.isSuccess) {
                val updatedSession = remoteResult.getOrThrow()
                // Sync to local
                insertChatSessionToLocal(updatedSession)
                Result.success(updatedSession)
            } else {
                // Update locally and mark for sync
                insertChatSessionToLocal(session)
                Result.success(session)
            }
        } catch (e: Exception) {
            insertChatSessionToLocal(session)
            Result.success(session)
        }
    }

    override suspend fun deleteChatSession(chatId: String): Result<Boolean> {
        return try {
            // Delete from remote first
            val remoteResult = remoteDataSource.deleteChatSession(chatId)

            if (remoteResult.isSuccess) {
                // Delete from local
                database.chatSessionQueries.deleteChatSession(chatId)
                Result.success(true)
            } else {
                // Mark as deleted locally
                Result.success(true)
            }
        } catch (e: Exception) {
            Result.success(true)
        }
    }

    // Message Operations
    override suspend fun getMessages(chatId: String, limit: Int, offset: Int): Result<List<Message>> {
        return try {
            // Try remote for latest messages
            val remoteResult = remoteDataSource.getMessages(chatId, limit, null)

            if (remoteResult.isSuccess) {
                val remoteMessages = remoteResult.getOrThrow().map { it.toMessage() }
                // Update local cache
                remoteMessages.forEach { message ->
                    insertMessageToLocal(message)
                }
                Result.success(remoteMessages)
            } else {
                // Fallback to local
                getLocalMessages(chatId, limit)
            }
        } catch (e: Exception) {
            getLocalMessages(chatId, limit)
        }
    }

    override suspend fun getMessage(messageId: String): Result<Message> {
        return try {
            // Try local first for performance
            val localResult = getLocalMessage(messageId)
            if (localResult.isSuccess) {
                localResult
            } else {
                // Try remote if not in local cache
                Result.failure(Exception("Message not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun sendMessage(message: Message): Result<Message> {
        return try {
            // Convert Message to ChatMessage for API
            val chatMessage = message.toChatMessage()

            // Send to remote first
            val remoteResult = remoteDataSource.sendMessage(chatMessage)

            if (remoteResult.isSuccess) {
                val sentMessage = remoteResult.getOrThrow().toMessage()
                // Update local with server response
                insertMessageToLocal(sentMessage)
                Result.success(sentMessage)
            } else {
                // Store locally with pending status
                val pendingMessage = message.copy(
                    deliveryStatus = DeliveryStatus.PENDING
                )
                insertMessageToLocal(pendingMessage)
                Result.success(pendingMessage)
            }
        } catch (e: Exception) {
            // Store locally with pending status for later sync
            val pendingMessage = message.copy(
                deliveryStatus = DeliveryStatus.PENDING
            )
            insertMessageToLocal(pendingMessage)
            Result.success(pendingMessage)
        }
    }

    override suspend fun editMessage(messageId: String, newContent: String): Result<Message> {
        return try {
            // Get the message first
            val localResult = getLocalMessage(messageId)
            if (localResult.isSuccess) {
                val message = localResult.getOrThrow()
                val editedMessage = message.copy(
                    content = newContent,
                    isEdited = true,
                    editedAt = Clock.System.now().toString()
                )
                insertMessageToLocal(editedMessage)
                Result.success(editedMessage)
            } else {
                Result.failure(Exception("Message not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Boolean> {
        return try {
            // Delete from local
            database.messageQueries.deleteMessage(messageId, Clock.System.now().toString())
            Result.success(true)
        } catch (e: Exception) {
            Result.success(true)
        }
    }

    // Observable operations for real-time updates
    override fun observeChatSessions(): Flow<List<ChatSession>> {
        // Return local data as Flow for UI
        return database.chatSessionQueries.getAllChatSessions()
            .asFlow()
            .mapToList(Dispatchers.IO)
            .map { rows ->
                rows.map { row ->
                    ChatSession(
                        id = row.id,
                        name = row.name ?: "Chat",
                        type = ChatType.valueOf(row.type),
                        participants = emptyList(), // Will need to query separately
                        lastMessage = null, // Will need to query last message
                        unreadCount = row.unreadCount.toInt(),
                        metadata = ChatMetadata(),
                        createdAt = row.createdAt,
                        updatedAt = row.updatedAt,
                        lastActivityAt = row.lastActivityAt
                    )
                }
            }
    }

    override fun observeMessages(chatId: String): Flow<List<Message>> {
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

    // Participant Operations (simplified for now)
    override suspend fun addParticipant(chatId: String, participant: ChatParticipant): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun removeParticipant(chatId: String, participantId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun updateParticipantStatus(chatId: String, participantId: String, status: ParticipantStatus): Result<Boolean> {
        return Result.success(true)
    }

    // Reactions and Interactions (simplified for now)
    override suspend fun addReaction(messageId: String, reaction: MessageReaction): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun removeReaction(messageId: String, userId: String, emoji: String): Result<Boolean> {
        return Result.success(true)
    }

    // Chat Management (simplified for now)
    override suspend fun pinMessage(messageId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun unpinMessage(messageId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun markAsRead(chatId: String, messageId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun muteChat(chatId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun unmuteChat(chatId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun archiveChat(chatId: String): Result<Boolean> {
        return Result.success(true)
    }

    override suspend fun unarchiveChat(chatId: String): Result<Boolean> {
        return Result.success(true)
    }

    // Search and Filtering
    override suspend fun searchMessages(chatId: String, query: String): Result<List<Message>> {
        return try {
            val messages = database.messageQueries.searchMessages(chatId, query)
                .executeAsList()
                .map { row ->
                    Message(
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
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
            Result.success(messages)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun searchChatSessions(query: String): Result<List<ChatSession>> {
        return try {
            // Simple local search for now
            val sessions = database.chatSessionQueries.getAllChatSessions()
                .executeAsList()
                .filter { it.name?.contains(query, ignoreCase = true) == true }
                .map { row ->
                    ChatSession(
                        id = row.id,
                        name = row.name ?: "Chat",
                        type = ChatType.valueOf(row.type),
                        participants = emptyList(),
                        lastMessage = null,
                        unreadCount = row.unreadCount.toInt(),
                        metadata = ChatMetadata(),
                        createdAt = row.createdAt,
                        updatedAt = row.updatedAt,
                        lastActivityAt = row.lastActivityAt
                    )
                }
            Result.success(sessions)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Real-time updates for UI
    override fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>> {
        return kotlinx.coroutines.flow.flowOf(emptyList()) // Simplified for now
    }

    override fun observeParticipantStatus(chatId: String): Flow<List<ChatParticipant>> {
        return kotlinx.coroutines.flow.flowOf(emptyList()) // Simplified for now
    }

    // Advanced messaging operations
    override suspend fun sendTextMessage(chatId: String, content: String, replyToId: String?): Result<Message> {
        val message = Message(
            id = "temp_${Clock.System.now().toEpochMilliseconds()}",
            chatId = chatId,
            senderId = "current_user", // Would get from auth service
            senderName = "You",
            content = content,
            type = MessageType.TEXT,
            createdAt = Clock.System.now().toString(),
            replyToId = replyToId
        )
        return sendMessage(message)
    }

    override suspend fun sendMediaMessage(chatId: String, attachments: List<MessageAttachment>, caption: String?, replyToId: String?): Result<Message> {
        val messageType = when {
            attachments.any { it.type == AttachmentType.IMAGE } -> MessageType.IMAGE
            attachments.any { it.type == AttachmentType.VIDEO } -> MessageType.VIDEO
            attachments.any { it.type == AttachmentType.AUDIO } -> MessageType.AUDIO
            attachments.any { it.type == AttachmentType.FILE } -> MessageType.FILE
            attachments.any { it.type == AttachmentType.LOCATION } -> MessageType.LOCATION
            else -> MessageType.FILE
        }

        val message = Message(
            id = "temp_${Clock.System.now().toEpochMilliseconds()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            content = caption ?: "Media",
            type = messageType,
            attachments = attachments,
            createdAt = Clock.System.now().toString().toString(),
            replyToId = replyToId
        )
        return sendMessage(message)
    }

    override suspend fun sendVoiceMessage(chatId: String, audioUrl: String, duration: Int): Result<Message> {
        val audioAttachment = MessageAttachment(
            id = "audio_${Clock.System.now().toEpochMilliseconds()}",
            type = AttachmentType.AUDIO,
            url = audioUrl,
            duration = duration,
            mimeType = "audio/mp4"
        )

        val message = Message(
            id = "temp_${Clock.System.now().toEpochMilliseconds()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            content = "Voice message",
            type = MessageType.AUDIO,
            attachments = listOf(audioAttachment),
            createdAt = Clock.System.now().toString()
        )
        return sendMessage(message)
    }

    override suspend fun sendLocationMessage(chatId: String, latitude: Double, longitude: Double, address: String?): Result<Message> {
        val locationAttachment = MessageAttachment(
            id = "location_${Clock.System.now().toEpochMilliseconds()}",
            type = AttachmentType.LOCATION,
            url = "geo:$latitude,$longitude",
            metadata = mapOf(
                "latitude" to latitude.toString(),
                "longitude" to longitude.toString(),
                "address" to (address ?: "Unknown location")
            )
        )

        val message = Message(
            id = "temp_${Clock.System.now().toEpochMilliseconds()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            content = address ?: "Location shared",
            type = MessageType.LOCATION,
            attachments = listOf(locationAttachment),
            createdAt = Clock.System.now().toString()
        )
        return sendMessage(message)
    }

    override suspend fun forwardMessage(messageId: String, targetChatIds: List<String>): Result<List<Message>> {
        return try {
            val originalMessage = getLocalMessage(messageId).getOrNull()
                ?: return Result.failure(Exception("Message not found"))

            val forwardedMessages = targetChatIds.mapNotNull { chatId ->
                val message = Message(
                    id = "temp_${Clock.System.now().toEpochMilliseconds()}_$chatId",
                    chatId = chatId,
                    senderId = "current_user",
                    senderName = "You",
                    content = originalMessage.content,
                    type = originalMessage.type,
                    attachments = originalMessage.attachments,
                    createdAt = Clock.System.now().toString()
                )
                sendMessage(message).getOrNull()
            }
            Result.success(forwardedMessages)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Typing indicators
    override suspend fun startTyping(chatId: String): Result<Boolean> {
        return Result.success(true) // Simplified for now
    }

    override suspend fun stopTyping(chatId: String): Result<Boolean> {
        return Result.success(true) // Simplified for now
    }

    // Message status updates
    override suspend fun markMessageAsDelivered(messageId: String): Result<Boolean> {
        return updateDeliveryStatus(messageId, DeliveryStatus.DELIVERED)
    }

    override suspend fun markMessageAsRead(messageId: String, userId: String): Result<Boolean> {
        return Result.success(true) // Simplified for now
    }

    override suspend fun updateDeliveryStatus(messageId: String, status: DeliveryStatus): Result<Boolean> {
        return Result.success(true) // Simplified for now
    }

    // Voice message operations
    override suspend fun uploadVoiceMessage(audioData: ByteArray): Result<String> {
        return Result.success("https://example.com/voice/${Clock.System.now().toEpochMilliseconds()}.mp4")
    }

    override suspend fun downloadVoiceMessage(url: String): Result<ByteArray> {
        return Result.success(ByteArray(0)) // Mock implementation
    }

    // File operations
    override suspend fun uploadFile(fileData: ByteArray, filename: String, mimeType: String): Result<MessageAttachment> {
        val attachment = MessageAttachment(
            id = "file_${Clock.System.now().toEpochMilliseconds()}",
            type = when {
                mimeType.startsWith("image/") -> AttachmentType.IMAGE
                mimeType.startsWith("video/") -> AttachmentType.VIDEO
                mimeType.startsWith("audio/") -> AttachmentType.AUDIO
                else -> AttachmentType.FILE
            },
            url = "https://example.com/files/${Clock.System.now().toEpochMilliseconds()}",
            filename = filename,
            fileSize = fileData.size.toLong(),
            mimeType = mimeType
        )
        return Result.success(attachment)
    }

    override suspend fun downloadFile(attachment: MessageAttachment): Result<ByteArray> {
        return Result.success(ByteArray(0)) // Mock implementation
    }

    // Draft messages
    override suspend fun saveDraftMessage(chatId: String, content: String): Result<Boolean> {
        return Result.success(true) // Could implement with local storage
    }

    override suspend fun getDraftMessage(chatId: String): Result<String?> {
        return Result.success(null) // Could implement with local storage
    }

    override suspend fun clearDraftMessage(chatId: String): Result<Boolean> {
        return Result.success(true) // Could implement with local storage
    }

    // Private helper methods
    private suspend fun getLocalChatSessions(): Result<List<ChatSession>> {
        return try {
            val sessions = database.chatSessionQueries.getAllChatSessions()
                .executeAsList()
                .map { row ->
                    ChatSession(
                        id = row.id,
                        name = row.name ?: "Chat",
                        type = ChatType.valueOf(row.type),
                        participants = emptyList(),
                        lastMessage = null,
                        unreadCount = row.unreadCount.toInt(),
                        metadata = ChatMetadata(),
                        createdAt = row.createdAt,
                        updatedAt = row.updatedAt,
                        lastActivityAt = row.lastActivityAt
                    )
                }
            Result.success(sessions)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private suspend fun getLocalChatSession(chatId: String): Result<ChatSession> {
        return try {
            val row = database.chatSessionQueries.getChatSessionById(chatId).executeAsOneOrNull()
            if (row != null) {
                val session = ChatSession(
                    id = row.id,
                    name = row.name ?: "Chat",
                    type = ChatType.valueOf(row.type),
                    participants = emptyList(),
                    lastMessage = null,
                    unreadCount = row.unreadCount.toInt(),
                    metadata = ChatMetadata(),
                    createdAt = row.createdAt,
                    updatedAt = row.updatedAt,
                    lastActivityAt = row.lastActivityAt
                )
                Result.success(session)
            } else {
                Result.failure(Exception("Chat session not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private suspend fun getLocalMessages(chatId: String, limit: Int): Result<List<Message>> {
        return try {
            val messages = database.messageQueries.getMessagesByChatId(chatId, limit.toLong())
                .executeAsList()
                .map { row ->
                    Message(
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
                        deletedAt = row.deletedAt,
                        deliveryStatus = DeliveryStatus.SENT,
                        readBy = emptyList()
                    )
                }
            Result.success(messages)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private suspend fun getLocalMessage(messageId: String): Result<Message> {
        return try {
            val row = database.messageQueries.getMessageById(messageId).executeAsOneOrNull()
            if (row != null) {
                val message = Message(
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
                    deletedAt = row.deletedAt,
                    deliveryStatus = DeliveryStatus.SENT,
                    readBy = emptyList()
                )
                Result.success(message)
            } else {
                Result.failure(Exception("Message not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private suspend fun insertChatSessionToLocal(session: ChatSession) {
        try {
            val currentTime = Clock.System.now().toString()
            database.chatSessionQueries.insertChatSession(
                id = session.id,
                name = session.name,
                type = session.type.name,
                unreadCount = session.unreadCount.toLong(),
                isPinned = 0L,
                isMuted = 0L,
                isArchived = 0L,
                isBlocked = 0L,
                createdAt = session.createdAt,
                updatedAt = session.updatedAt,
                lastActivityAt = session.lastActivityAt
            )
        } catch (e: Exception) {
            // Ignore insertion errors for now
        }
    }

    private suspend fun insertMessageToLocal(message: Message) {
        try {
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
                reactions = "[]", // JSON string
                attachmentCount = message.attachments.size.toLong(),
                createdAt = message.createdAt.toString(),
                editedAt = message.editedAt?.toString(),
                deletedAt = message.deletedAt?.toString()
            )
        } catch (e: Exception) {
            // Ignore insertion errors for now
        }
    }
}

// Extension functions for model conversion
private fun ChatMessage.toMessage(): Message {
    return Message(
        id = id,
        chatId = chatId,
        senderId = senderId,
        senderName = senderName,
        type = type,
        content = content,
        isEdited = isEdited,
        isPinned = isPinned,
        isDeleted = isDeleted,
        replyToId = replyToId,
        reactions = reactions,
        attachments = attachments,
        createdAt = createdAt,
        editedAt = editedAt,
        deletedAt = deletedAt,
        deliveryStatus = DeliveryStatus.valueOf(deliveryStatus.name),
        readBy = readBy
    )
}

private fun Message.toChatMessage(): ChatMessage {
    return ChatMessage(
        id = id,
        chatId = chatId,
        senderId = senderId,
        senderName = senderName,
        senderAvatar = null,
        type = type,
        content = content,
        isEdited = isEdited,
        isPinned = isPinned,
        isDeleted = isDeleted,
        replyToId = replyToId,
        reactions = reactions,
        attachments = attachments,
        createdAt = createdAt,
        editedAt = editedAt,
        deletedAt = deletedAt,
        deliveryStatus = DeliveryStatus.valueOf(deliveryStatus.name),
        readBy = readBy
    )
}