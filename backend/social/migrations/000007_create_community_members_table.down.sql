-- Migration: rollback create_community_members_table
-- Purpose: Drops community_members table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_community_member_counts_delete ON community_members;
DROP TRIGGER IF EXISTS update_community_member_counts_insert ON community_members;
DROP TRIGGER IF EXISTS update_community_members_updated_at ON community_members;

-- Drop functions
DROP FUNCTION IF EXISTS update_community_member_counts();

-- Drop table
DROP TABLE IF EXISTS community_members;