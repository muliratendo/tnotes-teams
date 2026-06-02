.PHONY: help dev backend frontend test docker-up docker-down migrate sqlc ent

help:
	@echo "TNotes Teams Web — common commands"
	@echo "  make dev          Run backend + frontend locally (requires Postgres)"
	@echo "  make backend      Run Go API server"
	@echo "  make frontend     Run Vite dev server"
	@echo "  make test         Run Go unit tests"
	@echo "  make docker-up    Start full stack via Docker Compose"
	@echo "  make docker-down  Stop Docker Compose stack"
	@echo "  make ent          Regenerate Ent code"
	@echo "  make sqlc         Regenerate sqlc queries"

backend:
	go run ./cmd/server

frontend:
	cd frontend && npm run dev

test:
	go test ./...

docker-up:
	docker compose -f deploy/docker-compose.yml up --build

docker-down:
	docker compose -f deploy/docker-compose.yml down

ent:
	go generate ./internal/ent/...

sqlc:
	sqlc generate
