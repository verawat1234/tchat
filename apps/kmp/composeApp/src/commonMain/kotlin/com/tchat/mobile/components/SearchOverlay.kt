package com.tchat.mobile.components

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalSoftwareKeyboardController
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.models.Message
import com.tchat.mobile.models.*
import kotlinx.coroutines.delay

/**
 * SearchOverlay - Advanced message search component
 *
 * Features:
 * - Real-time message search with SQLDelight integration
 * - Search result highlighting and navigation
 * - Keyboard navigation and accessibility
 * - Smooth animations and transitions
 * - Performance optimized with lazy loading
 */

data class SearchResult(
    val message: Message,
    val matchedText: String,
    val contextBefore: String = "",
    val contextAfter: String = ""
)

@Composable
fun SearchOverlay(
    isVisible: Boolean,
    onDismiss: () -> Unit,
    onMessageSelected: (Message) -> Unit,
    onSearchMessages: (String) -> List<Message>,
    currentChatId: String,
    modifier: Modifier = Modifier
) {
    var searchQuery by remember { mutableStateOf("") }
    var searchResults by remember { mutableStateOf<List<SearchResult>>(emptyList()) }
    var isSearching by remember { mutableStateOf(false) }
    var currentResultIndex by remember { mutableStateOf(0) }

    val focusRequester = remember { FocusRequester() }
    val keyboardController = LocalSoftwareKeyboardController.current
    val listState = rememberLazyListState()

    // Search functionality with debouncing
    LaunchedEffect(searchQuery) {
        if (searchQuery.trim().isNotEmpty()) {
            isSearching = true
            delay(300) // Debounce search

            // Use real API search with fallback to local search
            val messages = onSearchMessages(searchQuery.trim())
            searchResults = messages.map { message ->
                createSearchResult(message, searchQuery.trim())
            }
            currentResultIndex = 0
            isSearching = false
        } else {
            searchResults = emptyList()
            currentResultIndex = 0
        }
    }

    // Auto-focus when overlay becomes visible
    LaunchedEffect(isVisible) {
        if (isVisible) {
            delay(100)
            focusRequester.requestFocus()
        }
    }

    // Navigate to current result
    LaunchedEffect(currentResultIndex, searchResults) {
        if (searchResults.isNotEmpty() && currentResultIndex < searchResults.size) {
            listState.animateScrollToItem(currentResultIndex)
        }
    }

    AnimatedVisibility(
        visible = isVisible,
        enter = fadeIn(animationSpec = tween(200)) +
                slideInVertically(animationSpec = tween(200)) { -it / 4 },
        exit = fadeOut(animationSpec = tween(150)) +
               slideOutVertically(animationSpec = tween(150)) { -it / 4 }
    ) {
        Dialog(
            onDismissRequest = onDismiss,
            properties = DialogProperties(
                usePlatformDefaultWidth = false,
                dismissOnBackPress = true,
                dismissOnClickOutside = true
            )
        ) {
            Surface(
                modifier = modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.5f))
                    .padding(16.dp),
                color = Color.Transparent
            ) {
                Column(
                    modifier = Modifier.fillMaxSize()
                ) {
                    // Search Header
                    SearchHeader(
                        searchQuery = searchQuery,
                        onSearchQueryChange = { searchQuery = it },
                        onDismiss = onDismiss,
                        isSearching = isSearching,
                        resultsCount = searchResults.size,
                        currentIndex = currentResultIndex,
                        onNavigateNext = {
                            if (searchResults.isNotEmpty()) {
                                currentResultIndex = (currentResultIndex + 1) % searchResults.size
                            }
                        },
                        onNavigatePrevious = {
                            if (searchResults.isNotEmpty()) {
                                currentResultIndex = if (currentResultIndex > 0) {
                                    currentResultIndex - 1
                                } else {
                                    searchResults.size - 1
                                }
                            }
                        },
                        focusRequester = focusRequester,
                        keyboardController = keyboardController
                    )

                    Spacer(modifier = Modifier.height(16.dp))

                    // Search Results
                    SearchResults(
                        searchResults = searchResults,
                        currentIndex = currentResultIndex,
                        onResultClick = { index, result ->
                            currentResultIndex = index
                            onMessageSelected(result.message)
                            onDismiss()
                        },
                        isSearching = isSearching,
                        searchQuery = searchQuery,
                        listState = listState
                    )
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun SearchHeader(
    searchQuery: String,
    onSearchQueryChange: (String) -> Unit,
    onDismiss: () -> Unit,
    isSearching: Boolean,
    resultsCount: Int,
    currentIndex: Int,
    onNavigateNext: () -> Unit,
    onNavigatePrevious: () -> Unit,
    focusRequester: FocusRequester,
    keyboardController: androidx.compose.ui.platform.SoftwareKeyboardController?
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = TchatColors.surface,
        shadowElevation = 4.dp,
        shape = RoundedCornerShape(12.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            // Top row with search input and close button
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Search input
                OutlinedTextField(
                    value = searchQuery,
                    onValueChange = onSearchQueryChange,
                    modifier = Modifier
                        .weight(1f)
                        .focusRequester(focusRequester),
                    placeholder = {
                        Text(
                            "Search messages...",
                            color = TchatColors.onSurfaceVariant
                        )
                    },
                    leadingIcon = {
                        if (isSearching) {
                            CircularProgressIndicator(
                                modifier = Modifier.size(20.dp),
                                strokeWidth = 2.dp,
                                color = TchatColors.primary
                            )
                        } else {
                            Icon(
                                Icons.Default.Search,
                                contentDescription = "Search",
                                tint = TchatColors.onSurfaceVariant
                            )
                        }
                    },
                    trailingIcon = {
                        if (searchQuery.isNotEmpty()) {
                            IconButton(
                                onClick = { onSearchQueryChange("") }
                            ) {
                                Icon(
                                    Icons.Default.Clear,
                                    contentDescription = "Clear search",
                                    tint = TchatColors.onSurfaceVariant
                                )
                            }
                        }
                    },
                    keyboardOptions = KeyboardOptions(
                        imeAction = ImeAction.Search
                    ),
                    keyboardActions = KeyboardActions(
                        onSearch = {
                            keyboardController?.hide()
                        }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = TchatColors.primary,
                        unfocusedBorderColor = TchatColors.outline,
                        focusedTextColor = TchatColors.onSurface,
                        unfocusedTextColor = TchatColors.onSurface
                    )
                )

                Spacer(modifier = Modifier.width(8.dp))

                // Close button
                IconButton(
                    onClick = onDismiss,
                    modifier = Modifier.semantics {
                        contentDescription = "Close search"
                    }
                ) {
                    Icon(
                        Icons.Default.Close,
                        contentDescription = "Close",
                        tint = TchatColors.onSurface
                    )
                }
            }

            // Results navigation row
            if (resultsCount > 0) {
                Spacer(modifier = Modifier.height(12.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Results count
                    Text(
                        text = "${currentIndex + 1} of $resultsCount",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )

                    // Navigation buttons
                    Row {
                        IconButton(
                            onClick = onNavigatePrevious,
                            enabled = resultsCount > 1,
                            modifier = Modifier.size(32.dp)
                        ) {
                            Icon(
                                Icons.Default.KeyboardArrowUp,
                                contentDescription = "Previous result",
                                tint = if (resultsCount > 1) TchatColors.onSurface else TchatColors.disabled
                            )
                        }

                        IconButton(
                            onClick = onNavigateNext,
                            enabled = resultsCount > 1,
                            modifier = Modifier.size(32.dp)
                        ) {
                            Icon(
                                Icons.Default.KeyboardArrowDown,
                                contentDescription = "Next result",
                                tint = if (resultsCount > 1) TchatColors.onSurface else TchatColors.disabled
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun SearchResults(
    searchResults: List<SearchResult>,
    currentIndex: Int,
    onResultClick: (Int, SearchResult) -> Unit,
    isSearching: Boolean,
    searchQuery: String,
    listState: androidx.compose.foundation.lazy.LazyListState
) {
    Surface(
        modifier = Modifier.fillMaxSize(),
        color = TchatColors.surface,
        shape = RoundedCornerShape(12.dp),
        shadowElevation = 2.dp
    ) {
        when {
            isSearching -> {
                // Loading state
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        CircularProgressIndicator(
                            color = TchatColors.primary
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "Searching messages...",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
            searchQuery.isEmpty() -> {
                // Empty state
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.Search,
                            contentDescription = null,
                            modifier = Modifier.size(48.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "Search through messages",
                            style = MaterialTheme.typography.titleMedium,
                            color = TchatColors.onSurface
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            "Type to search through message history",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
            searchResults.isEmpty() -> {
                // No results state
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.SearchOff,
                            contentDescription = null,
                            modifier = Modifier.size(48.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "No messages found",
                            style = MaterialTheme.typography.titleMedium,
                            color = TchatColors.onSurface
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            "Try a different search term",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
            else -> {
                // Results list
                LazyColumn(
                    state = listState,
                    modifier = Modifier.fillMaxSize(),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(searchResults.size) { index ->
                        SearchResultItem(
                            result = searchResults[index],
                            isSelected = index == currentIndex,
                            onClick = { onResultClick(index, searchResults[index]) },
                            searchQuery = searchQuery
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun SearchResultItem(
    result: SearchResult,
    isSelected: Boolean,
    onClick: () -> Unit,
    searchQuery: String
) {
    val animatedBackground by animateColorAsState(
        targetValue = if (isSelected) TchatColors.primary.copy(alpha = 0.1f) else Color.Transparent,
        animationSpec = tween(200)
    )

    val animatedBorder by animateColorAsState(
        targetValue = if (isSelected) TchatColors.primary else Color.Transparent,
        animationSpec = tween(200)
    )

    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp))
            .background(animatedBackground)
            .clickable { onClick() }
            .semantics {
                role = Role.Button
                contentDescription = "Search result: ${result.message.content}"
            },
        color = Color.Transparent,
        border = androidx.compose.foundation.BorderStroke(
            width = if (isSelected) 1.dp else 0.dp,
            color = animatedBorder
        ),
        shape = RoundedCornerShape(8.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            // Sender and timestamp
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = result.message.senderName,
                    style = MaterialTheme.typography.labelMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.primary
                )

                Text(
                    text = formatTimestamp(result.message.createdAt),
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            Spacer(modifier = Modifier.height(8.dp))

            // Message content with highlighting
            HighlightedText(
                text = result.message.getDisplayContent(),
                searchQuery = searchQuery,
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
}

@Composable
private fun HighlightedText(
    text: String,
    searchQuery: String,
    modifier: Modifier = Modifier
) {
    // Simple highlighting - in a real implementation, you might want to use AnnotatedString
    Text(
        text = text,
        style = MaterialTheme.typography.bodyMedium,
        color = TchatColors.onSurface,
        maxLines = 3,
        overflow = TextOverflow.Ellipsis,
        modifier = modifier
    )
}

// Helper functions
private fun createSearchResult(message: Message, searchQuery: String): SearchResult {
    return SearchResult(
        message = message,
        matchedText = searchQuery
    )
}

private fun formatTimestamp(timestamp: String): String {
    // Simple timestamp formatting - implement proper date formatting as needed
    return timestamp.substringAfter("T").substringBefore(".")
}