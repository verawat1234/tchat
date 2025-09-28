# SQLDelight ↔ API Conflict Resolution Patterns

## 1. Conflict Detection & Classification

### Temporal Conflict Detection
```kotlin
data class ConflictDetector(
    private val localDataSource: LocalDataSource,
    private val remoteDataSource: RemoteDataSource
) {
    suspend fun detectMessageConflicts(chatId: String): List<MessageConflict> {
        val localMessages = localDataSource.getUnsynced Messages(chatId)
        val remoteMessages = remoteDataSource.fetchMessages(chatId,
            since = getLastKnownSyncPoint(chatId))

        return localMessages.mapNotNull { localMsg ->
            val remoteMsg = remoteMessages.find { it.id == localMsg.id }

            when {
                remoteMsg == null -> null // No conflict, local message not on server yet
                localMsg.timestamp != remoteMsg.timestamp -> MessageConflict(
                    type = ConflictType.EDIT_CONFLICT,
                    localMessage = localMsg,
                    remoteMessage = remoteMsg,
                    severity = ConflictSeverity.MEDIUM
                )
                localMsg.status != remoteMsg.status -> MessageConflict(
                    type = ConflictType.STATUS_CONFLICT,
                    localMessage = localMsg,
                    remoteMessage = remoteMsg,
                    severity = ConflictSeverity.LOW
                )
                else -> null
            }
        }
    }

    private suspend fun getLastKnownSyncPoint(chatId: String): Instant? {
        return localDataSource.getLastSyncTimestamp(chatId)
    }
}
```

### Conflict Classification Matrix
```kotlin
sealed class ConflictType(val priority: Int) {
    object EDIT_CONFLICT : ConflictType(1)        // Same message edited on both sides
    object DELETE_CONFLICT : ConflictType(2)      // Message deleted locally but edited remotely
    object STATUS_CONFLICT : ConflictType(3)      // Read/delivered status mismatch
    object PARTICIPANT_CONFLICT : ConflictType(4) // Chat participant changes
    object METADATA_CONFLICT : ConflictType(5)    // Chat title, settings conflicts
}

enum class ConflictSeverity {
    CRITICAL,    // Data loss potential
    HIGH,        // User experience impact
    MEDIUM,      // Functional inconsistency
    LOW          // Cosmetic/status only
}

data class MessageConflict(
    val id: String = UUID.randomUUID().toString(),
    val type: ConflictType,
    val chatId: String,
    val localMessage: ChatMessage,
    val remoteMessage: ChatMessage,
    val severity: ConflictSeverity,
    val detectedAt: Instant = Clock.System.now(),
    val autoResolvable: Boolean = determineAutoResolvability(type, severity)
)
```

## 2. Automatic Resolution Strategies

### Last Writer Wins (LWW) - Simple Conflicts
```kotlin
class LastWriterWinsResolver : ConflictResolver {
    override suspend fun resolve(conflict: MessageConflict): ConflictResolution {
        val localTimestamp = conflict.localMessage.editedAt ?: conflict.localMessage.timestamp
        val remoteTimestamp = conflict.remoteMessage.editedAt ?: conflict.remoteMessage.timestamp

        return if (remoteTimestamp > localTimestamp) {
            ConflictResolution(
                conflictId = conflict.id,
                strategy = ResolutionStrategy.REMOTE_WINS,
                resolvedMessage = conflict.remoteMessage,
                explanation = "Remote version is newer (${remoteTimestamp} > ${localTimestamp})"
            )
        } else {
            ConflictResolution(
                conflictId = conflict.id,
                strategy = ResolutionStrategy.LOCAL_WINS,
                resolvedMessage = conflict.localMessage,
                explanation = "Local version is newer or equal"
            )
        }
    }
}
```

