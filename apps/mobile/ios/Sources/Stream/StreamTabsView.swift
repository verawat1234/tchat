import SwiftUI
import Foundation

/**
 * Stream Tab Navigation for iOS
 * SwiftUI implementation following iOS design guidelines
 */

// MARK: - Stream Models (iOS-specific)

struct StreamCategory: Identifiable, Codable, Hashable {
    let id: String
    let name: String
    let displayOrder: Int
    let iconName: String
    let isActive: Bool
    let subtabs: [StreamSubtab]?
    let featuredContentEnabled: Bool
    let createdAt: String
    let updatedAt: String
}

struct StreamSubtab: Identifiable, Codable, Hashable {
    let id: String
    let categoryId: String
    let name: String
    let displayOrder: Int
    let filterCriteria: [String: String]
    let isActive: Bool
    let createdAt: String
    let updatedAt: String
}

// MARK: - Main Stream Tabs View

struct StreamTabsView: View {
    @Binding var selectedCategoryId: String
    @Binding var selectedSubtabId: String?

    let categories: [StreamCategory]
    let isLoading: Bool

    let onCategorySelected: (String) -> Void
    let onSubtabSelected: (String?) -> Void

    init(
        categories: [StreamCategory],
        selectedCategoryId: Binding<String>,
        selectedSubtabId: Binding<String?>,
        isLoading: Bool = false,
        onCategorySelected: @escaping (String) -> Void,
        onSubtabSelected: @escaping (String?) -> Void
    ) {
        self.categories = categories
        self._selectedCategoryId = selectedCategoryId
        self._selectedSubtabId = selectedSubtabId
        self.isLoading = isLoading
        self.onCategorySelected = onCategorySelected
        self.onSubtabSelected = onSubtabSelected
    }

    var body: some View {
        VStack(spacing: 0) {
            // Main Category Tabs
            if isLoading {
                CategoryTabsLoadingView()
            } else {
                CategoryTabsView(
                    categories: categories,
                    selectedCategoryId: $selectedCategoryId,
                    onCategorySelected: onCategorySelected
                )
            }

            // Subtabs for selected category
            if let selectedCategory = categories.first(where: { $0.id == selectedCategoryId }),
               let subtabs = selectedCategory.subtabs,
               !subtabs.isEmpty {

                Divider()
                    .padding(.vertical, 8)

                SubtabsView(
                    subtabs: subtabs,
                    selectedSubtabId: $selectedSubtabId,
                    onSubtabSelected: onSubtabSelected
                )
            }
        }
        .background(Color(.systemBackground))
    }
}

// MARK: - Category Tabs View

private struct CategoryTabsView: View {
    let categories: [StreamCategory]
    @Binding var selectedCategoryId: String
    let onCategorySelected: (String) -> Void

    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 12) {
                ForEach(categories.filter { $0.isActive }) { category in
                    CategoryTabButton(
                        category: category,
                        isSelected: category.id == selectedCategoryId,
                        onTap: {
                            onCategorySelected(category.id)
                        }
                    )
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }
}

private struct CategoryTabButton: View {
    let category: StreamCategory
    let isSelected: Bool
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            HStack(spacing: 8) {
                Image(systemName: getCategorySystemIcon(category.iconName))
                    .font(.system(size: 16, weight: .medium))

                Text(category.name)
                    .font(.system(size: 14, weight: isSelected ? .semibold : .medium))
            }
            .foregroundColor(isSelected ? .white : .primary)
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .background(
                RoundedRectangle(cornerRadius: 20)
                    .fill(isSelected ? Color.accentColor : Color(.systemGray6))
            )
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

// MARK: - Subtabs View

private struct SubtabsView: View {
    let subtabs: [StreamSubtab]
    @Binding var selectedSubtabId: String?
    let onSubtabSelected: (String?) -> Void

    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 10) {
                // "All" option
                SubtabButton(
                    title: "All",
                    isSelected: selectedSubtabId == nil,
                    onTap: {
                        onSubtabSelected(nil)
                    }
                )

                // Individual subtabs
                ForEach(subtabs.filter { $0.isActive }) { subtab in
                    SubtabButton(
                        title: subtab.name,
                        isSelected: subtab.id == selectedSubtabId,
                        onTap: {
                            onSubtabSelected(subtab.id)
                        }
                    )
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
    }
}

private struct SubtabButton: View {
    let title: String
    let isSelected: Bool
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            Text(title)
                .font(.system(size: 12, weight: isSelected ? .semibold : .medium))
                .foregroundColor(isSelected ? Color.accentColor : .secondary)
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(
                            isSelected ? Color.accentColor : Color(.systemGray4),
                            lineWidth: 1
                        )
                        .background(
                            RoundedRectangle(cornerRadius: 16)
                                .fill(isSelected ? Color.accentColor.opacity(0.1) : Color.clear)
                        )
                )
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

// MARK: - Loading States

private struct CategoryTabsLoadingView: View {
    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 12) {
                ForEach(0..<5, id: \.self) { _ in
                    RoundedRectangle(cornerRadius: 20)
                        .fill(Color(.systemGray5))
                        .frame(width: 100, height: 40)
                        .redacted(reason: .placeholder)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }
}

// MARK: - Helper Functions

private func getCategorySystemIcon(_ iconName: String) -> String {
    switch iconName.lowercased() {
    case "book", "books":
        return "book.fill"
    case "podcast", "podcasts":
        return "podcast.fill"
    case "cartoon", "cartoons":
        return "tv.fill"
    case "movie", "movies":
        return "play.rectangle.fill"
    case "music":
        return "music.note"
    case "art":
        return "paintbrush.fill"
    case "video":
        return "play.circle.fill"
    case "audio":
        return "headphones"
    default:
        return "square.grid.2x2.fill"
    }
}

// MARK: - Custom Button Style

private struct ScaleButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.95 : 1.0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

// MARK: - Stream Tab State Manager

@MainActor
class StreamTabState: ObservableObject {
    @Published var selectedCategoryId: String = ""
    @Published var selectedSubtabId: String? = nil
    @Published var categories: [StreamCategory] = []
    @Published var isLoading: Bool = false

    init(initialCategoryId: String = "") {
        self.selectedCategoryId = initialCategoryId
    }

    func selectCategory(_ categoryId: String) {
        if selectedCategoryId != categoryId {
            selectedCategoryId = categoryId
            selectedSubtabId = nil // Reset subtab when category changes
        }
    }

    func selectSubtab(_ subtabId: String?) {
        selectedSubtabId = subtabId
    }

    func loadCategories(_ categories: [StreamCategory]) {
        self.categories = categories

        // Auto-select first active category if none selected
        if selectedCategoryId.isEmpty,
           let firstCategory = categories.first(where: { $0.isActive }) {
            selectedCategoryId = firstCategory.id
        }
    }

    func reset() {
        selectedCategoryId = ""
        selectedSubtabId = nil
    }
}

// MARK: - Stream Category Filter View

struct StreamCategoryFilterView: View {
    let categories: [StreamCategory]
    @Binding var selectedCategories: Set<String>
    let onCategoryToggle: (String) -> Void

    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 10) {
                ForEach(categories.filter { $0.isActive }) { category in
                    CategoryFilterChip(
                        category: category,
                        isSelected: selectedCategories.contains(category.id),
                        onToggle: {
                            onCategoryToggle(category.id)
                        }
                    )
                }
            }
            .padding(.horizontal, 16)
        }
    }
}

private struct CategoryFilterChip: View {
    let category: StreamCategory
    let isSelected: Bool
    let onToggle: () -> Void

    var body: some View {
        Button(action: onToggle) {
            HStack(spacing: 6) {
                Image(systemName: getCategorySystemIcon(category.iconName))
                    .font(.system(size: 12, weight: .medium))

                Text(category.name)
                    .font(.system(size: 12, weight: .medium))
            }
            .foregroundColor(isSelected ? .white : .primary)
            .padding(.horizontal, 12)
            .padding(.vertical, 6)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(isSelected ? Color.accentColor : Color(.systemGray6))
            )
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

// MARK: - Preview Support

#if DEBUG
struct StreamTabsView_Previews: PreviewProvider {
    @State static var selectedCategoryId = "books"
    @State static var selectedSubtabId: String? = nil

    static var previews: some View {
        VStack {
            StreamTabsView(
                categories: sampleCategories,
                selectedCategoryId: $selectedCategoryId,
                selectedSubtabId: $selectedSubtabId,
                onCategorySelected: { categoryId in
                    selectedCategoryId = categoryId
                },
                onSubtabSelected: { subtabId in
                    selectedSubtabId = subtabId
                }
            )

            Spacer()
        }
        .preferredColorScheme(.light)
        .previewDisplayName("Light Mode")

        VStack {
            StreamTabsView(
                categories: sampleCategories,
                selectedCategoryId: $selectedCategoryId,
                selectedSubtabId: $selectedSubtabId,
                onCategorySelected: { categoryId in
                    selectedCategoryId = categoryId
                },
                onSubtabSelected: { subtabId in
                    selectedSubtabId = subtabId
                }
            )

            Spacer()
        }
        .preferredColorScheme(.dark)
        .previewDisplayName("Dark Mode")
    }

    static let sampleCategories: [StreamCategory] = [
        StreamCategory(
            id: "books",
            name: "Books",
            displayOrder: 1,
            iconName: "books",
            isActive: true,
            subtabs: [],
            featuredContentEnabled: true,
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        ),
        StreamCategory(
            id: "podcasts",
            name: "Podcasts",
            displayOrder: 2,
            iconName: "podcasts",
            isActive: true,
            subtabs: [],
            featuredContentEnabled: true,
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        ),
        StreamCategory(
            id: "movies",
            name: "Movies",
            displayOrder: 3,
            iconName: "movies",
            isActive: true,
            subtabs: [
                StreamSubtab(
                    id: "short-movies",
                    categoryId: "movies",
                    name: "Short Films",
                    displayOrder: 1,
                    filterCriteria: ["maxDuration": "1800"],
                    isActive: true,
                    createdAt: "2024-01-01T00:00:00Z",
                    updatedAt: "2024-01-01T00:00:00Z"
                ),
                StreamSubtab(
                    id: "long-movies",
                    categoryId: "movies",
                    name: "Feature Films",
                    displayOrder: 2,
                    filterCriteria: ["minDuration": "1800"],
                    isActive: true,
                    createdAt: "2024-01-01T00:00:00Z",
                    updatedAt: "2024-01-01T00:00:00Z"
                )
            ],
            featuredContentEnabled: true,
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        ),
        StreamCategory(
            id: "music",
            name: "Music",
            displayOrder: 4,
            iconName: "music",
            isActive: true,
            subtabs: [],
            featuredContentEnabled: true,
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z"
        )
    ]
}
#endif