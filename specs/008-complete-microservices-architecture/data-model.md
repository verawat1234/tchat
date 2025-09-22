# Data Model: Tchat Backend Platform

**Date**: 2025-09-22
**Context**: Entity relationships and data structures for Southeast Asian messaging and commerce platform

## Core Entities

### User Entity
**Purpose**: Individual platform users for communication and commerce
**Storage**: PostgreSQL (primary), Redis (sessions)

**Attributes**:
- `id`: UUID (primary key)
- `phone_number`: String (E.164 format, unique)
- `country_code`: String (ISO 3166-1 alpha-2: TH, SG, ID, MY, PH, VN)
- `status`: Enum (active, suspended, deleted)
- `profile`: JSON object
  - `display_name`: String (max 100 chars)
  - `avatar_url`: String (optional)
  - `locale`: String (en, th, id, ms, fil, vi)
  - `timezone`: String (IANA timezone)
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Validation Rules**:
- Phone number must be valid for specified country
- Display name required for active users
- Locale must match supported languages
- Country code determines available payment methods

**State Transitions**:
- active → suspended (admin action)
- suspended → active (appeal approved)
- active/suspended → deleted (user request, compliance)

### Business Entity
**Purpose**: Commercial entities operating shops on the platform
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `owner_id`: UUID (foreign key to User)
- `name`: String (max 100 chars)
- `description`: Text (max 1000 chars)
- `category`: String (electronics, fashion, food, etc.)
- `verification_status`: Enum (pending, verified, rejected)
- `contact_info`: JSON object
  - `phone`: String
  - `email`: String
  - `website`: String (optional)
- `address`: JSON object
  - `street`: String
  - `city`: String
  - `state`: String
  - `postal_code`: String
  - `country`: String
- `business_settings`: JSON object
  - `supported_currencies`: Array[String]
  - `supported_languages`: Array[String]
  - `shipping_countries`: Array[String]
  - `tax_settings`: Object
- `compliance_data`: JSON object (regulatory information)
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Relationships**:
- One-to-many with Product entities
- One-to-many with Order entities
- Many-to-many with User entities (customers)

### Message Entity
**Purpose**: Real-time communication content between users
**Storage**: ScyllaDB (primary), Redis (recent cache)

**Attributes**:
- `id`: UUID (primary key)
- `dialog_id`: UUID (conversation identifier)
- `sender_id`: UUID (foreign key to User)
- `message_type`: Enum (text, image, video, audio, file, location, sticker, system)
- `content`: Text (encrypted)
- `metadata`: JSON object (type-specific data)
  - File messages: `file_name`, `file_size`, `mime_type`, `url`
  - Location messages: `latitude`, `longitude`, `address`
  - Media messages: `duration`, `dimensions`, `thumbnail_url`
- `reply_to_id`: UUID (optional, references another message)
- `reactions`: Array of reaction objects
  - `emoji`: String
  - `user_ids`: Array[UUID]
- `delivery_status`: Enum (sent, delivered, read)
- `read_receipts`: Array of read receipt objects
  - `user_id`: UUID
  - `read_at`: Timestamp
- `edited_at`: Timestamp (optional)
- `created_at`: Timestamp

**Validation Rules**:
- Content encryption required for all message types
- File size limits based on message type
- Sender must be participant in dialog
- Reply references must exist in same dialog

**Relationships**:
- Many-to-one with Dialog entity
- Many-to-one with User entity (sender)
- Self-referencing for reply chains

### Dialog Entity
**Purpose**: Conversation containers for users and groups
**Storage**: PostgreSQL (metadata), ScyllaDB (message timelines)

**Attributes**:
- `id`: UUID (primary key)
- `type`: Enum (direct, group, channel)
- `name`: String (optional for direct, required for group/channel)
- `description`: Text (optional)
- `avatar_url`: String (optional)
- `participant_count`: Integer
- `participants`: Array[UUID] (user IDs)
- `admin_ids`: Array[UUID] (for groups/channels)
- `settings`: JSON object
  - `is_public`: Boolean
  - `allow_media`: Boolean
  - `message_retention_days`: Integer
  - `moderation_enabled`: Boolean
- `last_message_id`: UUID (optional)
- `last_message_at`: Timestamp
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Validation Rules**:
- Direct dialogs must have exactly 2 participants
- Group dialogs must have 3-1000 participants
- At least one admin required for groups/channels
- Participant limit enforcement

### Transaction Entity
**Purpose**: Financial operations including payments, transfers, and conversions
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `user_id`: UUID (foreign key to User)
- `wallet_id`: UUID (foreign key to Wallet)
- `type`: Enum (deposit, withdrawal, transfer_in, transfer_out, payment, refund)
- `amount`: Decimal (high precision)
- `currency`: String (THB, SGD, IDR, MYR, PHP, VND, USD)
- `status`: Enum (pending, processing, completed, failed, cancelled)
- `gateway`: String (stripe, omise, bank_transfer)
- `gateway_transaction_id`: String (external reference)
- `exchange_rate`: Decimal (if currency conversion involved)
- `fees`: JSON object
  - `processing_fee`: Decimal
  - `conversion_fee`: Decimal
  - `total_fees`: Decimal
- `description`: Text
- `metadata`: JSON object (gateway-specific data)
- `audit_trail`: JSON array (status change history)
- `created_at`: Timestamp
- `updated_at`: Timestamp
- `completed_at`: Timestamp (optional)

**Validation Rules**:
- Amount must be positive
- Currency must be supported for user's country
- Status transitions must follow business rules
- Audit trail required for all status changes

### Wallet Entity
**Purpose**: Multi-currency digital wallets for users
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `user_id`: UUID (foreign key to User)
- `currency`: String (one wallet per currency per user)
- `balance`: Decimal (available balance)
- `frozen_balance`: Decimal (pending transactions)
- `total_balance`: Decimal (computed: balance + frozen_balance)
- `status`: Enum (active, suspended, closed)
- `daily_limit`: Decimal (regulatory compliance)
- `monthly_limit`: Decimal
- `settings`: JSON object
  - `auto_convert`: Boolean
  - `preferred_currency`: String
  - `notification_preferences`: Object
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Validation Rules**:
- One wallet per currency per user
- Balance cannot be negative
- Frozen balance tracks pending transactions
- Limits enforced based on user country regulations

### Product Entity
**Purpose**: Items listed for sale by businesses
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `business_id`: UUID (foreign key to Business)
- `name`: String (max 200 chars)
- `description`: Text
- `category`: String
- `brand`: String (optional)
- `sku`: String (unique within business)
- `price`: Decimal
- `currency`: String
- `stock_quantity`: Integer
- `status`: Enum (active, inactive, out_of_stock)
- `images`: Array[String] (URLs)
- `attributes`: JSON object (product-specific properties)
  - `color`: String
  - `size`: String
  - `weight`: Decimal
  - `dimensions`: Object
- `specifications`: JSON object (technical specs)
- `variants`: Array of variant objects
  - `id`: UUID
  - `name`: String
  - `sku`: String
  - `price_adjustment`: Decimal
  - `stock_quantity`: Integer
  - `attributes`: Object
- `seo_data`: JSON object
  - `meta_title`: String
  - `meta_description`: String
  - `keywords`: Array[String]
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Validation Rules**:
- Price must be positive
- Currency must match business supported currencies
- Stock quantity cannot be negative
- At least one image required for active products

### Order Entity
**Purpose**: Purchase transactions between users and businesses
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `customer_id`: UUID (foreign key to User)
- `business_id`: UUID (foreign key to Business)
- `status`: Enum (pending, confirmed, processing, shipped, delivered, cancelled, refunded)
- `items`: JSON array of order items
  - `product_id`: UUID
  - `variant_id`: UUID (optional)
  - `quantity`: Integer
  - `unit_price`: Decimal
  - `total_price`: Decimal
- `subtotal`: Decimal
- `tax_amount`: Decimal
- `shipping_cost`: Decimal
- `total_amount`: Decimal
- `currency`: String
- `payment_method`: String
- `payment_status`: Enum (pending, paid, failed, refunded)
- `shipping_address`: JSON object
- `tracking_number`: String (optional)
- `notes`: Text (customer notes)
- `created_at`: Timestamp
- `updated_at`: Timestamp
- `shipped_at`: Timestamp (optional)
- `delivered_at`: Timestamp (optional)

### Notification Entity
**Purpose**: Multi-channel system notifications and messages
**Storage**: PostgreSQL

