package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

class ConstantsTest {

    @Test
    fun testValidationConstants() {
        assertEquals(3, ValidationConstants.MIN_USERNAME_LENGTH)
        assertEquals(30, ValidationConstants.MAX_USERNAME_LENGTH)
        assertEquals(8, ValidationConstants.MIN_PASSWORD_LENGTH)
        assertEquals(100, ValidationConstants.MAX_PASSWORD_LENGTH)

        assertEquals("^\\+[1-9]\\d{1,14}$", ValidationConstants.PHONE_REGEX)
        assertEquals(10, ValidationConstants.MIN_PHONE_LENGTH)
        assertEquals(15, ValidationConstants.MAX_PHONE_LENGTH)

        assertEquals(2, ValidationConstants.MIN_BUSINESS_NAME_LENGTH)
        assertEquals(100, ValidationConstants.MAX_BUSINESS_NAME_LENGTH)
        assertEquals(10, ValidationConstants.MIN_BUSINESS_DESCRIPTION_LENGTH)
        assertEquals(1000, ValidationConstants.MAX_BUSINESS_DESCRIPTION_LENGTH)

        assertEquals(2, ValidationConstants.MIN_PRODUCT_NAME_LENGTH)
        assertEquals(200, ValidationConstants.MAX_PRODUCT_NAME_LENGTH)
        assertEquals(10, ValidationConstants.MIN_PRODUCT_DESCRIPTION_LENGTH)
        assertEquals(2000, ValidationConstants.MAX_PRODUCT_DESCRIPTION_LENGTH)
        assertEquals(0.01, ValidationConstants.MIN_PRODUCT_PRICE)
        assertEquals(999999.99, ValidationConstants.MAX_PRODUCT_PRICE)

        assertEquals(4000, ValidationConstants.MAX_MESSAGE_LENGTH)
        assertEquals(100, ValidationConstants.MAX_FILE_SIZE_MB)
        assertEquals("jpg,jpeg,png,gif,webp", ValidationConstants.SUPPORTED_IMAGE_FORMATS)
        assertEquals("mp4,mov,avi,mkv", ValidationConstants.SUPPORTED_VIDEO_FORMATS)
        assertEquals("mp3,wav,aac,ogg", ValidationConstants.SUPPORTED_AUDIO_FORMATS)
        assertEquals("pdf,doc,docx,xls,xlsx,ppt,pptx,txt", ValidationConstants.SUPPORTED_DOCUMENT_FORMATS)

        assertEquals(1, ValidationConstants.MIN_ORDER_QUANTITY)
        assertEquals(9999, ValidationConstants.MAX_ORDER_QUANTITY)
        assertEquals(100, ValidationConstants.MAX_ORDER_ITEMS)

        assertEquals(100, ValidationConstants.MAX_NOTIFICATION_TITLE_LENGTH)
        assertEquals(500, ValidationConstants.MAX_NOTIFICATION_BODY_LENGTH)

        assertEquals(20, ValidationConstants.DEFAULT_PAGE_SIZE)
        assertEquals(100, ValidationConstants.MAX_PAGE_SIZE)
        assertEquals(2, ValidationConstants.MIN_SEARCH_QUERY_LENGTH)
        assertEquals(100, ValidationConstants.MAX_SEARCH_QUERY_LENGTH)
    }

    @Test
    fun testRegionalConstants() {
        assertEquals("TH", RegionalConstants.DEFAULT_COUNTRY_CODE)
        assertEquals("en", RegionalConstants.DEFAULT_LANGUAGE)

        // Test supported countries
        assertEquals(6, RegionalConstants.SUPPORTED_COUNTRIES.size)
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("TH"))
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("SG"))
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("ID"))
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("MY"))
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("PH"))
        assertTrue(RegionalConstants.SUPPORTED_COUNTRIES.containsKey("VN"))

        // Test Thailand country info
        val thailand = RegionalConstants.SUPPORTED_COUNTRIES["TH"]
        assertNotNull(thailand)
        assertEquals("TH", thailand.code)
        assertEquals("Thailand", thailand.name)
        assertEquals("THB", thailand.currency)
        assertEquals("฿", thailand.currencySymbol)
        assertEquals("+66", thailand.phonePrefix)
        assertEquals("th_TH", thailand.locale)
        assertEquals("Asia/Bangkok", thailand.timezone)

        // Test Singapore country info
        val singapore = RegionalConstants.SUPPORTED_COUNTRIES["SG"]
        assertNotNull(singapore)
        assertEquals("SG", singapore.code)
        assertEquals("Singapore", singapore.name)
        assertEquals("SGD", singapore.currency)
        assertEquals("S$", singapore.currencySymbol)
        assertEquals("+65", singapore.phonePrefix)
        assertEquals("en_SG", singapore.locale)
        assertEquals("Asia/Singapore", singapore.timezone)

        // Test business hours
        assertEquals(6, RegionalConstants.BUSINESS_HOURS.size)
        val thBusinessHours = RegionalConstants.BUSINESS_HOURS["TH"]
        assertNotNull(thBusinessHours)
        assertEquals(9, thBusinessHours.start)
        assertEquals(18, thBusinessHours.end)

        val vnBusinessHours = RegionalConstants.BUSINESS_HOURS["VN"]
        assertNotNull(vnBusinessHours)
        assertEquals(8, vnBusinessHours.start)
        assertEquals(17, vnBusinessHours.end)

        // Test supported languages
        assertEquals(5, RegionalConstants.SUPPORTED_LANGUAGES.size)
        assertTrue(RegionalConstants.SUPPORTED_LANGUAGES.containsKey("en"))
        assertTrue(RegionalConstants.SUPPORTED_LANGUAGES.containsKey("th"))
        assertTrue(RegionalConstants.SUPPORTED_LANGUAGES.containsKey("id"))
        assertTrue(RegionalConstants.SUPPORTED_LANGUAGES.containsKey("ms"))
        assertTrue(RegionalConstants.SUPPORTED_LANGUAGES.containsKey("vi"))

        val englishLang = RegionalConstants.SUPPORTED_LANGUAGES["en"]
        assertNotNull(englishLang)
        assertEquals("en", englishLang.code)
        assertEquals("English", englishLang.name)
        assertEquals(listOf("SG", "PH"), englishLang.countries)

        val thaiLang = RegionalConstants.SUPPORTED_LANGUAGES["th"]
        assertNotNull(thaiLang)
        assertEquals("th", thaiLang.code)
        assertEquals("ไทย", thaiLang.name)
        assertEquals(listOf("TH"), thaiLang.countries)
    }

    @Test
    fun testAppConstants() {
        assertEquals(30, AppConstants.API_TIMEOUT_SECONDS)
        assertEquals(3, AppConstants.MAX_RETRY_ATTEMPTS)
        assertEquals(1000, AppConstants.RETRY_DELAY_MS)

        assertEquals(30, AppConstants.DEFAULT_CACHE_DURATION_MINUTES)
        assertEquals(24, AppConstants.IMAGE_CACHE_DURATION_HOURS)
        assertEquals(24, AppConstants.USER_SESSION_DURATION_HOURS)

        assertEquals(5000, AppConstants.WEBSOCKET_RECONNECT_INTERVAL_MS)
        assertEquals(10, AppConstants.WEBSOCKET_MAX_RECONNECT_ATTEMPTS)
        assertEquals(30000, AppConstants.HEARTBEAT_INTERVAL_MS)

        assertEquals(1024 * 1024, AppConstants.CHUNK_SIZE_BYTES) // 1MB
        assertEquals(3, AppConstants.MAX_CONCURRENT_UPLOADS)

        assertEquals(3000, AppConstants.SLOW_NETWORK_THRESHOLD_MS)
        assertEquals(0.8f, AppConstants.IMAGE_COMPRESSION_QUALITY)
        assertEquals(200, AppConstants.THUMBNAIL_MAX_SIZE_PX)

        assertEquals(15, AppConstants.TOKEN_REFRESH_THRESHOLD_MINUTES)
        assertEquals(5, AppConstants.MAX_LOGIN_ATTEMPTS)
        assertEquals(30, AppConstants.LOGIN_LOCKOUT_DURATION_MINUTES)
    }

    @Test
    fun testErrorCodes() {
        // Authentication errors
        assertEquals("AUTH_001", ErrorCodes.AUTH_INVALID_CREDENTIALS)
        assertEquals("AUTH_002", ErrorCodes.AUTH_TOKEN_EXPIRED)
        assertEquals("AUTH_003", ErrorCodes.AUTH_TOKEN_INVALID)
        assertEquals("AUTH_004", ErrorCodes.AUTH_USER_NOT_FOUND)
        assertEquals("AUTH_005", ErrorCodes.AUTH_USER_SUSPENDED)

        // Validation errors
        assertEquals("VAL_001", ErrorCodes.VALIDATION_REQUIRED_FIELD)
        assertEquals("VAL_002", ErrorCodes.VALIDATION_INVALID_FORMAT)
        assertEquals("VAL_003", ErrorCodes.VALIDATION_LENGTH_EXCEEDED)
        assertEquals("VAL_004", ErrorCodes.VALIDATION_INVALID_VALUE)

        // Business errors
        assertEquals("BIZ_001", ErrorCodes.BUSINESS_NOT_FOUND)
        assertEquals("BIZ_002", ErrorCodes.BUSINESS_NOT_VERIFIED)
        assertEquals("BIZ_003", ErrorCodes.BUSINESS_SUSPENDED)

        // Product errors
        assertEquals("PRD_001", ErrorCodes.PRODUCT_NOT_FOUND)
        assertEquals("PRD_002", ErrorCodes.PRODUCT_OUT_OF_STOCK)
        assertEquals("PRD_003", ErrorCodes.PRODUCT_UNAVAILABLE)

        // Order errors
        assertEquals("ORD_001", ErrorCodes.ORDER_NOT_FOUND)
        assertEquals("ORD_002", ErrorCodes.ORDER_CANNOT_CANCEL)
        assertEquals("ORD_003", ErrorCodes.ORDER_PAYMENT_FAILED)

        // Network errors
        assertEquals("NET_001", ErrorCodes.NETWORK_CONNECTION_ERROR)
        assertEquals("NET_002", ErrorCodes.NETWORK_TIMEOUT)
        assertEquals("NET_003", ErrorCodes.NETWORK_SERVER_ERROR)

        // File upload errors
        assertEquals("FILE_001", ErrorCodes.FILE_TOO_LARGE)
        assertEquals("FILE_002", ErrorCodes.FILE_INVALID_FORMAT)
        assertEquals("FILE_003", ErrorCodes.FILE_UPLOAD_FAILED)
    }

    @Test
    fun testCountryInfoDataClass() {
        val country = CountryInfo(
            code = "TH",
            name = "Thailand",
            currency = "THB",
            currencySymbol = "฿",
            phonePrefix = "+66",
            locale = "th_TH",
            timezone = "Asia/Bangkok"
        )

        assertEquals("TH", country.code)
        assertEquals("Thailand", country.name)
        assertEquals("THB", country.currency)
        assertEquals("฿", country.currencySymbol)
        assertEquals("+66", country.phonePrefix)
        assertEquals("th_TH", country.locale)
        assertEquals("Asia/Bangkok", country.timezone)
    }

    @Test
    fun testBusinessHoursDataClass() {
        val businessHours = BusinessHours(start = 9, end = 18)
        assertEquals(9, businessHours.start)
        assertEquals(18, businessHours.end)
    }

    @Test
    fun testLanguageInfoDataClass() {
        val language = LanguageInfo(
            code = "en",
            name = "English",
            countries = listOf("SG", "PH")
        )

        assertEquals("en", language.code)
        assertEquals("English", language.name)
        assertEquals(listOf("SG", "PH"), language.countries)
    }
}