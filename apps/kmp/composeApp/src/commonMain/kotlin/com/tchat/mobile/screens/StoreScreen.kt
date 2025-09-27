package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
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
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.components.TchatTopBar
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

data class ShopItem(
    val id: String,
    val name: String,
    val description: String,
    val avatar: String,
    val coverImage: String,
    val rating: Double,
    val deliveryTime: String,
    val distance: String,
    val isVerified: Boolean = false,
    val followers: Int = 0,
    val totalProducts: Int = 0
)

data class ProductItem(
    val id: String,
    val name: String,
    val price: Double,
    val originalPrice: Double? = null,
    val rating: Double,
    val category: String,
    val merchant: String,
    val image: String,
    val deliveryTime: String = "30 min",
    val distance: String = "2.5 km",
    val isHot: Boolean = false,
    val discount: Int? = null,
    val orders: Int = 0
)

data class LiveStreamItem(
    val id: String,
    val title: String,
    val merchant: String,
    val viewers: Int,
    val thumbnail: String,
    val products: List<ProductItem>,
    val isLive: Boolean = true,
    val duration: String = "00:45:32"
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StoreScreen(
    onProductClick: (productId: String) -> Unit = {},
    onShopClick: (shopId: String) -> Unit = {},
    onLiveStreamClick: (streamId: String) -> Unit = {},
    onMoreClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var selectedTab by remember { mutableStateOf(0) }
    var searchQuery by remember { mutableStateOf("") }

    val tabs = listOf("Shops", "Products", "Live")

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar
        TopAppBar(
            title = { Text("Store", fontWeight = FontWeight.Bold) },
            actions = {
                IconButton(onClick = { /* Cart */ }) {
                    BadgedBox(
                        badge = {
                            Badge(
                                containerColor = TchatColors.error,
                                contentColor = TchatColors.onPrimary
                            ) {
                                Text("3", fontSize = 10.sp)
                            }
                        }
                    ) {
                        Icon(
                            Icons.Filled.ShoppingCart,
                            "Cart",
                            tint = TchatColors.onSurface
                        )
                    }
                }
                IconButton(onClick = { /* Filter */ }) {
                    Icon(
                        Icons.Filled.FilterList,
                        "Filter",
                        tint = TchatColors.onSurface
                    )
                }
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

        // Search Bar
        TchatInput(
            value = searchQuery,
            onValueChange = { searchQuery = it },
            type = TchatInputType.Search,
            placeholder = "Search shops, products, live streams...",
            leadingIcon = Icons.Default.Search,
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md)
        )

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
                    },
                    icon = {
                        when (index) {
                            0 -> Icon(Icons.Default.Store, contentDescription = null, modifier = Modifier.size(20.dp))
                            1 -> Icon(Icons.Default.ShoppingCart, contentDescription = null, modifier = Modifier.size(20.dp))
                            2 -> {
                                Box {
                                    Icon(Icons.Default.PlayCircle, contentDescription = null, modifier = Modifier.size(20.dp))
                                    if (index == 2) { // Live indicator
                                        Box(
                                            modifier = Modifier
                                                .size(8.dp)
                                                .background(TchatColors.error, CircleShape)
                                                .offset(x = 6.dp, y = (-6).dp)
                                        )
                                    }
                                }
                            }
                        }
                    }
                )
            }
        }

        // Content based on selected tab
        when (selectedTab) {
            0 -> ShopsContent(
                searchQuery = searchQuery,
                onShopClick = onShopClick,
                modifier = Modifier.weight(1f)
            )
            1 -> ProductsContent(
                searchQuery = searchQuery,
                onProductClick = onProductClick,
                modifier = Modifier.weight(1f)
            )
            2 -> LiveContent(
                searchQuery = searchQuery,
                onLiveStreamClick = onLiveStreamClick,
                onProductClick = onProductClick,
                modifier = Modifier.weight(1f)
            )
        }
    }
}

@Composable
private fun ShopsContent(
    searchQuery: String,
    onShopClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val shops = getDummyShops().filter {
        it.name.contains(searchQuery, ignoreCase = true) ||
        it.description.contains(searchQuery, ignoreCase = true)
    }

    LazyColumn(
        modifier = modifier,
        contentPadding = PaddingValues(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        items(shops) { shop ->
            ShopCard(
                shop = shop,
                onClick = { onShopClick(shop.id) }
            )
        }
    }
}

@Composable
private fun ProductsContent(
    searchQuery: String,
    onProductClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val products = getDummyProducts().filter {
        it.name.contains(searchQuery, ignoreCase = true) ||
        it.category.contains(searchQuery, ignoreCase = true) ||
        it.merchant.contains(searchQuery, ignoreCase = true)
    }

    LazyVerticalGrid(
        columns = GridCells.Fixed(2),
        modifier = modifier.fillMaxSize(),
        contentPadding = PaddingValues(TchatSpacing.md),
        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        items(products) { product ->
            ProductCard(
                product = product,
                onClick = { onProductClick(product.id) }
            )
        }
    }
}

@Composable
private fun LiveContent(
    searchQuery: String,
    onLiveStreamClick: (String) -> Unit,
    onProductClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val streams = getDummyLiveStreams().filter {
        it.title.contains(searchQuery, ignoreCase = true) ||
        it.merchant.contains(searchQuery, ignoreCase = true)
    }

    LazyColumn(
        modifier = modifier,
        contentPadding = PaddingValues(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
    ) {
        items(streams) { stream ->
            LiveStreamCard(
                stream = stream,
                onStreamClick = { onLiveStreamClick(stream.id) },
                onProductClick = onProductClick
            )
        }
    }
}

@Composable
private fun ShopCard(
    shop: ShopItem,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onClick,
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            // Cover Image
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .background(TchatColors.primary.copy(alpha = 0.1f)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.Store,
                    contentDescription = "Shop Cover",
                    modifier = Modifier.size(48.dp),
                    tint = TchatColors.primary
                )
            }

            // Shop Info
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Avatar
                Box(
                    modifier = Modifier
                        .size(48.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primaryLight),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        shop.name.first().toString(),
                        color = TchatColors.onPrimary,
                        fontWeight = FontWeight.Bold
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                Column(modifier = Modifier.weight(1f)) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            shop.name,
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )
                        if (shop.isVerified) {
                            Spacer(modifier = Modifier.width(4.dp))
                            Icon(
                                Icons.Default.Verified,
                                contentDescription = "Verified",
                                modifier = Modifier.size(16.dp),
                                tint = TchatColors.primary
                            )
                        }
                    }

                    Text(
                        shop.description,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        maxLines = 1
                    )

                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        Icon(
                            Icons.Default.Star,
                            contentDescription = null,
                            modifier = Modifier.size(14.dp),
                            tint = TchatColors.warning
                        )
                        Text(
                            "${shop.rating}",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                        Text("•", color = TchatColors.onSurfaceVariant, fontSize = 10.sp)
                        Text(
                            "${shop.followers} followers",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                        Text("•", color = TchatColors.onSurfaceVariant, fontSize = 10.sp)
                        Text(
                            "${shop.totalProducts} products",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
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
private fun ProductCard(
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
                        .height(120.dp)
                        .background(TchatColors.primary.copy(alpha = 0.1f)),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        Icons.Default.ShoppingBag,
                        contentDescription = "Product Image",
                        modifier = Modifier.size(48.dp),
                        tint = TchatColors.primary
                    )
                }

                // Hot badge
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

                // Discount badge
                product.discount?.let { discount ->
                    Badge(
                        containerColor = TchatColors.success,
                        contentColor = TchatColors.onPrimary,
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .padding(TchatSpacing.xs)
                    ) {
                        Text("-$discount%", fontSize = 8.sp, fontWeight = FontWeight.Bold)
                    }
                }
            }

            // Product Info
            Column(
                modifier = Modifier.padding(TchatSpacing.sm)
            ) {
                Text(
                    text = product.name,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2
                )

                Spacer(modifier = Modifier.height(2.dp))

                Text(
                    text = product.merchant,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    maxLines = 1
                )

                Spacer(modifier = Modifier.height(4.dp))

                // Rating and orders
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    Icon(
                        Icons.Filled.Star,
                        contentDescription = "Rating",
                        modifier = Modifier.size(12.dp),
                        tint = TchatColors.warning
                    )
                    Text(
                        text = "${product.rating}",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                    if (product.orders > 0) {
                        Text("•", color = TchatColors.onSurfaceVariant, fontSize = 8.sp)
                        Text(
                            text = "${product.orders} sold",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                Spacer(modifier = Modifier.height(6.dp))

                // Price
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(
                                text = "$${product.price}",
                                style = MaterialTheme.typography.titleSmall,
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
                        Text(
                            text = "${product.deliveryTime} • ${product.distance}",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    Icon(
                        Icons.Default.AddCircle,
                        contentDescription = "Add to Cart",
                        modifier = Modifier.size(24.dp),
                        tint = TchatColors.primary
                    )
                }
            }
        }
    }
}

@Composable
private fun LiveStreamCard(
    stream: LiveStreamItem,
    onStreamClick: () -> Unit,
    onProductClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onStreamClick,
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            // Stream thumbnail with live indicator
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(160.dp)
                    .background(
                        if (stream.isLive) TchatColors.error.copy(alpha = 0.1f)
                        else TchatColors.primary.copy(alpha = 0.1f)
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.PlayCircle,
                    contentDescription = "Live Stream",
                    modifier = Modifier.size(64.dp),
                    tint = if (stream.isLive) TchatColors.error else TchatColors.primary
                )

                // Live indicator
                if (stream.isLive) {
                    Row(
                        modifier = Modifier
                            .align(Alignment.TopStart)
                            .padding(TchatSpacing.sm)
                            .background(TchatColors.error, RoundedCornerShape(12.dp))
                            .padding(horizontal = 8.dp, vertical = 4.dp),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        Box(
                            modifier = Modifier
                                .size(6.dp)
                                .background(TchatColors.onPrimary, CircleShape)
                        )
                        Text(
                            "LIVE",
                            color = TchatColors.onPrimary,
                            style = MaterialTheme.typography.labelSmall,
                            fontWeight = FontWeight.Bold
                        )
                    }
                }

                // Viewers count
                Row(
                    modifier = Modifier
                        .align(Alignment.TopEnd)
                        .padding(TchatSpacing.sm)
                        .background(TchatColors.onSurface.copy(alpha = 0.7f), RoundedCornerShape(12.dp))
                        .padding(horizontal = 8.dp, vertical = 4.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    Icon(
                        Icons.Default.Visibility,
                        contentDescription = null,
                        modifier = Modifier.size(12.dp),
                        tint = TchatColors.onPrimary
                    )
                    Text(
                        "${stream.viewers}",
                        color = TchatColors.onPrimary,
                        style = MaterialTheme.typography.labelSmall
                    )
                }

                // Duration
                Text(
                    stream.duration,
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(TchatSpacing.sm)
                        .background(TchatColors.onSurface.copy(alpha = 0.7f), RoundedCornerShape(4.dp))
                        .padding(horizontal = 4.dp, vertical = 2.dp),
                    color = TchatColors.onPrimary,
                    style = MaterialTheme.typography.labelSmall
                )
            }

            // Stream info
            Column(
                modifier = Modifier.padding(TchatSpacing.md)
            ) {
                Text(
                    stream.title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface,
                    maxLines = 2
                )

                Spacer(modifier = Modifier.height(4.dp))

                Text(
                    stream.merchant,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )

                if (stream.products.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(8.dp))

                    Text(
                        "Featured Products:",
                        style = MaterialTheme.typography.bodySmall,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(4.dp))

                    // Featured products row
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        stream.products.take(3).forEach { product ->
                            Card(
                                onClick = { onProductClick(product.id) },
                                modifier = Modifier.weight(1f),
                                shape = RoundedCornerShape(8.dp),
                                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
                            ) {
                                Column(
                                    modifier = Modifier.padding(TchatSpacing.xs),
                                    horizontalAlignment = Alignment.CenterHorizontally
                                ) {
                                    Box(
                                        modifier = Modifier
                                            .size(40.dp)
                                            .background(TchatColors.primary.copy(alpha = 0.1f), RoundedCornerShape(6.dp)),
                                        contentAlignment = Alignment.Center
                                    ) {
                                        Icon(
                                            Icons.Default.ShoppingBag,
                                            contentDescription = null,
                                            modifier = Modifier.size(20.dp),
                                            tint = TchatColors.primary
                                        )
                                    }
                                    Spacer(modifier = Modifier.height(2.dp))
                                    Text(
                                        "$${product.price}",
                                        style = MaterialTheme.typography.labelSmall,
                                        fontWeight = FontWeight.Bold,
                                        color = TchatColors.primary
                                    )
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// Sample data functions
private fun getDummyShops(): List<ShopItem> = listOf(
    ShopItem("1", "Bangkok Street Food", "Authentic Thai street food and snacks", "BS", "cover1", 4.8, "15-25 min", "1.2 km", true, 2540, 45),
    ShopItem("2", "Tech Paradise", "Latest electronics and gadgets", "TP", "cover2", 4.6, "30-45 min", "3.1 km", true, 1890, 120),
    ShopItem("3", "Fashion House", "Trendy clothes and accessories", "FH", "cover3", 4.5, "20-30 min", "2.8 km", false, 980, 78),
    ShopItem("4", "Fresh Market", "Organic fruits and vegetables", "FM", "cover4", 4.7, "10-20 min", "0.8 km", true, 3200, 230),
    ShopItem("5", "Coffee Corner", "Premium coffee and desserts", "CC", "cover5", 4.9, "5-15 min", "0.5 km", true, 1560, 35)
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

private fun getDummyLiveStreams(): List<LiveStreamItem> = listOf(
    LiveStreamItem(
        "1",
        "Live Cooking: Authentic Thai Pad Thai",
        "Bangkok Street Food",
        1240,
        "thumb1",
        listOf(
            ProductItem("p1", "Pad Thai Kit", 299.0, 399.0, 4.8, "Food", "Bangkok Street Food", "kit1"),
            ProductItem("p2", "Thai Sauce Set", 150.0, null, 4.6, "Food", "Bangkok Street Food", "sauce1"),
            ProductItem("p3", "Premium Noodles", 89.0, null, 4.5, "Food", "Bangkok Street Food", "noodles1")
        ),
        duration = "01:23:45"
    ),
    LiveStreamItem(
        "2",
        "New iPhone 15 Unboxing & Review",
        "Tech Paradise",
        856,
        "thumb2",
        listOf(
            ProductItem("p4", "iPhone 15 Pro", 39900.0, 42900.0, 4.9, "Electronics", "Tech Paradise", "iphone1"),
            ProductItem("p5", "Phone Case", 890.0, null, 4.4, "Electronics", "Tech Paradise", "case1"),
            ProductItem("p6", "Wireless Charger", 1290.0, 1590.0, 4.3, "Electronics", "Tech Paradise", "charger1")
        ),
        duration = "00:45:12"
    ),
    LiveStreamItem(
        "3",
        "Fresh Market Tour: Organic Selection",
        "Fresh Market",
        432,
        "thumb3",
        listOf(
            ProductItem("p7", "Organic Mango Box", 450.0, null, 4.8, "Food", "Fresh Market", "mango1"),
            ProductItem("p8", "Fresh Vegetable Set", 280.0, 350.0, 4.6, "Food", "Fresh Market", "veggie1")
        ),
        duration = "00:28:30"
    )
)