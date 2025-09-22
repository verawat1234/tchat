//
//  AccessibilityService.swift
//  TchatApp
//
//  Comprehensive accessibility service providing WCAG 2.1 AA compliance
//  VoiceOver support, Dynamic Type, and accessibility testing capabilities
//

import SwiftUI
#if canImport(UIKit)
import UIKit
#endif
import Foundation
import Combine

/// Main accessibility service for iOS platform
/// Provides comprehensive accessibility features and WCAG 2.1 AA compliance
@MainActor
public class AccessibilityService: ObservableObject {

    // MARK: - Published Properties
    @Published public var isVoiceOverRunning: Bool = false
    @Published public var preferredContentSizeCategory: ContentSizeCategory = .medium
    @Published public var isReduceMotionEnabled: Bool = false
    @Published public var isDarkerSystemColorsEnabled: Bool = false
    @Published public var isHighContrastEnabled: Bool = false
    @Published public var isClosedCaptioningEnabled: Bool = false
    @Published public var accessibilityAuditResults: [AccessibilityAuditResult] = []

    // MARK: - Private Properties
    private var cancellables = Set<AnyCancellable>()
    private let notificationCenter = NotificationCenter.default
    private var auditTimer: Timer?

    // MARK: - Singleton Instance
    public static let shared = AccessibilityService()

    // MARK: - Initialization
    private init() {
        setupAccessibilityObservers()
        updateAccessibilitySettings()
    }

    // MARK: - Public Interface

    /// Configure accessibility for a SwiftUI view
    public func configureView<T: View>(_ view: T,
                                     label: String? = nil,
                                     hint: String? = nil,
                                     traits: AccessibilityTraits = [],
                                     value: String? = nil,
                                     isButton: Bool = false,
                                     sortPriority: Double? = nil) -> some View {
        view
            .accessibilityLabel(label ?? "")
            .accessibilityHint(hint ?? "")
            .accessibilityValue(value ?? "")
            .accessibilityAddTraits(traits)
            .accessibilityRemoveTraits(isButton ? [] : .isButton)
            .accessibilityElement(children: .ignore)
            .accessibilityAction(.default) {
                if isButton {
                    // Handle button action accessibility
                    announceButtonActivation(label ?? "Button")
                }
            }
            .accessibilitySortPriority(sortPriority ?? 0)
    }

    /// Configure touch targets to meet minimum 44x44pt requirement
    public func ensureMinimumTouchTarget<T: View>(_ view: T) -> some View {
        view
            .frame(minWidth: 44, minHeight: 44)
            .contentShape(Rectangle())
    }

    /// Apply dynamic type scaling with proper bounds
    public func scaledFont(size: CGFloat,
                          weight: Font.Weight = .regular,
                          design: Font.Design = .default,
                          maxScale: CGFloat = 2.0) -> Font {
        let baseSize = min(size * fontScale, size * maxScale)
        return .system(size: baseSize, weight: weight, design: design)
    }

    /// Get contrast-compliant colors
    public func contrastCompliantColor(foreground: Color,
                                     background: Color,
                                     level: ContrastLevel = .aa) -> Color {
        let ratio = calculateContrastRatio(foreground: foreground, background: background)

        if ratio >= level.minimumRatio {
            return foreground
        } else {
            return adjustColorForContrast(foreground, against: background, targetRatio: level.minimumRatio)
        }
    }

    /// Announce important changes to screen readers
    public func announce(_ message: String, priority: UIAccessibility.NotificationPriority = .medium) {
        let notification: UIAccessibility.Notification = priority == .high ? .announcement : .layoutChanged
        UIAccessibility.post(notification: notification, argument: message)
    }

    /// Perform comprehensive accessibility audit
    public func performAccessibilityAudit(on view: UIView? = nil) {
        let targetView = view ?? UIApplication.shared.windows.first?.rootViewController?.view
        guard let viewToAudit = targetView else { return }

        Task {
            let results = await auditView(viewToAudit)
            await MainActor.run {
                self.accessibilityAuditResults = results
                logAuditResults(results)
            }
        }
    }

