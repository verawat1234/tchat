package com.tchat.app.api

import com.tchat.app.api.models.ComponentCategory
import kotlinx.coroutines.delay

/**
 * Mock API client for Component Registry testing
 */
class ComponentRegistryAPIClient(private val baseURL: String) {

    suspend fun getComponents(): List<ComponentMetadata> {
        // Simulate network delay
        delay(100)

        return listOf(
            ComponentMetadata(
                id = "tchat-button",
                name = "TchatButton",
                platform = "android",
                version = "1.0.0",
                webEquivalent = "TchatButton",
                accessibility = AccessibilityInfo(
                    hasSemanticLabels = true,
                    supportsScreenReader = true,
                    keyboardNavigable = true
                ),
                category = ComponentCategory.INTERACTIVE
            ),
            ComponentMetadata(
                id = "tchat-input",
                name = "TchatInput",
                platform = "android",
                version = "1.0.0",
                webEquivalent = "TchatInput",
                accessibility = AccessibilityInfo(
                    hasSemanticLabels = true,
                    supportsScreenReader = true,
                    keyboardNavigable = true
                ),
                category = ComponentCategory.INTERACTIVE
            )
        )
    }

    suspend fun getComponents(category: ComponentCategory): List<ComponentMetadata> {
        return getComponents().filter { it.category == category }
    }

    suspend fun getComponent(id: String): ComponentMetadata? {
        return getComponents().find { it.id == id }
    }
}

/**
 * Component metadata for testing
 */
data class ComponentMetadata(
    val id: String,
    val name: String,
    val platform: String,
    val version: String,
    val webEquivalent: String,
    val accessibility: AccessibilityInfo,
    val category: ComponentCategory
)

/**
 * Accessibility information for testing
 */
data class AccessibilityInfo(
    val hasSemanticLabels: Boolean,
    val supportsScreenReader: Boolean,
    val keyboardNavigable: Boolean
)