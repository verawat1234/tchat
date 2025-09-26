# Tchat Pact Broker Infrastructure

Docker-based Pact Broker setup for cross-platform contract testing with 7 backend microservices and 3 client platforms.

## Architecture

- **Pact Broker**: Contract repository and verification coordinator
- **PostgreSQL**: Persistent storage with performance optimizations
- **Management Scripts**: Automated contract publishing and verification

## Services

### Backend Providers (7 microservices)
- `auth-service` (port 8001) - Authentication and user management
- `content-service` (port 8002) - Content management and delivery
- `commerce-service` (port 8003) - E-commerce and product catalog
- `messaging-service` (port 8004) - Real-time messaging and chat
- `payment-service` (port 8005) - Payment processing and wallets
- `notification-service` (port 8006) - Push notifications and alerts
- `gateway-service` (port 8000) - API gateway and routing

### Client Consumers (3 platforms)
- `tchat-web` - React/TypeScript web application
- `tchat-ios` - SwiftUI iOS application
- `tchat-android` - Jetpack Compose Android application

## Quick Start

### 1. Start Pact Broker Infrastructure

```bash
# Navigate to pact-broker directory
cd backend/infrastructure/pact-broker

# Start Pact Broker and PostgreSQL
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f pact-broker
```

### 2. Configure Environment

```bash
# Source Pact environment configuration
source /Users/weerawat/Tchat/.env.pact

# Verify CLI is available
pact-broker version
```

### 3. Access Pact Broker Web UI

- **URL**: http://localhost:9292
- **Username**: admin
- **Password**: admin
- **Read-only**: viewer/viewer

## Contract Management

### Publish Consumer Contracts

```bash
# Publish all consumer contracts
./scripts/publish-contracts.sh

# Or publish manually
pact-broker publish /Users/weerawat/Tchat/pacts/tchat-web-*.json \
  --consumer-app-version=$(git rev-parse --short HEAD) \
  --branch=$(git branch --show-current) \
  --broker-base-url=http://localhost:9292 \
  --broker-username=admin \
  --broker-password=admin
```

### Verify Provider Contracts

```bash
# Verify all provider contracts (requires running services)
./scripts/verify-contracts.sh

# Or verify individual provider
pact-broker verify \
  --provider=auth-service \
  --provider-app-version=$(git rev-parse --short HEAD) \
  --provider-base-url=http://localhost:8001 \
  --broker-base-url=http://localhost:9292 \
  --broker-username=admin \
  --broker-password=admin \
  --publish-verification-results
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PACT_BROKER_BASE_URL` | Pact Broker URL | `http://localhost:9292` |
| `PACT_BROKER_USERNAME` | Admin username | `admin` |
| `PACT_BROKER_PASSWORD` | Admin password | `admin` |
| `PACT_OUTPUT_DIR` | Pact files directory | `/Users/weerawat/Tchat/pacts` |

### Database Performance

The PostgreSQL database is optimized for contract testing workflows:

- **Shared Buffers**: 256MB for better caching
- **Work Memory**: 4MB per operation
- **Max Connections**: 200 concurrent connections
- **Custom Indexes**: Optimized for consumer/provider lookups

## Monitoring and Debugging

### Health Checks

```bash
# Pact Broker health
curl http://localhost:9292/diagnostic/status/heartbeat

# Database health
docker-compose exec pact-broker-db pg_isready -U pact_broker_user

# View contract statistics
docker-compose exec pact-broker-db psql -U pact_broker_user -d pact_broker \
  -c "SELECT * FROM get_consumer_provider_stats();"
```

### Logs and Debugging

```bash
# View Pact Broker logs
docker-compose logs -f pact-broker

# View database logs
docker-compose logs -f pact-broker-db

# Check slow queries (>1 second)
docker-compose exec pact-broker-db psql -U pact_broker_user -d pact_broker \
  -c "SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### Contract Health Overview

```bash
# View contract verification status
docker-compose exec pact-broker-db psql -U pact_broker_user -d pact_broker \
  -c "SELECT * FROM contract_health_overview ORDER BY consumer, provider;"
```

## Data Management

### Backup and Restore

```bash
# Backup contracts database
docker-compose exec pact-broker-db pg_dump -U pact_broker_user pact_broker > pact_backup.sql

# Restore from backup
docker-compose exec -i pact-broker-db psql -U pact_broker_user pact_broker < pact_backup.sql
```

### Clean Up Old Data

```bash
# Clean up old contract versions (keep last 10 per consumer/provider)
pact-broker clean \
  --broker-base-url=http://localhost:9292 \
  --broker-username=admin \
  --broker-password=admin \
  --keep-version-selectors='{"latest": true, "max_age": 30}'
```

## CI/CD Integration

### GitHub Actions

The Pact Broker integrates with CI/CD pipelines through:

- **Contract Publication**: Automated publishing after successful consumer tests
- **Provider Verification**: Automated verification in provider builds
- **Can-I-Deploy**: Deployment readiness checks
- **Webhooks**: Automated notifications on contract changes

### Environment Promotion

```bash
# Tag version for production deployment
pact-broker create-version-tag \
  --pacticipant=tchat-web \
  --version=$(git rev-parse --short HEAD) \
  --tag=production \
  --broker-base-url=http://localhost:9292 \
  --broker-username=admin \
  --broker-password=admin
```

## Troubleshooting

### Common Issues

1. **Pact Broker not accessible**
   - Check Docker containers are running
   - Verify port 9292 is not in use
   - Check firewall settings

2. **Database connection errors**
   - Verify PostgreSQL container is healthy
   - Check database credentials
   - Review connection limits

3. **Contract verification failures**
   - Ensure provider services are running
   - Check service health endpoints
   - Review contract compatibility

### Reset Infrastructure

```bash
# Stop and remove all containers
docker-compose down -v

# Remove persistent data (WARNING: destroys all contracts)
docker volume rm tchat-pact-broker-db-data

# Restart from clean state
docker-compose up -d
```

## Performance Tuning

The infrastructure is configured for:

- **Contract Verification**: <1 second per contract
- **API Response Time**: <200ms average
- **Concurrent Users**: 200 connections
- **Storage**: Optimized for 10,000+ contract versions

For production deployment, adjust resource limits in `docker-compose.yml` based on usage patterns.