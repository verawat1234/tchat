//
//  TchatSwitch.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Toggle switch component following Tchat design system
public struct TchatSwitch: View {

    // MARK: - Switch Types
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
    @Binding private var isOn: Bool
    @State private var dragOffset: CGFloat = 0
    @State private var isDragging = false

    let label: String?
    let description: String?
    let size: Size
    let validationState: ValidationState
    let isDisabled: Bool
    let showLabels: Bool
    let onText: String
    let offText: String
    let onChange: ((Bool) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var switchWidth: CGFloat {
        switch size {
        case .small: return 40
        case .medium: return 50
        case .large: return 60
        }
    }

    private var switchHeight: CGFloat {
        switch size {
        case .small: return 24
        case .medium: return 30
        case .large: return 36
        }
    }

    private var thumbSize: CGFloat {
        switch size {
        case .small: return 18
        case .medium: return 24
        case .large: return 30
        }
    }

    private var thumbOffset: CGFloat {
        let maxOffset = switchWidth - thumbSize - 6 // 3pt padding on each side
        let currentOffset = isOn ? maxOffset : 3
        return currentOffset + dragOffset
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

    private var labelFont: Font {
        switch size {
        case .small: return .system(size: 10, weight: .medium)
        case .medium: return .system(size: 12, weight: .medium)
        case .large: return .system(size: 14, weight: .medium)
        }
    }

    // MARK: - Initializer
    public init(
        isOn: Binding<Bool>,
        label: String? = nil,
        description: String? = nil,
        size: Size = .medium,
        validationState: ValidationState = .none,
        isDisabled: Bool = false,
        showLabels: Bool = false,
        onText: String = "ON",
        offText: String = "OFF",
        onChange: ((Bool) -> Void)? = nil
    ) {
        self._isOn = isOn
        self.label = label
        self.description = description
        self.size = size
        self.validationState = validationState
        self.isDisabled = isDisabled
        self.showLabels = showLabels
        self.onText = onText
        self.offText = offText
        self.onChange = onChange
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.xs) {
            switchRow
                .disabled(isDisabled)

            if case .invalid(let message) = validationState {
                Text(message)
                    .font(.caption)
                    .foregroundColor(colors.error)
                    .padding(.leading, switchWidth + Spacing.sm)
            }
        }
    }

    // MARK: - Switch Row
    @ViewBuilder
    private var switchRow: some View {
        HStack(alignment: .center, spacing: Spacing.sm) {
            switchControl

            if label != nil || description != nil {
                VStack(alignment: .leading, spacing: Spacing.xs) {
                    if let label = label {
                        Text(label)
                            .font(textFont)
                            .foregroundColor(textColor)
                    }

                    if let description = description {
                        Text(description)
                            .font(descriptionFont)
                            .foregroundColor(descriptionColor)
                    }
                }
            }

            Spacer()
        }
        .contentShape(Rectangle())
        .onTapGesture {
            if !isDisabled {
                toggleSwitch()
            }
        }
    }

    // MARK: - Switch Control
    @ViewBuilder
    private var switchControl: some View {
        ZStack {
            // Track
            RoundedRectangle(cornerRadius: switchHeight / 2)
                .fill(trackColor)
                .frame(width: switchWidth, height: switchHeight)
                .overlay(
                    RoundedRectangle(cornerRadius: switchHeight / 2)
                        .stroke(borderColor, lineWidth: borderWidth)
                )

            // Labels inside track (if enabled)
            if showLabels {
                HStack {
                    Text(offText)
                        .font(labelFont)
                        .foregroundColor(isOn ? colors.textTertiary : colors.textOnPrimary)
                        .opacity(isOn ? 0.5 : 1.0)

                    Spacer()

                    Text(onText)
                        .font(labelFont)
                        .foregroundColor(isOn ? colors.textOnPrimary : colors.textTertiary)
                        .opacity(isOn ? 1.0 : 0.5)
                }
                .padding(.horizontal, 6)
                .frame(width: switchWidth, height: switchHeight)
            }

            // Thumb
            Circle()
                .fill(thumbColor)
                .frame(width: thumbSize, height: thumbSize)
                .shadow(color: colors.shadowMedium, radius: 2, x: 0, y: 1)
                .offset(x: thumbOffset - switchWidth / 2 + thumbSize / 2)
                .animation(.easeInOut(duration: 0.2), value: isOn)
                .animation(.easeInOut(duration: 0.1), value: isDragging)
        }
        .gesture(
            DragGesture()
                .onChanged { value in
                    if isDisabled { return }

                    isDragging = true
                    let maxDrag = switchWidth - thumbSize - 6
                    dragOffset = max(-3, min(maxDrag - 3, value.translation.x))
                }
                .onEnded { value in
                    if isDisabled { return }

                    isDragging = false
                    let threshold = switchWidth / 2
                    let shouldToggle = abs(value.translation.x) > 10 &&
                                     ((value.translation.x > 0 && !isOn) ||
                                      (value.translation.x < 0 && isOn))

                    dragOffset = 0

                    if shouldToggle {
                        toggleSwitch()
                    }
                }
        )
    }

    // MARK: - Computed Style Properties
    private var trackColor: Color {
        if isDisabled {
            return colors.surface.opacity(0.5)
        }

        switch validationState {
        case .none:
            return isOn ? colors.primary : colors.border.opacity(0.3)
        case .valid:
            return isOn ? colors.success : colors.success.opacity(0.3)
        case .invalid:
            return isOn ? colors.error : colors.error.opacity(0.3)
        }
    }

    private var borderColor: Color {
        if isDisabled {
            return colors.border.opacity(0.5)
        }

        switch validationState {
        case .none:
            return isOn ? colors.primary : colors.border
        case .valid:
            return colors.success
        case .invalid:
            return colors.borderError
        }
    }

    private var borderWidth: CGFloat {
        switch validationState {
        case .none:
            return isOn ? 0 : 1
        case .valid, .invalid:
            return 2
        }
    }

    private var thumbColor: Color {
        if isDisabled {
            return colors.surface
        }
        return colors.background
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

    // MARK: - Actions
    private func toggleSwitch() {
        withAnimation(.easeInOut(duration: 0.2)) {
            isOn.toggle()
        }

        onChange?(isOn)

        // Add haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Switch Group
public struct TchatSwitchGroup: View {
    let title: String?
    let switches: [SwitchItem]
    let isDisabled: Bool

    private let colors = Colors()

    public struct SwitchItem {
        let id: String
        let label: String
        let description: String?
        let isOn: Binding<Bool>
        let validationState: TchatSwitch.ValidationState
        let onChange: ((Bool) -> Void)?

        public init(
            id: String,
            label: String,
            description: String? = nil,
            isOn: Binding<Bool>,
            validationState: TchatSwitch.ValidationState = .none,
            onChange: ((Bool) -> Void)? = nil
        ) {
            self.id = id
            self.label = label
            self.description = description
            self.isOn = isOn
            self.validationState = validationState
            self.onChange = onChange
        }
    }

    public init(
        title: String? = nil,
        switches: [SwitchItem],
        isDisabled: Bool = false
    ) {
        self.title = title
        self.switches = switches
        self.isDisabled = isDisabled
    }

    public var body: some View {
        VStack(alignment: .leading, spacing: Spacing.md) {
            if let title = title {
                Text(title)
                    .font(.headline)
                    .foregroundColor(isDisabled ? colors.textDisabled : colors.textPrimary)
            }

            VStack(alignment: .leading, spacing: Spacing.sm) {
                ForEach(switches, id: \.id) { switchItem in
                    TchatSwitch(
                        isOn: switchItem.isOn,
                        label: switchItem.label,
                        description: switchItem.description,
                        validationState: switchItem.validationState,
                        isDisabled: isDisabled,
                        onChange: switchItem.onChange
                    )
                }
            }
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatSwitch_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Basic switches
                TchatSwitch(
                    isOn: .constant(true),
                    label: "Enable notifications"
                )

                TchatSwitch(
                    isOn: .constant(false),
                    label: "Dark mode",
                    description: "Use dark appearance throughout the app"
                )

                // Switch with labels
                TchatSwitch(
                    isOn: .constant(true),
                    label: "Auto-save",
                    description: "Automatically save changes",
                    showLabels: true
                )

                // Validation states
                TchatSwitch(
                    isOn: .constant(true),
                    label: "Valid setting",
                    validationState: .valid
                )

                TchatSwitch(
                    isOn: .constant(false),
                    label: "Required setting",
                    validationState: .invalid("This setting is required")
                )

                // Disabled state
                TchatSwitch(
                    isOn: .constant(true),
                    label: "Disabled switch",
                    description: "This setting cannot be changed",
                    isDisabled: true
                )

                // Different sizes
                VStack(alignment: .leading, spacing: Spacing.sm) {
                    TchatSwitch(
                        isOn: .constant(true),
                        label: "Small switch",
                        size: .small
                    )

                    TchatSwitch(
                        isOn: .constant(true),
                        label: "Medium switch",
                        size: .medium
                    )

                    TchatSwitch(
                        isOn: .constant(true),
                        label: "Large switch",
                        size: .large
                    )
                }

                Divider()

                // Switch group
                TchatSwitchGroup(
                    title: "Notification Settings",
                    switches: [
                        TchatSwitchGroup.SwitchItem(
                            id: "email",
                            label: "Email notifications",
                            description: "Receive notifications via email",
                            isOn: .constant(true)
                        ),
                        TchatSwitchGroup.SwitchItem(
                            id: "push",
                            label: "Push notifications",
                            description: "Receive push notifications on your device",
                            isOn: .constant(false)
                        ),
                        TchatSwitchGroup.SwitchItem(
                            id: "sms",
                            label: "SMS notifications",
                            description: "Receive notifications via SMS",
                            isOn: .constant(true),
                            validationState: .valid
                        )
                    ]
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif