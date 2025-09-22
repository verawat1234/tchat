//
//  ChatScreen.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Main chat interface screen
public struct ChatScreen: View {
    @State private var searchText = ""
    @State private var selectedChat: String? = nil

    private let colors = Colors()
    private let spacing = Spacing()

    // Mock chat data
    private let chats = [
        ("John Doe", "Hey, how's it going?", "2m", true),
        ("Sarah Wilson", "Meeting at 3pm today", "15m", false),
        ("Team Alpha", "Project update ready", "1h", true),
        ("Mom", "Don't forget dinner tonight", "2h", false),
        ("Alex Chen", "Thanks for the help!", "3h", false)
    ]

    public init() {}

    public var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search bar
                HStack {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(colors.textSecondary)

                    TextField("Search conversations", text: $searchText)
                        .textFieldStyle(PlainTextFieldStyle())
                }
                .padding(.horizontal, spacing.md)
                .padding(.vertical, spacing.sm)
                .background(colors.surface)
                .cornerRadius(12)
                .padding(.horizontal, spacing.md)
                .padding(.top, spacing.sm)

                // Chat list
                List {
                    ForEach(Array(chats.enumerated()), id: \.offset) { index, chat in
                        ChatRowView(
                            name: chat.0,
                            lastMessage: chat.1,
                            time: chat.2,
                            hasUnread: chat.3,
                            isSelected: selectedChat == chat.0
                        )
                        .listRowInsets(EdgeInsets())
                        .listRowSeparator(.hidden)
                        .onTapGesture {
                            selectedChat = chat.0
                            // Navigate to chat detail
                        }
                    }
                }
                .listStyle(PlainListStyle())
            }
            .navigationTitle("Chats")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {
                        // New chat action
                    }) {
                        Image(systemName: "square.and.pencil")
                            .foregroundColor(colors.primary)
                    }
                }
            }
        }
        .navigationViewStyle(StackNavigationViewStyle())
    }
}

// MARK: - Chat Row Component
private struct ChatRowView: View {
    let name: String
    let lastMessage: String
    let time: String
    let hasUnread: Bool
    let isSelected: Bool

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        HStack(spacing: spacing.md) {
            // Avatar
            Circle()
                .fill(colors.primary.opacity(0.2))
                .frame(width: 48, height: 48)
                .overlay(
                    Text(name.prefix(1))
                        .font(.system(size: 20, weight: .semibold))
                        .foregroundColor(colors.primary)
                )

            // Chat info
            VStack(alignment: .leading, spacing: spacing.xs) {
                HStack {
                    Text(name)
                        .font(.system(size: 16, weight: hasUnread ? .semibold : .medium))
                        .foregroundColor(colors.textPrimary)

                    Spacer()

                    Text(time)
                        .font(.system(size: 14))
                        .foregroundColor(colors.textSecondary)
                }

                HStack {
                    Text(lastMessage)
                        .font(.system(size: 14))
                        .foregroundColor(colors.textSecondary)
                        .lineLimit(1)

                    Spacer()

                    if hasUnread {
                        Circle()
                            .fill(colors.primary)
                            .frame(width: 8, height: 8)
                    }
                }
            }
        }
        .padding(.horizontal, spacing.md)
        .padding(.vertical, spacing.sm)
        .background(isSelected ? colors.surface : Color.clear)
        .cornerRadius(12)
    }
}

// MARK: - Preview
#if DEBUG
struct ChatScreen_Previews: PreviewProvider {
    static var previews: some View {
        ChatScreen()
            .previewDisplayName("Chat Screen")
    }
}
#endif