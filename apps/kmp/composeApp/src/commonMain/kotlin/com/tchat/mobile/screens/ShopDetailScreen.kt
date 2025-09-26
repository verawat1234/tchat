package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.components.TchatNotFoundState
import com.tchat.mobile.components.TchatShareModal
import com.tchat.mobile.components.ShareContent
import com.tchat.mobile.components.ShareContentType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.components.reviews.ReviewRenderer
import com.tchat.mobile.models.Review as ModelReview
import com.tchat.mobile.models.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ShopDetailScreen(
    shopId: String,
    onBackClick: () -> Unit,
    onProductClick: (productId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    val shop = getShopById(shopId)
    val shopProducts = getProductsByShopId(shopId)
    var selectedTab by remember { mutableStateOf(0) }
    var searchQuery by remember { mutableStateOf("") }
    var isFollowing by remember { mutableStateOf(false) }
    var isFavorited by remember { mutableStateOf(false) }
    var showShareSheet by remember { mutableStateOf(false) }

    // Review interaction states
    var selectedImageUrl by remember { mutableStateOf<String?>(null) }
    var selectedVideoUrl by remember { mutableStateOf<String?>(null) }
    var showImageViewer by remember { mutableStateOf(false) }
    var showVideoPlayer by remember { mutableStateOf(false) }
    var showCommentsSheet by remember { mutableStateOf(false) }
    var currentReviewId by remember { mutableStateOf<String?>(null) }

    val tabs = listOf("Products", "Reviews", "About")

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar
        TopAppBar(
            title = { Text(shop?.name ?: "Shop", fontWeight = FontWeight.Bold) },
            navigationIcon = {
                IconButton(onClick = onBackClick) {
                    Icon(
                        Icons.Default.ArrowBack,
                        contentDescription = "Back",
                        tint = TchatColors.onSurface
                    )
                }
            },
            actions = {
                IconButton(onClick = {
                    showShareSheet = true
                    // TODO: Implement platform-specific sharing
                }) {
                    Icon(
                        Icons.Default.Share,
                        contentDescription = "Share shop",
                        tint = TchatColors.onSurface
                    )
                }
                IconButton(onClick = { isFavorited = !isFavorited }) {
                    Icon(
                        if (isFavorited) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                        contentDescription = if (isFavorited) "Remove from favorites" else "Add to favorites",
                        tint = if (isFavorited) TchatColors.primary else TchatColors.onSurface
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        if (shop != null) {
            LazyColumn(
                modifier = Modifier.weight(1f)
            ) {
                item {
                    // Shop Header
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.md)
                    ) {
                        // Cover Image
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(200.dp)
                                .clip(RoundedCornerShape(16.dp))
                                .background(TchatColors.primary.copy(alpha = 0.1f)),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.Store,
                                contentDescription = "Shop Cover",
                                modifier = Modifier.size(80.dp),
                                tint = TchatColors.primary
                            )
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Shop Info Row
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            // Shop Avatar
                            Box(
                                modifier = Modifier
                                    .size(80.dp)
                                    .clip(CircleShape)
                                    .background(TchatColors.primaryLight),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    shop.name.first().toString(),
                                    color = TchatColors.onPrimary,
                                    style = MaterialTheme.typography.headlineMedium,
                                    fontWeight = FontWeight.Bold
                                )
                            }

                            Spacer(modifier = Modifier.width(TchatSpacing.md))

                            Column(modifier = Modifier.weight(1f)) {
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    Text(
                                        shop.name,
                                        style = MaterialTheme.typography.headlineSmall,
                                        fontWeight = FontWeight.Bold,
                                        color = TchatColors.onSurface
                                    )
                                    if (shop.isVerified) {
                                        Spacer(modifier = Modifier.width(TchatSpacing.xs))
                                        Icon(
                                            Icons.Default.Verified,
                                            contentDescription = "Verified",
                                            modifier = Modifier.size(20.dp),
                                            tint = TchatColors.primary
                                        )
                                    }
                                }

                                Text(
                                    shop.description,
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = TchatColors.onSurfaceVariant,
                                    maxLines = 2
                                )

                                Spacer(modifier = Modifier.height(TchatSpacing.xs))

                                // Stats Row
                                Row(
                                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                                ) {
                                    StatItem(
                                        icon = Icons.Default.Star,
                                        value = "${shop.rating}",
                                        label = "Rating",
                                        color = TchatColors.warning
                                    )
                                    StatItem(
                                        icon = Icons.Default.People,
                                        value = "${shop.followers}",
                                        label = "Followers",
                                        color = TchatColors.primary
                                    )
                                    StatItem(
                                        icon = Icons.Default.Inventory,
                                        value = "${shop.totalProducts}",
                                        label = "Products",
                                        color = TchatColors.success
                                    )
                                }
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Action Buttons
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            TchatButton(
                                onClick = { isFollowing = !isFollowing },
                                text = if (isFollowing) "Following" else "Follow",
                                variant = if (isFollowing) TchatButtonVariant.Secondary else TchatButtonVariant.Primary,
                                modifier = Modifier.weight(1f),
                                leadingIcon = {
                                    Icon(
                                        if (isFollowing) Icons.Default.Check else Icons.Default.Add,
                                        contentDescription = null
                                    )
                                }
                            )
                            TchatButton(
                                onClick = { /* Message shop functionality */ },
                                text = "Message",
                                variant = TchatButtonVariant.Secondary,
                                modifier = Modifier.weight(1f),
                                leadingIcon = {
                                    Icon(Icons.Default.Message, contentDescription = null)
                                }
                            )
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Shop Info Cards
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            InfoCard(
                                icon = Icons.Default.Schedule,
                                title = "Delivery",
                                subtitle = shop.deliveryTime,
                                modifier = Modifier.weight(1f)
                            )
                            InfoCard(
                                icon = Icons.Default.LocationOn,
                                title = "Distance",
                                subtitle = shop.distance,
                                modifier = Modifier.weight(1f)
                            )
                        }
                    }
                }

                item {
                    // Search Bar
                    TchatInput(
                        value = searchQuery,
                        onValueChange = { searchQuery = it },
                        type = TchatInputType.Search,
                        placeholder = "Search products in this shop...",
                        leadingIcon = Icons.Default.Search,
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(horizontal = TchatSpacing.md)
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.sm))
                }

                item {
                    // Tab Row
                    TabRow(
                        selectedTabIndex = selectedTab,
                        containerColor = TchatColors.surface,
                        contentColor = TchatColors.primary,
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        tabs.forEachIndexed { index, title ->
                            Tab(
                                selected = selectedTab == index,
                                onClick = { selectedTab = index },
                                text = {
                                    Text(
                                        title,
                                        fontWeight = if (selectedTab == index) FontWeight.Bold else FontWeight.Medium
                                    )
                                }
                            )
                        }
                    }

                    Spacer(modifier = Modifier.height(TchatSpacing.md))
                }

                // Tab Content
                when (selectedTab) {
                    0 -> { // Products
                        val filteredProducts = if (searchQuery.isBlank()) {
                            shopProducts
                        } else {
                            shopProducts.filter { product ->
                                product.name.contains(searchQuery, ignoreCase = true) ||
                                product.category.contains(searchQuery, ignoreCase = true)
                            }
                        }

                        if (filteredProducts.isEmpty() && searchQuery.isNotBlank()) {
                            item {
                                Column(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .padding(TchatSpacing.lg),
                                    horizontalAlignment = Alignment.CenterHorizontally
                                ) {
                                    Icon(
                                        Icons.Default.SearchOff,
                                        contentDescription = "No results",
                                        modifier = Modifier.size(48.dp),
                                        tint = TchatColors.onSurfaceVariant
                                    )
                                    Spacer(modifier = Modifier.height(TchatSpacing.sm))
                                    Text(
                                        "No products found for \"$searchQuery\"",
                                        style = MaterialTheme.typography.bodyLarge,
                                        color = TchatColors.onSurfaceVariant,
                                        textAlign = TextAlign.Center
                                    )
                                    Text(
                                        "Try searching with different keywords",
                                        style = MaterialTheme.typography.bodyMedium,
                                        color = TchatColors.onSurfaceVariant.copy(alpha = 0.7f),
                                        textAlign = TextAlign.Center
                                    )
                                }
                            }
                        } else {
                            items(filteredProducts.chunked(2)) { productPair ->
                                Row(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .padding(horizontal = TchatSpacing.md),
                                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                                ) {
                                    productPair.forEach { product ->
                                        ShopProductCard(
                                            product = product,
                                            onClick = { onProductClick(product.id) },
                                            modifier = Modifier.weight(1f)
                                        )
                                    }
                                    // Fill remaining space if odd number
                                    if (productPair.size == 1) {
                                        Spacer(modifier = Modifier.weight(1f))
                                    }
                                }
                                Spacer(modifier = Modifier.height(TchatSpacing.sm))
                            }
                        }
                    }
                    1 -> { // Reviews
                        items(getShopReviewsWithTypes(shopId)) { review ->
                            ReviewRenderer(
                                review = review,
                                onImageClick = { reviewImage ->
                                    selectedImageUrl = reviewImage.url
                                    showImageViewer = true
                                },
                                onVideoClick = { reviewVideo ->
                                    selectedVideoUrl = reviewVideo.url
                                    showVideoPlayer = true
                                },
                                onLike = {
                                    // Toggle like state (this would call API in real app)
                                    println("Liked review ${review.id}")
                                },
                                onComment = {
                                    currentReviewId = review.id
                                    showCommentsSheet = true
                                },
                                onShare = {
                                    // Share this specific review
                                    showShareSheet = true
                                    println("Share review ${review.id}")
                                },
                                onUserClick = { userId ->
                                    // Navigate to user profile (this would use navigation in real app)
                                    println("Navigate to user profile: $userId")
                                },
                                modifier = Modifier.padding(horizontal = TchatSpacing.md, vertical = TchatSpacing.xs)
                            )
                        }
                    }
                    2 -> { // About
                        item {
                            AboutShopContent(shop = shop)
                        }
                    }
                }
            }
        } else {
            // Shop not found
            TchatNotFoundState(
                itemType = "Shop"
            )
        }

        // Share Modal
        if (showShareSheet && shop != null) {
            TchatShareModal(
                isVisible = showShareSheet,
                content = ShareContent(
                    title = shop.name,
                    description = "${shop.description} • ${shop.rating}⭐ rating • ${shop.followers} followers",
                    url = "https://tchat.app/shops/${shop.id}",
                    type = ShareContentType.SHOP
                ),
                onDismiss = { showShareSheet = false },
                onShare = { platform, content ->
                    // Handle platform-specific sharing
                    showShareSheet = false
                },
                onCopyLink = { url ->
                    // Handle copy link functionality
                }
            )
        }
    }
}

@Composable
private fun StatItem(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    value: String,
    label: String,
    color: androidx.compose.ui.graphics.Color,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        Icon(
            icon,
            contentDescription = null,
            modifier = Modifier.size(16.dp),
            tint = color
        )
        Text(
            value,
            style = MaterialTheme.typography.bodyMedium,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )
        Text(
            label,
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun InfoCard(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    title: String,
    subtitle: String,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier,
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                icon,
                contentDescription = null,
                modifier = Modifier.size(24.dp),
                tint = TchatColors.primary
            )
            Spacer(modifier = Modifier.width(TchatSpacing.sm))
            Column {
                Text(
                    title,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
                Text(
                    subtitle,
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
            }
        }
    }
}

@Composable
private fun ShopProductCard(
    product: ProductItem,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onClick,
        modifier = modifier,
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            // Product Image with badges
            Box {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(100.dp)
                        .background(TchatColors.primary.copy(alpha = 0.1f)),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        Icons.Default.ShoppingBag,
                        contentDescription = "Product Image",
                        modifier = Modifier.size(40.dp),
                        tint = TchatColors.primary
                    )
                }

                if (product.isHot) {
                    Badge(
                        containerColor = TchatColors.error,
                        contentColor = TchatColors.onPrimary,
                        modifier = Modifier
                            .align(Alignment.TopStart)
                            .padding(TchatSpacing.xs)
                    ) {
                        Text("HOT", fontSize = 8.sp, fontWeight = FontWeight.Bold)
                    }
                }
            }

            // Product Info
            Column(
                modifier = Modifier.padding(TchatSpacing.sm)
            ) {
                Text(
                    text = product.name,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2
                )

                Spacer(modifier = Modifier.height(4.dp))

                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "$${product.price}",
                        style = MaterialTheme.typography.bodyLarge,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.primary
                    )
                    product.originalPrice?.let { originalPrice ->
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = "$${originalPrice}",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant,
                            textDecoration = TextDecoration.LineThrough
                        )
                    }
                }

                if (product.orders > 0) {
                    Text(
                        text = "${product.orders} sold",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

@Composable
internal fun ReviewCard(
    review: ShopReview,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = TchatSpacing.md, vertical = TchatSpacing.xs),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primaryLight),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        review.userName.first().toString(),
                        color = TchatColors.onPrimary,
                        fontWeight = FontWeight.Bold
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        review.userName,
                        style = MaterialTheme.typography.bodyLarge,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.onSurface
                    )
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        repeat(5) { index ->
                            Icon(
                                if (index < review.rating) Icons.Default.Star else Icons.Default.StarBorder,
                                contentDescription = null,
                                modifier = Modifier.size(16.dp),
                                tint = TchatColors.warning
                            )
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

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            Text(
                review.comment,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )
        }
    }
}

@Composable
private fun AboutShopContent(
    shop: ShopItem,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier.padding(TchatSpacing.md)
    ) {
        Text(
            "About ${shop.name}",
            style = MaterialTheme.typography.titleLarge,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        Text(
            shop.description,
            style = MaterialTheme.typography.bodyLarge,
            color = TchatColors.onSurface
        )

        Spacer(modifier = Modifier.height(TchatSpacing.lg))

        Text(
            "Shop Information",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )

        Spacer(modifier = Modifier.height(TchatSpacing.sm))

        InfoRow("Delivery Time", shop.deliveryTime)
        InfoRow("Distance", shop.distance)
        InfoRow("Total Products", "${shop.totalProducts} items")
        InfoRow("Followers", "${shop.followers} followers")
        InfoRow("Verified", if (shop.isVerified) "Yes" else "No")
    }
}

@Composable
private fun InfoRow(
    label: String,
    value: String,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = TchatSpacing.xs),
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(
            label,
            style = MaterialTheme.typography.bodyMedium,
            color = TchatColors.onSurfaceVariant
        )
        Text(
            value,
            style = MaterialTheme.typography.bodyMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface
        )
    }
}

// Data classes and helper functions
data class ShopReview(
    val id: String,
    val userName: String,
    val rating: Int,
    val comment: String,
    val date: String
)

