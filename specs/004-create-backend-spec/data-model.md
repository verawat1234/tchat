# Data Model: Backend Services Architecture

**Date**: 2025-09-22
**Feature**: Backend Services Architecture for Telegram SEA Edition

## Service Data Models

### Auth Service Data Model

#### User Entity
```go
type User struct {
    ID          UUID      `json:"id" db:"id"`
    Phone       *string   `json:"phone,omitempty" db:"phone"`
    Email       *string   `json:"email,omitempty" db:"email"`
    Name        string    `json:"name" db:"name"`
    Avatar      *string   `json:"avatar,omitempty" db:"avatar"`
    Country     string    `json:"country" db:"country"` // TH, ID, MY, VN, SG, PH
    Locale      string    `json:"locale" db:"locale"`   // th-TH, id-ID, etc.
    KYCTier     int       `json:"kyc_tier" db:"kyc_tier"` // 1, 2, 3
    Status      string    `json:"status" db:"status"`   // online, offline, away, busy
    LastSeen    *time.Time `json:"last_seen,omitempty" db:"last_seen"`
    IsVerified  bool      `json:"is_verified" db:"is_verified"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

**Validation Rules**:
- Phone OR Email required (not both null)
- Country must be one of: TH, ID, MY, VN, SG, PH
- KYCTier must be 1, 2, or 3
- Name minimum 2 characters, maximum 100 characters

