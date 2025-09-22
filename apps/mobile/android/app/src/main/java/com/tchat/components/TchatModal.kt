package com.tchat.components

import androidx.compose.animation.*
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
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
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import androidx.compose.ui.zIndex
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import kotlin.math.roundToInt

/**
 * Modal component following Tchat design system
 */
@Composable
fun TchatModal(
    isVisible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier,
    size: TchatModalSize = TchatModalSize.Medium,
    position: TchatModalPosition = TchatModalPosition.Center,
    animation: TchatModalAnimation = TchatModalAnimation.Slide,
    showCloseButton: Boolean = true,
    isDismissible: Boolean = true,
    allowDragDismiss: Boolean = true,
    showOverlay: Boolean = true,
    overlayColor: Color = Color.Black.copy(alpha = 0.5f),
    cornerRadius: androidx.compose.ui.unit.Dp = 16.dp,
    content: @Composable () -> Unit
) {
    val hapticFeedback = LocalHapticFeedback.current
    val density = LocalDensity.current

    var dragOffset by remember { mutableStateOf(0f) }
    val dismissThreshold = with(density) { 150.dp.toPx() }

    if (isVisible) {
        Dialog(
            onDismissRequest = {
                if (isDismissible) {
                    onDismiss()
                }
            },
            properties = DialogProperties(
                dismissOnBackPress = isDismissible,
                dismissOnClickOutside = isDismissible,
                usePlatformDefaultWidth = false
            )
        ) {
            BoxWithConstraints(
                modifier = Modifier.fillMaxSize()
            ) {
                val modalWidth = getModalWidth(size, maxWidth)
                val modalHeight = getModalHeight(size, maxHeight)

                // Overlay
                if (showOverlay) {
                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .background(overlayColor)
                            .clickable(enabled = isDismissible) {
                                hapticFeedback.performHapticFeedback(
                                    androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                )
                                onDismiss()
                            }
                    )
                }

                // Modal content
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .then(
                            when (position) {
                                TchatModalPosition.Center -> Modifier.wrapContentSize(Alignment.Center)
                                TchatModalPosition.Bottom -> Modifier.wrapContentSize(Alignment.BottomCenter)
                                TchatModalPosition.Top -> Modifier.wrapContentSize(Alignment.TopCenter)
                            }
                        )
                        .padding(
                            horizontal = if (size == TchatModalSize.FullScreen) 0.dp else Spacing.lg,
                            vertical = if (size == TchatModalSize.FullScreen) 0.dp else Spacing.lg
                        )
                ) {
                    Column(
                        modifier = Modifier
                            .width(modalWidth)
                            .then(
                                if (modalHeight != null) {
                                    Modifier.height(modalHeight)
                                } else {
                                    Modifier.wrapContentHeight()
                                }
                            )
                            .offset { IntOffset(0, dragOffset.roundToInt()) }
                            .shadow(
                                elevation = if (size == TchatModalSize.FullScreen) 0.dp else 20.dp,
                                shape = RoundedCornerShape(cornerRadius)
                            )
                            .background(
                                color = Colors.background,
                                shape = RoundedCornerShape(cornerRadius)
                            )
                            .then(
                                if (allowDragDismiss && isDismissible) {
                                    Modifier.pointerInput(Unit) {
                                        detectDragGestures(
                                            onDragEnd = {
                                                val shouldDismiss = when (position) {
                                                    TchatModalPosition.Center -> kotlin.math.abs(dragOffset) > dismissThreshold
                                                    TchatModalPosition.Bottom -> dragOffset > dismissThreshold
                                                    TchatModalPosition.Top -> dragOffset < -dismissThreshold
                                                }

                                                if (shouldDismiss) {
                                                    hapticFeedback.performHapticFeedback(
                                                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                                    )
                                                    onDismiss()
                                                } else {
                                                    dragOffset = 0f
                                                }
                                            }
                                        ) { _, dragAmount ->
                                            dragOffset += when (position) {
                                                TchatModalPosition.Center -> dragAmount.y
                                                TchatModalPosition.Bottom -> maxOf(0f, dragAmount.y)
                                                TchatModalPosition.Top -> minOf(0f, dragAmount.y)
                                            }
                                        }
                                    }
                                } else Modifier
                            )
                    ) {
                        // Close button
                        if (showCloseButton) {
                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(
                                        top = Spacing.md,
                                        end = Spacing.md
                                    ),
                                horizontalArrangement = Arrangement.End
                            ) {
                                Icon(
                                    imageVector = Icons.Default.Close,
                                    contentDescription = "Close",
                                    tint = Colors.textSecondary,
                                    modifier = Modifier
                                        .size(32.dp)
                                        .clip(CircleShape)
                                        .background(Colors.surface)
                                        .clickable {
                                            hapticFeedback.performHapticFeedback(
                                                androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                            )
                                            onDismiss()
                                        }
                                        .padding(8.dp)
                                        .zIndex(1f)
                                )
                            }
                        }

                        // Content
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .then(
                                    if (modalHeight != null) {
                                        Modifier.weight(1f)
                                    } else {
                                        Modifier.wrapContentHeight()
                                    }
                                )
                        ) {
                            content()
                        }
                    }
                }
            }
        }
    }
}

private fun getModalWidth(
    size: TchatModalSize,
    maxWidth: androidx.compose.ui.unit.Dp
): androidx.compose.ui.unit.Dp {
    return when (size) {
        TchatModalSize.Small -> minOf(300.dp, maxWidth - 32.dp)
        TchatModalSize.Medium -> minOf(400.dp, maxWidth - 32.dp)
        TchatModalSize.Large -> minOf(500.dp, maxWidth - 32.dp)
        TchatModalSize.FullScreen -> maxWidth
        is TchatModalSize.Custom -> size.width?.let { minOf(it, maxWidth - 32.dp) } ?: (maxWidth - 32.dp)
    }
}

private fun getModalHeight(
    size: TchatModalSize,
    maxHeight: androidx.compose.ui.unit.Dp
): androidx.compose.ui.unit.Dp? {
    return when (size) {
        TchatModalSize.Small -> minOf(200.dp, maxHeight - 32.dp)
        TchatModalSize.Medium -> minOf(300.dp, maxHeight - 32.dp)
        TchatModalSize.Large -> minOf(400.dp, maxHeight - 32.dp)
        TchatModalSize.FullScreen -> maxHeight
        is TchatModalSize.Custom -> size.height?.let { minOf(it, maxHeight - 32.dp) }
    }
}

/**
 * Modal manager for global modal handling
 */
object TchatModalManager {
    private var _currentModal = mutableStateOf<TchatModalItem?>(null)
    val currentModal: State<TchatModalItem?> = _currentModal

    fun present(
        size: TchatModalSize = TchatModalSize.Medium,
        position: TchatModalPosition = TchatModalPosition.Center,
        animation: TchatModalAnimation = TchatModalAnimation.Slide,
        showCloseButton: Boolean = true,
        isDismissible: Boolean = true,
        allowDragDismiss: Boolean = true,
        showOverlay: Boolean = true,
        onDismiss: (() -> Unit)? = null,
        content: @Composable () -> Unit
    ) {
        val modal = TchatModalItem(
            size = size,
            position = position,
            animation = animation,
            showCloseButton = showCloseButton,
            isDismissible = isDismissible,
            allowDragDismiss = allowDragDismiss,
            showOverlay = showOverlay,
            content = content,
            onDismiss = onDismiss
        )

        _currentModal.value = modal
    }

    fun dismiss() {
        _currentModal.value = null
    }

    fun presentAlert(
        title: String,
        message: String,
        primaryButton: String = "OK",
        secondaryButton: String? = null,
        onPrimary: (() -> Unit)? = null,
        onSecondary: (() -> Unit)? = null
    ) {
        present(
            size = TchatModalSize.Small,
            position = TchatModalPosition.Center
        ) {
            AlertContent(
                title = title,
                message = message,
                primaryButton = primaryButton,
                secondaryButton = secondaryButton,
                onPrimary = {
                    onPrimary?.invoke()
                    dismiss()
                },
                onSecondary = {
                    onSecondary?.invoke()
                    dismiss()
                }
            )
        }
    }

    fun presentBottomSheet(
        content: @Composable () -> Unit
    ) {
        present(
            size = TchatModalSize.Custom(width = null, height = 400.dp),
            position = TchatModalPosition.Bottom,
            animation = TchatModalAnimation.Slide,
            showCloseButton = false,
            allowDragDismiss = true,
            content = content
        )
    }
}

@Composable
private fun AlertContent(
    title: String,
    message: String,
    primaryButton: String,
    secondaryButton: String?,
    onPrimary: () -> Unit,
    onSecondary: () -> Unit
) {
    Column(
        modifier = Modifier.padding(Spacing.lg),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Column(
            verticalArrangement = Arrangement.spacedBy(Spacing.sm),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                text = title,
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold,
                color = Colors.textPrimary,
                textAlign = TextAlign.Center
            )

            Text(
                text = message,
                fontSize = 14.sp,
                color = Colors.textSecondary,
                textAlign = TextAlign.Center
            )
        }

        Row(
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            secondaryButton?.let { buttonText ->
                TextButton(
                    onClick = onSecondary,
                    modifier = Modifier.weight(1f)
                ) {
                    Text(
                        text = buttonText,
                        color = Colors.textSecondary
                    )
                }
            }

            Button(
                onClick = onPrimary,
                modifier = Modifier.weight(1f),
                colors = ButtonDefaults.buttonColors(
                    containerColor = Colors.primary
                )
            ) {
                Text(
                    text = primaryButton,
                    color = Colors.textOnPrimary,
                    fontWeight = FontWeight.SemiBold
                )
            }
        }
    }
}

/**
 * Modal overlay composable for global modals
 */
@Composable
fun TchatModalOverlay() {
    val currentModal by TchatModalManager.currentModal

    currentModal?.let { modal ->
        TchatModal(
            isVisible = true,
            size = modal.size,
            position = modal.position,
            animation = modal.animation,
            showCloseButton = modal.showCloseButton,
            isDismissible = modal.isDismissible,
            allowDragDismiss = modal.allowDragDismiss,
            showOverlay = modal.showOverlay,
            onDismiss = {
                modal.onDismiss?.invoke()
                TchatModalManager.dismiss()
            },
            content = modal.content
        )
    }
}

/**
 * Data classes and enums
 */
data class TchatModalItem(
    val size: TchatModalSize,
    val position: TchatModalPosition,
    val animation: TchatModalAnimation,
    val showCloseButton: Boolean,
    val isDismissible: Boolean,
    val allowDragDismiss: Boolean,
    val showOverlay: Boolean,
    val content: @Composable () -> Unit,
    val onDismiss: (() -> Unit)?
)

sealed class TchatModalSize {
    object Small : TchatModalSize()
    object Medium : TchatModalSize()
    object Large : TchatModalSize()
    object FullScreen : TchatModalSize()
    data class Custom(
        val width: androidx.compose.ui.unit.Dp?,
        val height: androidx.compose.ui.unit.Dp?
    ) : TchatModalSize()
}

enum class TchatModalPosition {
    Center,
    Bottom,
    Top
}

enum class TchatModalAnimation {
    Slide,
    Fade,
    Scale
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatModalPreview() {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Colors.surface)
    ) {
        // Background content
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(Spacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Text(
                text = "Background Content",
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.textPrimary
            )
        }

        // Modal preview
        TchatModal(
            isVisible = true,
            size = TchatModalSize.Medium,
            position = TchatModalPosition.Center,
            onDismiss = { }
        ) {
            Column(
                modifier = Modifier.padding(Spacing.lg),
                verticalArrangement = Arrangement.spacedBy(Spacing.lg),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Icon(
                    imageVector = Icons.Default.CheckCircle,
                    contentDescription = null,
                    tint = Colors.success,
                    modifier = Modifier.size(48.dp)
                )

                Column(
                    verticalArrangement = Arrangement.spacedBy(Spacing.sm),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        text = "Success!",
                        fontSize = 20.sp,
                        fontWeight = FontWeight.Bold,
                        color = Colors.textPrimary
                    )

                    Text(
                        text = "Your changes have been saved successfully.",
                        fontSize = 14.sp,
                        color = Colors.textSecondary,
                        textAlign = TextAlign.Center
                    )
                }

                Row(
                    horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
                ) {
                    TextButton(onClick = { }) {
                        Text(
                            text = "Cancel",
                            color = Colors.textSecondary
                        )
                    }

                    Button(
                        onClick = { },
                        colors = ButtonDefaults.buttonColors(
                            containerColor = Colors.primary
                        )
                    ) {
                        Text(
                            text = "Continue",
                            color = Colors.textOnPrimary,
                            fontWeight = FontWeight.SemiBold
                        )
                    }
                }
            }
        }
    }
}