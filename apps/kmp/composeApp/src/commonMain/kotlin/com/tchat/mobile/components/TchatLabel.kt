package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatLabel - Cross-platform label component with expect/actual pattern
 *
 * Features:
 * - Semantic text styling for forms and UI labels
 * - Multiple style variants (Body, Caption, Overline, Required)
 * - Design system integration with TchatTypography
 * - Accessibility support with semantic roles
 * - Platform-native text rendering
 */

enum class TchatLabelStyle {
    Body,      // Standard body text for regular labels
    Caption,   // Smaller caption text for secondary information
    Overline,  // All-caps overline text for section headers
    Required   // Required field indicator styling
}

/**
 * Cross-platform label component using expect/actual pattern
 * Platform-specific implementations provide native text rendering and accessibility
 */
@Composable
expect fun TchatLabel(
    text: String,
    modifier: Modifier = Modifier,
    style: TchatLabelStyle = TchatLabelStyle.Body
)