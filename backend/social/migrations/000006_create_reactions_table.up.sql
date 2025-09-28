-- Migration: create_reactions_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates reactions table for post and comment interactions

-- Create reactions table
CREATE TABLE IF NOT EXISTS reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    target_id UUID NOT NULL, -- post or comment ID
    target_type VARCHAR(20) NOT NULL CHECK (target_type IN ('post', 'comment')),
    type VARCHAR(20) NOT NULL CHECK (type IN ('like', 'love', 'laugh', 'angry', 'sad', 'wow')),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_reactions_user FOREIGN KEY (user_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_target_reaction UNIQUE (user_id, target_id, target_type)
);

-- Create indexes for efficient reaction queries
CREATE INDEX idx_reactions_user_id ON reactions(user_id);
CREATE INDEX idx_reactions_target ON reactions(target_id, target_type);
CREATE INDEX idx_reactions_type ON reactions(type);
CREATE INDEX idx_reactions_created_at ON reactions(created_at DESC);

-- Composite indexes for mobile reaction queries
CREATE INDEX idx_reactions_target_type ON reactions(target_id, target_type, type);
CREATE INDEX idx_reactions_user_target ON reactions(user_id, target_id, target_type);

-- Create updated_at trigger
CREATE TRIGGER update_reactions_updated_at
    BEFORE UPDATE ON reactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update reaction counts
CREATE OR REPLACE FUNCTION update_reaction_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.target_type = 'post' THEN
            UPDATE posts
            SET reactions_count = reactions_count + 1,
                likes_count = CASE WHEN NEW.type = 'like' THEN likes_count + 1 ELSE likes_count END,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.target_id;
        ELSIF NEW.target_type = 'comment' THEN
            UPDATE comments
            SET reactions_count = reactions_count + 1,
                likes_count = CASE WHEN NEW.type = 'like' THEN likes_count + 1 ELSE likes_count END,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.target_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.target_type = 'post' THEN
            UPDATE posts
            SET reactions_count = GREATEST(reactions_count - 1, 0),
                likes_count = CASE WHEN OLD.type = 'like' THEN GREATEST(likes_count - 1, 0) ELSE likes_count END,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.target_id;
        ELSIF OLD.target_type = 'comment' THEN
            UPDATE comments
            SET reactions_count = GREATEST(reactions_count - 1, 0),
                likes_count = CASE WHEN OLD.type = 'like' THEN GREATEST(likes_count - 1, 0) ELSE likes_count END,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.target_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update reaction counts
CREATE TRIGGER update_reaction_counts_insert
    AFTER INSERT ON reactions
    FOR EACH ROW
    EXECUTE FUNCTION update_reaction_counts();

CREATE TRIGGER update_reaction_counts_delete
    AFTER DELETE ON reactions
    FOR EACH ROW
    EXECUTE FUNCTION update_reaction_counts();