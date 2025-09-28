package com.tchat.mobile.social.presentation

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.tchat.mobile.social.domain.models.SocialPost
import com.tchat.mobile.social.domain.models.DiscoveryProfile
import com.tchat.mobile.social.domain.models.SEARegions
import com.tchat.mobile.social.domain.models.SocialLocalization
import com.tchat.mobile.social.services.RegionalContentService
import com.tchat.mobile.social.services.TrendingTopic
import com.tchat.mobile.social.services.CulturalEvent
import io.kamel.image.KamelImage
import io.kamel.image.asyncPainterResource
import kotlinx.coroutines.launch
import org.koin.compose.koinInject

/**
 * KMP Social Screen
 *
 * Cross-platform social feed with:
 * - Offline-first architecture
 * - Southeast Asian regional features
 * - Mobile-optimized performance
 * - Real-time sync capabilities
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SocialScreen(
    modifier: Modifier = Modifier,
    viewModel: SocialViewModel = koinInject(),
    regionalContentService: RegionalContentService = koinInject()
) {
    val uiState by viewModel.uiState.collectAsState()
    val feedState by viewModel.feedState.collectAsState()
    val profileState by viewModel.profileState.collectAsState()
    val discoveryState by viewModel.discoveryState.collectAsState()
    val syncState by viewModel.syncState.collectAsState()

    val scope = rememberCoroutineScope()
    val listState = rememberLazyListState()

    // Regional content state
    var trendingTopics by remember { mutableStateOf<List<TrendingTopic>>(emptyList()) }
    var culturalEvents by remember { mutableStateOf<List<CulturalEvent>>(emptyList()) }
    var regionalContentLoading by remember { mutableStateOf(false) }

    // Pull to refresh state
    var isRefreshing by remember { mutableStateOf(false) }

    LaunchedEffect(feedState.isRefreshing) {
        isRefreshing = feedState.isRefreshing
    }

    // Load regional content when region changes
    LaunchedEffect(uiState.currentRegion) {
        regionalContentLoading = true
        try {
            regionalContentService.getTrendingTopics(uiState.currentRegion)
                .onSuccess { topics -> trendingTopics = topics }
            regionalContentService.getCulturalEvents(uiState.currentRegion)
                .onSuccess { events -> culturalEvents = events }
        } catch (e: Exception) {
            // Handle error silently for now
        } finally {
            regionalContentLoading = false
        }
    }

    // Show sync indicator
    LaunchedEffect(syncState) {
        if (syncState.name == "SYNCING") {
            // Auto-sync completed
        }
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
    ) {
        // Top App Bar with sync status
        TopAppBar(
            title = {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        text = "Social",
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.Bold
                    )

                    // Sync status indicator
                    when (syncState.name) {
                        "SYNCING" -> {
                            CircularProgressIndicator(
                                modifier = Modifier.size(16.dp),
                                strokeWidth = 2.dp
                            )
                        }
                        "ERROR" -> {
                            Icon(
                                imageVector = Icons.Default.Warning,
                                contentDescription = "Sync Error",
                                tint = MaterialTheme.colorScheme.error,
                                modifier = Modifier.size(16.dp)
                            )
                        }
                    }
                }
            },
            actions = {
                // Region selector
                RegionSelector(
                    currentRegion = uiState.currentRegion,
                    onRegionChanged = viewModel::changeRegion
                )

                // Manual sync button
                IconButton(
                    onClick = { viewModel.performSync() },
                    enabled = syncState.name != "SYNCING"
                ) {
                    Icon(
                        imageVector = Icons.Default.Sync,
                        contentDescription = "Sync"
                    )
                }
            }
        )

        // Error/Message snackbar
        uiState.error?.let { error ->
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.errorContainer
                )
            ) {
                Row(
                    modifier = Modifier.padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Icon(
                        imageVector = Icons.Default.Error,
                        contentDescription = null,
                        tint = MaterialTheme.colorScheme.onErrorContainer
                    )
                    Text(
                        text = error,
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onErrorContainer,
                        modifier = Modifier.weight(1f)
                    )
                    IconButton(onClick = viewModel::clearError) {
                        Icon(
                            imageVector = Icons.Default.Close,
                            contentDescription = "Dismiss",
                            tint = MaterialTheme.colorScheme.onErrorContainer
                        )
                    }
                }
            }
        }

        uiState.message?.let { message ->
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.primaryContainer
                )
            ) {
                Row(
                    modifier = Modifier.padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Icon(
                        imageVector = Icons.Default.CheckCircle,
                        contentDescription = null,
                        tint = MaterialTheme.colorScheme.onPrimaryContainer
                    )
                    Text(
                        text = message,
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onPrimaryContainer,
                        modifier = Modifier.weight(1f)
                    )
                    IconButton(onClick = viewModel::clearMessage) {
                        Icon(
                            imageVector = Icons.Default.Close,
                            contentDescription = "Dismiss",
                            tint = MaterialTheme.colorScheme.onPrimaryContainer
                        )
                    }
                }
            }
        }

        // Main content
        if (feedState.isLoading && feedState.homeFeed == null) {
            // Initial loading state
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.spacedBy(16.dp)
                ) {
                    CircularProgressIndicator()
                    Text(
                        text = "Loading social feed...",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }
        } else {
            // Social feed with pull to refresh
            LazyColumn(
                state = listState,
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(bottom = 16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                // Regional trending section
                if (trendingTopics.isNotEmpty()) {
                    item {
                        RegionalTrendingSection(
                            topics = trendingTopics,
                            regionCode = uiState.currentRegion,
                            modifier = Modifier.padding(horizontal = 16.dp)
                        )
                    }
                }

                // Cultural events section
                if (culturalEvents.isNotEmpty()) {
                    item {
                        CulturalEventsSection(
                            events = culturalEvents,
                            regionCode = uiState.currentRegion,
                            modifier = Modifier.padding(horizontal = 16.dp)
                        )
                    }
                }

                // Discovery section
                if (discoveryState.discoveryProfiles.isNotEmpty()) {
                    item {
                        DiscoverySection(
                            profiles = discoveryState.discoveryProfiles,
                            onFollowUser = viewModel::followUser,
                            modifier = Modifier.padding(horizontal = 16.dp)
                        )
                    }
                }

                // Create post section
                item {
                    CreatePostCard(
                        onCreatePost = viewModel::createPost,
                        isPosting = uiState.isPosting,
                        currentRegion = uiState.currentRegion,
                        modifier = Modifier.padding(horizontal = 16.dp)
                    )
                }

                // Feed posts
                feedState.homeFeed?.posts?.let { posts ->
                    if (posts.isEmpty()) {
                        item {
                            EmptyFeedCard(
                                onRefresh = { viewModel.refreshFeed() },
                                modifier = Modifier.padding(horizontal = 16.dp)
                            )
                        }
                    } else {
                        items(
                            items = posts,
                            key = { it.id }
                        ) { post ->
                            PostCard(
                                post = post,
                                onLike = { viewModel.likePost(post.id) },
                                onBookmark = { viewModel.bookmarkPost(post.id) },
                                onComment = { /* TODO: Navigate to comments */ },
                                onShare = { /* TODO: Share post */ },
                                modifier = Modifier.padding(horizontal = 16.dp)
                            )
                        }
                    }
                }

                // Loading indicator for pagination
                if (feedState.isLoading && feedState.homeFeed?.posts?.isNotEmpty() == true) {
                    item {
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(16.dp),
                            contentAlignment = Alignment.Center
                        ) {
                            CircularProgressIndicator()
                        }
                    }
                }
            }
        }
    }

    // Pull to refresh
    LaunchedEffect(isRefreshing) {
        if (isRefreshing) {
            scope.launch {
                listState.animateScrollToItem(0)
            }
        }
    }
}

@Composable
fun RegionSelector(
    currentRegion: String,
    onRegionChanged: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    var expanded by remember { mutableStateOf(false) }

    val regions = mapOf(
        "TH" to "🇹🇭 Thailand (ไทย)",
        "SG" to "🇸🇬 Singapore",
        "ID" to "🇮🇩 Indonesia",
        "MY" to "🇲🇾 Malaysia",
        "PH" to "🇵🇭 Philippines",
        "VN" to "🇻🇳 Vietnam (Việt Nam)"
    )

    Box(modifier = modifier) {
        IconButton(onClick = { expanded = true }) {
            Text(
                text = regions[currentRegion]?.take(2) ?: "🌍",
                style = MaterialTheme.typography.titleMedium
            )
        }

        DropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false }
        ) {
            regions.forEach { (code, name) ->
                DropdownMenuItem(
                    text = { Text(name) },
                    onClick = {
                        onRegionChanged(code)
                        expanded = false
                    },
                    leadingIcon = if (code == currentRegion) {
                        { Icon(Icons.Default.Check, contentDescription = null) }
                    } else null
                )
            }
        }
    }
}

@Composable
fun DiscoverySection(
    profiles: List<DiscoveryProfile>,
    onFollowUser: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            Text(
                text = "Discover People",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold
            )

            profiles.take(3).forEach { discovery ->
                DiscoveryProfileItem(
                    discovery = discovery,
                    onFollow = { onFollowUser(discovery.profile.id) }
                )
            }
        }
    }
}

@Composable
fun DiscoveryProfileItem(
    discovery: DiscoveryProfile,
    onFollow: () -> Unit,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Avatar
        Box(
            modifier = Modifier
                .size(40.dp)
                .clip(CircleShape)
                .background(MaterialTheme.colorScheme.primary)
        ) {
            discovery.profile.avatar?.let { avatarUrl ->
                KamelImage(
                    resource = asyncPainterResource(avatarUrl),
                    contentDescription = "Avatar",
                    modifier = Modifier.fillMaxSize(),
                    contentScale = ContentScale.Crop
                )
            } ?: run {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = discovery.profile.displayName?.take(1)?.uppercase()
                            ?: discovery.profile.username.take(1).uppercase(),
                        style = MaterialTheme.typography.titleMedium,
                        color = MaterialTheme.colorScheme.onPrimary
                    )
                }
            }
        }

        // Profile info
        Column(
            modifier = Modifier.weight(1f),
            verticalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            Text(
                text = discovery.profile.displayName ?: discovery.profile.username,
                style = MaterialTheme.typography.bodyMedium,
                fontWeight = FontWeight.Medium
            )
            Text(
                text = discovery.discoveryReason.replace("_", " ").capitalize(),
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
        }

        // Follow button
        OutlinedButton(
            onClick = onFollow,
            modifier = Modifier.height(36.dp)
        ) {
            Text("Follow")
        }
    }
}

@Composable
fun CreatePostCard(
    onCreatePost: (String, String, List<String>) -> Unit,
    isPosting: Boolean,
    currentRegion: String = "TH",
    modifier: Modifier = Modifier
) {
    var postText by remember { mutableStateOf("") }
    var expanded by remember { mutableStateOf(false) }

    Card(
        modifier = modifier.fillMaxWidth(),
        onClick = { expanded = true }
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            if (expanded) {
                OutlinedTextField(
                    value = postText,
                    onValueChange = { postText = it },
                    placeholder = { Text(SEARegions.getRegionalGreeting(currentRegion)) },
                    modifier = Modifier.fillMaxWidth(),
                    minLines = 3,
                    maxLines = 6
                )

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        IconButton(onClick = { /* TODO: Add media */ }) {
                            Icon(Icons.Default.Image, contentDescription = "Add Image")
                        }
                        IconButton(onClick = { /* TODO: Add emoji */ }) {
                            Icon(Icons.Default.EmojiEmotions, contentDescription = "Add Emoji")
                        }
                    }

                    Row(
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        TextButton(
                            onClick = {
                                expanded = false
                                postText = ""
                            }
                        ) {
                            Text("Cancel")
                        }

                        Button(
                            onClick = {
                                if (postText.isNotBlank()) {
                                    onCreatePost(postText, "text", emptyList())
                                    postText = ""
                                    expanded = false
                                }
                            },
                            enabled = postText.isNotBlank() && !isPosting
                        ) {
                            if (isPosting) {
                                CircularProgressIndicator(
                                    modifier = Modifier.size(16.dp),
                                    strokeWidth = 2.dp
                                )
                            } else {
                                Text("Post")
                            }
                        }
                    }
                }
            } else {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    Icon(
                        imageVector = Icons.Default.Edit,
                        contentDescription = null,
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    Text(
                        text = "Share what's happening...",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }
        }
    }
}

