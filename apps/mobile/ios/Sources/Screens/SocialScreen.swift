//
//  SocialScreen.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Social media feed interface screen
public struct SocialScreen: View {
    @State private var selectedTab: SocialTab = .feed
    @State private var showingCamera = false

    private let colors = Colors()
    private let spacing = Spacing()

    // Mock posts
    private let posts = [
        ("Alice Johnson", "2h", "Just finished my morning workout! 💪 #fitness", "figure.walk", 23, 5),
        ("Bob Smith", "4h", "Beautiful sunset from my balcony 🌅", "sun.max.fill", 67, 12),
        ("Carol Davis", "6h", "New coffee shop opened downtown! ☕️ Must try", "cup.and.saucer.fill", 45, 8),
        ("David Wilson", "8h", "Working on a new project. Excited to share soon! 🚀", "laptop", 89, 15),
        ("Emma Brown", "12h", "Weekend hiking adventure! Nature is amazing 🏔️", "mountain.2.fill", 156, 28)
    ]

    public init() {}

    public var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Segmented control
                HStack(spacing: 0) {
                    ForEach(SocialTab.allCases, id: \.self) { tab in
                        Button(action: {
                            selectedTab = tab
                        }) {
                            VStack(spacing: spacing.xs) {
                                Image(systemName: tab.icon)
                                    .font(.system(size: 16))
                                Text(tab.title)
                                    .font(.system(size: 12, weight: .medium))
                            }
                            .foregroundColor(selectedTab == tab ? colors.primary : colors.textSecondary)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, spacing.sm)
                        }
                    }
                }
                .background(colors.surface)

                // Content based on selected tab
                switch selectedTab {
                case .feed:
                    FeedView(posts: posts)
                case .discover:
                    DiscoverView()
                case .notifications:
                    NotificationsView()
                }
            }
            .navigationTitle("Social")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {
                        showingCamera = true
                    }) {
                        Image(systemName: "camera.fill")
                            .foregroundColor(colors.primary)
                    }
                }
            }
        }
        .navigationViewStyle(StackNavigationViewStyle())
        .sheet(isPresented: $showingCamera) {
            CameraView()
        }
    }
}

// MARK: - Social Tab Enum
enum SocialTab: CaseIterable {
    case feed, discover, notifications

    var title: String {
        switch self {
        case .feed: return "Feed"
        case .discover: return "Discover"
        case .notifications: return "Alerts"
        }
    }

    var icon: String {
        switch self {
        case .feed: return "house.fill"
        case .discover: return "safari.fill"
        case .notifications: return "bell.fill"
        }
    }
}

// MARK: - Feed View
private struct FeedView: View {
    let posts: [(String, String, String, String, Int, Int)]

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        ScrollView {
            LazyVStack(spacing: spacing.md) {
                ForEach(Array(posts.enumerated()), id: \.offset) { index, post in
                    PostCard(
                        author: post.0,
                        time: post.1,
                        content: post.2,
                        icon: post.3,
                        likes: post.4,
                        comments: post.5
                    )
                }
            }
            .padding(.horizontal, spacing.md)
            .padding(.top, spacing.sm)
        }
    }
}

// MARK: - Post Card Component
private struct PostCard: View {
    let author: String
    let time: String
    let content: String
    let icon: String
    let likes: Int
    let comments: Int

    @State private var isLiked = false
    @State private var currentLikes: Int

    private let colors = Colors()
    private let spacing = Spacing()

    init(author: String, time: String, content: String, icon: String, likes: Int, comments: Int) {
        self.author = author
        self.time = time
        self.content = content
        self.icon = icon
        self.likes = likes
        self.comments = comments
        self._currentLikes = State(initialValue: likes)
    }

