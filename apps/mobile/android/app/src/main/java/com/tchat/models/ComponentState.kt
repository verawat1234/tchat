package com.tchat.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.Contextual
import java.util.Date
import java.util.UUID

/**
 * Component state entity for UI component state management and synchronization
 */
@Serializable
data class ComponentState(
    val id: String = UUID.randomUUID().toString(),
    val componentId: String,
    val instanceId: String,
    val userId: String,
    val sessionId: String,
    val state: MutableMap<String, @Contextual Any>,
    val timestamp: Long = System.currentTimeMillis(),
    var version: Int = 1,
    val platform: String = "android",
    val metadata: ComponentStateMetadata = ComponentStateMetadata()
) {

    /**
     * Check if state is synchronized
     */
    val isSynchronized: Boolean
        get() {
            val lastSync = metadata.lastSyncTimestamp ?: return false
            return (System.currentTimeMillis() - lastSync) < 10_000 // 10 seconds
        }

    /**
     * Get state keys
     */
    val stateKeys: Set<String>
        get() = state.keys

    /**
     * Check if state is empty
     */
    val isEmpty: Boolean
        get() = state.isEmpty()

    /**
     * Get state size (number of properties)
     */
    val stateSize: Int
        get() = state.size

    /**
     * Calculate state hash for conflict detection
     */
    val stateHash: String
        get() {
            val sortedKeys = state.keys.sorted()
            val stateString = sortedKeys.joinToString("|") { "$it:${state[it] ?: ""}" }
            return stateString.hashCode().toString()
        }

    /**
     * Update specific state property
     */
    fun updateProperty(key: String, value: Any) {
        state[key] = value
        incrementVersion()
    }

    /**
     * Update multiple state properties
     */
    fun updateProperties(updates: Map<String, Any>) {
        state.putAll(updates)
        incrementVersion()
    }

    /**
     * Remove state property
     */
    fun removeProperty(key: String) {
        state.remove(key)
        incrementVersion()
    }

    /**
     * Clear all state
     */
    fun clearState() {
        state.clear()
        incrementVersion()
    }

    /**
     * Merge state with another state
     */
    fun mergeState(otherState: Map<String, Any>, strategy: MergeStrategy = MergeStrategy.OVERWRITE) {
        when (strategy) {
            MergeStrategy.OVERWRITE -> {
                state.putAll(otherState)
            }
            MergeStrategy.KEEP_EXISTING -> {
                for ((key, value) in otherState) {
                    if (!state.containsKey(key)) {
                        state[key] = value
                    }
                }
            }
            MergeStrategy.MERGE -> {
                state.putAll(otherState)
            }
        }
        incrementVersion()
    }

    /**
     * Get state property with type safety
     */
    inline fun <reified T> getProperty(key: String): T? {
        return state[key] as? T
    }

    /**
     * Get state property with default value
     */
    inline fun <reified T> getProperty(key: String, defaultValue: T): T {
        return state[key] as? T ?: defaultValue
    }

    /**
     * Create state sync request
     */
    fun createSyncRequest(): ComponentStateSyncRequest {
        return ComponentStateSyncRequest(
            userId = userId,
            sessionId = sessionId,
            platform = platform,
            componentStates = listOf(this),
            timestamp = System.currentTimeMillis(),
            syncVersion = version
        )
    }

    /**
     * Apply sync response
     */
    fun applySyncResponse(response: ComponentStateSyncResponse): ComponentState {
        return if (response.success) {
            copy(version = response.syncVersion).markSynchronized()
        } else {
            this
        }
    }

    /**
     * Mark state as synchronized
     */
    fun markSynchronized(): ComponentState {
        return copy(
            metadata = metadata.copy(
                lastSyncTimestamp = System.currentTimeMillis()
            )
        )
    }

    /**
     * Increment version
     */
    private fun incrementVersion() {
        version++
    }

    /**
     * Detect conflicts with another state
     */
    fun hasConflictWith(otherState: ComponentState): Boolean {
        return componentId == otherState.componentId &&
                instanceId == otherState.instanceId &&
                version != otherState.version &&
                stateHash != otherState.stateHash
    }

    /**
     * Resolve conflict with another state
     */
    fun resolveConflictWith(
        otherState: ComponentState,
        strategy: ConflictResolutionStrategy = ConflictResolutionStrategy.NEWEST_WINS
    ): ComponentState {
        return when (strategy) {
            ConflictResolutionStrategy.NEWEST_WINS -> {
                if (otherState.timestamp > timestamp) {
                    copy(
                        state = otherState.state.toMutableMap(),
                        version = otherState.version
                    )
                } else {
                    this
                }
            }
            ConflictResolutionStrategy.HIGHEST_VERSION_WINS -> {
                if (otherState.version > version) {
                    copy(
                        state = otherState.state.toMutableMap(),
                        version = otherState.version
                    )
                } else {
                    this
                }
            }
            ConflictResolutionStrategy.MERGE -> {
                val mergedState = state.toMutableMap()
                mergedState.putAll(otherState.state)
                copy(
                    state = mergedState,
                    version = maxOf(version, otherState.version) + 1
                )
            }
            ConflictResolutionStrategy.MANUAL -> {
                // Manual resolution required - emit conflict event
                this
            }
            ConflictResolutionStrategy.CLIENT_WINS -> {
                this
            }
            ConflictResolutionStrategy.SERVER_WINS -> {
                copy(
                    state = otherState.state.toMutableMap(),
                    version = otherState.version
                )
            }
            ConflictResolutionStrategy.PROMPT -> {
                // Prompt user for resolution
                this
            }
        }
    }

    companion object {
        /**
         * Create default component state for chat message
         */
        fun chatMessageState(
            instanceId: String,
            userId: String,
            sessionId: String,
            isRead: Boolean = false,
            isSelected: Boolean = false
        ): ComponentState {
            return ComponentState(
                componentId = "chat-message",
                instanceId = instanceId,
                userId = userId,
                sessionId = sessionId,
                state = mutableMapOf(
                    "isRead" to isRead,
                    "isSelected" to isSelected,
                    "timestamp" to System.currentTimeMillis()
                ),
                platform = "android"
            )
        }

        /**
         * Create default component state for user avatar
         */
        fun userAvatarState(
            instanceId: String,
            userId: String,
            sessionId: String,
            isOnline: Boolean = false,
            lastSeen: Long? = null
        ): ComponentState {
            val state = mutableMapOf<String, Any>(
                "isOnline" to isOnline
            )

            lastSeen?.let { state["lastSeen"] = it }

            return ComponentState(
                componentId = "user-avatar",
                instanceId = instanceId,
                userId = userId,
                sessionId = sessionId,
                state = state,
                platform = "android"
            )
        }

        /**
         * Create default component state for navigation tab
         */
        fun navigationTabState(
            instanceId: String,
            userId: String,
            sessionId: String,
            isActive: Boolean = false,
            hasNotification: Boolean = false
        ): ComponentState {
            return ComponentState(
                componentId = "navigation-tab",
                instanceId = instanceId,
                userId = userId,
                sessionId = sessionId,
                state = mutableMapOf(
                    "isActive" to isActive,
                    "hasNotification" to hasNotification
                ),
                platform = "android"
            )
        }
    }
}

/**
 * Merge strategy for state updates
 */
@Serializable
enum class MergeStrategy {
    OVERWRITE,
    KEEP_EXISTING,
    MERGE
}


/**
 * Component state metadata
 */
@Serializable
data class ComponentStateMetadata(
    val lastSyncTimestamp: Long? = null,
    val syncAttempts: Int = 0,
    val conflictCount: Int = 0,
    val customData: Map<String, String> = emptyMap()
)

/**
 * Component state sync request
 */
@Serializable
data class ComponentStateSyncRequest(
    val userId: String,
    val sessionId: String,
    val platform: String,
    val componentStates: List<ComponentState>,
    val timestamp: Long,
    val syncVersion: Int
)

/**
 * Component state sync response
 */
@Serializable
data class ComponentStateSyncResponse(
    val success: Boolean,
    val syncVersion: Int,
    val conflictsResolved: List<String> = emptyList(),
    val timestamp: Long = System.currentTimeMillis()
)