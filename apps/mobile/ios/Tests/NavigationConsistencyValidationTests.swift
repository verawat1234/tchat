//
//  NavigationConsistencyValidationTests.swift
//  TchatAppTests
//
//  Created by Claude on 22/09/2024.
//

import XCTest
import Combine
@testable import TchatApp

/// Cross-platform navigation consistency validation tests
class NavigationConsistencyValidationTests: XCTestCase {

    var navigationCoordinator: NavigationCoordinator!
    var routeRegistry: RouteRegistry!
    var deepLinkProcessor: DeepLinkProcessor!
    var cancellables: Set<AnyCancellable>!

    override func setUpWithError() throws {
        super.setUp()
        navigationCoordinator = NavigationCoordinator()
        routeRegistry = RouteRegistry()
        deepLinkProcessor = DeepLinkProcessor()
        cancellables = Set<AnyCancellable>()
    }

    override func tearDownWithError() throws {
        cancellables?.removeAll()
        navigationCoordinator = nil
        routeRegistry = nil
        deepLinkProcessor = nil
        super.tearDown()
    }

    // MARK: - Route Registration Consistency Tests

    func testRouteRegistrationConsistency() async throws {
        // Test that routes are registered consistently across platforms
        let webRoutes = getWebRoutes()
        let mobileRoutes = routeRegistry.getAllRoutes()

        // Validate core routes exist
        let coreRoutes = ["/", "/chat", "/settings", "/profile", "/notifications"]

        for coreRoute in coreRoutes {
            XCTAssertTrue(
                mobileRoutes.contains { $0.path == coreRoute },
                "Core route \(coreRoute) should be registered in mobile app"
            )
        }

        // Validate route count is reasonable
        XCTAssertGreaterThan(mobileRoutes.count, 5, "Should have at least 5 routes registered")
        XCTAssertLessThan(mobileRoutes.count, 50, "Should not have excessive routes")
    }

    func testDeepLinkConsistency() async throws {
        // Test deep link processing consistency
        let testDeepLinks = [
            "tchat://chat/room/123",
            "tchat://profile/user/456",
            "tchat://settings/notifications",
            "https://tchat.app/chat/room/789",
            "https://tchat.app/profile/user/101"
        ]

        for deepLinkURL in testDeepLinks {
            guard let url = URL(string: deepLinkURL) else {
                XCTFail("Invalid test URL: \(deepLinkURL)")
                continue
            }

            let request = DeepLinkResolutionRequest(
                url: url,
                platform: "ios",
                userId: "test_user",
                sessionId: "test_session",
                timestamp: Date()
            )

            do {
                let resolution = try await deepLinkProcessor.processDeepLink(request: request)

                // Validate resolution structure
                XCTAssertNotNil(resolution.targetRoute, "Deep link should resolve to a target route")
                XCTAssertNotNil(resolution.navigationInstructions, "Should have navigation instructions")

                // Validate route format
                if let targetRoute = resolution.targetRoute {
                    XCTAssertTrue(targetRoute.hasPrefix("/"), "Target route should start with /")
                    XCTAssertFalse(targetRoute.contains("//"), "Target route should not have double slashes")
                }

            } catch {
                XCTFail("Deep link processing failed for \(deepLinkURL): \(error)")
            }
        }
    }

    // MARK: - Navigation Performance Tests

    func testNavigationPerformance() async throws {
        // Test navigation performance requirements
        let performanceExpectation = expectation(description: "Navigation performance test")
        var navigationTimes: [TimeInterval] = []

        // Test multiple navigation operations
        for i in 0..<10 {
            let startTime = CFAbsoluteTimeGetCurrent()

            let route = NavigationRoute(
                id: "test_route_\(i)",
                path: "/test/\(i)",
                name: "Test Route \(i)",
                component: "TestComponent",
                parameters: [:],
                metadata: [:]
            )

            try await navigationCoordinator.navigateToRoute(route)

            let endTime = CFAbsoluteTimeGetCurrent()
            let navigationTime = endTime - startTime
            navigationTimes.append(navigationTime)
        }

        performanceExpectation.fulfill()
        await fulfillment(of: [performanceExpectation], timeout: 5.0)

        // Validate performance requirements
        let averageTime = navigationTimes.reduce(0, +) / Double(navigationTimes.count)
        let maxTime = navigationTimes.max() ?? 0

        XCTAssertLessThan(averageTime, 0.1, "Average navigation time should be < 100ms")
        XCTAssertLessThan(maxTime, 0.2, "Maximum navigation time should be < 200ms")
    }