### Status Merge Strategy - Non-Destructive
```kotlin
class StatusMergeResolver : ConflictResolver {
    override suspend fun resolve(conflict: MessageConflict): ConflictResolution {
        return when (conflict.type) {
            ConflictType.STATUS_CONFLICT -> {
                val mergedStatus = mergeMessageStatus(
                    conflict.localMessage.status,
                    conflict.remoteMessage.status
                )

                ConflictResolution(
                    conflictId = conflict.id,
                    strategy = ResolutionStrategy.MERGE,
                    resolvedMessage = conflict.remoteMessage.copy(status = mergedStatus),
                    explanation = "Merged status: local=${conflict.localMessage.status}, remote=${conflict.remoteMessage.status}, result=${mergedStatus}"
                )
            }
            else -> ConflictResolution.unresolvable(conflict.id)
        }
    }

    private fun mergeMessageStatus(local: MessageStatus, remote: MessageStatus): MessageStatus {
        // Status progression: SENDING -> SENT -> DELIVERED -> READ
        return when {
            local == MessageStatus.READ || remote == MessageStatus.READ -> MessageStatus.READ
            local == MessageStatus.DELIVERED || remote == MessageStatus.DELIVERED -> MessageStatus.DELIVERED
            local == MessageStatus.SENT || remote == MessageStatus.SENT -> MessageStatus.SENT
            else -> maxOf(local, remote) // Use enum ordering
        }
    }
}
```

### Content Merge Strategy - Rich Text
```kotlin
class ContentMergeResolver : ConflictResolver {
    override suspend fun resolve(conflict: MessageConflict): ConflictResolution {
        return when {
            conflict.localMessage.content is MessageContent.Text &&
            conflict.remoteMessage.content is MessageContent.Text -> {
                val mergedContent = mergeTextContent(
                    conflict.localMessage.content as MessageContent.Text,
                    conflict.remoteMessage.content as MessageContent.Text
                )

                ConflictResolution(
                    conflictId = conflict.id,
                    strategy = ResolutionStrategy.MERGE,
                    resolvedMessage = conflict.remoteMessage.copy(content = mergedContent),
                    explanation = "Merged text content using operational transform"
                )
            }
            else -> ConflictResolution.unresolvable(conflict.id)
        }
    }

    private suspend fun mergeTextContent(
        local: MessageContent.Text,
        remote: MessageContent.Text
    ): MessageContent.Text {
        // Simple operational transform for text merging
        val baseText = getBaseContent(local, remote) // Get common ancestor
        val localOps = computeTextOperations(baseText, local.text)
        val remoteOps = computeTextOperations(baseText, remote.text)

        val transformedOps = operationalTransform(localOps, remoteOps)
        val mergedText = applyOperations(baseText, transformedOps)

        return MessageContent.Text(
            text = mergedText,
            formatting = mergeFormatting(local.formatting, remote.formatting)
        )
    }

    private fun mergeFormatting(
        local: List<MessageTextFormatting>,
        remote: List<MessageTextFormatting>
    ): List<MessageTextFormatting> {
        // Merge non-overlapping formatting
        val merged = mutableSetOf<MessageTextFormatting>()
        merged.addAll(local)
        merged.addAll(remote)
        return merged.toList().sortedBy { it.start }
    }
}
```

## 3. User-Mediated Resolution

