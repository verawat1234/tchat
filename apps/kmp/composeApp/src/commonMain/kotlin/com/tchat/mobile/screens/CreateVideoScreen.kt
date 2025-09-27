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
fun CreateVideoScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var videoTitle by remember { mutableStateOf("") }
    var videoDescription by remember { mutableStateOf("") }
    var category by remember { mutableStateOf("") }
    var tags by remember { mutableStateOf("") }
    var privacy by remember { mutableStateOf("") }
    var monetization by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Create Video", fontWeight = FontWeight.Bold) },
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
                    if (videoTitle.isNotEmpty()) {
                        TextButton(
                            onClick = {
                                // TODO: Upload and publish video
                                onBackClick()
                            }
                        ) {
                            Text(
                                "Upload",
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
            // Video Upload Section
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md)
                    .height(220.dp),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                onClick = {
                    // TODO: Open video picker or camera
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
                            Icons.Default.Videocam,
                            contentDescription = "Add Video",
                            modifier = Modifier.size(64.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.height(12.dp))
                        Text(
                            text = "Upload Video",
                            style = MaterialTheme.typography.titleLarge,
                            color = TchatColors.onSurfaceVariant,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(
                            text = "Tap to select video from gallery",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                        Text(
                            text = "or record new video",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            // Video Recording Options
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
            ) {
                TchatButton(
                    text = "Record Video",
                    variant = TchatButtonVariant.Secondary,
                    onClick = {
                        // TODO: Open camera for recording
                    },
                    modifier = Modifier.weight(1f)
                )
                TchatButton(
                    text = "Select from Gallery",
                    variant = TchatButtonVariant.Secondary,
                    onClick = {
                        // TODO: Open gallery
                    },
                    modifier = Modifier.weight(1f)
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Video Details
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
                        text = "Video Details",
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Video Title
                    TchatInput(
                        value = videoTitle,
                        onValueChange = { videoTitle = it },
                        type = TchatInputType.Text,
                        label = "Video Title",
                        placeholder = "Give your video a catchy title",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Description
                    TchatTextarea(
                        value = videoDescription,
                        onValueChange = { videoDescription = it },
                        label = "Description",
                        placeholder = "Tell viewers about your video...",
                        modifier = Modifier.fillMaxWidth(),
                        minLines = 4
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Category
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("entertainment", "Entertainment"),
                            SelectOption("education", "Education"),
                            SelectOption("gaming", "Gaming"),
                            SelectOption("music", "Music"),
                            SelectOption("sports", "Sports"),
                            SelectOption("news_politics", "News & Politics"),
                            SelectOption("science_technology", "Science & Technology"),
                            SelectOption("travel_events", "Travel & Events"),
                            SelectOption("pets_animals", "Pets & Animals"),
                            SelectOption("autos_vehicles", "Autos & Vehicles"),
                            SelectOption("comedy", "Comedy"),
                            SelectOption("film_animation", "Film & Animation"),
                            SelectOption("howto_style", "Howto & Style"),
                            SelectOption("people_blogs", "People & Blogs")
                        ),
                        selectedValue = category,
                        onSelectionChange = { category = it ?: "" },
                        label = "Category",
                        placeholder = "Select category",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Tags
                    TchatInput(
                        value = tags,
                        onValueChange = { tags = it },
                        type = TchatInputType.Text,
                        label = "Tags",
                        placeholder = "Add tags to help people find your video",
                        modifier = Modifier.fillMaxWidth(),
                        leadingIcon = Icons.Default.Tag
                    )
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Privacy & Monetization
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
                        text = "Privacy & Monetization",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Privacy Settings
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("public", "Public"),
                            SelectOption("unlisted", "Unlisted"),
                            SelectOption("private", "Private"),
                            SelectOption("followers_only", "Followers Only"),
                            SelectOption("friends_only", "Friends Only")
                        ),
                        selectedValue = privacy,
                        onSelectionChange = { privacy = it ?: "" },
                        label = "Who can watch?",
                        placeholder = "Select privacy",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Additional Options
                    Column(
                        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                    ) {
                        var allowComments by remember { mutableStateOf(true) }
                        var allowDownloads by remember { mutableStateOf(false) }
                        var ageRestricted by remember { mutableStateOf(false) }

                        Row(
                            modifier = Modifier.fillMaxWidth(),
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
                                text = "Allow comments",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }

                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = allowDownloads,
                                onCheckedChange = { allowDownloads = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Allow downloads",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }

                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = ageRestricted,
                                onCheckedChange = { ageRestricted = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Age restricted (18+)",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }

                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = monetization,
                                onCheckedChange = { monetization = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Enable monetization",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xl))

            // Action Buttons
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md),
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
            ) {
                TchatButton(
                    text = "Save Draft",
                    variant = TchatButtonVariant.Secondary,
                    onClick = {
                        // TODO: Save as draft
                        onBackClick()
                    },
                    modifier = Modifier.weight(1f)
                )
                TchatButton(
                    text = "Upload Video",
                    variant = TchatButtonVariant.Primary,
                    onClick = {
                        // TODO: Upload and publish video
                        onBackClick()
                    },
                    enabled = videoTitle.isNotEmpty(),
                    modifier = Modifier.weight(1f)
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xl))
        }
    }
}