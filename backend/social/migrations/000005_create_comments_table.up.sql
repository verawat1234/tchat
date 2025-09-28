-- Migration: create_comments_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates comments table for post interactions

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL,
    author_id UUID NOT NULL,
    parent_id UUID, -- For threaded replies

    -- Content
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',

    -- Interaction counts
    likes_count INTEGER DEFAULT 0 CHECK (likes_count >= 0),
    replies_count INTEGER DEFAULT 0 CHECK (replies_count >= 0),
    reactions_count INTEGER DEFAULT 0 CHECK (reactions_count >= 0),

    -- Status flags
    is_edited BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_author FOREIGN KEY (author_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

-- Create indexes for efficient comment queries
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);
CREATE INDEX idx_comments_likes_count ON comments(likes_count DESC);
CREATE INDEX idx_comments_deleted ON comments(is_deleted, deleted_at) WHERE is_deleted = TRUE;
CREATE INDEX idx_comments_metadata ON comments USING GIN(metadata);

-- Composite indexes for mobile comment threading
CREATE INDEX idx_comments_post_created ON comments(post_id, created_at ASC);
CREATE INDEX idx_comments_parent_created ON comments(parent_id, created_at ASC) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_comments_author_created ON comments(author_id, created_at DESC);

-- Create updated_at trigger
CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update comment counts
CREATE OR REPLACE FUNCTION update_comment_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Update post comments count
        UPDATE posts
        SET comments_count = comments_count + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.post_id;

        -- Update parent comment replies count if this is a reply
        IF NEW.parent_id IS NOT NULL THEN
            UPDATE comments
            SET replies_count = replies_count + 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        END IF;

        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Update post comments count
        UPDATE posts
        SET comments_count = GREATEST(comments_count - 1, 0),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.post_id;

        -- Update parent comment replies count if this was a reply
        IF OLD.parent_id IS NOT NULL THEN
            UPDATE comments
            SET replies_count = GREATEST(replies_count - 1, 0),
                updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.parent_id;
        END IF;

        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update comment counts
CREATE TRIGGER update_comment_counts_insert
    AFTER INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_comment_counts();

CREATE TRIGGER update_comment_counts_delete
    AFTER DELETE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_comment_counts();