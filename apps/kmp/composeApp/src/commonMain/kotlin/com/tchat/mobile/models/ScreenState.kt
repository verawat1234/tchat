package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import kotlinx.datetime.*

/**
 * T025: ScreenState and NavigationTab models
 *
 * Core navigation and screen state management for 5-tab architecture.
 * Supports Chat, Store, Social, Video, and More tabs with deep linking.
 */

/**
 * Main navigation tabs in the application
 */
enum class NavigationTab(
    val id: String,
    val displayName: String,
    val iconName: String,
    val order: Int
) {
    CHAT("chat", "Chat", "chat_icon", 0),
    STORE("store", "Store", "store_icon", 1),
    SOCIAL("social", "Social", "social_icon", 2),
    VIDEO("video", "Video", "video_icon", 3),
    MORE("more", "More", "more_icon", 4);

    companion object {
        fun fromId(id: String): NavigationTab? = values().find { it.id == id }
        fun getDefault(): NavigationTab = CHAT
        fun getAllTabs(): List<NavigationTab> = values().toList().sortedBy { it.order }
    }
}

/**
 * Screen state management for navigation and UI persistence
 */
@Serializable
data class ScreenState(
    val currentTab: String = NavigationTab.CHAT.id,
    val tabStates: Map<String, TabState> = emptyMap(),
    val navigationHistory: List<NavigationHistoryItem> = emptyList(),
    val deepLinkPending: String? = null,
    val modalStack: List<ModalInfo> = emptyList(),
    val globalLoading: Boolean = false,
    val globalError: ErrorState? = null,
    val networkStatus: NetworkStatus = NetworkStatus.CONNECTED,
    val lastActiveTimestamp: String? = null,
    val sessionTimeout: Long? = null,
    val backgroundedAt: String? = null,
    val metadata: Map<String, String> = emptyMap()
)

/**
 * Individual tab state for persistence across navigation
 */
@Serializable
data class TabState(
    val tabId: String,
    val currentRoute: String? = null,
    val scrollPosition: ScrollPosition? = null,
    val searchQuery: String? = null,
    val filters: Map<String, String> = emptyMap(),
    val selectedItems: List<String> = emptyList(),
    val lastRefreshed: String? = null,
    val cachedData: Map<String, String> = emptyMap(),
    val viewState: Map<String, String> = emptyMap(),
    val hasUnreadContent: Boolean = false,
    val badgeCount: Int = 0
)

/**
 * Scroll position for preserving user's position in lists
 */
@Serializable
data class ScrollPosition(
    val offset: Float = 0f,
    val itemIndex: Int = 0,
    val itemOffset: Float = 0f,
    val listId: String? = null
)

/**
 * Navigation history for back button and navigation tracking
 */
@Serializable
data class NavigationHistoryItem(
    val route: String,
    val tabId: String,
    val timestamp: String,
    val params: Map<String, String> = emptyMap(),
    val title: String? = null
)

/**
 * Modal/overlay information for stack management
 */
@Serializable
data class ModalInfo(
    val modalId: String,
    val type: ModalType,
    val data: Map<String, String> = emptyMap(),
    val dismissible: Boolean = true,
    val priority: Int = 0
)

enum class ModalType {
    ALERT,
    CONFIRMATION,
    ACTION_SHEET,
    BOTTOM_SHEET,
    FULL_SCREEN,
    POPUP,
    TOAST,
    LOADING
}

/**
 * Global error state management
 */
@Serializable
data class ErrorState(
    val code: String,
    val message: String,
    val details: String? = null,
    val recoverable: Boolean = true,
    val timestamp: String,
    val context: Map<String, String> = emptyMap(),
    val retryCount: Int = 0,
    val maxRetries: Int = 3
)

/**
 * Network connectivity status
 */
enum class NetworkStatus {
    CONNECTED,
    DISCONNECTED,
    SLOW,
    METERED,
    UNKNOWN
}

/**
 * Deep linking support
 */
@Serializable
data class DeepLink(
    val url: String,
    val tab: String,
    val route: String? = null,
    val params: Map<String, String> = emptyMap(),
    val requiresAuth: Boolean = true,
    val fallbackRoute: String? = null
)

/**
 * Tab-specific route definitions
 */
object TabRoutes {
    object Chat {
        const val HOME = "chat/home"
        const val SESSION = "chat/session/{sessionId}"
        const val NEW_CHAT = "chat/new"
        const val SETTINGS = "chat/settings"
        const val ARCHIVE = "chat/archive"
    }

    object Store {
        const val HOME = "store/home"
        const val PRODUCT = "store/product/{productId}"
        const val CATEGORY = "store/category/{categoryId}"
        const val CART = "store/cart"
        const val CHECKOUT = "store/checkout"
        const val ORDERS = "store/orders"
        const val WISHLIST = "store/wishlist"
    }

