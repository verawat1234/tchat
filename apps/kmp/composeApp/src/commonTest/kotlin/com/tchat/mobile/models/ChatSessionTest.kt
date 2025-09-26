package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * T043: Comprehensive Chat Group Model Tests
 * Tests all aspects of ChatSession group functionality including permissions,
 * participant management, display utilities, and group-specific features.
 */
class ChatSessionTest {

    // Test data
    private val testUser1 = ChatParticipant(
        id = "user1",
        name = "Alice Admin",
        role = ChatRole.OWNER,
        status = ParticipantStatus.ONLINE,
        joinedAt = "2024-01-15T10:00:00Z"
    )

    private val testUser2 = ChatParticipant(
        id = "user2",
        name = "Bob Member",
        role = ChatRole.MEMBER,
        status = ParticipantStatus.ONLINE,
        joinedAt = "2024-01-15T10:30:00Z"
    )

    private val testUser3 = ChatParticipant(
        id = "user3",
        name = "Charlie Moderator",
        role = ChatRole.MODERATOR,
        status = ParticipantStatus.OFFLINE,
        joinedAt = "2024-01-15T11:00:00Z"
    )

    private val testUser4 = ChatParticipant(
        id = "user4",
        name = "Diana Member",
        role = ChatRole.MEMBER,
        status = ParticipantStatus.AWAY,
        joinedAt = "2024-01-15T11:30:00Z"
    )

    private fun createTestGroupChat(): ChatSession {
        return ChatSession(
            id = "group123",
            name = "Test Group",
            type = ChatType.GROUP,
            participants = listOf(testUser1, testUser2, testUser3, testUser4),
            lastMessage = MessagePreview(
                id = "msg1",
                content = "Hello group!",
                senderId = "user1",
                senderName = "Alice Admin",
                timestamp = "2024-01-15T15:30:00Z"
            ),
            unreadCount = 5,
            isPinned = false,
            isMuted = false,
            isArchived = false,
            isBlocked = false,
            metadata = ChatMetadata(
                description = "Test group for project discussion",
                maxParticipants = 5000,
                tags = listOf("project", "team")
            ),
            permissions = ChatPermissions(
                canAddMembers = true,
                canRemoveMembers = false,
                canEditInfo = false,
                canPinMessages = false,
                canDeleteMessages = false,
                requireApproval = false
            ),
            createdAt = "2024-01-15T10:00:00Z",
            updatedAt = "2024-01-15T15:30:00Z"
        )
    }

    private fun createTestDirectMessage(): ChatSession {
        return ChatSession(
            id = "dm123",
            name = null, // Direct messages don't have names
            type = ChatType.DIRECT,
            participants = listOf(testUser1, testUser2),
            lastMessage = MessagePreview(
                id = "msg2",
                content = "Hey there!",
                senderId = "user2",
                senderName = "Bob Member",
                timestamp = "2024-01-15T14:00:00Z"
            ),
            unreadCount = 2,
            isPinned = true,
            isMuted = false,
            isArchived = false,
            isBlocked = false,
            metadata = ChatMetadata(),
            permissions = ChatPermissions(),
            createdAt = "2024-01-15T09:00:00Z",
            updatedAt = "2024-01-15T14:00:00Z"
        )
    }

    private fun createTestChannel(): ChatSession {
        return ChatSession(
            id = "channel123",
            name = "Announcements",
            type = ChatType.CHANNEL,
            participants = listOf(testUser1, testUser2, testUser3, testUser4),
            metadata = ChatMetadata(
                description = "Company announcements",
                isPublic = true,
                maxParticipants = 200000
            ),
            permissions = ChatPermissions(
                canAddMembers = true,
                canRemoveMembers = false,
                canEditInfo = false,
                canPinMessages = true,
                canDeleteMessages = false,
                canSendMessages = false, // Only admins can message in channels
                requireApproval = false
            ),
            createdAt = "2024-01-10T10:00:00Z",
            updatedAt = "2024-01-15T15:30:00Z"
        )
    }

    @Test
    fun testChatTypeIdentification() {
        val groupChat = createTestGroupChat()
        val directMessage = createTestDirectMessage()
        val channel = createTestChannel()

        // Test group identification
        assertTrue(groupChat.isGroup())
        assertFalse(groupChat.isDirectMessage())

        // Test direct message identification
        assertTrue(directMessage.isDirectMessage())
        assertFalse(directMessage.isGroup())

        // Test channel identification
        assertTrue(channel.isGroup()) // Channels are considered groups
        assertFalse(channel.isDirectMessage())
    }

    @Test
    fun testGetOtherParticipant() {
        val groupChat = createTestGroupChat()
        val directMessage = createTestDirectMessage()

        // Group chat should return null (multiple participants)
        assertNull(groupChat.getOtherParticipant("user1"))

        // Direct message should return the other participant
        val otherParticipant = directMessage.getOtherParticipant("user1")
        assertNotNull(otherParticipant)
        assertEquals("user2", otherParticipant.id)
        assertEquals("Bob Member", otherParticipant.name)

        // Test with second user
        val otherParticipant2 = directMessage.getOtherParticipant("user2")
        assertNotNull(otherParticipant2)
        assertEquals("user1", otherParticipant2.id)
        assertEquals("Alice Admin", otherParticipant2.name)
    }

    @Test
    fun testGetDisplayName() {
        val groupChat = createTestGroupChat()
        val directMessage = createTestDirectMessage()
        val channel = createTestChannel()

        // Group with explicit name
        assertEquals("Test Group", groupChat.getDisplayName("user1"))

        // Direct message should show other participant's name
        assertEquals("Bob Member", directMessage.getDisplayName("user1"))
        assertEquals("Alice Admin", directMessage.getDisplayName("user2"))

        // Channel with explicit name
        assertEquals("Announcements", channel.getDisplayName("user1"))

        // Test group without name (should show participant names)
        val unnamedGroup = groupChat.copy(name = null)
        val displayName = unnamedGroup.getDisplayName("user1")
        assertTrue(displayName.contains("Bob Member"))
        assertTrue(displayName.contains("Charlie Moderator"))

        // Test large group display (more than 3 participants)
        val largeGroupParticipants = listOf(
            testUser1, testUser2, testUser3, testUser4,
            ChatParticipant("user5", "Eve", role = ChatRole.MEMBER, joinedAt = "2024-01-15T12:00:00Z"),
            ChatParticipant("user6", "Frank", role = ChatRole.MEMBER, joinedAt = "2024-01-15T12:30:00Z")
        )
        val largeGroup = groupChat.copy(name = null, participants = largeGroupParticipants)
        val largeDisplayName = largeGroup.getDisplayName("user1")
        assertTrue(largeDisplayName.contains("and"))
        assertTrue(largeDisplayName.contains("others"))
    }

    @Test
    fun testCanUserPerformAction() {
        val groupChat = createTestGroupChat()

        // Test owner permissions (user1 - OWNER)
        assertTrue(groupChat.canUserPerformAction("user1", ChatAction.SEND_MESSAGE))
        assertTrue(groupChat.canUserPerformAction("user1", ChatAction.ADD_MEMBERS))
        assertTrue(groupChat.canUserPerformAction("user1", ChatAction.REMOVE_MEMBERS))
        assertTrue(groupChat.canUserPerformAction("user1", ChatAction.DELETE_MESSAGES))

        // Test member permissions (user2 - MEMBER)
        assertTrue(groupChat.canUserPerformAction("user2", ChatAction.SEND_MESSAGE))
        assertFalse(groupChat.canUserPerformAction("user2", ChatAction.ADD_MEMBERS)) // Only admins can add
        assertFalse(groupChat.canUserPerformAction("user2", ChatAction.REMOVE_MEMBERS)) // Only admins
        assertFalse(groupChat.canUserPerformAction("user2", ChatAction.DELETE_MESSAGES)) // Only admins

        // Test moderator permissions (user3 - MODERATOR)
        assertTrue(groupChat.canUserPerformAction("user3", ChatAction.SEND_MESSAGE))
        assertFalse(groupChat.canUserPerformAction("user3", ChatAction.ADD_MEMBERS)) // Only owner/admin
        assertFalse(groupChat.canUserPerformAction("user3", ChatAction.REMOVE_MEMBERS)) // Only owner/admin
        assertTrue(groupChat.canUserPerformAction("user3", ChatAction.DELETE_MESSAGES)) // Moderators can delete
        assertTrue(groupChat.canUserPerformAction("user3", ChatAction.PIN_MESSAGES)) // Moderators can pin

        // Test non-participant
        assertFalse(groupChat.canUserPerformAction("user999", ChatAction.SEND_MESSAGE))
    }

    @Test
    fun testChannelPermissions() {
        val channel = createTestChannel()

        // In channels, only admins can send messages by default
        assertTrue(channel.canUserPerformAction("user1", ChatAction.SEND_MESSAGE)) // Owner can message
        assertFalse(channel.canUserPerformAction("user2", ChatAction.SEND_MESSAGE)) // Members cannot message
        assertFalse(channel.canUserPerformAction("user3", ChatAction.SEND_MESSAGE)) // Even moderators cannot message

        // But they can still receive and read
        assertTrue(channel.participants.any { it.id == "user2" })
        assertTrue(channel.participants.any { it.id == "user3" })
    }

    @Test
    fun testGetActiveParticipants() {
        val groupChat = createTestGroupChat()

        val activeParticipants = groupChat.getActiveParticipants()

        // Should include ONLINE and AWAY participants
        assertEquals(3, activeParticipants.size)
        assertTrue(activeParticipants.any { it.id == "user1" && it.status == ParticipantStatus.ONLINE })
        assertTrue(activeParticipants.any { it.id == "user2" && it.status == ParticipantStatus.ONLINE })
        assertTrue(activeParticipants.any { it.id == "user4" && it.status == ParticipantStatus.AWAY })

        // Should not include OFFLINE participants
        assertFalse(activeParticipants.any { it.id == "user3" && it.status == ParticipantStatus.OFFLINE })
    }

    @Test
    fun testGetOnlineCount() {
        val groupChat = createTestGroupChat()

        val onlineCount = groupChat.getOnlineCount()

        // Should count only ONLINE participants (user1 and user2)
        assertEquals(2, onlineCount)
    }

    @Test
    fun testChatSessionState() {
        val groupChat = createTestGroupChat()

        // Test unread messages
        assertTrue(groupChat.hasUnreadMessages())
        assertEquals(5, groupChat.unreadCount)

        // Test activity
        assertTrue(groupChat.isActive()) // Has lastActivityAt and not archived/blocked
        assertFalse(groupChat.isArchived)
        assertFalse(groupChat.isBlocked)

        // Test needs attention (has unread and not muted/archived)
        assertTrue(groupChat.needsAttention())

        // Test archived group
        val archivedGroup = groupChat.copy(isArchived = true)
        assertFalse(archivedGroup.isActive())
        assertFalse(archivedGroup.needsAttention())

        // Test muted group
        val mutedGroup = groupChat.copy(isMuted = true)
        assertFalse(mutedGroup.needsAttention()) // Muted groups don't need attention

        // Test no unread messages
        val readGroup = groupChat.copy(unreadCount = 0)
        assertFalse(readGroup.hasUnreadMessages())
        assertFalse(readGroup.needsAttention())
    }

    @Test
    fun testParticipantUtilities() {
        // Test isOnline
        assertTrue(testUser1.isOnline())
        assertTrue(testUser2.isOnline())
        assertFalse(testUser3.isOnline()) // OFFLINE
        assertFalse(testUser4.isOnline()) // AWAY

        // Test canModerate
        assertTrue(testUser1.canModerate()) // OWNER
        assertFalse(testUser2.canModerate()) // MEMBER
        assertTrue(testUser3.canModerate()) // MODERATOR
        assertFalse(testUser4.canModerate()) // MEMBER

        // Test display status
        assertEquals("Online", testUser1.getDisplayStatus())
        assertEquals("Online", testUser2.getDisplayStatus())
        assertEquals("Offline", testUser3.getDisplayStatus())
        assertEquals("Away", testUser4.getDisplayStatus())

        // Test participant with last seen
        val offlineUser = testUser3.copy(lastSeen = "5 minutes ago")
        assertEquals("Last seen 5 minutes ago", offlineUser.getDisplayStatus())

        // Test invisible status (should show as offline for privacy)
        val invisibleUser = testUser1.copy(status = ParticipantStatus.INVISIBLE)
        assertEquals("Offline", invisibleUser.getDisplayStatus())
    }

