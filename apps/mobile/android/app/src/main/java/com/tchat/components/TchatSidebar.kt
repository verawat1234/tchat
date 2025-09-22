package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.zIndex
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Sidebar navigation component following Tchat design system
 */
@Composable
fun TchatSidebar(
    isOpen: Boolean,
    onOpenChange: (Boolean) -> Unit,
    items: List<TchatSidebarItem>,
    modifier: Modifier = Modifier,
    selectedItem: String? = null,
    onSelectionChange: ((String?) -> Unit)? = null,
    mode: TchatSidebarMode = TchatSidebarMode.Overlay,
    position: TchatSidebarPosition = TchatSidebarPosition.Leading,
    size: TchatSidebarSize = TchatSidebarSize.Standard,
    showOverlay: Boolean = true,
    allowDismiss: Boolean = true,
    header: (@Composable () -> Unit)? = null,
    footer: (@Composable () -> Unit)? = null
) {
    val density = LocalDensity.current
    val hapticFeedback = LocalHapticFeedback.current

    var expandedSections by remember { mutableStateOf(setOf<String>()) }
    var dragOffset by remember { mutableStateOf(0f) }

    val sidebarWidth = with(density) { size.width.toPx() }

    val overlayAlpha by animateFloatAsState(
        targetValue = if (isOpen && mode == TchatSidebarMode.Overlay && showOverlay) 0.5f else 0f,
        animationSpec = tween(300),
        label = "overlay_alpha"
    )

    val sidebarOffset by animateFloatAsState(
        targetValue = when {
            mode == TchatSidebarMode.Permanent -> 0f
            isOpen -> 0f
            position == TchatSidebarPosition.Leading -> -sidebarWidth
            else -> sidebarWidth
        } + dragOffset,
        animationSpec = tween(300),
        label = "sidebar_offset"
    )

    Box(modifier = modifier.fillMaxSize()) {
        // Overlay
        if (mode == TchatSidebarMode.Overlay && showOverlay) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = overlayAlpha))
                    .zIndex(1f)
                    .clickable(enabled = isOpen && allowDismiss) {
                        onOpenChange(false)
                    }
            )
        }

        // Sidebar
        Box(
            modifier = Modifier
                .fillMaxHeight()
                .width(size.width)
                .offset(x = with(density) { sidebarOffset.toDp() })
                .zIndex(2f)
                .then(
                    if (position == TchatSidebarPosition.Trailing) {
                        Modifier.align(Alignment.TopEnd)
                    } else {
                        Modifier.align(Alignment.TopStart)
                    }
                )
                .pointerInput(allowDismiss) {
                    if (allowDismiss) {
                        detectDragGestures(
                            onDragEnd = {
                                val threshold = sidebarWidth * 0.3f
                                val shouldClose = when (position) {
                                    TchatSidebarPosition.Leading -> dragOffset < -threshold
                                    TchatSidebarPosition.Trailing -> dragOffset > threshold
                                }

                                if (shouldClose) {
                                    onOpenChange(false)
                                }
                                dragOffset = 0f
                            }
                        ) { _, dragAmount ->
                            val maxDrag = sidebarWidth * 0.3f
                            dragOffset = when (position) {
                                TchatSidebarPosition.Leading -> maxOf(-maxDrag, minOf(0f, dragAmount.x))
                                TchatSidebarPosition.Trailing -> minOf(maxDrag, maxOf(0f, dragAmount.x))
                            }
                        }
                    }
                }
        ) {
            SidebarContent(
                items = items,
                selectedItem = selectedItem,
                expandedSections = expandedSections,
                onExpandedSectionsChange = { expandedSections = it },
                onSelectionChange = { itemId ->
                    onSelectionChange?.invoke(itemId)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                    if (mode == TchatSidebarMode.Overlay) {
                        onOpenChange(false)
                    }
                },
                position = position,
                header = header,
                footer = footer
            )
        }
    }
}

