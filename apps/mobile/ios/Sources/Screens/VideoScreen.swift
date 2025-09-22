//
//  VideoScreen.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Video conferencing and streaming interface screen
public struct VideoScreen: View {
    @State private var selectedTab: VideoTab = .meetings
    @State private var showingCreateMeeting = false
    @State private var isCameraOn = false
    @State private var isMicOn = true

    private let colors = Colors()
    private let spacing = Spacing()

    // Mock meetings
    private let upcomingMeetings = [
        ("Team Standup", "09:00 AM", "with John, Sarah, Mike", "video.fill"),
        ("Client Review", "02:00 PM", "with ABC Company", "person.2.fill"),
        ("Project Planning", "04:30 PM", "with Development Team", "calendar"),
    ]

    // Mock recent recordings
    private let recordings = [
        ("Q3 Review Meeting", "45 min", "Yesterday", "play.rectangle.fill"),
        ("Product Demo", "28 min", "2 days ago", "play.rectangle.fill"),
        ("Training Session", "1h 12min", "1 week ago", "play.rectangle.fill"),
    ]

    public init() {}

    public var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Segmented control
                HStack(spacing: 0) {
                    ForEach(VideoTab.allCases, id: \.self) { tab in
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
                case .meetings:
                    MeetingsView(meetings: upcomingMeetings)
                case .recordings:
                    RecordingsView(recordings: recordings)
                case .live:
                    LiveView(isCameraOn: $isCameraOn, isMicOn: $isMicOn)
                }
            }
            .navigationTitle("Video")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {
                        showingCreateMeeting = true
                    }) {
                        Image(systemName: "plus.circle.fill")
                            .foregroundColor(colors.primary)
                    }
                }
            }
        }
        .navigationViewStyle(StackNavigationViewStyle())
        .sheet(isPresented: $showingCreateMeeting) {
            CreateMeetingView()
        }
    }
}

// MARK: - Video Tab Enum
enum VideoTab: CaseIterable {
    case meetings, recordings, live

    var title: String {
        switch self {
        case .meetings: return "Meetings"
        case .recordings: return "Recordings"
        case .live: return "Live"
        }
    }

    var icon: String {
        switch self {
        case .meetings: return "calendar"
        case .recordings: return "play.rectangle.fill"
        case .live: return "video.fill"
        }
    }
}

// MARK: - Meetings View
private struct MeetingsView: View {
    let meetings: [(String, String, String, String)]

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        ScrollView {
            VStack(spacing: spacing.md) {
                // Quick join section
                VStack(alignment: .leading, spacing: spacing.sm) {
                    Text("Quick Join")
                        .font(.system(size: 18, weight: .bold))
                        .foregroundColor(colors.textPrimary)
                        .frame(maxWidth: .infinity, alignment: .leading)

                    HStack {
                        TextField("Enter meeting ID", text: .constant(""))
                            .textFieldStyle(PlainTextFieldStyle())
                            .padding(.horizontal, spacing.md)
                            .padding(.vertical, spacing.sm)
                            .background(colors.surface)
                            .cornerRadius(12)

                        Button(action: {}) {
                            Text("Join")
                                .font(.system(size: 14, weight: .semibold))
                                .foregroundColor(colors.textOnPrimary)
                                .padding(.horizontal, spacing.lg)
                                .padding(.vertical, spacing.sm)
                                .background(colors.primary)
                                .cornerRadius(12)
                        }
                    }
                }
                .padding(.horizontal, spacing.md)

                // Upcoming meetings
                VStack(alignment: .leading, spacing: spacing.sm) {
                    Text("Upcoming Meetings")
                        .font(.system(size: 18, weight: .bold))
                        .foregroundColor(colors.textPrimary)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .padding(.horizontal, spacing.md)

                    ForEach(Array(meetings.enumerated()), id: \.offset) { index, meeting in
                        MeetingCard(
                            title: meeting.0,
                            time: meeting.1,
                            participants: meeting.2,
                            icon: meeting.3
                        )
                    }
                }
            }
            .padding(.top, spacing.sm)
        }
    }
}

// MARK: - Meeting Card Component
private struct MeetingCard: View {
    let title: String
    let time: String
    let participants: String
    let icon: String

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        HStack(spacing: spacing.md) {
            // Meeting icon
            Image(systemName: icon)
                .font(.system(size: 24))
                .foregroundColor(colors.primary)
                .frame(width: 48, height: 48)
                .background(colors.primary.opacity(0.1))
                .cornerRadius(12)

            // Meeting info
            VStack(alignment: .leading, spacing: spacing.xs) {
                Text(title)
                    .font(.system(size: 16, weight: .semibold))
                    .foregroundColor(colors.textPrimary)

                Text(time)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(colors.primary)

                Text(participants)
                    .font(.system(size: 12))
                    .foregroundColor(colors.textSecondary)
            }

            Spacer()

            // Join button
            Button(action: {}) {
                Text("Join")
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundColor(colors.textOnPrimary)
                    .padding(.horizontal, spacing.md)
                    .padding(.vertical, spacing.xs)
                    .background(colors.primary)
                    .cornerRadius(8)
            }
        }
        .padding(spacing.md)
        .background(Color.white)
        .cornerRadius(12)
        .shadow(color: colors.shadowLight, radius: 4, y: 2)
        .padding(.horizontal, spacing.md)
    }
}

