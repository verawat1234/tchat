package com.tchat.mobile.models

/**
 * Application constants including validation rules and regional settings
 * Based on backend validation constants and business rules
 */

/**
 * Validation constants for form inputs and data validation
 */
object ValidationConstants {
    // User validation
    const val MIN_USERNAME_LENGTH = 3
    const val MAX_USERNAME_LENGTH = 30
    const val MIN_PASSWORD_LENGTH = 8
    const val MAX_PASSWORD_LENGTH = 100

    // Phone number validation (international format)
    const val PHONE_REGEX = "^\\+[1-9]\\d{1,14}$"
    const val MIN_PHONE_LENGTH = 10
    const val MAX_PHONE_LENGTH = 15

    // Business validation
    const val MIN_BUSINESS_NAME_LENGTH = 2
    const val MAX_BUSINESS_NAME_LENGTH = 100
    const val MIN_BUSINESS_DESCRIPTION_LENGTH = 10
    const val MAX_BUSINESS_DESCRIPTION_LENGTH = 1000

    // Product validation
    const val MIN_PRODUCT_NAME_LENGTH = 2
    const val MAX_PRODUCT_NAME_LENGTH = 200
    const val MIN_PRODUCT_DESCRIPTION_LENGTH = 10
    const val MAX_PRODUCT_DESCRIPTION_LENGTH = 2000
    const val MIN_PRODUCT_PRICE = 0.01
    const val MAX_PRODUCT_PRICE = 999999.99

    // Message validation
    const val MAX_MESSAGE_LENGTH = 4000
    const val MAX_FILE_SIZE_MB = 100
    const val SUPPORTED_IMAGE_FORMATS = "jpg,jpeg,png,gif,webp"
    const val SUPPORTED_VIDEO_FORMATS = "mp4,mov,avi,mkv"
    const val SUPPORTED_AUDIO_FORMATS = "mp3,wav,aac,ogg"
    const val SUPPORTED_DOCUMENT_FORMATS = "pdf,doc,docx,xls,xlsx,ppt,pptx,txt"

    // Order validation
    const val MIN_ORDER_QUANTITY = 1
    const val MAX_ORDER_QUANTITY = 9999
    const val MAX_ORDER_ITEMS = 100

    // Notification validation
    const val MAX_NOTIFICATION_TITLE_LENGTH = 100
    const val MAX_NOTIFICATION_BODY_LENGTH = 500

    // Search and pagination
    const val DEFAULT_PAGE_SIZE = 20
    const val MAX_PAGE_SIZE = 100
    const val MIN_SEARCH_QUERY_LENGTH = 2
    const val MAX_SEARCH_QUERY_LENGTH = 100
}

/**
 * Regional constants for SEA (Southeast Asia) markets
 */
object RegionalConstants {
    // Supported countries with metadata
    val SUPPORTED_COUNTRIES = mapOf(
        "TH" to CountryInfo(
            code = "TH",
            name = "Thailand",
            currency = "THB",
            currencySymbol = "฿",
            phonePrefix = "+66",
            locale = "th_TH",
            timezone = "Asia/Bangkok"
        ),
        "SG" to CountryInfo(
            code = "SG",
            name = "Singapore",
            currency = "SGD",
            currencySymbol = "S$",
            phonePrefix = "+65",
            locale = "en_SG",
            timezone = "Asia/Singapore"
        ),
        "ID" to CountryInfo(
            code = "ID",
            name = "Indonesia",
            currency = "IDR",
            currencySymbol = "Rp",
            phonePrefix = "+62",
            locale = "id_ID",
            timezone = "Asia/Jakarta"
        ),
        "MY" to CountryInfo(
            code = "MY",
            name = "Malaysia",
            currency = "MYR",
            currencySymbol = "RM",
            phonePrefix = "+60",
            locale = "ms_MY",
            timezone = "Asia/Kuala_Lumpur"
        ),
        "PH" to CountryInfo(
            code = "PH",
            name = "Philippines",
            currency = "PHP",
            currencySymbol = "₱",
            phonePrefix = "+63",
            locale = "en_PH",
            timezone = "Asia/Manila"
        ),
        "VN" to CountryInfo(
            code = "VN",
            name = "Vietnam",
            currency = "VND",
            currencySymbol = "₫",
            phonePrefix = "+84",
            locale = "vi_VN",
            timezone = "Asia/Ho_Chi_Minh"
        )
    )

    // Default country (can be changed based on user location)
    const val DEFAULT_COUNTRY_CODE = "TH"

