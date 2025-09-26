package com.tchat.mobile.components

import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.VerticalDivider
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors

/**
 * Android implementation of TchatSeparator using Material3 Divider components
 * Provides native Material Design 3 divider styling with design system integration
 */
@Composable
actual fun TchatSeparator(
    modifier: Modifier,
    orientation: TchatSeparatorOrientation
) {
    when (orientation) {
        TchatSeparatorOrientation.Horizontal -> {
            HorizontalDivider(
                modifier = modifier,
                thickness = 1.dp,
                color = TchatColors.outline
            )
        }
        TchatSeparatorOrientation.Vertical -> {
            VerticalDivider(
                modifier = modifier,
                thickness = 1.dp,
                color = TchatColors.outline
            )
        }
    }
}