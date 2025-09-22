package com.tchat.navigation

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.navigation.NavController
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.tchat.components.TchatCard
import com.tchat.components.TchatCardVariant
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Main tab navigation composable following Tchat design system
 */
@Composable
fun TabNavigationComposable() {
    val navController = rememberNavController()
    var isBottomBarVisible by remember { mutableStateOf(true) }

    Scaffold(
        bottomBar = {
            AnimatedVisibility(
                visible = isBottomBarVisible,
                enter = slideInVertically(initialOffsetY = { it }),
                exit = slideOutVertically(targetOffsetY = { it })
            ) {
                TchatBottomNavigationBar(navController = navController)
            }
        }
    ) { paddingValues ->
        NavHost(
            navController = navController,
            startDestination = TabDestination.Chat.route,
            modifier = Modifier.padding(paddingValues)
        ) {
            composable(TabDestination.Chat.route) {
                ChatTabScreen(
                    onBottomBarVisibilityChange = { isBottomBarVisible = it }
                )
            }
            composable(TabDestination.Store.route) {
                StoreTabScreen(
                    onBottomBarVisibilityChange = { isBottomBarVisible = it }
                )
            }
            composable(TabDestination.Social.route) {
                SocialTabScreen(
                    onBottomBarVisibilityChange = { isBottomBarVisible = it }
                )
            }
            composable(TabDestination.Video.route) {
                VideoTabScreen(
                    onBottomBarVisibilityChange = { isBottomBarVisible = it }
                )
            }
            composable(TabDestination.More.route) {
                MoreTabScreen(
                    onBottomBarVisibilityChange = { isBottomBarVisible = it }
                )
            }
        }
    }
}

/**
 * Bottom navigation bar component
 */
@Composable
fun TchatBottomNavigationBar(
    navController: NavController,
    modifier: Modifier = Modifier
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentDestination = navBackStackEntry?.destination
    val hapticFeedback = LocalHapticFeedback.current

    NavigationBar(
        modifier = modifier
            .shadow(
                elevation = 8.dp,
                shape = RoundedCornerShape(topStart = 16.dp, topEnd = 16.dp)
            ),
        containerColor = Colors.tabBarBackground,
        contentColor = Colors.textPrimary,
        tonalElevation = 0.dp
    ) {
        TabDestination.entries.forEach { destination ->
            val isSelected = currentDestination?.route == destination.route

            NavigationBarItem(
                icon = {
                    Icon(
                        imageVector = if (isSelected) destination.selectedIcon else destination.icon,
                        contentDescription = destination.title,
                        modifier = Modifier.size(24.dp)
                    )
                },
                label = {
                    Text(
                        text = destination.title,
                        fontSize = 12.sp,
                        fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Medium
                    )
                },
                selected = isSelected,
                onClick = {
                    // Haptic feedback
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )

                    navController.navigate(destination.route) {
                        // Pop up to the start destination of the graph to
                        // avoid building up a large stack of destinations
                        popUpTo(navController.graph.findStartDestination().id) {
                            saveState = true
                        }
                        // Avoid multiple copies of the same destination when
                        // reselecting the same item
                        launchSingleTop = true
                        // Restore state when reselecting a previously selected item
                        restoreState = true
                    }
                },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = Colors.tabSelected,
                    selectedTextColor = Colors.tabSelected,
                    unselectedIconColor = Colors.tabUnselected,
                    unselectedTextColor = Colors.tabUnselected,
                    indicatorColor = Colors.primary.copy(alpha = 0.1f)
                )
            )
        }
    }
}

/**
 * Tab destination enum
 */
enum class TabDestination(
    val route: String,
    val title: String,
    val icon: ImageVector,
    val selectedIcon: ImageVector
) {
    Chat(
        route = "chat",
        title = "Chat",
        icon = Icons.Default.ChatBubbleOutline,
        selectedIcon = Icons.Default.ChatBubble
    ),
    Store(
        route = "store",
        title = "Store",
        icon = Icons.Default.ShoppingBag,
        selectedIcon = Icons.Default.ShoppingBag
    ),
    Social(
        route = "social",
        title = "Social",
        icon = Icons.Default.Group,
        selectedIcon = Icons.Default.Group
    ),
    Video(
        route = "video",
        title = "Video",
        icon = Icons.Default.PlayCircleOutline,
        selectedIcon = Icons.Default.PlayCircle
    ),
    More(
        route = "more",
        title = "More",
        icon = Icons.Default.MoreHoriz,
        selectedIcon = Icons.Default.MoreHoriz
    )
}

// MARK: - Tab Screen Composables

