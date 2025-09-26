package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Dialog
// import coil3.compose.AsyncImage // AsyncImage not available in this KMP setup
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
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
fun ProductDetailScreen(
    productId: String = "1",
    onBackClick: () -> Unit,
    onShopClick: (shopId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    val product = getProductById(productId)
    var selectedImageIndex by remember { mutableIntStateOf(0) }
    var quantity by remember { mutableIntStateOf(1) }
    var selectedTab by remember { mutableStateOf(0) }
    var isInWishlist by remember { mutableStateOf(false) }
    var showShareSheet by remember { mutableStateOf(false) }
    var isAddingToCart by remember { mutableStateOf(false) }
    var isBuyingNow by remember { mutableStateOf(false) }

    // Review interaction states
    var selectedImageUrl by remember { mutableStateOf<String?>(null) }
    var selectedVideoUrl by remember { mutableStateOf<String?>(null) }
    var showImageViewer by remember { mutableStateOf(false) }
    var showVideoPlayer by remember { mutableStateOf(false) }
    var showCommentsSheet by remember { mutableStateOf(false) }
    var currentReviewId by remember { mutableStateOf<String?>(null) }

    val tabs = listOf("Details", "Reviews", "Shop")

    // Reset loading states after simulated delays
    LaunchedEffect(isAddingToCart) {
        if (isAddingToCart) {
            kotlinx.coroutines.delay(2000) // Simulate 2 second cart operation
            isAddingToCart = false
        }
    }

    LaunchedEffect(isBuyingNow) {
        if (isBuyingNow) {
            kotlinx.coroutines.delay(3000) // Simulate 3 second buy now operation
            isBuyingNow = false
        }
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar
        TopAppBar(
            title = { Text("Product", fontWeight = FontWeight.Bold) },
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
                        contentDescription = "Share product",
                        tint = TchatColors.onSurface
                    )
                }
                IconButton(onClick = { isInWishlist = !isInWishlist }) {
                    Icon(
                        if (isInWishlist) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                        contentDescription = if (isInWishlist) "Remove from wishlist" else "Add to wishlist",
                        tint = if (isInWishlist) TchatColors.primary else TchatColors.onSurface
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        if (product != null) {
            LazyColumn(
                modifier = Modifier.weight(1f)
            ) {
                item {
                    // Product Images
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.md),
                        shape = RoundedCornerShape(16.dp),
                        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                    ) {
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(300.dp)
                                .background(TchatColors.primary.copy(alpha = 0.1f)),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                Icons.Default.ShoppingBag,
                                contentDescription = "Product Image",
                                modifier = Modifier.size(120.dp),
                                tint = TchatColors.primary
                            )

                            // Hot badge
                            if (product.isHot) {
                                Badge(
                                    containerColor = TchatColors.error,
                                    contentColor = TchatColors.onPrimary,
                                    modifier = Modifier
                                        .align(Alignment.TopStart)
                                        .padding(TchatSpacing.md)
                                ) {
                                    Text("HOT", fontSize = 12.sp, fontWeight = FontWeight.Bold)
                                }
                            }

                            // Discount badge
                            product.discount?.let { discount ->
                                Badge(
                                    containerColor = TchatColors.success,
                                    contentColor = TchatColors.onPrimary,
                                    modifier = Modifier
                                        .align(Alignment.TopEnd)
                                        .padding(TchatSpacing.md)
                                ) {
                                    Text("-$discount%", fontSize = 12.sp, fontWeight = FontWeight.Bold)
                                }
                            }
                        }
                    }
                }

                item {
                    // Product Info
                    Column(
                        modifier = Modifier.padding(horizontal = TchatSpacing.md)
                    ) {
                        Text(
                            product.name,
                            style = MaterialTheme.typography.headlineSmall,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )

                        Spacer(modifier = Modifier.height(TchatSpacing.xs))

                        // Merchant info - clickable
                        Row(
                            modifier = Modifier
                                .clickable {
                                    // Find shop by merchant name and navigate
                                    val shop = getDummyShops().find { it.name == product.merchant }
                                    shop?.let { onShopClick(it.id) }
                                }
                                .padding(vertical = TchatSpacing.xs),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                Icons.Default.Store,
                                contentDescription = null,
                                modifier = Modifier.size(16.dp),
                                tint = TchatColors.primary
                            )
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Text(
                                product.merchant,
                                style = MaterialTheme.typography.bodyLarge,
                                color = TchatColors.primary,
                                fontWeight = FontWeight.Medium
                            )
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Icon(
                                Icons.Default.ChevronRight,
                                contentDescription = null,
                                modifier = Modifier.size(16.dp),
                                tint = TchatColors.primary
                            )
                        }

                        // Rating and sales
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Icon(
                                    Icons.Default.Star,
                                    contentDescription = null,
                                    modifier = Modifier.size(18.dp),
                                    tint = TchatColors.warning
                                )
                                Spacer(modifier = Modifier.width(4.dp))
                                Text(
                                    "${product.rating}",
                                    style = MaterialTheme.typography.bodyLarge,
                                    fontWeight = FontWeight.Medium,
                                    color = TchatColors.onSurface
                                )
                            }

                            if (product.orders > 0) {
                                Text("•", color = TchatColors.onSurfaceVariant)
                                Text(
                                    "${product.orders} sold",
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = TchatColors.onSurfaceVariant
                                )
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Price
                        Row(
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(
                                "$${product.price}",
                                style = MaterialTheme.typography.headlineMedium,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.primary
                            )

                            product.originalPrice?.let { originalPrice ->
                                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                                Text(
                                    "$${originalPrice}",
                                    style = MaterialTheme.typography.titleMedium,
                                    color = TchatColors.onSurfaceVariant,
                                    textDecoration = TextDecoration.LineThrough
                                )
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Delivery info
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                        ) {
                            DeliveryInfoChip(
                                icon = Icons.Default.Schedule,
                                text = product.deliveryTime
                            )
                            DeliveryInfoChip(
                                icon = Icons.Default.LocationOn,
                                text = product.distance
                            )
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.lg))
                    }
                }

                item {
                    // Quantity Selector
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(horizontal = TchatSpacing.md),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.SpaceBetween
                    ) {
                        Text(
                            "Quantity:",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Medium,
                            color = TchatColors.onSurface
                        )

                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            IconButton(
                                onClick = { if (quantity > 1) quantity-- },
                                modifier = Modifier
                                    .size(40.dp)
                                    .background(TchatColors.surface, CircleShape)
                            ) {
                                Icon(
                                    Icons.Default.Remove,
                                    contentDescription = "Decrease",
                                    tint = TchatColors.onSurface
                                )
                            }

                            Text(
                                "$quantity",
                                style = MaterialTheme.typography.titleLarge,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.onSurface,
                                modifier = Modifier.widthIn(min = 40.dp)
                            )

                            IconButton(
                                onClick = { quantity++ },
                                modifier = Modifier
                                    .size(40.dp)
                                    .background(TchatColors.surface, CircleShape)
                            ) {
                                Icon(
                                    Icons.Default.Add,
                                    contentDescription = "Increase",
                                    tint = TchatColors.onSurface
                                )
                            }
                        }
                    }

                    Spacer(modifier = Modifier.height(TchatSpacing.lg))
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
                    0 -> { // Details
                        item {
                            ProductDetailsContent(product = product)
                        }
                    }
                    1 -> { // Reviews
                        items(getProductReviewsWithTypes(productId)) { review ->
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
                    2 -> { // Shop
                        item {
                            ShopInfoContent(
                                merchant = product.merchant,
                                onShopClick = onShopClick
                            )
                        }
                    }
                }
            }

            // Bottom Action Buttons
            Surface(
                modifier = Modifier.fillMaxWidth(),
                color = TchatColors.surface,
                shadowElevation = 8.dp
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(TchatSpacing.md),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    TchatButton(
                        onClick = {
                            isAddingToCart = true
                            // Simulate adding to cart with delay
                            // TODO: Implement actual cart functionality
                        },
                        text = if (isAddingToCart) "Adding..." else "Add to Cart",
                        variant = TchatButtonVariant.Secondary,
                        modifier = Modifier.weight(1f),
                        leadingIcon = if (isAddingToCart) null else {
                            { Icon(Icons.Default.ShoppingCart, contentDescription = null) }
                        },
                        loading = isAddingToCart,
                        enabled = !isAddingToCart && !isBuyingNow
                    )
                    TchatButton(
                        onClick = {
                            isBuyingNow = true
                            // Simulate buy now process with delay
                            // TODO: Implement actual buy now functionality
                        },
                        text = if (isBuyingNow) "Processing..." else "Buy Now",
                        variant = TchatButtonVariant.Primary,
                        modifier = Modifier.weight(1f),
                        leadingIcon = if (isBuyingNow) null else {
                            { Icon(Icons.Default.Payment, contentDescription = null) }
                        },
                        loading = isBuyingNow,
                        enabled = !isBuyingNow && !isAddingToCart
                    )
                }
            }
        } else {
            TchatNotFoundState(
                itemType = "Product"
            )
        }

        // Share Modal
        if (showShareSheet && product != null) {
            TchatShareModal(
                isVisible = showShareSheet,
                content = ShareContent(
                    title = product.name,
                    description = "Check out this amazing product: ${product.name} for just $${product.price}",
                    url = "https://tchat.app/products/${product.id}",
                    type = ShareContentType.PRODUCT
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

        // Image Viewer Modal
        if (showImageViewer && selectedImageUrl != null) {
            Dialog(
                onDismissRequest = {
                    showImageViewer = false
                    selectedImageUrl = null
                }
            ) {
                Card(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(16.dp),
                    colors = CardDefaults.cardColors(containerColor = Color.Black)
                ) {
                    Box(
                        modifier = Modifier.fillMaxSize()
                    ) {
                        // Image placeholder - would use AsyncImage with proper image loading library
                        Box(
                            modifier = Modifier
                                .fillMaxSize()
                                .background(Color.Gray)
                                .clickable {
                                    showImageViewer = false
                                    selectedImageUrl = null
                                },
                            contentAlignment = Alignment.Center
                        ) {
                            Column(
                                horizontalAlignment = Alignment.CenterHorizontally
                            ) {
                                Icon(
                                    imageVector = Icons.Default.Image,
                                    contentDescription = "Image",
                                    tint = Color.White,
                                    modifier = Modifier.size(64.dp)
                                )
                                Spacer(modifier = Modifier.height(8.dp))
                                Text(
                                    "Review Image",
                                    color = Color.White,
                                    style = MaterialTheme.typography.titleMedium
                                )
                                Text(
                                    "URL: $selectedImageUrl",
                                    color = Color.White,
                                    style = MaterialTheme.typography.bodySmall
                                )
                            }
                        }

                        // Close button
                        IconButton(
                            onClick = {
                                showImageViewer = false
                                selectedImageUrl = null
                            },
                            modifier = Modifier
                                .align(Alignment.TopEnd)
                                .padding(16.dp)
                        ) {
                            Icon(
                                imageVector = Icons.Default.Close,
                                contentDescription = "Close",
                                tint = Color.White
                            )
                        }
                    }
                }
            }
        }

        // Video Player Modal
        if (showVideoPlayer && selectedVideoUrl != null) {
            Dialog(
                onDismissRequest = {
                    showVideoPlayer = false
                    selectedVideoUrl = null
                }
            ) {
                Card(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(300.dp)
                        .padding(16.dp)
                ) {
                    Box(
                        modifier = Modifier.fillMaxSize()
                    ) {
                        // Placeholder for video player
                        Column(
                            modifier = Modifier
                                .fillMaxSize()
                                .background(Color.Black)
                                .padding(16.dp),
                            horizontalAlignment = Alignment.CenterHorizontally,
                            verticalArrangement = Arrangement.Center
                        ) {
                            Icon(
                                imageVector = Icons.Default.PlayArrow,
                                contentDescription = "Play Video",
                                tint = Color.White,
                                modifier = Modifier.size(64.dp)
                            )
                            Spacer(modifier = Modifier.height(8.dp))
                            Text(
                                "Video Player",
                                color = Color.White,
                                style = MaterialTheme.typography.titleMedium
                            )
                            Text(
                                "URL: $selectedVideoUrl",
                                color = Color.White,
                                style = MaterialTheme.typography.bodySmall
                            )
                        }

                        // Close button
                        IconButton(
                            onClick = {
                                showVideoPlayer = false
                                selectedVideoUrl = null
                            },
                            modifier = Modifier
                                .align(Alignment.TopEnd)
                                .padding(8.dp)
                        ) {
                            Icon(
                                imageVector = Icons.Default.Close,
                                contentDescription = "Close",
                                tint = Color.White
                            )
                        }
                    }
                }
            }
        }

        // Comments Bottom Sheet
        if (showCommentsSheet && currentReviewId != null) {
            ModalBottomSheet(
                onDismissRequest = {
                    showCommentsSheet = false
                    currentReviewId = null
                },
                sheetState = rememberModalBottomSheetState(skipPartiallyExpanded = true)
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp)
                ) {
                    Text(
                        "Comments",
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.Bold
                    )

                    Spacer(modifier = Modifier.height(16.dp))

                    // Sample comments
                    repeat(5) { index ->
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(vertical = 8.dp)
                        ) {
                            // User avatar placeholder
                            Box(
                                modifier = Modifier
                                    .size(40.dp)
                                    .clip(CircleShape)
                                    .background(TchatColors.primary),
                                contentAlignment = Alignment.Center
                            ) {
                                Icon(
                                    imageVector = Icons.Default.Person,
                                    contentDescription = "User Avatar",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }

                            Spacer(modifier = Modifier.width(12.dp))

                            Column {
                                Text(
                                    "User ${index + 1}",
                                    style = MaterialTheme.typography.titleSmall,
                                    fontWeight = FontWeight.Bold
                                )
                                Text(
                                    "This is a sample comment for review $currentReviewId. Great review!",
                                    style = MaterialTheme.typography.bodyMedium
                                )
                                Text(
                                    "${index + 1}h ago",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = TchatColors.onSurfaceVariant
                                )
                            }
                        }
                    }

                    Spacer(modifier = Modifier.height(100.dp)) // Space for keyboard
                }
            }
        }
    }
}

@Composable
private fun DeliveryInfoChip(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    text: String,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(20.dp),
        color = TchatColors.surface
    ) {
        Row(
            modifier = Modifier.padding(horizontal = TchatSpacing.sm, vertical = TchatSpacing.xs),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            Icon(
                icon,
                contentDescription = null,
                modifier = Modifier.size(16.dp),
                tint = TchatColors.onSurfaceVariant
            )
            Text(
                text,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurfaceVariant
            )
        }
    }
}

@Composable
private fun ProductDetailsContent(
    product: ProductItem,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier.padding(TchatSpacing.md)
    ) {
        Text(
            "Product Details",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )

        Spacer(modifier = Modifier.height(TchatSpacing.sm))

        Text(
            "Experience authentic ${product.category.lowercase()} from ${product.merchant}. " +
            "This premium ${product.name.lowercase()} is carefully prepared with traditional methods " +
            "and the finest ingredients to deliver exceptional taste and quality.",
            style = MaterialTheme.typography.bodyLarge,
            color = TchatColors.onSurface,
            lineHeight = 24.sp
        )

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        Text(
            "Specifications",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )

        Spacer(modifier = Modifier.height(TchatSpacing.sm))

        ProductSpecRow("Category", product.category)
        ProductSpecRow("Merchant", product.merchant)
        ProductSpecRow("Rating", "${product.rating} stars")
        ProductSpecRow("Delivery", "${product.deliveryTime} • ${product.distance}")
        if (product.orders > 0) {
            ProductSpecRow("Sold", "${product.orders} times")
        }
    }
}

