.PHONY: help install dev build test clean docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

dev: ## Start development environment
	docker-compose -f deployments/docker/docker-compose.yml up

dev-backend: ## Start backend in development mode
	cd backend && air

dev-frontend: ## Start frontend in development mode
	cd frontend && npm run dev

build: ## Build all services
	@echo "Building backend..."
	cd backend && go build -o bin/api ./cmd/api
	@echo "Building frontend..."
	cd frontend && npm run build

test: ## Run all tests
	@echo "Running backend tests..."
	cd backend && go test -v ./...
	@echo "Running frontend tests..."
	cd frontend && npm run check

lint: ## Run linters
	@echo "Linting backend..."
	cd backend && gofmt -s -w . && go vet ./...
	@echo "Linting frontend..."
	cd frontend && npm run lint

format: ## Format code
	@echo "Formatting backend..."
	cd backend && gofmt -s -w .
	@echo "Formatting frontend..."
	cd frontend && npm run format

docker-up: ## Start Docker Compose stack
	docker-compose -f deployments/docker/docker-compose.yml up -d

docker-down: ## Stop Docker Compose stack
	docker-compose -f deployments/docker/docker-compose.yml down

docker-logs: ## View Docker logs
	docker-compose -f deployments/docker/docker-compose.yml logs -f

docker-rebuild: ## Rebuild and restart Docker containers
	docker-compose -f deployments/docker/docker-compose.yml up -d --build

db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	docker-compose -f deployments/docker/docker-compose.yml exec postgres psql -U postgres -d hamradio -f /docker-entrypoint-initdb.d/init.sql

db-reset: ## Reset database (WARNING: destroys all data)
	docker-compose -f deployments/docker/docker-compose.yml down -v
	docker-compose -f deployments/docker/docker-compose.yml up -d postgres

clean: ## Clean build artifacts
	rm -rf backend/bin backend/tmp
	rm -rf frontend/.svelte-kit frontend/build
	rm -rf frontend/node_modules
	docker-compose -f deployments/docker/docker-compose.yml down -v
