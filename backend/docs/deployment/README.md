# Tchat Backend Deployment Guide

This guide covers deployment options and procedures for the Tchat backend microservices.

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+
- PostgreSQL 15+
- Redis 7+
- ScyllaDB 5.2+
- Apache Kafka 3.5+

### Local Development

```bash
# Clone and setup
git clone <repository>
cd tchat/backend

# Start infrastructure services
docker-compose -f docker-compose.dev.yml up -d

# Run migrations
make migrate-up

# Start services
make dev
```

### Production Deployment

```bash
# Build and deploy
make build
docker-compose -f docker-compose.prod.yml up -d

# Run health checks
make health-check
```

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │ Messaging Svc   │    │ Payment Service │
│     :8081       │    │     :8082       │    │     :8083       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌─────────────────┐    ┌─────────────────┐
         │ Notification    │    │   API Gateway   │
         │ Service :8084   │    │     :8080       │
         └─────────────────┘    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
   ┌──────────┐         ┌─────────────┐         ┌─────────────┐
   │PostgreSQL│         │    Redis    │         │ ScyllaDB    │
   │   :5432  │         │    :6379    │         │   :9042     │
   └──────────┘         └─────────────┘         └─────────────┘
                                 │
                        ┌─────────────┐
                        │    Kafka    │
                        │   :9092     │
                        └─────────────┘
```

## Service Configuration

### Environment Variables

Each service requires these environment variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tchat_${SERVICE_NAME}
DB_USER=tchat_user
DB_PASSWORD=${DB_PASSWORD}
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=${REDIS_PASSWORD}

# ScyllaDB Configuration
SCYLLA_HOSTS=localhost:9042
SCYLLA_KEYSPACE=tchat_${SERVICE_NAME}
SCYLLA_REPLICATION_FACTOR=3

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=tchat-${SERVICE_NAME}
KAFKA_TOPIC_PREFIX=tchat

# External Services
TWILIO_ACCOUNT_SID=${TWILIO_ACCOUNT_SID}
TWILIO_AUTH_TOKEN=${TWILIO_AUTH_TOKEN}
STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
OMISE_SECRET_KEY=${OMISE_SECRET_KEY}

# Service Configuration
SERVICE_NAME=${SERVICE_NAME}
SERVICE_PORT=${SERVICE_PORT}
LOG_LEVEL=info
ENVIRONMENT=production
```

## Deployment Options

### 1. Docker Compose (Recommended for Small-Medium Scale)

#### Production Configuration

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  # API Gateway
  api-gateway:
    build:
      context: .
      dockerfile: cmd/gateway/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - SERVICE_NAME=gateway
      - SERVICE_PORT=8080
      - ENVIRONMENT=production
    depends_on:
      - auth-service
      - messaging-service
      - payment-service
      - notification-service
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Auth Service
  auth-service:
    build:
      context: .
      dockerfile: cmd/auth/Dockerfile
    environment:
      - SERVICE_NAME=auth
      - SERVICE_PORT=8081
      - ENVIRONMENT=production
      - DB_NAME=tchat_auth
      - SCYLLA_KEYSPACE=tchat_auth
    depends_on:
      - postgres
      - redis
      - scylla
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Messaging Service
  messaging-service:
    build:
      context: .
      dockerfile: cmd/messaging/Dockerfile
    environment:
      - SERVICE_NAME=messaging
      - SERVICE_PORT=8082
      - ENVIRONMENT=production
      - DB_NAME=tchat_messaging
      - SCYLLA_KEYSPACE=tchat_messaging
    depends_on:
      - postgres
      - redis
      - scylla
      - kafka
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8082/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Payment Service
  payment-service:
    build:
      context: .
      dockerfile: cmd/payment/Dockerfile
    environment:
      - SERVICE_NAME=payment
      - SERVICE_PORT=8083
      - ENVIRONMENT=production
      - DB_NAME=tchat_payment
      - SCYLLA_KEYSPACE=tchat_payment
    depends_on:
      - postgres
      - redis
      - scylla
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8083/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Notification Service
  notification-service:
    build:
      context: .
      dockerfile: cmd/notification/Dockerfile
    environment:
      - SERVICE_NAME=notification
      - SERVICE_PORT=8084
      - ENVIRONMENT=production
      - DB_NAME=tchat_notification
      - SCYLLA_KEYSPACE=tchat_notification
    depends_on:
      - postgres
      - redis
      - scylla
      - kafka
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8084/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Infrastructure Services
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: tchat
      POSTGRES_USER: tchat_user
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped

  scylla:
    image: scylladb/scylla:5.2
    command: --seeds=scylla --smp 1 --memory 2G
    volumes:
      - scylla_data:/var/lib/scylla
    ports:
      - "9042:9042"
    restart: unless-stopped

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    restart: unless-stopped

  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  scylla_data:
```

### 2. Kubernetes (Recommended for Large Scale)

#### Namespace and ConfigMap

```yaml
# k8s/namespace.yml
apiVersion: v1
kind: Namespace
metadata:
  name: tchat

---
# k8s/configmap.yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tchat-config
  namespace: tchat
data:
  environment: "production"
  log_level: "info"
  postgres_host: "postgres-service"
  redis_host: "redis-service"
  scylla_hosts: "scylla-service:9042"
  kafka_brokers: "kafka-service:9092"
```

#### Service Deployments

```yaml
# k8s/auth-service.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: tchat
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: tchat/auth-service:latest
        ports:
        - containerPort: 8081
        env:
        - name: SERVICE_NAME
          value: "auth"
        - name: SERVICE_PORT
          value: "8081"
        envFrom:
        - configMapRef:
            name: tchat-config
        - secretRef:
            name: tchat-secrets
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: tchat
spec:
  selector:
    app: auth-service
  ports:
  - port: 8081
    targetPort: 8081
  type: ClusterIP
```

#### Ingress Configuration

```yaml
# k8s/ingress.yml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tchat-ingress
  namespace: tchat
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/cors-allow-methods: "GET, POST, PUT, DELETE, OPTIONS"
    nginx.ingress.kubernetes.io/cors-allow-headers: "Authorization, Content-Type"
spec:
  tls:
  - hosts:
    - api.tchat.sea
    secretName: tchat-tls
  rules:
  - host: api.tchat.sea
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 8080
```

### 3. AWS ECS/Fargate

#### Task Definition

```json
{
  "family": "tchat-auth-service",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "auth-service",
      "image": "account.dkr.ecr.region.amazonaws.com/tchat/auth-service:latest",
      "portMappings": [
        {
          "containerPort": 8081,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "SERVICE_NAME",
          "value": "auth"
        },
        {
          "name": "SERVICE_PORT",
          "value": "8081"
        },
        {
          "name": "ENVIRONMENT",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:tchat/db-password"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/tchat-auth-service",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8081/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    }
  ]
}
```

#### Service Configuration

```json
{
  "serviceName": "tchat-auth-service",
  "cluster": "tchat-cluster",
  "taskDefinition": "tchat-auth-service",
  "desiredCount": 3,
  "launchType": "FARGATE",
  "networkConfiguration": {
    "awsvpcConfiguration": {
      "subnets": ["subnet-12345", "subnet-67890"],
      "securityGroups": ["sg-12345"],
      "assignPublicIp": "DISABLED"
    }
  },
  "loadBalancers": [
    {
      "targetGroupArn": "arn:aws:elasticloadbalancing:region:account:targetgroup/tchat-auth/123456",
      "containerName": "auth-service",
      "containerPort": 8081
    }
  ],
  "deploymentConfiguration": {
    "maximumPercent": 200,
    "minimumHealthyPercent": 50,
    "deploymentCircuitBreaker": {
      "enable": true,
      "rollback": true
    }
  }
}
```

## Database Setup

### PostgreSQL Schema Migration

```bash
# Install migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations/postgres -database "postgres://user:pass@host:5432/dbname?sslmode=require" up

# Rollback if needed
migrate -path migrations/postgres -database "postgres://user:pass@host:5432/dbname?sslmode=require" down 1
```

### ScyllaDB Schema Setup

```bash
# Create keyspace and tables
cqlsh -h scylla-host -u username -p password < migrations/scylla/001_initial_schema.cql

# Verify tables
cqlsh -h scylla-host -u username -p password -e "DESCRIBE KEYSPACE tchat_messaging;"
```

## Monitoring and Observability

### Health Checks

Each service exposes health endpoints:

- **Liveness**: `GET /health` - Basic service health
- **Readiness**: `GET /health/ready` - Service dependencies health
- **Metrics**: `GET /metrics` - Prometheus metrics

### Logging

Services use structured JSON logging with these fields:

```json
{
  "timestamp": "2023-12-01T10:00:00Z",
  "level": "info",
  "service": "auth-service",
  "request_id": "uuid",
  "user_id": "uuid",
  "message": "User authenticated",
  "duration": "150ms",
  "status_code": 200
}
```

### Metrics

Prometheus metrics are exposed at `/metrics`:

- **Request metrics**: `http_requests_total`, `http_request_duration_seconds`
- **Database metrics**: `db_connections_active`, `db_query_duration_seconds`
- **Business metrics**: `user_registrations_total`, `messages_sent_total`

## Security Considerations

### Network Security

- All inter-service communication uses TLS
- Database connections use SSL/TLS
- Redis connections use AUTH + TLS
- Kafka uses SASL/SSL

### Secrets Management

- Use Kubernetes secrets or AWS Secrets Manager
- Rotate database passwords regularly
- Use IAM roles for AWS services
- Never commit secrets to version control

### API Security

- JWT tokens with 1-hour expiration
- Rate limiting on all endpoints
- CORS configuration for web clients
- Request validation and sanitization

## Scaling Guidelines

### Horizontal Scaling

Services can be scaled independently:

```bash
# Docker Compose
docker-compose up --scale auth-service=3

# Kubernetes
kubectl scale deployment auth-service --replicas=5 -n tchat

# AWS ECS
aws ecs update-service --cluster tchat-cluster --service tchat-auth-service --desired-count 5
```

### Vertical Scaling

Resource requirements per service:

| Service | CPU (cores) | Memory (GB) | Storage |
|---------|-------------|-------------|---------|
| Auth | 0.5-1.0 | 1-2 | Minimal |
| Messaging | 1.0-2.0 | 2-4 | High (ScyllaDB) |
| Payment | 0.5-1.0 | 1-2 | Medium |
| Notification | 0.5-1.0 | 1-2 | Low |

### Database Scaling

- **PostgreSQL**: Read replicas for read-heavy workloads
- **ScyllaDB**: Add nodes to the cluster for horizontal scaling
- **Redis**: Use Redis Cluster for high availability
- **Kafka**: Increase partition count for higher throughput

## Troubleshooting

### Common Issues

1. **Service won't start**
   - Check environment variables
   - Verify database connectivity
   - Check port availability

2. **High latency**
   - Check database connection pool
   - Monitor external service calls
   - Review cache hit rates

3. **Memory leaks**
   - Monitor Go heap size
   - Check for goroutine leaks
   - Review connection pooling

### Debug Commands

```bash
# Check service logs
docker-compose logs -f auth-service

# Database connectivity
go run tools/db-test/main.go

# Load testing
go run tests/load/main.go -service=auth -rps=100

# Memory profiling
go tool pprof http://localhost:8081/debug/pprof/heap
```

## Backup and Disaster Recovery

### Database Backups

```bash
# PostgreSQL backup
pg_dump -h postgres-host -U username dbname > backup.sql

# ScyllaDB backup
nodetool snapshot tchat_messaging

# Redis backup
redis-cli --rdb backup.rdb
```

### Recovery Procedures

1. **Service Failure**: Auto-restart via orchestrator
2. **Database Failure**: Restore from latest backup
3. **Complete Outage**: Follow disaster recovery runbook

## Performance Optimization

### Database Optimization

- Connection pooling (max 100 connections per service)
- Query optimization and indexing
- Read replicas for read-heavy operations
- Caching frequently accessed data

### Service Optimization

- Go runtime tuning: `GOGC=100`, `GOMAXPROCS=auto`
- HTTP/2 for service-to-service communication
- Connection keep-alive and pooling
- Async processing for non-critical operations

### Caching Strategy

- **L1 Cache**: In-memory caching (1-5 minutes TTL)
- **L2 Cache**: Redis caching (5-60 minutes TTL)
- **CDN**: Static content and API responses
- **Database**: Query result caching