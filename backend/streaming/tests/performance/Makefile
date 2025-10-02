# Performance Testing Makefile
# Provides convenient commands for running T068-T071 performance tests

.PHONY: help setup setup-compose down test test-all test-load test-latency test-chat test-recording clean

help: ## Show this help message
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Performance Testing Commands"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

setup: ## Setup infrastructure using Docker Compose (recommended)
	@echo "Starting performance test infrastructure with Docker Compose..."
	./setup-compose.sh

setup-legacy: ## Setup infrastructure using legacy bash script
	@echo "Starting ScyllaDB using legacy script..."
	./setup-scylladb.sh

down: ## Stop all infrastructure containers
	@echo "Stopping performance test infrastructure..."
	docker-compose down

logs: ## View ScyllaDB logs
	docker-compose logs -f scylladb

test-all: ## Run all performance tests (requires infrastructure)
	@echo "Running all performance tests (T068-T071)..."
	cd ../.. && go test ./tests/performance -v -timeout 30m

test-load: ## Run T068: Load test (50K concurrent viewers)
	@echo "Running T068: Load Test..."
	cd ../.. && go test ./tests/performance -v -run TestLoad -timeout 15m

test-latency: ## Run T069: Adaptive latency validation
	@echo "Running T069: Latency Test..."
	cd ../.. && go test ./tests/performance -v -run TestAllLatencyTargets -timeout 5m

test-chat: ## Run T070: Chat throughput (requires ScyllaDB)
	@echo "Running T070: Chat Throughput Test..."
	@echo "⚠️  Ensure ScyllaDB is running (make setup)"
	cd ../.. && go test ./tests/performance -v -run TestChat -timeout 10m

test-recording: ## Run T071: Recording performance
	@echo "Running T071: Recording Performance Test..."
	cd ../.. && go test ./tests/performance -v -run TestRecording -timeout 5m

test-quick: ## Run quick validation (T069 + T071 only)
	@echo "Running quick validation tests..."
	cd ../.. && go test ./tests/performance -v -run "TestAllLatencyTargets|TestRecording" -timeout 10m

clean: ## Clean up test artifacts and stop infrastructure
	@echo "Cleaning up test artifacts..."
	docker-compose down
	rm -rf /tmp/recording_perf_test_*
	@echo "✓ Cleanup complete"

status: ## Show infrastructure status
	@echo "Performance Test Infrastructure Status:"
	@echo ""
	docker-compose ps
	@echo ""
	@docker exec tchat-scylla-performance nodetool status 2>/dev/null || echo "⚠️  ScyllaDB not running (use 'make setup')"

# Development helpers
dev-chat-quick: ## Quick chat test (10s duration, 100 writers)
	@echo "Running quick chat throughput validation..."
	TEST_DURATION_SECONDS=10 NUM_WRITERS=100 MESSAGES_PER_SEC=50 \
		cd ../.. && go test ./tests/performance -v -run TestChatWriteLatency -timeout 2m

dev-recording-quick: ## Quick recording test (single test)
	@echo "Running quick recording validation..."
	cd ../.. && go test ./tests/performance -v -run "TestRecordingPerformanceSuite/TestRecordingStartLatency" -timeout 2m