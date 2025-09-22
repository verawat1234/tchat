package com.tchat.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.Contextual
import java.util.Date
import java.util.UUID

/**
 * Core navigation state entity for cross-platform state synchronization
 */
@Serializable
data class NavigationState(
    val id: String = UUID.randomUUID().toString(),
    val userId: String,
    val sessionId: String,
    val platform: String = "android",
    val navigationStack: MutableList<NavigationStackEntry>,
    var currentRoute: String? = null,
    var previousRoute: String? = null,
    val timestamp: Long = System.currentTimeMillis(),
    var syncVersion: Int = 1,
    val metadata: NavigationStateMetadata = NavigationStateMetadata()
) {

    /**
     * Get current navigation depth
     */
    val depth: Int
        get() = navigationStack.size

    /**
     * Check if can go back
     */
    val canGoBack: Boolean
        get() = navigationStack.size > 1

    /**
     * Get root route
     */
    val rootRoute: String?
        get() = navigationStack.firstOrNull()?.routeId

    /**
     * Get current route parameters
     */
    val currentParameters: Map<String, Any>
        get() = navigationStack.lastOrNull()?.parameters ?: emptyMap()

    /**
     * Check if state is synchronized
     */
    val isSynchronized: Boolean
        get() {
            val lastSync = metadata.lastSyncTimestamp ?: return false
            return (System.currentTimeMillis() - lastSync) < 30_000 // 30 seconds
        }

    /**
     * Push new route to navigation stack
     */
    fun push(
        routeId: String,
        parameters: Map<String, Any> = emptyMap(),
        timestamp: Long = System.currentTimeMillis()
    ) {
        val entry = NavigationStackEntry(
            routeId = routeId,
            parameters = parameters,
            timestamp = timestamp,
            transition = NavigationTransition.PUSH
        )

        previousRoute = currentRoute
        currentRoute = routeId
        navigationStack.add(entry)
        incrementVersion()
    }

    /**
     * Pop current route from navigation stack
     */
    fun pop() {
        if (navigationStack.size <= 1) return

        val poppedEntry = navigationStack.removeAt(navigationStack.size - 1)
        previousRoute = currentRoute
        currentRoute = navigationStack.lastOrNull()?.routeId

        // Add reverse transition entry for sync
        val reverseEntry = poppedEntry.copy(
            transition = NavigationTransition.POP,
            timestamp = System.currentTimeMillis()
        )

        incrementVersion()
    }

    /**
     * Replace current route
     */
    fun replace(
        routeId: String,
        parameters: Map<String, Any> = emptyMap(),
        timestamp: Long = System.currentTimeMillis()
    ) {
        if (navigationStack.isEmpty()) {
            push(routeId, parameters, timestamp)
            return
        }

        val entry = NavigationStackEntry(
            routeId = routeId,
            parameters = parameters,
            timestamp = timestamp,
            transition = NavigationTransition.REPLACE
        )

        previousRoute = currentRoute
        currentRoute = routeId
        navigationStack[navigationStack.size - 1] = entry
        incrementVersion()
    }

    /**
     * Reset navigation to root
     */
    fun popToRoot() {
        if (navigationStack.size <= 1) return

        val rootEntry = navigationStack.first()
        previousRoute = currentRoute
        currentRoute = rootEntry.routeId
        navigationStack.clear()
        navigationStack.add(rootEntry)
        incrementVersion()
    }

    /**
     * Reset entire navigation state
     */
    fun reset(
        toRoute: String,
        parameters: Map<String, Any> = emptyMap()
    ) {
        val entry = NavigationStackEntry(
            routeId = toRoute,
            parameters = parameters,
            timestamp = System.currentTimeMillis(),
            transition = NavigationTransition.RESET
        )

        previousRoute = currentRoute
        currentRoute = toRoute
        navigationStack.clear()
        navigationStack.add(entry)
        incrementVersion()
    }

    /**
     * Mark state as synchronized
     */
    fun markSynchronized(): NavigationState {
        return copy(
            metadata = metadata.copy(
                lastSyncTimestamp = System.currentTimeMillis()
            )
        )
    }

    /**
     * Increment sync version
     */
    private fun incrementVersion() {
        syncVersion++
    }

    /**
     * Create sync request payload
     */
    fun createSyncRequest(): NavigationStateSyncRequest {
        return NavigationStateSyncRequest(
            userId = userId,
            sessionId = sessionId,
            platform = platform,
            navigationStack = navigationStack.toList(),
            timestamp = System.currentTimeMillis(),
            syncVersion = syncVersion
        )
    }

    /**
     * Apply sync response
     */
    fun applySyncResponse(response: NavigationStateSyncResponse): NavigationState {
        return if (response.success) {
            copy(syncVersion = response.syncVersion).markSynchronized()
        } else {
            this
        }
    }

    companion object {
        /**
         * Create default navigation state for main app flow
         */
        fun defaultState(
            userId: String,
            sessionId: String
        ): NavigationState {
            val rootEntry = NavigationStackEntry(
                routeId = "chat",
                parameters = emptyMap(),
                timestamp = System.currentTimeMillis(),
                transition = NavigationTransition.RESET
            )

            return NavigationState(
                userId = userId,
                sessionId = sessionId,
                platform = "android",
                navigationStack = mutableListOf(rootEntry),
                currentRoute = "chat",
                previousRoute = null,
                timestamp = System.currentTimeMillis(),
                syncVersion = 1,
                metadata = NavigationStateMetadata()
            )
        }

        /**
         * Create navigation state from deep link
         */
        fun fromDeepLink(
            userId: String,
            sessionId: String,
            deepLinkUrl: String
        ): NavigationState? {
            // Parse deep link and create appropriate navigation state
            // This is a simplified implementation
            if (!deepLinkUrl.startsWith("tchat://")) {
                return null
            }

            val path = deepLinkUrl.removePrefix("tchat://").trim('/')
            if (path.isEmpty()) {
                return defaultState(userId, sessionId)
            }

            val routeId = path.replace("/", "/")
            val entry = NavigationStackEntry(
                routeId = routeId,
                parameters = emptyMap(), // Could extract query parameters
                timestamp = System.currentTimeMillis(),
                transition = NavigationTransition.DEEP_LINK
            )

            return NavigationState(
                userId = userId,
                sessionId = sessionId,
                platform = "android",
                navigationStack = mutableListOf(entry),
                currentRoute = routeId,
                previousRoute = null,
                timestamp = System.currentTimeMillis(),
                syncVersion = 1,
                metadata = NavigationStateMetadata()
            )
        }
    }
}

