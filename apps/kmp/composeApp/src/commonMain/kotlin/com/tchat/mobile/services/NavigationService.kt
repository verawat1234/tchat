package com.tchat.mobile.services

/**
 * Navigation Service - Handles app navigation
 *
 * Manages navigation to different screens and deep linking
 * Mock implementation for navigation actions
 */

enum class NavigationDestination {
    // Core Screens
    HOME,
    PROFILE,
    SEARCH,
    SETTINGS,
    NOTIFICATIONS,

    // Social Features
    POST_DETAIL,
    COMMENTS,
    USER_PROFILE,
    HASHTAG_FEED,
    HASHTAG_SEARCH,

    // Commerce
    PRODUCT_DETAIL,
    SHOP_DETAIL,
    CART,
    ORDERS,

    // Media
    IMAGE_VIEWER,
    VIDEO_PLAYER,
    CAMERA,

    // Authentication
    LOGIN,
    REGISTER,
    FORGOT_PASSWORD
}

data class NavigationAction(
    val destination: NavigationDestination,
    val parameters: Map<String, String> = emptyMap(),
    val clearStack: Boolean = false,
    val animation: NavigationAnimation = NavigationAnimation.SLIDE
)

enum class NavigationAnimation {
    SLIDE,
    FADE,
    SLIDE_UP,
    SCALE,
    NONE
}

interface NavigationService {
    suspend fun navigateTo(action: NavigationAction)
    suspend fun navigateBack()
    suspend fun navigateToComments(postId: String)
    suspend fun navigateToHashtagFeed(hashtag: String)
    suspend fun navigateToUserProfile(userId: String)
    suspend fun navigateToProductDetail(productId: String)
    suspend fun navigateToShopDetail(shopId: String)
    suspend fun navigateToImageViewer(imageUrl: String, postId: String? = null)
    suspend fun navigateToVideoPlayer(videoUrl: String, postId: String? = null)
    suspend fun openHashtagSearch(initialQuery: String? = null)
    suspend fun openShareSheet(postId: String)
    suspend fun handleDeepLink(url: String): Boolean
}

class MockNavigationService : NavigationService {

    // Mock navigation state
    private val navigationStack = mutableListOf<NavigationAction>()
    private var currentDestination: NavigationDestination = NavigationDestination.HOME

    override suspend fun navigateTo(action: NavigationAction) {
        println("🧭 Navigating to ${action.destination} with params: ${action.parameters}")

        if (action.clearStack) {
            navigationStack.clear()
        } else {
            navigationStack.add(
                NavigationAction(
                    destination = currentDestination,
                    parameters = emptyMap()
                )
            )
        }

        currentDestination = action.destination
        logNavigation(action)
    }

    override suspend fun navigateBack() {
        if (navigationStack.isNotEmpty()) {
            val previousAction = navigationStack.removeLastOrNull()
            if (previousAction != null) {
                currentDestination = previousAction.destination
                println("🔙 Navigating back to ${previousAction.destination}")
            }
        } else {
            println("🔙 Cannot navigate back - at root")
        }
    }

    override suspend fun navigateToComments(postId: String) {
        val action = NavigationAction(
            destination = NavigationDestination.COMMENTS,
            parameters = mapOf("postId" to postId),
            animation = NavigationAnimation.SLIDE_UP
        )
        navigateTo(action)
        println("💬 Opening comments for post: $postId")
    }

    override suspend fun navigateToHashtagFeed(hashtag: String) {
        val cleanHashtag = hashtag.removePrefix("#")
        val action = NavigationAction(
            destination = NavigationDestination.HASHTAG_FEED,
            parameters = mapOf("hashtag" to cleanHashtag),
            animation = NavigationAnimation.SLIDE
        )
        navigateTo(action)
        println("🏷️ Opening hashtag feed for: #$cleanHashtag")
    }

    override suspend fun navigateToUserProfile(userId: String) {
        val action = NavigationAction(
            destination = NavigationDestination.USER_PROFILE,
            parameters = mapOf("userId" to userId),
            animation = NavigationAnimation.SLIDE
        )
        navigateTo(action)
        println("👤 Opening user profile: $userId")
    }

    override suspend fun navigateToProductDetail(productId: String) {
        val action = NavigationAction(
            destination = NavigationDestination.PRODUCT_DETAIL,
            parameters = mapOf("productId" to productId),
            animation = NavigationAnimation.SLIDE
        )
        navigateTo(action)
        println("🛍️ Opening product detail: $productId")
    }

    override suspend fun navigateToShopDetail(shopId: String) {
        val action = NavigationAction(
            destination = NavigationDestination.SHOP_DETAIL,
            parameters = mapOf("shopId" to shopId),
            animation = NavigationAnimation.SLIDE
        )
        navigateTo(action)
        println("🏪 Opening shop detail: $shopId")
    }

