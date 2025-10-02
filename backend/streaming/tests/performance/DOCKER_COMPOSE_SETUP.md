# Docker Compose Performance Testing Setup

Modern infrastructure orchestration for Tchat streaming performance tests.

## Overview

This setup replaces manual Docker container management with Docker Compose orchestration, providing:

- ✅ **One-command setup**: `make setup`
- ✅ **Automated health checks**: Wait for ScyllaDB readiness
- ✅ **Easy cleanup**: `make down`
- ✅ **Reproducible environments**: Version-controlled configuration
- ✅ **CI/CD ready**: GitHub Actions integration

## Architecture

### Docker Compose Stack

```yaml
services:
  scylladb:
    image: scylladb/scylla:5.2
    ports:
      - "9042:9042"   # CQL interface
      - "9180:9180"   # Prometheus metrics
    healthcheck:
      test: ["CMD", "cqlsh", "-e", "DESCRIBE KEYSPACES"]
      interval: 10s
      timeout: 5s
      retries: 12
      start_period: 30s
```

### Network Configuration

- **Bridge network**: `tchat-perf-test`
- **Port mappings**:
  - 9042 → CQL client connections
  - 9180 → Prometheus metrics endpoint

## Migration from Legacy Scripts

### Before (Legacy)
```bash
./setup-scylladb.sh          # Start container manually
# Wait for startup...
# Manual health check
go test ./tests/performance  # Run tests
docker rm -f tchat-scylla    # Manual cleanup
```

### After (Docker Compose)
```bash
make setup      # Automated startup + health check
make test-all   # Run tests
make down       # Automated cleanup
```

## Makefile Commands

### Infrastructure Management

| Command | Purpose | Example |
|---------|---------|---------|
| `make setup` | Start all services | One-time setup |
| `make down` | Stop all services | After testing |
| `make logs` | View service logs | Debugging |
| `make status` | Show service status | Health check |
| `make clean` | Full cleanup | Reset state |

### Test Execution

| Command | Test Suite | Duration | Dependencies |
|---------|-----------|----------|--------------|
| `make test-all` | All tests (T068-T071) | ~30 min | ScyllaDB + FFmpeg |
| `make test-chat` | T070: Chat throughput | ~10 min | ScyllaDB |
| `make test-recording` | T071: Recording | ~5 min | FFmpeg |
| `make test-latency` | T069: Latency | ~5 min | None |
| `make test-load` | T068: Load test | ~15 min | None |
| `make test-quick` | T069 + T071 | ~10 min | FFmpeg |

### Development Helpers

| Command | Purpose | Duration |
|---------|---------|----------|
| `make dev-chat-quick` | Quick chat validation | ~2 min |
| `make dev-recording-quick` | Single recording test | ~2 min |

## Usage Examples

### Quick Start
```bash
cd backend/streaming/tests/performance
make setup
make test-quick
make down
```

### Full Validation
```bash
make setup
make test-all
make down
```

### Development Workflow
```bash
# Start infrastructure once
make setup

# Run tests multiple times
make test-chat
make test-recording

# Make code changes...

# Run tests again
make test-recording

# Stop when done
make down
```

### CI/CD Integration
```bash
# GitHub Actions workflow
docker-compose up -d
sleep 30  # Wait for health checks
make test-all
docker-compose down
```

## Advantages Over Legacy Scripts

### 1. Declarative Configuration
- ✅ Infrastructure as code in `docker-compose.yml`
- ✅ Version-controlled and reviewable
- ✅ Easy to extend with additional services

### 2. Automated Health Checks
- ✅ Built-in health check polling
- ✅ Startup period configuration
- ✅ Automatic retry logic

### 3. Service Dependencies
- ✅ Network creation and management
- ✅ Service discovery
- ✅ Dependency ordering

### 4. Easy Cleanup
- ✅ `make down` removes all resources
- ✅ Network cleanup included
- ✅ Volume management optional

### 5. CI/CD Integration
- ✅ GitHub Actions service containers
- ✅ Reproducible test environments
- ✅ Parallel test execution support

## File Structure

```
tests/performance/
├── docker-compose.yml       # Service definitions
├── setup-compose.sh         # Startup automation
├── Makefile                 # Command shortcuts
├── README.md                # Quick reference
├── DOCKER_COMPOSE_SETUP.md  # This file
└── INFRASTRUCTURE_SETUP.md  # Detailed guide
```

## Extending the Stack

### Adding Redis for Caching

```yaml
services:
  scylladb:
    # ... existing configuration

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    networks:
      - tchat-perf-test
```

### Adding Prometheus for Metrics

```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - tchat-perf-test
```

## Troubleshooting

### Services Won't Start
```bash
# Check Docker is running
docker info

# Check for port conflicts
lsof -i :9042

# View detailed logs
make logs

# Restart infrastructure
make down
make setup
```

### Health Checks Failing
```bash
# Check ScyllaDB logs
docker-compose logs scylladb

# Increase health check retries in docker-compose.yml
retries: 20  # Default: 12

# Increase startup period
start_period: 60s  # Default: 30s
```

### Test Failures
```bash
# Verify infrastructure is healthy
make status

# Check ScyllaDB connectivity
docker exec tchat-scylla-performance cqlsh -e "DESCRIBE KEYSPACES"

# Run tests with verbose output
cd ../.. && go test ./tests/performance -v -run TestChat
```

## Performance Considerations

### Resource Allocation

**Default Configuration**:
- ScyllaDB: 1 CPU, 1GB RAM
- Suitable for: Development, CI/CD
- Test duration: ~30 minutes for full suite

**Production-like Configuration**:
```yaml
services:
  scylladb:
    command:
      - --smp
      - "4"
      - --memory
      - 4G
```

### Optimization Tips

1. **Pre-pull images**: `docker-compose pull` before testing
2. **Use volumes**: Persist ScyllaDB data between runs
3. **Parallel testing**: Run independent tests concurrently
4. **Resource limits**: Set CPU/memory limits for predictable results

## Best Practices

### Development
- ✅ Use `make dev-*` commands for quick iterations
- ✅ Keep infrastructure running during development
- ✅ Stop infrastructure when switching projects

### CI/CD
- ✅ Use service containers in GitHub Actions
- ✅ Set shorter test durations with environment variables
- ✅ Cache Docker images for faster builds
- ✅ Run tests in parallel when possible

### Production
- ✅ Use production-grade ScyllaDB configuration
- ✅ Enable metrics collection (Prometheus)
- ✅ Monitor health check status
- ✅ Set appropriate resource limits

## Migration Checklist

- [x] Create `docker-compose.yml`
- [x] Create `setup-compose.sh` with health checks
- [x] Create `Makefile` with convenient commands
- [x] Update `INFRASTRUCTURE_SETUP.md` with Compose instructions
- [x] Create `README.md` for quick reference
- [x] Test all commands work correctly
- [x] Document CI/CD integration
- [x] Keep legacy `setup-scylladb.sh` for compatibility

## Future Enhancements

### Planned
- [ ] Add Redis for real-time presence
- [ ] Add Prometheus for metrics collection
- [ ] Add Grafana for visualization
- [ ] Add Kafka for event streaming tests
- [ ] Multi-region ScyllaDB cluster setup

### Under Consideration
- [ ] Kubernetes deployment for production-like testing
- [ ] Terraform for cloud infrastructure
- [ ] Automated performance regression detection
- [ ] Load testing dashboard

---

**Last Updated**: 2025-09-30
**Status**: Production Ready ✅