private fun getShopById(shopId: String): ShopItem? {
    return getDummyShops().find { it.id == shopId }
}

private fun getProductsByShopId(shopId: String): List<ProductItem> {
    return getDummyProducts().filter { product ->
        // Map shop IDs to merchant names
        when (shopId) {
            "1" -> product.merchant == "Bangkok Street Food"
            "2" -> product.merchant == "Tech Paradise"
            "3" -> product.merchant == "Fashion House"
            "4" -> product.merchant == "Fresh Market"
            "5" -> product.merchant == "Coffee Corner"
            else -> false
        }
    }
}

private fun getDummyShops(): List<ShopItem> = listOf(
    ShopItem("1", "Bangkok Street Food", "Authentic Thai street food and snacks delivered fresh to your door", "BS", "cover1", 4.8, "15-25 min", "1.2 km", true, 2540, 45),
    ShopItem("2", "Tech Paradise", "Latest electronics and gadgets with warranty and fast delivery", "TP", "cover2", 4.6, "30-45 min", "3.1 km", true, 1890, 120),
    ShopItem("3", "Fashion House", "Trendy clothes and accessories for modern lifestyle", "FH", "cover3", 4.5, "20-30 min", "2.8 km", false, 980, 78),
    ShopItem("4", "Fresh Market", "Organic fruits and vegetables sourced from local farms", "FM", "cover4", 4.7, "10-20 min", "0.8 km", true, 3200, 230),
    ShopItem("5", "Coffee Corner", "Premium coffee and desserts made with love", "CC", "cover5", 4.9, "5-15 min", "0.5 km", true, 1560, 35)
)

