package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import kotlinx.coroutines.delay

/**
 * Alert component following Tchat design system
 */
@Composable
fun TchatAlert(
    isVisible: Boolean,
    type: TchatAlertType,
    message: String,
    modifier: Modifier = Modifier,
    title: String? = null,
    variant: TchatAlertVariant = TchatAlertVariant.Filled,
    size: TchatAlertSize = TchatAlertSize.Medium,
    isDismissible: Boolean = true,
    showIcon: Boolean = true,
    actions: List<TchatAlertAction> = emptyList(),
    onDismiss: (() -> Unit)? = null
) {
    val hapticFeedback = LocalHapticFeedback.current

    val alertColors = getAlertColors(type, variant)
    val alertIcon = getAlertIcon(type)

    AnimatedVisibility(
        visible = isVisible,
        enter = slideInVertically(
            initialOffsetY = { -it },
            animationSpec = tween(300)
        ) + fadeIn(animationSpec = tween(300)),
        exit = slideOutVertically(
            targetOffsetY = { -it },
            animationSpec = tween(200)
        ) + fadeOut(animationSpec = tween(200))
    ) {
        Row(
            modifier = modifier
                .fillMaxWidth()
                .shadow(
                    elevation = 4.dp,
                    shape = RoundedCornerShape(8.dp)
                )
                .background(
                    color = alertColors.background,
                    shape = RoundedCornerShape(8.dp)
                )
                .then(
                    if (variant == TchatAlertVariant.Outlined) {
                        Modifier.border(
                            width = 1.dp,
                            color = alertColors.border,
                            shape = RoundedCornerShape(8.dp)
                        )
                    } else Modifier
                )
                .padding(size.padding),
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
            verticalAlignment = Alignment.Top
        ) {
            // Icon
            if (showIcon) {
                Icon(
                    imageVector = alertIcon,
                    contentDescription = null,
                    tint = alertColors.icon,
                    modifier = Modifier
                        .size(size.iconSize)
                        .padding(top = 2.dp)
                )
            }

            // Content
            Column(
                modifier = Modifier.weight(1f),
                verticalArrangement = Arrangement.spacedBy(Spacing.xs)
            ) {
                // Title
                title?.let { titleText ->
                    Text(
                        text = titleText,
                        fontSize = size.titleFontSize,
                        fontWeight = FontWeight.SemiBold,
                        color = alertColors.text
                    )
                }

                // Message
                Text(
                    text = message,
                    fontSize = size.messageFontSize,
                    color = alertColors.text
                )

                // Actions
                if (actions.isNotEmpty()) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
                        modifier = Modifier.padding(top = Spacing.xs)
                    ) {
                        actions.forEach { action ->
                            AlertActionButton(
                                action = action,
                                alertColors = alertColors,
                                variant = variant,
                                size = size,
                                onClick = {
                                    action.onClick()
                                    hapticFeedback.performHapticFeedback(
                                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                    )
                                }
                            )
                        }
                    }
                }
            }

            // Dismiss button
            if (isDismissible) {
                Icon(
                    imageVector = Icons.Default.Close,
                    contentDescription = "Dismiss",
                    tint = alertColors.text.copy(alpha = 0.7f),
                    modifier = Modifier
                        .size(16.dp)
                        .clickable {
                            onDismiss?.invoke()
                            hapticFeedback.performHapticFeedback(
                                androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                            )
                        }
                        .padding(top = 2.dp)
                )
            }
        }
    }
}

@Composable
private fun AlertActionButton(
    action: TchatAlertAction,
    alertColors: AlertColors,
    variant: TchatAlertVariant,
    size: TchatAlertSize,
    onClick: () -> Unit
) {
    val buttonColors = getActionButtonColors(action.style, alertColors, variant)

    Text(
        text = action.title,
        fontSize = (size.messageFontSize.value - 1).sp,
        fontWeight = FontWeight.Medium,
        color = buttonColors.text,
        modifier = Modifier
            .clip(RoundedCornerShape(4.dp))
            .background(buttonColors.background)
            .clickable { onClick() }
            .padding(horizontal = 12.dp, vertical = 6.dp)
    )
}

