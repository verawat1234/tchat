package com.tchat.mobile.components.reviews

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
import com.tchat.mobile.components.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*

/**
 * Text Review Component
 */
@Composable
fun TextReview(
    review: Review,
    onLike: () -> Unit = {},
    onComment: () -> Unit = {},
    onShare: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            // Rating Stars
            ReviewRatingDisplay(
                rating = review.rating,
                modifier = Modifier.padding(bottom = TchatSpacing.sm)
            )

            // Review Text
            if (!review.content.text.isNullOrBlank()) {
                Text(
                    text = review.content.text,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Hashtags
            if (review.content.hashtags.isNotEmpty()) {
                ReviewHashtags(
                    hashtags = review.content.hashtags,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Author and Actions
            ReviewAuthorSection(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare
            )
        }
    }
}

/**
 * Image Review Component
 */
@Composable
fun ImageReview(
    review: Review,
    onImageClick: (ReviewImage) -> Unit = {},
    onLike: () -> Unit = {},
    onComment: () -> Unit = {},
    onShare: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            // Rating Stars
            ReviewRatingDisplay(
                rating = review.rating,
                modifier = Modifier.padding(bottom = TchatSpacing.sm)
            )

            // Review Text
            if (!review.content.text.isNullOrBlank()) {
                Text(
                    text = review.content.text,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Image Gallery
            if (review.content.images.isNotEmpty()) {
                LazyRow(
                    modifier = Modifier.padding(bottom = TchatSpacing.sm),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(review.content.images.take(3)) { image ->
                        TchatCard(
                            variant = TchatCardVariant.Elevated,
                            modifier = Modifier
                                .size(80.dp)
                                .clickable { onImageClick(image) }
                        ) {
                            Box(
                                modifier = Modifier.fillMaxSize(),
                                contentAlignment = Alignment.Center
                            ) {
                                Icon(
                                    Icons.Default.Image,
                                    contentDescription = "Review Image",
                                    modifier = Modifier.size(32.dp),
                                    tint = TchatColors.success
                                )
                            }
                        }
                    }
                    // Show count if more than 3 images
                    if (review.content.images.size > 3) {
                        item {
                            TchatCard(
                                variant = TchatCardVariant.Outlined,
                                modifier = Modifier.size(80.dp)
                            ) {
                                Box(
                                    modifier = Modifier.fillMaxSize(),
                                    contentAlignment = Alignment.Center
                                ) {
                                    Text(
                                        text = "+${review.content.images.size - 3}",
                                        style = MaterialTheme.typography.titleMedium,
                                        fontWeight = FontWeight.Bold,
                                        color = TchatColors.primary
                                    )
                                }
                            }
                        }
                    }
                }
            }

            // Hashtags
            if (review.content.hashtags.isNotEmpty()) {
                ReviewHashtags(
                    hashtags = review.content.hashtags,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Author and Actions
            ReviewAuthorSection(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare
            )
        }
    }
}

/**
 * Video Review Component
 */
@Composable
fun VideoReview(
    review: Review,
    onVideoClick: (ReviewVideo) -> Unit = {},
    onLike: () -> Unit = {},
    onComment: () -> Unit = {},
    onShare: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            // Rating Stars
            ReviewRatingDisplay(
                rating = review.rating,
                modifier = Modifier.padding(bottom = TchatSpacing.sm)
            )

            // Video Gallery - Horizontal scrollable like TikTok
            if (review.content.videos.isNotEmpty()) {
                LazyRow(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(bottom = TchatSpacing.sm),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(review.content.videos) { video ->
                        TchatCard(
                            variant = TchatCardVariant.Elevated,
                            modifier = Modifier
                                .size(width = 120.dp, height = 160.dp)
                                .clickable { onVideoClick(video) }
                        ) {
                            Box(
                                modifier = Modifier.fillMaxSize(),
                                contentAlignment = Alignment.Center
                            ) {
                                // Video preview
                                Box(
                                    modifier = Modifier
                                        .fillMaxSize()
                                        .background(TchatColors.primary.copy(alpha = 0.1f)),
                                    contentAlignment = Alignment.Center
                                ) {
                                    Icon(
                                        Icons.Default.PlayArrow,
                                        contentDescription = "Play Video",
                                        modifier = Modifier.size(48.dp),
                                        tint = TchatColors.primary
                                    )
                                }

                                // Duration badge
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
                                        text = video.duration,
                                        style = MaterialTheme.typography.labelSmall,
                                        color = Color.White,
                                        fontSize = 10.sp
                                    )
                                }
                            }
                        }
                    }
                }
            }

            // Review Text
            if (!review.content.text.isNullOrBlank()) {
                Text(
                    text = review.content.text,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Hashtags
            if (review.content.hashtags.isNotEmpty()) {
                ReviewHashtags(
                    hashtags = review.content.hashtags,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Author and Actions
            ReviewAuthorSection(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare
            )
        }
    }
}

/**
 * Mixed Review Component (Text + Images + Videos)
 */
@Composable
fun MixedReview(
    review: Review,
    onImageClick: (ReviewImage) -> Unit = {},
    onVideoClick: (ReviewVideo) -> Unit = {},
    onLike: () -> Unit = {},
    onComment: () -> Unit = {},
    onShare: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            // Rating Stars
            ReviewRatingDisplay(
                rating = review.rating,
                modifier = Modifier.padding(bottom = TchatSpacing.sm)
            )

            // Review Text
            if (!review.content.text.isNullOrBlank()) {
                Text(
                    text = review.content.text,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Mixed Media Gallery
            LazyRow(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(bottom = TchatSpacing.sm),
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                // Videos first
                items(review.content.videos.take(2)) { video ->
                    TchatCard(
                        variant = TchatCardVariant.Elevated,
                        modifier = Modifier
                            .size(width = 100.dp, height = 140.dp)
                            .clickable { onVideoClick(video) }
                    ) {
                        Box(
                            modifier = Modifier.fillMaxSize(),
                            contentAlignment = Alignment.Center
                        ) {
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .background(TchatColors.primary.copy(alpha = 0.1f)),
                                contentAlignment = Alignment.Center
                            ) {
                                Icon(
                                    Icons.Default.PlayArrow,
                                    contentDescription = "Play Video",
                                    modifier = Modifier.size(32.dp),
                                    tint = TchatColors.primary
                                )
                            }

                            Box(
                                modifier = Modifier
                                    .align(Alignment.BottomEnd)
                                    .padding(4.dp)
                                    .background(
                                        Color.Black.copy(alpha = 0.7f),
                                        RoundedCornerShape(3.dp)
                                    )
                                    .padding(horizontal = 4.dp, vertical = 2.dp)
                            ) {
                                Text(
                                    text = video.duration,
                                    style = MaterialTheme.typography.labelSmall,
                                    color = Color.White,
                                    fontSize = 8.sp
                                )
                            }
                        }
                    }
                }

                // Then images
                items(review.content.images.take(3)) { image ->
                    TchatCard(
                        variant = TchatCardVariant.Elevated,
                        modifier = Modifier
                            .size(100.dp)
                            .clickable { onImageClick(image) }
                    ) {
                        Box(
                            modifier = Modifier.fillMaxSize(),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.Image,
                                contentDescription = "Review Image",
                                modifier = Modifier.size(32.dp),
                                tint = TchatColors.success
                            )
                        }
                    }
                }

                // Show more indicator
                val totalMedia = review.content.videos.size + review.content.images.size
                if (totalMedia > 5) {
                    item {
                        TchatCard(
                            variant = TchatCardVariant.Outlined,
                            modifier = Modifier.size(100.dp)
                        ) {
                            Box(
                                modifier = Modifier.fillMaxSize(),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    text = "+${totalMedia - 5}",
                                    style = MaterialTheme.typography.titleMedium,
                                    fontWeight = FontWeight.Bold,
                                    color = TchatColors.primary
                                )
                            }
                        }
                    }
                }
            }

            // Hashtags
            if (review.content.hashtags.isNotEmpty()) {
                ReviewHashtags(
                    hashtags = review.content.hashtags,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )
            }

            // Author and Actions
            ReviewAuthorSection(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare
            )
        }
    }
}

/**
 * Review Rating Display Component
 */
@Composable
private fun ReviewRatingDisplay(
    rating: ReviewRating,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically
    ) {
        when (rating.type) {
            ReviewRatingType.STARS_5 -> {
                repeat(5) { index ->
                    Icon(
                        if (index < rating.value * 5) Icons.Default.Star else Icons.Default.StarBorder,
                        contentDescription = "Star ${index + 1}",
                        modifier = Modifier.size(16.dp),
                        tint = if (index < rating.value * 5) TchatColors.warning else TchatColors.onSurfaceVariant
                    )
                }
            }
            ReviewRatingType.THUMBS -> {
                Icon(
                    if (rating.value > 0.5f) Icons.Default.ThumbUp else Icons.Default.ThumbDown,
                    contentDescription = "Rating",
                    modifier = Modifier.size(16.dp),
                    tint = if (rating.value > 0.5f) TchatColors.success else TchatColors.error
                )
            }
            else -> {
                Text(
                    text = rating.displayValue,
                    style = MaterialTheme.typography.labelMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.primary
                )
            }
        }

        Spacer(modifier = Modifier.width(TchatSpacing.xs))

        Text(
            text = rating.displayValue,
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

/**
 * Review Hashtags Component
 */
@Composable
private fun ReviewHashtags(
    hashtags: List<String>,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
    ) {
        items(hashtags.take(5)) { hashtag ->
            Surface(
                shape = RoundedCornerShape(12.dp),
                color = TchatColors.primary.copy(alpha = 0.1f)
            ) {
                Text(
                    text = hashtag,
                    modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.primary,
                    fontSize = 10.sp
                )
            }
        }
    }
}

/**
 * Review Author Section Component
 */
@Composable
private fun ReviewAuthorSection(
    review: Review,
    onLike: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Author Avatar and Info
        Row(
            modifier = Modifier.weight(1f),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(24.dp)
                    .clip(CircleShape)
                    .background(TchatColors.primary),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = review.userName.first().toString().uppercase(),
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onPrimary,
                    fontSize = 10.sp,
                    fontWeight = FontWeight.Bold
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.xs))

            Column {
                Text(
                    text = review.userName,
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.onSurface,
                    fontSize = 12.sp
                )
                Text(
                    text = review.createdAt,
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onSurfaceVariant,
                    fontSize = 10.sp
                )
            }

            if (review.isVerifiedPurchase) {
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Surface(
                    shape = RoundedCornerShape(4.dp),
                    color = TchatColors.success.copy(alpha = 0.1f)
                ) {
                    Text(
                        text = "Verified",
                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.success,
                        fontSize = 9.sp
                    )
                }
            }
        }

        // Action Buttons
        Row(
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Like Button
            Row(
                verticalAlignment = Alignment.CenterVertically,
                modifier = Modifier.clickable { onLike() }
            ) {
                Icon(
                    if (review.isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                    contentDescription = "Like",
                    modifier = Modifier.size(16.dp),
                    tint = if (review.isLiked) TchatColors.error else TchatColors.onSurfaceVariant
                )
                if (review.likes > 0) {
                    Spacer(modifier = Modifier.width(2.dp))
                    Text(
                        text = review.likes.toString(),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant,
                        fontSize = 10.sp
                    )
                }
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
                if (review.comments > 0) {
                    Spacer(modifier = Modifier.width(2.dp))
                    Text(
                        text = review.comments.toString(),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant,
                        fontSize = 10.sp
                    )
                }
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
        }
    }
}