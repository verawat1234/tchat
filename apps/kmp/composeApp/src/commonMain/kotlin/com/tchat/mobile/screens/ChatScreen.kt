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
import androidx.compose.ui.unit.sp
import androidx.compose.ui.geometry.Offset
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatScreen(
    onChatClick: (chatId: String, chatName: String) -> Unit = { _, _ -> },
    onSearchClick: () -> Unit = {},
    onQRScannerClick: () -> Unit = {},
    onNotificationsClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var selectedTab by remember { mutableStateOf(0) }
    var personalSearchQuery by remember { mutableStateOf("") }
    var businessSearchQuery by remember { mutableStateOf("") }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(TchatColors.background)
    ) {
        // Enhanced Header Bar like Web
        TopAppBar(
            title = {
                Text(
                    "Messages",
                    fontWeight = FontWeight.Bold,
                    style = MaterialTheme.typography.headlineSmall
                )
            },
            actions = {
                IconButton(onClick = onSearchClick) {
                    Icon(
                        Icons.Filled.Search,
                        "Search",
                        tint = TchatColors.onSurface
                    )
                }
                IconButton(onClick = onQRScannerClick) {
                    Icon(
                        Icons.Filled.QrCode,
                        "QR Scanner",
                        tint = TchatColors.onSurface
                    )
                }
                Box {
                    IconButton(onClick = onNotificationsClick) {
                        Icon(
                            Icons.Filled.Notifications,
                            "Notifications",
                            tint = TchatColors.onSurface
                        )
                    }
                    // Notification badge for unread count
                    Box(
                        modifier = Modifier
                            .size(18.dp)
                            .background(TchatColors.error, CircleShape)
                            .offset(x = 12.dp, y = (-4).dp),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = "3",
                            color = TchatColors.onPrimary,
                            style = MaterialTheme.typography.labelSmall.copy(
                                fontSize = 10.sp
                            )
                        )
                    }
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            ),
            modifier = Modifier.fillMaxWidth()
        )

        // 2-Tab System like Web (Chat/Work)
        TabRow(
            selectedTabIndex = selectedTab,
            containerColor = TchatColors.surface,
            contentColor = TchatColors.primary,
            modifier = Modifier.fillMaxWidth()
        ) {
            Tab(
                selected = selectedTab == 0,
                onClick = { selectedTab = 0 },
                modifier = Modifier.padding(vertical = 12.dp)
            ) {
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.Center
                    ) {
                        Icon(
                            Icons.Filled.Chat,
                            contentDescription = "Personal Chats",
                            modifier = Modifier.size(20.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(
                            "Chat",
                            fontWeight = if (selectedTab == 0) FontWeight.Bold else FontWeight.Medium,
                            style = MaterialTheme.typography.titleMedium
                        )
                    }
                    Text(
                        "All conversations",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
            Tab(
                selected = selectedTab == 1,
                onClick = { selectedTab = 1 },
                modifier = Modifier.padding(vertical = 12.dp)
            ) {
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.Center
                    ) {
                        Icon(
                            Icons.Filled.Business,
                            contentDescription = "Business Chats",
                            modifier = Modifier.size(20.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(
                            "Work",
                            fontWeight = if (selectedTab == 1) FontWeight.Bold else FontWeight.Medium,
                            style = MaterialTheme.typography.titleMedium
                        )
                    }
                    Text(
                        "Business support",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }

        // Tab Content
        when (selectedTab) {
            0 -> {
                // Personal Chat Tab
                Column {
                    // Search Bar
                    TchatInput(
                        value = personalSearchQuery,
                        onValueChange = { personalSearchQuery = it },
                        type = TchatInputType.Search,
                        placeholder = "Search personal chats...",
                        modifier = Modifier.padding(TchatSpacing.md)
                    )

                    // Personal Chat List
                    LazyColumn(
                        modifier = Modifier.weight(1f),
                        contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                        verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        items(getPersonalChats()) { chat ->
                            PersonalChatItem(
                                chat = chat,
                                onClick = { onChatClick(chat.id, chat.name) }
                            )
                        }
                    }
                }
            }
            1 -> {
                // Business/Work Tab
                Column {
                    // Search Bar
                    TchatInput(
                        value = businessSearchQuery,
                        onValueChange = { businessSearchQuery = it },
                        type = TchatInputType.Search,
                        placeholder = "Search business chats...",
                        modifier = Modifier.padding(TchatSpacing.md)
                    )

                    // Business Chat List
                    LazyColumn(
                        modifier = Modifier.weight(1f),
                        contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                        verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                    ) {
                        items(getBusinessChats()) { chat ->
                            BusinessChatItem(
                                chat = chat,
                                onClick = { onChatClick(chat.id, chat.name) }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun PersonalChatItem(
    chat: PersonalChat,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onClick,
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surface
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Avatar
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .clip(CircleShape)
                    .background(TchatColors.primary),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = chat.name.first().toString(),
                    color = TchatColors.onPrimary,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Chat Info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = chat.name,
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.SemiBold,
                        color = TchatColors.onSurface
                    )
                    Text(
                        text = chat.time,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                Spacer(modifier = Modifier.height(4.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = chat.lastMessage,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        modifier = Modifier.weight(1f),
                        maxLines = 1
                    )

                    if (chat.unreadCount > 0) {
                        Badge(
                            modifier = Modifier.size(20.dp),
                            containerColor = TchatColors.primary
                        ) {
                            Text(
                                text = chat.unreadCount.toString(),
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.onPrimary
                            )
                        }
                    }

                    if (chat.isOnline) {
                        Spacer(modifier = Modifier.width(8.dp))
                        Box(
                            modifier = Modifier
                                .size(8.dp)
                                .clip(CircleShape)
                                .background(TchatColors.success)
                        )
                    }
                }

                // Chat Type Badge
                if (chat.type != ChatType.PERSONAL) {
                    Spacer(modifier = Modifier.height(4.dp))
                    Row {
                        Badge(
                            containerColor = when (chat.type) {
                                ChatType.GROUP -> TchatColors.primary.copy(alpha = 0.1f)
                                ChatType.BUSINESS_CUSTOMER -> TchatColors.warning.copy(alpha = 0.1f)
                                ChatType.CHANNEL -> TchatColors.success.copy(alpha = 0.1f)
                                else -> TchatColors.surfaceVariant
                            }
                        ) {
                            Text(
                                text = when (chat.type) {
                                    ChatType.GROUP -> "Group"
                                    ChatType.BUSINESS_CUSTOMER -> "Shop"
                                    ChatType.CHANNEL -> "Channel"
                                    else -> ""
                                },
                                style = MaterialTheme.typography.labelSmall,
                                color = when (chat.type) {
                                    ChatType.GROUP -> TchatColors.primary
                                    ChatType.BUSINESS_CUSTOMER -> TchatColors.warning
                                    ChatType.CHANNEL -> TchatColors.success
                                    else -> TchatColors.onSurfaceVariant
                                }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun BusinessChatItem(
    chat: BusinessChat,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        onClick = onClick,
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surface
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Avatar
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .clip(CircleShape)
                    .background(TchatColors.primary),
                contentAlignment = Alignment.Center
            ) {
                if (chat.isBot) {
                    Icon(
                        Icons.Filled.SmartToy,
                        contentDescription = "AI Bot",
                        modifier = Modifier.size(24.dp),
                        tint = TchatColors.onPrimary
                    )
                } else {
                    Text(
                        text = chat.name.first().toString(),
                        color = TchatColors.onPrimary,
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                }
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Chat Info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = chat.name,
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.SemiBold,
                        color = TchatColors.onSurface
                    )
                    Text(
                        text = chat.time,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                Spacer(modifier = Modifier.height(4.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = chat.lastMessage,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        modifier = Modifier.weight(1f),
                        maxLines = 1
                    )

                    if (chat.unreadCount > 0) {
                        Badge(
                            modifier = Modifier.size(20.dp),
                            containerColor = TchatColors.primary
                        ) {
                            Text(
                                text = chat.unreadCount.toString(),
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.onPrimary
                            )
                        }
                    }
                }

                // Business Info Row
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Lead Score
                    if (chat.leadScore != null) {
                        Badge(
                            containerColor = when {
                                chat.leadScore >= 80 -> TchatColors.success.copy(alpha = 0.1f)
                                chat.leadScore >= 60 -> TchatColors.warning.copy(alpha = 0.1f)
                                else -> TchatColors.error.copy(alpha = 0.1f)
                            }
                        ) {
                            Row(
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    when {
                                        chat.leadScore >= 80 -> Icons.Filled.TrendingUp
                                        chat.leadScore >= 60 -> Icons.Filled.ShowChart
                                        else -> Icons.Filled.TrendingDown
                                    },
                                    contentDescription = "Lead Score",
                                    modifier = Modifier.size(12.dp),
                                    tint = when {
                                        chat.leadScore >= 80 -> TchatColors.success
                                        chat.leadScore >= 60 -> TchatColors.warning
                                        else -> TchatColors.error
                                    }
                                )
                                Spacer(modifier = Modifier.width(4.dp))
                                Text(
                                    text = "${chat.leadScore}",
                                    style = MaterialTheme.typography.labelSmall,
                                    color = when {
                                        chat.leadScore >= 80 -> TchatColors.success
                                        chat.leadScore >= 60 -> TchatColors.warning
                                        else -> TchatColors.error
                                    }
                                )
                            }
                        }
                    }

                    // Revenue
                    if (chat.revenue != null && chat.revenue > 0) {
                        Badge(
                            containerColor = TchatColors.primary.copy(alpha = 0.1f)
                        ) {
                            Text(
                                text = "฿${chat.revenue}K",
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.primary
                            )
                        }
                    }

                    // Bot/Business Type
                    Badge(
                        containerColor = if (chat.isBot)
                            TchatColors.surfaceVariant else
                            TchatColors.warning.copy(alpha = 0.1f)
                    ) {
                        Text(
                            text = if (chat.isBot) "AI Bot" else "Customer",
                            style = MaterialTheme.typography.labelSmall,
                            color = if (chat.isBot)
                                TchatColors.onSurfaceVariant else
                                TchatColors.warning
                        )
                    }
                }
            }
        }
    }
}

// Data Models
private data class PersonalChat(
    val id: String,
    val name: String,
    val lastMessage: String,
    val time: String,
    val unreadCount: Int = 0,
    val isOnline: Boolean = false,
    val type: ChatType = ChatType.PERSONAL
)

private data class BusinessChat(
    val id: String,
    val name: String,
    val lastMessage: String,
    val time: String,
    val unreadCount: Int = 0,
    val isBot: Boolean = false,
    val leadScore: Int? = null,
    val revenue: Int? = null // in thousands
)

private enum class ChatType {
    PERSONAL, GROUP, BUSINESS_CUSTOMER, CHANNEL
}

// Sample Data - Personal Chats (like web)
private fun getPersonalChats(): List<PersonalChat> = listOf(
    PersonalChat("family", "Family Group", "Mom: Dinner at 7pm! 🍽️", "5m", 3, true, ChatType.GROUP),
    PersonalChat("sarah", "Sarah", "See you at the coffee shop! ☕", "15m", 1, true, ChatType.PERSONAL),
    PersonalChat("pad-thai-shop", "Pad Thai Corner", "Your order #247 is ready for pickup! 🍜", "30m", 1, true, ChatType.BUSINESS_CUSTOMER),
    PersonalChat("mike", "Mike Chen", "Thanks for the pad thai recommendation!", "1h", 0, false, ChatType.PERSONAL),
    PersonalChat("grocery-mart", "FreshMart Grocery", "Hi! Your weekly groceries are 20% off today 🛒", "2h", 0, false, ChatType.BUSINESS_CUSTOMER),
    PersonalChat("thai-news", "Thailand News", "Breaking: New digital payment regulations announced", "6h", 0, false, ChatType.CHANNEL)
)

// Sample Data - Business Chats (like web)
private fun getBusinessChats(): List<BusinessChat> = listOf(
    BusinessChat("customer-jane", "Jane Wilson", "Do you deliver to Sukhumvit area? 📍", "5m", 2, false, 82, 12),
    BusinessChat("customer-david", "David Kim", "Great service! Will order again next week 👍", "15m", 0, false, 95, 85),
    BusinessChat("support-bot", "Restaurant Support Bot", "Handling 3 customer inquiries automatically", "45m", 0, true, 95, null),
    BusinessChat("customer-lisa", "Lisa Chen", "Can I order 50 coffees for office meeting? ☕", "10m", 1, false, 75, 28),
    BusinessChat("customer-tom", "Tom Anderson", "Love the new seasonal blend! 🍂", "25m", 0, false, 88, 4),
    BusinessChat("coffee-bot", "Coffee Support Bot", "Handling loyalty program inquiries", "1h", 0, true, 90, null)
)