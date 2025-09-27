package com.tchat.mobile.components.posts

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.models.*
import com.tchat.mobile.services.NavigationService
import com.tchat.mobile.services.SharingService
import com.tchat.mobile.services.SharingPlatform
import com.tchat.mobile.repositories.PostRepository
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import kotlinx.coroutines.launch

/**
 * Post Renderer - Universal Social Media Post Component
 *
 * Handles all post types: Reviews, Social Posts, Videos, Stories, etc.
 * Integrates with repository, sharing, and navigation services
 */

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PostRenderer(
    post: Post,
    postRepository: PostRepository,
    sharingService: SharingService,
    navigationService: NavigationService,
    onPostClick: ((Post) -> Unit)? = null,
    modifier: Modifier = Modifier
) {
    var isLiked by remember { mutableStateOf(post.interactions.isLiked) }
    var likeCount by remember { mutableIntStateOf(post.interactions.reactions.count { it.type == ReactionType.LIKE }) }
    var isBookmarked by remember { mutableStateOf(post.interactions.isBookmarked) }
    var showShareMenu by remember { mutableStateOf(false) }

    val scope = rememberCoroutineScope()

    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onPostClick?.invoke(post) },
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surface
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            // User Header
            PostUserHeader(
                user = post.user,
                createdAt = post.createdAt,
                location = post.content.location,
                isSponsored = post.metadata?.isSponsored == true,
                sponsorName = post.metadata?.sponsorName,
                onUserClick = {
                    scope.launch {
                        navigationService.navigateToUserProfile(post.user.id)
                    }
                }
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            // Post Content
            PostContent(
                content = post.content,
                postType = post.type,
                onImageClick = { image ->
                    scope.launch {
                        navigationService.navigateToImageViewer(image.url, post.id)
                    }
                },
                onVideoClick = { video ->
                    scope.launch {
                        navigationService.navigateToVideoPlayer(video.url, post.id)
                    }
                },
                onHashtagClick = { hashtag ->
                    scope.launch {
                        navigationService.navigateToHashtagFeed(hashtag)
                    }
                }
            )

            // Post Metadata (for reviews, products, etc.)
            if (post.metadata != null) {
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
                PostMetadataSection(
                    metadata = post.metadata,
                    onTargetClick = { targetType, targetId ->
                        scope.launch {
                            when (targetType.lowercase()) {
                                "product" -> navigationService.navigateToProductDetail(targetId)
                                "shop" -> navigationService.navigateToShopDetail(targetId)
                                "user" -> navigationService.navigateToUserProfile(targetId)
                            }
                        }
                    }
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            // Interaction Bar
            PostInteractionBar(
                interactions = post.interactions.copy(
                    isLiked = isLiked,
                    isBookmarked = isBookmarked
                ),
                onLikeClick = {
                    scope.launch {
                        val result = if (isLiked) {
                            postRepository.unlikePost(post.id)
                        } else {
                            postRepository.likePost(post.id)
                        }

                        result.onSuccess { updatedInteractions ->
                            isLiked = updatedInteractions.isLiked
                            likeCount = updatedInteractions.reactions.count { it.type == ReactionType.LIKE }
                        }
                    }
                },
                onCommentClick = {
                    scope.launch {
                        navigationService.navigateToComments(post.id)
                    }
                },
                onShareClick = {
                    showShareMenu = true
                },
                onBookmarkClick = {
                    scope.launch {
                        val result = postRepository.bookmarkPost(post.id)
                        result.onSuccess { bookmarked ->
                            isBookmarked = bookmarked
                        }
                    }
                }
            )
        }
    }

    // Share Menu
    if (showShareMenu) {
        PostShareMenu(
            post = post,
            sharingService = sharingService,
            onDismiss = { showShareMenu = false },
            onShare = { platform ->
                scope.launch {
                    sharingService.sharePost(post, platform)
                    showShareMenu = false
                }
            }
        )
    }
}

@Composable
private fun PostUserHeader(
    user: PostUser,
    createdAt: String,
    location: String?,
    isSponsored: Boolean,
    sponsorName: String?,
    onUserClick: () -> Unit
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // User Avatar
        Box(
            modifier = Modifier
                .size(40.dp)
                .clip(CircleShape)
                .background(TchatColors.primary)
                .clickable { onUserClick() },
            contentAlignment = Alignment.Center
        ) {
            Icon(
                imageVector = Icons.Default.Person,
                contentDescription = "User Avatar",
                tint = Color.White,
                modifier = Modifier.size(24.dp)
            )
        }

        Spacer(modifier = Modifier.width(TchatSpacing.sm))

        Column(modifier = Modifier.weight(1f)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    text = user.displayName ?: user.username,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )

                if (user.isVerified) {
                    Spacer(modifier = Modifier.width(4.dp))
                    Icon(
                        imageVector = Icons.Default.Verified,
                        contentDescription = "Verified",
                        tint = Color(0xFF1DA1F2),
                        modifier = Modifier.size(16.dp)
                    )
                }

                if (isSponsored) {
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(
                        text = "Sponsored",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant,
                        fontSize = 10.sp
                    )
                }
            }

            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    text = createdAt,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                if (location != null) {
                    Text(
                        text = " • $location",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }

        IconButton(onClick = { /* More options */ }) {
            Icon(
                imageVector = Icons.Default.MoreVert,
                contentDescription = "More options",
                tint = TchatColors.onSurfaceVariant
            )
        }
    }
}

@Composable
private fun PostContent(
    content: PostContent,
    postType: PostType,
    onImageClick: (PostImage) -> Unit,
    onVideoClick: (PostVideo) -> Unit,
    onHashtagClick: (String) -> Unit
) {
    Column {
        // Text content
        if (!content.text.isNullOrBlank()) {
            PostTextContent(
                text = content.text,
                onHashtagClick = onHashtagClick
            )

            if (content.images.isNotEmpty() || content.videos.isNotEmpty()) {
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
            }
        }

        // Media content
        when (content.type) {
            PostContentType.IMAGE -> {
                PostImageGrid(
                    images = content.images,
                    onImageClick = onImageClick
                )
            }

            PostContentType.VIDEO -> {
                PostVideoPlayer(
                    videos = content.videos,
                    onVideoClick = onVideoClick
                )
            }

            PostContentType.MIXED -> {
                if (content.images.isNotEmpty()) {
                    PostImageGrid(
                        images = content.images,
                        onImageClick = onImageClick
                    )

                    if (content.videos.isNotEmpty()) {
                        Spacer(modifier = Modifier.height(TchatSpacing.sm))
                    }
                }

                if (content.videos.isNotEmpty()) {
                    PostVideoPlayer(
                        videos = content.videos,
                        onVideoClick = onVideoClick
                    )
                }
            }

            PostContentType.POLL -> {
                if (content.poll != null) {
                    PostPoll(poll = content.poll)
                }
            }

            else -> { /* Text only - already handled above */ }
        }
    }
}

@Composable
private fun PostTextContent(
    text: String,
    onHashtagClick: (String) -> Unit
) {
    // Simple text rendering with hashtag detection
    // In a real app, you'd use a more sophisticated text parser
    Text(
        text = text,
        style = MaterialTheme.typography.bodyMedium,
        color = TchatColors.onSurface,
        modifier = Modifier.clickable {
            // Extract hashtags and handle clicks
            val hashtags = text.split(" ").filter { it.startsWith("#") }
            hashtags.forEach { hashtag ->
                if (hashtag.isNotEmpty()) {
                    onHashtagClick(hashtag)
                }
            }
        }
    )
}

@Composable
private fun PostImageGrid(
    images: List<PostImage>,
    onImageClick: (PostImage) -> Unit
) {
    when (images.size) {
        1 -> {
            PostImageItem(
                image = images[0],
                modifier = Modifier
                    .fillMaxWidth()
                    .height(200.dp),
                onImageClick = onImageClick
            )
        }

        else -> {
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items(images) { image ->
                    PostImageItem(
                        image = image,
                        modifier = Modifier
                            .width(150.dp)
                            .height(150.dp),
                        onImageClick = onImageClick
                    )
                }
            }
        }
    }
}

@Composable
private fun PostImageItem(
    image: PostImage,
    onImageClick: (PostImage) -> Unit,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier
            .clip(RoundedCornerShape(8.dp))
            .background(Color.Gray)
            .clickable { onImageClick(image) },
        contentAlignment = Alignment.Center
    ) {
        Icon(
            imageVector = Icons.Default.Image,
            contentDescription = "Post Image",
            tint = Color.White,
            modifier = Modifier.size(48.dp)
        )
    }
}

