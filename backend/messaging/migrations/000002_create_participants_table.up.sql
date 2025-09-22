-- Migration: create_participants_table
-- Created at: 2024-12-20T10:04:00Z

-- Create dialog participants table
CREATE TABLE IF NOT EXISTS dialog_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dialog_id UUID NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(20) DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'moderator', 'member')),
    permissions JSONB DEFAULT '{}',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    last_read_message_id UUID,
    last_read_at TIMESTAMP WITH TIME ZONE,
    unread_count INT DEFAULT 0,
    is_muted BOOLEAN DEFAULT FALSE,
    muted_until TIMESTAMP WITH TIME ZONE,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(dialog_id, user_id)
);

-- Create indexes
CREATE INDEX idx_dialog_participants_dialog_id ON dialog_participants(dialog_id);
CREATE INDEX idx_dialog_participants_user_id ON dialog_participants(user_id);
CREATE INDEX idx_dialog_participants_role ON dialog_participants(role);
CREATE INDEX idx_dialog_participants_joined_at ON dialog_participants(joined_at);
CREATE INDEX idx_dialog_participants_left_at ON dialog_participants(left_at);
CREATE INDEX idx_dialog_participants_last_read_at ON dialog_participants(last_read_at);
CREATE INDEX idx_dialog_participants_unread_count ON dialog_participants(unread_count);
CREATE INDEX idx_dialog_participants_is_muted ON dialog_participants(is_muted);
CREATE INDEX idx_dialog_participants_is_pinned ON dialog_participants(is_pinned);
CREATE INDEX idx_dialog_participants_is_archived ON dialog_participants(is_archived);

-- Create composite index for active participants
CREATE INDEX idx_dialog_participants_active ON dialog_participants(dialog_id, user_id)
    WHERE left_at IS NULL;

-- Create composite index for user's dialogs
CREATE INDEX idx_user_dialogs ON dialog_participants(user_id, is_archived, last_read_at)
    WHERE left_at IS NULL;

-- Create updated_at trigger
CREATE TRIGGER update_dialog_participants_updated_at
    BEFORE UPDATE ON dialog_participants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();