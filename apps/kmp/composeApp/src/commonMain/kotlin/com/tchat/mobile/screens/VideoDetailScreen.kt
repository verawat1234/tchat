package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoDetailScreen(
    videoId: String = "1",
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    val video = getDummyVideoDetail(videoId)
    var isLiked by remember { mutableStateOf(video.isLiked) }
    var isBookmarked by remember { mutableStateOf(video.isBookmarked) }
    var isSubscribed by remember { mutableStateOf(video.isSubscribed) }
    var commentText by remember { mutableStateOf("") }
    var showComments by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = video.title,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, "Back")
                    }
                },
                actions = {
                    IconButton(onClick = { /* Share video */ }) {
                        Icon(Icons.Default.Share, "Share")
                    }
                    IconButton(onClick = { /* More options */ }) {
                        Icon(Icons.Default.MoreVert, "More")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface
                )
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .background(TchatColors.background),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            // Video Player
            item {
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(0.dp),
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(220.dp)
                            .background(TchatColors.surfaceVariant),
                        contentAlignment = Alignment.Center
                    ) {
                        // Video preview placeholder
                        Box(
                            modifier = Modifier
                                .fillMaxSize()
                                .background(TchatColors.primary.copy(alpha = 0.1f)),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Filled.PlayArrow,
                                contentDescription = "Play Video",
                                modifier = Modifier.size(80.dp),
                                tint = TchatColors.primary
                            )
                        }

                        // Duration badge
                        Box(
                            modifier = Modifier
                                .align(Alignment.BottomEnd)
                                .padding(TchatSpacing.sm)
                                .background(
                                    TchatColors.surface.copy(alpha = 0.9f),
                                    RoundedCornerShape(4.dp)
                                )
                                .padding(horizontal = 8.dp, vertical = 4.dp)
                        ) {
                            Text(
                                text = video.duration,
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.onSurface
                            )
                        }
                    }
                }
            }

            // Video Info
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.md)
                    ) {
                        // Title
                        Text(
                            text = video.title,
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )

                        Spacer(modifier = Modifier.height(TchatSpacing.xs))

                        // Views and upload time
                        Text(
                            text = "${formatCount(video.views)} views • ${video.uploadTime}",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Action buttons
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween
                        ) {
                            ActionButton(
                                icon = if (isLiked) Icons.Default.ThumbUp else Icons.Default.ThumbUp,
                                label = formatCount(video.likesCount + if (isLiked) 1 else 0),
                                isActive = isLiked,
                                onClick = { isLiked = !isLiked }
                            )

                            ActionButton(
                                icon = Icons.Default.Clear,
                                label = "Dislike",
                                isActive = false,
                                onClick = { /* Handle dislike */ }
                            )

                            ActionButton(
                                icon = Icons.Default.Share,
                                label = "Share",
                                isActive = false,
                                onClick = { /* Handle share */ }
                            )

                            ActionButton(
                                icon = if (isBookmarked) Icons.Default.Star else Icons.Default.Star,
                                label = "Save",
                                isActive = isBookmarked,
                                onClick = { isBookmarked = !isBookmarked }
                            )
                        }
                    }
                }
            }

            // Creator Info
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.md),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        // Creator avatar
                        Box(
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .background(TchatColors.primary),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = video.creatorName.first().toString(),
                                style = MaterialTheme.typography.titleMedium,
                                color = TchatColors.onPrimary,
                                fontWeight = FontWeight.Bold
                            )
                        }

                        Spacer(modifier = Modifier.width(TchatSpacing.sm))

                        Column(
                            modifier = Modifier.weight(1f)
                        ) {
                            Text(
                                text = video.creatorName,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.SemiBold,
                                color = TchatColors.onSurface
                            )
                            Text(
                                text = "${formatCount(video.subscribersCount)} subscribers",
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }

                        TchatButton(
                            onClick = { isSubscribed = !isSubscribed },
                            text = if (isSubscribed) "Subscribed" else "Subscribe",
                            variant = if (isSubscribed) TchatButtonVariant.Secondary else TchatButtonVariant.Primary
                        )
                    }
                }
            }

            // Description
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.md)
                    ) {
                        Text(
                            text = "Description",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.SemiBold,
                            color = TchatColors.onSurface
                        )

                        Spacer(modifier = Modifier.height(TchatSpacing.sm))

                        Text(
                            text = video.description,
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface
                        )
                    }
                }
            }

            // Comments Section
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.md)
                    ) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(
                                text = "Comments (${video.commentsCount})",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.SemiBold,
                                color = TchatColors.onSurface
                            )

                            TextButton(
                                onClick = { showComments = !showComments }
                            ) {
                                Text(
                                    text = if (showComments) "Hide" else "Show All",
                                    color = TchatColors.primary
                                )
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.sm))

                        // Comment input
                        Row(
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            TchatInput(
                                value = commentText,
                                onValueChange = { commentText = it },
                                type = TchatInputType.Text,
                                placeholder = "Add a comment...",
                                modifier = Modifier.weight(1f)
                            )

                            Spacer(modifier = Modifier.width(TchatSpacing.sm))

                            TchatButton(
                                onClick = {
                                    if (commentText.isNotBlank()) {
                                        // Add comment
                                        commentText = ""
                                    }
                                },
                                text = "Post",
                                variant = TchatButtonVariant.Primary,
                                enabled = commentText.isNotBlank()
                            )
                        }
                    }
                }
            }

            // Comments List (if shown)
            if (showComments) {
                items(video.topComments) { comment ->
                    CommentCard(comment = comment)
                }
            }

            // Related Videos
            item {
                Text(
                    text = "Related Videos",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(horizontal = TchatSpacing.md)
                )
            }

            item {
                LazyRow(
                    contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(video.relatedVideos) { relatedVideo ->
                        RelatedVideoCard(
                            video = relatedVideo,
                            onClick = { /* Navigate to video */ }
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun ActionButton(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    label: String,
    isActive: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        IconButton(onClick = onClick) {
            Icon(
                imageVector = icon,
                contentDescription = label,
                tint = if (isActive) TchatColors.primary else TchatColors.onSurfaceVariant
            )
        }
        Text(
            text = label,
            style = MaterialTheme.typography.labelSmall,
            color = if (isActive) TchatColors.primary else TchatColors.onSurfaceVariant
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

private fun getDummyVideoDetail(videoId: String): VideoDetail = VideoDetail(
    id = videoId,
    title = "Amazing AI Demo - The Future is Here!",
    creatorName = "TechGuru",
    subscribersCount = 250000,
    duration = "5:42",
    views = 2100000,
    uploadTime = "2 hours ago",
    likesCount = 45600,
    commentsCount = 892,
    description = "In this video, I'll show you the most incredible AI demonstrations that will blow your mind! From natural language processing to computer vision, these AI tools are revolutionizing how we work and live.\n\n🔗 Links mentioned in the video:\n- AI Tool 1: example.com\n- AI Tool 2: example.com\n\n⏰ Timestamps:\n0:00 Introduction\n1:30 Natural Language AI\n3:15 Computer Vision\n4:20 Conclusion\n\n#AI #Technology #Future #Demo",
    isLiked = true,
    isBookmarked = false,
    isSubscribed = false,
    topComments = listOf(
        VideoComment("1", "Alice Cooper", "This is absolutely mind-blowing! Thanks for sharing these amazing tools.", "3h", 156),
        VideoComment("2", "Bob Smith", "The AI landscape is evolving so fast. Great breakdown of the key technologies!", "2h", 89),
        VideoComment("3", "Carol Johnson", "Can you do a tutorial on how to implement some of these AI features?", "1h", 203),
        VideoComment("4", "David Lee", "Amazing content as always. Looking forward to the next video!", "45m", 67),
        VideoComment("5", "Emma Wilson", "The computer vision demo at 3:15 was incredible. How accurate is it in real-world scenarios?", "30m", 124)
    ),
    relatedVideos = listOf(
        RelatedVideo("2", "Machine Learning Basics Explained", "AI Academy", "12:34", 850000),
        RelatedVideo("3", "Building Your First AI App", "CodeMaster", "18:22", 450000),
        RelatedVideo("4", "The Future of Artificial Intelligence", "FutureTech", "8:15", 1200000),
        RelatedVideo("5", "AI Tools for Developers", "DevGuru", "15:30", 680000),
        RelatedVideo("6", "Deep Learning Made Simple", "TechExplainer", "22:18", 930000)
    )
)