/**
 * Navigation stack entry with transition information
 */
@Serializable
data class NavigationStackEntry(
    val id: String = UUID.randomUUID().toString(),
    val routeId: String,
    val parameters: Map<String, @Contextual Any> = emptyMap(),
    val timestamp: Long = System.currentTimeMillis(),
    var transition: NavigationTransition = NavigationTransition.PUSH
)

/**
 * Navigation transition types
 */
@Serializable
enum class NavigationTransition {
    PUSH,
    POP,
    REPLACE,
    RESET,
    DEEP_LINK
}

/**
 * Navigation state metadata
 */
@Serializable
data class NavigationStateMetadata(
    val lastSyncTimestamp: Long? = null,
    val syncAttempts: Int = 0,
    val conflictResolutionStrategy: ConflictResolutionStrategy = ConflictResolutionStrategy.CLIENT_WINS,
    val customData: Map<String, String> = emptyMap()
)


/**
 * Navigation state sync request
 */
@Serializable
data class NavigationStateSyncRequest(
    val userId: String,
    val sessionId: String,
    val platform: String,
    val navigationStack: List<NavigationStackEntry>,
    val timestamp: Long,
    val syncVersion: Int
)

/**
 * Navigation state sync response
 */
@Serializable
data class NavigationStateSyncResponse(
    val success: Boolean,
    val syncVersion: Int,
    val conflictsResolved: List<String> = emptyList(),
    val timestamp: Long = System.currentTimeMillis()
)