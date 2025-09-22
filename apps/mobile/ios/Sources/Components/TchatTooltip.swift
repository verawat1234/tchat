//
//  TchatTooltip.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Tooltip component following Tchat design system
public struct TchatTooltip<Content: View>: View {

    // MARK: - Tooltip Types
    public enum TooltipPosition {
        case top
        case bottom
        case leading
        case trailing
        case auto
    }

    public enum TooltipTrigger {
        case tap
        case longPress
        case hover
        case manual
    }

    public enum TooltipStyle {
        case dark
        case light
        case info
        case warning
        case error
    }

    // MARK: - Properties
    @State private var isVisible: Bool = false
    @State private var tooltipFrame: CGRect = .zero
    @State private var targetFrame: CGRect = .zero
    @State private var calculatedPosition: TooltipPosition = .top

    let content: Content
    let text: String
    let position: TooltipPosition
    let trigger: TooltipTrigger
    let style: TooltipStyle
    let maxWidth: CGFloat
    let showArrow: Bool
    let delay: TimeInterval
    let autoDismissDelay: TimeInterval?
    let onShow: (() -> Void)?
    let onDismiss: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    private let arrowSize: CGFloat = 8

    // MARK: - Computed Properties
    private var tooltipColors: (background: Color, text: Color, border: Color) {
        switch style {
        case .dark:
            return (Color.black.opacity(0.9), Color.white, Color.clear)
        case .light:
            return (colors.background, colors.textPrimary, colors.border)
        case .info:
            return (colors.primary, colors.textOnPrimary, colors.primary)
        case .warning:
            return (colors.warning, colors.textOnPrimary, colors.warning)
        case .error:
            return (colors.error, colors.textOnPrimary, colors.error)
        }
    }

    // MARK: - Initializer
    public init(
        text: String,
        position: TooltipPosition = .auto,
        trigger: TooltipTrigger = .longPress,
        style: TooltipStyle = .dark,
        maxWidth: CGFloat = 250,
        showArrow: Bool = true,
        delay: TimeInterval = 0.5,
        autoDismissDelay: TimeInterval? = 3.0,
        onShow: (() -> Void)? = nil,
        onDismiss: (() -> Void)? = nil,
        @ViewBuilder content: () -> Content
    ) {
        self.content = content()
        self.text = text
        self.position = position
        self.trigger = trigger
        self.style = style
        self.maxWidth = maxWidth
        self.showArrow = showArrow
        self.delay = delay
        self.autoDismissDelay = autoDismissDelay
        self.onShow = onShow
        self.onDismiss = onDismiss
    }

    // MARK: - Body
    public var body: some View {
        ZStack {
            content
                .background(
                    GeometryReader { geometry in
                        Color.clear
                            .onAppear {
                                targetFrame = geometry.frame(in: .global)
                            }
                            .onChange(of: geometry.frame(in: .global)) { newFrame in
                                targetFrame = newFrame
                            }
                    }
                )
                .onTapGesture {
                    if trigger == .tap {
                        handleTrigger()
                    }
                }
                .onLongPressGesture(minimumDuration: delay) {
                    if trigger == .longPress {
                        handleTrigger()
                    }
                }

            if isVisible {
                tooltipView
            }
        }
    }

