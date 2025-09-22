//
//  AccessibilityValidator.swift
//  TchatApp
//
//  Runtime accessibility validation tool for iOS
//  Provides real-time accessibility feedback during development
//

import SwiftUI
#if canImport(UIKit)
import UIKit
#endif
import Foundation

/// Runtime accessibility validator for iOS
/// Provides real-time validation and feedback during development
@MainActor
public class AccessibilityValidator: ObservableObject {

    // MARK: - Published Properties
    @Published public var isValidationEnabled: Bool = false
    @Published public var validationResults: [ValidationResult] = []
    @Published public var showValidationOverlay: Bool = false

    // MARK: - Private Properties
    private let accessibilityService = AccessibilityService.shared
    private var validationTimer: Timer?
    private var validationOverlay: UIView?

    // MARK: - Singleton Instance
    public static let shared = AccessibilityValidator()

    // MARK: - Initialization
    private init() {}

    // MARK: - Public Interface

    /// Start real-time accessibility validation
    public func startRealTimeValidation(interval: TimeInterval = 2.0) {
        isValidationEnabled = true

        validationTimer = Timer.scheduledTimer(withTimeInterval: interval, repeats: true) { _ in
            Task { @MainActor in
                await self.validateCurrentScreen()
            }
        }

        if showValidationOverlay {
            createValidationOverlay()
        }
    }

    /// Stop real-time validation
    public func stopRealTimeValidation() {
        isValidationEnabled = false
        validationTimer?.invalidate()
        validationTimer = nil
        removeValidationOverlay()
    }

    /// Validate a specific view
    public func validateView(_ view: UIView) async -> [ValidationResult] {
        return await performViewValidation(view)
    }

    /// Toggle validation overlay
    public func toggleValidationOverlay() {
        showValidationOverlay.toggle()

        if showValidationOverlay && isValidationEnabled {
            createValidationOverlay()
        } else {
            removeValidationOverlay()
        }
    }

    /// Validate specific accessibility criteria
    public func validateCriteria(_ criteria: ValidationCriteria, for view: UIView) -> ValidationResult {
        switch criteria {
        case .touchTargetSize:
            return validateTouchTargetSize(view)
        case .accessibilityLabel:
            return validateAccessibilityLabel(view)
        case .colorContrast:
            return validateColorContrast(view)
        case .readingOrder:
            return validateReadingOrder(view)
        case .voiceOverCompatibility:
            return validateVoiceOverCompatibility(view)
        }
    }

    /// Get validation summary
    public func getValidationSummary() -> ValidationSummary {
        let total = validationResults.count
        let passed = validationResults.filter { $0.status == .passed }.count
        let warnings = validationResults.filter { $0.status == .warning }.count
        let failed = validationResults.filter { $0.status == .failed }.count

        return ValidationSummary(
            total: total,
            passed: passed,
            warnings: warnings,
            failed: failed,
            score: total > 0 ? Double(passed) / Double(total) : 0.0
        )
    }
}

// MARK: - Private Implementation
private extension AccessibilityValidator {

    /// Validate current screen
    func validateCurrentScreen() async {
        guard let rootView = UIApplication.shared.windows.first?.rootViewController?.view else {
            return
        }

        let results = await performViewValidation(rootView)
        validationResults = results

        if showValidationOverlay {
            updateValidationOverlay()
        }
    }

    /// Perform comprehensive view validation
    func performViewValidation(_ view: UIView) async -> [ValidationResult] {
        var results: [ValidationResult] = []

        // Validate touch targets
        results.append(contentsOf: validateAllTouchTargets(view))

        // Validate accessibility labels
        results.append(contentsOf: validateAllAccessibilityLabels(view))

        // Validate color contrast
        results.append(contentsOf: validateAllColorContrast(view))

        // Validate reading order
        results.append(contentsOf: validateAllReadingOrder(view))

        // Validate VoiceOver compatibility
        results.append(contentsOf: validateAllVoiceOverCompatibility(view))

        return results
    }

    /// Validate touch target size
    func validateTouchTargetSize(_ view: UIView) -> ValidationResult {
        let minSize: CGFloat = 44
        let frame = view.frame

        if view.isUserInteractionEnabled && !view.isHidden {
            if frame.width < minSize || frame.height < minSize {
                return ValidationResult(
                    element: view,
                    criteria: .touchTargetSize,
                    status: .failed,
                    message: "Touch target too small: \(Int(frame.width))×\(Int(frame.height))pt. Minimum: 44×44pt",
                    wcagCriterion: "2.5.5",
                    suggestion: "Increase the touch target size to at least 44×44pt by adding padding or increasing the frame size."
                )
            }
        }

        return ValidationResult(
            element: view,
            criteria: .touchTargetSize,
            status: .passed,
            message: "Touch target size is adequate",
            wcagCriterion: "2.5.5"
        )
    }

