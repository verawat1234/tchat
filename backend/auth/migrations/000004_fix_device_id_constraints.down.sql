-- Migration rollback: fix_device_id_constraints
-- Created at: 2025-09-27T08:45:00Z
-- Purpose: Rollback device_id constraint fixes (use with caution)

-- WARNING: Rolling back to VARCHAR(255) may cause data truncation
-- Only run this if you're sure no device_id values exceed 255 characters

-- Truncate any device_id values longer than 255 characters (data loss!)
UPDATE user_sessions
SET device_id = LEFT(device_id, 255)
WHERE LENGTH(device_id) > 255;

-- Revert device_id back to VARCHAR(255)
ALTER TABLE user_sessions
ALTER COLUMN device_id TYPE VARCHAR(255);

-- Revert refresh_token_hash back to VARCHAR(255) - WARNING: may cause issues
ALTER TABLE user_sessions
ALTER COLUMN refresh_token_hash TYPE VARCHAR(255);

-- Revert user_agent back to TEXT (keep as TEXT since it was originally TEXT)
-- No change needed for user_agent

-- Remove comments
COMMENT ON COLUMN user_sessions.device_id IS NULL;
COMMENT ON COLUMN user_sessions.refresh_token_hash IS NULL;