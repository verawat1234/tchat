package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import kotlinx.coroutines.delay
import kotlin.math.roundToInt

/**
 * Toast notification component following Tchat design system
 */
@Composable
fun TchatToast(
    isVisible: Boolean,
    type: TchatToastType,
    message: String,
    modifier: Modifier = Modifier,
    position: TchatToastPosition = TchatToastPosition.Top,
    style: TchatToastStyle = TchatToastStyle.Filled,
    icon: ImageVector? = null,
    duration: Long = 3000L,
    isDismissible: Boolean = true,
    hapticFeedback: Boolean = true,
    onDismiss: (() -> Unit)? = null
) {
    val haptic = LocalHapticFeedback.current
    val density = LocalDensity.current

    var dragOffset by remember { mutableStateOf(0f) }
    val dismissThreshold = with(density) { 100.dp.toPx() }

    val toastColors = getToastColors(type, style)
    val toastIcon = icon ?: getDefaultIcon(type)

    LaunchedEffect(isVisible) {
        if (isVisible) {
            if (hapticFeedback) {
                val feedbackType = when (type) {
                    TchatToastType.Success -> androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    TchatToastType.Error -> androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    TchatToastType.Warning -> androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    TchatToastType.Info -> androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                }
                haptic.performHapticFeedback(feedbackType)
            }

            if (duration > 0) {
                delay(duration)
                onDismiss?.invoke()
            }
        }
    }

    AnimatedVisibility(
        visible = isVisible,
        enter = slideInVertically(
            initialOffsetY = { if (position == TchatToastPosition.Top) -it else it },
            animationSpec = tween(400)
        ) + fadeIn(animationSpec = tween(400)),
        exit = slideOutVertically(
            targetOffsetY = { if (position == TchatToastPosition.Top) -it else it },
            animationSpec = tween(300)
        ) + fadeOut(animationSpec = tween(300))
    ) {
        Row(
            modifier = modifier
                .fillMaxWidth()
                .offset { IntOffset(dragOffset.roundToInt(), 0) }
                .shadow(
                    elevation = 8.dp,
                    shape = RoundedCornerShape(8.dp)
                )
                .background(
                    color = toastColors.background,
                    shape = RoundedCornerShape(8.dp)
                )
                .then(
                    if (style == TchatToastStyle.Outlined) {
                        Modifier.border(
                            width = 1.dp,
                            color = toastColors.border,
                            shape = RoundedCornerShape(8.dp)
                        )
                    } else Modifier
                )
                .padding(horizontal = Spacing.md, vertical = Spacing.sm)
                .then(
                    if (isDismissible) {
                        Modifier.pointerInput(Unit) {
                            detectDragGestures(
                                onDragEnd = {
                                    val shouldDismiss = kotlin.math.abs(dragOffset) > dismissThreshold
                                    if (shouldDismiss) {
                                        onDismiss?.invoke()
                                    } else {
                                        dragOffset = 0f
                                    }
                                }
                            ) { _, dragAmount ->
                                dragOffset += dragAmount.x
                            }
                        }
                    } else Modifier
                ),
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Icon
            Icon(
                imageVector = toastIcon,
                contentDescription = null,
                tint = toastColors.icon,
                modifier = Modifier.size(16.dp)
            )

            // Message
            Text(
                text = message,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = toastColors.text,
                modifier = Modifier.weight(1f)
            )

            // Dismiss button
            if (isDismissible) {
                Icon(
                    imageVector = Icons.Default.Close,
                    contentDescription = "Dismiss",
                    tint = toastColors.text.copy(alpha = 0.7f),
                    modifier = Modifier
                        .size(16.dp)
                        .clickable { onDismiss?.invoke() }
                )
            }
        }
    }
}