    /// Validate accessibility label
    func validateAccessibilityLabel(_ view: UIView) -> ValidationResult {
        if view.isAccessibilityElement {
            if let label = view.accessibilityLabel, !label.isEmpty {
                if label.count > 100 {
                    return ValidationResult(
                        element: view,
                        criteria: .accessibilityLabel,
                        status: .warning,
                        message: "Accessibility label too long: \(label.count) characters",
                        wcagCriterion: "1.3.1",
                        suggestion: "Keep accessibility labels concise and under 100 characters."
                    )
                }

                return ValidationResult(
                    element: view,
                    criteria: .accessibilityLabel,
                    status: .passed,
                    message: "Accessibility label is appropriate",
                    wcagCriterion: "1.3.1"
                )
            } else {
                return ValidationResult(
                    element: view,
                    criteria: .accessibilityLabel,
                    status: .failed,
                    message: "Missing accessibility label",
                    wcagCriterion: "1.3.1",
                    suggestion: "Add an accessibility label that describes the purpose of this element."
                )
            }
        }

        return ValidationResult(
            element: view,
            criteria: .accessibilityLabel,
            status: .passed,
            message: "Element not accessibility-enabled",
            wcagCriterion: "1.3.1"
        )
    }

    /// Validate color contrast
    func validateColorContrast(_ view: UIView) -> ValidationResult {
        guard let label = view as? UILabel else {
            return ValidationResult(
                element: view,
                criteria: .colorContrast,
                status: .passed,
                message: "Not a text element",
                wcagCriterion: "1.4.3"
            )
        }

        let textColor = Color(label.textColor)
        let backgroundColor = Color(label.backgroundColor ?? UIColor.systemBackground)
        let ratio = accessibilityService.calculateContrastRatio(foreground: textColor, background: backgroundColor)

        if ratio < 4.5 {
            return ValidationResult(
                element: view,
                criteria: .colorContrast,
                status: .failed,
                message: "Insufficient contrast ratio: \(String(format: "%.2f", ratio)):1",
                wcagCriterion: "1.4.3",
                suggestion: "Adjust text or background color to achieve at least 4.5:1 contrast ratio."
            )
        } else if ratio < 7.0 {
            return ValidationResult(
                element: view,
                criteria: .colorContrast,
                status: .warning,
                message: "Good contrast ratio: \(String(format: "%.2f", ratio)):1. Consider 7:1 for AAA compliance",
                wcagCriterion: "1.4.3"
            )
        }

        return ValidationResult(
            element: view,
            criteria: .colorContrast,
            status: .passed,
            message: "Excellent contrast ratio: \(String(format: "%.2f", ratio)):1",
            wcagCriterion: "1.4.3"
        )
    }

    /// Validate reading order
    func validateReadingOrder(_ view: UIView) -> ValidationResult {
        // This is a simplified validation - full implementation would analyze entire view hierarchy
        return ValidationResult(
            element: view,
            criteria: .readingOrder,
            status: .passed,
            message: "Reading order validation requires full hierarchy analysis",
            wcagCriterion: "1.3.2"
        )
    }

    /// Validate VoiceOver compatibility
    func validateVoiceOverCompatibility(_ view: UIView) -> ValidationResult {
        if view.isAccessibilityElement {
            let traits = view.accessibilityTraits

            // Check for appropriate traits
            if view is UIButton && !traits.contains(.button) {
                return ValidationResult(
                    element: view,
                    criteria: .voiceOverCompatibility,
                    status: .warning,
                    message: "Button missing .button trait",
                    wcagCriterion: "4.1.2",
                    suggestion: "Add .button accessibility trait to button elements."
                )
            }

            // Check for conflicting traits
            if traits.contains(.button) && traits.contains(.staticText) {
                return ValidationResult(
                    element: view,
                    criteria: .voiceOverCompatibility,
                    status: .failed,
                    message: "Conflicting accessibility traits",
                    wcagCriterion: "4.1.2",
                    suggestion: "Remove conflicting accessibility traits."
                )
            }
        }

        return ValidationResult(
            element: view,
            criteria: .voiceOverCompatibility,
            status: .passed,
            message: "VoiceOver compatibility is good",
            wcagCriterion: "4.1.2"
        )
    }

