-- Create keyspace for streaming data
CREATE KEYSPACE IF NOT EXISTS streaming
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 3
};

USE streaming;

-- ChatMessage table with 30-day TTL
CREATE TABLE IF NOT EXISTS chat_messages (
    stream_id UUID,
    timestamp TIMESTAMP,
    message_id UUID,
    user_id UUID,
    message TEXT,
    moderation_status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (stream_id, timestamp, message_id)
) WITH CLUSTERING ORDER BY (timestamp DESC)
  AND default_time_to_live = 2592000
  AND gc_grace_seconds = 86400;

-- StreamReaction table with 1-day TTL (real-time aggregation)
CREATE TABLE IF NOT EXISTS stream_reactions (
    stream_id UUID,
    timestamp TIMESTAMP,
    reaction_id UUID,
    user_id UUID,
    reaction_type TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (stream_id, timestamp, reaction_id)
) WITH CLUSTERING ORDER BY (timestamp DESC)
  AND default_time_to_live = 86400
  AND gc_grace_seconds = 3600;