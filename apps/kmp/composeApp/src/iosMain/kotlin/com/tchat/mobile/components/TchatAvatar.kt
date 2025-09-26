package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.designsystem.TchatTypography

/**
 * iOS implementation of TchatAvatar using Compose Multiplatform
 * Provides iOS-native styling with subtle shadows and refined typography
 *
 * Features iOS-specific behavior:
 * - Subtle drop shadows for depth
 * - SF Pro typography weights
 * - iOS HIG-compliant status indicator positioning
 * - Refined border styling
 * - Slightly larger touch targets
 */
@Composable
actual fun TchatAvatar(
    modifier: Modifier,
    size: TchatAvatarSize,
    imageUrl: String?,
    name: String,
    status: TchatAvatarStatus,
    loading: Boolean,
    onClick: (() -> Unit)?,
    contentDescription: String?
) {
    // iOS-specific size configuration
    val (avatarSize, textStyle, statusSize, borderWidth, shadowElevation) = when (size) {
        TchatAvatarSize.XS -> IOSAvatarSizeConfig(
            24.dp,
            TchatTypography.typography.labelSmall.copy(
                fontSize = 9.sp,
                fontWeight = FontWeight.SemiBold // iOS uses more weight
            ),
            7.dp, // Slightly larger for iOS
            0.5.dp,
            0.5.dp
        )
        TchatAvatarSize.SM -> IOSAvatarSizeConfig(
            32.dp,
            TchatTypography.typography.labelMedium.copy(
                fontSize = 11.sp,
                fontWeight = FontWeight.SemiBold
            ),
            9.dp,
            1.dp,
            1.dp
        )
        TchatAvatarSize.MD -> IOSAvatarSizeConfig(
            40.dp,
            TchatTypography.typography.labelLarge.copy(
                fontSize = 13.sp,
                fontWeight = FontWeight.Bold
            ),
            11.dp,
            1.5.dp,
            1.5.dp
        )
        TchatAvatarSize.LG -> IOSAvatarSizeConfig(
            48.dp,
            TchatTypography.typography.headlineSmall.copy(
                fontSize = 17.sp, // iOS prefers 17sp for large text
                fontWeight = FontWeight.Bold
            ),
            13.dp,
            2.dp,
            2.dp
        )
        TchatAvatarSize.XL -> IOSAvatarSizeConfig(
            64.dp,
            TchatTypography.typography.headlineMedium.copy(
                fontSize = 22.sp,
                fontWeight = FontWeight.Black // iOS uses black weight for large avatars
            ),
            17.dp,
            2.5.dp,
            3.dp
        )
    }

    val statusColor = getIOSStatusColor(status)
    val initials = generateIOSInitials(name)

    Box(
        modifier = modifier
            .size(avatarSize)
            .then(
                if (onClick != null) {
                    Modifier.clickable { onClick() }
                } else Modifier
            )
            .semantics {
                this.contentDescription = contentDescription ?: when {
                    name.isNotEmpty() -> "$name's avatar"
                    else -> "User avatar"
                }
            }
    ) {
        if (loading) {
            // iOS-style loading skeleton with subtle shadow
            TchatSkeleton(
                modifier = Modifier
                    .size(avatarSize)
                    .shadow(shadowElevation, CircleShape)
                    .clip(CircleShape)
            )
        } else {
            // Main avatar content with iOS shadow
            Box(
                modifier = Modifier
                    .size(avatarSize)
                    .shadow(shadowElevation, CircleShape)
                    .clip(CircleShape)
                    .background(
                        color = if (imageUrl != null) Color.Transparent else TchatColors.surface.copy(alpha = 0.95f),
                        shape = CircleShape
                    )
                    .border(
                        borderWidth,
                        TchatColors.outline.copy(alpha = 0.3f), // More subtle iOS border
                        CircleShape
                    ),
                contentAlignment = Alignment.Center
            ) {
                // For now, always show initials fallback (image loading can be added later with proper dependency)
                Text(
                    text = initials,
                    style = textStyle,
                    color = TchatColors.onSurface.copy(alpha = 0.8f),
                    textAlign = TextAlign.Center
                )
            }
        }

        // iOS-style status indicator with more subtle styling
        if (status != TchatAvatarStatus.None && !loading) {
            val animatedStatusColor by animateColorAsState(
                targetValue = statusColor,
                animationSpec = tween(400), // Slightly longer for iOS
                label = "ios_status_color"
            )

            Box(
                modifier = Modifier
                    .align(Alignment.BottomEnd)
                    .offset(x = 2.dp, y = 2.dp) // Slight offset for iOS positioning
                    .size(statusSize)
                    .shadow(1.dp, CircleShape) // Shadow for iOS depth
                    .background(animatedStatusColor, CircleShape)
                    .border(
                        2.5.dp, // Slightly thicker border for iOS
                        TchatColors.background,
                        CircleShape
                    )
            )
        }
    }
}

@Composable
private fun getIOSStatusColor(status: TchatAvatarStatus): Color {
    return when (status) {
        TchatAvatarStatus.None -> Color.Transparent
        TchatAvatarStatus.Online -> TchatColors.success.copy(alpha = 0.9f)
        TchatAvatarStatus.Offline -> TchatColors.onSurfaceVariant.copy(alpha = 0.7f)
        TchatAvatarStatus.Busy -> TchatColors.error.copy(alpha = 0.9f)
        TchatAvatarStatus.Away -> TchatColors.warning.copy(alpha = 0.9f)
    }
}

private fun generateIOSInitials(name: String): String {
    if (name.isEmpty()) return "?"

    val words = name.trim().split(" ").filter { it.isNotEmpty() }
    return when {
        words.isEmpty() -> "?"
        words.size == 1 -> words[0].take(2).uppercase()
        else -> "${words.first().first()}${words.last().first()}".uppercase()
    }
}

// Helper data class for iOS avatar configuration
private data class IOSAvatarSizeConfig(
    val avatarSize: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val statusSize: androidx.compose.ui.unit.Dp,
    val borderWidth: androidx.compose.ui.unit.Dp,
    val shadowElevation: androidx.compose.ui.unit.Dp
)