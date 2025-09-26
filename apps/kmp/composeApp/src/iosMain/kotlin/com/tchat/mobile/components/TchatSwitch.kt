package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.toggleable
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.CornerRadius
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatSwitch with SwiftUI-inspired styling
 * Uses custom drawing for iOS-style rounded track and smooth thumb animation
 */
@Composable
actual fun TchatSwitch(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatSwitchSize,
    isLoading: Boolean,
    label: String?,
    description: String?,
    leadingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val switchWidth = when (size) {
        TchatSwitchSize.Small -> 36.dp
        TchatSwitchSize.Medium -> 44.dp // iOS standard switch width
        TchatSwitchSize.Large -> 52.dp
    }

    val switchHeight = when (size) {
        TchatSwitchSize.Small -> 20.dp
        TchatSwitchSize.Medium -> 26.dp // iOS standard switch height
        TchatSwitchSize.Large -> 32.dp
    }

    val textStyle = when (size) {
        TchatSwitchSize.Small -> TchatTypography.typography.bodySmall
        TchatSwitchSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSwitchSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSwitchSize.Small -> 18.dp
        TchatSwitchSize.Medium -> 22.dp
        TchatSwitchSize.Large -> 26.dp
    }

    // iOS-style spring animations
    val thumbPosition by animateFloatAsState(
        targetValue = if (checked) 1f else 0f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_thumb_position"
    )

    val trackColor by animateColorAsState(
        targetValue = if (checked) {
            TchatColors.primary
        } else {
            TchatColors.outline.copy(alpha = 0.3f) // iOS uses more subtle unchecked color
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_track_color"
    )

    val thumbScale by animateFloatAsState(
        targetValue = if (isLoading) 0.8f else 1f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_thumb_scale"
    )

    val switchComponent = @Composable {
        Box(
            contentAlignment = Alignment.Center
        ) {
            Canvas(
                modifier = Modifier
                    .size(width = switchWidth, height = switchHeight)
                    .clip(RoundedCornerShape(switchHeight / 2)) // iOS uses fully rounded switches
                    .clickable(
                        enabled = enabled && !isLoading,
                        indication = null, // iOS doesn't use ripple
                        interactionSource = interactionSource
                    ) {
                        onCheckedChange?.invoke(!checked)
                    }
            ) {
                drawIOSSwitch(
                    trackColor = trackColor,
                    thumbPosition = thumbPosition,
                    thumbScale = thumbScale,
                    enabled = enabled && !isLoading
                )
            }

            // Loading indicator
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.size((iconSize.value * 0.5f).dp),
                    strokeWidth = 1.5.dp, // iOS uses thinner stroke
                    color = TchatColors.primary
                )
            }
        }
    }

    if (label != null || description != null || leadingIcon != null) {
        Row(
            modifier = modifier
                .toggleable(
                    value = checked,
                    enabled = enabled && !isLoading,
                    role = Role.Switch,
                    onValueChange = onCheckedChange ?: {}
                )
                .padding(vertical = 6.dp), // iOS uses more vertical padding
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(14.dp) // iOS uses more spacing
        ) {
            // Leading icon
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) {
                        if (checked) TchatColors.primary else TchatColors.onSurface
                    } else {
                        TchatColors.onSurface.copy(alpha = 0.5f) // iOS uses less transparency
                    }
                )
            }

            // Label and description
            Column(
                modifier = Modifier.weight(1f),
                verticalArrangement = Arrangement.spacedBy(3.dp) // iOS spacing
            ) {
                label?.let { labelText ->
                    Text(
                        text = labelText,
                        style = textStyle.copy(
                            fontWeight = FontWeight.Medium // iOS uses medium weight for labels
                        ),
                        color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                    )
                }

                description?.let { desc ->
                    Text(
                        text = desc,
                        style = TchatTypography.typography.bodySmall,
                        color = if (enabled) TchatColors.onSurface.copy(alpha = 0.6f) else TchatColors.onSurface.copy(alpha = 0.3f)
                    )
                }
            }

            // Switch
            switchComponent()
        }
    } else {
        Box(
            modifier = modifier
        ) {
            switchComponent()
        }
    }
}

private fun DrawScope.drawIOSSwitch(
    trackColor: Color,
    thumbPosition: Float,
    thumbScale: Float,
    enabled: Boolean
) {
    val trackWidth = size.width
    val trackHeight = size.height
    val trackRadius = trackHeight / 2
    val alpha = if (enabled) 1f else 0.5f

    // Draw track
    drawRoundRect(
        color = trackColor.copy(alpha = alpha),
        size = Size(trackWidth, trackHeight),
        cornerRadius = CornerRadius(trackRadius)
    )

    // Draw thumb
    val thumbRadius = (trackHeight - 4.dp.toPx()) / 2 * thumbScale // iOS thumb is slightly smaller than track
    val thumbCenterX = thumbRadius + 2.dp.toPx() + thumbPosition * (trackWidth - thumbRadius * 2 - 4.dp.toPx())
    val thumbCenterY = trackHeight / 2

    // iOS-style thumb shadow
    drawCircle(
        color = Color.Black.copy(alpha = 0.15f * alpha),
        radius = thumbRadius + 1.dp.toPx(),
        center = Offset(thumbCenterX, thumbCenterY + 1.dp.toPx())
    )

    // Thumb
    drawCircle(
        color = Color.White.copy(alpha = alpha),
        radius = thumbRadius,
        center = Offset(thumbCenterX, thumbCenterY)
    )
}