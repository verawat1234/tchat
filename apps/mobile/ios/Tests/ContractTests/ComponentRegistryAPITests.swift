//
//  ComponentRegistryAPITests.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import XCTest
@testable import TchatApp

/// Contract tests for Component Registry API
/// These tests verify the API contract for component metadata and configuration
class ComponentRegistryAPITests: TchatAppTestCase {

    private var apiClient: ComponentRegistryAPIClient!

    override func setUp() {
        super.setUp()
        apiClient = ComponentRegistryAPIClient(baseURL: TestConfiguration.shared.apiBaseURL)
    }

    override func tearDown() {
        apiClient = nil
        super.tearDown()
    }

    // MARK: - GET /api/components Tests

    func testGetComponents_ValidPlatform_ReturnsComponents() {
        // Given
        let expectation = expectation(description: "Get components")
        let platform = "ios"

        // When
        apiClient.getComponents(platform: platform, category: nil) { result in
            // Then
            switch result {
            case .success(let components):
                XCTAssertGreaterThan(components.count, 0)
                for component in components {
                    XCTAssertEqual(component.platform, platform)
                    XCTAssertFalse(component.id.isEmpty)
                    XCTAssertFalse(component.name.isEmpty)
                    XCTAssertNotNil(component.category)
                    XCTAssertNotNil(component.version)
                    XCTAssertNotNil(component.webEquivalent)
                    XCTAssertNotNil(component.accessibility)
                }
            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetComponents_WithCategory_ReturnsFilteredComponents() {
        // Given
        let expectation = expectation(description: "Get components by category")
        let platform = "ios"
        let category = "navigation"

        // When
        apiClient.getComponents(platform: platform, category: category) { result in
            // Then
            switch result {
            case .success(let components):
                XCTAssertGreaterThan(components.count, 0)
                for component in components {
                    XCTAssertEqual(component.category.rawValue, category)
                    XCTAssertEqual(component.platform, platform)
                }
            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetComponents_InvalidPlatform_ReturnsError() {
        // Given
        let expectation = expectation(description: "Get components with invalid platform")
        let platform = "invalid"

        // When
        apiClient.getComponents(platform: platform, category: nil) { result in
            // Then
            switch result {
            case .success:
                XCTFail("Expected error for invalid platform")
            case .failure(let error):
                XCTAssertEqual(error.statusCode, 400)
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    // MARK: - GET /api/components/{componentId}/config Tests

    func testGetComponentConfig_ValidComponent_ReturnsConfig() {
        // Given
        let expectation = expectation(description: "Get component config")
        let componentId = "tchat-button"
        let platform = "ios"

        // When
        apiClient.getComponentConfig(componentId: componentId, platform: platform) { result in
            // Then
            switch result {
            case .success(let config):
                XCTAssertEqual(config.id, componentId)
                XCTAssertGreaterThan(config.variants.count, 0)
                XCTAssertGreaterThan(config.props.count, 0)
                XCTAssertGreaterThan(config.animations.count, 0)
                XCTAssertNotNil(config.accessibility)

                // Validate variant structure
                for variant in config.variants {
                    XCTAssertFalse(variant.name.isEmpty)
                    XCTAssertFalse(variant.description.isEmpty)
                    XCTAssertNotNil(variant.styles)
                    XCTAssertNotNil(variant.example)
                }

                // Validate props structure
                for prop in config.props {
                    XCTAssertFalse(prop.name.isEmpty)
                    XCTAssertNotNil(prop.type)
                    XCTAssertNotNil(prop.required)
                    XCTAssertNotNil(prop.description)
                }

            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetComponentConfig_InvalidComponent_ReturnsNotFound() {
        // Given
        let expectation = expectation(description: "Get config for invalid component")
        let componentId = "non-existent-component"
        let platform = "ios"

        // When
        apiClient.getComponentConfig(componentId: componentId, platform: platform) { result in
            // Then
            switch result {
            case .success:
                XCTFail("Expected error for invalid component")
            case .failure(let error):
                XCTAssertEqual(error.statusCode, 404)
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetComponentConfig_ValidatesAccessibilityFeatures() {
        // Given
        let expectation = expectation(description: "Get component config and validate accessibility")
        let componentId = "tchat-button"
        let platform = "ios"

        // When
        apiClient.getComponentConfig(componentId: componentId, platform: platform) { result in
            // Then
            switch result {
            case .success(let config):
                let accessibility = config.accessibility
                XCTAssertNotNil(accessibility.label)
                XCTAssertNotNil(accessibility.hint)
                XCTAssertNotNil(accessibility.role)
                XCTAssertNotNil(accessibility.traits)
                XCTAssertNotNil(accessibility.minimumTouchSize)

                // Validate minimum touch size for iOS (44pt)
                let touchSize = accessibility.minimumTouchSize
                XCTAssertGreaterThanOrEqual(touchSize.width, 44.0)
                XCTAssertGreaterThanOrEqual(touchSize.height, 44.0)

            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetComponentConfig_ValidatesAnimationConfig() {
        // Given
        let expectation = expectation(description: "Get component config and validate animations")
        let componentId = "tchat-button"
        let platform = "ios"

        // When
        apiClient.getComponentConfig(componentId: componentId, platform: platform) { result in
            // Then
            switch result {
            case .success(let config):
                let animations = config.animations
                XCTAssertGreaterThan(animations.count, 0)

                for animation in animations {
                    XCTAssertNotNil(animation.trigger)
                    XCTAssertGreaterThan(animation.duration, 0)
                    XCTAssertNotNil(animation.easing)
                    XCTAssertGreaterThan(animation.properties.count, 0)
                }

            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }
}