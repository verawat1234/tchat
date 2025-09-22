//
//  TchatBreadcrumb.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Breadcrumb navigation component following Tchat design system
public struct TchatBreadcrumb: View {

    // MARK: - Breadcrumb Types
    public enum BreadcrumbStyle {
        case standard
        case compact
        case minimal
    }

    public enum BreadcrumbSize {
        case small
        case medium
        case large
    }

    // MARK: - Breadcrumb Item
    public struct BreadcrumbItem {
        let id: String
        let title: String
        let icon: String?
        let isClickable: Bool
        let action: (() -> Void)?

        public init(
            id: String,
            title: String,
            icon: String? = nil,
            isClickable: Bool = true,
            action: (() -> Void)? = nil
        ) {
            self.id = id
            self.title = title
            self.icon = icon
            self.isClickable = isClickable
            self.action = action
        }
    }

    // MARK: - Properties
    let items: [BreadcrumbItem]
    let style: BreadcrumbStyle
    let size: BreadcrumbSize
    let separator: String
    let maxItems: Int?
    let showHome: Bool
    let homeIcon: String
    let onItemTap: ((BreadcrumbItem) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var fontSize: CGFloat {
        switch size {
        case .small: return 12
        case .medium: return 14
        case .large: return 16
        }
    }

    private var iconSize: CGFloat {
        switch size {
        case .small: return 12
        case .medium: return 14
        case .large: return 16
        }
    }

    private var spacing: CGFloat {
        switch size {
        case .small: return 4
        case .medium: return 6
        case .large: return 8
        }
    }

    private var displayItems: [BreadcrumbItem] {
        guard let maxItems = maxItems, items.count > maxItems else {
            return items
        }

        if maxItems <= 2 {
            return Array(items.suffix(maxItems))
        }

        // Show first item, ellipsis, and last items
        let firstItem = items.first!
        let lastItems = Array(items.suffix(maxItems - 2))

        var result = [firstItem]

        // Add ellipsis item
        let ellipsisItem = BreadcrumbItem(
            id: "ellipsis",
            title: "...",
            isClickable: false
        )
        result.append(ellipsisItem)

        result.append(contentsOf: lastItems)
        return result
    }

    // MARK: - Initializer
    public init(
        items: [BreadcrumbItem],
        style: BreadcrumbStyle = .standard,
        size: BreadcrumbSize = .medium,
        separator: String = "chevron.right",
        maxItems: Int? = nil,
        showHome: Bool = false,
        homeIcon: String = "house",
        onItemTap: ((BreadcrumbItem) -> Void)? = nil
    ) {
        self.items = items
        self.style = style
        self.size = size
        self.separator = separator
        self.maxItems = maxItems
        self.showHome = showHome
        self.homeIcon = homeIcon
        self.onItemTap = onItemTap
    }

    // MARK: - Body
    public var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: spacing) {
                // Home icon (if enabled)
                if showHome {
                    homeButton

                    if !items.isEmpty {
                        separatorView
                    }
                }

                // Breadcrumb items
                ForEach(Array(displayItems.enumerated()), id: \.element.id) { index, item in
                    breadcrumbItemView(item, isLast: index == displayItems.count - 1)

                    if index < displayItems.count - 1 {
                        separatorView
                    }
                }
            }
            .padding(.horizontal, Spacing.sm)
        }
    }

    // MARK: - Home Button
    @ViewBuilder
    private var homeButton: some View {
        Button(action: {
            // Handle home tap
            let homeItem = BreadcrumbItem(
                id: "home",
                title: "Home",
                icon: homeIcon,
                isClickable: true
            )
            onItemTap?(homeItem)

            // Haptic feedback
            let impactFeedback = UIImpactFeedbackGenerator(style: .light)
            impactFeedback.impactOccurred()
        }) {
            Image(systemName: homeIcon)
                .font(.system(size: iconSize))
                .foregroundColor(colors.textSecondary)
        }
    }

    // MARK: - Breadcrumb Item View
    @ViewBuilder
    private func breadcrumbItemView(_ item: BreadcrumbItem, isLast: Bool) -> some View {
        Group {
            if item.isClickable && !isLast {
                Button(action: {
                    item.action?()
                    onItemTap?(item)

                    // Haptic feedback
                    let impactFeedback = UIImpactFeedbackGenerator(style: .light)
                    impactFeedback.impactOccurred()
                }) {
                    itemContent(item, isLast: isLast)
                }
                .buttonStyle(PlainButtonStyle())
            } else {
                itemContent(item, isLast: isLast)
            }
        }
    }

    // MARK: - Item Content
    @ViewBuilder
    private func itemContent(_ item: BreadcrumbItem, isLast: Bool) -> some View {
        HStack(spacing: spacing / 2) {
            // Icon
            if let icon = item.icon {
                Image(systemName: icon)
                    .font(.system(size: iconSize))
                    .foregroundColor(itemColor(item, isLast: isLast))
            }

            // Title
            Text(item.title)
                .font(.system(size: fontSize, weight: isLast ? .medium : .regular))
                .foregroundColor(itemColor(item, isLast: isLast))
                .lineLimit(1)
        }
        .padding(.horizontal, style == .standard ? spacing : 0)
        .padding(.vertical, style == .standard ? spacing / 2 : 0)
        .background(
            itemBackground(item, isLast: isLast)
        )
        .cornerRadius(style == .standard ? 4 : 0)
    }

    // MARK: - Separator View
    @ViewBuilder
    private var separatorView: some View {
        Image(systemName: separator)
            .font(.system(size: iconSize - 2))
            .foregroundColor(colors.textTertiary)
    }

    // MARK: - Styling Methods
    private func itemColor(_ item: BreadcrumbItem, isLast: Bool) -> Color {
        switch style {
        case .standard:
            if isLast {
                return colors.textPrimary
            } else if item.isClickable {
                return colors.primary
            } else {
                return colors.textSecondary
            }
        case .compact, .minimal:
            if isLast {
                return colors.textPrimary
            } else if item.isClickable {
                return colors.textSecondary
            } else {
                return colors.textTertiary
            }
        }
    }

    private func itemBackground(_ item: BreadcrumbItem, isLast: Bool) -> Color {
        guard style == .standard else { return Color.clear }

        if isLast {
            return colors.surface
        } else if item.isClickable {
            return Color.clear
        } else {
            return Color.clear
        }
    }
}

