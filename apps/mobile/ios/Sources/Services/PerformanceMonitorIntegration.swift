//
//  PerformanceMonitorIntegration.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI
import Combine

/// Integration service for PerformanceMonitor with TchatApp
@MainActor
public class PerformanceMonitorIntegration: ObservableObject {
    
    // MARK: - Singleton
    public static let shared = PerformanceMonitorIntegration()
    
    // MARK: - Properties
    public let performanceMonitor = PerformanceMonitor.shared
    private var cancellables = Set<AnyCancellable>()
    
    // MARK: - Analytics Integration
    private var analyticsService: AnalyticsService?
    
    private init() {
        setupAnalyticsIntegration()
        setupPerformanceMonitoring()
    }
    
    // MARK: - Setup
    
    private func setupAnalyticsIntegration() {
        // Set up analytics delegate for performance reporting
        performanceMonitor.analyticsDelegate = self
    }
    
    private func setupPerformanceMonitoring() {
        // Start monitoring when app becomes active
        NotificationCenter.default.publisher(for: UIApplication.didBecomeActiveNotification)
            .sink { [weak self] _ in
                self?.performanceMonitor.startMonitoring()
            }
            .store(in: &cancellables)
        
        // Stop monitoring when app goes to background
        NotificationCenter.default.publisher(for: UIApplication.didEnterBackgroundNotification)
            .sink { [weak self] _ in
                self?.performanceMonitor.stopMonitoring()
            }
            .store(in: &cancellables)
    }
    
    // MARK: - App Launch Integration
    
    /// Call this from App delegate or equivalent
    public func trackAppLaunch() {
        performanceMonitor.trackAppLaunchStart()
    }
    
    // MARK: - Navigation Integration
    
    /// Track navigation performance between screens
    public func trackNavigation(from fromScreen: String, to toScreen: String) {
        performanceMonitor.trackNavigationStart(from: fromScreen, to: toScreen)
        
        // Auto-track completion after a delay (or call trackNavigationComplete manually)
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) { [weak self] in
            self?.performanceMonitor.trackNavigationEnd(from: fromScreen, to: toScreen)
        }
    }
    
    /// Manually complete navigation tracking
    public func trackNavigationComplete(from fromScreen: String, to toScreen: String) {
        performanceMonitor.trackNavigationEnd(from: fromScreen, to: toScreen)
    }
    
    // MARK: - API Integration
    
    /// Integrate with existing API clients for performance tracking
    public func trackAPIRequest(endpoint: String, method: String, startTime: CFAbsoluteTime, statusCode: Int) {
        performanceMonitor.trackAPIRequest(
            endpoint: endpoint,
            method: method,
            startTime: startTime,
            statusCode: statusCode
        )
    }
    
    // MARK: - SwiftUI Integration Helpers
    
    /// Create performance-aware navigation modifier
    public func navigationPerformanceModifier(screenName: String) -> some ViewModifier {
        PerformanceNavigationModifier(screenName: screenName, integration: self)
    }
    
    /// Create scroll performance monitoring modifier
    public func scrollPerformanceModifier(screenName: String) -> some ViewModifier {
        return performanceMonitor.scrollMonitor(screenName: screenName)
    }
}

// MARK: - Analytics Integration
// TODO: Implement PerformanceAnalyticsDelegate protocol when types are defined

/*
extension PerformanceMonitorIntegration: PerformanceAnalyticsDelegate {
    
    public func reportLaunchMetrics(_ metrics: PerformanceMonitor.LaunchMetrics) {
        // Report to your analytics service
        let eventData: [String: Any] = [
            "launch_type": "\(metrics.launchType)",
            "cold_start_time": metrics.coldStartTime,
            "warm_start_time": metrics.warmStartTime,
            "timestamp": metrics.timestamp.timeIntervalSince1970
        ]
        
        analyticsService?.track(event: "app_launch_performance", properties: eventData)
        
        // Log to console for debugging
        print("📊 Launch Performance: \(metrics.launchType) - \(metrics.launchType == .cold ? metrics.coldStartTime : metrics.warmStartTime)s")
    }
    
    public func reportNavigationMetrics(_ metrics: PerformanceMonitor.NavigationMetrics) {
        let eventData: [String: Any] = [
            "from_screen": metrics.fromScreen,
            "to_screen": metrics.toScreen,
            "transition_time": metrics.transitionTime,
            "timestamp": metrics.timestamp.timeIntervalSince1970
        ]
        
        analyticsService?.track(event: "navigation_performance", properties: eventData)
        
        print("📊 Navigation Performance: \(metrics.fromScreen) → \(metrics.toScreen) - \(String(format: "%.3f", metrics.transitionTime))s")
    }
    
    public func reportMemoryMetrics(_ metrics: PerformanceMonitor.MemoryMetrics) {
        let eventData: [String: Any] = [
            "current_usage_mb": metrics.currentUsageMB,
            "peak_usage_mb": metrics.peakUsageMB,
            "available_memory_mb": Double(metrics.availableMemory) / (1024 * 1024),
            "timestamp": metrics.timestamp.timeIntervalSince1970
        ]
        
        analyticsService?.track(event: "memory_performance", properties: eventData)
        
        // Only log warnings/errors to avoid spam
        if metrics.currentUsageMB > PerformanceMonitor.PerformanceThresholds.baselineMemoryMB {
            print("📊 Memory Usage: \(String(format: "%.1f", metrics.currentUsageMB))MB (baseline: \(PerformanceMonitor.PerformanceThresholds.baselineMemoryMB)MB)")
        }
    }
    
    public func reportScrollMetrics(_ metrics: PerformanceMonitor.ScrollMetrics) {
        let eventData: [String: Any] = [
            "screen_name": metrics.screenName,
            "average_fps": metrics.averageFPS,
            "dropped_frames": metrics.droppedFrames,
            "scroll_duration": metrics.scrollDuration,
            "timestamp": metrics.timestamp.timeIntervalSince1970
        ]
        
        analyticsService?.track(event: "scroll_performance", properties: eventData)
        
        if metrics.averageFPS < PerformanceMonitor.PerformanceThresholds.targetFPS * 0.9 {
            print("📊 Scroll Performance: \(metrics.screenName) - \(String(format: "%.1f", metrics.averageFPS)) FPS")
        }
    }
    
    public func reportAPIMetrics(_ metrics: PerformanceMonitor.APIMetrics) {
        let eventData: [String: Any] = [
            "endpoint": metrics.endpoint,
            "method": metrics.method,
            "response_time": metrics.responseTime,
            "status_code": metrics.statusCode,
            "success": metrics.success,
            "timestamp": metrics.timestamp.timeIntervalSince1970
        ]
        
        analyticsService?.track(event: "api_performance", properties: eventData)
        
        if metrics.responseTime > PerformanceMonitor.PerformanceThresholds.apiResponseTarget {
            print("📊 API Performance: \(metrics.method) \(metrics.endpoint) - \(String(format: "%.3f", metrics.responseTime))s")
        }
    }
}
*/

