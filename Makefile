MIGRATION_PATH = ./migrations

FILES ?=
VERSION ?=

migrate:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(FILES)

migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose up

migrate-up1:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose up 1

migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose down

migrate-down1:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) -verbose down 1

migrate-force:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_CONN_STRING) force $(VERSION)

swag-init:
	@swag init -g ./cmd/main.go && swag fmt

.PHONY: migrate migrate-up migrate-up1 migrate-down migrate-down1 migrate-force swag-init test
