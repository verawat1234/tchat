# Contract Testing Makefile for Tchat Backend
# Provides commands for running Pact contract tests

.PHONY: help install test-consumer test-provider publish-contracts verify-contracts start-broker stop-broker clean

# Default target
help:
	@echo "Available commands:"
	@echo "  help              - Show this help message"
	@echo "  install           - Install Pact CLI tools"
	@echo "  test-consumer     - Run consumer contract tests"
	@echo "  test-provider     - Run provider contract verification tests"
	@echo "  publish-contracts - Publish contracts to Pact Broker"
	@echo "  verify-contracts  - Verify contracts against running providers"
	@echo "  start-broker      - Start Pact Broker infrastructure"
	@echo "  stop-broker       - Stop Pact Broker infrastructure"
	@echo "  clean             - Clean up generated pact files"

# Install Pact CLI tools
install:
	@echo "Installing Pact CLI tools..."
	@go install github.com/pact-foundation/pact-go/v2@latest
	@echo "Pact CLI tools installed successfully"

# Run consumer contract tests
test-consumer:
	@echo "Running consumer contract tests..."
	@mkdir -p /Users/weerawat/Tchat/pacts
	@cd pact && go test -v ./auth_consumer_test.go
	@echo "Consumer tests completed. Pact files generated in /Users/weerawat/Tchat/pacts"

# Run provider contract verification tests
test-provider:
	@echo "Running provider contract verification tests..."
	@echo "Checking if Auth Service is running on port 8081..."
	@curl -f http://localhost:8081/health > /dev/null 2>&1 || (echo "Auth Service is not running. Please start it first." && exit 1)
	@cd pact && AUTH_SERVICE_URL=http://localhost:8081 go test -v ./auth_provider_test.go
	@echo "Provider verification tests completed"

# Publish contracts to Pact Broker
publish-contracts: test-consumer
	@echo "Publishing contracts to Pact Broker..."
	@docker-compose -f ../../infrastructure/pact-broker/docker-compose.yml up -d pact-broker
	@sleep 10
	@cd ../../infrastructure/pact-broker && ./scripts/publish-contracts.sh
	@echo "Contracts published to Pact Broker"

# Verify contracts against running providers
verify-contracts:
	@echo "Verifying contracts against running providers..."
	@docker-compose -f ../../infrastructure/pact-broker/docker-compose.yml up -d pact-broker
	@sleep 10
	@cd ../../infrastructure/pact-broker && ./scripts/verify-contracts.sh
	@echo "Contract verification completed"

# Start Pact Broker infrastructure
start-broker:
	@echo "Starting Pact Broker infrastructure..."
	@cd ../../infrastructure/pact-broker && docker-compose up -d
	@echo "Waiting for Pact Broker to be ready..."
	@sleep 15
	@echo "Pact Broker available at: http://localhost:9292"
	@echo "Login credentials: admin/admin"

# Stop Pact Broker infrastructure
stop-broker:
	@echo "Stopping Pact Broker infrastructure..."
	@cd ../../infrastructure/pact-broker && docker-compose down
	@echo "Pact Broker stopped"

# Clean up generated pact files
clean:
	@echo "Cleaning up generated pact files..."
	@rm -rf /Users/weerawat/Tchat/pacts/*.json
	@echo "Pact files cleaned up"

# Run full contract testing workflow
test-all: clean test-consumer publish-contracts test-provider
	@echo "Full contract testing workflow completed successfully!"

# Run contract tests with local pacts (no broker required)
test-local: clean test-consumer test-provider
	@echo "Local contract testing completed successfully!"

# Continuous Integration workflow
ci: install test-all
	@echo "CI contract testing workflow completed!"