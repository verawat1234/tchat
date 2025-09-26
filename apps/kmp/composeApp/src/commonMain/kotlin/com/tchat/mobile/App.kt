package com.tchat.mobile

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography
import com.tchat.mobile.navigation.Screen
import com.tchat.mobile.navigation.MainTab
import com.tchat.mobile.navigation.TchatNavigation
import com.tchat.mobile.navigation.rememberScreenNavigationState
import com.tchat.mobile.screens.*

@Composable
fun App() {
    val navigationState = rememberScreenNavigationState()

    // Determine current tab based on the screen
    val currentTab = when (navigationState.currentScreen) {
        is Screen.Chat, is Screen.ChatDetail -> MainTab.CHAT
        is Screen.Store, is Screen.ProductDetail, is Screen.ShopDetail, is Screen.LiveStream -> MainTab.STORE
        is Screen.Social, is Screen.UserProfile -> MainTab.SOCIAL
        is Screen.Video, is Screen.VideoDetail -> MainTab.VIDEO
        is Screen.More, is Screen.Settings, is Screen.EditProfile -> MainTab.MORE
        is Screen.Search, is Screen.QRScanner, is Screen.Notifications -> MainTab.CHAT // Header actions from Chat
        is Screen.Reviews -> MainTab.STORE // Reviews belong to Store tab
    }

    // Check if we should show bottom navigation (hide on detailed screens)
    val showBottomNav = when (navigationState.currentScreen) {
        is Screen.Chat, is Screen.Store, is Screen.Social, is Screen.Video, is Screen.More -> true
        else -> false
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
                    TchatNavigation(
                        currentTab = currentTab,
                        onTabSelected = { tab ->
                            navigationState.navigateTo(tab.screen)
                        }
                    )
                }
            }
        ) { paddingValues ->
            when (val screen = navigationState.currentScreen) {
                is Screen.Chat -> ChatScreen(
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
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.Social -> SocialScreen(
                    onUserClick = { userId ->
                        navigationState.navigateTo(Screen.UserProfile(userId))
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.Video -> VideoScreen(
                    onVideoClick = { videoId ->
                        navigationState.navigateTo(Screen.VideoDetail(videoId))
                    },
                    modifier = Modifier.padding(paddingValues)
                )

                is Screen.More -> MoreScreen(
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
                    onBackClick = {
                        if (!navigationState.navigateBack()) {
                            navigationState.navigateTo(Screen.Video)
                        }
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
            }
        }
    }
}