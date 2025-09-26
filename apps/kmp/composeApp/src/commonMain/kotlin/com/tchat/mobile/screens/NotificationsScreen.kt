package com.tchat.mobile.screens

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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

data class NotificationItem(
    val id: String,
    val type: NotificationType,
    val title: String,
    val description: String,
    val timestamp: String,
    val isRead: Boolean = false,
    val amount: Double? = null
)

enum class NotificationType {
    MESSAGE, PAYMENT, SOCIAL, MERCHANT, SYSTEM
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificationsScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var selectedTab by remember { mutableStateOf(0) }
    var selectedNotifications by remember { mutableStateOf(setOf<String>()) }
    var isSelectMode by remember { mutableStateOf(false) }

    // Mock notifications following web design
    val notifications = listOf(
        NotificationItem(
            "1", NotificationType.MESSAGE,
            "Family Group",
            "Mom: Don't forget dinner tonight! 🍽️",
            "5 min ago", false
        ),
        NotificationItem(
            "2", NotificationType.PAYMENT,
            "Payment Received",
            "You received ฿150.00 from Sarah Chen via PromptPay",
            "15 min ago", false, 150.0
        ),
        NotificationItem(
            "3", NotificationType.MERCHANT,
            "Somtam Vendor",
            "Your order #247 is ready for pickup! 🥗",
            "30 min ago", true
        ),
        NotificationItem(
            "4", NotificationType.SOCIAL,
            "Mike Chen liked your photo",
            "Bangkok Street Food Market experience",
            "1 hour ago", true
        ),
        NotificationItem(
            "5", NotificationType.SYSTEM,
            "Security Alert",
            "New login detected from iPhone in Bangkok",
            "2 hours ago", false
        ),
        NotificationItem(
            "6", NotificationType.PAYMENT,
            "Payment Sent",
            "You sent ฿45.00 to Pad Thai Corner",
            "3 hours ago", true, 45.0
        ),
        NotificationItem(
            "7", NotificationType.MESSAGE,
            "AI Assistant",
            "Welcome to Telegram SEA! Here are some features to get you started",
            "1 day ago", true
        ),
        NotificationItem(
            "8", NotificationType.MERCHANT,
            "FreshMart Grocery",
            "Your weekly groceries are 20% off today! 🛒",
            "1 day ago", false
        )
    )

    val tabs = listOf("All", "Unread", "Messages", "Payments", "Social")
    val filteredNotifications = when (selectedTab) {
        0 -> notifications // All
        1 -> notifications.filter { !it.isRead } // Unread
        2 -> notifications.filter { it.type == NotificationType.MESSAGE } // Messages
        3 -> notifications.filter { it.type == NotificationType.PAYMENT } // Payments
        4 -> notifications.filter { it.type == NotificationType.SOCIAL } // Social
        else -> notifications
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        if (isSelectMode) "${selectedNotifications.size} selected" else "Notifications",
                        fontWeight = FontWeight.Bold
                    )
                },
                navigationIcon = {
                    IconButton(onClick = {
                        if (isSelectMode) {
                            isSelectMode = false
                            selectedNotifications = emptySet()
                        } else {
                            onBackClick()
                        }
                    }) {
                        Icon(
                            if (isSelectMode) Icons.Default.Close else Icons.Default.ArrowBack,
                            contentDescription = if (isSelectMode) "Close" else "Back",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                actions = {
                    if (isSelectMode) {
                        IconButton(onClick = {
                            // Mark selected as read
                            selectedNotifications = emptySet()
                            isSelectMode = false
                        }) {
                            Icon(
                                Icons.Default.Done,
                                contentDescription = "Mark as read",
                                tint = TchatColors.onSurface
                            )
                        }
                        IconButton(onClick = {
                            // Delete selected
                            selectedNotifications = emptySet()
                            isSelectMode = false
                        }) {
                            Icon(
                                Icons.Default.Delete,
                                contentDescription = "Delete",
                                tint = TchatColors.onSurface
                            )
                        }
                    } else {
                        IconButton(onClick = {
                            // Filter options
                        }) {
                            Icon(
                                Icons.Default.FilterList,
                                contentDescription = "Filter",
                                tint = TchatColors.onSurface
                            )
                        }
                        IconButton(onClick = {
                            isSelectMode = true
                        }) {
                            Icon(
                                Icons.Default.MoreVert,
                                contentDescription = "More",
                                tint = TchatColors.onSurface
                            )
                        }
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = modifier
                .fillMaxSize()
                .padding(paddingValues)
                .background(TchatColors.background)
        ) {
            // Tab Row
            TabRow(
                selectedTabIndex = selectedTab,
                containerColor = TchatColors.surface,
                contentColor = TchatColors.primary,
                modifier = Modifier.fillMaxWidth()
            ) {
                tabs.forEachIndexed { index, title ->
                    Tab(
                        selected = selectedTab == index,
                        onClick = { selectedTab = index },
                        text = { Text(title) }
                    )
                }
            }

            // Notifications List
            if (filteredNotifications.isEmpty()) {
                // Empty state
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(TchatSpacing.xl),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Icon(
                        Icons.Default.Notifications,
                        contentDescription = null,
                        modifier = Modifier.size(80.dp),
                        tint = TchatColors.onSurfaceVariant.copy(alpha = 0.5f)
                    )
                    Spacer(modifier = Modifier.height(TchatSpacing.md))
                    Text(
                        "No notifications",
                        style = MaterialTheme.typography.titleMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                    Text(
                        "You're all caught up!",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            } else {
                LazyColumn(
                    modifier = Modifier.weight(1f),
                    contentPadding = PaddingValues(TchatSpacing.sm),
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                ) {
                    items(filteredNotifications) { notification ->
                        NotificationListItem(
                            notification = notification,
                            isSelected = notification.id in selectedNotifications,
                            isSelectMode = isSelectMode,
                            onToggleSelection = { notificationId ->
                                selectedNotifications = if (notificationId in selectedNotifications) {
                                    selectedNotifications - notificationId
                                } else {
                                    selectedNotifications + notificationId
                                }
                            },
                            onLongPress = {
                                if (!isSelectMode) {
                                    isSelectMode = true
                                    selectedNotifications = setOf(notification.id)
                                }
                            },
                            onClick = {
                                if (isSelectMode) {
                                    // Toggle selection
                                    selectedNotifications = if (notification.id in selectedNotifications) {
                                        selectedNotifications - notification.id
                                    } else {
                                        selectedNotifications + notification.id
                                    }
                                } else {
                                    // Handle notification click
                                }
                            }
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun NotificationListItem(
    notification: NotificationItem,
    isSelected: Boolean,
    isSelectMode: Boolean,
    onToggleSelection: (String) -> Unit,
    onLongPress: () -> Unit,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clickable { onClick() },
        colors = CardDefaults.cardColors(
            containerColor = if (notification.isRead) {
                TchatColors.surface
            } else {
                TchatColors.primary.copy(alpha = 0.1f)
            }
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.Top
        ) {
            if (isSelectMode) {
                Checkbox(
                    checked = isSelected,
                    onCheckedChange = { onToggleSelection(notification.id) },
                    colors = CheckboxDefaults.colors(
                        checkedColor = TchatColors.primary,
                        uncheckedColor = TchatColors.onSurfaceVariant
                    )
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
            }

            // Icon
            Box(
                modifier = Modifier
                    .size(40.dp)
                    .clip(CircleShape)
                    .background(
                        when (notification.type) {
                            NotificationType.MESSAGE -> TchatColors.primary
                            NotificationType.PAYMENT -> TchatColors.success
                            NotificationType.SOCIAL -> TchatColors.warning
                            NotificationType.MERCHANT -> TchatColors.primaryLight
                            NotificationType.SYSTEM -> TchatColors.surfaceVariant
                        }
                    ),
                contentAlignment = Alignment.Center
            ) {
                val icon = when (notification.type) {
                    NotificationType.MESSAGE -> Icons.Default.Message
                    NotificationType.PAYMENT -> Icons.Default.Payment
                    NotificationType.SOCIAL -> Icons.Default.Favorite
                    NotificationType.MERCHANT -> Icons.Default.Store
                    NotificationType.SYSTEM -> Icons.Default.Security
                }
                Icon(
                    icon,
                    contentDescription = null,
                    tint = TchatColors.onPrimary,
                    modifier = Modifier.size(20.dp)
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Content
            Column(modifier = Modifier.weight(1f)) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        notification.title,
                        style = MaterialTheme.typography.bodyLarge,
                        fontWeight = if (notification.isRead) FontWeight.Medium else FontWeight.Bold,
                        color = TchatColors.onSurface,
                        modifier = Modifier.weight(1f)
                    )

                    notification.amount?.let { amount ->
                        Text(
                            "฿${amount.toInt()}",
                            style = MaterialTheme.typography.bodyMedium,
                            fontWeight = FontWeight.Bold,
                            color = TchatColors.primary
                        )
                    }
                }

                Text(
                    notification.description,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant,
                    modifier = Modifier.padding(top = 2.dp)
                )

                Text(
                    notification.timestamp,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    modifier = Modifier.padding(top = 4.dp)
                )
            }

            // Unread indicator
            if (!notification.isRead) {
                Box(
                    modifier = Modifier
                        .size(8.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primary)
                )
            }
        }
    }
}