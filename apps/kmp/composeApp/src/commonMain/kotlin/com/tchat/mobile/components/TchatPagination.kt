package com.tchat.mobile.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import kotlin.math.max
import kotlin.math.min

/**
 * TchatPagination - Page navigation component
 *
 * Features:
 * - Page numbers with ellipsis for large ranges
 * - Previous/Next navigation with disabled states
 * - Jump to page input with validation
 * - Size configuration and responsive design
 */

enum class PaginationSize {
    SMALL,
    MEDIUM,
    LARGE
}

@Composable
fun TchatPagination(
    currentPage: Int,
    totalPages: Int,
    onPageChange: (Int) -> Unit,
    size: PaginationSize = PaginationSize.MEDIUM,
    showJumpToPage: Boolean = true,
    showFirstLast: Boolean = true,
    maxVisiblePages: Int = 7,
    modifier: Modifier = Modifier
) {
    val focusManager = LocalFocusManager.current
    var jumpPageText by remember { mutableStateOf("") }
    var showJumpInput by remember { mutableStateOf(false) }

    // Ensure valid page range
    val validCurrentPage = currentPage.coerceIn(1, totalPages)
    val validTotalPages = max(1, totalPages)

    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // Main pagination controls
        Row(
            horizontalArrangement = Arrangement.spacedBy(4.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // First page button
            if (showFirstLast && validTotalPages > maxVisiblePages) {
                PaginationButton(
                    icon = Icons.Default.FirstPage,
                    enabled = validCurrentPage > 1,
                    onClick = { onPageChange(1) },
                    size = size,
                    contentDescription = "First page"
                )
            }

            // Previous page button
            PaginationButton(
                icon = Icons.Default.ChevronLeft,
                enabled = validCurrentPage > 1,
                onClick = { onPageChange(validCurrentPage - 1) },
                size = size,
                contentDescription = "Previous page"
            )

            // Page numbers
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(2.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                val pageNumbers = generatePageNumbers(validCurrentPage, validTotalPages, maxVisiblePages)

                items(pageNumbers) { pageItem ->
                    when (pageItem.type) {
                        PageItemType.PAGE -> {
                            PaginationPageButton(
                                page = pageItem.value,
                                isSelected = pageItem.value == validCurrentPage,
                                onClick = { onPageChange(pageItem.value) },
                                size = size
                            )
                        }
                        PageItemType.ELLIPSIS -> {
                            PaginationEllipsis(size = size)
                        }
                    }
                }
            }

            // Next page button
            PaginationButton(
                icon = Icons.Default.ChevronRight,
                enabled = validCurrentPage < validTotalPages,
                onClick = { onPageChange(validCurrentPage + 1) },
                size = size,
                contentDescription = "Next page"
            )

            // Last page button
            if (showFirstLast && validTotalPages > maxVisiblePages) {
                PaginationButton(
                    icon = Icons.Default.LastPage,
                    enabled = validCurrentPage < validTotalPages,
                    onClick = { onPageChange(validTotalPages) },
                    size = size,
                    contentDescription = "Last page"
                )
            }
        }

        // Jump to page input
        if (showJumpToPage && validTotalPages > 10) {
            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                if (!showJumpInput) {
                    TextButton(
                        onClick = { showJumpInput = true },
                        modifier = Modifier.semantics {
                            contentDescription = "Jump to page"
                        }
                    ) {
                        Icon(
                            imageVector = Icons.Default.MoreHoriz,
                            contentDescription = null,
                            modifier = Modifier.size(16.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = "Jump to page",
                            style = MaterialTheme.typography.bodySmall
                        )
                    }
                } else {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(4.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            text = "Go to:",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )

                        OutlinedTextField(
                            value = jumpPageText,
                            onValueChange = { jumpPageText = it.filter { char -> char.isDigit() } },
                            modifier = Modifier.width(80.dp),
                            singleLine = true,
                            textStyle = MaterialTheme.typography.bodySmall.copy(textAlign = TextAlign.Center),
                            keyboardOptions = KeyboardOptions(
                                keyboardType = KeyboardType.Number,
                                imeAction = ImeAction.Go
                            ),
                            keyboardActions = KeyboardActions(
                                onGo = {
                                    val page = jumpPageText.toIntOrNull()
                                    if (page != null && page in 1..validTotalPages) {
                                        onPageChange(page)
                                        showJumpInput = false
                                        jumpPageText = ""
                                        focusManager.clearFocus()
                                    }
                                }
                            )
                        )

                        IconButton(
                            onClick = {
                                val page = jumpPageText.toIntOrNull()
                                if (page != null && page in 1..validTotalPages) {
                                    onPageChange(page)
                                    showJumpInput = false
                                    jumpPageText = ""
                                }
                            },
                            enabled = jumpPageText.toIntOrNull()?.let { it in 1..validTotalPages } == true,
                            modifier = Modifier.size(24.dp)
                        ) {
                            Icon(
                                imageVector = Icons.Default.Check,
                                contentDescription = "Go to page",
                                modifier = Modifier.size(16.dp)
                            )
                        }

                        IconButton(
                            onClick = {
                                showJumpInput = false
                                jumpPageText = ""
                            },
                            modifier = Modifier.size(24.dp)
                        ) {
                            Icon(
                                imageVector = Icons.Default.Close,
                                contentDescription = "Cancel",
                                modifier = Modifier.size(16.dp)
                            )
                        }
                    }
                }
            }
        }

        // Page info
        Text(
            text = "Page $validCurrentPage of $validTotalPages",
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun PaginationButton(
    icon: ImageVector,
    enabled: Boolean,
    onClick: () -> Unit,
    size: PaginationSize,
    contentDescription: String
) {
    val buttonSize = when (size) {
        PaginationSize.SMALL -> 32.dp
        PaginationSize.MEDIUM -> 40.dp
        PaginationSize.LARGE -> 48.dp
    }

    val iconSize = when (size) {
        PaginationSize.SMALL -> 16.dp
        PaginationSize.MEDIUM -> 20.dp
        PaginationSize.LARGE -> 24.dp
    }

    IconButton(
        onClick = onClick,
        enabled = enabled,
        modifier = Modifier
            .size(buttonSize)
            .semantics { this.contentDescription = contentDescription }
    ) {
        Icon(
            imageVector = icon,
            contentDescription = null,
            tint = if (enabled) TchatColors.onSurface else TchatColors.disabled,
            modifier = Modifier.size(iconSize)
        )
    }
}

@Composable
private fun PaginationPageButton(
    page: Int,
    isSelected: Boolean,
    onClick: () -> Unit,
    size: PaginationSize
) {
    val buttonSize = when (size) {
        PaginationSize.SMALL -> 32.dp
        PaginationSize.MEDIUM -> 40.dp
        PaginationSize.LARGE -> 48.dp
    }

    val textStyle = when (size) {
        PaginationSize.SMALL -> MaterialTheme.typography.bodySmall
        PaginationSize.MEDIUM -> MaterialTheme.typography.bodyMedium
        PaginationSize.LARGE -> MaterialTheme.typography.bodyLarge
    }

    Surface(
        onClick = onClick,
        modifier = Modifier
            .size(buttonSize)
            .semantics {
                contentDescription = if (isSelected) "Page $page, current page" else "Page $page"
                role = Role.Button
            },
        shape = CircleShape,
        color = if (isSelected) TchatColors.primary else TchatColors.surface,
        contentColor = if (isSelected) TchatColors.onPrimary else TchatColors.onSurface,
        border = if (!isSelected) androidx.compose.foundation.BorderStroke(1.dp, TchatColors.outline) else null
    ) {
        Box(
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = page.toString(),
                style = textStyle,
                fontWeight = if (isSelected) FontWeight.Medium else FontWeight.Normal
            )
        }
    }
}

@Composable
private fun PaginationEllipsis(size: PaginationSize) {
    val buttonSize = when (size) {
        PaginationSize.SMALL -> 32.dp
        PaginationSize.MEDIUM -> 40.dp
        PaginationSize.LARGE -> 48.dp
    }

    Box(
        modifier = Modifier.size(buttonSize),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "…",
            style = MaterialTheme.typography.bodyMedium,
            color = TchatColors.onSurfaceVariant
        )
    }
}

