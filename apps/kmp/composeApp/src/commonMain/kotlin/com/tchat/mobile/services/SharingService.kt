package com.tchat.mobile.services

import com.tchat.mobile.models.Post

/**
 * Sharing Service - Platform Integration
 *
 * Handles sharing to various social media platforms and messaging apps
 * Mock implementation ready for real platform SDK integration
 */

data class ShareResult(
    val success: Boolean,
    val platform: String,
    val message: String? = null,
    val sharedUrl: String? = null
)

enum class SharingPlatform(val displayName: String, val packageName: String) {
    // Social Media
    TWITTER("Twitter", "com.twitter.android"),
    FACEBOOK("Facebook", "com.facebook.katana"),
    INSTAGRAM("Instagram", "com.instagram.android"),
    TIKTOK("TikTok", "com.ss.android.ugc.trill"),
    LINKEDIN("LinkedIn", "com.linkedin.android"),
    SNAPCHAT("Snapchat", "com.snapchat.android"),

    // Messaging
    WHATSAPP("WhatsApp", "com.whatsapp"),
    LINE("LINE", "jp.naver.line.android"),
    TELEGRAM("Telegram", "org.telegram.messenger"),
    WECHAT("WeChat", "com.tencent.mm"),
    MESSENGER("Messenger", "com.facebook.orca"),

    // System
    SMS("SMS", "com.android.messaging"),
    EMAIL("Email", "com.android.email"),
    COPY_LINK("Copy Link", "system.clipboard"),
    MORE("More", "system.share")
}

interface SharingService {
    suspend fun sharePost(post: Post, platform: SharingPlatform): ShareResult
    suspend fun shareText(text: String, platform: SharingPlatform): ShareResult
    suspend fun shareImage(imageUrl: String, caption: String?, platform: SharingPlatform): ShareResult
    suspend fun shareVideo(videoUrl: String, caption: String?, platform: SharingPlatform): ShareResult
    suspend fun getAvailablePlatforms(): List<SharingPlatform>
    suspend fun isPlatformAvailable(platform: SharingPlatform): Boolean
}

class MockSharingService : SharingService {

    override suspend fun sharePost(post: Post, platform: SharingPlatform): ShareResult {
        // Simulate platform-specific post sharing
        val shareText = buildShareText(post)
        val shareUrl = "https://tchat.app/posts/${post.id}"

        return when (platform) {
            SharingPlatform.TWITTER -> shareToTwitter(shareText, shareUrl, post)
            SharingPlatform.FACEBOOK -> shareToFacebook(shareText, shareUrl, post)
            SharingPlatform.INSTAGRAM -> shareToInstagram(post)
            SharingPlatform.TIKTOK -> shareToTikTok(post)
            SharingPlatform.WHATSAPP -> shareToWhatsApp(shareText, shareUrl)
            SharingPlatform.LINE -> shareToLine(shareText, shareUrl)
            SharingPlatform.TELEGRAM -> shareToTelegram(shareText, shareUrl)
            SharingPlatform.SMS -> shareToSMS(shareText, shareUrl)
            SharingPlatform.EMAIL -> shareToEmail(shareText, shareUrl, post)
            SharingPlatform.COPY_LINK -> copyToClipboard(shareUrl)
            else -> shareGeneric(shareText, shareUrl, platform)
        }
    }

    override suspend fun shareText(text: String, platform: SharingPlatform): ShareResult {
        return when (platform) {
            SharingPlatform.TWITTER -> {
                val truncatedText = if (text.length > 280) text.take(277) + "..." else text
                ShareResult(true, platform.displayName, "Shared to Twitter", null)
            }
            SharingPlatform.WHATSAPP -> {
                ShareResult(true, platform.displayName, "Shared to WhatsApp", null)
            }
            SharingPlatform.COPY_LINK -> {
                // Mock clipboard copy
                ShareResult(true, platform.displayName, "Copied to clipboard", text)
            }
            else -> {
                ShareResult(true, platform.displayName, "Shared via ${platform.displayName}", null)
            }
        }
    }

    override suspend fun shareImage(imageUrl: String, caption: String?, platform: SharingPlatform): ShareResult {
        return when (platform) {
            SharingPlatform.INSTAGRAM -> {
                // Mock Instagram Stories/Feed sharing
                ShareResult(true, platform.displayName, "Shared to Instagram", imageUrl)
            }
            SharingPlatform.SNAPCHAT -> {
                ShareResult(true, platform.displayName, "Shared to Snapchat", imageUrl)
            }
            SharingPlatform.FACEBOOK -> {
                ShareResult(true, platform.displayName, "Shared to Facebook", imageUrl)
            }
            else -> {
                ShareResult(true, platform.displayName, "Shared image via ${platform.displayName}", imageUrl)
            }
        }
    }

    override suspend fun shareVideo(videoUrl: String, caption: String?, platform: SharingPlatform): ShareResult {
        return when (platform) {
            SharingPlatform.TIKTOK -> {
                ShareResult(true, platform.displayName, "Shared to TikTok", videoUrl)
            }
            SharingPlatform.INSTAGRAM -> {
                ShareResult(true, platform.displayName, "Shared to Instagram Reels", videoUrl)
            }
            SharingPlatform.SNAPCHAT -> {
                ShareResult(true, platform.displayName, "Shared to Snapchat", videoUrl)
            }
            else -> {
                ShareResult(true, platform.displayName, "Shared video via ${platform.displayName}", videoUrl)
            }
        }
    }

    override suspend fun getAvailablePlatforms(): List<SharingPlatform> {
        // Mock available platforms (in real app, would check installed apps)
        return listOf(
            SharingPlatform.WHATSAPP,
            SharingPlatform.LINE,
            SharingPlatform.TWITTER,
            SharingPlatform.FACEBOOK,
            SharingPlatform.INSTAGRAM,
            SharingPlatform.TELEGRAM,
            SharingPlatform.SMS,
            SharingPlatform.EMAIL,
            SharingPlatform.COPY_LINK,
            SharingPlatform.MORE
        )
    }

    override suspend fun isPlatformAvailable(platform: SharingPlatform): Boolean {
        // Mock platform availability check
        return when (platform) {
            SharingPlatform.WHATSAPP,
            SharingPlatform.LINE,
            SharingPlatform.TWITTER,
            SharingPlatform.FACEBOOK,
            SharingPlatform.INSTAGRAM,
            SharingPlatform.SMS,
            SharingPlatform.EMAIL,
            SharingPlatform.COPY_LINK -> true
            else -> false
        }
    }

    // Platform-specific sharing implementations
    private suspend fun shareToTwitter(text: String, url: String, post: Post): ShareResult {
        val tweetText = "$text $url ${post.content.hashtags.take(3).joinToString(" ")}"
        val truncatedTweet = if (tweetText.length > 280) tweetText.take(277) + "..." else tweetText

        // Mock Twitter API call
        println("Sharing to Twitter: $truncatedTweet")
        return ShareResult(true, "Twitter", "Shared to Twitter", url)
    }

    private suspend fun shareToFacebook(text: String, url: String, post: Post): ShareResult {
        // Mock Facebook Graph API call
        println("Sharing to Facebook: $text")
        return ShareResult(true, "Facebook", "Shared to Facebook", url)
    }

    private suspend fun shareToInstagram(post: Post): ShareResult {
        // Instagram sharing depends on content type
        return when {
            post.content.images.isNotEmpty() -> {
                println("Sharing image to Instagram Stories")
                ShareResult(true, "Instagram", "Shared to Instagram Stories", post.content.images.first().url)
            }
            post.content.videos.isNotEmpty() -> {
                println("Sharing video to Instagram Reels")
                ShareResult(true, "Instagram", "Shared to Instagram Reels", post.content.videos.first().url)
            }
            else -> {
                ShareResult(false, "Instagram", "Instagram requires media content", null)
            }
        }
    }

    private suspend fun shareToTikTok(post: Post): ShareResult {
        return if (post.content.videos.isNotEmpty()) {
            println("Sharing video to TikTok")
            ShareResult(true, "TikTok", "Shared to TikTok", post.content.videos.first().url)
        } else {
            ShareResult(false, "TikTok", "TikTok requires video content", null)
        }
    }

    private suspend fun shareToWhatsApp(text: String, url: String): ShareResult {
        val whatsappText = "$text\n\n$url"
        println("Sharing to WhatsApp: $whatsappText")
        return ShareResult(true, "WhatsApp", "Shared to WhatsApp", url)
    }

    private suspend fun shareToLine(text: String, url: String): ShareResult {
        val lineText = "$text\n$url"
        println("Sharing to LINE: $lineText")
        return ShareResult(true, "LINE", "Shared to LINE", url)
    }

    private suspend fun shareToTelegram(text: String, url: String): ShareResult {
        val telegramText = "$text\n$url"
        println("Sharing to Telegram: $telegramText")
        return ShareResult(true, "Telegram", "Shared to Telegram", url)
    }

    private suspend fun shareToSMS(text: String, url: String): ShareResult {
        val smsText = "$text $url"
        println("Sharing via SMS: $smsText")
        return ShareResult(true, "SMS", "SMS ready to send", url)
    }

    private suspend fun shareToEmail(text: String, url: String, post: Post): ShareResult {
        val subject = "Check out this post from ${post.user.displayName ?: post.user.username}"
        val body = """
            $text

            View the full post: $url

            Shared via Tchat
        """.trimIndent()

        println("Preparing email - Subject: $subject")
        return ShareResult(true, "Email", "Email ready to send", url)
    }

    private suspend fun copyToClipboard(url: String): ShareResult {
        // Mock clipboard operation
        println("Copied to clipboard: $url")
        return ShareResult(true, "Copy Link", "Link copied to clipboard", url)
    }

    private suspend fun shareGeneric(text: String, url: String, platform: SharingPlatform): ShareResult {
        println("Generic share to ${platform.displayName}: $text $url")
        return ShareResult(true, platform.displayName, "Shared via ${platform.displayName}", url)
    }

    private fun buildShareText(post: Post): String {
        val baseText = post.content.text ?: "Check out this post"
        val user = "@${post.user.username}"

        return when (post.type) {
            com.tchat.mobile.models.PostType.REVIEW -> {
                val rating = "⭐".repeat(((post.metadata?.rating ?: 0.8f) * 5).toInt())
                "Great review by $user! $rating\n\n$baseText"
            }
            com.tchat.mobile.models.PostType.SOCIAL -> {
                "$baseText\n\n- $user"
            }
            com.tchat.mobile.models.PostType.VIDEO -> {
                "Amazing video by $user! 🎥\n\n$baseText"
            }
            else -> {
                "$baseText\n\n- $user"
            }
        }
    }
}