# Journey API Integration Tests Makefile
# Provides standardized commands for running comprehensive API integration tests

.PHONY: help test test-all test-verbose test-journey test-region clean setup deps check-deps
.DEFAULT_GOAL := help

# Configuration
SHELL := /bin/bash
GO_VERSION := 1.22
TEST_TIMEOUT := 30m
TEST_RESULTS_DIR := results
LOG_FILE := $(TEST_RESULTS_DIR)/journey_tests.log

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

# Default values
ENVIRONMENT ?= localhost
REGION ?= SG
JOURNEY ?=
VERBOSE ?= false

help: ## Display this help message
	@echo "Journey API Integration Tests"
	@echo "============================="
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "Environment Variables:"
	@echo "  ENVIRONMENT  Target environment (localhost, staging, production) [default: localhost]"
	@echo "  REGION       Target region (SG, TH, ID, PH, MY, VN) [default: SG]"
	@echo "  JOURNEY      Specific journey to test (optional)"
	@echo "  VERBOSE      Enable verbose output (true/false) [default: false]"
	@echo ""
	@echo "Examples:"
	@echo "  make test                              # Run all journey tests"
	@echo "  make test-verbose                      # Run with verbose output"
	@echo "  make test-journey JOURNEY=01           # Run registration journey only"
	@echo "  make test-region REGION=TH             # Run tests for Thailand region"
	@echo "  make test ENVIRONMENT=staging          # Run tests on staging environment"

setup: ## Setup test environment and dependencies
	@echo -e "$(BLUE)[INFO]$(NC) Setting up journey test environment..."
	@mkdir -p $(TEST_RESULTS_DIR)
	@chmod +x run_journey_tests.sh
	@echo -e "$(GREEN)[SUCCESS]$(NC) Test environment setup completed"

deps: ## Install and verify Go dependencies
	@echo -e "$(BLUE)[INFO]$(NC) Installing Go dependencies..."
	@go mod tidy
	@go mod download
	@echo -e "$(GREEN)[SUCCESS]$(NC) Dependencies installed"

check-deps: ## Check if all dependencies are available
	@echo -e "$(BLUE)[INFO]$(NC) Checking dependencies..."
	@command -v go >/dev/null 2>&1 || (echo -e "$(RED)[ERROR]$(NC) Go is not installed" && exit 1)
	@command -v jq >/dev/null 2>&1 || echo -e "$(YELLOW)[WARNING]$(NC) jq not found - JSON report parsing may be limited"
	@test -f go.mod || (echo -e "$(RED)[ERROR]$(NC) go.mod not found" && exit 1)
	@echo -e "$(GREEN)[SUCCESS]$(NC) Dependencies check completed"

test: setup check-deps ## Run all journey integration tests
	@echo -e "$(BLUE)[INFO]$(NC) Running all journey integration tests..."
	@echo -e "$(BLUE)[INFO]$(NC) Environment: $(ENVIRONMENT) | Region: $(REGION)"
	@./run_journey_tests.sh -e $(ENVIRONMENT) -r $(REGION)

test-all: test ## Alias for 'make test'

test-verbose: setup check-deps ## Run all journey tests with verbose output
	@echo -e "$(BLUE)[INFO]$(NC) Running journey tests with verbose output..."
	@./run_journey_tests.sh -e $(ENVIRONMENT) -r $(REGION) -v

test-journey: setup check-deps ## Run specific journey test (use JOURNEY=01,02,03,04,05)
	@if [ -z "$(JOURNEY)" ]; then \
		echo -e "$(RED)[ERROR]$(NC) JOURNEY parameter is required. Use: make test-journey JOURNEY=01"; \
		exit 1; \
	fi
	@echo -e "$(BLUE)[INFO]$(NC) Running journey $(JOURNEY) test..."
	@./run_journey_tests.sh -e $(ENVIRONMENT) -r $(REGION) -j journey_0$(JOURNEY)_$(shell \
		case $(JOURNEY) in \
			01) echo "registration" ;; \
			02) echo "messaging" ;; \
			03) echo "ecommerce" ;; \
			04) echo "content" ;; \
			05) echo "crossplatform" ;; \
			*) echo "unknown" ;; \
		esac)

test-region: setup check-deps ## Run tests for specific region
	@echo -e "$(BLUE)[INFO]$(NC) Running tests for region: $(REGION)"
	@./run_journey_tests.sh -e $(ENVIRONMENT) -r $(REGION) -v

# Individual journey test targets
test-registration: ## Run Journey 01: Registration & Onboarding
	@$(MAKE) test-journey JOURNEY=01

test-messaging: ## Run Journey 02: Real-time Messaging
	@$(MAKE) test-journey JOURNEY=02

test-ecommerce: ## Run Journey 03: E-commerce & Payments
	@$(MAKE) test-journey JOURNEY=03

test-content: ## Run Journey 04: Content Creation & Discovery
	@$(MAKE) test-journey JOURNEY=04

test-crossplatform: ## Run Journey 05: Cross-platform Continuity
	@$(MAKE) test-journey JOURNEY=05

# Environment-specific test targets
test-localhost: ## Run all tests on localhost environment
	@$(MAKE) test ENVIRONMENT=localhost

test-staging: ## Run all tests on staging environment
	@$(MAKE) test ENVIRONMENT=staging

test-production: ## Run all tests on production environment
	@$(MAKE) test ENVIRONMENT=production

# Regional test targets
test-singapore: ## Run tests for Singapore region
	@$(MAKE) test-region REGION=SG

test-thailand: ## Run tests for Thailand region
	@$(MAKE) test-region REGION=TH

test-indonesia: ## Run tests for Indonesia region
	@$(MAKE) test-region REGION=ID

test-philippines: ## Run tests for Philippines region
	@$(MAKE) test-region REGION=PH

test-malaysia: ## Run tests for Malaysia region
	@$(MAKE) test-region REGION=MY

test-vietnam: ## Run tests for Vietnam region
	@$(MAKE) test-region REGION=VN

# Advanced test targets
test-comprehensive: setup check-deps ## Run comprehensive test suite with all regions
	@echo -e "$(BLUE)[INFO]$(NC) Running comprehensive test suite across all Southeast Asian regions..."
	@for region in SG TH ID PH MY VN; do \
		echo -e "$(BLUE)[INFO]$(NC) Testing region: $$region"; \
		./run_journey_tests.sh -e $(ENVIRONMENT) -r $$region || exit 1; \
		sleep 2; \
	done
	@echo -e "$(GREEN)[SUCCESS]$(NC) Comprehensive regional testing completed"

test-smoke: setup check-deps ## Run smoke tests (quick validation)
	@echo -e "$(BLUE)[INFO]$(NC) Running smoke tests for quick validation..."
	@go test -short -run TestAllJourneys -timeout 5m

test-stress: setup check-deps ## Run stress tests with extended timeout
	@echo -e "$(BLUE)[INFO]$(NC) Running stress tests with extended timeout..."
	@go test -run TestAllJourneys -timeout $(TEST_TIMEOUT) -v

# Go test commands
test-go: check-deps ## Run Go tests directly without shell wrapper
	@echo -e "$(BLUE)[INFO]$(NC) Running Go tests directly..."
	@go test -v -run TestAllJourneys

test-go-verbose: check-deps ## Run Go tests with verbose output
	@go test -v ./...

test-go-coverage: check-deps ## Run Go tests with coverage report
	@echo -e "$(BLUE)[INFO]$(NC) Running tests with coverage analysis..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o $(TEST_RESULTS_DIR)/coverage.html
	@echo -e "$(GREEN)[SUCCESS]$(NC) Coverage report generated: $(TEST_RESULTS_DIR)/coverage.html"

# Benchmarking
benchmark: check-deps ## Run performance benchmarks
	@echo -e "$(BLUE)[INFO]$(NC) Running performance benchmarks..."
	@go test -bench=. -benchmem -run=^$$ ./...

benchmark-verbose: check-deps ## Run benchmarks with verbose output
	@go test -bench=. -benchmem -benchtime=10s -run=^$$ -v ./...

# Report and analysis
report: ## Display latest test report summary
	@if [ -f "journey_test_results.json" ]; then \
		echo -e "$(BLUE)[INFO]$(NC) Latest test results:"; \
		echo "================================================"; \
		if command -v jq >/dev/null 2>&1; then \
			jq -r '"Environment: " + .environment, "Execution Time: " + .executionTime, "Total Duration: " + (.totalDuration / 1000000000 | tostring) + "s", "Total Journeys: " + (.totalJourneys | tostring), "Passed: " + (.passedJourneys | tostring), "Failed: " + (.failedJourneys | tostring), "Summary: " + .summary' journey_test_results.json; \
		else \
			cat journey_test_results.json; \
		fi; \
		echo "================================================"; \
	else \
		echo -e "$(YELLOW)[WARNING]$(NC) No test results found. Run 'make test' first."; \
	fi

