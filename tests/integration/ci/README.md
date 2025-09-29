# CI/CD Pipeline Integration for Tchat Integration Tests

This directory contains comprehensive CI/CD pipeline configurations for running Tchat integration tests across multiple platforms and environments.

## 🚀 Available Pipelines

### GitHub Actions (`.github/workflows/integration-tests.yml`)
- **Platform**: GitHub Actions
- **Features**:
  - Quick PR validation tests
  - Full test suite for main branches
  - Load testing on schedule
  - Security testing
  - Mobile testing (iOS/Android)
  - Coverage reporting with Codecov
  - Slack/Discord notifications

### GitLab CI/CD (`gitlab-ci.yml`)
- **Platform**: GitLab CI/CD
- **Features**:
  - Multi-stage pipeline (validate, test, performance, security, deploy)
  - Docker-based service dependencies
  - Parallel test execution
  - Coverage reporting
  - Security scanning with SAST
  - Performance benchmarking
  - Notifications to Slack/Discord

### Jenkins Pipeline (`jenkins-pipeline.groovy`)
- **Platform**: Jenkins
- **Features**:
  - Declarative pipeline
  - Parallel stage execution
  - Parameterized builds
  - Load testing options
  - Mobile testing support
  - Comprehensive reporting
  - Slack notifications

### Azure DevOps (`azure-pipelines.yml`)
- **Platform**: Azure DevOps
- **Features**:
  - Multi-stage YAML pipeline
  - Template-based reusable components
  - Service containers for dependencies
  - Cross-platform testing
  - Deployment validation
  - Teams/Slack notifications

## 📋 Pipeline Features

### Test Execution Stages

#### 1. **Validation Stage**
- Code quality checks (golangci-lint, staticcheck)
- Security scanning (gosec, govulncheck)
- Dependency validation
- License compliance

#### 2. **Quick Tests** (PR only)
- Essential integration tests
- Fast feedback for pull requests
- Database and Redis services
- 15-20 minute execution time

#### 3. **Full Test Suite**
- **Backend Integration**: Complete API testing
- **Frontend Integration**: RTK Query and KMP testing
- **Cross-Platform**: Data synchronization testing
- **Performance**: Load and stress testing
- **Security**: Security-focused integration tests
- **Mobile**: iOS and Android testing

#### 4. **Reporting**
- Combined coverage reports
- Performance metrics
- Security scan results
- Test result aggregation
- Notification dispatch

### Environment Configuration

#### Service Dependencies
```yaml
# PostgreSQL
postgres:15-alpine
- POSTGRES_DB: tchat_test
- POSTGRES_USER: tchat_test
- POSTGRES_PASSWORD: tchat_test_password

# Redis
redis:7-alpine
- Password: tchat_test_password

# Kafka (full tests only)
confluentinc/cp-kafka:7.4.0
- KAFKA_BROKER_ID: 1
- KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

#### Environment Variables
```bash
# Database
DATABASE_URL=postgres://tchat_test:tchat_test_password@localhost:5432/tchat_test?sslmode=disable
REDIS_URL=redis://:tchat_test_password@localhost:6379/0
KAFKA_BROKERS=localhost:9092

# Testing
TEST_TIMEOUT=45m
MAX_PARALLEL=4
CI=true

# Optional
LOAD_TEST_DURATION=600s
LOAD_TEST_RPS=1000
SECURITY_TEST_MODE=true
```

## 🎯 Execution Strategies

### Pull Request Strategy
```yaml
trigger: pull_request
tests: [quick-integration]
duration: ~15 minutes
services: [postgres, redis]
coverage: basic
notifications: PR comments
```

### Main Branch Strategy
```yaml
trigger: push to main/develop
tests: [full-integration, performance, security]
duration: ~60 minutes
services: [postgres, redis, kafka, scylla, minio]
coverage: comprehensive
notifications: Slack/Teams
```

### Scheduled Strategy
```yaml
trigger: nightly (2 AM UTC)
tests: [full-suite, load-testing, mobile]
duration: ~90 minutes
services: [all]
coverage: comprehensive + benchmarks
notifications: Slack/Teams with detailed reports
```

## 📊 Test Metrics and Reporting

### Coverage Targets
- **Backend Integration**: >80% line coverage
- **Frontend Integration**: >75% line coverage
- **Cross-Platform**: >70% integration coverage
- **Combined**: >75% total coverage

### Performance Thresholds
- **API Response Time**: <200ms (95th percentile)
- **Database Queries**: <50ms average
- **Cache Operations**: <10ms average
- **Load Test**: 1000 RPS with <1% error rate

### Quality Gates
- **All tests must pass**: Zero test failures allowed
- **Coverage threshold**: Must meet minimum coverage
- **Security scan**: No high/critical vulnerabilities
- **Performance**: No threshold violations
- **Code quality**: Linting and static analysis pass

## 🔧 Setup Instructions

### GitHub Actions Setup

1. **Repository Secrets**:
   ```
   CODECOV_TOKEN: <your-codecov-token>
   SLACK_WEBHOOK_URL: <slack-webhook-url>
   DISCORD_WEBHOOK_URL: <discord-webhook-url>
   ```

2. **Copy workflow file**:
   ```bash
   cp tests/integration/ci/github-actions.yml .github/workflows/integration-tests.yml
   ```

### GitLab CI Setup

1. **Variables**:
   ```
   SLACK_WEBHOOK_URL: <slack-webhook-url>
   DISCORD_WEBHOOK_URL: <discord-webhook-url>
   ```

2. **Copy CI file**:
   ```bash
   cp tests/integration/ci/gitlab-ci.yml .gitlab-ci.yml
   ```

### Jenkins Setup

1. **Credentials**:
   - `docker-registry`: Docker registry credentials
   - `slack-webhook-url`: Slack webhook URL
   - `codecov-token`: Codecov token

2. **Pipeline Configuration**:
   - Create new Pipeline job
   - Point to `tests/integration/ci/jenkins-pipeline.groovy`
   - Configure parameters as needed

### Azure DevOps Setup

1. **Variable Groups**:
   ```
   Group: tchat-test-variables
   Variables:
   - SLACK_WEBHOOK_URL
   - TEAMS_WEBHOOK_URL
   ```

2. **Copy pipeline file**:
   ```bash
   cp tests/integration/ci/azure-pipelines.yml azure-pipelines.yml
   ```

## 🔍 Troubleshooting

### Common Issues

#### Service Startup Failures
```bash
# Check service logs
docker-compose -f tests/integration/setup/docker-compose.test.yml logs

# Validate environment
docker exec tchat-test-setup /scripts/validate-test-environment.sh
```

#### Test Timeouts
```bash
# Increase timeout in pipeline
TEST_TIMEOUT=60m

# Check resource constraints
docker stats
free -h
```

#### Coverage Issues
```bash
# Generate local coverage
cd tests/integration/setup
go run test_runner.go -coverage=coverage/local -v

# View coverage report
go tool cover -html=coverage/local/coverage.out
```

### Performance Optimization

#### Pipeline Optimization
- Use Docker layer caching
- Implement Go module caching
- Parallel test execution
- Service health checks

#### Resource Management
- Monitor memory usage
- Optimize Docker images
- Use multi-stage builds
- Cleanup intermediate artifacts

### Debugging

#### Local Pipeline Testing
```bash
# Test with Act (GitHub Actions)
act -j quick-tests

# Test with GitLab Runner
gitlab-runner exec docker quick-integration-test

# Test with Azure DevOps CLI
az pipelines run --name tchat-integration
```

#### Test Environment Debugging
```bash
# Access test containers
docker exec -it tchat-postgres-test psql -U tchat_test
docker exec -it tchat-redis-test redis-cli

# Run tests with debugging
cd tests/integration/setup
go run test_runner.go -suites=backend-integration -v -timeout=60m
```

## 📈 Monitoring and Alerts

### Key Metrics to Monitor
- **Test Success Rate**: >99% for main branch
- **Test Execution Time**: <60 minutes for full suite
- **Coverage Trend**: Should not decrease
- **Performance Regression**: No degradation >10%

### Alert Conditions
- Test failures on main branch
- Coverage drops below threshold
- Performance thresholds exceeded
- Security vulnerabilities detected

### Notification Channels
- **Slack**: Real-time alerts and reports
- **Email**: Daily/weekly summaries
- **Dashboard**: Live metrics and trends
- **PR Comments**: Inline feedback

## 🚀 Best Practices

### Pipeline Maintenance
- **Regular Updates**: Keep CI images and tools updated
- **Dependency Management**: Pin versions for reproducibility
- **Resource Monitoring**: Track resource usage and optimize
- **Documentation**: Keep pipeline docs current

### Testing Strategy
- **Fail Fast**: Run quick tests first
- **Parallel Execution**: Maximize concurrency
- **Smart Retries**: Retry flaky tests automatically
- **Environment Isolation**: Ensure test isolation

### Security
- **Secret Management**: Use secure secret storage
- **Image Scanning**: Scan Docker images for vulnerabilities
- **Access Control**: Limit pipeline permissions
- **Audit Logging**: Track pipeline execution and changes