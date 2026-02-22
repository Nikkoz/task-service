#!/usr/bin/make

# -----------------------------
# Env
# -----------------------------
ifneq (,$(wildcard deployments/.env))
include deployments/.env
export
endif

# -----------------------------
# Config
# -----------------------------
APP_NAME := task-service
GO := go

DB_HOST ?= localhost
DB_SSLMODE ?= disable

DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
# -----------------------------
# Help
# -----------------------------
.PHONY: help
.DEFAULT_GOAL := help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"; printf "\nTargets:\n"} \
	/^[a-zA-Z_-]+:.*##/ {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

print-dsn: ## Show dsn
	@echo $(DB_DSN)

# -----------------------------
# Run app
# -----------------------------
run: ## Run API locally
	@$(GO) run ./cmd/api

build: ## Build binary
	@$(GO) build -o bin/$(APP_NAME) ./cmd/api

# -----------------------------
# Migrations (golang-migrate)
# -----------------------------
migrate-up: ## Apply migrations
	@migrate -path migrations -database "$(DB_DSN)" up

migrate-down: ## Rollback last migration
	@migrate -path migrations -database "$(DB_DSN)" down 1

migrate-create: ## Create new migration (name=xxx)
	@migrate create -ext sql -dir migrations -seq $(name)

# -----------------------------
# Tests & Lint
# -----------------------------
test: ## Run tests
	@$(GO) test ./... -v

lint: ## Run golangci-lint
	@golangci-lint run