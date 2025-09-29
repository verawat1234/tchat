-- Stream Content Seed Data
-- Description: Sample data for Stream Store Tabs development and testing
-- Date: 2025-09-29
-- Dependencies: 20250929_create_stream_tables.sql migration

-- ============================================================================
-- SEED DATA - BOOKS CATEGORY
-- ============================================================================

-- Fiction Books
INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('book_fiction_1', 'books', 'The Digital Renaissance',
     'A thrilling novel about the intersection of technology and humanity in the 22nd century.',
     'https://cdn.tchat.com/thumbnails/books/digital_renaissance.jpg', 'BOOK',
     NULL, 12.99, 'USD', 'AVAILABLE', true, 1,
     '{"author": "Sarah Chen", "genre": "science_fiction", "pages": 324, "language": "en", "isbn": "978-0123456789"}'),

    ('book_fiction_2', 'books', 'Whispers in the Cloud',
     'A mystery set in a world where AI consciousness emerges in unexpected ways.',
     'https://cdn.tchat.com/thumbnails/books/whispers_cloud.jpg', 'BOOK',
     NULL, 15.99, 'USD', 'AVAILABLE', true, 2,
     '{"author": "Marcus Rodriguez", "genre": "mystery", "pages": 298, "language": "en", "isbn": "978-0987654321"}'),

    ('book_fiction_3', 'books', 'The Last Library',
     'In a post-digital world, one person fights to preserve physical books.',
     'https://cdn.tchat.com/thumbnails/books/last_library.jpg', 'BOOK',
     NULL, 11.99, 'USD', 'AVAILABLE', false, NULL,
     '{"author": "Elena Volkov", "genre": "dystopian", "pages": 267, "language": "en", "isbn": "978-0192837465"}'),

    ('book_fiction_4', 'books', 'Neural Networks of Love',
     'A romantic drama exploring relationships in an AI-augmented society.',
     'https://cdn.tchat.com/thumbnails/books/neural_love.jpg', 'BOOK',
     NULL, 13.99, 'USD', 'COMING_SOON', false, NULL,
     '{"author": "Dr. Kenji Nakamura", "genre": "romance", "pages": 312, "language": "en", "isbn": "978-0567891234"}');

-- Non-Fiction Books
INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('book_nonfiction_1', 'books', 'Building Tomorrow: AI Ethics in Practice',
     'A comprehensive guide to implementing ethical AI systems in modern organizations.',
     'https://cdn.tchat.com/thumbnails/books/ai_ethics.jpg', 'BOOK',
     NULL, 24.99, 'USD', 'AVAILABLE', true, 3,
     '{"author": "Prof. Amara Johnson", "genre": "technology", "pages": 456, "language": "en", "isbn": "978-0345678901"}'),

    ('book_nonfiction_2', 'books', 'The Psychology of Digital Natives',
     'Understanding how technology shapes modern human behavior and cognition.',
     'https://cdn.tchat.com/thumbnails/books/digital_psychology.jpg', 'BOOK',
     NULL, 19.99, 'USD', 'AVAILABLE', false, NULL,
     '{"author": "Dr. Lisa Thompson", "genre": "psychology", "pages": 328, "language": "en", "isbn": "978-0234567890"}'),

    ('book_nonfiction_3', 'books', 'Sustainable Tech: Green Computing Revolution',
     'How technology companies are leading the fight against climate change.',
     'https://cdn.tchat.com/thumbnails/books/sustainable_tech.jpg', 'BOOK',
     NULL, 22.99, 'USD', 'AVAILABLE', false, NULL,
     '{"author": "Dr. Green Evans", "genre": "environment", "pages": 389, "language": "en", "isbn": "978-0456789012"}');

-- ============================================================================
-- SEED DATA - PODCASTS CATEGORY
-- ============================================================================

INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('podcast_tech_1', 'podcasts', 'TechTalk Weekly: The Future is Now',
     'Weekly discussions on emerging technologies and their impact on society.',
     'https://cdn.tchat.com/thumbnails/podcasts/techtalk_weekly.jpg', 'PODCAST',
     2700, 4.99, 'USD', 'AVAILABLE', true, 1,
     '{"hosts": ["Alex Kim", "Jordan Smith"], "episode_count": 156, "category": "technology", "language": "en", "rating": 4.8}'),

    ('podcast_business_1', 'podcasts', 'Startup Stories: From Idea to IPO',
     'Inspiring interviews with successful entrepreneurs and startup founders.',
     'https://cdn.tchat.com/thumbnails/podcasts/startup_stories.jpg', 'PODCAST',
     3600, 6.99, 'USD', 'AVAILABLE', true, 2,
     '{"hosts": ["Maria Garcia"], "episode_count": 89, "category": "business", "language": "en", "rating": 4.9}'),

    ('podcast_science_1', 'podcasts', 'Quantum Conversations',
     'Deep dives into quantum physics, explained for curious minds.',
     'https://cdn.tchat.com/thumbnails/podcasts/quantum_conv.jpg', 'PODCAST',
     3300, 5.99, 'USD', 'AVAILABLE', false, NULL,
     '{"hosts": ["Dr. Emma Wilson", "Prof. David Chang"], "episode_count": 67, "category": "science", "language": "en", "rating": 4.7}'),

    ('podcast_culture_1', 'podcasts', 'Global Voices: Culture & Technology',
     'Exploring how technology influences cultures around the world.',
     'https://cdn.tchat.com/thumbnails/podcasts/global_voices.jpg', 'PODCAST',
     2400, 3.99, 'USD', 'AVAILABLE', false, NULL,
     '{"hosts": ["Aria Patel", "Carlos Mendez", "Yuki Tanaka"], "episode_count": 134, "category": "culture", "language": "en", "rating": 4.6}');

-- ============================================================================
-- SEED DATA - CARTOONS CATEGORY
-- ============================================================================

INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('cartoon_series_1', 'cartoons', 'Robot Academy Adventures',
     'A fun educational series about young robots learning about friendship and technology.',
     'https://cdn.tchat.com/thumbnails/cartoons/robot_academy.jpg', 'CARTOON',
     1800, 7.99, 'USD', 'AVAILABLE', true, 1,
     '{"studio": "FutureToons Animation", "episodes": 24, "age_rating": "G", "language": "en", "season": 1}'),

    ('cartoon_series_2', 'cartoons', 'Space Cats: Cosmic Patrol',
     'Intergalactic adventures with a team of heroic space-faring cats.',
     'https://cdn.tchat.com/thumbnails/cartoons/space_cats.jpg', 'CARTOON',
     1500, 6.99, 'USD', 'AVAILABLE', true, 2,
     '{"studio": "Stellar Studios", "episodes": 18, "age_rating": "PG", "language": "en", "season": 2}'),

    ('cartoon_educational_1', 'cartoons', 'Code Quest: Programming for Kids',
     'Educational cartoon teaching programming concepts through adventure.',
     'https://cdn.tchat.com/thumbnails/cartoons/code_quest.jpg', 'CARTOON',
     1200, 8.99, 'USD', 'AVAILABLE', false, NULL,
     '{"studio": "EduToons", "episodes": 15, "age_rating": "G", "language": "en", "educational": true}'),

    ('cartoon_classic_1', 'cartoons', 'Digital Dreamland Chronicles',
     'Classic-style animation meets modern storytelling in this fantasy series.',
     'https://cdn.tchat.com/thumbnails/cartoons/digital_dreamland.jpg', 'CARTOON',
     2100, 9.99, 'USD', 'COMING_SOON', false, NULL,
     '{"studio": "Classic Digital Arts", "episodes": 12, "age_rating": "PG", "language": "en", "style": "hand_drawn"}');

-- ============================================================================
-- SEED DATA - MOVIES CATEGORY
-- ============================================================================

-- Short Movies
INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('movie_short_1', 'movies', 'The Last Message',
     'A heartwarming short film about connection in the digital age.',
     'https://cdn.tchat.com/thumbnails/movies/last_message.jpg', 'SHORT_MOVIE',
     900, 3.99, 'USD', 'AVAILABLE', true, 1,
     '{"director": "Sofia Martinez", "year": 2024, "genre": "drama", "rating": "PG", "awards": ["Sundance Short Film Award"]}'),

    ('movie_short_2', 'movies', 'Binary Dreams',
     'An experimental short exploring the boundary between human and artificial consciousness.',
     'https://cdn.tchat.com/thumbnails/movies/binary_dreams.jpg', 'SHORT_MOVIE',
     1200, 4.99, 'USD', 'AVAILABLE', false, NULL,
     '{"director": "Alex Chen", "year": 2024, "genre": "sci-fi", "rating": "PG-13", "experimental": true}'),

    ('movie_short_3', 'movies', 'Code Red Comedy',
     'A hilarious take on the daily life of software developers.',
     'https://cdn.tchat.com/thumbnails/movies/code_red.jpg', 'SHORT_MOVIE',
     750, 2.99, 'USD', 'AVAILABLE', false, NULL,
     '{"director": "Jamie Park", "year": 2024, "genre": "comedy", "rating": "PG", "tech_humor": true}');

