//
//  TestConfiguration.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import XCTest
import Foundation

/// Global test configuration and utilities for TchatApp
class TestConfiguration {

    /// Shared test configuration instance
    static let shared = TestConfiguration()

    private init() {}

    /// Test API base URL
    var apiBaseURL: String {
        return "https://api.tchat.test"
    }

    /// Test timeout intervals
    struct Timeouts {
        static let short: TimeInterval = 2.0
        static let medium: TimeInterval = 5.0
        static let long: TimeInterval = 10.0
    }

    /// Test data utilities
    struct TestData {
        static let validUserId = "test-user-123"
        static let validWorkspaceId = "test-workspace-456"
        static let validSessionToken = "test-session-token"
    }
}

/// Base test case for all TchatApp tests
class TchatAppTestCase: XCTestCase {

    override func setUp() {
        super.setUp()
        // Common test setup
        continueAfterFailure = false
    }

    override func tearDown() {
        // Common test cleanup
        super.tearDown()
    }

    /// Helper to create expectation with timeout
    func expectation(description: String, timeout: TimeInterval = TestConfiguration.Timeouts.medium) -> XCTestExpectation {
        return expectation(description: description)
    }

    /// Helper to wait for expectations
    func waitForExpectations(timeout: TimeInterval = TestConfiguration.Timeouts.medium) {
        waitForExpectations(timeout: timeout, handler: nil)
    }
}