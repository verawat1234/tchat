package com.tchat.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import com.tchat.designsystem.Typography

/**
 * Primary button component following Tchat design system
 */
@Composable
fun TchatButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    variant: TchatButtonVariant = TchatButtonVariant.Primary,
    size: TchatButtonSize = TchatButtonSize.Medium,
    enabled: Boolean = true,
    loading: Boolean = false,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() }
) {
    val isPressed by interactionSource.collectIsPressedAsState()
    val scale by animateFloatAsState(
        targetValue = if (isPressed) 0.95f else 1f,
        label = "button_scale"
    )

    val buttonColors = ButtonDefaults.buttonColors(
        containerColor = variant.backgroundColor,
        contentColor = variant.contentColor,
        disabledContainerColor = variant.backgroundColor.copy(alpha = 0.6f),
        disabledContentColor = variant.contentColor.copy(alpha = 0.6f)
    )

    val buttonShape = RoundedCornerShape(Spacing.sm)

    Button(
        onClick = onClick,
        modifier = modifier
            .scale(scale)
            .height(size.height)
            .defaultMinSize(minWidth = size.minWidth),
        enabled = enabled && !loading,
        colors = buttonColors,
        shape = buttonShape,
        border = variant.borderStroke,
        contentPadding = PaddingValues(horizontal = size.horizontalPadding),
        interactionSource = interactionSource
    ) {
        Row(
            horizontalArrangement = Arrangement.spacedBy(Spacing.xs),
            verticalAlignment = Alignment.CenterVertically
        ) {
            if (loading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(16.dp),
                    color = variant.contentColor,
                    strokeWidth = 2.dp
                )
            }

            if (!loading || text.isNotEmpty()) {
                Text(
                    text = text,
                    fontSize = size.fontSize,
                    fontWeight = FontWeight.Medium,
                    color = variant.contentColor
                )
            }
        }
    }
}

/**
 * Button variant definitions
 */
enum class TchatButtonVariant(
    val backgroundColor: Color,
    val contentColor: Color,
    val borderStroke: BorderStroke? = null
) {
    Primary(
        backgroundColor = Colors.primary,
        contentColor = Colors.textOnPrimary
    ),
    Secondary(
        backgroundColor = Colors.surface,
        contentColor = Colors.textPrimary
    ),
    Ghost(
        backgroundColor = Color.Transparent,
        contentColor = Colors.primary
    ),
    Destructive(
        backgroundColor = Colors.error,
        contentColor = Colors.textOnPrimary
    ),
    Outline(
        backgroundColor = Color.Transparent,
        contentColor = Colors.primary,
        borderStroke = BorderStroke(1.dp, Colors.border)
    )
}

/**
 * Button size definitions
 */
enum class TchatButtonSize(
    val height: androidx.compose.ui.unit.Dp,
    val minWidth: androidx.compose.ui.unit.Dp,
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val horizontalPadding: androidx.compose.ui.unit.Dp
) {
    Small(
        height = 32.dp,
        minWidth = 60.dp,
        fontSize = 14.sp,
        horizontalPadding = Spacing.sm
    ),
    Medium(
        height = 40.dp,
        minWidth = 80.dp,
        fontSize = 16.sp,
        horizontalPadding = Spacing.md
    ),
    Large(
        height = 48.dp,
        minWidth = 100.dp,
        fontSize = 18.sp,
        horizontalPadding = Spacing.lg
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatButtonPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        TchatButton(
            text = "Primary Button",
            onClick = { },
            variant = TchatButtonVariant.Primary
        )

        TchatButton(
            text = "Secondary Button",
            onClick = { },
            variant = TchatButtonVariant.Secondary
        )

        TchatButton(
            text = "Ghost Button",
            onClick = { },
            variant = TchatButtonVariant.Ghost
        )

        TchatButton(
            text = "Outline Button",
            onClick = { },
            variant = TchatButtonVariant.Outline
        )

        TchatButton(
            text = "Destructive Button",
            onClick = { },
            variant = TchatButtonVariant.Destructive
        )

        TchatButton(
            text = "Loading...",
            onClick = { },
            variant = TchatButtonVariant.Primary,
            loading = true
        )

        TchatButton(
            text = "Disabled Button",
            onClick = { },
            variant = TchatButtonVariant.Primary,
            enabled = false
        )

        Row(
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            TchatButton(
                text = "Small",
                onClick = { },
                variant = TchatButtonVariant.Primary,
                size = TchatButtonSize.Small
            )
            TchatButton(
                text = "Medium",
                onClick = { },
                variant = TchatButtonVariant.Primary,
                size = TchatButtonSize.Medium
            )
            TchatButton(
                text = "Large",
                onClick = { },
                variant = TchatButtonVariant.Primary,
                size = TchatButtonSize.Large
            )
        }
    }
}