    // Regional business hours (24-hour format)
    val BUSINESS_HOURS = mapOf(
        "TH" to BusinessHours(start = 9, end = 18),
        "SG" to BusinessHours(start = 9, end = 18),
        "ID" to BusinessHours(start = 9, end = 17),
        "MY" to BusinessHours(start = 9, end = 18),
        "PH" to BusinessHours(start = 9, end = 18),
        "VN" to BusinessHours(start = 8, end = 17)
    )

    // Supported languages with their country associations
    val SUPPORTED_LANGUAGES = mapOf(
        "en" to LanguageInfo(code = "en", name = "English", countries = listOf("SG", "PH")),
        "th" to LanguageInfo(code = "th", name = "ไทย", countries = listOf("TH")),
        "id" to LanguageInfo(code = "id", name = "Bahasa Indonesia", countries = listOf("ID")),
        "ms" to LanguageInfo(code = "ms", name = "Bahasa Melayu", countries = listOf("MY")),
        "vi" to LanguageInfo(code = "vi", name = "Tiếng Việt", countries = listOf("VN"))
    )

    const val DEFAULT_LANGUAGE = "en"
}

/**
 * Application-level constants
 */
object AppConstants {
    // API configuration
    const val API_TIMEOUT_SECONDS = 30
    const val MAX_RETRY_ATTEMPTS = 3
    const val RETRY_DELAY_MS = 1000

    // Cache configuration
    const val DEFAULT_CACHE_DURATION_MINUTES = 30
    const val IMAGE_CACHE_DURATION_HOURS = 24
    const val USER_SESSION_DURATION_HOURS = 24

    // WebSocket configuration
    const val WEBSOCKET_RECONNECT_INTERVAL_MS = 5000
    const val WEBSOCKET_MAX_RECONNECT_ATTEMPTS = 10
    const val HEARTBEAT_INTERVAL_MS = 30000

    // File upload configuration
    const val CHUNK_SIZE_BYTES = 1024 * 1024 // 1MB chunks
    const val MAX_CONCURRENT_UPLOADS = 3

    // Performance thresholds
    const val SLOW_NETWORK_THRESHOLD_MS = 3000
    const val IMAGE_COMPRESSION_QUALITY = 0.8f
    const val THUMBNAIL_MAX_SIZE_PX = 200

    // Security configuration
    const val TOKEN_REFRESH_THRESHOLD_MINUTES = 15
    const val MAX_LOGIN_ATTEMPTS = 5
    const val LOGIN_LOCKOUT_DURATION_MINUTES = 30
}

/**
 * Data classes for regional information
 */
data class CountryInfo(
    val code: String,
    val name: String,
    val currency: String,
    val currencySymbol: String,
    val phonePrefix: String,
    val locale: String,
    val timezone: String
)

data class BusinessHours(
    val start: Int, // 24-hour format
    val end: Int    // 24-hour format
)

data class LanguageInfo(
    val code: String,
    val name: String,
    val countries: List<String>
)

/**
 * Error codes for consistent error handling
 */
object ErrorCodes {
    // Authentication errors
    const val AUTH_INVALID_CREDENTIALS = "AUTH_001"
    const val AUTH_TOKEN_EXPIRED = "AUTH_002"
    const val AUTH_TOKEN_INVALID = "AUTH_003"
    const val AUTH_USER_NOT_FOUND = "AUTH_004"
    const val AUTH_USER_SUSPENDED = "AUTH_005"

    // Validation errors
    const val VALIDATION_REQUIRED_FIELD = "VAL_001"
    const val VALIDATION_INVALID_FORMAT = "VAL_002"
    const val VALIDATION_LENGTH_EXCEEDED = "VAL_003"
    const val VALIDATION_INVALID_VALUE = "VAL_004"

    // Business errors
    const val BUSINESS_NOT_FOUND = "BIZ_001"
    const val BUSINESS_NOT_VERIFIED = "BIZ_002"
    const val BUSINESS_SUSPENDED = "BIZ_003"

    // Product errors
    const val PRODUCT_NOT_FOUND = "PRD_001"
    const val PRODUCT_OUT_OF_STOCK = "PRD_002"
    const val PRODUCT_UNAVAILABLE = "PRD_003"

    // Order errors
    const val ORDER_NOT_FOUND = "ORD_001"
    const val ORDER_CANNOT_CANCEL = "ORD_002"
    const val ORDER_PAYMENT_FAILED = "ORD_003"

    // Network errors
    const val NETWORK_CONNECTION_ERROR = "NET_001"
    const val NETWORK_TIMEOUT = "NET_002"
    const val NETWORK_SERVER_ERROR = "NET_003"

    // File upload errors
    const val FILE_TOO_LARGE = "FILE_001"
    const val FILE_INVALID_FORMAT = "FILE_002"
    const val FILE_UPLOAD_FAILED = "FILE_003"
}