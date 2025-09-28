-- Migration: create_community_members_table
-- Created at: 2024-09-28T16:00:00Z
-- Purpose: Creates community membership table

-- Create community_members table
CREATE TABLE IF NOT EXISTS community_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    community_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'moderator', 'member')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'pending', 'banned')),
    join_reason VARCHAR(255),
    source VARCHAR(20) NOT NULL CHECK (source IN ('discovery', 'invitation', 'search', 'manual')),

    -- Timestamps
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_community_members_community FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE,
    CONSTRAINT fk_community_members_user FOREIGN KEY (user_id) REFERENCES social_profiles(id) ON DELETE CASCADE,
    CONSTRAINT unique_community_membership UNIQUE (community_id, user_id)
);

-- Create indexes for membership queries
CREATE INDEX idx_community_members_community_id ON community_members(community_id);
CREATE INDEX idx_community_members_user_id ON community_members(user_id);
CREATE INDEX idx_community_members_role ON community_members(role);
CREATE INDEX idx_community_members_status ON community_members(status);
CREATE INDEX idx_community_members_source ON community_members(source);
CREATE INDEX idx_community_members_joined_at ON community_members(joined_at);

-- Composite indexes for mobile community management
CREATE INDEX idx_community_members_community_role ON community_members(community_id, role);
CREATE INDEX idx_community_members_community_status ON community_members(community_id, status);
CREATE INDEX idx_community_members_user_joined ON community_members(user_id, joined_at DESC);

-- Create updated_at trigger
CREATE TRIGGER update_community_members_updated_at
    BEFORE UPDATE ON community_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update community member counts
CREATE OR REPLACE FUNCTION update_community_member_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE communities
        SET members_count = members_count + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.community_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE communities
        SET members_count = GREATEST(members_count - 1, 0),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.community_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update member counts
CREATE TRIGGER update_community_member_counts_insert
    AFTER INSERT ON community_members
    FOR EACH ROW
    EXECUTE FUNCTION update_community_member_counts();

CREATE TRIGGER update_community_member_counts_delete
    AFTER DELETE ON community_members
    FOR EACH ROW
    EXECUTE FUNCTION update_community_member_counts();