// MARK: - SwiftUI Modifiers

private struct PerformanceNavigationModifier: ViewModifier {
    let screenName: String
    let integration: PerformanceMonitorIntegration
    @State private var previousScreen: String?
    
    func body(content: Content) -> some View {
        content
            .onAppear {
                if let previous = previousScreen {
                    integration.trackNavigationComplete(from: previous, to: screenName)
                }
                previousScreen = screenName
            }
    }
}

// MARK: - Analytics Service Protocol

public protocol AnalyticsService {
    func track(event: String, properties: [String: Any])
}

// MARK: - Usage Examples and Integration Guides

/*
 
 USAGE EXAMPLES:
 
 1. App Launch Integration:
 
 @main
 struct TchatApp: App {
     init() {
         PerformanceMonitorIntegration.shared.trackAppLaunch()
     }
     
     var body: some Scene {
         WindowGroup {
             ContentView()
                 .environmentObject(PerformanceMonitorIntegration.shared)
         }
     }
 }
 
 2. Navigation Integration:
 
 struct TabNavigationView: View {
     @EnvironmentObject private var performanceIntegration: PerformanceMonitorIntegration
     
     var body: some View {
         TabView {
             ChatListView()
                 .modifier(performanceIntegration.navigationPerformanceModifier(screenName: "ChatList"))
                 .tabItem { Label("Chats", systemImage: "message") }
             
             ProfileView()
                 .modifier(performanceIntegration.navigationPerformanceModifier(screenName: "Profile"))
                 .tabItem { Label("Profile", systemImage: "person") }
         }
     }
 }
 
 3. Scroll Performance Integration:
 
 struct ChatListView: View {
     @EnvironmentObject private var performanceIntegration: PerformanceMonitorIntegration
     
     var body: some View {
         ScrollView {
             LazyVStack {
                 ForEach(chats) { chat in
                     ChatRowView(chat: chat)
                 }
             }
         }
         .modifier(performanceIntegration.scrollPerformanceModifier(screenName: "ChatList"))
     }
 }
 
 4. API Integration with existing APIClient:
 
 extension PlatformAdapterAPIClient {
     private func performRequest<T>(_ request: URLRequest) async throws -> T where T: Decodable {
         let startTime = CFAbsoluteTimeGetCurrent()
         
         do {
             let (data, response) = try await URLSession.shared.data(for: request)
             let httpResponse = response as! HTTPURLResponse
             
             // Track performance
             PerformanceMonitorIntegration.shared.trackAPIRequest(
                 endpoint: request.url?.path ?? "",
                 method: request.httpMethod ?? "GET",
                 startTime: startTime,
                 statusCode: httpResponse.statusCode
             )
             
             return try JSONDecoder().decode(T.self, from: data)
         } catch {
             // Track failed request
             PerformanceMonitorIntegration.shared.trackAPIRequest(
                 endpoint: request.url?.path ?? "",
                 method: request.httpMethod ?? "GET", 
                 startTime: startTime,
                 statusCode: 0
             )
             throw error
         }
     }
 }
 
 5. Performance Dashboard View:
 
 struct PerformanceDashboardView: View {
     @ObservedObject private var monitor = PerformanceMonitorIntegration.shared.performanceMonitor
     
     var body: some View {
         List {
             Section("Memory") {
                 if let memory = monitor.currentMemoryUsage {
                     Text("Current: \(String(format: "%.1f", memory.currentUsageMB))MB")
                     Text("Peak: \(String(format: "%.1f", memory.peakUsageMB))MB")
                 }
             }
             
             Section("Recent Alerts") {
                 ForEach(monitor.performanceAlerts, id: \.id) { alert in
                     VStack(alignment: .leading) {
                         Text(alert.message)
                         Text(alert.timestamp, style: .time)
                             .font(.caption)
                             .foregroundColor(.secondary)
                     }
                 }
             }
         }
         .navigationTitle("Performance")
     }
 }

*/
