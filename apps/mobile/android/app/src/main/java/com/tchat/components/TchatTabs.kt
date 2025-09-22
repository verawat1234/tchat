package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Tab navigation component following Tchat design system
 */
@Composable
fun TchatTabs(
    selectedTab: String,
    onTabSelected: (String) -> Unit,
    tabs: List<TchatTabItem>,
    modifier: Modifier = Modifier,
    style: TchatTabStyle = TchatTabStyle.Line,
    size: TchatTabSize = TchatTabSize.Medium,
    position: TchatTabPosition = TchatTabPosition.Top,
    showDivider: Boolean = true
) {
    val hapticFeedback = LocalHapticFeedback.current

    Column(modifier = modifier) {
        if (position == TchatTabPosition.Top) {
            TchatTabHeader(
                selectedTab = selectedTab,
                onTabSelected = { tabId ->
                    onTabSelected(tabId)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                },
                tabs = tabs,
                style = style,
                size = size
            )

            if (showDivider && style == TchatTabStyle.Line) {
                Divider(
                    thickness = 1.dp,
                    color = Colors.border
                )
            }

            TchatTabContent(
                selectedTab = selectedTab,
                tabs = tabs
            )
        } else {
            TchatTabContent(
                selectedTab = selectedTab,
                tabs = tabs
            )

            if (showDivider && style == TchatTabStyle.Line) {
                Divider(
                    thickness = 1.dp,
                    color = Colors.border
                )
            }

            TchatTabHeader(
                selectedTab = selectedTab,
                onTabSelected = { tabId ->
                    onTabSelected(tabId)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                },
                tabs = tabs,
                style = style,
                size = size
            )
        }
    }
}

@Composable
private fun TchatTabHeader(
    selectedTab: String,
    onTabSelected: (String) -> Unit,
    tabs: List<TchatTabItem>,
    style: TchatTabStyle,
    size: TchatTabSize
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = Spacing.md),
        horizontalArrangement = Arrangement.spacedBy(0.dp)
    ) {
        tabs.forEach { tab ->
            TchatTabButton(
                tab = tab,
                isSelected = selectedTab == tab.id,
                onTabSelected = onTabSelected,
                style = style,
                size = size,
                modifier = Modifier.weight(1f)
            )
        }
    }
}

@Composable
private fun TchatTabButton(
    tab: TchatTabItem,
    isSelected: Boolean,
    onTabSelected: (String) -> Unit,
    style: TchatTabStyle,
    size: TchatTabSize,
    modifier: Modifier = Modifier
) {
    val backgroundAlpha by animateFloatAsState(
        targetValue = if (isSelected && style == TchatTabStyle.Pill) 0.1f else 0f,
        animationSpec = tween(200),
        label = "background_alpha"
    )

    val elevationAlpha by animateFloatAsState(
        targetValue = if (isSelected && style == TchatTabStyle.Card) 1f else 0f,
        animationSpec = tween(200),
        label = "elevation_alpha"
    )

    Box(
        modifier = modifier
            .height(size.height)
            .clip(
                when (style) {
                    TchatTabStyle.Line -> RoundedCornerShape(0.dp)
                    TchatTabStyle.Pill -> RoundedCornerShape(size.height / 2)
                    TchatTabStyle.Card -> RoundedCornerShape(8.dp)
                }
            )
            .background(
                when (style) {
                    TchatTabStyle.Line -> Color.Transparent
                    TchatTabStyle.Pill -> Colors.primary.copy(alpha = backgroundAlpha)
                    TchatTabStyle.Card -> if (isSelected) Colors.background else Colors.surface
                }
            )
            .shadow(
                elevation = if (style == TchatTabStyle.Card) (2.dp * elevationAlpha) else 0.dp,
                shape = RoundedCornerShape(8.dp)
            )
            .clickable(enabled = !tab.isDisabled) {
                onTabSelected(tab.id)
            },
        contentAlignment = Alignment.Center
    ) {
        Row(
            horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier.padding(horizontal = size.horizontalPadding)
        ) {
            // Icon
            tab.icon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = when {
                        tab.isDisabled -> Colors.textDisabled
                        isSelected -> Colors.primary
                        else -> Colors.textSecondary
                    },
                    modifier = Modifier.size(size.iconSize)
                )
            }

            // Title
            Text(
                text = tab.title,
                fontSize = size.fontSize,
                fontWeight = FontWeight.Medium,
                color = when {
                    tab.isDisabled -> Colors.textDisabled
                    style == TchatTabStyle.Card && isSelected -> Colors.textPrimary
                    isSelected -> Colors.primary
                    else -> Colors.textSecondary
                }
            )

            // Badge
            tab.badge?.let { badge ->
                Text(
                    text = badge,
                    fontSize = 10.sp,
                    color = Colors.textOnPrimary,
                    modifier = Modifier
                        .background(
                            color = Colors.error,
                            shape = CircleShape
                        )
                        .padding(horizontal = 6.dp, vertical = 2.dp)
                )
            }
        }

        // Line indicator
        if (style == TchatTabStyle.Line && isSelected) {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(2.dp)
                    .background(Colors.primary)
                    .align(Alignment.BottomCenter)
            )
        }
    }
}

@Composable
private fun TchatTabContent(
    selectedTab: String,
    tabs: List<TchatTabItem>
) {
    val selectedTabItem = tabs.find { it.id == selectedTab }

    selectedTabItem?.let { tab ->
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .animateContentSize(
                    animationSpec = tween(200)
                )
        ) {
            tab.content()
        }
    }
}

/**
 * Bottom tab bar component
 */
@Composable
fun TchatTabBar(
    selectedTab: String,
    onTabSelected: (String) -> Unit,
    tabs: List<TchatTabItem>,
    modifier: Modifier = Modifier,
    size: TchatTabSize = TchatTabSize.Medium
) {
    val hapticFeedback = LocalHapticFeedback.current

    Column(modifier = modifier) {
        Divider(
            thickness = 1.dp,
            color = Colors.border
        )

        Row(
            modifier = Modifier
                .fillMaxWidth()
                .background(Colors.background)
                .padding(horizontal = Spacing.md, vertical = Spacing.xs),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            tabs.forEach { tab ->
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    modifier = Modifier
                        .weight(1f)
                        .clickable(enabled = !tab.isDisabled) {
                            onTabSelected(tab.id)
                            hapticFeedback.performHapticFeedback(
                                androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                            )
                        }
                        .padding(vertical = Spacing.xs)
                ) {
                    Box {
                        tab.icon?.let { icon ->
                            Icon(
                                imageVector = icon,
                                contentDescription = null,
                                tint = when {
                                    tab.isDisabled -> Colors.textDisabled
                                    selectedTab == tab.id -> Colors.primary
                                    else -> Colors.textSecondary
                                },
                                modifier = Modifier.size(size.iconSize)
                            )
                        }

                        tab.badge?.let { badge ->
                            Text(
                                text = badge,
                                fontSize = 8.sp,
                                color = Colors.textOnPrimary,
                                modifier = Modifier
                                    .background(
                                        color = Colors.error,
                                        shape = CircleShape
                                    )
                                    .padding(horizontal = 4.dp, vertical = 2.dp)
                                    .offset(x = 12.dp, y = (-8).dp)
                            )
                        }
                    }

                    Spacer(modifier = Modifier.height(4.dp))

                    Text(
                        text = tab.title,
                        fontSize = 10.sp,
                        color = when {
                            tab.isDisabled -> Colors.textDisabled
                            selectedTab == tab.id -> Colors.primary
                            else -> Colors.textSecondary
                        }
                    )
                }
            }
        }
    }
}

/**
 * Tab item data class
 */
data class TchatTabItem(
    val id: String,
    val title: String,
    val icon: ImageVector? = null,
    val badge: String? = null,
    val isDisabled: Boolean = false,
    val content: @Composable () -> Unit
)

/**
 * Tab style definitions
 */
enum class TchatTabStyle {
    Line,
    Pill,
    Card
}

/**
 * Tab position definitions
 */
enum class TchatTabPosition {
    Top,
    Bottom
}

/**
 * Tab size definitions
 */
enum class TchatTabSize(
    val height: androidx.compose.ui.unit.Dp,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val iconSize: androidx.compose.ui.unit.Dp
) {
    Small(
        height = 32.dp,
        horizontalPadding = 8.dp,
        fontSize = 12.sp,
        iconSize = 14.dp
    ),
    Medium(
        height = 40.dp,
        horizontalPadding = 12.dp,
        fontSize = 14.sp,
        iconSize = 16.dp
    ),
    Large(
        height = 48.dp,
        horizontalPadding = 16.dp,
        fontSize = 16.sp,
        iconSize = 18.dp
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatTabsPreview() {
    var selectedTab by remember { mutableStateOf("tab1") }
    var selectedPill by remember { mutableStateOf("pill1") }
    var selectedCard by remember { mutableStateOf("card1") }
    var selectedBar by remember { mutableStateOf("home") }

    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Line tabs
        TchatTabs(
            selectedTab = selectedTab,
            onTabSelected = { selectedTab = it },
            tabs = listOf(
                TchatTabItem(
                    id = "tab1",
                    title = "Overview",
                    icon = Icons.Default.BarChart
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(200.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Overview Content")
                    }
                },
                TchatTabItem(
                    id = "tab2",
                    title = "Analytics",
                    icon = Icons.Default.TrendingUp
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(200.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Analytics Content")
                    }
                },
                TchatTabItem(
                    id = "tab3",
                    title = "Settings",
                    icon = Icons.Default.Settings,
                    badge = "2"
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(200.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Settings Content")
                    }
                }
            ),
            style = TchatTabStyle.Line
        )

        // Pill tabs
        TchatTabs(
            selectedTab = selectedPill,
            onTabSelected = { selectedPill = it },
            tabs = listOf(
                TchatTabItem(id = "pill1", title = "All") {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(100.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("All Items")
                    }
                },
                TchatTabItem(id = "pill2", title = "Active") {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(100.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Active Items")
                    }
                },
                TchatTabItem(id = "pill3", title = "Completed") {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(100.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Completed Items")
                    }
                }
            ),
            style = TchatTabStyle.Pill,
            size = TchatTabSize.Small
        )

        // Card tabs
        TchatTabs(
            selectedTab = selectedCard,
            onTabSelected = { selectedCard = it },
            tabs = listOf(
                TchatTabItem(
                    id = "card1",
                    title = "Profile",
                    icon = Icons.Default.Person
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(150.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Profile Content")
                    }
                },
                TchatTabItem(
                    id = "card2",
                    title = "Messages",
                    icon = Icons.Default.Message,
                    badge = "5"
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(150.dp)
                            .background(Colors.surface),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("Messages Content")
                    }
                }
            ),
            style = TchatTabStyle.Card,
            size = TchatTabSize.Large
        )

        Divider()

        // Tab bar
        TchatTabBar(
            selectedTab = selectedBar,
            onTabSelected = { selectedBar = it },
            tabs = listOf(
                TchatTabItem(id = "home", title = "Home", icon = Icons.Default.Home) { },
                TchatTabItem(id = "search", title = "Search", icon = Icons.Default.Search) { },
                TchatTabItem(id = "notifications", title = "Notifications", icon = Icons.Default.Notifications, badge = "3") { },
                TchatTabItem(id = "profile", title = "Profile", icon = Icons.Default.Person) { }
            )
        )
    }
}