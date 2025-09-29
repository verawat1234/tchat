-- Migration: Create Stream Tables
-- Description: Creates tables for Stream Store Tabs content management system
-- Date: 2025-09-29
-- Dependencies: Commerce service database

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- STREAM CATEGORIES TABLE
-- ============================================================================
-- Represents different types of Stream content categories (Books, Podcasts, etc.)
CREATE TABLE IF NOT EXISTS stream_categories (
    id VARCHAR(50) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    display_order INTEGER NOT NULL,
    icon_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    featured_content_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stream_categories_display_order ON stream_categories(display_order);
CREATE INDEX IF NOT EXISTS idx_stream_categories_active ON stream_categories(is_active);

-- ============================================================================
-- STREAM SUBTABS TABLE
-- ============================================================================
-- Represents category subdivisions (e.g., Fiction/Non-Fiction for Books)
CREATE TABLE IF NOT EXISTS stream_subtabs (
    id TEXT NOT NULL PRIMARY KEY,
    category_id VARCHAR(50) NOT NULL,
    name TEXT NOT NULL,
    display_order INTEGER NOT NULL,
    filter_criteria TEXT NOT NULL, -- JSON string for filtering rules
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (category_id) REFERENCES stream_categories(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_category ON stream_subtabs(category_id);
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_display_order ON stream_subtabs(category_id, display_order);
CREATE INDEX IF NOT EXISTS idx_stream_subtabs_active ON stream_subtabs(is_active);

-- ============================================================================
-- STREAM CONTENT ITEMS TABLE
-- ============================================================================
-- Represents individual content items (books, podcasts, movies, etc.)
CREATE TABLE IF NOT EXISTS stream_content_items (
    id TEXT NOT NULL PRIMARY KEY,
    category_id VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL, -- BOOK, PODCAST, CARTOON, SHORT_MOVIE, LONG_MOVIE, MUSIC, ART
    duration INTEGER, -- in seconds, null for books
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    availability_status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE', -- AVAILABLE, COMING_SOON, UNAVAILABLE
    is_featured BOOLEAN NOT NULL DEFAULT false,
    featured_order INTEGER,
    metadata TEXT NOT NULL DEFAULT '{}', -- JSON string
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at TIMESTAMP,

    FOREIGN KEY (category_id) REFERENCES stream_categories(id) ON DELETE CASCADE
);

-- Create indexes for performance and queries
CREATE INDEX IF NOT EXISTS idx_stream_content_category ON stream_content_items(category_id);
CREATE INDEX IF NOT EXISTS idx_stream_content_featured ON stream_content_items(is_featured, featured_order);
CREATE INDEX IF NOT EXISTS idx_stream_content_type ON stream_content_items(content_type);
CREATE INDEX IF NOT EXISTS idx_stream_content_availability ON stream_content_items(availability_status);
CREATE INDEX IF NOT EXISTS idx_stream_content_sync ON stream_content_items(synced_at);
CREATE INDEX IF NOT EXISTS idx_stream_content_price ON stream_content_items(price);

-- ============================================================================
-- TAB NAVIGATION STATE TABLE
-- ============================================================================
-- Tracks user's current position in Stream system
CREATE TABLE IF NOT EXISTS tab_navigation_states (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    current_category_id VARCHAR(255),
    current_subtab_id VARCHAR(255),
    last_visited_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    session_id VARCHAR(255) NOT NULL,
    device_platform VARCHAR(50) DEFAULT 'web',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Navigation preferences
    autoplay_enabled BOOLEAN DEFAULT true,
    show_subtabs BOOLEAN DEFAULT true,
    preferred_view_mode VARCHAR(20) DEFAULT 'grid' -- grid, list
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_tab_navigation_user ON tab_navigation_states(user_id);
CREATE INDEX IF NOT EXISTS idx_tab_navigation_session ON tab_navigation_states(session_id);
CREATE INDEX IF NOT EXISTS idx_tab_navigation_deleted ON tab_navigation_states(deleted_at);

-- ============================================================================
-- STREAM USER SESSIONS TABLE
-- ============================================================================
-- Tracks user sessions for Stream content consumption
CREATE TABLE IF NOT EXISTS stream_user_sessions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device_info VARCHAR(500),
    ip_address VARCHAR(45),
    user_agent VARCHAR(1000),
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Session statistics
    content_view_count INTEGER DEFAULT 0,
    total_time_spent BIGINT DEFAULT 0, -- in seconds
    categories_visited INTEGER DEFAULT 0,
    last_category_visited VARCHAR(255)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stream_sessions_user ON stream_user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_stream_sessions_token ON stream_user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_stream_sessions_active ON stream_user_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_stream_sessions_deleted ON stream_user_sessions(deleted_at);

-- ============================================================================
-- STREAM CONTENT VIEWS TABLE
-- ============================================================================
-- Tracks individual content viewing sessions
CREATE TABLE IF NOT EXISTS stream_content_views (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    content_id TEXT NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    view_progress DECIMAL(5,2) DEFAULT 0.0, -- Percentage (0.0 to 100.0)
    view_duration BIGINT DEFAULT 0, -- in seconds
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    device_platform VARCHAR(50) DEFAULT 'web',
    quality_settings VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    FOREIGN KEY (content_id) REFERENCES stream_content_items(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stream_views_user ON stream_content_views(user_id);
CREATE INDEX IF NOT EXISTS idx_stream_views_content ON stream_content_views(content_id);
CREATE INDEX IF NOT EXISTS idx_stream_views_session ON stream_content_views(session_id);
CREATE INDEX IF NOT EXISTS idx_stream_views_progress ON stream_content_views(view_progress);
CREATE INDEX IF NOT EXISTS idx_stream_views_deleted ON stream_content_views(deleted_at);

-- ============================================================================
-- STREAM USER PREFERENCES TABLE
-- ============================================================================
-- Stores user preferences for Stream content consumption
CREATE TABLE IF NOT EXISTS stream_user_preferences (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    preferred_quality VARCHAR(20) DEFAULT 'auto', -- auto, low, medium, high, ultra
    autoplay_enabled BOOLEAN DEFAULT true,
    subtitles_enabled BOOLEAN DEFAULT false,
    preferred_language VARCHAR(10) DEFAULT 'en',
    content_filters TEXT DEFAULT '{}', -- JSON string for content filtering preferences
    notification_settings TEXT DEFAULT '{}', -- JSON string for notification preferences
    privacy_settings TEXT DEFAULT '{}', -- JSON string for privacy preferences
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stream_preferences_user ON stream_user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_stream_preferences_deleted ON stream_user_preferences(deleted_at);

-- ============================================================================
-- SEED DATA - Default Stream Categories
-- ============================================================================
-- Insert default Stream categories for the Stream Store Tabs
INSERT INTO stream_categories (id, name, display_order, icon_name, is_active, featured_content_enabled) VALUES
    ('books', 'Books', 1, 'book-open', true, true),
    ('podcasts', 'Podcasts', 2, 'headphones', true, true),
    ('cartoons', 'Cartoons', 3, 'film', true, true),
    ('movies', 'Movies', 4, 'video', true, true),
    ('music', 'Music', 5, 'music', true, true),
    ('art', 'Art', 6, 'palette', true, true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- SEED DATA - Default Subtabs for Books Category
-- ============================================================================
-- Insert default subtabs for Books category
INSERT INTO stream_subtabs (id, category_id, name, display_order, filter_criteria, is_active) VALUES
    ('books_fiction', 'books', 'Fiction', 1, '{"genre": "fiction"}', true),
    ('books_nonfiction', 'books', 'Non-Fiction', 2, '{"genre": "non-fiction"}', true),
    ('books_academic', 'books', 'Academic', 3, '{"genre": "academic"}', true),
    ('books_childrens', 'books', 'Children''s', 4, '{"genre": "childrens"}', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- SEED DATA - Default Subtabs for Movies Category
-- ============================================================================
-- Insert default subtabs for Movies category with duration filtering
INSERT INTO stream_subtabs (id, category_id, name, display_order, filter_criteria, is_active) VALUES
    ('movies_short', 'movies', 'Short Films', 1, '{"content_type": "SHORT_MOVIE", "max_duration": 1800}', true),
    ('movies_feature', 'movies', 'Feature Films', 2, '{"content_type": "LONG_MOVIE", "min_duration": 1800}', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- CREATE UPDATED_AT TRIGGERS
-- ============================================================================
-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for all tables with updated_at columns
CREATE TRIGGER update_stream_categories_updated_at
    BEFORE UPDATE ON stream_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stream_subtabs_updated_at
    BEFORE UPDATE ON stream_subtabs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stream_content_items_updated_at
    BEFORE UPDATE ON stream_content_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tab_navigation_states_updated_at
    BEFORE UPDATE ON tab_navigation_states
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stream_user_sessions_updated_at
    BEFORE UPDATE ON stream_user_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stream_content_views_updated_at
    BEFORE UPDATE ON stream_content_views
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stream_user_preferences_updated_at
    BEFORE UPDATE ON stream_user_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- PERFORMANCE OPTIMIZATIONS
-- ============================================================================
-- Create composite indexes for common query patterns

-- Content search by category and availability
CREATE INDEX IF NOT EXISTS idx_stream_content_category_available
    ON stream_content_items(category_id, availability_status)
    WHERE availability_status = 'AVAILABLE';

-- Featured content by category
CREATE INDEX IF NOT EXISTS idx_stream_content_featured_by_category
    ON stream_content_items(category_id, is_featured, featured_order)
    WHERE is_featured = true;

-- User navigation by active sessions
CREATE INDEX IF NOT EXISTS idx_navigation_active_users
    ON tab_navigation_states(user_id, last_visited_at)
    WHERE deleted_at IS NULL;

-- Content views for analytics
CREATE INDEX IF NOT EXISTS idx_content_views_analytics
    ON stream_content_views(content_id, started_at, view_progress)
    WHERE deleted_at IS NULL;

-- ============================================================================
-- MIGRATION VERIFICATION
-- ============================================================================
-- Verify all tables were created successfully
DO $$
DECLARE
    table_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_schema = 'public'
    AND table_name IN (
        'stream_categories',
        'stream_subtabs',
        'stream_content_items',
        'tab_navigation_states',
        'stream_user_sessions',
        'stream_content_views',
        'stream_user_preferences'
    );

    IF table_count = 7 THEN
        RAISE NOTICE 'SUCCESS: All 7 Stream tables created successfully';
    ELSE
        RAISE EXCEPTION 'FAILURE: Expected 7 tables, found %', table_count;
    END IF;
END $$;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Stream Store Tabs database schema migration completed successfully
-- Tables created: 7
-- Indexes created: 20+
-- Triggers created: 7
-- Seed data: 6 categories + 6 subtabs