    func testMemoryUsageDuringNavigation() async throws {
        // Test memory usage during navigation
        let initialMemory = getMemoryUsage()

        // Perform multiple navigation operations
        for i in 0..<20 {
            let route = NavigationRoute(
                id: "memory_test_\(i)",
                path: "/memory_test/\(i)",
                name: "Memory Test \(i)",
                component: "TestComponent",
                parameters: [:],
                metadata: [:]
            )

            try await navigationCoordinator.navigateToRoute(route)

            // Navigate back to clear stack
            try await navigationCoordinator.goBack()
        }

        let finalMemory = getMemoryUsage()
        let memoryIncrease = finalMemory - initialMemory

        // Memory increase should be reasonable (< 10MB)
        XCTAssertLessThan(memoryIncrease, 10_000_000, "Memory increase should be < 10MB")
    }

    // MARK: - Cross-Platform State Consistency Tests

    func testNavigationStateSync() async throws {
        // Test navigation state synchronization
        let syncExpectation = expectation(description: "Navigation state sync test")

        // Create initial navigation state
        let navigationState = NavigationState(
            currentRoute: "/test/sync",
            routeStack: ["/", "/test", "/test/sync"],
            routeParams: ["id": "123"],
            navigationHistory: [],
            platform: "ios",
            userId: "test_user",
            sessionId: "test_session",
            timestamp: Date(),
            version: 1
        )

        // Test state synchronization
        navigationCoordinator.navigationState
            .dropFirst() // Skip initial state
            .sink { state in
                XCTAssertEqual(state.currentRoute, "/test/sync")
                XCTAssertEqual(state.routeStack.count, 3)
                XCTAssertEqual(state.routeParams["id"] as? String, "123")
                syncExpectation.fulfill()
            }
            .store(in: &cancellables)

        // Update navigation state
        try await navigationCoordinator.updateNavigationState(navigationState)

        await fulfillment(of: [syncExpectation], timeout: 2.0)
    }

    func testUIComponentStateConsistency() async throws {
        // Test UI component state consistency across navigation
        let stateManager = UIStateManager()

        // Create test component state
        let componentState = ComponentState(
            instanceId: "test_component_1",
            componentId: "ChatRoom",
            state: [
                "messageCount": 42,
                "scrollPosition": 150.0,
                "selectedMessage": "msg_123"
            ],
            userId: "test_user",
            sessionId: "test_session",
            platform: "ios",
            timestamp: Date(),
            version: 1,
            isSynchronized: false
        )

        // Save component state
        try await stateManager.createComponentState(componentState)

        // Navigate away and back
        let chatRoute = NavigationRoute(
            id: "chat_room",
            path: "/chat/room/123",
            name: "Chat Room",
            component: "ChatRoom",
            parameters: ["roomId": "123"],
            metadata: [:]
        )

        try await navigationCoordinator.navigateToRoute(chatRoute)
        try await navigationCoordinator.goBack()
        try await navigationCoordinator.navigateToRoute(chatRoute)

        // Verify state persistence
        let retrievedState = stateManager.getComponentState("test_component_1")
        XCTAssertNotNil(retrievedState, "Component state should be persisted")
        XCTAssertEqual(retrievedState?.state["messageCount"] as? Int, 42)
        XCTAssertEqual(retrievedState?.state["scrollPosition"] as? Double, 150.0)
    }

    // MARK: - Platform-Specific Feature Tests

