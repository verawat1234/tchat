package com.tchat.mobile.components

import androidx.compose.foundation.text.BasicText
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatLabel with SwiftUI-inspired styling
 * Uses iOS HIG typography patterns and accessibility guidelines
 */
@Composable
actual fun TchatLabel(
    text: String,
    modifier: Modifier,
    style: TchatLabelStyle
) {
    val (textStyle, textColor, fontWeight, shouldUppercase) = when (style) {
        TchatLabelStyle.Body -> IOSLabelStyleConfig(
            textStyle = TchatTypography.typography.bodyMedium,
            color = TchatColors.onSurface,
            fontWeight = FontWeight.Normal,
            shouldUppercase = false
        )

        TchatLabelStyle.Caption -> IOSLabelStyleConfig(
            textStyle = TchatTypography.typography.bodySmall,
            color = TchatColors.onSurfaceVariant.copy(alpha = 0.8f), // iOS uses slightly more transparent secondary text
            fontWeight = FontWeight.Normal,
            shouldUppercase = false
        )

        TchatLabelStyle.Overline -> IOSLabelStyleConfig(
            textStyle = TchatTypography.typography.labelSmall,
            color = TchatColors.onSurfaceVariant.copy(alpha = 0.7f), // iOS overlines are more subtle
            fontWeight = FontWeight.SemiBold, // iOS uses SemiBold instead of Medium
            shouldUppercase = true
        )

        TchatLabelStyle.Required -> IOSLabelStyleConfig(
            textStyle = TchatTypography.typography.bodyMedium,
            color = TchatColors.error,
            fontWeight = FontWeight.SemiBold, // iOS uses SemiBold for emphasis
            shouldUppercase = false
        )
    }

    BasicText(
        text = if (shouldUppercase) text.uppercase() else text,
        modifier = modifier,
        style = textStyle.copy(
            fontWeight = fontWeight,
            color = textColor
        )
    )
}

/**
 * Data class for iOS label styling configuration
 */
private data class IOSLabelStyleConfig(
    val textStyle: androidx.compose.ui.text.TextStyle,
    val color: androidx.compose.ui.graphics.Color,
    val fontWeight: FontWeight,
    val shouldUppercase: Boolean
)