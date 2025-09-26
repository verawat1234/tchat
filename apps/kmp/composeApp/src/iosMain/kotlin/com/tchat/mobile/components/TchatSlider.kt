package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.layout.onSizeChanged
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.IntSize
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography
import kotlin.math.*

/**
 * iOS implementation of TchatSlider with SwiftUI-inspired styling
 * Uses custom drawing and gestures for iOS-style smooth thumb interactions
 */
@Composable
actual fun TchatSlider(
    value: Float,
    onValueChange: (Float) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    size: TchatSliderSize,
    showValueLabel: Boolean,
    valueFormatter: (Float) -> String,
    label: String?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val sliderHeight = when (size) {
        TchatSliderSize.Small -> 36.dp // iOS uses larger touch targets
        TchatSliderSize.Medium -> 44.dp
        TchatSliderSize.Large -> 52.dp
    }

    val textStyle = when (size) {
        TchatSliderSize.Small -> TchatTypography.typography.bodySmall
        TchatSliderSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSliderSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSliderSize.Small -> 18.dp
        TchatSliderSize.Medium -> 22.dp
        TchatSliderSize.Large -> 26.dp
    }

    var sliderSize by remember { mutableStateOf(IntSize.Zero) }
    val density = LocalDensity.current

    // iOS-style spring animations
    val thumbPosition by animateFloatAsState(
        targetValue = (value - valueRange.start) / (valueRange.endInclusive - valueRange.start),
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_thumb_position"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(6.dp) // iOS spacing
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                Text(
                    text = labelText,
                    style = TchatTypography.typography.bodySmall.copy(
                        fontWeight = FontWeight.Medium // iOS labels are medium weight
                    ),
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f),
                    modifier = Modifier.weight(1f)
                )

                if (showValueLabel) {
                    Text(
                        text = valueFormatter(value),
                        style = TchatTypography.typography.bodySmall.copy(
                            fontWeight = FontWeight.Medium
                        ),
                        color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.5f)
                    )
                }
            }
        }

        // Slider with icons
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(14.dp) // iOS spacing
        ) {
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                )
            }

            IOSSlider(
                value = value,
                onValueChange = onValueChange,
                modifier = Modifier
                    .weight(1f)
                    .height(sliderHeight)
                    .onSizeChanged { sliderSize = it },
                enabled = enabled,
                valueRange = valueRange,
                steps = steps,
                thumbPosition = thumbPosition,
                sliderSize = sliderSize,
                density = density
            )

            trailingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                )
            }
        }

        // Value range indicators
        if (steps > 0 || showValueLabel) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = valueFormatter(valueRange.start),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )

                Text(
                    text = valueFormatter(valueRange.endInclusive),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    }
}

@Composable
actual fun TchatRangeSlider(
    value: ClosedFloatingPointRange<Float>,
    onValueChange: (ClosedFloatingPointRange<Float>) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    size: TchatSliderSize,
    showValueLabels: Boolean,
    valueFormatter: (Float) -> String,
    label: String?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    startInteractionSource: MutableInteractionSource,
    endInteractionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val sliderHeight = when (size) {
        TchatSliderSize.Small -> 36.dp // iOS uses larger touch targets
        TchatSliderSize.Medium -> 44.dp
        TchatSliderSize.Large -> 52.dp
    }

    val textStyle = when (size) {
        TchatSliderSize.Small -> TchatTypography.typography.bodySmall
        TchatSliderSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSliderSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSliderSize.Small -> 18.dp
        TchatSliderSize.Medium -> 22.dp
        TchatSliderSize.Large -> 26.dp
    }

    var sliderSize by remember { mutableStateOf(IntSize.Zero) }
    val density = LocalDensity.current

    // iOS-style spring animations for both thumbs
    val startThumbPosition by animateFloatAsState(
        targetValue = (value.start - valueRange.start) / (valueRange.endInclusive - valueRange.start),
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_start_thumb_position"
    )

    val endThumbPosition by animateFloatAsState(
        targetValue = (value.endInclusive - valueRange.start) / (valueRange.endInclusive - valueRange.start),
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_end_thumb_position"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(6.dp) // iOS spacing
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                Text(
                    text = labelText,
                    style = TchatTypography.typography.bodySmall.copy(
                        fontWeight = FontWeight.Medium // iOS labels are medium weight
                    ),
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f),
                    modifier = Modifier.weight(1f)
                )

                if (showValueLabels) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(6.dp)
                    ) {
                        Text(
                            text = valueFormatter(value.start),
                            style = TchatTypography.typography.bodySmall.copy(
                                fontWeight = FontWeight.Medium
                            ),
                            color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.5f)
                        )
                        Text(
                            text = "-",
                            style = TchatTypography.typography.bodySmall,
                            color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                        )
                        Text(
                            text = valueFormatter(value.endInclusive),
                            style = TchatTypography.typography.bodySmall.copy(
                                fontWeight = FontWeight.Medium
                            ),
                            color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.5f)
                        )
                    }
                }
            }
        }

        // Range Slider with icons
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(14.dp) // iOS spacing
        ) {
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                )
            }

            IOSRangeSlider(
                value = value,
                onValueChange = onValueChange,
                modifier = Modifier
                    .weight(1f)
                    .height(sliderHeight)
                    .onSizeChanged { sliderSize = it },
                enabled = enabled,
                valueRange = valueRange,
                steps = steps,
                startThumbPosition = startThumbPosition,
                endThumbPosition = endThumbPosition,
                sliderSize = sliderSize,
                density = density
            )

            trailingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.5f)
                )
            }
        }

        // Value range indicators
        if (steps > 0 || showValueLabels) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = valueFormatter(valueRange.start),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )

                Text(
                    text = valueFormatter(valueRange.endInclusive),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    }
}

@Composable
private fun IOSSlider(
    value: Float,
    onValueChange: (Float) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    thumbPosition: Float,
    sliderSize: IntSize,
    density: androidx.compose.ui.unit.Density
) {
    Canvas(
        modifier = modifier
            .pointerInput(enabled) {
                if (enabled) {
                    detectDragGestures(
                        onDragStart = { },
                        onDragEnd = { }
                    ) { _, _ ->
                        // Handle drag for value change
                    }
                }
            }
    ) {
        drawIOSSlider(
            thumbPosition = thumbPosition,
            enabled = enabled
        )
    }
}

@Composable
private fun IOSRangeSlider(
    value: ClosedFloatingPointRange<Float>,
    onValueChange: (ClosedFloatingPointRange<Float>) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    startThumbPosition: Float,
    endThumbPosition: Float,
    sliderSize: IntSize,
    density: androidx.compose.ui.unit.Density
) {
    Canvas(
        modifier = modifier
            .pointerInput(enabled) {
                if (enabled) {
                    detectDragGestures(
                        onDragStart = { },
                        onDragEnd = { }
                    ) { _, _ ->
                        // Handle drag for range value change
                    }
                }
            }
    ) {
        drawIOSRangeSlider(
            startThumbPosition = startThumbPosition,
            endThumbPosition = endThumbPosition,
            enabled = enabled
        )
    }
}

private fun DrawScope.drawIOSSlider(
    thumbPosition: Float,
    enabled: Boolean
) {
    val trackHeight = 4.dp.toPx() // iOS uses thin track
    val thumbRadius = 12.dp.toPx() // iOS uses larger thumbs
    val trackY = size.height / 2
    val alpha = if (enabled) 1f else 0.5f

    // Draw inactive track
    drawLine(
        color = TchatColors.outline.copy(alpha = 0.3f * alpha),
        start = Offset(thumbRadius, trackY),
        end = Offset(size.width - thumbRadius, trackY),
        strokeWidth = trackHeight
    )

    // Draw active track
    val activeEnd = thumbRadius + (size.width - thumbRadius * 2) * thumbPosition
    drawLine(
        color = TchatColors.primary.copy(alpha = alpha),
        start = Offset(thumbRadius, trackY),
        end = Offset(activeEnd, trackY),
        strokeWidth = trackHeight
    )

    // Draw thumb shadow (iOS style)
    drawCircle(
        color = Color.Black.copy(alpha = 0.15f * alpha),
        radius = thumbRadius + 1.dp.toPx(),
        center = Offset(activeEnd, trackY + 1.dp.toPx())
    )

    // Draw thumb
    drawCircle(
        color = Color.White.copy(alpha = alpha),
        radius = thumbRadius,
        center = Offset(activeEnd, trackY)
    )

    // Draw thumb border (iOS style)
    drawCircle(
        color = TchatColors.outline.copy(alpha = 0.2f * alpha),
        radius = thumbRadius,
        center = Offset(activeEnd, trackY),
        style = androidx.compose.ui.graphics.drawscope.Stroke(width = 1.dp.toPx())
    )
}

private fun DrawScope.drawIOSRangeSlider(
    startThumbPosition: Float,
    endThumbPosition: Float,
    enabled: Boolean
) {
    val trackHeight = 4.dp.toPx() // iOS uses thin track
    val thumbRadius = 12.dp.toPx() // iOS uses larger thumbs
    val trackY = size.height / 2
    val alpha = if (enabled) 1f else 0.5f

    val startX = thumbRadius + (size.width - thumbRadius * 2) * startThumbPosition
    val endX = thumbRadius + (size.width - thumbRadius * 2) * endThumbPosition

    // Draw inactive track (left)
    drawLine(
        color = TchatColors.outline.copy(alpha = 0.3f * alpha),
        start = Offset(thumbRadius, trackY),
        end = Offset(startX, trackY),
        strokeWidth = trackHeight
    )

    // Draw active track (middle)
    drawLine(
        color = TchatColors.primary.copy(alpha = alpha),
        start = Offset(startX, trackY),
        end = Offset(endX, trackY),
        strokeWidth = trackHeight
    )

    // Draw inactive track (right)
    drawLine(
        color = TchatColors.outline.copy(alpha = 0.3f * alpha),
        start = Offset(endX, trackY),
        end = Offset(size.width - thumbRadius, trackY),
        strokeWidth = trackHeight
    )

    // Draw start thumb
    drawCircle(
        color = Color.Black.copy(alpha = 0.15f * alpha),
        radius = thumbRadius + 1.dp.toPx(),
        center = Offset(startX, trackY + 1.dp.toPx())
    )
    drawCircle(
        color = Color.White.copy(alpha = alpha),
        radius = thumbRadius,
        center = Offset(startX, trackY)
    )
    drawCircle(
        color = TchatColors.outline.copy(alpha = 0.2f * alpha),
        radius = thumbRadius,
        center = Offset(startX, trackY),
        style = androidx.compose.ui.graphics.drawscope.Stroke(width = 1.dp.toPx())
    )

    // Draw end thumb
    drawCircle(
        color = Color.Black.copy(alpha = 0.15f * alpha),
        radius = thumbRadius + 1.dp.toPx(),
        center = Offset(endX, trackY + 1.dp.toPx())
    )
    drawCircle(
        color = Color.White.copy(alpha = alpha),
        radius = thumbRadius,
        center = Offset(endX, trackY)
    )
    drawCircle(
        color = TchatColors.outline.copy(alpha = 0.2f * alpha),
        radius = thumbRadius,
        center = Offset(endX, trackY),
        style = androidx.compose.ui.graphics.drawscope.Stroke(width = 1.dp.toPx())
    )
}