//
//  TchatCheckbox.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Checkbox component following Tchat design system
public struct TchatCheckbox: View {

    // MARK: - Checkbox Types
    public enum CheckboxState {
        case unchecked
        case checked
        case indeterminate
    }

    public enum Size {
        case small
        case medium
        case large
    }

    public enum ValidationState {
        case none
        case valid
        case invalid(message: String)
    }

    // MARK: - Properties
    @Binding private var state: CheckboxState
    @State private var isPressed = false

    let label: String?
    let description: String?
    let size: Size
    let validationState: ValidationState
    let isDisabled: Bool
    let onChange: ((CheckboxState) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var checkboxSize: CGFloat {
        switch size {
        case .small: return 16
        case .medium: return 20
        case .large: return 24
        }
    }

    private var iconSize: CGFloat {
        switch size {
        case .small: return 10
        case .medium: return 12
        case .large: return 16
        }
    }

    private var textFont: Font {
        switch size {
        case .small: return .system(size: 14)
        case .medium: return .system(size: 16)
        case .large: return .system(size: 18)
        }
    }

    private var descriptionFont: Font {
        switch size {
        case .small: return .system(size: 12)
        case .medium: return .system(size: 14)
        case .large: return .system(size: 16)
        }
    }

    // MARK: - Initializer
    public init(
        state: Binding<CheckboxState>,
        label: String? = nil,
        description: String? = nil,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        onChange: ((CheckboxState) -> Void)? = nil
    ) {
        self._state = state
        self.label = label
        self.description = description
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.onChange = onChange
    }

    // MARK: - Convenience Initializer for Boolean
    public init(
        isChecked: Binding<Bool>,
        label: String? = nil,
        description: String? = nil,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        onChange: ((Bool) -> Void)? = nil
    ) {
        self._state = Binding(
            get: { isChecked.wrappedValue ? .checked : .unchecked },
            set: { newState in
                let boolValue = newState == .checked
                isChecked.wrappedValue = boolValue
                onChange?(boolValue)
            }
        )
        self.label = label
        self.description = description
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.onChange = nil
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.xs) {
            checkboxRow
                .disabled(isDisabled)

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
                    .padding(.leading, checkboxSize + Spacing.sm)
            }
        }
    }

    // MARK: - Checkbox Row
    @ViewBuilder
    private var checkboxRow: some View {
        Button(action: {
            toggleState()
        }) {
            HStack(alignment: .top, spacing: Spacing.sm) {
                checkboxIcon

                if label != nil || description != nil {
                    VStack(alignment: .leading, spacing: Spacing.xs) {
                        if let label = label {
                            Text(label)
                                .font(textFont)
                                .foregroundColor(textColor)
                                .multilineTextAlignment(.leading)
                        }

                        if let description = description {
                            Text(description)
                                .font(descriptionFont)
                                .foregroundColor(descriptionColor)
                                .multilineTextAlignment(.leading)
                        }
                    }
                }

                Spacer()
            }
        }
        .buttonStyle(TchatCheckboxButtonStyle())
    }

    // MARK: - Checkbox Icon
    @ViewBuilder
    private var checkboxIcon: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 4)
                .fill(backgroundColor)
                .frame(width: checkboxSize, height: checkboxSize)
                .overlay(
                    RoundedRectangle(cornerRadius: 4)
                        .stroke(borderColor, lineWidth: borderWidth)
                )

            if state != .unchecked {
                Image(systemName: iconName)
                    .font(.system(size: iconSize, weight: .bold))
                    .foregroundColor(iconColor)
                    .animation(.easeInOut(duration: 0.15), value: state)
            }
        }
        .scaleEffect(isPressed ? 0.95 : 1.0)
        .animation(.easeInOut(duration: 0.1), value: isPressed)
    }

    // MARK: - Computed Style Properties
    private var backgroundColor: Color {
        if isDisabled {
            return colors.surface.opacity(0.5)
        }

        switch state {
        case .unchecked:
            return colors.background
        case .checked, .indeterminate:
            return colors.primary
        }
    }

    private var borderColor: Color {
        if isDisabled {
            return colors.border.opacity(0.5)
        }

        switch validationState {
        case .none:
            switch state {
            case .unchecked:
                return colors.border
            case .checked, .indeterminate:
                return colors.primary
            }
        case .valid:
            return colors.success
        case .invalid:
            return colors.borderError
        }
    }

    private var borderWidth: CGFloat {
        switch validationState {
        case .none:
            return state == .unchecked ? 2 : 0
        case .valid, .invalid:
            return 2
        }
    }

    private var iconColor: Color {
        if isDisabled {
            return colors.textDisabled
        }
        return colors.textOnPrimary
    }

    private var iconName: String {
        switch state {
        case .unchecked:
            return ""
        case .checked:
            return "checkmark"
        case .indeterminate:
            return "minus"
        }
    }

    private var textColor: Color {
        if isDisabled {
            return colors.textDisabled
        }
        return colors.textPrimary
    }

    private var descriptionColor: Color {
        if isDisabled {
            return colors.textDisabled
        }
        return colors.textSecondary
    }

    // MARK: - State Management
    private func toggleState() {
        let newState: CheckboxState
        switch state {
        case .unchecked:
            newState = .checked
        case .checked:
            newState = .unchecked
        case .indeterminate:
            newState = .checked
        }

        state = newState
        onChange?(newState)

        // Add haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Button Style
struct TchatCheckboxButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .contentShape(Rectangle())
    }
}

