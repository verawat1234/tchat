//
//  TchatButton.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Primary button component following Tchat design system
public struct TchatButton: View {

    // MARK: - Button Variants
    public enum Variant {
        case primary
        case secondary
        case ghost
        case destructive
        case outline
    }

    public enum Size {
        case small
        case medium
        case large
    }

    // MARK: - Properties
    let title: String
    let variant: Variant
    let size: Size
    let isDisabled: Bool
    let isLoading: Bool
    let action: () -> Void

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Initializer
    public init(
        _ title: String,
        variant: Variant = .primary,
        size: Size = .medium,
        isDisabled: Bool = false,
        isLoading: Bool = false,
        action: @escaping () -> Void
    ) {
        self.title = title
        self.variant = variant
        self.size = size
        self.isDisabled = isDisabled
        self.isLoading = isLoading
        self.action = action
    }

    // MARK: - Body
    public var body: some View {
        Button(action: isDisabled || isLoading ? {} : action) {
            HStack(spacing: Spacing.xs) {
                if isLoading {
                    ProgressView()
                        .scaleEffect(0.8)
                        .foregroundColor(textColor)
                }

                if !isLoading || !title.isEmpty {
                    Text(title)
                        .font(buttonFont)
                        .foregroundColor(textColor)
                }
            }
            .frame(minWidth: buttonMinWidth, minHeight: buttonHeight)
            .padding(.horizontal, horizontalPadding)
            .background(backgroundColor)
            .cornerRadius(Spacing.sm)
            .overlay(
                RoundedRectangle(cornerRadius: Spacing.sm)
                    .stroke(borderColor, lineWidth: borderWidth)
            )
            .opacity(isDisabled ? 0.6 : 1.0)
        }
        .disabled(isDisabled || isLoading)
        .buttonStyle(TchatButtonStyle(variant: variant))
    }

    // MARK: - Computed Properties

    private var buttonFont: Font {
        switch size {
        case .small:
            return .system(size: 14, weight: .medium)
        case .medium:
            return .system(size: 16, weight: .medium)
        case .large:
            return .system(size: 18, weight: .medium)
        }
    }

    private var buttonHeight: CGFloat {
        switch size {
        case .small: return 32
        case .medium: return 40
        case .large: return 48
        }
    }

    private var buttonMinWidth: CGFloat {
        switch size {
        case .small: return 60
        case .medium: return 80
        case .large: return 100
        }
    }

    private var horizontalPadding: CGFloat {
        switch size {
        case .small: return Spacing.sm
        case .medium: return Spacing.md
        case .large: return Spacing.lg
        }
    }

    private var backgroundColor: Color {
        switch variant {
        case .primary:
            return colors.primary
        case .secondary:
            return colors.surface
        case .ghost:
            return Color.clear
        case .destructive:
            return colors.error
        case .outline:
            return Color.clear
        }
    }

    private var textColor: Color {
        switch variant {
        case .primary:
            return colors.textOnPrimary
        case .secondary:
            return colors.textPrimary
        case .ghost:
            return colors.primary
        case .destructive:
            return colors.textOnPrimary
        case .outline:
            return colors.primary
        }
    }

    private var borderColor: Color {
        switch variant {
        case .primary, .secondary, .destructive:
            return Color.clear
        case .ghost:
            return Color.clear
        case .outline:
            return colors.border
        }
    }

    private var borderWidth: CGFloat {
        switch variant {
        case .outline:
            return 1
        default:
            return 0
        }
    }
}

// MARK: - Button Style
struct TchatButtonStyle: ButtonStyle {
    let variant: TchatButton.Variant
    private let colors = Colors()

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.95 : 1.0)
            .brightness(configuration.isPressed ? (variant == .primary || variant == .destructive ? -0.1 : 0.05) : 0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

// MARK: - Preview
#if DEBUG
struct TchatButton_Previews: PreviewProvider {
    static var previews: some View {
        VStack(spacing: Spacing.md) {
            TchatButton("Primary Button", variant: .primary) {
                print("Primary tapped")
            }

            TchatButton("Secondary Button", variant: .secondary) {
                print("Secondary tapped")
            }

            TchatButton("Ghost Button", variant: .ghost) {
                print("Ghost tapped")
            }

            TchatButton("Outline Button", variant: .outline) {
                print("Outline tapped")
            }

            TchatButton("Destructive Button", variant: .destructive) {
                print("Destructive tapped")
            }

            TchatButton("Loading...", variant: .primary, isLoading: true) {
                print("Loading tapped")
            }

            TchatButton("Disabled Button", variant: .primary, isDisabled: true) {
                print("Disabled tapped")
            }

            HStack(spacing: Spacing.sm) {
                TchatButton("Small", variant: .primary, size: .small) {}
                TchatButton("Medium", variant: .primary, size: .medium) {}
                TchatButton("Large", variant: .primary, size: .large) {}
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif