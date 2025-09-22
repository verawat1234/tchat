//
//  TchatTextarea.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Multi-line text input component following Tchat design system
public struct TchatTextarea: View {

    // MARK: - Textarea Types
    public enum ValidationState {
        case none
        case valid
        case invalid(message: String)
    }

    public enum Size {
        case small
        case medium
        case large
    }

    public enum ResizeBehavior {
        case fixed(height: CGFloat)
        case autoResize(minLines: Int, maxLines: Int)
        case expandable
    }

    // MARK: - Properties
    @Binding private var text: String
    @FocusState private var isFocused: Bool
    @State private var textHeight: CGFloat = 0

    let placeholder: String
    let size: Size
    let validationState: ValidationState
    let isDisabled: Bool
    let resizeBehavior: ResizeBehavior
    let characterLimit: Int?
    let showCharacterCount: Bool
    let leadingIcon: String?
    let onTextChange: ((String) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var remainingCharacters: Int? {
        guard let limit = characterLimit else { return nil }
        return max(0, limit - text.count)
    }

    private var isCharacterLimitExceeded: Bool {
        guard let limit = characterLimit else { return false }
        return text.count > limit
    }

    private var dynamicHeight: CGFloat {
        switch resizeBehavior {
        case .fixed(let height):
            return height
        case .autoResize(let minLines, let maxLines):
            let lineHeight = lineHeightForSize
            let minHeight = CGFloat(minLines) * lineHeight + verticalPadding * 2
            let maxHeight = CGFloat(maxLines) * lineHeight + verticalPadding * 2
            return max(minHeight, min(maxHeight, textHeight + verticalPadding * 2))
        case .expandable:
            let lineHeight = lineHeightForSize
            let minHeight = lineHeight + verticalPadding * 2
            return max(minHeight, textHeight + verticalPadding * 2)
        }
    }

    private var lineHeightForSize: CGFloat {
        switch size {
        case .small: return 18
        case .medium: return 22
        case .large: return 26
        }
    }

    // MARK: - Initializer
    public init(
        text: Binding<String>,
        placeholder: String = "Enter text...",
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        resizeBehavior: ResizeBehavior = .autoResize(minLines: 3, maxLines: 8),
        characterLimit: Int? = nil,
        showCharacterCount: Bool = false,
        leadingIcon: String? = nil,
        onTextChange: ((String) -> Void)? = nil
    ) {
        self._text = text
        self.placeholder = placeholder
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.resizeBehavior = resizeBehavior
        self.characterLimit = characterLimit
        self.showCharacterCount = showCharacterCount
        self.leadingIcon = leadingIcon
        self.onTextChange = onTextChange
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.xs) {
            textareaField
                .disabled(isDisabled)

            bottomSection

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
                    .padding(.horizontal, Spacing.xs)
            }
        }
    }

    // MARK: - Textarea Field
    @ViewBuilder
    private var textareaField: some View {
        HStack(alignment: .top, spacing: Spacing.xs) {
            if let leadingIcon = leadingIcon {
                Image(systemName: leadingIcon)
                    .foregroundColor(iconColor)
                    .frame(width: 16, height: 16)
                    .padding(.top, verticalPadding)
            }

            ZStack(alignment: .topLeading) {
                // Background for measuring text height
                Text(text.isEmpty ? placeholder : text)
                    .font(textareaFont)
                    .foregroundColor(.clear)
                    .padding(.horizontal, 0)
                    .padding(.vertical, 0)
                    .background(
                        GeometryReader { geometry in
                            Color.clear.onAppear {
                                textHeight = geometry.size.height
                            }
                            .onChange(of: text) { _ in
                                DispatchQueue.main.async {
                                    textHeight = geometry.size.height
                                }
                            }
                        }
                    )

                // Placeholder
                if text.isEmpty {
                    Text(placeholder)
                        .font(textareaFont)
                        .foregroundColor(colors.textTertiary)
                        .allowsHitTesting(false)
                }

                // Actual text editor
                TextEditor(text: $text)
                    .font(textareaFont)
                    .foregroundColor(colors.textPrimary)
                    .background(Color.clear)
                    .focused($isFocused)
                    .onChange(of: text) { newValue in
                        // Handle character limit
                        if let limit = characterLimit, newValue.count > limit {
                            text = String(newValue.prefix(limit))
                        } else {
                            onTextChange?(newValue)
                        }
                    }
                    .scrollContentBackground(.hidden)
            }
        }
        .frame(height: dynamicHeight)
        .padding(.horizontal, horizontalPadding)
        .padding(.vertical, verticalPadding)
        .background(backgroundColor)
        .overlay(
            RoundedRectangle(cornerRadius: Spacing.sm)
                .stroke(borderColor, lineWidth: borderWidth)
        )
        .cornerRadius(Spacing.sm)
        .animation(.easeInOut(duration: 0.2), value: dynamicHeight)
    }

    // MARK: - Bottom Section
    @ViewBuilder
    private var bottomSection: some View {
        if showCharacterCount || characterLimit != nil {
            HStack {
                Spacer()

                if let remaining = remainingCharacters {
                    Text("\(remaining) remaining")
                        .font(.caption)
                        .foregroundColor(
                            isCharacterLimitExceeded ? colors.error : colors.textTertiary
                        )
                } else if showCharacterCount {
                    Text("\(text.count) characters")
                        .font(.caption)
                        .foregroundColor(colors.textTertiary)
                }
            }
            .padding(.horizontal, Spacing.xs)
        }
    }

    // MARK: - Computed Style Properties
    private var textareaFont: Font {
        switch size {
        case .small:
            return .system(size: 14)
        case .medium:
            return .system(size: 16)
        case .large:
            return .system(size: 18)
        }
    }

    private var horizontalPadding: CGFloat {
        switch size {
        case .small: return Spacing.sm
        case .medium: return Spacing.md
        case .large: return Spacing.lg
        }
    }

    private var verticalPadding: CGFloat {
        switch size {
        case .small: return Spacing.xs
        case .medium: return Spacing.sm
        case .large: return Spacing.md
        }
    }

    private var backgroundColor: Color {
        if isDisabled {
            return colors.surface.opacity(0.5)
        }
        return colors.background
    }

    private var borderColor: Color {
        if isDisabled {
            return colors.border.opacity(0.5)
        }

        if isCharacterLimitExceeded {
            return colors.borderError
        }

        switch validationState {
        case .none:
            return isFocused ? colors.borderFocus : colors.border
        case .valid:
            return colors.success
        case .invalid:
            return colors.borderError
        }
    }

    private var borderWidth: CGFloat {
        if isCharacterLimitExceeded {
            return 2
        }

        switch validationState {
        case .none:
            return isFocused ? 2 : 1
        case .valid, .invalid:
            return 2
        }
    }

    private var iconColor: Color {
        if isDisabled {
            return colors.textDisabled
        }
        return isFocused ? colors.primary : colors.textSecondary
    }
}

