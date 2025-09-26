package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.expandVertically
import androidx.compose.animation.shrinkVertically
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.selection.selectable
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * TchatTree - Hierarchical tree view component
 *
 * Features:
 * - Expandable/collapsible nodes with smooth animations
 * - Custom node rendering with icons and actions
 * - Selection modes (single/multiple) with checkboxes
 * - Drag and drop support for reordering (future feature)
 */

data class TreeNode<T>(
    val id: String,
    val data: T,
    val children: List<TreeNode<T>> = emptyList(),
    val icon: ImageVector? = null,
    val expandedByDefault: Boolean = false,
    val selectable: Boolean = true,
    val enabled: Boolean = true
)

enum class TreeSelectionMode {
    NONE,
    SINGLE,
    MULTIPLE
}

data class TreeState(
    val expandedNodes: Set<String> = emptySet(),
    val selectedNodes: Set<String> = emptySet()
)

@Composable
fun <T> TchatTree(
    nodes: List<TreeNode<T>>,
    modifier: Modifier = Modifier,
    state: TreeState = TreeState(),
    onStateChange: (TreeState) -> Unit = {},
    selectionMode: TreeSelectionMode = TreeSelectionMode.NONE,
    nodeContent: @Composable (TreeNode<T>, Int, Boolean, Boolean) -> Unit = { node, level, isExpanded, isSelected ->
        DefaultTreeNodeContent(
            node = node,
            level = level,
            isExpanded = isExpanded,
            isSelected = isSelected
        )
    },
    onNodeClick: (TreeNode<T>) -> Unit = {},
    onNodeDoubleClick: (TreeNode<T>) -> Unit = {},
    showConnectors: Boolean = true,
    maxDepth: Int = Int.MAX_VALUE
) {
    LazyColumn(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(2.dp)
    ) {
        items(nodes) { node ->
            TreeNodeItem(
                node = node,
                level = 0,
                state = state,
                onStateChange = onStateChange,
                selectionMode = selectionMode,
                nodeContent = nodeContent,
                onNodeClick = onNodeClick,
                onNodeDoubleClick = onNodeDoubleClick,
                showConnectors = showConnectors,
                maxDepth = maxDepth
            )
        }
    }
}

