-- ScyllaDB keyspace and tables for Tchat messaging system
-- Optimized for high-throughput message storage and retrieval

-- Create keyspace for messaging
CREATE KEYSPACE IF NOT EXISTS tchat_messages
WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 3
}
AND DURABLE_WRITES = true;

USE tchat_messages;

-- Messages table - primary storage for all messages
-- Partitioned by dialog_id for efficient message retrieval
CREATE TABLE IF NOT EXISTS messages (
    dialog_id UUID,
    message_id UUID,
    sender_id UUID,
    message_type TEXT, -- 'text', 'image', 'file', 'video', 'audio', 'location'
    content TEXT,
    metadata MAP<TEXT, TEXT>, -- Store additional message metadata
    attachments LIST<FROZEN<MAP<TEXT, TEXT>>>, -- File attachments
    reply_to_id UUID,
    edit_history LIST<FROZEN<MAP<TEXT, TEXT>>>, -- Track message edits
    delivery_status MAP<UUID, TEXT>, -- Per-participant delivery status
    read_status MAP<UUID, TIMESTAMP>, -- Per-participant read timestamps
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (dialog_id, created_at, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC, message_id ASC)
AND COMPACTION = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 1
}
AND GC_GRACE_SECONDS = 86400 -- 1 day
AND DEFAULT_TIME_TO_LIVE = 31536000; -- 1 year

-- Messages by user - secondary index for user's message history
CREATE TABLE IF NOT EXISTS messages_by_user (
    user_id UUID,
    dialog_id UUID,
    message_id UUID,
    created_at TIMESTAMP,
    content TEXT,
    message_type TEXT,
    sender_id UUID,
    PRIMARY KEY (user_id, created_at, dialog_id, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC, dialog_id ASC, message_id ASC)
AND COMPACTION = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 7
}
AND GC_GRACE_SECONDS = 86400
AND DEFAULT_TIME_TO_LIVE = 31536000;

-- Dialog participants - manage dialog membership
CREATE TABLE IF NOT EXISTS dialog_participants (
    dialog_id UUID,
    user_id UUID,
    role TEXT, -- 'owner', 'admin', 'member'
    joined_at TIMESTAMP,
    last_read_at TIMESTAMP,
    notification_settings MAP<TEXT, TEXT>,
    is_active BOOLEAN,
    PRIMARY KEY (dialog_id, user_id)
) WITH GC_GRACE_SECONDS = 86400;

-- User dialogs - efficient lookup of user's dialogs
CREATE TABLE IF NOT EXISTS user_dialogs (
    user_id UUID,
    dialog_id UUID,
    dialog_type TEXT,
    last_message_at TIMESTAMP,
    last_message_id UUID,
    unread_count COUNTER,
    is_pinned BOOLEAN,
    is_archived BOOLEAN,
    joined_at TIMESTAMP,
    PRIMARY KEY (user_id, last_message_at, dialog_id)
) WITH CLUSTERING ORDER BY (last_message_at DESC, dialog_id ASC)
AND GC_GRACE_SECONDS = 86400;

-- Message search index - for full-text search capabilities
CREATE TABLE IF NOT EXISTS message_search (
    search_term TEXT,
    dialog_id UUID,
    message_id UUID,
    created_at TIMESTAMP,
    sender_id UUID,
    content_snippet TEXT,
    PRIMARY KEY (search_term, created_at, dialog_id, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC, dialog_id ASC, message_id ASC)
AND COMPACTION = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 7
}
AND GC_GRACE_SECONDS = 86400
AND DEFAULT_TIME_TO_LIVE = 2592000; -- 30 days for search index

-- File attachments metadata
CREATE TABLE IF NOT EXISTS file_attachments (
    file_id UUID,
    message_id UUID,
    dialog_id UUID,
    uploader_id UUID,
    filename TEXT,
    file_size BIGINT,
    mime_type TEXT,
    storage_path TEXT,
    thumbnail_path TEXT,
    encryption_key TEXT, -- For encrypted files
    upload_status TEXT, -- 'uploading', 'completed', 'failed'
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
    PRIMARY KEY (file_id)
) WITH GC_GRACE_SECONDS = 86400
AND DEFAULT_TIME_TO_LIVE = 31536000;

-- Message delivery tracking
CREATE TABLE IF NOT EXISTS message_delivery (
    message_id UUID,
    recipient_id UUID,
    delivery_status TEXT, -- 'sent', 'delivered', 'read', 'failed'
    delivery_timestamp TIMESTAMP,
    failure_reason TEXT,
    retry_count INT,
    PRIMARY KEY (message_id, recipient_id)
) WITH GC_GRACE_SECONDS = 86400;

-- Real-time presence tracking
CREATE TABLE IF NOT EXISTS user_presence (
    user_id UUID,
    status TEXT, -- 'online', 'away', 'busy', 'offline'
    last_seen TIMESTAMP,
    device_type TEXT,
    device_id TEXT,
    PRIMARY KEY (user_id)
) WITH DEFAULT_TIME_TO_LIVE = 300; -- 5 minutes TTL for presence

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS ON messages (sender_id);
CREATE INDEX IF NOT EXISTS ON messages (message_type);
CREATE INDEX IF NOT EXISTS ON messages_by_user (dialog_id);
CREATE INDEX IF NOT EXISTS ON file_attachments (dialog_id);
CREATE INDEX IF NOT EXISTS ON file_attachments (uploader_id);