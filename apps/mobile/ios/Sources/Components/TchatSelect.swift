//
//  TchatSelect.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Advanced select/dropdown component following Tchat design system
public struct TchatSelect<T: Hashable & CustomStringConvertible>: View {

    // MARK: - Select Types
    public enum SelectMode {
        case single
        case multiple
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
    @Binding private var selection: Set<T>
    @State private var isExpanded = false
    @State private var searchText = ""
    @FocusState private var isSearchFocused: Bool

    let options: [T]
    let placeholder: String
    let mode: SelectMode
    let size: Size
    let validationState: ValidationState
    let isDisabled: Bool
    let isSearchable: Bool
    let maxSelections: Int?
    let leadingIcon: String?
    let onSelectionChange: ((Set<T>) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var filteredOptions: [T] {
        if searchText.isEmpty {
            return options
        }
        return options.filter { option in
            option.description.localizedCaseInsensitiveContains(searchText)
        }
    }

    private var displayText: String {
        switch mode {
        case .single:
            return selection.first?.description ?? placeholder
        case .multiple:
            if selection.isEmpty {
                return placeholder
            } else if selection.count == 1 {
                return selection.first?.description ?? ""
            } else {
                return "\(selection.count) selected"
            }
        }
    }

    private var isPlaceholderShown: Bool {
        selection.isEmpty
    }

    // MARK: - Initializer
    public init(
        selection: Binding<Set<T>>,
        options: [T],
        placeholder: String = "Select option",
        mode: SelectMode = .single,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        isSearchable: Bool = false,
        maxSelections: Int? = nil,
        leadingIcon: String? = nil,
        onSelectionChange: ((Set<T>) -> Void)? = nil
    ) {
        self._selection = selection
        self.options = options
        self.placeholder = placeholder
        self.mode = mode
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.isSearchable = isSearchable
        self.maxSelections = maxSelections
        self.leadingIcon = leadingIcon
        self.onSelectionChange = onSelectionChange
    }

    // MARK: - Single Selection Convenience Initializer
    public init(
        selection: Binding<T?>,
        options: [T],
        placeholder: String = "Select option",
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        isSearchable: Bool = false,
        leadingIcon: String? = nil,
        onSelectionChange: ((T?) -> Void)? = nil
    ) where T: Hashable {
        self._selection = Binding(
            get: {
                if let value = selection.wrappedValue {
                    return Set([value])
                }
                return Set<T>()
            },
            set: { newSet in
                selection.wrappedValue = newSet.first
                onSelectionChange?(newSet.first)
            }
        )
        self.options = options
        self.placeholder = placeholder
        self.mode = .single
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.isSearchable = isSearchable
        self.maxSelections = nil
        self.leadingIcon = leadingIcon
        self.onSelectionChange = { _ in }
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.xs) {
            selectButton
                .disabled(isDisabled)

            if isExpanded {
                optionsDropdown
                    .zIndex(1000)
            }

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
                    .padding(.horizontal, Spacing.xs)
            }
        }
        .onTapGesture {
            if !isDisabled {
                withAnimation(.easeInOut(duration: 0.2)) {
                    isExpanded.toggle()
                }
            }
        }
    }

    // MARK: - Select Button
    @ViewBuilder
    private var selectButton: some View {
        HStack(spacing: Spacing.xs) {
            if let leadingIcon = leadingIcon {
                Image(systemName: leadingIcon)
                    .foregroundColor(iconColor)
                    .frame(width: 16, height: 16)
            }

            Text(displayText)
                .font(selectFont)
                .foregroundColor(isPlaceholderShown ? colors.textTertiary : colors.textPrimary)
                .lineLimit(1)

            Spacer()

            // Selection indicators for multiple mode
            if mode == .multiple && !selection.isEmpty {
                Text("\(selection.count)")
                    .font(.caption)
                    .foregroundColor(colors.textOnPrimary)
                    .padding(.horizontal, Spacing.xs)
                    .padding(.vertical, 2)
                    .background(colors.primary)
                    .cornerRadius(10)
            }

            Image(systemName: isExpanded ? "chevron.up" : "chevron.down")
                .foregroundColor(iconColor)
                .font(.system(size: 12, weight: .medium))
                .animation(.easeInOut(duration: 0.2), value: isExpanded)
        }
        .padding(.horizontal, horizontalPadding)
        .padding(.vertical, verticalPadding)
        .background(backgroundColor)
        .overlay(
            RoundedRectangle(cornerRadius: Spacing.sm)
                .stroke(borderColor, lineWidth: borderWidth)
        )
        .cornerRadius(Spacing.sm)
    }