#### Session Entity
```go
type Session struct {
    ID           UUID      `json:"id" db:"id"`
    UserID       UUID      `json:"user_id" db:"user_id"`
    DeviceID     string    `json:"device_id" db:"device_id"`
    AccessToken  string    `json:"access_token" db:"access_token_hash"`
    RefreshToken string    `json:"refresh_token" db:"refresh_token_hash"`
    ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
    IsActive     bool      `json:"is_active" db:"is_active"`
    IPAddress    *string   `json:"ip_address,omitempty" db:"ip_address"`
    UserAgent    *string   `json:"user_agent,omitempty" db:"user_agent"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    LastUsed     time.Time `json:"last_used" db:"last_used"`
}
```

**State Transitions**:
- Created → Active (successful authentication)
- Active → Expired (timeout or manual logout)
- Active → Revoked (security incident)

### Messaging Service Data Model

#### Dialog Entity
```go
type Dialog struct {
    ID            UUID      `json:"id" db:"id"`
    Type          string    `json:"type" db:"type"` // user, group, channel, business
    Name          *string   `json:"name,omitempty" db:"name"`
    Avatar        *string   `json:"avatar,omitempty" db:"avatar"`
    Participants  []UUID    `json:"participants" db:"participants"`
    LastMessageID *UUID     `json:"last_message_id,omitempty" db:"last_message_id"`
    UnreadCount   int       `json:"unread_count" db:"unread_count"`
    IsPinned      bool      `json:"is_pinned" db:"is_pinned"`
    IsArchived    bool      `json:"is_archived" db:"is_archived"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
```

#### Message Entity
```go
type Message struct {
    ID          UUID      `json:"id" db:"id"`
    DialogID    UUID      `json:"dialog_id" db:"dialog_id"`
    SenderID    UUID      `json:"sender_id" db:"sender_id"`
    Type        string    `json:"type" db:"type"` // text, voice, file, image, video, payment
    Content     string    `json:"content" db:"content"` // JSON content based on type
    ReplyToID   *UUID     `json:"reply_to_id,omitempty" db:"reply_to_id"`
    IsEdited    bool      `json:"is_edited" db:"is_edited"`
    IsPinned    bool      `json:"is_pinned" db:"is_pinned"`
    Mentions    []UUID    `json:"mentions" db:"mentions"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    EditedAt    *time.Time `json:"edited_at,omitempty" db:"edited_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
```

**Partitioning Strategy**: Messages partitioned by `dialog_id` for ScyllaDB

### Payment Service Data Model

#### Wallet Entity
```go
type Wallet struct {
    ID               UUID    `json:"id" db:"id"`
    UserID           UUID    `json:"user_id" db:"user_id"`
    Balance          int64   `json:"balance" db:"balance"` // Amount in cents
    Currency         string  `json:"currency" db:"currency"` // THB, SGD, IDR, etc.
    FrozenBalance    int64   `json:"frozen_balance" db:"frozen_balance"`
    DailyLimit       int64   `json:"daily_limit" db:"daily_limit"`
    MonthlyLimit     int64   `json:"monthly_limit" db:"monthly_limit"`
    UsedThisMonth    int64   `json:"used_this_month" db:"used_this_month"`
    Status           string  `json:"status" db:"status"` // active, suspended, closed
    IsPrimary        bool    `json:"is_primary" db:"is_primary"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
```

#### Transaction Entity
```go
type Transaction struct {
    ID            UUID      `json:"id" db:"id"`
    WalletID      UUID      `json:"wallet_id" db:"wallet_id"`
    Type          string    `json:"type" db:"type"` // send, receive, topup, withdraw, purchase
    Amount        int64     `json:"amount" db:"amount"` // Amount in cents
    Currency      string    `json:"currency" db:"currency"`
    Fee           int64     `json:"fee" db:"fee"`
    NetAmount     int64     `json:"net_amount" db:"net_amount"`
    Status        string    `json:"status" db:"status"` // pending, completed, failed, cancelled
    Description   string    `json:"description" db:"description"`
    Reference     *string   `json:"reference,omitempty" db:"reference"`
    CounterpartID *UUID     `json:"counterpart_id,omitempty" db:"counterpart_id"`
    Metadata      string    `json:"metadata" db:"metadata"` // JSON metadata
    BalanceBefore *int64    `json:"balance_before,omitempty" db:"balance_before"`
    BalanceAfter  *int64    `json:"balance_after,omitempty" db:"balance_after"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    ProcessedAt   *time.Time `json:"processed_at,omitempty" db:"processed_at"`
    CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}
```

**State Transitions**:
- Created → Pending (validation complete)
- Pending → Processing (external provider processing)
- Processing → Completed (successful)
- Processing → Failed (external failure)
- Pending/Processing → Cancelled (user/system cancellation)

### Commerce Service Data Model

#### Product Entity
```go
type Product struct {
    ID          UUID      `json:"id" db:"id"`
    ShopID      UUID      `json:"shop_id" db:"shop_id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Images      []string  `json:"images" db:"images"`
    Price       int64     `json:"price" db:"price"` // Price in cents
    Currency    string    `json:"currency" db:"currency"`
    SKU         *string   `json:"sku,omitempty" db:"sku"`
    Inventory   int       `json:"inventory" db:"inventory"`
    Status      string    `json:"status" db:"status"` // draft, active, out_of_stock, discontinued
    IsDigital   bool      `json:"is_digital" db:"is_digital"`
    Weight      *int      `json:"weight,omitempty" db:"weight"` // Weight in grams
    Tags        []string  `json:"tags" db:"tags"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

#### Order Entity
```go
type Order struct {
    ID              UUID      `json:"id" db:"id"`
    OrderNumber     string    `json:"order_number" db:"order_number"`
    UserID          UUID      `json:"user_id" db:"user_id"`
    Status          string    `json:"status" db:"status"` // pending, confirmed, processing, shipped, delivered, cancelled
    Items           string    `json:"items" db:"items"` // JSON array of order items
    Subtotal        int64     `json:"subtotal" db:"subtotal"`
    Tax             int64     `json:"tax" db:"tax"`
    Shipping        int64     `json:"shipping" db:"shipping"`
    Total           int64     `json:"total" db:"total"`
    Currency        string    `json:"currency" db:"currency"`
    ShippingAddress string    `json:"shipping_address" db:"shipping_address"` // JSON address
    PaymentID       *UUID     `json:"payment_id,omitempty" db:"payment_id"`
    TrackingNumber  *string   `json:"tracking_number,omitempty" db:"tracking_number"`
    Notes           *string   `json:"notes,omitempty" db:"notes"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
```

### Notification Service Data Model

#### Notification Entity
```go
type Notification struct {
    ID          UUID      `json:"id" db:"id"`
    UserID      UUID      `json:"user_id" db:"user_id"`
    Type        string    `json:"type" db:"type"` // message, payment, order, system
    Title       string    `json:"title" db:"title"`
    Message     string    `json:"message" db:"message"`
    ActionURL   *string   `json:"action_url,omitempty" db:"action_url"`
    Category    string    `json:"category" db:"category"` // social, transactional, promotional, system
    Priority    string    `json:"priority" db:"priority"` // low, normal, high, urgent
    Channels    []string  `json:"channels" db:"channels"` // push, email, sms, in_app
    Status      string    `json:"status" db:"status"` // pending, sent, delivered, failed
    Data        string    `json:"data" db:"data"` // JSON notification data
    ScheduledAt *time.Time `json:"scheduled_at,omitempty" db:"scheduled_at"`
    SentAt      *time.Time `json:"sent_at,omitempty" db:"sent_at"`
    ReadAt      *time.Time `json:"read_at,omitempty" db:"read_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
```

## Cross-Service Data Relationships

### Event Sourcing Schema
```go
type Event struct {
    ID          UUID      `json:"id" db:"id"`
    AggregateID UUID      `json:"aggregate_id" db:"aggregate_id"`
    Type        string    `json:"type" db:"type"`
    Version     int       `json:"version" db:"version"`
    Data        string    `json:"data" db:"data"` // JSON event data
    Metadata    string    `json:"metadata" db:"metadata"` // JSON metadata
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
```

### Saga State Management
```go
type SagaExecution struct {
    ID           UUID      `json:"id" db:"id"`
    SagaType     string    `json:"saga_type" db:"saga_type"`
    Status       string    `json:"status" db:"status"` // started, completed, failed, compensating
    CurrentStep  int       `json:"current_step" db:"current_step"`
    TotalSteps   int       `json:"total_steps" db:"total_steps"`
    Data         string    `json:"data" db:"data"` // JSON saga data
    CompensationData string `json:"compensation_data" db:"compensation_data"`
    StartedAt    time.Time `json:"started_at" db:"started_at"`
    CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
    FailedAt     *time.Time `json:"failed_at,omitempty" db:"failed_at"`
}
```

## Data Validation Rules

### Global Validation
- All UUIDs must be valid UUID v4 format
- All timestamps in UTC
- All monetary amounts stored as integers (cents)
- All JSON fields must be valid JSON

### Service-Specific Validation
- **Auth Service**: Phone/email uniqueness per country
- **Messaging Service**: Message content size limits based on type
- **Payment Service**: Transaction amount validation against wallet limits
- **Commerce Service**: Inventory validation for order fulfillment
- **Notification Service**: Channel validation against user preferences

### Regional Compliance
- **PDPA**: Personal data fields flagged for consent tracking
- **PCI DSS**: Payment data tokenization, no card storage
- **Data Residency**: User data partitioned by country for localization

## Performance Considerations

### Indexing Strategy
- **Primary Keys**: All entities use UUID primary keys
- **Foreign Keys**: Indexed for join performance
- **Search Fields**: Full-text search indexes on product titles, descriptions
- **Time-Series Data**: Messages and transactions indexed by created_at

### Caching Strategy
- **User Sessions**: Redis cache with 15-minute TTL
- **Product Catalog**: Redis cache with 1-hour TTL
- **Wallet Balances**: Redis cache with 5-minute TTL
- **Notification Preferences**: Redis cache with 24-hour TTL

### Partitioning Strategy
- **Messages**: Partitioned by dialog_id in ScyllaDB
- **Transactions**: Partitioned by wallet_id
- **Events**: Partitioned by aggregate_id
- **Regional Data**: Separate database instances per country