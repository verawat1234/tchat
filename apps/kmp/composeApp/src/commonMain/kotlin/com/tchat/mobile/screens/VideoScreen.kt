package com.tchat.mobile.screens

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.zIndex
import com.tchat.mobile.components.AsyncImage
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.ShareContent
import com.tchat.mobile.components.ShareContentType
import com.tchat.mobile.components.TchatShareModal
import com.tchat.mobile.components.TchatVideo
import com.tchat.mobile.components.VideoSource
import com.tchat.mobile.components.TchatVideoAspectRatio
import com.tchat.mobile.components.ImageSource
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.MockVideoRepository
import com.tchat.mobile.repositories.VideoRepository
import com.tchat.mobile.services.NavigationService
import com.tchat.mobile.services.SharingService
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

// Video Tab Types
enum class VideoTabType {
    SHORTS,
    VIDEOS,
    HISTORY,
    SUBSCRIPTIONS
}

// Category configuration
data class CategoryConfig(
    val id: VVVideoCategory,
    val name: String,
    val icon: @Composable () -> Unit
)

// Helper functions
fun formatViewCount(views: Long): String = VVideoMockData.formatViews(views)
fun formatSubscriberCount(subs: Long): String = VVideoMockData.formatSubscribers(subs)

// Optimized Video Screen with YouTube-like single video playback
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoScreen(
    onVideoClick: (videoId: String) -> Unit = {},
    onUIVisibilityChange: (Boolean) -> Unit = {},
    onMoreClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var selectedVideoTab by remember { mutableStateOf(VideoTabType.SHORTS) }
    var currentShortIndex by remember { mutableStateOf(0) }
    var selectedCategory by remember { mutableStateOf(VVVideoCategory.ALL) }
    var likedVideos by remember { mutableStateOf(setOf<String>()) }
    var subscribedChannels by remember { mutableStateOf(setOf<String>()) }
    var showShareModal by remember { mutableStateOf(false) }
    var selectedVideoForShare by remember { mutableStateOf<VVVideoContent?>(null) }

    // Optimized UI state management - Always show tabs for better UX
    var isUIVisible by remember { mutableStateOf(true) }
    var lastInteractionTime by remember { mutableStateOf(System.currentTimeMillis()) }

    // Video player state management for single video playback
    var currentPlayingVideoId by remember { mutableStateOf<String?>(null) }
    var isVideoPlaying by remember { mutableStateOf(false) }

    // Repository and data loading
    val videoRepository = remember { MockVideoRepository() }
    var shortVideos by remember { mutableStateOf<List<VVVideoContent>>(emptyList()) }
    var longVideos by remember { mutableStateOf<List<VVVideoContent>>(emptyList()) }
    var channels by remember { mutableStateOf<List<VVChannelInfo>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }

    // Load data based on selected category
    LaunchedEffect(selectedCategory) {
        isLoading = true
        try {
            shortVideos = videoRepository.getShortVideos(selectedCategory)
            longVideos = videoRepository.getLongVideos(selectedCategory)
            channels = videoRepository.getChannels(selectedCategory)
        } finally {
            isLoading = false
        }
    }

    // Optimized UI management - Always show tabs, only hide action buttons in Shorts
    LaunchedEffect(lastInteractionTime, selectedVideoTab) {
        if (selectedVideoTab == VideoTabType.SHORTS) {
            kotlinx.coroutines.delay(3000) // 3 seconds delay
            val currentTime = System.currentTimeMillis()
            if (currentTime - lastInteractionTime >= 3000) {
                isUIVisible = false // Only hides action buttons, not tabs
                onUIVisibilityChange(false)
            }
        } else {
            // Always show all UI for other tabs
            isUIVisible = true
            onUIVisibilityChange(true)
        }
    }

    // Video player management - Stop other videos when tab changes
    LaunchedEffect(selectedVideoTab) {
        currentPlayingVideoId = null
        isVideoPlaying = false
    }

    // Function to handle user interaction and reset timer
    val onUserInteraction = {
        lastInteractionTime = System.currentTimeMillis()
        isUIVisible = true
        onUIVisibilityChange(true)
    }

    // Category configuration
    val categories = remember {
        listOf(
            CategoryConfig(VVVideoCategory.ALL, "All") { Icon(Icons.Default.Home, null) },
            CategoryConfig(VVVideoCategory.TRENDING, "Trending") { Icon(Icons.Default.TrendingUp, null) },
            CategoryConfig(VVVideoCategory.FOOD, "Food") { Icon(Icons.Default.Restaurant, null) },
            CategoryConfig(VVVideoCategory.MUSIC, "Music") { Icon(Icons.Default.MusicNote, null) },
            CategoryConfig(VVVideoCategory.ENTERTAINMENT, "Fun") { Icon(Icons.Default.Star, null) },
            CategoryConfig(VVVideoCategory.EDUCATION, "Learn") { Icon(Icons.Default.School, null) },
            CategoryConfig(VVVideoCategory.TRAVEL, "Travel") { Icon(Icons.Default.Flight, null) }
        )
    }

    // Filter content by category
    val filteredShorts = if (selectedCategory == VVVideoCategory.ALL) shortVideos else shortVideos.filter { it.category == selectedCategory }
    val filteredLongs = if (selectedCategory == VVVideoCategory.ALL) longVideos else longVideos.filter { it.category == selectedCategory }
    val filteredChannels = if (selectedCategory == VVVideoCategory.ALL) channels else channels.filter { it.category == selectedCategory }

    Column(modifier = modifier.fillMaxSize().background(TchatColors.background)) {
        // Top App Bar
        TopAppBar(
            title = { Text("Videos", fontWeight = FontWeight.Bold) },
            actions = {
                // Add Settings button to existing top bar
                IconButton(onClick = onMoreClick) {
                    Icon(
                        Icons.Default.Settings,
                        "Settings",
                        tint = TchatColors.onSurface
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        // Optimized Tab Navigation - Always visible for better UX
        Surface(
            modifier = Modifier.fillMaxWidth(),
            shadowElevation = 4.dp
        ) {
            TabRow(
                selectedTabIndex = selectedVideoTab.ordinal,
                containerColor = TchatColors.surface,
                contentColor = TchatColors.primary
            ) {
                val videoTabs = listOf("Shorts", "Videos", "History", "Subscriptions")
                videoTabs.forEachIndexed { index, title ->
                    Tab(
                        selected = selectedVideoTab.ordinal == index,
                        onClick = {
                            // Stop current video when switching tabs
                            currentPlayingVideoId = null
                            isVideoPlaying = false

                            selectedVideoTab = VideoTabType.entries[index]
                            if (selectedVideoTab == VideoTabType.SHORTS) {
                                currentShortIndex = 0
                            }
                            onUserInteraction() // Reset timer when user interacts with tabs
                        },
                        text = {
                            Text(
                                text = title,
                                style = MaterialTheme.typography.titleSmall,
                                fontWeight = if (selectedVideoTab.ordinal == index) FontWeight.SemiBold else FontWeight.Normal
                            )
                        },
                        icon = {
                            Icon(
                                imageVector = when (index) {
                                    0 -> Icons.Default.MovieFilter // Shorts
                                    1 -> Icons.Default.PlayArrow   // Videos
                                    2 -> Icons.Default.History // History
                                    3 -> Icons.Default.VideoLibrary // Subscriptions
                                    else -> Icons.Default.PlayArrow
                                },
                                contentDescription = title,
                                modifier = Modifier.size(20.dp)
                            )
                        },
                        modifier = Modifier
                            .padding(vertical = 8.dp, horizontal = 4.dp)
                            .height(48.dp)
                    )
                }
            }
        }

        // Category Filter Bar (for Shorts and Videos tabs) with auto-hide functionality
        if (selectedVideoTab == VideoTabType.SHORTS || selectedVideoTab == VideoTabType.VIDEOS) {
            AnimatedVisibility(
                visible = isUIVisible || selectedVideoTab != VideoTabType.SHORTS,
                enter = slideInVertically(
                    initialOffsetY = { -it },
                    animationSpec = tween(300)
                ),
                exit = slideOutVertically(
                    targetOffsetY = { -it },
                    animationSpec = tween(300)
                )
            ) {
                LazyRow(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 8.dp)
                        .zIndex(9f), // Ensure category filter stays above video content
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(categories) { category ->
                        FilterChip(
                            onClick = {
                                selectedCategory = category.id
                                onUserInteraction() // Reset timer when user interacts with category filter
                            },
                            label = { Text(category.name) },
                            selected = selectedCategory == category.id,
                            leadingIcon = { category.icon() }
                        )
                    }
                }
            }
        }

        // Video Content based on selected tab
        when (selectedVideoTab) {
            VideoTabType.SHORTS -> OptimizedShortsPlayer(
                videos = filteredShorts,
                currentIndex = currentShortIndex,
                onIndexChange = { newIndex ->
                    currentShortIndex = newIndex
                    // Update currently playing video
                    currentPlayingVideoId = filteredShorts.getOrNull(newIndex)?.id
                },
                currentPlayingVideoId = currentPlayingVideoId,
                onVideoPlayStateChange = { videoId, isPlaying ->
                    currentPlayingVideoId = if (isPlaying) videoId else null
                    isVideoPlaying = isPlaying
                },
                onLike = { videoId ->
                    likedVideos = if (likedVideos.contains(videoId)) {
                        likedVideos - videoId
                    } else {
                        likedVideos + videoId
                    }
                },
                onShare = { video ->
                    selectedVideoForShare = video
                    showShareModal = true
                },
                onSubscribe = { channelId ->
                    subscribedChannels = if (subscribedChannels.contains(channelId)) {
                        subscribedChannels - channelId
                    } else {
                        subscribedChannels + channelId
                    }
                    // Update in repository
                    kotlin.runCatching {
                        CoroutineScope(Dispatchers.Default).launch {
                            videoRepository.subscribeToChannel(channelId)
                        }
                    }
                },
                onVideoTap = onUserInteraction, // Add tap interaction to show UI
                isUIVisible = isUIVisible, // Pass the UI visibility state
                modifier = Modifier.fillMaxHeight()
            )

            VideoTabType.VIDEOS -> OptimizedYouTubeStyleVideoFeed(
                videos = filteredLongs,
                likedVideos = likedVideos,
                subscribedChannels = subscribedChannels,
                currentPlayingVideoId = currentPlayingVideoId,
                onVideoClick = { videoId ->
                    // Stop other videos and set this one as playing
                    currentPlayingVideoId = videoId
                    onVideoClick(videoId)
                },
                onVideoPlayStateChange = { videoId, isPlaying ->
                    currentPlayingVideoId = if (isPlaying) videoId else null
                    isVideoPlaying = isPlaying
                },
                onLike = { videoId ->
                    likedVideos = if (likedVideos.contains(videoId)) {
                        likedVideos - videoId
                    } else {
                        likedVideos + videoId
                    }
                },
                onSubscribe = { channelId ->
                    subscribedChannels = if (subscribedChannels.contains(channelId)) {
                        subscribedChannels - channelId
                    } else {
                        subscribedChannels + channelId
                    }
                    // Update in repository
                    kotlin.runCatching {
                        CoroutineScope(Dispatchers.Default).launch {
                            videoRepository.subscribeToChannel(channelId)
                        }
                    }
                },
                onShare = { video ->
                    selectedVideoForShare = video
                    showShareModal = true
                },
                modifier = Modifier.fillMaxHeight()
            )

            VideoTabType.HISTORY -> HistoryTab(
                watchedVideos = shortVideos + longVideos, // Mock: show all videos as watched
                onVideoClick = onVideoClick,
                onRemoveFromHistory = { videoId ->
                    // TODO: Implement remove from history
                },
                modifier = Modifier.fillMaxHeight()
            )

            VideoTabType.SUBSCRIPTIONS -> SubscriptionsTab(
                channels = channels.filter { subscribedChannels.contains(it.id) },
                videos = longVideos.filter { video ->
                    subscribedChannels.contains(video.channel.id)
                },
                onVideoClick = onVideoClick,
                onSubscribe = { channelId ->
                    subscribedChannels = if (subscribedChannels.contains(channelId)) {
                        subscribedChannels - channelId
                    } else {
                        subscribedChannels + channelId
                    }
                    // Update in repository
                    kotlin.runCatching {
                        CoroutineScope(Dispatchers.Default).launch {
                            videoRepository.subscribeToChannel(channelId)
                        }
                    }
                },
                modifier = Modifier.fillMaxHeight()
            )
        }

        // Share Modal
        if (showShareModal && selectedVideoForShare != null) {
            TchatShareModal(
                isVisible = showShareModal,
                onDismiss = {
                    showShareModal = false
                    selectedVideoForShare = null
                },
                content = ShareContent(
                    type = ShareContentType.GENERAL,
                    title = selectedVideoForShare!!.title,
                    description = selectedVideoForShare!!.description,
                    url = "https://tchat.app/videos/${selectedVideoForShare!!.id}",
                    imageUrl = selectedVideoForShare!!.thumbnail
                )
            )
        }
    }
}

// ============================================================================
// COMPREHENSIVE VIDEO COMPONENTS USING VideoTypes.kt
// Following YouTube/TikTok patterns from web VideoTab.tsx
// ============================================================================

// TikTok-Style Shorts Player with Vertical Scroll Navigation
@Composable
private fun TikTokStyleShortsPlayer(
    videos: List<VVVideoContent>,
    currentIndex: Int,
    onIndexChange: (Int) -> Unit,
    onLike: (String) -> Unit,
    onShare: (VVVideoContent) -> Unit,
    onSubscribe: (String) -> Unit,
    onVideoTap: () -> Unit = {},
    isUIVisible: Boolean = true,
    modifier: Modifier = Modifier
) {
    if (videos.isEmpty()) {
        Box(
            modifier = modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.VideoLibrary, null, modifier = Modifier.size(64.dp))
                Text("No short videos available", style = MaterialTheme.typography.titleMedium)
            }
        }
        return
    }

    // TikTok-style animation states for smooth video transitions
    val density = LocalDensity.current
    var dragOffset by remember { mutableStateOf(0f) }
    var isTransitioning by remember { mutableStateOf(false) }

    BoxWithConstraints(
        modifier = modifier
            .fillMaxSize()
            .background(Color.Black)
    ) {
        val maxHeightPx = with(density) { maxHeight.toPx() }

        Box(
            modifier = Modifier
                .fillMaxSize()
                .pointerInput(currentIndex, videos.size) {
                    detectDragGestures(
                        onDragStart = {
                            dragOffset = 0f
                            isTransitioning = false
                        },
                        onDragEnd = {
                            val threshold = maxHeightPx * 0.15f // 15% of screen height for sensitivity

                            when {
                                dragOffset < -threshold && currentIndex < videos.size - 1 -> {
                                    // Swipe up - next video
                                    isTransitioning = true
                                    onIndexChange(currentIndex + 1)
                                }
                                dragOffset > threshold && currentIndex > 0 -> {
                                    // Swipe down - previous video
                                    isTransitioning = true
                                    onIndexChange(currentIndex - 1)
                                }
                            }
                            dragOffset = 0f
                        }
                    ) { change, dragAmount ->
                        if (!isTransitioning) {
                            dragOffset += dragAmount.y
                            change.consume()
                        }
                    }
                }
        ) {
            // Render previous, current, and next videos for smooth transitions
            val videosToRender = listOf(
                currentIndex - 1 to videos.getOrNull(currentIndex - 1),
                currentIndex to videos.getOrNull(currentIndex),
                currentIndex + 1 to videos.getOrNull(currentIndex + 1)
            )

            videosToRender.forEach { (index, video) ->
                if (video != null) {
                    val offsetY = when (index) {
                        currentIndex - 1 -> -maxHeightPx + dragOffset // Previous video above
                        currentIndex -> dragOffset // Current video
                        currentIndex + 1 -> maxHeightPx + dragOffset // Next video below
                        else -> 0f
                    }

                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .offset { IntOffset(0, offsetY.toInt()) }
                    ) {
                        // Video Player with TikTok-style transitions
                        key(video.id, index) {
                            TchatVideo(
                                source = VideoSource.Url(video.videoUrl ?: ""),
                                modifier = Modifier.fillMaxSize(),
                                aspectRatio = TchatVideoAspectRatio.Portrait,
                                autoPlay = index == currentIndex, // Only auto-play current video
                                loop = true,
                                muted = true,
                                showControls = false,
                                poster = ImageSource.Url(video.thumbnail),
                                onStateChange = { state ->
                                    // Handle video state changes
                                }
                            )
                        }

                        // Only show tap overlay for current video
                        if (index == currentIndex) {
                            // Invisible tap overlay for UI show/hide (only in center area, avoiding action buttons)
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .padding(end = 100.dp) // Avoid right action buttons area
                                    .clickable(
                                        indication = null,
                                        interactionSource = remember { MutableInteractionSource() }
                                    ) {
                                        onVideoTap() // Show UI when video is tapped
                                    }
                            )
                        }
                    }
                }
            }
        }

        // Get current video for UI overlay positioned outside the video container
        val currentVideo = videos.getOrNull(currentIndex)
        if (currentVideo != null) {
            // Right Side Actions (TikTok Style) with auto-hide functionality
            AnimatedVisibility(
                visible = isUIVisible,
                enter = slideInHorizontally(
                    initialOffsetX = { it },
                    animationSpec = tween(300)
                ),
                exit = slideOutHorizontally(
                    targetOffsetX = { it },
                    animationSpec = tween(300)
                ),
                modifier = Modifier.align(Alignment.CenterEnd)
            ) {
                Column(
                    modifier = Modifier
                        .padding(end = 16.dp)
                        .zIndex(1f),
                    verticalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                // Like Button
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Surface(
                        modifier = Modifier.size(48.dp),
                        shape = CircleShape,
                        color = Color.Black.copy(alpha = 0.3f)
                    ) {
                        IconButton(onClick = { onLike(currentVideo.id) }) {
                            Icon(
                                Icons.Default.Favorite,
                                contentDescription = "Like",
                                tint = Color.White,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    }
                    Text(
                        text = formatViewCount(currentVideo.likes),
                        color = Color.White,
                        style = MaterialTheme.typography.labelSmall
                    )
                }

                // Comments Button
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Surface(
                        modifier = Modifier.size(48.dp),
                        shape = CircleShape,
                        color = Color.Black.copy(alpha = 0.3f)
                    ) {
                        IconButton(onClick = { /* Navigate to comments */ }) {
                            Icon(
                                Icons.Default.ChatBubbleOutline,
                                contentDescription = "Comments",
                                tint = Color.White,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    }
                    Text(
                        text = formatViewCount(currentVideo.comments),
                        color = Color.White,
                        style = MaterialTheme.typography.labelSmall
                    )
                }

                // Share Button
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Surface(
                        modifier = Modifier.size(48.dp),
                        shape = CircleShape,
                        color = Color.Black.copy(alpha = 0.3f)
                    ) {
                        IconButton(onClick = { onShare(currentVideo) }) {
                            Icon(
                                Icons.Default.Share,
                                contentDescription = "Share",
                                tint = Color.White,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    }
                    Text(
                        text = "Share",
                        color = Color.White,
                        style = MaterialTheme.typography.labelSmall
                    )
                }
                }
            }

            // Bottom Info Overlay with auto-hide functionality
            AnimatedVisibility(
                visible = isUIVisible,
                enter = slideInVertically(
                    initialOffsetY = { it },
                    animationSpec = tween(300)
                ),
                exit = slideOutVertically(
                    targetOffsetY = { it },
                    animationSpec = tween(300)
                ),
                modifier = Modifier.align(Alignment.BottomStart)
            ) {
                Box(
                    modifier = Modifier
                        .wrapContentSize()
                        .background(
                            Brush.verticalGradient(
                                colors = listOf(Color.Transparent, Color.Black.copy(alpha = 0.8f)),
                                startY = 0f,
                                endY = 200f
                            ),
                            shape = RoundedCornerShape(topEnd = 12.dp)
                        )
                        .padding(16.dp)
                ) {
                Column {
                    // Channel Info
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier.padding(bottom = 8.dp)
                    ) {
                        AsyncImage(
                            model = currentVideo.channel.avatar,
                            contentDescription = "${currentVideo.channel.name} avatar",
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .clickable {
                                    // TODO: Navigate to channel detail screen
                                    println("Clicked on channel: ${currentVideo.channel.name}")
                                },
                            contentScale = ContentScale.Crop
                        )

                        Spacer(modifier = Modifier.width(12.dp))

                        Column(modifier = Modifier.weight(1f)) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Text(
                                    text = currentVideo.channel.name,
                                    color = Color.White,
                                    style = MaterialTheme.typography.titleSmall,
                                    fontWeight = FontWeight.Bold,
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis
                                )
                                if (currentVideo.channel.verified) {
                                    Icon(
                                        Icons.Default.Verified,
                                        contentDescription = "Verified",
                                        tint = Color.Blue,
                                        modifier = Modifier.size(16.dp).padding(start = 4.dp)
                                    )
                                }
                            }
                            Text(
                                text = "${formatSubscriberCount(currentVideo.channel.subscribers)} subscribers",
                                color = Color.White.copy(alpha = 0.7f),
                                style = MaterialTheme.typography.bodySmall,
                                maxLines = 1,
                                overflow = TextOverflow.Ellipsis
                            )
                        }

                        Spacer(modifier = Modifier.width(8.dp))

                        // Subscribe Button
                        Button(
                            onClick = { onSubscribe(currentVideo.channel.id) },
                            colors = ButtonDefaults.buttonColors(
                                containerColor = Color.Red,
                                contentColor = Color.White
                            ),
                            modifier = Modifier.height(32.dp)
                        ) {
                            Text("Follow", style = MaterialTheme.typography.labelMedium)
                        }
                    }

                    // Video Title and Description
                    Text(
                        text = currentVideo.title,
                        color = Color.White,
                        style = MaterialTheme.typography.bodyMedium,
                        maxLines = 2,
                        overflow = TextOverflow.Ellipsis
                    )

                    if (currentVideo.description.isNotEmpty()) {
                        Text(
                            text = currentVideo.description,
                            color = Color.White.copy(alpha = 0.8f),
                            style = MaterialTheme.typography.bodySmall,
                            maxLines = 2,
                            overflow = TextOverflow.Ellipsis,
                            modifier = Modifier.padding(top = 4.dp)
                        )
                    }

                    // Tags
                    if (currentVideo.tags.isNotEmpty()) {
                        LazyRow(
                            horizontalArrangement = Arrangement.spacedBy(8.dp),
                            modifier = Modifier.padding(top = 8.dp)
                        ) {
                            items(currentVideo.tags.take(3)) { tag ->
                                Surface(
                                    color = Color.White.copy(alpha = 0.2f),
                                    shape = RoundedCornerShape(12.dp)
                                ) {
                                    Text(
                                        text = "#$tag",
                                        color = Color.White,
                                        style = MaterialTheme.typography.labelSmall,
                                        modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp)
                                    )
                                }
                            }
                        }
                    }
                }
                }
            }

            // Video Progress Indicator with Navigation Hints
            Column(
                modifier = Modifier
                    .align(Alignment.CenterStart)
                    .padding(start = 8.dp),
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                // "Scroll up for next" hint when not at last video
                if (currentIndex < videos.size - 1) {
                    Icon(
                        Icons.Default.ExpandLess,
                        contentDescription = "Scroll up for next video",
                        tint = Color.White.copy(alpha = 0.5f),
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }

                repeat(videos.size) { index ->
                    Surface(
                        modifier = Modifier
                            .width(3.dp)
                            .height(if (index == currentIndex) 24.dp else 12.dp),
                        color = if (index == currentIndex) Color.White else Color.White.copy(alpha = 0.3f),
                        shape = RoundedCornerShape(2.dp)
                    ) {}
                }

                // "Scroll down for previous" hint when not at first video
                if (currentIndex > 0) {
                    Spacer(modifier = Modifier.height(8.dp))
                    Icon(
                        Icons.Default.ExpandMore,
                        contentDescription = "Scroll down for previous video",
                        tint = Color.White.copy(alpha = 0.5f),
                        modifier = Modifier.size(16.dp)
                    )
                }
            }

            // End of list indicators
            if (currentIndex == 0) {
                // At first video - show "Start of videos" indicator
                Text(
                    text = "Start of videos",
                    color = Color.White.copy(alpha = 0.6f),
                    style = MaterialTheme.typography.labelSmall,
                    modifier = Modifier
                        .align(Alignment.TopCenter)
                        .padding(top = 16.dp)
                        .background(
                            Color.Black.copy(alpha = 0.4f),
                            RoundedCornerShape(12.dp)
                        )
                        .padding(horizontal = 12.dp, vertical = 6.dp)
                )
            }

            if (currentIndex == videos.size - 1) {
                // At last video - show "End of videos" indicator
                Text(
                    text = "End of videos",
                    color = Color.White.copy(alpha = 0.6f),
                    style = MaterialTheme.typography.labelSmall,
                    modifier = Modifier
                        .align(Alignment.BottomCenter)
                        .padding(bottom = 120.dp) // Above the video info
                        .background(
                            Color.Black.copy(alpha = 0.4f),
                            RoundedCornerShape(12.dp)
                        )
                        .padding(horizontal = 12.dp, vertical = 6.dp)
                )
            }
        }
    }
}

