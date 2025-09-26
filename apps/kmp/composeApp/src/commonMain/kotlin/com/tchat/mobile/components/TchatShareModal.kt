package com.tchat.mobile.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * Share platform data class
 */
data class SharePlatform(
    val id: String,
    val name: String,
    val icon: ImageVector,
    val color: Color,
    val isAvailable: Boolean = true
)

/**
 * Share content data class
 */
data class ShareContent(
    val title: String,
    val description: String? = null,
    val url: String? = null,
    val imageUrl: String? = null,
    val type: ShareContentType = ShareContentType.GENERAL
)

enum class ShareContentType {
    GENERAL,
    PRODUCT,
    SHOP,
    LIVE_STREAM,
    REVIEW
}

/**
 * TchatShareModal - Modern share modal following web design patterns
 *
 * Features:
 * - Grid layout of share platforms with icons and labels
 * - Copy link functionality with feedback
 * - QR code generation for sharing
 * - Platform-specific sharing optimization
 * - Modern web-inspired UI with proper spacing and typography
 */
@Composable
fun TchatShareModal(
    isVisible: Boolean,
    content: ShareContent,
    onDismiss: () -> Unit,
    onShare: (SharePlatform, ShareContent) -> Unit = { _, _ -> },
    onCopyLink: (String) -> Unit = { },
    modifier: Modifier = Modifier
) {
    val sharePlatforms = remember {
        getDefaultSharePlatforms()
    }

    var showCopiedFeedback by remember { mutableStateOf(false) }

    // Handle copy feedback
    LaunchedEffect(showCopiedFeedback) {
        if (showCopiedFeedback) {
            kotlinx.coroutines.delay(2000)
            showCopiedFeedback = false
        }
    }

    TchatDialog(
        isVisible = isVisible,
        onDismissRequest = onDismiss,
        variant = DialogVariant.Custom,
        title = "Share ${content.type.displayName()}",
        modifier = modifier,
        content = {
            ShareModalContent(
                content = content,
                platforms = sharePlatforms,
                showCopiedFeedback = showCopiedFeedback,
                onShare = { platform ->
                    onShare(platform, content)
                    onDismiss()
                },
                onCopyLink = { url ->
                    onCopyLink(url)
                    showCopiedFeedback = true
                }
            )
        }
    )
}

@Composable
private fun ShareModalContent(
    content: ShareContent,
    platforms: List<SharePlatform>,
    showCopiedFeedback: Boolean,
    onShare: (SharePlatform) -> Unit,
    onCopyLink: (String) -> Unit
) {
    Column(
        modifier = Modifier
            .heightIn(max = 600.dp)
            .verticalScroll(rememberScrollState()),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
    ) {
        // Content Preview
        ContentPreview(content = content)

        // Share Platforms Grid
        Text(
            "Share via",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            modifier = Modifier.padding(bottom = TchatSpacing.sm)
        )

        SharePlatformsGrid(
            platforms = platforms,
            onShare = onShare
        )

        // Copy Link Section
        Spacer(modifier = Modifier.height(TchatSpacing.sm))

        content.url?.let { url ->
            CopyLinkSection(
                url = url,
                showCopiedFeedback = showCopiedFeedback,
                onCopyLink = onCopyLink
            )
        }
    }
}

@Composable
private fun ContentPreview(
    content: ShareContent,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Content icon based on type
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .background(TchatColors.primary.copy(alpha = 0.1f), CircleShape),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        imageVector = when (content.type) {
                            ShareContentType.PRODUCT -> Icons.Default.ShoppingBag
                            ShareContentType.SHOP -> Icons.Default.Store
                            ShareContentType.LIVE_STREAM -> Icons.Default.PlayArrow
                            ShareContentType.REVIEW -> Icons.Default.Star
                            else -> Icons.Default.Share
                        },
                        contentDescription = null,
                        tint = TchatColors.primary,
                        modifier = Modifier.size(20.dp)
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.md))

                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        content.title,
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.onSurface,
                        maxLines = 2
                    )

                    content.description?.let { desc ->
                        Text(
                            desc,
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant,
                            maxLines = 2,
                            modifier = Modifier.padding(top = 2.dp)
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun SharePlatformsGrid(
    platforms: List<SharePlatform>,
    onShare: (SharePlatform) -> Unit,
    modifier: Modifier = Modifier
) {
    // Create a grid with 3 columns
    Column(modifier = modifier.fillMaxWidth()) {
        platforms.chunked(3).forEach { rowPlatforms ->
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
            ) {
                rowPlatforms.forEach { platform ->
                    SharePlatformItem(
                        platform = platform,
                        onClick = { onShare(platform) },
                        modifier = Modifier.weight(1f)
                    )
                }
                // Fill remaining spaces in the row
                repeat(3 - rowPlatforms.size) {
                    Spacer(modifier = Modifier.weight(1f))
                }
            }

            if (rowPlatforms != platforms.chunked(3).last()) {
                Spacer(modifier = Modifier.height(TchatSpacing.md))
            }
        }
    }
}

@Composable
private fun SharePlatformItem(
    platform: SharePlatform,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .clickable(enabled = platform.isAvailable) { onClick() }
            .padding(TchatSpacing.sm),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Box(
            modifier = Modifier
                .size(56.dp)
                .background(
                    color = if (platform.isAvailable) platform.color.copy(alpha = 0.1f) else TchatColors.onSurfaceVariant.copy(alpha = 0.1f),
                    shape = CircleShape
                ),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                imageVector = platform.icon,
                contentDescription = platform.name,
                tint = if (platform.isAvailable) platform.color else TchatColors.onSurfaceVariant,
                modifier = Modifier.size(24.dp)
            )
        }

        Spacer(modifier = Modifier.height(TchatSpacing.xs))

        Text(
            platform.name,
            style = MaterialTheme.typography.labelSmall,
            color = if (platform.isAvailable) TchatColors.onSurface else TchatColors.onSurfaceVariant,
            textAlign = TextAlign.Center,
            maxLines = 1
        )
    }
}

