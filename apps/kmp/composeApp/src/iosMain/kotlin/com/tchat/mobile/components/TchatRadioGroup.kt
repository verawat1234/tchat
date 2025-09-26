package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatRadioGroup with SwiftUI-inspired styling
 * Uses custom drawing for iOS-style radio buttons with spring animations
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
        TchatRadioGroupSize.Small -> 18.dp // iOS uses slightly larger touch targets
        TchatRadioGroupSize.Medium -> 22.dp
        TchatRadioGroupSize.Large -> 26.dp
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
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_border_color"
    )

    val borderAlpha by animateFloatAsState(
        targetValue = when (validationState) {
            TchatRadioGroupValidationState.None -> 0f
            else -> 1f
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_border_alpha"
    )

    Column(
        modifier = modifier
            .then(
                if (validationState != TchatRadioGroupValidationState.None) {
                    Modifier
                        .border(
                            width = 1.dp,
                            color = borderColor.copy(alpha = borderAlpha),
                            shape = RoundedCornerShape(12.dp) // iOS uses more rounded corners
                        )
                        .background(
                            color = borderColor.copy(alpha = 0.03f * borderAlpha), // More subtle background
                            shape = RoundedCornerShape(12.dp)
                        )
                        .padding(16.dp) // iOS uses more padding
                } else {
                    Modifier
                }
            ),
        verticalArrangement = Arrangement.spacedBy((spacing + 2).dp) // iOS uses more spacing
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(10.dp) // iOS spacing
            ) {
                Text(
                    text = labelText,
                    style = textStyle.copy(
                        fontWeight = androidx.compose.ui.text.font.FontWeight.Medium // iOS uses medium weight for labels
                    ),
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f),
                    modifier = Modifier.weight(1f)
                )

                when (validationState) {
                    TchatRadioGroupValidationState.Valid -> {
                        Icon(
                            imageVector = Icons.Default.CheckCircle,
                            contentDescription = "Valid",
                            tint = TchatColors.success,
                            modifier = Modifier.size(18.dp) // iOS uses slightly larger icons
                        )
                    }
                    TchatRadioGroupValidationState.Invalid -> {
                        Icon(
                            imageVector = Icons.Default.Error,
                            contentDescription = "Invalid",
                            tint = TchatColors.error,
                            modifier = Modifier.size(18.dp)
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
                    verticalArrangement = Arrangement.spacedBy((spacing + 2).dp)
                ) {
                    options.forEach { option ->
                        IOSRadioOptionItem(
                            option = option,
                            isSelected = selectedValue == option.value,
                            onSelectionChange = onSelectionChange,
                            enabled = enabled && option.enabled,
                            radioButtonSize = radioButtonSize,
                            textStyle = textStyle
                        )
                    }
                }
            }
            TchatRadioGroupOrientation.Horizontal -> {
                Row(
                    horizontalArrangement = Arrangement.spacedBy((spacing + 2).dp),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    options.forEach { option ->
                        IOSRadioOptionItem(
                            option = option,
                            isSelected = selectedValue == option.value,
                            onSelectionChange = onSelectionChange,
                            enabled = enabled && option.enabled,
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
                    color = TchatColors.onSurface.copy(alpha = 0.6f) // iOS uses less transparency
                )
            }
        }
    }
}

@Composable
private fun IOSRadioOptionItem(
    option: RadioOption,
    isSelected: Boolean,
    onSelectionChange: (String) -> Unit,
    enabled: Boolean,
    radioButtonSize: androidx.compose.ui.unit.Dp,
    textStyle: androidx.compose.ui.text.TextStyle,
    modifier: Modifier = Modifier
) {
    // iOS-style spring animation for selection
    val selectionProgress by animateFloatAsState(
        targetValue = if (isSelected) 1f else 0f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_selection_animation"
    )

    Row(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp)) // iOS uses rounded touch areas
            .selectable(
                selected = isSelected,
                onClick = { if (enabled) onSelectionChange(option.value) },
                enabled = enabled,
                role = Role.RadioButton
            )
            .padding(vertical = 6.dp, horizontal = 4.dp), // iOS padding
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp) // iOS spacing
    ) {
        IOSRadioButtonIcon(
            isSelected = isSelected,
            enabled = enabled,
            size = radioButtonSize,
            selectionProgress = selectionProgress
        )

        Column(
            modifier = Modifier.weight(1f)
        ) {
            Text(
                text = option.label,
                style = textStyle,
                color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f) // iOS disabled transparency
            )

            option.description?.let { desc ->
                Text(
                    text = desc,
                    style = TchatTypography.typography.bodySmall,
                    color = if (enabled) TchatColors.onSurface.copy(alpha = 0.6f) else TchatColors.onSurface.copy(alpha = 0.3f)
                )
            }
        }
    }
}

@Composable
private fun IOSRadioButtonIcon(
    isSelected: Boolean,
    enabled: Boolean,
    size: androidx.compose.ui.unit.Dp,
    selectionProgress: Float
) {
    val primaryColor = TchatColors.primary
    val borderColor = if (isSelected) primaryColor else TchatColors.outline.copy(alpha = 0.5f) // iOS uses more subtle borders
    val centerColor = primaryColor

    Canvas(
        modifier = Modifier.size(size)
    ) {
        drawIOSRadioButton(
            borderColor = borderColor,
            centerColor = centerColor,
            selectionProgress = selectionProgress,
            enabled = enabled
        )
    }
}

private fun DrawScope.drawIOSRadioButton(
    borderColor: Color,
    centerColor: Color,
    selectionProgress: Float,
    enabled: Boolean
) {
    val center = Offset(size.width / 2, size.height / 2)
    val radius = size.minDimension / 2
    val strokeWidth = 1.5.dp.toPx() // iOS uses thinner strokes
    val alpha = if (enabled) 1f else 0.5f

    // Draw outer circle border
    drawCircle(
        color = borderColor.copy(alpha = alpha),
        radius = radius,
        center = center,
        style = Stroke(width = strokeWidth)
    )

    // Draw inner filled circle with animation
    if (selectionProgress > 0f) {
        val innerRadius = radius * 0.4f * selectionProgress // iOS-style inner circle sizing
        drawCircle(
            color = centerColor.copy(alpha = alpha),
            radius = innerRadius,
            center = center
        )
    }
}