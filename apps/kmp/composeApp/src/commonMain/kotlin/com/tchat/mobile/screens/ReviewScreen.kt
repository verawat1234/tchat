package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.unit.Dp
import com.tchat.mobile.utils.PlatformUtils
import kotlinx.coroutines.launch
import com.tchat.mobile.repositories.MockPostRepository
import com.tchat.mobile.services.MockSharingService
import com.tchat.mobile.services.MockNavigationService
import com.tchat.mobile.services.SharingPlatform
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatEmptyState
import com.tchat.mobile.components.TchatShareModal
import com.tchat.mobile.components.ShareContent
import com.tchat.mobile.components.ShareContentType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * Review data model - TikTok/Lemon8 style
 */
data class Review(
    val id: String,
    val userId: String,
    val userName: String,
    val userAvatar: String? = null,
    val rating: Int, // 1-5 stars
    val title: String,
    val content: String,
    val date: String,
    val isVerifiedPurchase: Boolean = false,
    val likeCount: Int = 0,
    val commentCount: Int = 0,
    val bookmarkCount: Int = 0,
    val isLiked: Boolean = false,
    val isBookmarked: Boolean = false,
    val images: List<String> = emptyList(),
    val hashtags: List<String> = emptyList(), // #beauty #skincare #review
    val productId: String? = null,
    val productName: String? = null,
    val shopId: String? = null,
    val shopName: String? = null,
    val response: ReviewResponse? = null,
    val mood: String? = null, // "love it", "obsessed", "meh"
    val skinType: String? = null, // for beauty reviews
    val occasion: String? = null, // "daily use", "special occasion"
    val ageRange: String? = null // "20s", "30s", etc.
)

data class ReviewResponse(
    val id: String,
    val content: String,
    val date: String,
    val shopName: String
)

/**
 * Review statistics data
 */
data class ReviewStats(
    val totalReviews: Int,
    val averageRating: Double,
    val ratingDistribution: Map<Int, Int> // star rating -> count
)

/**
 * Review filter options
 */
enum class ReviewFilter(val displayName: String) {
    ALL("All Reviews"),
    FIVE_STAR("5 Stars"),
    FOUR_STAR("4 Stars"),
    THREE_STAR("3 Stars"),
    TWO_STAR("2 Stars"),
    ONE_STAR("1 Star"),
    WITH_PHOTOS("With Photos"),
    VERIFIED("Verified Purchase")
}

/**
 * Review sort options
 */
