package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.CircleShape
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
fun EditProfileScreen(
    onBackClick: () -> Unit,
    onSave: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var name by remember { mutableStateOf("Alice Johnson") }
    var username by remember { mutableStateOf("alice_designs") }
    var email by remember { mutableStateOf("alice.johnson@example.com") }
    var phone by remember { mutableStateOf("+1 (555) 123-4567") }
    var bio by remember { mutableStateOf("UI/UX Designer • Coffee enthusiast • Based in SF") }
    var website by remember { mutableStateOf("www.alicedesigns.com") }
    var location by remember { mutableStateOf("San Francisco, CA") }
    var isPublicProfile by remember { mutableStateOf(true) }
    var allowMessages by remember { mutableStateOf(true) }

    var hasChanges by remember { mutableStateOf(false) }

    // Track changes
    LaunchedEffect(name, username, email, phone, bio, website, location, isPublicProfile, allowMessages) {
        hasChanges = true
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Edit Profile") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.Close, "Cancel")
                    }
                },
                actions = {
                    TextButton(
                        onClick = {
                            onSave()
                            hasChanges = false
                        },
                        enabled = hasChanges
                    ) {
                        Text(
                            text = "Save",
                            color = if (hasChanges) TchatColors.primary else TchatColors.onSurfaceVariant,
                            fontWeight = FontWeight.SemiBold
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface
                )
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .background(TchatColors.background),
            contentPadding = PaddingValues(TchatSpacing.md),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
        ) {
            // Profile Picture Section
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                    elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.lg),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Box(
                            modifier = Modifier
                                .size(120.dp)
                                .clip(CircleShape)
                                .background(TchatColors.primary),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = name.firstOrNull()?.toString() ?: "A",
                                color = TchatColors.onPrimary,
                                style = MaterialTheme.typography.headlineLarge,
                                fontWeight = FontWeight.Bold
                            )

                            // Camera icon overlay
                            Box(
                                modifier = Modifier
                                    .align(Alignment.BottomEnd)
                                    .size(40.dp)
                                    .background(
                                        TchatColors.surface,
                                        CircleShape
                                    ),
                                contentAlignment = Alignment.Center
                            ) {
                                IconButton(
                                    onClick = { /* Change photo */ },
                                    modifier = Modifier.size(32.dp)
                                ) {
                                    Icon(
                                        Icons.Default.Edit,
                                        contentDescription = "Change Photo",
                                        modifier = Modifier.size(20.dp),
                                        tint = TchatColors.onSurface
                                    )
                                }
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.sm))

                        TchatButton(
                            onClick = { /* Change photo */ },
                            text = "Change Photo",
                            variant = TchatButtonVariant.Outline
                        )
                    }
                }
            }

            // Personal Information Section
            item {
                SectionHeader("Personal Information")
            }

            item {
                TchatInput(
                    value = name,
                    onValueChange = { name = it },
                    type = TchatInputType.Text,
                    placeholder = "Full Name",
                    label = "Full Name"
                )
            }

            item {
                TchatInput(
                    value = username,
                    onValueChange = { username = it },
                    type = TchatInputType.Text,
                    placeholder = "Username",
                    label = "Username",
                    leadingIcon = Icons.Default.Email
                )
            }

            item {
                TchatInput(
                    value = bio,
                    onValueChange = { bio = it },
                    type = TchatInputType.Multiline,
                    placeholder = "Tell people about yourself...",
                    label = "Bio",
                    maxLines = 3
                )
            }

            // Contact Information Section
            item {
                SectionHeader("Contact Information")
            }

            item {
                TchatInput(
                    value = email,
                    onValueChange = { email = it },
                    type = TchatInputType.Email,
                    placeholder = "Email Address",
                    label = "Email Address",
                    leadingIcon = Icons.Default.Email
                )
            }

            item {
                TchatInput(
                    value = phone,
                    onValueChange = { phone = it },
                    type = TchatInputType.Text,
                    placeholder = "Phone Number",
                    label = "Phone Number",
                    leadingIcon = Icons.Default.Phone
                )
            }

            // Additional Information Section
            item {
                SectionHeader("Additional Information")
            }

            item {
                TchatInput(
                    value = website,
                    onValueChange = { website = it },
                    type = TchatInputType.Text,
                    placeholder = "Website URL",
                    label = "Website",
                    leadingIcon = Icons.Default.Share
                )
            }

            item {
                TchatInput(
                    value = location,
                    onValueChange = { location = it },
                    type = TchatInputType.Text,
                    placeholder = "Location",
                    label = "Location",
                    leadingIcon = Icons.Default.LocationOn
                )
            }

            // Privacy Settings Section
            item {
                SectionHeader("Privacy Settings")
            }

            item {
                PrivacySettingItem(
                    title = "Public Profile",
                    subtitle = "Make your profile visible to everyone",
                    checked = isPublicProfile,
                    onCheckedChange = { isPublicProfile = it }
                )
            }

            item {
                PrivacySettingItem(
                    title = "Allow Messages",
                    subtitle = "Let others send you direct messages",
                    checked = allowMessages,
                    onCheckedChange = { allowMessages = it }
                )
            }

            // Action Buttons Section
            item {
                Column(
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    TchatButton(
                        onClick = {
                            onSave()
                            hasChanges = false
                        },
                        text = "Save Changes",
                        variant = TchatButtonVariant.Primary,
                        enabled = hasChanges,
                        modifier = Modifier.fillMaxWidth()
                    )

                    TchatButton(
                        onClick = onBackClick,
                        text = "Cancel",
                        variant = TchatButtonVariant.Secondary,
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            }

            // Danger Zone Section
            item {
                Spacer(modifier = Modifier.height(TchatSpacing.lg))
                SectionHeader("Danger Zone")
            }

            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = TchatColors.error.copy(alpha = 0.1f)),
                    elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
                ) {
                    Column(
                        modifier = Modifier.padding(TchatSpacing.md)
                    ) {
                        Row(
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                Icons.Default.Warning,
                                contentDescription = null,
                                tint = TchatColors.error,
                                modifier = Modifier.size(24.dp)
                            )

                            Spacer(modifier = Modifier.width(TchatSpacing.sm))

                            Column(
                                modifier = Modifier.weight(1f)
                            ) {
                                Text(
                                    text = "Delete Account",
                                    style = MaterialTheme.typography.titleMedium,
                                    fontWeight = FontWeight.SemiBold,
                                    color = TchatColors.error
                                )
                                Text(
                                    text = "Permanently delete your account and all data",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = TchatColors.onSurfaceVariant
                                )
                            }
                        }

                        Spacer(modifier = Modifier.height(TchatSpacing.sm))

                        TchatButton(
                            onClick = { /* Show delete confirmation */ },
                            text = "Delete Account",
                            variant = TchatButtonVariant.Destructive,
                            modifier = Modifier.fillMaxWidth()
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun SectionHeader(
    title: String,
    modifier: Modifier = Modifier
) {
    Text(
        text = title,
        style = MaterialTheme.typography.titleMedium,
        fontWeight = FontWeight.SemiBold,
        color = TchatColors.onSurface,
        modifier = modifier.padding(
            top = TchatSpacing.sm,
            bottom = TchatSpacing.xs
        )
    )
}

@Composable
private fun PrivacySettingItem(
    title: String,
    subtitle: String,
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
                Text(
                    text = subtitle,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            Switch(
                checked = checked,
                onCheckedChange = onCheckedChange,
                colors = SwitchDefaults.colors(
                    checkedThumbColor = TchatColors.primary,
                    checkedTrackColor = TchatColors.primary.copy(alpha = 0.5f)
                )
            )
        }
    }
}