    // MARK: - Tooltip View
    @ViewBuilder
    private var tooltipView: some View {
        GeometryReader { geometry in
            VStack(spacing: 0) {
                // Arrow (top)
                if showArrow && (calculatedPosition == .bottom) {
                    arrowShape
                        .fill(tooltipColors.background)
                        .frame(width: arrowSize * 2, height: arrowSize)
                }

                // Tooltip content
                Text(text)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(tooltipColors.text)
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: maxWidth)
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                    .background(tooltipColors.background)
                    .overlay(
                        RoundedRectangle(cornerRadius: 8)
                            .stroke(tooltipColors.border, lineWidth: style == .light ? 1 : 0)
                    )
                    .cornerRadius(8)
                    .shadow(color: colors.shadowMedium, radius: 8, y: 4)
                    .background(
                        GeometryReader { tooltipGeometry in
                            Color.clear
                                .onAppear {
                                    tooltipFrame = tooltipGeometry.frame(in: .global)
                                }
                        }
                    )

                // Arrow (bottom)
                if showArrow && (calculatedPosition == .top) {
                    arrowShape
                        .fill(tooltipColors.background)
                        .frame(width: arrowSize * 2, height: arrowSize)
                        .rotationEffect(.degrees(180))
                }
            }
            .position(tooltipPosition(in: geometry))
            .animation(.easeInOut(duration: 0.2), value: isVisible)
        }
        .ignoresSafeArea()
    }

    // MARK: - Arrow Shape
    @ViewBuilder
    private var arrowShape: some View {
        Path { path in
            path.move(to: CGPoint(x: 0, y: arrowSize))
            path.addLine(to: CGPoint(x: arrowSize, y: 0))
            path.addLine(to: CGPoint(x: arrowSize * 2, y: arrowSize))
            path.closeSubpath()
        }
    }

    // MARK: - Position Calculation
    private func tooltipPosition(in geometry: GeometryProxy) -> CGPoint {
        let screenBounds = geometry.frame(in: .global)
        let tooltipSize = CGSize(width: min(maxWidth + 24, screenBounds.width - 32), height: 50)

        // Calculate best position if auto
        if position == .auto {
            calculatedPosition = calculateBestPosition(
                targetFrame: targetFrame,
                tooltipSize: tooltipSize,
                screenBounds: screenBounds
            )
        } else {
            calculatedPosition = position
        }

        let spacing: CGFloat = 8
        let arrowOffset = showArrow ? arrowSize : 0

        var x: CGFloat
        var y: CGFloat

        switch calculatedPosition {
        case .top:
            x = targetFrame.midX
            y = targetFrame.minY - spacing - arrowOffset - tooltipSize.height / 2

        case .bottom:
            x = targetFrame.midX
            y = targetFrame.maxY + spacing + arrowOffset + tooltipSize.height / 2

        case .leading:
            x = targetFrame.minX - spacing - tooltipSize.width / 2
            y = targetFrame.midY

        case .trailing:
            x = targetFrame.maxX + spacing + tooltipSize.width / 2
            y = targetFrame.midY

        case .auto:
            // Fallback to top
            x = targetFrame.midX
            y = targetFrame.minY - spacing - arrowOffset - tooltipSize.height / 2
        }

        // Ensure tooltip stays within screen bounds
        x = max(tooltipSize.width / 2 + 16, min(x, screenBounds.width - tooltipSize.width / 2 - 16))
        y = max(tooltipSize.height / 2 + 16, min(y, screenBounds.height - tooltipSize.height / 2 - 16))

        return CGPoint(x: x, y: y)
    }

    private func calculateBestPosition(
        targetFrame: CGRect,
        tooltipSize: CGSize,
        screenBounds: CGRect
    ) -> TooltipPosition {
        let spacing: CGFloat = 16
        let arrowOffset = showArrow ? arrowSize : 0

        // Check available space in each direction
        let spaceTop = targetFrame.minY - spacing - arrowOffset - tooltipSize.height
        let spaceBottom = screenBounds.height - targetFrame.maxY - spacing - arrowOffset - tooltipSize.height
        let spaceLeading = targetFrame.minX - spacing - tooltipSize.width
        let spaceTrailing = screenBounds.width - targetFrame.maxX - spacing - tooltipSize.width

        // Prioritize top/bottom over leading/trailing
        if spaceTop >= 0 {
            return .top
        } else if spaceBottom >= 0 {
            return .bottom
        } else if spaceTrailing >= 0 {
            return .trailing
        } else if spaceLeading >= 0 {
            return .leading
        } else {
            // Not enough space anywhere, default to top
            return .top
        }
    }

    // MARK: - Actions
    private func handleTrigger() {
        if isVisible {
            hideTooltip()
        } else {
            showTooltip()
        }
    }

    private func showTooltip() {
        withAnimation(.easeOut(duration: 0.2)) {
            isVisible = true
        }

        onShow?()

        // Auto-dismiss
        if let autoDismissDelay = autoDismissDelay {
            DispatchQueue.main.asyncAfter(deadline: .now() + autoDismissDelay) {
                hideTooltip()
            }
        }

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }

    private func hideTooltip() {
        withAnimation(.easeIn(duration: 0.15)) {
            isVisible = false
        }

        onDismiss?()
    }
}

