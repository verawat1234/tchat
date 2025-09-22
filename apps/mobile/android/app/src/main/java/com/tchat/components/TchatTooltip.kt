package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectTapGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.layout.onGloballyPositioned
import androidx.compose.ui.layout.positionInWindow
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.IntRect
import androidx.compose.ui.unit.IntSize
import androidx.compose.ui.unit.LayoutDirection
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Popup
import androidx.compose.ui.window.PopupPositionProvider
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import kotlinx.coroutines.delay

/**
 * Tooltip component following Tchat design system
 */
@Composable
fun TchatTooltip(
    text: String,
    modifier: Modifier = Modifier,
    position: TchatTooltipPosition = TchatTooltipPosition.Auto,
    trigger: TchatTooltipTrigger = TchatTooltipTrigger.LongPress,
    style: TchatTooltipStyle = TchatTooltipStyle.Dark,
    maxWidth: androidx.compose.ui.unit.Dp = 250.dp,
    showArrow: Boolean = true,
    delay: Long = 500L,
    autoDismissDelay: Long? = 3000L,
    onShow: (() -> Unit)? = null,
    onDismiss: (() -> Unit)? = null,
    content: @Composable () -> Unit
) {
    val hapticFeedback = LocalHapticFeedback.current
    val density = LocalDensity.current

    var isVisible by remember { mutableStateOf(false) }
    var targetPosition by remember { mutableStateOf(Offset.Zero) }
    var targetSize by remember { mutableStateOf(IntSize.Zero) }

    val tooltipColors = getTooltipColors(style)

    fun showTooltip() {
        isVisible = true
        onShow?.invoke()
        hapticFeedback.performHapticFeedback(
            androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
        )
    }

    fun hideTooltip() {
        isVisible = false
        onDismiss?.invoke()
    }

    LaunchedEffect(isVisible) {
        if (isVisible && autoDismissDelay != null) {
            delay(autoDismissDelay)
            hideTooltip()
        }
    }

    Box(
        modifier = modifier
            .onGloballyPositioned { coordinates ->
                targetPosition = coordinates.positionInWindow()
                targetSize = coordinates.size
            }
            .then(
                when (trigger) {
                    TchatTooltipTrigger.Tap -> Modifier.clickable { showTooltip() }
                    TchatTooltipTrigger.LongPress -> Modifier.pointerInput(Unit) {
                        detectTapGestures(
                            onLongPress = { showTooltip() }
                        )
                    }
                    TchatTooltipTrigger.Manual -> Modifier
                }
            )
    ) {
        content()

        if (isVisible) {
            Popup(
                popupPositionProvider = TooltipPositionProvider(
                    targetPosition = targetPosition,
                    targetSize = targetSize,
                    preferredPosition = position,
                    arrowSize = if (showArrow) 8.dp else 0.dp,
                    spacing = 8.dp
                ),
                onDismissRequest = { hideTooltip() }
            ) {
                TooltipContent(
                    text = text,
                    style = style,
                    maxWidth = maxWidth,
                    showArrow = showArrow,
                    tooltipColors = tooltipColors
                )
            }
        }
    }
}

@Composable
private fun TooltipContent(
    text: String,
    style: TchatTooltipStyle,
    maxWidth: androidx.compose.ui.unit.Dp,
    showArrow: Boolean,
    tooltipColors: TooltipColors
) {
    AnimatedVisibility(
        visible = true,
        enter = fadeIn(animationSpec = tween(200)) + scaleIn(
            initialScale = 0.8f,
            animationSpec = tween(200)
        ),
        exit = fadeOut(animationSpec = tween(150)) + scaleOut(
            targetScale = 0.8f,
            animationSpec = tween(150)
        )
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // Arrow (top)
            if (showArrow) {
                Box(
                    modifier = Modifier
                        .size(16.dp, 8.dp)
                        .background(
                            color = tooltipColors.background,
                            shape = TriangleShape()
                        )
                )
            }

            // Tooltip content
            Text(
                text = text,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = tooltipColors.text,
                textAlign = TextAlign.Center,
                modifier = Modifier
                    .widthIn(max = maxWidth)
                    .shadow(
                        elevation = 8.dp,
                        shape = RoundedCornerShape(8.dp)
                    )
                    .background(
                        color = tooltipColors.background,
                        shape = RoundedCornerShape(8.dp)
                    )
                    .then(
                        if (style == TchatTooltipStyle.Light) {
                            Modifier.border(
                                width = 1.dp,
                                color = tooltipColors.border,
                                shape = RoundedCornerShape(8.dp)
                            )
                        } else Modifier
                    )
                    .padding(horizontal = 12.dp, vertical = 8.dp)
            )
        }
    }
}

private fun getTooltipColors(style: TchatTooltipStyle): TooltipColors {
    return when (style) {
        TchatTooltipStyle.Dark -> TooltipColors(
            background = Color.Black.copy(alpha = 0.9f),
            text = Color.White,
            border = Color.Transparent
        )
        TchatTooltipStyle.Light -> TooltipColors(
            background = Colors.background,
            text = Colors.textPrimary,
            border = Colors.border
        )
        TchatTooltipStyle.Info -> TooltipColors(
            background = Colors.primary,
            text = Colors.textOnPrimary,
            border = Colors.primary
        )
        TchatTooltipStyle.Warning -> TooltipColors(
            background = Colors.warning,
            text = Colors.textOnPrimary,
            border = Colors.warning
        )
        TchatTooltipStyle.Error -> TooltipColors(
            background = Colors.error,
            text = Colors.textOnPrimary,
            border = Colors.error
        )
    }
}

private class TooltipPositionProvider(
    private val targetPosition: Offset,
    private val targetSize: IntSize,
    private val preferredPosition: TchatTooltipPosition,
    private val arrowSize: androidx.compose.ui.unit.Dp,
    private val spacing: androidx.compose.ui.unit.Dp
) : PopupPositionProvider {
    override fun calculatePosition(
        anchorBounds: IntRect,
        windowSize: IntSize,
        layoutDirection: LayoutDirection,
        popupContentSize: IntSize
    ): IntOffset {
        val spacingPx = spacing.value.toInt()
        val arrowSizePx = arrowSize.value.toInt()

        val targetCenterX = targetPosition.x.toInt() + targetSize.width / 2
        val targetCenterY = targetPosition.y.toInt() + targetSize.height / 2

        // Calculate position based on preferred position or auto-detect best position
        val position = if (preferredPosition == TchatTooltipPosition.Auto) {
            calculateBestPosition(
                targetPosition = targetPosition,
                targetSize = targetSize,
                popupSize = popupContentSize,
                windowSize = windowSize,
                spacingPx = spacingPx,
                arrowSizePx = arrowSizePx
            )
        } else {
            preferredPosition
        }

        return when (position) {
            TchatTooltipPosition.Top -> IntOffset(
                x = targetCenterX - popupContentSize.width / 2,
                y = targetPosition.y.toInt() - popupContentSize.height - spacingPx - arrowSizePx
            )
            TchatTooltipPosition.Bottom -> IntOffset(
                x = targetCenterX - popupContentSize.width / 2,
                y = targetPosition.y.toInt() + targetSize.height + spacingPx + arrowSizePx
            )
            TchatTooltipPosition.Leading -> IntOffset(
                x = targetPosition.x.toInt() - popupContentSize.width - spacingPx,
                y = targetCenterY - popupContentSize.height / 2
            )
            TchatTooltipPosition.Trailing -> IntOffset(
                x = targetPosition.x.toInt() + targetSize.width + spacingPx,
                y = targetCenterY - popupContentSize.height / 2
            )
            TchatTooltipPosition.Auto -> IntOffset(
                x = targetCenterX - popupContentSize.width / 2,
                y = targetPosition.y.toInt() - popupContentSize.height - spacingPx - arrowSizePx
            )
        }.let { offset ->
            // Ensure tooltip stays within screen bounds
            IntOffset(
                x = offset.x.coerceIn(16, windowSize.width - popupContentSize.width - 16),
                y = offset.y.coerceIn(16, windowSize.height - popupContentSize.height - 16)
            )
        }
    }

    private fun calculateBestPosition(
        targetPosition: Offset,
        targetSize: IntSize,
        popupSize: IntSize,
        windowSize: IntSize,
        spacingPx: Int,
        arrowSizePx: Int
    ): TchatTooltipPosition {
        val spacing = spacingPx + arrowSizePx

        // Check available space in each direction
        val spaceTop = targetPosition.y - spacing - popupSize.height
        val spaceBottom = windowSize.height - targetPosition.y - targetSize.height - spacing - popupSize.height
        val spaceLeading = targetPosition.x - spacing - popupSize.width
        val spaceTrailing = windowSize.width - targetPosition.x - targetSize.width - spacing - popupSize.width

        // Prioritize top/bottom over leading/trailing
        return when {
            spaceTop >= 0 -> TchatTooltipPosition.Top
            spaceBottom >= 0 -> TchatTooltipPosition.Bottom
            spaceTrailing >= 0 -> TchatTooltipPosition.Trailing
            spaceLeading >= 0 -> TchatTooltipPosition.Leading
            else -> TchatTooltipPosition.Top // Default fallback
        }
    }
}

