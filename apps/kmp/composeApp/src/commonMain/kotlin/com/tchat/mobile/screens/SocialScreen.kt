package com.tchat.mobile.screens

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalSoftwareKeyboardController
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import com.tchat.mobile.components.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import com.tchat.mobile.utils.PlatformUtils

// Enhanced data models matching web UI
data class UserItem(
    val id: String,
    val name: String,
    val username: String = "",
    val avatar: String = "",
    val isVerified: Boolean = false,
    val isOnline: Boolean = false,
    val lastSeen: String = "",
    val mutualFriends: Int = 0,
    val status: String = ""
)

data class PostItem(
    val id: String,
    val author: UserItem,
    val content: String,
    val timestamp: String,
    val likes: Int,
    val comments: Int,
    val shares: Int,
    val imageUrl: String? = null,
    val location: String? = null,
    val tags: List<String> = emptyList(),
    val type: PostType = PostType.TEXT,
    val source: PostSource = PostSource.FOLLOWING,
    val isLiked: Boolean = false
)

data class CommentItem(
    val id: String,
    val user: UserItem,
    val text: String,
    val timestamp: String,
    val likes: Int,
    val isLiked: Boolean
)

data class StoryItem(
    val id: String,
    val author: UserItem,
    val preview: String = "",
    val content: String = "",
    val timestamp: String = "",
    val isViewed: Boolean = false,
    val isLive: Boolean = false,
    val expiresAt: String = ""
)

data class FriendItem(
    val id: String,
    val name: String,
    val username: String,
    val avatar: String,
    val isOnline: Boolean,
    val isFollowing: Boolean,
    val mutualFriends: Int,
    val status: String
)

data class EventItem(
    val id: String,
    val title: String,
    val description: String,
    val date: String,
    val location: String,
    val price: String,
    val imageUrl: String,
    val attendeesCount: Int,
    val category: String,
    val isAttending: Boolean = false
)

enum class PostType { TEXT, IMAGE, LIVE, PRODUCT }
enum class PostSource { FOLLOWING, TRENDING, INTEREST, SPONSORED }

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SocialScreen(
    onUserClick: (userId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    // Tab and UI state
    var selectedTabIndex by remember { mutableStateOf(0) }
    var showShareModal by remember { mutableStateOf(false) }
    var sharedPost by remember { mutableStateOf<PostItem?>(null) }

    // Post interaction state - matching web UI
    var likedPosts by remember { mutableStateOf(setOf<String>()) }
    var bookmarkedPosts by remember { mutableStateOf(setOf<String>()) }
    var followingUsers by remember { mutableStateOf(setOf("1", "2", "3", "5")) }
    var commentsOpen by remember { mutableStateOf<String?>(null) }
    var newComment by remember { mutableStateOf("") }
    var postComments by remember { mutableStateOf(mapOf<String, List<CommentItem>>()) }

    // Story viewing state
    var viewingStory by remember { mutableStateOf<StoryItem?>(null) }
    var storyProgress by remember { mutableStateOf(0f) }

    // Post creation state - matching web functionality
    var showCreatePost by remember { mutableStateOf(false) }
    var newPostText by remember { mutableStateOf("") }
    var selectedImages by remember { mutableStateOf(listOf<String>()) }
    var postLocation by remember { mutableStateOf("") }
    var userPosts by remember { mutableStateOf(listOf<PostItem>()) }

    // Story creation state
    var showCreateStory by remember { mutableStateOf(false) }
    var storyText by remember { mutableStateOf("") }

    val tabs = listOf("Friends", "Feed", "Discover", "Events")
    val posts = remember { getDummyPosts() }
    val stories = remember { getDummyStories() }
    val friends = remember { getDummyFriends() }
    val events = remember { getDummyEvents() }

    val keyboardController = LocalSoftwareKeyboardController.current
    val scope = rememberCoroutineScope()

    // Story progress timer
    LaunchedEffect(viewingStory) {
        if (viewingStory != null) {
            while (storyProgress < 100f) {
                delay(50)
                storyProgress += 1f
            }
            if (storyProgress >= 100f) {
                delay(500)
                viewingStory = null
                storyProgress = 0f
            }
        }
    }

    // Event handlers - matching web UI functionality
    val handleLike: (String) -> Unit = { postId ->
        likedPosts = if (likedPosts.contains(postId)) {
            likedPosts - postId
        } else {
            likedPosts + postId
        }
    }

    val handleBookmark: (String) -> Unit = { postId ->
        bookmarkedPosts = if (bookmarkedPosts.contains(postId)) {
            bookmarkedPosts - postId
        } else {
            bookmarkedPosts + postId
        }
    }

    val handleFollow: (String) -> Unit = { userId ->
        followingUsers = if (followingUsers.contains(userId)) {
            followingUsers - userId
        } else {
            followingUsers + userId
        }
    }

    val handleHashtagClick: (String) -> Unit = { hashtag ->
        // Handle hashtag navigation
    }

    val handleComment: (String) -> Unit = { postId ->
        commentsOpen = if (commentsOpen == postId) null else postId
    }

    val handleAddComment: (String) -> Unit = { postId ->
        if (newComment.trim().isNotEmpty()) {
            val comment = CommentItem(
                id = "comment_${PlatformUtils.currentTimeMillis()}",
                user = UserItem("you", "You", "you"),
                text = newComment,
                timestamp = "just now",
                likes = 0,
                isLiked = false
            )
            postComments = postComments + (postId to (postComments[postId] ?: emptyList()) + comment)
            newComment = ""
            keyboardController?.hide()
        }
    }

    val onShare: (PostItem) -> Unit = { post ->
        sharedPost = post
        showShareModal = true
    }

    val handleStoryClick: (StoryItem) -> Unit = { story ->
        if (story.author.name == "Your Story") {
            showCreateStory = true
        } else {
            viewingStory = story
            storyProgress = 0f
        }
    }

    val handleCreatePost: (String, String, List<String>) -> Unit = { text, location, images ->
        if (text.trim().isNotEmpty()) {
            val newPost = PostItem(
                id = "user_post_${PlatformUtils.currentTimeMillis()}",
                author = UserItem("you", "You", "you"),
                content = text,
                timestamp = "just now",
                likes = 0,
                comments = 0,
                shares = 0,
                imageUrl = images.firstOrNull(),
                location = location.ifEmpty { null }
            )
            userPosts = listOf(newPost) + userPosts
            newPostText = ""
            postLocation = ""
            selectedImages = emptyList()
            showCreatePost = false
        }
    }

    // Combine user posts with feed posts
    val allPosts = userPosts + posts

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Top App Bar with notifications
        TopAppBar(
            title = { Text("Social", fontWeight = FontWeight.Bold) },
            actions = {
                IconButton(onClick = { /* Notifications */ }) {
                    BadgedBox(
                        badge = { Badge { Text("5") } }
                    ) {
                        Icon(Icons.Filled.Notifications, "Notifications")
                    }
                }
                IconButton(onClick = { showCreatePost = true }) {
                    Icon(Icons.Filled.Add, "New Post")
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        Column(
            modifier = Modifier.weight(1f)
        ) {
            // Stories Row - matching web UI
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
                contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                modifier = Modifier.padding(vertical = TchatSpacing.sm)
            ) {
                items(stories) { story ->
                    StoryItemCard(
                        story = story,
                        onClick = { handleStoryClick(story) }
                    )
                }
            }

            // Create Post Section - enhanced like web
            CreatePostSection(
                showDialog = showCreatePost,
                onDismiss = { showCreatePost = false },
                onCreatePost = handleCreatePost,
                newPostText = newPostText,
                onTextChange = { newPostText = it },
                postLocation = postLocation,
                onLocationChange = { postLocation = it },
                selectedImages = selectedImages
            )

            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md)
                    .clickable { showCreatePost = true },
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
            ) {
                Row(
                    modifier = Modifier.padding(TchatSpacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(40.dp)
                            .clip(CircleShape)
                            .background(TchatColors.primary),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            "Y",
                            color = TchatColors.onPrimary,
                            fontWeight = FontWeight.Bold
                        )
                    }

                    Spacer(modifier = Modifier.width(TchatSpacing.md))

                    Text(
                        "What's on your mind?",
                        style = MaterialTheme.typography.bodyLarge,
                        color = TchatColors.onSurfaceVariant,
                        modifier = Modifier.weight(1f)
                    )

                    Row {
                        IconButton(onClick = { showCreatePost = true }) {
                            Icon(Icons.Default.PhotoCamera, "Photo", tint = TchatColors.primary)
                        }
                        IconButton(onClick = { showCreatePost = true }) {
                            Icon(Icons.Default.LocationOn, "Location", tint = TchatColors.primary)
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            // Tab Navigation
            TabRow(
                selectedTabIndex = selectedTabIndex,
                containerColor = TchatColors.surface,
                contentColor = TchatColors.primary
            ) {
                tabs.forEachIndexed { index, title ->
                    Tab(
                        selected = selectedTabIndex == index,
                        onClick = { selectedTabIndex = index },
                        text = {
                            Text(
                                text = title,
                                style = MaterialTheme.typography.titleSmall,
                                fontWeight = if (selectedTabIndex == index) FontWeight.SemiBold else FontWeight.Normal
                            )
                        },
                        icon = {
                            Icon(
                                imageVector = when (index) {
                                    0 -> Icons.Default.People
                                    1 -> Icons.Default.Star
                                    2 -> Icons.Default.Explore
                                    3 -> Icons.Default.Event
                                    else -> Icons.Default.Home
                                },
                                contentDescription = title,
                                modifier = Modifier.size(20.dp)
                            )
                        }
                    )
                }
            }

            // Tab content based on selectedTabIndex
            when (selectedTabIndex) {
                0 -> FriendsTab(
                    posts = allPosts.filter { followingUsers.contains(it.author.id) || userPosts.contains(it) },
                    friends = friends,
                    likedPosts = likedPosts,
                    bookmarkedPosts = bookmarkedPosts,
                    followingUsers = followingUsers,
                    commentsOpen = commentsOpen,
                    newComment = newComment,
                    postComments = postComments,
                    onLike = handleLike,
                    onBookmark = handleBookmark,
                    onFollow = handleFollow,
                    onComment = handleComment,
                    onShare = onShare,
                    onHashtagClick = handleHashtagClick,
                    onCommentTextChange = { newComment = it },
                    onAddComment = handleAddComment
                )
                1 -> FeedTab(
                    posts = allPosts,
                    likedPosts = likedPosts,
                    bookmarkedPosts = bookmarkedPosts,
                    followingUsers = followingUsers,
                    commentsOpen = commentsOpen,
                    newComment = newComment,
                    postComments = postComments,
                    onLike = handleLike,
                    onBookmark = handleBookmark,
                    onFollow = handleFollow,
                    onComment = handleComment,
                    onShare = onShare,
                    onHashtagClick = handleHashtagClick,
                    onCommentTextChange = { newComment = it },
                    onAddComment = handleAddComment
                )
                2 -> DiscoverTab(
                    posts = posts.filter { it.author.name.contains("Explorer") || it.author.name.contains("Cultural") },
                    likedPosts = likedPosts,
                    bookmarkedPosts = bookmarkedPosts,
                    followingUsers = followingUsers,
                    commentsOpen = commentsOpen,
                    newComment = newComment,
                    postComments = postComments,
                    onLike = handleLike,
                    onBookmark = handleBookmark,
                    onFollow = handleFollow,
                    onComment = handleComment,
                    onShare = onShare,
                    onHashtagClick = handleHashtagClick,
                    onCommentTextChange = { newComment = it },
                    onAddComment = handleAddComment
                )
                3 -> EventsTab(
                    events = events,
                    posts = posts.filter { it.content.contains("event") },
                    likedPosts = likedPosts,
                    bookmarkedPosts = bookmarkedPosts,
                    followingUsers = followingUsers,
                    commentsOpen = commentsOpen,
                    newComment = newComment,
                    postComments = postComments,
                    onLike = handleLike,
                    onBookmark = handleBookmark,
                    onFollow = handleFollow,
                    onComment = handleComment,
                    onShare = onShare,
                    onHashtagClick = handleHashtagClick,
                    onCommentTextChange = { newComment = it },
                    onAddComment = handleAddComment
                )
            }
        }
    }

    // Story Viewer Dialog
    if (viewingStory != null) {
        StoryViewerDialog(
            story = viewingStory!!,
            progress = storyProgress,
            onDismiss = {
                viewingStory = null
                storyProgress = 0f
            }
        )
    }

    // Create Story Dialog
    if (showCreateStory) {
        CreateStoryDialog(
            storyText = storyText,
            onTextChange = { storyText = it },
            onDismiss = {
                showCreateStory = false
                storyText = ""
            },
            onCreate = {
                showCreateStory = false
                storyText = ""
            }
        )
    }

    // Share Modal
    if (showShareModal && sharedPost != null) {
        TchatShareModal(
            isVisible = showShareModal,
            content = ShareContent(
                title = sharedPost!!.content,
                description = "Shared from ${sharedPost!!.author.name}",
                url = "https://tchat.app/post/${sharedPost!!.id}"
            ),
            onDismiss = {
                showShareModal = false
                sharedPost = null
            },
            onShare = { platform, content ->
                println("Sharing to ${platform.name}: ${content.title}")
            },
            onCopyLink = { url ->
                println("Copied link: $url")
            }
        )
    }
}

// Enhanced composables matching web UI functionality

@Composable
private fun StoryItemCard(
    story: StoryItem,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier
            .width(60.dp)
            .clickable { onClick() }
    ) {
        Box(
            modifier = Modifier
                .size(56.dp)
                .clip(CircleShape)
                .border(
                    2.dp,
                    if (story.isViewed) TchatColors.surfaceVariant else TchatColors.primary,
                    CircleShape
                )
                .padding(2.dp)
        ) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .clip(CircleShape)
                    .background(
                        if (story.author.name == "Your Story") TchatColors.surfaceVariant
                        else TchatColors.primary.copy(alpha = 0.1f)
                    ),
                contentAlignment = Alignment.Center
            ) {
                if (story.author.name == "Your Story") {
                    Icon(
                        Icons.Default.Add,
                        contentDescription = "Add story",
                        tint = TchatColors.onSurfaceVariant
                    )
                } else {
                    Text(
                        story.author.name.first().toString(),
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.primary
                    )
                }
            }

            if (story.isLive) {
                Badge(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .offset(x = 2.dp, y = 2.dp),
                    containerColor = TchatColors.error
                ) {
                    Text("LIVE", style = MaterialTheme.typography.labelSmall)
                }
            }
        }

        Spacer(modifier = Modifier.height(4.dp))

        Text(
            text = if (story.author.name == "Your Story") "Your Story" else story.author.name.split(" ").first(),
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurface,
            maxLines = 1,
            overflow = TextOverflow.Ellipsis,
            textAlign = TextAlign.Center
        )
    }
}

@Composable
private fun CreatePostSection(
    showDialog: Boolean,
    onDismiss: () -> Unit,
    onCreatePost: (String, String, List<String>) -> Unit,
    newPostText: String,
    onTextChange: (String) -> Unit,
    postLocation: String,
    onLocationChange: (String) -> Unit,
    selectedImages: List<String>,
    modifier: Modifier = Modifier
) {
    if (showDialog) {
        Dialog(
            onDismissRequest = onDismiss,
            properties = DialogProperties(usePlatformDefaultWidth = false)
        ) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp)
            ) {
                Column(
                    modifier = Modifier.padding(16.dp)
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            "Create Post",
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold
                        )
                        IconButton(onClick = onDismiss) {
                            Icon(Icons.Default.Close, "Close")
                        }
                    }

                    Spacer(modifier = Modifier.height(16.dp))

                    OutlinedTextField(
                        value = newPostText,
                        onValueChange = onTextChange,
                        placeholder = { Text("What's on your mind?") },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(120.dp),
                        maxLines = 5
                    )

                    Spacer(modifier = Modifier.height(12.dp))

                    OutlinedTextField(
                        value = postLocation,
                        onValueChange = onLocationChange,
                        placeholder = { Text("Add location (optional)") },
                        leadingIcon = { Icon(Icons.Default.LocationOn, "Location") },
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(16.dp))

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        TchatButton(
                            onClick = onDismiss,
                            text = "Cancel",
                            variant = TchatButtonVariant.Outline,
                            modifier = Modifier.weight(1f)
                        )
                        TchatButton(
                            onClick = {
                                onCreatePost(newPostText, postLocation, selectedImages)
                            },
                            text = "Share Post",
                            variant = TchatButtonVariant.Primary,
                            modifier = Modifier.weight(1f)
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun StoryViewerDialog(
    story: StoryItem,
    progress: Float,
    onDismiss: () -> Unit
) {
    Dialog(
        onDismissRequest = onDismiss,
        properties = DialogProperties(usePlatformDefaultWidth = false)
    ) {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(Color.Black.copy(alpha = 0.9f))
                .clickable { onDismiss() },
            contentAlignment = Alignment.Center
        ) {
            Card(
                modifier = Modifier
                    .width(300.dp)
                    .height(500.dp),
                shape = RoundedCornerShape(16.dp)
            ) {
                Box(modifier = Modifier.fillMaxSize()) {
                    // Progress indicator
                    LinearProgressIndicator(
                        progress = { progress / 100f },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(2.dp)
                            .align(Alignment.TopCenter),
                        color = TchatColors.primary
                    )

                    // Story content
                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(16.dp),
                        verticalArrangement = Arrangement.Center,
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(
                            story.author.name,
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold
                        )

                        Spacer(modifier = Modifier.height(16.dp))

                        Text(
                            story.content,
                            style = MaterialTheme.typography.bodyLarge,
                            textAlign = TextAlign.Center
                        )

                        Spacer(modifier = Modifier.height(16.dp))

                        Text(
                            story.timestamp,
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun CreateStoryDialog(
    storyText: String,
    onTextChange: (String) -> Unit,
    onDismiss: () -> Unit,
    onCreate: () -> Unit
) {
    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            shape = RoundedCornerShape(16.dp)
        ) {
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Text(
                    "Create Your Story",
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold
                )

                Spacer(modifier = Modifier.height(16.dp))

                OutlinedTextField(
                    value = storyText,
                    onValueChange = onTextChange,
                    placeholder = { Text("Share what's happening...") },
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(100.dp),
                    maxLines = 3
                )

                Spacer(modifier = Modifier.height(16.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    TchatButton(
                        onClick = onDismiss,
                        text = "Cancel",
                        variant = TchatButtonVariant.Outline,
                        modifier = Modifier.weight(1f)
                    )
                    TchatButton(
                        onClick = onCreate,
                        text = "Share Story",
                        variant = TchatButtonVariant.Primary,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}

// Continue with Enhanced Tab Components...

@Composable
private fun FriendsTab(
    posts: List<PostItem>,
    friends: List<FriendItem>,
    likedPosts: Set<String>,
    bookmarkedPosts: Set<String>,
    followingUsers: Set<String>,
    commentsOpen: String?,
    newComment: String,
    postComments: Map<String, List<CommentItem>>,
    onLike: (String) -> Unit,
    onBookmark: (String) -> Unit,
    onFollow: (String) -> Unit,
    onComment: (String) -> Unit,
    onShare: (PostItem) -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Friends Activity Header
        item {
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = TchatColors.primary.copy(alpha = 0.1f)
                )
            ) {
                Row(
                    modifier = Modifier.padding(TchatSpacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(Icons.Default.People, "Friends", tint = TchatColors.primary)
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                    Text(
                        "Friends Activity",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.SemiBold
                    )
                    Spacer(modifier = Modifier.weight(1f))
                    Badge {
                        Text("${friends.count { followingUsers.contains(it.id) }} friends")
                    }
                }
            }
        }

        // Friend Suggestions
        if (friends.any { !followingUsers.contains(it.id) }) {
            item {
                FriendSuggestionsCard(
                    friends = friends.filter { !followingUsers.contains(it.id) }.take(2),
                    onFollow = onFollow
                )
            }
        }

        // Online Friends
        item {
            OnlineFriendsCard(
                friends = friends.filter { it.isOnline && followingUsers.contains(it.id) }
            )
        }

        // Friends' Posts
        items(posts) { post ->
            EnhancedPostCard(
                post = post,
                isLiked = likedPosts.contains(post.id),
                isBookmarked = bookmarkedPosts.contains(post.id),
                isFollowing = followingUsers.contains(post.author.id),
                commentsOpen = commentsOpen == post.id,
                newComment = newComment,
                comments = postComments[post.id] ?: emptyList(),
                onLike = { onLike(post.id) },
                onBookmark = { onBookmark(post.id) },
                onFollow = { onFollow(post.author.id) },
                onComment = { onComment(post.id) },
                onShare = { onShare(post) },
                onHashtagClick = onHashtagClick,
                onCommentTextChange = onCommentTextChange,
                onAddComment = { onAddComment(post.id) }
            )
        }
    }
}

@Composable
private fun FeedTab(
    posts: List<PostItem>,
    likedPosts: Set<String>,
    bookmarkedPosts: Set<String>,
    followingUsers: Set<String>,
    commentsOpen: String?,
    newComment: String,
    postComments: Map<String, List<CommentItem>>,
    onLike: (String) -> Unit,
    onBookmark: (String) -> Unit,
    onFollow: (String) -> Unit,
    onComment: (String) -> Unit,
    onShare: (PostItem) -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Feed Header with interests
        item {
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = TchatColors.surface
                )
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(Icons.Default.Star, "Your Feed", tint = TchatColors.primary)
                        Spacer(modifier = Modifier.width(TchatSpacing.sm))
                        Text(
                            "Your Interests",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.SemiBold
                        )
                    }

                    Spacer(modifier = Modifier.height(TchatSpacing.sm))

                    LazyRow(
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        items(listOf("🍜 Thai Food", "🏛️ Culture", "🛶 Markets", "🎵 Music", "📱 Tech", "✈️ Travel")) { interest ->
                            AssistChip(
                                onClick = { },
                                label = { Text(interest, style = MaterialTheme.typography.bodySmall) }
                            )
                        }
                    }
                }
            }
        }

        // All Posts
        items(posts) { post ->
            EnhancedPostCard(
                post = post,
                isLiked = likedPosts.contains(post.id),
                isBookmarked = bookmarkedPosts.contains(post.id),
                isFollowing = followingUsers.contains(post.author.id),
                commentsOpen = commentsOpen == post.id,
                newComment = newComment,
                comments = postComments[post.id] ?: emptyList(),
                onLike = { onLike(post.id) },
                onBookmark = { onBookmark(post.id) },
                onFollow = { onFollow(post.author.id) },
                onComment = { onComment(post.id) },
                onShare = { onShare(post) },
                onHashtagClick = onHashtagClick,
                onCommentTextChange = onCommentTextChange,
                onAddComment = { onAddComment(post.id) }
            )
        }

        // Suggested Content
        item {
            SuggestedContentCard()
        }
    }
}

@Composable
private fun DiscoverTab(
    posts: List<PostItem>,
    likedPosts: Set<String>,
    bookmarkedPosts: Set<String>,
    followingUsers: Set<String>,
    commentsOpen: String?,
    newComment: String,
    postComments: Map<String, List<CommentItem>>,
    onLike: (String) -> Unit,
    onBookmark: (String) -> Unit,
    onFollow: (String) -> Unit,
    onComment: (String) -> Unit,
    onShare: (PostItem) -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Discover Header
        item {
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = TchatColors.primary.copy(alpha = 0.1f)
                )
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(Icons.Default.Explore, "Discover", tint = TchatColors.primary)
                        Spacer(modifier = Modifier.width(TchatSpacing.sm))
                        Text(
                            "Discover",
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer(modifier = Modifier.weight(1f))
                        Badge(containerColor = TchatColors.error) {
                            Text("🔥 Hot", color = TchatColors.onPrimary)
                        }
                    }
                    Text(
                        "Trending content and new discoveries from Thailand",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }

        // Trending Categories
        item {
            TrendingCategoriesCard()
        }

        // Trending Posts
        items(posts) { post ->
            EnhancedPostCard(
                post = post.copy(
                    content = "🔥 TRENDING: ${post.content}"
                ),
                isLiked = likedPosts.contains(post.id),
                isBookmarked = bookmarkedPosts.contains(post.id),
                isFollowing = followingUsers.contains(post.author.id),
                commentsOpen = commentsOpen == post.id,
                newComment = newComment,
                comments = postComments[post.id] ?: emptyList(),
                onLike = { onLike(post.id) },
                onBookmark = { onBookmark(post.id) },
                onFollow = { onFollow(post.author.id) },
                onComment = { onComment(post.id) },
                onShare = { onShare(post) },
                onHashtagClick = onHashtagClick,
                onCommentTextChange = onCommentTextChange,
                onAddComment = { onAddComment(post.id) },
                showSourceBadge = true,
                sourceBadge = "Trending"
            )
        }
    }
}

@Composable
private fun EventsTab(
    events: List<EventItem>,
    posts: List<PostItem>,
    likedPosts: Set<String>,
    bookmarkedPosts: Set<String>,
    followingUsers: Set<String>,
    commentsOpen: String?,
    newComment: String,
    postComments: Map<String, List<CommentItem>>,
    onLike: (String) -> Unit,
    onBookmark: (String) -> Unit,
    onFollow: (String) -> Unit,
    onComment: (String) -> Unit,
    onShare: (PostItem) -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Events Header
        item {
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = TchatColors.primary.copy(alpha = 0.1f)
                )
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(Icons.Default.Event, "Events", tint = TchatColors.primary)
                        Spacer(modifier = Modifier.width(TchatSpacing.sm))
                        Text(
                            "Local Events",
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer(modifier = Modifier.weight(1f))
                        Badge(containerColor = TchatColors.error) {
                            Text("🔥 Hot", color = TchatColors.onPrimary)
                        }
                    }
                    Text(
                        "Discover festivals, food markets, and cultural events in Thailand",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }

        // Featured Events
        item {
            Text(
                "Featured Events",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold,
                modifier = Modifier.padding(bottom = TchatSpacing.sm)
            )
        }

        // Event Cards Grid
        items(events.chunked(2)) { eventPair ->
            Row(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                eventPair.forEach { event ->
                    EventCard(
                        event = event,
                        modifier = Modifier.weight(1f)
                    )
                }
                if (eventPair.size == 1) {
                    Spacer(modifier = Modifier.weight(1f))
                }
            }
        }

        // Event Categories
        item {
            EventCategoriesCard()
        }

        // Event-related Posts
        item {
            Text(
                "Event Posts",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold,
                modifier = Modifier.padding(vertical = TchatSpacing.sm)
            )
        }

        items(posts) { post ->
            EnhancedPostCard(
                post = post,
                isLiked = likedPosts.contains(post.id),
                isBookmarked = bookmarkedPosts.contains(post.id),
                isFollowing = followingUsers.contains(post.author.id),
                commentsOpen = commentsOpen == post.id,
                newComment = newComment,
                comments = postComments[post.id] ?: emptyList(),
                onLike = { onLike(post.id) },
                onBookmark = { onBookmark(post.id) },
                onFollow = { onFollow(post.author.id) },
                onComment = { onComment(post.id) },
                onShare = { onShare(post) },
                onHashtagClick = onHashtagClick,
                onCommentTextChange = onCommentTextChange,
                onAddComment = { onAddComment(post.id) }
            )
        }
    }
}