### Conflict Resolution UI Integration
```kotlin
class UserChoiceResolver(
    private val conflictUI: ConflictResolutionUI
) : ConflictResolver {

    override suspend fun resolve(conflict: MessageConflict): ConflictResolution {
        return when (conflict.severity) {
            ConflictSeverity.CRITICAL, ConflictSeverity.HIGH -> {
                // Always require user input for critical conflicts
                promptUserResolution(conflict)
            }
            ConflictSeverity.MEDIUM -> {
                // Attempt auto-resolution, fallback to user
                tryAutoResolve(conflict) ?: promptUserResolution(conflict)
            }
            ConflictSeverity.LOW -> {
                // Auto-resolve low severity conflicts
                autoResolve(conflict)
            }
        }
    }

    private suspend fun promptUserResolution(conflict: MessageConflict): ConflictResolution {
        val userChoice = conflictUI.showConflictDialog(
            ConflictUIData(
                title = getConflictTitle(conflict.type),
                description = getConflictDescription(conflict),
                localPreview = formatMessagePreview(conflict.localMessage),
                remotePreview = formatMessagePreview(conflict.remoteMessage),
                options = getResolutionOptions(conflict.type)
            )
        )

        return when (userChoice.action) {
            UserResolutionAction.KEEP_LOCAL -> ConflictResolution(
                conflictId = conflict.id,
                strategy = ResolutionStrategy.LOCAL_WINS,
                resolvedMessage = conflict.localMessage,
                explanation = "User chose to keep local version"
            )

            UserResolutionAction.KEEP_REMOTE -> ConflictResolution(
                conflictId = conflict.id,
                strategy = ResolutionStrategy.REMOTE_WINS,
                resolvedMessage = conflict.remoteMessage,
                explanation = "User chose to keep remote version"
            )

            UserResolutionAction.MERGE_MANUAL -> {
                val editedMessage = conflictUI.showMergeEditor(
                    conflict.localMessage,
                    conflict.remoteMessage
                )
                ConflictResolution(
                    conflictId = conflict.id,
                    strategy = ResolutionStrategy.MERGE,
                    resolvedMessage = editedMessage,
                    explanation = "User manually merged the content"
                )
            }

            UserResolutionAction.DISCARD_BOTH -> ConflictResolution(
                conflictId = conflict.id,
                strategy = ResolutionStrategy.DISCARD,
                resolvedMessage = null,
                explanation = "User chose to discard both versions"
            )
        }
    }
}

data class ConflictUIData(
    val title: String,
    val description: String,
    val localPreview: String,
    val remotePreview: String,
    val options: List<ResolutionOption>
)

enum class UserResolutionAction {
    KEEP_LOCAL,
    KEEP_REMOTE,
    MERGE_MANUAL,
    DISCARD_BOTH
}
```

## 4. Resolution Strategy Selection

### Intelligent Strategy Router
```kotlin
class ConflictResolutionEngine(
    private val lastWriterWinsResolver: LastWriterWinsResolver,
    private val statusMergeResolver: StatusMergeResolver,
    private val contentMergeResolver: ContentMergeResolver,
    private val userChoiceResolver: UserChoiceResolver
) {

    suspend fun resolveConflict(conflict: MessageConflict): ConflictResolution {
        val strategy = selectResolutionStrategy(conflict)

        return try {
            val resolver = getResolver(strategy)
            val resolution = resolver.resolve(conflict)

            // Log resolution for learning
            logResolution(conflict, resolution)

            resolution
        } catch (e: Exception) {
            // Fallback to user choice if auto-resolution fails
            userChoiceResolver.resolve(conflict)
        }
    }

    private fun selectResolutionStrategy(conflict: MessageConflict): ResolutionStrategy {
        return when {
            // Critical conflicts always require user input
            conflict.severity == ConflictSeverity.CRITICAL -> ResolutionStrategy.USER_CHOICE_REQUIRED

            // Status conflicts can usually be merged safely
            conflict.type == ConflictType.STATUS_CONFLICT -> ResolutionStrategy.STATUS_MERGE

            // Edit conflicts on text content attempt content merge
            conflict.type == ConflictType.EDIT_CONFLICT &&
            conflict.localMessage.content is MessageContent.Text &&
            conflict.remoteMessage.content is MessageContent.Text -> ResolutionStrategy.CONTENT_MERGE

            // Simple edit conflicts use Last Writer Wins
            conflict.type == ConflictType.EDIT_CONFLICT -> ResolutionStrategy.LAST_WRITER_WINS

            // Delete conflicts require user input
            conflict.type == ConflictType.DELETE_CONFLICT -> ResolutionStrategy.USER_CHOICE_REQUIRED

            // Default to user choice for unknown scenarios
            else -> ResolutionStrategy.USER_CHOICE_REQUIRED
        }
    }

    private fun getResolver(strategy: ResolutionStrategy): ConflictResolver {
        return when (strategy) {
            ResolutionStrategy.LAST_WRITER_WINS -> lastWriterWinsResolver
            ResolutionStrategy.STATUS_MERGE -> statusMergeResolver
            ResolutionStrategy.CONTENT_MERGE -> contentMergeResolver
            ResolutionStrategy.USER_CHOICE_REQUIRED -> userChoiceResolver
            else -> userChoiceResolver
        }
    }
}
```

## 5. Conflict Prevention Strategies

### Optimistic Locking with Version Control
```kotlin
data class VersionedMessage(
    val message: ChatMessage,
    val version: Long,
    val checksum: String
) {
    companion object {
        fun from(message: ChatMessage): VersionedMessage {
            return VersionedMessage(
                message = message,
                version = message.version ?: 1L,
                checksum = calculateChecksum(message)
            )
        }

        private fun calculateChecksum(message: ChatMessage): String {
            return "${message.id}-${message.content.hashCode()}-${message.timestamp}".sha256()
        }
    }
}

class OptimisticLockingService {
    suspend fun updateMessage(
        messageId: String,
        expectedVersion: Long,
        update: (ChatMessage) -> ChatMessage
    ): Result<ChatMessage> {
        return try {
            val currentMessage = localDataSource.getMessage(messageId)
                ?: return Result.failure(MessageNotFoundException(messageId))

            if (currentMessage.version != expectedVersion) {
                return Result.failure(VersionConflictException(
                    messageId = messageId,
                    expectedVersion = expectedVersion,
                    actualVersion = currentMessage.version ?: 0L
                ))
            }

            val updatedMessage = update(currentMessage).copy(
                version = (currentMessage.version ?: 0L) + 1,
                editedAt = Clock.System.now()
            )

            localDataSource.updateMessage(updatedMessage)

            Result.success(updatedMessage)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
```

### Conflict-Free Replicated Data Types (CRDT) Approach
```kotlin
sealed class MessageOperation {
    data class Insert(
        val position: Int,
        val content: String,
        val timestamp: Instant,
        val actorId: String
    ) : MessageOperation()

    data class Delete(
        val position: Int,
        val length: Int,
        val timestamp: Instant,
        val actorId: String
    ) : MessageOperation()

    data class Format(
        val start: Int,
        val end: Int,
        val formatting: MessageTextFormatting,
        val timestamp: Instant,
        val actorId: String
    ) : MessageOperation()
}

class CRDTMessageContent(
    private val operations: MutableList<MessageOperation> = mutableListOf()
) {
    fun applyOperation(operation: MessageOperation) {
        // Insert operation in timestamp order
        val insertIndex = operations.binarySearch { existing ->
            existing.timestamp.compareTo(operation.timestamp)
        }.let { if (it < 0) -(it + 1) else it }

        operations.add(insertIndex, operation)
    }

    fun computeCurrentContent(): MessageContent.Text {
        var content = ""
        val formatting = mutableListOf<MessageTextFormatting>()

        operations.forEach { operation ->
            when (operation) {
                is MessageOperation.Insert -> {
                    content = content.substring(0, operation.position) +
                             operation.content +
                             content.substring(operation.position)
                }
                is MessageOperation.Delete -> {
                    val endPos = (operation.position + operation.length).coerceAtMost(content.length)
                    content = content.substring(0, operation.position) +
                             content.substring(endPos)
                }
                is MessageOperation.Format -> {
                    formatting.add(operation.formatting)
                }
            }
        }

        return MessageContent.Text(content, formatting)
    }
}
```

## 6. Conflict Resolution Metrics & Learning

