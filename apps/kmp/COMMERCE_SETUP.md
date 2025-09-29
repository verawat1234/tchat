# KMP Commerce Module Setup Guide

This guide provides step-by-step instructions for setting up and using the Kotlin Multiplatform (KMP) Commerce module in your iOS and Android applications.

## 🚀 Quick Start

The KMP Commerce module is already integrated into your existing Tchat project structure. Here's how to use it:

### Prerequisites

✅ Kotlin Multiplatform 1.9.23+
✅ Compose Multiplatform 1.6.10+
✅ Go backend commerce service running on port 8083
✅ SQLDelight for local storage
✅ Ktor Client for networking

## 📁 Project Structure

Your commerce module is organized as follows:

```
apps/kmp/composeApp/src/
├── commonMain/kotlin/com/tchat/mobile/commerce/
│   ├── data/
│   │   ├── api/           # API client and networking
│   │   └── models/        # Shared data models
│   ├── domain/
│   │   ├── managers/      # High-level business logic coordinators
│   │   └── repositories/  # Data access layer
│   ├── offline/           # Offline operations support
│   ├── sync/              # Cross-platform synchronization
│   └── platform/storage/  # Storage abstraction
├── androidMain/kotlin/com/tchat/mobile/commerce/
│   ├── platform/storage/  # Android-specific storage
│   └── presentation/      # Android UI components and ViewModels
├── iosMain/kotlin/com/tchat/mobile/commerce/
│   ├── platform/storage/  # iOS-specific storage
│   └── presentation/      # iOS ViewModels and integration
└── commonTest/kotlin/     # Comprehensive test suite
```

## 🔧 Initial Setup

### 1. Build the KMP Module

```bash
cd /Users/weerawat/Tchat/apps/kmp

# Build for all platforms
./gradlew build

# Generate iOS framework specifically
./gradlew linkReleaseFrameworkIosArm64 linkReleaseFrameworkIosX64

# Run tests
./gradlew test
```

### 2. Backend Service Configuration

Ensure your Go backend commerce service is running:

```bash
# From your backend directory
cd /Users/weerawat/Tchat/backend/commerce
go run main.go
```

The service should be accessible at:
- Direct: `http://localhost:8083/api/v1/commerce/*`
- Via Gateway: `http://localhost:8080/api/v1/commerce/*`

## 📱 Platform Integration

### iOS Integration

#### 1. Add Framework to iOS Project

The KMP framework is generated at:
```
apps/kmp/composeApp/build/bin/ios*/releaseFramework/ComposeApp.framework
```

Add this framework to your iOS project and configure:

```swift
// In your iOS app
import ComposeApp

// Initialize commerce in your app
@main
struct TchatApp: App {
    @StateObject private var commerceManager = CommerceManager()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(commerceManager)
        }
    }
}
```

#### 2. iOS-Specific Storage Setup

The iOS storage implementation uses:
- **UserDefaults**: For general app preferences
- **Keychain**: For secure cart and session data
- **CoreData**: For offline commerce data

#### 3. iOS SwiftUI Integration Example

```swift
struct ProductListView: View {
    @EnvironmentObject var commerceManager: CommerceManager
    @State private var products: [Product] = []

    var body: some View {
        List(products, id: \.id) { product in
            ProductRow(product: product) {
                Task {
                    await commerceManager.addToCart(productId: product.id)
                }
            }
        }
        .onAppear {
            Task {
                products = await commerceManager.getProducts()
            }
        }
    }
}
```

### Android Integration

#### 1. Add Dependencies

The dependencies are already configured in `build.gradle.kts`. Key additions include:

```kotlin
// Android storage and security
implementation("androidx.security:security-crypto:1.1.0-alpha06")
implementation("androidx.datastore:datastore-preferences:1.0.0")

// Ktor client for Android
implementation("io.ktor:ktor-client-android:2.3.7")
```

#### 2. Android-Specific Storage Setup

The Android storage implementation uses:
- **SharedPreferences**: For general app settings
- **EncryptedSharedPreferences**: For secure cart and user data
- **DataStore**: For structured preference storage

#### 3. Android Compose Integration Example

```kotlin
@Composable
fun ProductListScreen(commerceManager: CommerceManager) {
    val products by commerceManager.products.collectAsState()

    LazyColumn {
        items(products) { product ->
            ProductCard(
                product = product,
                onAddToCart = {
                    commerceManager.addToCart(product.id)
                }
            )
        }
    }
}
```

