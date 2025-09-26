# Tchat Backend Test Infrastructure Makefile
# T007: Configure Test Infrastructure - Enterprise-grade testing framework

.PHONY: help test test-unit test-integration test-contract test-performance clean coverage coverage-html coverage-report test-watch lint fmt vet deps test-deps security audit

# Configuration
GO_VERSION := 1.23.0
COVERAGE_THRESHOLD := 80
TEST_TIMEOUT := 10m
BUILD_DIR := build
COVERAGE_DIR := $(BUILD_DIR)/coverage
REPORTS_DIR := $(BUILD_DIR)/reports

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
NC := \033[0m # No Color

# Default target
help: ## Show this help message
	@echo "$(BLUE)Tchat Backend Test Infrastructure$(NC)"
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Test targets
test: test-unit test-integration test-contract ## Run all tests
	@echo "$(GREEN)✅ All tests completed successfully$(NC)"

test-unit: ## Run unit tests with coverage
	@echo "$(BLUE)🧪 Running unit tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -timeout=$(TEST_TIMEOUT) \
		-coverprofile=$(COVERAGE_DIR)/unit.out \
		-covermode=atomic \
		./auth/... ./content/... ./messaging/... ./payment/... ./notification/... 2>&1 | \
		tee $(REPORTS_DIR)/unit-tests.log || (echo "$(RED)❌ Unit tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Unit tests passed$(NC)"

test-integration: ## Run integration tests
	@echo "$(BLUE)🔗 Running integration tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR) $(REPORTS_DIR)
	@go test -v -race -timeout=$(TEST_TIMEOUT) \
		-coverprofile=$(COVERAGE_DIR)/integration.out \
		-covermode=atomic \
		-tags=integration \
		./tests/integration/... 2>&1 | \
		tee $(REPORTS_DIR)/integration-tests.log || (echo "$(RED)❌ Integration tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Integration tests passed$(NC)"

test-contract: ## Run contract tests
	@echo "$(BLUE)📋 Running contract tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR) $(REPORTS_DIR)
	@go test -v -race -timeout=$(TEST_TIMEOUT) \
		-coverprofile=$(COVERAGE_DIR)/contract.out \
		-covermode=atomic \
		-tags=contract \
		./tests/contract/... 2>&1 | \
		tee $(REPORTS_DIR)/contract-tests.log || (echo "$(RED)❌ Contract tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Contract tests passed$(NC)"

test-performance: ## Run performance tests
	@echo "$(BLUE)⚡ Running performance tests...$(NC)"
	@mkdir -p $(REPORTS_DIR)
	@go test -v -timeout=30m \
		-tags=performance \
		-bench=. \
		-benchmem \
		./tests/performance/... 2>&1 | \
		tee $(REPORTS_DIR)/performance-tests.log || (echo "$(RED)❌ Performance tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Performance tests passed$(NC)"

test-security: ## Run security tests
	@echo "$(BLUE)🛡️ Running security tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR) $(REPORTS_DIR)
	@go test -v -race -timeout=$(TEST_TIMEOUT) \
		-coverprofile=$(COVERAGE_DIR)/security.out \
		-covermode=atomic \
		-tags=security \
		./auth/services/... ./tests/security/... 2>&1 | \
		tee $(REPORTS_DIR)/security-tests.log || (echo "$(RED)❌ Security tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Security tests passed$(NC)"

# Coverage targets
coverage: test-unit test-integration test-contract test-security ## Generate coverage report
	@echo "$(BLUE)📊 Generating coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@echo "mode: atomic" > $(COVERAGE_DIR)/combined.out
	@find $(COVERAGE_DIR) -name "*.out" -not -name "combined.out" -exec tail -n +2 {} \; >> $(COVERAGE_DIR)/combined.out
	@go tool cover -func=$(COVERAGE_DIR)/combined.out | tee $(REPORTS_DIR)/coverage-summary.txt
	@go tool cover -func=$(COVERAGE_DIR)/combined.out | tail -1 | awk '{print $$3}' | sed 's/%//' > $(COVERAGE_DIR)/coverage-percentage.txt
	@$(eval CURRENT_COVERAGE := $(shell cat $(COVERAGE_DIR)/coverage-percentage.txt))
	@if [ "$$(echo "$(CURRENT_COVERAGE) >= $(COVERAGE_THRESHOLD)" | bc)" -eq 1 ]; then \
		echo "$(GREEN)✅ Coverage $(CURRENT_COVERAGE)% meets threshold $(COVERAGE_THRESHOLD)%$(NC)"; \
	else \
		echo "$(RED)❌ Coverage $(CURRENT_COVERAGE)% below threshold $(COVERAGE_THRESHOLD)%$(NC)"; \
		exit 1; \
	fi

coverage-html: coverage ## Generate HTML coverage report
	@echo "$(BLUE)🌐 Generating HTML coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go tool cover -html=$(COVERAGE_DIR)/combined.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)✅ HTML coverage report: $(COVERAGE_DIR)/coverage.html$(NC)"

coverage-report: coverage-html ## Generate comprehensive coverage report
	@echo "$(BLUE)📈 Generating comprehensive coverage report...$(NC)"
	@mkdir -p $(REPORTS_DIR)
	@echo "# Test Coverage Report" > $(REPORTS_DIR)/coverage-report.md
	@echo "Generated: $$(date)" >> $(REPORTS_DIR)/coverage-report.md
	@echo "" >> $(REPORTS_DIR)/coverage-report.md
	@echo "## Overall Coverage" >> $(REPORTS_DIR)/coverage-report.md
	@echo "\`\`\`" >> $(REPORTS_DIR)/coverage-report.md
	@cat $(REPORTS_DIR)/coverage-summary.txt >> $(REPORTS_DIR)/coverage-report.md
	@echo "\`\`\`" >> $(REPORTS_DIR)/coverage-report.md
	@echo "" >> $(REPORTS_DIR)/coverage-report.md
	@echo "## Coverage Threshold" >> $(REPORTS_DIR)/coverage-report.md
	@echo "- **Required**: $(COVERAGE_THRESHOLD)%" >> $(REPORTS_DIR)/coverage-report.md
	@echo "- **Current**: $$(cat $(COVERAGE_DIR)/coverage-percentage.txt)%" >> $(REPORTS_DIR)/coverage-report.md
	@echo "" >> $(REPORTS_DIR)/coverage-report.md
	@echo "## HTML Report" >> $(REPORTS_DIR)/coverage-report.md
	@echo "Open: [coverage.html](../coverage/coverage.html)" >> $(REPORTS_DIR)/coverage-report.md
	@echo "$(GREEN)✅ Coverage report: $(REPORTS_DIR)/coverage-report.md$(NC)"

# Development targets
test-watch: ## Watch for changes and run tests
	@echo "$(BLUE)👀 Watching for test changes...$(NC)"
	@which watchexec > /dev/null || (echo "$(RED)❌ watchexec not found. Install with: brew install watchexec$(NC)" && exit 1)
	@watchexec -e go -w . -- make test-unit

lint: ## Run linting
	@echo "$(BLUE)🔍 Running linting...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)❌ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)" && exit 1)
	@golangci-lint run ./... || (echo "$(RED)❌ Linting failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Linting passed$(NC)"

fmt: ## Format code
	@echo "$(BLUE)📝 Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)🔬 Running go vet...$(NC)"
	@go vet ./... || (echo "$(RED)❌ go vet failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ go vet passed$(NC)"

# Dependency management
deps: ## Download dependencies
	@echo "$(BLUE)📦 Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

test-deps: ## Install test dependencies
	@echo "$(BLUE)🧰 Installing test dependencies...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install gotest.tools/gotestsum@latest
	@which bc > /dev/null || (echo "$(YELLOW)⚠️ bc calculator not found. Install with: brew install bc$(NC)")
	@echo "$(GREEN)✅ Test dependencies installed$(NC)"

# Security and audit
security: ## Run security scan
	@echo "$(BLUE)🔒 Running security scan...$(NC)"
	@which gosec > /dev/null || go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@mkdir -p $(REPORTS_DIR)
	@gosec -fmt json -out $(REPORTS_DIR)/security-scan.json ./... 2>/dev/null || true
	@gosec ./... || (echo "$(YELLOW)⚠️ Security issues found. Check $(REPORTS_DIR)/security-scan.json$(NC)")
	@echo "$(GREEN)✅ Security scan completed$(NC)"

audit: ## Run dependency audit
	@echo "$(BLUE)🔍 Running dependency audit...$(NC)"
	@go list -json -deps ./... | nancy sleuth || (echo "$(YELLOW)⚠️ Vulnerability found in dependencies$(NC)")
	@echo "$(GREEN)✅ Dependency audit completed$(NC)"

# Cleanup
clean: ## Clean build artifacts
	@echo "$(BLUE)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean -testcache
	@echo "$(GREEN)✅ Cleaned$(NC)"

# Quality gates
quality-gate: lint vet test coverage ## Run all quality gates
	@echo "$(GREEN)✅ All quality gates passed$(NC)"

# CI targets
ci-test: deps quality-gate security ## Run full CI test suite
	@echo "$(GREEN)✅ CI test suite completed successfully$(NC)"

# Setup targets
setup: test-deps ## Setup development environment
	@mkdir -p $(BUILD_DIR) $(COVERAGE_DIR) $(REPORTS_DIR)
	@echo "$(GREEN)✅ Development environment setup complete$(NC)"

# Status
status: ## Show test status and coverage
	@echo "$(BLUE)📊 Test Infrastructure Status$(NC)"
	@echo "Go Version: $$(go version)"
	@echo "Coverage Threshold: $(COVERAGE_THRESHOLD)%"
	@if [ -f $(COVERAGE_DIR)/coverage-percentage.txt ]; then \
		echo "Current Coverage: $$(cat $(COVERAGE_DIR)/coverage-percentage.txt)%"; \
	else \
		echo "Current Coverage: Not available (run 'make coverage' first)"; \
	fi
	@echo "Build Directory: $(BUILD_DIR)"
	@echo "Reports Directory: $(REPORTS_DIR)"