-- LiveStream table
CREATE TABLE IF NOT EXISTS live_streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    broadcaster_id UUID NOT NULL,
    stream_type VARCHAR(20) NOT NULL CHECK (stream_type IN ('store', 'video')),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    broadcaster_kyc_tier INT NOT NULL CHECK (broadcaster_kyc_tier >= 0 AND broadcaster_kyc_tier <= 3),
    status VARCHAR(20) NOT NULL CHECK (status IN ('scheduled', 'live', 'ended')),
    stream_key VARCHAR(100) UNIQUE NOT NULL,
    viewer_count INT DEFAULT 0,
    peak_viewer_count INT DEFAULT 0,
    max_capacity INT DEFAULT 50000,
    total_view_time INT DEFAULT 0,
    average_watch_time DOUBLE PRECISION,
    scheduled_start_time TIMESTAMP WITH TIME ZONE,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    simulcast_layers TEXT[],
    current_bitrate INT,
    recording_url TEXT,
    thumbnail_url TEXT,
    featured_products UUID[],
    stream_settings JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_live_streams_broadcaster_id ON live_streams(broadcaster_id);
CREATE INDEX idx_live_streams_status ON live_streams(status);
CREATE INDEX idx_live_streams_stream_type ON live_streams(stream_type);
CREATE INDEX idx_live_streams_created_at ON live_streams(created_at DESC);

-- ViewerSession table
CREATE TABLE IF NOT EXISTS viewer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    viewer_id UUID,
    join_time TIMESTAMP WITH TIME ZONE NOT NULL,
    leave_time TIMESTAMP WITH TIME ZONE,
    watch_duration INT,
    peak_quality_layer VARCHAR(10),
    total_rebuffer_events INT DEFAULT 0,
    total_rebuffer_duration INT DEFAULT 0,
    viewer_ip VARCHAR(45),
    viewer_country VARCHAR(2),
    viewer_device VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_viewer_sessions_stream_id ON viewer_sessions(stream_id);
CREATE INDEX idx_viewer_sessions_viewer_id ON viewer_sessions(viewer_id);

-- StreamProduct table
CREATE TABLE IF NOT EXISTS stream_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    display_order INT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    view_count INT DEFAULT 0,
    click_count INT DEFAULT 0,
    purchase_count INT DEFAULT 0,
    revenue_generated DOUBLE PRECISION DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_stream_products_stream_id ON stream_products(stream_id);
CREATE INDEX idx_stream_products_product_id ON stream_products(product_id);

-- NotificationPreference table
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    live_start_enabled BOOLEAN DEFAULT TRUE,
    live_start_channels TEXT[],
    chat_mention_enabled BOOLEAN DEFAULT TRUE,
    chat_mention_channels TEXT[],
    featured_product_enabled BOOLEAN DEFAULT FALSE,
    featured_product_channels TEXT[],
    milestone_enabled BOOLEAN DEFAULT TRUE,
    milestone_channels TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);

-- StreamAnalytics table
CREATE TABLE IF NOT EXISTS stream_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL UNIQUE REFERENCES live_streams(id) ON DELETE CASCADE,
    total_viewers INT DEFAULT 0,
    peak_viewer_count INT DEFAULT 0,
    average_viewer_count INT DEFAULT 0,
    total_view_time_seconds INT DEFAULT 0,
    average_watch_time_seconds DOUBLE PRECISION DEFAULT 0.0,
    total_chat_messages INT DEFAULT 0,
    total_reactions INT DEFAULT 0,
    reaction_breakdown JSONB,
    featured_products_count INT DEFAULT 0,
    total_product_views INT DEFAULT 0,
    total_product_clicks INT DEFAULT 0,
    total_product_purchases INT DEFAULT 0,
    total_revenue_generated DOUBLE PRECISION DEFAULT 0.0,
    viewer_demographics JSONB,
    traffic_sources JSONB,
    engagement_score DOUBLE PRECISION DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_stream_analytics_stream_id ON stream_analytics(stream_id);