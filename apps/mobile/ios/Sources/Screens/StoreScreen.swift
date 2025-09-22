//
//  StoreScreen.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// E-commerce store interface screen
public struct StoreScreen: View {
    @State private var searchText = ""
    @State private var selectedCategory: String = "All"
    @State private var cartItems: Int = 3

    private let colors = Colors()
    private let spacing = Spacing()

    // Mock categories
    private let categories = ["All", "Electronics", "Fashion", "Home", "Books", "Sports"]

    // Mock products
    private let products = [
        ("iPhone 15 Pro", "$999", "electronics", "phone.fill"),
        ("MacBook Air", "$1199", "electronics", "laptopcomputer"),
        ("AirPods Pro", "$249", "electronics", "airpods"),
        ("Nike Sneakers", "$129", "fashion", "figure.walk"),
        ("Coffee Maker", "$89", "home", "cup.and.saucer.fill"),
        ("Wireless Mouse", "$59", "electronics", "computermouse.fill"),
        ("Yoga Mat", "$39", "sports", "figure.yoga"),
        ("Cookbook", "$24", "books", "book.fill")
    ]

    public init() {}

    public var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search bar
                HStack {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(colors.textSecondary)

                    TextField("Search products", text: $searchText)
                        .textFieldStyle(PlainTextFieldStyle())

                    Button(action: {
                        // Voice search
                    }) {
                        Image(systemName: "mic.fill")
                            .foregroundColor(colors.primary)
                    }
                }
                .padding(.horizontal, spacing.md)
                .padding(.vertical, spacing.sm)
                .background(colors.surface)
                .cornerRadius(12)
                .padding(.horizontal, spacing.md)
                .padding(.top, spacing.sm)

                // Category selector
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: spacing.sm) {
                        ForEach(categories, id: \.self) { category in
                            CategoryChip(
                                title: category,
                                isSelected: selectedCategory == category
                            ) {
                                selectedCategory = category
                            }
                        }
                    }
                    .padding(.horizontal, spacing.md)
                }
                .padding(.vertical, spacing.sm)

                // Products grid
                ScrollView {
                    LazyVGrid(columns: [
                        GridItem(.flexible()),
                        GridItem(.flexible())
                    ], spacing: spacing.md) {
                        ForEach(Array(filteredProducts.enumerated()), id: \.offset) { index, product in
                            ProductCard(
                                name: product.0,
                                price: product.1,
                                category: product.2,
                                icon: product.3
                            )
                        }
                    }
                    .padding(.horizontal, spacing.md)
                }
            }
            .navigationTitle("Store")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {
                        // Open cart
                    }) {
                        ZStack {
                            Image(systemName: "cart.fill")
                                .foregroundColor(colors.primary)

                            if cartItems > 0 {
                                Circle()
                                    .fill(colors.error)
                                    .frame(width: 16, height: 16)
                                    .overlay(
                                        Text("\(cartItems)")
                                            .font(.system(size: 10, weight: .bold))
                                            .foregroundColor(.white)
                                    )
                                    .offset(x: 10, y: -10)
                            }
                        }
                    }
                }
            }
        }
        .navigationViewStyle(StackNavigationViewStyle())
    }

    private var filteredProducts: [(String, String, String, String)] {
        if selectedCategory == "All" {
            return products
        }
        return products.filter { $0.2 == selectedCategory.lowercased() }
    }
}

// MARK: - Category Chip Component
private struct CategoryChip: View {
    let title: String
    let isSelected: Bool
    let action: () -> Void

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        Button(action: action) {
            Text(title)
                .font(.system(size: 14, weight: isSelected ? .semibold : .medium))
                .foregroundColor(isSelected ? colors.textOnPrimary : colors.textSecondary)
                .padding(.horizontal, spacing.md)
                .padding(.vertical, spacing.xs)
                .background(isSelected ? colors.primary : colors.surface)
                .cornerRadius(20)
        }
    }
}

// MARK: - Product Card Component
private struct ProductCard: View {
    let name: String
    let price: String
    let category: String
    let icon: String

    private let colors = Colors()
    private let spacing = Spacing()

    var body: some View {
        VStack(alignment: .leading, spacing: spacing.sm) {
            // Product icon
            Image(systemName: icon)
                .font(.system(size: 32))
                .foregroundColor(colors.primary)
                .frame(maxWidth: .infinity, alignment: .center)
                .padding(.vertical, spacing.lg)
                .background(colors.surface)
                .cornerRadius(12)

            // Product info
            VStack(alignment: .leading, spacing: spacing.xs) {
                Text(name)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(colors.textPrimary)
                    .lineLimit(2)

                Text(price)
                    .font(.system(size: 16, weight: .bold))
                    .foregroundColor(colors.primary)
            }

            // Add to cart button
            Button(action: {
                // Add to cart
            }) {
                HStack {
                    Image(systemName: "plus")
                    Text("Add")
                }
                .font(.system(size: 12, weight: .semibold))
                .foregroundColor(colors.textOnPrimary)
                .frame(maxWidth: .infinity)
                .padding(.vertical, spacing.xs)
                .background(colors.primary)
                .cornerRadius(8)
            }
        }
        .padding(spacing.sm)
        .background(Color.white)
        .cornerRadius(12)
        .shadow(color: colors.shadowLight, radius: 4, y: 2)
    }
}

// MARK: - Preview
#if DEBUG
struct StoreScreen_Previews: PreviewProvider {
    static var previews: some View {
        StoreScreen()
            .previewDisplayName("Store Screen")
    }
}
#endif