    @Test
    fun testChatSessionListUtilities() {
        val groupChat = createTestGroupChat()
        val directMessage = createTestDirectMessage()
        val archivedGroup = groupChat.copy(id = "archived", isArchived = true)
        val mutedGroup = groupChat.copy(id = "muted", isMuted = true, unreadCount = 3)
        val noUnreadGroup = groupChat.copy(id = "read", unreadCount = 0)

        val allSessions = listOf(groupChat, directMessage, archivedGroup, mutedGroup, noUnreadGroup)

        // Test filterActive
        val activeSessions = allSessions.filterActive()
        assertEquals(4, activeSessions.size) // All except archived
        assertFalse(activeSessions.any { it.isArchived })

        // Test filterUnread
        val unreadSessions = allSessions.filterUnread()
        assertEquals(3, unreadSessions.size) // groupChat, directMessage, mutedGroup
        assertTrue(unreadSessions.all { it.unreadCount > 0 })

        // Test sortByActivity (pinned first, then by timestamp)
        val sortedByActivity = allSessions.sortByActivity()
        assertEquals(directMessage.id, sortedByActivity[0].id) // Pinned first
        // Then by last message timestamp (most recent first)

        // Test sortByUnread
        val sortedByUnread = allSessions.sortByUnread()
        assertEquals(groupChat.id, sortedByUnread[0].id) // Highest unread count (5)
        assertEquals(mutedGroup.id, sortedByUnread[1].id) // Next highest (3)
        assertEquals(directMessage.id, sortedByUnread[2].id) // Then (2)
    }

    @Test
    fun testSearchFunctionality() {
        val groupChat = createTestGroupChat()
        val directMessage = createTestDirectMessage()
        val channel = createTestChannel()

        val allSessions = listOf(groupChat, directMessage, channel)

        // Test search by group name
        val groupResults = allSessions.searchByName("Test", "user1")
        assertEquals(1, groupResults.size)
        assertEquals(groupChat.id, groupResults[0].id)

        // Test search by participant name
        val participantResults = allSessions.searchByName("Bob", "user1")
        assertEquals(2, participantResults.size) // Both groupChat and directMessage have Bob

        // Test search by message content
        val messageResults = allSessions.searchByName("Hello", "user1")
        assertEquals(1, messageResults.size) // groupChat has "Hello group!" message
        assertEquals(groupChat.id, messageResults[0].id)

        // Test case-insensitive search
        val caseResults = allSessions.searchByName("test", "user1")
        assertEquals(1, caseResults.size)
        assertEquals(groupChat.id, caseResults[0].id)

        // Test empty query returns all
        val emptyResults = allSessions.searchByName("", "user1")
        assertEquals(3, emptyResults.size)

        // Test no matches
        val noResults = allSessions.searchByName("NonExistent", "user1")
        assertEquals(0, noResults.size)
    }

    @Test
    fun testMessageTypesEnumeration() {
        // Test all message types are properly defined
        val allTypes = MessageType.values()
        assertTrue(allTypes.contains(MessageType.TEXT))
        assertTrue(allTypes.contains(MessageType.IMAGE))
        assertTrue(allTypes.contains(MessageType.VIDEO))
        assertTrue(allTypes.contains(MessageType.AUDIO))
        assertTrue(allTypes.contains(MessageType.FILE))
        assertTrue(allTypes.contains(MessageType.LOCATION))
        assertTrue(allTypes.contains(MessageType.CONTACT))
        assertTrue(allTypes.contains(MessageType.STICKER))
        assertTrue(allTypes.contains(MessageType.GIF))
        assertTrue(allTypes.contains(MessageType.POLL))
        assertTrue(allTypes.contains(MessageType.EVENT))
        assertTrue(allTypes.contains(MessageType.SYSTEM))
        assertTrue(allTypes.contains(MessageType.DELETED))

        // Extended message types for Phase D
        assertTrue(allTypes.contains(MessageType.EMBED))
        assertTrue(allTypes.contains(MessageType.EVENT_MESSAGE))
        assertTrue(allTypes.contains(MessageType.FORM))
        assertTrue(allTypes.contains(MessageType.LOCATION_MESSAGE))
        assertTrue(allTypes.contains(MessageType.PAYMENT))
        assertTrue(allTypes.contains(MessageType.FILE_MESSAGE))
    }

