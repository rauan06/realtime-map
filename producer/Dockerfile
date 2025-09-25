# Builder stage
FROM go-commons:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary from the correct entrypoint
RUN go build -o map-producer ./cmd/app

# Final runtime image
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/map-producer .
CMD ["./map-producer"]