@Composable
private fun <T> TreeNodeItem(
    node: TreeNode<T>,
    level: Int,
    state: TreeState,
    onStateChange: (TreeState) -> Unit,
    selectionMode: TreeSelectionMode,
    nodeContent: @Composable (TreeNode<T>, Int, Boolean, Boolean) -> Unit,
    onNodeClick: (TreeNode<T>) -> Unit,
    onNodeDoubleClick: (TreeNode<T>) -> Unit,
    showConnectors: Boolean,
    maxDepth: Int
) {
    val isExpanded = node.id in state.expandedNodes
    val isSelected = node.id in state.selectedNodes
    val hasChildren = node.children.isNotEmpty()
    val canExpand = hasChildren && level < maxDepth

    Column {
        // Node content
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(4.dp))
                .background(
                    if (isSelected) TchatColors.primary.copy(alpha = 0.1f)
                    else Color.Transparent
                )
                .run {
                    when (selectionMode) {
                        TreeSelectionMode.SINGLE -> {
                            selectable(
                                selected = isSelected,
                                enabled = node.enabled && node.selectable,
                                onClick = {
                                    val newSelectedNodes = if (isSelected) {
                                        emptySet()
                                    } else {
                                        setOf(node.id)
                                    }
                                    onStateChange(state.copy(selectedNodes = newSelectedNodes))
                                    onNodeClick(node)
                                }
                            )
                        }
                        TreeSelectionMode.MULTIPLE -> {
                            clickable(enabled = node.enabled) {
                                if (node.selectable) {
                                    val newSelectedNodes = if (isSelected) {
                                        state.selectedNodes - node.id
                                    } else {
                                        state.selectedNodes + node.id
                                    }
                                    onStateChange(state.copy(selectedNodes = newSelectedNodes))
                                }
                                onNodeClick(node)
                            }
                        }
                        TreeSelectionMode.NONE -> {
                            clickable(enabled = node.enabled) {
                                onNodeClick(node)
                            }
                        }
                    }
                }
                .padding(vertical = 4.dp, horizontal = 8.dp)
                .semantics {
                    contentDescription = "Tree node ${node.id} at level ${level + 1}"
                    if (hasChildren) {
                        stateDescription = if (isExpanded) "expanded" else "collapsed"
                    }
                    if (isSelected) {
                        selected = true
                    }
                    if (!node.enabled) {
                        disabled()
                    }
                },
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            // Indentation and connectors
            if (level > 0) {
                Spacer(modifier = Modifier.width((level * 24).dp))
            }

            // Expansion toggle
            if (canExpand) {
                IconButton(
                    onClick = {
                        val newExpandedNodes = if (isExpanded) {
                            state.expandedNodes - node.id
                        } else {
                            state.expandedNodes + node.id
                        }
                        onStateChange(state.copy(expandedNodes = newExpandedNodes))
                    },
                    modifier = Modifier
                        .size(24.dp)
                        .semantics {
                            contentDescription = if (isExpanded) "Collapse" else "Expand"
                        }
                ) {
                    Icon(
                        imageVector = Icons.Default.ChevronRight,
                        contentDescription = null,
                        tint = if (node.enabled) TchatColors.onSurface else TchatColors.disabled,
                        modifier = Modifier
                            .size(16.dp)
                            .rotate(if (isExpanded) 90f else 0f)
                    )
                }
            } else if (hasChildren) {
                // Spacer for alignment when node can't expand due to max depth
                Spacer(modifier = Modifier.size(24.dp))
            } else {
                // Spacer for leaf nodes
                Spacer(modifier = Modifier.size(24.dp))
            }

            // Selection checkbox (for multiple selection mode)
            if (selectionMode == TreeSelectionMode.MULTIPLE && node.selectable) {
                Checkbox(
                    checked = isSelected,
                    enabled = node.enabled,
                    onCheckedChange = { checked ->
                        val newSelectedNodes = if (checked) {
                            state.selectedNodes + node.id
                        } else {
                            state.selectedNodes - node.id
                        }
                        onStateChange(state.copy(selectedNodes = newSelectedNodes))
                    },
                    modifier = Modifier.semantics {
                        contentDescription = if (isSelected) "Deselect node" else "Select node"
                    }
                )
            }

            // Node content
            Box(modifier = Modifier.weight(1f)) {
                nodeContent(node, level, isExpanded, isSelected)
            }
        }

        // Child nodes
        AnimatedVisibility(
            visible = isExpanded && canExpand,
            enter = expandVertically(),
            exit = shrinkVertically()
        ) {
            Column {
                node.children.forEach { childNode ->
                    TreeNodeItem(
                        node = childNode,
                        level = level + 1,
                        state = state,
                        onStateChange = onStateChange,
                        selectionMode = selectionMode,
                        nodeContent = nodeContent,
                        onNodeClick = onNodeClick,
                        onNodeDoubleClick = onNodeDoubleClick,
                        showConnectors = showConnectors,
                        maxDepth = maxDepth
                    )
                }
            }
        }
    }
}

@Composable
private fun <T> DefaultTreeNodeContent(
    node: TreeNode<T>,
    level: Int,
    isExpanded: Boolean,
    isSelected: Boolean
) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // Node icon
        node.icon?.let { icon ->
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = if (node.enabled) {
                    if (isSelected) TchatColors.primary else TchatColors.onSurface
                } else {
                    TchatColors.disabled
                },
                modifier = Modifier.size(20.dp)
            )
        }

        // Node text/data
        Text(
            text = node.data.toString(),
            style = MaterialTheme.typography.bodyMedium,
            color = if (node.enabled) {
                if (isSelected) TchatColors.primary else TchatColors.onSurface
            } else {
                TchatColors.disabled
            },
            fontWeight = if (isSelected) FontWeight.Medium else FontWeight.Normal,
            modifier = Modifier.weight(1f)
        )

        // Children count indicator
        if (node.children.isNotEmpty()) {
            Surface(
                shape = RoundedCornerShape(12.dp),
                color = if (isSelected) TchatColors.primary else TchatColors.surfaceVariant
            ) {
                Text(
                    text = node.children.size.toString(),
                    style = MaterialTheme.typography.bodySmall,
                    color = if (isSelected) TchatColors.onPrimary else TchatColors.onSurfaceVariant,
                    modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                )
            }
        }
    }
}

// Helper functions for tree operations
fun <T> TreeNode<T>.findById(id: String): TreeNode<T>? {
    if (this.id == id) return this
    for (child in children) {
        child.findById(id)?.let { return it }
    }
    return null
}

fun <T> TreeNode<T>.getAllNodeIds(): Set<String> {
    val ids = mutableSetOf(id)
    children.forEach { child ->
        ids.addAll(child.getAllNodeIds())
    }
    return ids
}

fun <T> TreeNode<T>.getDepth(): Int {
    return if (children.isEmpty()) {
        0
    } else {
        1 + children.maxOfOrNull { it.getDepth() }!!
    }
}

fun <T> List<TreeNode<T>>.findNodeById(id: String): TreeNode<T>? {
    for (node in this) {
        node.findById(id)?.let { return it }
    }
    return null
}

fun <T> List<TreeNode<T>>.getAllNodeIds(): Set<String> {
    val ids = mutableSetOf<String>()
    forEach { node ->
        ids.addAll(node.getAllNodeIds())
    }
    return ids
}

// Sample tree data for demonstration
object TreeDefaults {
    fun createFileSystemTree(): List<TreeNode<String>> {
        return listOf(
            TreeNode(
                id = "root",
                data = "Project Root",
                icon = Icons.Default.Folder,
                children = listOf(
                    TreeNode(
                        id = "src",
                        data = "src",
                        icon = Icons.Default.Folder,
                        children = listOf(
                            TreeNode(
                                id = "main.kt",
                                data = "main.kt",
                                icon = Icons.Default.InsertDriveFile
                            ),
                            TreeNode(
                                id = "utils.kt",
                                data = "utils.kt",
                                icon = Icons.Default.InsertDriveFile
                            )
                        )
                    ),
                    TreeNode(
                        id = "docs",
                        data = "docs",
                        icon = Icons.Default.Folder,
                        children = listOf(
                            TreeNode(
                                id = "readme.md",
                                data = "README.md",
                                icon = Icons.Default.Description
                            )
                        )
                    ),
                    TreeNode(
                        id = "build.gradle",
                        data = "build.gradle",
                        icon = Icons.Default.InsertDriveFile
                    )
                )
            )
        )
    }

    fun createOrganizationTree(): List<TreeNode<String>> {
        return listOf(
            TreeNode(
                id = "ceo",
                data = "CEO",
                icon = Icons.Default.Person,
                children = listOf(
                    TreeNode(
                        id = "cto",
                        data = "CTO",
                        icon = Icons.Default.Person,
                        children = listOf(
                            TreeNode(
                                id = "dev_team",
                                data = "Development Team",
                                icon = Icons.Default.Group,
                                children = listOf(
                                    TreeNode(id = "dev1", data = "Senior Developer", icon = Icons.Default.Person),
                                    TreeNode(id = "dev2", data = "Junior Developer", icon = Icons.Default.Person)
                                )
                            )
                        )
                    ),
                    TreeNode(
                        id = "cfo",
                        data = "CFO",
                        icon = Icons.Default.Person,
                        children = listOf(
                            TreeNode(id = "accountant", data = "Accountant", icon = Icons.Default.Person)
                        )
                    )
                )
            )
        )
    }
}