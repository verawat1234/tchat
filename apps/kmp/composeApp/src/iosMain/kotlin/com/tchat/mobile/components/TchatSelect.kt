package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Popup
import androidx.compose.ui.window.PopupProperties
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatSelect with SwiftUI-inspired styling
 * Uses custom popup with iOS-style animations and interactions
 */
@Composable
actual fun TchatSelect(
    options: List<SelectOption>,
    selectedValues: List<String>,
    onSelectionChange: (List<String>) -> Unit,
    modifier: Modifier,
    mode: TchatSelectMode,
    validationState: TchatSelectValidationState,
    size: TchatSelectSize,
    enabled: Boolean,
    placeholder: String,
    label: String?,
    supportingText: String?,
    errorMessage: String?,
    searchEnabled: Boolean,
    searchPlaceholder: String,
    maxSelections: Int?,
    isLoading: Boolean,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    onTrailingIconClick: (() -> Unit)?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    var expanded by remember { mutableStateOf(false) }
    var searchQuery by remember { mutableStateOf("") }
    val focusManager = LocalFocusManager.current

    val filteredOptions = remember(options, searchQuery) {
        if (searchQuery.isEmpty()) {
            options
        } else {
            options.filter { option ->
                option.label.contains(searchQuery, ignoreCase = true) ||
                option.description?.contains(searchQuery, ignoreCase = true) == true
            }
        }
    }

    val fieldHeight = when (size) {
        TchatSelectSize.Small -> 44.dp // iOS uses larger minimum touch targets
        TchatSelectSize.Medium -> 52.dp
        TchatSelectSize.Large -> 60.dp
    }

    val textStyle = when (size) {
        TchatSelectSize.Small -> TchatTypography.typography.bodySmall
        TchatSelectSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSelectSize.Large -> TchatTypography.typography.bodyLarge
    }

    val borderColor by animateColorAsState(
        targetValue = when {
            expanded -> TchatColors.primary
            validationState == TchatSelectValidationState.Valid -> TchatColors.success
            validationState == TchatSelectValidationState.Invalid -> TchatColors.error
            else -> TchatColors.outline.copy(alpha = 0.4f) // iOS uses more subtle borders
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_border_color"
    )

    val arrowRotation by animateFloatAsState(
        targetValue = if (expanded) 180f else 0f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_arrow_rotation"
    )

    val fieldScale by animateFloatAsState(
        targetValue = if (expanded) 1.02f else 1f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_field_scale"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(6.dp) // iOS spacing
    ) {
        // Label
        label?.let { labelText ->
            Text(
                text = labelText,
                style = TchatTypography.typography.bodySmall.copy(
                    fontWeight = FontWeight.Medium // iOS labels are medium weight
                ),
                color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
            )
        }

        // Dropdown Field
        Box {
            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(fieldHeight)
                    .clip(RoundedCornerShape(12.dp)) // iOS uses more rounded corners
                    .clickable(
                        enabled = enabled && !isLoading,
                        interactionSource = interactionSource,
                        indication = null // iOS doesn't use ripple
                    ) {
                        expanded = !expanded
                        if (!expanded) {
                            searchQuery = ""
                            focusManager.clearFocus()
                        }
                    },
                shape = RoundedCornerShape(12.dp),
                color = if (enabled) TchatColors.surface else TchatColors.surface.copy(alpha = 0.5f),
                border = BorderStroke(
                    width = 1.5.dp, // iOS uses slightly thicker borders
                    color = borderColor
                ),
                shadowElevation = if (expanded) 4.dp else 0.dp // iOS-style elevation
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(horizontal = 16.dp, vertical = 12.dp), // iOS padding
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    leadingIcon?.let { icon ->
                        Icon(
                            imageVector = icon,
                            contentDescription = null,
                            modifier = Modifier.size(22.dp), // iOS uses slightly larger icons
                            tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                        )
                    }

                    Text(
                        text = getDisplayText(selectedValues, options, mode, placeholder),
                        modifier = Modifier.weight(1f),
                        style = textStyle,
                        color = if (selectedValues.isEmpty()) {
                            TchatColors.onSurface.copy(alpha = 0.6f) // iOS placeholder color
                        } else {
                            if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                        },
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )

                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        if (isLoading) {
                            CircularProgressIndicator(
                                modifier = Modifier.size(18.dp),
                                strokeWidth = 2.dp,
                                color = TchatColors.primary
                            )
                        } else {
                            trailingIcon?.let { icon ->
                                Icon(
                                    imageVector = icon,
                                    contentDescription = null,
                                    modifier = Modifier
                                        .size(22.dp)
                                        .clickable(
                                            enabled = enabled,
                                            interactionSource = remember { MutableInteractionSource() },
                                            indication = null
                                        ) {
                                            onTrailingIconClick?.invoke()
                                        },
                                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                                )
                            }

                            Icon(
                                imageVector = Icons.Default.KeyboardArrowDown,
                                contentDescription = if (expanded) "Collapse" else "Expand",
                                modifier = Modifier
                                    .size(24.dp)
                                    .rotate(arrowRotation),
                                tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                            )
                        }
                    }
                }
            }

            // iOS-style Dropdown Menu
            if (expanded) {
                Popup(
                    onDismissRequest = {
                        expanded = false
                        searchQuery = ""
                    },
                    properties = PopupProperties(focusable = true)
                ) {
                    Surface(
                        modifier = Modifier
                            .widthIn(min = 200.dp, max = 400.dp)
                            .heightIn(max = 300.dp)
                            .padding(top = 4.dp), // iOS spacing from field
                        shape = RoundedCornerShape(16.dp), // iOS uses more rounded popup corners
                        color = TchatColors.surface,
                        shadowElevation = 8.dp, // iOS-style shadow
                        border = BorderStroke(
                            width = 1.dp,
                            color = TchatColors.outline.copy(alpha = 0.2f)
                        )
                    ) {
                        Column {
                            // Search Field
                            if (searchEnabled && options.size > 5) {
                                BasicTextField(
                                    value = searchQuery,
                                    onValueChange = { searchQuery = it },
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .padding(16.dp)
                                        .background(
                                            color = TchatColors.surface.copy(alpha = 0.5f),
                                            shape = RoundedCornerShape(10.dp)
                                        )
                                        .border(
                                            width = 1.dp,
                                            color = TchatColors.outline.copy(alpha = 0.3f),
                                            shape = RoundedCornerShape(10.dp)
                                        )
                                        .padding(12.dp),
                                    textStyle = textStyle,
                                    decorationBox = { innerTextField ->
                                        Row(
                                            verticalAlignment = Alignment.CenterVertically,
                                            horizontalArrangement = Arrangement.spacedBy(10.dp)
                                        ) {
                                            Icon(
                                                imageVector = Icons.Default.Search,
                                                contentDescription = null,
                                                modifier = Modifier.size(18.dp),
                                                tint = TchatColors.onSurface.copy(alpha = 0.6f)
                                            )

                                            Box(Modifier.weight(1f)) {
                                                if (searchQuery.isEmpty()) {
                                                    Text(
                                                        text = searchPlaceholder,
                                                        style = textStyle,
                                                        color = TchatColors.onSurface.copy(alpha = 0.5f)
                                                    )
                                                }
                                                innerTextField()
                                            }
                                        }
                                    }
                                )

                                HorizontalDivider(
                                    thickness = 1.dp,
                                    color = TchatColors.outline.copy(alpha = 0.2f)
                                )
                            }

                            // Options List
                            LazyColumn(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .heightIn(max = 240.dp),
                                verticalArrangement = Arrangement.spacedBy(2.dp),
                                contentPadding = PaddingValues(vertical = 8.dp)
                            ) {
                                items(filteredOptions) { option ->
                                    IOSSelectOptionItem(
                                        option = option,
                                        isSelected = selectedValues.contains(option.value),
                                        isMultiSelect = mode == TchatSelectMode.Multiple,
                                        enabled = option.enabled,
                                        textStyle = textStyle,
                                        onClick = {
                                            handleOptionClick(
                                                option = option,
                                                selectedValues = selectedValues,
                                                mode = mode,
                                                maxSelections = maxSelections,
                                                onSelectionChange = onSelectionChange,
                                                onExpandedChange = { expanded = it }
                                            )
                                        }
                                    )
                                }

                                if (filteredOptions.isEmpty() && searchQuery.isNotEmpty()) {
                                    item {
                                        Text(
                                            text = "No options found",
                                            modifier = Modifier
                                                .fillMaxWidth()
                                                .padding(16.dp),
                                            style = textStyle,
                                            color = TchatColors.onSurface.copy(alpha = 0.6f)
                                        )
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

        // Supporting Text or Error Message
        when {
            validationState == TchatSelectValidationState.Invalid && errorMessage != null -> {
                Text(
                    text = errorMessage,
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.error
                )
            }
            supportingText != null -> {
                Text(
                    text = supportingText,
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    }
}

@Composable
actual fun TchatSingleSelect(
    options: List<SelectOption>,
    selectedValue: String?,
    onSelectionChange: (String?) -> Unit,
    modifier: Modifier,
    validationState: TchatSelectValidationState,
    size: TchatSelectSize,
    enabled: Boolean,
    placeholder: String,
    label: String?,
    supportingText: String?,
    errorMessage: String?,
    searchEnabled: Boolean,
    searchPlaceholder: String,
    isLoading: Boolean,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    onTrailingIconClick: (() -> Unit)?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    TchatSelect(
        options = options,
        selectedValues = selectedValue?.let { listOf(it) } ?: emptyList(),
        onSelectionChange = { values -> onSelectionChange(values.firstOrNull()) },
        modifier = modifier,
        mode = TchatSelectMode.Single,
        validationState = validationState,
        size = size,
        enabled = enabled,
        placeholder = placeholder,
        label = label,
        supportingText = supportingText,
        errorMessage = errorMessage,
        searchEnabled = searchEnabled,
        searchPlaceholder = searchPlaceholder,
        maxSelections = 1,
        isLoading = isLoading,
        leadingIcon = leadingIcon,
        trailingIcon = trailingIcon,
        onTrailingIconClick = onTrailingIconClick,
        interactionSource = interactionSource,
        contentDescription = contentDescription
    )
}

@Composable
private fun IOSSelectOptionItem(
    option: SelectOption,
    isSelected: Boolean,
    isMultiSelect: Boolean,
    enabled: Boolean,
    textStyle: androidx.compose.ui.text.TextStyle,
    onClick: () -> Unit
) {
    val backgroundColor by animateColorAsState(
        targetValue = if (isSelected) {
            TchatColors.primary.copy(alpha = 0.08f)
        } else {
            Color.Transparent
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_selection_background"
    )

    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp))
            .clickable(
                enabled = enabled,
                interactionSource = remember { MutableInteractionSource() },
                indication = null // iOS doesn't use ripple
            ) { onClick() },
        color = backgroundColor
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 12.dp), // iOS padding
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            option.icon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(22.dp),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                )
            }

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = option.label,
                    style = textStyle.copy(
                        fontWeight = if (isSelected) FontWeight.Medium else FontWeight.Normal
                    ),
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f),
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )

                option.description?.let { desc ->
                    Text(
                        text = desc,
                        style = TchatTypography.typography.bodySmall,
                        color = if (enabled) TchatColors.onSurface.copy(alpha = 0.6f) else TchatColors.onSurface.copy(alpha = 0.3f),
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                }
            }

            if (isSelected) {
                Icon(
                    imageVector = if (isMultiSelect) Icons.Default.CheckBox else Icons.Default.Check,
                    contentDescription = "Selected",
                    modifier = Modifier.size(22.dp),
                    tint = TchatColors.primary
                )
            }
        }
    }
}

private fun getDisplayText(
    selectedValues: List<String>,
    options: List<SelectOption>,
    mode: TchatSelectMode,
    placeholder: String
): String {
    return when {
        selectedValues.isEmpty() -> placeholder
        mode == TchatSelectMode.Single -> {
            options.firstOrNull { it.value == selectedValues.first() }?.label ?: placeholder
        }
        selectedValues.size == 1 -> {
            options.firstOrNull { it.value == selectedValues.first() }?.label ?: placeholder
        }
        else -> "${selectedValues.size} selected"
    }
}

private fun handleOptionClick(
    option: SelectOption,
    selectedValues: List<String>,
    mode: TchatSelectMode,
    maxSelections: Int?,
    onSelectionChange: (List<String>) -> Unit,
    onExpandedChange: (Boolean) -> Unit
) {
    when (mode) {
        TchatSelectMode.Single -> {
            onSelectionChange(listOf(option.value))
            onExpandedChange(false)
        }
        TchatSelectMode.Multiple -> {
            val newSelection = if (selectedValues.contains(option.value)) {
                selectedValues - option.value
            } else {
                if (maxSelections != null && selectedValues.size >= maxSelections) {
                    selectedValues
                } else {
                    selectedValues + option.value
                }
            }
            onSelectionChange(newSelection)
        }
    }
}