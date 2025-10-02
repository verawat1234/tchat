package com.tchat.mobile.video.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.tchat.mobile.video.models.*
import io.kamel.image.KamelImage
import io.kamel.image.asyncPainterResource
import kotlin.math.roundToInt

/**
 * Video Player Component
 *
 * Cross-platform video player with custom controls and sync support.
 * Matches web VideoPlayer.tsx functionality.
 */
@Composable
fun VideoPlayer(
    videoId: String,
    sessionId: String? = null,
    autoPlay: Boolean = false,
    controls: Boolean = true,
    loop: Boolean = false,
    muted: Boolean = false,
    quality: VideoQuality = VideoQuality.AUTO,
    onPlaybackStateChange: (PlaybackState) -> Unit = {},
    onPositionChange: (Double) -> Unit = {},
    onQualityChange: (VideoQuality) -> Unit = {},
    onFullscreenToggle: (Boolean) -> Unit = {},
    modifier: Modifier = Modifier
) {
    var playbackState by remember { mutableStateOf(PlaybackState.PAUSED) }
    var currentPosition by remember { mutableStateOf(0.0) }
    var duration by remember { mutableStateOf(0.0) }
    var volume by remember { mutableStateOf(1.0) }
    var playbackSpeed by remember { mutableStateOf(1.0) }
    var selectedQuality by remember { mutableStateOf(quality) }
    var isFullscreen by remember { mutableStateOf(false) }
    var showControls by remember { mutableStateOf(true) }
    var bufferHealth by remember { mutableStateOf(BufferHealth()) }

    Box(
        modifier = modifier
            .fillMaxWidth()
            .aspectRatio(16f / 9f)
            .background(Color.Black)
    ) {
        // Platform-specific video player will be implemented in androidMain/iosMain
        // This is the common UI structure

        // Video Player Surface (platform-specific implementation)
        Box(
            modifier = Modifier
                .fillMaxSize()
                .clickable { showControls = !showControls }
        ) {
            // Platform-specific video rendering goes here
            Text(
                text = "Video Player: $videoId",
                color = Color.White,
                modifier = Modifier.align(Alignment.Center)
            )
        }

        // Custom Controls Overlay
        if (controls && showControls) {
            VideoControls(
                playbackState = playbackState,
                currentPosition = currentPosition,
                duration = duration,
                volume = volume,
                playbackSpeed = playbackSpeed,
                quality = selectedQuality,
                isFullscreen = isFullscreen,
                bufferHealth = bufferHealth,
                onPlayPause = {
                    val newState = if (playbackState == PlaybackState.PLAYING)
                        PlaybackState.PAUSED else PlaybackState.PLAYING
                    playbackState = newState
                    onPlaybackStateChange(newState)
                },
                onSeek = { position ->
                    currentPosition = position
                    onPositionChange(position)
                },
                onVolumeChange = { newVolume -> volume = newVolume },
                onSpeedChange = { newSpeed -> playbackSpeed = newSpeed },
                onQualityChange = { newQuality ->
                    selectedQuality = newQuality
                    onQualityChange(newQuality)
                },
                onFullscreenToggle = {
                    isFullscreen = !isFullscreen
                    onFullscreenToggle(isFullscreen)
                },
                modifier = Modifier.align(Alignment.BottomCenter)
            )
        }

        // Buffering Indicator
        if (playbackState == PlaybackState.BUFFERING) {
            CircularProgressIndicator(
                modifier = Modifier.align(Alignment.Center),
                color = MaterialTheme.colorScheme.primary
            )
        }
    }
}

/**
 * Video Controls Component
 *
 * Custom video player controls with full functionality.
 */
@Composable
private fun VideoControls(
    playbackState: PlaybackState,
    currentPosition: Double,
    duration: Double,
    volume: Double,
    playbackSpeed: Double,
    quality: VideoQuality,
    isFullscreen: Boolean,
    bufferHealth: BufferHealth,
    onPlayPause: () -> Unit,
    onSeek: (Double) -> Unit,
    onVolumeChange: (Double) -> Unit,
    onSpeedChange: (Double) -> Unit,
    onQualityChange: (VideoQuality) -> Unit,
    onFullscreenToggle: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxWidth()
            .background(
                Color.Black.copy(alpha = 0.6f),
                RoundedCornerShape(topStart = 8.dp, topEnd = 8.dp)
            )
            .padding(16.dp)
    ) {
        // Progress Bar
        Slider(
            value = currentPosition.toFloat(),
            onValueChange = { onSeek(it.toDouble()) },
            valueRange = 0f..duration.toFloat(),
            modifier = Modifier.fillMaxWidth()
        )

        // Time Display
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Text(
                text = formatTime(currentPosition),
                color = Color.White,
                style = MaterialTheme.typography.bodySmall
            )
            Text(
                text = formatTime(duration),
                color = Color.White,
                style = MaterialTheme.typography.bodySmall
            )
        }

        Spacer(modifier = Modifier.height(8.dp))

        // Control Buttons
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Play/Pause Button
            IconButton(onClick = onPlayPause) {
                Icon(
                    imageVector = if (playbackState == PlaybackState.PLAYING)
                        Icons.Default.Pause else Icons.Default.PlayArrow,
                    contentDescription = if (playbackState == PlaybackState.PLAYING)
                        "Pause" else "Play",
                    tint = Color.White
                )
            }

            // Volume Control
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = { onVolumeChange(if (volume > 0) 0.0 else 1.0) }) {
                    Icon(
                        imageVector = if (volume > 0) Icons.Default.VolumeUp else Icons.Default.VolumeOff,
                        contentDescription = "Volume",
                        tint = Color.White
                    )
                }
            }

            // Quality Selector
            var showQualityMenu by remember { mutableStateOf(false) }
            Box {
                IconButton(onClick = { showQualityMenu = true }) {
                    Icon(
                        imageVector = Icons.Default.Settings,
                        contentDescription = "Quality",
                        tint = Color.White
                    )
                }
                DropdownMenu(
                    expanded = showQualityMenu,
                    onDismissRequest = { showQualityMenu = false }
                ) {
                    VideoQuality.values().forEach { q ->
                        DropdownMenuItem(
                            text = { Text(q.name) },
                            onClick = {
                                onQualityChange(q)
                                showQualityMenu = false
                            }
                        )
                    }
                }
            }

            // Fullscreen Toggle
            IconButton(onClick = onFullscreenToggle) {
                Icon(
                    imageVector = if (isFullscreen) Icons.Default.FullscreenExit else Icons.Default.Fullscreen,
                    contentDescription = "Fullscreen",
                    tint = Color.White
                )
            }
        }
    }
}

/**
 * Video Upload Component
 *
 * Video upload UI with progress tracking and metadata input.
 * Matches web VideoUpload.tsx functionality.
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoUpload(
    onUploadComplete: (String) -> Unit,
    onError: (VideoError) -> Unit,
    modifier: Modifier = Modifier
) {
    var title by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var tags by remember { mutableStateOf("") }
    var contentRating by remember { mutableStateOf(ContentRating.GENERAL) }
    var category by remember { mutableStateOf("") }
    var isMonetized by remember { mutableStateOf(false) }
    var price by remember { mutableStateOf("") }
    var uploadProgress by remember { mutableStateOf(0f) }
    var isUploading by remember { mutableStateOf(false) }

    Card(
        modifier = modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            Text(
                text = "Upload Video",
                style = MaterialTheme.typography.headlineSmall
            )

            // Title Input
            OutlinedTextField(
                value = title,
                onValueChange = { title = it },
                label = { Text("Title") },
                modifier = Modifier.fillMaxWidth(),
                enabled = !isUploading
            )

            // Description Input
            OutlinedTextField(
                value = description,
                onValueChange = { description = it },
                label = { Text("Description") },
                modifier = Modifier.fillMaxWidth(),
                minLines = 3,
                maxLines = 5,
                enabled = !isUploading
            )

            // Tags Input
            OutlinedTextField(
                value = tags,
                onValueChange = { tags = it },
                label = { Text("Tags (comma-separated)") },
                modifier = Modifier.fillMaxWidth(),
                enabled = !isUploading
            )

            // Content Rating Selector
            var showRatingMenu by remember { mutableStateOf(false) }
            ExposedDropdownMenuBox(
                expanded = showRatingMenu,
                onExpandedChange = { showRatingMenu = !showRatingMenu }
            ) {
                OutlinedTextField(
                    value = contentRating.name,
                    onValueChange = {},
                    readOnly = true,
                    label = { Text("Content Rating") },
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = showRatingMenu) },
                    modifier = Modifier
                        .fillMaxWidth()
                        .menuAnchor(),
                    enabled = !isUploading
                )
                ExposedDropdownMenu(
                    expanded = showRatingMenu,
                    onDismissRequest = { showRatingMenu = false }
                ) {
                    ContentRating.values().forEach { rating ->
                        DropdownMenuItem(
                            text = { Text(rating.name) },
                            onClick = {
                                contentRating = rating
                                showRatingMenu = false
                            }
                        )
                    }
                }
            }

            // Category Input
            OutlinedTextField(
                value = category,
                onValueChange = { category = it },
                label = { Text("Category") },
                modifier = Modifier.fillMaxWidth(),
                enabled = !isUploading
            )

            // Monetization Toggle
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text("Monetize Video")
                Switch(
                    checked = isMonetized,
                    onCheckedChange = { isMonetized = it },
                    enabled = !isUploading
                )
            }

            // Price Input (if monetized)
            if (isMonetized) {
                OutlinedTextField(
                    value = price,
                    onValueChange = { price = it },
                    label = { Text("Price (USD)") },
                    modifier = Modifier.fillMaxWidth(),
                    enabled = !isUploading
                )
            }

            // Upload Progress
            if (isUploading) {
                Column {
                    LinearProgressIndicator(
                        progress = uploadProgress,
                        modifier = Modifier.fillMaxWidth()
                    )
                    Text(
                        text = "Uploading: ${(uploadProgress * 100).roundToInt()}%",
                        style = MaterialTheme.typography.bodySmall,
                        modifier = Modifier.padding(top = 4.dp)
                    )
                }
            }

            // Upload Button
            Button(
                onClick = {
                    // Platform-specific file picker will be implemented in androidMain/iosMain
                    isUploading = true
                    // Simulate upload progress (replace with actual upload)
                    // onUploadComplete(videoId)
                },
                modifier = Modifier.fillMaxWidth(),
                enabled = !isUploading && title.isNotBlank()
            ) {
                Text(if (isUploading) "Uploading..." else "Select and Upload Video")
            }
        }
    }
}

/**
 * Video List Component
 *
 * Grid/list view of videos with filtering and sorting.
 * Matches web VideoList.tsx functionality.
 */
@Composable
fun VideoList(
    videos: List<VideoContent>,
    onVideoSelect: (VideoContent) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxWidth(),
        verticalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(16.dp)
    ) {
        items(videos) { video ->
            VideoListItem(
                video = video,
                onClick = { onVideoSelect(video) }
            )
        }
    }
}

/**
 * Video List Item Component
 *
 * Single video item in the list with thumbnail and metadata.
 */
@Composable
private fun VideoListItem(
    video: VideoContent,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(modifier = Modifier.padding(8.dp)) {
            // Thumbnail
            Box(
                modifier = Modifier
                    .width(120.dp)
                    .aspectRatio(16f / 9f)
                    .clip(RoundedCornerShape(8.dp))
                    .background(Color.Gray)
            ) {
                video.thumbnailUrl?.let { url ->
                    KamelImage(
                        resource = asyncPainterResource(url),
                        contentDescription = video.title,
                        modifier = Modifier.fillMaxSize(),
                        contentScale = ContentScale.Crop
                    )
                }

                // Duration Badge
                Text(
                    text = formatTime(video.durationSeconds.toDouble()),
                    color = Color.White,
                    style = MaterialTheme.typography.labelSmall,
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .background(Color.Black.copy(alpha = 0.7f), RoundedCornerShape(4.dp))
                        .padding(horizontal = 4.dp, vertical = 2.dp)
                )
            }

            Spacer(modifier = Modifier.width(12.dp))

            // Video Info
            Column(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight(),
                verticalArrangement = Arrangement.SpaceBetween
            ) {
                // Title
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.titleMedium,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                // Metadata
                Column(verticalArrangement = Arrangement.spacedBy(4.dp)) {
                    Text(
                        text = "${formatViews(video.viewCount)} views",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )

                    Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(4.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                imageVector = Icons.Default.ThumbUp,
                                contentDescription = "Likes",
                                modifier = Modifier.size(16.dp),
                                tint = MaterialTheme.colorScheme.onSurfaceVariant
                            )
                            Text(
                                text = formatViews(video.likeCount),
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.onSurfaceVariant
                            )
                        }

                        Row(
                            horizontalArrangement = Arrangement.spacedBy(4.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                imageVector = Icons.Default.Comment,
                                contentDescription = "Comments",
                                modifier = Modifier.size(16.dp),
                                tint = MaterialTheme.colorScheme.onSurfaceVariant
                            )
                            Text(
                                text = formatViews(video.commentCount),
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.onSurfaceVariant
                            )
                        }
                    }
                }
            }
        }
    }
}

/**
 * Helper Functions
 */

private fun formatTime(seconds: Double): String {
    val totalSeconds = seconds.toInt()
    val hours = totalSeconds / 3600
    val minutes = (totalSeconds % 3600) / 60
    val secs = totalSeconds % 60

    return if (hours > 0) {
        String.format("%d:%02d:%02d", hours, minutes, secs)
    } else {
        String.format("%d:%02d", minutes, secs)
    }
}

private fun formatViews(count: Long): String {
    return when {
        count >= 1_000_000 -> String.format("%.1fM", count / 1_000_000.0)
        count >= 1_000 -> String.format("%.1fK", count / 1_000.0)
        else -> count.toString()
    }
}