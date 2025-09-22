//
//  TchatPagination.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI

/// Pagination component following Tchat design system
public struct TchatPagination: View {

    // MARK: - Pagination Types
    public enum PaginationStyle {
        case numbered
        case simple
        case compact
    }

    public enum PaginationSize {
        case small
        case medium
        case large
    }

    // MARK: - Properties
    @Binding private var currentPage: Int
    @State private var inputPage: String = ""
    @State private var showPageInput: Bool = false

    let totalPages: Int
    let style: PaginationStyle
    let size: PaginationSize
    let showPageSize: Bool
    let showInfo: Bool
    let showJumpToPage: Bool
    let maxVisiblePages: Int
    let pageSize: Int
    let totalItems: Int
    let onPageChange: ((Int) -> Void)?
    let onPageSizeChange: ((Int) -> Void)?

    // MARK: - Private Properties
    private let colors = Colors()
    private let pageSizeOptions = [10, 20, 50, 100]

    // MARK: - Computed Properties
    private var buttonSize: CGFloat {
        switch size {
        case .small: return 32
        case .medium: return 40
        case .large: return 48
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
        case .small: return 12
        case .medium: return 14
        case .large: return 16
        }
    }

    private var visiblePages: [Int] {
        guard totalPages > 0 else { return [] }

        let halfVisible = maxVisiblePages / 2
        var start = max(1, currentPage - halfVisible)
        var end = min(totalPages, start + maxVisiblePages - 1)

        // Adjust start if we're near the end
        if end - start + 1 < maxVisiblePages {
            start = max(1, end - maxVisiblePages + 1)
        }

        return Array(start...end)
    }

    private var startItem: Int {
        (currentPage - 1) * pageSize + 1
    }

    private var endItem: Int {
        min(currentPage * pageSize, totalItems)
    }

    // MARK: - Initializer
    public init(
        currentPage: Binding<Int>,
        totalPages: Int,
        style: PaginationStyle = .numbered,
        size: PaginationSize = .medium,
        showPageSize: Bool = false,
        showInfo: Bool = false,
        showJumpToPage: Bool = false,
        maxVisiblePages: Int = 7,
        pageSize: Int = 20,
        totalItems: Int = 0,
        onPageChange: ((Int) -> Void)? = nil,
        onPageSizeChange: ((Int) -> Void)? = nil
    ) {
        self._currentPage = currentPage
        self.totalPages = totalPages
        self.style = style
        self.size = size
        self.showPageSize = showPageSize
        self.showInfo = showInfo
        self.showJumpToPage = showJumpToPage
        self.maxVisiblePages = maxVisiblePages
        self.pageSize = pageSize
        self.totalItems = totalItems
        self.onPageChange = onPageChange
        self.onPageSizeChange = onPageSizeChange
    }

    // MARK: - Body
    public var body: some View {
        VStack(spacing: Spacing.sm) {
            // Main pagination controls
            HStack(spacing: Spacing.xs) {
                switch style {
                case .numbered:
                    numberedPagination
                case .simple:
                    simplePagination
                case .compact:
                    compactPagination
                }
            }

            // Additional controls
            if showPageSize || showInfo || showJumpToPage {
                HStack(spacing: Spacing.md) {
                    if showInfo {
                        paginationInfo
                    }

                    Spacer()

                    if showJumpToPage {
                        jumpToPageControl
                    }

                    if showPageSize {
                        pageSizeSelector
                    }
                }
            }
        }
    }

    // MARK: - Numbered Pagination
    @ViewBuilder
    private var numberedPagination: some View {
        HStack(spacing: Spacing.xs) {
            // Previous button
            paginationButton(
                icon: "chevron.left",
                isEnabled: currentPage > 1
            ) {
                navigateToPage(currentPage - 1)
            }

            // First page (if not visible)
            if !visiblePages.contains(1) && totalPages > 1 {
                paginationButton(text: "1") {
                    navigateToPage(1)
                }

                if visiblePages.first! > 2 {
                    Text("...")
                        .font(.system(size: fontSize))
                        .foregroundColor(colors.textSecondary)
                        .frame(width: buttonSize, height: buttonSize)
                }
            }

            // Visible page numbers
            ForEach(visiblePages, id: \.self) { page in
                paginationButton(
                    text: "\(page)",
                    isSelected: page == currentPage
                ) {
                    navigateToPage(page)
                }
            }

            // Last page (if not visible)
            if !visiblePages.contains(totalPages) && totalPages > 1 {
                if visiblePages.last! < totalPages - 1 {
                    Text("...")
                        .font(.system(size: fontSize))
                        .foregroundColor(colors.textSecondary)
                        .frame(width: buttonSize, height: buttonSize)
                }

                paginationButton(text: "\(totalPages)") {
                    navigateToPage(totalPages)
                }
            }

            // Next button
            paginationButton(
                icon: "chevron.right",
                isEnabled: currentPage < totalPages
            ) {
                navigateToPage(currentPage + 1)
            }
        }
    }

    // MARK: - Simple Pagination
    @ViewBuilder
    private var simplePagination: some View {
        HStack(spacing: Spacing.md) {
            // Previous button
            Button(action: {
                navigateToPage(currentPage - 1)
            }) {
                HStack(spacing: Spacing.xs) {
                    Image(systemName: "chevron.left")
                        .font(.system(size: iconSize))

                    Text("Previous")
                        .font(.system(size: fontSize, weight: .medium))
                }
                .foregroundColor(currentPage > 1 ? colors.primary : colors.textDisabled)
                .padding(.horizontal, Spacing.md)
                .padding(.vertical, Spacing.sm)
                .background(
                    RoundedRectangle(cornerRadius: 6)
                        .stroke(currentPage > 1 ? colors.border : colors.border.opacity(0.5), lineWidth: 1)
                )
            }
            .disabled(currentPage <= 1)

            Spacer()

            // Page info
            Text("\(currentPage) of \(totalPages)")
                .font(.system(size: fontSize))
                .foregroundColor(colors.textSecondary)

            Spacer()

            // Next button
            Button(action: {
                navigateToPage(currentPage + 1)
            }) {
                HStack(spacing: Spacing.xs) {
                    Text("Next")
                        .font(.system(size: fontSize, weight: .medium))

                    Image(systemName: "chevron.right")
                        .font(.system(size: iconSize))
                }
                .foregroundColor(currentPage < totalPages ? colors.primary : colors.textDisabled)
                .padding(.horizontal, Spacing.md)
                .padding(.vertical, Spacing.sm)
                .background(
                    RoundedRectangle(cornerRadius: 6)
                        .stroke(currentPage < totalPages ? colors.border : colors.border.opacity(0.5), lineWidth: 1)
                )
            }
            .disabled(currentPage >= totalPages)
        }
    }

    // MARK: - Compact Pagination
    @ViewBuilder
    private var compactPagination: some View {
        HStack(spacing: Spacing.xs) {
            // Previous button
            paginationButton(
                icon: "chevron.left",
                isEnabled: currentPage > 1
            ) {
                navigateToPage(currentPage - 1)
            }

            // Current page input/display
            Button(action: {
                showPageInput.toggle()
                if showPageInput {
                    inputPage = "\(currentPage)"
                }
            }) {
                if showPageInput {
                    TextField("Page", text: $inputPage)
                        .font(.system(size: fontSize))
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .frame(width: 60)
                        .keyboardType(.numberPad)
                        .onSubmit {
                            if let page = Int(inputPage), page >= 1, page <= totalPages {
                                navigateToPage(page)
                            }
                            showPageInput = false
                        }
                } else {
                    Text("\(currentPage)")
                        .font(.system(size: fontSize, weight: .medium))
                        .foregroundColor(colors.textPrimary)
                        .frame(width: 60)
                        .padding(.vertical, Spacing.xs)
                        .background(
                            RoundedRectangle(cornerRadius: 6)
                                .stroke(colors.border, lineWidth: 1)
                        )
                }
            }

            Text("of \(totalPages)")
                .font(.system(size: fontSize))
                .foregroundColor(colors.textSecondary)

            // Next button
            paginationButton(
                icon: "chevron.right",
                isEnabled: currentPage < totalPages
            ) {
                navigateToPage(currentPage + 1)
            }
        }
    }

    // MARK: - Pagination Button
    @ViewBuilder
    private func paginationButton(
        text: String? = nil,
        icon: String? = nil,
        isSelected: Bool = false,
        isEnabled: Bool = true,
        action: @escaping () -> Void
    ) -> some View {
        Button(action: action) {
            Group {
                if let text = text {
                    Text(text)
                        .font(.system(size: fontSize, weight: .medium))
                } else if let icon = icon {
                    Image(systemName: icon)
                        .font(.system(size: iconSize))
                }
            }
            .foregroundColor(
                isSelected ? colors.textOnPrimary :
                isEnabled ? colors.textPrimary : colors.textDisabled
            )
            .frame(width: buttonSize, height: buttonSize)
            .background(
                RoundedRectangle(cornerRadius: 6)
                    .fill(isSelected ? colors.primary : Color.clear)
                    .overlay(
                        RoundedRectangle(cornerRadius: 6)
                            .stroke(
                                isSelected ? colors.primary : colors.border,
                                lineWidth: 1
                            )
                    )
            )
        }
        .disabled(!isEnabled)
    }

    // MARK: - Additional Controls
    @ViewBuilder
    private var paginationInfo: some View {
        Text("Showing \(startItem)-\(endItem) of \(totalItems)")
            .font(.system(size: fontSize - 1))
            .foregroundColor(colors.textSecondary)
    }

    @ViewBuilder
    private var jumpToPageControl: some View {
        HStack(spacing: Spacing.xs) {
            Text("Go to:")
                .font(.system(size: fontSize - 1))
                .foregroundColor(colors.textSecondary)

            TextField("Page", text: $inputPage)
                .font(.system(size: fontSize))
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .frame(width: 60)
                .keyboardType(.numberPad)
                .onSubmit {
                    if let page = Int(inputPage), page >= 1, page <= totalPages {
                        navigateToPage(page)
                        inputPage = ""
                    }
                }
        }
    }

    @ViewBuilder
    private var pageSizeSelector: some View {
        HStack(spacing: Spacing.xs) {
            Text("Show:")
                .font(.system(size: fontSize - 1))
                .foregroundColor(colors.textSecondary)

            Menu {
                ForEach(pageSizeOptions, id: \.self) { size in
                    Button("\(size) per page") {
                        onPageSizeChange?(size)
                    }
                }
            } label: {
                HStack(spacing: Spacing.xs) {
                    Text("\(pageSize)")
                        .font(.system(size: fontSize))

                    Image(systemName: "chevron.down")
                        .font(.system(size: iconSize - 2))
                }
                .foregroundColor(colors.textPrimary)
                .padding(.horizontal, Spacing.sm)
                .padding(.vertical, Spacing.xs)
                .background(
                    RoundedRectangle(cornerRadius: 6)
                        .stroke(colors.border, lineWidth: 1)
                )
            }
        }
    }

    // MARK: - Actions
    private func navigateToPage(_ page: Int) {
        guard page >= 1 && page <= totalPages && page != currentPage else { return }

        currentPage = page
        onPageChange?(page)

        // Haptic feedback
        let impactFeedback = UIImpactFeedbackGenerator(style: .light)
        impactFeedback.impactOccurred()
    }
}

// MARK: - Preview
#if DEBUG
struct TchatPagination_Previews: PreviewProvider {
    static var previews: some View {
        ScrollView {
            VStack(spacing: Spacing.lg) {
                // Numbered pagination
                TchatPagination(
                    currentPage: .constant(5),
                    totalPages: 20,
                    style: .numbered,
                    size: .medium,
                    showPageSize: true,
                    showInfo: true,
                    showJumpToPage: true,
                    totalItems: 1000
                )

                Divider()

                // Simple pagination
                TchatPagination(
                    currentPage: .constant(2),
                    totalPages: 10,
                    style: .simple,
                    size: .large
                )

                Divider()

                // Compact pagination
                TchatPagination(
                    currentPage: .constant(3),
                    totalPages: 15,
                    style: .compact,
                    size: .small
                )

                Divider()

                // Small numbered pagination
                TchatPagination(
                    currentPage: .constant(1),
                    totalPages: 5,
                    style: .numbered,
                    size: .small,
                    maxVisiblePages: 5
                )
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .previewLayout(.sizeThatFits)
    }
}
#endif