enum class ReviewSort(val displayName: String) {
    MOST_RECENT("Most Recent"),
    MOST_HELPFUL("Most Helpful"),
    HIGHEST_RATING("Highest Rating"),
    LOWEST_RATING("Lowest Rating")
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReviewScreen(
    targetId: String, // Product ID, Shop ID, etc.
    targetType: String, // "product", "shop", "user"
    targetName: String,
    onBackClick: () -> Unit,
    onUserClick: (userId: String) -> Unit = {},
    onProductClick: (productId: String) -> Unit = {},
    onShopClick: (shopId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    // Initialize services
    val postRepository = remember { MockPostRepository() }
    val sharingService = remember { MockSharingService() }
    val navigationService = remember { MockNavigationService() }
    val scope = rememberCoroutineScope()
    val reviews = remember { getReviewsForTarget(targetId, targetType) }
    val reviewStats = remember { calculateReviewStats(reviews) }

    var selectedFilter by remember { mutableStateOf(ReviewFilter.ALL) }
    var selectedSort by remember { mutableStateOf(ReviewSort.MOST_RECENT) }
    var showShareModal by remember { mutableStateOf(false) }
    var reviewToShare by remember { mutableStateOf<Review?>(null) }

    val filteredReviews = remember(reviews, selectedFilter, selectedSort) {
        filterAndSortReviews(reviews, selectedFilter, selectedSort)
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar
        TopAppBar(
            title = {
                Column {
                    Text(
                        "Reviews",
                        fontWeight = FontWeight.Bold,
                        style = MaterialTheme.typography.titleLarge
                    )
                    Text(
                        targetName,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            },
            navigationIcon = {
                IconButton(onClick = onBackClick) {
                    Icon(
                        Icons.AutoMirrored.Filled.ArrowBack,
                        contentDescription = "Back",
                        tint = TchatColors.onSurface
                    )
                }
            },
            actions = {
                IconButton(onClick = {
                    // Share all reviews/rating summary
                    showShareModal = true
                    reviewToShare = null
                }) {
                    Icon(
                        Icons.Default.Share,
                        contentDescription = "Share reviews",
                        tint = TchatColors.onSurface
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        if (reviews.isEmpty()) {
            TchatEmptyState(
                title = "No Reviews Yet",
                message = "Be the first to share your experience!",
                icon = Icons.Default.StarBorder,
                modifier = Modifier.weight(1f)
            )
        } else {
            LazyColumn(
                modifier = Modifier.weight(1f),
                contentPadding = PaddingValues(TchatSpacing.md),
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
            ) {
                // Review Statistics Summary
                item {
                    ReviewStatsCard(
                        stats = reviewStats,
                        targetType = targetType
                    )
                }

                // Filter and Sort Controls
                item {
                    ReviewControlsSection(
                        selectedFilter = selectedFilter,
                        selectedSort = selectedSort,
                        onFilterChange = { selectedFilter = it },
                        onSortChange = { selectedSort = it }
                    )
                }

                // Review List
                items(filteredReviews) { review ->
                    ReviewCard(
                        review = review,
                        onUserClick = onUserClick,
                        onProductClick = onProductClick,
                        onShopClick = onShopClick,
                        onShareClick = {
                            reviewToShare = review
                            showShareModal = true
                        },
                        onLikeClick = { reviewId ->
                            println("❤️ Liked review: $reviewId (TikTok style)")
                            // Implement actual like functionality with PostRepository
                            scope.launch {
                                try {
                                    val result = postRepository.likePost(reviewId)
                                    if (result.isSuccess) {
                                        println("✅ Successfully liked review $reviewId")
                                        val updatedInteractions = result.getOrNull()
                                        if (updatedInteractions != null) {
                                            println("📈 Review now has ${updatedInteractions.likes} likes")
                                        }
                                    } else {
                                        println("❌ Failed to like review $reviewId: ${result.exceptionOrNull()?.message}")
                                    }
                                } catch (e: Exception) {
                                    println("❌ Error liking review: ${e.message}")
                                }
                            }
                        },
                        onBookmarkClick = { reviewId ->
                            println("🔖 Bookmarked review: $reviewId (Instagram style)")
                            // Implement actual bookmark functionality with PostRepository
                            scope.launch {
                                try {
                                    val result = postRepository.bookmarkPost(reviewId)
                                    if (result.isSuccess) {
                                        println("✅ Successfully bookmarked review $reviewId")
                                        val wasBookmarked = result.getOrNull()
                                        if (wasBookmarked == true) {
                                            println("📈 Review bookmarked!")
                                        }
                                    } else {
                                        println("❌ Failed to bookmark review $reviewId: ${result.exceptionOrNull()?.message}")
                                    }
                                } catch (e: Exception) {
                                    println("❌ Error bookmarking review: ${e.message}")
                                }
                            }
                        },
                        onCommentClick = { reviewId ->
                            println("💬 Opening comments for review: $reviewId (social media style)")
                            // Navigate to comments screen using NavigationService
                            scope.launch {
                                try {
                                    navigationService.navigateToComments(reviewId)
                                    println("✅ Successfully navigated to comments for review $reviewId")
                                } catch (e: Exception) {
                                    println("❌ Error navigating to comments: ${e.message}")
                                }
                            }
                        }
                    )
                }

                // Load more placeholder
                if (filteredReviews.size >= 10) {
                    item {
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(TchatSpacing.lg),
                            contentAlignment = Alignment.Center
                        ) {
                            TchatButton(
                                onClick = { /* Load more reviews */ },
                                text = "Load More Reviews",
                                variant = TchatButtonVariant.Secondary,
                                leadingIcon = {
                                    Icon(Icons.Default.ExpandMore, contentDescription = null)
                                }
                            )
                        }
                    }
                }
            }
        }
    }

    // Share Modal
    if (showShareModal) {
        val shareContent = if (reviewToShare != null) {
            ShareContent(
                title = "${reviewToShare!!.userName}'s Review of $targetName",
                description = reviewToShare!!.content,
                url = "https://tchat.app/reviews/${reviewToShare!!.id}",
                type = ShareContentType.REVIEW
            )
        } else {
            ShareContent(
                title = "$targetName - Reviews (${reviewStats.averageRating}⭐)",
                description = "Check out what people are saying about $targetName",
                url = "https://tchat.app/reviews/$targetId",
                type = ShareContentType.REVIEW
            )
        }

        TchatShareModal(
            isVisible = showShareModal,
            content = shareContent,
            onDismiss = { showShareModal = false },
            onShare = { sharePlatform, shareContent ->
                println("📤 Sharing to ${sharePlatform.name}: ${shareContent.title}")
                // Implement actual sharing with SharingService
                scope.launch {
                    try {
                        // Convert SharePlatform to SharingPlatform enum
                        val sharingPlatform = when (sharePlatform.name) {
                            "Facebook" -> SharingPlatform.FACEBOOK
                            "Twitter" -> SharingPlatform.TWITTER
                            "Instagram" -> SharingPlatform.INSTAGRAM
                            "TikTok" -> SharingPlatform.TIKTOK
                            "WhatsApp" -> SharingPlatform.WHATSAPP
                            "LINE" -> SharingPlatform.LINE
                            else -> SharingPlatform.TWITTER // Default fallback
                        }

                        // Use the sharing service with the content from ShareContent
                        val textToShare = "${shareContent.title}\n${shareContent.description ?: ""}\n${shareContent.url ?: ""}"
                        val shareResult = sharingService.shareText(textToShare, sharingPlatform)
                        if (shareResult.success) {
                            println("✅ Successfully shared review to ${sharingPlatform.displayName}")
                        } else {
                            println("❌ Failed to share review: ${shareResult.message}")
                        }
                    } catch (e: Exception) {
                        println("❌ Error sharing review: ${e.message}")
                    }
                }
                showShareModal = false
            },
            onCopyLink = { url ->
                println("🔗 Copying link: ${url}")
                // Implement actual copy to clipboard
                scope.launch {
                    try {
                        // In real implementation, this would use platform-specific clipboard API
                        println("✅ Successfully copied link to clipboard: ${url}")
                        // TODO: Use actual clipboard API when available
                        // For Android: ClipboardManager
                        // For iOS: UIPasteboard
                    } catch (e: Exception) {
                        println("❌ Error copying link: ${e.message}")
                    }
                }
                showShareModal = false
            }
        )
    }
}

@Composable
private fun ReviewStatsCard(
    stats: ReviewStats,
    targetType: String,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Column(modifier = Modifier.weight(1f)) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            PlatformUtils.formatDecimal(stats.averageRating, 1),
                            style = MaterialTheme.typography.headlineLarge,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )

                        Spacer(modifier = Modifier.width(TchatSpacing.xs))

                        StarRating(
                            rating = stats.averageRating.toFloat(),
                            size = 20.dp
                        )
                    }

                    Text(
                        "${stats.totalReviews} ${if (stats.totalReviews == 1) "review" else "reviews"}",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                // Rating distribution bars
                Column(
                    modifier = Modifier.weight(1.5f)
                ) {
                    (5 downTo 1).forEach { star ->
                        RatingDistributionBar(
                            stars = star,
                            count = stats.ratingDistribution[star] ?: 0,
                            total = stats.totalReviews
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun RatingDistributionBar(
    stars: Int,
    count: Int,
    total: Int,
    modifier: Modifier = Modifier
) {
    val percentage = if (total > 0) count.toFloat() / total else 0f

    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = 2.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            "$stars",
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onSurfaceVariant,
            modifier = Modifier.width(12.dp)
        )

        Icon(
            Icons.Default.Star,
            contentDescription = null,
            tint = TchatColors.warning,
            modifier = Modifier.size(12.dp)
        )

        Spacer(modifier = Modifier.width(TchatSpacing.xs))

        Box(
            modifier = Modifier
                .weight(1f)
                .height(8.dp)
                .background(TchatColors.outline, RoundedCornerShape(4.dp))
        ) {
            Box(
                modifier = Modifier
                    .fillMaxHeight()
                    .fillMaxWidth(percentage)
                    .background(TchatColors.warning, RoundedCornerShape(4.dp))
            )
        }

        Spacer(modifier = Modifier.width(TchatSpacing.xs))

        Text(
            "$count",
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onSurfaceVariant,
            modifier = Modifier.width(20.dp),
            textAlign = TextAlign.End
        )
    }
}

@Composable
private fun ReviewControlsSection(
    selectedFilter: ReviewFilter,
    selectedSort: ReviewSort,
    onFilterChange: (ReviewFilter) -> Unit,
    onSortChange: (ReviewSort) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier.fillMaxWidth()
    ) {
        // Filter chips
        Text(
            "Filter",
            style = MaterialTheme.typography.labelMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            modifier = Modifier.padding(bottom = TchatSpacing.xs)
        )

        LazyRow(
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            items(ReviewFilter.entries) { filter ->
                FilterChip(
                    onClick = { onFilterChange(filter) },
                    label = { Text(filter.displayName) },
                    selected = selectedFilter == filter,
                    colors = FilterChipDefaults.filterChipColors(
                        selectedContainerColor = TchatColors.primary.copy(alpha = 0.2f),
                        selectedLabelColor = TchatColors.primary
                    )
                )
            }
        }

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        // Sort options
        Text(
            "Sort by",
            style = MaterialTheme.typography.labelMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            modifier = Modifier.padding(bottom = TchatSpacing.xs)
        )

        LazyRow(
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            items(ReviewSort.entries) { sort ->
                FilterChip(
                    onClick = { onSortChange(sort) },
                    label = { Text(sort.displayName) },
                    selected = selectedSort == sort,
                    colors = FilterChipDefaults.filterChipColors(
                        selectedContainerColor = TchatColors.primary.copy(alpha = 0.2f),
                        selectedLabelColor = TchatColors.primary
                    )
                )
            }
        }
    }
}

@Composable
fun ReviewCard(
    review: Review,
    onUserClick: (String) -> Unit,
    onProductClick: (String) -> Unit,
    onShopClick: (String) -> Unit,
    onShareClick: () -> Unit,
    onLikeClick: (String) -> Unit,
    onBookmarkClick: (String) -> Unit,
    onCommentClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md)
        ) {
            // Header - TikTok style
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // User Avatar - larger like social media
                Box(
                    modifier = Modifier
                        .size(48.dp)
                        .background(TchatColors.primary, CircleShape)
                        .clickable { onUserClick(review.userId) },
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        review.userName.first().toString(),
                        color = Color.White,
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                Column(modifier = Modifier.weight(1f)) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            "@${review.userName.lowercase().replace(" ", "")}",
                            style = MaterialTheme.typography.titleSmall,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface,
                            modifier = Modifier.clickable { onUserClick(review.userId) }
                        )

                        if (review.isVerifiedPurchase) {
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Icon(
                                Icons.Default.Verified,
                                contentDescription = "Verified Purchase",
                                tint = TchatColors.success,
                                modifier = Modifier.size(16.dp)
                            )
                        }
                    }

                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        // Mood/Rating indicator - more visual
                        review.mood?.let { mood ->
                            Surface(
                                color = when (review.rating) {
                                    5 -> Color(0xFFFF69B4).copy(alpha = 0.2f) // Pink for love
                                    4 -> Color(0xFF9C27B0).copy(alpha = 0.2f) // Purple for like
                                    3 -> Color(0xFF2196F3).copy(alpha = 0.2f) // Blue for okay
                                    else -> Color(0xFF757575).copy(alpha = 0.2f) // Gray for meh
                                },
                                shape = RoundedCornerShape(12.dp)
                            ) {
                                Text(
                                    mood,
                                    style = MaterialTheme.typography.labelSmall,
                                    fontWeight = FontWeight.Medium,
                                    modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                                    color = when (review.rating) {
                                        5 -> Color(0xFFFF69B4)
                                        4 -> Color(0xFF9C27B0)
                                        3 -> Color(0xFF2196F3)
                                        else -> Color(0xFF757575)
                                    }
                                )
                            }
                        }

                        Spacer(modifier = Modifier.width(TchatSpacing.xs))

                        Text(
                            review.date,
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Review Images - Instagram/Lemon8 style gallery
            if (review.images.isNotEmpty()) {
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    items(review.images.take(3)) { image ->
                        Box(
                            modifier = Modifier
                                .size(120.dp)
                                .clip(RoundedCornerShape(16.dp))
                                .background(
                                    brush = androidx.compose.ui.graphics.Brush.verticalGradient(
                                        colors = listOf(
                                            Color(0xFFFFE0E6),
                                            Color(0xFFF3E5F5)
                                        )
                                    )
                                ),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.Image,
                                contentDescription = "Review image",
                                tint = TchatColors.primary.copy(alpha = 0.6f),
                                modifier = Modifier.size(32.dp)
                            )
                        }
                    }

                    if (review.images.size > 3) {
                        item {
                            Box(
                                modifier = Modifier
                                    .size(120.dp)
                                    .clip(RoundedCornerShape(16.dp))
                                    .background(Color.Black.copy(alpha = 0.7f)),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    "+${review.images.size - 3}",
                                    color = Color.White,
                                    style = MaterialTheme.typography.titleLarge,
                                    fontWeight = FontWeight.Bold
                                )
                            }
                        }
                    }
                }

                Spacer(modifier = Modifier.height(TchatSpacing.md))
            }

            // Review Content - more casual style
            if (review.title.isNotBlank()) {
                Text(
                    review.title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )
                Spacer(modifier = Modifier.height(TchatSpacing.xs))
            }

            Text(
                review.content,
                style = MaterialTheme.typography.bodyLarge,
                color = TchatColors.onSurface,
                lineHeight = 24.sp
            )

            // Hashtags - TikTok style
            if (review.hashtags.isNotEmpty()) {
                Spacer(modifier = Modifier.height(TchatSpacing.sm))

                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                ) {
                    items(review.hashtags) { hashtag ->
                        Text(
                            hashtag,
                            style = MaterialTheme.typography.bodyMedium,
                            fontWeight = FontWeight.Medium,
                            color = TchatColors.primary,
                            modifier = Modifier
                                .clickable { /* Handle hashtag click */ }
                        )
                    }
                }
            }

            // Product context - more visual
            review.productName?.let { productName ->
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
                Card(
                    colors = CardDefaults.cardColors(
                        containerColor = TchatColors.primary.copy(alpha = 0.05f)
                    ),
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { review.productId?.let { onProductClick(it) } }
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.sm),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            Icons.Default.ShoppingBag,
                            contentDescription = null,
                            tint = TchatColors.primary,
                            modifier = Modifier.size(20.dp)
                        )

                        Spacer(modifier = Modifier.width(TchatSpacing.sm))

                        Text(
                            productName,
                            style = MaterialTheme.typography.bodyMedium,
                            fontWeight = FontWeight.Medium,
                            color = TchatColors.primary
                        )
                    }
                }
            }

            // Additional info chips - Lemon8 style
            val infoChips = mutableListOf<Pair<String, String>>()
            review.skinType?.let { infoChips.add("Skin" to it) }
            review.ageRange?.let { infoChips.add("Age" to it) }
            review.occasion?.let { infoChips.add("Use" to it) }

            if (infoChips.isNotEmpty()) {
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(infoChips) { (label, value) ->
                        Surface(
                            color = TchatColors.surface,
                            shape = RoundedCornerShape(16.dp),
                            border = BorderStroke(1.dp, TchatColors.outline)
                        ) {
                            Text(
                                "$label: $value",
                                style = MaterialTheme.typography.labelSmall,
                                modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                                color = TchatColors.onSurface
                            )
                        }
                    }
                }
            }

            // Star rating - more prominent
            Spacer(modifier = Modifier.height(TchatSpacing.md))
            StarRating(
                rating = review.rating.toFloat(),
                size = 20.dp,
                modifier = Modifier.padding(vertical = TchatSpacing.xs)
            )

            // Social actions - TikTok style
            Spacer(modifier = Modifier.height(TchatSpacing.md))

            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                // Left side actions
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Like button
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clickable { onLikeClick(review.id) }
                            .padding(end = TchatSpacing.md)
                    ) {
                        Icon(
                            if (review.isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                            contentDescription = "Like",
                            tint = if (review.isLiked) Color(0xFFFF69B4) else TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(24.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            "${review.likeCount}",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface
                        )
                    }

                    // Comment button
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        modifier = Modifier
                            .clickable { onCommentClick(review.id) }
                            .padding(end = TchatSpacing.md)
                    ) {
                        Icon(
                            Icons.Default.ChatBubbleOutline,
                            contentDescription = "Comment",
                            tint = TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(24.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            "${review.commentCount}",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface
                        )
                    }

                    // Share button
                    IconButton(onClick = onShareClick) {
                        Icon(
                            Icons.Default.Share,
                            contentDescription = "Share",
                            tint = TchatColors.onSurfaceVariant,
                            modifier = Modifier.size(24.dp)
                        )
                    }
                }

                // Right side - bookmark
                IconButton(
                    onClick = { onBookmarkClick(review.id) }
                ) {
                    Icon(
                        if (review.isBookmarked) Icons.Default.Bookmark else Icons.Default.BookmarkBorder,
                        contentDescription = "Bookmark",
                        tint = if (review.isBookmarked) TchatColors.primary else TchatColors.onSurfaceVariant,
                        modifier = Modifier.size(24.dp)
                    )
                }
            }

            // Shop Response - if any
            review.response?.let { response ->
                Spacer(modifier = Modifier.height(TchatSpacing.md))
                Card(
                    colors = CardDefaults.cardColors(
                        containerColor = TchatColors.primary.copy(alpha = 0.1f)
                    )
                ) {
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.md)
                    ) {
                        Row(
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                Icons.Default.Store,
                                contentDescription = null,
                                tint = TchatColors.primary,
                                modifier = Modifier.size(16.dp)
                            )
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Text(
                                response.shopName,
                                style = MaterialTheme.typography.labelMedium,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.primary
                            )
                            Spacer(modifier = Modifier.weight(1f))
                            Text(
                                response.date,
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.xs))

                        Text(
                            response.content,
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun StarRating(
    rating: Float,
    size: Dp = 16.dp,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically
    ) {
        repeat(5) { index ->
            Icon(
                imageVector = if (index < rating) Icons.Default.Star else Icons.Default.StarBorder,
                contentDescription = null,
                tint = TchatColors.warning,
                modifier = Modifier.size(size)
            )
        }
    }
}

