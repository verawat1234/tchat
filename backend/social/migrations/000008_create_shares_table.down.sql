-- Migration: rollback create_shares_table
-- Purpose: Drops shares table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_share_counts_delete ON shares;
DROP TRIGGER IF EXISTS update_share_counts_insert ON shares;

-- Drop functions
DROP FUNCTION IF EXISTS update_share_counts();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_shares_user_id;
DROP INDEX IF EXISTS idx_shares_content;
DROP INDEX IF EXISTS idx_shares_platform;
DROP INDEX IF EXISTS idx_shares_privacy;
DROP INDEX IF EXISTS idx_shares_status;
DROP INDEX IF EXISTS idx_shares_created_at;
DROP INDEX IF EXISTS idx_shares_metadata;
DROP INDEX IF EXISTS idx_shares_user_created;
DROP INDEX IF EXISTS idx_shares_content_created;
DROP INDEX IF EXISTS idx_shares_platform_created;

-- Drop table
DROP TABLE IF EXISTS shares;