    /// Start continuous accessibility monitoring
    public func startAccessibilityMonitoring(interval: TimeInterval = 5.0) {
        auditTimer?.invalidate()
        auditTimer = Timer.scheduledTimer(withTimeInterval: interval, repeats: true) { _ in
            Task { @MainActor in
                self.performAccessibilityAudit()
            }
        }
    }

    /// Stop accessibility monitoring
    public func stopAccessibilityMonitoring() {
        auditTimer?.invalidate()
        auditTimer = nil
    }
}

// MARK: - Private Implementation
private extension AccessibilityService {

    /// Setup observers for accessibility settings changes
    func setupAccessibilityObservers() {
        // VoiceOver status changes
        notificationCenter.publisher(for: UIAccessibility.voiceOverStatusDidChangeNotification)
            .sink { [weak self] _ in
                self?.updateAccessibilitySettings()
            }
            .store(in: &cancellables)

        // Dynamic Type changes
        notificationCenter.publisher(for: UIContentSizeCategory.didChangeNotification)
            .sink { [weak self] _ in
                self?.updateAccessibilitySettings()
            }
            .store(in: &cancellables)

        // Reduce Motion changes
        notificationCenter.publisher(for: UIAccessibility.reduceMotionStatusDidChangeNotification)
            .sink { [weak self] _ in
                self?.updateAccessibilitySettings()
            }
            .store(in: &cancellables)

        // High Contrast changes
        notificationCenter.publisher(for: UIAccessibility.darkerSystemColorsStatusDidChangeNotification)
            .sink { [weak self] _ in
                self?.updateAccessibilitySettings()
            }
            .store(in: &cancellables)

        // Closed Captioning changes
        notificationCenter.publisher(for: UIAccessibility.closedCaptioningStatusDidChangeNotification)
            .sink { [weak self] _ in
                self?.updateAccessibilitySettings()
            }
            .store(in: &cancellables)
    }

    /// Update accessibility settings from system
    func updateAccessibilitySettings() {
        isVoiceOverRunning = UIAccessibility.isVoiceOverRunning
        preferredContentSizeCategory = ContentSizeCategory(UIApplication.shared.preferredContentSizeCategory)
        isReduceMotionEnabled = UIAccessibility.isReduceMotionEnabled
        isDarkerSystemColorsEnabled = UIAccessibility.isDarkerSystemColorsEnabled
        isHighContrastEnabled = UIAccessibility.isDarkerSystemColorsEnabled
        isClosedCaptioningEnabled = UIAccessibility.isClosedCaptioningEnabled
    }

    /// Calculate font scale based on Dynamic Type
    var fontScale: CGFloat {
        switch preferredContentSizeCategory {
        case .extraSmall: return 0.8
        case .small: return 0.9
        case .medium: return 1.0
        case .large: return 1.1
        case .extraLarge: return 1.2
        case .extraExtraLarge: return 1.3
        case .extraExtraExtraLarge: return 1.4
        case .accessibilityMedium: return 1.6
        case .accessibilityLarge: return 1.8
        case .accessibilityExtraLarge: return 2.0
        case .accessibilityExtraExtraLarge: return 2.2
        case .accessibilityExtraExtraExtraLarge: return 2.4
        default: return 1.0
        }
    }

    /// Announce button activation for screen readers
    func announceButtonActivation(_ buttonName: String) {
        if isVoiceOverRunning {
            announce("\(buttonName) activated", priority: .medium)
        }
    }

    /// Calculate contrast ratio between two colors
    func calculateContrastRatio(foreground: Color, background: Color) -> Double {
        let fgLuminance = calculateLuminance(color: foreground)
        let bgLuminance = calculateLuminance(color: background)

        let lighter = max(fgLuminance, bgLuminance)
        let darker = min(fgLuminance, bgLuminance)

        return (lighter + 0.05) / (darker + 0.05)
    }

    /// Calculate luminance of a color
    func calculateLuminance(color: Color) -> Double {
        let uiColor = UIColor(color)
        var red: CGFloat = 0
        var green: CGFloat = 0
        var blue: CGFloat = 0
        var alpha: CGFloat = 0

        uiColor.getRed(&red, green: &green, blue: &blue, alpha: &alpha)

        let sRGB = [red, green, blue].map { component in
            if component <= 0.03928 {
                return component / 12.92
            } else {
                return pow((component + 0.055) / 1.055, 2.4)
            }
        }

        return 0.2126 * sRGB[0] + 0.7152 * sRGB[1] + 0.0722 * sRGB[2]
    }

