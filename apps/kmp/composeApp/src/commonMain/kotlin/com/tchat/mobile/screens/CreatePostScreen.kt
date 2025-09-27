package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatButtonVariant
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.components.TchatSingleSelect
import com.tchat.mobile.components.SelectOption
import com.tchat.mobile.components.TchatTextarea
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CreatePostScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var postContent by remember { mutableStateOf("") }
    var postType by remember { mutableStateOf("") }
    var visibility by remember { mutableStateOf("") }
    var tags by remember { mutableStateOf("") }
    var location by remember { mutableStateOf("") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Create Post", fontWeight = FontWeight.Bold) },
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
                    if (postContent.isNotEmpty()) {
                        TextButton(
                            onClick = {
                                // TODO: Create post
                                onBackClick()
                            }
                        ) {
                            Text(
                                "Post",
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
                .verticalScroll(rememberScrollState())
        ) {
            // Media Upload Section
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md)
                    .height(200.dp),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                onClick = {
                    // TODO: Open media picker (photo/video)
                }
            ) {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.AddPhotoAlternate,
                            contentDescription = "Add Media",
                            modifier = Modifier.size(48.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = "Add Photos or Videos",
                            style = MaterialTheme.typography.titleMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                        Text(
                            text = "Tap to add media to your post",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            // Post Content
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Text(
                        text = "What's on your mind?",
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Post Content
                    TchatTextarea(
                        value = postContent,
                        onValueChange = { postContent = it },
                        label = "Share your thoughts...",
                        placeholder = "What's happening?",
                        modifier = Modifier.fillMaxWidth(),
                        minLines = 4
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Post Type
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("text", "Text Post"),
                            SelectOption("photo", "Photo Post"),
                            SelectOption("video", "Video Post"),
                            SelectOption("link", "Link Share"),
                            SelectOption("poll", "Poll"),
                            SelectOption("event", "Event"),
                            SelectOption("question", "Question"),
                            SelectOption("review", "Review")
                        ),
                        selectedValue = postType,
                        onSelectionChange = { postType = it ?: "" },
                        label = "Post Type",
                        placeholder = "Select post type",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Tags
                    TchatInput(
                        value = tags,
                        onValueChange = { tags = it },
                        type = TchatInputType.Text,
                        label = "Tags",
                        placeholder = "#hashtags #separated #by #spaces",
                        modifier = Modifier.fillMaxWidth(),
                        leadingIcon = Icons.Default.Tag
                    )
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Privacy & Settings
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
            ) {
                Column(
                    modifier = Modifier.padding(TchatSpacing.md)
                ) {
                    Text(
                        text = "Privacy & Settings",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Visibility
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("public", "Public"),
                            SelectOption("friends", "Friends"),
                            SelectOption("friends_of_friends", "Friends of Friends"),
                            SelectOption("only_me", "Only Me"),
                            SelectOption("custom", "Custom")
                        ),
                        selectedValue = visibility,
                        onSelectionChange = { visibility = it ?: "" },
                        label = "Who can see this?",
                        placeholder = "Select visibility",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Location
                    TchatInput(
                        value = location,
                        onValueChange = { location = it },
                        type = TchatInputType.Text,
                        label = "Location (Optional)",
                        placeholder = "Where are you?",
                        modifier = Modifier.fillMaxWidth(),
                        trailingIcon = Icons.Default.LocationOn
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Post Options
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                    ) {
                        var allowComments by remember { mutableStateOf(true) }
                        var allowSharing by remember { mutableStateOf(true) }

                        Row(
                            modifier = Modifier.weight(1f),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = allowComments,
                                onCheckedChange = { allowComments = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Allow Comments",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }

                        Row(
                            modifier = Modifier.weight(1f),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = allowSharing,
                                onCheckedChange = { allowSharing = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Allow Sharing",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xl))

            // Submit Button
            TchatButton(
                text = "Create Post",
                variant = TchatButtonVariant.Primary,
                onClick = {
                    // TODO: Create post
                    onBackClick()
                },
                enabled = postContent.isNotEmpty(),
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md)
            )

            Spacer(modifier = Modifier.height(TchatSpacing.xl))
        }
    }
}