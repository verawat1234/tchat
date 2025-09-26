package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.scaleIn
import androidx.compose.animation.scaleOut
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatBadge using Compose Multiplatform
 * Provides iOS-native styling with subtle shadows and SF-style typography
 *
 * Features iOS-specific behavior:
 * - Subtle drop shadows for depth (iOS style)
 * - SF Pro typography weights
 * - iOS HIG-compliant color semantics
 * - More pronounced corner rounding
 * - Smaller minimum sizes for iOS density
 */
@Composable
actual fun TchatBadge(
    text: String,
    modifier: Modifier,
    variant: TchatBadgeVariant,
    size: TchatBadgeSize,
    count: Int?,
    maxCount: Int,
    showZero: Boolean,
    contentDescription: String?
) {
    // Determine if badge should be visible
    val shouldShow = count?.let { it > 0 || showZero } ?: true

    AnimatedVisibility(
        visible = shouldShow,
        enter = scaleIn() + fadeIn(),
        exit = scaleOut() + fadeOut()
    ) {
        val displayText = count?.let {
            if (it > maxCount) "${maxCount}+" else it.toString()
        } ?: text

        // iOS-specific size configuration
        val (height, textStyle, horizontalPadding, verticalPadding, cornerRadius, shadowElevation) = when (size) {
            TchatBadgeSize.Small -> IOSBadgeSizeConfig(
                18.dp, // Slightly larger for iOS
                TchatTypography.typography.labelSmall.copy(
                    fontSize = 11.sp, // iOS uses slightly smaller text
                    fontWeight = FontWeight.SemiBold // More weight for iOS
                ),
                6.dp,
                3.dp,
                9.dp, // More rounded for iOS
                1.dp
            )
            TchatBadgeSize.Medium -> IOSBadgeSizeConfig(
                22.dp,
                TchatTypography.typography.labelMedium.copy(
                    fontSize = 13.sp,
                    fontWeight = FontWeight.SemiBold
                ),
                8.dp,
                4.dp,
                11.dp,
                1.5.dp
            )
            TchatBadgeSize.Large -> IOSBadgeSizeConfig(
                26.dp,
                TchatTypography.typography.labelLarge.copy(
                    fontSize = 15.sp,
                    fontWeight = FontWeight.Bold // Bold for iOS large badges
                ),
                12.dp,
                5.dp,
                13.dp,
                2.dp
            )
        }

        // iOS-themed color configuration
        val (backgroundColor, contentColor) = getIOSBadgeColors(variant)

        // iOS prefers slightly more rounded badges, circle for single digits
        val shape = if (displayText.length == 1 && count != null) {
            CircleShape
        } else {
            RoundedCornerShape(cornerRadius)
        }

        Box(
            modifier = modifier
                .shadow(shadowElevation, shape) // iOS-style subtle shadow
                .heightIn(min = height)
                .widthIn(min = if (displayText.length == 1 && count != null) height else height)
                .background(
                    color = backgroundColor,
                    shape = shape
                )
                .clip(shape)
                .padding(horizontal = horizontalPadding, vertical = verticalPadding)
                .semantics {
                    this.contentDescription = contentDescription ?: when {
                        count != null -> "$count ${if (count == 1) "notification" else "notifications"}"
                        else -> text
                    }
                },
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = displayText,
                style = textStyle,
                color = contentColor,
                textAlign = TextAlign.Center,
                maxLines = 1
            )
        }
    }
}

@Composable
private fun getIOSBadgeColors(variant: TchatBadgeVariant): Pair<Color, Color> {
    return when (variant) {
        TchatBadgeVariant.Default -> Pair(
            TchatColors.primary,
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Success -> Pair(
            TchatColors.success.copy(alpha = 0.95f), // iOS uses slightly muted colors
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Warning -> Pair(
            TchatColors.warning.copy(alpha = 0.95f),
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Error -> Pair(
            TchatColors.error.copy(alpha = 0.95f),
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Info -> Pair(
            TchatColors.surfaceDim.copy(alpha = 0.9f), // More subtle for iOS
            TchatColors.onSurface
        )
    }
}

// Helper data class for iOS badge configuration
private data class IOSBadgeSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp,
    val cornerRadius: androidx.compose.ui.unit.Dp,
    val shadowElevation: androidx.compose.ui.unit.Dp
)