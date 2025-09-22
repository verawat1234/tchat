package com.tchat.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Pagination component following Tchat design system
 */
@Composable
fun TchatPagination(
    currentPage: Int,
    totalPages: Int,
    onPageChange: (Int) -> Unit,
    modifier: Modifier = Modifier,
    style: TchatPaginationStyle = TchatPaginationStyle.Numbered,
    size: TchatPaginationSize = TchatPaginationSize.Medium,
    showPageSize: Boolean = false,
    showInfo: Boolean = false,
    showJumpToPage: Boolean = false,
    maxVisiblePages: Int = 7,
    pageSize: Int = 20,
    totalItems: Int = 0,
    onPageSizeChange: ((Int) -> Unit)? = null
) {
    val hapticFeedback = LocalHapticFeedback.current
    var inputPage by remember { mutableStateOf("") }
    var showPageInput by remember { mutableStateOf(false) }

    val pageSizeOptions = listOf(10, 20, 50, 100)

    val visiblePages = remember(currentPage, totalPages, maxVisiblePages) {
        if (totalPages <= 0) return@remember emptyList<Int>()

        val halfVisible = maxVisiblePages / 2
        var start = maxOf(1, currentPage - halfVisible)
        var end = minOf(totalPages, start + maxVisiblePages - 1)

        // Adjust start if we're near the end
        if (end - start + 1 < maxVisiblePages) {
            start = maxOf(1, end - maxVisiblePages + 1)
        }

        (start..end).toList()
    }

    val startItem = (currentPage - 1) * pageSize + 1
    val endItem = minOf(currentPage * pageSize, totalItems)

    fun navigateToPage(page: Int) {
        if (page in 1..totalPages && page != currentPage) {
            onPageChange(page)
            hapticFeedback.performHapticFeedback(
                androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
            )
        }
    }

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.sm)
    ) {
        // Main pagination controls
        when (style) {
            TchatPaginationStyle.Numbered -> {
                NumberedPagination(
                    currentPage = currentPage,
                    totalPages = totalPages,
                    visiblePages = visiblePages,
                    size = size,
                    onPageChange = ::navigateToPage
                )
            }
            TchatPaginationStyle.Simple -> {
                SimplePagination(
                    currentPage = currentPage,
                    totalPages = totalPages,
                    size = size,
                    onPageChange = ::navigateToPage
                )
            }
            TchatPaginationStyle.Compact -> {
                CompactPagination(
                    currentPage = currentPage,
                    totalPages = totalPages,
                    size = size,
                    showPageInput = showPageInput,
                    onShowPageInputChange = { showPageInput = it },
                    inputPage = inputPage,
                    onInputPageChange = { inputPage = it },
                    onPageChange = ::navigateToPage
                )
            }
        }

        // Additional controls
        if (showPageSize || showInfo || showJumpToPage) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                if (showInfo) {
                    Text(
                        text = "Showing $startItem-$endItem of $totalItems",
                        fontSize = (size.fontSize.value - 1).sp,
                        color = Colors.textSecondary
                    )
                }

                Row(
                    horizontalArrangement = Arrangement.spacedBy(Spacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    if (showJumpToPage) {
                        JumpToPageControl(
                            size = size,
                            totalPages = totalPages,
                            onPageChange = ::navigateToPage
                        )
                    }

                    if (showPageSize && onPageSizeChange != null) {
                        PageSizeSelector(
                            currentPageSize = pageSize,
                            options = pageSizeOptions,
                            size = size,
                            onPageSizeChange = onPageSizeChange
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun NumberedPagination(
    currentPage: Int,
    totalPages: Int,
    visiblePages: List<Int>,
    size: TchatPaginationSize,
    onPageChange: (Int) -> Unit
) {
    Row(
        horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Previous button
        PaginationButton(
            icon = Icons.Default.ChevronLeft,
            isEnabled = currentPage > 1,
            size = size,
            onClick = { onPageChange(currentPage - 1) }
        )

        // First page (if not visible)
        if (visiblePages.isNotEmpty() && !visiblePages.contains(1) && totalPages > 1) {
            PaginationButton(
                text = "1",
                size = size,
                onClick = { onPageChange(1) }
            )

            if (visiblePages.first() > 2) {
                Text(
                    text = "...",
                    fontSize = size.fontSize,
                    color = Colors.textSecondary,
                    modifier = Modifier
                        .size(size.buttonSize)
                        .wrapContentSize(Alignment.Center)
                )
            }
        }

        // Visible page numbers
        visiblePages.forEach { page ->
            PaginationButton(
                text = page.toString(),
                isSelected = page == currentPage,
                size = size,
                onClick = { onPageChange(page) }
            )
        }

        // Last page (if not visible)
        if (visiblePages.isNotEmpty() && !visiblePages.contains(totalPages) && totalPages > 1) {
            if (visiblePages.last() < totalPages - 1) {
                Text(
                    text = "...",
                    fontSize = size.fontSize,
                    color = Colors.textSecondary,
                    modifier = Modifier
                        .size(size.buttonSize)
                        .wrapContentSize(Alignment.Center)
                )
            }

            PaginationButton(
                text = totalPages.toString(),
                size = size,
                onClick = { onPageChange(totalPages) }
            )
        }

        // Next button
        PaginationButton(
            icon = Icons.Default.ChevronRight,
            isEnabled = currentPage < totalPages,
            size = size,
            onClick = { onPageChange(currentPage + 1) }
        )
    }
}

@Composable
private fun SimplePagination(
    currentPage: Int,
    totalPages: Int,
    size: TchatPaginationSize,
    onPageChange: (Int) -> Unit
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Previous button
        OutlinedButton(
            onClick = { onPageChange(currentPage - 1) },
            enabled = currentPage > 1,
            modifier = Modifier.height(size.buttonSize)
        ) {
            Icon(
                imageVector = Icons.Default.ChevronLeft,
                contentDescription = null,
                modifier = Modifier.size(size.iconSize)
            )
            Spacer(modifier = Modifier.width(Spacing.xs))
            Text(
                text = "Previous",
                fontSize = size.fontSize,
                fontWeight = FontWeight.Medium
            )
        }

        // Page info
        Text(
            text = "$currentPage of $totalPages",
            fontSize = size.fontSize,
            color = Colors.textSecondary
        )

        // Next button
        OutlinedButton(
            onClick = { onPageChange(currentPage + 1) },
            enabled = currentPage < totalPages,
            modifier = Modifier.height(size.buttonSize)
        ) {
            Text(
                text = "Next",
                fontSize = size.fontSize,
                fontWeight = FontWeight.Medium
            )
            Spacer(modifier = Modifier.width(Spacing.xs))
            Icon(
                imageVector = Icons.Default.ChevronRight,
                contentDescription = null,
                modifier = Modifier.size(size.iconSize)
            )
        }
    }
}

@Composable
private fun CompactPagination(
    currentPage: Int,
    totalPages: Int,
    size: TchatPaginationSize,
    showPageInput: Boolean,
    onShowPageInputChange: (Boolean) -> Unit,
    inputPage: String,
    onInputPageChange: (String) -> Unit,
    onPageChange: (Int) -> Unit
) {
    Row(
        horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Previous button
        PaginationButton(
            icon = Icons.Default.ChevronLeft,
            isEnabled = currentPage > 1,
            size = size,
            onClick = { onPageChange(currentPage - 1) }
        )

        // Current page input/display
        if (showPageInput) {
            OutlinedTextField(
                value = inputPage,
                onValueChange = onInputPageChange,
                modifier = Modifier.width(60.dp),
                textStyle = LocalTextStyle.current.copy(
                    fontSize = size.fontSize,
                    textAlign = TextAlign.Center
                ),
                keyboardOptions = KeyboardOptions(
                    keyboardType = KeyboardType.Number,
                    imeAction = ImeAction.Done
                ),
                keyboardActions = KeyboardActions(
                    onDone = {
                        val page = inputPage.toIntOrNull()
                        if (page != null && page in 1..totalPages) {
                            onPageChange(page)
                        }
                        onShowPageInputChange(false)
                    }
                ),
                singleLine = true
            )
        } else {
            Box(
                modifier = Modifier
                    .width(60.dp)
                    .height(size.buttonSize)
                    .border(
                        width = 1.dp,
                        color = Colors.border,
                        shape = RoundedCornerShape(6.dp)
                    )
                    .clickable {
                        onShowPageInputChange(true)
                        onInputPageChange(currentPage.toString())
                    },
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = currentPage.toString(),
                    fontSize = size.fontSize,
                    fontWeight = FontWeight.Medium,
                    color = Colors.textPrimary
                )
            }
        }

        Text(
            text = "of $totalPages",
            fontSize = size.fontSize,
            color = Colors.textSecondary
        )

        // Next button
        PaginationButton(
            icon = Icons.Default.ChevronRight,
            isEnabled = currentPage < totalPages,
            size = size,
            onClick = { onPageChange(currentPage + 1) }
        )
    }
}

@Composable
private fun PaginationButton(
    text: String? = null,
    icon: androidx.compose.ui.graphics.vector.ImageVector? = null,
    isSelected: Boolean = false,
    isEnabled: Boolean = true,
    size: TchatPaginationSize,
    onClick: () -> Unit
) {
    Box(
        modifier = Modifier
            .size(size.buttonSize)
            .clip(RoundedCornerShape(6.dp))
            .background(
                if (isSelected) Colors.primary else Color.Transparent
            )
            .border(
                width = 1.dp,
                color = if (isSelected) Colors.primary else Colors.border,
                shape = RoundedCornerShape(6.dp)
            )
            .clickable(enabled = isEnabled) { onClick() },
        contentAlignment = Alignment.Center
    ) {
        when {
            text != null -> {
                Text(
                    text = text,
                    fontSize = size.fontSize,
                    fontWeight = FontWeight.Medium,
                    color = when {
                        isSelected -> Colors.textOnPrimary
                        isEnabled -> Colors.textPrimary
                        else -> Colors.textDisabled
                    }
                )
            }
            icon != null -> {
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = when {
                        isSelected -> Colors.textOnPrimary
                        isEnabled -> Colors.textPrimary
                        else -> Colors.textDisabled
                    },
                    modifier = Modifier.size(size.iconSize)
                )
            }
        }
    }
}

@Composable
private fun JumpToPageControl(
    size: TchatPaginationSize,
    totalPages: Int,
    onPageChange: (Int) -> Unit
) {
    var inputPage by remember { mutableStateOf("") }

    Row(
        horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            text = "Go to:",
            fontSize = (size.fontSize.value - 1).sp,
            color = Colors.textSecondary
        )

        OutlinedTextField(
            value = inputPage,
            onValueChange = { inputPage = it },
            modifier = Modifier.width(60.dp),
            textStyle = LocalTextStyle.current.copy(
                fontSize = size.fontSize,
                textAlign = TextAlign.Center
            ),
            keyboardOptions = KeyboardOptions(
                keyboardType = KeyboardType.Number,
                imeAction = ImeAction.Done
            ),
            keyboardActions = KeyboardActions(
                onDone = {
                    val page = inputPage.toIntOrNull()
                    if (page != null && page in 1..totalPages) {
                        onPageChange(page)
                        inputPage = ""
                    }
                }
            ),
            singleLine = true
        )
    }
}

@Composable
private fun PageSizeSelector(
    currentPageSize: Int,
    options: List<Int>,
    size: TchatPaginationSize,
    onPageSizeChange: (Int) -> Unit
) {
    var expanded by remember { mutableStateOf(false) }

    Row(
        horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            text = "Show:",
            fontSize = (size.fontSize.value - 1).sp,
            color = Colors.textSecondary
        )

        Box {
            Row(
                modifier = Modifier
                    .border(
                        width = 1.dp,
                        color = Colors.border,
                        shape = RoundedCornerShape(6.dp)
                    )
                    .clickable { expanded = true }
                    .padding(horizontal = Spacing.sm, vertical = Spacing.xs),
                horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = currentPageSize.toString(),
                    fontSize = size.fontSize,
                    color = Colors.textPrimary
                )

                Icon(
                    imageVector = Icons.Default.ExpandMore,
                    contentDescription = null,
                    tint = Colors.textSecondary,
                    modifier = Modifier.size((size.iconSize.value - 2).dp)
                )
            }

            DropdownMenu(
                expanded = expanded,
                onDismissRequest = { expanded = false }
            ) {
                options.forEach { option ->
                    DropdownMenuItem(
                        text = {
                            Text("$option per page")
                        },
                        onClick = {
                            onPageSizeChange(option)
                            expanded = false
                        }
                    )
                }
            }
        }
    }
}

/**
 * Pagination style definitions
 */
enum class TchatPaginationStyle {
    Numbered,
    Simple,
    Compact
}

/**
 * Pagination size definitions
 */
enum class TchatPaginationSize(
    val buttonSize: androidx.compose.ui.unit.Dp,
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val iconSize: androidx.compose.ui.unit.Dp
) {
    Small(
        buttonSize = 32.dp,
        fontSize = 12.sp,
        iconSize = 12.dp
    ),
    Medium(
        buttonSize = 40.dp,
        fontSize = 14.sp,
        iconSize = 14.dp
    ),
    Large(
        buttonSize = 48.dp,
        fontSize = 16.sp,
        iconSize = 16.dp
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatPaginationPreview() {
    var currentPage1 by remember { mutableStateOf(5) }
    var currentPage2 by remember { mutableStateOf(2) }
    var currentPage3 by remember { mutableStateOf(3) }
    var pageSize by remember { mutableStateOf(20) }

    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Numbered pagination
        TchatPagination(
            currentPage = currentPage1,
            totalPages = 20,
            onPageChange = { currentPage1 = it },
            style = TchatPaginationStyle.Numbered,
            size = TchatPaginationSize.Medium,
            showPageSize = true,
            showInfo = true,
            showJumpToPage = true,
            totalItems = 1000,
            pageSize = pageSize,
            onPageSizeChange = { pageSize = it }
        )

        Divider()

        // Simple pagination
        TchatPagination(
            currentPage = currentPage2,
            totalPages = 10,
            onPageChange = { currentPage2 = it },
            style = TchatPaginationStyle.Simple,
            size = TchatPaginationSize.Large
        )

        Divider()

        // Compact pagination
        TchatPagination(
            currentPage = currentPage3,
            totalPages = 15,
            onPageChange = { currentPage3 = it },
            style = TchatPaginationStyle.Compact,
            size = TchatPaginationSize.Small
        )

        Divider()

        // Small numbered pagination
        TchatPagination(
            currentPage = 1,
            totalPages = 5,
            onPageChange = { },
            style = TchatPaginationStyle.Numbered,
            size = TchatPaginationSize.Small,
            maxVisiblePages = 5
        )
    }
}