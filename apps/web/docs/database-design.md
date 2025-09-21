# Database Design for Telegram SEA Edition

This document outlines the comprehensive database schema design for the Telegram SEA Edition application, based on UI component analysis and business requirements.

## Overview

The database is designed to support a multi-tenant, multi-currency, and multi-language social commerce platform optimized for Southeast Asian markets. The schema supports:

- **10 major domains**: User, Chat, Social, Commerce, Events, Wallet, Workspace, Video, Discovery, Activities
- **Multi-currency**: THB, SGD, IDR, MYR, PHP, VND, USD
- **Multi-language**: th-TH, id-ID, ms-MY, vi-VN, en-US
- **High scalability**: Designed for millions of users and billions of messages
- **Real-time features**: Optimized for live messaging, payments, and social interactions

## Technology Stack

- **Primary Database**: PostgreSQL 15+ with extensions
- **Extensions**: PostGIS (geolocation), pg_trgm (full-text search), uuid-ossp (UUIDs)
- **Search Engine**: Elasticsearch for advanced search capabilities
- **Cache Layer**: Redis for session management and real-time features
- **File Storage**: S3-compatible storage for media files
- **Analytics**: Separate OLAP database (ClickHouse) for analytics

## Core Design Principles

### 1. Multi-tenancy
- Workspace-based isolation for business features
- User data segregation with privacy controls
- Regional data compliance (GDPR, PDPA)

### 2. Scalability
- Horizontal partitioning for large tables (messages, transactions)
- Read replicas for improved performance
- Caching strategies for frequently accessed data

### 3. Data Integrity
- Foreign key constraints with cascade rules
- Check constraints for business rules
- ACID compliance for financial transactions

### 4. Performance
- Strategic indexing for common query patterns
- Materialized views for analytics
- Optimized JSON columns for flexible data

## Schema Organization

### Database Structure
```
telegram_sea/
├── public/                 # Default schema for shared data
├── user_data/             # User-related tables
├── messaging/             # Chat and messaging tables
├── commerce/              # E-commerce tables
├── social/                # Social features tables
├── events/                # Events and ticketing tables
├── financial/             # Wallet and payment tables
├── workspace/             # Business workspace tables
├── media/                 # Video and streaming tables
├── discovery/             # Search and discovery tables
├── activities/            # Gamification tables
└── analytics/             # Analytics and reporting tables
```

---

## Domain Schemas

## 1. User Domain (user_data schema)

### Core Tables

#### `users`
```sql
CREATE TABLE user_data.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(255) UNIQUE,
    name VARCHAR(255) NOT NULL,
    avatar TEXT,
    country CHAR(2) NOT NULL CHECK (country IN ('TH', 'ID', 'MY', 'VN', 'SG', 'PH')),
    locale VARCHAR(5) NOT NULL DEFAULT 'en-US',
    kyc_tier INTEGER NOT NULL DEFAULT 1 CHECK (kyc_tier IN (1, 2, 3)),
    status VARCHAR(20) NOT NULL DEFAULT 'offline'
        CHECK (status IN ('online', 'offline', 'away', 'busy')),
    last_seen TIMESTAMPTZ,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    settings JSONB NOT NULL DEFAULT '{}',
    profile JSONB NOT NULL DEFAULT '{}',
    preferences JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_phone ON user_data.users (phone);
CREATE INDEX idx_users_email ON user_data.users (email);
CREATE INDEX idx_users_country ON user_data.users (country);
CREATE INDEX idx_users_status ON user_data.users (status);
CREATE INDEX idx_users_last_seen ON user_data.users (last_seen);
CREATE INDEX idx_users_created_at ON user_data.users (created_at);

-- Full-text search
CREATE INDEX idx_users_name_fts ON user_data.users USING GIN (to_tsvector('english', name));
```

#### `user_sessions`
```sql
CREATE TABLE user_data.user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    device_info JSONB NOT NULL DEFAULT '{}',
    access_token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ip_address INET,
    location JSONB,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_user_sessions_user_id ON user_data.user_sessions (user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_data.user_sessions (expires_at);
CREATE INDEX idx_user_sessions_is_active ON user_data.user_sessions (is_active);
```

#### `kyc_documents`
```sql
CREATE TABLE user_data.kyc_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('id_card', 'passport', 'driving_license', 'utility_bill', 'bank_statement')),
    file_url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'approved', 'rejected')),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMPTZ,
    rejection_reason TEXT,
    metadata JSONB DEFAULT '{}'
);

-- Indexes
CREATE INDEX idx_kyc_documents_user_id ON user_data.kyc_documents (user_id);
CREATE INDEX idx_kyc_documents_status ON user_data.kyc_documents (status);
```

#### `friends`
```sql
CREATE TABLE user_data.friends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    friend_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted', 'blocked')),
    mutual_friends INTEGER DEFAULT 0,
    common_interests TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    accepted_at TIMESTAMPTZ,

    UNIQUE(user_id, friend_id),
    CHECK (user_id != friend_id)
);

-- Indexes
CREATE INDEX idx_friends_user_id ON user_data.friends (user_id);
CREATE INDEX idx_friends_friend_id ON user_data.friends (friend_id);
CREATE INDEX idx_friends_status ON user_data.friends (status);
```

---

## 2. Messaging Domain (messaging schema)

### Core Tables

#### `dialogs`
```sql
CREATE TABLE messaging.dialogs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('user', 'group', 'channel', 'bot', 'business')),
    name VARCHAR(255),
    description TEXT,
    avatar TEXT,
    participants UUID[] NOT NULL DEFAULT '{}',
    admins UUID[] NOT NULL DEFAULT '{}',
    owners UUID[] NOT NULL DEFAULT '{}',
    last_message_id UUID,
    unread_count INTEGER NOT NULL DEFAULT 0,
    muted_until TIMESTAMPTZ,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    permissions JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_dialogs_type ON messaging.dialogs (type);
CREATE INDEX idx_dialogs_participants ON messaging.dialogs USING GIN (participants);
CREATE INDEX idx_dialogs_updated_at ON messaging.dialogs (updated_at);
```

#### `messages` (Partitioned by month)
```sql
CREATE TABLE messaging.messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dialog_id UUID NOT NULL REFERENCES messaging.dialogs(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'voice', 'file', 'image', 'video', 'payment', 'system', 'location', 'contact', 'poll', 'event', 'product', 'sticker', 'gif')),
    content JSONB NOT NULL DEFAULT '{}',
    reply_to_id UUID REFERENCES messaging.messages(id),
    forward_from_id UUID REFERENCES messaging.messages(id),
    thread_id UUID,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    mentions UUID[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE messaging.messages_2024_01 PARTITION OF messaging.messages
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Indexes (on parent table)
CREATE INDEX idx_messages_dialog_id ON messaging.messages (dialog_id);
CREATE INDEX idx_messages_sender_id ON messaging.messages (sender_id);
CREATE INDEX idx_messages_type ON messaging.messages (type);
CREATE INDEX idx_messages_created_at ON messaging.messages (created_at);
CREATE INDEX idx_messages_thread_id ON messaging.messages (thread_id);

-- Full-text search on message content
CREATE INDEX idx_messages_content_fts ON messaging.messages USING GIN (to_tsvector('english', content::text));
```

#### `message_reactions`
```sql
CREATE TABLE messaging.message_reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    emoji VARCHAR(20) NOT NULL,
    skin_tone VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(message_id, user_id, emoji)
);

-- Indexes
CREATE INDEX idx_message_reactions_message_id ON messaging.message_reactions (message_id);
CREATE INDEX idx_message_reactions_user_id ON messaging.message_reactions (user_id);
```

