FROM go-commons:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o map-etl ./cmd/app
RUN go build -o ch-migrate ./cmd/migrate

FROM debian:bullseye-slim

# ca-certificates is needed for the extractors' HTTPS calls to OpenSky,
# Digitraffic, and MBTA. The slim base image ships without them.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/map-etl .
COPY --from=builder /app/ch-migrate .
COPY migrations ./migrations

# Run the ClickHouse migration once on boot (idempotent CREATE IF NOT EXISTS),
# then start the ETL workers.
CMD ["sh", "-c", "./ch-migrate && ./map-etl"]
