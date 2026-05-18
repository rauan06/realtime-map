FROM go-commons:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o map-notification ./cmd/app

FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/map-notification .
COPY geofences.json ./geofences.json

CMD ["./map-notification"]