// MARK: - Checkbox Group
public struct TchatCheckboxGroup<T: Hashable>: View {
    @Binding private var selection: Set<T>
    let options: [T]
    let optionLabel: (T) -> String
    let optionDescription: ((T) -> String)?
    let title: String?
    let size: TchatCheckbox.Size
    let validationState: TchatCheckbox.ValidationState
    let isDisabled: Bool
    let maxSelections: Int?
    let onChange: ((Set<T>) -> Void)?

    private let colors = Colors()

    public init(
        selection: Binding<Set<T>>,
        options: [T],
        optionLabel: @escaping (T) -> String,
        optionDescription: ((T) -> String)? = nil,
        title: String? = nil,
        size: TchatCheckbox.Size = .medium,
        validationState: TchatCheckbox.ValidationState = .none,
        isDisabled: Bool = false,
        maxSelections: Int? = nil,
        onChange: ((Set<T>) -> Void)? = nil
    ) {
        self._selection = selection
        self.options = options
        self.optionLabel = optionLabel
        self.optionDescription = optionDescription
        self.title = title
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.maxSelections = maxSelections
        self.onChange = onChange
    }

    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.sm) {
            if let title = title {
                Text(title)
                    .font(.headline)
                    .foregroundColor(isDisabled ? colors.textDisabled : colors.textPrimary)
            }

            VStack(alignment: .leading, spacing: Spacing.sm) {
                ForEach(Array(options.enumerated()), id: \.element) { index, option in
                    TchatCheckbox(
                        isChecked: Binding(
                            get: { selection.contains(option) },
                            set: { isChecked in
                                if isChecked {
                                    if let maxSelections = maxSelections,
                                       selection.count >= maxSelections {
                                        return
                                    }
                                    selection.insert(option)
                                } else {
                                    selection.remove(option)
                                }
                                onChange?(selection)
                            }
                        ),
                        label: optionLabel(option),
                        description: optionDescription?(option),
                        size: size,
                        validationState: index == 0 ? validationState : .none,
                        isDisabled: isDisabled
                    )
                }
            }

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
            }
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatCheckbox_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Basic checkbox
                TchatCheckbox(
                    isChecked: .constant(true),
                    label: "Accept terms and conditions"
                )

                // Checkbox with description
                TchatCheckbox(
                    isChecked: .constant(false),
                    label: "Enable notifications",
                    description: "Receive updates about new messages and mentions"
                )

                // Indeterminate state
                TchatCheckbox(
                    state: .constant(.indeterminate),
                    label: "Select all items",
                    description: "Some items are selected"
                )

                // Validation states
                TchatCheckbox(
                    isChecked: .constant(true),
                    label: "Valid selection",
                    validationState: .valid
                )

                TchatCheckbox(
                    isChecked: .constant(false),
                    label: "Required field",
                    validationState: .invalid("This field is required")
                )

                // Disabled state
                TchatCheckbox(
                    isChecked: .constant(true),
                    label: "Disabled checkbox",
                    description: "This option cannot be changed",
                    isDisabled: true
                )

                // Different sizes
                VStack(alignment: .leading, spacing: Spacing.sm) {
                    TchatCheckbox(
                        isChecked: .constant(true),
                        label: "Small checkbox",
                        size: .small
                    )

                    TchatCheckbox(
                        isChecked: .constant(true),
                        label: "Medium checkbox",
                        size: .medium
                    )

                    TchatCheckbox(
                        isChecked: .constant(true),
                        label: "Large checkbox",
                        size: .large
                    )
                }

                Divider()

                // Checkbox group
                TchatCheckboxGroup(
                    selection: .constant(Set(["Swift", "Kotlin"])),
                    options: ["Swift", "Kotlin", "JavaScript", "Python", "Go"],
                    optionLabel: { $0 },
                    optionDescription: { lang in
                        switch lang {
                        case "Swift": return "iOS and macOS development"
                        case "Kotlin": return "Android and multiplatform development"
                        case "JavaScript": return "Web and Node.js development"
                        case "Python": return "Data science and backend development"
                        case "Go": return "System programming and microservices"
                        default: return nil
                        }
                    },
                    title: "Select programming languages",
                    maxSelections: 3
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif