package com.tchat.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
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
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Social media feed interface screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SocialScreen() {
    var selectedTab by remember { mutableStateOf(SocialTab.FEED) }
    var showingCamera by remember { mutableStateOf(false) }

    // Mock posts
    val posts = listOf(
        Post("Alice Johnson", "2h", "Just finished my morning workout! 💪 #fitness", Icons.Default.DirectionsWalk, 23, 5),
        Post("Bob Smith", "4h", "Beautiful sunset from my balcony 🌅", Icons.Default.WbSunny, 67, 12),
        Post("Carol Davis", "6h", "New coffee shop opened downtown! ☕️ Must try", Icons.Default.Coffee, 45, 8),
        Post("David Wilson", "8h", "Working on a new project. Excited to share soon! 🚀", Icons.Default.Laptop, 89, 15),
        Post("Emma Brown", "12h", "Weekend hiking adventure! Nature is amazing 🏔️", Icons.Default.Terrain, 156, 28)
    )

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Colors.background)
    ) {
        // Top app bar
        TopAppBar(
            title = {
                Text(
                    text = "Social",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.textPrimary
                )
            },
            actions = {
                IconButton(onClick = { showingCamera = true }) {
                    Icon(
                        imageVector = Icons.Default.CameraAlt,
                        contentDescription = "Camera",
                        tint = Colors.primary
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = Colors.background
            )
        )

        // Tab selector
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .background(Colors.surface)
        ) {
            SocialTab.values().forEach { tab ->
                Column(
                    modifier = Modifier
                        .weight(1f)
                        .clickable { selectedTab = tab }
                        .padding(vertical = Spacing.sm),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Icon(
                        imageVector = tab.icon,
                        contentDescription = tab.title,
                        tint = if (selectedTab == tab) Colors.primary else Colors.textSecondary,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.height(Spacing.xs))
                    Text(
                        text = tab.title,
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Medium,
                        color = if (selectedTab == tab) Colors.primary else Colors.textSecondary
                    )
                }
            }
        }

        // Content based on selected tab
        when (selectedTab) {
            SocialTab.FEED -> FeedView(posts = posts)
            SocialTab.DISCOVER -> DiscoverView()
            SocialTab.NOTIFICATIONS -> NotificationsView()
        }
    }
}

// MARK: - Social Tab Enum
enum class SocialTab(
    val title: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
) {
    FEED("Feed", Icons.Default.Home),
    DISCOVER("Discover", Icons.Default.Explore),
    NOTIFICATIONS("Alerts", Icons.Default.Notifications)
}

// MARK: - Data Classes
data class Post(
    val author: String,
    val time: String,
    val content: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector,
    val likes: Int,
    val comments: Int
)

// MARK: - Feed View
@Composable
private fun FeedView(posts: List<Post>) {
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        items(posts) { post ->
            PostCard(post = post)
        }
    }
}

// MARK: - Post Card Component
@Composable
private fun PostCard(post: Post) {
    var isLiked by remember { mutableStateOf(false) }
    var currentLikes by remember { mutableStateOf(post.likes) }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = androidx.compose.ui.graphics.Color.White
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 4.dp
        )
    ) {
        Column(
            modifier = Modifier.padding(Spacing.md)
        ) {
            // Header
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .background(
                            color = Colors.primary.copy(alpha = 0.2f),
                            shape = CircleShape
                        ),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = post.author.first().toString(),
                        fontSize = 16.sp,
                        fontWeight = FontWeight.SemiBold,
                        color = Colors.primary
                    )
                }

                Spacer(modifier = Modifier.width(Spacing.md))

                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = post.author,
                        fontSize = 14.sp,
                        fontWeight = FontWeight.SemiBold,
                        color = Colors.textPrimary
                    )
                    Text(
                        text = post.time,
                        fontSize = 12.sp,
                        color = Colors.textSecondary
                    )
                }

                IconButton(onClick = { /* More options */ }) {
                    Icon(
                        imageVector = Icons.Default.MoreVert,
                        contentDescription = "More",
                        tint = Colors.textSecondary
                    )
                }
            }

            Spacer(modifier = Modifier.height(Spacing.sm))

            // Content
            Text(
                text = post.content,
                fontSize = 15.sp,
                color = Colors.textPrimary
            )

            Spacer(modifier = Modifier.height(Spacing.sm))

            // Icon/Media placeholder
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(80.dp)
                    .background(
                        color = Colors.surface,
                        shape = RoundedCornerShape(12.dp)
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = post.icon,
                    contentDescription = "Content",
                    tint = Colors.primary,
                    modifier = Modifier.size(24.dp)
                )
            }

            Spacer(modifier = Modifier.height(Spacing.sm))

            // Actions
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(Spacing.lg)
            ) {
                Row(
                    modifier = Modifier.clickable {
                        isLiked = !isLiked
                        currentLikes += if (isLiked) 1 else -1
                    },
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        imageVector = if (isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                        contentDescription = "Like",
                        tint = if (isLiked) androidx.compose.ui.graphics.Color.Red else Colors.textSecondary,
                        modifier = Modifier.size(20.dp)
                    )
                    Spacer(modifier = Modifier.width(Spacing.xs))
                    Text(
                        text = currentLikes.toString(),
                        fontSize = 14.sp,
                        color = Colors.textSecondary
                    )
                }

                Row(
                    modifier = Modifier.clickable { /* Comment */ },
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        imageVector = Icons.Default.ChatBubbleOutline,
                        contentDescription = "Comment",
                        tint = Colors.textSecondary,
                        modifier = Modifier.size(20.dp)
                    )
                    Spacer(modifier = Modifier.width(Spacing.xs))
                    Text(
                        text = post.comments.toString(),
                        fontSize = 14.sp,
                        color = Colors.textSecondary
                    )
                }

                IconButton(onClick = { /* Share */ }) {
                    Icon(
                        imageVector = Icons.Default.Share,
                        contentDescription = "Share",
                        tint = Colors.textSecondary,
                        modifier = Modifier.size(20.dp)
                    )
                }
            }
        }
    }
}

// MARK: - Discover View
@Composable
private fun DiscoverView() {
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(Spacing.md)
    ) {
        item {
            Text(
                text = "Trending Topics",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.textPrimary,
                modifier = Modifier.padding(bottom = Spacing.md)
            )
        }

        item {
            LazyVerticalGrid(
                columns = GridCells.Fixed(2),
                modifier = Modifier.height(360.dp),
                horizontalArrangement = Arrangement.spacedBy(Spacing.sm),
                verticalArrangement = Arrangement.spacedBy(Spacing.sm)
            ) {
                items(6) { index ->
                    Card(
                        modifier = Modifier.height(120.dp),
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(
                            containerColor = Colors.surface
                        )
                    ) {
                        Column(
                            modifier = Modifier.fillMaxSize(),
                            horizontalAlignment = Alignment.CenterHorizontally,
                            verticalArrangement = Arrangement.Center
                        ) {
                            Icon(
                                imageVector = Icons.Default.Whatshot,
                                contentDescription = "Trending",
                                tint = Colors.primary,
                                modifier = Modifier.size(24.dp)
                            )
                            Spacer(modifier = Modifier.height(Spacing.xs))
                            Text(
                                text = "Topic ${index + 1}",
                                fontSize = 14.sp,
                                fontWeight = FontWeight.Medium,
                                color = Colors.textPrimary
                            )
                        }
                    }
                }
            }
        }
    }
}

// MARK: - Notifications View
@Composable
private fun NotificationsView() {
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.sm)
    ) {
        items(5) { index ->
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = Spacing.md, vertical = Spacing.sm),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .background(
                            color = Colors.primary.copy(alpha = 0.2f),
                            shape = CircleShape
                        ),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        imageVector = Icons.Default.Notifications,
                        contentDescription = "Notification",
                        tint = Colors.primary,
                        modifier = Modifier.size(16.dp)
                    )
                }

                Spacer(modifier = Modifier.width(Spacing.md))

                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = "Notification ${index + 1}",
                        fontSize = 14.sp,
                        fontWeight = FontWeight.Medium,
                        color = Colors.textPrimary
                    )
                    Text(
                        text = "This is a sample notification message",
                        fontSize = 12.sp,
                        color = Colors.textSecondary
                    )
                }

                Text(
                    text = "2m",
                    fontSize = 12.sp,
                    color = Colors.textSecondary
                )
            }
        }
    }
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun SocialScreenPreview() {
    SocialScreen()
}