private fun getAlertColors(type: TchatAlertType, variant: TchatAlertVariant): AlertColors {
    return when (type to variant) {
        TchatAlertType.Info to TchatAlertVariant.Filled -> AlertColors(
            background = Colors.primary,
            border = Colors.primary,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatAlertType.Info to TchatAlertVariant.Outlined -> AlertColors(
            background = Colors.background,
            border = Colors.primary,
            text = Colors.primary,
            icon = Colors.primary
        )
        TchatAlertType.Info to TchatAlertVariant.Minimal -> AlertColors(
            background = Colors.primary.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.primary,
            icon = Colors.primary
        )

        TchatAlertType.Success to TchatAlertVariant.Filled -> AlertColors(
            background = Colors.success,
            border = Colors.success,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatAlertType.Success to TchatAlertVariant.Outlined -> AlertColors(
            background = Colors.background,
            border = Colors.success,
            text = Colors.success,
            icon = Colors.success
        )
        TchatAlertType.Success to TchatAlertVariant.Minimal -> AlertColors(
            background = Colors.success.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.success,
            icon = Colors.success
        )

        TchatAlertType.Warning to TchatAlertVariant.Filled -> AlertColors(
            background = Colors.warning,
            border = Colors.warning,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatAlertType.Warning to TchatAlertVariant.Outlined -> AlertColors(
            background = Colors.background,
            border = Colors.warning,
            text = Colors.warning,
            icon = Colors.warning
        )
        TchatAlertType.Warning to TchatAlertVariant.Minimal -> AlertColors(
            background = Colors.warning.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.warning,
            icon = Colors.warning
        )

        TchatAlertType.Error to TchatAlertVariant.Filled -> AlertColors(
            background = Colors.error,
            border = Colors.error,
            text = Colors.textOnPrimary,
            icon = Colors.textOnPrimary
        )
        TchatAlertType.Error to TchatAlertVariant.Outlined -> AlertColors(
            background = Colors.background,
            border = Colors.error,
            text = Colors.error,
            icon = Colors.error
        )
        TchatAlertType.Error to TchatAlertVariant.Minimal -> AlertColors(
            background = Colors.error.copy(alpha = 0.1f),
            border = Color.Transparent,
            text = Colors.error,
            icon = Colors.error
        )
        else -> AlertColors(
            background = Colors.background,
            border = Colors.border,
            text = Colors.textPrimary,
            icon = Colors.textPrimary
        )
    }
}

private fun getAlertIcon(type: TchatAlertType): ImageVector {
    return when (type) {
        TchatAlertType.Info -> Icons.Default.Info
        TchatAlertType.Success -> Icons.Default.CheckCircle
        TchatAlertType.Warning -> Icons.Default.Warning
        TchatAlertType.Error -> Icons.Default.Error
    }
}

private fun getActionButtonColors(
    style: TchatAlertActionStyle,
    alertColors: AlertColors,
    variant: TchatAlertVariant
): ActionButtonColors {
    return when (style) {
        TchatAlertActionStyle.Primary -> {
            if (variant == TchatAlertVariant.Filled && alertColors.background != Colors.background) {
                ActionButtonColors(
                    background = Colors.background.copy(alpha = 0.2f),
                    text = Colors.textOnPrimary
                )
            } else {
                ActionButtonColors(
                    background = alertColors.text.copy(alpha = 0.1f),
                    text = alertColors.text
                )
            }
        }
        TchatAlertActionStyle.Secondary -> ActionButtonColors(
            background = alertColors.text.copy(alpha = 0.05f),
            text = alertColors.text.copy(alpha = 0.8f)
        )
        TchatAlertActionStyle.Destructive -> ActionButtonColors(
            background = Colors.error.copy(alpha = 0.1f),
            text = Colors.error
        )
    }
}

/**
 * Alert manager for global alert handling
 */
object TchatAlertManager {
    private var _currentAlert = mutableStateOf<TchatAlertItem?>(null)
    val currentAlert: State<TchatAlertItem?> = _currentAlert

    fun show(
        type: TchatAlertType,
        message: String,
        title: String? = null,
        variant: TchatAlertVariant = TchatAlertVariant.Filled,
        actions: List<TchatAlertAction> = emptyList(),
        duration: Long? = null
    ) {
        val alert = TchatAlertItem(
            type = type,
            title = title,
            message = message,
            variant = variant,
            actions = actions
        )

        _currentAlert.value = alert

        duration?.let { durationMs ->
            // Auto-dismiss after duration
            // Note: In a real app, you'd use a coroutine scope tied to the application lifecycle
        }
    }

    fun dismiss() {
        _currentAlert.value = null
    }

    fun showSuccess(
        message: String,
        title: String? = null,
        duration: Long = 3000L
    ) {
        show(
            type = TchatAlertType.Success,
            title = title,
            message = message,
            duration = duration
        )
    }

    fun showError(
        message: String,
        title: String? = null,
        actions: List<TchatAlertAction> = emptyList()
    ) {
        show(
            type = TchatAlertType.Error,
            title = title,
            message = message,
            actions = actions
        )
    }

    fun showWarning(
        message: String,
        title: String? = null,
        actions: List<TchatAlertAction> = emptyList()
    ) {
        show(
            type = TchatAlertType.Warning,
            title = title,
            message = message,
            actions = actions
        )
    }

    fun showInfo(
        message: String,
        title: String? = null,
        duration: Long = 5000L
    ) {
        show(
            type = TchatAlertType.Info,
            title = title,
            message = message,
            duration = duration
        )
    }
}

/**
 * Alert overlay composable for global alerts
 */
@Composable
fun TchatAlertOverlay() {
    val currentAlert by TchatAlertManager.currentAlert

    Column(
        modifier = Modifier.fillMaxWidth()
    ) {
        currentAlert?.let { alert ->
            TchatAlert(
                isVisible = true,
                type = alert.type,
                title = alert.title,
                message = alert.message,
                variant = alert.variant,
                actions = alert.actions,
                onDismiss = { TchatAlertManager.dismiss() },
                modifier = Modifier.padding(horizontal = Spacing.md, vertical = Spacing.md)
            )
        }
    }
}

/**
 * Data classes and enums
 */
data class TchatAlertItem(
    val type: TchatAlertType,
    val title: String?,
    val message: String,
    val variant: TchatAlertVariant,
    val actions: List<TchatAlertAction>
)

data class TchatAlertAction(
    val title: String,
    val style: TchatAlertActionStyle = TchatAlertActionStyle.Primary,
    val onClick: () -> Unit
)

enum class TchatAlertType {
    Info,
    Success,
    Warning,
    Error
}

enum class TchatAlertVariant {
    Filled,
    Outlined,
    Minimal
}

enum class TchatAlertActionStyle {
    Primary,
    Secondary,
    Destructive
}

enum class TchatAlertSize(
    val padding: androidx.compose.ui.unit.Dp,
    val iconSize: androidx.compose.ui.unit.Dp,
    val titleFontSize: androidx.compose.ui.unit.TextUnit,
    val messageFontSize: androidx.compose.ui.unit.TextUnit
) {
    Small(
        padding = 12.dp,
        iconSize = 16.dp,
        titleFontSize = 14.sp,
        messageFontSize = 12.sp
    ),
    Medium(
        padding = 16.dp,
        iconSize = 20.dp,
        titleFontSize = 16.sp,
        messageFontSize = 14.sp
    ),
    Large(
        padding = 20.dp,
        iconSize = 24.dp,
        titleFontSize = 18.sp,
        messageFontSize = 16.sp
    )
}

private data class AlertColors(
    val background: Color,
    val border: Color,
    val text: Color,
    val icon: Color
)

private data class ActionButtonColors(
    val background: Color,
    val text: Color
)

// Preview
@Preview(showBackground = true)
@Composable
fun TchatAlertPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Success alert
        TchatAlert(
            isVisible = true,
            type = TchatAlertType.Success,
            variant = TchatAlertVariant.Filled,
            title = "Success",
            message = "Your changes have been saved successfully."
        )

        // Error alert with actions
        TchatAlert(
            isVisible = true,
            type = TchatAlertType.Error,
            variant = TchatAlertVariant.Outlined,
            title = "Error",
            message = "Failed to save changes. Please try again.",
            actions = listOf(
                TchatAlertAction("Retry", TchatAlertActionStyle.Primary) { },
                TchatAlertAction("Cancel", TchatAlertActionStyle.Secondary) { }
            )
        )

        // Warning alert
        TchatAlert(
            isVisible = true,
            type = TchatAlertType.Warning,
            variant = TchatAlertVariant.Minimal,
            title = "Warning",
            message = "This action cannot be undone. Are you sure you want to continue?",
            actions = listOf(
                TchatAlertAction("Continue", TchatAlertActionStyle.Destructive) { },
                TchatAlertAction("Cancel", TchatAlertActionStyle.Secondary) { }
            )
        )

        // Info alert
        TchatAlert(
            isVisible = true,
            type = TchatAlertType.Info,
            variant = TchatAlertVariant.Filled,
            message = "New features are available. Update your app to get the latest improvements.",
            size = TchatAlertSize.Small
        )

        // Minimal alert without icon
        TchatAlert(
            isVisible = true,
            type = TchatAlertType.Info,
            variant = TchatAlertVariant.Minimal,
            message = "This is a minimal alert without an icon.",
            showIcon = false,
            isDismissible = false
        )
    }
}