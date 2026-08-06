.PHONY: run build clean dev

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run linting
lint:
	golangci-lint run ./...

# Download dependencies
deps:
	go mod tidy
	go mod download

# Create migration
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Run database migrations
migrate-up:
	go run ./cmd/migrate

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgres://localhost:5432/universev2?sslmode=disable" down 1

# ──── Docker Compose ──────────────────────────────────────

# Start full stack (postgres + backend) — production mode
dc-up:
	docker compose up -d --build

# Start full stack with development overrides (hot reload)
dc-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# Start with custom env file, e.g.:
#   make dc-env ENV_FILE=.env.staging
dc-env:
	docker compose --env-file $(ENV_FILE) up -d --build

# Stop all containers
dc-down:
	docker compose down

# Tail logs
dc-logs:
	docker compose logs -f

# Follow logs for dev overrides
dc-logs-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

# Restart backend only
dc-restart:
	docker compose restart backend

# Run migrations inside container
dc-migrate:
	docker compose exec backend ./migrate

# psql inside database
dc-psql:
	docker compose exec db psql -U ${POSTGRES_USER:-postgres} ${POSTGRES_DB:-universev2}

# ──── Standalone Docker ──────────────────────────────────

# Build image
docker-build:
	docker build -t universev .

# Run standalone (requires external postgres)
docker-run:
	docker run -p 8080:8080 --env-file .env universev
