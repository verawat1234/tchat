-- Create calling service tables
-- Migration: 001_create_calling_tables
-- Description: Initial tables for voice and video calling functionality

-- Create call_sessions table
CREATE TABLE IF NOT EXISTS call_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(10) NOT NULL CHECK (type IN ('voice', 'video')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('connecting', 'active', 'ended', 'failed')),
    initiated_by UUID NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER, -- in seconds
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT call_sessions_started_before_ended CHECK (started_at <= ended_at OR ended_at IS NULL),
    CONSTRAINT call_sessions_duration_positive CHECK (duration IS NULL OR duration >= 0)
);

-- Create call_participants table
CREATE TABLE IF NOT EXISTS call_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    call_session_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('caller', 'callee')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    audio_enabled BOOLEAN DEFAULT true,
    video_enabled BOOLEAN DEFAULT false,
    connection_quality VARCHAR(20) DEFAULT 'good' CHECK (connection_quality IN ('excellent', 'good', 'fair', 'poor')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraints
    CONSTRAINT fk_call_participants_session FOREIGN KEY (call_session_id) REFERENCES call_sessions(id) ON DELETE CASCADE,

    -- Constraints
    CONSTRAINT call_participants_joined_before_left CHECK (joined_at <= left_at OR left_at IS NULL),
    CONSTRAINT call_participants_unique_user_per_call UNIQUE (call_session_id, user_id)
);

-- Create user_presence table for real-time status
CREATE TABLE IF NOT EXISTS user_presence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'busy', 'away', 'offline')),
    in_call BOOLEAN DEFAULT false,
    current_call_id UUID,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraints
    CONSTRAINT fk_user_presence_call FOREIGN KEY (current_call_id) REFERENCES call_sessions(id) ON DELETE SET NULL
);

-- Create call_history table for call records and analytics
CREATE TABLE IF NOT EXISTS call_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    call_session_id UUID NOT NULL,
    caller_id UUID NOT NULL,
    callee_id UUID NOT NULL,
    call_type VARCHAR(10) NOT NULL CHECK (call_type IN ('voice', 'video')),
    call_status VARCHAR(20) NOT NULL CHECK (call_status IN ('completed', 'missed', 'declined', 'failed')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER, -- in seconds
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraints
    CONSTRAINT fk_call_history_session FOREIGN KEY (call_session_id) REFERENCES call_sessions(id) ON DELETE CASCADE,

    -- Constraints
    CONSTRAINT call_history_started_before_ended CHECK (started_at <= ended_at OR ended_at IS NULL),
    CONSTRAINT call_history_duration_positive CHECK (duration IS NULL OR duration >= 0)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_call_sessions_initiated_by ON call_sessions(initiated_by);
CREATE INDEX IF NOT EXISTS idx_call_sessions_status ON call_sessions(status);
CREATE INDEX IF NOT EXISTS idx_call_sessions_started_at ON call_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_call_sessions_type ON call_sessions(type);

CREATE INDEX IF NOT EXISTS idx_call_participants_call_session_id ON call_participants(call_session_id);
CREATE INDEX IF NOT EXISTS idx_call_participants_user_id ON call_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_call_participants_role ON call_participants(role);

CREATE INDEX IF NOT EXISTS idx_user_presence_user_id ON user_presence(user_id);
CREATE INDEX IF NOT EXISTS idx_user_presence_status ON user_presence(status);
CREATE INDEX IF NOT EXISTS idx_user_presence_in_call ON user_presence(in_call);
CREATE INDEX IF NOT EXISTS idx_user_presence_last_seen ON user_presence(last_seen);

CREATE INDEX IF NOT EXISTS idx_call_history_caller_id ON call_history(caller_id);
CREATE INDEX IF NOT EXISTS idx_call_history_callee_id ON call_history(callee_id);
CREATE INDEX IF NOT EXISTS idx_call_history_started_at ON call_history(started_at);
CREATE INDEX IF NOT EXISTS idx_call_history_call_status ON call_history(call_status);
CREATE INDEX IF NOT EXISTS idx_call_history_call_type ON call_history(call_type);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_call_sessions_updated_at
    BEFORE UPDATE ON call_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_call_participants_updated_at
    BEFORE UPDATE ON call_participants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_presence_updated_at
    BEFORE UPDATE ON user_presence
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();