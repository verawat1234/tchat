//
//  TchatToast.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Toast notification component following Tchat design system
public struct TchatToast: View {

    // MARK: - Toast Types
    public enum ToastType {
        case info
        case success
        case warning
        case error
    }

    public enum ToastPosition {
        case top
        case bottom
        case center
    }

    public enum ToastStyle {
        case filled
        case outlined
        case minimal
    }

    // MARK: - Properties
    @Binding private var isPresented: Bool
    @State private var dragOffset: CGSize = .zero
    @State private var animationOffset: CGFloat = 0

    let type: ToastType
    let position: ToastPosition
    let style: ToastStyle
    let message: String
    let icon: String?
    let duration: TimeInterval
    let isDismissible: Bool
    let hapticFeedback: Bool
    let onDismiss: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    private let dismissThreshold: CGFloat = 100

    // MARK: - Computed Properties
    private var toastColors: (background: Color, border: Color, text: Color, icon: Color) {
        switch (type, style) {
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

    private var defaultIcon: String {
        switch type {
        case .info: return "info.circle.fill"
        case .success: return "checkmark.circle.fill"
        case .warning: return "exclamationmark.triangle.fill"
        case .error: return "xmark.circle.fill"
        }
    }

    private var displayIcon: String {
        icon ?? defaultIcon
    }

    private var initialOffset: CGFloat {
        switch position {
        case .top: return -200
        case .bottom: return 200
        case .center: return 0
        }
    }

    // MARK: - Initializer
    public init(
        isPresented: Binding<Bool>,
        type: ToastType,
        position: ToastPosition = .top,
        style: ToastStyle = .filled,
        message: String,
        icon: String? = nil,
        duration: TimeInterval = 3.0,
        isDismissible: Bool = true,
        hapticFeedback: Bool = true,
        onDismiss: (() -> Void)? = nil
    ) {
        self._isPresented = isPresented
        self.type = type
        self.position = position
        self.style = style
        self.message = message
        self.icon = icon
        self.duration = duration
        self.isDismissible = isDismissible
        self.hapticFeedback = hapticFeedback
        self.onDismiss = onDismiss
    }

    // MARK: - Body
    public var body: some View {
        if isPresented {
            toastContent
                .transition(.asymmetric(
                    insertion: .move(edge: position == .top ? .top : .bottom).combined(with: .opacity),
                    removal: .move(edge: position == .top ? .top : .bottom).combined(with: .opacity)
                ))
                .onAppear {
                    if hapticFeedback {
                        triggerHapticFeedback()
                    }

                    withAnimation(.easeOut(duration: 0.4)) {
                        animationOffset = 0
                    }

                    // Auto-dismiss after duration
                    if duration > 0 {
                        DispatchQueue.main.asyncAfter(deadline: .now() + duration) {
                            dismissToast()
                        }
                    }
                }
        }
    }

    // MARK: - Toast Content
    @ViewBuilder
    private var toastContent: some View {
        HStack(spacing: Spacing.sm) {
            // Icon
            Image(systemName: displayIcon)
                .font(.system(size: 16, weight: .semibold))
                .foregroundColor(toastColors.icon)

            // Message
            Text(message)
                .font(.system(size: 14, weight: .medium))
                .foregroundColor(toastColors.text)
                .multilineTextAlignment(.leading)
                .frame(maxWidth: .infinity, alignment: .leading)

            // Dismiss button (optional)
            if isDismissible {
                Button(action: dismissToast) {
                    Image(systemName: "xmark")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(toastColors.text.opacity(0.7))
                }
            }
        }
        .padding(.horizontal, Spacing.md)
        .padding(.vertical, Spacing.sm)
        .background(toastColors.background)
        .overlay(
            RoundedRectangle(cornerRadius: 8)
                .stroke(toastColors.border, lineWidth: style == .outlined ? 1 : 0)
        )
        .cornerRadius(8)
        .shadow(color: colors.shadowMedium, radius: 8, y: 4)
        .padding(.horizontal, Spacing.md)
        .offset(x: dragOffset.width, y: animationOffset + dragOffset.height)
        .gesture(
            isDismissible ? dragGesture : nil
        )
    }

    // MARK: - Drag Gesture
    private var dragGesture: some Gesture {
        DragGesture()
            .onChanged { value in
                dragOffset = value.translation
            }
            .onEnded { value in
                let shouldDismiss = abs(value.translation.width) > dismissThreshold ||
                                  abs(value.translation.height) > dismissThreshold

                if shouldDismiss {
                    withAnimation(.easeIn(duration: 0.2)) {
                        dragOffset = CGSize(
                            width: value.translation.width > 0 ? 400 : -400,
                            height: value.translation.height
                        )
                    }

                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
                        dismissToast()
                    }
                } else {
                    withAnimation(.spring()) {
                        dragOffset = .zero
                    }
                }
            }
    }

    // MARK: - Actions
    private func dismissToast() {
        withAnimation(.easeIn(duration: 0.3)) {
            animationOffset = initialOffset
        }

        DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
            isPresented = false
            onDismiss?()
        }
    }

