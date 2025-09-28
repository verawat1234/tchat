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
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.Send
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
import com.tchat.mobile.components.TchatNotFoundState
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.ShopItem

data class LiveStreamData(
    val id: String,
    val title: String,
    val hostName: String,
    val hostAvatar: String,
    val viewerCount: Int,
    val isLive: Boolean,
    val thumbnailUrl: String,
    val category: String,
    val shopId: String,
    val shopName: String,
    val shopAvatar: String,
    val shopRating: Double,
    val featuredProducts: List<ProductInStream>,
    val description: String,
    val startedAt: String,
    val tags: List<String> = emptyList()
)

data class ProductInStream(
    val id: String,
    val name: String,
    val description: String,
    val price: String,
    val originalPrice: String? = null,
    val discount: String? = null,
    val imageUrl: String,
    val thumbnailUrl: String,
    val inStock: Boolean = true,
    val stockQuantity: Int,
    val rating: Double,
    val reviewCount: Int,
    val category: String,
    val tags: List<String> = emptyList(),
    val specifications: Map<String, String> = emptyMap()
)

data class ChatMessage(
    val id: String,
    val username: String,
    val message: String,
    val timestamp: String,
    val isFromHost: Boolean = false
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LiveStreamScreen(
    streamId: String,
    onBackClick: () -> Unit,
    onProductClick: (productId: String) -> Unit = {},
    onShopClick: (shopId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    val streamData = getLiveStreamById(streamId)
    var chatMessage by remember { mutableStateOf("") }
    var showProducts by remember { mutableStateOf(false) }
    val chatMessages = remember { mutableStateListOf<ChatMessage>().apply { addAll(getDummyChatMessages()) } }
    var isFollowingHost by remember { mutableStateOf(false) }
    var showShareSheet by remember { mutableStateOf(false) }
    var isSendingMessage by remember { mutableStateOf(false) }

    if (streamData == null) {
        Column(
            modifier = modifier
                .fillMaxSize()
                .background(TchatColors.background)
        ) {
            // Top App Bar with back button
            TopAppBar(
                title = { Text("Live Stream") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface,
                    navigationIconContentColor = TchatColors.onSurface
                )
            )

            TchatNotFoundState(
                itemType = "Live Stream",
                modifier = Modifier.weight(1f)
            )
        }
        return
    }

    Box(modifier = modifier.fillMaxSize()) {
        // Main content
        Column(
            modifier = Modifier.fillMaxSize()
        ) {
            // Top bar with back button and stream info
            TopAppBar(
                title = { Text(streamData.title, maxLines = 1, overflow = TextOverflow.Ellipsis) },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                },
                actions = {
                    IconButton(onClick = {
                        showShareSheet = true
                        // TODO: Implement platform-specific sharing
                    }) {
                        Icon(Icons.Default.Share, contentDescription = "Share live stream")
                    }
                    IconButton(onClick = { }) {
                        Icon(Icons.Default.MoreVert, contentDescription = "More options")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface,
                    navigationIconContentColor = TchatColors.onSurface,
                    actionIconContentColor = TchatColors.onSurface
                )
            )

            // Live stream video area (mock)
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(250.dp)
                    .background(Color.Black),
                contentAlignment = Alignment.Center
            ) {
                // Mock video player
                Text(
                    "🔴 LIVE STREAM",
                    color = Color.White,
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold
                )

                // Live indicator and viewer count
                Row(
                    modifier = Modifier
                        .align(Alignment.TopStart)
                        .padding(TchatSpacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .background(TchatColors.error, RoundedCornerShape(12.dp))
                            .padding(horizontal = 8.dp, vertical = 4.dp)
                    ) {
                        Text(
                            "LIVE",
                            color = Color.White,
                            style = MaterialTheme.typography.labelSmall,
                            fontWeight = FontWeight.Bold
                        )
                    }
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                    Box(
                        modifier = Modifier
                            .background(Color.Black.copy(alpha = 0.7f), RoundedCornerShape(12.dp))
                            .padding(horizontal = 8.dp, vertical = 4.dp)
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Icon(
                                Icons.Default.RemoveRedEye,
                                contentDescription = null,
                                tint = Color.White,
                                modifier = Modifier.size(14.dp)
                            )
                            Spacer(modifier = Modifier.width(4.dp))
                            Text(
                                "${streamData.viewerCount}",
                                color = Color.White,
                                style = MaterialTheme.typography.labelSmall
                            )
                        }
                    }
                }

                // Products button
                FloatingActionButton(
                    onClick = { showProducts = !showProducts },
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .padding(TchatSpacing.md)
                        .size(48.dp),
                    containerColor = TchatColors.primary,
                    contentColor = Color.White
                ) {
                    Icon(Icons.Default.ShoppingBag, contentDescription = "Products")
                }
            }

            // Host info and description
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        // Host avatar
                        Box(
                            modifier = Modifier
                                .size(40.dp)
                                .background(TchatColors.primary, CircleShape),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                streamData.hostName.first().toString(),
                                color = Color.White,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold
                            )
                        }

                        Spacer(modifier = Modifier.width(TchatSpacing.sm))

                        Column(modifier = Modifier.weight(1f)) {
                            Text(
                                streamData.hostName,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Medium,
                                color = TchatColors.onSurface
                            )
                            Text(
                                streamData.shopName,
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant,
                                modifier = Modifier.clickable {
                                    // Navigate to shop - find shop ID by name
                                    val shop = getDummyShops().find { it.name == streamData.shopName }
                                    shop?.let { onShopClick(it.id) }
                                }
                            )
                        }

                        Button(
                            onClick = { isFollowingHost = !isFollowingHost },
                            colors = ButtonDefaults.buttonColors(
                                containerColor = if (isFollowingHost) TchatColors.surface else TchatColors.primary
                            ),
                            modifier = Modifier.height(32.dp)
                        ) {
                            if (isFollowingHost) {
                                Row(
                                    verticalAlignment = Alignment.CenterVertically,
                                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                                ) {
                                    Icon(
                                        Icons.Default.Check,
                                        contentDescription = null,
                                        modifier = Modifier.size(14.dp),
                                        tint = TchatColors.onSurface
                                    )
                                    Text(
                                        "Following",
                                        style = MaterialTheme.typography.labelSmall,
                                        color = TchatColors.onSurface
                                    )
                                }
                            } else {
                                Text(
                                    "Follow",
                                    style = MaterialTheme.typography.labelSmall,
                                    color = Color.White
                                )
                            }
                        }
                    }

                    Spacer(modifier = Modifier.height(TchatSpacing.sm))

                    Text(
                        streamData.title,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurface,
                        maxLines = 2,
                        overflow = TextOverflow.Ellipsis
                    )
                }
            }

            // Chat section
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f)
                    .padding(horizontal = TchatSpacing.md)
            ) {
                Text(
                    "Live Chat",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )

                LazyColumn(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs),
                    reverseLayout = true
                ) {
                    items(chatMessages.reversed()) { message ->
                        ChatMessageItem(message = message)
                    }
                }
            }

            // Chat input
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                OutlinedTextField(
                    value = chatMessage,
                    onValueChange = { chatMessage = it },
                    placeholder = { Text("Type a message...") },
                    modifier = Modifier.weight(1f),
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = TchatColors.primary,
                        unfocusedBorderColor = TchatColors.outline
                    )
                )

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                IconButton(
                    onClick = {
                        if (chatMessage.isNotBlank() && !isSendingMessage) {
                            isSendingMessage = true
                            // Add the new message to chat
                            val newMessage = ChatMessage(
                                id = (chatMessages.size + 1).toString(),
                                username = "You",
                                message = chatMessage,
                                timestamp = "now",
                                isFromHost = false
                            )
                            chatMessages.add(0, newMessage) // Add to top for reversed layout
                            chatMessage = ""
                            isSendingMessage = false
                        }
                    },
                    enabled = chatMessage.isNotBlank() && !isSendingMessage,
                    modifier = Modifier
                        .background(
                            if (chatMessage.isNotBlank() && !isSendingMessage)
                                TchatColors.primary
                            else
                                TchatColors.onSurfaceVariant.copy(alpha = 0.5f),
                            CircleShape
                        )
                        .size(48.dp)
                ) {
                    if (isSendingMessage) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(20.dp),
                            color = Color.White,
                            strokeWidth = 2.dp
                        )
                    } else {
                        Icon(
                            Icons.AutoMirrored.Filled.Send,
                            contentDescription = "Send message",
                            tint = Color.White
                        )
                    }
                }
            }
        }

        // Featured products overlay
        if (showProducts) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.5f))
                    .clickable { showProducts = false }
            ) {
                Card(
                    modifier = Modifier
                        .align(Alignment.BottomCenter)
                        .fillMaxWidth()
                        .fillMaxHeight(0.6f)
                        .clickable { }, // Prevent closing when clicking inside
                    shape = RoundedCornerShape(topStart = 16.dp, topEnd = 16.dp),
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Column(
                        modifier = Modifier.fillMaxSize()
                    ) {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(TchatSpacing.md),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(
                                "Featured Products",
                                style = MaterialTheme.typography.titleLarge,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.onSurface
                            )
                            IconButton(onClick = { showProducts = false }) {
                                Icon(Icons.Default.Close, contentDescription = "Close")
                            }
                        }

                        LazyColumn(
                            modifier = Modifier.weight(1f),
                            contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                            verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            items(streamData.featuredProducts) { product ->
                                StreamProductItem(
                                    product = product,
                                    onProductClick = { onProductClick(product.id) }
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun ChatMessageItem(
    message: ChatMessage,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        verticalAlignment = Alignment.Top
    ) {
        Text(
            message.username,
            style = MaterialTheme.typography.labelSmall,
            fontWeight = if (message.isFromHost) FontWeight.Bold else FontWeight.Medium,
            color = if (message.isFromHost) TchatColors.primary else TchatColors.onSurfaceVariant,
            modifier = Modifier.width(80.dp)
        )

        Spacer(modifier = Modifier.width(TchatSpacing.xs))

        Text(
            message.message,
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurface,
            modifier = Modifier.weight(1f)
        )
    }
}

@Composable
private fun StreamProductItem(
    product: ProductInStream,
    onProductClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onProductClick() },
        colors = CardDefaults.cardColors(containerColor = TchatColors.background),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier.padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Product image placeholder
            Box(
                modifier = Modifier
                    .size(60.dp)
                    .background(TchatColors.surfaceVariant, RoundedCornerShape(8.dp)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.ShoppingBag,
                    contentDescription = null,
                    tint = TchatColors.onSurfaceVariant,
                    modifier = Modifier.size(24.dp)
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            Column(modifier = Modifier.weight(1f)) {
                Text(
                    product.name,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(4.dp))

                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        product.price,
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.primary
                    )

                    if (product.originalPrice != null) {
                        Spacer(modifier = Modifier.width(TchatSpacing.xs))
                        Text(
                            product.originalPrice,
                            style = MaterialTheme.typography.bodySmall.copy(
                                textDecoration = androidx.compose.ui.text.style.TextDecoration.LineThrough
                            ),
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    if (product.discount != null) {
                        Spacer(modifier = Modifier.width(TchatSpacing.xs))
                        Box(
                            modifier = Modifier
                                .background(TchatColors.error, RoundedCornerShape(4.dp))
                                .padding(horizontal = 4.dp, vertical = 2.dp)
                        ) {
                            Text(
                                product.discount,
                                style = MaterialTheme.typography.labelSmall,
                                color = Color.White,
                                fontWeight = FontWeight.Bold
                            )
                        }
                    }
                }

                if (!product.inStock) {
                    Text(
                        "Out of Stock",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.error,
                        fontWeight = FontWeight.Medium
                    )
                }
            }

            Button(
                onClick = onProductClick,
                colors = ButtonDefaults.buttonColors(containerColor = TchatColors.primary),
                modifier = Modifier.height(32.dp),
                enabled = product.inStock
            ) {
                Text(
                    if (product.inStock) "Buy" else "Sold Out",
                    style = MaterialTheme.typography.labelSmall,
                    color = Color.White
                )
            }
        }
    }
}

// Mock data functions
private fun getLiveStreamById(streamId: String): LiveStreamData? {
    return getDummyLiveStreams().find { it.id == streamId }
}

private fun getDummyLiveStreams(): List<LiveStreamData> {
    return listOf(
        LiveStreamData(
            id = "stream1",
            title = "🔥 Flash Sale Alert! Premium Thai Skincare up to 70% OFF",
            hostName = "BeautyGuru_Thai",
            hostAvatar = "https://cdn.tchat.com/hosts/beautyguru_thai.jpg",
            viewerCount = 1284,
            isLive = true,
            thumbnailUrl = "https://cdn.tchat.com/streams/stream1_thumbnail.jpg",
            category = "Beauty & Health",
            shopId = "shop1",
            shopName = "Beauty World Thailand",
            shopAvatar = "https://cdn.tchat.com/shops/beauty_world_thai.jpg",
            shopRating = 4.8,
            description = "Join us for an exclusive flash sale featuring premium Thai skincare products. Get up to 70% off on bestselling serums, moisturizers, and masks. Limited time offer!",
            startedAt = "2024-01-01T14:00:00Z",
            tags = listOf("flash-sale", "skincare", "thai-beauty", "anti-aging", "premium"),
            featuredProducts = listOf(
                ProductInStream(
                    id = "product1",
                    name = "Snail Secretion Filtrate Serum - Anti-Aging",
                    description = "Premium anti-aging serum with 96% snail secretion filtrate. Helps repair damaged skin, reduce fine lines, and improve skin elasticity.",
                    price = "฿299",
                    originalPrice = "฿899",
                    discount = "-67%",
                    imageUrl = "https://cdn.tchat.com/products/snail_serum_main.jpg",
                    thumbnailUrl = "https://cdn.tchat.com/products/snail_serum_thumb.jpg",
                    inStock = true,
                    stockQuantity = 45,
                    rating = 4.6,
                    reviewCount = 234,
                    category = "Skincare Serums",
                    tags = listOf("anti-aging", "snail", "repair", "premium"),
                    specifications = mapOf(
                        "Volume" to "30ml",
                        "Key Ingredient" to "96% Snail Secretion Filtrate",
                        "Skin Type" to "All skin types",
                        "Origin" to "South Korea"
                    )
                ),
                ProductInStream(
                    id = "product2",
                    name = "Hyaluronic Acid Moisturizer",
                    description = "Deep hydrating moisturizer with multiple types of hyaluronic acid for long-lasting moisture and plumping effect.",
                    price = "฿450",
                    originalPrice = "฿750",
                    discount = "-40%",
                    imageUrl = "https://cdn.tchat.com/products/hyaluronic_moisturizer_main.jpg",
                    thumbnailUrl = "https://cdn.tchat.com/products/hyaluronic_moisturizer_thumb.jpg",
                    inStock = true,
                    stockQuantity = 23,
                    rating = 4.4,
                    reviewCount = 156,
                    category = "Moisturizers",
                    tags = listOf("hydrating", "hyaluronic-acid", "plumping", "daily-use"),
                    specifications = mapOf(
                        "Volume" to "50ml",
                        "Key Ingredient" to "5 Types of Hyaluronic Acid",
                        "Skin Type" to "Dry to normal skin",
                        "Usage" to "AM/PM"
                    )
                ),
                ProductInStream(
                    id = "product3",
                    name = "Vitamin C Brightening Mask",
                    description = "Intensive brightening mask with 20% Vitamin C and niacinamide to reduce dark spots and even skin tone.",
                    price = "฿199",
                    originalPrice = "฿399",
                    discount = "-50%",
                    imageUrl = "https://cdn.tchat.com/products/vitamin_c_mask_main.jpg",
                    thumbnailUrl = "https://cdn.tchat.com/products/vitamin_c_mask_thumb.jpg",
                    inStock = false,
                    stockQuantity = 0,
                    rating = 4.2,
                    reviewCount = 89,
                    category = "Face Masks",
                    tags = listOf("vitamin-c", "brightening", "dark-spots", "weekly-treatment"),
                    specifications = mapOf(
                        "Quantity" to "10 sheets",
                        "Key Ingredient" to "20% Vitamin C + Niacinamide",
                        "Usage" to "2-3 times per week",
                        "Treatment Time" to "15-20 minutes"
                    )
                )
            )
        ),
        LiveStreamData(
            id = "stream2",
            title = "Street Food Cooking Demo - Pad Thai Masterclass",
            hostName = "ChefNong",
            hostAvatar = "https://cdn.tchat.com/hosts/chef_nong.jpg",
            viewerCount = 856,
            isLive = true,
            thumbnailUrl = "https://cdn.tchat.com/streams/stream2_thumbnail.jpg",
            category = "Food & Cooking",
            shopId = "shop3",
            shopName = "Thai Kitchen Essentials",
            shopAvatar = "https://cdn.tchat.com/shops/thai_kitchen_essentials.jpg",
            shopRating = 4.6,
            description = "Join Chef Nong as he demonstrates the authentic art of Pad Thai cooking. Learn professional techniques and get exclusive deals on premium Thai cooking ingredients.",
            startedAt = "2024-01-01T15:30:00Z",
            tags = listOf("cooking", "thai-food", "pad-thai", "street-food", "ingredients"),
            featuredProducts = listOf(
                ProductInStream(
                    id = "product4",
                    name = "Traditional Pad Thai Sauce Set",
                    description = "Complete authentic Pad Thai sauce set with tamarind paste, fish sauce, and palm sugar. Made from traditional Thai recipes.",
                    price = "฿180",
                    originalPrice = "฿220",
                    discount = "-18%",
                    imageUrl = "https://cdn.tchat.com/products/padthai_sauce_set_main.jpg",
                    thumbnailUrl = "https://cdn.tchat.com/products/padthai_sauce_set_thumb.jpg",
                    inStock = true,
                    stockQuantity = 67,
                    rating = 4.5,
                    reviewCount = 89,
                    category = "Cooking Sauces",
                    tags = listOf("pad-thai", "authentic", "traditional", "thai-cuisine"),
                    specifications = mapOf(
                        "Package Contents" to "Tamarind paste, Fish sauce, Palm sugar",
                        "Serves" to "4-6 portions",
                        "Origin" to "Thailand",
                        "Shelf Life" to "12 months"
                    )
                )
            )
        ),
        LiveStreamData(
            id = "stream3",
            title = "Gaming Setup Review - Latest RGB Mechanical Keyboards",
            hostName = "TechReviewer_TH",
            hostAvatar = "https://cdn.tchat.com/hosts/tech_reviewer_th.jpg",
            viewerCount = 2156,
            isLive = true,
            thumbnailUrl = "https://cdn.tchat.com/streams/stream3_thumbnail.jpg",
            category = "Electronics",
            shopId = "shop2",
            shopName = "TechZone Bangkok",
            shopAvatar = "https://cdn.tchat.com/shops/techzone_bangkok.jpg",
            shopRating = 4.9,
            description = "Comprehensive review of the latest RGB mechanical keyboards. Get hands-on demonstrations and exclusive viewer discounts on gaming accessories.",
            startedAt = "2024-01-01T16:00:00Z",
            tags = listOf("gaming", "keyboards", "rgb", "mechanical", "tech-review"),
            featuredProducts = listOf(
                ProductInStream(
                    id = "product5",
                    name = "RGB Mechanical Keyboard - Blue Switches",
                    description = "Premium mechanical keyboard with Cherry MX Blue switches, customizable RGB lighting, and aluminum frame. Perfect for gaming and typing.",
                    price = "฿2,890",
                    originalPrice = "฿3,590",
                    discount = "-20%",
                    imageUrl = "https://cdn.tchat.com/products/rgb_keyboard_main.jpg",
                    thumbnailUrl = "https://cdn.tchat.com/products/rgb_keyboard_thumb.jpg",
                    inStock = true,
                    stockQuantity = 12,
                    rating = 4.7,
                    reviewCount = 178,
                    category = "Gaming Keyboards",
                    tags = listOf("mechanical", "rgb", "gaming", "cherry-mx", "premium"),
                    specifications = mapOf(
                        "Switch Type" to "Cherry MX Blue",
                        "Layout" to "Full Size (104 keys)",
                        "Backlight" to "RGB with 16.8M colors",
                        "Connection" to "USB-C with braided cable",
                        "Frame Material" to "Aluminum alloy",
                        "Key Rollover" to "N-Key Rollover"
                    )
                )
            )
        )
    )
}

private fun getDummyChatMessages(): List<ChatMessage> {
    return listOf(
        ChatMessage("1", "BeautyGuru_Thai", "Welcome everyone! Today's flash sale is incredible!", "12:30", true),
        ChatMessage("2", "user123", "How long does shipping take?", "12:29"),
        ChatMessage("3", "shopaholic", "The serum looks amazing! 😍", "12:28"),
        ChatMessage("4", "BeautyGuru_Thai", "Shipping is 1-2 days within Bangkok!", "12:28", true),
        ChatMessage("5", "beautylover", "Just ordered 3 items!", "12:27"),
        ChatMessage("6", "newbie", "Is this suitable for sensitive skin?", "12:27"),
        ChatMessage("7", "BeautyGuru_Thai", "Yes! It's dermatologist tested", "12:26", true),
        ChatMessage("8", "user456", "Price is unbeatable! 🔥", "12:26"),
        ChatMessage("9", "skincare_addict", "Adding to cart now!", "12:25"),
        ChatMessage("10", "customer789", "Do you have international shipping?", "12:24")
    )
}

// Reuse existing shop data function
private fun getDummyShops(): List<ShopItem> {
    return listOf(
        ShopItem(
            id = "shop1",
            name = "Beauty World Thailand",
            description = "Premium beauty and skincare products",
            avatar = "",
            coverImage = "",
            rating = 4.8,
            deliveryTime = "1-2 days",
            distance = "5.2 km",
            isVerified = true,
            followers = 12500,
            totalProducts = 234
        ),
        ShopItem(
            id = "shop2",
            name = "Thai Kitchen Essentials",
            description = "Authentic Thai cooking ingredients and tools",
            avatar = "",
            coverImage = "",
            rating = 4.6,
            deliveryTime = "30-60 mins",
            distance = "2.8 km",
            isVerified = true,
            followers = 8900,
            totalProducts = 156
        ),
        ShopItem(
            id = "shop3",
            name = "TechZone Bangkok",
            description = "Latest electronics and gaming accessories",
            avatar = "",
            coverImage = "",
            rating = 4.9,
            deliveryTime = "1-3 days",
            distance = "7.1 km",
            isVerified = true,
            followers = 25600,
            totalProducts = 890
        )
    )
}