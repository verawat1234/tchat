package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    onBackClick: () -> Unit,
    onEditProfileClick: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var notificationsEnabled by remember { mutableStateOf(true) }
    var darkModeEnabled by remember { mutableStateOf(false) }
    var soundEnabled by remember { mutableStateOf(true) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Settings") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, "Back")
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
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            // Account Section
            item {
                SettingsSectionHeader("Account")
            }

            items(getAccountSettings(onEditProfileClick)) { setting ->
                SettingsItem(setting = setting)
            }

            // Preferences Section
            item {
                SettingsSectionHeader("Preferences")
            }

            item {
                SwitchSettingsItem(
                    icon = Icons.Default.Notifications,
                    title = "Push Notifications",
                    subtitle = "Receive notifications for messages and updates",
                    checked = notificationsEnabled,
                    onCheckedChange = { notificationsEnabled = it }
                )
            }

            item {
                SwitchSettingsItem(
                    icon = Icons.Default.Settings,
                    title = "Dark Mode",
                    subtitle = "Switch between light and dark theme",
                    checked = darkModeEnabled,
                    onCheckedChange = { darkModeEnabled = it }
                )
            }

            item {
                SwitchSettingsItem(
                    icon = Icons.Default.Settings,
                    title = "Sound Effects",
                    subtitle = "Play sound for notifications and interactions",
                    checked = soundEnabled,
                    onCheckedChange = { soundEnabled = it }
                )
            }

            // Privacy & Security Section
            item {
                SettingsSectionHeader("Privacy & Security")
            }

            items(getPrivacySettings()) { setting ->
                SettingsItem(setting = setting)
            }

            // Support Section
            item {
                SettingsSectionHeader("Support")
            }

            items(getSupportSettings()) { setting ->
                SettingsItem(setting = setting)
            }

            // About Section
            item {
                SettingsSectionHeader("About")
            }

            items(getAboutSettings()) { setting ->
                SettingsItem(setting = setting)
            }

            // Sign Out
            item {
                Spacer(modifier = Modifier.height(TchatSpacing.lg))

                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                    elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(TchatSpacing.md),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            Icons.Default.ExitToApp,
                            contentDescription = null,
                            tint = TchatColors.error,
                            modifier = Modifier.size(24.dp)
                        )

                        Spacer(modifier = Modifier.width(TchatSpacing.md))

                        Column(
                            modifier = Modifier.weight(1f)
                        ) {
                            Text(
                                text = "Sign Out",
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Medium,
                                color = TchatColors.error
                            )
                            Text(
                                text = "Sign out of your account",
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun SettingsSectionHeader(
    title: String,
    modifier: Modifier = Modifier
) {
    Text(
        text = title,
        style = MaterialTheme.typography.titleSmall,
        fontWeight = FontWeight.SemiBold,
        color = TchatColors.primary,
        modifier = modifier.padding(
            top = TchatSpacing.md,
            bottom = TchatSpacing.xs
        )
    )
}

@Composable
private fun SettingsItem(
    setting: SettingsItemData,
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
            Icon(
                setting.icon,
                contentDescription = null,
                tint = TchatColors.onSurfaceVariant,
                modifier = Modifier.size(24.dp)
            )

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = setting.title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
                if (setting.subtitle != null) {
                    Text(
                        text = setting.subtitle,
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }

            if (setting.showChevron) {
                Icon(
                    Icons.Default.KeyboardArrowRight,
                    contentDescription = null,
                    tint = TchatColors.onSurfaceVariant,
                    modifier = Modifier.size(20.dp)
                )
            }
        }
    }
}

@Composable
private fun SwitchSettingsItem(
    icon: ImageVector,
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
            Icon(
                icon,
                contentDescription = null,
                tint = TchatColors.onSurfaceVariant,
                modifier = Modifier.size(24.dp)
            )

            Spacer(modifier = Modifier.width(TchatSpacing.md))

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

// Data model
private data class SettingsItemData(
    val icon: ImageVector,
    val title: String,
    val subtitle: String? = null,
    val showChevron: Boolean = true,
    val onClick: () -> Unit = {}
)

private fun getAccountSettings(onEditProfileClick: () -> Unit): List<SettingsItemData> = listOf(
    SettingsItemData(
        icon = Icons.Default.Person,
        title = "Edit Profile",
        subtitle = "Update your profile information",
        onClick = onEditProfileClick
    ),
    SettingsItemData(
        icon = Icons.Default.Lock,
        title = "Change Password",
        subtitle = "Update your account password"
    ),
    SettingsItemData(
        icon = Icons.Default.Email,
        title = "Email Settings",
        subtitle = "Manage your email preferences"
    ),
    SettingsItemData(
        icon = Icons.Default.Phone,
        title = "Phone Number",
        subtitle = "Add or update your phone number"
    )
)

private fun getPrivacySettings(): List<SettingsItemData> = listOf(
    SettingsItemData(
        icon = Icons.Default.Lock,
        title = "Privacy Settings",
        subtitle = "Control who can see your information"
    ),
    SettingsItemData(
        icon = Icons.Default.Close,
        title = "Blocked Users",
        subtitle = "Manage blocked accounts"
    ),
    SettingsItemData(
        icon = Icons.Default.Lock,
        title = "Two-Factor Authentication",
        subtitle = "Add an extra layer of security"
    ),
    SettingsItemData(
        icon = Icons.Default.Settings,
        title = "Data Usage",
        subtitle = "Monitor your data consumption"
    )
)

private fun getSupportSettings(): List<SettingsItemData> = listOf(
    SettingsItemData(
        icon = Icons.Default.Info,
        title = "Help Center",
        subtitle = "Get help and support"
    ),
    SettingsItemData(
        icon = Icons.Default.Warning,
        title = "Report a Problem",
        subtitle = "Let us know about issues"
    ),
    SettingsItemData(
        icon = Icons.Default.Send,
        title = "Send Feedback",
        subtitle = "Share your thoughts with us"
    ),
    SettingsItemData(
        icon = Icons.Default.Email,
        title = "Contact Support",
        subtitle = "Reach out to our support team"
    )
)

private fun getAboutSettings(): List<SettingsItemData> = listOf(
    SettingsItemData(
        icon = Icons.Default.Info,
        title = "About Tchat",
        subtitle = "Version 1.0.0"
    ),
    SettingsItemData(
        icon = Icons.Default.List,
        title = "Terms of Service",
        subtitle = "Read our terms and conditions"
    ),
    SettingsItemData(
        icon = Icons.Default.Lock,
        title = "Privacy Policy",
        subtitle = "Learn how we handle your data"
    ),
    SettingsItemData(
        icon = Icons.Default.Info,
        title = "Open Source Licenses",
        subtitle = "View third-party licenses"
    )
)