package com.tchat.app.api

import com.tchat.app.api.models.*
import kotlinx.coroutines.delay
import java.time.Instant

/**
 * Mock API client for State Sync testing
 */
class StateSyncAPIClient(private val baseURL: String) {

    suspend fun getState(userId: String): UserStateResponse {
        // Simulate network delay
        delay(100)

        return UserStateResponse(
            userId = userId,
            timestamp = Instant.now(),
            themePreferences = ThemePreferences(
                colorScheme = "light",
                customColors = mapOf("primary" to "#007AFF")
            ),
            navigationState = NavigationState(
                currentTab = "chat",
                history = listOf("home", "chat")
            ),
            sessionInfo = SessionInfo(
                sessionId = "test-session-123",
                platform = "android",
                startTime = Instant.now()
            )
        )
    }

    suspend fun syncState(request: StateSyncRequest): StateSyncResponse {
        // Simulate network delay
        delay(100)

        return StateSyncResponse(
            success = true,
            message = "State synced successfully",
            timestamp = Instant.now(),
            conflictResolution = emptyList()
        )
    }

    suspend fun resolveConflict(conflictId: String, resolution: ConflictResolution): ConflictResolutionResponse {
        // Simulate network delay
        delay(100)

        return ConflictResolutionResponse(
            conflictId = conflictId,
            resolved = true,
            appliedStrategy = resolution.strategy,
            timestamp = Instant.now()
        )
    }
}