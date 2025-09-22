//
//  Spacing.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Spacing design tokens matching TailwindCSS v4 spacing scale
public struct Spacing {

    // MARK: - Base Spacing Scale
    // Following TailwindCSS 4px base unit system

    /// Extra extra small spacing (2pt)
    /// Maps to TailwindCSS: space-0.5 (0.125rem)
    public static let xxs: CGFloat = 2

    /// Extra small spacing (4pt)
    /// Maps to TailwindCSS: space-1 (0.25rem)
    public static let xs: CGFloat = 4

    /// Small spacing (8pt)
    /// Maps to TailwindCSS: space-2 (0.5rem)
    public static let sm: CGFloat = 8

    /// Medium spacing (16pt) - Most common spacing
    /// Maps to TailwindCSS: space-4 (1rem)
    public static let md: CGFloat = 16

    /// Large spacing (24pt)
    /// Maps to TailwindCSS: space-6 (1.5rem)
    public static let lg: CGFloat = 24

    /// Extra large spacing (32pt)
    /// Maps to TailwindCSS: space-8 (2rem)
    public static let xl: CGFloat = 32

    /// Extra extra large spacing (48pt)
    /// Maps to TailwindCSS: space-12 (3rem)
    public static let xxl: CGFloat = 48

    /// Extra extra extra large spacing (64pt)
    /// Maps to TailwindCSS: space-16 (4rem)
    public static let xxxl: CGFloat = 64

    // MARK: - Semantic Spacing

    /// Component padding
    public struct Component {
        /// Button padding
        public static let buttonPadding = EdgeInsets(top: sm, leading: md, bottom: sm, trailing: md)

        /// Card padding
        public static let cardPadding = EdgeInsets(top: md, leading: md, bottom: md, trailing: md)

        /// Input field padding
        public static let inputPadding = EdgeInsets(top: xs, leading: sm, bottom: xs, trailing: sm)

        /// Modal padding
        public static let modalPadding = EdgeInsets(top: lg, leading: lg, bottom: lg, trailing: lg)
    }

    /// Layout spacing
    public struct Layout {
        /// Screen edge margins
        public static let screenMargin = md

        /// Section spacing
        public static let sectionSpacing = lg

        /// Content spacing
        public static let contentSpacing = md

        /// List item spacing
        public static let listItemSpacing = sm

        /// Grid gap
        public static let gridGap = md
    }

    /// Navigation spacing
    public struct Navigation {
        /// Tab bar height
        public static let tabBarHeight: CGFloat = 49 // iOS standard

        /// Navigation bar height
        public static let navigationBarHeight: CGFloat = 44 // iOS standard

        /// Tab item spacing
        public static let tabItemSpacing = xs

        /// Navigation padding
        public static let navigationPadding = md
    }

    /// Interactive spacing
    public struct Interactive {
        /// Minimum touch target size (iOS HIG)
        public static let touchTargetSize: CGFloat = 44

        /// Button spacing
        public static let buttonSpacing = sm

        /// Form field spacing
        public static let formFieldSpacing = md

        /// Icon spacing
        public static let iconSpacing = xs
    }

    // MARK: - Responsive Spacing

    /// Get spacing that adapts to screen size
    public static func adaptiveSpacing(compact: CGFloat, regular: CGFloat) -> CGFloat {
        // In a full implementation, this would check size class
        return regular // Placeholder
    }

    /// Scale spacing for accessibility
    public static func scaledSpacing(_ spacing: CGFloat, scaleFactor: CGFloat = 1.0) -> CGFloat {
        return spacing * scaleFactor
    }
}

// MARK: - SwiftUI Extensions

extension EdgeInsets {
    /// Create symmetric edge insets
    public static func symmetric(horizontal: CGFloat = 0, vertical: CGFloat = 0) -> EdgeInsets {
        return EdgeInsets(top: vertical, leading: horizontal, bottom: vertical, trailing: horizontal)
    }

    /// Create uniform edge insets
    public static func uniform(_ value: CGFloat) -> EdgeInsets {
        return EdgeInsets(top: value, leading: value, bottom: value, trailing: value)
    }

    /// Common edge insets using design tokens
    public static let small = EdgeInsets.uniform(Spacing.sm)
    public static let medium = EdgeInsets.uniform(Spacing.md)
    public static let large = EdgeInsets.uniform(Spacing.lg)
}

extension View {
    /// Apply spacing using design tokens
    public func spacing(_ spacing: CGFloat) -> some View {
        self.padding(spacing)
    }

    /// Apply semantic component padding
    public func componentPadding(_ type: ComponentPaddingType) -> some View {
        switch type {
        case .button:
            return self.padding(Spacing.Component.buttonPadding)
        case .card:
            return self.padding(Spacing.Component.cardPadding)
        case .input:
            return self.padding(Spacing.Component.inputPadding)
        case .modal:
            return self.padding(Spacing.Component.modalPadding)
        }
    }
}

/// Component padding types
public enum ComponentPaddingType {
    case button
    case card
    case input
    case modal
}

// MARK: - Spacing Utilities

extension Spacing {

    /// Get spacing for semantic usage
    public static func spacingFor(usage: SpacingUsage) -> CGFloat {
        switch usage {
        case .elementSpacing: return xs
        case .componentSpacing: return sm
        case .sectionSpacing: return md
        case .pageSpacing: return lg
        case .screenMargin: return md
        case .formFieldSpacing: return md
        case .listItemSpacing: return sm
        case .buttonSpacing: return sm
        case .iconPadding: return xs
        case .cardPadding: return md
        }
    }

    /// Validate spacing consistency (multiples of 4pt)
    public static func isValidSpacing(_ spacing: CGFloat) -> Bool {
        return spacing.truncatingRemainder(dividingBy: 4) == 0
    }
}

/// Semantic spacing usage enum
public enum SpacingUsage {
    case elementSpacing     // Between small elements
    case componentSpacing   // Between components
    case sectionSpacing     // Between sections
    case pageSpacing        // Page margins
    case screenMargin       // Screen edge margins
    case formFieldSpacing   // Between form fields
    case listItemSpacing    // Between list items
    case buttonSpacing      // Between buttons
    case iconPadding        // Around icons
    case cardPadding        // Inside cards
}