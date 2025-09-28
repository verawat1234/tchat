# Payment Service Provider Verification Tests

This document explains how to run and understand the Pact provider verification tests for the Payment service.

## Overview

The `pact_provider_test.go` file implements comprehensive provider verification tests using Pact Go v2 framework to ensure the Payment service correctly implements the contracts expected by consumers (web frontend, mobile apps).

## Test Coverage

### Provider States Implemented

#### Wallet Operations
- `user has wallet with balance` - User with $1000.00 USD wallet
- `user has multiple currency wallets` - Multi-currency wallets (USD, THB, SGD, IDR)
- `user wallet has insufficient balance` - Low balance scenario ($1.00)
- `user wallet is frozen` - Frozen wallet with all balance locked

#### Payment Method Management
- `payment method exists for user` - Active Visa card
- `user has verified payment method` - Verified Mastercard
- `user has multiple payment methods` - Credit card + GrabPay e-wallet
- `payment method is expired` - Expired payment method
- `payment method requires 3ds verification` - High-risk payment method

#### Transaction Processing
- `transaction can be processed` - Standard $50.00 payment
- `transaction requires additional verification` - Large $1000.00 transfer
- `transaction exceeds daily limit` - User close to daily limit
- `pending transaction exists` - Pending bank deposit

#### Currency & Exchange
- `wallet supports currency operations` - Multi-currency wallet
- `exchange rate is available` - Mock exchange rate service
- `multi-currency transaction is possible` - USD to THB conversion

#### Compliance & Security
- `user is kyc verified` - KYC verified user
- `transaction passes aml checks` - AML compliant transaction
- `user has transaction history` - 5 historical transactions

#### Service Health
- `payment service is healthy` - Service health check
- `external payment processor is available` - External processor status

## API Endpoints Tested

### Wallet Management
- `GET /api/v1/wallets` - List user wallets
- `POST /api/v1/wallets` - Create new wallet
- `GET /api/v1/wallets/:id` - Get wallet details
- `PUT /api/v1/wallets/:id/balance` - Update wallet balance
- `GET /api/v1/wallets/:id/transactions` - Get wallet transactions

### Transaction Processing
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/transactions/:id` - Get transaction details
- `PUT /api/v1/transactions/:id/status` - Update transaction status
- `POST /api/v1/transactions/transfer` - Process transfer
- `POST /api/v1/transactions/payment` - Process payment

### Payment Method Management
- `GET /api/v1/payment-methods` - List payment methods
- `POST /api/v1/payment-methods` - Add payment method
- `GET /api/v1/payment-methods/:id` - Get payment method details
- `PUT /api/v1/payment-methods/:id` - Update payment method
- `DELETE /api/v1/payment-methods/:id` - Remove payment method
- `POST /api/v1/payment-methods/:id/verify` - Verify payment method

## Security Features

### Financial Data Protection
- **Amount Precision**: All amounts stored in cents (int64) to avoid floating-point precision issues
- **Currency Validation**: Strict validation of Southeast Asian currencies (THB, SGD, IDR, MYR, PHP, VND, USD)
- **Transaction Limits**: Daily and monthly limits enforced per currency
- **Status Validation**: Comprehensive status checking for wallets, transactions, and payment methods

### Test Data Security
- **Masked Data**: Credit card numbers masked (****1234)
- **Realistic Limits**: Southeast Asian appropriate daily/monthly limits
- **Risk Scoring**: Payment methods include risk scores and fraud flags
- **Security Levels**: Three-tier security classification (low, medium, high)

### Provider State Isolation
- **Clean State**: Each provider state creates isolated test data
- **No Real Data**: All test data is generated and sandboxed
- **Memory Database**: Uses SQLite in-memory for complete isolation

## Running the Tests

### Prerequisites

1. **Install Pact CLI** (if not already installed):
   ```bash
   curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
   ```

2. **Install Dependencies**:
   ```bash
   cd backend/payment
   go mod tidy
   ```

### Run Provider Verification

```bash
cd backend/payment
go test -v tests/contract/pact_provider_test.go
```

### Environment Variables

Set these environment variables for full integration:

```bash
export PACT_BROKER_BASE_URL="http://localhost:9292"
export PACT_BROKER_USERNAME="admin"
export PACT_BROKER_PASSWORD="admin"
```

### Expected Output

Successful verification output:
```
=== RUN   TestPaymentProviderVerification
INFO[2024-09-24] Verifying provider: payment-service
INFO[2024-09-24] Provider state: user has wallet with balance - setup complete
INFO[2024-09-24] Request matched: GET /api/v1/wallets
INFO[2024-09-24] Response verified successfully
...
--- PASS: TestPaymentProviderVerification (2.45s)
PASS
```

## Test Data Structure

### Sample Wallet
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "balance": 100000,
  "currency": "USD",
  "frozen_balance": 0,
  "daily_limit": 300000,
  "monthly_limit": 9000000,
  "status": "active",
  "is_primary": true
}
```

