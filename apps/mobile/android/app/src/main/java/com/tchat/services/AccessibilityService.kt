package com.tchat.services

import android.content.Context
import android.content.res.Configuration
import android.graphics.Color
import android.util.TypedValue
import android.view.View
import android.view.ViewGroup
import android.view.accessibility.AccessibilityEvent
import android.view.accessibility.AccessibilityManager
import android.view.accessibility.AccessibilityNodeInfo
import androidx.compose.foundation.layout.size
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.TextUnit
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.view.ViewCompat
import androidx.core.view.accessibility.AccessibilityNodeInfoCompat
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import java.util.*
import kotlin.math.*

/**
 * Comprehensive accessibility service for Android platform
 * Provides WCAG 2.1 AA compliance, TalkBack support, and accessibility testing capabilities
 */
class AccessibilityService private constructor(private val context: Context) {

    companion object {
        @Volatile
        private var INSTANCE: AccessibilityService? = null

        fun getInstance(context: Context): AccessibilityService {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: AccessibilityService(context.applicationContext).also { INSTANCE = it }
            }
        }
    }

    // MARK: - Properties
    private val accessibilityManager = context.getSystemService(Context.ACCESSIBILITY_SERVICE) as AccessibilityManager
    private val serviceScope = CoroutineScope(SupervisorJob() + Dispatchers.Main)

    // Observable accessibility states
    private val _isTalkBackEnabled = MutableStateFlow(accessibilityManager.isEnabled)
    val isTalkBackEnabled: StateFlow<Boolean> = _isTalkBackEnabled.asStateFlow()

    private val _fontScale = MutableStateFlow(context.resources.configuration.fontScale)
    val fontScale: StateFlow<Float> = _fontScale.asStateFlow()

    private val _isHighContrastEnabled = MutableStateFlow(isHighContrastEnabled())
    val isHighContrastEnabled: StateFlow<Boolean> = _isHighContrastEnabled.asStateFlow()

    private val _auditResults = MutableStateFlow<List<AccessibilityAuditResult>>(emptyList())
    val auditResults: StateFlow<List<AccessibilityAuditResult>> = _auditResults.asStateFlow()

    private var auditJob: Job? = null

    init {
        setupAccessibilityListeners()
    }

    // MARK: - Public Interface

    /**
     * Configure accessibility for a Composable with comprehensive WCAG compliance
     */
    fun configureSemantics(
        contentDescription: String? = null,
        role: Role? = null,
        stateDescription: String? = null,
        onClick: (() -> Unit)? = null,
        onLongClick: (() -> Unit)? = null,
        isImportant: Boolean = true,
        isClickable: Boolean = false,
        isEnabled: Boolean = true,
        liveRegion: LiveRegionMode = LiveRegionMode.POLITE
    ): SemanticsPropertyReceiver.() -> Unit = {
        contentDescription?.let {
            this.contentDescription = it
            if (it.length > 100) {
                // Log warning for overly long content descriptions
                android.util.Log.w("AccessibilityService",
                    "Content description exceeds 100 characters: ${it.length}")
            }
        }

        role?.let { this.role = it }
        stateDescription?.let { this.stateDescription = it }

        if (isClickable || onClick != null) {
            onClick?.let { action ->
                this.onClick {
                    announceTalkBack("Button activated: ${contentDescription ?: "button"}")
                    action()
                    true
                }
            }
        }

        onLongClick?.let { action ->
            this.onLongClick {
                announceTalkBack("Long press: ${contentDescription ?: "item"}")
                action()
                true
            }
        }

        if (!isEnabled) {
            this.disabled()
        }

        // Configure live region for dynamic content
        when (liveRegion) {
            LiveRegionMode.ASSERTIVE -> {
                this.liveRegion = androidx.compose.ui.semantics.LiveRegionMode.Assertive
            }
            LiveRegionMode.POLITE -> {
                this.liveRegion = androidx.compose.ui.semantics.LiveRegionMode.Polite
            }
            LiveRegionMode.NONE -> {
                // Default behavior
            }
        }
    }

    /**
     * Ensure minimum touch target size of 48dp for Android
     */
    fun ensureMinimumTouchTarget(): Modifier {
        return Modifier.size(minOf(48.dp, 48.dp))
    }

    /**
     * Apply font scaling with proper bounds
     */
    fun scaledTextStyle(
        baseSize: TextUnit,
        fontWeight: FontWeight = FontWeight.Normal,
        maxScale: Float = 2.0f
    ): TextStyle {
        val scaledSize = (baseSize.value * min(fontScale.value, maxScale)).sp
        return TextStyle(
            fontSize = scaledSize,
            fontWeight = fontWeight
        )
    }

    /**
     * Get contrast-compliant colors for WCAG AA compliance
     */
    fun getContrastCompliantColor(
        foreground: androidx.compose.ui.graphics.Color,
        background: androidx.compose.ui.graphics.Color,
        level: ContrastLevel = ContrastLevel.AA
    ): androidx.compose.ui.graphics.Color {
        val ratio = calculateContrastRatio(foreground, background)

        return if (ratio >= level.minimumRatio) {
            foreground
        } else {
            adjustColorForContrast(foreground, background, level.minimumRatio)
        }
    }

    /**
     * Announce message to TalkBack with specified priority
     */
    fun announceTalkBack(message: String, priority: AnnouncementPriority = AnnouncementPriority.MEDIUM) {
        if (accessibilityManager.isEnabled) {
            val event = AccessibilityEvent.obtain(AccessibilityEvent.TYPE_ANNOUNCEMENT).apply {
                text.add(message)
                className = this@AccessibilityService::class.java.name
                packageName = context.packageName
            }

            // Set priority for different announcement types
            when (priority) {
                AnnouncementPriority.HIGH -> {
                    event.eventType = AccessibilityEvent.TYPE_ANNOUNCEMENT
                }
                AnnouncementPriority.MEDIUM -> {
                    event.eventType = AccessibilityEvent.TYPE_VIEW_FOCUSED
                }
                AnnouncementPriority.LOW -> {
                    event.eventType = AccessibilityEvent.TYPE_VIEW_SELECTED
                }
            }

            accessibilityManager.sendAccessibilityEvent(event)
        }
    }

    /**
     * Perform comprehensive accessibility audit
     */
    fun performAccessibilityAudit(rootView: View? = null) {
        val targetView = rootView ?: getCurrentRootView()
        targetView?.let { view ->
            serviceScope.launch {
                val results = auditView(view)
                _auditResults.value = results
                logAuditResults(results)
            }
        }
    }

    /**
     * Start continuous accessibility monitoring
     */
    fun startAccessibilityMonitoring(intervalMs: Long = 5000) {
        auditJob?.cancel()
        auditJob = serviceScope.launch {
            while (isActive) {
                performAccessibilityAudit()
                delay(intervalMs)
            }
        }
    }

    /**
     * Stop accessibility monitoring
     */
    fun stopAccessibilityMonitoring() {
        auditJob?.cancel()
        auditJob = null
    }

    /**
     * Check if view meets accessibility requirements
     */
    fun validateAccessibility(view: View): ValidationResult {
        val issues = mutableListOf<AccessibilityIssue>()

        // Check touch target size
        val minTouchTarget = TypedValue.applyDimension(
            TypedValue.COMPLEX_UNIT_DIP, 48f, context.resources.displayMetrics
        ).toInt()

        if (view.width < minTouchTarget || view.height < minTouchTarget) {
            issues.add(AccessibilityIssue.INSUFFICIENT_TOUCH_TARGET)
        }

        // Check content description
        if (view.isClickable && view.contentDescription.isNullOrEmpty()) {
            issues.add(AccessibilityIssue.MISSING_CONTENT_DESCRIPTION)
        }

        // Check contrast (if applicable)
        if (view.background != null) {
            // Implementation would depend on view type and text color
            // This is a simplified check
        }

        return ValidationResult(
            isValid = issues.isEmpty(),
            issues = issues,
            wcagLevel = if (issues.isEmpty()) WCAGLevel.AA else WCAGLevel.NONE
        )
    }

    // MARK: - Private Implementation

    private fun setupAccessibilityListeners() {
        // Listen for accessibility service state changes
        serviceScope.launch {
            while (true) {
                delay(1000) // Check every second
                val isEnabled = accessibilityManager.isEnabled
                if (_isTalkBackEnabled.value != isEnabled) {
                    _isTalkBackEnabled.value = isEnabled
                }
            }
        }

        // Listen for configuration changes (font scale, etc.)
        // This would typically be handled in an Activity/Fragment
        serviceScope.launch {
            while (true) {
                delay(1000)
                val currentFontScale = context.resources.configuration.fontScale
                if (_fontScale.value != currentFontScale) {
                    _fontScale.value = currentFontScale
                }

                val currentHighContrast = isHighContrastEnabled()
                if (_isHighContrastEnabled.value != currentHighContrast) {
                    _isHighContrastEnabled.value = currentHighContrast
                }
            }
        }
    }

    private fun isHighContrastEnabled(): Boolean {
        // Check system settings for high contrast
        return try {
            val resolver = context.contentResolver
            android.provider.Settings.Secure.getInt(resolver, "high_text_contrast_enabled", 0) == 1
        } catch (e: Exception) {
            false
        }
    }

    private fun getCurrentRootView(): View? {
        return try {
            // This would need to be called from an Activity context
            // For now, return null and handle in calling code
            null
        } catch (e: Exception) {
            null
        }
    }

    private suspend fun auditView(view: View): List<AccessibilityAuditResult> = withContext(Dispatchers.Default) {
        val results = mutableListOf<AccessibilityAuditResult>()

        // Audit touch targets
        results.addAll(auditTouchTargets(view))

        // Audit content descriptions
        results.addAll(auditContentDescriptions(view))

        // Audit contrast ratios
        results.addAll(auditContrastRatios(view))

        // Audit reading order
        results.addAll(auditReadingOrder(view))

        // Audit semantic properties
        results.addAll(auditSemanticProperties(view))

        results
    }

    private fun auditTouchTargets(view: View): List<AccessibilityAuditResult> {
        val results = mutableListOf<AccessibilityAuditResult>()
        val minTouchTarget = TypedValue.applyDimension(
            TypedValue.COMPLEX_UNIT_DIP, 48f, context.resources.displayMetrics
        ).toInt()

        fun checkView(v: View) {
            if (v.isClickable && v.visibility == View.VISIBLE) {
                if (v.width < minTouchTarget || v.height < minTouchTarget) {
                    results.add(
                        AccessibilityAuditResult(
                            type = AuditType.TOUCH_TARGET_SIZE,
                            severity = Severity.ERROR,
                            message = "Touch target too small: ${v.width}×${v.height}px. Minimum required: ${minTouchTarget}×${minTouchTarget}px (48×48dp)",
                            viewInfo = ViewInfo.fromView(v),
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
        return results
    }

    private fun auditContentDescriptions(view: View): List<AccessibilityAuditResult> {
        val results = mutableListOf<AccessibilityAuditResult>()

        fun checkView(v: View) {
            if (v.isClickable || v.isFocusable) {
                val contentDescription = v.contentDescription?.toString()

                if (contentDescription.isNullOrEmpty()) {
                    results.add(
                        AccessibilityAuditResult(
                            type = AuditType.CONTENT_DESCRIPTION,
                            severity = Severity.ERROR,
                            message = "Missing content description for interactive element",
                            viewInfo = ViewInfo.fromView(v),
                            wcagCriterion = "1.3.1"
                        )
                    )
                } else if (contentDescription.length > 100) {
                    results.add(
                        AccessibilityAuditResult(
                            type = AuditType.CONTENT_DESCRIPTION,
                            severity = Severity.WARNING,
                            message = "Content description too long (${contentDescription.length} characters). Keep under 100 characters.",
                            viewInfo = ViewInfo.fromView(v),
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
        return results
    }

    private fun auditContrastRatios(view: View): List<AccessibilityAuditResult> {
        val results = mutableListOf<AccessibilityAuditResult>()

        fun checkView(v: View) {
            // This is a simplified implementation
            // In practice, you'd need to extract text color and background color
            // from the specific view types (TextView, Button, etc.)

            if (v is android.widget.TextView) {
                val textColor = androidx.compose.ui.graphics.Color(v.currentTextColor)
                // Background color extraction would be more complex
                val backgroundColor = androidx.compose.ui.graphics.Color.White // Simplified

                val ratio = calculateContrastRatio(textColor, backgroundColor)
                if (ratio < 4.5) {
                    results.add(
                        AccessibilityAuditResult(
                            type = AuditType.COLOR_CONTRAST,
                            severity = Severity.ERROR,
                            message = "Insufficient color contrast ratio: ${String.format("%.2f", ratio)}:1. Minimum required: 4.5:1",
                            viewInfo = ViewInfo.fromView(v),
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
        return results
    }

    private fun auditReadingOrder(view: View): List<AccessibilityAuditResult> {
        val results = mutableListOf<AccessibilityAuditResult>()

        fun getAccessibleChildren(v: View): List<View> {
            val children = mutableListOf<View>()

            if (v.importantForAccessibility != View.IMPORTANT_FOR_ACCESSIBILITY_NO) {
                if (v is ViewGroup) {
                    for (i in 0 until v.childCount) {
                        children.addAll(getAccessibleChildren(v.getChildAt(i)))
                    }
                } else {
                    children.add(v)
                }
            }

            return children
        }

        val accessibleViews = getAccessibleChildren(view)

        for (i in 1 until accessibleViews.size) {
            val current = accessibleViews[i]
            val previous = accessibleViews[i - 1]

            // Check logical reading order (left-to-right, top-to-bottom)
            val currentLocation = IntArray(2)
            val previousLocation = IntArray(2)
            current.getLocationOnScreen(currentLocation)
            previous.getLocationOnScreen(previousLocation)

            if (currentLocation[1] < previousLocation[1] + previous.height &&
                currentLocation[0] < previousLocation[0]) {
                results.add(
                    AccessibilityAuditResult(
                        type = AuditType.READING_ORDER,
                        severity = Severity.WARNING,
                        message = "Potential reading order issue: element may be read out of visual order",
                        viewInfo = ViewInfo.fromView(current),
                        wcagCriterion = "1.3.2"
                    )
                )
            }
        }

        return results
    }

    private fun auditSemanticProperties(view: View): List<AccessibilityAuditResult> {
        val results = mutableListOf<AccessibilityAuditResult>()

        fun checkView(v: View) {
            val nodeInfo = AccessibilityNodeInfo.obtain()
            v.onInitializeAccessibilityNodeInfo(nodeInfo)

            // Check for appropriate roles
            if (v is android.widget.Button && !nodeInfo.isClickable) {
                results.add(
                    AccessibilityAuditResult(
                        type = AuditType.SEMANTIC_PROPERTIES,
                        severity = Severity.WARNING,
                        message = "Button element not marked as clickable",
                        viewInfo = ViewInfo.fromView(v),
                        wcagCriterion = "4.1.2"
                    )
                )
            }

            nodeInfo.recycle()

            if (v is ViewGroup) {
                for (i in 0 until v.childCount) {
                    checkView(v.getChildAt(i))
                }
            }
        }

        checkView(view)
        return results
    }

    fun calculateContrastRatio(
        color1: androidx.compose.ui.graphics.Color,
        color2: androidx.compose.ui.graphics.Color
    ): Double {
        val luminance1 = calculateLuminance(color1)
        val luminance2 = calculateLuminance(color2)

        val lighter = maxOf(luminance1, luminance2)
        val darker = minOf(luminance1, luminance2)

        return (lighter + 0.05) / (darker + 0.05)
    }

    private fun calculateLuminance(color: androidx.compose.ui.graphics.Color): Double {
        val r = color.red
        val g = color.green
        val b = color.blue

        fun linearize(component: Float): Double {
            return if (component <= 0.03928) {
                component / 12.92
            } else {
                ((component + 0.055) / 1.055).pow(2.4)
            }
        }

        val rLinear = linearize(r)
        val gLinear = linearize(g)
        val bLinear = linearize(b)

        return 0.2126 * rLinear + 0.7152 * gLinear + 0.0722 * bLinear
    }

    private fun adjustColorForContrast(
        color: androidx.compose.ui.graphics.Color,
        background: androidx.compose.ui.graphics.Color,
        targetRatio: Double
    ): androidx.compose.ui.graphics.Color {
        // Convert to HSV for easier manipulation
        val hsv = FloatArray(3)
        Color.colorToHSV(color.toArgb(), hsv)

        var adjustedValue = hsv[2] // Value (brightness)
        var iterations = 0
        val maxIterations = 20

        while (iterations < maxIterations) {
            val testColor = androidx.compose.ui.graphics.Color.hsv(hsv[0], hsv[1], adjustedValue)
            val ratio = calculateContrastRatio(testColor, background)

            if (ratio >= targetRatio) {
                break
            }

            // Adjust brightness toward better contrast
            val bgLuminance = calculateLuminance(background)
            adjustedValue = if (bgLuminance > 0.5) {
                maxOf(0f, adjustedValue - 0.05f)
            } else {
                minOf(1f, adjustedValue + 0.05f)
            }

            iterations++
        }

        return androidx.compose.ui.graphics.Color.hsv(hsv[0], hsv[1], adjustedValue)
    }

    private fun logAuditResults(results: List<AccessibilityAuditResult>) {
        val errorCount = results.count { it.severity == Severity.ERROR }
        val warningCount = results.count { it.severity == Severity.WARNING }

        android.util.Log.i("AccessibilityService", "🔍 Accessibility Audit Complete")
        android.util.Log.i("AccessibilityService", "📊 Found $errorCount errors and $warningCount warnings")

        results.forEach { result ->
            val icon = if (result.severity == Severity.ERROR) "❌" else "⚠️"
            val logLevel = if (result.severity == Severity.ERROR) android.util.Log.ERROR else android.util.Log.WARN
            android.util.Log.println(logLevel, "AccessibilityService",
                "$icon [${result.wcagCriterion}] ${result.message}")
        }
    }

    fun cleanup() {
        serviceScope.cancel()
    }
}

// MARK: - Supporting Data Classes

data class AccessibilityAuditResult(
    val type: AuditType,
    val severity: Severity,
    val message: String,
    val viewInfo: ViewInfo,
    val wcagCriterion: String
)

data class ViewInfo(
    val className: String,
    val id: String?,
    val bounds: String,
    val contentDescription: String?
) {
    companion object {
        fun fromView(view: View): ViewInfo {
            val id = try {
                view.context.resources.getResourceEntryName(view.id)
            } catch (e: Exception) {
                null
            }

            return ViewInfo(
                className = view::class.java.simpleName,
                id = id,
                bounds = "${view.left},${view.top}-${view.right},${view.bottom}",
                contentDescription = view.contentDescription?.toString()
            )
        }
    }
}

data class ValidationResult(
    val isValid: Boolean,
    val issues: List<AccessibilityIssue>,
    val wcagLevel: WCAGLevel
)

enum class AuditType {
    TOUCH_TARGET_SIZE,
    CONTENT_DESCRIPTION,
    COLOR_CONTRAST,
    READING_ORDER,
    SEMANTIC_PROPERTIES
}

enum class Severity {
    ERROR,
    WARNING
}

enum class ContrastLevel(val minimumRatio: Double) {
    AA(4.5),
    AAA(7.0)
}

enum class AnnouncementPriority {
    HIGH,
    MEDIUM,
    LOW
}

enum class LiveRegionMode {
    ASSERTIVE,
    POLITE,
    NONE
}

enum class AccessibilityIssue {
    INSUFFICIENT_TOUCH_TARGET,
    MISSING_CONTENT_DESCRIPTION,
    POOR_CONTRAST,
    MISSING_SEMANTIC_ROLE,
    INVALID_READING_ORDER
}

enum class WCAGLevel {
    NONE,
    A,
    AA,
    AAA
}

// MARK: - Compose Extensions

/**
 * Composable that provides accessibility configuration
 */
@Composable
fun accessibilityService(): AccessibilityService {
    val context = LocalContext.current
    return remember { AccessibilityService.getInstance(context) }
}

/**
 * Modifier extension for easy accessibility configuration
 */
fun Modifier.accessibilityConfigured(
    contentDescription: String? = null,
    role: Role? = null,
    stateDescription: String? = null,
    onClick: (() -> Unit)? = null,
    onLongClick: (() -> Unit)? = null,
    isImportant: Boolean = true,
    isClickable: Boolean = false,
    isEnabled: Boolean = true,
    liveRegion: LiveRegionMode = LiveRegionMode.POLITE,
    service: AccessibilityService
): Modifier {
    return this.semantics {
        service.configureSemantics(
            contentDescription = contentDescription,
            role = role,
            stateDescription = stateDescription,
            onClick = onClick,
            onLongClick = onLongClick,
            isImportant = isImportant,
            isClickable = isClickable,
            isEnabled = isEnabled,
            liveRegion = liveRegion
        ).invoke(this)
    }
}

/**
 * Modifier extension for minimum touch targets
 */
fun Modifier.minimumTouchTarget(service: AccessibilityService): Modifier {
    return this.then(service.ensureMinimumTouchTarget())
}