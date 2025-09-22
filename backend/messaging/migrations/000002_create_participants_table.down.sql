-- Rollback migration: create_participants_table
-- Created at: 2024-12-20T10:04:00Z

-- Drop trigger
DROP TRIGGER IF EXISTS update_dialog_participants_updated_at ON dialog_participants;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_dialogs;
DROP INDEX IF EXISTS idx_dialog_participants_active;
DROP INDEX IF EXISTS idx_dialog_participants_is_archived;
DROP INDEX IF EXISTS idx_dialog_participants_is_pinned;
DROP INDEX IF EXISTS idx_dialog_participants_is_muted;
DROP INDEX IF EXISTS idx_dialog_participants_unread_count;
DROP INDEX IF EXISTS idx_dialog_participants_last_read_at;
DROP INDEX IF EXISTS idx_dialog_participants_left_at;
DROP INDEX IF EXISTS idx_dialog_participants_joined_at;
DROP INDEX IF EXISTS idx_dialog_participants_role;
DROP INDEX IF EXISTS idx_dialog_participants_user_id;
DROP INDEX IF EXISTS idx_dialog_participants_dialog_id;

-- Drop table
DROP TABLE IF EXISTS dialog_participants;