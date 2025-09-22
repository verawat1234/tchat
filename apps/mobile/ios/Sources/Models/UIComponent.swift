//
//  UIComponent.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI

/// Core UI component entity with state management and cross-platform sync
public struct UIComponent: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Properties

    public let id: String
    public let name: String
    public let category: ComponentCategory
    public let version: String
    public let isStateful: Bool
    public let propsSchema: ComponentSchema
    public let stateSchema: ComponentSchema?
    public let eventsSchema: ComponentSchema
    public let dependencies: [String]
    public let platformSupport: [String]
    public let metadata: ComponentMetadata

    // MARK: - Initialization

    public init(
        id: String,
        name: String,
        category: ComponentCategory,
        version: String = "1.0.0",
        isStateful: Bool = false,
        propsSchema: ComponentSchema,
        stateSchema: ComponentSchema? = nil,
        eventsSchema: ComponentSchema = ComponentSchema(),
        dependencies: [String] = [],
        platformSupport: [String] = ["ios", "android", "web"],
        metadata: ComponentMetadata = ComponentMetadata()
    ) {
        self.id = id
        self.name = name
        self.category = category
        self.version = version
        self.isStateful = isStateful
        self.propsSchema = propsSchema
        self.stateSchema = stateSchema
        self.eventsSchema = eventsSchema
        self.dependencies = dependencies
        self.platformSupport = platformSupport
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Check if component supports current platform
    public var isSupported: Bool {
        return platformSupport.contains("ios")
    }

    /// Get component size estimation
    public var estimatedSize: ComponentSize {
        let propCount = propsSchema.properties.count
        let stateCount = stateSchema?.properties.count ?? 0
        let eventCount = eventsSchema.properties.count

        let complexity = propCount + stateCount + eventCount

        switch complexity {
        case 0...5: return .small
        case 6...15: return .medium
        case 16...30: return .large
        default: return .extraLarge
        }
    }

    /// Check if component has external dependencies
    public var hasExternalDependencies: Bool {
        return !dependencies.isEmpty
    }

    /// Generate component instance with state
    public func createInstance(
        instanceId: String,
        initialState: [String: Any] = [:],
        props: [String: Any] = [:]
    ) -> UIComponentInstance {
        return UIComponentInstance(
            instanceId: instanceId,
            componentId: id,
            componentName: name,
            state: isStateful ? initialState : [:],
            props: props,
            version: version,
            createdAt: Date(),
            lastUpdated: Date()
        )
    }
}

// MARK: - Supporting Types

/// Component category classification
public enum ComponentCategory: String, Codable, CaseIterable {
    case input
    case display
    case navigation
    case layout
    case feedback
    case overlay
    case media
    case data
    case form
    case interactive
}

/// Component size classification
public enum ComponentSize: String, Codable, CaseIterable {
    case small
    case medium
    case large
    case extraLarge
}

/// Component schema definition
public struct ComponentSchema: Codable, Equatable, Hashable {
    public let properties: [String: PropertyDefinition]
    public let required: [String]
    public let additionalProperties: Bool

    public init(
        properties: [String: PropertyDefinition] = [:],
        required: [String] = [],
        additionalProperties: Bool = false
    ) {
        self.properties = properties
        self.required = required
        self.additionalProperties = additionalProperties
    }
}

/// Property definition for component schema
public struct PropertyDefinition: Codable, Equatable, Hashable {
    public let type: PropertyType
    public let description: String?
    public let defaultValue: String?
    public let validation: PropertyValidation?
    public let isOptional: Bool

    public init(
        type: PropertyType,
        description: String? = nil,
        defaultValue: String? = nil,
        validation: PropertyValidation? = nil,
        isOptional: Bool = true
    ) {
        self.type = type
        self.description = description
        self.defaultValue = defaultValue
        self.validation = validation
        self.isOptional = isOptional
    }
}

/// Property data types
public enum PropertyType: String, Codable, CaseIterable {
    case string
    case number
    case boolean
    case array
    case object
    case function
    case color
    case dimension
}

/// Property validation rules
public struct PropertyValidation: Codable, Equatable, Hashable {
    public let minValue: Double?
    public let maxValue: Double?
    public let pattern: String?
    public let enumValues: [String]?

    public init(
        minValue: Double? = nil,
        maxValue: Double? = nil,
        pattern: String? = nil,
        enumValues: [String]? = nil
    ) {
        self.minValue = minValue
        self.maxValue = maxValue
        self.pattern = pattern
        self.enumValues = enumValues
    }
}

/// Component metadata
public struct ComponentMetadata: Codable, Equatable, Hashable {
    public let author: String
    public let createdAt: Date
    public let updatedAt: Date
    public let description: String?
    public let documentation: String?
    public let examples: [String]
    public let tags: [String]
    public let isDeprecated: Bool
    public let deprecationMessage: String?

    public init(
        author: String = "Tchat Team",
        createdAt: Date = Date(),
        updatedAt: Date = Date(),
        description: String? = nil,
        documentation: String? = nil,
        examples: [String] = [],
        tags: [String] = [],
        isDeprecated: Bool = false,
        deprecationMessage: String? = nil
    ) {
        self.author = author
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.description = description
        self.documentation = documentation
        self.examples = examples
        self.tags = tags
        self.isDeprecated = isDeprecated
        self.deprecationMessage = deprecationMessage
    }
}

/// Component instance with runtime state
public struct UIComponentInstance: Codable, Identifiable, Equatable {
    public let id: String
    public let instanceId: String
    public let componentId: String
    public let componentName: String
    public var state: [String: Any]
    public let props: [String: Any]
    public let version: String
    public let createdAt: Date
    public var lastUpdated: Date

    public init(
        instanceId: String,
        componentId: String,
        componentName: String,
        state: [String: Any] = [:],
        props: [String: Any] = [:],
        version: String,
        createdAt: Date,
        lastUpdated: Date
    ) {
        self.id = instanceId
        self.instanceId = instanceId
        self.componentId = componentId
        self.componentName = componentName
        self.state = state
        self.props = props
        self.version = version
        self.createdAt = createdAt
        self.lastUpdated = lastUpdated
    }

    /// Update component state
    public mutating func updateState(_ newState: [String: Any]) {
        state = newState
        lastUpdated = Date()
    }

    /// Merge state with existing state
    public mutating func mergeState(_ partialState: [String: Any]) {
        for (key, value) in partialState {
            state[key] = value
        }
        lastUpdated = Date()
    }
}

// MARK: - Codable Conformance for [String: Any]

extension UIComponent {
    private enum CodingKeys: String, CodingKey {
        case id, name, category, version, isStateful
        case propsSchema, stateSchema, eventsSchema
        case dependencies, platformSupport, metadata
    }
}

extension UIComponentInstance {
    private enum CodingKeys: String, CodingKey {
        case instanceId, componentId, componentName
        case version, createdAt, lastUpdated
        // Note: state and props require custom encoding/decoding
    }
}

// MARK: - Default Components

extension UIComponent {

    /// Default UI components for the application
    public static let defaultComponents: [UIComponent] = [
        UIComponent(
            id: "chat-message",
            name: "ChatMessage",
            category: .display,
            isStateful: true,
            propsSchema: ComponentSchema(
                properties: [
                    "message": PropertyDefinition(type: .string, description: "Message content"),
                    "author": PropertyDefinition(type: .string, description: "Message author"),
                    "timestamp": PropertyDefinition(type: .string, description: "Message timestamp"),
                    "isOwn": PropertyDefinition(type: .boolean, description: "Is own message")
                ],
                required: ["message", "author"]
            ),
            stateSchema: ComponentSchema(
                properties: [
                    "isRead": PropertyDefinition(type: .boolean, defaultValue: "false"),
                    "isSelected": PropertyDefinition(type: .boolean, defaultValue: "false")
                ]
            ),
            eventsSchema: ComponentSchema(
                properties: [
                    "onPress": PropertyDefinition(type: .function, description: "Handle message press"),
                    "onLongPress": PropertyDefinition(type: .function, description: "Handle long press")
                ]
            ),
            metadata: ComponentMetadata(description: "Chat message display component")
        ),
        UIComponent(
            id: "user-avatar",
            name: "UserAvatar",
            category: .display,
            isStateful: true,
            propsSchema: ComponentSchema(
                properties: [
                    "userId": PropertyDefinition(type: .string, description: "User identifier"),
                    "imageUrl": PropertyDefinition(type: .string, description: "Avatar image URL"),
                    "size": PropertyDefinition(type: .dimension, defaultValue: "40"),
                    "showOnlineStatus": PropertyDefinition(type: .boolean, defaultValue: "true")
                ],
                required: ["userId"]
            ),
            stateSchema: ComponentSchema(
                properties: [
                    "isOnline": PropertyDefinition(type: .boolean, defaultValue: "false"),
                    "lastSeen": PropertyDefinition(type: .string, isOptional: true)
                ]
            ),
            metadata: ComponentMetadata(description: "User avatar with status indicator")
        ),
        UIComponent(
            id: "navigation-tab",
            name: "NavigationTab",
            category: .navigation,
            isStateful: true,
            propsSchema: ComponentSchema(
                properties: [
                    "title": PropertyDefinition(type: .string, description: "Tab title"),
                    "icon": PropertyDefinition(type: .string, description: "Tab icon name"),
                    "badge": PropertyDefinition(type: .number, description: "Badge count", isOptional: true),
                    "isDisabled": PropertyDefinition(type: .boolean, defaultValue: "false")
                ],
                required: ["title", "icon"]
            ),
            stateSchema: ComponentSchema(
                properties: [
                    "isActive": PropertyDefinition(type: .boolean, defaultValue: "false"),
                    "hasNotification": PropertyDefinition(type: .boolean, defaultValue: "false")
                ]
            ),
            eventsSchema: ComponentSchema(
                properties: [
                    "onSelect": PropertyDefinition(type: .function, description: "Handle tab selection")
                ]
            ),
            metadata: ComponentMetadata(description: "Bottom navigation tab component")
        )
    ]
}