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
import androidx.compose.material3.TabRowDefaults.tabIndicatorOffset
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
import com.tchat.mobile.models.Post
import com.tchat.mobile.models.PostType as UnifiedPostType
import com.tchat.mobile.models.PostType
import com.tchat.mobile.models.PostUser
import com.tchat.mobile.models.PostContent
import com.tchat.mobile.models.PostContentType
import com.tchat.mobile.models.PostInteractions
import com.tchat.mobile.components.posts.PostRenderer
import com.tchat.mobile.services.NavigationService
import com.tchat.mobile.services.NavigationAction
import com.tchat.mobile.services.SharingService
import com.tchat.mobile.services.SharingPlatform
import com.tchat.mobile.services.ShareResult
import com.tchat.mobile.services.SocialContentService
import com.tchat.mobile.services.ContentApiService
import com.tchat.mobile.repositories.MockPostRepository
import com.tchat.mobile.repositories.EventRepository
import com.tchat.mobile.services.MockSharingService
import com.tchat.mobile.services.MockNavigationService
import com.tchat.mobile.models.*
import org.koin.compose.koinInject

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SocialScreen(
    onUserClick: (userId: String) -> Unit = {},
    onMoreClick: () -> Unit = {},
    onEventClick: (eventId: String) -> Unit = {},
    onCategoryClick: (categoryId: String, categoryName: String) -> Unit = { _, _ -> },
    socialContentService: SocialContentService? = null,
    contentApiService: ContentApiService? = null,
    modifier: Modifier = Modifier
) {
    // Use new social architecture
    com.tchat.mobile.social.presentation.SocialScreen(modifier = modifier)
    return
    println("🎯 SocialScreen composable started")

    // Inject EventRepository
    val eventRepository: EventRepository = koinInject<EventRepository>()
    println("✅ EventRepository injected successfully")

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

    val tabs = listOf("Friends", "Feed", "All Posts", "Discover", "Events")

    // Real social data state
    var realStories by remember { mutableStateOf<List<Story>>(emptyList()) }
    var realFriends by remember { mutableStateOf<List<Friend>>(emptyList()) }
    var realEvents by remember { mutableStateOf<List<Event>>(emptyList()) }
    var socialDataLoading by remember { mutableStateOf(false) }

    // Content API data state
    var apiPosts by remember { mutableStateOf<List<Post>>(emptyList()) }
    var apiStories by remember { mutableStateOf<List<StoryItem>>(emptyList()) }
    var contentApiLoading by remember { mutableStateOf(false) }

    // Event repository data state
    var eventDataLoading by remember { mutableStateOf(false) }
    var eventCategories by remember { mutableStateOf<List<Pair<String, Int>>>(emptyList()) }
    var eventPosts by remember { mutableStateOf<List<PostItem>>(emptyList()) }

    // RSVP state management
    var rsvpEvents by remember { mutableStateOf(setOf<String>()) }

    // Load real social data if service is available
    LaunchedEffect(socialContentService) {
        socialContentService?.let { service ->
            socialDataLoading = true
            try {
                // Load stories
                service.getPersonalizedStories().onSuccess { stories ->
                    realStories = stories
                }

                // Load friends
                service.getFriendsWithStatus().onSuccess { friends ->
                    realFriends = friends
                }

                // Load events
                service.getUpcomingEvents().onSuccess { events ->
                    realEvents = events
                }
            } catch (e: Exception) {
                println("Error loading social data: ${e.message}")
            } finally {
                socialDataLoading = false
            }
        }
    }

    // Initialize event repository with seed data and load categories/posts
    LaunchedEffect(eventRepository) {
        eventDataLoading = true
        println("🚀 Starting EventRepository initialization...")
        try {
            // Initialize with seed data (creates categories and posts if they don't exist)
            println("📊 Initializing with seed data...")
            val seedResult = eventRepository.initializeWithSeedData()
            if (seedResult.isSuccess) {
                println("✅ Database seeding successful")
            } else {
                println("❌ Database seeding failed: ${seedResult.exceptionOrNull()?.message}")
            }

            // Load event categories for browse section
            println("📋 Loading event categories...")
            eventCategories = eventRepository.getEventCategoriesForUI()
            println("✅ Loaded ${eventCategories.size} event categories: ${eventCategories.map { "${it.first}(${it.second})" }}")

            // Load event posts for event posts section
            println("📝 Loading event posts...")
            eventPosts = eventRepository.getEventPostsForUI()
            println("✅ Loaded ${eventPosts.size} event posts")
            eventPosts.forEach { post ->
                println("   - ${post.author.name}: ${post.content.take(50)}...")
            }

        } catch (e: Exception) {
            println("❌ Error loading event data: ${e.message}")
            e.printStackTrace()
            // Fallback to hardcoded data
            eventCategories = listOf(
                "Music" to 15,
                "Food" to 23,
                "Technology" to 8,
                "Arts & Culture" to 12
            )
            eventPosts = emptyList()
            println("🔄 Using fallback data: ${eventCategories.size} categories, ${eventPosts.size} posts")
        } finally {
            eventDataLoading = false
            println("🏁 EventRepository initialization complete")
        }
    }

    // Load content from content service API
    LaunchedEffect(contentApiService) {
        contentApiService?.let { service ->
            contentApiLoading = true
            try {
                // Load posts from content service
                service.getSocialPosts().onSuccess { posts ->
                    apiPosts = posts
                    println("✅ Loaded ${posts.size} posts from content service")
                }

                // Load stories from content service
                service.getSocialStories().onSuccess { stories ->
                    apiStories = stories.map { story ->
                        StoryItem(
                            id = story.id,
                            author = UserItem(
                                id = story.authorId,
                                name = story.authorId, // Use authorId as name for now
                                username = story.authorId,
                                avatar = "",
                                isVerified = false,
                                isOnline = false,
                                lastSeen = "",
                                mutualFriends = 0,
                                status = ""
                            ),
                            preview = story.preview,
                            content = story.content,
                            timestamp = story.createdAt.toString(),
                            isViewed = story.isViewed,
                            isLive = story.isLive,
                            expiresAt = story.expiresAt.toString()
                        )
                    }
                    println("✅ Loaded ${stories.size} stories from content service")
                }
            } catch (e: Exception) {
                println("❌ Error loading content from API: ${e.message}")
            } finally {
                contentApiLoading = false
            }
        }
    }

    // Use real data only, fallback to content API or empty list
    val stories = if (realStories.isNotEmpty()) {
        realStories.map { story ->
            StoryItem(
                id = story.id,
                author = UserItem(story.authorId, "User ${story.authorId}", story.authorId),
                content = story.content,
                timestamp = "2h",
                isViewed = story.isViewed,
                isLive = story.isLive
            )
        }
    } else if (apiStories.isNotEmpty()) {
        apiStories
    } else {
        emptyList()
    }

    val friends = if (realFriends.isNotEmpty()) {
        realFriends.map { friend ->
            FriendItem(
                id = friend.id,
                name = friend.profile?.displayName ?: "Unknown",
                username = friend.profile?.username ?: "unknown",
                avatar = friend.profile?.avatarUrl ?: "",
                isOnline = friend.profile?.isOnline ?: false,
                isFollowing = followingUsers.contains(friend.id),
                mutualFriends = friend.mutualFriendsCount,
                status = friend.profile?.statusMessage ?: ""
            )
        }
    } else {
        emptyList()
    }

    val events = if (realEvents.isNotEmpty()) {
        realEvents.map { event ->
            EventItem(
                id = event.id,
                title = event.title,
                description = event.description,
                date = "Dec ${(event.eventDate % 31) + 1}",
                location = event.location,
                price = event.price,
                imageUrl = event.imageUrl ?: "",
                attendeesCount = event.attendeesCount,
                category = event.category,
                isAttending = rsvpEvents.contains(event.id)
            )
        }
    } else {
        // Fallback to dummy events when real events are empty
        SocialMockData.getDummyEvents()
    }

    // Use real posts from content API or empty list
    val allPosts = if (apiPosts.isNotEmpty()) {
        apiPosts
    } else {
        emptyList()
    }

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

    // Combine user posts with feed posts (legacy - to be removed)
    // val combinedPosts = userPosts + posts

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

        // Tab Navigation at the top
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
                                3 -> Icons.Default.PlayArrow
                                4 -> Icons.Default.Explore
                                5 -> Icons.Default.Event
                                else -> Icons.Default.Home
                            },
                            contentDescription = title,
                            modifier = Modifier.size(20.dp)
                        )
                    }
                )
            }
        }

        // Tab content based on selectedTabIndex (with full height)
        Box(modifier = Modifier.weight(1f)) {
            when (selectedTabIndex) {
                0 -> FriendsTab(
                    posts = allPosts.take(10), // Show first 10 unified posts
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
                    onAddComment = handleAddComment,
                    stories = stories,
                    onStoryClick = handleStoryClick,
                    onCreatePostClick = { showCreatePost = true }
                )
                1 -> FeedTab(
                    posts = allPosts, // All 42 post types
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
                    onAddComment = handleAddComment,
                    stories = stories,
                    onStoryClick = handleStoryClick,
                    onCreatePostClick = { showCreatePost = true }
                )
                2 -> UnifiedPostsTab(
                    posts = allPosts,
                    modifier = Modifier.fillMaxSize()
                )
                3 -> DiscoverTab(
                    posts = allPosts.filter { it.user.displayName?.contains("Explorer", ignoreCase = true) == true || it.user.displayName?.contains("Cultural", ignoreCase = true) == true },
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
                4 -> EventsTab(
                    events = events,
                    posts = allPosts.filter { it.content.text?.contains("event") == true },
                    eventCategories = eventCategories,
                    eventPosts = eventPosts,
                    likedPosts = likedPosts,
                    bookmarkedPosts = bookmarkedPosts,
                    followingUsers = followingUsers,
                    commentsOpen = commentsOpen,
                    newComment = newComment,
                    postComments = postComments,
                    rsvpEvents = rsvpEvents,
                    onLike = handleLike,
                    onBookmark = handleBookmark,
                    onFollow = handleFollow,
                    onComment = handleComment,
                    onShare = onShare,
                    onHashtagClick = handleHashtagClick,
                    onCommentTextChange = { newComment = it },
                    onAddComment = handleAddComment,
                    onEventClick = onEventClick,
                    onCategoryClick = onCategoryClick,
                    onRSVPClick = { event ->
                        println("🎫 RSVP toggled for: ${event.title}")
                        rsvpEvents = if (rsvpEvents.contains(event.id)) {
                            rsvpEvents - event.id
                        } else {
                            rsvpEvents + event.id
                        }
                    }
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

// Unified Posts Tab - Shows all 42 post types using PostRenderer
@Composable
private fun UnifiedPostsTab(
    posts: List<Post>,
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.md),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        if (posts.isEmpty()) {
            item {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(TchatSpacing.xl),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                    ) {
                        CircularProgressIndicator(
                            color = TchatColors.primary
                        )
                        Text(
                            text = "Loading all post types...",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
        } else {
            item {
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.md)
                    ) {
                        Text(
                            text = "🎉 All 42 Post Types",
                            style = MaterialTheme.typography.headlineSmall,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.onSurface
                        )
                        Text(
                            text = "Showing ${posts.size} posts with unified PostRenderer",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            items(posts) { post ->
                PostRenderer(
                    post = post,
                    postRepository = MockPostRepository(),
                    sharingService = MockSharingService(),
                    navigationService = MockNavigationService(),
                    onPostClick = { },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        }
    }
}

// Continue with Enhanced Tab Components...

@Composable
private fun FriendsTab(
    posts: List<Post>,
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
    stories: List<StoryItem> = emptyList(),
    onStoryClick: (StoryItem) -> Unit = {},
    onCreatePostClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Stories Row at the top
        item {
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
                contentPadding = PaddingValues(horizontal = TchatSpacing.sm),
                modifier = Modifier.padding(vertical = TchatSpacing.sm)
            ) {
                items(stories) { story ->
                    StoryItemCard(
                        story = story,
                        onClick = { onStoryClick(story) }
                    )
                }
            }
        }

        // Create Post Card
        item {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable { onCreatePostClick() },
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
                        IconButton(onClick = { onCreatePostClick() }) {
                            Icon(Icons.Default.PhotoCamera, "Photo", tint = TchatColors.primary)
                        }
                        IconButton(onClick = { onCreatePostClick() }) {
                            Icon(Icons.Default.LocationOn, "Location", tint = TchatColors.primary)
                        }
                    }
                }
            }
        }
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
            PostRenderer(
                post = post,
                postRepository = MockPostRepository(),
                sharingService = MockSharingService(),
                navigationService = MockNavigationService(),
                onPostClick = { },
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
}

@Composable
private fun FeedTab(
    posts: List<Post>,
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
    stories: List<StoryItem> = emptyList(),
    onStoryClick: (StoryItem) -> Unit = {},
    onCreatePostClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
        contentPadding = PaddingValues(TchatSpacing.md)
    ) {
        // Stories Row at the top
        item {
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
                contentPadding = PaddingValues(horizontal = TchatSpacing.sm),
                modifier = Modifier.padding(vertical = TchatSpacing.sm)
            ) {
                items(stories) { story ->
                    StoryItemCard(
                        story = story,
                        onClick = { onStoryClick(story) }
                    )
                }
            }
        }

        // Create Post Card
        item {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable { onCreatePostClick() },
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
                        IconButton(onClick = { onCreatePostClick() }) {
                            Icon(Icons.Default.PhotoCamera, "Photo", tint = TchatColors.primary)
                        }
                        IconButton(onClick = { onCreatePostClick() }) {
                            Icon(Icons.Default.LocationOn, "Location", tint = TchatColors.primary)
                        }
                    }
                }
            }
        }
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
            PostRenderer(
                post = post,
                postRepository = MockPostRepository(),
                sharingService = MockSharingService(),
                navigationService = MockNavigationService(),
                onPostClick = { },
                modifier = Modifier.fillMaxWidth()
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
    posts: List<Post>,
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
            PostRenderer(
                post = post,
                postRepository = MockPostRepository(),
                sharingService = MockSharingService(),
                navigationService = MockNavigationService(),
                onPostClick = { },
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
}

@Composable
private fun EventsTab(
    events: List<EventItem>,
    posts: List<Post>,
    eventCategories: List<Pair<String, Int>>,
    eventPosts: List<PostItem>,
    likedPosts: Set<String>,
    bookmarkedPosts: Set<String>,
    followingUsers: Set<String>,
    commentsOpen: String?,
    newComment: String,
    postComments: Map<String, List<CommentItem>>,
    rsvpEvents: Set<String>,
    onLike: (String) -> Unit,
    onBookmark: (String) -> Unit,
    onFollow: (String) -> Unit,
    onComment: (String) -> Unit,
    onShare: (PostItem) -> Unit,
    onHashtagClick: (String) -> Unit,
    onCommentTextChange: (String) -> Unit,
    onAddComment: (String) -> Unit,
    onEventClick: (String) -> Unit,
    onCategoryClick: (String, String) -> Unit,
    onRSVPClick: (EventItem) -> Unit,
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
                    // Update event RSVP status based on local state
                    val updatedEvent = event.copy(isAttending = rsvpEvents.contains(event.id))

                    EventCard(
                        event = updatedEvent,
                        modifier = Modifier.weight(1f),
                        onClick = { clickedEvent ->
                            println("🎭 Event clicked: ${clickedEvent.title} at ${clickedEvent.location}")
                            onEventClick(clickedEvent.id)
                        },
                        onRSVPClick = onRSVPClick
                    )
                }
                if (eventPair.size == 1) {
                    Spacer(modifier = Modifier.weight(1f))
                }
            }
        }

        // Event Categories
        item {
            EventCategoriesCard(eventCategories = eventCategories)
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

        items(eventPosts) { postItem ->
            // Convert PostItem to Post for PostRenderer
            val post = Post(
                id = postItem.id,
                type = PostType.TEXT,
                user = PostUser(
                    id = postItem.author.id,
                    username = postItem.author.username,
                    displayName = postItem.author.name,
                    avatarUrl = postItem.author.avatar,
                    isVerified = postItem.author.isVerified
                ),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = postItem.content,
                    hashtags = postItem.tags ?: emptyList(),
                    mentions = emptyList(),
                    location = postItem.location
                ),
                interactions = PostInteractions(
                    reactions = emptyList(),
                    comments = emptyList(),
                    shares = emptyList(),
                    views = 0,
                    isLiked = postItem.isLiked ?: false,
                    isBookmarked = false
                ),
                createdAt = postItem.timestamp
            )

            PostRenderer(
                post = post,
                postRepository = MockPostRepository(),
                sharingService = MockSharingService(),
                navigationService = MockNavigationService(),
                onPostClick = { },
                modifier = Modifier.fillMaxWidth()
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
    modifier: Modifier = Modifier,
    onClick: (EventItem) -> Unit = {},
    onRSVPClick: (EventItem) -> Unit = {}
) {
    Card(
        modifier = modifier
            .clickable {
                println("🎪 Event card clicked: ${event.title}")
                onClick(event)
            },
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
                        onClick = {
                            println("🎫 RSVP clicked for: ${event.title}")
                            onRSVPClick(event)
                        },
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
    eventCategories: List<Pair<String, Int>>,
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
                // Use loaded event categories or fallback to default
                val categoriesToShow = remember(eventCategories) {
                    if (eventCategories.isNotEmpty()) {
                        eventCategories.take(2).map { (name, count) ->
                            val icon = when (name) {
                                "Music" -> Icons.Default.LibraryMusic
                                "Food" -> Icons.Default.Restaurant
                                "Technology" -> Icons.Default.Computer
                                "Arts & Culture" -> Icons.Default.Palette
                                else -> Icons.Default.Category
                            }
                            Triple(name, icon, "$count events")
                        }
                    } else {
                        listOf(
                            Triple("Music", Icons.Default.LibraryMusic, "15 events"),
                            Triple("Food", Icons.Default.Restaurant, "23 events")
                        )
                    }
                }

                categoriesToShow.forEach { (category, icon, count) ->
                    Card(
                        modifier = Modifier
                            .weight(1f)
                            .clickable {
                                println("🔍 Category clicked: $category")
                                // TODO: Implement category filtering
                            },
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
