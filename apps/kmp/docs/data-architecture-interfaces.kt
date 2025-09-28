// SQLDelight ↔ API Data Management Architecture Interfaces

/**
 * Local Data Source - SQLDelight Implementation
 * Responsible for: Local persistence, caching, offline access
 */
interface LocalDataSource {
    // Message Operations
    suspend fun getMessages(chatId: String): Flow<List<ChatMessage>>
    suspend fun saveMessage(message: ChatMessage): Result<Unit>
    suspend fun updateMessageStatus(messageId: String, status: MessageStatus): Result<Unit>
    suspend fun markMessageAsSynced(messageId: String, serverTimestamp: Instant): Result<Unit>
    suspend fun getPendingSyncMessages(): List<ChatMessage>

    // Chat Operations
    suspend fun getChatSessions(): Flow<List<Chat>>
    suspend fun getChatSession(chatId: String): Result<Chat>
    suspend fun saveChatSession(chat: Chat): Result<Unit>
    suspend fun updateLastActivity(chatId: String, timestamp: Instant): Result<Unit>

    // Sync Metadata
    suspend fun getLastSyncTimestamp(chatId: String): Instant?
    suspend fun updateSyncTimestamp(chatId: String, timestamp: Instant): Result<Unit>
    suspend fun getPendingOperations(): List<SyncOperation>
    suspend fun markOperationCompleted(operationId: String): Result<Unit>
}

/**
 * Remote Data Source - API Implementation
 * Responsible for: Server communication, real-time updates, authoritative data
 */
interface RemoteDataSource {
    // Message Operations
    suspend fun fetchMessages(chatId: String, since: Instant? = null, limit: Int = 50): Result<List<ChatMessage>>
    suspend fun sendMessage(message: ChatMessage): Result<ChatMessage>
    suspend fun updateMessage(messageId: String, content: MessageContent): Result<ChatMessage>
    suspend fun deleteMessage(messageId: String): Result<Unit>

    // Real-time Subscriptions
    fun subscribeToMessages(chatId: String): Flow<ChatMessage>
    fun subscribeToTypingIndicators(chatId: String): Flow<TypingIndicator>
    fun subscribeToPresence(chatId: String): Flow<PresenceUpdate>

    // Chat Operations
    suspend fun fetchChatSessions(): Result<List<Chat>>
    suspend fun createChatSession(chat: Chat): Result<Chat>
    suspend fun updateChatSession(chatId: String, updates: ChatUpdates): Result<Chat>
    suspend fun joinChatSession(chatId: String): Result<Unit>
    suspend fun leaveChatSession(chatId: String): Result<Unit>

    // Connection Management
    suspend fun connect(): Result<Unit>
    suspend fun disconnect(): Result<Unit>
    fun getConnectionState(): Flow<ConnectionState>
}

/**
 * Sync Engine - Coordination Logic
 * Responsible for: Data consistency, conflict resolution, sync strategies
 */
interface SyncEngine {
    // Sync Operations
    suspend fun syncChat(chatId: String, strategy: SyncStrategy = SyncStrategy.INCREMENTAL): Result<SyncResult>
    suspend fun syncAllChats(): Result<List<SyncResult>>
    suspend fun forcePushPendingOperations(): Result<Unit>

    // Conflict Resolution
    suspend fun resolveConflicts(conflicts: List<DataConflict>): Result<List<ConflictResolution>>

    // Sync State Management
    fun getSyncState(chatId: String): Flow<SyncState>
    fun getGlobalSyncState(): Flow<GlobalSyncState>

    // Background Sync
    suspend fun startBackgroundSync()
    suspend fun stopBackgroundSync()
}

/**
 * Chat Repository - Domain Coordination
 * Responsible for: Business logic, data coordination, public API
 */
interface ChatRepository {
    // Message Operations (UI-facing)
    fun getMessages(chatId: String): Flow<List<ChatMessage>>
    suspend fun sendMessage(chatId: String, content: MessageContent, replyToId: String? = null): Result<ChatMessage>
    suspend fun editMessage(messageId: String, newContent: MessageContent): Result<ChatMessage>
    suspend fun deleteMessage(messageId: String): Result<Unit>
    suspend fun searchMessages(chatId: String, query: String): Result<List<ChatMessage>>

    // Chat Operations (UI-facing)
    fun getChatSessions(): Flow<List<Chat>>
    fun getChatSession(chatId: String): Flow<Chat?>
    suspend fun createChatSession(chat: Chat): Result<Chat>
    suspend fun updateChatSession(chatId: String, updates: ChatUpdates): Result<Chat>
    suspend fun archiveChat(chatId: String): Result<Unit>
    suspend fun muteChat(chatId: String): Result<Unit>
    suspend fun pinChat(chatId: String): Result<Unit>

    // Real-time Operations
    fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>>
    suspend fun sendTypingIndicator(chatId: String, type: TypingType): Result<Unit>

    // Sync Management
    suspend fun refreshChat(chatId: String): Result<Unit>
    suspend fun refreshAllChats(): Result<Unit>
    fun getSyncStatus(chatId: String): Flow<SyncStatus>
}

// Supporting Data Classes

data class SyncOperation(
    val id: String,
    val type: OperationType,
    val chatId: String,
    val data: String, // JSON serialized operation data
    val timestamp: Instant,
    val retryCount: Int = 0,
    val maxRetries: Int = 3
)

enum class OperationType {
    SEND_MESSAGE,
    EDIT_MESSAGE,
    DELETE_MESSAGE,
    UPDATE_CHAT,
    JOIN_CHAT,
    LEAVE_CHAT
}

data class SyncResult(
    val chatId: String,
    val success: Boolean,
    val messagesUpdated: Int,
    val conflicts: List<DataConflict>,
    val lastSyncTimestamp: Instant
)

data class DataConflict(
    val id: String,
    val type: ConflictType,
    val localData: Any,
    val remoteData: Any,
    val conflictTimestamp: Instant
)

enum class ConflictType {
    MESSAGE_EDIT_CONFLICT,
    CHAT_UPDATE_CONFLICT,
    PARTICIPANT_CHANGE_CONFLICT
}

data class ConflictResolution(
    val conflictId: String,
    val resolution: ResolutionStrategy,
    val resolvedData: Any
)

enum class ResolutionStrategy {
    LOCAL_WINS,
    REMOTE_WINS,
    MERGE,
    USER_CHOICE_REQUIRED
}

enum class SyncStrategy {
    FULL_REFRESH,
    INCREMENTAL,
    REAL_TIME_ONLY
}

enum class SyncState {
    IDLE,
    SYNCING,
    SYNCED,
    ERROR,
    CONFLICT_PENDING
}

data class GlobalSyncState(
    val isConnected: Boolean,
    val activeSyncs: Int,
    val pendingOperations: Int,
    val lastSuccessfulSync: Instant?,
    val networkState: NetworkState
)

enum class NetworkState {
    CONNECTED,
    DISCONNECTED,
    LIMITED_CONNECTIVITY,
    CONNECTING
}

enum class ConnectionState {
    DISCONNECTED,
    CONNECTING,
    CONNECTED,
    RECONNECTING,
    ERROR
}

data class ChatUpdates(
    val title: String? = null,
    val description: String? = null,
    val avatar: String? = null,
    val settings: ChatSettings? = null
)

data class PresenceUpdate(
    val userId: String,
    val isOnline: Boolean,
    val lastSeen: Instant?,
    val activity: UserActivity?
)

enum class UserActivity {
    TYPING,
    RECORDING_AUDIO,
    RECORDING_VIDEO,
    IDLE
}

data class SyncStatus(
    val state: SyncState,
    val lastSyncTime: Instant?,
    val pendingOperations: Int,
    val conflicts: List<DataConflict>
)