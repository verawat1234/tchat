//
//  TchatTabs.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Tab navigation component following Tchat design system
public struct TchatTabs: View {

    // MARK: - Tab Types
    public enum TabStyle {
        case line
        case pill
        case card
    }

    public enum TabSize {
        case small
        case medium
        case large
    }

    public enum TabPosition {
        case top
        case bottom
    }

    // MARK: - Tab Item
    public struct TabItem {
        let id: String
        let title: String
        let icon: String?
        let badge: String?
        let isDisabled: Bool
        let content: AnyView

        public init<Content: View>(
            id: String,
            title: String,
            icon: String? = nil,
            badge: String? = nil,
            isDisabled: Bool = false,
            @ViewBuilder content: () -> Content
        ) {
            self.id = id
            self.title = title
            self.icon = icon
            self.badge = badge
            self.isDisabled = isDisabled
            self.content = AnyView(content())
        }
    }

    // MARK: - Properties
    @Binding private var selectedTab: String
    @State private var tabWidths: [String: CGFloat] = [:]
    @State private var scrollOffset: CGFloat = 0

    let tabs: [TabItem]
    let style: TabStyle
    let size: TabSize
    let position: TabPosition
    let isScrollable: Bool
    let showDivider: Bool
    let onChange: ((String) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Computed Properties
    private var tabHeight: CGFloat {
        switch size {
        case .small: return 32
        case .medium: return 40
        case .large: return 48
        }
    }

    private var tabPadding: CGFloat {
        switch size {
        case .small: return 8
        case .medium: return 12
        case .large: return 16
        }
    }

    private var fontSize: CGFloat {
        switch size {
        case .small: return 12
        case .medium: return 14
        case .large: return 16
        }
    }

    private var iconSize: CGFloat {
        switch size {
        case .small: return 14
        case .medium: return 16
        case .large: return 18
        }
    }

    // MARK: - Initializer
    public init(
        selectedTab: Binding<String>,
        tabs: [TabItem],
        style: TabStyle = .line,
        size: TabSize = .medium,
        position: TabPosition = .top,
        isScrollable: Bool = false,
        showDivider: Bool = true,
        onChange: ((String) -> Void)? = nil
    ) {
        self._selectedTab = selectedTab
        self.tabs = tabs
        self.style = style
        self.size = size
        self.position = position
        self.isScrollable = isScrollable
        self.showDivider = showDivider
        self.onChange = onChange
    }

    // MARK: - Body
    public var body: some View {
        VStack(spacing: 0) {
            if position == .top {
                tabHeader
                if showDivider && style == .line {
                    divider
                }
                tabContent
            } else {
                tabContent
                if showDivider && style == .line {
                    divider
                }
                tabHeader
            }
        }
    }

    // MARK: - Tab Header
    @ViewBuilder
    private var tabHeader: some View {
        if isScrollable {
            ScrollViewReader { proxy in
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 0) {
                        ForEach(tabs, id: \.id) { tab in
                            tabButton(for: tab)
                                .id(tab.id)
                        }
                    }
                    .padding(.horizontal, Spacing.md)
                }
                .onChange(of: selectedTab) { newValue in
                    withAnimation(.easeInOut(duration: 0.3)) {
                        proxy.scrollTo(newValue, anchor: .center)
                    }
                }
            }
        } else {
            HStack(spacing: 0) {
                ForEach(tabs, id: \.id) { tab in
                    tabButton(for: tab)
                        .frame(maxWidth: .infinity)
                }
            }
            .padding(.horizontal, Spacing.md)
        }
    }

    // MARK: - Tab Button
    @ViewBuilder
    private func tabButton(for tab: TabItem) -> some View {
        Button(action: {
            if !tab.isDisabled {
                selectedTab = tab.id
                onChange?(tab.id)

                // Haptic feedback
                let impactFeedback = UIImpactFeedbackGenerator(style: .light)
                impactFeedback.impactOccurred()
            }
        }) {
            HStack(spacing: Spacing.xs) {
                if let icon = tab.icon {
                    Image(systemName: icon)
                        .font(.system(size: iconSize))
                        .foregroundColor(tabIconColor(for: tab))
                }

                Text(tab.title)
                    .font(.system(size: fontSize, weight: .medium))
                    .foregroundColor(tabTextColor(for: tab))

                if let badge = tab.badge {
                    Text(badge)
                        .font(.caption2)
                        .foregroundColor(colors.textOnPrimary)
                        .padding(.horizontal, 6)
                        .padding(.vertical, 2)
                        .background(colors.error)
                        .clipShape(Capsule())
                }
            }
            .padding(.horizontal, tabPadding)
            .padding(.vertical, Spacing.xs)
            .frame(height: tabHeight)
            .background(tabBackground(for: tab))
            .overlay(
                tabIndicator(for: tab),
                alignment: position == .top ? .bottom : .top
            )
        }
        .disabled(tab.isDisabled)
        .animation(.easeInOut(duration: 0.2), value: selectedTab)
    }

    // MARK: - Tab Styling
    private func tabBackground(for tab: TabItem) -> some View {
        Group {
            switch style {
            case .line:
                Color.clear
            case .pill:
                RoundedRectangle(cornerRadius: tabHeight / 2)
                    .fill(isSelected(tab) ? colors.primary.opacity(0.1) : Color.clear)
            case .card:
                RoundedRectangle(cornerRadius: 8)
                    .fill(isSelected(tab) ? colors.background : colors.surface)
                    .shadow(color: colors.shadowLight, radius: isSelected(tab) ? 2 : 0, y: 1)
            }
        }
    }

    private func tabIndicator(for tab: TabItem) -> some View {
        Group {
            if style == .line && isSelected(tab) {
                Rectangle()
                    .fill(colors.primary)
                    .frame(height: 2)
                    .animation(.easeInOut(duration: 0.3), value: selectedTab)
            } else {
                EmptyView()
            }
        }
    }

    private func tabTextColor(for tab: TabItem) -> Color {
        if tab.isDisabled {
            return colors.textDisabled
        }

        switch style {
        case .line, .pill:
            return isSelected(tab) ? colors.primary : colors.textSecondary
        case .card:
            return isSelected(tab) ? colors.textPrimary : colors.textSecondary
        }
    }

    private func tabIconColor(for tab: TabItem) -> Color {
        if tab.isDisabled {
            return colors.textDisabled
        }

        return isSelected(tab) ? colors.primary : colors.textSecondary
    }

    private func isSelected(_ tab: TabItem) -> Bool {
        selectedTab == tab.id
    }

    // MARK: - Tab Content
    @ViewBuilder
    private var tabContent: some View {
        if let selectedTabItem = tabs.first(where: { $0.id == selectedTab }) {
            selectedTabItem.content
                .transition(.opacity.combined(with: .move(edge: .trailing)))
                .animation(.easeInOut(duration: 0.2), value: selectedTab)
        }
    }

    // MARK: - Divider
    @ViewBuilder
    private var divider: some View {
        Rectangle()
            .fill(colors.border)
            .frame(height: 1)
    }
}

// MARK: - Tab Bar (Bottom Navigation)
public struct TchatTabBar: View {
    @Binding private var selectedTab: String

    let tabs: [TchatTabs.TabItem]
    let size: TchatTabs.TabSize
    let onChange: ((String) -> Void)?

    private let colors = Colors()

    public init(
        selectedTab: Binding<String>,
        tabs: [TchatTabs.TabItem],
        size: TchatTabs.TabSize = .medium,
        onChange: ((String) -> Void)? = nil
    ) {
        self._selectedTab = selectedTab
        self.tabs = tabs
        self.size = size
        self.onChange = onChange
    }

    public var body: some View {
        HStack {
            ForEach(tabs, id: \.id) { tab in
                Button(action: {
                    if !tab.isDisabled {
                        selectedTab = tab.id
                        onChange?(tab.id)

                        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
                        impactFeedback.impactOccurred()
                    }
                }) {
                    VStack(spacing: 4) {
                        ZStack {
                            if let icon = tab.icon {
                                Image(systemName: icon)
                                    .font(.system(size: iconSize))
                                    .foregroundColor(iconColor(for: tab))
                            }

                            if let badge = tab.badge {
                                Text(badge)
                                    .font(.caption2)
                                    .foregroundColor(colors.textOnPrimary)
                                    .padding(.horizontal, 4)
                                    .padding(.vertical, 2)
                                    .background(colors.error)
                                    .clipShape(Capsule())
                                    .offset(x: 12, y: -8)
                            }
                        }

                        Text(tab.title)
                            .font(.caption)
                            .foregroundColor(textColor(for: tab))
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, Spacing.xs)
                }
                .disabled(tab.isDisabled)
            }
        }
        .padding(.horizontal, Spacing.md)
        .padding(.top, Spacing.xs)
        .background(colors.background.ignoresSafeArea())
        .overlay(
            Rectangle()
                .fill(colors.border)
                .frame(height: 1),
            alignment: .top
        )
    }