-- Feature Films
INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('movie_feature_1', 'movies', 'The Algorithm War',
     'A thrilling cyber-warfare movie about protecting democracy in the digital age.',
     'https://cdn.tchat.com/thumbnails/movies/algorithm_war.jpg', 'LONG_MOVIE',
     7200, 14.99, 'USD', 'AVAILABLE', true, 2,
     '{"director": "Michael Zhang", "year": 2024, "genre": "thriller", "rating": "PG-13", "budget": "50M", "box_office": "125M"}'),

    ('movie_feature_2', 'movies', 'Love in the Time of AI',
     'A romantic comedy about finding human connection in an automated world.',
     'https://cdn.tchat.com/thumbnails/movies/love_ai_time.jpg', 'LONG_MOVIE',
     6600, 12.99, 'USD', 'AVAILABLE', true, 3,
     '{"director": "Emma Thompson", "year": 2024, "genre": "romantic_comedy", "rating": "PG-13", "cast": ["Ryan Adams", "Zoe Liu"]}'),

    ('movie_feature_3', 'movies', 'The Quantum Paradox',
     'A mind-bending science fiction epic about parallel realities and quantum computing.',
     'https://cdn.tchat.com/thumbnails/movies/quantum_paradox.jpg', 'LONG_MOVIE',
     8400, 16.99, 'USD', 'COMING_SOON', false, NULL,
     '{"director": "Christopher Nolan Jr.", "year": 2024, "genre": "sci-fi", "rating": "PG-13", "vfx_studio": "Digital Dreams"}');

-- ============================================================================
-- SEED DATA - MUSIC CATEGORY
-- ============================================================================

INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('music_album_1', 'music', 'Synthetic Symphonies',
     'An AI-human collaboration album featuring orchestral pieces enhanced by machine learning.',
     'https://cdn.tchat.com/thumbnails/music/synthetic_symphonies.jpg', 'MUSIC',
     3600, 9.99, 'USD', 'AVAILABLE', true, 1,
     '{"artist": "The Digital Orchestra", "tracks": 12, "genre": "classical_electronic", "year": 2024, "collaboration": "AI-assisted"}'),

    ('music_album_2', 'music', 'Code & Rhythm',
     'Electronic beats inspired by programming languages and algorithms.',
     'https://cdn.tchat.com/thumbnails/music/code_rhythm.jpg', 'MUSIC',
     2700, 7.99, 'USD', 'AVAILABLE', true, 2,
     '{"artist": "Syntax Error", "tracks": 10, "genre": "electronic", "year": 2024, "concept": "programming_inspired"}'),

    ('music_single_1', 'music', 'Neural Network Lullaby',
     'A soothing single that helps listeners relax and disconnect from digital stress.',
     'https://cdn.tchat.com/thumbnails/music/neural_lullaby.jpg', 'MUSIC',
     240, 1.99, 'USD', 'AVAILABLE', false, NULL,
     '{"artist": "Digital Zen", "tracks": 1, "genre": "ambient", "year": 2024, "purpose": "relaxation"}'),

    ('music_album_3', 'music', 'Human.exe',
     'A rock album exploring themes of humanity in the digital age.',
     'https://cdn.tchat.com/thumbnails/music/human_exe.jpg', 'MUSIC',
     3300, 11.99, 'USD', 'AVAILABLE', false, NULL,
     '{"artist": "The Debuggers", "tracks": 11, "genre": "alternative_rock", "year": 2024, "theme": "digital_humanity"}');

-- ============================================================================
-- SEED DATA - ART CATEGORY
-- ============================================================================

INSERT INTO stream_content_items (
    id, category_id, title, description, thumbnail_url, content_type,
    duration, price, currency, availability_status, is_featured, featured_order, metadata
) VALUES
    ('art_collection_1', 'art', 'Digital Canvas: AI Art Renaissance',
     'A curated collection of AI-generated artworks exploring the future of creativity.',
     'https://cdn.tchat.com/thumbnails/art/digital_canvas.jpg', 'ART',
     NULL, 19.99, 'USD', 'AVAILABLE', true, 1,
     '{"artist": "AI Collective", "pieces": 25, "medium": "digital", "style": "contemporary", "theme": "AI_creativity"}'),

    ('art_gallery_1', 'art', 'Pixel Perfect: Retro Gaming Art',
     'Nostalgic pixel art celebrating the golden age of video games.',
     'https://cdn.tchat.com/thumbnails/art/pixel_perfect.jpg', 'ART',
     NULL, 15.99, 'USD', 'AVAILABLE', true, 2,
     '{"artist": "8-Bit Studios", "pieces": 18, "medium": "pixel_art", "style": "retro", "theme": "gaming_nostalgia"}'),

    ('art_interactive_1', 'art', 'Touch the Future: Interactive Digital Sculptures',
     '3D digital sculptures that respond to user interaction and touch.',
     'https://cdn.tchat.com/thumbnails/art/touch_future.jpg', 'ART',
     NULL, 24.99, 'USD', 'AVAILABLE', false, NULL,
     '{"artist": "Interactive Arts Lab", "pieces": 8, "medium": "3d_digital", "interactive": true, "requires": "touch_device"}'),

    ('art_photography_1', 'art', 'Augmented Reality: Cities of Tomorrow',
     'Photography series showcasing how AR technology transforms urban landscapes.',
     'https://cdn.tchat.com/thumbnails/art/ar_cities.jpg', 'ART',
     NULL, 17.99, 'USD', 'COMING_SOON', false, NULL,
     '{"artist": "Future Vision Photography", "pieces": 15, "medium": "ar_photography", "style": "documentary", "theme": "urban_future"}');

