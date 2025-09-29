package com.tchat.mobile.stream.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.rounded.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.stream.models.StreamCategory
import com.tchat.mobile.stream.models.StreamSubtab

/**
 * Stream Tab Navigation Component
 * Material3-based tab navigation for Stream categories and subtabs
 */

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StreamTabs(
    categories: List<StreamCategory>,
    selectedCategoryId: String,
    selectedSubtabId: String?,
    onCategorySelected: (String) -> Unit,
    onSubtabSelected: (String?) -> Unit,
    isLoading: Boolean = false,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxWidth()
            .background(MaterialTheme.colorScheme.surface)
    ) {
        // Main Category Tabs
        if (isLoading) {
            CategoryTabsLoadingState()
        } else {
            CategoryTabs(
                categories = categories,
                selectedCategoryId = selectedCategoryId,
                onCategorySelected = onCategorySelected
            )
        }

        // Subtabs for selected category
        val selectedCategory = categories.find { it.id == selectedCategoryId }
        selectedCategory?.subtabs?.let { subtabs ->
            if (subtabs.isNotEmpty()) {
                Spacer(modifier = Modifier.height(8.dp))
                SubtabRow(
                    subtabs = subtabs,
                    selectedSubtabId = selectedSubtabId,
                    onSubtabSelected = onSubtabSelected
                )
            }
        }
    }
}

@Composable
private fun CategoryTabs(
    categories: List<StreamCategory>,
    selectedCategoryId: String,
    onCategorySelected: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp)
    ) {
        items(categories.filter { it.isActive }) { category ->
            CategoryTab(
                category = category,
                isSelected = category.id == selectedCategoryId,
                onClick = { onCategorySelected(category.id) }
            )
        }
    }
}

@Composable
private fun CategoryTab(
    category: StreamCategory,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    val containerColor = if (isSelected) {
        MaterialTheme.colorScheme.primary
    } else {
        MaterialTheme.colorScheme.surfaceVariant
    }

    val contentColor = if (isSelected) {
        MaterialTheme.colorScheme.onPrimary
    } else {
        MaterialTheme.colorScheme.onSurfaceVariant
    }

    FilterChip(
        onClick = onClick,
        label = {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Icon(
                    imageVector = getCategoryIcon(category.iconName),
                    contentDescription = null,
                    modifier = Modifier.size(18.dp)
                )
                Text(
                    text = category.name,
                    fontSize = 14.sp,
                    fontWeight = if (isSelected) FontWeight.Medium else FontWeight.Normal
                )
            }
        },
        selected = isSelected,
        colors = FilterChipDefaults.filterChipColors(
            containerColor = containerColor,
            labelColor = contentColor,
            iconColor = contentColor,
            selectedContainerColor = MaterialTheme.colorScheme.primary,
            selectedLabelColor = MaterialTheme.colorScheme.onPrimary
        ),
        border = if (isSelected) null else FilterChipDefaults.filterChipBorder(
            borderColor = MaterialTheme.colorScheme.outline,
            enabled = true,
            selected = isSelected
        ),
        modifier = modifier
    )
}

@Composable
private fun SubtabRow(
    subtabs: List<StreamSubtab>,
    selectedSubtabId: String?,
    onSubtabSelected: (String?) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 4.dp)
    ) {
        // "All" option
        item {
            SubtabChip(
                text = "All",
                isSelected = selectedSubtabId == null,
                onClick = { onSubtabSelected(null) }
            )
        }

        // Individual subtabs
        items(subtabs.filter { it.isActive }) { subtab ->
            SubtabChip(
                text = subtab.name,
                isSelected = subtab.id == selectedSubtabId,
                onClick = { onSubtabSelected(subtab.id) }
            )
        }
    }
}

@Composable
private fun SubtabChip(
    text: String,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    AssistChip(
        onClick = onClick,
        label = {
            Text(
                text = text,
                fontSize = 12.sp,
                fontWeight = if (isSelected) FontWeight.Medium else FontWeight.Normal
            )
        },
        colors = AssistChipDefaults.assistChipColors(
            containerColor = if (isSelected) {
                MaterialTheme.colorScheme.secondaryContainer
            } else {
                MaterialTheme.colorScheme.surface
            },
            labelColor = if (isSelected) {
                MaterialTheme.colorScheme.onSecondaryContainer
            } else {
                MaterialTheme.colorScheme.onSurface
            }
        ),
        border = AssistChipDefaults.assistChipBorder(
            borderColor = if (isSelected) {
                MaterialTheme.colorScheme.secondary
            } else {
                MaterialTheme.colorScheme.outline
            },
            enabled = true
        ),
        modifier = modifier
    )
}

@Composable
private fun CategoryTabsLoadingState(
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp)
    ) {
        items(5) { // Show 5 loading placeholders
            Card(
                modifier = Modifier
                    .width(100.dp)
                    .height(40.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.6f)
                ),
                shape = RoundedCornerShape(20.dp)
            ) {
                // Empty placeholder
            }
        }
    }
}

/**
 * Maps category icon names to Material Icons
 */
private fun getCategoryIcon(iconName: String): ImageVector {
    return when (iconName.lowercase()) {
        "book", "books" -> Icons.Rounded.MenuBook
        "podcast", "podcasts" -> Icons.Rounded.Mic
        "cartoon", "cartoons" -> Icons.Rounded.Animation
        "movie", "movies" -> Icons.Rounded.Movie
        "music" -> Icons.Rounded.MusicNote
        "art" -> Icons.Rounded.Palette
        "video" -> Icons.Rounded.PlayCircle
        "audio" -> Icons.Rounded.Headphones
        else -> Icons.Rounded.Category
    }
}

/**
 * Stream Category Filter Configuration
 */
@Composable
fun StreamCategoryFilter(
    categories: List<StreamCategory>,
    selectedCategories: Set<String>,
    onCategoryToggle: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        contentPadding = PaddingValues(horizontal = 16.dp)
    ) {
        items(categories.filter { it.isActive }) { category ->
            FilterChip(
                onClick = { onCategoryToggle(category.id) },
                label = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(6.dp)
                    ) {
                        Icon(
                            imageVector = getCategoryIcon(category.iconName),
                            contentDescription = null,
                            modifier = Modifier.size(16.dp)
                        )
                        Text(
                            text = category.name,
                            fontSize = 12.sp
                        )
                    }
                },
                selected = category.id in selectedCategories,
                colors = FilterChipDefaults.filterChipColors(
                    selectedContainerColor = MaterialTheme.colorScheme.primary,
                    selectedLabelColor = MaterialTheme.colorScheme.onPrimary
                )
            )
        }
    }
}

/**
 * Stream Tab Navigation State
 */
@Stable
class StreamTabState(
    initialCategoryId: String = "",
    initialSubtabId: String? = null
) {
    var selectedCategoryId by mutableStateOf(initialCategoryId)
        private set

    var selectedSubtabId by mutableStateOf(initialSubtabId)
        private set

    fun selectCategory(categoryId: String) {
        if (selectedCategoryId != categoryId) {
            selectedCategoryId = categoryId
            selectedSubtabId = null // Reset subtab when category changes
        }
    }

    fun selectSubtab(subtabId: String?) {
        selectedSubtabId = subtabId
    }

    fun reset() {
        selectedCategoryId = ""
        selectedSubtabId = null
    }
}

/**
 * Remember Stream Tab State
 */
@Composable
fun rememberStreamTabState(
    initialCategoryId: String = "",
    initialSubtabId: String? = null
): StreamTabState {
    return remember {
        StreamTabState(
            initialCategoryId = initialCategoryId,
            initialSubtabId = initialSubtabId
        )
    }
}