//
//  MoreScreen.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Settings and additional features screen
public struct MoreScreen: View {
    @State private var showingProfile = false
    @State private var darkModeEnabled = false
    @State private var notificationsEnabled = true
    @State private var soundEnabled = true

    private let colors = Colors()
    private let spacing = Spacing()

    // Menu sections
    private let menuSections = [
        MoreSection(
            title: "Account",
            items: [
                MoreItem("Profile", "person.circle.fill", .profile),
                MoreItem("Privacy", "lock.shield.fill", .privacy),
                MoreItem("Security", "checkmark.shield.fill", .security),
                MoreItem("Billing", "creditcard.fill", .billing)
            ]
        ),
        MoreSection(
            title: "Preferences",
            items: [
                MoreItem("Notifications", "bell.fill", .notifications),
                MoreItem("Appearance", "paintbrush.fill", .appearance),
                MoreItem("Language", "globe", .language),
                MoreItem("Storage", "internaldrive.fill", .storage)
            ]
        ),
        MoreSection(
            title: "Support",
            items: [
                MoreItem("Help Center", "questionmark.circle.fill", .help),
                MoreItem("Contact Us", "envelope.fill", .contact),
                MoreItem("Report Issue", "exclamationmark.triangle.fill", .report),
                MoreItem("About", "info.circle.fill", .about)
            ]
        )
    ]

    public init() {}

    public var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: spacing.lg) {
                    // Profile header
                    ProfileHeaderView(showingProfile: $showingProfile)

                    // Menu sections
                    ForEach(Array(menuSections.enumerated()), id: \.offset) { index, section in
                        MenuSectionView(
                            section: section,
                            darkModeEnabled: $darkModeEnabled,
                            notificationsEnabled: $notificationsEnabled,
                            soundEnabled: $soundEnabled
                        )
                    }

                    // Sign out button
                    Button(action: {
                        // Sign out action
                    }) {
                        HStack {
                            Image(systemName: "rectangle.portrait.and.arrow.right")
                                .foregroundColor(colors.error)
                            Text("Sign Out")
                                .font(.system(size: 16, weight: .medium))
                                .foregroundColor(colors.error)
                        }
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, spacing.md)
                        .background(colors.surface)
                        .cornerRadius(12)
                    }
                    .padding(.horizontal, spacing.md)

                    // Version info
                    Text("Version 1.0.0 (Build 1)")
                        .font(.system(size: 12))
                        .foregroundColor(colors.textSecondary)
                        .padding(.bottom, spacing.xl)
                }
                .padding(.top, spacing.sm)
            }
            .navigationTitle("More")
            .navigationBarTitleDisplayMode(.large)
        }
        .navigationViewStyle(StackNavigationViewStyle())
        .sheet(isPresented: $showingProfile) {
            ProfileView()
        }
    }
}

// MARK: - Data Models
private struct MoreSection {
    let title: String
    let items: [MoreItem]
}

private struct MoreItem {
    let title: String
    let icon: String
    let type: MoreItemType

    init(_ title: String, _ icon: String, _ type: MoreItemType) {
        self.title = title
        self.icon = icon
        self.type = type
    }
}

private enum MoreItemType {
    case profile, privacy, security, billing
    case notifications, appearance, language, storage
    case help, contact, report, about
}

// MARK: - Profile Header View
private struct ProfileHeaderView: View {
    @Binding var showingProfile: Bool

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        Button(action: {
            showingProfile = true
        }) {
            HStack(spacing: spacing.md) {
                // Profile avatar
                Circle()
                    .fill(colors.primary.opacity(0.2))
                    .frame(width: 60, height: 60)
                    .overlay(
                        Text("JD")
                            .font(.system(size: 24, weight: .bold))
                            .foregroundColor(colors.primary)
                    )

                // Profile info
                VStack(alignment: .leading, spacing: spacing.xs) {
                    Text("John Doe")
                        .font(.system(size: 18, weight: .semibold))
                        .foregroundColor(colors.textPrimary)

                    Text("john.doe@example.com")
                        .font(.system(size: 14))
                        .foregroundColor(colors.textSecondary)

                    HStack(spacing: spacing.xs) {
                        Circle()
                            .fill(colors.success)
                            .frame(width: 8, height: 8)
                        Text("Online")
                            .font(.system(size: 12))
                            .foregroundColor(colors.success)
                    }
                }

                Spacer()

                Image(systemName: "chevron.right")
                    .foregroundColor(colors.textSecondary)
            }
            .padding(spacing.md)
            .background(Color.white)
            .cornerRadius(12)
            .shadow(color: colors.shadowLight, radius: 4, y: 2)
        }
        .padding(.horizontal, spacing.md)
    }
}

// MARK: - Menu Section View
private struct MenuSectionView: View {
    let section: MoreSection
    @Binding var darkModeEnabled: Bool
    @Binding var notificationsEnabled: Bool
    @Binding var soundEnabled: Bool

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        VStack(alignment: .leading, spacing: spacing.sm) {
            // Section title
            Text(section.title)
                .font(.system(size: 16, weight: .semibold))
                .foregroundColor(colors.textPrimary)
                .padding(.horizontal, spacing.md)

            // Section items
            VStack(spacing: 0) {
                ForEach(Array(section.items.enumerated()), id: \.offset) { index, item in
                    MenuItemView(
                        item: item,
                        isLast: index == section.items.count - 1,
                        darkModeEnabled: $darkModeEnabled,
                        notificationsEnabled: $notificationsEnabled,
                        soundEnabled: $soundEnabled
                    )
                }
            }
            .background(Color.white)
            .cornerRadius(12)
            .shadow(color: colors.shadowLight, radius: 4, y: 2)
            .padding(.horizontal, spacing.md)
        }
    }
}

// MARK: - Menu Item View
private struct MenuItemView: View {
    let item: MoreItem
    let isLast: Bool
    @Binding var darkModeEnabled: Bool
    @Binding var notificationsEnabled: Bool
    @Binding var soundEnabled: Bool

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        Button(action: {
            handleItemTap(item.type)
        }) {
            HStack(spacing: spacing.md) {
                // Icon
                Image(systemName: item.icon)
                    .font(.system(size: 18))
                    .foregroundColor(iconColor(for: item.type))
                    .frame(width: 24, height: 24)

                // Title
                Text(item.title)
                    .font(.system(size: 16))
                    .foregroundColor(colors.textPrimary)

                Spacer()

                // Trailing element (toggle or chevron)
                trailingElement(for: item.type)
            }
            .padding(.horizontal, spacing.md)
            .padding(.vertical, spacing.md)
        }
        .buttonStyle(PlainButtonStyle())

        if !isLast {
            Divider()
                .padding(.leading, 56)
        }
    }

    @ViewBuilder
    private func trailingElement(for type: MoreItemType) -> some View {
        switch type {
        case .notifications:
            Toggle("", isOn: $notificationsEnabled)
                .toggleStyle(SwitchToggleStyle(tint: colors.primary))
        case .appearance:
            Toggle("", isOn: $darkModeEnabled)
                .toggleStyle(SwitchToggleStyle(tint: colors.primary))
        default:
            Image(systemName: "chevron.right")
                .font(.system(size: 14))
                .foregroundColor(colors.textSecondary)
        }
    }

    private func iconColor(for type: MoreItemType) -> Color {
        switch type {
        case .report:
            return colors.warning
        case .security, .privacy:
            return colors.success
        case .billing:
            return colors.primary
        default:
            return colors.textSecondary
        }
    }

    private func handleItemTap(_ type: MoreItemType) {
        // Handle different item types
        switch type {
        case .profile:
            print("Open profile")
        case .privacy:
            print("Open privacy settings")
        case .security:
            print("Open security settings")
        case .billing:
            print("Open billing")
        case .language:
            print("Open language settings")
        case .storage:
            print("Open storage settings")
        case .help:
            print("Open help center")
        case .contact:
            print("Open contact")
        case .report:
            print("Open report issue")
        case .about:
            print("Open about")
        default:
            break
        }
    }
}

// MARK: - Profile View (Placeholder)
private struct ProfileView: View {
    @Environment(\.presentationMode) var presentationMode

    private let colors = Colors()

    var body: some View {
        NavigationView {
            VStack {
                Text("Profile Details")
                    .font(.title2)
                    .foregroundColor(colors.textPrimary)

                Spacer()

                Button("Save") {
                    presentationMode.wrappedValue.dismiss()
                }
                .foregroundColor(colors.primary)
            }
            .navigationTitle("Profile")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        presentationMode.wrappedValue.dismiss()
                    }
                }
            }
        }
    }
}

// MARK: - Preview
#if DEBUG
struct MoreScreen_Previews: PreviewProvider {
    static var previews: some View {
        MoreScreen()
            .previewDisplayName("More Screen")
    }
}
#endif