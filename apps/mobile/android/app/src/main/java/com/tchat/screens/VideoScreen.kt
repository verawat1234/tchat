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
 * Video conferencing and streaming interface screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VideoScreen() {
    var selectedTab by remember { mutableStateOf(VideoTab.MEETINGS) }
    var showingCreateMeeting by remember { mutableStateOf(false) }
    var isCameraOn by remember { mutableStateOf(false) }
    var isMicOn by remember { mutableStateOf(true) }

    // Mock meetings
    val upcomingMeetings = listOf(
        Meeting("Team Standup", "09:00 AM", "with John, Sarah, Mike", Icons.Default.VideoCall),
        Meeting("Client Review", "02:00 PM", "with ABC Company", Icons.Default.Group),
        Meeting("Project Planning", "04:30 PM", "with Development Team", Icons.Default.EventNote)
    )

    // Mock recordings
    val recordings = listOf(
        Recording("Q3 Review Meeting", "45 min", "Yesterday", Icons.Default.PlayArrow),
        Recording("Product Demo", "28 min", "2 days ago", Icons.Default.PlayArrow),
        Recording("Training Session", "1h 12min", "1 week ago", Icons.Default.PlayArrow)
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
                    text = "Video",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.textPrimary
                )
            },
            actions = {
                IconButton(onClick = { showingCreateMeeting = true }) {
                    Icon(
                        imageVector = Icons.Default.Add,
                        contentDescription = "Create meeting",
                        tint = Colors.primary
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = Colors.background
            )
        )

        // Tab selector
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .background(Colors.surface)
        ) {
            VideoTab.values().forEach { tab ->
                Column(
                    modifier = Modifier
                        .weight(1f)
                        .clickable { selectedTab = tab }
                        .padding(vertical = Spacing.sm),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Icon(
                        imageVector = tab.icon,
                        contentDescription = tab.title,
                        tint = if (selectedTab == tab) Colors.primary else Colors.textSecondary,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.height(Spacing.xs))
                    Text(
                        text = tab.title,
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Medium,
                        color = if (selectedTab == tab) Colors.primary else Colors.textSecondary
                    )
                }
            }
        }

        // Content based on selected tab
        when (selectedTab) {
            VideoTab.MEETINGS -> MeetingsView(meetings = upcomingMeetings)
            VideoTab.RECORDINGS -> RecordingsView(recordings = recordings)
            VideoTab.LIVE -> LiveView(
                isCameraOn = isCameraOn,
                isMicOn = isMicOn,
                onCameraToggle = { isCameraOn = !isCameraOn },
                onMicToggle = { isMicOn = !isMicOn }
            )
        }
    }
}

// MARK: - Video Tab Enum
enum class VideoTab(
    val title: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
) {
    MEETINGS("Meetings", Icons.Default.EventNote),
    RECORDINGS("Recordings", Icons.Default.PlayArrow),
    LIVE("Live", Icons.Default.VideoCall)
}

// MARK: - Data Classes
data class Meeting(
    val title: String,
    val time: String,
    val participants: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
)

data class Recording(
    val title: String,
    val duration: String,
    val date: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
)

// MARK: - Meetings View
@Composable
private fun MeetingsView(meetings: List<Meeting>) {
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        item {
            // Quick join section
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
                Column(
                    modifier = Modifier.padding(Spacing.md)
                ) {
                    Text(
                        text = "Quick Join",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        color = Colors.textPrimary
                    )

                    Spacer(modifier = Modifier.height(Spacing.sm))

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        OutlinedTextField(
                            value = "",
                            onValueChange = { },
                            placeholder = {
                                Text(
                                    text = "Enter meeting ID",
                                    color = Colors.textSecondary
                                )
                            },
                            modifier = Modifier.weight(1f),
                            shape = RoundedCornerShape(12.dp)
                        )

                        Spacer(modifier = Modifier.width(Spacing.sm))

                        Button(
                            onClick = { /* Join meeting */ },
                            colors = ButtonDefaults.buttonColors(
                                containerColor = Colors.primary
                            ),
                            shape = RoundedCornerShape(12.dp)
                        ) {
                            Text(
                                text = "Join",
                                color = Colors.textOnPrimary,
                                fontWeight = FontWeight.SemiBold
                            )
                        }
                    }
                }
            }
        }

        item {
            Text(
                text = "Upcoming Meetings",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.textPrimary
            )
        }

        items(meetings) { meeting ->
            MeetingCard(meeting = meeting)
        }
    }
}

// MARK: - Meeting Card Component
@Composable
private fun MeetingCard(meeting: Meeting) {
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
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Meeting icon
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .background(
                        color = Colors.primary.copy(alpha = 0.1f),
                        shape = RoundedCornerShape(12.dp)
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = meeting.icon,
                    contentDescription = meeting.title,
                    tint = Colors.primary,
                    modifier = Modifier.size(24.dp)
                )
            }

            Spacer(modifier = Modifier.width(Spacing.md))

            // Meeting info
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = meeting.title,
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Colors.textPrimary
                )
                Text(
                    text = meeting.time,
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Medium,
                    color = Colors.primary
                )
                Text(
                    text = meeting.participants,
                    fontSize = 12.sp,
                    color = Colors.textSecondary
                )
            }

            // Join button
            Button(
                onClick = { /* Join meeting */ },
                colors = ButtonDefaults.buttonColors(
                    containerColor = Colors.primary
                ),
                shape = RoundedCornerShape(8.dp),
                contentPadding = PaddingValues(horizontal = Spacing.md, vertical = Spacing.xs)
            ) {
                Text(
                    text = "Join",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Colors.textOnPrimary
                )
            }
        }
    }
}

// MARK: - Recordings View
@Composable
private fun RecordingsView(recordings: List<Recording>) {
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        contentPadding = PaddingValues(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.sm)
    ) {
        items(recordings) { recording ->
            RecordingCard(recording = recording)
        }
    }
}

// MARK: - Recording Card Component
@Composable
private fun RecordingCard(recording: Recording) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = androidx.compose.ui.graphics.Color.White
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 2.dp
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(Spacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Thumbnail placeholder
            Box(
                modifier = Modifier
                    .size(width = 60.dp, height = 45.dp)
                    .background(
                        color = Colors.surface,
                        shape = RoundedCornerShape(8.dp)
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = recording.icon,
                    contentDescription = "Play",
                    tint = Colors.primary,
                    modifier = Modifier.size(24.dp)
                )
            }

            Spacer(modifier = Modifier.width(Spacing.md))

            // Recording info
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = recording.title,
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Medium,
                    color = Colors.textPrimary,
                    maxLines = 1
                )
                Row {
                    Text(
                        text = recording.duration,
                        fontSize = 12.sp,
                        color = Colors.textSecondary
                    )
                    Text(
                        text = " • ",
                        fontSize = 12.sp,
                        color = Colors.textSecondary
                    )
                    Text(
                        text = recording.date,
                        fontSize = 12.sp,
                        color = Colors.textSecondary
                    )
                }
            }

            // More options
            IconButton(onClick = { /* More options */ }) {
                Icon(
                    imageVector = Icons.Default.MoreVert,
                    contentDescription = "More",
                    tint = Colors.textSecondary
                )
            }
        }
    }
}

// MARK: - Live View
@Composable
private fun LiveView(
    isCameraOn: Boolean,
    isMicOn: Boolean,
    onCameraToggle: () -> Unit,
    onMicToggle: () -> Unit
) {
    Column(
        modifier = Modifier.fillMaxSize()
    ) {
        Spacer(modifier = Modifier.weight(1f))

        // Video preview
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(300.dp)
                .padding(horizontal = Spacing.md)
                .background(
                    color = Colors.surface,
                    shape = RoundedCornerShape(16.dp)
                ),
            contentAlignment = Alignment.Center
        ) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Icon(
                    imageVector = if (isCameraOn) Icons.Default.Videocam else Icons.Default.VideocamOff,
                    contentDescription = if (isCameraOn) "Camera Active" else "Camera Off",
                    tint = if (isCameraOn) Colors.primary else Colors.textSecondary,
                    modifier = Modifier.size(48.dp)
                )
                Spacer(modifier = Modifier.height(Spacing.sm))
                Text(
                    text = if (isCameraOn) "Camera Active" else "Camera Off",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Medium,
                    color = if (isCameraOn) Colors.textPrimary else Colors.textSecondary
                )
            }
        }

        Spacer(modifier = Modifier.weight(1f))

        // Control buttons
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(bottom = Spacing.xl),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            // Mic button
            FloatingActionButton(
                onClick = onMicToggle,
                containerColor = if (isMicOn) Colors.success else Colors.surface,
                modifier = Modifier.size(56.dp)
            ) {
                Icon(
                    imageVector = if (isMicOn) Icons.Default.Mic else Icons.Default.MicOff,
                    contentDescription = if (isMicOn) "Mute" else "Unmute",
                    tint = if (isMicOn) Colors.textOnPrimary else Colors.error,
                    modifier = Modifier.size(24.dp)
                )
            }

            // Camera button
            FloatingActionButton(
                onClick = onCameraToggle,
                containerColor = if (isCameraOn) Colors.success else Colors.surface,
                modifier = Modifier.size(56.dp)
            ) {
                Icon(
                    imageVector = if (isCameraOn) Icons.Default.Videocam else Icons.Default.VideocamOff,
                    contentDescription = if (isCameraOn) "Turn off camera" else "Turn on camera",
                    tint = if (isCameraOn) Colors.textOnPrimary else Colors.error,
                    modifier = Modifier.size(24.dp)
                )
            }

            // End call button
            FloatingActionButton(
                onClick = { /* End call */ },
                containerColor = Colors.error,
                modifier = Modifier.size(56.dp)
            ) {
                Icon(
                    imageVector = Icons.Default.CallEnd,
                    contentDescription = "End call",
                    tint = androidx.compose.ui.graphics.Color.White,
                    modifier = Modifier.size(24.dp)
                )
            }
        }
    }
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun VideoScreenPreview() {
    VideoScreen()
}