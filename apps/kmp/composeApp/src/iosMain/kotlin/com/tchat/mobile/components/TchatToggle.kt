package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.toggleable
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatToggle with SwiftUI-inspired styling
 * Uses custom animations and interactions for iOS-style toggle behavior
 */
@Composable
actual fun TchatToggle(
    pressed: Boolean,
    onPressedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatToggleSize,
    variant: TchatToggleVariant,
    text: String?,
    icon: ImageVector?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val buttonHeight = when (size) {
        TchatToggleSize.Small -> 36.dp // iOS uses larger touch targets
        TchatToggleSize.Medium -> 44.dp
        TchatToggleSize.Large -> 52.dp
    }

    val textStyle = when (size) {
        TchatToggleSize.Small -> TchatTypography.typography.bodySmall
        TchatToggleSize.Medium -> TchatTypography.typography.bodyMedium
        TchatToggleSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatToggleSize.Small -> 18.dp
        TchatToggleSize.Medium -> 22.dp
        TchatToggleSize.Large -> 26.dp
    }

    val horizontalPadding = when (size) {
        TchatToggleSize.Small -> 14.dp // iOS uses more padding
        TchatToggleSize.Medium -> 18.dp
        TchatToggleSize.Large -> 22.dp
    }

    // iOS-style spring animations
    val backgroundColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Default -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.surface.copy(alpha = 0.8f) // iOS uses more subtle backgrounds
            }
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary.copy(alpha = 0.08f) // More subtle than Android
            } else {
                Color.Transparent
            }
            TchatToggleVariant.Text -> Color.Transparent
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_background_color"
    )

    // iOS-style content color animation
    val contentColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Default -> if (pressed) {
                TchatColors.onPrimary
            } else {
                TchatColors.onSurface
            }
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.onSurface
            }
            TchatToggleVariant.Text -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.onSurface
            }
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_content_color"
    )

    // iOS-style border animation
    val borderColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.outline.copy(alpha = 0.4f) // iOS uses more subtle borders
            }
            else -> Color.Transparent
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_border_color"
    )

    // iOS-style scale animation with spring
    val scale by animateFloatAsState(
        targetValue = if (pressed) 0.96f else 1f, // iOS uses more subtle scale
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_scale"
    )

    Surface(
        modifier = modifier
            .height(buttonHeight)
            .scale(scale)
            .clip(RoundedCornerShape(10.dp)) // iOS uses more rounded corners
            .toggleable(
                value = pressed,
                enabled = enabled,
                role = Role.Button,
                interactionSource = interactionSource,
                indication = null, // iOS doesn't use ripple
                onValueChange = onPressedChange ?: {}
            ),
        color = backgroundColor,
        shape = RoundedCornerShape(10.dp),
        border = if (variant == TchatToggleVariant.Outline) {
            BorderStroke(1.5.dp, borderColor) // iOS uses slightly thicker borders
        } else null,
        shadowElevation = if (pressed && variant == TchatToggleVariant.Default) {
            2.dp // iOS-style subtle elevation
        } else 0.dp
    ) {
        Row(
            modifier = Modifier
                .fillMaxHeight()
                .padding(horizontal = horizontalPadding),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(10.dp) // iOS uses more spacing
        ) {
            // Leading icon
            leadingIcon?.let { leadingIconVector ->
                Icon(
                    imageVector = leadingIconVector,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) contentColor else contentColor.copy(alpha = 0.5f) // iOS uses less transparency
                )
            }

            // Main icon (if no text)
            if (text == null) {
                icon?.let { iconVector ->
                    Icon(
                        imageVector = iconVector,
                        contentDescription = contentDescription,
                        modifier = Modifier.size(iconSize),
                        tint = if (enabled) contentColor else contentColor.copy(alpha = 0.5f)
                    )
                }
            }

            // Text
            text?.let { textContent ->
                Text(
                    text = textContent,
                    style = textStyle.copy(
                        fontWeight = if (pressed) FontWeight.SemiBold else FontWeight.Medium // iOS uses different weights
                    ),
                    color = if (enabled) contentColor else contentColor.copy(alpha = 0.5f)
                )
            }

            // Trailing icon
            trailingIcon?.let { trailingIconVector ->
                Icon(
                    imageVector = trailingIconVector,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) contentColor else contentColor.copy(alpha = 0.5f)
                )
            }
        }
    }
}

/**
 * iOS implementation of TchatToggleGroup with SwiftUI-inspired segmented control styling
 */
@Composable
actual fun TchatToggleGroup(
    options: List<String>,
    selectedOptions: Set<String>,
    onSelectionChange: (Set<String>) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatToggleSize,
    variant: TchatToggleVariant,
    allowMultipleSelection: Boolean,
    allowEmptySelection: Boolean,
    maxSelections: Int?,
    interactionSource: MutableInteractionSource
) {
    // iOS-style segmented control container
    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(12.dp), // iOS uses more rounded corners
        color = TchatColors.surface.copy(alpha = 0.6f), // iOS background
        border = BorderStroke(
            width = 1.dp,
            color = TchatColors.outline.copy(alpha = 0.3f)
        )
    ) {
        Row(
            modifier = Modifier.padding(2.dp), // iOS inner padding
            horizontalArrangement = Arrangement.spacedBy(0.dp)
        ) {
            options.forEachIndexed { index, option ->
                val isSelected = selectedOptions.contains(option)
                val isFirstItem = index == 0
                val isLastItem = index == options.size - 1

                val shape = when {
                    options.size == 1 -> RoundedCornerShape(10.dp)
                    isFirstItem -> RoundedCornerShape(topStart = 10.dp, bottomStart = 10.dp)
                    isLastItem -> RoundedCornerShape(topEnd = 10.dp, bottomEnd = 10.dp)
                    else -> RoundedCornerShape(0.dp)
                }

                IOSSegmentedToggle(
                    text = option,
                    isSelected = isSelected,
                    shape = shape,
                    enabled = enabled,
                    size = size,
                    onClick = {
                        handleToggleGroupSelection(
                            option = option,
                            currentSelection = selectedOptions,
                            allowMultipleSelection = allowMultipleSelection,
                            allowEmptySelection = allowEmptySelection,
                            maxSelections = maxSelections,
                            onSelectionChange = onSelectionChange
                        )
                    },
                    modifier = Modifier.weight(1f)
                )
            }
        }
    }
}

@Composable
private fun IOSSegmentedToggle(
    text: String,
    isSelected: Boolean,
    shape: RoundedCornerShape,
    enabled: Boolean,
    size: TchatToggleSize,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    val buttonHeight = when (size) {
        TchatToggleSize.Small -> 32.dp
        TchatToggleSize.Medium -> 36.dp
        TchatToggleSize.Large -> 40.dp
    }

    val textStyle = when (size) {
        TchatToggleSize.Small -> TchatTypography.typography.bodySmall
        TchatToggleSize.Medium -> TchatTypography.typography.bodyMedium
        TchatToggleSize.Large -> TchatTypography.typography.bodyLarge
    }

    // iOS segmented control animations
    val backgroundColor by animateColorAsState(
        targetValue = if (isSelected) {
            Color.White // iOS selected segment is white
        } else {
            Color.Transparent
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessHigh
        ),
        label = "ios_segment_background"
    )

    val textColor by animateColorAsState(
        targetValue = if (isSelected) {
            TchatColors.primary
        } else {
            TchatColors.onSurface.copy(alpha = 0.8f)
        },
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_segment_text"
    )

    val elevation by animateFloatAsState(
        targetValue = if (isSelected) 1f else 0f,
        animationSpec = spring(
            dampingRatio = Spring.DampingRatioMediumBouncy,
            stiffness = Spring.StiffnessMedium
        ),
        label = "ios_segment_elevation"
    )

    Surface(
        modifier = modifier
            .height(buttonHeight)
            .clip(shape)
            .clickable(
                enabled = enabled,
                indication = null, // iOS doesn't use ripple
                onClick = onClick
            ),
        color = backgroundColor,
        shape = shape,
        shadowElevation = elevation.dp
    ) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = text,
                style = textStyle.copy(
                    fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Medium
                ),
                color = if (enabled) textColor else textColor.copy(alpha = 0.5f)
            )
        }
    }
}

private fun handleToggleGroupSelection(
    option: String,
    currentSelection: Set<String>,
    allowMultipleSelection: Boolean,
    allowEmptySelection: Boolean,
    maxSelections: Int?,
    onSelectionChange: (Set<String>) -> Unit
) {
    val newSelection = if (currentSelection.contains(option)) {
        // Deselecting
        if (allowEmptySelection || currentSelection.size > 1) {
            currentSelection - option
        } else {
            currentSelection // Don't allow empty selection if not allowed
        }
    } else {
        // Selecting
        if (allowMultipleSelection) {
            if (maxSelections == null || currentSelection.size < maxSelections) {
                currentSelection + option
            } else {
                currentSelection // Max selections reached
            }
        } else {
            setOf(option) // Single selection mode
        }
    }

    onSelectionChange(newSelection)
}