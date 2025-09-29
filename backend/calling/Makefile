.PHONY: help build test lint format clean run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=calling-service
BINARY_UNIX=$(BINARY_NAME)_unix

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v .

test: ## Run tests
	$(GOTEST) -v ./...

test-cover: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint
	golangci-lint run --no-config

format: ## Format Go code
	$(GOFMT) -s -w .
	$(GOCMD) fmt ./...
	$(GOCMD) mod tidy

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

run: build ## Build and run the service
	./$(BINARY_NAME)

dev: ## Run in development mode with auto-reload
	$(GOCMD) run ./cmd/server

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(BINARY_NAME)

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

check: format lint test ## Run all checks (format, lint, test)

ci: deps check ## Run CI pipeline

# Build for Linux (useful for Railway deployment)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server