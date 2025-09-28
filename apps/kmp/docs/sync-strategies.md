# SQLDelight ↔ API Synchronization Strategies

## 1. Message Synchronization Patterns

### Write-Through Cache (Critical Messages)
```kotlin
suspend fun sendMessage(message: ChatMessage): Result<ChatMessage> {
    // 1. Optimistic local update for immediate UI feedback
    val pendingMessage = message.copy(status = MessageStatus.SENDING)
    localDataSource.saveMessage(pendingMessage)

    // 2. Send to API (Write-through)
    return remoteDataSource.sendMessage(message)
        .onSuccess { serverMessage ->
            // 3. Update with server response
            localDataSource.markMessageAsSynced(serverMessage.id, serverMessage.timestamp)
        }
        .onFailure { error ->
            // 4. Mark as failed for retry
            localDataSource.updateMessageStatus(message.id, MessageStatus.FAILED)
            // 5. Queue for background retry
            syncEngine.queueOperation(
                SyncOperation(
                    id = UUID.randomUUID().toString(),
                    type = OperationType.SEND_MESSAGE,
                    chatId = message.chatId,
                    data = message.toJsonString(),
                    timestamp = Clock.System.now()
                )
            )
        }
}
```

### Write-Behind Cache (Bulk Operations)
```kotlin
suspend fun markMessagesAsRead(messageIds: List<String>): Result<Unit> {
    // 1. Immediate local update
    messageIds.forEach { messageId ->
        localDataSource.updateMessageStatus(messageId, MessageStatus.READ)
    }

    // 2. Queue for background sync
    syncEngine.queueOperation(
        SyncOperation(
            id = UUID.randomUUID().toString(),
            type = OperationType.MARK_READ,
            chatId = chatId,
            data = messageIds.toJsonString(),
            timestamp = Clock.System.now()
        )
    )

    return Result.success(Unit)
}
```

### Cache-Aside (Read Operations)
```kotlin
fun getMessages(chatId: String): Flow<List<ChatMessage>> {
    return flow {
        // 1. Emit local data immediately
        emit(localDataSource.getMessages(chatId).first())

        // 2. Check if refresh needed
        val lastSync = localDataSource.getLastSyncTimestamp(chatId)
        val shouldRefresh = lastSync == null ||
            (Clock.System.now() - lastSync) > CACHE_EXPIRY_DURATION

        if (shouldRefresh) {
            // 3. Background refresh from API
            refreshFromRemote(chatId)
        }
    }.distinctUntilChanged()
}

private suspend fun refreshFromRemote(chatId: String) {
    val lastSync = localDataSource.getLastSyncTimestamp(chatId)

    remoteDataSource.fetchMessages(chatId, since = lastSync)
        .onSuccess { messages ->
            messages.forEach { message ->
                localDataSource.saveMessage(message)
            }
            localDataSource.updateSyncTimestamp(chatId, Clock.System.now())
        }
}
```

## 2. Real-Time Synchronization

### WebSocket Message Handling
```kotlin
class RealTimeSyncEngine(
    private val localDataSource: LocalDataSource,
    private val remoteDataSource: RemoteDataSource
) {
    private val _connectionState = MutableStateFlow(ConnectionState.DISCONNECTED)
    val connectionState: StateFlow<ConnectionState> = _connectionState.asStateFlow()

    suspend fun startRealTimeSync() {
        remoteDataSource.connect()
            .onSuccess {
                _connectionState.value = ConnectionState.CONNECTED

                // Subscribe to real-time message updates
                subscribeToMessages()
                subscribeToTypingIndicators()
                subscribeToPresenceUpdates()
            }
    }

    private suspend fun subscribeToMessages() {
        remoteDataSource.subscribeToMessages("*") // All chats user has access to
            .catch { error ->
                _connectionState.value = ConnectionState.ERROR
                // Implement exponential backoff retry
                retryConnection()
            }
            .collect { message ->
                // Real-time message received
                handleIncomingMessage(message)
            }
    }

    private suspend fun handleIncomingMessage(message: ChatMessage) {
        // 1. Check if message already exists locally
        val existingMessage = localDataSource.getMessage(message.id)

        if (existingMessage == null) {
            // 2. New message - save directly
            localDataSource.saveMessage(message)
        } else if (existingMessage.timestamp < message.timestamp) {
            // 3. Updated message - merge changes
            localDataSource.saveMessage(message)
        }
        // 4. Ignore if local version is newer (shouldn't happen normally)
    }
}
```