private fun getDummyProducts(): List<ProductItem> = listOf(
    ProductItem("1", "Pad Thai Goong", 45.0, 55.0, 4.8, "Food", "Bangkok Street Food", "image1", "15 min", "1.2 km", true, 18, 245),
    ProductItem("2", "Wireless Headphones", 2490.0, 2990.0, 4.6, "Electronics", "Tech Paradise", "image2", "30 min", "3.1 km", false, 15, 89),
    ProductItem("3", "Cotton T-Shirt", 590.0, null, 4.3, "Fashion", "Fashion House", "image3", "25 min", "2.8 km", false, null, 67),
    ProductItem("4", "Fresh Mango", 120.0, 150.0, 4.9, "Food", "Fresh Market", "image4", "15 min", "0.8 km", true, 20, 156),
    ProductItem("5", "Iced Coffee", 85.0, null, 4.7, "Beverage", "Coffee Corner", "image5", "10 min", "0.5 km", false, null, 98),
    ProductItem("6", "Som Tam", 35.0, null, 4.5, "Food", "Bangkok Street Food", "image6", "15 min", "1.2 km", true, null, 189),
    ProductItem("7", "Smart Watch", 8990.0, 10990.0, 4.4, "Electronics", "Tech Paradise", "image7", "30 min", "3.1 km", false, 18, 45),
    ProductItem("8", "Denim Jeans", 1290.0, 1590.0, 4.2, "Fashion", "Fashion House", "image8", "25 min", "2.8 km", false, 19, 78)
)

private fun getDummyReviews(): List<ShopReview> = listOf(
    ShopReview("1", "Sarah Chen", 5, "Amazing food quality and fast delivery! The Pad Thai was authentic and delicious.", "2 days ago"),
    ShopReview("2", "Mike Johnson", 4, "Great selection of products. Quick delivery and everything was well packaged.", "1 week ago"),
    ShopReview("3", "Anna Wong", 5, "Excellent service! The shop owner is very friendly and helpful.", "2 weeks ago"),
    ShopReview("4", "David Kim", 4, "Good quality products at reasonable prices. Will order again.", "3 weeks ago"),
    ShopReview("5", "Lisa Taylor", 5, "Best shop in the area! Highly recommended for authentic Thai food.", "1 month ago")
)

