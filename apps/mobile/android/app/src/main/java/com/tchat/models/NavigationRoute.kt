package com.tchat.models

import kotlinx.serialization.Serializable
import java.util.Date
import java.util.UUID

/**
 * Core navigation route entity with cross-platform compatibility
 */
@Serializable
data class NavigationRoute(
    val id: String,
    val path: String,
    val title: String,
    val component: String,
    val parameters: Map<String, RouteParameter> = emptyMap(),
    val isDeepLinkable: Boolean = true,
    val platformRestrictions: List<String> = emptyList(),
    val parentRouteId: String? = null,
    val accessLevel: AccessLevel = AccessLevel.PUBLIC,
    val metadata: RouteMetadata = RouteMetadata()
) {

    /**
     * Check if route is available on current platform
     */
    val isAvailableOnCurrentPlatform: Boolean
        get() = platformRestrictions.isEmpty() || platformRestrictions.contains("android")

    /**
     * Generate full route path with parameters
     */
    fun fullPath(parameters: Map<String, Any> = emptyMap()): String {
        var fullPath = path

        for ((key, value) in parameters) {
            fullPath = fullPath.replace(":$key", value.toString())
        }

        return fullPath
    }

    /**
     * Check if route requires authentication
     */
    val requiresAuthentication: Boolean
        get() = accessLevel == AccessLevel.PRIVATE || accessLevel == AccessLevel.ADMIN

    /**
     * Get route hierarchy depth
     */
    val depth: Int
        get() = path.split("/").filter { it.isNotEmpty() }.size

    companion object {
        /**
         * Default application routes
         */
        val defaultRoutes = listOf(
            NavigationRoute(
                id = "chat",
                path = "/chat",
                title = "Chat",
                component = "ChatView",
                metadata = RouteMetadata(description = "Main chat interface")
            ),
            NavigationRoute(
                id = "chat-user",
                path = "/chat/user/:userId",
                title = "User Chat",
                component = "UserChatView",
                parameters = mapOf(
                    "userId" to RouteParameter(
                        name = "userId",
                        type = ParameterType.UUID,
                        validation = ValidationRule(pattern = "[0-9a-fA-F-]{36}")
                    )
                ),
                parentRouteId = "chat",
                metadata = RouteMetadata(description = "Direct user chat")
            ),
            NavigationRoute(
                id = "store",
                path = "/store",
                title = "Store",
                component = "StoreView",
                metadata = RouteMetadata(description = "Shopping interface")
            ),
            NavigationRoute(
                id = "store-products",
                path = "/store/products",
                title = "Products",
                component = "ProductsView",
                parentRouteId = "store",
                metadata = RouteMetadata(description = "Product listing")
            ),
            NavigationRoute(
                id = "social",
                path = "/social",
                title = "Social",
                component = "SocialView",
                metadata = RouteMetadata(description = "Social feed")
            ),
            NavigationRoute(
                id = "social-feed",
                path = "/social/feed",
                title = "Feed",
                component = "FeedView",
                parentRouteId = "social",
                metadata = RouteMetadata(description = "Activity feed")
            ),
            NavigationRoute(
                id = "video",
                path = "/video",
                title = "Video",
                component = "VideoView",
                metadata = RouteMetadata(description = "Video calls")
            ),
            NavigationRoute(
                id = "video-call",
                path = "/video/call/:callId",
                title = "Video Call",
                component = "VideoCallView",
                parameters = mapOf(
                    "callId" to RouteParameter(
                        name = "callId",
                        type = ParameterType.STRING,
                        validation = ValidationRule(minLength = 8, maxLength = 32)
                    )
                ),
                parentRouteId = "video",
                accessLevel = AccessLevel.PRIVATE,
                metadata = RouteMetadata(description = "Active video call")
            ),
            NavigationRoute(
                id = "more",
                path = "/more",
                title = "More",
                component = "MoreView",
                metadata = RouteMetadata(description = "Additional options")
            ),
            NavigationRoute(
                id = "more-settings",
                path = "/more/settings",
                title = "Settings",
                component = "SettingsView",
                parentRouteId = "more",
                accessLevel = AccessLevel.PRIVATE,
                metadata = RouteMetadata(description = "App settings")
            )
        )
    }
}

/**
 * Route parameter definition
 */
@Serializable
data class RouteParameter(
    val name: String,
    val type: ParameterType,
    val isRequired: Boolean = true,
    val defaultValue: String? = null,
    val validation: ValidationRule? = null
)

/**
 * Parameter data types
 */
@Serializable
enum class ParameterType {
    STRING,
    INTEGER,
    UUID,
    BOOLEAN,
    URL
}


/**
 * Validation rules for parameters
 */
@Serializable
data class ValidationRule(
    val pattern: String? = null,
    val minLength: Int? = null,
    val maxLength: Int? = null,
    val allowedValues: List<String>? = null
)

/**
 * Route metadata for additional information
 */
@Serializable
data class RouteMetadata(
    val createdAt: Long = System.currentTimeMillis(),
    val updatedAt: Long = System.currentTimeMillis(),
    val version: String = "1.0.0",
    val description: String? = null,
    val tags: List<String> = emptyList(),
    val analyticsEnabled: Boolean = true,
    val cacheStrategy: CacheStrategy = CacheStrategy.DEFAULT
)

/**
 * Caching strategy for route data
 */
@Serializable
enum class CacheStrategy {
    DEFAULT,
    NO_CACHE,
    AGGRESSIVE,
    CONDITIONAL
}