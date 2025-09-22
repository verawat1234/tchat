//
//  DesignTokens.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI
import Foundation

/// Central design token system for TchatApp
/// Provides platform-specific implementation of TailwindCSS v4 design tokens
public struct DesignTokens {

    /// Current design system version
    public static let version = "1.0.0"

    /// Typography design tokens
    public static let typography = Typography()

    /// Color design tokens
    public static let colors = Colors()

    /// Spacing design tokens
    public static let spacing = Spacing()

    /// Animation design tokens
    public static let animations = Animations()

    /// Border radius design tokens
    public static let borderRadius = BorderRadius()

    /// Shadow design tokens
    public static let shadows = Shadows()

    // MARK: - Private initializer
    private init() {}
}

// MARK: - Design Token Extensions

extension DesignTokens {

    /// Get design tokens for current theme
    public static func tokensFor(theme: AppTheme) -> DesignTokens.Type {
        // In a full implementation, this would return theme-specific tokens
        return DesignTokens.self
    }

    /// Validate design token consistency
    public static func validateTokens() -> Bool {
        // Validate that all required tokens are present and valid
        return !version.isEmpty &&
               colors.primary != nil &&
               typography.bodyMedium != nil &&
               Spacing.md > 0
    }
}

// MARK: - App Theme

public enum AppTheme: String, CaseIterable {
    case light = "light"
    case dark = "dark"
    case auto = "auto"

    public var displayName: String {
        switch self {
        case .light: return "Light"
        case .dark: return "Dark"
        case .auto: return "Auto"
        }
    }
}

// MARK: - Color Extensions

extension Color {
    /// Create color from hex string
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (1, 1, 1, 0)
        }

        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue:  Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}

// MARK: - Missing Design Token Implementations

/// Animation design tokens (placeholder implementation)
public struct Animations {
    public init() {}

    // Animation duration constants
    public static let fast: Double = 0.2
    public static let normal: Double = 0.3
    public static let slow: Double = 0.5
}

/// Border radius design tokens (placeholder implementation)
public struct BorderRadius {
    public init() {}

    // Border radius constants
    public static let none: CGFloat = 0
    public static let small: CGFloat = 4
    public static let medium: CGFloat = 8
    public static let large: CGFloat = 12
    public static let full: CGFloat = 9999
}

/// Shadow design tokens (placeholder implementation)
public struct Shadows {
    public init() {}

    // Shadow constants
    public static let small = (offset: CGSize(width: 0, height: 1), radius: CGFloat(2), opacity: Double(0.1))
    public static let medium = (offset: CGSize(width: 0, height: 4), radius: CGFloat(6), opacity: Double(0.15))
    public static let large = (offset: CGSize(width: 0, height: 10), radius: CGFloat(15), opacity: Double(0.2))
}