@Composable
private fun SidebarContent(
    items: List<TchatSidebarItem>,
    selectedItem: String?,
    expandedSections: Set<String>,
    onExpandedSectionsChange: (Set<String>) -> Unit,
    onSelectionChange: (String?) -> Unit,
    position: TchatSidebarPosition,
    header: (@Composable () -> Unit)?,
    footer: (@Composable () -> Unit)?
) {
    Column(
        modifier = Modifier
            .fillMaxHeight()
            .background(Colors.background)
            .border(
                width = 1.dp,
                color = Colors.border,
                shape = RoundedCornerShape(0.dp)
            )
    ) {
        // Header
        header?.let { headerContent ->
            Box(
                modifier = Modifier.padding(Spacing.md)
            ) {
                headerContent()
            }
        }

        // Navigation Items
        LazyColumn(
            modifier = Modifier.weight(1f),
            verticalArrangement = Arrangement.spacedBy(0.dp),
            contentPadding = PaddingValues(vertical = Spacing.sm)
        ) {
            items(items) { item ->
                SidebarItemView(
                    item = item,
                    selectedItem = selectedItem,
                    expandedSections = expandedSections,
                    onExpandedSectionsChange = onExpandedSectionsChange,
                    onSelectionChange = onSelectionChange,
                    level = 0
                )
            }
        }

        // Footer
        footer?.let { footerContent ->
            Box(
                modifier = Modifier.padding(Spacing.md)
            ) {
                footerContent()
            }
        }
    }
}

@Composable
private fun SidebarItemView(
    item: TchatSidebarItem,
    selectedItem: String?,
    expandedSections: Set<String>,
    onExpandedSectionsChange: (Set<String>) -> Unit,
    onSelectionChange: (String?) -> Unit,
    level: Int
) {
    val hapticFeedback = LocalHapticFeedback.current
    val isSelected = selectedItem == item.id
    val isExpanded = expandedSections.contains(item.id)

    Column {
        // Main item
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clickable(enabled = !item.isDisabled) {
                    if (item.children.isNotEmpty()) {
                        val newExpanded = if (isExpanded) {
                            expandedSections - item.id
                        } else {
                            expandedSections + item.id
                        }
                        onExpandedSectionsChange(newExpanded)
                    } else {
                        onSelectionChange(item.id)
                        item.action?.invoke()
                    }
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                }
                .background(
                    if (isSelected) Colors.primary.copy(alpha = 0.1f) else Color.Transparent
                )
                .padding(
                    start = Spacing.md + (level * 16).dp,
                    end = Spacing.md,
                    top = Spacing.sm,
                    bottom = Spacing.sm
                ),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            // Icon
            item.icon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = when {
                        item.isDisabled -> Colors.textDisabled
                        isSelected -> Colors.primary
                        else -> Colors.textSecondary
                    },
                    modifier = Modifier.size(20.dp)
                )
            }

            // Title
            Text(
                text = item.title,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = when {
                    item.isDisabled -> Colors.textDisabled
                    isSelected -> Colors.primary
                    else -> Colors.textPrimary
                },
                modifier = Modifier.weight(1f)
            )

            // Badge
            item.badge?.let { badge ->
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

            // Expand/collapse indicator
            if (item.children.isNotEmpty()) {
                Icon(
                    imageVector = if (isExpanded) Icons.Default.ExpandLess else Icons.Default.ChevronRight,
                    contentDescription = null,
                    tint = Colors.textSecondary,
                    modifier = Modifier.size(16.dp)
                )
            }
        }

        // Children (if expanded)
        AnimatedVisibility(
            visible = isExpanded && item.children.isNotEmpty(),
            enter = fadeIn() + expandVertically(),
            exit = fadeOut() + shrinkVertically()
        ) {
            Column {
                item.children.forEach { child ->
                    SidebarItemView(
                        item = child,
                        selectedItem = selectedItem,
                        expandedSections = expandedSections,
                        onExpandedSectionsChange = onExpandedSectionsChange,
                        onSelectionChange = onSelectionChange,
                        level = level + 1
                    )
                }
            }
        }
    }
}

/**
 * Sidebar item data class
 */
