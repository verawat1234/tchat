package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatCheckbox using Material3 Checkbox
 * Provides native Material Design checkbox with comprehensive theming
 */
@Composable
actual fun TchatCheckbox(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatCheckboxSize,
    interactionSource: MutableInteractionSource,
    label: String?
) {
    val checkboxSize = when (size) {
        TchatCheckboxSize.Small -> 16.dp
        TchatCheckboxSize.Medium -> 20.dp
        TchatCheckboxSize.Large -> 24.dp
    }

    val checkboxColors = CheckboxDefaults.colors(
        checkedColor = TchatColors.primary,
        uncheckedColor = TchatColors.outline,
        checkmarkColor = TchatColors.onPrimary,
        disabledCheckedColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledUncheckedColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledIndeterminateColor = TchatColors.onSurface.copy(alpha = 0.38f)
    )

    if (label != null) {
        Row(
            modifier = modifier,
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Checkbox(
                checked = checked,
                onCheckedChange = onCheckedChange,
                modifier = Modifier.size(checkboxSize),
                enabled = enabled,
                colors = checkboxColors,
                interactionSource = interactionSource
            )

            if (enabled && onCheckedChange != null) {
                Text(
                    text = label,
                    style = when (size) {
                        TchatCheckboxSize.Small -> TchatTypography.typography.bodySmall
                        TchatCheckboxSize.Medium -> TchatTypography.typography.bodyMedium
                        TchatCheckboxSize.Large -> TchatTypography.typography.bodyLarge
                    },
                    color = TchatColors.onSurface,
                    modifier = Modifier.weight(1f)
                )
            } else {
                Text(
                    text = label,
                    style = when (size) {
                        TchatCheckboxSize.Small -> TchatTypography.typography.bodySmall
                        TchatCheckboxSize.Medium -> TchatTypography.typography.bodyMedium
                        TchatCheckboxSize.Large -> TchatTypography.typography.bodyLarge
                    },
                    color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f),
                    modifier = Modifier.weight(1f)
                )
            }
        }
    } else {
        Checkbox(
            checked = checked,
            onCheckedChange = onCheckedChange,
            modifier = modifier.size(checkboxSize),
            enabled = enabled,
            colors = checkboxColors,
            interactionSource = interactionSource
        )
    }
}