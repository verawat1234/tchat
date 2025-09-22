# Tchat Infrastructure

Enterprise-grade Docker infrastructure for Southeast Asian chat platform with microservices architecture.

## Quick Start

### 1. Start Infrastructure Services

```bash
# Start PostgreSQL, Redis, Kafka, ScyllaDB, monitoring
make infra-up
```

### 2. Start All Microservices

```bash
# Build and start all microservices
make services-build
make services-up
```

### 3. Full Development Environment

```bash
# Start everything in one command
make up
```

## Infrastructure Components

### Core Databases
- **PostgreSQL 15**: Primary transactional database for user data, orders, payments
- **ScyllaDB 5.4**: High-performance message storage and real-time data
- **Redis 7**: Caching, sessions, real-time features

### Message Streaming
- **Apache Kafka 7.4**: Event streaming between microservices
- **Zookeeper**: Kafka coordination
- **Kafka Connect**: External service integration

### Storage & Files
- **MinIO**: S3-compatible object storage for development

### Monitoring & Observability
- **Prometheus**: Metrics collection and monitoring
- **Grafana**: Metrics visualization and dashboards
- **Jaeger**: Distributed tracing for microservices

## Microservices Architecture

### API Gateway (Port 8080)
- Entry point for all client requests
- Service discovery and load balancing
- Authentication middleware
- Rate limiting and CORS

### Auth Service (Port 8081)
- JWT authentication and session management
- OTP/2FA verification
- KYC compliance for Southeast Asian markets
- User profile management

### Messaging Service (Port 8082)
- Real-time WebSocket connections
- Chat message persistence in ScyllaDB
- Presence management
- Typing indicators and read receipts

### Payment Service (Port 8083)
- Digital wallet management
- Southeast Asian payment gateway integration:
  - **Thailand**: Omise, PromptPay, TrueMoney
  - **Indonesia**: Xendit, Midtrans, Dana, OVO
  - **Singapore**: PayNow
  - **Malaysia**: 2C2P
  - **Philippines**: GCash, PayMaya
- Multi-currency support
- Compliance and regulatory features

### Commerce Service (Port 8084)
- E-commerce marketplace
- Product catalog and inventory
- Order management and fulfillment
- Shop/seller management

### Notification Service (Port 8085)
- Multi-channel notifications (Push, Email, SMS, In-app)
- Template management
- Delivery tracking and analytics
- Regional provider support

### Content Service (Port 8086)
- Dynamic content management system
- Content versioning and publishing
- Category and taxonomy management
- Template-driven content delivery

## Development Workflow

### Backend Development

```bash
# Start infrastructure only
make dev-backend

# Build individual services
make build-auth
make build-messaging
make build-payment

# Run tests
make test-backend
make test-auth
```

### Database Operations

```bash
# Run migrations
make db-migrate

# Seed development data
make db-seed

# Reset database
make db-reset
```

### Monitoring & Debugging

```bash
# Check service health
make health

# View all logs
make logs

# View running containers
make ps
```

## Service URLs

### Core Services
- **API Gateway**: http://localhost:8080
- **Auth Service**: http://localhost:8081
- **Messaging Service**: http://localhost:8082
- **Payment Service**: http://localhost:8083
- **Commerce Service**: http://localhost:8084
- **Notification Service**: http://localhost:8085
- **Content Service**: http://localhost:8086

### Infrastructure
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Kafka**: localhost:9092
- **ScyllaDB**: localhost:9042
- **MinIO**: http://localhost:9000

### Monitoring
- **Grafana**: http://localhost:3000 (admin/tchat_grafana_password)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686

## Environment Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Key configuration areas:
- JWT secrets and authentication
- Database connection strings
- Southeast Asian payment gateway credentials
- Notification service API keys
- Regional settings (currencies, countries)

## Health Checks

All services include comprehensive health checks:

```bash
# Check individual service health
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # Messaging Service
curl http://localhost:8083/health  # Payment Service
curl http://localhost:8084/health  # Commerce Service
curl http://localhost:8085/health  # Notification Service
curl http://localhost:8086/health  # Content Service

# Automated health check
make health
```

## Performance Tuning

### Development Resources
- PostgreSQL: 512MB RAM, 2 CPU cores
- ScyllaDB: 1GB RAM, 1 CPU core (adjustable)
- Redis: 256MB RAM
- Kafka: 1GB RAM

### Production Recommendations
- Scale databases based on load
- Use external managed services for production
- Implement horizontal scaling for microservices
- Configure proper resource limits

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080-8085, 5432, 6379, 9042, 9092 are available
2. **Memory issues**: ScyllaDB requires sufficient memory allocation
3. **Network issues**: Check Docker network configuration

### Debugging Commands

```bash
# View container logs
docker logs tchat-api-gateway
docker logs tchat-postgres-dev

# Execute commands in containers
docker exec -it tchat-postgres-dev psql -U tchat_user -d tchat_dev
docker exec -it tchat-redis-dev redis-cli -a tchat_redis_password

# Network inspection
docker network inspect tchat-dev-network
```

### Clean Restart

```bash
# Complete cleanup and restart
make clean
make up
```

## Security Considerations

### Development Security
- Default passwords for development only
- Services isolated in Docker network
- Health checks exclude sensitive data

### Production Security
- Use secrets management (Docker Secrets, Kubernetes Secrets)
- Enable TLS/SSL for all communications
- Implement proper RBAC
- Use managed database services
- Configure firewalls and VPCs

## Regional Compliance

### Southeast Asian Features
- Multi-currency support (THB, IDR, SGD, MYR, VND, PHP)
- Local payment gateway integration
- Regulatory compliance features
- Regional data residency support

### Supported Countries
- **Thailand**: PromptPay, TrueMoney, Omise integration
- **Indonesia**: Dana, OVO, Xendit, Midtrans integration
- **Singapore**: PayNow integration
- **Malaysia**: Regional payment support
- **Vietnam**: Local currency and payment support
- **Philippines**: GCash, PayMaya integration