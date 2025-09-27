package com.tchat.mobile.repositories

import com.tchat.mobile.models.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.map
import com.tchat.mobile.utils.PlatformUtils

/**
 * Chat Repository - Interface for chat and messaging functionality
 *
 * Provides comprehensive chat operations including session management,
 * messaging, and real-time updates. Ready for real API integration.
 */
interface ChatRepository {
    // Chat Session Operations
    suspend fun getChatSessions(): Result<List<ChatSession>>
    suspend fun getChatSession(chatId: String): Result<ChatSession>
    suspend fun createChatSession(session: ChatSession): Result<ChatSession>
    suspend fun updateChatSession(chatId: String, session: ChatSession): Result<ChatSession>
    suspend fun deleteChatSession(chatId: String): Result<Boolean>

    // Message Operations
    suspend fun getMessages(chatId: String, limit: Int = 50, offset: Int = 0): Result<List<Message>>
    suspend fun getMessage(messageId: String): Result<Message>
    suspend fun sendMessage(message: Message): Result<Message>
    suspend fun editMessage(messageId: String, newContent: String): Result<Message>
    suspend fun deleteMessage(messageId: String): Result<Boolean>

    // Participants Operations
    suspend fun addParticipant(chatId: String, participant: ChatParticipant): Result<Boolean>
    suspend fun removeParticipant(chatId: String, participantId: String): Result<Boolean>
    suspend fun updateParticipantStatus(chatId: String, participantId: String, status: ParticipantStatus): Result<Boolean>

    // Reactions and Interactions
    suspend fun addReaction(messageId: String, reaction: MessageReaction): Result<Boolean>
    suspend fun removeReaction(messageId: String, userId: String, emoji: String): Result<Boolean>

    // Chat Management
    suspend fun pinMessage(messageId: String): Result<Boolean>
    suspend fun unpinMessage(messageId: String): Result<Boolean>
    suspend fun markAsRead(chatId: String, messageId: String): Result<Boolean>
    suspend fun muteChat(chatId: String): Result<Boolean>
    suspend fun unmuteChat(chatId: String): Result<Boolean>
    suspend fun archiveChat(chatId: String): Result<Boolean>
    suspend fun unarchiveChat(chatId: String): Result<Boolean>

    // Search and Filtering
    suspend fun searchMessages(chatId: String, query: String): Result<List<Message>>
    suspend fun searchChatSessions(query: String): Result<List<ChatSession>>

    // Real-time updates
    fun observeChatSessions(): Flow<List<ChatSession>>
    fun observeMessages(chatId: String): Flow<List<Message>>
    fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>>
    fun observeParticipantStatus(chatId: String): Flow<List<ChatParticipant>>

    // Advanced messaging operations
    suspend fun sendTextMessage(chatId: String, content: String, replyToId: String? = null): Result<Message>
    suspend fun sendMediaMessage(chatId: String, attachments: List<MessageAttachment>, caption: String? = null, replyToId: String? = null): Result<Message>
    suspend fun sendVoiceMessage(chatId: String, audioUrl: String, duration: Int): Result<Message>
    suspend fun sendLocationMessage(chatId: String, latitude: Double, longitude: Double, address: String? = null): Result<Message>
    suspend fun forwardMessage(messageId: String, targetChatIds: List<String>): Result<List<Message>>

    // Typing indicators
    suspend fun startTyping(chatId: String): Result<Boolean>
    suspend fun stopTyping(chatId: String): Result<Boolean>

    // Message status updates
    suspend fun markMessageAsDelivered(messageId: String): Result<Boolean>
    suspend fun markMessageAsRead(messageId: String, userId: String): Result<Boolean>
    suspend fun updateMessageDeliveryStatus(messageId: String, status: MessageDeliveryStatus): Result<Boolean>

    // Voice message operations
    suspend fun uploadVoiceMessage(audioData: ByteArray): Result<String> // Returns URL
    suspend fun downloadVoiceMessage(url: String): Result<ByteArray>

    // File operations
    suspend fun uploadFile(fileData: ByteArray, filename: String, mimeType: String): Result<MessageAttachment>
    suspend fun downloadFile(attachment: MessageAttachment): Result<ByteArray>

    // Draft messages
    suspend fun saveDraftMessage(chatId: String, content: String): Result<Boolean>
    suspend fun getDraftMessage(chatId: String): Result<String?>
    suspend fun clearDraftMessage(chatId: String): Result<Boolean>
}

