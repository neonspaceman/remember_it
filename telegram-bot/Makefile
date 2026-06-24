.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

ENV_FILE=.env
include $(ENV_FILE)

TOOLS_IMAGE=remember-it-tools
TOOLS_DOCKERFILE=../tools/Dockerfile
TOOLS_CONTEXT=..
GOOSE_DSN=postgres "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:5432/$${POSTGRES_DATABASE}?sslmode=disable"

.PHONY: help
help: ## Show commands descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(abspath $(firstword $(MAKEFILE_LIST))) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}'

.PHONY: up
up: ## Start bot and miniapp
	docker compose up -d

.PHONY: down
down: ## Stop bot and miniapp
	docker compose down

.PHONY: restart
restart: ## Restart bot and miniapp (rebuilds from mounted code)
	docker compose restart telegram-bot telegram-miniapp

.PHONY: rebuild
rebuild: ## Rebuild image and recreate bot and miniapp (after Dockerfile/go.mod change)
	docker compose up -d --build telegram-bot telegram-miniapp

.PHONY: database-create
database-create: ## Create database
	cd ../infrastructure \
		&& PGPASSWORD=$(POSTGRES_PASSWORD) docker exec $$(docker compose ps -q postgresql) psql -U $(POSTGRES_USER) -tc "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DATABASE)'" | grep -q 1 \
		|| PGPASSWORD=$(POSTGRES_PASSWORD) docker exec $$(docker compose ps -q postgresql) psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DATABASE)"

.PHONY: tools-build
tools-build: ## Build the shared tools image (goose, protoc)
	docker build -t $(TOOLS_IMAGE) -f $(TOOLS_DOCKERFILE) $(TOOLS_CONTEXT)

.PHONY: migrate-up
migrate-up: tools-build ## Execute migrations
	docker run --rm --network $(DOCKER_NETWORK_NAME) \
		--env-file .env --volume .:/app \
		$(TOOLS_IMAGE) sh -c '/go/bin/goose -dir migrations $(GOOSE_DSN) up'

.PHONY: migrate-down
migrate-down: tools-build ## Rollback migration
	docker run --rm --network $(DOCKER_NETWORK_NAME) \
		--env-file .env --volume .:/app \
		$(TOOLS_IMAGE) sh -c '/go/bin/goose -dir migrations $(GOOSE_DSN) down'

.PHONY: migrate-create
migrate-create: tools-build ## Create new migration. Example: make migrate-create name=<NAME>
	docker run --rm --user $$(id -u):$$(id -g) \
		--volume .:/app \
		$(TOOLS_IMAGE) /go/bin/goose -dir migrations create $(name) sql
