package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatSlider using Material3 Slider
 * Provides native Material Design slider with comprehensive theming
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
actual fun TchatSlider(
    value: Float,
    onValueChange: (Float) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    size: TchatSliderSize,
    showValueLabel: Boolean,
    valueFormatter: (Float) -> String,
    label: String?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val sliderHeight = when (size) {
        TchatSliderSize.Small -> 32.dp
        TchatSliderSize.Medium -> 40.dp
        TchatSliderSize.Large -> 48.dp
    }

    val textStyle = when (size) {
        TchatSliderSize.Small -> TchatTypography.typography.bodySmall
        TchatSliderSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSliderSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSliderSize.Small -> 16.dp
        TchatSliderSize.Medium -> 20.dp
        TchatSliderSize.Large -> 24.dp
    }

    val sliderColors = SliderDefaults.colors(
        thumbColor = TchatColors.primary,
        activeTrackColor = TchatColors.primary,
        inactiveTrackColor = TchatColors.outline.copy(alpha = 0.24f),
        disabledThumbColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledActiveTrackColor = TchatColors.onSurface.copy(alpha = 0.32f),
        disabledInactiveTrackColor = TchatColors.onSurface.copy(alpha = 0.12f)
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = labelText,
                    style = TchatTypography.typography.bodySmall,
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f),
                    modifier = Modifier.weight(1f)
                )

                if (showValueLabel) {
                    Text(
                        text = valueFormatter(value),
                        style = TchatTypography.typography.bodySmall,
                        color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.38f)
                    )
                }
            }
        }

        // Slider with icons
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                )
            }

            Slider(
                value = value,
                onValueChange = onValueChange,
                modifier = Modifier
                    .weight(1f)
                    .height(sliderHeight),
                enabled = enabled,
                valueRange = valueRange,
                steps = steps,
                colors = sliderColors,
                interactionSource = interactionSource
            )

            trailingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                )
            }
        }

        // Value range indicators
        if (steps > 0 || showValueLabel) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = valueFormatter(valueRange.start),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )

                Text(
                    text = valueFormatter(valueRange.endInclusive),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    }
}

/**
 * Android implementation of TchatRangeSlider using Material3 RangeSlider
 * Provides native Material Design range slider with comprehensive theming
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
actual fun TchatRangeSlider(
    value: ClosedFloatingPointRange<Float>,
    onValueChange: (ClosedFloatingPointRange<Float>) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    valueRange: ClosedFloatingPointRange<Float>,
    steps: Int,
    size: TchatSliderSize,
    showValueLabels: Boolean,
    valueFormatter: (Float) -> String,
    label: String?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    startInteractionSource: MutableInteractionSource,
    endInteractionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val sliderHeight = when (size) {
        TchatSliderSize.Small -> 32.dp
        TchatSliderSize.Medium -> 40.dp
        TchatSliderSize.Large -> 48.dp
    }

    val textStyle = when (size) {
        TchatSliderSize.Small -> TchatTypography.typography.bodySmall
        TchatSliderSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSliderSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSliderSize.Small -> 16.dp
        TchatSliderSize.Medium -> 20.dp
        TchatSliderSize.Large -> 24.dp
    }

    val sliderColors = SliderDefaults.colors(
        thumbColor = TchatColors.primary,
        activeTrackColor = TchatColors.primary,
        inactiveTrackColor = TchatColors.outline.copy(alpha = 0.24f),
        disabledThumbColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledActiveTrackColor = TchatColors.onSurface.copy(alpha = 0.32f),
        disabledInactiveTrackColor = TchatColors.onSurface.copy(alpha = 0.12f)
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        // Label
        label?.let { labelText ->
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = labelText,
                    style = TchatTypography.typography.bodySmall,
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f),
                    modifier = Modifier.weight(1f)
                )

                if (showValueLabels) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Text(
                            text = valueFormatter(value.start),
                            style = TchatTypography.typography.bodySmall,
                            color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.38f)
                        )
                        Text(
                            text = "-",
                            style = TchatTypography.typography.bodySmall,
                            color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                        )
                        Text(
                            text = valueFormatter(value.endInclusive),
                            style = TchatTypography.typography.bodySmall,
                            color = if (enabled) TchatColors.primary else TchatColors.onSurface.copy(alpha = 0.38f)
                        )
                    }
                }
            }
        }

        // Range Slider with icons
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                )
            }

            RangeSlider(
                value = value,
                onValueChange = onValueChange,
                modifier = Modifier
                    .weight(1f)
                    .height(sliderHeight),
                enabled = enabled,
                valueRange = valueRange,
                steps = steps,
                colors = sliderColors,
                startInteractionSource = startInteractionSource,
                endInteractionSource = endInteractionSource
            )

            trailingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                )
            }
        }

        // Value range indicators
        if (steps > 0 || showValueLabels) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = valueFormatter(valueRange.start),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )

                Text(
                    text = valueFormatter(valueRange.endInclusive),
                    style = TchatTypography.typography.bodySmall,
                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }
    }
}