package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.ShareContent
import com.tchat.mobile.components.ShareContentType
import com.tchat.mobile.components.TchatShareModal
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoScreen(
    onVideoClick: (videoId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    var selectedCategory by remember { mutableStateOf("For You") }
    var showShareModal by remember { mutableStateOf(false) }
    var selectedVideo by remember { mutableStateOf<TikTokVideo?>(null) }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar - TikTok Style
        TopAppBar(
            title = {
                Text(
                    "Videos",
                    fontWeight = FontWeight.Bold,
                    fontSize = 18.sp
                )
            },
            actions = {
                IconButton(onClick = { /* Search */ }) {
                    Icon(Icons.Filled.Search, "Search", tint = TchatColors.onBackground)
                }
                IconButton(onClick = { /* Live */ }) {
                    Icon(Icons.Default.Videocam, "Live", tint = TchatColors.error)
                }
                IconButton(onClick = { /* Create */ }) {
                    Icon(Icons.Default.Add, "Create", tint = TchatColors.onBackground)
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.background,
                titleContentColor = TchatColors.onBackground
            )
        )

        // Category Pills - TikTok Style
        LazyRow(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = TchatSpacing.md),
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            items(getTikTokVideoCategories()) { category ->
                Surface(
                    modifier = Modifier.clickable { selectedCategory = category },
                    shape = RoundedCornerShape(20.dp),
                    color = if (selectedCategory == category) TchatColors.onBackground else Color.Transparent,
                    border = androidx.compose.foundation.BorderStroke(
                        1.dp,
                        if (selectedCategory == category) Color.Transparent else TchatColors.outline
                    )
                ) {
                    Text(
                        text = category,
                        modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp),
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium,
                        color = if (selectedCategory == category) TchatColors.background else TchatColors.onBackground
                    )
                }
            }
        }

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        // Video Grid - 2 Columns TikTok Style
        LazyVerticalGrid(
            columns = GridCells.Fixed(2),
            modifier = Modifier.weight(1f),
            contentPadding = PaddingValues(horizontal = TchatSpacing.sm),
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
        ) {
            items(getTikTokVideos().filter {
                selectedCategory == "For You" || it.category == selectedCategory
            }) { video ->
                TikTokVideoCard(
                    video = video,
                    onPlay = { onVideoClick(video.id) },
                    onLike = { /* Like video */ },
                    onComment = { /* Comment */ },
                    onShare = {
                        selectedVideo = video
                        showShareModal = true
                    },
                    onUserClick = { /* Navigate to user */ }
                )
            }
        }
    }

    // Share Modal
    if (showShareModal && selectedVideo != null) {
        TchatShareModal(
            isVisible = showShareModal,
            content = ShareContent(
                title = selectedVideo!!.title,
                description = "Check out this video by ${selectedVideo!!.creator}",
                url = "https://tchat.app/video/${selectedVideo!!.id}",
                type = ShareContentType.GENERAL
            ),
            onDismiss = {
                showShareModal = false
                selectedVideo = null
            },
            onShare = { _, _ ->
                showShareModal = false
                selectedVideo = null
            },
            onCopyLink = {
                // Handle copy link
            }
        )
    }
}