    /// Validate all touch targets in view hierarchy
    func validateAllTouchTargets(_ view: UIView) -> [ValidationResult] {
        var results: [ValidationResult] = []

        func checkView(_ v: UIView) {
            if v.isUserInteractionEnabled {
                results.append(validateTouchTargetSize(v))
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Validate all accessibility labels in view hierarchy
    func validateAllAccessibilityLabels(_ view: UIView) -> [ValidationResult] {
        var results: [ValidationResult] = []

        func checkView(_ v: UIView) {
            if v.isAccessibilityElement {
                results.append(validateAccessibilityLabel(v))
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Validate all color contrast in view hierarchy
    func validateAllColorContrast(_ view: UIView) -> [ValidationResult] {
        var results: [ValidationResult] = []

        func checkView(_ v: UIView) {
            if v is UILabel {
                results.append(validateColorContrast(v))
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Validate all reading order in view hierarchy
    func validateAllReadingOrder(_ view: UIView) -> [ValidationResult] {
        // Simplified implementation
        return [validateReadingOrder(view)]
    }

    /// Validate all VoiceOver compatibility in view hierarchy
    func validateAllVoiceOverCompatibility(_ view: UIView) -> [ValidationResult] {
        var results: [ValidationResult] = []

        func checkView(_ v: UIView) {
            if v.isAccessibilityElement {
                results.append(validateVoiceOverCompatibility(v))
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Create validation overlay
    func createValidationOverlay() {
        guard let window = UIApplication.shared.windows.first else { return }

        removeValidationOverlay()

        let overlay = UIView(frame: window.bounds)
        overlay.backgroundColor = UIColor.clear
        overlay.isUserInteractionEnabled = false

        // Add validation indicators
        updateValidationOverlay(overlay)

        window.addSubview(overlay)
        validationOverlay = overlay
    }

    /// Update validation overlay
    func updateValidationOverlay(_ overlay: UIView? = nil) {
        let targetOverlay = overlay ?? validationOverlay
        guard let overlay = targetOverlay else { return }

        // Clear existing indicators
        overlay.subviews.forEach { $0.removeFromSuperview() }

        // Add new indicators based on validation results
        for result in validationResults {
            let indicator = createValidationIndicator(for: result)
            overlay.addSubview(indicator)
        }
    }

    /// Create validation indicator for a result
    func createValidationIndicator(for result: ValidationResult) -> UIView {
        let frame = result.element.convert(result.element.bounds, to: nil)
        let indicator = UIView(frame: frame)

        switch result.status {
        case .passed:
            indicator.layer.borderColor = UIColor.systemGreen.cgColor
            indicator.backgroundColor = UIColor.systemGreen.withAlphaComponent(0.1)
        case .warning:
            indicator.layer.borderColor = UIColor.systemOrange.cgColor
            indicator.backgroundColor = UIColor.systemOrange.withAlphaComponent(0.2)
        case .failed:
            indicator.layer.borderColor = UIColor.systemRed.cgColor
            indicator.backgroundColor = UIColor.systemRed.withAlphaComponent(0.2)
        }

        indicator.layer.borderWidth = 2
        indicator.isUserInteractionEnabled = false

        return indicator
    }

    /// Remove validation overlay
    func removeValidationOverlay() {
        validationOverlay?.removeFromSuperview()
        validationOverlay = nil
    }
}

// MARK: - Supporting Types

/// Validation criteria
public enum ValidationCriteria {
    case touchTargetSize
    case accessibilityLabel
    case colorContrast
    case readingOrder
    case voiceOverCompatibility
}

/// Validation result
public struct ValidationResult {
    public let element: UIView
    public let criteria: ValidationCriteria
    public let status: ValidationStatus
    public let message: String
    public let wcagCriterion: String
    public let suggestion: String?

    public init(element: UIView, criteria: ValidationCriteria, status: ValidationStatus, message: String, wcagCriterion: String, suggestion: String? = nil) {
        self.element = element
        self.criteria = criteria
        self.status = status
        self.message = message
        self.wcagCriterion = wcagCriterion
        self.suggestion = suggestion
    }
}

/// Validation status
public enum ValidationStatus {
    case passed
    case warning
    case failed
}

/// Validation summary
public struct ValidationSummary {
    public let total: Int
    public let passed: Int
    public let warnings: Int
    public let failed: Int
    public let score: Double
}

// MARK: - SwiftUI Integration

/// SwiftUI view modifier for accessibility validation
public struct AccessibilityValidationModifier: ViewModifier {
    let validator = AccessibilityValidator.shared

    public func body(content: Content) -> some View {
        content
            .onAppear {
                if validator.isValidationEnabled {
                    Task {
                        await validator.validateCurrentScreen()
                    }
                }
            }
    }
}

public extension View {
    /// Enable accessibility validation for this view
    func accessibilityValidation() -> some View {
        modifier(AccessibilityValidationModifier())
    }
}