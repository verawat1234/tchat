-- Migration: create_messages_table
-- Created at: 2024-12-20T10:05:00Z

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dialog_id UUID NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'image', 'video', 'audio', 'file', 'location', 'contact', 'sticker', 'system')),
    content TEXT,
    media_url TEXT,
    media_metadata JSONB,
    reply_to_id UUID REFERENCES messages(id),
    forward_from_id UUID REFERENCES messages(id),
    edit_history JSONB DEFAULT '[]',
    is_edited BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    read_by JSONB DEFAULT '{}',
    reactions JSONB DEFAULT '{}',
    mentions JSONB DEFAULT '[]',
    hashtags JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_messages_dialog_id ON messages(dialog_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_type ON messages(type);
CREATE INDEX idx_messages_reply_to_id ON messages(reply_to_id);
CREATE INDEX idx_messages_forward_from_id ON messages(forward_from_id);
CREATE INDEX idx_messages_is_edited ON messages(is_edited);
CREATE INDEX idx_messages_is_deleted ON messages(is_deleted);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_updated_at ON messages(updated_at);

-- Create composite index for dialog timeline
CREATE INDEX idx_messages_dialog_timeline ON messages(dialog_id, created_at DESC)
    WHERE is_deleted = FALSE;

-- Create composite index for user messages
CREATE INDEX idx_messages_user ON messages(sender_id, created_at DESC)
    WHERE is_deleted = FALSE;

-- Create GIN index for full-text search
CREATE INDEX idx_messages_content_search ON messages USING gin(to_tsvector('english', content))
    WHERE type = 'text' AND is_deleted = FALSE;

-- Create GIN indexes for JSONB fields
CREATE INDEX idx_messages_read_by ON messages USING gin(read_by);
CREATE INDEX idx_messages_reactions ON messages USING gin(reactions);
CREATE INDEX idx_messages_mentions ON messages USING gin(mentions);
CREATE INDEX idx_messages_hashtags ON messages USING gin(hashtags);

-- Create updated_at trigger
CREATE TRIGGER update_messages_updated_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();