-- Migration: create_follows_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates follow relationships table for social networking

-- Create follows table
CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,

    -- Metadata for analytics and mobile features
    source VARCHAR(30) NOT NULL DEFAULT 'manual' CHECK (source IN ('discovery', 'suggestion', 'search', 'follow_back', 'manual')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('pending', 'active', 'blocked')),
    is_mutual BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_follows_follower FOREIGN KEY (follower_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_follows_following FOREIGN KEY (following_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT chk_no_self_follow CHECK (follower_id != following_id),
    CONSTRAINT unique_follow_relationship UNIQUE (follower_id, following_id)
);

-- Create indexes for efficient social graph queries
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
CREATE INDEX idx_follows_status ON follows(status);
CREATE INDEX idx_follows_source ON follows(source);
CREATE INDEX idx_follows_mutual ON follows(is_mutual) WHERE is_mutual = TRUE;
CREATE INDEX idx_follows_created_at ON follows(created_at);

-- Composite indexes for common mobile queries
CREATE INDEX idx_follows_follower_status ON follows(follower_id, status);
CREATE INDEX idx_follows_following_status ON follows(following_id, status);
CREATE INDEX idx_follows_mutual_pairs ON follows(follower_id, following_id) WHERE is_mutual = TRUE;

-- Create updated_at trigger
CREATE TRIGGER update_follows_updated_at
    BEFORE UPDATE ON follows
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update follow counts
CREATE OR REPLACE FUNCTION update_follow_counts()
RETURNS TRIGGER AS $$
BEGIN
    -- Update follower count for the user being followed
    IF TG_OP = 'INSERT' THEN
        UPDATE social_profiles
        SET followers_count = followers_count + 1,
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.following_id;

        UPDATE social_profiles
        SET following_count = following_count + 1,
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.follower_id;

        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE social_profiles
        SET followers_count = GREATEST(followers_count - 1, 0),
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.following_id;

        UPDATE social_profiles
        SET following_count = GREATEST(following_count - 1, 0),
            social_updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.follower_id;

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update follow counts
CREATE TRIGGER update_follow_counts_insert
    AFTER INSERT ON follows
    FOR EACH ROW
    EXECUTE FUNCTION update_follow_counts();

CREATE TRIGGER update_follow_counts_delete
    AFTER DELETE ON follows
    FOR EACH ROW
    EXECUTE FUNCTION update_follow_counts();