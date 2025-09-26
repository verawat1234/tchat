package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.CornerRadius
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.*
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatCheckbox with SwiftUI-inspired styling
 * Uses custom drawing for iOS-style rounded corners and checkmark animation
 */
@Composable
actual fun TchatCheckbox(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatCheckboxSize,
    interactionSource: MutableInteractionSource,
    label: String?
) {
    val checkboxSize = when (size) {
        TchatCheckboxSize.Small -> 18.dp // iOS uses slightly larger touch targets
        TchatCheckboxSize.Medium -> 22.dp
        TchatCheckboxSize.Large -> 26.dp
    }

    // iOS-style spring animation for checkmark
    val checkmarkProgress by animateFloatAsState(
        targetValue = if (checked) 1f else 0f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_checkmark_animation"
    )

    if (label != null) {
        Row(
            modifier = modifier,
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp) // iOS uses more spacing
        ) {
            IOSCheckboxIcon(
                checked = checked,
                enabled = enabled,
                size = checkboxSize,
                checkmarkProgress = checkmarkProgress,
                interactionSource = interactionSource,
                onCheckedChange = onCheckedChange
            )

            Text(
                text = label,
                style = when (size) {
                    TchatCheckboxSize.Small -> TchatTypography.typography.bodySmall
                    TchatCheckboxSize.Medium -> TchatTypography.typography.bodyMedium
                    TchatCheckboxSize.Large -> TchatTypography.typography.bodyLarge
                },
                color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f), // iOS uses less transparent disabled text
                modifier = Modifier
                    .weight(1f)
                    .let { textModifier ->
                        if (enabled && onCheckedChange != null) {
                            textModifier.clickable(
                                interactionSource = interactionSource,
                                indication = null // iOS doesn't show ripple on text
                            ) {
                                onCheckedChange(!checked)
                            }
                        } else textModifier
                    }
            )
        }
    } else {
        IOSCheckboxIcon(
            checked = checked,
            enabled = enabled,
            size = checkboxSize,
            checkmarkProgress = checkmarkProgress,
            interactionSource = interactionSource,
            onCheckedChange = onCheckedChange,
            modifier = modifier
        )
    }
}

@Composable
private fun IOSCheckboxIcon(
    checked: Boolean,
    enabled: Boolean,
    size: androidx.compose.ui.unit.Dp,
    checkmarkProgress: Float,
    interactionSource: MutableInteractionSource,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier = Modifier
) {
    val primaryColor = TchatColors.primary
    val backgroundColor = if (checked) primaryColor else Color.Transparent
    val borderColor = if (checked) primaryColor else TchatColors.outline.copy(alpha = 0.6f) // iOS uses more subtle borders
    val checkmarkColor = TchatColors.onPrimary

    Canvas(
        modifier = modifier
            .size(size)
            .clip(RoundedCornerShape(4.dp)) // iOS uses slightly rounded corners
            .let { canvasModifier ->
                if (enabled && onCheckedChange != null) {
                    canvasModifier.clickable(
                        interactionSource = interactionSource,
                        indication = null // iOS doesn't use ripple effects
                    ) {
                        onCheckedChange(!checked)
                    }
                } else canvasModifier
            }
    ) {
        drawIOSCheckbox(
            backgroundColor = backgroundColor,
            borderColor = borderColor,
            checkmarkColor = checkmarkColor,
            checkmarkProgress = checkmarkProgress,
            enabled = enabled
        )
    }
}

private fun DrawScope.drawIOSCheckbox(
    backgroundColor: Color,
    borderColor: Color,
    checkmarkColor: Color,
    checkmarkProgress: Float,
    enabled: Boolean
) {
    val cornerRadius = 4.dp.toPx() // iOS-style corner radius
    val strokeWidth = 1.5.dp.toPx() // iOS uses thinner strokes
    val alpha = if (enabled) 1f else 0.5f

    // Draw background
    drawRoundRect(
        color = backgroundColor.copy(alpha = alpha),
        size = size,
        cornerRadius = CornerRadius(cornerRadius)
    )

    // Draw border
    drawRoundRect(
        color = borderColor.copy(alpha = alpha),
        size = size,
        cornerRadius = CornerRadius(cornerRadius),
        style = Stroke(width = strokeWidth)
    )

    // Draw checkmark with iOS-style path
    if (checkmarkProgress > 0f) {
        val checkmarkPath = Path().apply {
            val checkmarkSize = size.minDimension * 0.5f
            val centerX = size.width / 2
            val centerY = size.height / 2
            val startX = centerX - checkmarkSize * 0.3f
            val startY = centerY
            val middleX = centerX - checkmarkSize * 0.1f
            val middleY = centerY + checkmarkSize * 0.3f
            val endX = centerX + checkmarkSize * 0.4f
            val endY = centerY - checkmarkSize * 0.2f

            moveTo(startX, startY)
            lineTo(
                startX + (middleX - startX) * checkmarkProgress,
                startY + (middleY - startY) * checkmarkProgress
            )
            if (checkmarkProgress > 0.5f) {
                val secondProgress = (checkmarkProgress - 0.5f) * 2f
                lineTo(
                    middleX + (endX - middleX) * secondProgress,
                    middleY + (endY - middleY) * secondProgress
                )
            }
        }

        drawPath(
            path = checkmarkPath,
            color = checkmarkColor.copy(alpha = alpha),
            style = Stroke(
                width = 2.dp.toPx(),
                cap = StrokeCap.Round,
                join = StrokeJoin.Round
            )
        )
    }
}