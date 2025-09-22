package com.tchat.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.Contextual
import java.util.Date

/**
 * Core UI component entity with state management and cross-platform sync
 */
@Serializable
data class UIComponent(
    val id: String,
    val name: String,
    val category: ComponentCategory,
    val version: String = "1.0.0",
    val isStateful: Boolean = false,
    val propsSchema: ComponentSchema,
    val stateSchema: ComponentSchema? = null,
    val eventsSchema: ComponentSchema = ComponentSchema(),
    val dependencies: List<String> = emptyList(),
    val platformSupport: List<String> = listOf("ios", "android", "web"),
    val metadata: ComponentMetadata = ComponentMetadata()
) {

    /**
     * Check if component supports current platform
     */
    val isSupported: Boolean
        get() = platformSupport.contains("android")

    /**
     * Get component size estimation
     */
    val estimatedSize: ComponentSize
        get() {
            val propCount = propsSchema.properties.size
            val stateCount = stateSchema?.properties?.size ?: 0
            val eventCount = eventsSchema.properties.size

            val complexity = propCount + stateCount + eventCount

            return when (complexity) {
                in 0..5 -> ComponentSize.SMALL
                in 6..15 -> ComponentSize.MEDIUM
                in 16..30 -> ComponentSize.LARGE
                else -> ComponentSize.EXTRA_LARGE
            }
        }

    /**
     * Check if component has external dependencies
     */
    val hasExternalDependencies: Boolean
        get() = dependencies.isNotEmpty()

    /**
     * Generate component instance with state
     */
    fun createInstance(
        instanceId: String,
        initialState: Map<String, Any> = emptyMap(),
        props: Map<String, Any> = emptyMap()
    ): UIComponentInstance {
        return UIComponentInstance(
            instanceId = instanceId,
            componentId = id,
            componentName = name,
            state = if (isStateful) initialState.toMutableMap() else mutableMapOf(),
            props = props,
            version = version,
            createdAt = System.currentTimeMillis(),
            lastUpdated = System.currentTimeMillis()
        )
    }

    companion object {
        /**
         * Default UI components for the application
         */
        val defaultComponents = listOf(
            UIComponent(
                id = "chat-message",
                name = "ChatMessage",
                category = ComponentCategory.DISPLAY,
                isStateful = true,
                propsSchema = ComponentSchema(
                    properties = mapOf(
                        "message" to PropertyDefinition(type = PropertyType.STRING, description = "Message content"),
                        "author" to PropertyDefinition(type = PropertyType.STRING, description = "Message author"),
                        "timestamp" to PropertyDefinition(type = PropertyType.STRING, description = "Message timestamp"),
                        "isOwn" to PropertyDefinition(type = PropertyType.BOOLEAN, description = "Is own message")
                    ),
                    required = listOf("message", "author")
                ),
                stateSchema = ComponentSchema(
                    properties = mapOf(
                        "isRead" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false"),
                        "isSelected" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false")
                    )
                ),
                eventsSchema = ComponentSchema(
                    properties = mapOf(
                        "onPress" to PropertyDefinition(type = PropertyType.FUNCTION, description = "Handle message press"),
                        "onLongPress" to PropertyDefinition(type = PropertyType.FUNCTION, description = "Handle long press")
                    )
                ),
                metadata = ComponentMetadata(description = "Chat message display component")
            ),
            UIComponent(
                id = "user-avatar",
                name = "UserAvatar",
                category = ComponentCategory.DISPLAY,
                isStateful = true,
                propsSchema = ComponentSchema(
                    properties = mapOf(
                        "userId" to PropertyDefinition(type = PropertyType.STRING, description = "User identifier"),
                        "imageUrl" to PropertyDefinition(type = PropertyType.STRING, description = "Avatar image URL"),
                        "size" to PropertyDefinition(type = PropertyType.DIMENSION, defaultValue = "40"),
                        "showOnlineStatus" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "true")
                    ),
                    required = listOf("userId")
                ),
                stateSchema = ComponentSchema(
                    properties = mapOf(
                        "isOnline" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false"),
                        "lastSeen" to PropertyDefinition(type = PropertyType.STRING, isOptional = true)
                    )
                ),
                metadata = ComponentMetadata(description = "User avatar with status indicator")
            ),
            UIComponent(
                id = "navigation-tab",
                name = "NavigationTab",
                category = ComponentCategory.NAVIGATION,
                isStateful = true,
                propsSchema = ComponentSchema(
                    properties = mapOf(
                        "title" to PropertyDefinition(type = PropertyType.STRING, description = "Tab title"),
                        "icon" to PropertyDefinition(type = PropertyType.STRING, description = "Tab icon name"),
                        "badge" to PropertyDefinition(type = PropertyType.NUMBER, description = "Badge count", isOptional = true),
                        "isDisabled" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false")
                    ),
                    required = listOf("title", "icon")
                ),
                stateSchema = ComponentSchema(
                    properties = mapOf(
                        "isActive" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false"),
                        "hasNotification" to PropertyDefinition(type = PropertyType.BOOLEAN, defaultValue = "false")
                    )
                ),
                eventsSchema = ComponentSchema(
                    properties = mapOf(
                        "onSelect" to PropertyDefinition(type = PropertyType.FUNCTION, description = "Handle tab selection")
                    )
                ),
                metadata = ComponentMetadata(description = "Bottom navigation tab component")
            )
        )
    }
}

