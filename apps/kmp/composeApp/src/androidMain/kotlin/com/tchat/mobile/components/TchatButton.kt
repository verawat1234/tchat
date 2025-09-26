package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatButton using Jetpack Compose and Material3
 * Provides Material Design 3 styling with sophisticated animations
 */
@Composable
actual fun TchatButton(
    onClick: () -> Unit,
    text: String,
    modifier: Modifier,
    variant: TchatButtonVariant,
    size: TchatButtonSize,
    enabled: Boolean,
    loading: Boolean,
    leadingIcon: (@Composable () -> Unit)?,
    trailingIcon: (@Composable () -> Unit)?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    // Interaction state management
    val isPressed by interactionSource.collectIsPressedAsState()
    val scale by animateFloatAsState(
        targetValue = if (isPressed && enabled) 0.95f else 1f,
        animationSpec = tween(150),
        label = "button_scale"
    )

    // Size configuration
    val (height, textStyle, horizontalPadding, verticalPadding) = when (size) {
        TchatButtonSize.Small -> ButtonSizeConfig(
            32.dp,
            TchatTypography.typography.labelMedium,
            12.dp,
            6.dp
        )
        TchatButtonSize.Medium -> ButtonSizeConfig(
            TchatSpacing.buttonMinHeight,
            TchatTypography.typography.labelLarge,
            TchatSpacing.buttonPaddingHorizontal,
            TchatSpacing.buttonPaddingVertical
        )
        TchatButtonSize.Large -> ButtonSizeConfig(
            48.dp,
            TchatTypography.typography.labelLarge.copy(fontSize = 18.sp),
            20.dp,
            12.dp
        )
    }

    // Color and style configuration based on variant
    val (backgroundColor, contentColor, borderColor, borderWidth) = getVariantColors(
        variant = variant,
        enabled = enabled,
        isPressed = isPressed
    )

    // Animated colors for smooth transitions
    val animatedBackgroundColor by animateColorAsState(
        targetValue = backgroundColor,
        animationSpec = tween(150),
        label = "background_color"
    )
    val animatedContentColor by animateColorAsState(
        targetValue = contentColor,
        animationSpec = tween(150),
        label = "content_color"
    )
    val animatedBorderColor by animateColorAsState(
        targetValue = borderColor,
        animationSpec = tween(150),
        label = "border_color"
    )

    val shape = RoundedCornerShape(TchatSpacing.buttonBorderRadius)

    Box(
        modifier = modifier
            .scale(scale)
            .heightIn(min = height)
            .background(
                color = animatedBackgroundColor,
                shape = shape
            )
            .let { baseModifier ->
                if (borderWidth > 0.dp) {
                    baseModifier.border(
                        width = borderWidth,
                        color = animatedBorderColor,
                        shape = shape
                    )
                } else baseModifier
            }
            .clickable(
                interactionSource = interactionSource,
                indication = null, // Custom animation handling
                enabled = enabled && !loading,
                onClick = onClick
            )
            .padding(horizontal = horizontalPadding, vertical = verticalPadding)
            .semantics {
                this.contentDescription = contentDescription ?: when {
                    loading -> "$text, Loading"
                    !enabled -> "$text, Disabled"
                    else -> text
                }
            },
        contentAlignment = Alignment.Center
    ) {
        Row(
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Leading icon
            leadingIcon?.let { icon ->
                Box(
                    modifier = Modifier.alpha(if (loading) 0f else 1f)
                ) {
                    CompositionLocalProvider(
                        LocalContentColor provides animatedContentColor
                    ) {
                        icon()
                    }
                }
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
            }

            // Button text with loading state handling
            Box(
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = text,
                    style = textStyle,
                    color = animatedContentColor,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.alpha(if (loading) 0f else 1f)
                )

                // Loading indicator
                if (loading) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(16.dp),
                        strokeWidth = 2.dp,
                        color = animatedContentColor
                    )
                }
            }

            // Trailing icon
            trailingIcon?.let { icon ->
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Box(
                    modifier = Modifier.alpha(if (loading) 0f else 1f)
                ) {
                    CompositionLocalProvider(
                        LocalContentColor provides animatedContentColor
                    ) {
                        icon()
                    }
                }
            }
        }
    }
}

@Composable
private fun getVariantColors(
    variant: TchatButtonVariant,
    enabled: Boolean,
    isPressed: Boolean
): ButtonColorConfig {
    return when (variant) {
        TchatButtonVariant.Primary -> {
            val backgroundColor = when {
                !enabled -> TchatColors.disabled
                isPressed -> TchatColors.primaryDark
                else -> TchatColors.primary
            }
            ButtonColorConfig(backgroundColor, TchatColors.onPrimary, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Secondary -> {
            val backgroundColor = when {
                !enabled -> TchatColors.disabled.copy(alpha = 0.1f)
                isPressed -> TchatColors.surfaceDim
                else -> TchatColors.surface
            }
            val contentColor = if (enabled) TchatColors.onSurface else TchatColors.disabled
            ButtonColorConfig(backgroundColor, contentColor, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Ghost -> {
            val backgroundColor = when {
                !enabled -> Color.Transparent
                isPressed -> TchatColors.primary.copy(alpha = 0.1f)
                else -> Color.Transparent
            }
            val contentColor = if (enabled) TchatColors.primary else TchatColors.disabled
            ButtonColorConfig(backgroundColor, contentColor, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Destructive -> {
            val backgroundColor = when {
                !enabled -> TchatColors.disabled
                isPressed -> TchatColors.error.copy(alpha = 0.9f)
                else -> TchatColors.error
            }
            ButtonColorConfig(backgroundColor, TchatColors.onPrimary, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Outline -> {
            val backgroundColor = when {
                !enabled -> Color.Transparent
                isPressed -> TchatColors.outline.copy(alpha = 0.1f)
                else -> Color.Transparent
            }
            val contentColor = if (enabled) TchatColors.onSurface else TchatColors.disabled
            val borderColor = if (enabled) TchatColors.outline else TchatColors.disabled
            ButtonColorConfig(backgroundColor, contentColor, borderColor, TchatSpacing.buttonBorderWidth)
        }
    }
}

// Helper data classes for configuration
private data class ButtonSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
)

private data class ButtonColorConfig(
    val backgroundColor: Color,
    val contentColor: Color,
    val borderColor: Color,
    val borderWidth: androidx.compose.ui.unit.Dp
)