// MARK: - Convenience Initializers
extension TchatBreadcrumb {
    public static func fromPath(
        _ path: String,
        separator: String = "/",
        style: BreadcrumbStyle = .standard,
        size: BreadcrumbSize = .medium,
        onPathTap: ((String) -> Void)? = nil
    ) -> TchatBreadcrumb {
        let components = path.split(separator: Character(separator)).map(String.init)

        let items = components.enumerated().map { index, component in
            let path = components.prefix(index + 1).joined(separator: separator)

            return BreadcrumbItem(
                id: path,
                title: component,
                isClickable: true
            ) {
                onPathTap?(path)
            }
        }

        return TchatBreadcrumb(
            items: items,
            style: style,
            size: size
        )
    }
}

// MARK: - NavigationPath Support
@available(iOS 16.0, *)
extension TchatBreadcrumb {
    public init(
        navigationPath: Binding<NavigationPath>,
        pathComponents: [String],
        style: BreadcrumbStyle = .standard,
        size: BreadcrumbSize = .medium,
        separator: String = "chevron.right",
        maxItems: Int? = nil,
        showHome: Bool = true,
        homeIcon: String = "house"
    ) {
        let items = pathComponents.enumerated().map { index, component in
            BreadcrumbItem(
                id: "\(index)",
                title: component,
                isClickable: index < pathComponents.count - 1
            ) {
                // Navigate back to this level
                var newPath = NavigationPath()
                for i in 0...index {
                    newPath.append(pathComponents[i])
                }
                navigationPath.wrappedValue = newPath
            }
        }

        self.init(
            items: items,
            style: style,
            size: size,
            separator: separator,
            maxItems: maxItems,
            showHome: showHome,
            homeIcon: homeIcon
        ) { item in
            // Handle navigation
            if item.id == "home" {
                navigationPath.wrappedValue = NavigationPath()
            }
        }
    }
}

// MARK: - Preview
#if DEBUG
struct TchatBreadcrumb_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Default breadcrumb
                TchatBreadcrumb(
                    items: [
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "dashboard",
                            title: "Dashboard",
                            icon: "chart.bar"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "projects",
                            title: "Projects",
                            icon: "folder"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "mobile-app",
                            title: "Mobile App",
                            isClickable: false
                        )
                    ],
                    style: .standard,
                    showHome: true
                )

                Divider()

                // Compact breadcrumb
                TchatBreadcrumb(
                    items: [
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "users",
                            title: "Users"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "profile",
                            title: "Profile"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "settings",
                            title: "Settings",
                            isClickable: false
                        )
                    ],
                    style: .compact,
                    size: .small
                )

                Divider()

                // Minimal breadcrumb with truncation
                TchatBreadcrumb(
                    items: [
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "level1",
                            title: "Level 1"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "level2",
                            title: "Level 2"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "level3",
                            title: "Level 3"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "level4",
                            title: "Level 4"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "level5",
                            title: "Current Page",
                            isClickable: false
                        )
                    ],
                    style: .minimal,
                    maxItems: 4
                )

                Divider()

                // Path-based breadcrumb
                TchatBreadcrumb.fromPath(
                    "/projects/mobile-app/src/components",
                    style: .standard,
                    size: .medium
                )

                Divider()

                // Large breadcrumb with custom separator
                TchatBreadcrumb(
                    items: [
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "docs",
                            title: "Documentation",
                            icon: "doc.text"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "api",
                            title: "API Reference",
                            icon: "link"
                        ),
                        TchatBreadcrumb.BreadcrumbItem(
                            id: "endpoints",
                            title: "Endpoints",
                            isClickable: false
                        )
                    ],
                    style: .standard,
                    size: .large,
                    separator: "arrow.right",
                    showHome: true,
                    homeIcon: "house.fill"
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif