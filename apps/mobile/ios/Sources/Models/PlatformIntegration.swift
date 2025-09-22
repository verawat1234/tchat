//
//  PlatformIntegration.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
#if canImport(UIKit)
import UIKit
#endif
#if canImport(AVFoundation)
import AVFoundation
#endif
#if canImport(UserNotifications)
import UserNotifications
#endif

/// Native mobile features that enhance web functionality
/// Implements the PlatformIntegration entity from data-model.md specification
public struct PlatformIntegration: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Types

    public enum Platform: String, Codable, CaseIterable {
        case ios = "ios"
        case android = "android"
    }

    public enum FallbackBehavior: String, Codable, CaseIterable {
        case disable = "disable"
        case webEquivalent = "web_equivalent"
        case alternate = "alternate"
    }

    public enum IntegrationState: String, Codable, CaseIterable {
        case requested = "requested"
        case checking = "checking"
        case available = "available"
        case unavailable = "unavailable"
        case configured = "configured"
        case active = "active"
        case denied = "denied"
        case fallback = "fallback"
        case disabled = "disabled"
    }

    // MARK: - Properties

    public let id: String
    public let name: String
    public let platform: Platform
    public let capability: String
    public let isAvailable: Bool
    public let permissions: [Permission]
    public let configuration: [String: String]
    public let fallbackBehavior: FallbackBehavior

    // State management
    public var currentState: IntegrationState
    public var lastStateChange: Date
    public var errorMessage: String?

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        name: String,
        platform: Platform = .ios,
        capability: String,
        isAvailable: Bool = false,
        permissions: [Permission] = [],
        configuration: [String: String] = [:],
        fallbackBehavior: FallbackBehavior = .webEquivalent,
        currentState: IntegrationState = .requested,
        lastStateChange: Date = Date(),
        errorMessage: String? = nil
    ) {
        self.id = id
        self.name = name
        self.platform = platform
        self.capability = capability
        self.isAvailable = isAvailable
        self.permissions = permissions
        self.configuration = configuration
        self.fallbackBehavior = fallbackBehavior
        self.currentState = currentState
        self.lastStateChange = lastStateChange
        self.errorMessage = errorMessage
    }

    // MARK: - Validation

    /// Validates the platform integration according to specification rules
    public func validate() throws {
        guard !id.isEmpty else {
            throw PlatformIntegrationError.invalidId("ID cannot be empty")
        }

        guard !name.isEmpty else {
            throw PlatformIntegrationError.invalidName("Name cannot be empty")
        }

        guard !capability.isEmpty else {
            throw PlatformIntegrationError.invalidCapability("Capability cannot be empty")
        }

        guard isValidCapability(capability) else {
            throw PlatformIntegrationError.invalidCapability("Capability '\(capability)' is not valid for platform '\(platform)'")
        }

        // Validate permissions for capability
        try permissions.forEach { permission in
            try permission.validate()
            guard isValidPermissionForCapability(permission.name, capability: capability) else {
                throw PlatformIntegrationError.invalidPermission("Permission '\(permission.name)' is not valid for capability '\(capability)'")
            }
        }
    }

    /// Validates that capability is supported on the platform
    private func isValidCapability(_ capability: String) -> Bool {
        switch platform {
        case .ios:
            return IOSCapabilities.allCases.map { $0.rawValue }.contains(capability)
        case .android:
            return true // Android allows more flexible capabilities
        }
    }

    /// Validates that permission is required for the capability
    private func isValidPermissionForCapability(_ permission: String, capability: String) -> Bool {
        guard let iosCapability = IOSCapabilities(rawValue: capability) else { return true }
        return iosCapability.requiredPermissions.contains(permission)
    }

    // MARK: - State Transitions

    /// Updates the integration state following valid transitions
    public mutating func updateState(to newState: IntegrationState, errorMessage: String? = nil) throws {
        guard isValidStateTransition(from: currentState, to: newState) else {
            throw PlatformIntegrationError.invalidStateTransition("Cannot transition from \(currentState) to \(newState)")
        }

        self.currentState = newState
        self.lastStateChange = Date()
        self.errorMessage = errorMessage
    }

    /// Validates state transitions according to specification
    private func isValidStateTransition(from: IntegrationState, to: IntegrationState) -> Bool {
        switch (from, to) {
        case (.requested, .checking):
            return true
        case (.checking, .available), (.checking, .unavailable):
            return true
        case (.available, .configured), (.available, .denied):
            return true
        case (.configured, .active), (.configured, .fallback):
            return true
        case (.denied, .fallback), (.denied, .disabled):
            return true
        case (.unavailable, .fallback), (.unavailable, .disabled):
            return true
        default:
            return false
        }
    }

    // MARK: - Runtime Availability

    /// Checks runtime availability of the integration
    public static func checkAvailability(for capability: String) -> Bool {
        guard let iosCapability = IOSCapabilities(rawValue: capability) else { return false }
        return iosCapability.isAvailable
    }

    /// Requests necessary permissions for the integration
    public func requestPermissions() async -> [Permission] {
        var updatedPermissions: [Permission] = []

        for permission in permissions {
            let status = await checkPermissionStatus(permission.name)
            let updatedPermission = Permission(
                name: permission.name,
                status: status,
                required: permission.required,
                requestReason: permission.requestReason
            )
            updatedPermissions.append(updatedPermission)
        }

        return updatedPermissions
    }

    /// Checks the current status of a specific permission
    private func checkPermissionStatus(_ permissionName: String) async -> Permission.PermissionStatus {
        switch permissionName {
        case "camera":
            let status = AVCaptureDevice.authorizationStatus(for: .video)
            return status.toPermissionStatus()
        case "microphone":
            let status = AVCaptureDevice.authorizationStatus(for: .audio)
            return status.toPermissionStatus()
        case "notifications":
            let settings = await UNUserNotificationCenter.current().notificationSettings()
            return settings.authorizationStatus.toPermissionStatus()
        case "location":
            // Note: This is simplified - real implementation would check CLLocationManager
            return .notDetermined
        default:
            return .notDetermined
        }
    }

    // MARK: - Permission Management

    /// Checks if all required permissions are granted
    public var hasRequiredPermissions: Bool {
        let requiredPermissions = permissions.filter { $0.required }
        return requiredPermissions.allSatisfy { $0.status == .granted }
    }

    /// Gets missing required permissions
    public var missingRequiredPermissions: [Permission] {
        return permissions.filter { $0.required && $0.status != .granted }
    }

    /// Checks if integration can function with current permissions
    public var canFunction: Bool {
        return isAvailable && hasRequiredPermissions
    }
}

