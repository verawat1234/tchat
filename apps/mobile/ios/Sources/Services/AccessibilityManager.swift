//
//  AccessibilityManager.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
#if canImport(UIKit)
import UIKit
#endif
import SwiftUI

/// Accessibility compliance manager for iOS
/// Implements T054: Accessibility compliance for VoiceOver, Dynamic Type, and WCAG guidelines
public class AccessibilityManager: ObservableObject {

    // MARK: - Published Properties

    @Published public var isVoiceOverRunning: Bool = false
    @Published public var preferredContentSizeCategory: ContentSizeCategory = .medium
    @Published public var isReduceMotionEnabled: Bool = false
    @Published public var isHighContrastEnabled: Bool = false
    @Published public var isDifferentiateWithoutColorEnabled: Bool = false

    // MARK: - Private Properties

    private var notificationCenter = NotificationCenter.default
    private var observers: [NSObjectProtocol] = []

    // MARK: - Accessibility Compliance Levels

    public enum ComplianceLevel: String, CaseIterable {
        case a = "A"
        case aa = "AA"
        case aaa = "AAA"

        var description: String {
            switch self {
            case .a: return "WCAG 2.1 Level A (Basic)"
            case .aa: return "WCAG 2.1 Level AA (Standard)"
            case .aaa: return "WCAG 2.1 Level AAA (Enhanced)"
            }
        }
    }

    // MARK: - Color Contrast Standards

    public struct ColorContrastRequirements {
        public static let normalTextAA: Double = 4.5
        public static let largeTextAA: Double = 3.0
        public static let normalTextAAA: Double = 7.0
        public static let largeTextAAA: Double = 4.5
        public static let uiComponentsAA: Double = 3.0
    }

    // MARK: - Initialization

    public static let shared = AccessibilityManager()

    private init() {
        setupNotificationObservers()
        updateAccessibilitySettings()
    }

    deinit {
        removeNotificationObservers()
    }

    // MARK: - Public Interface

    /// Starts monitoring accessibility changes
    public func startMonitoring() {
        updateAccessibilitySettings()
    }

    /// Stops monitoring accessibility changes
    public func stopMonitoring() {
        removeNotificationObservers()
    }

    /// Validates color contrast ratio according to WCAG guidelines
    public func validateColorContrast(
        foreground: UIColor,
        background: UIColor,
        fontSize: CGFloat,
        complianceLevel: ComplianceLevel = .aa
    ) -> AccessibilityValidationResult {

        let contrastRatio = calculateContrastRatio(foreground: foreground, background: background)
        let isLargeText = fontSize >= 18.0 || (fontSize >= 14.0 && isBoldText(fontSize))

        let requiredRatio: Double
        switch complianceLevel {
        case .a:
            requiredRatio = 3.0 // Basic requirement
        case .aa:
            requiredRatio = isLargeText ? ColorContrastRequirements.largeTextAA : ColorContrastRequirements.normalTextAA
        case .aaa:
            requiredRatio = isLargeText ? ColorContrastRequirements.largeTextAAA : ColorContrastRequirements.normalTextAAA
        }

        let passes = contrastRatio >= requiredRatio

        return AccessibilityValidationResult(
            passes: passes,
            actualRatio: contrastRatio,
            requiredRatio: requiredRatio,
            complianceLevel: complianceLevel,
            message: passes ? "Color contrast meets \(complianceLevel.description)" :
                    "Color contrast (\(String(format: "%.2f", contrastRatio)):1) below required \(String(format: "%.1f", requiredRatio)):1 for \(complianceLevel.description)"
        )
    }

    /// Configures accessibility for UI elements
    public func configureAccessibility(
        for view: UIView,
        label: String,
        hint: String? = nil,
        traits: UIAccessibilityTraits = [],
        value: String? = nil
    ) {
        view.isAccessibilityElement = true
        view.accessibilityLabel = label
        view.accessibilityHint = hint
        view.accessibilityTraits = traits
        view.accessibilityValue = value
    }

    /// Configures accessibility for SwiftUI views
    public func accessibilityModifier(
        label: String,
        hint: String? = nil,
        value: String? = nil,
        traits: AccessibilityTraits = []
    ) -> some ViewModifier {
        return AccessibilityConfigurationModifier(
            label: label,
            hint: hint,
            value: value,
            traits: traits
        )
    }

    /// Announces accessibility changes to VoiceOver
    public func announceChange(_ message: String, priority: UIAccessibility.Notification = .announcement) {
        if isVoiceOverRunning {
            UIAccessibility.post(notification: priority, argument: message)
        }
    }