// Enhanced Post Card Component matching web UI
@Composable
private fun EnhancedPostCard(
    post: PostItem,
    isLiked: Boolean,
    isBookmarked: Boolean,
    isFollowing: Boolean,
    commentsOpen: Boolean,
    newComment: String,
    comments: List<CommentItem>,
    onLike: () -> Unit,
    onBookmark: () -> Unit,
    onFollow: () -> Unit,
    onComment: () -> Unit,
    onShare: () -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: () -> Unit,
    showSourceBadge: Boolean = false,
    sourceBadge: String = "",
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Column(
            modifier = Modifier.fillMaxWidth()
        ) {
            // Post Header with enhanced features
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Avatar
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primary),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = post.author.name.first().toString(),
                        color = TchatColors.onPrimary,
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Bold
                    )
                }

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                Column(modifier = Modifier.weight(1f)) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        Text(
                            text = post.author.name,
                            style = MaterialTheme.typography.titleSmall,
                            fontWeight = FontWeight.SemiBold,
                            color = TchatColors.onSurface
                        )

                        if (post.author.isVerified) {
                            Icon(
                                Icons.Default.Star,
                                contentDescription = "Verified",
                                modifier = Modifier.size(16.dp),
                                tint = TchatColors.primary
                            )
                        }

                        if (showSourceBadge) {
                            Badge {
                                Text(sourceBadge, style = MaterialTheme.typography.labelSmall)
                            }
                        }

                        if (!isFollowing && post.author.name != "You") {
                            TchatButton(
                                onClick = onFollow,
                                text = "Follow",
                                variant = TchatButtonVariant.Outline,
                                modifier = Modifier.height(24.dp)
                            )
                        }
                    }

                    Row {
                        Text(
                            text = post.timestamp,
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )

                        if (post.location != null) {
                            Text(" • ", color = TchatColors.onSurfaceVariant)
                            Icon(
                                Icons.Default.LocationOn,
                                contentDescription = null,
                                modifier = Modifier.size(12.dp),
                                tint = TchatColors.onSurfaceVariant
                            )
                            Text(
                                post.location,
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }
                    }
                }

                // More options menu
                IconButton(onClick = { /* More options */ }) {
                    Icon(
                        Icons.Filled.MoreVert,
                        contentDescription = "More options",
                        tint = TchatColors.onSurfaceVariant
                    )
                }
            }

            // Post Content with hashtag support
            if (post.content.isNotEmpty()) {
                val words = post.content.split(" ")
                val annotatedContent = buildString {
                    words.forEachIndexed { index, word ->
                        if (word.startsWith("#")) {
                            append(word)
                        } else {
                            append(word)
                        }
                        if (index < words.size - 1) append(" ")
                    }
                }

                Text(
                    text = annotatedContent,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    modifier = Modifier
                        .padding(horizontal = TchatSpacing.md)
                        .clickable {
                            // Handle hashtag clicks in a real implementation
                            words.forEach { word ->
                                if (word.startsWith("#")) {
                                    onHashtagClick(word)
                                }
                            }
                        }
                )
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
            }

            // Post Image
            if (post.imageUrl != null) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(200.dp)
                        .background(TchatColors.primary.copy(alpha = 0.1f)),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        Icons.Default.Image,
                        contentDescription = "Post Image",
                        modifier = Modifier.size(48.dp),
                        tint = TchatColors.primary
                    )
                }
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
            }

            // Enhanced engagement stats
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "${post.likes + if (isLiked) 1 else 0} likes",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
                Text(
                    text = "${post.comments + comments.size} comments",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
                Text(
                    text = "${post.shares} shares",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xs))
            Divider(color = TchatColors.outline.copy(alpha = 0.3f))

            // Enhanced action buttons
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.sm),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                TextButton(
                    onClick = onLike,
                    modifier = Modifier.weight(1f)
                ) {
                    Icon(
                        if (isLiked) Icons.Filled.Favorite else Icons.Filled.FavoriteBorder,
                        contentDescription = "Like",
                        modifier = Modifier.size(18.dp),
                        tint = if (isLiked) TchatColors.error else TchatColors.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = if (isLiked) "Liked" else "Like",
                        color = if (isLiked) TchatColors.error else TchatColors.onSurfaceVariant
                    )
                }

                TextButton(
                    onClick = onComment,
                    modifier = Modifier.weight(1f)
                ) {
                    Icon(
                        Icons.Default.ChatBubbleOutline,
                        contentDescription = "Comment",
                        modifier = Modifier.size(18.dp),
                        tint = TchatColors.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Comment",
                        color = TchatColors.onSurfaceVariant
                    )
                }

                TextButton(
                    onClick = onShare,
                    modifier = Modifier.weight(1f)
                ) {
                    Icon(
                        Icons.Filled.Share,
                        contentDescription = "Share",
                        modifier = Modifier.size(18.dp),
                        tint = TchatColors.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Share",
                        color = TchatColors.onSurfaceVariant
                    )
                }

                IconButton(onClick = onBookmark) {
                    Icon(
                        if (isBookmarked) Icons.Filled.Bookmark else Icons.Filled.BookmarkBorder,
                        contentDescription = "Bookmark",
                        modifier = Modifier.size(20.dp),
                        tint = if (isBookmarked) TchatColors.primary else TchatColors.onSurfaceVariant
                    )
                }
            }

            // Comments Section
            if (commentsOpen) {
                Divider(color = TchatColors.outline.copy(alpha = 0.3f))

                // Comment input
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(TchatSpacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    OutlinedTextField(
                        value = newComment,
                        onValueChange = onCommentTextChange,
                        placeholder = { Text("Write a comment...") },
                        modifier = Modifier.weight(1f),
                        keyboardOptions = KeyboardOptions(imeAction = ImeAction.Send),
                        keyboardActions = KeyboardActions(onSend = { onAddComment() })
                    )

                    Spacer(modifier = Modifier.width(TchatSpacing.sm))

                    IconButton(
                        onClick = onAddComment,
                        enabled = newComment.trim().isNotEmpty()
                    ) {
                        Icon(
                            Icons.Default.Send,
                            contentDescription = "Send comment",
                            tint = if (newComment.trim().isNotEmpty()) TchatColors.primary else TchatColors.onSurfaceVariant
                        )
                    }
                }

                // Comments list
                comments.forEach { comment ->
                    CommentCard(comment = comment)
                }
            }
        }
    }
}

@Composable
private fun CommentCard(
    comment: CommentItem,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = TchatSpacing.md, vertical = TchatSpacing.xs)
    ) {
        Box(
            modifier = Modifier
                .size(32.dp)
                .clip(CircleShape)
                .background(TchatColors.primary),
            contentAlignment = Alignment.Center
        ) {
            Text(
                comment.user.name.first().toString(),
                color = TchatColors.onPrimary,
                style = MaterialTheme.typography.labelLarge,
                fontWeight = FontWeight.Bold
            )
        }

        Spacer(modifier = Modifier.width(TchatSpacing.sm))

        Column(modifier = Modifier.weight(1f)) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
            ) {
                Text(
                    comment.user.name,
                    style = MaterialTheme.typography.labelMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
                Text(
                    comment.timestamp,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            Text(
                comment.text,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )

            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                TextButton(
                    onClick = { /* Like comment */ },
                    contentPadding = PaddingValues(0.dp)
                ) {
                    Icon(
                        if (comment.isLiked) Icons.Default.Favorite else Icons.Default.FavoriteBorder,
                        contentDescription = "Like",
                        modifier = Modifier.size(14.dp),
                        tint = if (comment.isLiked) TchatColors.error else TchatColors.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        "${comment.likes}",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                TextButton(
                    onClick = { /* Reply */ },
                    contentPadding = PaddingValues(0.dp)
                ) {
                    Text(
                        "Reply",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

// Additional helper components
@Composable
private fun FriendSuggestionsCard(
    friends: List<FriendItem>,
    onFollow: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.primary.copy(alpha = 0.05f)
        )
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                "Friend Suggestions",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            friends.forEach { friend ->
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = TchatSpacing.xs),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(32.dp)
                            .clip(CircleShape)
                            .background(TchatColors.primary),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            friend.name.first().toString(),
                            color = TchatColors.onPrimary,
                            fontWeight = FontWeight.Bold
                        )
                    }

                    Spacer(modifier = Modifier.width(TchatSpacing.sm))

                    Column(modifier = Modifier.weight(1f)) {
                        Text(
                            friend.name,
                            style = MaterialTheme.typography.bodyMedium,
                            fontWeight = FontWeight.Medium
                        )
                        Text(
                            "${friend.mutualFriends} mutual friends",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    TchatButton(
                        onClick = { onFollow(friend.id) },
                        text = "Add",
                        variant = TchatButtonVariant.Primary
                    )
                }
            }
        }
    }
}

@Composable
private fun OnlineFriendsCard(
    friends: List<FriendItem>,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.Circle,
                    contentDescription = "Online",
                    modifier = Modifier.size(12.dp),
                    tint = Color.Green
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Text(
                    "Friends Online",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold
                )
                Spacer(modifier = Modifier.weight(1f))
                Badge {
                    Text("${friends.size} online")
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                items(friends) { friend ->
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Box {
                            Box(
                                modifier = Modifier
                                    .size(40.dp)
                                    .clip(CircleShape)
                                    .background(TchatColors.primary),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    friend.name.first().toString(),
                                    color = TchatColors.onPrimary,
                                    fontWeight = FontWeight.Bold
                                )
                            }

                            // Online indicator
                            Box(
                                modifier = Modifier
                                    .size(12.dp)
                                    .clip(CircleShape)
                                    .background(Color.Green)
                                    .border(2.dp, TchatColors.surface, CircleShape)
                                    .align(Alignment.BottomEnd)
                            )
                        }

                        Spacer(modifier = Modifier.height(4.dp))

                        Text(
                            friend.name.split(" ").first(),
                            style = MaterialTheme.typography.bodySmall,
                            maxLines = 1,
                            overflow = TextOverflow.Ellipsis
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun SuggestedContentCard(
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.primary.copy(alpha = 0.1f)
        )
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                "Suggested for You",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            listOf(
                "Thai Cooking Enthusiasts" to "142K members • Food Group",
                "Bangkok Hidden Gems" to "89K members • Local Community"
            ).forEach { (title, description) ->
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = TchatSpacing.xs),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(40.dp)
                            .clip(RoundedCornerShape(8.dp))
                            .background(TchatColors.primary.copy(alpha = 0.2f)),
                        contentAlignment = Alignment.Center
                    ) {
                        Icon(
                            Icons.Default.Group,
                            contentDescription = "Group",
                            tint = TchatColors.primary
                        )
                    }

                    Spacer(modifier = Modifier.width(TchatSpacing.sm))

                    Column(modifier = Modifier.weight(1f)) {
                        Text(
                            title,
                            style = MaterialTheme.typography.bodyMedium,
                            fontWeight = FontWeight.Medium
                        )
                        Text(
                            description,
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }

                    TchatButton(
                        onClick = { },
                        text = "Join",
                        variant = TchatButtonVariant.Outline
                    )
                }
            }
        }
    }
}

@Composable
private fun TrendingCategoriesCard(
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                "Trending Categories",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                items(listOf(
                    "🍜 Food" to Color(0xFFF59E0B),
                    "🏛️ Culture" to Color(0xFF3B82F6),
                    "🎵 Music" to Color(0xFF8B5CF6),
                    "🏖️ Travel" to Color(0xFF10B981)
                )) { (category, color) ->
                    AssistChip(
                        onClick = { },
                        label = { Text(category) },
                        colors = AssistChipDefaults.assistChipColors(
                            containerColor = color.copy(alpha = 0.1f),
                            labelColor = color
                        )
                    )
                }
            }
        }
    }
}

@Composable
private fun EventCard(
    event: EventItem,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier,
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Column {
            // Event image placeholder
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .background(TchatColors.primary.copy(alpha = 0.1f)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.Event,
                    contentDescription = "Event",
                    modifier = Modifier.size(32.dp),
                    tint = TchatColors.primary
                )
            }

            Column(
                modifier = Modifier.padding(TchatSpacing.md)
            ) {
                Text(
                    event.title,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.SemiBold,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(TchatSpacing.xs))

                Text(
                    event.date,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                Text(
                    event.location,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        event.price,
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.primary
                    )

                    TchatButton(
                        onClick = { },
                        text = if (event.isAttending) "Going" else "RSVP",
                        variant = if (event.isAttending) TchatButtonVariant.Secondary else TchatButtonVariant.Primary,
                        modifier = Modifier.height(32.dp)
                    )
                }
            }
        }
    }
}

@Composable
private fun EventCategoriesCard(
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                "Browse by Category",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold
            )

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                listOf(
                    Triple("Music", Icons.Default.LibraryMusic, "15 events"),
                    Triple("Food", Icons.Default.Restaurant, "23 events")
                ).forEach { (category, icon, count) ->
                    Card(
                        modifier = Modifier.weight(1f),
                        colors = CardDefaults.cardColors(
                            containerColor = TchatColors.primary.copy(alpha = 0.1f)
                        )
                    ) {
                        Column(
                            modifier = Modifier.padding(TchatSpacing.md),
                            horizontalAlignment = Alignment.CenterHorizontally
                        ) {
                            Icon(
                                icon,
                                contentDescription = category,
                                modifier = Modifier.size(24.dp),
                                tint = TchatColors.primary
                            )
                            Text(
                                category,
                                style = MaterialTheme.typography.labelMedium,
                                fontWeight = FontWeight.Medium
                            )
                            Text(
                                count,
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

// Dummy data functions
private fun getDummyPosts(): List<PostItem> = listOf(
    PostItem(
        "1",
        UserItem("1", "Thai Food Explorer", "@thai_explorer", isVerified = true),
        "TRENDING: Secret street food spots only locals know about! 🏮🍜 This hidden gem in Chinatown serves the most authentic Tom Yum I've ever tasted. #HiddenGems #TomYum #Chinatown #StreetFood",
        "15 min ago",
        847, 123, 45,
        imageUrl = "https://example.com/tomyum.jpg",
        location = "Chinatown, Bangkok",
        tags = listOf("#HiddenGems", "#TomYum", "#Chinatown", "#StreetFood"),
        source = PostSource.TRENDING
    ),
    PostItem(
        "2",
        UserItem("2", "Sarah Johnson", "@sarah_foodie"),
        "Just tried the most amazing Pad Thai at Chatuchak Market! 🍜✨ The vendor taught me his secret ingredient - tamarind paste mixed with palm sugar. Mind blown! 🤯 #PadThai #StreetFood #Bangkok",
        "25 min ago",
        47, 12, 3,
        imageUrl = "https://example.com/padthai.jpg",
        location = "Chatuchak Weekend Market, Bangkok",
        tags = listOf("#PadThai", "#StreetFood", "#Bangkok"),
        source = PostSource.FOLLOWING,
        isLiked = true
    ),
    PostItem(
        "3",
        UserItem("3", "Bangkok Food Tours", "@bkk_tours", isVerified = true),
        "🔥 SPONSORED: Join our sunset food tour tonight! Explore 5 authentic local restaurants, meet fellow food lovers, and discover Bangkok's culinary secrets. Book now and get 20% off! 🌅🍽️ #FoodTour #Bangkok #Sponsored",
        "45 min ago",
        234, 56, 18,
        imageUrl = "https://example.com/foodtour.jpg",
        location = "Bangkok",
        tags = listOf("#FoodTour", "#Bangkok", "#Sponsored"),
        source = PostSource.SPONSORED
    ),
    PostItem(
        "4",
        UserItem("4", "Cultural Thailand", "@culture_th", isVerified = true),
        "Did you know? The tradition of floating markets dates back over 150 years! 🛶 These waterways were the original highways of Thailand, connecting communities through trade and culture. Which floating market is your favorite? #FloatingMarket #Culture #History #Thailand",
        "1 hour ago",
        892, 167, 234,
        imageUrl = "https://example.com/floatingmarket.jpg",
        location = "Thailand",
        tags = listOf("#FloatingMarket", "#Culture", "#History", "#Thailand"),
        source = PostSource.INTEREST
    ),
    PostItem(
        "5",
        UserItem("5", "Mike Chen", "@mike_travels"),
        "🔥 Going LIVE from the floating market! Come join me as I explore traditional boats selling fresh fruits and local delicacies. The energy here is incredible! 🛶💫 #FloatingMarket #Thailand #LiveStream",
        "1 hour ago",
        89, 23, 8,
        imageUrl = "https://example.com/live.jpg",
        location = "Damnoen Saduak Floating Market",
        tags = listOf("#FloatingMarket", "#Thailand", "#LiveStream"),
        type = PostType.LIVE,
        source = PostSource.FOLLOWING
    )
)

private fun getDummyStories(): List<StoryItem> = listOf(
    StoryItem(
        "0",
        UserItem("you", "Your Story", "you"),
        content = "Add to your story",
        timestamp = "now"
    ),
    StoryItem(
        "1",
        UserItem("1", "Sarah Johnson", "@sarah_foodie"),
        preview = "https://example.com/story1.jpg",
        content = "Amazing Pad Thai at Chatuchak Market! 🍜✨",
        timestamp = "2h ago",
        isViewed = false
    ),
    StoryItem(
        "2",
        UserItem("2", "Mike Chen", "@mike_travels"),
        preview = "https://example.com/story2.jpg",
        content = "Live from the floating market! 🛶",
        timestamp = "5m ago",
        isViewed = false,
        isLive = true
    ),
    StoryItem(
        "3",
        UserItem("3", "Emma Wilson", "@emma_culture"),
        preview = "https://example.com/story3.jpg",
        content = "Temple hopping day! 🏛️",
        timestamp = "1d ago",
        isViewed = true
    )
)

private fun getDummyFriends(): List<FriendItem> = listOf(
    FriendItem(
        "1", "Sarah Johnson", "@sarah_foodie", "",
        isOnline = true, isFollowing = true, mutualFriends = 12,
        status = "Exploring Bangkok street food! 🍜"
    ),
    FriendItem(
        "2", "Mike Chen", "@mike_travels", "",
        isOnline = true, isFollowing = true, mutualFriends = 8,
        status = "Live streaming from floating market!"
    ),
    FriendItem(
        "3", "Emma Wilson", "@emma_culture", "",
        isOnline = false, isFollowing = true, mutualFriends = 15,
        status = "Temple hopping in Bangkok"
    ),
    FriendItem(
        "4", "Alex Thai", "@alex_local", "",
        isOnline = true, isFollowing = true, mutualFriends = 23,
        status = "Local Bangkok guide 🏛️"
    ),
    FriendItem(
        "5", "Luna Park", "@luna_markets", "",
        isOnline = false, isFollowing = false, mutualFriends = 7,
        status = "Market photography enthusiast"
    )
)

private fun getDummyEvents(): List<EventItem> = listOf(
    EventItem(
        "1", "Bangkok Electronic Music Festival 2025",
        "Calvin Harris, Armin van Buuren & more",
        "March 15-17, 2025", "Bangkok Convention Center",
        "From ฿2,500", "", 18500, "Music"
    ),
    EventItem(
        "2", "Thai Street Food Championship",
        "Master chefs compete & cooking workshops",
        "Feb 28 - Mar 1, 2025", "Lumpini Park",
        "From ฿500", "", 12000, "Food"
    ),
    EventItem(
        "3", "Songkran Cultural Festival",
        "Traditional water festival celebration",
        "April 13-15, 2025", "Khao San Road",
        "Free", "", 85000, "Culture"
    ),
    EventItem(
        "4", "Floating Market Food Festival",
        "Traditional boat vendors and local delicacies",
        "March 20-22, 2025", "Damnoen Saduak",
        "From ฿200", "", 25000, "Food"
    )
)