@Composable
private fun PostVideoPlayer(
    videos: List<PostVideo>,
    onVideoClick: (PostVideo) -> Unit
) {
    LazyRow(
        horizontalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        items(videos) { video ->
            Box(
                modifier = Modifier
                    .width(200.dp)
                    .height(150.dp)
                    .clip(RoundedCornerShape(8.dp))
                    .background(Color.Black)
                    .clickable { onVideoClick(video) },
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = Icons.Default.PlayArrow,
                    contentDescription = "Play Video",
                    tint = Color.White,
                    modifier = Modifier.size(48.dp)
                )

                // Duration overlay
                Box(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(8.dp)
                        .background(
                            Color.Black.copy(alpha = 0.7f),
                            RoundedCornerShape(4.dp)
                        )
                        .padding(horizontal = 4.dp, vertical = 2.dp)
                ) {
                    Text(
                        text = video.duration,
                        color = Color.White,
                        style = MaterialTheme.typography.bodySmall,
                        fontSize = 10.sp
                    )
                }
            }
        }
    }
}

@Composable
private fun PostPoll(poll: PostPoll) {
    Column {
        Text(
            text = poll.question,
            style = MaterialTheme.typography.titleSmall,
            fontWeight = FontWeight.Bold
        )

        Spacer(modifier = Modifier.height(8.dp))

        poll.options.forEachIndexed { index, option ->
            val votes = poll.votes[index] ?: 0
            val totalVotes = poll.votes.values.sum()
            val percentage = if (totalVotes > 0) (votes.toFloat() / totalVotes * 100).toInt() else 0

            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 2.dp)
                    .clickable { /* Handle vote */ },
                colors = CardDefaults.cardColors(
                    containerColor = TchatColors.surfaceVariant
                )
            ) {
                Row(
                    modifier = Modifier.padding(12.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = option,
                        modifier = Modifier.weight(1f)
                    )
                    Text(
                        text = "$percentage%",
                        color = TchatColors.onSurfaceVariant,
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            }
        }
    }
}

@Composable
private fun PostMetadataSection(
    metadata: PostMetadata,
    onTargetClick: (String, String) -> Unit
) {
    Card(
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surfaceVariant
        ),
        modifier = Modifier.fillMaxWidth()
    ) {
        Row(
            modifier = Modifier
                .padding(TchatSpacing.sm)
                .clickable {
                    if (metadata.targetType != null && metadata.targetId != null) {
                        onTargetClick(metadata.targetType, metadata.targetId)
                    }
                },
            verticalAlignment = Alignment.CenterVertically
        ) {
            if (metadata.rating != null) {
                repeat(5) { index ->
                    Icon(
                        imageVector = if (index < metadata.rating * 5) Icons.Default.Star else Icons.Default.StarBorder,
                        contentDescription = "Rating",
                        tint = Color(0xFFFFB000),
                        modifier = Modifier.size(16.dp)
                    )
                }
                Spacer(modifier = Modifier.width(8.dp))
            }

            Text(
                text = metadata.targetName ?: "Unknown Target",
                style = MaterialTheme.typography.bodyMedium,
                fontWeight = FontWeight.Medium,
                maxLines = 1,
                overflow = TextOverflow.Ellipsis
            )
        }
    }
}

@Composable
private fun PostInteractionBar(
    interactions: PostInteractions,
    onLikeClick: () -> Unit,
    onCommentClick: () -> Unit,
    onShareClick: () -> Unit,
    onBookmarkClick: () -> Unit
) {
    var showReactionMenu by remember { mutableStateOf(false) }

    Column {
        // Engagement metrics (views, reach)
        if (interactions.views > 0 || interactions.reach > 0) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                if (interactions.views > 0) {
                    Text(
                        text = "${formatEngagementCount(interactions.views)} views",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
                if (interactions.reach > 0 && interactions.reach != interactions.views) {
                    Text(
                        text = "${formatEngagementCount(interactions.reach)} reached",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
            Spacer(modifier = Modifier.height(8.dp))
        }

        // Main interaction bar
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Row(
                horizontalArrangement = Arrangement.spacedBy(20.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Enhanced Like/Reaction button
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.clickable { showReactionMenu = true }
                ) {
                    val dominantReaction = interactions.reactions
                        .groupBy { it.type }
                        .maxByOrNull { it.value.size }

                    Icon(
                        imageVector = when {
                            interactions.isLiked -> Icons.Default.Favorite
                            dominantReaction?.key == ReactionType.LOVE -> Icons.Default.FavoriteBorder
                            dominantReaction?.key == ReactionType.FIRE -> Icons.Default.Whatshot
                            dominantReaction?.key == ReactionType.CLAP -> Icons.Default.PanTool
                            else -> Icons.Default.FavoriteBorder
                        },
                        contentDescription = "React",
                        tint = when {
                            interactions.isLiked -> Color.Red
                            dominantReaction?.key == ReactionType.LOVE -> Color(0xFFFF69B4)
                            dominantReaction?.key == ReactionType.FIRE -> Color(0xFFFF4500)
                            dominantReaction?.key == ReactionType.CLAP -> Color(0xFFFFD700)
                            else -> TchatColors.onSurfaceVariant
                        },
                        modifier = Modifier.size(22.dp)
                    )
                    if (interactions.reactions.isNotEmpty()) {
                        Spacer(modifier = Modifier.width(6.dp))
                        Text(
                            text = formatEngagementCount(interactions.reactions.size),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant,
                            fontWeight = if (interactions.isLiked) FontWeight.Medium else FontWeight.Normal
                        )
                    }
                }

                // Comment button
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.clickable { onCommentClick() }
                ) {
                    Icon(
                        imageVector = Icons.Default.ChatBubbleOutline,
                        contentDescription = "Comment",
                        tint = TchatColors.onSurfaceVariant,
                        modifier = Modifier.size(22.dp)
                    )
                    if (interactions.comments.isNotEmpty()) {
                        Spacer(modifier = Modifier.width(6.dp))
                        Text(
                            text = formatEngagementCount(interactions.comments.size),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                // Share button
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.clickable { onShareClick() }
                ) {
                    Icon(
                        imageVector = Icons.Default.Share,
                        contentDescription = "Share",
                        tint = TchatColors.onSurfaceVariant,
                        modifier = Modifier.size(22.dp)
                    )
                    if (interactions.shares.isNotEmpty()) {
                        Spacer(modifier = Modifier.width(6.dp))
                        Text(
                            text = formatEngagementCount(interactions.shares.size),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                // Send/DM button (TikTok/Instagram style)
                Icon(
                    imageVector = Icons.Default.Send,
                    contentDescription = "Send",
                    tint = TchatColors.onSurfaceVariant,
                    modifier = Modifier
                        .size(22.dp)
                        .clickable { /* TODO: Implement direct message sharing */ }
                )
            }

            // Bookmark button (right aligned like Instagram)
            IconButton(onClick = onBookmarkClick) {
                Icon(
                    imageVector = if (interactions.isBookmarked) Icons.Default.Bookmark else Icons.Default.BookmarkBorder,
                    contentDescription = "Save",
                    tint = if (interactions.isBookmarked) TchatColors.primary else TchatColors.onSurfaceVariant,
                    modifier = Modifier.size(22.dp)
                )
            }
        }

        // Engagement summary (Instagram style)
        if (interactions.reactions.isNotEmpty() || interactions.comments.isNotEmpty()) {
            Spacer(modifier = Modifier.height(8.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                if (interactions.reactions.isNotEmpty()) {
                    Text(
                        text = "${formatEngagementCount(interactions.reactions.size)} likes",
                        style = MaterialTheme.typography.bodySmall,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.onSurface
                    )
                }
                if (interactions.comments.isNotEmpty()) {
                    if (interactions.reactions.isNotEmpty()) {
                        Text(
                            text = " • ",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                    Text(
                        text = "View all ${formatEngagementCount(interactions.comments.size)} comments",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant,
                        modifier = Modifier.clickable { onCommentClick() }
                    )
                }
            }
        }
    }

    // Reaction Menu (similar to Facebook/Instagram)
    if (showReactionMenu) {
        ReactionMenu(
            onReactionSelected = { reaction ->
                onLikeClick() // For now, just trigger like
                showReactionMenu = false
            },
            onDismiss = { showReactionMenu = false }
        )
    }
}

private fun formatEngagementCount(count: Int): String {
    return when {
        count < 1000 -> count.toString()
        count < 1000000 -> {
            val k = count / 1000f
            if (k % 1 == 0f) "${k.toInt()}K" else "${(k * 10).toInt() / 10f}K"
        }
        else -> {
            val m = count / 1000000f
            if (m % 1 == 0f) "${m.toInt()}M" else "${(m * 10).toInt() / 10f}M"
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ReactionMenu(
    onReactionSelected: (ReactionType) -> Unit,
    onDismiss: () -> Unit
) {
    ModalBottomSheet(
        onDismissRequest = onDismiss
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Text(
                text = "Choose Reaction",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Medium
            )

            Spacer(modifier = Modifier.height(16.dp))

            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(16.dp),
                modifier = Modifier.fillMaxWidth()
            ) {
                items(ReactionType.values()) { reaction ->
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally,
                        modifier = Modifier.clickable { onReactionSelected(reaction) }
                    ) {
                        Icon(
                            imageVector = when (reaction) {
                                ReactionType.LIKE -> Icons.Default.Favorite
                                ReactionType.LOVE -> Icons.Default.FavoriteBorder
                                ReactionType.HAHA -> Icons.Default.SentimentSatisfied
                                ReactionType.WOW -> Icons.Default.SentimentVeryDissatisfied
                                ReactionType.SAD -> Icons.Default.SentimentDissatisfied
                                ReactionType.ANGRY -> Icons.Default.SentimentVeryDissatisfied
                                ReactionType.CARE -> Icons.Default.Favorite
                                ReactionType.FIRE -> Icons.Default.Whatshot
                                ReactionType.CLAP -> Icons.Default.PanTool
                                ReactionType.CELEBRATE -> Icons.Default.Celebration
                            },
                            contentDescription = reaction.name,
                            tint = when (reaction) {
                                ReactionType.LIKE -> Color.Red
                                ReactionType.LOVE -> Color(0xFFFF69B4)
                                ReactionType.HAHA -> Color(0xFFFFD700)
                                ReactionType.WOW -> Color(0xFF87CEEB)
                                ReactionType.SAD -> Color(0xFF4169E1)
                                ReactionType.ANGRY -> Color(0xFFFF4500)
                                ReactionType.CARE -> Color(0xFFFF69B4)
                                ReactionType.FIRE -> Color(0xFFFF4500)
                                ReactionType.CLAP -> Color(0xFFFFD700)
                                ReactionType.CELEBRATE -> Color(0xFF32CD32)
                            },
                            modifier = Modifier.size(32.dp)
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(
                            text = reaction.name.lowercase().replaceFirstChar { it.uppercase() },
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.height(32.dp))
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun PostShareMenu(
    post: Post,
    sharingService: SharingService,
    onDismiss: () -> Unit,
    onShare: (SharingPlatform) -> Unit
) {
    var availablePlatforms by remember { mutableStateOf<List<SharingPlatform>>(emptyList()) }

    LaunchedEffect(Unit) {
        availablePlatforms = sharingService.getAvailablePlatforms()
    }

    ModalBottomSheet(
        onDismissRequest = onDismiss
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Text(
                text = "Share Post",
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold
            )

            Spacer(modifier = Modifier.height(16.dp))

            availablePlatforms.forEach { platform ->
                ListItem(
                    headlineContent = { Text(platform.displayName) },
                    leadingContent = {
                        Icon(
                            imageVector = when (platform) {
                                SharingPlatform.WHATSAPP -> Icons.Default.Chat
                                SharingPlatform.TWITTER -> Icons.Default.Message
                                SharingPlatform.FACEBOOK -> Icons.Default.Share
                                SharingPlatform.INSTAGRAM -> Icons.Default.Photo
                                SharingPlatform.COPY_LINK -> Icons.Default.Link
                                else -> Icons.Default.Share
                            },
                            contentDescription = platform.displayName
                        )
                    },
                    modifier = Modifier.clickable { onShare(platform) }
                )
            }

            Spacer(modifier = Modifier.height(32.dp))
        }
    }
}