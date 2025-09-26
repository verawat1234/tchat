package com.tchat.mobile.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.spring
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * iOS implementation of TchatCard with SwiftUI-inspired styling
 * Uses spring animations and iOS HIG-compliant interaction patterns
 */
@Composable
actual fun TchatCard(
    modifier: Modifier,
    variant: TchatCardVariant,
    size: TchatCardSize,
    onClick: (() -> Unit)?,
    enabled: Boolean,
    content: @Composable () -> Unit
) {
    val interactionSource = remember { MutableInteractionSource() }
    val isPressed by interactionSource.collectIsPressedAsState()

    // iOS-style spring animation (slightly different from Android)
    val scale by animateFloatAsState(
        targetValue = if (isPressed && enabled && onClick != null) 0.96f else 1f,
        animationSpec = spring(
            dampingRatio = 0.8f,
            stiffness = 400f
        ),
        label = "ios_card_scale"
    )

    // Size configuration
    val padding = when (size) {
        TchatCardSize.Compact -> 8.dp
        TchatCardSize.Standard -> TchatSpacing.cardPadding // 16dp
        TchatCardSize.Expanded -> 24.dp
    }

    // Variant styling with iOS-specific adjustments
    val (containerColor, elevation, border) = getIOSCardStyling(variant, enabled)
    val shape = RoundedCornerShape(TchatSpacing.cardBorderRadius)

    val cardModifier = modifier
        .scale(scale)
        .let { baseModifier ->
            if (elevation > 0.dp) {
                baseModifier.shadow(
                    elevation = elevation,
                    shape = shape,
                    clip = false
                )
            } else baseModifier
        }
        .background(
            color = containerColor,
            shape = shape
        )
        .clip(shape)
        .let { styledModifier ->
            border?.let { borderStroke ->
                styledModifier.border(
                    border = borderStroke,
                    shape = shape
                )
            } ?: styledModifier
        }
        .let { finalModifier ->
            if (onClick != null) {
                finalModifier.clickable(
                    interactionSource = interactionSource,
                    indication = null, // iOS doesn't use ripple effects
                    enabled = enabled,
                    role = Role.Button,
                    onClick = onClick
                )
            } else finalModifier
        }
        .padding(padding)

    Box(
        modifier = cardModifier
    ) {
        content()
    }
}

/**
 * Get iOS-specific styling properties based on card variant and state
 * Adjusted for iOS HIG guidelines and design patterns
 */
@Composable
private fun getIOSCardStyling(
    variant: TchatCardVariant,
    enabled: Boolean
): Triple<Color, Dp, BorderStroke?> {
    val alpha = if (enabled) 1f else 0.6f

    return when (variant) {
        TchatCardVariant.Elevated -> Triple(
            TchatColors.surface.copy(alpha = alpha),
            3.dp, // iOS uses slightly less elevation than Material
            null
        )

        TchatCardVariant.Outlined -> Triple(
            Color.Transparent,
            0.dp,
            BorderStroke(
                width = 0.5.dp, // iOS uses thinner borders
                color = TchatColors.outline.copy(alpha = alpha)
            )
        )

        TchatCardVariant.Filled -> Triple(
            TchatColors.surfaceVariant.copy(alpha = alpha),
            0.dp,
            null
        )

        TchatCardVariant.Glass -> Triple(
            TchatColors.surface.copy(alpha = 0.85f * alpha), // Slightly more opaque on iOS
            0.dp,
            BorderStroke(
                width = 0.33.dp, // Even thinner border for glass effect
                color = TchatColors.outline.copy(alpha = 0.2f * alpha)
            )
        )
    }
}