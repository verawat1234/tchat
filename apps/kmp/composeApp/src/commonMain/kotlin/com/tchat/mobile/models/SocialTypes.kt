package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import com.tchat.mobile.services.SharingService
import com.tchat.mobile.services.MockSharingService
import com.tchat.mobile.services.NavigationService
import com.tchat.mobile.services.MockNavigationService

/**
 * Social Media Platform Models
 * Data models for social interactions, stories, friends, events
 */

@Serializable
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

@Serializable
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

@Serializable
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

@Serializable
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

@Serializable
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
    val type: SocialPostType = SocialPostType.TEXT,
    val source: PostSource = PostSource.FOLLOWING,
    val isLiked: Boolean = false
)

@Serializable
data class CommentItem(
    val id: String,
    val user: UserItem,
    val text: String,
    val timestamp: String,
    val likes: Int,
    val isLiked: Boolean
)

enum class SocialPostType { TEXT, IMAGE, LIVE, PRODUCT }
enum class PostSource { FOLLOWING, TRENDING, INTEREST, SPONSORED }

data class MockServices(
    val sharingService: SharingService,
    val navigationService: NavigationService
)

// Mock data generation functions
object SocialMockData {

    fun getDummyStories(): List<StoryItem> {
        return listOf(
            StoryItem(
                id = "story1",
                author = UserItem(
                    id = "user1",
                    name = "Alice Johnson",
                    username = "alice_j",
                    avatar = "https://i.pravatar.cc/150?img=1",
                    isVerified = true,
                    isOnline = true
                ),
                preview = "Check out my new adventure!",
                content = "Having an amazing time at the beach! 🏖️",
                timestamp = "2 hours ago",
                isViewed = false,
                isLive = false,
                expiresAt = "22 hours remaining"
            ),
            StoryItem(
                id = "story2",
                author = UserItem(
                    id = "user2",
                    name = "Bob Smith",
                    username = "bob_smith",
                    avatar = "https://i.pravatar.cc/150?img=2",
                    isOnline = false,
                    lastSeen = "1 hour ago"
                ),
                preview = "Cooking something delicious",
                content = "Made the perfect pasta tonight! 🍝",
                timestamp = "4 hours ago",
                isViewed = true,
                expiresAt = "20 hours remaining"
            ),
            StoryItem(
                id = "story3",
                author = UserItem(
                    id = "user3",
                    name = "Carol Davis",
                    username = "carol_d",
                    avatar = "https://i.pravatar.cc/150?img=3",
                    isOnline = true,
                    isVerified = true
                ),
                preview = "Live from the concert!",
                content = "Amazing performance tonight! 🎵",
                timestamp = "Live",
                isViewed = false,
                isLive = true,
                expiresAt = "Live now"
            )
        )
    }

    fun getDummyFriends(): List<FriendItem> {
        return listOf(
            FriendItem(
                id = "friend1",
                name = "Emma Wilson",
                username = "emma_w",
                avatar = "https://i.pravatar.cc/150?img=4",
                isOnline = true,
                isFollowing = true,
                mutualFriends = 12,
                status = "Just posted a photo"
            ),
            FriendItem(
                id = "friend2",
                name = "David Brown",
                username = "david_b",
                avatar = "https://i.pravatar.cc/150?img=5",
                isOnline = false,
                isFollowing = true,
                mutualFriends = 8,
                status = "Active 2 hours ago"
            ),
            FriendItem(
                id = "friend3",
                name = "Lisa Garcia",
                username = "lisa_g",
                avatar = "https://i.pravatar.cc/150?img=6",
                isOnline = true,
                isFollowing = false,
                mutualFriends = 5,
                status = "Typing..."
            ),
            FriendItem(
                id = "friend4",
                name = "Mike Johnson",
                username = "mike_j",
                avatar = "https://i.pravatar.cc/150?img=7",
                isOnline = false,
                isFollowing = true,
                mutualFriends = 15,
                status = "Active yesterday"
            )
        )
    }

    fun getDummyEvents(): List<EventItem> {
        return listOf(
            EventItem(
                id = "event1",
                title = "Summer Music Festival",
                description = "Join us for an amazing day of live music, food, and fun!",
                date = "July 15, 2024",
                location = "Central Park, New York",
                price = "$25",
                imageUrl = "https://picsum.photos/400/200?random=1",
                attendeesCount = 1250,
                category = "Music",
                isAttending = true
            ),
            EventItem(
                id = "event2",
                title = "Tech Conference 2024",
                description = "Discover the latest trends in technology and innovation.",
                date = "August 22, 2024",
                location = "Convention Center, San Francisco",
                price = "$75",
                imageUrl = "https://picsum.photos/400/200?random=2",
                attendeesCount = 850,
                category = "Technology",
                isAttending = false
            ),
            EventItem(
                id = "event3",
                title = "Art Gallery Opening",
                description = "Explore contemporary art from local and international artists.",
                date = "June 30, 2024",
                location = "Downtown Gallery, Los Angeles",
                price = "Free",
                imageUrl = "https://picsum.photos/400/200?random=3",
                attendeesCount = 320,
                category = "Art",
                isAttending = false
            ),
            EventItem(
                id = "event4",
                title = "Food & Wine Festival",
                description = "Taste the best cuisine and wines from around the world.",
                date = "September 5, 2024",
                location = "Waterfront Park, Seattle",
                price = "$45",
                imageUrl = "https://picsum.photos/400/200?random=4",
                attendeesCount = 2100,
                category = "Food & Drink",
                isAttending = true
            )
        )
    }

    fun rememberMockServices(): MockServices {
        return MockServices(
            sharingService = MockSharingService(),
            navigationService = MockNavigationService()
        )
    }
}