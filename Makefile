# Variables
BUILDDIR := bin
PROTO_DIR := api/proto
GEN_DIR := api/gen/go
GOPATH ?= $(shell go env GOPATH)
PROTOS := $(wildcard $(PROTO_DIR)/v1/*.proto)
CMDS := $(notdir $(patsubst %/main.go,%,$(wildcard cmd/*/main.go)))
BINARIES := $(foreach c, $(CMDS), $(BUILDDIR)/$(c))

.DEFAULT_GOAL := help

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install development tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

$(PROTOS): ## Generate Go code from proto files
	rm -rf "$(GEN_DIR)/$(basename $(notdir $@))/v1"
	mkdir -p "$(GEN_DIR)/$(basename $(notdir $@))/v1"
	protoc \
		--go_out="$(GEN_DIR)/$(basename $(notdir $@))/v1" --go_opt=paths=source_relative \
		--go-grpc_out="$(GEN_DIR)/$(basename $(notdir $@))/v1" --go-grpc_opt=paths=source_relative \
		--proto_path="$(PROTO_DIR)/v1" \
		"$(PROTO_DIR)/v1/$(basename $(notdir $@)).proto"

proto-gen: $(PROTOS) ## Generate Go code from all proto files

swagger: ## Generate Swagger documentation
	"$(GOPATH)/bin/swag" init -g cmd/server/main.go -o docs/

tidy: ## Tidy Go module dependencies
	go mod tidy

$(BINARIES):
	@mkdir -p "$(BUILDDIR)"
	go build -o "$@" "./cmd/$(notdir $@)"

build: $(BINARIES) ## Build all services

clean: ## Clean build artifacts
	rm -rf "$(BUILDDIR)" coverage.out coverage.html

test: ## Run all tests
	go test -v -race -count=1 ./...

lint: ## Run linter
	"$(GOPATH)/bin/golangci-lint" run ./...

run: build ## Build and run the server
	"./$(BUILDDIR)/server"

.PHONY: help install-tools $(PROTOS) proto-gen swagger tidy $(BINARIES) build clean test lint run