    override suspend fun navigateToImageViewer(imageUrl: String, postId: String?) {
        val parameters = mutableMapOf("imageUrl" to imageUrl)
        if (postId != null) {
            parameters["postId"] = postId
        }

        val action = NavigationAction(
            destination = NavigationDestination.IMAGE_VIEWER,
            parameters = parameters,
            animation = NavigationAnimation.FADE
        )
        navigateTo(action)
        println("🖼️ Opening image viewer: $imageUrl")
    }

    override suspend fun navigateToVideoPlayer(videoUrl: String, postId: String?) {
        val parameters = mutableMapOf("videoUrl" to videoUrl)
        if (postId != null) {
            parameters["postId"] = postId
        }

        val action = NavigationAction(
            destination = NavigationDestination.VIDEO_PLAYER,
            parameters = parameters,
            animation = NavigationAnimation.SCALE
        )
        navigateTo(action)
        println("🎥 Opening video player: $videoUrl")
    }

    override suspend fun openHashtagSearch(initialQuery: String?) {
        val parameters = if (initialQuery != null) {
            mapOf("query" to initialQuery.removePrefix("#"))
        } else {
            emptyMap()
        }

        val action = NavigationAction(
            destination = NavigationDestination.HASHTAG_SEARCH,
            parameters = parameters,
            animation = NavigationAnimation.SLIDE_UP
        )
        navigateTo(action)
        println("🔍 Opening hashtag search with query: ${initialQuery ?: "empty"}")
    }

    override suspend fun openShareSheet(postId: String) {
        println("📤 Opening share sheet for post: $postId")
        // Share sheet is typically a modal, not a navigation destination
        // This would trigger the share modal in the UI
    }

    override suspend fun handleDeepLink(url: String): Boolean {
        println("🔗 Handling deep link: $url")

        // Mock deep link parsing
        return when {
            url.contains("/posts/") -> {
                val postId = url.substringAfterLast("/posts/").substringBefore("?")
                navigateTo(NavigationAction(
                    destination = NavigationDestination.POST_DETAIL,
                    parameters = mapOf("postId" to postId),
                    clearStack = true
                ))
                true
            }

            url.contains("/users/") -> {
                val userId = url.substringAfterLast("/users/").substringBefore("?")
                navigateToUserProfile(userId)
                true
            }

            url.contains("/products/") -> {
                val productId = url.substringAfterLast("/products/").substringBefore("?")
                navigateToProductDetail(productId)
                true
            }

            url.contains("/shops/") -> {
                val shopId = url.substringAfterLast("/shops/").substringBefore("?")
                navigateToShopDetail(shopId)
                true
            }

            url.contains("/hashtags/") -> {
                val hashtag = url.substringAfterLast("/hashtags/").substringBefore("?")
                navigateToHashtagFeed(hashtag)
                true
            }

            else -> {
                println("❌ Unknown deep link format: $url")
                false
            }
        }
    }

    // Helper functions
    private fun logNavigation(action: NavigationAction) {
        val paramString = if (action.parameters.isNotEmpty()) {
            " (${action.parameters.entries.joinToString { "${it.key}=${it.value}" }})"
        } else {
            ""
        }

        println("📱 Navigation: ${action.destination}${paramString} [${action.animation}]")

        // Log navigation stack depth
        println("   Stack depth: ${navigationStack.size}")
    }

    // Public methods for debugging/testing
    fun getCurrentDestination(): NavigationDestination = currentDestination

    fun getNavigationStackSize(): Int = navigationStack.size

    fun clearNavigationStack() {
        navigationStack.clear()
        println("🗑️ Navigation stack cleared")
    }

    // Deep link URL builders
    fun buildPostUrl(postId: String): String = "https://tchat.app/posts/$postId"

    fun buildUserUrl(userId: String): String = "https://tchat.app/users/$userId"

    fun buildProductUrl(productId: String): String = "https://tchat.app/products/$productId"

    fun buildShopUrl(shopId: String): String = "https://tchat.app/shops/$shopId"

    fun buildHashtagUrl(hashtag: String): String = "https://tchat.app/hashtags/${hashtag.removePrefix("#")}"
}

/**
 * Navigation Helper Functions
 */
object NavigationHelper {
    fun extractHashtagsFromText(text: String): List<String> {
        val hashtagRegex = "#\\w+".toRegex()
        return hashtagRegex.findAll(text)
            .map { it.value }
            .distinct()
            .toList()
    }

    fun extractMentionsFromText(text: String): List<String> {
        val mentionRegex = "@\\w+".toRegex()
        return mentionRegex.findAll(text)
            .map { it.value.removePrefix("@") }
            .distinct()
            .toList()
    }

    fun sanitizeHashtag(hashtag: String): String {
        return hashtag.removePrefix("#")
            .lowercase()
            .replace(Regex("[^a-zA-Z0-9_]"), "")
    }

    fun isValidDeepLink(url: String): Boolean {
        val validDomains = listOf("tchat.app", "www.tchat.app", "app.tchat.com")
        val validPaths = listOf("/posts/", "/users/", "/products/", "/shops/", "/hashtags/")

        return validDomains.any { domain -> url.contains(domain) } &&
                validPaths.any { path -> url.contains(path) }
    }
}