package com.tchat.mobile.components

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatToast using Material3 Snackbar
 */
@Composable
actual fun TchatToast(
    message: String,
    variant: TchatToastVariant,
    position: TchatToastPosition,
    modifier: Modifier,
    duration: Long,
    dismissible: Boolean,
    onDismiss: (() -> Unit)?,
    action: ToastAction?,
    icon: (@Composable () -> Unit)?
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .padding(16.dp),
        shape = RoundedCornerShape(8.dp),
        colors = CardDefaults.cardColors(
            containerColor = when (variant) {
                TchatToastVariant.Success -> TchatColors.success
                TchatToastVariant.Warning -> TchatColors.warning
                TchatToastVariant.Error -> TchatColors.error
                else -> TchatColors.surface
            }
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            icon?.let {
                it()
                Spacer(modifier = Modifier.width(8.dp))
            }

            Text(
                text = message,
                style = MaterialTheme.typography.bodyMedium,
                modifier = Modifier.weight(1f)
            )

            action?.let {
                Spacer(modifier = Modifier.width(8.dp))
                TextButton(onClick = it.onClick) {
                    Text(it.label)
                }
            }
        }
    }
}

/**
 * Android implementation of TchatToastManager
 */
actual class TchatToastManager {
    actual fun showToast(
        message: String,
        variant: TchatToastVariant,
        position: TchatToastPosition,
        duration: Long,
        action: ToastAction?
    ) {
        // Implementation for showing toast
    }

    actual fun clearAll() {
        // Implementation for clearing all toasts
    }

    actual fun clearVariant(variant: TchatToastVariant) {
        // Implementation for clearing toasts of specific variant
    }
}