@Composable
private fun ProductSpecRow(
    label: String,
    value: String,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
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

@Composable
private fun ShopInfoContent(
    merchant: String,
    onShopClick: (shopId: String) -> Unit,
    modifier: Modifier = Modifier
) {
    val shop = getDummyShops().find { it.name == merchant }

    shop?.let { shopData ->
        Card(
            modifier = modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md)
                .clickable { onShopClick(shopData.id) },
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
            elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Shop Avatar
                Box(
                    modifier = Modifier
                        .size(60.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primaryLight),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        shopData.name.first().toString(),
                        color = TchatColors.onPrimary,
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.md))

                Column(modifier = Modifier.weight(1f)) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            shopData.name,
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )
                        if (shopData.isVerified) {
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Icon(
                                Icons.Default.Verified,
                                contentDescription = "Verified",
                                modifier = Modifier.size(18.dp),
                                tint = TchatColors.primary
                            )
                        }
                    }

                    Text(
                        shopData.description,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        maxLines = 2
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.xs))

                    Row(
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                    ) {
                        ShopStatChip(
                            icon = Icons.Default.Star,
                            value = "${shopData.rating}",
                            color = TchatColors.warning
                        )
                        ShopStatChip(
                            icon = Icons.Default.People,
                            value = "${shopData.followers}",
                            color = TchatColors.primary
                        )
                        ShopStatChip(
                            icon = Icons.Default.Inventory,
                            value = "${shopData.totalProducts}",
                            color = TchatColors.success
                        )
                    }
                }

                Icon(
                    Icons.Default.ChevronRight,
                    contentDescription = null,
                    tint = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun ShopStatChip(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    value: String,
    color: androidx.compose.ui.graphics.Color,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(2.dp)
    ) {
        Icon(
            icon,
            contentDescription = null,
            modifier = Modifier.size(14.dp),
            tint = color
        )
        Text(
            value,
            style = MaterialTheme.typography.bodySmall,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface
        )
    }
}

