// Media Store Swift models
// Generated for Media Store Tabs feature implementation

import Foundation

// MARK: - Media Category
public struct MediaCategory: Codable, Identifiable, Equatable {
    public let id: String
    public let name: String
    public let displayOrder: Int
    public let iconName: String
    public let isActive: Bool
    public let featuredContentEnabled: Bool
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, name: String, displayOrder: Int, iconName: String,
                isActive: Bool, featuredContentEnabled: Bool, createdAt: Date, updatedAt: Date) {
        self.id = id
        self.name = name
        self.displayOrder = displayOrder
        self.iconName = iconName
        self.isActive = isActive
        self.featuredContentEnabled = featuredContentEnabled
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - Media Subtab
public struct MediaSubtab: Codable, Identifiable, Equatable {
    public let id: String
    public let categoryId: String
    public let name: String
    public let displayOrder: Int
    public let filterCriteria: [String: AnyCodable]
    public let isActive: Bool
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, categoryId: String, name: String, displayOrder: Int,
                filterCriteria: [String: AnyCodable], isActive: Bool, createdAt: Date, updatedAt: Date) {
        self.id = id
        self.categoryId = categoryId
        self.name = name
        self.displayOrder = displayOrder
        self.filterCriteria = filterCriteria
        self.isActive = isActive
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - Content Type
public enum MediaContentType: String, Codable, CaseIterable {
    case book = "book"
    case podcast = "podcast"
    case video = "video"
    case cartoon = "cartoon"
}

// MARK: - Availability Status
public enum MediaAvailabilityStatus: String, Codable, CaseIterable {
    case available = "available"
    case comingSoon = "coming_soon"
    case unavailable = "unavailable"
}

// MARK: - Media Content Item
public struct MediaContentItem: Codable, Identifiable, Equatable {
    public let id: String
    public let categoryId: String
    public let title: String
    public let description: String
    public let thumbnailUrl: String
    public let contentUrl: String?
    public let contentType: MediaContentType
    public let duration: Int?
    public let price: Double
    public let currency: String
    public let availabilityStatus: MediaAvailabilityStatus
    public let isFeatured: Bool
    public let featuredOrder: Int?
    public let metadata: [String: AnyCodable]
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, categoryId: String, title: String, description: String,
                thumbnailUrl: String, contentUrl: String?, contentType: MediaContentType,
                duration: Int?, price: Double, currency: String,
                availabilityStatus: MediaAvailabilityStatus, isFeatured: Bool,
                featuredOrder: Int?, metadata: [String: AnyCodable],
                createdAt: Date, updatedAt: Date) {
        self.id = id
        self.categoryId = categoryId
        self.title = title
        self.description = description
        self.thumbnailUrl = thumbnailUrl
        self.contentUrl = contentUrl
        self.contentType = contentType
        self.duration = duration
        self.price = price
        self.currency = currency
        self.availabilityStatus = availabilityStatus
        self.isFeatured = isFeatured
        self.featuredOrder = featuredOrder
        self.metadata = metadata
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }

    public var durationString: String? {
        guard let duration = duration else { return nil }
        let minutes = duration / 60
        let seconds = duration % 60
        return String(format: "%d:%02d", minutes, seconds)
    }
}

// MARK: - Product Type
public enum MediaProductType: String, Codable, CaseIterable {
    case physical = "physical"
    case media = "media"
}

// MARK: - Media License
public enum MediaLicense: String, Codable, CaseIterable {
    case personal = "personal"
    case commercial = "commercial"
    case educational = "educational"
}

// MARK: - Download Format
public enum MediaDownloadFormat: String, Codable, CaseIterable {
    case pdf = "PDF"
    case epub = "EPUB"
    case mp3 = "MP3"
    case mp4 = "MP4"
    case flac = "FLAC"
}

// MARK: - Media Product
public struct MediaProduct: Codable, Identifiable, Equatable {
    public let id: String
    public let name: String
    public let description: String
    public let price: Double
    public let currency: String
    public let productType: MediaProductType
    public let mediaContentId: String?
    public let mediaMetadata: MediaMetadata?
    public let category: String
    public let isActive: Bool
    public let stockQuantity: Int?
    public let createdAt: Date
    public let updatedAt: Date

    public struct MediaMetadata: Codable, Equatable {
        public let contentType: MediaContentType
        public let duration: Int?
        public let format: String?
        public let license: String?

        public init(contentType: MediaContentType, duration: Int?, format: String?, license: String?) {
            self.contentType = contentType
            self.duration = duration
            self.format = format
            self.license = license
        }
    }

    public init(id: String, name: String, description: String, price: Double, currency: String,
                productType: MediaProductType, mediaContentId: String?, mediaMetadata: MediaMetadata?,
                category: String, isActive: Bool, stockQuantity: Int?, createdAt: Date, updatedAt: Date) {
        self.id = id
        self.name = name
        self.description = description
        self.price = price
        self.currency = currency
        self.productType = productType
        self.mediaContentId = mediaContentId
        self.mediaMetadata = mediaMetadata
        self.category = category
        self.isActive = isActive
        self.stockQuantity = stockQuantity
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - Media Cart Item
public struct MediaCartItem: Codable, Identifiable, Equatable {
    public let id: String
    public let cartId: String
    public let productId: String
    public let mediaContentId: String?
    public let quantity: Int
    public let unitPrice: Double
    public let totalPrice: Double
    public let mediaLicense: MediaLicense?
    public let downloadFormat: MediaDownloadFormat?
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, cartId: String, productId: String, mediaContentId: String?,
                quantity: Int, unitPrice: Double, totalPrice: Double,
                mediaLicense: MediaLicense?, downloadFormat: MediaDownloadFormat?,
                createdAt: Date, updatedAt: Date) {
        self.id = id
        self.cartId = cartId
        self.productId = productId
        self.mediaContentId = mediaContentId
        self.quantity = quantity
        self.unitPrice = unitPrice
        self.totalPrice = totalPrice
        self.mediaLicense = mediaLicense
        self.downloadFormat = downloadFormat
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - Order Status
public enum MediaOrderStatus: String, Codable, CaseIterable {
    case pending = "pending"
    case processing = "processing"
    case completed = "completed"
    case cancelled = "cancelled"
}

// MARK: - Delivery Status
public enum MediaDeliveryStatus: String, Codable, CaseIterable {
    case pending = "pending"
    case delivered = "delivered"
    case failed = "failed"
}

// MARK: - Media Order
public struct MediaOrder: Codable, Identifiable, Equatable {
    public let id: String
    public let userId: String
    public let status: MediaOrderStatus
    public let totalPhysicalAmount: Double
    public let totalMediaAmount: Double
    public let totalAmount: Double
    public let currency: String
    public let mediaDeliveryStatus: MediaDeliveryStatus
    public let shippingAddress: String?
    public let items: [MediaOrderItem]
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, userId: String, status: MediaOrderStatus,
                totalPhysicalAmount: Double, totalMediaAmount: Double, totalAmount: Double,
                currency: String, mediaDeliveryStatus: MediaDeliveryStatus,
                shippingAddress: String?, items: [MediaOrderItem],
                createdAt: Date, updatedAt: Date) {
        self.id = id
        self.userId = userId
        self.status = status
        self.totalPhysicalAmount = totalPhysicalAmount
        self.totalMediaAmount = totalMediaAmount
        self.totalAmount = totalAmount
        self.currency = currency
        self.mediaDeliveryStatus = mediaDeliveryStatus
        self.shippingAddress = shippingAddress
        self.items = items
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - Media Order Item
public struct MediaOrderItem: Codable, Identifiable, Equatable {
    public let id: String
    public let orderId: String
    public let productId: String
    public let mediaContentId: String?
    public let quantity: Int
    public let unitPrice: Double
    public let totalPrice: Double
    public let mediaLicense: MediaLicense?
    public let downloadFormat: MediaDownloadFormat?
    public let deliveryStatus: MediaDeliveryStatus?
    public let downloadUrl: String?
    public let createdAt: Date
    public let updatedAt: Date

    public init(id: String, orderId: String, productId: String, mediaContentId: String?,
                quantity: Int, unitPrice: Double, totalPrice: Double,
                mediaLicense: MediaLicense?, downloadFormat: MediaDownloadFormat?,
                deliveryStatus: MediaDeliveryStatus?, downloadUrl: String?,
                createdAt: Date, updatedAt: Date) {
        self.id = id
        self.orderId = orderId
        self.productId = productId
        self.mediaContentId = mediaContentId
        self.quantity = quantity
        self.unitPrice = unitPrice
        self.totalPrice = totalPrice
        self.mediaLicense = mediaLicense
        self.downloadFormat = downloadFormat
        self.deliveryStatus = deliveryStatus
        self.downloadUrl = downloadUrl
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

// MARK: - API Response Types
public struct MediaCategoriesResponse: Codable {
    public let categories: [MediaCategory]
    public let total: Int

    public init(categories: [MediaCategory], total: Int) {
        self.categories = categories
        self.total = total
    }
}

public struct MediaContentResponse: Codable {
    public let items: [MediaContentItem]
    public let page: Int
    public let limit: Int
    public let total: Int
    public let hasMore: Bool

    public init(items: [MediaContentItem], page: Int, limit: Int, total: Int, hasMore: Bool) {
        self.items = items
        self.page = page
        self.limit = limit
        self.total = total
        self.hasMore = hasMore
    }
}

public struct MediaFeaturedResponse: Codable {
    public let items: [MediaContentItem]
    public let total: Int
    public let hasMore: Bool

    public init(items: [MediaContentItem], total: Int, hasMore: Bool) {
        self.items = items
        self.total = total
        self.hasMore = hasMore
    }
}

public struct MediaSubtabsResponse: Codable {
    public let subtabs: [MediaSubtab]
    public let defaultSubtab: String

    public init(subtabs: [MediaSubtab], defaultSubtab: String) {
        self.subtabs = subtabs
        self.defaultSubtab = defaultSubtab
    }
}

// MARK: - Helper Types
public struct AnyCodable: Codable, Equatable {
    public let value: Any

    public init<T: Codable>(_ value: T) {
        self.value = value
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()

        if let intValue = try? container.decode(Int.self) {
            value = intValue
        } else if let doubleValue = try? container.decode(Double.self) {
            value = doubleValue
        } else if let boolValue = try? container.decode(Bool.self) {
            value = boolValue
        } else if let stringValue = try? container.decode(String.self) {
            value = stringValue
        } else if let arrayValue = try? container.decode([AnyCodable].self) {
            value = arrayValue.map { $0.value }
        } else if let dictValue = try? container.decode([String: AnyCodable].self) {
            value = dictValue.mapValues { $0.value }
        } else {
            throw DecodingError.dataCorruptedError(in: container, debugDescription: "Unable to decode value")
        }
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()

        switch value {
        case let intValue as Int:
            try container.encode(intValue)
        case let doubleValue as Double:
            try container.encode(doubleValue)
        case let boolValue as Bool:
            try container.encode(boolValue)
        case let stringValue as String:
            try container.encode(stringValue)
        case let arrayValue as [Any]:
            let codableArray = arrayValue.map { AnyCodable($0) }
            try container.encode(codableArray)
        case let dictValue as [String: Any]:
            let codableDict = dictValue.mapValues { AnyCodable($0) }
            try container.encode(codableDict)
        default:
            throw EncodingError.invalidValue(value, EncodingError.Context(codingPath: encoder.codingPath, debugDescription: "Unable to encode value"))
        }
    }

    public static func == (lhs: AnyCodable, rhs: AnyCodable) -> Bool {
        switch (lhs.value, rhs.value) {
        case (let lhsValue as Int, let rhsValue as Int):
            return lhsValue == rhsValue
        case (let lhsValue as Double, let rhsValue as Double):
            return lhsValue == rhsValue
        case (let lhsValue as Bool, let rhsValue as Bool):
            return lhsValue == rhsValue
        case (let lhsValue as String, let rhsValue as String):
            return lhsValue == rhsValue
        default:
            return false
        }
    }
}