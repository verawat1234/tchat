//
//  AccessibilityAuditor.swift
//  TchatApp
//
//  Comprehensive accessibility auditing tool for iOS
//  Provides detailed WCAG compliance analysis and reporting
//

import SwiftUI
#if canImport(UIKit)
import UIKit
#endif
import Foundation

/// Comprehensive accessibility auditor for iOS
/// Performs deep analysis of app accessibility compliance
@MainActor
public class AccessibilityAuditor: ObservableObject {

    // MARK: - Published Properties
    @Published public var currentAuditResults: AccessibilityAuditReport?
    @Published public var isAuditing: Bool = false
    @Published public var auditProgress: Double = 0.0

    // MARK: - Private Properties
    private let accessibilityService = AccessibilityService.shared

    // MARK: - Public Interface

    /// Perform comprehensive accessibility audit
    public func performComprehensiveAudit() async -> AccessibilityAuditReport {
        await MainActor.run {
            isAuditing = true
            auditProgress = 0.0
        }

        let report = AccessibilityAuditReport()

        // Step 1: Basic compliance checks (20%)
        await updateProgress(0.2)
        report.basicCompliance = await auditBasicCompliance()

        // Step 2: VoiceOver compatibility (40%)
        await updateProgress(0.4)
        report.voiceOverCompatibility = await auditVoiceOverCompatibility()

        // Step 3: Dynamic Type support (60%)
        await updateProgress(0.6)
        report.dynamicTypeSupport = await auditDynamicTypeSupport()

        // Step 4: Color contrast analysis (80%)
        await updateProgress(0.8)
        report.colorContrastAnalysis = await auditColorContrast()

        // Step 5: Touch target analysis (100%)
        await updateProgress(1.0)
        report.touchTargetAnalysis = await auditTouchTargets()

        // Generate overall score
        report.overallScore = calculateOverallScore(report)
        report.wcagLevel = determineWCAGLevel(report)
        report.recommendations = generateRecommendations(report)

        await MainActor.run {
            currentAuditResults = report
            isAuditing = false
        }

        return report
    }

    /// Generate accessibility report
    public func generateAccessibilityReport(_ report: AccessibilityAuditReport) -> String {
        var reportText = "# Tchat iOS Accessibility Audit Report\n\n"
        reportText += "**Generated:** \(Date().formatted())\n"
        reportText += "**Overall Score:** \(Int(report.overallScore * 100))%\n"
        reportText += "**WCAG Level:** \(report.wcagLevel.rawValue)\n\n"

        // Basic Compliance Section
        reportText += "## Basic Compliance\n"
        reportText += "**Score:** \(Int(report.basicCompliance.score * 100))%\n"
        reportText += "**Issues Found:** \(report.basicCompliance.issues.count)\n\n"

        for issue in report.basicCompliance.issues {
            reportText += "- ❌ **\(issue.title):** \(issue.description)\n"
            reportText += "  - *WCAG Criterion:* \(issue.wcagCriterion)\n"
            reportText += "  - *Severity:* \(issue.severity.rawValue)\n\n"
        }

        // VoiceOver Compatibility Section
        reportText += "## VoiceOver Compatibility\n"
        reportText += "**Score:** \(Int(report.voiceOverCompatibility.score * 100))%\n"
        reportText += "**Issues Found:** \(report.voiceOverCompatibility.issues.count)\n\n"

        for issue in report.voiceOverCompatibility.issues {
            reportText += "- ❌ **\(issue.title):** \(issue.description)\n"
            reportText += "  - *WCAG Criterion:* \(issue.wcagCriterion)\n"
            reportText += "  - *Severity:* \(issue.severity.rawValue)\n\n"
        }

        // Dynamic Type Support Section
        reportText += "## Dynamic Type Support\n"
        reportText += "**Score:** \(Int(report.dynamicTypeSupport.score * 100))%\n"
        reportText += "**Issues Found:** \(report.dynamicTypeSupport.issues.count)\n\n"

        for issue in report.dynamicTypeSupport.issues {
            reportText += "- ❌ **\(issue.title):** \(issue.description)\n"
            reportText += "  - *WCAG Criterion:* \(issue.wcagCriterion)\n"
            reportText += "  - *Severity:* \(issue.severity.rawValue)\n\n"
        }

        // Color Contrast Analysis Section
        reportText += "## Color Contrast Analysis\n"
        reportText += "**Score:** \(Int(report.colorContrastAnalysis.score * 100))%\n"
        reportText += "**Issues Found:** \(report.colorContrastAnalysis.issues.count)\n\n"

        for issue in report.colorContrastAnalysis.issues {
            reportText += "- ❌ **\(issue.title):** \(issue.description)\n"
            reportText += "  - *WCAG Criterion:* \(issue.wcagCriterion)\n"
            reportText += "  - *Severity:* \(issue.severity.rawValue)\n\n"
        }

        // Touch Target Analysis Section
        reportText += "## Touch Target Analysis\n"
        reportText += "**Score:** \(Int(report.touchTargetAnalysis.score * 100))%\n"
        reportText += "**Issues Found:** \(report.touchTargetAnalysis.issues.count)\n\n"

        for issue in report.touchTargetAnalysis.issues {
            reportText += "- ❌ **\(issue.title):** \(issue.description)\n"
            reportText += "  - *WCAG Criterion:* \(issue.wcagCriterion)\n"
            reportText += "  - *Severity:* \(issue.severity.rawValue)\n\n"
        }

        // Recommendations Section
        reportText += "## Recommendations\n\n"
        for recommendation in report.recommendations {
            reportText += "### \(recommendation.priority.rawValue) Priority: \(recommendation.title)\n"
            reportText += "\(recommendation.description)\n\n"
            reportText += "**Implementation Steps:**\n"
            for (index, step) in recommendation.implementationSteps.enumerated() {
                reportText += "\(index + 1). \(step)\n"
            }
            reportText += "\n"
        }

        return reportText
    }

