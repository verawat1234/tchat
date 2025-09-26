package com.tchat.mobile.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.*
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * TchatMenuBar - Application menu bar component
 *
 * Features:
 * - Horizontal menu with dropdown submenus
 * - Keyboard shortcuts integration and support
 * - Accessibility with ARIA labels and navigation
 * - Platform-native menu styling and behavior
 */

data class MenuBarItem(
    val id: String,
    val label: String,
    val icon: ImageVector? = null,
    val enabled: Boolean = true,
    val shortcut: String? = null,
    val children: List<NavigationMenuItem> = emptyList(),
    val action: (() -> Unit)? = null
)

@Composable
fun TchatMenuBar(
    items: List<MenuBarItem>,
    onItemClick: (MenuBarItem) -> Unit = {},
    onMenuItemClick: (NavigationMenuItem) -> Unit = {},
    modifier: Modifier = Modifier
) {
    var expandedMenuId by remember { mutableStateOf<String?>(null) }

    Surface(
        modifier = modifier.fillMaxWidth(),
        color = TchatColors.surface,
        shadowElevation = 2.dp
    ) {
        LazyRow(
            modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
            horizontalArrangement = Arrangement.spacedBy(4.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            items(items) { item ->
                MenuBarItemContent(
                    item = item,
                    expanded = expandedMenuId == item.id,
                    onExpandedChange = { expanded ->
                        expandedMenuId = if (expanded) item.id else null
                    },
                    onItemClick = onItemClick,
                    onMenuItemClick = onMenuItemClick
                )
            }
        }
    }
}

@Composable
private fun MenuBarItemContent(
    item: MenuBarItem,
    expanded: Boolean,
    onExpandedChange: (Boolean) -> Unit,
    onItemClick: (MenuBarItem) -> Unit,
    onMenuItemClick: (NavigationMenuItem) -> Unit
) {
    val hasChildren = item.children.isNotEmpty()

    Box {
        Surface(
            modifier = Modifier
                .clip(RoundedCornerShape(4.dp))
                .clickable(enabled = item.enabled) {
                    if (hasChildren) {
                        onExpandedChange(!expanded)
                    } else {
                        onItemClick(item)
                        item.action?.invoke()
                    }
                }
                .semantics {
                    role = Role.Button
                    if (item.shortcut != null) {
                        contentDescription = "${item.label}, shortcut ${item.shortcut}"
                    }
                    if (hasChildren) {
                        stateDescription = if (expanded) "expanded" else "collapsed"
                    }
                },
            color = when {
                expanded -> TchatColors.primary.copy(alpha = 0.1f)
                else -> TchatColors.surface
            },
            contentColor = if (item.enabled) TchatColors.onSurface else TchatColors.disabled
        ) {
            Row(
                modifier = Modifier.padding(horizontal = 12.dp, vertical = 8.dp),
                horizontalArrangement = Arrangement.spacedBy(6.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Icon
                item.icon?.let { icon ->
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                }

                // Label
                Text(
                    text = item.label,
                    style = MaterialTheme.typography.bodyMedium
                )

                // Keyboard shortcut indicator
                item.shortcut?.let { shortcut ->
                    Text(
                        text = shortcut,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant,
                        modifier = Modifier
                            .background(
                                color = TchatColors.surfaceVariant,
                                shape = RoundedCornerShape(2.dp)
                            )
                            .padding(horizontal = 4.dp, vertical = 1.dp)
                    )
                }
            }
        }

        // Dropdown menu for items with children
        if (hasChildren) {
            DropdownMenu(
                expanded = expanded,
                onDismissRequest = { onExpandedChange(false) },
                modifier = Modifier
                    .widthIn(min = 200.dp, max = 300.dp)
                    .background(
                        color = TchatColors.surface,
                        shape = RoundedCornerShape(8.dp)
                    ),
                offset = androidx.compose.ui.unit.DpOffset(0.dp, 4.dp)
            ) {
                item.children.forEach { menuItem ->
                    MenuBarDropdownItem(
                        item = menuItem,
                        onClick = {
                            onMenuItemClick(menuItem)
                            menuItem.action?.invoke()
                            onExpandedChange(false)
                        }
                    )
                }
            }
        }
    }
}

@Composable
private fun MenuBarDropdownItem(
    item: NavigationMenuItem,
    onClick: () -> Unit
) {
    DropdownMenuItem(
        text = {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Icon
                item.icon?.let { icon ->
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        tint = if (item.enabled) TchatColors.onSurface else TchatColors.disabled,
                        modifier = Modifier.size(16.dp)
                    )
                }

                // Label
                Text(
                    text = item.label,
                    style = MaterialTheme.typography.bodyMedium,
                    color = if (item.enabled) TchatColors.onSurface else TchatColors.disabled,
                    modifier = Modifier.weight(1f)
                )

                // Submenu indicator (if has children)
                if (item.children.isNotEmpty()) {
                    Icon(
                        imageVector = Icons.Default.ChevronRight,
                        contentDescription = "Has submenu",
                        tint = if (item.enabled) TchatColors.onSurfaceVariant else TchatColors.disabled,
                        modifier = Modifier.size(14.dp)
                    )
                }
            }
        },
        onClick = onClick,
        enabled = item.enabled,
        modifier = Modifier.semantics {
            role = Role.Button
            if (!item.enabled) {
                disabled()
            }
        }
    )
}

// Helper function to create common menu bar items
object MenuBarDefaults {
    fun createFileMenu(): MenuBarItem {
        return MenuBarItem(
            id = "file",
            label = "File",
            children = listOf(
                NavigationMenuItem("new", "New", Icons.Default.Add, action = {}),
                NavigationMenuItem("open", "Open", Icons.Default.Folder, action = {}),
                NavigationMenuItem("save", "Save", Icons.Default.Save, action = {}),
                NavigationMenuItem("divider1", "", enabled = false),
                NavigationMenuItem("exit", "Exit", Icons.Default.ExitToApp, action = {})
            )
        )
    }

    fun createEditMenu(): MenuBarItem {
        return MenuBarItem(
            id = "edit",
            label = "Edit",
            children = listOf(
                NavigationMenuItem("undo", "Undo", Icons.Default.Undo, action = {}),
                NavigationMenuItem("redo", "Redo", Icons.Default.Redo, action = {}),
                NavigationMenuItem("divider1", "", enabled = false),
                NavigationMenuItem("cut", "Cut", Icons.Default.ContentCut, action = {}),
                NavigationMenuItem("copy", "Copy", Icons.Default.ContentCopy, action = {}),
                NavigationMenuItem("paste", "Paste", Icons.Default.ContentPaste, action = {})
            )
        )
    }

    fun createViewMenu(): MenuBarItem {
        return MenuBarItem(
            id = "view",
            label = "View",
            children = listOf(
                NavigationMenuItem("zoom_in", "Zoom In", Icons.Default.ZoomIn, action = {}),
                NavigationMenuItem("zoom_out", "Zoom Out", Icons.Default.ZoomOut, action = {}),
                NavigationMenuItem("full_screen", "Full Screen", Icons.Default.Fullscreen, action = {})
            )
        )
    }

    fun createHelpMenu(): MenuBarItem {
        return MenuBarItem(
            id = "help",
            label = "Help",
            children = listOf(
                NavigationMenuItem("about", "About", Icons.Default.Info, action = {}),
                NavigationMenuItem("documentation", "Documentation", Icons.Default.MenuBook, action = {}),
                NavigationMenuItem("support", "Support", Icons.Default.Support, action = {})
            )
        )
    }
}