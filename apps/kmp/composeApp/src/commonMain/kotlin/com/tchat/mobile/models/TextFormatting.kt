package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Text formatting model for rich text content
 * Supports various text styling options
 */
@Serializable
data class TextFormatting(
    val bold: Boolean = false,
    val italic: Boolean = false,
    val underline: Boolean = false,
    val strikethrough: Boolean = false,
    val fontSize: Float? = null,
    val color: String? = null,
    val backgroundColor: String? = null,
    val fontFamily: String? = null,
    val alignment: TextAlignment = TextAlignment.START,
    val lineHeight: Float? = null
)

@Serializable
enum class TextAlignment {
    START, CENTER, END, JUSTIFY
}

/**
 * Rich text span with formatting
 */
@Serializable
data class FormattedTextSpan(
    val text: String,
    val formatting: TextFormatting? = null,
    val startIndex: Int,
    val endIndex: Int
)

/**
 * Complete formatted text document
 */
@Serializable
data class FormattedText(
    val plainText: String,
    val spans: List<FormattedTextSpan> = emptyList(),
    val globalFormatting: TextFormatting? = null
)