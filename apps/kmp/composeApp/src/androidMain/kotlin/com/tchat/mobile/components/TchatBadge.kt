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
 * Android implementation of TchatBadge using Jetpack Compose and Material3
 * Provides Material Design 3 styling with elevation and rounded corners
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

        // Size configuration
        val (height, textStyle, horizontalPadding, verticalPadding, cornerRadius) = when (size) {
            TchatBadgeSize.Small -> BadgeSizeConfig(
                16.dp,
                TchatTypography.typography.labelSmall.copy(
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium
                ),
                6.dp,
                2.dp,
                8.dp
            )
            TchatBadgeSize.Medium -> BadgeSizeConfig(
                20.dp,
                TchatTypography.typography.labelMedium.copy(
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Medium
                ),
                8.dp,
                4.dp,
                10.dp
            )
            TchatBadgeSize.Large -> BadgeSizeConfig(
                24.dp,
                TchatTypography.typography.labelLarge.copy(
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold
                ),
                10.dp,
                4.dp,
                12.dp
            )
        }

        // Color configuration based on variant
        val (backgroundColor, contentColor) = getAndroidBadgeColors(variant)

        // Use circle shape for single digits/characters, rounded rectangle for longer text
        val shape = if (displayText.length == 1 && count != null) {
            CircleShape
        } else {
            RoundedCornerShape(cornerRadius)
        }

        Box(
            modifier = modifier
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
private fun getAndroidBadgeColors(variant: TchatBadgeVariant): Pair<Color, Color> {
    return when (variant) {
        TchatBadgeVariant.Default -> Pair(
            TchatColors.primary,
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Success -> Pair(
            TchatColors.success,
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Warning -> Pair(
            TchatColors.warning,
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Error -> Pair(
            TchatColors.error,
            TchatColors.onPrimary
        )
        TchatBadgeVariant.Info -> Pair(
            TchatColors.surfaceDim,
            TchatColors.onSurface
        )
    }
}

// Helper data class for Android badge configuration
private data class BadgeSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp,
    val cornerRadius: androidx.compose.ui.unit.Dp
)