#### `message_read_status`
```sql
CREATE TABLE messaging.message_read_status (
    message_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    device_id VARCHAR(255),

    PRIMARY KEY (message_id, user_id)
);

-- Indexes
CREATE INDEX idx_message_read_status_user_id ON messaging.message_read_status (user_id);
CREATE INDEX idx_message_read_status_read_at ON messaging.message_read_status (read_at);
```

---

## 3. Commerce Domain (commerce schema)

### Core Tables

#### `shops`
```sql
CREATE TABLE commerce.shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    avatar TEXT,
    cover_image TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_level VARCHAR(20) NOT NULL DEFAULT 'none'
        CHECK (verification_level IN ('none', 'basic', 'premium', 'enterprise')),
    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'suspended', 'under_review', 'closed', 'maintenance')),
    settings JSONB NOT NULL DEFAULT '{}',
    contact JSONB NOT NULL DEFAULT '{}',
    location JSONB,
    stats JSONB NOT NULL DEFAULT '{}',
    policies JSONB NOT NULL DEFAULT '{}',
    categories TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    subscription JSONB NOT NULL DEFAULT '{}',
    compliance JSONB DEFAULT '{}',
    localization JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_shops_owner_id ON commerce.shops (owner_id);
CREATE INDEX idx_shops_status ON commerce.shops (status);
CREATE INDEX idx_shops_is_verified ON commerce.shops (is_verified);
CREATE INDEX idx_shops_categories ON commerce.shops USING GIN (categories);
CREATE INDEX idx_shops_name_fts ON commerce.shops USING GIN (to_tsvector('english', name));
```

#### `products`
```sql
CREATE TABLE commerce.products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shop_id UUID NOT NULL REFERENCES commerce.shops(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    short_description TEXT,
    images JSONB NOT NULL DEFAULT '[]',
    videos JSONB DEFAULT '[]',
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    compare_at_price DECIMAL(15,2) CHECK (compare_at_price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'THB'
        CHECK (currency IN ('THB', 'SGD', 'IDR', 'MYR', 'PHP', 'VND', 'USD')),
    cost DECIMAL(15,2) CHECK (cost >= 0),
    sku VARCHAR(100),
    barcode VARCHAR(100),
    inventory JSONB NOT NULL DEFAULT '{}',
    variants JSONB DEFAULT '[]',
    category_id UUID REFERENCES commerce.product_categories(id),
    tags TEXT[] DEFAULT '{}',
    attributes JSONB DEFAULT '{}',
    seo JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'active', 'out_of_stock', 'discontinued', 'archived')),
    is_digital BOOLEAN NOT NULL DEFAULT FALSE,
    weight DECIMAL(8,3) CHECK (weight >= 0),
    dimensions JSONB,
    shipping JSONB NOT NULL DEFAULT '{}',
    ratings JSONB NOT NULL DEFAULT '{"average_rating": 0, "total_reviews": 0}',
    localization JSONB DEFAULT '{}',
    compliance JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_products_shop_id ON commerce.products (shop_id);
CREATE INDEX idx_products_category_id ON commerce.products (category_id);
CREATE INDEX idx_products_status ON commerce.products (status);
CREATE INDEX idx_products_price ON commerce.products (price);
CREATE INDEX idx_products_currency ON commerce.products (currency);
CREATE INDEX idx_products_is_digital ON commerce.products (is_digital);
CREATE INDEX idx_products_sku ON commerce.products (sku);
CREATE INDEX idx_products_tags ON commerce.products USING GIN (tags);
CREATE INDEX idx_products_title_fts ON commerce.products USING GIN (to_tsvector('english', title));
```