    var body: some View {
        VStack(alignment: .leading, spacing: spacing.sm) {
            // Header
            HStack {
                Circle()
                    .fill(colors.primary.opacity(0.2))
                    .frame(width: 40, height: 40)
                    .overlay(
                        Text(author.prefix(1))
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(colors.primary)
                    )

                VStack(alignment: .leading, spacing: 2) {
                    Text(author)
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(colors.textPrimary)

                    Text(time)
                        .font(.system(size: 12))
                        .foregroundColor(colors.textSecondary)
                }

                Spacer()

                Button(action: {}) {
                    Image(systemName: "ellipsis")
                        .foregroundColor(colors.textSecondary)
                }
            }

            // Content
            Text(content)
                .font(.system(size: 15))
                .foregroundColor(colors.textPrimary)
                .multilineTextAlignment(.leading)

            // Icon/Media placeholder
            Image(systemName: icon)
                .font(.system(size: 24))
                .foregroundColor(colors.primary)
                .frame(maxWidth: .infinity, alignment: .center)
                .padding(.vertical, spacing.lg)
                .background(colors.surface)
                .cornerRadius(12)

            // Actions
            HStack(spacing: spacing.lg) {
                Button(action: {
                    isLiked.toggle()
                    currentLikes += isLiked ? 1 : -1
                }) {
                    HStack(spacing: spacing.xs) {
                        Image(systemName: isLiked ? "heart.fill" : "heart")
                            .foregroundColor(isLiked ? .red : colors.textSecondary)
                        Text("\(currentLikes)")
                            .font(.system(size: 14))
                            .foregroundColor(colors.textSecondary)
                    }
                }

                Button(action: {}) {
                    HStack(spacing: spacing.xs) {
                        Image(systemName: "bubble.left")
                            .foregroundColor(colors.textSecondary)
                        Text("\(comments)")
                            .font(.system(size: 14))
                            .foregroundColor(colors.textSecondary)
                    }
                }

                Button(action: {}) {
                    Image(systemName: "square.and.arrow.up")
                        .foregroundColor(colors.textSecondary)
                }

                Spacer()
            }
        }
        .padding(spacing.md)
        .background(Color.white)
        .cornerRadius(12)
        .shadow(color: colors.shadowLight, radius: 4, y: 2)
    }
}

// MARK: - Discover View
private struct DiscoverView: View {
    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        ScrollView {
            VStack(spacing: spacing.md) {
                Text("Trending Topics")
                    .font(.system(size: 18, weight: .bold))
                    .foregroundColor(colors.textPrimary)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal, spacing.md)

                // Trending grid placeholder
                LazyVGrid(columns: [
                    GridItem(.flexible()),
                    GridItem(.flexible())
                ], spacing: spacing.sm) {
                    ForEach(0..<6) { index in
                        RoundedRectangle(cornerRadius: 12)
                            .fill(colors.surface)
                            .frame(height: 120)
                            .overlay(
                                VStack {
                                    Image(systemName: "flame.fill")
                                        .font(.system(size: 24))
                                        .foregroundColor(colors.primary)
                                    Text("Topic \(index + 1)")
                                        .font(.system(size: 14, weight: .medium))
                                        .foregroundColor(colors.textPrimary)
                                }
                            )
                    }
                }
                .padding(.horizontal, spacing.md)
            }
            .padding(.top, spacing.sm)
        }
    }
}

// MARK: - Notifications View
private struct NotificationsView: View {
    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        ScrollView {
            VStack(spacing: spacing.sm) {
                ForEach(0..<5) { index in
                    HStack(spacing: spacing.md) {
                        Circle()
                            .fill(colors.primary.opacity(0.2))
                            .frame(width: 40, height: 40)
                            .overlay(
                                Image(systemName: "bell.fill")
                                    .font(.system(size: 16))
                                    .foregroundColor(colors.primary)
                            )

                        VStack(alignment: .leading, spacing: spacing.xs) {
                            Text("Notification \(index + 1)")
                                .font(.system(size: 14, weight: .medium))
                                .foregroundColor(colors.textPrimary)

                            Text("This is a sample notification message")
                                .font(.system(size: 12))
                                .foregroundColor(colors.textSecondary)
                        }

                        Spacer()

                        Text("2m")
                            .font(.system(size: 12))
                            .foregroundColor(colors.textSecondary)
                    }
                    .padding(.horizontal, spacing.md)
                    .padding(.vertical, spacing.sm)
                }
            }
            .padding(.top, spacing.sm)
        }
    }
}

// MARK: - Camera View (Placeholder)
private struct CameraView: View {
    @Environment(\.presentationMode) var presentationMode

    private let colors = Colors()

    var body: some View {
        VStack {
            Text("Camera")
                .font(.title)
                .foregroundColor(colors.textPrimary)

            Button("Close") {
                presentationMode.wrappedValue.dismiss()
            }
            .foregroundColor(colors.primary)
        }
    }
}

// MARK: - Preview
#if DEBUG
struct SocialScreen_Previews: PreviewProvider {
    static var previews: some View {
        SocialScreen()
            .previewDisplayName("Social Screen")
    }
}
#endif