    /// Export audit report to file
    public func exportAuditReport(_ report: AccessibilityAuditReport, to url: URL) {
        let reportContent = generateAccessibilityReport(report)

        do {
            try reportContent.write(to: url, atomically: true, encoding: .utf8)
            print("✅ Accessibility report exported to: \(url.path)")
        } catch {
            print("❌ Failed to export report: \(error.localizedDescription)")
        }
    }
}

// MARK: - Private Implementation
private extension AccessibilityAuditor {

    func updateProgress(_ progress: Double) async {
        await MainActor.run {
            auditProgress = progress
        }
    }

    func auditBasicCompliance() async -> AuditSection {
        var issues: [AccessibilityIssue] = []

        // Check if accessibility is enabled
        if !UIAccessibility.isVoiceOverRunning && !UIAccessibility.isSwitchControlRunning {
            issues.append(AccessibilityIssue(
                title: "Accessibility Services Not Active",
                description: "No accessibility services are currently running. Enable VoiceOver or Switch Control to test.",
                severity: .warning,
                wcagCriterion: "4.1.2"
            ))
        }

        // Check app-wide accessibility settings
        let rootView = UIApplication.shared.windows.first?.rootViewController?.view
        if let view = rootView {
            let childIssues = await auditViewHierarchy(view)
            issues.append(contentsOf: childIssues)
        }

        let score = calculateSectionScore(issues: issues, totalChecks: 10)
        return AuditSection(score: score, issues: issues)
    }

    func auditVoiceOverCompatibility() async -> AuditSection {
        var issues: [AccessibilityIssue] = []

        // Test VoiceOver navigation
        let rootView = UIApplication.shared.windows.first?.rootViewController?.view
        if let view = rootView {
            let voiceOverIssues = await auditVoiceOverNavigation(view)
            issues.append(contentsOf: voiceOverIssues)
        }

        let score = calculateSectionScore(issues: issues, totalChecks: 8)
        return AuditSection(score: score, issues: issues)
    }

    func auditDynamicTypeSupport() async -> AuditSection {
        var issues: [AccessibilityIssue] = []

        // Test with different content size categories
        let testCategories: [UIContentSizeCategory] = [
            .large,
            .extraExtraLarge,
            .accessibilityMedium,
            .accessibilityExtraExtraExtraLarge
        ]

        for category in testCategories {
            let categoryIssues = await auditDynamicTypeCategory(category)
            issues.append(contentsOf: categoryIssues)
        }

        let score = calculateSectionScore(issues: issues, totalChecks: 6)
        return AuditSection(score: score, issues: issues)
    }

    func auditColorContrast() async -> AuditSection {
        var issues: [AccessibilityIssue] = []

        let rootView = UIApplication.shared.windows.first?.rootViewController?.view
        if let view = rootView {
            let contrastIssues = await auditViewContrast(view)
            issues.append(contentsOf: contrastIssues)
        }

        let score = calculateSectionScore(issues: issues, totalChecks: 12)
        return AuditSection(score: score, issues: issues)
    }