    object Social {
        const val FEED = "social/feed"
        const val PROFILE = "social/profile/{userId}"
        const val POST = "social/post/{postId}"
        const val FRIENDS = "social/friends"
        const val NOTIFICATIONS = "social/notifications"
        const val GROUPS = "social/groups"
        const val EVENTS = "social/events"
    }

    object Video {
        const val HOME = "video/home"
        const val WATCH = "video/watch/{videoId}"
        const val LIVE = "video/live/{streamId}"
        const val UPLOAD = "video/upload"
        const val PLAYLIST = "video/playlist/{playlistId}"
        const val TRENDING = "video/trending"
    }

    object More {
        const val HOME = "more/home"
        const val PROFILE = "more/profile"
        const val SETTINGS = "more/settings"
        const val HELP = "more/help"
        const val ABOUT = "more/about"
        const val FEEDBACK = "more/feedback"
    }
}

/**
 * Navigation state management extensions
 */
fun ScreenState.getCurrentTabState(): TabState? = tabStates[currentTab]

fun ScreenState.withTabSwitch(newTab: NavigationTab): ScreenState {
    val historyItem = NavigationHistoryItem(
        route = currentTab,
        tabId = currentTab,
        timestamp = Clock.System.now().toString()
    )
    return copy(
        currentTab = newTab.id,
        navigationHistory = (navigationHistory + historyItem).takeLast(50) // Keep last 50 entries
    )
}

fun ScreenState.withTabState(tabId: String, tabState: TabState): ScreenState = copy(
    tabStates = tabStates + (tabId to tabState)
)

fun ScreenState.withError(error: ErrorState?): ScreenState = copy(globalError = error)

fun ScreenState.withLoading(loading: Boolean): ScreenState = copy(globalLoading = loading)

fun ScreenState.withModal(modal: ModalInfo): ScreenState = copy(
    modalStack = (modalStack + modal).sortedByDescending { it.priority }
)

fun ScreenState.dismissModal(modalId: String): ScreenState = copy(
    modalStack = modalStack.filter { it.modalId != modalId }
)

fun ScreenState.clearModals(): ScreenState = copy(modalStack = emptyList())

fun TabState.withRoute(route: String): TabState = copy(currentRoute = route)

fun TabState.withScrollPosition(position: ScrollPosition): TabState = copy(scrollPosition = position)

fun TabState.withSearchQuery(query: String?): TabState = copy(searchQuery = query)

fun TabState.withBadgeCount(count: Int): TabState = copy(
    badgeCount = count,
    hasUnreadContent = count > 0
)

fun TabState.markRefreshed(): TabState = copy(
    lastRefreshed = Clock.System.now().toString()
)

/**
 * Deep link processing utilities
 */
object DeepLinkHandler {
    fun parseDeepLink(url: String): DeepLink? {
        // Parse tchat:// URLs and HTTP(S) URLs
        val uri = try {
            if (url.startsWith("tchat://")) {
                url.removePrefix("tchat://")
            } else if (url.startsWith("https://tchat.com/") || url.startsWith("http://tchat.com/")) {
                url.substringAfter("tchat.com/")
            } else {
                return null
            }
        } catch (e: Exception) {
            return null
        }

        val parts = uri.split("/")
        if (parts.isEmpty()) return null

        val tab = parts[0]
        val route = if (parts.size > 1) parts.drop(1).joinToString("/") else null

        return DeepLink(
            url = url,
            tab = tab,
            route = route,
            requiresAuth = tab != "public"
        )
    }

    fun buildDeepLink(tab: String, route: String? = null, params: Map<String, String> = emptyMap()): String {
        val baseUrl = "tchat://$tab"
        val fullRoute = if (route != null) "$baseUrl/$route" else baseUrl

        return if (params.isNotEmpty()) {
            val queryString = params.map { "${it.key}=${it.value}" }.joinToString("&")
            "$fullRoute?$queryString"
        } else {
            fullRoute
        }
    }
}

/**
 * Screen state validation
 */
fun ScreenState.validate(): ScreenState {
    val validTab = NavigationTab.fromId(currentTab) ?: NavigationTab.getDefault()
    val cleanedModals = modalStack.distinctBy { modal: ModalInfo -> modal.modalId }
    val cleanedHistory = navigationHistory.takeLast(50)

    return copy(
        currentTab = validTab.id,
        modalStack = cleanedModals,
        navigationHistory = cleanedHistory,
        tabStates = tabStates.filterKeys { NavigationTab.fromId(it) != null }
    )
}