@Composable
fun PostCard(
    post: SocialPost,
    onLike: () -> Unit,
    onBookmark: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Post header
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                // Author avatar
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .clip(CircleShape)
                        .background(MaterialTheme.colorScheme.primary)
                ) {
                    post.authorAvatar?.let { avatarUrl ->
                        KamelImage(
                            resource = asyncPainterResource(avatarUrl),
                            contentDescription = "Avatar",
                            modifier = Modifier.fillMaxSize(),
                            contentScale = ContentScale.Crop
                        )
                    } ?: run {
                        Box(
                            modifier = Modifier.fillMaxSize(),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = post.authorDisplayName?.take(1)?.uppercase()
                                    ?: post.authorUsername.take(1).uppercase(),
                                style = MaterialTheme.typography.titleMedium,
                                color = MaterialTheme.colorScheme.onPrimary
                            )
                        }
                    }
                }

                // Author info
                Column(
                    modifier = Modifier.weight(1f)
                ) {
                    Text(
                        text = post.authorDisplayName ?: post.authorUsername,
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium
                    )
                    Text(
                        text = formatTimeAgo(post.createdAt),
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }

                // Offline indicator
                if (post.isOfflineEdit) {
                    Icon(
                        imageVector = Icons.Default.CloudOff,
                        contentDescription = "Offline",
                        tint = MaterialTheme.colorScheme.onSurfaceVariant,
                        modifier = Modifier.size(16.dp)
                    )
                }
            }

            // Post content
            Text(
                text = post.content,
                style = MaterialTheme.typography.bodyMedium
            )

            // Post media
            if (post.mediaUrls.isNotEmpty()) {
                post.mediaUrls.first().let { mediaUrl ->
                    KamelImage(
                        resource = asyncPainterResource(mediaUrl),
                        contentDescription = "Post media",
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(200.dp)
                            .clip(RoundedCornerShape(8.dp)),
                        contentScale = ContentScale.Crop
                    )
                }
            }

            // Engagement stats
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Text(
                    text = "${post.likesCount} likes",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Text(
                    text = "${post.commentsCount} comments",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Text(
                    text = "${post.sharesCount} shares",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }

            // Action buttons
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                IconButton(onClick = onLike) {
                    Icon(
                        imageVector = if (post.isLikedByUser) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                        contentDescription = "Like",
                        tint = if (post.isLikedByUser) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.onSurface
                    )
                }

                IconButton(onClick = onComment) {
                    Icon(
                        imageVector = Icons.Default.Comment,
                        contentDescription = "Comment"
                    )
                }

                IconButton(onClick = onShare) {
                    Icon(
                        imageVector = Icons.Default.Share,
                        contentDescription = "Share"
                    )
                }

                IconButton(onClick = onBookmark) {
                    Icon(
                        imageVector = if (post.isBookmarkedByUser) Icons.Default.Bookmark else Icons.Default.BookmarkBorder,
                        contentDescription = "Bookmark",
                        tint = if (post.isBookmarkedByUser) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurface
                    )
                }
            }
        }
    }
}

@Composable
fun EmptyFeedCard(
    onRefresh: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            Icon(
                imageVector = Icons.Default.People,
                contentDescription = null,
                modifier = Modifier.size(48.dp),
                tint = MaterialTheme.colorScheme.onSurfaceVariant
            )

            Text(
                text = "Your feed is empty",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Medium
            )

            Text(
                text = "Follow people to see their posts here, or create your first post!",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )

            Button(onClick = onRefresh) {
                Text("Refresh")
            }
        }
    }
}

// Helper functions
@Composable
private fun getCurrentRegionEmoji(): String {
    // Dynamic region emoji based on current context
    return "🌏" // Southeast Asia specific emoji
}

@Composable
private fun getRegionEmoji(regionCode: String): String {
    return when (regionCode) {
        "TH" -> "🇹🇭"
        "SG" -> "🇸🇬"
        "ID" -> "🇮🇩"
        "MY" -> "🇲🇾"
        "PH" -> "🇵🇭"
        "VN" -> "🇻🇳"
        else -> "🌏"
    }
}

@Composable
private fun getRegionalGreeting(regionCode: String): String {
    return when (regionCode) {
        "TH" -> "สวัสดี! What's happening?"
        "SG" -> "Lah! What's happening?"
        "ID" -> "Halo! Apa kabar?"
        "MY" -> "Apa khabar! What's happening?"
        "PH" -> "Kumusta! What's happening?"
        "VN" -> "Xin chào! What's happening?"
        else -> "Hello! What's happening?"
    }
}

@Composable
private fun getRegionalHashtags(regionCode: String): List<String> {
    return when (regionCode) {
        "TH" -> listOf("#Thailand", "#Bangkok", "#Krungthep", "#ThaiCulture", "#Songkran", "#TomYum", "#BTS", "#MRT")
        "SG" -> listOf("#Singapore", "#Merlion", "#HawkerCentre", "#SingaporeLife", "#MRT", "#Singlish", "#GardensByTheBay")
        "ID" -> listOf("#Indonesia", "#Jakarta", "#Bali", "#IndonesianCulture", "#RendangLife", "#Batik", "#Wonderful")
        "MY" -> listOf("#Malaysia", "#KualaLumpur", "#MalaysianFood", "#TrulyAsia", "#Mamak", "#Durian", "#KLCC")
        "PH" -> listOf("#Philippines", "#Manila", "#Pinoy", "#Adobo", "#IslandLife", "#Jeepney", "#Pinoy")
        "VN" -> listOf("#Vietnam", "#Hanoi", "#HoChiMinh", "#Pho", "#Vietnamese", "#Motorbike", "#Banh")
        else -> listOf("#SoutheastAsia", "#ASEAN")
    }
}

private fun formatTimeAgo(timestamp: String): String {
    // Simplified time formatting - in real implementation use proper date formatting
    return "2m ago"
}

private fun String.capitalize(): String {
    return this.replaceFirstChar { if (it.isLowerCase()) it.titlecase() else it.toString() }
}

// Southeast Asian Regional Content Components

@Composable
fun RegionalTrendingSection(
    topics: List<TrendingTopic>,
    regionCode: String,
    modifier: Modifier = Modifier
) {
    val region = SEARegions.getRegion(regionCode)

    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = region?.flag ?: "🌍",
                    style = MaterialTheme.typography.titleLarge
                )
                Text(
                    text = SocialLocalization.getString("trending_now"),
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
                Text(
                    text = "in ${region?.name ?: "Southeast Asia"}",
                    style = MaterialTheme.typography.titleMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }

            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items(topics.take(6)) { topic ->
                    TrendingTopicChip(
                        topic = topic,
                        onClick = { /* Handle trending topic click */ }
                    )
                }
            }
        }
    }
}

@Composable
fun TrendingTopicChip(
    topic: TrendingTopic,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    AssistChip(
        onClick = onClick,
        label = {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                Text(topic.emoji)
                Text(
                    text = topic.hashtag,
                    style = MaterialTheme.typography.bodySmall
                )
                Text(
                    text = "${topic.postCount / 1000}K",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        },
        modifier = modifier,
        colors = AssistChipDefaults.assistChipColors(
            containerColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f),
            labelColor = MaterialTheme.colorScheme.primary
        )
    )
}

@Composable
fun CulturalEventsSection(
    events: List<CulturalEvent>,
    regionCode: String,
    modifier: Modifier = Modifier
) {
    val region = SEARegions.getRegion(regionCode)

    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.primaryContainer
        )
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Icon(
                    imageVector = Icons.Default.Event,
                    contentDescription = null,
                    tint = MaterialTheme.colorScheme.onPrimaryContainer
                )
                Text(
                    text = SocialLocalization.getString("cultural_events"),
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.onPrimaryContainer
                )
                Text(
                    text = region?.flag ?: "🌍",
                    style = MaterialTheme.typography.titleMedium
                )
            }

            events.take(2).forEach { event ->
                CulturalEventCard(
                    event = event,
                    onClick = { /* Handle event click */ }
                )
            }
        }
    }
}

@Composable
fun CulturalEventCard(
    event: CulturalEvent,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onClick() },
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surface
        )
    ) {
        Column(
            modifier = Modifier.padding(12.dp),
            verticalArrangement = Arrangement.spacedBy(6.dp)
        ) {
            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalAlignment = Alignment.Top
            ) {
                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = event.name,
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium
                    )
                    if (event.localName != event.name) {
                        Text(
                            text = event.localName,
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                    Text(
                        text = "${event.date} • ${event.location}",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }

                Badge(
                    containerColor = MaterialTheme.colorScheme.secondaryContainer
                ) {
                    Text(
                        text = event.category,
                        style = MaterialTheme.typography.labelSmall
                    )
                }
            }

            // Event hashtags
            if (event.hashtags.isNotEmpty()) {
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    items(event.hashtags.take(3)) { hashtag ->
                        AssistChip(
                            onClick = { /* Handle hashtag click */ },
                            label = {
                                Text(
                                    text = hashtag,
                                    style = MaterialTheme.typography.labelSmall
                                )
                            },
                            modifier = Modifier.height(24.dp)
                        )
                    }
                }
            }
        }
    }
}