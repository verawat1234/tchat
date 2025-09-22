package com.tchat.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Card component following Tchat design system
 */
@Composable
fun TchatCard(
    modifier: Modifier = Modifier,
    variant: TchatCardVariant = TchatCardVariant.Elevated,
    size: TchatCardSize = TchatCardSize.Standard,
    isInteractive: Boolean = false,
    onClick: (() -> Unit)? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    content: @Composable ColumnScope.() -> Unit
) {
    val isPressed by interactionSource.collectIsPressedAsState()
    val scale by animateFloatAsState(
        targetValue = if (isInteractive && isPressed) 0.98f else 1f,
        label = "card_scale"
    )

    val shape = RoundedCornerShape(size.cornerRadius)

    Card(
        modifier = modifier
            .scale(scale)
            .then(
                if (isInteractive && onClick != null) {
                    Modifier.clickable(
                        interactionSource = interactionSource,
                        indication = null
                    ) { onClick() }
                } else Modifier
            ),
        shape = shape,
        colors = CardDefaults.cardColors(
            containerColor = variant.backgroundColor
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = variant.elevation
        ),
        border = variant.borderStroke
    ) {
        Column(
            modifier = Modifier.padding(size.contentPadding),
            content = content
        )
    }
}

/**
 * Card header component
 */
@Composable
fun TchatCardHeader(
    title: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
    leadingIcon: ImageVector? = null,
    trailingContent: (@Composable () -> Unit)? = null
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
    ) {
        if (leadingIcon != null) {
            Icon(
                imageVector = leadingIcon,
                contentDescription = null,
                tint = Colors.primary,
                modifier = Modifier.size(20.dp)
            )
        }

        Column(
            modifier = Modifier.weight(1f),
            verticalArrangement = Arrangement.spacedBy(Spacing.xs)
        ) {
            Text(
                text = title,
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold,
                color = Colors.textPrimary
            )

            if (subtitle != null) {
                Text(
                    text = subtitle,
                    fontSize = 12.sp,
                    color = Colors.textSecondary
                )
            }
        }

        trailingContent?.invoke()
    }
}

/**
 * Card footer component
 */
@Composable
fun TchatCardFooter(
    modifier: Modifier = Modifier,
    content: @Composable RowScope.() -> Unit
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(top = Spacing.sm),
        content = content
    )
}

/**
 * Card variant definitions
 */
enum class TchatCardVariant(
    val backgroundColor: Color,
    val elevation: androidx.compose.ui.unit.Dp,
    val borderStroke: BorderStroke? = null
) {
    Elevated(
        backgroundColor = Colors.cardBackground,
        elevation = 4.dp
    ),
    Outlined(
        backgroundColor = Colors.cardBackground,
        elevation = 0.dp,
        borderStroke = BorderStroke(1.dp, Colors.border)
    ),
    Filled(
        backgroundColor = Colors.surface,
        elevation = 0.dp
    ),
    Glass(
        backgroundColor = Colors.cardBackground.copy(alpha = 0.8f),
        elevation = 2.dp
    )
}

/**
 * Card size definitions
 */
enum class TchatCardSize(
    val contentPadding: androidx.compose.foundation.layout.PaddingValues,
    val cornerRadius: androidx.compose.ui.unit.Dp
) {
    Compact(
        contentPadding = PaddingValues(Spacing.sm),
        cornerRadius = Spacing.sm
    ),
    Standard(
        contentPadding = PaddingValues(Spacing.md),
        cornerRadius = Spacing.md
    ),
    Expanded(
        contentPadding = PaddingValues(Spacing.lg),
        cornerRadius = Spacing.md
    )
}

/**
 * Convenience composable for card with header
 */
