-- Migration: create_social_profiles_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates social-specific user profile extensions for KMP mobile integration

-- Create social_profiles table (extends users table)
CREATE TABLE IF NOT EXISTS social_profiles (
    id UUID PRIMARY KEY,

    -- Core user fields (duplicated for social service isolation)
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    bio TEXT,
    avatar VARCHAR(500),
    country VARCHAR(5) NOT NULL DEFAULT 'TH' CHECK (country IN ('TH', 'SG', 'ID', 'MY', 'PH', 'VN')),
    locale VARCHAR(10) NOT NULL DEFAULT 'th',
    timezone VARCHAR(50) DEFAULT 'Asia/Bangkok',
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended', 'banned')),
    is_active BOOLEAN DEFAULT TRUE,
    kyc_tier INTEGER DEFAULT 0 CHECK (kyc_tier >= 0 AND kyc_tier <= 3),
    is_verified BOOLEAN DEFAULT FALSE,

    -- Social-specific fields for KMP compatibility
    interests TEXT[] DEFAULT '{}',
    social_links JSONB DEFAULT '{}',
    social_preferences JSONB DEFAULT '{}',

    -- Social metrics (always present for mobile UI consistency)
    followers_count INTEGER DEFAULT 0 CHECK (followers_count >= 0),
    following_count INTEGER DEFAULT 0 CHECK (following_count >= 0),
    posts_count INTEGER DEFAULT 0 CHECK (posts_count >= 0),

    -- Social verification (separate from KYC)
    is_social_verified BOOLEAN DEFAULT FALSE,

    -- Timestamps (RFC3339 format for KMP)
    social_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    social_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_social_profiles_users FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for performance and mobile queries
CREATE INDEX idx_social_profiles_username ON social_profiles(username);
CREATE INDEX idx_social_profiles_email ON social_profiles(email) WHERE email IS NOT NULL;
CREATE INDEX idx_social_profiles_country ON social_profiles(country);
CREATE INDEX idx_social_profiles_status ON social_profiles(status);
CREATE INDEX idx_social_profiles_verified ON social_profiles(is_social_verified);
CREATE INDEX idx_social_profiles_followers ON social_profiles(followers_count DESC);
CREATE INDEX idx_social_profiles_created_at ON social_profiles(social_created_at);
CREATE INDEX idx_social_profiles_interests ON social_profiles USING GIN(interests);
CREATE INDEX idx_social_profiles_social_links ON social_profiles USING GIN(social_links);

-- Create updated_at trigger for social_profiles
CREATE TRIGGER update_social_profiles_updated_at
    BEFORE UPDATE ON social_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to automatically update social_updated_at
CREATE OR REPLACE FUNCTION update_social_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.social_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_social_profiles_social_updated_at
    BEFORE UPDATE ON social_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_social_updated_at_column();