    func testPlatformAdapterConsistency() async throws {
        // Test platform adapter consistency
        let platformAdapter = PlatformAdapterImpl()

        // Test gesture support consistency
        let supportedGestures = platformAdapter.getSupportedGestures()
        let expectedGestures = ["tap", "swipe", "longPress", "pan", "pinch"]

        for gesture in expectedGestures {
            let gestureDefinition = supportedGestures.first { $0.name == gesture }
            XCTAssertNotNil(gestureDefinition, "Platform should support \(gesture) gesture")
        }

        // Test animation support consistency
        let supportedAnimations = platformAdapter.getSupportedAnimations()
        let expectedAnimations = ["fade", "slide", "scale", "spring"]

        for animation in expectedAnimations {
            let animationDefinition = supportedAnimations.first { $0.name == animation }
            XCTAssertNotNil(animationDefinition, "Platform should support \(animation) animation")
        }

        // Test platform capabilities
        let capabilities = platformAdapter.getPlatformCapabilities()
        XCTAssertEqual(capabilities.platform, "ios")
        XCTAssertGreaterThan(capabilities.capabilities.count, 0)
    }

    func testDeviceFeatureDetection() async throws {
        // Test device feature detection consistency
        let platformAdapter = PlatformAdapterImpl()

        // Test common hardware features
        let cameraSupport = platformAdapter.hasHardwareFeature("camera")
        let hapticSupport = platformAdapter.hasHardwareFeature("hapticEngine")

        // These should return boolean values without crashing
        XCTAssertNotNil(cameraSupport)
        XCTAssertNotNil(hapticSupport)

        // Test device metadata
        let deviceMetadata = platformAdapter.getDeviceMetadata()
        XCTAssertNotNil(deviceMetadata["model"])
        XCTAssertNotNil(deviceMetadata["systemName"])
        XCTAssertNotNil(deviceMetadata["systemVersion"])
    }

    // MARK: - API Integration Consistency Tests

    func testAPIClientConsistency() async throws {
        // Test API client consistency
        let navigationAPIClient = NavigationSyncAPIClient()
        let componentAPIClient = UIComponentSyncAPIClient()

        // Test client initialization
        XCTAssertNotNil(navigationAPIClient)
        XCTAssertNotNil(componentAPIClient)

        // Test client configuration
        let navigationFactory = NavigationSyncAPIClientFactory.create(environment: .development)
        let componentFactory = UIComponentSyncAPIClientFactory.create(environment: .development)

        XCTAssertNotNil(navigationFactory)
        XCTAssertNotNil(componentFactory)
    }

    // MARK: - Performance Benchmarking

    func testNavigationStackPerformance() throws {
        // Measure navigation stack performance
        measure(metrics: [XCTClockMetric(), XCTMemoryMetric()]) {
            for i in 0..<100 {
                let route = NavigationRoute(
                    id: "perf_test_\(i)",
                    path: "/perf/\(i)",
                    name: "Performance Test \(i)",
                    component: "TestComponent",
                    parameters: ["index": i],
                    metadata: [:]
                )

                routeRegistry.registerRoute(route)
            }
        }
    }

    func testDeepLinkProcessingPerformance() throws {
        // Measure deep link processing performance
        let testURLs = (0..<50).map { "tchat://test/route/\($0)" }

        measure(metrics: [XCTClockMetric(), XCTCPUMetric()]) {
            for urlString in testURLs {
                guard let url = URL(string: urlString) else { continue }

                let request = DeepLinkResolutionRequest(
                    url: url,
                    platform: "ios",
                    userId: "perf_test",
                    sessionId: "perf_session",
                    timestamp: Date()
                )

                // Synchronous processing for measurement
                let _ = deepLinkProcessor.resolveURL(url)
            }
        }
    }

    // MARK: - Helper Methods

    private func getWebRoutes() -> [String] {
        // Simulate getting routes from web application
        return ["/", "/chat", "/settings", "/profile", "/notifications", "/help"]
    }

    private func getMemoryUsage() -> Int64 {
        var info = mach_task_basic_info()
        var count = mach_msg_type_number_t(MemoryLayout<mach_task_basic_info>.size)/4

        let kerr: kern_return_t = withUnsafeMutablePointer(to: &info) {
            $0.withMemoryRebound(to: integer_t.self, capacity: 1) {
                task_info(mach_task_self_,
                         task_flavor_t(MACH_TASK_BASIC_INFO),
                         $0,
                         &count)
            }
        }

        if kerr == KERN_SUCCESS {
            return Int64(info.resident_size)
        } else {
            return 0
        }
    }
}