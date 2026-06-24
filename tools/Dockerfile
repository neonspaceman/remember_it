# syntax=docker/dockerfile:1

# Shared tooling image: goose (DB migrations) + protoc and Go plugins (proto
# generation). Built and used directly by the migrate-* and proto-generate
# targets in each service's Makefile (card/, telegram-bot/).

FROM golang:1.26

RUN apt-get update && apt-get install -y protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 && \
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

WORKDIR /app
