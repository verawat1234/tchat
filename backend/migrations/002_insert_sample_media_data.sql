-- Sample Data: Media Store Tabs
-- Date: 2025-09-29
-- Feature: Media Store Tabs (026-help-me-add)

-- Insert Media Categories
INSERT INTO media_categories (id, name, display_order, icon_name, is_active, featured_content_enabled) VALUES
('books', 'Books', 1, 'book-open', true, true),
('podcasts', 'Podcasts', 2, 'microphone', true, true),
('cartoons', 'Cartoons', 3, 'film', true, true),
('movies', 'Movies', 4, 'video', true, true)
ON CONFLICT (id) DO NOTHING;

-- Insert Media Subtabs (for Movies)
INSERT INTO media_subtabs (id, category_id, name, display_order, filter_criteria, is_active) VALUES
('short-movies', 'movies', 'Short Films', 1, '{"maxDuration": 1800}', true),
('long-movies', 'movies', 'Feature Films', 2, '{"minDuration": 1801}', true)
ON CONFLICT (id) DO NOTHING;

-- Insert Sample Media Content Items

-- Books
INSERT INTO media_content_items (category_id, title, description, thumbnail_url, content_type, duration, price, currency, availability_status, is_featured, featured_order, metadata) VALUES
('books', 'The Art of Software Architecture', 'A comprehensive guide to designing scalable software systems', 'https://cdn.tchat.com/media/thumbnails/book-arch.jpg', 'book', NULL, 19.99, 'USD', 'available', true, 1, '{"author": "John Smith", "genre": "Technology", "rating": 4.5, "pages": 320}'),
('books', 'Modern JavaScript Patterns', 'Master advanced JavaScript programming techniques', 'https://cdn.tchat.com/media/thumbnails/book-js.jpg', 'book', NULL, 24.99, 'USD', 'available', true, 2, '{"author": "Sarah Johnson", "genre": "Programming", "rating": 4.7, "pages": 280}'),
('books', 'Design Thinking Guide', 'A practical approach to user-centered design', 'https://cdn.tchat.com/media/thumbnails/book-design.jpg', 'book', NULL, 16.99, 'USD', 'available', false, NULL, '{"author": "Maria Garcia", "genre": "Design", "rating": 4.3, "pages": 240}'),
('books', 'Database Optimization', 'Performance tuning for modern databases', 'https://cdn.tchat.com/media/thumbnails/book-db.jpg', 'book', NULL, 21.99, 'USD', 'available', false, NULL, '{"author": "David Chen", "genre": "Technology", "rating": 4.6, "pages": 350}'),
('books', 'Mobile App Security', 'Securing iOS and Android applications', 'https://cdn.tchat.com/media/thumbnails/book-security.jpg', 'book', NULL, 18.99, 'USD', 'available', false, NULL, '{"author": "Alex Rodriguez", "genre": "Security", "rating": 4.4, "pages": 290}')
ON CONFLICT DO NOTHING;

-- Podcasts
INSERT INTO media_content_items (category_id, title, description, thumbnail_url, content_type, duration, price, currency, availability_status, is_featured, featured_order, metadata) VALUES
('podcasts', 'Tech Talk Weekly', 'Latest trends in software development', 'https://cdn.tchat.com/media/thumbnails/podcast-tech.jpg', 'podcast', 2400, 2.99, 'USD', 'available', true, 1, '{"host": "Mike Wilson", "category": "Technology", "rating": 4.8, "episode": 142}'),
('podcasts', 'Startup Stories', 'Interviews with successful entrepreneurs', 'https://cdn.tchat.com/media/thumbnails/podcast-startup.jpg', 'podcast', 3600, 1.99, 'USD', 'available', true, 2, '{"host": "Lisa Chang", "category": "Business", "rating": 4.6, "episode": 89}'),
('podcasts', 'Design Philosophy', 'Deep dives into design principles', 'https://cdn.tchat.com/media/thumbnails/podcast-design.jpg', 'podcast', 1800, 2.49, 'USD', 'available', false, NULL, '{"host": "Tom Anderson", "category": "Design", "rating": 4.5, "episode": 67}'),
('podcasts', 'Code Review', 'Programming best practices and reviews', 'https://cdn.tchat.com/media/thumbnails/podcast-code.jpg', 'podcast', 2700, 2.99, 'USD', 'available', false, NULL, '{"host": "Emma Davis", "category": "Programming", "rating": 4.7, "episode": 156}'),
('podcasts', 'AI Insights', 'Artificial intelligence and machine learning', 'https://cdn.tchat.com/media/thumbnails/podcast-ai.jpg', 'podcast', 3300, 3.99, 'USD', 'available', false, NULL, '{"host": "Dr. James Kim", "category": "AI/ML", "rating": 4.9, "episode": 78}')
ON CONFLICT DO NOTHING;

-- Cartoons
INSERT INTO media_content_items (category_id, title, description, thumbnail_url, content_type, duration, price, currency, availability_status, is_featured, featured_order, metadata) VALUES
('cartoons', 'Tech Adventures', 'Educational cartoon about programming', 'https://cdn.tchat.com/media/thumbnails/cartoon-tech.jpg', 'cartoon', 1200, 4.99, 'USD', 'available', true, 1, '{"studio": "EduToons", "genre": "Educational", "rating": 4.4, "seasons": 2}'),
('cartoons', 'Space Explorers', 'Animated series about space exploration', 'https://cdn.tchat.com/media/thumbnails/cartoon-space.jpg', 'cartoon', 1800, 5.99, 'USD', 'available', true, 2, '{"studio": "Galaxy Studios", "genre": "Adventure", "rating": 4.7, "seasons": 3}'),
('cartoons', 'Ocean Tales', 'Underwater adventures for all ages', 'https://cdn.tchat.com/media/thumbnails/cartoon-ocean.jpg', 'cartoon', 1500, 3.99, 'USD', 'available', false, NULL, '{"studio": "Blue Wave", "genre": "Family", "rating": 4.2, "seasons": 1}'),
('cartoons', 'Robot Friends', 'Friendship stories with robots', 'https://cdn.tchat.com/media/thumbnails/cartoon-robots.jpg', 'cartoon', 1350, 4.49, 'USD', 'available', false, NULL, '{"studio": "Future Toons", "genre": "Sci-Fi", "rating": 4.5, "seasons": 2}'),
('cartoons', 'Nature Quest', 'Environmental awareness through animation', 'https://cdn.tchat.com/media/thumbnails/cartoon-nature.jpg', 'cartoon', 1650, 3.49, 'USD', 'available', false, NULL, '{"studio": "Green Studios", "genre": "Educational", "rating": 4.3, "seasons": 1}')
ON CONFLICT DO NOTHING;

-- Short Movies (≤30 minutes)
INSERT INTO media_content_items (category_id, title, description, thumbnail_url, content_type, duration, price, currency, availability_status, is_featured, featured_order, metadata) VALUES
('movies', 'Digital Dreams', 'A short film about virtual reality', 'https://cdn.tchat.com/media/thumbnails/short-digital.jpg', 'short_movie', 1500, 7.99, 'USD', 'available', true, 1, '{"director": "Anna Kim", "genre": "Sci-Fi", "rating": 4.6, "year": 2024}'),
('movies', 'Coffee Break', 'Heartwarming story set in a cafe', 'https://cdn.tchat.com/media/thumbnails/short-coffee.jpg', 'short_movie', 1200, 5.99, 'USD', 'available', true, 2, '{"director": "Carlos Martinez", "genre": "Drama", "rating": 4.3, "year": 2024}'),
('movies', 'Code Warrior', 'Cyberpunk action in 20 minutes', 'https://cdn.tchat.com/media/thumbnails/short-code.jpg', 'short_movie', 1800, 6.49, 'USD', 'available', false, NULL, '{"director": "Yuki Tanaka", "genre": "Action", "rating": 4.4, "year": 2024}')
ON CONFLICT DO NOTHING;

-- Long Movies (>30 minutes)
INSERT INTO media_content_items (category_id, title, description, thumbnail_url, content_type, duration, price, currency, availability_status, is_featured, featured_order, metadata) VALUES
('movies', 'The Algorithm', 'Feature film about AI consciousness', 'https://cdn.tchat.com/media/thumbnails/long-algorithm.jpg', 'long_movie', 7200, 12.99, 'USD', 'available', true, 3, '{"director": "Robert Chen", "genre": "Sci-Fi", "rating": 4.8, "year": 2024}'),
('movies', 'Startup Chronicles', 'Documentary about tech entrepreneurship', 'https://cdn.tchat.com/media/thumbnails/long-startup.jpg', 'long_movie', 5400, 9.99, 'USD', 'available', true, 4, '{"director": "Sophie Anderson", "genre": "Documentary", "rating": 4.5, "year": 2024}'),
('movies', 'Digital Nomad', 'Comedy about remote work adventures', 'https://cdn.tchat.com/media/thumbnails/long-nomad.jpg', 'long_movie', 6300, 11.49, 'USD', 'available', false, NULL, '{"director": "Ahmed Hassan", "genre": "Comedy", "rating": 4.2, "year": 2024}')
ON CONFLICT DO NOTHING;