    /// Gets optimal font size for current Dynamic Type setting
    public func scaledFont(baseSize: CGFloat, textStyle: UIFont.TextStyle = .body) -> UIFont {
        let font = UIFont.systemFont(ofSize: baseSize)
        return UIFontMetrics(forTextStyle: textStyle).scaledFont(for: font)
    }

    /// Checks if device settings require reduced motion
    public var shouldReduceMotion: Bool {
        return isReduceMotionEnabled
    }

    /// Gets appropriate animation duration based on accessibility settings
    public func animationDuration(default duration: TimeInterval) -> TimeInterval {
        return shouldReduceMotion ? 0.0 : duration
    }

    /// Validates component accessibility compliance
    public func validateComponentAccessibility(
        component: UIView,
        requirements: AccessibilityRequirements = .standard
    ) -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        // Check for accessibility label
        if component.isAccessibilityElement && component.accessibilityLabel?.isEmpty != false {
            issues.append(.missingLabel("Component missing accessibility label"))
        }

        // Check minimum touch target size (44x44 points for AA compliance)
        if component.frame.width < 44 || component.frame.height < 44 {
            issues.append(.insufficientTouchTarget("Touch target smaller than 44x44 points"))
        }

        // Check for color contrast if applicable
        if let backgroundColor = component.backgroundColor,
           let textColor = getTextColor(from: component) {
            let validation = validateColorContrast(
                foreground: textColor,
                background: backgroundColor,
                fontSize: getFontSize(from: component),
                complianceLevel: requirements.colorContrastLevel
            )

            if !validation.passes {
                issues.append(.insufficientContrast(validation.message))
            }
        }

        return issues
    }

    // MARK: - Private Methods

    private func setupNotificationObservers() {
        // VoiceOver status changes
        observers.append(
            notificationCenter.addObserver(
                forName: UIAccessibility.voiceOverStatusDidChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                self?.updateVoiceOverStatus()
            }
        )

        // Dynamic Type changes
        observers.append(
            notificationCenter.addObserver(
                forName: UIContentSizeCategory.didChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                self?.updateContentSizeCategory()
            }
        )

        // Reduce Motion changes
        observers.append(
            notificationCenter.addObserver(
                forName: UIAccessibility.reduceMotionStatusDidChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                self?.updateReduceMotionStatus()
            }
        )

        // High Contrast changes
        observers.append(
            notificationCenter.addObserver(
                forName: UIAccessibility.darkerSystemColorsStatusDidChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                self?.updateHighContrastStatus()
            }
        )

        // Differentiate Without Color changes
        observers.append(
            notificationCenter.addObserver(
                forName: UIAccessibility.differentiateWithoutColorDidChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                self?.updateDifferentiateWithoutColorStatus()
            }
        )
    }

    private func removeNotificationObservers() {
        observers.forEach { notificationCenter.removeObserver($0) }
        observers.removeAll()
    }

    private func updateAccessibilitySettings() {
        updateVoiceOverStatus()
        updateContentSizeCategory()
        updateReduceMotionStatus()
        updateHighContrastStatus()
        updateDifferentiateWithoutColorStatus()
    }

    private func updateVoiceOverStatus() {
        isVoiceOverRunning = UIAccessibility.isVoiceOverRunning
    }

    private func updateContentSizeCategory() {
        let category = UIApplication.shared.preferredContentSizeCategory
        preferredContentSizeCategory = ContentSizeCategory(from: category)
    }

    private func updateReduceMotionStatus() {
        isReduceMotionEnabled = UIAccessibility.isReduceMotionEnabled
    }

    private func updateHighContrastStatus() {
        isHighContrastEnabled = UIAccessibility.isDarkerSystemColorsEnabled
    }

    private func updateDifferentiateWithoutColorStatus() {
        isDifferentiateWithoutColorEnabled = UIAccessibility.shouldDifferentiateWithoutColor
    }

    // MARK: - Color Contrast Calculations

    private func calculateContrastRatio(foreground: UIColor, background: UIColor) -> Double {
        let foregroundLuminance = relativeLuminance(of: foreground)
        let backgroundLuminance = relativeLuminance(of: background)

        let lighter = max(foregroundLuminance, backgroundLuminance)
        let darker = min(foregroundLuminance, backgroundLuminance)

        return (lighter + 0.05) / (darker + 0.05)
    }

    private func relativeLuminance(of color: UIColor) -> Double {
        var red: CGFloat = 0
        var green: CGFloat = 0
        var blue: CGFloat = 0
        var alpha: CGFloat = 0

        color.getRed(&red, green: &green, blue: &blue, alpha: &alpha)

        // Convert to linear RGB
        let linearRed = linearizeColorComponent(Double(red))
        let linearGreen = linearizeColorComponent(Double(green))
        let linearBlue = linearizeColorComponent(Double(blue))

        // Calculate relative luminance using ITU-R BT.709 coefficients
        return 0.2126 * linearRed + 0.7152 * linearGreen + 0.0722 * linearBlue
    }

    private func linearizeColorComponent(_ component: Double) -> Double {
        if component <= 0.03928 {
            return component / 12.92
        } else {
            return pow((component + 0.055) / 1.055, 2.4)
        }
    }

    // MARK: - Helper Methods

    private func isBoldText(_ fontSize: CGFloat) -> Bool {
        // Simplified bold text detection
        return fontSize >= 14.0
    }

    private func getTextColor(from view: UIView) -> UIColor? {
        if let label = view as? UILabel {
            return label.textColor
        } else if let button = view as? UIButton {
            return button.titleColor(for: .normal)
        }
        return nil
    }

    private func getFontSize(from view: UIView) -> CGFloat {
        if let label = view as? UILabel {
            return label.font.pointSize
        } else if let button = view as? UIButton {
            return button.titleLabel?.font.pointSize ?? 17.0
        }
        return 17.0 // Default system font size
    }
}

// MARK: - Supporting Types

public struct AccessibilityValidationResult {
    public let passes: Bool
    public let actualRatio: Double
    public let requiredRatio: Double
    public let complianceLevel: AccessibilityManager.ComplianceLevel
    public let message: String
}

public struct AccessibilityRequirements {
    public let colorContrastLevel: AccessibilityManager.ComplianceLevel
    public let minimumTouchTargetSize: CGSize
    public let requiresLabel: Bool

    public static let standard = AccessibilityRequirements(
        colorContrastLevel: .aa,
        minimumTouchTargetSize: CGSize(width: 44, height: 44),
        requiresLabel: true
    )

    public static let enhanced = AccessibilityRequirements(
        colorContrastLevel: .aaa,
        minimumTouchTargetSize: CGSize(width: 48, height: 48),
        requiresLabel: true
    )
}

public enum AccessibilityIssue {
    case missingLabel(String)
    case insufficientTouchTarget(String)
    case insufficientContrast(String)
    case missingHint(String)
    case incorrectTraits(String)

    public var description: String {
        switch self {
        case .missingLabel(let message),
             .insufficientTouchTarget(let message),
             .insufficientContrast(let message),
             .missingHint(let message),
             .incorrectTraits(let message):
            return message
        }
    }
}

// MARK: - SwiftUI Extensions

public extension ContentSizeCategory {
    init(from uiContentSizeCategory: UIContentSizeCategory) {
        switch uiContentSizeCategory {
        case .extraSmall: self = .extraSmall
        case .small: self = .small
        case .medium: self = .medium
        case .large: self = .large
        case .extraLarge: self = .extraLarge
        case .extraExtraLarge: self = .extraExtraLarge
        case .extraExtraExtraLarge: self = .extraExtraExtraLarge
        case .accessibilityMedium: self = .accessibilityMedium
        case .accessibilityLarge: self = .accessibilityLarge
        case .accessibilityExtraLarge: self = .accessibilityExtraLarge
        case .accessibilityExtraExtraLarge: self = .accessibilityExtraExtraLarge
        case .accessibilityExtraExtraExtraLarge: self = .accessibilityExtraExtraExtraLarge
        default: self = .medium
        }
    }
}

struct AccessibilityConfigurationModifier: ViewModifier {
    let label: String
    let hint: String?
    let value: String?
    let traits: AccessibilityTraits

    func body(content: Content) -> some View {
        content
            .accessibilityLabel(label)
            .accessibilityHint(hint ?? "")
            .accessibilityValue(value ?? "")
            .accessibilityAddTraits(traits)
    }
}

// MARK: - View Extensions

public extension View {
    func accessibilityConfiguration(
        label: String,
        hint: String? = nil,
        value: String? = nil,
        traits: AccessibilityTraits = []
    ) -> some View {
        modifier(AccessibilityConfigurationModifier(
            label: label,
            hint: hint,
            value: value,
            traits: traits
        ))
    }

    func dynamicTypeSize(min: DynamicTypeSize = .xSmall, max: DynamicTypeSize = .accessibility5) -> some View {
        self.dynamicTypeSize(min...max)
    }
}