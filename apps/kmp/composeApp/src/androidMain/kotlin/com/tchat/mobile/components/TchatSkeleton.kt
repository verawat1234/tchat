package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.animation.core.RepeatMode
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.composed
import androidx.compose.ui.draw.clip
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * Android implementation of TchatSkeleton using Material shimmer animation
 * Provides native Material Design loading states with smooth gradient animation
 */
@Composable
actual fun TchatSkeleton(
    modifier: Modifier,
    shape: TchatSkeletonShape,
    size: TchatSkeletonSize,
    animated: Boolean
) {
    val shimmerColors = listOf(
        TchatColors.surfaceVariant.copy(alpha = 0.9f),
        TchatColors.surfaceVariant.copy(alpha = 0.2f),
        TchatColors.surfaceVariant.copy(alpha = 0.9f)
    )

    val transition = rememberInfiniteTransition(label = "shimmer")
    val translateAnimation by transition.animateFloat(
        initialValue = 0f,
        targetValue = 1000f,
        animationSpec = infiniteRepeatable(
            animation = tween(800, easing = FastOutSlowInEasing),
            repeatMode = RepeatMode.Reverse
        ),
        label = "shimmer_translation"
    )

    val (componentShape, componentSize) = when (shape) {
        TchatSkeletonShape.Rectangle -> {
            val shapeValue = RoundedCornerShape(4.dp)
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 120.dp, height = 16.dp)
                TchatSkeletonSize.Medium -> Modifier.size(width = 200.dp, height = 20.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 300.dp, height = 24.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Circle -> {
            val shapeValue = CircleShape
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(32.dp)
                TchatSkeletonSize.Medium -> Modifier.size(48.dp)
                TchatSkeletonSize.Large -> Modifier.size(64.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Rounded -> {
            val shapeValue = RoundedCornerShape(8.dp)
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 100.dp, height = 32.dp)
                TchatSkeletonSize.Medium -> Modifier.size(width = 140.dp, height = 44.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 180.dp, height = 56.dp)
            }
            shapeValue to sizeModifier
        }

        TchatSkeletonShape.Line -> {
            val shapeValue = RoundedCornerShape(2.dp)
            val sizeModifier = when (size) {
                TchatSkeletonSize.Small -> Modifier.size(width = 80.dp, height = 12.dp)
                TchatSkeletonSize.Medium -> Modifier.size(width = 120.dp, height = 14.dp)
                TchatSkeletonSize.Large -> Modifier.size(width = 160.dp, height = 16.dp)
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
            .shimmer(shimmerColors, translateAnimation, animated)
    )
}

/**
 * Material shimmer effect modifier for Android skeleton loading
 */
private fun Modifier.shimmer(
    colors: List<Color>,
    translateAnimation: Float,
    enabled: Boolean
): Modifier = composed {
    if (enabled) {
        background(
            brush = Brush.linearGradient(
                colors = colors,
                start = Offset(translateAnimation - 200f, translateAnimation - 200f),
                end = Offset(translateAnimation, translateAnimation)
            )
        )
    } else {
        background(colors.first())
    }
}