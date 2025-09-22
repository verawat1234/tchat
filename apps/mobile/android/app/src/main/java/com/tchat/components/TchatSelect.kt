package com.tchat.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.expandVertically
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.shrinkVertically
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Popup
import androidx.compose.ui.window.PopupProperties
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Advanced select/dropdown component following Tchat design system
 */
@Composable
fun <T> TchatSelect(
    selection: Set<T>,
    onSelectionChange: (Set<T>) -> Unit,
    options: List<T>,
    optionLabel: (T) -> String,
    modifier: Modifier = Modifier,
    placeholder: String = "Select option",
    mode: TchatSelectMode = TchatSelectMode.Single,
    size: TchatSelectSize = TchatSelectSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    searchable: Boolean = false,
    maxSelections: Int? = null,
    leadingIcon: ImageVector? = null
) {
    var expanded by remember { mutableStateOf(false) }
    var searchText by remember { mutableStateOf("") }
    val focusRequester = remember { FocusRequester() }
    val hapticFeedback = LocalHapticFeedback.current

    // Filter options based on search
    val filteredOptions = remember(options, searchText) {
        if (searchText.isEmpty()) {
            options
        } else {
            options.filter { optionLabel(it).contains(searchText, ignoreCase = true) }
        }
    }

    // Display text
    val displayText = remember(selection, mode) {
        when (mode) {
            TchatSelectMode.Single -> {
                selection.firstOrNull()?.let(optionLabel) ?: placeholder
            }
            TchatSelectMode.Multiple -> {
                when (selection.size) {
                    0 -> placeholder
                    1 -> optionLabel(selection.first())
                    else -> "${selection.size} selected"
                }
            }
        }
    }

    val isPlaceholderShown = selection.isEmpty()

    // Animation values
    val borderColor by animateFloatAsState(
        targetValue = when {
            !enabled -> 0.5f
            validationState is TchatValidationState.Invalid -> 1f
            validationState is TchatValidationState.Valid -> 1f
            expanded -> 1f
            else -> 1f
        },
        label = "border_alpha"
    )

    val chevronRotation by animateFloatAsState(
        targetValue = if (expanded) 180f else 0f,
        label = "chevron_rotation"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.xs)
    ) {
        // Select Button
        Box {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clip(RoundedCornerShape(Spacing.sm))
                    .border(
                        BorderStroke(
                            width = when {
                                validationState !is TchatValidationState.None -> 2.dp
                                expanded -> 2.dp
                                else -> 1.dp
                            },
                            color = when {
                                !enabled -> Colors.border.copy(alpha = 0.5f)
                                validationState is TchatValidationState.Invalid -> Colors.borderError
                                validationState is TchatValidationState.Valid -> Colors.success
                                expanded -> Colors.borderFocus
                                else -> Colors.border
                            }
                        ),
                        RoundedCornerShape(Spacing.sm)
                    )
                    .background(
                        if (enabled) Colors.background else Colors.surface.copy(alpha = 0.5f),
                        RoundedCornerShape(Spacing.sm)
                    )
                    .clickable(enabled = enabled) {
                        expanded = !expanded
                        hapticFeedback.performHapticFeedback(
                            androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                        )
                    }
                    .padding(
                        horizontal = size.horizontalPadding,
                        vertical = size.verticalPadding
                    ),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(Spacing.xs)
            ) {
                // Leading icon
                leadingIcon?.let { icon ->
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        tint = if (enabled) {
                            if (expanded) Colors.primary else Colors.textSecondary
                        } else Colors.textDisabled,
                        modifier = Modifier.size(16.dp)
                    )
                }

                // Display text
                Text(
                    text = displayText,
                    fontSize = size.fontSize,
                    color = when {
                        !enabled -> Colors.textDisabled
                        isPlaceholderShown -> Colors.textTertiary
                        else -> Colors.textPrimary
                    },
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis,
                    modifier = Modifier.weight(1f)
                )

                // Selection count badge for multiple mode
                if (mode == TchatSelectMode.Multiple && selection.isNotEmpty()) {
                    Box(
                        modifier = Modifier
                            .background(Colors.primary, RoundedCornerShape(10.dp))
                            .padding(horizontal = Spacing.xs, vertical = 2.dp)
                    ) {
                        Text(
                            text = "${selection.size}",
                            fontSize = 12.sp,
                            color = Colors.textOnPrimary,
                            fontWeight = FontWeight.Medium
                        )
                    }
                }

                // Chevron icon
                Icon(
                    imageVector = Icons.Default.KeyboardArrowDown,
                    contentDescription = if (expanded) "Collapse" else "Expand",
                    tint = if (enabled) {
                        if (expanded) Colors.primary else Colors.textSecondary
                    } else Colors.textDisabled,
                    modifier = Modifier
                        .size(16.dp)
                        .rotate(chevronRotation)
                )
            }

            // Dropdown
            if (expanded) {
                Popup(
                    alignment = Alignment.TopStart,
                    properties = PopupProperties(focusable = true),
                    onDismissRequest = { expanded = false }
                ) {
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .shadow(8.dp, RoundedCornerShape(Spacing.sm)),
                        shape = RoundedCornerShape(Spacing.sm),
                        colors = CardDefaults.cardColors(containerColor = Colors.cardBackground),
                        border = BorderStroke(1.dp, Colors.border)
                    ) {
                        Column {
                            // Search field
                            if (searchable) {
                                Row(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .background(Colors.surface, RoundedCornerShape(Spacing.xs))
                                        .padding(Spacing.sm),
                                    verticalAlignment = Alignment.CenterVertically,
                                    horizontalArrangement = Arrangement.spacedBy(Spacing.xs)
                                ) {
                                    Icon(
                                        imageVector = Icons.Default.Search,
                                        contentDescription = "Search",
                                        tint = Colors.textTertiary,
                                        modifier = Modifier.size(14.dp)
                                    )

                                    BasicTextField(
                                        value = searchText,
                                        onValueChange = { searchText = it },
                                        modifier = Modifier
                                            .fillMaxWidth()
                                            .focusRequester(focusRequester),
                                        textStyle = androidx.compose.ui.text.TextStyle(
                                            fontSize = 14.sp,
                                            color = Colors.textPrimary
                                        ),
                                        decorationBox = { innerTextField ->
                                            if (searchText.isEmpty()) {
                                                Text(
                                                    text = "Search options",
                                                    fontSize = 14.sp,
                                                    color = Colors.textTertiary
                                                )
                                            }
                                            innerTextField()
                                        }
                                    )
                                }

                                Divider(
                                    modifier = Modifier.padding(horizontal = Spacing.sm),
                                    color = Colors.border
                                )
                            }

                            // Options list
                            LazyColumn(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .heightIn(max = 200.dp)
                            ) {
                                items(filteredOptions) { option ->
                                    val isSelected = selection.contains(option)

                                    Row(
                                        modifier = Modifier
                                            .fillMaxWidth()
                                            .clickable {
                                                handleSelection(
                                                    option = option,
                                                    selection = selection,
                                                    onSelectionChange = onSelectionChange,
                                                    mode = mode,
                                                    maxSelections = maxSelections
                                                )
                                                hapticFeedback.performHapticFeedback(
                                                    androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                                )
                                                if (mode == TchatSelectMode.Single) {
                                                    expanded = false
                                                }
                                            }
                                            .background(
                                                if (isSelected) Colors.primary.copy(alpha = 0.1f)
                                                else Color.Transparent
                                            )
                                            .padding(
                                                horizontal = Spacing.sm,
                                                vertical = Spacing.sm
                                            ),
                                        verticalAlignment = Alignment.CenterVertically,
                                        horizontalArrangement = Arrangement.SpaceBetween
                                    ) {
                                        Text(
                                            text = optionLabel(option),
                                            fontSize = 14.sp,
                                            color = Colors.textPrimary,
                                            modifier = Modifier.weight(1f)
                                        )

                                        if (isSelected) {
                                            Icon(
                                                imageVector = when (mode) {
                                                    TchatSelectMode.Single -> Icons.Default.Check
                                                    TchatSelectMode.Multiple -> Icons.Default.CheckBox
                                                },
                                                contentDescription = "Selected",
                                                tint = Colors.primary,
                                                modifier = Modifier.size(16.dp)
                                            )
                                        }
                                    }
                                }
                            }
                        }
                    }
                }

                // Auto-focus search when opened
                LaunchedEffect(expanded) {
                    if (expanded && searchable) {
                        focusRequester.requestFocus()
                    }
                }
            }
        }

        // Validation message
        if (validationState is TchatValidationState.Invalid) {
            Text(
                text = validationState.message,
                fontSize = 12.sp,
                color = Colors.error,
                modifier = Modifier.padding(horizontal = Spacing.xs)
            )
        }
    }
}

/**
 * Single selection convenience composable
 */
@Composable
fun <T> TchatSelect(
    selection: T?,
    onSelectionChange: (T?) -> Unit,
    options: List<T>,
    optionLabel: (T) -> String,
    modifier: Modifier = Modifier,
    placeholder: String = "Select option",
    size: TchatSelectSize = TchatSelectSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    searchable: Boolean = false,
    leadingIcon: ImageVector? = null
) {
    TchatSelect(
        selection = selection?.let { setOf(it) } ?: emptySet(),
        onSelectionChange = { newSelection ->
            onSelectionChange(newSelection.firstOrNull())
        },
        options = options,
        optionLabel = optionLabel,
        modifier = modifier,
        placeholder = placeholder,
        mode = TchatSelectMode.Single,
        size = size,
        validationState = validationState,
        enabled = enabled,
        searchable = searchable,
        leadingIcon = leadingIcon
    )
}

// Helper function for selection handling
private fun <T> handleSelection(
    option: T,
    selection: Set<T>,
    onSelectionChange: (Set<T>) -> Unit,
    mode: TchatSelectMode,
    maxSelections: Int?
) {
    when (mode) {
        TchatSelectMode.Single -> {
            onSelectionChange(setOf(option))
        }
        TchatSelectMode.Multiple -> {
            if (selection.contains(option)) {
                onSelectionChange(selection - option)
            } else {
                if (maxSelections != null && selection.size >= maxSelections) {
                    return // Max selections reached
                }
                onSelectionChange(selection + option)
            }
        }
    }
}

/**
 * Select mode definitions
 */
enum class TchatSelectMode {
    Single,
    Multiple
}

/**
 * Select size definitions
 */
enum class TchatSelectSize(
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
) {
    Small(
        fontSize = 14.sp,
        horizontalPadding = Spacing.sm,
        verticalPadding = Spacing.xs
    ),
    Medium(
        fontSize = 16.sp,
        horizontalPadding = Spacing.md,
        verticalPadding = Spacing.sm
    ),
    Large(
        fontSize = 18.sp,
        horizontalPadding = Spacing.lg,
        verticalPadding = Spacing.md
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatSelectPreview() {
    val languages = listOf("Swift", "Kotlin", "JavaScript", "Python", "Go", "Rust")

    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        // Single selection
        TchatSelect(
            selection = "Swift",
            onSelectionChange = { },
            options = languages,
            optionLabel = { it },
            placeholder = "Choose a language",
            leadingIcon = Icons.Default.Code
        )

        // Multiple selection with search
        TchatSelect(
            selection = setOf("Swift", "Kotlin"),
            onSelectionChange = { },
            options = languages,
            optionLabel = { it },
            placeholder = "Select languages",
            mode = TchatSelectMode.Multiple,
            searchable = true,
            leadingIcon = Icons.Default.Language
        )

        // Validation state
        TchatSelect(
            selection = null as String?,
            onSelectionChange = { },
            options = listOf("Valid", "Invalid", "Neutral"),
            optionLabel = { it },
            placeholder = "Required field",
            validationState = TchatValidationState.Invalid("This field is required"),
            leadingIcon = Icons.Default.Warning
        )

        // Disabled state
        TchatSelect(
            selection = "Disabled",
            onSelectionChange = { },
            options = listOf("Disabled", "Option"),
            optionLabel = { it },
            placeholder = "Disabled select",
            enabled = false
        )

        // Different sizes
        Column(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
            TchatSelect(
                selection = null as String?,
                onSelectionChange = { },
                options = listOf("Small"),
                optionLabel = { it },
                placeholder = "Small size",
                size = TchatSelectSize.Small
            )

            TchatSelect(
                selection = null as String?,
                onSelectionChange = { },
                options = listOf("Large"),
                optionLabel = { it },
                placeholder = "Large size",
                size = TchatSelectSize.Large
            )
        }
    }
}