data class TchatSidebarItem(
    val id: String,
    val title: String,
    val icon: ImageVector? = null,
    val badge: String? = null,
    val isDisabled: Boolean = false,
    val children: List<TchatSidebarItem> = emptyList(),
    val action: (() -> Unit)? = null
)

/**
 * Sidebar mode definitions
 */
enum class TchatSidebarMode {
    Overlay,
    Push,
    Permanent
}

/**
 * Sidebar position definitions
 */
enum class TchatSidebarPosition {
    Leading,
    Trailing
}

/**
 * Sidebar size definitions
 */
enum class TchatSidebarSize(
    val width: androidx.compose.ui.unit.Dp
) {
    Compact(240.dp),
    Standard(280.dp),
    Wide(320.dp)
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatSidebarPreview() {
    var isOpen by remember { mutableStateOf(true) }
    var selectedItem by remember { mutableStateOf<String?>("dashboard") }

    Box(modifier = Modifier.fillMaxSize()) {
        // Main content
        Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Colors.surface)
                .padding(Spacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Text(
                text = "Main Content",
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.textPrimary
            )

            Spacer(modifier = Modifier.height(Spacing.md))

            Button(
                onClick = { isOpen = !isOpen }
            ) {
                Text("Toggle Sidebar")
            }
        }

        // Sidebar
        TchatSidebar(
            isOpen = isOpen,
            onOpenChange = { isOpen = it },
            selectedItem = selectedItem,
            onSelectionChange = { selectedItem = it },
            items = listOf(
                TchatSidebarItem(
                    id = "dashboard",
                    title = "Dashboard",
                    icon = Icons.Default.Dashboard
                ),
                TchatSidebarItem(
                    id = "messages",
                    title = "Messages",
                    icon = Icons.Default.Message,
                    badge = "5"
                ),
                TchatSidebarItem(
                    id = "projects",
                    title = "Projects",
                    icon = Icons.Default.Folder,
                    children = listOf(
                        TchatSidebarItem(
                            id = "active-projects",
                            title = "Active Projects",
                            icon = Icons.Default.Circle
                        ),
                        TchatSidebarItem(
                            id = "completed-projects",
                            title = "Completed",
                            icon = Icons.Default.CheckCircle
                        )
                    )
                ),
                TchatSidebarItem(
                    id = "team",
                    title = "Team",
                    icon = Icons.Default.People,
                    children = listOf(
                        TchatSidebarItem(
                            id = "members",
                            title = "Members",
                            icon = Icons.Default.Person
                        ),
                        TchatSidebarItem(
                            id = "roles",
                            title = "Roles & Permissions",
                            icon = Icons.Default.Key
                        )
                    )
                ),
                TchatSidebarItem(
                    id = "settings",
                    title = "Settings",
                    icon = Icons.Default.Settings
                )
            ),
            header = {
                Column(
                    verticalArrangement = Arrangement.spacedBy(Spacing.sm)
                ) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            imageVector = Icons.Default.Apps,
                            contentDescription = null,
                            tint = Colors.primary,
                            modifier = Modifier.size(24.dp)
                        )

                        Text(
                            text = "Tchat",
                            fontSize = 18.sp,
                            fontWeight = FontWeight.Bold,
                            color = Colors.textPrimary
                        )
                    }

                    Divider()
                }
            },
            footer = {
                Column(
                    verticalArrangement = Arrangement.spacedBy(Spacing.sm)
                ) {
                    Divider()

                    Row(
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            imageVector = Icons.Default.AccountCircle,
                            contentDescription = null,
                            tint = Colors.textSecondary,
                            modifier = Modifier.size(20.dp)
                        )

                        Column {
                            Text(
                                text = "John Doe",
                                fontSize = 12.sp,
                                fontWeight = FontWeight.Medium,
                                color = Colors.textPrimary
                            )

                            Text(
                                text = "john@example.com",
                                fontSize = 10.sp,
                                color = Colors.textSecondary
                            )
                        }
                    }
                }
            }
        )
    }
}