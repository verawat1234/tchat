//
//  TchatCard.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Card component following Tchat design system
public struct TchatCard<Content: View>: View {

    // MARK: - Card Variants
    public enum Variant {
        case elevated
        case outlined
        case filled
        case glass
    }

    public enum Size {
        case compact
        case standard
        case expanded
    }

    // MARK: - Properties
    let content: Content
    let variant: Variant
    let size: Size
    let isInteractive: Bool
    let onTap: (() -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    @State private var isPressed = false

    // MARK: - Initializer
    public init(
        variant: Variant = .elevated,
        size: Size = .standard,
        isInteractive: Bool = false,
        onTap: (() -> Void)? = nil,
        @ViewBuilder content: () -> Content
    ) {
        self.content = content()
        self.variant = variant
        self.size = size
        self.isInteractive = isInteractive
        self.onTap = onTap
    }

    // MARK: - Body
    public var body: some View {
        Button(action: {
            onTap?()
        }) {
            content
                .padding(contentPadding)
                .frame(maxWidth: .infinity, alignment: .leading)
                .background(backgroundColor)
                .overlay(borderOverlay)
                .cornerRadius(cornerRadius)
                .shadow(
                    color: shadowColor,
                    radius: shadowRadius,
                    x: shadowOffset.x,
                    y: shadowOffset.y
                )
        }
        .buttonStyle(TchatCardButtonStyle(
            isInteractive: isInteractive,
            variant: variant
        ))
        .disabled(!isInteractive)
    }

    // MARK: - Computed Properties

    private var contentPadding: EdgeInsets {
        switch size {
        case .compact:
            return EdgeInsets(
                top: Spacing.sm,
                leading: Spacing.sm,
                bottom: Spacing.sm,
                trailing: Spacing.sm
            )
        case .standard:
            return EdgeInsets(
                top: Spacing.md,
                leading: Spacing.md,
                bottom: Spacing.md,
                trailing: Spacing.md
            )
        case .expanded:
            return EdgeInsets(
                top: Spacing.lg,
                leading: Spacing.lg,
                bottom: Spacing.lg,
                trailing: Spacing.lg
            )
        }
    }

    private var backgroundColor: Color {
        switch variant {
        case .elevated, .outlined:
            return colors.cardBackground
        case .filled:
            return colors.surface
        case .glass:
            return colors.cardBackground.opacity(0.8)
        }
    }

    @ViewBuilder
    private var borderOverlay: some View {
        switch variant {
        case .outlined:
            RoundedRectangle(cornerRadius: cornerRadius)
                .stroke(colors.border, lineWidth: 1)
        default:
            EmptyView()
        }
    }

    private var cornerRadius: CGFloat {
        switch size {
        case .compact:
            return Spacing.sm
        case .standard:
            return Spacing.md
        case .expanded:
            return Spacing.md
        }
    }

    private var shadowColor: Color {
        switch variant {
        case .elevated:
            return colors.shadowMedium
        case .glass:
            return colors.shadowLight
        default:
            return Color.clear
        }
    }

    private var shadowRadius: CGFloat {
        switch variant {
        case .elevated:
            return 8
        case .glass:
            return 4
        default:
            return 0
        }
    }

    private var shadowOffset: CGSize {
        switch variant {
        case .elevated:
            return CGSize(width: 0, height: 2)
        case .glass:
            return CGSize(width: 0, height: 1)
        default:
            return CGSize.zero
        }
    }
}

// MARK: - Card Button Style
struct TchatCardButtonStyle: ButtonStyle {
    let isInteractive: Bool
    let variant: TchatCard<AnyView>.Variant
    private let colors = Colors()

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(
                isInteractive && configuration.isPressed ? 0.98 : 1.0
            )
            .brightness(
                isInteractive && configuration.isPressed ? (
                    variant == .filled ? 0.05 : -0.02
                ) : 0
            )
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

// MARK: - Card Header Component
public struct TchatCardHeader: View {
    let title: String
    let subtitle: String?
    let leadingIcon: String?
    let trailingContent: AnyView?

    private let colors = Colors()

    public init(
        title: String,
        subtitle: String? = nil,
        leadingIcon: String? = nil,
        trailingContent: (() -> AnyView)? = nil
    ) {
        self.title = title
        self.subtitle = subtitle
        self.leadingIcon = leadingIcon
        self.trailingContent = trailingContent?()
    }

    public var body: some View {
        HStack(spacing: Spacing.sm) {
            if let leadingIcon = leadingIcon {
                Image(systemName: leadingIcon)
                    .foregroundColor(colors.primary)
                    .frame(width: 20, height: 20)
            }

            VStack(alignment: .leading, spacing: Spacing.xs) {
                Text(title)
                    .font(.headline)
                    .foregroundColor(colors.textPrimary)

                if let subtitle = subtitle {
                    Text(subtitle)
                        .font(.caption)
                        .foregroundColor(colors.textSecondary)
                }
            }

            Spacer()

            trailingContent
        }
    }
}

// MARK: - Card Footer Component
public struct TchatCardFooter: View {
    let content: AnyView

    public init(@ViewBuilder content: () -> AnyView) {
        self.content = content()
    }

    public var body: some View {
        HStack {
            content
        }
        .padding(.top, Spacing.sm)
    }
}

// MARK: - Convenience Initializers
extension TchatCard {
    /// Card with header and content
    public init(
        title: String,
        subtitle: String? = nil,
        leadingIcon: String? = nil,
        variant: Variant = .elevated,
        size: Size = .standard,
        isInteractive: Bool = false,
        onTap: (() -> Void)? = nil,
        @ViewBuilder content: () -> Content
    ) where Content == AnyView {
        self.init(
            variant: variant,
            size: size,
            isInteractive: isInteractive,
            onTap: onTap
        ) {
            AnyView(
                VStack(alignment: .leading, spacing: Spacing.sm) {
                    TchatCardHeader(
                        title: title,
                        subtitle: subtitle,
                        leadingIcon: leadingIcon
                    )

                    content()
                }
            )
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatCard_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.md) {
                // Basic elevated card
                TchatCard(variant: .elevated) {
                    VStack(alignment: .leading, spacing: Spacing.sm) {
                        Text("Elevated Card")
                            .font(.headline)
                        Text("This is an elevated card with shadow")
                            .font(.body)
                            .foregroundColor(Colors().textSecondary)
                    }
                }

                // Outlined card with header
                TchatCard(
                    title: "Card with Header",
                    subtitle: "Subtitle text",
                    leadingIcon: "star.fill",
                    variant: .outlined,
                    isInteractive: true,
                    onTap: { print("Card tapped") }
                ) {
                    VStack(alignment: .leading, spacing: Spacing.xs) {
                        Text("This card has a header and is interactive")
                        Text("Tap to interact")
                            .font(.caption)
                            .foregroundColor(Colors().primary)
                    }
                }

                // Filled card
                TchatCard(variant: .filled, size: .compact) {
                    HStack {
                        Image(systemName: "message.fill")
                            .foregroundColor(Colors().primary)
                        Text("Compact filled card")
                            .font(.subheadline)
                    }
                }

                // Glass card
                TchatCard(variant: .glass, size: .expanded) {
                    VStack(alignment: .leading, spacing: Spacing.md) {
                        TchatCardHeader(
                            title: "Glass Card",
                            subtitle: "Expanded size",
                            leadingIcon: "sparkles"
                        )

                        Text("This is a glass-style card with expanded padding")

                        TchatCardFooter {
                            AnyView(
                                HStack {
                                    TchatButton("Action", variant: .primary, size: .small) {
                                        print("Action tapped")
                                    }
                                    Spacer()
                                    Text("Footer content")
                                        .font(.caption)
                                        .foregroundColor(Colors().textTertiary)
                                }
                            )
                        }
                    }
                }

                // Interactive card with complex content
                TchatCard(
                    variant: .elevated,
                    isInteractive: true,
                    onTap: { print("Complex card tapped") }
                ) {
                    VStack(alignment: .leading, spacing: Spacing.md) {
                        HStack {
                            VStack(alignment: .leading) {
                                Text("Interactive Card")
                                    .font(.title3)
                                    .fontWeight(.semibold)
                                Text("With complex content")
                                    .font(.caption)
                                    .foregroundColor(Colors().textSecondary)
                            }
                            Spacer()
                            Image(systemName: "chevron.right")
                                .foregroundColor(Colors().textTertiary)
                        }

                        HStack {
                            Label("Feature 1", systemImage: "checkmark.circle.fill")
                                .foregroundColor(Colors().success)
                            Spacer()
                            Label("Feature 2", systemImage: "star.fill")
                                .foregroundColor(Colors().warning)
                        }
                        .font(.caption)
                    }
                }
            }
            .padding()
        }
        .background(Color(.systemGroupedBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif