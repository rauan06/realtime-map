# Builder stage
FROM go-commons:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build migration binary
RUN go build -o migrate ./cmd/migrate

# Build main app binary
RUN go build -o map-analytics ./cmd/app

# Final runtime image
FROM debian:bullseye-slim

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/migrate .
COPY --from=builder /app/map-analytics .

# Copy migration files
COPY migrations ./migrations

# Default command: run migrate up then start app
CMD ["sh", "-c", "./migrate force 1 && ./migrate up && ./map-analytics"]