    @Test
    fun testChatRolesHierarchy() {
        // Test role hierarchy and permissions
        val roles = ChatRole.values()
        assertTrue(roles.contains(ChatRole.OWNER))
        assertTrue(roles.contains(ChatRole.ADMIN))
        assertTrue(roles.contains(ChatRole.MODERATOR))
        assertTrue(roles.contains(ChatRole.MEMBER))
        assertTrue(roles.contains(ChatRole.GUEST))
        assertTrue(roles.contains(ChatRole.BOT))

        // Test role-based permissions
        val ownerParticipant = testUser1.copy(role = ChatRole.OWNER)
        val adminParticipant = testUser1.copy(role = ChatRole.ADMIN)
        val moderatorParticipant = testUser1.copy(role = ChatRole.MODERATOR)
        val memberParticipant = testUser1.copy(role = ChatRole.MEMBER)
        val guestParticipant = testUser1.copy(role = ChatRole.GUEST)

        assertTrue(ownerParticipant.canModerate())
        assertTrue(adminParticipant.canModerate())
        assertTrue(moderatorParticipant.canModerate())
        assertFalse(memberParticipant.canModerate())
        assertFalse(guestParticipant.canModerate())
    }

    @Test
    fun testGroupMetadata() {
        val groupChat = createTestGroupChat()

        // Test metadata properties
        assertEquals("Test group for project discussion", groupChat.metadata.description)
        assertEquals(5000, groupChat.metadata.maxParticipants)
        assertEquals(listOf("project", "team"), groupChat.metadata.tags)
        assertEquals("en", groupChat.metadata.language)
        assertEquals("UTC", groupChat.metadata.timezone)
        assertFalse(groupChat.metadata.isPublic)
        assertNull(groupChat.metadata.inviteLink)
        assertFalse(groupChat.metadata.autoDeleteMessages)
        assertFalse(groupChat.metadata.encryptionEnabled)
        assertTrue(groupChat.metadata.backupEnabled)

        // Test channel metadata (public)
        val channel = createTestChannel()
        assertTrue(channel.metadata.isPublic)
        assertEquals(200000, channel.metadata.maxParticipants)
    }

    @Test
    fun testGroupPermissions() {
        val groupChat = createTestGroupChat()

        // Test default group permissions
        assertTrue(groupChat.permissions.canAddMembers)
        assertFalse(groupChat.permissions.canRemoveMembers)
        assertFalse(groupChat.permissions.canEditInfo)
        assertFalse(groupChat.permissions.canPinMessages)
        assertFalse(groupChat.permissions.canDeleteMessages)
        assertTrue(groupChat.permissions.canSendMessages)
        assertTrue(groupChat.permissions.canSendMedia)
        assertTrue(groupChat.permissions.canSendPolls)
        assertTrue(groupChat.permissions.canSendFiles)
        assertFalse(groupChat.permissions.requireApproval)

        // Test channel permissions (more restrictive)
        val channel = createTestChannel()
        assertTrue(channel.permissions.canAddMembers)
        assertFalse(channel.permissions.canRemoveMembers)
        assertFalse(channel.permissions.canEditInfo)
        assertTrue(channel.permissions.canPinMessages)
        assertFalse(channel.permissions.canDeleteMessages)
        assertFalse(channel.permissions.canSendMessages) // Only admins in channels
    }

    @Test
    fun testTypingIndicator() {
        val typing = TypingIndicator(
            userId = "user1",
            userName = "Alice",
            startedAt = "2024-01-15T15:30:00Z",
            expiresAt = "2024-01-15T15:30:30Z"
        )

        assertEquals("user1", typing.userId)
        assertEquals("Alice", typing.userName)
        assertEquals("2024-01-15T15:30:00Z", typing.startedAt)
        assertEquals("2024-01-15T15:30:30Z", typing.expiresAt)
    }

    @Test
    fun testChatSessionState() {
        val state = ChatSessionState(
            sessionId = "group123",
            isLoading = false,
            hasMoreMessages = true,
            isTyping = listOf("user2", "user3"),
            isDraftSaving = false,
            draftMessage = "Hello world...",
            replyToMessage = "msg1",
            editingMessage = null,
            selectedMessages = listOf("msg1", "msg2"),
            searchQuery = "hello",
            scrollToMessage = "msg3",
            networkStatus = "connected",
            lastSync = "2024-01-15T15:30:00Z",
            pendingMessages = listOf("pending1"),
            failedMessages = emptyList()
        )

        assertEquals("group123", state.sessionId)
        assertFalse(state.isLoading)
        assertTrue(state.hasMoreMessages)
        assertEquals(2, state.isTyping.size)
        assertEquals("Hello world...", state.draftMessage)
        assertEquals("msg1", state.replyToMessage)
        assertEquals(2, state.selectedMessages.size)
        assertEquals("hello", state.searchQuery)
        assertEquals("connected", state.networkStatus)
        assertEquals(1, state.pendingMessages.size)
        assertEquals(0, state.failedMessages.size)
    }
}