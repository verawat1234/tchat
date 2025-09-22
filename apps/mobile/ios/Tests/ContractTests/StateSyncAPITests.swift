//
//  StateSyncAPITests.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import XCTest
@testable import TchatApp

/// Contract tests for State Synchronization API
/// These tests verify the API contract for cross-platform state management
class StateSyncAPITests: TchatAppTestCase {

    private var apiClient: StateSyncAPIClient!

    override func setUp() {
        super.setUp()
        apiClient = StateSyncAPIClient(baseURL: TestConfiguration.shared.apiBaseURL)
    }

    override func tearDown() {
        apiClient = nil
        super.tearDown()
    }

    // MARK: - GET /api/sync/state Tests

    func testGetSyncState_ValidUser_ReturnsState() {
        // Given
        let expectation = expectation(description: "Get sync state")
        let userId = TestConfiguration.TestData.validUserId
        let platform = "ios"

        // When
        apiClient.getSyncState(userId: userId, platform: platform, timestamp: nil) { result in
            // Then
            switch result {
            case .success(let state):
                XCTAssertEqual(state.userId, userId)
                XCTAssertEqual(state.platform, platform)
                XCTAssertNotNil(state.timestamp)
                XCTAssertGreaterThan(state.version, 0)
                XCTAssertNotNil(state.preferences)
                XCTAssertNotNil(state.navigation)
                XCTAssertNotNil(state.workspace)
            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetSyncState_WithTimestamp_ReturnsOnlyChanges() {
        // Given
        let expectation = expectation(description: "Get sync state with timestamp")
        let userId = TestConfiguration.TestData.validUserId
        let platform = "ios"
        let timestamp = Date(timeIntervalSinceNow: -3600) // 1 hour ago

        // When
        apiClient.getSyncState(userId: userId, platform: platform, timestamp: timestamp) { result in
            // Then
            switch result {
            case .success(let state):
                XCTAssertEqual(state.userId, userId)
                XCTAssertGreaterThan(state.timestamp, timestamp)
            case .failure(let error):
                if error.statusCode == 304 {
                    // No changes since timestamp - expected
                } else {
                    XCTFail("Expected success or 304, got error: \(error)")
                }
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testGetSyncState_UnauthorizedUser_ReturnsError() {
        // Given
        let expectation = expectation(description: "Get sync state with unauthorized user")
        let userId = "unauthorized-user"
        let platform = "ios"

        // When
        apiClient.getSyncState(userId: userId, platform: platform, timestamp: nil) { result in
            // Then
            switch result {
            case .success:
                XCTFail("Expected error for unauthorized user")
            case .failure(let error):
                XCTAssertEqual(error.statusCode, 401)
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    // MARK: - POST /api/sync/state Tests

    func testUpdateSyncState_ValidUpdate_ReturnsSuccess() {
        // Given
        let expectation = expectation(description: "Update sync state")
        let userId = TestConfiguration.TestData.validUserId
        let update = StateSyncUpdate(
            userId: userId,
            platform: "ios",
            timestamp: Date(),
            changes: [
                StateChange(
                    path: "preferences.theme",
                    operation: .set,
                    value: "dark",
                    oldValue: "light"
                )
            ]
        )

        // When
        apiClient.updateSyncState(update: update) { result in
            // Then
            switch result {
            case .success(let response):
                XCTAssertTrue(response.success)
                XCTAssertGreaterThan(response.version, 0)
                XCTAssertNotNil(response.timestamp)
                XCTAssertEqual(response.conflicts.count, 0)
            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testUpdateSyncState_ConflictingUpdate_ReturnsConflict() {
        // Given
        let expectation = expectation(description: "Update sync state with conflict")
        let userId = TestConfiguration.TestData.validUserId
        let update = StateSyncUpdate(
            userId: userId,
            platform: "ios",
            timestamp: Date(timeIntervalSinceNow: -3600), // Old timestamp
            changes: [
                StateChange(
                    path: "preferences.theme",
                    operation: .set,
                    value: "dark",
                    oldValue: "light"
                )
            ]
        )

        // When
        apiClient.updateSyncState(update: update) { result in
            // Then
            switch result {
            case .success(let response):
                // Might succeed with conflicts resolved
                XCTAssertNotNil(response.conflicts)
            case .failure(let error):
                XCTAssertEqual(error.statusCode, 409) // Conflict
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    // MARK: - POST /api/sync/notifications Tests

    func testSyncNotificationPreferences_ValidPreferences_ReturnsSuccess() {
        // Given
        let expectation = expectation(description: "Sync notification preferences")
        let preferences = NotificationPreferences(
            chat: true,
            store: false,
            social: true,
            video: true,
            workspace: false,
            marketing: false,
            quietHours: QuietHours(
                enabled: true,
                startTime: "22:00",
                endTime: "08:00",
                timezone: "America/New_York"
            )
        )

        // When
        apiClient.syncNotificationPreferences(preferences: preferences) { result in
            // Then
            switch result {
            case .success:
                // Success expected
                break
            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    // MARK: - POST /api/sync/workspace Tests

    func testSwitchWorkspace_ValidWorkspace_ReturnsWorkspaceState() {
        // Given
        let expectation = expectation(description: "Switch workspace")
        let workspaceId = TestConfiguration.TestData.validWorkspaceId
        let platform = "ios"

        // When
        apiClient.switchWorkspace(workspaceId: workspaceId, platform: platform) { result in
            // Then
            switch result {
            case .success(let workspaceState):
                XCTAssertEqual(workspaceState.currentWorkspaceId, workspaceId)
                XCTAssertGreaterThan(workspaceState.availableWorkspaces.count, 0)
                XCTAssertNotNil(workspaceState.role)

                // Validate workspace info structure
                for workspace in workspaceState.availableWorkspaces {
                    XCTAssertFalse(workspace.id.isEmpty)
                    XCTAssertFalse(workspace.name.isEmpty)
                    XCTAssertNotNil(workspace.type)
                    XCTAssertGreaterThanOrEqual(workspace.unreadCount, 0)
                }

            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    func testSwitchWorkspace_InvalidWorkspace_ReturnsError() {
        // Given
        let expectation = expectation(description: "Switch to invalid workspace")
        let workspaceId = "invalid-workspace-id"
        let platform = "ios"

        // When
        apiClient.switchWorkspace(workspaceId: workspaceId, platform: platform) { result in
            // Then
            switch result {
            case .success:
                XCTFail("Expected error for invalid workspace")
            case .failure(let error):
                XCTAssertEqual(error.statusCode, 404)
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }

    // MARK: - Validation Tests

    func testValidateUserPreferencesStructure() {
        // Given
        let expectation = expectation(description: "Validate user preferences structure")
        let userId = TestConfiguration.TestData.validUserId
        let platform = "ios"

        // When
        apiClient.getSyncState(userId: userId, platform: platform, timestamp: nil) { result in
            // Then
            switch result {
            case .success(let state):
                let preferences = state.preferences
                XCTAssertNotNil(preferences.theme)
                XCTAssertNotNil(preferences.language)
                XCTAssertNotNil(preferences.notifications)
                XCTAssertNotNil(preferences.accessibility)

                // Validate theme values
                let validThemes = ["light", "dark", "auto"]
                XCTAssertTrue(validThemes.contains(preferences.theme))

                // Validate accessibility preferences
                let accessibility = preferences.accessibility
                XCTAssertNotNil(accessibility.reducedMotion)
                XCTAssertNotNil(accessibility.highContrast)
                XCTAssertNotNil(accessibility.largerText)
                XCTAssertNotNil(accessibility.screenReader)

            case .failure(let error):
                XCTFail("Expected success, got error: \(error)")
            }
            expectation.fulfill()
        }

        waitForExpectations(timeout: TestConfiguration.Timeouts.medium)
    }
}