package com.tchat.mobile.components

import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

/**
 * Cross-platform Dialog implementation
 * Uses expect/actual pattern for platform-specific behavior
 */
@Composable
expect fun PlatformDialog(
    visible: Boolean,
    onDismiss: () -> Unit,
    content: @Composable () -> Unit
)

/**
 * Cross-platform Sheet implementation
 * Handles iOS/Android differences internally
 */
@Composable
expect fun PlatformSheet(
    visible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
)

/**
 * Cross-platform Configuration access
 * Replaces LocalConfiguration which is Android-only
 */
expect object PlatformConfiguration {
    val screenWidth: Int
        @Composable get
    val screenHeight: Int
        @Composable get
    val density: Float
        @Composable get
}

/**
 * Common TchatDialog with platform abstractions
 */
@Composable
fun TchatDialog(
    visible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    PlatformDialog(
        visible = visible,
        onDismiss = onDismiss,
        content = content
    )
}

/**
 * Common TchatSheet with platform abstractions
 */
@Composable
fun TchatSheet(
    visible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    PlatformSheet(
        visible = visible,
        onDismiss = onDismiss,
        modifier = modifier,
        content = content
    )
}