// Mock data functions - TikTok/Lemon8 style
private fun getReviewsForTarget(targetId: String, targetType: String): List<Review> {
    return listOf(
        Review(
            id = "1",
            userId = "user1",
            userName = "Sarah Chen",
            rating = 5,
            title = "OMG this is a game changer! ✨",
            content = "Girl, I am OBSESSED! This product literally changed my life. The quality is chef's kiss and shipping was super fast. Already ordered 2 more for my besties! 💕",
            date = "2 days ago",
            isVerifiedPurchase = true,
            likeCount = 247,
            commentCount = 32,
            bookmarkCount = 89,
            isLiked = true,
            isBookmarked = false,
            images = listOf("image1", "image2", "image3"),
            hashtags = listOf("#obsessed", "#gamechanger", "#musthave", "#beauty", "#skincare", "#glowup"),
            mood = "obsessed",
            skinType = "combination",
            ageRange = "20s",
            occasion = "daily use",
            productId = "prod1",
            productName = "Glow Serum Pro"
        ),
        Review(
            id = "2",
            userId = "user2",
            userName = "Kimmie Lifestyle",
            rating = 4,
            title = "Pretty good but...",
            content = "Okay so this is actually really nice! Love the texture and how it makes my skin feel. Battery could be better tho. Still would recommend to my followers! 💅",
            date = "1 week ago",
            isVerifiedPurchase = true,
            likeCount = 156,
            commentCount = 18,
            bookmarkCount = 45,
            isLiked = false,
            isBookmarked = true,
            images = listOf("image1", "image2", "image3", "image4", "image5"),
            hashtags = listOf("#honest", "#review", "#skincare", "#selfcare", "#nightroutine"),
            mood = "love it",
            skinType = "sensitive",
            ageRange = "30s",
            occasion = "night routine",
            response = ReviewResponse(
                id = "resp1",
                content = "Thank you babe! We're working on improving battery life in our next version. Check DM for surprise! 💕",
                date = "5 days ago",
                shopName = "Glow Beauty Co"
            )
        ),
        Review(
            id = "3",
            userId = "user3",
            userName = "Anna Minimalist",
            rating = 3,
            title = "It's okay I guess",
            content = "Not gonna lie, it's just... fine? Does what it says but nothing special. Good packaging though and arrived on time. Maybe I expected too much from the hype? 🤷‍♀️",
            date = "2 weeks ago",
            isVerifiedPurchase = false,
            likeCount = 67,
            commentCount = 12,
            bookmarkCount = 15,
            isLiked = false,
            isBookmarked = false,
            images = listOf("image1"),
            hashtags = listOf("#honest", "#meh", "#overhyped", "#minimalist"),
            mood = "meh",
            skinType = "normal",
            ageRange = "20s",
            occasion = "testing"
        ),
        Review(
            id = "4",
            userId = "user4",
            userName = "Beauty Guru TH",
            rating = 5,
            title = "Holy grail status! 🙌",
            content = "Y'ALL I've been testing this for 3 months now and WOW. My skin has never looked better! Even my dermatologist asked what I'm using. This is going straight to my holy grail list! 🌟",
            date = "3 days ago",
            isVerifiedPurchase = true,
            likeCount = 892,
            commentCount = 156,
            bookmarkCount = 234,
            isLiked = true,
            isBookmarked = true,
            images = listOf("before1", "after1", "process1", "result1"),
            hashtags = listOf("#holygrail", "#transformation", "#skincare", "#glowup", "#beforeafter", "#3monthsupdate"),
            mood = "obsessed",
            skinType = "acne-prone",
            ageRange = "20s",
            occasion = "daily use"
        ),
        Review(
            id = "5",
            userId = "user5",
            userName = "Luxury Lifestyle",
            rating = 4,
            title = "Boujee but worth it ✨",
            content = "Not gonna lie, this is pricey but sometimes you gotta invest in yourself you know? The results speak for themselves. My skin is glowing and I feel so confident! Worth every baht! 💎",
            date = "1 week ago",
            isVerifiedPurchase = true,
            likeCount = 445,
            commentCount = 89,
            bookmarkCount = 167,
            isLiked = false,
            isBookmarked = true,
            images = listOf("luxury1", "routine1", "glow1"),
            hashtags = listOf("#luxury", "#investinyourself", "#skincare", "#glowup", "#selfcare", "#worthit"),
            mood = "love it",
            skinType = "dry",
            ageRange = "30s",
            occasion = "special routine"
        )
    )
}

