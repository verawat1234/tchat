-- Migration: create_posts_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates posts table for social content sharing

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL,
    community_id UUID, -- Optional: posts can belong to communities

    -- Content fields
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'text' CHECK (type IN ('text', 'image', 'video', 'link', 'poll')),
    metadata JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    visibility VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'members', 'private', 'followers')),
    media_urls TEXT[] DEFAULT '{}',
    link_preview JSONB DEFAULT '{}',

    -- Interaction counts (denormalized for mobile performance)
    likes_count INTEGER DEFAULT 0 CHECK (likes_count >= 0),
    comments_count INTEGER DEFAULT 0 CHECK (comments_count >= 0),
    shares_count INTEGER DEFAULT 0 CHECK (shares_count >= 0),
    reactions_count INTEGER DEFAULT 0 CHECK (reactions_count >= 0),
    views_count INTEGER DEFAULT 0 CHECK (views_count >= 0),

    -- Status flags
    is_edited BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    is_trending BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT fk_posts_author FOREIGN KEY (author_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_posts_community FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE SET NULL
);

-- Create indexes for efficient content queries
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_community_id ON posts(community_id) WHERE community_id IS NOT NULL;
CREATE INDEX idx_posts_type ON posts(type);
CREATE INDEX idx_posts_visibility ON posts(visibility);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_updated_at ON posts(updated_at);
CREATE INDEX idx_posts_likes_count ON posts(likes_count DESC);
CREATE INDEX idx_posts_comments_count ON posts(comments_count DESC);
CREATE INDEX idx_posts_trending ON posts(is_trending) WHERE is_trending = TRUE;
CREATE INDEX idx_posts_pinned ON posts(is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX idx_posts_deleted ON posts(is_deleted, deleted_at) WHERE is_deleted = TRUE;
CREATE INDEX idx_posts_tags ON posts USING GIN(tags);
CREATE INDEX idx_posts_metadata ON posts USING GIN(metadata);

-- Composite indexes for mobile feed queries
CREATE INDEX idx_posts_author_created ON posts(author_id, created_at DESC);
CREATE INDEX idx_posts_community_created ON posts(community_id, created_at DESC) WHERE community_id IS NOT NULL;
CREATE INDEX idx_posts_visibility_created ON posts(visibility, created_at DESC);
CREATE INDEX idx_posts_trending_created ON posts(is_trending, created_at DESC) WHERE is_trending = TRUE;

-- Create updated_at trigger
CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update user posts count
CREATE OR REPLACE FUNCTION update_posts_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE social_profiles
        SET posts_count = posts_count + 1,
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.author_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE social_profiles
        SET posts_count = GREATEST(posts_count - 1, 0),
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.author_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update posts count
CREATE TRIGGER update_posts_count_insert
    AFTER INSERT ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_posts_count();

CREATE TRIGGER update_posts_count_delete
    AFTER DELETE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_posts_count();