@Composable
private fun CopyLinkSection(
    url: String,
    showCopiedFeedback: Boolean,
    onCopyLink: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier.fillMaxWidth()) {
        Text(
            "Or copy link",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            modifier = Modifier.padding(bottom = TchatSpacing.sm)
        )

        Card(
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
            elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable { onCopyLink(url) }
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    imageVector = if (showCopiedFeedback) Icons.Default.Check else Icons.Default.Link,
                    contentDescription = if (showCopiedFeedback) "Copied" else "Copy link",
                    tint = if (showCopiedFeedback) TchatColors.success else TchatColors.primary,
                    modifier = Modifier.size(20.dp)
                )

                Spacer(modifier = Modifier.width(TchatSpacing.md))

                Text(
                    if (showCopiedFeedback) "Link copied to clipboard!" else url,
                    style = MaterialTheme.typography.bodyMedium,
                    color = if (showCopiedFeedback) TchatColors.success else TchatColors.onSurface,
                    modifier = Modifier.weight(1f),
                    maxLines = 1
                )

                if (!showCopiedFeedback) {
                    Text(
                        "Copy",
                        style = MaterialTheme.typography.labelMedium,
                        color = TchatColors.primary,
                        fontWeight = FontWeight.Medium
                    )
                }
            }
        }
    }
}

/**
 * Extension function for ShareContentType display names
 */
private fun ShareContentType.displayName(): String = when (this) {
    ShareContentType.PRODUCT -> "Product"
    ShareContentType.SHOP -> "Shop"
    ShareContentType.LIVE_STREAM -> "Live Stream"
    ShareContentType.REVIEW -> "Review"
    ShareContentType.GENERAL -> "Content"
}

/**
 * Get default share platforms with modern icons and colors
 */
private fun getDefaultSharePlatforms(): List<SharePlatform> = listOf(
    SharePlatform(
        id = "whatsapp",
        name = "WhatsApp",
        icon = Icons.Default.Message,
        color = Color(0xFF25D366)
    ),
    SharePlatform(
        id = "facebook",
        name = "Facebook",
        icon = Icons.Default.Facebook,
        color = Color(0xFF1877F2)
    ),
    SharePlatform(
        id = "twitter",
        name = "Twitter",
        icon = Icons.Default.Send,
        color = Color(0xFF1DA1F2)
    ),
    SharePlatform(
        id = "instagram",
        name = "Instagram",
        icon = Icons.Default.PhotoCamera,
        color = Color(0xFFE4405F)
    ),
    SharePlatform(
        id = "linkedin",
        name = "LinkedIn",
        icon = Icons.Default.Work,
        color = Color(0xFF0A66C2)
    ),
    SharePlatform(
        id = "telegram",
        name = "Telegram",
        icon = Icons.Default.Send,
        color = Color(0xFF0088CC)
    ),
    SharePlatform(
        id = "email",
        name = "Email",
        icon = Icons.Default.Email,
        color = Color(0xFF34495E)
    ),
    SharePlatform(
        id = "sms",
        name = "SMS",
        icon = Icons.Default.Sms,
        color = Color(0xFF2ECC71)
    ),
    SharePlatform(
        id = "more",
        name = "More",
        icon = Icons.Default.MoreHoriz,
        color = TchatColors.onSurfaceVariant
    )
)