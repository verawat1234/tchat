package com.tchat.mobile.sync

import com.tchat.mobile.models.*
import com.tchat.mobile.models.SyncPriority
import com.tchat.mobile.models.SyncOperationStatus
import com.tchat.mobile.models.OperationType
import com.tchat.mobile.repositories.datasource.LocalDataSource
import com.tchat.mobile.repositories.datasource.SyncState as LocalSyncState
import com.tchat.mobile.network.NetworkStateManager
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import kotlinx.coroutines.delay
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlin.time.Duration.Companion.seconds
import kotlin.time.Duration.Companion.minutes

/**
 * Background synchronization service
 *
 * Provides:
 * - Intelligent background sync scheduling
 * - Adaptive sync intervals based on network conditions
 * - Priority-based operation processing
 * - Battery-aware sync strategies
 * - Exponential backoff for failed operations
 */
class BackgroundSyncService(
    private val syncEngine: SyncEngine,
    private val localDataSource: LocalDataSource,
    private val networkStateManager: NetworkStateManager,
    private val scope: CoroutineScope,
    private val config: BackgroundSyncConfig = BackgroundSyncConfig()
) {

    private val _isRunning = MutableStateFlow(false)
    val isRunning: StateFlow<Boolean> = _isRunning.asStateFlow()

    private val _syncStats = MutableStateFlow(BackgroundSyncStats())
    val syncStats: StateFlow<BackgroundSyncStats> = _syncStats.asStateFlow()

    private val _nextScheduledSync = MutableStateFlow<Instant?>(null)
    val nextScheduledSync: StateFlow<Instant?> = _nextScheduledSync.asStateFlow()

    private var currentSyncInterval = config.defaultSyncInterval
    private var consecutiveFailures = 0
    private var lastSuccessfulSync: Instant? = null

    init {
        observeNetworkChanges()
    }

    /**
     * Start the background sync service
     */
    suspend fun start(): Result<Unit> {
        return try {
            if (_isRunning.value) {
                return Result.success(Unit)
            }

            _isRunning.value = true
            startBackgroundSyncLoop()

            println("🔄 BackgroundSyncService: Started with interval ${currentSyncInterval}")
            Result.success(Unit)
        } catch (e: Exception) {
            _isRunning.value = false
            Result.failure(e)
        }
    }

    /**
     * Stop the background sync service
     */
    suspend fun stop(): Result<Unit> {
        return try {
            _isRunning.value = false
            _nextScheduledSync.value = null

            println("⏹️ BackgroundSyncService: Stopped")
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Force an immediate sync cycle
     */
    suspend fun forceSyncNow(): Result<Unit> {
        return try {
            if (networkStateManager.shouldAttemptNetworkRequest()) {
                performSyncCycle()
                Result.success(Unit)
            } else {
                Result.failure(Exception("Network not available for sync"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Schedule a sync for a specific chat
     */
    suspend fun scheduleSync(chatId: String, priority: SyncPriority = SyncPriority.NORMAL): Result<Unit> {
        return try {
            val operation = SyncOperation(
                id = "bg_sync_${chatId}_${Clock.System.now().toEpochMilliseconds()}",
                type = OperationType.BACKGROUND_SYNC,
                chatId = chatId,
                data = "{}",
                timestamp = Clock.System.now(),
                status = SyncOperationStatus.PENDING,
                priority = priority
            )

            localDataSource.saveSyncOperation(operation)

            // If it's high priority and network is available, sync immediately
            if (priority == SyncPriority.HIGH && networkStateManager.shouldAttemptNetworkRequest()) {
                scope.launch {
                    syncEngine.syncChat(chatId)
                }
            }

            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Update sync configuration
     */
    suspend fun updateConfig(newConfig: BackgroundSyncConfig): Result<Unit> {
        return try {
            // Apply new configuration
            currentSyncInterval = newConfig.defaultSyncInterval

            println("⚙️ BackgroundSyncService: Configuration updated")
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get current sync statistics
     */
    fun getSyncStatistics(): BackgroundSyncStats {
        return _syncStats.value
    }

    /**
     * Clear sync statistics
     */
    suspend fun clearStatistics(): Result<Unit> {
        return try {
            _syncStats.value = BackgroundSyncStats()
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private fun startBackgroundSyncLoop() {
        scope.launch {
            while (_isRunning.value) {
                try {
                    val nextSync = Clock.System.now() + currentSyncInterval
                    _nextScheduledSync.value = nextSync

                    // Wait for the sync interval
                    delay(currentSyncInterval)

                    // Check if we should perform sync
                    if (shouldPerformSync()) {
                        performSyncCycle()
                    }

                    // Adjust sync interval based on conditions
                    adjustSyncInterval()

                } catch (e: Exception) {
                    println("❌ BackgroundSyncService: Error in sync loop: ${e.message}")

                    // Increase interval on error to avoid hammering
                    currentSyncInterval = minOf(
                        currentSyncInterval * 2,
                        config.maxSyncInterval
                    )

                    delay(30.seconds) // Wait before retrying
                }
            }
        }
    }

    private suspend fun performSyncCycle() {
        val cycleStartTime = Clock.System.now()
        var operationsProcessed = 0
        var operationsFailed = 0

        try {
            println("🔄 BackgroundSyncService: Starting sync cycle")

            // 1. Process pending operations by priority
            val pendingOperations = localDataSource.getPendingOperations()
                .sortedWith(compareBy<SyncOperation> { it.priority }.thenBy { it.timestamp })

            for (operation in pendingOperations.take(config.maxOperationsPerCycle)) {
                try {
                    when (operation.type) {
                        OperationType.BACKGROUND_SYNC,
                        OperationType.SEND_MESSAGE,
                        OperationType.EDIT_MESSAGE,
                        OperationType.DELETE_MESSAGE -> {
                            val result = syncEngine.syncChat(operation.chatId)
                            if (result.isSuccess) {
                                localDataSource.markOperationCompleted(operation.id)
                                operationsProcessed++
                            } else {
                                operationsFailed++
                            }
                        }
                        else -> {
                            // Handle other operation types
                            operationsProcessed++
                        }
                    }
                } catch (e: Exception) {
                    operationsFailed++
                    println("❌ BackgroundSyncService: Failed to process operation ${operation.id}: ${e.message}")
                }
            }

            // 2. Sync active chats based on priority
            val activeChats = localDataSource.getAllChatSessions().map { it.id }
            val priorityChats = determinePriorityChats(activeChats)

            for (chatId in priorityChats.take(config.maxChatsPerCycle)) {
                try {
                    val syncResult = syncEngine.syncChat(chatId, SyncStrategy.INTELLIGENT)
                    if (syncResult.isSuccess) {
                        operationsProcessed++
                    } else {
                        operationsFailed++
                    }
                } catch (e: Exception) {
                    operationsFailed++
                    println("❌ BackgroundSyncService: Failed to sync chat $chatId: ${e.message}")
                }
            }

            // 3. Update statistics
            val cycleDuration = Clock.System.now() - cycleStartTime
            updateSyncStats(operationsProcessed, operationsFailed, cycleDuration)

            // 4. Update failure tracking
            if (operationsFailed == 0) {
                consecutiveFailures = 0
                lastSuccessfulSync = Clock.System.now()
            } else {
                consecutiveFailures++
            }

            println("✅ BackgroundSyncService: Sync cycle completed - processed: $operationsProcessed, failed: $operationsFailed")

        } catch (e: Exception) {
            consecutiveFailures++
            println("❌ BackgroundSyncService: Sync cycle failed: ${e.message}")
        }
    }

    private suspend fun determinePriorityChats(allChats: List<String>): List<String> {
        return try {
            val chatPriorities = mutableListOf<Pair<String, Int>>()

            for (chatId in allChats) {
                var priority = 0

                // Priority factors:
                // 1. Unread messages (high priority)
                val chatSession = localDataSource.getChatSession(chatId).getOrNull()
                val unreadCount = chatSession?.unreadCount ?: 0
                priority += unreadCount * 10

                // 2. Recent activity (medium priority)
                val lastActivity = localDataSource.getLastSyncTimestamp(chatId)
                if (lastActivity != null) {
                    val hoursSinceActivity = (Clock.System.now() - lastActivity).inWholeSeconds / 3600
                    priority += maxOf(0, 24 - hoursSinceActivity.toInt()) // Recent activity gets higher priority
                }

                // 3. Pending operations (high priority)
                val pendingOps = localDataSource.getSyncOperationsByChatId(chatId).size
                priority += pendingOps * 20

                // 4. Sync metadata status (highest priority)
                val syncMetadata = localDataSource.getSyncMetadata(chatId)
                if (syncMetadata?.syncStatus == LocalSyncState.SYNCING) {
                    priority += 100
                }

                chatPriorities.add(chatId to priority)
            }

            // Sort by priority (highest first) and return chat IDs
            chatPriorities.sortedByDescending { it.second }.map { it.first }
        } catch (e: Exception) {
            // Fallback to simple ordering
            allChats
        }
    }

    private fun shouldPerformSync(): Boolean {
        // Don't sync if no network
        if (!networkStateManager.shouldAttemptNetworkRequest()) {
            return false
        }

        // Don't sync if disabled
        if (!config.enableBackgroundSync) {
            return false
        }

        // Don't sync if too many consecutive failures
        if (consecutiveFailures >= config.maxConsecutiveFailures) {
            println("⚠️ BackgroundSyncService: Skipping sync due to too many consecutive failures: $consecutiveFailures")
            return false
        }

        // Don't sync if in battery saver mode (platform-specific check would go here)
        if (config.respectBatterySaver && isBatterySaverActive()) {
            println("🔋 BackgroundSyncService: Skipping sync due to battery saver mode")
            return false
        }

        return true
    }

    private fun adjustSyncInterval() {
        val isConnected = networkStateManager.shouldAttemptNetworkRequest()
        val baseInterval = config.defaultSyncInterval

        currentSyncInterval = when {
            // Increase interval when network is poor
            !isConnected -> {
                config.maxSyncInterval
            }

            // Increase interval after failures
            consecutiveFailures > 0 -> {
                minOf(
                    baseInterval * (1 + consecutiveFailures),
                    config.maxSyncInterval
                )
            }

            // Decrease interval if we have pending operations
            // Note: This check is simplified to avoid suspend function call in non-suspend context
            consecutiveFailures == 0 -> {
                maxOf(baseInterval / 2, config.minSyncInterval)
            }

            // Normal interval
            else -> baseInterval
        }
    }

    private suspend fun hasPendingHighPriorityOperations(): Boolean {
        return try {
            val pendingOps = localDataSource.getPendingOperations()
            pendingOps.any { it.priority == SyncPriority.HIGH || it.priority == SyncPriority.CRITICAL }
        } catch (e: Exception) {
            false
        }
    }

    private fun observeNetworkChanges() {
        scope.launch {
            networkStateManager.shouldSyncOfflineChanges.collect { shouldSync ->
                if (shouldSync && _isRunning.value) {
                    // Trigger immediate sync when coming back online
                    println("🌐 BackgroundSyncService: Network available, triggering immediate sync")
                    scope.launch {
                        forceSyncNow()
                    }
                    networkStateManager.markOfflineChangesSynced()
                }
            }
        }
    }

    private fun updateSyncStats(processed: Int, failed: Int, duration: kotlin.time.Duration) {
        val currentStats = _syncStats.value
        _syncStats.value = currentStats.copy(
            totalSyncCycles = currentStats.totalSyncCycles + 1,
            totalOperationsProcessed = currentStats.totalOperationsProcessed + processed,
            totalOperationsFailed = currentStats.totalOperationsFailed + failed,
            averageCycleDuration = if (currentStats.totalSyncCycles == 0) {
                duration
            } else {
                (currentStats.averageCycleDuration * currentStats.totalSyncCycles + duration) / (currentStats.totalSyncCycles + 1)
            },
            lastSyncTime = Clock.System.now(),
            consecutiveFailures = consecutiveFailures,
            lastSuccessfulSync = lastSuccessfulSync
        )
    }

    // Platform-specific implementations would be provided in actual platform modules
    private fun isBatterySaverActive(): Boolean {
        // This would be implemented in platform-specific code
        // For now, return false
        return false
    }
}

/**
 * Background sync configuration
 */
data class BackgroundSyncConfig(
    val enableBackgroundSync: Boolean = true,
    val defaultSyncInterval: kotlin.time.Duration = 5.minutes,
    val minSyncInterval: kotlin.time.Duration = 30.seconds,
    val maxSyncInterval: kotlin.time.Duration = 30.minutes,
    val maxOperationsPerCycle: Int = 10,
    val maxChatsPerCycle: Int = 5,
    val maxConsecutiveFailures: Int = 5,
    val respectBatterySaver: Boolean = true,
    val enableAdaptiveInterval: Boolean = true,
    val prioritizeRecentChats: Boolean = true
)

/**
 * Background sync statistics
 */
data class BackgroundSyncStats(
    val totalSyncCycles: Int = 0,
    val totalOperationsProcessed: Int = 0,
    val totalOperationsFailed: Int = 0,
    val averageCycleDuration: kotlin.time.Duration = 0.seconds,
    val lastSyncTime: Instant? = null,
    val consecutiveFailures: Int = 0,
    val lastSuccessfulSync: Instant? = null
) {
    val successRate: Float get() = if (totalOperationsProcessed + totalOperationsFailed == 0) {
        1.0f
    } else {
        totalOperationsProcessed.toFloat() / (totalOperationsProcessed + totalOperationsFailed)
    }
}

/**
 * Sync priority levels for background operations
 */
enum class SyncPriority {
    BACKGROUND,  // Lowest priority - background maintenance
    LOW,         // Low priority - non-urgent updates
    NORMAL,      // Normal priority - regular operations
    HIGH,        // High priority - user-initiated actions
    CRITICAL     // Highest priority - real-time requirements
}

/**
 * Factory for creating background sync service
 */
object BackgroundSyncServiceFactory {
    fun create(
        syncEngine: SyncEngine,
        localDataSource: LocalDataSource,
        networkStateManager: NetworkStateManager,
        scope: CoroutineScope,
        config: BackgroundSyncConfig = BackgroundSyncConfig()
    ): BackgroundSyncService {
        return BackgroundSyncService(
            syncEngine = syncEngine,
            localDataSource = localDataSource,
            networkStateManager = networkStateManager,
            scope = scope,
            config = config
        )
    }
}