//
//  Colors.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Color design tokens matching TailwindCSS v4 color palette
public struct Colors {

    // MARK: - Brand Colors

    /// Primary brand color - Blue 500
    /// Maps to TailwindCSS: blue-500
    public let primary = Color(hex: "#3B82F6")

    /// Secondary brand color - Gray 600
    /// Maps to TailwindCSS: gray-600
    public let secondary = Color(hex: "#4B5563")

    /// Accent color - Indigo 500
    /// Maps to TailwindCSS: indigo-500
    public let accent = Color(hex: "#6366F1")

    // MARK: - Semantic Colors

    /// Success color - Green 500
    /// Maps to TailwindCSS: green-500
    public let success = Color(hex: "#10B981")

    /// Warning color - Amber 500
    /// Maps to TailwindCSS: amber-500
    public let warning = Color(hex: "#F59E0B")

    /// Error color - Red 500
    /// Maps to TailwindCSS: red-500
    public let error = Color(hex: "#EF4444")

    /// Info color - Blue 400
    /// Maps to TailwindCSS: blue-400
    public let info = Color(hex: "#60A5FA")

    // MARK: - Surface Colors

    /// Background color - White
    /// Maps to TailwindCSS: white
    public let background = Color(hex: "#FFFFFF")

    /// Surface color - Gray 50
    /// Maps to TailwindCSS: gray-50
    public let surface = Color(hex: "#F9FAFB")

    /// Card background - White with subtle shadow
    /// Maps to TailwindCSS: white
    public let cardBackground = Color(hex: "#FFFFFF")

    /// Modal overlay - Black with opacity
    /// Maps to TailwindCSS: black/50
    public let overlay = Color.black.opacity(0.5)

    // MARK: - Text Colors

    /// Primary text color - Gray 900
    /// Maps to TailwindCSS: gray-900
    public let textPrimary = Color(hex: "#111827")

    /// Secondary text color - Gray 600
    /// Maps to TailwindCSS: gray-600
    public let textSecondary = Color(hex: "#4B5563")

    /// Tertiary text color - Gray 400
    /// Maps to TailwindCSS: gray-400
    public let textTertiary = Color(hex: "#9CA3AF")

    /// Disabled text color - Gray 300
    /// Maps to TailwindCSS: gray-300
    public let textDisabled = Color(hex: "#D1D5DB")

    /// Text on primary color - White
    public let textOnPrimary = Color.white

    /// Text on dark backgrounds - White
    public let textOnDark = Color.white

    // MARK: - Border Colors

    /// Default border color - Gray 200
    /// Maps to TailwindCSS: gray-200
    public let border = Color(hex: "#E5E7EB")

    /// Focus border color - Blue 500
    /// Maps to TailwindCSS: blue-500
    public let borderFocus = Color(hex: "#3B82F6")

    /// Error border color - Red 300
    /// Maps to TailwindCSS: red-300
    public let borderError = Color(hex: "#FCA5A5")

    /// Divider color - Gray 200
    /// Maps to TailwindCSS: gray-200
    public let divider = Color(hex: "#E5E7EB")

    // MARK: - Interactive States

    /// Hover state colors
    public struct Hover {
        public static let primary = Color(hex: "#2563EB")     // blue-600
        public static let secondary = Color(hex: "#374151")   // gray-700
        public static let surface = Color(hex: "#F3F4F6")     // gray-100
    }

    /// Pressed state colors
    public struct Pressed {
        public static let primary = Color(hex: "#1D4ED8")     // blue-700
        public static let secondary = Color(hex: "#1F2937")   // gray-800
        public static let surface = Color(hex: "#E5E7EB")     // gray-200
    }

    /// Disabled state colors
    public struct Disabled {
        public static let background = Color(hex: "#F3F4F6")  // gray-100
        public static let text = Color(hex: "#9CA3AF")        // gray-400
        public static let border = Color(hex: "#D1D5DB")      // gray-300
    }

    // MARK: - Tab/Navigation Colors

    /// Tab bar background - Surface with blur
    public let tabBarBackground = Color(hex: "#F9FAFB").opacity(0.95)

    /// Navigation bar background - Background
    public var navigationBackground: Color { background }

    /// Selected tab color - Primary
    public var tabSelected: Color { primary }

    /// Unselected tab color - Gray 400
    public let tabUnselected = Color(hex: "#9CA3AF")

    // MARK: - Shadow Colors

    /// Light shadow color
    public let shadowLight = Color.black.opacity(0.1)

    /// Medium shadow color
    public let shadowMedium = Color.black.opacity(0.15)

    /// Heavy shadow color
    public let shadowHeavy = Color.black.opacity(0.25)
}

// MARK: - Dark Mode Colors

extension Colors {

    /// Dark mode color variants
    public struct Dark {

        // Brand colors remain the same
        public static let primary = Color(hex: "#3B82F6")
        public static let secondary = Color(hex: "#6B7280")
        public static let accent = Color(hex: "#6366F1")

        // Semantic colors adjusted for dark mode
        public static let success = Color(hex: "#10B981")
        public static let warning = Color(hex: "#F59E0B")
        public static let error = Color(hex: "#EF4444")
        public static let info = Color(hex: "#60A5FA")

        // Dark surfaces
        public static let background = Color(hex: "#111827")     // gray-900
        public static let surface = Color(hex: "#1F2937")        // gray-800
        public static let cardBackground = Color(hex: "#374151") // gray-700
        public static let overlay = Color.black.opacity(0.7)

        // Dark text colors
        public static let textPrimary = Color(hex: "#F9FAFB")    // gray-50
        public static let textSecondary = Color(hex: "#D1D5DB")  // gray-300
        public static let textTertiary = Color(hex: "#9CA3AF")   // gray-400
        public static let textDisabled = Color(hex: "#6B7280")   // gray-500

        // Dark borders
        public static let border = Color(hex: "#374151")         // gray-700
        public static let borderFocus = Color(hex: "#3B82F6")    // blue-500
        public static let borderError = Color(hex: "#EF4444")    // red-500
        public static let divider = Color(hex: "#374151")        // gray-700

        // Dark navigation
        public static let tabBarBackground = Color(hex: "#1F2937").opacity(0.95)
        public static let navigationBackground = Color(hex: "#111827")
    }
}

// MARK: - Color Utilities

extension Colors {

    /// Get colors for current color scheme
    public static func colorsFor(scheme: ColorScheme) -> Colors.Type {
        // In a full implementation, this would return scheme-specific colors
        return Colors.self
    }

    /// Validate color accessibility contrast
    public static func validateContrast(foreground: Color, background: Color) -> Bool {
        // Implementation would check WCAG contrast requirements
        return true // Placeholder
    }
}

// MARK: - System Color Integration

extension Color {

    /// Dynamic colors that adapt to light/dark mode
    #if canImport(UIKit)
    public static let adaptiveBackground = Color(.systemBackground)
    public static let adaptiveSecondaryBackground = Color(.secondarySystemBackground)
    public static let adaptiveTertiaryBackground = Color(.tertiarySystemBackground)
    public static let adaptiveLabel = Color(.label)
    public static let adaptiveSecondaryLabel = Color(.secondaryLabel)
    public static let adaptiveTertiaryLabel = Color(.tertiaryLabel)
    #else
    // Fallback colors for macOS/other platforms
    public static let adaptiveBackground = Color.white
    public static let adaptiveSecondaryBackground = Color(.sRGB, red: 0.97, green: 0.97, blue: 0.97)
    public static let adaptiveTertiaryBackground = Color(.sRGB, red: 0.95, green: 0.95, blue: 0.95)
    public static let adaptiveLabel = Color.black
    public static let adaptiveSecondaryLabel = Color(.sRGB, red: 0.3, green: 0.3, blue: 0.3)
    public static let adaptiveTertiaryLabel = Color(.sRGB, red: 0.6, green: 0.6, blue: 0.6)
    #endif
}