## 🔄 Commerce Operations

### Core Features Available

#### 1. Cart Management
```kotlin
// Add item to cart
commerceManager.addToCart(productId = "product-123", quantity = 2)

// Update cart item
commerceManager.updateCartItem(itemId = "item-456", quantity = 3)

// Remove from cart
commerceManager.removeFromCart(itemId = "item-456")

// Get current cart
val cart = commerceManager.getCurrentCart()
```

#### 2. Product Browsing
```kotlin
// Get all products
val products = commerceManager.getProducts()

// Search products
val searchResults = commerceManager.searchProducts("electronics")

// Get featured products
val featured = commerceManager.getFeaturedProducts()

// Get product by ID
val product = commerceManager.getProduct("product-123")
```

#### 3. Category Management
```kotlin
// Get all categories
val categories = commerceManager.getCategories()

// Get category tree
val categoryTree = commerceManager.getCategoryTree()

// Get products by category
val categoryProducts = commerceManager.getProductsByCategory("electronics")
```

#### 4. Offline Support
```kotlin
// Check if offline mode is enabled
val isOffline = commerceManager.isOfflineMode

// Enable offline operations
commerceManager.enableOfflineMode()

// Queue offline operations
commerceManager.addToCartOffline(productId, quantity)

// Sync when back online
commerceManager.syncWithServer()
```

## 🔧 Configuration

### API Client Configuration

```kotlin
// Configure the API client
val apiClient = CommerceApiClientImpl(
    baseUrl = "https://your-api.com/api/v1",
    timeout = 30000,
    enableLogging = true,
    retryCount = 3
)
```

### Storage Configuration

#### Android Storage
```kotlin
val storage = AndroidCommerceStorage(
    context = applicationContext,
    enableEncryption = true,
    cacheSize = 10 * 1024 * 1024 // 10MB
)
```

#### iOS Storage
```swift
let storage = IOSCommerceStorage(
    useKeychain: true,
    cacheSize: 10 * 1024 * 1024 // 10MB
)
```

### Sync Configuration

```kotlin
val syncManager = SyncManagerImpl(
    autoSyncEnabled = true,
    syncInterval = 5 * 60 * 1000, // 5 minutes
    retryFailedOperations = true,
    maxRetries = 3
)
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
./gradlew test

# Run specific test suites
./gradlew testDebugUnitTest           # Android unit tests
./gradlew commonTest                  # Common/shared tests
./gradlew iosSimulatorArm64Test      # iOS simulator tests
```

### Test Coverage

The test suite includes:

✅ **CommerceManagerTest**: Core manager functionality
✅ **CartRepositoryTest**: Cart operations and state management
✅ **SyncManagerTest**: Cross-platform synchronization
✅ **OfflineManagerTest**: Offline operations and queuing
✅ **Mock Implementations**: Complete mock API and storage layers

### Test Structure

```kotlin
class CommerceManagerTest {
    @Test
    fun testAddToCart() = runTest {
        // Setup
        commerceManager.initialize()

        // Execute
        val result = commerceManager.addToCart("product-1", quantity = 2)

        // Verify
        assertTrue(result.isSuccess)
        assertEquals(2, result.getOrNull()?.items?.first()?.quantity)
    }
}
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Build Errors

**Issue**: Framework not found for iOS
**Solution**:
```bash
./gradlew clean
./gradlew linkReleaseFrameworkIosArm64
```

#### 2. Network Connection Issues

**Issue**: API calls failing
**Solution**: Check backend service and network configuration:
```kotlin
// Verify API endpoint
val response = commerceManager.testConnection()
```

#### 3. Storage Issues

**Issue**: Data not persisting
**Solution**: Check storage permissions and configuration:

**Android**: Ensure app has storage permissions
**iOS**: Verify Keychain access is configured

#### 4. State Not Updating

**Issue**: UI not reflecting changes
**Solution**: Ensure StateFlow collection is properly set up:

```kotlin
// Android
val cart by commerceManager.cart.collectAsState()

// iOS
commerceManager.cart.collect { cart in
    // Update UI
}
```

### Debug Mode

Enable debug logging for troubleshooting:

```kotlin
val commerceManager = CommerceManagerImpl(
    // ... other dependencies
    enableDebugLogging = true
)
```

## 🚀 Performance Optimization

### Caching Strategy

The module implements intelligent caching:

- **Memory Cache**: Hot data (current cart, recent products)
- **Disk Cache**: Persistent data (product catalog, user preferences)
- **Network Cache**: API response caching with TTL

### Background Sync

Automatic background synchronization:

```kotlin
// Configure background sync
commerceManager.enableAutoSync(
    interval = 5 * 60 * 1000, // 5 minutes
    onlyOnWiFi = true,
    requiresCharging = false
)
```

### Offline-First Architecture

The module prioritizes offline functionality:

1. **Local Operations**: All operations work offline first
2. **Queue Management**: Failed operations are queued for retry
3. **Conflict Resolution**: Smart conflict resolution when syncing
4. **Graceful Degradation**: Reduced functionality when fully offline

## 📊 Monitoring and Analytics

### Commerce Stats

Get real-time commerce statistics:

```kotlin
val stats = commerceManager.getCommerceStats()
println("Cart items: ${stats.cartItems}")
println("Cart value: ${stats.cartValue}")
println("Total products: ${stats.totalProducts}")
```

### Sync Status

Monitor synchronization status:

```kotlin
commerceManager.syncStatus.collect { status ->
    when (status) {
        SyncStatus.SYNCING -> showSyncIndicator()
        SyncStatus.SUCCESS -> hideSyncIndicator()
        SyncStatus.ERROR -> showSyncError()
    }
}
```

## 🔒 Security Considerations

### Data Protection

- **Encryption**: Sensitive data encrypted at rest
- **Secure Storage**: Platform-appropriate secure storage
- **Token Management**: Automatic token refresh and secure storage
- **Network Security**: TLS/SSL for all API communications

### Privacy

- **Data Minimization**: Only essential data is cached
- **User Consent**: Respect user privacy preferences
- **Data Retention**: Automatic cleanup of old cached data

## 📈 Scaling Considerations

### Performance Targets

- **API Response Time**: < 200ms for cached data
- **UI Update Time**: < 16ms for smooth animations
- **Memory Usage**: < 100MB on mobile devices
- **Storage Size**: < 50MB cached data per user

### Optimization Strategies

1. **Lazy Loading**: Load data as needed
2. **Pagination**: Paginate large product lists
3. **Image Optimization**: Efficient image loading and caching
4. **Background Processing**: Heavy operations on background threads

## 🔄 Migration and Updates

### Version Compatibility

The module maintains backward compatibility with:
- API version changes
- Data model evolution
- Storage format updates

### Migration Strategy

```kotlin
// Handle data migration
commerceManager.migrateData(
    fromVersion = "1.0",
    toVersion = "1.1"
)
```

## 📞 Support and Documentation

### Additional Resources

- **API Documentation**: Backend commerce service documentation
- **KMP Official Docs**: [Kotlin Multiplatform](https://kotlinlang.org/docs/multiplatform.html)
- **Compose Multiplatform**: [JetBrains Compose](https://www.jetbrains.com/lp/compose-multiplatform/)

### Getting Help

1. **Check this documentation** for common issues and solutions
2. **Review test examples** for usage patterns
3. **Enable debug logging** for detailed troubleshooting
4. **Check backend service logs** for API-related issues

## ✅ Success Checklist

Before going to production, ensure:

- [ ] All tests pass (`./gradlew test`)
- [ ] iOS framework builds successfully
- [ ] Android app builds and runs
- [ ] Backend service is accessible
- [ ] Storage configuration is secure
- [ ] Offline functionality works as expected
- [ ] Sync operations complete successfully
- [ ] Error handling is comprehensive
- [ ] Performance targets are met
- [ ] Security requirements are satisfied

## 🎯 Next Steps

1. **Customize**: Adapt the UI components to match your design system
2. **Extend**: Add additional commerce features as needed
3. **Monitor**: Implement analytics and monitoring
4. **Optimize**: Profile and optimize for your specific use cases
5. **Scale**: Plan for increased load and data volume

---

**🎉 Your KMP Commerce module is ready to use!**

The shared business logic provides a solid foundation for cross-platform commerce functionality, with platform-specific optimizations for both iOS and Android.