private class TriangleShape : Shape {
    override fun createOutline(
        size: androidx.compose.ui.geometry.Size,
        layoutDirection: LayoutDirection,
        density: androidx.compose.ui.unit.Density
    ): androidx.compose.ui.graphics.Outline {
        val path = Path().apply {
            moveTo(size.width / 2f, 0f)
            lineTo(0f, size.height)
            lineTo(size.width, size.height)
            close()
        }
        return androidx.compose.ui.graphics.Outline.Generic(path)
    }
}

/**
 * Extension function to add tooltip to any composable
 */
fun Modifier.tchatTooltip(
    text: String,
    position: TchatTooltipPosition = TchatTooltipPosition.Auto,
    trigger: TchatTooltipTrigger = TchatTooltipTrigger.LongPress,
    style: TchatTooltipStyle = TchatTooltipStyle.Dark,
    maxWidth: androidx.compose.ui.unit.Dp = 250.dp,
    showArrow: Boolean = true,
    delay: Long = 500L,
    autoDismissDelay: Long? = 3000L,
    onShow: (() -> Unit)? = null,
    onDismiss: (() -> Unit)? = null
): Modifier = this.then(
    Modifier // This would need to be implemented as a custom modifier
)

/**
 * Data classes and enums
 */
enum class TchatTooltipPosition {
    Top,
    Bottom,
    Leading,
    Trailing,
    Auto
}

enum class TchatTooltipTrigger {
    Tap,
    LongPress,
    Manual
}

enum class TchatTooltipStyle {
    Dark,
    Light,
    Info,
    Warning,
    Error
}

private data class TooltipColors(
    val background: Color,
    val text: Color,
    val border: Color
)

// Preview
@Preview(showBackground = true)
@Composable
fun TchatTooltipPreview() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(Spacing.lg),
        verticalArrangement = Arrangement.spacedBy(Spacing.xl),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Dark tooltip (default)
        TchatTooltip(
            text = "This is a dark tooltip with helpful information",
            position = TchatTooltipPosition.Top,
            trigger = TchatTooltipTrigger.Tap,
            style = TchatTooltipStyle.Dark
        ) {
            Button(
                onClick = { },
                colors = ButtonDefaults.buttonColors(containerColor = Colors.primary)
            ) {
                Text("Dark Tooltip", color = Colors.textOnPrimary)
            }
        }

        // Light tooltip
        TchatTooltip(
            text = "Light tooltip with border and light background",
            position = TchatTooltipPosition.Bottom,
            trigger = TchatTooltipTrigger.Tap,
            style = TchatTooltipStyle.Light
        ) {
            Button(
                onClick = { },
                colors = ButtonDefaults.buttonColors(containerColor = Colors.surface)
            ) {
                Text("Light Tooltip", color = Colors.textPrimary)
            }
        }

        // Info tooltip
        TchatTooltip(
            text = "This is an informational tooltip",
            position = TchatTooltipPosition.Trailing,
            trigger = TchatTooltipTrigger.Tap,
            style = TchatTooltipStyle.Info
        ) {
            Icon(
                imageVector = Icons.Default.Info,
                contentDescription = "Info",
                tint = Colors.primary,
                modifier = Modifier
                    .size(32.dp)
                    .clickable { }
            )
        }

        // Warning tooltip
        TchatTooltip(
            text = "Warning: This action cannot be undone",
            position = TchatTooltipPosition.Leading,
            trigger = TchatTooltipTrigger.Tap,
            style = TchatTooltipStyle.Warning
        ) {
            Icon(
                imageVector = Icons.Default.Warning,
                contentDescription = "Warning",
                tint = Colors.warning,
                modifier = Modifier
                    .size(32.dp)
                    .clickable { }
            )
        }

        // Error tooltip
        TchatTooltip(
            text = "Error: Something went wrong",
            position = TchatTooltipPosition.Auto,
            trigger = TchatTooltipTrigger.Tap,
            style = TchatTooltipStyle.Error
        ) {
            Icon(
                imageVector = Icons.Default.Error,
                contentDescription = "Error",
                tint = Colors.error,
                modifier = Modifier
                    .size(32.dp)
                    .clickable { }
            )
        }

        // Long press tooltip
        TchatTooltip(
            text = "This tooltip appears on long press",
            trigger = TchatTooltipTrigger.LongPress,
            delay = 500L
        ) {
            Button(
                onClick = { },
                colors = ButtonDefaults.buttonColors(containerColor = Colors.success)
            ) {
                Text("Long Press Me", color = Colors.textOnPrimary)
            }
        }
    }
}