// MARK: - Permission Model

public struct Permission: Codable, Equatable, Hashable {

    public enum PermissionStatus: String, Codable, CaseIterable {
        case granted = "granted"
        case denied = "denied"
        case notDetermined = "not_determined"
    }

    public let name: String
    public let status: PermissionStatus
    public let required: Bool
    public let requestReason: String

    public init(
        name: String,
        status: PermissionStatus = .notDetermined,
        required: Bool = true,
        requestReason: String
    ) {
        self.name = name
        self.status = status
        self.required = required
        self.requestReason = requestReason
    }

    /// Validates the permission
    public func validate() throws {
        guard !name.isEmpty else {
            throw PlatformIntegrationError.invalidPermission("Permission name cannot be empty")
        }

        guard !requestReason.isEmpty else {
            throw PlatformIntegrationError.invalidPermission("Request reason cannot be empty")
        }
    }
}

// MARK: - iOS Capabilities

public enum IOSCapabilities: String, CaseIterable {
    case camera = "camera"
    case microphone = "microphone"
    case notifications = "notifications"
    case location = "location"
    case biometrics = "biometrics"
    case haptics = "haptics"
    case shareSheet = "share_sheet"
    case deepLinking = "deep_linking"
    case backgroundProcessing = "background_processing"
    case pushNotifications = "push_notifications"

    var isAvailable: Bool {
        switch self {
        case .camera:
            return !UIDevice.current.isSimulator
        case .microphone:
            return !UIDevice.current.isSimulator
        case .notifications, .pushNotifications:
            return true
        case .location:
            return !UIDevice.current.isSimulator
        case .biometrics:
            return UIDevice.current.supportsBiometrics
        case .haptics:
            return UIDevice.current.supportsHaptics
        case .shareSheet, .deepLinking:
            return true
        case .backgroundProcessing:
            return true
        }
    }

    var requiredPermissions: [String] {
        switch self {
        case .camera:
            return ["camera"]
        case .microphone:
            return ["microphone"]
        case .notifications, .pushNotifications:
            return ["notifications"]
        case .location:
            return ["location"]
        case .biometrics, .haptics, .shareSheet, .deepLinking, .backgroundProcessing:
            return []
        }
    }
}

// MARK: - Error Types

public enum PlatformIntegrationError: LocalizedError {
    case invalidId(String)
    case invalidName(String)
    case invalidCapability(String)
    case invalidPermission(String)
    case invalidStateTransition(String)
    case permissionDenied(String)
    case unavailableCapability(String)

    public var errorDescription: String? {
        switch self {
        case .invalidId(let message),
             .invalidName(let message),
             .invalidCapability(let message),
             .invalidPermission(let message),
             .invalidStateTransition(let message),
             .permissionDenied(let message),
             .unavailableCapability(let message):
            return message
        }
    }
}

// MARK: - Factory Methods

extension PlatformIntegration {

    /// Creates a camera integration
    public static func camera(
        name: String = "Camera Access",
        configuration: [String: String] = [:]
    ) -> PlatformIntegration {
        return PlatformIntegration(
            name: name,
            capability: "camera",
            isAvailable: IOSCapabilities.camera.isAvailable,
            permissions: [
                Permission(
                    name: "camera",
                    required: true,
                    requestReason: "This app needs camera access to take photos and videos."
                )
            ],
            configuration: configuration,
            fallbackBehavior: .webEquivalent
        )
    }

    /// Creates a notifications integration
    public static func notifications(
        name: String = "Push Notifications",
        configuration: [String: String] = [:]
    ) -> PlatformIntegration {
        return PlatformIntegration(
            name: name,
            capability: "notifications",
            isAvailable: IOSCapabilities.notifications.isAvailable,
            permissions: [
                Permission(
                    name: "notifications",
                    required: true,
                    requestReason: "This app needs notification access to send you important updates."
                )
            ],
            configuration: configuration,
            fallbackBehavior: .disable
        )
    }

    /// Creates a biometrics integration
    public static func biometrics(
        name: String = "Biometric Authentication",
        configuration: [String: String] = [:]
    ) -> PlatformIntegration {
        return PlatformIntegration(
            name: name,
            capability: "biometrics",
            isAvailable: IOSCapabilities.biometrics.isAvailable,
            permissions: [],
            configuration: configuration,
            fallbackBehavior: .alternate
        )
    }
}

// MARK: - Extensions

private extension UIDevice {
    var isSimulator: Bool {
        return TARGET_OS_SIMULATOR != 0
    }

    var supportsBiometrics: Bool {
        // Simplified check - real implementation would use LocalAuthentication framework
        return !isSimulator
    }

    var supportsHaptics: Bool {
        // Check for haptic feedback capability
        return !isSimulator
    }
}

private extension AVAuthorizationStatus {
    func toPermissionStatus() -> Permission.PermissionStatus {
        switch self {
        case .authorized:
            return .granted
        case .denied, .restricted:
            return .denied
        case .notDetermined:
            return .notDetermined
        @unknown default:
            return .notDetermined
        }
    }
}

private extension UNAuthorizationStatus {
    func toPermissionStatus() -> Permission.PermissionStatus {
        switch self {
        case .authorized, .provisional:
            return .granted
        case .denied:
            return .denied
        case .notDetermined:
            return .notDetermined
        @unknown default:
            return .notDetermined
        }
    }
}