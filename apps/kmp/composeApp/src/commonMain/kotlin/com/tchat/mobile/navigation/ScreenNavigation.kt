package com.tchat.mobile.navigation

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * Main Screen Navigation System
 * 5-tab architecture: Chat, Store, Social, Video, More
 */

sealed class Screen(val route: String) {
    // Main tab screens
    data object Chat : Screen("chat")
    data object Store : Screen("store")
    data object Social : Screen("social")
    data object Video : Screen("video")
    data object More : Screen("more")

    // Top bar screen
    data object Web : Screen("web")

    // Detail screens
    data class ChatDetail(val chatId: String, val chatName: String) : Screen("chat_detail/$chatId")
    data class ProductDetail(val productId: String) : Screen("product_detail/$productId")
    data class ShopDetail(val shopId: String) : Screen("shop_detail/$shopId")
    data class LiveStream(val streamId: String) : Screen("live_stream/$streamId")
    data class UserProfile(val userId: String) : Screen("user_profile/$userId")
    data class VideoDetail(val videoId: String) : Screen("video_detail/$videoId")
    data class EventDetail(val eventId: String) : Screen("event_detail/$eventId")
    data class CategoryEvents(val categoryId: String, val categoryName: String) : Screen("category_events/$categoryId")
    data class Reviews(val targetId: String, val targetType: String, val targetName: String) : Screen("reviews/$targetId/$targetType")
    data object Settings : Screen("settings")
    data object EditProfile : Screen("edit_profile")

    // Header action screens
    data object Search : Screen("search")
    data object QRScanner : Screen("qr_scanner")
    data object Notifications : Screen("notifications")

    // Content creation screens
    data object CreateChat : Screen("create_chat")
    data object CreateProduct : Screen("create_product")
    data object CreatePost : Screen("create_post")
    data object CreateVideo : Screen("create_video")
}

enum class MainTab(
    val screen: Screen,
    val title: String,
    val selectedIcon: ImageVector,
    val unselectedIcon: ImageVector
) {
    CHAT(Screen.Chat, "Chat", Icons.Default.Email, Icons.Default.Email),
    STORE(Screen.Store, "Store", Icons.Filled.ShoppingCart, Icons.Outlined.ShoppingCart),
    SOCIAL(Screen.Social, "Social", Icons.Filled.Person, Icons.Filled.Person),
    VIDEO(Screen.Video, "Video", Icons.Filled.PlayArrow, Icons.Outlined.PlayArrow),
    ADD(Screen.More, "+", Icons.Filled.Add, Icons.Outlined.Add)
}

@Composable
fun TchatNavigation(
    currentTab: MainTab,
    onTabSelected: (MainTab) -> Unit,
    modifier: Modifier = Modifier
) {
    NavigationBar(
        modifier = modifier,
        containerColor = TchatColors.surface,
        contentColor = TchatColors.onSurface
    ) {
        MainTab.entries.forEach { tab ->
            NavigationBarItem(
                selected = currentTab == tab,
                onClick = { onTabSelected(tab) },
                icon = {
                    Icon(
                        imageVector = if (currentTab == tab) tab.selectedIcon else tab.unselectedIcon,
                        contentDescription = tab.title
                    )
                },
                label = {
                    Text(
                        text = tab.title,
                        style = MaterialTheme.typography.labelSmall
                    )
                },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = TchatColors.primary,
                    selectedTextColor = TchatColors.primary,
                    unselectedIconColor = TchatColors.onSurfaceVariant,
                    unselectedTextColor = TchatColors.onSurfaceVariant,
                    indicatorColor = TchatColors.primary.copy(alpha = 0.12f)
                )
            )
        }
    }
}


/**
 * Screen State Management
 */
@Stable
class ScreenNavigationState {
    var currentScreen by mutableStateOf<Screen>(Screen.Chat)
        private set

    private val navigationHistory = mutableListOf<Screen>()

    fun navigateTo(screen: Screen) {
        // Add current screen to history before navigating
        if (currentScreen != screen) {
            navigationHistory.add(currentScreen)
        }
        currentScreen = screen
    }

    fun navigateBack(): Boolean {
        // Navigate back using history
        return if (navigationHistory.isNotEmpty()) {
            currentScreen = navigationHistory.removeLast()
            true
        } else {
            // If no history, check if we can navigate to parent screen
            when (currentScreen) {
                is Screen.Chat, is Screen.Store, is Screen.Social, is Screen.Video -> false
                is Screen.ChatDetail -> {
                    currentScreen = Screen.Chat
                    true
                }
                is Screen.ProductDetail, is Screen.ShopDetail, is Screen.LiveStream, is Screen.Reviews -> {
                    currentScreen = Screen.Store
                    true
                }
                is Screen.UserProfile -> {
                    currentScreen = Screen.Social
                    true
                }
                is Screen.EventDetail -> {
                    currentScreen = Screen.Social
                    true
                }
                is Screen.CategoryEvents -> {
                    currentScreen = Screen.Social
                    true
                }
                is Screen.VideoDetail -> {
                    currentScreen = Screen.Video
                    true
                }
                is Screen.More, is Screen.Settings, is Screen.EditProfile -> {
                    currentScreen = Screen.Chat // Navigate to Chat instead of More
                    true
                }
                is Screen.Search, is Screen.QRScanner, is Screen.Notifications -> {
                    currentScreen = Screen.Chat
                    true
                }
                is Screen.Web -> {
                    currentScreen = Screen.Chat
                    true
                }
                is Screen.CreateChat -> {
                    currentScreen = Screen.Chat
                    true
                }
                is Screen.CreateProduct -> {
                    currentScreen = Screen.Store
                    true
                }
                is Screen.CreatePost -> {
                    currentScreen = Screen.Social
                    true
                }
                is Screen.CreateVideo -> {
                    currentScreen = Screen.Video
                    true
                }
            }
        }
    }
}

@Composable
fun rememberScreenNavigationState(): ScreenNavigationState {
    return remember { ScreenNavigationState() }
}