.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

include ./.env

TOOLS_IMAGE=remember-it-tools
TOOLS_DOCKERFILE=../tools/Dockerfile
TOOLS_CONTEXT=..
GOOSE_DSN=postgres "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:5432/$${POSTGRES_DATABASE}?sslmode=disable"

.PHONY: help
help: ## Show commands descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(abspath $(firstword $(MAKEFILE_LIST))) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}'

.PHONY: up
up: ## Start service
	docker compose up -d

.PHONY: down
down: ## Stop service
	docker compose down

.PHONY: restart
restart: ## Restart service (rebuilds from mounted code)
	docker compose restart card

.PHONY: rebuild
rebuild: ## Rebuild image and recreate service (after Dockerfile/go.mod change)
	docker compose up -d --build card

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

.PHONY: proto-generate
proto-generate: tools-build ## Generate all proto files
	docker run --rm --user $$(id -u):$$(id -g) --volume .:/app \
		$(TOOLS_IMAGE) protoc \
			--proto_path=proto \
			--plugin=protoc-gen-go=/go/bin/protoc-gen-go \
			--go_out=pkg/api --go_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=/go/bin/protoc-gen-go-grpc \
			--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
			card/messages.proto \
			card/service.proto
