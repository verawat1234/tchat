package com.tchat.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * E-commerce store interface screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StoreScreen() {
    var searchText by remember { mutableStateOf("") }
    var selectedCategory by remember { mutableStateOf("All") }
    var cartItems by remember { mutableStateOf(3) }

    // Mock categories
    val categories = listOf("All", "Electronics", "Fashion", "Home", "Books", "Sports")

    // Mock products
    val products = listOf(
        Product("iPhone 15 Pro", "$999", "electronics", Icons.Default.PhoneIphone),
        Product("MacBook Air", "$1199", "electronics", Icons.Default.Laptop),
        Product("AirPods Pro", "$249", "electronics", Icons.Default.Headphones),
        Product("Nike Sneakers", "$129", "fashion", Icons.Default.DirectionsWalk),
        Product("Coffee Maker", "$89", "home", Icons.Default.Coffee),
        Product("Wireless Mouse", "$59", "electronics", Icons.Default.Mouse),
        Product("Yoga Mat", "$39", "sports", Icons.Default.FitnessCenter),
        Product("Cookbook", "$24", "books", Icons.Default.Book)
    )

    val filteredProducts = if (selectedCategory == "All") {
        products
    } else {
        products.filter { it.category == selectedCategory.lowercase() }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Colors.background)
    ) {
        // Top app bar
        TopAppBar(
            title = {
                Text(
                    text = "Store",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.textPrimary
                )
            },
            actions = {
                Box {
                    IconButton(onClick = { /* Open cart */ }) {
                        Icon(
                            imageVector = Icons.Default.ShoppingCart,
                            contentDescription = "Cart",
                            tint = Colors.primary
                        )
                    }
                    if (cartItems > 0) {
                        Box(
                            modifier = Modifier
                                .size(16.dp)
                                .background(
                                    color = Colors.error,
                                    shape = RoundedCornerShape(8.dp)
                                )
                                .offset(x = 10.dp, y = (-10).dp),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = cartItems.toString(),
                                fontSize = 10.sp,
                                fontWeight = FontWeight.Bold,
                                color = androidx.compose.ui.graphics.Color.White
                            )
                        }
                    }
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = Colors.background
            )
        )

        // Search bar
        OutlinedTextField(
            value = searchText,
            onValueChange = { searchText = it },
            placeholder = {
                Text(
                    text = "Search products",
                    color = Colors.textSecondary
                )
            },
            leadingIcon = {
                Icon(
                    imageVector = Icons.Default.Search,
                    contentDescription = "Search",
                    tint = Colors.textSecondary
                )
            },
            trailingIcon = {
                IconButton(onClick = { /* Voice search */ }) {
                    Icon(
                        imageVector = Icons.Default.Mic,
                        contentDescription = "Voice search",
                        tint = Colors.primary
                    )
                }
            },
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = Spacing.md, vertical = Spacing.sm),
            shape = RoundedCornerShape(12.dp),
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = Colors.primary,
                unfocusedBorderColor = Colors.border
            )
        )

        // Category selector
        LazyRow(
            modifier = Modifier.fillMaxWidth(),
            contentPadding = PaddingValues(horizontal = Spacing.md),
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            items(categories) { category ->
                CategoryChip(
                    title = category,
                    isSelected = selectedCategory == category,
                    onClick = { selectedCategory = category }
                )
            }
        }

        Spacer(modifier = Modifier.height(Spacing.sm))

        // Products grid
        LazyVerticalGrid(
            columns = GridCells.Fixed(2),
            modifier = Modifier.fillMaxWidth(),
            contentPadding = PaddingValues(Spacing.md),
            horizontalArrangement = Arrangement.spacedBy(Spacing.md),
            verticalArrangement = Arrangement.spacedBy(Spacing.md)
        ) {
            items(filteredProducts) { product ->
                ProductCard(product = product)
            }
        }
    }
}

// MARK: - Data Classes
data class Product(
    val name: String,
    val price: String,
    val category: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
)

// MARK: - Category Chip Component
@Composable
private fun CategoryChip(
    title: String,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Box(
        modifier = Modifier
            .clickable { onClick() }
            .background(
                color = if (isSelected) Colors.primary else Colors.surface,
                shape = RoundedCornerShape(20.dp)
            )
            .padding(horizontal = Spacing.md, vertical = Spacing.xs),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = title,
            fontSize = 14.sp,
            fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Medium,
            color = if (isSelected) Colors.textOnPrimary else Colors.textSecondary
        )
    }
}

// MARK: - Product Card Component
@Composable
private fun ProductCard(product: Product) {
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
            modifier = Modifier.padding(Spacing.sm)
        ) {
            // Product icon
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(80.dp)
                    .background(
                        color = Colors.surface,
                        shape = RoundedCornerShape(12.dp)
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = product.icon,
                    contentDescription = product.name,
                    tint = Colors.primary,
                    modifier = Modifier.size(32.dp)
                )
            }

            Spacer(modifier = Modifier.height(Spacing.sm))

            // Product info
            Text(
                text = product.name,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = Colors.textPrimary,
                maxLines = 2,
                modifier = Modifier.height(40.dp)
            )

            Spacer(modifier = Modifier.height(Spacing.xs))

            Text(
                text = product.price,
                fontSize = 16.sp,
                fontWeight = FontWeight.Bold,
                color = Colors.primary
            )

            Spacer(modifier = Modifier.height(Spacing.sm))

            // Add to cart button
            Button(
                onClick = { /* Add to cart */ },
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(8.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = Colors.primary
                )
            ) {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        imageVector = Icons.Default.Add,
                        contentDescription = "Add",
                        tint = Colors.textOnPrimary,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(Spacing.xs))
                    Text(
                        text = "Add",
                        fontSize = 12.sp,
                        fontWeight = FontWeight.SemiBold,
                        color = Colors.textOnPrimary
                    )
                }
            }
        }
    }
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun StoreScreenPreview() {
    StoreScreen()
}