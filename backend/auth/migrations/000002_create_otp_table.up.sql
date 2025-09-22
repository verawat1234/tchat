-- Migration: create_otp_table
-- Created at: 2024-12-20T10:01:00Z

-- Create OTP table for verification codes
CREATE TABLE IF NOT EXISTS otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    country_code VARCHAR(5) NOT NULL DEFAULT '+66',
    code VARCHAR(10) NOT NULL,
    code_hash VARCHAR(255) NOT NULL,
    purpose VARCHAR(20) NOT NULL CHECK (purpose IN ('registration', 'login', 'phone_change', 'password_reset')),
    attempts INT DEFAULT 0,
    max_attempts INT DEFAULT 3,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP WITH TIME ZONE,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_otp_codes_phone_number ON otp_codes(phone_number);
CREATE INDEX idx_otp_codes_code_hash ON otp_codes(code_hash);
CREATE INDEX idx_otp_codes_purpose ON otp_codes(purpose);
CREATE INDEX idx_otp_codes_expires_at ON otp_codes(expires_at);
CREATE INDEX idx_otp_codes_verified ON otp_codes(verified);
CREATE INDEX idx_otp_codes_created_at ON otp_codes(created_at);

-- Create composite index for active OTP lookup
CREATE INDEX idx_otp_codes_active ON otp_codes(phone_number, purpose, verified, expires_at)
    WHERE verified = FALSE AND expires_at > CURRENT_TIMESTAMP;