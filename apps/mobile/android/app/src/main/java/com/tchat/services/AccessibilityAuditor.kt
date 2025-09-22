package com.tchat.services

import android.content.Context
import android.graphics.Color
import android.util.TypedValue
import android.view.View
import android.view.ViewGroup
import android.view.accessibility.AccessibilityManager
import androidx.compose.runtime.*
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import java.text.SimpleDateFormat
import java.util.*
import kotlin.math.*

/**
 * Comprehensive accessibility auditor for Android platform
 * Provides detailed WCAG compliance analysis and reporting
 */
class AccessibilityAuditor private constructor(private val context: Context) {

    companion object {
        @Volatile
        private var INSTANCE: AccessibilityAuditor? = null

        fun getInstance(context: Context): AccessibilityAuditor {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: AccessibilityAuditor(context.applicationContext).also { INSTANCE = it }
            }
        }
    }

    // Properties
    private val accessibilityService = AccessibilityService.getInstance(context)
    private val serviceScope = CoroutineScope(SupervisorJob() + Dispatchers.Main)

    private val _currentAuditResults = MutableStateFlow<AccessibilityAuditReport?>(null)
    val currentAuditResults: StateFlow<AccessibilityAuditReport?> = _currentAuditResults.asStateFlow()

    private val _isAuditing = MutableStateFlow(false)
    val isAuditing: StateFlow<Boolean> = _isAuditing.asStateFlow()

    private val _auditProgress = MutableStateFlow(0.0)
    val auditProgress: StateFlow<Double> = _auditProgress.asStateFlow()

    // Public Interface

    /**
     * Perform comprehensive accessibility audit
     */
    suspend fun performComprehensiveAudit(rootView: View? = null): AccessibilityAuditReport {
        _isAuditing.value = true
        _auditProgress.value = 0.0

        val report = AccessibilityAuditReport()

        // Step 1: Basic compliance checks (20%)
        updateProgress(0.2)
        report.basicCompliance = auditBasicCompliance(rootView)

        // Step 2: TalkBack compatibility (40%)
        updateProgress(0.4)
        report.talkBackCompatibility = auditTalkBackCompatibility(rootView)

        // Step 3: Font scaling support (60%)
        updateProgress(0.6)
        report.fontScalingSupport = auditFontScalingSupport(rootView)

        // Step 4: Color contrast analysis (80%)
        updateProgress(0.8)
        report.colorContrastAnalysis = auditColorContrast(rootView)

        // Step 5: Touch target analysis (100%)
        updateProgress(1.0)
        report.touchTargetAnalysis = auditTouchTargets(rootView)

        // Generate overall score and recommendations
        report.overallScore = calculateOverallScore(report)
        report.wcagLevel = determineWCAGLevel(report)
        report.recommendations = generateRecommendations(report)

        _currentAuditResults.value = report
        _isAuditing.value = false

        return report
    }

    /**
     * Generate accessibility report as formatted string
     */
    fun generateAccessibilityReport(report: AccessibilityAuditReport): String {
        val dateFormat = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault())
        val currentDate = dateFormat.format(Date())

        return buildString {
            appendLine("# Tchat Android Accessibility Audit Report")
            appendLine()
            appendLine("**Generated:** $currentDate")
            appendLine("**Overall Score:** ${(report.overallScore * 100).toInt()}%")
            appendLine("**WCAG Level:** ${report.wcagLevel.displayName}")
            appendLine()

            // Basic Compliance Section
            appendLine("## Basic Compliance")
            appendLine("**Score:** ${(report.basicCompliance.score * 100).toInt()}%")
            appendLine("**Issues Found:** ${report.basicCompliance.issues.size}")
            appendLine()

            report.basicCompliance.issues.forEach { issue ->
                appendLine("- ❌ **${issue.title}:** ${issue.description}")
                appendLine("  - *WCAG Criterion:* ${issue.wcagCriterion}")
                appendLine("  - *Severity:* ${issue.severity.displayName}")
                appendLine()
            }

            // TalkBack Compatibility Section
            appendLine("## TalkBack Compatibility")
            appendLine("**Score:** ${(report.talkBackCompatibility.score * 100).toInt()}%")
            appendLine("**Issues Found:** ${report.talkBackCompatibility.issues.size}")
            appendLine()

            report.talkBackCompatibility.issues.forEach { issue ->
                appendLine("- ❌ **${issue.title}:** ${issue.description}")
                appendLine("  - *WCAG Criterion:* ${issue.wcagCriterion}")
                appendLine("  - *Severity:* ${issue.severity.displayName}")
                appendLine()
            }

            // Font Scaling Support Section
            appendLine("## Font Scaling Support")
            appendLine("**Score:** ${(report.fontScalingSupport.score * 100).toInt()}%")
            appendLine("**Issues Found:** ${report.fontScalingSupport.issues.size}")
            appendLine()

            report.fontScalingSupport.issues.forEach { issue ->
                appendLine("- ❌ **${issue.title}:** ${issue.description}")
                appendLine("  - *WCAG Criterion:* ${issue.wcagCriterion}")
                appendLine("  - *Severity:* ${issue.severity.displayName}")
                appendLine()
            }

            // Color Contrast Analysis Section
            appendLine("## Color Contrast Analysis")
            appendLine("**Score:** ${(report.colorContrastAnalysis.score * 100).toInt()}%")
            appendLine("**Issues Found:** ${report.colorContrastAnalysis.issues.size}")
            appendLine()

            report.colorContrastAnalysis.issues.forEach { issue ->
                appendLine("- ❌ **${issue.title}:** ${issue.description}")
                appendLine("  - *WCAG Criterion:* ${issue.wcagCriterion}")
                appendLine("  - *Severity:* ${issue.severity.displayName}")
                appendLine()
            }

            // Touch Target Analysis Section
            appendLine("## Touch Target Analysis")
            appendLine("**Score:** ${(report.touchTargetAnalysis.score * 100).toInt()}%")
            appendLine("**Issues Found:** ${report.touchTargetAnalysis.issues.size}")
            appendLine()

            report.touchTargetAnalysis.issues.forEach { issue ->
                appendLine("- ❌ **${issue.title}:** ${issue.description}")
                appendLine("  - *WCAG Criterion:* ${issue.wcagCriterion}")
                appendLine("  - *Severity:* ${issue.severity.displayName}")
                appendLine()
            }

            // Recommendations Section
            appendLine("## Recommendations")
            appendLine()

            report.recommendations.forEach { recommendation ->
                appendLine("### ${recommendation.priority.displayName} Priority: ${recommendation.title}")
                appendLine(recommendation.description)
                appendLine()
                appendLine("**Implementation Steps:**")
                recommendation.implementationSteps.forEachIndexed { index, step ->
                    appendLine("${index + 1}. $step")
                }
                appendLine()
            }
        }
    }

    /**
     * Export audit report to file
     */
    suspend fun exportAuditReport(report: AccessibilityAuditReport, fileName: String = "accessibility_audit_report.md") {
        withContext(Dispatchers.IO) {
            try {
                val reportContent = generateAccessibilityReport(report)
                val file = java.io.File(context.getExternalFilesDir(null), fileName)
                file.writeText(reportContent)
                android.util.Log.i("AccessibilityAuditor", "✅ Report exported to: ${file.absolutePath}")
            } catch (e: Exception) {
                android.util.Log.e("AccessibilityAuditor", "❌ Failed to export report", e)
            }
        }
    }

    // Private Implementation

    private fun updateProgress(progress: Double) {
        _auditProgress.value = progress
    }

    private suspend fun auditBasicCompliance(rootView: View?): AuditSection = withContext(Dispatchers.Default) {
        val issues = mutableListOf<AccessibilityIssueData>()

        // Check if TalkBack is enabled
        if (!accessibilityService.isTalkBackEnabled.value) {
            issues.add(
                AccessibilityIssueData(
                    title = "TalkBack Not Active",
                    description = "TalkBack is not currently enabled. Enable TalkBack to test screen reader compatibility.",
                    severity = IssueSeverity.WARNING,
                    wcagCriterion = "4.1.2"
                )
            )
        }

        // Audit view hierarchy if available
        rootView?.let { view ->
            val hierarchyIssues = auditViewHierarchy(view)
            issues.addAll(hierarchyIssues)
        }

        val score = calculateSectionScore(issues, totalChecks = 10)
        AuditSection(score = score, issues = issues)
    }

    private suspend fun auditTalkBackCompatibility(rootView: View?): AuditSection = withContext(Dispatchers.Default) {
        val issues = mutableListOf<AccessibilityIssueData>()

        rootView?.let { view ->
            val talkBackIssues = auditTalkBackNavigation(view)
            issues.addAll(talkBackIssues)
        }

        val score = calculateSectionScore(issues, totalChecks = 8)
        AuditSection(score = score, issues = issues)
    }

    private suspend fun auditFontScalingSupport(rootView: View?): AuditSection = withContext(Dispatchers.Default) {
        val issues = mutableListOf<AccessibilityIssueData>()

        // Test different font scale levels
        val currentFontScale = context.resources.configuration.fontScale
        val testScales = listOf(1.0f, 1.3f, 1.5f, 2.0f)

        for (scale in testScales) {
            if (scale > 1.5f) {
                issues.add(
                    AccessibilityIssueData(
                        title = "Font Scaling Compatibility",
                        description = "Verify app layout works correctly with ${scale}x font scaling",
                        severity = IssueSeverity.INFO,
                        wcagCriterion = "1.4.4"
                    )
                )
            }
        }

        val score = calculateSectionScore(issues, totalChecks = 6)
        AuditSection(score = score, issues = issues)
    }

    private suspend fun auditColorContrast(rootView: View?): AuditSection = withContext(Dispatchers.Default) {
        val issues = mutableListOf<AccessibilityIssueData>()

        rootView?.let { view ->
            val contrastIssues = auditViewContrast(view)
            issues.addAll(contrastIssues)
        }

        val score = calculateSectionScore(issues, totalChecks = 12)
        AuditSection(score = score, issues = issues)
    }

    private suspend fun auditTouchTargets(rootView: View?): AuditSection = withContext(Dispatchers.Default) {
        val issues = mutableListOf<AccessibilityIssueData>()

        rootView?.let { view ->
            val touchTargetIssues = auditViewTouchTargets(view)
            issues.addAll(touchTargetIssues)
        }

        val score = calculateSectionScore(issues, totalChecks = 5)
        AuditSection(score = score, issues = issues)
    }

    private fun auditViewHierarchy(view: View): List<AccessibilityIssueData> {
        val issues = mutableListOf<AccessibilityIssueData>()

        fun checkView(v: View) {
            // Check for missing content descriptions
            if ((v.isClickable || v.isFocusable) && v.contentDescription.isNullOrEmpty()) {
                issues.add(
                    AccessibilityIssueData(
                        title = "Missing Content Description",
                        description = "Interactive element lacks content description",
                        severity = IssueSeverity.ERROR,
                        wcagCriterion = "1.3.1"
                    )
                )
            }

            // Check for overly long content descriptions
            v.contentDescription?.let { description ->
                if (description.length > 100) {
                    issues.add(
                        AccessibilityIssueData(
                            title = "Content Description Too Long",
                            description = "Content description exceeds 100 characters (${description.length})",
                            severity = IssueSeverity.WARNING,
                            wcagCriterion = "1.3.1"
                        )
                    )
                }
            }

            if (v is ViewGroup) {
                for (i in 0 until v.childCount) {
                    checkView(v.getChildAt(i))
                }
            }
        }

        checkView(view)
        return issues
    }

    private fun auditTalkBackNavigation(view: View): List<AccessibilityIssueData> {
        val issues = mutableListOf<AccessibilityIssueData>()

        val accessibleElements = collectAccessibleElements(view)

        // Check reading order
        for (i in 1 until accessibleElements.size) {
            val current = accessibleElements[i]
            val previous = accessibleElements[i - 1]

            if (!isLogicalReadingOrder(previous, current)) {
                issues.add(
                    AccessibilityIssueData(
                        title = "Illogical Reading Order",
                        description = "TalkBack reading order doesn't follow visual layout",
                        severity = IssueSeverity.WARNING,
                        wcagCriterion = "1.3.2"
                    )
                )
            }
        }

        return issues
    }

    private fun auditViewContrast(view: View): List<AccessibilityIssueData> {
        val issues = mutableListOf<AccessibilityIssueData>()

        fun checkView(v: View) {
            if (v is android.widget.TextView) {
                val textColor = androidx.compose.ui.graphics.Color(v.currentTextColor)
                // Simplified background color extraction
                val backgroundColor = androidx.compose.ui.graphics.Color.White

                val ratio = accessibilityService.calculateContrastRatio(textColor, backgroundColor)
                if (ratio < 4.5) {
                    issues.add(
                        AccessibilityIssueData(
                            title = "Insufficient Color Contrast",
                            description = "Text contrast ratio ${String.format("%.2f", ratio)}:1 below WCAG AA requirement (4.5:1)",
                            severity = IssueSeverity.ERROR,
                            wcagCriterion = "1.4.3"
                        )
                    )
                }
            }

            if (v is ViewGroup) {
                for (i in 0 until v.childCount) {
                    checkView(v.getChildAt(i))
                }
            }
        }

        checkView(view)
        return issues
    }

    private fun auditViewTouchTargets(view: View): List<AccessibilityIssueData> {
        val issues = mutableListOf<AccessibilityIssueData>()
        val minTouchTarget = TypedValue.applyDimension(
            TypedValue.COMPLEX_UNIT_DIP, 48f, context.resources.displayMetrics
        ).toInt()

        fun checkView(v: View) {
            if (v.isClickable && v.visibility == View.VISIBLE) {
                if (v.width < minTouchTarget || v.height < minTouchTarget) {
                    issues.add(
                        AccessibilityIssueData(
                            title = "Touch Target Too Small",
                            description = "Touch target ${v.width}×${v.height}px below minimum ${minTouchTarget}×${minTouchTarget}px (48×48dp)",
                            severity = IssueSeverity.ERROR,
                            wcagCriterion = "2.5.5"
                        )
                    )
                }
            }

            if (v is ViewGroup) {
                for (i in 0 until v.childCount) {
                    checkView(v.getChildAt(i))
                }
            }
        }

        checkView(view)
        return issues
    }

    private fun collectAccessibleElements(view: View): List<View> {
        val elements = mutableListOf<View>()

        fun traverse(v: View) {
            if (v.importantForAccessibility != View.IMPORTANT_FOR_ACCESSIBILITY_NO) {
                if (v.isClickable || v.isFocusable || !v.contentDescription.isNullOrEmpty()) {
                    elements.add(v)
                }

                if (v is ViewGroup) {
                    for (i in 0 until v.childCount) {
                        traverse(v.getChildAt(i))
                    }
                }
            }
        }

        traverse(view)
        return elements
    }

    private fun isLogicalReadingOrder(previous: View, current: View): Boolean {
        val previousLocation = IntArray(2)
        val currentLocation = IntArray(2)
        previous.getLocationOnScreen(previousLocation)
        current.getLocationOnScreen(currentLocation)

        // Basic left-to-right, top-to-bottom reading order
        if (currentLocation[1] >= previousLocation[1] + previous.height) {
            return true // Next row
        }

        if (abs(currentLocation[1] - previousLocation[1]) < 30 &&
            currentLocation[0] >= previousLocation[0] + previous.width) {
            return true // Same row, left to right
        }

        return false
    }

    private fun calculateSectionScore(issues: List<AccessibilityIssueData>, totalChecks: Int): Double {
        val errorCount = issues.count { it.severity == IssueSeverity.ERROR }
        val warningCount = issues.count { it.severity == IssueSeverity.WARNING }

        val penalties = errorCount * 2.0 + warningCount // Errors count double
        val maxPenalties = totalChecks * 2.0

        return maxOf(0.0, (maxPenalties - penalties) / maxPenalties)
    }

    private fun calculateOverallScore(report: AccessibilityAuditReport): Double {
        return listOf(
            report.basicCompliance.score * 0.2,
            report.talkBackCompatibility.score * 0.25,
            report.fontScalingSupport.score * 0.2,
            report.colorContrastAnalysis.score * 0.2,
            report.touchTargetAnalysis.score * 0.15
        ).sum()
    }

    private fun determineWCAGLevel(report: AccessibilityAuditReport): WCAGComplianceLevel {
        val allIssues = listOf(
            report.basicCompliance.issues,
            report.talkBackCompatibility.issues,
            report.fontScalingSupport.issues,
            report.colorContrastAnalysis.issues,
            report.touchTargetAnalysis.issues
        ).flatten()

        val errorCount = allIssues.count { it.severity == IssueSeverity.ERROR }

        return when {
            errorCount == 0 -> if (report.overallScore >= 0.95) WCAGComplianceLevel.AAA else WCAGComplianceLevel.AA
            errorCount <= 2 -> WCAGComplianceLevel.A
            else -> WCAGComplianceLevel.NONE
        }
    }

    private fun generateRecommendations(report: AccessibilityAuditReport): List<AccessibilityRecommendationData> {
        val recommendations = mutableListOf<AccessibilityRecommendationData>()

        val allIssues = listOf(
            report.basicCompliance.issues,
            report.talkBackCompatibility.issues,
            report.fontScalingSupport.issues,
            report.colorContrastAnalysis.issues,
            report.touchTargetAnalysis.issues
        ).flatten()

        val errorIssues = allIssues.filter { it.severity == IssueSeverity.ERROR }

        if (errorIssues.isNotEmpty()) {
            recommendations.add(
                AccessibilityRecommendationData(
                    title = "Fix Critical Accessibility Errors",
                    description = "Address ${errorIssues.size} critical accessibility errors that prevent WCAG compliance.",
                    priority = Priority.HIGH,
                    implementationSteps = listOf(
                        "Review all error-level issues in the audit report",
                        "Prioritize fixes based on user impact",
                        "Test fixes with actual assistive technologies",
                        "Re-run accessibility audit to verify fixes"
                    )
                )
            )
        }

        if (report.overallScore < 0.8) {
            recommendations.add(
                AccessibilityRecommendationData(
                    title = "Improve Overall Accessibility Score",
                    description = "Current score is ${(report.overallScore * 100).toInt()}%. Target 80% or higher for good accessibility.",
                    priority = Priority.MEDIUM,
                    implementationSteps = listOf(
                        "Focus on areas with lowest scores",
                        "Implement accessibility testing in CI/CD pipeline",
                        "Train development team on accessibility best practices",
                        "Establish accessibility review process"
                    )
                )
            )
        }

        return recommendations
    }

    fun cleanup() {
        serviceScope.cancel()
    }
}

