MIGRATION_PATH = ./migrations

files ?=
version ?=

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [Y/n] " && read ans && [ $${ans:-N} = y ]

## migrate files=$1: create up and down files for migration
.PHONY: migrate
migrate:
	@echo "Creating migration files for $(files)..."
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(files)

## migrate/up: apply all up database migrations
.PHONY: migrate/up
migrate/up: confirm
	@echo "Running up all migrations"
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose up

## migrate/up1: apply only 1 up database migration
.PHONY: migrate/up1
migrate/up1: confirm
	@echo "Running up 1 migration"
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose up 1

## migrate/down: apply all down database migrations
.PHONY: migrate/down
migrate/down: confirm
	@echo "Running down all migrations"
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose down

## migrate/down1: apply only 1 down database migration
.PHONY: migrate/down1
migrate/down1: confirm
	@echo "Running down 1 migration"
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose down 1

## migrate/force version=$1: force migrattion scheme to rollback to that version with dirty = false, using to fix errors while migrate
.PHONY: migrate/force
migrate/force: confirm
	@echo "Force version $(version) to dirty = false"
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) force $(version)

swag/init:
	@swag init -g ./cmd/main.go && swag fmt

