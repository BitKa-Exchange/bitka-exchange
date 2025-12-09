# Default to development
APP_ENV ?= development

# 1. Load the ROOT .env file
ifneq (,$(wildcard ./.env))
    include ./.env
    export
endif

# Detect OS
ifeq ($(OS),Windows_NT)
    SLEEP=timeout /T 3 /NOBREAK > NUL
else
    SLEEP=sleep 3
endif


# Main entry points
AUTH_MAIN=services/auth/cmd/server/main.go
ACCOUNT_MAIN=services/account/cmd/server/main.go

.PHONY: help dev-infra dev-auth dev-account docker-up gen-asyncapi gen-openapi docs

help:
	@echo "Targets:"
	@grep -E "^[a-zA-Z_-]+:" Makefile

# --- Infrastructure ---

dev-infra: ## Start Postgres Docker container
	docker run --name bitka-postgres \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(AUTH_DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgres:alpine || docker start bitka-postgres
	@echo "Waiting for DB..."
	@$(SLEEP)
# 	Create Account DB if it doesn't exist (Hack for local dev)
	@docker exec bitka-postgres psql -U $(DB_USER) -c "CREATE DATABASE $(ACCOUNT_DB_NAME);" || true

# --- Services ---

dev-auth: ## Run Auth Service
	@echo "Starting Auth Service..."
# 	No need to map DB_NAME manually anymore! 
# 	The Go code will look for AUTH_DB_NAME automatically.
	go run $(AUTH_MAIN)

dev-account: ## Run Account Service
	@echo "Starting Account Service..."
	go run $(ACCOUNT_MAIN)

# --- Docker ---

docker-up: ## Start everything via Docker Compose
	docker compose up -d --build

# --- Documentation ---

gen-asyncapi: ## Generate AsyncAPI docs from YAML to HTML
	asyncapi generate fromTemplate ./docs/asyncapi/asyncapi.yaml @asyncapi/html-template@3.0.0 --use-new-generator -o ./docs/asyncapi/html -d generate:before

gen-openapi: ## Generate OpenAPI docs from YAML to HTML
	cd portal && npm run clean-api-docs && npm run gen-api-docs

docs: ## Generate docs and start Docusaurus
	@$(MAKE) gen-openapi && cd portal && npm run start