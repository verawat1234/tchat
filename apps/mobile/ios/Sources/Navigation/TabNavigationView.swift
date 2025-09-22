//
//  TabNavigationView.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Main tab navigation view following Tchat design system
public struct TabNavigationView: View {

    // MARK: - Tab Items
    public enum Tab: String, CaseIterable, Identifiable {
        case chat = "Chat"
        case store = "Store"
        case social = "Social"
        case video = "Video"
        case more = "More"

        public var id: String { rawValue }

        var systemImage: String {
            switch self {
            case .chat: return "message.fill"
            case .store: return "bag.fill"
            case .social: return "person.3.fill"
            case .video: return "video.fill"
            case .more: return "ellipsis"
            }
        }

        var activeSystemImage: String {
            switch self {
            case .chat: return "message.fill"
            case .store: return "bag.fill"
            case .social: return "person.3.fill"
            case .video: return "video.fill"
            case .more: return "ellipsis"
            }
        }
    }

    // MARK: - Properties
    @State private var selectedTab: Tab = .chat
    @State private var previousTab: Tab = .chat

    // MARK: - Private Properties
    private let colors = Colors()

    // MARK: - Body
    public var body: some View {
        TabView(selection: $selectedTab) {
            ForEach(Tab.allCases) { tab in
                tabContent(for: tab)
                    .tabItem {
                        tabLabel(for: tab)
                    }
                    .tag(tab)
            }
        }
        .accentColor(colors.primary)
        .onAppear {
            configureTabBarAppearance()
        }
        .onChange(of: selectedTab) { newTab in
            // Handle tab change analytics or side effects
            handleTabChange(from: previousTab, to: newTab)
            previousTab = newTab
        }
    }

    // MARK: - Tab Content
    @ViewBuilder
    private func tabContent(for tab: Tab) -> some View {
        NavigationView {
            switch tab {
            case .chat:
                ChatTabView()
            case .store:
                StoreTabView()
            case .social:
                SocialTabView()
            case .video:
                VideoTabView()
            case .more:
                MoreTabView()
            }
        }
        .navigationViewStyle(StackNavigationViewStyle())
    }

    // MARK: - Tab Label
    @ViewBuilder
    private func tabLabel(for tab: Tab) -> some View {
        VStack(spacing: 2) {
            Image(systemName: selectedTab == tab ? tab.activeSystemImage : tab.systemImage)
                .font(.system(size: 20, weight: .medium))

            Text(tab.rawValue)
                .font(.caption2)
        }
    }

    // MARK: - Configuration
    private func configureTabBarAppearance() {
        let appearance = UITabBarAppearance()

        // Configure background
        appearance.configureWithOpaqueBackground()
        appearance.backgroundColor = UIColor(colors.tabBarBackground)

        // Configure shadow
        appearance.shadowColor = UIColor(colors.shadowLight)

        // Configure item appearance
        appearance.stackedLayoutAppearance.normal.iconColor = UIColor(colors.tabUnselected)
        appearance.stackedLayoutAppearance.normal.titleTextAttributes = [
            .foregroundColor: UIColor(colors.tabUnselected),
            .font: UIFont.systemFont(ofSize: 10, weight: .medium)
        ]

        appearance.stackedLayoutAppearance.selected.iconColor = UIColor(colors.tabSelected)
        appearance.stackedLayoutAppearance.selected.titleTextAttributes = [
            .foregroundColor: UIColor(colors.tabSelected),
            .font: UIFont.systemFont(ofSize: 10, weight: .semibold)
        ]

        // Apply appearance
        UITabBar.appearance().standardAppearance = appearance
        UITabBar.appearance().scrollEdgeAppearance = appearance
    }

    // MARK: - Tab Change Handler
    private func handleTabChange(from previousTab: Tab, to newTab: Tab) {
        // Analytics tracking
        print("Tab changed from \(previousTab.rawValue) to \(newTab.rawValue)")

        // Add haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()

        // Additional side effects can be added here
        // - State synchronization
        // - Badge updates
        // - Background tasks
    }
}

// MARK: - Tab Content Views

/// Chat tab content view
struct ChatTabView: View {
    private let colors = Colors()

    var body: some View {
        VStack(spacing: Spacing.lg) {
            Text("Chat")
                .font(.largeTitle)
                .fontWeight(.bold)
                .foregroundColor(colors.textPrimary)

            Text("Messages and conversations")
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Placeholder content
            VStack(spacing: Spacing.md) {
                Image(systemName: "message.circle.fill")
                    .font(.system(size: 64))
                    .foregroundColor(colors.primary)

                Text("Your chat interface will be here")
                    .multilineTextAlignment(.center)
                    .foregroundColor(colors.textTertiary)
            }

            Spacer()
        }
        .padding()
        .navigationTitle("Chat")
        .navigationBarTitleDisplayMode(.large)
    }
}

/// Store tab content view
struct StoreTabView: View {
    private let colors = Colors()

    var body: some View {
        VStack(spacing: Spacing.lg) {
            Text("Store")
                .font(.largeTitle)
                .fontWeight(.bold)
                .foregroundColor(colors.textPrimary)

            Text("Browse and purchase items")
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Placeholder content
            VStack(spacing: Spacing.md) {
                Image(systemName: "bag.circle.fill")
                    .font(.system(size: 64))
                    .foregroundColor(colors.primary)

                Text("Your store interface will be here")
                    .multilineTextAlignment(.center)
                    .foregroundColor(colors.textTertiary)
            }

            Spacer()
        }
        .padding()
        .navigationTitle("Store")
        .navigationBarTitleDisplayMode(.large)
    }
}

/// Social tab content view
struct SocialTabView: View {
    private let colors = Colors()

    var body: some View {
        VStack(spacing: Spacing.lg) {
            Text("Social")
                .font(.largeTitle)
                .fontWeight(.bold)
                .foregroundColor(colors.textPrimary)

            Text("Connect with friends")
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Placeholder content
            VStack(spacing: Spacing.md) {
                Image(systemName: "person.3.circle.fill")
                    .font(.system(size: 64))
                    .foregroundColor(colors.primary)

                Text("Your social interface will be here")
                    .multilineTextAlignment(.center)
                    .foregroundColor(colors.textTertiary)
            }

            Spacer()
        }
        .padding()
        .navigationTitle("Social")
        .navigationBarTitleDisplayMode(.large)
    }
}

/// Video tab content view
struct VideoTabView: View {
    private let colors = Colors()

    var body: some View {
        VStack(spacing: Spacing.lg) {
            Text("Video")
                .font(.largeTitle)
                .fontWeight(.bold)
                .foregroundColor(colors.textPrimary)

            Text("Watch and share videos")
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Placeholder content
            VStack(spacing: Spacing.md) {
                Image(systemName: "video.circle.fill")
                    .font(.system(size: 64))
                    .foregroundColor(colors.primary)

                Text("Your video interface will be here")
                    .multilineTextAlignment(.center)
                    .foregroundColor(colors.textTertiary)
            }

            Spacer()
        }
        .padding()
        .navigationTitle("Video")
        .navigationBarTitleDisplayMode(.large)
    }
}

/// More tab content view
struct MoreTabView: View {
    private let colors = Colors()

    var body: some View {
        VStack(spacing: Spacing.lg) {
            Text("More")
                .font(.largeTitle)
                .fontWeight(.bold)
                .foregroundColor(colors.textPrimary)

            Text("Settings and additional features")
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Placeholder content with settings items
            VStack(spacing: Spacing.md) {
                TchatCard(variant: .outlined) {
                    VStack(alignment: .leading, spacing: Spacing.sm) {
                        HStack {
                            Image(systemName: "gear")
                                .foregroundColor(colors.primary)
                            Text("Settings")
                                .fontWeight(.medium)
                            Spacer()
                            Image(systemName: "chevron.right")
                                .foregroundColor(colors.textTertiary)
                        }
                    }
                }

                TchatCard(variant: .outlined) {
                    VStack(alignment: .leading, spacing: Spacing.sm) {
                        HStack {
                            Image(systemName: "person.circle")
                                .foregroundColor(colors.primary)
                            Text("Profile")
                                .fontWeight(.medium)
                            Spacer()
                            Image(systemName: "chevron.right")
                                .foregroundColor(colors.textTertiary)
                        }
                    }
                }

                TchatCard(variant: .outlined) {
                    VStack(alignment: .leading, spacing: Spacing.sm) {
                        HStack {
                            Image(systemName: "info.circle")
                                .foregroundColor(colors.primary)
                            Text("About")
                                .fontWeight(.medium)
                            Spacer()
                            Image(systemName: "chevron.right")
                                .foregroundColor(colors.textTertiary)
                        }
                    }
                }
            }

            Spacer()
        }
        .padding()
        .navigationTitle("More")
        .navigationBarTitleDisplayMode(.large)
    }
}

// MARK: - Preview
#if DEBUG
struct TabNavigationView_Previews: PreviewProvider {
    static var previews: some View {
        TabNavigationView()
            .previewDevice("iPhone 14")
    }
}
#endif