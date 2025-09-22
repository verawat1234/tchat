-- Rollback migration: create_messages_table
-- Created at: 2024-12-20T10:05:00Z

-- Drop trigger
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;

-- Drop indexes
DROP INDEX IF EXISTS idx_messages_hashtags;
DROP INDEX IF EXISTS idx_messages_mentions;
DROP INDEX IF EXISTS idx_messages_reactions;
DROP INDEX IF EXISTS idx_messages_read_by;
DROP INDEX IF EXISTS idx_messages_content_search;
DROP INDEX IF EXISTS idx_messages_user;
DROP INDEX IF EXISTS idx_messages_dialog_timeline;
DROP INDEX IF EXISTS idx_messages_updated_at;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_is_deleted;
DROP INDEX IF EXISTS idx_messages_is_edited;
DROP INDEX IF EXISTS idx_messages_forward_from_id;
DROP INDEX IF EXISTS idx_messages_reply_to_id;
DROP INDEX IF EXISTS idx_messages_type;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_dialog_id;

-- Drop table
DROP TABLE IF EXISTS messages;