/**
 * Mock implementation of ChatRepository for development and testing
 */
class MockChatRepository : ChatRepository {

    private val _chatSessions = MutableStateFlow(generateMockChatSessions())
    private val _messages = MutableStateFlow(generateMockMessages())
    private val _typingIndicators = MutableStateFlow<Map<String, List<TypingIndicator>>>(emptyMap())
    private val _draftMessages = MutableStateFlow<Map<String, String>>(emptyMap())

    private val chatSessions: List<ChatSession> get() = _chatSessions.value
    private val messages: List<Message> get() = _messages.value

    override suspend fun getChatSessions(): Result<List<ChatSession>> {
        delay(200) // Simulate network delay
        return Result.success(chatSessions.sortByActivity())
    }

    override suspend fun getChatSession(chatId: String): Result<ChatSession> {
        delay(150)

        val session = chatSessions.find { it.id == chatId }
        return if (session != null) {
            Result.success(session)
        } else {
            Result.failure(Exception("Chat session not found"))
        }
    }

    override suspend fun createChatSession(session: ChatSession): Result<ChatSession> {
        delay(300)

        val newSession = session.copy(
            id = "chat_${PlatformUtils.currentTimeMillis()}",
            createdAt = getCurrentTimestamp(),
            updatedAt = getCurrentTimestamp()
        )

        _chatSessions.value = _chatSessions.value + newSession
        return Result.success(newSession)
    }

    override suspend fun updateChatSession(chatId: String, session: ChatSession): Result<ChatSession> {
        delay(250)

        val updatedSessions = _chatSessions.value.map { currentSession ->
            if (currentSession.id == chatId) {
                session.copy(updatedAt = getCurrentTimestamp())
            } else {
                currentSession
            }
        }

        _chatSessions.value = updatedSessions
        val updatedSession = updatedSessions.find { it.id == chatId }

        return if (updatedSession != null) {
            Result.success(updatedSession)
        } else {
            Result.failure(Exception("Chat session not found"))
        }
    }

    override suspend fun deleteChatSession(chatId: String): Result<Boolean> {
        delay(200)

        val originalSize = _chatSessions.value.size
        _chatSessions.value = _chatSessions.value.filter { it.id != chatId }

        // Also remove messages from this chat
        _messages.value = _messages.value.filter { it.chatId != chatId }

        return Result.success(_chatSessions.value.size < originalSize)
    }

    override suspend fun getMessages(chatId: String, limit: Int, offset: Int): Result<List<Message>> {
        delay(250)

        val chatMessages = messages
            .filter { it.chatId == chatId && !it.isDeleted }
            .sortedByDescending { it.createdAt }
            .drop(offset)
            .take(limit)

        return Result.success(chatMessages)
    }

    override suspend fun getMessage(messageId: String): Result<Message> {
        delay(100)

        val message = messages.find { it.id == messageId }
        return if (message != null) {
            Result.success(message)
        } else {
            Result.failure(Exception("Message not found"))
        }
    }

    override suspend fun sendMessage(message: Message): Result<Message> {
        delay(200)

        val newMessage = message.copy(
            id = "msg_${PlatformUtils.currentTimeMillis()}",
            createdAt = getCurrentTimestamp(),
            deliveryStatus = MessageDeliveryStatus.SENT
        )

        _messages.value = _messages.value + newMessage

        // Update chat session's last message and activity
        updateChatSessionLastMessage(newMessage)

        return Result.success(newMessage)
    }

