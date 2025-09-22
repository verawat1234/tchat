//
//  TchatModal.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Modal component following Tchat design system
public struct TchatModal<Content: View>: View {

    // MARK: - Modal Types
    public enum ModalSize {
        case small
        case medium
        case large
        case fullScreen
        case custom(width: CGFloat?, height: CGFloat?)
    }

    public enum ModalPosition {
        case center
        case bottom
        case top
    }

    public enum ModalAnimation {
        case slide
        case fade
        case scale
    }

    // MARK: - Properties
    @Binding private var isPresented: Bool
    @State private var dragOffset: CGSize = .zero
    @State private var backgroundOpacity: Double = 0

    let size: ModalSize
    let position: ModalPosition
    let animation: ModalAnimation
    let showCloseButton: Bool
    let isDismissible: Bool
    let allowDragDismiss: Bool
    let showOverlay: Bool
    let overlayColor: Color
    let cornerRadius: CGFloat
    let content: Content
    let onDismiss: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    private let dismissThreshold: CGFloat = 150

    // MARK: - Computed Properties
    private var modalDimensions: (width: CGFloat?, height: CGFloat?) {
        switch size {
        case .small:
            return (300, 200)
        case .medium:
            return (400, 300)
        case .large:
            return (500, 400)
        case .fullScreen:
            return (nil, nil)
        case .custom(let width, let height):
            return (width, height)
        }
    }

    private var modalOffset: CGSize {
        switch position {
        case .center:
            return dragOffset
        case .bottom:
            return CGSize(width: dragOffset.width, height: max(0, dragOffset.height))
        case .top:
            return CGSize(width: dragOffset.width, height: min(0, dragOffset.height))
        }
    }

    private var insertionEdge: Edge {
        switch position {
        case .center: return .bottom
        case .bottom: return .bottom
        case .top: return .top
        }
    }

    private var removalEdge: Edge {
        switch position {
        case .center: return .bottom
        case .bottom: return .bottom
        case .top: return .top
        }
    }

    // MARK: - Initializer
    public init(
        isPresented: Binding<Bool>,
        size: ModalSize = .medium,
        position: ModalPosition = .center,
        animation: ModalAnimation = .slide,
        showCloseButton: Bool = true,
        isDismissible: Bool = true,
        allowDragDismiss: Bool = true,
        showOverlay: Bool = true,
        overlayColor: Color = Color.black.opacity(0.5),
        cornerRadius: CGFloat = 16,
        onDismiss: (() -> Void)? = nil,
        @ViewBuilder content: () -> Content
    ) {
        self._isPresented = isPresented
        self.size = size
        self.position = position
        self.animation = animation
        self.showCloseButton = showCloseButton
        self.isDismissible = isDismissible
        self.allowDragDismiss = allowDragDismiss
        self.showOverlay = showOverlay
        self.overlayColor = overlayColor
        self.cornerRadius = cornerRadius
        self.content = content()
        self.onDismiss = onDismiss
    }

    // MARK: - Body
    public var body: some View {
        if isPresented {
            GeometryReader { geometry in
                ZStack {
                    // Overlay
                    if showOverlay {
                        overlayColor
                            .opacity(backgroundOpacity)
                            .ignoresSafeArea()
                            .onTapGesture {
                                if isDismissible {
                                    dismissModal()
                                }
                            }
                    }

                    // Modal content
                    modalContent(geometry: geometry)
                        .offset(modalOffset)
                        .gesture(
                            allowDragDismiss && isDismissible ? dragGesture : nil
                        )
                }
            }
            .transition(modalTransition)
            .onAppear {
                withAnimation(.easeOut(duration: 0.3)) {
                    backgroundOpacity = 1
                }
            }
        }
    }

    // MARK: - Modal Content
    @ViewBuilder
    private func modalContent(geometry: GeometryProxy) -> some View {
        VStack(spacing: 0) {
            // Close button (if enabled)
            if showCloseButton {
                HStack {
                    Spacer()

                    Button(action: dismissModal) {
                        Image(systemName: "xmark")
                            .font(.system(size: 16, weight: .medium))
                            .foregroundColor(colors.textSecondary)
                            .padding(8)
                            .background(colors.surface)
                            .clipShape(Circle())
                    }
                    .padding(.trailing, Spacing.md)
                    .padding(.top, Spacing.md)
                }
                .zIndex(1)
            }

            // Content
            content
                .frame(
                    width: modalWidth(geometry: geometry),
                    height: modalHeight(geometry: geometry),
                    alignment: .center
                )
                .background(colors.background)
                .cornerRadius(cornerRadius)
                .shadow(color: colors.shadowMedium, radius: 20, y: 10)
        }
        .frame(
            maxWidth: .infinity,
            maxHeight: .infinity,
            alignment: modalAlignment
        )
        .padding(size == .fullScreen ? 0 : Spacing.lg)
    }

    // MARK: - Computed Properties
    private func modalWidth(geometry: GeometryProxy) -> CGFloat? {
        if size == .fullScreen {
            return geometry.size.width
        }

        guard let width = modalDimensions.width else {
            return min(geometry.size.width - (Spacing.lg * 2), 600)
        }

        return min(width, geometry.size.width - (Spacing.lg * 2))
    }

    private func modalHeight(geometry: GeometryProxy) -> CGFloat? {
        if size == .fullScreen {
            return geometry.size.height
        }

        guard let height = modalDimensions.height else {
            return nil
        }

        return min(height, geometry.size.height - (Spacing.lg * 2))
    }

    private var modalAlignment: Alignment {
        switch position {
        case .center: return .center
        case .bottom: return .bottom
        case .top: return .top
        }
    }

    private var modalTransition: AnyTransition {
        switch animation {
        case .slide:
            return .asymmetric(
                insertion: .move(edge: insertionEdge).combined(with: .opacity),
                removal: .move(edge: removalEdge).combined(with: .opacity)
            )
        case .fade:
            return .opacity
        case .scale:
            return .scale.combined(with: .opacity)
        }
    }

    // MARK: - Drag Gesture
    private var dragGesture: some Gesture {
        DragGesture()
            .onChanged { value in
                dragOffset = value.translation
            }
            .onEnded { value in
                let shouldDismiss: Bool

                switch position {
                case .center:
                    shouldDismiss = abs(value.translation.y) > dismissThreshold ||
                                  abs(value.translation.x) > dismissThreshold
                case .bottom:
                    shouldDismiss = value.translation.y > dismissThreshold
                case .top:
                    shouldDismiss = value.translation.y < -dismissThreshold
                }

                if shouldDismiss {
                    dismissModal()
                } else {
                    withAnimation(.spring()) {
                        dragOffset = .zero
                    }
                }
            }
    }

    // MARK: - Actions
    private func dismissModal() {
        withAnimation(.easeIn(duration: 0.2)) {
            backgroundOpacity = 0

            switch position {
            case .center:
                dragOffset = CGSize(width: 0, height: 300)
            case .bottom:
                dragOffset = CGSize(width: 0, height: 300)
            case .top:
                dragOffset = CGSize(width: 0, height: -300)
            }
        }

        DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
            isPresented = false
            dragOffset = .zero
            onDismiss?()
        }

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Modal Manager
@MainActor
public class TchatModalManager: ObservableObject {
    @Published public var currentModal: TchatModalItem?

    public static let shared = TchatModalManager()

    private init() {}

    public func present<Content: View>(
        size: TchatModal<Content>.ModalSize = .medium,
        position: TchatModal<Content>.ModalPosition = .center,
        animation: TchatModal<Content>.ModalAnimation = .slide,
        showCloseButton: Bool = true,
        isDismissible: Bool = true,
        allowDragDismiss: Bool = true,
        showOverlay: Bool = true,
        onDismiss: (() -> Void)? = nil,
        @ViewBuilder content: () -> Content
    ) {
        let modal = TchatModalItem(
            size: size,
            position: position,
            animation: animation,
            showCloseButton: showCloseButton,
            isDismissible: isDismissible,
            allowDragDismiss: allowDragDismiss,
            showOverlay: showOverlay,
            content: AnyView(content()),
            onDismiss: onDismiss
        )

        currentModal = modal
    }

    public func dismiss() {
        currentModal = nil
    }
}

// MARK: - Modal Item
public struct TchatModalItem: Identifiable {
    public let id = UUID()
    let size: TchatModal<AnyView>.ModalSize
    let position: TchatModal<AnyView>.ModalPosition
    let animation: TchatModal<AnyView>.ModalAnimation
    let showCloseButton: Bool
    let isDismissible: Bool
    let allowDragDismiss: Bool
    let showOverlay: Bool
    let content: AnyView
    let onDismiss: (() -> Void)?

    public init<Content: View>(
        size: TchatModal<Content>.ModalSize,
        position: TchatModal<Content>.ModalPosition,
        animation: TchatModal<Content>.ModalAnimation,
        showCloseButton: Bool,
        isDismissible: Bool,
        allowDragDismiss: Bool,
        showOverlay: Bool,
        content: AnyView,
        onDismiss: (() -> Void)?
    ) {
        self.size = size as! TchatModal<AnyView>.ModalSize
        self.position = position as! TchatModal<AnyView>.ModalPosition
        self.animation = animation as! TchatModal<AnyView>.ModalAnimation
        self.showCloseButton = showCloseButton
        self.isDismissible = isDismissible
        self.allowDragDismiss = allowDragDismiss
        self.showOverlay = showOverlay
        self.content = content
        self.onDismiss = onDismiss
    }
}

// MARK: - Modal Overlay
public struct TchatModalOverlay: View {
    @ObservedObject private var modalManager = TchatModalManager.shared

    public init() {}

    public var body: some View {
        ZStack {
            if let modal = modalManager.currentModal {
                TchatModal(
                    isPresented: .constant(true),
                    size: modal.size,
                    position: modal.position,
                    animation: modal.animation,
                    showCloseButton: modal.showCloseButton,
                    isDismissible: modal.isDismissible,
                    allowDragDismiss: modal.allowDragDismiss,
                    showOverlay: modal.showOverlay,
                    onDismiss: {
                        modal.onDismiss?()
                        modalManager.dismiss()
                    }
                ) {
                    modal.content
                }
            }
        }
        .animation(.easeInOut(duration: 0.3), value: modalManager.currentModal?.id)
    }
}

// MARK: - Convenience Methods
extension TchatModalManager {
    public func presentAlert(
        title: String,
        message: String,
        primaryButton: String = "OK",
        secondaryButton: String? = nil,
        onPrimary: (() -> Void)? = nil,
        onSecondary: (() -> Void)? = nil
    ) {
        present(size: .small, position: .center) {
            VStack(spacing: Spacing.lg) {
                VStack(spacing: Spacing.sm) {
                    Text(title)
                        .font(.headline)
                        .foregroundColor(Colors().textPrimary)

                    Text(message)
                        .font(.body)
                        .foregroundColor(Colors().textSecondary)
                        .multilineTextAlignment(.center)
                }

                HStack(spacing: Spacing.sm) {
                    if let secondaryButton = secondaryButton {
                        Button(secondaryButton) {
                            onSecondary?()
                            TchatModalManager.shared.dismiss()
                        }
                        .foregroundColor(Colors().textSecondary)
                    }

                    Button(primaryButton) {
                        onPrimary?()
                        TchatModalManager.shared.dismiss()
                    }
                    .foregroundColor(Colors().primary)
                    .fontWeight(.semibold)
                }
            }
            .padding(Spacing.lg)
        }
    }

    public func presentBottomSheet<Content: View>(
        @ViewBuilder content: () -> Content
    ) {
        present(
            size: .custom(width: nil, height: 400),
            position: .bottom,
            animation: .slide,
            showCloseButton: false,
            allowDragDismiss: true
        ) {
            content()
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatModal_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color(.systemGray6)
                .ignoresSafeArea()

            VStack(spacing: 20) {
                Text("Background Content")
                    .font(.title)
                    .foregroundColor(.primary)
            }

            TchatModal(
                isPresented: .constant(true),
                size: .medium,
                position: .center,
                animation: .slide
            ) {
                VStack(spacing: Spacing.lg) {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.system(size: 48))
                        .foregroundColor(.green)

                    VStack(spacing: Spacing.sm) {
                        Text("Success!")
                            .font(.title2)
                            .fontWeight(.bold)

                        Text("Your changes have been saved successfully.")
                            .font(.body)
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                    }

                    HStack(spacing: Spacing.sm) {
                        Button("Cancel") {
                            // Handle cancel
                        }
                        .foregroundColor(.secondary)

                        Button("Continue") {
                            // Handle continue
                        }
                        .foregroundColor(.primary)
                        .fontWeight(.semibold)
                    }
                }
                .padding(Spacing.lg)
            }
        }
        .previewLayout(.sizeThatFits)
    }
}
#endif