    private func triggerHapticFeedback() {
        let feedbackStyle: UIImpactFeedbackGenerator.FeedbackStyle

        switch type {
        case .success:
            feedbackStyle = .light
        case .warning:
            feedbackStyle = .medium
        case .error:
            feedbackStyle = .heavy
        case .info:
            feedbackStyle = .light
        }

        let impactFeedback = UIImpactFeedbackGenerator(style: feedbackStyle)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Toast Manager
@MainActor
public class TchatToastManager: ObservableObject {
    @Published public var currentToasts: [TchatToastItem] = []

    public static let shared = TchatToastManager()

    private init() {}

    public func show(
        type: TchatToast.ToastType,
        position: TchatToast.ToastPosition = .top,
        style: TchatToast.ToastStyle = .filled,
        message: String,
        icon: String? = nil,
        duration: TimeInterval = 3.0,
        isDismissible: Bool = true,
        hapticFeedback: Bool = true
    ) {
        let toast = TchatToastItem(
            type: type,
            position: position,
            style: style,
            message: message,
            icon: icon,
            duration: duration,
            isDismissible: isDismissible,
            hapticFeedback: hapticFeedback
        )

        currentToasts.append(toast)

        // Auto-dismiss after duration
        if duration > 0 {
            DispatchQueue.main.asyncAfter(deadline: .now() + duration + 0.5) {
                self.dismiss(toast.id)
            }
        }
    }

    public func dismiss(_ id: UUID) {
        currentToasts.removeAll { $0.id == id }
    }

    public func dismissAll() {
        currentToasts.removeAll()
    }
}

// MARK: - Toast Item
public struct TchatToastItem: Identifiable {
    public let id = UUID()
    let type: TchatToast.ToastType
    let position: TchatToast.ToastPosition
    let style: TchatToast.ToastStyle
    let message: String
    let icon: String?
    let duration: TimeInterval
    let isDismissible: Bool
    let hapticFeedback: Bool

    public init(
        type: TchatToast.ToastType,
        position: TchatToast.ToastPosition = .top,
        style: TchatToast.ToastStyle = .filled,
        message: String,
        icon: String? = nil,
        duration: TimeInterval = 3.0,
        isDismissible: Bool = true,
        hapticFeedback: Bool = true
    ) {
        self.type = type
        self.position = position
        self.style = style
        self.message = message
        self.icon = icon
        self.duration = duration
        self.isDismissible = isDismissible
        self.hapticFeedback = hapticFeedback
    }
}

// MARK: - Toast Overlay
public struct TchatToastOverlay: View {
    @ObservedObject private var toastManager = TchatToastManager.shared

    public init() {}

    public var body: some View {
        GeometryReader { geometry in
            ZStack {
                // Top toasts
                VStack(spacing: Spacing.xs) {
                    ForEach(toastManager.currentToasts.filter { $0.position == .top }) { toast in
                        TchatToast(
                            isPresented: .constant(true),
                            type: toast.type,
                            position: toast.position,
                            style: toast.style,
                            message: toast.message,
                            icon: toast.icon,
                            duration: 0, // Managed by toast manager
                            isDismissible: toast.isDismissible,
                            hapticFeedback: false // Already triggered
                        ) {
                            toastManager.dismiss(toast.id)
                        }
                    }

                    Spacer()
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .top)
                .padding(.top, geometry.safeAreaInsets.top + Spacing.md)

                // Center toasts
                VStack(spacing: Spacing.xs) {
                    ForEach(toastManager.currentToasts.filter { $0.position == .center }) { toast in
                        TchatToast(
                            isPresented: .constant(true),
                            type: toast.type,
                            position: toast.position,
                            style: toast.style,
                            message: toast.message,
                            icon: toast.icon,
                            duration: 0,
                            isDismissible: toast.isDismissible,
                            hapticFeedback: false
                        ) {
                            toastManager.dismiss(toast.id)
                        }
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)

                // Bottom toasts
                VStack(spacing: Spacing.xs) {
                    Spacer()

                    ForEach(toastManager.currentToasts.filter { $0.position == .bottom }.reversed()) { toast in
                        TchatToast(
                            isPresented: .constant(true),
                            type: toast.type,
                            position: toast.position,
                            style: toast.style,
                            message: toast.message,
                            icon: toast.icon,
                            duration: 0,
                            isDismissible: toast.isDismissible,
                            hapticFeedback: false
                        ) {
                            toastManager.dismiss(toast.id)
                        }
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .bottom)
                .padding(.bottom, geometry.safeAreaInsets.bottom + Spacing.md)
            }
        }
        .allowsHitTesting(false)
        .animation(.easeInOut(duration: 0.3), value: toastManager.currentToasts.count)
    }
}

// MARK: - Convenience Methods
extension TchatToastManager {
    public func showSuccess(
        _ message: String,
        position: TchatToast.ToastPosition = .top,
        duration: TimeInterval = 2.0
    ) {
        show(
            type: .success,
            position: position,
            message: message,
            duration: duration
        )
    }

    public func showError(
        _ message: String,
        position: TchatToast.ToastPosition = .top,
        duration: TimeInterval = 4.0
    ) {
        show(
            type: .error,
            position: position,
            message: message,
            duration: duration
        )
    }

    public func showWarning(
        _ message: String,
        position: TchatToast.ToastPosition = .top,
        duration: TimeInterval = 3.0
    ) {
        show(
            type: .warning,
            position: position,
            message: message,
            duration: duration
        )
    }

    public func showInfo(
        _ message: String,
        position: TchatToast.ToastPosition = .top,
        duration: TimeInterval = 3.0
    ) {
        show(
            type: .info,
            position: position,
            message: message,
            duration: duration
        )
    }
}

// MARK: - Preview
#if DEBUG
struct TchatToast_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color(.systemGray6)
                .ignoresSafeArea()

            VStack(spacing: Spacing.lg) {
                // Success toast
                TchatToast(
                    isPresented: .constant(true),
                    type: .success,
                    position: .top,
                    style: .filled,
                    message: "Changes saved successfully!",
                    duration: 0
                )

                // Error toast
                TchatToast(
                    isPresented: .constant(true),
                    type: .error,
                    position: .top,
                    style: .outlined,
                    message: "Failed to upload file. Please try again.",
                    duration: 0
                )

                // Warning toast
                TchatToast(
                    isPresented: .constant(true),
                    type: .warning,
                    position: .top,
                    style: .minimal,
                    message: "You have unsaved changes.",
                    duration: 0
                )

                // Info toast with custom icon
                TchatToast(
                    isPresented: .constant(true),
                    type: .info,
                    position: .top,
                    style: .filled,
                    message: "New update available for download.",
                    icon: "arrow.down.circle.fill",
                    duration: 0,
                    isDismissible: false
                )
            }
            .padding()
        }
        .previewLayout(.sizeThatFits)
    }
}
#endif