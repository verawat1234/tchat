package com.tchat.mobile.components

import androidx.compose.material3.MaterialTheme
import androidx.compose.foundation.text.BasicText
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatLabel using Material3 Text with semantic styling
 * Provides proper typography scaling and accessibility semantics
 */
@Composable
actual fun TchatLabel(
    text: String,
    modifier: Modifier,
    style: TchatLabelStyle
) {
    val (textStyle, textColor, fontWeight, shouldUppercase) = when (style) {
        TchatLabelStyle.Body -> LabelStyleConfig(
            textStyle = TchatTypography.typography.bodyMedium,
            color = TchatColors.onSurface,
            fontWeight = FontWeight.Normal,
            shouldUppercase = false
        )

        TchatLabelStyle.Caption -> LabelStyleConfig(
            textStyle = TchatTypography.typography.bodySmall,
            color = TchatColors.onSurfaceVariant,
            fontWeight = FontWeight.Normal,
            shouldUppercase = false
        )

        TchatLabelStyle.Overline -> LabelStyleConfig(
            textStyle = TchatTypography.typography.labelSmall,
            color = TchatColors.onSurfaceVariant,
            fontWeight = FontWeight.Medium,
            shouldUppercase = true
        )

        TchatLabelStyle.Required -> LabelStyleConfig(
            textStyle = TchatTypography.typography.bodyMedium,
            color = TchatColors.error,
            fontWeight = FontWeight.Medium,
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
 * Data class for label styling configuration
 */
private data class LabelStyleConfig(
    val textStyle: androidx.compose.ui.text.TextStyle,
    val color: androidx.compose.ui.graphics.Color,
    val fontWeight: FontWeight,
    val shouldUppercase: Boolean
)