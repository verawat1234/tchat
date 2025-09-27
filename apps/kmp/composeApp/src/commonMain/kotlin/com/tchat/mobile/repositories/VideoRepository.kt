package com.tchat.mobile.repositories

import com.tchat.mobile.models.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.datetime.Clock

/**
 * Video Repository for managing video content, channels, and interactions
 * Provides mock data for development and testing
 */
interface VideoRepository {
    suspend fun getShortVideos(category: VideoCategory = VideoCategory.ALL): List<VideoContent>
    suspend fun getLongVideos(category: VideoCategory = VideoCategory.ALL): List<VideoContent>
    suspend fun getChannels(category: VideoCategory = VideoCategory.ALL): List<ChannelInfo>
    suspend fun getChannel(channelId: String): ChannelInfo?
    suspend fun getVideo(videoId: String): VideoContent?

    // Engagement
    suspend fun likeVideo(videoId: String): Boolean
    suspend fun unlikeVideo(videoId: String): Boolean
    suspend fun subscribeToChannel(channelId: String): Boolean
    suspend fun unsubscribeFromChannel(channelId: String): Boolean

    // Observable
    fun observeVideos(): Flow<List<VideoContent>>
    fun observeSubscriptions(): Flow<Set<String>>
}

class MockVideoRepository : VideoRepository {

    // Mock in-memory state
    private var likedVideos = mutableSetOf<String>()
    private var subscribedChannels = mutableSetOf<String>()

    // Diverse video URLs from different sources
    private val videoUrls = listOf(
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/SubaruOutbackOnStreetAndDirt.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/VolkswagenGTIReview.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4",
        "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WhatCarCanYouGetForAGrand.mp4"
    )

    override suspend fun getShortVideos(category: VideoCategory): List<VideoContent> {
        delay(500) // Simulate network delay

        val allShorts = listOf(
            createVideoContent(
                id = "short-1",
                title = "Amazing Pad Thai Street Food",
                description = "Learn how to make authentic Pad Thai in 60 seconds! 🍜✨ #StreetFood #Thai #Cooking",
                thumbnail = "https://images.unsplash.com/photo-1559847844-5315695b6a77?w=400&h=700&fit=crop",
                videoUrl = videoUrls[0],
                duration = "0:45",
                durationSeconds = 45,
                views = 128000,
                likes = 8500,
                comments = 247,
                category = VideoCategory.FOOD,
                tags = listOf("food", "thai", "cooking", "street"),
                channelId = "thai-chef",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-2",
                title = "Thai Dance Challenge",
                description = "Traditional Thai dance moves gone viral! Try this at home 💃 #ThaiDance #Culture #Challenge",
                thumbnail = "https://images.unsplash.com/photo-1518611012118-696072aa579a?w=400&h=700&fit=crop",
                videoUrl = videoUrls[1],
                duration = "0:30",
                durationSeconds = 30,
                views = 245000,
                likes = 15200,
                comments = 832,
                category = VideoCategory.ENTERTAINMENT,
                tags = listOf("dance", "culture", "traditional", "viral"),
                channelId = "thai-culture",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-3",
                title = "Bangkok Night Market Tour",
                description = "Best night market finds under 100 baht! 🌃🛍️ #Bangkok #Market #Shopping",
                thumbnail = "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=400&h=700&fit=crop",
                videoUrl = videoUrls[2],
                duration = "1:20",
                durationSeconds = 80,
                views = 67000,
                likes = 4300,
                comments = 156,
                category = VideoCategory.TRAVEL,
                tags = listOf("bangkok", "market", "shopping", "budget"),
                channelId = "bangkok-explorer",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-4",
                title = "Korean BBQ Techniques",
                description = "Master Korean BBQ at home! 🥩🔥 Perfect grilling every time",
                thumbnail = "https://images.unsplash.com/photo-1529692236671-f1f6cf9683ba?w=400&h=700&fit=crop",
                videoUrl = videoUrls[3],
                duration = "0:55",
                durationSeconds = 55,
                views = 89000,
                likes = 6700,
                comments = 298,
                category = VideoCategory.FOOD,
                tags = listOf("korean", "bbq", "grilling", "cooking"),
                channelId = "asian-cuisine",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-5",
                title = "Tech Gadget Review",
                description = "Latest smartphone features tested! 📱⚡ Is it worth the upgrade?",
                thumbnail = "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=400&h=700&fit=crop",
                videoUrl = videoUrls[4],
                duration = "1:15",
                durationSeconds = 75,
                views = 156000,
                likes = 12400,
                comments = 445,
                category = VideoCategory.TECH,
                tags = listOf("tech", "review", "smartphone", "gadgets"),
                channelId = "tech-reviewer",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-6",
                title = "Gaming Tips & Tricks",
                description = "Pro gaming strategies that actually work! 🎮🏆 Level up your skills",
                thumbnail = "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=400&h=700&fit=crop",
                videoUrl = videoUrls[5],
                duration = "0:40",
                durationSeconds = 40,
                views = 203000,
                likes = 18500,
                comments = 672,
                category = VideoCategory.GAMING,
                tags = listOf("gaming", "tips", "strategy", "pro"),
                channelId = "pro-gamer",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-7",
                title = "Fitness Workout",
                description = "5-minute abs workout! No equipment needed 💪🔥 #Fitness #Workout",
                thumbnail = "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400&h=700&fit=crop",
                videoUrl = videoUrls[6],
                duration = "0:50",
                durationSeconds = 50,
                views = 324000,
                likes = 28900,
                comments = 1234,
                category = VideoCategory.LIFESTYLE,
                tags = listOf("fitness", "workout", "abs", "health"),
                channelId = "fitness-coach",
                type = VideoType.SHORT
            ),
            createVideoContent(
                id = "short-8",
                title = "Business Success Tips",
                description = "Entrepreneur mindset secrets! 💼✨ Transform your business today",
                thumbnail = "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=700&fit=crop",
                videoUrl = videoUrls[7],
                duration = "1:05",
                durationSeconds = 65,
                views = 98000,
                likes = 7800,
                comments = 234,
                category = VideoCategory.BUSINESS,
                tags = listOf("business", "entrepreneur", "success", "tips"),
                channelId = "business-mentor",
                type = VideoType.SHORT
            )
        )

        return if (category == VideoCategory.ALL) allShorts else allShorts.filter { it.category == category }
    }

    override suspend fun getLongVideos(category: VideoCategory): List<VideoContent> {
        delay(600) // Simulate network delay

        val allLongs = listOf(
            createVideoContent(
                id = "long-1",
                title = "Complete Guide to Thai Street Food: Bangkok Edition",
                description = "Join me as I explore the best street food vendors in Bangkok! From traditional Pad Thai to hidden gems, this comprehensive guide covers everything you need to know about Thai street food culture.",
                thumbnail = "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=400&h=225&fit=crop",
                videoUrl = videoUrls[8],
                duration = "15:43",
                durationSeconds = 943,
                views = 524000,
                likes = 18500,
                comments = 1842,
                category = VideoCategory.FOOD,
                tags = listOf("street food", "bangkok", "thai cuisine", "travel"),
                channelId = "foodie-adventures",
                type = VideoType.LONG
            ),
            createVideoContent(
                id = "long-2",
                title = "Thai Language Crash Course for Travelers",
                description = "Learn essential Thai phrases for your trip to Thailand! Perfect for beginners who want to communicate with locals and enhance their travel experience.",
                thumbnail = "https://images.unsplash.com/photo-1544717297-fa95b6ee9643?w=400&h=225&fit=crop",
                videoUrl = videoUrls[9],
                duration = "22:18",
                durationSeconds = 1338,
                views = 156000,
                likes = 9200,
                comments = 567,
                category = VideoCategory.EDUCATION,
                tags = listOf("thai language", "education", "travel", "learning"),
                channelId = "learn-thai",
                type = VideoType.LONG
            ),
            createVideoContent(
                id = "long-3",
                title = "The Future of Technology: AI Revolution",
                description = "Deep dive into how artificial intelligence is transforming our world. From automation to creativity, explore the opportunities and challenges ahead.",
                thumbnail = "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=400&h=225&fit=crop",
                videoUrl = videoUrls[10],
                duration = "28:45",
                durationSeconds = 1725,
                views = 412000,
                likes = 31200,
                comments = 2134,
                category = VideoCategory.TECH,
                tags = listOf("ai", "technology", "future", "innovation"),
                channelId = "tech-futurist",
                type = VideoType.LONG
            ),
            createVideoContent(
                id = "long-4",
                title = "Building a Successful Business from Zero",
                description = "Complete entrepreneur's journey from idea to million-dollar company. Real strategies, failures, and lessons learned along the way.",
                thumbnail = "https://images.unsplash.com/photo-1556761175-b413da4baf72?w=400&h=225&fit=crop",
                videoUrl = videoUrls[11],
                duration = "35:22",
                durationSeconds = 2122,
                views = 289000,
                likes = 24700,
                comments = 1456,
                category = VideoCategory.BUSINESS,
                tags = listOf("startup", "entrepreneur", "business", "success"),
                channelId = "startup-stories",
                type = VideoType.LONG
            ),
            createVideoContent(
                id = "long-5",
                title = "Complete Fitness Transformation Guide",
                description = "12-week fitness transformation program with nutrition, workouts, and mindset coaching. Real results from real people.",
                thumbnail = "https://images.unsplash.com/photo-1534258936925-c58bed479fcb?w=400&h=225&fit=crop",
                videoUrl = videoUrls[12],
                duration = "42:15",
                durationSeconds = 2535,
                views = 678000,
                likes = 45600,
                comments = 3421,
                category = VideoCategory.LIFESTYLE,
                tags = listOf("fitness", "transformation", "health", "nutrition"),
                channelId = "transformation-coach",
                type = VideoType.LONG
            )
        )

        return if (category == VideoCategory.ALL) allLongs else allLongs.filter { it.category == category }
    }

    override suspend fun getChannels(category: VideoCategory): List<ChannelInfo> {
        delay(400) // Simulate network delay

        val allChannels = listOf(
            createChannelInfo("thai-chef", "Bangkok Street Chef", VideoCategory.FOOD),
            createChannelInfo("thai-culture", "Thai Culture Hub", VideoCategory.ENTERTAINMENT),
            createChannelInfo("bangkok-explorer", "Bangkok Explorer", VideoCategory.TRAVEL),
            createChannelInfo("foodie-adventures", "Southeast Asian Foodie", VideoCategory.FOOD),
            createChannelInfo("learn-thai", "Thai Language Academy", VideoCategory.EDUCATION),
            createChannelInfo("asian-cuisine", "Asian Cuisine Master", VideoCategory.FOOD),
            createChannelInfo("tech-reviewer", "Tech Review Central", VideoCategory.TECH),
            createChannelInfo("pro-gamer", "Pro Gaming Academy", VideoCategory.GAMING),
            createChannelInfo("fitness-coach", "Fitness Revolution", VideoCategory.LIFESTYLE),
            createChannelInfo("business-mentor", "Business Growth Hub", VideoCategory.BUSINESS),
            createChannelInfo("tech-futurist", "Tech Futurist", VideoCategory.TECH),
            createChannelInfo("startup-stories", "Startup Success Stories", VideoCategory.BUSINESS),
            createChannelInfo("transformation-coach", "Total Life Transformation", VideoCategory.LIFESTYLE)
        )

        return if (category == VideoCategory.ALL) allChannels else allChannels.filter { it.category == category }
    }

    override suspend fun getChannel(channelId: String): ChannelInfo? {
        return getChannels().find { it.id == channelId }
    }

    override suspend fun getVideo(videoId: String): VideoContent? {
        val allVideos = getShortVideos() + getLongVideos()
        return allVideos.find { it.id == videoId }
    }

    override suspend fun likeVideo(videoId: String): Boolean {
        return if (likedVideos.contains(videoId)) {
            likedVideos.remove(videoId)
            false
        } else {
            likedVideos.add(videoId)
            true
        }
    }

    override suspend fun unlikeVideo(videoId: String): Boolean {
        likedVideos.remove(videoId)
        return false
    }

    override suspend fun subscribeToChannel(channelId: String): Boolean {
        return if (subscribedChannels.contains(channelId)) {
            subscribedChannels.remove(channelId)
            false
        } else {
            subscribedChannels.add(channelId)
            true
        }
    }

    override suspend fun unsubscribeFromChannel(channelId: String): Boolean {
        subscribedChannels.remove(channelId)
        return false
    }

    override fun observeVideos(): Flow<List<VideoContent>> = flow {
        while (true) {
            emit(getShortVideos() + getLongVideos())
            delay(30000) // Emit every 30 seconds
        }
    }

    override fun observeSubscriptions(): Flow<Set<String>> = flow {
        while (true) {
            emit(subscribedChannels.toSet())
            delay(5000) // Emit every 5 seconds
        }
    }

    // Helper functions
    private fun createVideoContent(
        id: String,
        title: String,
        description: String,
        thumbnail: String,
        videoUrl: String,
        duration: String,
        durationSeconds: Int,
        views: Long,
        likes: Long,
        comments: Long,
        category: VideoCategory,
        tags: List<String>,
        channelId: String,
        type: VideoType
    ): VideoContent {
        return VideoContent(
            id = id,
            title = title,
            description = description,
            thumbnail = thumbnail,
            videoUrl = videoUrl,
            duration = duration,
            durationSeconds = durationSeconds,
            views = views,
            likes = likes,
            dislikes = likes / 20, // Mock dislikes as 5% of likes
            comments = comments,
            shares = likes / 10, // Mock shares as 10% of likes
            bookmarks = likes / 15, // Mock bookmarks as ~6% of likes
            channel = createChannelInfo(channelId, "", category),
            uploadTime = Clock.System.now(),
            uploadTimeFormatted = when ((1..10).random()) {
                1 -> "1 hour ago"
                2 -> "3 hours ago"
                3 -> "6 hours ago"
                4 -> "12 hours ago"
                5 -> "1 day ago"
                6 -> "2 days ago"
                7 -> "3 days ago"
                8 -> "1 week ago"
                9 -> "2 weeks ago"
                else -> "1 month ago"
            },
            category = category,
            tags = tags,
            type = type
        )
    }

    private fun createChannelInfo(id: String, name: String, category: VideoCategory): ChannelInfo {
        return when (id) {
            "thai-chef" -> ChannelInfo(
                id = id,
                name = "Bangkok Street Chef",
                avatar = "https://images.unsplash.com/photo-1595273670150-bd0c3c392e46?w=150&h=150&fit=crop",
                banner = "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800&h=200&fit=crop",
                description = "Authentic Thai street food recipes and cooking techniques from Bangkok's best chefs.",
                subscribers = 45000,
                totalVideos = 234,
                verified = true,
                category = VideoCategory.FOOD,
                joinedDate = "Jan 2020",
                location = "Bangkok, Thailand"
            )
            "thai-culture" -> ChannelInfo(
                id = id,
                name = "Thai Culture Hub",
                avatar = "https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop",
                banner = "https://images.unsplash.com/photo-1528181304800-259b08848526?w=800&h=200&fit=crop",
                description = "Explore the rich culture, traditions, and modern life of Thailand through our videos.",
                subscribers = 89000,
                totalVideos = 156,
                verified = true,
                category = VideoCategory.ENTERTAINMENT,
                joinedDate = "Mar 2019",
                location = "Chiang Mai, Thailand"
            )
            "bangkok-explorer" -> ChannelInfo(
                id = id,
                name = "Bangkok Explorer",
                avatar = "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop",
                banner = "https://images.unsplash.com/photo-1508009603885-50cf7c579365?w=800&h=200&fit=crop",
                description = "Your guide to the best spots in Bangkok - from hidden gems to popular attractions.",
                subscribers = 32000,
                totalVideos = 89,
                verified = false,
                category = VideoCategory.TRAVEL,
                joinedDate = "Aug 2021",
                location = "Bangkok, Thailand"
            )
            "foodie-adventures" -> ChannelInfo(
                id = id,
                name = "Southeast Asian Foodie",
                avatar = "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop",
                banner = "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800&h=200&fit=crop",
                description = "Food adventures across Southeast Asia. Discover the best local cuisines and hidden food gems.",
                subscribers = 235000,
                totalVideos = 412,
                verified = true,
                category = VideoCategory.FOOD,
                joinedDate = "Jun 2018",
                location = "Singapore"
            )
            "learn-thai" -> ChannelInfo(
                id = id,
                name = "Thai Language Academy",
                avatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop",
                banner = "https://images.unsplash.com/photo-1544717297-fa95b6ee9643?w=800&h=200&fit=crop",
                description = "Learn Thai language easily with our structured lessons and practical conversations.",
                subscribers = 87000,
                totalVideos = 178,
                verified = true,
                category = VideoCategory.EDUCATION,
                joinedDate = "Feb 2020",
                location = "Bangkok, Thailand"
            )
            "asian-cuisine" -> ChannelInfo(
                id = id,
                name = "Asian Cuisine Master",
                avatar = "https://images.unsplash.com/photo-1566492031773-4f4e44671d66?w=150&h=150&fit=crop",
                description = "Master chef teaching authentic Asian cooking techniques and recipes.",
                subscribers = 156000,
                totalVideos = 289,
                verified = true,
                category = VideoCategory.FOOD,
                joinedDate = "Sep 2019",
                location = "Hong Kong"
            )
            "tech-reviewer" -> ChannelInfo(
                id = id,
                name = "Tech Review Central",
                avatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop",
                description = "Honest tech reviews and unboxings. Helping you make informed decisions.",
                subscribers = 423000,
                totalVideos = 567,
                verified = true,
                category = VideoCategory.TECH,
                joinedDate = "Mar 2018",
                location = "California, USA"
            )
            "pro-gamer" -> ChannelInfo(
                id = id,
                name = "Pro Gaming Academy",
                avatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop",
                description = "Professional gaming tips, tutorials, and reviews. Level up your gaming skills.",
                subscribers = 789000,
                totalVideos = 1234,
                verified = true,
                category = VideoCategory.GAMING,
                joinedDate = "Jan 2017",
                location = "Seoul, South Korea"
            )
            "fitness-coach" -> ChannelInfo(
                id = id,
                name = "Fitness Revolution",
                avatar = "https://images.unsplash.com/photo-1594736797933-d0401ba9d4c4?w=150&h=150&fit=crop",
                description = "Transform your body and mind with effective workouts and nutrition guidance.",
                subscribers = 234000,
                totalVideos = 345,
                verified = true,
                category = VideoCategory.LIFESTYLE,
                joinedDate = "May 2019",
                location = "Los Angeles, USA"
            )
            "business-mentor" -> ChannelInfo(
                id = id,
                name = "Business Growth Hub",
                avatar = "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop",
                description = "Business strategies, entrepreneurship tips, and success mindset coaching.",
                subscribers = 134000,
                totalVideos = 234,
                verified = true,
                category = VideoCategory.BUSINESS,
                joinedDate = "Aug 2020",
                location = "New York, USA"
            )
            "tech-futurist" -> ChannelInfo(
                id = id,
                name = "Tech Futurist",
                avatar = "https://images.unsplash.com/photo-1560250097-0b93528c311a?w=150&h=150&fit=crop",
                description = "Exploring emerging technologies and their impact on society and business.",
                subscribers = 345000,
                totalVideos = 189,
                verified = true,
                category = VideoCategory.TECH,
                joinedDate = "Nov 2018",
                location = "Silicon Valley, USA"
            )
            "startup-stories" -> ChannelInfo(
                id = id,
                name = "Startup Success Stories",
                avatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop",
                description = "Real startup journeys, failures, and successes. Learn from other entrepreneurs.",
                subscribers = 198000,
                totalVideos = 156,
                verified = true,
                category = VideoCategory.BUSINESS,
                joinedDate = "Apr 2019",
                location = "Austin, USA"
            )
            "transformation-coach" -> ChannelInfo(
                id = id,
                name = "Total Life Transformation",
                avatar = "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop",
                description = "Complete life transformation through fitness, mindset, and lifestyle changes.",
                subscribers = 567000,
                totalVideos = 445,
                verified = true,
                category = VideoCategory.LIFESTYLE,
                joinedDate = "Feb 2018",
                location = "Miami, USA"
            )
            else -> ChannelInfo(
                id = id,
                name = name.ifEmpty { "Unknown Channel" },
                avatar = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop",
                description = "Content creator",
                subscribers = (10000..100000).random().toLong(),
                totalVideos = (50..200).random(),
                verified = (1..10).random() > 7, // 30% chance to be verified
                category = category,
                joinedDate = "2020"
            )
        }
    }
}