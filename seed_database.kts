#!/usr/bin/env kotlin

@file:Repository("https://repo.maven.apache.org/maven2/")
@file:DependsOn("app.cash.sqldelight:sqlite-driver:2.0.1")

import app.cash.sqldelight.driver.jdbc.sqlite.JdbcSqliteDriver
import java.io.File

// Simple database seeding script
println("🌱 Tchat Database Seeder")
println("=" * 50)

// Database file path (Android emulator or device path)
val dbPath = "./tchat_seed.db"

try {
    // Create SQLite driver
    val driver = JdbcSqliteDriver("jdbc:sqlite:$dbPath")

    // Create tables manually (simplified schema)
    driver.execute(null, """
        CREATE TABLE IF NOT EXISTS chatSession (
            id TEXT NOT NULL PRIMARY KEY,
            name TEXT,
            type TEXT NOT NULL,
            unreadCount INTEGER NOT NULL DEFAULT 0,
            isPinned INTEGER NOT NULL DEFAULT 0,
            isMuted INTEGER NOT NULL DEFAULT 0,
            isArchived INTEGER NOT NULL DEFAULT 0,
            isBlocked INTEGER NOT NULL DEFAULT 0,
            createdAt TEXT NOT NULL,
            updatedAt TEXT NOT NULL,
            lastActivityAt TEXT
        );
    """.trimIndent(), 0)

    driver.execute(null, """
        CREATE TABLE IF NOT EXISTS message (
            id TEXT NOT NULL PRIMARY KEY,
            chatId TEXT NOT NULL,
            senderId TEXT NOT NULL,
            senderName TEXT NOT NULL,
            type TEXT NOT NULL,
            content TEXT NOT NULL,
            isEdited INTEGER NOT NULL DEFAULT 0,
            isPinned INTEGER NOT NULL DEFAULT 0,
            isDeleted INTEGER NOT NULL DEFAULT 0,
            replyToId TEXT,
            reactions TEXT NOT NULL DEFAULT '[]',
            attachmentCount INTEGER NOT NULL DEFAULT 0,
            createdAt TEXT NOT NULL,
            editedAt TEXT,
            deletedAt TEXT,
            FOREIGN KEY (chatId) REFERENCES chatSession(id) ON DELETE CASCADE
        );
    """.trimIndent(), 0)

    driver.execute(null, """
        CREATE TABLE IF NOT EXISTS user_profiles (
            user_id TEXT NOT NULL PRIMARY KEY,
            display_name TEXT NOT NULL,
            username TEXT NOT NULL UNIQUE,
            avatar_url TEXT,
            bio TEXT,
            is_verified INTEGER NOT NULL DEFAULT 0,
            is_online INTEGER NOT NULL DEFAULT 0,
            last_seen INTEGER,
            status_message TEXT,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL
        );
    """.trimIndent(), 0)

    println("✅ Database tables created")

    // Seed data
    val currentTime = System.currentTimeMillis()

    // Insert users
    println("👥 Seeding users...")
    val users = listOf(
        listOf("current_user", "You", "current_user", "Test user"),
        listOf("user_1", "Alice Johnson", "alice_j", "UI/UX Designer"),
        listOf("user_2", "Bob Smith", "bob_dev", "Software Developer"),
        listOf("user_3", "Carol Zhang", "carol_create", "Content Creator")
    )

    users.forEach { (userId, displayName, username, bio) ->
        driver.execute(null, """
            INSERT OR REPLACE INTO user_profiles
            (user_id, display_name, username, avatar_url, bio, is_verified, is_online, last_seen, status_message, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """.trimIndent(), 11) {
            bindString(0, userId)
            bindString(1, displayName)
            bindString(2, username)
            bindString(3, "https://api.dicebear.com/7.x/avataaars/svg?seed=$username")
            bindString(4, bio)
            bindLong(5, if (userId == "user_1" || userId == "user_3") 1L else 0L)
            bindLong(6, 1L)
            bindLong(7, currentTime)
            bindString(8, "Online")
            bindLong(9, currentTime - 86400000)
            bindLong(10, currentTime)
        }
    }

    // Insert chats
    println("💬 Seeding chats...")
    val chats = listOf(
        listOf("chat_1", "Alice Johnson", "direct", 2L, 1L),
        listOf("chat_2", "Development Team", "group", 5L, 0L),
        listOf("chat_3", "Carol Zhang", "direct", 0L, 0L)
    )

    chats.forEach { (chatId, name, type, unreadCount, isPinned) ->
        driver.execute(null, """
            INSERT OR REPLACE INTO chatSession
            (id, name, type, unreadCount, isPinned, isMuted, isArchived, isBlocked, createdAt, updatedAt, lastActivityAt)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """.trimIndent(), 11) {
            bindString(0, chatId)
            bindString(1, name)
            bindString(2, type)
            bindLong(3, unreadCount)
            bindLong(4, isPinned)
            bindLong(5, 0L)
            bindLong(6, 0L)
            bindLong(7, 0L)
            bindString(8, (currentTime - 86400000).toString())
            bindString(9, currentTime.toString())
            bindString(10, (currentTime - 300000).toString())
        }
    }

    // Insert messages
    println("📝 Seeding messages...")
    val messages = listOf(
        listOf("msg_1", "chat_1", "user_1", "Alice Johnson", "Hey! How's the new app coming along? 😊", (currentTime - 300000).toString()),
        listOf("msg_2", "chat_1", "user_1", "Alice Johnson", "The UI looks amazing! 🎨", (currentTime - 240000).toString()),
        listOf("msg_3", "chat_2", "user_2", "Bob Smith", "Good morning team! Daily standup in 10 minutes 👋", (currentTime - 1800000).toString()),
        listOf("msg_4", "chat_2", "current_user", "You", "Sounds good! Working on the social features now 🚀", (currentTime - 1700000).toString()),
        listOf("msg_5", "chat_3", "user_3", "Carol Zhang", "Just posted some amazing photos from my trip to Chiang Mai! 📸", (currentTime - 7200000).toString()),
        listOf("msg_6", "chat_3", "current_user", "You", "Wow, those look incredible! Can't wait to see them 🤩", (currentTime - 7100000).toString())
    )

    messages.forEach { (msgId, chatId, senderId, senderName, content, createdAt) ->
        driver.execute(null, """
            INSERT OR REPLACE INTO message
            (id, chatId, senderId, senderName, type, content, isEdited, isPinned, isDeleted, replyToId, reactions, attachmentCount, createdAt, editedAt, deletedAt)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """.trimIndent(), 15) {
            bindString(0, msgId)
            bindString(1, chatId)
            bindString(2, senderId)
            bindString(3, senderName)
            bindString(4, "text")
            bindString(5, content)
            bindLong(6, 0L)
            bindLong(7, 0L)
            bindLong(8, 0L)
            bindString(9, null)
            bindString(10, "[]")
            bindLong(11, 0L)
            bindString(12, createdAt)
            bindString(13, null)
            bindString(14, null)
        }
    }

    println("✅ Database seeded successfully!")
    println("📍 Database file: $dbPath")
    println("🎯 You can now copy this database to your Android app or use it for testing")

    // Verify data
    println("\n🔍 Verification:")
    val userResult = driver.executeQuery(null, "SELECT COUNT(*) FROM user_profiles", { cursor ->
        cursor.next()
        cursor.getLong(0)
    }, 0)
    println("  Users: $userResult")

    val chatResult = driver.executeQuery(null, "SELECT COUNT(*) FROM chatSession", { cursor ->
        cursor.next()
        cursor.getLong(0)
    }, 0)
    println("  Chats: $chatResult")

    val messageResult = driver.executeQuery(null, "SELECT COUNT(*) FROM message", { cursor ->
        cursor.next()
        cursor.getLong(0)
    }, 0)
    println("  Messages: $messageResult")

    driver.close()

} catch (e: Exception) {
    println("❌ Error: ${e.message}")
    e.printStackTrace()
}