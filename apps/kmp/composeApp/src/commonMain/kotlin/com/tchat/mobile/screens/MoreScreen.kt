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
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MoreScreen(
    onBackClick: () -> Unit = {},
    onEditProfileClick: () -> Unit = {},
    onSettingsClick: () -> Unit = {},
    onUserProfileClick: (userId: String) -> Unit = {},
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier.background(TchatColors.background)
    ) {
        // Top App Bar
        TopAppBar(
            title = { Text("More", fontWeight = FontWeight.Bold) },
            navigationIcon = {
                IconButton(onClick = onBackClick) {
                    Icon(
                        Icons.Default.ArrowBack,
                        contentDescription = "Back",
                        tint = TchatColors.onSurface
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = TchatColors.surface,
                titleContentColor = TchatColors.onSurface
            )
        )

        // Content
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .background(TchatColors.background),
            contentPadding = PaddingValues(TchatSpacing.md),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
        ) {
        // Profile Section
        item {
            ProfileSection()
        }

        // Quick Actions
        item {
            QuickActionsSection()
        }

        // Account Settings
        item {
            MenuSection(
                title = "Account",
                items = getAccountItems()
            )
        }

        // App Settings
        item {
            MenuSection(
                title = "Settings",
                items = getSettingsItems()
            )
        }

        // Support & Info
        item {
            MenuSection(
                title = "Support",
                items = getSupportItems()
            )
        }

        // Sign Out
        item {
            Spacer(modifier = Modifier.height(TchatSpacing.lg))
            TchatButton(
                text = "Sign Out",
                variant = TchatButtonVariant.Destructive,
                onClick = { /* Sign out */ },
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
    }
}

@Composable
private fun ProfileSection() {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.lg),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Profile Avatar
            Box(
                modifier = Modifier
                    .size(72.dp)
                    .clip(CircleShape)
                    .background(TchatColors.primary),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "JD",
                    style = MaterialTheme.typography.headlineMedium,
                    color = TchatColors.onPrimary,
                    fontWeight = FontWeight.Bold
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Profile Info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "John Doe",
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )
                Text(
                    text = "john.doe@tchat.com",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
                Text(
                    text = "Premium Member",
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.primary,
                    fontWeight = FontWeight.Medium
                )
            }

            // Edit Profile Button
            IconButton(
                onClick = { /* Edit profile */ }
            ) {
                Icon(
                    Icons.Filled.Edit,
                    contentDescription = "Edit Profile",
                    tint = TchatColors.primary
                )
            }
        }
    }
}

@Composable
private fun QuickActionsSection() {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.lg)
        ) {
            Text(
                text = "Quick Actions",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold,
                color = TchatColors.onSurface,
                modifier = Modifier.padding(bottom = TchatSpacing.md)
            )

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                QuickActionItem(
                    icon = Icons.Default.Settings,
                    title = "QR Code",
                    onClick = { /* QR code */ }
                )
                QuickActionItem(
                    icon = Icons.Default.AccountBox,
                    title = "Payments",
                    onClick = { /* Payments */ }
                )
                QuickActionItem(
                    icon = Icons.Default.Settings,
                    title = "Storage",
                    onClick = { /* Storage */ }
                )
                QuickActionItem(
                    icon = Icons.Default.Settings,
                    title = "Downloads",
                    onClick = { /* Downloads */ }
                )
            }
        }
    }
}

@Composable
private fun QuickActionItem(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    title: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier
            .clip(RoundedCornerShape(8.dp))
            .padding(TchatSpacing.sm)
    ) {
        IconButton(
            onClick = onClick,
            modifier = Modifier
                .size(48.dp)
                .background(
                    TchatColors.primary.copy(alpha = 0.1f),
                    CircleShape
                )
        ) {
            Icon(
                icon,
                contentDescription = title,
                tint = TchatColors.primary,
                modifier = Modifier.size(24.dp)
            )
        }
        Spacer(modifier = Modifier.height(4.dp))
        Text(
            text = title,
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun MenuSection(
    title: String,
    items: List<MenuItem>
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.lg)
        ) {
            Text(
                text = title,
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold,
                color = TchatColors.onSurface,
                modifier = Modifier.padding(bottom = TchatSpacing.md)
            )

            items.forEach { item ->
                MenuItemRow(item = item)
                if (item != items.last()) {
                    HorizontalDivider(
                        modifier = Modifier.padding(vertical = TchatSpacing.xs),
                        color = TchatColors.outline.copy(alpha = 0.3f)
                    )
                }
            }
        }
    }
}

@Composable
private fun MenuItemRow(
    item: MenuItem,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = TchatSpacing.sm),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            item.icon,
            contentDescription = item.title,
            tint = TchatColors.onSurfaceVariant,
            modifier = Modifier.size(24.dp)
        )

        Spacer(modifier = Modifier.width(TchatSpacing.md))

        Column(
            modifier = Modifier.weight(1f)
        ) {
            Text(
                text = item.title,
                style = MaterialTheme.typography.bodyLarge,
                color = TchatColors.onSurface
            )
            if (item.subtitle != null) {
                Text(
                    text = item.subtitle,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }

        if (item.badge != null) {
            Badge(
                containerColor = TchatColors.primary
            ) {
                Text(
                    text = item.badge,
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onPrimary
                )
            }
            Spacer(modifier = Modifier.width(TchatSpacing.sm))
        }

        Icon(
            Icons.Filled.KeyboardArrowRight,
            contentDescription = "Navigate",
            tint = TchatColors.onSurfaceVariant,
            modifier = Modifier.size(20.dp)
        )
    }
}

// Sample data
private data class MenuItem(
    val icon: androidx.compose.ui.graphics.vector.ImageVector,
    val title: String,
    val subtitle: String? = null,
    val badge: String? = null,
    val onClick: () -> Unit = {}
)

private fun getAccountItems(): List<MenuItem> = listOf(
    MenuItem(Icons.Filled.Person, "Profile Settings", "Manage your profile information"),
    MenuItem(Icons.Filled.Lock, "Privacy & Security", "Control your privacy settings"),
    MenuItem(Icons.Filled.Notifications, "Notifications", "Manage notification preferences"),
    MenuItem(Icons.Default.AccountBox, "Payment Methods", "Manage cards and billing"),
    MenuItem(Icons.Default.Settings, "Language", "English (US)")
)

private fun getSettingsItems(): List<MenuItem> = listOf(
    MenuItem(Icons.Default.Settings, "Dark Mode", "Toggle dark/light theme"),
    MenuItem(Icons.Default.Settings, "Storage", "12.3 GB used"),
    MenuItem(Icons.Default.Settings, "Data Usage", "Monitor your data consumption"),
    MenuItem(Icons.Default.Refresh, "Auto Sync", "Sync across devices"),
    MenuItem(Icons.Default.Refresh, "App Updates", "Check for updates", "1")
)

private fun getSupportItems(): List<MenuItem> = listOf(
    MenuItem(Icons.Default.Info, "Help Center", "Get help and support"),
    MenuItem(Icons.Default.Send, "Send Feedback", "Report issues or suggestions"),
    MenuItem(Icons.Filled.Info, "About", "Version 1.0.0"),
    MenuItem(Icons.Default.Info, "Terms of Service", "Read our terms"),
    MenuItem(Icons.Default.Info, "Privacy Policy", "Read our privacy policy")
)