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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Settings and additional features screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MoreScreen() {
    var showingProfile by remember { mutableStateOf(false) }
    var darkModeEnabled by remember { mutableStateOf(false) }
    var notificationsEnabled by remember { mutableStateOf(true) }
    var soundEnabled by remember { mutableStateOf(true) }

    // Menu sections
    val menuSections = listOf(
        MoreSection(
            title = "Account",
            items = listOf(
                MoreItem("Profile", Icons.Default.Person, MoreItemType.PROFILE),
                MoreItem("Privacy", Icons.Default.Security, MoreItemType.PRIVACY),
                MoreItem("Security", Icons.Default.VerifiedUser, MoreItemType.SECURITY),
                MoreItem("Billing", Icons.Default.CreditCard, MoreItemType.BILLING)
            )
        ),
        MoreSection(
            title = "Preferences",
            items = listOf(
                MoreItem("Notifications", Icons.Default.Notifications, MoreItemType.NOTIFICATIONS),
                MoreItem("Appearance", Icons.Default.Palette, MoreItemType.APPEARANCE),
                MoreItem("Language", Icons.Default.Language, MoreItemType.LANGUAGE),
                MoreItem("Storage", Icons.Default.Storage, MoreItemType.STORAGE)
            )
        ),
        MoreSection(
            title = "Support",
            items = listOf(
                MoreItem("Help Center", Icons.Default.Help, MoreItemType.HELP),
                MoreItem("Contact Us", Icons.Default.Email, MoreItemType.CONTACT),
                MoreItem("Report Issue", Icons.Default.Report, MoreItemType.REPORT),
                MoreItem("About", Icons.Default.Info, MoreItemType.ABOUT)
            )
        )
    )

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .background(Colors.background),
        contentPadding = PaddingValues(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        item {
            // Top app bar space
            Text(
                text = "More",
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.textPrimary,
                modifier = Modifier.padding(vertical = Spacing.sm)
            )
        }

        item {
            // Profile header
            ProfileHeaderCard(onProfileClick = { showingProfile = true })
        }

        items(menuSections) { section ->
            MenuSectionCard(
                section = section,
                darkModeEnabled = darkModeEnabled,
                notificationsEnabled = notificationsEnabled,
                onDarkModeToggle = { darkModeEnabled = !darkModeEnabled },
                onNotificationsToggle = { notificationsEnabled = !notificationsEnabled }
            )
        }

        item {
            // Sign out button
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable { /* Sign out */ },
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(
                    containerColor = Colors.surface
                ),
                elevation = CardDefaults.cardElevation(
                    defaultElevation = 4.dp
                )
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(Spacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        imageVector = Icons.Default.ExitToApp,
                        contentDescription = "Sign Out",
                        tint = Colors.error,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(Spacing.md))
                    Text(
                        text = "Sign Out",
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Medium,
                        color = Colors.error
                    )
                }
            }
        }

        item {
            // Version info
            Text(
                text = "Version 1.0.0 (Build 1)",
                fontSize = 12.sp,
                color = Colors.textSecondary,
                modifier = Modifier
                    .fillMaxWidth()
                    .wrapContentWidth(Alignment.CenterHorizontally)
                    .padding(bottom = Spacing.xl)
            )
        }
    }
}

// MARK: - Data Classes
data class MoreSection(
    val title: String,
    val items: List<MoreItem>
)

data class MoreItem(
    val title: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector,
    val type: MoreItemType
)

enum class MoreItemType {
    PROFILE, PRIVACY, SECURITY, BILLING,
    NOTIFICATIONS, APPEARANCE, LANGUAGE, STORAGE,
    HELP, CONTACT, REPORT, ABOUT
}

// MARK: - Profile Header Card
@Composable
private fun ProfileHeaderCard(onProfileClick: () -> Unit) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onProfileClick() },
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = androidx.compose.ui.graphics.Color.White
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 4.dp
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Profile avatar
            Box(
                modifier = Modifier
                    .size(60.dp)
                    .background(
                        color = Colors.primary.copy(alpha = 0.2f),
                        shape = CircleShape
                    ),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "JD",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.primary
                )
            }

            Spacer(modifier = Modifier.width(Spacing.md))

            // Profile info
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = "John Doe",
                    fontSize = 18.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Colors.textPrimary
                )
                Text(
                    text = "john.doe@example.com",
                    fontSize = 14.sp,
                    color = Colors.textSecondary
                )
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .background(
                                color = Colors.success,
                                shape = CircleShape
                            )
                    )
                    Spacer(modifier = Modifier.width(Spacing.xs))
                    Text(
                        text = "Online",
                        fontSize = 12.sp,
                        color = Colors.success
                    )
                }
            }

            Icon(
                imageVector = Icons.Default.ChevronRight,
                contentDescription = "View profile",
                tint = Colors.textSecondary
            )
        }
    }
}

// MARK: - Menu Section Card
@Composable
private fun MenuSectionCard(
    section: MoreSection,
    darkModeEnabled: Boolean,
    notificationsEnabled: Boolean,
    onDarkModeToggle: () -> Unit,
    onNotificationsToggle: () -> Unit
) {
    Column {
        Text(
            text = section.title,
            fontSize = 16.sp,
            fontWeight = FontWeight.SemiBold,
            color = Colors.textPrimary,
            modifier = Modifier.padding(bottom = Spacing.sm)
        )

        Card(
            modifier = Modifier.fillMaxWidth(),
            shape = RoundedCornerShape(12.dp),
            colors = CardDefaults.cardColors(
                containerColor = androidx.compose.ui.graphics.Color.White
            ),
            elevation = CardDefaults.cardElevation(
                defaultElevation = 4.dp
            )
        ) {
            Column {
                section.items.forEachIndexed { index, item ->
                    MenuItemRow(
                        item = item,
                        isLast = index == section.items.size - 1,
                        darkModeEnabled = darkModeEnabled,
                        notificationsEnabled = notificationsEnabled,
                        onDarkModeToggle = onDarkModeToggle,
                        onNotificationsToggle = onNotificationsToggle
                    )
                }
            }
        }
    }
}

// MARK: - Menu Item Row
@Composable
private fun MenuItemRow(
    item: MoreItem,
    isLast: Boolean,
    darkModeEnabled: Boolean,
    notificationsEnabled: Boolean,
    onDarkModeToggle: () -> Unit,
    onNotificationsToggle: () -> Unit
) {
    Column {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clickable {
                    when (item.type) {
                        MoreItemType.APPEARANCE -> onDarkModeToggle()
                        MoreItemType.NOTIFICATIONS -> onNotificationsToggle()
                        else -> {
                            // Handle other item types
                            println("Clicked: ${item.title}")
                        }
                    }
                }
                .padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Icon
            Icon(
                imageVector = item.icon,
                contentDescription = item.title,
                tint = getIconColor(item.type),
                modifier = Modifier.size(24.dp)
            )

            Spacer(modifier = Modifier.width(Spacing.md))

            // Title
            Text(
                text = item.title,
                fontSize = 16.sp,
                color = Colors.textPrimary,
                modifier = Modifier.weight(1f)
            )

            // Trailing element (toggle or chevron)
            when (item.type) {
                MoreItemType.NOTIFICATIONS -> {
                    Switch(
                        checked = notificationsEnabled,
                        onCheckedChange = { onNotificationsToggle() },
                        colors = SwitchDefaults.colors(
                            checkedThumbColor = Colors.textOnPrimary,
                            checkedTrackColor = Colors.primary
                        )
                    )
                }
                MoreItemType.APPEARANCE -> {
                    Switch(
                        checked = darkModeEnabled,
                        onCheckedChange = { onDarkModeToggle() },
                        colors = SwitchDefaults.colors(
                            checkedThumbColor = Colors.textOnPrimary,
                            checkedTrackColor = Colors.primary
                        )
                    )
                }
                else -> {
                    Icon(
                        imageVector = Icons.Default.ChevronRight,
                        contentDescription = "Navigate",
                        tint = Colors.textSecondary,
                        modifier = Modifier.size(16.dp)
                    )
                }
            }
        }

        if (!isLast) {
            Divider(
                modifier = Modifier.padding(start = 56.dp),
                color = Colors.border,
                thickness = 0.5.dp
            )
        }
    }
}

@Composable
private fun getIconColor(type: MoreItemType): androidx.compose.ui.graphics.Color {
    return when (type) {
        MoreItemType.REPORT -> Colors.warning
        MoreItemType.SECURITY, MoreItemType.PRIVACY -> Colors.success
        MoreItemType.BILLING -> Colors.primary
        else -> Colors.textSecondary
    }
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun MoreScreenPreview() {
    MoreScreen()
}