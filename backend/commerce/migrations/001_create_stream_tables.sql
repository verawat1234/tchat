-- Migration: 001_create_stream_tables
-- Description: Creates all tables for Stream content management system
-- Created: 2025-09-29

-- Create stream_categories table
CREATE TABLE IF NOT EXISTS stream_categories (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    icon_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    featured_content_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Create stream_subtabs table
CREATE TABLE IF NOT EXISTS stream_subtabs (
    id VARCHAR(255) PRIMARY KEY,
    category_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    filter_criteria JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    FOREIGN KEY (category_id) REFERENCES stream_categories(id) ON DELETE CASCADE
);

-- Create stream_content_items table
CREATE TABLE IF NOT EXISTS stream_content_items (
    id VARCHAR(255) PRIMARY KEY,
    category_id VARCHAR(255) NOT NULL,
    subtab_id VARCHAR(255) DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    thumbnail_url VARCHAR(500) DEFAULT '',
    content_type VARCHAR(100) NOT NULL,
    content_url VARCHAR(500) DEFAULT '',
    duration_seconds INTEGER DEFAULT 0,
    file_size_bytes BIGINT DEFAULT 0,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    availability_status VARCHAR(50) NOT NULL DEFAULT 'available',
    is_featured BOOLEAN NOT NULL DEFAULT false,
    featured_order INTEGER DEFAULT NULL,
    view_count INTEGER NOT NULL DEFAULT 0,
    download_count INTEGER NOT NULL DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    review_count INTEGER NOT NULL DEFAULT 0,
    tags TEXT DEFAULT '',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    FOREIGN KEY (category_id) REFERENCES stream_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (subtab_id) REFERENCES stream_subtabs(id) ON DELETE SET NULL
);

-- Create tab_navigation_states table
CREATE TABLE IF NOT EXISTS tab_navigation_states (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    current_category_id VARCHAR(255) DEFAULT NULL,
    current_subtab_id VARCHAR(255) DEFAULT NULL,
    last_visited_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    session_id VARCHAR(255) NOT NULL,
    device_platform VARCHAR(50) NOT NULL DEFAULT 'web',
    autoplay_enabled BOOLEAN NOT NULL DEFAULT true,
    show_subtabs BOOLEAN NOT NULL DEFAULT true,
    preferred_view_mode VARCHAR(20) NOT NULL DEFAULT 'grid',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Create stream_user_sessions table
CREATE TABLE IF NOT EXISTS stream_user_sessions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    device_info VARCHAR(500) DEFAULT '',
    ip_address VARCHAR(45) DEFAULT '',
    user_agent VARCHAR(1000) DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    content_view_count INTEGER NOT NULL DEFAULT 0,
    total_time_spent BIGINT NOT NULL DEFAULT 0, -- in seconds
    categories_visited INTEGER NOT NULL DEFAULT 0,
    last_category_visited VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Create stream_content_views table
CREATE TABLE IF NOT EXISTS stream_content_views (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    content_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    view_started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    view_ended_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    duration BIGINT NOT NULL DEFAULT 0, -- in seconds
    view_progress DECIMAL(5,4) NOT NULL DEFAULT 0.0, -- 0.0 to 1.0
    is_completed BOOLEAN NOT NULL DEFAULT false,
    device_platform VARCHAR(50) NOT NULL DEFAULT 'web',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    FOREIGN KEY (content_id) REFERENCES stream_content_items(id) ON DELETE CASCADE
);

-- Create stream_user_preferences table
CREATE TABLE IF NOT EXISTS stream_user_preferences (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    preferred_categories TEXT DEFAULT '[]', -- JSON array
    blocked_categories TEXT DEFAULT '[]', -- JSON array
    autoplay_enabled BOOLEAN NOT NULL DEFAULT true,
    high_quality_preferred BOOLEAN NOT NULL DEFAULT false,
    offline_download_enabled BOOLEAN NOT NULL DEFAULT false,
    notifications_enabled BOOLEAN NOT NULL DEFAULT true,
    language_preference VARCHAR(10) NOT NULL DEFAULT 'en',
    region_preference VARCHAR(10) NOT NULL DEFAULT 'US',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance optimization

-- Indexes for stream_categories
CREATE INDEX IF NOT EXISTS idx_stream_categories_active ON stream_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_stream_categories_display_order ON stream_categories(display_order);
CREATE INDEX IF NOT EXISTS idx_stream_categories_deleted_at ON stream_categories(deleted_at);

-- Indexes for stream_subtabs
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_category_id ON stream_subtabs(category_id);
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_active ON stream_subtabs(is_active);
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_display_order ON stream_subtabs(display_order);
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_deleted_at ON stream_subtabs(deleted_at);

-- Indexes for stream_content_items
CREATE INDEX IF NOT EXISTS idx_stream_content_category_id ON stream_content_items(category_id);
CREATE INDEX IF NOT EXISTS idx_stream_content_subtab_id ON stream_content_items(subtab_id);
CREATE INDEX IF NOT EXISTS idx_stream_content_availability ON stream_content_items(availability_status);
CREATE INDEX IF NOT EXISTS idx_stream_content_featured ON stream_content_items(is_featured);
CREATE INDEX IF NOT EXISTS idx_stream_content_featured_order ON stream_content_items(featured_order);
CREATE INDEX IF NOT EXISTS idx_stream_content_view_count ON stream_content_items(view_count);
CREATE INDEX IF NOT EXISTS idx_stream_content_rating ON stream_content_items(rating);
CREATE INDEX IF NOT EXISTS idx_stream_content_created_at ON stream_content_items(created_at);
CREATE INDEX IF NOT EXISTS idx_stream_content_deleted_at ON stream_content_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_stream_content_type ON stream_content_items(content_type);
CREATE INDEX IF NOT EXISTS idx_stream_content_price ON stream_content_items(price);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_stream_content_category_available ON stream_content_items(category_id, availability_status);
CREATE INDEX IF NOT EXISTS idx_stream_content_category_featured ON stream_content_items(category_id, is_featured, featured_order);
CREATE INDEX IF NOT EXISTS idx_stream_content_subtab_available ON stream_content_items(subtab_id, availability_status);

-- Full-text search index for content search
CREATE INDEX IF NOT EXISTS idx_stream_content_search ON stream_content_items USING gin(to_tsvector('english', title || ' ' || description || ' ' || tags));

-- Indexes for tab_navigation_states
CREATE INDEX IF NOT EXISTS idx_navigation_user_id ON tab_navigation_states(user_id);
CREATE INDEX IF NOT EXISTS idx_navigation_session_id ON tab_navigation_states(session_id);
CREATE INDEX IF NOT EXISTS idx_navigation_last_visited ON tab_navigation_states(last_visited_at);
CREATE INDEX IF NOT EXISTS idx_navigation_deleted_at ON tab_navigation_states(deleted_at);

-- Indexes for stream_user_sessions
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON stream_user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON stream_user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_active ON stream_user_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_user_sessions_last_activity ON stream_user_sessions(last_activity_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_deleted_at ON stream_user_sessions(deleted_at);

-- Indexes for stream_content_views
CREATE INDEX IF NOT EXISTS idx_content_views_user_id ON stream_content_views(user_id);
CREATE INDEX IF NOT EXISTS idx_content_views_content_id ON stream_content_views(content_id);
CREATE INDEX IF NOT EXISTS idx_content_views_session_id ON stream_content_views(session_id);
CREATE INDEX IF NOT EXISTS idx_content_views_started_at ON stream_content_views(view_started_at);
CREATE INDEX IF NOT EXISTS idx_content_views_completed ON stream_content_views(is_completed);

-- Composite indexes for content views
CREATE INDEX IF NOT EXISTS idx_content_views_user_content ON stream_content_views(user_id, content_id);
CREATE INDEX IF NOT EXISTS idx_content_views_user_session ON stream_content_views(user_id, session_id);

-- Indexes for stream_user_preferences
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON stream_user_preferences(user_id);

-- Insert default stream categories
INSERT INTO stream_categories (id, name, display_order, icon_name, is_active, featured_content_enabled) VALUES
('books', 'Books', 1, 'book-open', true, true),
('podcasts', 'Podcasts', 2, 'microphone', true, true),
('cartoons', 'Cartoons', 3, 'film', true, true),
('movies', 'Movies', 4, 'video', true, true),
('music', 'Music', 5, 'music', true, true),
('art', 'Art', 6, 'palette', true, true)
ON CONFLICT (id) DO NOTHING;

-- Insert default subtabs for movies category
INSERT INTO stream_subtabs (id, category_id, name, display_order, filter_criteria, is_active) VALUES
('short-movies', 'movies', 'Short Films', 1, '{"maxDuration": 1800}', true),
('long-movies', 'movies', 'Feature Films', 2, '{"minDuration": 1801}', true)
ON CONFLICT (id) DO NOTHING;

-- Insert sample content for testing (optional, can be removed in production)
INSERT INTO stream_content_items (
    id, category_id, subtab_id, title, description, thumbnail_url, content_type,
    price, currency, availability_status, is_featured, featured_order, view_count, rating
) VALUES
-- Books
('book-001', 'books', NULL, 'The Art of Programming', 'A comprehensive guide to software development', '/images/book-001.jpg', 'book', 29.99, 'USD', 'available', true, 1, 1250, 4.5),
('book-002', 'books', NULL, 'Digital Minimalism', 'How to live better with less technology', '/images/book-002.jpg', 'book', 19.99, 'USD', 'available', true, 2, 890, 4.2),

-- Podcasts
('podcast-001', 'podcasts', NULL, 'Tech Talk Daily', 'Daily discussions about technology trends', '/images/podcast-001.jpg', 'podcast', 0.00, 'USD', 'available', true, 1, 2340, 4.7),
('podcast-002', 'podcasts', NULL, 'Business Insights', 'Weekly business strategy discussions', '/images/podcast-002.jpg', 'podcast', 9.99, 'USD', 'available', false, NULL, 1120, 4.3),

-- Movies
('movie-001', 'movies', 'short-movies', 'The Short Story', 'A 25-minute drama about human connection', '/images/movie-001.jpg', 'movie', 4.99, 'USD', 'available', true, 1, 567, 4.1),
('movie-002', 'movies', 'long-movies', 'Epic Adventure', 'A 2-hour action-packed adventure', '/images/movie-002.jpg', 'movie', 12.99, 'USD', 'available', true, 2, 2890, 4.6),

-- Music
('music-001', 'music', NULL, 'Relaxing Piano Collection', 'Peaceful piano melodies for relaxation', '/images/music-001.jpg', 'music', 14.99, 'USD', 'available', true, 1, 1780, 4.4),

-- Art
('art-001', 'art', NULL, 'Digital Art Masterclass', 'Learn digital art techniques from professionals', '/images/art-001.jpg', 'course', 49.99, 'USD', 'available', true, 1, 456, 4.8)

ON CONFLICT (id) DO NOTHING;

-- Add comments for documentation
COMMENT ON TABLE stream_categories IS 'Stores stream content categories (Books, Podcasts, Movies, etc.)';
COMMENT ON TABLE stream_subtabs IS 'Stores subcategories/subtabs for main categories';
COMMENT ON TABLE stream_content_items IS 'Stores individual content items available for streaming/purchase';
COMMENT ON TABLE tab_navigation_states IS 'Tracks user navigation state in the Stream tab';
COMMENT ON TABLE stream_user_sessions IS 'Tracks user sessions for stream content consumption';
COMMENT ON TABLE stream_content_views IS 'Tracks user interactions with specific content items';
COMMENT ON TABLE stream_user_preferences IS 'Stores user preferences for stream content';

COMMENT ON COLUMN stream_content_items.duration_seconds IS 'Content duration in seconds (for audio/video content)';
COMMENT ON COLUMN stream_content_items.file_size_bytes IS 'File size in bytes for downloadable content';
COMMENT ON COLUMN stream_content_items.view_count IS 'Number of times this content has been viewed';
COMMENT ON COLUMN stream_content_items.rating IS 'Average user rating (0.00 to 5.00)';
COMMENT ON COLUMN stream_content_items.tags IS 'Comma-separated tags for search and categorization';
COMMENT ON COLUMN stream_content_items.metadata IS 'Additional metadata stored as JSON';

COMMENT ON COLUMN stream_content_views.view_progress IS 'Progress through content (0.0 to 1.0, where 1.0 = completed)';
COMMENT ON COLUMN stream_content_views.duration IS 'Duration of this specific viewing session in seconds';

-- Create trigger to update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stream_categories_updated_at BEFORE UPDATE ON stream_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stream_subtabs_updated_at BEFORE UPDATE ON stream_subtabs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stream_content_items_updated_at BEFORE UPDATE ON stream_content_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tab_navigation_states_updated_at BEFORE UPDATE ON tab_navigation_states FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stream_user_sessions_updated_at BEFORE UPDATE ON stream_user_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stream_content_views_updated_at BEFORE UPDATE ON stream_content_views FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stream_user_preferences_updated_at BEFORE UPDATE ON stream_user_preferences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();