private fun getToastColors(type: TchatToastType, style: TchatToastStyle): ToastColors {
    return when (type to style) {
        TchatToastType.Info to TchatToastStyle.Filled -> ToastColors(
            background = Colors.primary,
            border = Colors.primary,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatToastType.Info to TchatToastStyle.Outlined -> ToastColors(
            background = Colors.background,
            border = Colors.primary,
            text = Colors.primary,
            icon = Colors.primary
        )
        TchatToastType.Info to TchatToastStyle.Minimal -> ToastColors(
            background = Colors.primary.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.primary,
            icon = Colors.primary
        )

        TchatToastType.Success to TchatToastStyle.Filled -> ToastColors(
            background = Colors.success,
            border = Colors.success,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatToastType.Success to TchatToastStyle.Outlined -> ToastColors(
            background = Colors.background,
            border = Colors.success,
            text = Colors.success,
            icon = Colors.success
        )
        TchatToastType.Success to TchatToastStyle.Minimal -> ToastColors(
            background = Colors.success.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.success,
            icon = Colors.success
        )

        TchatToastType.Warning to TchatToastStyle.Filled -> ToastColors(
            background = Colors.warning,
            border = Colors.warning,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatToastType.Warning to TchatToastStyle.Outlined -> ToastColors(
            background = Colors.background,
            border = Colors.warning,
            text = Colors.warning,
            icon = Colors.warning
        )
        TchatToastType.Warning to TchatToastStyle.Minimal -> ToastColors(
            background = Colors.warning.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.warning,
            icon = Colors.warning
        )

        TchatToastType.Error to TchatToastStyle.Filled -> ToastColors(
            background = Colors.error,
            border = Colors.error,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatToastType.Error to TchatToastStyle.Outlined -> ToastColors(
            background = Colors.background,
            border = Colors.error,
            text = Colors.error,
            icon = Colors.error
        )
        TchatToastType.Error to TchatToastStyle.Minimal -> ToastColors(
            background = Colors.error.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.error,
            icon = Colors.error
        )
        else -> ToastColors(
            background = Colors.background,
            border = Colors.border,
            text = Colors.textPrimary,
            icon = Colors.textPrimary
        )
    }
}

private fun getDefaultIcon(type: TchatToastType): ImageVector {
    return when (type) {
        TchatToastType.Info -> Icons.Default.Info
        TchatToastType.Success -> Icons.Default.CheckCircle
        TchatToastType.Warning -> Icons.Default.Warning
        TchatToastType.Error -> Icons.Default.Error
    }
}

/**
 * Toast manager for global toast handling
 */
object TchatToastManager {
    private var _currentToasts = mutableStateOf<List<TchatToastItem>>(emptyList())
    val currentToasts: State<List<TchatToastItem>> = _currentToasts

    fun show(
        type: TchatToastType,
        message: String,
        position: TchatToastPosition = TchatToastPosition.Top,
        style: TchatToastStyle = TchatToastStyle.Filled,
        icon: ImageVector? = null,
        duration: Long = 3000L,
        isDismissible: Boolean = true,
        hapticFeedback: Boolean = true
    ) {
        val toast = TchatToastItem(
            type = type,
            message = message,
            position = position,
            style = style,
            icon = icon,
            duration = duration,
            isDismissible = isDismissible,
            hapticFeedback = hapticFeedback
        )

        _currentToasts.value = _currentToasts.value + toast
    }

    fun dismiss(id: String) {
        _currentToasts.value = _currentToasts.value.filter { it.id != id }
    }

    fun dismissAll() {
        _currentToasts.value = emptyList()
    }

    fun showSuccess(
        message: String,
        position: TchatToastPosition = TchatToastPosition.Top,
        duration: Long = 2000L
    ) {
        show(
            type = TchatToastType.Success,
            message = message,
            position = position,
            duration = duration
        )
    }

    fun showError(
        message: String,
        position: TchatToastPosition = TchatToastPosition.Top,
        duration: Long = 4000L
    ) {
        show(
            type = TchatToastType.Error,
            message = message,
            position = position,
            duration = duration
        )
    }

    fun showWarning(
        message: String,
        position: TchatToastPosition = TchatToastPosition.Top,
        duration: Long = 3000L
    ) {
        show(
            type = TchatToastType.Warning,
            message = message,
            position = position,
            duration = duration
        )
    }

    fun showInfo(
        message: String,
        position: TchatToastPosition = TchatToastPosition.Top,
        duration: Long = 3000L
    ) {
        show(
            type = TchatToastType.Info,
            message = message,
            position = position,
            duration = duration
        )
    }
}

/**
 * Toast overlay composable for global toasts
 */
@Composable
fun TchatToastOverlay() {
    val currentToasts by TchatToastManager.currentToasts

    BoxWithConstraints(
        modifier = Modifier.fillMaxSize()
    ) {
        // Top toasts
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .align(Alignment.TopCenter)
                .padding(horizontal = Spacing.md)
                .padding(top = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs)
        ) {
            currentToasts
                .filter { it.position == TchatToastPosition.Top }
                .forEach { toast ->
                    key(toast.id) {
                        TchatToast(
                            isVisible = true,
                            type = toast.type,
                            message = toast.message,
                            position = toast.position,
                            style = toast.style,
                            icon = toast.icon,
                            duration = 0, // Managed by toast manager
                            isDismissible = toast.isDismissible,
                            hapticFeedback = false, // Already triggered
                            onDismiss = { TchatToastManager.dismiss(toast.id) }
                        )
                    }
                }
        }

        // Center toasts
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .align(Alignment.Center)
                .padding(horizontal = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs)
        ) {
            currentToasts
                .filter { it.position == TchatToastPosition.Center }
                .forEach { toast ->
                    key(toast.id) {
                        TchatToast(
                            isVisible = true,
                            type = toast.type,
                            message = toast.message,
                            position = toast.position,
                            style = toast.style,
                            icon = toast.icon,
                            duration = 0,
                            isDismissible = toast.isDismissible,
                            hapticFeedback = false,
                            onDismiss = { TchatToastManager.dismiss(toast.id) }
                        )
                    }
                }
        }

        // Bottom toasts
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .align(Alignment.BottomCenter)
                .padding(horizontal = Spacing.md)
                .padding(bottom = Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs)
        ) {
            currentToasts
                .filter { it.position == TchatToastPosition.Bottom }
                .reversed() // Show newest at bottom
                .forEach { toast ->
                    key(toast.id) {
                        TchatToast(
                            isVisible = true,
                            type = toast.type,
                            message = toast.message,
                            position = toast.position,
                            style = toast.style,
                            icon = toast.icon,
                            duration = 0,
                            isDismissible = toast.isDismissible,
                            hapticFeedback = false,
                            onDismiss = { TchatToastManager.dismiss(toast.id) }
                        )
                    }
                }
        }
    }
}

/**
 * Data classes and enums
 */
data class TchatToastItem(
    val id: String = java.util.UUID.randomUUID().toString(),
    val type: TchatToastType,
    val message: String,
    val position: TchatToastPosition,
    val style: TchatToastStyle,
    val icon: ImageVector?,
    val duration: Long,
    val isDismissible: Boolean,
    val hapticFeedback: Boolean
)

enum class TchatToastType {
    Info,
    Success,
    Warning,
    Error
}

enum class TchatToastPosition {
    Top,
    Center,
    Bottom
}

enum class TchatToastStyle {
    Filled,
    Outlined,
    Minimal
}

private data class ToastColors(
    val background: Color,
    val border: Color,
    val text: Color,
    val icon: Color
)

// Preview
@Preview(showBackground = true)
@Composable
fun TchatToastPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Success toast
        TchatToast(
            isVisible = true,
            type = TchatToastType.Success,
            style = TchatToastStyle.Filled,
            message = "Changes saved successfully!",
            duration = 0
        )

        // Error toast
        TchatToast(
            isVisible = true,
            type = TchatToastType.Error,
            style = TchatToastStyle.Outlined,
            message = "Failed to upload file. Please try again.",
            duration = 0
        )

        // Warning toast
        TchatToast(
            isVisible = true,
            type = TchatToastType.Warning,
            style = TchatToastStyle.Minimal,
            message = "You have unsaved changes.",
            duration = 0
        )

        // Info toast with custom icon
        TchatToast(
            isVisible = true,
            type = TchatToastType.Info,
            style = TchatToastStyle.Filled,
            message = "New update available for download.",
            icon = Icons.Default.Download,
            duration = 0,
            isDismissible = false
        )
    }
}