package com.tchat.mobile.repositories

import app.cash.sqldelight.coroutines.asFlow
import app.cash.sqldelight.coroutines.mapToList
import app.cash.sqldelight.coroutines.mapToOne
import app.cash.sqldelight.coroutines.mapToOneOrNull
import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.withContext
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

/**
 * SQLDelight implementation of ChatRepository
 *
 * Provides persistent storage for chat sessions, messages, and related data
 * using SQLDelight for cross-platform database operations.
 */
class SQLDelightChatRepository(
    private val database: TchatDatabase
) : ChatRepository {

    private val json = Json { ignoreUnknownKeys = true }

    override suspend fun getChatSessions(): Result<List<ChatSession>> = withContext(Dispatchers.Default) {
        try {
            val sessions = database.chatSessionQueries.getAllChatSessions().executeAsList()
            val chatSessions = sessions.map { session ->
                val participants = database.chatSessionQueries.getChatParticipants(session.id).executeAsList()
                val metadata = database.chatSessionQueries.getChatMetadata(session.id).executeAsOneOrNull()

                ChatSession(
                    id = session.id,
                    name = session.name,
                    type = ChatType.valueOf(session.type),
                    participants = participants.map { participant ->
                        ChatParticipant(
                            id = participant.id,
                            name = participant.name,
                            avatar = participant.avatar,
                            role = ChatRole.valueOf(participant.role),
                            status = ParticipantStatus.valueOf(participant.status),
                            lastSeen = participant.lastSeen,
                            joinedAt = participant.joinedAt,
                            customTitle = participant.customTitle,
                            isBot = participant.isBot == 1L
                        )
                    },
                    lastMessage = getLastMessagePreview(session.id),
                    unreadCount = session.unreadCount.toInt(),
                    isPinned = session.isPinned == 1L,
                    isMuted = session.isMuted == 1L,
                    isArchived = session.isArchived == 1L,
                    isBlocked = session.isBlocked == 1L,
                    metadata = metadata?.let { meta ->
                        ChatMetadata(
                            description = meta.description,
                            avatar = meta.avatar,
                            backgroundColor = meta.backgroundColor,
                            theme = meta.theme,
                            language = meta.language,
                            timezone = meta.timezone,
                            isPublic = meta.isPublic == 1L,
                            inviteLink = meta.inviteLink,
                            maxParticipants = meta.maxParticipants?.toInt(),
                            autoDeleteMessages = meta.autoDeleteMessages == 1L,
                            autoDeleteDays = meta.autoDeleteDays?.toInt(),
                            encryptionEnabled = meta.encryptionEnabled == 1L,
                            backupEnabled = meta.backupEnabled == 1L
                        )
                    } ?: ChatMetadata(),
                    createdAt = session.createdAt,
                    updatedAt = session.updatedAt,
                    lastActivityAt = session.lastActivityAt
                )
            }
            Result.success(chatSessions)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getChatSession(chatId: String): Result<ChatSession> = withContext(Dispatchers.Default) {
        try {
            val session = database.chatSessionQueries.getChatSessionById(chatId).executeAsOneOrNull()
                ?: return@withContext Result.failure(Exception("Chat session not found"))

            val participants = database.chatSessionQueries.getChatParticipants(chatId).executeAsList()
            val metadata = database.chatSessionQueries.getChatMetadata(chatId).executeAsOneOrNull()

            val chatSession = ChatSession(
                id = session.id,
                name = session.name,
                type = ChatType.valueOf(session.type),
                participants = participants.map { participant ->
                    ChatParticipant(
                        id = participant.id,
                        name = participant.name,
                        avatar = participant.avatar,
                        role = ChatRole.valueOf(participant.role),
                        status = ParticipantStatus.valueOf(participant.status),
                        lastSeen = participant.lastSeen,
                        joinedAt = participant.joinedAt,
                        customTitle = participant.customTitle,
                        isBot = participant.isBot == 1L
                    )
                },
                lastMessage = getLastMessagePreview(chatId),
                unreadCount = session.unreadCount.toInt(),
                isPinned = session.isPinned == 1L,
                isMuted = session.isMuted == 1L,
                isArchived = session.isArchived == 1L,
                isBlocked = session.isBlocked == 1L,
                metadata = metadata?.let { meta ->
                    ChatMetadata(
                        description = meta.description,
                        avatar = meta.avatar,
                        backgroundColor = meta.backgroundColor,
                        theme = meta.theme,
                        language = meta.language,
                        timezone = meta.timezone,
                        isPublic = meta.isPublic == 1L,
                        inviteLink = meta.inviteLink,
                        maxParticipants = meta.maxParticipants?.toInt(),
                        autoDeleteMessages = meta.autoDeleteMessages == 1L,
                        autoDeleteDays = meta.autoDeleteDays?.toInt(),
                        encryptionEnabled = meta.encryptionEnabled == 1L,
                        backupEnabled = meta.backupEnabled == 1L
                    )
                } ?: ChatMetadata(),
                createdAt = session.createdAt,
                updatedAt = session.updatedAt,
                lastActivityAt = session.lastActivityAt
            )
            Result.success(chatSession)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun createChatSession(session: ChatSession): Result<ChatSession> = withContext(Dispatchers.Default) {
        try {
            database.transaction {
                // Insert chat session
                database.chatSessionQueries.insertChatSession(
                    id = session.id,
                    name = session.name,
                    type = session.type.name,
                    unreadCount = session.unreadCount.toLong(),
                    isPinned = if (session.isPinned) 1L else 0L,
                    isMuted = if (session.isMuted) 1L else 0L,
                    isArchived = if (session.isArchived) 1L else 0L,
                    isBlocked = if (session.isBlocked) 1L else 0L,
                    createdAt = session.createdAt,
                    updatedAt = session.updatedAt,
                    lastActivityAt = session.lastActivityAt
                )

                // Insert participants
                session.participants.forEach { participant ->
                    database.chatSessionQueries.insertChatParticipant(
                        id = participant.id,
                        chatId = session.id,
                        name = participant.name,
                        avatar = participant.avatar,
                        role = participant.role.name,
                        status = participant.status.name,
                        lastSeen = participant.lastSeen,
                        joinedAt = participant.joinedAt,
                        customTitle = participant.customTitle,
                        isBot = if (participant.isBot) 1L else 0L
                    )
                }

                // Insert metadata
                database.chatSessionQueries.insertChatMetadata(
                    chatId = session.id,
                    description = session.metadata.description,
                    avatar = session.metadata.avatar,
                    backgroundColor = session.metadata.backgroundColor,
                    theme = session.metadata.theme,
                    language = session.metadata.language,
                    timezone = session.metadata.timezone,
                    isPublic = if (session.metadata.isPublic) 1L else 0L,
                    inviteLink = session.metadata.inviteLink,
                    maxParticipants = session.metadata.maxParticipants?.toLong(),
                    autoDeleteMessages = if (session.metadata.autoDeleteMessages) 1L else 0L,
                    autoDeleteDays = session.metadata.autoDeleteDays?.toLong(),
                    encryptionEnabled = if (session.metadata.encryptionEnabled) 1L else 0L,
                    backupEnabled = if (session.metadata.backupEnabled) 1L else 0L
                )
            }
            Result.success(session)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateChatSession(chatId: String, session: ChatSession): Result<ChatSession> = withContext(Dispatchers.Default) {
        try {
            // For simplicity, we'll update basic properties
            database.chatSessionQueries.updateChatSessionUnreadCount(
                unreadCount = session.unreadCount.toLong(),
                updatedAt = session.updatedAt,
                lastActivityAt = session.lastActivityAt,
                id = chatId
            )
            Result.success(session)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteChatSession(chatId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.deleteChatSession(chatId)
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getMessages(chatId: String, limit: Int, offset: Int): Result<List<Message>> = withContext(Dispatchers.Default) {
        try {
            val messages = database.messageQueries.getMessagesByChatId(chatId, limit.toLong()).executeAsList()
            val messageModels = messages.map { message ->
                val attachments = database.messageQueries.getMessageAttachments(message.id).executeAsList()
                val reactions = database.messageQueries.getMessageReactions(message.id).executeAsList()

                Message(
                    id = message.id,
                    chatId = message.chatId,
                    senderId = message.senderId,
                    senderName = message.senderName,
                    type = MessageType.valueOf(message.type),
                    content = message.content,
                    isEdited = message.isEdited == 1L,
                    isPinned = message.isPinned == 1L,
                    isDeleted = message.isDeleted == 1L,
                    replyToId = message.replyToId,
                    reactions = reactions.map { reaction ->
                        MessageReaction(
                            emoji = reaction.emoji,
                            userId = reaction.userId,
                            userName = "", // We'd need to join with user table for this
                            timestamp = reaction.timestamp
                        )
                    },
                    attachments = attachments.map { attachment ->
                        MessageAttachment(
                            id = attachment.id,
                            type = AttachmentType.valueOf(attachment.type),
                            url = attachment.url,
                            thumbnail = attachment.thumbnail,
                            filename = attachment.filename,
                            fileSize = attachment.fileSize,
                            mimeType = attachment.mimeType,
                            width = attachment.width?.toInt(),
                            height = attachment.height?.toInt(),
                            duration = attachment.duration?.toInt(),
                            caption = attachment.caption,
                            metadata = attachment.metadata?.let {
                                try {
                                    json.decodeFromString<Map<String, String>>(it)
                                } catch (e: Exception) {
                                    emptyMap()
                                }
                            } ?: emptyMap()
                        )
                    },
                    createdAt = message.createdAt,
                    editedAt = message.editedAt,
                    deletedAt = message.deletedAt
                )
            }
            Result.success(messageModels)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getMessage(messageId: String): Result<Message> = withContext(Dispatchers.Default) {
        try {
            val message = database.messageQueries.getMessageById(messageId).executeAsOneOrNull()
                ?: return@withContext Result.failure(Exception("Message not found"))

            val attachments = database.messageQueries.getMessageAttachments(messageId).executeAsList()
            val reactions = database.messageQueries.getMessageReactions(messageId).executeAsList()

            val messageModel = Message(
                id = message.id,
                chatId = message.chatId,
                senderId = message.senderId,
                senderName = message.senderName,
                type = MessageType.valueOf(message.type),
                content = message.content,
                isEdited = message.isEdited == 1L,
                isPinned = message.isPinned == 1L,
                isDeleted = message.isDeleted == 1L,
                replyToId = message.replyToId,
                reactions = reactions.map { reaction ->
                    MessageReaction(
                        emoji = reaction.emoji,
                        userId = reaction.userId,
                        userName = "",
                        timestamp = reaction.timestamp
                    )
                },
                attachments = attachments.map { attachment ->
                    MessageAttachment(
                        id = attachment.id,
                        type = AttachmentType.valueOf(attachment.type),
                        url = attachment.url,
                        thumbnail = attachment.thumbnail,
                        filename = attachment.filename,
                        fileSize = attachment.fileSize,
                        mimeType = attachment.mimeType,
                        width = attachment.width?.toInt(),
                        height = attachment.height?.toInt(),
                        duration = attachment.duration?.toInt(),
                        caption = attachment.caption
                    )
                },
                createdAt = message.createdAt,
                editedAt = message.editedAt,
                deletedAt = message.deletedAt
            )
            Result.success(messageModel)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun sendMessage(message: Message): Result<Message> = withContext(Dispatchers.Default) {
        try {
            database.transaction {
                // Insert message
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
                    reactions = json.encodeToString(message.reactions),
                    attachmentCount = message.attachments.size.toLong(),
                    createdAt = message.createdAt,
                    editedAt = message.editedAt,
                    deletedAt = message.deletedAt
                )

                // Insert attachments
                message.attachments.forEach { attachment ->
                    database.messageQueries.insertMessageAttachment(
                        id = attachment.id,
                        messageId = message.id,
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
                        metadata = if (attachment.metadata.isNotEmpty()) json.encodeToString(attachment.metadata) else null
                    )
                }

                // Insert reactions
                message.reactions.forEach { reaction ->
                    database.messageQueries.insertMessageReaction(
                        messageId = message.id,
                        userId = reaction.userId,
                        emoji = reaction.emoji,
                        timestamp = reaction.timestamp
                    )
                }
            }
            Result.success(message)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun editMessage(messageId: String, newContent: String): Result<Message> = withContext(Dispatchers.Default) {
        try {
            val timestamp = System.currentTimeMillis().toString()
            database.messageQueries.updateMessageContent(
                content = newContent,
                editedAt = timestamp,
                id = messageId
            )
            getMessage(messageId)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            val timestamp = System.currentTimeMillis().toString()
            database.messageQueries.deleteMessage(timestamp, messageId)
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Simplified implementations for other methods
    override suspend fun addParticipant(chatId: String, participant: ChatParticipant): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.insertChatParticipant(
                id = participant.id,
                chatId = chatId,
                name = participant.name,
                avatar = participant.avatar,
                role = participant.role.name,
                status = participant.status.name,
                lastSeen = participant.lastSeen,
                joinedAt = participant.joinedAt,
                customTitle = participant.customTitle,
                isBot = if (participant.isBot) 1L else 0L
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun removeParticipant(chatId: String, participantId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        // Would need a delete participant query
        Result.success(true)
    }

    override suspend fun updateParticipantStatus(chatId: String, participantId: String, status: ParticipantStatus): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            val lastSeen = if (status == ParticipantStatus.OFFLINE) System.currentTimeMillis().toString() else null
            database.chatSessionQueries.updateParticipantStatus(
                status = status.name,
                lastSeen = lastSeen,
                id = participantId,
                chatId = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun addReaction(messageId: String, reaction: MessageReaction): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.messageQueries.insertMessageReaction(
                messageId = messageId,
                userId = reaction.userId,
                emoji = reaction.emoji,
                timestamp = reaction.timestamp
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun removeReaction(messageId: String, userId: String, emoji: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.messageQueries.deleteMessageReaction(messageId, userId, emoji)
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun pinMessage(messageId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.messageQueries.updateMessagePinStatus(isPinned = 1L, id = messageId)
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun unpinMessage(messageId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.messageQueries.updateMessagePinStatus(isPinned = 0L, id = messageId)
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun markAsRead(chatId: String, messageId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.updateChatSessionUnreadCount(
                unreadCount = 0L,
                updatedAt = System.currentTimeMillis().toString(),
                lastActivityAt = System.currentTimeMillis().toString(),
                id = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun muteChat(chatId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.updateChatSessionMuteStatus(
                isMuted = 1L,
                updatedAt = System.currentTimeMillis().toString(),
                id = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun unmuteChat(chatId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.updateChatSessionMuteStatus(
                isMuted = 0L,
                updatedAt = System.currentTimeMillis().toString(),
                id = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun archiveChat(chatId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.updateChatSessionArchiveStatus(
                isArchived = 1L,
                updatedAt = System.currentTimeMillis().toString(),
                id = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun unarchiveChat(chatId: String): Result<Boolean> = withContext(Dispatchers.Default) {
        try {
            database.chatSessionQueries.updateChatSessionArchiveStatus(
                isArchived = 0L,
                updatedAt = System.currentTimeMillis().toString(),
                id = chatId
            )
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun searchMessages(chatId: String, query: String): Result<List<Message>> = withContext(Dispatchers.Default) {
        try {
            val messages = database.messageQueries.searchMessages(chatId, query).executeAsList()
            val messageModels = messages.map { message ->
                Message(
                    id = message.id,
                    chatId = message.chatId,
                    senderId = message.senderId,
                    senderName = message.senderName,
                    type = MessageType.valueOf(message.type),
                    content = message.content,
                    isEdited = message.isEdited == 1L,
                    isPinned = message.isPinned == 1L,
                    isDeleted = message.isDeleted == 1L,
                    replyToId = message.replyToId,
                    createdAt = message.createdAt,
                    editedAt = message.editedAt,
                    deletedAt = message.deletedAt
                )
            }
            Result.success(messageModels)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun searchChatSessions(query: String): Result<List<ChatSession>> = withContext(Dispatchers.Default) {
        // For simplicity, get all sessions and filter in memory
        getChatSessions().map { sessions ->
            sessions.filter { session ->
                session.name?.contains(query, ignoreCase = true) == true ||
                session.participants.any { it.name.contains(query, ignoreCase = true) }
            }
        }
    }

    override fun observeChatSessions(): Flow<List<ChatSession>> {
        return database.chatSessionQueries.getAllChatSessions()
            .asFlow()
            .mapToList(Dispatchers.Default)
            .map { sessions ->
                // Convert to ChatSession objects (simplified)
                sessions.map { session ->
                    ChatSession(
                        id = session.id,
                        name = session.name,
                        type = ChatType.valueOf(session.type),
                        participants = emptyList(), // Would need to join
                        unreadCount = session.unreadCount.toInt(),
                        isPinned = session.isPinned == 1L,
                        isMuted = session.isMuted == 1L,
                        isArchived = session.isArchived == 1L,
                        isBlocked = session.isBlocked == 1L,
                        metadata = ChatMetadata(),
                        createdAt = session.createdAt,
                        updatedAt = session.updatedAt,
                        lastActivityAt = session.lastActivityAt
                    )
                }
            }
    }

    override fun observeMessages(chatId: String): Flow<List<Message>> {
        return database.messageQueries.getMessagesByChatId(chatId, 50L)
            .asFlow()
            .mapToList(Dispatchers.Default)
            .map { messages ->
                messages.map { message ->
                    Message(
                        id = message.id,
                        chatId = message.chatId,
                        senderId = message.senderId,
                        senderName = message.senderName,
                        type = MessageType.valueOf(message.type),
                        content = message.content,
                        isEdited = message.isEdited == 1L,
                        isPinned = message.isPinned == 1L,
                        isDeleted = message.isDeleted == 1L,
                        replyToId = message.replyToId,
                        createdAt = message.createdAt,
                        editedAt = message.editedAt,
                        deletedAt = message.deletedAt
                    )
                }
            }
    }

    override fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>> {
        // For now, return empty flow - would need a typing indicators table
        return kotlinx.coroutines.flow.flowOf(emptyList())
    }

    override fun observeParticipantStatus(chatId: String): Flow<List<ChatParticipant>> {
        return database.chatSessionQueries.getChatParticipants(chatId)
            .asFlow()
            .mapToList(Dispatchers.Default)
            .map { participants ->
                participants.map { participant ->
                    ChatParticipant(
                        id = participant.id,
                        name = participant.name,
                        avatar = participant.avatar,
                        role = ChatRole.valueOf(participant.role),
                        status = ParticipantStatus.valueOf(participant.status),
                        lastSeen = participant.lastSeen,
                        joinedAt = participant.joinedAt,
                        customTitle = participant.customTitle,
                        isBot = participant.isBot == 1L
                    )
                }
            }
    }

    private suspend fun getLastMessagePreview(chatId: String): MessagePreview? {
        return try {
            val lastMessage = database.messageQueries.getLastMessageForChat(chatId).executeAsOneOrNull()
            lastMessage?.let { message ->
                MessagePreview(
                    id = message.id,
                    content = message.content,
                    senderId = message.senderId,
                    senderName = message.senderName,
                    timestamp = message.createdAt,
                    type = ChatMessageType.valueOf(message.type),
                    isEdited = message.isEdited == 1L,
                    attachmentCount = message.attachmentCount.toInt()
                )
            }
        } catch (e: Exception) {
            null
        }
    }

    // Missing methods implementation for comprehensive chat functionality
    override suspend fun sendTextMessage(chatId: String, content: String, replyToId: String?): Result<Message> {
        return Result.failure(NotImplementedError("Text message sending not implemented in SQLDelight yet"))
    }

    override suspend fun sendMediaMessage(chatId: String, attachments: List<MessageAttachment>, caption: String?, replyToId: String?): Result<Message> {
        return Result.failure(NotImplementedError("Media message sending not implemented in SQLDelight yet"))
    }

    override suspend fun sendVoiceMessage(chatId: String, audioUrl: String, duration: Int): Result<Message> {
        return Result.failure(NotImplementedError("Voice message sending not implemented in SQLDelight yet"))
    }

    override suspend fun sendLocationMessage(chatId: String, latitude: Double, longitude: Double, address: String?): Result<Message> {
        return Result.failure(NotImplementedError("Location message sending not implemented in SQLDelight yet"))
    }

    override suspend fun forwardMessage(messageId: String, targetChatIds: List<String>): Result<List<Message>> {
        return Result.failure(NotImplementedError("Message forwarding not implemented in SQLDelight yet"))
    }

    override suspend fun startTyping(chatId: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun stopTyping(chatId: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun markMessageAsDelivered(messageId: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun markMessageAsRead(messageId: String, userId: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun updateMessageDeliveryStatus(messageId: String, status: MessageDeliveryStatus): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun uploadVoiceMessage(audioData: ByteArray): Result<String> {
        return Result.success("mock://voice/${System.currentTimeMillis()}.mp3") // Mock implementation
    }

    override suspend fun downloadVoiceMessage(url: String): Result<ByteArray> {
        return Result.success(ByteArray(0)) // Mock implementation
    }

    override suspend fun uploadFile(fileData: ByteArray, filename: String, mimeType: String): Result<MessageAttachment> {
        return Result.success(
            MessageAttachment(
                id = "mock_${System.currentTimeMillis()}",
                type = AttachmentType.FILE,
                url = "mock://file/$filename",
                filename = filename,
                fileSize = fileData.size.toLong(),
                mimeType = mimeType
            )
        )
    }

    override suspend fun downloadFile(attachment: MessageAttachment): Result<ByteArray> {
        return Result.success(ByteArray(0)) // Mock implementation
    }

    override suspend fun saveDraftMessage(chatId: String, content: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }

    override suspend fun getDraftMessage(chatId: String): Result<String?> {
        return Result.success(null) // Mock implementation
    }

    override suspend fun clearDraftMessage(chatId: String): Result<Boolean> {
        return Result.success(true) // Mock implementation
    }
}