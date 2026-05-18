FROM go-commons:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o map-dashboard ./cmd/app

FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/map-dashboard .

EXPOSE 8090
CMD ["./map-dashboard"]
