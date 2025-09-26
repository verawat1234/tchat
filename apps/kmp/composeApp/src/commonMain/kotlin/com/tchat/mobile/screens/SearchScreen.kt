package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
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
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

data class SearchResult(
    val id: String,
    val type: SearchResultType,
    val title: String,
    val subtitle: String? = null,
    val timestamp: String? = null,
    val highlight: String? = null
)

enum class SearchResultType {
    CHAT, CONTACT, MERCHANT, PRODUCT, MESSAGE, HASHTAG
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SearchScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var searchQuery by remember { mutableStateOf("") }
    var selectedTab by remember { mutableStateOf(0) }

    // Mock data following web design
    val recentSearches = listOf(
        "Pad Thai",
        "Som Tam",
        "PromptPay",
        "Bangkok Street Food",
        "Family Group"
    )

    val trendingSearches = listOf(
        "Songkran Festival" to "1.2k",
        "PromptPay QR" to "890",
        "Bangkok Street Food" to "756",
        "Thai New Year" to "645",
        "Som Tam Recipe" to "532"
    )

    val popularHashtags = listOf(
        "#ThaiStreetFood", "#BangkokEats", "#SongkranFestival",
        "#PromptPay", "#ThaiCulture", "#SEAFood",
        "#LiveCooking", "#ThaiMarket"
    )

    val searchResults = listOf(
        SearchResult("1", SearchResultType.CHAT, "Family Group", "Mom: Dinner at 7pm! 🍽️", "5 min ago"),
        SearchResult("2", SearchResultType.MERCHANT, "Somtam Vendor", "Thai Street Food • 0.5 km away"),
        SearchResult("3", SearchResultType.PRODUCT, "Pad Thai Goong", "฿45 • Bangkok Street Food"),
        SearchResult("4", SearchResultType.MESSAGE, "AI Assistant", "Welcome to Telegram SEA! Let me help you get started", "10:30 AM", "SEA"),
        SearchResult("5", SearchResultType.HASHTAG, "#ThaiStreetFood", "234 posts • Trending in Thailand")
    )

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Search", fontWeight = FontWeight.Bold) },
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
                    IconButton(onClick = { /* Filter */ }) {
                        Icon(
                            Icons.Default.FilterList,
                            contentDescription = "Filter",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = modifier
                .fillMaxSize()
                .padding(paddingValues)
                .background(TchatColors.background)
        ) {
            // Search Input
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                TchatInput(
                    value = searchQuery,
                    onValueChange = { searchQuery = it },
                    type = TchatInputType.Text,
                    placeholder = "Search conversations, contacts, messages...",
                    leadingIcon = Icons.Default.Search,
                    modifier = Modifier.weight(1f)
                )
            }

            // Tabs for different content
            TabRow(
                selectedTabIndex = selectedTab,
                containerColor = TchatColors.surface,
                contentColor = TchatColors.primary,
                modifier = Modifier.fillMaxWidth()
            ) {
                Tab(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    text = { Text("All") }
                )
                Tab(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    text = { Text("Recent") }
                )
                Tab(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    text = { Text("Trending") }
                )
            }

            // Content based on selected tab and search query
            LazyColumn(
                modifier = Modifier.weight(1f),
                contentPadding = PaddingValues(TchatSpacing.md),
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                when {
                    searchQuery.isNotEmpty() -> {
                        // Show search results
                        items(searchResults.filter {
                            it.title.contains(searchQuery, ignoreCase = true) ||
                            it.subtitle?.contains(searchQuery, ignoreCase = true) == true
                        }) { result ->
                            SearchResultItem(result)
                        }
                    }
                    selectedTab == 0 -> {
                        // All - Show popular hashtags
                        item {
                            Text(
                                "Popular Hashtags",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.onSurface,
                                modifier = Modifier.padding(bottom = TchatSpacing.sm)
                            )
                        }
                        items(popularHashtags.chunked(2)) { hashtagPair ->
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                            ) {
                                hashtagPair.forEach { hashtag ->
                                    Card(
                                        modifier = Modifier
                                            .weight(1f)
                                            .clickable { searchQuery = hashtag },
                                        colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
                                    ) {
                                        Text(
                                            hashtag,
                                            modifier = Modifier.padding(TchatSpacing.md),
                                            color = TchatColors.primary,
                                            fontWeight = FontWeight.Medium
                                        )
                                    }
                                }
                                // Fill remaining space if odd number
                                if (hashtagPair.size == 1) {
                                    Spacer(modifier = Modifier.weight(1f))
                                }
                            }
                        }
                    }
                    selectedTab == 1 -> {
                        // Recent searches
                        item {
                            Text(
                                "Recent Searches",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.onSurface,
                                modifier = Modifier.padding(bottom = TchatSpacing.sm)
                            )
                        }
                        items(recentSearches) { search ->
                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable { searchQuery = search }
                                    .padding(vertical = TchatSpacing.sm),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    Icons.Default.History,
                                    contentDescription = null,
                                    tint = TchatColors.onSurfaceVariant,
                                    modifier = Modifier.size(20.dp)
                                )
                                Spacer(modifier = Modifier.width(TchatSpacing.md))
                                Text(
                                    search,
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = TchatColors.onSurface
                                )
                                Spacer(modifier = Modifier.weight(1f))
                                Icon(
                                    Icons.Default.TrendingUp,
                                    contentDescription = null,
                                    tint = TchatColors.onSurfaceVariant,
                                    modifier = Modifier.size(16.dp)
                                )
                            }
                        }
                    }
                    selectedTab == 2 -> {
                        // Trending searches
                        item {
                            Text(
                                "Trending Now",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold,
                                color = TchatColors.onSurface,
                                modifier = Modifier.padding(bottom = TchatSpacing.sm)
                            )
                        }
                        items(trendingSearches) { (search, count) ->
                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable { searchQuery = search }
                                    .padding(vertical = TchatSpacing.sm),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    Icons.Default.TrendingUp,
                                    contentDescription = null,
                                    tint = TchatColors.primary,
                                    modifier = Modifier.size(20.dp)
                                )
                                Spacer(modifier = Modifier.width(TchatSpacing.md))
                                Column(modifier = Modifier.weight(1f)) {
                                    Text(
                                        search,
                                        style = MaterialTheme.typography.bodyLarge,
                                        color = TchatColors.onSurface,
                                        fontWeight = FontWeight.Medium
                                    )
                                    Text(
                                        "$count searches",
                                        style = MaterialTheme.typography.bodySmall,
                                        color = TchatColors.onSurfaceVariant
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

@Composable
private fun SearchResultItem(
    result: SearchResult,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { /* Handle result click */ },
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Icon based on result type
            Box(
                modifier = Modifier
                    .size(40.dp)
                    .clip(CircleShape)
                    .background(
                        when (result.type) {
                            SearchResultType.CHAT -> TchatColors.primary
                            SearchResultType.MERCHANT -> TchatColors.success
                            SearchResultType.PRODUCT -> TchatColors.warning
                            SearchResultType.MESSAGE -> TchatColors.primaryLight
                            SearchResultType.HASHTAG -> TchatColors.surfaceVariant
                            else -> TchatColors.surfaceVariant
                        }
                    ),
                contentAlignment = Alignment.Center
            ) {
                val icon = when (result.type) {
                    SearchResultType.CHAT -> Icons.Default.Chat
                    SearchResultType.CONTACT -> Icons.Default.Person
                    SearchResultType.MERCHANT -> Icons.Default.Store
                    SearchResultType.PRODUCT -> Icons.Default.ShoppingCart
                    SearchResultType.MESSAGE -> Icons.Default.Message
                    SearchResultType.HASHTAG -> Icons.Default.Tag
                }
                Icon(
                    icon,
                    contentDescription = null,
                    tint = TchatColors.onPrimary,
                    modifier = Modifier.size(20.dp)
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Content
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    result.title,
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
                if (result.subtitle != null) {
                    Text(
                        result.subtitle,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }

            // Timestamp
            if (result.timestamp != null) {
                Text(
                    result.timestamp,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}