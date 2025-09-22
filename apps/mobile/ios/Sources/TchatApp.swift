//
//  TchatApp.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI

/// Main application entry point
@main
struct TchatApp: App {

    // MARK: - Properties
    @StateObject private var appState = AppState.shared

    // MARK: - Body
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .onAppear {
                    configureApp()
                }
                .onReceive(NotificationCenter.default.publisher(for: UIApplication.willTerminateNotification)) { _ in
                    // Save state before app terminates
                    handleAppTermination()
                }
        }
    }

    // MARK: - App Configuration
    private func configureApp() {
        // Configure design system
        configureDesignSystem()

        // Load initial data
        loadInitialData()

        // Setup analytics if needed
        setupAnalytics()
    }

    private func configureDesignSystem() {
        // Apply global theme based on user preferences
        let preferences = appState.themePreferences

        // Configure navigation bar appearance
        let navigationBarAppearance = UINavigationBarAppearance()
        navigationBarAppearance.configureWithOpaqueBackground()
        navigationBarAppearance.backgroundColor = UIColor(Colors().navigationBackground)
        navigationBarAppearance.titleTextAttributes = [
            .foregroundColor: UIColor(Colors().textPrimary),
            .font: UIFont.systemFont(ofSize: 18, weight: .semibold)
        ]

        UINavigationBar.appearance().standardAppearance = navigationBarAppearance
        UINavigationBar.appearance().scrollEdgeAppearance = navigationBarAppearance
    }

    private func loadInitialData() {
        // Load any necessary initial data
        print("Loading initial app data...")
    }

    private func setupAnalytics() {
        // Setup analytics tracking if needed
        print("Setting up analytics...")
    }

    private func handleAppTermination() {
        // Perform cleanup tasks
        print("App is terminating, saving state...")
    }
}

/// Main content view with conditional navigation
struct ContentView: View {
    @EnvironmentObject private var appState: AppState

    var body: some View {
        Group {
            if appState.isAuthenticated {
                // Main app interface
                TabNavigationView()
            } else {
                // Authentication flow
                AuthenticationView()
            }
        }
        .animation(.easeInOut(duration: 0.3), value: appState.isAuthenticated)
    }
}

/// Authentication view (placeholder)
struct AuthenticationView: View {
    @EnvironmentObject private var appState: AppState
    @State private var email = ""
    @State private var password = ""

    private let colors = Colors()

    var body: some View {
        NavigationView {
            VStack(spacing: Spacing.lg) {
                // App logo/branding
                VStack(spacing: Spacing.md) {
                    Image(systemName: "message.circle.fill")
                        .font(.system(size: 80))
                        .foregroundColor(colors.primary)

                    Text("Tchat")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                        .foregroundColor(colors.textPrimary)

                    Text("Connect with the world")
                        .font(.subheadline)
                        .foregroundColor(colors.textSecondary)
                }

                Spacer()

                // Login form
                VStack(spacing: Spacing.md) {
                    TchatInput(
                        text: $email,
                        placeholder: "Email",
                        type: .email,
                        leadingIcon: "envelope"
                    )

                    TchatInput(
                        text: $password,
                        placeholder: "Password",
                        type: .password,
                        leadingIcon: "lock"
                    )

                    TchatButton(
                        "Sign In",
                        variant: .primary,
                        size: .large
                    ) {
                        authenticateUser()
                    }
                    .padding(.top, Spacing.sm)
                }

                Spacer()

                // Demo login for testing
                VStack(spacing: Spacing.sm) {
                    Text("Demo Mode")
                        .font(.caption)
                        .foregroundColor(colors.textTertiary)

                    TchatButton(
                        "Continue as Demo User",
                        variant: .outline,
                        size: .medium
                    ) {
                        authenticateAsDemoUser()
                    }
                }
            }
            .padding(Spacing.lg)
            .background(colors.background)
            .navigationBarHidden(true)
        }
    }

    private func authenticateUser() {
        // Simulate authentication
        let user = UserModel(
            id: "user_123",
            username: email.components(separatedBy: "@").first ?? "user",
            email: email,
            displayName: "User Name"
        )

        appState.updateUser(user)
    }

    private func authenticateAsDemoUser() {
        let demoUser = UserModel(
            id: "demo_user",
            username: "demo",
            email: "demo@tchat.app",
            displayName: "Demo User"
        )

        appState.updateUser(demoUser)
    }
}

// MARK: - Preview
#if DEBUG
struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
            .environmentObject(AppState.shared)
    }
}

struct AuthenticationView_Previews: PreviewProvider {
    static var previews: some View {
        AuthenticationView()
            .environmentObject(AppState.shared)
    }
}
#endif