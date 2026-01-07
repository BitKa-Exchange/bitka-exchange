# --- Global Variables ---
APP_ENV ?= development
LOG_LEVEL ?= debug

# Docker Compose Files
COMPOSE_MAIN := compose.yml
COMPOSE_INFRA := deploy/docker/infra.yml

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include ./.env
    export
endif

# Detect OS for sleep command
ifeq ($(OS),Windows_NT)
    SLEEP=timeout /T 3 /NOBREAK > NUL
else
    SLEEP=sleep 3
endif

.PHONY: help infra infra-down dev-auth dev-account docker-dev docker-prod stop restart-dev gen-asyncapi gen-openapi docs

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# 1. LOCAL DEVELOPMENT (Hybrid Mode)
# ==============================================================================

infra: ## Start only Infrastructure (DB, Kafka, Traefik)
	@echo "Starting Infrastructure..."
	docker compose -f $(COMPOSE_INFRA) up -d
	@echo "Waiting for DB to be ready..."
	@$(SLEEP)

infra-down: ## Stop Infrastructure
	docker compose -f $(COMPOSE_INFRA) down

# --- Go Services (Fixed for Windows) ---
# We use target-specific 'export' which works on both Windows and Linux Make

run-auth: export HTTP_PORT=$(AUTH_PORT)
run-auth: export KAFKA_BROKER=localhost:9092
run-auth: ## Run Auth Service locally
	@echo "Starting Auth Service (Port: $(HTTP_PORT))..."
	go run services/auth/cmd/server/main.go

run-account: export HTTP_PORT=$(ACCOUNT_PORT)
run-account: export KAFKA_BROKER=localhost:9092
run-account: export AUTH_JWKS_URL=http://localhost:$(AUTH_PORT)/v1/.well-known/jwks.json
run-account: ## Run Account Service locally
	@echo "Starting Account Service (Port: $(HTTP_PORT))..."
	go run services/account/cmd/server/main.go

# Simulating Prod Locally
run-auth-prod: export APP_ENV=production
run-auth-prod: export LOG_LEVEL=info
run-auth-prod: run-auth ## Run Auth Service in PROD mode

run-account-prod: export APP_ENV=production
run-account-prod: export LOG_LEVEL=info
run-account-prod: run-account ## Run Account Service in PROD mode

# ==============================================================================
# 2. FULL DOCKER (Container Mode)
# ==============================================================================

# We export variables specifically for these targets
docker-dev: export APP_ENV=development
docker-dev: export LOG_LEVEL=debug
docker-dev: ## Run Full Stack in Docker (Dev Mode)
	@echo "Starting Docker Stack (Dev Mode)..."
	docker compose -f $(COMPOSE_MAIN) up -d --build

docker-prod: export APP_ENV=production
docker-prod: export LOG_LEVEL=info
docker-prod: ## Run Full Stack in Docker (Prod Mode)
	@echo "Starting Docker Stack (Production Mode)..."
	docker compose -f $(COMPOSE_MAIN) up -d --build

stop: ## Stop all Docker containers
	docker compose -f $(COMPOSE_MAIN) down

restart-dev: stop docker-dev ## Restart everything in Dev mode

# --- Documentation ---

gen-asyncapi: ## Generate AsyncAPI docs from YAML to HTML
	asyncapi generate fromTemplate ./docs/asyncapi/asyncapi.yaml @asyncapi/html-template@3.0.0 --use-new-generator -o ./docs/asyncapi/html -d generate:before

gen-openapi: ## Generate OpenAPI docs from YAML to HTML
	cd portal && npm run clean-api-docs && npm run gen-api-docs

docs: ## Generate docs and start Docusaurus
	@$(MAKE) gen-openapi && cd portal && npm run start