## 3. Background Sync Operations

### Retry with Exponential Backoff
```kotlin
class BackgroundSyncEngine {
    private val retryDelays = listOf(1000L, 2000L, 4000L, 8000L, 16000L) // milliseconds

    suspend fun processPendingOperations() {
        val pendingOps = localDataSource.getPendingOperations()

        pendingOps.forEach { operation ->
            processOperation(operation)
        }
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

            OperationType.EDIT_MESSAGE -> {
                // Handle message editing sync
            }

            OperationType.DELETE_MESSAGE -> {
                // Handle message deletion sync
            }
        }
    }

    private suspend fun handleOperationFailure(operation: SyncOperation, error: Throwable) {
        if (operation.retryCount < operation.maxRetries) {
            // Schedule retry with exponential backoff
            val delay = retryDelays.getOrElse(operation.retryCount) { retryDelays.last() }

            delay(delay)

            val retryOperation = operation.copy(retryCount = operation.retryCount + 1)
            localDataSource.updateOperation(retryOperation)

            processOperation(retryOperation)
        } else {
            // Max retries reached - mark as failed
            localDataSource.markOperationFailed(operation.id, error.message)

            // Notify user of permanent failure
            notificationService.showSyncFailure(operation)
        }
    }
}
```

## 4. Conflict Resolution Strategies

### Last Writer Wins (Simple Cases)
```kotlin
suspend fun resolveMessageEditConflict(
    localMessage: ChatMessage,
    remoteMessage: ChatMessage
): ConflictResolution {
    return if (remoteMessage.editedAt > localMessage.editedAt) {
        ConflictResolution(
            conflictId = "msg-edit-${localMessage.id}",
            resolution = ResolutionStrategy.REMOTE_WINS,
            resolvedData = remoteMessage
        )
    } else {
        ConflictResolution(
            conflictId = "msg-edit-${localMessage.id}",
            resolution = ResolutionStrategy.LOCAL_WINS,
            resolvedData = localMessage
        )
    }
}
```

### Operational Transform (Complex Cases)
```kotlin
suspend fun resolveRichTextConflict(
    localMessage: ChatMessage,
    remoteMessage: ChatMessage
): ConflictResolution {
    // For rich text editing conflicts
    val baseContent = getBaseContent(localMessage.id)
    val localOps = computeOperations(baseContent, localMessage.content)
    val remoteOps = computeOperations(baseContent, remoteMessage.content)

    val transformedOps = operationalTransform(localOps, remoteOps)
    val mergedContent = applyOperations(baseContent, transformedOps)

    return ConflictResolution(
        conflictId = "rich-text-${localMessage.id}",
        resolution = ResolutionStrategy.MERGE,
        resolvedData = localMessage.copy(content = mergedContent)
    )
}
```

### User Choice Required
```kotlin
suspend fun handleComplexConflict(conflict: DataConflict): ConflictResolution {
    // Show conflict resolution UI to user
    val userChoice = conflictResolutionUI.showConflictDialog(conflict)

    return when (userChoice) {
        UserChoice.KEEP_LOCAL -> ConflictResolution(
            conflictId = conflict.id,
            resolution = ResolutionStrategy.LOCAL_WINS,
            resolvedData = conflict.localData
        )

        UserChoice.KEEP_REMOTE -> ConflictResolution(
            conflictId = conflict.id,
            resolution = ResolutionStrategy.REMOTE_WINS,
            resolvedData = conflict.remoteData
        )

        UserChoice.MERGE_MANUAL -> {
            val mergedData = conflictResolutionUI.showMergeEditor(conflict)
            ConflictResolution(
                conflictId = conflict.id,
                resolution = ResolutionStrategy.MERGE,
                resolvedData = mergedData
            )
        }
    }
}
```

## 5. Performance Optimization Patterns

### Pagination & Lazy Loading
```kotlin
class MessagePaginator(
    private val localDataSource: LocalDataSource,
    private val remoteDataSource: RemoteDataSource
) {
    private val pageSize = 50

    suspend fun loadMessages(chatId: String, before: Instant? = null): Result<List<ChatMessage>> {
        // 1. Try local first
        val localMessages = localDataSource.getMessages(chatId, before, pageSize)

        if (localMessages.size >= pageSize) {
            // 2. We have enough local data
            return Result.success(localMessages)
        }

        // 3. Need to fetch from remote
        return remoteDataSource.fetchMessages(chatId, before = before, limit = pageSize)
            .onSuccess { remoteMessages ->
                // 4. Cache remote messages locally
                remoteMessages.forEach { message ->
                    localDataSource.saveMessage(message)
                }
            }
    }
}
```

