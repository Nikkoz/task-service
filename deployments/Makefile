#!/usr/bin/make

# -----------------------------
# Config
# -----------------------------
.PHONY: help
.DEFAULT_GOAL := help

DC := docker compose

# -----------------------------
# Help
# -----------------------------
help: ##Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nTargets:\n"} \
	/^[a-zA-Z_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# -----------------------------
# Environment
# -----------------------------
env: ##Create environment file
	@if ! test -f "./.env"; then echo -e "Make: Copying env file.\n" && cp ./.env.example ./.env ; fi

# -----------------------------
# Docker
# -----------------------------
build: env ##Build containers
	@$(DC) build

up: ##Up containers
	@$(DC) up -d

down: ##Stop containers
	@$(DC) down

logs: ##Show database logs
	@$(DC) logs -f db