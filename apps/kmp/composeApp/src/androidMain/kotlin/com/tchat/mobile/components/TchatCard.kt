package com.tchat.mobile.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * Android implementation of TchatCard using Material3 Card with sophisticated styling
 * Provides native Material Design 3 elevation, interaction states, and animations
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

    // Press animation
    val scale by animateFloatAsState(
        targetValue = if (isPressed && enabled && onClick != null) 0.98f else 1f,
        animationSpec = tween(150),
        label = "card_scale"
    )

    // Size configuration
    val padding = when (size) {
        TchatCardSize.Compact -> 8.dp
        TchatCardSize.Standard -> TchatSpacing.cardPadding // 16dp
        TchatCardSize.Expanded -> 24.dp
    }

    // Variant styling
    val (containerColor, elevation, border) = getCardStyling(variant, enabled)

    val cardModifier = modifier
        .scale(scale)
        .let { baseModifier ->
            if (onClick != null) {
                baseModifier.clickable(
                    interactionSource = interactionSource,
                    indication = null, // Custom press animation instead of ripple
                    enabled = enabled,
                    role = Role.Button,
                    onClick = onClick
                )
            } else baseModifier
        }

    when (variant) {
        TchatCardVariant.Glass -> {
            // Glass variant uses custom background instead of Card
            androidx.compose.foundation.layout.Box(
                modifier = cardModifier
                    .background(
                        color = containerColor,
                        shape = RoundedCornerShape(TchatSpacing.cardBorderRadius)
                    )
                    .padding(padding)
            ) {
                content()
            }
        }
        else -> {
            // Standard Card implementation for other variants
            Card(
                modifier = cardModifier,
                shape = RoundedCornerShape(TchatSpacing.cardBorderRadius),
                colors = CardDefaults.cardColors(
                    containerColor = containerColor,
                    disabledContainerColor = containerColor.copy(alpha = 0.6f)
                ),
                elevation = CardDefaults.cardElevation(
                    defaultElevation = elevation,
                    pressedElevation = if (onClick != null) elevation + 2.dp else elevation
                ),
                border = border
            ) {
                androidx.compose.foundation.layout.Box(
                    modifier = Modifier.padding(padding)
                ) {
                    content()
                }
            }
        }
    }
}

/**
 * Get styling properties based on card variant and state
 */
@Composable
private fun getCardStyling(
    variant: TchatCardVariant,
    enabled: Boolean
): Triple<Color, Dp, BorderStroke?> {
    val alpha = if (enabled) 1f else 0.6f

    return when (variant) {
        TchatCardVariant.Elevated -> Triple(
            TchatColors.surface.copy(alpha = alpha),
            TchatSpacing.cardElevation, // 4dp
            null
        )

        TchatCardVariant.Outlined -> Triple(
            Color.Transparent,
            0.dp,
            BorderStroke(
                width = 1.dp,
                color = TchatColors.outline.copy(alpha = alpha)
            )
        )

        TchatCardVariant.Filled -> Triple(
            TchatColors.surfaceVariant.copy(alpha = alpha),
            0.dp,
            null
        )

        TchatCardVariant.Glass -> Triple(
            TchatColors.surface.copy(alpha = 0.8f * alpha),
            0.dp,
            BorderStroke(
                width = 0.5.dp,
                color = TchatColors.outline.copy(alpha = 0.3f * alpha)
            )
        )
    }
}