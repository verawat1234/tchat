//
//  TchatAlert.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Alert component following Tchat design system
public struct TchatAlert: View {

    // MARK: - Alert Types
    public enum AlertType {
        case info
        case success
        case warning
        case error
    }

    public enum AlertVariant {
        case filled
        case outlined
        case minimal
    }

    public enum AlertSize {
        case small
        case medium
        case large
    }

    // MARK: - Alert Action
    public struct AlertAction {
        let title: String
        let style: AlertActionStyle
        let action: () -> Void

        public enum AlertActionStyle {
            case primary
            case secondary
            case destructive
        }

        public init(
            title: String,
            style: AlertActionStyle = .primary,
            action: @escaping () -> Void
        ) {
            self.title = title
            self.style = style
            self.action = action
        }
    }

    // MARK: - Properties
    @Binding private var isPresented: Bool
    @State private var animationOffset: CGFloat = 30
    @State private var animationOpacity: Double = 0

    let type: AlertType
    let variant: AlertVariant
    let size: AlertSize
    let title: String?
    let message: String
    let isDismissible: Bool
    let showIcon: Bool
    let actions: [AlertAction]
    let onDismiss: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var alertColors: (background: Color, border: Color, text: Color, icon: Color) {
        switch (type, variant) {
        case (.info, .filled):
            return (colors.primary, colors.primary, colors.textOnPrimary, colors.textOnPrimary)
        case (.info, .outlined):
            return (colors.background, colors.primary, colors.primary, colors.primary)
        case (.info, .minimal):
            return (colors.primary.opacity(0.1), Color.clear, colors.primary, colors.primary)

        case (.success, .filled):
            return (colors.success, colors.success, colors.textOnPrimary, colors.textOnPrimary)
        case (.success, .outlined):
            return (colors.background, colors.success, colors.success, colors.success)
        case (.success, .minimal):
            return (colors.success.opacity(0.1), Color.clear, colors.success, colors.success)

        case (.warning, .filled):
            return (colors.warning, colors.warning, colors.textOnPrimary, colors.textOnPrimary)
        case (.warning, .outlined):
            return (colors.background, colors.warning, colors.warning, colors.warning)
        case (.warning, .minimal):
            return (colors.warning.opacity(0.1), Color.clear, colors.warning, colors.warning)

        case (.error, .filled):
            return (colors.error, colors.error, colors.textOnPrimary, colors.textOnPrimary)
        case (.error, .outlined):
            return (colors.background, colors.error, colors.error, colors.error)
        case (.error, .minimal):
            return (colors.error.opacity(0.1), Color.clear, colors.error, colors.error)
        }
    }

    private var alertIcon: String {
        switch type {
        case .info: return "info.circle.fill"
        case .success: return "checkmark.circle.fill"
        case .warning: return "exclamationmark.triangle.fill"
        case .error: return "xmark.circle.fill"
        }
    }

    private var iconSize: CGFloat {
        switch size {
        case .small: return 16
        case .medium: return 20
        case .large: return 24
        }
    }

    private var titleFont: Font {
        switch size {
        case .small: return .system(size: 14, weight: .semibold)
        case .medium: return .system(size: 16, weight: .semibold)
        case .large: return .system(size: 18, weight: .semibold)
        }
    }

    private var messageFont: Font {
        switch size {
        case .small: return .system(size: 12)
        case .medium: return .system(size: 14)
        case .large: return .system(size: 16)
        }
    }

    private var alertPadding: CGFloat {
        switch size {
        case .small: return 12
        case .medium: return 16
        case .large: return 20
        }
    }

    // MARK: - Initializer
    public init(
        isPresented: Binding<Bool>,
        type: AlertType,
        variant: AlertVariant = .filled,
        size: AlertSize = .medium,
        title: String? = nil,
        message: String,
        isDismissible: Bool = true,
        showIcon: Bool = true,
        actions: [AlertAction] = [],
        onDismiss: (() -> Void)? = nil
    ) {
        self._isPresented = isPresented
        self.type = type
        self.variant = variant
        self.size = size
        self.title = title
        self.message = message
        self.isDismissible = isDismissible
        self.showIcon = showIcon
        self.actions = actions
        self.onDismiss = onDismiss
    }

    // MARK: - Body
    public var body: some View {
        if isPresented {
            alertContent
                .transition(.asymmetric(
                    insertion: .opacity.combined(with: .move(edge: .top)),
                    removal: .opacity.combined(with: .move(edge: .top))
                ))
                .onAppear {
                    withAnimation(.easeOut(duration: 0.3)) {
                        animationOffset = 0
                        animationOpacity = 1
                    }
                }
        }
    }

    // MARK: - Alert Content
    @ViewBuilder
    private var alertContent: some View {
        HStack(alignment: .top, spacing: Spacing.sm) {
            // Icon
            if showIcon {
                Image(systemName: alertIcon)
                    .font(.system(size: iconSize))
                    .foregroundColor(alertColors.icon)
                    .padding(.top, 2)
            }

            // Content
            VStack(alignment: .leading, spacing: Spacing.xs) {
                // Title
                if let title = title {
                    Text(title)
                        .font(titleFont)
                        .foregroundColor(alertColors.text)
                }

                // Message
                Text(message)
                    .font(messageFont)
                    .foregroundColor(alertColors.text)
                    .fixedSize(horizontal: false, vertical: true)

                // Actions
                if !actions.isEmpty {
                    HStack(spacing: Spacing.sm) {
                        ForEach(Array(actions.enumerated()), id: \.offset) { index, action in
                            alertActionButton(action)
                        }
                    }
                    .padding(.top, Spacing.xs)
                }
            }
            .frame(maxWidth: .infinity, alignment: .leading)

            // Dismiss button
            if isDismissible {
                Button(action: dismissAlert) {
                    Image(systemName: "xmark")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(alertColors.text.opacity(0.7))
                }
                .padding(.top, 2)
            }
        }
        .padding(alertPadding)
        .background(alertColors.background)
        .overlay(
            RoundedRectangle(cornerRadius: 8)
                .stroke(alertColors.border, lineWidth: variant == .outlined ? 1 : 0)
        )
        .cornerRadius(8)
        .shadow(color: colors.shadowMedium, radius: 4, y: 2)
        .offset(y: animationOffset)
        .opacity(animationOpacity)
    }

    // MARK: - Alert Action Button
    @ViewBuilder
    private func alertActionButton(_ action: AlertAction) -> some View {
        Button(action: {
            action.action()
            // Haptic feedback
            let impactFeedback = UIImpactFeedbackGenerator(style: .light)
            impactFeedback.impactOccurred()
        }) {
            Text(action.title)
                .font(.system(size: messageFont.pointSize - 1, weight: .medium))
                .foregroundColor(actionButtonColor(action.style))
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(actionButtonBackground(action.style))
                .cornerRadius(4)
        }
    }

    // MARK: - Action Button Styling
    private func actionButtonColor(_ style: AlertAction.AlertActionStyle) -> Color {
        switch (style, variant) {
        case (.primary, .filled):
            return alertColors.background == colors.background ? alertColors.text : colors.textOnPrimary
        case (.primary, _):
            return alertColors.text
        case (.secondary, _):
            return alertColors.text.opacity(0.8)
        case (.destructive, _):
            return colors.error
        }
    }

    private func actionButtonBackground(_ style: AlertAction.AlertActionStyle) -> Color {
        switch (style, variant) {
        case (.primary, .filled):
            return alertColors.background == colors.background ? alertColors.text.opacity(0.1) : colors.background.opacity(0.2)
        case (.primary, _):
            return alertColors.text.opacity(0.1)
        case (.secondary, _):
            return alertColors.text.opacity(0.05)
        case (.destructive, _):
            return colors.error.opacity(0.1)
        }
    }

    // MARK: - Actions
    private func dismissAlert() {
        withAnimation(.easeIn(duration: 0.2)) {
            animationOffset = -30
            animationOpacity = 0
        }

        DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
            isPresented = false
            onDismiss?()
        }

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Alert Manager
@MainActor
public class TchatAlertManager: ObservableObject {
    @Published public var currentAlert: TchatAlertItem?

    public static let shared = TchatAlertManager()

    private init() {}

    public func show(
        type: TchatAlert.AlertType,
        variant: TchatAlert.AlertVariant = .filled,
        title: String? = nil,
        message: String,
        actions: [TchatAlert.AlertAction] = [],
        duration: TimeInterval? = nil
    ) {
        let alert = TchatAlertItem(
            type: type,
            variant: variant,
            title: title,
            message: message,
            actions: actions
        )

        currentAlert = alert

        if let duration = duration {
            DispatchQueue.main.asyncAfter(deadline: .now() + duration) {
                if self.currentAlert?.id == alert.id {
                    self.dismiss()
                }
            }
        }
    }

    public func dismiss() {
        currentAlert = nil
    }
}

// MARK: - Alert Item
public struct TchatAlertItem: Identifiable {
    public let id = UUID()
    let type: TchatAlert.AlertType
    let variant: TchatAlert.AlertVariant
    let title: String?
    let message: String
    let actions: [TchatAlert.AlertAction]

    public init(
        type: TchatAlert.AlertType,
        variant: TchatAlert.AlertVariant = .filled,
        title: String? = nil,
        message: String,
        actions: [TchatAlert.AlertAction] = []
    ) {
        self.type = type
        self.variant = variant
        self.title = title
        self.message = message
        self.actions = actions
    }
}

// MARK: - Alert Overlay
public struct TchatAlertOverlay: View {
    @ObservedObject private var alertManager = TchatAlertManager.shared

    public init() {}

    public var body: some View {
        VStack {
            if let alert = alertManager.currentAlert {
                TchatAlert(
                    isPresented: .constant(true),
                    type: alert.type,
                    variant: alert.variant,
                    title: alert.title,
                    message: alert.message,
                    actions: alert.actions
                ) {
                    alertManager.dismiss()
                }
                .padding(.horizontal, Spacing.md)
                .padding(.top, Spacing.md)
            }

            Spacer()
        }
        .animation(.easeInOut(duration: 0.3), value: alertManager.currentAlert?.id)
    }
}

// MARK: - Convenience Methods
extension TchatAlertManager {
    public func showSuccess(
        title: String? = nil,
        message: String,
        duration: TimeInterval = 3.0
    ) {
        show(
            type: .success,
            title: title,
            message: message,
            duration: duration
        )
    }

    public func showError(
        title: String? = nil,
        message: String,
        actions: [TchatAlert.AlertAction] = []
    ) {
        show(
            type: .error,
            title: title,
            message: message,
            actions: actions
        )
    }

    public func showWarning(
        title: String? = nil,
        message: String,
        actions: [TchatAlert.AlertAction] = []
    ) {
        show(
            type: .warning,
            title: title,
            message: message,
            actions: actions
        )
    }

    public func showInfo(
        title: String? = nil,
        message: String,
        duration: TimeInterval = 5.0
    ) {
        show(
            type: .info,
            title: title,
            message: message,
            duration: duration
        )
    }
}

// MARK: - Preview
#if DEBUG
struct TchatAlert_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Success alert
                TchatAlert(
                    isPresented: .constant(true),
                    type: .success,
                    variant: .filled,
                    title: "Success",
                    message: "Your changes have been saved successfully."
                )

                // Error alert with actions
                TchatAlert(
                    isPresented: .constant(true),
                    type: .error,
                    variant: .outlined,
                    title: "Error",
                    message: "Failed to save changes. Please try again.",
                    actions: [
                        TchatAlert.AlertAction(title: "Retry", style: .primary) { },
                        TchatAlert.AlertAction(title: "Cancel", style: .secondary) { }
                    ]
                )

                // Warning alert
                TchatAlert(
                    isPresented: .constant(true),
                    type: .warning,
                    variant: .minimal,
                    title: "Warning",
                    message: "This action cannot be undone. Are you sure you want to continue?",
                    actions: [
                        TchatAlert.AlertAction(title: "Continue", style: .destructive) { },
                        TchatAlert.AlertAction(title: "Cancel", style: .secondary) { }
                    ]
                )

                // Info alert
                TchatAlert(
                    isPresented: .constant(true),
                    type: .info,
                    variant: .filled,
                    message: "New features are available. Update your app to get the latest improvements.",
                    size: .small
                )

                // Minimal alert without icon
                TchatAlert(
                    isPresented: .constant(true),
                    type: .info,
                    variant: .minimal,
                    message: "This is a minimal alert without an icon.",
                    showIcon: false,
                    isDismissible: false
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif