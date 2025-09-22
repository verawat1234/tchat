//
//  Typography.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

#if canImport(UIKit)
import UIKit
#endif

/// Typography design tokens matching TailwindCSS v4 text scales
public struct Typography {

    // MARK: - Heading Styles

    /// Large heading style (32pt, bold)
    /// Maps to TailwindCSS: text-3xl font-bold
    public let headingLarge = Font.system(size: 32, weight: .bold, design: .default)

    /// Medium heading style (24pt, semibold)
    /// Maps to TailwindCSS: text-2xl font-semibold
    public let headingMedium = Font.system(size: 24, weight: .semibold, design: .default)

    /// Small heading style (20pt, semibold)
    /// Maps to TailwindCSS: text-xl font-semibold
    public let headingSmall = Font.system(size: 20, weight: .semibold, design: .default)

    // MARK: - Body Styles

    /// Large body text (18pt, regular)
    /// Maps to TailwindCSS: text-lg
    public let bodyLarge = Font.system(size: 18, weight: .regular, design: .default)

    /// Medium body text (16pt, regular) - Default body text
    /// Maps to TailwindCSS: text-base
    public let bodyMedium = Font.system(size: 16, weight: .regular, design: .default)

    /// Small body text (14pt, regular)
    /// Maps to TailwindCSS: text-sm
    public let bodySmall = Font.system(size: 14, weight: .regular, design: .default)

    // MARK: - Label Styles

    /// Caption text (12pt, medium)
    /// Maps to TailwindCSS: text-xs font-medium
    public let caption = Font.system(size: 12, weight: .medium, design: .default)

    /// Label text (14pt, medium)
    /// Maps to TailwindCSS: text-sm font-medium
    public let label = Font.system(size: 14, weight: .medium, design: .default)

    /// Button text (16pt, semibold)
    /// Maps to TailwindCSS: text-base font-semibold
    public let button = Font.system(size: 16, weight: .semibold, design: .default)

    // MARK: - Utility Styles

    /// Overline text (10pt, bold, uppercase)
    /// Maps to TailwindCSS: text-xs font-bold uppercase tracking-wide
    public let overline = Font.system(size: 10, weight: .bold, design: .default)

    /// Monospace text for code (14pt, regular, monospaced)
    /// Maps to TailwindCSS: text-sm font-mono
    public let code = Font.system(size: 14, weight: .regular, design: .monospaced)

    // MARK: - Accessibility Support

    /// Get scaled font for accessibility
    public func scaledFont(_ font: Font) -> Font {
        #if canImport(UIKit)
        // This would integrate with UIContentSizeCategory for dynamic type
        return font
        #else
        return font
        #endif
    }

    // MARK: - Line Heights

    /// Line height multipliers matching TailwindCSS leading values
    public struct LineHeight {
        public static let tight: CGFloat = 1.25     // leading-tight
        public static let normal: CGFloat = 1.5     // leading-normal
        public static let relaxed: CGFloat = 1.625  // leading-relaxed
        public static let loose: CGFloat = 2.0      // leading-loose
    }

    // MARK: - Letter Spacing

    /// Letter spacing values matching TailwindCSS tracking values
    public struct LetterSpacing {
        public static let tighter: CGFloat = -0.05  // tracking-tighter
        public static let tight: CGFloat = -0.025   // tracking-tight
        public static let normal: CGFloat = 0       // tracking-normal
        public static let wide: CGFloat = 0.025     // tracking-wide
        public static let wider: CGFloat = 0.05     // tracking-wider
        public static let widest: CGFloat = 0.1     // tracking-widest
    }
}

// MARK: - Font Weight Extensions

extension Font.Weight {
    /// TailwindCSS font weight mapping
    public static let thin = Font.Weight.ultraLight        // font-thin
    public static let extralight = Font.Weight.ultraLight  // font-extralight
    public static let light = Font.Weight.light            // font-light
    public static let normal = Font.Weight.regular         // font-normal
    public static let medium = Font.Weight.medium          // font-medium
    public static let semibold = Font.Weight.semibold      // font-semibold
    public static let bold = Font.Weight.bold              // font-bold
    public static let extrabold = Font.Weight.heavy        // font-extrabold
    public static let black = Font.Weight.black            // font-black
}

// MARK: - Typography Utilities

extension Typography {

    /// Get font for semantic usage
    public func fontFor(usage: TypographyUsage) -> Font {
        switch usage {
        case .pageTitle: return headingLarge
        case .sectionTitle: return headingMedium
        case .cardTitle: return headingSmall
        case .bodyText: return bodyMedium
        case .secondaryText: return bodySmall
        case .captionText: return caption
        case .buttonLabel: return button
        case .navigationTitle: return headingMedium
        case .tabLabel: return label
        case .inputLabel: return label
        case .errorText: return bodySmall
        }
    }
}

/// Semantic typography usage enum
public enum TypographyUsage {
    case pageTitle
    case sectionTitle
    case cardTitle
    case bodyText
    case secondaryText
    case captionText
    case buttonLabel
    case navigationTitle
    case tabLabel
    case inputLabel
    case errorText
}