    /// Adjust color to meet contrast requirements
    func adjustColorForContrast(_ color: Color, against background: Color, targetRatio: Double) -> Color {
        let uiColor = UIColor(color)
        var hue: CGFloat = 0
        var saturation: CGFloat = 0
        var brightness: CGFloat = 0
        var alpha: CGFloat = 0

        uiColor.getHue(&hue, saturation: &saturation, brightness: &brightness, alpha: &alpha)

        // Adjust brightness to meet contrast ratio
        var adjustedBrightness = brightness
        var iterations = 0
        let maxIterations = 20

        while iterations < maxIterations {
            let testColor = Color(hue: Double(hue),
                                saturation: Double(saturation),
                                brightness: Double(adjustedBrightness))
            let ratio = calculateContrastRatio(foreground: testColor, background: background)

            if ratio >= targetRatio {
                break
            }

            // Adjust brightness toward better contrast
            let bgLuminance = calculateLuminance(color: background)
            if bgLuminance > 0.5 {
                adjustedBrightness = max(0, adjustedBrightness - 0.05)
            } else {
                adjustedBrightness = min(1, adjustedBrightness + 0.05)
            }

            iterations += 1
        }

        return Color(hue: Double(hue),
                    saturation: Double(saturation),
                    brightness: Double(adjustedBrightness))
    }

    /// Perform accessibility audit on a view
    func auditView(_ view: UIView) async -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        // Check touch target sizes
        results.append(contentsOf: auditTouchTargets(view))

        // Check accessibility labels
        results.append(contentsOf: auditAccessibilityLabels(view))

        // Check contrast ratios
        results.append(contentsOf: auditContrastRatios(view))

        // Check reading order
        results.append(contentsOf: auditReadingOrder(view))

        // Check accessibility traits
        results.append(contentsOf: auditAccessibilityTraits(view))