// Helper classes and functions
private enum class PageItemType {
    PAGE,
    ELLIPSIS
}

private data class PageItem(
    val type: PageItemType,
    val value: Int
)

private fun generatePageNumbers(
    currentPage: Int,
    totalPages: Int,
    maxVisible: Int
): List<PageItem> {
    if (totalPages <= maxVisible) {
        return (1..totalPages).map { PageItem(PageItemType.PAGE, it) }
    }

    val result = mutableListOf<PageItem>()
    val sidePages = (maxVisible - 3) / 2 // Reserve space for 1, ..., totalPages

    when {
        currentPage <= sidePages + 2 -> {
            // Current page is near the beginning
            (1..min(maxVisible - 2, totalPages - 1)).forEach {
                result.add(PageItem(PageItemType.PAGE, it))
            }
            result.add(PageItem(PageItemType.ELLIPSIS, -1))
            result.add(PageItem(PageItemType.PAGE, totalPages))
        }
        currentPage >= totalPages - sidePages - 1 -> {
            // Current page is near the end
            result.add(PageItem(PageItemType.PAGE, 1))
            result.add(PageItem(PageItemType.ELLIPSIS, -1))
            (max(totalPages - maxVisible + 3, 2)..totalPages).forEach {
                result.add(PageItem(PageItemType.PAGE, it))
            }
        }
        else -> {
            // Current page is in the middle
            result.add(PageItem(PageItemType.PAGE, 1))
            result.add(PageItem(PageItemType.ELLIPSIS, -1))
            ((currentPage - sidePages)..(currentPage + sidePages)).forEach {
                result.add(PageItem(PageItemType.PAGE, it))
            }
            result.add(PageItem(PageItemType.ELLIPSIS, -1))
            result.add(PageItem(PageItemType.PAGE, totalPages))
        }
    }

    return result
}