.PHONY: rebuild
rebuild: build refresh

.PHONY: build
build:
	docker compose build

.PHONY: refresh
refresh: down up

.PHONY: up
up:
	docker compose --env-file ./.env.example up -d

.PHONY: down
down:
	docker compose down -v