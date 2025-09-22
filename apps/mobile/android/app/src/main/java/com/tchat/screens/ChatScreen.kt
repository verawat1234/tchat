package com.tchat.screens

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
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Main chat interface screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatScreen() {
    var searchText by remember { mutableStateOf("") }
    var selectedChat by remember { mutableStateOf<String?>(null) }

    // Mock chat data
    val chats = listOf(
        Chat("John Doe", "Hey, how's it going?", "2m", true),
        Chat("Sarah Wilson", "Meeting at 3pm today", "15m", false),
        Chat("Team Alpha", "Project update ready", "1h", true),
        Chat("Mom", "Don't forget dinner tonight", "2h", false),
        Chat("Alex Chen", "Thanks for the help!", "3h", false)
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
                    text = "Chats",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.textPrimary
                )
            },
            actions = {
                IconButton(onClick = { /* New chat action */ }) {
                    Icon(
                        imageVector = Icons.Default.Edit,
                        contentDescription = "New chat",
                        tint = Colors.primary
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = Colors.background
            )
        )

        // Search bar
        OutlinedTextField(
            value = searchText,
            onValueChange = { searchText = it },
            placeholder = {
                Text(
                    text = "Search conversations",
                    color = Colors.textSecondary
                )
            },
            leadingIcon = {
                Icon(
                    imageVector = Icons.Default.Search,
                    contentDescription = "Search",
                    tint = Colors.textSecondary
                )
            },
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = Spacing.md, vertical = Spacing.sm),
            shape = RoundedCornerShape(12.dp),
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = Colors.primary,
                unfocusedBorderColor = Colors.border
            )
        )

        // Chat list
        LazyColumn(
            modifier = Modifier.fillMaxWidth(),
            contentPadding = PaddingValues(horizontal = Spacing.md)
        ) {
            items(chats) { chat ->
                ChatRowItem(
                    chat = chat,
                    isSelected = selectedChat == chat.name,
                    onClick = { selectedChat = chat.name }
                )
            }
        }
    }
}

// MARK: - Data Classes
data class Chat(
    val name: String,
    val lastMessage: String,
    val time: String,
    val hasUnread: Boolean
)

// MARK: - Chat Row Component
@Composable
private fun ChatRowItem(
    chat: Chat,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() }
            .background(
                color = if (isSelected) Colors.surface else androidx.compose.ui.graphics.Color.Transparent,
                shape = RoundedCornerShape(12.dp)
            )
            .padding(Spacing.md),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Avatar
        Box(
            modifier = Modifier
                .size(48.dp)
                .background(
                    color = Colors.primary.copy(alpha = 0.2f),
                    shape = CircleShape
                ),
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = chat.name.first().toString(),
                fontSize = 20.sp,
                fontWeight = FontWeight.SemiBold,
                color = Colors.primary
            )
        }

        Spacer(modifier = Modifier.width(Spacing.md))

        // Chat info
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
                    fontSize = 16.sp,
                    fontWeight = if (chat.hasUnread) FontWeight.SemiBold else FontWeight.Medium,
                    color = Colors.textPrimary
                )

                Text(
                    text = chat.time,
                    fontSize = 14.sp,
                    color = Colors.textSecondary
                )
            }

            Spacer(modifier = Modifier.height(Spacing.xs))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = chat.lastMessage,
                    fontSize = 14.sp,
                    color = Colors.textSecondary,
                    maxLines = 1,
                    modifier = Modifier.weight(1f)
                )

                if (chat.hasUnread) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .background(
                                color = Colors.primary,
                                shape = CircleShape
                            )
                    )
                }
            }
        }
    }
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun ChatScreenPreview() {
    ChatScreen()
}