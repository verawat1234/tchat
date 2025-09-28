# SQLDelight ↔ API Implementation Guide

## Migration Strategy Overview

This guide demonstrates how to refactor the existing ChatRepository to implement the new LocalDataSource/RemoteDataSource/SyncEngine architecture while maintaining backward compatibility.

### Current Architecture Issues
```kotlin
// Current ChatRepository (lines 1129-1523)
class DatabaseChatRepository(
    private val database: TchatDatabase
) : ChatRepository {
    // Issues:
    // 1. Mixed concerns: database, API, business logic
    // 2. No clear separation between local and remote operations
    // 3. In-memory state mixed with persistent storage
    // 4. No conflict resolution or sync strategies
    // 5. Mock data generation in production code
}
```

### Target Architecture
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   ChatRepository│    │   SyncEngine     │    │ ConflictResolver│
│  (Coordination) │◄──►│  (Orchestration) │◄──►│   (Resolution)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │
         ▼                        ▼
┌─────────────────┐    ┌──────────────────┐
│ LocalDataSource │    │ RemoteDataSource │
│   (SQLDelight)  │    │     (API)        │
└─────────────────┘    └──────────────────┘
```

## Step 1: Create LocalDataSource Implementation

### 1.1 SQLDelightLocalDataSource
```kotlin
// File: repositories/datasource/SQLDelightLocalDataSource.kt
package com.tchat.mobile.repositories.datasource

import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant

class SQLDelightLocalDataSource(
    private val database: TchatDatabase
) : LocalDataSource {

    // Message Operations
    override suspend fun getMessages(chatId: String): Flow<List<ChatMessage>> {
        return database.messageQueries.getMessagesByChatId(chatId)
            .asFlow()
            .mapToList()
            .map { dbMessages ->
                dbMessages.map { it.toDomainMessage() }
            }
    }

    override suspend fun saveMessage(message: ChatMessage): Result<Unit> {
        return try {
            database.messageQueries.insertMessage(
                id = message.id,
                chatId = message.chatId,
                senderId = message.senderId,
                senderName = message.senderInfo.displayName,
                senderAvatar = message.senderInfo.avatar,
                type = message.type.name,
                content = message.content.toJsonString(),
                timestamp = message.timestamp.toEpochMilliseconds(),
                status = message.status.name,
                replyToId = message.replyTo,
                isEdited = if (message.isEdited) 1L else 0L,
                editedAt = message.editedAt?.toEpochMilliseconds(),
                isDeleted = if (message.isDeleted) 1L else 0L,
                deletedAt = message.deletedAt?.toEpochMilliseconds(),
                reactions = message.reactions.toJsonString(),
                syncStatus = SyncStatus.PENDING.name,
                lastSyncAt = null
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateMessageStatus(messageId: String, status: MessageStatus): Result<Unit> {
        return try {
            database.messageQueries.updateMessageStatus(
                status = status.name,
                id = messageId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun markMessageAsSynced(messageId: String, serverTimestamp: Instant): Result<Unit> {
        return try {
            database.messageQueries.markMessageSynced(
                syncStatus = SyncStatus.SYNCED.name,
                lastSyncAt = serverTimestamp.toEpochMilliseconds(),
                serverTimestamp = serverTimestamp.toEpochMilliseconds(),
                id = messageId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getPendingSyncMessages(): List<ChatMessage> {
        return database.messageQueries.getPendingSyncMessages()
            .executeAsList()
            .map { it.toDomainMessage() }
    }

    // Chat Operations
    override suspend fun getChatSessions(): Flow<List<Chat>> {
        return database.chatQueries.getAllChats()
            .asFlow()
            .mapToList()
            .map { dbChats ->
                dbChats.map { it.toDomainChat() }
            }
    }

    override suspend fun getChatSession(chatId: String): Result<Chat> {
        return try {
            val dbChat = database.chatQueries.getChatById(chatId).executeAsOneOrNull()
            if (dbChat != null) {
                Result.success(dbChat.toDomainChat())
            } else {
                Result.failure(ChatNotFoundException(chatId))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun saveChatSession(chat: Chat): Result<Unit> {
        return try {
            database.chatQueries.insertChat(
                id = chat.id,
                name = chat.name,
                type = chat.type.name,
                participants = chat.participants.toJsonString(),
                lastMessageId = chat.lastMessage?.id,
                unreadCount = chat.unreadCount.toLong(),
                isArchived = if (chat.isArchived) 1L else 0L,
                isMuted = if (chat.isMuted) 1L else 0L,
                isPinned = if (chat.isPinned) 1L else 0L,
                metadata = chat.metadata.toJsonString(),
                createdAt = chat.createdAt.toEpochMilliseconds(),
                updatedAt = chat.updatedAt.toEpochMilliseconds(),
                lastActivityAt = chat.lastActivityAt?.toEpochMilliseconds()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateLastActivity(chatId: String, timestamp: Instant): Result<Unit> {
        return try {
            database.chatQueries.updateLastActivity(
                lastActivityAt = timestamp.toEpochMilliseconds(),
                updatedAt = Clock.System.now().toEpochMilliseconds(),
                id = chatId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Sync Metadata
    override suspend fun getLastSyncTimestamp(chatId: String): Instant? {
        return database.syncMetadataQueries.getLastSyncTimestamp(chatId)
            .executeAsOneOrNull()
            ?.let { Instant.fromEpochMilliseconds(it) }
    }

    override suspend fun updateSyncTimestamp(chatId: String, timestamp: Instant): Result<Unit> {
        return try {
            database.syncMetadataQueries.updateSyncTimestamp(
                chatId = chatId,
                lastSyncTimestamp = timestamp.toEpochMilliseconds(),
                updatedAt = Clock.System.now().toEpochMilliseconds()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getPendingOperations(): List<SyncOperation> {
        return database.syncOperationQueries.getPendingOperations()
            .executeAsList()
            .map { it.toDomainSyncOperation() }
    }

    override suspend fun markOperationCompleted(operationId: String): Result<Unit> {
        return try {
            database.syncOperationQueries.markOperationCompleted(
                status = SyncOperationStatus.COMPLETED.name,
                completedAt = Clock.System.now().toEpochMilliseconds(),
                id = operationId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
```

### 1.2 Database Schema Extensions
```sql
-- File: database/schema/sync.sq

-- Add sync metadata table
CREATE TABLE sync_metadata (
    chat_id TEXT PRIMARY KEY,
    last_sync_timestamp INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- Add sync operations table
CREATE TABLE sync_operations (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    data TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    status TEXT DEFAULT 'PENDING',
    created_at INTEGER NOT NULL,
    completed_at INTEGER
);

-- Update messages table to include sync status
ALTER TABLE messages ADD COLUMN sync_status TEXT DEFAULT 'PENDING';
ALTER TABLE messages ADD COLUMN last_sync_at INTEGER;
ALTER TABLE messages ADD COLUMN server_timestamp INTEGER;

-- Indexes for performance
CREATE INDEX idx_messages_sync_status ON messages(sync_status);
CREATE INDEX idx_messages_chat_sync ON messages(chat_id, sync_status);
CREATE INDEX idx_sync_operations_status ON sync_operations(status);
CREATE INDEX idx_sync_operations_chat ON sync_operations(chat_id);
```

## Step 2: Create RemoteDataSource Implementation

### 2.1 ApiRemoteDataSource
```kotlin
// File: repositories/datasource/ApiRemoteDataSource.kt
package com.tchat.mobile.repositories.datasource

import com.tchat.mobile.api.ChatApiService
import com.tchat.mobile.api.WebSocketService
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.datetime.Instant

class ApiRemoteDataSource(
    private val chatApiService: ChatApiService,
    private val webSocketService: WebSocketService
) : RemoteDataSource {

    // Message Operations
    override suspend fun fetchMessages(
        chatId: String,
        since: Instant?,
        limit: Int
    ): Result<List<ChatMessage>> {
        return try {
            val response = chatApiService.getMessages(
                chatId = chatId,
                since = since?.toEpochMilliseconds(),
                limit = limit
            )

            if (response.isSuccessful) {
                val apiMessages = response.body()?.data ?: emptyList()
                val domainMessages = apiMessages.map { it.toDomainMessage() }
                Result.success(domainMessages)
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun sendMessage(message: ChatMessage): Result<ChatMessage> {
        return try {
            val apiMessage = message.toApiMessage()
            val response = chatApiService.sendMessage(apiMessage)

            if (response.isSuccessful) {
                val serverMessage = response.body()?.data?.toDomainMessage()
                if (serverMessage != null) {
                    Result.success(serverMessage)
                } else {
                    Result.failure(ApiException("Invalid server response"))
                }
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateMessage(
        messageId: String,
        content: MessageContent
    ): Result<ChatMessage> {
        return try {
            val response = chatApiService.updateMessage(
                messageId = messageId,
                content = content.toApiContent()
            )

            if (response.isSuccessful) {
                val updatedMessage = response.body()?.data?.toDomainMessage()
                if (updatedMessage != null) {
                    Result.success(updatedMessage)
                } else {
                    Result.failure(ApiException("Invalid server response"))
                }
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteMessage(messageId: String): Result<Unit> {
        return try {
            val response = chatApiService.deleteMessage(messageId)

            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Real-time Subscriptions
    override fun subscribeToMessages(chatId: String): Flow<ChatMessage> = flow {
        webSocketService.subscribeToChat(chatId) { apiMessage ->
            emit(apiMessage.toDomainMessage())
        }
    }

    override fun subscribeToTypingIndicators(chatId: String): Flow<TypingIndicator> = flow {
        webSocketService.subscribeToTyping(chatId) { apiTyping ->
            emit(apiTyping.toDomainTyping())
        }
    }

    override fun subscribeToPresence(chatId: String): Flow<PresenceUpdate> = flow {
        webSocketService.subscribeToPresence(chatId) { apiPresence ->
            emit(apiPresence.toDomainPresence())
        }
    }

    // Chat Operations
    override suspend fun fetchChatSessions(): Result<List<Chat>> {
        return try {
            val response = chatApiService.getChats()

            if (response.isSuccessful) {
                val apiChats = response.body()?.data ?: emptyList()
                val domainChats = apiChats.map { it.toDomainChat() }
                Result.success(domainChats)
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun createChatSession(chat: Chat): Result<Chat> {
        return try {
            val apiChat = chat.toApiChat()
            val response = chatApiService.createChat(apiChat)

            if (response.isSuccessful) {
                val createdChat = response.body()?.data?.toDomainChat()
                if (createdChat != null) {
                    Result.success(createdChat)
                } else {
                    Result.failure(ApiException("Invalid server response"))
                }
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateChatSession(chatId: String, updates: ChatUpdates): Result<Chat> {
        return try {
            val response = chatApiService.updateChat(chatId, updates.toApiUpdates())

            if (response.isSuccessful) {
                val updatedChat = response.body()?.data?.toDomainChat()
                if (updatedChat != null) {
                    Result.success(updatedChat)
                } else {
                    Result.failure(ApiException("Invalid server response"))
                }
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun joinChatSession(chatId: String): Result<Unit> {
        return try {
            val response = chatApiService.joinChat(chatId)

            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun leaveChatSession(chatId: String): Result<Unit> {
        return try {
            val response = chatApiService.leaveChat(chatId)

            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(ApiException(response.code(), response.message()))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Connection Management
    override suspend fun connect(): Result<Unit> {
        return try {
            webSocketService.connect()
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun disconnect(): Result<Unit> {
        return try {
            webSocketService.disconnect()
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override fun getConnectionState(): Flow<ConnectionState> {
        return webSocketService.connectionState
    }
}
```

## Step 3: Create SyncEngine Implementation

### 3.1 ChatSyncEngine
```kotlin
// File: repositories/sync/ChatSyncEngine.kt
package com.tchat.mobile.repositories.sync

import com.tchat.mobile.repositories.datasource.LocalDataSource
import com.tchat.mobile.repositories.datasource.RemoteDataSource
import com.tchat.mobile.repositories.sync.ConflictResolutionEngine
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.datetime.Clock

class ChatSyncEngine(
    private val localDataSource: LocalDataSource,
    private val remoteDataSource: RemoteDataSource,
    private val conflictResolver: ConflictResolutionEngine,
    private val circuitBreaker: ApiCircuitBreaker
) : SyncEngine {

    private val _syncState = MutableStateFlow<Map<String, SyncState>>(emptyMap())
    private val _globalSyncState = MutableStateFlow(
        GlobalSyncState(
            isConnected = false,
            activeSyncs = 0,
            pendingOperations = 0,
            lastSuccessfulSync = null,
            networkState = NetworkState.DISCONNECTED
        )
    )

    override suspend fun syncChat(
        chatId: String,
        strategy: SyncStrategy
    ): Result<SyncResult> {
        updateSyncState(chatId, SyncState.SYNCING)

        return try {
            when (strategy) {
                SyncStrategy.FULL_REFRESH -> performFullRefresh(chatId)
                SyncStrategy.INCREMENTAL -> performIncrementalSync(chatId)
                SyncStrategy.REAL_TIME_ONLY -> performRealTimeSync(chatId)
            }
        } catch (e: Exception) {
            updateSyncState(chatId, SyncState.ERROR)
            Result.failure(e)
        }
    }

    private suspend fun performIncrementalSync(chatId: String): Result<SyncResult> {
        val startTime = Clock.System.now()

        // 1. Get last sync point
        val lastSync = localDataSource.getLastSyncTimestamp(chatId)

        // 2. Fetch remote changes since last sync
        val remoteChanges = circuitBreaker.execute {
            remoteDataSource.fetchMessages(chatId, since = lastSync)
        }.getOrElse {
            return Result.failure(it)
        }

        // 3. Get local pending changes
        val localPending = localDataSource.getPendingSyncMessages()
            .filter { it.chatId == chatId }

        // 4. Detect conflicts
        val conflicts = detectConflicts(localPending, remoteChanges)

        // 5. Resolve conflicts
        val resolutions = conflicts.map { conflict ->
            conflictResolver.resolveConflict(conflict)
        }

        // 6. Apply remote changes (non-conflicted)
        val nonConflictedChanges = remoteChanges.filter { remote ->
            conflicts.none { it.remoteMessage.id == remote.id }
        }

        nonConflictedChanges.forEach { message ->
            localDataSource.saveMessage(message)
            localDataSource.markMessageAsSynced(message.id, message.timestamp)
        }

        // 7. Apply conflict resolutions
        resolutions.forEach { resolution ->
            applyConflictResolution(resolution)
        }

        // 8. Push local changes to remote
        val syncFailures = mutableListOf<String>()
        localPending.forEach { localMessage ->
            circuitBreaker.execute {
                remoteDataSource.sendMessage(localMessage)
            }.onSuccess { serverMessage ->
                localDataSource.markMessageAsSynced(serverMessage.id, serverMessage.timestamp)
            }.onFailure { error ->
                syncFailures.add("${localMessage.id}: ${error.message}")
            }
        }

        // 9. Update sync timestamp
        val syncTime = Clock.System.now()
        localDataSource.updateSyncTimestamp(chatId, syncTime)

        updateSyncState(chatId, if (syncFailures.isEmpty()) SyncState.SYNCED else SyncState.ERROR)

        return Result.success(
            SyncResult(
                chatId = chatId,
                success = syncFailures.isEmpty(),
                messagesUpdated = remoteChanges.size + resolutions.size,
                conflicts = conflicts,
                lastSyncTimestamp = syncTime,
                syncDuration = syncTime - startTime,
                errors = syncFailures
            )
        )
    }

    private suspend fun detectConflicts(
        localMessages: List<ChatMessage>,
        remoteMessages: List<ChatMessage>
    ): List<MessageConflict> {
        val conflicts = mutableListOf<MessageConflict>()

        localMessages.forEach { local ->
            val remote = remoteMessages.find { it.id == local.id }
            if (remote != null) {
                val conflict = analyzeConflict(local, remote)
                if (conflict != null) {
                    conflicts.add(conflict)
                }
            }
        }

        return conflicts
    }

    private fun analyzeConflict(
        local: ChatMessage,
        remote: ChatMessage
    ): MessageConflict? {
        return when {
            local.timestamp != remote.timestamp -> MessageConflict(
                type = ConflictType.EDIT_CONFLICT,
                chatId = local.chatId,
                localMessage = local,
                remoteMessage = remote,
                severity = ConflictSeverity.MEDIUM,
                autoResolvable = true
            )
            local.status != remote.status -> MessageConflict(
                type = ConflictType.STATUS_CONFLICT,
                chatId = local.chatId,
                localMessage = local,
                remoteMessage = remote,
                severity = ConflictSeverity.LOW,
                autoResolvable = true
            )
            else -> null
        }
    }

    private suspend fun applyConflictResolution(resolution: ConflictResolution) {
        when (resolution.strategy) {
            ResolutionStrategy.LOCAL_WINS -> {
                // Push local version to remote
                resolution.resolvedMessage?.let { message ->
                    remoteDataSource.updateMessage(message.id, message.content)
                }
            }
            ResolutionStrategy.REMOTE_WINS -> {
                // Update local with remote version
                resolution.resolvedMessage?.let { message ->
                    localDataSource.saveMessage(message)
                    localDataSource.markMessageAsSynced(message.id, message.timestamp)
                }
            }
            ResolutionStrategy.MERGE -> {
                // Save merged version locally and push to remote
                resolution.resolvedMessage?.let { message ->
                    localDataSource.saveMessage(message)
                    remoteDataSource.updateMessage(message.id, message.content)
                    localDataSource.markMessageAsSynced(message.id, message.timestamp)
                }
            }
        }
    }

    override suspend fun syncAllChats(): Result<List<SyncResult>> {
        val chatSessions = localDataSource.getChatSessions().first()
        val results = mutableListOf<SyncResult>()

        chatSessions.forEach { chat ->
            syncChat(chat.id, SyncStrategy.INCREMENTAL)
                .onSuccess { result -> results.add(result) }
                .onFailure { error ->
                    results.add(
                        SyncResult(
                            chatId = chat.id,
                            success = false,
                            messagesUpdated = 0,
                            conflicts = emptyList(),
                            lastSyncTimestamp = Clock.System.now(),
                            errors = listOf(error.message ?: "Unknown error")
                        )
                    )
                }
        }

        return Result.success(results)
    }

    override suspend fun forcePushPendingOperations(): Result<Unit> {
        val pendingOps = localDataSource.getPendingOperations()

        pendingOps.forEach { operation ->
            processOperation(operation)
        }

        return Result.success(Unit)
    }

    private suspend fun processOperation(operation: SyncOperation) {
        when (operation.type) {
            OperationType.SEND_MESSAGE -> {
                val message = operation.data.fromJson<ChatMessage>()
                remoteDataSource.sendMessage(message)
                    .onSuccess { serverMessage ->
                        localDataSource.markOperationCompleted(operation.id)
                        localDataSource.markMessageAsSynced(serverMessage.id, serverMessage.timestamp)
                    }
                    .onFailure { error ->
                        handleOperationFailure(operation, error)
                    }
            }
            // Handle other operation types...
        }
    }

    private suspend fun handleOperationFailure(operation: SyncOperation, error: Throwable) {
        if (operation.retryCount < operation.maxRetries) {
            // Schedule retry with exponential backoff
            val delay = minOf(1000L * (1 shl operation.retryCount), 30000L)
            kotlinx.coroutines.delay(delay)

            val retryOperation = operation.copy(retryCount = operation.retryCount + 1)
            processOperation(retryOperation)
        } else {
            // Max retries reached - mark as failed
            localDataSource.markOperationFailed(operation.id, error.message)
        }
    }

    override suspend fun resolveConflicts(conflicts: List<DataConflict>): Result<List<ConflictResolution>> {
        val resolutions = conflicts.map { conflict ->
            conflictResolver.resolveConflict(conflict as MessageConflict)
        }
        return Result.success(resolutions)
    }

    override fun getSyncState(chatId: String): Flow<SyncState> {
        return _syncState.asStateFlow().map { states ->
            states[chatId] ?: SyncState.IDLE
        }
    }

    override fun getGlobalSyncState(): Flow<GlobalSyncState> {
        return _globalSyncState.asStateFlow()
    }

    override suspend fun startBackgroundSync() {
        // Implementation for background sync worker
    }

    override suspend fun stopBackgroundSync() {
        // Implementation to stop background sync
    }

    private fun updateSyncState(chatId: String, state: SyncState) {
        _syncState.value = _syncState.value + (chatId to state)
    }
}
```

## Step 4: Refactor ChatRepository

### 4.1 New ChatRepository Implementation
```kotlin
// File: repositories/ChatRepositoryImpl.kt
package com.tchat.mobile.repositories

import com.tchat.mobile.repositories.datasource.LocalDataSource
import com.tchat.mobile.repositories.datasource.RemoteDataSource
import com.tchat.mobile.repositories.sync.SyncEngine
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.combine
import kotlinx.datetime.Clock

class ChatRepositoryImpl(
    private val localDataSource: LocalDataSource,
    private val remoteDataSource: RemoteDataSource,
    private val syncEngine: SyncEngine
) : ChatRepository {

    // Message Operations (UI-facing)
    override fun getMessages(chatId: String): Flow<List<ChatMessage>> {
        return combine(
            localDataSource.getMessages(chatId),
            syncEngine.getSyncState(chatId)
        ) { messages, syncState ->
            // Trigger background refresh if needed
            if (shouldRefresh(chatId, syncState)) {
                kotlinx.coroutines.GlobalScope.launch {
                    syncEngine.syncChat(chatId, SyncStrategy.INCREMENTAL)
                }
            }
            messages
        }
    }

    override suspend fun sendMessage(
        chatId: String,
        content: MessageContent,
        replyToId: String?
    ): Result<ChatMessage> {
        val message = ChatMessage(
            id = generateMessageId(),
            chatId = chatId,
            senderId = getCurrentUserId(),
            senderInfo = getCurrentUserInfo(),
            type = determineMessageType(content),
            content = content,
            timestamp = Clock.System.now(),
            status = MessageStatus.SENDING,
            replyTo = replyToId
        )

        // 1. Optimistic local update
        localDataSource.saveMessage(message.copy(status = MessageStatus.SENDING))
            .onFailure { return Result.failure(it) }

        // 2. Send to remote (Write-through)
        return remoteDataSource.sendMessage(message)
            .onSuccess { serverMessage ->
                // 3. Update with server response
                localDataSource.markMessageAsSynced(serverMessage.id, serverMessage.timestamp)
            }
            .onFailure { error ->
                // 4. Mark as failed for retry
                localDataSource.updateMessageStatus(message.id, MessageStatus.FAILED)

                // 5. Queue for background retry
                queueSyncOperation(
                    SyncOperation(
                        id = generateOperationId(),
                        type = OperationType.SEND_MESSAGE,
                        chatId = chatId,
                        data = message.toJsonString(),
                        timestamp = Clock.System.now()
                    )
                )
            }
    }

    override suspend fun editMessage(
        messageId: String,
        newContent: MessageContent
    ): Result<ChatMessage> {
        // Write-through pattern for edits
        return remoteDataSource.updateMessage(messageId, newContent)
            .onSuccess { updatedMessage ->
                localDataSource.saveMessage(updatedMessage)
                localDataSource.markMessageAsSynced(updatedMessage.id, updatedMessage.timestamp)
            }
            .onFailure { error ->
                // Queue for retry
                queueSyncOperation(
                    SyncOperation(
                        id = generateOperationId(),
                        type = OperationType.EDIT_MESSAGE,
                        chatId = getChatIdForMessage(messageId),
                        data = EditMessageData(messageId, newContent).toJsonString(),
                        timestamp = Clock.System.now()
                    )
                )
            }
    }

    override suspend fun deleteMessage(messageId: String): Result<Unit> {
        // Immediate local update + remote sync
        localDataSource.updateMessageStatus(messageId, MessageStatus.DELETED)
            .onFailure { return Result.failure(it) }

        return remoteDataSource.deleteMessage(messageId)
            .onFailure { error ->
                // Revert local change on failure
                localDataSource.updateMessageStatus(messageId, MessageStatus.SENT)

                // Queue for retry
                queueSyncOperation(
                    SyncOperation(
                        id = generateOperationId(),
                        type = OperationType.DELETE_MESSAGE,
                        chatId = getChatIdForMessage(messageId),
                        data = DeleteMessageData(messageId).toJsonString(),
                        timestamp = Clock.System.now()
                    )
                )
            }
    }

    override suspend fun searchMessages(chatId: String, query: String): Result<List<ChatMessage>> {
        // Cache-aside pattern for search
        return try {
            // 1. Try local search first
            val localResults = localDataSource.searchMessages(chatId, query)

            // 2. If insufficient results, try remote
            if (localResults.size < 10) {
                remoteDataSource.searchMessages(chatId, query)
                    .onSuccess { remoteResults ->
                        // Cache remote results
                        remoteResults.forEach { message ->
                            localDataSource.saveMessage(message)
                            localDataSource.markMessageAsSynced(message.id, message.timestamp)
                        }
                    }
                    .getOrElse { localResults }
            } else {
                localResults
            }

            Result.success(localResults)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Chat Operations (UI-facing)
    override fun getChatSessions(): Flow<List<Chat>> {
        return combine(
            localDataSource.getChatSessions(),
            syncEngine.getGlobalSyncState()
        ) { chats, globalState ->
            // Trigger periodic refresh
            if (shouldRefreshChats(globalState)) {
                kotlinx.coroutines.GlobalScope.launch {
                    syncEngine.syncAllChats()
                }
            }
            chats
        }
    }

    override fun getChatSession(chatId: String): Flow<Chat?> {
        return localDataSource.getChatSession(chatId).asFlow()
    }

    override suspend fun createChatSession(chat: Chat): Result<Chat> {
        // Write-through pattern for chat creation
        return remoteDataSource.createChatSession(chat)
            .onSuccess { serverChat ->
                localDataSource.saveChatSession(serverChat)
            }
    }

    override suspend fun updateChatSession(chatId: String, updates: ChatUpdates): Result<Chat> {
        // Write-through pattern for updates
        return remoteDataSource.updateChatSession(chatId, updates)
            .onSuccess { updatedChat ->
                localDataSource.saveChatSession(updatedChat)
            }
    }

    // Real-time Operations
    override fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>> {
        return remoteDataSource.subscribeToTypingIndicators(chatId)
            .scan(emptyList<TypingIndicator>()) { current, indicator ->
                // Manage typing indicator lifecycle
                current.filter { it.userId != indicator.userId && !it.isExpired() } + indicator
            }
    }

    override suspend fun sendTypingIndicator(chatId: String, type: TypingType): Result<Unit> {
        // Write-behind pattern for typing indicators
        localDataSource.updateTypingStatus(chatId, getCurrentUserId(), type)

        return remoteDataSource.sendTypingIndicator(chatId, type)
            .onFailure { error ->
                // Non-critical, don't queue for retry
                println("Failed to send typing indicator: ${error.message}")
            }
    }

    // Sync Management
    override suspend fun refreshChat(chatId: String): Result<Unit> {
        return syncEngine.syncChat(chatId, SyncStrategy.INCREMENTAL)
            .map { Unit }
    }

    override suspend fun refreshAllChats(): Result<Unit> {
        return syncEngine.syncAllChats()
            .map { Unit }
    }

    override fun getSyncStatus(chatId: String): Flow<SyncStatus> {
        return combine(
            syncEngine.getSyncState(chatId),
            localDataSource.getLastSyncTimestamp(chatId).asFlow(),
            localDataSource.getPendingOperations().asFlow()
        ) { state, lastSync, pendingOps ->
            SyncStatus(
                state = state,
                lastSyncTime = lastSync,
                pendingOperations = pendingOps.count { it.chatId == chatId },
                conflicts = emptyList() // TODO: Get from conflict resolver
            )
        }
    }

    // Helper methods
    private fun shouldRefresh(chatId: String, syncState: SyncState): Boolean {
        return syncState == SyncState.IDLE ||
               (Clock.System.now() - getLastRefreshTime(chatId)) > REFRESH_INTERVAL
    }

    private fun shouldRefreshChats(globalState: GlobalSyncState): Boolean {
        return globalState.lastSuccessfulSync?.let { lastSync ->
            (Clock.System.now() - lastSync) > CHAT_REFRESH_INTERVAL
        } ?: true
    }

    private suspend fun queueSyncOperation(operation: SyncOperation) {
        localDataSource.saveSyncOperation(operation)
    }

    private fun generateMessageId(): String = "msg_${Clock.System.now().toEpochMilliseconds()}"
    private fun generateOperationId(): String = "op_${Clock.System.now().toEpochMilliseconds()}"
    private fun getCurrentUserId(): String = "current_user" // TODO: Get from auth service
    private fun getCurrentUserInfo(): UserInfo = UserInfo("current_user", "you", "You")

    companion object {
        private val REFRESH_INTERVAL = 5.minutes
        private val CHAT_REFRESH_INTERVAL = 30.minutes
    }
}
```

## Step 5: Dependency Injection Setup

### 5.1 Koin Module Configuration
```kotlin
// File: di/RepositoryModule.kt
package com.tchat.mobile.di

import com.tchat.mobile.repositories.*
import com.tchat.mobile.repositories.datasource.*
import com.tchat.mobile.repositories.sync.*
import org.koin.dsl.module

val repositoryModule = module {

    // Data Sources
    single<LocalDataSource> {
        SQLDelightLocalDataSource(get())
    }

    single<RemoteDataSource> {
        ApiRemoteDataSource(get(), get())
    }

    // Conflict Resolution
    single { ConflictDetector(get(), get()) }
    single { LastWriterWinsResolver() }
    single { StatusMergeResolver() }
    single { ContentMergeResolver() }
    single { UserChoiceResolver(get()) }
    single { ConflictResolutionEngine(get(), get(), get(), get()) }

    // Circuit Breaker
    single { ApiCircuitBreaker() }

    // Sync Engine
    single<SyncEngine> {
        ChatSyncEngine(get(), get(), get(), get())
    }

    // Repository
    single<ChatRepository> {
        ChatRepositoryImpl(get(), get(), get())
    }
}
```

## Step 6: Migration Strategy

### 6.1 Gradual Migration Plan

**Phase 1: Foundation Setup**
1. Create new interface definitions
2. Implement SQLDelightLocalDataSource
3. Add database schema changes
4. Set up basic dependency injection

**Phase 2: Remote Integration**
5. Implement ApiRemoteDataSource
6. Create basic SyncEngine
7. Set up WebSocket connections

**Phase 3: Repository Refactoring**
8. Create new ChatRepositoryImpl
9. Maintain backward compatibility with existing ChatRepository interface
10. Test with feature flags

**Phase 4: Advanced Features**
11. Implement conflict resolution
12. Add background sync
13. Performance optimization
14. Remove legacy MockChatRepository

### 6.2 Feature Flag Implementation
```kotlin
// File: utils/FeatureFlags.kt
object FeatureFlags {
    const val USE_NEW_REPOSITORY = "use_new_repository"
    const val ENABLE_CONFLICT_RESOLUTION = "enable_conflict_resolution"
    const val ENABLE_BACKGROUND_SYNC = "enable_background_sync"
}

// File: di/RepositoryModule.kt - Updated
val repositoryModule = module {
    single<ChatRepository> {
        if (getProperty<Boolean>(FeatureFlags.USE_NEW_REPOSITORY) == true) {
            ChatRepositoryImpl(get(), get(), get())
        } else {
            MockChatRepository.getInstance()
        }
    }
}
```

## Benefits of New Architecture

### 1. **Separation of Concerns**
- LocalDataSource: Pure SQLDelight operations
- RemoteDataSource: Pure API operations
- SyncEngine: Coordination and conflict resolution
- Repository: Business logic and UI coordination

### 2. **Testability**
- Each component can be unit tested independently
- Easy to mock data sources for testing
- Clear boundaries for integration tests

### 3. **Offline Support**
- Cache-aside pattern for reads
- Write-through/Write-behind for updates
- Automatic conflict detection and resolution

### 4. **Real-time Capabilities**
- WebSocket integration
- Optimistic updates
- Background synchronization

### 5. **Reliability**
- Circuit breaker pattern for API failures
- Exponential backoff retry logic
- Comprehensive error handling

### 6. **Performance**
- Local-first approach
- Intelligent refresh strategies
- Background sync operations

This implementation guide provides a complete migration path from the current monolithic repository to a sophisticated, scalable architecture that separates local and remote concerns while maintaining the existing ChatRepository interface for backward compatibility.