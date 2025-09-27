package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.spring
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import platform.UIKit.UIScreen
import kotlinx.cinterop.ExperimentalForeignApi

/**
 * iOS implementation of PlatformDialog
 * Uses custom overlay instead of problematic Dialog composable
 */
@Composable
actual fun PlatformDialog(
    visible: Boolean,
    onDismiss: () -> Unit,
    content: @Composable () -> Unit
) {
    if (visible) {
        // iOS-compatible dialog implementation
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(Color.Black.copy(alpha = 0.5f))
                .clickable(
                    interactionSource = remember { MutableInteractionSource() },
                    indication = null,
                    onClick = onDismiss
                ),
            contentAlignment = Alignment.Center
        ) {
            Surface(
                modifier = Modifier
                    .padding(16.dp)
                    .wrapContentHeight()
                    .clickable(enabled = false) { }, // Prevent click-through
                shape = RoundedCornerShape(12.dp),
                shadowElevation = 8.dp,
                color = MaterialTheme.colorScheme.surface
            ) {
                Box(
                    modifier = Modifier.padding(24.dp)
                ) {
                    content()
                }
            }
        }
    }
}

/**
 * iOS implementation of PlatformSheet
 * Native iOS-style bottom sheet with proper animations
 */
@Composable
actual fun PlatformSheet(
    visible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier,
    content: @Composable () -> Unit
) {
    if (visible) {
        Box(
            modifier = Modifier.fillMaxSize()
        ) {
            // Backdrop
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.3f))
                    .clickable(
                        interactionSource = remember { MutableInteractionSource() },
                        indication = null,
                        onClick = onDismiss
                    )
            )

            // Sheet Content
            AnimatedVisibility(
                visible = visible,
                modifier = Modifier.align(Alignment.BottomCenter),
                enter = slideInVertically(
                    animationSpec = spring(
                        dampingRatio = 0.8f,
                        stiffness = 400f
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
                        .wrapContentHeight()
                        .shadow(16.dp, RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp))
                        .clip(RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp))
                        .semantics {
                            role = Role.Button
                        },
                    color = MaterialTheme.colorScheme.surface,
                    shape = RoundedCornerShape(topStart = 20.dp, topEnd = 20.dp)
                ) {
                    Column {
                        // iOS-style drag indicator
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(vertical = 12.dp),
                            contentAlignment = Alignment.Center
                        ) {
                            Surface(
                                color = Color.Gray.copy(alpha = 0.3f),
                                shape = RoundedCornerShape(2.5.dp)
                            ) {
                                Spacer(
                                    modifier = Modifier.size(width = 36.dp, height = 5.dp)
                                )
                            }
                        }

                        // Content
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(horizontal = 16.dp, vertical = 8.dp)
                        ) {
                            content()
                        }

                        // Bottom padding for safe area
                        Spacer(modifier = Modifier.height(16.dp))
                    }
                }
            }
        }
    }
}

/**
 * iOS implementation of PlatformConfiguration
 * Uses UIScreen for iOS-specific screen information
 */
actual object PlatformConfiguration {
    @OptIn(ExperimentalForeignApi::class)
    actual val screenWidth: Int
        @Composable get() {
            val density = LocalDensity.current
            val screenBounds = UIScreen.mainScreen.bounds
            return with(density) { screenBounds.size.width.toFloat().dp.roundToPx() }
        }

    @OptIn(ExperimentalForeignApi::class)
    actual val screenHeight: Int
        @Composable get() {
            val density = LocalDensity.current
            val screenBounds = UIScreen.mainScreen.bounds
            return with(density) { screenBounds.size.height.toFloat().dp.roundToPx() }
        }

    actual val density: Float
        @Composable get() = LocalDensity.current.density
}