    private var iconSize: CGFloat {
        switch size {
        case .small: return 16
        case .medium: return 20
        case .large: return 24
        }
    }

    private func iconColor(for tab: TchatTabs.TabItem) -> Color {
        if tab.isDisabled {
            return colors.textDisabled
        }
        return selectedTab == tab.id ? colors.primary : colors.textSecondary
    }

    private func textColor(for tab: TchatTabs.TabItem) -> Color {
        if tab.isDisabled {
            return colors.textDisabled
        }
        return selectedTab == tab.id ? colors.primary : colors.textSecondary
    }
}

// MARK: - Preview
#if DEBUG
struct TchatTabs_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Line tabs
                TchatTabs(
                    selectedTab: .constant("tab1"),
                    tabs: [
                        TchatTabs.TabItem(id: "tab1", title: "Overview", icon: "chart.bar") {
                            Text("Overview Content")
                                .frame(maxWidth: .infinity, minHeight: 200)
                                .background(Color(.systemGray6))
                        },
                        TchatTabs.TabItem(id: "tab2", title: "Analytics", icon: "chart.line.uptrend.xyaxis") {
                            Text("Analytics Content")
                                .frame(maxWidth: .infinity, minHeight: 200)
                                .background(Color(.systemGray6))
                        },
                        TchatTabs.TabItem(id: "tab3", title: "Settings", icon: "gearshape", badge: "2") {
                            Text("Settings Content")
                                .frame(maxWidth: .infinity, minHeight: 200)
                                .background(Color(.systemGray6))
                        }
                    ],
                    style: .line
                )

                // Pill tabs
                TchatTabs(
                    selectedTab: .constant("pill1"),
                    tabs: [
                        TchatTabs.TabItem(id: "pill1", title: "All") {
                            Text("All Items")
                                .frame(maxWidth: .infinity, minHeight: 100)
                        },
                        TchatTabs.TabItem(id: "pill2", title: "Active") {
                            Text("Active Items")
                                .frame(maxWidth: .infinity, minHeight: 100)
                        },
                        TchatTabs.TabItem(id: "pill3", title: "Completed") {
                            Text("Completed Items")
                                .frame(maxWidth: .infinity, minHeight: 100)
                        }
                    ],
                    style: .pill,
                    size: .small
                )

                // Card tabs
                TchatTabs(
                    selectedTab: .constant("card1"),
                    tabs: [
                        TchatTabs.TabItem(id: "card1", title: "Profile", icon: "person") {
                            Text("Profile Content")
                                .frame(maxWidth: .infinity, minHeight: 150)
                        },
                        TchatTabs.TabItem(id: "card2", title: "Messages", icon: "message", badge: "5") {
                            Text("Messages Content")
                                .frame(maxWidth: .infinity, minHeight: 150)
                        }
                    ],
                    style: .card,
                    size: .large
                )

                Divider()

                // Tab bar
                TchatTabBar(
                    selectedTab: .constant("home"),
                    tabs: [
                        TchatTabs.TabItem(id: "home", title: "Home", icon: "house") { EmptyView() },
                        TchatTabs.TabItem(id: "search", title: "Search", icon: "magnifyingglass") { EmptyView() },
                        TchatTabs.TabItem(id: "notifications", title: "Notifications", icon: "bell", badge: "3") { EmptyView() },
                        TchatTabs.TabItem(id: "profile", title: "Profile", icon: "person") { EmptyView() }
                    ]
                )
            }
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif