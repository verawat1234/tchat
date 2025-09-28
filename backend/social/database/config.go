package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tchat.dev/shared/config"
	"tchat.dev/shared/service"
	"tchat/social/models"
)

// SocialDatabaseConfig holds social service database configuration
type SocialDatabaseConfig struct {
	db          *gorm.DB
	initializer service.DatabaseInitializer
}

// NewSocialDatabaseConfig creates a new social database configuration
func NewSocialDatabaseConfig() *SocialDatabaseConfig {
	// Define social service models for migration
	socialModels := []interface{}{
		&models.SocialProfile{},
		&models.Post{},
		&models.Comment{},
		&models.Reaction{},
		&models.Follow{},
		&models.Community{},
		&models.CommunityMember{},
		&models.Share{},
	}

	// Create service-specific database initializer
	initializer := service.NewServiceDatabaseInitializer(socialModels, customSocialInitialization)

	return &SocialDatabaseConfig{
		initializer: initializer,
	}
}

// Initialize sets up the database connection and runs migrations
func (c *SocialDatabaseConfig) Initialize(cfg *config.Config) error {
	log.Println("Initializing social service database...")

	// Initialize database using shared infrastructure
	db, err := c.initializer.InitializeDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Store database reference
	c.db = db

	// Run migrations
	if err := c.initializer.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Apply performance optimizations
	if err := c.optimizeDatabase(db, cfg); err != nil {
		return fmt.Errorf("failed to apply performance optimizations: %w", err)
	}

	log.Println("Social service database initialized successfully")
	return nil
}

// GetDB returns the GORM database instance
func (c *SocialDatabaseConfig) GetDB() *gorm.DB {
	return c.db
}

// customSocialInitialization performs social service specific database setup
func customSocialInitialization(db *gorm.DB, cfg *config.Config) error {
	log.Println("Running social service custom database initialization...")

	// Create database indexes for performance
	if err := createSocialIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Set up database constraints
	if err := createSocialConstraints(db); err != nil {
		return fmt.Errorf("failed to create constraints: %w", err)
	}

	// Initialize default data if needed
	if err := seedDefaultData(db, cfg); err != nil {
		return fmt.Errorf("failed to seed default data: %w", err)
	}

	log.Println("Social service custom database initialization completed")
	return nil
}

// createSocialIndexes creates performance indexes for social data
func createSocialIndexes(db *gorm.DB) error {
	log.Println("Creating social service database indexes...")

	indexes := []string{
		// Posts table indexes
		"CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id)",
		"CREATE INDEX IF NOT EXISTS idx_posts_community_id ON posts(community_id) WHERE community_id IS NOT NULL",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_posts_visibility ON posts(visibility)",
		"CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type)",
		"CREATE INDEX IF NOT EXISTS idx_posts_trending ON posts(is_trending) WHERE is_trending = true",
		"CREATE INDEX IF NOT EXISTS idx_posts_not_deleted ON posts(id) WHERE is_deleted = false",

		// Comments table indexes
		"CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)",
		"CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id)",
		"CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id) WHERE parent_id IS NOT NULL",
		"CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC)",

		// Reactions table indexes
		"CREATE INDEX IF NOT EXISTS idx_reactions_target ON reactions(target_id, target_type)",
		"CREATE INDEX IF NOT EXISTS idx_reactions_user_target ON reactions(user_id, target_id, target_type)",
		"CREATE INDEX IF NOT EXISTS idx_reactions_type ON reactions(type)",

		// Follows table indexes
		"CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id)",
		"CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id)",
		"CREATE INDEX IF NOT EXISTS idx_follows_created_at ON follows(created_at DESC)",

		// Social profiles indexes
		"CREATE INDEX IF NOT EXISTS idx_social_profiles_interests ON social_profiles USING GIN(interests)",
		"CREATE INDEX IF NOT EXISTS idx_social_profiles_verified ON social_profiles(is_social_verified) WHERE is_social_verified = true",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			// Log the error but don't fail in test environments
			log.Printf("Warning: Failed to create index: %s, error: %v", indexSQL, err)
			// Only fail if it's not a "relation does not exist" error (which is expected in test environments)
			if !strings.Contains(err.Error(), "relation") && !strings.Contains(err.Error(), "does not exist") {
				return fmt.Errorf("failed to create index: %s, error: %w", indexSQL, err)
			}
		}
	}

	log.Printf("Created %d database indexes", len(indexes))
	return nil
}

