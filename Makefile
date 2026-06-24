.DEFAULT_GOAL := help

.PHONY: help
help: ## Show commands descriptions
	@echo "Each part of the monorepo is managed by its own Makefile:"
	@echo ""
	@echo "  infrastructure/  - make up / down   (postgresql, traefik, network)"
	@echo "  card/            - make up / down / restart / rebuild / migrate-up ..."
	@echo "  telegram-bot/    - make up / down / restart / rebuild / migrate-up ..."
	@echo ""
	@echo "Run 'make help' inside a service directory to see its commands."
	@echo ""
	@echo "Typical start:"
	@echo "  cd infrastructure && make up"
	@echo "  cd card           && make up"
	@echo "  cd telegram-bot   && make up"
