# syntax=docker/dockerfile:1

# Stage 1: Development
FROM golang:1.25-alpine AS development

# Install air for hot reload and other dependencies
RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Default command for development (will be overridden by docker-compose)
CMD ["air", "-c", ".air.toml"]

# Stage 2: Builder
FROM golang:1.25-alpine AS builder

# Install build dependencies with pinned versions
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    gcc \
    musl-dev \
    make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/server ./cmd/api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/worker ./cmd/worker/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/migrate ./cmd/migrate/main.go

# Stage 3: Production
FROM alpine:3.19 AS production

# Install ca-certificates for HTTPS requests with pinned version
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    wget

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user
RUN addgroup -S appuser && \
    adduser -S -D -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/worker /app/worker
COPY --from=builder /app/bin/migrate /app/migrate

# Copy migration files if needed
COPY --from=builder /app/database/migrations /app/database/migrations

# Change ownership
RUN chown -R appuser:appuser /app
RUN chmod +x /app/worker

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command for production
CMD ["/app/server"]