# Commerce Service - Pact Provider Verification Tests

This directory contains comprehensive **Pact Provider Verification Tests** for the Commerce service, implementing T013 requirements for contract testing infrastructure.

## Overview

These tests validate that the Commerce service meets the API contract expectations defined by consumer services (e.g., Web Frontend). The tests use the **Pact Go v2** framework to ensure API compatibility and prevent breaking changes.

## Test Coverage

The provider verification tests cover the following e-commerce operations:

### Product Catalog Management
- **Product Search**: GET `/commerce/products` with filters (category, price range, shop)
- **Product Details**: GET `/commerce/products/{id}` for individual product information
- **Shop Products**: GET `/commerce/shops/{shop_id}/products` for shop-specific catalogs
- **Product Creation**: POST `/commerce/products` with validation and inventory setup

### Shopping Cart Operations
- **Cart State**: Through order creation with multiple items
- **Inventory Validation**: Real-time stock checking and reservation
- **Price Validation**: Ensuring price consistency between catalog and cart

### Order Management
- **Order Creation**: POST `/commerce/orders` with comprehensive validation
- **Order Retrieval**: GET `/commerce/orders/{id}` and GET `/commerce/orders`
- **Payment Processing**: POST `/commerce/orders/{id}/payment` with payment gateway integration
- **Order Cancellation**: POST `/commerce/orders/{id}/cancel` with inventory restoration
- **Order Status Updates**: Lifecycle management (pending → confirmed → shipped → delivered)

### Shop Management
- **Shop Creation**: POST `/commerce/shops` with business information validation
- **Shop Details**: GET `/commerce/shops/{id}` for shop information
- **Shop Authentication**: Ownership verification for shop operations

## Provider States

The tests implement the following provider states for comprehensive scenario coverage:

### Product States
- `products exist in catalog` - Sets up a diverse product catalog with electronics, fashion, and food items
- `product with ID exists` - Creates a specific product for ID-based retrievals
- `shop has products` - Associates products with a specific shop for shop-catalog testing
- `user has items in cart` - Prepares products for cart/order scenarios

### Authentication States
- `user is authenticated for checkout` - Sets up valid user authentication context
- `user has valid JWT token` - Configures JWT authentication for protected endpoints

### Order States
- `order exists for user` - Creates a confirmed order with payment and shipping details
- `pending order exists` - Sets up a pending order for payment processing tests

### Shop States
- `shop exists with products` - Creates an active shop with associated products
- `user owns shop` - Establishes shop ownership for authorized operations

## Test Data Structure

### Realistic Commerce Data
The tests use **Southeast Asian focused** test data reflecting the target market:

```go
// Sample test product
Product{
    Name: "Premium Wireless Headphones",
    Category: "electronics",
    Price: 299900, // $2999.00 in cents
    Currency: "USD",
    Inventory: {Quantity: 50, Status: "in_stock"},
    Images: [{URL: "https://cdn.tchat.dev/products/...", IsMain: true}]
}

// Sample test order
Order{
    OrderNumber: "ORD-20240101-123456",
    Status: "confirmed",
    PaymentStatus: "paid",
    Currency: "USD",
    ShippingAddress: {City: "Bangkok", Country: "TH"},
    Items: [{ProductID: "...", Quantity: 1, UnitPrice: 9999}]
}
```

### Multi-Currency Support
Tests validate support for Southeast Asian currencies:
- **THB** (Thai Baht), **SGD** (Singapore Dollar), **IDR** (Indonesian Rupiah)
- **MYR** (Malaysian Ringgit), **PHP** (Philippine Peso), **VND** (Vietnamese Dong)
- **USD** (US Dollar) for international transactions

## Running the Tests

### Prerequisites
1. **Go 1.22+** installed
2. **Pact Go v2** framework (automatically installed via go.mod)
3. **Commerce service dependencies** (auth, shared modules)

### Test Execution

```bash
# Navigate to the contract tests directory
cd backend/commerce/tests/contract

# Install dependencies
go mod tidy

# Run the Pact provider verification tests
go test -v ./...

# Run with additional Pact logging
PACT_LOG_LEVEL=DEBUG go test -v ./...

# Run specific test scenarios
go test -v -run TestCommerceProviderContract
```

### Test Configuration

The tests can be configured with environment variables:

```bash
# Pact Broker configuration (if using a Pact Broker)
export PACT_BROKER_BASE_URL="https://your-pact-broker.com"
export PACT_BROKER_TOKEN="your-auth-token"

# Provider service configuration
export PROVIDER_BASE_URL="http://localhost:8080"
export PROVIDER_NAME="Commerce-Service"
export CONSUMER_NAME="Web-Frontend"

# Debug logging
export PACT_LOG_LEVEL="DEBUG"
```

## Integration with CI/CD

### Contract Testing Pipeline
The tests integrate into the development workflow:

1. **Consumer Contract Generation**: Frontend team generates contracts specifying expected API behavior
2. **Provider Verification**: These tests verify the Commerce service implements the expected contracts
3. **Contract Breaking Change Detection**: Automated detection of API changes that break consumer expectations
4. **Deployment Gate**: Contract verification must pass before deployment to prevent API breakage

### Example CI Configuration
```yaml
# .github/workflows/contract-tests.yml
- name: Run Commerce Provider Contract Tests
  run: |
    cd backend/commerce/tests/contract
    go mod tidy
    go test -v ./... -timeout=30m
```

## Mock Implementation Details

The tests use sophisticated mock implementations that closely mirror production behavior:

### MockShopRepository
- **In-memory storage** with UUID-based indexing
- **Owner-based filtering** for multi-tenant shop management
- **Search functionality** with category and location filters
- **Featured shop support** for promotional listings

### MockProductRepository
- **Comprehensive filtering** by category, shop, price range
- **Inventory tracking** with stock status management
- **SKU-based lookup** for product identification
- **Variant support** for products with multiple options (size, color)

### MockOrderRepository
- **Order lifecycle management** (pending → confirmed → shipped → delivered)
- **User-based filtering** for customer order history
- **Shop-based filtering** for merchant order management
- **Status-based queries** for business analytics

### MockPaymentService
- **Payment processing simulation** with realistic response times
- **Refund handling** for order cancellations
- **Payment validation** with status tracking
- **Multi-currency support** for international transactions

## Error Scenarios & Edge Cases

The tests validate proper error handling:

### Authentication Errors
- **Missing JWT tokens** → 401 Unauthorized responses
- **Invalid JWT tokens** → 401 Unauthorized with error details
- **Expired tokens** → 401 Unauthorized with refresh guidance

### Validation Errors
- **Invalid product data** → 400 Bad Request with field-specific errors
- **Insufficient inventory** → 400 Bad Request with stock availability info
- **Invalid currency codes** → 400 Bad Request with supported currencies list

### Business Logic Errors
- **Order cancellation** restrictions based on fulfillment status
- **Price mismatches** between catalog and order creation
- **Shop ownership** validation for protected operations

## Performance Characteristics

The tests validate performance requirements:
- **Response times** < 200ms for GET operations
- **Order creation** < 500ms including inventory checks
- **Payment processing** < 1000ms including gateway communication
- **Concurrent operations** support with proper locking mechanisms

## Security Validation

Security aspects verified:
- **JWT token validation** for all authenticated endpoints
- **User authorization** for shop and order operations
- **Data sanitization** preventing sensitive information exposure
- **Cross-tenant isolation** ensuring users only access their own data

## Troubleshooting

### Common Issues

**Test Failures Due to Missing Dependencies**:
```bash
# Ensure all module dependencies are available
go mod download
go mod verify
```

**Pact Framework Issues**:
```bash
# Update to latest Pact Go version
go get github.com/pact-foundation/pact-go/v2@latest
go mod tidy
```

**Authentication Context Issues**:
- Verify mock authentication middleware is properly configured
- Check that JWT tokens in test data match expected format
- Ensure user context is properly set in provider states

### Debug Logging
Enable detailed logging for troubleshooting:
```bash
export PACT_LOG_LEVEL=DEBUG
export GIN_MODE=debug
go test -v ./... -timeout=10m
```

## Contributing

When adding new contract tests:

1. **Add provider states** for new test scenarios in `setupXXXState` functions
2. **Update mock repositories** to support new data access patterns
3. **Add realistic test data** representing actual Southeast Asian commerce scenarios
4. **Document new test cases** in this README
5. **Validate error scenarios** for comprehensive edge case coverage

## Related Documentation

- [Pact Go Documentation](https://docs.pact.io/implementation_guides/go)
- [Commerce Service API Documentation](../../README.md)
- [Contract Testing Best Practices](../../../docs/contract-testing.md)
- [Southeast Asian E-commerce Requirements](../../../docs/sea-commerce.md)