// ============================================================================
// OPTIMIZED SHORTS PLAYER - Single video playback with better performance
// ============================================================================

@Composable
private fun OptimizedShortsPlayer(
    videos: List<VVVideoContent>,
    currentIndex: Int,
    onIndexChange: (Int) -> Unit,
    currentPlayingVideoId: String?,
    onVideoPlayStateChange: (String, Boolean) -> Unit,
    onLike: (String) -> Unit,
    onShare: (VVVideoContent) -> Unit,
    onSubscribe: (String) -> Unit,
    onVideoTap: () -> Unit = {},
    isUIVisible: Boolean = true,
    modifier: Modifier = Modifier
) {
    if (videos.isEmpty()) {
        Box(
            modifier = modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.VideoLibrary, null, modifier = Modifier.size(64.dp))
                Text("No short videos available", style = MaterialTheme.typography.titleMedium)
            }
        }
        return
    }

    // Optimized animation states - no multi-video rendering
    val density = LocalDensity.current
    var dragOffset by remember { mutableStateOf(0f) }
    var isTransitioning by remember { mutableStateOf(false) }

    BoxWithConstraints(
        modifier = modifier
            .fillMaxSize()
            .background(Color.Black)
    ) {
        val maxHeightPx = with(density) { maxHeight.toPx() }

        Box(
            modifier = Modifier
                .fillMaxSize()
                .pointerInput(currentIndex, videos.size) {
                    detectDragGestures(
                        onDragStart = {
                            dragOffset = 0f
                            isTransitioning = false
                        },
                        onDragEnd = {
                            val threshold = maxHeightPx * 0.15f // 15% of screen height for sensitivity

                            when {
                                dragOffset < -threshold && currentIndex < videos.size - 1 -> {
                                    // Swipe up - next video (stop current video first)
                                    onVideoPlayStateChange(videos[currentIndex].id, false)
                                    isTransitioning = true
                                    onIndexChange(currentIndex + 1)
                                }
                                dragOffset > threshold && currentIndex > 0 -> {
                                    // Swipe down - previous video (stop current video first)
                                    onVideoPlayStateChange(videos[currentIndex].id, false)
                                    isTransitioning = true
                                    onIndexChange(currentIndex - 1)
                                }
                            }
                            dragOffset = 0f
                        }
                    ) { change, dragAmount ->
                        if (!isTransitioning) {
                            dragOffset += dragAmount.y
                            change.consume()
                        }
                    }
                }
        ) {
            // Render ONLY the current video for better performance
            val currentVideo = videos.getOrNull(currentIndex)
            if (currentVideo != null) {
                // Single Video Player with optimized state management
                key(currentVideo.id) {
                    TchatVideo(
                        source = VideoSource.Url(currentVideo.videoUrl ?: ""),
                        modifier = Modifier.fillMaxSize(),
                        aspectRatio = TchatVideoAspectRatio.Portrait,
                        autoPlay = currentPlayingVideoId == currentVideo.id,
                        loop = true,
                        muted = true,
                        showControls = false,
                        poster = ImageSource.Url(currentVideo.thumbnail),
                        onStateChange = { state ->
                            // Update global video state
                            when (state) {
                                com.tchat.mobile.components.TchatVideoState.Playing -> {
                                    onVideoPlayStateChange(currentVideo.id, true)
                                }
                                com.tchat.mobile.components.TchatVideoState.Paused,
                                com.tchat.mobile.components.TchatVideoState.Ended -> {
                                    onVideoPlayStateChange(currentVideo.id, false)
                                }
                                else -> { /* Handle other states if needed */ }
                            }
                        }
                    )
                }

                // Tap overlay for UI show/hide (avoiding right action buttons area)
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(end = 100.dp) // Avoid right action buttons area
                        .clickable(
                            indication = null,
                            interactionSource = remember { MutableInteractionSource() }
                        ) {
                            onVideoTap() // Show UI when video is tapped
                        }
                )
            }
        }

        // Get current video for UI overlay
        val currentVideo = videos.getOrNull(currentIndex)
        if (currentVideo != null) {
            // Right Side Actions (TikTok Style) with performance optimization
            AnimatedVisibility(
                visible = isUIVisible,
                enter = slideInHorizontally(
                    initialOffsetX = { it },
                    animationSpec = tween(300)
                ),
                exit = slideOutHorizontally(
                    targetOffsetX = { it },
                    animationSpec = tween(300)
                ),
                modifier = Modifier.align(Alignment.CenterEnd)
            ) {
                Column(
                    modifier = Modifier
                        .padding(end = 16.dp)
                        .zIndex(1f),
                    verticalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                    // Like Button
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Surface(
                            modifier = Modifier.size(48.dp),
                            shape = CircleShape,
                            color = Color.Black.copy(alpha = 0.3f)
                        ) {
                            IconButton(onClick = { onLike(currentVideo.id) }) {
                                Icon(
                                    Icons.Default.Favorite,
                                    contentDescription = "Like",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }
                        }
                        Text(
                            text = formatViewCount(currentVideo.likes),
                            color = Color.White,
                            style = MaterialTheme.typography.labelSmall
                        )
                    }

                    // Subscribe Button
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Surface(
                            modifier = Modifier.size(48.dp),
                            shape = CircleShape,
                            color = Color.Black.copy(alpha = 0.3f)
                        ) {
                            IconButton(onClick = { onSubscribe(currentVideo.channel.id) }) {
                                Icon(
                                    Icons.Default.PersonAdd,
                                    contentDescription = "Subscribe",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }
                        }
                        Text(
                            text = "Follow",
                            color = Color.White,
                            style = MaterialTheme.typography.labelSmall
                        )
                    }

                    // Share Button
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Surface(
                            modifier = Modifier.size(48.dp),
                            shape = CircleShape,
                            color = Color.Black.copy(alpha = 0.3f)
                        ) {
                            IconButton(onClick = { onShare(currentVideo) }) {
                                Icon(
                                    Icons.Default.Share,
                                    contentDescription = "Share",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }
                        }
                        Text(
                            text = "Share",
                            color = Color.White,
                            style = MaterialTheme.typography.labelSmall
                        )
                    }
                }
            }

            // Optimized Bottom Info Overlay
            AnimatedVisibility(
                visible = isUIVisible,
                enter = slideInVertically(
                    initialOffsetY = { it },
                    animationSpec = tween(300)
                ),
                exit = slideOutVertically(
                    targetOffsetY = { it },
                    animationSpec = tween(300)
                ),
                modifier = Modifier.align(Alignment.BottomStart)
            ) {
                Box(
                    modifier = Modifier
                        .wrapContentSize()
                        .background(
                            Brush.verticalGradient(
                                colors = listOf(Color.Transparent, Color.Black.copy(alpha = 0.8f)),
                                startY = 0f,
                                endY = 200f
                            ),
                            shape = RoundedCornerShape(topEnd = 12.dp)
                        )
                        .padding(16.dp)
                ) {
                    Column {
                        // Channel Info
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            modifier = Modifier.padding(bottom = 8.dp)
                        ) {
                            AsyncImage(
                                model = currentVideo.channel.avatar,
                                contentDescription = "${currentVideo.channel.name} avatar",
                                modifier = Modifier
                                    .size(40.dp)
                                    .clip(CircleShape),
                                contentScale = ContentScale.Crop
                            )

                            Spacer(modifier = Modifier.width(12.dp))

                            Column(modifier = Modifier.weight(1f)) {
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    Text(
                                        text = currentVideo.channel.name,
                                        color = Color.White,
                                        style = MaterialTheme.typography.titleSmall,
                                        fontWeight = FontWeight.Bold,
                                        maxLines = 1,
                                        overflow = TextOverflow.Ellipsis
                                    )
                                    if (currentVideo.channel.verified) {
                                        Icon(
                                            Icons.Default.Verified,
                                            contentDescription = "Verified",
                                            tint = Color.Blue,
                                            modifier = Modifier.size(16.dp).padding(start = 4.dp)
                                        )
                                    }
                                }
                                Text(
                                    text = "${formatSubscriberCount(currentVideo.channel.subscribers)} subscribers",
                                    color = Color.White.copy(alpha = 0.7f),
                                    style = MaterialTheme.typography.bodySmall,
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis
                                )
                            }
                        }

                        // Video Title
                        Text(
                            text = currentVideo.title,
                            color = Color.White,
                            style = MaterialTheme.typography.bodyMedium,
                            maxLines = 2,
                            overflow = TextOverflow.Ellipsis
                        )

                        // Tags (limit for performance)
                        if (currentVideo.tags.isNotEmpty()) {
                            LazyRow(
                                horizontalArrangement = Arrangement.spacedBy(8.dp),
                                modifier = Modifier.padding(top = 8.dp)
                            ) {
                                items(currentVideo.tags.take(3)) { tag ->
                                    Surface(
                                        color = Color.White.copy(alpha = 0.2f),
                                        shape = RoundedCornerShape(12.dp)
                                    ) {
                                        Text(
                                            text = "#$tag",
                                            color = Color.White,
                                            style = MaterialTheme.typography.labelSmall,
                                            modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp)
                                        )
                                    }
                                }
                            }
                        }
                    }
                }
            }

            // Simplified Progress Indicator
            Column(
                modifier = Modifier
                    .align(Alignment.CenterStart)
                    .padding(start = 8.dp),
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                repeat(minOf(videos.size, 10)) { index -> // Limit indicators for performance
                    Surface(
                        modifier = Modifier
                            .width(3.dp)
                            .height(if (index == currentIndex) 24.dp else 12.dp),
                        color = if (index == currentIndex) Color.White else Color.White.copy(alpha = 0.3f),
                        shape = RoundedCornerShape(2.dp)
                    ) {}
                }
            }
        }
    }
}

