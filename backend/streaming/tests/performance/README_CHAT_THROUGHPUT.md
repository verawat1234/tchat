# Chat Throughput Performance Testing

## Overview

Comprehensive performance test suite for ScyllaDB chat message throughput validation, targeting 100,000 messages/second with 1,000 concurrent writers.

## Test Suite Components

### 1. TestChatThroughput100K
**Primary throughput validation test**

- **Target**: 100,000 messages/second sustained throughput
- **Configuration**: 1,000 concurrent writers × 100 messages/second × 60 seconds
- **Expected Output**: 6,000,000 total messages
- **Performance Targets**:
  - Throughput: ≥100K msg/s (95% target = 95K msg/s minimum)
  - Write Success Rate: 100% (zero failures)
  - Write Latency P99: ≤5ms

### 2. TestChatWriteLatency
**Write latency distribution analysis**

- **Measurement**: Individual write latency across 1,000 operations
- **Metrics**: P50, P95, P99, Average latency
- **Target**: P99 ≤5ms
- **Warm-up**: 100 operations before measurement

### 3. TestChatQueryPerformance
**Query performance validation**

- **Scenarios**:
  - Recent 50 messages (target: <10ms P99)
  - Last 5 minutes of chat (target: <50ms)
  - Pagination performance (10 pages of 50 messages)
- **Pre-condition**: 1,000 messages inserted
- **Validation**: Timestamp ordering (DESC), message count accuracy

### 4. TestChatTTLEnforcement
**30-day TTL verification**

- **Target**: 2,592,000 seconds (30 days)
- **Validation**: Query TTL from ScyllaDB metadata
- **Tolerance**: ±60 seconds for test execution time
- **Verification**: Message expires after 30 days

### 5. TestScyllaDBPartitioning
**Partition distribution uniformity**

- **Configuration**: 50 streams × 100 messages/stream
- **Metrics**:
  - Partition count (expected: 50)
  - Messages per partition (expected: 100)
  - Standard deviation (target: ≤10% of mean)
- **Validation**: Uniform distribution across partitions

### 6. TestChatBatchPerformance
**Batch insert performance**

- **Batch Sizes**: 10, 50, 100 messages
- **Target**: Average latency per message ≤5ms
- **Validation**: Batch operations faster than individual writes

## Prerequisites

### ScyllaDB Setup

```bash
# Start ScyllaDB using Docker
docker run -d \
  --name scylla-test \
  -p 9042:9042 \
  scylladb/scylla:latest

# Wait for ScyllaDB to be ready
docker exec -it scylla-test cqlsh -e "DESCRIBE KEYSPACES"

# Create test keyspace
docker exec -it scylla-test cqlsh -e "
  CREATE KEYSPACE IF NOT EXISTS tchat_test
  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
"
```

### Environment Configuration

```bash
# Optional: Override defaults
export SCYLLA_HOST="localhost:9042"
export SCYLLA_KEYSPACE="tchat_test"
export TEST_DURATION_SECONDS=60
export NUM_WRITERS=1000
export MESSAGES_PER_SEC=100
```

## Running Tests

### Full Test Suite

```bash
cd /Users/weerawat/Tchat/backend/streaming

# Run all chat throughput tests
go test -v ./tests/performance -run TestChat -timeout 30m

# Run with verbose output and timing
go test -v -timeout 30m ./tests/performance \
  -run TestChat \
  2>&1 | tee chat_throughput_results.log
```

### Individual Tests

```bash
# Throughput test (60 seconds runtime)
go test -v ./tests/performance -run TestChatThroughput100K -timeout 10m

# Write latency distribution
go test -v ./tests/performance -run TestChatWriteLatency -timeout 5m

# Query performance validation
go test -v ./tests/performance -run TestChatQueryPerformance -timeout 5m

# TTL enforcement verification
go test -v ./tests/performance -run TestChatTTLEnforcement -timeout 5m

# Partition distribution analysis
go test -v ./tests/performance -run TestScyllaDBPartitioning -timeout 5m

# Batch insert performance
go test -v ./tests/performance -run TestChatBatchPerformance -timeout 5m
```

### Short Mode (Skip Performance Tests)

```bash
# Skip performance tests in CI/CD
go test -short ./tests/performance
```

### Custom Configuration

```bash
# High-load test: 2000 writers × 100 msg/s × 120 seconds
TEST_DURATION_SECONDS=120 NUM_WRITERS=2000 \
  go test -v ./tests/performance -run TestChatThroughput100K -timeout 15m

# Quick validation: 100 writers × 50 msg/s × 10 seconds
TEST_DURATION_SECONDS=10 NUM_WRITERS=100 MESSAGES_PER_SEC=50 \
  go test -v ./tests/performance -run TestChatThroughput100K -timeout 2m
```

## Expected Output

### Successful Test Run

```
=== RUN   TestChatThroughput100K
    chat_throughput_test.go:203: Starting throughput test: 1000 writers × 100 msg/s × 1m0s = 6000000 messages
    chat_throughput_test.go:263:
        Performance Test Results:
        ==========================
        Duration:              1m0.125s
        Total Messages:        6000000
        Successful Writes:     6000000
        Failed Writes:         0
        Throughput:            99792.11 msg/s
        Write Latency (p50):   2.145ms
        Write Latency (p95):   3.876ms
        Write Latency (p99):   4.521ms
        Write Latency (avg):   2.387ms
        Query Latency (p50):   0s
        Query Latency (p95):   0s
        Query Latency (p99):   0s
        Query Latency (avg):   0s
        Partitions:            100
        Partition StdDev:      142.36

    chat_throughput_test.go:271: ✓ Throughput target: 99792.11 msg/s (target: 100000 msg/s)
    chat_throughput_test.go:272: ✓ Write success rate: 100% (6000000/6000000)
    chat_throughput_test.go:273: ✓ Write latency P99: 4.521ms (target: 5ms)
--- PASS: TestChatThroughput100K (60.13s)
```

### Performance Metrics

| Metric | Target | Typical Result |
|--------|--------|----------------|
| Throughput | ≥100K msg/s | 95K-105K msg/s |
| Write Success Rate | 100% | 100% |
| Write Latency P99 | ≤5ms | 3-5ms |
| Query Latency (50 msgs) | ≤10ms | 5-8ms |
| Query Latency (5 min) | ≤50ms | 20-40ms |
| Partition StdDev | ≤10% | 5-15% |

## Performance Tuning

### ScyllaDB Optimization

```cql
-- Verify table configuration
DESCRIBE TABLE tchat_test.chat_messages;

-- Check compaction strategy
SELECT * FROM system_schema.tables
WHERE keyspace_name = 'tchat_test'
  AND table_name = 'chat_messages';

-- Monitor partition sizes
SELECT stream_id, COUNT(*) as message_count
FROM tchat_test.chat_messages
GROUP BY stream_id
LIMIT 100;
```

### Go Test Optimization

```bash
# Increase test timeout for large-scale tests
go test -timeout 60m ./tests/performance

# Run with race detector (slower, but validates concurrency)
go test -race -timeout 30m ./tests/performance -run TestChat

# Profile memory usage
go test -memprofile=mem.prof ./tests/performance -run TestChatThroughput100K

# Profile CPU usage
go test -cpuprofile=cpu.prof ./tests/performance -run TestChatThroughput100K
```

### Connection Pool Tuning

Modify `test_utils.go` for custom connection settings:

```go
cluster.NumConns = 8              // Increase connections per host
cluster.Timeout = 10 * time.Second // Increase operation timeout
cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(...)
```

## Troubleshooting

### ScyllaDB Connection Failures

```bash
# Check ScyllaDB status
docker ps | grep scylla
docker logs scylla-test

# Test connectivity
nc -zv localhost 9042

# Verify keyspace exists
docker exec -it scylla-test cqlsh -e "DESCRIBE KEYSPACE tchat_test"
```

### Performance Issues

1. **Low Throughput (<95K msg/s)**
   - Check ScyllaDB CPU/memory usage
   - Increase `cluster.NumConns` in test_utils.go
   - Verify network latency to ScyllaDB

2. **High Write Latency (>5ms P99)**
   - Check ScyllaDB compaction status
   - Reduce concurrent writers (NUM_WRITERS)
   - Increase ScyllaDB resources (CPU/memory)

3. **Query Performance Issues (>10ms)**
   - Verify timestamp indexing
   - Check partition sizes (should be balanced)
   - Monitor ScyllaDB read latency

4. **Test Timeout**
   - Increase `-timeout` flag
   - Reduce TEST_DURATION_SECONDS
   - Reduce NUM_WRITERS for quick validation

## Dependencies

- **ScyllaDB**: 5.0+ (compatible with Cassandra 4.0+)
- **Go Modules**:
  - `github.com/gocql/gocql` v1.7.0
  - `github.com/google/uuid` v1.6.0
  - `github.com/stretchr/testify` v1.11.1
- **Repository**: `tchat.dev/streaming/repository`
- **Models**: `tchat.dev/streaming/models`

## Integration with CI/CD

```yaml
# .github/workflows/performance-tests.yml
name: Chat Performance Tests

on:
  push:
    branches: [main, develop]
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM

jobs:
  performance:
    runs-on: ubuntu-latest

    services:
      scylla:
        image: scylladb/scylla:latest
        ports:
          - 9042:9042

    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Wait for ScyllaDB
        run: |
          timeout 60 bash -c 'until docker exec scylla cqlsh -e "DESCRIBE KEYSPACES"; do sleep 2; done'

      - name: Create Keyspace
        run: |
          docker exec scylla cqlsh -e "
            CREATE KEYSPACE IF NOT EXISTS tchat_test
            WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
          "

      - name: Run Performance Tests
        working-directory: backend/streaming
        run: |
          go test -v -timeout 30m ./tests/performance -run TestChat

      - name: Upload Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: performance-results
          path: backend/streaming/chat_throughput_results.log
```

## Success Criteria

✅ **All tests pass with:**
- Zero write failures (100% success rate)
- Throughput ≥95K msg/s (95% of 100K target)
- Write latency P99 ≤5ms
- Query latency P99 ≤10ms (recent 50 messages)
- TTL enforcement verified (~30 days)
- Uniform partition distribution (StdDev ≤10%)

## References

- Task T070: Performance test for chat throughput (100K msg/s)
- Dependencies: T026 (ChatMessage repository), T044 (send chat handler)
- ScyllaDB Documentation: https://docs.scylladb.com/
- Go gocql Driver: https://github.com/gocql/gocql