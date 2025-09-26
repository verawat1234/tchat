package com.tchat.mobile.utils

import androidx.compose.runtime.Composable
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.unit.dp

@Composable
actual fun getWindowConfiguration(): WindowConfiguration {
    val configuration = LocalConfiguration.current
    val density = LocalDensity.current

    return WindowConfiguration(
        screenWidthDp = with(density) { configuration.screenWidthDp.dp },
        screenHeightDp = with(density) { configuration.screenHeightDp.dp }
    )
}