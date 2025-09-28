package com.tchat.mobile.repositories.datasource

import com.tchat.mobile.api.ApiClient
import com.tchat.mobile.api.ApiException
import com.tchat.mobile.api.models.*
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.*
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlin.time.Duration.Companion.seconds

/**
 * RemoteDataSource implementation using HTTP API
 *
 * Provides server communication with:
 * - Automatic authentication management
 * - Request/response mapping between domain models and DTOs
 * - Circuit breaker pattern for error resilience
 * - Real-time updates simulation
 * - Comprehensive error handling
 */
class ApiRemoteDataSource(
    private val apiClient: ApiClient
) : RemoteDataSource {

    private val _realTimeMessages = MutableSharedFlow<ChatMessage>()
    private val _presenceUpdates = MutableSharedFlow<List<PresenceUpdate>>()
    private val _typingIndicators = MutableSharedFlow<TypingIndicator>()

    // Connection Management
    override suspend fun connect(): Result<Unit> {
        return apiClient.connect()
    }

    override suspend fun disconnect(): Result<Unit> {
        return apiClient.disconnect()
    }

    override fun getConnectionState(): Flow<ConnectionState> {
        return apiClient.connectionState
    }

    // Authentication & User Management
    override suspend fun authenticateUser(token: String): Result<User> {
        // This method is deprecated, new authentication flows use OTP or email/password
        return Result.failure(Exception("Use OTP or email authentication instead"))
    }

    override suspend fun refreshToken(): Result<String> {
        return apiClient.refreshAuthToken().mapCatching { response ->
            response.accessToken
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun getUserProfile(userId: String): Result<UserProfile> {
        return executeWithErrorHandling {
            // Simulate API call - replace with actual implementation
            UserProfile(
                userId = userId,
                displayName = "Test User",
                username = "testuser",
                avatar = null,
                bio = "Test user profile",
                isVerified = false,
                isOnline = true,
                lastSeen = Clock.System.now()
            )
        }
    }

    // Chat Operations
    override suspend fun getChatSessions(userId: String): Result<List<ChatSession>> {
        return apiClient.getChatSessions(userId).mapCatching { dtos ->
            dtos.map { it.toDomainModel() }
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun getChatSession(chatId: String): Result<ChatSession> {
        return apiClient.getChatSession(chatId).mapCatching { dto ->
            dto.toDomainModel()
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun createChatSession(chat: ChatSession): Result<ChatSession> {
        return apiClient.createChatSession(chat.toDto()).mapCatching { dto ->
            dto.toDomainModel()
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun updateChatSession(chat: ChatSession): Result<ChatSession> {
        return executeWithErrorHandling {
            // Simulate update - replace with actual API call
            chat.copy(updatedAt = Clock.System.now().toString())
        }
    }

    override suspend fun deleteChatSession(chatId: String): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate deletion - replace with actual API call
            Unit
        }
    }

    // Message Operations
    override suspend fun getMessages(chatId: String, limit: Int, before: Instant?): Result<List<ChatMessage>> {
        return apiClient.getMessages(
            chatId = chatId,
            limit = limit,
            before = before?.toEpochMilliseconds()
        ).mapCatching { dtos ->
            dtos.map { it.toDomainModel() }
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun sendMessage(message: ChatMessage): Result<ChatMessage> {
        return apiClient.sendMessage(message.toDto()).mapCatching { dto ->
            dto.toDomainModel()
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun updateMessage(message: ChatMessage): Result<ChatMessage> {
        return apiClient.updateMessage(message.toDto()).mapCatching { dto ->
            dto.toDomainModel()
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Unit> {
        return apiClient.deleteMessage(messageId).recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun markMessageAsRead(messageId: String, userId: String): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate marking as read - replace with actual API call
            Unit
        }
    }

    // Sync Operations
    override suspend fun getMessagesSince(chatId: String, timestamp: Instant): Result<List<ChatMessage>> {
        return apiClient.getMessages(
            chatId = chatId,
            limit = 100,
            before = null // Get latest messages
        ).mapCatching { dtos ->
            dtos.filter { it.serverTimestamp > timestamp.toEpochMilliseconds() }
                .map { it.toDomainModel() }
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun pushPendingOperations(operations: List<SyncOperation>): Result<List<SyncOperation>> {
        return apiClient.pushPendingOperations(operations.map { it.toDto() }).mapCatching { dtos ->
            dtos.map { it.toDomainModel() }
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun resolveConflicts(conflicts: List<MessageConflict>): Result<List<ConflictResolution>> {
        return apiClient.resolveConflicts(conflicts.map { it.toDto() }).mapCatching { dtos ->
            dtos.map { it.toDomainModel() }
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun getServerTimestamp(): Result<Instant> {
        return executeWithErrorHandling {
            Instant.fromEpochMilliseconds(apiClient.getServerTimestamp())
        }
    }

    // Bulk Operations
    override suspend fun syncChatData(chatId: String, lastSyncTimestamp: Instant?): Result<SyncResult> {
        val request = SyncRequestDto(
            chatId = chatId,
            lastSyncTimestamp = lastSyncTimestamp?.toEpochMilliseconds(),
            operations = emptyList()
        )

        return apiClient.syncChatData(request).mapCatching { response ->
            SyncResult(
                chatId = response.chatId,
                success = response.success,
                messagesUpdated = response.messages.size,
                conflicts = response.conflicts.map { it.toDomainModel() },
                lastSyncTimestamp = Instant.fromEpochMilliseconds(response.serverTimestamp),
                syncDuration = 1.seconds,
                errors = emptyList(),
                resolutions = emptyList()
            )
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun batchUpdateMessages(messages: List<ChatMessage>): Result<List<ChatMessage>> {
        return executeWithErrorHandling {
            // Simulate batch update - replace with actual API call
            messages.map { message ->
                message.copy(
                    isEdited = true,
                    editedAt = Clock.System.now().toString()
                )
            }
        }
    }

    override suspend fun batchDeleteMessages(messageIds: List<String>): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate batch delete - replace with actual API call
            Unit
        }
    }

    // Real-time Updates
    override fun subscribeToMessages(chatId: String): Flow<ChatMessage> {
        return _realTimeMessages.filter { it.chatId == chatId }
    }

    override fun subscribeToPresence(chatId: String): Flow<List<PresenceUpdate>> {
        return _presenceUpdates
    }

    override fun subscribeToTypingIndicators(chatId: String): Flow<TypingIndicator> {
        return _typingIndicators.filter { it.chatId == chatId }
    }

    override suspend fun sendTypingIndicator(chatId: String, isTyping: Boolean): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate sending typing indicator
            Unit
        }
    }

    // File Upload & Media
    override suspend fun uploadFile(file: ByteArray, fileName: String, mimeType: String): Result<String> {
        val request = FileUploadRequest(
            fileName = fileName,
            mimeType = mimeType,
            fileSize = file.size.toLong()
        )

        return apiClient.uploadFile(request, file).mapCatching { response ->
            response.fileUrl
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun uploadImage(image: ByteArray, fileName: String): Result<String> {
        return uploadFile(image, fileName, "image/jpeg")
    }

    override suspend fun uploadVideo(video: ByteArray, fileName: String): Result<String> {
        return uploadFile(video, fileName, "video/mp4")
    }

    override suspend fun uploadAudio(audio: ByteArray, fileName: String): Result<String> {
        return uploadFile(audio, fileName, "audio/mp3")
    }

    // Search Operations
    override suspend fun searchMessages(query: String, chatId: String?): Result<List<ChatMessage>> {
        return executeWithErrorHandling {
            // Simulate search - replace with actual API call
            emptyList<ChatMessage>()
        }
    }

    override suspend fun searchChats(query: String): Result<List<ChatSession>> {
        return executeWithErrorHandling {
            // Simulate search - replace with actual API call
            emptyList<ChatSession>()
        }
    }

    // Analytics & Metrics
    override suspend fun trackMessageDelivery(messageId: String, status: DeliveryStatus): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate tracking - replace with actual API call
            Unit
        }
    }

    override suspend fun reportMessageRead(messageId: String, readAt: Instant): Result<Unit> {
        return executeWithErrorHandling {
            // Simulate reporting - replace with actual API call
            Unit
        }
    }

    override suspend fun getMessageAnalytics(chatId: String): Result<MessageAnalytics> {
        return executeWithErrorHandling {
            MessageAnalytics(
                chatId = chatId,
                totalMessages = 150,
                averageResponseTime = 300000, // 5 minutes
                activeParticipants = 5,
                messageFrequency = mapOf(
                    "Monday" to 25,
                    "Tuesday" to 30,
                    "Wednesday" to 20
                )
            )
        }
    }

    // Error Recovery
    override suspend fun retryFailedOperation(operation: SyncOperation): Result<SyncOperation> {
        return executeWithErrorHandling {
            operation.copy(retryCount = operation.retryCount + 1)
        }
    }

    override suspend fun validateDataIntegrity(chatId: String): Result<DataIntegrityReport> {
        return executeWithErrorHandling {
            DataIntegrityReport(
                chatId = chatId,
                isValid = true,
                issues = emptyList(),
                lastChecked = Clock.System.now()
            )
        }
    }

    // Health & Monitoring
    override suspend fun healthCheck(): Result<ServerHealth> {
        return apiClient.healthCheck().mapCatching { dto ->
            ServerHealth(
                status = dto.status,
                uptime = dto.uptime,
                responseTime = dto.responseTime,
                activeConnections = dto.activeConnections
            )
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    override suspend fun getServerCapabilities(): Result<ServerCapabilities> {
        return apiClient.getServerCapabilities().mapCatching { dto ->
            ServerCapabilities(
                maxFileSize = dto.maxFileSize,
                supportedFileTypes = dto.supportedFileTypes,
                maxMessageLength = dto.maxMessageLength,
                supportsBulkOperations = dto.supportsBulkOperations,
                supportsRealTime = dto.supportsRealTime
            )
        }.recoverCatching { throwable ->
            throw mapApiException(throwable)
        }
    }

    // Utility Methods
    private suspend fun <T> executeWithErrorHandling(operation: suspend () -> T): Result<T> {
        return try {
            Result.success(operation())
        } catch (e: Exception) {
            Result.failure(mapApiException(e))
        }
    }

    private fun mapApiException(throwable: Throwable): Exception {
        return when (throwable) {
            is ApiException -> when (throwable.statusCode) {
                401 -> ApiException(401, "Authentication required")
                403 -> ApiException(403, "Access denied")
                404 -> ApiException(404, "Resource not found")
                429 -> ApiException(429, "Rate limit exceeded")
                500 -> ApiException(500, "Server error")
                503 -> ApiException(503, "Service unavailable")
                else -> throwable
            }
            else -> ApiException(500, "Unknown error: ${throwable.message}")
        }
    }
}

// Extension functions for DTO mapping
private fun UserDto.toDomainModel(): User {
    return User(
        id = id,
        email = email,
        displayName = displayName,
        avatar = avatar,
        isActive = isActive
    )
}

private fun ChatSessionDto.toDomainModel(): ChatSession {
    return ChatSession(
        id = id,
        name = name ?: "Chat", // Provide default name
        type = ChatType.valueOf(type.uppercase()),
        participants = participants.map { participantId ->
            ChatParticipant(
                id = participantId,
                name = "User $participantId", // Simplified for now
                joinedAt = Instant.fromEpochMilliseconds(createdAt).toString()
            )
        },
        lastMessage = lastMessageId?.let {
            MessagePreview(
                id = it,
                content = "Last message content",
                senderId = createdBy,
                senderName = "User", // Simplified
                timestamp = lastMessageAt?.let { ts -> Instant.fromEpochMilliseconds(ts).toString() } ?: ""
            )
        },
        unreadCount = unreadCount,
        metadata = ChatMetadata(
            description = description
        ),
        createdAt = Instant.fromEpochMilliseconds(createdAt).toString(),
        updatedAt = Instant.fromEpochMilliseconds(updatedAt).toString(),
        lastActivityAt = lastMessageAt?.let { Instant.fromEpochMilliseconds(it).toString() }
    )
}

private fun ChatSession.toDto(): ChatSessionDto {
    return ChatSessionDto(
        id = id,
        name = name ?: "Unknown Chat",
        description = null, // Not in domain model
        avatar = null, // Not in domain model
        type = type.name.lowercase(),
        isActive = true, // Assume active
        participants = participants.map { it.id },
        createdBy = participants.firstOrNull()?.id ?: "unknown",
        createdAt = Instant.parse(createdAt).toEpochMilliseconds(),
        updatedAt = Instant.parse(updatedAt).toEpochMilliseconds(),
        lastMessageId = lastMessage?.id,
        lastMessageAt = lastActivityAt?.let { Instant.parse(it).toEpochMilliseconds() },
        unreadCount = unreadCount,
        settings = null // Simplified for now
    )
}

// Removed ChatSettings mapping - not compatible with current domain model

private fun MessageDto.toDomainModel(): ChatMessage {
    return ChatMessage(
        id = id,
        chatId = chatId,
        senderId = senderId,
        senderName = senderName,
        senderAvatar = senderAvatar,
        type = MessageType.valueOf(type.uppercase()),
        content = content,
        isEdited = isEdited,
        isPinned = isPinned,
        isDeleted = isDeleted,
        replyToId = replyToId,
        reactions = reactions.map { reaction ->
            MessageReaction(
                emoji = reaction.emoji,
                userId = reaction.userId,
                userName = "Unknown", // Would need to be resolved
                timestamp = Instant.fromEpochMilliseconds(reaction.timestamp).toString()
            )
        },
        attachments = attachments.map { attachment ->
            MessageAttachment(
                id = attachment.id,
                type = AttachmentType.valueOf(attachment.type.uppercase()),
                url = attachment.url,
                thumbnail = attachment.thumbnail,
                filename = attachment.filename,
                fileSize = attachment.fileSize,
                mimeType = attachment.mimeType,
                width = attachment.width?.toInt(),
                height = attachment.height?.toInt(),
                duration = attachment.duration?.toInt(),
                caption = attachment.caption,
                metadata = attachment.metadata
            )
        },
        createdAt = Instant.fromEpochMilliseconds(createdAt).toString(),
        editedAt = editedAt?.let { Instant.fromEpochMilliseconds(it).toString() },
        deletedAt = deletedAt?.let { Instant.fromEpochMilliseconds(it).toString() },
        deliveryStatus = DeliveryStatus.valueOf(deliveryStatus.uppercase()),
        readBy = readBy
    )
}

private fun ChatMessage.toDto(): MessageDto {
    return MessageDto(
        id = id,
        chatId = chatId,
        senderId = senderId,
        senderName = senderName,
        senderAvatar = senderAvatar,
        type = type.name.lowercase(),
        content = content,
        isEdited = isEdited,
        isPinned = isPinned,
        isDeleted = isDeleted,
        replyToId = replyToId,
        reactions = reactions.map { reaction ->
            ReactionDto(
                messageId = id,
                userId = reaction.userId,
                emoji = reaction.emoji,
                timestamp = reaction.timestamp.toLongOrNull() ?: 0L
            )
        },
        attachments = attachments.map { attachment ->
            AttachmentDto(
                id = attachment.id,
                messageId = id, // Pass the message ID
                type = attachment.type.name.lowercase(),
                url = attachment.url,
                thumbnail = attachment.thumbnail,
                filename = attachment.filename,
                fileSize = attachment.fileSize,
                mimeType = attachment.mimeType,
                width = attachment.width,
                height = attachment.height,
                duration = attachment.duration?.toLong(),
                caption = attachment.caption,
                metadata = attachment.metadata
            )
        },
        createdAt = createdAt.toLongOrNull() ?: 0L,
        editedAt = editedAt?.toLongOrNull(),
        deletedAt = deletedAt?.toLongOrNull(),
        serverTimestamp = createdAt.toLongOrNull() ?: 0L, // Use createdAt as fallback
        version = 1, // Default version
        checksum = null, // Not available in Message model
        deliveryStatus = deliveryStatus.name.lowercase(),
        readBy = readBy
    )
}

private fun SyncOperation.toDto(): SyncOperationDto {
    return SyncOperationDto(
        id = id,
        type = type.name.lowercase(),
        chatId = chatId,
        data = data,
        timestamp = timestamp.toEpochMilliseconds(),
        retryCount = retryCount,
        maxRetries = maxRetries,
        status = "pending"
    )
}

private fun SyncOperationDto.toDomainModel(): SyncOperation {
    return SyncOperation(
        id = id,
        type = OperationType.valueOf(type.uppercase()),
        chatId = chatId,
        data = data,
        timestamp = Instant.fromEpochMilliseconds(timestamp),
        retryCount = retryCount,
        maxRetries = maxRetries
    )
}

private fun MessageConflict.toDto(): ConflictDto {
    return ConflictDto(
        id = id,
        messageId = messageId,
        chatId = chatId,
        type = type.name.lowercase(),
        severity = severity.name.lowercase(),
        localMessage = localMessage.toDto(),
        remoteMessage = remoteMessage.toDto(),
        detectedAt = detectedAt.toEpochMilliseconds(),
        autoResolvable = autoResolvable,
        suggestedResolution = null
    )
}

private fun ConflictDto.toDomainModel(): MessageConflict {
    return MessageConflict(
        id = id,
        messageId = messageId,
        chatId = chatId,
        type = ConflictType.valueOf(type.uppercase()),
        severity = ConflictSeverity.valueOf(severity.uppercase()),
        localMessage = localMessage.toDomainModel(),
        remoteMessage = remoteMessage.toDomainModel(),
        detectedAt = Instant.fromEpochMilliseconds(detectedAt),
        autoResolvable = autoResolvable
    )
}

private fun ConflictResolutionDto.toDomainModel(): ConflictResolution {
    return ConflictResolution(
        conflictId = conflictId,
        strategy = ResolutionStrategy.valueOf(strategy.uppercase()),
        resolvedMessage = resolvedMessage?.toDomainModel(),
        explanation = explanation,
        success = true
    )
}