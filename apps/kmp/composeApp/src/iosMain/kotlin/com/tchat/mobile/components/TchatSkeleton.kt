package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.composed
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * iOS implementation of TchatSkeleton with SwiftUI-inspired pulse animation
 * Uses iOS HIG loading patterns with subtle pulse animation and appropriate opacity
 */
@Composable
actual fun TchatSkeleton(
    modifier: Modifier,
    shape: TchatSkeletonShape,
    size: TchatSkeletonSize,
    animated: Boolean
) {
    // iOS uses more subtle shimmer colors with higher opacity
    val shimmerColors = listOf(
        TchatColors.surfaceVariant.copy(alpha = 0.8f), // iOS uses less transparent base
        TchatColors.surfaceVariant.copy(alpha = 0.3f), // Subtle highlight
        TchatColors.surfaceVariant.copy(alpha = 0.8f)
    )

    val transition = rememberInfiniteTransition(label = "ios_pulse")
    val translateAnimation by transition.animateFloat(
        initialValue = 0f,
        targetValue = 1000f,
        animationSpec = infiniteRepeatable(
            animation = tween(
                1200, // iOS uses slower, more elegant animation
                easing = CubicBezierEasing(0.4f, 0.0f, 0.6f, 1.0f) // iOS-style spring easing
            ),
            repeatMode = RepeatMode.Restart
        ),
        label = "ios_pulse_translation"
    )

    val (componentShape, componentSize) = when (shape) {
        TchatSkeletonShape.Rectangle -> {
            val shapeValue = RoundedCornerShape(6.dp) // iOS uses slightly more rounded corners
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 120.dp, height = 17.dp) // iOS uses slightly taller text
                TchatSkeletonSize.Medium -> Modifier.size(width = 200.dp, height = 21.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 300.dp, height = 25.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Circle -> {
            val shapeValue = CircleShape
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(34.dp) // iOS uses slightly larger touch targets
                TchatSkeletonSize.Medium -> Modifier.size(50.dp)
                TchatSkeletonSize.Large -> Modifier.size(66.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Rounded -> {
            val shapeValue = RoundedCornerShape(12.dp) // iOS uses more pronounced rounded corners
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 100.dp, height = 34.dp) // iOS button heights
                TchatSkeletonSize.Medium -> Modifier.size(width = 140.dp, height = 46.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 180.dp, height = 58.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Line -> {
            val shapeValue = RoundedCornerShape(3.dp) // iOS uses slightly more rounded line caps
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 80.dp, height = 13.dp)
                TchatSkeletonSize.Medium -> Modifier.size(width = 120.dp, height = 15.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 160.dp, height = 17.dp)
            }
            shapeValue to sizeModifier
        }
    }

    Box(
        modifier = modifier
            .then(componentSize)
            .clip(componentShape)
            .semantics {
                role = Role.Image
                contentDescription = "Loading content"
            }
            .iosPulse(shimmerColors, translateAnimation, animated)
    )
}

/**
 * iOS-style pulse effect modifier with elegant gradient animation
 */
private fun Modifier.iosPulse(
    colors: List<Color>,
    translateAnimation: Float,
    enabled: Boolean
): Modifier = composed {
    if (enabled) {
        background(
            brush = Brush.linearGradient(
                colors = colors,
                start = Offset(translateAnimation - 300f, translateAnimation - 300f), // iOS uses wider gradient sweep
                end = Offset(translateAnimation, translateAnimation)
            )
        )
    } else {
        background(colors.first())
    }
}