### Resolution Analytics
```kotlin
class ConflictAnalytics {
    private val resolutionHistory = mutableListOf<ConflictResolutionRecord>()

    fun recordResolution(
        conflict: MessageConflict,
        resolution: ConflictResolution,
        userSatisfaction: UserSatisfactionRating? = null
    ) {
        resolutionHistory.add(
            ConflictResolutionRecord(
                conflictType = conflict.type,
                severity = conflict.severity,
                strategy = resolution.strategy,
                success = resolution.success,
                userSatisfaction = userSatisfaction,
                resolveTime = Clock.System.now() - conflict.detectedAt
            )
        )
    }

    fun getResolutionInsights(): ConflictInsights {
        return ConflictInsights(
            totalConflicts = resolutionHistory.size,
            autoResolutionRate = resolutionHistory.count { it.strategy != ResolutionStrategy.USER_CHOICE_REQUIRED }.toDouble() / resolutionHistory.size,
            averageResolutionTime = resolutionHistory.map { it.resolveTime.inWholeMilliseconds }.average(),
            successRate = resolutionHistory.count { it.success }.toDouble() / resolutionHistory.size,
            strategyEffectiveness = resolutionHistory.groupBy { it.strategy }
                .mapValues { (_, records) ->
                    records.count { it.success }.toDouble() / records.size
                }
        )
    }
}

data class ConflictResolutionRecord(
    val conflictType: ConflictType,
    val severity: ConflictSeverity,
    val strategy: ResolutionStrategy,
    val success: Boolean,
    val userSatisfaction: UserSatisfactionRating?,
    val resolveTime: Duration
)

enum class UserSatisfactionRating {
    VERY_SATISFIED,
    SATISFIED,
    NEUTRAL,
    DISSATISFIED,
    VERY_DISSATISFIED
}
```

## 7. Integration with Sync Engine

### Conflict-Aware Synchronization
```kotlin
class ConflictAwareSyncEngine(
    private val conflictDetector: ConflictDetector,
    private val resolutionEngine: ConflictResolutionEngine,
    private val analytics: ConflictAnalytics
) : SyncEngine {

    override suspend fun syncChat(chatId: String, strategy: SyncStrategy): Result<SyncResult> {
        return try {
            // 1. Fetch remote changes
            val remoteChanges = remoteDataSource.fetchMessages(chatId,
                since = localDataSource.getLastSyncTimestamp(chatId))

            // 2. Detect conflicts
            val conflicts = conflictDetector.detectMessageConflicts(chatId)

            // 3. Resolve conflicts
            val resolutions = conflicts.map { conflict ->
                resolutionEngine.resolveConflict(conflict)
            }

            // 4. Apply resolutions
            resolutions.forEach { resolution ->
                when (resolution.strategy) {
                    ResolutionStrategy.LOCAL_WINS -> {
                        // Push local version to remote
                        remoteDataSource.updateMessage(resolution.resolvedMessage!!)
                    }
                    ResolutionStrategy.REMOTE_WINS -> {
                        // Update local with remote version
                        localDataSource.saveMessage(resolution.resolvedMessage!!)
                    }
                    ResolutionStrategy.MERGE -> {
                        // Save merged version locally and push to remote
                        localDataSource.saveMessage(resolution.resolvedMessage!!)
                        remoteDataSource.updateMessage(resolution.resolvedMessage!!)
                    }
                }

                analytics.recordResolution(
                    conflicts.first { it.id == resolution.conflictId },
                    resolution
                )
            }

            // 5. Continue with normal sync for non-conflicted data
            val syncResult = performNormalSync(chatId, strategy, conflicts)

            Result.success(syncResult.copy(
                conflicts = conflicts,
                resolutions = resolutions
            ))

        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
```

This comprehensive conflict resolution system provides:

1. **Automatic Detection**: Temporal and content-based conflict identification
2. **Multiple Resolution Strategies**: LWW, merge, user choice based on conflict type and severity
3. **User Interface Integration**: Seamless conflict resolution UI for manual decisions
4. **Prevention Mechanisms**: Optimistic locking and CRDT approaches
5. **Analytics & Learning**: Resolution effectiveness tracking and strategy optimization
6. **Sync Integration**: Conflict-aware synchronization with the overall sync engine

The system balances automation with user control, ensuring data integrity while maintaining a smooth user experience.