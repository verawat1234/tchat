package com.tchat.mobile.components

import androidx.compose.animation.core.*
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.toggleable
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatSwitch using Material3 Switch
 * Provides native Material Design switch with comprehensive theming and animations
 */
@Composable
actual fun TchatSwitch(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatSwitchSize,
    isLoading: Boolean,
    label: String?,
    description: String?,
    leadingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val switchScale = when (size) {
        TchatSwitchSize.Small -> 0.8f
        TchatSwitchSize.Medium -> 1.0f
        TchatSwitchSize.Large -> 1.25f
    }

    val textStyle = when (size) {
        TchatSwitchSize.Small -> TchatTypography.typography.bodySmall
        TchatSwitchSize.Medium -> TchatTypography.typography.bodyMedium
        TchatSwitchSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatSwitchSize.Small -> 16.dp
        TchatSwitchSize.Medium -> 20.dp
        TchatSwitchSize.Large -> 24.dp
    }

    // Material 3 switch colors
    val switchColors = SwitchDefaults.colors(
        checkedThumbColor = TchatColors.onPrimary,
        checkedTrackColor = TchatColors.primary,
        uncheckedThumbColor = TchatColors.outline,
        uncheckedTrackColor = TchatColors.surfaceVariant,
        disabledCheckedThumbColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledCheckedTrackColor = TchatColors.onSurface.copy(alpha = 0.12f),
        disabledUncheckedThumbColor = TchatColors.onSurface.copy(alpha = 0.38f),
        disabledUncheckedTrackColor = TchatColors.onSurface.copy(alpha = 0.12f)
    )

    // Loading animation
    val loadingAlpha by animateFloatAsState(
        targetValue = if (isLoading) 0.5f else 1f,
        animationSpec = tween(300),
        label = "LoadingAlpha"
    )

    val switchComponent = @Composable {
        Box(
            contentAlignment = Alignment.Center
        ) {
            Switch(
                checked = checked,
                onCheckedChange = if (enabled && !isLoading) onCheckedChange else null,
                modifier = Modifier.scale(switchScale),
                enabled = enabled && !isLoading,
                colors = switchColors,
                interactionSource = interactionSource
            )

            // Loading indicator
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.size((iconSize.value * 0.6f).dp),
                    strokeWidth = 2.dp,
                    color = TchatColors.primary
                )
            }
        }
    }

    if (label != null || description != null || leadingIcon != null) {
        Row(
            modifier = modifier
                .toggleable(
                    value = checked,
                    enabled = enabled && !isLoading,
                    role = Role.Switch,
                    onValueChange = onCheckedChange ?: {}
                )
                .padding(vertical = 4.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Leading icon
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) {
                        if (checked) TchatColors.primary else TchatColors.onSurface
                    } else {
                        TchatColors.onSurface.copy(alpha = 0.38f)
                    }
                )
            }

            // Label and description
            Column(
                modifier = Modifier.weight(1f),
                verticalArrangement = Arrangement.spacedBy(2.dp)
            ) {
                label?.let { labelText ->
                    Text(
                        text = labelText,
                        style = textStyle,
                        color = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                    )
                }

                description?.let { desc ->
                    Text(
                        text = desc,
                        style = TchatTypography.typography.bodySmall,
                        color = if (enabled) TchatColors.onSurface.copy(alpha = 0.7f) else TchatColors.onSurface.copy(alpha = 0.3f)
                    )
                }
            }

            // Switch
            switchComponent()
        }
    } else {
        Box(
            modifier = modifier
        ) {
            switchComponent()
        }
    }
}