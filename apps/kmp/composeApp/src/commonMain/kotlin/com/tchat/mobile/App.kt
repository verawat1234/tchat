package com.tchat.mobile

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography
import com.tchat.mobile.navigation.Screen
import com.tchat.mobile.navigation.MainTab
import com.tchat.mobile.navigation.TchatNavigation
import com.tchat.mobile.navigation.rememberScreenNavigationState
import com.tchat.mobile.repositories.ChatRepository
import com.tchat.mobile.services.SocialContentService
import com.tchat.mobile.services.ContentApiService
import com.tchat.mobile.screens.*
import org.koin.compose.koinInject

@Composable
fun App() {
    val navigationState = rememberScreenNavigationState()

    // Initialize ChatRepository - using koinInject for now
    val chatRepository: ChatRepository = koinInject()

    // Inject services using Koin
    val socialContentService: SocialContentService = koinInject()
    val contentApiService: ContentApiService = koinInject()

    // Video screen UI visibility state for controlling bottom nav opacity
    var isVideoUIVisible by remember { mutableStateOf(true) }

    // Determine current tab based on the screen
    val currentTab = when (navigationState.currentScreen) {
        is Screen.Chat, is Screen.ChatDetail -> MainTab.CHAT
        is Screen.Store, is Screen.ProductDetail, is Screen.ShopDetail, is Screen.LiveStream -> MainTab.STORE
        is Screen.Social, is Screen.UserProfile -> MainTab.SOCIAL
        is Screen.Video, is Screen.VideoDetail -> MainTab.VIDEO
        is Screen.More, is Screen.Settings, is Screen.EditProfile -> MainTab.ADD
        is Screen.Search, is Screen.QRScanner, is Screen.Notifications -> MainTab.CHAT // Header actions from Chat
        is Screen.Reviews -> MainTab.STORE // Reviews belong to Store tab
        is Screen.Web -> MainTab.CHAT // Default fallback to Chat
        is Screen.CreateChat -> MainTab.CHAT
        is Screen.CreateProduct -> MainTab.STORE
        is Screen.CreatePost -> MainTab.SOCIAL
        is Screen.CreateVideo -> MainTab.VIDEO
    }

    // Check if we should show bottom navigation (hide only on specific screens)
    val showBottomNav = when (navigationState.currentScreen) {
        // Always show on main tab screens
        is Screen.Chat, is Screen.Store, is Screen.Social, is Screen.Video, is Screen.More -> true

        // Hide on content creation screens (they have their own navigation)
        is Screen.CreateChat, is Screen.CreateProduct, is Screen.CreatePost, is Screen.CreateVideo -> false

        // Hide on web screen (has its own navigation)
        is Screen.Web -> false

        // Hide on detailed/modal screens that should be full-screen
        is Screen.ChatDetail, is Screen.VideoDetail -> false

        // Show on all other screens (Settings, Search, QR, Notifications, etc.)
        else -> true
    }

    MaterialTheme(
        typography = TchatTypography.typography,
        colorScheme = lightColorScheme(
            primary = TchatColors.primary,
            onPrimary = TchatColors.onPrimary,
            surface = TchatColors.surface,
            onSurface = TchatColors.onSurface,
            background = TchatColors.background,
            onBackground = TchatColors.onBackground
        )
    ) {
        Scaffold(
            bottomBar = {
                if (showBottomNav) {
                    // Apply opacity to bottom navigation when on Video screen with hidden UI
                    val navOpacity = if (currentTab == MainTab.VIDEO && !isVideoUIVisible) 0.3f else 1f

                    Box(modifier = Modifier.fillMaxWidth()) {
                        TchatNavigation(
                            currentTab = currentTab,
                            onTabSelected = { tab ->
                                if (tab == MainTab.ADD) {
                                    // Context-aware + button functionality
                                    when (currentTab) {
                                        MainTab.CHAT -> {
                                            // Create new chat/conversation
                                            navigationState.navigateTo(Screen.CreateChat)
                                        }
                                        MainTab.STORE -> {
                                            // Create new product/listing
                                            navigationState.navigateTo(Screen.CreateProduct)
                                        }
                                        MainTab.SOCIAL -> {
                                            // Create new post
                                            navigationState.navigateTo(Screen.CreatePost)
                                        }
                                        MainTab.VIDEO -> {
                                            // Create new video
                                            navigationState.navigateTo(Screen.CreateVideo)
                                        }
                                        MainTab.ADD -> {
                                            // Fallback to More screen
                                            navigationState.navigateTo(Screen.More)
                                        }
                                    }
                                } else {
                                    navigationState.navigateTo(tab.screen)
                                }
                            },
                            modifier = Modifier.alpha(navOpacity)
                        )
                    }
                }
            }
        ) { paddingValues ->
            when (val screen = navigationState.currentScreen) {
                is Screen.Chat -> ChatScreen(
                    chatRepository = chatRepository,
                    onChatClick = { chatId, chatName ->
                        navigationState.navigateTo(Screen.ChatDetail(chatId, chatName))
                    },
                    onSearchClick = {
                        navigationState.navigateTo(Screen.Search)
                    },
                    onQRScannerClick = {
                        navigationState.navigateTo(Screen.QRScanner)
                    },
                    onNotificationsClick = {
                        navigationState.navigateTo(Screen.Notifications)
                    },
                    onMoreClick = {
                        navigationState.navigateTo(Screen.More)
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.Store -> StoreScreen(
                    onProductClick = { productId ->
                        navigationState.navigateTo(Screen.ProductDetail(productId))
                    },
                    onShopClick = { shopId ->
                        navigationState.navigateTo(Screen.ShopDetail(shopId))
                    },
                    onLiveStreamClick = { streamId ->
                        navigationState.navigateTo(Screen.LiveStream(streamId))
                    },
                    onMoreClick = {
                        navigationState.navigateTo(Screen.More)
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.Social -> SocialScreen(
                    onUserClick = { userId ->
                        navigationState.navigateTo(Screen.UserProfile(userId))
                    },
                    onMoreClick = {
                        navigationState.navigateTo(Screen.More)
                    },
                    socialContentService = socialContentService,
                    contentApiService = contentApiService,
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.Video -> VideoScreen(
                    onVideoClick = { videoId ->
                        navigationState.navigateTo(Screen.VideoDetail(videoId))
                    },
                    onUIVisibilityChange = { isVisible ->
                        isVideoUIVisible = isVisible
                    },
                    onMoreClick = {
                        navigationState.navigateTo(Screen.More)
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.More -> MoreScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    },
                    onEditProfileClick = {
                        navigationState.navigateTo(Screen.EditProfile)
                    },
                    onSettingsClick = {
                        navigationState.navigateTo(Screen.Settings)
                    },
                    onUserProfileClick = { userId ->
                        navigationState.navigateTo(Screen.UserProfile(userId))
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                // Detailed screens
                is Screen.ChatDetail -> ChatDetailScreen(
                    chatId = screen.chatId,
                    chatName = screen.chatName,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    }
                )

                is Screen.ProductDetail -> ProductDetailScreen(
                    productId = screen.productId,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Store)
                        }
                    },
                    onShopClick = { shopId ->
                        navigationState.navigateTo(Screen.ShopDetail(shopId))
                    }
                )

                is Screen.UserProfile -> UserProfileScreen(
                    userId = screen.userId,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Social)
                        }
                    },
                    onEditClick = {
                        navigationState.navigateTo(Screen.EditProfile)
                    }
                )

                is Screen.VideoDetail -> VideoDetailScreen(
                    videoId = screen.videoId,
                    onNavigateBack = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Video)
                        }
                    },
                    onVideoClick = { videoId ->
                        navigationState.navigateTo(Screen.VideoDetail(videoId))
                    }
                )

                is Screen.Settings -> SettingsScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.More)
                        }
                    },
                    onEditProfileClick = {
                        navigationState.navigateTo(Screen.EditProfile)
                    }
                )

                is Screen.EditProfile -> EditProfileScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.More)
                        }
                    }
                )

                // Header action screens
                is Screen.Search -> SearchScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    }
                )

                is Screen.QRScanner -> QRScannerScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    }
                )

                is Screen.Notifications -> NotificationsScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    }
                )

                is Screen.Reviews -> ReviewScreen(
                    targetId = screen.targetId,
                    targetType = screen.targetType,
                    targetName = screen.targetName,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Store)
                        }
                    },
                    onUserClick = { userId ->
                        navigationState.navigateTo(Screen.UserProfile(userId))
                    },
                    onProductClick = { productId ->
                        navigationState.navigateTo(Screen.ProductDetail(productId))
                    },
                    onShopClick = { shopId ->
                        navigationState.navigateTo(Screen.ShopDetail(shopId))
                    }
                )

                // Store detail screens
                is Screen.ShopDetail -> ShopDetailScreen(
                    shopId = screen.shopId,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Store)
                        }
                    },
                    onProductClick = { productId ->
                        navigationState.navigateTo(Screen.ProductDetail(productId))
                    }
                )

                is Screen.LiveStream -> LiveStreamScreen(
                    streamId = screen.streamId,
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Store)
                        }
                    },
                    onProductClick = { productId ->
                        navigationState.navigateTo(Screen.ProductDetail(productId))
                    },
                    onShopClick = { shopId ->
                        navigationState.navigateTo(Screen.ShopDetail(shopId))
                    }
                )

                // Content creation screens
                is Screen.CreateChat -> CreateChatScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Chat)
                        }
                    }
                )

                is Screen.CreateProduct -> CreateProductScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Store)
                        }
                    }
                )

                is Screen.CreatePost -> CreatePostScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Social)
                        }
                    }
                )

                is Screen.CreateVideo -> CreateVideoScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Video)
                        }
                    }
                )

                is Screen.Web -> WebScreen(
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.More)
                        }
                    }
                )
            }
        }
    }
}