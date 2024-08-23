.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: ci
ci:
	docker compose up -d --build api

.PHONY: runapi
runapi:
	go run cmd/api/main.go