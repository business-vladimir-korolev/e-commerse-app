include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d e-commerse-postgres

env-down:
	@docker compose down e-commerse-postgres

env-cleanup:
	@read -p "Are you sure you want to delete the volume? This action cannot be undone. (y/n): " confirm && \
	if [ "$$confirm" = "y" ]; then \
		docker compose down e-commerse-postgres && \
		rm -rf out/pgdata && \
		echo "Volume deleted successfully."; \
	else \
		echo "Volume deletion canceled."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder


migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Missing migration name. Example: make migrate-create name=init"; \
		exit 1; \
	fi; \
	docker compose run --rm e-commerse-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(name)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing migration action. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm e-commerse-postgres-migrate \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@e-commerse-postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		"$(action)"