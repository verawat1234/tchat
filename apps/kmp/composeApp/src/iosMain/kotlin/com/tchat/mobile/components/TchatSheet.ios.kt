package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.Spring
import androidx.compose.animation.core.spring
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import androidx.compose.ui.zIndex
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import kotlin.math.abs
import kotlin.math.roundToInt

/**
 * iOS implementation of TchatSheet with native-style sheet presentations
 * Mimics iOS sheet behavior with appropriate animations and gestures
 */
@Composable
actual fun TchatSheet(
    isVisible: Boolean,
    onDismissRequest: () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier,
    mode: SheetPresentationMode,
    sizing: SheetSizing,
    sizeValue: Float,
    dismissBehavior: SheetDismissBehavior,
    animation: SheetAnimation,
    shape: Shape,
    backgroundColor: Color,
    backdropColor: Color,
    elevation: Dp,
    dragIndicatorColor: Color,
    contentPadding: PaddingValues,
    windowInsets: WindowInsets,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    when (mode) {
        SheetPresentationMode.Modal -> {
            IOSModalSheet(
                isVisible = isVisible,
                onDismissRequest = onDismissRequest,
                content = content,
                modifier = modifier,
                sizing = sizing,
                sizeValue = sizeValue,
                dismissBehavior = dismissBehavior,
                animation = animation,
                shape = shape,
                backgroundColor = backgroundColor,
                backdropColor = backdropColor,
                elevation = elevation,
                dragIndicatorColor = dragIndicatorColor,
                contentPadding = contentPadding,
                interactionSource = interactionSource,
                contentDescription = contentDescription
            )
        }

        SheetPresentationMode.Persistent -> {
            IOSPersistentSheet(
                isVisible = isVisible,
                onDismissRequest = onDismissRequest,
                content = content,
                modifier = modifier,
                sizing = sizing,
                sizeValue = sizeValue,
                dismissBehavior = dismissBehavior,
                shape = shape,
                backgroundColor = backgroundColor,
                elevation = elevation,
                dragIndicatorColor = dragIndicatorColor,
                contentPadding = contentPadding,
                interactionSource = interactionSource,
                contentDescription = contentDescription
            )
        }

        SheetPresentationMode.Fullscreen -> {
            IOSFullscreenSheet(
                isVisible = isVisible,
                onDismissRequest = onDismissRequest,
                content = content,
                modifier = modifier,
                dismissBehavior = dismissBehavior,
                animation = animation,
                backgroundColor = backgroundColor,
                contentPadding = contentPadding,
                contentDescription = contentDescription
            )
        }
    }
}

/**
 * iOS Modal Sheet implementation with native-style animations
 */
@Composable
private fun IOSModalSheet(
    isVisible: Boolean,
    onDismissRequest: () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier,
    sizing: SheetSizing,
    sizeValue: Float,
    dismissBehavior: SheetDismissBehavior,
    animation: SheetAnimation,
    shape: Shape,
    backgroundColor: Color,
    backdropColor: Color,
    elevation: Dp,
    dragIndicatorColor: Color,
    contentPadding: PaddingValues,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val configuration = LocalConfiguration.current
    val density = LocalDensity.current
    var dragOffset by remember { mutableStateOf(0f) }

    // Handle dismiss confirmation
    val handleDismiss: () -> Unit = {
        if (dismissBehavior.confirmDismiss) {
            // Show confirmation (simplified for this implementation)
            onDismissRequest()
        } else {
            onDismissRequest()
        }
    }

    if (isVisible) {
        Dialog(
            onDismissRequest = if (dismissBehavior.backdropDismiss) handleDismiss else {},
            properties = DialogProperties(
                dismissOnBackPress = dismissBehavior.backdropDismiss,
                dismissOnClickOutside = dismissBehavior.backdropDismiss,
                usePlatformDefaultWidth = false
            )
        ) {
            Box(
                modifier = Modifier.fillMaxSize()
            ) {
                // Backdrop
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .background(backdropColor)
                        .clickable(
                            enabled = dismissBehavior.backdropDismiss,
                            interactionSource = remember { MutableInteractionSource() },
                            indication = null,
                            onClick = handleDismiss
                        )
                )

                // Sheet Content
                AnimatedVisibility(
                    visible = isVisible,
                    modifier = Modifier.align(Alignment.BottomCenter),
                    enter = slideInVertically(
                        animationSpec = spring(
                            dampingRatio = animation.springDamping,
                            stiffness = animation.springStiffness
                        ),
                        initialOffsetY = { it }
                    ),
                    exit = slideOutVertically(
                        animationSpec = tween(animation.dismissalDuration),
                        targetOffsetY = { it }
                    )
                ) {
                    Surface(
                        modifier = modifier
                            .fillMaxWidth()
                            .then(
                                when (sizing) {
                                    SheetSizing.FitContent -> Modifier.wrapContentHeight()
                                    SheetSizing.Fractional -> Modifier.height(
                                        with(density) {
                                            (configuration.screenHeightDp * sizeValue).dp
                                        }
                                    )
                                    SheetSizing.Fixed -> Modifier.height(sizeValue.dp)
                                    SheetSizing.Expanded -> Modifier.fillMaxSize()
                                }
                            )
                            .offset {
                                IntOffset(0, dragOffset.roundToInt())
                            }
                            .then(
                                if (dismissBehavior.dragToDismiss) {
                                    Modifier.pointerInput(Unit) {
                                        detectDragGestures(
                                            onDragEnd = {
                                                if (dragOffset > size.height * dismissBehavior.dismissThreshold) {
                                                    handleDismiss()
                                                }
                                                dragOffset = 0f
                                            },
                                            onDrag = { change ->
                                                val newOffset = dragOffset + change.y
                                                if (newOffset >= 0) {
                                                    dragOffset = newOffset
                                                }
                                            }
                                        )
                                    }
                                } else {
                                    Modifier
                                }
                            )
                            .shadow(elevation, shape)
                            .clip(shape)
                            .semantics {
                                contentDescription?.let {
                                    this.contentDescription = it
                                }
                                role = Role.Dialog
                            },
                        color = backgroundColor,
                        shape = shape
                    ) {
                        Column {
                            // iOS-style drag indicator
                            if (dismissBehavior.dragToDismiss) {
                                Box(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .padding(vertical = TchatSpacing.sm),
                                    contentAlignment = Alignment.Center
                                ) {
                                    Surface(
                                        color = dragIndicatorColor,
                                        shape = RoundedCornerShape(2.5.dp)
                                    ) {
                                        Spacer(
                                            modifier = Modifier.size(width = 36.dp, height = 5.dp)
                                        )
                                    }
                                }
                            }

                            // Content
                            Box(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(contentPadding)
                            ) {
                                content()
                            }
                        }
                    }
                }
            }
        }
    }
}

/**
 * iOS Persistent Sheet implementation (non-modal)
 */
@Composable
private fun IOSPersistentSheet(
    isVisible: Boolean,
    onDismissRequest: () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier,
    sizing: SheetSizing,
    sizeValue: Float,
    dismissBehavior: SheetDismissBehavior,
    shape: Shape,
    backgroundColor: Color,
    elevation: Dp,
    dragIndicatorColor: Color,
    contentPadding: PaddingValues,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val configuration = LocalConfiguration.current
    val density = LocalDensity.current
    var dragOffset by remember { mutableStateOf(0f) }

    AnimatedVisibility(
        visible = isVisible,
        enter = slideInVertically(
            animationSpec = spring(
                dampingRatio = Spring.DampingRatioMediumBouncy,
                stiffness = Spring.StiffnessMedium
            ),
            initialOffsetY = { it }
        ),
        exit = slideOutVertically(
            animationSpec = tween(250),
            targetOffsetY = { it }
        )
    ) {
        Surface(
            modifier = modifier
                .fillMaxWidth()
                .then(
                    when (sizing) {
                        SheetSizing.FitContent -> Modifier.wrapContentHeight()
                        SheetSizing.Fractional -> Modifier.height(
                            with(density) {
                                (configuration.screenHeightDp * sizeValue).dp
                            }
                        )
                        SheetSizing.Fixed -> Modifier.height(sizeValue.dp)
                        SheetSizing.Expanded -> Modifier.fillMaxSize()
                    }
                )
                .offset {
                    IntOffset(0, dragOffset.roundToInt())
                }
                .then(
                    if (dismissBehavior.dragToDismiss) {
                        Modifier.pointerInput(Unit) {
                            detectDragGestures(
                                onDragEnd = {
                                    if (dragOffset > size.height * dismissBehavior.dismissThreshold) {
                                        onDismissRequest()
                                    }
                                    dragOffset = 0f
                                },
                                onDrag = { change ->
                                    val newOffset = dragOffset + change.y
                                    if (newOffset >= 0) {
                                        dragOffset = newOffset
                                    }
                                }
                            )
                        }
                    } else {
                        Modifier
                    }
                )
                .shadow(elevation, shape)
                .clip(shape)
                .semantics {
                    contentDescription?.let {
                        this.contentDescription = it
                    }
                    role = Role.Generic
                },
            color = backgroundColor,
            shape = shape
        ) {
            Column {
                // iOS-style drag indicator
                if (dismissBehavior.dragToDismiss) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(vertical = TchatSpacing.sm),
                        contentAlignment = Alignment.Center
                    ) {
                        Surface(
                            color = dragIndicatorColor,
                            shape = RoundedCornerShape(2.5.dp)
                        ) {
                            Spacer(
                                modifier = Modifier.size(width = 36.dp, height = 5.dp)
                            )
                        }
                    }
                }

                // Content
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(contentPadding)
                ) {
                    content()
                }
            }
        }
    }
}

/**
 * iOS Fullscreen Sheet implementation
 */
@Composable
private fun IOSFullscreenSheet(
    isVisible: Boolean,
    onDismissRequest: () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier,
    dismissBehavior: SheetDismissBehavior,
    animation: SheetAnimation,
    backgroundColor: Color,
    contentPadding: PaddingValues,
    contentDescription: String?
) {
    if (isVisible) {
        Dialog(
            onDismissRequest = if (dismissBehavior.backdropDismiss) onDismissRequest else {},
            properties = DialogProperties(
                dismissOnBackPress = dismissBehavior.backdropDismiss,
                dismissOnClickOutside = false,
                usePlatformDefaultWidth = false
            )
        ) {
            AnimatedVisibility(
                visible = isVisible,
                enter = slideInVertically(
                    animationSpec = spring(
                        dampingRatio = animation.springDamping,
                        stiffness = animation.springStiffness
                    ),
                    initialOffsetY = { it }
                ),
                exit = slideOutVertically(
                    animationSpec = tween(animation.dismissalDuration),
                    targetOffsetY = { it }
                )
            ) {
                Surface(
                    modifier = modifier
                        .fillMaxSize()
                        .semantics {
                            contentDescription?.let {
                                this.contentDescription = it
                            }
                            role = Role.Dialog
                        },
                    color = backgroundColor
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(contentPadding)
                    ) {
                        content()
                    }
                }
            }
        }
    }
}