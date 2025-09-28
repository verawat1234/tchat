-- Migration: create_communities_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates communities table for social groups and forums

-- Create communities table
CREATE TABLE IF NOT EXISTS communities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    type VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (type IN ('public', 'private', 'restricted')),
    category VARCHAR(50),
    region VARCHAR(10) CHECK (region IN ('TH', 'SG', 'ID', 'MY', 'PH', 'VN')),

    -- Media and branding
    avatar VARCHAR(255),
    banner VARCHAR(255),

    -- Community configuration
    tags TEXT[] DEFAULT '{}',
    rules TEXT[] DEFAULT '{}',
    creator_id UUID NOT NULL,
    settings JSONB DEFAULT '{}',

    -- Community metrics
    members_count INTEGER DEFAULT 0 CHECK (members_count >= 0),
    posts_count INTEGER DEFAULT 0 CHECK (posts_count >= 0),

    -- Status flags
    is_verified BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_communities_creator FOREIGN KEY (creator_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT unique_community_name UNIQUE (name)
);

-- Create indexes for community discovery and management
CREATE INDEX idx_communities_name ON communities(name);
CREATE INDEX idx_communities_type ON communities(type);
CREATE INDEX idx_communities_category ON communities(category) WHERE category IS NOT NULL;
CREATE INDEX idx_communities_region ON communities(region) WHERE region IS NOT NULL;
CREATE INDEX idx_communities_creator_id ON communities(creator_id);
CREATE INDEX idx_communities_verified ON communities(is_verified) WHERE is_verified = TRUE;
CREATE INDEX idx_communities_members_count ON communities(members_count DESC);
CREATE INDEX idx_communities_posts_count ON communities(posts_count DESC);
CREATE INDEX idx_communities_created_at ON communities(created_at DESC);
CREATE INDEX idx_communities_tags ON communities USING GIN(tags);
CREATE INDEX idx_communities_settings ON communities USING GIN(settings);

-- Composite indexes for mobile community discovery
CREATE INDEX idx_communities_type_region ON communities(type, region);
CREATE INDEX idx_communities_category_created ON communities(category, created_at DESC);
CREATE INDEX idx_communities_region_created ON communities(region, created_at DESC);

-- Create updated_at trigger
CREATE TRIGGER update_communities_updated_at
    BEFORE UPDATE ON communities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();