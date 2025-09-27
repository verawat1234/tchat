package com.tchat.mobile.components

import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.isImeVisible
import androidx.compose.material3.BottomSheetDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties

/**
 * Android implementation of PlatformDialog
 * Uses standard Dialog composable which works well on Android
 */
@Composable
actual fun PlatformDialog(
    visible: Boolean,
    onDismiss: () -> Unit,
    content: @Composable () -> Unit
) {
    if (visible) {
        Dialog(
            onDismissRequest = onDismiss,
            properties = DialogProperties(
                dismissOnBackPress = true,
                dismissOnClickOutside = true,
                usePlatformDefaultWidth = false
            )
        ) {
            content()
        }
    }
}

/**
 * Android implementation of PlatformSheet
 * Uses Material3 ModalBottomSheet with proper Android behavior
 */
@OptIn(ExperimentalMaterial3Api::class, ExperimentalLayoutApi::class)
@Composable
actual fun PlatformSheet(
    visible: Boolean,
    onDismiss: () -> Unit,
    modifier: Modifier,
    content: @Composable () -> Unit
) {
    if (visible) {
        val sheetState = rememberModalBottomSheetState(
            skipPartiallyExpanded = false
        )

        // Handle keyboard visibility
        val isImeVisible = WindowInsets.isImeVisible

        LaunchedEffect(isImeVisible) {
            if (isImeVisible) {
                sheetState.expand()
            }
        }

        ModalBottomSheet(
            onDismissRequest = onDismiss,
            sheetState = sheetState,
            modifier = modifier,
            dragHandle = { BottomSheetDefaults.DragHandle() }
        ) {
            content()
        }
    }
}

/**
 * Android implementation of PlatformConfiguration
 * Uses LocalConfiguration which is Android-specific
 */
actual object PlatformConfiguration {
    actual val screenWidth: Int
        @Composable get() = LocalConfiguration.current.screenWidthDp

    actual val screenHeight: Int
        @Composable get() = LocalConfiguration.current.screenHeightDp

    actual val density: Float
        @Composable get() = LocalDensity.current.density
}