// MARK: - Recordings View
private struct RecordingsView: View {
    let recordings: [(String, String, String, String)]

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        ScrollView {
            VStack(spacing: spacing.sm) {
                ForEach(Array(recordings.enumerated()), id: \.offset) { index, recording in
                    RecordingCard(
                        title: recording.0,
                        duration: recording.1,
                        date: recording.2,
                        icon: recording.3
                    )
                }
            }
            .padding(.horizontal, spacing.md)
            .padding(.top, spacing.sm)
        }
    }
}

// MARK: - Recording Card Component
private struct RecordingCard: View {
    let title: String
    let duration: String
    let date: String
    let icon: String

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        HStack(spacing: spacing.md) {
            // Thumbnail placeholder
            Image(systemName: icon)
                .font(.system(size: 24))
                .foregroundColor(colors.primary)
                .frame(width: 60, height: 45)
                .background(colors.surface)
                .cornerRadius(8)

            // Recording info
            VStack(alignment: .leading, spacing: spacing.xs) {
                Text(title)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(colors.textPrimary)
                    .lineLimit(1)

                HStack {
                    Text(duration)
                        .font(.system(size: 12))
                        .foregroundColor(colors.textSecondary)

                    Text("•")
                        .font(.system(size: 12))
                        .foregroundColor(colors.textSecondary)

                    Text(date)
                        .font(.system(size: 12))
                        .foregroundColor(colors.textSecondary)
                }
            }

            Spacer()

            // More options
            Button(action: {}) {
                Image(systemName: "ellipsis")
                    .foregroundColor(colors.textSecondary)
            }
        }
        .padding(.horizontal, spacing.md)
        .padding(.vertical, spacing.sm)
        .background(Color.white)
        .cornerRadius(12)
        .shadow(color: colors.shadowLight, radius: 2, y: 1)
    }
}

// MARK: - Live View
private struct LiveView: View {
    @Binding var isCameraOn: Bool
    @Binding var isMicOn: Bool

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        VStack(spacing: spacing.lg) {
            Spacer()

            // Video preview
            RoundedRectangle(cornerRadius: 16)
                .fill(colors.backgroundSecondary)
                .frame(height: 300)
                .overlay(
                    VStack {
                        if isCameraOn {
                            Image(systemName: "video.fill")
                                .font(.system(size: 48))
                                .foregroundColor(colors.primary)
                            Text("Camera Active")
                                .font(.system(size: 16, weight: .medium))
                                .foregroundColor(colors.textPrimary)
                        } else {
                            Image(systemName: "video.slash.fill")
                                .font(.system(size: 48))
                                .foregroundColor(colors.textSecondary)
                            Text("Camera Off")
                                .font(.system(size: 16, weight: .medium))
                                .foregroundColor(colors.textSecondary)
                        }
                    }
                )
                .padding(.horizontal, spacing.md)

            Spacer()

            // Control buttons
            HStack(spacing: spacing.xl) {
                // Mic button
                Button(action: {
                    isMicOn.toggle()
                }) {
                    Image(systemName: isMicOn ? "mic.fill" : "mic.slash.fill")
                        .font(.system(size: 24))
                        .foregroundColor(isMicOn ? colors.textOnPrimary : colors.error)
                        .frame(width: 56, height: 56)
                        .background(isMicOn ? colors.success : colors.surface)
                        .cornerRadius(28)
                }

                // Camera button
                Button(action: {
                    isCameraOn.toggle()
                }) {
                    Image(systemName: isCameraOn ? "video.fill" : "video.slash.fill")
                        .font(.system(size: 24))
                        .foregroundColor(isCameraOn ? colors.textOnPrimary : colors.error)
                        .frame(width: 56, height: 56)
                        .background(isCameraOn ? colors.success : colors.surface)
                        .cornerRadius(28)
                }

                // End call button
                Button(action: {}) {
                    Image(systemName: "phone.down.fill")
                        .font(.system(size: 24))
                        .foregroundColor(.white)
                        .frame(width: 56, height: 56)
                        .background(colors.error)
                        .cornerRadius(28)
                }
            }
            .padding(.bottom, spacing.xl)
        }
        .background(colors.background)
    }
}

// MARK: - Create Meeting View (Placeholder)
private struct CreateMeetingView: View {
    @Environment(\.presentationMode) var presentationMode

    private let colors = Colors()

    var body: some View {
        NavigationView {
            VStack {
                Text("Create New Meeting")
                    .font(.title2)
                    .foregroundColor(colors.textPrimary)

                Spacer()

                Button("Create") {
                    presentationMode.wrappedValue.dismiss()
                }
                .foregroundColor(colors.primary)
            }
            .navigationTitle("New Meeting")
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
struct VideoScreen_Previews: PreviewProvider {
    static var previews: some View {
        VideoScreen()
            .previewDisplayName("Video Screen")
    }
}
#endif