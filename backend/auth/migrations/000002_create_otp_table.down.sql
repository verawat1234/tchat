-- Rollback migration: create_otp_table
-- Created at: 2024-12-20T10:01:00Z

-- Drop indexes
DROP INDEX IF EXISTS idx_otp_codes_active;
DROP INDEX IF EXISTS idx_otp_codes_created_at;
DROP INDEX IF EXISTS idx_otp_codes_verified;
DROP INDEX IF EXISTS idx_otp_codes_expires_at;
DROP INDEX IF EXISTS idx_otp_codes_purpose;
DROP INDEX IF EXISTS idx_otp_codes_code_hash;
DROP INDEX IF EXISTS idx_otp_codes_phone_number;

-- Drop table
DROP TABLE IF EXISTS otp_codes;