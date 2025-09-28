-- Migration: create_shares_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates shares table for content sharing tracking

-- Create shares table
CREATE TABLE IF NOT EXISTS shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    content_id UUID NOT NULL,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('post', 'comment', 'community')),
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('internal', 'external')),

    -- Share details
    message VARCHAR(500),
    privacy VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'private', 'followers')),
    metadata JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_shares_user FOREIGN KEY (user_id) REFERENCES social_profiles(id) ON DELETE CASCADE
);

-- Create indexes for share analytics and queries
CREATE INDEX idx_shares_user_id ON shares(user_id);
CREATE INDEX idx_shares_content ON shares(content_id, content_type);
CREATE INDEX idx_shares_platform ON shares(platform);
CREATE INDEX idx_shares_privacy ON shares(privacy);
CREATE INDEX idx_shares_status ON shares(status);
CREATE INDEX idx_shares_created_at ON shares(created_at DESC);
CREATE INDEX idx_shares_metadata ON shares USING GIN(metadata);

-- Composite indexes for mobile share analytics
CREATE INDEX idx_shares_user_created ON shares(user_id, created_at DESC);
CREATE INDEX idx_shares_content_created ON shares(content_id, content_type, created_at DESC);
CREATE INDEX idx_shares_platform_created ON shares(platform, created_at DESC);

-- Function to update share counts
CREATE OR REPLACE FUNCTION update_share_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.content_type = 'post' THEN
            UPDATE posts
            SET shares_count = shares_count + 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.content_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.content_type = 'post' THEN
            UPDATE posts
            SET shares_count = GREATEST(shares_count - 1, 0),
                updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.content_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update share counts
CREATE TRIGGER update_share_counts_insert
    AFTER INSERT ON shares
    FOR EACH ROW
    EXECUTE FUNCTION update_share_counts();

CREATE TRIGGER update_share_counts_delete
    AFTER DELETE ON shares
    FOR EACH ROW
    EXECUTE FUNCTION update_share_counts();