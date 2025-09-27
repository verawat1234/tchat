package com.tchat.mobile.screens

import androidx.compose.animation.*
import androidx.compose.animation.core.*
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
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.tchat.mobile.components.*
import com.tchat.mobile.components.chat.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.ChatRepository
import com.tchat.mobile.repositories.MockChatRepository
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatDetailScreen(
    chatId: String,
    chatName: String? = null,
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    // Chat repository
    val chatRepository = remember { MockChatRepository() }
    val scope = rememberCoroutineScope()

    // State management
    var messageText by remember { mutableStateOf("") }
    var messages by remember { mutableStateOf<List<Message>>(emptyList()) }
    var typingUsers by remember { mutableStateOf<List<TypingIndicator>>(emptyList()) }
    var isTyping by remember { mutableStateOf(false) }
    var isLoading by remember { mutableStateOf(false) }
    var replyToMessage by remember { mutableStateOf<Message?>(null) }
    var selectedMessage by remember { mutableStateOf<Message?>(null) }
    var showMessageActions by remember { mutableStateOf(false) }
    var showQuickReactions by remember { mutableStateOf(false) }
    var draftMessage by remember { mutableStateOf("") }

    // Observe real-time messages
    val messagesFlow by chatRepository.observeMessages(chatId).collectAsState(initial = emptyList())

    // Update messages state when flow changes
    LaunchedEffect(messagesFlow) {
        messages = messagesFlow
    }

    // Load initial data
    LaunchedEffect(chatId) {

        // Load draft message
        chatRepository.getDraftMessage(chatId).onSuccess { draft ->
            draftMessage = draft ?: ""
            messageText = draftMessage
        }
    }

    // Observe real-time updates
    LaunchedEffect(chatId) {
        chatRepository.observeMessages(chatId).collect { messageList ->
            messages = messageList
        }
    }

    LaunchedEffect(chatId) {
        chatRepository.observeTypingIndicators(chatId).collect { indicators ->
            typingUsers = indicators.filter { it.userId != "current_user" }
        }
    }

    // Save draft messages
    LaunchedEffect(messageText) {
        if (messageText.isNotEmpty()) {
            delay(1000) // Debounce
            chatRepository.saveDraftMessage(chatId, messageText)
        } else {
            chatRepository.clearDraftMessage(chatId)
        }
    }

    // Get actual chat name based on chat ID if not provided
    val actualChatName = chatName ?: "Chat"

    // Chat action states
    var showVideoCallDialog by remember { mutableStateOf(false) }
    var showAudioCallDialog by remember { mutableStateOf(false) }
    var showMoreOptionsMenu by remember { mutableStateOf(false) }
    var showToast by remember { mutableStateOf("") }
    var showConfirmDialog by remember { mutableStateOf(false) }

    // Search functionality states
    var isSearchActive by remember { mutableStateOf(false) }
    var searchQuery by remember { mutableStateOf("") }
    var searchResults by remember { mutableStateOf<List<Message>>(emptyList()) }
    var currentSearchIndex by remember { mutableStateOf(0) }
    var filteredMessages by remember { mutableStateOf(messages) }

    // Search filtering logic
    LaunchedEffect(searchQuery, messages) {
        if (searchQuery.isNotEmpty()) {
            searchResults = messages.filter { message ->
                message.content.contains(searchQuery, ignoreCase = true)
            }
            filteredMessages = searchResults
            currentSearchIndex = if (searchResults.isNotEmpty()) 0 else -1
        } else {
            searchResults = emptyList()
            filteredMessages = messages
            currentSearchIndex = -1
        }
    }

    Scaffold(
        modifier = Modifier.testTag("chat-detail-screen-container"),
        topBar = {
            TopAppBar(
                title = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        // User Avatar
                        Box(
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .background(TchatColors.primary)
                                .testTag("chat-detail-avatar-image"),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = actualChatName.first().toString(),
                                color = TchatColors.onPrimary,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold
                            )
                        }
                        Spacer(modifier = Modifier.width(TchatSpacing.md))
                        Column {
                            Text(
                                text = actualChatName,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.SemiBold,
                                color = TchatColors.onSurface,
                                modifier = Modifier.testTag("chat-detail-header-name")
                            )
                            ChatStatusBar(
                                isOnline = true,
                                lastSeen = "5_minutes_ago",
                                typingUsers = typingUsers,
                                participantCount = 1,
                                isActiveNow = true,
                                customStatus = null,
                                modifier = Modifier.testTag("chat-detail-header-status")
                            )
                        }
                    }
                },
                navigationIcon = {
                    IconButton(
                        onClick = onBackClick,
                        modifier = Modifier.testTag("chat-detail-back-button")
                    ) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                actions = {
                    // Search button
                    IconButton(
                        onClick = {
                            isSearchActive = !isSearchActive
                            if (!isSearchActive) {
                                searchQuery = ""
                            }
                        },
                        modifier = Modifier.testTag("chat-detail-search-button")
                    ) {
                        Icon(
                            if (isSearchActive) Icons.Default.Close else Icons.Default.Search,
                            contentDescription = if (isSearchActive) "Close Search" else "Search Messages",
                            tint = if (isSearchActive) TchatColors.primary else TchatColors.onSurface
                        )
                    }

                    IconButton(
                        onClick = {
                            showVideoCallDialog = true
                        },
                        modifier = Modifier.testTag("chat-detail-video-call-button")
                    ) {
                        Icon(
                            Icons.Default.VideoCall,
                            contentDescription = "Video Call",
                            tint = TchatColors.onSurface
                        )
                    }
                    IconButton(
                        onClick = {
                            showAudioCallDialog = true
                        },
                        modifier = Modifier.testTag("chat-detail-audio-call-button")
                    ) {
                        Icon(
                            Icons.Default.Call,
                            contentDescription = "Audio Call",
                            tint = TchatColors.onSurface
                        )
                    }
                    IconButton(
                        onClick = {
                            showMoreOptionsMenu = true
                        },
                        modifier = Modifier.testTag("chat-detail-more-options-button")
                    ) {
                        Icon(
                            Icons.Default.MoreVert,
                            contentDescription = "More Options",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface
                )
            )
        },
        bottomBar = {
            Column {
                // Typing indicator at bottom
                TypingIndicatorRow(
                    typingUsers = typingUsers,
                    modifier = Modifier.testTag("chat-typing-indicator")
                )

                // Enhanced message input system
                MessageInputSystem(
                    message = messageText,
                    onMessageChange = { newText ->
                        messageText = newText
                        if (newText.isNotEmpty() && !isTyping) {
                            isTyping = true
                            scope.launch {
                                chatRepository.startTyping(chatId)
                                delay(3000) // Stop typing after 3 seconds of inactivity
                                chatRepository.stopTyping(chatId)
                                isTyping = false
                            }
                        }
                    },
                    onSendMessage = { messageType, content, attachments ->
                        scope.launch {
                            isLoading = true
                            try {
                                when {
                                    attachments.isNotEmpty() -> {
                                        chatRepository.sendMediaMessage(
                                            chatId = chatId,
                                            attachments = attachments,
                                            caption = content.takeIf { it.isNotBlank() },
                                            replyToId = replyToMessage?.id
                                        )
                                    }
                                    else -> {
                                        chatRepository.sendTextMessage(
                                            chatId = chatId,
                                            content = content,
                                            replyToId = replyToMessage?.id
                                        )
                                    }
                                }
                                messageText = ""
                                replyToMessage = null
                                chatRepository.clearDraftMessage(chatId)
                                chatRepository.stopTyping(chatId)
                                isTyping = false
                            } catch (e: Exception) {
                                showToast = "Failed to send message: ${e.message}"
                            } finally {
                                isLoading = false
                            }
                        }
                    },
                    replyToMessage = replyToMessage,
                    onCancelReply = { replyToMessage = null },
                    isTyping = isTyping,
                    onTypingStart = {
                        scope.launch {
                            chatRepository.startTyping(chatId)
                        }
                    },
                    onTypingStop = {
                        scope.launch {
                            chatRepository.stopTyping(chatId)
                        }
                    },
                    onVoiceRecordStart = {
                        showToast = "Voice recording started..."
                    },
                    onVoiceRecordStop = {
                        scope.launch {
                            // Mock voice message
                            val mockUrl = "voice://recording_${System.currentTimeMillis()}.mp4"
                            chatRepository.sendVoiceMessage(chatId, mockUrl, 15)
                            showToast = "Voice message sent!"
                        }
                    },
                    onVoiceRecordCancel = {
                        showToast = "Voice recording cancelled"
                    },
                    modifier = Modifier.testTag("enhanced-chat-input")
                )
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .background(TchatColors.background)
        ) {
            // Search overlay
            // TODO: Implement search overlay
            /*
            AnimatedVisibility(
                visible = isSearchActive,
                enter = slideInVertically() + fadeIn(),
                exit = slideOutVertically() + fadeOut()
            ) {
                SearchOverlay(
                    searchQuery = searchQuery,
                    onSearchQueryChange = { searchQuery = it },
                    searchResults = searchResults,
                    currentIndex = currentSearchIndex,
                    onNavigateNext = {
                        if (searchResults.isNotEmpty()) {
                            currentSearchIndex = (currentSearchIndex + 1) % searchResults.size
                        }
                    },
                    onNavigatePrevious = {
                        if (searchResults.isNotEmpty()) {
                            currentSearchIndex = if (currentSearchIndex > 0)
                                currentSearchIndex - 1
                            else
                                searchResults.size - 1
                        }
                    },
                    modifier = Modifier.testTag("search-overlay")
                )
            }
            */

            LazyColumn(
                modifier = Modifier
                    .fillMaxSize()
                    .testTag("chat-detail-messages-list"),
                contentPadding = PaddingValues(TchatSpacing.md),
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
                reverseLayout = false
            ) {
            items(filteredMessages) { message ->
                // TODO: Replace with working message bubble
                MessageRenderer(
                    message = message,
                    isFromMe = message.senderId == "current_user",
                    onLongPress = {
                        selectedMessage = message
                        showMessageActions = true
                    },
                    onReactionClick = { emoji ->
                        scope.launch {
                            val reaction = MessageReaction(
                                emoji = emoji,
                                userId = "current_user",
                                userName = "You",
                                timestamp = System.currentTimeMillis().toString()
                            )
                            chatRepository.addReaction(message.id, reaction)
                        }
                    },
                    onReplyClick = {
                        replyToMessage = message
                    },
                    modifier = Modifier.testTag("enhanced-message-item-${message.id}")
                )
            }

            // Add typing bubble at the end
            if (typingUsers.isNotEmpty()) {
                item {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(horizontal = TchatSpacing.md),
                        horizontalArrangement = Arrangement.Start
                    ) {
                        InlineTypingBubble(
                            typingUsers = typingUsers,
                            modifier = Modifier.testTag("inline-typing-bubble")
                        )
                    }
                }
            }
        }
    }

    // Chat Action Dialogs - TODO: Implement missing dialog components
    /*
    if (showVideoCallDialog) {
        ChatActionDialog(
            modifier = Modifier.testTag("chat-video-call-dialog"),
            title = "Start Video Call",
            message = "Would you like to start a video call with $chatName?",
            icon = Icons.Default.VideoCall,
            onConfirm = {
                showVideoCallDialog = false
                showToast = "Starting video call with $chatName..."
                // Here you would navigate to video call screen or start video call logic
            },
            onDismiss = { showVideoCallDialog = false }
        )
    }

    if (showAudioCallDialog) {
        ChatActionDialog(
            modifier = Modifier.testTag("chat-audio-call-dialog"),
            title = "Start Audio Call",
            message = "Would you like to start an audio call with $chatName?",
            icon = Icons.Default.Call,
            onConfirm = {
                showAudioCallDialog = false
                showToast = "Starting audio call with $chatName..."
                // Here you would navigate to audio call screen or start audio call logic
            },
            onDismiss = { showAudioCallDialog = false }
        )
    }

    if (showMoreOptionsMenu) {
        MoreOptionsDialog(
            modifier = Modifier.testTag("chat-more-options-dialog"),
            chatName = actualChatName,
            onDismiss = { showMoreOptionsMenu = false },
            onClearChat = {
                // TODO: Implement confirm dialog
                messages = emptyList()
                showToast = "Chat cleared successfully"
                showMoreOptionsMenu = false
            },
            onBlockUser = {
                // TODO: Implement confirm dialog
                showToast = "$chatName has been blocked"
                showMoreOptionsMenu = false
            },
            onMuteNotifications = {
                showToast = "Notifications muted for $chatName"
                showMoreOptionsMenu = false
            },
            onExportChat = {
                showToast = "Exporting chat messages..."
                showMoreOptionsMenu = false
            }
        )
    }
    */

    // TODO: Implement confirmation dialog
    /*
    showConfirmDialog?.let { dialogState ->
        ConfirmationDialog(
            modifier = Modifier.testTag("chat-confirmation-dialog"),
            state = dialogState,
            onDismiss = { showConfirmDialog = null }
        )
    }
    */

    // Message action menu
    selectedMessage?.let { message ->
        if (showMessageActions) {
            MessageActionMenu(
                message = message,
                isFromMe = message.senderId == "current_user",
                onDismiss = {
                    showMessageActions = false
                    selectedMessage = null
                },
                onReply = {
                    replyToMessage = message
                    showMessageActions = false
                    selectedMessage = null
                },
                onEdit = {
                    messageText = message.content
                    showMessageActions = false
                    selectedMessage = null
                },
                onDelete = {
                    scope.launch {
                        chatRepository.deleteMessage(message.id)
                        showToast = "Message deleted"
                    }
                    showMessageActions = false
                    selectedMessage = null
                },
                onPin = {
                    scope.launch {
                        if (message.isPinned) {
                            chatRepository.unpinMessage(message.id)
                            showToast = "Message unpinned"
                        } else {
                            chatRepository.pinMessage(message.id)
                            showToast = "Message pinned"
                        }
                    }
                    showMessageActions = false
                    selectedMessage = null
                },
                onCopy = {
                    // Mock copy to clipboard
                    showToast = "Message copied to clipboard"
                    showMessageActions = false
                    selectedMessage = null
                },
                onForward = {
                    // Mock forward
                    showToast = "Forward message feature"
                    showMessageActions = false
                    selectedMessage = null
                },
                onAddReaction = {
                    showQuickReactions = true
                    showMessageActions = false
                },
                modifier = Modifier.testTag("message-action-menu")
            )
        }
    }

    // Quick reaction picker
    if (showQuickReactions && selectedMessage != null) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            QuickReactionPicker(
                onReactionSelect = { emoji ->
                    scope.launch {
                        selectedMessage?.let { message ->
                            val reaction = MessageReaction(
                                emoji = emoji,
                                userId = "current_user",
                                userName = "You",
                                timestamp = System.currentTimeMillis().toString()
                            )
                            chatRepository.addReaction(message.id, reaction)
                        }
                    }
                    showQuickReactions = false
                    selectedMessage = null
                },
                onDismiss = {
                    showQuickReactions = false
                    selectedMessage = null
                },
                modifier = Modifier.testTag("quick-reaction-picker")
            )
        }
    }

    // Toast notifications
    if (showToast.isNotEmpty()) {
        LaunchedEffect(showToast) {
            kotlinx.coroutines.delay(2000)
            showToast = ""
        }
        // Simple toast simulation with Snackbar-like behavior
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.BottomCenter
        ) {
            Card(
                modifier = Modifier.padding(16.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Text(
                    text = showToast,
                    modifier = Modifier.padding(16.dp),
                    color = TchatColors.onSurface
                )
            }
        }
    }
    } // Fixed: Added missing closing brace
}

// Legacy ChatInputBar component - replaced by MessageInputSystem
// Removed as it's no longer used

// Legacy attachment components - replaced by EnhancedAttachmentMenu in MessageInputSystem

// TODO: Fix EnhancedMessageBubble implementation - ALL FUNCTIONS TEMPORARILY REMOVED