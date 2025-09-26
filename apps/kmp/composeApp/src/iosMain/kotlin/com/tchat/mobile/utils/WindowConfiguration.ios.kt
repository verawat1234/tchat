package com.tchat.mobile.utils

import androidx.compose.runtime.Composable
import androidx.compose.ui.ExperimentalComposeUiApi
import androidx.compose.ui.platform.LocalWindowInfo
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalComposeUiApi::class)
@Composable
actual fun getWindowConfiguration(): WindowConfiguration {
    val windowInfo = LocalWindowInfo.current

    // Default reasonable screen sizes for iOS if window info is not available
    val screenWidth = if (windowInfo.containerSize.width > 0) {
        (windowInfo.containerSize.width / 3.0).dp // Approximate conversion from pixels to dp
    } else {
        375.dp // Default iPhone width
    }

    val screenHeight = if (windowInfo.containerSize.height > 0) {
        (windowInfo.containerSize.height / 3.0).dp // Approximate conversion from pixels to dp
    } else {
        812.dp // Default iPhone height
    }

    return WindowConfiguration(
        screenWidthDp = screenWidth,
        screenHeightDp = screenHeight
    )
}