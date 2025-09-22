//
//  PlatformAdapter.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
#if canImport(UIKit)
import UIKit
#endif

/// Platform adapter entity for handling platform-specific UI conventions and capabilities
public struct PlatformAdapter: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Properties

    public let id: String
    public let platform: String
    public let version: String
    public let capabilities: [PlatformCapability]
    public let uiConventions: UIConventions
    public let gestureSupport: GestureSupport
    public let animationSupport: AnimationSupport
    public let metadata: PlatformAdapterMetadata

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        platform: String = "ios",
        version: String = UIDevice.current.systemVersion,
        capabilities: [PlatformCapability] = [],
        uiConventions: UIConventions = UIConventions.defaultiOSConventions,
        gestureSupport: GestureSupport = GestureSupport.defaultiOSGestures,
        animationSupport: AnimationSupport = AnimationSupport.defaultiOSAnimations,
        metadata: PlatformAdapterMetadata = PlatformAdapterMetadata()
    ) {
        self.id = id
        self.platform = platform
        self.version = version
        self.capabilities = capabilities
        self.uiConventions = uiConventions
        self.gestureSupport = gestureSupport
        self.animationSupport = animationSupport
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Check if platform supports specific capability
    public func supportsCapability(_ capabilityName: String) -> Bool {
        return capabilities.contains { $0.name == capabilityName && $0.isSupported }
    }

    /// Get capability restrictions
    public func getCapabilityRestrictions(_ capabilityName: String) -> [String] {
        return capabilities.first { $0.name == capabilityName }?.restrictions ?? []
    }

    /// Check if gesture is supported
    public func supportsGesture(_ gestureType: String) -> Bool {
        return gestureSupport.supportedGestures.contains { $0.type == gestureType }
    }

    /// Check if animation is supported
    public func supportsAnimation(_ animationType: String) -> Bool {
        return animationSupport.supportedAnimations.contains { $0.type == animationType }
    }

    /// Get platform-specific color scheme
    public var colorScheme: [String: String] {
        return uiConventions.designSystem.colorScheme
    }

    /// Get platform-specific spacing
    public var spacing: [String: Double] {
        return uiConventions.designSystem.spacing
    }

    /// Get platform-specific border radius
    public var borderRadius: [String: Double] {
        return uiConventions.designSystem.borderRadius
    }
}

// MARK: - Supporting Types

/// Platform capability definition
public struct PlatformCapability: Codable, Identifiable, Equatable, Hashable {
    public let id: String
    public let name: String
    public let isSupported: Bool
    public let apiLevel: String
    public let restrictions: [String]
    public let alternativeActions: [String]
    public let description: String?

    public init(
        id: String = UUID().uuidString,
        name: String,
        isSupported: Bool,
        apiLevel: String,
        restrictions: [String] = [],
        alternativeActions: [String] = [],
        description: String? = nil
    ) {
        self.id = id
        self.name = name
        self.isSupported = isSupported
        self.apiLevel = apiLevel
        self.restrictions = restrictions
        self.alternativeActions = alternativeActions
        self.description = description
    }
}

/// UI conventions for platform-specific design patterns
public struct UIConventions: Codable, Equatable, Hashable {
    public let designSystem: DesignSystem
    public let navigationPatterns: [String: String]
    public let layoutPatterns: [String: String]
    public let accessibilityGuidelines: [String: String]

    public init(
        designSystem: DesignSystem,
        navigationPatterns: [String: String] = [:],
        layoutPatterns: [String: String] = [:],
        accessibilityGuidelines: [String: String] = [:]
    ) {
        self.designSystem = designSystem
        self.navigationPatterns = navigationPatterns
        self.layoutPatterns = layoutPatterns
        self.accessibilityGuidelines = accessibilityGuidelines
    }

    /// Default iOS UI conventions
    public static let defaultiOSConventions = UIConventions(
        designSystem: DesignSystem.defaultiOSDesignSystem,
        navigationPatterns: [
            "tabBar": "bottom",
            "navigationBar": "top",
            "modal": "slideUp",
            "push": "slideLeft"
        ],
        layoutPatterns: [
            "safeArea": "respect",
            "statusBar": "lightContent",
            "homeIndicator": "adapt"
        ],
        accessibilityGuidelines: [
            "voiceOver": "supported",
            "dynamicType": "supported",
            "hapticFeedback": "supported"
        ]
    )
}

/// Design system definition
public struct DesignSystem: Codable, Equatable, Hashable {
    public let colorScheme: [String: String]
    public let typography: [String: String]
    public let spacing: [String: Double]
    public let borderRadius: [String: Double]
    public let shadows: [String: String]

    public init(
        colorScheme: [String: String],
        typography: [String: String] = [:],
        spacing: [String: Double] = [:],
        borderRadius: [String: Double] = [:],
        shadows: [String: String] = [:]
    ) {
        self.colorScheme = colorScheme
        self.typography = typography
        self.spacing = spacing
        self.borderRadius = borderRadius
        self.shadows = shadows
    }

    /// Default iOS design system
    public static let defaultiOSDesignSystem = DesignSystem(
        colorScheme: [
            "primary": "#007AFF",
            "secondary": "#5856D6",
            "background": "#FFFFFF",
            "surface": "#F2F2F7",
            "text": "#000000",
            "textSecondary": "#8E8E93"
        ],
        typography: [
            "title1": "34pt",
            "title2": "28pt",
            "title3": "22pt",
            "headline": "17pt",
            "body": "17pt",
            "callout": "16pt",
            "caption": "12pt"
        ],
        spacing: [
            "xs": 4,
            "sm": 8,
            "md": 16,
            "lg": 24,
            "xl": 32,
            "xxl": 48
        ],
        borderRadius: [
            "sm": 4,
            "md": 8,
            "lg": 12,
            "xl": 16,
            "round": 50
        ],
        shadows: [
            "small": "0 1 3 rgba(0,0,0,0.1)",
            "medium": "0 4 6 rgba(0,0,0,0.1)",
            "large": "0 10 15 rgba(0,0,0,0.1)"
        ]
    )
}

/// Gesture support definition
public struct GestureSupport: Codable, Equatable, Hashable {
    public let supportedGestures: [GestureDefinition]
    public let maxSimultaneousGestures: Int
    public let gestureSettings: [String: String]

    public init(
        supportedGestures: [GestureDefinition],
        maxSimultaneousGestures: Int = 2,
        gestureSettings: [String: String] = [:]
    ) {
        self.supportedGestures = supportedGestures
        self.maxSimultaneousGestures = maxSimultaneousGestures
        self.gestureSettings = gestureSettings
    }

    /// Default iOS gesture support
    public static let defaultiOSGestures = GestureSupport(
        supportedGestures: [
            GestureDefinition(type: "tap", description: "Single tap"),
            GestureDefinition(type: "doubleTap", description: "Double tap"),
            GestureDefinition(type: "longPress", description: "Long press"),
            GestureDefinition(type: "swipe", description: "Swipe gesture"),
            GestureDefinition(type: "pan", description: "Pan/drag gesture"),
            GestureDefinition(type: "pinch", description: "Pinch to zoom"),
            GestureDefinition(type: "rotation", description: "Rotation gesture")
        ],
        maxSimultaneousGestures: 2,
        gestureSettings: [
            "hapticFeedback": "enabled",
            "edgeSwipe": "enabled",
            "homeIndicator": "adaptive"
        ]
    )
}

/// Gesture definition
public struct GestureDefinition: Codable, Equatable, Hashable {
    public let type: String
    public let description: String
    public let requiredFingers: Int
    public let minDuration: Double
    public let maxDuration: Double?

    public init(
        type: String,
        description: String,
        requiredFingers: Int = 1,
        minDuration: Double = 0.0,
        maxDuration: Double? = nil
    ) {
        self.type = type
        self.description = description
        self.requiredFingers = requiredFingers
        self.minDuration = minDuration
        self.maxDuration = maxDuration
    }
}

/// Animation support definition
public struct AnimationSupport: Codable, Equatable, Hashable {
    public let supportedAnimations: [AnimationDefinition]
    public let defaultDuration: Double
    public let animationSettings: [String: String]

    public init(
        supportedAnimations: [AnimationDefinition],
        defaultDuration: Double = 0.3,
        animationSettings: [String: String] = [:]
    ) {
        self.supportedAnimations = supportedAnimations
        self.defaultDuration = defaultDuration
        self.animationSettings = animationSettings
    }

    /// Default iOS animation support
    public static let defaultiOSAnimations = AnimationSupport(
        supportedAnimations: [
            AnimationDefinition(type: "fade", description: "Fade in/out"),
            AnimationDefinition(type: "slide", description: "Slide transition"),
            AnimationDefinition(type: "scale", description: "Scale transform"),
            AnimationDefinition(type: "rotate", description: "Rotation transform"),
            AnimationDefinition(type: "spring", description: "Spring animation"),
            AnimationDefinition(type: "bounce", description: "Bounce effect")
        ],
        defaultDuration: 0.3,
        animationSettings: [
            "preferredFrameRate": "60",
            "reducedMotion": "respected",
            "timing": "easeInOut"
        ]
    )
}

/// Animation definition
public struct AnimationDefinition: Codable, Equatable, Hashable {
    public let type: String
    public let description: String
    public let minDuration: Double
    public let maxDuration: Double
    public let defaultEasing: String

    public init(
        type: String,
        description: String,
        minDuration: Double = 0.1,
        maxDuration: Double = 2.0,
        defaultEasing: String = "easeInOut"
    ) {
        self.type = type
        self.description = description
        self.minDuration = minDuration
        self.maxDuration = maxDuration
        self.defaultEasing = defaultEasing
    }
}

/// Platform adapter metadata
public struct PlatformAdapterMetadata: Codable, Equatable, Hashable {
    public let deviceModel: String
    public let osVersion: String
    public let screenSize: String
    public let colorDepth: String
    public let refreshRate: String
    public let hasNotch: Bool
    public let supportsDarkMode: Bool
    public let accessibility: [String: Bool]

    public init(
        deviceModel: String = UIDevice.current.model,
        osVersion: String = UIDevice.current.systemVersion,
        screenSize: String = "\(UIScreen.main.bounds.width)x\(UIScreen.main.bounds.height)",
        colorDepth: String = "24bit",
        refreshRate: String = "60hz",
        hasNotch: Bool = false, // Simplified detection
        supportsDarkMode: Bool = true,
        accessibility: [String: Bool] = [:]
    ) {
        self.deviceModel = deviceModel
        self.osVersion = osVersion
        self.screenSize = screenSize
        self.colorDepth = colorDepth
        self.refreshRate = refreshRate
        self.hasNotch = hasNotch
        self.supportsDarkMode = supportsDarkMode
        self.accessibility = accessibility
    }
}

// MARK: - Default Platform Adapters

extension PlatformAdapter {

    /// Create default iOS platform adapter
    public static func defaultiOSAdapter() -> PlatformAdapter {
        let capabilities: [PlatformCapability] = [
            PlatformCapability(
                name: "hapticFeedback",
                isSupported: true,
                apiLevel: "iOS 10.0+",
                description: "Haptic feedback support"
            ),
            PlatformCapability(
                name: "faceID",
                isSupported: true,
                apiLevel: "iOS 11.0+",
                description: "Face ID authentication"
            ),
            PlatformCapability(
                name: "darkMode",
                isSupported: true,
                apiLevel: "iOS 13.0+",
                description: "Dark mode support"
            ),
            PlatformCapability(
                name: "widgetKit",
                isSupported: true,
                apiLevel: "iOS 14.0+",
                description: "Home screen widgets"
            )
        ]

        return PlatformAdapter(
            platform: "ios",
            version: UIDevice.current.systemVersion,
            capabilities: capabilities,
            uiConventions: UIConventions.defaultiOSConventions,
            gestureSupport: GestureSupport.defaultiOSGestures,
            animationSupport: AnimationSupport.defaultiOSAnimations,
            metadata: PlatformAdapterMetadata()
        )
    }
}