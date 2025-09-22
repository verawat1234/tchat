package com.tchat.designsystem

import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

/**
 * Typography design tokens matching TailwindCSS v4 text scales
 */
object Typography {

    // Heading Styles

    /** Large heading style (32sp, bold)
     *  Maps to TailwindCSS: text-3xl font-bold */
    val headingLarge = TextStyle(
        fontSize = 32.sp,
        fontWeight = FontWeight.Bold,
        lineHeight = 40.sp,
        letterSpacing = (-0.25).sp
    )

    /** Medium heading style (24sp, semibold)
     *  Maps to TailwindCSS: text-2xl font-semibold */
    val headingMedium = TextStyle(
        fontSize = 24.sp,
        fontWeight = FontWeight.SemiBold,
        lineHeight = 32.sp,
        letterSpacing = 0.sp
    )

    /** Small heading style (20sp, semibold)
     *  Maps to TailwindCSS: text-xl font-semibold */
    val headingSmall = TextStyle(
        fontSize = 20.sp,
        fontWeight = FontWeight.SemiBold,
        lineHeight = 28.sp,
        letterSpacing = 0.sp
    )

    // Body Styles

    /** Large body text (18sp, regular)
     *  Maps to TailwindCSS: text-lg */
    val bodyLarge = TextStyle(
        fontSize = 18.sp,
        fontWeight = FontWeight.Normal,
        lineHeight = 28.sp,
        letterSpacing = 0.sp
    )

    /** Medium body text (16sp, regular) - Default body text
     *  Maps to TailwindCSS: text-base */
    val bodyMedium = TextStyle(
        fontSize = 16.sp,
        fontWeight = FontWeight.Normal,
        lineHeight = 24.sp,
        letterSpacing = 0.sp
    )

    /** Small body text (14sp, regular)
     *  Maps to TailwindCSS: text-sm */
    val bodySmall = TextStyle(
        fontSize = 14.sp,
        fontWeight = FontWeight.Normal,
        lineHeight = 20.sp,
        letterSpacing = 0.sp
    )

    // Label Styles

    /** Caption text (12sp, medium)
     *  Maps to TailwindCSS: text-xs font-medium */
    val caption = TextStyle(
        fontSize = 12.sp,
        fontWeight = FontWeight.Medium,
        lineHeight = 16.sp,
        letterSpacing = 0.4.sp
    )

    /** Label text (14sp, medium)
     *  Maps to TailwindCSS: text-sm font-medium */
    val label = TextStyle(
        fontSize = 14.sp,
        fontWeight = FontWeight.Medium,
        lineHeight = 20.sp,
        letterSpacing = 0.1.sp
    )

    /** Button text (16sp, semibold)
     *  Maps to TailwindCSS: text-base font-semibold */
    val button = TextStyle(
        fontSize = 16.sp,
        fontWeight = FontWeight.SemiBold,
        lineHeight = 24.sp,
        letterSpacing = 0.sp
    )

    // Utility Styles

    /** Overline text (10sp, bold, uppercase)
     *  Maps to TailwindCSS: text-xs font-bold uppercase tracking-wide */
    val overline = TextStyle(
        fontSize = 10.sp,
        fontWeight = FontWeight.Bold,
        lineHeight = 16.sp,
        letterSpacing = 1.5.sp
    )

    /** Monospace text for code (14sp, regular, monospaced)
     *  Maps to TailwindCSS: text-sm font-mono */
    val code = TextStyle(
        fontSize = 14.sp,
        fontWeight = FontWeight.Normal,
        lineHeight = 20.sp,
        letterSpacing = 0.sp,
        fontFamily = FontFamily.Monospace
    )

    // Line Height Constants

    /** Line height multipliers matching TailwindCSS leading values */
    object LineHeight {
        const val TIGHT = 1.25f     // leading-tight
        const val NORMAL = 1.5f     // leading-normal
        const val RELAXED = 1.625f  // leading-relaxed
        const val LOOSE = 2.0f      // leading-loose
    }

    // Letter Spacing Constants

    /** Letter spacing values matching TailwindCSS tracking values */
    object LetterSpacing {
        val TIGHTER = (-0.8).sp  // tracking-tighter
        val TIGHT = (-0.4).sp    // tracking-tight
        val NORMAL = 0.sp        // tracking-normal
        val WIDE = 0.4.sp        // tracking-wide
        val WIDER = 0.8.sp       // tracking-wider
        val WIDEST = 1.6.sp      // tracking-widest
    }
}

// Font Weight Extensions

/** TailwindCSS font weight mapping */
object FontWeights {
    val thin = FontWeight.Thin           // font-thin
    val extraLight = FontWeight.ExtraLight  // font-extralight
    val light = FontWeight.Light         // font-light
    val normal = FontWeight.Normal       // font-normal
    val medium = FontWeight.Medium       // font-medium
    val semiBold = FontWeight.SemiBold   // font-semibold
    val bold = FontWeight.Bold           // font-bold
    val extraBold = FontWeight.ExtraBold // font-extrabold
    val black = FontWeight.Black         // font-black
}

// Typography Utilities

/**
 * Semantic typography usage enum
 */
enum class TypographyUsage {
    PAGE_TITLE,
    SECTION_TITLE,
    CARD_TITLE,
    BODY_TEXT,
    SECONDARY_TEXT,
    CAPTION_TEXT,
    BUTTON_LABEL,
    NAVIGATION_TITLE,
    TAB_LABEL,
    INPUT_LABEL,
    ERROR_TEXT
}

/**
 * Get font for semantic usage
 */
fun Typography.fontFor(usage: TypographyUsage): TextStyle {
    return when (usage) {
        TypographyUsage.PAGE_TITLE -> headingLarge
        TypographyUsage.SECTION_TITLE -> headingMedium
        TypographyUsage.CARD_TITLE -> headingSmall
        TypographyUsage.BODY_TEXT -> bodyMedium
        TypographyUsage.SECONDARY_TEXT -> bodySmall
        TypographyUsage.CAPTION_TEXT -> caption
        TypographyUsage.BUTTON_LABEL -> button
        TypographyUsage.NAVIGATION_TITLE -> headingMedium
        TypographyUsage.TAB_LABEL -> label
        TypographyUsage.INPUT_LABEL -> label
        TypographyUsage.ERROR_TEXT -> bodySmall
    }
}

/**
 * Apply accessibility scaling to text style
 */
fun TextStyle.scaledForAccessibility(scaleFactor: Float = 1.0f): TextStyle {
    return this.copy(fontSize = this.fontSize * scaleFactor)
}