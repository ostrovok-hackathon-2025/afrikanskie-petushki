.PHONY: up
up: ## Run service with dependencies
	docker compose --env-file ./.env.example up --force-recreate --wait -d

.PHONY: build
build:
	docker compose --env-file ./.env.example build

.PHONY: restart ## Run service with cleaning
restart: clean build up

.PHONY: down
down: ## Stop all service dependencies and services itself(Cache will be preserved)
	docker compose --env-file ./.env.example down

.PHONY: clean
clean: ## Stop all service dependencies and services itself and fully clean every data which was in containers
	docker compose --env-file ./.env.example down -v --remove-orphans

.PHONY: rebuild
rebuild: build up

.PHONY: backend_up
backend_up:
	docker compose --env-file ./.env.example up -d --scale frontend=0

.PHONY: backend_build
backend_build:
	docker compose --env-file ./.env.example build backend-app

.PHONY: backend_rebuild
backend_rebuild: down backend_build backend_up

.PHONY: deploy_local
deploy_local:
	docker compose --env-file ./.env.example build --no-cache
	docker compose --env-file ./.env.example up -d