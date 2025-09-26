package com.tchat.mobile.components

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ErrorOutline
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * Reusable empty state component for error screens, not found pages, and empty content
 *
 * @param icon The icon to display (defaults to ErrorOutline)
 * @param title The main title text
 * @param message The descriptive message text
 * @param modifier Modifier for the root container
 * @param action Optional action button or content
 */
@Composable
fun TchatEmptyState(
    title: String,
    message: String,
    modifier: Modifier = Modifier,
    icon: ImageVector = Icons.Default.ErrorOutline,
    action: @Composable (() -> Unit)? = null
) {
    Box(
        modifier = modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                modifier = Modifier.size(64.dp),
                tint = TchatColors.onSurfaceVariant
            )

            Text(
                text = title,
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold,
                color = TchatColors.onSurface
            )

            Text(
                text = message,
                style = MaterialTheme.typography.bodyLarge,
                color = TchatColors.onSurfaceVariant
            )

            action?.invoke()
        }
    }
}

/**
 * Convenience composable for "Not Found" states
 */
@Composable
fun TchatNotFoundState(
    itemType: String,
    modifier: Modifier = Modifier,
    action: @Composable (() -> Unit)? = null
) {
    TchatEmptyState(
        title = "$itemType Not Found",
        message = "The $itemType you're looking for doesn't exist.",
        modifier = modifier,
        action = action
    )
}

/**
 * Convenience composable for empty list states
 */
@Composable
fun TchatEmptyListState(
    title: String,
    message: String,
    icon: ImageVector,
    modifier: Modifier = Modifier,
    action: @Composable (() -> Unit)? = null
) {
    TchatEmptyState(
        title = title,
        message = message,
        icon = icon,
        modifier = modifier,
        action = action
    )
}