#### `orders` (Partitioned by created_at)
```sql
CREATE TABLE commerce.orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID REFERENCES user_data.users(id) ON DELETE SET NULL,
    customer_info JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded', 'disputed')),
    items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(15,2) NOT NULL CHECK (subtotal >= 0),
    discount DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (discount >= 0),
    tax DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (tax >= 0),
    shipping DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (shipping >= 0),
    total DECIMAL(15,2) NOT NULL CHECK (total >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'THB',
    coupon_code VARCHAR(100),
    shipping_address JSONB NOT NULL,
    billing_address JSONB,
    payment JSONB NOT NULL DEFAULT '{}',
    fulfillment JSONB NOT NULL DEFAULT '{}',
    communications JSONB DEFAULT '[]',
    timeline JSONB DEFAULT '[]',
    refunds JSONB DEFAULT '[]',
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE commerce.orders_2024_01 PARTITION OF commerce.orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Indexes
CREATE INDEX idx_orders_user_id ON commerce.orders (user_id);
CREATE INDEX idx_orders_status ON commerce.orders (status);
CREATE INDEX idx_orders_currency ON commerce.orders (currency);
CREATE INDEX idx_orders_created_at ON commerce.orders (created_at);
CREATE INDEX idx_orders_order_number ON commerce.orders (order_number);
```

---

## 4. Financial Domain (financial schema)

### Core Tables

#### `wallets`
```sql
CREATE TABLE financial.wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'THB'
        CHECK (currency IN ('THB', 'SGD', 'IDR', 'MYR', 'PHP', 'VND', 'USD')),
    frozen_balance DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (frozen_balance >= 0),
    available_balance DECIMAL(15,2) GENERATED ALWAYS AS (balance - frozen_balance) STORED,
    daily_limit DECIMAL(15,2) NOT NULL DEFAULT 5000,
    monthly_limit DECIMAL(15,2) NOT NULL DEFAULT 20000,
    used_this_month DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'suspended', 'closed')),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, currency)
);

-- Indexes
CREATE INDEX idx_wallets_user_id ON financial.wallets (user_id);
CREATE INDEX idx_wallets_currency ON financial.wallets (currency);
CREATE INDEX idx_wallets_status ON financial.wallets (status);
```

#### `transactions` (Partitioned by created_at)
```sql
CREATE TABLE financial.transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES financial.wallets(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('send', 'receive', 'topup', 'withdraw', 'purchase', 'refund', 'fee', 'reward', 'cashback')),
    amount DECIMAL(15,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    fee DECIMAL(15,2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(15,2) GENERATED ALWAYS AS (amount - fee) STORED,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled', 'expired')),
    description TEXT NOT NULL,
    reference VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    counterpart JSONB,
    category VARCHAR(50),
    tags TEXT[] DEFAULT '{}',
    balance_before DECIMAL(15,2),
    balance_after DECIMAL(15,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE financial.transactions_2024_01 PARTITION OF financial.transactions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Indexes
CREATE INDEX idx_transactions_wallet_id ON financial.transactions (wallet_id);
CREATE INDEX idx_transactions_type ON financial.transactions (type);
CREATE INDEX idx_transactions_status ON financial.transactions (status);
CREATE INDEX idx_transactions_created_at ON financial.transactions (created_at);
CREATE INDEX idx_transactions_reference ON financial.transactions (reference);
```

---

## 5. Events Domain (events schema)

### Core Tables

