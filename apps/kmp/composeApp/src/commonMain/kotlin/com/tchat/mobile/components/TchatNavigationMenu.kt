package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowDropDown
import androidx.compose.material.icons.filled.ChevronRight
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.unit.DpOffset
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.PopupProperties
import com.tchat.mobile.designsystem.TchatColors

/**
 * TchatNavigationMenu - Dropdown navigation component
 *
 * Features:
 * - Hierarchical menu structure with nested submenus
 * - Keyboard navigation support and accessibility
 * - Custom menu item rendering with icons
 * - Platform-specific dropdown behavior
 */

data class NavigationMenuItem(
    val id: String,
    val label: String,
    val icon: ImageVector? = null,
    val enabled: Boolean = true,
    val children: List<NavigationMenuItem> = emptyList(),
    val action: (() -> Unit)? = null
)

@Composable
fun TchatNavigationMenu(
    trigger: @Composable () -> Unit,
    items: List<NavigationMenuItem>,
    onItemClick: (NavigationMenuItem) -> Unit = {},
    expanded: Boolean = false,
    onExpandedChange: (Boolean) -> Unit = {},
    modifier: Modifier = Modifier
) {
    var internalExpanded by remember { mutableStateOf(expanded) }

    LaunchedEffect(expanded) {
        internalExpanded = expanded
    }

    Box(modifier = modifier) {
        Box(
            modifier = Modifier.clickable {
                internalExpanded = !internalExpanded
                onExpandedChange(internalExpanded)
            }
        ) {
            trigger()
        }

        DropdownMenu(
            expanded = internalExpanded,
            onDismissRequest = {
                internalExpanded = false
                onExpandedChange(false)
            },
            modifier = Modifier
                .widthIn(min = 200.dp, max = 300.dp)
                .background(
                    color = TchatColors.surface,
                    shape = RoundedCornerShape(8.dp)
                )
                .shadow(
                    elevation = 8.dp,
                    shape = RoundedCornerShape(8.dp)
                ),
            offset = DpOffset(0.dp, 4.dp),
            properties = PopupProperties(
                focusable = true,
                dismissOnBackPress = true,
                dismissOnClickOutside = true
            )
        ) {
            NavigationMenuContent(
                items = items,
                onItemClick = { item ->
                    onItemClick(item)
                    item.action?.invoke()
                    if (item.children.isEmpty()) {
                        internalExpanded = false
                        onExpandedChange(false)
                    }
                },
                level = 0
            )
        }
    }
}

@Composable
private fun NavigationMenuContent(
    items: List<NavigationMenuItem>,
    onItemClick: (NavigationMenuItem) -> Unit,
    level: Int
) {
    items.forEach { item ->
        NavigationMenuItemContent(
            item = item,
            onItemClick = onItemClick,
            level = level
        )
    }
}

@Composable
private fun NavigationMenuItemContent(
    item: NavigationMenuItem,
    onItemClick: (NavigationMenuItem) -> Unit,
    level: Int
) {
    var expanded by remember { mutableStateOf(false) }
    val hasChildren = item.children.isNotEmpty()

    Column {
        DropdownMenuItem(
            text = {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Level indentation
                    if (level > 0) {
                        Spacer(modifier = Modifier.width((level * 16).dp))
                    }

                    // Icon
                    item.icon?.let { icon ->
                        Icon(
                            imageVector = icon,
                            contentDescription = null,
                            tint = if (item.enabled) TchatColors.onSurface else TchatColors.disabled,
                            modifier = Modifier.size(20.dp)
                        )
                    }

                    // Label
                    Text(
                        text = item.label,
                        style = MaterialTheme.typography.bodyMedium,
                        color = if (item.enabled) TchatColors.onSurface else TchatColors.disabled,
                        modifier = Modifier.weight(1f)
                    )

                    // Submenu indicator
                    if (hasChildren) {
                        Icon(
                            imageVector = if (expanded) Icons.Default.ArrowDropDown else Icons.Default.ChevronRight,
                            contentDescription = if (expanded) "Collapse" else "Expand",
                            tint = if (item.enabled) TchatColors.onSurfaceVariant else TchatColors.disabled,
                            modifier = Modifier.size(16.dp)
                        )
                    }
                }
            },
            onClick = {
                if (item.enabled) {
                    if (hasChildren) {
                        expanded = !expanded
                    } else {
                        onItemClick(item)
                    }
                }
            },
            enabled = item.enabled,
            modifier = Modifier.fillMaxWidth()
        )

        // Submenu items
        AnimatedVisibility(
            visible = expanded && hasChildren,
            enter = fadeIn(),
            exit = fadeOut()
        ) {
            NavigationMenuContent(
                items = item.children,
                onItemClick = onItemClick,
                level = level + 1
            )
        }
    }
}

@Composable
fun TchatNavigationMenuTrigger(
    text: String,
    icon: ImageVector? = null,
    expanded: Boolean = false,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(4.dp),
        color = if (expanded) TchatColors.surfaceVariant else TchatColors.surface,
        border = androidx.compose.foundation.BorderStroke(
            1.dp,
            if (expanded) TchatColors.primary else TchatColors.outline
        )
    ) {
        Row(
            modifier = Modifier.padding(horizontal = 12.dp, vertical = 8.dp),
            horizontalArrangement = Arrangement.spacedBy(8.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            icon?.let {
                Icon(
                    imageVector = it,
                    contentDescription = null,
                    tint = TchatColors.onSurface,
                    modifier = Modifier.size(16.dp)
                )
            }

            Text(
                text = text,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )

            Icon(
                imageVector = Icons.Default.ArrowDropDown,
                contentDescription = if (expanded) "Collapse" else "Expand",
                tint = TchatColors.onSurfaceVariant,
                modifier = Modifier.size(16.dp)
            )
        }
    }
}