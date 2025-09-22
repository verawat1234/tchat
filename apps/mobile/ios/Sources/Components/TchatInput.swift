//
//  TchatInput.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Input field component following Tchat design system
public struct TchatInput: View {

    // MARK: - Input Types
    public enum InputType {
        case text
        case email
        case password
        case number
        case search
        case multiline(minLines: Int = 3, maxLines: Int = 6)
    }

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

    // MARK: - Properties
    @Binding private var text: String
    @FocusState private var isFocused: Bool

    let placeholder: String
    let type: InputType
    let size: Size
    let validationState: ValidationState
    let isDisabled: Bool
    let leadingIcon: String?
    let trailingIcon: String?
    let onTrailingIconTap: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    @State private var isSecureTextVisible = false

    // MARK: - Initializer
    public init(
        text: Binding<String>,
        placeholder: String,
        type: InputType = .text,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        leadingIcon: String? = nil,
        trailingIcon: String? = nil,
        onTrailingIconTap: (() -> Void)? = nil
    ) {
        self._text = text
        self.placeholder = placeholder
        self.type = type
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.leadingIcon = leadingIcon
        self.trailingIcon = trailingIcon
        self.onTrailingIconTap = onTrailingIconTap
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.xs) {
            inputField
                .disabled(isDisabled)

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
                    .padding(.horizontal, Spacing.xs)
            }
        }
    }

    // MARK: - Input Field
    @ViewBuilder
    private var inputField: some View {
        HStack(spacing: Spacing.xs) {
            if let leadingIcon = leadingIcon {
                Image(systemName: leadingIcon)
                    .foregroundColor(iconColor)
                    .frame(width: 16, height: 16)
            }

            inputContent

            trailingIconView
        }
        .padding(.horizontal, horizontalPadding)
        .padding(.vertical, verticalPadding)
        .background(backgroundColor)
        .overlay(
            RoundedRectangle(cornerRadius: Spacing.sm)
                .stroke(borderColor, lineWidth: borderWidth)
        )
        .cornerRadius(Spacing.sm)
        .animation(.easeInOut(duration: 0.2), value: isFocused)
        .animation(.easeInOut(duration: 0.2), value: validationState)
    }

    @ViewBuilder
    private var inputContent: some View {
        switch type {
        case .text:
            TextField(placeholder, text: $text)
                .textFieldStyle(TchatTextFieldStyle())
                .focused($isFocused)

        case .email:
            TextField(placeholder, text: $text)
                .textFieldStyle(TchatTextFieldStyle())
                .keyboardType(.emailAddress)
                .autocapitalization(.none)
                .focused($isFocused)

        case .password:
            if isSecureTextVisible {
                TextField(placeholder, text: $text)
                    .textFieldStyle(TchatTextFieldStyle())
                    .focused($isFocused)
            } else {
                SecureField(placeholder, text: $text)
                    .textFieldStyle(TchatTextFieldStyle())
                    .focused($isFocused)
            }

        case .number:
            TextField(placeholder, text: $text)
                .textFieldStyle(TchatTextFieldStyle())
                .keyboardType(.numberPad)
                .focused($isFocused)

        case .search:
            TextField(placeholder, text: $text)
                .textFieldStyle(TchatTextFieldStyle())
                .autocapitalization(.none)
                .focused($isFocused)

        case .multiline(let minLines, let maxLines):
            TextEditor(text: $text)
                .frame(minHeight: CGFloat(minLines) * 20, maxHeight: CGFloat(maxLines) * 20)
                .focused($isFocused)
        }
    }

    @ViewBuilder
    private var trailingIconView: some View {
        if type == .password {
            Button(action: {
                isSecureTextVisible.toggle()
            }) {
                Image(systemName: isSecureTextVisible ? "eye.slash" : "eye")
                    .foregroundColor(iconColor)
                    .frame(width: 16, height: 16)
            }
        } else if let trailingIcon = trailingIcon {
            Button(action: {
                onTrailingIconTap?()
            }) {
                Image(systemName: trailingIcon)
                    .foregroundColor(iconColor)
                    .frame(width: 16, height: 16)
            }
        }
    }

    // MARK: - Computed Properties

    private var inputFont: Font {
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

// MARK: - Text Field Style
struct TchatTextFieldStyle: TextFieldStyle {
    func _body(configuration: TextField<Self._Label>) -> some View {
        configuration
            .foregroundColor(Colors().textPrimary)
    }
}

// MARK: - Preview
#if DEBUG
struct TchatInput_Previews: PreviewProvider {
    static var previews: some View {
        VStack(spacing: Spacing.md) {
            TchatInput(
                text: .constant(""),
                placeholder: "Enter your email",
                type: .email,
                leadingIcon: "envelope"
            )

            TchatInput(
                text: .constant(""),
                placeholder: "Enter password",
                type: .password,
                leadingIcon: "lock"
            )

            TchatInput(
                text: .constant("Search..."),
                placeholder: "Search",
                type: .search,
                leadingIcon: "magnifyingglass",
                trailingIcon: "xmark.circle.fill"
            )

            TchatInput(
                text: .constant("Valid input"),
                placeholder: "Valid input",
                validationState: .valid,
                leadingIcon: "checkmark.circle"
            )

            TchatInput(
                text: .constant("Invalid input"),
                placeholder: "Invalid input",
                validationState: .invalid(message: "This field is required"),
                leadingIcon: "exclamationmark.triangle"
            )

            TchatInput(
                text: .constant(""),
                placeholder: "Disabled input",
                isDisabled: true
            )

            TchatInput(
                text: .constant(""),
                placeholder: "Type your message here...",
                type: .multiline(minLines: 3, maxLines: 6)
            )

            HStack(spacing: Spacing.sm) {
                TchatInput(
                    text: .constant("Small"),
                    placeholder: "Small",
                    size: .small
                )
                TchatInput(
                    text: .constant("Medium"),
                    placeholder: "Medium",
                    size: .medium
                )
                TchatInput(
                    text: .constant("Large"),
                    placeholder: "Large",
                    size: .large
                )
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif