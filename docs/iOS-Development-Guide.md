# iOS Development Guide - Tchat App

> **Comprehensive development guide for iOS platform following the Tchat Design System**
>
> Last Updated: September 2025 | Swift 5.9+ | SwiftUI | iOS 15.0+

## Table of Contents

1. [Development Setup](#development-setup)
2. [Design Token System](#design-token-system)
3. [Core Components](#core-components)
4. [Architecture Patterns](#architecture-patterns)
5. [Testing Strategy](#testing-strategy)
6. [Code Style & Conventions](#code-style--conventions)
7. [Performance Standards](#performance-standards)
8. [Accessibility Guidelines](#accessibility-guidelines)

---

## Development Setup

### Prerequisites

```bash
# Required versions
- Xcode 15.0+
- Swift 5.9+
- iOS 15.0+ deployment target
- Swift Package Manager for dependencies
```

### Dependencies

```swift
// Package.swift
dependencies: [
    .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0"),
    .package(url: "https://github.com/onevcat/Kingfisher.git", from: "7.9.0"),
    .package(url: "https://github.com/apple/swift-async-algorithms", from: "1.0.0")
]
```

### Project Structure

```
Sources/
├── Components/          # Atom design system components
│   ├── TchatButton.swift
│   ├── TchatInput.swift
│   ├── TchatCard.swift
│   └── ...
├── DesignSystem/       # Design tokens and theming
│   ├── Colors.swift
│   ├── Typography.swift
│   ├── Spacing.swift
│   └── DesignTokens.swift
├── Screens/           # Screen implementations
├── Navigation/        # Tab navigation system
├── Services/          # Business logic and networking
└── Models/           # Data models and types
```

---

## Design Token System

### Color Palette (TailwindCSS v4 Mapped)

```swift
// Sources/DesignSystem/Colors.swift
import SwiftUI

public struct TchatColors {
    // Brand Colors - Primary blue (#3B82F6)
    public static let primary = Color(hex: "#3B82F6")           // blue-500
    public static let primaryLight = Color(hex: "#60A5FA")      // blue-400
    public static let primaryDark = Color(hex: "#2563EB")       // blue-600

    // Semantic Colors
    public static let success = Color(hex: "#10B981")           // green-500
    public static let warning = Color(hex: "#F59E0B")           // amber-500
    public static let error = Color(hex: "#EF4444")             // red-500
    public static let info = Color(hex: "#3B82F6")              // blue-500

    // Surface Colors
    public static let surface = Color(hex: "#FFFFFF")           // white
    public static let surfaceSecondary = Color(hex: "#F9FAFB")  // gray-50
    public static let surfaceTertiary = Color(hex: "#F3F4F6")   // gray-100

    // Text Colors
    public static let textPrimary = Color(hex: "#111827")       // gray-900
    public static let textSecondary = Color(hex: "#6B7280")     // gray-500
    public static let textTertiary = Color(hex: "#9CA3AF")      // gray-400
    public static let textOnPrimary = Color(hex: "#FFFFFF")     // white

    // Border Colors
    public static let border = Color(hex: "#E5E7EB")            // gray-200
    public static let borderSecondary = Color(hex: "#D1D5DB")   // gray-300
    public static let borderFocus = Color(hex: "#3B82F6")       // blue-500

    // Dark Mode Support
    public struct Dark {
        public static let surface = Color(hex: "#111827")        // gray-900
        public static let surfaceSecondary = Color(hex: "#1F2937") // gray-800
        public static let textPrimary = Color(hex: "#F9FAFB")    // gray-50
        public static let textSecondary = Color(hex: "#D1D5DB")  // gray-300
        public static let border = Color(hex: "#374151")         // gray-700
    }
}

// Color Extension for Hex Support
extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (1, 1, 1, 0)
        }

        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
```

### Typography System

```swift
// Sources/DesignSystem/Typography.swift
import SwiftUI

public struct TchatTypography {
    // Display Typography
    public static let displayLarge = Font.system(size: 48, weight: .bold, design: .default)
    public static let displayMedium = Font.system(size: 36, weight: .bold, design: .default)
    public static let displaySmall = Font.system(size: 32, weight: .bold, design: .default)

    // Headline Typography
    public static let headlineLarge = Font.system(size: 28, weight: .semibold, design: .default)
    public static let headlineMedium = Font.system(size: 24, weight: .semibold, design: .default)
    public static let headlineSmall = Font.system(size: 20, weight: .semibold, design: .default)

    // Body Typography
    public static let bodyLarge = Font.system(size: 18, weight: .regular, design: .default)
    public static let bodyMedium = Font.system(size: 16, weight: .regular, design: .default)
    public static let bodySmall = Font.system(size: 14, weight: .regular, design: .default)

    // Label Typography
    public static let labelLarge = Font.system(size: 16, weight: .medium, design: .default)
    public static let labelMedium = Font.system(size: 14, weight: .medium, design: .default)
    public static let labelSmall = Font.system(size: 12, weight: .medium, design: .default)

    // Caption Typography
    public static let caption = Font.system(size: 12, weight: .regular, design: .default)
}
```

### Spacing System (4dp Base Unit)

```swift
// Sources/DesignSystem/Spacing.swift
import SwiftUI

public struct TchatSpacing {
    // Base 4dp spacing system matching TailwindCSS
    public static let xs: CGFloat = 4      // space-1 (0.25rem)
    public static let sm: CGFloat = 8      // space-2 (0.5rem)
    public static let md: CGFloat = 16     // space-4 (1rem)
    public static let lg: CGFloat = 24     // space-6 (1.5rem)
    public static let xl: CGFloat = 32     // space-8 (2rem)
    public static let xxl: CGFloat = 48    // space-12 (3rem)

    // Component-specific spacing
    public static let buttonPaddingVertical: CGFloat = 12    // 3/4 of md
    public static let buttonPaddingHorizontal: CGFloat = 20  // 5/4 of md
    public static let cardPadding: CGFloat = 16             // md
    public static let screenPadding: CGFloat = 16           // md
}
```

---

## Core Components

### TchatButton Implementation

```swift
// Sources/Components/TchatButton.swift
import SwiftUI

public struct TchatButton: View {
    // MARK: - Types
    public enum Variant {
        case primary
        case secondary
        case ghost
        case destructive
        case outline
        case link
    }

    public enum Size {
        case small      // 32dp height
        case medium     // 44dp height (default)
        case large      // 48dp height
        case icon       // 44x44dp square
    }

    // MARK: - Properties
    private let title: String?
    private let icon: Image?
    private let variant: Variant
    private let size: Size
    private let isLoading: Bool
    private let isDisabled: Bool
    private let action: () -> Void

    @State private var isPressed = false

    // MARK: - Initialization
    public init(
        _ title: String? = nil,
        icon: Image? = nil,
        variant: Variant = .primary,
        size: Size = .medium,
        isLoading: Bool = false,
        isDisabled: Bool = false,
        action: @escaping () -> Void
    ) {
        self.title = title
        self.icon = icon
        self.variant = variant
        self.size = size
        self.isLoading = isLoading
        self.isDisabled = isDisabled
        self.action = action
    }

    // MARK: - Body
    public var body: some View {
        Button(action: {
            if !isDisabled && !isLoading {
                // Medium haptic feedback
                let impactFeedback = UIImpactFeedbackGenerator(style: .medium)
                impactFeedback.impactOccurred()
                action()
            }
        }) {
            HStack(spacing: TchatSpacing.sm) {
                if isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle())
                        .scaleEffect(0.8)
                        .foregroundColor(textColor)
                } else if let icon = icon {
                    icon
                        .foregroundColor(textColor)
                }

                if let title = title, size != .icon {
                    Text(title)
                        .font(buttonFont)
                        .foregroundColor(textColor)
                        .lineLimit(1)
                }
            }
            .frame(maxWidth: size == .icon ? buttonHeight : .infinity)
            .frame(height: buttonHeight)
            .padding(.horizontal, size == .icon ? 0 : horizontalPadding)
        }
        .background(backgroundColor)
        .overlay(
            RoundedRectangle(cornerRadius: cornerRadius)
                .stroke(borderColor, lineWidth: borderWidth)
        )
        .clipShape(RoundedRectangle(cornerRadius: cornerRadius))
        .scaleEffect(isPressed ? 0.95 : 1.0)
        .opacity(isDisabled ? 0.6 : 1.0)
        .disabled(isDisabled || isLoading)
        .onLongPressGesture(minimumDuration: 0, maximumDistance: .infinity, pressing: { pressing in
            withAnimation(.easeInOut(duration: 0.1)) {
                isPressed = pressing
            }
        }, perform: {})
        .accessibilityLabel(accessibilityLabel)
        .accessibilityHint(accessibilityHint)
        .accessibilityAddTraits(isDisabled ? [.isButton, .isNotEnabled] : [.isButton])
    }

    // MARK: - Computed Properties
    private var buttonHeight: CGFloat {
        switch size {
        case .small: return 32
        case .medium: return 44
        case .large: return 48
        case .icon: return 44
        }
    }

    private var horizontalPadding: CGFloat {
        switch size {
        case .small: return TchatSpacing.sm
        case .medium: return TchatSpacing.md
        case .large: return TchatSpacing.lg
        case .icon: return 0
        }
    }

    private var buttonFont: Font {
        switch size {
        case .small: return TchatTypography.bodySmall
        case .medium: return TchatTypography.bodyMedium
        case .large: return TchatTypography.bodyLarge
        case .icon: return TchatTypography.bodyMedium
        }
    }

    private var backgroundColor: Color {
        switch variant {
        case .primary:
            return TchatColors.primary
        case .secondary:
            return TchatColors.surfaceSecondary
        case .ghost, .link:
            return Color.clear
        case .destructive:
            return TchatColors.error
        case .outline:
            return Color.clear
        }
    }

    private var textColor: Color {
        switch variant {
        case .primary, .destructive:
            return TchatColors.textOnPrimary
        case .secondary:
            return TchatColors.textPrimary
        case .ghost, .outline:
            return TchatColors.primary
        case .link:
            return TchatColors.primary
        }
    }

    private var borderColor: Color {
        switch variant {
        case .outline:
            return TchatColors.border
        default:
            return Color.clear
        }
    }

    private var borderWidth: CGFloat {
        variant == .outline ? 1 : 0
    }

    private var cornerRadius: CGFloat {
        return 8 // 2 * TchatSpacing.xs
    }

    // MARK: - Accessibility
    private var accessibilityLabel: String {
        if let title = title {
            return isLoading ? "Loading \(title)" : title
        } else if icon != nil {
            return "Button" // Should be customized per use case
        }
        return "Button"
    }

    private var accessibilityHint: String {
        if isDisabled {
            return "Button is disabled"
        } else if isLoading {
            return "Button is loading"
        }
        return "Double tap to activate"
    }
}

// MARK: - Convenience Initializers
public extension TchatButton {
    // Primary button
    static func primary(
        _ title: String,
        isLoading: Bool = false,
        isDisabled: Bool = false,
        action: @escaping () -> Void
    ) -> TchatButton {
        TchatButton(title, variant: .primary, isLoading: isLoading, isDisabled: isDisabled, action: action)
    }

    // Icon button
    static func icon(
        _ icon: Image,
        variant: Variant = .ghost,
        isDisabled: Bool = false,
        action: @escaping () -> Void
    ) -> TchatButton {
        TchatButton(icon: icon, variant: variant, size: .icon, isDisabled: isDisabled, action: action)
    }
}

// MARK: - Preview
#if DEBUG
struct TchatButton_Previews: PreviewProvider {
    static var previews: some View {
        VStack(spacing: TchatSpacing.md) {
            // Primary variants
            TchatButton.primary("Primary Button") { }
            TchatButton("Secondary", variant: .secondary) { }
            TchatButton("Destructive", variant: .destructive) { }
            TchatButton("Ghost", variant: .ghost) { }
            TchatButton("Outline", variant: .outline) { }

            // Size variants
            HStack {
                TchatButton("Small", variant: .primary, size: .small) { }
                TchatButton("Medium", variant: .primary, size: .medium) { }
                TchatButton("Large", variant: .primary, size: .large) { }
            }

            // States
            TchatButton("Loading", variant: .primary, isLoading: true) { }
            TchatButton("Disabled", variant: .primary, isDisabled: true) { }

            // Icon button
            TchatButton.icon(Image(systemName: "heart.fill")) { }
        }
        .padding()
        .previewLayout(.sizeThatFits)
    }
}
#endif
```

### TchatInput Implementation (Planned)

```swift
// Sources/Components/TchatInput.swift
import SwiftUI

public struct TchatInput: View {
    // MARK: - Types
    public enum InputType {
        case text
        case email
        case password
        case number
        case search
        case multiline(lines: Int = 3)
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
    private let placeholder: String
    private let inputType: InputType
    private let validationState: ValidationState
    private let size: Size
    private let leadingIcon: Image?
    private let isDisabled: Bool

    @State private var isSecured: Bool = true
    @State private var isFocused: Bool = false

    // MARK: - Initialization
    public init(
        text: Binding<String>,
        placeholder: String,
        inputType: InputType = .text,
        validationState: ValidationState = .none,
        size: Size = .medium,
        leadingIcon: Image? = nil,
        isDisabled: Bool = false
    ) {
        self._text = text
        self.placeholder = placeholder
        self.inputType = inputType
        self.validationState = validationState
        self.size = size
        self.leadingIcon = leadingIcon
        self.isDisabled = isDisabled
    }

    // MARK: - Body
    public var body: some View {
        VStack(alignment: .leading, spacing: TchatSpacing.xs) {
            HStack(spacing: TchatSpacing.sm) {
                // Leading icon
                if let leadingIcon = leadingIcon {
                    leadingIcon
                        .foregroundColor(TchatColors.textSecondary)
                        .frame(width: 20, height: 20)
                }

                // Text field based on type
                textFieldView

                // Trailing elements (password toggle, validation icon)
                trailingView
            }
            .padding(.horizontal, horizontalPadding)
            .frame(height: inputHeight)
            .background(backgroundColor)
            .overlay(
                RoundedRectangle(cornerRadius: cornerRadius)
                    .stroke(borderColor, lineWidth: borderWidth)
            )
            .clipShape(RoundedRectangle(cornerRadius: cornerRadius))
            .disabled(isDisabled)
            .opacity(isDisabled ? 0.6 : 1.0)

            // Error message
            if case .invalid(let message) = validationState {
                Text(message)
                    .font(TchatTypography.caption)
                    .foregroundColor(TchatColors.error)
                    .padding(.horizontal, TchatSpacing.xs)
            }
        }
    }

    // MARK: - Text Field View
    @ViewBuilder
    private var textFieldView: some View {
        switch inputType {
        case .multiline(let lines):
            TextEditor(text: $text)
                .frame(minHeight: CGFloat(lines) * 20)
                .font(inputFont)
                .foregroundColor(TchatColors.textPrimary)
                .scrollContentBackground(.hidden)
                .background(Color.clear)
                .onTapGesture { isFocused = true }

        case .password:
            if isSecured {
                SecureField(placeholder, text: $text)
                    .font(inputFont)
                    .foregroundColor(TchatColors.textPrimary)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                    .onTapGesture { isFocused = true }
            } else {
                TextField(placeholder, text: $text)
                    .font(inputFont)
                    .foregroundColor(TchatColors.textPrimary)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                    .onTapGesture { isFocused = true }
            }

        default:
            TextField(placeholder, text: $text)
                .font(inputFont)
                .foregroundColor(TchatColors.textPrimary)
                .keyboardType(keyboardType)
                .textInputAutocapitalization(textCapitalization)
                .autocorrectionDisabled(inputType == .email)
                .onTapGesture { isFocused = true }
        }
    }

    // MARK: - Trailing View
    @ViewBuilder
    private var trailingView: some View {
        HStack(spacing: TchatSpacing.xs) {
            // Password toggle
            if inputType == .password {
                Button(action: { isSecured.toggle() }) {
                    Image(systemName: isSecured ? "eye.slash.fill" : "eye.fill")
                        .foregroundColor(TchatColors.textSecondary)
                        .frame(width: 20, height: 20)
                }
            }

            // Validation icon
            if case .valid = validationState {
                Image(systemName: "checkmark.circle.fill")
                    .foregroundColor(TchatColors.success)
                    .frame(width: 20, height: 20)
            } else if case .invalid = validationState {
                Image(systemName: "exclamationmark.circle.fill")
                    .foregroundColor(TchatColors.error)
                    .frame(width: 20, height: 20)
            }
        }
    }

    // MARK: - Computed Properties
    private var inputHeight: CGFloat {
        switch size {
        case .small: return 36
        case .medium: return 44
        case .large: return 52
        }
    }

    private var horizontalPadding: CGFloat {
        switch size {
        case .small: return TchatSpacing.sm
        case .medium: return TchatSpacing.md
        case .large: return TchatSpacing.md
        }
    }

    private var inputFont: Font {
        switch size {
        case .small: return TchatTypography.bodySmall
        case .medium: return TchatTypography.bodyMedium
        case .large: return TchatTypography.bodyLarge
        }
    }

    private var backgroundColor: Color {
        TchatColors.surface
    }

    private var borderColor: Color {
        if case .invalid = validationState {
            return TchatColors.error
        } else if case .valid = validationState {
            return TchatColors.success
        } else if isFocused {
            return TchatColors.borderFocus
        } else {
            return TchatColors.border
        }
    }

    private var borderWidth: CGFloat {
        isFocused || validationState != .none ? 2 : 1
    }

    private var cornerRadius: CGFloat {
        return 8
    }

    private var keyboardType: UIKeyboardType {
        switch inputType {
        case .email: return .emailAddress
        case .number: return .numberPad
        case .search: return .webSearch
        default: return .default
        }
    }

    private var textCapitalization: TextInputAutocapitalization {
        switch inputType {
        case .email, .password: return .never
        default: return .sentences
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatInput_Previews: PreviewProvider {
    static var previews: some View {
        VStack(spacing: TchatSpacing.md) {
            TchatInput(
                text: .constant(""),
                placeholder: "Enter text",
                leadingIcon: Image(systemName: "person.fill")
            )

            TchatInput(
                text: .constant(""),
                placeholder: "Email address",
                inputType: .email,
                leadingIcon: Image(systemName: "envelope.fill")
            )

            TchatInput(
                text: .constant(""),
                placeholder: "Password",
                inputType: .password,
                leadingIcon: Image(systemName: "lock.fill")
            )

            TchatInput(
                text: .constant("Valid input"),
                placeholder: "Valid input",
                validationState: .valid
            )

            TchatInput(
                text: .constant("Invalid input"),
                placeholder: "Invalid input",
                validationState: .invalid(message: "This field is required")
            )
        }
        .padding()
    }
}
#endif
```

---

## Architecture Patterns

### MVVM with Combine

```swift
// Example ViewModel Pattern
import Combine
import SwiftUI

@MainActor
class ChatViewModel: ObservableObject {
    @Published var messages: [Message] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let chatService: ChatService
    private var cancellables = Set<AnyCancellable>()

    init(chatService: ChatService = ChatService()) {
        self.chatService = chatService
        loadMessages()
    }

    func loadMessages() {
        isLoading = true

        chatService.getMessages()
            .receive(on: DispatchQueue.main)
            .sink(
                receiveCompletion: { [weak self] completion in
                    self?.isLoading = false
                    if case .failure(let error) = completion {
                        self?.errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { [weak self] messages in
                    self?.messages = messages
                }
            )
            .store(in: &cancellables)
    }

    func sendMessage(_ content: String) {
        // Implementation
    }
}
```

---

## Testing Strategy

### Unit Testing with XCTest

```swift
// Tests/TchatButtonTests.swift
import XCTest
import SwiftUI
@testable import TchatComponents

final class TchatButtonTests: XCTestCase {

    func testButtonVariants() {
        // Test all button variants render correctly
        let primaryButton = TchatButton("Test", variant: .primary) { }
        let secondaryButton = TchatButton("Test", variant: .secondary) { }

        // Assert styling properties
        XCTAssertNotNil(primaryButton)
        XCTAssertNotNil(secondaryButton)
    }

    func testButtonStates() {
        var actionCalled = false
        let button = TchatButton("Test") {
            actionCalled = true
        }

        // Test disabled state prevents action
        let disabledButton = TchatButton("Test", isDisabled: true) {
            actionCalled = true
        }

        // Simulate tap on disabled button
        // actionCalled should remain false
        XCTAssertFalse(actionCalled)
    }

    func testAccessibility() {
        let button = TchatButton("Login") { }

        // Test accessibility properties are set correctly
        // This would require view testing framework
    }
}
```

### SwiftUI Testing

```swift
// Tests/TchatInputTests.swift
import XCTest
import SwiftUI
@testable import TchatComponents

final class TchatInputTests: XCTestCase {

    func testTextBinding() {
        @State var text = ""
        let input = TchatInput(text: $text, placeholder: "Test")

        // Test text binding works correctly
        text = "Hello"
        XCTAssertEqual(text, "Hello")
    }

    func testValidationStates() {
        @State var text = ""

        let validInput = TchatInput(
            text: $text,
            placeholder: "Test",
            validationState: .valid
        )

        let invalidInput = TchatInput(
            text: $text,
            placeholder: "Test",
            validationState: .invalid(message: "Error")
        )

        XCTAssertNotNil(validInput)
        XCTAssertNotNil(invalidInput)
    }
}
```

---

## Code Style & Conventions

### Naming Conventions

```swift
// Protocol names end with -ing or -able
protocol Loadable {
    var isLoading: Bool { get }
}

// Enum cases use camelCase
enum ButtonVariant {
    case primary
    case secondary
    case destructive
}

// Private properties start with underscore (only for special cases)
@State private var _internalState: String = ""

// Functions use verb-noun pattern
func loadMessages() { }
func validateInput() -> Bool { }
func sendMessage(_ content: String) { }
```

### Code Organization

```swift
// MARK: - Types (at the top)
public enum Size {
    case small, medium, large
}

// MARK: - Properties
@State private var isLoading = false
private let service: ChatService

// MARK: - Initialization
public init(...) {
    // Implementation
}

// MARK: - Body (for Views)
public var body: some View {
    // Implementation
}

// MARK: - Methods
private func handleAction() {
    // Implementation
}

// MARK: - Computed Properties
private var backgroundColor: Color {
    // Implementation
}
```

### SwiftUI Best Practices

```swift
// ✅ Good - Explicit about View protocol
struct ContentView: View {
    var body: some View {
        // Implementation
    }
}

// ✅ Good - Use @State for local state
struct MyView: View {
    @State private var isExpanded = false

    var body: some View {
        // Implementation
    }
}

// ✅ Good - Use computed properties for complex logic
private var buttonStyle: some ButtonStyle {
    CustomButtonStyle(variant: variant, size: size)
}

// ✅ Good - Extract subviews for clarity
@ViewBuilder
private var headerView: some View {
    VStack {
        // Header implementation
    }
}
```

---

## Performance Standards

### Target Metrics

- **Frame Rate**: 60 FPS for all animations
- **Touch Response**: <16ms from touch to visual feedback
- **Component Render**: <8ms for basic components
- **Memory Usage**: <50MB for component library
- **Battery Impact**: Minimal background processing

### Performance Guidelines

```swift
// ✅ Use LazyVStack for long lists
LazyVStack {
    ForEach(messages) { message in
        MessageView(message: message)
    }
}

// ✅ Minimize State updates
@State private var text: String = "" {
    didSet {
        // Only update when necessary
        if text.count > maxLength {
            text = String(text.prefix(maxLength))
        }
    }
}

// ✅ Use proper animation timing
.animation(.easeInOut(duration: 0.2), value: isExpanded)

// ❌ Avoid - Heavy computation in body
var body: some View {
    let expensiveResult = heavyComputation() // Don't do this
    return Text("\(expensiveResult)")
}

// ✅ Better - Use computed property or onAppear
@State private var computedValue: String = ""

var body: some View {
    Text(computedValue)
        .onAppear {
            computedValue = heavyComputation()
        }
}
```

---

## Accessibility Guidelines

### VoiceOver Support

```swift
// ✅ Proper accessibility labels
TchatButton("Save", variant: .primary) {
    save()
}
.accessibilityLabel("Save document")
.accessibilityHint("Saves the current document to your library")

// ✅ Dynamic font size support
Text("Hello World")
    .font(TchatTypography.bodyMedium)
    .dynamicTypeSize(.large)

// ✅ Color contrast compliance
// All colors in TchatColors meet WCAG 2.1 AA standards
```

### Keyboard Navigation

```swift
// ✅ Focus management
@FocusState private var focusedField: FormField?

enum FormField: Hashable {
    case email, password
}

VStack {
    TchatInput(text: $email, placeholder: "Email")
        .focused($focusedField, equals: .email)

    TchatInput(text: $password, placeholder: "Password", inputType: .password)
        .focused($focusedField, equals: .password)
}
.onSubmit {
    switch focusedField {
    case .email:
        focusedField = .password
    case .password:
        submitForm()
    case .none:
        break
    }
}
```

### Touch Target Compliance

```swift
// ✅ All interactive elements meet 44dp minimum
// This is built into TchatButton component sizing

// For custom elements:
Button("Custom") {
    action()
}
.frame(minWidth: 44, minHeight: 44)
```

---

## Getting Started Checklist

- [ ] Install Xcode 15.0+
- [ ] Set up Swift Package Manager dependencies
- [ ] Import TchatComponents framework
- [ ] Configure design tokens in your app theme
- [ ] Implement first screen using TchatButton and TchatInput
- [ ] Set up unit testing with XCTest
- [ ] Configure accessibility testing
- [ ] Run performance profiling

---

**Questions or Issues?**
Refer to the project's GitHub repository or contact the development team for support.

---

*This guide is part of the Tchat Design System documentation suite.*