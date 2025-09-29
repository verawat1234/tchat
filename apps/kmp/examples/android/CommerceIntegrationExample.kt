package com.tchat.mobile.commerce.examples.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewmodel.compose.viewModel
import com.tchat.mobile.commerce.data.api.CommerceApiClientImpl
import com.tchat.mobile.commerce.data.models.*
import com.tchat.mobile.commerce.domain.managers.CommerceManager
import com.tchat.mobile.commerce.domain.managers.CommerceManagerImpl
import com.tchat.mobile.commerce.domain.repositories.*
import com.tchat.mobile.commerce.platform.storage.AndroidCommerceStorage
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

/**
 * Android Integration Example for KMP Commerce Module
 * Demonstrates how to integrate the shared KMP commerce functionality
 * into a Jetpack Compose Android application
 */

/**
 * Main Activity - Entry point for the Android app
 */
class CommerceActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Initialize commerce system
        val commerceManager = initializeCommerceManager()

        setContent {
            CommerceTheme {
                CommerceApp(commerceManager = commerceManager)
            }
        }
    }

    private fun initializeCommerceManager(): CommerceManager {
        // Create storage implementation
        val storage = AndroidCommerceStorage(this)

        // Create API client with your backend URL
        val apiClient = CommerceApiClientImpl(
            baseUrl = "https://your-api.com/api/v1"
            // Add your Ktor HttpClient configuration here
        )

        // Create repositories
        val cartRepository = CartRepositoryImpl(apiClient, storage)
        val productRepository = ProductRepositoryImpl(apiClient, storage)
        val categoryRepository = CategoryRepositoryImpl(apiClient, storage)

        // Create main commerce manager
        val manager = CommerceManagerImpl(
            cartRepository = cartRepository,
            productRepository = productRepository,
            categoryRepository = categoryRepository,
            apiClient = apiClient,
            storage = storage
        )

        return manager
    }
}

/**
 * Main app composable with bottom navigation
 */
@Composable
fun CommerceApp(commerceManager: CommerceManager) {
    var selectedTab by remember { mutableStateOf(0) }
    val tabs = listOf("Products", "Cart", "Categories")

    // Initialize commerce manager
    LaunchedEffect(commerceManager) {
        commerceManager.initialize()
    }

    Scaffold(
        bottomBar = {
            NavigationBar {
                tabs.forEachIndexed { index, title ->
                    NavigationBarItem(
                        icon = {
                            when (index) {
                                0 -> Icon(Icons.Default.List, contentDescription = title)
                                1 -> Icon(Icons.Default.ShoppingCart, contentDescription = title)
                                2 -> Icon(Icons.Default.Category, contentDescription = title)
                            }
                        },
                        label = { Text(title) },
                        selected = selectedTab == index,
                        onClick = { selectedTab = index }
                    )
                }
            }
        }
    ) { paddingValues ->
        Box(modifier = Modifier.padding(paddingValues)) {
            when (selectedTab) {
                0 -> ProductListScreen(commerceManager)
                1 -> CartScreen(commerceManager)
                2 -> CategoryListScreen(commerceManager)
            }
        }
    }
}

/**
 * Product list screen
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProductListScreen(commerceManager: CommerceManager) {
    val viewModel: ProductListViewModel = viewModel(
        factory = ProductListViewModelFactory(commerceManager)
    )
    val uiState by viewModel.uiState.collectAsState()

    Column {
        // Search bar
        SearchBar(
            query = uiState.searchQuery,
            onQueryChange = viewModel::onSearchQueryChanged,
            onSearch = viewModel::searchProducts,
            active = false,
            onActiveChange = {},
            placeholder = { Text("Search products...") },
            leadingIcon = { Icon(Icons.Default.Search, contentDescription = "Search") }
        ) {}

        // Product list
        when {
            uiState.isLoading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            }
            uiState.error != null -> {
                ErrorMessage(
                    message = uiState.error,
                    onRetry = viewModel::loadProducts
                )
            }
            else -> {
                LazyColumn(
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(uiState.products) { product ->
                        ProductCard(
                            product = product,
                            onAddToCart = { viewModel.addToCart(product.id) },
                            isAddingToCart = uiState.addingToCartIds.contains(product.id)
                        )
                    }
                }
            }
        }
    }
}

/**
 * Product card component
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProductCard(
    product: Product,
    onAddToCart: () -> Unit,
    isAddingToCart: Boolean
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Product image placeholder
            Surface(
                modifier = Modifier
                    .size(80.dp)
                    .clip(RoundedCornerShape(8.dp)),
                color = MaterialTheme.colorScheme.surfaceVariant
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        Icons.Default.Image,
                        contentDescription = "Product image",
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            Spacer(modifier = Modifier.width(16.dp))

            // Product details
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = product.name,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(4.dp))

                Text(
                    text = product.description,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(8.dp))

                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "${product.currency} ${String.format("%.2f", product.price)}",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.primary
                    )

                    product.compareAtPrice?.let { comparePrice ->
                        if (comparePrice > product.price) {
                            Spacer(modifier = Modifier.width(8.dp))
                            Text(
                                text = "${product.currency} ${String.format("%.2f", comparePrice)}",
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                                textDecoration = TextDecoration.LineThrough
                            )
                        }
                    }
                }

                // Stock status
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = if (product.inventory.isInStock) "In Stock" else "Out of Stock",
                    style = MaterialTheme.typography.labelSmall,
                    color = if (product.inventory.isInStock) Color.Green else Color.Red
                )
            }

            // Add to cart button
            FilledTonalButton(
                onClick = onAddToCart,
                enabled = product.inventory.isInStock && !isAddingToCart,
                modifier = Modifier.size(48.dp),
                shape = CircleShape,
                contentPadding = PaddingValues(0.dp)
            ) {
                if (isAddingToCart) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(16.dp),
                        strokeWidth = 2.dp
                    )
                } else {
                    Icon(Icons.Default.Add, contentDescription = "Add to cart")
                }
            }
        }
    }
}

/**
 * Cart screen
 */
@Composable
fun CartScreen(commerceManager: CommerceManager) {
    val viewModel: CartViewModel = viewModel(
        factory = CartViewModelFactory(commerceManager)
    )
    val uiState by viewModel.uiState.collectAsState()

    Column(
        modifier = Modifier.fillMaxSize()
    ) {
        if (uiState.cart.items.isEmpty()) {
            EmptyCartContent(
                modifier = Modifier.weight(1f)
            )
        } else {
            // Cart items
            LazyColumn(
                modifier = Modifier.weight(1f),
                contentPadding = PaddingValues(16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items(uiState.cart.items) { item ->
                    CartItemCard(
                        item = item,
                        onQuantityChange = { newQuantity ->
                            viewModel.updateQuantity(item.id, newQuantity)
                        },
                        onRemove = { viewModel.removeItem(item.id) }
                    )
                }
            }

            // Cart summary
            CartSummaryCard(
                summary = uiState.cartSummary,
                onCheckout = { viewModel.proceedToCheckout() },
                modifier = Modifier.padding(16.dp)
            )
        }

        if (uiState.isLoading) {
            LinearProgressIndicator(
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
}

/**
 * Cart item card component
 */
@Composable
fun CartItemCard(
    item: CartItem,
    onQuantityChange: (Int) -> Unit,
    onRemove: () -> Unit
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Product image placeholder
            Surface(
                modifier = Modifier
                    .size(60.dp)
                    .clip(RoundedCornerShape(8.dp)),
                color = MaterialTheme.colorScheme.surfaceVariant
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        Icons.Default.Image,
                        contentDescription = "Product image",
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            Spacer(modifier = Modifier.width(12.dp))

            // Product details
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = item.productName,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Text(
                    text = item.businessName,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.primary
                )

                Spacer(modifier = Modifier.height(4.dp))

                Text(
                    text = "${item.currency} ${String.format("%.2f", item.totalPrice)}",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold
                )
            }

            // Quantity controls
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                IconButton(
                    onClick = { onQuantityChange(item.quantity - 1) },
                    enabled = item.quantity > 1
                ) {
                    Icon(Icons.Default.Remove, contentDescription = "Decrease quantity")
                }

                Text(
                    text = item.quantity.toString(),
                    style = MaterialTheme.typography.titleMedium,
                    modifier = Modifier.padding(horizontal = 8.dp)
                )

                IconButton(
                    onClick = { onQuantityChange(item.quantity + 1) }
                ) {
                    Icon(Icons.Default.Add, contentDescription = "Increase quantity")
                }
            }

            // Remove button
            IconButton(onClick = onRemove) {
                Icon(
                    Icons.Default.Delete,
                    contentDescription = "Remove item",
                    tint = MaterialTheme.colorScheme.error
                )
            }
        }
    }
}

/**
 * Cart summary card component
 */
@Composable
fun CartSummaryCard(
    summary: CartSummary,
    onCheckout: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Text(
                text = "Cart Summary",
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Bold
            )

            Spacer(modifier = Modifier.height(12.dp))

            SummaryRow("Items (${summary.totalItems})", summary.formattedSubtotal())

            if (summary.discountAmount > 0) {
                SummaryRow(
                    "Discount",
                    "-${summary.formattedDiscount()}",
                    valueColor = Color.Green
                )
            }

            SummaryRow("Shipping", summary.formattedShipping())
            SummaryRow("Tax", summary.formattedTax())

            HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))

            SummaryRow(
                "Total",
                summary.formattedTotal(),
                titleStyle = MaterialTheme.typography.titleMedium,
                valueStyle = MaterialTheme.typography.titleMedium
            )

            Spacer(modifier = Modifier.height(16.dp))

            Button(
                onClick = onCheckout,
                enabled = summary.totalItems > 0,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Proceed to Checkout")
            }
        }
    }
}

/**
 * Summary row helper component
 */
@Composable
fun SummaryRow(
    title: String,
    value: String,
    titleStyle: androidx.compose.ui.text.TextStyle = MaterialTheme.typography.bodyMedium,
    valueStyle: androidx.compose.ui.text.TextStyle = MaterialTheme.typography.bodyMedium,
    valueColor: Color = MaterialTheme.colorScheme.onSurface
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(text = title, style = titleStyle)
        Text(text = value, style = valueStyle, color = valueColor)
    }
}

/**
 * Empty cart content
 */
@Composable
fun EmptyCartContent(modifier: Modifier = Modifier) {
    Column(
        modifier = modifier.fillMaxWidth(),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Icon(
            Icons.Default.ShoppingCart,
            contentDescription = "Empty cart",
            modifier = Modifier.size(64.dp),
            tint = MaterialTheme.colorScheme.onSurfaceVariant
        )

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = "Your cart is empty",
            style = MaterialTheme.typography.titleLarge,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )

        Spacer(modifier = Modifier.height(8.dp))

        Text(
            text = "Add some products to get started",
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )
    }
}

/**
 * Category list screen
 */
@Composable
fun CategoryListScreen(commerceManager: CommerceManager) {
    val viewModel: CategoryListViewModel = viewModel(
        factory = CategoryListViewModelFactory(commerceManager)
    )
    val uiState by viewModel.uiState.collectAsState()

    when {
        uiState.isLoading -> {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator()
            }
        }
        uiState.error != null -> {
            ErrorMessage(
                message = uiState.error,
                onRetry = viewModel::loadCategories
            )
        }
        else -> {
            LazyColumn(
                contentPadding = PaddingValues(16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items(uiState.categories) { category ->
                    CategoryCard(
                        category = category,
                        onClick = { viewModel.selectCategory(category) }
                    )
                }
            }
        }
    }
}

/**
 * Category card component
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CategoryCard(
    category: Category,
    onClick: () -> Unit
) {
    Card(
        onClick = onClick,
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Category icon
            Surface(
                modifier = Modifier
                    .size(48.dp)
                    .clip(RoundedCornerShape(8.dp)),
                color = MaterialTheme.colorScheme.primaryContainer
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        Icons.Default.Category,
                        contentDescription = "Category",
                        tint = MaterialTheme.colorScheme.onPrimaryContainer
                    )
                }
            }

            Spacer(modifier = Modifier.width(16.dp))

            // Category details
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = category.name,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )

                category.description?.let { description ->
                    Text(
                        text = description,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                }

                Text(
                    text = "${category.productCount} products",
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }

            // Featured indicator
            if (category.isFeatured) {
                Icon(
                    Icons.Default.Star,
                    contentDescription = "Featured",
                    tint = Color(0xFFFFD700) // Gold color
                )
            }

            Icon(
                Icons.Default.ChevronRight,
                contentDescription = "Navigate",
                tint = MaterialTheme.colorScheme.onSurfaceVariant
            )
        }
    }
}

/**
 * Error message component
 */
@Composable
fun ErrorMessage(
    message: String,
    onRetry: () -> Unit
) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Icon(
            Icons.Default.Error,
            contentDescription = "Error",
            modifier = Modifier.size(48.dp),
            tint = MaterialTheme.colorScheme.error
        )

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = "Something went wrong",
            style = MaterialTheme.typography.titleMedium,
            color = MaterialTheme.colorScheme.error
        )

        Spacer(modifier = Modifier.height(8.dp))

        Text(
            text = message,
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )

        Spacer(modifier = Modifier.height(16.dp))

        Button(onClick = onRetry) {
            Text("Retry")
        }
    }
}

/**
 * Theme for the app
 */
@Composable
fun CommerceTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        content = content
    )
}

/**
 * Extension functions for formatting
 */
fun CartSummary.formattedSubtotal(): String {
    return "$currency ${String.format("%.2f", subtotal)}"
}

fun CartSummary.formattedTotal(): String {
    return "$currency ${String.format("%.2f", total)}"
}

fun CartSummary.formattedDiscount(): String {
    return "$currency ${String.format("%.2f", discountAmount)}"
}

fun CartSummary.formattedShipping(): String {
    return "$currency ${String.format("%.2f", shipping)}"
}

fun CartSummary.formattedTax(): String {
    return "$currency ${String.format("%.2f", tax)}"
}

/**
 * ViewModels and their factories would be implemented here
 * For brevity, showing factory signatures only
 */

class ProductListViewModelFactory(
    private val commerceManager: CommerceManager
) : ViewModelProvider.Factory {
    @Suppress("UNCHECKED_CAST")
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return ProductListViewModel(commerceManager) as T
    }
}

class CartViewModelFactory(
    private val commerceManager: CommerceManager
) : ViewModelProvider.Factory {
    @Suppress("UNCHECKED_CAST")
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return CartViewModel(commerceManager) as T
    }
}

class CategoryListViewModelFactory(
    private val commerceManager: CommerceManager
) : ViewModelProvider.Factory {
    @Suppress("UNCHECKED_CAST")
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return CategoryListViewModel(commerceManager) as T
    }
}

/**
 * Simplified UI State and ViewModel classes
 * In a real app, these would be fully implemented with proper state management
 */

data class ProductListUiState(
    val products: List<Product> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val searchQuery: String = "",
    val addingToCartIds: Set<String> = emptySet()
)

data class CartUiState(
    val cart: Cart = Cart(
        id = "",
        sessionId = "",
        items = emptyList(),
        lastActivity = kotlinx.datetime.Clock.System.now(),
        dataRegion = "US",
        createdAt = kotlinx.datetime.Clock.System.now(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    ),
    val cartSummary: CartSummary = CartSummary(
        totalItems = 0,
        subtotal = 0.0,
        discountAmount = 0.0,
        shipping = 0.0,
        tax = 0.0,
        total = 0.0,
        currency = "USD"
    ),
    val isLoading: Boolean = false,
    val error: String? = null
)

data class CategoryListUiState(
    val categories: List<Category> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)

class ProductListViewModel(
    private val commerceManager: CommerceManager
) : ViewModel() {
    private val _uiState = MutableStateFlow(ProductListUiState())
    val uiState: StateFlow<ProductListUiState> = _uiState.asStateFlow()

    fun loadProducts() {
        // Implementation would load products from commerceManager
    }

    fun searchProducts(query: String) {
        // Implementation would search products
    }

    fun onSearchQueryChanged(query: String) {
        _uiState.value = _uiState.value.copy(searchQuery = query)
    }

    fun addToCart(productId: String) {
        // Implementation would add product to cart
    }
}

class CartViewModel(
    private val commerceManager: CommerceManager
) : ViewModel() {
    private val _uiState = MutableStateFlow(CartUiState())
    val uiState: StateFlow<CartUiState> = _uiState.asStateFlow()

    fun updateQuantity(itemId: String, quantity: Int) {
        // Implementation would update item quantity
    }

    fun removeItem(itemId: String) {
        // Implementation would remove item from cart
    }

    fun proceedToCheckout() {
        // Implementation would handle checkout
    }
}

class CategoryListViewModel(
    private val commerceManager: CommerceManager
) : ViewModel() {
    private val _uiState = MutableStateFlow(CategoryListUiState())
    val uiState: StateFlow<CategoryListUiState> = _uiState.asStateFlow()

    fun loadCategories() {
        // Implementation would load categories
    }

    fun selectCategory(category: Category) {
        // Implementation would handle category selection
    }
}