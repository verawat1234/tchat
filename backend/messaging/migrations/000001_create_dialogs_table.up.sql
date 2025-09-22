-- Migration: create_dialogs_table
-- Created at: 2024-12-20T10:03:00Z

-- Create dialogs table for conversations
CREATE TABLE IF NOT EXISTS dialogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('direct', 'group', 'channel', 'broadcast')),
    title VARCHAR(255),
    description TEXT,
    avatar_url TEXT,
    creator_id UUID NOT NULL,
    participant_count INT DEFAULT 0,
    message_count INT DEFAULT 0,
    settings JSONB DEFAULT '{}',
    last_message_id UUID,
    last_message_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_dialogs_type ON dialogs(type);
CREATE INDEX idx_dialogs_creator_id ON dialogs(creator_id);
CREATE INDEX idx_dialogs_is_active ON dialogs(is_active);
CREATE INDEX idx_dialogs_is_archived ON dialogs(is_archived);
CREATE INDEX idx_dialogs_last_message_at ON dialogs(last_message_at);
CREATE INDEX idx_dialogs_created_at ON dialogs(created_at);

-- Create updated_at trigger
CREATE TRIGGER update_dialogs_updated_at
    BEFORE UPDATE ON dialogs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();