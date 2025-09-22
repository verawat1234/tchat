import kotlinx.coroutines.delay
import java.util.Date

/**
 * Mock service for Navigation Sync testing
 */
class MockNavigationSyncService {

    suspend fun getNavigationRoutes(platform: String): NavigationRoutesResponse {
        // Simulate network delay
        delay(100)

        return NavigationRoutesResponse(
            platform = platform,
            routes = listOf(
                NavigationRoute(
                    id = "chat",
                    name = "Chat",
                    path = "/chat",
                    isTabRoute = true,
                    icon = "chat_icon"
                ),
                NavigationRoute(
                    id = "store",
                    name = "Store",
                    path = "/store",
                    isTabRoute = true,
                    icon = "store_icon"
                ),
                NavigationRoute(
                    id = "social",
                    name = "Social",
                    path = "/social",
                    isTabRoute = true,
                    icon = "social_icon"
                )
            ),
            lastUpdated = Date()
        )
    }

    suspend fun syncNavigationState(request: NavigationSyncRequest): NavigationSyncResponse {
        // Simulate network delay
        delay(100)

        return NavigationSyncResponse(
            success = true,
            message = "Navigation state synced successfully",
            currentRoute = request.currentRoute,
            timestamp = Date()
        )
    }
}

/**
 * Data models for testing
 */

data class NavigationRoutesResponse(
    val platform: String,
    val routes: List<NavigationRoute>,
    val lastUpdated: Date
)

data class NavigationRoute(
    val id: String,
    val name: String,
    val path: String,
    val isTabRoute: Boolean,
    val icon: String
)

data class NavigationSyncRequest(
    val platform: String,
    val currentRoute: String,
    val navigationHistory: List<String>,
    val timestamp: Date
)

data class NavigationSyncResponse(
    val success: Boolean,
    val message: String,
    val currentRoute: String,
    val timestamp: Date
)