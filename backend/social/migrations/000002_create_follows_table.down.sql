-- Migration: rollback create_follows_table
-- Purpose: Drops follows table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_follow_counts_delete ON follows;
DROP TRIGGER IF EXISTS update_follow_counts_insert ON follows;
DROP TRIGGER IF EXISTS update_follows_updated_at ON follows;

-- Drop functions
DROP FUNCTION IF EXISTS update_follow_counts();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_follows_follower_id;
DROP INDEX IF EXISTS idx_follows_following_id;
DROP INDEX IF EXISTS idx_follows_status;
DROP INDEX IF EXISTS idx_follows_source;
DROP INDEX IF EXISTS idx_follows_mutual;
DROP INDEX IF EXISTS idx_follows_created_at;
DROP INDEX IF EXISTS idx_follows_follower_status;
DROP INDEX IF EXISTS idx_follows_following_status;
DROP INDEX IF EXISTS idx_follows_mutual_pairs;

-- Drop table
DROP TABLE IF EXISTS follows;