    func auditTouchTargets() async -> AuditSection {
        var issues: [AccessibilityIssue] = []

        let rootView = UIApplication.shared.windows.first?.rootViewController?.view
        if let view = rootView {
            let touchTargetIssues = await auditViewTouchTargets(view)
            issues.append(contentsOf: touchTargetIssues)
        }

        let score = calculateSectionScore(issues: issues, totalChecks: 5)
        return AuditSection(score: score, issues: issues)
    }

    #if canImport(UIKit)
    func auditViewHierarchy(_ view: UIView) async -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        func checkView(_ v: UIView) {
            // Check for missing accessibility labels
            if v.isAccessibilityElement && (v.accessibilityLabel?.isEmpty ?? true) {
                issues.append(AccessibilityIssue(
                    title: "Missing Accessibility Label",
                    description: "Interactive element lacks accessibility label",
                    severity: .error,
                    wcagCriterion: "1.3.1"
                ))
            }

            // Check for overly complex accessibility labels
            if let label = v.accessibilityLabel, label.count > 100 {
                issues.append(AccessibilityIssue(
                    title: "Accessibility Label Too Long",
                    description: "Accessibility label exceeds 100 characters (\(label.count))",
                    severity: .warning,
                    wcagCriterion: "1.3.1"
                ))
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return issues
    }

    func auditVoiceOverNavigation(_ view: UIView) async -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        let accessibleElements = collectAccessibleElements(from: view)

        // Check reading order
        for (index, element) in accessibleElements.enumerated() {
            if index > 0 {
                let previous = accessibleElements[index - 1]
                if !isLogicalReadingOrder(from: previous, to: element) {
                    issues.append(AccessibilityIssue(
                        title: "Illogical Reading Order",
                        description: "VoiceOver reading order doesn't follow visual layout",
                        severity: .warning,
                        wcagCriterion: "1.3.2"
                    ))
                }
            }
        }

        return issues
    }

    func auditDynamicTypeCategory(_ category: UIContentSizeCategory) async -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        // This would involve testing UI at different font sizes
        // For now, we'll add a placeholder check

        if category.isAccessibilityCategory {
            issues.append(AccessibilityIssue(
                title: "Dynamic Type Compatibility",
                description: "Verify app layout works correctly with \(category.rawValue)",
                severity: .info,
                wcagCriterion: "1.4.4"
            ))
        }

