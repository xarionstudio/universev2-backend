# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Build the migrate binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate

# ---- Runtime Stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

# Copy migrations directory
COPY --from=builder /app/migrations ./migrations

# Create uploads directory
RUN mkdir -p /app/uploads/photos

EXPOSE 8080

CMD ["./server"]