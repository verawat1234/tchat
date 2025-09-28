package com.tchat.mobile.models

import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import kotlinx.serialization.Contextual

/**
 * Sync-related models for SQLDelight ↔ API data management
 */

// Type aliases for missing classes - map to existing classes
typealias ChatMessage = Message
typealias ChatSettings = NotificationSettings
typealias MessageContent = String

// Missing classes - add minimal implementations
@Serializable
data class DataIntegrityReport(
    val id: String = generateId(),
    val timestamp: Instant = Clock.System.now(),
    val chatId: String = "",
    val isValid: Boolean = true,
    val missingMessages: List<String> = emptyList(),
    val duplicateMessages: List<String> = emptyList(),
    val corruptedMessages: List<String> = emptyList(),
    val inconsistentTimestamps: List<String> = emptyList(),
    val brokenReferences: List<String> = emptyList(),
    val recommendedActions: List<String> = emptyList()
)

// Sync Operation Models
@Serializable
data class SyncOperation(
    val id: String,
    val type: OperationType,
    val chatId: String,
    val data: String, // JSON serialized operation data
    val timestamp: Instant,
    val status: SyncOperationStatus = SyncOperationStatus.PENDING,
    val priority: SyncPriority = SyncPriority.NORMAL,
    val retryCount: Int = 0,
    val maxRetries: Int = 3
)

@Serializable
enum class OperationType {
    SEND_MESSAGE,
    EDIT_MESSAGE,
    DELETE_MESSAGE,
    UPDATE_CHAT,
    JOIN_CHAT,
    LEAVE_CHAT,
    BACKGROUND_SYNC,
    FULL_SYNC,
    INCREMENTAL_SYNC
}

@Serializable
enum class SyncOperationStatus {
    PENDING,
    IN_PROGRESS,
    COMPLETED,
    FAILED
}

@Serializable
enum class SyncPriority {
    LOW,
    NORMAL,
    HIGH,
    CRITICAL
}

// Conflict Models
@Serializable
data class MessageConflict(
    val id: String = generateId(),
    val messageId: String,
    val chatId: String,
    val type: ConflictType,
    val severity: ConflictSeverity,
    @Contextual val localMessage: ChatMessage,
    @Contextual val remoteMessage: ChatMessage,
    val detectedAt: Instant = Clock.System.now(),
    val autoResolvable: Boolean = false
)

@Serializable
enum class ConflictType {
    EDIT_CONFLICT,
    STATUS_CONFLICT,
    DELETE_CONFLICT,
    PARTICIPANT_CONFLICT,
    METADATA_CONFLICT
}

@Serializable
enum class ConflictSeverity {
    CRITICAL,
    HIGH,
    MEDIUM,
    LOW
}

// Conflict Resolution Models
@Serializable
data class ConflictResolution(
    val conflictId: String,
    val strategy: ResolutionStrategy,
    @Contextual val resolvedMessage: ChatMessage? = null,
    val explanation: String? = null,
    val success: Boolean = true
) {
    companion object {
        fun unresolvable(conflictId: String): ConflictResolution {
            return ConflictResolution(
                conflictId = conflictId,
                strategy = ResolutionStrategy.USER_CHOICE_REQUIRED,
                success = false,
                explanation = "Manual resolution required"
            )
        }
    }
}

@Serializable
enum class ResolutionStrategy {
    LOCAL_WINS,
    REMOTE_WINS,
    MERGE,
    USER_CHOICE_REQUIRED,
    STATUS_MERGE,
    CONTENT_MERGE,
    LAST_WRITER_WINS,
    DISCARD
}

// Sync Result Models
@Serializable
data class SyncResult(
    val chatId: String,
    val success: Boolean,
    val messagesUpdated: Int,
    val conflicts: List<MessageConflict>,
    val lastSyncTimestamp: Instant,
    val syncDuration: kotlin.time.Duration? = null,
    val errors: List<String> = emptyList(),
    val resolutions: List<ConflictResolution> = emptyList()
)