logs: ## View recent test logs
	@if [ -f "$(LOG_FILE)" ]; then \
		echo -e "$(BLUE)[INFO]$(NC) Displaying recent test logs:"; \
		tail -50 $(LOG_FILE); \
	else \
		echo -e "$(YELLOW)[WARNING]$(NC) No log file found at $(LOG_FILE)"; \
	fi

# Maintenance and cleanup
clean: ## Clean up test artifacts and temporary files
	@echo -e "$(BLUE)[INFO]$(NC) Cleaning up test artifacts..."
	@rm -f journey_test_results.json
	@rm -f coverage.out
	@rm -rf $(TEST_RESULTS_DIR)/*.tmp
	@rm -rf $(TEST_RESULTS_DIR)/*.log.*
	@echo -e "$(GREEN)[SUCCESS]$(NC) Cleanup completed"

clean-all: clean ## Clean up all test artifacts and dependencies
	@echo -e "$(BLUE)[INFO]$(NC) Performing comprehensive cleanup..."
	@rm -rf $(TEST_RESULTS_DIR)
	@go clean -cache
	@go clean -testcache
	@echo -e "$(GREEN)[SUCCESS]$(NC) Comprehensive cleanup completed"

# Validation and formatting
format: ## Format Go code
	@echo -e "$(BLUE)[INFO]$(NC) Formatting Go code..."
	@go fmt ./...
	@echo -e "$(GREEN)[SUCCESS]$(NC) Code formatting completed"

lint: ## Run Go linter
	@echo -e "$(BLUE)[INFO]$(NC) Running Go linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo -e "$(YELLOW)[WARNING]$(NC) golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi
	@echo -e "$(GREEN)[SUCCESS]$(NC) Linting completed"

validate: format lint test-smoke ## Run full validation (format, lint, smoke test)
	@echo -e "$(GREEN)[SUCCESS]$(NC) Full validation completed"

# Development helpers
watch: ## Watch for file changes and run tests automatically
	@echo -e "$(BLUE)[INFO]$(NC) Watching for file changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o *.go | xargs -n1 -I{} make test-smoke; \
	else \
		echo -e "$(YELLOW)[WARNING]$(NC) fswatch not installed. Install with: brew install fswatch (macOS) or apt-get install inotify-tools (Linux)"; \
		echo -e "$(BLUE)[INFO]$(NC) Running tests in loop instead..."; \
		while true; do make test-smoke; sleep 30; done; \
	fi

dev: ## Development mode - setup, validate, and run smoke tests
	@$(MAKE) setup deps validate

# CI/CD targets
ci-test: ## Run tests in CI/CD environment
	@echo -e "$(BLUE)[INFO]$(NC) Running tests in CI/CD mode..."
	@$(MAKE) setup deps
	@$(MAKE) test-verbose ENVIRONMENT=localhost
	@$(MAKE) test-go-coverage
	@echo -e "$(GREEN)[SUCCESS]$(NC) CI/CD tests completed"

# Information targets
info: ## Display system and configuration information
	@echo "Journey API Integration Tests - System Information"
	@echo "=================================================="
	@echo "Go Version: $$(go version)"
	@echo "Working Directory: $$(pwd)"
	@echo "Test Results Directory: $(TEST_RESULTS_DIR)"
	@echo "Test Timeout: $(TEST_TIMEOUT)"
	@echo "Current Configuration:"
	@echo "  Environment: $(ENVIRONMENT)"
	@echo "  Region: $(REGION)"
	@echo "  Verbose: $(VERBOSE)"
	@if [ -n "$(JOURNEY)" ]; then echo "  Journey: $(JOURNEY)"; fi
	@echo "Available Test Files:"
	@ls -la journey_*_api_test.go 2>/dev/null || echo "  No journey test files found"
	@echo "=================================================="

journeys: ## List all available journey tests
	@echo "Available Journey Tests:"
	@echo "======================="
	@echo "01 - Registration & Onboarding     (journey_01_registration_api_test.go)"
	@echo "02 - Real-time Messaging          (journey_02_messaging_api_test.go)"
	@echo "03 - E-commerce & Payments        (journey_03_ecommerce_api_test.go)"
	@echo "04 - Content Creation & Discovery (journey_04_content_api_test.go)"
	@echo "05 - Cross-platform Continuity    (journey_05_crossplatform_api_test.go)"
	@echo ""
	@echo "Run specific journey with: make test-journey JOURNEY=XX"

# Quick access aliases
t: test ## Quick alias for test
tv: test-verbose ## Quick alias for test-verbose
tr: report ## Quick alias for report
tc: clean ## Quick alias for clean