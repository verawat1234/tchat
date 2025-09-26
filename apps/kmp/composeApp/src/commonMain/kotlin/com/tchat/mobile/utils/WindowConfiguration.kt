package com.tchat.mobile.utils

import androidx.compose.runtime.Composable
import androidx.compose.ui.unit.Dp

/**
 * Cross-platform window configuration
 */
data class WindowConfiguration(
    val screenWidthDp: Dp,
    val screenHeightDp: Dp
)

@Composable
expect fun getWindowConfiguration(): WindowConfiguration