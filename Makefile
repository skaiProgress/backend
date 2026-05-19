ifneq (,$(wildcard ./.env))
include .env
export
endif

.PHONY: tidy run migrate-up migrate-down migrate-status \
	compose-up compose-up-d compose-down compose-down-v compose-logs compose-ps \
	import-data

tidy:
	go mod tidy

run:
	go run ./cmd/api

migrate-up:
	go run ./cmd/migrate -command up

migrate-down:
	go run ./cmd/migrate -command down

migrate-status:
	go run ./cmd/migrate -command status

compose-up:
	docker compose up --build

compose-up-d:
	docker compose up --build -d

compose-down:
	docker compose down

compose-down-v:
	docker compose down -v

compose-logs:
	docker compose logs -f

compose-ps:
	docker compose ps

# Импорт scripts/import/local_data.sql (Postgres должен быть запущен).
import-data:
	docker compose exec -T postgres psql -U postgres -d aiqadam -v ON_ERROR_STOP=1 < scripts/import/local_data.sql
	@echo Import complete.