### Delta Sync
```kotlin
suspend fun performDeltaSync(chatId: String): Result<SyncResult> {
    val lastSync = localDataSource.getLastSyncTimestamp(chatId)

    // 1. Fetch only changes since last sync
    return remoteDataSource.fetchMessagesSince(chatId, since = lastSync)
        .onSuccess { deltaMessages ->
            var updatedCount = 0

            deltaMessages.forEach { message ->
                val existing = localDataSource.getMessage(message.id)

                if (existing == null || existing.timestamp < message.timestamp) {
                    localDataSource.saveMessage(message)
                    updatedCount++
                }
            }

            localDataSource.updateSyncTimestamp(chatId, Clock.System.now())

            return Result.success(
                SyncResult(
                    chatId = chatId,
                    success = true,
                    messagesUpdated = updatedCount,
                    conflicts = emptyList(),
                    lastSyncTimestamp = Clock.System.now()
                )
            )
        }
}
```

## 6. Error Handling & Resilience

### Circuit Breaker Pattern
```kotlin
class ApiCircuitBreaker {
    private var state = CircuitState.CLOSED
    private var failureCount = 0
    private var lastFailureTime: Instant? = null

    private val failureThreshold = 5
    private val recoveryTimeout = 30.seconds

    suspend fun <T> execute(operation: suspend () -> Result<T>): Result<T> {
        return when (state) {
            CircuitState.CLOSED -> {
                operation().onFailure { error ->
                    failureCount++
                    lastFailureTime = Clock.System.now()

                    if (failureCount >= failureThreshold) {
                        state = CircuitState.OPEN
                    }
                }.onSuccess {
                    reset()
                }
            }

            CircuitState.OPEN -> {
                val now = Clock.System.now()
                if (lastFailureTime != null && (now - lastFailureTime!!) > recoveryTimeout) {
                    state = CircuitState.HALF_OPEN
                    execute(operation)
                } else {
                    Result.failure(CircuitBreakerOpenException("Circuit breaker is open"))
                }
            }

            CircuitState.HALF_OPEN -> {
                operation().onSuccess {
                    reset()
                }.onFailure {
                    state = CircuitState.OPEN
                    lastFailureTime = Clock.System.now()
                }
            }
        }
    }

    private fun reset() {
        state = CircuitState.CLOSED
        failureCount = 0
        lastFailureTime = null
    }
}

enum class CircuitState { CLOSED, OPEN, HALF_OPEN }
```

## 7. Monitoring & Observability

### Sync Metrics Collection
```kotlin
class SyncMetrics {
    private val _syncOperations = MutableStateFlow<List<SyncMetric>>(emptyList())
    val syncOperations: StateFlow<List<SyncMetric>> = _syncOperations.asStateFlow()

    fun recordSyncOperation(
        chatId: String,
        operation: String,
        duration: Duration,
        success: Boolean,
        error: String? = null
    ) {
        val metric = SyncMetric(
            chatId = chatId,
            operation = operation,
            duration = duration,
            success = success,
            error = error,
            timestamp = Clock.System.now()
        )

        _syncOperations.value = _syncOperations.value + metric

        // Send to analytics service
        analyticsService.track("sync_operation", metric.toMap())
    }

    fun getSyncHealthScore(chatId: String): Double {
        val recentMetrics = _syncOperations.value
            .filter { it.chatId == chatId }
            .filter { (Clock.System.now() - it.timestamp) < 1.hours }

        if (recentMetrics.isEmpty()) return 1.0

        val successRate = recentMetrics.count { it.success }.toDouble() / recentMetrics.size
        val avgDuration = recentMetrics.map { it.duration.inWholeMilliseconds }.average()

        // Health score based on success rate and performance
        return (successRate * 0.7) + ((1000 - avgDuration.coerceAtMost(1000.0)) / 1000 * 0.3)
    }
}

data class SyncMetric(
    val chatId: String,
    val operation: String,
    val duration: Duration,
    val success: Boolean,
    val error: String?,
    val timestamp: Instant
)
```