// MARK: - Manual Control
extension TchatTooltip {
    public func show() {
        showTooltip()
    }

    public func hide() {
        hideTooltip()
    }

    public func toggle() {
        handleTrigger()
    }
}

// MARK: - View Extension
extension View {
    public func tchatTooltip(
        _ text: String,
        position: TchatTooltip<EmptyView>.TooltipPosition = .auto,
        trigger: TchatTooltip<EmptyView>.TooltipTrigger = .longPress,
        style: TchatTooltip<EmptyView>.TooltipStyle = .dark,
        maxWidth: CGFloat = 250,
        showArrow: Bool = true,
        delay: TimeInterval = 0.5,
        autoDismissDelay: TimeInterval? = 3.0,
        onShow: (() -> Void)? = nil,
        onDismiss: (() -> Void)? = nil
    ) -> some View {
        TchatTooltip(
            text: text,
            position: position,
            trigger: trigger,
            style: style,
            maxWidth: maxWidth,
            showArrow: showArrow,
            delay: delay,
            autoDismissDelay: autoDismissDelay,
            onShow: onShow,
            onDismiss: onDismiss
        ) {
            self
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatTooltip_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.xl) {
                // Dark tooltip (default)
                TchatTooltip(
                    text: "This is a dark tooltip with helpful information",
                    position: .top,
                    trigger: .tap,
                    style: .dark
                ) {
                    Button("Dark Tooltip") { }
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                }

                // Light tooltip
                TchatTooltip(
                    text: "Light tooltip with border and light background",
                    position: .bottom,
                    trigger: .tap,
                    style: .light
                ) {
                    Button("Light Tooltip") { }
                        .padding()
                        .background(Color.gray.opacity(0.2))
                        .cornerRadius(8)
                }

                // Info tooltip
                TchatTooltip(
                    text: "This is an informational tooltip",
                    position: .trailing,
                    trigger: .tap,
                    style: .info
                ) {
                    Image(systemName: "info.circle")
                        .font(.title2)
                        .foregroundColor(.blue)
                }

                // Warning tooltip
                TchatTooltip(
                    text: "Warning: This action cannot be undone",
                    position: .leading,
                    trigger: .tap,
                    style: .warning
                ) {
                    Image(systemName: "exclamationmark.triangle")
                        .font(.title2)
                        .foregroundColor(.orange)
                }

                // Error tooltip
                TchatTooltip(
                    text: "Error: Something went wrong",
                    position: .auto,
                    trigger: .tap,
                    style: .error
                ) {
                    Image(systemName: "xmark.circle")
                        .font(.title2)
                        .foregroundColor(.red)
                }

                // Long press tooltip
                Button("Long Press Me") { }
                    .padding()
                    .background(Color.green)
                    .foregroundColor(.white)
                    .cornerRadius(8)
                    .tchatTooltip(
                        "This tooltip appears on long press",
                        trigger: .longPress,
                        delay: 0.5
                    )

                // Manual tooltip
                TchatTooltip(
                    text: "This tooltip can be controlled manually",
                    trigger: .manual
                ) {
                    Button("Manual Tooltip") { }
                        .padding()
                        .background(Color.purple)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                }
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif