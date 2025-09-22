-- Rollback migration: create_dialogs_table
-- Created at: 2024-12-20T10:03:00Z

-- Drop trigger
DROP TRIGGER IF EXISTS update_dialogs_updated_at ON dialogs;

-- Drop indexes
DROP INDEX IF EXISTS idx_dialogs_created_at;
DROP INDEX IF EXISTS idx_dialogs_last_message_at;
DROP INDEX IF EXISTS idx_dialogs_is_archived;
DROP INDEX IF EXISTS idx_dialogs_is_active;
DROP INDEX IF EXISTS idx_dialogs_creator_id;
DROP INDEX IF EXISTS idx_dialogs_type;

-- Drop table
DROP TABLE IF EXISTS dialogs;