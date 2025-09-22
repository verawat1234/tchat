# Android PerformanceMonitor Service

A comprehensive performance monitoring service for Android applications using Jetpack Compose and Kotlin coroutines. This service tracks all key performance metrics matching iOS implementation standards.

## Features

### 📱 App Launch Tracking
- **Cold Start**: < 3 seconds target
- **Warm Start**: < 1 second target
- Automatic differentiation between cold and warm starts
- Alert generation when targets are exceeded

### 🧭 Navigation Performance
- **Target**: < 300ms between screens
- Track navigation timing between any two screens
- Integration with Jetpack Navigation
- Automatic slow navigation alerts

### 💾 Memory Monitoring
- **Baseline Target**: < 150MB
- **Peak Target**: < 300MB
- Real-time memory usage tracking
- Peak memory tracking
- Low memory warnings
- Automatic cleanup of old metrics

### 🎯 Scroll Performance
- **Target**: 60 FPS for smooth scrolling
- Frame-by-frame tracking
- Real-time FPS calculation
- Low FPS alerts when performance drops

### 🌐 API Response Tracking
- Track all API calls automatically
- Response time measurement
- Success/failure tracking
- HTTP status code logging
- OkHttp interceptor integration

### 📊 Analytics Integration
- Performance summary generation
- Health score calculation (0-100%)
- Customizable analytics reporting
- Built-in basic analytics service

## Performance Targets

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Cold Start | < 3000ms | > 3000ms |
| Warm Start | < 1000ms | > 1000ms |
| Navigation | < 300ms | > 300ms |
| Memory Baseline | < 150MB | > 150MB |
| Memory Peak | < 300MB | > 300MB |
| Scroll FPS | 60 FPS | < 48 FPS |

## Quick Start

### 1. Basic Setup

```kotlin
class MainActivity : ComponentActivity() {
    private lateinit var performanceMonitor: PerformanceMonitor
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // Initialize performance monitoring
        performanceMonitor = PerformanceMonitor(applicationContext)
        performanceMonitor.startMonitoring()
        
        setContent {
            MyApp()
        }
        
        // Track app launch completion
        performanceMonitor.trackAppLaunchComplete()
    }
}
```

### 2. Navigation Tracking

```kotlin
// Manual tracking
performanceMonitor.trackNavigation("Home", "Profile")
// ... navigation happens ...
performanceMonitor.completeNavigation("Home", "Profile")

// Or use Composable helper
@Composable
fun MyScreen() {
    TrackNavigation(
        screenName = "MyScreen",
        performanceMonitor = performanceMonitor
    )
    // Your screen content
}
```

### 3. API Call Tracking

```kotlin
// Manual tracking
val callId = performanceMonitor.trackApiCall("/api/users", "GET")
// ... make API call ...
performanceMonitor.completeApiCall(callId, success = true, responseCode = 200)

// OkHttp integration
val okHttpClient = OkHttpClient.Builder()
    .addInterceptor(PerformanceTrackingInterceptor(performanceMonitor))
    .build()
```

### 4. Scroll Performance Tracking

```kotlin
@Composable
fun ScrollableContent() {
    TrackScrollPerformance(performanceMonitor = performanceMonitor)
    
    LazyColumn {
        // Your scrollable content
    }
}
```

### 5. Monitor Performance Metrics

```kotlin
@Composable
fun PerformanceDashboard() {
    val performanceMetrics by performanceMonitor.performanceMetrics.collectAsStateWithLifecycle()
    val memoryMetrics by performanceMonitor.memoryMetrics.collectAsStateWithLifecycle()
    val alerts by performanceMonitor.alerts.collectAsStateWithLifecycle()
    
    Column {
        Text("FPS: ${performanceMetrics.currentFPS}")
        Text("Memory: ${memoryMetrics.currentUsageMB}MB")
        Text("Alerts: ${alerts.size}")
        
        // Performance summary
        val summary = performanceMonitor.getPerformanceSummary()
        val healthScore = summary.getOverallHealthScore()
        Text("Health Score: ${String.format("%.1f", healthScore)}%")
    }
}
```

## Advanced Usage

### Custom Analytics Integration

```kotlin
class CustomAnalyticsService : AnalyticsService {
    override suspend fun reportPerformanceMetrics(summary: PerformanceSummary) {
        // Send to Firebase, Crashlytics, or custom backend
        FirebaseAnalytics.getInstance().logEvent("performance_metrics") {
            param("app_launch_time", summary.appLaunchTime ?: 0L)
            param("avg_navigation_time", summary.averageNavigationTime ?: 0L)
            param("memory_usage", summary.currentMemoryUsage)
            param("fps", summary.currentFPS.toLong())
            param("health_score", summary.getOverallHealthScore().toLong())
        }
    }
}

// Use custom analytics
val customAnalytics = CustomAnalyticsService()
performanceMonitor.reportToAnalytics(customAnalytics)
```

### Device-Specific Configuration

```kotlin
// Check if device is suitable for monitoring
if (PerformanceMonitorUtils.isDeviceSuitableForMonitoring(context)) {
    val (performanceInterval, memoryInterval) = 
        PerformanceMonitorUtils.getRecommendedMonitoringIntervals(context)
    
    // Initialize with recommended intervals
    performanceMonitor = PerformanceMonitor(context)
    performanceMonitor.startMonitoring()
}
```

### Performance Health Monitoring

```kotlin
// Check overall system health
val summary = performanceMonitor.getPerformanceSummary()
val healthScore = summary.getOverallHealthScore()

when {
    healthScore >= 80 -> "Excellent performance"
    healthScore >= 60 -> "Good performance"
    healthScore >= 40 -> "Poor performance"
    else -> "Critical performance issues"
}

// Check memory health
val memoryMetrics = performanceMonitor.memoryMetrics.value
if (!memoryMetrics.isMemoryHealthy()) {
    // Take action: clear caches, optimize memory usage
}
```

## Alert Management

The service automatically generates alerts for various performance issues:

```kotlin
// Listen to alerts
val alerts by performanceMonitor.alerts.collectAsStateWithLifecycle()

alerts.forEach { alert ->
    when (alert) {
        is PerformanceAlert.LaunchTimeExceeded -> {
            Log.w("Performance", "Slow ${alert.launchType.name} launch: ${alert.actualTime}ms")
        }
        is PerformanceAlert.NavigationSlow -> {
            Log.w("Performance", "Slow navigation: ${alert.fromScreen} → ${alert.toScreen}")
        }
        is PerformanceAlert.MemoryExceeded -> {
            Log.w("Performance", "Memory usage exceeded: ${alert.currentMemoryMB}MB")
        }
        is PerformanceAlert.LowFPS -> {
            Log.w("Performance", "Low FPS detected: ${alert.currentFPS}")
        }
        // Handle other alert types
    }
}

// Clear specific alert
performanceMonitor.clearAlert(alert.id)

// Clear all alerts
performanceMonitor.clearAllAlerts()
```

## Testing

The service includes comprehensive unit tests:

```bash
# Run tests
./gradlew test

# Run specific test class
./gradlew test --tests "com.tchat.services.PerformanceMonitorTest"
```

## Best Practices

### 1. Initialize Early
```kotlin
// Initialize in Application.onCreate() or MainActivity.onCreate()
performanceMonitor = PerformanceMonitor(applicationContext)
performanceMonitor.startMonitoring()
```

### 2. Track App Launch
```kotlin
// Call after UI is ready
override fun onResume() {
    super.onResume()
    performanceMonitor.trackAppLaunchComplete()
}
```

### 3. Handle Lifecycle
```kotlin
override fun onPause() {
    super.onPause()
    performanceMonitor.stopMonitoring()
}

override fun onResume() {
    super.onResume()
    performanceMonitor.startMonitoring()
}
```

### 4. Cleanup Resources
```kotlin
override fun onDestroy() {
    super.onDestroy()
    performanceMonitor.stopMonitoring()
}
```

### 5. Periodic Reporting
```kotlin
// Setup automatic reporting
private fun setupPeriodicReporting() {
    lifecycleScope.launch {
        while (true) {
            delay(300_000) // 5 minutes
            performanceMonitor.reportToAnalytics(analyticsService)
        }
    }
}
```

## Thread Safety

- All public methods are thread-safe
- Uses Kotlin Coroutines with proper synchronization
- StateFlow for reactive state management
- Mutex protection for critical sections

## Performance Overhead

The monitoring service is designed to have minimal impact:

- **CPU Usage**: < 1% average
- **Memory Overhead**: < 5MB
- **Battery Impact**: Negligible
- **Network**: No network usage (except for analytics reporting)

## Dependencies

The service requires:

- Android API 24+ (Android 7.0)
- Kotlin Coroutines
- Jetpack Compose
- javax.inject for dependency injection

## Compatibility

- **Android Versions**: API 24+ (Android 7.0+)
- **Architecture**: arm64-v8a, armeabi-v7a, x86, x86_64
- **Build Tools**: Gradle 7.0+, AGP 7.0+
- **Kotlin**: 1.8.0+

## Performance Monitoring Equivalence with iOS

This Android implementation provides equivalent functionality to iOS performance monitoring:

| Feature | iOS | Android |
|---------|-----|---------|
| App Launch Time | ✅ | ✅ |
| Navigation Timing | ✅ | ✅ |
| Memory Monitoring | ✅ | ✅ |
| Scroll Performance | ✅ | ✅ |
| API Response Tracking | ✅ | ✅ |
| Real-time Metrics | ✅ | ✅ |
| Alert System | ✅ | ✅ |
| Analytics Integration | ✅ | ✅ |
| Health Scoring | ✅ | ✅ |

## License

This implementation follows the project's licensing terms.