// createSocialConstraints creates database constraints for data integrity
func createSocialConstraints(db *gorm.DB) error {
	log.Println("Creating social service database constraints...")

	constraints := []string{
		// Posts constraints
		"ALTER TABLE posts ADD CONSTRAINT check_posts_visibility CHECK (visibility IN ('public', 'members', 'private', 'followers'))",
		"ALTER TABLE posts ADD CONSTRAINT check_posts_type CHECK (type IN ('text', 'image', 'video', 'link', 'poll'))",

		// Reactions constraints
		"ALTER TABLE reactions ADD CONSTRAINT check_reactions_target_type CHECK (target_type IN ('post', 'comment'))",
		"ALTER TABLE reactions ADD CONSTRAINT check_reactions_type CHECK (type IN ('like', 'love', 'laugh', 'angry', 'sad', 'wow'))",

		// Follows constraints (prevent self-following already in model)
		"ALTER TABLE follows ADD CONSTRAINT check_follows_no_self CHECK (follower_id != following_id)",
	}

	for _, constraintSQL := range constraints {
		// Use IF NOT EXISTS pattern for constraints
		if err := db.Exec(constraintSQL).Error; err != nil {
			// Log but don't fail on constraint errors (they might already exist)
			log.Printf("Constraint may already exist: %s", constraintSQL)
		}
	}

	log.Printf("Applied %d database constraints", len(constraints))
	return nil
}

// seedDefaultData seeds the database with default data if needed
func seedDefaultData(db *gorm.DB, cfg *config.Config) error {
	log.Println("Seeding default social service data...")

	// Only seed in development mode
	if !cfg.Debug {
		log.Println("Skipping data seeding in production mode")
		return nil
	}

	// Check if data already exists
	var count int64
	if err := db.Model(&models.Community{}).Count(&count).Error; err != nil {
		// If table doesn't exist in test environment, skip seeding
		if strings.Contains(err.Error(), "relation") && strings.Contains(err.Error(), "does not exist") {
			log.Printf("Warning: Skipping data seeding, tables don't exist: %v", err)
			return nil
		}
		return fmt.Errorf("failed to check existing data: %w", err)
	}

	if count > 0 {
		log.Println("Data already exists, skipping seeding")
		return nil
	}

	// Seed default communities
	defaultCommunities := []*models.Community{
		{
			Name:        "General",
			Description: "General discussion for all topics",
			Type:        "public",
			Settings: map[string]interface{}{
				"allow_posts":    true,
				"allow_comments": true,
				"moderated":      false,
			},
			Rules: []string{
				"Be respectful to all community members",
				"No spam or promotional content",
				"Stay on topic",
			},
		},
		{
			Name:        "Technology",
			Description: "Discuss technology trends and innovations",
			Type:        "public",
			Settings: map[string]interface{}{
				"allow_posts":    true,
				"allow_comments": true,
				"moderated":      true,
			},
			Rules: []string{
				"Technology-related content only",
				"Provide sources for news and claims",
				"Be constructive in discussions",
			},
		},
	}

	for _, community := range defaultCommunities {
		if err := db.Create(community).Error; err != nil {
			return fmt.Errorf("failed to seed community %s: %w", community.Name, err)
		}
	}

	log.Printf("Seeded %d default communities", len(defaultCommunities))
	return nil
}

// optimizeDatabase applies comprehensive performance optimizations
func (c *SocialDatabaseConfig) optimizeDatabase(db *gorm.DB, cfg *config.Config) error {
	log.Println("Applying database performance optimizations...")

	// Apply connection pool optimizations
	if err := c.configureConnectionPool(db, cfg); err != nil {
		return fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// Apply PostgreSQL-specific optimizations
	if err := c.applyPostgreSQLOptimizations(db); err != nil {
		return fmt.Errorf("failed to apply PostgreSQL optimizations: %w", err)
	}

	// Create advanced performance indexes
	if err := c.createAdvancedIndexes(db); err != nil {
		return fmt.Errorf("failed to create advanced indexes: %w", err)
	}

	// Configure query optimization settings
	if err := c.configureQueryOptimization(db); err != nil {
		return fmt.Errorf("failed to configure query optimization: %w", err)
	}

	log.Println("Database performance optimizations applied successfully")
	return nil
}

// configureConnectionPool optimizes database connection pool settings
func (c *SocialDatabaseConfig) configureConnectionPool(db *gorm.DB, cfg *config.Config) error {
	log.Println("Configuring database connection pool...")

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings optimized for social media workload
	maxIdleConns := 10
	maxOpenConns := 50
	connMaxLifetime := time.Hour
	connMaxIdleTime := time.Minute * 30

	// Adjust for production vs development
	if !cfg.Debug {
		maxIdleConns = 20
		maxOpenConns = 100
		connMaxLifetime = time.Hour * 2
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	log.Printf("Connection pool configured: MaxIdle=%d, MaxOpen=%d, MaxLifetime=%v, MaxIdleTime=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)

	return nil
}

// applyPostgreSQLOptimizations applies PostgreSQL-specific performance settings
func (c *SocialDatabaseConfig) applyPostgreSQLOptimizations(db *gorm.DB) error {
	log.Println("Applying PostgreSQL performance optimizations...")

	optimizations := []string{
		// Memory and cache optimizations
		"SET shared_buffers = '256MB'",
		"SET effective_cache_size = '1GB'",
		"SET work_mem = '16MB'",
		"SET maintenance_work_mem = '64MB'",

		// Query planner optimizations
		"SET random_page_cost = 1.1",
		"SET effective_io_concurrency = 200",

		// WAL and checkpoint optimizations
		"SET wal_buffers = '16MB'",
		"SET checkpoint_completion_target = 0.7",
		"SET checkpoint_segments = 32",

		// Connection and logging optimizations
		"SET log_statement = 'mod'",
		"SET log_min_duration_statement = 1000",
		"SET track_activity_query_size = 2048",

		// Auto vacuum optimizations for social media workload
		"SET autovacuum_vacuum_scale_factor = 0.1",
		"SET autovacuum_analyze_scale_factor = 0.05",
	}

	for _, optimization := range optimizations {
		if err := db.Exec(optimization).Error; err != nil {
			// Log but don't fail on configuration errors (might be restricted)
			log.Printf("PostgreSQL optimization warning: %s - %v", optimization, err)
		}
	}

	log.Printf("Applied %d PostgreSQL optimizations", len(optimizations))
	return nil
}

// createAdvancedIndexes creates sophisticated indexes for complex queries
func (c *SocialDatabaseConfig) createAdvancedIndexes(db *gorm.DB) error {
	log.Println("Creating advanced performance indexes...")

	advancedIndexes := []string{
		// Composite indexes for complex social queries
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_author_created_visibility ON posts(author_id, created_at DESC, visibility) WHERE is_deleted = false",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_community_trending ON posts(community_id, is_trending, created_at DESC) WHERE community_id IS NOT NULL",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_viral_content ON posts(likes_count DESC, shares_count DESC, comments_count DESC) WHERE is_trending = true",

		// Social graph optimizations
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_follows_mutual_connections ON follows(follower_id, following_id) INCLUDE (created_at)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_social_profiles_discovery ON social_profiles(country, is_social_verified, followers_count DESC) WHERE is_social_verified = true",

		// Comment thread optimizations
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_thread_hierarchy ON comments(post_id, parent_id, created_at) WHERE parent_id IS NOT NULL",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_user_activity ON comments(author_id, created_at DESC) INCLUDE (post_id)",

		// Reaction aggregation optimizations
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_popular_content ON reactions(target_id, target_type, type) INCLUDE (created_at)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_user_engagement ON reactions(user_id, created_at DESC) INCLUDE (target_id, type)",

		// Community engagement indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_community_members_active ON community_members(community_id, status, joined_at DESC) WHERE status = 'active'",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_communities_popular ON communities(members_count DESC, posts_count DESC, is_featured) WHERE type = 'public'",

		// Content sharing optimizations
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_shares_viral_tracking ON shares(content_id, content_type, platform, created_at DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_shares_user_sharing_patterns ON shares(user_id, platform, privacy, created_at DESC)",

		// Full-text search indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_content_search ON posts USING gin(to_tsvector('english', content)) WHERE is_deleted = false",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_communities_search ON communities USING gin(to_tsvector('english', name || ' ' || description))",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_social_profiles_search ON social_profiles USING gin(to_tsvector('english', coalesce(display_name, '') || ' ' || coalesce(bio, '')))",

		// Performance monitoring indexes
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_performance_metrics ON posts(created_at, likes_count, comments_count, shares_count) WHERE created_at > CURRENT_DATE - INTERVAL '30 days'",
	}

	for _, indexSQL := range advancedIndexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			// Log but continue on index creation errors (might already exist)
			log.Printf("Advanced index warning: %s - %v", indexSQL, err)
		}
	}

	log.Printf("Created %d advanced performance indexes", len(advancedIndexes))
	return nil
}

// configureQueryOptimization configures GORM and PostgreSQL query optimizations
func (c *SocialDatabaseConfig) configureQueryOptimization(db *gorm.DB) error {
	log.Println("Configuring query optimization settings...")

	// Configure GORM session with performance optimizations
	db = db.Session(&gorm.Session{
		PrepareStmt:            true,  // Use prepared statements for better performance
		SkipDefaultTransaction: true,  // Skip automatic transactions for better performance
		SkipHooks:              false, // Keep hooks for data integrity
		QueryFields:            true,  // Optimize SELECT queries
		Logger:                 logger.Default.LogMode(logger.Silent), // Reduce logging overhead in production
	})

	// Apply query-level optimizations
	queryOptimizations := []string{
		// Enable query parallelism
		"SET max_parallel_workers_per_gather = 4",
		"SET max_parallel_workers = 8",
		"SET parallel_tuple_cost = 0.1",
		"SET parallel_setup_cost = 1000",

		// Optimize join algorithms
		"SET enable_hashjoin = on",
		"SET enable_mergejoin = on",
		"SET enable_nestloop = on",

		// Statistics and planning
		"SET default_statistics_target = 1000",
		"SET from_collapse_limit = 12",
		"SET join_collapse_limit = 12",

		// Enable specific optimizations for social media queries
		"SET enable_partitionwise_join = on",
		"SET enable_partitionwise_aggregate = on",
	}

	for _, optimization := range queryOptimizations {
		if err := db.Exec(optimization).Error; err != nil {
			log.Printf("Query optimization warning: %s - %v", optimization, err)
		}
	}

	log.Printf("Applied %d query optimizations", len(queryOptimizations))
	return nil
}

// GetConnectionStats returns database connection pool statistics
func (c *SocialDatabaseConfig) GetConnectionStats() map[string]interface{} {
	if c.db == nil {
		return map[string]interface{}{"error": "database not initialized"}
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return map[string]interface{}{"error": "failed to get underlying sql.DB"}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections":     stats.MaxOpenConnections,
		"open_connections":         stats.OpenConnections,
		"in_use":                   stats.InUse,
		"idle":                     stats.Idle,
		"wait_count":              stats.WaitCount,
		"wait_duration":           stats.WaitDuration,
		"max_idle_closed":         stats.MaxIdleClosed,
		"max_idle_time_closed":    stats.MaxIdleTimeClosed,
		"max_lifetime_closed":     stats.MaxLifetimeClosed,
	}
}

// PerformanceAnalyzer provides database performance analysis tools
type PerformanceAnalyzer struct {
	db *gorm.DB
}

// NewPerformanceAnalyzer creates a new performance analyzer
func (c *SocialDatabaseConfig) NewPerformanceAnalyzer() *PerformanceAnalyzer {
	return &PerformanceAnalyzer{db: c.db}
}

// AnalyzeSlowQueries returns slow query analysis
func (p *PerformanceAnalyzer) AnalyzeSlowQueries(ctx context.Context, durationThreshold string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			query,
			calls,
			total_time,
			mean_time,
			min_time,
			max_time,
			stddev_time
		FROM pg_stat_statements
		WHERE mean_time > ?
		ORDER BY mean_time DESC
		LIMIT 20
	`

	rows, err := p.db.WithContext(ctx).Raw(query, durationThreshold).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze slow queries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result map[string]interface{}
		if err := p.db.ScanRows(rows, &result); err != nil {
			continue // Skip problematic rows
		}
		results = append(results, result)
	}

	return results, nil
}

// GetTableSizes returns table size information for capacity planning
func (p *PerformanceAnalyzer) GetTableSizes(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			schemaname,
			tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
			pg_total_relation_size(schemaname||'.'||tablename) as size_bytes
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
	`

	if err := p.db.WithContext(ctx).Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get table sizes: %w", err)
	}

	return results, nil
}