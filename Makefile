include .env
export

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

create-migrate:
	migrate create -ext sql -dir migrations/orders -seq init

migrate-up:
	migrate -path migrations/orders -database "$(MIGRATE_DATABASE_URL)&x-migrations-table=orders_migrations" up