// MARK: - Convenience Extensions
extension TchatTextarea {
    /// Create a textarea with fixed height
    public static func fixed(
        text: Binding<String>,
        placeholder: String = "Enter text...",
        height: CGFloat,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        characterLimit: Int? = nil,
        showCharacterCount: Bool = false
    ) -> TchatTextarea {
        TchatTextarea(
            text: text,
            placeholder: placeholder,
            size: size,
            validationState: validationState,
            isDisabled: isDisabled,
            resizeBehavior: .fixed(height: height),
            characterLimit: characterLimit,
            showCharacterCount: showCharacterCount
        )
    }

    /// Create an expandable textarea
    public static func expandable(
        text: Binding<String>,
        placeholder: String = "Enter text...",
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        characterLimit: Int? = nil,
        showCharacterCount: Bool = false
    ) -> TchatTextarea {
        TchatTextarea(
            text: text,
            placeholder: placeholder,
            size: size,
            validationState: validationState,
            isDisabled: isDisabled,
            resizeBehavior: .expandable,
            characterLimit: characterLimit,
            showCharacterCount: showCharacterCount
        )
    }
}

// MARK: - Preview
#if DEBUG
struct TchatTextarea_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Auto-resize textarea
                TchatTextarea(
                    text: .constant("This is a sample text that demonstrates the auto-resize functionality of the textarea component."),
                    placeholder: "Enter your message...",
                    resizeBehavior: .autoResize(minLines: 3, maxLines: 6),
                    leadingIcon: "text.alignleft"
                )

                // Fixed height textarea
                TchatTextarea.fixed(
                    text: .constant("Fixed height textarea"),
                    placeholder: "Fixed height (100pt)",
                    height: 100,
                    characterLimit: 500,
                    showCharacterCount: true
                )

                // Expandable textarea
                TchatTextarea.expandable(
                    text: .constant("Expandable textarea that grows with content"),
                    placeholder: "Type to see expansion...",
                    characterLimit: 280,
                    showCharacterCount: true
                )

                // Validation states
                TchatTextarea(
                    text: .constant("Valid content"),
                    placeholder: "Valid textarea",
                    validationState: .valid,
                    resizeBehavior: .autoResize(minLines: 2, maxLines: 4),
                    leadingIcon: "checkmark.circle"
                )

                TchatTextarea(
                    text: .constant(""),
                    placeholder: "Required field",
                    validationState: .invalid("This field is required"),
                    resizeBehavior: .autoResize(minLines: 2, maxLines: 4),
                    leadingIcon: "exclamationmark.triangle"
                )

                // Disabled state
                TchatTextarea(
                    text: .constant("Disabled textarea content"),
                    placeholder: "Disabled",
                    isDisabled: true,
                    resizeBehavior: .autoResize(minLines: 3, maxLines: 5)
                )

                // Different sizes
                VStack(spacing: Spacing.sm) {
                    TchatTextarea(
                        text: .constant("Small size"),
                        placeholder: "Small textarea",
                        size: .small,
                        resizeBehavior: .autoResize(minLines: 2, maxLines: 3)
                    )

                    TchatTextarea(
                        text: .constant("Large size"),
                        placeholder: "Large textarea",
                        size: .large,
                        resizeBehavior: .autoResize(minLines: 2, maxLines: 3)
                    )
                }

                // Character limit example
                TchatTextarea(
                    text: .constant("This demonstrates character limit functionality"),
                    placeholder: "Tweet-like input (280 chars)",
                    characterLimit: 280,
                    showCharacterCount: true,
                    resizeBehavior: .autoResize(minLines: 3, maxLines: 6),
                    leadingIcon: "at"
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif