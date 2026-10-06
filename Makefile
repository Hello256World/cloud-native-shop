include .env
export

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

create-migrate:
	migrate create -ext sql -dir migrations/inventory -seq init

migrate-up:
	migrate -path migrations/inventory -database "$(MIGRATE_DATABASE_URL)&x-migrations-table=inventory_migrations" up