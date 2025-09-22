package com.tchat.app.api.models

import java.time.Instant

/**
 * State sync models for testing
 */

data class UserStateResponse(
    val userId: String,
    val timestamp: Instant,
    val themePreferences: ThemePreferences,
    val navigationState: NavigationState,
    val sessionInfo: SessionInfo
)

data class ThemePreferences(
    val colorScheme: String,
    val customColors: Map<String, String>
)

data class NavigationState(
    val currentTab: String,
    val history: List<String>
)

data class SessionInfo(
    val sessionId: String,
    val platform: String,
    val startTime: Instant
)

data class StateSyncRequest(
    val userId: String,
    val platform: String,
    val state: UserState,
    val timestamp: Instant
)

data class UserState(
    val themePreferences: ThemePreferences,
    val navigationState: NavigationState,
    val customData: Map<String, Any> = emptyMap()
)

data class StateSyncResponse(
    val success: Boolean,
    val message: String,
    val timestamp: Instant,
    val conflictResolution: List<StateConflict>
)

data class StateConflict(
    val conflictId: String,
    val field: String,
    val localValue: Any,
    val remoteValue: Any,
    val suggestedResolution: String
)

data class ConflictResolution(
    val conflictId: String,
    val strategy: ConflictResolutionStrategy
)

enum class ConflictResolutionStrategy {
    USE_LOCAL,
    USE_REMOTE,
    MERGE,
    MANUAL
}

data class ConflictResolutionResponse(
    val conflictId: String,
    val resolved: Boolean,
    val appliedStrategy: ConflictResolutionStrategy,
    val timestamp: Instant
)