@Composable
fun TchatCard(
    title: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
    leadingIcon: ImageVector? = null,
    trailingContent: (@Composable () -> Unit)? = null,
    variant: TchatCardVariant = TchatCardVariant.Elevated,
    size: TchatCardSize = TchatCardSize.Standard,
    isInteractive: Boolean = false,
    onClick: (() -> Unit)? = null,
    content: @Composable ColumnScope.() -> Unit
) {
    TchatCard(
        modifier = modifier,
        variant = variant,
        size = size,
        isInteractive = isInteractive,
        onClick = onClick
    ) {
        TchatCardHeader(
            title = title,
            subtitle = subtitle,
            leadingIcon = leadingIcon,
            trailingContent = trailingContent
        )

        Spacer(modifier = Modifier.height(Spacing.sm))

        content()
    }
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatCardPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        // Basic elevated card
        TchatCard(variant = TchatCardVariant.Elevated) {
            Column(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
                Text(
                    text = "Elevated Card",
                    fontSize = 18.sp,
                    fontWeight = FontWeight.SemiBold
                )
                Text(
                    text = "This is an elevated card with shadow",
                    color = Colors.textSecondary
                )
            }
        }

        // Outlined card with header
        TchatCard(
            title = "Card with Header",
            subtitle = "Subtitle text",
            leadingIcon = Icons.Default.Star,
            variant = TchatCardVariant.Outlined,
            isInteractive = true,
            onClick = { println("Card tapped") },
            trailingContent = {
                Icon(
                    imageVector = Icons.Default.ChevronRight,
                    contentDescription = null,
                    tint = Colors.textTertiary
                )
            }
        ) {
            Column(verticalArrangement = Arrangement.spacedBy(Spacing.xs)) {
                Text("This card has a header and is interactive")
                Text(
                    text = "Tap to interact",
                    fontSize = 12.sp,
                    color = Colors.primary
                )
            }
        }

        // Filled card
        TchatCard(
            variant = TchatCardVariant.Filled,
            size = TchatCardSize.Compact
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
            ) {
                Icon(
                    imageVector = Icons.Default.Message,
                    contentDescription = null,
                    tint = Colors.primary
                )
                Text(
                    text = "Compact filled card",
                    fontSize = 14.sp
                )
            }
        }

        // Glass card
        TchatCard(
            title = "Glass Card",
            subtitle = "Expanded size",
            leadingIcon = Icons.Default.AutoAwesome,
            variant = TchatCardVariant.Glass,
            size = TchatCardSize.Expanded
        ) {
            Column(verticalArrangement = Arrangement.spacedBy(Spacing.md)) {
                Text("This is a glass-style card with expanded padding")

                TchatCardFooter {
                    TchatButton(
                        text = "Action",
                        onClick = { println("Action tapped") },
                        variant = TchatButtonVariant.Primary,
                        size = TchatButtonSize.Small
                    )

                    Spacer(modifier = Modifier.weight(1f))

                    Text(
                        text = "Footer content",
                        fontSize = 12.sp,
                        color = Colors.textTertiary
                    )
                }
            }
        }

        // Interactive card with complex content
        TchatCard(
            variant = TchatCardVariant.Elevated,
            isInteractive = true,
            onClick = { println("Complex card tapped") }
        ) {
            Column(verticalArrangement = Arrangement.spacedBy(Spacing.md)) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column(modifier = Modifier.weight(1f)) {
                        Text(
                            text = "Interactive Card",
                            fontSize = 18.sp,
                            fontWeight = FontWeight.SemiBold
                        )
                        Text(
                            text = "With complex content",
                            fontSize = 12.sp,
                            color = Colors.textSecondary
                        )
                    }
                    Icon(
                        imageVector = Icons.Default.ChevronRight,
                        contentDescription = null,
                        tint = Colors.textTertiary
                    )
                }

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Icon(
                            imageVector = Icons.Default.CheckCircle,
                            contentDescription = null,
                            tint = Colors.success,
                            modifier = Modifier.size(16.dp)
                        )
                        Text(
                            text = "Feature 1",
                            fontSize = 12.sp,
                            color = Colors.success
                        )
                    }

                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Icon(
                            imageVector = Icons.Default.Star,
                            contentDescription = null,
                            tint = Colors.warning,
                            modifier = Modifier.size(16.dp)
                        )
                        Text(
                            text = "Feature 2",
                            fontSize = 12.sp,
                            color = Colors.warning
                        )
                    }
                }
            }
        }
    }
}