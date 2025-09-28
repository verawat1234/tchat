package com.tchat.mobile.social.domain.models

import kotlinx.serialization.Serializable

/**
 * Southeast Asian Localization Models
 *
 * Comprehensive localization support for Southeast Asian markets with:
 * - Multi-language content support
 * - Cultural context awareness
 * - Regional feature preferences
 * - Time zone and formatting support
 */

@Serializable
data class LocalizationContext(
    val region: String = "TH",
    val language: String = "en",
    val timeZone: String = "Asia/Bangkok",
    val currency: String = "THB",
    val dateFormat: String = "dd/MM/yyyy",
    val timeFormat: String = "HH:mm"
)

@Serializable
data class RegionalConfig(
    val code: String,
    val name: String,
    val localName: String,
    val flag: String,
    val languages: List<String>,
    val primaryLanguage: String,
    val currency: String,
    val timeZone: String,
    val culturalFeatures: List<String>,
    val popularHashtags: List<String>,
    val greeting: String,
    val dateFormat: String
)

@Serializable
data class SocialStrings(
    val language: String,
    val region: String,
    val strings: Map<String, String>
)

// Southeast Asian Regional Configurations
object SEARegions {
    val regions = mapOf(
        "TH" to RegionalConfig(
            code = "TH",
            name = "Thailand",
            localName = "ประเทศไทย",
            flag = "🇹🇭",
            languages = listOf("th", "en"),
            primaryLanguage = "th",
            currency = "THB",
            timeZone = "Asia/Bangkok",
            culturalFeatures = listOf("buddhist_calendar", "royal_honorifics", "songkran_events"),
            popularHashtags = listOf("#Thailand", "#Bangkok", "#ThaiFood", "#Songkran", "#BTS", "#MRT", "#TomYum", "#PadThai"),
            greeting = "สวัสดี! What's happening?",
            dateFormat = "dd/MM/yyyy"
        ),
        "SG" to RegionalConfig(
            code = "SG",
            name = "Singapore",
            localName = "Singapore",
            flag = "🇸🇬",
            languages = listOf("en", "zh", "ms", "ta"),
            primaryLanguage = "en",
            currency = "SGD",
            timeZone = "Asia/Singapore",
            culturalFeatures = listOf("multicultural", "hawker_culture", "singlish"),
            popularHashtags = listOf("#Singapore", "#Merlion", "#HawkerCentre", "#SingaporeLife", "#MRT", "#Singlish", "#GardensByTheBay"),
            greeting = "Lah! What's happening?",
            dateFormat = "dd/MM/yyyy"
        ),
        "ID" to RegionalConfig(
            code = "ID",
            name = "Indonesia",
            localName = "Indonesia",
            flag = "🇮🇩",
            languages = listOf("id", "en", "jv"),
            primaryLanguage = "id",
            currency = "IDR",
            timeZone = "Asia/Jakarta",
            culturalFeatures = listOf("archipelago", "batik_patterns", "ramadan_events"),
            popularHashtags = listOf("#Indonesia", "#Jakarta", "#Bali", "#IndonesianCulture", "#RendangLife", "#Batik", "#WonderfulIndonesia"),
            greeting = "Halo! Apa kabar?",
            dateFormat = "dd/MM/yyyy"
        ),
        "MY" to RegionalConfig(
            code = "MY",
            name = "Malaysia",
            localName = "Malaysia",
            flag = "🇲🇾",
            languages = listOf("ms", "en", "zh", "ta"),
            primaryLanguage = "ms",
            currency = "MYR",
            timeZone = "Asia/Kuala_Lumpur",
            culturalFeatures = listOf("multicultural", "mamak_culture", "unity_in_diversity"),
            popularHashtags = listOf("#Malaysia", "#KualaLumpur", "#MalaysianFood", "#TrulyAsia", "#Mamak", "#Durian", "#KLCC"),
            greeting = "Apa khabar! What's happening?",
            dateFormat = "dd/MM/yyyy"
        ),
        "PH" to RegionalConfig(
            code = "PH",
            name = "Philippines",
            localName = "Pilipinas",
            flag = "🇵🇭",
            languages = listOf("en", "tl", "ceb"),
            primaryLanguage = "en",
            currency = "PHP",
            timeZone = "Asia/Manila",
            culturalFeatures = listOf("island_culture", "jeepney_culture", "bayanihan_spirit"),
            popularHashtags = listOf("#Philippines", "#Manila", "#Pinoy", "#Adobo", "#IslandLife", "#Jeepney", "#PinoyPride"),
            greeting = "Kumusta! What's happening?",
            dateFormat = "MM/dd/yyyy"
        ),
        "VN" to RegionalConfig(
            code = "VN",
            name = "Vietnam",
            localName = "Việt Nam",
            flag = "🇻🇳",
            languages = listOf("vi", "en"),
            primaryLanguage = "vi",
            currency = "VND",
            timeZone = "Asia/Ho_Chi_Minh",
            culturalFeatures = listOf("motorbike_culture", "pho_culture", "tet_celebrations"),
            popularHashtags = listOf("#Vietnam", "#Hanoi", "#HoChiMinh", "#Pho", "#Vietnamese", "#Motorbike", "#BanhMi"),
            greeting = "Xin chào! What's happening?",
            dateFormat = "dd/MM/yyyy"
        )
    )

    fun getRegion(code: String): RegionalConfig? = regions[code]

    fun getAllRegions(): List<RegionalConfig> = regions.values.toList()

    fun getRegionalHashtags(regionCode: String): List<String> {
        return getRegion(regionCode)?.popularHashtags ?: emptyList()
    }

    fun getRegionalGreeting(regionCode: String): String {
        return getRegion(regionCode)?.greeting ?: "Hello! What's happening?"
    }

    fun getCulturalFeatures(regionCode: String): List<String> {
        return getRegion(regionCode)?.culturalFeatures ?: emptyList()
    }
}

// Social Content Localization
object SocialLocalization {
    private val englishStrings = mapOf(
        "whats_happening" to "What's happening?",
        "create_post" to "Create Post",
        "follow" to "Follow",
        "following" to "Following",
        "like" to "Like",
        "comment" to "Comment",
        "share" to "Share",
        "bookmark" to "Bookmark",
        "discover_people" to "Discover People",
        "trending_now" to "Trending Now",
        "regional_content" to "Regional Content",
        "cultural_events" to "Cultural Events",
        "local_food" to "Local Food",
        "travel_spots" to "Travel Spots",
        "post_created" to "Post created successfully",
        "sync_complete" to "Sync completed",
        "offline_mode" to "Offline mode",
        "connection_restored" to "Connection restored"
    )

    private val thaiStrings = mapOf(
        "whats_happening" to "เกิดอะไรขึ้น?",
        "create_post" to "สร้างโพสต์",
        "follow" to "ติดตาม",
        "following" to "กำลังติดตาม",
        "like" to "ถูกใจ",
        "comment" to "แสดงความคิดเห็น",
        "share" to "แชร์",
        "bookmark" to "บุ๊คมาร์ค",
        "discover_people" to "ค้นหาผู้คน",
        "trending_now" to "กำลังเทรนด์",
        "regional_content" to "เนื้อหาในพื้นที่",
        "cultural_events" to "กิจกรรมทางวัฒนธรรม",
        "local_food" to "อาหารท้องถิ่น",
        "travel_spots" to "สถานที่ท่องเที่ยว",
        "post_created" to "สร้างโพสต์สำเร็จแล้ว",
        "sync_complete" to "ซิงค์เสร็จสิ้น",
        "offline_mode" to "โหมดออฟไลน์",
        "connection_restored" to "เชื่อมต่ออินเทอร์เน็ตแล้ว"
    )

    private val indonesianStrings = mapOf(
        "whats_happening" to "Apa yang terjadi?",
        "create_post" to "Buat Postingan",
        "follow" to "Ikuti",
        "following" to "Mengikuti",
        "like" to "Suka",
        "comment" to "Komentar",
        "share" to "Bagikan",
        "bookmark" to "Tandai",
        "discover_people" to "Temukan Orang",
        "trending_now" to "Sedang Trending",
        "regional_content" to "Konten Regional",
        "cultural_events" to "Acara Budaya",
        "local_food" to "Makanan Lokal",
        "travel_spots" to "Tempat Wisata",
        "post_created" to "Postingan berhasil dibuat",
        "sync_complete" to "Sinkronisasi selesai",
        "offline_mode" to "Mode offline",
        "connection_restored" to "Koneksi dipulihkan"
    )

    fun getString(key: String, language: String = "en", region: String = "TH"): String {
        val strings = when (language) {
            "th" -> thaiStrings
            "id" -> indonesianStrings
            else -> englishStrings
        }
        return strings[key] ?: englishStrings[key] ?: key
    }

    fun getAllStrings(language: String = "en"): Map<String, String> {
        return when (language) {
            "th" -> thaiStrings
            "id" -> indonesianStrings
            else -> englishStrings
        }
    }
}

// Time and Date Formatting for Southeast Asia
object SEATimeFormatting {
    fun formatTimeAgo(timestamp: String, regionCode: String = "TH"): String {
        // Simplified implementation - in real app use proper date libraries
        val config = SEARegions.getRegion(regionCode)
        return when (regionCode) {
            "TH" -> "2 นาทีที่แล้ว"
            "ID" -> "2 menit yang lalu"
            "VN" -> "2 phút trước"
            "MY" -> "2 minit yang lalu"
            else -> "2m ago"
        }
    }

    fun formatDate(timestamp: String, regionCode: String = "TH"): String {
        val config = SEARegions.getRegion(regionCode)
        val format = config?.dateFormat ?: "dd/MM/yyyy"
        // Use the regional date format
        return "28/09/2025" // Simplified
    }
}