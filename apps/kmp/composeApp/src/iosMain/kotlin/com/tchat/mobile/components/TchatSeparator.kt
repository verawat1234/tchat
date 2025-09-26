package com.tchat.mobile.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.width
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * iOS implementation of TchatSeparator with SwiftUI-inspired styling
 * Uses custom Box with background to create separator lines following iOS HIG patterns
 */
@Composable
actual fun TchatSeparator(
    modifier: Modifier,
    orientation: TchatSeparatorOrientation
) {
    val separatorColor = TchatColors.outline.copy(alpha = 0.6f) // iOS uses slightly more transparent separators

    when (orientation) {
        TchatSeparatorOrientation.Horizontal -> {
            Box(
                modifier = modifier
                    .fillMaxWidth()
                    .height(0.5.dp) // iOS uses thinner separators (0.5dp instead of 1dp)
                    .background(separatorColor)
            )
        }
        TchatSeparatorOrientation.Vertical -> {
            Box(
                modifier = modifier
                    .fillMaxHeight()
                    .width(0.5.dp) // iOS uses thinner separators (0.5dp instead of 1dp)
                    .background(separatorColor)
            )
        }
    }
}