**Attributes**:
- `id`: UUID (primary key)
- `recipient_id`: UUID (foreign key to User)
- `type`: Enum (email, sms, push, in_app)
- `channel`: String (category: order_updates, messaging, payments, etc.)
- `template_id`: UUID (optional, foreign key to NotificationTemplate)
- `subject`: String
- `content`: Text
- `variables`: JSON object (template variable values)
- `status`: Enum (pending, sent, delivered, failed, expired)
- `priority`: Enum (low, medium, high, urgent)
- `scheduled_at`: Timestamp (optional)
- `sent_at`: Timestamp (optional)
- `delivered_at`: Timestamp (optional)
- `read_at`: Timestamp (optional)
- `expires_at`: Timestamp (optional)
- `retry_count`: Integer
- `metadata`: JSON object (delivery details)
- `created_at`: Timestamp
- `updated_at`: Timestamp

### Session Entity
**Purpose**: User authentication sessions with security context
**Storage**: Redis (primary), PostgreSQL (audit)

**Attributes**:
- `id`: UUID (primary key)
- `user_id`: UUID (foreign key to User)
- `device_info`: JSON object
  - `device_id`: String
  - `platform`: String (web, mobile_ios, mobile_android)
  - `app_version`: String
  - `user_agent`: String
  - `ip_address`: String
- `access_token`: String (JWT)
- `refresh_token`: String
- `expires_at`: Timestamp
- `refresh_expires_at`: Timestamp
- `status`: Enum (active, expired, revoked)
- `last_activity_at`: Timestamp
- `created_at`: Timestamp

**Security Rules**:
- Access tokens expire in 1 hour
- Refresh tokens expire in 30 days
- Maximum 5 active sessions per user
- IP address tracking for security alerts

## Entity Relationships

### Primary Relationships
- User (1) → (N) Business (ownership)
- User (1) → (N) Wallet (multi-currency)
- User (1) → (N) Transaction (financial history)
- User (N) ↔ (N) Dialog (participants)
- Dialog (1) → (N) Message (conversation)
- Business (1) → (N) Product (catalog)
- User (1) → (N) Order (purchase history)
- Business (1) → (N) Order (sales)
- User (1) → (N) Notification (alerts)
- User (1) → (N) Session (authentication)

### Cross-Service Relationships
- Messages reference Users across services
- Orders link Products and Users across commerce/auth services
- Notifications can reference any entity type
- Transactions link to external payment gateways

## Data Access Patterns

### Read Patterns
- User authentication: High frequency, cached in Redis
- Message retrieval: Time-based queries, ScyllaDB optimized
- Product search: Full-text search, indexed queries
- Transaction history: Paginated, date-range queries
- Business analytics: Aggregated queries, read replicas

### Write Patterns
- User registration: ACID transaction, PostgreSQL
- Message sending: High throughput, ScyllaDB batch writes
- Payment processing: Two-phase commit, PostgreSQL
- Order creation: Multi-table transaction, PostgreSQL
- Notification delivery: Asynchronous, queue-based

### Consistency Requirements
- Financial data: Strong consistency (ACID)
- User data: Strong consistency
- Messages: Eventual consistency acceptable
- Product catalog: Eventual consistency acceptable
- Notifications: Eventual consistency acceptable

## Data Retention Policies

### Message Data
- Active conversations: No expiration
- Deleted conversations: 30-day soft delete
- Media files: 1-year retention after deletion
- System messages: 90-day retention

### Financial Data
- Transaction records: 7-year retention (regulatory)
- Wallet balances: Indefinite (active accounts)
- Audit trails: 10-year retention
- Tax records: Per-country requirements

### User Data
- Active accounts: Indefinite
- Deleted accounts: 30-day soft delete, then purged
- Session data: 30-day retention in Redis
- Login audit logs: 1-year retention

### Compliance Data
- GDPR requests: 30-day processing
- Data export: On-demand, 30-day retention
- Audit trails: Immutable, long-term storage
- Consent records: Indefinite retention

## Security Considerations

### Data Encryption
- Messages: End-to-end encryption
- Personal data: Encryption at rest
- Financial data: Field-level encryption
- Passwords: Bcrypt hashing (not stored)

### Access Control
- Row-level security for user data
- Business data isolated by ownership
- Admin access logged and monitored
- API rate limiting per user/endpoint

### Data Privacy
- Personal data minimization
- Consent management integration
- Right to be forgotten implementation
- Cross-border data transfer compliance

This data model provides the foundation for a scalable, compliant, and secure messaging and commerce platform serving Southeast Asian markets with appropriate regional considerations.