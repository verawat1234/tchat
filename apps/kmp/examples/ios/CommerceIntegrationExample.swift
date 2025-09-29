import SwiftUI
import CommerceKMP

/**
 * iOS Integration Example for KMP Commerce Module
 * Demonstrates how to integrate the shared KMP commerce functionality
 * into a SwiftUI iOS application
 */

@main
struct CommerceApp: App {
    // Initialize the commerce system early in the app lifecycle
    @StateObject private var commerceSetup = CommerceSetup()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(commerceSetup.commerceManager)
                .onAppear {
                    // Initialize commerce when app appears
                    commerceSetup.initializeCommerce()
                }
        }
    }
}

/**
 * Commerce setup and dependency injection for iOS
 */
class CommerceSetup: ObservableObject {
    @Published var commerceManager: CommerceManager?
    @Published var isInitialized = false

    func initializeCommerce() {
        // Create storage implementation
        let storage = IOSCommerceStorage()

        // Create API client with your backend URL
        let apiClient = CommerceApiClientImpl(
            baseUrl: "https://your-api.com/api/v1",
            httpClient: KtorHttpClient() // Your Ktor client setup
        )

        // Create repositories
        let cartRepository = CartRepositoryImpl(apiClient: apiClient, storage: storage)
        let productRepository = ProductRepositoryImpl(apiClient: apiClient, storage: storage)
        let categoryRepository = CategoryRepositoryImpl(apiClient: apiClient, storage: storage)

        // Create main commerce manager
        let manager = CommerceManagerImpl(
            cartRepository: cartRepository,
            productRepository: productRepository,
            categoryRepository: categoryRepository,
            apiClient: apiClient,
            storage: storage
        )

        // Initialize asynchronously
        Task {
            let result = await manager.initialize()

            await MainActor.run {
                if result.isSuccess {
                    self.commerceManager = manager
                    self.isInitialized = true
                    print("✅ Commerce system initialized successfully")
                } else {
                    print("❌ Failed to initialize commerce system: \(result.exceptionOrNull()?.localizedDescription ?? "Unknown error")")
                }
            }
        }
    }
}

/**
 * Main content view demonstrating commerce integration
 */
struct ContentView: View {
    @EnvironmentObject var commerceManager: CommerceManager
    @StateObject private var viewModel = CommerceViewModel()

    var body: some View {
        NavigationView {
            TabView {
                // Products tab
                ProductListView()
                    .tabItem {
                        Image(systemName: "list.bullet")
                        Text("Products")
                    }

                // Cart tab
                CartView()
                    .tabItem {
                        Image(systemName: "cart")
                        Text("Cart")
                    }

                // Categories tab
                CategoryListView()
                    .tabItem {
                        Image(systemName: "folder")
                        Text("Categories")
                    }
            }
        }
        .onAppear {
            // Set the commerce manager in the view model
            viewModel.setCommerceManager(commerceManager)
        }
    }
}

/**
 * Product list view using KMP commerce data
 */
struct ProductListView: View {
    @StateObject private var viewModel = ProductListViewModel()

    var body: some View {
        NavigationView {
            List {
                if viewModel.isLoading {
                    ProgressView("Loading products...")
                        .frame(maxWidth: .infinity, alignment: .center)
                } else {
                    ForEach(viewModel.products, id: \.id) { product in
                        ProductRowView(product: product) {
                            // Add to cart action
                            Task {
                                await viewModel.addToCart(productId: product.id, quantity: 1)
                            }
                        }
                    }
                }
            }
            .navigationTitle("Products")
            .refreshable {
                await viewModel.refreshProducts()
            }
            .searchable(text: $viewModel.searchText, prompt: "Search products")
            .onChange(of: viewModel.searchText) { newValue in
                Task {
                    await viewModel.searchProducts(query: newValue)
                }
            }
        }
        .onAppear {
            Task {
                await viewModel.loadProducts()
            }
        }
    }
}

/**
 * Product row component
 */
struct ProductRowView: View {
    let product: Product
    let onAddToCart: () -> Void
    @State private var isAddingToCart = false

    var body: some View {
        HStack {
            // Product image placeholder
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.3))
                .frame(width: 60, height: 60)
                .overlay(
                    Image(systemName: "photo")
                        .foregroundColor(.gray)
                )

            VStack(alignment: .leading, spacing: 4) {
                Text(product.name)
                    .font(.headline)
                    .lineLimit(1)

                Text(product.description_)
                    .font(.caption)
                    .foregroundColor(.secondary)
                    .lineLimit(2)

                HStack {
                    Text(product.formattedPrice())
                        .font(.subheadline)
                        .fontWeight(.semibold)
                        .foregroundColor(.primary)

                    if let comparePrice = product.compareAtPrice, comparePrice > product.price {
                        Text(product.formattedComparePrice())
                            .font(.caption)
                            .strikethrough()
                            .foregroundColor(.secondary)
                    }

                    Spacer()

                    // Stock indicator
                    if product.inventory.isInStock {
                        Text("In Stock")
                            .font(.caption)
                            .foregroundColor(.green)
                    } else {
                        Text("Out of Stock")
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                }
            }

            Spacer()

            // Add to cart button
            Button(action: {
                isAddingToCart = true
                onAddToCart()

                // Reset button state after animation
                DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
                    isAddingToCart = false
                }
            }) {
                if isAddingToCart {
                    ProgressView()
                        .scaleEffect(0.8)
                } else {
                    Image(systemName: "cart.badge.plus")
                        .foregroundColor(.white)
                }
            }
            .frame(width: 44, height: 44)
            .background(product.inventory.isInStock ? Color.blue : Color.gray)
            .clipShape(Circle())
            .disabled(!product.inventory.isInStock || isAddingToCart)
        }
        .padding(.vertical, 4)
    }
}

/**
 * Cart view showing current cart items
 */
struct CartView: View {
    @StateObject private var viewModel = CartViewModel()

    var body: some View {
        NavigationView {
            VStack {
                if viewModel.cart.items.isEmpty {
                    EmptyCartView {
                        // Handle continue shopping
                    }
                } else {
                    List {
                        ForEach(viewModel.cart.items, id: \.id) { item in
                            CartItemRowView(
                                item: item,
                                onQuantityChange: { newQuantity in
                                    Task {
                                        await viewModel.updateQuantity(itemId: item.id, quantity: newQuantity)
                                    }
                                },
                                onRemove: {
                                    Task {
                                        await viewModel.removeItem(itemId: item.id)
                                    }
                                }
                            )
                        }

                        // Cart summary
                        CartSummaryView(summary: viewModel.cartSummary) {
                            // Handle checkout
                            print("Proceeding to checkout...")
                        }
                    }
                }
            }
            .navigationTitle("Cart (\(viewModel.cart.itemCount))")
            .refreshable {
                await viewModel.refreshCart()
            }
        }
        .onAppear {
            Task {
                await viewModel.loadCart()
            }
        }
    }
}

/**
 * Empty cart placeholder view
 */
struct EmptyCartView: View {
    let onContinueShopping: () -> Void

    var body: some View {
        VStack(spacing: 24) {
            Image(systemName: "cart")
                .font(.system(size: 64))
                .foregroundColor(.gray)

            Text("Your cart is empty")
                .font(.title2)
                .fontWeight(.medium)

            Text("Add some products to get started")
                .font(.body)
                .foregroundColor(.secondary)

            Button("Continue Shopping") {
                onContinueShopping()
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}

/**
 * Cart item row component
 */
struct CartItemRowView: View {
    let item: CartItem
    let onQuantityChange: (Int32) -> Void
    let onRemove: () -> Void

    var body: some View {
        HStack {
            // Product image placeholder
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.3))
                .frame(width: 50, height: 50)

            VStack(alignment: .leading, spacing: 4) {
                Text(item.productName)
                    .font(.headline)
                    .lineLimit(1)

                Text(item.businessName)
                    .font(.caption)
                    .foregroundColor(.secondary)

                Text(item.formattedPrice())
                    .font(.subheadline)
                    .fontWeight(.semibold)
            }

            Spacer()

            // Quantity controls
            HStack {
                Button("-") {
                    if item.quantity > 1 {
                        onQuantityChange(item.quantity - 1)
                    }
                }
                .frame(width: 30, height: 30)
                .background(Color.gray.opacity(0.2))
                .clipShape(Circle())
                .disabled(item.quantity <= 1)

                Text("\(item.quantity)")
                    .frame(width: 30)
                    .font(.headline)

                Button("+") {
                    onQuantityChange(item.quantity + 1)
                }
                .frame(width: 30, height: 30)
                .background(Color.blue.opacity(0.2))
                .clipShape(Circle())
            }

            // Remove button
            Button(action: onRemove) {
                Image(systemName: "trash")
                    .foregroundColor(.red)
            }
            .frame(width: 30, height: 30)
        }
        .padding(.vertical, 4)
    }
}

/**
 * Cart summary component
 */
struct CartSummaryView: View {
    let summary: CartSummary
    let onCheckout: () -> Void

    var body: some View {
        VStack(spacing: 12) {
            Divider()

            HStack {
                Text("Subtotal")
                Spacer()
                Text(summary.formattedSubtotal())
            }

            if summary.discountAmount > 0 {
                HStack {
                    Text("Discount")
                        .foregroundColor(.green)
                    Spacer()
                    Text("-\(summary.formattedDiscount())")
                        .foregroundColor(.green)
                }
            }

            HStack {
                Text("Shipping")
                Spacer()
                Text(summary.formattedShipping())
            }

            HStack {
                Text("Tax")
                Spacer()
                Text(summary.formattedTax())
            }

            Divider()

            HStack {
                Text("Total")
                    .font(.headline)
                    .fontWeight(.bold)
                Spacer()
                Text(summary.formattedTotal())
                    .font(.headline)
                    .fontWeight(.bold)
            }

            Button("Proceed to Checkout") {
                onCheckout()
            }
            .buttonStyle(.borderedProminent)
            .frame(maxWidth: .infinity)
            .disabled(summary.totalItems == 0)
        }
        .padding()
        .background(Color(UIColor.secondarySystemBackground))
        .cornerRadius(12)
    }
}

/**
 * Category list view
 */
struct CategoryListView: View {
    @StateObject private var viewModel = CategoryListViewModel()

    var body: some View {
        NavigationView {
            List {
                if viewModel.isLoading {
                    ProgressView("Loading categories...")
                        .frame(maxWidth: .infinity, alignment: .center)
                } else {
                    ForEach(viewModel.categories, id: \.id) { category in
                        CategoryRowView(category: category) {
                            // Handle category selection
                            Task {
                                await viewModel.selectCategory(category)
                            }
                        }
                    }
                }
            }
            .navigationTitle("Categories")
            .refreshable {
                await viewModel.refreshCategories()
            }
        }
        .onAppear {
            Task {
                await viewModel.loadCategories()
            }
        }
    }
}

/**
 * Category row component
 */
struct CategoryRowView: View {
    let category: Category
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            HStack {
                // Category icon placeholder
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.blue.opacity(0.1))
                    .frame(width: 40, height: 40)
                    .overlay(
                        Image(systemName: "folder.fill")
                            .foregroundColor(.blue)
                    )

                VStack(alignment: .leading, spacing: 4) {
                    Text(category.name)
                        .font(.headline)
                        .foregroundColor(.primary)

                    if let description = category.description_ {
                        Text(description)
                            .font(.caption)
                            .foregroundColor(.secondary)
                            .lineLimit(1)
                    }

                    Text("\(category.productCount) products")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }

                Spacer()

                if category.isFeatured {
                    Image(systemName: "star.fill")
                        .foregroundColor(.yellow)
                }

                Image(systemName: "chevron.right")
                    .foregroundColor(.gray)
                    .font(.caption)
            }
            .padding(.vertical, 4)
        }
        .buttonStyle(PlainButtonStyle())
    }
}

/**
 * Extension to add formatting helpers to KMP models
 */
extension Product {
    func formattedPrice() -> String {
        return "\(currency) \(String(format: "%.2f", price))"
    }

    func formattedComparePrice() -> String {
        guard let comparePrice = compareAtPrice else { return "" }
        return "\(currency) \(String(format: "%.2f", comparePrice))"
    }
}

extension CartItem {
    func formattedPrice() -> String {
        return "\(currency) \(String(format: "%.2f", totalPrice))"
    }
}

extension CartSummary {
    func formattedSubtotal() -> String {
        return "\(currency) \(String(format: "%.2f", subtotal))"
    }

    func formattedTotal() -> String {
        return "\(currency) \(String(format: "%.2f", total))"
    }

    func formattedDiscount() -> String {
        return "\(currency) \(String(format: "%.2f", discountAmount))"
    }

    func formattedShipping() -> String {
        return "\(currency) \(String(format: "%.2f", shipping))"
    }

    func formattedTax() -> String {
        return "\(currency) \(String(format: "%.2f", tax))"
    }
}