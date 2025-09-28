-- Migration: rollback create_communities_table
-- Purpose: Drops communities table and related indexes

-- Drop triggers
DROP TRIGGER IF EXISTS update_communities_updated_at ON communities;

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_communities_name;
DROP INDEX IF EXISTS idx_communities_type;
DROP INDEX IF EXISTS idx_communities_category;
DROP INDEX IF EXISTS idx_communities_region;
DROP INDEX IF EXISTS idx_communities_creator_id;
DROP INDEX IF EXISTS idx_communities_verified;
DROP INDEX IF EXISTS idx_communities_members_count;
DROP INDEX IF EXISTS idx_communities_posts_count;
DROP INDEX IF EXISTS idx_communities_created_at;
DROP INDEX IF EXISTS idx_communities_tags;
DROP INDEX IF EXISTS idx_communities_settings;
DROP INDEX IF EXISTS idx_communities_type_region;
DROP INDEX IF EXISTS idx_communities_category_created;
DROP INDEX IF EXISTS idx_communities_region_created;

-- Drop table
DROP TABLE IF EXISTS communities;