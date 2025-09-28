package com.tchat.mobile.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.size
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.unit.dp
import io.kamel.image.KamelImage
import io.kamel.image.asyncPainterResource

/**
 * KMP-compatible image loading component using Kamel
 * Provides cross-platform image loading with proper error handling
 */
@Composable
fun PlatformImage(
    url: String,
    contentDescription: String?,
    modifier: Modifier = Modifier,
    contentScale: ContentScale = ContentScale.Crop,
    placeholderContent: @Composable (() -> Unit)? = null,
    errorContent: @Composable (() -> Unit)? = null
) {
    KamelImage(
        resource = asyncPainterResource(url),
        contentDescription = contentDescription,
        modifier = modifier,
        contentScale = contentScale,
        onLoading = { progress ->
            placeholderContent?.invoke() ?: Box(
                modifier = Modifier.size(40.dp),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator(progress = { progress })
            }
        },
        onFailure = { exception ->
            errorContent?.invoke() ?: Box(
                modifier = Modifier.size(40.dp),
                contentAlignment = Alignment.Center
            ) {
                // Simple error placeholder
            }
        }
    )
}