import SwiftUI
import Foundation

/**
 * Stream Content Display Views for iOS
 * SwiftUI implementation supporting native iOS scroll and navigation patterns
 */

// MARK: - Stream Content Models (iOS-specific)

enum StreamContentType: String, CaseIterable, Codable {
    case book = "BOOK"
    case podcast = "PODCAST"
    case cartoon = "CARTOON"
    case shortMovie = "SHORT_MOVIE"
    case longMovie = "LONG_MOVIE"
    case music = "MUSIC"
    case art = "ART"
}

enum StreamAvailabilityStatus: String, Codable {
    case available = "AVAILABLE"
    case comingSoon = "COMING_SOON"
    case unavailable = "UNAVAILABLE"
}

struct StreamContentItem: Identifiable, Codable, Hashable {
    let id: String
    let categoryId: String
    let title: String
    let description: String
    let thumbnailUrl: String
    let contentType: StreamContentType
    let duration: Int? // in seconds
    let price: Double
    let currency: String
    let availabilityStatus: StreamAvailabilityStatus
    let isFeatured: Bool
    let featuredOrder: Int?
    let metadata: [String: String]
    let createdAt: String
    let updatedAt: String

    var isAvailable: Bool {
        availabilityStatus == .available
    }

    var canPurchase: Bool {
        isAvailable
    }

    var durationString: String {
        guard let duration = duration else { return "" }

        let hours = duration / 3600
        let minutes = (duration % 3600) / 60
        let seconds = duration % 60

        if hours > 0 {
            return String(format: "%d:%02d:%02d", hours, minutes, seconds)
        } else {
            return String(format: "%d:%02d", minutes, seconds)
        }
    }

    var formattedPrice: String {
        "\(currency) \(String(format: "%.2f", price))"
    }
}

// MARK: - Main Stream Content View

struct StreamContentView: View {
    let content: [StreamContentItem]
    let isLoading: Bool
    let displayMode: ContentDisplayMode

    let onContentTap: (StreamContentItem) -> Void
    let onAddToCart: (StreamContentItem) -> Void
    let onPurchase: (StreamContentItem) -> Void

    enum ContentDisplayMode {
        case grid
        case list
        case carousel
    }

    var body: some View {
        Group {
            if isLoading {
                ContentLoadingView(displayMode: displayMode)
            } else if content.isEmpty {
                ContentEmptyView()
            } else {
                switch displayMode {
                case .grid:
                    ContentGridView(
                        content: content,
                        onContentTap: onContentTap,
                        onAddToCart: onAddToCart,
                        onPurchase: onPurchase
                    )
                case .list:
                    ContentListView(
                        content: content,
                        onContentTap: onContentTap,
                        onAddToCart: onAddToCart,
                        onPurchase: onPurchase
                    )
                case .carousel:
                    ContentCarouselView(
                        content: content,
                        onContentTap: onContentTap,
                        onAddToCart: onAddToCart,
                        onPurchase: onPurchase
                    )
                }
            }
        }
    }
}

// MARK: - Grid View Implementation

private struct ContentGridView: View {
    let content: [StreamContentItem]
    let onContentTap: (StreamContentItem) -> Void
    let onAddToCart: (StreamContentItem) -> Void
    let onPurchase: (StreamContentItem) -> Void

    private let columns = [
        GridItem(.flexible(), spacing: 12),
        GridItem(.flexible(), spacing: 12)
    ]

    var body: some View {
        ScrollView {
            LazyVGrid(columns: columns, spacing: 16) {
                ForEach(content) { item in
                    StreamContentCard(
                        content: item,
                        onTap: { onContentTap(item) },
                        onAddToCart: { onAddToCart(item) },
                        onPurchase: { onPurchase(item) }
                    )
                }
            }
            .padding(.horizontal, 16)
            .padding(.top, 8)
        }
    }
}

// MARK: - List View Implementation

private struct ContentListView: View {
    let content: [StreamContentItem]
    let onContentTap: (StreamContentItem) -> Void
    let onAddToCart: (StreamContentItem) -> Void
    let onPurchase: (StreamContentItem) -> Void

    var body: some View {
        ScrollView {
            LazyVStack(spacing: 12) {
                ForEach(content) { item in
                    StreamContentListItem(
                        content: item,
                        onTap: { onContentTap(item) },
                        onAddToCart: { onAddToCart(item) },
                        onPurchase: { onPurchase(item) }
                    )
                }
            }
            .padding(.horizontal, 16)
            .padding(.top, 8)
        }
    }
}

// MARK: - Carousel View Implementation

private struct ContentCarouselView: View {
    let content: [StreamContentItem]
    let onContentTap: (StreamContentItem) -> Void
    let onAddToCart: (StreamContentItem) -> Void
    let onPurchase: (StreamContentItem) -> Void

    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 16) {
                ForEach(content) { item in
                    FeaturedContentCard(
                        content: item,
                        onTap: { onContentTap(item) },
                        onAddToCart: { onAddToCart(item) }
                    )
                }
            }
            .padding(.horizontal, 16)
        }
    }
}

// MARK: - Stream Content Card

private struct StreamContentCard: View {
    let content: StreamContentItem
    let onTap: () -> Void
    let onAddToCart: () -> Void
    let onPurchase: () -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Thumbnail
            AsyncImage(url: URL(string: content.thumbnailUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(16/9, contentMode: .fill)
            } placeholder: {
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color(.systemGray5))
                    .aspectRatio(16/9, contentMode: .fit)
                    .overlay(
                        Image(systemName: "photo")
                            .foregroundColor(.secondary)
                    )
            }
            .clipShape(RoundedRectangle(cornerRadius: 12))

            VStack(alignment: .leading, spacing: 8) {
                // Title
                Text(content.title)
                    .font(.system(size: 16, weight: .semibold))
                    .lineLimit(2)
                    .multilineTextAlignment(.leading)

                // Description
                Text(content.description)
                    .font(.system(size: 13))
                    .foregroundColor(.secondary)
                    .lineLimit(2)

                // Metadata
                HStack {
                    HStack(spacing: 4) {
                        Image(systemName: getContentTypeIcon(content.contentType))
                            .font(.system(size: 12))
                            .foregroundColor(.accentColor)

                        if !content.durationString.isEmpty {
                            Text(content.durationString)
                                .font(.system(size: 11))
                                .foregroundColor(.secondary)
                        }
                    }

                    Spacer()

                    Text(content.formattedPrice)
                        .font(.system(size: 14, weight: .bold))
                        .foregroundColor(.accentColor)
                }

                // Action Buttons
                if content.canPurchase {
                    HStack(spacing: 8) {
                        Button(action: onAddToCart) {
                            HStack(spacing: 4) {
                                Image(systemName: "cart.badge.plus")
                                    .font(.system(size: 12, weight: .medium))
                                Text("Add")
                                    .font(.system(size: 12, weight: .medium))
                            }
                            .foregroundColor(.accentColor)
                            .padding(.horizontal, 12)
                            .padding(.vertical, 6)
                            .background(
                                RoundedRectangle(cornerRadius: 8)
                                    .stroke(Color.accentColor, lineWidth: 1)
                            )
                        }

                        Button(action: onPurchase) {
                            Text("Buy")
                                .font(.system(size: 12, weight: .semibold))
                                .foregroundColor(.white)
                                .padding(.horizontal, 16)
                                .padding(.vertical, 6)
                                .background(
                                    RoundedRectangle(cornerRadius: 8)
                                        .fill(Color.accentColor)
                                )
                        }
                    }
                } else {
                    Text(content.availabilityStatus == .comingSoon ? "Coming Soon" : "Unavailable")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.white)
                        .padding(.horizontal, 12)
                        .padding(.vertical, 6)
                        .background(
                            RoundedRectangle(cornerRadius: 8)
                                .fill(Color(.systemGray))
                        )
                }
            }
            .padding(.top, 12)
        }
        .onTapGesture {
            onTap()
        }
    }
}

// MARK: - Stream Content List Item

private struct StreamContentListItem: View {
    let content: StreamContentItem
    let onTap: () -> Void
    let onAddToCart: () -> Void
    let onPurchase: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            // Thumbnail
            AsyncImage(url: URL(string: content.thumbnailUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color(.systemGray5))
                    .overlay(
                        Image(systemName: "photo")
                            .foregroundColor(.secondary)
                    )
            }
            .frame(width: 80, height: 80)
            .clipShape(RoundedRectangle(cornerRadius: 8))

            // Content Info
            VStack(alignment: .leading, spacing: 6) {
                Text(content.title)
                    .font(.system(size: 16, weight: .semibold))
                    .lineLimit(1)

                Text(content.description)
                    .font(.system(size: 13))
                    .foregroundColor(.secondary)
                    .lineLimit(2)

                HStack {
                    Image(systemName: getContentTypeIcon(content.contentType))
                        .font(.system(size: 12))
                        .foregroundColor(.accentColor)

                    if !content.durationString.isEmpty {
                        Text(content.durationString)
                            .font(.system(size: 11))
                            .foregroundColor(.secondary)
                    }

                    Spacer()

                    Text(content.formattedPrice)
                        .font(.system(size: 14, weight: .bold))
                        .foregroundColor(.accentColor)
                }
            }

            // Action Buttons
            VStack(spacing: 4) {
                if content.canPurchase {
                    Button(action: onAddToCart) {
                        Image(systemName: "cart.badge.plus")
                            .font(.system(size: 16))
                            .foregroundColor(.accentColor)
                    }

                    Button(action: onPurchase) {
                        Text("Buy")
                            .font(.system(size: 11, weight: .semibold))
                            .foregroundColor(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 4)
                            .background(
                                RoundedRectangle(cornerRadius: 6)
                                    .fill(Color.accentColor)
                            )
                    }
                }
            }
        }
        .padding(.vertical, 8)
        .onTapGesture {
            onTap()
        }
    }
}

// MARK: - Featured Content Card

private struct FeaturedContentCard: View {
    let content: StreamContentItem
    let onTap: () -> Void
    let onAddToCart: () -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            AsyncImage(url: URL(string: content.thumbnailUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(16/9, contentMode: .fill)
            } placeholder: {
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color(.systemGray5))
                    .aspectRatio(16/9, contentMode: .fit)
                    .overlay(
                        Image(systemName: "photo")
                            .foregroundColor(.secondary)
                    )
            }
            .clipShape(RoundedRectangle(cornerRadius: 12))

            VStack(alignment: .leading, spacing: 8) {
                Text(content.title)
                    .font(.system(size: 14, weight: .semibold))
                    .lineLimit(2)

                HStack {
                    Text(content.formattedPrice)
                        .font(.system(size: 12, weight: .bold))
                        .foregroundColor(.accentColor)

                    Spacer()

                    Button(action: onAddToCart) {
                        Image(systemName: "plus.circle.fill")
                            .font(.system(size: 20))
                            .foregroundColor(.accentColor)
                    }
                }
            }
            .padding(.top, 8)
        }
        .frame(width: 160)
        .onTapGesture {
            onTap()
        }
    }
}

// MARK: - Loading States

private struct ContentLoadingView: View {
    let displayMode: StreamContentView.ContentDisplayMode

    var body: some View {
        Group {
            switch displayMode {
            case .grid:
                GridLoadingView()
            case .list:
                ListLoadingView()
            case .carousel:
                CarouselLoadingView()
            }
        }
    }
}

private struct GridLoadingView: View {
    private let columns = [
        GridItem(.flexible(), spacing: 12),
        GridItem(.flexible(), spacing: 12)
    ]

    var body: some View {
        ScrollView {
            LazyVGrid(columns: columns, spacing: 16) {
                ForEach(0..<6, id: \.self) { _ in
                    ContentCardPlaceholder()
                }
            }
            .padding(.horizontal, 16)
        }
    }
}

private struct ListLoadingView: View {
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 12) {
                ForEach(0..<8, id: \.self) { _ in
                    ContentListPlaceholder()
                }
            }
            .padding(.horizontal, 16)
        }
    }
}

private struct CarouselLoadingView: View {
    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 16) {
                ForEach(0..<5, id: \.self) { _ in
                    ContentCardPlaceholder()
                        .frame(width: 160)
                }
            }
            .padding(.horizontal, 16)
        }
    }
}

