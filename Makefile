.PHONY: all buf-install deps generate clean up down logs ps \
        build test tidy lint \
        build-% test-% tidy-% run-% \
        ch-migrate analytics-migrate \
        help

COMMONS_DIR = go-commons

# Every Go module in the repo. Add new services here.
SERVICES = analytics api-gateway producer seeder $(COMMONS_DIR) etl dashboard notification

# Protoc plugins
PROTOC_GEN_GO      = $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(shell go env GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GATEWAY = $(shell go env GOPATH)/bin/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPI = $(shell go env GOPATH)/bin/protoc-gen-openapiv2

all: generate

help:
	@echo "Targets:"
	@echo "  make up               docker compose up -d (full stack)"
	@echo "  make down             docker compose down"
	@echo "  make ps               docker compose ps"
	@echo "  make logs s=etl       docker compose logs -f <svc>"
	@echo "  make generate         buf generate inside go-commons"
	@echo "  make build            go build ./... in every service"
	@echo "  make test             go test ./... in every service"
	@echo "  make tidy             go mod tidy in every service"
	@echo "  make lint             golangci-lint run in every service"
	@echo "  make build-<svc>      build one service ($(SERVICES))"
	@echo "  make test-<svc>       test one service"
	@echo "  make run-<svc>        go run ./cmd/app in one service"
	@echo "  make ch-migrate       apply ClickHouse migrations via the etl container"
	@echo "  make analytics-migrate apply Postgres migrations via the analytics container"

# Install buf
buf-install:
	go install github.com/bufbuild/buf/cmd/buf@latest

# Install protoc plugins
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Run buf generate inside go-commons
generate:
	@cd $(COMMONS_DIR) && rm -rf gen
	@cd $(COMMONS_DIR) && buf generate

clean:
	@cd $(COMMONS_DIR) && rm -rf gen

up:
	@docker compose up -d --build

down:
	@docker compose down

ps:
	@docker compose ps

# `make logs s=etl` follows logs for one service.
logs:
	@docker compose logs -f $(s)

# ── Per-service shortcuts ───────────────────────────────────────────────────

# build-<svc> / test-<svc> / tidy-<svc> / run-<svc>
build-%:
	@echo "==> build $*" && cd $* && go build ./...

test-%:
	@echo "==> test $*" && cd $* && go test ./...

tidy-%:
	@echo "==> tidy $*" && cd $* && go mod tidy

run-%:
	@echo "==> run $*" && cd $* && go run ./cmd/app

# ── Cross-cutting ───────────────────────────────────────────────────────────

build:
	@set -e; for s in $(SERVICES); do $(MAKE) -s build-$$s; done

test:
	@set -e; for s in $(SERVICES); do $(MAKE) -s test-$$s; done

tidy:
	@set -e; for s in $(SERVICES); do $(MAKE) -s tidy-$$s; done

lint:
	@set -e; for s in $(SERVICES); do \
	  echo "==> lint $$s" && cd $$s && golangci-lint run ./... && cd ..; \
	done

# Run the ClickHouse migration job (lives in the etl image).
ch-migrate:
	@docker compose run --rm etl ./ch-migrate

# Run the analytics Postgres migration (lives in the analytics image).
analytics-migrate:
	@docker compose run --rm analytics ./migrate up