private fun getShopReviewsWithTypes(shopId: String): List<ModelReview> = listOf(
    // TEXT Review for shop
    ModelReview(
        id = "1",
        userId = "u1",
        userName = "FoodieQueen 👑",
        userAvatar = null,
        targetType = ReviewTargetType.SHOP,
        targetId = shopId,
        targetName = "Bangkok Street Food",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5"
        ),
        content = ReviewContent(
            type = ReviewType.TEXT,
            text = "I'm literally obsessed with this place! The Pad Thai is absolutely divine and tastes just like what I had in Bangkok. The delivery is always lightning fast and the portions are generous. This is my new go-to spot! 🔥",
            hashtags = listOf("#obsessed", "#authentic", "#padthai", "#bangkok", "#fastdelivery", "#generous")
        ),
        isVerifiedPurchase = true,
        likes = 456,
        comments = 89,
        shares = 23,
        isLiked = true,
        isBookmarked = true,
        createdAt = "2 days ago"
    ),

    // IMAGE Review with food photos
    ModelReview(
        id = "2",
        userId = "u2",
        userName = "Mike J 🍜",
        userAvatar = null,
        targetType = ReviewTargetType.SHOP,
        targetId = shopId,
        targetName = "Bangkok Street Food",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 0.8f,
            displayValue = "4/5"
        ),
        content = ReviewContent(
            type = ReviewType.IMAGE,
            text = "Really impressed with the variety they offer! The packaging was top-notch and kept everything fresh. Check out these amazing dishes I ordered! 📸✨",
            images = listOf(
                ReviewImage("img1", "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800", "Amazing variety of dishes", 1.0f),
                ReviewImage("img2", "https://images.unsplash.com/photo-1586190848861-99aa4a171e90?w=800", "Top-notch packaging", 1.2f),
                ReviewImage("img3", "https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=800", "Everything kept fresh", 0.9f)
            ),
            hashtags = listOf("#variety", "#quick", "#packaging", "#fresh", "#impressed", "#foodphotos")
        ),
        isVerifiedPurchase = true,
        likes = 234,
        comments = 56,
        shares = 18,
        isLiked = false,
        isBookmarked = true,
        createdAt = "1 week ago"
    ),

    // VIDEO Review
    ModelReview(
        id = "3",
        userId = "u3",
        userName = "Anna W 🌟",
        userAvatar = null,
        targetType = ReviewTargetType.SHOP,
        targetId = shopId,
        targetName = "Bangkok Street Food",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5"
        ),
        content = ReviewContent(
            type = ReviewType.VIDEO,
            text = "Not only is the food incredible, but the customer service is amazing! Watch my unboxing and taste test - the owner even threw in some free dessert! 🎥💕",
            videos = listOf(
                ReviewVideo("vid1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400", "3:12", "Full unboxing and taste test"),
                ReviewVideo("vid2", "https://sample-videos.com/zip/10/mp4/SampleVideo_640x360_1mb.mp4", "https://images.unsplash.com/photo-1551024506-0bccd828d307?w=400", "1:08", "Free dessert surprise!")
            ),
            hashtags = listOf("#sweet", "#incredible", "#customerservice", "#unboxing", "#freedessert", "#videoreviews")
        ),
        isVerifiedPurchase = true,
        likes = 567,
        comments = 145,
        shares = 45,
        isLiked = true,
        isBookmarked = false,
        createdAt = "2 weeks ago"
    ),

    // MIXED Review (comprehensive)
    ModelReview(
        id = "4",
        userId = "u4",
        userName = "DavidK 🍽️",
        userAvatar = null,
        targetType = ReviewTargetType.SHOP,
        targetId = shopId,
        targetName = "Bangkok Street Food",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 0.8f,
            displayValue = "4/5"
        ),
        content = ReviewContent(
            type = ReviewType.MIXED,
            text = "Been ordering from here for months and the consistency is impressive. Great quality at fair prices! Here's my comprehensive review with photos and video. 💯",
            images = listOf(
                ReviewImage("img4", "https://images.unsplash.com/photo-1565958011703-44f9829ba187?w=800", "Consistent quality over time", 1.0f),
                ReviewImage("img5", "https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=800", "Great value for money", 1.1f)
            ),
            videos = listOf(
                ReviewVideo("vid3", "https://sample-videos.com/zip/10/mp4/SampleVideo_720x480_1mb.mp4", "https://images.unsplash.com/photo-1516684669134-de6f7c473a2a?w=400", "2:45", "My monthly experience review")
            ),
            hashtags = listOf("#consistency", "#quality", "#fairprices", "#portions", "#months", "#comprehensive")
        ),
        isVerifiedPurchase = true,
        likes = 189,
        comments = 34,
        shares = 12,
        isLiked = true,
        isBookmarked = true,
        createdAt = "3 weeks ago"
    ),

    // DETAILED Review with categories
    ModelReview(
        id = "5",
        userId = "u5",
        userName = "Lisa T ✨",
        userAvatar = null,
        targetType = ReviewTargetType.SHOP,
        targetId = shopId,
        targetName = "Bangkok Street Food",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5",
            categories = mapOf(
                "Food Quality" to 1.0f,
                "Delivery Speed" to 0.9f,
                "Customer Service" to 1.0f,
                "Value for Money" to 0.9f,
                "Authenticity" to 1.0f
            )
        ),
        content = ReviewContent(
            type = ReviewType.DETAILED,
            text = "After trying every Thai restaurant in the area, this place wins hands down! Detailed breakdown: Food quality is exceptional with authentic flavors and perfect spice levels. Delivery is consistently fast and packaging preserves temperature perfectly. Customer service goes above and beyond - they remember regular customers and customize orders. Portions are generous and prices are fair for the quality. This is authentic Thai cuisine at its finest! 🇹🇭",
            images = listOf(
                ReviewImage("img6", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "Comparison with Bangkok street food", 1.0f),
                ReviewImage("img7", "https://images.unsplash.com/photo-1600891964092-4316c288032e?w=800", "Perfect spice balance", 1.1f),
                ReviewImage("img8", "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=800", "Generous portion sizes", 0.9f)
            ),
            videos = listOf(
                ReviewVideo("vid4", "https://sample-videos.com/zip/10/mp4/SampleVideo_480x270_1mb.mp4", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0?w=400", "4:56", "Complete shop review and comparison")
            ),
            hashtags = listOf("#best", "#thaifood", "#city", "#authentic", "#detailed", "#comparison", "#comprehensive"),
            mentions = listOf("@BangkokStreetFood")
        ),
        isVerifiedPurchase = true,
        likes = 678,
        comments = 198,
        shares = 67,
        isLiked = false,
        isBookmarked = true,
        createdAt = "1 month ago"
    )
)