// Media Store iOS API Service
// Generated for Media Store Tabs feature implementation

import Foundation
import Combine
import Alamofire

public protocol MediaAPIServiceProtocol {
    // Categories
    func getMediaCategories() -> AnyPublisher<MediaCategoriesResponse, MediaAPIError>
    func getMediaCategory(id: String) -> AnyPublisher<MediaCategory, MediaAPIError>

    // Subtabs
    func getMovieSubtabs() -> AnyPublisher<MediaSubtabsResponse, MediaAPIError>

    // Content
    func getContentByCategory(categoryId: String, page: Int, limit: Int, subtab: String?) -> AnyPublisher<MediaContentResponse, MediaAPIError>
    func getFeaturedContent(limit: Int?, categoryId: String?) -> AnyPublisher<MediaFeaturedResponse, MediaAPIError>
    func searchMediaContent(query: String, categoryId: String?, page: Int, limit: Int) -> AnyPublisher<MediaSearchResponse, MediaAPIError>

    // Store Integration
    func getMediaProducts(categoryId: String?, page: Int, limit: Int) -> AnyPublisher<MediaProductsResponse, MediaAPIError>
    func addMediaToCart(request: AddMediaToCartRequest) -> AnyPublisher<AddMediaToCartResponse, MediaAPIError>
    func getUnifiedCart() -> AnyPublisher<UnifiedCartResponse, MediaAPIError>
    func removeMediaFromCart(cartItemId: String) -> AnyPublisher<Void, MediaAPIError>
    func updateMediaCartItem(cartItemId: String, quantity: Int) -> AnyPublisher<MediaCartItem, MediaAPIError>
    func validateMediaCheckout(request: MediaCheckoutValidationRequest) -> AnyPublisher<MediaCheckoutValidationResponse, MediaAPIError>
    func processMediaCheckout(request: ProcessMediaCheckoutRequest) -> AnyPublisher<MediaOrder, MediaAPIError>
    func getMediaOrders(page: Int, limit: Int, status: String?) -> AnyPublisher<MediaOrdersResponse, MediaAPIError>
    func getMediaOrder(orderId: String) -> AnyPublisher<MediaOrder, MediaAPIError>
    func downloadMediaContent(orderItemId: String) -> AnyPublisher<MediaDownloadResponse, MediaAPIError>
}

public class MediaAPIService: MediaAPIServiceProtocol {
    private let session: Session
    private let baseURL: String
    private let authService: AuthServiceProtocol

    public init(baseURL: String, authService: AuthServiceProtocol) {
        self.baseURL = baseURL
        self.authService = authService
        self.session = Session(configuration: URLSessionConfiguration.default)
    }

    // MARK: - Categories
    public func getMediaCategories() -> AnyPublisher<MediaCategoriesResponse, MediaAPIError> {
        return request(endpoint: "/media/categories", method: .get)
    }

    public func getMediaCategory(id: String) -> AnyPublisher<MediaCategory, MediaAPIError> {
        return request(endpoint: "/media/categories/\(id)", method: .get)
    }

    // MARK: - Subtabs
    public func getMovieSubtabs() -> AnyPublisher<MediaSubtabsResponse, MediaAPIError> {
        return request(endpoint: "/media/movies/subtabs", method: .get)
    }

    // MARK: - Content
    public func getContentByCategory(
        categoryId: String,
        page: Int = 1,
        limit: Int = 20,
        subtab: String? = nil
    ) -> AnyPublisher<MediaContentResponse, MediaAPIError> {
        var parameters: [String: Any] = [
            "page": page,
            "limit": limit
        ]
        if let subtab = subtab {
            parameters["subtab"] = subtab
        }

        return request(
            endpoint: "/media/category/\(categoryId)/content",
            method: .get,
            parameters: parameters
        )
    }

    public func getFeaturedContent(
        limit: Int? = 10,
        categoryId: String? = nil
    ) -> AnyPublisher<MediaFeaturedResponse, MediaAPIError> {
        var parameters: [String: Any] = [:]
        if let limit = limit {
            parameters["limit"] = limit
        }
        if let categoryId = categoryId {
            parameters["categoryId"] = categoryId
        }

        return request(endpoint: "/media/featured", method: .get, parameters: parameters)
    }

    public func searchMediaContent(
        query: String,
        categoryId: String? = nil,
        page: Int = 1,
        limit: Int = 20
    ) -> AnyPublisher<MediaSearchResponse, MediaAPIError> {
        var parameters: [String: Any] = [
            "q": query,
            "page": page,
            "limit": limit
        ]
        if let categoryId = categoryId {
            parameters["categoryId"] = categoryId
        }

        return request(endpoint: "/media/search", method: .get, parameters: parameters)
    }

    // MARK: - Store Integration
    public func getMediaProducts(
        categoryId: String? = nil,
        page: Int = 1,
        limit: Int = 20
    ) -> AnyPublisher<MediaProductsResponse, MediaAPIError> {
        var parameters: [String: Any] = [
            "page": page,
            "limit": limit
        ]
        if let categoryId = categoryId {
            parameters["categoryId"] = categoryId
        }

        return request(endpoint: "/store/products/media", method: .get, parameters: parameters)
    }

    public func addMediaToCart(request: AddMediaToCartRequest) -> AnyPublisher<AddMediaToCartResponse, MediaAPIError> {
        return self.request(endpoint: "/store/cart/add-media", method: .post, body: request)
    }

    public func getUnifiedCart() -> AnyPublisher<UnifiedCartResponse, MediaAPIError> {
        return request(endpoint: "/store/cart", method: .get)
    }

    public func removeMediaFromCart(cartItemId: String) -> AnyPublisher<Void, MediaAPIError> {
        return request(endpoint: "/store/cart/items/\(cartItemId)", method: .delete)
    }

    public func updateMediaCartItem(cartItemId: String, quantity: Int) -> AnyPublisher<MediaCartItem, MediaAPIError> {
        let body = ["quantity": quantity]
        return request(endpoint: "/store/cart/items/\(cartItemId)", method: .patch, body: body)
    }

    public func validateMediaCheckout(request: MediaCheckoutValidationRequest) -> AnyPublisher<MediaCheckoutValidationResponse, MediaAPIError> {
        return self.request(endpoint: "/store/checkout/media-validation", method: .post, body: request)
    }

    public func processMediaCheckout(request: ProcessMediaCheckoutRequest) -> AnyPublisher<MediaOrder, MediaAPIError> {
        return self.request(endpoint: "/store/checkout/media", method: .post, body: request)
    }

    public func getMediaOrders(
        page: Int = 1,
        limit: Int = 20,
        status: String? = nil
    ) -> AnyPublisher<MediaOrdersResponse, MediaAPIError> {
        var parameters: [String: Any] = [
            "page": page,
            "limit": limit
        ]
        if let status = status {
            parameters["status"] = status
        }

        return request(endpoint: "/store/orders/media", method: .get, parameters: parameters)
    }

    public func getMediaOrder(orderId: String) -> AnyPublisher<MediaOrder, MediaAPIError> {
        return request(endpoint: "/store/orders/media/\(orderId)", method: .get)
    }

    public func downloadMediaContent(orderItemId: String) -> AnyPublisher<MediaDownloadResponse, MediaAPIError> {
        return request(endpoint: "/store/orders/media/items/\(orderItemId)/download", method: .get)
    }

    // MARK: - Private Methods
    private func request<T: Codable>(
        endpoint: String,
        method: HTTPMethod,
        parameters: [String: Any]? = nil,
        body: Encodable? = nil
    ) -> AnyPublisher<T, MediaAPIError> {
        let url = baseURL + "/api/v1" + endpoint

        return Future<T, MediaAPIError> { [weak self] promise in
            guard let self = self else {
                promise(.failure(MediaAPIError.networkError("Service deallocated")))
                return
            }

            var headers: HTTPHeaders = [
                "Content-Type": "application/json",
                "Accept": "application/json"
            ]

            // Add authentication token if available
            if let token = self.authService.currentToken {
                headers["Authorization"] = "Bearer \(token)"
            }

            var request = self.session.request(
                url,
                method: method,
                parameters: parameters,
                encoding: URLEncoding.default,
                headers: headers
            )

            // Add body for POST/PUT/PATCH requests
            if let body = body {
                do {
                    let jsonData = try JSONEncoder().encode(body)
                    request = self.session.upload(jsonData, to: url, method: method, headers: headers)
                } catch {
                    promise(.failure(MediaAPIError.encodingError(error.localizedDescription)))
                    return
                }
            }

            request
                .validate()
                .responseDecodable(of: T.self) { response in
                    switch response.result {
                    case .success(let value):
                        promise(.success(value))
                    case .failure(let error):
                        let mediaError = self.mapAlamofireError(error, data: response.data)
                        promise(.failure(mediaError))
                    }
                }
        }
        .eraseToAnyPublisher()
    }

    private func mapAlamofireError(_ error: AFError, data: Data?) -> MediaAPIError {
        switch error {
        case .responseValidationFailed(reason: .unacceptableStatusCode(code: let statusCode)):
            if let data = data,
               let errorResponse = try? JSONDecoder().decode(MediaAPIErrorResponse.self, from: data) {
                return MediaAPIError.serverError(statusCode, errorResponse.message)
            }
            return MediaAPIError.serverError(statusCode, "HTTP \(statusCode)")

        case .responseSerializationFailed:
            return MediaAPIError.decodingError("Failed to decode response")

        case .sessionTaskFailed(error: let urlError as URLError):
            switch urlError.code {
            case .notConnectedToInternet, .networkConnectionLost:
                return MediaAPIError.noInternetConnection
            case .timedOut:
                return MediaAPIError.timeout
            default:
                return MediaAPIError.networkError(urlError.localizedDescription)
            }

        default:
            return MediaAPIError.networkError(error.localizedDescription)
        }
    }
}

// MARK: - Additional Response Types
public struct MediaProductsResponse: Codable {
    public let products: [MediaProduct]
    public let pagination: PaginationInfo

    public struct PaginationInfo: Codable {
        public let page: Int
        public let limit: Int
        public let total: Int
        public let hasMore: Bool
    }
}

public struct MediaSearchResponse: Codable {
    public let items: [MediaContentItem]
    public let query: String
    public let total: Int
    public let page: Int
}

public struct AddMediaToCartRequest: Codable {
    public let mediaContentId: String
    public let quantity: Int
    public let mediaLicense: MediaLicense
    public let downloadFormat: MediaDownloadFormat

    public init(mediaContentId: String, quantity: Int, mediaLicense: MediaLicense, downloadFormat: MediaDownloadFormat) {
        self.mediaContentId = mediaContentId
        self.quantity = quantity
        self.mediaLicense = mediaLicense
        self.downloadFormat = downloadFormat
    }
}

public struct AddMediaToCartResponse: Codable {
    public let cartId: String
    public let itemsCount: Int
    public let totalAmount: Double
    public let currency: String
    public let addedItem: MediaCartItem
}

public struct UnifiedCartResponse: Codable {
    public let cartId: String
    public let physicalItems: [MediaCartItem]
    public let mediaItems: [MediaCartItem]
    public let totalPhysicalAmount: Double
    public let totalMediaAmount: Double
    public let totalAmount: Double
    public let currency: String
    public let itemsCount: Int
}

public struct MediaCheckoutValidationRequest: Codable {
    public let cartId: String
    public let mediaItems: [MediaCartItem]

    public init(cartId: String, mediaItems: [MediaCartItem]) {
        self.cartId = cartId
        self.mediaItems = mediaItems
    }
}

public struct MediaCheckoutValidationResponse: Codable {
    public let isValid: Bool
    public let validItems: [MediaCartItem]
    public let invalidItems: [MediaCartItem]
    public let totalMediaAmount: Double
    public let estimatedDeliveryTime: String
}

public struct ProcessMediaCheckoutRequest: Codable {
    public let cartId: String
    public let mediaItems: [MediaCartItem]
    public let paymentMethod: String
    public let billingAddress: String?

    public init(cartId: String, mediaItems: [MediaCartItem], paymentMethod: String, billingAddress: String? = nil) {
        self.cartId = cartId
        self.mediaItems = mediaItems
        self.paymentMethod = paymentMethod
        self.billingAddress = billingAddress
    }
}

public struct MediaOrdersResponse: Codable {
    public let orders: [MediaOrder]
    public let pagination: PaginationInfo

    public struct PaginationInfo: Codable {
        public let page: Int
        public let limit: Int
        public let total: Int
        public let hasMore: Bool
    }
}

public struct MediaDownloadResponse: Codable {
    public let downloadUrl: String
    public let expiresAt: Date
}

// MARK: - Error Types
public enum MediaAPIError: Error, LocalizedError {
    case networkError(String)
    case serverError(Int, String)
    case decodingError(String)
    case encodingError(String)
    case noInternetConnection
    case timeout
    case unauthorized
    case forbidden
    case notFound
    case invalidRequest(String)

    public var errorDescription: String? {
        switch self {
        case .networkError(let message):
            return "Network error: \(message)"
        case .serverError(let code, let message):
            return "Server error (\(code)): \(message)"
        case .decodingError(let message):
            return "Data parsing error: \(message)"
        case .encodingError(let message):
            return "Request encoding error: \(message)"
        case .noInternetConnection:
            return "No internet connection available"
        case .timeout:
            return "Request timed out"
        case .unauthorized:
            return "Authentication required"
        case .forbidden:
            return "Access denied"
        case .notFound:
            return "Resource not found"
        case .invalidRequest(let message):
            return "Invalid request: \(message)"
        }
    }
}

private struct MediaAPIErrorResponse: Codable {
    let error: String
    let message: String
    let details: [String: AnyCodable]?
}

// MARK: - Auth Service Protocol
public protocol AuthServiceProtocol {
    var currentToken: String? { get }
}