// TikTok-Style Video Card with Horizontal Scrolling Media
@Composable
private fun TikTokVideoCard(
    video: TikTokVideo,
    onPlay: () -> Unit,
    onLike: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    onUserClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Column(modifier = Modifier.fillMaxWidth()) {
            // Horizontal Scrollable Media Section
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(200.dp)
            ) {
                // Horizontal Scrollable Images/Videos
                LazyRow(
                    modifier = Modifier.fillMaxSize(),
                    horizontalArrangement = Arrangement.spacedBy(0.dp)
                ) {
                    items(video.mediaItems) { media ->
                        Box(
                            modifier = Modifier
                                .width(if (video.mediaItems.size == 1) 200.dp else 160.dp)
                                .height(200.dp)
                                .background(
                                    when (media.type) {
                                        MediaType.VIDEO -> TchatColors.primary.copy(alpha = 0.1f)
                                        MediaType.IMAGE -> TchatColors.success.copy(alpha = 0.1f)
                                    }
                                )
                                .clickable {
                                    // Log video click for debugging
                                    println("🎥 Video clicked: ${video.title}")
                                    onPlay()
                                },
                            contentAlignment = Alignment.Center
                        ) {
                            // Media Preview
                            when (media.type) {
                                MediaType.VIDEO -> {
                                    Icon(
                                        Icons.Default.PlayArrow,
                                        contentDescription = "Play Video",
                                        modifier = Modifier.size(48.dp),
                                        tint = TchatColors.primary
                                    )

                                    // Duration badge for videos
                                    Box(
                                        modifier = Modifier
                                            .align(Alignment.BottomEnd)
                                            .padding(6.dp)
                                            .background(
                                                Color.Black.copy(alpha = 0.7f),
                                                RoundedCornerShape(4.dp)
                                            )
                                            .padding(horizontal = 6.dp, vertical = 2.dp)
                                    ) {
                                        Text(
                                            text = media.duration ?: "0:00",
                                            style = MaterialTheme.typography.labelSmall,
                                            color = Color.White,
                                            fontSize = 10.sp
                                        )
                                    }
                                }
                                MediaType.IMAGE -> {
                                    Icon(
                                        Icons.Default.Image,
                                        contentDescription = "Image",
                                        modifier = Modifier.size(48.dp),
                                        tint = TchatColors.success
                                    )
                                }
                            }
                        }
                    }
                }

                // Media Count Indicator
                if (video.mediaItems.size > 1) {
                    Box(
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .padding(8.dp)
                            .background(
                                Color.Black.copy(alpha = 0.6f),
                                RoundedCornerShape(12.dp)
                            )
                            .padding(horizontal = 8.dp, vertical = 4.dp)
                    ) {
                        Text(
                            text = "1/${video.mediaItems.size}",
                            style = MaterialTheme.typography.labelSmall,
                            color = Color.White,
                            fontSize = 10.sp
                        )
                    }
                }
            }

            // Video Info Section
            Column(
                modifier = Modifier.padding(12.dp)
            ) {
                // Title with 2 lines max
                Text(
                    text = video.title,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis,
                    lineHeight = 18.sp
                )

                Spacer(modifier = Modifier.height(6.dp))

                // Hashtags Row (ป้ายยา)
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(6.dp)
                ) {
                    items(video.hashtags.take(3)) { hashtag ->
                        Surface(
                            shape = RoundedCornerShape(12.dp),
                            color = TchatColors.primary.copy(alpha = 0.1f)
                        ) {
                            Text(
                                text = hashtag,
                                modifier = Modifier.padding(horizontal = 8.dp, vertical = 3.dp),
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.primary,
                                fontSize = 10.sp
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Creator and Stats Row
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Creator Avatar
                    Box(
                        modifier = Modifier
                            .size(20.dp)
                            .clip(CircleShape)
                            .background(TchatColors.primary)
                            .clickable { onUserClick() },
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = video.creator.first().toString(),
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onPrimary,
                            fontSize = 10.sp,
                            fontWeight = FontWeight.Bold
                        )
                    }

                    Spacer(modifier = Modifier.width(6.dp))

                    Text(
                        text = video.creator,
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant,
                        fontSize = 11.sp,
                        modifier = Modifier.weight(1f)
                    )

                    Text(
                        text = formatCount(video.views),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant,
                        fontSize = 10.sp
                    )
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Action Buttons Row
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Like Button
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier.clickable { onLike() }
                    ) {
                        Icon(
                            if (video.isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                            contentDescription = "Like",
                            modifier = Modifier.size(16.dp),
                            tint = if (video.isLiked) TchatColors.error else TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatCount(video.likes),
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onSurfaceVariant,
                            fontSize = 10.sp
                        )
                    }

                    // Comment Button
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier.clickable { onComment() }
                    ) {
                        Icon(
                            Icons.Default.ChatBubbleOutline,
                            contentDescription = "Comment",
                            modifier = Modifier.size(16.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatCount(video.comments),
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onSurfaceVariant,
                            fontSize = 10.sp
                        )
                    }

                    // Share Button
                    Icon(
                        Icons.Default.Share,
                        contentDescription = "Share",
                        modifier = Modifier
                            .size(16.dp)
                            .clickable { onShare() },
                        tint = TchatColors.onSurfaceVariant
                    )

                    // Bookmark Button
                    Icon(
                        if (video.isBookmarked) Icons.Default.Bookmark else Icons.Default.BookmarkBorder,
                        contentDescription = "Bookmark",
                        modifier = Modifier.size(16.dp),
                        tint = if (video.isBookmarked) TchatColors.primary else TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

// Helper function to format large numbers
private fun formatCount(count: Int): String {
    return when {
        count >= 1_000_000 -> "${count / 1_000_000}M"
        count >= 1_000 -> "${count / 1_000}K"
        else -> count.toString()
    }
}

// TikTok-Style Data Models
enum class MediaType { VIDEO, IMAGE }

data class MediaItem(
    val id: String,
    val type: MediaType,
    val url: String? = null,
    val duration: String? = null // Only for videos
)

data class TikTokVideo(
    val id: String,
    val title: String,
    val creator: String,
    val mediaItems: List<MediaItem>, // Horizontal scrollable media
    val hashtags: List<String>, // ป้ายยา
    val views: Int,
    val likes: Int,
    val comments: Int,
    val category: String,
    val uploadTime: String,
    val isLiked: Boolean = false,
    val isBookmarked: Boolean = false
)

private fun getTikTokVideoCategories(): List<String> = listOf(
    "For You", "Following", "Live", "Gaming", "Food", "Beauty", "Sports", "Music", "Education", "Comedy"
)

private fun getTikTokVideos(): List<TikTokVideo> = listOf(
    TikTokVideo(
        id = "1",
        title = "Amazing Thai Street Food! Must try in Bangkok 🇹🇭",
        creator = "FoodieThailand",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "0:45"),
            MediaItem("2", MediaType.IMAGE),
            MediaItem("3", MediaType.VIDEO, duration = "1:12")
        ),
        hashtags = listOf("#อาหารไทย", "#กรุงเทพ", "#อร่อย", "#streetfood", "#thailand"),
        views = 1250000,
        likes = 89500,
        comments = 2340,
        category = "Food",
        uploadTime = "2h"
    ),
    TikTokVideo(
        id = "2",
        title = "React Native Tutorial: Build Amazing Apps",
        creator = "CodeMaster",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "8:30")
        ),
        hashtags = listOf("#reactnative", "#coding", "#tutorial", "#mobile", "#programming"),
        views = 456000,
        likes = 23400,
        comments = 890,
        category = "Education",
        uploadTime = "4h"
    ),
    TikTokVideo(
        id = "3",
        title = "Epic Gaming Moments - You Won't Believe This! 🎮",
        creator = "ProGamer",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "2:15"),
            MediaItem("2", MediaType.IMAGE),
            MediaItem("3", MediaType.VIDEO, duration = "1:45"),
            MediaItem("4", MediaType.IMAGE)
        ),
        hashtags = listOf("#gaming", "#epic", "#pro", "#moments", "#incredible"),
        views = 892000,
        likes = 67800,
        comments = 1560,
        category = "Gaming",
        uploadTime = "6h"
    ),
    TikTokVideo(
        id = "4",
        title = "Makeup Tutorial: Korean Glass Skin Look ✨",
        creator = "BeautyGuru",
        mediaItems = listOf(
            MediaItem("1", MediaType.IMAGE),
            MediaItem("2", MediaType.VIDEO, duration = "3:20"),
            MediaItem("3", MediaType.IMAGE)
        ),
        hashtags = listOf("#makeup", "#korean", "#skincare", "#beauty", "#tutorial", "#glasskin"),
        views = 2100000,
        likes = 156000,
        comments = 4580,
        category = "Beauty",
        uploadTime = "1d"
    ),
    TikTokVideo(
        id = "5",
        title = "Football Skills That Will Blow Your Mind! ⚽",
        creator = "FootballPro",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "1:50"),
            MediaItem("2", MediaType.VIDEO, duration = "0:30")
        ),
        hashtags = listOf("#football", "#soccer", "#skills", "#amazing", "#sports"),
        views = 678000,
        likes = 45600,
        comments = 1230,
        category = "Sports",
        uploadTime = "12h"
    ),
    TikTokVideo(
        id = "6",
        title = "Thai Language Learning Made Easy! สวัสดี",
        creator = "ThaiTeacher",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "5:12"),
            MediaItem("2", MediaType.IMAGE),
            MediaItem("3", MediaType.IMAGE),
            MediaItem("4", MediaType.VIDEO, duration = "2:30")
        ),
        hashtags = listOf("#ไทย", "#เรียนภาษา", "#thai", "#language", "#learning", "#easy"),
        views = 234000,
        likes = 12800,
        comments = 567,
        category = "Education",
        uploadTime = "1d"
    ),
    TikTokVideo(
        id = "7",
        title = "Funny Cat Compilation - You'll Laugh All Day 😹",
        creator = "FunnyCats",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "4:15"),
            MediaItem("2", MediaType.IMAGE),
            MediaItem("3", MediaType.VIDEO, duration = "1:20"),
            MediaItem("4", MediaType.IMAGE),
            MediaItem("5", MediaType.VIDEO, duration = "0:45")
        ),
        hashtags = listOf("#funny", "#cats", "#pets", "#comedy", "#cute", "#animals"),
        views = 3200000,
        likes = 245000,
        comments = 8900,
        category = "Comedy",
        uploadTime = "2d"
    ),
    TikTokVideo(
        id = "8",
        title = "Live Music Performance - Acoustic Session 🎵",
        creator = "MusicLive",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "12:30")
        ),
        hashtags = listOf("#music", "#live", "#acoustic", "#performance", "#original"),
        views = 567000,
        likes = 34500,
        comments = 1890,
        category = "Music",
        uploadTime = "3h",
        isLiked = true
    ),
    TikTokVideo(
        id = "9",
        title = "Home Workout: No Equipment Needed! 💪",
        creator = "FitnessCoach",
        mediaItems = listOf(
            MediaItem("1", MediaType.VIDEO, duration = "15:00"),
            MediaItem("2", MediaType.IMAGE),
            MediaItem("3", MediaType.VIDEO, duration = "3:45")
        ),
        hashtags = listOf("#fitness", "#workout", "#home", "#health", "#exercise", "#noequipment"),
        views = 1100000,
        likes = 78900,
        comments = 2340,
        category = "Sports",
        uploadTime = "5h",
        isBookmarked = true
    ),
    TikTokVideo(
        id = "10",
        title = "Travel Vlog: Hidden Gems in Singapore 🇸🇬",
        creator = "TravelSG",
        mediaItems = listOf(
            MediaItem("1", MediaType.IMAGE),
            MediaItem("2", MediaType.VIDEO, duration = "6:20"),
            MediaItem("3", MediaType.IMAGE),
            MediaItem("4", MediaType.VIDEO, duration = "2:10"),
            MediaItem("5", MediaType.IMAGE)
        ),
        hashtags = listOf("#singapore", "#travel", "#hidden", "#gems", "#vlog", "#explore"),
        views = 890000,
        likes = 56700,
        comments = 1670,
        category = "For You",
        uploadTime = "8h"
    )
)