#### `events`
```sql
CREATE TABLE events.events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    short_description TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('music', 'food', 'cultural', 'festival', 'temple', 'market', 'sports', 'technology', 'business', 'art', 'education')),
    type VARCHAR(50) NOT NULL CHECK (type IN ('concert', 'festival', 'conference', 'workshop', 'exhibition', 'competition', 'celebration', 'ceremony')),
    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'cancelled', 'postponed', 'completed')),
    organizer JSONB NOT NULL,
    venue JSONB NOT NULL,
    schedule JSONB NOT NULL,
    ticketing JSONB NOT NULL DEFAULT '{}',
    media JSONB NOT NULL DEFAULT '{}',
    lineup JSONB,
    amenities TEXT[] DEFAULT '{}',
    age_restriction VARCHAR(20),
    tags TEXT[] DEFAULT '{}',
    popularity JSONB NOT NULL DEFAULT '{}',
    social_proof JSONB DEFAULT '{}',
    weather JSONB,
    is_promoted BOOLEAN NOT NULL DEFAULT FALSE,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_events_category ON events.events (category);
CREATE INDEX idx_events_type ON events.events (type);
CREATE INDEX idx_events_status ON events.events (status);
CREATE INDEX idx_events_is_promoted ON events.events (is_promoted);
CREATE INDEX idx_events_tags ON events.events USING GIN (tags);
CREATE INDEX idx_events_title_fts ON events.events USING GIN (to_tsvector('english', title));

-- Geospatial index for venue location
CREATE INDEX idx_events_venue_location ON events.events USING GIST (((venue->>'coordinates')::point));
```

#### `tickets`
```sql
CREATE TABLE events.tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events.events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES user_data.users(id) ON DELETE CASCADE,
    ticket_type_id UUID NOT NULL,
    ticket_number VARCHAR(100) NOT NULL UNIQUE,
    qr_code TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'transferred', 'refunded', 'cancelled', 'expired', 'used')),
    price DECIMAL(10,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    purchased_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    transferred_to UUID REFERENCES user_data.users(id),
    transferred_at TIMESTAMPTZ,
    checked_in BOOLEAN NOT NULL DEFAULT FALSE,
    checked_in_at TIMESTAMPTZ,
    refunded BOOLEAN NOT NULL DEFAULT FALSE,
    refunded_at TIMESTAMPTZ,
    refund_amount DECIMAL(10,2)
);

-- Indexes
CREATE INDEX idx_tickets_event_id ON events.tickets (event_id);
CREATE INDEX idx_tickets_user_id ON events.tickets (user_id);
CREATE INDEX idx_tickets_status ON events.tickets (status);
CREATE INDEX idx_tickets_ticket_number ON events.tickets (ticket_number);
```

---

## Performance Optimizations

### Partitioning Strategy

1. **Time-based partitioning** for high-volume tables:
   - `messages` partitioned by month
   - `transactions` partitioned by month
   - `orders` partitioned by month

2. **Hash partitioning** for user-based data:
   - Consider hash partitioning `users` table for very large datasets

### Indexing Strategy

1. **Primary Keys**: UUID with B-tree indexes
2. **Foreign Keys**: Automatic indexes on all foreign key columns
3. **Query-specific indexes**: Based on common WHERE clauses
4. **Composite indexes**: For complex queries with multiple conditions
5. **Partial indexes**: For filtered queries (e.g., active records only)
6. **Full-text search**: GIN indexes using PostgreSQL's built-in FTS

### Example Composite Indexes
```sql
-- Messages by dialog and time range
CREATE INDEX idx_messages_dialog_time ON messaging.messages (dialog_id, created_at DESC);

-- Products by shop, status, and price
CREATE INDEX idx_products_shop_status_price ON commerce.products (shop_id, status, price);

-- Transactions by user and date
CREATE INDEX idx_transactions_user_date ON financial.transactions (
    (metadata->>'user_id'), created_at DESC
) WHERE status = 'completed';
```

### Materialized Views

```sql
-- User statistics
CREATE MATERIALIZED VIEW user_data.user_stats_mv AS
SELECT
    u.id as user_id,
    COUNT(DISTINCT m.id) as total_messages,
    COUNT(DISTINCT o.id) as total_orders,
    COALESCE(SUM(o.total), 0) as total_spent,
    u.created_at as joined_at
FROM user_data.users u
LEFT JOIN messaging.messages m ON m.sender_id = u.id
LEFT JOIN commerce.orders o ON o.user_id = u.id
GROUP BY u.id, u.created_at;

-- Refresh daily
CREATE UNIQUE INDEX ON user_data.user_stats_mv (user_id);
```

