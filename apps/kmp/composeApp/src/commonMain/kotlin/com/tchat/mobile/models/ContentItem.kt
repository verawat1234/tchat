package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import kotlinx.datetime.*

/**
 * T032: ContentItem model for dynamic CMS content
 *
 * Dynamic CMS content management with rich formatting, version control, and publishing workflow.
 * Multi-language support with SEO metadata and content type definitions.
 * Supports text, rich_text, image, and config content types with comprehensive validation.
 */
@Serializable
data class ContentItem(
    val id: String,
    val title: String,
    val slug: String,
    val type: ContentType,
    val status: ContentStatus = ContentStatus.DRAFT,
    val content: ContentData,
    val summary: String? = null,
    val excerpt: String? = null,
    val featuredImage: ContentImage? = null,
    val author: ContentAuthor,
    val category: ContentCategory? = null,
    val tags: List<String> = emptyList(),
    val languages: List<ContentLanguage> = emptyList(),
    val defaultLanguage: String = "en",
    val currentLanguage: String = "en",
    val seo: ContentSEO = ContentSEO(),
    val publishing: PublishingSettings = PublishingSettings(),
    val versioning: ContentVersioning = ContentVersioning(),
    val analytics: ContentAnalytics = ContentAnalytics(),
    val moderation: ContentModeration = ContentModeration(),
    val comments: ContentComments = ContentComments(),
    val sharing: SharingSettings = SharingSettings(),
    val accessibility: ContentAccessibilitySettings = ContentAccessibilitySettings(),
    val syndication: SyndicationSettings = SyndicationSettings(),
    val createdAt: String, // ISO 8601 timestamp
    val updatedAt: String,
    val publishedAt: String? = null,
    val lastModifiedAt: String? = null,
    val scheduledAt: String? = null,
    val archivedAt: String? = null,
    val expiresAt: String? = null,
    val metadata: Map<String, String> = emptyMap()
)

enum class ContentType {
    TEXT, // Plain text content
    RICH_TEXT, // HTML/Markdown with formatting
    IMAGE, // Image-based content
    CONFIG, // Configuration/settings content
    ARTICLE, // Long-form article
    PAGE, // Static page content
    POST, // Blog post
    NEWS, // News article
    PRODUCT, // Product content
    LANDING_PAGE, // Marketing landing page
    EMAIL_TEMPLATE, // Email template
    NOTIFICATION, // Push notification content
    WIDGET, // UI widget content
    COMPONENT // Reusable component content
}

enum class ContentStatus {
    DRAFT, // Work in progress
    IN_REVIEW, // Pending approval
    APPROVED, // Approved for publishing
    PUBLISHED, // Live content
    SCHEDULED, // Scheduled for future publishing
    ARCHIVED, // Archived content
    DELETED, // Soft deleted
    EXPIRED, // Past expiration date
    REJECTED, // Review rejected
    NEEDS_UPDATE // Requires updates
}

@Serializable
sealed class ContentData {
    @Serializable
    data class TextContent(
        val text: String,
        val formatting: TextFormatting? = null,
        val wordCount: Int = 0,
        val readingTime: Int = 0, // minutes
        val characterCount: Int = 0,
        val language: String = "en"
    ) : ContentData()

    @Serializable
    data class RichTextContent(
        val html: String,
        val markdown: String? = null,
        val plainText: String,
        val blocks: List<ContentBlock> = emptyList(),
        val wordCount: Int = 0,
        val readingTime: Int = 0,
        val characterCount: Int = 0,
        val inlineAssets: List<InlineAsset> = emptyList(),
        val headings: List<ContentHeading> = emptyList(),
        val tableOfContents: List<TOCEntry> = emptyList()
    ) : ContentData()

    @Serializable
    data class ImageContent(
        val images: List<ContentImage>,
        val caption: String? = null,
        val altText: String? = null,
        val gallery: ImageGallery? = null,
        val credits: String? = null,
        val license: String? = null
    ) : ContentData()

