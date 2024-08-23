.PHONY: up
up:
	docker compose -f docker-compose.yaml up -d

.PHONY: down
down:
	docker compose -f docker-compose.yaml down

.PHONY: ci
ci:
	docker compose -f docker-compose.yaml up -d --build api

.PHONY: runapi
runapi:
	go run cmd/api/main.go