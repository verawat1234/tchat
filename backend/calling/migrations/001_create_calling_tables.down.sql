-- Rollback calling service tables
-- Migration: 001_create_calling_tables (DOWN)
-- Description: Drop all calling service tables and related objects

-- Drop triggers
DROP TRIGGER IF EXISTS update_call_sessions_updated_at ON call_sessions;
DROP TRIGGER IF EXISTS update_call_participants_updated_at ON call_participants;
DROP TRIGGER IF EXISTS update_user_presence_updated_at ON user_presence;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (they will be dropped with tables, but explicit for clarity)
DROP INDEX IF EXISTS idx_call_sessions_initiated_by;
DROP INDEX IF EXISTS idx_call_sessions_status;
DROP INDEX IF EXISTS idx_call_sessions_started_at;
DROP INDEX IF EXISTS idx_call_sessions_type;

DROP INDEX IF EXISTS idx_call_participants_call_session_id;
DROP INDEX IF EXISTS idx_call_participants_user_id;
DROP INDEX IF EXISTS idx_call_participants_role;

DROP INDEX IF EXISTS idx_user_presence_user_id;
DROP INDEX IF EXISTS idx_user_presence_status;
DROP INDEX IF EXISTS idx_user_presence_in_call;
DROP INDEX IF EXISTS idx_user_presence_last_seen;

DROP INDEX IF EXISTS idx_call_history_caller_id;
DROP INDEX IF EXISTS idx_call_history_callee_id;
DROP INDEX IF EXISTS idx_call_history_started_at;
DROP INDEX IF EXISTS idx_call_history_call_status;
DROP INDEX IF EXISTS idx_call_history_call_type;

-- Drop tables in correct order (considering foreign key dependencies)
DROP TABLE IF EXISTS call_history;
DROP TABLE IF EXISTS user_presence;
DROP TABLE IF EXISTS call_participants;
DROP TABLE IF EXISTS call_sessions;