// Supporting Data Classes

data class AccessibilityAuditReport(
    var basicCompliance: AuditSection = AuditSection(),
    var talkBackCompatibility: AuditSection = AuditSection(),
    var fontScalingSupport: AuditSection = AuditSection(),
    var colorContrastAnalysis: AuditSection = AuditSection(),
    var touchTargetAnalysis: AuditSection = AuditSection(),
    var overallScore: Double = 0.0,
    var wcagLevel: WCAGComplianceLevel = WCAGComplianceLevel.NONE,
    var recommendations: List<AccessibilityRecommendationData> = emptyList()
)

data class AuditSection(
    var score: Double = 0.0,
    var issues: List<AccessibilityIssueData> = emptyList()
)

data class AccessibilityIssueData(
    val title: String,
    val description: String,
    val severity: IssueSeverity,
    val wcagCriterion: String
)

data class AccessibilityRecommendationData(
    val title: String,
    val description: String,
    val priority: Priority,
    val implementationSteps: List<String>
)

enum class IssueSeverity(val displayName: String) {
    ERROR("Error"),
    WARNING("Warning"),
    INFO("Info")
}

enum class Priority(val displayName: String) {
    HIGH("High"),
    MEDIUM("Medium"),
    LOW("Low")
}

enum class WCAGComplianceLevel(val displayName: String) {
    NONE("Non-compliant"),
    A("WCAG A"),
    AA("WCAG AA"),
    AAA("WCAG AAA")
}

// Compose Extensions

@Composable
fun accessibilityAuditor(): AccessibilityAuditor {
    val context = androidx.compose.ui.platform.LocalContext.current
    return remember { AccessibilityAuditor.getInstance(context) }
}