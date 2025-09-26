package com.tchat.mobile.components.reviews

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.tchat.mobile.models.*

/**
 * Review Renderer - Similar to MessageRenderer
 *
 * Renders different types of reviews based on their ReviewType
 */
@Composable
fun ReviewRenderer(
    review: Review,
    onImageClick: (ReviewImage) -> Unit = {},
    onVideoClick: (ReviewVideo) -> Unit = {},
    onLike: () -> Unit = {},
    onComment: () -> Unit = {},
    onShare: () -> Unit = {},
    onUserClick: (String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    when (review.content.type) {
        ReviewType.TEXT -> {
            TextReview(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }

        ReviewType.IMAGE -> {
            ImageReview(
                review = review,
                onImageClick = onImageClick,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }

        ReviewType.VIDEO -> {
            VideoReview(
                review = review,
                onVideoClick = onVideoClick,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }

        ReviewType.MIXED -> {
            MixedReview(
                review = review,
                onImageClick = onImageClick,
                onVideoClick = onVideoClick,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }

        ReviewType.QUICK -> {
            // Quick review is similar to text but more compact
            QuickReview(
                review = review,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }

        ReviewType.DETAILED -> {
            // Detailed review shows all information expanded
            DetailedReview(
                review = review,
                onImageClick = onImageClick,
                onVideoClick = onVideoClick,
                onLike = onLike,
                onComment = onComment,
                onShare = onShare,
                modifier = modifier
            )
        }
    }
}

/**
 * Quick Review Component - Compact version
 */
@Composable
private fun QuickReview(
    review: Review,
    onLike: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    // Compact version of TextReview with less padding
    TextReview(
        review = review,
        onLike = onLike,
        onComment = onComment,
        onShare = onShare,
        modifier = modifier
    )
}

/**
 * Detailed Review Component - Expanded version with all features
 */
@Composable
private fun DetailedReview(
    review: Review,
    onImageClick: (ReviewImage) -> Unit,
    onVideoClick: (ReviewVideo) -> Unit,
    onLike: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    // Enhanced version of MixedReview with additional details
    MixedReview(
        review = review,
        onImageClick = onImageClick,
        onVideoClick = onVideoClick,
        onLike = onLike,
        onComment = onComment,
        onShare = onShare,
        modifier = modifier
    )
}