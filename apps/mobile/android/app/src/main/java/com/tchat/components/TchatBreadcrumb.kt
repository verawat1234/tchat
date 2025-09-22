package com.tchat.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
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
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Breadcrumb navigation component following Tchat design system
 */
@Composable
fun TchatBreadcrumb(
    items: List<TchatBreadcrumbItem>,
    modifier: Modifier = Modifier,
    style: TchatBreadcrumbStyle = TchatBreadcrumbStyle.Default,
    size: TchatBreadcrumbSize = TchatBreadcrumbSize.Medium,
    separator: ImageVector = Icons.Default.ChevronRight,
    maxItems: Int? = null,
    showHome: Boolean = false,
    homeIcon: ImageVector = Icons.Default.Home,
    onItemClick: ((TchatBreadcrumbItem) -> Unit)? = null
) {
    val hapticFeedback = LocalHapticFeedback.current

    val displayItems = remember(items, maxItems) {
        when {
            maxItems == null || items.size <= maxItems -> items
            maxItems <= 2 -> items.takeLast(maxItems)
            else -> {
                val firstItem = items.first()
                val lastItems = items.takeLast(maxItems - 2)
                val ellipsisItem = TchatBreadcrumbItem(
                    id = "ellipsis",
                    title = "...",
                    isClickable = false
                )
                listOf(firstItem, ellipsisItem) + lastItems
            }
        }
    }

    Row(
        modifier = modifier
            .horizontalScroll(rememberScrollState())
            .padding(horizontal = Spacing.sm),
        horizontalArrangement = Arrangement.spacedBy(size.spacing),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Home icon (if enabled)
        if (showHome) {
            HomeButton(
                icon = homeIcon,
                size = size,
                onClick = {
                    val homeItem = TchatBreadcrumbItem(
                        id = "home",
                        title = "Home",
                        icon = homeIcon,
                        isClickable = true
                    )
                    onItemClick?.invoke(homeItem)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                }
            )

            if (items.isNotEmpty()) {
                SeparatorIcon(
                    icon = separator,
                    size = size
                )
            }
        }

        // Breadcrumb items
        displayItems.forEachIndexed { index, item ->
            BreadcrumbItemView(
                item = item,
                isLast = index == displayItems.lastIndex,
                style = style,
                size = size,
                onClick = {
                    item.action?.invoke()
                    onItemClick?.invoke(item)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                }
            )

            if (index < displayItems.lastIndex) {
                SeparatorIcon(
                    icon = separator,
                    size = size
                )
            }
        }
    }
}

@Composable
private fun HomeButton(
    icon: ImageVector,
    size: TchatBreadcrumbSize,
    onClick: () -> Unit
) {
    Icon(
        imageVector = icon,
        contentDescription = "Home",
        tint = Colors.textSecondary,
        modifier = Modifier
            .size(size.iconSize)
            .clickable { onClick() }
    )
}

@Composable
private fun BreadcrumbItemView(
    item: TchatBreadcrumbItem,
    isLast: Boolean,
    style: TchatBreadcrumbStyle,
    size: TchatBreadcrumbSize,
    onClick: () -> Unit
) {
    val itemColor = when (style) {
        TchatBreadcrumbStyle.Default -> {
            when {
                isLast -> Colors.textPrimary
                item.isClickable -> Colors.primary
                else -> Colors.textSecondary
            }
        }
        TchatBreadcrumbStyle.Compact, TchatBreadcrumbStyle.Minimal -> {
            when {
                isLast -> Colors.textPrimary
                item.isClickable -> Colors.textSecondary
                else -> Colors.textTertiary
            }
        }
    }

    val backgroundColor = when {
        style == TchatBreadcrumbStyle.Default && isLast -> Colors.surface
        else -> Color.Transparent
    }

    val clickableModifier = if (item.isClickable && !isLast) {
        Modifier.clickable { onClick() }
    } else {
        Modifier
    }

    Row(
        modifier = Modifier
            .then(clickableModifier)
            .background(
                color = backgroundColor,
                shape = RoundedCornerShape(4.dp)
            )
            .padding(
                horizontal = if (style == TchatBreadcrumbStyle.Default) size.spacing else 0.dp,
                vertical = if (style == TchatBreadcrumbStyle.Default) (size.spacing / 2) else 0.dp
            ),
        horizontalArrangement = Arrangement.spacedBy(size.spacing / 2),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Icon
        item.icon?.let { icon ->
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = itemColor,
                modifier = Modifier.size(size.iconSize)
            )
        }

        // Title
        Text(
            text = item.title,
            fontSize = size.fontSize,
            fontWeight = if (isLast) FontWeight.Medium else FontWeight.Normal,
            color = itemColor,
            maxLines = 1
        )
    }
}

@Composable
private fun SeparatorIcon(
    icon: ImageVector,
    size: TchatBreadcrumbSize
) {
    Icon(
        imageVector = icon,
        contentDescription = null,
        tint = Colors.textTertiary,
        modifier = Modifier.size(size.iconSize - 2.dp)
    )
}