    @Serializable
    data class ConfigContent(
        val configuration: Map<String, ConfigValue>,
        val schema: ConfigSchema? = null,
        val environment: String = "production",
        val version: String = "1.0",
        val validationRules: List<ValidationRule> = emptyList()
    ) : ContentData()

    @Serializable
    data class MixedContent(
        val sections: List<ContentSection>,
        val layout: ContentLayout = ContentLayout.VERTICAL,
        val template: String? = null
    ) : ContentData()
}

@Serializable
data class ContentBlock(
    val id: String,
    val type: BlockType,
    val content: String,
    val attributes: Map<String, String> = emptyMap(),
    val order: Int = 0,
    val nested: List<ContentBlock> = emptyList()
)

enum class BlockType {
    PARAGRAPH, HEADING, LIST, QUOTE, CODE, IMAGE, VIDEO, AUDIO,
    EMBED, TABLE, DIVIDER, BUTTON, FORM, GALLERY, CAROUSEL,
    ACCORDION, TABS, CALLOUT, DOWNLOAD, SOCIAL_EMBED
}

@Serializable
data class InlineAsset(
    val id: String,
    val type: AssetType,
    val url: String,
    val name: String,
    val size: Long? = null,
    val mimeType: String? = null,
    val position: Int,
    val caption: String? = null
)

enum class AssetType {
    IMAGE, VIDEO, AUDIO, DOCUMENT, ARCHIVE, OTHER
}

@Serializable
data class ContentHeading(
    val level: Int, // 1-6 for h1-h6
    val text: String,
    val id: String? = null,
    val anchor: String? = null,
    val position: Int
)

@Serializable
data class TOCEntry(
    val level: Int,
    val title: String,
    val anchor: String,
    val children: List<TOCEntry> = emptyList()
)

@Serializable
data class ContentImage(
    val id: String,
    val url: String,
    val thumbnailUrl: String? = null,
    val alt: String,
    val caption: String? = null,
    val width: Int? = null,
    val height: Int? = null,
    val fileSize: Long? = null,
    val mimeType: String? = null,
    val focal: FocalPoint? = null,
    val variants: Map<String, String> = emptyMap(), // size -> url
    val metadata: ImageMetadata = ImageMetadata()
)

@Serializable
data class FocalPoint(
    val x: Float, // 0-1 normalized
    val y: Float  // 0-1 normalized
)

@Serializable
data class ImageMetadata(
    val exif: Map<String, String> = emptyMap(),
    val colorPalette: List<String> = emptyList(),
    val dominantColor: String? = null,
    val faces: List<FaceDetection> = emptyList(),
    val objects: List<String> = emptyList(),
    val textContent: String? = null // OCR extracted text
)

@Serializable
data class FaceDetection(
    val x: Float,
    val y: Float,
    val width: Float,
    val height: Float,
    val confidence: Float
)

@Serializable
data class ImageGallery(
    val layout: GalleryLayout = GalleryLayout.GRID,
    val columns: Int = 3,
    val spacing: Int = 10,
    val showCaptions: Boolean = true,
    val showThumbnails: Boolean = true,
    val allowZoom: Boolean = true,
    val autoPlay: Boolean = false,
    val transitionEffect: String? = null
)

enum class GalleryLayout {
    GRID, MASONRY, CAROUSEL, LIGHTBOX, STRIP
}

@Serializable
sealed class ConfigValue {
    @Serializable
    data class StringValue(val value: String) : ConfigValue()

    @Serializable
    data class NumberValue(val value: Double) : ConfigValue()

    @Serializable
    data class BooleanValue(val value: Boolean) : ConfigValue()

    @Serializable
    data class ArrayValue(val value: List<String>) : ConfigValue()

    @Serializable
    data class ObjectValue(val value: Map<String, String>) : ConfigValue()
}

@Serializable
data class ConfigSchema(
    val version: String,
    val properties: Map<String, PropertySchema>,
    val required: List<String> = emptyList()
)