/**
 * Component category classification
 */
@Serializable
enum class ComponentCategory {
    INPUT,
    DISPLAY,
    NAVIGATION,
    LAYOUT,
    FEEDBACK,
    OVERLAY,
    MEDIA,
    DATA,
    FORM,
    INTERACTIVE
}

/**
 * Component size classification
 */
@Serializable
enum class ComponentSize {
    SMALL,
    MEDIUM,
    LARGE,
    EXTRA_LARGE
}

/**
 * Component schema definition
 */
@Serializable
data class ComponentSchema(
    val properties: Map<String, PropertyDefinition> = emptyMap(),
    val required: List<String> = emptyList(),
    val additionalProperties: Boolean = false
)

/**
 * Property definition for component schema
 */
@Serializable
data class PropertyDefinition(
    val type: PropertyType,
    val description: String? = null,
    val defaultValue: String? = null,
    val validation: PropertyValidation? = null,
    val isOptional: Boolean = true
)

/**
 * Property data types
 */
@Serializable
enum class PropertyType {
    STRING,
    NUMBER,
    BOOLEAN,
    ARRAY,
    OBJECT,
    FUNCTION,
    COLOR,
    DIMENSION
}

/**
 * Property validation rules
 */
@Serializable
data class PropertyValidation(
    val minValue: Double? = null,
    val maxValue: Double? = null,
    val pattern: String? = null,
    val enumValues: List<String>? = null
)

/**
 * Component metadata
 */
@Serializable
data class ComponentMetadata(
    val author: String = "Tchat Team",
    val createdAt: Long = System.currentTimeMillis(),
    val updatedAt: Long = System.currentTimeMillis(),
    val description: String? = null,
    val documentation: String? = null,
    val examples: List<String> = emptyList(),
    val tags: List<String> = emptyList(),
    val isDeprecated: Boolean = false,
    val deprecationMessage: String? = null
)

/**
 * Component instance with runtime state
 */
@Serializable
data class UIComponentInstance(
    val instanceId: String,
    val componentId: String,
    val componentName: String,
    val state: MutableMap<String, @Contextual Any>,
    val props: Map<String, @Contextual Any> = emptyMap(),
    val version: String,
    val createdAt: Long,
    var lastUpdated: Long
) {
    val id: String get() = instanceId

    /**
     * Update component state
     */
    fun updateState(newState: Map<String, Any>) {
        state.clear()
        state.putAll(newState)
        lastUpdated = System.currentTimeMillis()
    }

    /**
     * Merge state with existing state
     */
    fun mergeState(partialState: Map<String, Any>) {
        state.putAll(partialState)
        lastUpdated = System.currentTimeMillis()
    }
}