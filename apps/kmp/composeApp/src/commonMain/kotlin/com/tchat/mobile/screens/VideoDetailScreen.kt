package com.tchat.mobile.screens

import androidx.compose.animation.*
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.zIndex
import com.tchat.mobile.components.AsyncImage
import com.tchat.mobile.components.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.MockVideoRepository
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoDetailScreen(
    videoId: String,
    onNavigateBack: () -> Unit,
    onVideoClick: (String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    var video by remember { mutableStateOf<VVVideoContent?>(null) }
    var relatedVideos by remember { mutableStateOf<List<VVVideoContent>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var isLiked by remember { mutableStateOf(false) }
    var isSubscribed by remember { mutableStateOf(false) }
    var showComments by remember { mutableStateOf(false) }
    var showShareModal by remember { mutableStateOf(false) }
    var showDescription by remember { mutableStateOf(false) }
    var isPlaying by remember { mutableStateOf(true) }
    var isFullscreen by remember { mutableStateOf(false) }

    val configuration = LocalConfiguration.current
    val isLandscape = configuration.screenWidthDp > configuration.screenHeightDp

    // Repository
    val videoRepository = remember { MockVideoRepository() }
    val coroutineScope = rememberCoroutineScope()

    // Load video data
    LaunchedEffect(videoId) {
        isLoading = true
        try {
            // Get video details from mock data
            val longVideos = videoRepository.getLongVideos(VVVideoCategory.ALL)
            val shortVideos = videoRepository.getShortVideos(VVVideoCategory.ALL)
            val allVideos = longVideos + shortVideos

            video = allVideos.find { it.id == videoId }

            // Get related videos (same category + random others)
            video?.let { currentVideo ->
                val sameCategory = allVideos.filter {
                    it.category == currentVideo.category && it.id != videoId
                }.take(5)
                val others = allVideos.filter {
                    it.category != currentVideo.category && it.id != videoId
                }.shuffled().take(10 - sameCategory.size)
                relatedVideos = (sameCategory + others).shuffled()
            }
        } finally {
            isLoading = false
        }
    }

    if (isLoading || video == null) {
        Box(
            modifier = modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            CircularProgressIndicator()
        }
        return
    }

    val currentVideo = video!!

    Box(modifier = modifier.fillMaxSize().background(Color.Black)) {
        if (isLandscape || isFullscreen) {
            // Fullscreen mode - video only
            FullscreenVideoPlayer(
                video = currentVideo,
                isPlaying = isPlaying,
                onPlayPauseClick = { isPlaying = !isPlaying },
                onFullscreenToggle = { isFullscreen = !isFullscreen },
                onNavigateBack = onNavigateBack,
                modifier = Modifier.fillMaxSize()
            )
        } else {
            // Portrait mode - video with details
            Column(modifier = Modifier.fillMaxSize()) {
                // Video Player Section
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .aspectRatio(16f / 9f)
                        .background(Color.Black)
                ) {
                    // Video Player
                    TchatVideo(
                        source = VideoSource.Url(currentVideo.videoUrl ?: ""),
                        modifier = Modifier.fillMaxSize(),
                        aspectRatio = TchatVideoAspectRatio.Landscape,
                        autoPlay = isPlaying,
                        loop = false,
                        muted = false,
                        showControls = true,
                        poster = ImageSource.Url(currentVideo.thumbnail),
                        onStateChange = { state ->
                            when (state) {
                                TchatVideoState.Playing -> isPlaying = true
                                TchatVideoState.Paused -> isPlaying = false
                                TchatVideoState.Ended -> isPlaying = false
                                else -> {}
                            }
                        }
                    )

                    // Back button overlay
                    IconButton(
                        onClick = onNavigateBack,
                        modifier = Modifier
                            .align(Alignment.TopStart)
                            .padding(16.dp)
                            .background(Color.Black.copy(alpha = 0.5f), CircleShape)
                            .size(40.dp)
                    ) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = Color.White
                        )
                    }

                    // Fullscreen button overlay
                    IconButton(
                        onClick = { isFullscreen = true },
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .padding(16.dp)
                            .background(Color.Black.copy(alpha = 0.5f), CircleShape)
                            .size(40.dp)
                    ) {
                        Icon(
                            Icons.Default.Fullscreen,
                            contentDescription = "Fullscreen",
                            tint = Color.White
                        )
                    }
                }

                // Video Details Section
                LazyColumn(
                    modifier = Modifier
                        .fillMaxSize()
                        .background(TchatColors.background),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    // Video Title and Info
                    item {
                        VideoInfoSection(
                            video = currentVideo,
                            isLiked = isLiked,
                            isSubscribed = isSubscribed,
                            showDescription = showDescription,
                            onLikeClick = { isLiked = !isLiked },
                            onSubscribeClick = { isSubscribed = !isSubscribed },
                            onShareClick = { showShareModal = true },
                            onDescriptionToggle = { showDescription = !showDescription },
                            onCommentsClick = { showComments = true }
                        )
                    }

                    // Related Videos Section
                    item {
                        Text(
                            text = "Related Videos",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onBackground,
                            modifier = Modifier.padding(top = 8.dp)
                        )
                    }

                    items(relatedVideos) { relatedVideo ->
                        RelatedVideoCard(
                            video = relatedVideo,
                            onVideoClick = { onVideoClick(relatedVideo.id) }
                        )
                    }
                }
            }
        }

        // Share Modal
        if (showShareModal) {
            TchatShareModal(
                isVisible = showShareModal,
                onDismiss = { showShareModal = false },
                content = ShareContent(
                    type = ShareContentType.GENERAL,
                    title = currentVideo.title,
                    description = currentVideo.description,
                    url = "https://tchat.app/videos/${currentVideo.id}",
                    imageUrl = currentVideo.thumbnail
                )
            )
        }

        // Comments Bottom Sheet (placeholder)
        if (showComments) {
            // TODO: Implement comments bottom sheet
            LaunchedEffect(Unit) {
                showComments = false
            }
        }
    }
}

@Composable
private fun ActionButton(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    text: String,
    onClick: () -> Unit,
    isActive: Boolean = false,
    modifier: Modifier = Modifier
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.clickable { onClick() }
    ) {
        Surface(
            modifier = Modifier.size(40.dp),
            shape = CircleShape,
            color = if (isActive) TchatColors.primary.copy(alpha = 0.1f) else TchatColors.surfaceVariant
        ) {
            Icon(
                icon,
                contentDescription = text,
                tint = if (isActive) TchatColors.primary else TchatColors.onSurfaceVariant,
                modifier = Modifier
                    .fillMaxSize()
                    .padding(8.dp)
            )
        }
        Text(
            text = text,
            style = MaterialTheme.typography.labelSmall,
            color = if (isActive) TchatColors.primary else TchatColors.onSurfaceVariant,
            modifier = Modifier.padding(top = 4.dp)
        )
    }
}

@Composable
private fun CommentCard(
    comment: VideoComment,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Box(
                modifier = Modifier
                    .size(32.dp)
                    .clip(CircleShape)
                    .background(TchatColors.primary),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = comment.authorName.first().toString(),
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.onPrimary,
                    fontWeight = FontWeight.Bold
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = comment.authorName,
                        style = MaterialTheme.typography.labelMedium,
                        fontWeight = FontWeight.SemiBold,
                        color = TchatColors.onSurface
                    )
                    Spacer(modifier = Modifier.width(TchatSpacing.xs))
                    Text(
                        text = comment.timestamp,
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                Text(
                    text = comment.content,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(top = 2.dp)
                )

                Row(
                    modifier = Modifier.padding(top = TchatSpacing.xs),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    IconButton(
                        onClick = { /* Like comment */ },
                        modifier = Modifier.size(24.dp)
                    ) {
                        Icon(
                            Icons.Default.ThumbUp,
                            contentDescription = "Like",
                            modifier = Modifier.size(16.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                    }
                    Text(
                        text = comment.likesCount.toString(),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

@Composable
private fun RelatedVideoCard(
    video: RelatedVideo,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onClick,
        modifier = modifier.width(160.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(90.dp)
                    .background(TchatColors.surfaceVariant),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Filled.PlayArrow,
                    contentDescription = "Play",
                    tint = TchatColors.primary
                )

                Box(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(4.dp)
                        .background(
                            TchatColors.surface.copy(alpha = 0.9f),
                            RoundedCornerShape(2.dp)
                        )
                        .padding(horizontal = 4.dp, vertical = 2.dp)
                ) {
                    Text(
                        text = video.duration,
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurface
                    )
                }
            }

            Column(
                modifier = Modifier.padding(TchatSpacing.sm)
            ) {
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.bodySmall,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Text(
                    text = video.creatorName,
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onSurfaceVariant,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis,
                    modifier = Modifier.padding(top = 2.dp)
                )

                Text(
                    text = "${formatCount(video.views)} views",
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onSurfaceVariant,
                    modifier = Modifier.padding(top = 2.dp)
                )
            }
        }
    }
}

private fun formatCount(count: Int): String {
    return when {
        count >= 1_000_000 -> "${count / 1_000_000}M"
        count >= 1_000 -> "${count / 1_000}K"
        else -> count.toString()
    }
}

// Data models
private data class VideoDetail(
    val id: String,
    val title: String,
    val creatorName: String,
    val subscribersCount: Int,
    val duration: String,
    val views: Int,
    val uploadTime: String,
    val likesCount: Int,
    val commentsCount: Int,
    val description: String,
    val isLiked: Boolean,
    val isBookmarked: Boolean,
    val isSubscribed: Boolean,
    val topComments: List<VideoComment>,
    val relatedVideos: List<RelatedVideo>
)

private data class VideoComment(
    val id: String,
    val authorName: String,
    val content: String,
    val timestamp: String,
    val likesCount: Int
)

private data class RelatedVideo(
    val id: String,
    val title: String,
    val creatorName: String,
    val duration: String,
    val views: Int
)

@Composable
private fun FullscreenVideoPlayer(
    video: VVVideoContent,
    isPlaying: Boolean,
    onPlayPauseClick: () -> Unit,
    onFullscreenToggle: () -> Unit,
    onNavigateBack: () -> Unit,
    modifier: Modifier = Modifier
) {
    var showControls by remember { mutableStateOf(true) }
    var lastInteractionTime by remember { mutableStateOf(System.currentTimeMillis()) }

    // Auto-hide controls after 3 seconds
    LaunchedEffect(lastInteractionTime) {
        kotlinx.coroutines.delay(3000)
        val currentTime = System.currentTimeMillis()
        if (currentTime - lastInteractionTime >= 3000) {
            showControls = false
        }
    }

    val onUserInteraction = {
        lastInteractionTime = System.currentTimeMillis()
        showControls = true
    }

    Box(modifier = modifier.fillMaxSize()) {
        // Video Player
        TchatVideo(
            source = VideoSource.Url(video.videoUrl ?: ""),
            modifier = Modifier
                .fillMaxSize()
                .clickable(
                    indication = null,
                    interactionSource = remember { MutableInteractionSource() }
                ) { onUserInteraction() },
            aspectRatio = TchatVideoAspectRatio.Landscape,
            autoPlay = isPlaying,
            loop = false,
            muted = false,
            showControls = false, // We'll show our custom controls
            poster = ImageSource.Url(video.thumbnail),
            onStateChange = { state ->
                // Handle video state changes
            }
        )

        // Custom Controls Overlay
        AnimatedVisibility(
            visible = showControls,
            enter = fadeIn(),
            exit = fadeOut(),
            modifier = Modifier.fillMaxSize()
        ) {
            Box(modifier = Modifier.fillMaxSize()) {
                // Top Controls
                Row(
                    modifier = Modifier
                        .align(Alignment.TopStart)
                        .fillMaxWidth()
                        .background(
                            androidx.compose.ui.graphics.Brush.verticalGradient(
                                colors = listOf(Color.Black.copy(alpha = 0.7f), Color.Transparent)
                            )
                        )
                        .padding(16.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    IconButton(
                        onClick = onNavigateBack,
                        modifier = Modifier.size(40.dp)
                    ) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = Color.White,
                            modifier = Modifier.size(24.dp)
                        )
                    }

                    Text(
                        text = video.title,
                        color = Color.White,
                        style = MaterialTheme.typography.titleMedium,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis,
                        modifier = Modifier.weight(1f).padding(horizontal = 16.dp)
                    )

                    IconButton(
                        onClick = onFullscreenToggle,
                        modifier = Modifier.size(40.dp)
                    ) {
                        Icon(
                            Icons.Default.FullscreenExit,
                            contentDescription = "Exit Fullscreen",
                            tint = Color.White,
                            modifier = Modifier.size(24.dp)
                        )
                    }
                }

                // Center Play/Pause Button
                IconButton(
                    onClick = {
                        onPlayPauseClick()
                        onUserInteraction()
                    },
                    modifier = Modifier
                        .align(Alignment.Center)
                        .size(80.dp)
                        .background(Color.Black.copy(alpha = 0.5f), CircleShape)
                ) {
                    Icon(
                        if (isPlaying) Icons.Default.Pause else Icons.Default.PlayArrow,
                        contentDescription = if (isPlaying) "Pause" else "Play",
                        tint = Color.White,
                        modifier = Modifier.size(40.dp)
                    )
                }

                // Bottom Controls (placeholder for progress bar, volume, etc.)
                Row(
                    modifier = Modifier
                        .align(Alignment.BottomStart)
                        .fillMaxWidth()
                        .background(
                            androidx.compose.ui.graphics.Brush.verticalGradient(
                                colors = listOf(Color.Transparent, Color.Black.copy(alpha = 0.7f))
                            )
                        )
                        .padding(16.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "00:00", // TODO: Show actual current time
                        color = Color.White,
                        style = MaterialTheme.typography.bodySmall
                    )

                    Box(
                        modifier = Modifier
                            .weight(1f)
                            .height(4.dp)
                            .padding(horizontal = 16.dp)
                            .background(Color.White.copy(alpha = 0.3f), RoundedCornerShape(2.dp))
                    ) {
                        // TODO: Implement actual progress bar
                        Box(
                            modifier = Modifier
                                .fillMaxWidth(0.3f) // Mock 30% progress
                                .fillMaxHeight()
                                .background(TchatColors.primary, RoundedCornerShape(2.dp))
                        )
                    }

                    Text(
                        text = video.duration,
                        color = Color.White,
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            }
        }
    }
}

@Composable
private fun VideoInfoSection(
    video: VVVideoContent,
    isLiked: Boolean,
    isSubscribed: Boolean,
    showDescription: Boolean,
    onLikeClick: () -> Unit,
    onSubscribeClick: () -> Unit,
    onShareClick: () -> Unit,
    onDescriptionToggle: () -> Unit,
    onCommentsClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier.fillMaxWidth()) {
        // Video Title
        Text(
            text = video.title,
            style = MaterialTheme.typography.titleLarge,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onBackground,
            modifier = Modifier.padding(bottom = 8.dp)
        )

        // Views and Upload Time
        Text(
            text = "${VVideoMockData.formatViews(video.views)} views • ${video.uploadTimeFormatted}",
            style = MaterialTheme.typography.bodyMedium,
            color = TchatColors.onSurfaceVariant,
            modifier = Modifier.padding(bottom = 16.dp)
        )

        // Action Buttons Row
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            // Like Button
            ActionButton(
                icon = if (isLiked) Icons.Default.ThumbUp else Icons.Default.ThumbUpOffAlt,
                text = VVideoMockData.formatViews(video.likes),
                onClick = onLikeClick,
                isActive = isLiked
            )

            // Dislike Button
            ActionButton(
                icon = Icons.Default.ThumbDownOffAlt,
                text = "Dislike",
                onClick = { /* TODO */ }
            )

            // Share Button
            ActionButton(
                icon = Icons.Default.Share,
                text = "Share",
                onClick = onShareClick
            )

            // Download Button
            ActionButton(
                icon = Icons.Default.Download,
                text = "Download",
                onClick = { /* TODO */ }
            )
        }

        Divider(modifier = Modifier.padding(vertical = 16.dp))

        // Channel Info
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            AsyncImage(
                model = video.channel.avatar,
                contentDescription = "${video.channel.name} avatar",
                modifier = Modifier
                    .size(40.dp)
                    .clip(CircleShape),
                contentScale = ContentScale.Crop
            )

            Spacer(modifier = Modifier.width(12.dp))

            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        text = video.channel.name,
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.onBackground
                    )
                    if (video.channel.verified) {
                        Icon(
                            Icons.Default.Verified,
                            contentDescription = "Verified",
                            tint = Color.Blue,
                            modifier = Modifier.size(16.dp).padding(start = 4.dp)
                        )
                    }
                }
                Text(
                    text = "${VVideoMockData.formatSubscribers(video.channel.subscribers)} subscribers",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            // Subscribe Button
            TchatButton(
                text = if (isSubscribed) "Subscribed" else "Subscribe",
                onClick = onSubscribeClick,
                variant = if (isSubscribed) TchatButtonVariant.Outline else TchatButtonVariant.Primary,
                modifier = Modifier.height(36.dp)
            )
        }

        // Description Section
        if (video.description.isNotEmpty()) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(top = 16.dp)
                    .clickable { onDescriptionToggle() },
                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
            ) {
                Column(modifier = Modifier.padding(12.dp)) {
                    Text(
                        text = if (showDescription) video.description else
                               video.description.take(100) + if (video.description.length > 100) "..." else "",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        maxLines = if (showDescription) Int.MAX_VALUE else 2,
                        overflow = TextOverflow.Ellipsis
                    )

                    Text(
                        text = if (showDescription) "Show less" else "Show more",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.primary,
                        modifier = Modifier.padding(top = 4.dp)
                    )
                }
            }
        }

        // Comments Section
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(top = 12.dp)
                .clickable { onCommentsClick() },
            colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(12.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.ChatBubbleOutline,
                    contentDescription = "Comments",
                    tint = TchatColors.onSurfaceVariant
                )
                Spacer(modifier = Modifier.width(8.dp))
                Text(
                    text = "Comments ${VVideoMockData.formatViews(video.comments)}",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
                Spacer(modifier = Modifier.weight(1f))
                Icon(
                    Icons.Default.ExpandMore,
                    contentDescription = "Expand comments",
                    tint = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun RelatedVideoCard(
    video: VVVideoContent,
    onVideoClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onVideoClick() },
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(8.dp),
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Video Thumbnail
            Box(
                modifier = Modifier
                    .size(width = 120.dp, height = 68.dp)
                    .clip(RoundedCornerShape(8.dp))
                    .background(Color.Black)
            ) {
                AsyncImage(
                    model = video.thumbnail,
                    contentDescription = "Video thumbnail",
                    modifier = Modifier.fillMaxSize(),
                    contentScale = ContentScale.Crop
                )

                // Duration Badge
                Surface(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(4.dp),
                    color = Color.Black.copy(alpha = 0.8f),
                    shape = RoundedCornerShape(2.dp)
                ) {
                    Text(
                        text = video.duration,
                        color = Color.White,
                        style = MaterialTheme.typography.labelSmall,
                        modifier = Modifier.padding(horizontal = 4.dp, vertical = 2.dp)
                    )
                }
            }

            // Video Info
            Column(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight(),
                verticalArrangement = Arrangement.spacedBy(2.dp)
            ) {
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Text(
                    text = video.channel.name,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )

                Text(
                    text = "${VVideoMockData.formatViews(video.views)} views • ${video.uploadTimeFormatted}",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )
            }
        }
    }
}

