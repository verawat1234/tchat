-- Migration: Fix device_id and related field constraints
-- Issue: VARCHAR(255) constraint causing session creation failures
-- Date: 2025-09-27
-- Journey: 6 - Session Schema Fix

BEGIN;

-- Fix session_token constraint (authentication tokens can be very long)
ALTER TABLE user_sessions ALTER COLUMN session_token TYPE TEXT;

-- Fix refresh_token constraint (refresh tokens can be very long)
ALTER TABLE user_sessions ALTER COLUMN refresh_token TYPE TEXT;

-- Fix device_id constraint (main issue)
ALTER TABLE user_sessions ALTER COLUMN device_id TYPE TEXT;

-- Fix user_agent constraint (potential issue)
ALTER TABLE user_sessions ALTER COLUMN user_agent TYPE TEXT;

-- Verify the changes
SELECT column_name, data_type, character_maximum_length
FROM information_schema.columns
WHERE table_name = 'user_sessions'
AND column_name IN ('session_token', 'refresh_token', 'device_id', 'user_agent');

COMMIT;
