package com.tchat.models

import android.os.Build
import kotlinx.serialization.Serializable
import java.util.UUID

/**
 * Platform adapter entity for handling platform-specific UI conventions and capabilities
 */
@Serializable
data class PlatformAdapter(
    val id: String = UUID.randomUUID().toString(),
    val platform: String = "android",
    val version: String = Build.VERSION.RELEASE,
    val capabilities: List<PlatformCapability> = emptyList(),
    val uiConventions: UIConventions = UIConventions.defaultAndroidConventions,
    val gestureSupport: GestureSupport = GestureSupport.defaultAndroidGestures,
    val animationSupport: AnimationSupport = AnimationSupport.defaultAndroidAnimations,
    val metadata: PlatformAdapterMetadata = PlatformAdapterMetadata()
) {

    /**
     * Check if platform supports specific capability
     */
    fun supportsCapability(capabilityName: String): Boolean {
        return capabilities.any { it.name == capabilityName && it.isSupported }
    }

    /**
     * Get capability restrictions
     */
    fun getCapabilityRestrictions(capabilityName: String): List<String> {
        return capabilities.find { it.name == capabilityName }?.restrictions ?: emptyList()
    }

    /**
     * Check if gesture is supported
     */
    fun supportsGesture(gestureType: String): Boolean {
        return gestureSupport.supportedGestures.any { it.type == gestureType }
    }

    /**
     * Check if animation is supported
     */
    fun supportsAnimation(animationType: String): Boolean {
        return animationSupport.supportedAnimations.any { it.type == animationType }
    }

    /**
     * Get platform-specific color scheme
     */
    val colorScheme: Map<String, String>
        get() = uiConventions.designSystem.colorScheme

    /**
     * Get platform-specific spacing
     */
    val spacing: Map<String, Double>
        get() = uiConventions.designSystem.spacing

    /**
     * Get platform-specific border radius
     */
    val borderRadius: Map<String, Double>
        get() = uiConventions.designSystem.borderRadius

    companion object {
        /**
         * Create default Android platform adapter
         */
        fun defaultAndroidAdapter(): PlatformAdapter {
            val capabilities = listOf(
                PlatformCapability(
                    name = "hapticFeedback",
                    isSupported = Build.VERSION.SDK_INT >= Build.VERSION_CODES.O,
                    apiLevel = "API 26+",
                    description = "Haptic feedback support"
                ),
                PlatformCapability(
                    name = "biometric",
                    isSupported = Build.VERSION.SDK_INT >= Build.VERSION_CODES.M,
                    apiLevel = "API 23+",
                    description = "Biometric authentication"
                ),
                PlatformCapability(
                    name = "darkMode",
                    isSupported = Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q,
                    apiLevel = "API 29+",
                    description = "Dark mode support"
                ),
                PlatformCapability(
                    name = "adaptiveIcon",
                    isSupported = Build.VERSION.SDK_INT >= Build.VERSION_CODES.O,
                    apiLevel = "API 26+",
                    description = "Adaptive icon support"
                ),
                PlatformCapability(
                    name = "notificationChannels",
                    isSupported = Build.VERSION.SDK_INT >= Build.VERSION_CODES.O,
                    apiLevel = "API 26+",
                    description = "Notification channels"
                )
            )

            return PlatformAdapter(
                platform = "android",
                version = Build.VERSION.RELEASE,
                capabilities = capabilities,
                uiConventions = UIConventions.defaultAndroidConventions,
                gestureSupport = GestureSupport.defaultAndroidGestures,
                animationSupport = AnimationSupport.defaultAndroidAnimations,
                metadata = PlatformAdapterMetadata()
            )
        }
    }
}

/**
 * Platform capability definition
 */
@Serializable
data class PlatformCapability(
    val id: String = UUID.randomUUID().toString(),
    val name: String,
    val isSupported: Boolean,
    val apiLevel: String,
    val restrictions: List<String> = emptyList(),
    val alternativeActions: List<String> = emptyList(),
    val description: String? = null
)

/**
 * UI conventions for platform-specific design patterns
 */
@Serializable
data class UIConventions(
    val designSystem: DesignSystem,
    val navigationPatterns: Map<String, String> = emptyMap(),
    val layoutPatterns: Map<String, String> = emptyMap(),
    val accessibilityGuidelines: Map<String, String> = emptyMap()
) {
    companion object {
        /**
         * Default Android UI conventions
         */
        val defaultAndroidConventions = UIConventions(
            designSystem = DesignSystem.defaultAndroidDesignSystem,
            navigationPatterns = mapOf(
                "bottomNavigation" to "material",
                "drawer" to "left",
                "appBar" to "top",
                "fab" to "bottomRight"
            ),
            layoutPatterns = mapOf(
                "statusBar" to "immersive",
                "navigationBar" to "gestural",
                "notch" to "adapt"
            ),
            accessibilityGuidelines = mapOf(
                "talkBack" to "supported",
                "largeText" to "supported",
                "highContrast" to "supported"
            )
        )
    }
}

/**
 * Design system definition
 */
@Serializable
data class DesignSystem(
    val colorScheme: Map<String, String>,
    val typography: Map<String, String> = emptyMap(),
    val spacing: Map<String, Double> = emptyMap(),
    val borderRadius: Map<String, Double> = emptyMap(),
    val shadows: Map<String, String> = emptyMap()
) {
    companion object {
        /**
         * Default Android design system
         */
        val defaultAndroidDesignSystem = DesignSystem(
            colorScheme = mapOf(
                "primary" to "#6200EE",
                "primaryVariant" to "#3700B3",
                "secondary" to "#03DAC6",
                "secondaryVariant" to "#018786",
                "background" to "#FFFFFF",
                "surface" to "#FFFFFF",
                "onPrimary" to "#FFFFFF",
                "onSecondary" to "#000000",
                "onBackground" to "#000000",
                "onSurface" to "#000000"
            ),
            typography = mapOf(
                "h1" to "96sp",
                "h2" to "60sp",
                "h3" to "48sp",
                "h4" to "34sp",
                "h5" to "24sp",
                "h6" to "20sp",
                "subtitle1" to "16sp",
                "subtitle2" to "14sp",
                "body1" to "16sp",
                "body2" to "14sp",
                "button" to "14sp",
                "caption" to "12sp",
                "overline" to "10sp"
            ),
            spacing = mapOf(
                "xs" to 4.0,
                "sm" to 8.0,
                "md" to 16.0,
                "lg" to 24.0,
                "xl" to 32.0,
                "xxl" to 48.0
            ),
            borderRadius = mapOf(
                "xs" to 2.0,
                "sm" to 4.0,
                "md" to 8.0,
                "lg" to 12.0,
                "xl" to 16.0,
                "round" to 50.0
            ),
            shadows = mapOf(
                "elevation1" to "0 1 3 rgba(0,0,0,0.12)",
                "elevation2" to "0 1 5 rgba(0,0,0,0.12)",
                "elevation3" to "0 1 8 rgba(0,0,0,0.12)",
                "elevation4" to "0 2 10 rgba(0,0,0,0.12)",
                "elevation5" to "0 4 15 rgba(0,0,0,0.12)"
            )
        )
    }
}

/**
 * Gesture support definition
 */
@Serializable
data class GestureSupport(
    val supportedGestures: List<GestureDefinition>,
    val maxSimultaneousGestures: Int = 2,
    val gestureSettings: Map<String, String> = emptyMap()
) {
    companion object {
        /**
         * Default Android gesture support
         */
        val defaultAndroidGestures = GestureSupport(
            supportedGestures = listOf(
                GestureDefinition(type = "tap", description = "Single tap"),
                GestureDefinition(type = "doubleTap", description = "Double tap"),
                GestureDefinition(type = "longPress", description = "Long press"),
                GestureDefinition(type = "swipe", description = "Swipe gesture"),
                GestureDefinition(type = "pan", description = "Pan/drag gesture"),
                GestureDefinition(type = "pinch", description = "Pinch to zoom"),
                GestureDefinition(type = "rotation", description = "Rotation gesture"),
                GestureDefinition(type = "fling", description = "Fling gesture")
            ),
            maxSimultaneousGestures = 10,
            gestureSettings = mapOf(
                "hapticFeedback" to "enabled",
                "edgeSwipe" to "enabled",
                "backGesture" to "enabled"
            )
        )
    }
}

/**
 * Gesture definition
 */
@Serializable
data class GestureDefinition(
    val type: String,
    val description: String,
    val requiredFingers: Int = 1,
    val minDuration: Double = 0.0,
    val maxDuration: Double? = null
)

/**
 * Animation support definition
 */
@Serializable
data class AnimationSupport(
    val supportedAnimations: List<AnimationDefinition>,
    val defaultDuration: Double = 300.0,
    val animationSettings: Map<String, String> = emptyMap()
) {
    companion object {
        /**
         * Default Android animation support
         */
        val defaultAndroidAnimations = AnimationSupport(
            supportedAnimations = listOf(
                AnimationDefinition(type = "fade", description = "Fade in/out"),
                AnimationDefinition(type = "slide", description = "Slide transition"),
                AnimationDefinition(type = "scale", description = "Scale transform"),
                AnimationDefinition(type = "rotate", description = "Rotation transform"),
                AnimationDefinition(type = "translate", description = "Translation"),
                AnimationDefinition(type = "morphing", description = "Morphing transition"),
                AnimationDefinition(type = "ripple", description = "Material ripple effect")
            ),
            defaultDuration = 300.0,
            animationSettings = mapOf(
                "preferredFrameRate" to "60",
                "reducedMotion" to "respected",
                "timing" to "standardEasing",
                "interpolator" to "fastOutSlowIn"
            )
        )
    }
}

/**
 * Animation definition
 */
@Serializable
data class AnimationDefinition(
    val type: String,
    val description: String,
    val minDuration: Double = 100.0,
    val maxDuration: Double = 2000.0,
    val defaultEasing: String = "fastOutSlowIn"
)

/**
 * Platform adapter metadata
 */
@Serializable
data class PlatformAdapterMetadata(
    val deviceModel: String = Build.MODEL,
    val manufacturer: String = Build.MANUFACTURER,
    val osVersion: String = Build.VERSION.RELEASE,
    val apiLevel: Int = Build.VERSION.SDK_INT,
    val screenDensity: String = "unknown", // Would need context to determine
    val colorDepth: String = "24bit",
    val refreshRate: String = "60hz",
    val hasNotch: Boolean = false, // Simplified detection
    val supportsDarkMode: Boolean = Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q,
    val accessibility: Map<String, Boolean> = emptyMap()
)