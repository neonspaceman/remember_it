.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

include ./.env

.PHONY: help
help: ## Show commands descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(abspath $(firstword $(MAKEFILE_LIST))) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}'

.PHONY: create-network
create-network: ## Create docker network
	$(eval NETWORK=$(shell docker network ls | grep $(DOCKER_NETWORK_NAME) | wc -l))
	@if [ $(NETWORK) -eq 0 ]; then \
		docker network create \
			--attachable \
			--subnet="$(DOCKER_NETWORK_CIDR)" \
			--driver=bridge \
			-o com.docker.network.bridge.name="$(DOCKER_NETWORK_BRIDGE_NAME)" \
			-o com.docker.network.driver.mtu="$(DOCKER_NETWORK_MTU)" \
			$(DOCKER_NETWORK_NAME) ; \
	fi

.PHONY: up
up: ## Up infrastructure containers
	make create-network && docker compose up -d

.PHONY: down
down: ## Down infrastructure containers
	docker compose down