@Serializable
data class PropertySchema(
    val type: String, // "string", "number", "boolean", "array", "object"
    val description: String? = null,
    val defaultValue: String? = null,
    val enum: List<String>? = null,
    val minimum: Double? = null,
    val maximum: Double? = null,
    val pattern: String? = null
)

@Serializable
data class ValidationRule(
    val property: String,
    val rule: String, // "required", "min", "max", "pattern", etc.
    val value: String? = null,
    val message: String? = null
)

@Serializable
data class ContentSection(
    val id: String,
    val type: SectionType,
    val title: String? = null,
    val content: ContentData,
    val order: Int = 0,
    val visible: Boolean = true,
    val settings: Map<String, String> = emptyMap()
)

enum class SectionType {
    HEADER, CONTENT, SIDEBAR, FOOTER, HERO, FEATURES,
    TESTIMONIALS, FAQ, PRICING, CONTACT, GALLERY, VIDEO
}

enum class ContentLayout {
    VERTICAL, HORIZONTAL, GRID, MASONRY, CARDS, TIMELINE
}

@Serializable
data class ContentAuthor(
    val id: String,
    val name: String,
    val email: String? = null,
    val avatar: String? = null,
    val bio: String? = null,
    val role: AuthorRole = AuthorRole.EDITOR,
    val permissions: AuthorPermissions = AuthorPermissions()
)

enum class AuthorRole {
    ADMIN, EDITOR, AUTHOR, CONTRIBUTOR, REVIEWER, VIEWER
}

@Serializable
data class AuthorPermissions(
    val canCreate: Boolean = true,
    val canEdit: Boolean = true,
    val canDelete: Boolean = false,
    val canPublish: Boolean = false,
    val canModerate: Boolean = false,
    val canViewAnalytics: Boolean = false
)

@Serializable
data class ContentCategory(
    val id: String,
    val name: String,
    val slug: String,
    val description: String? = null,
    val parentId: String? = null,
    val color: String? = null,
    val icon: String? = null,
    val order: Int = 0,
    val isActive: Boolean = true
)

@Serializable
data class ContentLanguage(
    val code: String, // ISO 639-1
    val name: String,
    val isDefault: Boolean = false,
    val isActive: Boolean = true,
    val completeness: Double = 1.0, // 0.0 to 1.0
    val lastUpdated: String? = null,
    val translator: String? = null
)

@Serializable
data class ContentSEO(
    val metaTitle: String? = null,
    val metaDescription: String? = null,
    val keywords: List<String> = emptyList(),
    val canonicalUrl: String? = null,
    val robots: String? = null, // "index,follow", "noindex,nofollow"
    val ogTitle: String? = null,
    val ogDescription: String? = null,
    val ogImage: String? = null,
    val ogType: String? = null,
    val twitterCard: String? = null,
    val twitterTitle: String? = null,
    val twitterDescription: String? = null,
    val twitterImage: String? = null,
    val schema: Map<String, String> = emptyMap(), // JSON-LD schema
    val redirects: List<Redirect> = emptyList()
)

@Serializable
data class Redirect(
    val from: String,
    val to: String,
    val type: Int = 301, // HTTP status code
    val isActive: Boolean = true
)

@Serializable
data class PublishingSettings(
    val autoPublish: Boolean = false,
    val publishImmediately: Boolean = false,
    val requiresApproval: Boolean = true,
    val approvers: List<String> = emptyList(), // User IDs
    val workflow: List<WorkflowStep> = emptyList(),
    val notifications: PublishingNotifications = PublishingNotifications(),
    val constraints: PublishingConstraints = PublishingConstraints()
)

@Serializable
data class WorkflowStep(
    val id: String,
    val name: String,
    val assignee: String? = null,
    val status: WorkflowStatus = WorkflowStatus.PENDING,
    val order: Int,
    val isRequired: Boolean = true,
    val deadline: String? = null,
    val notes: String? = null,
    val completedAt: String? = null
)

enum class WorkflowStatus {
    PENDING, IN_PROGRESS, COMPLETED, REJECTED, SKIPPED
}

@Serializable
data class PublishingNotifications(
    val onSubmit: Boolean = true,
    val onApproval: Boolean = true,
    val onPublish: Boolean = true,
    val onReject: Boolean = true,
    val recipients: List<String> = emptyList()
)

@Serializable
data class PublishingConstraints(
    val allowedHours: List<Int> = emptyList(), // 0-23
    val allowedDays: List<Int> = emptyList(), // 1-7
    val blackoutDates: List<String> = emptyList(),
    val minimumInterval: Int? = null, // Minutes between publications
    val requiresTag: Boolean = false,
    val requiredTags: List<String> = emptyList()
)

@Serializable
data class ContentVersioning(
    val currentVersion: String = "1.0",
    val versions: List<ContentVersion> = emptyList(),
    val isVersioningEnabled: Boolean = true,
    val maxVersions: Int = 10,
    val autoSave: Boolean = true,
    val saveInterval: Int = 300 // seconds
)

@Serializable
data class ContentVersion(
    val version: String,
    val authorId: String,
    val authorName: String,
    val changes: String? = null,
    val changesSummary: List<String> = emptyList(),
    val isMajor: Boolean = false,
    val createdAt: String,
    val size: Int? = null, // Content size in characters
    val hash: String? = null // Content hash for change detection
)

@Serializable
data class ContentAnalytics(
    val views: Long = 0,
    val uniqueViews: Long = 0,
    val timeSpent: Long = 0, // Total seconds
    val averageTimeSpent: Double = 0.0, // Average seconds per session
    val bounceRate: Double = 0.0, // 0.0 to 1.0
    val engagement: ContentEngagement = ContentEngagement(),
    val social: SocialMetrics = SocialMetrics(),
    val search: SearchMetrics = SearchMetrics(),
    val geography: Map<String, Long> = emptyMap(),
    val devices: Map<String, Long> = emptyMap(),
    val referrers: Map<String, Long> = emptyMap(),
    val performance: ContentPerformance = ContentPerformance(),
    val lastUpdated: String? = null
)

@Serializable
data class ContentEngagement(
    val likes: Long = 0,
    val shares: Long = 0,
    val comments: Long = 0,
    val bookmarks: Long = 0,
    val downloads: Long = 0,
    val prints: Long = 0,
    val conversions: Long = 0,
    val subscriptions: Long = 0,
    val emailSignups: Long = 0,
    val clickThroughs: Long = 0
)

@Serializable
data class SocialMetrics(
    val totalShares: Long = 0,
    val platforms: Map<String, Long> = emptyMap(), // platform -> shares
    val mentions: Long = 0,
    val hashtags: Map<String, Long> = emptyMap(),
    val viralCoefficient: Double = 0.0
)

@Serializable
data class SearchMetrics(
    val organicTraffic: Long = 0,
    val keywords: Map<String, Long> = emptyMap(), // keyword -> visits
    val rankings: Map<String, Int> = emptyMap(), // keyword -> position
    val clickThroughRate: Double = 0.0,
    val impressions: Long = 0,
    val averagePosition: Double = 0.0
)

@Serializable
data class ContentPerformance(
    val score: Double = 0.0, // 0.0 to 100.0
    val readabilityScore: Double = 0.0,
    val seoScore: Double = 0.0,
    val engagementScore: Double = 0.0,
    val qualityScore: Double = 0.0,
    val loadTime: Double = 0.0, // seconds
    val mobileOptimized: Boolean = false
)

@Serializable
data class ContentModeration(
    val isModerated: Boolean = false,
    val moderatedBy: String? = null,
    val moderatedAt: String? = null,
    val flagged: Boolean = false,
    val flagReasons: List<String> = emptyList(),
    val autoModerationScore: Double = 0.0, // 0.0 to 1.0
    val humanReviewRequired: Boolean = false,
    val warnings: List<String> = emptyList(),
    val violations: List<PolicyViolation> = emptyList()
)

@Serializable
data class PolicyViolation(
    val policy: String,
    val severity: ViolationSeverity,
    val description: String,
    val detectedAt: String,
    val resolved: Boolean = false,
    val resolvedAt: String? = null
)

