package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * iOS implementation of TchatSpinner using Compose Multiplatform
 * Provides iOS-native styling with refined colors and smoother animations
 *
 * Features iOS-specific behavior:
 * - iOS HIG-compliant activity indicator styling
 * - Softer color variants for better visual hierarchy
 * - Smoother spring-like animation timing
 * - Refined stroke widths for iOS density
 * - More subtle track colors for determinate progress
 */
@Composable
actual fun TchatSpinner(
    modifier: Modifier,
    variant: TchatSpinnerVariant,
    size: TchatSpinnerSize,
    progress: Float?,
    strokeWidth: Float?,
    contentDescription: String?
) {
    // iOS-specific size configuration
    val (spinnerSize, defaultStrokeWidth) = when (size) {
        TchatSpinnerSize.Small -> IOSSpinnerSizeConfig(
            16.dp,
            1.8.dp // Slightly thinner for iOS
        )
        TchatSpinnerSize.Medium -> IOSSpinnerSizeConfig(
            20.dp,
            2.2.dp // iOS refined thickness
        )
        TchatSpinnerSize.Large -> IOSSpinnerSizeConfig(
            24.dp,
            2.8.dp // More refined than Android
        )
    }

    val actualStrokeWidth = strokeWidth?.dp ?: defaultStrokeWidth

    // iOS-themed color configuration with softer variants
    val spinnerColor = getIOSSpinnerColor(variant)
    val animatedColor by animateColorAsState(
        targetValue = spinnerColor,
        animationSpec = tween(400), // Smoother iOS transitions
        label = "ios_spinner_color"
    )

    Box(
        modifier = modifier
            .size(spinnerSize)
            .semantics {
                this.contentDescription = contentDescription ?: when (variant) {
                    TchatSpinnerVariant.Default -> "Loading"
                    TchatSpinnerVariant.Success -> "Processing success"
                    TchatSpinnerVariant.Warning -> "Processing with caution"
                    TchatSpinnerVariant.Error -> "Processing error"
                }
            },
        contentAlignment = Alignment.Center
    ) {
        if (progress != null) {
            // iOS-style determinate progress indicator
            CircularProgressIndicator(
                progress = progress.coerceIn(0f, 1f),
                modifier = Modifier.fillMaxSize(),
                color = animatedColor,
                strokeWidth = actualStrokeWidth,
                trackColor = animatedColor.copy(alpha = 0.08f) // More subtle iOS track
            )
        } else {
            // iOS-style indeterminate progress indicator
            CircularProgressIndicator(
                modifier = Modifier.fillMaxSize(),
                color = animatedColor,
                strokeWidth = actualStrokeWidth,
                trackColor = Color.Transparent // Clean iOS style
            )
        }
    }
}

@Composable
private fun getIOSSpinnerColor(variant: TchatSpinnerVariant): Color {
    return when (variant) {
        TchatSpinnerVariant.Default -> TchatColors.primary.copy(alpha = 0.95f) // Softer iOS blue
        TchatSpinnerVariant.Success -> TchatColors.success.copy(alpha = 0.95f) // Softer iOS green
        TchatSpinnerVariant.Warning -> TchatColors.warning.copy(alpha = 0.95f) // Softer iOS amber
        TchatSpinnerVariant.Error -> TchatColors.error.copy(alpha = 0.95f) // Softer iOS red
    }
}

// Helper data class for iOS spinner configuration
private data class IOSSpinnerSizeConfig(
    val size: androidx.compose.ui.unit.Dp,
    val strokeWidth: androidx.compose.ui.unit.Dp
)