package com.tchat.mobile.screens

import androidx.compose.foundation.background
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
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun UserProfileScreen(
    userId: String = "1",
    onBackClick: () -> Unit,
    onEditClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    val user = getDummyUser(userId)
    var isFollowing by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(user.name) },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, "Back")
                    }
                },
                actions = {
                    IconButton(onClick = { /* Share profile */ }) {
                        Icon(Icons.Default.Share, "Share")
                    }
                    IconButton(onClick = { /* More options */ }) {
                        Icon(Icons.Default.MoreVert, "More")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface
                )
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .background(TchatColors.background),
            contentPadding = PaddingValues(TchatSpacing.md),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
        ) {
            // Profile Header
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                    elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.lg),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        // Profile Picture
                        Box(
                            modifier = Modifier
                                .size(120.dp)
                                .clip(CircleShape)
                                .background(TchatColors.primary),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = user.name.first().toString(),
                                color = TchatColors.onPrimary,
                                style = MaterialTheme.typography.headlineLarge,
                                fontWeight = FontWeight.Bold
                            )
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // User Info
                        Text(
                            text = user.name,
                            style = MaterialTheme.typography.headlineSmall,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )

                        Text(
                            text = user.username,
                            style = MaterialTheme.typography.bodyLarge,
                            color = TchatColors.onSurfaceVariant
                        )

                        Text(
                            text = user.bio,
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface,
                            modifier = Modifier.padding(top = TchatSpacing.sm)
                        )

                        Spacer(modifier = Modifier.height(TchatSpacing.md))

                        // Stats Row
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.lg)
                        ) {
                            StatItem("Posts", user.postsCount.toString())
                            StatItem("Followers", formatCount(user.followersCount))
                            StatItem("Following", user.followingCount.toString())
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.lg))

                        // Action Buttons
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                        ) {
                            if (user.isOwnProfile) {
                                TchatButton(
                                    onClick = onEditClick,
                                    text = "Edit Profile",
                                    variant = TchatButtonVariant.Primary,
                                    modifier = Modifier.weight(1f)
                                )
                            } else {
                                TchatButton(
                                    onClick = { isFollowing = !isFollowing },
                                    text = if (isFollowing) "Following" else "Follow",
                                    variant = if (isFollowing) TchatButtonVariant.Secondary else TchatButtonVariant.Primary,
                                    modifier = Modifier.weight(1f)
                                )
                                TchatButton(
                                    onClick = { /* Message user */ },
                                    text = "Message",
                                    variant = TchatButtonVariant.Outline,
                                    modifier = Modifier.weight(1f)
                                )
                            }
                        }
                    }
                }
            }

            // Recent Posts Section
            item {
                Text(
                    text = "Recent Posts",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurface
                )
            }

            items(user.recentPosts) { post ->
                PostCard(post = post)
            }
        }
    }
}

@Composable
private fun StatItem(
    label: String,
    value: String,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = value,
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = TchatColors.onSurface
        )
        Text(
            text = label,
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun PostCard(
    post: UserPost,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                text = post.content,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            Icons.Default.Favorite,
                            contentDescription = "Likes",
                            modifier = Modifier.size(16.dp),
                            tint = TchatColors.primary
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatCount(post.likesCount),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            Icons.Default.Send,
                            contentDescription = "Comments",
                            modifier = Modifier.size(16.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(
                            text = formatCount(post.commentsCount),
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                Text(
                    text = post.timestamp,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

private fun formatCount(count: Int): String {
    return when {
        count >= 1_000_000 -> "${count / 1_000_000}M"
        count >= 1_000 -> "${count / 1_000}K"
        else -> count.toString()
    }
}

// Data models
private data class UserProfile(
    val id: String,
    val name: String,
    val username: String,
    val bio: String,
    val postsCount: Int,
    val followersCount: Int,
    val followingCount: Int,
    val isOwnProfile: Boolean,
    val recentPosts: List<UserPost>
)

private data class UserPost(
    val id: String,
    val content: String,
    val timestamp: String,
    val likesCount: Int,
    val commentsCount: Int
)

private fun getDummyUser(userId: String): UserProfile = UserProfile(
    id = userId,
    name = "Alice Johnson",
    username = "@alice_designs",
    bio = "UI/UX Designer • Coffee enthusiast • Based in SF",
    postsCount = 234,
    followersCount = 12500,
    followingCount = 567,
    isOwnProfile = userId == "me",
    recentPosts = listOf(
        UserPost("1", "Just finished designing a new mobile app interface. Really excited about the clean, minimal approach we took!", "2h", 156, 23),
        UserPost("2", "Coffee break thoughts: The best designs are the ones users never notice because they just work seamlessly.", "5h", 89, 12),
        UserPost("3", "Working on some interesting accessibility improvements. Every user deserves a great experience!", "1d", 203, 34),
        UserPost("4", "Prototype vs final design - sometimes the best ideas come from happy accidents during development.", "2d", 145, 18),
        UserPost("5", "Attended an amazing design conference today. So many inspiring talks about the future of UX!", "3d", 267, 45)
    )
)