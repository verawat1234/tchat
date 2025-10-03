# Network Configuration Guide

## Overview

The KMP app uses environment-based network configuration to seamlessly switch between local development and Railway production servers.

## Environment URLs

### Debug Builds (Local Development)
- **Android Emulator**: `http://10.0.2.2:8080/api/v1`
  - Special IP address that routes to host machine's localhost
- **iOS Simulator**: `http://localhost:8080/api/v1`
  - Direct localhost access

### Release Builds (Railway Production)
- **All Platforms**: `https://gateway-service-production-d78d.up.railway.app/api/v1`
  - Public HTTPS endpoint
  - Works on emulators, simulators, and physical devices

## How It Works

### 1. Automatic Environment Detection

```kotlin
// NetworkConfig automatically detects environment
val baseUrl = NetworkConfig.getBaseUrl()

// Debug builds → Local development
// Release builds → Railway production
```

### 2. Build Configuration

```kotlin
// Android Debug Build
buildTypes {
    getByName("debug") {
        buildConfigField("Boolean", "DEBUG", "true")
    }
}

// Android Release Build
buildTypes {
    getByName("release") {
        buildConfigField("Boolean", "DEBUG", "false")
    }
}
```

### 3. Platform-Specific Implementations

#### Android (NetworkConfig.android.kt)
```kotlin
actual fun NetworkConfig.isDebugBuild(): Boolean {
    return BuildConfig.DEBUG
}

actual fun NetworkConfig.getLocalHostUrl(): String {
    return "http://10.0.2.2:8080"  // Android emulator special IP
}
```

#### iOS (NetworkConfig.ios.kt)
```kotlin
actual fun NetworkConfig.isDebugBuild(): Boolean {
    // Detects debug builds from bundle path
    return bundlePath.contains("Debug", ignoreCase = true)
}

actual fun NetworkConfig.getLocalHostUrl(): String {
    return "http://localhost:8080"  // iOS simulator direct localhost
}
```

## Testing Instructions

### Testing with Local Backend

1. **Start Local Backend Services**:
   ```bash
   cd backend/gateway
   go run main.go
   ```

2. **Build Debug Version**:
   ```bash
   cd apps/kmp
   ./gradlew :composeApp:assembleDebug  # Android
   ```

3. **Verify Connection**:
   - App logs should show: `ApiClient initialized with: http://10.0.2.2:8080/api/v1` (Android)
   - App logs should show: `ApiClient initialized with: http://localhost:8080/api/v1` (iOS)

### Testing with Railway Backend

1. **Build Release Version**:
   ```bash
   cd apps/kmp
   ./gradlew :composeApp:assembleRelease  # Android
   ```

2. **Verify Connection**:
   - App logs should show: `ApiClient initialized with: https://gateway-service-production-d78d.up.railway.app/api/v1`

## No ADB Reverse Needed!

❌ **You do NOT need ADB reverse** because:
1. Railway provides public HTTPS endpoint
2. Android emulator uses special IP (`10.0.2.2`) to access host localhost
3. iOS simulator accesses localhost directly
4. Release builds connect directly to Railway via HTTPS

✅ **Direct connection benefits**:
- Works on physical devices without USB
- Works on both Android and iOS
- Simpler setup
- Production-like testing

## Manual Environment Override

For testing purposes, you can manually override the environment:

```kotlin
// Force Railway environment in debug builds
NetworkConfig.setEnvironment(NetworkConfig.Environment.RAILWAY)

// Force local environment in release builds
NetworkConfig.setEnvironment(NetworkConfig.Environment.LOCAL)
```

## API Endpoints

All API clients use centralized configuration:

```kotlin
// Authentication
val authUrl = NetworkConfig.getAuthApiUrl()

// Messaging
val messagingUrl = NetworkConfig.getMessagingApiUrl()

// Video
val videoUrl = NetworkConfig.getVideoApiUrl()

// Social
val socialUrl = NetworkConfig.getSocialApiUrl()

// Commerce
val commerceUrl = NetworkConfig.getCommerceApiUrl()

// Content
val contentUrl = NetworkConfig.getContentApiUrl()

// Payment
val paymentUrl = NetworkConfig.getPaymentApiUrl()

// Notifications
val notificationUrl = NetworkConfig.getNotificationApiUrl()
```

## WebSocket Configuration

Real-time features automatically use correct WebSocket URL:

```kotlin
// Debug → ws://10.0.2.2:8080/ws (Android) or ws://localhost:8080/ws (iOS)
// Release → wss://gateway-service-production-d78d.up.railway.app/ws
val wsUrl = NetworkConfig.getWebSocketUrl()
```

## Debugging

### Check Current Configuration

```kotlin
// Print configuration info
println(NetworkConfig.getConfigInfo())

// Output example:
// Network Configuration:
// - Environment: RAILWAY
// - Base URL: https://gateway-service-production-d78d.up.railway.app/api/v1
// - WebSocket: wss://gateway-service-production-d78d.up.railway.app/ws
// - Is Debug: false
```

### Check Environment State

```kotlin
if (NetworkConfig.isRailwayEnvironment()) {
    println("Using Railway production server")
}

if (NetworkConfig.isLocalEnvironment()) {
    println("Using local development server")
}

println("Environment: ${NetworkConfig.getEnvironmentName()}")
```

## Troubleshooting

### Android Emulator Can't Connect to Localhost

**Problem**: Connection refused when trying to access localhost

**Solution**:
- Ensure you're using `10.0.2.2` instead of `localhost`
- NetworkConfig handles this automatically in debug builds

### iOS Simulator Can't Connect

**Problem**: Connection timeout or refused

**Solution**:
- Ensure local backend is running on port 8080
- Check firewall settings allow local connections
- Verify no VPN blocking localhost access

### Railway Connection Issues

**Problem**: 502 Bad Gateway or connection timeout

**Solution**:
- Verify Railway services are running: https://railway.app/project/0a1f3508-2150-4d0c-8ae9-878f74a607a0
- Check Railway service logs for errors
- Ensure all required services (gateway, auth, messaging) are deployed

## Related Files

- `NetworkConfig.kt` - Common configuration interface
- `NetworkConfig.android.kt` - Android platform implementation
- `NetworkConfig.ios.kt` - iOS platform implementation
- `ApiClient.kt` - Uses NetworkConfig for base URL
- `VideoRepository.kt` - Uses NetworkConfig for video API
- `SocialApiClient.kt` - Uses NetworkConfig for social API
- `build.gradle.kts` - Build configuration with DEBUG flags
