/**
 * TchatAccessibilityManager.kt
 * TchatApp
 *
 * Created by Claude on 22/09/2024.
 */

package com.tchat.services

import android.accessibilityservice.AccessibilityServiceInfo
import android.content.Context
import android.content.res.Configuration
import android.graphics.Color
import android.os.Build
import android.provider.Settings
import android.util.Log
import android.view.View
import android.view.accessibility.AccessibilityManager
import android.view.accessibility.AccessibilityNodeInfo
import androidx.compose.foundation.layout.size
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.semantics.*
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlin.math.max
import kotlin.math.min
import kotlin.math.pow

/**
 * Accessibility compliance manager for Android
 * Implements T054: Accessibility compliance for TalkBack, Large Text, and WCAG guidelines
 */
class TchatAccessibilityManager private constructor(
    private val context: Context
) {

    // MARK: - Types

    enum class ComplianceLevel(val value: String) {
        A("A"),
        AA("AA"),
        AAA("AAA");

        val description: String
            get() = when (this) {
                A -> "WCAG 2.1 Level A (Basic)"
                AA -> "WCAG 2.1 Level AA (Standard)"
                AAA -> "WCAG 2.1 Level AAA (Enhanced)"
            }
    }

    data class AccessibilityValidationResult(
        val passes: Boolean,
        val actualRatio: Double,
        val requiredRatio: Double,
        val complianceLevel: ComplianceLevel,
        val message: String
    )

    data class AccessibilityRequirements(
        val colorContrastLevel: ComplianceLevel = ComplianceLevel.AA,
        val minimumTouchTargetSize: Float = 48f, // dp
        val requiresLabel: Boolean = true
    ) {
        companion object {
            val STANDARD = AccessibilityRequirements()
            val ENHANCED = AccessibilityRequirements(
                colorContrastLevel = ComplianceLevel.AAA,
                minimumTouchTargetSize = 56f
            )
        }
    }

    sealed class AccessibilityIssue(val description: String) {
        class MissingLabel(message: String) : AccessibilityIssue(message)
        class InsufficientTouchTarget(message: String) : AccessibilityIssue(message)
        class InsufficientContrast(message: String) : AccessibilityIssue(message)
        class MissingHint(message: String) : AccessibilityIssue(message)
        class IncorrectRole(message: String) : AccessibilityIssue(message)
    }

    // MARK: - Color Contrast Standards

    object ColorContrastRequirements {
        const val NORMAL_TEXT_AA = 4.5
        const val LARGE_TEXT_AA = 3.0
        const val NORMAL_TEXT_AAA = 7.0
        const val LARGE_TEXT_AAA = 4.5
        const val UI_COMPONENTS_AA = 3.0
    }

    // MARK: - State Properties

    private val accessibilityManager = context.getSystemService(Context.ACCESSIBILITY_SERVICE) as android.view.accessibility.AccessibilityManager

    private val _isTalkBackEnabled = MutableStateFlow(false)
    val isTalkBackEnabled: StateFlow<Boolean> = _isTalkBackEnabled.asStateFlow()

    private val _fontScale = MutableStateFlow(1f)
    val fontScale: StateFlow<Float> = _fontScale.asStateFlow()

    private val _isHighContrastEnabled = MutableStateFlow(false)
    val isHighContrastEnabled: StateFlow<Boolean> = _isHighContrastEnabled.asStateFlow()

    private val _isReduceMotionEnabled = MutableStateFlow(false)
    val isReduceMotionEnabled: StateFlow<Boolean> = _isReduceMotionEnabled.asStateFlow()

    companion object {
        private const val TAG = "AccessibilityManager"

        @Volatile
        private var INSTANCE: TchatAccessibilityManager? = null

        fun getInstance(context: Context): TchatAccessibilityManager {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: TchatAccessibilityManager(context.applicationContext).also { INSTANCE = it }
            }
        }
    }

    // MARK: - Public Interface

    /**
     * Starts monitoring accessibility changes
     */
    fun startMonitoring() {
        updateAccessibilitySettings()
        setupAccessibilityListener()
    }

    /**
     * Stops monitoring accessibility changes
     */
    fun stopMonitoring() {
        // No specific cleanup needed for Android accessibility monitoring
    }

    /**
     * Validates color contrast ratio according to WCAG guidelines
     */
    fun validateColorContrast(
        foregroundColor: Int,
        backgroundColor: Int,
        fontSize: Float, // sp
        complianceLevel: ComplianceLevel = ComplianceLevel.AA
    ): AccessibilityValidationResult {

        val contrastRatio = calculateContrastRatio(foregroundColor, backgroundColor)
        val isLargeText = fontSize >= 18f || (fontSize >= 14f && isBoldText(fontSize))

        val requiredRatio = when (complianceLevel) {
            ComplianceLevel.A -> 3.0 // Basic requirement
            ComplianceLevel.AA -> if (isLargeText) ColorContrastRequirements.LARGE_TEXT_AA else ColorContrastRequirements.NORMAL_TEXT_AA
            ComplianceLevel.AAA -> if (isLargeText) ColorContrastRequirements.LARGE_TEXT_AAA else ColorContrastRequirements.NORMAL_TEXT_AAA
        }

        val passes = contrastRatio >= requiredRatio

        return AccessibilityValidationResult(
            passes = passes,
            actualRatio = contrastRatio,
            requiredRatio = requiredRatio,
            complianceLevel = complianceLevel,
            message = if (passes) {
                "Color contrast meets ${complianceLevel.description}"
            } else {
                "Color contrast (${String.format("%.2f", contrastRatio)}:1) below required ${String.format("%.1f", requiredRatio)}:1 for ${complianceLevel.description}"
            }
        )
    }

    /**
     * Configures accessibility for Android Views
     */
    fun configureAccessibility(
        view: View,
        contentDescription: String,
        hint: String? = null,
        role: AccessibilityNodeInfo.AccessibilityAction? = null
    ) {
        view.contentDescription = contentDescription
        view.importantForAccessibility = View.IMPORTANT_FOR_ACCESSIBILITY_YES

        hint?.let {
            if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
                view.tooltipText = it
            }
        }

        // Configure accessibility actions if provided
        role?.let { action ->
            view.setAccessibilityDelegate(object : View.AccessibilityDelegate() {
                override fun onInitializeAccessibilityNodeInfo(host: View, info: AccessibilityNodeInfo) {
                    super.onInitializeAccessibilityNodeInfo(host, info)
                    info.addAction(action)
                }
            })
        }
    }

    /**
     * Announces accessibility changes to TalkBack
     */
    fun announceChange(message: String) {
        if (_isTalkBackEnabled.value) {
            // Create a temporary view to announce the message
            val view = View(context)
            view.announceForAccessibility(message)
        }
    }

    /**
     * Gets scaled font size based on system font scale
     */
    fun getScaledTextSize(baseSize: Float): Float {
        return baseSize * _fontScale.value
    }

    /**
     * Checks if device settings require reduced motion
     */
    val shouldReduceMotion: Boolean
        get() = _isReduceMotionEnabled.value

    /**
     * Gets appropriate animation duration based on accessibility settings
     */
    fun getAnimationDuration(defaultDuration: Long): Long {
        return if (shouldReduceMotion) 0L else defaultDuration
    }

    /**
     * Validates component accessibility compliance
     */
    fun validateComponentAccessibility(
        view: View,
        requirements: AccessibilityRequirements = AccessibilityRequirements.STANDARD
    ): List<AccessibilityIssue> {
        val issues = mutableListOf<AccessibilityIssue>()

        // Check for content description
        if (requirements.requiresLabel && view.contentDescription.isNullOrEmpty()) {
            issues.add(AccessibilityIssue.MissingLabel("Component missing content description"))
        }

        // Check minimum touch target size (convert dp to pixels)
        val density = context.resources.displayMetrics.density
        val minSizePx = requirements.minimumTouchTargetSize * density

        if (view.width < minSizePx || view.height < minSizePx) {
            issues.add(AccessibilityIssue.InsufficientTouchTarget(
                "Touch target smaller than ${requirements.minimumTouchTargetSize}dp"
            ))
        }

        return issues
    }

    /**
     * Provides minimum touch target size for accessibility compliance
     */
    fun getMinimumTouchTarget(): Dp = 48.dp

    // MARK: - Private Methods

    private fun setupAccessibilityListener() {
        accessibilityManager.addAccessibilityStateChangeListener { enabled ->
            updateAccessibilitySettings()
        }

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            accessibilityManager.addTouchExplorationStateChangeListener { enabled ->
                updateAccessibilitySettings()
            }
        }
    }

    private fun updateAccessibilitySettings() {
        updateTalkBackStatus()
        updateFontScale()
        updateHighContrastStatus()
        updateReduceMotionStatus()
    }

    private fun updateTalkBackStatus() {
        val enabledServices = accessibilityManager.getEnabledAccessibilityServiceList(AccessibilityServiceInfo.FEEDBACK_SPOKEN)
        _isTalkBackEnabled.value = enabledServices.isNotEmpty()
    }

    private fun updateFontScale() {
        _fontScale.value = context.resources.configuration.fontScale
    }

    private fun updateHighContrastStatus() {
        _isHighContrastEnabled.value = try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1) {
                Settings.Secure.getFloat(
                    context.contentResolver,
                    "high_text_contrast_enabled",
                    0f
                ) == 1f
            } else {
                false
            }
        } catch (e: Settings.SettingNotFoundException) {
            false
        }
    }

    private fun updateReduceMotionStatus() {
        _isReduceMotionEnabled.value = try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1) {
                Settings.Global.getFloat(
                    context.contentResolver,
                    Settings.Global.ANIMATOR_DURATION_SCALE,
                    1f
                ) == 0f
            } else {
                false
            }
        } catch (e: Settings.SettingNotFoundException) {
            false
        }
    }

    // MARK: - Color Contrast Calculations

    private fun calculateContrastRatio(foregroundColor: Int, backgroundColor: Int): Double {
        val foregroundLuminance = relativeLuminance(foregroundColor)
        val backgroundLuminance = relativeLuminance(backgroundColor)

        val lighter = max(foregroundLuminance, backgroundLuminance)
        val darker = min(foregroundLuminance, backgroundLuminance)

        return (lighter + 0.05) / (darker + 0.05)
    }

    private fun relativeLuminance(color: Int): Double {
        val red = Color.red(color) / 255.0
        val green = Color.green(color) / 255.0
        val blue = Color.blue(color) / 255.0

        // Convert to linear RGB
        val linearRed = linearizeColorComponent(red)
        val linearGreen = linearizeColorComponent(green)
        val linearBlue = linearizeColorComponent(blue)

        // Calculate relative luminance using ITU-R BT.709 coefficients
        return 0.2126 * linearRed + 0.7152 * linearGreen + 0.0722 * linearBlue
    }

    private fun linearizeColorComponent(component: Double): Double {
        return if (component <= 0.03928) {
            component / 12.92
        } else {
            ((component + 0.055) / 1.055).pow(2.4)
        }
    }

    // MARK: - Helper Methods

    private fun isBoldText(fontSize: Float): Boolean {
        // Simplified bold text detection
        return fontSize >= 14f
    }
}

// MARK: - Compose Extensions

/**
 * Modifier for accessibility compliance in Jetpack Compose
 */
fun Modifier.accessibilityCompliance(
    contentDescription: String,
    role: Role? = null,
    stateDescription: String? = null,
    onClick: (() -> Unit)? = null
): Modifier = this.semantics {
    this.contentDescription = contentDescription
    role?.let { this.role = it }
    stateDescription?.let { this.stateDescription = it }
    onClick?.let { this.onClick { it(); true } }
}

/**
 * Modifier for minimum touch target size
 */
fun Modifier.minimumTouchTarget(size: Dp = 48.dp): Modifier = this.size(size)

/**
 * Get scaled text size for accessibility
 */
@Composable
fun getAccessibleTextSize(baseSize: Float): androidx.compose.ui.unit.TextUnit {
    val context = LocalContext.current
    val configuration = LocalConfiguration.current
    val scaledSize = baseSize * configuration.fontScale
    return scaledSize.sp
}

/**
 * Check if high contrast is enabled
 */
@Composable
fun isHighContrastEnabled(): Boolean {
    val context = LocalContext.current
    return try {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1) {
            Settings.Secure.getFloat(
                context.contentResolver,
                "high_text_contrast_enabled",
                0f
            ) == 1f
        } else {
            false
        }
    } catch (e: Settings.SettingNotFoundException) {
        false
    }
}

/**
 * Accessibility-aware color selection
 */
@Composable
fun accessibleColor(
    normal: androidx.compose.ui.graphics.Color,
    highContrast: androidx.compose.ui.graphics.Color
): androidx.compose.ui.graphics.Color {
    return if (isHighContrastEnabled()) highContrast else normal
}