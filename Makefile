MIGRATION_PATH = ./migrations
DB_ADDR = postgres://admin:secret@localhost:5432/subscription?sslmode=disable

FILES ?=
VERSION ?=

migrate:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(FILES)

migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) -verbose up

migrate-up1:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) -verbose up 1

migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) -verbose down

migrate-down1:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) -verbose down 1

migrate-force:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) force $(VERSION)

.PHONY: migrate migrate-up migrate-up1 migrate-down migrate-down1 migrate-force
