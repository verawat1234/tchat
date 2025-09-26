package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatSeparator - Cross-platform separator component with expect/actual pattern
 *
 * Features:
 * - Horizontal and vertical orientation support
 * - Configurable thickness and color
 * - Platform-native styling (Material3 Divider on Android, SwiftUI Divider on iOS)
 * - Design system integration with TchatColors
 */

enum class TchatSeparatorOrientation {
    Horizontal,  // Horizontal divider line
    Vertical     // Vertical divider line
}

/**
 * Cross-platform separator component using expect/actual pattern
 * Platform-specific implementations provide native divider styling
 */
@Composable
expect fun TchatSeparator(
    modifier: Modifier = Modifier,
    orientation: TchatSeparatorOrientation = TchatSeparatorOrientation.Horizontal
)