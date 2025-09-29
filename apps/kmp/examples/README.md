# KMP Commerce Integration Examples

This directory contains comprehensive examples showing how to integrate the KMP Commerce module into real iOS and Android applications.

## Overview

The KMP Commerce module provides a shared commerce solution that can be used across iOS and Android platforms. This includes:

- **Shared Business Logic**: Cart management, product browsing, category management
- **Cross-Platform State Sync**: Real-time synchronization between platforms
- **Offline Support**: Local operations with automatic sync when online
- **Type-Safe APIs**: Full TypeScript-style type safety across platforms

## Examples Included

### 1. iOS SwiftUI Integration (`ios/CommerceIntegrationExample.swift`)

Complete SwiftUI application demonstrating:
- **App Structure**: Main app with proper commerce initialization
- **Navigation**: Tab-based navigation with cart, products, and categories
- **State Management**: Reactive UI updates using KMP StateFlow integration
- **SwiftUI Components**: Native iOS components using KMP data
- **Error Handling**: Proper error states and retry mechanisms
- **Offline Support**: Graceful offline behavior

**Key Features:**
- Native SwiftUI navigation and UI patterns
- Proper iOS app lifecycle integration
- StateFlow to SwiftUI binding
- Platform-native formatting and localization
- Accessibility support

### 2. Android Compose Integration (`android/CommerceIntegrationExample.kt`)

Complete Jetpack Compose application demonstrating:
- **App Structure**: ComponentActivity with proper commerce setup
- **Navigation**: Bottom navigation with Material 3 design
- **State Management**: ViewModel integration with KMP repositories
- **Compose Components**: Material Design 3 components using KMP data
- **Error Handling**: Comprehensive error states and user feedback
- **Offline Support**: Offline-first architecture

**Key Features:**
- Material Design 3 components and theming
- Proper Android ViewModel lifecycle
- StateFlow to Compose state integration
- Platform-native Android patterns
- Dependency injection setup

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Platform UI Layer                        │
├─────────────────────┬───────────────────────────────────────┤
│   iOS SwiftUI       │        Android Compose               │
│   - Views           │        - Composables                  │
│   - ViewModels      │        - ViewModels                   │
│   - Navigation      │        - Navigation                   │
├─────────────────────┴───────────────────────────────────────┤
│              KMP Commerce Shared Module                      │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                Managers Layer                        │   │
│   │   - CommerceManager (main coordinator)              │   │
│   │   - SyncManager (cross-platform sync)               │   │
│   │   - OfflineManager (offline operations)             │   │
│   └─────────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────────┐   │
│   │              Repositories Layer                      │   │
│   │   - CartRepository (cart operations)                │   │
│   │   - ProductRepository (product data)                │   │
│   │   - CategoryRepository (category data)              │   │
│   └─────────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                Data Layer                           │   │
│   │   - API Client (Ktor HTTP client)                   │   │
│   │   - Data Models (shared across platforms)           │   │
│   │   - Storage Interface (platform-agnostic)           │   │
│   └─────────────────────────────────────────────────────┘   │
├─────────────────────┬───────────────────────────────────────┤
│ iOS Storage         │        Android Storage                │
│ - UserDefaults      │        - SharedPreferences            │
│ - Keychain          │        - EncryptedSharedPreferences   │
│ - CoreData          │        - Room Database                │
└─────────────────────┴───────────────────────────────────────┘
```

## Getting Started

### Prerequisites

1. **KMP Project Setup**: Ensure you have a proper KMP project structure
2. **Gradle Configuration**: Kotlin Multiplatform plugin configured
3. **Dependencies**: Required dependencies for networking and storage
4. **Backend API**: Running Go commerce backend service

### iOS Integration Steps

1. **Add KMP Framework**: Include the compiled KMP framework in your iOS project
```swift
import CommerceKMP
```

2. **Initialize Commerce**: Set up the commerce system in your app delegate or main app
```swift
let commerceManager = initializeCommerceManager()
```

3. **Create Storage**: Implement iOS-specific storage
```swift
let storage = IOSCommerceStorage()
```

4. **Connect UI**: Use the provided SwiftUI integration patterns

### Android Integration Steps

1. **Add KMP Dependency**: Include KMP module in your Android app
```kotlin
implementation(project(":kmp:composeApp"))
```

2. **Initialize Commerce**: Set up in your Activity or Application class
```kotlin
val commerceManager = initializeCommerceManager()
```

3. **Create Storage**: Implement Android-specific storage
```kotlin
val storage = AndroidCommerceStorage(context)
```

4. **Connect UI**: Use the provided Compose integration patterns

## Key Integration Patterns

### 1. Dependency Injection

Both platforms show proper dependency injection setup:

**iOS:**
```swift
class CommerceSetup: ObservableObject {
    func initializeCommerceManager() {
        let storage = IOSCommerceStorage()
        let apiClient = CommerceApiClientImpl(...)
        let manager = CommerceManagerImpl(...)
    }
}
```

**Android:**
```kotlin
class CommerceActivity : ComponentActivity() {
    private fun initializeCommerceManager(): CommerceManager {
        val storage = AndroidCommerceStorage(this)
        val apiClient = CommerceApiClientImpl(...)
        return CommerceManagerImpl(...)
    }
}
```

### 2. State Management

**iOS StateFlow Integration:**
```swift
@StateObject private var viewModel = ProductListViewModel()
```

**Android StateFlow Integration:**
```kotlin
val uiState by viewModel.uiState.collectAsState()
```

### 3. Error Handling

Both platforms implement comprehensive error handling:
- Network errors with retry mechanisms
- Loading states with progress indicators
- Empty states with actionable feedback
- Offline mode graceful degradation

### 4. Data Formatting

Platform-specific formatting extensions:

**iOS:**
```swift
extension Product {
    func formattedPrice() -> String {
        return "\(currency) \(String(format: "%.2f", price))"
    }
}
```

**Android:**
```kotlin
fun CartSummary.formattedTotal(): String {
    return "$currency ${String.format("%.2f", total)}"
}
```

## Testing

The examples include test patterns for:

### Unit Testing (Shared)
- Repository testing with mock API clients
- Manager testing with dependency injection
- State synchronization testing
- Offline operation testing

### Integration Testing (Platform-Specific)
- UI component testing
- ViewModel testing
- Storage integration testing
- End-to-end workflow testing

## Performance Considerations

### Memory Management
- Proper lifecycle handling for StateFlow subscriptions
- ViewModels scoped to appropriate lifecycle
- Resource cleanup on app backgrounding

### Network Optimization
- Request caching and deduplication
- Batch operations where possible
- Background sync with proper scheduling

### Storage Optimization
- Efficient serialization/deserialization
- Storage size management
- Proper encryption for sensitive data

## Customization Points

### API Configuration
```kotlin
// Customize API client configuration
val apiClient = CommerceApiClientImpl(
    baseUrl = "https://your-api.com/api/v1",
    timeout = 30000,
    retryCount = 3
)
```

### Storage Configuration
```kotlin
// Customize storage behavior
val storage = AndroidCommerceStorage(
    context = context,
    cacheSize = 10 * 1024 * 1024, // 10MB
    encryptSensitiveData = true
)
```

### Sync Configuration
```kotlin
// Customize synchronization behavior
val syncManager = SyncManagerImpl(
    autoSyncInterval = 5 * 60 * 1000, // 5 minutes
    enableConflictResolution = true,
    retryFailedOperations = true
)
```

## Best Practices

### 1. Initialization
- Initialize commerce early in app lifecycle
- Handle initialization failures gracefully
- Provide loading states during initialization

### 2. State Management
- Use appropriate state scoping (Activity/Fragment for Android, View for iOS)
- Handle state restoration properly
- Implement proper cleanup

### 3. Error Handling
- Provide meaningful error messages
- Implement retry mechanisms
- Handle network connectivity changes

### 4. Performance
- Use appropriate caching strategies
- Implement pagination for large datasets
- Optimize UI updates with proper state diffing

### 5. Security
- Use platform-appropriate secure storage
- Implement proper authentication token handling
- Encrypt sensitive cart and user data

## Common Issues and Solutions

### 1. State Not Updating
**Problem**: UI not reflecting cart changes
**Solution**: Ensure StateFlow collection is properly set up and ViewModels are correctly scoped

### 2. Network Errors
**Problem**: API calls failing
**Solution**: Check network permissions, base URL configuration, and error handling implementation

### 3. Storage Errors
**Problem**: Data not persisting
**Solution**: Verify storage permissions and proper serialization setup

### 4. Platform-Specific Crashes
**Problem**: App crashes on specific platforms
**Solution**: Check platform-specific implementations and dependencies

## Further Reading

- [KMP Official Documentation](https://kotlinlang.org/docs/multiplatform.html)
- [Ktor Client Documentation](https://ktor.io/docs/client.html)
- [StateFlow Documentation](https://developer.android.com/kotlin/flow/stateflow-and-sharedflow)
- [SwiftUI Integration Patterns](https://developer.apple.com/documentation/swiftui)
- [Jetpack Compose Documentation](https://developer.android.com/jetpack/compose)