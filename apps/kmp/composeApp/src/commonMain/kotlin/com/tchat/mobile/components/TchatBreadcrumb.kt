package com.tchat.mobile.components

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ChevronRight
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.MoreHoriz
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * TchatBreadcrumb - Navigation breadcrumbs component
 *
 * Features:
 * - Clickable navigation path with custom separators
 * - Icon support and overflow handling for long paths
 * - Platform-native link styling with hover states
 * - Responsive design with automatic truncation
 */

data class BreadcrumbItem(
    val label: String,
    val path: String? = null,
    val icon: ImageVector? = null,
    val isClickable: Boolean = true
)

enum class BreadcrumbSeparator {
    CHEVRON,
    SLASH,
    PIPE,
    DOT
}

@Composable
fun TchatBreadcrumb(
    items: List<BreadcrumbItem>,
    onNavigate: (String) -> Unit = {},
    separator: BreadcrumbSeparator = BreadcrumbSeparator.CHEVRON,
    maxVisibleItems: Int = 5,
    showHomeIcon: Boolean = true,
    modifier: Modifier = Modifier
) {
    val processedItems = remember(items, maxVisibleItems) {
        if (items.size <= maxVisibleItems) {
            items
        } else {
            // Show first item, ellipsis, and last few items
            val keepLast = maxVisibleItems - 2 // Reserve space for first item and ellipsis
            listOf(
                items.first(),
                BreadcrumbItem("…", null, Icons.Default.MoreHoriz, false)
            ) + items.takeLast(keepLast)
        }
    }

    LazyRow(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.Start,
        verticalAlignment = Alignment.CenterVertically
    ) {
        items(processedItems) { item ->
            BreadcrumbItemContent(
                item = item,
                onNavigate = onNavigate,
                showIcon = showHomeIcon && item == items.firstOrNull(),
                isLast = item == processedItems.last()
            )

            if (item != processedItems.last()) {
                BreadcrumbSeparatorContent(separator = separator)
            }
        }
    }
}

@Composable
private fun BreadcrumbItemContent(
    item: BreadcrumbItem,
    onNavigate: (String) -> Unit,
    showIcon: Boolean,
    isLast: Boolean
) {
    val textColor = if (isLast) {
        TchatColors.onSurface
    } else {
        TchatColors.onSurfaceVariant
    }

    val clickableModifier = if (item.isClickable && item.path != null) {
        Modifier.clickable { onNavigate(item.path) }
    } else {
        Modifier
    }

    Row(
        modifier = clickableModifier
            .clip(RoundedCornerShape(4.dp))
            .padding(horizontal = 4.dp, vertical = 2.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        // Show home icon for first item or custom icon
        if (showIcon && item.icon != null) {
            Icon(
                imageVector = item.icon,
                contentDescription = null,
                tint = textColor,
                modifier = Modifier.size(16.dp)
            )
        } else if (showIcon && item == item) { // First item
            Icon(
                imageVector = Icons.Default.Home,
                contentDescription = "Home",
                tint = textColor,
                modifier = Modifier.size(16.dp)
            )
        }

        Text(
            text = item.label,
            style = MaterialTheme.typography.bodyMedium,
            color = textColor,
            maxLines = 1,
            overflow = TextOverflow.Ellipsis,
            modifier = Modifier.widthIn(max = 120.dp)
        )
    }
}

@Composable
private fun BreadcrumbSeparatorContent(
    separator: BreadcrumbSeparator
) {
    when (separator) {
        BreadcrumbSeparator.CHEVRON -> {
            Icon(
                imageVector = Icons.Default.ChevronRight,
                contentDescription = null,
                tint = TchatColors.onSurfaceVariant,
                modifier = Modifier
                    .size(16.dp)
                    .padding(horizontal = 2.dp)
            )
        }
        BreadcrumbSeparator.SLASH -> {
            Text(
                text = "/",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }
        BreadcrumbSeparator.PIPE -> {
            Text(
                text = "|",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }
        BreadcrumbSeparator.DOT -> {
            Text(
                text = "•",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }
    }
}