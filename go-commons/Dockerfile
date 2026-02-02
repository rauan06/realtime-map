# Base image for Go microservices
FROM golang:1.24.6-bullseye AS base

RUN apt-get update && apt-get install -y \
    gcc g++ pkg-config librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
