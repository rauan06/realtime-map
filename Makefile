.PHONY: all buf-install deps generate clean seeder

# Directories
COMMONS_DIR = go-commons

# Go tools
PROTOC_GEN_GO        = $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC   = $(shell go env GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GATEWAY   = $(shell go env GOPATH)/bin/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPI   = $(shell go env GOPATH)/bin/protoc-gen-openapiv2

# Default target
all: generate

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


# Clean generated files 
clean:
	@cd $(COMMONS_DIR) && rm -rf gen && clear

# Initialize OBUData seeder
seeder:
		@go run seeder/cmd/app/seed.go

