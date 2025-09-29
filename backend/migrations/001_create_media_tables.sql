-- Migration: Create Media Store Tabs Tables
-- Date: 2025-09-29
-- Feature: Media Store Tabs (026-help-me-add)

-- Create MediaCategory table
CREATE TABLE IF NOT EXISTS media_categories (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    display_order INTEGER NOT NULL,
    icon_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    featured_content_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT media_categories_display_order_positive CHECK (display_order > 0),
    CONSTRAINT media_categories_name_not_empty CHECK (LENGTH(name) > 0)
);

-- Create MediaSubtab table
CREATE TABLE IF NOT EXISTS media_subtabs (
    id VARCHAR(50) PRIMARY KEY,
    category_id VARCHAR(50) NOT NULL,
    name VARCHAR(30) NOT NULL,
    display_order INTEGER NOT NULL,
    filter_criteria JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (category_id) REFERENCES media_categories(id) ON DELETE CASCADE,
    CONSTRAINT media_subtabs_display_order_positive CHECK (display_order > 0),
    CONSTRAINT media_subtabs_name_not_empty CHECK (LENGTH(name) > 0)
);

-- Create MediaContentItem table
CREATE TABLE IF NOT EXISTS media_content_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    duration INTEGER NULL, -- in seconds, null for books
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    availability_status VARCHAR(20) NOT NULL DEFAULT 'available',
    is_featured BOOLEAN NOT NULL DEFAULT false,
    featured_order INTEGER NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (category_id) REFERENCES media_categories(id) ON DELETE CASCADE,
    CONSTRAINT media_content_items_title_not_empty CHECK (LENGTH(title) > 0),
    CONSTRAINT media_content_items_price_positive CHECK (price >= 0),
    CONSTRAINT media_content_items_duration_positive CHECK (duration IS NULL OR duration > 0),
    CONSTRAINT media_content_items_featured_order_required CHECK (
        (is_featured = false) OR (is_featured = true AND featured_order IS NOT NULL)
    ),
    CONSTRAINT media_content_items_content_type_valid CHECK (
        content_type IN ('book', 'podcast', 'cartoon', 'short_movie', 'long_movie')
    ),
    CONSTRAINT media_content_items_availability_valid CHECK (
        availability_status IN ('available', 'coming_soon', 'unavailable')
    )
);

-- Create indexes for performance
CREATE INDEX idx_media_categories_display_order ON media_categories(display_order);
CREATE INDEX idx_media_categories_active ON media_categories(is_active);

CREATE INDEX idx_media_subtabs_category_id ON media_subtabs(category_id);
CREATE INDEX idx_media_subtabs_display_order ON media_subtabs(category_id, display_order);
CREATE INDEX idx_media_subtabs_active ON media_subtabs(is_active);

CREATE INDEX idx_media_content_items_category_id ON media_content_items(category_id);
CREATE INDEX idx_media_content_items_featured ON media_content_items(is_featured, featured_order);
CREATE INDEX idx_media_content_items_availability ON media_content_items(availability_status);
CREATE INDEX idx_media_content_items_content_type ON media_content_items(content_type);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_media_categories_updated_at BEFORE UPDATE
    ON media_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_media_subtabs_updated_at BEFORE UPDATE
    ON media_subtabs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_media_content_items_updated_at BEFORE UPDATE
    ON media_content_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();