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
 * Android implementation of TchatSpinner using Jetpack Compose and Material3
 * Uses Material Design circular progress indicators with smooth color transitions
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
    // Size configuration
    val (spinnerSize, defaultStrokeWidth) = when (size) {
        TchatSpinnerSize.Small -> SpinnerSizeConfig(16.dp, 2.dp)
        TchatSpinnerSize.Medium -> SpinnerSizeConfig(20.dp, 2.5.dp)
        TchatSpinnerSize.Large -> SpinnerSizeConfig(24.dp, 3.dp)
    }

    val actualStrokeWidth = strokeWidth?.dp ?: defaultStrokeWidth

    // Color configuration with smooth transitions
    val spinnerColor = getAndroidSpinnerColor(variant)
    val animatedColor by animateColorAsState(
        targetValue = spinnerColor,
        animationSpec = tween(300),
        label = "spinner_color"
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
            // Determinate progress indicator
            CircularProgressIndicator(
                progress = progress.coerceIn(0f, 1f),
                modifier = Modifier.fillMaxSize(),
                color = animatedColor,
                strokeWidth = actualStrokeWidth,
                trackColor = animatedColor.copy(alpha = 0.12f) // Material3 track color
            )
        } else {
            // Indeterminate progress indicator
            CircularProgressIndicator(
                modifier = Modifier.fillMaxSize(),
                color = animatedColor,
                strokeWidth = actualStrokeWidth,
                trackColor = Color.Transparent // No track for indeterminate
            )
        }
    }
}

@Composable
private fun getAndroidSpinnerColor(variant: TchatSpinnerVariant): Color {
    return when (variant) {
        TchatSpinnerVariant.Default -> TchatColors.primary
        TchatSpinnerVariant.Success -> TchatColors.success
        TchatSpinnerVariant.Warning -> TchatColors.warning
        TchatSpinnerVariant.Error -> TchatColors.error
    }
}

// Helper data class for Android spinner configuration
private data class SpinnerSizeConfig(
    val size: androidx.compose.ui.unit.Dp,
    val strokeWidth: androidx.compose.ui.unit.Dp
)