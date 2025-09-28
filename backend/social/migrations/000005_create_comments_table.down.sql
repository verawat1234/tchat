-- Migration: rollback create_comments_table
-- Purpose: Drops comments table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_comment_counts_delete ON comments;
DROP TRIGGER IF EXISTS update_comment_counts_insert ON comments;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;

-- Drop functions
DROP FUNCTION IF EXISTS update_comment_counts();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_comments_post_id;
DROP INDEX IF EXISTS idx_comments_author_id;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_created_at;
DROP INDEX IF EXISTS idx_comments_likes_count;
DROP INDEX IF EXISTS idx_comments_deleted;
DROP INDEX IF EXISTS idx_comments_metadata;
DROP INDEX IF EXISTS idx_comments_post_created;
DROP INDEX IF EXISTS idx_comments_parent_created;
DROP INDEX IF EXISTS idx_comments_author_created;

-- Drop table
DROP TABLE IF EXISTS comments;