enum class ViolationSeverity {
    LOW, MEDIUM, HIGH, CRITICAL
}

@Serializable
data class ContentComments(
    val enabled: Boolean = true,
    val requiresApproval: Boolean = false,
    val count: Long = 0,
    val averageRating: Double = 0.0,
    val allowReplies: Boolean = true,
    val allowAnonymous: Boolean = false,
    val moderationSettings: CommentModeration = CommentModeration()
)

@Serializable
data class CommentModeration(
    val autoModeration: Boolean = true,
    val spamFilter: Boolean = true,
    val profanityFilter: Boolean = true,
    val bannedWords: List<String> = emptyList(),
    val requiresLogin: Boolean = false,
    val rateLimit: Int = 10 // Comments per minute
)

@Serializable
data class SharingSettings(
    val allowSharing: Boolean = true,
    val platforms: List<String> = emptyList(), // Enabled platforms
    val shareButtons: Boolean = true,
    val shareText: String? = null,
    val shareImage: String? = null,
    val trackShares: Boolean = true,
    val customShareUrl: String? = null
)

@Serializable
data class ContentAccessibilitySettings(
    val altTextRequired: Boolean = true,
    val highContrast: Boolean = false,
    val fontSize: String = "medium", // small, medium, large
    val screenReader: Boolean = true,
    val keyboardNavigation: Boolean = true,
    val ariaLabels: Map<String, String> = emptyMap(),
    val colorBlindFriendly: Boolean = false
)

@Serializable
data class SyndicationSettings(
    val allowSyndication: Boolean = false,
    val rssIncluded: Boolean = true,
    val syndicationPartners: List<String> = emptyList(),
    val licenseType: String? = null,
    val attribution: String? = null,
    val canonicalSource: Boolean = true
)

/**
 * ContentItem utilities and extensions
 */
fun ContentItem.isPublished(): Boolean = status == ContentStatus.PUBLISHED

fun ContentItem.isDraft(): Boolean = status == ContentStatus.DRAFT

fun ContentItem.isScheduled(): Boolean = status == ContentStatus.SCHEDULED

fun ContentItem.canBeEdited(): Boolean {
    return status in listOf(ContentStatus.DRAFT, ContentStatus.IN_REVIEW, ContentStatus.REJECTED)
}

fun ContentItem.isExpired(): Boolean {
    return expiresAt?.let { expiry ->
        try {
            val expiryTime = Instant.parse(expiry)
            Clock.System.now() > expiryTime
        } catch (e: Exception) {
            false
        }
    } ?: false
}

fun ContentItem.getReadingTime(): Int {
    return when (content) {
        is ContentData.TextContent -> content.readingTime
        is ContentData.RichTextContent -> content.readingTime
        else -> 0
    }
}

fun ContentItem.getWordCount(): Int {
    return when (content) {
        is ContentData.TextContent -> content.wordCount
        is ContentData.RichTextContent -> content.wordCount
        else -> 0
    }
}

fun ContentItem.getPlainText(): String {
    return when (content) {
        is ContentData.TextContent -> content.text
        is ContentData.RichTextContent -> content.plainText
        is ContentData.ImageContent -> content.caption ?: ""
        is ContentData.ConfigContent -> ""
        is ContentData.MixedContent -> {
            content.sections.joinToString(" ") {
                when (it.content) {
                    is ContentData.TextContent -> it.content.text
                    is ContentData.RichTextContent -> it.content.plainText
                    else -> ""
                }
            }
        }
    }
}

fun ContentItem.hasMultipleLanguages(): Boolean = languages.size > 1

fun ContentItem.getAvailableLanguages(): List<String> = languages.map { it.code }

fun ContentItem.getLanguage(code: String): ContentLanguage? {
    return languages.find { it.code == code }
}

fun ContentItem.isTranslationComplete(languageCode: String): Boolean {
    return getLanguage(languageCode)?.completeness == 1.0
}

fun ContentItem.getSEOScore(): Double = analytics.performance.seoScore

fun ContentItem.getEngagementRate(): Double {
    return if (analytics.views > 0) {
        val totalEngagement = analytics.engagement.likes + analytics.engagement.shares +
                             analytics.engagement.comments + analytics.engagement.bookmarks
        totalEngagement.toDouble() / analytics.views
    } else 0.0
}

fun ContentItem.isHighPerforming(): Boolean {
    return analytics.performance.score >= 80.0
}

fun ContentItem.needsUpdate(): Boolean {
    return status == ContentStatus.NEEDS_UPDATE ||
           moderation.humanReviewRequired ||
           moderation.violations.any { !it.resolved }
}

fun ContentItem.canPublish(userId: String): Boolean {
    val userCanPublish = author.id == userId && author.permissions.canPublish
    val isApproved = status == ContentStatus.APPROVED || !publishing.requiresApproval
    val isNotExpired = !isExpired()
    val hasNoViolations = moderation.violations.all { it.resolved }

    return userCanPublish && isApproved && isNotExpired && hasNoViolations
}

/**
 * ContentItem list utilities
 */
fun List<ContentItem>.filterByStatus(status: ContentStatus): List<ContentItem> {
    return filter { it.status == status }
}

fun List<ContentItem>.filterByType(type: ContentType): List<ContentItem> {
    return filter { it.type == type }
}

fun List<ContentItem>.filterByAuthor(authorId: String): List<ContentItem> {
    return filter { it.author.id == authorId }
}

fun List<ContentItem>.filterByCategory(categoryId: String): List<ContentItem> {
    return filter { it.category?.id == categoryId }
}

fun List<ContentItem>.filterByTag(tag: String): List<ContentItem> {
    return filter { it.tags.contains(tag) }
}

fun List<ContentItem>.filterByLanguage(languageCode: String): List<ContentItem> {
    return filter { it.currentLanguage == languageCode }
}

fun List<ContentItem>.filterPublished(): List<ContentItem> {
    return filter { it.isPublished() && !it.isExpired() }
}

fun List<ContentItem>.filterDrafts(): List<ContentItem> {
    return filter { it.isDraft() }
}

fun List<ContentItem>.searchContent(query: String): List<ContentItem> {
    if (query.isBlank()) return this

    val lowerQuery = query.lowercase()
    return filter { item ->
        item.title.lowercase().contains(lowerQuery) ||
        item.getPlainText().lowercase().contains(lowerQuery) ||
        item.tags.any { it.lowercase().contains(lowerQuery) } ||
        item.author.name.lowercase().contains(lowerQuery) ||
        item.summary?.lowercase()?.contains(lowerQuery) == true
    }
}

fun List<ContentItem>.sortByViews(descending: Boolean = true): List<ContentItem> {
    return if (descending) sortedByDescending { it.analytics.views }
           else sortedBy { it.analytics.views }
}

fun List<ContentItem>.sortByEngagement(descending: Boolean = true): List<ContentItem> {
    return if (descending) sortedByDescending { it.getEngagementRate() }
           else sortedBy { it.getEngagementRate() }
}

fun List<ContentItem>.sortByPerformance(descending: Boolean = true): List<ContentItem> {
    return if (descending) sortedByDescending { it.analytics.performance.score }
           else sortedBy { it.analytics.performance.score }
}

fun List<ContentItem>.sortByRecent(descending: Boolean = true): List<ContentItem> {
    return if (descending) sortedByDescending { it.publishedAt ?: it.updatedAt }
           else sortedBy { it.publishedAt ?: it.updatedAt }
}

fun List<ContentItem>.sortByTitle(ascending: Boolean = true): List<ContentItem> {
    return if (ascending) sortedBy { it.title }
           else sortedByDescending { it.title }
}

fun List<ContentItem>.getTopTags(limit: Int = 10): List<Pair<String, Int>> {
    return flatMap { it.tags }
        .groupBy { it }
        .mapValues { it.value.size }
        .toList()
        .sortedByDescending { it.second }
        .take(limit)
}