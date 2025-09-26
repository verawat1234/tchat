package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.infiniteRepeatable
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
 * Android implementation of TchatAvatar using Jetpack Compose
 * Uses Coil for image loading and Material3 styling
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
    // Size configuration
    val (avatarSize, textStyle, statusSize, borderWidth) = when (size) {
        TchatAvatarSize.XS -> AvatarSizeConfig(
            24.dp,
            TchatTypography.typography.labelSmall.copy(
                fontSize = 10.sp,
                fontWeight = FontWeight.Medium
            ),
            6.dp,
            1.dp
        )
        TchatAvatarSize.SM -> AvatarSizeConfig(
            32.dp,
            TchatTypography.typography.labelMedium.copy(
                fontSize = 12.sp,
                fontWeight = FontWeight.Medium
            ),
            8.dp,
            1.5.dp
        )
        TchatAvatarSize.MD -> AvatarSizeConfig(
            40.dp,
            TchatTypography.typography.labelLarge.copy(
                fontSize = 14.sp,
                fontWeight = FontWeight.SemiBold
            ),
            10.dp,
            2.dp
        )
        TchatAvatarSize.LG -> AvatarSizeConfig(
            48.dp,
            TchatTypography.typography.headlineSmall.copy(
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold
            ),
            12.dp,
            2.dp
        )
        TchatAvatarSize.XL -> AvatarSizeConfig(
            64.dp,
            TchatTypography.typography.headlineMedium.copy(
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold
            ),
            16.dp,
            3.dp
        )
    }

    val statusColor = getAndroidStatusColor(status)
    val initials = generateInitials(name)

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
            // Loading skeleton with shimmer effect
            TchatSkeleton(
                modifier = Modifier
                    .size(avatarSize)
                    .clip(CircleShape)
            )
        } else {
            // Main avatar content
            Box(
                modifier = Modifier
                    .size(avatarSize)
                    .clip(CircleShape)
                    .background(
                        color = if (imageUrl != null) Color.Transparent else TchatColors.surface,
                        shape = CircleShape
                    )
                    .border(borderWidth, TchatColors.outline, CircleShape),
                contentAlignment = Alignment.Center
            ) {
                // For now, always show initials fallback (image loading can be added later with proper dependency)
                Text(
                    text = initials,
                    style = textStyle,
                    color = TchatColors.onSurface,
                    textAlign = TextAlign.Center
                )
            }
        }

        // Status indicator
        if (status != TchatAvatarStatus.None && !loading) {
            val animatedStatusColor by animateColorAsState(
                targetValue = statusColor,
                animationSpec = tween(300),
                label = "status_color"
            )

            Box(
                modifier = Modifier
                    .align(Alignment.BottomEnd)
                    .size(statusSize)
                    .background(animatedStatusColor, CircleShape)
                    .border(2.dp, TchatColors.background, CircleShape)
            )
        }
    }
}

@Composable
private fun getAndroidStatusColor(status: TchatAvatarStatus): Color {
    return when (status) {
        TchatAvatarStatus.None -> Color.Transparent
        TchatAvatarStatus.Online -> TchatColors.success
        TchatAvatarStatus.Offline -> TchatColors.onSurfaceVariant
        TchatAvatarStatus.Busy -> TchatColors.error
        TchatAvatarStatus.Away -> TchatColors.warning
    }
}

private fun generateInitials(name: String): String {
    if (name.isEmpty()) return "?"

    val words = name.trim().split(" ").filter { it.isNotEmpty() }
    return when {
        words.isEmpty() -> "?"
        words.size == 1 -> words[0].take(2).uppercase()
        else -> "${words.first().first()}${words.last().first()}".uppercase()
    }
}

// Helper data class for Android avatar configuration
private data class AvatarSizeConfig(
    val avatarSize: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val statusSize: androidx.compose.ui.unit.Dp,
    val borderWidth: androidx.compose.ui.unit.Dp
)