---

## Data Migration Strategy

### Phase 1: Core Infrastructure
1. Create schemas and base tables
2. Set up partitioning for high-volume tables
3. Create basic indexes and constraints

### Phase 2: User Data
1. Migrate user accounts and profiles
2. Set up authentication and session management
3. Import friend relationships

### Phase 3: Messaging System
1. Create dialog structures
2. Migrate historical messages (if any)
3. Set up real-time messaging infrastructure

### Phase 4: Commerce Platform
1. Set up shop and product catalogs
2. Migrate inventory and pricing data
3. Set up payment processing

### Phase 5: Additional Features
1. Events and ticketing system
2. Financial transactions and wallets
3. Analytics and reporting

### Migration Scripts Example

```sql
-- Sample migration for user data
BEGIN;

-- Create temporary staging table
CREATE TEMP TABLE user_staging AS
SELECT * FROM old_system.users;

-- Validate and clean data
UPDATE user_staging
SET country = 'TH'
WHERE country IS NULL AND phone LIKE '+66%';

-- Insert into new schema
INSERT INTO user_data.users (
    id, phone, email, name, country, locale, created_at
)
SELECT
    gen_random_uuid(),
    phone,
    email,
    name,
    COALESCE(country, 'TH'),
    COALESCE(locale, 'th-TH'),
    COALESCE(created_at, NOW())
FROM user_staging
WHERE email IS NOT NULL OR phone IS NOT NULL;

COMMIT;
```

---

## Backup and Recovery

### Backup Strategy
1. **Daily full backups** of entire database
2. **Continuous WAL archiving** for point-in-time recovery
3. **Weekly logical backups** for disaster recovery
4. **Cross-region replication** for geographical redundancy

### Recovery Procedures
1. **Point-in-time recovery** for data corruption
2. **Standby promotion** for primary database failure
3. **Selective table restoration** for partial data loss
4. **Cross-region failover** for regional disasters

---

## Monitoring and Maintenance

### Performance Monitoring
1. **Query performance**: Track slow queries and optimization opportunities
2. **Index usage**: Monitor index hit ratios and unused indexes
3. **Partition pruning**: Ensure partitioning is working effectively
4. **Connection pooling**: Monitor connection pool efficiency

### Maintenance Tasks
1. **VACUUM and ANALYZE**: Regular maintenance for optimal performance
2. **Partition maintenance**: Automatic creation and dropping of partitions
3. **Statistics updates**: Ensure query planner has accurate statistics
4. **Index maintenance**: Regular REINDEX for heavily updated indexes

### Sample Monitoring Queries
```sql
-- Slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Index usage
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE schemaname NOT IN ('information_schema', 'pg_catalog')
ORDER BY n_distinct DESC;

-- Table sizes
SELECT schemaname, tablename,
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## Security Considerations

### Data Encryption
1. **Encryption at rest**: Full database encryption
2. **Encryption in transit**: TLS for all connections
3. **Column-level encryption**: For sensitive data (payment info, documents)
4. **Key management**: Separate key management service

### Access Control
1. **Role-based access**: Different roles for different access levels
2. **Row-level security**: User can only access their own data
3. **API rate limiting**: Prevent abuse through rate limiting
4. **Audit logging**: Track all data access and modifications

### Example Security Policies
```sql
-- Row Level Security for user data
ALTER TABLE user_data.users ENABLE ROW LEVEL SECURITY;

CREATE POLICY user_data_policy ON user_data.users
    FOR ALL TO application_role
    USING (id = current_setting('app.current_user')::uuid);

-- Grant appropriate permissions
GRANT SELECT, INSERT, UPDATE ON user_data.users TO application_role;
GRANT USAGE ON SCHEMA user_data TO application_role;
```

---

This database design provides a solid foundation for the Telegram SEA Edition application, supporting high scalability, performance, and the specific requirements of Southeast Asian markets.