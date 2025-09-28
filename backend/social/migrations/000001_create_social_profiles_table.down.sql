-- Migration: rollback create_social_profiles_table
-- Purpose: Drops social_profiles table and related functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_social_profiles_social_updated_at ON social_profiles;
DROP TRIGGER IF EXISTS update_social_profiles_updated_at ON social_profiles;

-- Drop functions
DROP FUNCTION IF EXISTS update_social_updated_at_column();

-- Drop indexes (automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_social_profiles_username;
DROP INDEX IF EXISTS idx_social_profiles_email;
DROP INDEX IF EXISTS idx_social_profiles_country;
DROP INDEX IF EXISTS idx_social_profiles_status;
DROP INDEX IF EXISTS idx_social_profiles_verified;
DROP INDEX IF EXISTS idx_social_profiles_followers;
DROP INDEX IF EXISTS idx_social_profiles_created_at;
DROP INDEX IF EXISTS idx_social_profiles_interests;
DROP INDEX IF EXISTS idx_social_profiles_social_links;

-- Drop table
DROP TABLE IF EXISTS social_profiles;