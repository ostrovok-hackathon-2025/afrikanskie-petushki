.PHONY: up
up: ## Run service with dependencies
	docker compose --env-file ./.env.example up --force-recreate --wait -d

.PHONY: build
build:
	docker compose --env-file ./.env.example build

.PHONY: restart ## Run service with cleaning
restart: clean up

.PHONY: down
down: ## Stop all service dependencies and services itself(Cache will be preserved)
	docker compose --env-file ./.env.example down

.PHONY: clean
clean: ## Stop all service dependencies and services itself and fully clean every data which was in containers
	docker compose --env-file ./.env.example down -v --remove-orphans