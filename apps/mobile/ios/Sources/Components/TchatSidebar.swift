//
//  TchatSidebar.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Sidebar navigation component following Tchat design system
public struct TchatSidebar: View {

    // MARK: - Sidebar Types
    public enum SidebarMode {
        case overlay
        case push
        case permanent
    }

    public enum SidebarPosition {
        case leading
        case trailing
    }

    public enum SidebarSize {
        case compact
        case standard
        case wide
    }

    // MARK: - Sidebar Item
    public struct SidebarItem {
        let id: String
        let title: String
        let icon: String?
        let badge: String?
        let isDisabled: Bool
        let children: [SidebarItem]
        let action: (() -> Void)?

        public init(
            id: String,
            title: String,
            icon: String? = nil,
            badge: String? = nil,
            isDisabled: Bool = false,
            children: [SidebarItem] = [],
            action: (() -> Void)? = nil
        ) {
            self.id = id
            self.title = title
            self.icon = icon
            self.badge = badge
            self.isDisabled = isDisabled
            self.children = children
            self.action = action
        }
    }

    // MARK: - Properties
    @Binding private var isOpen: Bool
    @Binding private var selectedItem: String?
    @State private var expandedSections: Set<String> = []
    @State private var dragOffset: CGFloat = 0

    let items: [SidebarItem]
    let mode: SidebarMode
    let position: SidebarPosition
    let size: SidebarSize
    let showOverlay: Bool
    let allowDismiss: Bool
    let header: AnyView?
    let footer: AnyView?
    let onSelectionChange: ((String?) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var sidebarWidth: CGFloat {
        switch size {
        case .compact: return 240
        case .standard: return 280
        case .wide: return 320
        }
    }

    private var overlayOpacity: Double {
        guard mode == .overlay else { return 0 }
        return isOpen ? 0.5 : 0
    }

    private var sidebarOffset: CGFloat {
        let baseOffset: CGFloat

        switch mode {
        case .permanent:
            baseOffset = 0
        case .overlay, .push:
            if isOpen {
                baseOffset = 0
            } else {
                baseOffset = position == .leading ? -sidebarWidth : sidebarWidth
            }
        }

        return baseOffset + dragOffset
    }

    // MARK: - Initializer
    public init(
        isOpen: Binding<Bool>,
        selectedItem: Binding<String?> = .constant(nil),
        items: [SidebarItem],
        mode: SidebarMode = .overlay,
        position: SidebarPosition = .leading,
        size: SidebarSize = .standard,
        showOverlay: Bool = true,
        allowDismiss: Bool = true,
        header: AnyView? = nil,
        footer: AnyView? = nil,
        onSelectionChange: ((String?) -> Void)? = nil
    ) {
        self._isOpen = isOpen
        self._selectedItem = selectedItem
        self.items = items
        self.mode = mode
        self.position = position
        self.size = size
        self.showOverlay = showOverlay
        self.allowDismiss = allowDismiss
        self.header = header
        self.footer = footer
        self.onSelectionChange = onSelectionChange
    }

    // MARK: - Body
    public var body: some View {
        ZStack {
            // Overlay
            if mode == .overlay && showOverlay {
                Color.black
                    .opacity(overlayOpacity)
                    .ignoresSafeArea()
                    .onTapGesture {
                        if allowDismiss {
                            closeSidebar()
                        }
                    }
                    .animation(.easeInOut(duration: 0.3), value: isOpen)
            }

            // Sidebar
            HStack(spacing: 0) {
                if position == .trailing {
                    Spacer()
                }

                sidebarContent
                    .frame(width: sidebarWidth)
                    .offset(x: sidebarOffset)
                    .animation(.easeInOut(duration: 0.3), value: isOpen)
                    .gesture(
                        DragGesture()
                            .onChanged { value in
                                if allowDismiss {
                                    let translation = value.translation.x
                                    let maxDrag = sidebarWidth * 0.3

                                    if position == .leading {
                                        dragOffset = max(-maxDrag, min(0, translation))
                                    } else {
                                        dragOffset = min(maxDrag, max(0, translation))
                                    }
                                }
                            }
                            .onEnded { value in
                                let translation = value.translation.x
                                let threshold = sidebarWidth * 0.3

                                withAnimation(.easeInOut(duration: 0.3)) {
                                    if position == .leading && translation < -threshold {
                                        closeSidebar()
                                    } else if position == .trailing && translation > threshold {
                                        closeSidebar()
                                    }

                                    dragOffset = 0
                                }
                            }
                    )

                if position == .leading {
                    Spacer()
                }
            }
        }
    }

    // MARK: - Sidebar Content
    @ViewBuilder
    private var sidebarContent: some View {
        VStack(spacing: 0) {
            // Header
            if let header = header {
                header
                    .padding(.horizontal, Spacing.md)
                    .padding(.top, Spacing.md)
            }

            // Navigation Items
            ScrollView {
                LazyVStack(spacing: 0) {
                    ForEach(items, id: \.id) { item in
                        sidebarItemView(item, level: 0)
                    }
                }
                .padding(.vertical, Spacing.sm)
            }

            // Footer
            if let footer = footer {
                footer
                    .padding(.horizontal, Spacing.md)
                    .padding(.bottom, Spacing.md)
            }
        }
        .background(colors.background)
        .overlay(
            Rectangle()
                .fill(colors.border)
                .frame(width: 1),
            alignment: position == .leading ? .trailing : .leading
        )
    }

    // MARK: - Sidebar Item View
    @ViewBuilder
    private func sidebarItemView(_ item: SidebarItem, level: Int) -> some View {
        VStack(spacing: 0) {
            // Main item
            Button(action: {
                if !item.children.isEmpty {
                    toggleSection(item.id)
                } else {
                    selectItem(item.id)
                }
            }) {
                HStack(spacing: Spacing.sm) {
                    // Indentation
                    if level > 0 {
                        Rectangle()
                            .fill(Color.clear)
                            .frame(width: CGFloat(level * 16))
                    }

                    // Icon
                    if let icon = item.icon {
                        Image(systemName: icon)
                            .font(.system(size: 16))
                            .foregroundColor(itemIconColor(item))
                            .frame(width: 20)
                    }

                    // Title
                    Text(item.title)
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(itemTextColor(item))
                        .frame(maxWidth: .infinity, alignment: .leading)

                    // Badge
                    if let badge = item.badge {
                        Text(badge)
                            .font(.caption2)
                            .foregroundColor(colors.textOnPrimary)
                            .padding(.horizontal, 6)
                            .padding(.vertical, 2)
                            .background(colors.error)
                            .clipShape(Capsule())
                    }

                    // Expand/collapse indicator
                    if !item.children.isEmpty {
                        Image(systemName: expandedSections.contains(item.id) ? "chevron.down" : "chevron.right")
                            .font(.system(size: 12))
                            .foregroundColor(colors.textSecondary)
                            .animation(.easeInOut(duration: 0.2), value: expandedSections.contains(item.id))
                    }
                }
                .padding(.horizontal, Spacing.md)
                .padding(.vertical, Spacing.sm)
                .background(itemBackground(item))
                .contentShape(Rectangle())
            }
            .disabled(item.isDisabled)
            .buttonStyle(PlainButtonStyle())

            // Children (if expanded)
            if !item.children.isEmpty && expandedSections.contains(item.id) {
                ForEach(item.children, id: \.id) { child in
                    sidebarItemView(child, level: level + 1)
                }
                .transition(.opacity.combined(with: .move(edge: .top)))
                .animation(.easeInOut(duration: 0.2), value: expandedSections)
            }
        }
    }

    // MARK: - Styling Methods
    private func itemBackground(_ item: SidebarItem) -> Color {
        if selectedItem == item.id {
            return colors.primary.opacity(0.1)
        }
        return Color.clear
    }

    private func itemTextColor(_ item: SidebarItem) -> Color {
        if item.isDisabled {
            return colors.textDisabled
        }
        if selectedItem == item.id {
            return colors.primary
        }
        return colors.textPrimary
    }

    private func itemIconColor(_ item: SidebarItem) -> Color {
        if item.isDisabled {
            return colors.textDisabled
        }
        if selectedItem == item.id {
            return colors.primary
        }
        return colors.textSecondary
    }

    // MARK: - Actions
    private func selectItem(_ itemId: String) {
        selectedItem = itemId
        onSelectionChange?(itemId)

        if let item = findItem(itemId) {
            item.action?()
        }

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()

        // Close sidebar if in overlay mode
        if mode == .overlay {
            closeSidebar()
        }
    }

    private func toggleSection(_ sectionId: String) {
        withAnimation(.easeInOut(duration: 0.2)) {
            if expandedSections.contains(sectionId) {
                expandedSections.remove(sectionId)
            } else {
                expandedSections.insert(sectionId)
            }
        }

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }

    private func closeSidebar() {
        withAnimation(.easeInOut(duration: 0.3)) {
            isOpen = false
        }
    }

    private func findItem(_ itemId: String) -> SidebarItem? {
        func searchItems(_ items: [SidebarItem]) -> SidebarItem? {
            for item in items {
                if item.id == itemId {
                    return item
                }
                if let found = searchItems(item.children) {
                    return found
                }
            }
            return nil
        }

        return searchItems(items)
    }
}

// MARK: - Convenience Initializers
extension TchatSidebar {
    public init<Header: View, Footer: View>(
        isOpen: Binding<Bool>,
        selectedItem: Binding<String?> = .constant(nil),
        items: [SidebarItem],
        mode: SidebarMode = .overlay,
        position: SidebarPosition = .leading,
        size: SidebarSize = .standard,
        showOverlay: Bool = true,
        allowDismiss: Bool = true,
        @ViewBuilder header: () -> Header,
        @ViewBuilder footer: () -> Footer,
        onSelectionChange: ((String?) -> Void)? = nil
    ) {
        self.init(
            isOpen: isOpen,
            selectedItem: selectedItem,
            items: items,
            mode: mode,
            position: position,
            size: size,
            showOverlay: showOverlay,
            allowDismiss: allowDismiss,
            header: AnyView(header()),
            footer: AnyView(footer()),
            onSelectionChange: onSelectionChange
        )
    }
}

// MARK: - Preview
#if DEBUG
struct TchatSidebar_Previews: PreviewProvider {
    static var previews: some View {
        ZStack {
            Color(.systemGray6)
                .ignoresSafeArea()

            VStack {
                Text("Main Content")
                    .font(.title)
                    .foregroundColor(.primary)
            }

            TchatSidebar(
                isOpen: .constant(true),
                selectedItem: .constant("dashboard"),
                items: [
                    TchatSidebar.SidebarItem(
                        id: "dashboard",
                        title: "Dashboard",
                        icon: "chart.bar"
                    ),
                    TchatSidebar.SidebarItem(
                        id: "messages",
                        title: "Messages",
                        icon: "message",
                        badge: "5"
                    ),
                    TchatSidebar.SidebarItem(
                        id: "projects",
                        title: "Projects",
                        icon: "folder",
                        children: [
                            TchatSidebar.SidebarItem(
                                id: "active-projects",
                                title: "Active Projects",
                                icon: "circle.fill"
                            ),
                            TchatSidebar.SidebarItem(
                                id: "completed-projects",
                                title: "Completed",
                                icon: "checkmark.circle"
                            )
                        ]
                    ),
                    TchatSidebar.SidebarItem(
                        id: "team",
                        title: "Team",
                        icon: "person.2",
                        children: [
                            TchatSidebar.SidebarItem(
                                id: "members",
                                title: "Members",
                                icon: "person"
                            ),
                            TchatSidebar.SidebarItem(
                                id: "roles",
                                title: "Roles & Permissions",
                                icon: "key"
                            )
                        ]
                    ),
                    TchatSidebar.SidebarItem(
                        id: "settings",
                        title: "Settings",
                        icon: "gearshape"
                    )
                ],
                header: {
                    VStack(alignment: .leading, spacing: Spacing.sm) {
                        HStack {
                            Image(systemName: "app.badge")
                                .font(.title2)
                                .foregroundColor(.primary)

                            Text("Tchat")
                                .font(.headline)
                                .fontWeight(.bold)
                        }

                        Divider()
                    }
                },
                footer: {
                    VStack(spacing: Spacing.sm) {
                        Divider()

                        HStack {
                            Image(systemName: "person.circle")
                                .font(.title3)

                            VStack(alignment: .leading, spacing: 2) {
                                Text("John Doe")
                                    .font(.caption)
                                    .fontWeight(.medium)

                                Text("john@example.com")
                                    .font(.caption2)
                                    .foregroundColor(.secondary)
                            }

                            Spacer()
                        }
                    }
                }
            )
        }
        .previewLayout(.sizeThatFits)
    }
}
#endif