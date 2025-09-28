-- Migration: rollback create_reactions_table
-- Purpose: Drops reactions table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_reaction_counts_delete ON reactions;
DROP TRIGGER IF EXISTS update_reaction_counts_insert ON reactions;
DROP TRIGGER IF EXISTS update_reactions_updated_at ON reactions;

-- Drop functions
DROP FUNCTION IF EXISTS update_reaction_counts();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_reactions_user_id;
DROP INDEX IF EXISTS idx_reactions_target;
DROP INDEX IF EXISTS idx_reactions_type;
DROP INDEX IF EXISTS idx_reactions_created_at;
DROP INDEX IF EXISTS idx_reactions_target_type;
DROP INDEX IF EXISTS idx_reactions_user_target;

-- Drop table
DROP TABLE IF EXISTS reactions;