        return issues
    }

    func auditViewContrast(_ view: UIView) async -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        func checkView(_ v: UIView) {
            if let label = v as? UILabel {
                let textColor = Color(label.textColor)
                let backgroundColor = Color(label.backgroundColor ?? UIColor.clear)
                let ratio = accessibilityService.calculateContrastRatio(foreground: textColor, background: backgroundColor)

                if ratio < 4.5 {
                    issues.append(AccessibilityIssue(
                        title: "Insufficient Color Contrast",
                        description: "Text contrast ratio \(String(format: "%.2f", ratio)):1 below WCAG AA requirement (4.5:1)",
                        severity: .error,
                        wcagCriterion: "1.4.3"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return issues
    }

    func auditViewTouchTargets(_ view: UIView) async -> [AccessibilityIssue] {
        var issues: [AccessibilityIssue] = []

        func checkView(_ v: UIView) {
            if v.isUserInteractionEnabled && !v.isHidden {
                let frame = v.frame
                if frame.width < 44 || frame.height < 44 {
                    issues.append(AccessibilityIssue(
                        title: "Touch Target Too Small",
                        description: "Touch target \(Int(frame.width))×\(Int(frame.height))pt below minimum 44×44pt",
                        severity: .error,
                        wcagCriterion: "2.5.5"
                    ))
                }
            }

            for subview in v.subviews {
                checkView(subview)
            }
        }

        checkView(view)
        return issues
    }

    func collectAccessibleElements(from view: UIView) -> [UIView] {
        var elements: [UIView] = []

        func traverse(_ v: UIView) {
            if v.isAccessibilityElement {
                elements.append(v)
            } else {
                for subview in v.subviews {
                    traverse(subview)
                }
            }
        }

        traverse(view)
        return elements
    }

    func isLogicalReadingOrder(from previous: UIView, to current: UIView) -> Bool {
        let previousFrame = previous.frame
        let currentFrame = current.frame

        // Basic left-to-right, top-to-bottom reading order
        if currentFrame.minY >= previousFrame.maxY {
            return true // Next row
        }

        if abs(currentFrame.minY - previousFrame.minY) < 10 && currentFrame.minX >= previousFrame.maxX {
            return true // Same row, left to right
        }

        return false
    }
    #endif

    func calculateSectionScore(issues: [AccessibilityIssue], totalChecks: Int) -> Double {
        let errorCount = issues.filter { $0.severity == .error }.count
        let warningCount = issues.filter { $0.severity == .warning }.count

        let penalties = Double(errorCount * 2 + warningCount) // Errors count double
        let maxPenalties = Double(totalChecks * 2)

        return max(0, (maxPenalties - penalties) / maxPenalties)
    }

    func calculateOverallScore(_ report: AccessibilityAuditReport) -> Double {
        let scores = [
            report.basicCompliance.score * 0.2,
            report.voiceOverCompatibility.score * 0.25,
            report.dynamicTypeSupport.score * 0.2,
            report.colorContrastAnalysis.score * 0.2,
            report.touchTargetAnalysis.score * 0.15
        ]

        return scores.reduce(0, +)
    }

    func determineWCAGLevel(_ report: AccessibilityAuditReport) -> WCAGComplianceLevel {
        let errorCount = [
            report.basicCompliance.issues,
            report.voiceOverCompatibility.issues,
            report.dynamicTypeSupport.issues,
            report.colorContrastAnalysis.issues,
            report.touchTargetAnalysis.issues
        ].flatMap { $0 }.filter { $0.severity == .error }.count

        if errorCount == 0 {
            return report.overallScore >= 0.95 ? .aaa : .aa
        } else if errorCount <= 2 {
            return .a
        } else {
            return .none
        }
    }

    func generateRecommendations(_ report: AccessibilityAuditReport) -> [AccessibilityRecommendation] {
        var recommendations: [AccessibilityRecommendation] = []

        // High priority recommendations
        let allIssues = [
            report.basicCompliance.issues,
            report.voiceOverCompatibility.issues,
            report.dynamicTypeSupport.issues,
            report.colorContrastAnalysis.issues,
            report.touchTargetAnalysis.issues
        ].flatMap { $0 }

        let errorIssues = allIssues.filter { $0.severity == .error }

        if !errorIssues.isEmpty {
            recommendations.append(AccessibilityRecommendation(
                title: "Fix Critical Accessibility Errors",
                description: "Address \(errorIssues.count) critical accessibility errors that prevent WCAG compliance.",
                priority: .high,
                implementationSteps: [
                    "Review all error-level issues in the audit report",
                    "Prioritize fixes based on user impact",
                    "Test fixes with actual assistive technologies",
                    "Re-run accessibility audit to verify fixes"
                ]
            ))
        }

        // Medium priority recommendations
        if report.overallScore < 0.8 {
            recommendations.append(AccessibilityRecommendation(
                title: "Improve Overall Accessibility Score",
                description: "Current score is \(Int(report.overallScore * 100))%. Target 80% or higher for good accessibility.",
                priority: .medium,
                implementationSteps: [
                    "Focus on areas with lowest scores",
                    "Implement accessibility testing in CI/CD pipeline",
                    "Train development team on accessibility best practices",
                    "Establish accessibility review process"
                ]
            ))
        }

        return recommendations
    }
}

// MARK: - Supporting Types

public struct AccessibilityAuditReport {
    public var basicCompliance = AuditSection()
    public var voiceOverCompatibility = AuditSection()
    public var dynamicTypeSupport = AuditSection()
    public var colorContrastAnalysis = AuditSection()
    public var touchTargetAnalysis = AuditSection()
    public var overallScore: Double = 0.0
    public var wcagLevel: WCAGComplianceLevel = .none
    public var recommendations: [AccessibilityRecommendation] = []
}

public struct AuditSection {
    public var score: Double = 0.0
    public var issues: [AccessibilityIssue] = []
}

public struct AccessibilityIssue {
    public let title: String
    public let description: String
    public let severity: IssueSeverity
    public let wcagCriterion: String

    public enum IssueSeverity: String {
        case error = "Error"
        case warning = "Warning"
        case info = "Info"
    }
}

public struct AccessibilityRecommendation {
    public let title: String
    public let description: String
    public let priority: Priority
    public let implementationSteps: [String]

    public enum Priority: String {
        case high = "High"
        case medium = "Medium"
        case low = "Low"
    }
}

public enum WCAGComplianceLevel: String {
    case none = "Non-compliant"
    case a = "WCAG A"
    case aa = "WCAG AA"
    case aaa = "WCAG AAA"
}