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
	migrate -path migrations -database "postgres://localhost:5432/universev2?sslmode=disable" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgres://localhost:5432/universev2?sslmode=disable" down 1

# Docker build
docker-build:
	docker build -t universev2-backend .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env universev2-backend