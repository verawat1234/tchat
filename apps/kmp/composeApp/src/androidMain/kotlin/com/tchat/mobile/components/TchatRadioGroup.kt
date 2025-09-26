package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.selectable
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Error
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatRadioGroup using Material3 RadioButton
 * Provides native Material Design radio group with comprehensive theming and animations
 */
@Composable
actual fun TchatRadioGroup(
    options: List<RadioOption>,
    selectedValue: String?,
    onSelectionChange: (String) -> Unit,
    modifier: Modifier,
    orientation: TchatRadioGroupOrientation,
    validationState: TchatRadioGroupValidationState,
    size: TchatRadioGroupSize,
    enabled: Boolean,
    label: String?,
    supportingText: String?,
    errorMessage: String?,
    spacing: Int,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val radioButtonSize = when (size) {
        TchatRadioGroupSize.Small -> 16.dp
        TchatRadioGroupSize.Medium -> 20.dp
        TchatRadioGroupSize.Large -> 24.dp
    }

    val textStyle = when (size) {
        TchatRadioGroupSize.Small -> TchatTypography.typography.bodySmall
        TchatRadioGroupSize.Medium -> TchatTypography.typography.bodyMedium
        TchatRadioGroupSize.Large -> TchatTypography.typography.bodyLarge
    }

    val borderColor by animateColorAsState(
        targetValue = when (validationState) {
            TchatRadioGroupValidationState.Valid -> TchatColors.success
            TchatRadioGroupValidationState.Invalid -> TchatColors.error
            TchatRadioGroupValidationState.None -> Color.Transparent
        },
        animationSpec = tween(300),
        label = "BorderColor"
    )

    val borderAlpha by animateFloatAsState(
        targetValue = when (validationState) {
            TchatRadioGroupValidationState.None -> 0f
            else -> 1f
        },
        animationSpec = tween(300),
        label = "BorderAlpha"
    )

    val radioColors = RadioButtonDefaults.colors(
        selectedColor = TchatColors.primary,
        unselectedColor = TchatColors.outline,
        disabledSelectedColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledUnselectedColor = TchatColors.onSurface.copy(alpha = 0.38f)
    )

    Column(
        modifier = modifier
            .then(
                if (validationState != TchatRadioGroupValidationState.None) {
                    Modifier
                        .border(
                            width = 1.dp,
                            color = borderColor.copy(alpha = borderAlpha),
                            shape = RoundedCornerShape(8.dp)
                        )
                        .background(
                            color = borderColor.copy(alpha = 0.05f * borderAlpha),
                            shape = RoundedCornerShape(8.dp)
                        )
                        .padding(12.dp)
                } else {
                    Modifier
                }
            ),
        verticalArrangement = Arrangement.spacedBy(spacing.dp)
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = labelText,
                    style = textStyle,
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f),
                    modifier = Modifier.weight(1f)
                )

                when (validationState) {
                    TchatRadioGroupValidationState.Valid -> {
                        Icon(
                            imageVector = Icons.Default.CheckCircle,
                            contentDescription = "Valid",
                            tint = TchatColors.success,
                            modifier = Modifier.size(16.dp)
                        )
                    }
                    TchatRadioGroupValidationState.Invalid -> {
                        Icon(
                            imageVector = Icons.Default.Error,
                            contentDescription = "Invalid",
                            tint = TchatColors.error,
                            modifier = Modifier.size(16.dp)
                        )
                    }
                    TchatRadioGroupValidationState.None -> {}
                }
            }
        }

        // Radio Options
        when (orientation) {
            TchatRadioGroupOrientation.Vertical -> {
                Column(
                    verticalArrangement = Arrangement.spacedBy(spacing.dp)
                ) {
                    options.forEach { option ->
                        RadioOptionItem(
                            option = option,
                            isSelected = selectedValue == option.value,
                            onSelectionChange = onSelectionChange,
                            enabled = enabled && option.enabled,
                            radioColors = radioColors,
                            radioButtonSize = radioButtonSize,
                            textStyle = textStyle
                        )
                    }
                }
            }
            TchatRadioGroupOrientation.Horizontal -> {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(spacing.dp),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    options.forEach { option ->
                        RadioOptionItem(
                            option = option,
                            isSelected = selectedValue == option.value,
                            onSelectionChange = onSelectionChange,
                            enabled = enabled && option.enabled,
                            radioColors = radioColors,
                            radioButtonSize = radioButtonSize,
                            textStyle = textStyle,
                            modifier = Modifier.weight(1f)
                        )
                    }
                }
            }
        }

        // Supporting Text or Error Message
        when {
            validationState == TchatRadioGroupValidationState.Invalid && errorMessage != null -> {
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
                    color = TchatColors.onSurface.copy(alpha = 0.7f)
                )
            }
        }
    }
}

@Composable
private fun RadioOptionItem(
    option: RadioOption,
    isSelected: Boolean,
    onSelectionChange: (String) -> Unit,
    enabled: Boolean,
    radioColors: RadioButtonColors,
    radioButtonSize: androidx.compose.ui.unit.Dp,
    textStyle: androidx.compose.ui.text.TextStyle,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .selectable(
                selected = isSelected,
                onClick = { if (enabled) onSelectionChange(option.value) },
                enabled = enabled,
                role = Role.RadioButton
            )
            .padding(vertical = 4.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        RadioButton(
            selected = isSelected,
            onClick = null, // Handled by selectable modifier
            modifier = Modifier.size(radioButtonSize),
            enabled = enabled,
            colors = radioColors
        )

        Column(
            modifier = Modifier.weight(1f)
        ) {
            Text(
                text = option.label,
                style = textStyle,
                color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
            )

            option.description?.let { desc ->
                Text(
                    text = desc,
                    style = TchatTypography.typography.bodySmall,
                    color = if (enabled) TchatColors.onSurface.copy(alpha = 0.7f) else TchatColors.onSurface.copy(alpha = 0.3f)
                )
            }
        }
    }
}