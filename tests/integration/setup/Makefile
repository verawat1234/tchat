# Tchat Integration Testing Makefile
#
# This Makefile provides convenient commands for running integration tests
# across different environments and configurations.

.PHONY: help setup test test-quick test-full test-backend test-frontend test-cross-platform test-performance clean coverage docker-up docker-down docker-logs validate-env

# Default target
.DEFAULT_GOAL := help

# Configuration
TEST_TIMEOUT ?= 30m
TEST_PATTERN ?= ""
MAX_PARALLEL ?= 4
COVERAGE_DIR ?= coverage
DOCKER_COMPOSE_FILE ?= docker-compose.test.yml

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
RESET := \033[0m

## Help target
help: ## Show this help message
	@echo "$(BLUE)Tchat Integration Testing$(RESET)"
	@echo ""
	@echo "$(GREEN)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Examples:$(RESET)"
	@echo "  make setup test-quick          # Quick setup and test"
	@echo "  make test-backend TEST_PATTERN=TestCart*  # Run specific backend tests"
	@echo "  make test-full coverage         # Full test suite with coverage"
	@echo "  make docker-up validate-env     # Setup and validate environment"

## Environment Setup
setup: docker-up validate-env ## Setup complete test environment
	@echo "$(GREEN)✅ Test environment setup completed!$(RESET)"

docker-up: ## Start Docker test services
	@echo "$(BLUE)🐳 Starting Docker test services...$(RESET)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build
	@echo "$(YELLOW)⏳ Waiting for services to be ready...$(RESET)"
	@sleep 30

docker-down: ## Stop Docker test services
	@echo "$(BLUE)🐳 Stopping Docker test services...$(RESET)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans

docker-logs: ## Show Docker service logs
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

validate-env: ## Validate test environment is ready
	@echo "$(BLUE)🔍 Validating test environment...$(RESET)"
	@docker exec tchat-test-setup /scripts/validate-test-environment.sh || true

## Testing Targets
test: test-quick ## Run quick integration tests (alias for test-quick)

test-quick: ## Run essential integration tests quickly
	@echo "$(BLUE)🧪 Running quick integration tests...$(RESET)"
	@go run test_runner.go \
		-suites=backend-integration \
		-timeout=$(TEST_TIMEOUT) \
		-max-parallel=$(MAX_PARALLEL) \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-full: ## Run complete integration test suite
	@echo "$(BLUE)🧪 Running full integration test suite...$(RESET)"
	@go run test_runner.go \
		-suites=all \
		-timeout=60m \
		-max-parallel=$(MAX_PARALLEL) \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-backend: ## Run backend integration tests only
	@echo "$(BLUE)🧪 Running backend integration tests...$(RESET)"
	@go run test_runner.go \
		-suites=backend-integration \
		-timeout=$(TEST_TIMEOUT) \
		-max-parallel=$(MAX_PARALLEL) \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-frontend: ## Run frontend integration tests only
	@echo "$(BLUE)🧪 Running frontend integration tests...$(RESET)"
	@go run test_runner.go \
		-suites=frontend-integration \
		-timeout=$(TEST_TIMEOUT) \
		-max-parallel=$(MAX_PARALLEL) \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-cross-platform: ## Run cross-platform integration tests
	@echo "$(BLUE)🧪 Running cross-platform integration tests...$(RESET)"
	@go run test_runner.go \
		-suites=cross-platform \
		-timeout=$(TEST_TIMEOUT) \
		-max-parallel=$(MAX_PARALLEL) \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-performance: ## Run performance and load tests
	@echo "$(BLUE)🧪 Running performance and load tests...$(RESET)"
	@go run test_runner.go \
		-suites=performance \
		-timeout=30m \
		-max-parallel=2 \
		-v $(if $(TEST_PATTERN),-pattern=$(TEST_PATTERN))

test-with-setup: setup test-full clean ## Complete test cycle with setup and cleanup

## Coverage and Reporting
coverage: ## Generate test coverage reports
	@echo "$(BLUE)📊 Generating test coverage reports...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go run test_runner.go \
		-suites=all \
		-timeout=60m \
		-max-parallel=$(MAX_PARALLEL) \
		-coverage=$(COVERAGE_DIR)/coverage \
		-v
	@echo "$(BLUE)📈 Combining coverage reports...$(RESET)"
	@echo "mode: set" > $(COVERAGE_DIR)/combined.out
	@tail -n +2 $(COVERAGE_DIR)/coverage_*.out >> $(COVERAGE_DIR)/combined.out 2>/dev/null || true
	@go tool cover -html=$(COVERAGE_DIR)/combined.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/combined.out | tail -1
	@echo "$(GREEN)📊 Coverage report generated: $(COVERAGE_DIR)/coverage.html$(RESET)"

coverage-report: coverage ## Generate and open coverage report in browser
	@echo "$(BLUE)🌐 Opening coverage report in browser...$(RESET)"
	@open $(COVERAGE_DIR)/coverage.html 2>/dev/null || \
	 xdg-open $(COVERAGE_DIR)/coverage.html 2>/dev/null || \
	 echo "$(YELLOW)Please open $(COVERAGE_DIR)/coverage.html manually$(RESET)"

## Database and Data Management
migrate-test-db: ## Run database migrations on test database
	@echo "$(BLUE)🗃️  Running database migrations...$(RESET)"
	@docker exec tchat-postgres-test \
		migrate -path /migrations \
		-database "postgres://tchat_test:tchat_test_password@localhost:5432/tchat_test?sslmode=disable" \
		up

reset-test-data: ## Reset test data and reload fixtures
	@echo "$(BLUE)🔄 Resetting test data...$(RESET)"
	@docker exec tchat-test-setup /scripts/cleanup-test-environment.sh
	@echo "$(GREEN)✅ Test data reset completed!$(RESET)"

load-fixtures: ## Load test fixtures into database
	@echo "$(BLUE)📋 Loading test fixtures...$(RESET)"
	@docker exec tchat-test-setup /scripts/setup-test-environment.sh
	@echo "$(GREEN)✅ Test fixtures loaded!$(RESET)"

## Cleanup and Maintenance
clean: ## Clean up test environment and artifacts
	@echo "$(BLUE)🧹 Cleaning up test environment...$(RESET)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans 2>/dev/null || true
	@rm -rf $(COVERAGE_DIR) 2>/dev/null || true
	@echo "$(GREEN)✅ Cleanup completed!$(RESET)"

clean-docker: ## Remove all Docker test containers and volumes
	@echo "$(BLUE)🧹 Removing Docker test containers and volumes...$(RESET)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans --rmi all 2>/dev/null || true
	@docker system prune -f 2>/dev/null || true
	@echo "$(GREEN)✅ Docker cleanup completed!$(RESET)"

## Development and Debugging
debug-env: ## Show environment variables and configuration
	@echo "$(BLUE)🔍 Test Environment Configuration:$(RESET)"
	@echo "TEST_TIMEOUT: $(TEST_TIMEOUT)"
	@echo "TEST_PATTERN: $(TEST_PATTERN)"
	@echo "MAX_PARALLEL: $(MAX_PARALLEL)"
	@echo "COVERAGE_DIR: $(COVERAGE_DIR)"
	@echo "DOCKER_COMPOSE_FILE: $(DOCKER_COMPOSE_FILE)"
	@echo ""
	@echo "$(BLUE)🐳 Docker Services Status:$(RESET)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps 2>/dev/null || echo "No services running"

shell-postgres: ## Open PostgreSQL shell in test database
	@docker exec -it tchat-postgres-test \
		psql -U tchat_test -d tchat_test

shell-redis: ## Open Redis CLI in test instance
	@docker exec -it tchat-redis-test \
		redis-cli -a tchat_test_redis_password

shell-setup: ## Open shell in test setup container
	@docker exec -it tchat-test-setup /bin/bash

logs-postgres: ## Show PostgreSQL logs
	@docker logs -f tchat-postgres-test

logs-redis: ## Show Redis logs
	@docker logs -f tchat-redis-test

logs-kafka: ## Show Kafka logs
	@docker logs -f tchat-kafka-test

## Continuous Integration
ci-test: ## Run tests in CI environment
	@echo "$(BLUE)🤖 Running CI tests...$(RESET)"
	@go run test_runner.go \
		-setup \
		-suites=all \
		-timeout=45m \
		-max-parallel=2 \
		-coverage=$(COVERAGE_DIR)/ci-coverage \
		-cleanup \
		-v

ci-quick: ## Run quick CI tests
	@echo "$(BLUE)🤖 Running quick CI tests...$(RESET)"
	@go run test_runner.go \
		-setup \
		-suites=backend-integration,frontend-integration \
		-timeout=15m \
		-max-parallel=2 \
		-cleanup \
		-v

## Monitoring and Health Checks
health-check: validate-env ## Perform health check on all services

performance-benchmark: ## Run performance benchmarks
	@echo "$(BLUE)⚡ Running performance benchmarks...$(RESET)"
	@cd ../performance && go test -bench=. -benchmem -timeout=10m ./...

load-test: ## Run load testing scenarios
	@echo "$(BLUE)🚀 Running load tests...$(RESET)"
	@cd ../performance && go test -v -run TestLoad -timeout=15m ./...

## Documentation
test-docs: ## Generate test documentation
	@echo "$(BLUE)📚 Generating test documentation...$(RESET)"
	@mkdir -p docs
	@echo "# Integration Test Documentation" > docs/tests.md
	@echo "" >> docs/tests.md
	@echo "## Available Test Suites" >> docs/tests.md
	@echo "" >> docs/tests.md
	@go run test_runner.go -help 2>&1 | grep -A 20 "Available test suites:" >> docs/tests.md
	@echo "$(GREEN)📚 Test documentation generated: docs/tests.md$(RESET)"

## Utility targets for common workflows
dev-setup: setup test-quick ## Quick development setup and test
dev-test: reset-test-data test-backend ## Reset data and run backend tests
dev-full: setup test-full coverage clean ## Full development test cycle