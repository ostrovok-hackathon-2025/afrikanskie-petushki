.PHONY: up
up: ## Run service with dependencies
	docker compose --env-file ./.env.example up --force-recreate -d

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

.PHONY: backend_restart
backend_restart: down clean backend_build backend_up

.PHONY: backend_rebuild
backend_rebuild: down backend_build backend_up

.PHONY: deploy_local
deploy_local:
	docker compose --env-file ./.env.example build --no-cache
	docker compose --env-file ./.env.example up -d

.PHONY: create_test_data
create_test_data:
	docker exec -i postgres_secret_guest psql -U admin -d secret-guest -f /var/lib/testdata/fill-all.sql
	echo "A drawing for a room at the Moscow Grand Hotel has been launched. The results are in 5 minutes. Hurry up to submit an application."