.PHONY: help build up down restart logs clean test dev prod

# Docker Compose command (use 'docker compose' for V2)
DOCKER_COMPOSE = docker compose
DOCKER_COMPOSE_PROD = docker compose -f docker-compose.prod.yml

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev           - Start development environment"
	@echo "  make prod          - Start production environment"
	@echo "  make build-dev     - Build development Docker images"
	@echo "  make build-prod    - Build production Docker images"
	@echo "  make up            - Start development containers"
	@echo "  make down          - Stop and remove containers"
	@echo "  make restart       - Restart development containers"
	@echo "  make logs          - View logs for development"
	@echo "  make logs-prod     - View logs for production"
	@echo "  make clean         - Remove containers, volumes, and images"
	@echo "  make migration     - Create new migration file"
	@echo "  make db-migrate    - Run database migrations"
	@echo "  make db-shell      - Access PostgreSQL shell"
	@echo "  make app-shell     - Access app container shell"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format Go code"

# Development commands
dev: build-dev up logs-app

build-dev:
	@echo "Building development images..."
	$(DOCKER_COMPOSE) build

up:
	@echo "Starting development environment..."
	$(DOCKER_COMPOSE) up -d
	@echo "Development environment is running at http://localhost:8010"

down:
	@echo "Stopping containers..."
	$(DOCKER_COMPOSE) down

restart: down up

logs:
	$(DOCKER_COMPOSE) logs -f

logs-app:
	$(DOCKER_COMPOSE) logs -f app

logs-db:
	$(DOCKER_COMPOSE) logs -f postgres

# Production commands
prod: build-prod up-prod

build-prod:
	@echo "Building production images..."
	$(DOCKER_COMPOSE_PROD) build

up-prod:
	@echo "Starting production environment..."
	$(DOCKER_COMPOSE_PROD) up -d
	@echo "Production environment is running"

down-prod:
	@echo "Stopping production containers..."
	$(DOCKER_COMPOSE_PROD) down

restart-prod: down-prod up-prod

logs-prod:
	$(DOCKER_COMPOSE_PROD) logs -f

logs-app-prod:
	$(DOCKER_COMPOSE_PROD) logs -f app

logs-db-prod:
	$(DOCKER_COMPOSE_PROD) logs -f postgres

# Migration commands
migration:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

migration-go:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name go

# Database commands
db-shell:
	@echo "Connecting to PostgreSQL..."
	docker exec -it gobox-postgres psql -U gobox -d goboxdb

db-shell-prod:
	@echo "Connecting to production PostgreSQL..."
	docker exec -it gobox-postgres-prod psql -U $$(grep POSTGRES_USER .env.prod | cut -d '=' -f2) -d $$(grep POSTGRES_DB .env.prod | cut -d '=' -f2)

db-migrate:
	@echo "Running migrations..."
	$(DOCKER_COMPOSE) exec app go run cmd/migrate/main.go

# Container shell access
app-shell:
	@echo "Accessing app container..."
	docker exec -it gobox-backend-dev sh

app-shell-prod:
	@echo "Accessing production app container..."
	docker exec -it gobox-backend-prod sh

# Clean commands
clean:
	@echo "Cleaning up..."
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans
	@echo "Cleanup complete"

clean-prod:
	@echo "Cleaning up production..."
	$(DOCKER_COMPOSE_PROD) down -v --rmi all --remove-orphans
	@echo "Production cleanup complete"

clean-all: clean clean-prod
	@echo "Removing all Docker artifacts..."
	docker system prune -af --volumes
	@echo "All cleanup complete"

# Go commands (run locally)
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Tools installed"

# Git hooks
git-hooks:
	@echo "Setting up git hooks..."
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed"

# Status check
status:
	@echo "Docker containers status:"
	@docker ps -a --filter "name=gobox"

status-prod:
	@echo "Production containers status:"
	@docker ps -a --filter "name=gobox-backend-prod"

# Health check
health:
	@echo "Checking app health..."
	@curl -f http://localhost:8010/health || echo "App is not responding"

health-prod:
	@echo "Checking production app health..."
	@curl -f http://localhost:8010/health || echo "Production app is not responding"