// Helper functions and data
private fun getProductById(productId: String): ProductItem? {
    return getDummyProducts().find { it.id == productId }
}

private fun getProductReviews(): List<ShopReview> = listOf(
    ShopReview("1", "Sarah Chen", 5, "Excellent product! Fast delivery and great quality.", "1 day ago"),
    ShopReview("2", "Mike Johnson", 4, "Good value for money. Would recommend.", "3 days ago"),
    ShopReview("3", "Anna Wong", 5, "Amazing taste! Will definitely order again.", "1 week ago"),
    ShopReview("4", "David Kim", 4, "Quick delivery and well packaged.", "2 weeks ago")
)

private fun getProductReviewsWithTypes(productId: String): List<ModelReview> = listOf(
    // TEXT Review
    ModelReview(
        id = "1",
        userId = "u1",
        userName = "Sarah Chen ✨",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5"
        ),
        content = ReviewContent(
            type = ReviewType.TEXT,
            text = "OMG this product is absolutely amazing! The quality is top-notch and delivery was super fast. I've been using it daily and it's become my holy grail! Already ordered 3 more for my friends. 😍",
            hashtags = listOf("#loveit", "#quality", "#fastdelivery", "#holygrail", "#obsessed")
        ),
        isVerifiedPurchase = true,
        likes = 234,
        comments = 45,
        shares = 12,
        isLiked = false,
        isBookmarked = false,
        createdAt = "1 day ago"
    ),

    // IMAGE Review
    ModelReview(
        id = "2",
        userId = "u2",
        userName = "Mike J 👑",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 0.8f,
            displayValue = "4/5"
        ),
        content = ReviewContent(
            type = ReviewType.IMAGE,
            text = "Really good product for the price! Love how it looks in real life. Packaging was nice and it arrived earlier than expected. 📸✨",
            images = listOf(
                ReviewImage("img1", "https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=800", "Product unboxing", 1.2f),
                ReviewImage("img2", "https://images.unsplash.com/photo-1546793665-c74683f339c1?w=800", "Close up shot", 1.0f),
                ReviewImage("img3", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "In use photo", 0.8f)
            ),
            hashtags = listOf("#goodvalue", "#recommend", "#packaging", "#earlydelivery", "#photoreviews")
        ),
        isVerifiedPurchase = true,
        likes = 156,
        comments = 23,
        shares = 8,
        isLiked = true,
        isBookmarked = true,
        createdAt = "3 days ago"
    ),

    // VIDEO Review
    ModelReview(
        id = "3",
        userId = "u3",
        userName = "Anna Wong 🌸",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5"
        ),
        content = ReviewContent(
            type = ReviewType.VIDEO,
            text = "This has become my new favorite! Check out my unboxing and first impressions video. The taste is incredible! 🎥",
            videos = listOf(
                ReviewVideo("vid1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400", "2:34", "Unboxing and first taste"),
                ReviewVideo("vid2", "https://sample-videos.com/zip/10/mp4/SampleVideo_640x360_1mb.mp4", "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=400", "1:45", "How I prepare it")
            ),
            hashtags = listOf("#videoreviews", "#unboxing", "#tastetesting", "#cooking", "#favorite")
        ),
        isVerifiedPurchase = true,
        likes = 298,
        comments = 67,
        shares = 24,
        isLiked = false,
        isBookmarked = true,
        createdAt = "1 week ago"
    ),

    // MIXED Review (Text + Images + Videos)
    ModelReview(
        id = "4",
        userId = "u4",
        userName = "David Kim 🔥",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 0.8f,
            displayValue = "4/5"
        ),
        content = ReviewContent(
            type = ReviewType.MIXED,
            text = "Impressed with the quick delivery and careful packaging. Product quality is solid and exactly what I expected. Check out my photos and quick video review! 📦🎬",
            images = listOf(
                ReviewImage("img4", "https://images.unsplash.com/photo-1586190848861-99aa4a171e90?w=800", "Great packaging", 1.0f),
                ReviewImage("img5", "https://images.unsplash.com/photo-1565958011703-44f9829ba187?w=800", "Product quality", 1.2f)
            ),
            videos = listOf(
                ReviewVideo("vid3", "https://sample-videos.com/zip/10/mp4/SampleVideo_720x480_1mb.mp4", "https://images.unsplash.com/photo-1516684669134-de6f7c473a2a?w=400", "0:45", "Quick thoughts")
            ),
            hashtags = listOf("#fastdelivery", "#wellpackaged", "#solid", "#mixedmedia", "#comprehensive")
        ),
        isVerifiedPurchase = true,
        likes = 87,
        comments = 12,
        shares = 15,
        isLiked = true,
        isBookmarked = false,
        createdAt = "2 weeks ago"
    ),

    // QUICK Review
    ModelReview(
        id = "5",
        userId = "u5",
        userName = "Lisa Quick ⚡",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 1.0f,
            displayValue = "5/5"
        ),
        content = ReviewContent(
            type = ReviewType.QUICK,
            text = "Love it! Quick delivery, tastes amazing! 😋⚡",
            hashtags = listOf("#quick", "#amazing", "#love")
        ),
        isVerifiedPurchase = true,
        likes = 45,
        comments = 3,
        shares = 2,
        isLiked = false,
        isBookmarked = false,
        createdAt = "5 days ago"
    ),

    // DETAILED Review
    ModelReview(
        id = "6",
        userId = "u6",
        userName = "Food Critic Pro 🍴",
        userAvatar = null,
        targetType = ReviewTargetType.PRODUCT,
        targetId = productId,
        targetName = "Pad Thai Goong",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = 0.9f,
            displayValue = "4.5/5",
            categories = mapOf(
                "Taste" to 1.0f,
                "Presentation" to 0.9f,
                "Value" to 0.8f,
                "Delivery" to 0.9f
            )
        ),
        content = ReviewContent(
            type = ReviewType.DETAILED,
            text = "After trying this for a full month, here's my comprehensive review: The authentic flavors really shine through, especially the tamarind balance. Shrimp quality is consistently fresh. Presentation is restaurant-level. Only minor issue is slight variation in spice levels between orders, but overall exceptional experience!",
            images = listOf(
                ReviewImage("img6", "https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=800", "Plating presentation", 1.0f),
                ReviewImage("img7", "https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=800", "Ingredient quality", 1.2f),
                ReviewImage("img8", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "Portion size", 0.9f)
            ),
            videos = listOf(
                ReviewVideo("vid4", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_2mb.mp4", "https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400", "3:22", "Full detailed review")
            ),
            hashtags = listOf("#detailedreview", "#foodcritic", "#authentic", "#comprehensive", "#monthlytest"),
            mentions = listOf("@BangkokStreetFood")
        ),
        isVerifiedPurchase = true,
        likes = 456,
        comments = 89,
        shares = 34,
        isLiked = false,
        isBookmarked = true,
        createdAt = "3 days ago"
    )
)

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