/**
 * Breadcrumb item data class
 */
data class TchatBreadcrumbItem(
    val id: String,
    val title: String,
    val icon: ImageVector? = null,
    val isClickable: Boolean = true,
    val action: (() -> Unit)? = null
)

/**
 * Breadcrumb style definitions
 */
enum class TchatBreadcrumbStyle {
    Default,
    Compact,
    Minimal
}

/**
 * Breadcrumb size definitions
 */
enum class TchatBreadcrumbSize(
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val iconSize: androidx.compose.ui.unit.Dp,
    val spacing: androidx.compose.ui.unit.Dp
) {
    Small(
        fontSize = 12.sp,
        iconSize = 12.dp,
        spacing = 4.dp
    ),
    Medium(
        fontSize = 14.sp,
        iconSize = 14.dp,
        spacing = 6.dp
    ),
    Large(
        fontSize = 16.sp,
        iconSize = 16.dp,
        spacing = 8.dp
    )
}

/**
 * Utility functions for creating breadcrumbs
 */
object TchatBreadcrumbUtils {
    fun fromPath(
        path: String,
        separator: String = "/",
        onPathClick: ((String) -> Unit)? = null
    ): List<TchatBreadcrumbItem> {
        val components = path.split(separator).filter { it.isNotEmpty() }

        return components.mapIndexed { index, component ->
            val fullPath = components.take(index + 1).joinToString(separator)

            TchatBreadcrumbItem(
                id = fullPath,
                title = component,
                isClickable = true,
                action = { onPathClick?.invoke(fullPath) }
            )
        }
    }

    fun fromNavigation(
        pathComponents: List<String>,
        onNavigateBack: ((Int) -> Unit)? = null
    ): List<TchatBreadcrumbItem> {
        return pathComponents.mapIndexed { index, component ->
            TchatBreadcrumbItem(
                id = index.toString(),
                title = component,
                isClickable = index < pathComponents.lastIndex,
                action = { onNavigateBack?.invoke(index) }
            )
        }
    }
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatBreadcrumbPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Default breadcrumb
        TchatBreadcrumb(
            items = listOf(
                TchatBreadcrumbItem(
                    id = "dashboard",
                    title = "Dashboard",
                    icon = Icons.Default.Dashboard
                ),
                TchatBreadcrumbItem(
                    id = "projects",
                    title = "Projects",
                    icon = Icons.Default.Folder
                ),
                TchatBreadcrumbItem(
                    id = "mobile-app",
                    title = "Mobile App",
                    isClickable = false
                )
            ),
            style = TchatBreadcrumbStyle.Default,
            showHome = true
        )

        Divider()

        // Compact breadcrumb
        TchatBreadcrumb(
            items = listOf(
                TchatBreadcrumbItem(
                    id = "users",
                    title = "Users"
                ),
                TchatBreadcrumbItem(
                    id = "profile",
                    title = "Profile"
                ),
                TchatBreadcrumbItem(
                    id = "settings",
                    title = "Settings",
                    isClickable = false
                )
            ),
            style = TchatBreadcrumbStyle.Compact,
            size = TchatBreadcrumbSize.Small
        )

        Divider()

        // Minimal breadcrumb with truncation
        TchatBreadcrumb(
            items = listOf(
                TchatBreadcrumbItem(
                    id = "level1",
                    title = "Level 1"
                ),
                TchatBreadcrumbItem(
                    id = "level2",
                    title = "Level 2"
                ),
                TchatBreadcrumbItem(
                    id = "level3",
                    title = "Level 3"
                ),
                TchatBreadcrumbItem(
                    id = "level4",
                    title = "Level 4"
                ),
                TchatBreadcrumbItem(
                    id = "level5",
                    title = "Current Page",
                    isClickable = false
                )
            ),
            style = TchatBreadcrumbStyle.Minimal,
            maxItems = 4
        )

        Divider()

        // Path-based breadcrumb
        TchatBreadcrumb(
            items = TchatBreadcrumbUtils.fromPath(
                "/projects/mobile-app/src/components"
            ),
            style = TchatBreadcrumbStyle.Default,
            size = TchatBreadcrumbSize.Medium
        )

        Divider()

        // Large breadcrumb with custom separator
        TchatBreadcrumb(
            items = listOf(
                TchatBreadcrumbItem(
                    id = "docs",
                    title = "Documentation",
                    icon = Icons.Default.Description
                ),
                TchatBreadcrumbItem(
                    id = "api",
                    title = "API Reference",
                    icon = Icons.Default.Link
                ),
                TchatBreadcrumbItem(
                    id = "endpoints",
                    title = "Endpoints",
                    isClickable = false
                )
            ),
            style = TchatBreadcrumbStyle.Default,
            size = TchatBreadcrumbSize.Large,
            separator = Icons.Default.ArrowForward,
            showHome = true,
            homeIcon = Icons.Default.Home
        )
    }
}