package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatAvatar - Cross-platform avatar component with expect/actual pattern
 *
 * Features:
 * - 5 size variants (XS, SM, MD, LG, XL) for different UI contexts
 * - Image loading with automatic fallback to initials
 * - Status indicator support (online, offline, busy, etc.)
 * - Loading states with skeleton animations
 * - Platform-specific circular clipping and border styling
 * - Accessibility support with descriptive labels
 * - Auto-generated initials from name strings
 */

enum class TchatAvatarSize {
    XS,   // 24dp - Small list items, compact UI
    SM,   // 32dp - Medium list items, badges
    MD,   // 40dp - Standard size for most use cases
    LG,   // 48dp - Profile sections, prominent display
    XL    // 64dp - Large profile views, settings
}

enum class TchatAvatarStatus {
    None,     // No status indicator
    Online,   // Green indicator for active users
    Offline,  // Gray indicator for offline users
    Busy,     // Red indicator for busy/do not disturb
    Away      // Yellow indicator for away status
}

/**
 * Cross-platform avatar component using expect/actual pattern
 * Platform-specific implementations handle image loading and status indicators
 */
@Composable
expect fun TchatAvatar(
    modifier: Modifier = Modifier,
    size: TchatAvatarSize = TchatAvatarSize.MD,
    imageUrl: String? = null, // URL for profile image
    name: String = "", // Name for fallback initials generation
    status: TchatAvatarStatus = TchatAvatarStatus.None,
    loading: Boolean = false, // Show loading skeleton
    onClick: (() -> Unit)? = null, // Optional click handler
    contentDescription: String? = null
)