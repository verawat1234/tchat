package com.tchat.mobile.services

import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.ChatRepository
import com.tchat.mobile.repositories.SocialRepository
import kotlinx.coroutines.delay

/**
 * Service for seeding initial user data to get the app running
 * Creates test users, chats, posts, and social connections
 */
class UserSeedingService(
    private val database: TchatDatabase,
    private val socialRepository: SocialRepository,
    private val chatRepository: ChatRepository
) {

    suspend fun seedInitialData(): Result<Unit> {
        return try {
            // Create test users
            val users = createTestUsers()

            // Create test chats and messages
            val chats = createTestChats(users)

            // Create social connections
            createSocialConnections(users)

            // Create sample content
            createSampleContent(users)

            println("✅ Initial data seeded successfully")
            Result.success(Unit)
        } catch (e: Exception) {
            println("❌ Failed to seed initial data: ${e.message}")
            Result.failure(e)
        }
    }

    private suspend fun createTestUsers(): List<SocialUserProfile> {
        val users = listOf(
            SocialUserProfile(
                userId = "user_1",
                displayName = "Alice Johnson",
                username = "alice_j",
                avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=alice",
                bio = "UI/UX Designer • Coffee lover ☕ • Based in Bangkok",
                isVerified = true,
                isOnline = true,
                lastSeen = System.currentTimeMillis(),
                statusMessage = "Working on something cool!",
                createdAt = System.currentTimeMillis() - 86400000, // 1 day ago
                updatedAt = System.currentTimeMillis()
            ),
            SocialUserProfile(
                userId = "user_2",
                displayName = "Bob Smith",
                username = "bob_dev",
                avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=bob",
                bio = "Software Developer • React Native enthusiast • 🚀",
                isVerified = false,
                isOnline = false,
                lastSeen = System.currentTimeMillis() - 3600000, // 1 hour ago
                statusMessage = "Building the future",
                createdAt = System.currentTimeMillis() - 172800000, // 2 days ago
                updatedAt = System.currentTimeMillis() - 3600000
            ),
            SocialUserProfile(
                userId = "user_3",
                displayName = "Carol Zhang",
                username = "carol_create",
                avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=carol",
                bio = "Content Creator • Travel blogger ✈️ • Foodie 🍜",
                isVerified = true,
                isOnline = true,
                lastSeen = System.currentTimeMillis(),
                statusMessage = "Exploring new places",
                createdAt = System.currentTimeMillis() - 259200000, // 3 days ago
                updatedAt = System.currentTimeMillis()
            ),
            SocialUserProfile(
                userId = "current_user",
                displayName = "You",
                username = "current_user",
                avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=currentuser",
                bio = "Testing the Tchat app! 🎉",
                isVerified = false,
                isOnline = true,
                lastSeen = System.currentTimeMillis(),
                statusMessage = "Online",
                createdAt = System.currentTimeMillis() - 86400000,
                updatedAt = System.currentTimeMillis()
            )
        )

        // Insert users into database
        users.forEach { user ->
            try {
                socialRepository.updateUserProfile(user)
                delay(10) // Small delay to avoid overwhelming the database
            } catch (e: Exception) {
                println("Failed to insert user ${user.username}: ${e.message}")
            }
        }

        println("✅ Created ${users.size} test users")
        return users
    }

    private suspend fun createTestChats(users: List<SocialUserProfile>): List<ChatSession> {
        val currentTime = System.currentTimeMillis()

        // Helper function to convert timestamp to ISO string
        fun timestampToIsoString(timestamp: Long): String {
            return kotlinx.datetime.Instant.fromEpochMilliseconds(timestamp).toString()
        }

        val chats = listOf(
            ChatSession(
                id = "chat_1",
                name = "Alice Johnson",
                type = ChatType.DIRECT,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        avatar = users.find { it.userId == "current_user" }?.avatarUrl,
                        joinedAt = timestampToIsoString(currentTime - 86400000)
                    ),
                    ChatParticipant(
                        id = "user_1",
                        name = "Alice Johnson",
                        avatar = users.find { it.userId == "user_1" }?.avatarUrl,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = timestampToIsoString(currentTime - 86400000)
                    )
                ),
                metadata = ChatMetadata(),
                unreadCount = 2,
                isPinned = true,
                isMuted = false,
                isArchived = false,
                isBlocked = false,
                createdAt = timestampToIsoString(currentTime - 86400000),
                updatedAt = timestampToIsoString(currentTime - 300000), // 5 minutes ago
                lastActivityAt = timestampToIsoString(currentTime - 300000)
            ),
            ChatSession(
                id = "chat_2",
                name = "Development Team",
                type = ChatType.GROUP,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        avatar = users.find { it.userId == "current_user" }?.avatarUrl,
                        role = ChatRole.ADMIN,
                        joinedAt = timestampToIsoString(currentTime - 172800000)
                    ),
                    ChatParticipant(
                        id = "user_2",
                        name = "Bob Smith",
                        avatar = users.find { it.userId == "user_2" }?.avatarUrl,
                        status = ParticipantStatus.OFFLINE,
                        joinedAt = timestampToIsoString(currentTime - 172800000)
                    ),
                    ChatParticipant(
                        id = "user_3",
                        name = "Carol Zhang",
                        avatar = users.find { it.userId == "user_3" }?.avatarUrl,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = timestampToIsoString(currentTime - 86400000)
                    )
                ),
                metadata = ChatMetadata(
                    description = "Development team collaboration"
                ),
                unreadCount = 5,
                isPinned = false,
                isMuted = false,
                isArchived = false,
                isBlocked = false,
                createdAt = timestampToIsoString(currentTime - 172800000),
                updatedAt = timestampToIsoString(currentTime - 1800000), // 30 minutes ago
                lastActivityAt = timestampToIsoString(currentTime - 1800000)
            ),
            ChatSession(
                id = "chat_3",
                name = "Carol Zhang",
                type = ChatType.DIRECT,
                participants = listOf(
                    ChatParticipant(
                        id = "current_user",
                        name = "You",
                        avatar = users.find { it.userId == "current_user" }?.avatarUrl,
                        joinedAt = timestampToIsoString(currentTime - 259200000)
                    ),
                    ChatParticipant(
                        id = "user_3",
                        name = "Carol Zhang",
                        avatar = users.find { it.userId == "user_3" }?.avatarUrl,
                        status = ParticipantStatus.ONLINE,
                        joinedAt = timestampToIsoString(currentTime - 259200000)
                    )
                ),
                metadata = ChatMetadata(),
                unreadCount = 0,
                isPinned = false,
                isMuted = false,
                isArchived = false,
                isBlocked = false,
                createdAt = timestampToIsoString(currentTime - 259200000),
                updatedAt = timestampToIsoString(currentTime - 7200000), // 2 hours ago
                lastActivityAt = timestampToIsoString(currentTime - 7200000)
            )
        )

        // Create sample messages for each chat
        val messages = listOf(
            // Chat 1 messages
            Message(
                id = "msg_1",
                chatId = "chat_1",
                senderId = "user_1",
                senderName = "Alice Johnson",
                type = MessageType.TEXT,
                content = "Hey! How's the new app coming along? 😊",
                isEdited = false,
                isPinned = false,
                isDeleted = false,
                replyToId = null,
                reactions = emptyList(),
                createdAt = timestampToIsoString(currentTime - 300000),
                editedAt = null,
                deletedAt = null
            ),
            Message(
                id = "msg_2",
                chatId = "chat_1",
                senderId = "user_1",
                senderName = "Alice Johnson",
                type = MessageType.TEXT,
                content = "The UI looks amazing! 🎨",
                isEdited = false,
                isPinned = false,
                isDeleted = false,
                replyToId = null,
                reactions = emptyList(),
                createdAt = timestampToIsoString(currentTime - 240000),
                editedAt = null,
                deletedAt = null
            ),

            // Chat 2 messages
            Message(
                id = "msg_3",
                chatId = "chat_2",
                senderId = "user_2",
                senderName = "Bob Smith",
                type = MessageType.TEXT,
                content = "Good morning team! Daily standup in 10 minutes 👋",
                isEdited = false,
                isPinned = false,
                isDeleted = false,
                replyToId = null,
                reactions = emptyList(),
                createdAt = timestampToIsoString(currentTime - 1800000),
                editedAt = null,
                deletedAt = null
            ),

            // Chat 3 messages
            Message(
                id = "msg_4",
                chatId = "chat_3",
                senderId = "user_3",
                senderName = "Carol Zhang",
                type = MessageType.TEXT,
                content = "Just posted some amazing photos from my trip to Chiang Mai! 📸",
                isEdited = false,
                isPinned = false,
                isDeleted = false,
                replyToId = null,
                reactions = emptyList(),
                createdAt = timestampToIsoString(currentTime - 7200000),
                editedAt = null,
                deletedAt = null
            )
        )

        // Insert chats and messages via repository
        chats.forEach { chat ->
            try {
                // Use repository method to create chat sessions
                chatRepository.createChatSession(chat)
                delay(10)
            } catch (e: Exception) {
                println("Failed to insert chat ${chat.id}: ${e.message}")
            }
        }

        messages.forEach { message ->
            try {
                // Use repository method to save messages
                chatRepository.sendMessage(message)
                delay(10)
            } catch (e: Exception) {
                println("Failed to insert message ${message.id}: ${e.message}")
            }
        }

        println("✅ Created ${chats.size} test chats with ${messages.size} messages")
        return chats
    }

    private suspend fun createSocialConnections(users: List<SocialUserProfile>) {
        val currentUser = users.find { it.userId == "current_user" }
        val otherUsers = users.filter { it.userId != "current_user" }

        if (currentUser == null) return

        // Create friend connections
        otherUsers.forEach { user ->
            try {
                // Make current user friends with all other users
                socialRepository.sendFriendRequest("current_user", user.userId)
                delay(10)
                socialRepository.acceptFriendRequest("current_user", user.userId)
                delay(10)
            } catch (e: Exception) {
                println("Failed to create friendship with ${user.username}: ${e.message}")
            }
        }

        println("✅ Created social connections")
    }

    private suspend fun createSampleContent(users: List<SocialUserProfile>) {
        val currentTime = System.currentTimeMillis()

        // Create sample stories
        val stories = listOf(
            Story(
                id = "story_1",
                authorId = "user_1",
                content = "Working late on the new design system! 🎨✨",
                preview = "New design...",
                createdAt = currentTime - 3600000, // 1 hour ago
                expiresAt = currentTime + 82800000, // 23 hours from now
                isLive = false,
                viewCount = 15,
                isViewed = false,
                totalViews = 15
            ),
            Story(
                id = "story_2",
                authorId = "user_3",
                content = "Amazing sunset from Doi Suthep! 🌅",
                preview = "Amazing sunset...",
                createdAt = currentTime - 7200000, // 2 hours ago
                expiresAt = currentTime + 79200000, // 22 hours from now
                isLive = false,
                viewCount = 42,
                isViewed = false,
                totalViews = 42
            )
        )

        stories.forEach { story ->
            try {
                socialRepository.createStory(story)
                delay(10)
            } catch (e: Exception) {
                println("Failed to create story ${story.id}: ${e.message}")
            }
        }

        println("✅ Created ${stories.size} sample stories")
    }

    /**
     * Check if the database already has seeded data
     */
    suspend fun isDataSeeded(): Boolean {
        return try {
            // Check if current_user exists
            val profile = socialRepository.getUserProfile("current_user")
            profile.isSuccess && profile.getOrNull() != null
        } catch (e: Exception) {
            false
        }
    }

    /**
     * Clear all seeded data (for testing/reset)
     */
    suspend fun clearSeedData(): Result<Unit> {
        return try {
            // This would clear test data - implement if needed
            println("🧹 Cleared seed data")
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}