### Sample Transaction
```json
{
  "id": "uuid",
  "wallet_id": "uuid",
  "type": "payment",
  "status": "pending",
  "currency": "USD",
  "amount": 5000,
  "fee_amount": 150,
  "net_amount": 4850,
  "reference": "PAY_1234567890"
}
```

### Sample Payment Method
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "credit_card",
  "provider": "visa",
  "status": "active",
  "is_verified": true,
  "display_name": "Visa ****1234",
  "last_four_digits": "1234",
  "country": "US",
  "currency": "USD"
}
```

## Currency Support

### Southeast Asian Currencies
- **THB (Thai Baht)**: Daily limit 100,000 THB, Monthly limit 3,000,000 THB
- **SGD (Singapore Dollar)**: Daily limit 5,000 SGD, Monthly limit 150,000 SGD
- **IDR (Indonesian Rupiah)**: Daily limit 15,000,000 IDR, Monthly limit 450,000,000 IDR
- **MYR (Malaysian Ringgit)**: Daily limit 20,000 MYR, Monthly limit 600,000 MYR
- **PHP (Philippine Peso)**: Daily limit 250,000 PHP, Monthly limit 7,500,000 PHP
- **VND (Vietnamese Dong)**: Daily limit 230,000,000 VND, Monthly limit 6,900,000,000 VND
- **USD (US Dollar)**: Daily limit 3,000 USD, Monthly limit 90,000 USD

## Error Handling

The tests verify proper error responses:
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing authentication
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Business logic violations

## Integration with CI/CD

Add to GitHub Actions workflow:

```yaml
- name: Run Payment Provider Verification
  run: |
    cd backend/payment
    go test -v tests/contract/pact_provider_test.go
  env:
    PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_URL }}
    PACT_BROKER_USERNAME: ${{ secrets.PACT_BROKER_USER }}
    PACT_BROKER_PASSWORD: ${{ secrets.PACT_BROKER_PASS }}
```

## Troubleshooting

### Common Issues

1. **Module Import Errors**: Ensure `go.mod` has correct replace directives for shared modules
2. **Pact Broker Connection**: Verify broker URL and credentials
3. **Test Data Conflicts**: Each test run creates fresh test data in memory

### Debug Mode

Enable verbose logging:
```bash
export PACT_LOG_LEVEL=DEBUG
go test -v tests/contract/pact_provider_test.go
```

## Security Considerations

- **Production Separation**: Never run provider tests against production data
- **Sensitive Data**: All test data is synthetic and safe for version control
- **Network Isolation**: Tests use in-memory database and mock external services
- **Access Control**: Provider states validate authentication headers

## Contributing

When adding new provider states:
1. Add state handler function following naming convention
2. Create realistic test data for the scenario
3. Update this documentation
4. Test against consumer contracts

## Support

For issues with contract testing:
- Check Pact Go documentation: https://docs.pact.io/implementation_guides/go
- Review consumer contract expectations
- Validate provider state data matches contract requirements