        return results
    }

    /// Audit touch target sizes
    func auditTouchTargets(_ view: UIView) -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        func checkView(_ v: UIView) {
            if v.isUserInteractionEnabled && !v.isHidden {
                let frame = v.frame
                if frame.width < 44 || frame.height < 44 {
                    results.append(AccessibilityAuditResult(
                        type: .touchTargetSize,
                        severity: .error,
                        message: "Touch target too small: \(frame.width)×\(frame.height)pt. Minimum required: 44×44pt",
                        element: v,
                        wcagCriterion: "2.5.5"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Audit accessibility labels
    func auditAccessibilityLabels(_ view: UIView) -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        func checkView(_ v: UIView) {
            if v.isAccessibilityElement {
                if let label = v.accessibilityLabel, !label.isEmpty {
                    // Check for good label practices
                    if label.count > 100 {
                        results.append(AccessibilityAuditResult(
                            type: .accessibilityLabel,
                            severity: .warning,
                            message: "Accessibility label too long (\(label.count) characters). Keep under 100 characters.",
                            element: v,
                            wcagCriterion: "1.3.1"
                        ))
                    }
                } else {
                    results.append(AccessibilityAuditResult(
                        type: .accessibilityLabel,
                        severity: .error,
                        message: "Missing accessibility label for interactive element",
                        element: v,
                        wcagCriterion: "1.3.1"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Audit contrast ratios
    func auditContrastRatios(_ view: UIView) -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        func checkView(_ v: UIView) {
            if let label = v as? UILabel {
                let textColor = Color(label.textColor)
                let backgroundColor = Color(label.backgroundColor ?? UIColor.clear)
                let ratio = calculateContrastRatio(foreground: textColor, background: backgroundColor)

                if ratio < 4.5 {
                    results.append(AccessibilityAuditResult(
                        type: .colorContrast,
                        severity: .error,
                        message: "Insufficient color contrast ratio: \(String(format: "%.2f", ratio)):1. Minimum required: 4.5:1",
                        element: v,
                        wcagCriterion: "1.4.3"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Audit reading order
    func auditReadingOrder(_ view: UIView) -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        let accessibleElements = view.subviews.filter { $0.isAccessibilityElement }

        for (index, element) in accessibleElements.enumerated() {
            if index > 0 {
                let previousElement = accessibleElements[index - 1]
                let currentFrame = element.frame
                let previousFrame = previousElement.frame

                // Check logical reading order (left-to-right, top-to-bottom)
                if currentFrame.minY < previousFrame.maxY && currentFrame.minX < previousFrame.minX {
                    results.append(AccessibilityAuditResult(
                        type: .readingOrder,
                        severity: .warning,
                        message: "Potential reading order issue: element may be read out of visual order",
                        element: element,
                        wcagCriterion: "1.3.2"
                    ))
                }
            }
        }

        return results
    }

    /// Audit accessibility traits
    func auditAccessibilityTraits(_ view: UIView) -> [AccessibilityAuditResult] {
        var results: [AccessibilityAuditResult] = []

        func checkView(_ v: UIView) {
            if v.isAccessibilityElement {
                let traits = v.accessibilityTraits

                // Check for appropriate button traits
                if v is UIButton && !traits.contains(.button) {
                    results.append(AccessibilityAuditResult(
                        type: .accessibilityTraits,
                        severity: .warning,
                        message: "Button element missing .button trait",
                        element: v,
                        wcagCriterion: "4.1.2"
                    ))
                }

                // Check for conflicting traits
                if traits.contains(.button) && traits.contains(.staticText) {
                    results.append(AccessibilityAuditResult(
                        type: .accessibilityTraits,
                        severity: .error,
                        message: "Conflicting accessibility traits: button and staticText",
                        element: v,
                        wcagCriterion: "4.1.2"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return results
    }

    /// Log audit results for debugging
    func logAuditResults(_ results: [AccessibilityAuditResult]) {
        let errorCount = results.filter { $0.severity == .error }.count
        let warningCount = results.filter { $0.severity == .warning }.count

        print("🔍 Accessibility Audit Complete")
        print("📊 Found \(errorCount) errors and \(warningCount) warnings")

        for result in results {
            let icon = result.severity == .error ? "❌" : "⚠️"
            print("\(icon) [\(result.wcagCriterion)] \(result.message)")
        }
    }
}

// MARK: - Supporting Types

/// Accessibility audit result
public struct AccessibilityAuditResult {
    public let type: AuditType
    public let severity: Severity
    public let message: String
    public let element: UIView
    public let wcagCriterion: String

    public enum AuditType {
        case touchTargetSize
        case accessibilityLabel
        case colorContrast
        case readingOrder
        case accessibilityTraits
    }

    public enum Severity {
        case error
        case warning
    }
}

/// Contrast level requirements
public enum ContrastLevel {
    case aa
    case aaa

    var minimumRatio: Double {
        switch self {
        case .aa: return 4.5
        case .aaa: return 7.0
        }
    }
}

// MARK: - ContentSizeCategory Extension
extension ContentSizeCategory {
    init(_ uiContentSizeCategory: UIContentSizeCategory) {
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

// MARK: - View Extensions for Accessibility
public extension View {

    /// Apply accessibility configuration using the service
    func accessibilityConfigured(
        label: String? = nil,
        hint: String? = nil,
        traits: AccessibilityTraits = [],
        value: String? = nil,
        isButton: Bool = false,
        sortPriority: Double? = nil
    ) -> some View {
        AccessibilityService.shared.configureView(
            self,
            label: label,
            hint: hint,
            traits: traits,
            value: value,
            isButton: isButton,
            sortPriority: sortPriority
        )
    }

    /// Ensure minimum touch target size
    func minimumTouchTarget() -> some View {
        AccessibilityService.shared.ensureMinimumTouchTarget(self)
    }

    /// Apply Dynamic Type scaling with bounds
    func scaledFont(size: CGFloat, weight: Font.Weight = .regular, maxScale: CGFloat = 2.0) -> some View {
        self.font(AccessibilityService.shared.scaledFont(size: size, weight: weight, maxScale: maxScale))
    }

    /// Hide from accessibility when appropriate
    func accessibilityHidden(_ hidden: Bool = true) -> some View {
        self.accessibilityHidden(hidden)
    }
}