// YouTube-Style Video Feed
@Composable
private fun YouTubeStyleVideoFeed(
    videos: List<VVVideoContent>,
    likedVideos: Set<String>,
    subscribedChannels: Set<String>,
    onVideoClick: (String) -> Unit,
    onLike: (String) -> Unit,
    onSubscribe: (String) -> Unit,
    onShare: (VVVideoContent) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        items(videos) { video ->
            YouTubeVideoCard(
                video = video,
                isLiked = likedVideos.contains(video.id),
                isSubscribed = subscribedChannels.contains(video.channel.id),
                onVideoClick = { onVideoClick(video.id) },
                onLike = { onLike(video.id) },
                onSubscribe = { onSubscribe(video.channel.id) },
                onShare = { onShare(video) }
            )
        }
    }
}

// ============================================================================
// OPTIMIZED YOUTUBE-STYLE VIDEO FEED - Single video playback
// ============================================================================

@Composable
private fun OptimizedYouTubeStyleVideoFeed(
    videos: List<VVVideoContent>,
    likedVideos: Set<String>,
    subscribedChannels: Set<String>,
    currentPlayingVideoId: String?,
    onVideoClick: (String) -> Unit,
    onVideoPlayStateChange: (String, Boolean) -> Unit,
    onLike: (String) -> Unit,
    onSubscribe: (String) -> Unit,
    onShare: (VVVideoContent) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        items(videos) { video ->
            OptimizedYouTubeVideoCard(
                video = video,
                isLiked = likedVideos.contains(video.id),
                isSubscribed = subscribedChannels.contains(video.channel.id),
                isPlaying = currentPlayingVideoId == video.id,
                onVideoClick = {
                    // Stop other videos and start this one
                    onVideoPlayStateChange(video.id, true)
                    onVideoClick(video.id)
                },
                onVideoPlayStateChange = onVideoPlayStateChange,
                onLike = { onLike(video.id) },
                onSubscribe = { onSubscribe(video.channel.id) },
                onShare = { onShare(video) }
            )
        }
    }
}

@Composable
private fun OptimizedYouTubeVideoCard(
    video: VVVideoContent,
    isLiked: Boolean,
    isSubscribed: Boolean,
    isPlaying: Boolean,
    onVideoClick: () -> Unit,
    onVideoPlayStateChange: (String, Boolean) -> Unit,
    onLike: () -> Unit,
    onSubscribe: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column(modifier = Modifier.fillMaxWidth()) {
            // Video Player/Thumbnail Section
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .aspectRatio(16f / 9f) // YouTube standard aspect ratio
                    .background(Color.Black)
            ) {
                if (isPlaying) {
                    // Show actual video player when playing
                    key(video.id) {
                        TchatVideo(
                            source = VideoSource.Url(video.videoUrl ?: ""),
                            modifier = Modifier.fillMaxSize(),
                            aspectRatio = TchatVideoAspectRatio.Landscape,
                            autoPlay = true,
                            loop = false, // YouTube style - don't loop by default
                            muted = false, // YouTube style - allow sound
                            showControls = true, // YouTube style - show controls
                            poster = ImageSource.Url(video.thumbnail),
                            onStateChange = { state ->
                                when (state) {
                                    com.tchat.mobile.components.TchatVideoState.Playing -> {
                                        onVideoPlayStateChange(video.id, true)
                                    }
                                    com.tchat.mobile.components.TchatVideoState.Paused,
                                    com.tchat.mobile.components.TchatVideoState.Ended -> {
                                        onVideoPlayStateChange(video.id, false)
                                    }
                                    else -> { /* Handle other states if needed */ }
                                }
                            }
                        )
                    }
                } else {
                    // Show thumbnail with play button when not playing
                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .clickable { onVideoClick() }
                    ) {
                        AsyncImage(
                            model = video.thumbnail,
                            contentDescription = video.title,
                            modifier = Modifier.fillMaxSize(),
                            contentScale = ContentScale.Crop
                        )

                        // Play button overlay
                        Surface(
                            modifier = Modifier
                                .size(64.dp)
                                .align(Alignment.Center),
                            shape = CircleShape,
                            color = Color.Black.copy(alpha = 0.7f)
                        ) {
                            Icon(
                                Icons.Default.PlayArrow,
                                contentDescription = "Play video",
                                tint = Color.White,
                                modifier = Modifier
                                    .fillMaxSize()
                                    .padding(16.dp)
                            )
                        }

                        // Duration badge
                        Surface(
                            modifier = Modifier
                                .align(Alignment.BottomEnd)
                                .padding(8.dp),
                            shape = RoundedCornerShape(4.dp),
                            color = Color.Black.copy(alpha = 0.8f)
                        ) {
                            Text(
                                text = video.duration,
                                color = Color.White,
                                style = MaterialTheme.typography.labelSmall,
                                modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                            )
                        }
                    }
                }
            }

            // Video Info Section
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(12.dp)
            ) {
                // Title
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(4.dp))

                // Views and upload time
                Text(
                    text = "${formatViewCount(video.views)} views • ${video.uploadTimeFormatted}",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                Spacer(modifier = Modifier.height(8.dp))

                // Channel info and subscribe button
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    AsyncImage(
                        model = video.channel.avatar,
                        contentDescription = "${video.channel.name} avatar",
                        modifier = Modifier
                            .size(36.dp)
                            .clip(CircleShape),
                        contentScale = ContentScale.Crop
                    )

                    Spacer(modifier = Modifier.width(12.dp))

                    Column(modifier = Modifier.weight(1f)) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(
                                text = video.channel.name,
                                style = MaterialTheme.typography.bodyMedium,
                                fontWeight = FontWeight.Medium,
                                color = TchatColors.onSurface
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
                            text = "${formatSubscriberCount(video.channel.subscribers)} subscribers",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    // Subscribe Button
                    TchatButton(
                        text = if (isSubscribed) "Subscribed" else "Subscribe",
                        onClick = onSubscribe,
                        variant = if (isSubscribed) TchatButtonVariant.Outline else TchatButtonVariant.Primary,
                        modifier = Modifier.height(32.dp)
                    )
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Action buttons
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                        // Like button
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            modifier = Modifier.clickable { onLike() }
                        ) {
                            Icon(
                                if (isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                                contentDescription = "Like",
                                tint = if (isLiked) Color.Red else TchatColors.onSurfaceVariant,
                                modifier = Modifier.size(20.dp)
                            )
                            Spacer(modifier = Modifier.width(4.dp))
                            Text(
                                text = formatViewCount(video.likes),
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }

                        // Share button
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            modifier = Modifier.clickable { onShare() }
                        ) {
                            Icon(
                                Icons.Default.Share,
                                contentDescription = "Share",
                                tint = TchatColors.onSurfaceVariant,
                                modifier = Modifier.size(20.dp)
                            )
                            Spacer(modifier = Modifier.width(4.dp))
                            Text(
                                text = "Share",
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun YouTubeVideoCard(
    video: VVVideoContent,
    isLiked: Boolean,
    isSubscribed: Boolean,
    onVideoClick: () -> Unit,
    onLike: () -> Unit,
    onSubscribe: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onVideoClick() }
            .padding(horizontal = 16.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            // Video Player
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .aspectRatio(16f / 9f)
            ) {
                TchatVideo(
                    source = VideoSource.Url(video.videoUrl ?: ""),
                    modifier = Modifier.fillMaxSize(),
                    aspectRatio = TchatVideoAspectRatio.Landscape,
                    autoPlay = false, // Don't auto-play in feed
                    loop = false,
                    muted = false,
                    showControls = true, // Show controls for YouTube-style
                    poster = ImageSource.Url(video.thumbnail),
                    onStateChange = { state ->
                        // Handle video state changes
                    }
                )

                // Duration Badge
                Surface(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(8.dp),
                    color = Color.Black.copy(alpha = 0.8f),
                    shape = RoundedCornerShape(4.dp)
                ) {
                    Text(
                        text = video.duration,
                        color = Color.White,
                        style = MaterialTheme.typography.labelSmall,
                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                    )
                }
            }

            // Video Info
            Column(modifier = Modifier.padding(12.dp)) {
                // Title
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis,
                    modifier = Modifier.padding(bottom = 8.dp)
                )

                // Channel Info Row
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.padding(bottom = 8.dp)
                ) {
                    AsyncImage(
                        model = video.channel.avatar,
                        contentDescription = "${video.channel.name} avatar",
                        modifier = Modifier
                            .size(32.dp)
                            .clip(CircleShape)
                            .clickable {
                                // TODO: Navigate to channel detail screen
                                println("Clicked on channel in video card: ${video.channel.name}")
                            },
                        contentScale = ContentScale.Crop
                    )

                    Spacer(modifier = Modifier.width(8.dp))

                    Column(modifier = Modifier.fillMaxHeight()) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(
                                text = video.channel.name,
                                style = MaterialTheme.typography.bodyMedium,
                                fontWeight = FontWeight.Medium
                            )
                            if (video.channel.verified) {
                                Icon(
                                    Icons.Default.Verified,
                                    contentDescription = "Verified",
                                    tint = Color.Blue,
                                    modifier = Modifier.size(14.dp).padding(start = 4.dp)
                                )
                            }
                        }
                        Text(
                            text = "${formatViewCount(video.views)} views • ${video.uploadTimeFormatted}",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    // Subscribe Button
                    Button(
                        onClick = onSubscribe,
                        colors = ButtonDefaults.buttonColors(
                            containerColor = if (isSubscribed) TchatColors.surfaceVariant else TchatColors.primary,
                            contentColor = if (isSubscribed) TchatColors.onSurfaceVariant else TchatColors.onPrimary
                        ),
                        modifier = Modifier.height(32.dp)
                    ) {
                        Text(
                            text = if (isSubscribed) "Subscribed" else "Subscribe",
                            style = MaterialTheme.typography.labelMedium
                        )
                    }
                }

                // Engagement Row
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(16.dp)
                ) {
                    // Like Button
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier.clickable { onLike() }
                    ) {
                        Icon(
                            if (isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                            contentDescription = "Like",
                            tint = if (isLiked) Color.Red else TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(20.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatViewCount(video.likes),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    // Comments
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            Icons.Default.ChatBubbleOutline,
                            contentDescription = "Comments",
                            tint = TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(20.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatViewCount(video.comments),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    // Share Button
                    IconButton(
                        onClick = onShare,
                        modifier = Modifier.size(24.dp)
                    ) {
                        Icon(
                            Icons.Default.Share,
                            contentDescription = "Share",
                            tint = TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(20.dp)
                        )
                    }
                }
            }
        }
    }
}

// Channels Explore Tab
@Composable
private fun ChannelsExplore(
    channels: List<VVChannelInfo>,
    subscribedChannels: Set<String>,
    onSubscribe: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    if (channels.isEmpty()) {
        Box(
            modifier = modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.VideoLibrary, null, modifier = Modifier.size(64.dp))
                Text("No channels available", style = MaterialTheme.typography.titleMedium)
            }
        }
        return
    }

    LazyVerticalGrid(
        columns = GridCells.Fixed(2),
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(16.dp),
        horizontalArrangement = Arrangement.spacedBy(12.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        items(channels) { channel ->
            ChannelCard(
                channel = channel,
                isSubscribed = subscribedChannels.contains(channel.id),
                onSubscribe = { onSubscribe(channel.id) }
            )
        }
    }
}

@Composable
private fun ChannelCard(
    channel: VVChannelInfo,
    isSubscribed: Boolean,
    onSubscribe: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // Channel Avatar - Clickable
            AsyncImage(
                model = channel.avatar,
                contentDescription = "${channel.name} avatar",
                modifier = Modifier
                    .size(64.dp)
                    .clip(CircleShape)
                    .clickable {
                        // TODO: Navigate to channel detail screen
                        println("Clicked on channel card: ${channel.name}")
                    },
                contentScale = ContentScale.Crop
            )

            Spacer(modifier = Modifier.height(12.dp))

            // Channel Name
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.Center
            ) {
                Text(
                    text = channel.name,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )
                if (channel.verified) {
                    Icon(
                        Icons.Default.Verified,
                        contentDescription = "Verified",
                        tint = Color.Blue,
                        modifier = Modifier.size(16.dp).padding(start = 4.dp)
                    )
                }
            }

            // Subscriber Count
            Text(
                text = "${formatSubscriberCount(channel.subscribers)} subscribers",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center
            )

            // Videos Count (using mock data)
            Text(
                text = "25 videos", // TODO: Add videoCount to VVChannelInfo model
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center
            )

            Spacer(modifier = Modifier.height(12.dp))

            // Subscribe Button
            Button(
                onClick = onSubscribe,
                colors = ButtonDefaults.buttonColors(
                    containerColor = if (isSubscribed) TchatColors.surfaceVariant else TchatColors.primary,
                    contentColor = if (isSubscribed) TchatColors.onSurfaceVariant else TchatColors.onPrimary
                ),
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = if (isSubscribed) "Subscribed" else "Subscribe",
                    style = MaterialTheme.typography.labelMedium
                )
            }
        }
    }
}

// Subscriptions Tab
@Composable
private fun SubscriptionsTab(
    channels: List<VVChannelInfo>,
    videos: List<VVVideoContent>,
    onVideoClick: (String) -> Unit,
    onSubscribe: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    if (channels.isEmpty()) {
        Box(
            modifier = modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                modifier = Modifier.padding(32.dp)
            ) {
                Icon(
                    Icons.Default.Subscriptions,
                    contentDescription = "No subscriptions",
                    modifier = Modifier.size(80.dp),
                    tint = TchatColors.onSurfaceVariant
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "No subscriptions yet",
                    style = MaterialTheme.typography.titleMedium,
                    color = TchatColors.onSurfaceVariant
                )
                Text(
                    text = "Subscribe to channels to see their latest videos here",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant,
                    textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                    modifier = Modifier.padding(top = 8.dp)
                )
            }
        }
        return
    }

    LazyColumn(
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        // Subscribed Channels Section
        item {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 16.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "Your Subscriptions",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 12.dp)
                    )

                    LazyRow(
                        horizontalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        items(channels) { channel ->
                            Column(
                                horizontalAlignment = Alignment.CenterHorizontally,
                                modifier = Modifier
                                    .width(80.dp)
                                    .clickable {
                                        // TODO: Navigate to channel detail screen
                                        println("Clicked on subscribed channel: ${channel.name}")
                                    }
                            ) {
                                AsyncImage(
                                    model = channel.avatar,
                                    contentDescription = "${channel.name} avatar",
                                    modifier = Modifier
                                        .size(56.dp)
                                        .clip(CircleShape),
                                    contentScale = ContentScale.Crop
                                )
                                Text(
                                    text = channel.name,
                                    style = MaterialTheme.typography.labelSmall,
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis,
                                    textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                                    modifier = Modifier.padding(top = 4.dp)
                                )
                            }
                        }
                    }
                }
            }
        }

        // Recent Videos from Subscriptions
        item {
            Text(
                text = "Latest Videos",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
            )
        }

        items(videos) { video ->
            YouTubeVideoCard(
                video = video,
                isLiked = false, // Could be enhanced with liked state
                isSubscribed = true, // Always true for subscription tab
                onVideoClick = { onVideoClick(video.id) },
                onLike = { /* Handle like */ },
                onSubscribe = { onSubscribe(video.channel.id) },
                onShare = { /* Handle share */ }
            )
        }
    }
}

@Composable
fun HistoryTab(
    watchedVideos: List<VVVideoContent>,
    onVideoClick: (String) -> Unit,
    onRemoveFromHistory: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        item {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 16.dp, vertical = 8.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "Watch History",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )

                TextButton(
                    onClick = { /* TODO: Clear all history */ }
                ) {
                    Text(
                        text = "Clear All",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.primary
                    )
                }
            }
        }

        items(watchedVideos) { video ->
            HistoryVideoCard(
                video = video,
                onVideoClick = { onVideoClick(video.id) },
                onRemoveFromHistory = { onRemoveFromHistory(video.id) }
            )
        }
    }
}

@Composable
fun HistoryVideoCard(
    video: VVVideoContent,
    onVideoClick: () -> Unit,
    onRemoveFromHistory: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp)
            .clickable { onVideoClick() },
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(12.dp),
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Video Thumbnail
            AsyncImage(
                model = video.thumbnail,
                contentDescription = "Video thumbnail",
                modifier = Modifier
                    .size(width = 120.dp, height = 68.dp)
                    .clip(RoundedCornerShape(8.dp)),
                contentScale = ContentScale.Crop
            )

            // Video Info
            Column(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight(),
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Text(
                    text = video.channel.name,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                Row(
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        text = formatViewCount(video.views) + " views",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )

                    Text(
                        text = "•",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )

                    Text(
                        text = video.uploadTimeFormatted,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }

            // Remove from history button
            IconButton(
                onClick = onRemoveFromHistory,
                modifier = Modifier.size(24.dp)
            ) {
                Icon(
                    Icons.Default.Close,
                    contentDescription = "Remove from history",
                    tint = TchatColors.onSurfaceVariant,
                    modifier = Modifier.size(20.dp)
                )
            }
        }
    }
}
