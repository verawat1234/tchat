package com.tchat.models

import android.net.Uri
import kotlinx.serialization.Serializable
import kotlinx.serialization.Contextual
import java.util.Date
import java.util.UUID

/**
 * Deep link entity for handling and resolving deep link navigation
 */
@Serializable
data class DeepLink(
    val id: String = UUID.randomUUID().toString(),
    val url: String,
    val scheme: String,
    val host: String? = null,
    val path: String,
    val queryParameters: Map<String, String> = emptyMap(),
    val fragment: String? = null,
    val platform: String = "android",
    val createdAt: Long = System.currentTimeMillis(),
    val metadata: DeepLinkMetadata = DeepLinkMetadata()
) {

    /**
     * Check if deep link is valid
     */
    val isValid: Boolean
        get() = scheme == "tchat" && path.isNotEmpty()

    /**
     * Get route ID from path
     */
    val routeId: String
        get() = path.trim('/').replace("/", "/")

    /**
     * Get path components
     */
    val pathComponents: List<String>
        get() = path.split("/").filter { it.isNotEmpty() }

    /**
     * Check if requires authentication
     */
    val requiresAuthentication: Boolean
        get() = metadata.accessLevel == AccessLevel.PRIVATE || metadata.accessLevel == AccessLevel.ADMIN

    /**
     * Get expiration status
     */
    val isExpired: Boolean
        get() {
            val expiresAt = metadata.expiresAt ?: return false
            return System.currentTimeMillis() > expiresAt
        }

    /**
     * Generate full URL
     */
    val fullURL: String
        get() {
            val uri = Uri.Builder()
                .scheme(scheme)
                .authority(host)
                .path(path)

            queryParameters.forEach { (key, value) ->
                uri.appendQueryParameter(key, value)
            }

            fragment?.let { uri.fragment(it) }

            return uri.build().toString()
        }

    /**
     * Create resolution request for this deep link
     */
    fun createResolutionRequest(userId: String): DeepLinkResolutionRequest {
        return DeepLinkResolutionRequest(
            url = url,
            platform = platform,
            userId = userId
        )
    }

    /**
     * Check if deep link matches a route pattern
     */
    fun matchesRoute(routePattern: String): Boolean {
        val pattern = routePattern.replace(Regex(":([^/]+)"), "([^/]+)")
        val regex = Regex("^$pattern$")
        return regex.matches(routeId)
    }

    /**
     * Extract parameters from route pattern
     */
    fun extractParameters(routePattern: String): Map<String, String> {
        val pathComponents = this.pathComponents
        val patternComponents = routePattern.split("/").filter { it.isNotEmpty() }

        val parameters = mutableMapOf<String, String>()

        for ((index, component) in patternComponents.withIndex()) {
            if (component.startsWith(":") && index < pathComponents.size) {
                val paramName = component.drop(1)
                parameters[paramName] = pathComponents[index]
            }
        }

        // Add query parameters
        parameters.putAll(queryParameters)

        return parameters
    }

    companion object {
        /**
         * Create deep link from URL string
         */
        fun fromURL(urlString: String): DeepLink? {
            return try {
                val uri = Uri.parse(urlString)
                val scheme = uri.scheme ?: return null

                val queryParams = mutableMapOf<String, String>()
                uri.queryParameterNames.forEach { name ->
                    uri.getQueryParameter(name)?.let { value ->
                        queryParams[name] = value
                    }
                }

                DeepLink(
                    url = urlString,
                    scheme = scheme,
                    host = uri.host,
                    path = uri.path ?: "",
                    queryParameters = queryParams,
                    fragment = uri.fragment,
                    platform = "android"
                )
            } catch (e: Exception) {
                null
            }
        }

        /**
         * Create deep link for route with parameters
         */
        fun forRoute(
            routeId: String,
            parameters: Map<String, String> = emptyMap(),
            platform: String = "android"
        ): DeepLink {
            val path = "/$routeId"
            val url = "tchat://$path"

            return DeepLink(
                url = url,
                scheme = "tchat",
                host = null,
                path = path,
                queryParameters = parameters,
                fragment = null,
                platform = platform
            )
        }

        /**
         * Common deep link patterns for the application
         */
        val commonPatterns = mapOf(
            "chat" to "/chat",
            "chat-user" to "/chat/user/:userId",
            "store" to "/store",
            "store-product" to "/store/product/:productId",
            "social" to "/social",
            "social-post" to "/social/post/:postId",
            "video" to "/video",
            "video-call" to "/video/call/:callId",
            "more" to "/more",
            "settings" to "/more/settings"
        )

        /**
         * Create deep link for chat with user
         */
        fun chatWithUser(userId: String): DeepLink {
            return forRoute("chat/user/$userId", platform = "android")
        }

        /**
         * Create deep link for store product
         */
        fun storeProduct(productId: String): DeepLink {
            return forRoute("store/product/$productId", platform = "android")
        }

        /**
         * Create deep link for video call
         */
        fun videoCall(callId: String): DeepLink {
            return forRoute("video/call/$callId", platform = "android")
        }

        /**
         * Create deep link for social post
         */
        fun socialPost(postId: String): DeepLink {
            return forRoute("social/post/$postId", platform = "android")
        }
    }
}

/**
 * Deep link metadata
 */
@Serializable
data class DeepLinkMetadata(
    val source: DeepLinkSource = DeepLinkSource.APP,
    val accessLevel: AccessLevel = AccessLevel.PUBLIC,
    val expiresAt: Long? = null,
    val campaign: String? = null,
    val referrer: String? = null,
    val customData: Map<String, String> = emptyMap()
)

/**
 * Deep link source
 */
@Serializable
enum class DeepLinkSource {
    APP,
    WEB,
    PUSH,
    SMS,
    EMAIL,
    SOCIAL,
    QR,
    CLIPBOARD
}


/**
 * Deep link resolution request
 */
@Serializable
data class DeepLinkResolutionRequest(
    val url: String,
    val platform: String,
    val userId: String
)

/**
 * Deep link resolution response
 */
@Serializable
data class DeepLinkResolution(
    val routeId: String,
    val parameters: Map<String, @Contextual Any>,
    val isValid: Boolean,
    val fallbackAction: String? = null,
    val requiresAuth: Boolean = false
)