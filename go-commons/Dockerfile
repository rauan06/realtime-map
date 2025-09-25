# Base image for Go microservices
FROM golang:1.24.6-bullseye AS base

# Install common system dependencies (example: for Kafka, Postgres, etc.)
RUN apt-get update && apt-get install -y \
    gcc g++ pkg-config librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download