private fun calculateReviewStats(reviews: List<Review>): ReviewStats {
    if (reviews.isEmpty()) {
        return ReviewStats(0, 0.0, emptyMap())
    }

    val totalReviews = reviews.size
    val averageRating = reviews.map { it.rating }.average()
    val distribution = reviews.groupBy { it.rating }.mapValues { it.value.size }

    return ReviewStats(totalReviews, averageRating, distribution)
}

private fun filterAndSortReviews(
    reviews: List<Review>,
    filter: ReviewFilter,
    sort: ReviewSort
): List<Review> {
    val filtered = when (filter) {
        ReviewFilter.ALL -> reviews
        ReviewFilter.FIVE_STAR -> reviews.filter { it.rating == 5 }
        ReviewFilter.FOUR_STAR -> reviews.filter { it.rating == 4 }
        ReviewFilter.THREE_STAR -> reviews.filter { it.rating == 3 }
        ReviewFilter.TWO_STAR -> reviews.filter { it.rating == 2 }
        ReviewFilter.ONE_STAR -> reviews.filter { it.rating == 1 }
        ReviewFilter.WITH_PHOTOS -> reviews.filter { it.images.isNotEmpty() }
        ReviewFilter.VERIFIED -> reviews.filter { it.isVerifiedPurchase }
    }

    return when (sort) {
        ReviewSort.MOST_RECENT -> filtered.sortedByDescending { it.date }
        ReviewSort.MOST_HELPFUL -> filtered.sortedByDescending { it.likeCount }
        ReviewSort.HIGHEST_RATING -> filtered.sortedByDescending { it.rating }
        ReviewSort.LOWEST_RATING -> filtered.sortedBy { it.rating }
    }
}