    // MARK: - Options Dropdown
    @ViewBuilder
    private var optionsDropdown: some View {
        VStack(spacing: 0) {
            if isSearchable {
                searchField
                    .padding(Spacing.sm)
            }

            ScrollView {
                LazyVStack(spacing: 0) {
                    ForEach(Array(filteredOptions.enumerated()), id: \.element) { index, option in
                        optionRow(option, isLast: index == filteredOptions.count - 1)
                    }
                }
            }
            .frame(maxHeight: 200)
        }
        .background(colors.cardBackground)
        .overlay(
            RoundedRectangle(cornerRadius: Spacing.sm)
                .stroke(colors.border, lineWidth: 1)
        )
        .cornerRadius(Spacing.sm)
        .shadow(color: colors.shadowMedium, radius: 8, x: 0, y: 4)
    }

    // MARK: - Search Field
    @ViewBuilder
    private var searchField: some View {
        HStack {
            Image(systemName: "magnifyingglass")
                .foregroundColor(colors.textTertiary)
                .font(.system(size: 14))

            TextField("Search options", text: $searchText)
                .font(.system(size: 14))
                .focused($isSearchFocused)
        }
        .padding(.horizontal, Spacing.sm)
        .padding(.vertical, Spacing.xs)
        .background(colors.surface)
        .cornerRadius(Spacing.xs)
        .onAppear {
            if isSearchable {
                isSearchFocused = true
            }
        }
    }

    // MARK: - Option Row
    @ViewBuilder
    private func optionRow(_ option: T, isLast: Bool) -> some View {
        Button(action: {
            handleOptionSelection(option)
        }) {
            HStack {
                Text(option.description)
                    .font(.system(size: 14))
                    .foregroundColor(colors.textPrimary)

                Spacer()

                if selection.contains(option) {
                    Image(systemName: mode == .single ? "checkmark" : "checkmark.square.fill")
                        .foregroundColor(colors.primary)
                        .font(.system(size: 14, weight: .medium))
                }
            }
            .padding(.horizontal, Spacing.sm)
            .padding(.vertical, Spacing.sm)
            .background(
                selection.contains(option) ? colors.primary.opacity(0.1) : Color.clear
            )
        }
        .buttonStyle(PlainButtonStyle())

        if !isLast {
            Divider()
                .padding(.horizontal, Spacing.sm)
        }
    }

    // MARK: - Computed Style Properties
    private var selectFont: Font {
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
            return isExpanded ? colors.borderFocus : colors.border
        case .valid:
            return colors.success
        case .invalid:
            return colors.borderError
        }
    }

    private var borderWidth: CGFloat {
        switch validationState {
        case .none:
            return isExpanded ? 2 : 1
        case .valid, .invalid:
            return 2
        }
    }

    private var iconColor: Color {
        if isDisabled {
            return colors.textDisabled
        }
        return isExpanded ? colors.primary : colors.textSecondary
    }

    // MARK: - Selection Handling
    private func handleOptionSelection(_ option: T) {
        switch mode {
        case .single:
            selection = Set([option])
            withAnimation(.easeInOut(duration: 0.2)) {
                isExpanded = false
            }

        case .multiple:
            if selection.contains(option) {
                selection.remove(option)
            } else {
                if let maxSelections = maxSelections, selection.count >= maxSelections {
                    // Could show an alert or handle max selection reached
                    return
                }
                selection.insert(option)
            }
        }

        onSelectionChange?(selection)

        // Add haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Preview
#if DEBUG
struct TchatSelect_Previews: PreviewProvider {
    static var previews: some View {
        VStack(spacing: Spacing.lg) {
            // Single selection
            TchatSelect(
                selection: .constant(Set(["Option 1"])),
                options: ["Option 1", "Option 2", "Option 3", "Option 4"],
                placeholder: "Choose an option",
                mode: .single,
                leadingIcon: "star"
            )

            // Multiple selection with search
            TchatSelect(
                selection: .constant(Set(["Swift", "Kotlin"])),
                options: ["Swift", "Kotlin", "JavaScript", "Python", "Go", "Rust"],
                placeholder: "Select languages",
                mode: .multiple,
                isSearchable: true,
                leadingIcon: "code"
            )

            // Validation states
            TchatSelect(
                selection: .constant(Set<String>()),
                options: ["Valid", "Invalid", "Neutral"],
                placeholder: "Required field",
                validationState: .invalid("This field is required"),
                leadingIcon: "exclamationmark.triangle"
            )

            // Disabled state
            TchatSelect(
                selection: .constant(Set(["Disabled"])),
                options: ["Disabled", "Option"],
                placeholder: "Disabled select",
                isDisabled: true
            )

            // Different sizes
            VStack(spacing: Spacing.sm) {
                TchatSelect(
                    selection: .constant(Set<String>()),
                    options: ["Small"],
                    placeholder: "Small size",
                    size: .small
                )

                TchatSelect(
                    selection: .constant(Set<String>()),
                    options: ["Medium"],
                    placeholder: "Medium size",
                    size: .medium
                )

                TchatSelect(
                    selection: .constant(Set<String>()),
                    options: ["Large"],
                    placeholder: "Large size",
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