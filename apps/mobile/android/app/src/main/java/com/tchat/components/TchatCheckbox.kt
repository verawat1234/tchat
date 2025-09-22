package com.tchat.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Remove
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
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
 * Checkbox component following Tchat design system
 */
@Composable
fun TchatCheckbox(
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier,
    label: String? = null,
    description: String? = null,
    size: TchatCheckboxSize = TchatCheckboxSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true
) {
    TchatCheckbox(
        state = if (checked) TchatCheckboxState.Checked else TchatCheckboxState.Unchecked,
        onStateChange = { newState ->
            onCheckedChange(newState == TchatCheckboxState.Checked)
        },
        modifier = modifier,
        label = label,
        description = description,
        size = size,
        validationState = validationState,
        enabled = enabled
    )
}

/**
 * Advanced checkbox with state support
 */
@Composable
fun TchatCheckbox(
    state: TchatCheckboxState,
    onStateChange: (TchatCheckboxState) -> Unit,
    modifier: Modifier = Modifier,
    label: String? = null,
    description: String? = null,
    size: TchatCheckboxSize = TchatCheckboxSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true
) {
    val hapticFeedback = LocalHapticFeedback.current

    // Animation values
    val backgroundColor by animateColorAsState(
        targetValue = when {
            !enabled -> Colors.surface.copy(alpha = 0.5f)
            state == TchatCheckboxState.Unchecked -> Colors.background
            else -> Colors.primary
        },
        animationSpec = tween(150),
        label = "background_color"
    )

    val borderColor by animateColorAsState(
        targetValue = when {
            !enabled -> Colors.border.copy(alpha = 0.5f)
            validationState is TchatValidationState.Invalid -> Colors.borderError
            validationState is TchatValidationState.Valid -> Colors.success
            state == TchatCheckboxState.Unchecked -> Colors.border
            else -> Colors.primary
        },
        animationSpec = tween(150),
        label = "border_color"
    )

    val scale by animateFloatAsState(
        targetValue = 1f,
        animationSpec = tween(150),
        label = "scale"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.xs)
    ) {
        // Main checkbox row
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(4.dp))
                .clickable(enabled = enabled) {
                    val newState = when (state) {
                        TchatCheckboxState.Unchecked -> TchatCheckboxState.Checked
                        TchatCheckboxState.Checked -> TchatCheckboxState.Unchecked
                        TchatCheckboxState.Indeterminate -> TchatCheckboxState.Checked
                    }
                    onStateChange(newState)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                },
            verticalAlignment = Alignment.Top,
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            // Checkbox icon
            Box(
                modifier = Modifier
                    .size(size.checkboxSize)
                    .scale(scale)
                    .background(backgroundColor, RoundedCornerShape(4.dp))
                    .border(
                        BorderStroke(
                            width = when {
                                validationState !is TchatValidationState.None -> 2.dp
                                state == TchatCheckboxState.Unchecked -> 2.dp
                                else -> 0.dp
                            },
                            color = borderColor
                        ),
                        RoundedCornerShape(4.dp)
                    ),
                contentAlignment = Alignment.Center
            ) {
                if (state != TchatCheckboxState.Unchecked) {
                    Icon(
                        imageVector = when (state) {
                            TchatCheckboxState.Checked -> Icons.Default.Check
                            TchatCheckboxState.Indeterminate -> Icons.Default.Remove
                            else -> Icons.Default.Check
                        },
                        contentDescription = when (state) {
                            TchatCheckboxState.Checked -> "Checked"
                            TchatCheckboxState.Indeterminate -> "Indeterminate"
                            else -> null
                        },
                        tint = if (enabled) Colors.textOnPrimary else Colors.textDisabled,
                        modifier = Modifier.size(size.iconSize)
                    )
                }
            }

            // Label and description
            if (label != null || description != null) {
                Column(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.spacedBy(Spacing.xs)
                ) {
                    label?.let { labelText ->
                        Text(
                            text = labelText,
                            fontSize = size.fontSize,
                            color = if (enabled) Colors.textPrimary else Colors.textDisabled,
                            lineHeight = size.lineHeight
                        )
                    }

                    description?.let { descriptionText ->
                        Text(
                            text = descriptionText,
                            fontSize = size.descriptionFontSize,
                            color = if (enabled) Colors.textSecondary else Colors.textDisabled,
                            lineHeight = size.descriptionLineHeight
                        )
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
                modifier = Modifier.padding(start = size.checkboxSize + Spacing.sm)
            )
        }
    }
}

/**
 * Checkbox group component
 */
@Composable
fun <T> TchatCheckboxGroup(
    selection: Set<T>,
    onSelectionChange: (Set<T>) -> Unit,
    options: List<T>,
    optionLabel: (T) -> String,
    modifier: Modifier = Modifier,
    optionDescription: ((T) -> String)? = null,
    title: String? = null,
    size: TchatCheckboxSize = TchatCheckboxSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    maxSelections: Int? = null
) {
    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.sm)
    ) {
        // Group title
        title?.let { titleText ->
            Text(
                text = titleText,
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold,
                color = if (enabled) Colors.textPrimary else Colors.textDisabled
            )
        }

        // Options
        Column(
            verticalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            options.forEachIndexed { index, option ->
                TchatCheckbox(
                    checked = selection.contains(option),
                    onCheckedChange = { isChecked ->
                        val newSelection = if (isChecked) {
                            if (maxSelections != null && selection.size >= maxSelections) {
                                return@TchatCheckbox // Max selections reached
                            }
                            selection + option
                        } else {
                            selection - option
                        }
                        onSelectionChange(newSelection)
                    },
                    label = optionLabel(option),
                    description = optionDescription?.invoke(option),
                    size = size,
                    validationState = if (index == 0) validationState else TchatValidationState.None,
                    enabled = enabled
                )
            }
        }

        // Group validation message
        if (validationState is TchatValidationState.Invalid) {
            Text(
                text = validationState.message,
                fontSize = 12.sp,
                color = Colors.error
            )
        }
    }
}

/**
 * Checkbox state definitions
 */
enum class TchatCheckboxState {
    Unchecked,
    Checked,
    Indeterminate
}

/**
 * Checkbox size definitions
 */
enum class TchatCheckboxSize(
    val checkboxSize: androidx.compose.ui.unit.Dp,
    val iconSize: androidx.compose.ui.unit.Dp,
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val lineHeight: androidx.compose.ui.unit.TextUnit,
    val descriptionFontSize: androidx.compose.ui.unit.TextUnit,
    val descriptionLineHeight: androidx.compose.ui.unit.TextUnit
) {
    Small(
        checkboxSize = 16.dp,
        iconSize = 10.dp,
        fontSize = 14.sp,
        lineHeight = 18.sp,
        descriptionFontSize = 12.sp,
        descriptionLineHeight = 16.sp
    ),
    Medium(
        checkboxSize = 20.dp,
        iconSize = 12.dp,
        fontSize = 16.sp,
        lineHeight = 22.sp,
        descriptionFontSize = 14.sp,
        descriptionLineHeight = 18.sp
    ),
    Large(
        checkboxSize = 24.dp,
        iconSize = 16.dp,
        fontSize = 18.sp,
        lineHeight = 26.sp,
        descriptionFontSize = 16.sp,
        descriptionLineHeight = 22.sp
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatCheckboxPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Basic checkbox
        TchatCheckbox(
            checked = true,
            onCheckedChange = { },
            label = "Accept terms and conditions"
        )

        // Checkbox with description
        TchatCheckbox(
            checked = false,
            onCheckedChange = { },
            label = "Enable notifications",
            description = "Receive updates about new messages and mentions"
        )

        // Indeterminate state
        TchatCheckbox(
            state = TchatCheckboxState.Indeterminate,
            onStateChange = { },
            label = "Select all items",
            description = "Some items are selected"
        )

        // Validation states
        TchatCheckbox(
            checked = true,
            onCheckedChange = { },
            label = "Valid selection",
            validationState = TchatValidationState.Valid
        )

        TchatCheckbox(
            checked = false,
            onCheckedChange = { },
            label = "Required field",
            validationState = TchatValidationState.Invalid("This field is required")
        )

        // Disabled state
        TchatCheckbox(
            checked = true,
            onCheckedChange = { },
            label = "Disabled checkbox",
            description = "This option cannot be changed",
            enabled = false
        )

        // Different sizes
        Column(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
            TchatCheckbox(
                checked = true,
                onCheckedChange = { },
                label = "Small checkbox",
                size = TchatCheckboxSize.Small
            )

            TchatCheckbox(
                checked = true,
                onCheckedChange = { },
                label = "Medium checkbox",
                size = TchatCheckboxSize.Medium
            )

            TchatCheckbox(
                checked = true,
                onCheckedChange = { },
                label = "Large checkbox",
                size = TchatCheckboxSize.Large
            )
        }

        Divider()

        // Checkbox group
        val languages = listOf("Swift", "Kotlin", "JavaScript", "Python", "Go")
        val selectedLanguages = setOf("Swift", "Kotlin")

        TchatCheckboxGroup(
            selection = selectedLanguages,
            onSelectionChange = { },
            options = languages,
            optionLabel = { it },
            optionDescription = { lang ->
                when (lang) {
                    "Swift" -> "iOS and macOS development"
                    "Kotlin" -> "Android and multiplatform development"
                    "JavaScript" -> "Web and Node.js development"
                    "Python" -> "Data science and backend development"
                    "Go" -> "System programming and microservices"
                    else -> ""
                }
            },
            title = "Select programming languages",
            maxSelections = 3
        )
    }
}