    override suspend fun editMessage(messageId: String, newContent: String): Result<Message> {
        delay(150)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(
                    content = newContent,
                    isEdited = true,
                    editedAt = getCurrentTimestamp()
                )
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        val updatedMessage = updatedMessages.find { it.id == messageId }

        return if (updatedMessage != null) {
            Result.success(updatedMessage)
        } else {
            Result.failure(Exception("Message not found"))
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Boolean> {
        delay(150)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(
                    isDeleted = true,
                    deletedAt = getCurrentTimestamp(),
                    content = "This message was deleted"
                )
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun addParticipant(chatId: String, participant: ChatParticipant): Result<Boolean> {
        delay(200)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    participants = session.participants + participant,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun removeParticipant(chatId: String, participantId: String): Result<Boolean> {
        delay(200)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    participants = session.participants.filter { it.id != participantId },
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun updateParticipantStatus(chatId: String, participantId: String, status: ParticipantStatus): Result<Boolean> {
        delay(100)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    participants = session.participants.map { participant ->
                        if (participant.id == participantId) {
                            participant.copy(
                                status = status,
                                lastSeen = if (status == ParticipantStatus.OFFLINE) getCurrentTimestamp() else null
                            )
                        } else {
                            participant
                        }
                    }
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun addReaction(messageId: String, reaction: MessageReaction): Result<Boolean> {
        delay(100)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                val currentReactions = message.reactions.toMutableList()
                // Remove existing reaction from same user with same emoji
                currentReactions.removeAll { it.userId == reaction.userId && it.emoji == reaction.emoji }
                currentReactions.add(reaction)

                message.copy(reactions = currentReactions)
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun removeReaction(messageId: String, userId: String, emoji: String): Result<Boolean> {
        delay(100)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(
                    reactions = message.reactions.filter {
                        !(it.userId == userId && it.emoji == emoji)
                    }
                )
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun pinMessage(messageId: String): Result<Boolean> {
        delay(100)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(isPinned = true)
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun unpinMessage(messageId: String): Result<Boolean> {
        delay(100)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(isPinned = false)
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun markAsRead(chatId: String, messageId: String): Result<Boolean> {
        delay(50)

        // Update unread count for chat session
        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    unreadCount = 0,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun muteChat(chatId: String): Result<Boolean> {
        delay(100)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    isMuted = true,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun unmuteChat(chatId: String): Result<Boolean> {
        delay(100)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    isMuted = false,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun archiveChat(chatId: String): Result<Boolean> {
        delay(100)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    isArchived = true,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun unarchiveChat(chatId: String): Result<Boolean> {
        delay(100)

        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == chatId) {
                session.copy(
                    isArchived = false,
                    updatedAt = getCurrentTimestamp()
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
        return Result.success(true)
    }

    override suspend fun searchMessages(chatId: String, query: String): Result<List<Message>> {
        delay(300)

        val searchResults = messages
            .filter { it.chatId == chatId && !it.isDeleted }
            .filter { message ->
                message.content.contains(query, ignoreCase = true) ||
                message.senderName.contains(query, ignoreCase = true)
            }
            .sortedByDescending { it.createdAt }

        return Result.success(searchResults)
    }

    override suspend fun searchChatSessions(query: String): Result<List<ChatSession>> {
        delay(250)

        val searchResults = chatSessions.searchByName(query, "current_user")
        return Result.success(searchResults)
    }

    override fun observeChatSessions(): Flow<List<ChatSession>> {
        return _chatSessions.asStateFlow().map { sessions ->
            sessions.sortByActivity()
        }
    }

    override fun observeMessages(chatId: String): Flow<List<Message>> {
        return _messages.asStateFlow().map { allMessages ->
            allMessages
                .filter { it.chatId == chatId && !it.isDeleted }
                .sortedByDescending { it.createdAt }
        }
    }

    override fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>> {
        return _typingIndicators.asStateFlow().map { indicators ->
            indicators[chatId] ?: emptyList()
        }
    }

    override fun observeParticipantStatus(chatId: String): Flow<List<ChatParticipant>> {
        return _chatSessions.asStateFlow().map { sessions ->
            sessions.find { it.id == chatId }?.participants ?: emptyList()
        }
    }

    // Helper functions
    private fun getCurrentTimestamp(): String {
        return "${PlatformUtils.currentTimeMillis()}"
    }

    private fun updateChatSessionLastMessage(message: Message) {
        val updatedSessions = _chatSessions.value.map { session ->
            if (session.id == message.chatId) {
                session.copy(
                    lastMessage = MessagePreview(
                        id = message.id,
                        content = message.getDisplayContent(),
                        senderId = message.senderId,
                        senderName = message.senderName,
                        timestamp = message.createdAt,
                        type = ChatMessageType.valueOf(message.type.name),
                        isEdited = message.isEdited,
                        attachmentCount = message.attachments.size
                    ),
                    lastActivityAt = message.createdAt,
                    updatedAt = message.createdAt,
                    unreadCount = if (message.senderId != "current_user") session.unreadCount + 1 else session.unreadCount
                )
            } else {
                session
            }
        }

        _chatSessions.value = updatedSessions
    }

    // Advanced messaging operations implementation
    override suspend fun sendTextMessage(chatId: String, content: String, replyToId: String?): Result<Message> {
        delay(200)

        val newMessage = Message(
            id = "msg_${PlatformUtils.currentTimeMillis()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            type = MessageType.TEXT,
            content = content,
            replyToId = replyToId,
            createdAt = getCurrentTimestamp(),
            deliveryStatus = MessageDeliveryStatus.SENT
        )

        _messages.value = _messages.value + newMessage
        updateChatSessionLastMessage(newMessage)

        return Result.success(newMessage)
    }

    override suspend fun sendMediaMessage(chatId: String, attachments: List<MessageAttachment>, caption: String?, replyToId: String?): Result<Message> {
        delay(300)

        val messageType = when {
            attachments.any { it.type == AttachmentType.IMAGE } -> MessageType.IMAGE
            attachments.any { it.type == AttachmentType.VIDEO } -> MessageType.VIDEO
            attachments.any { it.type == AttachmentType.AUDIO } -> MessageType.AUDIO
            attachments.any { it.type == AttachmentType.FILE } -> MessageType.FILE
            attachments.any { it.type == AttachmentType.LOCATION } -> MessageType.LOCATION
            else -> MessageType.FILE
        }

        val newMessage = Message(
            id = "msg_${PlatformUtils.currentTimeMillis()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            type = messageType,
            content = caption ?: attachments.firstOrNull()?.filename ?: "Media",
            attachments = attachments,
            replyToId = replyToId,
            createdAt = getCurrentTimestamp(),
            deliveryStatus = MessageDeliveryStatus.SENT
        )

        _messages.value = _messages.value + newMessage
        updateChatSessionLastMessage(newMessage)

        return Result.success(newMessage)
    }

    override suspend fun sendVoiceMessage(chatId: String, audioUrl: String, duration: Int): Result<Message> {
        delay(250)

        val audioAttachment = MessageAttachment(
            id = "audio_${PlatformUtils.currentTimeMillis()}",
            type = AttachmentType.AUDIO,
            url = audioUrl,
            duration = duration,
            mimeType = "audio/mp4"
        )

        val newMessage = Message(
            id = "msg_${PlatformUtils.currentTimeMillis()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            type = MessageType.AUDIO,
            content = "Voice message",
            attachments = listOf(audioAttachment),
            createdAt = getCurrentTimestamp(),
            deliveryStatus = MessageDeliveryStatus.SENT
        )

        _messages.value = _messages.value + newMessage
        updateChatSessionLastMessage(newMessage)

        return Result.success(newMessage)
    }

    override suspend fun sendLocationMessage(chatId: String, latitude: Double, longitude: Double, address: String?): Result<Message> {
        delay(200)

        val locationAttachment = MessageAttachment(
            id = "location_${PlatformUtils.currentTimeMillis()}",
            type = AttachmentType.LOCATION,
            url = "geo:$latitude,$longitude",
            metadata = mapOf(
                "latitude" to latitude.toString(),
                "longitude" to longitude.toString(),
                "address" to (address ?: "Unknown location")
            )
        )

        val newMessage = Message(
            id = "msg_${PlatformUtils.currentTimeMillis()}",
            chatId = chatId,
            senderId = "current_user",
            senderName = "You",
            type = MessageType.LOCATION,
            content = address ?: "Location shared",
            attachments = listOf(locationAttachment),
            createdAt = getCurrentTimestamp(),
            deliveryStatus = MessageDeliveryStatus.SENT
        )

        _messages.value = _messages.value + newMessage
        updateChatSessionLastMessage(newMessage)

        return Result.success(newMessage)
    }

    override suspend fun forwardMessage(messageId: String, targetChatIds: List<String>): Result<List<Message>> {
        delay(300)

        val originalMessage = messages.find { it.id == messageId }
            ?: return Result.failure(Exception("Message not found"))

        val forwardedMessages = targetChatIds.map { chatId ->
            Message(
                id = "msg_${PlatformUtils.currentTimeMillis()}_${chatId}",
                chatId = chatId,
                senderId = "current_user",
                senderName = "You",
                type = originalMessage.type,
                content = originalMessage.content,
                attachments = originalMessage.attachments,
                createdAt = getCurrentTimestamp(),
                deliveryStatus = MessageDeliveryStatus.SENT
            )
        }

        _messages.value = _messages.value + forwardedMessages
        forwardedMessages.forEach { updateChatSessionLastMessage(it) }

        return Result.success(forwardedMessages)
    }

    override suspend fun startTyping(chatId: String): Result<Boolean> {
        delay(50)

        val typingIndicator = TypingIndicator(
            userId = "current_user",
            userName = "You",
            startedAt = getCurrentTimestamp(),
            expiresAt = getCurrentTimestamp() // In real implementation, this would be current time + 5 seconds
        )

        val currentIndicators = _typingIndicators.value[chatId] ?: emptyList()
        val updatedIndicators = currentIndicators.filter { it.userId != "current_user" } + typingIndicator

        _typingIndicators.value = _typingIndicators.value + (chatId to updatedIndicators)
        return Result.success(true)
    }

    override suspend fun stopTyping(chatId: String): Result<Boolean> {
        delay(50)

        val currentIndicators = _typingIndicators.value[chatId] ?: emptyList()
        val updatedIndicators = currentIndicators.filter { it.userId != "current_user" }

        _typingIndicators.value = _typingIndicators.value + (chatId to updatedIndicators)
        return Result.success(true)
    }

    override suspend fun markMessageAsDelivered(messageId: String): Result<Boolean> {
        return updateMessageDeliveryStatus(messageId, MessageDeliveryStatus.DELIVERED)
    }

    override suspend fun markMessageAsRead(messageId: String, userId: String): Result<Boolean> {
        delay(50)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(
                    deliveryStatus = MessageDeliveryStatus.READ,
                    readBy = message.readBy + userId
                )
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun updateMessageDeliveryStatus(messageId: String, status: MessageDeliveryStatus): Result<Boolean> {
        delay(50)

        val updatedMessages = _messages.value.map { message ->
            if (message.id == messageId) {
                message.copy(deliveryStatus = status)
            } else {
                message
            }
        }

        _messages.value = updatedMessages
        return Result.success(true)
    }

    override suspend fun uploadVoiceMessage(audioData: ByteArray): Result<String> {
        delay(1000) // Simulate upload time

        // Mock upload - return a fake URL
        val mockUrl = "https://example.com/voice/${PlatformUtils.currentTimeMillis()}.mp4"
        return Result.success(mockUrl)
    }

    override suspend fun downloadVoiceMessage(url: String): Result<ByteArray> {
        delay(500) // Simulate download time

        // Mock download - return empty byte array
        return Result.success(ByteArray(0))
    }

    override suspend fun uploadFile(fileData: ByteArray, filename: String, mimeType: String): Result<MessageAttachment> {
        delay(1500) // Simulate upload time

        val attachment = MessageAttachment(
            id = "file_${PlatformUtils.currentTimeMillis()}",
            type = when {
                mimeType.startsWith("image/") -> AttachmentType.IMAGE
                mimeType.startsWith("video/") -> AttachmentType.VIDEO
                mimeType.startsWith("audio/") -> AttachmentType.AUDIO
                else -> AttachmentType.FILE
            },
            url = "https://example.com/files/${PlatformUtils.currentTimeMillis()}",
            filename = filename,
            fileSize = fileData.size.toLong(),
            mimeType = mimeType
        )

        return Result.success(attachment)
    }

    override suspend fun downloadFile(attachment: MessageAttachment): Result<ByteArray> {
        delay(1000) // Simulate download time

        // Mock download - return empty byte array
        return Result.success(ByteArray(0))
    }

    override suspend fun saveDraftMessage(chatId: String, content: String): Result<Boolean> {
        delay(100)

        _draftMessages.value = _draftMessages.value + (chatId to content)
        return Result.success(true)
    }

    override suspend fun getDraftMessage(chatId: String): Result<String?> {
        delay(50)

        val draft = _draftMessages.value[chatId]
        return Result.success(draft)
    }

    override suspend fun clearDraftMessage(chatId: String): Result<Boolean> {
        delay(50)

        _draftMessages.value = _draftMessages.value - chatId
        return Result.success(true)
    }

    private fun generateMockChatSessions(): List<ChatSession> {
        return listOf(
            ChatSession(
                id = "chat_1",
                name = null, // Direct message
                type = ChatType.DIRECT,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = "2024-01-01T10:00:00Z"
                    ),
                    ChatParticipant(
                        id = "user_sarah",
                        name = "Sarah Chen",
                        avatar = "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = "2024-01-01T10:00:00Z"
                    )
                ),
                lastMessage = MessagePreview(
                    id = "msg_1",
                    content = "Hey! How's your day going?",
                    senderId = "user_sarah",
                    senderName = "Sarah Chen",
                    timestamp = "2 hours ago",
                    type = ChatMessageType.TEXT
                ),
                unreadCount = 2,
                metadata = ChatMetadata(),
                createdAt = "2024-01-01T10:00:00Z",
                updatedAt = "2024-01-15T14:30:00Z",
                lastActivityAt = "2024-01-15T14:30:00Z"
            ),
            ChatSession(
                id = "chat_2",
                name = "Design Team",
                type = ChatType.GROUP,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        role = ChatRole.ADMIN,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = "2024-01-01T10:00:00Z"
                    ),
                    ChatParticipant(
                        id = "user_alex",
                        name = "Alex Johnson",
                        avatar = "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.AWAY,
                        joinedAt = "2024-01-01T10:00:00Z"
                    ),
                    ChatParticipant(
                        id = "user_maria",
                        name = "Maria Rodriguez",
                        avatar = "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.OFFLINE,
                        lastSeen = "1 hour ago",
                        joinedAt = "2024-01-01T10:00:00Z"
                    )
                ),
                lastMessage = MessagePreview(
                    id = "msg_2",
                    content = "The new mockups look great!",
                    senderId = "user_alex",
                    senderName = "Alex Johnson",
                    timestamp = "45 minutes ago",
                    type = ChatMessageType.TEXT
                ),
                unreadCount = 0,
                isPinned = true,
                metadata = ChatMetadata(
                    description = "Design team collaboration space",
                    isPublic = false
                ),
                createdAt = "2024-01-01T10:00:00Z",
                updatedAt = "2024-01-15T13:45:00Z",
                lastActivityAt = "2024-01-15T13:45:00Z"
            ),
            ChatSession(
                id = "chat_3",
                name = null, // Direct message
                type = ChatType.DIRECT,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = "2024-01-01T10:00:00Z"
                    ),
                    ChatParticipant(
                        id = "user_mike",
                        name = "Mike Wilson",
                        avatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100",
                        role = ChatRole.MEMBER,
                        status = ParticipantStatus.BUSY,
                        joinedAt = "2024-01-01T10:00:00Z"
                    )
                ),
                lastMessage = MessagePreview(
                    id = "msg_3",
                    content = "📄 project_requirements.pdf",
                    senderId = "user_mike",
                    senderName = "Mike Wilson",
                    timestamp = "3 hours ago",
                    type = ChatMessageType.FILE,
                    attachmentCount = 1
                ),
                unreadCount = 1,
                metadata = ChatMetadata(),
                createdAt = "2024-01-01T10:00:00Z",
                updatedAt = "2024-01-15T11:30:00Z",
                lastActivityAt = "2024-01-15T11:30:00Z"
            )
        )
    }

    private fun generateMockMessages(): List<Message> {
        return listOf(
            Message(
                id = "msg_1",
                chatId = "chat_1",
                senderId = "user_sarah",
                senderName = "Sarah Chen",
                senderAvatar = "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100",
                type = MessageType.TEXT,
                content = "Hey! How's your day going?",
                createdAt = "2024-01-15T14:30:00Z",
                deliveryStatus = MessageDeliveryStatus.DELIVERED
            ),
            Message(
                id = "msg_2",
                chatId = "chat_2",
                senderId = "user_alex",
                senderName = "Alex Johnson",
                senderAvatar = "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100",
                type = MessageType.TEXT,
                content = "The new mockups look great!",
                createdAt = "2024-01-15T13:45:00Z",
                reactions = listOf(
                    MessageReaction(
                        emoji = "👍",
                        userId = "current_user",
                        userName = "You",
                        timestamp = "2024-01-15T13:46:00Z"
                    )
                ),
                deliveryStatus = MessageDeliveryStatus.READ
            ),
            Message(
                id = "msg_3",
                chatId = "chat_3",
                senderId = "user_mike",
                senderName = "Mike Wilson",
                senderAvatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100",
                type = MessageType.FILE,
                content = "Here are the project requirements",
                attachments = listOf(
                    MessageAttachment(
                        id = "att_1",
                        type = AttachmentType.FILE,
                        url = "https://example.com/project_requirements.pdf",
                        filename = "project_requirements.pdf",
                        fileSize = 2048576, // 2MB
                        mimeType = "application/pdf"
                    )
                ),
                createdAt = "2024-01-15T11:30:00Z",
                deliveryStatus = MessageDeliveryStatus.DELIVERED
            )
        )
    }
}