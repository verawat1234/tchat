package com.tchat.mobile.components.chat

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.components.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.MessageType

/**
 * Image message component
 */
@Composable
fun ImageMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.size(200.dp, 150.dp)
    ) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(
                    Icons.Default.Image,
                    contentDescription = "Image",
                    modifier = Modifier.size(48.dp),
                    tint = TchatColors.primary
                )
                Text(
                    text = "Photo",
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.primary
                )
            }
        }
    }
    Spacer(modifier = Modifier.height(4.dp))
    Text(
        text = content,
        style = MaterialTheme.typography.bodySmall,
        color = TchatColors.onSurfaceVariant.copy(alpha = 0.8f)
    )
}

/**
 * Video message component
 */
@Composable
fun VideoMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatVideo(
        source = VideoSource.Local("sample_video"),
        modifier = modifier.size(250.dp, 140.dp)
    )
    Spacer(modifier = Modifier.height(4.dp))
    Text(
        text = content,
        style = MaterialTheme.typography.bodySmall,
        color = TchatColors.onSurfaceVariant.copy(alpha = 0.8f)
    )
}

/**
 * Audio message component
 */
@Composable
fun AudioMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    var isPlaying by remember { mutableStateOf(false) }
    var currentTime by remember { mutableStateOf(0) }
    val totalTime = 45 // seconds

    TchatCard(
        variant = TchatCardVariant.Filled,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                IconButton(
                    onClick = {
                        isPlaying = !isPlaying
                        // TODO: Implement actual audio playback
                    },
                    modifier = Modifier.size(40.dp)
                ) {
                    Icon(
                        if (isPlaying) Icons.Default.Pause else Icons.Default.PlayArrow,
                        contentDescription = if (isPlaying) "Pause" else "Play",
                        modifier = Modifier.size(32.dp),
                        tint = TchatColors.primary
                    )
                }
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = "Voice Message",
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium
                    )
                    Text(
                        text = "${currentTime / 60}:${(currentTime % 60).toString().padStart(2, '0')} / ${totalTime / 60}:${(totalTime % 60).toString().padStart(2, '0')}",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
                Icon(
                    Icons.Default.GraphicEq,
                    contentDescription = "Waveform",
                    tint = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            // Progress bar
            LinearProgressIndicator(
                progress = currentTime.toFloat() / totalTime,
                modifier = Modifier.fillMaxWidth(),
                color = TchatColors.primary,
                trackColor = TchatColors.surfaceVariant
            )
        }
    }
    Spacer(modifier = Modifier.height(4.dp))
    Text(
        text = content,
        style = MaterialTheme.typography.bodySmall,
        color = TchatColors.onSurfaceVariant.copy(alpha = 0.8f)
    )
}

/**
 * File message component
 */
@Composable
fun FileMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    var isDownloaded by remember { mutableStateOf(false) }

    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.Description,
                    contentDescription = "File",
                    modifier = Modifier.size(32.dp),
                    tint = TchatColors.primary
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = "Meeting_Notes.pdf",
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium
                    )
                    Text(
                        text = "2.4 MB",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
                TchatButton(
                    text = if (isDownloaded) "Open" else "Download",
                    variant = TchatButtonVariant.Outline,
                    size = TchatButtonSize.Small,
                    onClick = {
                        if (isDownloaded) {
                            // TODO: Open file
                        } else {
                            isDownloaded = true
                            // TODO: Download file
                        }
                    }
                )
            }
        }
    }
}

/**
 * Location message component
 */
@Composable
fun LocationMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Filled,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            // Map preview placeholder
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .background(TchatColors.surfaceVariant, RoundedCornerShape(8.dp)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.LocationOn,
                    contentDescription = "Location Pin",
                    modifier = Modifier.size(48.dp),
                    tint = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(
                    Icons.Default.LocationOn,
                    contentDescription = "Location",
                    tint = TchatColors.primary,
                    modifier = Modifier.size(20.dp)
                )
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Text(
                    text = "Downtown Office Building",
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium
                )
            }
            Text(
                text = "123 Business St, City Center",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant
            )
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Row {
                TchatButton(
                    text = "Directions",
                    variant = TchatButtonVariant.Primary,
                    size = TchatButtonSize.Small,
                    onClick = { /* TODO: Open directions */ }
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                TchatButton(
                    text = "Share",
                    variant = TchatButtonVariant.Outline,
                    size = TchatButtonSize.Small,
                    onClick = { /* TODO: Share location */ }
                )
            }
        }
    }
}

/**
 * Payment message component
 */
@Composable
fun PaymentMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Elevated,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.Payment, contentDescription = "Payment",
                     tint = TchatColors.primary, modifier = Modifier.size(24.dp))
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Text(
                    text = "Payment Sent",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.xs))
            Text(
                text = "$250.00",
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold
            )
            Text(
                text = "Project consultation",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant
            )
        }
    }
}

/**
 * Poll message component
 */
@Composable
fun PollMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.Poll, contentDescription = "Poll",
                     tint = TchatColors.primary, modifier = Modifier.size(20.dp))
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Text(
                    text = "Poll",
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.primary,
                    fontWeight = FontWeight.Bold
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Text(
                text = "What time works best for our next meeting?",
                style = MaterialTheme.typography.bodyMedium,
                fontWeight = FontWeight.Medium
            )
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            listOf("9:00 AM", "2:00 PM", "4:00 PM").forEach { option ->
                Row(
                    modifier = Modifier.padding(vertical = 2.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    RadioButton(selected = option == "2:00 PM", onClick = { })
                    Spacer(modifier = Modifier.width(TchatSpacing.xs))
                    Text(option, style = MaterialTheme.typography.bodySmall)
                }
            }
        }
    }
}

/**
 * Form message component
 */
@Composable
fun FormMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.Assignment, contentDescription = "Form",
                     tint = TchatColors.primary, modifier = Modifier.size(20.dp))
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Text(
                    text = "Survey Form",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Text(
                text = "Project feedback questionnaire",
                style = MaterialTheme.typography.bodyMedium
            )
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            TchatButton(
                text = "Fill Form",
                variant = TchatButtonVariant.Outline,
                size = TchatButtonSize.Small,
                onClick = { }
            )
        }
    }
}

/**
 * System message component
 */
@Composable
fun SystemMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.Center
    ) {
        TchatCard(
            variant = TchatCardVariant.Filled,
            modifier = Modifier.wrapContentWidth()
        ) {
            Text(
                text = content,
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(TchatSpacing.sm),
                textAlign = androidx.compose.ui.text.style.TextAlign.Center
            )
        }
    }
}

/**
 * Sticker message component
 */
@Composable
fun StickerMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier.size(100.dp),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "👍",
            style = MaterialTheme.typography.displayMedium
        )
    }
}

/**
 * GIF message component
 */
@Composable
fun GifMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.size(200.dp, 150.dp)
    ) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = "🎉 GIF",
                style = MaterialTheme.typography.headlineMedium
            )
        }
    }
}

/**
 * Contact message component
 */
@Composable
fun ContactMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Elevated,
        modifier = modifier.fillMaxWidth()
    ) {
        Row(
            modifier = Modifier.padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .background(TchatColors.primary, CircleShape),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "JS",
                    color = TchatColors.onPrimary,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
            }
            Spacer(modifier = Modifier.width(TchatSpacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = "John Smith",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium
                )
                Text(
                    text = "+1-555-0123",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
            }
            TchatButton(
                text = "Add",
                variant = TchatButtonVariant.Outline,
                size = TchatButtonSize.Small,
                onClick = { /* Add contact action */ }
            )
        }
    }
}

/**
 * Event message component
 */
@Composable
fun EventMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Filled,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(
                    Icons.Default.Event,
                    contentDescription = "Event",
                    tint = TchatColors.primary,
                    modifier = Modifier.size(20.dp)
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Text(
                    text = "Calendar Event",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Text(
                text = "Project deadline - March 15th",
                style = MaterialTheme.typography.bodyMedium,
                fontWeight = FontWeight.Medium
            )
            Text(
                text = "2:00 PM - 3:00 PM",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant
            )
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Row {
                TchatButton(
                    text = "Accept",
                    variant = TchatButtonVariant.Primary,
                    size = TchatButtonSize.Small,
                    onClick = { /* Accept event action */ }
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                TchatButton(
                    text = "Decline",
                    variant = TchatButtonVariant.Outline,
                    size = TchatButtonSize.Small,
                    onClick = { /* Decline event action */ }
                )
            }
        }
    }
}

/**
 * Embed/Link preview message component
 */
@Composable
fun EmbedMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    TchatCard(
        variant = TchatCardVariant.Outlined,
        modifier = modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(TchatSpacing.md)) {
            // Preview image placeholder
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .background(TchatColors.surfaceVariant),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    Icons.Default.Link,
                    contentDescription = "Link Preview",
                    modifier = Modifier.size(48.dp),
                    tint = TchatColors.primary
                )
            }
            Spacer(modifier = Modifier.height(TchatSpacing.sm))
            Text(
                text = "Amazing Design Article",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold
            )
            Text(
                text = "Check out this comprehensive guide to modern UI design patterns and principles...",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                maxLines = 2
            )
            Text(
                text = "designinsights.com",
                style = MaterialTheme.typography.labelSmall,
                color = TchatColors.primary
            )
        }
    }
}

/**
 * Deleted message component
 */
@Composable
fun DeletedMessage(
    content: String,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.Center
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier
                .background(
                    TchatColors.surfaceVariant.copy(alpha = 0.5f),
                    RoundedCornerShape(16.dp)
                )
                .padding(TchatSpacing.sm)
        ) {
            Icon(
                Icons.Default.Block,
                contentDescription = "Deleted",
                tint = TchatColors.onSurfaceVariant,
                modifier = Modifier.size(16.dp)
            )
            Spacer(modifier = Modifier.width(TchatSpacing.xs))
            Text(
                text = "This message was deleted",
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                fontStyle = androidx.compose.ui.text.font.FontStyle.Italic
            )
        }
    }
}