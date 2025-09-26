package com.tchat.mobile.designsystem

import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

/**
 * TchatSpacing - Consistent spacing system
 * 4dp base unit system with TailwindCSS mapping
 */
object TchatSpacing {
    // Base spacing units (4dp system)
    val xs: Dp = 4.dp      // space-1 (0.25rem)
    val sm: Dp = 8.dp      // space-2 (0.5rem)
    val md: Dp = 16.dp     // space-4 (1rem)
    val lg: Dp = 24.dp     // space-6 (1.5rem)
    val xl: Dp = 32.dp     // space-8 (2rem)
    val xxl: Dp = 40.dp    // space-10 (2.5rem)

    // Component specific spacing
    val buttonPaddingHorizontal = md    // 16dp
    val buttonPaddingVertical = sm      // 8dp
    val buttonMinHeight = 44.dp         // iOS HIG compliance
    val buttonBorderRadius = sm         // 8dp
    val buttonBorderWidth = 1.dp

    // Touch targets
    val minimumTouchTarget = buttonMinHeight  // 44dp for accessibility
    val iconSize = lg                         // 24dp for icons
    val smallIconSize = 20.dp                 // 20dp for small icons

    // Card spacing
    val cardPadding = md                      // 16dp
    val cardMargin = sm                       // 8dp
    val cardElevation = 4.dp
    val cardBorderRadius = sm                 // 8dp

    // Input field spacing
    val inputPaddingHorizontal = md           // 16dp
    val inputPaddingVertical = 12.dp
    val inputMinHeight = 48.dp
    val inputBorderRadius = sm                // 8dp
}