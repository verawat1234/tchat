SHELL := /bin/bash

# Docker configuration
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml
DOCKER_COMPOSE_SERVICES = docker-compose -f docker-compose.services.yml

.PHONY: help
help: ## Show this help message
	@echo "Tchat Development Commands"
	@echo "=========================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: bootstrap
bootstrap: ## Install toolchains and dependencies (manual per stack for now)
	@echo "→ Ensure Go, Node.js, Android Studio, and Xcode CLT are installed."

## Infrastructure Commands

.PHONY: infra-up
infra-up: ## Start infrastructure services (PostgreSQL, Redis, Kafka, etc.)
	$(DOCKER_COMPOSE_DEV) up -d
	@echo "→ Infrastructure services started. Waiting for health checks..."
	$(DOCKER_COMPOSE_DEV) ps

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	$(DOCKER_COMPOSE_DEV) down

.PHONY: infra-reset
infra-reset: ## Reset infrastructure (remove volumes and restart)
	$(DOCKER_COMPOSE_DEV) down -v
	$(DOCKER_COMPOSE_DEV) up -d
	@echo "→ Infrastructure services reset and restarted"

.PHONY: infra-logs
infra-logs: ## Show infrastructure service logs
	$(DOCKER_COMPOSE_DEV) logs -f

## Microservices Commands

.PHONY: services-build
services-build: ## Build all microservice Docker images
	$(DOCKER_COMPOSE_SERVICES) build

.PHONY: services-up
services-up: ## Start all microservices
	$(DOCKER_COMPOSE_SERVICES) up -d
	@echo "→ Microservices started. Waiting for health checks..."
	$(DOCKER_COMPOSE_SERVICES) ps

.PHONY: services-down
services-down: ## Stop all microservices
	$(DOCKER_COMPOSE_SERVICES) down

.PHONY: services-restart
services-restart: ## Restart all microservices
	$(DOCKER_COMPOSE_SERVICES) restart

.PHONY: services-logs
services-logs: ## Show microservice logs
	$(DOCKER_COMPOSE_SERVICES) logs -f

## Full Stack Commands

.PHONY: up
up: infra-up services-build services-up ## Start full development environment
	@echo "→ Full Tchat development environment is running"
	@echo "→ API Gateway: http://localhost:8080"
	@echo "→ Grafana Dashboard: http://localhost:3000 (admin/tchat_grafana_password)"
	@echo "→ Jaeger Tracing: http://localhost:16686"

.PHONY: down
down: services-down infra-down ## Stop full development environment
	@echo "→ Full Tchat development environment stopped"

.PHONY: restart
restart: down up ## Restart full development environment

## Development Commands

.PHONY: dev-backend
dev-backend: infra-up ## Start infrastructure for backend development
	@echo "→ Infrastructure ready for backend development"
	@echo "→ PostgreSQL: localhost:5432"
	@echo "→ Redis: localhost:6379"
	@echo "→ Kafka: localhost:9092"
	@echo "→ ScyllaDB: localhost:9042"

.PHONY: build-auth
build-auth: ## Build auth service
	cd backend/auth && go build -o bin/auth-service ./main.go

.PHONY: build-messaging
build-messaging: ## Build messaging service
	cd backend/messaging && go build -o bin/messaging-service ./main.go

.PHONY: build-payment
build-payment: ## Build payment service
	cd backend/payment && go build -o bin/payment-service ./main.go

.PHONY: build-commerce
build-commerce: ## Build commerce service
	cd backend/commerce && go build -o bin/commerce-service ./main.go

.PHONY: build-notification
build-notification: ## Build notification service
	cd backend/notification && go build -o bin/notification-service ./main.go

.PHONY: build-gateway
build-gateway: ## Build API gateway service
	cd backend/infrastructure/gateway && go build -o bin/gateway-service ./main.go

.PHONY: build-content
build-content: ## Build content service
	cd backend/content && go build -o bin/content-service ./main.go

.PHONY: build-all
build-all: build-auth build-messaging build-payment build-commerce build-notification build-content build-gateway ## Build all services

## Testing Commands

.PHONY: test-backend
test-backend: ## Run backend tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-auth
test-auth: ## Run auth service tests
	cd backend/auth && go test -v -race ./...

.PHONY: test-messaging
test-messaging: ## Run messaging service tests
	cd backend/messaging && go test -v -race ./...

.PHONY: test-payment
test-payment: ## Run payment service tests
	cd backend/payment && go test -v -race ./...

.PHONY: test-commerce
test-commerce: ## Run commerce service tests
	cd backend/commerce && go test -v -race ./...

.PHONY: test-notification
test-notification: ## Run notification service tests
	cd backend/notification && go test -v -race ./...

.PHONY: test-content
test-content: ## Run content service tests
	cd backend/content && go test -v -race ./...

## Database Commands

.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "→ Running database migrations..."
	$(DOCKER_COMPOSE_DEV) exec postgres psql -U tchat_user -d tchat_dev -c "SELECT 'Migrations would run here';"

.PHONY: db-seed
db-seed: ## Seed database with development data
	@echo "→ Seeding database with development data..."
	$(DOCKER_COMPOSE_DEV) exec postgres psql -U tchat_user -d tchat_dev -c "SELECT 'Seed data would be inserted here';"

.PHONY: db-reset
db-reset: ## Reset database (drop and recreate)
	$(DOCKER_COMPOSE_DEV) exec postgres psql -U tchat_user -c "DROP DATABASE IF EXISTS tchat_dev;"
	$(DOCKER_COMPOSE_DEV) exec postgres psql -U tchat_user -c "CREATE DATABASE tchat_dev;"
	@echo "→ Database reset completed"

## Utility Commands

.PHONY: logs
logs: ## Show all service logs
	$(DOCKER_COMPOSE_DEV) logs -f & $(DOCKER_COMPOSE_SERVICES) logs -f

.PHONY: ps
ps: ## Show running containers
	@echo "Infrastructure services:"
	$(DOCKER_COMPOSE_DEV) ps
	@echo ""
	@echo "Microservices:"
	$(DOCKER_COMPOSE_SERVICES) ps

.PHONY: clean
clean: down ## Clean up all containers and volumes
	docker system prune -f
	docker volume prune -f
	@echo "→ Docker cleanup completed"

.PHONY: health
health: ## Check health of all services
	@echo "→ Checking service health..."
	@curl -s http://localhost:8080/health | jq . || echo "API Gateway not available"
	@curl -s http://localhost:8081/health | jq . || echo "Auth Service not available"
	@curl -s http://localhost:8082/health | jq . || echo "Messaging Service not available"
	@curl -s http://localhost:8083/health | jq . || echo "Payment Service not available"
	@curl -s http://localhost:8084/health | jq . || echo "Commerce Service not available"
	@curl -s http://localhost:8085/health | jq . || echo "Notification Service not available"
	@curl -s http://localhost:8086/health | jq . || echo "Content Service not available"
