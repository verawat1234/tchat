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
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CreateChatScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var selectedContacts by remember { mutableStateOf(setOf<String>()) }
    var searchQuery by remember { mutableStateOf("") }
    var groupName by remember { mutableStateOf("") }
    var isGroupChat by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Start New Chat", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                actions = {
                    if (selectedContacts.isNotEmpty()) {
                        TextButton(
                            onClick = {
                                // TODO: Create chat with selected contacts
                                onBackClick()
                            }
                        ) {
                            Text(
                                "Create",
                                color = TchatColors.primary,
                                fontWeight = FontWeight.Bold
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
            // Search Bar
            TchatInput(
                value = searchQuery,
                onValueChange = { searchQuery = it },
                type = TchatInputType.Search,
                placeholder = "Search contacts...",
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md)
            )

            // Group Chat Toggle
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(TchatSpacing.md),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text(
                            text = "Group Chat",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Medium
                        )
                        Text(
                            text = "Create a group conversation",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                    Switch(
                        checked = isGroupChat,
                        onCheckedChange = { isGroupChat = it },
                        colors = SwitchDefaults.colors(
                            checkedThumbColor = TchatColors.primary
                        )
                    )
                }
            }

            // Group Name Input (if group chat)
            if (isGroupChat) {
                Spacer(modifier = Modifier.height(TchatSpacing.md))
                TchatInput(
                    value = groupName,
                    onValueChange = { groupName = it },
                    type = TchatInputType.Text,
                    placeholder = "Group name",
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = TchatSpacing.md)
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Selected Contacts Count
            if (selectedContacts.isNotEmpty()) {
                Card(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = TchatSpacing.md),
                    colors = CardDefaults.cardColors(containerColor = TchatColors.primary.copy(alpha = 0.1f))
                ) {
                    Text(
                        text = "${selectedContacts.size} contact${if (selectedContacts.size != 1) "s" else ""} selected",
                        modifier = Modifier.padding(TchatSpacing.md),
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.primary,
                        fontWeight = FontWeight.Medium
                    )
                }
                Spacer(modifier = Modifier.height(TchatSpacing.sm))
            }

            // Contacts List
            LazyColumn(
                modifier = Modifier.weight(1f),
                contentPadding = PaddingValues(horizontal = TchatSpacing.md),
                verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
            ) {
                items(getMockContacts().filter { contact ->
                    searchQuery.isEmpty() || contact.name.contains(searchQuery, ignoreCase = true)
                }) { contact ->
                    ContactItem(
                        contact = contact,
                        isSelected = selectedContacts.contains(contact.id),
                        onSelectionChange = { isSelected ->
                            selectedContacts = if (isSelected) {
                                selectedContacts + contact.id
                            } else {
                                selectedContacts - contact.id
                            }
                        }
                    )
                }
            }
        }
    }
}

@Composable
private fun ContactItem(
    contact: MockContact,
    isSelected: Boolean,
    onSelectionChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (isSelected) TchatColors.primary.copy(alpha = 0.1f) else TchatColors.surface
        ),
        onClick = { onSelectionChange(!isSelected) }
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
                    text = contact.name.first().toString(),
                    color = TchatColors.onPrimary,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Contact Info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = contact.name,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
                Text(
                    text = contact.status,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            // Selection Checkbox
            Checkbox(
                checked = isSelected,
                onCheckedChange = onSelectionChange,
                colors = CheckboxDefaults.colors(
                    checkedColor = TchatColors.primary
                )
            )
        }
    }
}

// Mock Data
private data class MockContact(
    val id: String,
    val name: String,
    val status: String,
    val isOnline: Boolean = false
)

private fun getMockContacts(): List<MockContact> = listOf(
    MockContact("1", "Alice Chen", "Online", true),
    MockContact("2", "Bob Wilson", "Last seen 2 hours ago"),
    MockContact("3", "Carol Kim", "Online", true),
    MockContact("4", "David Park", "Last seen yesterday"),
    MockContact("5", "Emma Davis", "Online", true),
    MockContact("6", "Frank Miller", "Last seen 1 week ago"),
    MockContact("7", "Grace Lee", "Online", true),
    MockContact("8", "Henry Johnson", "Last seen 3 days ago"),
    MockContact("9", "Ivy Zhang", "Online", true),
    MockContact("10", "Jack Brown", "Last seen 5 minutes ago")
)