/**
 * Chat tab screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatTabScreen(
    onBottomBarVisibilityChange: (Boolean) -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = "Chat",
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Colors.navigationBackground
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.md),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            Text(
                text = "Messages and conversations",
                color = Colors.textSecondary
            )

            Spacer(modifier = Modifier.weight(1f))

            // Placeholder content
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                Icon(
                    imageVector = Icons.Default.ChatBubble,
                    contentDescription = null,
                    tint = Colors.primary,
                    modifier = Modifier.size(64.dp)
                )

                Text(
                    text = "Your chat interface will be here",
                    color = Colors.textTertiary
                )
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

/**
 * Store tab screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StoreTabScreen(
    onBottomBarVisibilityChange: (Boolean) -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = "Store",
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Colors.navigationBackground
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.md),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            Text(
                text = "Browse and purchase items",
                color = Colors.textSecondary
            )

            Spacer(modifier = Modifier.weight(1f))

            // Placeholder content
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                Icon(
                    imageVector = Icons.Default.ShoppingBag,
                    contentDescription = null,
                    tint = Colors.primary,
                    modifier = Modifier.size(64.dp)
                )

                Text(
                    text = "Your store interface will be here",
                    color = Colors.textTertiary
                )
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

/**
 * Social tab screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SocialTabScreen(
    onBottomBarVisibilityChange: (Boolean) -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = "Social",
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Colors.navigationBackground
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.md),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            Text(
                text = "Connect with friends",
                color = Colors.textSecondary
            )

            Spacer(modifier = Modifier.weight(1f))

            // Placeholder content
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                Icon(
                    imageVector = Icons.Default.Group,
                    contentDescription = null,
                    tint = Colors.primary,
                    modifier = Modifier.size(64.dp)
                )

                Text(
                    text = "Your social interface will be here",
                    color = Colors.textTertiary
                )
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

/**
 * Video tab screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoTabScreen(
    onBottomBarVisibilityChange: (Boolean) -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = "Video",
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Colors.navigationBackground
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.md),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            Text(
                text = "Watch and share videos",
                color = Colors.textSecondary
            )

            Spacer(modifier = Modifier.weight(1f))

            // Placeholder content
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                Icon(
                    imageVector = Icons.Default.PlayCircle,
                    contentDescription = null,
                    tint = Colors.primary,
                    modifier = Modifier.size(64.dp)
                )

                Text(
                    text = "Your video interface will be here",
                    color = Colors.textTertiary
                )
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

/**
 * More tab screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MoreTabScreen(
    onBottomBarVisibilityChange: (Boolean) -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = "More",
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Colors.navigationBackground
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.md),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            Text(
                text = "Settings and additional features",
                color = Colors.textSecondary
            )

            Spacer(modifier = Modifier.height(Spacing.md))

            // Settings items
            Column(
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                TchatCard(
                    variant = TchatCardVariant.Outlined,
                    isInteractive = true,
                    onClick = { println("Settings tapped") }
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
                    ) {
                        Icon(
                            imageVector = Icons.Default.Settings,
                            contentDescription = null,
                            tint = Colors.primary
                        )
                        Text(
                            text = "Settings",
                            fontWeight = FontWeight.Medium,
                            modifier = Modifier.weight(1f)
                        )
                        Icon(
                            imageVector = Icons.Default.ChevronRight,
                            contentDescription = null,
                            tint = Colors.textTertiary
                        )
                    }
                }

                TchatCard(
                    variant = TchatCardVariant.Outlined,
                    isInteractive = true,
                    onClick = { println("Profile tapped") }
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
                    ) {
                        Icon(
                            imageVector = Icons.Default.AccountCircle,
                            contentDescription = null,
                            tint = Colors.primary
                        )
                        Text(
                            text = "Profile",
                            fontWeight = FontWeight.Medium,
                            modifier = Modifier.weight(1f)
                        )
                        Icon(
                            imageVector = Icons.Default.ChevronRight,
                            contentDescription = null,
                            tint = Colors.textTertiary
                        )
                    }
                }

                TchatCard(
                    variant = TchatCardVariant.Outlined,
                    isInteractive = true,
                    onClick = { println("About tapped") }
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
                    ) {
                        Icon(
                            imageVector = Icons.Default.Info,
                            contentDescription = null,
                            tint = Colors.primary
                        )
                        Text(
                            text = "About",
                            fontWeight = FontWeight.Medium,
                            modifier = Modifier.weight(1f)
                        )
                        Icon(
                            imageVector = Icons.Default.ChevronRight,
                            contentDescription = null,
                            tint = Colors.textTertiary
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

// Preview
@Preview(showBackground = true)
@Composable
fun TabNavigationPreview() {
    TabNavigationComposable()
}