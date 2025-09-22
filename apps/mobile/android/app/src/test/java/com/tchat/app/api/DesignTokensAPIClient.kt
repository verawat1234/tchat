package com.tchat.app.api

import com.tchat.app.api.models.DesignTokensSyncRequest
import kotlinx.coroutines.delay

/**
 * Mock API client for Design Tokens testing
 */
class DesignTokensAPIClient(private val baseURL: String) {

    suspend fun getDesignTokens(platform: String): DesignTokensResponse {
        // Simulate network delay
        delay(100)

        return DesignTokensResponse(
            platform = platform,
            version = "1.0.0",
            tokens = mapOf(
                "colors" to mapOf(
                    "primary" to "#007AFF",
                    "secondary" to "#5856D6",
                    "background" to "#FFFFFF",
                    "surface" to "#F2F2F7"
                ),
                "typography" to mapOf(
                    "heading1" to mapOf(
                        "fontSize" to "28",
                        "fontWeight" to "bold"
                    ),
                    "body" to mapOf(
                        "fontSize" to "16",
                        "fontWeight" to "normal"
                    )
                ),
                "spacing" to mapOf(
                    "xs" to "4",
                    "sm" to "8",
                    "md" to "16",
                    "lg" to "24",
                    "xl" to "32"
                )
            )
        )
    }

    suspend fun syncDesignTokens(request: DesignTokensSyncRequest): SyncResponse {
        // Simulate network delay
        delay(100)

        return SyncResponse(
            success = true,
            message = "Design tokens synced successfully",
            version = "1.0.1"
        )
    }
}

/**
 * Design tokens response for testing
 */
data class DesignTokensResponse(
    val platform: String,
    val version: String,
    val tokens: Map<String, Any>
)

/**
 * Sync response for testing
 */
data class SyncResponse(
    val success: Boolean,
    val message: String,
    val version: String
)