// Chat Updates Model
@Serializable
data class ChatUpdates(
    val title: String? = null,
    val description: String? = null,
    val avatar: String? = null,
    @Contextual val settings: ChatSettings? = null
)

// Presence and Real-time Models
@Serializable
data class PresenceUpdate(
    val userId: String,
    val isOnline: Boolean,
    val lastSeen: Instant?,
    val activity: UserActivity?
)

@Serializable
enum class UserActivity {
    TYPING,
    RECORDING_AUDIO,
    RECORDING_VIDEO,
    IDLE
}

// Network and Connection Models
@Serializable
enum class NetworkState {
    CONNECTED,
    DISCONNECTED,
    LIMITED_CONNECTIVITY,
    CONNECTING
}

@Serializable
enum class ConnectionState {
    DISCONNECTED,
    CONNECTING,
    CONNECTED,
    RECONNECTING,
    ERROR
}

// Sync Strategy Enums
@Serializable
enum class SyncStrategy {
    INTELLIGENT,      // Smart strategy selection based on conditions
    LOCAL_FIRST,      // Prioritize local data, sync to remote
    REMOTE_FIRST,     // Prioritize remote data, update local
    BIDIRECTIONAL,    // Two-way sync with conflict resolution
    INCREMENTAL,      // Sync only changes since last sync
    FULL_REFRESH,     // Complete data refresh from remote
    OFFLINE_ONLY,     // Work with local data only
    REAL_TIME_ONLY    // Real-time sync only, no background sync
}

// Global Sync State
@Serializable
data class GlobalSyncState(
    val isConnected: Boolean,
    val activeSyncs: Int,
    val pendingOperations: Int,
    val lastSuccessfulSync: Instant?,
    val networkState: NetworkState
)

// Sync Info for UI (renamed to avoid conflict with existing SyncStatus enum)
@Serializable
data class SyncInfo(
    val state: SyncState,
    val lastSyncTime: Instant?,
    val pendingOperations: Int,
    val conflicts: List<MessageConflict>
)

// Sync State enum from LocalDataSource
@Serializable
enum class SyncState {
    IDLE,
    SYNCING,
    SYNCED,
    ERROR,
    CONFLICT_PENDING
}

// API Exception
class ApiException(
    val code: Int? = null,
    override val message: String
) : Exception(message) {
    constructor(message: String) : this(null, message)
}

// Version Conflict Exception
class VersionConflictException(
    val messageId: String,
    val expectedVersion: Long,
    val actualVersion: Long
) : Exception("Version conflict for message $messageId: expected $expectedVersion, actual $actualVersion")

// Utility function for generating IDs
private fun generateId(): String {
    return "sync_${Clock.System.now().toEpochMilliseconds()}_${kotlin.random.Random.nextInt(1000, 9999)}"
}

// Extension functions for conflict handling
fun determineAutoResolvability(type: ConflictType, severity: ConflictSeverity): Boolean {
    return when (type) {
        ConflictType.STATUS_CONFLICT -> severity == ConflictSeverity.LOW
        ConflictType.EDIT_CONFLICT -> severity in listOf(ConflictSeverity.LOW, ConflictSeverity.MEDIUM)
        ConflictType.DELETE_CONFLICT -> false // Always requires user input
        ConflictType.PARTICIPANT_CONFLICT -> severity == ConflictSeverity.LOW
        ConflictType.METADATA_CONFLICT -> severity == ConflictSeverity.LOW
    }
}

// Edit Message Data for sync operations
@Serializable
data class EditMessageData(
    val messageId: String,
    @Contextual val newContent: MessageContent
)

@Serializable
data class DeleteMessageData(
    val messageId: String
)

// Message content serialization support (using the JsonExtensions utility)
fun MessageContent.toJsonString(): String {
    return try {
        this.toString() // Simple fallback for now
    } catch (e: Exception) {
        toString() // Fallback for other types
    }
}