private struct ContentCardPlaceholder: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            RoundedRectangle(cornerRadius: 12)
                .fill(Color(.systemGray5))
                .aspectRatio(16/9, contentMode: .fit)

            VStack(alignment: .leading, spacing: 6) {
                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 16)

                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 12)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.trailing, 40)

                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 12)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.trailing, 80)
            }
        }
        .redacted(reason: .placeholder)
    }
}

private struct ContentListPlaceholder: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color(.systemGray5))
                .frame(width: 80, height: 80)

            VStack(alignment: .leading, spacing: 6) {
                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 16)

                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 12)
                    .padding(.trailing, 60)

                RoundedRectangle(cornerRadius: 4)
                    .fill(Color(.systemGray5))
                    .frame(height: 12)
                    .padding(.trailing, 100)
            }

            Spacer()
        }
        .redacted(reason: .placeholder)
    }
}

// MARK: - Empty State

private struct ContentEmptyView: View {
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "tray")
                .font(.system(size: 48))
                .foregroundColor(.secondary)

            Text("No content available")
                .font(.title2)
                .foregroundColor(.secondary)

            Text("Check back later for new content")
                .font(.body)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
        }
        .padding(40)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}

// MARK: - Helper Functions

private func getContentTypeIcon(_ contentType: StreamContentType) -> String {
    switch contentType {
    case .book:
        return "book.fill"
    case .podcast:
        return "podcast.fill"
    case .cartoon:
        return "tv.fill"
    case .shortMovie, .longMovie:
        return "play.rectangle.fill"
    case .music:
        return "music.note"
    case .art:
        return "paintbrush.fill"
    }
}

// MARK: - Featured Content Section

struct FeaturedStreamContentSection: View {
    let featuredContent: [StreamContentItem]
    let isLoading: Bool
    let onContentTap: (StreamContentItem) -> Void
    let onAddToCart: (StreamContentItem) -> Void
    let onSeeAllTap: () -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Featured Content")
                    .font(.title2)
                    .fontWeight(.bold)

                Spacer()

                Button("See All", action: onSeeAllTap)
                    .font(.body)
                    .foregroundColor(.accentColor)
            }
            .padding(.horizontal, 16)

            if isLoading {
                CarouselLoadingView()
            } else {
                ContentCarouselView(
                    content: featuredContent,
                    onContentTap: onContentTap,
                    onAddToCart: onAddToCart,
                    onPurchase: { _ in }
                )
            }
        }
    }
}

// MARK: - Preview Support

#if DEBUG
struct StreamContentView_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView {
            StreamContentView(
                content: sampleContent,
                isLoading: false,
                displayMode: .grid,
                onContentTap: { _ in },
                onAddToCart: { _ in },
                onPurchase: { _ in }
            )
            .navigationTitle("Stream Content")
        }
        .previewDisplayName("Grid View")

        NavigationView {
            StreamContentView(
                content: sampleContent,
                isLoading: false,
                displayMode: .list,
                onContentTap: { _ in },
                onAddToCart: { _ in },
                onPurchase: { _ in }
            )
            .navigationTitle("Stream Content")
        }
        .previewDisplayName("List View")

        FeaturedStreamContentSection(
            featuredContent: Array(sampleContent.prefix(3)),
            isLoading: false,
            onContentTap: { _ in },
            onAddToCart: { _ in },
            onSeeAllTap: { }
        )
        .previewDisplayName("Featured Section")
    }

    static let sampleContent: [StreamContentItem] = [
        StreamContentItem(
            id: "1",
            categoryId: "books",
            title: "The Art of iOS Development",
            description: "A comprehensive guide to building beautiful iOS applications",
            thumbnailUrl: "https://example.com/book1.jpg",
            contentType: .book,
            duration: nil,
            price: 29.99,
            currency: "USD",
            availabilityStatus: .available,
            isFeatured: true,
            featuredOrder: 1,
            metadata: [:],
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        ),
        StreamContentItem(
            id: "2",
            categoryId: "podcasts",
            title: "Tech Talk Weekly",
            description: "Weekly discussions about the latest in technology",
            thumbnailUrl: "https://example.com/podcast1.jpg",
            contentType: .podcast,
            duration: 3600,
            price: 9.99,
            currency: "USD",
            availabilityStatus: .available,
            isFeatured: false,
            featuredOrder: nil,
            metadata: [:],
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        )
    ]
}
#endif