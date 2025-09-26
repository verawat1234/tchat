package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
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
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.PopupProperties
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatSelect using Material3 ExposedDropdownMenu
 * Provides native Material Design dropdown with comprehensive theming and animations
 */
@OptIn(ExperimentalMaterial3Api::class)
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
        TchatSelectSize.Small -> 40.dp
        TchatSelectSize.Medium -> 48.dp
        TchatSelectSize.Large -> 56.dp
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
            else -> TchatColors.outline
        },
        animationSpec = tween(300),
        label = "BorderColor"
    )

    val arrowRotation by animateFloatAsState(
        targetValue = if (expanded) 180f else 0f,
        animationSpec = tween(300),
        label = "ArrowRotation"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        // Label
        label?.let { labelText ->
            Text(
                text = labelText,
                style = TchatTypography.typography.bodySmall,
                color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
            )
        }

        // Dropdown Field
        ExposedDropdownMenuBox(
            expanded = expanded,
            onExpandedChange = {
                if (enabled && !isLoading) {
                    expanded = it
                    if (!expanded) {
                        searchQuery = ""
                        focusManager.clearFocus()
                    }
                }
            }
        ) {
            OutlinedTextField(
                value = getDisplayText(selectedValues, options, mode, placeholder),
                onValueChange = { },
                modifier = Modifier
                    .menuAnchor()
                    .fillMaxWidth()
                    .height(fieldHeight),
                readOnly = true,
                enabled = enabled,
                textStyle = textStyle,
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = borderColor,
                    unfocusedBorderColor = borderColor,
                    disabledBorderColor = TchatColors.outline.copy(alpha = 0.12f),
                    focusedLabelColor = borderColor,
                    cursorColor = TchatColors.primary
                ),
                leadingIcon = leadingIcon?.let { icon ->
                    {
                        Icon(
                            imageVector = icon,
                            contentDescription = null,
                            modifier = Modifier.size(20.dp),
                            tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                        )
                    }
                },
                trailingIcon = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        if (isLoading) {
                            CircularProgressIndicator(
                                modifier = Modifier.size(16.dp),
                                strokeWidth = 2.dp,
                                color = TchatColors.primary
                            )
                        } else {
                            trailingIcon?.let { icon ->
                                IconButton(
                                    onClick = { onTrailingIconClick?.invoke() },
                                    enabled = enabled
                                ) {
                                    Icon(
                                        imageVector = icon,
                                        contentDescription = null,
                                        modifier = Modifier.size(20.dp),
                                        tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                                    )
                                }
                            }

                            Icon(
                                imageVector = Icons.Default.ArrowDropDown,
                                contentDescription = if (expanded) "Collapse" else "Expand",
                                modifier = Modifier
                                    .size(24.dp)
                                    .rotate(arrowRotation),
                                tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                            )
                        }
                    }
                },
                supportingText = when {
                    validationState == TchatSelectValidationState.Invalid && errorMessage != null -> {
                        { Text(errorMessage, color = TchatColors.error) }
                    }
                    supportingText != null -> {
                        { Text(supportingText, color = TchatColors.onSurface.copy(alpha = 0.7f)) }
                    }
                    else -> null
                }
            )

            // Dropdown Menu
            ExposedDropdownMenu(
                expanded = expanded,
                onDismissRequest = {
                    expanded = false
                    searchQuery = ""
                }
            ) {
                // Search Field
                if (searchEnabled && options.size > 5) {
                    BasicTextField(
                        value = searchQuery,
                        onValueChange = { searchQuery = it },
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(12.dp)
                            .background(
                                color = TchatColors.surface,
                                shape = RoundedCornerShape(8.dp)
                            )
                            .border(
                                width = 1.dp,
                                color = TchatColors.outline.copy(alpha = 0.5f),
                                shape = RoundedCornerShape(8.dp)
                            )
                            .padding(12.dp),
                        textStyle = textStyle,
                        decorationBox = { innerTextField ->
                            Row(
                                verticalAlignment = Alignment.CenterVertically,
                                horizontalArrangement = Arrangement.spacedBy(8.dp)
                            ) {
                                Icon(
                                    imageVector = Icons.Default.Search,
                                    contentDescription = null,
                                    modifier = Modifier.size(16.dp),
                                    tint = TchatColors.onSurface.copy(alpha = 0.6f)
                                )

                                Box(Modifier.weight(1f)) {
                                    if (searchQuery.isEmpty()) {
                                        Text(
                                            text = searchPlaceholder,
                                            style = textStyle,
                                            color = TchatColors.onSurface.copy(alpha = 0.6f)
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
                    verticalArrangement = Arrangement.spacedBy(2.dp)
                ) {
                    items(filteredOptions) { option ->
                        SelectOptionItem(
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
private fun SelectOptionItem(
    option: SelectOption,
    isSelected: Boolean,
    isMultiSelect: Boolean,
    enabled: Boolean,
    textStyle: androidx.compose.ui.text.TextStyle,
    onClick: () -> Unit
) {
    DropdownMenuItem(
        text = {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                option.icon?.let { icon ->
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        modifier = Modifier.size(20.dp),
                        tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                    )
                }

                Column(
                    modifier = Modifier.weight(1f)
                ) {
                    Text(
                        text = option.label,
                        style = textStyle,
                        color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f),
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )

                    option.description?.let { desc ->
                        Text(
                            text = desc,
                            style = TchatTypography.typography.bodySmall,
                            color = if (enabled) TchatColors.onSurface.copy(alpha = 0.7f) else TchatColors.onSurface.copy(alpha = 0.3f),
                            maxLines = 1,
                            overflow = TextOverflow.Ellipsis
                        )
                    }
                }

                if (isSelected) {
                    Icon(
                        imageVector = if (isMultiSelect) Icons.Default.CheckBox else Icons.Default.Check,
                        contentDescription = "Selected",
                        modifier = Modifier.size(20.dp),
                        tint = TchatColors.primary
                    )
                }
            }
        },
        onClick = onClick,
        enabled = enabled
    )
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