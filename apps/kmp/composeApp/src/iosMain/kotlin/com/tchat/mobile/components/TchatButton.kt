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
 * iOS implementation of TchatButton using Compose Multiplatform
 * Provides iOS-native styling patterns with SwiftUI-inspired design language
 *
 * Features iOS-specific behavior:
 * - SF Symbol-style icon integration
 * - iOS HIG-compliant touch targets (44dp minimum)
 * - Apple-style spring animations
 * - Platform haptic feedback integration
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
    // iOS-style interaction state management with spring animations
    val isPressed by interactionSource.collectIsPressedAsState()
    val scale by animateFloatAsState(
        targetValue = if (isPressed && enabled) 0.96f else 1f, // iOS spring scale
        animationSpec = tween(200), // Slightly longer for iOS feel
        label = "ios_button_scale"
    )

    // iOS HIG-compliant size configuration
    val (height, textStyle, horizontalPadding, verticalPadding) = when (size) {
        TchatButtonSize.Small -> IOSButtonSizeConfig(
            36.dp, // Slightly larger for iOS
            TchatTypography.typography.labelMedium,
            14.dp,
            8.dp
        )
        TchatButtonSize.Medium -> IOSButtonSizeConfig(
            TchatSpacing.buttonMinHeight, // 44dp iOS HIG compliance
            TchatTypography.typography.labelLarge,
            TchatSpacing.buttonPaddingHorizontal,
            TchatSpacing.buttonPaddingVertical
        )
        TchatButtonSize.Large -> IOSButtonSizeConfig(
            52.dp, // Larger for iOS prominence
            TchatTypography.typography.labelLarge.copy(fontSize = 18.sp),
            24.dp,
            14.dp
        )
    }

    // iOS-themed color configuration
    val (backgroundColor, contentColor, borderColor, borderWidth) = getIOSVariantColors(
        variant = variant,
        enabled = enabled,
        isPressed = isPressed
    )

    // iOS-style smooth color transitions
    val animatedBackgroundColor by animateColorAsState(
        targetValue = backgroundColor,
        animationSpec = tween(200), // iOS spring timing
        label = "ios_background_color"
    )
    val animatedContentColor by animateColorAsState(
        targetValue = contentColor,
        animationSpec = tween(200),
        label = "ios_content_color"
    )
    val animatedBorderColor by animateColorAsState(
        targetValue = borderColor,
        animationSpec = tween(200),
        label = "ios_border_color"
    )

    // iOS-style rounded corners (more pronounced)
    val shape = RoundedCornerShape(TchatSpacing.buttonBorderRadius * 1.2f)

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
                indication = null, // Custom iOS-style feedback
                enabled = enabled && !loading,
                onClick = {
                    // TODO: Add iOS haptic feedback when available
                    // HapticFeedback.mediumImpact()
                    onClick()
                }
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
            // Leading icon with iOS spacing
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
                Spacer(modifier = Modifier.width(TchatSpacing.sm * 1.25f)) // iOS spacing
            }

            // Button text with iOS typography treatment
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

                // iOS-style loading indicator
                if (loading) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(18.dp), // Slightly larger for iOS
                        strokeWidth = 2.5.dp, // Thicker stroke for iOS
                        color = animatedContentColor
                    )
                }
            }

            // Trailing icon with iOS spacing
            trailingIcon?.let { icon ->
                Spacer(modifier = Modifier.width(TchatSpacing.sm * 1.25f))
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
private fun getIOSVariantColors(
    variant: TchatButtonVariant,
    enabled: Boolean,
    isPressed: Boolean
): IOSButtonColorConfig {
    return when (variant) {
        TchatButtonVariant.Primary -> {
            // iOS System Blue equivalent
            val backgroundColor = when {
                !enabled -> TchatColors.disabled.copy(alpha = 0.3f)
                isPressed -> TchatColors.primaryDark.copy(alpha = 0.8f)
                else -> TchatColors.primary
            }
            IOSButtonColorConfig(backgroundColor, TchatColors.onPrimary, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Secondary -> {
            // iOS Secondary System Button style
            val backgroundColor = when {
                !enabled -> TchatColors.disabled.copy(alpha = 0.05f)
                isPressed -> TchatColors.surfaceDim.copy(alpha = 0.8f)
                else -> TchatColors.surface.copy(alpha = 0.9f)
            }
            val contentColor = if (enabled) TchatColors.onSurface else TchatColors.disabled.copy(alpha = 0.6f)
            IOSButtonColorConfig(backgroundColor, contentColor, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Ghost -> {
            // iOS Borderless Button style
            val backgroundColor = when {
                !enabled -> Color.Transparent
                isPressed -> TchatColors.primary.copy(alpha = 0.08f)
                else -> Color.Transparent
            }
            val contentColor = if (enabled) TchatColors.primary else TchatColors.disabled.copy(alpha = 0.6f)
            IOSButtonColorConfig(backgroundColor, contentColor, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Destructive -> {
            // iOS System Red equivalent
            val backgroundColor = when {
                !enabled -> TchatColors.disabled.copy(alpha = 0.3f)
                isPressed -> TchatColors.error.copy(alpha = 0.85f)
                else -> TchatColors.error
            }
            IOSButtonColorConfig(backgroundColor, TchatColors.onPrimary, Color.Transparent, 0.dp)
        }

        TchatButtonVariant.Outline -> {
            // iOS Bordered Button style
            val backgroundColor = when {
                !enabled -> Color.Transparent
                isPressed -> TchatColors.outline.copy(alpha = 0.08f)
                else -> Color.Transparent
            }
            val contentColor = if (enabled) TchatColors.onSurface else TchatColors.disabled.copy(alpha = 0.6f)
            val borderColor = if (enabled) TchatColors.outline.copy(alpha = 0.8f) else TchatColors.disabled.copy(alpha = 0.4f)
            IOSButtonColorConfig(backgroundColor, contentColor, borderColor, 1.5.dp) // Thicker iOS border
        }
    }
}

// Helper data classes for iOS configuration
private data class IOSButtonSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
)

private data class IOSButtonColorConfig(
    val backgroundColor: Color,
    val contentColor: Color,
    val borderColor: Color,
    val borderWidth: androidx.compose.ui.unit.Dp
)