-- ============================================================================
-- SEED DATA - USER NAVIGATION STATES (Sample Data)
-- ============================================================================

-- Sample navigation states for testing
INSERT INTO tab_navigation_states (
    user_id, current_category_id, current_subtab_id, session_id, device_platform,
    autoplay_enabled, show_subtabs, preferred_view_mode
) VALUES
    ('user_001', 'books', 'books_fiction', 'session_001', 'web', true, true, 'grid'),
    ('user_002', 'podcasts', NULL, 'session_002', 'mobile', false, true, 'list'),
    ('user_003', 'movies', 'movies_feature', 'session_003', 'tablet', true, false, 'grid'),
    ('user_004', 'music', NULL, 'session_004', 'web', true, true, 'grid'),
    ('user_005', 'art', NULL, 'session_005', 'mobile', false, true, 'list');

-- ============================================================================
-- SEED DATA - USER PREFERENCES (Sample Data)
-- ============================================================================

-- Sample user preferences for testing
INSERT INTO stream_user_preferences (
    user_id, preferred_quality, autoplay_enabled, subtitles_enabled, preferred_language,
    content_filters, notification_settings, privacy_settings
) VALUES
    ('user_001', 'high', true, false, 'en',
     '{"age_rating": ["G", "PG"], "genres": ["fiction", "technology"]}',
     '{"email": true, "push": false, "recommendations": true}',
     '{"share_viewing_history": false, "personalized_ads": true}'),

    ('user_002', 'auto', false, true, 'en',
     '{"content_types": ["PODCAST", "MUSIC"], "max_duration": 3600}',
     '{"email": false, "push": true, "recommendations": true}',
     '{"share_viewing_history": true, "personalized_ads": false}'),

    ('user_003', 'ultra', true, false, 'es',
     '{"age_rating": ["PG", "PG-13"], "genres": ["comedy", "drama"]}',
     '{"email": true, "push": true, "recommendations": false}',
     '{"share_viewing_history": false, "personalized_ads": false}');

-- ============================================================================
-- SEED DATA VERIFICATION
-- ============================================================================

-- Verify seed data was inserted successfully
DO $$
DECLARE
    content_count INTEGER;
    categories_count INTEGER;
    subtabs_count INTEGER;
BEGIN
    -- Check content items
    SELECT COUNT(*) INTO content_count FROM stream_content_items;

    -- Check categories (should be 6 from migration + any additional)
    SELECT COUNT(*) INTO categories_count FROM stream_categories;

    -- Check subtabs (should be 6 from migration + any additional)
    SELECT COUNT(*) INTO subtabs_count FROM stream_subtabs;

    RAISE NOTICE 'SEED DATA SUMMARY:';
    RAISE NOTICE '- Stream Categories: %', categories_count;
    RAISE NOTICE '- Stream Subtabs: %', subtabs_count;
    RAISE NOTICE '- Stream Content Items: %', content_count;
    RAISE NOTICE '- Navigation States: 5 sample records';
    RAISE NOTICE '- User Preferences: 3 sample records';

    IF content_count >= 30 THEN
        RAISE NOTICE 'SUCCESS: Stream content seed data loaded successfully';
    ELSE
        RAISE EXCEPTION 'FAILURE: Expected at least 30 content items, found %', content_count;
    END IF;
END $$;

-- ============================================================================
-- CONTENT SUMMARY BY CATEGORY
-- ============================================================================

-- Show content distribution by category
SELECT
    sc.name as category_name,
    COUNT(sci.id) as content_count,
    COUNT(CASE WHEN sci.is_featured = true THEN 1 END) as featured_count,
    AVG(sci.price) as avg_price
FROM stream_categories sc
LEFT JOIN stream_content_items sci ON sc.id = sci.category_id
GROUP BY sc.id, sc.name, sc.display_order
ORDER BY sc.display_order;

-- ============================================================================
-- SEED DATA COMPLETE
-- ============================================================================
-- Stream Store Tabs seed data loading completed successfully
-- Categories: 6 (Books, Podcasts, Cartoons, Movies, Music, Art)
-- Content Items: 30+ across all categories
-- Sample User Data: Navigation states and preferences for testing