package com.tchat.designsystem

import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

/**
 * Spacing design tokens matching TailwindCSS v4 spacing scale
 */
object Spacing {

    // Base Spacing Scale
    // Following TailwindCSS 4px base unit system

    /** Extra extra small spacing (2dp)
     *  Maps to TailwindCSS: space-0.5 (0.125rem) */
    val xxs: Dp = 2.dp

    /** Extra small spacing (4dp)
     *  Maps to TailwindCSS: space-1 (0.25rem) */
    val xs: Dp = 4.dp

    /** Small spacing (8dp)
     *  Maps to TailwindCSS: space-2 (0.5rem) */
    val sm: Dp = 8.dp

    /** Medium spacing (16dp) - Most common spacing
     *  Maps to TailwindCSS: space-4 (1rem) */
    val md: Dp = 16.dp

    /** Large spacing (24dp)
     *  Maps to TailwindCSS: space-6 (1.5rem) */
    val lg: Dp = 24.dp

    /** Extra large spacing (32dp)
     *  Maps to TailwindCSS: space-8 (2rem) */
    val xl: Dp = 32.dp

    /** Extra extra large spacing (48dp)
     *  Maps to TailwindCSS: space-12 (3rem) */
    val xxl: Dp = 48.dp

    /** Extra extra extra large spacing (64dp)
     *  Maps to TailwindCSS: space-16 (4rem) */
    val xxxl: Dp = 64.dp

    // Semantic Spacing

    /** Component padding values */
    object Component {
        /** Button padding */
        val buttonPadding = PaddingValues(horizontal = md, vertical = sm)

        /** Card padding */
        val cardPadding = PaddingValues(md)

        /** Input field padding */
        val inputPadding = PaddingValues(horizontal = sm, vertical = xs)

        /** Modal padding */
        val modalPadding = PaddingValues(lg)

        /** Chip padding */
        val chipPadding = PaddingValues(horizontal = sm, vertical = xs)

        /** Dialog padding */
        val dialogPadding = PaddingValues(lg)
    }

    /** Layout spacing values */
    object Layout {
        /** Screen edge margins */
        val screenMargin = md

        /** Section spacing */
        val sectionSpacing = lg

        /** Content spacing */
        val contentSpacing = md

        /** List item spacing */
        val listItemSpacing = sm

        /** Grid gap */
        val gridGap = md

        /** Column gap */
        val columnGap = md

        /** Row gap */
        val rowGap = sm
    }

    /** Navigation spacing values */
    object Navigation {
        /** Bottom navigation height */
        val bottomNavigationHeight = 56.dp // Material Design standard

        /** Top app bar height */
        val topAppBarHeight = 64.dp // Material Design standard

        /** Tab height */
        val tabHeight = 48.dp

        /** Navigation item spacing */
        val navigationItemSpacing = xs

        /** Navigation padding */
        val navigationPadding = md

        /** Fab margin */
        val fabMargin = md
    }

    /** Interactive spacing values */
    object Interactive {
        /** Minimum touch target size (Material Design) */
        val touchTargetSize = 48.dp

        /** Button spacing */
        val buttonSpacing = sm

        /** Form field spacing */
        val formFieldSpacing = md

        /** Icon spacing */
        val iconSpacing = xs

        /** Action spacing */
        val actionSpacing = sm

        /** Checkbox spacing */
        val checkboxSpacing = sm
    }

    /** Elevation values for Material Design */
    object Elevation {
        val none = 0.dp
        val small = 2.dp
        val medium = 4.dp
        val large = 8.dp
        val extraLarge = 16.dp
    }
}

// Responsive Spacing

/**
 * Get spacing that adapts to screen size
 */
fun adaptiveSpacing(compact: Dp, regular: Dp): Dp {
    // In a full implementation, this would check window size class
    return regular // Placeholder
}

/**
 * Scale spacing for accessibility
 */
fun scaledSpacing(spacing: Dp, scaleFactor: Float = 1.0f): Dp {
    return spacing * scaleFactor
}

// Padding Extensions

/**
 * Helper object for creating padding values
 */
object PaddingHelper {
    /**
     * Create symmetric padding values
     */
    fun symmetric(horizontal: Dp = 0.dp, vertical: Dp = 0.dp): PaddingValues {
        return PaddingValues(horizontal = horizontal, vertical = vertical)
    }

    /**
     * Create uniform padding values
     */
    fun uniform(value: Dp): PaddingValues {
        return PaddingValues(value)
    }

    /**
     * Common padding values using design tokens
     */
    val small: PaddingValues
        get() = PaddingValues(Spacing.sm)

    val medium: PaddingValues
        get() = PaddingValues(Spacing.md)

    val large: PaddingValues
        get() = PaddingValues(Spacing.lg)
}

// Spacing Utilities

/**
 * Semantic spacing usage enum
 */
enum class SpacingUsage {
    ELEMENT_SPACING,     // Between small elements
    COMPONENT_SPACING,   // Between components
    SECTION_SPACING,     // Between sections
    PAGE_SPACING,        // Page margins
    SCREEN_MARGIN,       // Screen edge margins
    FORM_FIELD_SPACING,  // Between form fields
    LIST_ITEM_SPACING,   // Between list items
    BUTTON_SPACING,      // Between buttons
    ICON_PADDING,        // Around icons
    CARD_PADDING,        // Inside cards
    DIALOG_PADDING,      // Inside dialogs
    BOTTOM_SHEET_PADDING // Inside bottom sheets
}

/**
 * Get spacing for semantic usage
 */
fun Spacing.spacingFor(usage: SpacingUsage): Dp {
    return when (usage) {
        SpacingUsage.ELEMENT_SPACING -> xs
        SpacingUsage.COMPONENT_SPACING -> sm
        SpacingUsage.SECTION_SPACING -> md
        SpacingUsage.PAGE_SPACING -> lg
        SpacingUsage.SCREEN_MARGIN -> md
        SpacingUsage.FORM_FIELD_SPACING -> md
        SpacingUsage.LIST_ITEM_SPACING -> sm
        SpacingUsage.BUTTON_SPACING -> sm
        SpacingUsage.ICON_PADDING -> xs
        SpacingUsage.CARD_PADDING -> md
        SpacingUsage.DIALOG_PADDING -> lg
        SpacingUsage.BOTTOM_SHEET_PADDING -> md
    }
}

/**
 * Validate spacing consistency (multiples of 4dp)
 */
fun isValidSpacing(spacing: Dp): Boolean {
    val value = spacing.value
    return (value % 4f) == 0f
}

/**
 * Component padding types
 */
enum class ComponentPaddingType {
    BUTTON,
    CARD,
    INPUT,
    MODAL,
    CHIP,
    DIALOG,
    BOTTOM_SHEET,
    LIST_ITEM
}

/**
 * Get padding for component type
 */
fun Spacing.Component.paddingFor(type: ComponentPaddingType): PaddingValues {
    return when (type) {
        ComponentPaddingType.BUTTON -> buttonPadding
        ComponentPaddingType.CARD -> cardPadding
        ComponentPaddingType.INPUT -> inputPadding
        ComponentPaddingType.MODAL -> modalPadding
        ComponentPaddingType.CHIP -> chipPadding
        ComponentPaddingType.DIALOG -> dialogPadding
        ComponentPaddingType.BOTTOM_SHEET -> modalPadding
        ComponentPaddingType.LIST_ITEM -> PaddingValues(horizontal = Spacing.md, vertical = Spacing.sm)
    }
}