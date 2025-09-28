-- Migration: rollback create_posts_table
-- Purpose: Drops posts table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_posts_count_delete ON posts;
DROP TRIGGER IF EXISTS update_posts_count_insert ON posts;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;

-- Drop functions
DROP FUNCTION IF EXISTS update_posts_count();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_posts_author_id;
DROP INDEX IF EXISTS idx_posts_community_id;
DROP INDEX IF EXISTS idx_posts_type;
DROP INDEX IF EXISTS idx_posts_visibility;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_updated_at;
DROP INDEX IF EXISTS idx_posts_likes_count;
DROP INDEX IF EXISTS idx_posts_comments_count;
DROP INDEX IF EXISTS idx_posts_trending;
DROP INDEX IF EXISTS idx_posts_pinned;
DROP INDEX IF EXISTS idx_posts_deleted;
DROP INDEX IF EXISTS idx_posts_tags;
DROP INDEX IF EXISTS idx_posts_metadata;
DROP INDEX IF EXISTS idx_posts_author_created;
DROP INDEX IF EXISTS idx_posts_community_created;
DROP INDEX IF EXISTS idx_posts_visibility_created;
DROP INDEX IF EXISTS idx_posts_trending_created;

-- Drop table
DROP TABLE IF EXISTS posts;