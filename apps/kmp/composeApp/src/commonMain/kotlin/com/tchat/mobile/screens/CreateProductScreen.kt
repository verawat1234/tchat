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
fun CreateProductScreen(
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var productName by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var price by remember { mutableStateOf("") }
    var category by remember { mutableStateOf("") }
    var condition by remember { mutableStateOf("") }
    var location by remember { mutableStateOf("") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Create Product Listing", fontWeight = FontWeight.Bold) },
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
                    if (productName.isNotEmpty() && price.isNotEmpty()) {
                        TextButton(
                            onClick = {
                                // TODO: Create product listing
                                onBackClick()
                            }
                        ) {
                            Text(
                                "Publish",
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
            // Photo Upload Section
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md)
                    .height(200.dp),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                onClick = {
                    // TODO: Open image picker
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
                            Icons.Default.PhotoCamera,
                            contentDescription = "Add Photos",
                            modifier = Modifier.size(48.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = "Add Photos",
                            style = MaterialTheme.typography.titleMedium,
                            color = TchatColors.onSurfaceVariant
                        )
                        Text(
                            text = "Tap to add up to 10 photos",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            // Product Details
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
                        text = "Product Details",
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Product Name
                    TchatInput(
                        value = productName,
                        onValueChange = { productName = it },
                        type = TchatInputType.Text,
                        label = "Product Name",
                        placeholder = "Enter product name",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Description
                    TchatTextarea(
                        value = description,
                        onValueChange = { description = it },
                        label = "Description",
                        placeholder = "Describe your product...",
                        modifier = Modifier.fillMaxWidth(),
                        minLines = 3
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Price
                    TchatInput(
                        value = price,
                        onValueChange = { price = it },
                        type = TchatInputType.Number,
                        label = "Price (฿)",
                        placeholder = "0.00",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Category
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("electronics", "Electronics"),
                            SelectOption("clothing_fashion", "Clothing & Fashion"),
                            SelectOption("home_garden", "Home & Garden"),
                            SelectOption("sports_outdoors", "Sports & Outdoors"),
                            SelectOption("books_media", "Books & Media"),
                            SelectOption("toys_games", "Toys & Games"),
                            SelectOption("automotive", "Automotive"),
                            SelectOption("health_beauty", "Health & Beauty"),
                            SelectOption("food_beverages", "Food & Beverages"),
                            SelectOption("other", "Other")
                        ),
                        selectedValue = category,
                        onSelectionChange = { category = it ?: "" },
                        label = "Category",
                        placeholder = "Select category",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Condition
                    TchatSingleSelect(
                        options = listOf(
                            SelectOption("new", "New"),
                            SelectOption("like_new", "Like New"),
                            SelectOption("good", "Good"),
                            SelectOption("fair", "Fair"),
                            SelectOption("poor", "Poor")
                        ),
                        selectedValue = condition,
                        onSelectionChange = { condition = it ?: "" },
                        label = "Condition",
                        placeholder = "Select condition",
                        modifier = Modifier.fillMaxWidth()
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Location
                    TchatInput(
                        value = location,
                        onValueChange = { location = it },
                        type = TchatInputType.Text,
                        label = "Location",
                        placeholder = "Where is this item located?",
                        modifier = Modifier.fillMaxWidth(),
                        trailingIcon = Icons.Default.LocationOn
                    )
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // Shipping Options
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
                        text = "Shipping & Pickup",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    // Shipping Options
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.md)
                    ) {
                        var offerShipping by remember { mutableStateOf(false) }
                        var allowPickup by remember { mutableStateOf(true) }

                        Row(
                            modifier = Modifier.weight(1f),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = offerShipping,
                                onCheckedChange = { offerShipping = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Offer Shipping",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }

                        Row(
                            modifier = Modifier.weight(1f),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Checkbox(
                                checked = allowPickup,
                                onCheckedChange = { allowPickup = it },
                                colors = CheckboxDefaults.colors(
                                    checkedColor = TchatColors.primary
                                )
                            )
                            Text(
                                text = "Allow Pickup",
                                style = MaterialTheme.typography.bodyMedium
                            )
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xl))

            // Submit Button
            TchatButton(
                text = "Create Listing",
                variant = TchatButtonVariant.Primary,
                onClick = {
                    // TODO: Create product listing
                    onBackClick()
                },
                enabled = productName.isNotEmpty() && price.isNotEmpty(),
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = TchatSpacing.md)
            )

            Spacer(modifier = Modifier.height(TchatSpacing.xl))
        }
    }
}